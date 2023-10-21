package alicloud

import (
	"fmt"
	"sort"
	"strings"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudKVStoreZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudKVStoreZoneRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaid"}, false),
			},
			"engine": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "redis",
				ValidateFunc: validation.StringInSlice([]string{"redis", "memcache"}, true),
			},
			"product_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Local", "Tair_rdb", "Tair_scm", "Tair_essd", "OnECS"}, false),
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

func dataSourceAlicloudKVStoreZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string
	instanceChargeType := d.Get("instance_charge_type").(string)

	request := r_kvstore.CreateDescribeAvailableResourceRequest()
	request.RegionId = client.RegionId
	request.InstanceChargeType = instanceChargeType
	request.Engine = d.Get("engine").(string)
	request.ProductType = d.Get("product_type").(string)

	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeAvailableResource(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_kvstore_zones", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	zones, _ := raw.(*r_kvstore.DescribeAvailableResourceResponse)
	if len(zones.AvailableZones.AvailableZone) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available zones for KVStore"))
	}
	for _, zone := range zones.AvailableZones.AvailableZone {
		if multi && strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zone.ZoneId)
			continue
		}
		if !multi && !strings.Contains(zone.ZoneId, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zone.ZoneId)
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
