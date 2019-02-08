package types

// MachineRole is the role for a machine
type MachineRole string

const (
	// ControlPlaneMachineRole is used for machines that comprise the control plane
	ControlPlaneMachineRole MachineRole = "control-plane"
	// ComputeMachineRole is used for machines that run work loads
	ComputeMachineRole MachineRole = "compute"
)

// ClusterAPIMachineRole returns the machine role used by clusterapi.
func (r MachineRole) ClusterAPIMachineRole() string {
	switch r {
	case ControlPlaneMachineRole:
		return "master"
	case ComputeMachineRole:
		return "worker"
	default:
		return ""
	}
}

// MachineConfigOperatorMachineRole returns the machine role used by machine-config-operator.
func (r MachineRole) MachineConfigOperatorMachineRole() string {
	switch r {
	case ControlPlaneMachineRole:
		return "master"
	case ComputeMachineRole:
		return "worker"
	default:
		return ""
	}
}
