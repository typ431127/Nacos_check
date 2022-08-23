package nacos

import (
	"net"
)

var cidrs []*net.IPNet

func ContainerdInit() {

	// 定义容器网络网段
	maxCidrBlocks := []string{
		"172.30.0.0/16",
		"172.17.0.0/16",
	}

	cidrs = make([]*net.IPNet, len(maxCidrBlocks))
	for i, maxCidrBlock := range maxCidrBlocks {
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
