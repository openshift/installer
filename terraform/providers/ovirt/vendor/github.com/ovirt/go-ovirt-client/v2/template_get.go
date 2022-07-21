// Code generated automatically using go:generate. DO NOT EDIT.

package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetTemplate(id TemplateID, retries ...RetryStrategy) (result Template, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting template %s", id),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().TemplatesService().TemplateService(string(id)).Get().Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Template()
			if !ok {
				return newError(
					ENotFound,
					"no template returned when getting template ID %s",
					id,
				)
			}
			result, err = convertSDKTemplate(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert template %s",
					id,
				)
			}
			return nil
		})
	return
}

func (m *mockClient) GetTemplate(id TemplateID, _ ...RetryStrategy) (Template, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if item, ok := m.templates[id]; ok {
		return item, nil
	}
	return nil, newError(ENotFound, "template with ID %s not found", id)
}
