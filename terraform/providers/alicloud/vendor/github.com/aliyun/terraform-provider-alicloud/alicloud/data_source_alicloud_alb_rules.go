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

func dataSourceAlicloudAlbRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAlbRulesRead,
		Schema: map[string]*schema.Schema{
			"listener_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"load_balancer_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				ValidateFunc: validation.StringInSlice([]string{"Available", "Configuring", "Provisioning"}, false),
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
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"fixed_response_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"content": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"content_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"http_code": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"forward_group_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"server_group_tuples": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"server_group_id": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
									"insert_header_config": {
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
												"value_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"order": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"redirect_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"http_code": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"port": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"protocol": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"query": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"rewrite_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"path": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"query": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"rule_conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cookie_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"values": {
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
											},
										},
									},
									"header_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"host_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"method_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"path_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"query_string_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"values": {
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
											},
										},
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_id": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlicloudAlbRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListRules"
	request := make(map[string]interface{})
	if m, ok := d.GetOk("listener_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("ListenerIds.%d", k+1)] = v.(string)
		}
	}
	if m, ok := d.GetOk("load_balancer_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("LoadBalancerIds.%d", k+1)] = v.(string)
		}
	}
	if m, ok := d.GetOk("rule_ids"); ok {
		for k, v := range m.([]interface{}) {
			request[fmt.Sprintf("RuleIds.%d", k+1)] = v.(string)
		}
	}
	request["MaxResults"] = PageSizeLarge
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
	conn, err := client.NewAlbClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-06-16"), StringPointer("AK"), nil, request, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_alb_rules", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.Rules", response)
		if formatInt(response["TotalCount"]) != 0 && err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.Rules", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if ruleNameRegex != nil && !ruleNameRegex.MatchString(fmt.Sprint(item["RuleName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["RuleId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["RuleStatus"].(string) {
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
			"listener_id":      object["ListenerId"],
			"load_balancer_id": object["LoadBalancerId"],
			"priority":         formatInt(object["Priority"]),
			"id":               fmt.Sprint(object["RuleId"]),
			"rule_id":          fmt.Sprint(object["RuleId"]),
			"rule_name":        object["RuleName"],
			"status":           object["RuleStatus"],
		}

		ruleActionsMaps := make([]map[string]interface{}, 0)
		if ruleActionsList, ok := object["RuleActions"]; ok {
			for _, ruleActions := range ruleActionsList.([]interface{}) {
				ruleActionsArg := ruleActions.(map[string]interface{})
				ruleActionsMap := map[string]interface{}{}
				ruleActionsMap["type"] = ruleActionsArg["Type"]
				ruleActionsMap["order"] = formatInt(ruleActionsArg["Order"])

				if forwardGroupConfig, ok := ruleActionsArg["ForwardGroupConfig"]; ok {
					forwardGroupConfigArg := forwardGroupConfig.(map[string]interface{})
					if len(forwardGroupConfigArg) > 0 {
						serverGroupTuplesMaps := make([]map[string]interface{}, 0)
						if forwardGroupConfigArgs, ok := forwardGroupConfigArg["ServerGroupTuples"].([]interface{}); ok {
							for _, serverGroupTuples := range forwardGroupConfigArgs {
								serverGroupTuplesArg := serverGroupTuples.(map[string]interface{})
								serverGroupTuplesMap := map[string]interface{}{}
								serverGroupTuplesMap["server_group_id"] = serverGroupTuplesArg["ServerGroupId"]
								serverGroupTuplesMaps = append(serverGroupTuplesMaps, serverGroupTuplesMap)
							}
						}
						if len(serverGroupTuplesMaps) > 0 {
							forwardGroupConfigMaps := make([]map[string]interface{}, 0)
							forwardGroupConfigMap := map[string]interface{}{}
							forwardGroupConfigMap["server_group_tuples"] = serverGroupTuplesMaps
							forwardGroupConfigMaps = append(forwardGroupConfigMaps, forwardGroupConfigMap)
							ruleActionsMap["forward_group_config"] = forwardGroupConfigMaps
							ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
						}
					}
				}

				if fixedResponseConfig, ok := ruleActionsArg["FixedResponseConfig"]; ok {
					fixedResponseConfigArg := fixedResponseConfig.(map[string]interface{})
					if len(fixedResponseConfigArg) > 0 {
						fixedResponseConfigMaps := make([]map[string]interface{}, 0)
						fixedResponseConfigMap := make(map[string]interface{}, 0)
						fixedResponseConfigMap["content"] = fixedResponseConfigArg["Content"]
						fixedResponseConfigMap["content_type"] = fixedResponseConfigArg["ContentType"]
						fixedResponseConfigMap["http_code"] = fixedResponseConfigArg["HttpCode"]
						fixedResponseConfigMaps = append(fixedResponseConfigMaps, fixedResponseConfigMap)
						ruleActionsMap["fixed_response_config"] = fixedResponseConfigMaps
						ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
					}
				}

				if insertHeaderConfig, ok := ruleActionsArg["InsertHeaderConfig"]; ok {
					insertHeaderConfigArg := insertHeaderConfig.(map[string]interface{})
					if len(insertHeaderConfigArg) > 0 {
						insertHeaderConfigMaps := make([]map[string]interface{}, 0)
						insertHeaderConfigMap := make(map[string]interface{}, 0)
						insertHeaderConfigMap["key"] = insertHeaderConfigArg["Key"]
						insertHeaderConfigMap["value"] = insertHeaderConfigArg["Value"]
						insertHeaderConfigMap["value_type"] = insertHeaderConfigArg["ValueType"]
						insertHeaderConfigMaps = append(insertHeaderConfigMaps, insertHeaderConfigMap)
						ruleActionsMap["insert_header_config"] = insertHeaderConfigMaps
						ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
					}
				}

				if redirectConfig, ok := ruleActionsArg["RedirectConfig"]; ok {
					redirectConfigArg := redirectConfig.(map[string]interface{})
					if len(redirectConfigArg) > 0 {
						redirectConfigMaps := make([]map[string]interface{}, 0)
						redirectConfigMap := make(map[string]interface{}, 0)
						redirectConfigMap["host"] = redirectConfigArg["Host"]
						redirectConfigMap["http_code"] = redirectConfigArg["HttpCode"]
						redirectConfigMap["path"] = redirectConfigArg["Path"]
						redirectConfigMap["port"] = formatInt(redirectConfigArg["Port"])
						redirectConfigMap["protocol"] = redirectConfigArg["Protocol"]
						redirectConfigMap["query"] = redirectConfigArg["Query"]
						redirectConfigMaps = append(redirectConfigMaps, redirectConfigMap)
						ruleActionsMap["redirect_config"] = redirectConfigMaps
						ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
					}
				}

				if rewriteConfig, ok := ruleActionsArg["RewriteConfig"]; ok {
					rewriteConfigArg := rewriteConfig.(map[string]interface{})
					if len(rewriteConfigArg) > 0 {
						rewriteConfigMaps := make([]map[string]interface{}, 0)
						rewriteConfigMap := make(map[string]interface{}, 0)
						rewriteConfigMap["host"] = rewriteConfigArg["Host"]
						rewriteConfigMap["path"] = rewriteConfigArg["Path"]
						rewriteConfigMap["query"] = rewriteConfigArg["Query"]
						rewriteConfigMaps = append(rewriteConfigMaps, rewriteConfigMap)
						ruleActionsMap["rewrite_config"] = rewriteConfigMaps
						ruleActionsMaps = append(ruleActionsMaps, ruleActionsMap)
					}
				}
			}
		}
		mapping["rule_actions"] = ruleActionsMaps

		ruleConditionsMaps := make([]map[string]interface{}, 0)
		if ruleConditionsList, ok := object["RuleConditions"]; ok {
			for _, ruleConditions := range ruleConditionsList.([]interface{}) {
				ruleConditionsArg := ruleConditions.(map[string]interface{})
				ruleConditionsMap := map[string]interface{}{}
				ruleConditionsMap["type"] = ruleConditionsArg["Type"]

				if cookieConfig, ok := ruleConditionsArg["CookieConfig"]; ok {
					cookieConfigArg := cookieConfig.(map[string]interface{})
					if len(cookieConfigArg) > 0 {
						cookieConfigMaps := make([]map[string]interface{}, 0)
						valuesMaps := make([]map[string]interface{}, 0)
						for _, values := range cookieConfigArg["Values"].([]interface{}) {
							valuesArg := values.(map[string]interface{})
							valuesMap := map[string]interface{}{}
							valuesMap["key"] = valuesArg["Key"]
							valuesMap["value"] = valuesArg["Value"]
							valuesMaps = append(valuesMaps, valuesMap)
						}
						cookieConfigMap := map[string]interface{}{}
						cookieConfigMap["values"] = valuesMaps
						cookieConfigMaps = append(cookieConfigMaps, cookieConfigMap)
						ruleConditionsMap["cookie_config"] = cookieConfigMaps
						ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
					}
				}

				if headerConfig, ok := ruleConditionsArg["HeaderConfig"]; ok {
					headerConfigArg := headerConfig.(map[string]interface{})
					if len(headerConfigArg) > 0 {
						headerConfigMaps := make([]map[string]interface{}, 0)
						headerConfigMap := map[string]interface{}{}
						headerConfigMap["values"] = headerConfigArg["Values"].([]interface{})
						headerConfigMap["key"] = headerConfigArg["Key"]
						headerConfigMaps = append(headerConfigMaps, headerConfigMap)
						ruleConditionsMap["header_config"] = headerConfigMaps
						ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
					}
				}

				if queryStringConfig, ok := ruleConditionsArg["QueryStringConfig"]; ok {
					queryStringConfigArg := queryStringConfig.(map[string]interface{})
					if len(queryStringConfigArg) > 0 {
						queryStringConfigMaps := make([]map[string]interface{}, 0)
						queryStringValuesMaps := make([]map[string]interface{}, 0)
						for _, values := range queryStringConfigArg["Values"].([]interface{}) {
							valuesArg := values.(map[string]interface{})
							valuesMap := map[string]interface{}{}
							valuesMap["key"] = valuesArg["Key"]
							valuesMap["value"] = valuesArg["Value"]
							queryStringValuesMaps = append(queryStringValuesMaps, valuesMap)
						}
						queryStringConfigMap := map[string]interface{}{}
						queryStringConfigMap["values"] = queryStringValuesMaps
						queryStringConfigMaps = append(queryStringConfigMaps, queryStringConfigMap)
						ruleConditionsMap["query_string_config"] = queryStringConfigMaps
						ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
					}
				}

				if hostConfig, ok := ruleConditionsArg["HostConfig"]; ok {
					hostConfigArg := hostConfig.(map[string]interface{})
					if len(hostConfigArg) > 0 {
						hostConfigMaps := make([]map[string]interface{}, 0)
						hostConfigMap := map[string]interface{}{}
						hostConfigMap["values"] = hostConfigArg["Values"].([]interface{})
						hostConfigMaps = append(hostConfigMaps, hostConfigMap)
						ruleConditionsMap["host_config"] = hostConfigMaps
						ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
					}
				}

				if methodConfig, ok := ruleConditionsArg["MethodConfig"]; ok {
					methodConfigArg := methodConfig.(map[string]interface{})
					if len(methodConfigArg) > 0 {
						methodConfigMaps := make([]map[string]interface{}, 0)
						methodConfigMap := map[string]interface{}{}
						methodConfigMap["values"] = methodConfigArg["Values"].([]interface{})
						methodConfigMaps = append(methodConfigMaps, methodConfigMap)
						ruleConditionsMap["method_config"] = methodConfigMaps
						ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
					}
				}

				if pathConfig, ok := ruleConditionsArg["PathConfig"]; ok {
					pathConfigArg := pathConfig.(map[string]interface{})
					if len(pathConfigArg) > 0 {
						pathConfigMaps := make([]map[string]interface{}, 0)
						pathConfigMap := map[string]interface{}{}
						pathConfigMap["values"] = pathConfigArg["Values"].([]interface{})
						pathConfigMaps = append(pathConfigMaps, pathConfigMap)
						ruleConditionsMap["path_config"] = pathConfigMaps
						ruleConditionsMaps = append(ruleConditionsMaps, ruleConditionsMap)
					}
				}
			}
		}
		mapping["rule_conditions"] = ruleConditionsMaps

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
