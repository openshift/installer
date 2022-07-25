// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

func (o *oVirtClient) ListTemplates(retries ...RetryStrategy) (result []Template, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	result = []Template{}
	err = retry(
		"listing templates",
		o.logger,
		retries,
		func() error {
			response, e := o.conn.SystemService().TemplatesService().List().Send()
			if e != nil {
				return e
			}
			sdkObjects, ok := response.Templates()
			if !ok {
				return nil
			}
			result = make([]Template, len(sdkObjects.Slice()))
			for i, sdkObject := range sdkObjects.Slice() {
				result[i], e = convertSDKTemplate(sdkObject, o)
				if e != nil {
					return wrap(e, EBug, "failed to convert template during listing item #%d", i)
				}
			}
			return nil
		})
	return
}

func (m *mockClient) ListTemplates(_ ...RetryStrategy) ([]Template, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]Template, len(m.templates))
	i := 0
	for _, item := range m.templates {
		result[i] = item
		i++
	}
	return result, nil
}
