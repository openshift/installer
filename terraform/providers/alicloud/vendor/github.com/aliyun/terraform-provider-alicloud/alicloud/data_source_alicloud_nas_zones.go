package alicloud

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudNasZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudNasZonesRead,
		Schema: map[string]*schema.Schema{
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"standard", "extreme", "cpfs"}, false),
			},
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
						"instance_types": {
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

func dataSourceAlicloudNasZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeZones"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewNasClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-26"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_nas_zones", action, AlibabaCloudSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.Zones.Zone", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Zones.Zone", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		objects = append(objects, item)
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"zone_id": object["ZoneId"],
		}
		InstanceTypes := make([]map[string]interface{}, 0)
		if InstanceTypesList, ok := object["InstanceTypes"]; ok {
			if data, ok := InstanceTypesList.(map[string]interface{}); ok {
				for _, v := range data {
					if m1, ok := v.([]interface{}); ok {
						for _, vv := range m1 {
							if res, ok := vv.(map[string]interface{}); ok {
								temp1 := map[string]interface{}{
									"storage_type":  res["StorageType"],
									"protocol_type": res["ProtocolType"],
								}
								InstanceTypes = append(InstanceTypes, temp1)
							}
						}
					}
				}
			}
		}
		mapping["instance_types"] = InstanceTypes
		s = append(s, mapping)
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
