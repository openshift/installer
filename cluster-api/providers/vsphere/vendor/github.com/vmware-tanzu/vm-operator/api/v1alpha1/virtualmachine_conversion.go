// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"fmt"
	"net"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	"github.com/vmware-tanzu/vm-operator/api/utilconversion"
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
)

const (
	// Well known device key used for the first disk.
	bootDiskDeviceKey = 2000
)

func Convert_v1alpha1_VirtualMachineVolume_To_v1alpha2_VirtualMachineVolume(
	in *VirtualMachineVolume, out *v1alpha2.VirtualMachineVolume, s apiconversion.Scope) error {

	if claim := in.PersistentVolumeClaim; claim != nil {
		out.PersistentVolumeClaim = &v1alpha2.PersistentVolumeClaimVolumeSource{
			PersistentVolumeClaimVolumeSource: claim.PersistentVolumeClaimVolumeSource,
		}

		if claim.InstanceVolumeClaim != nil {
			out.PersistentVolumeClaim.InstanceVolumeClaim = &v1alpha2.InstanceVolumeClaimVolumeSource{}

			if err := Convert_v1alpha1_InstanceVolumeClaimVolumeSource_To_v1alpha2_InstanceVolumeClaimVolumeSource(
				claim.InstanceVolumeClaim, out.PersistentVolumeClaim.InstanceVolumeClaim, s); err != nil {
				return err
			}
		}
	}

	// NOTE: in.VsphereVolume is dropped in v1a2. See filter_out_VirtualMachineVolumes_VsphereVolumes().

	return autoConvert_v1alpha1_VirtualMachineVolume_To_v1alpha2_VirtualMachineVolume(in, out, s)
}

func convert_v1alpha1_VirtualMachinePowerState_To_v1alpha2_VirtualMachinePowerState(
	in VirtualMachinePowerState) v1alpha2.VirtualMachinePowerState {

	switch in {
	case VirtualMachinePoweredOff:
		return v1alpha2.VirtualMachinePowerStateOff
	case VirtualMachinePoweredOn:
		return v1alpha2.VirtualMachinePowerStateOn
	case VirtualMachineSuspended:
		return v1alpha2.VirtualMachinePowerStateSuspended
	}

	return v1alpha2.VirtualMachinePowerState(in)
}

func convert_v1alpha2_VirtualMachinePowerState_To_v1alpha1_VirtualMachinePowerState(
	in v1alpha2.VirtualMachinePowerState) VirtualMachinePowerState {

	switch in {
	case v1alpha2.VirtualMachinePowerStateOff:
		return VirtualMachinePoweredOff
	case v1alpha2.VirtualMachinePowerStateOn:
		return VirtualMachinePoweredOn
	case v1alpha2.VirtualMachinePowerStateSuspended:
		return VirtualMachineSuspended
	}

	return VirtualMachinePowerState(in)
}

func convert_v1alpha1_VirtualMachinePowerOpMode_To_v1alpha2_VirtualMachinePowerOpMode(
	in VirtualMachinePowerOpMode) v1alpha2.VirtualMachinePowerOpMode {

	switch in {
	case VirtualMachinePowerOpModeHard:
		return v1alpha2.VirtualMachinePowerOpModeHard
	case VirtualMachinePowerOpModeSoft:
		return v1alpha2.VirtualMachinePowerOpModeSoft
	case VirtualMachinePowerOpModeTrySoft:
		return v1alpha2.VirtualMachinePowerOpModeTrySoft
	}

	return v1alpha2.VirtualMachinePowerOpMode(in)
}

func convert_v1alpha2_VirtualMachinePowerOpMode_To_v1alpha1_VirtualMachinePowerOpMode(
	in v1alpha2.VirtualMachinePowerOpMode) VirtualMachinePowerOpMode {

	switch in {
	case v1alpha2.VirtualMachinePowerOpModeHard:
		return VirtualMachinePowerOpModeHard
	case v1alpha2.VirtualMachinePowerOpModeSoft:
		return VirtualMachinePowerOpModeSoft
	case v1alpha2.VirtualMachinePowerOpModeTrySoft:
		return VirtualMachinePowerOpModeTrySoft
	}

	return VirtualMachinePowerOpMode(in)
}

func convert_v1alpha2_Conditions_To_v1alpha1_Phase(
	in []metav1.Condition) VMStatusPhase {

	// In practice, "Created" is the only really important value because some consumers
	// like CAPI use that as a part of their VM is-ready check.
	for _, c := range in {
		if c.Type == v1alpha2.VirtualMachineConditionCreated {
			if c.Status == metav1.ConditionTrue {
				return Created
			}
			return Creating
		}
	}

	return Unknown
}

