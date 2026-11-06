package vsphere

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25/mo"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"go.uber.org/mock/gomock"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	v1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func buildPermissionGroup(t *testing.T, authManagerMock *mock.MockAuthManager,
	managedObjectRef vim25types.ManagedObjectReference,
	username string,
	group PermissionGroupDefinition,
	groupName permissionGroup,
	overrideGroup *permissionGroup,
	permissionToExcludeSet sets.String) {
	t.Helper()
	permissionsToApply := group.Permissions

	if overrideGroup != nil && *overrideGroup == groupName {
		filteredPermissionsToApply := sets.NewString(group.Permissions...)
		filteredPermissions := filteredPermissionsToApply.Difference(permissionToExcludeSet)
		permissionsToApply = filteredPermissions.List()
	}
	authManagerMock.EXPECT().FetchUserPrivilegeOnEntities(gomock.Any(), []vim25types.ManagedObjectReference{
		managedObjectRef,
	}, username).Return([]vim25types.UserPrivilegeResult{
		{
			Privileges: permissionsToApply,
		},
	}, nil).AnyTimes()
}

func buildAuthManagerClient(ctx context.Context,
	t *testing.T,
	mockCtrl *gomock.Controller,
	finder Finder,
	username string,
	overrideGroup *permissionGroup,
	permissionsToRemoveFromResource sets.String,
	permissionsToRemoveFromAvailable sets.String) (*mock.MockAuthManager, error) {
	t.Helper()
	authManagerClient := mock.NewMockAuthManager(mockCtrl)
	authManagerMo := vim25types.ManagedObjectReference{
		Type:  "auth-manager",
		Value: "auth-manager",
	}
	authManagerClient.EXPECT().Reference().Return(authManagerMo).AnyTimes()

	privilegeListMap := map[string]vim25types.AuthorizationPrivilege{}

	for groupName, group := range permissions {
		availablePrivileges := sets.NewString(group.Permissions...)
		availablePrivileges = availablePrivileges.Difference(permissionsToRemoveFromAvailable)
		for _, availablePrivilege := range availablePrivileges.List() {
			privilegeListMap[availablePrivilege] = vim25types.AuthorizationPrivilege{
				PrivId: availablePrivilege,
			}
		}
		switch groupName {
		case permissionVcenter:
			vcenter, err := finder.Folder(ctx, "/")
			if err != nil {
				return nil, err
			}
			buildPermissionGroup(t, authManagerClient, vcenter.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
		case permissionDatacenter:
			datacenters, err := finder.DatacenterList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datacenter := range datacenters {
				buildPermissionGroup(t, authManagerClient, datacenter.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
			}
		case permissionDatastore:
			datastores, err := finder.DatastoreList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datastore := range datastores {
				buildPermissionGroup(t, authManagerClient, datastore.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
			}
		case permissionCluster:
			clusters, err := finder.ClusterComputeResourceList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, cluster := range clusters {
				// The mock account is sufficiently privileged for HostGroup
				// failure domains, which require the cluster reconfiguration
				// privilege on the compute cluster.
				clusterGroup := PermissionGroupDefinition{
					Permissions: append(append([]string(nil), group.Permissions...), "Host.Inventory.EditCluster"),
					Description: group.Description,
				}
				buildPermissionGroup(t, authManagerClient, cluster.Reference(), username, clusterGroup, groupName, overrideGroup, permissionsToRemoveFromResource)
			}
		case permissionPortgroup:
			networks, err := finder.NetworkList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, network := range networks {
				buildPermissionGroup(t, authManagerClient, network.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
			}
		case permissionResourcePool:
			resourcePools := []string{"/DC0/host/DC0_C0/Resources/test-resourcepool", "/DC0/host/DC0_C0/Resources"}
			for _, resourcePoolPath := range resourcePools {
				resourcePool, err := finder.ResourcePool(ctx, resourcePoolPath)
				if err != nil {
					return nil, err
				}
				buildPermissionGroup(t, authManagerClient, resourcePool.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
			}
		case permissionFolder:
			var folders = []string{"/DC0/vm", "/DC0/vm/my-folder"}
			for _, folder := range folders {
				folder, err := finder.Folder(ctx, folder)
				if err != nil {
					return nil, err
				}
				buildPermissionGroup(t, authManagerClient, folder.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
			}
		}
	}
	privilegeListSlice := make([]vim25types.AuthorizationPrivilege, 0, len(privilegeListMap))
	for _, authorizationPrivilege := range privilegeListMap {
		privilegeListSlice = append(privilegeListSlice, authorizationPrivilege)
	}

	authorizationManager := mo.AuthorizationManager{
		PrivilegeList: privilegeListSlice,
	}

	authManagerClient.EXPECT().Properties(gomock.Any(), authManagerMo, []string{"privilegeList"}, gomock.Any()).SetArg(3, authorizationManager).AnyTimes()
	return authManagerClient, nil
}

func validIPIMultiZoneInstallConfig() *types.InstallConfig {
	installConfig := validIPIInstallConfig()
	validMultiZonePlatform := validMultiVCenterPlatform()
	installConfig.VSphere.VCenters = validMultiZonePlatform.VCenters
	installConfig.VSphere.FailureDomains = validMultiZonePlatform.FailureDomains

	return installConfig
}

func TestComponentPermissionMatrix(t *testing.T) {
	tests := []struct {
		component      credentialComponent
		group          permissionGroup
		mustContain    []string
		mustNotContain []string
	}{
		{
			component:      componentStorage,
			group:          permissionDatastore,
			mustContain:    []string{"Datastore.FileManagement", "System.Read"},
			mustNotContain: []string{"VirtualMachine.Provisioning.Clone"},
		},
		{
			component:      componentCloudControllerManager,
			group:          permissionVcenter,
			mustContain:    []string{"System.Read"},
			mustNotContain: []string{"InventoryService.Tagging.CreateTag", "VirtualMachine.Inventory.Delete"},
		},
		{
			component:      componentVSphereProblemDetector,
			group:          permissionDatastore,
			mustContain:    []string{"Datastore.Browse"},
			mustNotContain: []string{"Datastore.FileManagement", "VirtualMachine.Inventory.Delete"},
		},
	}

	for _, test := range tests {
		t.Run(string(test.component)+"/"+string(test.group), func(t *testing.T) {
			group, found := componentPermissionGroup(test.component, test.group)
			if !found {
				t.Fatalf("component permission group not found")
			}
			for _, privilege := range test.mustContain {
				assert.Contains(t, group.Permissions, privilege)
			}
			for _, privilege := range test.mustNotContain {
				assert.NotContains(t, group.Permissions, privilege)
			}
		})
	}
}

func TestMachineManagementPermissions(t *testing.T) {
	legacyCluster := permissions[permissionCluster]

	withoutHostGroups := machineManagementPermissions(false)
	assert.NotContains(t, withoutHostGroups[permissionCluster].Permissions, "Host.Inventory.EditCluster")
	assert.ElementsMatch(t, withoutHostGroups[permissionCluster].Permissions, legacyCluster.Permissions)

	withHostGroups := machineManagementPermissions(true)
	assert.Contains(t, withHostGroups[permissionCluster].Permissions, "Host.Inventory.EditCluster")
	for group, definition := range permissions {
		if group == permissionCluster {
			continue
		}
		assert.ElementsMatch(t, withHostGroups[group].Permissions, definition.Permissions, string(group))
	}
}

func TestHostGroupTopologyUsed(t *testing.T) {
	tests := []struct {
		name           string
		failureDomains []vsphere.FailureDomain
		server         string
		expected       bool
	}{
		{
			name: "host group failure domain on the server",
			failureDomains: []vsphere.FailureDomain{
				{Server: "test-vcenter", ZoneType: vsphere.HostGroupFailureDomain},
			},
			server:   "test-vcenter",
			expected: true,
		},
		{
			name: "host group failure domain on another server",
			failureDomains: []vsphere.FailureDomain{
				{Server: "other-vcenter", ZoneType: vsphere.HostGroupFailureDomain},
			},
			server:   "test-vcenter",
			expected: false,
		},
		{
			name: "no zone type on the server",
			failureDomains: []vsphere.FailureDomain{
				{Server: "test-vcenter"},
			},
			server:   "test-vcenter",
			expected: false,
		},
		{
			name:           "no failure domains",
			failureDomains: nil,
			server:         "test-vcenter",
			expected:       false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, hostGroupTopologyUsed(test.failureDomains, test.server))
		})
	}
}

// TestMachineManagementValidation verifies that the machine-management
// privilege matrix requires Host.Inventory.EditCluster for HostGroup failure
// domains only when the component-scoped credentials feature gate is enabled;
// with the gate disabled the matrix matches the legacy behavior that existing
// global-credential install configs are validated against.
func TestMachineManagementValidation(t *testing.T) {
	tests := []struct {
		name                     string
		featureSet               v1.FeatureSet
		failureDomains           []vsphere.FailureDomain
		wantEditCluster          bool
		wantReconfigurationCheck bool
	}{
		{
			name:       "HostGroup topology with the gate disabled keeps the legacy matrix",
			featureSet: v1.Default,
			failureDomains: []vsphere.FailureDomain{
				{Server: "test-vcenter", ZoneType: vsphere.HostGroupFailureDomain},
			},
		},
		{
			name:       "HostGroup topology with the gate enabled requires EditCluster",
			featureSet: v1.DevPreviewNoUpgrade,
			failureDomains: []vsphere.FailureDomain{
				{Server: "test-vcenter", ZoneType: vsphere.HostGroupFailureDomain},
			},
			wantEditCluster:          true,
			wantReconfigurationCheck: true,
		},
		{
			name:       "compute cluster topology with the gate enabled keeps the legacy matrix",
			featureSet: v1.DevPreviewNoUpgrade,
			failureDomains: []vsphere.FailureDomain{
				{Server: "test-vcenter", ZoneType: vsphere.ComputeClusterFailureDomain},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ic := &types.InstallConfig{
				FeatureSet: test.featureSet,
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						FailureDomains: test.failureDomains,
					},
				},
			}
			componentGroups, checkReconfiguration := machineManagementValidation(ic, "test-vcenter")
			if test.wantEditCluster {
				assert.Contains(t, componentGroups[permissionCluster].Permissions, "Host.Inventory.EditCluster")
			} else {
				assert.NotContains(t, componentGroups[permissionCluster].Permissions, "Host.Inventory.EditCluster")
			}
			assert.Equal(t, test.wantReconfigurationCheck, checkReconfiguration)
		})
	}
}

func TestMachineManagementClusterPrivileges(t *testing.T) {
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

	cluster, err := finder.ClusterComputeResource(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Fatal(err)
	}
	clusterRef := cluster.Reference()

	legacyClusterPrivileges := append([]string(nil), permissions[permissionCluster].Permissions...)

	tests := []struct {
		name       string
		hostGroups bool
		privileges []string
		expectErr  string
	}{
		{
			name:       "host group topology without EditCluster fails",
			hostGroups: true,
			privileges: legacyClusterPrivileges,
			expectErr:  "Host.Inventory.EditCluster",
		},
		{
			name:       "host group topology with EditCluster passes",
			hostGroups: true,
			privileges: append(append([]string(nil), legacyClusterPrivileges...), "Host.Inventory.EditCluster"),
		},
		{
			name:       "compute cluster topology does not require EditCluster",
			hostGroups: false,
			privileges: legacyClusterPrivileges,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			authManager := mock.NewMockAuthManager(mockCtrl)
			authManager.EXPECT().FetchUserPrivilegeOnEntities(gomock.Any(), []vim25types.ManagedObjectReference{clusterRef}, "test-user").
				Return([]vim25types.UserPrivilegeResult{{Privileges: test.privileges}}, nil).AnyTimes()

			validationCtx := &validationContext{
				AuthManager:     authManager,
				Finder:          finder,
				Username:        "test-user",
				componentGroups: machineManagementPermissions(test.hostGroups),
			}

			errs := computeClusterExists(validationCtx, "/DC0/host/DC0_C0", field.NewPath("test").Child("computeCluster"), true, false)
			if test.expectErr == "" {
				assert.NoError(t, errs.ToAggregate())
			} else {
				assert.ErrorContains(t, errs.ToAggregate(), test.expectErr)
			}
		})
	}
}

// TestValidateFailureDomainHostGroupResourcePool verifies that a HostGroup
// failure domain using a custom resource pool still checks the cluster
// reconfiguration privilege when the component-scoped credentials feature
// gate is enabled: placement privileges move to the pool, but VM groups are
// created by reconfiguring the compute cluster. With the gate disabled the
// check must not run, keeping validation of existing global-credential
// install configs unchanged.
func TestValidateFailureDomainHostGroupResourcePool(t *testing.T) {
	ctx := context.TODO()
	vs := mock.NewSimulator("8.0.2", "8.0.2", 24321653, 23825572)
	server, err := vs.StartSimulator()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Close()

	client, _, err := mock.GetClient(server)
	if err != nil {
		t.Fatal(err)
	}

	finder, err := mock.GetFinder(server)
	if err != nil {
		t.Fatal(err)
	}

	resourcePools, err := finder.ResourcePoolList(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := resourcePools[0].Create(ctx, "test-resourcepool", vim25types.DefaultResourceConfigSpec()); err != nil {
		t.Fatal(err)
	}

	vmFolder, err := finder.Folder(ctx, "/DC0/vm")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := vmFolder.CreateFolder(ctx, "my-folder"); err != nil {
		t.Fatal(err)
	}

	sessionMgr := session.NewManager(client)
	userSession, err := sessionMgr.UserSession(ctx)
	if err != nil {
		t.Fatal(err)
	}

	failureDomain := validMultiVCenterPlatform().FailureDomains[0]
	failureDomain.ZoneType = vsphere.HostGroupFailureDomain
	failureDomain.Topology.HostGroup = "host-group-a"

	tests := []struct {
		name        string
		featureSet  v1.FeatureSet
		gateEnabled bool
	}{
		{
			name:        "feature gate enabled",
			featureSet:  v1.DevPreviewNoUpgrade,
			gateEnabled: true,
		},
		{
			name:       "feature gate disabled",
			featureSet: v1.Default,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			// The mock account has every privilege the installer checks except
			// Host.Inventory.EditCluster, which HostGroup failure domains with a
			// custom resource pool must still require on the compute cluster
			// when the feature gate is enabled.
			authManager, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, userSession.UserName, &permissionCluster, sets.NewString("Host.Inventory.EditCluster"), nil)
			if err != nil {
				t.Fatal(err)
			}

			ic := &types.InstallConfig{
				FeatureSet: test.featureSet,
				Platform: types.Platform{
					VSphere: &vsphere.Platform{
						FailureDomains: []vsphere.FailureDomain{
							{Server: failureDomain.Server, ZoneType: vsphere.HostGroupFailureDomain},
						},
					},
				},
			}
			componentGroups, checkReconfiguration := machineManagementValidation(ic, failureDomain.Server)

			validationCtx := &validationContext{
				AuthManager:                 authManager,
				Finder:                      finder,
				Client:                      client,
				Username:                    userSession.UserName,
				Component:                   componentMachineManagement,
				componentGroups:             componentGroups,
				checkVMGroupReconfiguration: checkReconfiguration,
			}

			// The host group does not exist in the simulator, which adds an
			// unrelated NotFound error; the assertion is whether the cluster
			// reconfiguration privilege is checked despite the custom resource
			// pool.
			errs := validateFailureDomain(validationCtx, &failureDomain, false).ToAggregate()
			if test.gateEnabled {
				assert.ErrorContains(t, errs, "Host.Inventory.EditCluster")
			} else {
				assert.NotContains(t, errs.Error(), "Host.Inventory.EditCluster")
			}
		})
	}
}

