// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineSnapshot to the Hub version.
func (src *VirtualMachineSnapshot) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineSnapshot)
	return Convert_v1alpha5_VirtualMachineSnapshot_To_v1alpha6_VirtualMachineSnapshot(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineSnapshot.
func (dst *VirtualMachineSnapshot) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineSnapshot)
	return Convert_v1alpha6_VirtualMachineSnapshot_To_v1alpha5_VirtualMachineSnapshot(src, dst, nil)
}

// ConvertTo converts this VirtualMachineSnapshotList to the Hub version.
func (src *VirtualMachineSnapshotList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineSnapshotList)
	return Convert_v1alpha5_VirtualMachineSnapshotList_To_v1alpha6_VirtualMachineSnapshotList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineSnapshotList.
func (dst *VirtualMachineSnapshotList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineSnapshotList)
	return Convert_v1alpha6_VirtualMachineSnapshotList_To_v1alpha5_VirtualMachineSnapshotList(src, dst, nil)
}
