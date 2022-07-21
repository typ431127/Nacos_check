package ext

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"go/types"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	nacosurl       string            // nacos url地址
	findstr        string            // 模糊匹配服务
	noconsole      bool              // 是否控制台输出
	writefile      string            // prometheus 字段 文件路径
	ipfile         string            // ip hostname 解析文件
	ipparse        bool              // 是否启用ip解析
	cluster_status bool              // 集群状态
	ipdata         map[string]string // 全部ip数据
	exitcode       int               // 全局退出状态码
	version        bool              // 版本
	watch          bool              // 监控
	second         time.Duration     // 监控服务间隔
	v2upgrade      bool              // 2.0版本升级详情
	_json          bool              // 导出json
)

type Nacos struct {
	Client         http.Client
	Namespaces     Namespaces
	DefaultUlr     string
	Host           string
	Scheme         string
	Healthydata    [][]string
	Healthydataerr [][]string
	AllInstance    [][]string
	Cluster        string
	clusterdata    map[string]ClusterStatus // 集群状态数据
}

type ClusterStatus struct {
	Ip               string
	Port             string
	state            string
	Version          string
	LastRefreshTime  string
	HealthInstance   [][]string
	UnHealthInstance [][]string
	V2Upgrade        V2Upgrade
}

type V2Upgrade struct {
	Upgraded             bool
	IsAll20XVersion      bool
	IsDoubleWriteEnabled bool
	ServiceCountV1       int64
	InstanceCountV1      int64
	ServiceCountV2       int64
	InstanceCountV2      int64
	SubscribeCountV2     int64
}
type NamespaceServer struct {
	Namespace         string `json:"namespace"`
	NamespaceShowName string `json:"namespaceShowName"`
	Quota             int    `json:"quota"`
	ConfigCount       int    `json:"configCount"`
	Type              int    `json:"type"`
}

type Namespaces struct {
	Code    int               `json:"code"`
	Message types.Nil         `json:"message"`
	Data    []NamespaceServer `json:"data"`
}

type service struct {
	Doms  []string `json:"doms"`
	Count int      `json:"count"`
}

type Instance struct {
	Hosts           []Instances `json:"hosts"`
	Dom             string      `json:"dom"`
	Name            string      `json:"name"`
	CacheMillis     int         `json:"cacheMillis"`
	LastRefTime     int64       `json:"lastRefTime"`
	Checksum        string      `json:"checksum"`
	UseSpecifiedURL bool        `json:"useSpecifiedURL"`
	Clusters        string      `json:"clusters"`
	Env             string      `json:"env"`
	Metadata        map[string]interface{}
}

type Instances struct {
	Ip                        string `json:"ip"`
	Port                      int    `json:"port"`
	Valid                     bool   `json:"valid"`
	Healthy                   bool   `json:"healthy"`
	Marked                    bool   `json:"marked"`
	InstanceId                string `json:"instanceId"`
	Metadata                  map[string]string
	Enabled                   bool    `json:"enabled"`
	Weight                    float32 `json:"weight"`
	ClusterName               string  `json:"clusterName"`
	ServiceName               string  `json:"serviceName"`
	Ephemeral                 bool    `json:"ephemeral"`
	InstanceHeartBeatInterval int64   `json:"instanceHeartBeatInterval"`
}

type NacosTarget struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels"`
}

type NacosFile struct {
	Data []NacosTarget
}

func init() {
	flag.StringVar(&nacosurl, "url", "http://dev-k8s-nacos:8848", "nacos地址")
	flag.StringVar(&writefile, "write", "", "prometheus 自动发现文件路径")
	flag.StringVar(&ipfile, "ipfile", "salt_ip.json", "ip解析文件")
	flag.StringVar(&findstr, "find", "", "查找服务")
	flag.BoolVar(&noconsole, "noconsole", false, "不输出console")
	flag.BoolVar(&_json, "json", false, "输出json")
	flag.BoolVar(&cluster_status, "cluster", false, "查看集群状态")
	flag.BoolVar(&v2upgrade, "v2upgrade", false, "查看2.0升级状态,和-cluster一起使用")
	flag.BoolVar(&version, "version", false, "查看版本")
	flag.BoolVar(&watch, "watch", false, "监控服务")
	flag.DurationVar(&second, "second", 2*time.Second, "监控服务间隔刷新时间")
}

