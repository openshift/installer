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
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/azure"
	"sigs.k8s.io/cluster-api-provider-azure/azure/converters"
	azureutil "sigs.k8s.io/cluster-api-provider-azure/util/azure"
	"sigs.k8s.io/cluster-api-provider-azure/util/tele"
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
	ImageGcHighThreshold *int32
	// ImageGcLowThreshold - The percent of disk usage before which image garbage collection is never run.
	ImageGcLowThreshold *int32
	// TopologyManagerPolicy - Topology Manager policy to use.
	TopologyManagerPolicy *string
	// AllowedUnsafeSysctls - Allowlist of unsafe sysctls or unsafe sysctl patterns (ending in `*`).
	AllowedUnsafeSysctls *[]string
	// FailSwapOn - If set to true it will make the Kubelet fail to start if swap is enabled on the node.
	FailSwapOn *bool
	// ContainerLogMaxSizeMB - The maximum size (e.g. 10Mi) of container log file before it is rotated.
	ContainerLogMaxSizeMB *int32
	// ContainerLogMaxFiles - The maximum number of container log files that can be present for a container. The number must be ≥ 2.
	ContainerLogMaxFiles *int32
	// PodMaxPids - The maximum number of processes per pod.
	PodMaxPids *int32
}

// AgentPoolSpec contains agent pool specification details.
type AgentPoolSpec struct {
	// Name is the name of agent pool.
	Name string

	// ResourceGroup is the name of the Azure resource group for the AKS Cluster.
	ResourceGroup string

	// Cluster is the name of the AKS cluster.
	Cluster string

	// Version defines the desired Kubernetes version.
	Version *string

	// SKU defines the Azure VM size for the agent pool VMs.
	SKU string

	// Replicas is the number of desired machines.
	Replicas int32

	// OSDiskSizeGB is the OS disk size in GB for every machine in this agent pool.
	OSDiskSizeGB int32

	// VnetSubnetID is the Azure Resource ID for the subnet which should contain nodes.
	VnetSubnetID string

	// Mode represents mode of an agent pool. Possible values include: 'System', 'User'.
	Mode string

	//  Maximum number of nodes for auto-scaling
	MaxCount *int32 `json:"maxCount,omitempty"`

	// Minimum number of nodes for auto-scaling
	MinCount *int32 `json:"minCount,omitempty"`

	// Node labels - labels for all of the nodes present in node pool
	NodeLabels map[string]*string `json:"nodeLabels,omitempty"`

	// NodeTaints specifies the taints for nodes present in this agent pool.
	NodeTaints []string `json:"nodeTaints,omitempty"`

	// EnableAutoScaling - Whether to enable auto-scaler
	EnableAutoScaling bool `json:"enableAutoScaling,omitempty"`

	// AvailabilityZones represents the Availability zones for nodes in the AgentPool.
	AvailabilityZones []string

	// MaxPods specifies the kubelet --max-pods configuration for the agent pool.
	MaxPods *int32 `json:"maxPods,omitempty"`

	// OsDiskType specifies the OS disk type for each node in the pool. Allowed values are 'Ephemeral' and 'Managed'.
	OsDiskType *string `json:"osDiskType,omitempty"`

	// EnableUltraSSD enables the storage type UltraSSD_LRS for the agent pool.
	EnableUltraSSD *bool `json:"enableUltraSSD,omitempty"`

	// OSType specifies the operating system for the node pool. Allowed values are 'Linux' and 'Windows'
	OSType *string `json:"osType,omitempty"`

	// Headers is the list of headers to add to the HTTP requests to update this resource.
	Headers map[string]string

	// EnableNodePublicIP controls whether or not nodes in the agent pool each have a public IP address.
	EnableNodePublicIP *bool `json:"enableNodePublicIP,omitempty"`

	// NodePublicIPPrefixID specifies the public IP prefix resource ID which VM nodes should use IPs from.
	NodePublicIPPrefixID *string `json:"nodePublicIPPrefixID,omitempty"`

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
}

// ResourceName returns the name of the agent pool.
func (s *AgentPoolSpec) ResourceName() string {
	return s.Name
}

// ResourceGroupName returns the name of the resource group.
func (s *AgentPoolSpec) ResourceGroupName() string {
	return s.ResourceGroup
}

// OwnerResourceName is a no-op for agent pools.
func (s *AgentPoolSpec) OwnerResourceName() string {
	return s.Cluster
}

