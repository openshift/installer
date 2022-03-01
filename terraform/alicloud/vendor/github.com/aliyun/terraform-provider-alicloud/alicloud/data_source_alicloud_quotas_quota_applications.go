package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudQuotasQuotaApplications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudQuotasQuotaApplicationsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"key_word": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"product_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_action_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"quota_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CommonQuota", "FlowControl"}, false),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Agree", "Disagree", "Process"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"approve_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"audit_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desire_value": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"dimensions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
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
						"effective_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notice_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"product_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_action_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
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

func dataSourceAlicloudQuotasQuotaApplicationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListQuotaApplications"
	request := make(map[string]interface{})
	request["SourceIp"] = client.SourceIp
	if v, ok := d.GetOk("dimensions"); ok {
		dimensionsMaps := make([]map[string]interface{}, 0)
		for _, dimensions := range v.([]interface{}) {
			dimensionsMap := make(map[string]interface{})
			dimensionsArg := dimensions.(map[string]interface{})
			dimensionsMap["Key"] = dimensionsArg["key"]
			dimensionsMap["Value"] = dimensionsArg["value"]
			dimensionsMaps = append(dimensionsMaps, dimensionsMap)
		}
		request["Dimensions"] = dimensionsMaps

	}
	if v, ok := d.GetOk("key_word"); ok {
		request["KeyWord"] = v
	}
	request["ProductCode"] = d.Get("product_code")
	if v, ok := d.GetOk("quota_action_code"); ok {
		request["QuotaActionCode"] = v
	}
	if v, ok := d.GetOk("quota_category"); ok {
		request["QuotaCategory"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["Status"] = v
	}
	request["MaxResults"] = PageSizeLarge
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
	conn, err := client.NewQuotasClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-10"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_quotas_quota_applications", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.QuotaApplications", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.QuotaApplications", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ApplicationId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if nextToken, ok := response["NextToken"].(string); ok && nextToken != "" {
			request["NextToken"] = nextToken
		} else {
			break
		}
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                fmt.Sprint(object["ApplicationId"]),
			"application_id":    fmt.Sprint(object["ApplicationId"]),
			"approve_value":     object["ApproveValue"],
			"audit_reason":      object["AuditReason"],
			"desire_value":      object["DesireValue"],
			"effective_time":    object["EffectiveTime"],
			"expire_time":       object["ExpireTime"],
			"notice_type":       formatInt(object["NoticeType"]),
			"product_code":      object["ProductCode"],
			"quota_action_code": object["QuotaActionCode"],
			"quota_description": object["QuotaDescription"],
			"quota_name":        object["QuotaName"],
			"quota_unit":        object["QuotaUnit"],
			"reason":            object["Reason"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(object["ApplicationId"]))
			s = append(s, mapping)
			continue
		}

		quotasService := QuotasService{client}
		id := fmt.Sprint(object["ApplicationId"])
		getResp, err := quotasService.DescribeQuotasQuotaApplication(id)
		if err != nil {
			return WrapError(err)
		}
		if statusOk && status != "" && status != getResp["Status"].(string) {
			continue
		}

		dimensionList := make([]map[string]interface{}, 0)
		if dimension, ok := getResp["Dimension"]; ok {
			for k, v := range dimension.(map[string]interface{}) {
				dimensionMap := make(map[string]interface{})
				dimensionMap["key"] = k
				dimensionMap["value"] = v
				dimensionList = append(dimensionList, dimensionMap)
			}
		}

		mapping["dimensions"] = dimensionList
		mapping["status"] = getResp["Status"]
		ids = append(ids, fmt.Sprint(object["ApplicationId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("applications", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
