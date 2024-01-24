package alicloud

import (
	"fmt"
	"sort"

	"github.com/aliyun/fc-go-sdk"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudFcZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudFcZonesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudFcZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var clientInfo *fc.Client
	raw, err := client.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		clientInfo = fcClient
		return fcClient.GetAccountSettings(fc.NewGetAccountSettingsInput())
	})
	if err != nil {
		return WrapError(fmt.Errorf("[API ERROR] FC GetAccountSettings: %#v", err))
	}
	addDebug("GetAccountSettings", raw, clientInfo)
	out, _ := raw.(*fc.GetAccountSettingsOutput)
	if out != nil && len(out.AvailableAZs) > 0 {
		sort.Strings(out.AvailableAZs)

		var zoneIds []string
		var s []map[string]interface{}
		for _, zoneId := range out.AvailableAZs {
			mapping := map[string]interface{}{"id": zoneId}
			s = append(s, mapping)
			zoneIds = append(zoneIds, zoneId)
		}

		d.SetId(dataResourceIdHash(out.AvailableAZs))
		if err := d.Set("zones", s); err != nil {
			return WrapError(err)
		}

		if err := d.Set("ids", zoneIds); err != nil {
			return WrapError(err)
		}
		if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
			writeToFile(output.(string), s)
		}
	}
	return nil
}
