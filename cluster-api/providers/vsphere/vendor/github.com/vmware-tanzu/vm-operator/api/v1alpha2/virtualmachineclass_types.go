// Copyright (c) 2022 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	vmClassLabel = "class." + GroupName + "/"

	// VirtualMachineClassCapabilityLabel is the prefix for a label that
	// advertises a VM class capability.
	//
	// For more information please see VirtualMachineClassStatus.Capabilities.
	VirtualMachineClassCapabilityLabel = "capability." + vmClassLabel

	// VirtualMachineClassReadyLabel is applied to VM class resources that can
	// be realized in the current cluster.
	//
	// For more information please see VirtualMachineClassStatus.Ready.
	VirtualMachineClassReadyLabel = vmClassLabel + "ready"
)

const (
	// VirtualMachineClassConditionReady is a condition on the VM class that
	// indicates whether or not the VM class's hardware can be realized on at
	// least one host in the cluster.
	//
	// For example, on Supervisor the VM Class / host compatibility is
	// performed by checking with the CheckVmConfig_Task API
	// (https://bit.ly/3CvoCc8).
	//
	// It is possible for this condition's Status to flip between True and False
	// over the lifetime of a VM class. For example, if a VM class requires a
	// specific vGPU and there is no host that provides that vGPU then the
	// Status for this condition is False. However, if a host is added with that
	// vGPU then the condition's status changes to True. And vice-versa -- if
	// the same host is removed then the condition's status flips back to False.
	//
	// Please note that a host in maintenance mode that would otherwise provide
	// compatibility will not be considered for this check. Only hosts where VMs
	// can be scheduled are part of the check.
	//
	VirtualMachineClassConditionReady = "VirtualMachineClassReady"
)

// VirtualMachineConfigSpec contains additional virtual machine
// configuration settings including hardware specifications for the
// VirtualMachine.
//
// We only support XML for now, but that may change in the future.
type VirtualMachineConfigSpec struct {
	// XML contains a vim.vm.ConfigSpec object that has been serialized to XML
	// and base64-encoded.
	//
	// +optional
	XML string `json:"xml,omitempty"`
}

// VGPUDevice contains the configuration corresponding to a vGPU device.
type VGPUDevice struct {
	ProfileName string `json:"profileName"`
}

// DynamicDirectPathIODevice contains the configuration corresponding to a
// Dynamic DirectPath I/O device.
type DynamicDirectPathIODevice struct {
	VendorID int64 `json:"vendorID"`
	DeviceID int64 `json:"deviceID"`

	// +optional
	CustomLabel string `json:"customLabel,omitempty"`
}

// InstanceStorage provides information used to configure instance
// storage volumes for a VirtualMachine.
type InstanceStorage struct {
	// StorageClass refers to the name of a StorageClass resource used to
	// provide the storage for the configured instance storage volumes.
	// The value of this field has no relationship to or bearing on the field
	// virtualMachine.spec.storageClass. Please note the referred StorageClass
	// must be available in the same namespace as the VirtualMachineClass that
	// uses it for configuring instance storage.
	StorageClass string `json:"storageClass,omitempty"`

	// Volumes describes instance storage volumes created for a VirtualMachine
	// instance that use this VirtualMachineClass.
	Volumes []InstanceStorageVolume `json:"volumes,omitempty"`
}

// InstanceStorageVolume contains information required to create an
// instance storage volume on a VirtualMachine.
type InstanceStorageVolume struct {
	Size resource.Quantity `json:"size"`
}

// VirtualDevices contains information about the virtual devices associated
// with a VirtualMachineClass.
type VirtualDevices struct {
	// +optional
	// +listType=map
	// +listMapKey=profileName
	VGPUDevices []VGPUDevice `json:"vgpuDevices,omitempty"`

	// +optional
	DynamicDirectPathIODevices []DynamicDirectPathIODevice `json:"dynamicDirectPathIODevices,omitempty"`
}

// VirtualMachineClassHardware describes a virtual hardware resource
// specification.
type VirtualMachineClassHardware struct {
	// +optional
	Cpus int64 `json:"cpus,omitempty"`

	// +optional
	Memory resource.Quantity `json:"memory,omitempty"`

	// +optional
	Devices VirtualDevices `json:"devices,omitempty"`

	// +optional
	InstanceStorage InstanceStorage `json:"instanceStorage,omitempty"`
}

// VirtualMachineResourceSpec describes a virtual hardware policy specification.
type VirtualMachineResourceSpec struct {
	// +optional
	Cpu resource.Quantity `json:"cpu,omitempty"` //nolint:stylecheck,revive

	// +optional
	Memory resource.Quantity `json:"memory,omitempty"`
}

