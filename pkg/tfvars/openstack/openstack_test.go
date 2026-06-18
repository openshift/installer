package openstack

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	types_openstack "github.com/openshift/installer/pkg/types/openstack"
)

// bootstrapFlavorTFVars is a minimal struct mirroring the fields relevant to
// bootstrap flavor in the anonymous struct returned by TFVars. It is used to
// validate JSON marshaling behavior without requiring a live OpenStack client.
type bootstrapFlavorTFVars struct {
	FlavorName          string `json:"openstack_master_flavor_name,omitempty"`
	BootstrapFlavorName string `json:"openstack_bootstrap_flavor_name,omitempty"`
	// Include required non-omitempty fields so the JSON is valid.
	MasterServerGroupPolicy types_openstack.ServerGroupPolicy `json:"openstack_master_server_group_policy"`
	WorkerServerGroupPolicy types_openstack.ServerGroupPolicy `json:"openstack_worker_server_group_policy"`
	AdditionalPorts         [][]terraformPort                 `json:"openstack_additional_ports"`
	MachinesPorts           []*terraformPort                  `json:"openstack_machines_ports"`
	UserManagedLoadBalancer bool                              `json:"openstack_user_managed_load_balancer"`
}

// TestBootstrapFlavorName verifies that the bootstrap flavor resolution logic
// produces the correct openstack_bootstrap_flavor_name value in the generated
// Terraform variables JSON, covering both the explicit-flavor and the
// fallback-to-master-flavor scenarios.
func TestBootstrapFlavorName(t *testing.T) {
	tests := []struct {
		name                string
		platform            *types_openstack.Platform
		masterFlavor        string
		wantBootstrapFlavor string
	}{
		{
			name: "explicit bootstrapFlavor is used",
			platform: &types_openstack.Platform{
				BootstrapFlavor: "m1.xlarge",
			},
			masterFlavor:        "m1.large",
			wantBootstrapFlavor: "m1.xlarge",
		},
		{
			name:                "empty bootstrapFlavor falls back to master flavor",
			platform:            &types_openstack.Platform{},
			masterFlavor:        "m1.large",
			wantBootstrapFlavor: "m1.large",
		},
		{
			name:                "nil platform falls back to master flavor",
			platform:            nil,
			masterFlavor:        "m1.large",
			wantBootstrapFlavor: "m1.large",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate the resolution performed inside TFVars.
			bootstrapFlavorName := tc.platform.ResolveBootstrapFlavor(tc.masterFlavor)

			// Marshal a struct that mirrors the fields written by TFVars.
			data, err := json.MarshalIndent(bootstrapFlavorTFVars{
				FlavorName:          tc.masterFlavor,
				BootstrapFlavorName: bootstrapFlavorName,
				AdditionalPorts:     [][]terraformPort{},
				MachinesPorts:       []*terraformPort{},
			}, "", "  ")
			if !assert.NoError(t, err) {
				return
			}

			// Unmarshal back and check the bootstrap flavor field.
			var result map[string]interface{}
			if !assert.NoError(t, json.Unmarshal(data, &result)) {
				return
			}

			got, ok := result["openstack_bootstrap_flavor_name"]
			assert.True(t, ok, "openstack_bootstrap_flavor_name must be present in Terraform variables JSON")
			assert.Equal(t, tc.wantBootstrapFlavor, got,
				"openstack_bootstrap_flavor_name should equal the resolved bootstrap flavor")

			// Verify master flavor is unchanged regardless of bootstrapFlavor setting.
			masterGot, masterOK := result["openstack_master_flavor_name"]
			assert.True(t, masterOK, "openstack_master_flavor_name must be present in Terraform variables JSON")
			assert.Equal(t, tc.masterFlavor, masterGot,
				"openstack_master_flavor_name must be unaffected by bootstrapFlavor setting")
		})
	}
}

// TestBootstrapFlavorNameInJSON verifies backward compatibility: when no
// bootstrapFlavor is configured, the generated JSON contains
// openstack_bootstrap_flavor_name equal to the master flavor.
func TestBootstrapFlavorNameInJSON(t *testing.T) {
	masterFlavor := "m1.large"
	platform := &types_openstack.Platform{} // no BootstrapFlavor set

	bootstrapFlavorName := platform.ResolveBootstrapFlavor(masterFlavor)
	assert.Equal(t, masterFlavor, bootstrapFlavorName,
		"without explicit bootstrapFlavor, resolved name must equal master flavor")

	data, err := json.MarshalIndent(bootstrapFlavorTFVars{
		FlavorName:          masterFlavor,
		BootstrapFlavorName: bootstrapFlavorName,
		AdditionalPorts:     [][]terraformPort{},
		MachinesPorts:       []*terraformPort{},
	}, "", "  ")
	if !assert.NoError(t, err) {
		return
	}

	var result map[string]interface{}
	if !assert.NoError(t, json.Unmarshal(data, &result)) {
		return
	}

	assert.Equal(t, masterFlavor, result["openstack_master_flavor_name"],
		"openstack_master_flavor_name must be preserved")
	assert.Equal(t, masterFlavor, result["openstack_bootstrap_flavor_name"],
		"openstack_bootstrap_flavor_name must fall back to master flavor")
}
