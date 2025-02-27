package vsphere

import (
	"fmt"

	"github.com/sirupsen/logrus"
	utilsnet "k8s.io/utils/net"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/api/features"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
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

			failureDomainSpec := configv1.VSpherePlatformFailureDomainSpec{
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
			}

			if ic.Config.EnabledFeatureGates().Enabled(features.FeatureGateVSphereHostVMGroupZonal) {
				logrus.Debug("Host VM Group based zonal feature gate enabled")

				if failureDomain.ZoneType == vsphere.HostGroupFailureDomain {
					vmGroupAndRuleName := fmt.Sprintf("%s-%s", clusterID, failureDomain.Name)
					failureDomainSpec.RegionAffinity = &configv1.VSphereFailureDomainRegionAffinity{
						Type: configv1.VSphereFailureDomainRegionType(failureDomain.RegionType),
					}
					failureDomainSpec.ZoneAffinity = &configv1.VSphereFailureDomainZoneAffinity{
						Type: configv1.VSphereFailureDomainZoneType(failureDomain.ZoneType),
						HostGroup: &configv1.VSphereFailureDomainHostGroup{
							HostGroup:  failureDomain.Topology.HostGroup,
							VMGroup:    vmGroupAndRuleName,
							VMHostRule: vmGroupAndRuleName,
						},
					}
				}
			}

			platformSpec.FailureDomains = append(platformSpec.FailureDomains, failureDomainSpec)
		}
	}

	platformSpec.APIServerInternalIPs = types.StringsToIPs(icPlatformSpec.APIVIPs)
	platformSpec.IngressIPs = types.StringsToIPs(icPlatformSpec.IngressVIPs)
	platformSpec.MachineNetworks = types.MachineNetworksToCIDRs(ic.Config.MachineNetwork)
	platformSpec.MachineNetworks = append(platformSpec.MachineNetworks, vipsToCIDRs(ic.Config.VSphere.APIVIPs)...)
	platformSpec.MachineNetworks = append(platformSpec.MachineNetworks, vipsToCIDRs(ic.Config.VSphere.IngressVIPs)...)

	if ic.Config.EnabledFeatureGates().Enabled(features.FeatureGateVSphereMultiNetworks) {
		logrus.Debug("Multi-networks feature gate enabled")
		if icPlatformSpec.NodeNetworking != nil {
			logrus.Debug("Multi-networks: node networking defined, copying to infrastructure spec")
			icPlatformSpec.NodeNetworking.DeepCopyInto(&platformSpec.NodeNetworking)
		} else {
			logrus.Debug("Multi-networks: node networking not defined, deriving from machineNetwork")
			var cidrs []string
			for _, machineNetwork := range ic.Config.MachineNetwork {
				cidrs = append(cidrs, machineNetwork.CIDR.String())
			}

			// if NodeNetworking is not defined, use the machine cidrs. the machine cidrs
			// should align with the VIP and should be a safe choice for inclusion in NodeNetworking.
			platformSpec.NodeNetworking.External.NetworkSubnetCIDR = cidrs
			platformSpec.NodeNetworking.Internal.NetworkSubnetCIDR = cidrs

			logrus.Debugf("Multi-networks appending cidrs: %v", cidrs)
		}
	}
	return &platformSpec
}

// vipsToCIDRs takes a single ip address and converts it to CIDR notation.
func vipsToCIDRs(vips []string) []configv1.CIDR {
	cidrs := make([]configv1.CIDR, len(vips))
	for i, vip := range vips {
		mask := "/32"
		if utilsnet.IsIPv6String(vip) {
			mask = "/128"
		}
		cidrs[i] = configv1.CIDR(vip + mask)
	}
	return cidrs
}
