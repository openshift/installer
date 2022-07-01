// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListHosts(retries ...RetryStrategy) (result []Host, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Host{}
	err = retry(
		"listing hosts",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().HostsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Hosts()
			if !ok {
				return nil
			}
			result = make([]Host, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKHost(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert host during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListHosts(_ ...RetryStrategy) ([]Host, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Host, len(m.hosts))
	i := 0
	for _, item := range m.hosts {
		result[i] = item
		i++
	}
	return result, nil
}
