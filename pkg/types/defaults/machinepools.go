package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/version"
)

// SetMachinePoolDefaults sets the defaults for the machine pool.
func SetMachinePoolDefaults(p *types.MachinePool, platform string) {
	defaultReplicaCount := int64(3)
	if platform == libvirt.Name {
		defaultReplicaCount = 1
	}
	if p.Name == types.MachinePoolEdgeRoleName {
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
}

// CreateEdgeMachinePoolDefaults create the edge compute pool when it is not already defined.
func CreateEdgeMachinePoolDefaults(pools []types.MachinePool, platform string, replicas int64) *types.MachinePool {
	edgePoolDefined := false
	for _, compute := range pools {
		if compute.Name == types.MachinePoolEdgeRoleName {
			edgePoolDefined = true
		}
	}
	if edgePoolDefined {
		return nil
	}
	pool := &types.MachinePool{
		Name:     types.MachinePoolEdgeRoleName,
		Replicas: &replicas,
	}
	SetMachinePoolDefaults(pool, platform)
	return pool
}
