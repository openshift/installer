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

// DomainResourceSpec contains the desired state of the resource.
type DomainResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *KeystoneName `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=255
	// +optional
	Description *string `json:"description,omitempty"`

	// enabled defines whether a domain is enabled or not. Default is true.
	// Note: Users can only authorize against an enabled domain (and any of its projects).
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// DomainFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type DomainFilter struct {
	// name of the existing resource
	// +optional
	Name *KeystoneName `json:"name,omitempty"`

	// enabled defines whether a domain is enabled or not. Default is true.
	// Note: Users can only authorize against an enabled domain (and any of its projects).
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// DomainResourceStatus represents the observed state of the resource.
type DomainResourceStatus struct {
	// name is a Human-readable name for the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Description string `json:"description,omitempty"`

	// enabled defines whether a domain is enabled or not. Default is true.
	// Note: Users can only authorize against an enabled domain (and any of its projects).
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}
