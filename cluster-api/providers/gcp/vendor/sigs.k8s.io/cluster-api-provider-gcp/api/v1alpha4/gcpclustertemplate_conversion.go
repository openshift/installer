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

package v1alpha4

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	infrav1beta1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this GCPClusterTemplate to the Hub version (v1beta1).
func (src *GCPClusterTemplate) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.GCPClusterTemplate)

	if err := Convert_v1alpha4_GCPClusterTemplate_To_v1beta1_GCPClusterTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1beta1.GCPClusterTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta

	for _, restoredSubnet := range restored.Spec.Template.Spec.Network.Subnets {
		for i, dstSubnet := range dst.Spec.Template.Spec.Network.Subnets {
			if dstSubnet.Name != restoredSubnet.Name {
				continue
			}
			dst.Spec.Template.Spec.Network.Subnets[i].Purpose = restoredSubnet.Purpose

			break
		}
	}

	if restored.Spec.Template.Spec.CredentialsRef != nil {
		dst.Spec.Template.Spec.CredentialsRef = restored.Spec.Template.Spec.CredentialsRef.DeepCopy()
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *GCPClusterTemplate) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.GCPClusterTemplate)
	if err := Convert_v1beta1_GCPClusterTemplate_To_v1alpha4_GCPClusterTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this GCPClusterTemplateList to the Hub version (v1beta1).
func (src *GCPClusterTemplateList) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.GCPClusterTemplateList)
	return Convert_v1alpha4_GCPClusterTemplateList_To_v1beta1_GCPClusterTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this version.
func (dst *GCPClusterTemplateList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.GCPClusterTemplateList)
	return Convert_v1beta1_GCPClusterTemplateList_To_v1alpha4_GCPClusterTemplateList(src, dst, nil)
}

func Convert_v1beta1_GCPClusterTemplateResource_To_v1alpha4_GCPClusterTemplateResource(in *infrav1beta1.GCPClusterTemplateResource, out *GCPClusterTemplateResource, s apiconversion.Scope) error {
	// NOTE: custom conversion func is required because spec.template.metadata has been added in v1beta1.
	return autoConvert_v1beta1_GCPClusterTemplateResource_To_v1alpha4_GCPClusterTemplateResource(in, out, s)
}
