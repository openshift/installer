/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// LabelAgentPoolMode represents mode of an agent pool. Possible values include: System, User.
	LabelAgentPoolMode = "azuremanagedmachinepool.infrastructure.cluster.x-k8s.io/agentpoolmode"

	// NodePoolModeSystem represents mode system for azuremachinepool.
	NodePoolModeSystem NodePoolMode = "System"

	// NodePoolModeUser represents mode user for azuremachinepool.
	NodePoolModeUser NodePoolMode = "User"

	// DefaultOSType represents the default operating system for azmachinepool.
	DefaultOSType string = LinuxOS
)

// NodePoolMode enumerates the values for agent pool mode.
type NodePoolMode string

// CPUManagerPolicy enumerates the values for KubeletConfig.CPUManagerPolicy.
type CPUManagerPolicy string

const (
	// CPUManagerPolicyNone ...
	CPUManagerPolicyNone CPUManagerPolicy = "none"
	// CPUManagerPolicyStatic ...
	CPUManagerPolicyStatic CPUManagerPolicy = "static"
)

// TopologyManagerPolicy enumerates the values for KubeletConfig.TopologyManagerPolicy.
type TopologyManagerPolicy string

// KubeletDiskType enumerates the values for the agent pool's KubeletDiskType.
type KubeletDiskType string

const (
	// KubeletDiskTypeOS ...
	KubeletDiskTypeOS KubeletDiskType = "OS"
	// KubeletDiskTypeTemporary ...
	KubeletDiskTypeTemporary KubeletDiskType = "Temporary"
)

const (
	// TopologyManagerPolicyNone ...
	TopologyManagerPolicyNone TopologyManagerPolicy = "none"
	// TopologyManagerPolicyBestEffort ...
	TopologyManagerPolicyBestEffort TopologyManagerPolicy = "best-effort"
	// TopologyManagerPolicyRestricted ...
	TopologyManagerPolicyRestricted TopologyManagerPolicy = "restricted"
	// TopologyManagerPolicySingleNumaNode ...
	TopologyManagerPolicySingleNumaNode TopologyManagerPolicy = "single-numa-node"
)

// TransparentHugePageOption enumerates the values for various modes of Transparent Hugepages.
type TransparentHugePageOption string

const (
	// TransparentHugePageOptionAlways ...
	TransparentHugePageOptionAlways TransparentHugePageOption = "always"

	// TransparentHugePageOptionMadvise ...
	TransparentHugePageOptionMadvise TransparentHugePageOption = "madvise"

	// TransparentHugePageOptionNever ...
	TransparentHugePageOptionNever TransparentHugePageOption = "never"

	// TransparentHugePageOptionDefer ...
	TransparentHugePageOptionDefer TransparentHugePageOption = "defer"

	// TransparentHugePageOptionDeferMadvise ...
	TransparentHugePageOptionDeferMadvise TransparentHugePageOption = "defer+madvise"
)

