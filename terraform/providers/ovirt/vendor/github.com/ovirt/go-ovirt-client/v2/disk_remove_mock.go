package ovirtclient

func (m *mockClient) RemoveDisk(diskID DiskID, _ ...RetryStrategy) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.disks[diskID]; !ok {
		return newError(ENotFound, "disk with ID %s not found", diskID)
	}

	// Check if disk is attached to a running VM
	if diskAttachment, ok := m.vmDiskAttachmentsByDisk[diskID]; ok {
		vm := m.vms[diskAttachment.vmid]
		if vm.status != VMStatusDown {
			return newError(
				EConflict,
				"Disk %s is attached to VM %s, which is \"%s\" not \"%s\".",
				diskID,
				vm.id,
				vm.status,
				VMStatusDown,
			)
		}
	}
	// Check if disk is attached to a template.
	if _, ok := m.templateDiskAttachmentsByDisk[diskID]; ok {
		return newError(EUnidentified, "Cannot remove disk attached to a template. Please specify storage domain to remove from.")
	}

	if diskAttachment, ok := m.vmDiskAttachmentsByDisk[diskID]; ok {
		vm := m.vms[diskAttachment.vmid]
		delete(m.vmDiskAttachmentsByVM[vm.id], diskAttachment.id)
	}

	delete(m.vmDiskAttachmentsByDisk, diskID)
	delete(m.disks, diskID)

	return nil
}
