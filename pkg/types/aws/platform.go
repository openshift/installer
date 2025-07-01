package aws

import (
	"os"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/dns"
)

const (
	// VolumeTypeGp2 is the type of EBS volume for General Purpose SSD gp2.
	VolumeTypeGp2 = "gp2"
	// VolumeTypeGp3 is the type of EBS volume for General Purpose SSD gp3.
	VolumeTypeGp3 = "gp3"
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// The field is deprecated. AMIID is the AMI that should be used to boot
	// machines for the cluster. If set, the AMI should belong to the same
	// region as the cluster.
	//
	// +optional
	AMIID string `json:"amiID,omitempty"`

	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	//
	// Deprecated: use platform.aws.vpc.subnets
	//
	// +optional
	DeprecatedSubnets []string `json:"subnets,omitempty"`

	// VPC specifies the VPC configuration for the cluster.
	//
	// +optional
	VPC VPC `json:"vpc,omitempty"`

	// HostedZone is the ID of an existing hosted zone into which to add DNS
	// records for the cluster's internal API. An existing hosted zone can
	// only be used when also using existing subnets. The hosted zone must be
	// associated with the VPC containing the subnets.
	// Leave the hosted zone unset to have the installer create the hosted zone
	// on your behalf.
	// +optional
	HostedZone string `json:"hostedZone,omitempty"`

	// HostedZoneRole is the ARN of an IAM role to be assumed when performing
	// operations on the provided HostedZone. HostedZoneRole can be used
	// in a shared VPC scenario when the private hosted zone belongs to a
	// different account than the rest of the cluster resources.
	// If HostedZoneRole is set, HostedZone must also be set.
	//
	// +optional
	HostedZoneRole string `json:"hostedZoneRole,omitempty"`

	// UserTags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// ServiceEndpoints list contains custom endpoints which will override default
	// service endpoint of AWS Services.
	// There must be only one ServiceEndpoint for a service.
	// +optional
	ServiceEndpoints []ServiceEndpoint `json:"serviceEndpoints,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on AWS for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// PropagateUserTags is a flag that directs in-cluster operators
	// to include the specified user tags in the tags of the
	// AWS resources that the operators create.
	// +optional
	PropagateUserTag bool `json:"propagateUserTags,omitempty"`

	// LBType is an optional field to specify a load balancer type.
	// When this field is specified, all ingresscontrollers (including the
	// default ingresscontroller) will be created using the specified load-balancer
	// type by default.
	//
	// Following are the accepted values:
	//
	// * "Classic": A Classic Load Balancer that makes routing decisions at
	// either the transport layer (TCP/SSL) or the application layer
	// (HTTP/HTTPS). See the following for additional details:
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#clb
	//
	// * "NLB": A Network Load Balancer that makes routing decisions at the
	// transport layer (TCP/SSL). See the following for additional details:
	// https://docs.aws.amazon.com/AmazonECS/latest/developerguide/load-balancer-types.html#nlb
	//
	// If this field is not set explicitly, it defaults to "Classic".  This
	// default is subject to change over time.
	//
	// +optional
	LBType configv1.AWSLBType `json:"lbType,omitempty"`

	// PreserveBootstrapIgnition is deprecated. Use bestEffortDeleteIgnition instead.
	// +optional
	PreserveBootstrapIgnition bool `json:"preserveBootstrapIgnition,omitempty"`

	// BestEffortDeleteIgnition is an optional field that can be used to ignore errors from S3 deletion of ignition
	// objects during cluster bootstrap. The default behavior is to fail the installation if ignition objects cannot be
	// deleted. Enable this functionality when there are known reasons disallowing their deletion.
	// +optional
	BestEffortDeleteIgnition bool `json:"bestEffortDeleteIgnition,omitempty"`

	// PublicIpv4Pool is an optional field that can be used to tell the installation process to use
	// Public IPv4 address that you bring to your AWS account with BYOIP.
	// +optional
	PublicIpv4Pool string `json:"publicIpv4Pool,omitempty"`

	// UserProvisionedDNS indicates if the customer is providing their own DNS solution in place of the default
	// provisioned by the Installer.
	// +kubebuilder:default:="Disabled"
	// +default="Disabled"
	// +kubebuilder:validation:Enum="Enabled";"Disabled"
	UserProvisionedDNS dns.UserProvisionedDNS `json:"userProvisionedDNS,omitempty"`
}

// ServiceEndpoint store the configuration for services to
// override existing defaults of AWS Services.
type ServiceEndpoint struct {
	// Name is the name of the AWS service.
	// This must be provided and cannot be empty.
	Name string `json:"name"`

	// URL is fully qualified URI with scheme https, that overrides the default generated
	// endpoint for a client.
	// This must be provided and cannot be empty.
	//
	// +kubebuilder:validation:Pattern=`^https://`
	URL string `json:"url"`
}

