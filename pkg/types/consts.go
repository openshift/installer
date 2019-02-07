package types

// MachineRole is the role for a machine
type MachineRole string

const (
	// ControlPlaneMachineRole is used for machines that comprise the control plane
	ControlPlaneMachineRole MachineRole = "control-plane"
	// ComputeMachineRole is used for machines that run work loads
	ComputeMachineRole MachineRole = "compute"
)
