package providers

import (
	"fmt"
	"strings"
)

var (
	// AliCloud is the provider for creating resources in the alibaba cloud.
	AliCloud = remoteProvider("hashicorp/alicloud")
	// AWS is the provider for creating resources in AWS.
	AWS = remoteProvider("hashicorp/aws")
	// AzurePrivateDNS is an internal provider for creating private DNS zones in Azure.
	AzurePrivateDNS = localProvider("azureprivatedns")
	// AzureRM is the provider for creating resources in the Azure clouds.
	AzureRM = remoteProvider("hashicorp/azurerm")
	// AzureStack is the provider for creating resources in Azure Stack.
	AzureStack = localProvider("azurestack")
	// Google is the provider for creating resources in GCP.
	Google = remoteProvider("hashicorp/google")
	// IBM is the provider for creating resources in IBM Cloud.
	IBM = localProvider("ibm")
	// Ignition is the provider for creating ignition config files.
	Ignition = localProvider("ignition")
	// Ironic is the provider for provisioning baremetal machines managed by Ironic.
	Ironic = localProvider("ironic")
	// Libvirt is the provider for provisioning VMs on a libvirt host.
	Libvirt = remoteProvider("dmacvicar/libvirt")
	// Local is the provider for creating local files.
	Local = remoteProvider("hashicorp/local")
	// OpenStack is the provider for creating resources in OpenStack.
	OpenStack = remoteProvider("terraform-provider-openstack/openstack")
	// OVirt is the provider for creating resources in oVirt.
	OVirt = localProvider("ovirt")
	// Random is the provider for generating randomness.
	Random = remoteProvider("hashicorp/random")
	// VSphere is the provider for creating resource in vSphere.
	VSphere = localProvider("vsphere")
	// VSpherePrivate is an internal provider augmenting the VSphere provider by adding functionality.
	VSpherePrivate = localProvider("vsphereprivate")
)

// Provider is a terraform provider.
type Provider struct {
	// Name of the provider.
	Name string
	// Source of the provider.
	Source string
	// Version of the provider.
	// This can be omitted for remote providers of which the installer is only embedding a single version.
	Version string
	// IsLocal is true if the provider is built locally as opposed to being downloaded from a remote registry.
	IsLocal bool
}

// remoteProvider configures a provider downloaded from a remote registry.
func remoteProvider(source string) Provider {
	sourceParts := strings.Split(source, "/")
	switch len(sourceParts) {
	case 1:
		source = "hashicorp/" + source
		fallthrough
	case 2:
		source = "registry.terraform.io/" + source
	}
	name := sourceParts[len(sourceParts)-1]
	return Provider{
		Name:   name,
		Source: source,
	}
}

// localProvider configures a provider built locally.
func localProvider(name string) Provider {
	return Provider{
		Name:    name,
		Source:  fmt.Sprintf("openshift/local/%s", name),
		Version: "1.0.0",
		IsLocal: true,
	}
}
