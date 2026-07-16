/*
Copyright The ORC Authors.

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

package v1alpha1

// RoleAssignmentResourceSpec defines the desired role assignment.
// A role assignment grants a role to a user or group on a project or domain.
// Role assignments are immutable once created and identified by the combination
// of (role, actor, scope) rather than a separate ID.
// +kubebuilder:validation:XValidation:rule="(has(self.userRef) && !has(self.groupRef)) || (!has(self.userRef) && has(self.groupRef))",message="exactly one of userRef or groupRef is required"
// +kubebuilder:validation:XValidation:rule="(has(self.projectRef) && !has(self.domainRef)) || (!has(self.projectRef) && has(self.domainRef))",message="exactly one of projectRef or domainRef is required"
// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="RoleAssignmentResourceSpec is immutable"
type RoleAssignmentResourceSpec struct {
	// roleRef references the Role being assigned.
	// +required
	RoleRef KubernetesNameRef `json:"roleRef,omitempty"`

	// userRef references the User receiving the role assignment.
	// Exactly one of userRef or groupRef must be specified.
	// +optional
	UserRef *KubernetesNameRef `json:"userRef,omitempty"`

	// groupRef references the Group receiving the role assignment.
	// Exactly one of userRef or groupRef must be specified.
	// +optional
	GroupRef *KubernetesNameRef `json:"groupRef,omitempty"`

	// projectRef references the Project scope for the assignment.
	// Exactly one of projectRef or domainRef must be specified.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// domainRef references the Domain scope for the assignment.
	// Exactly one of projectRef or domainRef must be specified.
	// +optional
	DomainRef *KubernetesNameRef `json:"domainRef,omitempty"`
}

// RoleAssignmentFilter defines import filter criteria for existing role assignments.
// +kubebuilder:validation:MinProperties:=1
type RoleAssignmentFilter struct {
	// roleRef filters by the referenced Role.
	// +optional
	RoleRef *KubernetesNameRef `json:"roleRef,omitempty"`

	// userRef filters by the referenced User.
	// +optional
	UserRef *KubernetesNameRef `json:"userRef,omitempty"`

	// groupRef filters by the referenced Group.
	// +optional
	GroupRef *KubernetesNameRef `json:"groupRef,omitempty"`

	// projectRef filters by the referenced Project scope.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// domainRef filters by the referenced Domain scope.
	// +optional
	DomainRef *KubernetesNameRef `json:"domainRef,omitempty"`
}

// RoleAssignmentResourceStatus represents the observed state of the role assignment.
// Note: Role assignments do not have a unique ID in OpenStack - they are identified
// by the combination of role, actor (user/group), and scope (project/domain).
type RoleAssignmentResourceStatus struct {
	// roleID is the OpenStack ID of the assigned role.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	RoleID string `json:"roleID,omitempty"`

	// userID is the OpenStack ID of the user (if actorType is User).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	UserID string `json:"userID,omitempty"`

	// groupID is the OpenStack ID of the group (if actorType is Group).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	GroupID string `json:"groupID,omitempty"`

	// projectID is the OpenStack ID of the project scope (if scopeType is Project).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// domainID is the OpenStack ID of the domain scope (if scopeType is Domain).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	DomainID string `json:"domainID,omitempty"`
}
