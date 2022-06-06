package appid

import (
	"context"
	"fmt"
	"sort"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDRoles() *schema.Resource {
	return &schema.Resource{
		Description: "A list of AppID roles",
		ReadContext: dataSourceIBMAppIDRolesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Description: "Role ID",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique role name",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional role description",
						},
						"access": {
							Type:     schema.TypeSet,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceIBMAppIDRolesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	roles, resp, err := appIDClient.ListRolesWithContext(ctx, &appid.ListRolesOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error listing AppID roles: %s\n%s", err, resp)
	}

	roleList := make([]interface{}, len(roles.Roles))

	for index, role := range roles.Roles {
		rMap := map[string]interface{}{
			"role_id": *role.ID,
		}

		if role.Name != nil {
			rMap["name"] = *role.Name
		}

		if role.Description != nil {
			rMap["description"] = *role.Description
		}

		rMap["access"] = flattenAppIDRoleAccess(role.Access)
		roleList[index] = rMap
	}

	// make this predictable for easier testing
	sort.Slice(roleList, func(a, b int) bool {
		roleA := roleList[a].(map[string]interface{})
		roleB := roleList[b].(map[string]interface{})
		return roleA["name"].(string) < roleB["name"].(string)
	})

	if err := d.Set("roles", roleList); err != nil {
		return diag.Errorf("Error setting AppID roles: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/roles", tenantID))
	return nil
}
