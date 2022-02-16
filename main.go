package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"go/types"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var nacosurl string          // nacos url地址
var findstr string           // 模糊匹配服务
var noconsole bool           // 是否控制台输出
var prometheusfile string    // prometheus 文件路径
var ipfile string            // ip hostname 解析文件
var ipparse bool             // 是否启用ip解析
var ipdata map[string]string // 全部ip数据

type nacos struct {
	Client     http.Client
	Namespaces namespaces
	Healthydata [][]string
	Healthydataerr [][]string
	AllInstance [][]string
}

type namespaceser struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	Quota             int    `json:"quota"`
	ConfigCount       int    `json:"configCount"`
	Type              int    `json:"type"`
}

type namespaces struct {
	Code    int            `json:"code"`
	Message types.Nil      `json:"message"`
	Data    []namespaceser `json:"data"`
}

type service struct {
	Doms  []string `json:"doms"`
	Count int      `json:"count"`
}

type Instance struct {
	Hosts []Instances `json:"hosts"`
	Dom string `json:"dom"`
	Name string `json:"name"`
	CacheMillis int `json:"cacheMillis"`
	LastRefTime int64 `json:"lastRefTime"`
	Checksum string `json:"checksum"`
	UseSpecifiedURL bool `json:"useSpecifiedURL"`
	Clusters string `json:"clusters"`
	Env string `json:"env"`
	Metadata map[string]interface{}
}

type Instances struct {
	Ip string `json:"ip"`
	Port int `json:"port"`
	Valid bool `json:"valid"`
	Healthy bool `json:"healthy"`
	Marked bool `json:"marked"`
	InstanceId string `json:"instanceId"`
	Metadata map[string]string
	Enabled bool `json:"enabled"`
	Weight float32 `json:"weight"`
	ClusterName string `json:"clusterName"`
	ServiceName string `json:"serviceName"`
	Ephemeral bool `json:"ephemeral"`
	InstanceHeartBeatInterval  int64 `json:"instanceHeartBeatInterval"`
}

type nacostarget struct {
	Targets []string `json:"targets"`
	Labels map[string]string `json:"labels"`
}

type nacosfile struct {
	Data []nacostarget
}

func init() {
	flag.StringVar(&nacosurl,"url", "http://nacos.ddn.svc.cluster.local:8848", "nacos地址")
	flag.StringVar(&prometheusfile,"write","/data/work/prometheus/discovery/nacos.json","prometheus 自动发现文件路径")
	flag.StringVar(&ipfile,"ipfile", "salt_ip.json", "ip解析文件")
	flag.StringVar(&findstr,"find","","查找服务")
	flag.BoolVar(&noconsole,"noconsole", false, "输出console")
}


func Filepathcheck (path string) bool{
	if _,err := os.Stat(path);err !=nil{
		if ! os.IsExist(err){
			ipparse = false
			return false
		}
	}
	ipparse = true
	fmt.Println("IP解析文件:",ipfile)
	return true
}

// Saltipparse ip解析
func Saltipparse() {
	file,err := os.OpenFile(ipfile,os.O_RDONLY,0644)
	if err != nil{
		fmt.Println("打开文件错误")
		os.Exit(2)
	}
	defer file.Close()
	fileb,_ :=ioutil.ReadAll(file)
	if err := json.Unmarshal(fileb,&ipdata);err != nil{
		fmt.Println("ip文件解析错误")
		os.Exit(2)
	}
}

func Gethostname(ip string) string{
	for hostname,i := range ipdata{
		if ip == i{
			return hostname
		}
	}
	return "None"
}


