/*
Copyright 2022 The Kubernetes Authors.

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

package agentpools

import (
	"context"

	asocontainerservicev1 "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001"
	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001/storage"
	asocontainerservicev1preview "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231102preview"
	"github.com/Azure/azure-service-operator/v2/pkg/genruntime"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/services/aso"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
	"sigs.k8s.io/cluster-api-provider-azure/util/versions"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// KubeletConfig defines the set of kubelet configurations for nodes in pools.
type KubeletConfig struct {
	// CPUManagerPolicy - CPU Manager policy to use.
	CPUManagerPolicy *string
	// CPUCfsQuota - Enable CPU CFS quota enforcement for containers that specify CPU limits.
	CPUCfsQuota *bool
	// CPUCfsQuotaPeriod - Sets CPU CFS quota period value.
	CPUCfsQuotaPeriod *string
	// ImageGcHighThreshold - The percent of disk usage after which image garbage collection is always run.
	ImageGcHighThreshold *int
	// ImageGcLowThreshold - The percent of disk usage before which image garbage collection is never run.
	ImageGcLowThreshold *int
	// TopologyManagerPolicy - Topology Manager policy to use.
	TopologyManagerPolicy *string
	// AllowedUnsafeSysctls - Allowlist of unsafe sysctls or unsafe sysctl patterns (ending in `*`).
	AllowedUnsafeSysctls []string
	// FailSwapOn - If set to true it will make the Kubelet fail to start if swap is enabled on the node.
	FailSwapOn *bool
	// ContainerLogMaxSizeMB - The maximum size (e.g. 10Mi) of container log file before it is rotated.
	ContainerLogMaxSizeMB *int
	// ContainerLogMaxFiles - The maximum number of container log files that can be present for a container. The number must be ≥ 2.
	ContainerLogMaxFiles *int
	// PodMaxPids - The maximum number of processes per pod.
	PodMaxPids *int
}

// AgentPoolSpec contains agent pool specification details.
type AgentPoolSpec struct {
	// Name is the name of the ASO ManagedClustersAgentPool resource.
	Name string

	// AzureName is the name of the agentpool resource in Azure.
	AzureName string

	// ResourceGroup is the name of the Azure resource group for the AKS Cluster.
	ResourceGroup string

	// Cluster is the name of the AKS cluster.
	Cluster string

	// Version defines the desired Kubernetes version.
	Version *string

	// SKU defines the Azure VM size for the agent pool VMs.
	SKU string

	// Replicas is the number of desired machines.
	Replicas int

	// OSDiskSizeGB is the OS disk size in GB for every machine in this agent pool.
	OSDiskSizeGB int

	// VnetSubnetID is the Azure Resource ID for the subnet which should contain nodes.
	VnetSubnetID string

	// Mode represents mode of an agent pool. Possible values include: 'System', 'User'.
	Mode string

	//  Maximum number of nodes for auto-scaling
	MaxCount *int `json:"maxCount,omitempty"`

	// Minimum number of nodes for auto-scaling
	MinCount *int `json:"minCount,omitempty"`

	// Node labels - labels for all of the nodes present in node pool
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`

	// NodeTaints specifies the taints for nodes present in this agent pool.
	NodeTaints []string `json:"nodeTaints,omitempty"`

	// EnableAutoScaling - Whether to enable auto-scaler
	EnableAutoScaling bool `json:"enableAutoScaling,omitempty"`

	// AvailabilityZones represents the Availability zones for nodes in the AgentPool.
	AvailabilityZones []string

	// MaxPods specifies the kubelet --max-pods configuration for the agent pool.
	MaxPods *int `json:"maxPods,omitempty"`

	// OsDiskType specifies the OS disk type for each node in the pool. Allowed values are 'Ephemeral' and 'Managed'.
	OsDiskType *string `json:"osDiskType,omitempty"`

	// EnableUltraSSD enables the storage type UltraSSD_LRS for the agent pool.
	EnableUltraSSD *bool `json:"enableUltraSSD,omitempty"`

	// OSType specifies the operating system for the node pool. Allowed values are 'Linux' and 'Windows'
	OSType *string `json:"osType,omitempty"`

	// EnableNodePublicIP controls whether or not nodes in the agent pool each have a public IP address.
	EnableNodePublicIP *bool `json:"enableNodePublicIP,omitempty"`

	// NodePublicIPPrefixID specifies the public IP prefix resource ID which VM nodes should use IPs from.
	NodePublicIPPrefixID string `json:"nodePublicIPPrefixID,omitempty"`

	// ScaleSetPriority specifies the ScaleSetPriority for the node pool. Allowed values are 'Spot' and 'Regular'
	ScaleSetPriority *string `json:"scaleSetPriority,omitempty"`

	// ScaleDownMode affects the cluster autoscaler behavior. Allowed values are 'Deallocate' and 'Delete'
	ScaleDownMode *string `json:"scaleDownMode,omitempty"`

	// SpotMaxPrice defines max price to pay for spot instance. Allowed values are any decimal value greater than zero or -1 which indicates the willingness to pay any on-demand price.
	SpotMaxPrice *resource.Quantity `json:"spotMaxPrice,omitempty"`

	// KubeletConfig specifies the kubelet configurations for nodes.
	KubeletConfig *KubeletConfig `json:"kubeletConfig,omitempty"`

	// KubeletDiskType specifies the kubelet disk type for each node in the pool. Allowed values are 'OS' and 'Temporary'
	KubeletDiskType *infrav1.KubeletDiskType `json:"kubeletDiskType,omitempty"`

	// AdditionalTags is an optional set of tags to add to Azure resources managed by the Azure provider, in addition to the ones added by default.
	AdditionalTags infrav1.Tags

	// LinuxOSConfig specifies the custom Linux OS settings and configurations
	LinuxOSConfig *infrav1.LinuxOSConfig

	// EnableFIPS indicates whether FIPS is enabled on the node pool
	EnableFIPS *bool

	// EnableEncryptionAtHost indicates whether host encryption is enabled on the node pool
	EnableEncryptionAtHost *bool

	// Patches are extra patches to be applied to the ASO resource.
	Patches []string

	// Preview indicates whether the agent pool is using a preview version of ASO.
	Preview bool
}

// ResourceRef implements azure.ASOResourceSpecGetter.
func (s *AgentPoolSpec) ResourceRef() genruntime.MetaObject {
	if s.Preview {
		return &asocontainerservicev1preview.ManagedClustersAgentPool{
			ObjectMeta: metav1.ObjectMeta{
				Name: azure.GetNormalizedKubernetesName(s.Name),
			},
		}
	}
	return &asocontainerservicev1.ManagedClustersAgentPool{
		ObjectMeta: metav1.ObjectMeta{
			Name: azure.GetNormalizedKubernetesName(s.Name),
		},
	}
}

// getManagedMachinePoolVersion gets the desired managed Kubernetes version.
// If the auto upgrade channel is set to patch, stable or rapid, clusters can be upgraded to a higher version by AKS.
// If auto upgrade is triggered, the existing Kubernetes version will be higher than the user's desired Kubernetes version.
// CAPZ should honour the upgrade and it should not downgrade to the lower desired version.
func (s *AgentPoolSpec) getManagedMachinePoolVersion(existing *asocontainerservicev1hub.ManagedClustersAgentPool) *string {
	if existing == nil || existing.Spec.OrchestratorVersion == nil {
		return s.Version
	}
	if s.Version == nil {
		return existing.Spec.OrchestratorVersion
	}
	v := versions.GetHigherK8sVersion(*s.Version, *existing.Spec.OrchestratorVersion)
	return ptr.To(v)
}

// Parameters returns the parameters for the agent pool.
func (s *AgentPoolSpec) Parameters(ctx context.Context, existingObj genruntime.MetaObject) (params genruntime.MetaObject, err error) {
	_, _, done := tele.StartSpanWithLogger(ctx, "agentpools.Service.Parameters")
	defer done()

	// If existing is preview, convert to stable then back to preview at the end of the function.
	var existing *asocontainerservicev1hub.ManagedClustersAgentPool
	if existingObj != nil {
		hub := &asocontainerservicev1hub.ManagedClustersAgentPool{}
		if err := existingObj.(conversion.Convertible).ConvertTo(hub); err != nil {
			return nil, err
		}
		existing = hub
	}

	agentPool := existing
	if agentPool == nil {
		agentPool = &asocontainerservicev1hub.ManagedClustersAgentPool{}
	}

	agentPool.Spec.AzureName = s.AzureName
	agentPool.Spec.Owner = &genruntime.KnownResourceReference{
		Name: s.Cluster,
	}
	agentPool.Spec.AvailabilityZones = s.AvailabilityZones
	agentPool.Spec.Count = &s.Replicas
	agentPool.Spec.EnableAutoScaling = ptr.To(s.EnableAutoScaling)
	agentPool.Spec.EnableUltraSSD = s.EnableUltraSSD
	agentPool.Spec.KubeletDiskType = azure.AliasOrNil[string]((*string)(s.KubeletDiskType))
	agentPool.Spec.MaxCount = s.MaxCount
	agentPool.Spec.MaxPods = s.MaxPods
	agentPool.Spec.MinCount = s.MinCount
	agentPool.Spec.Mode = ptr.To(string(asocontainerservicev1.AgentPoolMode(s.Mode)))
	agentPool.Spec.NodeLabels = s.NodeLabels
	agentPool.Spec.NodeTaints = s.NodeTaints
	agentPool.Spec.OsDiskSizeGB = ptr.To(int(asocontainerservicev1.ContainerServiceOSDisk(s.OSDiskSizeGB)))
	agentPool.Spec.OsDiskType = azure.AliasOrNil[string](s.OsDiskType)
	agentPool.Spec.OsType = azure.AliasOrNil[string](s.OSType)
	agentPool.Spec.ScaleSetPriority = azure.AliasOrNil[string](s.ScaleSetPriority)
	agentPool.Spec.ScaleDownMode = azure.AliasOrNil[string](s.ScaleDownMode)
	agentPool.Spec.Type = ptr.To(string(asocontainerservicev1.AgentPoolType_VirtualMachineScaleSets))
	agentPool.Spec.EnableNodePublicIP = s.EnableNodePublicIP
	agentPool.Spec.Tags = s.AdditionalTags
	agentPool.Spec.EnableFIPS = s.EnableFIPS
	agentPool.Spec.EnableEncryptionAtHost = s.EnableEncryptionAtHost
	if kubernetesVersion := s.getManagedMachinePoolVersion(existing); kubernetesVersion != nil {
		agentPool.Spec.OrchestratorVersion = kubernetesVersion
	}

	if s.KubeletConfig != nil {
		agentPool.Spec.KubeletConfig = &asocontainerservicev1hub.KubeletConfig{
			CpuManagerPolicy:      s.KubeletConfig.CPUManagerPolicy,
			CpuCfsQuota:           s.KubeletConfig.CPUCfsQuota,
			CpuCfsQuotaPeriod:     s.KubeletConfig.CPUCfsQuotaPeriod,
			ImageGcHighThreshold:  s.KubeletConfig.ImageGcHighThreshold,
			ImageGcLowThreshold:   s.KubeletConfig.ImageGcLowThreshold,
			TopologyManagerPolicy: s.KubeletConfig.TopologyManagerPolicy,
			FailSwapOn:            s.KubeletConfig.FailSwapOn,
			ContainerLogMaxSizeMB: s.KubeletConfig.ContainerLogMaxSizeMB,
			ContainerLogMaxFiles:  s.KubeletConfig.ContainerLogMaxFiles,
			PodMaxPids:            s.KubeletConfig.PodMaxPids,
			AllowedUnsafeSysctls:  s.KubeletConfig.AllowedUnsafeSysctls,
		}
	}

	if s.SKU != "" {
		agentPool.Spec.VmSize = &s.SKU
	}

	if s.SpotMaxPrice != nil {
		agentPool.Spec.SpotMaxPrice = ptr.To(s.SpotMaxPrice.AsApproximateFloat64())
	}

	if s.VnetSubnetID != "" {
		agentPool.Spec.VnetSubnetReference = &genruntime.ResourceReference{
			ARMID: s.VnetSubnetID,
		}
	}

	if s.NodePublicIPPrefixID != "" {
		agentPool.Spec.NodePublicIPPrefixReference = &genruntime.ResourceReference{
			ARMID: s.NodePublicIPPrefixID,
		}
	}

	if s.LinuxOSConfig != nil {
		agentPool.Spec.LinuxOSConfig = &asocontainerservicev1hub.LinuxOSConfig{
			SwapFileSizeMB:             s.LinuxOSConfig.SwapFileSizeMB,
			TransparentHugePageEnabled: (*string)(s.LinuxOSConfig.TransparentHugePageEnabled),
			TransparentHugePageDefrag:  (*string)(s.LinuxOSConfig.TransparentHugePageDefrag),
		}
		if s.LinuxOSConfig.Sysctls != nil {
			agentPool.Spec.LinuxOSConfig.Sysctls = &asocontainerservicev1hub.SysctlConfig{
				FsAioMaxNr:                     s.LinuxOSConfig.Sysctls.FsAioMaxNr,
				FsFileMax:                      s.LinuxOSConfig.Sysctls.FsFileMax,
				FsInotifyMaxUserWatches:        s.LinuxOSConfig.Sysctls.FsInotifyMaxUserWatches,
				FsNrOpen:                       s.LinuxOSConfig.Sysctls.FsNrOpen,
				KernelThreadsMax:               s.LinuxOSConfig.Sysctls.KernelThreadsMax,
				NetCoreNetdevMaxBacklog:        s.LinuxOSConfig.Sysctls.NetCoreNetdevMaxBacklog,
				NetCoreOptmemMax:               s.LinuxOSConfig.Sysctls.NetCoreOptmemMax,
				NetCoreRmemDefault:             s.LinuxOSConfig.Sysctls.NetCoreRmemDefault,
				NetCoreRmemMax:                 s.LinuxOSConfig.Sysctls.NetCoreRmemMax,
				NetCoreSomaxconn:               s.LinuxOSConfig.Sysctls.NetCoreSomaxconn,
				NetCoreWmemDefault:             s.LinuxOSConfig.Sysctls.NetCoreWmemDefault,
				NetCoreWmemMax:                 s.LinuxOSConfig.Sysctls.NetCoreWmemMax,
				NetIpv4IpLocalPortRange:        s.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
				NetIpv4NeighDefaultGcThresh1:   s.LinuxOSConfig.Sysctls.NetIpv4NeighDefaultGcThresh1,
				NetIpv4NeighDefaultGcThresh2:   s.LinuxOSConfig.Sysctls.NetIpv4NeighDefaultGcThresh2,
				NetIpv4NeighDefaultGcThresh3:   s.LinuxOSConfig.Sysctls.NetIpv4NeighDefaultGcThresh3,
				NetIpv4TcpFinTimeout:           s.LinuxOSConfig.Sysctls.NetIpv4TCPFinTimeout,
				NetIpv4TcpKeepaliveProbes:      s.LinuxOSConfig.Sysctls.NetIpv4TCPKeepaliveProbes,
				NetIpv4TcpKeepaliveTime:        s.LinuxOSConfig.Sysctls.NetIpv4TCPKeepaliveTime,
				NetIpv4TcpMaxSynBacklog:        s.LinuxOSConfig.Sysctls.NetIpv4TCPMaxSynBacklog,
				NetIpv4TcpMaxTwBuckets:         s.LinuxOSConfig.Sysctls.NetIpv4TCPMaxTwBuckets,
				NetIpv4TcpTwReuse:              s.LinuxOSConfig.Sysctls.NetIpv4TCPTwReuse,
				NetIpv4TcpkeepaliveIntvl:       s.LinuxOSConfig.Sysctls.NetIpv4TCPkeepaliveIntvl,
				NetNetfilterNfConntrackBuckets: s.LinuxOSConfig.Sysctls.NetNetfilterNfConntrackBuckets,
				NetNetfilterNfConntrackMax:     s.LinuxOSConfig.Sysctls.NetNetfilterNfConntrackMax,
				VmMaxMapCount:                  s.LinuxOSConfig.Sysctls.VMMaxMapCount,
				VmSwappiness:                   s.LinuxOSConfig.Sysctls.VMSwappiness,
				VmVfsCachePressure:             s.LinuxOSConfig.Sysctls.VMVfsCachePressure,
			}
		}
	}

	// When autoscaling is set, the count of the nodes differ based on the autoscaler and should not depend on the
	// count present in MachinePool or AzureManagedMachinePool, hence we should not make an update API call based
	// on difference in count.
	if s.EnableAutoScaling && agentPool.Status.Count != nil {
		agentPool.Spec.Count = agentPool.Status.Count
	}

	if s.Preview {
		prev := &asocontainerservicev1preview.ManagedClustersAgentPool{}
		if err := prev.ConvertFrom(agentPool); err != nil {
			return nil, err
		}
		return prev, nil
	}

	stable := &asocontainerservicev1.ManagedClustersAgentPool{}
	if err := stable.ConvertFrom(agentPool); err != nil {
		return nil, err
	}

	return stable, nil
}

// WasManaged implements azure.ASOResourceSpecGetter.
func (s *AgentPoolSpec) WasManaged(resource genruntime.MetaObject) bool {
	// CAPZ has never supported BYO agent pools.
	return true
}

var _ aso.Patcher = (*AgentPoolSpec)(nil)

// ExtraPatches implements aso.Patcher.
func (s *AgentPoolSpec) ExtraPatches() []string {
	return s.Patches
}
