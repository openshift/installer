package types

import (
	"fmt"

	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// InstallConfigVersion is the version supported by this package.
	// If you bump this, you must also update the list of convertable values in
	// pkg/conversion/installconfig.go
	InstallConfigVersion = "v1beta4"
)

var (
	// PlatformNames is a slice with all the visibly-supported
	// platform names in alphabetical order. This is the list of
	// platforms presented to the user in the interactive wizard.
	PlatformNames = []string{
		aws.Name,
		vsphere.Name,
	}
	// HiddenPlatformNames is a slice with all the
	// hidden-but-supported platform names. This list isn't presented
	// to the user in the interactive wizard.
	HiddenPlatformNames = []string{
		none.Name,
		openstack.Name,
	}
)

// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// SSHKey is the public ssh key to provide access to instances.
	// +optional
	SSHKey string `json:"sshKey,omitempty"`

	// BaseDomain is the base domain to which the cluster should belong.
	BaseDomain string `json:"baseDomain"`

	// Networking defines the pod network provider in the cluster.
	*Networking `json:"networking,omitempty"`

	// ControlPlane is the configuration for the machines that comprise the
	// control plane.
	// +optional
	ControlPlane *MachinePool `json:"controlPlane,omitempty"`

	// Compute is the list of compute MachinePools that need to be installed.
	// +optional
	Compute []MachinePool `json:"compute,omitempty"`

	// Platform is the configuration for the specific platform upon which to
	// perform the installation.
	Platform `json:"platform"`

	// PullSecret is the secret to use when pulling images.
	PullSecret string `json:"pullSecret"`
}

// ClusterDomain returns the DNS domain that all records for a cluster must belong to.
func (c *InstallConfig) ClusterDomain() string {
	return fmt.Sprintf("%s.%s", c.ObjectMeta.Name, c.BaseDomain)
}

// Platform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set.
type Platform struct {
	// AWS is the configuration used when installing on AWS.
	// +optional
	AWS *aws.Platform `json:"aws,omitempty"`

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
	case p.Libvirt != nil:
		return libvirt.Name
	case p.None != nil:
		return none.Name
	case p.OpenStack != nil:
		return openstack.Name
	case p.VSphere != nil:
		return vsphere.Name
	default:
		return ""
	}
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	// MachineCIDR is the IP address space from which to assign machine IPs.
	// +optional
	// Default is 10.0.0.0/16 for all platforms other than Libvirt.
	// For Libvirt, the default is 192.168.126.0/24.
	MachineCIDR *ipnet.IPNet `json:"machineCIDR,omitempty"`

	// NetworkType is the type of network to install.
	// +optional
	// Default is OpenShiftSDN.
	NetworkType string `json:"networkType,omitempty"`

	// ClusterNetwork is the IP address pool to use for pod IPs.
	// +optional
	// Default is 10.128.0.0/14 and a host prefix of /23
	ClusterNetwork []ClusterNetworkEntry `json:"clusterNetwork,omitempty"`

	// ServiceNetwork is the IP address pool to use for service IPs.
	// +optional
	// Default is 172.30.0.0/16
	// NOTE: currently only one entry is supported.
	ServiceNetwork []ipnet.IPNet `json:"serviceNetwork,omitempty"`

	// Deprected types, scheduled to be removed

	// Deprecated name for NetworkType
	// +optional
	DeprecatedType string `json:"type,omitempty"`

	// Depcreated name for ServiceNetwork
	// +optional
	DeprecatedServiceCIDR *ipnet.IPNet `json:"serviceCIDR,omitempty"`

	// Deprecated name for ClusterNetwork
	// +optional
	DeprecatedClusterNetworks []ClusterNetworkEntry `json:"clusterNetworks,omitempty"`
}

// ClusterNetworkEntry is a single IP address block for pod IP blocks. IP blocks
// are allocated with size 2^HostSubnetLength.
type ClusterNetworkEntry struct {
	// The IP block address pool
	CIDR ipnet.IPNet `json:"cidr"`

	// HostPrefix is the prefix size to allocate to each node from the CIDR.
	// For example, 24 would allocate 2^8=256 adresses to each node.
	HostPrefix int32 `json:"hostPrefix"`

	// The size of blocks to allocate from the larger pool.
	// This is the length in bits - so a 9 here will allocate a /23.
	DeprecatedHostSubnetLength int32 `json:"hostSubnetLength,omitempty"`
}
