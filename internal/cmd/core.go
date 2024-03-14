package cmd

import (
	"fmt"
	"nacos-check/internal/nacos"
	"nacos-check/pkg"
	"net/url"
	"os"
	"strings"
)

func preFunc() {
	nacosFilePathLoad()
	ipconfigLoad()
	pkg.ContainerdInit()
	for _, _url := range strings.Split(nacos.NACOSURL, ",") {
		u, _ := url.Parse(_url)
		if u.Scheme != "http" && u.Scheme != "https:" {
			fmt.Println("The URL is malformed:", _url)
			os.Exit(nacos.EXITCODE)
		}
		nacos.NACOSURLLIST = append(nacos.NACOSURLLIST, _url)
	}
	u, err := url.Parse(nacos.NACOSURLLIST[0])
	if err != nil {
		fmt.Println("url解析错误!")
		os.Exit(nacos.EXITCODE)
	}
	if !strings.HasPrefix(nacos.CONTEXTPATH, "/") {
		fmt.Println("CONTEXT-PATH解析错误!")
		os.Exit(nacos.EXITCODE)
	}
	nacos.FINDLIST = strings.Split(nacos.FIND, ",")
	for _, namespace := range strings.Split(nacos.NAMESPACE, ",") {
		if len(strings.TrimSpace(namespace)) != 0 {
			nacos.NAMESPACELIST = append(nacos.NAMESPACELIST, nacos.NamespaceServer{
				Namespace:         namespace,
				NamespaceShowName: namespace,
				Quota:             200,
				Type:              2,
			})
		}
	}
	for _, group := range strings.Split(nacos.GROUP, ",") {
		if !pkg.InString(group, nacos.GROUPLIST) {
			nacos.GROUPLIST = append(nacos.GROUPLIST, group)
		}
	}
	Nacos = nacos.NewNacosClint(nacos.NACOSURLLIST[0], u.Host, u.Scheme, u.Port())
	if len(nacos.USERNAME) != 0 && len(nacos.PASSWORD) != 0 {
		Nacos.WithAuth()
	}
	Nacos.GetNameSpace()
	Nacos.GetNacosInstance()
}
