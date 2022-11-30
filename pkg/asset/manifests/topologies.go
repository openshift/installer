package manifests

import (
	configv1 "github.com/openshift/api/config/v1"

	"github.com/openshift/installer/pkg/types"
)

// determineTopologies determines the Infrastructure CR's
// infrastructureTopology and controlPlaneTopology given an install config file
func determineTopologies(installConfig *types.InstallConfig) (controlPlaneTopology configv1.TopologyMode, infrastructureTopology configv1.TopologyMode) {
	if installConfig.ControlPlane.Replicas != nil && *installConfig.ControlPlane.Replicas < 3 {
		controlPlaneTopology = configv1.SingleReplicaTopologyMode
	} else {
		controlPlaneTopology = configv1.HighlyAvailableTopologyMode
	}

	numOfWorkers := int64(0)
	for _, mp := range installConfig.Compute {
		if mp.Replicas != nil {
			numOfWorkers += *mp.Replicas
		}
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
