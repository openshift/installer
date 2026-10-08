// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineClass to the Hub version.
func (src *VirtualMachineClass) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineClass)
	return Convert_v1alpha2_VirtualMachineClass_To_v1alpha6_VirtualMachineClass(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineClass.
func (dst *VirtualMachineClass) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineClass)
	return Convert_v1alpha6_VirtualMachineClass_To_v1alpha2_VirtualMachineClass(src, dst, nil)
}

// ConvertTo converts this VirtualMachineClassList to the Hub version.
func (src *VirtualMachineClassList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineClassList)
	return Convert_v1alpha2_VirtualMachineClassList_To_v1alpha6_VirtualMachineClassList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineClassList.
func (dst *VirtualMachineClassList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineClassList)
	return Convert_v1alpha6_VirtualMachineClassList_To_v1alpha2_VirtualMachineClassList(src, dst, nil)
}
