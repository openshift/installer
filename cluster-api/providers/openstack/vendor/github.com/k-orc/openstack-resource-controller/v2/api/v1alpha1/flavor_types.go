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

// FlavorResourceSpec contains the desired state of a flavor
// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="FlavorResourceSpec is immutable"
type FlavorResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// description contains a free form description of the flavor.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=65535
	// +optional
	Description *string `json:"description,omitempty"`

	// ram is the memory of the flavor, measured in MB.
	// +kubebuilder:validation:Minimum=1
	// +required
	RAM int32 `json:"ram,omitempty"`

	// vcpus is the number of vcpus for the flavor.
	// +kubebuilder:validation:Minimum=1
	// +required
	Vcpus int32 `json:"vcpus,omitempty"`

	// disk is the size of the root disk that will be created in GiB. If 0
	// the root disk will be set to exactly the size of the image used to
	// deploy the instance. However, in this case the scheduler cannot
	// select the compute host based on the virtual image size. Therefore,
	// 0 should only be used for volume booted instances or for testing
	// purposes. Volume-backed instances can be enforced for flavors with
	// zero root disk via the
	// os_compute_api:servers:create:zero_disk_flavor policy rule.
	// +kubebuilder:validation:Minimum=0
	// +required
	Disk int32 `json:"disk"`

	// swap is the size of a dedicated swap disk that will be allocated, in
	// MiB. If 0 (the default), no dedicated swap disk will be created.
	// +kubebuilder:validation:Minimum=0
	// +optional
	Swap int32 `json:"swap,omitempty"`

	// isPublic flags a flavor as being available to all projects or not.
	// +optional
	IsPublic *bool `json:"isPublic,omitempty"`

	// ephemeral is the size of the ephemeral disk that will be created, in GiB.
	// Ephemeral disks may be written over on server state changes. So should only
	// be used as a scratch space for applications that are aware of its
	// limitations. Defaults to 0.
	// +kubebuilder:validation:Minimum=0
	// +optional
	Ephemeral int32 `json:"ephemeral,omitempty"`
}

// FlavorFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type FlavorFilter struct {
	// name of the existing resource
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// ram is the memory of the flavor, measured in MB.
	// +kubebuilder:validation:Minimum=1
	// +optional
	RAM *int32 `json:"ram,omitempty"`

	// vcpus is the number of vcpus for the flavor.
	// +kubebuilder:validation:Minimum=1
	// +optional
	Vcpus *int32 `json:"vcpus,omitempty"`

	// disk is the size of the root disk in GiB.
	// +kubebuilder:validation:Minimum=0
	// +optional
	Disk *int32 `json:"disk,omitempty"`
}

// FlavorResourceStatus represents the observed state of the resource.
type FlavorResourceStatus struct {
	// name is a Human-readable name for the flavor. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength:=65535
	// +optional
	Description string `json:"description,omitempty"`

	// ram is the memory of the flavor, measured in MB.
	// +optional
	RAM *int32 `json:"ram,omitempty"`

	// vcpus is the number of vcpus for the flavor.
	// +optional
	Vcpus *int32 `json:"vcpus,omitempty"`

	// disk is the size of the root disk that will be created in GiB.
	// +optional
	Disk *int32 `json:"disk,omitempty"`

	// swap is the size of a dedicated swap disk that will be allocated, in
	// MiB.
	// +optional
	Swap *int32 `json:"swap,omitempty"`

	// isPublic flags a flavor as being available to all projects or not.
	// +optional
	IsPublic *bool `json:"isPublic,omitempty"`

	// ephemeral is the size of the ephemeral disk, in GiB.
	// +optional
	Ephemeral *int32 `json:"ephemeral,omitempty"`
}
