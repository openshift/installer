package vsphere

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/administrationroles"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/utils"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func resourceVsphereEntityPermissions() *schema.Resource {
	sch := map[string]*schema.Schema{
		"entity_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The managed object id or uuid of the entity.",
		},
		"entity_type": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The entity managed object type.",
		},
		"permissions": {
			Type:             schema.TypeList,
			Required:         true,
			MinItems:         1,
			Description:      "Permissions to be given to the entity.",
			Elem:             &schema.Resource{Schema: administrationroles.VspherePermissionSchema()},
			DiffSuppressFunc: permissionsDiffSuppressFunc,
		},
	}

	return &schema.Resource{
		Create:        resourceEntityPermissionsCreate,
		Read:          resourceEntityPermissionsRead,
		Update:        resourceEntityPermissionsUpdate,
		Delete:        resourceEntityPermissionsDelete,
		CustomizeDiff: resourceVSphereEntityPermissionsCustomizeDiff,
		Schema:        sch,
	}
}

func resourceEntityPermissionsCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning create permission for entity id %s", d.Get("entity_id").(string))
	client := meta.(*VSphereClient).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)
	permissions := d.Get("permissions").([]interface{})

	entityType := d.Get("entity_type").(string)
	entityId := d.Get("entity_id").(string)
	entityMoid, err := utils.GetMoid(client, entityType, entityId)
	if err != nil {
		return err
	}
	entityMor := types.ManagedObjectReference{
		Type:  entityType,
		Value: entityMoid,
	}

	usersAndGroupsMap := make(map[string]bool)
	var permissionObjs []types.Permission
	for _, permission := range permissions {
		userOrGroup := permission.(map[string]interface{})["user_or_group"].(string)
		if usersAndGroupsMap[userOrGroup] {
			return fmt.Errorf("user/group %s repeated, there is already a permission defined for the user/group", userOrGroup)
		}
		usersAndGroupsMap[userOrGroup] = true
		roleIdInt, err := strconv.ParseInt(permission.(map[string]interface{})["role_id"].(string), 10, 32)
		if err != nil {
			return fmt.Errorf("error while converting role id %s to integer", permission.(map[string]interface{})["role_id"].(string))
		}
		roleId := int32(roleIdInt)
		permissionObjs = append(permissionObjs, types.Permission{
			Principal: userOrGroup,
			Group:     permission.(map[string]interface{})["is_group"].(bool),
			Propagate: permission.(map[string]interface{})["propagate"].(bool),
			RoleId:    roleId,
		})
	}
	err = authorizationManager.SetEntityPermissions(context.Background(), entityMor, permissionObjs)
	if err != nil {
		return fmt.Errorf("error while creating permission for entity id %s %s", entityId, err)
	}
	d.SetId(entityMoid)
	return resourceEntityPermissionsRead(d, meta)
}

func resourceEntityPermissionsRead(d *schema.ResourceData, meta interface{}) error {
	entityType := d.Get("entity_type").(string)
	log.Printf(" [DEBUG] : Reading vm entity permissions for entity id %s and type %s", d.Id(), entityType)
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)

	entityMor := types.ManagedObjectReference{
		Type:  entityType,
		Value: d.Id(),
	}
	permissionsArr, err := authorizationManager.RetrieveEntityPermissions(context.Background(), entityMor, false)
	if err != nil {
		return fmt.Errorf("error while reading permissions for entity %s %s", d.Id(), err)
	}
	if len(permissionsArr) == 0 {
		log.Printf(" [DEBUG] :the permissions for entity with id %s and type %s is not found", d.Id(), entityType)
		d.SetId("")
		return nil
	}

	var permissionObjs []map[string]interface{}
	for _, permission := range permissionsArr {
		permissionObj := make(map[string]interface{})
		permissionObj["user_or_group"] = permission.Principal
		permissionObj["is_group"] = permission.Group
		permissionObj["propagate"] = permission.Propagate
		permissionObj["role_id"] = strconv.Itoa(int(permission.RoleId))
		permissionObjs = append(permissionObjs, permissionObj)
	}

	sort.Slice(permissionObjs, func(i, j int) bool {
		return strings.ToLower(permissionObjs[i]["user_or_group"].(string)) <
			strings.ToLower(permissionObjs[j]["user_or_group"].(string))
	})
	d.Set("permissions", permissionObjs)
	return nil
}