func Convert_v1alpha2_VirtualMachineVolume_To_v1alpha1_VirtualMachineVolume(
	in *v1alpha2.VirtualMachineVolume, out *VirtualMachineVolume, s apiconversion.Scope) error {

	if claim := in.PersistentVolumeClaim; claim != nil {
		out.PersistentVolumeClaim = &PersistentVolumeClaimVolumeSource{
			PersistentVolumeClaimVolumeSource: claim.PersistentVolumeClaimVolumeSource,
		}

		if claim.InstanceVolumeClaim != nil {
			out.PersistentVolumeClaim.InstanceVolumeClaim = &InstanceVolumeClaimVolumeSource{}

			if err := Convert_v1alpha2_InstanceVolumeClaimVolumeSource_To_v1alpha1_InstanceVolumeClaimVolumeSource(
				claim.InstanceVolumeClaim, out.PersistentVolumeClaim.InstanceVolumeClaim, s); err != nil {
				return err
			}
		}
	}

	return autoConvert_v1alpha2_VirtualMachineVolume_To_v1alpha1_VirtualMachineVolume(in, out, s)
}

func convert_v1alpha1_VmMetadata_To_v1alpha2_BootstrapSpec(
	in *VirtualMachineMetadata) *v1alpha2.VirtualMachineBootstrapSpec {

	if in == nil || apiequality.Semantic.DeepEqual(*in, VirtualMachineMetadata{}) {
		return nil
	}

	out := v1alpha2.VirtualMachineBootstrapSpec{}

	objectName := in.SecretName
	if objectName == "" {
		objectName = in.ConfigMapName
	}

	switch in.Transport {
	case VirtualMachineMetadataExtraConfigTransport:
		out.LinuxPrep = &v1alpha2.VirtualMachineBootstrapLinuxPrepSpec{
			HardwareClockIsUTC: true,
		}
		out.CloudInit = &v1alpha2.VirtualMachineBootstrapCloudInitSpec{}
		if objectName != "" {
			out.CloudInit.RawCloudConfig = &common.SecretKeySelector{
				Name: objectName,
				Key:  "guestinfo.userdata", // TODO: Is this good enough? v1a1 would include everything with the "guestinfo" prefix
			}
		}
	case VirtualMachineMetadataOvfEnvTransport:
		out.LinuxPrep = &v1alpha2.VirtualMachineBootstrapLinuxPrepSpec{
			HardwareClockIsUTC: true,
		}
		out.VAppConfig = &v1alpha2.VirtualMachineBootstrapVAppConfigSpec{
			RawProperties: objectName,
		}
	case VirtualMachineMetadataVAppConfigTransport:
		out.VAppConfig = &v1alpha2.VirtualMachineBootstrapVAppConfigSpec{
			RawProperties: objectName,
		}
	case VirtualMachineMetadataCloudInitTransport:
		out.CloudInit = &v1alpha2.VirtualMachineBootstrapCloudInitSpec{}
		if objectName != "" {
			out.CloudInit.RawCloudConfig = &common.SecretKeySelector{
				Name: objectName,
				Key:  "user-data",
			}
		}
	case VirtualMachineMetadataSysprepTransport:
		out.Sysprep = &v1alpha2.VirtualMachineBootstrapSysprepSpec{}
		if objectName != "" {
			out.Sysprep.RawSysprep = &common.SecretKeySelector{
				Name: objectName,
				Key:  "unattend",
			}
		}
	}

	return &out
}

func convert_v1alpha2_BootstrapSpec_To_v1alpha1_VmMetadata(
	in *v1alpha2.VirtualMachineBootstrapSpec) *VirtualMachineMetadata {

	if in == nil || apiequality.Semantic.DeepEqual(*in, v1alpha2.VirtualMachineBootstrapSpec{}) {
		return nil
	}

	// TODO: v1a2 only has a Secret bootstrap field so that's what we set in v1a1. If this was created
	// as v1a1, we need to store the serialized object to know to set either the ConfigMap or Secret field.
	out := &VirtualMachineMetadata{}

	if cloudInit := in.CloudInit; cloudInit != nil {
		if cloudInit.RawCloudConfig != nil {
			out.SecretName = cloudInit.RawCloudConfig.Name

			switch cloudInit.RawCloudConfig.Key {
			case "guestinfo.userdata":
				out.Transport = VirtualMachineMetadataExtraConfigTransport
			case "user-data":
				out.Transport = VirtualMachineMetadataCloudInitTransport
			}
		} else if cloudInit.CloudConfig != nil {
			out.Transport = VirtualMachineMetadataCloudInitTransport
		}
	} else if sysprep := in.Sysprep; sysprep != nil {
		out.Transport = VirtualMachineMetadataSysprepTransport
		if in.Sysprep.RawSysprep != nil {
			out.SecretName = sysprep.RawSysprep.Name
		}
	} else if in.VAppConfig != nil {
		out.SecretName = in.VAppConfig.RawProperties

		if in.LinuxPrep != nil {
			out.Transport = VirtualMachineMetadataOvfEnvTransport
		} else {
			out.Transport = VirtualMachineMetadataVAppConfigTransport
		}
	}

	return out
}

