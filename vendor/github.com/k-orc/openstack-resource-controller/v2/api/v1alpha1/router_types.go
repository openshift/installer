/*
Copyright 2024 The ORC Authors.

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

// RouterFilter specifies a query to select an OpenStack router. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type RouterFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

type ExternalGateway struct {
	// networkRef is a reference to the ORC Network which the external
	// gateway is on.
	// +required
	NetworkRef KubernetesNameRef `json:"networkRef,omitempty"`
}

type ExternalGatewayStatus struct {
	// networkID is the ID of the network the gateway is on.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NetworkID string `json:"networkID,omitempty"`
}

type RouterResourceSpec struct {
	// name is a human-readable name of the router. If not set, the
	// object's name will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// tags is a list of tags which will be applied to the router.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// adminStateUp represents the administrative state of the resource,
	// which is up (true) or down (false). Default is true.
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// externalGateways is a list of external gateways for the router.
	// Multiple gateways are not currently supported by ORC.
	// +kubebuilder:validation:MaxItems:=1
	// +listType=atomic
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="externalGateways is immutable"
	ExternalGateways []ExternalGateway `json:"externalGateways,omitempty"`

	// distributed indicates whether the router is distributed or not. It
	// is available when dvr extension is enabled.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="distributed is immutable"
	Distributed *bool `json:"distributed,omitempty"`

	// availabilityZoneHints is the availability zone candidate for the router.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="availabilityZoneHints is immutable"
	AvailabilityZoneHints []AvailabilityZoneHint `json:"availabilityZoneHints,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`
}

type RouterResourceStatus struct {
	// name is the human-readable name of the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// projectID is the project owner of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// status indicates the current status of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems:=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`

	// adminStateUp is the administrative state of the router,
	// which is up (true) or down (false).
	// +optional
	AdminStateUp *bool `json:"adminStateUp"`

	// externalGateways is a list of external gateways for the router.
	// +kubebuilder:validation:MaxItems:=32
	// +listType=atomic
	// +optional
	ExternalGateways []ExternalGatewayStatus `json:"externalGateways,omitempty"`

	// availabilityZoneHints is the availability zone candidate for the
	// router.
	// +kubebuilder:validation:MaxItems:=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	AvailabilityZoneHints []string `json:"availabilityZoneHints,omitempty"`
}
