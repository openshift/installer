package packet

// MachinePool stores the configuration for a machine pool installed
// on packet.
type MachinePool struct {
}

// Disk defines a BM disk
type Disk struct {
	// SizeGB size of the bootable disk in GiB.
	SizeGB int64 `json:"sizeGB"`
}

// CPU defines the BM cpu, made of (Sockets * Cores).
type CPU struct {
	// Sockets is the number of sockets for a BM.
	// Total CPUs is (Sockets * Cores)
	Sockets int32 `json:"sockets"`

	// Cores is the number of cores per socket.
	// Total CPUs is (Sockets * Cores)
	Cores int32 `json:"cores"`
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}
}
