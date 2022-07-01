package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ListVMTags(id VMID, retries ...RetryStrategy) (result []Tag, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("listing tags for vm %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().VmsService().VmService(string(id)).TagsService().List().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Tags()
			if !ok {
				return newError(
					ENotFound,
					"no tags returned when getting VM %s tags",
					id,
				)
			}
			result = make([]Tag, len(sdkObject.Slice()))
			for i, sdkTag := range sdkObject.Slice() {
				result[i], err = convertSDKTag(sdkTag, o)
				if err != nil {
					return err
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListVMTags(id VMID, _ ...RetryStrategy) (result []Tag, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.vms[id]; !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", id)
	}
	result = make([]Tag, len(m.vms[id].tagIDs))
	for i, tagID := range m.vms[id].tagIDs {
		result[i] = m.tags[tagID]
	}
	return result, nil
}