func (d *nacos) WriteFile() {
	var basedir string
	var filename string
	basedir = path.Dir(prometheusfile)
	filename = path.Base(prometheusfile)
	fmt.Println(filename)
	os.MkdirAll(basedir,os.ModePerm)
	file,err := os.OpenFile(basedir + "/nacos_tmp.json",os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		fmt.Printf("创建文件失败 %s",err)
		os.Exit(2)
	}
	defer file.Close()
	var nacos nacosfile
	for _,na := range d.AllInstance {
		var ta nacostarget
		ta.Labels = make(map[string]string)
		ta.Targets = append(ta.Targets,na[2])
		ta.Labels["namespace"] = na[0]
		ta.Labels["service"] = na[1]
		ta.Labels["hostname"] = na[4]
		ta.Labels["weight"] = na[5]
		ta.Labels["pid"] = na[6]
		nacos.Data = append(nacos.Data,ta)
	}
	data,err := json.MarshalIndent(&nacos.Data,"","  ")
	if err != nil{
		fmt.Println("json序列化失败!")
		os.Exit(2)
	}
	if _,err := file.Write(data);err != nil{
		fmt.Println("写入失败",err)
	}
	file.Close()
	os.Rename(basedir + "/nacos_tmp.json",basedir + "/nacos.json")
	fmt.Println("写入成功:",basedir + "/nacos.json")
}
func (d *nacos) HttpReq(url string, stu interface{}) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := d.Client.Do(req)
	if err != nil {
		fmt.Println("请求异常:",err)
		os.Exit(2)
	}
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	return resp

}
func (d *nacos) GetNamespace() {
	url := fmt.Sprintf("%s/nacos/v1/console/namespaces", nacosurl)
	res := d.HttpReq(url, namespaces{})
	err := json.Unmarshal(res, &d.Namespaces)
	if err != nil{
		fmt.Println("获取命名空间json异常");os.Exit(2)
	}
}
func (d *nacos) GetService(namespaceId string) []byte {
	url := fmt.Sprintf("%s/nacos/v1/ns/service/list?pageNo=1&pageSize=500&namespaceId=%s", nacosurl, namespaceId)
	res := d.HttpReq(url, service{})
	return res
}
func (d *nacos) GetInstance(servicename string,namespaceId string) []byte {
	url := fmt.Sprintf("%s/nacos/v1/ns/instance/list?serviceName=%s&namespaceId=%s", nacosurl, servicename,namespaceId)
	res := d.HttpReq(url, service{})
	return res
}
func main() {
	flag.Parse()
	Filepathcheck(ipfile)
	if ipparse{
		Saltipparse()
	}
	fmt.Println("Nacos:",nacosurl)
	tabletitle := []string{"命名空间", "服务名称", "实例","健康状态","主机名","权重","PID"}
	na := new(nacos)
	na.Client.Timeout = time.Second * 10
	na.GetNamespace()
	for _, namespace := range na.Namespaces.Data {
		res := na.GetService(namespace.Namespace)
		var ser service
		json.Unmarshal(res,&ser)
		for _,se := range ser.Doms {
			res := na.GetInstance(se,namespace.Namespace)
			var in Instance
			err := json.Unmarshal(res,&in)
			if err != nil{
				fmt.Println("json序列化错误:%s",err)
			}
			for _,host := range in.Hosts{
				metadataUrl := host.Metadata["dubbo.metadata-service.urls"]
				u,_ := regexp.Compile("pid=(.+?)&")
				_tmpmap := make([]string,0)
				ipinfo := fmt.Sprintf("%s:%d",host.Ip,host.Port)
				_tmpmap = append(_tmpmap,namespace.NamespaceShowName)
				_tmpmap = append(_tmpmap,se)
				_tmpmap = append(_tmpmap,ipinfo)
				_tmpmap = append(_tmpmap,strconv.FormatBool(host.Healthy))
				if ipparse{
					_tmpmap = append(_tmpmap,Gethostname(host.Ip))
				}else{
					_tmpmap = append(_tmpmap,"None")
				}
				_tmpmap = append(_tmpmap,fmt.Sprintf("%.0f",host.Weight))
				pid := u.FindStringSubmatch(metadataUrl)
				if len(pid) == 2{
					_tmpmap = append(_tmpmap,pid[1])
				}else{
					_tmpmap = append(_tmpmap,"")
				}
				if host.Healthy{
					if findstr == ""{
						na.Healthydata = append(na.Healthydata,_tmpmap)
					}else{
						if strings.Contains(se, findstr){
							na.Healthydata = append(na.Healthydata,_tmpmap)
						}
					}
				}else{
					if findstr == ""{
						na.Healthydataerr = append(na.Healthydataerr,_tmpmap)
					}else{
						if strings.Contains(se, findstr){
							na.Healthydataerr = append(na.Healthydataerr,_tmpmap)
						}
					}
				}
				na.AllInstance = append(na.AllInstance,_tmpmap)
			}
		}
	}
	if ! noconsole {
		// 正常实例
		fmt.Printf("健康实例:(%d 个)\n",len(na.Healthydata))
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(tabletitle)
		for _, v := range na.Healthydata {
			table.Append(v)
		}
		table.Render()
		// 异常实例
		fmt.Printf("异常实例:(%d 个)\n",len(na.Healthydataerr))
		tableerr := tablewriter.NewWriter(os.Stdout)
		tableerr.SetHeader(tabletitle)
		for _, v := range na.Healthydataerr {
			tableerr.Append(v)
		}
		tableerr.Render()
	}
	na.Client.CloseIdleConnections()
	na.WriteFile()

}
