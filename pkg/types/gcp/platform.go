package gcp

import (
	"fmt"

	"github.com/openshift/installer/pkg/types/dns"
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// ProjectID is the the project that will be used for the cluster.
	ProjectID string `json:"projectID"`

	// Region specifies the GCP region where the cluster will be created.
	Region string `json:"region"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on GCP for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network specifies an existing VPC where the cluster should be created
	// rather than provisioning a new one.
	// +optional
	Network string `json:"network,omitempty"`

	// NetworkProjectID specifies which project the network and subnets exist in when
	// they are not in the main ProjectID.
	// +optional
	NetworkProjectID string `json:"networkProjectID,omitempty"`

	// ControlPlaneSubnet is an existing subnet where the control plane will be deployed.
	// The value should be the name of the subnet.
	// +optional
	ControlPlaneSubnet string `json:"controlPlaneSubnet,omitempty"`

	// ComputeSubnet is an existing subnet where the compute nodes will be deployed.
	// The value should be the name of the subnet.
	// +optional
	ComputeSubnet string `json:"computeSubnet,omitempty"`

	// userLabels has additional keys and values that the installer will add as
	// labels to all resources that it creates on GCP. Resources created by the
	// cluster itself may not include these labels.
	UserLabels []UserLabel `json:"userLabels,omitempty"`

	// userTags has additional keys and values that the installer will add as
	// tags to all resources that it creates on GCP. Resources created by the
	// cluster itself may not include these tags. Tag key and tag value should
	// be the shortnames of the tag key and tag value resource.
	UserTags []UserTag `json:"userTags,omitempty"`

	// UserProvisionedDNS indicates if the customer is providing their own DNS solution in place of the default
	// provisioned by the Installer.
	// +kubebuilder:default:="Disabled"
	// +default="Disabled"
	// +kubebuilder:validation:Enum="Enabled";"Disabled"
	UserProvisionedDNS dns.UserProvisionedDNS `json:"userProvisionedDNS,omitempty"`

	// ServiceEndpoints list contains custom endpoints which will override default
	// service endpoint of GCP Services.
	// There must be only one ServiceEndpoint for a service.
	// +optional
	ServiceEndpoints []ServiceEndpoint `json:"serviceEndpoints,omitempty"`
}

// ServiceEndpoint stores the configuration for services to
// override existing defaults of GCP Services.
type ServiceEndpoint struct {
	// Name is the name of the GCP service.
	// This must be provided and cannot be empty.
	Name string `json:"name"`

	// URL is fully qualified URI with scheme https, that overrides the default generated
	// endpoint for a client.
	// This must be provided and cannot be empty.
	//
	// +kubebuilder:validation:Pattern=`^https://`
	URL string `json:"url"`
}

// UserLabel is a label to apply to GCP resources created for the cluster.
type UserLabel struct {
	// key is the key part of the label. A label key can have a maximum of 63 characters
	// and cannot be empty. Label must begin with a lowercase letter, and must contain
	// only lowercase letters, numeric characters, and the following special characters `_-`.
	Key string `json:"key"`

	// value is the value part of the label. A label value can have a maximum of 63 characters
	// and cannot be empty. Value must contain only lowercase letters, numeric characters, and
	// the following special characters `_-`.
	Value string `json:"value"`
}

// UserTag is a tag to apply to GCP resources created for the cluster.
type UserTag struct {
	// parentID is the ID of the hierarchical resource where the tags are defined,
	// e.g. at the Organization or the Project level. To find the Organization ID or Project ID refer to the following pages:
	// https://cloud.google.com/resource-manager/docs/creating-managing-organization#retrieving_your_organization_id,
	// https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects.
	// An OrganizationID must consist of decimal numbers, and cannot have leading zeroes.
	// A ProjectID must be 6 to 30 characters in length, can only contain lowercase letters,
	// numbers, and hyphens, and must start with a letter, and cannot end with a hyphen.
	ParentID string `json:"parentID"`

	// key is the key part of the tag. A tag key can have a maximum of 63 characters and
	// cannot be empty. Tag key must begin and end with an alphanumeric character, and
	// must contain only uppercase, lowercase alphanumeric characters, and the following
	// special characters `._-`.
	Key string `json:"key"`

	// value is the value part of the tag. A tag value can have a maximum of 63 characters
	// and cannot be empty. Tag value must begin and end with an alphanumeric character, and
	// must contain only uppercase, lowercase alphanumeric characters, and the following
	// special characters `_-.@%=+:,*#&(){}[]` and spaces.
	Value string `json:"value"`
}

// DefaultSubnetName sets a default name for the subnet.
func DefaultSubnetName(infraID, role string) string {
	return fmt.Sprintf("%s-%s-subnet", infraID, role)
}

// GetConfiguredServiceAccount returns the service account email from a configured service account for
// a control plane or compute node. Returns empty string if not configured.
func GetConfiguredServiceAccount(platform *Platform, mpool *MachinePool) string {
	if mpool != nil && mpool.ServiceAccount != "" {
		return mpool.ServiceAccount
	} else if platform.DefaultMachinePlatform != nil {
		return platform.DefaultMachinePlatform.ServiceAccount
	}

	return ""
}

// GetDefaultServiceAccount returns the default service account email to use based on role.
// The default should be used when an existing service account is not configured.
func GetDefaultServiceAccount(platform *Platform, clusterID string, role string) string {
	return fmt.Sprintf("%s-%s@%s.iam.gserviceaccount.com", clusterID, role[0:1], platform.ProjectID)
}
