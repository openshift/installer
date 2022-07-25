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

func dataSourceAlicloudConfigAggregateConfigRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigAggregateConfigRulesRead,
		Schema: map[string]*schema.Schema{
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
			"aggregate_config_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"aggregator_id": {
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
			"risk_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{1, 2, 3}),
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "DELETING", "EVALUATING", "INACTIVE"}, false),
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
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregate_config_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compliance": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compliance_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"compliance_pack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_rule_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_rule_trigger_types": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_resource_ids_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"input_parameters": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"maximum_execution_frequency": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_ids_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_ids_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_types_scope": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"risk_level": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"source_identifier": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_owner": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_key_scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_value_scope": {
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

func dataSourceAlicloudConfigAggregateConfigRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAggregateConfigRules"
	request := make(map[string]interface{})
	if v, ok := d.GetOk("aggregate_config_rule_name"); ok {
		request["ConfigRuleName"] = v
	}
	request["AggregatorId"] = d.Get("aggregator_id")
	if v, ok := d.GetOk("risk_level"); ok {
		request["RiskLevel"] = v
	}
	if v, ok := d.GetOk("status"); ok {
		request["ConfigRuleState"] = v
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var aggregateConfigRuleNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		aggregateConfigRuleNameRegex = r
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
	conn, err := client.NewConfigClient()
	if err != nil {
		return WrapError(err)
	}
	for {
		runtime := util.RuntimeOptions{}
		runtime.SetAutoretry(true)
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("GET"), StringPointer("2020-09-07"), StringPointer("AK"), request, nil, &runtime)
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_aggregate_config_rules", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.ConfigRules.ConfigRuleList", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ConfigRules.ConfigRuleList", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if aggregateConfigRuleNameRegex != nil {
				if !aggregateConfigRuleNameRegex.MatchString(fmt.Sprint(item["ConfigRuleName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["ConfigRuleId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_id":                 fmt.Sprint(object["AccountId"]),
			"aggregate_config_rule_name": object["ConfigRuleName"],
			"compliance_pack_id":         object["CreateBy"].(map[string]interface{})["CompliancePackId"],
			"config_rule_arn":            object["ConfigRuleArn"],
			"id":                         fmt.Sprint(object["ConfigRuleId"]),
			"config_rule_id":             fmt.Sprint(object["ConfigRuleId"]),
			"description":                object["Description"],
			"risk_level":                 formatInt(object["RiskLevel"]),
			"source_identifier":          object["SourceIdentifier"],
			"source_owner":               object["SourceOwner"],
			"status":                     object["ConfigRuleState"],
		}
		mapping["compliance"] = []map[string]interface{}{{
			"compliance_type": object["Compliance"].(map[string]interface{})["ComplianceType"],
			"count":           object["Compliance"].(map[string]interface{})["Count"],
		}}
		mapping["aggregator_id"] = object["CreateBy"].(map[string]interface{})["AggregatorId"]
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, object["ConfigRuleName"])
			s = append(s, mapping)
			continue
		}

		configService := ConfigService{client}
		id := fmt.Sprint(request["AggregatorId"], ":", object["ConfigRuleId"])
		getResp, err := configService.DescribeConfigAggregateConfigRule(id)
		if err != nil {
			return WrapError(err)
		}
		mapping["config_rule_trigger_types"] = getResp["ConfigRuleTriggerTypes"]
		mapping["exclude_resource_ids_scope"] = getResp["ExcludeResourceIdsScope"]
		mapping["input_parameters"] = getResp["InputParameters"]
		mapping["maximum_execution_frequency"] = getResp["MaximumExecutionFrequency"]
		mapping["modified_timestamp"] = fmt.Sprint(getResp["ModifiedTimestamp"])
		mapping["region_ids_scope"] = getResp["RegionIdsScope"]
		mapping["resource_group_ids_scope"] = getResp["ResourceGroupIdsScope"]
		mapping["resource_types_scope"] = getResp["Scope"].(map[string]interface{})["ComplianceResourceTypes"]
		mapping["tag_key_scope"] = getResp["TagKeyScope"]
		mapping["tag_value_scope"] = getResp["TagValueScope"]
		if v := getResp["Source"].(map[string]interface{})["SourceDetails"].([]interface{}); len(v) > 0 {
			mapping["event_source"] = v[0].(map[string]interface{})["EventSource"]
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["ConfigRuleName"])
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
