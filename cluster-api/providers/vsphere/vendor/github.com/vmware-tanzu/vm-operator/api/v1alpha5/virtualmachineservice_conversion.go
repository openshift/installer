// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/vmware-tanzu/vm-operator/api/utilconversion"
	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

func Convert_v1alpha6_VirtualMachineServiceSpec_To_v1alpha5_VirtualMachineServiceSpec(
	in *vmopv1.VirtualMachineServiceSpec, out *VirtualMachineServiceSpec, s apiconversion.Scope) error {
	return autoConvert_v1alpha6_VirtualMachineServiceSpec_To_v1alpha5_VirtualMachineServiceSpec(in, out, s)
}

func restoreV1alpha6VirtualMachineServiceIPFamilies(dst, restored *vmopv1.VirtualMachineService) {
	dst.Spec.IPFamilies = restored.Spec.IPFamilies
	dst.Spec.IPFamilyPolicy = restored.Spec.IPFamilyPolicy
}

// ConvertTo converts this VirtualMachineService to the Hub version.
func (src *VirtualMachineService) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineService)
	if err := Convert_v1alpha5_VirtualMachineService_To_v1alpha6_VirtualMachineService(src, dst, nil); err != nil {
		return err
	}

	restored := &vmopv1.VirtualMachineService{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	restoreV1alpha6VirtualMachineServiceIPFamilies(dst, restored)
	dst.Status = restored.Status
	return nil
}

// ConvertFrom converts the hub version to this VirtualMachineService.
func (dst *VirtualMachineService) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineService)
	if err := Convert_v1alpha6_VirtualMachineService_To_v1alpha5_VirtualMachineService(src, dst, nil); err != nil {
		return err
	}
	return utilconversion.MarshalData(src, dst)
}

// ConvertTo converts this VirtualMachineServiceList to the Hub version.
func (src *VirtualMachineServiceList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineServiceList)
	return Convert_v1alpha5_VirtualMachineServiceList_To_v1alpha6_VirtualMachineServiceList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineServiceList.
func (dst *VirtualMachineServiceList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineServiceList)
	return Convert_v1alpha6_VirtualMachineServiceList_To_v1alpha5_VirtualMachineServiceList(src, dst, nil)
}
