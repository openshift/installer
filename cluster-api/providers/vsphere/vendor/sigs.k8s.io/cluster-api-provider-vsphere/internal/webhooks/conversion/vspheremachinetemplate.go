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

// VSphereMachineTemplate is a HubSpokeConverter for the VSphereMachineTemplate API type.
var VSphereMachineTemplate = conversion.NewHubSpokeConverter(&infrav1.VSphereMachineTemplate{},
	conversion.NewSpokeConverter(&infrav1beta1.VSphereMachineTemplate{}, ConvertVSphereMachineTemplateHubToV1Beta1, ConvertVSphereMachineTemplateV1Beta1ToHub),
)

// ConvertVSphereMachineTemplateV1Beta1ToHub converts a v1beta1 VSphereMachineTemplate to a hub VSphereMachineTemplate.
func ConvertVSphereMachineTemplateV1Beta1ToHub(_ context.Context, src *infrav1beta1.VSphereMachineTemplate, dst *infrav1.VSphereMachineTemplate) error {
	if err := infrav1beta1.Convert_v1beta1_VSphereMachineTemplate_To_v1beta2_VSphereMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	restored := &infrav1.VSphereMachineTemplate{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}

	clusterv1.Convert_int32_To_Pointer_int32(src.Spec.Template.Spec.NumCoresPerSocket, ok, restored.Spec.Template.Spec.NumCoresPerSocket, &dst.Spec.Template.Spec.NumCoresPerSocket)

	if len(src.Spec.Template.Spec.Network.Routes) == len(dst.Spec.Template.Spec.Network.Routes) {
		for i, dstRoute := range dst.Spec.Template.Spec.Network.Routes {
			srcRoute := src.Spec.Template.Spec.Network.Routes[i]
			var restoredMetric *int32
			if len(dst.Spec.Template.Spec.Network.Routes) == len(restored.Spec.Template.Spec.Network.Routes) {
				restoredMetric = restored.Spec.Template.Spec.Network.Routes[i].Metric
			}
			clusterv1.Convert_int32_To_Pointer_int32(srcRoute.Metric, ok, restoredMetric, &dstRoute.Metric)
			dst.Spec.Template.Spec.Network.Routes[i] = dstRoute
		}
	}

	if len(src.Spec.Template.Spec.Network.Devices) == len(dst.Spec.Template.Spec.Network.Devices) {
		for i, dstDevice := range dst.Spec.Template.Spec.Network.Devices {
			srcDevice := src.Spec.Template.Spec.Network.Devices[i]
			var restoredDeviceDHCP4, restoredDeviceDHCP6, restoredSkipIPAllocation *bool
			if len(dst.Spec.Template.Spec.Network.Devices) == len(restored.Spec.Template.Spec.Network.Devices) {
				restoredDevice := restored.Spec.Template.Spec.Network.Devices[i]
				restoredDeviceDHCP4 = restoredDevice.DHCP4
				restoredDeviceDHCP6 = restoredDevice.DHCP6
				restoredSkipIPAllocation = restoredDevice.SkipIPAllocation

				if len(srcDevice.Routes) == len(dstDevice.Routes) {
					for i, dstRoute := range dstDevice.Routes {
						srcRoute := srcDevice.Routes[i]
						var restoredMetric *int32
						if len(dstDevice.Routes) == len(restoredDevice.Routes) {
							restoredMetric = restoredDevice.Routes[i].Metric
						}
						clusterv1.Convert_int32_To_Pointer_int32(srcRoute.Metric, ok, restoredMetric, &dstRoute.Metric)
						dstDevice.Routes[i] = dstRoute
					}
				}
			}

			clusterv1.Convert_bool_To_Pointer_bool(srcDevice.DHCP4, ok, restoredDeviceDHCP4, &dstDevice.DHCP4)
			clusterv1.Convert_bool_To_Pointer_bool(srcDevice.DHCP6, ok, restoredDeviceDHCP6, &dstDevice.DHCP6)
			clusterv1.Convert_bool_To_Pointer_bool(srcDevice.SkipIPAllocation, ok, restoredSkipIPAllocation, &dstDevice.SkipIPAllocation)

			dst.Spec.Template.Spec.Network.Devices[i] = dstDevice
		}
	}

	return nil
}

// ConvertVSphereMachineTemplateHubToV1Beta1 converts a hub VSphereMachineTemplate to a v1beta1 VSphereMachineTemplate.
func ConvertVSphereMachineTemplateHubToV1Beta1(_ context.Context, src *infrav1.VSphereMachineTemplate, dst *infrav1beta1.VSphereMachineTemplate) error {
	if err := infrav1beta1.Convert_v1beta2_VSphereMachineTemplate_To_v1beta1_VSphereMachineTemplate(src, dst, nil); err != nil {
		return err
	}

	if dst.Spec.Template.Spec.ProviderID != nil && *dst.Spec.Template.Spec.ProviderID == "" {
		dst.Spec.Template.Spec.ProviderID = nil
	}

	if dst.Spec.Template.Spec.FailureDomain != nil && *dst.Spec.Template.Spec.FailureDomain == "" {
		dst.Spec.Template.Spec.FailureDomain = nil
	}

	if dst.Spec.Template.Spec.NamingStrategy != nil && dst.Spec.Template.Spec.NamingStrategy.Template != nil && *dst.Spec.Template.Spec.NamingStrategy.Template == "" {
		dst.Spec.Template.Spec.NamingStrategy.Template = nil
	}

	return utilconversion.MarshalDataUnsafeNoCopy(src, dst)
}
