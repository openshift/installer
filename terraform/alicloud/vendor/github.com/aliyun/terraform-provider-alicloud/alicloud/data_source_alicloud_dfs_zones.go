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

func dataSourceAlicloudDfsZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDfsZonesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
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
						"options": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol_type": {
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

func dataSourceAlicloudDfsZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "GetRegion"
	request := make(map[string]interface{})
	request["InputRegionId"] = client.RegionId

	var response map[string]interface{}
	conn, err := client.NewAlidfsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dfs_zones", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.AvailableZones", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.AvailableZones", response)
	}

	s := make([]map[string]interface{}, 0)
	for _, object := range resp.([]interface{}) {
		if zoneMap, ok := object.(map[string]interface{}); ok {

			mapping := map[string]interface{}{
				"zone_id": fmt.Sprint(zoneMap["ZoneId"]),
			}

			options := make([]interface{}, 0)
			for _, optionArg := range zoneMap["Options"].([]interface{}) {
				optionMap := optionArg.(map[string]interface{})
				singleObj := make(map[string]interface{}, 0)
				singleObj["storage_type"] = fmt.Sprint(optionMap["StorageType"])
				singleObj["protocol_type"] = fmt.Sprint(optionMap["ProtocolType"])
				options = append(options, singleObj)
			}
			mapping["options"] = options

			s = append(s, mapping)
		}
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
