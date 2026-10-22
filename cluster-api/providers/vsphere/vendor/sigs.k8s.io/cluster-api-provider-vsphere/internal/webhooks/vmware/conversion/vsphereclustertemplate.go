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

	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	vmwarev1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
)

// VSphereClusterTemplate is a HubSpokeConverter for the VSphereClusterTemplate API type.
var VSphereClusterTemplate = conversion.NewHubSpokeConverter(&vmwarev1.VSphereClusterTemplate{},
	conversion.NewSpokeConverter(&vmwarev1beta1.VSphereClusterTemplate{}, ConvertVSphereClusterTemplateHubToV1Beta1, ConvertVSphereClusterTemplateV1Beta1ToHub),
)

// ConvertVSphereClusterTemplateV1Beta1ToHub converts a v1beta1 VSphereClusterTemplate to a hub VSphereClusterTemplate.
func ConvertVSphereClusterTemplateV1Beta1ToHub(_ context.Context, src *vmwarev1beta1.VSphereClusterTemplate, dst *vmwarev1.VSphereClusterTemplate) error {
	return vmwarev1beta1.Convert_v1beta1_VSphereClusterTemplate_To_v1beta2_VSphereClusterTemplate(src, dst, nil)
}

// ConvertVSphereClusterTemplateHubToV1Beta1 converts a hub VSphereClusterTemplate to a v1beta1 VSphereClusterTemplate.
func ConvertVSphereClusterTemplateHubToV1Beta1(_ context.Context, src *vmwarev1.VSphereClusterTemplate, dst *vmwarev1beta1.VSphereClusterTemplate) error {
	return vmwarev1beta1.Convert_v1beta2_VSphereClusterTemplate_To_v1beta1_VSphereClusterTemplate(src, dst, nil)
}
