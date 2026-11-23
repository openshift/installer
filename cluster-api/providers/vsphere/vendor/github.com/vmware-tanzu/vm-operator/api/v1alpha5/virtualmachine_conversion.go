// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha5

import (
	corev1 "k8s.io/api/core/v1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/vmware-tanzu/vm-operator/api/utilconversion"

	vmopv1 "github.com/vmware-tanzu/vm-operator/api/v1alpha6"
)

func Convert_v1alpha6_VirtualMachineNetworkSpec_To_v1alpha5_VirtualMachineNetworkSpec(
	in *vmopv1.VirtualMachineNetworkSpec, out *VirtualMachineNetworkSpec, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineNetworkSpec_To_v1alpha5_VirtualMachineNetworkSpec(in, out, s)
}

// Convert_v1alpha6_VirtualMachineStatus_To_v1alpha5_VirtualMachineStatus drops fields that do
// not exist in v1alpha5; they are fully restored via dst.Status = restored.Status in ConvertTo.
func Convert_v1alpha6_VirtualMachineStatus_To_v1alpha5_VirtualMachineStatus(
	in *vmopv1.VirtualMachineStatus, out *VirtualMachineStatus, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineStatus_To_v1alpha5_VirtualMachineStatus(in, out, s)
}

// Convert_v1alpha6_VirtualMachineAdvancedSpec_To_v1alpha5_VirtualMachineAdvancedSpec drops
// fields that do not exist in v1alpha5; they are preserved via MarshalData on ConvertFrom.
func Convert_v1alpha6_VirtualMachineAdvancedSpec_To_v1alpha5_VirtualMachineAdvancedSpec(
	in *vmopv1.VirtualMachineAdvancedSpec, out *VirtualMachineAdvancedSpec, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineAdvancedSpec_To_v1alpha5_VirtualMachineAdvancedSpec(in, out, s)
}

// Convert_v1alpha6_VirtualMachineNetworkInterfaceSpec_To_v1alpha5_VirtualMachineNetworkInterfaceSpec drops
// fields that do not exist in v1alpha5; they are preserved via MarshalData on ConvertFrom.
func Convert_v1alpha6_VirtualMachineNetworkInterfaceSpec_To_v1alpha5_VirtualMachineNetworkInterfaceSpec(
	in *vmopv1.VirtualMachineNetworkInterfaceSpec, out *VirtualMachineNetworkInterfaceSpec, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineNetworkInterfaceSpec_To_v1alpha5_VirtualMachineNetworkInterfaceSpec(in, out, s)
}

func restore_v1alpha6_VirtualMachineAdvanced(dst, src *vmopv1.VirtualMachine) {
	adv := src.Spec.Advanced
	if adv == nil {
		return
	}

	if dst.Spec.Advanced == nil {
		dst.Spec.Advanced = &vmopv1.VirtualMachineAdvancedSpec{}
	}
	dst.Spec.Advanced.PreferHTEnabled = adv.PreferHTEnabled
	dst.Spec.Advanced.HugePages1GEnabled = adv.HugePages1GEnabled
	dst.Spec.Advanced.TimeTrackerLowLatencyEnabled = adv.TimeTrackerLowLatencyEnabled
	dst.Spec.Advanced.CPUAffinityExclusiveNoStatsEnabled = adv.CPUAffinityExclusiveNoStatsEnabled
	dst.Spec.Advanced.VMXSwapEnabled = adv.VMXSwapEnabled
	dst.Spec.Advanced.PNUMANodeAffinity = adv.PNUMANodeAffinity
	dst.Spec.Advanced.ExtraConfig = adv.ExtraConfig
}

func restore_v1alpha6_VirtualMachineNetworkInterfaces(dst, src *vmopv1.VirtualMachine) {
	if src.Spec.Network == nil || len(src.Spec.Network.Interfaces) == 0 {
		return
	}
	if dst.Spec.Network == nil || len(dst.Spec.Network.Interfaces) == 0 {
		return
	}
	srcByName := make(map[string]*vmopv1.VirtualMachineNetworkInterfaceSpec, len(src.Spec.Network.Interfaces))
	for i := range src.Spec.Network.Interfaces {
		srcByName[src.Spec.Network.Interfaces[i].Name] = &src.Spec.Network.Interfaces[i]
	}
	for i := range dst.Spec.Network.Interfaces {
		srcIface, ok := srcByName[dst.Spec.Network.Interfaces[i].Name]
		if !ok {
			continue
		}
		dstIface := &dst.Spec.Network.Interfaces[i]
		dstIface.IPAMModes = append([]corev1.IPFamily(nil), srcIface.IPAMModes...)
		dstIface.Type = srcIface.Type
		dstIface.VNUMANodeID = srcIface.VNUMANodeID
		dstIface.VMXNet3 = srcIface.VMXNet3
		dstIface.AdvancedProperties = srcIface.AdvancedProperties
	}
}

// Convert_v1alpha6_VirtualMachineCPUAllocationStatus_To_v1alpha5_VirtualMachineCPUAllocationStatus drops
// v1alpha6 fields that do not exist in v1alpha5; they are preserved via MarshalData on ConvertFrom.
func Convert_v1alpha6_VirtualMachineCPUAllocationStatus_To_v1alpha5_VirtualMachineCPUAllocationStatus(
	in *vmopv1.VirtualMachineCPUAllocationStatus, out *VirtualMachineCPUAllocationStatus, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineCPUAllocationStatus_To_v1alpha5_VirtualMachineCPUAllocationStatus(in, out, s)
}

