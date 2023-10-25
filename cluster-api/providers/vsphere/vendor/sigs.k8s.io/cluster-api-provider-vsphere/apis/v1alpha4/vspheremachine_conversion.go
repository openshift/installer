/*
Copyright 2021 The Kubernetes Authors.

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

//nolint:forcetypeassert,golint,revive,stylecheck
package v1alpha4

import (
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1beta1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereMachine to the Hub version (v1beta1).
func (src *VSphereMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereMachine)
	if err := Convert_v1alpha4_VSphereMachine_To_v1beta1_VSphereMachine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1beta1.VSphereMachine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.AdditionalDisksGiB = restored.Spec.AdditionalDisksGiB
	dst.Spec.TagIDs = restored.Spec.TagIDs
	dst.Spec.PowerOffMode = restored.Spec.PowerOffMode
	dst.Spec.GuestSoftPowerOffTimeout = restored.Spec.GuestSoftPowerOffTimeout
	for i := range dst.Spec.Network.Devices {
		dst.Spec.Network.Devices[i].AddressesFromPools = restored.Spec.Network.Devices[i].AddressesFromPools
		dst.Spec.Network.Devices[i].DHCP4Overrides = restored.Spec.Network.Devices[i].DHCP4Overrides
		dst.Spec.Network.Devices[i].DHCP6Overrides = restored.Spec.Network.Devices[i].DHCP6Overrides
	}

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereMachine.
func (dst *VSphereMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta1.VSphereMachine)
	return Convert_v1beta1_VSphereMachine_To_v1alpha4_VSphereMachine(src, dst, nil)
}

// ConvertTo converts this VSphereMachineList to the Hub version (v1beta1).
func (src *VSphereMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1beta1.VSphereMachineList)
	return Convert_v1alpha4_VSphereMachineList_To_v1beta1_VSphereMachineList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereMachineList.
func (dst *VSphereMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1beta1.VSphereMachineList)
	return Convert_v1beta1_VSphereMachineList_To_v1alpha4_VSphereMachineList(src, dst, nil)
}
