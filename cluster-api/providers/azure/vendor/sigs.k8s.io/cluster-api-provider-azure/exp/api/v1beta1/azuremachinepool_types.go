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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/errors"
)

const (
	// MachinePoolNameLabel indicates the AzureMachinePool name the AzureMachinePoolMachine belongs.
	MachinePoolNameLabel = "azuremachinepool.infrastructure.cluster.x-k8s.io/machine-pool"

	// RollingUpdateAzureMachinePoolDeploymentStrategyType replaces AzureMachinePoolMachines with older models with
	// AzureMachinePoolMachines based on the latest model.
	// i.e. gradually scale down the old AzureMachinePoolMachines and scale up the new ones.
	RollingUpdateAzureMachinePoolDeploymentStrategyType AzureMachinePoolDeploymentStrategyType = "RollingUpdate"

	// OldestDeletePolicyType will delete machines with the oldest creation date first.
	OldestDeletePolicyType AzureMachinePoolDeletePolicyType = "Oldest"
	// NewestDeletePolicyType will delete machines with the newest creation date first.
	NewestDeletePolicyType AzureMachinePoolDeletePolicyType = "Newest"
	// RandomDeletePolicyType will delete machines in random order.
	RandomDeletePolicyType AzureMachinePoolDeletePolicyType = "Random"
)

