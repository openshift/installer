package vsphere

import (
	"bytes"
	"fmt"
	"strings"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	yaml "gopkg.in/yaml.v2"
	cloudconfig "k8s.io/cloud-provider-vsphere/pkg/common/config"
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

func appendTagCategory(tagCategory string, tagCategories []string) []string {
	tagDefined := false
	for _, regionTagCategory := range tagCategories {
		if regionTagCategory == tagCategory {
			tagDefined = true
			break
		}
	}
	if tagDefined == false {
		return append(tagCategories, tagCategory)
	}
	return tagCategories
}

// MultiZoneYamlCloudProviderConfig generates the yaml out of tree cloud provider config for the vSphere platform.
func MultiZoneYamlCloudProviderConfig(p *vspheretypes.Platform) (string, error) {
	vCenters := make(map[string]*cloudconfig.VirtualCenterConfigYAML)

	for _, vCenter := range p.VCenters {
		vCenterPort := uint(443)
		if vCenter.Port != 0 {
			vCenterPort = vCenter.Port
		}
		vCenterConfig := cloudconfig.VirtualCenterConfigYAML{
			VCenterIP:   vCenter.Server,
			VCenterPort: vCenterPort,
			Datacenters: vCenter.Datacenters,
		}
		vCenters[vCenter.Server] = &vCenterConfig
	}

	cloudProviderConfig := cloudconfig.CommonConfigYAML{
		Global: cloudconfig.GlobalYAML{
			SecretName:      "vsphere-creds",
			SecretNamespace: "kube-system",
		},
		Vcenter: vCenters,
		Labels: cloudconfig.LabelsYAML{
			Zone:   zoneTagCategory,
			Region: regionTagCategory,
		},
	}

	cloudProviderConfigYaml, err := yaml.Marshal(cloudProviderConfig)
	if err != nil {
		return "", err
	}
	return string(cloudProviderConfigYaml), nil
}

// MultiZoneIniCloudProviderConfig generates the multi-zone ini cloud provider config
// for the vSphere platform. folderPath is the absolute path to the VM folder that will be
// used for installation. p is the vSphere platform struct.
func MultiZoneIniCloudProviderConfig(folderPath string, p *vspheretypes.Platform) (string, error) {
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
		var datacenters []string
		for _, datacenter := range vcenter.Datacenters {
			datacenters = append(datacenters, datacenter)
		}
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
				if exists == false {
					datacenters = append(datacenters, failureDomainDatacenter)
				}
			}
		}
		printIfNotEmpty(buf, "datacenters", strings.Join(datacenters, ","))
	}
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "[Workspace]")
	printIfNotEmpty(buf, "server", p.VCenter)
	printIfNotEmpty(buf, "datacenter", p.Datacenter)
	printIfNotEmpty(buf, "default-datastore", p.DefaultDatastore)
	printIfNotEmpty(buf, "folder", folderPath)
	printIfNotEmpty(buf, "resourcepool-path", p.ResourcePool)
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "[Labels]")
	printIfNotEmpty(buf, "region", regionTagCategory)
	printIfNotEmpty(buf, "zone", zoneTagCategory)

	return buf.String(), nil
}

// InTreeCloudProviderConfig generates the in-tree cloud provider config for the vSphere platform.
// folderPath is the absolute path to the VM folder that will be used for installation.
// p is the vSphere platform struct.
func InTreeCloudProviderConfig(folderPath string, p *vspheretypes.Platform) (string, error) {
	buf := new(bytes.Buffer)

	fmt.Fprintln(buf, "[Global]")
	printIfNotEmpty(buf, "secret-name", "vsphere-creds")
	printIfNotEmpty(buf, "secret-namespace", "kube-system")
	printIfNotEmpty(buf, "insecure-flag", "1")
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "[Workspace]")
	printIfNotEmpty(buf, "server", p.VCenter)
	printIfNotEmpty(buf, "datacenter", p.Datacenter)
	printIfNotEmpty(buf, "default-datastore", p.DefaultDatastore)
	printIfNotEmpty(buf, "folder", folderPath)
	printIfNotEmpty(buf, "resourcepool-path", p.ResourcePool)
	fmt.Fprintln(buf, "")

	fmt.Fprintf(buf, "[VirtualCenter %q]\n", p.VCenter)
	printIfNotEmpty(buf, "datacenters", p.Datacenter)

	return buf.String(), nil
}
