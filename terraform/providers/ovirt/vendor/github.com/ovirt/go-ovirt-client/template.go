package ovirtclient

import (
	ovirtsdk "github.com/ovirt/go-ovirt"
)

//go:generate go run scripts/rest/rest.go -i "Template" -n "template" -T "TemplateID"

// TemplateClient represents the portion of the client that deals with VM templates.
type TemplateClient interface {
	// CreateTemplate creates a new template from an existing VM.
	CreateTemplate(vmID VMID, name string, params OptionalTemplateCreateParameters, retries ...RetryStrategy) (
		Template,
		error,
	)
	// ListTemplates returns all templates stored in the oVirt engine.
	ListTemplates(retries ...RetryStrategy) ([]Template, error)
	// GetTemplateByName returns a template by its Name.
	GetTemplateByName(templateName string, retries ...RetryStrategy) (Template, error)
	// GetTemplate returns a template by its ID.
	GetTemplate(id TemplateID, retries ...RetryStrategy) (Template, error)
	// GetBlankTemplate finds a blank template in the oVirt engine and returns it. If no blank template is present,
	// this function will return an error.
	GetBlankTemplate(retries ...RetryStrategy) (Template, error)
	// RemoveTemplate removes the template with the specified ID.
	RemoveTemplate(templateID TemplateID, retries ...RetryStrategy) error
	// WaitForTemplateStatus waits for a template to enter a specific status.
	WaitForTemplateStatus(templateID TemplateID, status TemplateStatus, retries ...RetryStrategy) (Template, error)
	// CopyTemplateDiskToStorageDomain copies template disk to the specified storage domain.
	CopyTemplateDiskToStorageDomain(diskID DiskID, storageDomainID StorageDomainID, retries ...RetryStrategy) (Disk, error)
}

// TemplateID is an identifier for a template. It has a special type so the compiler
// can catch errors when the template ID is erroneously passed elsewhere.
type TemplateID string

// TemplateData is a set of prepared configurations for VMs.
type TemplateData interface {
	// ID returns the identifier of the template. This is typically a UUID.
	ID() TemplateID
	// Name is the human-readable name for the template.
	Name() string
	// Description is a longer description for the template.
	Description() string
	// Status returns the status of the template.
	Status() TemplateStatus
	// CPU returns the CPU configuration of the template if any.
	CPU() VMCPU

	// IsBlank returns true, if the template either has the ID of all zeroes, or if the template has no settings, disks,
	// or other settings. This function only checks the details supported by go-ovirt-client.
	IsBlank(...RetryStrategy) (bool, error)
}

// Template incorporates the TemplateData to provide access to data in a template, but also
// adds convenience functions to work with templates.
type Template interface {
	TemplateData

	// WaitForStatus waits for a template to enter a specific status. It returns the updated
	// template as a result.
	WaitForStatus(status TemplateStatus, retries ...RetryStrategy) (Template, error)
	// ListDiskAttachments lists all disk attachments for the current template.
	ListDiskAttachments(retries ...RetryStrategy) ([]TemplateDiskAttachment, error)
	// Remove removes the specified template.
	Remove(retries ...RetryStrategy) error
}

// TemplateStatus represents the status the template is in.
type TemplateStatus string

const (
	// TemplateStatusOK indicates that the template is ready and can be used.
	TemplateStatusOK TemplateStatus = "ok"
	// TemplateStatusLocked means that an operation is taking place on the template and cannot
	// be currently modified.
	TemplateStatusLocked TemplateStatus = "locked"
	// TemplateStatusIllegal indicates that the template is invalid and cannot be used.
	TemplateStatusIllegal TemplateStatus = "illegal"
)

// OptionalTemplateCreateParameters contains the optional parameters for creating a template.
type OptionalTemplateCreateParameters interface {
	Description() *string
}

// BuildableTemplateCreateParameters is a buildable version of OptionalTemplateCreateParameters.
type BuildableTemplateCreateParameters interface {
	OptionalTemplateCreateParameters

	// WithDescription sets the description of a template.
	WithDescription(description string) (BuildableTemplateCreateParameters, error)
	// MustWithDescription is identical to WithDescription, but panics instead of returning an error.
	MustWithDescription(description string) BuildableTemplateCreateParameters
}

