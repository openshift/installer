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

	contolPlaneReplicas := ptr.Deref(installConfig.ControlPlane.Replicas, 100)
	if contolPlaneReplicas == 2 {
		controlPlaneTopology = configv1.DualReplicaTopologyMode

		if ptr.Deref(installConfig.Arbiter.Replicas, 0) != 0 {
			controlPlaneTopology = configv1.HighlyAvailableArbiterMode
		}
	} else if contolPlaneReplicas < 3 {
		controlPlaneTopology = configv1.SingleReplicaTopologyMode
	}

	numOfWorkers := int64(0)
	for _, mp := range installConfig.Compute {
		numOfWorkers += ptr.Deref(mp.Replicas, 0)
	}

	switch numOfWorkers {
	case 0:
		infrastructureTopology = controlPlaneTopology
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
