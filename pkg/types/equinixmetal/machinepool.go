package equinixmetal

// MachinePool stores the configuration for a machine pool installed
// on equinixmetal.
type MachinePool struct {
	// The Equinix Metal Plan defines the CPU, memory, and networking specs of the
	// provisioned node
	Plan string

	// CustomData is an arbitrary bit of json to make available within each
	// nodes metadata
	CustomData string

	// TODO(displague) Hardware reservation id?
	// TODO(displague) virtual network?
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}
}