// VirtualMachineClassResources describes the virtual hardware resource
// reservations and limits configuration to be used by a VirtualMachineClass.
type VirtualMachineClassResources struct {
	// +optional
	Requests VirtualMachineResourceSpec `json:"requests,omitempty"`

	// +optional
	Limits VirtualMachineResourceSpec `json:"limits,omitempty"`
}

// VirtualMachineClassPolicies describes the policy configuration to be used by
// a VirtualMachineClass.
type VirtualMachineClassPolicies struct {
	Resources VirtualMachineClassResources `json:"resources,omitempty"`
}

// VirtualMachineClassSpec defines the desired state of VirtualMachineClass.
type VirtualMachineClassSpec struct {
	// ControllerName describes the name of the controller responsible for
	// reconciling VirtualMachine resources that are realized from this
	// VirtualMachineClass.
	//
	// When omitted, controllers reconciling VirtualMachine resources determine
	// the default controller name from the environment variable
	// DEFAULT_VM_CLASS_CONTROLLER_NAME. If this environment variable is not
	// defined or empty, it defaults to vmoperator.vmware.com/vsphere.
	//
	// Once a non-empty value is assigned to this field, attempts to set this
	// field to an empty value will be silently ignored.
	//
	// +optional
	ControllerName string `json:"controllerName,omitempty"`

	// Hardware describes the configuration of the VirtualMachineClass
	// attributes related to virtual hardware. The configuration specified in
	// this field is used to customize the virtual hardware characteristics of
	// any VirtualMachine associated with this VirtualMachineClass.
	//
	// +optional
	Hardware VirtualMachineClassHardware `json:"hardware,omitempty"`

	// Policies describes the configuration of the VirtualMachineClass
	// attributes related to virtual infrastructure policy. The configuration
	// specified in this field is used to customize various policies related to
	// infrastructure resource consumption.
	//
	// +optional
	Policies VirtualMachineClassPolicies `json:"policies,omitempty"`

	// Description describes the configuration of the VirtualMachineClass which
	// is not related to virtual hardware or infrastructure policy. This field
	// is used to address remaining specs about this VirtualMachineClass.
	//
	// +optional
	Description string `json:"description,omitempty"`

	// ConfigSpec describes additional configuration information for a
	// VirtualMachine.
	// The contents of this field are the VirtualMachineConfigSpec data object
	// (https://bit.ly/3HDtiRu) marshaled to JSON using the discriminator
	// field "_typeName" to preserve type information.
	//
	// +optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	// +kubebuilder:pruning:PreserveUnknownFields
	ConfigSpec json.RawMessage `json:"configSpec,omitempty"`
}

// VirtualMachineClassStatus defines the observed state of VirtualMachineClass.
type VirtualMachineClassStatus struct {
	// Capabilities describes the class's observed capabilities.
	//
	// The capabilities are discerned when VM Operator reconciles a class
	// and inspects its specification. Well-known capabilities include:
	//
	// * instance-storage
	// * nvidia-gpu
	// * sriov-net
	//
	// In addition to "nvidia-gpu", a capability is added for every nVidia
	// profile name associated with the class.
	//
	// Every capability is also added to the resource's labels as
	// VirtualMachineClassCapabilityLabel + Value. For example, if the
	// capability is "nvidia-gpu" then the following label will be added to the
	// resource: capability.class.vmoperator.vmware.com/nvidia-gpu.
	//
	// +optional
	// +listType=set
	Capabilities []string `json:"capabilities,omitempty"`

	// Conditions describes the observed conditions of the VirtualMachineClass.
	//
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Ready describes whether the class's hardware can be realized in the
	// cluster.
	//
	// This field is only set to true if all of the class resource's conditions
	// have Status=True.
	//
	// +optional
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=vmclass
// +kubebuilder:storageversion:false
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="CPU",type="string",JSONPath=".spec.hardware.cpus"
// +kubebuilder:printcolumn:name="Memory",type="string",JSONPath=".spec.hardware.memory"
// +kubebuilder:printcolumn:name="Capabilities",type="string",priority=1,JSONPath=".status.capabilities"

// VirtualMachineClass is the schema for the virtualmachineclasses API and
// represents the desired state and observed status of a virtualmachineclasses
// resource.
type VirtualMachineClass struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineClassSpec   `json:"spec,omitempty"`
	Status VirtualMachineClassStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineClassList contains a list of VirtualMachineClass.
type VirtualMachineClassList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineClass `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&VirtualMachineClass{},
		&VirtualMachineClassList{},
	)
}
