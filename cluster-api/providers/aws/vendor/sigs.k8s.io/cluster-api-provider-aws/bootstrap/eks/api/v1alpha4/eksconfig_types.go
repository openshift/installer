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

package v1alpha4

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
)

// EKSConfigSpec defines the desired state of EKSConfig
type EKSConfigSpec struct {
	// Passes the kubelet args into the EKS bootstrap script
	// +optional
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`
}

// EKSConfigStatus defines the observed state of EKSConfig
type EKSConfigStatus struct {
	// Ready indicates the BootstrapData secret is ready to be consumed
	Ready bool `json:"ready,omitempty"`

	// DataSecretName is the name of the secret that stores the bootstrap data script.
	// +optional
	DataSecretName *string `json:"dataSecretName,omitempty"`

	// FailureReason will be set on non-retryable errors
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage will be set on non-retryable errors
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// ObservedGeneration is the latest generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions defines current service state of the EKSConfig.
	// +optional
	Conditions clusterv1alpha4.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=eksconfigs,scope=Namespaced,categories=cluster-api,shortName=eksc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Bootstrap configuration is ready"
// +kubebuilder:printcolumn:name="DataSecretName",type="string",JSONPath=".status.dataSecretName",description="Name of Secret containing bootstrap data"

// EKSConfig is the Schema for the eksconfigs API
type EKSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSConfigSpec   `json:"spec,omitempty"`
	Status EKSConfigStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the EKSConfig resource.
func (r *EKSConfig) GetConditions() clusterv1alpha4.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the EKSConfig to the predescribed clusterv1alpha4.Conditions.
func (r *EKSConfig) SetConditions(conditions clusterv1alpha4.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// EKSConfigList contains a list of EKSConfig.
type EKSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSConfig{}, &EKSConfigList{})
}
