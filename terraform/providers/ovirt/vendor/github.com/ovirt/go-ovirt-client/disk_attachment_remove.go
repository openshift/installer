package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveDiskAttachment(vmID VMID, diskAttachmentID DiskAttachmentID, retries ...RetryStrategy) error {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	return retry(
		fmt.Sprintf("removing disk attachment %s on VM %s", diskAttachmentID, vmID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.
				SystemService().
				VmsService().
				VmService(string(vmID)).
				DiskAttachmentsService().
				AttachmentService(string(diskAttachmentID)).
				Remove().
				Send()
			return err
		},
	)
}

func (m *mockClient) RemoveDiskAttachment(vmID VMID, diskAttachmentID DiskAttachmentID, _ ...RetryStrategy) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	vm, ok := m.vmDiskAttachmentsByVM[vmID]
	if !ok {
		return newError(ENotFound, "VM %s doesn't exist", vmID)
	}

	diskAttachment, ok := vm[diskAttachmentID]
	if !ok {
		return newError(ENotFound, "Disk attachment %s not found on VM %s", diskAttachmentID, vmID)
	}

	delete(m.vmDiskAttachmentsByDisk, diskAttachment.DiskID())
	delete(m.vmDiskAttachmentsByVM[vmID], diskAttachmentID)

	return nil
}
