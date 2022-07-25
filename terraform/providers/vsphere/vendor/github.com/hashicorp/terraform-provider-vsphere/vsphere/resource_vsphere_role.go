package vsphere

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/object"
)

const SystemRole = "System"

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
	client := meta.(*Client).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)

	name := d.Get("name").(string)
	rolePrivileges := structure.SliceInterfacesToStrings(d.Get("role_privileges").([]interface{}))

	roleID, err := authorizationManager.AddRole(context.Background(), name, rolePrivileges)
	if err != nil {
		return fmt.Errorf("error while creating role with name %s %s", name, err)
	}

	d.SetId(strconv.Itoa(int(roleID)))
	return resourceRoleRead(d, meta)
}

func resourceRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Reading vm role with id %s", d.Id())
	client := meta.(*Client).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)
	roleIDInt, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("error while coverting role id %s from string to int %s", d.Id(), err)
	}
	roleID := int32(roleIDInt)
	roleList, err := authorizationManager.RoleList(context.Background())

	if err != nil {
		return fmt.Errorf("error while reading the role list %s", err)
	}
	role := roleList.ById(roleID)
	if role == nil {
		log.Printf(" [DEBUG] Role %s doesn't exist", d.Get("name"))
		d.SetId("")
		return nil
	}

	_ = d.Set("name", role.Name)
	if role.Info != nil && role.Info.GetDescription() != nil {
		_ = d.Set("label", role.Info.GetDescription().Label)
	}

	var privilegesArr []string
	for _, str := range role.Privilege {
		if strings.Split(str, ".")[0] != SystemRole {
			privilegesArr = append(privilegesArr, str)
		}
	}
	_ = d.Set("role_privileges", privilegesArr)
	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Beginning update role %s", d.Get("name").(string))
	client := meta.(*Client).vimClient

	authorizationManager := object.NewAuthorizationManager(client.Client)
	name := d.Get("name").(string)
	roleIDInt, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("error while coverting role id %s from string to int %s", d.Id(), err)
	}
	roleID := int32(roleIDInt)
	rolePrivileges := structure.SliceInterfacesToStrings(d.Get("role_privileges").([]interface{}))

	err = authorizationManager.UpdateRole(context.Background(), roleID, name, rolePrivileges)
	if err != nil {
		return fmt.Errorf("error while updating role with is %d %s", roleIDInt, err)
	}
	return resourceRoleRead(d, meta)
}

func resourceRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Performing Delete of Role with ID %s", d.Id())
	client := meta.(*Client).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)

	roleIDInt, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("error while coverting role id %s from string to int %s", d.Id(), err)
	}
	roleID := int32(roleIDInt)
	err = authorizationManager.RemoveRole(context.Background(), roleID, true)
	if err != nil {
		return fmt.Errorf("error while deleting role with id %d %s", roleIDInt, err)
	}

	d.SetId("")
	log.Printf(" [DEBUG] %s: Delete complete", d.Id())
	return nil
}

func privilegesDiffCheck(_, _, _ string, d *schema.ResourceData) bool {
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
