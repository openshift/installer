/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha5

// OpenStackIdentityReference is a reference to an infrastructure
// provider identity to be used to provision cluster resources.
type OpenStackIdentityReference struct {
	// Kind of the identity. Must be supported by the infrastructure
	// provider and may be either cluster or namespace-scoped.
	// +kubebuilder:validation:MinLength=1
	Kind string `json:"kind"`

	// Name of the infrastructure identity to be used.
	// Must be either a cluster-scoped resource, or namespaced-scoped
	// resource the same namespace as the resource(s) being provisioned.
	Name string `json:"name"`
}
