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

// FloatingIPFilter specifies a query to select an OpenStack floatingip. At least one property must be set.
// +kubebuilder:validation:MinProperties:=1
type FloatingIPFilter struct {
	// floatingIP is the floatingip address.
	// +optional
	FloatingIP *IPvAny `json:"floatingIP,omitempty"`

	// description of the existing resource
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// floatingNetworkRef is a reference to the ORC Network which this resource is associated with.
	// +optional
	FloatingNetworkRef *KubernetesNameRef `json:"floatingNetworkRef,omitempty"`

	// portRef is a reference to the ORC Port which this resource is associated with.
	// +optional
	PortRef *KubernetesNameRef `json:"portRef,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// status is the status of the floatingip.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

// FloatingIPResourceSpec contains the desired state of a floating IP
// +kubebuilder:validation:XValidation:rule="has(self.floatingNetworkRef) != has(self.floatingSubnetRef)",message="Exactly one of 'floatingNetworkRef' or 'floatingSubnetRef' must be set"
type FloatingIPResourceSpec struct {
	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// tags is a list of tags which will be applied to the floatingip.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`

	// floatingNetworkRef references the network to which the floatingip is associated.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="floatingNetworkRef is immutable"
	FloatingNetworkRef *KubernetesNameRef `json:"floatingNetworkRef,omitempty"`

	// floatingSubnetRef references the subnet to which the floatingip is associated.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="floatingSubnetRef is immutable"
	FloatingSubnetRef *KubernetesNameRef `json:"floatingSubnetRef,omitempty"`

	// floatingIP is the IP that will be assigned to the floatingip. If not set, it will
	// be assigned automatically.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="floatingIP is immutable"
	FloatingIP *IPvAny `json:"floatingIP,omitempty"`

	// portRef is a reference to the ORC Port which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="portRef is immutable"
	PortRef *KubernetesNameRef `json:"portRef,omitempty"`

	// fixedIP is the IP address of the port to which the floatingip is associated.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="fixedIP is immutable"
	FixedIP *IPvAny `json:"fixedIP,omitempty"`

	// projectRef is a reference to the ORC Project this resource is associated with.
	// Typically, only used by admin.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`
}

type FloatingIPResourceStatus struct {
	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// floatingNetworkID is the ID of the network to which the floatingip is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	FloatingNetworkID string `json:"floatingNetworkID,omitempty"`

	// floatingIP is the IP address of the floatingip.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	FloatingIP string `json:"floatingIP,omitempty"`

	// portID is the ID of the port to which the floatingip is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	PortID string `json:"portID,omitempty"`

	// fixedIP is the IP address of the port to which the floatingip is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	FixedIP string `json:"fixedIP,omitempty"`

	// tenantID is the project owner of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	TenantID string `json:"tenantID,omitempty"`

	// projectID is the project owner of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// status indicates the current status of the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	// routerID is the ID of the router to which the floatingip is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	RouterID string `json:"routerID,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems:=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`

	NeutronStatusMetadata `json:",inline"`
}
