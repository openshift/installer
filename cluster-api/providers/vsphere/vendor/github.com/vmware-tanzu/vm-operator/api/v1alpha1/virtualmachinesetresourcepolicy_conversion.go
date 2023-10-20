// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/vmware-tanzu/vm-operator/api/v1alpha2"
)

func Convert_v1alpha1_VirtualMachineSetResourcePolicySpec_To_v1alpha2_VirtualMachineSetResourcePolicySpec(
	in *VirtualMachineSetResourcePolicySpec, out *v1alpha2.VirtualMachineSetResourcePolicySpec, s apiconversion.Scope) error {

	out.Folder = in.Folder.Name
	for _, mod := range in.ClusterModules {
		out.ClusterModuleGroups = append(out.ClusterModuleGroups, mod.GroupName)
	}

	return autoConvert_v1alpha1_VirtualMachineSetResourcePolicySpec_To_v1alpha2_VirtualMachineSetResourcePolicySpec(in, out, s)
}

func Convert_v1alpha2_VirtualMachineSetResourcePolicySpec_To_v1alpha1_VirtualMachineSetResourcePolicySpec(
	in *v1alpha2.VirtualMachineSetResourcePolicySpec, out *VirtualMachineSetResourcePolicySpec, s apiconversion.Scope) error {

	out.Folder.Name = in.Folder
	for _, name := range in.ClusterModuleGroups {
		out.ClusterModules = append(out.ClusterModules, ClusterModuleSpec{GroupName: name})
	}

	return autoConvert_v1alpha2_VirtualMachineSetResourcePolicySpec_To_v1alpha1_VirtualMachineSetResourcePolicySpec(in, out, s)
}

// ConvertTo converts this VirtualMachineSetResourcePolicy to the Hub version.
func (src *VirtualMachineSetResourcePolicy) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachineSetResourcePolicy)
	return Convert_v1alpha1_VirtualMachineSetResourcePolicy_To_v1alpha2_VirtualMachineSetResourcePolicy(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineSetResourcePolicy.
func (dst *VirtualMachineSetResourcePolicy) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachineSetResourcePolicy)
	return Convert_v1alpha2_VirtualMachineSetResourcePolicy_To_v1alpha1_VirtualMachineSetResourcePolicy(src, dst, nil)
}

// ConvertTo converts this VirtualMachineSetResourcePolicyList to the Hub version.
func (src *VirtualMachineSetResourcePolicyList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachineSetResourcePolicyList)
	return Convert_v1alpha1_VirtualMachineSetResourcePolicyList_To_v1alpha2_VirtualMachineSetResourcePolicyList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineSetResourcePolicyList.
func (dst *VirtualMachineSetResourcePolicyList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachineSetResourcePolicyList)
	return Convert_v1alpha2_VirtualMachineSetResourcePolicyList_To_v1alpha1_VirtualMachineSetResourcePolicyList(src, dst, nil)
}
