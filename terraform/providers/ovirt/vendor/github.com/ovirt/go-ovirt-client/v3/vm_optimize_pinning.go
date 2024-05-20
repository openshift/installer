package ovirtclient

import "fmt"

func (o *oVirtClient) AutoOptimizeVMCPUPinningSettings(id VMID, optimize bool, retries ...RetryStrategy) error {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	return retry(
		fmt.Sprintf("optimizing CPU pinning settings for VM %s", id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().
				VmsService().
				VmService(string(id)).
				AutoPinCpuAndNumaNodes().
				OptimizeCpuSettings(optimize).
				Send()
			return err
		})
}

func (m *mockClient) AutoOptimizeVMCPUPinningSettings(_ VMID, _ bool, _ ...RetryStrategy) error {
	// This function cannot be simulated as the VM object does not contain any observable return values apart from the
	// NUMA nodes being moved around. If you know of a way please add a mock and add a test for it.
	return nil
}
