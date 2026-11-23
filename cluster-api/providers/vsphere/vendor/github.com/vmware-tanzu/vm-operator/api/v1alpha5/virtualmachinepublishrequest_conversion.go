// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachinePublishRequest to the Hub version.
func (src *VirtualMachinePublishRequest) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachinePublishRequest)
	return Convert_v1alpha5_VirtualMachinePublishRequest_To_v1alpha6_VirtualMachinePublishRequest(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachinePublishRequest.
func (dst *VirtualMachinePublishRequest) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachinePublishRequest)
	return Convert_v1alpha6_VirtualMachinePublishRequest_To_v1alpha5_VirtualMachinePublishRequest(src, dst, nil)
}

// ConvertTo converts this VirtualMachinePublishRequestList to the Hub version.
func (src *VirtualMachinePublishRequestList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachinePublishRequestList)
	return Convert_v1alpha5_VirtualMachinePublishRequestList_To_v1alpha6_VirtualMachinePublishRequestList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachinePublishRequestList.
func (dst *VirtualMachinePublishRequestList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachinePublishRequestList)
	return Convert_v1alpha6_VirtualMachinePublishRequestList_To_v1alpha5_VirtualMachinePublishRequestList(src, dst, nil)
}
