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

// CloudCredentialsReference is a reference to a secret containing OpenStack credentials.
type CloudCredentialsReference struct {
	// SecretName is the name of a secret in the same namespace as the resource being provisioned.
	// The secret must contain a key named `clouds.yaml` which contains an OpenStack clouds.yaml file.
	// The secret may optionally contain a key named `cacert` containing a PEM-encoded CA certificate.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=253
	SecretName string `json:"secretName"`

	// CloudName specifies the name of the entry in the clouds.yaml file to use.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength:=1
	// +kubebuilder:validation:MaxLength:=256
	CloudName string `json:"cloudName"`
}

// CloudCredentialsRefProvider is an interface for obtaining OpenStack credentials from an API object
// +kubebuilder:object:generate:=false
type CloudCredentialsRefProvider interface {
	GetCloudCredentialsRef() (*string, *CloudCredentialsReference)
}
