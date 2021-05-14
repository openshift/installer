package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
)

func setMachinePool(p *types.MachinePool) {
	if p.Platform.Ovirt == nil {
		p.Platform.Ovirt = &ovirt.MachinePool{}
	}
}

func setDefaultAffinityGroups(p *ovirt.Platform, mp *types.MachinePool, agName string) {
	if mp.Platform.Ovirt.AffinityGroupsNames == nil {
		for _, ag := range p.AffinityGroups {
			if ag.Name == agName {
				mp.Platform.Ovirt.AffinityGroupsNames = []string{agName}
			}
		}
	}
}

// SetControlPlaneDefaults sets the defaults for the ControlPlane Machines.
func SetControlPlaneDefaults(p *ovirt.Platform, mp *types.MachinePool) {
	setMachinePool(mp)
	setDefaultAffinityGroups(p, mp, DefaultControlPlaneAffinityGroupName)
}

// SetComputeDefaults sets the defaults for the Compute Machines.
func SetComputeDefaults(p *ovirt.Platform, mp *types.MachinePool) {
	setMachinePool(mp)
	setDefaultAffinityGroups(p, mp, DefaultComputeAffinityGroupName)
}
