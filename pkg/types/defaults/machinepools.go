package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	awsdefaults "github.com/openshift/installer/pkg/types/aws/defaults"
	"github.com/openshift/installer/pkg/types/azure"
	azuredefaults "github.com/openshift/installer/pkg/types/azure/defaults"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpdefaults "github.com/openshift/installer/pkg/types/gcp/defaults"
	"github.com/openshift/installer/pkg/version"
)

// SetMachinePoolDefaults sets the defaults for the machine pool.
func SetMachinePoolDefaults(p *types.MachinePool, platform *types.Platform) {
	defaultReplicaCount := int64(3)
	if p.Name == types.MachinePoolEdgeRoleName || p.Name == types.MachinePoolArbiterRoleName {
		defaultReplicaCount = 0
	}
	if p.Replicas == nil {
		p.Replicas = &defaultReplicaCount
	}
	if p.Hyperthreading == "" {
		p.Hyperthreading = types.HyperthreadingEnabled
	}
	if p.Architecture == "" {
		p.Architecture = version.DefaultArch()
	}

	switch platform.Name() {
	case aws.Name:
		if p.Platform.AWS == nil && platform.AWS.DefaultMachinePlatform != nil {
			p.Platform.AWS = &aws.MachinePool{}
		}
		awsdefaults.Apply(platform.AWS.DefaultMachinePlatform, p.Platform.AWS)
		awsdefaults.SetMachinePoolDefaults(p.Platform.AWS, p.Name)
	case azure.Name:
		if p.Platform.Azure == nil && platform.Azure.DefaultMachinePlatform != nil {
			p.Platform.Azure = &azure.MachinePool{}
		}
		azuredefaults.Apply(platform.Azure.DefaultMachinePlatform, p.Platform.Azure)
	case gcp.Name:
		gcpdefaults.SetMachinePoolDefaults(platform, p.Platform.GCP)
	default:
	}
}

// hasEdgePoolConfig checks if the Edge compute pool has been defined on install-config.
func hasEdgePoolConfig(pools []types.MachinePool) bool {
	edgePoolDefined := false
	for _, compute := range pools {
		if compute.Name == types.MachinePoolEdgeRoleName {
			edgePoolDefined = true
		}
	}
	return edgePoolDefined
}

// CreateEdgeMachinePoolDefaults create the edge compute pool when it is not already defined.
func CreateEdgeMachinePoolDefaults(pools []types.MachinePool, platform *types.Platform, replicas int64) *types.MachinePool {
	if hasEdgePoolConfig(pools) {
		return nil
	}
	pool := &types.MachinePool{
		Name:     types.MachinePoolEdgeRoleName,
		Replicas: &replicas,
	}
	SetMachinePoolDefaults(pool, platform)
	return pool
}
