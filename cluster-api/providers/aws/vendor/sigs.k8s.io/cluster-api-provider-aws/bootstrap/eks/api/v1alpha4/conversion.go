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
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"sigs.k8s.io/cluster-api-provider-aws/bootstrap/eks/api/v1beta1"
)

// ConvertTo converts the v1alpha4 EKSConfig receiver to a v1beta1 EKSConfig.
func (r *EKSConfig) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.EKSConfig)

	if err := Convert_v1alpha4_EKSConfig_To_v1beta1_EKSConfig(r, dst, nil); err != nil {
		return err
	}

	restored := &v1beta1.EKSConfig{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	restoreSpec(&restored.Spec, &dst.Spec)

	return nil
}

func restoreSpec(rSpec, dSpec *v1beta1.EKSConfigSpec) {
	dSpec.ContainerRuntime = rSpec.ContainerRuntime
	dSpec.DNSClusterIP = rSpec.DNSClusterIP
	dSpec.DockerConfigJSON = rSpec.DockerConfigJSON
	dSpec.APIRetryAttempts = rSpec.APIRetryAttempts
	if rSpec.PauseContainer != nil {
		dSpec.PauseContainer = &v1beta1.PauseContainer{
			AccountNumber: rSpec.PauseContainer.AccountNumber,
			Version:       rSpec.PauseContainer.Version,
		}
	}
	dSpec.UseMaxPods = rSpec.UseMaxPods
}

// ConvertFrom converts the v1beta1 EKSConfig receiver to a v1alpha4 EKSConfig.
func (r *EKSConfig) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.EKSConfig)

	if err := Convert_v1beta1_EKSConfig_To_v1alpha4_EKSConfig(src, r, nil); err != nil {
		return err
	}

	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha4 EKSConfigList receiver to a v1beta1 EKSConfigList.
func (r *EKSConfigList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.EKSConfigList)

	return Convert_v1alpha4_EKSConfigList_To_v1beta1_EKSConfigList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 EKSConfigList receiver to a v1alpha4 EKSConfigList.
func (r *EKSConfigList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.EKSConfigList)

	return Convert_v1beta1_EKSConfigList_To_v1alpha4_EKSConfigList(src, r, nil)
}

// ConvertTo converts the v1alpha4 EKSConfigTemplate receiver to a v1beta1 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.EKSConfigTemplate)

	if err := Convert_v1alpha4_EKSConfigTemplate_To_v1beta1_EKSConfigTemplate(r, dst, nil); err != nil {
		return err
	}

	restored := &v1beta1.EKSConfigTemplate{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	restoreSpec(&restored.Spec.Template.Spec, &dst.Spec.Template.Spec)

	return nil
}

// ConvertFrom converts the v1beta1 EKSConfigTemplate receiver to a v1alpha4 EKSConfigTemplate.
func (r *EKSConfigTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.EKSConfigTemplate)

	if err := Convert_v1beta1_EKSConfigTemplate_To_v1alpha4_EKSConfigTemplate(src, r, nil); err != nil {
		return err
	}

	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1alpha4 EKSConfigTemplateList receiver to a v1beta1 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1beta1.EKSConfigTemplateList)

	return Convert_v1alpha4_EKSConfigTemplateList_To_v1beta1_EKSConfigTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1beta1 EKSConfigTemplateList receiver to a v1alpha4 EKSConfigTemplateList.
func (r *EKSConfigTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1beta1.EKSConfigTemplateList)

	return Convert_v1beta1_EKSConfigTemplateList_To_v1alpha4_EKSConfigTemplateList(src, r, nil)
}

func Convert_v1beta1_EKSConfigSpec_To_v1alpha4_EKSConfigSpec(in *v1beta1.EKSConfigSpec, out *EKSConfigSpec, s apiconversion.Scope) error {
	out.KubeletExtraArgs = in.KubeletExtraArgs

	return nil
}
