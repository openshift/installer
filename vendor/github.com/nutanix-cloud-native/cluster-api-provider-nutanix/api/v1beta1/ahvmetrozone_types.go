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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck // suppress complaining on Deprecated package
)

const (
	// Kind represents the Kind of AHVMetroZone
	AHVMetroZoneKind = "AHVMetroZone"

	// AHVMetroZoneFinalizer is the finalizer used by the AHVMetroZone controller to block
	// deletion of the AHVMetroZone object if there are references to this object by other resources.
	AHVMetroZoneFinalizer = "infrastructure.cluster.x-k8s.io/ahvmetrozone"
)

// AHVMetroZoneSpec defines the desired state of AHVMetroZone
type AHVMetroZoneSpec struct {
	// zones defines the Prism Element zones of the AHV Metro Domain.
	// +kubebuilder:validation:MinItems=2
	// +listType=map
	// +listMapKey=name
	Zones []MetroZone `json:"zones"`

	// placement defines the VM provisioning placement strategy
	Placement VMPlacement `json:"placement"`
}

// MetroZone defines a Prism Element zone of the AHV Metro Domain.
type MetroZone struct {
	// name is the unique name of the metro zone.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	Name string `json:"name"`

	// prismElement is the identifier of the Prism Element.
	// +kubebuilder:validation:Required
	PrismElement NutanixResourceIdentifier `json:"prismElement"`

	// subnets holds a list of identifiers (one or more) of the subnets.
	// The subnets should already exist in PC and there is no duplicate items configured.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=32
	Subnets []NutanixResourceIdentifier `json:"subnets"`
}

// VMPlacementStrategy is an enumeration of VM placement strategies.
type VMPlacementStrategy string

const (
	// PreferredStrategy is the VM placement strategy with the specified preferred zone.
	PreferredStrategy VMPlacementStrategy = "Preferred"

	// RandomStrategy is the VM placement strategy with random choice of zone.
	RandomStrategy VMPlacementStrategy = "Random"
)

// VMPlacement defines the placement strategy when provisioning a VM.
// +kubebuilder:validation:XValidation:rule=`self.strategy != "Preferred" || (has(self.preferredZone) && self.preferredZone != "")`,message="preferredZone is required for Preferred strategy."
type VMPlacement struct {
	// strategy defines the VM placement strategy of Preferred or Random,
	// with Preferred as the default. When the strategy is Preferred, the
	// preferredZone must be specified.
	// +kubebuilder:validation:Required
	// +kubebuilder:default="Preferred"
	// +kubebuilder:validation:Enum:=Preferred;Random
	Strategy VMPlacementStrategy `json:"strategy"`

	// preferredZone specifies the VM should be provisioned to the specified
	// preferred zone. When the preferred zone is not available, or being
	// evacuated (e.g. in maintainence), the VM will be provisioned to the
	// remaining zone. This field is required to set when the VM placement
	// strategy is set to Preferred.
	// +optional
	PreferredZone *string `json:"preferredZone,omitempty"`
}

// AHVMetroZoneStatus defines the observed state of AHVMetroZone resource.
type AHVMetroZoneStatus struct {
	// conditions represent the latest states of the AHVMetroZone.
	// +optional
	Conditions []capiv1beta1.Condition `json:"conditions,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in AHVMetroZone's status with the v1beta2 version.
	// +optional
	V1Beta2 *AHVMetroZoneV1Beta2Status `json:"v1beta2,omitempty"`
}

// AHVMetroZoneV1Beta2Status groups all the fields that will be added or modified in AHVMetroZoneStatus with the v1beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type AHVMetroZoneV1Beta2Status struct {
	// conditions represents the observations of a AHVMetroZone's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=ahvmetrozones,scope=Namespaced,categories=cluster-api
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:metadata:labels=clusterctl.cluster.x-k8s.io/move=

// AHVMetroZone is the Schema for the ahvmetrozones API.
type AHVMetroZone struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AHVMetroZoneSpec   `json:"spec,omitempty"`
	Status AHVMetroZoneStatus `json:"status,omitempty"`
}

// GetConditions returns the set of conditions for this object.
func (z *AHVMetroZone) GetConditions() capiv1beta1.Conditions {
	return z.Status.Conditions
}

// SetConditions sets the conditions on this object.
func (z *AHVMetroZone) SetConditions(conditions capiv1beta1.Conditions) {
	z.Status.Conditions = conditions
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (z *AHVMetroZone) GetV1Beta2Conditions() []metav1.Condition {
	if z.Status.V1Beta2 == nil {
		return nil
	}
	return z.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets the v1beta2 conditions on this object.
func (z *AHVMetroZone) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if z.Status.V1Beta2 == nil {
		z.Status.V1Beta2 = &AHVMetroZoneV1Beta2Status{}
	}
	z.Status.V1Beta2.Conditions = conditions
}

// +kubebuilder:object:root=true

// AHVMetroZoneList contains a list of AHVMetroZone resources
type AHVMetroZoneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AHVMetroZone `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AHVMetroZone{}, &AHVMetroZoneList{})
}