// KubeletConfig defines the supported subset of kubelet configurations for nodes in pools.
// See also [AKS doc], [K8s doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/custom-node-configuration
// [K8s doc]: https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/
type KubeletConfig struct {
	// CPUManagerPolicy - CPU Manager policy to use.
	// +kubebuilder:validation:Enum=none;static
	// +optional
	CPUManagerPolicy *CPUManagerPolicy `json:"cpuManagerPolicy,omitempty"`
	// CPUCfsQuota - Enable CPU CFS quota enforcement for containers that specify CPU limits.
	// +optional
	CPUCfsQuota *bool `json:"cpuCfsQuota,omitempty"`
	// CPUCfsQuotaPeriod - Sets CPU CFS quota period value.
	// Must end in "ms", e.g. "100ms"
	// +optional
	CPUCfsQuotaPeriod *string `json:"cpuCfsQuotaPeriod,omitempty"`
	// ImageGcHighThreshold - The percent of disk usage after which image garbage collection is always run.
	// Valid values are 0-100 (inclusive).
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	ImageGcHighThreshold *int `json:"imageGcHighThreshold,omitempty"`
	// ImageGcLowThreshold - The percent of disk usage before which image garbage collection is never run.
	// Valid values are 0-100 (inclusive) and must be less than `imageGcHighThreshold`.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	ImageGcLowThreshold *int `json:"imageGcLowThreshold,omitempty"`
	// TopologyManagerPolicy - Topology Manager policy to use.
	// +kubebuilder:validation:Enum=none;best-effort;restricted;single-numa-node
	// +optional
	TopologyManagerPolicy *TopologyManagerPolicy `json:"topologyManagerPolicy,omitempty"`
	// AllowedUnsafeSysctls - Allowlist of unsafe sysctls or unsafe sysctl patterns (ending in `*`).
	// Valid values match `kernel.shm*`, `kernel.msg*`, `kernel.sem`, `fs.mqueue.*`, or `net.*`.
	// +optional
	AllowedUnsafeSysctls []string `json:"allowedUnsafeSysctls,omitempty"`
	// FailSwapOn - If set to true it will make the Kubelet fail to start if swap is enabled on the node.
	// +optional
	FailSwapOn *bool `json:"failSwapOn,omitempty"`
	// ContainerLogMaxSizeMB - The maximum size in MB of a container log file before it is rotated.
	// +optional
	ContainerLogMaxSizeMB *int `json:"containerLogMaxSizeMB,omitempty"`
	// ContainerLogMaxFiles - The maximum number of container log files that can be present for a container. The number must be ≥ 2.
	// +kubebuilder:validation:Minimum=2
	// +optional
	ContainerLogMaxFiles *int `json:"containerLogMaxFiles,omitempty"`
	// PodMaxPids - The maximum number of processes per pod.
	// Must not exceed kernel PID limit. -1 disables the limit.
	// +kubebuilder:validation:Minimum=-1
	// +optional
	PodMaxPids *int `json:"podMaxPids,omitempty"`
}

