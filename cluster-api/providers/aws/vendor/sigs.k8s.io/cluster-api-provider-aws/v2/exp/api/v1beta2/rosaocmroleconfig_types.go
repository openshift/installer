/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// ROSAOCMRoleProfile defines the permission level for the OCM role
type ROSAOCMRoleProfile string

const (
	// ROSAOCMRoleProfileStandard provides standard OCM permissions
	ROSAOCMRoleProfileStandard ROSAOCMRoleProfile = "Standard"

	// ROSAOCMRoleProfileAdmin provides admin OCM permissions
	ROSAOCMRoleProfileAdmin ROSAOCMRoleProfile = "Admin"

	// ROSAOCMRoleProfileNoConsole provides minimal OCM permissions (cannot use console.redhat.com)
	ROSAOCMRoleProfileNoConsole ROSAOCMRoleProfile = "NoConsole"
)

// ROSAOCMRoleDeletionPolicy defines what happens to the OCM role when the CR is deleted.
type ROSAOCMRoleDeletionPolicy string

const (
	// ROSAOCMRoleDeletionPolicyDelete unlinks and deletes the OCM role when the CR is deleted.
	ROSAOCMRoleDeletionPolicyDelete ROSAOCMRoleDeletionPolicy = "Delete"

	// ROSAOCMRoleDeletionPolicyRetain keeps the OCM role intact when the CR is deleted.
	ROSAOCMRoleDeletionPolicyRetain ROSAOCMRoleDeletionPolicy = "Retain"
)

const (
	// ROSAOCMRoleConfigReadyCondition condition reports on the successful reconciliation of ROSAOCMRoleConfig.
	ROSAOCMRoleConfigReadyCondition = "ROSAOCMRoleConfigReady"

	// ROSAOCMRoleConfigDeletionFailedReason used to report failures while deleting ROSAOCMRoleConfig.
	ROSAOCMRoleConfigDeletionFailedReason = "DeletionFailed"

	// ROSAOCMRoleConfigReconciliationFailedReason used to report reconciliation failures.
	ROSAOCMRoleConfigReconciliationFailedReason = "ReconciliationFailed"

	// ROSAOCMRoleConfigDeletionStarted used to indicate that the deletion of ROSAOCMRoleConfig has started.
	ROSAOCMRoleConfigDeletionStarted = "DeletionStarted"

	// ROSAOCMRoleConfigCreatedReason used to indicate that the ROSAOCMRoleConfig has been created.
	ROSAOCMRoleConfigCreatedReason = "Created"

	// ROSAOCMRoleConfigLinkedReason used to indicate that the OCM role has been linked to the organization.
	ROSAOCMRoleConfigLinkedReason = "Linked"
)

// ROSAOCMRoleConfigSpec defines the desired state of ROSAOCMRoleConfig
type ROSAOCMRoleConfigSpec struct {
	// RolePrefix is the user-defined prefix for the OCM role name.
	// The final role name will be: {RolePrefix}-OCM-Role-{ExternalID}
	// where ExternalID is the organization's external identifier from OCM.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength:=12
	// +kubebuilder:validation:Pattern:=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="rolePrefix is immutable"
	RolePrefix string `json:"rolePrefix"`

	// Profile defines the permission level for the OCM role.
	// +kubebuilder:validation:Enum=Standard;Admin;NoConsole
	// +kubebuilder:default=Standard
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="profile is immutable"
	Profile ROSAOCMRoleProfile `json:"profile"`

	// PermissionsBoundaryARN is the ARN of the policy that is used to set the permissions boundary for the OCM role.
	// +optional
	PermissionsBoundaryARN string `json:"permissionsBoundaryARN,omitempty"`

	// Path is the IAM path for the OCM role.
	// +optional
	// +kubebuilder:validation:Pattern=`^\/.*\/$`
	Path string `json:"path,omitempty"`

	// IdentityRef is a reference to an identity to be used when reconciling the OCM Role Config.
	// If no identity is specified, the default identity for this controller will be used.
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// CredentialsSecretRef references a secret with necessary credentials to connect to the OCM API.
	// +optional
	CredentialsSecretRef *corev1.SecretReference `json:"credentialsSecretRef,omitempty"`

	// DeletionPolicy determines what happens to the OCM role when this CR is deleted.
	// Delete will unlink and delete the OCM role.
	// Retain will keep the OCM role intact.
	// This is useful when reprovisioning management clusters to avoid disrupting users in the same organization.
	// +kubebuilder:validation:Enum=Delete;Retain
	// +kubebuilder:default=Delete
	// +optional
	DeletionPolicy ROSAOCMRoleDeletionPolicy `json:"deletionPolicy,omitempty"`
}

// ROSAOCMRoleConfigStatus defines the observed state of ROSAOCMRoleConfig
type ROSAOCMRoleConfigStatus struct {
	// RoleARN is the ARN of the created OCM role.
	RoleARN string `json:"roleARN,omitempty"`

	// OrganizationID is the OCM organization ID that this role is linked to.
	OrganizationID string `json:"organizationID,omitempty"`

	// Conditions specifies the ROSAOCMRoleConfig conditions
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`
}

// ROSAOCMRoleConfig is the Schema for the rosaocmroleconfigs API
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosaocmroleconfigs,scope=Cluster,categories=cluster-api,shortName=rosaocmrole
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
type ROSAOCMRoleConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSAOCMRoleConfigSpec   `json:"spec,omitempty"`
	Status ROSAOCMRoleConfigStatus `json:"status,omitempty"`
}

// ROSAOCMRoleConfigList contains a list of ROSAOCMRoleConfig
// +kubebuilder:object:root=true
type ROSAOCMRoleConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSAOCMRoleConfig `json:"items"`
}

// SetConditions sets the conditions of the ROSAOCMRoleConfig.
func (r *ROSAOCMRoleConfig) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// GetConditions returns the observations of the operational state of the ROSAOCMRoleConfig resource.
func (r *ROSAOCMRoleConfig) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

func init() {
	SchemeBuilder.Register(&ROSAOCMRoleConfig{}, &ROSAOCMRoleConfigList{})
}
