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

package converters

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v4"
)

// AgentPoolToManagedClusterAgentPoolProfile converts a AgentPoolSpec to an Azure SDK ManagedClusterAgentPoolProfile used in managedcluster reconcile.
func AgentPoolToManagedClusterAgentPoolProfile(pool armcontainerservice.AgentPool) armcontainerservice.ManagedClusterAgentPoolProfile {
	properties := pool.Properties
	agentPool := armcontainerservice.ManagedClusterAgentPoolProfile{
		Name:                 pool.Name, // Note: if converting from agentPoolSpec.Parameters(), this field will not be set
		VMSize:               properties.VMSize,
		OSType:               properties.OSType,
		OSDiskSizeGB:         properties.OSDiskSizeGB,
		Count:                properties.Count,
		Type:                 properties.Type,
		OrchestratorVersion:  properties.OrchestratorVersion,
		VnetSubnetID:         properties.VnetSubnetID,
		Mode:                 properties.Mode,
		EnableAutoScaling:    properties.EnableAutoScaling,
		MaxCount:             properties.MaxCount,
		MinCount:             properties.MinCount,
		NodeTaints:           properties.NodeTaints,
		AvailabilityZones:    properties.AvailabilityZones,
		MaxPods:              properties.MaxPods,
		OSDiskType:           properties.OSDiskType,
		NodeLabels:           properties.NodeLabels,
		EnableUltraSSD:       properties.EnableUltraSSD,
		EnableNodePublicIP:   properties.EnableNodePublicIP,
		NodePublicIPPrefixID: properties.NodePublicIPPrefixID,
		ScaleSetPriority:     properties.ScaleSetPriority,
		ScaleDownMode:        properties.ScaleDownMode,
		SpotMaxPrice:         properties.SpotMaxPrice,
		Tags:                 properties.Tags,
		KubeletDiskType:      properties.KubeletDiskType,
		LinuxOSConfig:        properties.LinuxOSConfig,
		EnableFIPS:           properties.EnableFIPS,
	}
	if properties.KubeletConfig != nil {
		agentPool.KubeletConfig = properties.KubeletConfig
	}
	return agentPool
}
