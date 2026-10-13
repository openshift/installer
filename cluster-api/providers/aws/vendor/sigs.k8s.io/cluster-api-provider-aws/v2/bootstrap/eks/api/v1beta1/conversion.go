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
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"sigs.k8s.io/cluster-api-provider-aws/v2/bootstrap/eks/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

// ConvertTo converts the v1beta1 EKSConfig receiver to a v1beta2 EKSConfig.
func (r *EKSConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfig)

	if err := Convert_v1beta1_EKSConfig_To_v1beta2_EKSConfig(r, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &v1beta2.EKSConfig{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.PreBootstrapCommands != nil {
		dst.Spec.PreBootstrapCommands = restored.Spec.PreBootstrapCommands
	}
	if restored.Spec.PostBootstrapCommands != nil {
		dst.Spec.PostBootstrapCommands = restored.Spec.PostBootstrapCommands
	}
	if restored.Spec.BootstrapCommandOverride != nil {
		dst.Spec.BootstrapCommandOverride = restored.Spec.BootstrapCommandOverride
	}
	if restored.Spec.Files != nil {
		dst.Spec.Files = restored.Spec.Files
	}
	if restored.Spec.DiskSetup != nil {
		dst.Spec.DiskSetup = restored.Spec.DiskSetup
	}
	if restored.Spec.Mounts != nil {
		dst.Spec.Mounts = restored.Spec.Mounts
	}
	if restored.Spec.Users != nil {
		dst.Spec.Users = restored.Spec.Users
	}
	if restored.Spec.NTP != nil {
		dst.Spec.NTP = restored.Spec.NTP
	}

	return nil
}

// ConvertFrom converts the v1beta2 EKSConfig receiver to a v1beta1 EKSConfig.
func (r *EKSConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfig)

	if err := Convert_v1beta2_EKSConfig_To_v1beta1_EKSConfig(src, r, nil); err != nil {
		return err
	}

	return utilconversion.MarshalData(src, r)
}

// ConvertTo converts the v1beta1 EKSConfigList receiver to a v1beta2 EKSConfigList.
func (r *EKSConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfigList)

	return Convert_v1beta1_EKSConfigList_To_v1beta2_EKSConfigList(r, dst, nil)
}

// ConvertFrom converts the v1beta2 EKSConfigList receiver to a v1beta1 EKSConfigList.
func (r *EKSConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfigList)

	return Convert_v1beta2_EKSConfigList_To_v1beta1_EKSConfigList(src, r, nil)
}

// ConvertTo converts the v1beta1 EKSConfigTemplate receiver to a v1beta2 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfigTemplate)

	if err := Convert_v1beta1_EKSConfigTemplate_To_v1beta2_EKSConfigTemplate(r, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &v1beta2.EKSConfigTemplate{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.Template.Spec.PreBootstrapCommands != nil {
		dst.Spec.Template.Spec.PreBootstrapCommands = restored.Spec.Template.Spec.PreBootstrapCommands
	}
	if restored.Spec.Template.Spec.PostBootstrapCommands != nil {
		dst.Spec.Template.Spec.PostBootstrapCommands = restored.Spec.Template.Spec.PostBootstrapCommands
	}
	if restored.Spec.Template.Spec.BootstrapCommandOverride != nil {
		dst.Spec.Template.Spec.BootstrapCommandOverride = restored.Spec.Template.Spec.BootstrapCommandOverride
	}
	if restored.Spec.Template.Spec.Files != nil {
		dst.Spec.Template.Spec.Files = restored.Spec.Template.Spec.Files
	}
	if restored.Spec.Template.Spec.DiskSetup != nil {
		dst.Spec.Template.Spec.DiskSetup = restored.Spec.Template.Spec.DiskSetup
	}
	if restored.Spec.Template.Spec.Mounts != nil {
		dst.Spec.Template.Spec.Mounts = restored.Spec.Template.Spec.Mounts
	}
	if restored.Spec.Template.Spec.Users != nil {
		dst.Spec.Template.Spec.Users = restored.Spec.Template.Spec.Users
	}
	if restored.Spec.Template.Spec.NTP != nil {
		dst.Spec.Template.Spec.NTP = restored.Spec.Template.Spec.NTP
	}

	return nil
}

// ConvertFrom converts the v1beta2 EKSConfigTemplate receiver to a v1beta1 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfigTemplate)

	if err := Convert_v1beta2_EKSConfigTemplate_To_v1beta1_EKSConfigTemplate(src, r, nil); err != nil {
		return err
	}

	return utilconversion.MarshalData(src, r)
}

// ConvertTo converts the v1beta1 EKSConfigTemplateList receiver to a v1beta2 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta2.EKSConfigTemplateList)

	return Convert_v1beta1_EKSConfigTemplateList_To_v1beta2_EKSConfigTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1beta2 EKSConfigTemplateList receiver to a v1beta1 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta2.EKSConfigTemplateList)

	return Convert_v1beta2_EKSConfigTemplateList_To_v1beta1_EKSConfigTemplateList(src, r, nil)
}

// Convert_v1beta2_EKSConfigSpec_To_v1beta1_EKSConfigSpec converts a v1beta2 EKSConfigSpec receiver to a v1beta1 EKSConfigSpec.
func Convert_v1beta2_EKSConfigSpec_To_v1beta1_EKSConfigSpec(in *v1beta2.EKSConfigSpec, out *EKSConfigSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_EKSConfigSpec_To_v1beta1_EKSConfigSpec(in, out, s)
}