func FilePathCheck() {
	if _, err := os.Stat(ipfile); err != nil {
		if !os.IsExist(err) {
			ipparse = false
			return
		}
	} else {
		ipparse = true
	}
}

func IP_Parse() {
	file, err := os.OpenFile(ipfile, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("打开文件错误")
		os.Exit(exitcode)
	}
	defer file.Close()
	fileb, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(fileb, &ipdata); err != nil {
		fmt.Println("ip文件解析错误,请确认json格式")
		os.Exit(exitcode)
	}
}

func GetHostName(ip string) string {
	for hostname, i := range ipdata {
		if ip == i {
			return hostname
		}
	}
	return ip
}

func (d *Nacos) GetJson() []byte {
	nacos_server := d.clusterdata[d.Host]
	if len(nacos_server.HealthInstance) != 0 {
		var nacos NacosFile
		for _, na := range nacos_server.HealthInstance {
			var ta NacosTarget
			ta.Labels = make(map[string]string)
			ta.Targets = append(ta.Targets, na[2])
			ta.Labels["namespace"] = na[0]
			ta.Labels["service"] = na[1]
			ta.Labels["hostname"] = na[4]
			ta.Labels["weight"] = na[5]
			ta.Labels["pid"] = na[6]
			ta.Labels["ip"] = na[8]
			ta.Labels["port"] = na[9]
			ta.Labels["containerd"] = na[7]
			nacos.Data = append(nacos.Data, ta)
		}
		data, err := json.MarshalIndent(&nacos.Data, "", "  ")
		if err != nil {
			fmt.Println("json序列化失败!")
			os.Exit(exitcode)
		}
		return data
	}
	return []byte("[]")
}
func (d *Nacos) WriteFile() {
	var basedir string
	var filename string
	basedir = path.Dir(writefile)
	filename = path.Base(writefile)
	if err := os.MkdirAll(basedir, os.ModePerm); err != nil {
		os.Exit(exitcode)
	}
	file, err := os.OpenFile(basedir+"/.nacos_tmp.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("创建文件失败 %s", err)
		os.Exit(2)
	}
	defer file.Close()
	data := d.GetJson()
	if _, err := file.Write(data); err != nil {
		fmt.Println("写入失败", err)
		os.Exit(exitcode)
	}
	file.Close()
	os.Rename(basedir+"/.nacos_tmp.json", basedir+"/"+filename)
	fmt.Println("写入成功:", basedir+"/"+filename)
}

func (d *Nacos) HttpReq(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	res, err := d.Client.Do(req)
	if err != nil {
		fmt.Println("请求异常:", err)
		os.Exit(exitcode)
	}
	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	return resp

}

func (d *Nacos) GetCluster() {
	_url := fmt.Sprintf("%s/nacos/v1/ns/operator/servers", nacosurl)
	res := d.HttpReq(_url)
	d.Cluster = string(res)
}

func (d *Nacos) GetNameSpace() {
	_url := fmt.Sprintf("%s/nacos/v1/console/namespaces", nacosurl)
	res := d.HttpReq(_url)
	err := json.Unmarshal(res, &d.Namespaces)
	if err != nil {
		fmt.Println("获取命名空间json异常")
		os.Exit(2)
	}
}
func (d *Nacos) GetService(namespaceId string) []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/service/list?pageNo=1&pageSize=500&namespaceId=%s", nacosurl, namespaceId)
	res := d.HttpReq(_url)
	return res
}

func (d *Nacos) GetInstance(servicename string, namespaceId string) []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/instance/list?serviceName=%s&namespaceId=%s", nacosurl, servicename, namespaceId)
	res := d.HttpReq(_url)
	return res
}

func (d *Nacos) GetV2Upgrade() []byte {
	_url := fmt.Sprintf("%s/nacos/v1/ns/upgrade/ops/metrics", nacosurl)
	res := d.HttpReq(_url)
	return res
}

