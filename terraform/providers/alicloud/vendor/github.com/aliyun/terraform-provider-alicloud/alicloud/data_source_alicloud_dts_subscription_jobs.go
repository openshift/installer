package alicloud

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudDtsSubscriptionJobs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudDtsSubscriptionJobsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Abnormal", "Downgrade", "Locked", "Normal", "NotConfigured", "NotStarted", "PreCheckPass", "PrecheckFailed", "Prechecking", "Retrying", "Starting", "Upgrade"}, false),
			},
			"tags": tagsSchema(),
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"checkpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_list": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dts_job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_database_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_engine_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_oracle_sid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_owner_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_endpoint_user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_data_type_ddl": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"subscription_data_type_dml": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"subscription_host": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"private_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"public_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vpc_host": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"subscription_instance_network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_instance_vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_instance_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
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

func dataSourceAlicloudDtsSubscriptionJobsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "DescribeDtsJobs"
	request := make(map[string]interface{})
	request["JobType"] = "SUBSCRIBE"
	request["RegionId"] = client.RegionId
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var jobNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		jobNameRegex = r
	}

	idsMap := make(map[string]string)
	names := make([]interface{}, 0)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	status, statusOk := d.GetOk("status")
	tagsMap := make(map[string]interface{})
	if v, ok := d.GetOk("tags"); ok && len(v.(map[string]interface{})) > 0 {
		tagsMap = v.(map[string]interface{})
	}
	var response map[string]interface{}
	conn, err := client.NewDtsClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-01-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_dts_subscription_jobs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.DtsJobList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DtsJobList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if item["Status"] == "Starting" {
				item["Status"] = "Normal"
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DtsJobId"])]; !ok {
					continue
				}
			}
			if len(tagsMap) > 0 {
				if len(item["Tags"].(map[string]interface{})["Tag"].([]interface{})) != len(tagsMap) {
					continue
				}
				match := true
				for _, tag := range item["Tags"].([]interface{}) {
					if v, ok := tagsMap[tag.(map[string]interface{})["TagKey"].(string)]; !ok || v.(string) != tag.(map[string]interface{})["TagValue"].(string) {
						match = false
						break
					}
				}
				if !match {
					continue
				}
			}
			if jobNameRegex != nil && !jobNameRegex.MatchString(fmt.Sprint(item["DtsJobName"])) {
				continue
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
			"checkpoint":                    object["Checkpoint"],
			"create_time":                   object["CreateTime"],
			"db_list":                       object["DbObject"],
			"dts_instance_id":               object["DtsInstanceID"],
			"id":                            fmt.Sprint(object["DtsJobId"]),
			"dts_job_id":                    fmt.Sprint(object["DtsJobId"]),
			"dts_job_name":                  object["DtsJobName"],
			"expire_time":                   object["ExpireTime"],
			"payment_type":                  convertDtsPaymentTypeResponse(object["PayType"]),
			"source_endpoint_database_name": object["SourceEndpoint"].(map[string]interface{})["DatabaseName"],
			"source_endpoint_engine_name":   object["SourceEndpoint"].(map[string]interface{})["EngineName"],
			"source_endpoint_ip":            object["SourceEndpoint"].(map[string]interface{})["Ip"],
			"source_endpoint_instance_id":   object["SourceEndpoint"].(map[string]interface{})["InstanceID"],
			"source_endpoint_instance_type": object["SourceEndpoint"].(map[string]interface{})["InstanceType"],
			"source_endpoint_oracle_sid":    object["SourceEndpoint"].(map[string]interface{})["OracleSID"],
			"source_endpoint_port":          object["SourceEndpoint"].(map[string]interface{})["Port"],
			"source_endpoint_region":        object["SourceEndpoint"].(map[string]interface{})["Region"],
			"source_endpoint_user_name":     object["SourceEndpoint"].(map[string]interface{})["UserName"],
			"status":                        object["Status"],
		}

		var jsonData map[string]interface{}
		json.Unmarshal([]byte(object["Reserved"].(string)), &jsonData)
		if v, ok := jsonData["netType"].(string); ok {
			mapping["subscription_instance_network_type"] = strings.ToLower(v)
		}
		mapping["subscription_instance_vpc_id"] = jsonData["vpcId"]
		mapping["subscription_instance_vswitch_id"] = jsonData["vswitchId"]

		tags := make(map[string]interface{})
		t, _ := jsonpath.Get("$.TagList", object)
		if t != nil {
			for _, t := range t.([]interface{}) {
				key := t.(map[string]interface{})["TagKey"].(string)
				value := t.(map[string]interface{})["TagValue"].(string)
				if !ignoredTags(key, value) {
					tags[key] = value
				}
			}
		}
		mapping["tags"] = tags
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["DtsJobId"])
		dtsService := DtsService{client}
		getResp, err := dtsService.DescribeDtsSubscriptionJob(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["source_endpoint_owner_id"] = getResp["SourceEndpoint"].(map[string]interface{})["AliyunUid"]
		mapping["source_endpoint_role"] = getResp["SourceEndpoint"].(map[string]interface{})["RoleName"]
		mapping["subscription_data_type_ddl"] = getResp["SubscriptionDataType"].(map[string]interface{})["Ddl"]
		mapping["subscription_data_type_dml"] = getResp["SubscriptionDataType"].(map[string]interface{})["Dml"]
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
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
