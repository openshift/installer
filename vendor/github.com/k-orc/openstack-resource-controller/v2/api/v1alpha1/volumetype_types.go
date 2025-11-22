/*
Copyright 2025 The ORC Authors.

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

// VolumeTypeResourceSpec contains the desired state of the resource.
type VolumeTypeResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Description *string `json:"description,omitempty"`

	// extraSpecs is a map of key-value pairs that define extra specifications for the volume type.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=atomic
	// +optional
	ExtraSpecs []VolumeTypeExtraSpec `json:"extraSpecs,omitempty"`

	// isPublic indicates whether the volume type is public.
	// +optional
	IsPublic *bool `json:"isPublic,omitempty"`
}

// VolumeTypeFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type VolumeTypeFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description of the existing resource
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Description *string `json:"description,omitempty"`

	// isPublic indicates whether the VolumeType is public.
	// +optional
	IsPublic *bool `json:"isPublic,omitempty"`
}

// VolumeTypeResourceStatus represents the observed state of the resource.
type VolumeTypeResourceStatus struct {
	// name is a Human-readable name for the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// extraSpecs is a map of key-value pairs that define extra specifications for the volume type.
	// +kubebuilder:validation:MaxItems:=64
	// +listType=atomic
	// +optional
	ExtraSpecs []VolumeTypeExtraSpecStatus `json:"extraSpecs"`

	// isPublic indicates whether the VolumeType is public.
	// +optional
	IsPublic *bool `json:"isPublic"`
}

type VolumeTypeExtraSpec struct {
	// name is the name of the extraspec
	// +kubebuilder:validation:MaxLength:=255
	// +required
	Name string `json:"name"`

	// value is the value of the extraspec
	// +kubebuilder:validation:MaxLength:=255
	// +required
	Value string `json:"value"`
}

type VolumeTypeExtraSpecStatus struct {
	// name is the name of the extraspec
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Name string `json:"name,omitempty"`

	// value is the value of the extraspec
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Value string `json:"value,omitempty"`
}
