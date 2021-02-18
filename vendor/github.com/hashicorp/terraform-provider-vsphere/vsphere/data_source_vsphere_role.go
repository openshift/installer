package vsphere

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"log"
	"strconv"
	"strings"
)

func dataSourceVsphereRole() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereRoleRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the role.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the role.",
			},
			"role_privileges": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Privileges to be associated with the role",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"label": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display label of the role.",
			},
		},
	}
}

func dataSourceVSphereRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] : Reading vsphere role with label %s", d.Get("label"))
	client := meta.(*VSphereClient).vimClient
	authorizationManager := object.NewAuthorizationManager(client.Client)

	label := d.Get("label").(string)
	roleList, err := authorizationManager.RoleList(context.Background())
	if err != nil {
		return fmt.Errorf("error while fetching the role list %s", err)
	}
	var foundRole = types.AuthorizationRole{}
	for _, role := range roleList {
		if role.Info != nil && role.Info.GetDescription() != nil {
			if label == role.Info.GetDescription().Label {
				foundRole = role
			}
		}
	}

	if foundRole.RoleId == 0 {
		return fmt.Errorf("role with label %s not found", label)
	}

	d.SetId(strconv.Itoa(int(foundRole.RoleId)))
	d.Set("name", foundRole.Name)
	d.Set("description", foundRole.Info.GetDescription().Summary)

	var arr []string
	for _, str := range foundRole.Privilege {
		if strings.Split(str, ".")[0] != SYSTEM_ROLE {
			arr = append(arr, str)
		}
	}
	d.Set("role_privileges", arr)
	return nil
}
