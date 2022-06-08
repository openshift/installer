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

func ResourceIBMAppIDApplicationRoles() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAppIDApplicationRolesCreate,
		ReadContext:   resourceIBMAppIDApplicationRolesRead,
		DeleteContext: resourceIBMAppIDApplicationRolesDelete,
		UpdateContext: resourceIBMAppIDApplicationRolesUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"roles": {
				Description: "A list of role IDs for roles that you want to be assigned to the application (this is different from AppID role access)",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDApplicationRolesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	roles := flex.ExpandStringList(d.Get("roles").([]interface{}))

	roleOpts := &appid.PutApplicationsRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: roles,
		},
	}

	_, resp, err := appIDClient.PutApplicationsRolesWithContext(ctx, roleOpts)

	if err != nil {
		return diag.Errorf("Error setting application roles: %s\n%s", err, resp)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))

	return resourceIBMAppIDApplicationRolesRead(ctx, d, meta)
}

func resourceIBMAppIDApplicationRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	idParts := strings.Split(id, "/")

	if len(idParts) < 2 {
		return diag.Errorf("Incorrect ID %s: ID should be a combination of tenantID/clientID", d.Id())
	}

	tenantID := idParts[0]
	clientID := idParts[1]

	roles, resp, err := appIDClient.GetApplicationRolesWithContext(ctx, &appid.GetApplicationRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID application '%s' is not found, removing roles from state", clientID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID application roles: %s\n%s", err, resp)
	}

	var appRoles []interface{}

	if roles.Roles != nil {
		for _, v := range roles.Roles {
			appRoles = append(appRoles, *v.ID)
		}
	}

	if err := d.Set("roles", appRoles); err != nil {
		return diag.Errorf("Error setting application roles: %s", err)
	}

	d.Set("tenant_id", tenantID)
	d.Set("client_id", clientID)

	return nil
}

func resourceIBMAppIDApplicationRolesUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)
	roles := flex.ExpandStringList(d.Get("roles").([]interface{}))

	roleOpts := &appid.PutApplicationsRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: roles,
		},
	}

	_, resp, err := appIDClient.PutApplicationsRolesWithContext(ctx, roleOpts)

	if err != nil {
		return diag.Errorf("Error updating application roles: %s\n%s", err, resp)
	}

	return resourceIBMAppIDApplicationRolesRead(ctx, d, meta)
}

func resourceIBMAppIDApplicationRolesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	roleOpts := &appid.PutApplicationsRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
		Roles: &appid.UpdateUserRolesParamsRoles{
			Ids: []string{},
		},
	}

	_, resp, err := appIDClient.PutApplicationsRolesWithContext(ctx, roleOpts)

	if err != nil {
		return diag.Errorf("Error clearing application roles: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
