// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineReplicaSet to the Hub version.
func (src *VirtualMachineReplicaSet) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineReplicaSet)
	return Convert_v1alpha5_VirtualMachineReplicaSet_To_v1alpha6_VirtualMachineReplicaSet(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineReplicaSet.
func (dst *VirtualMachineReplicaSet) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineReplicaSet)
	return Convert_v1alpha6_VirtualMachineReplicaSet_To_v1alpha5_VirtualMachineReplicaSet(src, dst, nil)
}

// ConvertTo converts this VirtualMachineReplicaSetList to the Hub version.
func (src *VirtualMachineReplicaSetList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineReplicaSetList)
	return Convert_v1alpha5_VirtualMachineReplicaSetList_To_v1alpha6_VirtualMachineReplicaSetList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineReplicaSetList.
func (dst *VirtualMachineReplicaSetList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineReplicaSetList)
	return Convert_v1alpha6_VirtualMachineReplicaSetList_To_v1alpha5_VirtualMachineReplicaSetList(src, dst, nil)
}
