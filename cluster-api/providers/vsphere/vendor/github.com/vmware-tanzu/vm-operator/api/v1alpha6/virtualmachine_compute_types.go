// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

import "k8s.io/apimachinery/pkg/api/resource"

// VirtualMachineResourceQuantity holds a pair of CPU and memory resource
// quantities. It is used for spec.resources.size, spec.resources.requests,
// and spec.resources.limits.
//
// Requires the TelcoVMServiceAPI supervisor capability.
type VirtualMachineResourceQuantity struct {
	// +optional
	// +kubebuilder:validation:XValidation:rule="type(self) != string || !self.endsWith('m')",message="CPU must be a whole number (e.g. '4' for vCPUs or '2000' for 2000 MHz). The 'm' (milli) suffix is not supported."

	// CPU is a CPU resource quantity.
	//
	// For spec.resources.size, the value represents the guest-visible vCPU
	// count and must be a whole number (e.g. "4"). Maps to
	// ConfigSpec.NumCPUs.
	//
	// For spec.resources.requests and spec.resources.limits, the value
	// represents a host-level CPU allocation in MHz (e.g. "2000" for
	// 2 GHz). Maps to CpuAllocation.Reservation and CpuAllocation.Limit
	// respectively.
	CPU *resource.Quantity `json:"cpu,omitempty"`

	// +optional

	// Memory is a memory resource quantity in bytes (e.g. "8Gi").
	//
	// For spec.resources.size, maps to ConfigSpec.MemoryMB (guest-visible).
	// For spec.resources.requests and spec.resources.limits, maps to
	// MemoryAllocation.Reservation and MemoryAllocation.Limit respectively.
	Memory *resource.Quantity `json:"memory,omitempty"`
}

// VirtualMachineResourcesSpec describes the desired compute resource
// allocation for a VirtualMachine. All sub-fields are optional; a nil
// sub-field defers to the VirtualMachineClass value via field-level merge.
// When no VirtualMachineClass is referenced, spec.resources.size.cpu and
// spec.resources.size.memory are required.
//
// Requires the TelcoVMServiceAPI supervisor capability.
type VirtualMachineResourcesSpec struct {
	// +optional

	// Size is the guest-visible compute allocation.
	// size.cpu is the vCPU count; size.memory is the guest RAM.
	// Maps to ConfigSpec.NumCPUs and ConfigSpec.MemoryMB.
	// When set, overrides the corresponding VMClass hardware values.
	Size *VirtualMachineResourceQuantity `json:"size,omitempty"`

	// +optional

	// Requests is the host-level resource reservation (host guarantee).
	// requests.cpu is in MHz; requests.memory is in bytes.
	// Maps to CpuAllocation.Reservation and MemoryAllocation.Reservation.
	// Can be reconfigured while the VM is powered on.
	Requests *VirtualMachineResourceQuantity `json:"requests,omitempty"`

	// +optional

	// Limits is the host-level resource allocation ceiling.
	// limits.cpu is in MHz (nil = unlimited); limits.memory is in bytes
	// (nil = unlimited).
	// Maps to CpuAllocation.Limit and MemoryAllocation.Limit.
	// Can be reconfigured while the VM is powered on.
	Limits *VirtualMachineResourceQuantity `json:"limits,omitempty"`
}

// +kubebuilder:validation:Enum=Low;Normal;High;HighWithHyperthreading

// VirtualMachineLatencySensitivityLevel defines the vSphere CPU scheduler
// latency sensitivity level for a VirtualMachine.
type VirtualMachineLatencySensitivityLevel string

const (
	// VirtualMachineLatencySensitivityNormal is the default scheduling mode.
	VirtualMachineLatencySensitivityNormal VirtualMachineLatencySensitivityLevel = "Normal"

	// VirtualMachineLatencySensitivityHigh minimizes scheduling latency at
	// the cost of throughput. Requires full CPU and memory reservation
	// (requests = size). Recommended for latency-sensitive Telco VNF
	// workloads.
	VirtualMachineLatencySensitivityHigh VirtualMachineLatencySensitivityLevel = "High"

	// VirtualMachineLatencySensitivityHighWithHyperthreading combines High
	// latency sensitivity with simultaneous multi-threading (HT). Sets
	// ConfigSpec.LatencySensitivity.Level=High and
	// ConfigSpec.SimultaneousThreads=2. Requires full CPU and memory
	// reservation (requests = size).
	VirtualMachineLatencySensitivityHighWithHyperthreading VirtualMachineLatencySensitivityLevel = "HighWithHyperthreading"
)

