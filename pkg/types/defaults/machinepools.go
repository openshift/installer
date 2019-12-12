package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
	gcpdefaults "github.com/openshift/installer/pkg/types/gcp/defaults"
	"github.com/openshift/installer/pkg/types/libvirt"
)

// SetMachinePoolDefaults sets the defaults for the machine pool.
func SetMachinePoolDefaults(p *types.MachinePool, platform string) {
	defaultReplicaCount := int64(3)
	if platform == libvirt.Name {
		defaultReplicaCount = 1
	}
	if p.Replicas == nil {
		p.Replicas = &defaultReplicaCount
	}
	if p.Hyperthreading == "" {
		p.Hyperthreading = types.HyperthreadingEnabled
	}
	switch platform {
	case gcp.Name:
		gcpdefaults.SetMachinePoolDefaults(p)
	default:
	}
}
