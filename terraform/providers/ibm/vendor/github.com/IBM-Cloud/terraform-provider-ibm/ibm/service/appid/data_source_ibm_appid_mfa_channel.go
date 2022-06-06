package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDMFAChannel() *schema.Resource {
	return &schema.Resource{
		Description: "Get MFA channel configuration",
		ReadContext: dataSourceIBMAppIDMFAChannelRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"active": {
				Description: "Possible values: `email`, `sms`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"sms_config": {
				Description: "Configuration for `sms` channel",
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true, // terraform does not yet support nested sensitive attributes, this is temporary workaround
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Description: "API Key",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"secret": {
							Description: "API Secret",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"from": {
							Description: "Sender's phone number",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDMFAChannelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	ch, resp, err := appIDClient.ListChannelsWithContext(ctx, &appid.ListChannelsOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID MFA channels: %s\n%s", err, resp)
	}

	for _, channel := range ch.Channels {
		if *channel.IsActive {
			d.Set("active", *channel.Type)
		}

		if *channel.Type == "sms" && channel.Config != nil {
			config := map[string]interface{}{
				"key":    *channel.Config.Key,
				"secret": *channel.Config.Secret,
				"from":   *channel.Config.From,
			}

			if err := d.Set("sms_config", []interface{}{config}); err != nil {
				return diag.Errorf("Error setting AppID MFA channel config: %s", err)
			}
		}
	}

	d.SetId(tenantID)

	return nil
}
