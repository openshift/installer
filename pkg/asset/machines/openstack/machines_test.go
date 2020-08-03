package openstack_test

import (
	"fmt"
	"testing"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset/machines/openstack"
	"github.com/openshift/installer/pkg/types"
	openstackTypes "github.com/openshift/installer/pkg/types/openstack"
)

func newConfig() *types.InstallConfig {
	return &types.InstallConfig{
		AdditionalTrustBundle: "",
		SSHKey:                "",
		BaseDomain:            "basedomain.example.com",
		Networking:            nil,
		ControlPlane:          nil,
		Compute:               nil,
		Platform: types.Platform{
			OpenStack: &openstackTypes.Platform{
				Region:                 "RegionOne",
				DefaultMachinePlatform: nil,
				Cloud:                  "exampleCloud",
				ExternalNetwork:        "external",
				FlavorName:             "defaultFlavor",
				LbFloatingIP:           "203.0.113.23",
				IngressFloatingIP:      "203.0.113.19",
				ExternalDNS:            nil,
				ClusterOSImage:         "",
				APIVIP:                 "",
				IngressVIP:             "",
				MachinesSubnet:         "",
				// TrunkSupport
				// OctaviaSupport
			},
		},
		PullSecret:          "",
		Proxy:               nil,
		ImageContentSources: nil,
		// Publish
		// FIPS
		// CredentialsMode
	}
}

func newPool(options ...func(*types.MachinePool)) *types.MachinePool {
	pool := types.MachinePool{
		Name:     "MachinePoolName",
		Replicas: nil,
		Platform: types.MachinePoolPlatform{
			OpenStack: &openstackTypes.MachinePool{
				FlavorName:           "poolFlavor",
				RootVolume:           nil,
				AdditionalNetworkIDs: nil,
				Zones:                nil,
			},
		},
		// Hyperthreading
		// Architecture
	}

	for _, apply := range options {
		apply(&pool)
	}

	return &pool
}

func withZones(zones ...string) func(*types.MachinePool) {
	return func(pool *types.MachinePool) {
		pool.Platform.OpenStack.Zones = zones
	}
}

func withReplicas(replicas int64) func(*types.MachinePool) {
	return func(pool *types.MachinePool) {
		pool.Replicas = &replicas
	}
}

func TestMachines(t *testing.T) {
	type checkFunc func([]machineapi.Machine, error) error
	check := func(fns ...checkFunc) []checkFunc { return fns }

	hasMachines := func(want int) checkFunc {
		return func(machines []machineapi.Machine, _ error) error {
			if have := len(machines); have != want {
				return fmt.Errorf("expected %d machines, found %d", want, have)
			}
			return nil
		}
	}

	hasNilError := func(_ []machineapi.Machine, have error) error {
		if have != nil {
			return fmt.Errorf("found unexpected error")
		}
		return nil
	}

	for _, tc := range [...]struct {
		name           string
		clusterID      string
		config         *types.InstallConfig
		pool           *types.MachinePool
		osImage        string
		role           string
		userDataSecret string
		trunkSupport   bool
		checks         []checkFunc
	}{
		{
			name:      "Many replicas",
			clusterID: "cluster_id",
			config:    newConfig(),
			pool: newPool(
				withReplicas(1 << 15),
			),
			osImage:        "os_image",
			role:           "role",
			userDataSecret: "",
			trunkSupport:   true,
			checks: check(
				hasMachines(1<<15),
				hasNilError,
			),
		},
		{
			name:      "Two availability zones",
			clusterID: "cluster_id",
			config:    newConfig(),
			pool: newPool(
				withZones("zone1", "zone2"),
				withReplicas(2),
			),
			osImage:        "os_image",
			role:           "role",
			userDataSecret: "",
			trunkSupport:   true,
			checks: check(
				hasMachines(2),
				hasNilError,
			),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			machines, err := openstack.Machines(
				tc.clusterID,
				tc.config,
				tc.pool,
				tc.osImage,
				tc.role,
				tc.userDataSecret,
				tc.trunkSupport,
			)

			for _, check := range tc.checks {
				if e := check(machines, err); e != nil {
					t.Error(e)
				}
			}

			if l, c := len(machines), cap(machines); l != c {
				t.Errorf("the machines slice was not properly allocated. Len: %d, Cap: %d", l, c)
			}
		})
	}
}
