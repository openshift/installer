package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDIDPCustom() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDIDPCustomRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_key": {
				Description: "This is the public key used to validate your signed JWT. It is required to be a PEM in the RS256 or greater format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceIBMAppIDIDPCustomRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	config, resp, err := appIDClient.GetCustomIDPWithContext(ctx, &appid.GetCustomIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading AppID custom IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *config.IsActive)

	if config.Config != nil && config.Config.PublicKey != nil {
		if err := d.Set("public_key", *config.Config.PublicKey); err != nil {
			return diag.Errorf("failed setting config: %s", err)
		}
	}

	d.SetId(tenantID)

	return nil
}
