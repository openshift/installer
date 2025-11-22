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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=255
type NeutronDescription string

// NeutronTag represents a tag on a Neutron resource.
// It may not be empty and may not contain commas.
// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=255
type NeutronTag string

type FilterByNeutronTags struct {
	// tags is a list of tags to filter by. If specified, the resource must
	// have all of the tags specified to be included in the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=64
	Tags []NeutronTag `json:"tags,omitempty"`

	// tagsAny is a list of tags to filter by. If specified, the resource
	// must have at least one of the tags specified to be included in the
	// result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=64
	TagsAny []NeutronTag `json:"tagsAny,omitempty"`

	// notTags is a list of tags to filter by. If specified, resources which
	// contain all of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=64
	NotTags []NeutronTag `json:"notTags,omitempty"`

	// notTagsAny is a list of tags to filter by. If specified, resources
	// which contain any of the given tags will be excluded from the result.
	// +listType=set
	// +optional
	// +kubebuilder:validation:MaxItems:=64
	NotTagsAny []NeutronTag `json:"notTagsAny,omitempty"`
}

// +kubebuilder:validation:Enum:=4;6
type IPVersion int32

// +kubebuilder:validation:Format:=cidr
// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=49
type CIDR string

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=45
type IPvAny string

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=17
type MAC string

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=255
type AvailabilityZoneHint string

type NeutronStatusMetadata struct {
	// createdAt shows the date and time when the resource was created. The date and time stamp format is ISO 8601
	// +optional
	CreatedAt *metav1.Time `json:"createdAt,omitempty"`
	// updatedAt shows the date and time when the resource was updated. The date and time stamp format is ISO 8601
	// +optional
	UpdatedAt *metav1.Time `json:"updatedAt,omitempty"`

	// revisionNumber optionally set via extensions/standard-attr-revisions
	// +optional
	RevisionNumber *int64 `json:"revisionNumber,omitempty"`
}

// +kubebuilder:validation:MinLength:=1
// +kubebuilder:validation:MaxLength:=253
type KubernetesNameRef string
