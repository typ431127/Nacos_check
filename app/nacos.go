package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"nacos_check/nacos"
	"nacos_check/web"
	"net/url"
	"os"
	"strconv"
	"time"
)

func init() {
	flag.StringVar(&nacos.Nacosurl, "url", "http://dev-k8s-nacos:8848", "nacos地址")
	flag.StringVar(&nacos.Writefile, "write", "", "prometheus 自动发现文件路径")
	flag.StringVar(&nacos.Ipfile, "ipfile", "salt_ip.json", "ip解析文件")
	flag.StringVar(&nacos.Findstr, "find", "", "查找服务")
	flag.BoolVar(&nacos.Noconsole, "noconsole", false, "不输出console")
	flag.BoolVar(&nacos.Export_json, "json", false, "输出json")
	flag.BoolVar(&nacos.Web, "web", false, "开启Web api Prometheus http_sd_configs")
	flag.StringVar(&nacos.Port, "port", ":8099", "web 端口")
	flag.BoolVar(&nacos.Cluster_status, "cluster", false, "查看集群状态")
	flag.BoolVar(&nacos.V2upgrade, "v2upgrade", false, "查看2.0升级状态,和-cluster一起使用")
	flag.BoolVar(&nacos.Version, "version", false, "查看版本")
	flag.BoolVar(&nacos.Watch, "watch", false, "监控服务")
	flag.DurationVar(&nacos.Second, "second", 2*time.Second, "监控服务间隔刷新时间")
}

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
	if nacos.Version {
		fmt.Println("版本:0.4.2")
		os.Exit(0)
	}
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
	flag.Parse()
	FlagCheck()
	u, err := url.Parse(nacos.Nacosurl)
	if err != nil {
		fmt.Println("url解析错误!")
		os.Exit(nacos.Exitcode)
	}
	nacos.Na = nacos.Nacos{}
	nacos.Na.Client.Timeout = time.Second * 15
	nacos.Na.DefaultUlr = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	nacos.Na.Host = u.Host
	nacos.Na.Scheme = u.Scheme
	nacos.Na.GetNacosInstance()
}

func NacosRun() {
	if nacos.Web {
		web.Runwebserver()
	}
	if nacos.Export_json {
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
	}
	if !nacos.Noconsole && !nacos.Cluster_status {
		fmt.Println("Nacos:", nacos.Nacosurl)
		if nacos.Watch {
			for {
				nacos.Na.GetNacosInstance()
				nacos.Na.PrintInfo()
				time.Sleep(nacos.Second)
			}
		}
		nacos.Na.PrintInfo()
	}
	if nacos.Cluster_status {
		tablecluser := tablewriter.NewWriter(os.Stdout)
		tablecluser.SetHeader([]string{"节点", "端口", "状态", "版本", "刷新时间", "健康实例", "异常实例"})
		for _, key := range nacos.Na.Clusterdata {
			tablecluser.Append([]string{
				key.Ip,
				key.Port,
				key.State,
				key.Version,
				key.LastRefreshTime,
				strconv.Itoa(len(key.HealthInstance)),
				strconv.Itoa(len(key.UnHealthInstance)),
			})
		}
		leader := gjson.Get(nacos.Na.Cluster, "servers.0.extendInfo.raftMetaData.metaDataMap.naming_instance_metadata.leader")
		fmt.Printf("Nacos集群状态: (数量:%d)\n集群Master: %s\n", tablecluser.NumLines(), leader)
		tablecluser.Render()
	}
	if nacos.V2upgrade && nacos.Cluster_status {
		tablecluser := tablewriter.NewWriter(os.Stdout)
		tablecluser.SetHeader([]string{"节点", "双写", "v2服务", "v2实例", "v1服务", "v1实例", "upgraded", "全部V2"})
		for _, key := range nacos.Na.Clusterdata {
			tablecluser.Append([]string{
				key.Ip,
				strconv.FormatBool(key.V2Upgrade.IsDoubleWriteEnabled),
				strconv.FormatInt(key.V2Upgrade.ServiceCountV2, 10),
				strconv.FormatInt(key.V2Upgrade.InstanceCountV2, 10),
				strconv.FormatInt(key.V2Upgrade.ServiceCountV1, 10),
				strconv.FormatInt(key.V2Upgrade.InstanceCountV1, 10),
				strconv.FormatBool(key.V2Upgrade.Upgraded),
				strconv.FormatBool(key.V2Upgrade.IsAll20XVersion),
			})
		}
		fmt.Printf("v2版本升级接口详情\n")
		tablecluser.Render()
	}
	nacos.Na.Client.CloseIdleConnections()
	if nacos.Writefile != "" {
		nacos.Na.WriteFile()
	}
}
