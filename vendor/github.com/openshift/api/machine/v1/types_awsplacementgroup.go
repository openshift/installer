package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSPlacementGroup ensures that a placement group matching the given configuration exists within AWS
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.observedConfiguration.groupType",description="Placement Group Type"
// +kubebuilder:printcolumn:name="Management",type="string",JSONPath=".spec.managementSpec.managementState",description="Management State"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="EC2 Replicas within the Placement Group"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="AWSPlacementGroup age"
type AWSPlacementGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSPlacementGroupSpec   `json:"spec,omitempty"`
	Status AWSPlacementGroupStatus `json:"status,omitempty"`
}

type AWSPlacementGroupSpec struct {
	// AWSPlacementGroupManagementSpec defines the configuration for a managed or unmanaged placement group.
	// +kubebuilder:validation:Required
	ManagementSpec AWSPlacementGroupManagementSpec `json:"managementSpec"`

	// CredentialsSecret is a reference to the secret with AWS credentials. The secret must reside in the same namespace
	// as the AWSPlacementGroup resource. Otherwise, the controller will leverage the EC2 instance assigned IAM Role,
	// in OpenShift this will always be the Control Plane Machine IAM Role.
	// +optional
	CredentialsSecret *LocalSecretReference `json:"credentialsSecret,omitempty"`
}

// AWSPlacementGroupManagementSpec defines the configuration for a managed or unmanaged placement group.
// +union
type AWSPlacementGroupManagementSpec struct {
	// ManagementState determines whether the placement group is expected
	// to be managed by this CRD or whether it is user managed.
	// A managed placement group may be moved to unmanaged, however an unmanaged
	// group may not be moved back to managed.
	// +kubebuilder:validation:Required
	// +unionDiscriminator
	ManagementState ManagementState `json:"managementState"`

	// Managed defines the configuration for the placement groups to be created.
	// Updates to the configuration will not be observed as placement groups are immutable
	// after creation.
	// +optional
	Managed *ManagedAWSPlacementGroup `json:"managed,omitempty"`
}

// AWSPlacementGroupType represents the valid values for the Placement GroupType field.
type AWSPlacementGroupType string

const (
	// AWSClusterPlacementGroupType is the "Cluster" placement group type.
	// Cluster placement groups place instances close together to improve network latency and throughput.
	AWSClusterPlacementGroupType AWSPlacementGroupType = "Cluster"
	// AWSPartitionPlacementGroupType is the "Partition" placement group type.
	// Partition placement groups reduce the likelihood of hardware failures
	// disrupting your application's availability.
	// Partition placement groups are recommended for use with large scale
	// distributed and replicated workloads.
	AWSPartitionPlacementGroupType AWSPlacementGroupType = "Partition"
	// AWSSpreadPlacementGroupType is the "Spread" placement group type.
	// Spread placement groups place instances on distinct racks within the availability
	// zone. This ensures instances each have their own networking and power source
	// for maximum hardware fault tolerance.
	// Spread placement groups are recommended for a small number of critical instances
	// which must be kept separate from one another.
	// Using a Spread placement group imposes a limit of seven instances within
	// the placement group within a single availability zone.
	AWSSpreadPlacementGroupType AWSPlacementGroupType = "Spread"
)

// ManagedAWSPlacementGroup is a discriminated union of placement group configuration.
// +union
type ManagedAWSPlacementGroup struct {
	// GroupType specifies the type of AWS placement group to use for this Machine.
	// This parameter is only used when a Machine is being created and the named
	// placement group does not exist.
	// Valid values are "Cluster", "Partition", "Spread".
	// This value is required and, in case a placement group already exists, will be
	// validated against the existing placement group.
	// Note: If the value of this field is "Spread", Machines created within the group
	// may no have placement.tenancy set
	// to "dedicated".
	// +kubebuilder:validation:Enum:="Cluster";"Partition";"Spread"
	// +unionDiscriminator
	// +kubebuilder:validation:Required
	GroupType AWSPlacementGroupType `json:"groupType,omitempty"`

	// Partition defines the configuration of a partition placement group.
	// +optional
	Partition *AWSPartitionPlacement `json:"partition,omitempty"`
}

// AWSPartitionPlacement defines the configuration for partition placement groups.
type AWSPartitionPlacement struct {
	// Count specifies the number of partitions for a Partition placement
	// group. This value is only observed when creating a placement group and
	// only when the `groupType` is set to `Partition`.
	// Note the partition count of a placement group cannot be changed after creation.
	// If unset, AWS will provide a default partition count.
	// This default is currently 2.
	// Note: When using more than 2 partitions, the "dedicated" tenancy option on Machines
	// created within the group is unavailable.
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=7
	// +optional
	Count int32 `json:"count,omitempty"`
}

type AWSPlacementGroupStatus struct {
	// Conditions represents the observations of the AWSPlacementGroup's current state.
	// Known .status.conditions.type are: Ready, Deleting
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ExpiresAt identifies when the observed configuration is valid until.
	// The observed configuration should not be trusted if this time has passed.
	// The AWSPlacementGroup controller will attempt to update the status before it expires.
	// +optional
	ExpiresAt *metav1.Time `json:"expiresAt,omitempty"`

	// Replicas counts how many AWS EC2 instances are present within the placement group.
	// Note: This is a pointer to be able to distinguish between an empty placement group
	// and the status having not yet been observed.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// ManagementState determines whether the placement group is expected
	// to be managed by this CRD or whether it is user managed.
	// A managed placement group may be moved to unmanaged, however an unmanaged
	// group may not be moved back to managed.
	// This value is owned by the controller and may differ from the spec in cases
	// when a user attempts to manage a previously unmanaged placement group.
	// +optional
	ManagementState ManagementState `json:"managementState,omitempty"`

	// ObservedConfiguration represents the configuration present on the placement group on AWS.
	// +optional
	ObservedConfiguration ManagedAWSPlacementGroup `json:"observedConfiguration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSPlacementGroupList contains a list of AWSPlacementGroup
// Compatibility level 1: Stable within a major release for a minimum of 12 months or 3 minor releases (whichever is longer).
// +openshift:compatibility-gen:level=1
type AWSPlacementGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSPlacementGroup `json:"items"`
}
