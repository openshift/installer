package libvirt

// MachinePool stores the configuration for a machine pool installed
// on libvirt.
type MachinePool struct {
}

// Set sets the values from `required` to `a`.
func (l *MachinePool) Set(required *MachinePool) {
	if required == nil || l == nil {
		return
	}
}
