package vsphere

import (
	"bytes"
	"fmt"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

func printIfNotEmpty(buf *bytes.Buffer, k, v string) {
	if v != "" {
		fmt.Fprintf(buf, "%s = %q\n", k, v)
	}
}

func getDatacenters(vcenter vspheretypes.VCenter) string {
	datacenters := make([]string, 0)
	for _, region := range vcenter.Regions {
		if region.Datacenter != "" {
			datacenters = append(datacenters, region.Datacenter)
		}
	}
	datacenterString := ""
	for idx, dataCenter := range datacenters {
		if idx > 0 {
			datacenterString = datacenterString + "," + dataCenter
		} else {
			datacenterString = dataCenter
		}
	}
	return datacenterString
}

// CloudProviderConfig generates the cloud provider config for the vSphere platform.
// folderPath is the absolute path to the VM folder that will be used for installation.
// p is the vSphere platform struct.
func CloudProviderConfig(folderPath string, p *vspheretypes.Platform) (string, error) {
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

	regions := make([]string, 0)
	zones := make([]string, 0)

	for _, vcenter := range p.VCenters {
		fmt.Fprintln(buf, "")
		fmt.Fprintf(buf, "[VirtualCenter %q]\n", vcenter.Server)
		printIfNotEmpty(buf, "datacenters", getDatacenters(vcenter))
		for _, region := range vcenter.Regions {
			regions = append(regions, region.Name)
			for _, zone := range region.Zones {
				zones = append(zones, zone.Name)
			}
		}
	}

	if p.VCenter != "" {
		fmt.Fprintln(buf, "")
		fmt.Fprintf(buf, "[VirtualCenter %q]\n", p.VCenter)
		printIfNotEmpty(buf, "datacenters", p.Datacenter)
	}

	if len(zones) > 0 {
		fmt.Fprintln(buf, "")
		fmt.Fprintln(buf, "[Labels]")
		regionsString := ""
		zonesString := ""
		for idx, region := range regions {
			if idx > 0 {
				regionsString = regionsString + "," + region
			} else {
				regionsString = region
			}
		}
		printIfNotEmpty(buf, "region", regionsString)

		for idx, zone := range zones {
			if idx > 0 {
				zonesString = zonesString + "," + zone
			} else {
				zonesString = zone
			}
		}
		printIfNotEmpty(buf, "zone", zonesString)
	}

	return buf.String(), nil
}
