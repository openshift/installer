package types

import (
	"net"

	"github.com/openshift/installer/pkg/ipnet"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstallConfig is the configuration for an OpenShift install.
type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline" yaml:",inline"`

	metav1.ObjectMeta `json:"metadata" yaml:"metadata"`

	// ClusterID is the ID of the cluster.
	ClusterID string `json:"clusterID" yaml:"clusterID"`

	// Admin is the configuration for the admin user.
	Admin Admin `json:"admin" yaml:"admin"`

	// BaseDomain is the base domain to which the cluster should belong.
	BaseDomain string `json:"baseDomain" yaml:"baseDomain"`

	// Networking defines the pod network provider in the cluster.
	Networking `json:"networking" yaml:"networking"`

	// Machines is the list of MachinePools that need to be installed.
	Machines []MachinePool `json:"machines" yaml:"machines"`

	// Platform is the configuration for the specific platform upon which to
	// perform the installation.
	Platform `json:"platform" yaml:"platform"`

	// PullSecret is the secret to use when pulling images.
	PullSecret string `json:"pullSecret,omitempty" yaml:"pullSecret,omitempty"`
}

// Admin is the configuration for the admin user.
type Admin struct {
	// Email is the email address of the admin user.
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
	// Password is the password of the admin user.
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// SSHKey to use for the access to compute instances.
	SSHKey string `json:"sshKey,omitempty" yaml:"sshKey,omitempty"`
}

// Platform is the configuration for the specific platform upon which to perform
// the installation. Only one of the platform configuration should be set.
type Platform struct {
	// AWS is the configuration used when installing on AWS.
	AWS *AWSPlatform `json:"aws,omitempty" yaml:"aws,omitempty"`
	// Libvirt is the configuration used when installing on libvirt.
	Libvirt *LibvirtPlatform `json:"libvirt,omitempty" yaml:"libvirt,omitempty"`
}

// Networking defines the pod network provider in the cluster.
type Networking struct {
	Type        NetworkType `json:"type,omitempty" yaml:"type,omitempty"`
	ServiceCIDR ipnet.IPNet `json:"serviceCIDR" yaml:"serviceCIDR"`
	PodCIDR     ipnet.IPNet `json:"podCIDR" yaml:"podCIDR"`
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
	Region string `json:"region,omitempty" yaml:"region,omitempty"`

	// UserTags specifies additional tags for AWS resources created for the cluster.
	UserTags map[string]string `json:"userTags,omitempty" yaml:"userTags,omitempty"`

	// VPCID specifies the vpc to associate with the cluster.
	// If empty, new vpc will be created.
	// +optional
	VPCID string `json:"vpcID,omitempty" yaml:"vpcID,omitempty"`

	// VPCCIDRBlock
	// +optional
	VPCCIDRBlock string `json:"vpcCIDRBlock,omitempty" yaml:"vpcCIDRBlock,omitempty"`
}

// LibvirtPlatform stores all the global configuration that
// all machinesets use.
type LibvirtPlatform struct {
	// URI
	URI string `json:"URI,omitempty" yaml:"URI,omitempty"`

	// Network
	Network LibvirtNetwork `json:"network" yaml:"network"`

	// MasterIPs
	MasterIPs []net.IP `json:"masterIPs,omitempty" yaml:"masterIPs,omitempty"`
}

// LibvirtNetwork is the configuration of the libvirt network.
type LibvirtNetwork struct {
	// Name is the name of the nework.
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// IfName is the name of the network interface.
	IfName string `json:"if,omitempty" yaml:"if,omitempty"`
	// IPRange is the range of IPs to use.
	IPRange string `json:"ipRange,omitempty" yaml:"ipRange,omitempty"`
}
