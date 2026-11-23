/*
Copyright 2026 The Kubernetes Authors.

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

package conversion

import (
	"context"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
)

// VSphereClusterTemplate is a HubSpokeConverter for the VSphereClusterTemplate API type.
var VSphereClusterTemplate = conversion.NewHubSpokeConverter(&infrav1.VSphereClusterTemplate{},
	conversion.NewSpokeConverter(&infrav1beta1.VSphereClusterTemplate{}, ConvertVSphereClusterTemplateHubToV1Beta1, ConvertVSphereClusterTemplateV1Beta1ToHub),
)

// ConvertVSphereClusterTemplateV1Beta1ToHub converts a v1beta1 VSphereClusterTemplate to a hub VSphereClusterTemplate.
func ConvertVSphereClusterTemplateV1Beta1ToHub(_ context.Context, src *infrav1beta1.VSphereClusterTemplate, dst *infrav1.VSphereClusterTemplate) error {
	if err := infrav1beta1.Convert_v1beta1_VSphereClusterTemplate_To_v1beta2_VSphereClusterTemplate(src, dst, nil); err != nil {
		return err
	}

	restored := &infrav1.VSphereClusterTemplate{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}

	clusterv1.Convert_bool_To_Pointer_bool(src.Spec.Template.Spec.DisableClusterModule, ok, restored.Spec.Template.Spec.DisableClusterModule, &dst.Spec.Template.Spec.DisableClusterModule)

	if len(src.Spec.Template.Spec.ClusterModules) == len(dst.Spec.Template.Spec.ClusterModules) {
		for i, dstClusterModule := range dst.Spec.Template.Spec.ClusterModules {
			srcClusterModule := src.Spec.Template.Spec.ClusterModules[i]
			var restoredControlPlane *bool
			if len(dst.Spec.Template.Spec.ClusterModules) == len(restored.Spec.Template.Spec.ClusterModules) {
				restoredClusterModules := restored.Spec.Template.Spec.ClusterModules[i]
				restoredControlPlane = restoredClusterModules.ControlPlane
			}

			clusterv1.Convert_bool_To_Pointer_bool(srcClusterModule.ControlPlane, ok, restoredControlPlane, &dstClusterModule.ControlPlane)

			dst.Spec.Template.Spec.ClusterModules[i] = dstClusterModule
		}
	}

	return nil
}

// ConvertVSphereClusterTemplateHubToV1Beta1 converts a hub VSphereClusterTemplate to a v1beta1 VSphereClusterTemplate.
func ConvertVSphereClusterTemplateHubToV1Beta1(_ context.Context, src *infrav1.VSphereClusterTemplate, dst *infrav1beta1.VSphereClusterTemplate) error {
	if err := infrav1beta1.Convert_v1beta2_VSphereClusterTemplate_To_v1beta1_VSphereClusterTemplate(src, dst, nil); err != nil {
		return err
	}

	return utilconversion.MarshalDataUnsafeNoCopy(src, dst)
}
