package alicloud

import (
	"strings"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudNasProtocols() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasProtocolsRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Capacity",
					"Performance",
				}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocols": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlicloudNasProtocolsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DescribeZones"
	var response map[string]interface{}
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_protocols", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)
	resp, err := jsonpath.Get("$.Zones.Zone", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Zones.Zone", response)
	}
	var nasProtocol [][]interface{}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if v, ok := d.GetOk("zone_id"); ok && v.(string) != "" && item["ZoneId"].(string) != v.(string) {
			continue
		}
		if v, ok := d.GetOk("type"); ok {
			if Trim(v.(string)) == "Performance" {
				protocol, _ := item["Performance"].(map[string]interface{})["Protocol"].([]interface{})
				if len(protocol) == 0 {
					continue
				} else {
					nasProtocol = append(nasProtocol, protocol)
				}
			}
			if Trim(v.(string)) == "Capacity" {
				protocol, _ := item["Capacity"].(map[string]interface{})["Protocol"].([]interface{})
				if len(protocol) == 0 {
					continue
				} else {
					nasProtocol = append(nasProtocol, protocol)
				}
			}
		}
	}
	var s []string
	var ids []string
	for _, val := range nasProtocol {
		for _, protocol := range val {
			s = append(s, strings.ToUpper(protocol.(string)))
			ids = append(ids, protocol.(string))
		}
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("protocols", s); err != nil {
		return WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
