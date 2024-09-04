// Copyright (c) 2018-2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtualMachineImageProductInfo describes optional product-related information that can be added to an image
// template.  This information can be used by the image author to communicate details of the product contained in the
// image.
type VirtualMachineImageProductInfo struct {
	// Product typically describes the type of product contained in the image.
	// +optional
	Product string `json:"product,omitempty"`

	// Vendor typically describes the name of the vendor that is producing the image.
	// +optional
	Vendor string `json:"vendor,omitempty"`

	// Version typically describes a short-form version of the image.
	// +optional
	Version string `json:"version,omitempty"`

	// FullVersion typically describes a long-form version of the image.
	// +optional
	FullVersion string `json:"fullVersion,omitempty"`
}

// VirtualMachineImageOSInfo describes optional information related to the image operating system that can be added
// to an image template. This information can be used by the image author to communicate details of the operating
// system associated with the image.
type VirtualMachineImageOSInfo struct {
	// Version typically describes the version of the guest operating system.
	// +optional
	Version string `json:"version,omitempty"`

	// Type typically describes the type of the guest operating system.
	// +optional
	Type string `json:"type,omitempty"`
}

// OvfProperty describes information related to a user configurable property element that is supported by
// VirtualMachineImage and can be customized during VirtualMachine creation.
type OvfProperty struct {
	// Key describes the key of the ovf property.
	Key string `json:"key"`

	// Type describes the type of the ovf property.
	Type string `json:"type"`

	// Default describes the default value of the ovf key.
	// +optional
	Default *string `json:"default,omitempty"`

	// Description contains the value of the OVF property's optional
	// "Description" element.
	//
	// +optional
	Description string `json:"description,omitempty"`

	// Label contains the value of the OVF property's optional
	// "Label" element.
	//
	// +optional
	Label string `json:"label,omitempty"`
}

// VirtualMachineImageSpec defines the desired state of VirtualMachineImage.
type VirtualMachineImageSpec struct {
	// Type describes the type of the VirtualMachineImage. Currently, the only supported image is "OVF"
	Type string `json:"type"`

	// ImageSourceType describes the type of content source of the VirtualMachineImage.  The only Content Source
	// supported currently is the vSphere Content Library.
	// +optional
	ImageSourceType string `json:"imageSourceType,omitempty"`

	// ImageID is a unique identifier exposed by the provider of this VirtualMachineImage.
	ImageID string `json:"imageID"`

	// ProviderRef is a reference to a content provider object that describes a provider.
	ProviderRef ContentProviderReference `json:"providerRef"`

	// ProductInfo describes the attributes of the VirtualMachineImage relating to the product contained in the
	// image.
	// +optional
	ProductInfo VirtualMachineImageProductInfo `json:"productInfo,omitempty"`

	// OSInfo describes the attributes of the VirtualMachineImage relating to the Operating System contained in the
	// image.
	// +optional
	OSInfo VirtualMachineImageOSInfo `json:"osInfo,omitempty"`

	// OVFEnv describes the user configurable customization parameters of the VirtualMachineImage.
	// +optional
	OVFEnv map[string]OvfProperty `json:"ovfEnv,omitempty"`

	// HardwareVersion describes the virtual hardware version of the image
	// +optional
	HardwareVersion int32 `json:"hwVersion,omitempty"`
}

// VirtualMachineImageStatus defines the observed state of VirtualMachineImage.
type VirtualMachineImageStatus struct {
	// Deprecated
	Uuid string `json:"uuid,omitempty"` //nolint:revive,stylecheck

	// Deprecated
	InternalId string `json:"internalId,omitempty"` //nolint:revive,stylecheck

	// Deprecated
	PowerState string `json:"powerState,omitempty"`

	// ImageName describes the display name of this image.
	// +optional
	ImageName string `json:"imageName,omitempty"`

	// ImageSupported indicates whether the VirtualMachineImage is supported by VMService.
	// A VirtualMachineImage is supported by VMService if the following conditions are true:
	// - VirtualMachineImageV1Alpha1CompatibleCondition
	// +optional
	ImageSupported *bool `json:"imageSupported,omitempty"`

	// Conditions describes the current condition information of the VirtualMachineImage object. e.g. if the OS type
	// is supported or image is supported by VMService
	// +optional
	Conditions []Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// ContentLibraryRef is a reference to the source ContentLibrary/ClusterContentLibrary resource.
	//
	// Deprecated: This field is provider specific but the VirtualMachineImage types are intended to be provider generic.
	// This field does not exist in later API versions. Instead, the Spec.ProviderRef field should be used to look up the
	// provider. For images provided by a Content Library, the ProviderRef will point to either a ContentLibraryItem or
	// ClusterContentLibraryItem that contains a reference to the Content Library.
	// +optional
	ContentLibraryRef *corev1.TypedLocalObjectReference `json:"contentLibraryRef,omitempty"`

	// ContentVersion describes the observed content version of this VirtualMachineImage that was last successfully
	// synced with the vSphere content library item.
	// +optional
	ContentVersion string `json:"contentVersion,omitempty"`

	// Firmware describe the firmware type used by this VirtualMachineImage.
	// eg: bios, efi.
	// +optional
	Firmware string `json:"firmware,omitempty"`
}

func (vmImage *VirtualMachineImage) GetConditions() Conditions {
	return vmImage.Status.Conditions
}

func (vmImage *VirtualMachineImage) SetConditions(conditions Conditions) {
	vmImage.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=vmi;vmimage
// +kubebuilder:storageversion:false
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display-Name",type="string",JSONPath=".status.imageName"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.productInfo.version"
// +kubebuilder:printcolumn:name="Os-Type",type="string",JSONPath=".spec.osInfo.type"
// +kubebuilder:printcolumn:name="Format",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Image-Supported",type="boolean",priority=1,JSONPath=".status.imageSupported"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// VirtualMachineImage is the Schema for the virtualmachineimages API
// A VirtualMachineImage represents a VirtualMachine image (e.g. VM template) that can be used as the base image
// for creating a VirtualMachine instance.  The VirtualMachineImage is a required field of the VirtualMachine
// spec.  Currently, VirtualMachineImages are immutable to end users.
type VirtualMachineImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineImageSpec   `json:"spec,omitempty"`
	Status VirtualMachineImageStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VirtualMachineImageList contains a list of VirtualMachineImage.
type VirtualMachineImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtualMachineImage `json:"items"`
}

func (clusterVirtualMachineImage *ClusterVirtualMachineImage) GetConditions() Conditions {
	return clusterVirtualMachineImage.Status.Conditions
}

func (clusterVirtualMachineImage *ClusterVirtualMachineImage) SetConditions(conditions Conditions) {
	clusterVirtualMachineImage.Status.Conditions = conditions
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=cvmi;cvmimage;clustervmi;clustervmimage
// +kubebuilder:storageversion:false
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Display-Name",type="string",JSONPath=".status.imageName"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.productInfo.version"
// +kubebuilder:printcolumn:name="Os-Type",type="string",JSONPath=".spec.osInfo.type"
// +kubebuilder:printcolumn:name="Format",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Image-Supported",type="boolean",priority=1,JSONPath=".status.imageSupported"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// ClusterVirtualMachineImage is the schema for the clustervirtualmachineimage API
// A ClusterVirtualMachineImage represents the desired specification and the observed status of a
// ClusterVirtualMachineImage instance.
type ClusterVirtualMachineImage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualMachineImageSpec   `json:"spec,omitempty"`
	Status VirtualMachineImageStatus `json:"status,omitempty"`
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
		&ClusterVirtualMachineImageList{})
}
