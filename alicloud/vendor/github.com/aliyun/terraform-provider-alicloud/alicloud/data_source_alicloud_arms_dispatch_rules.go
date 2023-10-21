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

func dataSourceAlicloudArmsDispatchRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudArmsDispatchRulesRead,
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
			"dispatch_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dispatch_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group_id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_wait_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"group_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"grouping_fields": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"repeat_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"dispatch_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"label_match_expression_grid": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label_match_expression_groups": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"label_match_expressions": {
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
															"operator": {
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
							},
						},
						"notify_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"notify_objects": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"notify_object_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"notify_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"notify_channels": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"dispatch_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudArmsDispatchRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListDispatchRule"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("dispatch_rule_name"); ok {
		request["Name"] = v
	}
	request["RegionId"] = client.RegionId
	var objects []map[string]interface{}
	var dispatchRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		dispatchRuleNameRegex = r
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
	var response map[string]interface{}
	conn, err := client.NewArmsClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-08"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_arms_dispatch_rules", action, AlibabaCloudSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.DispatchRules", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.DispatchRules", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if dispatchRuleNameRegex != nil && !dispatchRuleNameRegex.MatchString(fmt.Sprint(item["Name"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["RuleId"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":                 fmt.Sprint(object["RuleId"]),
			"dispatch_rule_id":   fmt.Sprint(object["RuleId"]),
			"dispatch_rule_name": object["Name"],
			"status":             object["State"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["Name"])
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			s = append(s, mapping)
			continue
		}
		id := fmt.Sprint(object["RuleId"])
		armsService := ArmsService{client}
		getResp, err := armsService.DescribeArmsDispatchRule(id)
		if err != nil {
			return WrapError(err)
		}
		if groupRulesList, ok := getResp["GroupRules"]; ok && groupRulesList != nil {
			groupRulesMaps := make([]map[string]interface{}, 0)
			for _, groupRulesListItem := range groupRulesList.([]interface{}) {
				if groupRulesItemMap, ok := groupRulesListItem.(map[string]interface{}); ok {
					groupRulesMap := make(map[string]interface{}, 0)
					groupRulesMap["group_interval"] = groupRulesItemMap["GroupInterval"]
					groupRulesMap["group_wait_time"] = groupRulesItemMap["GroupWaitTime"]
					groupRulesMap["group_id"] = groupRulesItemMap["GroupId"]
					groupRulesMap["grouping_fields"] = groupRulesItemMap["GroupingFields"]
					groupRulesMap["repeat_interval"] = groupRulesItemMap["RepeatInterval"]
					groupRulesMaps = append(groupRulesMaps, groupRulesMap)
				}
			}
			mapping["group_rules"] = groupRulesMaps
		}
		if labelMatchExpressionGrid, ok := getResp["LabelMatchExpressionGrid"]; ok && labelMatchExpressionGrid != nil {
			labelMatchExpressionGridMaps := make([]map[string]interface{}, 0)
			labelMatchExpressionGridMap := make(map[string]interface{})

			labelMatchExpressionGroupsMaps := make([]map[string]interface{}, 0)
			if v, ok := labelMatchExpressionGrid.(map[string]interface{})["LabelMatchExpressionGroups"]; ok && v != nil {
				for _, labelMatchExpressionGroups := range v.([]interface{}) {
					labelMatchExpressionGroupsMap := make(map[string]interface{})
					labelMatchExpressionsMaps := make([]map[string]interface{}, 0)
					if v, ok := labelMatchExpressionGroups.(map[string]interface{})["LabelMatchExpressions"]; ok && v != nil {
						for _, labelMatchExpressions := range v.([]interface{}) {
							labelMatchExpressionsArg := labelMatchExpressions.(map[string]interface{})
							labelMatchExpressionsMap := make(map[string]interface{}, 0)
							labelMatchExpressionsMap["operator"] = labelMatchExpressionsArg["Operator"]
							labelMatchExpressionsMap["key"] = labelMatchExpressionsArg["Key"]
							labelMatchExpressionsMap["value"] = labelMatchExpressionsArg["Value"]
							labelMatchExpressionsMaps = append(labelMatchExpressionsMaps, labelMatchExpressionsMap)
						}
					}
					labelMatchExpressionGroupsMap["label_match_expressions"] = labelMatchExpressionsMaps
					labelMatchExpressionGroupsMaps = append(labelMatchExpressionGroupsMaps, labelMatchExpressionGroupsMap)
				}
			}
			labelMatchExpressionGridMap["label_match_expression_groups"] = labelMatchExpressionGroupsMaps
			labelMatchExpressionGridMaps = append(labelMatchExpressionGridMaps, labelMatchExpressionGridMap)
			mapping["label_match_expression_grid"] = labelMatchExpressionGridMaps
		}

		if notifyRulesList, ok := getResp["NotifyRules"]; ok && notifyRulesList != nil {
			notifyRulesMaps := make([]map[string]interface{}, 0)
			for _, notifyRulesListItem := range notifyRulesList.([]interface{}) {
				if notifyRulesItemMap, ok := notifyRulesListItem.(map[string]interface{}); ok {
					notifyRulesMap := make(map[string]interface{}, 0)
					notifyObjectsMaps := make([]map[string]interface{}, 0)
					for _, notifyObjects := range notifyRulesItemMap["NotifyObjects"].([]interface{}) {
						notifyObjectsArg := notifyObjects.(map[string]interface{})
						notifyObjectsMap := make(map[string]interface{}, 0)
						notifyObjectsMap["notify_type"] = convertArmsDispatchRuleNotifyTypeResponse(notifyObjectsArg["NotifyType"])
						notifyObjectsMap["notify_object_id"] = notifyObjectsArg["NotifyObjectId"]
						notifyObjectsMap["name"] = notifyObjectsArg["Name"]
						notifyObjectsMaps = append(notifyObjectsMaps, notifyObjectsMap)
					}
					notifyRulesMap["notify_objects"] = notifyObjectsMaps
					notifyRulesMap["notify_channels"] = notifyRulesItemMap["NotifyChannels"]
					notifyRulesMaps = append(notifyRulesMaps, notifyRulesMap)
				}
			}
			mapping["notify_rules"] = notifyRulesMaps
		}
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
