package ovirtclient

import (
	"strings"

	ovirtsdk4 "github.com/ovirt/go-ovirt"
)

// DiskAttachmentID is the identifier for the disk attachment.
type DiskAttachmentID string

// DiskAttachmentClient contains the methods required for handling disk attachments.
type DiskAttachmentClient interface {
	// CreateDiskAttachment attaches a disk to a VM.
	CreateDiskAttachment(
		vmID VMID,
		diskID DiskID,
		diskInterface DiskInterface,
		params CreateDiskAttachmentOptionalParams,
		retries ...RetryStrategy,
	) (DiskAttachment, error)
	// GetDiskAttachment returns a single disk attachment in a virtual machine.
	GetDiskAttachment(vmID VMID, id DiskAttachmentID, retries ...RetryStrategy) (DiskAttachment, error)
	// ListDiskAttachments lists all disk attachments for a virtual machine.
	ListDiskAttachments(vmID VMID, retries ...RetryStrategy) ([]DiskAttachment, error)
	// RemoveDiskAttachment removes the disk attachment in question.
	RemoveDiskAttachment(vmID VMID, diskAttachmentID DiskAttachmentID, retries ...RetryStrategy) error
}

// DiskInterface describes the means by which a disk will appear to the VM.
type DiskInterface string

const (
	// DiskInterfaceIDE is a legacy controller device. Works with almost all guest operating systems, so it is good for
	// compatibility. Performance is lower than with the other alternatives.
	DiskInterfaceIDE DiskInterface = "ide"
	// DiskInterfaceSATA is a SATA controller device.
	DiskInterfaceSATA DiskInterface = "sata"
	// DiskInterfacesPAPRvSCSI is a para-virtualized device supported by the IBM pSeries family of machines, using the
	// SCSI protocol.
	DiskInterfacesPAPRvSCSI DiskInterface = "spapr_vscsi"
	// DiskInterfaceVirtIO is a virtualization interface where just the guest's device driver knows it is running in a
	// virtual environment. Enables guests to get high performance disk operations.
	DiskInterfaceVirtIO DiskInterface = "virtio"
	// DiskInterfaceVirtIOSCSI is a para-virtualized SCSI controller device. Fast interface with the guest via direct
	// physical storage device address, using the SCSI protocol.
	DiskInterfaceVirtIOSCSI DiskInterface = "virtio_scsi"
)

// DiskInterfaceList is a list of DiskInterface.
type DiskInterfaceList []DiskInterface

// DiskInterfaceValues returns all possible DiskInterface values.
func DiskInterfaceValues() DiskInterfaceList {
	return []DiskInterface{
		DiskInterfaceIDE,
		DiskInterfaceSATA,
		DiskInterfacesPAPRvSCSI,
		DiskInterfaceVirtIO,
		DiskInterfaceVirtIOSCSI,
	}
}

// Strings creates a string list of the values.
func (l DiskInterfaceList) Strings() []string {
	result := make([]string, len(l))
	for i, status := range l {
		result[i] = string(status)
	}
	return result
}

// Validate checks if the DiskInterface actually has a valid value.
func (d DiskInterface) Validate() error {
	for _, format := range DiskInterfaceValues() {
		if format == d {
			return nil
		}
	}
	return newError(
		EBadArgument,
		"invalid disk interface: %s must be one of: %s",
		d,
		strings.Join(DiskInterfaceValues().Strings(), ", "),
	)
}

// CreateDiskAttachmentOptionalParams are the optional parameters for creating a disk attachment.
type CreateDiskAttachmentOptionalParams interface {
	// Bootable defines whether the disk is bootable.
	Bootable() *bool

	// Active defines whether the disk is active in the virtual machine it’s attached to.
	Active() *bool
}

// BuildableCreateDiskAttachmentParams is a buildable version of CreateDiskAttachmentOptionalParams.
type BuildableCreateDiskAttachmentParams interface {
	CreateDiskAttachmentOptionalParams
	// WithBootable sets whether the disk is bootable.
	WithBootable(bootable bool) (BuildableCreateDiskAttachmentParams, error)
	// MustWithBootable is the same as WithBootable, but panics instead of returning an error.
	MustWithBootable(bootable bool) BuildableCreateDiskAttachmentParams

	// WithActive sets whether the disk is active is visible to the virtual machine or not. default is true
	WithActive(active bool) (BuildableCreateDiskAttachmentParams, error)
	// MustWithActive is the same as WithActive, but panics instead of returning an error.
	MustWithActive(active bool) BuildableCreateDiskAttachmentParams
}

