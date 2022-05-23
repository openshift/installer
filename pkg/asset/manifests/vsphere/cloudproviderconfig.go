package vsphere

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"strings"

	"github.com/go-yaml/yaml"
	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	cloudconfig "k8s.io/cloud-provider-vsphere/pkg/common/config"
)

func printIfNotEmpty(buf *bytes.Buffer, k, v string) {
	if v != "" {
		fmt.Fprintf(buf, "%s = %q\n", k, v)
	}
}

func getFailureDomain(deploymentZone vspheretypes.DeploymentZoneSpec, p *vspheretypes.Platform) (*vspheretypes.FailureDomainSpec, error) {

	for _, failureDomain := range p.FailureDomains {
		if failureDomain.Name == deploymentZone.FailureDomain {
			return &failureDomain, nil
		}
	}
	return nil, fmt.Errorf("failure domain %s not found", deploymentZone.FailureDomain)
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

func getTagCategoriesForVcenter(p *vspheretypes.Platform) (string, string, error) {
	regionTagCategories := make([]string, 0)
	zoneTagCategories := make([]string, 0)
	for _, vCenter := range p.VCenters {
		for _, deploymentZone := range p.DeploymentZones {
			if deploymentZone.Server != vCenter.Server {
				continue
			}
			failureDomain, err := getFailureDomain(deploymentZone, p)
			if err != nil {
				return "", "", err
			}
			regionTagCategories = appendTagCategory(failureDomain.Region.TagCategory, regionTagCategories)
			zoneTagCategories = appendTagCategory(failureDomain.Zone.TagCategory, zoneTagCategories)
		}
	}
	return strings.Join(regionTagCategories, ","), strings.Join(zoneTagCategories, ","), nil
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

	regionTagCategory, zoneTagCategory, err := getTagCategoriesForVcenter(p)
	if err != nil {
		return "", err
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
		printIfNotEmpty(buf, "datacenters", strings.Join(vcenter.Datacenters, ","))
	}
	fmt.Fprintln(buf, "")

	fmt.Fprintln(buf, "[Workspace]")
	printIfNotEmpty(buf, "server", p.VCenter)
	printIfNotEmpty(buf, "datacenter", p.Datacenter)
	printIfNotEmpty(buf, "default-datastore", p.DefaultDatastore)
	printIfNotEmpty(buf, "folder", folderPath)
	printIfNotEmpty(buf, "resourcepool-path", p.ResourcePool)
	fmt.Fprintln(buf, "")

	regionTagCategory, zoneTagCategory, err := getTagCategoriesForVcenter(p)
	if err != nil {
		return "", errors.Wrap(err, "error adding zones to the cloud-config")
	}

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
