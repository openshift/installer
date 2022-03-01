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

func dataSourceAlicloudActiontrailHistoryDeliveryJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudActiontrailHistoryDeliveryJobsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 3}),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"history_delivery_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"home_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"trail_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_status": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAlicloudActiontrailHistoryDeliveryJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDeliveryHistoryJobs"
	request := make(map[string]interface{})
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
	status, statusOk := d.GetOkExists("status")
	var response map[string]interface{}
	conn, err := client.NewActiontrailClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-07-06"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_actiontrail_history_delivery_jobs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DeliveryHistoryJobs", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DeliveryHistoryJobs", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["JobId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(int) != formatInt(item["JobStatus"]) {
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
			"create_time":             object["CreatedTime"],
			"end_time":                object["EndTime"],
			"id":                      fmt.Sprint(object["JobId"]),
			"history_delivery_job_id": fmt.Sprint(object["JobId"]),
			"home_region":             object["HomeRegion"],
			"start_time":              object["StartTime"],
			"status":                  formatInt(object["JobStatus"]),
			"trail_name":              object["TrailName"],
			"updated_time":            object["UpdatedTime"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["JobId"])
		actiontrailService := ActiontrailService{client}
		getResp, err := actiontrailService.DescribeActiontrailHistoryDeliveryJob(id)
		if err != nil {
			return WrapError(err)
		}
		jobStatus := make([]map[string]interface{}, 0)
		if jobStatusList, ok := getResp["Status"]; ok {
			jobStatusListArg, _ := jobStatusList.([]interface{})
			for _, v := range jobStatusListArg {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"region": m1["Region"],
						"status": formatInt(m1["Status"]),
					}
					jobStatus = append(jobStatus, temp1)
				}
			}
		}
		mapping["job_status"] = jobStatus
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("jobs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
