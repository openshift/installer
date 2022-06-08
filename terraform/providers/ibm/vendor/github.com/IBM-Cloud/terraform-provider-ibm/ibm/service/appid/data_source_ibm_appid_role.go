package appid

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDRole() *schema.Resource {
	return &schema.Resource{
		Description: "A role is a collection of `scopes` that allow varying permissions to different types of app users",
		ReadContext: dataSourceIBMAppIDRoleRead,
		Schema: map[string]*schema.Schema{
			"role_id": {
				Description: "Role ID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "Unique role name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "Optional role description",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"access": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "A set of access policies that bind specific application scopes to the role",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scopes": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	id := d.Get("role_id").(string)

	role, resp, err := appIDClient.GetRoleWithContext(ctx, &appid.GetRoleOptions{
		RoleID:   &id,
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

	d.SetId(fmt.Sprintf("%s/%s", tenantID, *role.ID))

	return nil
}

func flattenAppIDRoleAccess(ra []appid.RoleAccessItem) []interface{} {
	var result []interface{}

	for _, a := range ra {
		access := map[string]interface{}{
			"scopes": flex.FlattenStringList(a.Scopes),
		}

		if a.ApplicationID != nil {
			access["application_id"] = *a.ApplicationID
		}

		result = append(result, access)
	}

	return result
}
