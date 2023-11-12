package vsphere

import (
	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
)

// GetInfraPlatformSpec constructs VSpherePlatformSpec for the infrastructure spec
func GetInfraPlatformSpec(ic *installconfig.InstallConfig) *configv1.VSpherePlatformSpec {
	var platformSpec configv1.VSpherePlatformSpec
	icPlatformSpec := ic.Config.VSphere

	for _, vcenter := range icPlatformSpec.VCenters {
		platformSpec.VCenters = append(platformSpec.VCenters, configv1.VSpherePlatformVCenterSpec{
			Server:      vcenter.Server,
			Port:        vcenter.Port,
			Datacenters: vcenter.Datacenters,
		})
	}

	for _, failureDomain := range icPlatformSpec.FailureDomains {
		topology := failureDomain.Topology
		if topology.ComputeCluster != "" && topology.Networks[0] != "" {
			platformSpec.FailureDomains = append(platformSpec.FailureDomains, configv1.VSpherePlatformFailureDomainSpec{
				Name:   failureDomain.Name,
				Region: failureDomain.Region,
				Zone:   failureDomain.Zone,
				Server: failureDomain.Server,
				Topology: configv1.VSpherePlatformTopology{
					Datacenter:     topology.Datacenter,
					ComputeCluster: topology.ComputeCluster,
					Networks:       topology.Networks,
					Datastore:      topology.Datastore,
					ResourcePool:   topology.ResourcePool,
					Folder:         topology.Folder,
				},
			})
		}
	}
	return &platformSpec
}
