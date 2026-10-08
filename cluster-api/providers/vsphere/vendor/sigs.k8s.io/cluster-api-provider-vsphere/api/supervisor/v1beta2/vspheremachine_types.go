/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta2

import (
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/api/deprecated/errors"
)

// VSphereMachineVolume defines a PVC attachment.
type VSphereMachineVolume struct {
	// name is the suffix used to name this PVC as: VSphereMachine.Name + "-" + Name
	// Note: Use short values for the name as the max length for PVC names is 253 characters.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name,omitempty"`

	// capacity is the PVC capacity
	// +required
	Capacity corev1.ResourceList `json:"capacity,omitempty,omitzero"`

	// storageClass defaults to VSphereMachineSpec.StorageClass
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	StorageClass string `json:"storageClass,omitempty"`
}

// VSphereMachineSpec defines the desired state of VSphereMachine.
type VSphereMachineSpec struct {
	// providerID is the virtual machine's BIOS UUID formatted as
	// vsphere://12345678-1234-1234-1234-123456789abc
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	ProviderID string `json:"providerID,omitempty"`

	// imageName is the name of the base image used when specifying the
	// underlying virtual machine
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ImageName string `json:"imageName,omitempty"`

	// className is the name of the class used when specifying the underlying
	// virtual machine
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	ClassName string `json:"className,omitempty"`

	// storageClass is the name of the storage class used when specifying the
	// underlying virtual machine.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	StorageClass string `json:"storageClass,omitempty"`

	// volumes is the set of PVCs to be created and attached to the VSphereMachine
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=1024
	Volumes []VSphereMachineVolume `json:"volumes,omitempty"`

	// network is the network configuration for the VSphereMachine
	// +optional
	Network VSphereMachineNetworkSpec `json:"network,omitempty,omitzero"`

	// powerOffMode describes the desired behavior when powering off a VM.
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
	PowerOffMode VirtualMachinePowerOpMode `json:"powerOffMode,omitempty"`

	// minHardwareVersion specifies the desired minimum hardware version
	// for this VM. Setting this field will ensure that the hardware version
	// of the VM is at least set to the specified value.
	// The expected format of the field is vmx-15.
	//
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	MinHardwareVersion string `json:"minHardwareVersion,omitempty"`

	// naming allows configuring the naming strategy used when calculating the name of the VirtualMachine.
	// +optional
	Naming VirtualMachineNamingSpec `json:"naming,omitempty,omitzero"`

	// policies specifies a list of optional infrastructure policies to be applied to the virtual machine.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	Policies []PolicyRef `json:"policies,omitempty"`
}

// PolicyRef identifies an optional infrastructure policy to specify for the virtual machine by name and kind.
type PolicyRef struct {
	// name of the infrastructure policy.
	// name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name,omitempty"`

	// kind of the infrastructure policy object being referenced.
	// kind must be a CamelCase Kubernetes Kind starting with an uppercase letter (e.g. ComputePolicy).
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
	Kind string `json:"kind,omitempty"`
}

// VSphereMachineNetworkSpec defines the network configuration of a VSphereMachine.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineNetworkSpec struct {
	// interfaces is the list of network interfaces attached to this VSphereMachine.
	//
	// +optional
	Interfaces InterfacesSpec `json:"interfaces,omitempty,omitzero"`

	// vlans is a list of VLAN sub-interfaces to be configured on the secondary
	// network interfaces. Each VLAN is linked to a specific secondary interface
	// via the link field.
	//
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +listType=atomic
	// +optional
	VLANs []VLANSpec `json:"vlans,omitempty"`
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
	// networkRef is the name of the network resource to which this interface is
	// connected.
	// +required
	NetworkRef InterfaceNetworkReference `json:"networkRef,omitempty,omitzero"`

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

