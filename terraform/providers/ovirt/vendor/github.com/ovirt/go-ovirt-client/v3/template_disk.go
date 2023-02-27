package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

// TemplateDiskClient contains the methods to work with template disk attachments.
type TemplateDiskClient interface {
	// ListTemplateDiskAttachments lists all disk attachments for a template.
	ListTemplateDiskAttachments(templateID TemplateID, retries ...RetryStrategy) ([]TemplateDiskAttachment, error)
}

// TemplateDiskAttachmentData contains the methods to get the details of a disk attached to a template.
type TemplateDiskAttachmentData interface {
	// ID returns the identifier of the attachment.
	ID() TemplateDiskAttachmentID
	// TemplateID returns the ID of the template the disk is attached to.
	TemplateID() TemplateID
	// DiskID returns the ID of the disk in this attachment.
	DiskID() DiskID
	// DiskInterface describes the means by which a disk will appear to the VM.
	DiskInterface() DiskInterface
	// Bootable defines whether the disk is bootable
	Bootable() bool
	// Active defines whether the disk is active in the virtual machine it’s attached to.
	Active() bool
}

// TemplateDiskAttachment contains all methods from TemplateDiskAttachmentData and also convenience functions to fetch
// and work with template disk attachments.
type TemplateDiskAttachment interface {
	TemplateDiskAttachmentData

	// Template fetches the template object this disk attachment is related to.
	Template(retries ...RetryStrategy) (Template, error)
	// Disk fetches the disk this attachment attaches.
	Disk(retries ...RetryStrategy) (Disk, error)
}

// TemplateDiskAttachmentID is a typed string to ensure that these IDs are only used for template disk attachments.
type TemplateDiskAttachmentID string

type templateDiskAttachment struct {
	client Client

	id            TemplateDiskAttachmentID
	templateID    TemplateID
	diskID        DiskID
	diskInterface DiskInterface
	bootable      bool
	active        bool
}

func (t templateDiskAttachment) ID() TemplateDiskAttachmentID {
	return t.id
}

func (t templateDiskAttachment) TemplateID() TemplateID {
	return t.templateID
}

func (t templateDiskAttachment) DiskID() DiskID {
	return t.diskID
}

func (t templateDiskAttachment) DiskInterface() DiskInterface {
	return t.diskInterface
}

func (t templateDiskAttachment) Bootable() bool {
	return t.bootable
}

func (t templateDiskAttachment) Active() bool {
	return t.active
}

func (t templateDiskAttachment) Template(retries ...RetryStrategy) (Template, error) {
	return t.client.GetTemplate(t.templateID, retries...)
}

func (t templateDiskAttachment) Disk(retries ...RetryStrategy) (Disk, error) {
	return t.client.GetDisk(t.diskID, retries...)
}

func convertSDKTemplateDiskAttachment(attachment *ovirtsdk.DiskAttachment, o *oVirtClient) (TemplateDiskAttachment, error) {
	id := attachment.MustId()
	disk, ok := attachment.Disk()
	if !ok {
		return nil, newFieldNotFound("template disk attachment", "disk")
	}
	diskID, ok := disk.Id()
	if !ok {
		return nil, newFieldNotFound("disk on template disk attachment", "disk")
	}
	template, ok := attachment.Template()
	if !ok {
		return nil, newFieldNotFound("template disk attachment", "template")
	}
	templateID, ok := template.Id()
	if !ok {
		return nil, newFieldNotFound("template on template disk attachment", "id")
	}
	diskInterface, ok := attachment.Interface()
	if !ok {
		return nil, newFieldNotFound("template disk attachment", "interface")
	}
	bootable, ok := attachment.Bootable()
	if !ok {
		return nil, newFieldNotFound("template disk attachment", "bootable")
	}
	active, ok := attachment.Bootable()
	if !ok {
		return nil, newFieldNotFound("template disk attachment", "active")
	}

	return &templateDiskAttachment{
		o,

		TemplateDiskAttachmentID(id),
		TemplateID(templateID),
		DiskID(diskID),
		DiskInterface(diskInterface),
		bootable,
		active,
	}, nil
}
