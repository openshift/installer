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
	"unsafe"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 AWSMachine receiver to a v1beta1 AWSMachine.
func (r *AWSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachine)
	if err := Convert_v1alpha3_AWSMachine_To_v1beta1_AWSMachine(r, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav1.AWSMachine{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	restoreSpec(&restored.Spec, &dst.Spec)

	dst.Spec.Ignition = restored.Spec.Ignition

	return nil
}

func restoreSpec(restored, dst *infrav1.AWSMachineSpec) {
	RestoreAMIReference(&restored.AMI, &dst.AMI)
	if restored.RootVolume != nil {
		if dst.RootVolume == nil {
			dst.RootVolume = &infrav1.Volume{}
		}
		RestoreRootVolume(restored.RootVolume, dst.RootVolume)
	}
	if restored.NonRootVolumes != nil {
		if dst.NonRootVolumes == nil {
			dst.NonRootVolumes = []infrav1.Volume{}
		}
		restoreNonRootVolumes(restored.NonRootVolumes, dst.NonRootVolumes)
	}
}

// ConvertFrom converts the v1beta1 AWSMachine receiver to a v1alpha3 AWSMachine.
func (r *AWSMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachine)

	if err := Convert_v1beta1_AWSMachine_To_v1alpha3_AWSMachine(src, r, nil); err != nil {
		return err
	}
	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}
	return nil
}

// ConvertTo converts the v1alpha3 AWSMachineList receiver to a v1beta1 AWSMachineList.
func (r *AWSMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachineList)

	return Convert_v1alpha3_AWSMachineList_To_v1beta1_AWSMachineList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineList receiver to a v1alpha3 AWSMachineList.
func (r *AWSMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachineList)

	return Convert_v1beta1_AWSMachineList_To_v1alpha3_AWSMachineList(src, r, nil)
}

// ConvertTo converts the v1alpha3 AWSMachineTemplate receiver to a v1beta1 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachineTemplate)
	if err := Convert_v1alpha3_AWSMachineTemplate_To_v1beta1_AWSMachineTemplate(r, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav1.AWSMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta
	dst.Spec.Template.Spec.Ignition = restored.Spec.Template.Spec.Ignition

	restoreSpec(&restored.Spec.Template.Spec, &dst.Spec.Template.Spec)

	return nil
}

// ConvertFrom converts the v1beta1 AWSMachineTemplate receiver to a v1alpha3 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachineTemplate)

	if err := Convert_v1beta1_AWSMachineTemplate_To_v1alpha3_AWSMachineTemplate(src, r, nil); err != nil {
		return err
	}
	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}
	return nil
}

// ConvertTo converts the v1alpha3 AWSMachineTemplateList receiver to a v1beta1 AWSMachineTemplateList.
func (r *AWSMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachineTemplateList)

	return Convert_v1alpha3_AWSMachineTemplateList_To_v1beta1_AWSMachineTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 AWSMachineTemplateList receiver to a v1alpha3 AWSMachineTemplateList.
func (r *AWSMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachineTemplateList)

	return Convert_v1beta1_AWSMachineTemplateList_To_v1alpha3_AWSMachineTemplateList(src, r, nil)
}

// Convert_v1beta1_Volume_To_v1alpha3_Volume .
func Convert_v1beta1_Volume_To_v1alpha3_Volume(in *infrav1.Volume, out *Volume, s apiconversion.Scope) error {
	return autoConvert_v1beta1_Volume_To_v1alpha3_Volume(in, out, s)
}

// Convert_v1beta1_AWSMachineSpec_To_v1alpha3_AWSMachineSpec .
func Convert_v1beta1_AWSMachineSpec_To_v1alpha3_AWSMachineSpec(in *infrav1.AWSMachineSpec, out *AWSMachineSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSMachineSpec_To_v1alpha3_AWSMachineSpec(in, out, s)
}