type (
	// AzureMachinePoolMachineTemplate defines the template for an AzureMachine.
	AzureMachinePoolMachineTemplate struct {
		// VMSize is the size of the Virtual Machine to build.
		// See https://learn.microsoft.com/rest/api/compute/virtualmachines/createorupdate#virtualmachinesizetypes
		VMSize string `json:"vmSize"`

		// Image is used to provide details of an image to use during VM creation.
		// If image details are omitted the image will default the Azure Marketplace "capi" offer,
		// which is based on Ubuntu.
		// +kubebuilder:validation:nullable
		// +optional
		Image *infrav1.Image `json:"image,omitempty"`

		// OSDisk contains the operating system disk information for a Virtual Machine
		OSDisk infrav1.OSDisk `json:"osDisk"`

		// DataDisks specifies the list of data disks to be created for a Virtual Machine
		// +optional
		DataDisks []infrav1.DataDisk `json:"dataDisks,omitempty"`

		// SSHPublicKey is the SSH public key string, base64-encoded to add to a Virtual Machine. Linux only.
		// Refer to documentation on how to set up SSH access on Windows instances.
		// +optional
		SSHPublicKey string `json:"sshPublicKey"`

		// Deprecated: AcceleratedNetworking should be set in the networkInterfaces field.
		// +optional
		AcceleratedNetworking *bool `json:"acceleratedNetworking,omitempty"`

		// Diagnostics specifies the diagnostics settings for a virtual machine.
		// If not specified then Boot diagnostics (Managed) will be enabled.
		// +optional
		Diagnostics *infrav1.Diagnostics `json:"diagnostics,omitempty"`

		// TerminateNotificationTimeout enables or disables VMSS scheduled events termination notification with specified timeout
		// allowed values are between 5 and 15 (mins)
		// +optional
		TerminateNotificationTimeout *int `json:"terminateNotificationTimeout,omitempty"`

		// SecurityProfile specifies the Security profile settings for a virtual machine.
		// +optional
		SecurityProfile *infrav1.SecurityProfile `json:"securityProfile,omitempty"`

		// SpotVMOptions allows the ability to specify the Machine should use a Spot VM
		// +optional
		SpotVMOptions *infrav1.SpotVMOptions `json:"spotVMOptions,omitempty"`

		// Deprecated: SubnetName should be set in the networkInterfaces field.
		// +optional
		SubnetName string `json:"subnetName,omitempty"`

		// VMExtensions specifies a list of extensions to be added to the scale set.
		// +optional
		VMExtensions []infrav1.VMExtension `json:"vmExtensions,omitempty"`

		// NetworkInterfaces specifies a list of network interface configurations.
		// If left unspecified, the VM will get a single network interface with a
		// single IPConfig in the subnet specified in the cluster's node subnet field.
		// The primary interface will be the first networkInterface specified (index 0) in the list.
		// +optional
		NetworkInterfaces []infrav1.NetworkInterface `json:"networkInterfaces,omitempty"`
	}

	// AzureMachinePoolSpec defines the desired state of AzureMachinePool.
	AzureMachinePoolSpec struct {
		// Location is the Azure region location e.g. westus2
		Location string `json:"location"`

		// Template contains the details used to build a replica virtual machine within the Machine Pool
		Template AzureMachinePoolMachineTemplate `json:"template"`

		// AdditionalTags is an optional set of tags to add to an instance, in addition to the ones added by default by the
		// Azure provider. If both the AzureCluster and the AzureMachine specify the same tag name with different values, the
		// AzureMachine's value takes precedence.
		// +optional
		AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

		// ProviderID is the identification ID of the Virtual Machine Scale Set
		// +optional
		ProviderID string `json:"providerID,omitempty"`

		// ProviderIDList are the identification IDs of machine instances provided by the provider.
		// This field must match the provider IDs as seen on the node objects corresponding to a machine pool's machine instances.
		// +optional
		ProviderIDList []string `json:"providerIDList,omitempty"`

		// Identity is the type of identity used for the Virtual Machine Scale Set.
		// The type 'SystemAssigned' is an implicitly created identity.
		// The generated identity will be assigned a Subscription contributor role.
		// The type 'UserAssigned' is a standalone Azure resource provided by the user
		// and assigned to the VM
		// +kubebuilder:default=None
		// +optional
		Identity infrav1.VMIdentity `json:"identity,omitempty"`

		// SystemAssignedIdentityRole defines the role and scope to assign to the system assigned identity.
		// +optional
		SystemAssignedIdentityRole *infrav1.SystemAssignedIdentityRole `json:"systemAssignedIdentityRole,omitempty"`

		// UserAssignedIdentities is a list of standalone Azure identities provided by the user
		// The lifecycle of a user-assigned identity is managed separately from the lifecycle of
		// the AzureMachinePool.
		// See https://learn.microsoft.com/azure/active-directory/managed-identities-azure-resources/how-to-manage-ua-identity-cli
		// +optional
		UserAssignedIdentities []infrav1.UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`

		// Deprecated: RoleAssignmentName should be set in the systemAssignedIdentityRole field.
		// +optional
		RoleAssignmentName string `json:"roleAssignmentName,omitempty"`

		// The deployment strategy to use to replace existing AzureMachinePoolMachines with new ones.
		// +optional
		// +kubebuilder:default={type: "RollingUpdate", rollingUpdate: {maxSurge: 1, maxUnavailable: 0, deletePolicy: Oldest}}
		Strategy AzureMachinePoolDeploymentStrategy `json:"strategy,omitempty"`

		// OrchestrationMode specifies the orchestration mode for the Virtual Machine Scale Set
		// +kubebuilder:default=Uniform
		OrchestrationMode infrav1.OrchestrationModeType `json:"orchestrationMode,omitempty"`

		// PlatformFaultDomainCount specifies the number of fault domains that the Virtual Machine Scale Set can use.
		// The count determines the spreading algorithm of the Azure fault domain.
		// +optional
		PlatformFaultDomainCount *int32 `json:"platformFaultDomainCount,omitempty"`

		// ZoneBalane dictates whether to force strictly even Virtual Machine distribution cross x-zones in case there is zone outage.
		// +optional
		ZoneBalance *bool `json:"zoneBalance,omitempty"`
	}

	// AzureMachinePoolDeploymentStrategyType is the type of deployment strategy employed to rollout a new version of
	// the AzureMachinePool.
	AzureMachinePoolDeploymentStrategyType string

	// AzureMachinePoolDeploymentStrategy describes how to replace existing machines with new ones.
	AzureMachinePoolDeploymentStrategy struct {
		// Type of deployment. Currently the only supported strategy is RollingUpdate
		// +optional
		// +kubebuilder:validation:Enum=RollingUpdate
		// +optional
		// +kubebuilder:default=RollingUpdate
		Type AzureMachinePoolDeploymentStrategyType `json:"type,omitempty"`

		// Rolling update config params. Present only if
		// MachineDeploymentStrategyType = RollingUpdate.
		// +optional
		RollingUpdate *MachineRollingUpdateDeployment `json:"rollingUpdate,omitempty"`
	}

	// AzureMachinePoolDeletePolicyType is the type of DeletePolicy employed to select machines to be deleted during an
	// upgrade.
	AzureMachinePoolDeletePolicyType string

	// MachineRollingUpdateDeployment is used to control the desired behavior of rolling update.
	MachineRollingUpdateDeployment struct {
		// The maximum number of machines that can be unavailable during the update.
		// Value can be an absolute number (ex: 5) or a percentage of desired
		// machines (ex: 10%).
		// Absolute number is calculated from percentage by rounding down.
		// This can not be 0 if MaxSurge is 0.
		// Defaults to 0.
		// Example: when this is set to 30%, the old MachineSet can be scaled
		// down to 70% of desired machines immediately when the rolling update
		// starts. Once new machines are ready, old MachineSet can be scaled
		// down further, followed by scaling up the new MachineSet, ensuring
		// that the total number of machines available at all times
		// during the update is at least 70% of desired machines.
		// +optional
		// +kubebuilder:default:=0
		MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

		// The maximum number of machines that can be scheduled above the
		// desired number of machines.
		// Value can be an absolute number (ex: 5) or a percentage of
		// desired machines (ex: 10%).
		// This can not be 0 if MaxUnavailable is 0.
		// Absolute number is calculated from percentage by rounding up.
		// Defaults to 1.
		// Example: when this is set to 30%, the new MachineSet can be scaled
		// up immediately when the rolling update starts, such that the total
		// number of old and new machines do not exceed 130% of desired
		// machines. Once old machines have been killed, new MachineSet can
		// be scaled up further, ensuring that total number of machines running
		// at any time during the update is at most 130% of desired machines.
		// +optional
		// +kubebuilder:default:=1
		MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`

		// DeletePolicy defines the policy used by the MachineDeployment to identify nodes to delete when downscaling.
		// Valid values are "Random, "Newest", "Oldest"
		// When no value is supplied, the default is Oldest
		// +optional
		// +kubebuilder:validation:Enum=Random;Newest;Oldest
		// +kubebuilder:default:=Oldest
		DeletePolicy AzureMachinePoolDeletePolicyType `json:"deletePolicy,omitempty"`
	}

	// AzureMachinePoolStatus defines the observed state of AzureMachinePool.
	AzureMachinePoolStatus struct {
		// Ready is true when the provider resource is ready.
		// +optional
		Ready bool `json:"ready"`

		// Replicas is the most recently observed number of replicas.
		// +optional
		Replicas int32 `json:"replicas"`

		// Instances is the VM instance status for each VM in the VMSS
		// +optional
		Instances []*AzureMachinePoolInstanceStatus `json:"instances,omitempty"`

		// Image is the current image used in the AzureMachinePool. When the spec image is nil, this image is populated
		// with the details of the defaulted Azure Marketplace "capi" offer.
		// +optional
		Image *infrav1.Image `json:"image,omitempty"`

		// Version is the Kubernetes version for the current VMSS model
		// +optional
		Version string `json:"version"`

		// ProvisioningState is the provisioning state of the Azure virtual machine.
		// +optional
		ProvisioningState *infrav1.ProvisioningState `json:"provisioningState,omitempty"`

		// FailureReason will be set in the event that there is a terminal problem
		// reconciling the MachinePool and will contain a succinct value suitable
		// for machine interpretation.
		//
		// This field should not be set for transitive errors that a controller
		// faces that are expected to be fixed automatically over
		// time (like service outages), but instead indicate that something is
		// fundamentally wrong with the MachinePool's spec or the configuration of
		// the controller, and that manual intervention is required. Examples
		// of terminal errors would be invalid combinations of settings in the
		// spec, values that are unsupported by the controller, or the
		// responsible controller itself being critically misconfigured.
		//
		// Any transient errors that occur during the reconciliation of MachinePools
		// can be added as events to the MachinePool object and/or logged in the
		// controller's output.
		// +optional
		FailureReason *errors.MachineStatusError `json:"failureReason,omitempty"`

		// FailureMessage will be set in the event that there is a terminal problem
		// reconciling the MachinePool and will contain a more verbose string suitable
		// for logging and human consumption.
		//
		// This field should not be set for transitive errors that a controller
		// faces that are expected to be fixed automatically over
		// time (like service outages), but instead indicate that something is
		// fundamentally wrong with the MachinePool's spec or the configuration of
		// the controller, and that manual intervention is required. Examples
		// of terminal errors would be invalid combinations of settings in the
		// spec, values that are unsupported by the controller, or the
		// responsible controller itself being critically misconfigured.
		//
		// Any transient errors that occur during the reconciliation of MachinePools
		// can be added as events to the MachinePool object and/or logged in the
		// controller's output.
		// +optional
		FailureMessage *string `json:"failureMessage,omitempty"`

		// Conditions defines current service state of the AzureMachinePool.
		// +optional
		Conditions clusterv1.Conditions `json:"conditions,omitempty"`

		// LongRunningOperationStates saves the state for Azure long-running operations so they can be continued on the
		// next reconciliation loop.
		// +optional
		LongRunningOperationStates infrav1.Futures `json:"longRunningOperationStates,omitempty"`

		// InfrastructureMachineKind is the kind of the infrastructure resources behind MachinePool Machines.
		// +optional
		InfrastructureMachineKind string `json:"infrastructureMachineKind,omitempty"`
	}

	// AzureMachinePoolInstanceStatus provides status information for each instance in the VMSS.
	AzureMachinePoolInstanceStatus struct {
		// Version defines the Kubernetes version for the VM Instance
		// +optional
		Version string `json:"version"`

		// ProvisioningState is the provisioning state of the Azure virtual machine instance.
		// +optional
		ProvisioningState *infrav1.ProvisioningState `json:"provisioningState"`

		// ProviderID is the provider identification of the VMSS Instance
		// +optional
		ProviderID string `json:"providerID"`

		// InstanceID is the identification of the Machine Instance within the VMSS
		// +optional
		InstanceID string `json:"instanceID"`

		// InstanceName is the name of the Machine Instance within the VMSS
		// +optional
		InstanceName string `json:"instanceName"`

		// LatestModelApplied indicates the instance is running the most up-to-date VMSS model. A VMSS model describes
		// the image version the VM is running. If the instance is not running the latest model, it means the instance
		// may not be running the version of Kubernetes the Machine Pool has specified and needs to be updated.
		LatestModelApplied bool `json:"latestModelApplied"`
	}

	// +kubebuilder:object:root=true
	// +kubebuilder:subresource:status
	// +kubebuilder:resource:path=azuremachinepools,scope=Namespaced,categories=cluster-api,shortName=amp
	// +kubebuilder:storageversion
	// +kubebuilder:printcolumn:name="Replicas",type="string",JSONPath=".status.replicas",description="AzureMachinePool replicas count"
	// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="AzureMachinePool replicas count"
	// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.provisioningState",description="Azure VMSS provisioning state"
	// +kubebuilder:printcolumn:name="Cluster",type="string",priority=1,JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AzureMachinePool belongs"
	// +kubebuilder:printcolumn:name="MachinePool",type="string",priority=1,JSONPath=".metadata.ownerReferences[?(@.kind==\"MachinePool\")].name",description="MachinePool object to which this AzureMachinePool belongs"
	// +kubebuilder:printcolumn:name="VMSS ID",type="string",priority=1,JSONPath=".spec.providerID",description="Azure VMSS ID"
	// +kubebuilder:printcolumn:name="VM Size",type="string",priority=1,JSONPath=".spec.template.vmSize",description="Azure VM Size"
	// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of this AzureMachinePool"

	// AzureMachinePool is the Schema for the azuremachinepools API.
	AzureMachinePool struct {
		metav1.TypeMeta   `json:",inline"`
		metav1.ObjectMeta `json:"metadata,omitempty"`

		Spec   AzureMachinePoolSpec   `json:"spec,omitempty"`
		Status AzureMachinePoolStatus `json:"status,omitempty"`
	}

	// +kubebuilder:object:root=true

	// AzureMachinePoolList contains a list of AzureMachinePools.
	AzureMachinePoolList struct {
		metav1.TypeMeta `json:",inline"`
		metav1.ListMeta `json:"metadata,omitempty"`
		Items           []AzureMachinePool `json:"items"`
	}
)

// GetConditions returns the list of conditions for an AzureMachinePool API object.
func (amp *AzureMachinePool) GetConditions() clusterv1.Conditions {
	return amp.Status.Conditions
}

// SetConditions will set the given conditions on an AzureMachinePool object.
func (amp *AzureMachinePool) SetConditions(conditions clusterv1.Conditions) {
	amp.Status.Conditions = conditions
}

// GetFutures returns the list of long running operation states for an AzureMachinePool API object.
func (amp *AzureMachinePool) GetFutures() infrav1.Futures {
	return amp.Status.LongRunningOperationStates
}

// SetFutures will set the given long running operation states on an AzureMachinePool object.
func (amp *AzureMachinePool) SetFutures(futures infrav1.Futures) {
	amp.Status.LongRunningOperationStates = futures
}

func init() {
	SchemeBuilder.Register(&AzureMachinePool{}, &AzureMachinePoolList{})
}
