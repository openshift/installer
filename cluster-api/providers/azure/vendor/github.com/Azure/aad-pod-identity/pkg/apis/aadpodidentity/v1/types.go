package v1

import (
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CRDLabelKey is the
	CRDLabelKey = "aadpodidbinding"

	// BehaviorKey is the key that describes the behavior of aad-pod-identity.
	// Supported values:
	// namespaced - used for running in namespaced mode. AzureIdentity,
	//              AzureIdentityBinding and pod in the same namespace
	//              will only be matched for this behavior.
	BehaviorKey = "aadpodidentity.k8s.io/Behavior"

	// BehaviorNamespaced indicates that aad-pod-identity is behaving in namespaced mode.
	BehaviorNamespaced = "namespaced"

	// AssignedIDCreated indicates that an AzureAssignedIdentity is created.
	AssignedIDCreated = "Created"

	// AssignedIDAssigned indicates that an identity has been assigned to the node.
	AssignedIDAssigned = "Assigned"

	// AssignedIDUnAssigned indicates that an identity has been unassigned from the node.
	AssignedIDUnAssigned = "Unassigned"
)

// AzureIdentity is the specification of the identity data structure.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type",description="",priority=0
// +kubebuilder:printcolumn:name="ClientID",type="string",JSONPath=".spec.clientID",description="",priority=0
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC."
type AzureIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureIdentitySpec   `json:"spec,omitempty"`
	Status AzureIdentityStatus `json:"status,omitempty"`
}

// AzureIdentityBinding brings together the spec of matching pods and the identity which they can use.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:printcolumn:name="AzureIdentity",type="string",JSONPath=".spec.azureIdentity",description="",priority=0
// +kubebuilder:printcolumn:name="Selector",type="string",JSONPath=".spec.selector",description="",priority=0
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC."
type AzureIdentityBinding struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureIdentityBindingSpec   `json:"spec,omitempty"`
	Status AzureIdentityBindingStatus `json:"status,omitempty"`
}

// AzureAssignedIdentity contains the identity <-> pod mapping which is matched.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureAssignedIdentity struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureAssignedIdentitySpec   `json:"spec,omitempty"`
	Status AzureAssignedIdentityStatus `json:"status,omitempty"`
}

// AzurePodIdentityException contains the pod selectors for all pods that don't require
// NMI to process and request token on their behalf.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzurePodIdentityException struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzurePodIdentityExceptionSpec   `json:"spec,omitempty"`
	Status AzurePodIdentityExceptionStatus `json:"status,omitempty"`
}

// AzureIdentityList contains a list of AzureIdentities.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AzureIdentity `json:"items"`
}

// AzureIdentityBindingList contains a list of AzureIdentityBindings.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureIdentityBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AzureIdentityBinding `json:"items"`
}

// AzureAssignedIdentityList contains a list of AzureAssignedIdentities.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzureAssignedIdentityList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AzureAssignedIdentity `json:"items"`
}

// AzurePodIdentityExceptionList contains a list of AzurePodIdentityExceptions.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AzurePodIdentityExceptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []AzurePodIdentityException `json:"items"`
}

// IdentityType represents different types of identities.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type IdentityType int

const (
	// UserAssignedMSI represents a user-assigned identity.
	UserAssignedMSI IdentityType = 0

	// ServicePrincipal represents a service principal.
	ServicePrincipal IdentityType = 1
)

// AzureIdentitySpec describes the credential specifications of an identity on Azure.
type AzureIdentitySpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// UserAssignedMSI or Service Principal
	Type IdentityType `json:"type,omitempty"`

	// User assigned MSI resource id.
	ResourceID string `json:"resourceID,omitempty"`
	// Both User Assigned MSI and SP can use this field.
	ClientID string `json:"clientID,omitempty"`

	// Used for service principal
	ClientPassword api.SecretReference `json:"clientPassword,omitempty"`
	// Service principal primary tenant id.
	TenantID string `json:"tenantID,omitempty"`
	// Service principal auxiliary tenant ids
	// +nullable
	AuxiliaryTenantIDs []string `json:"auxiliaryTenantIDs,omitempty"`
	// For service principal. Option param for specifying the  AD details.
	ADResourceID string `json:"adResourceID,omitempty"`
	ADEndpoint   string `json:"adEndpoint,omitempty"`

	// +nullable
	Replicas *int32 `json:"replicas,omitempty"`
}

// AzureIdentityStatus contains the replica status of the resource.
type AzureIdentityStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// AssignedIDState represents the state of an AzureAssignedIdentity
type AssignedIDState int

const (
	// Created - Default state of the assigned identity
	Created AssignedIDState = 0

	// Assigned - When the underlying platform assignment of
	// managed identity is complete, the state moves to assigned
	Assigned AssignedIDState = 1
)

const (
	// AzureIDResource is the name of AzureIdentity.
	AzureIDResource = "azureidentities"

	// AzureIDBindingResource is the name of AzureIdentityBinding.
	AzureIDBindingResource = "azureidentitybindings"

	// AzureAssignedIDResource is the name of AzureAssignedIdentity.
	AzureAssignedIDResource = "azureassignedidentities"

	// AzurePodIdentityExceptionResource is the name of AzureIdentityException.
	AzurePodIdentityExceptionResource = "azurepodidentityexceptions"
)

// AzureIdentityBindingSpec matches the pod with the Identity.
// Used to indicate the potential matches to look for between the pod/deployment
// and the identities present.
type AzureIdentityBindingSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AzureIdentity     string `json:"azureIdentity,omitempty"`
	Selector          string `json:"selector,omitempty"`
	// Weight is used to figure out which of the matching identities would be selected.
	Weight int `json:"weight,omitempty"`
}

// AzureIdentityBindingStatus contains the status of an AzureIdentityBinding.
type AzureIdentityBindingStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	AvailableReplicas int32 `json:"availableReplicas,omitempty"`
}

// AzureAssignedIdentitySpec contains the relationship
// between an AzureIdentity and an AzureIdentityBinding.
type AzureAssignedIdentitySpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// AzureIdentityRef is an embedded resource referencing the AzureIdentity used by the
	// AzureAssignedIdentity, which requires x-kubernetes-embedded-resource fields to be true
	// +kubebuilder:validation:XEmbeddedResource
	AzureIdentityRef *AzureIdentity `json:"azureIdentityRef,omitempty"`

	// AzureBindingRef is an embedded resource referencing the AzureIdentityBinding used by the
	// AzureAssignedIdentity, which requires x-kubernetes-embedded-resource fields to be true
	// +kubebuilder:validation:XEmbeddedResource
	AzureBindingRef *AzureIdentityBinding `json:"azureBindingRef,omitempty"`
	Pod             string                `json:"pod,omitempty"`
	PodNamespace    string                `json:"podNamespace,omitempty"`
	NodeName        string                `json:"nodename,omitempty"`

	// +nullable
	Replicas *int32 `json:"replicas,omitempty"`
}

// AzureAssignedIdentityStatus contains the replica status of the resource.
type AzureAssignedIdentityStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            string `json:"status,omitempty"`
	AvailableReplicas int32  `json:"availableReplicas,omitempty"`
}

// AzurePodIdentityExceptionSpec matches pods with the selector defined.
// If request originates from a pod that matches the selector, nmi will
// proxy the request and send response back without any validation.
type AzurePodIdentityExceptionSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	PodLabels         map[string]string `json:"podLabels,omitempty"`
}

// AzurePodIdentityExceptionStatus contains the status of an AzurePodIdentityException.
type AzurePodIdentityExceptionStatus struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            string `json:"status,omitempty"`
}
