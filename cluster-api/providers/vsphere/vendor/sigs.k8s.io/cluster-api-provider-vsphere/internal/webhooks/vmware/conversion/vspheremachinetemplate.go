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

	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	vmwarev1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta1"
	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
)

// VSphereMachineTemplate is a HubSpokeConverter for the VSphereMachineTemplate API type.
var VSphereMachineTemplate = conversion.NewHubSpokeConverter(&vmwarev1.VSphereMachineTemplate{},
	conversion.NewSpokeConverter(&vmwarev1beta1.VSphereMachineTemplate{}, ConvertVSphereMachineTemplateHubToV1Beta1, ConvertVSphereMachineTemplateV1Beta1ToHub),
)

// ConvertVSphereMachineTemplateV1Beta1ToHub converts a v1beta1 VSphereMachineTemplate to a hub VSphereMachineTemplate.
func ConvertVSphereMachineTemplateV1Beta1ToHub(_ context.Context, src *vmwarev1beta1.VSphereMachineTemplate, dst *vmwarev1.VSphereMachineTemplate) error {
	if err := vmwarev1beta1.Convert_v1beta1_VSphereMachineTemplate_To_v1beta2_VSphereMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	restored := &vmwarev1.VSphereMachineTemplate{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}

	if ok {
		dst.Status.NodeInfo = restored.Status.NodeInfo
	}

	return nil
}

// ConvertVSphereMachineTemplateHubToV1Beta1 converts a hub VSphereMachineTemplate to a v1beta1 VSphereMachineTemplate.
func ConvertVSphereMachineTemplateHubToV1Beta1(_ context.Context, src *vmwarev1.VSphereMachineTemplate, dst *vmwarev1beta1.VSphereMachineTemplate) error {
	if err := vmwarev1beta1.Convert_v1beta2_VSphereMachineTemplate_To_v1beta1_VSphereMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	if dst.Spec.Template.Spec.ProviderID != nil && *dst.Spec.Template.Spec.ProviderID == "" {
		dst.Spec.Template.Spec.ProviderID = nil
	}

	if dst.Spec.Template.Spec.FailureDomain != nil && *dst.Spec.Template.Spec.FailureDomain == "" {
		dst.Spec.Template.Spec.FailureDomain = nil
	}

	return utilconversion.MarshalDataUnsafeNoCopy(src, dst)
}
