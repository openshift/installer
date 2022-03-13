// Copyright 2021 Google LLC. All Rights Reserved.
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

func betaWorkerConfigToGaWorkerConfig(r *WorkerPool, c *WorkerPoolWorkerConfig) *WorkerPoolPrivatePoolV1ConfigWorkerConfig {
	if c == nil {
		return nil
	}
	cfg := &WorkerPoolPrivatePoolV1ConfigWorkerConfig{}
	cfg.DiskSizeGb = c.DiskSizeGb
	cfg.MachineType = c.MachineType
	// we need to *not* send this field, unfortunately.
	c.DiskSizeGb = nil
	c.MachineType = nil
	// Now, a little messy - we are going to IGNORE NoExternalIP.
	// It stays set to whatever its old value was.  We set
	// c.empty to true, trusting that a call to the function right below
	// this one, parsing out the network config, is going to read it.
	// We don't express an opinion about whether that should happen before
	// or after this current call.

	c.empty = true
	return cfg
}

func betaNetworkConfigToGaNetworkConfig(r *WorkerPool, c *WorkerPoolNetworkConfig) *WorkerPoolPrivatePoolV1ConfigNetworkConfig {
	if c == nil {
		return nil
	}
	cfg := &WorkerPoolPrivatePoolV1ConfigNetworkConfig{}
	cfg.PeeredNetwork = c.PeeredNetwork
	// The counterpart to the messy bit above - we need to translate this boolean in an unrelated
	// field into a tri-state enum in this field.  Gross!
	if r.WorkerConfig != nil {
		if r.WorkerConfig.NoExternalIP == nil {
			cfg.EgressOption = WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef("EGRESS_OPTION_UNSPECIFIED")
		} else if *r.WorkerConfig.NoExternalIP {
			cfg.EgressOption = WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef("NO_PUBLIC_EGRESS")
		} else {
			cfg.EgressOption = WorkerPoolPrivatePoolV1ConfigNetworkConfigEgressOptionEnumRef("PUBLIC_EGRESS")
		}
		r.WorkerConfig.NoExternalIP = nil
	}
	c.PeeredNetwork = nil
	c.empty = true
	return cfg
}

func gaNetworkConfigToBetaNetworkConfig(r *WorkerPool, c *WorkerPoolPrivatePoolV1ConfigNetworkConfig) *WorkerPoolNetworkConfig {
	if c == nil {
		return nil
	}
	if c.PeeredNetwork == nil {
		return EmptyWorkerPoolNetworkConfig
	}

	cfg := &WorkerPoolNetworkConfig{}
	cfg.PeeredNetwork = c.PeeredNetwork
	return cfg
}

func gaWorkerConfigToBetaWorkerConfig(r *WorkerPool, c *WorkerPoolPrivatePoolV1ConfigWorkerConfig) *WorkerPoolWorkerConfig {
	if c == nil {
		return nil
	}
	cfg := &WorkerPoolWorkerConfig{}
	cfg.DiskSizeGb = c.DiskSizeGb
	cfg.MachineType = c.MachineType
	if r.PrivatePoolV1Config != nil && r.PrivatePoolV1Config.NetworkConfig != nil {
		if r.PrivatePoolV1Config.NetworkConfig.EgressOption == nil {
			cfg.NoExternalIP = nil
		} else if string(*r.PrivatePoolV1Config.NetworkConfig.EgressOption) == "NO_PUBLIC_EGRESS" {
			cfg.NoExternalIP = dcl.Bool(true)
		} else if string(*r.PrivatePoolV1Config.NetworkConfig.EgressOption) == "PUBLIC_EGRESS" {
			cfg.NoExternalIP = dcl.Bool(false)
		}
	}
	return cfg
}
