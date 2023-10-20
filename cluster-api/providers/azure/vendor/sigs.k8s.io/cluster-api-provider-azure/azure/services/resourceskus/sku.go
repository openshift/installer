/*
Copyright 2020 The Kubernetes Authors.

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

package resourceskus

import (
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	"github.com/pkg/errors"
	"k8s.io/utils/ptr"
)

// SKU is a thin layer over the Azure resource SKU API to better introspect capabilities.
type SKU armcompute.ResourceSKU

// ResourceType models available resource types as a set of known string constants.
type ResourceType string

const (
	// VirtualMachines is a convenience constant to filter resource SKUs to only include VMs.
	VirtualMachines ResourceType = "virtualMachines"
	// Disks is a convenience constant to filter resource SKUs to only include disks.
	Disks ResourceType = "disks"
	// AvailabilitySets is a convenience constant to filter resource SKUs to only include availability sets.
	AvailabilitySets ResourceType = "availabilitySets"
)

// Supported models an enum of possible boolean values for resource support in the Azure API.
type Supported string

const (
	// CapabilitySupported is the value returned by this API from Azure when the capability is supported.
	CapabilitySupported Supported = "True"
	// CapabilityUnsupported is the value returned by this API from Azure when the capability is unsupported.
	CapabilityUnsupported Supported = "False"
)

const (
	// EphemeralOSDisk identifies the capability for ephemeral os support.
	EphemeralOSDisk = "EphemeralOSDiskSupported"
	// AcceleratedNetworking identifies the capability for accelerated networking support.
	AcceleratedNetworking = "AcceleratedNetworkingEnabled"
	// VCPUs identifies the capability for the number of vCPUS.
	VCPUs = "vCPUs"
	// MemoryGB identifies the capability for memory Size.
	MemoryGB = "MemoryGB"
	// MinimumVCPUS is the minimum vCPUS allowed.
	MinimumVCPUS = 2
	// MinimumMemory is the minimum memory allowed.
	MinimumMemory = 2
	// EncryptionAtHost identifies the capability for encryption at host.
	EncryptionAtHost = "EncryptionAtHostSupported"
	// MaximumPlatformFaultDomainCount identifies the maximum fault domain count for an availability set in a region.
	MaximumPlatformFaultDomainCount = "MaximumPlatformFaultDomainCount"
	// UltraSSDAvailable identifies the capability for the support of UltraSSD data disks.
	UltraSSDAvailable = "UltraSSDAvailable"
	// TrustedLaunchDisabled identifies the absence of the trusted launch capability.
	TrustedLaunchDisabled = "TrustedLaunchDisabled"
	// ConfidentialComputingType identifies the capability for confidentical computing.
	ConfidentialComputingType = "ConfidentialComputingType"
	// CPUArchitectureType identifies the capability for cpu architecture.
	CPUArchitectureType = "CpuArchitectureType"
)

// HasCapability return true for a capability which can be either
// supported or not. Examples include "EphemeralOSDiskSupported",
// "UltraSSDAvavailable" "EncryptionAtHostSupported",
// "AcceleratedNetworkingEnabled", and "RdmaEnabled".
func (s SKU) HasCapability(name string) bool {
	if s.Capabilities != nil {
		for _, capability := range s.Capabilities {
			if capability.Name != nil && *capability.Name == name {
				if capability.Value != nil && strings.EqualFold(*capability.Value, string(CapabilitySupported)) {
					return true
				}
			}
		}
	}
	return false
}

// HasCapabilityWithCapacity returns true when the provided resource
// exposes a numeric capability and the maximum value exposed by that
// capability exceeds the value requested by the user. Examples include
// "MaxResourceVolumeMB", "OSVhdSizeMB", "vCPUs",
// "MemoryGB","MaxDataDiskCount", "CombinedTempDiskAndCachedIOPS",
// "CombinedTempDiskAndCachedReadBytesPerSecond",
// "CombinedTempDiskAndCachedWriteBytesPerSecond", "UncachedDiskIOPS",
// and "UncachedDiskBytesPerSecond".
func (s SKU) HasCapabilityWithCapacity(name string, value int64) (bool, error) {
	if s.Capabilities == nil {
		return false, nil
	}

	for _, capability := range s.Capabilities {
		if capability.Name == nil || *capability.Name != name || capability.Value == nil {
			continue
		}

		intVal, err := strconv.ParseInt(*capability.Value, 10, 64)
		if err != nil {
			return false, errors.Wrapf(err, "failed to parse string '%s' as int64", *capability.Value)
		}

		if intVal >= value {
			return true, nil
		}
	}

	return false, nil
}

// GetCapability gets the value assigned to the given capability.
// Eg. MaximumPlatformFaultDomainCount -> "3" will return "3" for the capability "MaximumPlatformFaultDomainCount".
func (s SKU) GetCapability(name string) (string, bool) {
	if s.Capabilities != nil {
		for _, capability := range s.Capabilities {
			if capability.Name != nil && *capability.Name == name {
				return *capability.Value, true
			}
		}
	}
	return "", false
}

// HasLocationCapability returns true if the provided resource supports the location capability.
func (s SKU) HasLocationCapability(capabilityName, location, zone string) bool {
	if s.LocationInfo == nil {
		return false
	}

	for _, info := range s.LocationInfo {
		if info.Location == nil || *info.Location != location || info.ZoneDetails == nil {
			continue
		}

		for _, zoneDetail := range info.ZoneDetails {
			if zoneDetail.Capabilities == nil {
				continue
			}

			for _, capability := range zoneDetail.Capabilities {
				if capability.Name != nil && *capability.Name == capabilityName {
					for _, name := range zoneDetail.Name {
						if ptr.Deref(name, "") == zone {
							return true
						}
					}

					return false
				}
			}
		}
	}
	return false
}
