/*
Copyright 2025 Nutanix

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
	capiv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// NutanixFailureDomainKind represents the Kind of NutanixFailureDomain
	NutanixFailureDomainKind = "NutanixFailureDomain"

	// NutanixFailureDomainFinalizer is the finalizer used by the NutanixFailureDomain controller to block
	// deletion of the NutanixFailureDomain object if there are references to this object by other resources.
	NutanixFailureDomainFinalizer = "infrastructure.cluster.x-k8s.io/nutanixfailuredomain"
)

// NutanixFailureDomainSpec defines the desired state of NutanixFailureDomain.
type NutanixFailureDomainSpec struct {
	// prismElementCluster is to identify the Prism Element cluster in the Prism Central for the failure domain.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="prismElementCluster is immutable once set"
	// +kubebuilder:validation:Required
	PrismElementCluster NutanixResourceIdentifier `json:"prismElementCluster"`

	// subnets holds a list of identifiers (one or more) of the PE cluster's network subnets
	// for the Machine's VM to connect to.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="subnets is immutable once set"
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Subnets []NutanixResourceIdentifier `json:"subnets"`
}

// NutanixFailureDomainStatus defines the observed state of NutanixFailureDomain resource.
type NutanixFailureDomainStatus struct {
	// conditions represent the latest states of the failure domain.
	// +optional
	Conditions []capiv1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nutanixfailuredomains,shortName=nfd,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:labels=clusterctl.cluster.x-k8s.io/move=

// NutanixFailureDomain is the Schema for the NutanixFailureDomain API.
type NutanixFailureDomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NutanixFailureDomainSpec   `json:"spec,omitempty"`
	Status NutanixFailureDomainStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (nfd *NutanixFailureDomain) GetConditions() capiv1.Conditions {
	return nfd.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (nfd *NutanixFailureDomain) SetConditions(conditions capiv1.Conditions) {
	nfd.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// NutanixFailureDomainList contains a list of NutanixFailureDomain
type NutanixFailureDomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NutanixFailureDomain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NutanixFailureDomain{}, &NutanixFailureDomainList{})
}
