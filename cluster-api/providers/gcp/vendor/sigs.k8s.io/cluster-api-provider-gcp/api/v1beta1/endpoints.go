/*
Copyright 2025 The Kubernetes Authors.

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

// ServiceEndpoints contains all the gcp service endpoints that the user may override. Each field corresponds to
// a service where the expected value is the url that is used to override the default API endpoint.
type ServiceEndpoints struct {
	// ComputeServiceEndpoint is the custom endpoint url for the Compute Service
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:validation:Pattern=`^https://`
	// +optional
	ComputeServiceEndpoint string `json:"compute,omitempty"`

	// ContainerServiceEndpoint is the custom endpoint url for the Container Service
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:validation:Pattern=`^https://`
	// +optional
	ContainerServiceEndpoint string `json:"container,omitempty"`

	// IAMServiceEndpoint is the custom endpoint url for the IAM Service
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:validation:Pattern=`^https://`
	// +optional
	IAMServiceEndpoint string `json:"iam,omitempty"`

	// ResourceManagerServiceEndpoint is the custom endpoint url for the Resource Manager Service
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=uri
	// +kubebuilder:validation:Pattern=`^https://`
	// +optional
	ResourceManagerServiceEndpoint string `json:"resourceManager,omitempty"`
}
