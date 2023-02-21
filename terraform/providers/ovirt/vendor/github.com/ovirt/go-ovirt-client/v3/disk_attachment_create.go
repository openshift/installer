package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) CreateDiskAttachment(
	vmID VMID,
	diskID DiskID,
	diskInterface DiskInterface,
	params CreateDiskAttachmentOptionalParams,
	retries ...RetryStrategy,
) (result DiskAttachment, err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts(o))
	if err := diskInterface.Validate(); err != nil {
		return nil, wrap(err, EBadArgument, "failed to create disk attachment")
	}
	err = retry(
		fmt.Sprintf("attaching disk %s to vm %s", diskID, vmID),
		o.logger,
		retries,
		func() error {
			attachmentBuilder := ovirtsdk.NewDiskAttachmentBuilder()
			attachmentBuilder.Disk(ovirtsdk.NewDiskBuilder().Id(string(diskID)).MustBuild())
			attachmentBuilder.Interface(ovirtsdk.DiskInterface(diskInterface))
			attachmentBuilder.Vm(ovirtsdk.NewVmBuilder().Id(string(vmID)).MustBuild())
			if params != nil {
				if active := params.Active(); active != nil {
					attachmentBuilder.Active(*active)
				}
				if bootable := params.Bootable(); bootable != nil {
					attachmentBuilder.Bootable(*params.Bootable())
				}
			}
			attachment := attachmentBuilder.MustBuild()

			addRequest := o.conn.SystemService().VmsService().VmService(string(vmID)).DiskAttachmentsService().Add()
			addRequest.Attachment(attachment)
			response, err := addRequest.Send()
			if err != nil {
				return wrap(
					err,
					EUnidentified,
					"failed to attach disk %s to VM %s using %s",
					diskID,
					vmID,
					diskInterface,
				)
			}

			attachment, ok := response.Attachment()
			if !ok {
				return newFieldNotFound("attachment response", "attachment")
			}
			result, err = convertSDKDiskAttachment(attachment, o)
			if err != nil {
				return wrap(err, EUnidentified, "failed to convert SDK disk attachment")
			}
			return nil
		},
	)
	return result, err
}

func (m *mockClient) CreateDiskAttachment(
	vmID VMID,
	diskID DiskID,
	diskInterface DiskInterface,
	params CreateDiskAttachmentOptionalParams,
	_ ...RetryStrategy,
) (DiskAttachment, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if err := diskInterface.Validate(); err != nil {
		return nil, wrap(err, EBadArgument, "failed to create disk attachment")
	}

	vm, ok := m.vms[vmID]
	if !ok {
		return nil, newError(ENotFound, "VM with ID %s not found", vmID)
	}

	disk, ok := m.disks[diskID]
	if !ok {
		return nil, newError(ENotFound, "disk with ID %s not found", diskID)
	}

	attachment := &diskAttachment{
		client:        m,
		id:            DiskAttachmentID(m.GenerateUUID()),
		vmid:          vm.ID(),
		diskID:        disk.ID(),
		diskInterface: diskInterface,
	}
	if params != nil {
		if bootable := params.Bootable(); bootable != nil {
			attachment.bootable = *bootable
		}
		if active := params.Active(); active != nil {
			attachment.active = *active
		}
	}
	for _, diskAttachment := range m.vmDiskAttachmentsByVM[vm.ID()] {
		if diskAttachment.DiskID() == diskID {
			return nil, newError(EConflict, "disk %s is already attached to VM %s", diskID, vmID)
		}
	}

	if diskAttachment, ok := m.vmDiskAttachmentsByDisk[disk.ID()]; ok {
		return nil, newError(
			EConflict,
			"cannot attach disk %s to VM %s, already attached to VM %s",
			diskID,
			vmID,
			diskAttachment.VMID(),
		)
	}

	m.vmDiskAttachmentsByDisk[disk.ID()] = attachment
	m.vmDiskAttachmentsByVM[vm.ID()][attachment.ID()] = attachment

	return attachment, nil
}
