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

package hub

// VirtualMachineNetworkInterfaceType specifies the NIC device type.
type VirtualMachineNetworkInterfaceType string

const (
	// VirtualMachineNetworkInterfaceTypeVMXNet3 specifies a VMXNet3 paravirtual
	// NIC. This is the default and recommended type for most workloads.
	VirtualMachineNetworkInterfaceTypeVMXNet3 VirtualMachineNetworkInterfaceType = "VMXNet3"

	// VirtualMachineNetworkInterfaceTypeSRIOV specifies an SR-IOV NIC.
	VirtualMachineNetworkInterfaceTypeSRIOV VirtualMachineNetworkInterfaceType = "SRIOV"

	// VirtualMachineNetworkInterfaceTypeE1000 specifies an Intel e1000 emulated
	// NIC.
	VirtualMachineNetworkInterfaceTypeE1000 VirtualMachineNetworkInterfaceType = "E1000"

	// VirtualMachineNetworkInterfaceTypeE1000e specifies an Intel e1000e
	// emulated NIC.
	VirtualMachineNetworkInterfaceTypeE1000e VirtualMachineNetworkInterfaceType = "E1000e"

	// VirtualMachineNetworkInterfaceTypeVMXNet2 specifies a VMXNet2
	// paravirtual NIC.
	VirtualMachineNetworkInterfaceTypeVMXNet2 VirtualMachineNetworkInterfaceType = "VMXNet2"

	// VirtualMachineNetworkInterfaceTypePCNet32 specifies an AMD PCnet32
	// emulated NIC.
	VirtualMachineNetworkInterfaceTypePCNet32 VirtualMachineNetworkInterfaceType = "PCNet32"
)

// TxContextThreadingMode specifies the transmit context threading mode for a
// VMXNet3 interface.
// This is a "weak enum": constants are well-known values; the field accepts any string for forward compatibility.
//
// +kubebuilder:validation:MaxLength=50
// +kubebuilder:validation:XValidation:rule="self.matches(\"^(PerDevice|PerVM|PerQueue|[1-9])$\")",message="must be PerDevice, PerVM, PerQueue, or a number (1-9)"
type TxContextThreadingMode string

const (
	// TxContextThreadingModePerDevice configures one TX thread per vNIC.
	TxContextThreadingModePerDevice TxContextThreadingMode = "PerDevice"

	// TxContextThreadingModePerVM configures one TX thread for the whole VM (default).
	TxContextThreadingModePerVM TxContextThreadingMode = "PerVM"

	// TxContextThreadingModePerQueue configures 2-8 TX threads per vNIC queue
	// (scheduler-determined). Recommended for 100G workloads with pnicFeatures
	// including ReceiveSideScaling.
	TxContextThreadingModePerQueue TxContextThreadingMode = "PerQueue"
)

// CoalescingScheme specifies the interrupt coalescing scheme for a VMXNet3
// interface.
// This is a "weak enum": constants are well-known values; the field accepts any string for forward compatibility.
//
// +kubebuilder:validation:XValidation:rule="self == 'Disabled' || self == 'Adapt' || self == 'Static' || self == 'RateBasedCoalescing' || size(self) < 128",message="must be Disabled, Adapt, Static, RateBasedCoalescing, or any other string shorter than 128 characters"
type CoalescingScheme string

const (
	// CoalescingSchemeDisabled disables interrupt coalescing entirely.
	// Recommended for latency-sensitive (LS=High) non-DPDK workloads because it
	// ensures each packet triggers an immediate interrupt.
	CoalescingSchemeDisabled CoalescingScheme = "Disabled"

	// CoalescingSchemeAdapt uses adaptive coalescing, dynamically adjusting the
	// interrupt rate based on VM and system load. CoalescingParams is ignored
	// when this scheme is set.
	CoalescingSchemeAdapt CoalescingScheme = "Adapt"

	// CoalescingSchemeStatic queues a fixed number of packets before triggering
	// an interrupt. CoalescingParams sets the Tx,Rx packet queue limit (range 1-64,
	// default "64").
	CoalescingSchemeStatic CoalescingScheme = "Static"

	// CoalescingSchemeRateBasedCoalescing uses rate-based coalescing (RBC)
	// CoalescingParams sets the interrupt rate in interrupts/sec (range 100-100000, default "4000").
	CoalescingSchemeRateBasedCoalescing CoalescingScheme = "RateBasedCoalescing"
)

// PNICQueueFeature names one physical NIC queue offload feature for VMXNet3
// pnicFeatures.
// This is a "weak enum": constants are well-known values; the field accepts any string for forward compatibility.
//
// +kubebuilder:validation:MaxLength=50
// +kubebuilder:validation:XValidation:rule="self.matches(\"^(LargeReceiveOffload|ReceiveSideScaling|^[1-9][0-9]*)$\")",message="must be LargeReceiveOffload, ReceiveSideScaling, or a non-empty decimal digits string for VMX pNICFeatures bitmask values (powers of two such as 1, 2, 4, 8)"
type PNICQueueFeature string

const (
	// PNICQueueFeatureLargeReceiveOffload enables large receive offload (LRO).
	PNICQueueFeatureLargeReceiveOffload PNICQueueFeature = "LargeReceiveOffload"

	// PNICQueueFeatureReceiveSideScaling enables receive-side scaling (RSS)
	// hardware queues, allowing the physical NIC to distribute incoming packets
	// across multiple receive queues. Typically set alongside
	// ctxPerDev=PerQueue for maximum throughput on 100G workloads.
	PNICQueueFeatureReceiveSideScaling PNICQueueFeature = "ReceiveSideScaling"
)

