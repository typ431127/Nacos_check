package app

import (
	"encoding/json"
	"fmt"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"nacos_check/nacos"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func FilePathCheck() {
	if _, err := os.Stat(nacos.Ipfile); err != nil {
		if !os.IsExist(err) {
			nacos.Ipparse = false
			return
		}
	} else {
		nacos.Ipparse = true
	}
}

func NacosConfigCheck() {
	homedir, err := HomeDir()
	if err != nil {
		fmt.Println("获取系统家目录获取异常")
		homedir = "."
	}
	configfile := filepath.Join(homedir, ".nacos_url")
	if _, err := os.Stat(configfile); err != nil {
		if !os.IsExist(err) {
			return
		}
	} else {
		cfg, err := ini.Load(configfile)
		if err != nil {
			fmt.Println("读取配置文件发生错误", err)
			return
		}
		if nacos.Nacosurl == "http://dev-k8s-nacos:8848" {
			nacos.Nacosurl = cfg.Section("").Key("url").String()
		}
	}
}
func IP_Parse() {
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
		//os.Exit(nacos.Exitcode)
	}
}

func FlagCheck() {
	FilePathCheck()
	nacos.ContainerdInit()
	if nacos.Ipparse {
		IP_Parse()
	}
}

func NacosInit() {
	defer func() {
		if func_err := recover(); func_err != nil {
			fmt.Println("程序初始化错误.....")
			fmt.Println(func_err)
			os.Exit(2)
		}
	}()
	FlagCheck()
	NacosConfigCheck()
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
