package alicloud

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceAlicloudEssScalingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudEssScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"scaling_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.ValidateRegexp,
				ForceNew:     true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cooldown": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"adjustment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"adjustment_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_adjustment_magnitude": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scaling_rule_ari": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudEssScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := ess.CreateDescribeScalingRulesRequest()
	request.RegionId = client.RegionId
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

	if scalingGroupId, ok := d.GetOk("scaling_group_id"); ok && scalingGroupId.(string) != "" {
		request.ScalingGroupId = scalingGroupId.(string)
	}

	if ruleType, ok := d.GetOk("type"); ok && ruleType.(string) != "" {
		request.ScalingRuleType = ruleType.(string)
	}

	var allScalingRules []ess.ScalingRule

	for {
		raw, err := client.WithEssClient(func(essClient *ess.Client) (interface{}, error) {
			return essClient.DescribeScalingRules(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alicloud_ess_scalingrules", request.GetActionName(), AlibabaCloudSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*ess.DescribeScalingRulesResponse)
		if len(response.ScalingRules.ScalingRule) < 1 {
			break
		}

		allScalingRules = append(allScalingRules, response.ScalingRules.ScalingRule...)

		if len(response.ScalingRules.ScalingRule) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return err
		} else {
			request.PageNumber = page
		}
	}

	var filteredScalingRulesTemp = make([]ess.ScalingRule, 0)

	nameRegex, okNameRegex := d.GetOk("name_regex")
	idsMap := make(map[string]string)
	ids, okIds := d.GetOk("ids")
	if okIds {
		for _, i := range ids.([]interface{}) {
			if i == nil {
				continue
			}
			idsMap[i.(string)] = i.(string)
		}
	}

	if okNameRegex || okIds {
		for _, rule := range allScalingRules {
			if okNameRegex && nameRegex != "" {
				r, err := regexp.Compile(nameRegex.(string))
				if err != nil {
					return WrapError(err)
				}
				if r != nil && !r.MatchString(rule.ScalingRuleName) {
					continue
				}
			}
			if okIds && len(idsMap) > 0 {
				if _, ok := idsMap[rule.ScalingRuleId]; !ok {
					continue
				}
			}
			filteredScalingRulesTemp = append(filteredScalingRulesTemp, rule)
		}
	} else {
		filteredScalingRulesTemp = allScalingRules
	}
	return scalingRulesDescriptionAttribute(d, filteredScalingRulesTemp, meta)
}

func scalingRulesDescriptionAttribute(d *schema.ResourceData, scalingRules []ess.ScalingRule, meta interface{}) error {
	var ids []string
	var names []string
	var s = make([]map[string]interface{}, 0)
	for _, scalingRule := range scalingRules {
		mapping := map[string]interface{}{
			"id":                       scalingRule.ScalingRuleId,
			"scaling_group_id":         scalingRule.ScalingGroupId,
			"name":                     scalingRule.ScalingRuleName,
			"type":                     scalingRule.ScalingRuleType,
			"cooldown":                 scalingRule.Cooldown,
			"adjustment_type":          scalingRule.AdjustmentType,
			"adjustment_value":         scalingRule.AdjustmentValue,
			"min_adjustment_magnitude": scalingRule.MinAdjustmentMagnitude,
			"scaling_rule_ari":         scalingRule.ScalingRuleAri,
		}
		ids = append(ids, scalingRule.ScalingRuleId)
		names = append(names, scalingRule.ScalingRuleName)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("rules", s); err != nil {
		return WrapError(err)
	}

	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
