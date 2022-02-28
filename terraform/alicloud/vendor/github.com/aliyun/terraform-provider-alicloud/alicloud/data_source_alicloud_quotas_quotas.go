package alicloud

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudQuotasQuotas() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudQuotasQuotasRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
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
			"group_code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"sort_field": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"TIME", "TOTAL", "RESERVED"}, false),
			},
			"sort_order": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Ascending", "Descending"}, false),
			},
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"quotas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"applicable_range": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"applicable_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"consumable": {
							Type:     schema.TypeBool,
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
						"quota_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"quota_unit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_quota": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"total_usage": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"unadjustable_detail": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudQuotasQuotasRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListProductQuotas"
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
	if v, ok := d.GetOk("group_code"); ok {
		request["GroupCode"] = v
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
	if v, ok := d.GetOk("sort_field"); ok {
		request["SortField"] = v
	}
	if v, ok := d.GetOk("sort_order"); ok {
		request["SortOrder"] = v
	}
	request["MaxResults"] = PageSizeLarge
	var objects []map[string]interface{}
	var quotaNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		quotaNameRegex = r
	}
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_quotas_quotas", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.Quotas", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Quotas", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if quotaNameRegex != nil {
				if !quotaNameRegex.MatchString(fmt.Sprint(item["QuotaName"])) {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                  fmt.Sprint(object["ProductCode"], ":", object["QuotaActionCode"]),
			"adjustable":          object["Adjustable"],
			"applicable_range":    object["ApplicableRange"],
			"applicable_type":     object["ApplicableType"],
			"consumable":          object["Consumable"],
			"quota_action_code":   object["QuotaActionCode"],
			"quota_description":   object["QuotaDescription"],
			"quota_name":          object["QuotaName"],
			"quota_type":          object["QuotaType"],
			"quota_unit":          object["QuotaUnit"],
			"total_quota":         object["TotalQuota"],
			"total_usage":         object["TotalUsage"],
			"unadjustable_detail": object["UnadjustableDetail"],
		}
		ids = append(ids, fmt.Sprint(object["ProductCode"], ":", object["QuotaActionCode"]))
		names = append(names, object["QuotaName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("quotas", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
