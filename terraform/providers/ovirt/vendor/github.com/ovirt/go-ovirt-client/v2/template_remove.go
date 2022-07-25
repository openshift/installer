package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) RemoveTemplate(
	templateID TemplateID,
	retries ...RetryStrategy,
) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	err = retry(
		fmt.Sprintf("removing template %s", templateID),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().TemplatesService().TemplateService(string(templateID)).Remove().Send()
			return err
		})
	return
}

func (m *mockClient) RemoveTemplate(id TemplateID, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(m))
	err = retry(
		fmt.Sprintf("removing template %s", id),
		m.logger,
		retries,
		func() error {
			m.lock.Lock()
			defer m.lock.Unlock()
			tpl, ok := m.templates[id]
			if !ok {
				return newError(ENotFound, "Template with ID %s was not found", id)
			}

			for _, vm := range m.vms {
				if vm.templateID == id {
					return newError(
						EConflict,
						"Template %s cannot be removed because it is in use by VM %s.",
						id,
						vm.id,
					)
				}
			}

			if tpl.status == TemplateStatusLocked {
				return newError(EConflict, "Template %s is in status %s.", id, tpl.status)
			}

			delete(m.templates, id)
			return nil
		})
	return err
}
