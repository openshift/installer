// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachineGroupPublishRequest to the Hub version.
func (src *VirtualMachineGroupPublishRequest) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineGroupPublishRequest)
	return Convert_v1alpha5_VirtualMachineGroupPublishRequest_To_v1alpha6_VirtualMachineGroupPublishRequest(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineGroupPublishRequest.
func (dst *VirtualMachineGroupPublishRequest) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineGroupPublishRequest)
	return Convert_v1alpha6_VirtualMachineGroupPublishRequest_To_v1alpha5_VirtualMachineGroupPublishRequest(src, dst, nil)
}

// ConvertTo converts this VirtualMachineGroupPublishRequestList to the Hub version.
func (src *VirtualMachineGroupPublishRequestList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineGroupPublishRequestList)
	return Convert_v1alpha5_VirtualMachineGroupPublishRequestList_To_v1alpha6_VirtualMachineGroupPublishRequestList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineGroupPublishRequestList.
func (dst *VirtualMachineGroupPublishRequestList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineGroupPublishRequestList)
	return Convert_v1alpha6_VirtualMachineGroupPublishRequestList_To_v1alpha5_VirtualMachineGroupPublishRequestList(src, dst, nil)
}
