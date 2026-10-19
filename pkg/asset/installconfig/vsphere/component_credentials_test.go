package vsphere

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"go.uber.org/mock/gomock"

	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func componentScopedInstallConfig(vcenters ...vsphere.VCenter) *types.InstallConfig {
	return &types.InstallConfig{
		Platform: types.Platform{
			VSphere: &vsphere.Platform{
				VCenters: vcenters,
			},
		},
	}
}

func fullComponentCredentials() *vsphere.ComponentCredentials {
	return &vsphere.ComponentCredentials{
		MachineManagement:      vsphere.Credential{User: "machine", Password: "machine-pass"},
		Storage:                vsphere.Credential{User: "csi", Password: "csi-pass"},
		CloudControllerManager: vsphere.Credential{User: "ccm", Password: "ccm-pass"},
		VSphereProblemDetector: vsphere.Credential{User: "diag", Password: "diag-pass"},
	}
}

func TestGetComponentVCenterClientErrors(t *testing.T) {
	tests := []struct {
		name          string
		installConfig *types.InstallConfig
		failureDomain vsphere.FailureDomain
		component     credentialComponent
		expectErr     string
	}{
		{
			name: "vcenter not defined",
			installConfig: componentScopedInstallConfig(vsphere.VCenter{
				Server: "test-vcenter", Datacenters: []string{"DC0"},
				ComponentCredentials: fullComponentCredentials(),
			}),
			failureDomain: vsphere.FailureDomain{Server: "missing-vcenter"},
			component:     componentStorage,
			expectErr:     "vcenter missing-vcenter not defined in vcenters",
		},
		{
			name: "missing component credentials",
			installConfig: componentScopedInstallConfig(vsphere.VCenter{
				Server: "test-vcenter", Datacenters: []string{"DC0"},
			}),
			failureDomain: vsphere.FailureDomain{Server: "test-vcenter"},
			component:     componentStorage,
			expectErr:     "vcenter test-vcenter is missing componentCredentials",
		},
		{
			name: "unsupported component",
			installConfig: componentScopedInstallConfig(vsphere.VCenter{
				Server: "test-vcenter", Datacenters: []string{"DC0"},
				ComponentCredentials: fullComponentCredentials(),
			}),
			failureDomain: vsphere.FailureDomain{Server: "test-vcenter"},
			component:     componentMachineManagement,
			expectErr:     `unsupported component credential "machineManagement"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := getComponentVCenterClient(tt.failureDomain, tt.installConfig, tt.component)
			assert.ErrorContains(t, err, tt.expectErr)
		})
	}
}

// TestGetComponentVCenterClientLoginErrorNamesComponent verifies that a
// component connection failure identifies the component and vCenter, instead
// of surfacing an ambiguous vCenter login error.
func TestGetComponentVCenterClientLoginErrorNamesComponent(t *testing.T) {
	// Port 1 on localhost refuses connections deterministically.
	const unreachableServer = "https://127.0.0.1:1"
	installConfig := componentScopedInstallConfig(vsphere.VCenter{
		Server:               unreachableServer,
		Datacenters:          []string{"DC0"},
		ComponentCredentials: fullComponentCredentials(),
	})
	failureDomain := vsphere.FailureDomain{Server: unreachableServer}

	_, _, err := getComponentVCenterClient(failureDomain, installConfig, componentStorage)
	if err == nil {
		t.Fatal("expected a connection error")
	}
	assert.ErrorContains(t, err, "storage")
	assert.ErrorContains(t, err, unreachableServer)
}

func TestValidateComponentPrivileges(t *testing.T) {
	ctx := context.TODO()
	vs := mock.NewSimulator("8.0.2", "8.0.2", 24321653, 23825572)
	server, err := vs.StartSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	finder, err := mock.GetFinder(server)
	if err != nil {
		t.Fatal(err)
	}

	rootFolder, err := finder.Folder(ctx, "/")
	if err != nil {
		t.Fatal(err)
	}
	datastore, err := finder.Datastore(ctx, "/DC0/datastore/LocalDS_0")
	if err != nil {
		t.Fatal(err)
	}
	cluster, err := finder.ClusterComputeResource(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Fatal(err)
	}
	folder, err := finder.Folder(ctx, "/DC0/vm")
	if err != nil {
		t.Fatal(err)
	}

	// The storage component privilege matrix: vCenter root, datastore,
	// compute cluster and VM folder are checked; the datacenter and resource
	// pool groups are absent from the matrix and must be skipped.
	storagePrivileges := map[vim25types.ManagedObjectReference][]string{
		rootFolder.Reference(): {"Cns.Searchable", "StorageProfile.View"},
		datastore.Reference():  {"Datastore.FileManagement", "System.Read"},
		cluster.Reference():    {"Host.Config.Storage"},
		folder.Reference():     {"VirtualMachine.Config.AddExistingDisk", "VirtualMachine.Config.AddRemoveDevice"},
	}

	failureDomain := validMultiVCenterPlatform().FailureDomains[0]
	// A custom resource pool is not part of the storage matrix; drop it so the
	// lookup is not exercised against a pool the simulator does not create.
	failureDomain.Topology.ResourcePool = ""

	tests := []struct {
		name              string
		withheldPrivilege string
		expectErr         []string
	}{
		{
			name: "all component privileges present",
		},
		{
			name:              "missing vCenter privilege",
			withheldPrivilege: "Cns.Searchable",
			expectErr:         []string{"storage credentials", "Cns.Searchable"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			authManager := mock.NewMockAuthManager(mockCtrl)
			for ref, privileges := range storagePrivileges {
				granted := make([]string, 0, len(privileges))
				for _, privilege := range privileges {
					if privilege != tt.withheldPrivilege {
						granted = append(granted, privilege)
					}
				}
				authManager.EXPECT().FetchUserPrivilegeOnEntities(gomock.Any(), []vim25types.ManagedObjectReference{ref}, "test-user").
					Return([]vim25types.UserPrivilegeResult{{Privileges: granted}}, nil).AnyTimes()
			}

			validationCtx := &validationContext{
				AuthManager: authManager,
				Finder:      finder,
				Username:    "test-user",
				Component:   componentStorage,
			}

			errs := validateComponentPrivileges(validationCtx, &failureDomain)
			if len(tt.expectErr) == 0 {
				assert.Empty(t, errs)
			} else {
				aggregate := errs.ToAggregate()
				for _, expect := range tt.expectErr {
					assert.ErrorContains(t, aggregate, expect)
				}
			}
		})
	}
}
