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

func dataSourceAlicloudConfigAggregateCompliancePacks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudConfigAggregateCompliancePacksRead,
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
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ACTIVE", "CREATING", "INACTIVE"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"packs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregate_compliance_pack_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"aggregator_compliance_pack_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compliance_pack_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"config_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"config_rule_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"config_rule_parameters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"parameter_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"parameter_value": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"required": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"managed_rule_identifier": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"risk_level": {
							Type:     schema.TypeInt,
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

func dataSourceAlicloudConfigAggregateCompliancePacksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	action := "ListAggregateCompliancePacks"
	request := make(map[string]interface{})
	request["AggregatorId"] = d.Get("aggregator_id")
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	var objects []map[string]interface{}
	var aggregateCompliancePackNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return WrapError(err)
		}
		aggregateCompliancePackNameRegex = r
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
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_config_aggregate_compliance_packs", action, AlibabaCloudSdkGoERROR)
		}
		resp, err := jsonpath.Get("$.CompliancePacksResult.CompliancePacks", response)
		if err != nil {
			return WrapErrorf(err, FailedGetAttributeMsg, action, "$.CompliancePacksResult.CompliancePacks", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if aggregateCompliancePackNameRegex != nil {
				if !aggregateCompliancePackNameRegex.MatchString(fmt.Sprint(item["CompliancePackName"])) {
					continue
				}
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["CompliancePackId"])]; !ok {
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
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"account_id":                     fmt.Sprint(object["AccountId"]),
			"aggregate_compliance_pack_name": object["CompliancePackName"],
			"id":                             fmt.Sprint(object["CompliancePackId"]),
			"aggregator_compliance_pack_id":  fmt.Sprint(object["CompliancePackId"]),
			"compliance_pack_template_id":    object["CompliancePackTemplateId"],
			"description":                    object["Description"],
			"risk_level":                     formatInt(object["RiskLevel"]),
			"status":                         object["Status"],
		}
		if detailedEnabled := d.Get("enable_details"); !detailedEnabled.(bool) {
			ids = append(ids, fmt.Sprint(mapping["id"]))
			names = append(names, object["CompliancePackName"])
			s = append(s, mapping)
			continue
		}

		configService := ConfigService{client}
		id := fmt.Sprint(object["AggregatorId"], ":", object["CompliancePackId"])
		getResp, err := configService.DescribeConfigAggregateCompliancePack(id)
		if err != nil {
			return WrapError(err)
		}

		configRules := make([]map[string]interface{}, 0)
		if configRulesList, ok := getResp["ConfigRules"].([]interface{}); ok {
			for _, v := range configRulesList {
				if m1, ok := v.(map[string]interface{}); ok {
					temp1 := map[string]interface{}{
						"config_rule_id":          m1["ConfigRuleId"],
						"managed_rule_identifier": m1["ManagedRuleIdentifier"],
					}
					if m1["ConfigRuleParameters"] != nil {
						configRuleParametersMaps := make([]map[string]interface{}, 0)
						for _, configRuleParametersValue := range m1["ConfigRuleParameters"].([]interface{}) {
							configRuleParameters := configRuleParametersValue.(map[string]interface{})
							configRuleParametersMap := map[string]interface{}{
								"parameter_name":  configRuleParameters["ParameterName"],
								"parameter_value": configRuleParameters["ParameterValue"],
								"required":        configRuleParameters["Required"],
							}
							configRuleParametersMaps = append(configRuleParametersMaps, configRuleParametersMap)
						}
						temp1["config_rule_parameters"] = configRuleParametersMaps
					}
					configRules = append(configRules, temp1)
				}
			}
		}
		mapping["config_rules"] = configRules
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["CompliancePackName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	if err := d.Set("packs", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
