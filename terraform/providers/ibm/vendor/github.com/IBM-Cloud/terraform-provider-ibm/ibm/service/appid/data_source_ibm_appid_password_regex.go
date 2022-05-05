package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDPasswordRegex() *schema.Resource {
	return &schema.Resource{
		Description: "The regular expression used by App ID for password strength validation",
		ReadContext: dataSourceIBMAppIDPasswordRegexRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
			},
			"base64_encoded_regex": {
				Description: "The regex expression rule for acceptable password encoded in base64",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"error_message": {
				Description: "Custom error message",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"regex": {
				Description: "The escaped regex expression rule for acceptable password",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDPasswordRegexRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	pw, resp, err := appIDClient.GetCloudDirectoryPasswordRegexWithContext(ctx, &appid.GetCloudDirectoryPasswordRegexOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading AppID Cloud Directory password regex: %s\n%s", err, resp)
	}

	if pw.Base64EncodedRegex != nil {
		d.Set("base64_encoded_regex", *pw.Base64EncodedRegex)
	}

	if pw.Regex != nil {
		d.Set("regex", *pw.Regex)
	}

	if pw.ErrorMessage != nil {
		d.Set("error_message", *pw.ErrorMessage)
	}

	d.SetId(tenantID)

	return nil
}
