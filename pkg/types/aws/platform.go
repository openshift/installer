package aws

import (
	"github.com/aws/aws-sdk-go/aws/endpoints"

	configv1 "github.com/openshift/api/config/v1"
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
	// AMIID is the AMI that should be used to boot machines for the cluster.
	// If set, the AMI should belong to the same region as the cluster.
	//
	// +optional
	AMIID string `json:"amiID,omitempty"`

	// Region specifies the AWS region where the cluster will be created.
	Region string `json:"region"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	//
	// +optional
	Subnets []string `json:"subnets,omitempty"`

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

	// The field is deprecated. ExperimentalPropagateUserTags is an experimental
	// flag that directs in-cluster operators to include the specified
	// user tags in the tags of the AWS resources that the operators create.
	// +optional
	ExperimentalPropagateUserTag *bool `json:"experimentalPropagateUserTags,omitempty"`

	// PropagateUserTags is a flag that directs in-cluster operators
	// to include the specified user tags in the tags of the
	// AWS resources that the operators create.
	// +optional
	PropagateUserTag bool `json:"propagateUserTags,omitempty"`

	// LBType is an optional field to specify a load balancer type.
	//
	// When this field is specified, the default ingresscontroller will be
	// created using the specified load-balancer type.
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

// IsSecretRegion returns true if the region is part of either the ISO or ISOB partitions.
func IsSecretRegion(region string) bool {
	partition, ok := endpoints.PartitionForRegion(endpoints.DefaultPartitions(), region)
	if !ok {
		return false
	}
	switch partition.ID() {
	case endpoints.AwsIsoPartitionID, endpoints.AwsIsoBPartitionID:
		return true
	}
	return false
}
