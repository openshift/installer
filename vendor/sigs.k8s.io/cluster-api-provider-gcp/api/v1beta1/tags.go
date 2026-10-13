/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta1

// ResourceManagerTags is an slice of ResourceManagerTag structs.
type ResourceManagerTags []ResourceManagerTag

// ResourceManagerTagsMap defines a map of key value pairs as expected by compute.InstanceParams.ResourceManagerTags.
type ResourceManagerTagsMap map[string]string

// ResourceManagerTag is a tag to apply to GCP resources managed by the GCP provider.
type ResourceManagerTag struct {
	// ParentID is the ID of the hierarchical resource where the tags are defined
	// e.g. at the Organization or the Project level. To find the Organization or Project ID ref
	// https://cloud.google.com/resource-manager/docs/creating-managing-organization#retrieving_your_organization_id
	// https://cloud.google.com/resource-manager/docs/creating-managing-projects#identifying_projects
	// An OrganizationID must consist of decimal numbers, and cannot have leading zeroes.
	// A ProjectID must be 6 to 30 characters in length, can only contain lowercase letters,
	// numbers, and hyphens, and must start with a letter, and cannot end with a hyphen.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	// +kubebuilder:validation:Pattern=`(^[1-9][0-9]{0,31}$)|(^[a-z][a-z0-9-]{4,28}[a-z0-9]$)`
	ParentID string `json:"parentID"`

	// Key is the key part of the tag. A tag key can have a maximum of 63 characters and cannot
	// be empty. Tag key must begin and end with an alphanumeric character, and must contain
	// only uppercase, lowercase alphanumeric characters, and the following special
	// characters `._-`.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9]([0-9A-Za-z_.-]{0,61}[a-zA-Z0-9])?$`
	Key string `json:"key"`

	// Value is the value part of the tag. A tag value can have a maximum of 63 characters and
	// cannot be empty. Tag value must begin and end with an alphanumeric character, and must
	// contain only uppercase, lowercase alphanumeric characters, and the following special
	// characters `_-.@%=+:,*#&(){}[]` and spaces.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-zA-Z0-9]([0-9A-Za-z_.@%=+:,*#&()\[\]{}\-\s]{0,61}[a-zA-Z0-9])?$`
	Value string `json:"value"`
}

// Merge merges resource manager tags in receiver and other.
func (t *ResourceManagerTags) Merge(other ResourceManagerTags) {
	*t = append(*t, other...)
}
