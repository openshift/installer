package ovirt

// MachinePool stores the configuration for a machine pool installed
// on ovirt.
type MachinePool struct {
	// InstanceTypeID defines the VM instance type and overrides
	// the hardware parameters of the created VM, including cpu and memory.
	// If InstanceTypeID is passed, all memory and cpu variables will be ignored.
	InstanceTypeID string `json:"instanceTypeID,omitempty"`

	// CPU defines the VM CPU.
	// +optional
	CPU *CPU `json:"cpu,omitempty"`

	// MemoryMB is the size of a VM's memory in MiBs.
	// +optional
	MemoryMB int32 `json:"memoryMB,omitempty"`

	// OSDisk is the the root disk of the node.
	// +optional
	OSDisk *Disk `json:"osDisk,omitempty"`

	// VMType defines the workload type of the VM.
	// +kubebuilder:validation:Enum="";desktop;server;high_performance
	// +optional
	VMType VMType `json:"vmType,omitempty"`
}

// Disk defines a VM disk
type Disk struct {
	// SizeGB size of the bootable disk in GiB.
	SizeGB int64 `json:"sizeGB"`
}

// CPU defines the VM cpu, made of (Sockets * Cores).
type CPU struct {
	// Sockets is the number of sockets for a VM.
	// Total CPUs is (Sockets * Cores)
	Sockets int32 `json:"sockets"`

	// Cores is the number of cores per socket.
	// Total CPUs is (Sockets * Cores)
	Cores int32 `json:"cores"`
}

// VMType defines the type of the VM, which will change the VM configuration,
// like including or excluding devices (like excluding sound-card),
// device configuration (like using multi-queues for vNic), and several other
// configuration tweaks. This doesn't effect properties like CPU count and amount of memory.
type VMType string

const (
	// VMTypeDesktop set the VM type to desktop. Virtual machines optimized to act
	// as desktop machines do have a sound card, use an image (thin allocation),
	// and are stateless.
	VMTypeDesktop VMType = "desktop"
	// VMTypeServer sets the VM type to server. Virtual machines optimized to act
	// as servers have no sound card, use a cloned disk image, and are not stateless.
	VMTypeServer VMType = "server"
	// VMTypeHighPerformance sets a VM type to high_performance which sets various
	// properties of a VM to optimize for performance, like enabling headless mode,
	// disabling usb, smart-card, and sound devices, enabling host cpu pass-through,
	// multi-queues for vNics and several more items.
	// See https://www.ovirt.org/develop/release-management/features/virt/high-performance-vm.html.
	VMTypeHighPerformance VMType = "high_performance"
)

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.InstanceTypeID != "" {
		p.InstanceTypeID = required.InstanceTypeID
	}

	if required.VMType != "" {
		p.VMType = required.VMType
	}

	if required.CPU != nil {
		p.CPU = required.CPU
	}

	if required.MemoryMB != 0 {
		p.MemoryMB = required.MemoryMB
	}

	if required.OSDisk != nil {
		p.OSDisk = required.OSDisk
	}
}
