// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vmopv1common "github.com/vmware-tanzu/vm-operator/api/v1alpha6/common"
)

const (
	// VirtualMachineReplicaSetReplicaFailure is added in a replica set when one of
	// its VMs fails to be created due to insufficient quota, limit ranges,
	// security policy, host selection, or deleted due to the host being down or
	// finalizers are failing.
	VirtualMachineReplicaSetReplicaFailure = "ReplicaFailure"
)

const (
	// VirtualMachinesCreatedCondition documents that the virtual machines controlled by
	// the VirtualMachineReplicaSet are created. When this condition is false, it indicates that
	// there was an error when creating the VirtualMachine object.
	VirtualMachinesCreatedCondition = "VirtualMachinesCreated"

	// VirtualMachinesReadyCondition reports an aggregate of current status of
	// the virtual machines controlled by the VirtualMachineReplicaSet.
	VirtualMachinesReadyCondition = "VirtualMachinesReady"

	// VirtualMachineCreationFailedReason documents a VirtualMachineReplicaSet failing to
	// generate a VirtualMachine object.
	VirtualMachineCreationFailedReason = "VirtualMachineCreationFailed"

	// ResizedCondition documents a VirtualMachineReplicaSet is resizing the set of controlled VirtualMachines.
	ResizedCondition = "Resized"

	// ScalingUpReason documents a VirtualMachineReplicaSet is increasing the number of replicas.
	ScalingUpReason = "ScalingUp"

	// ScalingDownReason documents a VirtualMachineReplicaSet is decreasing the number of replicas.
	ScalingDownReason = "ScalingDown"
)

const (
	// VirtualMachineReplicaSetNameLabel is the key of the label applied on all the
	// replicas VirtualMachine objects that it owns.  The value of this label is the
	// name of the VirtualMachineReplicaSet.
	VirtualMachineReplicaSetNameLabel = "vmoperator.vmware.com/replicaset-name"
)

// VirtualMachineTemplateSpec describes the data needed to create a VirtualMachine
// from a template.
type VirtualMachineTemplateSpec struct {
	// +optional
	//
	// ObjectMeta contains the desired Labels and Annotations that must be
	// applied to each replica virtual machine.
	vmopv1common.ObjectMeta `json:"metadata,omitempty"`

	// +optional
	//
	// Specification of the desired behavior of each replica virtual machine.
	Spec VirtualMachineSpec `json:"spec,omitempty"`
}

// VirtualMachineReplicaSetSpec is the specification of a VirtualMachineReplicaSet.
type VirtualMachineReplicaSetSpec struct {
	// +optional
	// +kubebuilder:default=1
	//
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	Replicas *int32 `json:"replicas,omitempty"`

	// +optional
	// +kubebuilder:validation:Enum=Random
	//
	// DeletePolicy defines the policy used to identify nodes to delete when downscaling.
	// Only supported deletion policy is "Random".
	DeletePolicy string `json:"deletePolicy,omitempty"`

	// +optional
	//
	// Selector is a label to query over virtual machines that should match the
	// replica count. A virtual machine's label keys and values must match in order
	// to be controlled by this VirtualMachineReplicaSet.
	//
	// It must match the VirtualMachine template's labels.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	Selector *metav1.LabelSelector `json:"selector"`

	// +optional
	//
	// Template is the object that describes the virtual machine that will be
	// created if insufficient replicas are detected.
	Template VirtualMachineTemplateSpec `json:"template,omitempty"`
}

// VirtualMachineReplicaSetStatus represents the observed state of a
// VirtualMachineReplicaSet resource.
type VirtualMachineReplicaSetStatus struct {
	// +optional
	//
	// Replicas is the most recently observed number of replicas.
	Replicas int32 `json:"replicas"`

	// +optional
	//
	// FullyLabeledReplicas is the number of replicas that have labels matching the
	// labels of the virtual machine template of the VirtualMachineReplicaSet.
	FullyLabeledReplicas int32 `json:"fullyLabeledReplicas,omitempty"`

	// +optional
	//
	// ReadyReplicas is the number of ready replicas for this VirtualMachineReplicaSet. A
	// virtual machine is considered ready when it's "Ready" condition is marked as
	// true.
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// +optional
	//
	// ObservedGeneration reflects the generation of the most recently observed
	// VirtualMachineReplicaSet.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// +optional
	//
	// Conditions represents the latest available observations of a
	// VirtualMachineReplicaSet's current state.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

func (rs *VirtualMachineReplicaSet) GetConditions() []metav1.Condition {
	return rs.Status.Conditions
}

func (rs *VirtualMachineReplicaSet) SetConditions(conditions []metav1.Condition) {
	rs.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmrs;vmreplicaset
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="Total number of non-terminated virtual machines targeted by this VirtualMachineReplicaSet"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas",description="Total number of ready virtual machines targeted by this VirtualMachineReplicaSet"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VirtualMachineReplicaSet"

// VirtualMachineReplicaSet is the schema for the virtualmachinereplicasets API.
type VirtualMachineReplicaSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineReplicaSetSpec   `json:"spec,omitempty"`
	Status VirtualMachineReplicaSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineReplicaSetList contains a list of VirtualMachineReplicaSet.
type VirtualMachineReplicaSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineReplicaSet `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VirtualMachineReplicaSet{}, &VirtualMachineReplicaSetList{})
}
