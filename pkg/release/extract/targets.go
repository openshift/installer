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
			OS:      "linux", //TODO check how OS/ARCH are used and what needs to be done.
			Arch:    "amd64",
			Mapping: extract.Mapping{Image: "installer-artifacts", From: "usr/share/terraform/linux_amd64/terraform"},
		},
		// CAPI etcd & kas targets
		// TODO: expand to all arches, os, & static/dynamic linking
		{
			Command: "etcd",
			OS:      "linux",
			Arch:    "amd64",
			Mapping: extract.Mapping{
				Image: "etcd",
				From:  "usr/bin/etcd",
			},
		},
		{
			Command: "kube-apiserver",
			OS:      "linux",
			Arch:    "amd64",
			Mapping: extract.Mapping{
				Image: "hyperkube",
				From:  "usr/bin/kube-apiserver",
			},
		},
	}

	// TODO: also need to account for static vs dynamically linked binaries
	// I think this would mean changing the source image.
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

	targets = append(targets, []extractTarget{}...)

	return targets
}

func tfProvCommand(provName, os, arch string) string {
	return fmt.Sprintf("terraform-provider-%s_1.0.0_%s_%s.zip", provName, os, arch)
}

func srcImgPath(provName, os, arch string) string {
	return fmt.Sprintf("usr/share/terraform/%s_%s/terraform-provider-%s.zip", os, arch, provName)
}
