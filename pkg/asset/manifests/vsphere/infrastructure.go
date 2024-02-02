package vsphere

import (
	"fmt"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
)

// GetInfraPlatformSpec constructs VSpherePlatformSpec for the infrastructure spec
func GetInfraPlatformSpec(ic *installconfig.InstallConfig, clusterID string) *configv1.VSpherePlatformSpec {
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
			template := topology.Template
			if len(template) == 0 {
				template = fmt.Sprintf("/%s/vm/%s-rhcos-%s-%s", topology.Datacenter, clusterID, failureDomain.Region, failureDomain.Zone)
			}

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
					Template:       template,
				},
			})
		}
	}

	platformSpec.APIServerInternalIPs = types.StringsToIPs(icPlatformSpec.APIVIPs)
	platformSpec.IngressIPs = types.StringsToIPs(icPlatformSpec.IngressVIPs)
	platformSpec.MachineNetworks = types.MachineNetworksToCIDRs(ic.Config.MachineNetwork)

	return &platformSpec
}
