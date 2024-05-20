// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListDatacenters(retries ...RetryStrategy) (result []Datacenter, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Datacenter{}
	err = retry(
		"listing datacenters",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().DataCentersService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.DataCenters()
			if !ok {
				return nil
			}
			result = make([]Datacenter, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKDatacenter(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert datacenter during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListDatacenters(_ ...RetryStrategy) ([]Datacenter, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Datacenter, len(m.dataCenters))
	i := 0
	for _, item := range m.dataCenters {
		result[i] = item
		i++
	}
	return result, nil
}
