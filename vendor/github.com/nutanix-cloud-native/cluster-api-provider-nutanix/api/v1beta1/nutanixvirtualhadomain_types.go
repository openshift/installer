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
	// NutanixVirtualHADomainKind represents the Kind of NutanixVirtualHADomain
	NutanixVirtualHADomainKind = "NutanixVirtualHADomain"

	// NutanixVirtualHADomainFinalizer is the finalizer used by the NutanixVirtualHADomain controller to block
	// deletion of the NutanixVirtualHADomain object if there are references to this object by other resources.
	NutanixVirtualHADomainFinalizer = "infrastructure.cluster.x-k8s.io/nutanixvirtualhadomain"
)

// NutanixVirtualHADomainSpec defines the desired state of NutanixVirtualHADomain.
type NutanixVirtualHADomainSpec struct {
	// metroRef is a reference to the NutanixMetro object that this virtual HA domain belongs to.
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="metroRef is immutable once set"
	// +kubebuilder:validation:XValidation:rule=`self.name != ""`,message="metroRef.name must not be empty"
	// +kubebuilder:validation:Required
	MetroRef corev1.LocalObjectReference `json:"metroRef"`

	// protectionPolicy identifies the protection policy PC resource applied to this virtual HA domain.
	// +optional
	ProtectionPolicy *NutanixResourceIdentifier `json:"protectionPolicy"`

	// movementGroups defines the named groups of entities that move together within this virtual HA domain.
	// +listType=map
	// +listMapKey=name
	// +optional
	MovementGroups []NutanixMovementGroup `json:"movementGroups"`
}

// NutanixMovementGroup defines a named group of entities that are moved together as part of a
// virtual HA domain failover or migration. It maps each category to the recovery plan
// that protects the entities associated with that category.
type NutanixMovementGroup struct {
	// name is the name of the movement group
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=`[a-z0-9]([-a-z0-9]*[a-z0-9])?`
	Name string `json:"name"`

	// categoryRecoveryPlans is the list of category-to-recovery-plan mappings whose
	// member entities belong to this movement group.
	// +kubebuilder:validation:MinItems=2
	// +kubebuilder:validation:MaxItems=2
	// +listType=atomic
	CategoryRecoveryPlans []NutanixCategoryRecoveryPlan `json:"categoryRecoveryPlans"`
}

// NutanixCategoryRecoveryPlan maps a category to the recovery plan that protects the
// entities belonging to that category within a movement group, on a given Prism Element.
type NutanixCategoryRecoveryPlan struct {
	// category is the category (key/value) whose member entities are protected.
	// +kubebuilder:validation:Required
	Category NutanixCategoryIdentifier `json:"category"`

	// recoveryPlan is the recovery plan associated with the category.
	// +kubebuilder:validation:Required
	RecoveryPlan NutanixResourceIdentifier `json:"recoveryPlan"`

	// failureDomainRef is a reference to the NutanixFailureDomain object that identifies the
	// Prism Element cluster for this category-to-recovery-plan mapping.
	// +kubebuilder:validation:XValidation:rule=`self.name != ""`,message="failureDomainRef.name must not be empty"
	// +kubebuilder:validation:Required
	FailureDomainRef corev1.LocalObjectReference `json:"failureDomainRef"`
}

// NutanixVirtualHADomainStatus defines the observed state of NutanixVirtualHADomain.
type NutanixVirtualHADomainStatus struct {
	// conditions represent the latest states of the virtual HA domain.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// ready is set to true when the vHA domain PC resources (categories, protection
	// policy, and recovery plan) are valid and ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=nutanixvirtualhadomains,shortName=nvha,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:labels=clusterctl.cluster.x-k8s.io/move=
// +kubebuilder:printcolumn:name="Metro",type="string",JSONPath=".spec.metroRef.name",description="Reference of NutanixMetro object"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready",description="the vHA domain PC resources are ready or not"

// NutanixVirtualHADomain is the Schema for the NutanixVirtualHADomains API.
type NutanixVirtualHADomain struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NutanixVirtualHADomainSpec   `json:"spec,omitempty"`
	Status NutanixVirtualHADomainStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (d *NutanixVirtualHADomain) GetConditions() []metav1.Condition {
	return d.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (d *NutanixVirtualHADomain) SetConditions(conditions []metav1.Condition) {
	d.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// NutanixVirtualHADomainList contains a list of NutanixVirtualHADomain resources
type NutanixVirtualHADomainList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NutanixVirtualHADomain `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NutanixVirtualHADomain{}, &NutanixVirtualHADomainList{})
}
