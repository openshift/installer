package vsphere

import (
	"bytes"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

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

// CloudProviderConfigYaml generates the yaml out of tree cloud provider config for the vSphere platform.
func CloudProviderConfigYaml(infraID string, p *vspheretypes.Platform) (string, error) {
	vCenters := make(map[string]*cloudconfig.VirtualCenterConfig)

	for _, vCenter := range p.VCenters {
		vCenterPort := int32(443)
		if vCenter.Port != 0 {
			vCenterPort = vCenter.Port
		}
		vCenterConfig := cloudconfig.VirtualCenterConfig{
			VCenterIP:   vCenter.Server,
			VCenterPort: uint(vCenterPort),
			Datacenters: vCenter.Datacenters,
			// We are setting this in global so lets remove from here
			//InsecureFlag: true,
		}
		vCenters[vCenter.Server] = &vCenterConfig
	}

	cloudProviderConfig := cloudconfig.CommonConfig{
		Global: cloudconfig.Global{
			SecretName:      "vsphere-creds", // #nosec G101 -- this is the name of a Kubernetes secret, not a credential
			SecretNamespace: "kube-system",
			InsecureFlag:    true,
		},
		Vcenter: vCenters,
	}

	if len(p.FailureDomains) > 1 {
		cloudProviderConfig.Labels = &cloudconfig.Labels{
			Zone:   vspheretypes.TagCategoryZone,
			Region: vspheretypes.TagCategoryRegion,
		}
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
