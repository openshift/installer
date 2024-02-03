package types

import configv1 "github.com/openshift/api/config/v1"

// StringsToIPs is used to convert list of strings to list of IP addresses.
func StringsToIPs(ips []string) []configv1.IP {
	res := []configv1.IP{}

	if ips == nil {
		return res
	}

	for _, ip := range ips {
		res = append(res, configv1.IP(ip))
	}

	return res
}

// MachineNetworksToCIDRs is used to convert list of Machine Network Entries to
// list of CIDRs.
func MachineNetworksToCIDRs(nets []MachineNetworkEntry) []configv1.CIDR {
	res := []configv1.CIDR{}

	if nets == nil {
		return res
	}

	for _, net := range nets {
		res = append(res, configv1.CIDR(net.CIDR.String()))
	}

	return res
}
