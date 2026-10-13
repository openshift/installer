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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// MachineFinalizer allows ReconcileAzureMachine to clean up Azure resources associated with AzureMachine before
	// removing it from the apiserver.
	MachineFinalizer = "azuremachine.infrastructure.cluster.x-k8s.io"
)

// AzureMachineSpec defines the desired state of AzureMachine.
type AzureMachineSpec struct {
	// ProviderID is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderID *string `json:"providerID,omitempty"`

	VMSize string `json:"vmSize"`

	// FailureDomain is the failure domain unique identifier this Machine should be attached to,
	// as defined in Cluster API. This relates to an Azure Availability Zone
	// +optional
	FailureDomain *string `json:"failureDomain,omitempty"`

	// Image is used to provide details of an image to use during VM creation.
	// If image details are omitted, the default is to use an Azure Compute Gallery Image
	// from CAPZ's community gallery.
	// +kubebuilder:validation:nullable
	// +optional
	Image *Image `json:"image,omitempty"`

	// Identity is the type of identity used for the virtual machine.
	// The type 'SystemAssigned' is an implicitly created identity.
	// The generated identity will be assigned a Subscription contributor role.
	// The type 'UserAssigned' is a standalone Azure resource provided by the user
	// and assigned to the VM
	// +kubebuilder:default=None
	// +optional
	Identity VMIdentity `json:"identity,omitempty"`

	// UserAssignedIdentities is a list of standalone Azure identities provided by the user
	// The lifecycle of a user-assigned identity is managed separately from the lifecycle of
	// the AzureMachine.
	// See https://learn.microsoft.com/azure/active-directory/managed-identities-azure-resources/how-to-manage-ua-identity-cli
	// +optional
	UserAssignedIdentities []UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`

	// SystemAssignedIdentityRole defines the role and scope to assign to the system-assigned identity.
	// +optional
	SystemAssignedIdentityRole *SystemAssignedIdentityRole `json:"systemAssignedIdentityRole,omitempty"`

	// Deprecated: RoleAssignmentName should be set in the systemAssignedIdentityRole field.
	// +optional
	RoleAssignmentName string `json:"roleAssignmentName,omitempty"`

	// OSDisk specifies the parameters for the operating system disk of the machine
	OSDisk OSDisk `json:"osDisk"`

	// DataDisk specifies the parameters that are used to add one or more data disks to the machine
	// +optional
	DataDisks []DataDisk `json:"dataDisks,omitempty"`

	// SSHPublicKey is the SSH public key string, base64-encoded to add to a Virtual Machine. Linux only.
	// Refer to documentation on how to set up SSH access on Windows instances.
	// +optional
	SSHPublicKey string `json:"sshPublicKey,omitempty"`

	// AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
	// Azure provider. If both the AzureCluster and the AzureMachine specify the same tag name with different values, the
	// AzureMachine's value takes precedence.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// AdditionalCapabilities specifies additional capabilities enabled or disabled on the virtual machine.
	// +optional
	AdditionalCapabilities *AdditionalCapabilities `json:"additionalCapabilities,omitempty"`

	// AllocatePublicIP allows the ability to create dynamic public ips for machines where this value is true.
	// +optional
	AllocatePublicIP bool `json:"allocatePublicIP,omitempty"`

	// EnableIPForwarding enables IP Forwarding in Azure which is required for some CNI's to send traffic from a pods on one machine
	// to another. This is required for IpV6 with Calico in combination with User Defined Routes (set by the Azure Cloud Controller
	// manager). Default is false for disabled.
	// +optional
	EnableIPForwarding bool `json:"enableIPForwarding,omitempty"`

	// Deprecated: AcceleratedNetworking should be set in the networkInterfaces field.
	// +kubebuilder:validation:nullable
	// +optional
	AcceleratedNetworking *bool `json:"acceleratedNetworking,omitempty"`

	// Diagnostics specifies the diagnostics settings for a virtual machine.
	// If not specified then Boot diagnostics (Managed) will be enabled.
	// +optional
	Diagnostics *Diagnostics `json:"diagnostics,omitempty"`

	// SpotVMOptions allows the ability to specify the Machine should use a Spot VM
	// +optional
	SpotVMOptions *SpotVMOptions `json:"spotVMOptions,omitempty"`

	// SecurityProfile specifies the Security profile settings for a virtual machine.
	// +optional
	SecurityProfile *SecurityProfile `json:"securityProfile,omitempty"`

	// Deprecated: SubnetName should be set in the networkInterfaces field.
	// +optional
	SubnetName string `json:"subnetName,omitempty"`

	// DNSServers adds a list of DNS Server IP addresses to the VM NICs.
	// +optional
	DNSServers []string `json:"dnsServers,omitempty"`

	// DisableExtensionOperations specifies whether extension operations should be disabled on the virtual machine.
	// Use this setting only if VMExtensions are not supported by your image, as it disables CAPZ bootstrapping extension used for detecting Kubernetes bootstrap failure.
	// This may only be set to True when no extensions are configured on the virtual machine.
	// +optional
	DisableExtensionOperations *bool `json:"disableExtensionOperations,omitempty"`

	// DisableVMBootstrapExtension specifies whether the VM bootstrap extension should be disabled on the virtual machine.
	// Use this setting if you want to disable only the bootstrapping extension and not all extensions.
	// +optional
	DisableVMBootstrapExtension *bool `json:"disableVMBootstrapExtension,omitempty"`

	// VMExtensions specifies a list of extensions to be added to the virtual machine.
	// +optional
	VMExtensions []VMExtension `json:"vmExtensions,omitempty"`

	// NetworkInterfaces specifies a list of network interface configurations.
	// If left unspecified, the VM will get a single network interface with a
	// single IPConfig in the subnet specified in the cluster's node subnet field.
	// The primary interface will be the first networkInterface specified (index 0) in the list.
	// +optional
	NetworkInterfaces []NetworkInterface `json:"networkInterfaces,omitempty"`

	// CapacityReservationGroupID specifies the capacity reservation group resource id that should be
	// used for allocating the virtual machine.
	// The field size should be greater than 0 and the field input must start with '/'.
	// The input for capacityReservationGroupID must be similar to '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/capacityReservationGroups/{capacityReservationGroupName}'.
	// The keys which are used should be among 'subscriptions', 'providers' and 'resourcegroups' followed by valid ID or names respectively.
	// It is optional but may not be changed once set.
	// +optional
	CapacityReservationGroupID *string `json:"capacityReservationGroupID,omitempty"`
}

// SpotVMOptions defines the options relevant to running the Machine on Spot VMs.
type SpotVMOptions struct {
	// MaxPrice defines the maximum price the user is willing to pay for Spot VM instances
	// +optional
	MaxPrice *resource.Quantity `json:"maxPrice,omitempty"`

	// EvictionPolicy defines the behavior of the virtual machine when it is evicted. It can be either Delete or Deallocate.
	// +optional
	EvictionPolicy *SpotEvictionPolicy `json:"evictionPolicy,omitempty"`
}

// SystemAssignedIdentityRole defines the role and scope to assign to the system assigned identity.
type SystemAssignedIdentityRole struct {
	// Name is the name of the role assignment to create for a system assigned identity. It can be any valid UUID.
	// If not specified, a random UUID will be generated.
	// +optional
	Name string `json:"name,omitempty"`

	// DefinitionID is the ID of the role definition to create for a system assigned identity. It can be an Azure built-in role or a custom role.
	// Refer to built-in roles: https://learn.microsoft.com/en-us/azure/role-based-access-control/built-in-roles
	// +optional
	DefinitionID string `json:"definitionID,omitempty"`

	// Scope is the scope that the role assignment or definition applies to. The scope can be any REST resource instance.
	// If not specified, the scope will be the subscription.
	// +optional
	Scope string `json:"scope,omitempty"`
}

// AzureMachineStatus defines the observed state of AzureMachine.
type AzureMachineStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Addresses contains the Azure instance associated addresses.
	// +optional
	Addresses []corev1.NodeAddress `json:"addresses,omitempty"`

	// VMState is the provisioning state of the Azure virtual machine.
	// +optional
	VMState *ProvisioningState `json:"vmState,omitempty"`

	// ErrorReason will be set in the event that there is a terminal problem
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
	FailureReason *string `json:"failureReason,omitempty"`

	// ErrorMessage will be set in the event that there is a terminal problem
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

	// Conditions defines current service state of the AzureMachine.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// LongRunningOperationStates saves the states for Azure long-running operations so they can be continued on the
	// next reconciliation loop.
	// +optional
	LongRunningOperationStates Futures `json:"longRunningOperationStates,omitempty"`
}

// AdditionalCapabilities enables or disables a capability on the virtual machine.
type AdditionalCapabilities struct {
	// UltraSSDEnabled enables or disables Azure UltraSSD capability for the virtual machine.
	// Defaults to true if Ultra SSD data disks are specified,
	// otherwise it doesn't set the capability on the VM.
	// +optional
	UltraSSDEnabled *bool `json:"ultraSSDEnabled,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Cluster",type="string",priority=1,JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AzureMachine belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",priority=1,JSONPath=".status.conditions[?(@.type=='Ready')].message"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.vmState",description="Azure VM provisioning state"
// +kubebuilder:printcolumn:name="Machine",type="string",priority=1,JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object to which this AzureMachine belongs"
// +kubebuilder:printcolumn:name="VM ID",type="string",priority=1,JSONPath=".spec.providerID",description="Azure VM ID"
// +kubebuilder:printcolumn:name="VM Size",type="string",priority=1,JSONPath=".spec.vmSize",description="Azure VM Size"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of this AzureMachine"
// +kubebuilder:resource:path=azuremachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// AzureMachine is the Schema for the azuremachines API.
type AzureMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureMachineSpec   `json:"spec,omitempty"`
	Status AzureMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureMachineList contains a list of AzureMachine.
type AzureMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureMachine `json:"items"`
}

// GetConditions returns the list of conditions for an AzureMachine API object.
func (m *AzureMachine) GetConditions() clusterv1.Conditions {
	return m.Status.Conditions
}

// SetConditions will set the given conditions on an AzureMachine object.
func (m *AzureMachine) SetConditions(conditions clusterv1.Conditions) {
	m.Status.Conditions = conditions
}

// GetFutures returns the list of long running operation states for an AzureMachine API object.
func (m *AzureMachine) GetFutures() Futures {
	return m.Status.LongRunningOperationStates
}

// SetFutures will set the given long running operation states on an AzureMachine object.
func (m *AzureMachine) SetFutures(futures Futures) {
	m.Status.LongRunningOperationStates = futures
}

func init() {
	SchemeBuilder.Register(&AzureMachine{}, &AzureMachineList{})
}
