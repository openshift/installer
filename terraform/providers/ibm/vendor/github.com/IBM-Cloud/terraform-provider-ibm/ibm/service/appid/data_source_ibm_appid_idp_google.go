package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDIDPGoogle() *schema.Resource {
	return &schema.Resource{
		Description: "Returns the Google identity provider configuration.",
		ReadContext: dataSourceIBMAppIDIDPGoogleRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "`true` if Google IDP configuration is active",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"config": {
				Description: "Google IDP configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Google application id",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"application_secret": {
							Description: "Google application secret",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"redirect_url": {
				Description: "Paste the URI into the Authorized redirect URIs field in the Google Developer Console",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDIDPGoogleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	gg, resp, err := appIDClient.GetGoogleIDPWithContext(ctx, &appid.GetGoogleIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading AppID Google IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *gg.IsActive)

	if gg.RedirectURL != nil {
		d.Set("redirect_url", *gg.RedirectURL)
	}

	if gg.Config != nil {
		if err := d.Set("config", flattenIBMAppIDGoogleIDPConfig(gg.Config)); err != nil {
			return diag.Errorf("Failed setting AppID Google IDP config: %s", err)
		}
	}

	d.SetId(tenantID)

	return nil
}

func flattenIBMAppIDGoogleIDPConfig(config *appid.GoogleConfigParamsConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}
	mConfig["application_id"] = *config.IDPID
	mConfig["application_secret"] = *config.Secret

	return []interface{}{mConfig}
}
