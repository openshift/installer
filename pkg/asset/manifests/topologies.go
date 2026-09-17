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
		// Count control plane nodes as schedulable nodes to determine infrastructure topology.
		numOfWorkers += controlPlaneReplicas
	}

	// After accounting for schedulable control planes, determine infrastructure topology.
	// SNO (1 CP + 0 workers) becomes numOfWorkers=1 → Single
	// All other cases (2+ total schedulable nodes) → HA
	if numOfWorkers == 1 {
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
