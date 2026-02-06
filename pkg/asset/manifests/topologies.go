package manifests

import (
	"k8s.io/utils/ptr"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types"
)

// determineTopologies determines the Infrastructure CR's
// infrastructureTopology and controlPlaneTopology given an install config file
func determineTopologies(installConfig *types.InstallConfig) (controlPlaneTopology configv1.TopologyMode, infrastructureTopology configv1.TopologyMode) {
	controlPlaneReplicas := ptr.Deref(installConfig.ControlPlane.Replicas, 3)
	switch controlPlaneReplicas {
	case 1:
		controlPlaneTopology = configv1.SingleReplicaTopologyMode
	case 2:
		controlPlaneTopology = configv1.DualReplicaTopologyMode
	default:
		controlPlaneTopology = configv1.HighlyAvailableTopologyMode
	}

	if controlPlaneReplicas >= 2 && installConfig.Arbiter != nil && ptr.Deref(installConfig.Arbiter.Replicas, 0) != 0 {
		controlPlaneTopology = configv1.HighlyAvailableArbiterMode
	}

	numOfWorkers := int64(0)
	for _, mp := range installConfig.Compute {
		numOfWorkers += ptr.Deref(mp.Replicas, 0)
	}
	if numOfWorkers < 2 {
		// Control planes are schedulable when there are < 2 workers.
		// Adjust the number of schedulable nodes here to reflect.
		numOfWorkers += controlPlaneReplicas
	}

	switch numOfWorkers {
	case 1:
		infrastructureTopology = configv1.SingleReplicaTopologyMode
		// When there are less than 2 worker nodes, the masters are made scheduable,
		// effectively making the number of worker nodes greater than 1.
		// Setting infrastructureTopology accordingly.
		if controlPlaneTopology == configv1.HighlyAvailableTopologyMode {
			infrastructureTopology = configv1.HighlyAvailableTopologyMode
		}
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
