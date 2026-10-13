package vsphere

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	vim25types "github.com/vmware/govmomi/vim25/types"
	"k8s.io/apimachinery/pkg/util/sets"

	vspheretypes "github.com/openshift/installer/pkg/types/vsphere"
)

//go:generate mockgen -source=./permissions.go -destination=./mock/authmanager_generated.go -package=mock

// AuthManager defines an interface to an implementation of the AuthorizationManager to facilitate mocking
type AuthManager interface {
	FetchUserPrivilegeOnEntities(ctx context.Context, entities []vim25types.ManagedObjectReference, userName string) ([]vim25types.UserPrivilegeResult, error)
	Properties(ctx context.Context, r vim25types.ManagedObjectReference, ps []string, dst interface{}) error
	Reference() vim25types.ManagedObjectReference
}

// SessionManager defines an interface to an implementation of the SessionManager to facilitate mocking.
type SessionManager interface {
	UserSession(ctx context.Context) (*vim25types.UserSession, error)
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
			"InventoryService.Tagging.ObjectAttachable",
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
			"InventoryService.Tagging.ObjectAttachable",
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

// pruneToAvailablePermissions different versions of vCenter support different privileges.  the intent of this method
// is to prune privileges from the check that don't exist.
func pruneToAvailablePermissions(ctx context.Context, manager AuthManager) error {
	var authManagerMo mo.AuthorizationManager
	err := manager.Properties(ctx, manager.Reference(), []string{"privilegeList"}, &authManagerMo)
	if err != nil {
		return err
	}

	availablePermissions := sets.NewString()
	for _, availablePermission := range authManagerMo.PrivilegeList {
		availablePermissions.Insert(availablePermission.PrivId)
	}
	for permissionGroupKey, permissionGroup := range permissions {
		prunedPermissions := sets.NewString(permissionGroup.Permissions...)

		prunedPermissions = prunedPermissions.Intersection(availablePermissions)
		permissions[permissionGroupKey] = PermissionGroupDefinition{
			Permissions: prunedPermissions.List(),
			Description: permissionGroup.Description,
		}
	}
	return nil
}

func newAuthManager(client *vim25.Client) AuthManager {
	authManager := object.NewAuthorizationManager(client)
	return authManager
}

func comparePrivileges(ctx context.Context, validationCtx *validationContext, moRef vim25types.ManagedObjectReference, permissionGroup PermissionGroupDefinition) error {
	authManager := validationCtx.AuthManager
	sessionMgr := session.NewManager(validationCtx.Client)
	user, err := sessionMgr.UserSession(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to get user session")
	}
	derived, err := authManager.FetchUserPrivilegeOnEntities(ctx, []vim25types.ManagedObjectReference{moRef}, user.UserName)
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
		if !hasPrivilege {
			if missingPrivileges != "" {
				missingPrivileges += ", "
			}
			missingPrivileges += neededPrivilege
		}
	}
	if missingPrivileges != "" {
		return errors.Errorf("privileges missing for %s: %s", permissionGroup.Description, missingPrivileges)
	}
	return nil
}
