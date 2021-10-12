package alicloud

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/polardb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudPolarDBZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudZonesPolarDBRead,

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

func dataSourceAlicloudZonesPolarDBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := polardb.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	raw, err := client.WithPolarDBClient(func(polarDBClient *polardb.Client) (interface{}, error) {
		return polarDBClient.DescribeRegions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_polardb_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	regions, _ := raw.(*polardb.DescribeRegionsResponse)
	if len(regions.Regions.Region) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available region for PolarDB."))
	}
	for _, r := range regions.Regions.Region {
		for _, zone := range r.Zones.Zone {
			if multi && strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zone.ZoneId)
				continue
			}
			if !multi && !strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
				zoneIds = append(zoneIds, zone.ZoneId)
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