// UDPRSSMode controls whether UDP Receive Side Scaling is enabled on a VMXNet3
// interface.
type UDPRSSMode bool

const (
	// UDPRSSModeEnabled enables UDP RSS on the interface.
	UDPRSSModeEnabled UDPRSSMode = true

	// UDPRSSModeDisabled disables UDP RSS on the interface.
	UDPRSSModeDisabled UDPRSSMode = false
)

// VirtualMachineNetworkInterfaceVMXNet3Spec contains tuning options specific to
// VMXNet3 network interfaces. Fields with 'vmx' annotation map to ethernetX.* VMX keys, where X
// is the device index derived from the vSphere device key at runtime.
//
// These fields are only valid when the interface Type is VMXNet3. The CRD
// admission webhook rejects this struct when Type is set to an incompatible
// value.
//
// +kubebuilder:validation:XValidation:rule="!has(self.coalescingParams) || size(self.coalescingParams) < 128",message="coalescingParams must have length < 128"
type VirtualMachineNetworkInterfaceVMXNet3Spec struct {
	// +optional

	// UPTv2Enabled enables UPT v2 (Uniform Passthrough v2) for this interface.
	// UPT allows the guest to drive the physical NIC virtual function directly
	// via SR-IOV while preserving vMotion support by dynamically switching
	// between passthrough and emulation mode. UPTv1 is deprecated.
	//
	// Requires: spec.minHardwareVersion >= 20, SmartNIC with UPT support,
	// full VM memory reservation, and VMXNet3 v7 guest driver.
	UPTv2Enabled *bool `json:"uptv2Enabled,omitempty"`

	// +optional

	// CtxPerDev sets the TX context threading mode for this interface.
	// PerVM (default) gives one TX thread for the whole VM.
	// PerDevice gives one TX thread per vNIC.
	// PerQueue gives 2-8 TX threads per vNIC queue (scheduler-determined);
	// recommended for 100G workloads combined with pnicFeatures including ReceiveSideScaling.
	// Visible in esxtop as NetWorld-Dev-<name>-Tx threads.
	// Accepts known enum values (PerDevice, PerVM, PerQueue) or single digit integers (1-9)
	// for direct VMX values.
	CtxPerDev *TxContextThreadingMode `json:"ctxPerDev,omitempty" vmx:"ethernet%d.ctxPerDev"`

	// +optional

	// RSSOffloadEnabled enables RSS (Receive Side Scaling) offload, allowing
	// the physical NIC to distribute incoming packets across multiple receive
	// queues using a hardware-computed hash. Reduces hypervisor CPU overhead
	// and improves multi-core utilization for high-throughput workloads.
	// Requires pNIC RSS support.
	RSSOffloadEnabled *bool `json:"rssOffloadEnabled,omitempty" vmx:"ethernet%d.rssoffload"`

	// +optional

	// UDPRSSEnabled extends RSS to UDP traffic. By default RSS only distributes
	// TCP flows. Enabling this also distributes UDP flows, improving throughput
	// for UDP-heavy workloads such as GTP-U tunnels, QUIC, or media streaming.
	UDPRSSEnabled *UDPRSSMode `json:"udpRSSEnabled,omitempty" vmx:"ethernet%d.udpRSS"`

	// +optional
	// +kubebuilder:validation:MaxItems=16

	// +listType=set
	//
	// PNICFeatures lists physical NIC queue offload features to enable. The
	// primary use is including ReceiveSideScaling, which allows the vNIC to leverage physical
	// NIC RSS hardware queues. Typically set to ["ReceiveSideScaling"] alongside
	// ctxPerDev=PerQueue for maximum 100G throughput. Omitted or empty means no
	// extra pNIC queue features beyond defaults.
	// Accepts known enum values (LargeReceiveOffload, ReceiveSideScaling) or integer strings
	// representing powers of 2 (1,2,4,8,...) for direct VMX bitmask values.
	PNICFeatures []PNICQueueFeature `json:"pnicFeatures,omitempty" vmx:"ethernet%d.pnicfeatures"`

	// +optional

	// CoalescingScheme sets the interrupt coalescing scheme for this interface.
	// Use CoalescingSchemeDisabled for latency-sensitive (LS=High) non-DPDK
	// workloads to minimise interrupt latency.
	// Accepts known enum values (Disabled, Adapt, Static, RateBasedCoalescing) or any
	// string with length < 128 for forward compatibility.
	CoalescingScheme *CoalescingScheme `json:"coalescingScheme,omitempty" vmx:"ethernet%d.coalescingScheme"`

	// +optional

	// CoalescingParams sets the coalescing parameter when coalescingScheme is
	// RateBasedCoalescing or Static. The format depends on the scheme:
	//   - RateBasedCoalescing: integer string for interrupts/sec (e.g. "4000")
	//   - Static: integer string for packet queue limit (e.g. "64")
	// Ignored when coalescingScheme is Disabled or Adapt.
	// Must be length < 128 and valid 32-bit unsigned integer for RateBasedCoalescing.
	CoalescingParams *string `json:"coalescingParams,omitempty" vmx:"ethernet%d.coalescingParams"`
}
