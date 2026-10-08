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

	corev1 "k8s.io/api/core/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/webhook/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta1"
	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
)

// VSphereVM is a HubSpokeConverter for the VSphereVM API type.
var VSphereVM = conversion.NewHubSpokeConverter(&infrav1.VSphereVM{},
	conversion.NewSpokeConverter(&infrav1beta1.VSphereVM{}, ConvertVSphereVMHubToV1Beta1, ConvertVSphereVMV1Beta1ToHub),
)

// ConvertVSphereVMV1Beta1ToHub converts a v1beta1 VSphereVM to a hub VSphereVM.
func ConvertVSphereVMV1Beta1ToHub(_ context.Context, src *infrav1beta1.VSphereVM, dst *infrav1.VSphereVM) error {
	if err := infrav1beta1.Convert_v1beta1_VSphereVM_To_v1beta2_VSphereVM(src, dst, nil); err != nil {
		return err
	}

	restored := &infrav1.VSphereVM{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}

	clusterv1.Convert_int32_To_Pointer_int32(src.Spec.NumCoresPerSocket, ok, restored.Spec.NumCoresPerSocket, &dst.Spec.NumCoresPerSocket)

	if src.Spec.BootstrapRef != nil {
		dst.Spec.BootstrapRef.Name = src.Spec.BootstrapRef.Name
	}

	if len(src.Spec.Network.Routes) == len(dst.Spec.Network.Routes) {
		for i, dstRoute := range dst.Spec.Network.Routes {
			srcRoute := src.Spec.Network.Routes[i]
			var restoredMetric *int32
			if len(dst.Spec.Network.Routes) == len(restored.Spec.Network.Routes) {
				restoredMetric = restored.Spec.Network.Routes[i].Metric
			}
			clusterv1.Convert_int32_To_Pointer_int32(srcRoute.Metric, ok, restoredMetric, &dstRoute.Metric)
			dst.Spec.Network.Routes[i] = dstRoute
		}
	}

	if len(src.Spec.Network.Devices) == len(dst.Spec.Network.Devices) {
		for i, dstDevice := range dst.Spec.Network.Devices {
			srcDevice := src.Spec.Network.Devices[i]
			var restoredDeviceDHCP4, restoredDeviceDHCP6, restoredSkipIPAllocation *bool
			if len(dst.Spec.Network.Devices) == len(restored.Spec.Network.Devices) {
				restoredDevice := restored.Spec.Network.Devices[i]
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

			dst.Spec.Network.Devices[i] = dstDevice
		}
	}

	clusterv1.Convert_bool_To_Pointer_bool(src.Status.Ready, ok, restored.Status.Ready, &dst.Status.Ready)
	if len(src.Status.Network) == len(dst.Status.Network) {
		for i, dstNetwork := range dst.Status.Network {
			srcNetwork := src.Status.Network[i]
			var restoredConnected *bool
			if len(dst.Status.Network) == len(restored.Status.Network) {
				restoredNetwork := restored.Status.Network[i]
				restoredConnected = restoredNetwork.Connected
			}

			clusterv1.Convert_bool_To_Pointer_bool(srcNetwork.Connected, ok, restoredConnected, &dstNetwork.Connected)

			dst.Status.Network[i] = dstNetwork
		}
	}
	return nil
}

// ConvertVSphereVMHubToV1Beta1 converts a hub VSphereVM to a v1beta1 VSphereVM.
func ConvertVSphereVMHubToV1Beta1(_ context.Context, src *infrav1.VSphereVM, dst *infrav1beta1.VSphereVM) error {
	if err := infrav1beta1.Convert_v1beta2_VSphereVM_To_v1beta1_VSphereVM(src, dst, nil); err != nil {
		return err
	}

	if src.Spec.BootstrapRef.Name != "" {
		dst.Spec.BootstrapRef = &corev1.ObjectReference{
			APIVersion: "v1",
			Kind:       "Secret",
			Name:       src.Spec.BootstrapRef.Name,
			Namespace:  src.Namespace,
		}
	}

	return utilconversion.MarshalDataUnsafeNoCopy(src, dst)
}
