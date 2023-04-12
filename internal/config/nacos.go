package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tidwall/gjson"
	"gopkg.in/ini.v1"
	"nacos-check/pkg"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex sync.Mutex
var tablerow []string

func (d *Nacos) GetJson(resultType string, web bool) (result interface{}, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	defer func() {
		if funcErr := recover(); funcErr != nil {
			result = []byte("[]")
			err = errors.New("error")
		}
	}()
	if web {
		d.GetNacosInstance()
	}
	var nacos NacosFile
	for _, nacosServer := range d.Clusterdata {
		if len(nacosServer.HealthInstance) != 0 {
			for _, na := range nacosServer.HealthInstance {
				var ta NacosTarget
				ta.Labels = make(map[string]string)
				for key, value := range ADDLABEL {
					ta.Labels[key] = value
				}
				ta.Targets = append(ta.Targets, na.IpAddr)
				ta.Labels["namespace"] = na.NamespaceName
				ta.Labels["service"] = na.ServiceName
				ta.Labels["hostname"] = na.Hostname
				ta.Labels["weight"] = na.Weight
				ta.Labels["pid"] = na.Pid
				ta.Labels["ip"] = na.Ip
				ta.Labels["port"] = na.Port
				ta.Labels["group"] = na.GroupName
				ta.Labels["containerd"] = na.Container
				nacos.Data = append(nacos.Data, ta)
			}
		}
	}

	if resultType == "json" {
		result = nacos.Data
		return result, err
	}
	data, err := json.MarshalIndent(&nacos.Data, "", "  ")
	if err != nil {
		fmt.Println("json序列化失败!")
		result = []byte("[]")
		return result, err
	}
	result = data
	//result = []byte("[]")
	return result, err
}
func (d *Nacos) WriteFile() {
	var basedir string
	var filename string
	basedir = path.Dir(WRITEFILE)
	filename = path.Base(WRITEFILE)
	if err := os.MkdirAll(basedir, os.ModePerm); err != nil {
		os.Exit(EXITCODE)
	}
	file, err := os.OpenFile(basedir+"/.nacos_tmp.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("创建文件失败 %s", err)
		os.Exit(2)
	}
	defer file.Close()
	jsondata, err := d.GetJson("byte", false)
	data := make([]byte, 0)
	var check bool
	if data, check = jsondata.([]byte); !check {
		fmt.Println("转换失败")
	}
	if _, err := file.Write(data); err != nil {
		fmt.Println("写入失败", err)
		os.Exit(EXITCODE)
	}
	file.Close()
	if err := os.Rename(basedir+"/.nacos_tmp.json", basedir+"/"+filename); err != nil {
		fmt.Println("写入失败:", basedir+"/"+filename)
	} else {
		fmt.Println("写入成功:", basedir+"/"+filename)
	}
}

func (d *Nacos) Auth() {
	_url := fmt.Sprintf("%s%s/v1/auth/login", d.DefaultUlr, CONTEXTPATH)
	formData := map[string]string{
		"username": USERNAME,
		"password": PASSWORD,
	}
	res := d.POST(_url, formData)
	if len(gjson.GetBytes(res, "accessToken").String()) != 0 {
		fmt.Println("Authentication successful...")
		d.Token = gjson.GetBytes(res, "accessToken").String()
	} else {
		fmt.Println("Authentication failed!")
	}
}
func (d *Nacos) GetCluster() {
	_url := fmt.Sprintf("%s%s/v1/ns/operator/servers", d.DefaultUlr, CONTEXTPATH)
	res := d.GET(_url)
	d.Cluster = string(res)
}

func (d *Nacos) GetNameSpace() {
	if len(NAMESPACELIST) == 0 {
		_url := fmt.Sprintf("%s%s/v1/console/namespaces", d.DefaultUlr, CONTEXTPATH)
		res := d.GET(_url)
		err := json.Unmarshal(res, &d.Namespaces)
		if err != nil {
			fmt.Println("获取命名空间json异常")
		}
	} else {
		d.Namespaces.Data = NAMESPACELIST
	}
}
func (d *Nacos) GetService(url string, namespaceId string, group string) []byte {
	_url := fmt.Sprintf("%s%s/v1/ns/service/list?pageNo=1&pageSize=500&namespaceId=%s&groupName=%s", url, CONTEXTPATH, namespaceId, group)
	res := d.GET(_url)
	return res
}

