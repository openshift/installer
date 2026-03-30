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

package v1alpha4

import (
	"unsafe"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereMachineTemplate to the Hub version (v1beta1).
func (src *VSphereMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereMachineTemplate)
	if err := Convert_v1alpha4_VSphereMachineTemplate_To_v1beta1_VSphereMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1.VSphereMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	dst.Spec.Template.Spec.TagIDs = restored.Spec.Template.Spec.TagIDs
	dst.Spec.Template.Spec.AdditionalDisksGiB = restored.Spec.Template.Spec.AdditionalDisksGiB
	dst.Spec.Template.Spec.PowerOffMode = restored.Spec.Template.Spec.PowerOffMode
	dst.Spec.Template.Spec.GuestSoftPowerOffTimeout = restored.Spec.Template.Spec.GuestSoftPowerOffTimeout
	dst.Spec.Template.Spec.NamingStrategy = restored.Spec.Template.Spec.NamingStrategy
	for i := range dst.Spec.Template.Spec.Network.Devices {
		dst.Spec.Template.Spec.Network.Devices[i].AddressesFromPools = restored.Spec.Template.Spec.Network.Devices[i].AddressesFromPools
		dst.Spec.Template.Spec.Network.Devices[i].DHCP4Overrides = restored.Spec.Template.Spec.Network.Devices[i].DHCP4Overrides
		dst.Spec.Template.Spec.Network.Devices[i].DHCP6Overrides = restored.Spec.Template.Spec.Network.Devices[i].DHCP6Overrides
		dst.Spec.Template.Spec.Network.Devices[i].SkipIPAllocation = restored.Spec.Template.Spec.Network.Devices[i].SkipIPAllocation
	}
	dst.Spec.Template.Spec.DataDisks = restored.Spec.Template.Spec.DataDisks

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereMachineTemplate.
func (dst *VSphereMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereMachineTemplate)
	if err := Convert_v1beta1_VSphereMachineTemplate_To_v1alpha4_VSphereMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *VSphereMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereMachineTemplateList)
	return Convert_v1alpha4_VSphereMachineTemplateList_To_v1beta1_VSphereMachineTemplateList(src, dst, nil)
}

func (dst *VSphereMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereMachineTemplateList)
	return Convert_v1beta1_VSphereMachineTemplateList_To_v1alpha4_VSphereMachineTemplateList(src, dst, nil)
}

func Convert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in *ObjectMeta, out *clusterv1beta1.ObjectMeta, s apiconversion.Scope) error {
	// wrapping the conversion func to avoid having compile errors due to compileErrorOnMissingConversion()
	// more details at https://github.com/kubernetes/kubernetes/issues/98380
	return autoConvert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in, out, s)
}

// autoConvert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta was copied over from CAPI because it is now internal there.
func autoConvert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in *ObjectMeta, out *clusterv1beta1.ObjectMeta, s apiconversion.Scope) error {
	out.Labels = *(*map[string]string)(unsafe.Pointer(&in.Labels))
	out.Annotations = *(*map[string]string)(unsafe.Pointer(&in.Annotations))
	return nil
}

func Convert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in *clusterv1beta1.ObjectMeta, out *ObjectMeta, s apiconversion.Scope) error {
	// wrapping the conversion func to avoid having compile errors due to compileErrorOnMissingConversion()
	// more details at https://github.com/kubernetes/kubernetes/issues/98380
	return autoConvert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in, out, s)
}

// autoConvert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta was copied over from CAPI because it is now internal there.
func autoConvert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in *clusterv1beta1.ObjectMeta, out *ObjectMeta, s apiconversion.Scope) error {
	out.Labels = *(*map[string]string)(unsafe.Pointer(&in.Labels))
	out.Annotations = *(*map[string]string)(unsafe.Pointer(&in.Annotations))
	return nil
}
