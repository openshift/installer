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

package v1alpha2

import (
	"context"

	vmoprv1alpha2 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	vmoprv1alpha2common "github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
)

func convert_v1alpha2_VirtualMachineImage_To_hub_VirtualMachineImage(_ context.Context, src *vmoprv1alpha2.VirtualMachineImage, dst *vmoprvhub.VirtualMachineImage) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.ProviderRef != nil {
		dst.Spec.ProviderRef = &vmoprvhub.LocalObjectRef{
			APIVersion: src.Spec.ProviderRef.APIVersion,
			Kind:       src.Spec.ProviderRef.Kind,
			Name:       src.Spec.ProviderRef.Name,
		}
	}

	if src.Status.Conditions != nil {
		dst.Status.Conditions = []metav1.Condition{}
		for _, condition := range src.Status.Conditions {
			dst.Status.Conditions = append(dst.Status.Conditions, condition)
		}
	}
	dst.Status.Name = src.Status.Name
	dst.Status.OSInfo = vmoprvhub.VirtualMachineImageOSInfo{
		Type: src.Status.OSInfo.Type,
	}
	dst.Status.ProductInfo = vmoprvhub.VirtualMachineImageProductInfo{
		FullVersion: src.Status.ProductInfo.FullVersion,
	}
	dst.Status.ProviderItemID = src.Status.ProviderItemID

	// Convert VMwareSystemProperties
	if src.Status.VMwareSystemProperties != nil {
		dst.Status.VMwareSystemProperties = make([]vmoprvhub.KeyValuePair, len(src.Status.VMwareSystemProperties))
		for i, prop := range src.Status.VMwareSystemProperties {
			dst.Status.VMwareSystemProperties[i] = vmoprvhub.KeyValuePair{
				Key:   prop.Key,
				Value: prop.Value,
			}
		}
	}

	return nil
}

func convert_hub_VirtualMachineImage_To_v1alpha2_VirtualMachineImage(_ context.Context, src *vmoprvhub.VirtualMachineImage, dst *vmoprv1alpha2.VirtualMachineImage) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.ProviderRef != nil {
		dst.Spec.ProviderRef = &vmoprv1alpha2common.LocalObjectRef{
			APIVersion: src.Spec.ProviderRef.APIVersion,
			Kind:       src.Spec.ProviderRef.Kind,
			Name:       src.Spec.ProviderRef.Name,
		}
	}

	if src.Status.Conditions != nil {
		dst.Status.Conditions = []metav1.Condition{}
		for _, condition := range src.Status.Conditions {
			dst.Status.Conditions = append(dst.Status.Conditions, condition)
		}
	}
	dst.Status.Name = src.Status.Name
	dst.Status.OSInfo = vmoprv1alpha2.VirtualMachineImageOSInfo{
		Type: src.Status.OSInfo.Type,
	}
	dst.Status.ProductInfo = vmoprv1alpha2.VirtualMachineImageProductInfo{
		FullVersion: src.Status.ProductInfo.FullVersion,
	}
	dst.Status.ProviderItemID = src.Status.ProviderItemID

	// Convert VMwareSystemProperties
	if src.Status.VMwareSystemProperties != nil {
		dst.Status.VMwareSystemProperties = make([]vmoprv1alpha2common.KeyValuePair, len(src.Status.VMwareSystemProperties))
		for i, prop := range src.Status.VMwareSystemProperties {
			dst.Status.VMwareSystemProperties[i] = vmoprv1alpha2common.KeyValuePair{
				Key:   prop.Key,
				Value: prop.Value,
			}
		}
	}

	return nil
}

func convert_v1alpha2_ClusterVirtualMachineImage_To_hub_ClusterVirtualMachineImage(_ context.Context, src *vmoprv1alpha2.ClusterVirtualMachineImage, dst *vmoprvhub.ClusterVirtualMachineImage) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.ProviderRef != nil {
		dst.Spec.ProviderRef = &vmoprvhub.LocalObjectRef{
			APIVersion: src.Spec.ProviderRef.APIVersion,
			Kind:       src.Spec.ProviderRef.Kind,
			Name:       src.Spec.ProviderRef.Name,
		}
	}

	if src.Status.Conditions != nil {
		dst.Status.Conditions = []metav1.Condition{}
		for _, condition := range src.Status.Conditions {
			dst.Status.Conditions = append(dst.Status.Conditions, condition)
		}
	}
	dst.Status.Name = src.Status.Name
	dst.Status.OSInfo = vmoprvhub.VirtualMachineImageOSInfo{
		Type: src.Status.OSInfo.Type,
	}
	dst.Status.ProductInfo = vmoprvhub.VirtualMachineImageProductInfo{
		FullVersion: src.Status.ProductInfo.FullVersion,
	}
	dst.Status.ProviderItemID = src.Status.ProviderItemID

	// Convert VMwareSystemProperties
	if src.Status.VMwareSystemProperties != nil {
		dst.Status.VMwareSystemProperties = make([]vmoprvhub.KeyValuePair, len(src.Status.VMwareSystemProperties))
		for i, prop := range src.Status.VMwareSystemProperties {
			dst.Status.VMwareSystemProperties[i] = vmoprvhub.KeyValuePair{
				Key:   prop.Key,
				Value: prop.Value,
			}
		}
	}

	return nil
}

func convert_hub_ClusterVirtualMachineImage_To_v1alpha2_ClusterVirtualMachineImage(_ context.Context, src *vmoprvhub.ClusterVirtualMachineImage, dst *vmoprv1alpha2.ClusterVirtualMachineImage) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.ProviderRef != nil {
		dst.Spec.ProviderRef = &vmoprv1alpha2common.LocalObjectRef{
			APIVersion: src.Spec.ProviderRef.APIVersion,
			Kind:       src.Spec.ProviderRef.Kind,
			Name:       src.Spec.ProviderRef.Name,
		}
	}

	if src.Status.Conditions != nil {
		dst.Status.Conditions = []metav1.Condition{}
		for _, condition := range src.Status.Conditions {
			dst.Status.Conditions = append(dst.Status.Conditions, condition)
		}
	}
	dst.Status.Name = src.Status.Name
	dst.Status.OSInfo = vmoprv1alpha2.VirtualMachineImageOSInfo{
		Type: src.Status.OSInfo.Type,
	}
	dst.Status.ProductInfo = vmoprv1alpha2.VirtualMachineImageProductInfo{
		FullVersion: src.Status.ProductInfo.FullVersion,
	}
	dst.Status.ProviderItemID = src.Status.ProviderItemID

	// Convert VMwareSystemProperties
	if src.Status.VMwareSystemProperties != nil {
		dst.Status.VMwareSystemProperties = make([]vmoprv1alpha2common.KeyValuePair, len(src.Status.VMwareSystemProperties))
		for i, prop := range src.Status.VMwareSystemProperties {
			dst.Status.VMwareSystemProperties[i] = vmoprv1alpha2common.KeyValuePair{
				Key:   prop.Key,
				Value: prop.Value,
			}
		}
	}

	return nil
}

func init() {
	converterBuilder.AddConversion(
		conversion.NewAddConversionBuilder(convert_hub_VirtualMachineImage_To_v1alpha2_VirtualMachineImage, convert_v1alpha2_VirtualMachineImage_To_hub_VirtualMachineImage),
	)
	converterBuilder.AddConversion(
		conversion.NewAddConversionBuilder(convert_hub_ClusterVirtualMachineImage_To_v1alpha2_ClusterVirtualMachineImage, convert_v1alpha2_ClusterVirtualMachineImage_To_hub_ClusterVirtualMachineImage),
	)
}
