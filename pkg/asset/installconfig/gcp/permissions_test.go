package gcp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/types"
	dnstypes "github.com/openshift/installer/pkg/types/dns"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func kmsKey() *gcptypes.EncryptionKeyReference {
	return &gcptypes.EncryptionKeyReference{
		KMSKey: &gcptypes.KMSKeyReference{
			Name:      "test-key",
			KeyRing:   "test-keyring",
			ProjectID: "test-project",
			Location:  "us-east1",
		},
	}
}

func installConfig(mode types.CredentialsMode, opts ...func(*types.InstallConfig)) *types.InstallConfig {
	ic := &types.InstallConfig{
		CredentialsMode: mode,
		Platform: types.Platform{
			GCP: &gcptypes.Platform{
				ProjectID:               "test-project",
				Region:                  "us-east1",
				FirewallRulesManagement: gcptypes.ManagedFirewallRules,
			},
		},
		ControlPlane: &types.MachinePool{Name: "master"},
		Compute:      []types.MachinePool{{Name: "worker"}},
	}
	for _, opt := range opts {
		opt(ic)
	}
	return ic
}

// withServiceAccounts sets a service account everywhere it can be set.
func withServiceAccounts(ic *types.InstallConfig) {
	withPoolServiceAccounts(ic)
	withDefaultServiceAccount(ic)
}

// withPoolServiceAccounts gives every pool its own service account and leaves the default machine
// platform unset, so the pools alone have to satisfy the requirement.
func withPoolServiceAccounts(ic *types.InstallConfig) {
	ic.ControlPlane.Platform.GCP = &gcptypes.MachinePool{ServiceAccount: "master@test-project.iam.gserviceaccount.com"}
	ic.Compute[0].Platform.GCP = &gcptypes.MachinePool{ServiceAccount: "worker@test-project.iam.gserviceaccount.com"}
}

// withDefaultServiceAccount sets a service account only on the default machine platform, which
// every pool inherits.
func withDefaultServiceAccount(ic *types.InstallConfig) {
	ic.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{ServiceAccount: "nodes@test-project.iam.gserviceaccount.com"}
}

func withExistingVPC(ic *types.InstallConfig) {
	ic.GCP.Network = "test-network"
	ic.GCP.ControlPlaneSubnet = "test-master-subnet"
	ic.GCP.ComputeSubnet = "test-worker-subnet"
}

// The installer makes the same GCP API calls regardless of the credentials mode, so the
// permissions it needs must not vary with it. What varies by mode is who consumes the
// credentials once the cluster is up, which is out of scope for this permission set.
func TestRequiredPermissionGroupsIgnoresCredentialsMode(t *testing.T) {
	modes := []types.CredentialsMode{
		"",
		types.MintCredentialsMode,
		types.PassthroughCredentialsMode,
		types.ManualCredentialsMode,
	}

	configs := map[string][]func(*types.InstallConfig){
		"defaults":                         nil,
		"pre-created service accounts":     {withServiceAccounts},
		"pool service accounts":            {withPoolServiceAccounts},
		"default platform service account": {withDefaultServiceAccount},
		"existing vpc":                     {withExistingVPC},
	}

	for name, opts := range configs {
		t.Run(name, func(t *testing.T) {
			want := RequiredPermissionGroups(installConfig(modes[0], opts...))
			for _, mode := range modes[1:] {
				got := RequiredPermissionGroups(installConfig(mode, opts...))
				assert.Equal(t, want, got, "credentials mode %q changed the required permissions", mode)
			}
		})
	}
}

