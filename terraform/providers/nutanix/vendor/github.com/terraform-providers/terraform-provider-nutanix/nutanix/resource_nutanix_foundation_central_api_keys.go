package nutanix

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fc "github.com/terraform-providers/terraform-provider-nutanix/client/fc"
)

func resourceNutanixFCAPIKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixFCAPIKeysCreate,
		ReadContext:   resourceNutanixFCAPIKeysRead,
		DeleteContext: resourceNutanixFCAPIKeysDelete,
		Schema: map[string]*schema.Schema{
			"alias": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"created_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"api_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNutanixFCAPIKeysCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral
	req := &fc.CreateAPIKeysInput{}

	alias, ok := d.GetOk("alias")
	if ok {
		req.Alias = alias.(string)
	}

	resp, err := conn.Service.CreateAPIKey(ctx, req)
	if err != nil {
		return diag.Errorf("error creating API Keys with alias %s: %+v", (req.Alias), err)
	}

	d.SetId(resp.KeyUUID)
	return resourceNutanixFCAPIKeysRead(ctx, d, meta)
}

func resourceNutanixFCAPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral
	resp, err := conn.Service.GetAPIKey(ctx, d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	d.Set("created_timestamp", resp.CreatedTimestamp)
	d.Set("key_uuid", resp.KeyUUID)
	d.Set("api_key", resp.APIKey)
	d.Set("current_time", resp.CurrentTime)
	d.Set("alias", resp.Alias)

	return nil
}

func resourceNutanixFCAPIKeysDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
