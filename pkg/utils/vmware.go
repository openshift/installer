package utils

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"strings"

	"github.com/sirupsen/logrus"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"

	machinev1beta1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
)

// ConstructNetworkKargsFromMachine does something.
func ConstructNetworkKargsFromMachine(claims []ipamv1.IPAddressClaim, addresses []ipamv1.IPAddress, machine *machinev1beta1.Machine, network machinev1beta1.NetworkDeviceSpec) (string, error) {
	var ipAddresses []string
	var gateways []string
	for idx := range network.AddressesFromPools {
		for _, address := range addresses {
			logrus.Debugf("Checking IPAdress %v.  Does it match? %v", address.Name, fmt.Sprintf("%s-claim-%d-%d", machine.Name, 0, idx))
			if address.Name == fmt.Sprintf("%s-claim-%d-%d", machine.Name, 0, idx) {
				ipAddresses = append(ipAddresses, fmt.Sprintf("%v/%v", address.Spec.Address, address.Spec.Prefix))
				gateways = append(gateways, address.Spec.Gateway)
				break
			}
		}
	}
	return ConstructKargsFromNetworkConfig(ipAddresses, network.Nameservers, gateways)
}

func getSubnetMask(prefix netip.Prefix) (string, error) {
	prefixLength := net.IPv4len * 8
	if prefix.Addr().Is6() {
		prefixLength = net.IPv6len * 8
	}
	ipMask := net.CIDRMask(prefix.Masked().Bits(), prefixLength)
	maskBytes, err := hex.DecodeString(ipMask.String())
	if err != nil {
		return "", err
	}
	ip := net.IP(maskBytes)
	maskStr := ip.To16().String()
	return maskStr, nil
}

// ConstructKargsForBootstrap constructs the kargs string for a bootstrap node.
func ConstructKargsForBootstrap(installConfig *types.InstallConfig) (string, error) {
	for _, host := range installConfig.VSphere.Hosts {
		if host.Role != "bootstrap" {
			continue
		}
		return ConstructKargsFromNetworkConfig(host.NetworkDevice.IPAddrs, host.NetworkDevice.Nameservers, []string{host.NetworkDevice.Gateway})
	}
	return "", errors.New("unable to find host with bootstrap role")
}

// ConstructKargsFromNetworkConfig constructs the kargs string from the network configuration.
func ConstructKargsFromNetworkConfig(ipAddrs []string, nameservers []string, gateways []string) (string, error) {
	outKargs := ""

	for index, address := range ipAddrs {
		var gatewayIP netip.Addr
		gateway := gateways[index]
		if len(gateway) > 0 {
			ip, err := netip.ParseAddr(gateway)
			if err != nil {
				return "", err
			}
			if ip.Is6() {
				gateway = fmt.Sprintf("[%s]", gateway)
			}
			gatewayIP = ip
		}

		prefix, err := netip.ParsePrefix(address)
		if err != nil {
			return "", err
		}
		var ipStr, gatewayStr, maskStr string
		addr := prefix.Addr()
		switch {
		case addr.Is6():
			maskStr = fmt.Sprintf("%d", prefix.Bits())
			ipStr = fmt.Sprintf("[%s]", addr.String())
			if len(gateway) > 0 && gatewayIP.Is6() {
				gatewayStr = gateway
			}
		case addr.Is4():
			maskStr, err = getSubnetMask(prefix)
			if err != nil {
				return "", err
			}
			if len(gateway) > 0 && gatewayIP.Is4() {
				gatewayStr = gateway
			}
			ipStr = addr.String()
		default:
			return "", errors.New("IP address must adhere to IPv4 or IPv6 format")
		}
		outKargs += fmt.Sprintf("ip=%s::%s:%s:::none ", ipStr, gatewayStr, maskStr)
	}

	for _, nameserver := range nameservers {
		ip := net.ParseIP(nameserver)
		if ip.To4() == nil {
			nameserver = fmt.Sprintf("[%s]", nameserver)
		}
		outKargs += fmt.Sprintf("nameserver=%s ", nameserver)
	}

	outKargs = strings.Trim(outKargs, " ")
	logrus.Debugf("Generated karg: [%v].", outKargs)
	return outKargs, nil
}
