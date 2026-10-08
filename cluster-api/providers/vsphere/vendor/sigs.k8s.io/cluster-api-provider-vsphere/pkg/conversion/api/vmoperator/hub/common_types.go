/*
Copyright 2026 The Kubernetes Authors.

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

// +kubebuilder:object:generate=true

package hub

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

// KeyValuePair is useful when wanting to realize a map as a list of key/value
// pairs.
type KeyValuePair struct {
	// Key is the key part of the key/value pair.
	Key string `json:"key"`

	// +optional

	// Value is the optional value part of the key/value pair.
	Value string `json:"value,omitempty"`
}
