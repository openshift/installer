/*
Copyright 2026 Nutanix

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Kind represents the Kind of NutanixMetro
	NutanixMetroKind = "NutanixMetro"

	// NutanixMetroFinalizer is the finalizer used by the NutanixMetro controller to block
	// deletion of the NutanixMetro object if there are references to this object by other resources.
	NutanixMetroFinalizer = "infrastructure.cluster.x-k8s.io/nutanixmetro"
)

// NutanixMetroSpec defines the desired state of NutanixMetro
type NutanixMetroSpec struct {
	// failureDomains holds references to the two NutanixFailureDomain objects of the NutanixMetro.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="failureDomains is immutable once set"
	// +kubebuilder:validation:XValidation:rule=`self.all(x, x.name != "")`,message="each failureDomain reference name must not be empty"
	// +kubebuilder:validation:MinItems=2
	// +kubebuilder:validation:MaxItems=2
	// +listType=map
	// +listMapKey=name
	FailureDomains []corev1.LocalObjectReference `json:"failureDomains"`
}

// NutanixMetroStatus defines the state of the NutanixMetro resource.
type NutanixMetroStatus struct {
	// conditions represent the latest states of the metro
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nutanixmetros,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:annotations="clusterctl.cluster.x-k8s.io/skip-crd-name-preflight-check=true"
// +kubebuilder:metadata:labels=clusterctl.cluster.x-k8s.io/move=
// +kubebuilder:printcolumn:name="FailureDomains",type="string",JSONPath=".spec.failureDomains",description="References of NutanixFailureDomain objects"

// NutanixMetro is the Schema for the NutanixMetro API.
type NutanixMetro struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NutanixMetroSpec   `json:"spec,omitempty"`
	Status NutanixMetroStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (m *NutanixMetro) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (m *NutanixMetro) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// NutanixMetroList contains a list of NutanixMetro resources
type NutanixMetroList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NutanixMetro `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NutanixMetro{}, &NutanixMetroList{})
}
