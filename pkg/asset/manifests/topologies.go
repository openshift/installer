package manifests

import (
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

// determineTopologies determines the Infrastructure CR's
// infrastructureTopology and controlPlaneTopology given an install config file
func determineTopologies(installConfig *types.InstallConfig) (controlPlaneTopology configv1.TopologyMode, infrastructureTopology configv1.TopologyMode) {
	controlPlaneTopology = configv1.HighlyAvailableTopologyMode

	controlPlaneReplicas := ptr.Deref(installConfig.ControlPlane.Replicas, 3)
	if controlPlaneReplicas == 2 {
		controlPlaneTopology = configv1.DualReplicaTopologyMode

		if installConfig.Arbiter != nil && ptr.Deref(installConfig.Arbiter.Replicas, 0) != 0 {
			controlPlaneTopology = configv1.HighlyAvailableArbiterMode
		}
	} else if controlPlaneReplicas < 3 {
		controlPlaneTopology = configv1.SingleReplicaTopologyMode
	}

	numOfWorkers := int64(0)
	for _, mp := range installConfig.Compute {
		numOfWorkers += ptr.Deref(mp.Replicas, 0)
	}

	switch numOfWorkers {
	case 0:
		// Two node deployments with 0 workers mean that the control plane nodes are treated as workers
		// in that situation we have decided that it is appropriate to set the infrastructureTopology to HA.
		// All other configuration for different worker count are respected with the original intention.
		if controlPlaneReplicas == 2 {
			infrastructureTopology = configv1.HighlyAvailableTopologyMode
		} else {
			infrastructureTopology = controlPlaneTopology
		}
	case 1:
		infrastructureTopology = configv1.SingleReplicaTopologyMode
	default:
		infrastructureTopology = configv1.HighlyAvailableTopologyMode
	}

	return controlPlaneTopology, infrastructureTopology
}

func determineCPUPartitioning(installConfig *types.InstallConfig) configv1.CPUPartitioningMode {
	switch installConfig.CPUPartitioning {
	case types.CPUPartitioningAllNodes:
		return configv1.CPUPartitioningAllNodes
	default:
		return configv1.CPUPartitioningNone
	}
}
