/*
Copyright 2018 The Kubernetes Authors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
)

// Note: These types have been inlined from CAPI v1alpha4 as this package is not available/exported anymore.

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// users must create. This is a copy of customizable fields from metav1.ObjectMeta.
//
// ObjectMeta is embedded in `Machine.Spec`, `MachineDeployment.Template` and `MachineSet.Template`,
// which are not top-level Kubernetes objects. Given that metav1.ObjectMeta has lots of special cases
// and read-only fields which end up in the generated CRD validation, having it as a subset simplifies
// the API and some issues that can impact user experience.
//
// During the [upgrade to controller-tools@v2](https://github.com/kubernetes-sigs/cluster-api/pull/1054)
// for v1alpha2, we noticed a failure would occur running Cluster API test suite against the new CRDs,
// specifically `spec.metadata.creationTimestamp in body must be of type string: "null"`.
// The investigation showed that `controller-tools@v2` behaves differently than its previous version
// when handling types from [metav1](k8s.io/apimachinery/pkg/apis/meta/v1) package.
//
// In more details, we found that embedded (non-top level) types that embedded `metav1.ObjectMeta`
// had validation properties, including for `creationTimestamp` (metav1.Time).
// The `metav1.Time` type specifies a custom json marshaller that, when IsZero() is true, returns `null`
// which breaks validation because the field isn't marked as nullable.
//
// In future versions, controller-tools@v2 might allow overriding the type and validation for embedded
// types. When that happens, this hack should be revisited.
type ObjectMeta struct {
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
}

// ANCHOR: ConditionSeverity

// ConditionSeverity expresses the severity of a Condition Type failing.
type ConditionSeverity string

const (
	// ConditionSeverityError specifies that a condition with `Status=False` is an error.
	ConditionSeverityError ConditionSeverity = "Error"

	// ConditionSeverityWarning specifies that a condition with `Status=False` is a warning.
	ConditionSeverityWarning ConditionSeverity = "Warning"

	// ConditionSeverityInfo specifies that a condition with `Status=False` is informative.
	ConditionSeverityInfo ConditionSeverity = "Info"

	// ConditionSeverityNone should apply only to conditions with `Status=True`.
	ConditionSeverityNone ConditionSeverity = ""
)

// ANCHOR_END: ConditionSeverity

// ANCHOR: ConditionType

// ConditionType is a valid value for Condition.Type.
type ConditionType string

// ANCHOR_END: ConditionType

// ANCHOR: Condition

// Condition defines an observation of a Cluster API resource operational state.
type Condition struct {
	// Type of condition in CamelCase or in foo.example.com/CamelCase.
	// Many .condition.type values are consistent across resources like Available, but because arbitrary conditions
	// can be useful (see .node.status.conditions), the ability to deconflict is important.
	// +required
	Type ConditionType `json:"type"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	Status corev1.ConditionStatus `json:"status"`

	// Severity provides an explicit classification of Reason code, so the users or machines can immediately
	// understand the current situation and act accordingly.
	// The Severity field MUST be set only when Status=False.
	// +optional
	Severity ConditionSeverity `json:"severity,omitempty"`

	// Last time the condition transitioned from one status to another.
	// This should be when the underlying condition changed. If that is not known, then using the time when
	// the API field changed is acceptable.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`

	// The reason for the condition's last transition in CamelCase.
	// The specific API may choose whether or not this field is considered a guaranteed API.
	// This field may not be empty.
	// +optional
	Reason string `json:"reason,omitempty"`

	// A human readable message indicating details about the transition.
	// This field may be empty.
	// +optional
	Message string `json:"message,omitempty"`
}

// ANCHOR_END: Condition

// ANCHOR: Conditions

// Conditions provide observations of the operational state of a Cluster API resource.
type Conditions []Condition

// ANCHOR_END: Conditions

// FailureDomains is a slice of FailureDomains.
type FailureDomains map[string]FailureDomainSpec

// FailureDomainSpec is the Schema for Cluster API failure domains.
// It allows controllers to understand how many failure domains a cluster can optionally span across.
type FailureDomainSpec struct {
	// ControlPlane determines if this failure domain is suitable for use by control plane machines.
	// +optional
	ControlPlane bool `json:"controlPlane"`

	// Attributes is a free form map of attributes an infrastructure provider might use or require.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
}

// MachineAddressType describes a valid MachineAddress type.
type MachineAddressType string

// Define the MachineAddressType constants.
const (
	MachineHostName    MachineAddressType = "Hostname"
	MachineExternalIP  MachineAddressType = "ExternalIP"
	MachineInternalIP  MachineAddressType = "InternalIP"
	MachineExternalDNS MachineAddressType = "ExternalDNS"
	MachineInternalDNS MachineAddressType = "InternalDNS"
)

// MachineAddress contains information for the node's address.
type MachineAddress struct {
	// Machine address type, one of Hostname, ExternalIP or InternalIP.
	Type MachineAddressType `json:"type"`

	// The machine address.
	Address string `json:"address"`
}

// MachineAddresses is a slice of MachineAddress items to be used by infrastructure providers.
type MachineAddresses []MachineAddress

// Note: The following conversion functions are somehow required by the generated conversion code, but not used (tested via our conversion fuzz tests).

func Convert_v1alpha4_Condition_To_v1_Condition(in *Condition, out *metav1.Condition, s apiconversion.Scope) error {
	// in.Severity does not exists in v1.Condition.
	return autoConvert_v1alpha4_Condition_To_v1_Condition(in, out, s)
}

func Convert_v1_Condition_To_v1alpha4_Condition(in *metav1.Condition, out *Condition, s apiconversion.Scope) error {
	// in.ObservedGeneration does not exists in v1alpha3.Condition.
	return autoConvert_v1_Condition_To_v1alpha4_Condition(in, out, s)
}

func Convert_v1_ObjectMeta_To_v1alpha4_ObjectMeta(in *metav1.ObjectMeta, out *ObjectMeta, s apiconversion.Scope) error {
	// a few fields don't exist in v1alpha3.Condition.
	return autoConvert_v1_ObjectMeta_To_v1alpha4_ObjectMeta(in, out, s)
}
