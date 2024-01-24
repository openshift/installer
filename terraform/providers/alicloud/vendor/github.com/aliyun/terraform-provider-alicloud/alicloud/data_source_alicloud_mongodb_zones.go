package alicloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dds"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudMongoDBZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongoDBZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
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
						"multi_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongoDBZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := dds.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithDdsClient(func(ddsClient *dds.Client) (interface{}, error) {
		return ddsClient.DescribeRegions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_mongodb_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	regions, _ := raw.(*dds.DescribeRegionsResponse)
	if len(regions.Regions.DdsRegion) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available region for MongoDB."))
	}
	for _, r := range regions.Regions.DdsRegion {
		for _, zonid := range r.Zones.Zone {
			if multi && strings.Contains(zonid.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zonid.ZoneId)
				continue
			}
			if !multi && !strings.Contains(zonid.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zonid.ZoneId)
				continue
			}
		}
	}

	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
	}

	var s []map[string]interface{}
	if !multi {
		for _, zoneId := range zoneIds {
			mapping := map[string]interface{}{"id": zoneId}
			s = append(s, mapping)
		}
	} else {
		for _, zoneId := range zoneIds {
			mapping := map[string]interface{}{
				"id":             zoneId,
				"multi_zone_ids": splitMultiZoneId(zoneId),
			}
			s = append(s, mapping)
		}
	}
	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", zoneIds); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
