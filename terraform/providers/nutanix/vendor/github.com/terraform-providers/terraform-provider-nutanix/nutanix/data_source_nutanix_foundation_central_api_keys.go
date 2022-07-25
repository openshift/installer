package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNutanixFCAPIKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFCAPIKeysRead,
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_uuid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNutanixFCAPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral

	if uuid, uuidok := d.GetOk("key_uuid"); uuidok {
		resp, err := conn.Service.GetAPIKey(ctx, uuid.(string))
		if err != nil {
			return diag.Errorf("error reading API keys with error %s", err)
		}
		if err := d.Set("created_timestamp", resp.CreatedTimestamp); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("api_key", resp.APIKey); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("current_time", resp.CurrentTime); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("alias", resp.Alias); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(uuid.(string))
	} else {
		return diag.Errorf("please provide `key_uuid`")
	}
	return nil
}