// Convert_v1beta1_Instance_To_v1alpha3_Instance .
func Convert_v1beta1_Instance_To_v1alpha3_Instance(in *infrav1.Instance, out *Instance, s apiconversion.Scope) error {
	return autoConvert_v1beta1_Instance_To_v1alpha3_Instance(in, out, s)
}

// Manually restore the instance root device data.
// Assumes restored and dst are non-nil.
func restoreInstance(restored, dst *infrav1.Instance) {
	dst.VolumeIDs = restored.VolumeIDs

	if restored.RootVolume != nil {
		if dst.RootVolume == nil {
			dst.RootVolume = &infrav1.Volume{}
		}
		RestoreRootVolume(restored.RootVolume, dst.RootVolume)
	}

	if restored.NonRootVolumes != nil {
		if dst.NonRootVolumes == nil {
			dst.NonRootVolumes = []infrav1.Volume{}
		}
		restoreNonRootVolumes(restored.NonRootVolumes, dst.NonRootVolumes)
	}
}

// Convert_v1alpha3_AWSResourceReference_To_v1beta1_AMIReference is a conversion function.
func Convert_v1alpha3_AWSResourceReference_To_v1beta1_AMIReference(in *AWSResourceReference, out *infrav1.AMIReference, s apiconversion.Scope) error {
	out.ID = (*string)(unsafe.Pointer(in.ID))
	return nil
}

// Convert_v1beta1_AMIReference_To_v1alpha3_AWSResourceReference is a conversion function.
func Convert_v1beta1_AMIReference_To_v1alpha3_AWSResourceReference(in *infrav1.AMIReference, out *AWSResourceReference, s apiconversion.Scope) error {
	out.ID = (*string)(unsafe.Pointer(in.ID))
	return nil
}

// RestoreAMIReference manually restore the EKSOptimizedLookupType for AWSMachine and AWSMachineTemplate
// Assumes both restored and dst are non-nil.
func RestoreAMIReference(restored, dst *infrav1.AMIReference) {
	if restored.EKSOptimizedLookupType == nil {
		return
	}
	dst.EKSOptimizedLookupType = restored.EKSOptimizedLookupType
}

// restoreNonRootVolumes manually restores the non-root volumes
// Assumes both restoredVolumes and dstVolumes are non-nil.
func restoreNonRootVolumes(restoredVolumes, dstVolumes []infrav1.Volume) {
	// restoring the nonrootvolumes which are missing in dstVolumes
	// restoring dstVolumes[i].Encrypted to nil in order to avoid v1beta1 --> v1alpha3 --> v1beta1 round trip errors
	for i := range restoredVolumes {
		if restoredVolumes[i].Encrypted == nil {
			if len(dstVolumes) <= i {
				dstVolumes = append(dstVolumes, restoredVolumes[i])
			} else {
				dstVolumes[i].Encrypted = nil
			}
		}
		dstVolumes[i].Throughput = restoredVolumes[i].Throughput
	}
}

// RestoreRootVolume manually restores the root volumes.
// Assumes both restored and dst are non-nil.
// Volume.Encrypted type changed from bool in v1alpha3 to *bool in v1beta1
// Volume.Encrypted value as nil/&false in v1beta1 will convert to false in v1alpha3 by auto-conversion, so restoring it to nil in order to avoid v1beta1 --> v1alpha3 --> v1beta1 round trip errors
func RestoreRootVolume(restored, dst *infrav1.Volume) {
	if dst.Encrypted == pointer.BoolPtr(true) {
		return
	}
	if restored.Encrypted == nil {
		dst.Encrypted = nil
	}
	dst.Throughput = restored.Throughput
}

func Convert_v1beta1_AWSMachineTemplateResource_To_v1alpha3_AWSMachineTemplateResource(in *infrav1.AWSMachineTemplateResource, out *AWSMachineTemplateResource, s apiconversion.Scope) error {
	return autoConvert_v1beta1_AWSMachineTemplateResource_To_v1alpha3_AWSMachineTemplateResource(in, out, s)
}
