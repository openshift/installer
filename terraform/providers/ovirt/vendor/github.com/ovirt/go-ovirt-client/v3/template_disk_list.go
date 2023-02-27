package ovirtclient

import (
	"fmt"
)

func (o *oVirtClient) ListTemplateDiskAttachments(
	templateID TemplateID,
	retries ...RetryStrategy,
) (result []TemplateDiskAttachment, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	err = retry(
		fmt.Sprintf("listing disk attachments for template %s", templateID),
		o.logger,
		retries,
		func() error {
			res, err := o.conn.
				SystemService().
				TemplatesService().
				TemplateService(string(templateID)).
				DiskAttachmentsService().
				List().
				Send()
			if err != nil {
				return err
			}
			attachments, ok := res.Attachments()
			if !ok {
				return newFieldNotFound("template disk attachments result", "attachments")
			}
			result = make([]TemplateDiskAttachment, len(attachments.Slice()))
			for i, attachment := range attachments.Slice() {
				result[i], err = convertSDKTemplateDiskAttachment(attachment, o)
				if err != nil {
					return wrap(err, EBug, "failed to convert template disk attachment %d for template %s (%v)")
				}
			}
			return nil
		},
	)
	return result, err
}

func (m *mockClient) ListTemplateDiskAttachments(
	templateID TemplateID,
	_ ...RetryStrategy,
) ([]TemplateDiskAttachment, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	result := make([]TemplateDiskAttachment, len(m.templateDiskAttachmentsByTemplate[templateID]))
	i := 0
	for _, attachment := range m.templateDiskAttachmentsByTemplate[templateID] {
		result[i] = attachment
		i++
	}
	return result, nil
}
