package ovirtclient //nolint:dupl

func (o *oVirtClient) ListInstanceTypes(retries ...RetryStrategy) (result []InstanceType, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []InstanceType{}
	err = retry(
		"listing instance types",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().InstanceTypesService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.InstanceType()
			if !ok {
				return nil
			}
			result = make([]InstanceType, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKInstanceType(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert instance type during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListInstanceTypes(_ ...RetryStrategy) ([]InstanceType, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]InstanceType, len(m.instanceTypes))
	i := 0
	for _, item := range m.instanceTypes {
		result[i] = item
		i++
	}
	return result, nil
}