func convert_v1alpha1_NetworkInterface_To_v1alpha2_NetworkInterfaceSpec(
	idx int, in VirtualMachineNetworkInterface) v1alpha2.VirtualMachineNetworkInterfaceSpec {

	out := v1alpha2.VirtualMachineNetworkInterfaceSpec{}
	out.Name = fmt.Sprintf("eth%d", idx)
	out.Network.Name = in.NetworkName

	switch in.NetworkType {
	case "vsphere-distributed":
		out.Network.TypeMeta.APIVersion = "netoperator.vmware.com/v1alpha1"
		out.Network.TypeMeta.Kind = "Network"
	case "nsx-t":
		out.Network.TypeMeta.APIVersion = "vmware.com/v1alpha1"
		out.Network.TypeMeta.Kind = "VirtualNetwork"
	}

	return out
}

func convert_v1alpha2_NetworkInterfaceSpec_To_v1alpha1_NetworkInterface(
	in v1alpha2.VirtualMachineNetworkInterfaceSpec) VirtualMachineNetworkInterface {

	out := VirtualMachineNetworkInterface{
		NetworkName: in.Network.Name,
	}

	switch in.Network.TypeMeta.Kind {
	case "Network":
		out.NetworkType = "vsphere-distributed"
	case "VirtualNetwork":
		out.NetworkType = "nsx-t"
	}

	return out
}

func Convert_v1alpha1_Probe_To_v1alpha2_VirtualMachineReadinessProbeSpec(in *Probe, out *v1alpha2.VirtualMachineReadinessProbeSpec, s apiconversion.Scope) error {
	probeSpec := convert_v1alpha1_Probe_To_v1alpha2_ReadinessProbeSpec(in)
	if probeSpec != nil {
		*out = *probeSpec
	}

	return nil
}

func convert_v1alpha1_Probe_To_v1alpha2_ReadinessProbeSpec(in *Probe) *v1alpha2.VirtualMachineReadinessProbeSpec {

	if in == nil || apiequality.Semantic.DeepEqual(*in, Probe{}) {
		return nil
	}

	out := v1alpha2.VirtualMachineReadinessProbeSpec{}

	out.TimeoutSeconds = in.TimeoutSeconds
	out.PeriodSeconds = in.PeriodSeconds

	if in.TCPSocket != nil {
		out.TCPSocket = &v1alpha2.TCPSocketAction{
			Port: in.TCPSocket.Port,
			Host: in.TCPSocket.Host,
		}
	}

	if in.GuestHeartbeat != nil {
		out.GuestHeartbeat = &v1alpha2.GuestHeartbeatAction{
			ThresholdStatus: v1alpha2.GuestHeartbeatStatus(in.GuestHeartbeat.ThresholdStatus),
		}
	}

	// out.GuestInfo =

	return &out
}

func Convert_v1alpha2_VirtualMachineReadinessProbeSpec_To_v1alpha1_Probe(in *v1alpha2.VirtualMachineReadinessProbeSpec, out *Probe, s apiconversion.Scope) error {
	probe := convert_v1alpha2_ReadinessProbeSpec_To_v1alpha1_Probe(in)
	if probe != nil {
		*out = *probe
	}

	return nil
}

func convert_v1alpha2_ReadinessProbeSpec_To_v1alpha1_Probe(in *v1alpha2.VirtualMachineReadinessProbeSpec) *Probe {

	if in == nil || apiequality.Semantic.DeepEqual(*in, v1alpha2.VirtualMachineReadinessProbeSpec{}) {
		return nil
	}

	out := &Probe{
		TimeoutSeconds: in.TimeoutSeconds,
		PeriodSeconds:  in.PeriodSeconds,
	}

	if in.TCPSocket != nil {
		out.TCPSocket = &TCPSocketAction{
			Port: in.TCPSocket.Port,
			Host: in.TCPSocket.Host,
		}
	}

	if in.GuestHeartbeat != nil {
		out.GuestHeartbeat = &GuestHeartbeatAction{
			ThresholdStatus: GuestHeartbeatStatus(in.GuestHeartbeat.ThresholdStatus),
		}
	}

	return out
}

