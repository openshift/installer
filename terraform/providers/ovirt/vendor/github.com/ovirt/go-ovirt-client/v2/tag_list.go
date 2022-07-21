// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListTags(retries ...RetryStrategy) (result []Tag, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Tag{}
	err = retry(
		"listing tags",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().TagsService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Tags()
			if !ok {
				return nil
			}
			result = make([]Tag, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKTag(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert tag during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListTags(_ ...RetryStrategy) ([]Tag, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Tag, len(m.tags))
	i := 0
	for _, item := range m.tags {
		result[i] = item
		i++
	}
	return result, nil
}
