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

// TrunkSubportSpec represents a subport to attach to a trunk.
// It maps to gophercloud's trunks.Subport.
type TrunkSubportSpec struct {
	// portRef is a reference to the ORC Port that will be attached as a subport.
	// +required
	PortRef KubernetesNameRef `json:"portRef,omitempty"`

	// segmentationID is the segmentation ID for the subport (e.g. VLAN ID).
	// +required
	// +kubebuilder:validation:Minimum:=1
	// +kubebuilder:validation:Maximum:=4094
	SegmentationID int32 `json:"segmentationID,omitempty"`

	// segmentationType is the segmentation type for the subport (e.g. vlan).
	// +required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=32
	// +kubebuilder:validation:Enum:=inherit;vlan
	SegmentationType string `json:"segmentationType,omitempty"`
}

// TrunkSubportStatus represents an attached subport on a trunk.
// It maps to gophercloud's trunks.Subport.
type TrunkSubportStatus struct {
	// portID is the OpenStack ID of the Port attached as a subport.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	PortID string `json:"portID,omitempty"`

	// segmentationID is the segmentation ID for the subport (e.g. VLAN ID).
	// +optional
	SegmentationID int32 `json:"segmentationID,omitempty"`

	// segmentationType is the segmentation type for the subport (e.g. vlan).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	SegmentationType string `json:"segmentationType,omitempty"`
}

// TrunkResourceSpec contains the desired state of the resource.
type TrunkResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// portRef is a reference to the ORC Port which this resource is associated with.
	// +required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="portRef is immutable"
	PortRef KubernetesNameRef `json:"portRef,omitempty"`

	// projectRef is a reference to the ORC Project which this resource is associated with.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="projectRef is immutable"
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// adminStateUp is the administrative state of the trunk. If false (down),
	// the trunk does not forward packets.
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// subports is the list of ports to attach to the trunk.
	// +optional
	// +kubebuilder:validation:MaxItems:=1024
	// +listType=atomic
	Subports []TrunkSubportSpec `json:"subports,omitempty"`

	// tags is a list of Neutron tags to apply to the trunk.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=set
	// +optional
	Tags []NeutronTag `json:"tags,omitempty"`
}

// TrunkFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type TrunkFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +optional
	Description *NeutronDescription `json:"description,omitempty"`

	// portRef is a reference to the ORC Port which this resource is associated with.
	// +optional
	PortRef *KubernetesNameRef `json:"portRef,omitempty"`

	// projectRef is a reference to the ORC Project which this resource is associated with.
	// +optional
	ProjectRef *KubernetesNameRef `json:"projectRef,omitempty"`

	// Contrary to what the neutron doc say, we can't filter by status
	// https://github.com/gophercloud/gophercloud/issues/3626

	// adminStateUp is the administrative state of the trunk.
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	FilterByNeutronTags `json:",inline"`
}

// TrunkResourceStatus represents the observed state of the resource.
type TrunkResourceStatus struct {
	// name is a Human-readable name for the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// portID is the ID of the Port to which the resource is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	PortID string `json:"portID,omitempty"`

	// projectID is the ID of the Project to which the resource is associated.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	ProjectID string `json:"projectID,omitempty"`

	// tenantID is the project owner of the trunk (alias of projectID in some deployments).
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	TenantID string `json:"tenantID,omitempty"`

	// status indicates whether the trunk is currently operational.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Status string `json:"status,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`

	NeutronStatusMetadata `json:",inline"`

	// adminStateUp is the administrative state of the trunk.
	// +optional
	AdminStateUp *bool `json:"adminStateUp,omitempty"`

	// subports is a list of ports associated with the trunk.
	// +kubebuilder:validation:MaxItems=1024
	// +listType=atomic
	// +optional
	Subports []TrunkSubportStatus `json:"subports,omitempty"`
}
