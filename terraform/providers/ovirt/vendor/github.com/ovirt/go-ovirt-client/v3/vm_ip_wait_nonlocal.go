package ovirtclient

import (
	"net"
	"regexp"
)

func mustParseCIDR(cidr string) net.IPNet {
	_, result, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return *result
}

var nonLocalIPSearchParams = NewVMIPSearchParams().
	WithExcludedRange(mustParseCIDR("0.0.0.0/32")).
	WithExcludedRange(mustParseCIDR("127.0.0.0/8")).
	WithExcludedRange(mustParseCIDR("169.254.0.0/15")).
	WithExcludedRange(mustParseCIDR("224.0.0.0/4")).
	WithExcludedRange(mustParseCIDR("255.255.255.255/32")).
	WithExcludedRange(mustParseCIDR("::/128")).
	WithExcludedRange(mustParseCIDR("::1/128")).
	WithExcludedRange(mustParseCIDR("fe80::/64")).
	WithExcludedRange(mustParseCIDR("ff00::/8")).
	WithExcludedInterface("lo").
	WithExcludedInterfacePattern(regexp.MustCompile("^dummy[0-9]+$"))

func (m *mockClient) WaitForNonLocalVMIPAddress(id VMID, retries ...RetryStrategy) (map[string][]net.IP, error) {
	return m.WaitForVMIPAddresses(id, nonLocalIPSearchParams, retries...)
}

func (o *oVirtClient) WaitForNonLocalVMIPAddress(id VMID, retries ...RetryStrategy) (map[string][]net.IP, error) {
	return o.WaitForVMIPAddresses(id, nonLocalIPSearchParams, retries...)
}
