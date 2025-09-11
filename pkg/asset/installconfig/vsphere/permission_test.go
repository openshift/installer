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
				buildPermissionGroup(t, authManagerClient, cluster.Reference(), username, group, groupName, overrideGroup, permissionsToRemoveFromResource)
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

func TestPermissionValidate(t *testing.T) {
	ctx := context.TODO()
	vs := mock.NewSimulator("", "", 0, 0)
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
