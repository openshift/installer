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

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// ShareNetworkResourceSpec contains the desired state of the resource.
// +kubebuilder:validation:XValidation:rule="has(self.networkRef) == has(self.subnetRef)",message="networkRef and subnetRef must be specified together"
type ShareNetworkResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Description *string `json:"description,omitempty"`

	// networkRef is a reference to the ORC Network which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="networkRef is immutable"
	NetworkRef *KubernetesNameRef `json:"networkRef,omitempty"`

	// subnetRef is a reference to the ORC Subnet which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="subnetRef is immutable"
	SubnetRef *KubernetesNameRef `json:"subnetRef,omitempty"`
}

// ShareNetworkFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type ShareNetworkFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Description *string `json:"description,omitempty"`
}

// ShareNetworkResourceStatus represents the observed state of the resource.
type ShareNetworkResourceStatus struct {
	// name is a Human-readable name for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// neutronNetID is the Neutron network ID.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NeutronNetID string `json:"neutronNetID,omitempty"`

	// neutronSubnetID is the Neutron subnet ID.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NeutronSubnetID string `json:"neutronSubnetID,omitempty"`

	// networkType is the network type (e.g., vlan, vxlan, flat).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	NetworkType string `json:"networkType,omitempty"`

	// segmentationID is the segmentation ID of the network.
	// +optional
	SegmentationID *int32 `json:"segmentationID,omitempty"`

	// cidr is the CIDR of the subnet.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	CIDR string `json:"cidr"`

	// ipVersion is the IP version (4 or 6).
	// +optional
	IPVersion *int32 `json:"ipVersion,omitempty"`

	// projectID is the ID of the project that owns the share network.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// createdAt shows the date and time when the resource was created.
	// +optional
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`

	// updatedAt shows the date and time when the resource was updated.
	// +optional
	UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`
}
