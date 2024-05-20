package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	fc "github.com/terraform-providers/terraform-provider-nutanix/client/fc"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func dataSourceNutanixFCListAPIKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNutanixFCListAPIKeysRead,
		Schema: map[string]*schema.Schema{
			"length": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_matches": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"offset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"api_keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceNutanixFCListAPIKeysRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).FoundationCentral

	req := &fc.ListMetadataInput{}
	length, lok := d.GetOk("length")
	if lok {
		req.Length = utils.IntPtr(length.(int))
	}
	offset, ok := d.GetOk("offset")
	if ok {
		req.Offset = utils.IntPtr(offset.(int))
	}

	resp, err := conn.Service.ListAPIKeys(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.Metadata != nil {
		metalist := make([]map[string]interface{}, 0)
		meta := make(map[string]interface{})
		meta["length"] = (resp.Metadata.Length)
		meta["offset"] = (resp.Metadata.Offset)
		meta["total_matches"] = (resp.Metadata.TotalMatches)

		metalist = append(metalist, meta)
		d.Set("metadata", metalist)
	}

	list := flattenAPIKeysList(resp.APIKeys)
	d.Set("api_keys", list)

	d.SetId(resource.UniqueId())
	return nil
}

func flattenAPIKeysList(pr []*fc.CreateAPIKeysResponse) []map[string]interface{} {
	resp := make([]map[string]interface{}, len(pr))

	for k, v := range pr {
		manage := make(map[string]interface{})
		manage["alias"] = v.Alias
		manage["api_key"] = v.APIKey
		manage["created_timestamp"] = v.CreatedTimestamp
		manage["current_time"] = v.CurrentTime
		manage["key_uuid"] = v.KeyUUID

		resp[k] = manage
	}

	return resp
}