// Convert_v1alpha6_VirtualMachineMemoryAllocationStatus_To_v1alpha5_VirtualMachineMemoryAllocationStatus drops
// v1alpha6 fields that do not exist in v1alpha5; they are preserved via MarshalData on ConvertFrom.
func Convert_v1alpha6_VirtualMachineMemoryAllocationStatus_To_v1alpha5_VirtualMachineMemoryAllocationStatus(
	in *vmopv1.VirtualMachineMemoryAllocationStatus, out *VirtualMachineMemoryAllocationStatus, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineMemoryAllocationStatus_To_v1alpha5_VirtualMachineMemoryAllocationStatus(in, out, s)
}

func Convert_v1alpha6_VirtualMachineBootstrapSpec_To_v1alpha5_VirtualMachineBootstrapSpec(
	in *vmopv1.VirtualMachineBootstrapSpec, out *VirtualMachineBootstrapSpec, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineBootstrapSpec_To_v1alpha5_VirtualMachineBootstrapSpec(in, out, s)
}

func Convert_v1alpha6_VirtualMachineSpec_To_v1alpha5_VirtualMachineSpec(
	in *vmopv1.VirtualMachineSpec, out *VirtualMachineSpec, s apiconversion.Scope) error {

	return autoConvert_v1alpha6_VirtualMachineSpec_To_v1alpha5_VirtualMachineSpec(in, out, s)
}

func restore_v1alpha6_VirtualMachineBootstrapDisabled(dst, src *vmopv1.VirtualMachine) {
	if bs := src.Spec.Bootstrap; bs != nil {
		if bs.Disabled {
			if dst.Spec.Bootstrap == nil {
				dst.Spec.Bootstrap = &vmopv1.VirtualMachineBootstrapSpec{}
			}
			dst.Spec.Bootstrap.Disabled = true
		}
	}
}

func restore_v1alpha6_VirtualMachineVolumeAttributesClassName(dst, src *vmopv1.VirtualMachine) {
	if src.Spec.VolumeAttributesClassName != "" {
		dst.Spec.VolumeAttributesClassName = src.Spec.VolumeAttributesClassName
	}
}

func restore_v1alpha6_VirtualMachineNetworkVLANs(dst, src *vmopv1.VirtualMachine) {
	if src.Spec.Network == nil || len(src.Spec.Network.VLANs) == 0 {
		return
	}

	if dst.Spec.Network == nil {
		dst.Spec.Network = &vmopv1.VirtualMachineNetworkSpec{}
	}
	dst.Spec.Network.VLANs = src.Spec.Network.VLANs
}

func restore_v1alpha6_VirtualMachineResources(dst, src *vmopv1.VirtualMachine) {
	dst.Spec.Resources = src.Spec.Resources
}

func restore_v1alpha6_VirtualMachineCPUAdvanced(dst, src *vmopv1.VirtualMachine) {
	dst.Spec.CPUAdvanced = src.Spec.CPUAdvanced
}

func restore_v1alpha6_VirtualMachineMemoryAdvanced(dst, src *vmopv1.VirtualMachine) {
	dst.Spec.MemoryAdvanced = src.Spec.MemoryAdvanced
}

// ConvertTo converts this VirtualMachine to the Hub version.
func (src *VirtualMachine) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachine)
	if err := Convert_v1alpha5_VirtualMachine_To_v1alpha6_VirtualMachine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &vmopv1.VirtualMachine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	// BEGIN RESTORE

	restore_v1alpha6_VirtualMachineBootstrapDisabled(dst, restored)
	restore_v1alpha6_VirtualMachineVolumeAttributesClassName(dst, restored)
	restore_v1alpha6_VirtualMachineNetworkVLANs(dst, restored)
	restore_v1alpha6_VirtualMachineAdvanced(dst, restored)
	restore_v1alpha6_VirtualMachineNetworkInterfaces(dst, restored)
	restore_v1alpha6_VirtualMachineResources(dst, restored)
	restore_v1alpha6_VirtualMachineCPUAdvanced(dst, restored)
	restore_v1alpha6_VirtualMachineMemoryAdvanced(dst, restored)

	// END RESTORE

	dst.Status = restored.Status

	return nil
}

// ConvertFrom converts the hub version to this VirtualMachine.
func (dst *VirtualMachine) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachine)
	if err := Convert_v1alpha6_VirtualMachine_To_v1alpha5_VirtualMachine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	return utilconversion.MarshalData(src, dst)
}

// ConvertTo converts this VirtualMachineList to the Hub version.
func (src *VirtualMachineList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*vmopv1.VirtualMachineList)
	return Convert_v1alpha5_VirtualMachineList_To_v1alpha6_VirtualMachineList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineList.
func (dst *VirtualMachineList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*vmopv1.VirtualMachineList)
	return Convert_v1alpha6_VirtualMachineList_To_v1alpha5_VirtualMachineList(src, dst, nil)
}