// InterfaceNetworkReference describes a reference to a network object in the same
// namespace as the referrer.
type InterfaceNetworkReference struct {
	// kind of the network object.
	// kind must consist of alphanumeric characters or '-', start with an alphabetic character, and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`
	Kind string `json:"kind,omitempty"`

	// name of the network object.
	// name must consist of lower case alphanumeric characters, '-' or '.', and must start and end with an alphanumeric character.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	// +kubebuilder:validation:Pattern=`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`
	Name string `json:"name,omitempty"`

	// apiVersion of the network object.
	// apiVersion must be fully qualified domain name followed by / and a version.
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

// VirtualMachineNamingSpec defines the naming strategy for the VirtualMachines.
// +kubebuilder:validation:MinProperties=1
type VirtualMachineNamingSpec struct {
	// template defines the template to use for generating the name of the VirtualMachine object.
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
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=1024
	Template string `json:"template,omitempty"`
}

// VSphereMachineStatus defines the observed state of VSphereMachine.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineStatus struct {
	// conditions represents the observations of a VSphereMachine's current state.
	// Known condition types are Ready, VirtualMachineProvisioned and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization provides observations of the VSphereMachine initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Machine provisioning.
	// +optional
	Initialization VSphereMachineInitializationStatus `json:"initialization,omitempty,omitzero"`

	// addresses contains the instance associated addresses.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=atomic
	// +optional
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// failureDomain is the failure domain where the VirtualMachine has been scheduled.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	FailureDomain string `json:"failureDomain,omitempty"`

	// biosUUID is the biosUUID of the virtual machine.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	BiosUUID string `json:"biosUUID,omitempty"`

	// phase is the phase of the VSphereMachine.
	// +optional
	Phase VSphereMachinePhase `json:"phase,omitempty"`

	// network describes the observed state of the VM's network configuration.
	// Please note much of the network status information is only available if
	// the guest has VM Tools installed.
	// +optional
	Network VSphereMachineNetworkStatus `json:"network,omitempty,omitzero"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereMachineDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereMachineInitializationStatus provides observations of the VSphereMachine initialization process.
// +kubebuilder:validation:MinProperties=1
type VSphereMachineInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Machine's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Machine provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// VSphereMachineDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereMachineDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	V1Beta1 *VSphereMachineV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereMachineV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereMachineV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereMachine.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// failureReason will be set in the event that there is a terminal problem
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
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"` //nolint:kubeapilinter // field will be removed when v1beta1 is removed

	// failureMessage will be set in the event that there is a terminal problem
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
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"` //nolint:kubeapilinter // field will be removed when v1beta1 is removed
}

// VLANSpec defines a VLAN sub-interface configuration linked to a secondary interface.
type VLANSpec struct {
	// name is the name of the VLAN sub-interface as it appears inside
	// the guest operating system.
	//
	// The name must conform to Linux network interface naming rules:
	// it must be between 1 and 15 characters long, start with an
	// alphanumeric character, and contain only alphanumeric characters,
	// hyphens, underscores, or dots.
	//
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:Pattern="^[a-zA-Z0-9][a-zA-Z0-9._-]*$"
	Name string `json:"name,omitempty"`

	// id is the VLAN ID used to tag traffic on this sub-interface,
	// a number between 0 and 4094.
	//
	// +required
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=4094
	ID *int32 `json:"id,omitempty"`

	// link is the name of the secondary interface this VLAN is associated with.
	//
	// +required
	// +kubebuilder:validation:Pattern="^[a-z0-9]{2,}$"
	// +kubebuilder:validation:MinLength=2
	// +kubebuilder:validation:MaxLength=15
	Link string `json:"link,omitempty"`
}

// VSphereMachine is the Schema for the vspheremachines API
//
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vspheremachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels['cluster\\.x-k8s\\.io/cluster-name']",description="Cluster"
// +kubebuilder:printcolumn:name="Class",type="string",JSONPath=".spec.className",description="VirtualMachineClass name"
// +kubebuilder:printcolumn:name="Provider ID",type="string",JSONPath=".spec.providerID",description="Provider ID",priority=10
// +kubebuilder:printcolumn:name="Failure domain",type="string",JSONPath=".status.failureDomain",description="The failure domain where the VSphereMachine has been scheduled"
// +kubebuilder:printcolumn:name="IP",type="string",JSONPath=`.status.addresses[?(@.type=="InternalIP")].address`,description="IP of the Machine",priority=10
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.imageName",description="Image name",priority=10
// +kubebuilder:printcolumn:name="Paused",type="string",JSONPath=`.status.conditions[?(@.type=="Paused")].status`,description="Reconciliation paused",priority=10
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.initialization.provisioned",description="VSphereMachine is provisioned"
// +kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase",description="VSphereMachine phase"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VSphereMachine"
// +kubebuilder:validation:XValidation:rule="has(self.spec.policies) == has(oldSelf.spec.policies) && (!has(self.spec.policies) || self.spec.policies == oldSelf.spec.policies)",message="policies are immutable after creation"
type VSphereMachine struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereMachine.
	// +required
	Spec VSphereMachineSpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereMachine.
	// +optional
	Status VSphereMachineStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// VSphereMachineList contains a list of VSphereMachine.
type VSphereMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereMachine `json:"items"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereMachine) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereMachine) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereMachineDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereMachineV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereMachine) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereMachine) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

func init() {
	objectTypes = append(objectTypes, &VSphereMachine{}, &VSphereMachineList{})
}