func (d *Nacos) PrintInfo() {
	nacos_server := d.clusterdata[d.Host]
	tabletitle := []string{"命名空间", "服务名称", "实例", "健康状态", "主机名", "权重", "PID", "容器"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	for _, v := range nacos_server.HealthInstance {
		tabledata := v[0:8]
		if findstr == "" {
			table.Append(tabledata)
		} else {
			if strings.Contains(v[0], findstr) {
				table.Append(tabledata)
			}
			if strings.Contains(v[1], findstr) {
				table.Append(tabledata)
			}
			if strings.Contains(v[2], findstr) {
				table.Append(tabledata)
			}
		}
	}
	fmt.Printf("健康实例:(%d 个)\n", table.NumLines())
	table.Render()
	if len(nacos_server.UnHealthInstance) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(tabletitle)
		for _, v := range d.Healthydataerr {
			if strings.Contains(v[0], findstr) {
				table.Append(v)
			}
			if strings.Contains(v[1], findstr) {
				table.Append(v)
			}
			if strings.Contains(v[2], findstr) {
				table.Append(v)
			}
		}
		fmt.Printf("异常实例:(%d 个)\n", table.NumLines())
		table.Render()
	}
}
func (d *Nacos) GetNacosInstance() {
	d.GetCluster()
	d.clusterdata = make(map[string]ClusterStatus)
	results := gjson.GetMany(d.Cluster, "servers.#.ip", "servers.#.port", "servers.#.state", "servers.#.extendInfo.version", "servers.#.extendInfo.lastRefreshTime")
	cluster_list := []string{}
	for key := range results[0].Array() {
		timeStampStr := results[4].Array()[key].String()
		timeStamp, _ := strconv.Atoi(timeStampStr)
		formatTimeStr := time.Unix(int64(timeStamp/1000), 0).Format("2006-01-02 15:04:05")
		var cluster ClusterStatus
		cluster.Ip = results[0].Array()[key].String()
		cluster.Port = results[1].Array()[key].String()
		cluster.state = results[2].Array()[key].String()
		cluster.Version = results[3].Array()[key].String()
		cluster.LastRefreshTime = formatTimeStr
		key := fmt.Sprintf("%s:%s", results[0].Array()[key].String(), results[1].Array()[key].String())
		d.clusterdata[key] = cluster
		cluster_list = append(cluster_list, key)
	}
	if !cluster_status {
		for _, server := range cluster_list {
			url := fmt.Sprintf("%s://%s", d.Scheme, server)
			if url == d.DefaultUlr {
				cluster_list = []string{server}
			}
		}
	}
	if !cluster_status && len(cluster_list) != 1 {
		url := fmt.Sprintf("%s", d.Host)
		cluster_list = []string{url}
	}
	for _, server := range cluster_list {
		nacosurl = fmt.Sprintf("%s://%s", d.Scheme, server)
		d.GetNameSpace()
		for _, namespace := range d.Namespaces.Data {
			res := d.GetService(namespace.Namespace)
			var ser service
			var cluster ClusterStatus
			cluster = d.clusterdata[server]
			if v2upgrade {
				v2 := d.GetV2Upgrade()
				rep, _ := regexp.Compile(".*##.*")
				v2 = rep.ReplaceAll(v2, []byte(""))
				cfg, err := ini.Load(v2)
				if err != nil {
					fmt.Println("v2解析错误")
				}
				IsDoubleWriteEnabled, _ := cfg.Section("").Key("isDoubleWriteEnabled").Bool()
				Upgraded, _ := cfg.Section("").Key("upgraded").Bool()
				IsAll20XVersion, _ := cfg.Section("").Key("isAll20XVersion").Bool()
				ServiceCountV1, _ := cfg.Section("").Key("serviceCountV1").Int64()
				InstanceCountV1, _ := cfg.Section("").Key("instanceCountV1").Int64()
				ServiceCountV2, _ := cfg.Section("").Key("serviceCountV2").Int64()
				InstanceCountV2, _ := cfg.Section("").Key("instanceCountV2").Int64()
				SubscribeCountV2, _ := cfg.Section("").Key("subscribeCountV2").Int64()
				cluster.V2Upgrade.IsDoubleWriteEnabled = IsDoubleWriteEnabled
				cluster.V2Upgrade.Upgraded = Upgraded
				cluster.V2Upgrade.IsAll20XVersion = IsAll20XVersion
				cluster.V2Upgrade.ServiceCountV1 = ServiceCountV1
				cluster.V2Upgrade.InstanceCountV1 = InstanceCountV1
				cluster.V2Upgrade.ServiceCountV2 = ServiceCountV2
				cluster.V2Upgrade.InstanceCountV2 = InstanceCountV2
				cluster.V2Upgrade.SubscribeCountV2 = SubscribeCountV2
			}
			json.Unmarshal(res, &ser)
			for _, se := range ser.Doms {
				res := d.GetInstance(se, namespace.Namespace)
				var in Instance
				err := json.Unmarshal(res, &in)
				if err != nil {
					fmt.Println("json序列化错误:%s", err)
				}
				for _, host := range in.Hosts {
					metadataUrl := host.Metadata["dubbo.metadata-service.urls"]
					u, _ := regexp.Compile("pid=(.+?)&")
					_tmpmap := make([]string, 0)
					ipinfo := fmt.Sprintf("%s:%d", host.Ip, host.Port)
					_tmpmap = append(_tmpmap, namespace.NamespaceShowName)
					_tmpmap = append(_tmpmap, se)
					_tmpmap = append(_tmpmap, ipinfo)
					_tmpmap = append(_tmpmap, strconv.FormatBool(host.Healthy))
					if ipparse {
						_tmpmap = append(_tmpmap, GetHostName(host.Ip))
					} else {
						_tmpmap = append(_tmpmap, host.Ip)
					}
					_tmpmap = append(_tmpmap, fmt.Sprintf("%.0f", host.Weight))
					pid := u.FindStringSubmatch(metadataUrl)
					if len(pid) == 2 {
						_tmpmap = append(_tmpmap, pid[1])
					} else {
						_tmpmap = append(_tmpmap, "")
					}
					_tmpmap = append(_tmpmap, strconv.FormatBool(ContainerdIPCheck(host.Ip)))
					_tmpmap = append(_tmpmap, host.Ip)
					_tmpmap = append(_tmpmap, strconv.Itoa(host.Port))
					if host.Healthy {
						cluster.HealthInstance = append(cluster.HealthInstance, _tmpmap)
					} else {
						cluster.UnHealthInstance = append(cluster.UnHealthInstance, _tmpmap)
					}
					d.clusterdata[server] = cluster
				}
			}
		}
	}
}
func ParseCheck() {
	FilePathCheck()
	ContainerdInit()
	if ipparse {
		IP_Parse()
	}
}

