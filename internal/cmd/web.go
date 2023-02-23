package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/cobra"
	"nacos-check/internal/config"
	"nacos-check/pkg"
	"strconv"
	"strings"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "开启web api Prometheus http_sd_configs",
	Run: func(cmd *cobra.Command, args []string) {
		Webserver()
	},
}
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "注册本实例到Nacos并开启webapi",
	Run: func(cmd *cobra.Command, args []string) {
		go Register()
		Webserver()
	},
}
var svcname string
var ipaddr string

func init() {
	ips := pkg.GetIps()
	webCmd.Flags().StringVarP(&config.WEBPORT, "port", "p", ":8099", "web 端口")
	registerCmd.Flags().StringVarP(&config.WEBPORT, "port", "p", ":8099", "web 端口")
	registerCmd.Flags().StringVarP(&svcname, "name", "n", "nacos-check", "nacos注册名称")
	registerCmd.Flags().StringVarP(&ipaddr, "ip", "i", ips[0], "指定nacos注册客户端ip")
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(registerCmd)
}

func response(c *gin.Context) {
	if c.Request.RequestURI == "/health" {
		c.JSON(200, gin.H{"status": true})
		return
	}
	if c.Request.RequestURI == "/favicon.ico" {
		c.JSON(404, "404")
		return
	}
	result, err := Nacos.GetJson("json")
	if err != nil {
		c.JSON(500, []string{})
		return
	}
	c.JSON(200, result)
}

func Webserver() {
	fmt.Println("Nacos:", config.NACOSURL)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	v1 := r.Group("/")
	{
		v1.GET("/*route", response)
	}
	err := r.Run(config.WEBPORT)
	if err != nil {
		fmt.Println(err)
	}
}

func Register() {
	parseUint, err := strconv.ParseUint(Nacos.Port, 10, 64)
	if err != nil {
		return
	}
	webportUint, err := strconv.ParseUint(strings.Split(config.WEBPORT, ":")[1], 10, 64)
	if err != nil {
		return
	}
	_host := strings.Split(Nacos.Host, ":")[0]
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      _host,
			ContextPath: "/nacos",
			Port:        parseUint,
			Scheme:      Nacos.Scheme,
		},
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: serverConfigs,
		},
	)
	success, err := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ipaddr,
		Port:        webportUint,
		ServiceName: svcname,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"code": "golang"},
		ClusterName: "DEFAULT",       // default value is DEFAULT
		GroupName:   "DEFAULT_GROUP", // default value is DEFAULT_GROUP
	})
	if success {
		fmt.Println("Nacos注册成功")
	} else {
		fmt.Println("Nacos注册失败", err)
	}
}
