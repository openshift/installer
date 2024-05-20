package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) WaitForVMStatus(id VMID, status VMStatus, retries ...RetryStrategy) (vm VM, err error) {
	retries = defaultRetries(retries, defaultLongTimeouts(o))
	err = retry(
		fmt.Sprintf("waiting for VM %s status %s", id, status),
		o.logger,
		retries,
		func() error {
			vm, err = o.GetVM(id, retries...)
			if err != nil {
				return err
			}
			if vm.Status() != status {
				return newError(EPending, "VM status is %s, not %s", vm.Status(), status)
			}
			return nil
		})
	return
}

func (m *mockClient) WaitForVMStatus(id VMID, status VMStatus, retries ...RetryStrategy) (vm VM, err error) {
	retries = defaultRetries(retries, defaultLongTimeouts(m))
	err = retry(
		fmt.Sprintf("waiting for VM %s status %s", id, status),
		m.logger,
		retries,
		func() error {
			vm, err = m.GetVM(id, retries...)
			if err != nil {
				return err
			}
			if vm.Status() != status {
				return newError(EPending, "VM status is %s, not %s", vm.Status(), status)
			}
			return nil
		})
	return
}
