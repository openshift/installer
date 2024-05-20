package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveVM(id VMID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("removing VM %s", id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().VmsService().VmService(string(id)).Remove().Send()
			if err != nil {
				return err
			}
			return nil
		})
	return
}

func (m *mockClient) RemoveVM(id VMID, retries ...RetryStrategy) error {

	retries = defaultRetries(retries, defaultWriteTimeouts(m))

	return retry(
		fmt.Sprintf("removing VM %s", id),
		m.logger,
		retries,
		func() error {
			m.lock.Lock()
			defer m.lock.Unlock()

			if _, ok := m.vms[id]; !ok {
				return newError(ENotFound, "VM with ID %s not found", id)
			}

			for _, diskAttachment := range m.vmDiskAttachmentsByVM[id] {
				if m.disks[diskAttachment.DiskID()].status == DiskStatusLocked {
					return newError(EConflict, "Cannot delete VM, disk %s is locked.", diskAttachment.DiskID())
				}
				delete(m.disks, diskAttachment.DiskID())
				delete(m.vmDiskAttachmentsByDisk, diskAttachment.DiskID())
			}
			for nicID, nic := range m.nics {
				if nic.VMID() == id {
					delete(m.nics, nicID)
				}
			}
			delete(m.vmIPs, id)
			delete(m.vmDiskAttachmentsByVM, id)
			delete(m.graphicsConsolesByVM, id)
			delete(m.vms, id)

			return nil
		})
}
