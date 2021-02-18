package vsphere

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/object"
	"log"
	"sort"
	"strconv"
	"strings"
)

const SYSTEM_ROLE = "System"

func resourceVsphereRole() *schema.Resource {
	sch := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the storage policy.",
		},
		"role_privileges": {
			Type:             schema.TypeList,
			Optional:         true,
			Description:      "The privileges to be associated with the role.",
			Elem:             &schema.Schema{Type: schema.TypeString},
			DiffSuppressFunc: privilegesDiffCheck,
		},
		"label": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The display label of the role.",
		},
	}

	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Schema: sch,
	}
}

func resourceRoleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning create role %s", d.Get("name").(string))
	client := meta.(*VSphereClient).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)

	name := d.Get("name").(string)
	rolePrivileges := structure.SliceInterfacesToStrings(d.Get("role_privileges").([]interface{}))

	roleId, err := authorizationManager.AddRole(context.Background(), name, rolePrivileges)
	if err != nil {
		return fmt.Errorf("error while creating role with name %s %s", name, err)
	}

	d.SetId(strconv.Itoa(int(roleId)))
	return resourceRoleRead(d, meta)
}

func resourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading vm role with id %s", d.Id())
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)
	roleIdInt, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("error while coverting role id %s from string to int %s", d.Id(), err)
	}
	roleId := int32(roleIdInt)
	roleList, err := authorizationManager.RoleList(context.Background())

	if err != nil {
		return fmt.Errorf("error while reading the role list %s", err)
	}
	role := roleList.ById(roleId)
	if role == nil {
		log.Printf(" [DEBUG] Role %s doesn't exist", d.Get("name"))
		d.SetId("")
		return nil
	}

	d.Set("name", role.Name)
	if role.Info != nil && role.Info.GetDescription() != nil {
		d.Set("label", role.Info.GetDescription().Label)
	}

	var privilegesArr []string
	for _, str := range role.Privilege {
		if strings.Split(str, ".")[0] != SYSTEM_ROLE {
			privilegesArr = append(privilegesArr, str)
		}
	}
	d.Set("role_privileges", privilegesArr)
	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning update role %s", d.Get("name").(string))
	client := meta.(*VSphereClient).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)
	name := d.Get("name").(string)
	roleIdInt, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("error while coverting role id %s from string to int %s", d.Id(), err)
	}
	roleId := int32(roleIdInt)
	rolePrivileges := structure.SliceInterfacesToStrings(d.Get("role_privileges").([]interface{}))

	err = authorizationManager.UpdateRole(context.Background(), roleId, name, rolePrivileges)
	if err != nil {
		return fmt.Errorf("error while updating role with is %d %s", roleIdInt, err)
	}
	return resourceRoleRead(d, meta)
}

func resourceRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Performing Delete of Role with ID %s", d.Id())
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)

	roleIdInt, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("error while coverting role id %s from string to int %s", d.Id(), err)
	}
	roleId := int32(roleIdInt)
	err = authorizationManager.RemoveRole(context.Background(), roleId, true)
	if err != nil {
		return fmt.Errorf("error while deleting role with id %d %s", roleIdInt, err)
	}

	d.SetId("")
	log.Printf(" [DEBUG] %s: Delete complete", d.Id())
	return nil
}

func privilegesDiffCheck(k, old, new string, d *schema.ResourceData) bool {

	oldVal, newVal := d.GetChange("role_privileges")
	oldArr := structure.SliceInterfacesToStrings(oldVal.([]interface{}))
	newArr := structure.SliceInterfacesToStrings(newVal.([]interface{}))

	if len(oldArr) != len(newArr) {
		return false
	}
	sort.Strings(oldArr)
	sort.Strings(newArr)
	for i := 0; i < len(oldArr); i++ {
		if newArr[i] != oldArr[i] {
			return false
		}
	}
	return true
}
