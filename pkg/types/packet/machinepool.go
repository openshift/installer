package packet

// MachinePool stores the configuration for a machine pool installed
// on packet.
type MachinePool struct {
	// The Packet Plan defines the CPU, memory, and networking specs of the
	// provisioned node
	Plan string

	// TODO(displague) Hardware reservation id?
	// TODO(displague) virtual network?
	// TODO(displague) is userdata needed at this level?
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}
}
