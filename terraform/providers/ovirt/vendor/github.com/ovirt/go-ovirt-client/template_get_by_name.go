package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) GetTemplateByName(templateName string, retries ...RetryStrategy) (result Template, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("getting template by Name %s", templateName),
		o.logger,
		retries,
		func() error {
			response, err := o.conn.SystemService().TemplatesService().List().Search("name=" + templateName).Send()
			if err != nil {
				return err
			}
			for _, sdkObject := range response.MustTemplates().Slice() {
				if mTemplate, ok := sdkObject.Name(); ok {
					if templateName == mTemplate {
						result, err = convertSDKTemplate(sdkObject, o)
						if err != nil {
							return wrap(
								err,
								EBug,
								"failed to convert Template %s",
								templateName,
							)
						}
						return nil

					}
				}
			}
			return newError(ENotFound, "template with Name %s not found", templateName)
		})
	return result, err
}

func (m *mockClient) GetTemplateByName(templateName string, _ ...RetryStrategy) (result Template, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, template := range m.templates {
		if template.name == templateName {
			return template, nil
		}
	}
	return nil, newError(ENotFound, "template with Name %s not found", templateName)
}
