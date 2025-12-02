// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package cloudbuild contains utility methods to serve the DCL in handling odd situations in the cloudbuild API.
package cloudbuild

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// betaToGaPrivatePool is populating GA specific PrivatePoolV1Config values and setting WorkerConfig and NetworkConfig to nil.
// r.PrivatePoolV1Config and c points to the same object.
func betaToGaPrivatePool(r *WorkerPool, c *WorkerPoolPrivatePoolV1Config) *WorkerPoolPrivatePoolV1Config {
	cfgWorkerConfig := &WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	cfgNetworkConfig := &WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	if r.WorkerConfig != nil {
		cfgWorkerConfig.DiskSizeGb = r.WorkerConfig.DiskSizeGb
		cfgWorkerConfig.MachineType = r.WorkerConfig.MachineType
		cfgNetworkConfig.EgressOption = noExternalIPEnum(r.WorkerConfig.NoExternalIP)
	}
	if r.NetworkConfig != nil {
		cfgNetworkConfig.PeeredNetwork = r.NetworkConfig.PeeredNetwork
		cfgNetworkConfig.PeeredNetworkIPRange = r.NetworkConfig.PeeredNetworkIPRange
	}

	cfg := &WorkerPoolPrivatePoolV1Config{}
	cfg.WorkerConfig = cfgWorkerConfig
	cfg.NetworkConfig = cfgNetworkConfig

	r.WorkerConfig = nil
	r.NetworkConfig = nil
	return cfg
}

// gaToBetaPrivatePool is populating beta specific values (WorkerConfig and NetworkConfig) and setting PrivatePoolV1Config to nil.
// r.PrivatePoolV1Config and c points to the same object.
func gaToBetaPrivatePool(r *WorkerPool, c *WorkerPoolPrivatePoolV1Config) *WorkerPoolPrivatePoolV1Config {
	if c == nil {
		return nil
	}

	if c.WorkerConfig != nil && r.WorkerConfig == nil {
		r.WorkerConfig = &WorkerPoolWorkerConfig{
			DiskSizeGb:   c.WorkerConfig.DiskSizeGb,
			MachineType:  c.WorkerConfig.MachineType,
			NoExternalIP: noExternalIPBoolean(c.NetworkConfig),
		}

	}
	if c.NetworkConfig != nil && c.NetworkConfig.PeeredNetwork != nil && r.NetworkConfig == nil {
		r.NetworkConfig = &WorkerPoolNetworkConfig{
			PeeredNetwork:        c.NetworkConfig.PeeredNetwork,
			PeeredNetworkIPRange: c.NetworkConfig.PeeredNetworkIPRange,
		}
	}

	r.PrivatePoolV1Config = nil
	return nil
}

func noExternalIPBoolean(networkConfig *WorkerPoolPrivatePoolV1ConfigNetworkConfig) *bool {
	if networkConfig == nil || networkConfig.EgressOption == nil {
		return nil
	}
	if string(*networkConfig.EgressOption) == "NO_PUBLIC_EGRESS" {
		return dcl.Bool(true)
	}
	if string(*networkConfig.EgressOption) == "PUBLIC_EGRESS" {
		return dcl.Bool(false)
	}
	return nil
}

func noExternalIPEnum(noExternalIP *bool) *WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnum {
	if noExternalIP == nil {
		return WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef("EGRESS_OPTION_UNSPECIFIED")
	}
	if *noExternalIP {
		return WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef("NO_PUBLIC_EGRESS")
	}
	return WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef("PUBLIC_EGRESS")
}
