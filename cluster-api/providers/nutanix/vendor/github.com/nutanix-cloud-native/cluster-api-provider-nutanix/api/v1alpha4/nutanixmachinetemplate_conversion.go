/*
Copyright 2022 Nutanix

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
	infrav1beta1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	capiv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this NutanixMachineTemplate to the Hub version (v1beta1).
func (src *NutanixMachineTemplate) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.NutanixMachineTemplate)
	return Convert_v1alpha4_NutanixMachineTemplate_To_v1beta1_NutanixMachineTemplate(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this NutanixMachineTemplate.
func (dst *NutanixMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.NutanixMachineTemplate)
	if err := Convert_v1beta1_NutanixMachineTemplate_To_v1alpha4_NutanixMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this NutanixMachineTemplateList to the Hub version (v1beta1).
func (src *NutanixMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.NutanixMachineTemplateList)
	return Convert_v1alpha4_NutanixMachineTemplateList_To_v1beta1_NutanixMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version(v1beta1) to this NutanixMachineTemplateList.
func (dst *NutanixMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.NutanixMachineTemplateList)
	return Convert_v1beta1_NutanixMachineTemplateList_To_v1alpha4_NutanixMachineTemplateList(src, dst, nil)
}

// Convert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta converts ObjectMeta in NutanixMachineTemplateResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in *capiv1alpha4.ObjectMeta, out *capiv1beta1.ObjectMeta, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return capiv1alpha4.Convert_v1alpha4_ObjectMeta_To_v1beta1_ObjectMeta(in, out, s)
}

// Convert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta converts ObjectMeta in NutanixMachineTemplateResource from v1beta1 to v1alpha4 version.
//
//nolint:all
func Convert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in *capiv1beta1.ObjectMeta, out *capiv1alpha4.ObjectMeta, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return capiv1alpha4.Convert_v1beta1_ObjectMeta_To_v1alpha4_ObjectMeta(in, out, s)
}

// Convert_v1alpha4_NutanixMachineTemplateResource_To_v1beta1_NutanixMachineTemplateResource converts NutanixMachineTemplateResource in NutanixMachineTemplateResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_NutanixMachineTemplateResource_To_v1beta1_NutanixMachineTemplateResource(in *NutanixMachineTemplateResource, out *infrav1beta1.NutanixMachineTemplateResource, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1alpha4_NutanixMachineTemplateResource_To_v1beta1_NutanixMachineTemplateResource(in, out, s)
}

// Convert_v1beta1_NutanixMachineTemplateResource_To_v1alpha4_NutanixMachineTemplateResource converts NutanixMachineTemplateResource in NutanixMachineTemplateResource from v1beta1 to v1alpha4 version.
//
//nolint:all
func Convert_v1beta1_NutanixMachineTemplateResource_To_v1alpha4_NutanixMachineTemplateResource(in *infrav1beta1.NutanixMachineTemplateResource, out *NutanixMachineTemplateResource, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1beta1_NutanixMachineTemplateResource_To_v1alpha4_NutanixMachineTemplateResource(in, out, s)
}
