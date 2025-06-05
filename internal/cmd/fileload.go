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

func nacosFilePathLoad() {
	type Config struct {
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
		var newConfig Config
		_, err := toml.DecodeFile(configfilepath, &newConfig)
		for _, label := range newConfig.Label {
			nacos.ADDLABEL[label["name"]] = label["value"]
		}
		if err != nil {
			fmt.Println("配置文件错误格式错误", configfilepath)
			return
		}
		if len(nacos.USERNAME) == 0 {
			nacos.USERNAME = newConfig.Username
		}
		if len(nacos.PASSWORD) == 0 {
			nacos.PASSWORD = newConfig.Password
		}
		nacos.IPFILE = newConfig.Ipfile
		if nacos.NACOSURL == "http://dev-k8s-nacos:8848" {
			nacos.NACOSURL = newConfig.Url
		}
		if len(newConfig.Container_network) != 0 {
			pkg.MaxCidrBlocks = newConfig.Container_network
		}
		for _, namespace := range newConfig.Namespace {
			nacos.NAMESPACELIST = append(nacos.NAMESPACELIST, nacos.NamespaceServer{
				Namespace:         namespace,
				NamespaceShowName: namespace,
				Quota:             200,
				Type:              2,
			})
		}
		for _, group := range newConfig.Group {
			if !pkg.InString(group, nacos.GROUPLIST) {
				nacos.GROUPLIST = append(nacos.GROUPLIST, group)
			}
		}
		if newConfig.Nacos_sync_contextpath == "" {
			newConfig.Nacos_sync_contextpath = "/nacos"
		}
		nacos.FileConfig.ContextPath = newConfig.Nacos_sync_contextpath
		nacos.FileConfig.Sync = newConfig.Nacos_sync
	}
}

func ipconfigLoad() {
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
