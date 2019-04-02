package defaults

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// SetPlatformDefaults sets the defaults for the platform.
func SetPlatformDefaults(p *vsphere.Platform, installConfig *types.InstallConfig) {
	if p.Workspace.Server == "" {
		if len(p.VirtualCenters) == 1 {
			p.Workspace.Server = p.VirtualCenters[0].Name
		}
	}
	if p.Workspace.Datacenter == "" {
		for _, vc := range p.VirtualCenters {
			if p.Workspace.Server == vc.Name {
				if len(vc.Datacenters) == 1 {
					p.Workspace.Datacenter = vc.Datacenters[0]
				}
				break
			}
		}
	}
	if p.Workspace.Folder == "" {
		p.Workspace.Folder = installConfig.ObjectMeta.Name
	}
	if p.SCSIControllerType == "" {
		p.SCSIControllerType = "pvscsi"
	}
}
