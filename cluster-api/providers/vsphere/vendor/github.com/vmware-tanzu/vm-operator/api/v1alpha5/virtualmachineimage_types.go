// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vmopv1a5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
)

const (
	imageOSLabelPrefix = "os.image." + GroupName + "/"

	// VirtualMachineImageOSIDLabel is a label applied to images with
	// a value that specifies the operating system ID.
	VirtualMachineImageOSIDLabel = imageOSLabelPrefix + "id"

	// VirtualMachineImageOSTypeLabel is a label applied to images with a
	// value that specifies the operating system type.
	VirtualMachineImageOSTypeLabel = imageOSLabelPrefix + "type"

	// VirtualMachineImageOSVersionLabel is a label applied to images with
	// a value that specifies the operating system version.
	VirtualMachineImageOSVersionLabel = imageOSLabelPrefix + "version"
)

const (
	// VirtualMachineImageCapabilityLabel is the prefix for a label that
	// advertises an image capability.
	VirtualMachineImageCapabilityLabel = "capability.image." + GroupName + "/"
)

const (
	// VirtualMachineImageTypeLabel is the label key for the type of the image
	// (e.g. OVF, ISO, etc.).
	VirtualMachineImageTypeLabel = "image." + GroupName + "/type"
)

const (
	// VMIContentLibRefAnnotation is the key for the annotation that stores the content library
	// reference for VMI and CVMI down conversion.
	VMIContentLibRefAnnotation = "vmoperator.vmware.com/conversion-content-lib-ref"
)

const (
	// VirtualMachineImageV1Alpha1CompatibleCondition denotes that an image was prepared by
	// VMware specifically for compatibility with VMService.
	VirtualMachineImageV1Alpha1CompatibleCondition = "VirtualMachineImageV1Alpha1Compatible"
)

// Condition reasons for VirtualMachineImages.
const (
	// VirtualMachineImageNotSyncedReason documents that the VirtualMachineImage is not synced with
	// the vSphere content library item that contains the source of this image's information.
	VirtualMachineImageNotSyncedReason = "VirtualMachineImageNotSynced"

	// VirtualMachineImageProviderNotReadyReason documents that the VirtualMachineImage provider
	// is not in ready state.
	VirtualMachineImageProviderNotReadyReason = "VirtualMachineImageProviderNotReady"

	// VirtualMachineImageProviderSecurityNotCompliantReason documents that the
	// VirtualMachineImage provider doesn't meet security compliance requirements.
	VirtualMachineImageProviderSecurityNotCompliantReason = "VirtualMachineImageProviderSecurityNotCompliant"
)

// VirtualMachineImageProductInfo describes product information for an image.
type VirtualMachineImageProductInfo struct {
	// +optional

	// Product is a general descriptor for the image.
	Product string `json:"product,omitempty"`

	// +optional

	// Vendor describes the organization/user that produced the image.
	Vendor string `json:"vendor,omitempty"`

	// +optional

	// Version describes the short-form version of the image.
	Version string `json:"version,omitempty"`

	// +optional

	// FullVersion describes the long-form version of the image.
	FullVersion string `json:"fullVersion,omitempty"`
}

// VirtualMachineImageDiskInfo describes information about any disks associated with
// this image.
type VirtualMachineImageDiskInfo struct {

	// +required

	// Name is the unique identifier for the virtual disk.
	Name string `json:"name"`

	// +optional

	// Limit is the total virtual disk capacity.
	Limit *resource.Quantity `json:"limit,omitempty"`

	// +optional

	// Requested is the minimum storage requirements for the virtual disk.
	Requested *resource.Quantity `json:"requested,omitempty"`

	// +optional

	// ControllerType describes the type of controller to which this disk is
	// attached.
	ControllerType VirtualControllerType `json:"controllerType,omitempty"`

	// +optional

	// ControllerBusNumber describes bus number of the controller to which this
	// disk is attached.
	ControllerBusNumber *int32 `json:"controllerBusNumber,omitempty"`

	// +optional

	// UnitNumber describes the unit number for this disk.
	UnitNumber *int32 `json:"unitNumber,omitempty"`
}