// TestValidateClusterReconfigurationPrivileges verifies the targeted
// Host.Inventory.EditCluster check used for HostGroup failure domains with a
// custom resource pool.
func TestValidateClusterReconfigurationPrivileges(t *testing.T) {
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

	cluster, err := finder.ClusterComputeResource(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Fatal(err)
	}
	clusterRef := cluster.Reference()

	tests := []struct {
		name       string
		privileges []string
		expectErr  string
	}{
		{
			name:      "missing EditCluster fails",
			expectErr: "Host.Inventory.EditCluster",
		},
		{
			name:       "EditCluster granted passes",
			privileges: []string{"Host.Inventory.EditCluster"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			authManager := mock.NewMockAuthManager(mockCtrl)
			authManager.EXPECT().FetchUserPrivilegeOnEntities(gomock.Any(), []vim25types.ManagedObjectReference{clusterRef}, "test-user").
				Return([]vim25types.UserPrivilegeResult{{Privileges: tt.privileges}}, nil).AnyTimes()

			validationCtx := &validationContext{
				AuthManager: authManager,
				Finder:      finder,
				Username:    "test-user",
			}

			errs := validateClusterReconfigurationPrivileges(validationCtx, "/DC0/host/DC0_C0", field.NewPath("test").Child("computeCluster"))
			if tt.expectErr == "" {
				assert.NoError(t, errs.ToAggregate())
			} else {
				assert.ErrorContains(t, errs.ToAggregate(), tt.expectErr)
			}
		})
	}
}

