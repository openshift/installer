package vsphere

// MachinePool stores the configuration for a machine pool installed
// on vSphere.
type MachinePool struct {
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}
}
