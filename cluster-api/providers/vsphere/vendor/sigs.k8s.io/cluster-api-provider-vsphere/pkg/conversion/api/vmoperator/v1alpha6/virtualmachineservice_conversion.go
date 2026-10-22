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

package v1alpha6

import (
	"context"

	vmoprv1alpha6 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
)

func convert_v1alpha6_VirtualMachineService_To_hub_VirtualMachineService(_ context.Context, src *vmoprv1alpha6.VirtualMachineService, dst *vmoprvhub.VirtualMachineService) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.IPFamilies = src.Spec.IPFamilies
	if src.Spec.IPFamilyPolicy != nil {
		p := *src.Spec.IPFamilyPolicy
		dst.Spec.IPFamilyPolicy = &p
	}

	if src.Spec.Ports != nil {
		dst.Spec.Ports = []vmoprvhub.VirtualMachineServicePort{}
		for _, port := range src.Spec.Ports {
			dst.Spec.Ports = append(dst.Spec.Ports, vmoprvhub.VirtualMachineServicePort{
				Name:       port.Name,
				Protocol:   port.Protocol,
				Port:       port.Port,
				TargetPort: port.TargetPort,
			})
		}
	}
	dst.Spec.Selector = src.Spec.Selector
	dst.Spec.Type = vmoprvhub.VirtualMachineServiceType(src.Spec.Type)

	if src.Status.LoadBalancer.Ingress != nil {
		dst.Status.LoadBalancer.Ingress = []vmoprvhub.LoadBalancerIngress{}
		for _, ingress := range src.Status.LoadBalancer.Ingress {
			dst.Status.LoadBalancer.Ingress = append(dst.Status.LoadBalancer.Ingress, vmoprvhub.LoadBalancerIngress{
				IP:       ingress.IP,
				Hostname: ingress.Hostname,
			})
		}
	}

	return nil
}

func convert_hub_VirtualMachineService_To_v1alpha6_VirtualMachineService(_ context.Context, src *vmoprvhub.VirtualMachineService, dst *vmoprv1alpha6.VirtualMachineService) error {
	dst.ObjectMeta = src.ObjectMeta
	dst.Spec.IPFamilies = src.Spec.IPFamilies
	if src.Spec.IPFamilyPolicy != nil {
		p := *src.Spec.IPFamilyPolicy
		dst.Spec.IPFamilyPolicy = &p
	}

	if src.Spec.Ports != nil {
		dst.Spec.Ports = []vmoprv1alpha6.VirtualMachineServicePort{}
		for _, port := range src.Spec.Ports {
			dst.Spec.Ports = append(dst.Spec.Ports, vmoprv1alpha6.VirtualMachineServicePort{
				Name:       port.Name,
				Protocol:   port.Protocol,
				Port:       port.Port,
				TargetPort: port.TargetPort,
			})
		}
	}
	dst.Spec.Selector = src.Spec.Selector
	dst.Spec.Type = vmoprv1alpha6.VirtualMachineServiceType(src.Spec.Type)

	if src.Status.LoadBalancer.Ingress != nil {
		dst.Status.LoadBalancer.Ingress = []vmoprv1alpha6.LoadBalancerIngress{}
		for _, ingress := range src.Status.LoadBalancer.Ingress {
			dst.Status.LoadBalancer.Ingress = append(dst.Status.LoadBalancer.Ingress, vmoprv1alpha6.LoadBalancerIngress{
				IP:       ingress.IP,
				Hostname: ingress.Hostname,
			})
		}
	}

	return nil
}

func init() {
	converterBuilder.AddConversion(
		conversion.NewAddConversionBuilder(convert_hub_VirtualMachineService_To_v1alpha6_VirtualMachineService, convert_v1alpha6_VirtualMachineService_To_hub_VirtualMachineService),
	)
}