func convert_v1alpha1_VirtualMachineAdvancedOptions_To_v1alpha2_VirtualMachineAdvancedSpec(
	in *VirtualMachineAdvancedOptions) *v1alpha2.VirtualMachineAdvancedSpec {

	if in == nil || apiequality.Semantic.DeepEqual(*in, VirtualMachineAdvancedOptions{}) {
		return nil
	}

	out := v1alpha2.VirtualMachineAdvancedSpec{}

	if opts := in.DefaultVolumeProvisioningOptions; opts != nil {
		if opts.ThinProvisioned != nil {
			if *opts.ThinProvisioned {
				out.DefaultVolumeProvisioningMode = v1alpha2.VirtualMachineVolumeProvisioningModeThin
			} else {
				out.DefaultVolumeProvisioningMode = v1alpha2.VirtualMachineVolumeProvisioningModeThick
			}
		} else if opts.EagerZeroed != nil && *opts.EagerZeroed {
			out.DefaultVolumeProvisioningMode = v1alpha2.VirtualMachineVolumeProvisioningModeThickEagerZero
		}
	}

	if in.ChangeBlockTracking != nil {
		out.ChangeBlockTracking = *in.ChangeBlockTracking
	}

	return &out
}

func convert_v1alpha1_VsphereVolumes_To_v1alpah2_BootDiskCapacity(volumes []VirtualMachineVolume) *resource.Quantity {
	// The v1a1 VsphereVolume was never a great API as you had to know the DeviceKey upfront; at the time our
	// API was private - only used by CAPW - and predates the "VM Service" VMs; In v1a2, we only support resizing
	// the boot disk via an explicit field. As good as we can here, map v1a1 volume into the v1a2 specific field.

	for i := range volumes {
		vsVol := volumes[i].VsphereVolume

		if vsVol != nil && vsVol.DeviceKey != nil && *vsVol.DeviceKey == bootDiskDeviceKey {
			// This VsphereVolume has the well-known boot disk device key. Return that capacity if set.
			if capacity := vsVol.Capacity.StorageEphemeral(); capacity != nil {
				return capacity
			}
			break
		}
	}

	return nil
}

func convert_v1alpha2_VirtualMachineAdvancedSpec_To_v1alpha1_VirtualMachineAdvancedOptions(
	in *v1alpha2.VirtualMachineAdvancedSpec) *VirtualMachineAdvancedOptions {

	if in == nil || apiequality.Semantic.DeepEqual(*in, v1alpha2.VirtualMachineAdvancedSpec{}) {
		return nil
	}

	out := &VirtualMachineAdvancedOptions{}

	if in.ChangeBlockTracking {
		out.ChangeBlockTracking = pointer.Bool(true)
	}

	switch in.DefaultVolumeProvisioningMode {
	case v1alpha2.VirtualMachineVolumeProvisioningModeThin:
		out.DefaultVolumeProvisioningOptions = &VirtualMachineVolumeProvisioningOptions{
			ThinProvisioned: pointer.Bool(true),
		}
	case v1alpha2.VirtualMachineVolumeProvisioningModeThick:
		out.DefaultVolumeProvisioningOptions = &VirtualMachineVolumeProvisioningOptions{
			ThinProvisioned: pointer.Bool(false),
		}
	case v1alpha2.VirtualMachineVolumeProvisioningModeThickEagerZero:
		out.DefaultVolumeProvisioningOptions = &VirtualMachineVolumeProvisioningOptions{
			EagerZeroed: pointer.Bool(true),
		}
	}

	if reflect.DeepEqual(out, &VirtualMachineAdvancedOptions{}) {
		return nil
	}

	return out
}

func convert_v1alpha2_BootDiskCapacity_To_v1alpha1_VirtualMachineVolume(capacity *resource.Quantity) *VirtualMachineVolume {
	if capacity == nil || capacity.IsZero() {
		return nil
	}

	const name = "vmoperator-vm-boot-disk"

	return &VirtualMachineVolume{
		Name: name,
		VsphereVolume: &VsphereVolumeSource{
			Capacity: corev1.ResourceList{
				corev1.ResourceEphemeralStorage: *capacity,
			},
			DeviceKey: pointer.Int(bootDiskDeviceKey),
		},
	}
}

func convert_v1alpha1_Network_To_v1alpha2_NetworkStatus(
	vmIP string, in []NetworkInterfaceStatus) *v1alpha2.VirtualMachineNetworkStatus {

	if vmIP == "" && len(in) == 0 {
		return nil
	}

	out := &v1alpha2.VirtualMachineNetworkStatus{}

	if net.ParseIP(vmIP).To4() != nil {
		out.PrimaryIP4 = vmIP
	} else {
		out.PrimaryIP6 = vmIP
	}

	ipAddrsToAddrStatus := func(ipAddr []string) []v1alpha2.VirtualMachineNetworkInterfaceIPAddrStatus {
		statuses := make([]v1alpha2.VirtualMachineNetworkInterfaceIPAddrStatus, 0, len(ipAddr))
		for _, ip := range ipAddr {
			statuses = append(statuses, v1alpha2.VirtualMachineNetworkInterfaceIPAddrStatus{Address: ip})
		}
		return statuses
	}

	for _, inI := range in {
		interfaceStatus := v1alpha2.VirtualMachineNetworkInterfaceStatus{
			IP: v1alpha2.VirtualMachineNetworkInterfaceIPStatus{
				Addresses: ipAddrsToAddrStatus(inI.IpAddresses),
				MACAddr:   inI.MacAddress,
			},
		}
		out.Interfaces = append(out.Interfaces, interfaceStatus)
	}

	return out
}

