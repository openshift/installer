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

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=255
type KeystoneTag string

type FilterByKeystoneTags struct {
	// tags is a list of tags to filter by. If specified, the resource must
	// have all of the tags specified to be included in the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=80
	Tags []KeystoneTag `json:"tags,omitempty"`

	// tagsAny is a list of tags to filter by. If specified, the resource
	// must have at least one of the tags specified to be included in the
	// result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=80
	TagsAny []KeystoneTag `json:"tagsAny,omitempty"`

	// notTags is a list of tags to filter by. If specified, resources which
	// contain all of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=80
	NotTags []KeystoneTag `json:"notTags,omitempty"`

	// notTagsAny is a list of tags to filter by. If specified, resources
	// which contain any of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=80
	NotTagsAny []KeystoneTag `json:"notTagsAny,omitempty"`
}

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=64
type KeystoneName string

// ProjectResourceSpec contains the desired state of a project
type ProjectResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *KeystoneName `json:"name,omitempty"`

	// description contains a free form description of the project.
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=65535
	// +optional
	Description *string `json:"description,omitempty"`

	// enabled defines whether a project is enabled or not. Default is true.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// tags is list of simple strings assigned to a project.
	// Tags can be used to classify projects into groups.
	// +kubebuilder:validation:MaxItems:=80
	// +listType=set
	// +optional
	Tags []KeystoneTag `json:"tags,omitempty"`
}

// ProjectFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type ProjectFilter struct {
	// name of the existing resource
	// +optional
	Name *KeystoneName `json:"name,omitempty"`

	FilterByKeystoneTags `json:",inline"`
}

// ProjectResourceStatus represents the observed state of the resource.
type ProjectResourceStatus struct {
	// name is a Human-readable name for the project. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// description is a human-readable description for the resource.
	// +kubebuilder:validation:MaxLength:=65535
	// +optional
	Description string `json:"description,omitempty"`

	// enabled represents whether a project is enabled or not.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// tags is the list of tags on the resource.
	// +kubebuilder:validation:MaxItems=80
	// +kubebuilder:validation:items:MaxLength=1024
	// +listType=atomic
	// +optional
	Tags []string `json:"tags,omitempty"`
}
