package types

import (
	"net"

	"github.com/pborman/uuid"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InstallConfig struct {
	// +optional
	metav1.TypeMeta `json:",inline"`

	metav1.ObjectMeta `json:"metadata"`

	// ClusterID is the ID of the cluster.
	ClusterID uuid.UUID `json:"clusterID"`

	// Networking defines the pod network provider in the cluster.
	Networking `json:"networking"`

	// Machines is the list of MachinePools that need to be installed.
	Machines []MachinePools `json:"machines"`

	// only one of the platform configuration should be set
	Platform `json:"platform"`
}

type Platform struct {
	AWS     *AWSPlatform     `json:"aws,omitempty"`
	Libvirt *LibvirtPlatform `json:"libvirt,omitempty"`
}

type Networking struct {
	Type        NetworkType `json:"type"`
	ServiceCIDR net.IPNet   `json:"serviceCIDR"`
	PodCIDR     net.IPNet   `json:"podCIDR"`
}

// NetworkType defines the pod network provider in the cluster.
type NetworkType string

const (
	// NetworkTypeOpenshiftSDN
	NetworkTypeOpenshiftSDN NetworkType = "openshift-sdn"
	// NetworkTypeOpenshiftOVN
	NetworkTypeOpenshiftOVN NetworkType = "openshift-ovn"
)

// AWSPlatform stores all the global configuration that
// all machinesets use.
type AWSPlatform struct {
	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// KeyPairName is the name of the AWS key pair to use for SSH access to EC2
	// instances in this cluster.
	KeyPairName string `json:"keyPairName"`

	// VPCID specifies the vpc to associate with the cluster.
	// If empty, new vpc will be created.
	// +optional
	VPCID string `json:"vpcID"`

	// VPCCIDRBlock
	// +optional
	VPCCIDRBlock string `json:"vpcCIDRBlock"`
}

// LibvirtPlatform stores all the global configuration that
// all machinesets use.
type LibvirtPlatform struct {
	// URI
	URI string `json:"URI"`

	// SSHKey
	SSHKey string `json:"sshKey"`

	// Network
	Network LibvirtNetwork `json:"network"`

	// MasterIPs
	MasterIPs []net.IP `json:"masterIPs"`
}

type LibvirtNetwork struct {
	Name      string `json:"name"`
	IfName    string `json:"if"`
	DNSServer string `json:"resolver"`
	IPRange   string `json:"ipRange"`
}
