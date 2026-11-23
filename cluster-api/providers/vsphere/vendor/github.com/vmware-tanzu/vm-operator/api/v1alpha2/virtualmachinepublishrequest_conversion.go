// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

import (
	"k8s.io/apimachinery/pkg/conversion"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/vmware-tanzu/vm-operator/api/utilconversion"
	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

// ConvertTo converts this VirtualMachinePublishRequest to the Hub version.
func (src *VirtualMachinePublishRequest) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachinePublishRequest)
	if err := Convert_v1alpha2_VirtualMachinePublishRequest_To_v1alpha6_VirtualMachinePublishRequest(src, dst, nil); err != nil {
		return err
	}

	restored := &vmopv1.VirtualMachinePublishRequest{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.BackoffLimit = restored.Spec.BackoffLimit

	return nil
}

// ConvertFrom converts the hub version to this VirtualMachinePublishRequest.
func (dst *VirtualMachinePublishRequest) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachinePublishRequest)
	if err := Convert_v1alpha6_VirtualMachinePublishRequest_To_v1alpha2_VirtualMachinePublishRequest(src, dst, nil); err != nil {
		return err
	}

	return utilconversion.MarshalData(src, dst)
}

// ConvertTo converts this VirtualMachinePublishRequestList to the Hub version.
func (src *VirtualMachinePublishRequestList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachinePublishRequestList)
	return Convert_v1alpha2_VirtualMachinePublishRequestList_To_v1alpha6_VirtualMachinePublishRequestList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachinePublishRequestList.
func (dst *VirtualMachinePublishRequestList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachinePublishRequestList)
	return Convert_v1alpha6_VirtualMachinePublishRequestList_To_v1alpha2_VirtualMachinePublishRequestList(src, dst, nil)
}

func Convert_v1alpha6_VirtualMachinePublishRequestSpec_To_v1alpha2_VirtualMachinePublishRequestSpec(
	in *vmopv1.VirtualMachinePublishRequestSpec, out *VirtualMachinePublishRequestSpec, s conversion.Scope) error {

	if err := autoConvert_v1alpha6_VirtualMachinePublishRequestSpec_To_v1alpha2_VirtualMachinePublishRequestSpec(in, out, s); err != nil {
		return err
	}

	return nil
}
