/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzureMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It contains Azure-specific status information.
// +k8s:openapi-gen=true
type AzureMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// VMID is the ID of the virtual machine created in Azure.
	// +optional
	VMID *string `json:"vmId,omitempty"`

	// VMState is the provisioning state of the Azure virtual machine.
	// +optional
	VMState *VMState `json:"vmState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status.
	// +optional
	Conditions []AzureMachineProviderCondition `json:"conditions,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

func init() {
	SchemeBuilder.Register(&AzureMachineProviderStatus{})
}
