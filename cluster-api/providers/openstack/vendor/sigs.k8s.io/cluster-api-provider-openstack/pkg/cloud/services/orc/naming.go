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

package orc

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// hash6 returns a 6-character hex hash of the input string.
// Used to generate deterministic, unique suffixes for ORC resource names
// when multiple resources of the same type may be needed per server
// (e.g., multiple networks referenced by different ports).
func hash6(input string) string {
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:3])
}

// Deterministic ORC resource naming functions.
// All names are derived from the OpenStackServer name to ensure
// uniqueness and predictability.

// ImageName returns the ORC Image name for a server.
func ImageName(serverName string) string {
	return fmt.Sprintf("%s-image", serverName)
}

// FlavorName returns the ORC Flavor name for a server.
func FlavorName(serverName string) string {
	return fmt.Sprintf("%s-flavor", serverName)
}

// KeyPairName returns the ORC KeyPair name for a server.
func KeyPairName(serverName string) string {
	return fmt.Sprintf("%s-keypair", serverName)
}

// ServerGroupORCName returns the ORC ServerGroup name for a server.
func ServerGroupORCName(serverName string) string {
	return fmt.Sprintf("%s-servergroup", serverName)
}

// NetworkORCName returns the ORC Network name for a server, using a
// hash of the param key for deduplication.
func NetworkORCName(serverName, key string) string {
	return fmt.Sprintf("%s-net-%s", serverName, hash6(key))
}

// SubnetORCName returns the ORC Subnet name for a server, using a
// hash of the param key for deduplication.
func SubnetORCName(serverName, key string) string {
	return fmt.Sprintf("%s-subnet-%s", serverName, hash6(key))
}

// SecurityGroupORCName returns the ORC SecurityGroup name for a server,
// using a hash of the param key for deduplication.
func SecurityGroupORCName(serverName, key string) string {
	return fmt.Sprintf("%s-sg-%s", serverName, hash6(key))
}

// PortORCName returns the ORC Port name for a server port at the given index.
func PortORCName(serverName string, index int) string {
	return fmt.Sprintf("%s-port-%d", serverName, index)
}

// TrunkORCName returns the ORC Trunk name for a server trunk at the given index.
func TrunkORCName(serverName string, index int) string {
	return fmt.Sprintf("%s-trunk-%d", serverName, index)
}

// RootVolumeName returns the ORC Volume name for the root volume.
func RootVolumeName(serverName string) string {
	return fmt.Sprintf("%s-vol-root", serverName)
}

// AdditionalVolumeName returns the ORC Volume name for an additional block device.
func AdditionalVolumeName(serverName, deviceName string) string {
	return fmt.Sprintf("%s-vol-%s", serverName, deviceName)
}

// VolumeTypeORCName returns the ORC VolumeType name for a server,
// using a hash of the volume type name for deduplication.
func VolumeTypeORCName(serverName, key string) string {
	return fmt.Sprintf("%s-voltype-%s", serverName, hash6(key))
}

// ServerName returns the ORC Server name, which is the same as the
// OpenStackServer name.
func ServerName(serverName string) string {
	return serverName
}

// Key generation functions for deduplication.
// These produce a stable string key from a CAPO parameter type,
// used to generate deterministic ORC resource names via hash6().

// NetworkParamKey returns a stable string key for a NetworkParam.
func NetworkParamKey(param infrav1.NetworkParam) string {
	if param.ID != nil {
		return "id:" + *param.ID
	}
	if param.Filter != nil {
		b, _ := json.Marshal(param.Filter)
		return "filter:" + string(b)
	}
	return ""
}

// SubnetParamKey returns a stable string key for a SubnetParam.
func SubnetParamKey(param infrav1.SubnetParam) string {
	if param.ID != nil {
		return "id:" + *param.ID
	}
	if param.Filter != nil {
		b, _ := json.Marshal(param.Filter)
		return "filter:" + string(b)
	}
	return ""
}

// SecurityGroupParamKey returns a stable string key for a SecurityGroupParam.
func SecurityGroupParamKey(param infrav1.SecurityGroupParam) string {
	if param.ID != nil {
		return "id:" + *param.ID
	}
	if param.Filter != nil {
		b, _ := json.Marshal(param.Filter)
		return "filter:" + string(b)
	}
	return ""
}