func TestRequiredPermissionGroups(t *testing.T) {
	cases := []struct {
		name     string
		opts     []func(*types.InstallConfig)
		expected []PermissionGroup
		excluded []PermissionGroup
	}{{
		name: "installer creates the network, dns and the service accounts",
		expected: []PermissionGroup{
			PermissionCreateBase, PermissionDeleteBase,
			PermissionCreateDNS,
			PermissionCreateNetwork, PermissionDeleteNetwork,
			PermissionCreateFirewall, PermissionDeleteManagedFirewalls,
			PermissionCreateServiceAccounts, PermissionDeleteServiceAccounts,
		},
		excluded: []PermissionGroup{PermissionCreateKMSSupport},
	}, {
		name:     "pre-created service accounts drop the service account groups",
		opts:     []func(*types.InstallConfig){withServiceAccounts},
		excluded: []PermissionGroup{PermissionCreateServiceAccounts, PermissionDeleteServiceAccounts},
	}, {
		// Every pool carries its own service account, so there is nothing left for a default
		// machine platform to supply and its absence must not pull the groups back in.
		name:     "service accounts on the pools alone drop the service account groups",
		opts:     []func(*types.InstallConfig){withPoolServiceAccounts},
		excluded: []PermissionGroup{PermissionCreateServiceAccounts, PermissionDeleteServiceAccounts},
	}, {
		// Neither pool sets a service account, but both inherit one from the default machine
		// platform, so the installer still has none to mint.
		name:     "a default machine platform service account covers every pool",
		opts:     []func(*types.InstallConfig){withDefaultServiceAccount},
		excluded: []PermissionGroup{PermissionCreateServiceAccounts, PermissionDeleteServiceAccounts},
	}, {
		// One pool is left without a service account, so the installer has to mint that one.
		name: "a service account on only one pool keeps the service account groups",
		opts: []func(*types.InstallConfig){func(ic *types.InstallConfig) {
			ic.ControlPlane.Platform.GCP = &gcptypes.MachinePool{ServiceAccount: "master@test-project.iam.gserviceaccount.com"}
		}},
		expected: []PermissionGroup{PermissionCreateServiceAccounts, PermissionDeleteServiceAccounts},
	}, {
		name: "existing vpc drops the network and firewall groups",
		opts: []func(*types.InstallConfig){withExistingVPC},
		excluded: []PermissionGroup{
			PermissionCreateNetwork, PermissionDeleteNetwork,
			PermissionCreateFirewall, PermissionDeleteManagedFirewalls,
		},
	}, {
		name: "unmanaged firewall rules drop only the firewall groups",
		opts: []func(*types.InstallConfig){func(ic *types.InstallConfig) {
			ic.GCP.FirewallRulesManagement = gcptypes.UnmanagedFirewallRules
		}},
		expected: []PermissionGroup{PermissionCreateNetwork, PermissionDeleteNetwork},
		excluded: []PermissionGroup{PermissionCreateFirewall, PermissionDeleteManagedFirewalls},
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			groups := RequiredPermissionGroups(installConfig(types.ManualCredentialsMode, tc.opts...))
			for _, group := range tc.expected {
				assert.Contains(t, groups, group)
			}
			for _, group := range tc.excluded {
				assert.NotContains(t, groups, group)
			}
		})
	}
}

// The installer creates the private zone and the api/api-int records unless the user provides
// DNS, but destroy sweeps DNS unconditionally -- it reads metadata.json and cannot see
// userProvisionedDNS -- so the delete permissions must survive in every case.
func TestRequiredPermissionGroupsUserProvisionedDNS(t *testing.T) {
	cases := []struct {
		name            string
		dns             dnstypes.UserProvisionedDNS
		expectCreateDNS bool
	}{
		{name: "unset", dns: "", expectCreateDNS: true},
		{name: "disabled", dns: dnstypes.UserProvisionedDNSDisabled, expectCreateDNS: true},
		{name: "enabled", dns: dnstypes.UserProvisionedDNSEnabled, expectCreateDNS: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dnsMode := tc.dns
			groups := RequiredPermissionGroups(installConfig(types.ManualCredentialsMode, func(ic *types.InstallConfig) {
				ic.GCP.UserProvisionedDNS = dnsMode
			}))

			if tc.expectCreateDNS {
				assert.Contains(t, groups, PermissionCreateDNS)
			} else {
				assert.NotContains(t, groups, PermissionCreateDNS)
			}

			deletePerms, err := Permissions(PermissionDeleteBase)
			assert.NoError(t, err)
			assert.Contains(t, groups, PermissionDeleteBase)
			assert.Contains(t, deletePerms, "dns.managedZones.list", "destroy always lists DNS zones")
			assert.Contains(t, deletePerms, "dns.resourceRecordSets.list", "destroy always lists record sets")
		})
	}
}

// A KMS key can be set on the control plane pool, on a compute pool, or on the default machine
// platform, and validation checks it wherever it appears.
func TestRequiredPermissionGroupsKMSKey(t *testing.T) {
	cases := []struct {
		name     string
		opt      func(*types.InstallConfig)
		expected bool
	}{{
		name:     "no key",
		opt:      func(ic *types.InstallConfig) {},
		expected: false,
	}, {
		name: "default machine platform",
		opt: func(ic *types.InstallConfig) {
			ic.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{}
			ic.GCP.DefaultMachinePlatform.OSDisk.EncryptionKey = kmsKey()
		},
		expected: true,
	}, {
		name: "control plane only",
		opt: func(ic *types.InstallConfig) {
			ic.ControlPlane.Platform.GCP = &gcptypes.MachinePool{}
			ic.ControlPlane.Platform.GCP.OSDisk.EncryptionKey = kmsKey()
		},
		expected: true,
	}, {
		name: "compute only",
		opt: func(ic *types.InstallConfig) {
			ic.Compute[0].Platform.GCP = &gcptypes.MachinePool{}
			ic.Compute[0].Platform.GCP.OSDisk.EncryptionKey = kmsKey()
		},
		expected: true,
	}, {
		name: "encryption key without a kms key",
		opt: func(ic *types.InstallConfig) {
			ic.ControlPlane.Platform.GCP = &gcptypes.MachinePool{}
			ic.ControlPlane.Platform.GCP.OSDisk.EncryptionKey = &gcptypes.EncryptionKeyReference{}
		},
		expected: false,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			groups := RequiredPermissionGroups(installConfig(types.ManualCredentialsMode, tc.opt))
			if tc.expected {
				assert.Contains(t, groups, PermissionCreateKMSSupport)
			} else {
				assert.NotContains(t, groups, PermissionCreateKMSSupport)
			}
		})
	}
}