// SysctlConfig specifies the settings for Linux agent nodes.
type SysctlConfig struct {
	// FsAioMaxNr specifies the maximum number of system-wide asynchronous io requests.
	// Valid values are 65536-6553500 (inclusive).
	// Maps to fs.aio-max-nr.
	// +kubebuilder:validation:Minimum=65536
	// +kubebuilder:validation:Maximum=6553500
	// +optional
	FsAioMaxNr *int `json:"fsAioMaxNr,omitempty"`

	// FsFileMax specifies the max number of file-handles that the Linux kernel will allocate, by increasing increases the maximum number of open files permitted.
	// Valid values are 8192-12000500 (inclusive).
	// Maps to fs.file-max.
	// +kubebuilder:validation:Minimum=8192
	// +kubebuilder:validation:Maximum=12000500
	// +optional
	FsFileMax *int `json:"fsFileMax,omitempty"`

	// FsInotifyMaxUserWatches specifies the number of file watches allowed by the system. Each watch is roughly 90 bytes on a 32-bit kernel, and roughly 160 bytes on a 64-bit kernel.
	// Valid values are 781250-2097152 (inclusive).
	// Maps to fs.inotify.max_user_watches.
	// +kubebuilder:validation:Minimum=781250
	// +kubebuilder:validation:Maximum=2097152
	// +optional
	FsInotifyMaxUserWatches *int `json:"fsInotifyMaxUserWatches,omitempty"`

	// FsNrOpen specifies the maximum number of file-handles a process can allocate.
	// Valid values are 8192-20000500 (inclusive).
	// Maps to fs.nr_open.
	// +kubebuilder:validation:Minimum=8192
	// +kubebuilder:validation:Maximum=20000500
	// +optional
	FsNrOpen *int `json:"fsNrOpen,omitempty"`

	// KernelThreadsMax specifies the maximum number of all threads that can be created.
	// Valid values are 20-513785 (inclusive).
	// Maps to kernel.threads-max.
	// +kubebuilder:validation:Minimum=20
	// +kubebuilder:validation:Maximum=513785
	// +optional
	KernelThreadsMax *int `json:"kernelThreadsMax,omitempty"`

	// NetCoreNetdevMaxBacklog specifies maximum number of packets, queued on the INPUT side, when the interface receives packets faster than kernel can process them.
	// Valid values are 1000-3240000 (inclusive).
	// Maps to net.core.netdev_max_backlog.
	// +kubebuilder:validation:Minimum=1000
	// +kubebuilder:validation:Maximum=3240000
	// +optional
	NetCoreNetdevMaxBacklog *int `json:"netCoreNetdevMaxBacklog,omitempty"`

	// NetCoreOptmemMax specifies the maximum ancillary buffer size (option memory buffer) allowed per socket.
	// Socket option memory is used in a few cases to store extra structures relating to usage of the socket.
	// Valid values are 20480-4194304 (inclusive).
	// Maps to net.core.optmem_max.
	// +kubebuilder:validation:Minimum=20480
	// +kubebuilder:validation:Maximum=4194304
	// +optional
	NetCoreOptmemMax *int `json:"netCoreOptmemMax,omitempty"`

	// NetCoreRmemDefault specifies the default receive socket buffer size in bytes.
	// Valid values are 212992-134217728 (inclusive).
	// Maps to net.core.rmem_default.
	// +kubebuilder:validation:Minimum=212992
	// +kubebuilder:validation:Maximum=134217728
	// +optional
	NetCoreRmemDefault *int `json:"netCoreRmemDefault,omitempty"`

	// NetCoreRmemMax specifies the maximum receive socket buffer size in bytes.
	// Valid values are 212992-134217728 (inclusive).
	// Maps to net.core.rmem_max.
	// +kubebuilder:validation:Minimum=212992
	// +kubebuilder:validation:Maximum=134217728
	// +optional
	NetCoreRmemMax *int `json:"netCoreRmemMax,omitempty"`

	// NetCoreSomaxconn specifies maximum number of connection requests that can be queued for any given listening socket.
	// An upper limit for the value of the backlog parameter passed to the listen(2)(https://man7.org/linux/man-pages/man2/listen.2.html) function.
	// If the backlog argument is greater than the somaxconn, then it's silently truncated to this limit.
	// Valid values are 4096-3240000 (inclusive).
	// Maps to net.core.somaxconn.
	// +kubebuilder:validation:Minimum=4096
	// +kubebuilder:validation:Maximum=3240000
	// +optional
	NetCoreSomaxconn *int `json:"netCoreSomaxconn,omitempty"`

	// NetCoreWmemDefault specifies the default send socket buffer size in bytes.
	// Valid values are 212992-134217728 (inclusive).
	// Maps to net.core.wmem_default.
	// +kubebuilder:validation:Minimum=212992
	// +kubebuilder:validation:Maximum=134217728
	// +optional
	NetCoreWmemDefault *int `json:"netCoreWmemDefault,omitempty"`

	// NetCoreWmemMax specifies the maximum send socket buffer size in bytes.
	// Valid values are 212992-134217728 (inclusive).
	// Maps to net.core.wmem_max.
	// +kubebuilder:validation:Minimum=212992
	// +kubebuilder:validation:Maximum=134217728
	// +optional
	NetCoreWmemMax *int `json:"netCoreWmemMax,omitempty"`

	// NetIpv4IPLocalPortRange is used by TCP and UDP traffic to choose the local port on the agent node.
	// PortRange should be specified in the format "first last".
	// First, being an integer, must be between [1024 - 60999].
	// Last, being an integer, must be between [32768 - 65000].
	// Maps to net.ipv4.ip_local_port_range.
	// +optional
	NetIpv4IPLocalPortRange *string `json:"netIpv4IPLocalPortRange,omitempty"`

	// NetIpv4NeighDefaultGcThresh1 specifies the minimum number of entries that may be in the ARP cache.
	// Garbage collection won't be triggered if the number of entries is below this setting.
	// Valid values are 128-80000 (inclusive).
	// Maps to net.ipv4.neigh.default.gc_thresh1.
	// +kubebuilder:validation:Minimum=128
	// +kubebuilder:validation:Maximum=80000
	// +optional
	NetIpv4NeighDefaultGcThresh1 *int `json:"netIpv4NeighDefaultGcThresh1,omitempty"`

	// NetIpv4NeighDefaultGcThresh2 specifies soft maximum number of entries that may be in the ARP cache.
	// ARP garbage collection will be triggered about 5 seconds after reaching this soft maximum.
	// Valid values are 512-90000 (inclusive).
	// Maps to net.ipv4.neigh.default.gc_thresh2.
	// +kubebuilder:validation:Minimum=512
	// +kubebuilder:validation:Maximum=90000
	// +optional
	NetIpv4NeighDefaultGcThresh2 *int `json:"netIpv4NeighDefaultGcThresh2,omitempty"`

	// NetIpv4NeighDefaultGcThresh3 specified hard maximum number of entries in the ARP cache.
	// Valid values are 1024-100000 (inclusive).
	// Maps to net.ipv4.neigh.default.gc_thresh3.
	// +kubebuilder:validation:Minimum=1024
	// +kubebuilder:validation:Maximum=100000
	// +optional
	NetIpv4NeighDefaultGcThresh3 *int `json:"netIpv4NeighDefaultGcThresh3,omitempty"`

	// NetIpv4TCPFinTimeout specifies the length of time an orphaned connection will remain in the FIN_WAIT_2 state before it's aborted at the local end.
	// Valid values are 5-120 (inclusive).
	// Maps to net.ipv4.tcp_fin_timeout.
	// +kubebuilder:validation:Minimum=5
	// +kubebuilder:validation:Maximum=120
	// +optional
	NetIpv4TCPFinTimeout *int `json:"netIpv4TCPFinTimeout,omitempty"`

	// NetIpv4TCPKeepaliveProbes specifies the number of keepalive probes TCP sends out, until it decides the connection is broken.
	// Valid values are 1-15 (inclusive).
	// Maps to net.ipv4.tcp_keepalive_probes.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=15
	// +optional
	NetIpv4TCPKeepaliveProbes *int `json:"netIpv4TCPKeepaliveProbes,omitempty"`

	// NetIpv4TCPKeepaliveTime specifies the rate at which TCP sends out a keepalive message when keepalive is enabled.
	// Valid values are 30-432000 (inclusive).
	// Maps to net.ipv4.tcp_keepalive_time.
	// +kubebuilder:validation:Minimum=30
	// +kubebuilder:validation:Maximum=432000
	// +optional
	NetIpv4TCPKeepaliveTime *int `json:"netIpv4TCPKeepaliveTime,omitempty"`

	// NetIpv4TCPMaxSynBacklog specifies the maximum number of queued connection requests that have still not received an acknowledgment from the connecting client.
	// If this number is exceeded, the kernel will begin dropping requests.
	// Valid values are 128-3240000 (inclusive).
	// Maps to net.ipv4.tcp_max_syn_backlog.
	// +kubebuilder:validation:Minimum=128
	// +kubebuilder:validation:Maximum=3240000
	// +optional
	NetIpv4TCPMaxSynBacklog *int `json:"netIpv4TCPMaxSynBacklog,omitempty"`

	// NetIpv4TCPMaxTwBuckets specifies maximal number of timewait sockets held by system simultaneously.
	// If this number is exceeded, time-wait socket is immediately destroyed and warning is printed.
	// Valid values are 8000-1440000 (inclusive).
	// Maps to net.ipv4.tcp_max_tw_buckets.
	// +kubebuilder:validation:Minimum=8000
	// +kubebuilder:validation:Maximum=1440000
	// +optional
	NetIpv4TCPMaxTwBuckets *int `json:"netIpv4TCPMaxTwBuckets,omitempty"`

	// NetIpv4TCPTwReuse is used to allow to reuse TIME-WAIT sockets for new connections when it's safe from protocol viewpoint.
	// Maps to net.ipv4.tcp_tw_reuse.
	// +optional
	NetIpv4TCPTwReuse *bool `json:"netIpv4TCPTwReuse,omitempty"`

	// NetIpv4TCPkeepaliveIntvl specifies the frequency of the probes sent out.
	// Multiplied by tcpKeepaliveprobes, it makes up the time to kill a connection that isn't responding, after probes started.
	// Valid values are 1-75 (inclusive).
	// Maps to net.ipv4.tcp_keepalive_intvl.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=75
	// +optional
	NetIpv4TCPkeepaliveIntvl *int `json:"netIpv4TCPkeepaliveIntvl,omitempty"`

	// NetNetfilterNfConntrackBuckets specifies the size of hash table used by nf_conntrack module to record the established connection record of the TCP protocol.
	// Valid values are 65536-147456 (inclusive).
	// Maps to net.netfilter.nf_conntrack_buckets.
	// +kubebuilder:validation:Minimum=65536
	// +kubebuilder:validation:Maximum=147456
	// +optional
	NetNetfilterNfConntrackBuckets *int `json:"netNetfilterNfConntrackBuckets,omitempty"`

	// NetNetfilterNfConntrackMax specifies the maximum number of connections supported by the nf_conntrack module or the size of connection tracking table.
	// Valid values are 131072-1048576 (inclusive).
	// Maps to net.netfilter.nf_conntrack_max.
	// +kubebuilder:validation:Minimum=131072
	// +kubebuilder:validation:Maximum=1048576
	// +optional
	NetNetfilterNfConntrackMax *int `json:"netNetfilterNfConntrackMax,omitempty"`

	// VMMaxMapCount specifies the maximum number of memory map areas a process may have.
	// Maps to vm.max_map_count.
	// Valid values are 65530-262144 (inclusive).
	// +kubebuilder:validation:Minimum=65530
	// +kubebuilder:validation:Maximum=262144
	// +optional
	VMMaxMapCount *int `json:"vmMaxMapCount,omitempty"`

	// VMSwappiness specifies aggressiveness of the kernel in swapping memory pages.
	// Higher values will increase aggressiveness, lower values decrease the amount of swap.
	// Valid values are 0-100 (inclusive).
	// Maps to vm.swappiness.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	// +optional
	VMSwappiness *int `json:"vmSwappiness,omitempty"`

	// VMVfsCachePressure specifies the percentage value that controls tendency of the kernel to reclaim the memory, which is used for caching of directory and inode objects.
	// Valid values are 1-500 (inclusive).
	// Maps to vm.vfs_cache_pressure.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=500
	// +optional
	VMVfsCachePressure *int `json:"vmVfsCachePressure,omitempty"`
}

