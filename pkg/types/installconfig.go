package types

import (
	"fmt"
	"strings"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/kubevirt"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/types/vsphere"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// InstallConfigVersion is the version supported by this package.
	// If you bump this, you must also update the list of convertable values in
	// pkg/types/conversion/installconfig.go
	InstallConfigVersion = "v1"
)

var (
	// PlatformNames is a slice with all the visibly-supported
	// platform names in alphabetical order. This is the list of
	// platforms presented to the user in the interactive wizard.
	PlatformNames = []string{
		aws.Name,
		azure.Name,
		gcp.Name,
		openstack.Name,
		ovirt.Name,
		vsphere.Name,
	}
	// HiddenPlatformNames is a slice with all the
	// hidden-but-supported platform names. This list isn't presented
	// to the user in the interactive wizard.
	HiddenPlatformNames = []string{
		baremetal.Name,
		kubevirt.Name,
		none.Name,
	}
)

// PublishingStrategy is a strategy for how various endpoints for the cluster are exposed.
// +kubebuilder:validation:Enum="";External;Internal
type PublishingStrategy string

const (
	// ExternalPublishingStrategy exposes endpoints for the cluster to the Internet.
	ExternalPublishingStrategy PublishingStrategy = "External"
	// InternalPublishingStrategy exposes the endpoints for the cluster to the private network only.
	InternalPublishingStrategy PublishingStrategy = "Internal"
)

//go:generate go run ../../vendor/sigs.k8s.io/controller-tools/cmd/controller-gen crd:crdVersions=v1 paths=. output:dir=../../data/data/

// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// AdditionalTrustBundle is a PEM-encoded X.509 certificate bundle
	// that will be added to the nodes' trusted certificate store.
	//
	// +optional
	AdditionalTrustBundle string `json:"additionalTrustBundle,omitempty"`

	// SSHKey is the public Secure Shell (SSH) key to provide access to instances.
	// +optional
	SSHKey string `json:"sshKey,omitempty"`

	// BaseDomain is the base domain to which the cluster should belong.
	BaseDomain string `json:"baseDomain"`

	// Networking is the configuration for the pod network provider in
	// the cluster.
	*Networking `json:"networking,omitempty"`

	// ControlPlane is the configuration for the machines that comprise the
	// control plane.
	// +optional
	ControlPlane *MachinePool `json:"controlPlane,omitempty"`

	// Compute is the configuration for the machines that comprise the
	// compute nodes.
	// +optional
	Compute []MachinePool `json:"compute,omitempty"`

	// Platform is the configuration for the specific platform upon which to
	// perform the installation.
	Platform `json:"platform"`

	// PullSecret is the secret to use when pulling images.
	PullSecret string `json:"pullSecret"`

	// Proxy defines the proxy settings for the cluster.
	// If unset, the cluster will not be configured to use a proxy.
	// +optional
	Proxy *Proxy `json:"proxy,omitempty"`

	// ImageContentSources lists sources/repositories for the release-image content.
	// +optional
	ImageContentSources []ImageContentSource `json:"imageContentSources,omitempty"`

	// Publish controls how the user facing endpoints of the cluster like the Kubernetes API, OpenShift routes etc. are exposed.
	// When no strategy is specified, the strategy is "External".
	//
	// +kubebuilder:default=External
	// +optional
	Publish PublishingStrategy `json:"publish,omitempty"`

	// FIPS configures https://www.nist.gov/itl/fips-general-information
	//
	// +kubebuilder:default=false
	// +optional
	FIPS bool `json:"fips,omitempty"`

	// CredentialsMode is used to explicitly set the mode with which CredentialRequests are satisfied.
	//
	// If this field is set, then the installer will not attempt to query the cloud permissions before attempting
	// installation. If the field is not set or empty, then the installer will perform its normal verification that the
	// credentials provided are sufficient to perform an installation.
	//
	// There are three possible values for this field, but the valid values are dependent upon the platform being used.
	// "Mint": create new credentials with a subset of the overall permissions for each CredentialsRequest
	// "Passthrough": copy the credentials with all of the overall permissions for each CredentialsRequest
	// "Manual": CredentialsRequests must be handled manually by the user
	//
	// For each of the following platforms, the field can set to the specified values. For all other platforms, the
	// field must not be set.
	// AWS: "Mint", "Passthrough", "Manual"
	// Azure: "Mint", "Passthrough", "Manual"
	// GCP: "Mint", "Passthrough", "Manual"
	// +optional
	CredentialsMode CredentialsMode `json:"credentialsMode,omitempty"`

	// BootstrapInPlace is the configuration for installing a single node
	// with bootstrap in place installation.
	BootstrapInPlace *BootstrapInPlace `json:"bootstrapInPlace,omitempty"`
}

