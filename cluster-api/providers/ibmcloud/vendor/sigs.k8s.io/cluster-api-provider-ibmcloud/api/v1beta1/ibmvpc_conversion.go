/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	utilconversion "sigs.k8s.io/cluster-api/util/conversion"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

func (src *IBMVPCCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta2.IBMVPCCluster)

	return Convert_v1beta1_IBMVPCCluster_To_v1beta2_IBMVPCCluster(src, dst, nil)
}

func (dst *IBMVPCCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta2.IBMVPCCluster)

	return Convert_v1beta2_IBMVPCCluster_To_v1beta1_IBMVPCCluster(src, dst, nil)
}

func (src *IBMVPCClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta2.IBMVPCClusterList)

	return Convert_v1beta1_IBMVPCClusterList_To_v1beta2_IBMVPCClusterList(src, dst, nil)
}

func (dst *IBMVPCClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta2.IBMVPCClusterList)

	return Convert_v1beta2_IBMVPCClusterList_To_v1beta1_IBMVPCClusterList(src, dst, nil)
}

func (src *IBMVPCMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta2.IBMVPCMachine)

	if err := Convert_v1beta1_IBMVPCMachine_To_v1beta2_IBMVPCMachine(src, dst, nil); err != nil {
		return err
	}

	if src.Spec.Image != "" {
		dst.Spec.Image = &infrav1beta2.IBMVPCResourceReference{
			ID: &src.Spec.Image,
		}
	}

	if src.Spec.ImageName != "" {
		dst.Spec.Image = &infrav1beta2.IBMVPCResourceReference{
			Name: &src.Spec.ImageName,
		}
	}

	for _, sshKey := range src.Spec.SSHKeyNames {
		dst.Spec.SSHKeys = append(dst.Spec.SSHKeys, &infrav1beta2.IBMVPCResourceReference{
			Name: sshKey,
		})
	}

	return nil
}

func (dst *IBMVPCMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta2.IBMVPCMachine)

	if err := Convert_v1beta2_IBMVPCMachine_To_v1beta1_IBMVPCMachine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	if src.Spec.Image != nil && src.Spec.Image.ID != nil {
		dst.Spec.Image = *src.Spec.Image.ID
	}

	if src.Spec.Image != nil && src.Spec.Image.Name != nil {
		dst.Spec.ImageName = *src.Spec.Image.Name
	}

	for _, sshKey := range src.Spec.SSHKeys {
		if sshKey.Name != nil {
			dst.Spec.SSHKeyNames = append(dst.Spec.SSHKeyNames, sshKey.Name)
		}
	}

	return nil
}

func (src *IBMVPCMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta2.IBMVPCMachineList)

	return Convert_v1beta1_IBMVPCMachineList_To_v1beta2_IBMVPCMachineList(src, dst, nil)
}

func (dst *IBMVPCMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta2.IBMVPCMachineList)

	return Convert_v1beta2_IBMVPCMachineList_To_v1beta1_IBMVPCMachineList(src, dst, nil)
}

func (src *IBMVPCMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta2.IBMVPCMachineTemplate)

	if err := Convert_v1beta1_IBMVPCMachineTemplate_To_v1beta2_IBMVPCMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	if src.Spec.Template.Spec.Image != "" {
		dst.Spec.Template.Spec.Image = &infrav1beta2.IBMVPCResourceReference{
			ID: &src.Spec.Template.Spec.Image,
		}
	}

	if src.Spec.Template.Spec.ImageName != "" {
		dst.Spec.Template.Spec.Image = &infrav1beta2.IBMVPCResourceReference{
			Name: &src.Spec.Template.Spec.ImageName,
		}
	}

	for _, sshKey := range src.Spec.Template.Spec.SSHKeyNames {
		dst.Spec.Template.Spec.SSHKeys = append(dst.Spec.Template.Spec.SSHKeys, &infrav1beta2.IBMVPCResourceReference{
			Name: sshKey,
		})
	}

	return nil
}

func (dst *IBMVPCMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta2.IBMVPCMachineTemplate)

	if err := Convert_v1beta2_IBMVPCMachineTemplate_To_v1beta1_IBMVPCMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	if src.Spec.Template.Spec.Image != nil && src.Spec.Template.Spec.Image.ID != nil {
		dst.Spec.Template.Spec.Image = *src.Spec.Template.Spec.Image.ID
	}

	if src.Spec.Template.Spec.Image != nil && src.Spec.Template.Spec.Image.Name != nil {
		dst.Spec.Template.Spec.ImageName = *src.Spec.Template.Spec.Image.Name
	}

	for _, sshKey := range src.Spec.Template.Spec.SSHKeys {
		if sshKey.Name != nil {
			dst.Spec.Template.Spec.SSHKeyNames = append(dst.Spec.Template.Spec.SSHKeyNames, sshKey.Name)
		}
	}

	return nil
}

func (src *IBMVPCMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta2.IBMVPCMachineTemplateList)

	return Convert_v1beta1_IBMVPCMachineTemplateList_To_v1beta2_IBMVPCMachineTemplateList(src, dst, nil)
}

func (dst *IBMVPCMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta2.IBMVPCMachineTemplateList)

	return Convert_v1beta2_IBMVPCMachineTemplateList_To_v1beta1_IBMVPCMachineTemplateList(src, dst, nil)
}

func Convert_v1beta1_IBMVPCMachineSpec_To_v1beta2_IBMVPCMachineSpec(in *IBMVPCMachineSpec, out *infrav1beta2.IBMVPCMachineSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_IBMVPCMachineSpec_To_v1beta2_IBMVPCMachineSpec(in, out, s)
}

func Convert_v1beta2_IBMVPCMachineSpec_To_v1beta1_IBMVPCMachineSpec(in *infrav1beta2.IBMVPCMachineSpec, out *IBMVPCMachineSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_IBMVPCMachineSpec_To_v1beta1_IBMVPCMachineSpec(in, out, s)
}

func Convert_v1beta2_IBMVPCMachineTemplateStatus_To_v1beta1_IBMVPCMachineTemplateStatus(in *infrav1beta2.IBMVPCMachineTemplateStatus, out *IBMVPCMachineTemplateStatus, s apiconversion.Scope) error {
	return autoConvert_v1beta2_IBMVPCMachineTemplateStatus_To_v1beta1_IBMVPCMachineTemplateStatus(in, out, s)
}

func Convert_Slice_Pointer_string_To_Slice_Pointer_v1beta2_IBMVPCResourceReference(in *[]*string, out *[]*infrav1beta2.IBMVPCResourceReference, _ apiconversion.Scope) error {
	for _, sshKey := range *in {
		*out = append(*out, &infrav1beta2.IBMVPCResourceReference{
			ID: sshKey,
		})
	}
	return nil
}

func Convert_Slice_Pointer_v1beta2_IBMVPCResourceReference_To_Slice_Pointer_string(in *[]*infrav1beta2.IBMVPCResourceReference, out *[]*string, _ apiconversion.Scope) error {
	if in != nil {
		for _, sshKey := range *in {
			if sshKey.ID != nil {
				*out = append(*out, sshKey.ID)
			}
		}
	}
	return nil
}

func Convert_v1beta2_VPCLoadBalancerSpec_To_v1beta1_VPCLoadBalancerSpec(in *infrav1beta2.VPCLoadBalancerSpec, out *VPCLoadBalancerSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_VPCLoadBalancerSpec_To_v1beta1_VPCLoadBalancerSpec(in, out, s)
}
