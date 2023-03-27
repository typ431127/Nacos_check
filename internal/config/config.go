package config

import (
	"go/types"
	"net/http"
	"time"
)

var (
	NACOSURL      string // nacos url地址
	NACOSURLLIST  []string
	FIND          string // 模糊匹配服务
	FINDLIST      []string
	NAMESPACE     string
	NAMESPACELIST []NamespaceServer
	GROUP         string // 分组
	GROUPLIST     []string
	USERNAME      string            // 账户
	PASSWORD      string            // 密码
	WRITEFILE     string            // prometheus 字段 文件路径
	IPFILE        string            // ip hostname 解析文件
	PARSEIP       bool              // 是否启用ip解析
	CLUSTER       bool              // 集群状态
	IPDATA        map[string]string // 全部ip数据
	EXITCODE      int               // 全局退出状态码
	VERSION       bool              // 版本
	WATCH         bool              // 监控
	SECOND        time.Duration     // 监控服务间隔
	V2UPGRADE     bool              // 2.0版本升级详情
	EXPORTJSON    bool              // 导出json
	WEB           bool              // 开启webapi
	WEBPORT       string            // web端口
	ADDLABEL      map[string]string
	Na            *Nacos
)

type Nacos struct {
	Client         http.Client
	Namespaces     Namespaces
	DefaultUlr     string
	Host           string
	Scheme         string
	Token          string
	Port           string
	Healthydata    [][]string
	Healthydataerr [][]string
	AllInstance    [][]string
	Cluster        string
	Leader         string
	Clusterdata    map[string]ClusterStatus
}

type ClusterStatus struct {
	Ip               string
	Port             string
	State            string
	Version          string
	LastRefreshTime  string
	HealthInstance   []ServerInstance
	UnHealthInstance []ServerInstance
	V2Upgrade        V2Upgrade
}
type ServerInstance struct {
	NamespaceName string `json:"namespaceName"`
	ServiceName   string `json:"serviceName"`
	IpAddr        string `json:"ipAddr"`
	Health        string `json:"health"`
	Hostname      string `json:"hostname"`
	Weight        string `json:"weight"`
	Pid           string `json:"pid"`
	Container     string `json:"container"`
	Ip            string `json:"ip"`
	Port          string `json:"port"`
	GroupName     string `json:"groupName"`
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
	GroupName       string      `json:"groupName"`
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
