/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// ConvertTo converts this VSphereVM to the Hub version (v1beta1).
func (src *VSphereVM) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereVM)
	if err := Convert_v1alpha3_VSphereVM_To_v1beta1_VSphereVM(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1.VSphereVM{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	dst.Spec.TagIDs = restored.Spec.TagIDs
	dst.Spec.AdditionalDisksGiB = restored.Spec.AdditionalDisksGiB
	dst.Spec.PowerOffMode = restored.Spec.PowerOffMode
	dst.Spec.GuestSoftPowerOffTimeout = restored.Spec.GuestSoftPowerOffTimeout
	dst.Status.Host = restored.Status.Host
	for i := range dst.Spec.Network.Devices {
		dst.Spec.Network.Devices[i].AddressesFromPools = restored.Spec.Network.Devices[i].AddressesFromPools
		dst.Spec.Network.Devices[i].DHCP4Overrides = restored.Spec.Network.Devices[i].DHCP4Overrides
		dst.Spec.Network.Devices[i].DHCP6Overrides = restored.Spec.Network.Devices[i].DHCP6Overrides
		dst.Spec.Network.Devices[i].SkipIPAllocation = restored.Spec.Network.Devices[i].SkipIPAllocation
	}
	dst.Spec.DataDisks = restored.Spec.DataDisks
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

// ConvertFrom converts from the Hub version (v1beta1) to this VSphereVM.
func (dst *VSphereVM) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereVM)
	if err := Convert_v1beta1_VSphereVM_To_v1alpha3_VSphereVM(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this VSphereVMList to the Hub version (v1beta1).
func (src *VSphereVMList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.VSphereVMList)
	return Convert_v1alpha3_VSphereVMList_To_v1beta1_VSphereVMList(src, dst, nil)
}

// ConvertFrom converts this VSphereVM to the Hub version (v1beta1).
func (dst *VSphereVMList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.VSphereVMList)
	return Convert_v1beta1_VSphereVMList_To_v1alpha3_VSphereVMList(src, dst, nil)
}
