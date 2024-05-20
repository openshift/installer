package ovirtclient

import (
	"fmt"
)

//nolint:dupl
func (o *oVirtClient) ListDiskAttachments(
	vmid VMID,
	retries ...RetryStrategy,
) (result []DiskAttachment, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []DiskAttachment{}
	err = retry(
		fmt.Sprintf("listing disk attachments on VM %s", vmid),
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().VmsService().VmService(string(vmid)).DiskAttachmentsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Attachments()
			if !ok {
				return nil
			}
			result = make([]DiskAttachment, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKDiskAttachment(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert disk during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListDiskAttachments(vmID VMID, _ ...RetryStrategy) ([]DiskAttachment, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	diskAttachments, ok := m.vmDiskAttachmentsByVM[vmID]
	if !ok {
		return nil, newError(ENotFound, "VM %s doesn't exist", vmID)
	}

	result := make([]DiskAttachment, len(diskAttachments))
	i := 0
	for _, attachment := range diskAttachments {
		result[i] = attachment
		i++
	}

	return result, nil
}