// CreateDiskAttachmentParams creates a buildable set of parameters for creating a disk attachment.
func CreateDiskAttachmentParams() BuildableCreateDiskAttachmentParams {
	return &createDiskAttachmentParams{}
}

type createDiskAttachmentParams struct {
	bootable *bool
	active   *bool
}

func (c createDiskAttachmentParams) Bootable() *bool {
	return c.bootable
}

func (c createDiskAttachmentParams) Active() *bool {
	return c.active
}

func (c createDiskAttachmentParams) WithBootable(bootable bool) (BuildableCreateDiskAttachmentParams, error) {
	c.bootable = &bootable
	return c, nil
}

func (c createDiskAttachmentParams) MustWithBootable(bootable bool) BuildableCreateDiskAttachmentParams {
	builder, err := c.WithBootable(bootable)
	if err != nil {
		panic(err)
	}
	return builder
}

func (c createDiskAttachmentParams) WithActive(active bool) (BuildableCreateDiskAttachmentParams, error) {
	c.active = &active
	return c, nil
}

func (c createDiskAttachmentParams) MustWithActive(active bool) BuildableCreateDiskAttachmentParams {
	builder, err := c.WithActive(active)
	if err != nil {
		panic(err)
	}
	return builder
}

// DiskAttachment links together a Disk and a VM.
type DiskAttachment interface {
	// ID returns the identifier of the attachment.
	ID() DiskAttachmentID
	// VMID returns the ID of the virtual machine this attachment belongs to.
	VMID() VMID
	// DiskID returns the ID of the disk in this attachment.
	DiskID() DiskID
	// DiskInterface describes the means by which a disk will appear to the VM.
	DiskInterface() DiskInterface
	// Bootable defines whether the disk is bootable
	Bootable() bool
	// Active defines whether the disk is active in the virtual machine it’s attached to.
	Active() bool

	// VM fetches the virtual machine this attachment belongs to.
	VM(retries ...RetryStrategy) (VM, error)
	// Disk fetches the disk this attachment attaches.
	Disk(retries ...RetryStrategy) (Disk, error)

	// Remove removes the current disk attachment.
	Remove(retries ...RetryStrategy) error
}

type diskAttachment struct {
	client Client

	id            DiskAttachmentID
	vmid          VMID
	diskID        DiskID
	diskInterface DiskInterface
	active        bool
	bootable      bool
}

func (d *diskAttachment) DiskInterface() DiskInterface {
	return d.diskInterface
}

func (d *diskAttachment) Remove(retries ...RetryStrategy) error {
	return d.client.RemoveDiskAttachment(d.vmid, d.id, retries...)
}

func (d *diskAttachment) ID() DiskAttachmentID {
	return d.id
}

func (d *diskAttachment) VMID() VMID {
	return d.vmid
}

func (d *diskAttachment) DiskID() DiskID {
	return d.diskID
}

func (d *diskAttachment) Bootable() bool {
	return d.bootable
}

func (d *diskAttachment) Active() bool {
	return d.active
}

func (d *diskAttachment) VM(retries ...RetryStrategy) (VM, error) {
	return d.client.GetVM(d.vmid, retries...)
}

func (d *diskAttachment) Disk(retries ...RetryStrategy) (Disk, error) {
	return d.client.GetDisk(d.diskID, retries...)
}

func convertSDKDiskAttachment(object *ovirtsdk4.DiskAttachment, o *oVirtClient) (DiskAttachment, error) {
	id, ok := object.Id()
	if !ok {
		return nil, newFieldNotFound("disk attachment", "id")
	}
	vm, ok := object.Vm()
	if !ok {
		return nil, newFieldNotFound("disk attachment", "vm")
	}
	vmID, ok := vm.Id()
	if !ok {
		return nil, newFieldNotFound("vm on disk attachment", "id")
	}
	disk, ok := object.Disk()
	if !ok {
		return nil, newFieldNotFound("disk attachment", "disk")
	}
	diskID, ok := disk.Id()
	if !ok {
		return nil, newFieldNotFound("disk on disk attachment", "id")
	}
	diskInterface, ok := object.Interface()
	if !ok {
		return nil, newFieldNotFound("disk attachment", "disk interface")
	}
	bootable, ok := object.Bootable()
	if !ok {
		return nil, newFieldNotFound("bootable on disk attachment", "bootable")
	}
	active, ok := object.Active()
	if !ok {
		return nil, newFieldNotFound("active on disk attachment", "active")
	}
	return &diskAttachment{
		client: o,

		id:            DiskAttachmentID(id),
		vmid:          VMID(vmID),
		diskID:        DiskID(diskID),
		diskInterface: DiskInterface(diskInterface),
		bootable:      bootable,
		active:        active,
	}, nil
}