// ClusterDomain returns the DNS domain that all records for a cluster must belong to.
func (c *InstallConfig) ClusterDomain() string {
	return fmt.Sprintf("%s.%s", c.ObjectMeta.Name, strings.TrimSuffix(c.BaseDomain, "."))
}

// Platform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set.
type Platform struct {
	// AWS is the configuration used when installing on AWS.
	// +optional
	AWS *aws.Platform `json:"aws,omitempty"`

	// Azure is the configuration used when installing on Azure.
	// +optional
	Azure *azure.Platform `json:"azure,omitempty"`

	// BareMetal is the configuration used when installing on bare metal.
	// +optional
	BareMetal *baremetal.Platform `json:"baremetal,omitempty"`

	// GCP is the configuration used when installing on Google Cloud Platform.
	// +optional
	GCP *gcp.Platform `json:"gcp,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	// +optional
	Libvirt *libvirt.Platform `json:"libvirt,omitempty"`

	// None is the empty configuration used when installing on an unsupported
	// platform.
	None *none.Platform `json:"none,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	// +optional
	OpenStack *openstack.Platform `json:"openstack,omitempty"`

	// VSphere is the configuration used when installing on vSphere.
	// +optional
	VSphere *vsphere.Platform `json:"vsphere,omitempty"`

	// Ovirt is the configuration used when installing on oVirt.
	// +optional
	Ovirt *ovirt.Platform `json:"ovirt,omitempty"`

	// Kubevirt is the configuration used when installing on kubevirt.
	// +optional
	Kubevirt *kubevirt.Platform `json:"kubevirt,omitempty"`
}

// Name returns a string representation of the platform (e.g. "aws" if
// AWS is non-nil).  It returns an empty string if no platform is
// configured.
func (p *Platform) Name() string {
	switch {
	case p == nil:
		return ""
	case p.AWS != nil:
		return aws.Name
	case p.Azure != nil:
		return azure.Name
	case p.BareMetal != nil:
		return baremetal.Name
	case p.GCP != nil:
		return gcp.Name
	case p.Libvirt != nil:
		return libvirt.Name
	case p.None != nil:
		return none.Name
	case p.OpenStack != nil:
		return openstack.Name
	case p.VSphere != nil:
		return vsphere.Name
	case p.Ovirt != nil:
		return ovirt.Name
	case p.Kubevirt != nil:
		return kubevirt.Name
	default:
		return ""
	}
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	// NetworkType is the type of network to install. The default is OpenShiftSDN
	//
	// +kubebuilder:default=OpenShiftSDN
	// +optional
	NetworkType string `json:"networkType,omitempty"`

	// MachineNetwork is the list of IP address pools for machines.
	// This field replaces MachineCIDR, and if set MachineCIDR must
	// be empty or match the first entry in the list.
	// Default is 10.0.0.0/16 for all platforms other than libvirt.
	// For libvirt, the default is 192.168.126.0/24.
	//
	// +optional
	MachineNetwork []MachineNetworkEntry `json:"machineNetwork,omitempty"`

	// ClusterNetwork is the list of IP address pools for pods.
	// Default is 10.128.0.0/14 and a host prefix of /23.
	//
	// +optional
	ClusterNetwork []ClusterNetworkEntry `json:"clusterNetwork,omitempty"`

	// ServiceNetwork is the list of IP address pools for services.
	// Default is 172.30.0.0/16.
	// NOTE: currently only one entry is supported.
	//
	// +kubebuilder:validation:MaxItems=1
	// +optional
	ServiceNetwork []ipnet.IPNet `json:"serviceNetwork,omitempty"`

	// Deprecated types, scheduled to be removed

	// Deprecated name for MachineCIDRs. If set, MachineCIDRs must
	// be empty or the first index must match.
	// +optional
	DeprecatedMachineCIDR *ipnet.IPNet `json:"machineCIDR,omitempty"`

	// Deprecated name for NetworkType
	// +optional
	DeprecatedType string `json:"type,omitempty"`

	// Deprecated name for ServiceNetwork
	// +optional
	DeprecatedServiceCIDR *ipnet.IPNet `json:"serviceCIDR,omitempty"`

	// Deprecated name for ClusterNetwork
	// +optional
	DeprecatedClusterNetworks []ClusterNetworkEntry `json:"clusterNetworks,omitempty"`
}

