package cmd

import (
	"fmt"
	"nacos-check/internal/config"
	"nacos-check/pkg"
	"net/url"
	"os"
	"strings"
)

func PreFunc() {
	NacosFilePathLoad()
	for _, _url := range strings.Split(config.NACOSURL, ",") {
		config.NACOSURLLIST = append(config.NACOSURLLIST, _url)
	}
	u, err := url.Parse(config.NACOSURLLIST[0])
	if err != nil {
		fmt.Println("url解析错误!")
		os.Exit(config.EXITCODE)
	}
	config.FINDLIST = strings.Split(config.FIND, ",")
	for _, namespace := range strings.Split(config.NAMESPACE, ",") {
		if len(strings.TrimSpace(namespace)) != 0 {
			config.NAMESPACELIST = append(config.NAMESPACELIST, config.NamespaceServer{
				Namespace:         namespace,
				NamespaceShowName: namespace,
				Quota:             200,
				Type:              2,
			})
		}
	}
	for _, group := range strings.Split(config.GROUP, ",") {
		if !pkg.InString(group, config.GROUPLIST) {
			config.GROUPLIST = append(config.GROUPLIST, group)
		}
	}
	Nacos.DefaultUlr = config.NACOSURLLIST[0]
	Nacos.Host = u.Host
	Nacos.Scheme = u.Scheme
	Nacos.Port = u.Port()
	if len(config.USERNAME) != 0 && len(config.PASSWORD) != 0 {
		Nacos.Auth()
	}
	Nacos.GetNameSpace()
	Nacos.GetNacosInstance()
}
