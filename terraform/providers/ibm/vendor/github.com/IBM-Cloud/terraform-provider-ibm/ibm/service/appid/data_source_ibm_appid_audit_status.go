package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		Description: "Tenant audit status",
		ReadContext: dataSourceIBMAppIDAuditStatusRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The AppID instance GUID",
			},
			"is_active": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The auditing status of the tenant.",
			},
		},
	}
}

func dataSourceIBMAppIDAuditStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	auditStatus, resp, err := appIDClient.GetAuditStatusWithContext(ctx, &appid.GetAuditStatusOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID audit status: %s\n%s", err, resp)
	}

	d.Set("is_active", *auditStatus.IsActive)
	d.SetId(tenantID)

	return nil
}