func NacosInit() {
	flag.Parse()
	ParseCheck()
	u, err := url.Parse(nacosurl)
	if err != nil {
		fmt.Println("url解析错误!")
		os.Exit(exitcode)
	}
	na := new(Nacos)
	na.Client.Timeout = time.Second * 15
	na.DefaultUlr = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	na.Host = u.Host
	na.Scheme = u.Scheme
	na.GetNacosInstance()
	if _json {
		data := na.GetJson()
		fmt.Println(string(data))
		os.Exit(0)
	}
	if version {
		fmt.Println("版本:0.4.1")
		os.Exit(0)
	}
	if !noconsole && !cluster_status {
		fmt.Println("Nacos:", nacosurl)
		if watch {
			for {
				na.GetNacosInstance()
				na.PrintInfo()
				time.Sleep(second)
			}
		}
		na.PrintInfo()
	}
	if cluster_status {
		tablecluser := tablewriter.NewWriter(os.Stdout)
		tablecluser.SetHeader([]string{"节点", "端口", "状态", "版本", "刷新时间", "健康实例", "异常实例"})
		for _, key := range na.clusterdata {
			tablecluser.Append([]string{
				key.Ip,
				key.Port,
				key.state,
				key.Version,
				key.LastRefreshTime,
				strconv.Itoa(len(key.HealthInstance)),
				strconv.Itoa(len(key.UnHealthInstance)),
			})
		}
		leader := gjson.Get(na.Cluster, "servers.0.extendInfo.raftMetaData.metaDataMap.naming_instance_metadata.leader")
		fmt.Printf("Nacos集群状态: (数量:%d)\n集群Master: %s\n", tablecluser.NumLines(), leader)
		tablecluser.Render()
	}
	if v2upgrade && cluster_status {
		tablecluser := tablewriter.NewWriter(os.Stdout)
		tablecluser.SetHeader([]string{"节点", "双写", "v2服务", "v2实例", "v1服务", "v1实例", "upgraded", "全部V2"})
		for _, key := range na.clusterdata {
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
	na.Client.CloseIdleConnections()
	if writefile != "" {
		na.WriteFile()
	}

}
