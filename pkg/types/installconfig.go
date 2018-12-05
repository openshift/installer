package types

import (
	netopv1 "github.com/openshift/cluster-network-operator/pkg/apis/networkoperator/v1"
	"github.com/openshift/installer/pkg/ipnet"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/openstack"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	// PlatformNames is a slice with all the supported platform names in
	// alphabetical order.
	PlatformNames = []string{
		aws.Name,
		libvirt.Name,
		openstack.Name,
	}
)

// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// ClusterID is the ID of the cluster.
	ClusterID string `json:"clusterID"`

	// SSHKey is the public ssh key to provide access to instances.
	SSHKey string `json:"sshKey"`

	// BaseDomain is the base domain to which the cluster should belong.
	BaseDomain string `json:"baseDomain"`

	// Networking defines the pod network provider in the cluster.
	Networking `json:"networking"`

	// Machines is the list of MachinePools that need to be installed.
	Machines []MachinePool `json:"machines"`

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
	AWS *aws.Platform `json:"aws,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *libvirt.Platform `json:"libvirt,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
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
	if p.OpenStack != nil {
		return openstack.Name
	}
	return ""
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	// Type is the network type to install
	Type netopv1.NetworkType `json:"type"`

	// ServiceCIDR is the ip block from which to assign service IPs
	ServiceCIDR ipnet.IPNet `json:"serviceCIDR"`

	// ClusterNetworks is the IP address space from which to assign pod IPs.
	ClusterNetworks []netopv1.ClusterNetwork `json:"clusterNetworks,omitempty"`

	// PodCIDR is deprecated (and badly named; it should have always
	// been called ClusterCIDR. If no ClusterNetworks are specified,
	// we will fall back to the PodCIDR
	// TODO(cdc) remove this.
	PodCIDR *ipnet.IPNet `json:"podCIDR,omitempty"`
}
