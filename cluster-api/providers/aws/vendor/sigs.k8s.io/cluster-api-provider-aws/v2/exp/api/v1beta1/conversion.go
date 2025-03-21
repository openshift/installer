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

package v1beta1

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	infrav1exp "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

// ConvertTo converts the v1beta1 AWSMachinePool receiver to a v1beta2 AWSMachinePool.
func (src *AWSMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSMachinePool)
	if err := Convert_v1beta1_AWSMachinePool_To_v1beta2_AWSMachinePool(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1exp.AWSMachinePool{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.SuspendProcesses != nil {
		dst.Spec.SuspendProcesses = restored.Spec.SuspendProcesses
	}
	if restored.Spec.RefreshPreferences != nil {
		dst.Spec.RefreshPreferences.Disable = restored.Spec.RefreshPreferences.Disable
		dst.Spec.RefreshPreferences.MaxHealthyPercentage = restored.Spec.RefreshPreferences.MaxHealthyPercentage
	}
	if restored.Spec.AWSLaunchTemplate.InstanceMetadataOptions != nil {
		dst.Spec.AWSLaunchTemplate.InstanceMetadataOptions = restored.Spec.AWSLaunchTemplate.InstanceMetadataOptions
	}
	if restored.Spec.AvailabilityZoneSubnetType != nil {
		dst.Spec.AvailabilityZoneSubnetType = restored.Spec.AvailabilityZoneSubnetType
	}
	dst.Status.InfrastructureMachineKind = restored.Status.InfrastructureMachineKind

	if restored.Spec.AWSLaunchTemplate.PrivateDNSName != nil {
		dst.Spec.AWSLaunchTemplate.PrivateDNSName = restored.Spec.AWSLaunchTemplate.PrivateDNSName
	}

	if restored.Spec.AWSLaunchTemplate.CapacityReservationID != nil {
		dst.Spec.AWSLaunchTemplate.CapacityReservationID = restored.Spec.AWSLaunchTemplate.CapacityReservationID
	}

	if restored.Spec.AWSLaunchTemplate.MarketType != "" {
		dst.Spec.AWSLaunchTemplate.MarketType = restored.Spec.AWSLaunchTemplate.MarketType
	}

	dst.Spec.DefaultInstanceWarmup = restored.Spec.DefaultInstanceWarmup
	dst.Spec.AWSLaunchTemplate.NonRootVolumes = restored.Spec.AWSLaunchTemplate.NonRootVolumes
	return nil
}

// ConvertFrom converts the v1beta2 AWSMachinePool receiver to v1beta1 AWSMachinePool.
func (dst *AWSMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSMachinePool)

	if err := Convert_v1beta2_AWSMachinePool_To_v1beta1_AWSMachinePool(src, dst, nil); err != nil {
		return err
	}

	return utilconversion.MarshalData(src, dst)
}

// ConvertTo converts the v1beta1 AWSMachinePoolList receiver to a v1beta2 AWSMachinePoolList.
func (src *AWSMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSMachinePoolList)
	return Convert_v1beta1_AWSMachinePoolList_To_v1beta2_AWSMachinePoolList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSMachinePoolList receiver to v1beta1 AWSMachinePoolList.
func (r *AWSMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSMachinePoolList)
	return Convert_v1beta2_AWSMachinePoolList_To_v1beta1_AWSMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1beta1 AWSManagedMachinePool receiver to a v1beta2 AWSManagedMachinePool.
