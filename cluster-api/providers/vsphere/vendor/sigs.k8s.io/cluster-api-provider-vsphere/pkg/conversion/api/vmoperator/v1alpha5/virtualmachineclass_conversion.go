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

package v1alpha5

import (
	"context"

	vmoprv1alpha5 "github.com/vmware-tanzu/vm-operator/api/v1alpha5"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
)

func convert_v1alpha5_VirtualMachineClass_To_hub_VirtualMachineClass(_ context.Context, src *vmoprv1alpha5.VirtualMachineClass, dst *vmoprvhub.VirtualMachineClass) error {
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Hardware = vmoprvhub.VirtualMachineClassHardware{
		Cpus:   src.Spec.Hardware.Cpus,
		Memory: src.Spec.Hardware.Memory,
	}

	return nil
}

func convert_hub_VirtualMachineClass_To_v1alpha5_VirtualMachineClass(_ context.Context, src *vmoprvhub.VirtualMachineClass, dst *vmoprv1alpha5.VirtualMachineClass) error {
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Hardware = vmoprv1alpha5.VirtualMachineClassHardware{
		Cpus:   src.Spec.Hardware.Cpus,
		Memory: src.Spec.Hardware.Memory,
	}

	return nil
}

func init() {
	converterBuilder.AddConversion(
		conversion.NewAddConversionBuilder(convert_hub_VirtualMachineClass_To_v1alpha5_VirtualMachineClass, convert_v1alpha5_VirtualMachineClass_To_hub_VirtualMachineClass),
	)
}
