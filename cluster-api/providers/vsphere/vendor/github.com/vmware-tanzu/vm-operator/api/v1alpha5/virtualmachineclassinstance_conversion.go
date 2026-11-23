// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineClassInstance to the Hub version.
func (src *VirtualMachineClassInstance) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineClassInstance)
	return Convert_v1alpha5_VirtualMachineClassInstance_To_v1alpha6_VirtualMachineClassInstance(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineClassInstance.
func (dst *VirtualMachineClassInstance) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineClassInstance)
	return Convert_v1alpha6_VirtualMachineClassInstance_To_v1alpha5_VirtualMachineClassInstance(src, dst, nil)
}

// ConvertTo converts this VirtualMachineClassInstanceList to the Hub version.
func (src *VirtualMachineClassInstanceList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineClassInstanceList)
	return Convert_v1alpha5_VirtualMachineClassInstanceList_To_v1alpha6_VirtualMachineClassInstanceList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineClassInstanceList.
func (dst *VirtualMachineClassInstanceList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineClassInstanceList)
	return Convert_v1alpha6_VirtualMachineClassInstanceList_To_v1alpha5_VirtualMachineClassInstanceList(src, dst, nil)
}
