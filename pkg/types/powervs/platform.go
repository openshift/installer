package powervs

import "k8s.io/apimachinery/pkg/util/sets"

var (
	// C2SRegions are the C2S AWS regions.
	C2SRegions = sets.NewString("us-iso-east-1")
)

// Platform stores all the global configuration that all machinesets
// use.
/// used by the installconfig, and filled in by the installconfig/platform/powervs::Platform() func
type Platform struct {
	// Region specifies the IBM Cloud region where the cluster will be created.
	Region string `json:"region"`

	// Subnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	// @TODO: how will we handle networking?
	//
	// +optional ?
	Subnets []string `json:"subnets,omitempty"`

	// HostedZone is the ID of an existing hosted zone into which to add DNS
	// records for the cluster's internal API. An existing hosted zone can
	// only be used when also using existing subnets. The hosted zone must be
	// associated with the VPC containing the subnets.
	// Leave the hosted zone unset to have the installer create the hosted zone
	// on your behalf.
	// +optional
	HostedZone string `json:"hostedZone,omitempty"`

	// UserTags additional keys and values that the installer will add
	// as tags to all resources that it creates. Resources created by the
	// cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

}

// ServiceEndpoint store the configuration for services to
// override existing defaults
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
