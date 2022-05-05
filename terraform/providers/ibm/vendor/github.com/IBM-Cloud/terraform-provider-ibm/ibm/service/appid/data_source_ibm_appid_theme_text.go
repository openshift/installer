package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDThemeText() *schema.Resource {
	return &schema.Resource{
		Description: "The theme texts of the App ID login widget",
		ReadContext: dataSourceIBMAppIDThemeTextRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"tab_title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"footnote": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMAppIDThemeTextRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	text, resp, err := appIDClient.GetThemeTextWithContext(ctx, &appid.GetThemeTextOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID theme text: %s\n%s", err, resp)
	}

	if text.TabTitle != nil {
		d.Set("tab_title", *text.TabTitle)
	}

	if text.Footnote != nil {
		d.Set("footnote", *text.Footnote)
	}

	d.SetId(tenantID)

	return nil
}
