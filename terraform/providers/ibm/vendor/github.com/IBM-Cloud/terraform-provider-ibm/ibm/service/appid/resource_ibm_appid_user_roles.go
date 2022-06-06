package appid

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppIDUserRoles() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage AppID user roles",
		ReadContext:   resourceIBMAppIDUserRolesRead,
		CreateContext: resourceIBMAppIDUserRolesCreate,
		DeleteContext: resourceIBMAppIDUserRolesDelete,
		UpdateContext: resourceIBMAppIDUserRolesUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AppID instance GUID",
			},
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The user's identifier ('subject' in identity token)",
			},
			"role_ids": {
				Description: "A set of AppID role IDs that should be assigned to the user",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIBMAppIDUserRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/subject", id)
	}

	tenantID := idParts[0]
	subject := idParts[1]

	d.Set("tenant_id", tenantID)
	d.Set("subject", subject)

	roles, resp, err := appIDClient.GetUserRolesWithContext(ctx, &appid.GetUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
	})

	if err != nil {
		log.Printf("[DEBUG] Error getting AppID user roles: %s\n%s", err, resp)
		return diag.Errorf("Error getting AppID user roles: %s", err)
	}

	if roles.Roles != nil {
		if err := d.Set("role_ids", flattenAppIDUserRoleIDs(roles.Roles)); err != nil {
			return diag.Errorf("Error setting AppID user role_ids: %s", err)
		}
	}

	return nil
}

func resourceIBMAppIDUserRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	subject := d.Get("subject").(string)
	roleIds := d.Get("role_ids").(*schema.Set)

	input := &appid.UpdateUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: flex.ExpandStringList(roleIds.List()),
		},
	}

	_, resp, err := appIDClient.UpdateUserRolesWithContext(ctx, input)

	if err != nil {
		log.Printf("[DEBUG] Error updating AppID user roles: %s\n%s", err, resp)
		return diag.Errorf("Error updating AppID user roles: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, subject))
	return resourceIBMAppIDUserRolesRead(ctx, d, meta)
}

func resourceIBMAppIDUserRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	subject := d.Get("subject").(string)

	input := &appid.UpdateUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: []string{},
		},
	}

	_, resp, err := appIDClient.UpdateUserRolesWithContext(ctx, input)

	if err != nil {
		log.Printf("[DEBUG] Error deleting AppID user roles: %s\n%s", err, resp)
		return diag.Errorf("Error deleting AppID user roles: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppIDUserRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMAppIDUserRolesCreate(ctx, d, meta)
}

func flattenAppIDUserRoleIDs(r []appid.GetUserRolesResponseRolesItem) []string {
	result := make([]string, len(r))
	for i, role := range r {
		result[i] = *role.ID
	}
	return result
}
