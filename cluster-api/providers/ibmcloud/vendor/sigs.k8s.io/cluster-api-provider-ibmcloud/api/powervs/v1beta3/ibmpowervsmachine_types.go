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

package v1beta3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

func init() {
	objectTypes = append(objectTypes, &IBMPowerVSMachine{}, &IBMPowerVSMachineList{})
}

// PowerVSProcessorType enum attribute to identify the PowerVS instance processor type.
type PowerVSProcessorType string

const (
	// IBMPowerVSMachineFinalizer allows IBMPowerVSMachineReconciler to clean up resources associated with IBMPowerVSMachine before
	// removing it from the apiserver.
	IBMPowerVSMachineFinalizer = "ibmpowervsmachine.infrastructure.cluster.x-k8s.io"
	// PowerVSProcessorTypeDedicated enum property to identify a Dedicated Power VS processor type.
	PowerVSProcessorTypeDedicated PowerVSProcessorType = "Dedicated"
	// PowerVSProcessorTypeShared enum property to identify a Shared Power VS processor type.
	PowerVSProcessorTypeShared PowerVSProcessorType = "Shared"
	// PowerVSProcessorTypeCapped enum property to identify a Capped Power VS processor type.
	PowerVSProcessorTypeCapped PowerVSProcessorType = "Capped"
	// DefaultIgnitionVersion represents default Ignition version generated for machine userdata.
	DefaultIgnitionVersion = "2.3"
)

// ImageSourceType defines the method used to resolve the machine image.
// +kubebuilder:validation:Enum=Reference;Import
type ImageSourceType string

const (
	// ImageSourceTypeReference specifies that the machine should use an existing image already available in PowerVS.
	ImageSourceTypeReference ImageSourceType = "Reference"

	// ImageSourceTypeImport specifies that the machine should use an IBMPowerVSImage CRD to import an image from COS.
	ImageSourceTypeImport ImageSourceType = "Import"
)

// IBMPowerVSMachineSpec defines the desired state of IBMPowerVSMachine.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSMachineSpec struct {
	// workspace identifies the PowerVS workspace where the instance will be created.
	// If omitted, the workspace is inherited from the associated IBMPowerVSCluster.
	// Supported identifiers are name and id.
	// More details: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server
	// +optional
	Workspace ResourceIdentifier `json:"workspace,omitempty,omitzero"`

	// network is the reference to the Network to use for this instance.
	// Supported identifiers in ResourceIdentifier are Name, ID, and RegEx and can be obtained from IBM Cloud UI or IBM Cloud CLI.
	// +optional
	Network ResourceIdentifier `json:"network,omitempty,omitzero"`

	// image specifies how to resolve the OS image used to create the instance.
	// +required
	Image IBMPowerVSMachineImage `json:"image,omitempty,omitzero"`

	// sshKey is the name of the SSH key pair provided to the VM for authenticating users.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	SSHKey string `json:"sshKey,omitempty"`

	// systemType is the System type used to host the instance.
	// systemType determines the number of cores and memory that is available.
	// Few of the supported SystemTypes are e980,s1022,s1122,e1050,e1080.
	// When omitted, this means that the user has no opinion and the platform is left to choose a
	// reasonable default, which is subject to change over time. The current default is s1022 which is generally available.
	// + This is not an enum because we expect other values to be added later which should be supported implicitly.
	// + The pattern validation allows any systemType matching PowerVS naming convention (lowercase letter + numbers).
	// + Dynamic validation against PowerVS API is performed by the controller during reconciliation.
	// + The controller validates the systemType against current PowerVS datacenter capabilities.
	// + If the systemType is not supported, the machine will be marked with InvalidMachineConfiguration condition.
	// +kubebuilder:validation:Pattern=`^[a-z][0-9]+$`
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=16
	SystemType string `json:"systemType,omitempty"`

	// processorType is the VM instance processor type.
	// It must be set to one of the following values: Dedicated, Capped or Shared.
	// Dedicated: resources are allocated for a specific client, The hypervisor makes a 1:1 binding of a partition’s processor to a physical processor core.
	// Shared: Shared among other clients.
	// Capped: Shared, but resources do not expand beyond those that are requested, the amount of CPU time is Capped to the value specified for the entitlement.
	// if the processorType is selected as Dedicated, then processors value cannot be fractional.
	// When omitted, this means that the user has no opinion and the platform is left to choose a
	// reasonable default, which is subject to change over time. The current default is Shared.
	// +kubebuilder:validation:Enum:="Dedicated";"Shared";"Capped";""
	// +optional
	ProcessorType PowerVSProcessorType `json:"processorType,omitempty"`

	// processors is the number of virtual processors in a virtual machine.
	// when the processorType is selected as Dedicated the processors value cannot be fractional.
	// maximum value for the Processors depends on the selected SystemType,
	// and minimum value for Processors depends on the selected ProcessorType, which can be found here: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server-on-cloud.
	// when ProcessorType is set as Shared or Capped, The minimum processors is 0.25.
	// when ProcessorType is set as Dedicated, The minimum processors is 1.
	// When omitted, this means that the user has no opinion and the platform is left to choose a
	// reasonable default, which is subject to change over time. The default is set based on the selected ProcessorType.
	// when ProcessorType selected as Dedicated, the default is set to 1.
	// when ProcessorType selected as Shared or Capped, the default is set to 0.25.
	// +optional
	Processors intstr.IntOrString `json:"processors,omitempty"`

	// memoryGiB is the size of a virtual machine's memory, in GiB.
	// maximum value for the MemoryGiB depends on the selected SystemType, which can be found here: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server-on-cloud
	// The minimum memory is 2 GiB.
	// When omitted, this means the user has no opinion and the platform is left to choose a reasonable
	// default, which is subject to change over time. The current default is 2.
	// +optional
	// +kubebuilder:validation:Minimum=2
	MemoryGiB int32 `json:"memoryGiB,omitempty"`

	// providerID is the unique identifier as specified by the cloud provider.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	ProviderID string `json:"providerID,omitempty"`
}

