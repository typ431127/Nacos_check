package pkg

import (
	"fmt"
	"net"
	"strings"
)

var cidrs []*net.IPNet

var MaxCidrBlocks = []string{
	"172.30.0.0/16",
	"172.17.0.0/16",
}

func ContainerdInit() {

	cidrs = make([]*net.IPNet, len(MaxCidrBlocks))
	for i, maxCidrBlock := range MaxCidrBlocks {
		_, cidr, _ := net.ParseCIDR(maxCidrBlock)
		cidrs[i] = cidr
	}

}

func ContainerdIPCheck(ip string) bool {
	for i := range cidrs {
		if cidrs[i].Contains(net.ParseIP(ip)) {
			return true
		}
	}
	return false
}

//func IdcNetCIDR(ipaddr string) {
//	for _, cidr := range nacos.NETDATA {
//		isIPInCIDR(ipaddr, cidr)
//	}
//}

func GetIps() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interfaces ipAddress: %v\n", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if !strings.HasPrefix(ipNet.IP.To4().String(), "169.254") && ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

func InString(filed string, array []string) bool {
	for _, str := range array {
		if filed == str {
			return true
		}
	}
	return false
}
