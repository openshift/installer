package extract

import (
	"fmt"

	"github.com/openshift/oc/pkg/cli/image/extract"
)

const (
	// target architectures
	amd64 = "amd64"

	// target operating systems
	linux = "linux"

	// terraform providers #TODO: this list is duplicated between here and pkg/terraform/providers.providers.go: dedupe
	alicloud       = "alicloud"
	aws            = "aws"
	azurerm        = "azurerm"
	azurestack     = "azurestack"
	google         = "google"
	ibm            = "ibm"
	ignition       = "ignition"
	ironic         = "ironic"
	libvirt        = "libvirt"
	local          = "local"
	nutanix        = "nutanix"
	openstack      = "openstack"
	ovirt          = "ovirt"
	time           = "time"
	vsphere        = "vsphere"
	vsphereprivate = "vsphereprivate"
)

var (
	availableTargets = generate()
	architectures    = []string{amd64}
	oss              = []string{linux}
	tfProviders      = []string{alicloud, aws, azurerm, azurestack, google, ibm, ignition, ironic, libvirt, local, nutanix, openstack, ovirt, time, vsphere, vsphereprivate}
)

func generate() []extractTarget {
	targets := []extractTarget{
		{
			Command: "terraform",
			OS:      "linux",
			Arch:    "amd64",
			Mapping: extract.Mapping{Image: "installer-artifacts", From: "usr/share/openshift/linux_amd64/terraform/terraform"},
		},
	}

	for _, tfProvider := range tfProviders {
		for _, arch := range architectures {
			for _, os := range oss {
				target := extractTarget{
					Command: tfProvCommand(tfProvider, os, arch),
					OS:      os,
					Arch:    arch,
					Mapping: extract.Mapping{Image: "installer-artifacts", From: srcImgPath(tfProvider, os, arch)},
				}
				targets = append(targets, target)
			}
		}
	}

	return targets
}

func tfProvCommand(provName, os, arch string) string {
	return fmt.Sprintf("terraform-provider-%s_1.0.0_%s_%s.zip", provName, os, arch)
}

func srcImgPath(provName, os, arch string) string {
	return fmt.Sprintf("usr/share/openshift/%s_%s/terraform/terraform-provider-%s.zip", os, arch, provName)
}