// MachineNetworkEntry is a single IP address block for node IP blocks.
type MachineNetworkEntry struct {
	// CIDR is the IP block address pool for machines within the cluster.
	CIDR ipnet.IPNet `json:"cidr"`
}

// ClusterNetworkEntry is a single IP address block for pod IP blocks. IP blocks
// are allocated with size 2^HostSubnetLength.
type ClusterNetworkEntry struct {
	// CIDR is the IP block address pool.
	CIDR ipnet.IPNet `json:"cidr"`

	// HostPrefix is the prefix size to allocate to each node from the CIDR.
	// For example, 24 would allocate 2^8=256 adresses to each node. If this
	// field is not used by the plugin, it can be left unset.
	// +optional
	HostPrefix int32 `json:"hostPrefix,omitempty"`

	// The size of blocks to allocate from the larger pool.
	// This is the length in bits - so a 9 here will allocate a /23.
	// +optional
	DeprecatedHostSubnetLength int32 `json:"hostSubnetLength,omitempty"`
}

// Proxy defines the proxy settings for the cluster.
// At least one of HTTPProxy or HTTPSProxy is required.
type Proxy struct {
	// HTTPProxy is the URL of the proxy for HTTP requests.
	// +optional
	HTTPProxy string `json:"httpProxy,omitempty"`

	// HTTPSProxy is the URL of the proxy for HTTPS requests.
	// +optional
	HTTPSProxy string `json:"httpsProxy,omitempty"`

	// NoProxy is a comma-separated list of domains and CIDRs for which the proxy should not be used.
	// +optional
	NoProxy string `json:"noProxy,omitempty"`
}

// ImageContentSource defines a list of sources/repositories that can be used to pull content.
type ImageContentSource struct {
	// Source is the repository that users refer to, e.g. in image pull specifications.
	Source string `json:"source"`

	// Mirrors is one or more repositories that may also contain the same images.
	// +optional
	Mirrors []string `json:"mirrors,omitempty"`
}

// CredentialsMode is the mode by which CredentialsRequests will be satisfied.
// +kubebuilder:validation:Enum="";Mint;Passthrough;Manual
type CredentialsMode string

const (
	// ManualCredentialsMode indicates that cloud-credential-operator should not process any CredentialsRequests.
	ManualCredentialsMode CredentialsMode = "Manual"

	// MintCredentialsMode indicates that cloud-credential-operator should be creating users for each
	// CredentialsRequest.
	MintCredentialsMode CredentialsMode = "Mint"

	// PassthroughCredentialsMode indicates that cloud-credential-operator should just copy over the cluster's
	// cloud credentials for each CredentialsRequest.
	PassthroughCredentialsMode CredentialsMode = "Passthrough"
)

// BootstrapInPlace defines the configuration for bootstrap-in-place installation
type BootstrapInPlace struct {
	// InstallationDisk is the target disk drive for coreos-installer
	InstallationDisk string `json:"installationDisk"`
}
