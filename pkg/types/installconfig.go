package types

import (
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// InstallConfigVersion is the version supported by this package.
	InstallConfigVersion = "v1beta2"
)

var (
	// PlatformNames is a slice with all the visibly-supported
	// platform names in alphabetical order. This is the list of
	// platforms presented to the user in the interactive wizard.
	PlatformNames = []string{
		aws.Name,
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

	// Machines is the list of MachinePools that need to be installed.
	// +optional
	// Default on AWS and OpenStack is 3 masters and 3 workers.
	// Default on Libvirt is 1 master and 1 worker.
	Machines []MachinePool `json:"machines,omitempty"`

	// Platform is the configuration for the specific platform upon which to
	// perform the installation.
	Platform `json:"platform"`

	// PullSecret is the secret to use when pulling images.
	PullSecret string `json:"pullSecret"`
}

// MasterCount returns the number of replicas in the master machine pool,
// defaulting to one if no machine pool was found.
func (c *InstallConfig) MasterCount() int {
	for _, m := range c.Machines {
		if m.Name == "master" && m.Replicas != nil {
			return int(*m.Replicas)
		}
	}
	return 1
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
}

// Name returns a string representation of the platform (e.g. "aws" if
// AWS is non-nil).  It returns an empty string if no platform is
// configured.
func (p *Platform) Name() string {
	if p == nil {
		return ""
	}
	if p.AWS != nil {
		return aws.Name
	}
	if p.Libvirt != nil {
		return libvirt.Name
	}
	if p.None != nil {
		return none.Name
	}
	if p.OpenStack != nil {
		return openstack.Name
	}
	return ""
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	// MachineCIDR is the IP address space from which to assign machine IPs.
	// +optional
	// Default is 10.0.0.0/16 for all platforms other than Libvirt.
	// For Libvirt, the default is 192.168.126.0/24.
	MachineCIDR *ipnet.IPNet `json:"machineCIDR,omitempty"`

	// Type is the network type to install
	// +optional
	// Default is OpenShiftSDN.
	Type string `json:"type,omitempty"`

	// ServiceCIDR is the IP address space from which to assign service IPs.
	// +optional
	// Default is 172.30.0.0/16.
	ServiceCIDR *ipnet.IPNet `json:"serviceCIDR,omitempty"`

	// ClusterNetworks is the IP address space from which to assign pod IPs.
	// +optional
	// Default is a single cluster network with a CIDR of 10.128.0.0/14
	// and a host subnet length of 9.
	ClusterNetworks []ClusterNetworkEntry `json:"clusterNetworks,omitempty"`
}

// ClusterNetworkEntry is a single IP address block for pod IP blocks. IP blocks
// are allocated with size 2^HostSubnetLength.
type ClusterNetworkEntry struct {
	// The IP block address pool
	CIDR ipnet.IPNet `json:"cidr"`

	// The size of blocks to allocate from the larger pool.
	// This is the length in bits - so a 9 here will allocate a /23.
	HostSubnetLength int32 `json:"hostSubnetLength"`
}
