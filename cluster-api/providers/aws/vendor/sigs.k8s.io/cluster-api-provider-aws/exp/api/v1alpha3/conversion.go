/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha3

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-aws/api/v1alpha3"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/exp/api/v1beta1"
)

// ConvertTo converts the v1alpha3 AWSMachinePool receiver to a v1beta1 AWSMachinePool.
func (r *AWSMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSMachinePool)
	if err := Convert_v1alpha3_AWSMachinePool_To_v1beta1_AWSMachinePool(r, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav1exp.AWSMachinePool{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	infrav1alpha3.RestoreAMIReference(&restored.Spec.AWSLaunchTemplate.AMI, &dst.Spec.AWSLaunchTemplate.AMI)
	if restored.Spec.AWSLaunchTemplate.RootVolume != nil {
		if dst.Spec.AWSLaunchTemplate.RootVolume == nil {
			dst.Spec.AWSLaunchTemplate.RootVolume = &infrav1.Volume{}
		}
		infrav1alpha3.RestoreRootVolume(restored.Spec.AWSLaunchTemplate.RootVolume, dst.Spec.AWSLaunchTemplate.RootVolume)
	}
	return nil
}

// ConvertFrom converts the v1beta1 AWSMachinePool receiver to a v1alpha3 AWSMachinePool.
func (r *AWSMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSMachinePool)

	if err := Convert_v1beta1_AWSMachinePool_To_v1alpha3_AWSMachinePool(src, r, nil); err != nil {
		return err
	}
	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}
	return nil
}

// ConvertTo converts the v1alpha3 AWSMachinePoolList receiver to a v1beta1 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSMachinePoolList)

	return Convert_v1alpha3_AWSMachinePoolList_To_v1beta1_AWSMachinePoolList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachinePoolList receiver to a v1alpha3 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSMachinePoolList)

	return Convert_v1beta1_AWSMachinePoolList_To_v1alpha3_AWSMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSManagedMachinePool receiver to a v1beta1 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSManagedMachinePool)
	if err := Convert_v1alpha3_AWSManagedMachinePool_To_v1beta1_AWSManagedMachinePool(r, dst, nil); err != nil {
		return err
	}

	restored := &infrav1exp.AWSManagedMachinePool{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Taints = restored.Spec.Taints
	dst.Spec.CapacityType = restored.Spec.CapacityType
	dst.Spec.RoleAdditionalPolicies = restored.Spec.RoleAdditionalPolicies
	dst.Spec.UpdateConfig = restored.Spec.UpdateConfig

	return nil
}

// ConvertFrom converts the v1beta1 AWSManagedMachinePool receiver to a v1alpha3 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSManagedMachinePool)

	if err := Convert_v1beta1_AWSManagedMachinePool_To_v1alpha3_AWSManagedMachinePool(src, r, nil); err != nil {
		return err
	}

	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha3 AWSManagedMachinePoolList receiver to a v1beta1 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSManagedMachinePoolList)

	return Convert_v1alpha3_AWSManagedMachinePoolList_To_v1beta1_AWSManagedMachinePoolList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSManagedMachinePoolList receiver to a v1alpha3 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSManagedMachinePoolList)

	return Convert_v1beta1_AWSManagedMachinePoolList_To_v1alpha3_AWSManagedMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSFargateProfile receiver to a v1beta1 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSFargateProfile)

	return Convert_v1alpha3_AWSFargateProfile_To_v1beta1_AWSFargateProfile(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSFargateProfile receiver to a v1alpha3 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSFargateProfile)

	return Convert_v1beta1_AWSFargateProfile_To_v1alpha3_AWSFargateProfile(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSFargateProfileList receiver to a v1beta1 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSFargateProfileList)

	return Convert_v1alpha3_AWSFargateProfileList_To_v1beta1_AWSFargateProfileList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSFargateProfileList receiver to a v1alpha3 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSFargateProfileList)

	return Convert_v1beta1_AWSFargateProfileList_To_v1alpha3_AWSFargateProfileList(src, r, nil)
}

// Convert_v1alpha3_AWSResourceReference_To_v1beta1_AWSResourceReference is a conversion function.
func Convert_v1alpha3_AWSResourceReference_To_v1beta1_AWSResourceReference(in *infrav1alpha3.AWSResourceReference, out *infrav1.AWSResourceReference, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_AWSResourceReference_To_v1beta1_AWSResourceReference(in, out, s)
}

// Convert_v1beta1_AWSResourceReference_To_v1alpha3_AWSResourceReference conversion function.
func Convert_v1beta1_AWSResourceReference_To_v1alpha3_AWSResourceReference(in *infrav1.AWSResourceReference, out *infrav1alpha3.AWSResourceReference, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1beta1_AWSResourceReference_To_v1alpha3_AWSResourceReference(in, out, s)
}

// Convert_v1beta1_AWSManagedMachinePoolSpec_To_v1alpha3_AWSManagedMachinePoolSpec is a conversion function.
func Convert_v1beta1_AWSManagedMachinePoolSpec_To_v1alpha3_AWSManagedMachinePoolSpec(in *infrav1exp.AWSManagedMachinePoolSpec, out *AWSManagedMachinePoolSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSManagedMachinePoolSpec_To_v1alpha3_AWSManagedMachinePoolSpec(in, out, s)
}

// Convert_v1alpha3_Instance_To_v1beta1_Instance is a conversion function.
func Convert_v1alpha3_Instance_To_v1beta1_Instance(in *infrav1alpha3.Instance, out *infrav1.Instance, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_Instance_To_v1beta1_Instance(in, out, s)
}

// Convert_v1alpha3_Volume_To_v1beta1_Volume is a conversion function.
func Convert_v1alpha3_Volume_To_v1beta1_Volume(in *infrav1alpha3.Volume, out *infrav1.Volume, s apiconversion.Scope) error {
	return infrav1alpha3.Convert_v1alpha3_Volume_To_v1beta1_Volume(in, out, s)
}
