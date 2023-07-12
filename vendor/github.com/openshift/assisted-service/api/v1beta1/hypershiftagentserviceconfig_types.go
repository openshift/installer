/*
Copyright 2021.

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
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HypershiftAgentServiceConfigSpec defines the desired state of HypershiftAgentServiceConfig.
type HypershiftAgentServiceConfigSpec struct {
	AgentServiceConfigSpec `json:",inline"`

	// KubeconfigSecretRef is a reference to the secret containing the kubeconfig for the destination Hypershift instance.
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Hypershift kubeconfig reference"
	KubeconfigSecretRef corev1.LocalObjectReference `json:"kubeconfigSecretRef"`
}

type HypershiftAgentServiceConfigStatus struct {
	Conditions []conditionsv1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Namespaced

// HypershiftAgentServiceConfig represents an Assisted Service deployment over zero-worker
// hypershift cluster. Each resource represents a deployment scheme over hosted cluster
// that runs in that namespace.
// +kubebuilder:resource:shortName=hasc
// +operator-sdk:csv:customresourcedefinitions:displayName="Hypershift Service Config"
// +operator-sdk:csv:customresourcedefinitions:order=1
type HypershiftAgentServiceConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HypershiftAgentServiceConfigSpec   `json:"spec,omitempty"`
	Status HypershiftAgentServiceConfigStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// HypershiftAgentServiceConfigList contains a list of HypershiftAgentServiceConfigs
type HypershiftAgentServiceConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HypershiftAgentServiceConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HypershiftAgentServiceConfig{}, &HypershiftAgentServiceConfigList{})
}
