/*
Copyright 2022 Nutanix

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

package v1alpha4

import (
	infrav1beta1 "github.com/nutanix-cloud-native/cluster-api-provider-nutanix/api/v1beta1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts this NutanixMachine to the Hub version (v1beta1).
func (src *NutanixMachine) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.NutanixMachine)
	return Convert_v1alpha4_NutanixMachine_To_v1beta1_NutanixMachine(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this NutanixMachine.
func (dst *NutanixMachine) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.NutanixMachine)
	return Convert_v1beta1_NutanixMachine_To_v1alpha4_NutanixMachine(src, dst, nil)
}

// ConvertTo converts this NutanixMachineList to the Hub version (v1beta1).
func (src *NutanixMachineList) ConvertTo(dstRaw conversion.Hub) error { // nolint
	dst := dstRaw.(*infrav1beta1.NutanixMachineList)
	return Convert_v1alpha4_NutanixMachineList_To_v1beta1_NutanixMachineList(src, dst, nil)
}

// ConvertFrom converts from the Hub version(v1beta1) to this NutanixMachineList.
func (dst *NutanixMachineList) ConvertFrom(srcRaw conversion.Hub) error { // nolint
	src := srcRaw.(*infrav1beta1.NutanixMachineList)
	return Convert_v1beta1_NutanixMachineList_To_v1alpha4_NutanixMachineList(src, dst, nil)
}

// Convert_v1alpha4_NutanixMachineSpec_To_v1beta1_NutanixMachineSpec converts NutanixMachineSpec in NutanixMachineResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_NutanixMachineSpec_To_v1beta1_NutanixMachineSpec(in *NutanixMachineSpec, out *infrav1beta1.NutanixMachineSpec, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1alpha4_NutanixMachineSpec_To_v1beta1_NutanixMachineSpec(in, out, s)
}

// Convert_v1beta1_NutanixMachineSpec_To_v1alpha4_NutanixMachineSpec converts NutanixMachineSpec in NutanixMachineResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1beta1_NutanixMachineSpec_To_v1alpha4_NutanixMachineSpec(in *infrav1beta1.NutanixMachineSpec, out *NutanixMachineSpec, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1beta1_NutanixMachineSpec_To_v1alpha4_NutanixMachineSpec(in, out, s)
}

// Convert_v1alpha4_NutanixMachineStatus_To_v1beta1_NutanixMachineStatus converts NutanixMachineStatus in NutanixMachineResource from v1alpha4 to v1beta1 version.
//
//nolint:all
func Convert_v1alpha4_NutanixMachineStatus_To_v1beta1_NutanixMachineStatus(in *NutanixMachineStatus, out *infrav1beta1.NutanixMachineStatus, s apiconversion.Scope) error {
	// Wrapping the conversion function to avoid compilation errors due to compileErrorOnMissingConversion()
	// Ref: https://github.com/kubernetes/kubernetes/issues/98380
	return Convert_v1alpha4_NutanixMachineStatus_To_v1beta1_NutanixMachineStatus(in, out, s)
}
