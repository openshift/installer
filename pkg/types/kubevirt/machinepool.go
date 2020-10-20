package kubevirt

// MachinePool stores the configuration for a machine pool installed
// on kubevirt.
type MachinePool struct {
	// CPU is the amount of CPUs used.
	// +optional
	CPU uint32 `json:"cpu,omitempty"`

	// Memory is the size of a VM's memory.
	// Format: https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/api/resource/quantity.go
	// +optional
	Memory string `json:"memory,omitempty"`

	// StorageSize is the size of VM's boot volume.
	// Format: https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/api/resource/quantity.go
	// +optional
	StorageSize string `json:"storageSize,omitempty"`
}

// Set sets the values from `required` to `p`.
func (p *MachinePool) Set(required *MachinePool) {
	if required == nil || p == nil {
		return
	}

	if required.CPU != 0 {
		p.CPU = required.CPU
	}

	if required.Memory != "" {
		p.Memory = required.Memory
	}

	if required.StorageSize != "" {
		p.StorageSize = required.StorageSize
	}
}