// VirtualMachineCPUTopologySpec describes the guest CPU topology.
//
// Requires the TelcoVMServiceAPI supervisor capability.
// +kubebuilder:validation:XValidation:rule="!has(self.vnumaNodeCount) || has(self.coresPerSocket)",message="vnumaNodeCount requires coresPerSocket to also be set"
// +kubebuilder:validation:XValidation:rule="!has(self.numaFixedAutoAffinityEnabled) || !self.numaFixedAutoAffinityEnabled || !has(self.vnumaNodeCount)",message="vnumaNodeCount cannot be set when numaFixedAutoAffinityEnabled is true"
type VirtualMachineCPUTopologySpec struct {
	// +optional
	// +kubebuilder:validation:Minimum=1

	// CoresPerSocket controls the number of cores per virtual socket.
	// When unset, any explicit cores-per-socket configuration is cleared,
	// reverting to the vSphere default.
	// Maps to ConfigSpec.NumCoresPerSocket.
	// Requires power-off to apply.
	CoresPerSocket *int32 `json:"coresPerSocket,omitempty"`

	// +optional

	// NUMAFixedAutoAffinityEnabled enables fixed affinity between the VM's
	// vCPUs and physical NUMA nodes on the host. vCPUs are distributed equally
	// across (total vCPU count / CoresPerSocket) virtual NUMA nodes; when
	// CoresPerSocket is unset, all vCPUs are placed on a single physical NUMA
	// node. The affinity is established at each power-on and live migration.
	// Mutually exclusive with VNUMANodeCount.
	// When false or unset, fixed affinity is disabled.
	// Maps to ConfigSpec.NUMAFixedAutoAffinityEnabled.
	// Requires hardware version vmx-23 or later.
	// Requires power-off to apply.
	NUMAFixedAutoAffinityEnabled *bool `json:"numaFixedAutoAffinityEnabled,omitempty"`

	// +optional
	// +kubebuilder:validation:Minimum=1

	// VNUMANodeCount sets the number of virtual NUMA nodes. vSphere derives
	// coresPerNumaNode as total vCPU count / VNUMANodeCount; the vCPU count
	// should be evenly divisible by VNUMANodeCount for a balanced topology.
	// The derived coresPerNumaNode must be a multiple or divisor of
	// CoresPerSocket; the VM will fail to power on if this constraint is
	// violated.
	// When unset, any manual vNUMA configuration is cleared, enabling automatic
	// vNUMA sizing.
	// Must be set together with CoresPerSocket.
	// Maps to ConfigSpec.VirtualNuma.CoresPerNumaNode.
	// Requires hardware version vmx-20 or later.
	// Requires power-off to apply.
	VNUMANodeCount *int32 `json:"vnumaNodeCount,omitempty"`

	// +optional

	// ExposeVNUMAOnCPUHotAdd controls vNUMA exposure when CPU hot-add occurs.
	// When true, virtual NUMA topology is exposed to the guest when vCPUs are
	// hot-added.
	// When false or unset, a single virtual NUMA node is presented to the guest
	// during hot-add regardless of configured topology.
	// Only relevant when cpuAdvanced.hotAddEnabled is true.
	// Maps to ConfigSpec.VirtualNuma.ExposeVnumaOnCpuHotadd.
	// Requires hardware version vmx-20 or later.
	// Requires power-off to apply.
	ExposeVNUMAOnCPUHotAdd *bool `json:"exposeVnumaOnCpuHotadd,omitempty"`
}

