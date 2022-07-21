// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListDisks(retries ...RetryStrategy) (result []Disk, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Disk{}
	err = retry(
		"listing disks",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().DisksService().List().Send()
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

func (m *mockClient) ListDisks(_ ...RetryStrategy) ([]Disk, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Disk, len(m.disks))
	i := 0
	for _, item := range m.disks {
		result[i] = item
		i++
	}
	return result, nil
}