// A user-supplied OS image can be set on the control plane pool, on a compute pool, or on the
// default machine platform, and the installer has to read it wherever it appears.
func TestRequiredPermissionGroupsOSImage(t *testing.T) {
	osImage := func() *gcptypes.OSImage {
		return &gcptypes.OSImage{Name: "test-image", Project: "test-project"}
	}

	cases := []struct {
		name     string
		opt      func(*types.InstallConfig)
		expected bool
	}{{
		name:     "no image",
		opt:      func(ic *types.InstallConfig) {},
		expected: false,
	}, {
		name: "default machine platform",
		opt: func(ic *types.InstallConfig) {
			ic.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{OSImage: osImage()}
		},
		expected: true,
	}, {
		name: "control plane only",
		opt: func(ic *types.InstallConfig) {
			ic.ControlPlane.Platform.GCP = &gcptypes.MachinePool{OSImage: osImage()}
		},
		expected: true,
	}, {
		name: "compute only",
		opt: func(ic *types.InstallConfig) {
			ic.Compute[0].Platform.GCP = &gcptypes.MachinePool{OSImage: osImage()}
		},
		expected: true,
	}, {
		name: "image without a name",
		opt: func(ic *types.InstallConfig) {
			ic.ControlPlane.Platform.GCP = &gcptypes.MachinePool{OSImage: &gcptypes.OSImage{Project: "test-project"}}
		},
		expected: false,
	}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			groups := RequiredPermissionGroups(installConfig(types.ManualCredentialsMode, tc.opt))
			if tc.expected {
				assert.Contains(t, groups, PermissionBYOImage)
			} else {
				assert.NotContains(t, groups, PermissionBYOImage)
			}
		})
	}
}

// Every group returned by RequiredPermissionGroups must resolve to a permission list, otherwise
// writing the custom role fails at install time rather than here. Guards against a group constant
// outliving its entry in the permissions map.
func TestRequiredPermissionGroupsAreDefined(t *testing.T) {
	configs := [][]func(*types.InstallConfig){
		nil,
		{withServiceAccounts},
		{withPoolServiceAccounts},
		{withDefaultServiceAccount},
		{withExistingVPC},
		{func(ic *types.InstallConfig) {
			ic.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{
				OSImage: &gcptypes.OSImage{Name: "test-image", Project: "test-project"},
			}
		}},
		{func(ic *types.InstallConfig) {
			ic.GCP.UserProvisionedDNS = dnstypes.UserProvisionedDNSEnabled
		}},
		{func(ic *types.InstallConfig) {
			ic.GCP.DefaultMachinePlatform = &gcptypes.MachinePool{}
			ic.GCP.DefaultMachinePlatform.OSDisk.EncryptionKey = kmsKey()
		}},
	}

	for _, opts := range configs {
		for _, group := range RequiredPermissionGroups(installConfig(types.ManualCredentialsMode, opts...)) {
			perms, err := Permissions(group)
			assert.NoError(t, err, "group %q", group)
			assert.NotEmpty(t, perms, "group %q", group)
		}
	}
}

// The installer grants roles to the service accounts it creates in every credentials mode, so
// setIamPolicy has to travel with the service account group rather than with a mode.
func TestCreateServiceAccountsIncludesSetIamPolicy(t *testing.T) {
	perms, err := Permissions(PermissionCreateServiceAccounts)
	assert.NoError(t, err)
	assert.Contains(t, perms, "resourcemanager.projects.setIamPolicy")
	assert.Contains(t, perms, "iam.serviceAccounts.actAs")
}

// The installer neither creates nor deletes KMS resources, so this group must stay read-only.
func TestKMSSupportIsReadOnly(t *testing.T) {
	perms, err := Permissions(PermissionCreateKMSSupport)
	assert.NoError(t, err)
	for _, perm := range perms {
		assert.NotContains(t, perm, ".create", "the installer does not create KMS resources")
		assert.NotContains(t, perm, ".delete", "GCP does not permit deleting key rings or crypto keys")
		assert.NotContains(t, perm, ".destroy", "the installer does not destroy key versions")
	}
}