// GetControllerType returns the type of controller to which the disk is
// attached.
func (v VirtualMachineImageDiskInfo) GetControllerType() VirtualControllerType {
	return v.ControllerType
}

// GetControllerBusNumber returns the bus number of the controller to which the
// disk is connected.
func (v VirtualMachineImageDiskInfo) GetControllerBusNumber() *int32 {
	return v.ControllerBusNumber
}

// GetUnitNumber returns the disk's unit number.
func (v VirtualMachineImageDiskInfo) GetUnitNumber() *int32 {
	return v.UnitNumber
}

// VirtualMachineImageOSInfo describes the image's guest operating system.
type VirtualMachineImageOSInfo struct {
	// +optional

	// ID describes the operating system ID.
	//
	// This value is also added to the image resource's labels as
	// VirtualMachineImageOSIDLabel.
	ID string `json:"id,omitempty"`

	// +optional

	// Type describes the operating system type.
	//
	// This value is also added to the image resource's labels as
	// VirtualMachineImageOSTypeLabel.
	Type string `json:"type,omitempty"`

	// +optional

	// Version describes the operating system version.
	//
	// This value is also added to the image resource's labels as
	// VirtualMachineImageOSVersionLabel.
	Version string `json:"version,omitempty"`
}

// OVFProperty describes an OVF property associated with an image.
// OVF properties may be used in conjunction with the vAppConfig bootstrap
// provider to customize a VM during its creation.
type OVFProperty struct {
	// Key describes the OVF property's key.
	Key string `json:"key"`

	// Type describes the OVF property's type.
	Type string `json:"type"`

	// +optional

	// Default describes the OVF property's default value.
	Default *string `json:"default,omitempty"`
}

// VirtualMachineImageSpec defines the desired state of VirtualMachineImage.
type VirtualMachineImageSpec struct {
	// +optional

	// ProviderRef is a reference to the resource that contains the source of
	// this image's information.
	ProviderRef *vmopv1a5common.LocalObjectRef `json:"providerRef,omitempty"`
}

