package types

import (
	"net"

	"github.com/openshift/installer/pkg/ipnet"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// ClusterID is the ID of the cluster.
	ClusterID string `json:"clusterID"`

	// Admin is the configuration for the admin user.
	Admin Admin `json:"admin"`

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

// Admin is the configuration for the admin user.
type Admin struct {
	// Email is the email address of the admin user.
	Email string `json:"email"`
	// Password is the password of the admin user.
	Password string `json:"password"`
	// SSHKey to use for the access to compute instances.
	SSHKey string `json:"sshKey,omitempty"`
}

// Platform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set.
type Platform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *AWSPlatform `json:"aws,omitempty"`

	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *LibvirtPlatform `json:"libvirt,omitempty"`

	// OpenStack is the configuration used when installing on OpenStack.
	OpenStack *OpenStackPlatform `json:"openstack,omitempty"`
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	Type        NetworkType `json:"type"`
	ServiceCIDR ipnet.IPNet `json:"serviceCIDR"`
	PodCIDR     ipnet.IPNet `json:"podCIDR"`
}

// NetworkType defines the pod network provider in the cluster.
type NetworkType string

const (
	// NetworkTypeOpenshiftSDN is used to install with SDN.
	NetworkTypeOpenshiftSDN NetworkType = "openshift-sdn"
	// NetworkTypeOpenshiftOVN is used to install with OVN.
	NetworkTypeOpenshiftOVN NetworkType = "openshift-ovn"
)

// AWSPlatform stores all the global configuration that
// all machinesets use.
type AWSPlatform struct {
	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// UserTags specifies additional tags for AWS resources created for the cluster.
	UserTags map[string]string `json:"userTags,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	DefaultMachinePlatform *AWSMachinePoolPlatform `json:"defaultMachinePlatform,omitempty"`

	// VPCID specifies the vpc to associate with the cluster.
	// If empty, new vpc will be created.
	// +optional
	VPCID string `json:"vpcID"`

	// VPCCIDRBlock
	// +optional
	VPCCIDRBlock string `json:"vpcCIDRBlock"`
}

// OpenStackPlatform stores all the global configuration that
// all machinesets use.
type OpenStackPlatform struct {
	// Region specifies the OpenStack region where the cluster will be created.
	Region string `json:"region"`

	// VPCID specifies the vpc to associate with the cluster.
	// If empty, new vpc will be created.
	// +optional
	VPCID string `json:"vpcID"`

	// NetworkCIDRBlock
	// +optional
	NetworkCIDRBlock string `json:"NetworkCIDRBlock"`

	// BaseImage
	// Name of image to use from OpenStack cloud
	BaseImage string `json:"baseImage"`

	// Cloud
	// Name of OpenStack cloud to use from clouds.yaml
	Cloud string `json:"cloud"`

	// ExternalNetwork
	// The OpenStack external network to be used for installation.
	ExternalNetwork string `json:"externalNetwork"`
}

// LibvirtPlatform stores all the global configuration that
// all machinesets use.
type LibvirtPlatform struct {
	// URI is the identifier for the libvirtd connection.  It must be
	// reachable from both the host (where the installer is run) and the
	// cluster (where the cluster-API controller pod will be running).
	URI string `json:"URI"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	DefaultMachinePlatform *LibvirtMachinePoolPlatform `json:"defaultMachinePlatform,omitempty"`

	// Network
	Network LibvirtNetwork `json:"network"`

	// MasterIPs
	MasterIPs []net.IP `json:"masterIPs"`
}

// LibvirtNetwork is the configuration of the libvirt network.
type LibvirtNetwork struct {
	// Name is the name of the nework.
	Name string `json:"name"`
	// IfName is the name of the network interface.
	IfName string `json:"if"`
	// IPRange is the range of IPs to use.
	IPRange string `json:"ipRange"`
}
