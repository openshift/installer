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

	// The day-1 compute replica count is not a reliable signal on its own, since the
	// control plane may be schedulable and workers may be added after the install. Only
	// a single-node control plane with no workers is guaranteed to stay single-replica.
	if controlPlaneTopology == configv1.SingleReplicaTopologyMode && numOfWorkers == 0 {
		infrastructureTopology = configv1.SingleReplicaTopologyMode
	} else {
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