// VPC configures the VPC for the cluster.
type VPC struct {
	// Subnets defines the subnets in an existing VPC and can optionally specify their intended roles.
	// If no roles are specified on any subnet, then the subnet roles are decided automatically.
	// In this case, the VPC must not contain any other non-cluster subnets without the kubernetes.io/cluster/<cluster-id> tag.
	//
	// For manually specified subnet role selection, each subnet must have at least one assigned role,
	// and the ClusterNode, BootstrapNode, IngressControllerLB, ControlPlaneExternalLB, and ControlPlaneInternalLB roles must be assigned to at least one subnet.
	// However, if the cluster scope is internal, then ControlPlaneExternalLB is not required.
	//
	// Subnets must contain unique IDs, and can include no more than 10 subnets with the IngressController role.
	//
	// Leave this field unset to have the installer create subnets in a new VPC on your behalf.
	//
	// +listType=atomic
	// +optional
	Subnets []Subnet `json:"subnets,omitempty"`
}

// Subnet specifies a subnet in an existing VPC and can optionally specify their intended roles.
type Subnet struct {
	// ID specifies the subnet ID of an existing subnet.
	// The subnet ID must start with "subnet-", consist only of alphanumeric characters,
	// and must be exactly 24 characters long.
	//
	// +required
	ID AWSSubnetID `json:"id"`

	// Roles specifies the roles (aka functions) that the subnet will provide in the cluster.
	// If no roles are specified on any subnet, then the subnet roles are decided automatically.
	// Each role must be unique.
	//
	// +kubebuilder:validation:MaxItems=5
	// +optional
	Roles []SubnetRole `json:"roles,omitempty"`
}

// SubnetRole specifies the role (aka function) that the subnet will provide in the cluster.
type SubnetRole struct {
	// Type specifies the type of role (aka function) that the subnet will provide in the cluster.
	// Role types include ClusterNode, EdgeNode, BootstrapNode, IngressControllerLB, ControlPlaneExternalLB, and ControlPlaneInternalLB.
	//
	// +required
	Type SubnetRoleType `json:"type"`
}

// AWSSubnetID is a reference to an AWS subnet ID.
// +kubebuilder:validation:MinLength=24
// +kubebuilder:validation:MaxLength=24
// +kubebuilder:validation:Pattern=`^subnet-[0-9A-Za-z]+$`
type AWSSubnetID string // nolint:revive

// SubnetRoleType defines the type of role (aka function) that the subnet will provide in the cluster.
// +kubebuilder:validation:Enum:="ClusterNode";"EdgeNode";"BootstrapNode";"IngressControllerLB";"ControlPlaneExternalLB";"ControlPlaneInternalLB"
type SubnetRoleType string

const (
	// ClusterNodeSubnetRole specifies subnets that will be used as subnets for the
	// control plane and compute nodes.
	ClusterNodeSubnetRole SubnetRoleType = "ClusterNode"

	// EdgeNodeSubnetRole specifies subnets that will be used as edge subnets residing
	// in Local or Wavelength Zones for edge compute nodes.
	EdgeNodeSubnetRole SubnetRoleType = "EdgeNode"

	// BootstrapNodeSubnetRole specifies subnets that will be used as subnets for the
	// bootstrap node used to create the cluster.
	BootstrapNodeSubnetRole SubnetRoleType = "BootstrapNode"

	// IngressControllerLBSubnetRole specifies subnets used by the default IngressController.
	IngressControllerLBSubnetRole SubnetRoleType = "IngressControllerLB"

	// ControlPlaneExternalLBSubnetRole specifies subnets used by the external control plane
	// load balancer that serves the Kubernetes API server.
	ControlPlaneExternalLBSubnetRole SubnetRoleType = "ControlPlaneExternalLB"

	// ControlPlaneInternalLBSubnetRole specifies subnets used by the internal control plane
	// load balancer that serves the Kubernetes API server.
	ControlPlaneInternalLBSubnetRole SubnetRoleType = "ControlPlaneInternalLB"
)

// IsPublicOnlySubnetsEnabled returns whether the public-only subnets feature has been enabled via env var.
func IsPublicOnlySubnetsEnabled() bool {
	// Even though this looks too simple for a function, it's better than having to update the logic everywhere it's
	// used in case we decide to check for specific values set in the env var.
	return os.Getenv("OPENSHIFT_INSTALL_AWS_PUBLIC_ONLY") != ""
}
