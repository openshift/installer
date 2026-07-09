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

// UserResourceSpec contains the desired state of the resource.
// +kubebuilder:validation:XValidation:rule="!has(oldSelf.passwordRef) || has(self.passwordRef)",message="passwordRef may not be removed once set"
type UserResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Description *string `json:"description,omitempty"`

	// domainRef is a reference to the ORC Domain which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="domainRef is immutable"
	DomainRef *KubernetesNameRef `json:"domainRef,omitempty"`

	// defaultProjectRef is a reference to the Default Project which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="defaultProjectRef is immutable"
	DefaultProjectRef *KubernetesNameRef `json:"defaultProjectRef,omitempty"`

	// enabled defines whether a user is enabled or disabled
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// passwordRef is a reference to a Secret containing the password
	// for this user. The Secret must contain a key named "password".
	// If not specified, the user is created without a password.
	// +optional
	PasswordRef *KubernetesNameRef `json:"passwordRef,omitempty"`
}

// UserFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type UserFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// domainRef is a reference to the ORC Domain which this resource is associated with.
	// +optional
	DomainRef *KubernetesNameRef `json:"domainRef,omitempty"`
}

// UserResourceStatus represents the observed state of the resource.
type UserResourceStatus struct {
	// name is a Human-readable name for the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// domainID is the ID of the Domain to which the resource is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	DomainID string `json:"domainID,omitempty"`

	// defaultProjectID is the ID of the Default Project to which the user is associated with.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	DefaultProjectID string `json:"defaultProjectID,omitempty"`

	// enabled defines whether a user is enabled or disabled
	// +optional
	Enabled bool `json:"enabled,omitempty"`

	// passwordExpiresAt is the timestamp at which the user's password expires.
	// +kubebuilder:validation:MaxLength:=1024
	// +optional
	PasswordExpiresAt string `json:"passwordExpiresAt,omitempty"`

	// appliedPasswordRef is the name of the Secret containing the
	// password that was last applied to the OpenStack resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	AppliedPasswordRef string `json:"appliedPasswordRef,omitempty"`
}