type templateCreateParameters struct {
	description *string
}

func (t templateCreateParameters) Description() *string {
	return t.description
}

func (t templateCreateParameters) WithDescription(description string) (BuildableTemplateCreateParameters, error) {
	t.description = &description
	return t, nil
}

func (t templateCreateParameters) MustWithDescription(description string) BuildableTemplateCreateParameters {
	builder, err := t.WithDescription(description)
	if err != nil {
		panic(err)
	}
	return builder
}

// TemplateCreateParams creates a builder for the parameters of the template creation.
func TemplateCreateParams() BuildableTemplateCreateParameters {
	return &templateCreateParameters{}
}

func convertSDKTemplate(sdkTemplate *ovirtsdk.Template, client Client) (Template, error) {
	id, ok := sdkTemplate.Id()
	if !ok {
		return nil, newError(EFieldMissing, "template does not contain ID")
	}
	name, ok := sdkTemplate.Name()
	if !ok {
		return nil, newError(EFieldMissing, "template does not contain a name")
	}
	description, ok := sdkTemplate.Description()
	if !ok {
		return nil, newError(EFieldMissing, "template does not contain a description")
	}
	status, ok := sdkTemplate.Status()
	if !ok {
		return nil, newFieldNotFound("template", "status")
	}
	cpu, err := convertSDKTemplateCPU(sdkTemplate)
	if err != nil {
		return nil, err
	}
	return &template{
		client:      client,
		id:          TemplateID(id),
		name:        name,
		status:      TemplateStatus(status),
		description: description,
		cpu:         cpu,
	}, nil
}

func convertSDKTemplateCPU(sdkObject *ovirtsdk.Template) (*vmCPU, error) {
	sdkCPU, ok := sdkObject.Cpu()
	if !ok {
		return nil, newFieldNotFound("VM", "CPU")
	}
	cpuTopo, ok := sdkCPU.Topology()
	if !ok {
		return nil, newFieldNotFound("CPU in VM", "CPU topo")
	}
	cores, ok := cpuTopo.Cores()
	if !ok {
		return nil, newFieldNotFound("CPU topo in CPU in VM", "cores")
	}
	threads, ok := cpuTopo.Threads()
	if !ok {
		return nil, newFieldNotFound("CPU topo in CPU in VM", "threads")
	}
	sockets, ok := cpuTopo.Sockets()
	if !ok {
		return nil, newFieldNotFound("CPU topo in CPU in VM", "sockets")
	}
	cpu := &vmCPU{
		topo: &vmCPUTopo{
			uint(cores),
			uint(threads),
			uint(sockets),
		},
	}
	return cpu, nil
}

type template struct {
	client      Client
	id          TemplateID
	name        string
	description string
	status      TemplateStatus
	cpu         *vmCPU
}

func (t template) ListDiskAttachments(retries ...RetryStrategy) ([]TemplateDiskAttachment, error) {
	return t.client.ListTemplateDiskAttachments(t.id, retries...)
}

func (t template) CPU() VMCPU {
	return t.cpu
}

func (t template) Status() TemplateStatus {
	return t.status
}

func (t template) WaitForStatus(status TemplateStatus, retries ...RetryStrategy) (Template, error) {
	return t.client.WaitForTemplateStatus(t.id, status, retries...)
}

func (t template) Remove(retries ...RetryStrategy) error {
	return t.client.RemoveTemplate(t.id, retries...)
}

func (t template) IsBlank(retries ...RetryStrategy) (bool, error) {
	if t.cpu.topo.sockets != 1 || t.cpu.topo.cores != 1 || t.cpu.topo.threads != 1 {
		return false, nil
	}

	attachments, err := t.client.ListTemplateDiskAttachments(t.id, retries...)
	if err != nil {
		return false, wrap(err, EUnidentified, "failed to list template disk attachments")
	}

	if len(attachments) != 0 {
		return false, nil
	}

	// Assume it's a blank template. Further factors may be added later.
	return true, nil
}

func (t template) ID() TemplateID {
	return t.id
}

func (t template) Name() string {
	return t.name
}

func (t template) Description() string {
	return t.description
}
