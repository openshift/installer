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

// AddressScopeResourceSpec contains the desired state of the resource.
type AddressScopeResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// projectRef is a reference to the ORC Project which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// ipVersion is the IP protocol version.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="ipVersion is immutable"
	IPVersion IPVersion `json:"ipVersion"`

	// shared indicates whether this resource is shared across all
	// projects or not. By default, only admin users can change set
	// this value. We can't unshared a shared address scope; Neutron
	// enforces this.
	// +optional
	// +kubebuilder:validation:XValidation:rule="!(oldSelf && !self)",message="shared address scope can't be unshared"
	Shared *bool `json:"shared,omitempty"`
}

// AddressScopeFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type AddressScopeFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// projectRef is a reference to the ORC Project which this resource is associated with.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// ipVersion is the IP protocol version.
	// +optional
	IPVersion IPVersion `json:"ipVersion,omitempty"`

	// shared indicates whether this resource is shared across all
	// projects or not. By default, only admin users can change set
	// this value.
	// +optional
	Shared *bool `json:"shared,omitempty"`
}

// AddressScopeResourceStatus represents the observed state of the resource.
type AddressScopeResourceStatus struct {
	// name is a Human-readable name for the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// projectID is the ID of the Project to which the resource is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// ipVersion is the IP protocol version.
	// +optional
	IPVersion int32 `json:"ipVersion,omitempty"`

	// shared indicates whether this resource is shared across all
	// projects or not. By default, only admin users can change set
	// this value.
	// +optional
	Shared *bool `json:"shared,omitempty"`
}
