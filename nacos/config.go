package nacos

import (
	"go/types"
	"net/http"
	"time"
)

var (
	Nacosurl       string            // nacos url地址
	Findstr        string            // 模糊匹配服务
	Writefile      string            // prometheus 字段 文件路径
	Ipfile         string            // ip hostname 解析文件
	Ipparse        bool              // 是否启用ip解析
	Cluster_status bool              // 集群状态
	Ipdata         map[string]string // 全部ip数据
	Exitcode       int               // 全局退出状态码
	Version        bool              // 版本
	Watch          bool              // 监控
	Second         time.Duration     // 监控服务间隔
	V2upgrade      bool              // 2.0版本升级详情
	ExportJson     bool              // 导出json
	Web            bool              // 开启webapi
	Port           string            // web端口
	AddLable       map[string]string
	Na             *Nacos
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
	Clusterdata    map[string]ClusterStatus // 集群状态数据
}

type ClusterStatus struct {
	Ip               string
	Port             string
	State            string
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

type Service struct {
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
