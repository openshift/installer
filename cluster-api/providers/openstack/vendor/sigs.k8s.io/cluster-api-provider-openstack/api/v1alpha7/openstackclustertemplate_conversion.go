/*
Copyright 2023 The Kubernetes Authors.

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

package v1alpha7

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversion"
)

var _ ctrlconversion.Convertible = &OpenStackClusterTemplate{}

func (r *OpenStackClusterTemplate) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterTemplate)

	return conversion.ConvertAndRestore(
		r, dst,
		Convert_v1alpha7_OpenStackClusterTemplate_To_v1beta1_OpenStackClusterTemplate, Convert_v1beta1_OpenStackClusterTemplate_To_v1alpha7_OpenStackClusterTemplate,
		v1alpha7OpenStackClusterTemplateRestorer, v1beta1OpenStackClusterTemplateRestorer,
	)
}

func (r *OpenStackClusterTemplate) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterTemplate)

	return conversion.ConvertAndRestore(
		src, r,
		Convert_v1beta1_OpenStackClusterTemplate_To_v1alpha7_OpenStackClusterTemplate, Convert_v1alpha7_OpenStackClusterTemplate_To_v1beta1_OpenStackClusterTemplate,
		v1beta1OpenStackClusterTemplateRestorer, v1alpha7OpenStackClusterTemplateRestorer,
	)
}

var _ ctrlconversion.Convertible = &OpenStackClusterTemplateList{}

func (r *OpenStackClusterTemplateList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterTemplateList)
	return Convert_v1alpha7_OpenStackClusterTemplateList_To_v1beta1_OpenStackClusterTemplateList(r, dst, nil)
}

func (r *OpenStackClusterTemplateList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterTemplateList)
	return Convert_v1beta1_OpenStackClusterTemplateList_To_v1alpha7_OpenStackClusterTemplateList(src, r, nil)
}

/* Restorers */

var v1alpha7OpenStackClusterTemplateRestorer = conversion.RestorerFor[*OpenStackClusterTemplate]{
	"spec": conversion.HashedFieldRestorer(
		func(c *OpenStackClusterTemplate) *OpenStackClusterTemplateSpec {
			return &c.Spec
		},
		restorev1alpha7ClusterTemplateSpec,
	),
}

var v1beta1OpenStackClusterTemplateRestorer = conversion.RestorerFor[*infrav1.OpenStackClusterTemplate]{
	"spec": conversion.HashedFieldRestorer(
		func(c *infrav1.OpenStackClusterTemplate) *infrav1.OpenStackClusterTemplateSpec {
			return &c.Spec
		},
		restorev1beta1ClusterTemplateSpec,
	),
}

func restorev1alpha7ClusterTemplateSpec(previous *OpenStackClusterTemplateSpec, dst *OpenStackClusterTemplateSpec) {
	restorev1alpha7ClusterSpec(&previous.Template.Spec, &dst.Template.Spec)
	restorev1alpha7Bastion(&previous.Template.Spec.Bastion, &dst.Template.Spec.Bastion)
}

func restorev1beta1ClusterTemplateSpec(previous *infrav1.OpenStackClusterTemplateSpec, dst *infrav1.OpenStackClusterTemplateSpec) {
	restorev1beta1Bastion(&previous.Template.Spec.Bastion, &dst.Template.Spec.Bastion)
	restorev1beta1ClusterSpec(&previous.Template.Spec, &dst.Template.Spec)
}