// CustomHeaders returns custom headers to be added to the Azure API calls.
func (s *AgentPoolSpec) CustomHeaders() map[string]string {
	return s.Headers
}

// Parameters returns the parameters for the agent pool.
func (s *AgentPoolSpec) Parameters(ctx context.Context, existing interface{}) (params interface{}, err error) {
	_, log, done := tele.StartSpanWithLogger(ctx, "agentpools.Service.Parameters")
	defer done()

	nodeLabels := s.NodeLabels
	if existing != nil {
		existingPool, ok := existing.(armcontainerservice.AgentPool)
		if !ok {
			return nil, errors.Errorf("%T is not an armcontainerservice.AgentPool", existing)
		}

		// agent pool already exists
		ps := *existingPool.Properties.ProvisioningState
		if ps != string(infrav1.Canceled) && ps != string(infrav1.Failed) && ps != string(infrav1.Succeeded) {
			msg := fmt.Sprintf("Unable to update existing agent pool in non terminal state. Agent pool must be in one of the following provisioning states: Canceled, Failed, or Succeeded. Actual state: %s", ps)
			return nil, azure.WithTransientError(errors.New(msg), 20*time.Second)
		}

		// Normalize individual agent pools to diff in case we need to update
		existingProfile := armcontainerservice.AgentPool{
			Properties: &armcontainerservice.ManagedClusterAgentPoolProfileProperties{
				Count:               existingPool.Properties.Count,
				OrchestratorVersion: existingPool.Properties.OrchestratorVersion,
				Mode:                existingPool.Properties.Mode,
				EnableAutoScaling:   existingPool.Properties.EnableAutoScaling,
				MinCount:            existingPool.Properties.MinCount,
				MaxCount:            existingPool.Properties.MaxCount,
				NodeLabels:          existingPool.Properties.NodeLabels,
				NodeTaints:          existingPool.Properties.NodeTaints,
				Tags:                existingPool.Properties.Tags,
				ScaleDownMode:       existingPool.Properties.ScaleDownMode,
				SpotMaxPrice:        existingPool.Properties.SpotMaxPrice,
				KubeletConfig:       existingPool.Properties.KubeletConfig,
			},
		}

		normalizedProfile := armcontainerservice.AgentPool{
			Properties: &armcontainerservice.ManagedClusterAgentPoolProfileProperties{
				Count:               &s.Replicas,
				OrchestratorVersion: s.Version,
				Mode:                azure.AliasOrNil[armcontainerservice.AgentPoolMode](&s.Mode),
				EnableAutoScaling:   ptr.To(s.EnableAutoScaling),
				MinCount:            s.MinCount,
				MaxCount:            s.MaxCount,
				NodeLabels:          s.NodeLabels,
				NodeTaints:          azure.PtrSlice(&s.NodeTaints),
				ScaleDownMode:       azure.AliasOrNil[armcontainerservice.ScaleDownMode](s.ScaleDownMode),
				Tags:                converters.TagsToMap(s.AdditionalTags),
			},
		}
		if len(normalizedProfile.Properties.NodeTaints) == 0 {
			normalizedProfile.Properties.NodeTaints = nil
		}

		if s.SpotMaxPrice != nil {
			normalizedProfile.Properties.SpotMaxPrice = ptr.To[float32](float32(s.SpotMaxPrice.AsApproximateFloat64()))
		}

		if s.KubeletConfig != nil {
			normalizedProfile.Properties.KubeletConfig = &armcontainerservice.KubeletConfig{
				CPUManagerPolicy:      s.KubeletConfig.CPUManagerPolicy,
				CPUCfsQuota:           s.KubeletConfig.CPUCfsQuota,
				CPUCfsQuotaPeriod:     s.KubeletConfig.CPUCfsQuotaPeriod,
				ImageGcHighThreshold:  s.KubeletConfig.ImageGcHighThreshold,
				ImageGcLowThreshold:   s.KubeletConfig.ImageGcLowThreshold,
				TopologyManagerPolicy: s.KubeletConfig.TopologyManagerPolicy,
				FailSwapOn:            s.KubeletConfig.FailSwapOn,
				ContainerLogMaxSizeMB: s.KubeletConfig.ContainerLogMaxSizeMB,
				ContainerLogMaxFiles:  s.KubeletConfig.ContainerLogMaxFiles,
				PodMaxPids:            s.KubeletConfig.PodMaxPids,
				AllowedUnsafeSysctls:  azure.PtrSlice(s.KubeletConfig.AllowedUnsafeSysctls),
			}
		}

		// When autoscaling is set, the count of the nodes differ based on the autoscaler and should not depend on the
		// count present in MachinePool or AzureManagedMachinePool, hence we should not make an update API call based
		// on difference in count.
		if s.EnableAutoScaling {
			normalizedProfile.Properties.Count = existingProfile.Properties.Count
		}

		// We do a just-in-time merge of existent kubernetes.azure.com-prefixed labels
		// So that we don't unintentionally delete them
		// See https://github.com/Azure/AKS/issues/3152
		nodeLabels = mergeSystemNodeLabels(normalizedProfile.Properties.NodeLabels, existingPool.Properties.NodeLabels)
		normalizedProfile.Properties.NodeLabels = nodeLabels

		// Compute a diff to check if we require an update
		diff := cmp.Diff(normalizedProfile, existingProfile)
		if diff == "" {
			// agent pool is up to date, nothing to do
			log.V(4).Info("no changes found between user-updated spec and existing spec")
			return nil, nil
		}
		log.V(4).Info("found a diff between the desired spec and the existing agentpool", "difference", diff)
	}

	availabilityZones := azure.PtrSlice(&s.AvailabilityZones)
	nodeTaints := azure.PtrSlice(&s.NodeTaints)
	var sku *string
	if s.SKU != "" {
		sku = &s.SKU
	}
	var spotMaxPrice *float32
	if s.SpotMaxPrice != nil {
		spotMaxPrice = ptr.To[float32](float32(s.SpotMaxPrice.AsApproximateFloat64()))
	}
	tags := converters.TagsToMap(s.AdditionalTags)
	if tags == nil {
		// Make sure we send a non-nil, empty map if AdditionalTags are nil as this tells AKS to delete any existing tags.
		tags = make(map[string]*string, 0)
	}
	var vnetSubnetID *string
	if s.VnetSubnetID != "" {
		vnetSubnetID = &s.VnetSubnetID
	}

	var kubeletConfig *armcontainerservice.KubeletConfig
	if s.KubeletConfig != nil {
		kubeletConfig = &armcontainerservice.KubeletConfig{
			CPUManagerPolicy:      s.KubeletConfig.CPUManagerPolicy,
			CPUCfsQuota:           s.KubeletConfig.CPUCfsQuota,
			CPUCfsQuotaPeriod:     s.KubeletConfig.CPUCfsQuotaPeriod,
			ImageGcHighThreshold:  s.KubeletConfig.ImageGcHighThreshold,
			ImageGcLowThreshold:   s.KubeletConfig.ImageGcLowThreshold,
			TopologyManagerPolicy: s.KubeletConfig.TopologyManagerPolicy,
			FailSwapOn:            s.KubeletConfig.FailSwapOn,
			ContainerLogMaxSizeMB: s.KubeletConfig.ContainerLogMaxSizeMB,
			ContainerLogMaxFiles:  s.KubeletConfig.ContainerLogMaxFiles,
			PodMaxPids:            s.KubeletConfig.PodMaxPids,
			AllowedUnsafeSysctls:  azure.PtrSlice(s.KubeletConfig.AllowedUnsafeSysctls),
		}
	}

	var linuxOSConfig *armcontainerservice.LinuxOSConfig
	if s.LinuxOSConfig != nil {
		linuxOSConfig = &armcontainerservice.LinuxOSConfig{
			SwapFileSizeMB:             s.LinuxOSConfig.SwapFileSizeMB,
			TransparentHugePageEnabled: (*string)(s.LinuxOSConfig.TransparentHugePageEnabled),
			TransparentHugePageDefrag:  (*string)(s.LinuxOSConfig.TransparentHugePageDefrag),
		}
		if s.LinuxOSConfig.Sysctls != nil {
			linuxOSConfig.Sysctls = &armcontainerservice.SysctlConfig{
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
				NetIPv4IPLocalPortRange:        s.LinuxOSConfig.Sysctls.NetIpv4IPLocalPortRange,
				NetIPv4NeighDefaultGcThresh1:   s.LinuxOSConfig.Sysctls.NetIpv4NeighDefaultGcThresh1,
				NetIPv4NeighDefaultGcThresh2:   s.LinuxOSConfig.Sysctls.NetIpv4NeighDefaultGcThresh2,
				NetIPv4NeighDefaultGcThresh3:   s.LinuxOSConfig.Sysctls.NetIpv4NeighDefaultGcThresh3,
				NetIPv4TCPFinTimeout:           s.LinuxOSConfig.Sysctls.NetIpv4TCPFinTimeout,
				NetIPv4TCPKeepaliveProbes:      s.LinuxOSConfig.Sysctls.NetIpv4TCPKeepaliveProbes,
				NetIPv4TCPKeepaliveTime:        s.LinuxOSConfig.Sysctls.NetIpv4TCPKeepaliveTime,
				NetIPv4TCPMaxSynBacklog:        s.LinuxOSConfig.Sysctls.NetIpv4TCPMaxSynBacklog,
				NetIPv4TCPMaxTwBuckets:         s.LinuxOSConfig.Sysctls.NetIpv4TCPMaxTwBuckets,
				NetIPv4TCPTwReuse:              s.LinuxOSConfig.Sysctls.NetIpv4TCPTwReuse,
				NetIPv4TcpkeepaliveIntvl:       s.LinuxOSConfig.Sysctls.NetIpv4TCPkeepaliveIntvl,
				NetNetfilterNfConntrackBuckets: s.LinuxOSConfig.Sysctls.NetNetfilterNfConntrackBuckets,
				NetNetfilterNfConntrackMax:     s.LinuxOSConfig.Sysctls.NetNetfilterNfConntrackMax,
				VMMaxMapCount:                  s.LinuxOSConfig.Sysctls.VMMaxMapCount,
				VMSwappiness:                   s.LinuxOSConfig.Sysctls.VMSwappiness,
				VMVfsCachePressure:             s.LinuxOSConfig.Sysctls.VMVfsCachePressure,
			}
		}
	}

	agentPool := armcontainerservice.AgentPool{
		Properties: &armcontainerservice.ManagedClusterAgentPoolProfileProperties{
			AvailabilityZones:    availabilityZones,
			Count:                &s.Replicas,
			EnableAutoScaling:    ptr.To(s.EnableAutoScaling),
			EnableUltraSSD:       s.EnableUltraSSD,
			KubeletConfig:        kubeletConfig,
			KubeletDiskType:      azure.AliasOrNil[armcontainerservice.KubeletDiskType]((*string)(s.KubeletDiskType)),
			MaxCount:             s.MaxCount,
			MaxPods:              s.MaxPods,
			MinCount:             s.MinCount,
			Mode:                 ptr.To(armcontainerservice.AgentPoolMode(s.Mode)),
			NodeLabels:           nodeLabels,
			NodeTaints:           nodeTaints,
			OrchestratorVersion:  s.Version,
			OSDiskSizeGB:         &s.OSDiskSizeGB,
			OSDiskType:           azure.AliasOrNil[armcontainerservice.OSDiskType](s.OsDiskType),
			OSType:               azure.AliasOrNil[armcontainerservice.OSType](s.OSType),
			ScaleSetPriority:     azure.AliasOrNil[armcontainerservice.ScaleSetPriority](s.ScaleSetPriority),
			ScaleDownMode:        azure.AliasOrNil[armcontainerservice.ScaleDownMode](s.ScaleDownMode),
			SpotMaxPrice:         spotMaxPrice,
			Type:                 ptr.To(armcontainerservice.AgentPoolTypeVirtualMachineScaleSets),
			VMSize:               sku,
			VnetSubnetID:         vnetSubnetID,
			EnableNodePublicIP:   s.EnableNodePublicIP,
			NodePublicIPPrefixID: s.NodePublicIPPrefixID,
			Tags:                 tags,
			EnableFIPS:           s.EnableFIPS,
			LinuxOSConfig:        linuxOSConfig,
		},
	}

	return agentPool, nil
}

// mergeSystemNodeLabels appends any kubernetes.azure.com-prefixed labels from the AKS label set
// into the local capz label set.
func mergeSystemNodeLabels(capz, aks map[string]*string) map[string]*string {
	ret := capz
	if ret == nil {
		ret = make(map[string]*string)
	}
	// Look for labels returned from the AKS node pool API that begin with kubernetes.azure.com
	for aksNodeLabelKey := range aks {
		if azureutil.IsAzureSystemNodeLabelKey(aksNodeLabelKey) {
			ret[aksNodeLabelKey] = aks[aksNodeLabelKey]
		}
	}
	// Preserve nil-ness of capz
	if capz == nil && len(ret) == 0 {
		ret = nil
	}
	return ret
}
