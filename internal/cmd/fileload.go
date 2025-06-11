package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"nacos-check/internal/nacos"
	"nacos-check/pkg"
	"nacos-check/pkg/fmtd"
	"os"
	"path/filepath"
)

func getConfigFilePath() string {
	homedir, err := pkg.HomeDir()
	if err != nil {
		fmt.Println("获取系统家目录获取异常")
		homedir = "."
	}
	configfile := filepath.Join(homedir, ".nacos_conf.toml")
	return configfile
}

func LoadNacosConfig() {
	type TomlConfig struct {
		Url                    string   `toml:"url"`
		Username               string   `toml:"username"`
		Password               string   `toml:"password"`
		Namespace              []string `toml:"namespace"`
		Group                  []string `toml:"group"`
		Container_network      []string `toml:"container_network"`
		Label                  []map[string]string
		Nacos_sync             []map[string]string `toml:"nacos_sync"`
		Nacos_sync_contextpath string              `toml:"nacos_sync_contextpath"`
		Ipfile                 string              `toml:"ipfile"`
		NetworkFile            string              `toml:"networkfile"`
	}
	homedir, err := pkg.HomeDir()
	if err != nil {
		homedir = "."
	}
	configfilepath := ""
	if nacos.FileConfig.ConfigFile == "" {
		configfilepath = filepath.Join(homedir, ".nacos_conf.toml")
	} else {
		configfilepath = nacos.FileConfig.ConfigFile
	}
	defer func() {
		if err := recover(); err != nil {
			fmtd.Fatalln("配置文件错误格式错误", configfilepath, err)
			os.Exit(2)
		}
	}()
	if _, err := os.Stat(configfilepath); err != nil {
		if !os.IsExist(err) {
			return
		}
	} else {
		var tomlConfig TomlConfig
		_, err := toml.DecodeFile(configfilepath, &tomlConfig)
		for _, label := range tomlConfig.Label {
			nacos.ADDLABEL[label["name"]] = label["value"]
		}
		if err != nil {
			fmt.Println("配置文件错误格式错误", configfilepath)
			return
		}
		if len(nacos.USERNAME) == 0 {
			nacos.USERNAME = tomlConfig.Username
		}
		if len(nacos.PASSWORD) == 0 {
			nacos.PASSWORD = tomlConfig.Password
		}
		nacos.IPFILE = tomlConfig.Ipfile
		nacos.NETWORKFILE = tomlConfig.NetworkFile
		if nacos.NACOSURL == "http://dev-k8s-nacos:8848" {
			nacos.NACOSURL = tomlConfig.Url
		}
		if len(tomlConfig.Container_network) != 0 {
			pkg.MaxCidrBlocks = tomlConfig.Container_network
		}
		for _, namespace := range tomlConfig.Namespace {
			nacos.NAMESPACELIST = append(nacos.NAMESPACELIST, nacos.NamespaceServer{
				Namespace:         namespace,
				NamespaceShowName: namespace,
				Quota:             200,
				Type:              2,
			})
		}
		for _, group := range tomlConfig.Group {
			if !pkg.InString(group, nacos.GROUPLIST) {
				nacos.GROUPLIST = append(nacos.GROUPLIST, group)
			}
		}
		if tomlConfig.Nacos_sync_contextpath == "" {
			tomlConfig.Nacos_sync_contextpath = "/nacos"
		}
		nacos.FileConfig.ContextPath = tomlConfig.Nacos_sync_contextpath
		nacos.FileConfig.Sync = tomlConfig.Nacos_sync
	}
}

// LoadIPConfig 加载ip配置文件
func LoadIPConfig() {
	if _, err := os.Stat(nacos.IPFILE); err != nil {
		if !os.IsExist(err) {
			nacos.PARSEIP = false
			return
		}
	} else {
		nacos.PARSEIP = true
		file, err := os.OpenFile(nacos.IPFILE, os.O_RDONLY, 0644)
		if err != nil {
			fmtd.Fatalln("打开文件错误")
		}
		defer file.Close()
		data, _ := io.ReadAll(file)
		if err := json.Unmarshal(data, &nacos.IPDATA); err != nil {
			fmt.Println("ip文件解析错误,请确认json格式")
			nacos.PARSEIP = false
		}
	}
}

func LoadNetworkConfig() error {
	if _, err := os.Stat(nacos.NETWORKFILE); err != nil {
		if !os.IsExist(err) {
			nacos.PARSENET = false
			return fmt.Errorf("未找到网络配置文件")
		}
	} else {
		var filedata map[string][]string
		nacos.NETDATA = make(map[string]string)
		nacos.PARSENET = true
		file, err := os.OpenFile(nacos.NETWORKFILE, os.O_RDONLY, 0644)
		if err != nil {
			return fmt.Errorf("打开文件错误")
		}
		defer file.Close()
		data, _ := io.ReadAll(file)
		if err := json.Unmarshal(data, &filedata); err != nil {
			fmt.Println("network文件解析错误,请确认json格式")
			nacos.PARSENET = false
			return fmt.Errorf("network文件解析错误,请确认json格式")
		}
		for idc, cidr := range filedata {
			for _, cidr := range cidr {
				nacos.NETDATA[cidr] = idc
				nacos.NETCIDR = append(nacos.NETCIDR, cidr)
			}
		}
	}
	return nil
}