func (d *Nacos) GetInstance(url string, servicename string, namespaceId string, group string) []byte {
	_url := fmt.Sprintf("%s%s/v1/ns/instance/list?serviceName=%s&namespaceId=%s&groupName=%s", url, CONTEXTPATH, servicename, namespaceId, group)
	//fmt.Println(_url)
	res := d.GET(_url)
	return res
}

func (d *Nacos) GetV2Upgrade() []byte {
	_url := fmt.Sprintf("%s%s/v1/ns/upgrade/ops/metrics", d.DefaultUlr, CONTEXTPATH)
	res := d.GET(_url)
	return res
}
func (d *Nacos) tableAppend(table *tablewriter.Table, data []string) {
	datastr := strings.Join(data, "-")
	if !pkg.InString(datastr, tablerow) {
		tablerow = append(tablerow, datastr)
		table.Append(data)
	}
}
func (d *Nacos) TableRender() {
	tablerow = []string{}
	nacosServer := d.Clusterdata[d.Host]
	tabletitle := []string{"命名空间", "服务名称", "实例", "健康状态", "主机名", "权重", "容器", "组"}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(tabletitle)
	for _, nacosServer := range d.Clusterdata {
		for _, v := range nacosServer.HealthInstance {
			tabledata := []string{v.NamespaceName, v.ServiceName, v.IpAddr, v.Health, v.Hostname, v.Weight, v.Container, v.GroupName}
			if FIND == "" {
				d.tableAppend(table, tabledata)
			} else {
				for _, find := range FINDLIST {
					if strings.Contains(v.ServiceName, find) {
						d.tableAppend(table, tabledata)
					}
					if strings.Contains(v.ServiceName, find) {
						d.tableAppend(table, tabledata)
					}
					if strings.Contains(v.ServiceName, find) {
						d.tableAppend(table, tabledata)
					}
				}
			}
		}
	}
	fmt.Printf("健康实例:(%d 个)\n", table.NumLines())
	table.Render()
	if len(nacosServer.UnHealthInstance) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(tabletitle)
		for _, v := range nacosServer.UnHealthInstance {
			tabledata := []string{v.NamespaceName, v.ServiceName, v.IpAddr, v.Health, v.Hostname, v.Weight, v.Container, v.GroupName}
			for _, find := range FINDLIST {
				if strings.Contains(v.ServiceName, find) {
					d.tableAppend(table, tabledata)
				}
				if strings.Contains(v.ServiceName, find) {
					d.tableAppend(table, tabledata)
				}
				if strings.Contains(v.ServiceName, find) {
					d.tableAppend(table, tabledata)
				}
			}
		}
		fmt.Printf("异常实例:(%d 个)\n", table.NumLines())
		table.Render()
	}
}
func (d *Nacos) GetNacosInstance() {
	clusterList := []string{d.Host}
	d.Clusterdata = make(map[string]ClusterStatus)
	if CLUSTER {
		d.GetCluster()
		var results []gjson.Result
		var leader gjson.Result
		if len(gjson.Get(d.Cluster, "servers").String()) != 0 {
			leader = gjson.Get(d.Cluster, "servers.0.extendInfo.raftMetaData.metaDataMap.naming_instance_metadata.leader")
			results = gjson.GetMany(d.Cluster, "servers.#.ip", "servers.#.port", "servers.#.state", "servers.#.extendInfo.version", "servers.#.extendInfo.lastRefreshTime")
		} else {
			leader = gjson.Get(d.Cluster, "data.0.extendInfo.raftMetaData.metaDataMap.naming_instance_metadata.leader")
			results = gjson.GetMany(d.Cluster, "data.#.ip", "data.#.port", "data.#.state", "data.#.extendInfo.version", "data.#.extendInfo.lastRefreshTime")
		}
		d.Leader = leader.String()
		for key := range results[0].Array() {
			timeStampStr := results[4].Array()[key].String()
			timeStamp, _ := strconv.Atoi(timeStampStr)
			formatTimeStr := time.Unix(int64(timeStamp/1000), 0).Format("2006-01-02 15:04:05")
			var cluster ClusterStatus
			cluster.Ip = results[0].Array()[key].String()
			cluster.Port = results[1].Array()[key].String()
			cluster.State = results[2].Array()[key].String()
			cluster.Version = results[3].Array()[key].String()
			cluster.LastRefreshTime = formatTimeStr
			key := fmt.Sprintf("%s:%s", results[0].Array()[key].String(), results[1].Array()[key].String())
			d.Clusterdata[key] = cluster
			if !pkg.InString(key, clusterList) {
				clusterList = append(clusterList, key)
			}
		}
	} else {
		var cluster ClusterStatus
		cluster.Ip = d.Host
		cluster.Port = d.Port
		cluster.State = "UP"
		cluster.Version = ""
		cluster.LastRefreshTime = ""
		key := fmt.Sprintf("%s:%s", d.Host, d.Port)
		d.Clusterdata[key] = cluster
		if !pkg.InString(key, clusterList) {
			clusterList = append(clusterList, key)
		}
	}
	if !CLUSTER {
		for _, server := range clusterList {
			_url := fmt.Sprintf("%s://%s", d.Scheme, server)
			if _url == d.DefaultUlr {
				clusterList = []string{server}
			}
		}
	}
	if !CLUSTER && len(clusterList) != 1 {
		_url := fmt.Sprintf("%s", d.Host)
		clusterList = []string{_url}
	}
	for _, server := range clusterList {
		//fmt.Println(server)
		d.GetNameSpace()
		for _, namespace := range d.Namespaces.Data {
			//res := d.GetService(namespace.Namespace)
			_url := fmt.Sprintf("%s://%s", d.Scheme, server)
			var ser Service
			var cluster ClusterStatus
			cluster = d.Clusterdata[server]
			if V2UPGRADE {
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
			//
			for _, group := range GROUPLIST {
				res := d.GetService(_url, namespace.Namespace, group)
				err := json.Unmarshal(res, &ser)
				if err != nil {
					fmt.Println(err)
				}
				for _, se := range ser.Doms {
					res := d.GetInstance(_url, se, namespace.Namespace, group)
					var in Instance
					err := json.Unmarshal(res, &in)
					if err != nil {
						fmt.Printf("json序列化错误:%s\n", err)
					}
					for _, host := range in.Hosts {
						hostname := ""
						_pid := ""
						metadataUrl := host.Metadata["dubbo.metadata-service.urls"]
						u, _ := regexp.Compile("pid=(.+?)&")
						if PARSEIP {
							hostname = GetHostName(host.Ip)
						} else {
							hostname = host.Ip
						}
						pid := u.FindStringSubmatch(metadataUrl)
						if len(pid) == 2 {
							_pid = pid[1]
						}
						instance := ServerInstance{
							NamespaceName: namespace.NamespaceShowName,
							ServiceName:   se,
							IpAddr:        fmt.Sprintf("%s:%d", host.Ip, host.Port),
							Health:        strconv.FormatBool(host.Healthy),
							Hostname:      hostname,
							Weight:        fmt.Sprintf("%.1f", host.Weight),
							Pid:           _pid,
							Container:     strconv.FormatBool(pkg.ContainerdIPCheck(host.Ip)),
							Ip:            host.Ip,
							Port:          strconv.Itoa(host.Port),
							GroupName:     in.GroupName,
						}
						if host.Healthy {
							cluster.HealthInstance = append(cluster.HealthInstance, instance)
						} else {
							cluster.UnHealthInstance = append(cluster.UnHealthInstance, instance)
						}
						d.Clusterdata[server] = cluster
					}
				}
			}
		}
	}
}

func GetHostName(ip string) string {
	for hostname, i := range IPDATA {
		if ip == i {
			return hostname
		}
	}
	return ip
}
