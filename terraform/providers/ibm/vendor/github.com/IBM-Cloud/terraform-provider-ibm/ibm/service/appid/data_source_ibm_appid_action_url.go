package appid

import (
	"context"
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceIBMAppIDActionURL() *schema.Resource {
	return &schema.Resource{
		Description: "The custom url to redirect to when Cloud Directory action is executed.",
		ReadContext: dataSourceIBMAppIDActionURLRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"action": {
				Description:  "The type of the action: `on_user_verified` - the URL of your custom user verified page, `on_reset_password` - the URL of your custom reset password page",
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"on_user_verified", "on_reset_password"}, false),
				Required:     true,
			},
			"url": {
				Description: "The action URL",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDActionURLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	action := d.Get("action").(string)

	resp, rawResp, err := appIDClient.GetCloudDirectoryActionURLWithContext(ctx, &appid.GetCloudDirectoryActionURLOptions{
		TenantID: &tenantID,
		Action:   &action,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID actionURL: %s\n%s", err, rawResp)
	}

	if resp.ActionURL != nil {
		d.Set("url", *resp.ActionURL)
	}

	d.SetId(fmt.Sprintf("%s/%s", tenantID, action))

	return nil
}
