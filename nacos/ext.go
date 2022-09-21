package nacos

import (
	"net"
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
