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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// GKEConfigSpec defines the desired state of GCP GKE Bootstrap Configuration.
type GKEConfigSpec struct{}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=gkeconfigs,scope=Namespaced,categories=cluster-api,shortName=gkec
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Bootstrap configuration is ready"
// +kubebuilder:printcolumn:name="DataSecretName",type="string",JSONPath=".status.dataSecretName",description="Name of Secret containing bootstrap data"

// GKEConfig is the schema for the GCP GKE Bootstrap Configuration.
// this is a placeholder used for compliance with the CAPI contract.
type GKEConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GKEConfigSpec   `json:"spec,omitempty"`
	Status GKEConfigStatus `json:"status,omitempty"`
}

// GKEConfigStatus defines the observed state of the GCP GKE Bootstrap Configuration.
type GKEConfigStatus struct {
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

	// Conditions defines current service state of the GKEConfig.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true

// GKEConfigList contains a list of GCP GKE Bootstrap Configuration.
type GKEConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GKEConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GKEConfig{}, &GKEConfigList{})
}