func resourceEntityPermissionsUpdate(d *schema.ResourceData, meta interface{}) error {
	if d.HasChange("permissions") {
		oldPermissions, newPermissions := d.GetChange("permissions")
		log.Printf(" [DEBUG] : Beginning update Permission with entity id %s", d.Id())

		client := meta.(*VSphereClient).vimClient
		authorizationManager := object.NewAuthorizationManager(client.Client)

		entityType := d.Get("entity_type").(string)
		entityMor := types.ManagedObjectReference{
			Type:  entityType,
			Value: d.Id(),
		}

		usersAndGroups := make(map[string]bool)
		var permissionObjs []types.Permission
		for _, permission := range newPermissions.([]interface{}) {
			userOrGroup := permission.(map[string]interface{})["user_or_group"].(string)
			if usersAndGroups[strings.ToLower(userOrGroup)] {
				return fmt.Errorf("user/group %s repeated, there is already a permission defined for the user/group", userOrGroup)
			}
			usersAndGroups[strings.ToLower(userOrGroup)] = true
			roleIdInt, err := strconv.ParseInt(permission.(map[string]interface{})["role_id"].(string), 10, 32)
			if err != nil {
				return fmt.Errorf("error while converting role id %s to integer", permission.(map[string]interface{})["role_id"].(string))
			}
			roleId := int32(roleIdInt)
			permissionObjs = append(permissionObjs, types.Permission{
				Principal: userOrGroup,
				Group:     permission.(map[string]interface{})["is_group"].(bool),
				Propagate: permission.(map[string]interface{})["propagate"].(bool),
				RoleId:    roleId,
			})
		}
		err := authorizationManager.SetEntityPermissions(context.Background(), entityMor, permissionObjs)
		if err != nil {
			return fmt.Errorf("error while updating permissions for entity id %s %s", d.Id(), err)
		}

		// handle removed permissions
		for _, permission := range oldPermissions.([]interface{}) {

			userOrGroup := permission.(map[string]interface{})["user_or_group"].(string)
			isGroup := permission.(map[string]interface{})["is_group"].(bool)

			if !usersAndGroups[strings.ToLower(userOrGroup)] {
				log.Printf(" [DEBUG] Deleting permissions for user/group %s", userOrGroup)
				err = authorizationManager.RemoveEntityPermission(context.Background(), entityMor, userOrGroup, isGroup)
				if err != nil {
					return fmt.Errorf("error while deleting permission for the user/group %s %s", userOrGroup, err)
				}
			}
		}
	}
	return resourceEntityPermissionsRead(d, meta)
}

func resourceEntityPermissionsDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf(" [DEBUG] Performing Delete of Entity permission %s", d.Id())
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)
	entityType := d.Get("entity_type").(string)
	entityMor := types.ManagedObjectReference{
		Type:  entityType,
		Value: d.Id(),
	}

	permissions := d.Get("permissions").([]interface{})
	for _, permission := range permissions {

		userOrGroup := permission.(map[string]interface{})["user_or_group"].(string)
		isGroup := permission.(map[string]interface{})["is_group"].(bool)
		log.Printf(" [DEBUG] Deleting permissions for user/group %s", userOrGroup)
		err := authorizationManager.RemoveEntityPermission(context.Background(), entityMor, userOrGroup, isGroup)
		if err != nil {
			return fmt.Errorf("error while deleting permission for the user/group %s %s", userOrGroup, err)
		}
	}
	d.SetId("")
	log.Printf("[DEBUG] %s: Delete complete for Entity Permissions", d.Id())
	return nil
}

func resourceVSphereEntityPermissionsCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	if d.HasChange("entity_id") {
		oldEntityId, newEntityId := d.GetChange("entity_id")
		if oldEntityId.(string) != "" {
			return fmt.Errorf("change %s in entity id is not allowed post creation", newEntityId)
		}
	}
	if d.HasChange("entity_type") {
		oldEntityType, newEntityType := d.GetChange("entity_type")
		if oldEntityType.(string) != "" {
			return fmt.Errorf("change in entity type %s is not allowed post creation", newEntityType)
		}
	}
	return nil
}

func permissionsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	oldPermissions, newPermissions := d.GetChange("permissions")
	oldPermissionsArr := oldPermissions.([]interface{})
	newPermissionsArr := newPermissions.([]interface{})
	if len(oldPermissionsArr) != len(newPermissionsArr) {
		return false
	}
	for _, oldPermission := range oldPermissionsArr {
		oldPermission.(map[string]interface{})["user_or_group"] =
			strings.ToLower(oldPermission.(map[string]interface{})["user_or_group"].(string))
	}
	for _, newPermission := range newPermissionsArr {
		newPermission.(map[string]interface{})["user_or_group"] =
			strings.ToLower(newPermission.(map[string]interface{})["user_or_group"].(string))
	}
	sort.Slice(oldPermissionsArr, func(i, j int) bool {
		return oldPermissionsArr[i].(map[string]interface{})["user_or_group"].(string) <
			oldPermissionsArr[j].(map[string]interface{})["user_or_group"].(string)
	})
	sort.Slice(newPermissionsArr, func(i, j int) bool {
		return newPermissionsArr[i].(map[string]interface{})["user_or_group"].(string) <
			newPermissionsArr[j].(map[string]interface{})["user_or_group"].(string)
	})
	return reflect.DeepEqual(oldPermissionsArr, newPermissionsArr)
}
