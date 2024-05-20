package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ListDisksByAlias(alias string, retries ...RetryStrategy) (result []Disk, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Disk{}
	err = retry(
		fmt.Sprintf("listing disk by alias %s", alias),
		o.logger,
		retries,
		func() error {
			searchString := fmt.Sprintf("name=%s", alias)
			response, e := o.conn.SystemService().DisksService().List().Search(searchString).Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Disks()
			if !ok {
				return nil
			}
			result = make([]Disk, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKDisk(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert disk during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListDisksByAlias(alias string, _ ...RetryStrategy) ([]Disk, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Disk, 0)
	for _, d := range m.disks {
		if d.alias == alias {
			result = append(result, d)
		}
	}
	return result, nil
}
