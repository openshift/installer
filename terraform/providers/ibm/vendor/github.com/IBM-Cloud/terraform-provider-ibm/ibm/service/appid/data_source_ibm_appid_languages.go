package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDLanguages() *schema.Resource {
	return &schema.Resource{
		Description: "User localization configuration",
		ReadContext: dataSourceIBMAppIDLanguagesRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"languages": {
				Description: "The list of languages that can be used to customize email templates for Cloud Directory",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDLanguagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	langs, resp, err := appIDClient.GetLocalizationWithContext(ctx, &appid.GetLocalizationOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID languages: %s\n%s", err, resp)
	}

	d.Set("languages", langs.Languages)
	d.SetId(tenantID)

	return nil
}
