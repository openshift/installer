package alicloud

import (
	"sort"
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudEnhancedNatAvailableZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEnhancedNatAvailableZonesRead,

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
						"zone_id": {
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

func dataSourceAlicloudEnhancedNatAvailableZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "ListEnhanhcedNatGatewayAvailableZones"
	request := map[string]interface{}{
		"RegionId": client.RegionId,
	}

	conn, err := meta.(*connectivity.AliyunClient).NewVpcClient()
	if err != nil {
		return WrapError(err)
	}
	// If the API supports
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = resource.Retry(3*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2016-04-28"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"TaskConflict", "UnknownError", Throttling}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_enhanced_nat_available_zones", action, AlibabaCloudSdkGoERROR)
	}

	var s []map[string]interface{}
	var ids []string

	v, err := jsonpath.Get("$.Zones", response)
	if err != nil || len(v.([]interface{})) < 1 {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_enhanced_nat_available_zones", action, AlibabaCloudSdkGoERROR)
	}

	zones := v.([]interface{})
	sort.Slice(zones, func(i, j int) bool {
		return zones[i].(map[string]interface{})["ZoneId"].(string) < zones[j].(map[string]interface{})["ZoneId"].(string)
	})
	for _, val := range zones {
		value := val.(map[string]interface{})
		mapping := map[string]interface{}{
			"zone_id":    value["ZoneId"].(string),
			"local_name": value["LocalName"].(string),
		}
		s = append(s, mapping)
		ids = append(ids, value["LocalName"].(string))
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
