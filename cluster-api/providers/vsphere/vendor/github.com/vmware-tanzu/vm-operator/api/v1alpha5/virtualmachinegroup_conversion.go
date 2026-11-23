// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineGroup to the Hub version.
func (src *VirtualMachineGroup) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineGroup)
	return Convert_v1alpha5_VirtualMachineGroup_To_v1alpha6_VirtualMachineGroup(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineGroup.
func (dst *VirtualMachineGroup) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineGroup)
	return Convert_v1alpha6_VirtualMachineGroup_To_v1alpha5_VirtualMachineGroup(src, dst, nil)
}

// ConvertTo converts this VirtualMachineGroupList to the Hub version.
func (src *VirtualMachineGroupList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineGroupList)
	return Convert_v1alpha5_VirtualMachineGroupList_To_v1alpha6_VirtualMachineGroupList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineGroupList.
func (dst *VirtualMachineGroupList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineGroupList)
	return Convert_v1alpha6_VirtualMachineGroupList_To_v1alpha5_VirtualMachineGroupList(src, dst, nil)
}
