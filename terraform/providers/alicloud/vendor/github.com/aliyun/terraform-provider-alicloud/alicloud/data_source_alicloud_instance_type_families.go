package alicloud

import (
	"sort"

	"github.com/denverdino/aliyungo/common"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudInstanceTypeFamilies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudInstanceTypeFamiliesRead,

		Schema: map[string]*schema.Schema{
			"generation": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ecs-1", "ecs-2", "ecs-3", "ecs-4"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_charge_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  PostPaid,
				// %q must contain a valid InstanceChargeType, expected common.PrePaid, common.PostPaid
				ValidateFunc: validation.StringInSlice([]string{string(common.PrePaid), string(common.PostPaid)}, false),
			},
			"spot_strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      NoSpot,
				ValidateFunc: validation.StringInSlice([]string{"NoSpot", "SpotAsPriceGo", "SpotWithPriceLimit"}, false),
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
			// Computed values.
			"families": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"generation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_ids": {
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

func dataSourceAlicloudInstanceTypeFamiliesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	request := ecs.CreateDescribeInstanceTypeFamiliesRequest()
	request.RegionId = client.RegionId
	if v, ok := d.GetOk("generation"); ok {
		request.Generation = v.(string)
	}

	zones, err := ecsService.DescribeZones(d)
	families := make(map[string]map[string]string)
	for _, zone := range zones {
		for _, infos := range zone.AvailableResources.ResourcesInfo {
			for _, family := range infos.InstanceTypeFamilies.SupportedInstanceTypeFamily {
				if _, ok := families[family]; !ok {
					families[family] = make(map[string]string)
				}
				families[family][zone.ZoneId] = ""
			}
		}
	}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceTypeFamilies(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_instance_type_families", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	var instanceTypeFamilies []ecs.InstanceTypeFamily
	response, _ := raw.(*ecs.DescribeInstanceTypeFamiliesResponse)
	if response != nil {
		for _, family := range response.InstanceTypeFamilies.InstanceTypeFamily {
			if _, ok := families[family.InstanceTypeFamilyId]; !ok {
				continue
			}
			instanceTypeFamilies = append(instanceTypeFamilies, family)
		}
	}
	return instanceTypeFamiliesDescriptionAttributes(d, instanceTypeFamilies, families)
}

func instanceTypeFamiliesDescriptionAttributes(d *schema.ResourceData, typeFamilies []ecs.InstanceTypeFamily, families map[string]map[string]string) error {
	var ids []string
	var s []map[string]interface{}
	for _, f := range typeFamilies {

		mapping := map[string]interface{}{
			"id":         f.InstanceTypeFamilyId,
			"generation": f.Generation,
		}
		var zoneIds []string
		for zoneId := range families[f.InstanceTypeFamilyId] {
			zoneIds = append(zoneIds, zoneId)
		}
		sort.Strings(zoneIds)
		mapping["zone_ids"] = zoneIds

		ids = append(ids, f.InstanceTypeFamilyId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("families", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
