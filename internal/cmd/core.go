package cmd

import (
	"nacos-check/internal/nacos"
	"nacos-check/pkg"
	"nacos-check/pkg/fmtd"
	"net/url"
	"strings"
)

func preFunc() {
	nacosFilePathLoad()
	ipconfigLoad()
	pkg.ContainerdInit()
	for _, _url := range strings.Split(nacos.NACOSURL, ",") {
		u, err := url.Parse(_url)
		if err != nil {
			fmtd.Fatalln(err)
		}
		if u.Scheme != "http" && u.Scheme != "https:" {
			fmtd.Fatalln("The URL is malformed:", _url)
		}
		nacos.NACOSURLLIST = append(nacos.NACOSURLLIST, _url)
	}
	u, err := url.Parse(nacos.NACOSURLLIST[0])
	if err != nil {
		fmtd.Fatalln("url解析错误!")
	}
	if !strings.HasPrefix(nacos.CONTEXTPATH, "/") {
		fmtd.Fatalln("CONTEXT-PATH解析错误!")
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
