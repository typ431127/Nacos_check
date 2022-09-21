package app

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"nacos_check/nacos"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type FileConfig struct {
	Url               string   `toml:url`
	Container_network []string `toml:container_network`
	Label             []map[string]string
}
type FileConfigLable struct {
	Name  string
	Value string
}

func GetConfigFilePath() string {
	homedir, err := HomeDir()
	if err != nil {
		fmt.Println("获取系统家目录获取异常")
		homedir = "."
	}
	configfile := filepath.Join(homedir, ".nacos_conf.toml")
	return configfile
}

func IPFilePathLoad() {
	if _, err := os.Stat(nacos.Ipfile); err != nil {
		if !os.IsExist(err) {
			nacos.Ipparse = false
			return
		}
	} else {
		nacos.Ipparse = true
		file, err := os.OpenFile(nacos.Ipfile, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Println("打开文件错误")
			os.Exit(nacos.Exitcode)
		}
		defer file.Close()
		fileb, _ := ioutil.ReadAll(file)
		if err := json.Unmarshal(fileb, &nacos.Ipdata); err != nil {
			fmt.Println("ip文件解析错误,请确认json格式")
			nacos.Ipparse = false
		}
	}
}

func LoadConfig() {
	IPFilePathLoad()
	var config FileConfig
	configfile := GetConfigFilePath()
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("配置文件错误格式错误", configfile)
			fmt.Println(err)
			os.Exit(2)
		}
	}()
	if _, err := os.Stat(configfile); err != nil {
		if !os.IsExist(err) {
			return
		}
	} else {
		_, err := toml.DecodeFile(configfile, &config)
		for _, lable := range config.Label {
			nacos.AddLable[lable["name"]] = lable["value"]
		}
		if err != nil {
			fmt.Println("配置文件错误格式错误", configfile)
			return
		}
		if nacos.Nacosurl == "http://dev-k8s-nacos:8848" {
			nacos.Nacosurl = config.Url
		}
		if len(config.Container_network) != 0 {
			nacos.MaxCidrBlocks = config.Container_network
		}
	}
	nacos.ContainerdInit()
}

func NacosInit() {
	defer func() {
		if func_err := recover(); func_err != nil {
			fmt.Println("程序初始化错误.....")
			fmt.Println(func_err)
			os.Exit(2)
		}
	}()
	LoadConfig()
	u, err := url.Parse(nacos.Nacosurl)
	if err != nil {
		fmt.Println("url解析错误!")
		os.Exit(nacos.Exitcode)
	}
	nacos.Na = &nacos.Nacos{}
	nacos.Na.Client.Timeout = time.Second * 15
	nacos.Na.DefaultUlr = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	nacos.Na.Host = u.Host
	nacos.Na.Scheme = u.Scheme
	nacos.Na.GetNameSpace()
}

func NacosRun() {
	nacos.Na.GetNacosInstance()
	switch {
	case nacos.ExportJson:
		jsondata, err := nacos.Na.GetJson("byte")
		if err != nil {
			fmt.Println("获取json发生错误")
			os.Exit(2)
		}
		data := make([]byte, 0)
		var check bool
		if data, check = jsondata.([]byte); !check {
			fmt.Println("转换失败")
			os.Exit(2)
		}
		fmt.Println(string(data))
		os.Exit(0)
	case nacos.Writefile != "":
		nacos.Na.WriteFile()
	default:
		fmt.Println("Nacos:", nacos.Nacosurl)
		if nacos.Watch {
			fmt.Printf("监控模式 刷新时间:%s/次\n", nacos.Second)
			for {
				nacos.Na.GetNacosInstance()
				nacos.Na.TableRender()
				time.Sleep(nacos.Second)
			}
		}
		nacos.Na.TableRender()
	}
	nacos.Na.Client.CloseIdleConnections()
}
