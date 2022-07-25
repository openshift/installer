package appid

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDUserRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Get a list of AppID user roles",
		ReadContext: dataSourceIBMAppIDUserRolesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"subject": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The user's identifier ('subject' in identity token)",
			},
			"roles": {
				Description: "A set of user roles",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "Role ID",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Role name",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDUserRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	subject := d.Get("subject").(string)

	roles, resp, err := appIDClient.GetUserRolesWithContext(ctx, &appid.GetUserRolesOptions{
		TenantID: &tenantID,
		ID:       &subject,
	})

	if err != nil {
		log.Printf("[DEBUG] Error getting AppID user roles: %s\n%s", err, resp)
		return diag.Errorf("Error getting AppID user roles: %s", err)
	}

	if roles.Roles != nil {
		if err := d.Set("roles", flattenAppIDUserRoles(roles.Roles)); err != nil {
			return diag.Errorf("Error setting AppID user roles: %s", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, subject))
	return nil
}

func flattenAppIDUserRoles(r []appid.GetUserRolesResponseRolesItem) []interface{} {
	var result []interface{}

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
