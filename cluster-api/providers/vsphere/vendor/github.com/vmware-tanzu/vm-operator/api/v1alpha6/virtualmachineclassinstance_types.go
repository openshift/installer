// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// VMClassInstanceActiveLabelKey represents if a
	// VirtualMachineClassInstance can be used to deploy a new virtual
	// machine, or resize an existing one. The presence of this label
	// on a VirtualMachineClassInstance marks the instance as active.
	// The value of this key does not matter.
	VMClassInstanceActiveLabelKey = GroupName + "/active"
)

// VirtualMachineClassInstanceSpec defines the desired state of VirtualMachineClassInstance.
// It is a composite of VirtualMachineClassSpec.
type VirtualMachineClassInstanceSpec struct {
	VirtualMachineClassSpec `json:",inline"`
}

// VirtualMachineClassInstanceStatus defines the observed state of VirtualMachineClassInstance.
type VirtualMachineClassInstanceStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmclassinstance
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.hardware.cpus"
// +kubebuilder:printcolumn:name="Memory",type="string",JSONPath=".spec.hardware.memory"
// +kubebuilder:printcolumn:name="Active",type="string",JSONPath=".metadata.labels['vmoperator\\.vmware\\.com/active']"

// VirtualMachineClassInstance is the schema for the virtualmachineclassinstances API and
// represents the desired state and observed status of a virtualmachineclassinstance
// resource.
type VirtualMachineClassInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineClassInstanceSpec   `json:"spec,omitempty"`
	Status VirtualMachineClassInstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineClassInstanceList contains a list of VirtualMachineClassInstance.
type VirtualMachineClassInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineClassInstance `json:"items"`
}

func init() {
	objectTypes = append(objectTypes, &VirtualMachineClassInstance{}, &VirtualMachineClassInstanceList{})
}