func convert_v1alpha2_NetworkStatus_To_v1alpha1_Network(
	in *v1alpha2.VirtualMachineNetworkStatus) (string, []NetworkInterfaceStatus) {

	if in == nil {
		return "", nil
	}

	vmIP := in.PrimaryIP4
	if vmIP == "" {
		vmIP = in.PrimaryIP6
	}

	addrStatusToIPAddrs := func(addrStatus []v1alpha2.VirtualMachineNetworkInterfaceIPAddrStatus) []string {
		ipAddrs := make([]string, 0, len(addrStatus))
		for _, a := range addrStatus {
			ipAddrs = append(ipAddrs, a.Address)
		}
		return ipAddrs
	}

	out := make([]NetworkInterfaceStatus, 0, len(in.Interfaces))
	for _, i := range in.Interfaces {
		interfaceStatus := NetworkInterfaceStatus{
			Connected:   true,
			MacAddress:  i.IP.MACAddr,
			IpAddresses: addrStatusToIPAddrs(i.IP.Addresses),
		}
		out = append(out, interfaceStatus)
	}

	return vmIP, out
}

// In v1a2 we've dropped the v1a1 VsphereVolumes, and in its place we have a single field for the boot
// disk size. The Convert_v1alpha1_VirtualMachineVolume_To_v1alpha2_VirtualMachineVolume() stub does not
// allow us to not return something so filter those volumes - without a PersistentVolumeClaim set - here.
func filter_out_VirtualMachineVolumes_VsphereVolumes(in []v1alpha2.VirtualMachineVolume) []v1alpha2.VirtualMachineVolume {

	if len(in) == 0 {
		return nil
	}

	out := make([]v1alpha2.VirtualMachineVolume, 0, len(in))

	for _, v := range in {
		if v.PersistentVolumeClaim != nil {
			out = append(out, v)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}

func Convert_v1alpha1_VirtualMachineSpec_To_v1alpha2_VirtualMachineSpec(
	in *VirtualMachineSpec, out *v1alpha2.VirtualMachineSpec, s apiconversion.Scope) error {

	// The generated auto convert will convert the power modes as-is strings which breaks things, so keep
	// this first.
	if err := autoConvert_v1alpha1_VirtualMachineSpec_To_v1alpha2_VirtualMachineSpec(in, out, s); err != nil {
		return err
	}

	out.PowerState = convert_v1alpha1_VirtualMachinePowerState_To_v1alpha2_VirtualMachinePowerState(in.PowerState)
	out.PowerOffMode = convert_v1alpha1_VirtualMachinePowerOpMode_To_v1alpha2_VirtualMachinePowerOpMode(in.PowerOffMode)
	out.SuspendMode = convert_v1alpha1_VirtualMachinePowerOpMode_To_v1alpha2_VirtualMachinePowerOpMode(in.SuspendMode)
	out.NextRestartTime = in.NextRestartTime
	out.RestartMode = convert_v1alpha1_VirtualMachinePowerOpMode_To_v1alpha2_VirtualMachinePowerOpMode(in.RestartMode)
	out.Bootstrap = convert_v1alpha1_VmMetadata_To_v1alpha2_BootstrapSpec(in.VmMetadata)
	out.Volumes = filter_out_VirtualMachineVolumes_VsphereVolumes(out.Volumes)

	if len(in.NetworkInterfaces) > 0 {
		out.Network = &v1alpha2.VirtualMachineNetworkSpec{}
		for i, networkInterface := range in.NetworkInterfaces {
			networkInterfaceSpec := convert_v1alpha1_NetworkInterface_To_v1alpha2_NetworkInterfaceSpec(i, networkInterface)
			out.Network.Interfaces = append(out.Network.Interfaces, networkInterfaceSpec)
		}
	}

	out.ReadinessProbe = convert_v1alpha1_Probe_To_v1alpha2_ReadinessProbeSpec(in.ReadinessProbe)
	out.Advanced = convert_v1alpha1_VirtualMachineAdvancedOptions_To_v1alpha2_VirtualMachineAdvancedSpec(in.AdvancedOptions)
	if out.Advanced != nil {
		out.Advanced.BootDiskCapacity = convert_v1alpha1_VsphereVolumes_To_v1alpah2_BootDiskCapacity(in.Volumes)
	}

	if in.ResourcePolicyName != "" {
		if out.Reserved == nil {
			out.Reserved = &v1alpha2.VirtualMachineReservedSpec{}
		}
		out.Reserved.ResourcePolicyName = in.ResourcePolicyName
	}

	// Deprecated:
	// in.Ports

	return nil
}

func Convert_v1alpha2_VirtualMachineSpec_To_v1alpha1_VirtualMachineSpec(
	in *v1alpha2.VirtualMachineSpec, out *VirtualMachineSpec, s apiconversion.Scope) error {

	if err := autoConvert_v1alpha2_VirtualMachineSpec_To_v1alpha1_VirtualMachineSpec(in, out, s); err != nil {
		return err
	}

	out.PowerState = convert_v1alpha2_VirtualMachinePowerState_To_v1alpha1_VirtualMachinePowerState(in.PowerState)
	out.PowerOffMode = convert_v1alpha2_VirtualMachinePowerOpMode_To_v1alpha1_VirtualMachinePowerOpMode(in.PowerOffMode)
	out.SuspendMode = convert_v1alpha2_VirtualMachinePowerOpMode_To_v1alpha1_VirtualMachinePowerOpMode(in.SuspendMode)
	out.NextRestartTime = in.NextRestartTime
	out.RestartMode = convert_v1alpha2_VirtualMachinePowerOpMode_To_v1alpha1_VirtualMachinePowerOpMode(in.RestartMode)
	out.VmMetadata = convert_v1alpha2_BootstrapSpec_To_v1alpha1_VmMetadata(in.Bootstrap)

	if in.Network != nil {
		for _, networkInterfaceSpec := range in.Network.Interfaces {
			networkInterface := convert_v1alpha2_NetworkInterfaceSpec_To_v1alpha1_NetworkInterface(networkInterfaceSpec)
			out.NetworkInterfaces = append(out.NetworkInterfaces, networkInterface)
		}
	}

	out.ReadinessProbe = convert_v1alpha2_ReadinessProbeSpec_To_v1alpha1_Probe(in.ReadinessProbe)
	out.AdvancedOptions = convert_v1alpha2_VirtualMachineAdvancedSpec_To_v1alpha1_VirtualMachineAdvancedOptions(in.Advanced)

	if in.Reserved != nil {
		out.ResourcePolicyName = in.Reserved.ResourcePolicyName
	}

	if in.Advanced != nil {
		if bootDiskVol := convert_v1alpha2_BootDiskCapacity_To_v1alpha1_VirtualMachineVolume(in.Advanced.BootDiskCapacity); bootDiskVol != nil {
			out.Volumes = append(out.Volumes, *bootDiskVol)
		}
	}

	// TODO = in.ReadinessGates

	// Deprecated:
	// out.Ports

	return nil
}

func Convert_v1alpha1_VirtualMachineVolumeStatus_To_v1alpha2_VirtualMachineVolumeStatus(
	in *VirtualMachineVolumeStatus, out *v1alpha2.VirtualMachineVolumeStatus, s apiconversion.Scope) error {

	out.DiskUUID = in.DiskUuid

	return autoConvert_v1alpha1_VirtualMachineVolumeStatus_To_v1alpha2_VirtualMachineVolumeStatus(in, out, s)
}

func Convert_v1alpha2_VirtualMachineVolumeStatus_To_v1alpha1_VirtualMachineVolumeStatus(
	in *v1alpha2.VirtualMachineVolumeStatus, out *VirtualMachineVolumeStatus, s apiconversion.Scope) error {

	out.DiskUuid = in.DiskUUID

	return autoConvert_v1alpha2_VirtualMachineVolumeStatus_To_v1alpha1_VirtualMachineVolumeStatus(in, out, s)
}

func Convert_v1alpha1_VirtualMachineStatus_To_v1alpha2_VirtualMachineStatus(
	in *VirtualMachineStatus, out *v1alpha2.VirtualMachineStatus, s apiconversion.Scope) error {

	if err := autoConvert_v1alpha1_VirtualMachineStatus_To_v1alpha2_VirtualMachineStatus(in, out, s); err != nil {
		return err
	}

	out.PowerState = convert_v1alpha1_VirtualMachinePowerState_To_v1alpha2_VirtualMachinePowerState(in.PowerState)
	out.Network = convert_v1alpha1_Network_To_v1alpha2_NetworkStatus(in.VmIp, in.NetworkInterfaces)
	out.LastRestartTime = in.LastRestartTime

	// WARNING: in.Phase requires manual conversion: does not exist in peer-type

	return nil
}

func translate_v1alpha2_Conditions_To_v1alpha1_Conditions(conditions []Condition) []Condition {
	var preReqCond, vmClassCond, vmImageCond, vmSetResourcePolicy, vmBootstrap *Condition

	for i := range conditions {
		c := &conditions[i]

		switch c.Type {
		case VirtualMachinePrereqReadyCondition:
			preReqCond = c
		case v1alpha2.VirtualMachineConditionClassReady:
			vmClassCond = c
		case v1alpha2.VirtualMachineConditionImageReady:
			vmImageCond = c
		case v1alpha2.VirtualMachineConditionVMSetResourcePolicyReady:
			vmSetResourcePolicy = c
		case v1alpha2.VirtualMachineConditionBootstrapReady:
			vmBootstrap = c
		}
	}

	// Try to replicate how the v1a1 provider would set the singular prereqs condition. The class is checked
	// first, then the image. Note that the set resource policy and metadata (bootstrap) are not a part of
	// the v1a1 prereqs, and are optional.
	if vmClassCond != nil && vmClassCond.Status == corev1.ConditionTrue &&
		vmImageCond != nil && vmImageCond.Status == corev1.ConditionTrue &&
		(vmSetResourcePolicy == nil || vmSetResourcePolicy.Status == corev1.ConditionTrue) &&
		(vmBootstrap == nil || vmBootstrap.Status == corev1.ConditionTrue) {

		p := Condition{
			Type:   VirtualMachinePrereqReadyCondition,
			Status: corev1.ConditionTrue,
		}

		if preReqCond != nil {
			p.LastTransitionTime = preReqCond.LastTransitionTime
			*preReqCond = p
			return conditions
		}

		p.LastTransitionTime = vmImageCond.LastTransitionTime
		return append(conditions, p)
	}

	p := Condition{
		Type:     VirtualMachinePrereqReadyCondition,
		Status:   corev1.ConditionFalse,
		Severity: ConditionSeverityError,
	}

	if vmClassCond != nil && vmClassCond.Status == corev1.ConditionFalse {
		p.Reason = VirtualMachineClassNotFoundReason
		p.Message = vmClassCond.Message
		p.LastTransitionTime = vmClassCond.LastTransitionTime
	} else if vmImageCond != nil && vmImageCond.Status == corev1.ConditionFalse {
		p.Reason = VirtualMachineImageNotFoundReason
		p.Message = vmImageCond.Message
		p.LastTransitionTime = vmImageCond.LastTransitionTime
	}

	if p.Reason != "" {
		if preReqCond != nil {
			*preReqCond = p
			return conditions
		}

		return append(conditions, p)
	}

	if vmSetResourcePolicy != nil && vmSetResourcePolicy.Status == corev1.ConditionFalse &&
		vmBootstrap != nil && vmBootstrap.Status == corev1.ConditionFalse {

		// These are not a part of the v1a1 Prereqs. If either is false, the v1a1 provider would not
		// update the prereqs condition, but we don't set the condition to true either until all these
		// conditions are true. Just leave things as is to see how strict we really need to be here.
		return conditions
	}

	// TBD: For now, leave the v1a2 conditions if present since those provide more details.

	return conditions
}

func Convert_v1alpha2_VirtualMachineStatus_To_v1alpha1_VirtualMachineStatus(
	in *v1alpha2.VirtualMachineStatus, out *VirtualMachineStatus, s apiconversion.Scope) error {

	if err := autoConvert_v1alpha2_VirtualMachineStatus_To_v1alpha1_VirtualMachineStatus(in, out, s); err != nil {
		return err
	}

	out.PowerState = convert_v1alpha2_VirtualMachinePowerState_To_v1alpha1_VirtualMachinePowerState(in.PowerState)
	out.Phase = convert_v1alpha2_Conditions_To_v1alpha1_Phase(in.Conditions)
	out.VmIp, out.NetworkInterfaces = convert_v1alpha2_NetworkStatus_To_v1alpha1_Network(in.Network)
	out.LastRestartTime = in.LastRestartTime
	out.Conditions = translate_v1alpha2_Conditions_To_v1alpha1_Conditions(out.Conditions)

	// WARNING: in.Image requires manual conversion: does not exist in peer-type
	// WARNING: in.Class requires manual conversion: does not exist in peer-type

	return nil
}

func restore_v1alpha2_VirtualMachineBootstrapSpec(
	dst, src *v1alpha2.VirtualMachine) {

	dstBootstrap := dst.Spec.Bootstrap
	srcBootstrap := src.Spec.Bootstrap

	if dstBootstrap == nil || srcBootstrap == nil {
		return
	}

	mergeSecretKeySelector := func(dstSel, srcSel *common.SecretKeySelector) *common.SecretKeySelector {
		if dstSel == nil || srcSel == nil {
			return dstSel
		}

		// Restore with the new object name in case it was changed.
		newSel := *srcSel
		newSel.Name = dstSel.Name
		return &newSel
	}

	if dstCloudInit := dstBootstrap.CloudInit; dstCloudInit != nil {
		if srcCloudInit := srcBootstrap.CloudInit; srcCloudInit != nil {
			dstCloudInit.CloudConfig = srcCloudInit.CloudConfig
			dstCloudInit.RawCloudConfig = mergeSecretKeySelector(dstCloudInit.RawCloudConfig, srcCloudInit.RawCloudConfig)
			dstCloudInit.SSHAuthorizedKeys = srcCloudInit.SSHAuthorizedKeys
		}
	}

	if dstLinuxPrep := dstBootstrap.LinuxPrep; dstLinuxPrep != nil {
		if srcLinuxPrep := srcBootstrap.LinuxPrep; srcLinuxPrep != nil {
			dstLinuxPrep.HardwareClockIsUTC = srcLinuxPrep.HardwareClockIsUTC
			dstLinuxPrep.TimeZone = srcLinuxPrep.TimeZone
		}
	}

	if dstSysPrep := dstBootstrap.Sysprep; dstSysPrep != nil {
		if srcSysPrep := srcBootstrap.Sysprep; srcSysPrep != nil {
			dstSysPrep.Sysprep = srcSysPrep.Sysprep
			dstSysPrep.RawSysprep = mergeSecretKeySelector(srcSysPrep.RawSysprep, dstSysPrep.RawSysprep)
		}
	}

	if dstVAppConfig := dstBootstrap.VAppConfig; dstVAppConfig != nil {
		if srcVAppConfig := srcBootstrap.VAppConfig; srcVAppConfig != nil {
			dstVAppConfig.Properties = srcVAppConfig.Properties
			dstVAppConfig.RawProperties = srcVAppConfig.RawProperties
		}
	}
}

func restore_v1alpha2_VirtualMachineNetwork(
	dst, src *v1alpha2.VirtualMachine) {

	dstNetwork := dst.Spec.Network
	srcNetwork := src.Spec.Network

	if dstNetwork == nil || srcNetwork == nil {
		return
	}

	dstNetwork.HostName = srcNetwork.HostName
	dstNetwork.Disabled = srcNetwork.Disabled

	if len(dstNetwork.Interfaces) == 0 {
		// No interfaces so nothing to fixup (the interfaces were removed): we ignore the restored interfaces.
		return
	}

	// The exceedingly common case is there was and still is just one interface, and for the same
	// network. With multiple interfaces it gets much harder to line things up. We really just
	// don't have info in v1a1 so this is a little best effort. The API supports it, but  we really
	// never supported-supported it yet.
	//
	// Do the about the easiest thing and zip the interfaces together, using the network name to
	// determine if it is the "same" interface. v1a1 could not support multiple interfaces on the
	// same network, and it probably won't be common place in v1a2, so we could try to refine this
	// a little more later.
	for i := range dstNetwork.Interfaces {
		if i >= len(srcNetwork.Interfaces) {
			break
		}

		if dstNetwork.Interfaces[i].Network.Name == srcNetwork.Interfaces[i].Network.Name {
			dstNetwork.Interfaces[i] = srcNetwork.Interfaces[i]
		}
	}
}

// ConvertTo converts this VirtualMachine to the Hub version.
func (src *VirtualMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachine)
	if err := Convert_v1alpha1_VirtualMachine_To_v1alpha2_VirtualMachine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &v1alpha2.VirtualMachine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	restore_v1alpha2_VirtualMachineBootstrapSpec(dst, restored)
	restore_v1alpha2_VirtualMachineNetwork(dst, restored)

	if restored.Spec.ReadinessProbe != nil {
		if dst.Spec.ReadinessProbe == nil {
			dst.Spec.ReadinessProbe = &v1alpha2.VirtualMachineReadinessProbeSpec{}
		}
		dst.Spec.ReadinessProbe.GuestInfo = restored.Spec.ReadinessProbe.GuestInfo
	}

	return nil
}

// ConvertFrom converts the hub version to this VirtualMachine.
func (dst *VirtualMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachine)
	if err := Convert_v1alpha2_VirtualMachine_To_v1alpha1_VirtualMachine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	return utilconversion.MarshalData(src, dst)
}

// ConvertTo converts this VirtualMachineList to the Hub version.
func (src *VirtualMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachineList)
	return Convert_v1alpha1_VirtualMachineList_To_v1alpha2_VirtualMachineList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineList.
func (dst *VirtualMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachineList)
	return Convert_v1alpha2_VirtualMachineList_To_v1alpha1_VirtualMachineList(src, dst, nil)
}
