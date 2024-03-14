package cmd

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/cobra"
	"nacos-check/internal/nacos"
	"nacos-check/pkg"
	"net/url"
	"strconv"
	"strings"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "注册本实例到Nacos并开启webapi",
	Run: func(cmd *cobra.Command, args []string) {
		go Register()
		Webserver()
	},
}
var (
	svcname     string
	groupname   string
	namespaceid string
	ipaddr      string
)

func init() {
	ips := pkg.GetIps()
	registerCmd.Flags().StringVarP(&nacos.WEBPORT, "port", "p", ":8099", "web 端口")
	registerCmd.Flags().StringVarP(&svcname, "name", "n", "nacos-check", "nacos注册名称")
	registerCmd.Flags().StringVarP(&ipaddr, "ip", "i", ips[0], "指定nacos注册客户端ip")
	registerCmd.Flags().StringVarP(&namespaceid, "namespace", "", "", "指定要注册到的namespaceid")
	registerCmd.Flags().StringVarP(&groupname, "groupname", "g", "DEFAULT_GROUP", "指定注册的分组名称")
	rootCmd.AddCommand(registerCmd)
}

func Register() {
	var serverConfigs []constant.ServerConfig
	webportUint, err := strconv.ParseUint(strings.Split(nacos.WEBPORT, ":")[1], 10, 64)
	if err != nil {
		fmt.Println(err)
	}
	for _, _url := range strings.Split(nacos.NACOSURL, ",") {
		parse, _ := url.Parse(_url)
		_host := strings.Split(parse.Host, ":")[0]
		parseUint, err := strconv.ParseUint(parse.Port(), 10, 64)
		if err != nil {
			fmt.Println(err)
		}
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr:      _host,
			ContextPath: "/nacos",
			Port:        parseUint,
			Scheme:      parse.Scheme,
		})
	}
	var clientConfigs constant.ClientConfig
	if len(nacos.USERNAME) != 0 && len(nacos.PASSWORD) != 0 {
		clientConfigs = constant.ClientConfig{NamespaceId: namespaceid, Username: nacos.USERNAME, Password: nacos.PASSWORD}
	} else {
		clientConfigs = constant.ClientConfig{NamespaceId: namespaceid}
	}
	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ServerConfigs: serverConfigs,
			ClientConfig:  &clientConfigs,
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
		ClusterName: "DEFAULT", // default value is DEFAULT
		GroupName:   groupname, // default value is DEFAULT_GROUP
	})
	if success {
		fmt.Println("Nacos注册成功")
	} else {
		fmt.Println("Nacos注册失败", err)
	}
}
