package libvirt

// MachinePool stores the configuration for a machine pool installed
// on libvirt.
type MachinePool struct {
	// ImagePool is the name of the libvirt storage pool to which the storage
	// volume containing the OS image belongs.
	ImagePool string `json:"imagePool,omitempty"`
	// ImageVolume is the name of the libvirt storage volume containing the OS
	// image.
	ImageVolume string `json:"imageVolume,omitempty"`

	// Image is the URL to the OS image.
	// E.g. "http://aos-ostree.rhev-ci-vms.eng.rdu2.redhat.com/rhcos/images/cloud/latest/rhcos-qemu.qcow2.gz"
	Image string `json:"image"`
}

// Set sets the values from `required` to `a`.
func (l *MachinePool) Set(required *MachinePool) {
	if required == nil || l == nil {
		return
	}

	if required.ImagePool != "" {
		l.ImagePool = required.ImagePool
	}
	if required.ImageVolume != "" {
		l.ImageVolume = required.ImageVolume
	}
	if required.Image != "" {
		l.Image = required.Image
	}
}