func TestPermissionValidate(t *testing.T) {
	ctx := context.TODO()
	vs := mock.NewSimulator("8.0.2", "8.0.2", 24321653, 23825572)
	server, err := vs.StartSimulator()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Close()

	client, _, err := mock.GetClient(server)
	if err != nil {
		t.Error(err)
		return
	}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	finder, err := mock.GetFinder(server)
	if err != nil {
		t.Error(err)
		return
	}

	vmFolder, err := finder.Folder(ctx, "/DC0/vm")
	if err != nil {
		t.Error(err)
		return
	}

	_, err = vmFolder.CreateFolder(ctx, "my-folder")
	if err != nil {
		t.Error(err)
		return
	}

	resourcePools, err := finder.ResourcePoolList(ctx, "/DC0/host/DC0_C0")
	if err != nil {
		t.Error(err)
		return
	}
	_, err = resourcePools[0].Create(ctx, "test-resourcepool", vim25types.DefaultResourceConfigSpec())
	if err != nil {
		t.Error(err)
		return
	}

	validMultiZoneInstallConfig := validIPIMultiZoneInstallConfig()

	sessionMgr := session.NewManager(client)
	userSession, err := sessionMgr.UserSession(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	username := userSession.UserName

	validPermissionsAuthManagerClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, nil, nil, nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingObjectAttachableDatacenter, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionDatacenter, sets.NewString("InventoryService.Tagging.ObjectAttachable"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingPortgroupPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionPortgroup, sets.NewString("Network.Assign"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingVCenterPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionVcenter, sets.NewString("StorageProfile.View"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingClusterPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionCluster, sets.NewString("VirtualMachine.Config.AddNewDisk"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	readOnlyClusterPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionCluster, sets.NewString(permissions[permissionCluster].Permissions...), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingDatastorePermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionDatastore, sets.NewString("InventoryService.Tagging.ObjectAttachable"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingDatacenterPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionDatacenter, sets.NewString("Folder.Delete"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	readOnlyDatacenterPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionDatacenter, sets.NewString(permissions[permissionDatacenter].Permissions...), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingFolderPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionFolder, sets.NewString("VirtualMachine.Provisioning.DeployTemplate"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	missingResourcePoolPermissionsClient, err := buildAuthManagerClient(ctx, t, mockCtrl, finder, username, &permissionResourcePool, sets.NewString("VirtualMachine.Config.AddNewDisk"), nil)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(*validationContext, *vsphere.FailureDomain, bool) field.ErrorList
		failureDomain    *vsphere.FailureDomain
		expectErr        string
		authManager      AuthManager
	}{
		{
			name:             "multi-zone valid Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      validPermissionsAuthManagerClient,
		},
		{
			name:             "multi-zone missing portgroup Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      missingPortgroupPermissionsClient,
			expectErr:        "privileges missing for vSphere Port Group: Network.Assign",
		},
		{
			name:             "multi-zone missing vCenter Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      missingVCenterPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter: StorageProfile.View",
		},
		{
			name:             "multi-zone missing cluster Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				failureDomain.Topology.ResourcePool = ""
				return &failureDomain
			}(),
			authManager: missingClusterPermissionsClient,
			expectErr:   "privileges missing for vSphere vCenter Cluster: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:             "multi-zone resource pool provided, compute cluster can have read-only",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      readOnlyClusterPermissionsClient,
		},
		{
			name:             "multi-zone missing datacenter Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				failureDomain.Topology.Folder = ""
				return &failureDomain
			}(),
			authManager: missingDatacenterPermissionsClient,
			expectErr:   "privileges missing for vSphere vCenter Datacenter: Folder.Delete",
		},
		{
			name:             "multi-zone user-defined folder provided, datacenter can have read-only",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      readOnlyDatacenterPermissionsClient,
		},
		{
			name:             "multi-zone missing datastore Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      missingDatastorePermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Datastore: InventoryService.Tagging.ObjectAttachable",
		},
		{
			name:             "multi-zone missing resource pool Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      missingResourcePoolPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Resource Pool: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:             "multi-zone missing user-defined folder Permissions but no folder defined",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				failureDomain.Topology.Folder = ""
				failureDomain.Topology.ResourcePool = ""
				return &failureDomain
			}(),
			authManager: missingFolderPermissionsClient,
		},
		{
			name:             "multi-zone missing user-defined folder Permissions",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain:    &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:      missingFolderPermissionsClient,
			expectErr:        "privileges missing for Pre-existing Virtual Machine Folder: VirtualMachine.Provisioning.DeployTemplate",
		},
		{
			name:             "missing datacenter permission InventoryService.Tagging.ObjectAttachable",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				// If folder is empty permissions are checked at the folder level
				failureDomain.Topology.Folder = ""
				return &failureDomain
			}(),

			authManager: missingObjectAttachableDatacenter,
			expectErr:   "privileges missing for vSphere vCenter Datacenter: InventoryService.Tagging.ObjectAttachable",
		},
		{
			name:             "invalid defined datacenter",
			installConfig:    validMultiZoneInstallConfig,
			validationMethod: validateFailureDomain,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				// If folder is empty permissions are checked at the folder level
				failureDomain.Topology.Datacenter = "invalid"
				return &failureDomain
			}(),
			authManager: validPermissionsAuthManagerClient,
			expectErr:   `platform.vsphere.failureDomains.topology.datacenter: Invalid value: "invalid": datacenter 'invalid' not found`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationCtx := &validationContext{
				AuthManager: test.authManager,
				Finder:      finder,
				Client:      client,
			}
			pushPrivileges()
			err := pruneToAvailablePermissions(ctx, test.authManager)
			if err != nil {
				assert.NoError(t, err)
			}
			if test.validationMethod != nil {
				err = test.validationMethod(validationCtx, test.failureDomain, false).ToAggregate()
			} else {
				err = errors.New("no test method defined")
			}
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
			popPrivileges()
		})
	}
}

var permissionsBackup = map[permissionGroup]PermissionGroupDefinition{}

func pushPrivileges() {
	for permissionGroupKey, permissionGroup := range permissions {
		var permissions []string
		permissions = append(permissions, permissionGroup.Permissions...)

		permissionsBackup[permissionGroupKey] = PermissionGroupDefinition{
			Permissions: permissions,
			Description: permissionGroup.Description,
		}
	}
}

func popPrivileges() {
	for permissionGroupKey, permissionGroup := range permissionsBackup {
		var permissionsSet []string

		permissionsSet = append(permissionsSet, permissionGroup.Permissions...)
		permissions[permissionGroupKey] = PermissionGroupDefinition{
			Permissions: permissionsSet,
			Description: permissionGroup.Description,
		}
	}
}
