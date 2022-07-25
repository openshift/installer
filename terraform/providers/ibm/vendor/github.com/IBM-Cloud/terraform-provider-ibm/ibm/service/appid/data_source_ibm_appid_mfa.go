package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDMFA() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDMFARead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "`true` if MFA is active",
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDMFARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	mfa, resp, err := appIDClient.GetMFAConfigWithContext(ctx, &appid.GetMFAConfigOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting IBM AppID MFA configuration: %s\n%s", err, resp)
	}

	if mfa.IsActive != nil {
		d.Set("is_active", *mfa.IsActive)
	}

	d.SetId(tenantID)

	return nil
}
