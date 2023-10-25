// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this VirtualMachinePublishRequest to the Hub version.
func (src *VirtualMachinePublishRequest) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachinePublishRequest)
	return Convert_v1alpha1_VirtualMachinePublishRequest_To_v1alpha2_VirtualMachinePublishRequest(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachinePublishRequest.
func (dst *VirtualMachinePublishRequest) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachinePublishRequest)
	return Convert_v1alpha2_VirtualMachinePublishRequest_To_v1alpha1_VirtualMachinePublishRequest(src, dst, nil)
}

// ConvertTo converts this VirtualMachinePublishRequestList to the Hub version.
func (src *VirtualMachinePublishRequestList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachinePublishRequestList)
	return Convert_v1alpha1_VirtualMachinePublishRequestList_To_v1alpha2_VirtualMachinePublishRequestList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachinePublishRequestList.
func (dst *VirtualMachinePublishRequestList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachinePublishRequestList)
	return Convert_v1alpha2_VirtualMachinePublishRequestList_To_v1alpha1_VirtualMachinePublishRequestList(src, dst, nil)
}
