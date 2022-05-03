package appid

import (
	"context"
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppIDRole() *schema.Resource {
	return &schema.Resource{
		Description:   "A role is a collection of `scopes` that allow varying permissions to different types of app users",
		CreateContext: resourceIBMAppIDRoleCreate,
		ReadContext:   resourceIBMAppIDRoleRead,
		DeleteContext: resourceIBMAppIDRoleDelete,
		UpdateContext: resourceIBMAppIDRoleUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"role_id": {
				Description: "Role ID",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "Unique role name",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Optional role description",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"access": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Application `client_id`",
							Type:        schema.TypeString,
							Required:    true,
						},
						"scopes": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceIBMAppIDRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	roleName := d.Get("name").(string)

	input := &appid.CreateRoleOptions{
		Name:     &roleName,
		TenantID: &tenantID,
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = helpers.String(description.(string))
	}

	input.Access = expandAppIDRoleAccess(d.Get("access").(*schema.Set).List())

	role, resp, err := appIDClient.CreateRoleWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error creating AppID role: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *role.ID))

	return resourceIBMAppIDRoleRead(ctx, d, meta)
}

func resourceIBMAppIDRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/roleID", d.Id())
	}

	tenantID := idParts[0]
	roleID := idParts[1]

	role, resp, err := appIDClient.GetRoleWithContext(ctx, &appid.GetRoleOptions{
		RoleID:   &roleID,
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading AppID role: %s\n%s", err, resp)
	}

	d.Set("name", *role.Name)

	if role.Description != nil {
		d.Set("description", *role.Description)
	}

	if err := d.Set("access", flattenAppIDRoleAccess(role.Access)); err != nil {
		return diag.Errorf("Error setting AppID role access: %s", err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("role_id", roleID)

	return nil
}

func resourceIBMAppIDRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	tenantID := idParts[0]
	roleID := idParts[1]

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/roleID", d.Id())
	}

	resp, err := appIDClient.DeleteRoleWithContext(ctx, &appid.DeleteRoleOptions{
		TenantID: &tenantID,
		RoleID:   &roleID,
	})

	if err != nil {
		return diag.Errorf("Error deleting AppID role: %s\n%s", err, resp)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppIDRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// AppID role resource does not support partial updates, all inputs should be included
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	tenantID := idParts[0]
	roleID := idParts[1]
	roleName := d.Get("name").(string)

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/roleID", d.Id())
	}

	input := &appid.UpdateRoleOptions{
		TenantID: &tenantID,
		RoleID:   &roleID,
		Name:     &roleName,
	}

	if description, ok := d.GetOk("description"); ok {
		input.Description = helpers.String(description.(string))
	}

	input.Access = expandAppIDRoleAccess(d.Get("access").(*schema.Set).List())

	_, resp, err := appIDClient.UpdateRoleWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating AppID role: %s\n%s", err, resp)
	}

	return dataSourceIBMAppIDRoleRead(ctx, d, meta)
}

func expandAppIDRoleAccess(l []interface{}) []appid.RoleAccessItem {
	if len(l) == 0 {
		return []appid.RoleAccessItem{}
	}

	result := make([]appid.RoleAccessItem, len(l))

	for i, item := range l {
		aMap := item.(map[string]interface{})

		access := &appid.RoleAccessItem{
			ApplicationID: helpers.String(aMap["application_id"].(string)),
		}

		if scopes, ok := aMap["scopes"].([]interface{}); ok && len(scopes) > 0 {
			access.Scopes = []string{}

			for _, s := range scopes {
				access.Scopes = append(access.Scopes, s.(string))
			}
		}

		result[i] = *access
	}

	return result
}
