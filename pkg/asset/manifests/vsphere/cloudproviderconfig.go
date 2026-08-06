package vsphere

import (
	"bytes"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	cloudconfig "github.com/openshift/library-go/pkg/cloudprovider/vsphere"
)

const (
	regionTagCategory = "openshift-region"
	zoneTagCategory   = "openshift-zone"
)

func printIfNotEmpty(buf *bytes.Buffer, k, v string) {
	if v != "" {
		fmt.Fprintf(buf, "%s = %q\n", k, v)
	}
}

// setNodes sets Nodes section in vsphere-cloud-provider config according passed VSpherePlatformNodeNetworking spec.
func setNodes(cfg *cloudconfig.CPIConfig, nodeNetworking *configv1.VSpherePlatformNodeNetworking) {
	cfg.Nodes.ExternalVMNetworkName = nodeNetworking.External.Network
	cfg.Nodes.ExternalNetworkSubnetCIDR = strings.Join(nodeNetworking.External.NetworkSubnetCIDR, ",")
	cfg.Nodes.ExcludeExternalNetworkSubnetCIDR = strings.Join(nodeNetworking.External.ExcludeNetworkSubnetCIDR, ",")

	cfg.Nodes.InternalVMNetworkName = nodeNetworking.Internal.Network
	cfg.Nodes.InternalNetworkSubnetCIDR = strings.Join(nodeNetworking.Internal.NetworkSubnetCIDR, ",")
	cfg.Nodes.ExcludeInternalNetworkSubnetCIDR = strings.Join(nodeNetworking.Internal.ExcludeNetworkSubnetCIDR, ",")
}

// CloudProviderConfigYaml generates the yaml out of tree cloud provider config for the vSphere platform.
func CloudProviderConfigYaml(infraID string, ic *installconfig.InstallConfig) (string, error) {
	p := ic.Config.Platform.VSphere
	vCenters := make(map[string]*cloudconfig.VirtualCenterConfig)

	for _, vCenter := range p.VCenters {
		vCenterConfig := cloudconfig.VirtualCenterConfig{
			VCenterIP:   vCenter.Server,
			Datacenters: vCenter.Datacenters,
			// We are setting this in global so lets remove from here
			//InsecureFlag: true,
		}
		// Only set port if it is configured in the install-config.  infrastructure/cluster will do the same which results
		// in the 3CMO not setting port if not configured in infrastructure/cluster.
		if vCenter.Port != 0 {
			vCenterConfig.VCenterPort = uint(vCenter.Port)
		}
		vCenters[vCenter.Server] = &vCenterConfig
	}

	cloudProviderConfig := cloudconfig.CPIConfig{CommonConfig: cloudconfig.CommonConfig{
		Global: cloudconfig.Global{
			SecretName:      "vsphere-creds", // #nosec G101 -- this is the name of a Kubernetes secret, not a credential
			SecretNamespace: "kube-system",
			InsecureFlag:    true,
		},
		Vcenter: vCenters,
	}}

	if len(p.FailureDomains) > 1 {
		cloudProviderConfig.Labels = &cloudconfig.Labels{
			Zone:   vspheretypes.TagCategoryZone,
			Region: vspheretypes.TagCategoryRegion,
		}
	}

	// Populate the nodes section with networking information mirroring the logic
	// in GetInfraPlatformSpec. If nodeNetworking is explicitly set in the
	// install-config, use it directly; otherwise fall back to the machine
	// network CIDRs (which should encompass the VIPs).
	if p.NodeNetworking != nil {
		setNodes(&cloudProviderConfig, p.NodeNetworking)
	} else {
		var cidrs []string
		for _, machineNetwork := range ic.Config.MachineNetwork {
			cidrs = append(cidrs, machineNetwork.CIDR.String())
		}
		nodeNetworking := &configv1.VSpherePlatformNodeNetworking{
			External: configv1.VSpherePlatformNodeNetworkingSpec{
				NetworkSubnetCIDR: cidrs,
			},
			Internal: configv1.VSpherePlatformNodeNetworkingSpec{
				NetworkSubnetCIDR: cidrs,
			},
		}
		setNodes(&cloudProviderConfig, nodeNetworking)
	}

	cloudProviderConfigYaml, err := yaml.Marshal(cloudProviderConfig)
	if err != nil {
		return "", err
	}
	return string(cloudProviderConfigYaml), nil
}

// CloudProviderConfigIni generates the multi-zone ini cloud provider config
// for the vSphere platform. folderPath is the absolute path to the VM folder that will be
// used for installation. p is the vSphere platform struct.
func CloudProviderConfigIni(infraID string, p *vspheretypes.Platform) (string, error) {
	buf := new(bytes.Buffer)

	fmt.Fprintln(buf, "[Global]")
	printIfNotEmpty(buf, "secret-name", "vsphere-creds")
	printIfNotEmpty(buf, "secret-namespace", "kube-system")
	printIfNotEmpty(buf, "insecure-flag", "1")
	fmt.Fprintln(buf, "")

	for _, vcenter := range p.VCenters {
		fmt.Fprintf(buf, "[VirtualCenter %q]\n", vcenter.Server)
		if vcenter.Port != 0 {
			printIfNotEmpty(buf, "port", fmt.Sprintf("%d", vcenter.Port))
			fmt.Fprintln(buf, "")
		}
		datacenters := make([]string, 0, len(vcenter.Datacenters))
		datacenters = append(datacenters, vcenter.Datacenters...)
		for _, failureDomain := range p.FailureDomains {
			if failureDomain.Server == vcenter.Server {
				failureDomainDatacenter := failureDomain.Topology.Datacenter
				exists := false
				for _, existingDatacenter := range datacenters {
					if failureDomainDatacenter == existingDatacenter {
						exists = true
						break
					}
				}
				if !exists {
					datacenters = append(datacenters, failureDomainDatacenter)
				}
			}
		}
		printIfNotEmpty(buf, "datacenters", strings.Join(datacenters, ","))
	}
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "[Workspace]")
	printIfNotEmpty(buf, "server", p.FailureDomains[0].Server)
	printIfNotEmpty(buf, "datacenter", p.FailureDomains[0].Topology.Datacenter)
	printIfNotEmpty(buf, "default-datastore", p.FailureDomains[0].Topology.Datastore)

	folderPath := fmt.Sprintf("/%s/vm/%s", p.FailureDomains[0].Topology.Datacenter, infraID)
	if p.FailureDomains[0].Topology.Folder != "" {
		folderPath = p.FailureDomains[0].Topology.Folder
	}
	printIfNotEmpty(buf, "folder", folderPath)
	printIfNotEmpty(buf, "resourcepool-path", p.FailureDomains[0].Topology.ResourcePool)
	fmt.Fprintln(buf, "")

	if len(p.FailureDomains) > 1 {
		fmt.Fprintln(buf, "[Labels]")
		printIfNotEmpty(buf, "region", regionTagCategory)
		printIfNotEmpty(buf, "zone", zoneTagCategory)
	}

	return buf.String(), nil
}
