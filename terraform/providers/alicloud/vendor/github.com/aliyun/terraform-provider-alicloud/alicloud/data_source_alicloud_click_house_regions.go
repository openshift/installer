package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudClickHouseRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudClickHouseRegionsRead,
		Schema: map[string]*schema.Schema{
			"region_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"current"},
			},
			"current": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"region_id"},
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vpc_enabled": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudClickHouseRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeRegions"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var response map[string]interface{}
	update := false
	conn, err := client.NewClickhouseClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-11"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_click_house_regions", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.Regions.Region", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Accounts.Account", response)
	}
	result, _ := resp.([]interface{})
	if v, exist := d.GetOk("current"); exist {
		update = v.(bool)
	}
	for _, v := range result {
		obj := v.(map[string]interface{})
		if update {
			if client.RegionId != obj["RegionId"].(string) {
				continue
			}
		}
		if v, exist := d.GetOk("region_id"); exist {
			if v.(string) != obj["RegionId"].(string) {
				continue
			}
		}
		objects = append(objects, obj)
	}

	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"region_id": fmt.Sprint(object["RegionId"]),
		}

		zoneRes := make([]interface{}, 0)

		if _, exist := object["Zones"]; exist {
			for _, zonesmapslicSli := range object["Zones"].(map[string]interface{}) {
				for _, obj := range zonesmapslicSli.([]interface{}) {
					objMap := obj.(map[string]interface{})
					singleObj := make(map[string]interface{}, 0)
					singleObj["vpc_enabled"] = objMap["VpcEnabled"].(bool)
					singleObj["zone_id"] = fmt.Sprint(objMap["ZoneId"])
					zoneRes = append(zoneRes, singleObj)
				}
			}
		}

		mapping["zone_ids"] = zoneRes
		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("regions", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
