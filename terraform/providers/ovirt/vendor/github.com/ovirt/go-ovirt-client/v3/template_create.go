package ovirtclient

import (
	"fmt"
	"time"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateTemplate(
	vmID VMID,
	name string,
	params OptionalTemplateCreateParameters,
	retries ...RetryStrategy,
) (result Template, err error) {
	retries = defaultRetries(retries, defaultReadTimeouts(o))
	if params == nil {
		params = &templateCreateParameters{}
	}
	err = retry(
		fmt.Sprintf("creating template from VM %s", vmID),
		o.logger,
		retries,
		func() error {
			tpl := ovirtsdk.NewTemplateBuilder()
			tpl.VmBuilder(ovirtsdk.NewVmBuilder().Id(string(vmID)))
			tpl.Name(name)
			if desc := params.Description(); desc != nil {
				tpl.Description(*desc)
			}
			response, err := o.conn.SystemService().TemplatesService().Add().Template(tpl.MustBuild()).Send()
			if err != nil {
				return err
			}
			sdkObject, ok := response.Template()
			if !ok {
				return newError(
					ENotFound,
					"no template returned after creating template from VM %s",
					vmID,
				)
			}
			result, err = convertSDKTemplate(sdkObject, o)
			if err != nil {
				return wrap(
					err,
					EBug,
					"failed to convert template from VM %s",
					vmID,
				)
			}
			return nil
		})
	return result, err
}

func (m *mockClient) CreateTemplate(
	vmID VMID,
	name string,
	params OptionalTemplateCreateParameters,
	_ ...RetryStrategy,
) (Template, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	vm, ok := m.vms[vmID]
	if !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", vmID)
	}

	for _, tpl := range m.templates {
		if tpl.name == name {
			return nil, newError(EConflict, "A template with the name \"%s\" already exists.", name)
		}
	}

	if params == nil {
		params = &templateCreateParameters{}
	}

	description := ""
	if desc := params.Description(); desc != nil {
		description = *desc
	}
	tpl := &template{
		client:      m,
		id:          TemplateID(m.GenerateUUID()),
		name:        name,
		description: description,
		status:      TemplateStatusLocked,
		cpu:         vm.cpu.clone(),
	}
	m.templates[tpl.ID()] = tpl
	m.templateDiskAttachmentsByTemplate[tpl.ID()] = make(
		[]*templateDiskAttachment,
		len(m.vmDiskAttachmentsByVM[vmID]),
	)
	m.attachTemplateDisks(vmID, tpl)

	go m.handlePostTemplateCreation(tpl)
	return tpl, nil
}

func (m *mockClient) handlePostTemplateCreation(tpl *template) {
	func() {
		time.Sleep(2 * time.Second)
		m.lock.Lock()
		defer m.lock.Unlock()
		if tpl.status == TemplateStatusIllegal {
			return
		}
		for _, attachment := range m.templateDiskAttachmentsByTemplate[tpl.id] {
			disk := m.disks[attachment.diskID]
			disk.Unlock()
		}
		tpl.status = TemplateStatusOK
	}()
}

func (m *mockClient) attachTemplateDisks(vmID VMID, tpl *template) {
	i := 0
	for _, attachment := range m.vmDiskAttachmentsByVM[vmID] {
		disk := m.disks[attachment.diskID]
		newDisk := disk.clone(nil)
		_ = newDisk.Lock()
		newDisk.alias = fmt.Sprintf("disk-%s", generateRandomID(5, m.nonSecureRandom))
		m.disks[newDisk.ID()] = newDisk

		tplAttachment := &templateDiskAttachment{
			client:        m,
			id:            TemplateDiskAttachmentID(m.GenerateUUID()),
			templateID:    tpl.ID(),
			diskID:        newDisk.ID(),
			diskInterface: attachment.diskInterface,
			bootable:      attachment.bootable,
			active:        attachment.active,
		}
		m.templateDiskAttachmentsByDisk[newDisk.id] = tplAttachment
		m.templateDiskAttachmentsByTemplate[tpl.id][i] = tplAttachment
		i++
	}
}
