package vsphere

import (
	"context"
	"errors"
	"testing"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/govmomi/object"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/golang/mock/gomock"
	"github.com/openshift/installer/pkg/asset/installconfig/vsphere/mock"
)

func buildPermissionGroup(authManagerMock *mock.MockAuthManager,
	managedObjectRef vim25types.ManagedObjectReference,
	username string,
	group PermissionGroupDefinition,
	groupName permissionGroup,
	overrideGroup *permissionGroup,
	makeReadOnly bool) {
	permissionToApply := group.Permissions
	if overrideGroup != nil && *overrideGroup == groupName {
		if makeReadOnly {
			permissionToApply = []string{}
		} else {
			permissionToApply = permissionToApply[:len(permissionToApply)-1]
		}

	}
	authManagerMock.EXPECT().FetchUserPrivilegeOnEntities(gomock.Any(), []vim25types.ManagedObjectReference{
		managedObjectRef,
	}, username).Return([]vim25types.UserPrivilegeResult{
		{
			Privileges: permissionToApply,
		},
	}, nil).AnyTimes()
}

func buildAuthManagerClient(ctx context.Context, mockCtrl *gomock.Controller, finder Finder, username string, overrideGroup *permissionGroup, makeReadOnly bool) (*mock.MockAuthManager, error) {
	authManagerClient := mock.NewMockAuthManager(mockCtrl)
	for groupName, group := range permissions {
		switch groupName {
		case permissionVcenter:
			vcenter, err := finder.Folder(ctx, "/")
			if err != nil {
				return nil, err
			}
			buildPermissionGroup(authManagerClient, vcenter.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
		case permissionDatacenter:
			datacenters, err := finder.DatacenterList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datacenter := range datacenters {
				buildPermissionGroup(authManagerClient, datacenter.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
			}
		case permissionDatastore:
			datastores, err := finder.DatastoreList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, datastore := range datastores {
				buildPermissionGroup(authManagerClient, datastore.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
			}
		case permissionCluster:
			clusters, err := finder.ClusterComputeResourceList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, cluster := range clusters {
				buildPermissionGroup(authManagerClient, cluster.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
			}
		case permissionPortgroup:
			networks, err := finder.NetworkList(ctx, "/...")
			if err != nil {
				return nil, err
			}
			for _, network := range networks {
				buildPermissionGroup(authManagerClient, network.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
			}
		case permissionResourcePool:
			resourcePools := []string{"/DC0/host/DC0_C0/Resources/test-resourcepool", "/DC0/host/DC0_C0/Resources"}
			for _, resourcePoolPath := range resourcePools {
				resourcePool, err := finder.ResourcePool(ctx, resourcePoolPath)
				if err != nil {
					return nil, err
				}
				buildPermissionGroup(authManagerClient, resourcePool.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
			}
		case permissionFolder:
			var folders = []string{"/DC0/vm", "/DC0/vm/my-folder"}
			for _, folder := range folders {
				folder, err := finder.Folder(ctx, folder)
				if err != nil {
					return nil, err
				}
				buildPermissionGroup(authManagerClient, folder.Reference(), username, group, groupName, overrideGroup, makeReadOnly)
			}
		}
	}
	return authManagerClient, nil
}

func validIPIMultiZoneInstallConfig() *types.InstallConfig {
	installConfig := validIPIInstallConfig("DC0", "")
	validMultiZonePlatform := validMultiVCenterPlatform()
	installConfig.VSphere.VCenters = validMultiZonePlatform.VCenters
	installConfig.VSphere.FailureDomains = validMultiZonePlatform.FailureDomains

	return installConfig
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
	validMultiZoneInstallConfig := validIPIMultiZoneInstallConfig()

	userDefinedFolderInstallConfig := validIPIInstallConfig("DC0", "")
	userDefinedFolderInstallConfig.VSphere.Folder = "/DC0/vm/my-folder"

	invalidDatacenterInstallConfig := validIPIInstallConfig("DC0", "")
	invalidDatacenterInstallConfig.VSphere.Datacenter = "invalid"

	username := validInstallConfig.VSphere.Username

	validPermissionsAuthManagerClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, nil, false)
	if err != nil {
		t.Error(err)
		return
	}
	missingPortgroupPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionPortgroup, false)
	if err != nil {
		t.Error(err)
		return
	}

	missingVCenterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionVcenter, false)
	if err != nil {
		t.Error(err)
		return
	}

	missingClusterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionCluster, false)
	if err != nil {
		t.Error(err)
		return
	}

	readOnlyClusterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionCluster, true)
	if err != nil {
		t.Error(err)
		return
	}

	missingDatastorePermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionDatastore, false)
	if err != nil {
		t.Error(err)
		return
	}

	missingDatacenterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionDatacenter, false)
	if err != nil {
		t.Error(err)
		return
	}

	readOnlyDatacenterPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionDatacenter, true)
	if err != nil {
		t.Error(err)
		return
	}

	missingFolderPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionFolder, false)
	if err != nil {
		t.Error(err)
		return
	}

	missingResourcePoolPermissionsClient, err := buildAuthManagerClient(ctx, mockCtrl, finder, username, &permissionResourcePool, false)
	if err != nil {
		t.Error(err)
		return
	}

	tests := []struct {
		name                      string
		installConfig             *types.InstallConfig
		validationMethod          func(*validationContext, *types.InstallConfig) error
		multiZoneValidationMethod func(*validationContext, *vsphere.FailureDomain) field.ErrorList
		failureDomain             *vsphere.FailureDomain
		expectErr                 string
		authManager               AuthManager
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
			name: "missing cluster Permissions",
			installConfig: func() *types.InstallConfig {
				installConfig := validIPIInstallConfig("DC0", "")
				installConfig.VSphere.ResourcePool = ""
				return installConfig
			}(),
			validationMethod: validateProvisioning,
			authManager:      missingClusterPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Cluster: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:             "resource pool provided, compute cluster can have read-only",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      readOnlyClusterPermissionsClient,
		},
		{
			name:             "missing datacenter Permissions",
			installConfig:    validInstallConfig,
			validationMethod: validateProvisioning,
			authManager:      missingDatacenterPermissionsClient,
			expectErr:        "privileges missing for vSphere vCenter Datacenter: Folder.Delete",
		},
		{
			name: "user-defined folder provided, datacenter can have read-only",
			installConfig: func() *types.InstallConfig {
				installConfig := validIPIInstallConfig("DC0", "")
				installConfig.VSphere.Folder = "/DC0/vm/my-folder"
				return installConfig
			}(),
			validationMethod: validateProvisioning,
			authManager:      readOnlyDatacenterPermissionsClient,
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
		{
			name:                      "multi-zone valid Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               validPermissionsAuthManagerClient,
		},
		{
			name:                      "multi-zone missing portgroup Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               missingPortgroupPermissionsClient,
			expectErr:                 "privileges missing for vSphere Port Group: Network.Assign",
		},
		{
			name:                      "multi-zone missing vCenter Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               missingVCenterPermissionsClient,
			expectErr:                 "privileges missing for vSphere vCenter: StorageProfile.View",
		},
		{
			name:                      "multi-zone missing cluster Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				failureDomain.Topology.ResourcePool = ""
				return &failureDomain
			}(),
			authManager: missingClusterPermissionsClient,
			expectErr:   "privileges missing for vSphere vCenter Cluster: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:                      "multi-zone resource pool provided, compute cluster can have read-only",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               readOnlyClusterPermissionsClient,
		},
		{
			name:                      "multi-zone missing datacenter Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				failureDomain.Topology.Folder = ""
				return &failureDomain
			}(),
			authManager: missingDatacenterPermissionsClient,
			expectErr:   "privileges missing for vSphere vCenter Datacenter: Folder.Delete",
		},
		{
			name:                      "multi-zone user-defined folder provided, datacenter can have read-only",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               readOnlyDatacenterPermissionsClient,
		},
		{
			name:                      "multi-zone missing datastore Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               missingDatastorePermissionsClient,
			expectErr:                 "privileges missing for vSphere vCenter Datastore: InventoryService.Tagging.ObjectAttachable",
		},
		{
			name:                      "multi-zone missing resource pool Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               missingResourcePoolPermissionsClient,
			expectErr:                 "privileges missing for vSphere vCenter Resource Pool: VirtualMachine.Config.AddNewDisk",
		},
		{
			name:                      "multi-zone missing user-defined folder Permissions but no folder defined",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain: func() *vsphere.FailureDomain {
				failureDomain := validMultiZoneInstallConfig.VSphere.FailureDomains[0]
				failureDomain.Topology.Folder = ""
				failureDomain.Topology.ResourcePool = ""
				return &failureDomain
			}(),
			authManager: missingFolderPermissionsClient,
		},
		{
			name:                      "multi-zone missing user-defined folder Permissions",
			installConfig:             validMultiZoneInstallConfig,
			multiZoneValidationMethod: validateMultiZoneProvisioning,
			failureDomain:             &validMultiZoneInstallConfig.VSphere.FailureDomains[0],
			authManager:               missingFolderPermissionsClient,
			expectErr:                 "privileges missing for Pre-existing Virtual Machine Folder: VirtualMachine.Provisioning.DeployTemplate",
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
			if test.validationMethod != nil {
				err = test.validationMethod(validationCtx, test.installConfig)
			} else if test.multiZoneValidationMethod != nil {
				err = test.multiZoneValidationMethod(validationCtx, test.failureDomain).ToAggregate()
			} else {
				err = errors.New("no test method defined")
			}
			if test.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, test.expectErr, err)
			}
		})
	}
}
