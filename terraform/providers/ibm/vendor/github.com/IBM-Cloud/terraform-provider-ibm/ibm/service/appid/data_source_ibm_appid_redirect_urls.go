package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDRedirectURLs() *schema.Resource {
	return &schema.Resource{
		Description: "Redirect URIs that can be used as callbacks of App ID authentication flow",
		ReadContext: dataSourceIBMAppIDRedirectURLsRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The service `tenantId`",
			},
			"urls": {
				Description: "A list of redirect URLs",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDRedirectURLsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appidClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	urls, resp, err := appidClient.GetRedirectUrisWithContext(ctx, &appid.GetRedirectUrisOptions{
		TenantID: &tenantID,
	})
	if err != nil {
		return diag.Errorf("Error loading Cloud Directory AppID redirect urls: %s\n%s", err, resp)
	}

	if err := d.Set("urls", urls.RedirectUris); err != nil {
		return diag.Errorf("Error setting Cloud Directory AppID redirect URLs: %s", err)
	}

	d.SetId(tenantID)

	return nil
}
