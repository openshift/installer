package appid

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDApplicationRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDApplicationRolesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"client_id": {
				Description: "The `client_id` is a public identifier for applications",
				Type:        schema.TypeString,
				Required:    true,
			},
			"roles": {
				Description: "Defined roles for an application that is registered with an App ID instance",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Description: "Application role ID",
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Description: "Application role name",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDApplicationRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	roles, resp, err := appIDClient.GetApplicationRolesWithContext(ctx, &appid.GetApplicationRolesOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID application roles: %s\n%s", err, resp)
	}

	if err := d.Set("roles", flattenAppIDApplicationRoles(roles.Roles)); err != nil {
		return diag.Errorf("Error setting AppID application roles: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return nil
}

func flattenAppIDApplicationRoles(r []appid.GetUserRolesResponseRolesItem) []interface{} {
	var result []interface{}

	if r == nil {
		return result
	}

	for _, v := range r {
		role := map[string]interface{}{
			"id": *v.ID,
		}

		if v.Name != nil {
			role["name"] = *v.Name
		}

		result = append(result, role)
	}

	return result
}
