package ovirtclient //nolint:dupl

func (o *oVirtClient) ListVMs(retries ...RetryStrategy) (result []VM, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []VM{}
	err = retry(
		"listing vms",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().VmsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Vms()
			if !ok {
				return nil
			}
			result = make([]VM, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKVM(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert vm during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListVMs(_ ...RetryStrategy) ([]VM, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]VM, len(m.vms))
	i := 0
	for _, item := range m.vms {
		result[i] = item
		i++
	}
	return result, nil
}
