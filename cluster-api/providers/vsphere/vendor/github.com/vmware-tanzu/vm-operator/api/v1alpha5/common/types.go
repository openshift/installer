// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

// +kubebuilder:object:generate=true

package common

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// LocalObjectRef describes a reference to another object in the same
// namespace as the referrer.
type LocalObjectRef struct {
	// APIVersion defines the versioned schema of this representation of an
	// object. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	APIVersion string `json:"apiVersion"`

	// Kind is a string value representing the REST resource this object
	// represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	Kind string `json:"kind"`

	// Name refers to a unique resource in the current namespace.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name"`
}

// PartialObjectRef describes a reference to another object in the same
// namespace as the referrer. The reference can be just a name but may also
// include the referred resource's APIVersion and Kind.
type PartialObjectRef struct {
	metav1.TypeMeta `json:",inline"`

	// Name refers to a unique resource in the current namespace.
	// More info: http://kubernetes.io/docs/user-guide/identifiers#names
	Name string `json:"name"`
}

// SecretKeySelector references data from a Secret resource by a specific key.
type SecretKeySelector struct {
	// Name is the name of the secret.
	Name string `json:"name"`

	// Key is the key in the secret that specifies the requested data.
	Key string `json:"key"`
}

// ValueOrSecretKeySelector describes a value from either a SecretKeySelector
// or value directly in this object.
type ValueOrSecretKeySelector struct {
	// +optional

	// From is specified to reference a value from a Secret resource.
	//
	// Please note this field is mutually exclusive with the Value field.
	From *SecretKeySelector `json:"from,omitempty"`

	// +optional

	// Value is used to directly specify a value.
	//
	// Please note this field is mutually exclusive with the From field.
	Value *string `json:"value,omitempty"`
}

// KeyValuePair is useful when wanting to realize a map as a list of key/value
// pairs.
type KeyValuePair struct {
	// Key is the key part of the key/value pair.
	Key string `json:"key"`

	// +optional

	// Value is the optional value part of the key/value pair.
	Value string `json:"value,omitempty"`
}

// KeyValueOrSecretKeySelectorPair is useful when wanting to realize a map as a
// list of key/value pairs where each value could also reference data stored in
// a Secret resource.
type KeyValueOrSecretKeySelectorPair struct {
	// Key is the key part of the key/value pair.
	Key string `json:"key"`

	// +optional

	// Value is the optional value part of the key/value pair.
	Value ValueOrSecretKeySelector `json:"value,omitempty"`
}

// NameValuePair is useful when wanting to realize a map as a list of name/value
// pairs.
type NameValuePair struct {
	// Name is the name part of the name/value pair.
	Name string `json:"name"`

	// +optional

	// Value is the optional value part of the name/value pair.
	Value string `json:"value,omitempty"`
}

// PasswordSecretKeySelector references the password value from a Secret resource.
type PasswordSecretKeySelector struct {
	// Name is the name of the secret.
	Name string `json:"name"`

	// +optional
	// +kubebuilder:default=password

	// Key is the key in the secret that specifies the requested data.
	Key string `json:"key,omitempty"`
}

// ObjectMeta is metadata that all persisted resources must have, which includes
// all objects users must create. This is a copy of customizable fields from
// metav1.ObjectMeta.

// ObjectMeta is embedded in `VirtualMachineReplicaSet.Template`, which is not a
// top-level Kubernetes object. By default, controller-gen handles certain known
// types (e.g., metav1.ObjectMeta) differently by not including their properties.
// See: https://github.com/kubernetes-sigs/controller-tools/issues/385.

// With https://github.com/kubernetes-sigs/controller-tools/pull/557, this is
// somewhat fixed so that an embedded metav1.ObjectMeta is expanded while dropping
// some problematic fields. However, in our case, we are only interested
// in Labels and Annotations, so we still end up with extraneous fields (e.g.,
// Namespace). For that purpose, we introduce our copy of metav1.ObjectMeta that
// only contains the fields that we are interested in.

type ObjectMeta struct {
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Validate validates the labels and annotations in ObjectMeta.
func (metadata *ObjectMeta) Validate(parent *field.Path) field.ErrorList {
	allErrs := metav1validation.ValidateLabels(
		metadata.Labels,
		parent.Child("labels"),
	)
	allErrs = append(allErrs, apivalidation.ValidateAnnotations(
		metadata.Annotations,
		parent.Child("annotations"),
	)...)
	return allErrs
}
