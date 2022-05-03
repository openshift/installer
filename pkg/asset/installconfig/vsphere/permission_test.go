package vsphere

import (
	"context"
	"testing"

	"github.com/openshift/installer/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	vim25types "github.com/vmware/govmomi/vim25/types"

	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
)

func buildPermissionGroup(authManagerMock *mock.MockAuthManager,
	managedObjectRef vim25types.ManagedObjectReference,
	username string,
	group PermissionGroupDefinition,
	groupName permissionGroup,
	overrideGroup *permissionGroup) {
	permissionToApply := group.Permissions
	if overrideGroup != nil && *overrideGroup == groupName {
		permissionToApply = permissionToApply[:len(permissionToApply)-1]
	}
	authManagerMock.EXPECT().FetchUserPrivilegeOnEntities(gomock.Any(), []vim25types.ManagedObjectReference{
		managedObjectRef,
	}, username).Return([]vim25types.UserPrivilegeResult{
		{
			Privileges: permissionToApply,
		},
	}, nil).AnyTimes()
}

func buildAuthManagerClient(ctx context.Context, mockCtrl *gomock.Controller, finder Finder, username string, overrideGroup *permissionGroup) (*mock.MockAuthManager, error) {
	authManagerClient := mock.NewMockAuthManager(mockCtrl)
	for groupName, group := range permissions {
		switch groupName {
		case permissionVcenter:
			vcenter, err := finder.Folder(ctx, "/")
			if err != nil {
				return nil, err
			}
			buildPermissionGroup(authManagerClient, vcenter.Reference(), username, group, groupName, overrideGroup)
		case permissionDatacenter:
			datacenters, err := finder.DatacenterList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datacenter := range datacenters {
				buildPermissionGroup(authManagerClient, datacenter.Reference(), username, group, groupName, overrideGroup)
			}
		case permissionDatastore:
			datastores, err := finder.DatastoreList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datastore := range datastores {
				buildPermissionGroup(authManagerClient, datastore.Reference(), username, group, groupName, overrideGroup)
			}
		case permissionCluster:
			clusters, err := finder.ClusterComputeResourceList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, cluster := range clusters {
				buildPermissionGroup(authManagerClient, cluster.Reference(), username, group, groupName, overrideGroup)
			}
		case permissionPortgroup:
			networks, err := finder.NetworkList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, network := range networks {
				buildPermissionGroup(authManagerClient, network.Reference(), username, group, groupName, overrideGroup)
			}
		case permissionResourcePool:
			resourcePool, err := finder.ResourcePool(ctx, "/DC0/host/DC0_C0/Resources/test-resourcepool")
			if err != nil {
				return nil, err
			}
			buildPermissionGroup(authManagerClient, resourcePool.Reference(), username, group, groupName, overrideGroup)
		case permissionFolder:
			var folders = []string{"/DC0/vm", "/DC0/vm/my-folder"}
			for _, folder := range folders {
				folder, err := finder.Folder(ctx, folder)
				if err != nil {
					return nil, err
				}
				buildPermissionGroup(authManagerClient, folder.Reference(), username, group, groupName, overrideGroup)
			}
		}
	}
	return authManagerClient, nil
}

func TestPermissionValidate(t *testing.T) {
	ctx := context.TODO()
	server := mock.StartSimulator()
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

	rootFolder := object.NewRootFolder(client)
	_, err = rootFolder.CreateFolder(ctx, "/DC0/vm/my-folder")

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

	validInstallConfig := validIPIInstallConfig("DC0", "")
	userDefinedFolderInstallConfig := validIPIInstallConfig("DC0", "")
	userDefinedFolderInstallConfig.VSphere.Folder = "/DC0/vm/my-folder"

	invalidDatacenterInstallConfig := validIPIInstallConfig("DC0", "")
	invalidDatacenterInstallConfig.VSphere.Datacenter = "invalid"

	username := validInstallConfig.VSphere.Username

	validPermissionsAuthManagerClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, nil)
	if err != nil {
		t.Error(err)
		return
	}
	missingPortgroupPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionPortgroup)
	if err != nil {
		t.Error(err)
		return
	}

	missingVCenterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionVcenter)
	if err != nil {
		t.Error(err)
		return
	}

	missingClusterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionCluster)
	if err != nil {
		t.Error(err)
		return
	}

	missingDatastorePermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionDatastore)
	if err != nil {
		t.Error(err)
		return
	}

	missingDatacenterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionDatacenter)
	if err != nil {
		t.Error(err)
		return
	}

	missingFolderPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionFolder)
	if err != nil {
		t.Error(err)
		return
	}

	missingResourcePoolPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionResourcePool)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name             string
		installConfig    *types.InstallConfig
		validationMethod func(*validationContext, *types.InstallConfig) error
		expectErr        string
		authManager      AuthManager
	}{
		{
			name:             "valid Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      validPermissionsAuthManagerClient,
		},
		{
			name:             "missing portgroup Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingPortgroupPermissionsClient,
			expectErr:        "privileges missing for vSphere Port Group: Network.Assign",
		},
		{
			name:             "missing vCenter Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingVCenterPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter: StorageProfile.View",
		},
		{
			name:             "missing cluster Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingClusterPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Cluster: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:             "missing datacenter Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingDatacenterPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Datacenter: Folder.Delete",
		},
		{
			name:             "missing datastore Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingDatastorePermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Datastore: InventoryService.Tagging.ObjectAttachable",
		},
		{
			name:             "missing resource pool Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingResourcePoolPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Resource Pool: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:             "missing user-defined folder Permissions but no folder defined",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingFolderPermissionsClient,
		},
		{
			name:             "missing user-defined folder Permissions",
			installConfig:    userDefinedFolderInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingFolderPermissionsClient,
			expectErr:        "privileges missing for Pre-existing Virtual Machine Folder: VirtualMachine.Provisioning.DeployTemplate",
		},
		{
			name:             "invalid defined datacenter",
			installConfig:    invalidDatacenterInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      validPermissionsAuthManagerClient,
			expectErr:        "platform.vsphere.datacenter: Invalid value: \"invalid\": datacenter 'invalid' not found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validationCtx := &validationContext{
				User:        username,
				AuthManager: test.authManager,
				Finder:      finder,
				Client:      client,
			}
			err := test.validationMethod(validationCtx, test.installConfig)
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
