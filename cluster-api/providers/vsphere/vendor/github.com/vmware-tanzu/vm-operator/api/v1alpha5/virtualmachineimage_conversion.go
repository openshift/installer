// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineImage to the Hub version.
func (src *VirtualMachineImage) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineImage)
	return Convert_v1alpha5_VirtualMachineImage_To_v1alpha6_VirtualMachineImage(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineImage.
func (dst *VirtualMachineImage) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineImage)
	return Convert_v1alpha6_VirtualMachineImage_To_v1alpha5_VirtualMachineImage(src, dst, nil)
}

// ConvertTo converts this VirtualMachineImageList to the Hub version.
func (src *VirtualMachineImageList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineImageList)
	return Convert_v1alpha5_VirtualMachineImageList_To_v1alpha6_VirtualMachineImageList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineImageList.
func (dst *VirtualMachineImageList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineImageList)
	return Convert_v1alpha6_VirtualMachineImageList_To_v1alpha5_VirtualMachineImageList(src, dst, nil)
}

// ConvertTo converts this ClusterVirtualMachineImage to the Hub version.
func (src *ClusterVirtualMachineImage) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.ClusterVirtualMachineImage)
	return Convert_v1alpha5_ClusterVirtualMachineImage_To_v1alpha6_ClusterVirtualMachineImage(src, dst, nil)
}

// ConvertFrom converts the hub version to this ClusterVirtualMachineImage.
func (dst *ClusterVirtualMachineImage) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.ClusterVirtualMachineImage)
	return Convert_v1alpha6_ClusterVirtualMachineImage_To_v1alpha5_ClusterVirtualMachineImage(src, dst, nil)
}

// ConvertTo converts this ClusterVirtualMachineImageList to the Hub version.
func (src *ClusterVirtualMachineImageList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.ClusterVirtualMachineImageList)
	return Convert_v1alpha5_ClusterVirtualMachineImageList_To_v1alpha6_ClusterVirtualMachineImageList(src, dst, nil)
}

// ConvertFrom converts the hub version to this ClusterVirtualMachineImageList.
func (dst *ClusterVirtualMachineImageList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.ClusterVirtualMachineImageList)
	return Convert_v1alpha6_ClusterVirtualMachineImageList_To_v1alpha5_ClusterVirtualMachineImageList(src, dst, nil)
}