// IBMPowerVSMachineStatus defines the observed state of IBMPowerVSMachine.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSMachineStatus struct {
	// conditions represents the observations of a IBMPowerVSMachine's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization provides observations of the IBMPowerVSMachine initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Machine provisioning.
	// +optional
	Initialization IBMPowerVSMachineInitializationStatus `json:"initialization,omitempty,omitzero"`

	// instanceID is the instance ID.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	InstanceID string `json:"instanceID,omitempty"`

	// addresses contains the instance associated addresses.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	// +listType=atomic
	// +optional
	Addresses []clusterv1.MachineAddress `json:"addresses,omitempty"`

	// health is the health of the VM.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	Health string `json:"health,omitempty"`

	// instanceState is the status of the VM.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	InstanceState PowerVSInstanceState `json:"instanceState,omitempty"`

	// region specifies the Power VS Service instance region.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9\-_]+$
	Region string `json:"region,omitempty"`

	// zone specifies the Power VS Service instance zone.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9\-_]+$
	Zone string `json:"zone,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *IBMPowerVSMachineDeprecatedStatus `json:"deprecated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=ibmpowervsmachines,scope=Namespaced,categories=cluster-api
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this IBMPowerVSMachine belongs"
// +kubebuilder:printcolumn:name="Machine",type="string",priority=1,JSONPath=".metadata.ownerReferences[?(@.kind==\"Machine\")].name",description="Machine object to which this IBMPowerVSMachine belongs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMPowerVSMachine"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.initialization.provisioned	",description="Cluster infrastructure is ready for IBM PowerVS instances"
// +kubebuilder:printcolumn:name="Internal-IP",type="string",priority=1,JSONPath=".status.addresses[?(@.type==\"InternalIP\")].address",description="Instance Internal Addresses"
// +kubebuilder:printcolumn:name="External-IP",type="string",priority=1,JSONPath=".status.addresses[?(@.type==\"ExternalIP\")].address",description="Instance External Addresses"
// +kubebuilder:printcolumn:name="State",type="string",JSONPath=".status.instanceState",description="PowerVS instance state"
// +kubebuilder:printcolumn:name="Health",type="string",JSONPath=".status.health",description="PowerVS instance health"

// IBMPowerVSMachine is the Schema for the ibmpowervsmachines API.
type IBMPowerVSMachine struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of IBMPowerVSMachine
	// +required
	Spec IBMPowerVSMachineSpec `json:"spec,omitempty,omitzero"`

	// status defines the observed state of IBMPowerVSMachine
	// +optional
	Status IBMPowerVSMachineStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// IBMPowerVSMachineList contains a list of IBMPowerVSMachine.
type IBMPowerVSMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []IBMPowerVSMachine `json:"items"`
}

// IBMPowerVSMachineInitializationStatus provides observations of the IBMPowerVSMachine initialization process.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSMachineInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Machine's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Machine provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// IBMPowerVSMachineDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type IBMPowerVSMachineDeprecatedStatus struct {
	// v1beta2 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	V1Beta2 *IBMPowerVSMachineV1Beta2DeprecatedStatus `json:"v1beta2,omitempty"`
}

// IBMPowerVSMachineV1Beta2DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type IBMPowerVSMachineV1Beta2DeprecatedStatus struct {
	// conditions defines current service state of the VSphereMachine.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// IBMPowerVSMachineImage defines how to resolve the image for the machine.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Import' ? has(self.import) : !has(self.import)",message="import configuration is required when type is Import, and forbidden otherwise"
type IBMPowerVSMachineImage struct {
	// type defines whether to use an existing image in IBM Cloud or import a new one via the IBMPowerVSImage CRD.
	// +required
	Type ImageSourceType `json:"type,omitempty"`

	// reference contains the information to identify an existing image in the PowerVS workspace.
	// Supported identifiers are Name, ID, and RegEx.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// import is a reference to an IBMPowerVSImage CRD, which manages importing an image from an IBM COS Bucket.
	// +optional
	Import ImageReference `json:"import,omitempty,omitzero"`
}

// ImageReference is a reference to an IBMPowerVSImage resource.
type ImageReference struct {
	// name of the IBMPowerVSImage resource.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Name string `json:"name,omitempty"`
}

// GetConditions returns the observations of the operational state of the IBMPowerVSMachine resource.
func (r *IBMPowerVSMachine) GetConditions() []metav1.Condition {
	return r.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (r *IBMPowerVSMachine) SetConditions(conditions []metav1.Condition) {
	r.Status.Conditions = conditions
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (r *IBMPowerVSMachine) GetV1Beta1Conditions() clusterv1.Conditions {
	if r.Status.Deprecated == nil || r.Status.Deprecated.V1Beta2 == nil {
		return nil
	}
	return r.Status.Deprecated.V1Beta2.Conditions
}

// SetV1Beta1Conditions sets conditions for an API object.
func (r *IBMPowerVSMachine) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if r.Status.Deprecated == nil {
		r.Status.Deprecated = &IBMPowerVSMachineDeprecatedStatus{}
	}
	if r.Status.Deprecated.V1Beta2 == nil {
		r.Status.Deprecated.V1Beta2 = &IBMPowerVSMachineV1Beta2DeprecatedStatus{}
	}
	r.Status.Deprecated.V1Beta2.Conditions = conditions
}
