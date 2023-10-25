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
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereClusterTemplate to the Hub version (v1beta1).
func (src *VSphereClusterTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereClusterTemplate)
	return Convert_v1alpha4_VSphereClusterTemplate_To_v1beta1_VSphereClusterTemplate(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereClusterTemplate.
func (dst *VSphereClusterTemplate) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.VSphereClusterTemplate)
	return Convert_v1beta1_VSphereClusterTemplate_To_v1alpha4_VSphereClusterTemplate(src, dst, nil)
}

// ConvertTo converts this VSphereClusterIdentityList to the Hub version (v1beta1).
func (src *VSphereClusterTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereClusterTemplateList)
	return Convert_v1alpha4_VSphereClusterTemplateList_To_v1beta1_VSphereClusterTemplateList(src, dst, nil)
}

// ConvertFrom converts this VSphereClusterIdentityList to the Hub version (v1beta1).
func (dst *VSphereClusterTemplateList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.VSphereClusterTemplateList)
	return Convert_v1beta1_VSphereClusterTemplateList_To_v1alpha4_VSphereClusterTemplateList(src, dst, nil)
}
