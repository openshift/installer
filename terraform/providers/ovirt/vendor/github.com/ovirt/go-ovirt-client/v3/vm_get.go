package ovirtclient //nolint:dupl

import (
	"fmt"
)

func (o *oVirtClient) GetVM(id VMID, retries ...RetryStrategy) (result VM, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting vm %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().VmsService().VmService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Vm()
			if !ok {
				return newError(
					ENotFound,
					"no vm returned when getting vm ID %s",
					id,
				)
			}
			result, err = convertSDKVM(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert vm %s",
					id,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetVM(id VMID, _ ...RetryStrategy) (VM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.vms[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "vm with ID %s not found", id)
}
