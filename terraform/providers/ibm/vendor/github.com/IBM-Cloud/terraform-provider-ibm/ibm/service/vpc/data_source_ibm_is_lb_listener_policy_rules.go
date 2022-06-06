// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMISLBListenerPolicyRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsLbListenerPolicyRulesRead,

		Schema: map[string]*schema.Schema{
			isLBListenerPolicyRuleLBID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The load balancer identifier.",
			},
			isLBListenerPolicyRuleListenerID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The listener identifier.",
			},
			isLBListenerPolicyRulePolicyID: &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy identifier.",
			},
			"rules": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The condition of the rule.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this rule was created.",
						},
						"field": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The field. This is applicable to `header`, `query`, and `body` rule types.If the rule type is `header`, this property is required.If the rule type is `query`, this is optional. If specified and the rule condition is not`matches_regex`, the value must be percent-encoded.If the rule type is `body`, this is optional.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The rule's unique identifier.",
						},
						"provisioning_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning status of this rule.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the rule.Body rules are applied to form-encoded request bodies using the `UTF-8` character set.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value to be matched for rule condition.If the rule type is `query` and the rule condition is not `matches_regex`, the value must be percent-encoded.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsLbListenerPolicyRulesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listLoadBalancerListenerPolicyRulesOptions := &vpcv1.ListLoadBalancerListenerPolicyRulesOptions{}
	listLoadBalancerListenerPolicyRulesOptions.SetLoadBalancerID(d.Get(isLBListenerPolicyRuleLBID).(string))
	listLoadBalancerListenerPolicyRulesOptions.SetListenerID(d.Get(isLBListenerPolicyRuleListenerID).(string))
	listLoadBalancerListenerPolicyRulesOptions.SetPolicyID(d.Get(isLBListenerPolicyRulePolicyID).(string))

	loadBalancerListenerPolicyRuleCollection, response, err := vpcClient.ListLoadBalancerListenerPolicyRulesWithContext(context, listLoadBalancerListenerPolicyRulesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListLoadBalancerListenerPolicyRulesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListLoadBalancerListenerPolicyRulesWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsLbListenerPolicyRulesID(d))

	if loadBalancerListenerPolicyRuleCollection.Rules != nil {
		err = d.Set("rules", dataSourceLoadBalancerListenerPolicyRuleCollectionFlattenRules(loadBalancerListenerPolicyRuleCollection.Rules))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rules %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsLbListenerPolicyRulesID returns a reasonable ID for the list.
func dataSourceIBMIsLbListenerPolicyRulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceLoadBalancerListenerPolicyRuleCollectionFlattenRules(result []vpcv1.LoadBalancerListenerPolicyRule) (rules []map[string]interface{}) {
	for _, rulesItem := range result {
		rules = append(rules, dataSourceLoadBalancerListenerPolicyRuleCollectionRulesToMap(rulesItem))
	}

	return rules
}

func dataSourceLoadBalancerListenerPolicyRuleCollectionRulesToMap(rulesItem vpcv1.LoadBalancerListenerPolicyRule) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}
	if rulesItem.Condition != nil {
		rulesMap["condition"] = rulesItem.Condition
	}
	if rulesItem.CreatedAt != nil {
		rulesMap["created_at"] = rulesItem.CreatedAt.String()
	}
	if rulesItem.Field != nil {
		rulesMap["field"] = rulesItem.Field
	}
	if rulesItem.Href != nil {
		rulesMap["href"] = rulesItem.Href
	}
	if rulesItem.ID != nil {
		rulesMap["id"] = rulesItem.ID
	}
	if rulesItem.ProvisioningStatus != nil {
		rulesMap["provisioning_status"] = rulesItem.ProvisioningStatus
	}
	if rulesItem.Type != nil {
		rulesMap["type"] = rulesItem.Type
	}
	if rulesItem.Value != nil {
		rulesMap["value"] = rulesItem.Value
	}

	return rulesMap
}
