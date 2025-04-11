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
	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20231001/storage"
	"k8s.io/utils/ptr"
)

// AgentPoolToManagedClusterAgentPoolProfile converts a AgentPoolSpec to an Azure SDK ManagedClusterAgentPoolProfile used in managedcluster reconcile.
func AgentPoolToManagedClusterAgentPoolProfile(pool *asocontainerservicev1hub.ManagedClustersAgentPool) asocontainerservicev1hub.ManagedClusterAgentPoolProfile {
	properties := pool.Spec
	agentPool := asocontainerservicev1hub.ManagedClusterAgentPoolProfile{
		Name:                        ptr.To(pool.AzureName()),
		VmSize:                      properties.VmSize,
		OsType:                      properties.OsType,
		OsDiskSizeGB:                properties.OsDiskSizeGB,
		Count:                       properties.Count,
		Type:                        properties.Type,
		OrchestratorVersion:         properties.OrchestratorVersion,
		VnetSubnetReference:         properties.VnetSubnetReference,
		Mode:                        properties.Mode,
		EnableAutoScaling:           properties.EnableAutoScaling,
		MaxCount:                    properties.MaxCount,
		MinCount:                    properties.MinCount,
		NodeTaints:                  properties.NodeTaints,
		AvailabilityZones:           properties.AvailabilityZones,
		MaxPods:                     properties.MaxPods,
		OsDiskType:                  properties.OsDiskType,
		NodeLabels:                  properties.NodeLabels,
		EnableUltraSSD:              properties.EnableUltraSSD,
		EnableNodePublicIP:          properties.EnableNodePublicIP,
		NodePublicIPPrefixReference: properties.NodePublicIPPrefixReference,
		ScaleSetPriority:            properties.ScaleSetPriority,
		ScaleDownMode:               properties.ScaleDownMode,
		SpotMaxPrice:                properties.SpotMaxPrice,
		Tags:                        properties.Tags,
		KubeletDiskType:             properties.KubeletDiskType,
		LinuxOSConfig:               properties.LinuxOSConfig,
		EnableFIPS:                  properties.EnableFIPS,
		EnableEncryptionAtHost:      properties.EnableEncryptionAtHost,
	}
	if properties.KubeletConfig != nil {
		agentPool.KubeletConfig = properties.KubeletConfig
	}
	return agentPool
}
