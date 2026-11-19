package vsphere

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSet_PreservesZonesFromDefaultMachinePlatform tests that zones from defaultMachinePlatform
// are preserved when pool-specific platform config is set without zones.
// This reproduces the bug reported in OCPBUGS-62209.
func TestSet_PreservesZonesFromDefaultMachinePlatform(t *testing.T) {
	// Simulate the scenario from the bug report:
	// 1. defaultMachinePlatform has zones defined
	// 2. controlPlane/compute platform has CPU/memory but no zones
	// 3. Zones from defaultMachinePlatform should be preserved

	// Create default machine pool platform (simulating defaultVSphereMachinePoolPlatform())
	mpool := MachinePool{
		NumCPUs:           4,
		NumCoresPerSocket: 4,
		MemoryMiB:         16384,
		OSDisk: OSDisk{
			DiskSizeGB: 120,
		},
	}

	// Create defaultMachinePlatform with zones (from install-config)
	defaultMachinePlatform := &MachinePool{
		Zones: []string{"us-east-1a"},
	}

	// Create pool-specific config with CPU/memory but no zones
	// This simulates: controlPlane.platform.vsphere = {cpus: 8, memoryMB: 32768}
	poolSpecificConfig := &MachinePool{
		NumCPUs:   8,
		MemoryMiB: 32768,
	}

	// Apply defaults first (this should set zones)
	mpool.Set(defaultMachinePlatform)
	assert.Equal(t, []string{"us-east-1a"}, mpool.Zones, "Zones should be set from defaultMachinePlatform")

	// Apply pool-specific config (this should NOT clear zones)
	mpool.Set(poolSpecificConfig)

	// BUG: This test will fail if zones are being lost
	assert.Equal(t, []string{"us-east-1a"}, mpool.Zones, "Zones should be preserved after applying pool-specific config")
	assert.Equal(t, int32(8), mpool.NumCPUs, "NumCPUs should be updated from pool-specific config")
	assert.Equal(t, int64(32768), mpool.MemoryMiB, "MemoryMiB should be updated from pool-specific config")
}

// TestSet_ExplicitEmptyZonesOverwriteDefault tests that if a pool explicitly
// sets zones to empty, it should NOT overwrite zones from default.
func TestSet_ExplicitEmptyZones(t *testing.T) {
	mpool := MachinePool{}

	// Set zones via default
	defaultPlatform := &MachinePool{
		Zones: []string{"us-east-1a", "us-east-1b"},
	}
	mpool.Set(defaultPlatform)
	assert.Equal(t, []string{"us-east-1a", "us-east-1b"}, mpool.Zones)

	// Apply config with empty zones (nil)
	poolConfig := &MachinePool{
		NumCPUs: 8,
		Zones:   nil, // explicitly nil
	}
	mpool.Set(poolConfig)

	// Zones should be preserved because len(nil) == 0
	assert.Equal(t, []string{"us-east-1a", "us-east-1b"}, mpool.Zones, "Zones should be preserved when source has nil zones")

	// Apply config with empty array
	poolConfig2 := &MachinePool{
		NumCPUs: 16,
		Zones:   []string{}, // explicitly empty array
	}
	mpool.Set(poolConfig2)

	// Zones should still be preserved because len([]string{}) == 0
	assert.Equal(t, []string{"us-east-1a", "us-east-1b"}, mpool.Zones, "Zones should be preserved when source has empty zones array")
}

// TestSet_ZonesCanBeOverwritten tests that zones CAN be overwritten if explicitly set
func TestSet_ZonesCanBeOverwritten(t *testing.T) {
	mpool := MachinePool{
		Zones: []string{"us-east-1a"},
	}

	// Overwrite with different zones
	newConfig := &MachinePool{
		Zones: []string{"us-west-1a", "us-west-1b"},
	}
	mpool.Set(newConfig)

	assert.Equal(t, []string{"us-west-1a", "us-west-1b"}, mpool.Zones, "Zones should be overwritten when explicitly set")
}
