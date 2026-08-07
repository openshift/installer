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

// KeyPairResourceSpec contains the desired state of the resource.
type KeyPairResourceSpec struct {
	// name will be the name of the created resource. If not specified, the
	// name of the ORC object will be used.
	// +optional
	Name *OpenStackName `json:"name,omitempty"`

	// type specifies the type of the Keypair. Allowed values are ssh or x509.
	// If not specified, defaults to ssh.
	// +kubebuilder:validation:Enum=ssh;x509
	// +optional
	Type *string `json:"type,omitempty"`

	// publicKey is the public key to import.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=16384
	// +required
	PublicKey string `json:"publicKey,omitempty"`
}

// KeyPairFilter defines an existing resource by its properties
// +kubebuilder:validation:MinProperties:=1
type KeyPairFilter struct {
	// name of the existing Keypair
	// +optional
	Name *OpenStackName `json:"name,omitempty"`
}

// KeyPairResourceStatus represents the observed state of the resource.
type KeyPairResourceStatus struct {
	// name is a Human-readable name for the resource. Might not be unique.
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Name string `json:"name,omitempty"`

	// fingerprint is the fingerprint of the public key
	// +kubebuilder:validation:MaxLength=1024
	// +optional
	Fingerprint string `json:"fingerprint,omitempty"`

	// publicKey is the public key of the Keypair
	// +kubebuilder:validation:MaxLength=16384
	// +optional
	PublicKey string `json:"publicKey,omitempty"`

	// type is the type of the Keypair (ssh or x509)
	// +kubebuilder:validation:MaxLength=64
	// +optional
	Type string `json:"type,omitempty"`
}