// LinuxOSConfig specifies the custom Linux OS settings and configurations.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/custom-node-configuration#linux-os-custom-configuration
type LinuxOSConfig struct {
	// SwapFileSizeMB specifies size in MB of a swap file will be created on the agent nodes from this node pool.
	// Max value of SwapFileSizeMB should be the size of temporary disk(/dev/sdb).
	// Must be at least 1.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/azure/virtual-machines/managed-disks-overview#temporary-disk
	// +kubebuilder:validation:Minimum=1
	// +optional
	SwapFileSizeMB *int `json:"swapFileSizeMB,omitempty"`

	// Sysctl specifies the settings for Linux agent nodes.
	// +optional
	Sysctls *SysctlConfig `json:"sysctls,omitempty"`

	// TransparentHugePageDefrag specifies whether the kernel should make aggressive use of memory compaction to make more hugepages available.
	// See also [Linux doc].
	//
	// [Linux doc]: https://www.kernel.org/doc/html/latest/admin-guide/mm/transhuge.html#admin-guide-transhuge for more details.
	// +kubebuilder:validation:Enum=always;defer;defer+madvise;madvise;never
	// +optional
	TransparentHugePageDefrag *TransparentHugePageOption `json:"transparentHugePageDefrag,omitempty"`

	// TransparentHugePageEnabled specifies various modes of Transparent Hugepages.
	// See also [Linux doc].
	//
	// [Linux doc]: https://www.kernel.org/doc/html/latest/admin-guide/mm/transhuge.html#admin-guide-transhuge for more details.
	// +kubebuilder:validation:Enum=always;madvise;never
	// +optional
	TransparentHugePageEnabled *TransparentHugePageOption `json:"transparentHugePageEnabled,omitempty"`
}

