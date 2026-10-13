/*
Copyright 2021 The Kubernetes Authors.

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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	"sigs.k8s.io/cluster-api/errors"
)

// VSphereMachineVolume defines a PVC attachment.
type VSphereMachineVolume struct {
	// Name is suffix used to name this PVC as: VSphereMachine.Name + "-" + Name
	Name string `json:"name"`
	// Capacity is the PVC capacity
	Capacity corev1.ResourceList `json:"capacity"`
	// StorageClass defaults to VSphereMachineSpec.StorageClass
	// +optional
	StorageClass string `json:"storageClass,omitempty"`
}

// VSphereMachineSpec defines the desired state of VSphereMachine.
type VSphereMachineSpec struct {
	// ProviderID is the virtual machine's BIOS UUID formatted as
	// vsphere://12345678-1234-1234-1234-123456789abc.
	// This is required at runtime by CAPI. Do not remove this field.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	// FailureDomain is the failure domain the machine will be created in.
	// Must match a key in the FailureDomains map stored on the cluster object.
	// +optional
	FailureDomain *string `json:"failureDomain,omitempty"`

	// ImageName is the name of the base image used when specifying the
	// underlying virtual machine
	ImageName string `json:"imageName"`

	// ClassName is the name of the class used when specifying the underlying
	// virtual machine
	ClassName string `json:"className"`

	// StorageClass is the name of the storage class used when specifying the
	// underlying virtual machine.
	// +optional
	StorageClass string `json:"storageClass,omitempty"`

	// Volumes is the set of PVCs to be created and attached to the VSphereMachine
	// +optional
	Volumes []VSphereMachineVolume `json:"volumes,omitempty"`

	// network is the network configuration for the VSphereMachine
	// +optional
	Network VSphereMachineNetworkSpec `json:"network,omitempty,omitzero"`

	// PowerOffMode describes the desired behavior when powering off a VM.
	//
	// There are three, supported power off modes: hard, soft, and
	// trySoft. The first mode, hard, is the equivalent of a physical
	// system's power cord being ripped from the wall. The soft mode
	// requires the VM's guest to have VM Tools installed and attempts to
	// gracefully shut down the VM. Its variant, trySoft, first attempts
	// a graceful shutdown, and if that fails or the VM is not in a powered off
	// state after reaching 5 minutes timeout, the VM is halted.
	//
	// If omitted, the mode defaults to hard.
	//
	// +optional
	// +kubebuilder:default=hard
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// MinHardwareVersion specifies the desired minimum hardware version
	// for this VM. Setting this field will ensure that the hardware version
	// of the VM is at least set to the specified value.
	// The expected format of the field is vmx-15.
	//
	// +optional
	MinHardwareVersion string `json:"minHardwareVersion,omitempty"`

	// NamingStrategy allows configuring the naming strategy used when calculating the name of the VirtualMachine.
	// +optional
	NamingStrategy *VirtualMachineNamingStrategy `json:"namingStrategy,omitempty"`
}

// VSphereMachineNetworkSpec defines the network configuration of a VSphereMachine.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkSpec struct {
	// interfaces is the list of network interfaces attached to this VSphereMachine.
	//
	// +optional
	Interfaces InterfacesSpec `json:"interfaces,omitempty,omitzero"`
}

// IsDefined returns true if the VSphereMachineNetworkSpec is defined.
func (r *VSphereMachineNetworkSpec) IsDefined() bool {
	return !reflect.DeepEqual(r, &VSphereMachineNetworkSpec{})
}

// InterfacesSpec defines all the network interfaces of a VSphereMachine from Kubernetes perspective.
// +kubebuilder:validation:MinProperties=1
type InterfacesSpec struct {
	// primary is the primary network interface.
	//
	// It is used to connect the Kubernetes primary network for Load balancer,
	// Service discovery, Pod traffic and management traffic etc.
	// Leave it unset if you don't want to customize the primary network and interface.
	// Customization is only supported with network provider NSX-VPC.
	// It should be set only when VSphereCluster spec.network.nsxVPC.createSubnetSet is set to false.
	//
	// +optional
	Primary InterfaceSpec `json:"primary,omitempty,omitzero"`

	// secondary are the secondary network interfaces.
	//
	// It is used for any purpose like deploying Antrea secondary network,
	// Multus, mounting NFS etc.
	// Secondary network is supported with network provider NSX-VPC and vsphere-network.
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=9
	// +listType=atomic
	// +optional
	Secondary []SecondaryInterfaceSpec `json:"secondary,omitempty"`
}

// IsDefined returns true if the InterfacesSpec is defined.
func (r *InterfacesSpec) IsDefined() bool {
	return !reflect.DeepEqual(r, &InterfacesSpec{})
}

// SecondaryInterfaceSpec defines a secondary network interface for a VSphereMachine.
type SecondaryInterfaceSpec struct {
	// name describes the unique name of this network interface, used to
	// distinguish it from other network interfaces attached to this VSphereMachine.
	//
	// +kubebuilder:validation:Pattern="^[a-z0-9]{2,}$"
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=15
	// +required
	Name string `json:"name,omitempty"`

	InterfaceSpec `json:",inline"`
}

// InterfaceSpec defines properties of a network interface.
type InterfaceSpec struct {
	// network is the name of the network resource to which this interface is
	// connected.
	// +required
	Network InterfaceNetworkReference `json:"network,omitempty,omitzero"`

	// mtu is the Maximum Transmission Unit size in bytes.
	//
	// +kubebuilder:validation:Minimum=576
	// +kubebuilder:validation:Maximum=9000
	// +optional
	MTU int32 `json:"mtu,omitempty"`

	// routes is a list of optional, static routes.
	//
	// Please note this feature is available only with the following bootstrap
	// providers: CloudInit.
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +listType=atomic
	// +optional
	Routes []RouteSpec `json:"routes,omitempty"`
}

// IsDefined returns true if the InterfaceSpec is defined.
func (r *InterfaceSpec) IsDefined() bool {
	return !reflect.DeepEqual(r, &InterfaceSpec{})
}

// InterfaceNetworkReference describes a reference to another object in the same
// namespace as the referrer.
type InterfaceNetworkReference struct {
	// kind of the remediation template.
	// kind must consist of alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
	Kind string `json:"kind,omitempty"`

	// name of the remediation template.
	// name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name,omitempty"`

	// apiVersion of the remediation template.
	// apiVersion must be fully qualified domain name followed by / and a version.
	// NOTE: This field must be kept in sync with the APIVersion of the remediation template.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=317
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*\/[a-z]([-a-z0-9]*[a-z0-9])?$`
	APIVersion string `json:"apiVersion,omitempty"`
}

// GroupVersionKind gets the GroupVersionKind for an InterfaceNetworkReference.
func (r *InterfaceNetworkReference) GroupVersionKind() schema.GroupVersionKind {
	return schema.FromAPIVersionAndKind(r.APIVersion, r.Kind)
}

// RouteSpec defines a static route for a guest.
type RouteSpec struct {
	// to is an IP4 CIDR. IP6 is not supported yet.
	// Examples: 192.168.1.0/24, 192.168.100.100/32, 0.0.0.0/0
	//
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}\/[0-9]{1,2}$`
	// +kubebuilder:validation:MinLength=9
	// +kubebuilder:validation:MaxLength=18
	// +required
	To string `json:"to,omitempty"`

	// via is an IP4 address. IP6 is not supported yet.
	//
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}$`
	// +kubebuilder:validation:MinLength=7
	// +kubebuilder:validation:MaxLength=15
	// +required
	Via string `json:"via,omitempty"`
}

// VirtualMachineNamingStrategy defines the naming strategy for the VirtualMachines.
type VirtualMachineNamingStrategy struct {
	// Template defines the template to use for generating the name of the VirtualMachine object.
	// If not defined, it will fall back to `{{ .machine.name }}`.
	// The templating has the following data available:
	// * `.machine.name`: The name of the Machine object.
	// The templating also has the following funcs available:
	// * `trimSuffix`: same as strings.TrimSuffix
	// * `trunc`: truncates a string, e.g. `trunc 2 "hello"` or `trunc -2 "hello"`
	// Notes:
	// * While the template offers some flexibility, we would like the name to link to the Machine name
	//   to ensure better user experience when troubleshooting
	// * Generated names must be valid Kubernetes names as they are used to create a VirtualMachine object
	//   and usually also as the name of the Node object.
	// * Names are automatically truncated at 63 characters. Please note that this can lead to name conflicts,
	//   so we highly recommend to use a template which leads to a name shorter than 63 characters.
	// +optional
	Template *string `json:"template,omitempty"`
}

// VSphereMachineStatus defines the observed state of VSphereMachine.
type VSphereMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// This is required at runtime by CAPI. Do not remove this field.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the instance associated addresses.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=atomic
	// +optional
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// ID is used to identify the virtual machine.
	// +optional
	ID *string `json:"vmID,omitempty"`

	// IPAddr is the IP address used to access the virtual machine.
	// +optional
	IPAddr string `json:"vmIp,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the Machine and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the Machine's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// VMStatus is used to identify the virtual machine status.
	// +optional
	VMStatus VirtualMachineState `json:"vmstatus,omitempty"`

	// Conditions defines current service state of the VSphereMachine.
	// +optional
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in VSphereMachine's status with the V1Beta2 version.
	// +optional
	V1Beta2 *VSphereMachineV1Beta2Status `json:"v1beta2,omitempty"`

	// network describes the observed state of the VM's network configuration.
	// Please note much of the network status information is only available if
	// the guest has VM Tools installed.
	// +optional
	Network VSphereMachineNetworkStatus `json:"network,omitempty,omitzero"`
}

// VSphereMachineV1Beta2Status groups all the fields that will be added or modified in VSphereMachineStatus with the V1Beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereMachineV1Beta2Status struct {
	// conditions represents the observations of a VSphereMachine's current state.
	// Known condition types are Ready, VirtualMachineProvisioned and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// VSphereMachine is the Schema for the vspheremachines API
//
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Zone",type="string",JSONPath=".spec.failureDomain",description="Zone"
// +kubebuilder:printcolumn:name="ProviderID",type="string",JSONPath=".spec.providerID",description="Provider ID"
// +kubebuilder:printcolumn:name="IPAddr",type="string",JSONPath=".status.vmIp",description="IP address"
type VSphereMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereMachineSpec   `json:"spec,omitempty"`
	Status VSphereMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereMachineList contains a list of VSphereMachine.
type VSphereMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereMachine `json:"items"`
}

// GetConditions returns the conditions for the VSphereMachine.
func (r *VSphereMachine) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets conditions on the VSphereMachine.
func (r *VSphereMachine) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (r *VSphereMachine) GetV1Beta2Conditions() []metav1.Condition {
	if r.Status.V1Beta2 == nil {
		return nil
	}
	return r.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (r *VSphereMachine) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if r.Status.V1Beta2 == nil {
		r.Status.V1Beta2 = &VSphereMachineV1Beta2Status{}
	}
	r.Status.V1Beta2.Conditions = conditions
}

func init() {
	objectTypes = append(objectTypes, &VSphereMachine{}, &VSphereMachineList{})
}
