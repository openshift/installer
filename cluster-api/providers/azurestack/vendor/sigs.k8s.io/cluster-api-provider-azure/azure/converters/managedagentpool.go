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
	// NOTE: when the hub API version is updated, verify the
	// ManagedClusterAgentPoolProfile below has every field defined. If a field
	// isn't defined, the agent pool will be created with a zero/null value, and
	// then updated to the user-defined value. If the field is immutable, this
	// update will fail. The linter should catch if there are missing fields,
	// but verify that check is actually working.
	asocontainerservicev1hub "github.com/Azure/azure-service-operator/v2/api/containerservice/v1api20240901/storage"
	"k8s.io/utils/ptr"
)

// AgentPoolToManagedClusterAgentPoolProfile converts a AgentPoolSpec to an Azure SDK ManagedClusterAgentPoolProfile used in managedcluster reconcile.
func AgentPoolToManagedClusterAgentPoolProfile(pool *asocontainerservicev1hub.ManagedClustersAgentPool) asocontainerservicev1hub.ManagedClusterAgentPoolProfile {
	properties := pool.Spec
	agentPool := asocontainerservicev1hub.ManagedClusterAgentPoolProfile{
		AvailabilityZones:                 properties.AvailabilityZones,
		CapacityReservationGroupReference: properties.CapacityReservationGroupReference,
		Count:                             properties.Count,
		CreationData:                      properties.CreationData,
		EnableAutoScaling:                 properties.EnableAutoScaling,
		EnableEncryptionAtHost:            properties.EnableEncryptionAtHost,
		EnableFIPS:                        properties.EnableFIPS,
		EnableNodePublicIP:                properties.EnableNodePublicIP,
		EnableUltraSSD:                    properties.EnableUltraSSD,
		GpuInstanceProfile:                properties.GpuInstanceProfile,
		HostGroupReference:                properties.HostGroupReference,
		KubeletConfig:                     properties.KubeletConfig,
		KubeletDiskType:                   properties.KubeletDiskType,
		LinuxOSConfig:                     properties.LinuxOSConfig,
		MaxCount:                          properties.MaxCount,
		MaxPods:                           properties.MaxPods,
		MinCount:                          properties.MinCount,
		Mode:                              properties.Mode,
		Name:                              ptr.To(pool.AzureName()),
		NetworkProfile:                    properties.NetworkProfile,
		NodeLabels:                        properties.NodeLabels,
		NodePublicIPPrefixReference:       properties.NodePublicIPPrefixReference,
		NodeTaints:                        properties.NodeTaints,
		OrchestratorVersion:               properties.OrchestratorVersion,
		OsDiskSizeGB:                      properties.OsDiskSizeGB,
		OsDiskType:                        properties.OsDiskType,
		OsSKU:                             properties.OsSKU,
		OsType:                            properties.OsType,
		PodSubnetReference:                properties.PodSubnetReference,
		PowerState:                        properties.PowerState,
		PropertyBag:                       properties.PropertyBag,
		ProximityPlacementGroupReference:  properties.ProximityPlacementGroupReference,
		ScaleDownMode:                     properties.ScaleDownMode,
		ScaleSetEvictionPolicy:            properties.ScaleSetEvictionPolicy,
		ScaleSetPriority:                  properties.ScaleSetPriority,
		SecurityProfile:                   properties.SecurityProfile,
		SpotMaxPrice:                      properties.SpotMaxPrice,
		Tags:                              properties.Tags,
		Type:                              properties.Type,
		UpgradeSettings:                   properties.UpgradeSettings,
		VmSize:                            properties.VmSize,
		VnetSubnetReference:               properties.VnetSubnetReference,
		WindowsProfile:                    properties.WindowsProfile,
		WorkloadRuntime:                   properties.WorkloadRuntime,
	}
	return agentPool
}
