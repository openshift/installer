package powervs

import (
	"k8s.io/apimachinery/pkg/util/intstr"

	machinev1 "github.com/openshift/api/machine/v1"
)

// MachinePool stores the configuration for a machine pool installed on IBM Power VS.
type MachinePool struct {
	// VolumeIDs is the list of volumes attached to the instance.
	//
	// +optional
	VolumeIDs []string `json:"volumeIDs,omitempty"`

	// memoryGiB is the size of a virtual machine's memory, in GiB.
	//
	// +optional
	MemoryGiB int32 `json:"memoryGiB,omitempty"`

	// Processors defines the processing units for the instance.
	//
	// +optional
	Processors intstr.IntOrString `json:"processors,omitempty"`

	// ProcType defines the processor sharing model for the instance.
	// Must be one of {Capped, Dedicated, Shared}.
	//
	// +kubebuilder:validation:Enum:="Dedicated";"Shared";"Capped";""
	// +optional
	ProcType machinev1.PowerVSProcessorType `json:"procType,omitempty"`

	// SMTLevel specifies the level of SMT to set the control plane and worker nodes to.
	//
	// +optional
	SMTLevel string `json:"smtLevel,omitempty"`

	// SysType defines the system type for instance.
	//
	// +optional
	SysType string `json:"sysType,omitempty"`
}

// Set stores values from required into a
func (a *MachinePool) Set(required *MachinePool) {
	if required == nil || a == nil {
		return
	}
	if len(required.VolumeIDs) != 0 {
		a.VolumeIDs = required.VolumeIDs
	}
	if required.MemoryGiB != 0 {
		a.MemoryGiB = required.MemoryGiB
	}
	if required.Processors.StrVal != "" || required.Processors.IntVal != 0 {
		a.Processors = required.Processors
	}
	if required.ProcType != "" {
		a.ProcType = required.ProcType
	}
	if required.SMTLevel != "" {
		a.SMTLevel = required.SMTLevel
	}
	if required.SysType != "" {
		a.SysType = required.SysType
	}
}
