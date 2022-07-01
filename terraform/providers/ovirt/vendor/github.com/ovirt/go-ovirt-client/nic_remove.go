package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveNIC(vmid VMID, id NICID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("removing NIC %s from VM %s", id, vmid),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().VmsService().VmService(string(vmid)).NicsService().NicService(string(id)).Remove().Send()
			if err != nil {
				return err
			}
			return nil
		})
	return
}

func (m *mockClient) RemoveNIC(vmid VMID, id NICID, _ ...RetryStrategy) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.vms[vmid]; !ok {
		return newError(ENotFound, "NIC with ID %s not found", vmid)
	}
	if _, ok := m.nics[id]; !ok {
		return newError(ENotFound, "NIC with ID %s not found on VM with ID %s", id, vmid)
	}
	delete(m.nics, id)
	return nil
}
