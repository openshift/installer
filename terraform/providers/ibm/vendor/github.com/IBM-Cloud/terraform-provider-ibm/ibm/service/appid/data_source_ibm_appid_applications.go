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

func DataSourceIBMAppIDApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDApplicationsRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_id": {
							Description: "The `client_id` is a public identifier for applications",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "The application name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"secret": {
							Description: "The `secret` is a secret known only to the application and the authorization server",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"oauth_server_url": {
							Description: "Base URL for common OAuth endpoints, like `/authorization`, `/token` and `/publickeys`",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"profiles_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"discovery_endpoint": {
							Description: "This URL returns OAuth Authorization Server Metadata",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Description: "The type of application to be registered. Allowed types are `regularwebapp` and `singlepageapp`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDApplicationsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	apps, resp, err := appIDClient.ListApplicationsWithContext(ctx, &appid.ListApplicationsOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error listing AppID applications: %s\n%s", err, resp)
	}

	applicationList := make([]interface{}, len(apps.Applications))

	for index, app := range apps.Applications {
		application := map[string]interface{}{}
		application["client_id"] = *app.ClientID
		application["name"] = *app.Name

		if app.Secret != nil {
			application["secret"] = *app.Secret
		}

		if app.OAuthServerURL != nil {
			application["oauth_server_url"] = *app.OAuthServerURL
		}

		if app.ProfilesURL != nil {
			application["profiles_url"] = *app.ProfilesURL
		}

		if app.DiscoveryEndpoint != nil {
			application["discovery_endpoint"] = *app.DiscoveryEndpoint
		}

		if app.Type != nil {
			application["type"] = *app.Type
		}

		applicationList[index] = application
	}

	sort.Slice(applicationList, func(a, b int) bool {
		appA := applicationList[a].(map[string]interface{})
		appB := applicationList[b].(map[string]interface{})
		return appA["name"].(string) < appB["name"].(string)
	})

	if err := d.Set("applications", applicationList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/applications", tenantID))
	return nil
}
