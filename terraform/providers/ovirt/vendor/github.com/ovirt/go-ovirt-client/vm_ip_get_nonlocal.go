package ovirtclient

import "net"

func (m *mockClient) GetVMNonLocalIPAddresses(id VMID, retries ...RetryStrategy) (
	result map[string][]net.IP,
	err error,
) {
	return waitForIPAddresses(id, nonLocalIPSearchParams, retries, m.logger, m)
}

func (o *oVirtClient) GetVMNonLocalIPAddresses(id VMID, retries ...RetryStrategy) (map[string][]net.IP, error) {
	return waitForIPAddresses(id, nonLocalIPSearchParams, retries, o.logger, o)
}
