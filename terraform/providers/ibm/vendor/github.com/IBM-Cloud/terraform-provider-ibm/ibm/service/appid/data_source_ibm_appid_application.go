package appid

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDApplication() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDApplicationRead,
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
				Description: "The type of application. Allowed types are `regularwebapp` and `singlepageapp`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDApplicationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	clientID := d.Get("client_id").(string)

	app, resp, err := appIDClient.GetApplicationWithContext(ctx, &appid.GetApplicationOptions{
		TenantID: &tenantID,
		ClientID: &clientID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID application: %s\n%s", err, resp)
	}

	if app.Name != nil {
		d.Set("name", *app.Name)
	}

	if app.Secret != nil {
		d.Set("secret", *app.Secret)
	}

	if app.OAuthServerURL != nil {
		d.Set("oauth_server_url", *app.OAuthServerURL)
	}

	if app.ProfilesURL != nil {
		d.Set("profiles_url", *app.ProfilesURL)
	}

	if app.DiscoveryEndpoint != nil {
		d.Set("discovery_endpoint", *app.DiscoveryEndpoint)
	}

	if app.Type != nil {
		d.Set("type", *app.Type)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, clientID))
	return nil
}
