package vsphere

import (
	"context"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	vim25types "github.com/vmware/govmomi/vim25/types"
)

//go:generate mockgen -source=./permissions.go -destination=./mock/authmanager_generated.go -package=mock

// AuthManager defines an interface to an implementation of the AuthorizationManager to facilitate mocking
type AuthManager interface {
	FetchUserPrivilegeOnEntities(ctx context.Context, entities []vim25types.ManagedObjectReference, userName string) ([]vim25types.UserPrivilegeResult, error)
}

// permissionGroup is the group of permissions needed by cluster creation, operation, or teardown.
type permissionGroup string

// PermissionGroupDefinition defines a group of permissions and a related human friendly description
type PermissionGroupDefinition struct {
	/* Permissions array of privileges which correlate with the privileges listed in docs */
	Permissions []string
	/* Description friendly description of privilege group */
	Description string
}

// PrivilegeRelevanceFunc returns true if the associated privilege group privileges should be verified
type PrivilegeRelevanceFunc func(*vspheretypes.Platform) bool

var (
	permissionVcenter      permissionGroup = "vcenter"
	permissionCluster      permissionGroup = "cluster"
	permissionPortgroup    permissionGroup = "portgroup"
	permissionDatacenter   permissionGroup = "datacenter"
	permissionDatastore    permissionGroup = "datastore"
	permissionResourcePool permissionGroup = "resourcepool"
	permissionFolder       permissionGroup = "folder"
)

var permissions = map[permissionGroup]PermissionGroupDefinition{
	// Base set of permissions required for cluster creation
	permissionVcenter: {
		Permissions: []string{
			"Cns.Searchable",
			"InventoryService.Tagging.AttachTag",
			"InventoryService.Tagging.CreateCategory",
			"InventoryService.Tagging.CreateTag",
			"InventoryService.Tagging.DeleteCategory",
			"InventoryService.Tagging.DeleteTag",
			"InventoryService.Tagging.EditCategory",
			"InventoryService.Tagging.EditTag",
			"Sessions.ValidateSession",
			"StorageProfile.Update",
			"StorageProfile.View",
		},
		Description: "vSphere vCenter",
	},
	permissionCluster: {
		Permissions: []string{
			"Resource.AssignVMToPool",
			"VApp.AssignResourcePool",
			"VApp.Import",
			"VirtualMachine.Config.AddNewDisk",
		},
		Description: "vSphere vCenter Cluster",
	},
	permissionPortgroup: {
		Permissions: []string{
			"Network.Assign",
		},
		Description: "vSphere Port Group",
	},
	permissionFolder: {
		Permissions: []string{
			"Resource.AssignVMToPool",
			"VApp.Import",
			"VirtualMachine.Config.AddExistingDisk",
			"VirtualMachine.Config.AddNewDisk",
			"VirtualMachine.Config.AddRemoveDevice",
			"VirtualMachine.Config.AdvancedConfig",
			"VirtualMachine.Config.Annotation",
			"VirtualMachine.Config.CPUCount",
			"VirtualMachine.Config.DiskExtend",
			"VirtualMachine.Config.DiskLease",
			"VirtualMachine.Config.EditDevice",
			"VirtualMachine.Config.Memory",
			"VirtualMachine.Config.RemoveDisk",
			"VirtualMachine.Config.Rename",
			"VirtualMachine.Config.ResetGuestInfo",
			"VirtualMachine.Config.Resource",
			"VirtualMachine.Config.Settings",
			"VirtualMachine.Config.UpgradeVirtualHardware",
			"VirtualMachine.Interact.GuestControl",
			"VirtualMachine.Interact.PowerOff",
			"VirtualMachine.Interact.PowerOn",
			"VirtualMachine.Interact.Reset",
			"VirtualMachine.Inventory.Create",
			"VirtualMachine.Inventory.CreateFromExisting",
			"VirtualMachine.Inventory.Delete",
			"VirtualMachine.Provisioning.Clone",
			"VirtualMachine.Provisioning.MarkAsTemplate",
			"VirtualMachine.Provisioning.DeployTemplate",
		},
		Description: "Pre-existing Virtual Machine Folder",
	},
	permissionDatacenter: {
		Permissions: []string{
			"Resource.AssignVMToPool",
			"VApp.Import",
			"VirtualMachine.Config.AddExistingDisk",
			"VirtualMachine.Config.AddNewDisk",
			"VirtualMachine.Config.AddRemoveDevice",
			"VirtualMachine.Config.AdvancedConfig",
			"VirtualMachine.Config.Annotation",
			"VirtualMachine.Config.CPUCount",
			"VirtualMachine.Config.DiskExtend",
			"VirtualMachine.Config.DiskLease",
			"VirtualMachine.Config.EditDevice",
			"VirtualMachine.Config.Memory",
			"VirtualMachine.Config.RemoveDisk",
			"VirtualMachine.Config.Rename",
			"VirtualMachine.Config.ResetGuestInfo",
			"VirtualMachine.Config.Resource",
			"VirtualMachine.Config.Settings",
			"VirtualMachine.Config.UpgradeVirtualHardware",
			"VirtualMachine.Interact.GuestControl",
			"VirtualMachine.Interact.PowerOff",
			"VirtualMachine.Interact.PowerOn",
			"VirtualMachine.Interact.Reset",
			"VirtualMachine.Inventory.Create",
			"VirtualMachine.Inventory.CreateFromExisting",
			"VirtualMachine.Inventory.Delete",
			"VirtualMachine.Provisioning.Clone",
			"VirtualMachine.Provisioning.DeployTemplate",
			"VirtualMachine.Provisioning.MarkAsTemplate",
			"Folder.Create",
			"Folder.Delete",
		},
		Description: "vSphere vCenter Datacenter",
	},
	permissionDatastore: {
		Permissions: []string{
			"Datastore.AllocateSpace",
			"Datastore.Browse",
			"Datastore.FileManagement",
			"InventoryService.Tagging.ObjectAttachable",
		},
		Description: "vSphere vCenter Datastore",
	},
	permissionResourcePool: {
		Permissions: []string{
			"Resource.AssignVMToPool",
			"VApp.AssignResourcePool",
			"VApp.Import",
			"VirtualMachine.Config.AddNewDisk",
		},
		Description: "vSphere vCenter Resource Pool",
	},
}

func newAuthManager(client *vim25.Client) AuthManager {
	authManager := object.NewAuthorizationManager(client)
	return authManager
}

func comparePrivileges(ctx context.Context, validationCtx *validationContext, mo vim25types.ManagedObjectReference, permissionGroup PermissionGroupDefinition) error {
	authManager := validationCtx.AuthManager
	derived, err := authManager.FetchUserPrivilegeOnEntities(ctx, []vim25types.ManagedObjectReference{mo}, validationCtx.User)

	if err != nil {
		return errors.Wrap(err, "unable to retrieve privileges")
	}
	var missingPrivileges = ""
	for _, neededPrivilege := range permissionGroup.Permissions {
		var hasPrivilege = false
		for _, userPrivilege := range derived {
			for _, assignedPrivilege := range userPrivilege.Privileges {
				if assignedPrivilege == neededPrivilege {
					hasPrivilege = true
				}
			}
		}
		if hasPrivilege == false {
			if missingPrivileges != "" {
				missingPrivileges = missingPrivileges + ", "
			}
			missingPrivileges = missingPrivileges + neededPrivilege
		}
	}
	if missingPrivileges != "" {
		return errors.Errorf("privileges missing for %s: %s", permissionGroup.Description, missingPrivileges)
	}
	return nil
}
