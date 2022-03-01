package alicloud

import (
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudRegionsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"current": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
			//Computed value
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ecs.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	currentRegion := client.RegionId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_regions", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	resp, _ := raw.(*ecs.DescribeRegionsResponse)
	if resp == nil || len(resp.Regions.Region) == 0 {
		return fmt.Errorf("no matching regions found")
	}
	name, nameOk := d.GetOk("name")
	current := false
	if v, exist := d.GetOk("current"); exist {
		current = v.(bool)
	}
	var filterRegions []ecs.Region
	for _, region := range resp.Regions.Region {
		if current {
			if nameOk && name.(string) != currentRegion {
				return fmt.Errorf("name doesn't match current region: %#v, please input again.", currentRegion)
			}
			if region.RegionId == currentRegion {
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		if nameOk {
			if name.(string) == region.RegionId {
				filterRegions = append(filterRegions, region)
				break
			}
			continue
		}
		filterRegions = append(filterRegions, region)
	}
	if len(filterRegions) < 1 {
		return fmt.Errorf("Your query region returned no results. Please change your search criteria and try again.")
	}

	return regionsDescriptionAttributes(d, filterRegions)
}

func regionsDescriptionAttributes(d *schema.ResourceData, regions []ecs.Region) error {
	var ids []string
	var s []map[string]interface{}
	for _, region := range regions {
		mapping := map[string]interface{}{
			"id":         region.RegionId,
			"region_id":  region.RegionId,
			"local_name": region.LocalName,
		}

		log.Printf("[DEBUG] alicloud_regions - adding region mapping: %v", mapping)
		ids = append(ids, string(region.RegionId))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("regions", s); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
