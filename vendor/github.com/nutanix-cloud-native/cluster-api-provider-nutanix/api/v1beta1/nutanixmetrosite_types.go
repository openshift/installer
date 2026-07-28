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
	// Kind represents the Kind of NutanixMetroSite
	NutanixMetroSiteKind = "NutanixMetroSite"

	// NutanixMetroSiteFinalizer is the finalizer used by the NutanixMetroSite controller to block
	// deletion of the NutanixMetroSite object if there are references to this object by other resources.
	NutanixMetroSiteFinalizer = "infrastructure.cluster.x-k8s.io/nutanixmetrosite"
)

// NutanixMetroSiteSpec defines the desired state of NutanixMetroSite
type NutanixMetroSiteSpec struct {
	// preferredFailureDomain holds reference to the preferred NutanixFailureDomain object of the NutanixMetroSite.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="preferredFailureDomain is immutable once set"
	// +kubebuilder:validation:XValidation:rule=`self.name != ""`,message="preferredFailureDomain.name must not be empty"
	// +kubebuilder:validation:Required
	PreferredFailureDomain corev1.LocalObjectReference `json:"preferredFailureDomain"`

	// metroRef holds reference to the NutanixMetro object this site belongs to.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="metroRef is immutable once set"
	// +kubebuilder:validation:XValidation:rule=`self.name != ""`,message="metroRef.name must not be empty"
	// +kubebuilder:validation:Required
	MetroRef corev1.LocalObjectReference `json:"metroRef"`

	// groupNameLabel optionally provides name for a group of worker nodes normally anchored to a topology segment.
	// +optional
	GroupNameLabel *string `json:"groupNameLabel,omitempty"`
}

// NutanixMetroSiteStatus defines the state of the NutanixMetroSite resource.
type NutanixMetroSiteStatus struct {
	// conditions represent the latest states of the metro site,
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nutanixmetrosites,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:labels=clusterctl.cluster.x-k8s.io/move=
// +kubebuilder:printcolumn:name="PreferredFailureDomain",type="string",JSONPath=".spec.preferredFailureDomain",description="Reference of the preferred NutanixFailureDomain object"
// +kubebuilder:printcolumn:name="Metro",type="string",JSONPath=".spec.metroRef",description="Reference of NutanixMetro object this site belongs to"

// NutanixMetroSite is the Schema for the NutanixMetroSite API.
type NutanixMetroSite struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NutanixMetroSiteSpec   `json:"spec,omitempty"`
	Status NutanixMetroSiteStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (m *NutanixMetroSite) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (m *NutanixMetroSite) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// NutanixMetroSiteList contains a list of NutanixMetroSite resources
type NutanixMetroSiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NutanixMetroSite `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NutanixMetroSite{}, &NutanixMetroSiteList{})
}
