// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListStorageDomains(retries ...RetryStrategy) (result StorageDomainList, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []StorageDomain{}
	err = retry(
		"listing storage domains",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().StorageDomainsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.StorageDomains()
			if !ok {
				return nil
			}
			result = make([]StorageDomain, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKStorageDomain(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert storage domain during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListStorageDomains(_ ...RetryStrategy) (StorageDomainList, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]StorageDomain, len(m.storageDomains))
	i := 0
	for _, item := range m.storageDomains {
		result[i] = item
		i++
	}
	return result, nil
}
