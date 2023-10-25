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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this GCPMachineTemplate to the Hub version (v1beta1).
func (src *GCPMachineTemplate) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.GCPMachineTemplate)

	if err := Convert_v1alpha4_GCPMachineTemplate_To_v1beta1_GCPMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1beta1.GCPMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta

	if restored.Spec.Template.Spec.IPForwarding != nil {
		dst.Spec.Template.Spec.IPForwarding = restored.Spec.Template.Spec.IPForwarding
	}

	if restored.Spec.Template.Spec.ShieldedInstanceConfig != nil {
		dst.Spec.Template.Spec.ShieldedInstanceConfig = restored.Spec.Template.Spec.ShieldedInstanceConfig
	}
	if restored.Spec.Template.Spec.OnHostMaintenance != nil {
		dst.Spec.Template.Spec.OnHostMaintenance = restored.Spec.Template.Spec.OnHostMaintenance
	}

	if restored.Spec.Template.Spec.ConfidentialCompute != nil {
		dst.Spec.Template.Spec.ConfidentialCompute = restored.Spec.Template.Spec.ConfidentialCompute
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *GCPMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.GCPMachineTemplate)
	if err := Convert_v1beta1_GCPMachineTemplate_To_v1alpha4_GCPMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this GCPMachineTemplateList to the Hub version (v1beta1).
func (src *GCPMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.GCPMachineTemplateList)
	return Convert_v1alpha4_GCPMachineTemplateList_To_v1beta1_GCPMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *GCPMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.GCPMachineTemplateList)
	return Convert_v1beta1_GCPMachineTemplateList_To_v1alpha4_GCPMachineTemplateList(src, dst, nil)
}

func Convert_v1beta1_GCPMachineTemplateResource_To_v1alpha4_GCPMachineTemplateResource(in *infrav1beta1.GCPMachineTemplateResource, out *GCPMachineTemplateResource, s apiconversion.Scope) error {
	// NOTE: custom conversion func is required because spec.template.metadata has been added in v1beta1.
	return autoConvert_v1beta1_GCPMachineTemplateResource_To_v1alpha4_GCPMachineTemplateResource(in, out, s)
}
