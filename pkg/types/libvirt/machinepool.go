package libvirt

// MachinePool stores the configuration for a machine pool installed
// on libvirt.
type MachinePool struct {
	// DomainMemoryMiB is the amount of RAM each VM will have.
	DomainMemoryMiB int `json:"memory"`

	// DomainVcpCount is the number of VCPUs each VM will have.
	DomainVcpuCount int `json:"vcpus"`
}

// Set sets the values from `required` to `a`.
func (l *MachinePool) Set(required *MachinePool) {
	if required == nil || l == nil {
		return
	}

	if required.DomainMemoryMiB != 0 {
		l.DomainMemoryMiB = required.DomainMemoryMiB
	}
	if required.DomainVcpuCount != 0 {
		l.DomainVcpuCount = required.DomainVcpuCount
	}
}
