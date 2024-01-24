package alicloud

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudGaForwardingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudGaForwardingRulesRead,
		Schema: map[string]*schema.Schema{
			"accelerator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"active", "configuring"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"forwarding_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"forwarding_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarding_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forwarding_rule_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rule_conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rule_condition_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"path_config": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"values": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
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
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"rule_actions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"order": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_action_type": {
										Type:     schema.TypeString,
										Computed: true,
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
															"endpoint_group_id": {
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudGaForwardingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListForwardingRules"
	request := make(map[string]interface{})
	request["AcceleratorId"] = d.Get("accelerator_id")
	request["ListenerId"] = d.Get("listener_id")
	request["RegionId"] = client.RegionId
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
	conn, err := client.NewGaplusClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-11-20"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ga_forwarding_rules", action, AlibabaCloudSdkGoERROR)
		}
		addDebug(action, response, request)

		resp, err := jsonpath.Get("$.ForwardingRules", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ForwardingRules", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ForwardingRuleId"])]; !ok {
					continue
				}
			}
			if statusOk && status.(string) != "" && status.(string) != item["ForwardingRuleStatus"].(string) {
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
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		ruleConditions := make([]map[string]interface{}, 0)
		for _, v := range object["RuleConditions"].([]interface{}) {
			ruleCondition := map[string]interface{}{}
			v := v.(map[string]interface{})
			ruleCondition["rule_condition_type"] = v["RuleConditionType"]
			if v["PathConfig"].(map[string]interface{})["Values"] != nil {
				ruleCondition["path_config"] = []map[string]interface{}{
					{
						"values": v["PathConfig"].(map[string]interface{})["Values"],
					},
				}
			}
			if v["HostConfig"].(map[string]interface{})["Values"] != nil {
				ruleCondition["host_config"] = []map[string]interface{}{
					{
						"values": v["HostConfig"].(map[string]interface{})["Values"],
					},
				}
			}
			ruleConditions = append(ruleConditions, ruleCondition)
		}
		ruleActions := make([]map[string]interface{}, 0)
		for _, v := range object["RuleActions"].([]interface{}) {
			v := v.(map[string]interface{})
			ruleAction := map[string]interface{}{}
			ruleAction["order"] = v["Order"]
			ruleAction["rule_action_type"] = v["RuleActionType"]
			serverGroupTuples := make([]map[string]interface{}, 0)
			for _, serverGroupTuple := range v["ForwardGroupConfig"].(map[string]interface{})["ServerGroupTuples"].([]interface{}) {
				serverGroupTuples = append(serverGroupTuples, map[string]interface{}{
					"endpoint_group_id": serverGroupTuple.(map[string]interface{})["EndpointGroupId"],
				})
			}
			ruleAction["forward_group_config"] = []map[string]interface{}{
				{
					"server_group_tuples": serverGroupTuples,
				},
			}
			ruleActions = append(ruleActions, ruleAction)
		}
		mapping := map[string]interface{}{
			"priority":               object["Priority"],
			"forwarding_rule_id":     object["ForwardingRuleId"],
			"id":                     fmt.Sprint(d.Get("accelerator_id"), ":", object["ListenerId"], ":", object["ForwardingRuleId"]),
			"forwarding_rule_name":   object["ForwardingRuleName"],
			"forwarding_rule_status": object["ForwardingRuleStatus"],
			"rule_conditions":        ruleConditions,
			"rule_actions":           ruleActions,
			"listener_id":            object["ListenerId"],
		}
		ids = append(ids, fmt.Sprint(object["ForwardingRuleId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("forwarding_rules", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
