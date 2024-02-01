package gcp

import (
	"fmt"
)

// UserProvisionedDNS indicates whether the DNS solution is provisioned by the Installer or the user.
type UserProvisionedDNS string

const (
	// UserProvisionedDNSEnabled indicates that the DNS solution is provisioned and provided by the user.
	UserProvisionedDNSEnabled UserProvisionedDNS = "Enabled"

	// UserProvisionedDNSDisabled indicates that the DNS solution is provisioned by the Installer.
	UserProvisionedDNSDisabled UserProvisionedDNS = "Disabled"
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
	// cluster itself may not include these labels. This is a TechPreview feature
	// and requires setting CustomNoUpgrade featureSet with GCPLabelsTags featureGate
	// enabled or TechPreviewNoUpgrade featureSet to configure labels.
	UserLabels []UserLabel `json:"userLabels,omitempty"`

	// userTags has additional keys and values that the installer will add as
	// tags to all resources that it creates on GCP. Resources created by the
	// cluster itself may not include these tags. Tag key and tag value should
	// be the shortnames of the tag key and tag value resource. This is a TechPreview
	// feature and requires setting CustomNoUpgrade featureSet with GCPLabelsTags
	// featureGate enabled or TechPreviewNoUpgrade featureSet to configure tags.
	UserTags []UserTag `json:"userTags,omitempty"`

	// UserProvisionedDNS indicates if the customer is providing their own DNS solution in place of the default
	// provisioned by the Installer.
	// +kubebuilder:default:="Disabled"
	// +default="Disabled"
	// +kubebuilder:validation:Enum="Enabled";"Disabled"
	UserProvisionedDNS UserProvisionedDNS `json:"userProvisionedDNS,omitempty"`
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
