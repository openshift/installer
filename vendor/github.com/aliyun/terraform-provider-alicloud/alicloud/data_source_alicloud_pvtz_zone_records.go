package alicloud

import (
	"fmt"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudPvtzZoneRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudPvtzZoneRecordsRead,
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"lang": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"search_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"DISABLE", "ENABLE"}, false),
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"user_client_ip": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"record_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"remark": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ttl": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudPvtzZoneRecordsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeZoneRecords"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("keyword"); ok {
		request["Keyword"] = v
	}
	if v, ok := d.GetOk("lang"); ok {
		request["Lang"] = v
	}
	if v, ok := d.GetOk("search_mode"); ok {
		request["SearchMode"] = v
	}
	if v, ok := d.GetOk("tag"); ok {
		request["Tag"] = v
	}
	if v, ok := d.GetOk("user_client_ip"); ok {
		request["UserClientIp"] = v
	}
	request["ZoneId"] = d.Get("zone_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	var response map[string]interface{}
	conn, err := client.NewPvtzClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 5*time.Second)
		err := resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-01-01"), StringPointer("AK"), nil, request, &runtime)
			if err != nil {
				if IsExpectedErrors(err, []string{"System.Busy", "ServiceUnavailable", "Throttling.User"}) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_pvtz_zone_records", action, AlibabaCloudSdkGoERROR)
		}

		resp, err := jsonpath.Get("$.Records.Record", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Records.Record", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RecordId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":              fmt.Sprint(object["RecordId"], ":", request["ZoneId"]),
			"priority":        formatInt(object["Priority"]),
			"record_id":       fmt.Sprint(object["RecordId"]),
			"remark":          object["Remark"],
			"rr":              object["Rr"],
			"resource_record": object["Rr"],
			"status":          object["Status"],
			"ttl":             formatInt(object["Ttl"]),
			"type":            object["Type"],
			"value":           object["Value"],
		}
		ids = append(ids, fmt.Sprint(object["RecordId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("records", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