// AzureManagedMachinePoolSpec defines the desired state of AzureManagedMachinePool.
type AzureManagedMachinePoolSpec struct {
	AzureManagedMachinePoolClassSpec `json:",inline"`

	// ProviderIDList is the unique identifier as specified by the cloud provider.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`
}

// ManagedMachinePoolScaling specifies scaling options.
type ManagedMachinePoolScaling struct {
	// MinSize is the minimum number of nodes for auto-scaling.
	MinSize *int `json:"minSize,omitempty"`
	// MaxSize is the maximum number of nodes for auto-scaling.
	MaxSize *int `json:"maxSize,omitempty"`
}

// TaintEffect is the effect for a Kubernetes taint.
type TaintEffect string

// Taint represents a Kubernetes taint.
type Taint struct {
	// Effect specifies the effect for the taint
	// +kubebuilder:validation:Enum=NoSchedule;NoExecute;PreferNoSchedule
	Effect TaintEffect `json:"effect"`
	// Key is the key of the taint
	Key string `json:"key"`
	// Value is the value of the taint
	Value string `json:"value"`
}

// Taints is an array of Taints.
type Taints []Taint

// AzureManagedMachinePoolStatus defines the observed state of AzureManagedMachinePool.
type AzureManagedMachinePoolStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready"`

	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorReason *string `json:"errorReason,omitempty"`

	// Any transient errors that occur during the reconciliation of Machines
	// can be added as events to the Machine object and/or logged in the
	// controller's output.
	// +optional
	ErrorMessage *string `json:"errorMessage,omitempty"`

	// Conditions defines current service state of the AzureManagedControlPlane.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// LongRunningOperationStates saves the states for Azure long-running operations so they can be continued on the
	// next reconciliation loop.
	// +optional
	LongRunningOperationStates Futures `json:"longRunningOperationStates,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AzureManagedMachinePool belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",priority=1,JSONPath=".status.conditions[?(@.type=='Ready')].message"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of this AzureManagedMachinePool"
// +kubebuilder:printcolumn:name="Mode",type="string",JSONPath=".spec.mode"
// +kubebuilder:resource:path=azuremanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=ammp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// AzureManagedMachinePool is the Schema for the azuremanagedmachinepools API.
type AzureManagedMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureManagedMachinePoolSpec   `json:"spec,omitempty"`
	Status AzureManagedMachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureManagedMachinePoolList contains a list of AzureManagedMachinePools.
type AzureManagedMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureManagedMachinePool `json:"items"`
}

// GetConditions returns the list of conditions for an AzureManagedMachinePool API object.
func (m *AzureManagedMachinePool) GetConditions() clusterv1.Conditions {
	return m.Status.Conditions
}

// SetConditions will set the given conditions on an AzureManagedMachinePool object.
func (m *AzureManagedMachinePool) SetConditions(conditions clusterv1.Conditions) {
	m.Status.Conditions = conditions
}

// GetFutures returns the list of long running operation states for an AzureManagedMachinePool API object.
func (m *AzureManagedMachinePool) GetFutures() Futures {
	return m.Status.LongRunningOperationStates
}

// SetFutures will set the given long running operation states on an AzureManagedMachinePool object.
func (m *AzureManagedMachinePool) SetFutures(futures Futures) {
	m.Status.LongRunningOperationStates = futures
}

func init() {
	SchemeBuilder.Register(&AzureManagedMachinePool{}, &AzureManagedMachinePoolList{})
}
