// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
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
	// Product is a general descriptor for the image.
	// +optional
	Product string `json:"product,omitempty"`

	// Vendor describes the organization/user that produced the image.
	// +optional
	Vendor string `json:"vendor,omitempty"`

	// Version describes the short-form version of the image.
	// +optional
	Version string `json:"version,omitempty"`

	// FullVersion describes the long-form version of the image.
	// +optional
	FullVersion string `json:"fullVersion,omitempty"`
}

// VirtualMachineImageOSInfo describes the image's guest operating system.
type VirtualMachineImageOSInfo struct {
	// ID describes the operating system ID.
	//
	// This value is also added to the image resource's labels as
	// VirtualMachineImageOSIDLabel.
	//
	// +optional
	ID string `json:"id,omitempty"`

	// Type describes the operating system type.
	//
	// This value is also added to the image resource's labels as
	// VirtualMachineImageOSTypeLabel.
	//
	// +optional
	Type string `json:"type,omitempty"`

	// Version describes the operating system version.
	//
	// This value is also added to the image resource's labels as
	// VirtualMachineImageOSVersionLabel.
	//
	// +optional
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

	// Default describes the OVF property's default value.
	// +optional
	Default *string `json:"default,omitempty"`
}

// VirtualMachineImageSpec defines the desired state of VirtualMachineImage.
type VirtualMachineImageSpec struct {
	// ProviderRef is a reference to the resource that contains the source of
	// this image's information.
	//
	// +optional
	ProviderRef common.LocalObjectRef `json:"providerRef,omitempty"`
}

// VirtualMachineImageStatus defines the observed state of VirtualMachineImage.
type VirtualMachineImageStatus struct {

	// Name describes the display name of this image.
	//
	// +optional
	Name string `json:"name,omitempty"`

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
	//
	// +optional
	// +listType=set
	Capabilities []string `json:"capabilities,omitempty"`

	// Firmware describe the firmware type used by this image, ex. BIOS, EFI.
	// +optional
	Firmware string `json:"firmware,omitempty"`

	// HardwareVersion describes the observed hardware version of this image.
	//
	// +optional
	HardwareVersion *int32 `json:"hardwareVersion,omitempty"`

	// OSInfo describes the observed operating system information for this
	// image.
	//
	// The OS information is also added to the image resource's labels. Please
	// refer to VirtualMachineImageOSInfo for more information.
	//
	//
	// +optional
	OSInfo VirtualMachineImageOSInfo `json:"osInfo,omitempty"`

	// OVFProperties describes the observed user configurable OVF properties defined for this
	// image.
	//
	// +optional
	OVFProperties []OVFProperty `json:"ovfProperties,omitempty"`

	// VMwareSystemProperties describes the observed VMware system properties defined for
	// this image.
	//
	// +optional
	VMwareSystemProperties []common.KeyValuePair `json:"vmwareSystemProperties,omitempty"`

	// ProductInfo describes the observed product information for this image.
	// +optional
	ProductInfo VirtualMachineImageProductInfo `json:"productInfo,omitempty"`

	// ProviderContentVersion describes the content version from the provider item
	// that this image corresponds to. If the provider of this image is a Content
	// Library, this will be the version of the corresponding Content Library item.
	// +optional
	ProviderContentVersion string `json:"providerContentVersion,omitempty"`

	// ProviderItemID describes the ID of the provider item that this image corresponds to.
	// If the provider of this image is a Content Library, this ID will be that of the
	// corresponding Content Library item.
	// +optional
	ProviderItemID string `json:"providerItemID,omitempty"`

	// Conditions describes the observed conditions for this image.
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=vmi;vmimage
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".status.name"
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

func (i *VirtualMachineImage) GetConditions() []metav1.Condition {
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
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display Name",type="string",JSONPath=".status.name"
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

func (i *ClusterVirtualMachineImage) GetConditions() []metav1.Condition {
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
	SchemeBuilder.Register(
		&VirtualMachineImage{},
		&VirtualMachineImageList{},
		&ClusterVirtualMachineImage{},
		&ClusterVirtualMachineImageList{},
	)
}
