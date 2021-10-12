package alicloud

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEventBridgeRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEventBridgeRulesRead,
		Schema: map[string]*schema.Schema{
			"event_bus_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"rule_name_prefix": {
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
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_bus_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"filter_pattern": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
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

func dataSourceAlicloudEventBridgeRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListRules"
	request := make(map[string]interface{})
	request["Limit"] = PageSizeLarge
	request["EventBusName"] = d.Get("event_bus_name")
	if v, ok := d.GetOk("rule_name_prefix"); ok {
		request["RuleNamePrefix"] = v
	}
	var objects []map[string]interface{}
	var ruleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		ruleNameRegex = r
	}

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
	conn, err := client.NewEventbridgeClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-04-01"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_event_bridge_rules", action, AlibabaCloudSdkGoERROR)
		}
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("ListRules failed, response: %v", response))
		}
		resp, err := jsonpath.Get("$.Data.Rules", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Data.Rules", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if ruleNameRegex != nil {
				if !ruleNameRegex.MatchString(fmt.Sprint(item["RuleName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RuleName"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["Status"].(string) {
				continue
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
			"description":    object["Description"],
			"event_bus_name": object["EventBusName"],
			"filter_pattern": object["FilterPattern"],
			"id":             fmt.Sprint(object["RuleName"]),
			"rule_name":      fmt.Sprint(object["RuleName"]),
			"status":         object["Status"],
		}

		targets := make([]map[string]interface{}, 0)
		if targetsList, ok := object["Targets"].([]interface{}); ok {
			for _, v := range targetsList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"endpoint":  m1["Endpoint"],
						"target_id": m1["Id"],
						"type":      m1["Type"],
					}
					targets = append(targets, temp1)
				}
			}
		}
		mapping["targets"] = targets
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["RuleName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
