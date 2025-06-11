package nacos

import (
	"fmt"
	"net"
)

func isIPInCIDR(ipStr, cidrStr string) (bool, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, fmt.Errorf("无效的 IP 地址: %s", ipStr)
	}
	_, ipNet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return false, fmt.Errorf("无效的 CIDR 网段: %s", cidrStr)
	}
	return ipNet.Contains(ip), nil
}
func LookupIDC(ip string) string {
	for _, cidr := range NETCIDR {
		stat, _ := isIPInCIDR(ip, cidr)
		if stat {
			return NETDATA[cidr]
		}
	}
	return ""
}
