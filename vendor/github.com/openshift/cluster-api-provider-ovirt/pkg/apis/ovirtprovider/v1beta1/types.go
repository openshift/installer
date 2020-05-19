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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtMachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Ovirt VM. It is used by the Ovirt machine actuator to create a single machine instance.
type OvirtMachineProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with oVirt credentials.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	// Id is the UUID of the VM
	Id string `json:"id"`

	// Name is the VM name
	Name string `json:"name"`

	// The VM template this instance will be created from.
	TemplateName string `json:"template_name"`

	// the oVirt cluster this VM instance belongs too.
	ClusterId string `json:"cluster_id"`

	// InstanceTypeId defines the VM instance type and overrides
	// the hardware parameters of the created VM, including cpu and memory.
	// If InstanceTypeId is passed, all memory and cpu variables will be ignored.
	InstanceTypeId string `json:"instance_type_id,omitempty"`

	// CPU defines the VM CPU.
	CPU *CPU `json:"cpu,omitempty"`

	// MemoryMB is the size of a VM's memory in MiBs.
	MemoryMB int32 `json:"memory_mb,omitempty"`

	// OSDisk is the the root disk of the node.
	OSDisk *Disk `json:"os_disk,omitempty"`

	// VMType defines the workload type the instance will
	// be used for and this effects the instance parameters.
	// One of "desktop, server, high_performance"
	VMType string `json:"type,omitempty"`

	// NetworkInterfaces defines the list of the network interfaces of the VM.
	// All network interfaces from the template are discarded and new ones will
	// be created, unless the list is empty or nil
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces,omitempty"`
}

// CPU defines the VM cpu, made of (Sockets * Cores * Threads)
type CPU struct {
	// Sockets is the number of sockets for a VM.
	// Total CPUs is (Sockets * Cores * Threads)
	Sockets int32 `json:"sockets"`

	// Cores is the number of cores per socket.
	// Total CPUs is (Sockets * Cores * Threads)
	Cores int32 `json:"cores"`

	// Thread is the number of thread per core.
	// Total CPUs is (Sockets * Cores * Threads)
	Threads int32 `json:"threads"`
}

type Disk struct {
	// SizeGB size of the bootable disk in GiB.
	SizeGB int64 `json:"size_gb"`
}

// NetworkInterface defines a VM network interface
type NetworkInterface struct {
	// VNICProfileID the id of the vNic profile
	VNICProfileID string `json:"vnic_profile_id"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtClusterProviderSpec of an oVirt cluster
// +k8s:openapi-gen=true
type OvirtClusterProviderSpec struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtClusterProviderStatus
// +k8s:openapi-gen=true
type OvirtClusterProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// CACertificate is a PEM encoded CA Certificate for the control plane nodes.
	CACertificate []byte

	// CAPrivateKey is a PEM encoded PKCS1 CA PrivateKey for the control plane nodes.
	CAPrivateKey []byte
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OvirtMachineProviderStatus
type OvirtMachineProviderStatus struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// InstanceID is the ID of the instance in oVirt
	// +optional
	InstanceID *string `json:"instanceId,omitempty"`

	// InstanceState is the provisioning state of the oVirt Instance.
	// +optional
	InstanceState *string `json:"instanceState,omitempty"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []OvirtMachineProviderCondition `json:"conditions,omitempty"`
}

// OvirtMachineProviderConditionType is a valid value for OvirtMachineProviderCondition.Type
type OvirtMachineProviderConditionType string

// Valid conditions for an oVirt machine instance
const (
	// MachineCreated indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreated OvirtMachineProviderConditionType = "MachineCreated"
)

// OvirtMachineProviderCondition is a condition in a OvirtMachineProviderStatus
type OvirtMachineProviderCondition struct {
	// Type is the type of the condition.
	Type OvirtMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
}

func init() {
	SchemeBuilder.Register(&OvirtMachineProviderSpec{})
	SchemeBuilder.Register(&OvirtMachineProviderStatus{})
	SchemeBuilder.Register(&OvirtClusterProviderSpec{})
	SchemeBuilder.Register(&OvirtClusterProviderStatus{})
}