// VirtualMachineImageStatus defines the observed state of VirtualMachineImage.
type VirtualMachineImageStatus struct {

	// +optional

	// Name describes the display name of this image.
	Name string `json:"name,omitempty"`

	// +optional
	// +listType=set

	// Capabilities describes the image's observed capabilities.
	//
	// The capabilities are discerned when VM Operator reconciles an image.
	// If the source of an image is an OVF in Content Library, then the
	// capabilities are parsed from the OVF property
	// capabilities.image.vmoperator.vmware.com as a comma-separated list of
	// values. Well-known capabilities include:
	//
	// * cloud-init
	// * nvidia-gpu
	// * sriov-net
	//
	// Every capability is also added to the resource's labels as
	// VirtualMachineImageCapabilityLabel + Value. For example, if the
	// capability is "cloud-init" then the following label will be added to the
	// resource: capability.image.vmoperator.vmware.com/cloud-init.
	Capabilities []string `json:"capabilities,omitempty"`

	// +optional

	// Firmware describe the firmware type used by this image, ex. BIOS, EFI.
	Firmware string `json:"firmware,omitempty"`

	// +optional

	// HardwareVersion describes the observed hardware version of this image.
	HardwareVersion *int32 `json:"hardwareVersion,omitempty"`

	// +optional

	// OSInfo describes the observed operating system information for this
	// image.
	//
	// The OS information is also added to the image resource's labels. Please
	// refer to VirtualMachineImageOSInfo for more information.
	OSInfo VirtualMachineImageOSInfo `json:"osInfo,omitempty"`

	// +optional

	// OVFProperties describes the observed user configurable OVF properties defined for this
	// image.
	OVFProperties []OVFProperty `json:"ovfProperties,omitempty"`

	// +optional

	// VMwareSystemProperties describes the observed VMware system properties defined for
	// this image.
	VMwareSystemProperties []vmopv1a5common.KeyValuePair `json:"vmwareSystemProperties,omitempty"`

	// +optional

	// ProductInfo describes the observed product information for this image.
	ProductInfo VirtualMachineImageProductInfo `json:"productInfo,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=name

	// Disks describes the observed disk information for this image.
	Disks []VirtualMachineImageDiskInfo `json:"disks,omitempty"`

	// +optional

	// ProviderContentVersion describes the content version from the provider item
	// that this image corresponds to. If the provider of this image is a Content
	// Library, this will be the version of the corresponding Content Library item.
	ProviderContentVersion string `json:"providerContentVersion,omitempty"`

	// +optional

	// ProviderItemID describes the ID of the provider item that this image corresponds to.
	// If the provider of this image is a Content Library, this ID will be that of the
	// corresponding Content Library item.
	ProviderItemID string `json:"providerItemID,omitempty"`

	// +optional
	// +listType=map
	// +listMapKey=type

	// Conditions describes the observed conditions for this image.
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// +optional
	//
	// Type describes the content library item type (OVF, ISO, or VM) of the
	// image.
	Type string `json:"type,omitempty"`
}

func (i VirtualMachineImageStatus) GetConditions() []metav1.Condition {
	return i.Conditions
}

func (i *VirtualMachineImageStatus) SetConditions(conditions []metav1.Condition) {
	i.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmi;vmimage
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".status.name"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.type"
// +kubebuilder:printcolumn:name="Image Version",type="string",JSONPath=".status.productInfo.version"
// +kubebuilder:printcolumn:name="OS Name",type="string",JSONPath=".status.osInfo.type"
// +kubebuilder:printcolumn:name="OS Version",type="string",JSONPath=".status.osInfo.version"
// +kubebuilder:printcolumn:name="Hardware Version",type="string",JSONPath=".status.hardwareVersion"
// +kubebuilder:printcolumn:name="Capabilities",type="string",JSONPath=".status.capabilities"

// VirtualMachineImage is the schema for the virtualmachineimages API.
type VirtualMachineImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineImageSpec   `json:"spec,omitempty"`
	Status VirtualMachineImageStatus `json:"status,omitempty"`
}

func (i VirtualMachineImage) GetConditions() []metav1.Condition {
	return i.Status.Conditions
}

func (i *VirtualMachineImage) SetConditions(conditions []metav1.Condition) {
	i.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// VirtualMachineImageList contains a list of VirtualMachineImage.
type VirtualMachineImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineImage `json:"items"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=cvmi;cvmimage;clustervmi;clustervmimage
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".status.name"
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".status.type"
// +kubebuilder:printcolumn:name="Image Version",type="string",JSONPath=".status.productInfo.version"
// +kubebuilder:printcolumn:name="OS Name",type="string",JSONPath=".status.osInfo.type"
// +kubebuilder:printcolumn:name="OS Version",type="string",JSONPath=".status.osInfo.version"
// +kubebuilder:printcolumn:name="Hardware Version",type="string",JSONPath=".status.hardwareVersion"

// ClusterVirtualMachineImage is the schema for the clustervirtualmachineimages
// API.
type ClusterVirtualMachineImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineImageSpec   `json:"spec,omitempty"`
	Status VirtualMachineImageStatus `json:"status,omitempty"`
}

func (i ClusterVirtualMachineImage) GetConditions() []metav1.Condition {
	return i.Status.Conditions
}

func (i *ClusterVirtualMachineImage) SetConditions(conditions []metav1.Condition) {
	i.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// ClusterVirtualMachineImageList contains a list of ClusterVirtualMachineImage.
type ClusterVirtualMachineImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterVirtualMachineImage `json:"items"`
}

func init() {
	objectTypes = append(objectTypes,
		&VirtualMachineImage{},
		&VirtualMachineImageList{},
		&ClusterVirtualMachineImage{},
		&ClusterVirtualMachineImageList{})
}
