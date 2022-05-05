package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ListNICs(vmid VMID, retries ...RetryStrategy) (result []NIC, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("listing NICs for VM %s", vmid),
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().VmsService().VmService(string(vmid)).NicsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Nics()
			if !ok {
				return nil
			}
			result = make([]NIC, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKNIC(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert NIC during listing item #%d", i)
				}
			}
			return nil
		},
	)
	return
}

func (m *mockClient) ListNICs(vmid VMID, _ ...RetryStrategy) ([]NIC, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	var result []NIC
	for _, item := range m.nics {
		if item.vmid == vmid {
			result = append(result, item)
		}
	}
	return result, nil
}