func (src *AWSManagedMachinePool) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSManagedMachinePool)
	if err := Convert_v1beta1_AWSManagedMachinePool_To_v1beta2_AWSManagedMachinePool(src, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav1exp.AWSManagedMachinePool{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.AWSLaunchTemplate != nil {
		if dst.Spec.AWSLaunchTemplate == nil {
			dst.Spec.AWSLaunchTemplate = restored.Spec.AWSLaunchTemplate
		}
		dst.Spec.AWSLaunchTemplate.InstanceMetadataOptions = restored.Spec.AWSLaunchTemplate.InstanceMetadataOptions
		dst.Spec.AWSLaunchTemplate.NonRootVolumes = restored.Spec.AWSLaunchTemplate.NonRootVolumes

		if restored.Spec.AWSLaunchTemplate.PrivateDNSName != nil {
			dst.Spec.AWSLaunchTemplate.PrivateDNSName = restored.Spec.AWSLaunchTemplate.PrivateDNSName
		}

		if restored.Spec.AWSLaunchTemplate.CapacityReservationID != nil {
			dst.Spec.AWSLaunchTemplate.CapacityReservationID = restored.Spec.AWSLaunchTemplate.CapacityReservationID
		}

		if restored.Spec.AWSLaunchTemplate.MarketType != "" {
			dst.Spec.AWSLaunchTemplate.MarketType = restored.Spec.AWSLaunchTemplate.MarketType
		}

	}
	if restored.Spec.AvailabilityZoneSubnetType != nil {
		dst.Spec.AvailabilityZoneSubnetType = restored.Spec.AvailabilityZoneSubnetType
	}

	return nil
}

// ConvertFrom converts the v1beta2 AWSManagedMachinePool receiver to v1beta1 AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSManagedMachinePool)

	if err := Convert_v1beta2_AWSManagedMachinePool_To_v1beta1_AWSManagedMachinePool(src, r, nil); err != nil {
		return err
	}

	return utilconversion.MarshalData(src, r)
}

// Convert_v1beta2_AWSManagedMachinePoolSpec_To_v1beta1_AWSManagedMachinePoolSpec is a conversion function.
func Convert_v1beta2_AWSManagedMachinePoolSpec_To_v1beta1_AWSManagedMachinePoolSpec(in *infrav1exp.AWSManagedMachinePoolSpec, out *AWSManagedMachinePoolSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSManagedMachinePoolSpec_To_v1beta1_AWSManagedMachinePoolSpec(in, out, s)
}

func Convert_v1beta2_AWSMachinePoolStatus_To_v1beta1_AWSMachinePoolStatus(in *infrav1exp.AWSMachinePoolStatus, out *AWSMachinePoolStatus, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSMachinePoolStatus_To_v1beta1_AWSMachinePoolStatus(in, out, s)
}

// ConvertTo converts the v1beta1 AWSManagedMachinePoolList receiver to a v1beta2 AWSManagedMachinePoolList.
func (src *AWSManagedMachinePoolList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSManagedMachinePoolList)
	return Convert_v1beta1_AWSManagedMachinePoolList_To_v1beta2_AWSManagedMachinePoolList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSManagedMachinePoolList receiver to v1beta1 AWSManagedMachinePoolList.
func (r *AWSManagedMachinePoolList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSManagedMachinePoolList)

	return Convert_v1beta2_AWSManagedMachinePoolList_To_v1beta1_AWSManagedMachinePoolList(src, r, nil)
}

// ConvertTo converts the v1beta1 AWSFargateProfile receiver to a v1beta2 AWSFargateProfile.
func (src *AWSFargateProfile) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSFargateProfile)
	return Convert_v1beta1_AWSFargateProfile_To_v1beta2_AWSFargateProfile(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSFargateProfile receiver to v1beta1 AWSFargateProfile.
func (r *AWSFargateProfile) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSFargateProfile)

	return Convert_v1beta2_AWSFargateProfile_To_v1beta1_AWSFargateProfile(src, r, nil)
}

// ConvertTo converts the v1beta1 AWSFargateProfileList receiver to a v1beta2 AWSFargateProfileList.
func (src *AWSFargateProfileList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1exp.AWSFargateProfileList)
	return Convert_v1beta1_AWSFargateProfileList_To_v1beta2_AWSFargateProfileList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSFargateProfileList receiver to v1beta1 AWSFargateProfileList.
