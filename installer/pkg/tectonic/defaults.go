package tectonic

import (
	"fmt"
	"net"
)

// Kubernetes default constants for Tectonic
// Some values must match upstream bootkube.
// https://github.com/kubernetes-incubator/bootkube/blob/master/cmd/bootkube/render.go
const (
	apiOffset  = 1
	dnsOffset  = 10
	etcdOffset = 15
)

// APIServiceIP picks a default IP from the given service CIDR range.
func APIServiceIP(serviceCIDR string) (net.IP, error) {
	return offsetServiceIP(serviceCIDR, apiOffset)
}

// DNSServiceIP picks a default IP from the given service CIDR range.
func DNSServiceIP(serviceCIDR string) (net.IP, error) {
	return offsetServiceIP(serviceCIDR, dnsOffset)
}

// EtcdServiceIP picks a default IP from the given service CIDR range.
func EtcdServiceIP(serviceCIDR string) (net.IP, error) {
	return offsetServiceIP(serviceCIDR, etcdOffset)
}

func offsetServiceIP(serviceCIDR string, offset int) (net.IP, error) {
	_, ipnet, err := net.ParseCIDR(serviceCIDR)
	if err != nil {
		return net.ParseIP(""), err
	}
	ip := make(net.IP, len(ipnet.IP))
	copy(ip, ipnet.IP)
	for i := 0; i < offset; i++ {
		incIPv4(ip)
	}
	if ipnet.Contains(ip) {
		return ip, nil
	}
	return net.IP([]byte("")), fmt.Errorf("Service IP %v is not in %s", ip, ipnet)
}

func incIPv4(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
