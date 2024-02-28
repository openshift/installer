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

package v1alpha6

import (
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversion"
)

var _ ctrlconversion.Convertible = &OpenStackMachineTemplate{}

func (r *OpenStackMachineTemplate) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachineTemplate)

	return conversion.ConvertAndRestore(
		r, dst,
		Convert_v1alpha6_OpenStackMachineTemplate_To_v1beta1_OpenStackMachineTemplate, Convert_v1beta1_OpenStackMachineTemplate_To_v1alpha6_OpenStackMachineTemplate,
		v1alpha6OpenStackMachineTemplateRestorer, v1beta1OpenStackMachineTemplateRestorer,
	)
}

func (r *OpenStackMachineTemplate) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachineTemplate)

	return conversion.ConvertAndRestore(
		src, r,
		Convert_v1beta1_OpenStackMachineTemplate_To_v1alpha6_OpenStackMachineTemplate, Convert_v1alpha6_OpenStackMachineTemplate_To_v1beta1_OpenStackMachineTemplate,
		v1beta1OpenStackMachineTemplateRestorer, v1alpha6OpenStackMachineTemplateRestorer,
	)
}

var _ ctrlconversion.Convertible = &OpenStackMachineTemplateList{}

func (r *OpenStackMachineTemplateList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachineTemplateList)
	return Convert_v1alpha6_OpenStackMachineTemplateList_To_v1beta1_OpenStackMachineTemplateList(r, dst, nil)
}

func (r *OpenStackMachineTemplateList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachineTemplateList)
	return Convert_v1beta1_OpenStackMachineTemplateList_To_v1alpha6_OpenStackMachineTemplateList(src, r, nil)
}

/* Restorers */

var v1alpha6OpenStackMachineTemplateRestorer = conversion.RestorerFor[*OpenStackMachineTemplate]{
	"spec": conversion.HashedFieldRestorer(
		func(c *OpenStackMachineTemplate) *OpenStackMachineSpec {
			return &c.Spec.Template.Spec
		},
		restorev1alpha6MachineTemplateMachineSpec,
	),
}

var v1beta1OpenStackMachineTemplateRestorer = conversion.RestorerFor[*infrav1.OpenStackMachineTemplate]{
	"spec": conversion.HashedFieldRestorer(
		func(c *infrav1.OpenStackMachineTemplate) *infrav1.OpenStackMachineSpec {
			return &c.Spec.Template.Spec
		},
		restorev1beta1MachineSpec,
	),
}

// Restore a MachineSpec in the context of a MachineTemplate. Restores
// InstanceID, which will not otherwise be restored for a MachineTemplate.
func restorev1alpha6MachineTemplateMachineSpec(previous *OpenStackMachineSpec, dst *OpenStackMachineSpec) {
	if previous == nil || dst == nil {
		return
	}
	restorev1alpha6MachineSpec(previous, dst)
	dst.InstanceID = previous.InstanceID
}
