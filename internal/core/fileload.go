package core

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io"
	"nacos-check/internal/config"
	"nacos-check/pkg"
	"os"
	"path/filepath"
)

func GetConfigFilePath() string {
	homedir, err := pkg.HomeDir()
	if err != nil {
		fmt.Println("获取系统家目录获取异常")
		homedir = "."
	}
	configfile := filepath.Join(homedir, ".nacos_conf.toml")
	return configfile
}
func NacosFilePathLoad() {
	type NewConfig struct {
		Url              string   `toml:"url"`
		Username         string   `toml:"username"`
		Password         string   `toml:"password"`
		Namespace        []string `toml:"namespace"`
		Group            []string `toml:"group"`
		ContainerNetwork []string `toml:"container_network"`
		Label            []map[string]string
		Ipfile           string `toml:"ipfile"`
	}
	homedir, err := pkg.HomeDir()
	if err != nil {
		homedir = "."
	}
	configfile := filepath.Join(homedir, ".nacos_conf.toml")
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("配置文件错误格式错误", configfile, err)
			os.Exit(2)
		}
	}()
	if _, err := os.Stat(configfile); err != nil {
		if !os.IsExist(err) {
			return
		}
	} else {
		var newConfig NewConfig
		_, err := toml.DecodeFile(configfile, &newConfig)
		for _, label := range newConfig.Label {
			config.ADDLABEL[label["name"]] = label["value"]
		}
		if err != nil {
			fmt.Println("配置文件错误格式错误", configfile)
			return
		}
		if len(config.USERNAME) == 0 {
			config.USERNAME = newConfig.Username
		}
		if len(config.PASSWORD) == 0 {
			config.PASSWORD = newConfig.Password
		}
		config.IPFILE = newConfig.Ipfile
		if config.NACOSURL == "http://dev-k8s-nacos:8848" {
			config.NACOSURL = newConfig.Url
		}
		if len(newConfig.ContainerNetwork) != 0 {
			pkg.MaxCidrBlocks = newConfig.ContainerNetwork
		}
		for _, namespace := range newConfig.Namespace {
			config.NAMESPACELIST = append(config.NAMESPACELIST, config.NamespaceServer{
				Namespace:         namespace,
				NamespaceShowName: namespace,
				Quota:             200,
				Type:              2,
			})
		}
		for _, group := range newConfig.Group {
			if !pkg.InString(group, config.GROUPLIST) {
				config.GROUPLIST = append(config.GROUPLIST, group)
			}
		}
	}
}

func IPFilePathLoad() {
	if _, err := os.Stat(config.IPFILE); err != nil {
		if !os.IsExist(err) {
			config.PARSEIP = false
			return
		}
	} else {
		config.PARSEIP = true
		file, err := os.OpenFile(config.IPFILE, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println("打开文件错误")
			os.Exit(config.EXITCODE)
		}
		defer file.Close()
		bytefile, _ := io.ReadAll(file)
		if err := json.Unmarshal(bytefile, &config.IPDATA); err != nil {
			fmt.Println("ip文件解析错误,请确认json格式")
			config.PARSEIP = false
		}
	}
}