func (r *AWSFargateProfileList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1exp.AWSFargateProfileList)

	return Convert_v1beta2_AWSFargateProfileList_To_v1beta1_AWSFargateProfileList(src, r, nil)
}

// Convert_v1beta1_AMIReference_To_v1beta2_AMIReference converts the v1beta1 AMIReference receiver to a v1beta2 AMIReference.
func Convert_v1beta1_AMIReference_To_v1beta2_AMIReference(in *infrav1beta1.AMIReference, out *infrav1.AMIReference, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta1_AMIReference_To_v1beta2_AMIReference(in, out, s)
}

// Convert_v1beta2_AMIReference_To_v1beta1_AMIReference converts the v1beta2 AMIReference receiver to a v1beta1 AMIReference.
func Convert_v1beta2_AMIReference_To_v1beta1_AMIReference(in *infrav1.AMIReference, out *infrav1beta1.AMIReference, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta2_AMIReference_To_v1beta1_AMIReference(in, out, s)
}

// Convert_v1beta2_Instance_To_v1beta1_Instance is a conversion function.
func Convert_v1beta2_Instance_To_v1beta1_Instance(in *infrav1.Instance, out *infrav1beta1.Instance, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta2_Instance_To_v1beta1_Instance(in, out, s)
}

// Convert_v1beta1_Instance_To_v1beta2_Instance is a conversion function.
func Convert_v1beta1_Instance_To_v1beta2_Instance(in *infrav1beta1.Instance, out *infrav1.Instance, s apiconversion.Scope) error {
	return infrav1beta1.Convert_v1beta1_Instance_To_v1beta2_Instance(in, out, s)
}

// Convert_v1beta2_AWSLaunchTemplate_To_v1beta1_AWSLaunchTemplate converts the v1beta2 AWSLaunchTemplate receiver to a v1beta1 AWSLaunchTemplate.
func Convert_v1beta2_AWSLaunchTemplate_To_v1beta1_AWSLaunchTemplate(in *infrav1exp.AWSLaunchTemplate, out *AWSLaunchTemplate, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSLaunchTemplate_To_v1beta1_AWSLaunchTemplate(in, out, s)
}

func Convert_v1beta1_AWSMachinePoolSpec_To_v1beta2_AWSMachinePoolSpec(in *AWSMachinePoolSpec, out *infrav1exp.AWSMachinePoolSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSMachinePoolSpec_To_v1beta2_AWSMachinePoolSpec(in, out, s)
}

func Convert_v1beta2_AWSMachinePoolSpec_To_v1beta1_AWSMachinePoolSpec(in *infrav1exp.AWSMachinePoolSpec, out *AWSMachinePoolSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_AWSMachinePoolSpec_To_v1beta1_AWSMachinePoolSpec(in, out, s)
}

func Convert_v1beta1_AutoScalingGroup_To_v1beta2_AutoScalingGroup(in *AutoScalingGroup, out *infrav1exp.AutoScalingGroup, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AutoScalingGroup_To_v1beta2_AutoScalingGroup(in, out, s)
}

func Convert_v1beta2_AutoScalingGroup_To_v1beta1_AutoScalingGroup(in *infrav1exp.AutoScalingGroup, out *AutoScalingGroup, s apiconversion.Scope) error {
	// explicitly ignore CurrentlySuspended.
	return autoConvert_v1beta2_AutoScalingGroup_To_v1beta1_AutoScalingGroup(in, out, s)
}

// Convert_v1beta2_RefreshPreferences_To_v1beta1_RefreshPreferences converts the v1beta2 RefreshPreferences receiver to a v1beta1 RefreshPreferences.
func Convert_v1beta2_RefreshPreferences_To_v1beta1_RefreshPreferences(in *infrav1exp.RefreshPreferences, out *RefreshPreferences, s apiconversion.Scope) error {
	// spec.refreshPreferences.disable has been added to v1beta2.
	return autoConvert_v1beta2_RefreshPreferences_To_v1beta1_RefreshPreferences(in, out, s)
}