// VirtualMachineCPUAdvancedSpec describes advanced CPU scheduling and topology
// configuration for a VirtualMachine.
//
// Requires the TelcoVMServiceAPI supervisor capability.
type VirtualMachineCPUAdvancedSpec struct {
	// +optional

	// LatencySensitivity sets the vSphere CPU scheduler latency sensitivity
	// level. Maps to ConfigSpec.LatencySensitivity.Level and
	// ConfigSpec.SimultaneousThreads (for HighWithHyperthreading).
	// Requires full CPU and memory reservation (requests = size)
	// when set to High or HighWithHyperthreading.
	// When unset, the setting is cleared, reverting to the vSphere default.
	LatencySensitivity *VirtualMachineLatencySensitivityLevel `json:"latencySensitivity,omitempty"`

	// +optional

	// Topology describes CPU topology configuration (cores per socket,
	// cores per NUMA node, vNUMA exposure on hot-add). All sub-fields are
	// optional and used together to define the guest CPU layout.
	Topology *VirtualMachineCPUTopologySpec `json:"topology,omitempty"`

	// +optional

	// HotAddEnabled allows vCPUs to be added to the VM while it is powered
	// on (CPU hot-add). Maps to ConfigSpec.CpuHotAddEnabled.
	// Requires hardware version vmx-20 or later.
	// Requires power-off to apply.
	HotAddEnabled *bool `json:"hotAddEnabled,omitempty"`

	// +optional

	// IOMMUEnabled enables Intel Virtualization Technology for Directed I/O
	// (VT-d / IOMMU) for this VM. Required for SR-IOV and PCI passthrough
	// workloads. Requires EFI firmware (spec.bootOptions.firmware = "efi").
	// Maps to ConfigSpec.Flags.VvtdEnabled.
	// Requires power-off to apply.
	IOMMUEnabled *bool `json:"iommuEnabled,omitempty"`

	// +optional

	// NestedHardwareVirtualizationEnabled exposes hardware-assisted
	// virtualization to the guest OS, enabling the guest to run its own
	// hypervisor or use hardware VMX instructions. Required for nested
	// virtualization workloads.
	// Maps to ConfigSpec.NestedHVEnabled.
	// Requires power-off to apply.
	NestedHardwareVirtualizationEnabled *bool `json:"nestedHardwareVirtualizationEnabled,omitempty"`

	// +optional

	// PerformanceCountersEnabled enables virtualized CPU performance
	// counters (vPMC), allowing profiling tools inside the guest OS to
	// access hardware performance counter data.
	// Maps to ConfigSpec.VPMCEnabled.
	// Requires power-off to apply.
	PerformanceCountersEnabled *bool `json:"performanceCountersEnabled,omitempty"`

	// +optional

	// ReservationLockedToMax pins the host CPU reservation to guarantee full
	// CPU capacity on the placed host. Overrides any explicit value in
	// spec.resources.requests.cpu when set.
	// When false or unset, spec.resources.requests.cpu is used for the
	// CPU reservation.
	// Maps to ConfigSpec.CpuAllocation.ReservationLockedToMax.
	// Requires hardware version vmx-23 or later.
	// Requires power-off to apply.
	ReservationLockedToMax *bool `json:"reservationLockedToMax,omitempty"`
}

// VirtualMachineMemoryAdvancedSpec describes advanced memory configuration
// for a VirtualMachine.
//
// Requires the TelcoVMServiceAPI supervisor capability.
type VirtualMachineMemoryAdvancedSpec struct {
	// +optional

	// HotAddEnabled allows memory to be added to the VM while it is powered
	// on (memory hot-add). Maps to ConfigSpec.MemoryHotAddEnabled.
	// Requires hardware version vmx-20 or later.
	// Requires power-off to apply.
	HotAddEnabled *bool `json:"hotAddEnabled,omitempty"`

	// +optional

	// ReservationLockedToMax pins the host memory reservation to the full
	// guest-visible memory size. Required for SR-IOV workloads, which need
	// full guest RAM pinned on the host.
	// Maps to ConfigSpec.MemoryReservationLockedToMax.
	// Requires power-off to apply.
	ReservationLockedToMax *bool `json:"reservationLockedToMax,omitempty"`
}
