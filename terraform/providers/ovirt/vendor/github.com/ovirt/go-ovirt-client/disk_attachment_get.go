package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetDiskAttachment(
	vmid VMID,
	id DiskAttachmentID,
	retries ...RetryStrategy,
) (result DiskAttachment, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting disk attachment %s on VM %s", id, vmid),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.
				SystemService().
				VmsService().
				VmService(string(vmid)).
				DiskAttachmentsService().
				AttachmentService(string(id)).
				Get().
				Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Attachment()
			if !ok {
				return newFieldNotFound("disk attachment response", "attachment")
			}
			result, err = convertSDKDiskAttachment(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert disk attachment %s",
					id,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) GetDiskAttachment(vmID VMID, diskAttachmentID DiskAttachmentID, _ ...RetryStrategy) (DiskAttachment, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	vm, ok := m.vmDiskAttachmentsByVM[vmID]
	if !ok {
		return nil, newError(ENotFound, "VM %s doesn't exist", vmID)
	}

	diskAttachment, ok := vm[diskAttachmentID]
	if !ok {
		return nil, newError(ENotFound, "disk attachment %s not found on VM %s", diskAttachmentID, vmID)
	}

	return diskAttachment, nil
}
