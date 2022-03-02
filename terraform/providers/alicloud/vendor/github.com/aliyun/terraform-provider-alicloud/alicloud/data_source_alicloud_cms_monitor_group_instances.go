package alicloud

import (
	"strings"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceAlicloudCmsMonitorGroupInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudCmsMonitorGroupInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				MaxItems: 1,
				MinItems: 1,
			},
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
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

func dataSourceAlicloudCmsMonitorGroupInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeMonitorGroupInstances"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			request["GroupId"] = vv.(string)
		}
	}
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_cms_monitor_group_instances", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$", response)
		}
		objects = append(objects, resp.(map[string]interface{}))
		if len(resp.(map[string]interface{})["Resources"].(map[string]interface{})["Resource"].([]interface{})) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{}

		resourceMap := make([]map[string]interface{}, 0)
		if resourceMapList, ok := object["Resources"].(map[string]interface{})["Resource"].([]interface{}); ok {
			for _, v := range resourceMapList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"category":      strings.ToLower(m1["Category"].(string)),
						"instance_id":   m1["InstanceId"],
						"instance_name": m1["InstanceName"],
						"region_id":     m1["RegionId"],
					}
					resourceMap = append(resourceMap, temp1)
				}
			}
		}
		mapping["instances"] = resourceMap
		ids = append(ids, request["GroupId"].(string))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("instances", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
