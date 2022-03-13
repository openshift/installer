// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIBMIAMDynamicRule() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMIAMDynamicRuleCreate,
		Read:     resourceIBMIAMDynamicRuleRead,
		Update:   resourceIBMIAMDynamicRuleUpdate,
		Delete:   resourceIBMIAMDynamicRuleDelete,
		Exists:   resourceIBMIAMDynamicRuleExists,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Unique identifier of the access group",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Rule",
			},
			"expiration": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "The expiration in hours",
				ValidateFunc: validatePortRange(1, 24),
			},
			"identity_provider": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The realm name or identity proivider url",
			},
			"conditions": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "conditions info",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"claim": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateAllowedStringValue([]string{"EQUALS", "EQUALS_IGNORE_CASE", "IN", "NOT_EQUALS_IGNORE_CASE", "NOT_EQUALS", "CONTAINS"}),
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "id of the rule",
			},
		},
	}
}

func resourceIBMIAMDynamicRuleCreate(d *schema.ResourceData, meta interface{}) error {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return err
	}

	grpID := d.Get("access_group_id").(string)
	name := d.Get("name").(string)
	realm := d.Get("identity_provider").(string)
	expiration := int64(d.Get("expiration").(int))

	var cond []interface{}
	conditions := []iamaccessgroupsv2.RuleConditions{}
	if res, ok := d.GetOk("conditions"); ok {
		cond = res.([]interface{})
		for _, e := range cond {
			r, _ := e.(map[string]interface{})
			value := fmt.Sprintf("\"%s\"", r["value"].(string))
			claim := r["claim"].(string)
			operator := r["operator"].(string)
			conditionParam := iamaccessgroupsv2.RuleConditions{
				Claim:    &claim,
				Operator: &operator,
				Value:    &value,
			}
			conditions = append(conditions, conditionParam)
		}
	}

	addAccessGroupRuleOptions := &iamaccessgroupsv2.AddAccessGroupRuleOptions{
		AccessGroupID: &grpID,
		Name:          &name,
		RealmName:     &realm,
		Expiration:    &expiration,
		Conditions:    conditions,
	}
	rule, detailedResponse, err := iamAccessGroupsClient.AddAccessGroupRule(addAccessGroupRuleOptions)
	if err != nil || rule == nil {
		return fmt.Errorf("Error adding rule to Access Group(%s) %s. API Response: %s", grpID, err, detailedResponse)
	}
	ruleID := rule.ID
	d.SetId(fmt.Sprintf("%s/%s", grpID, *ruleID))

	return resourceIBMIAMDynamicRuleRead(d, meta)
}

func resourceIBMIAMDynamicRuleRead(d *schema.ResourceData, meta interface{}) error {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	grpID := parts[0]
	ruleID := parts[1]

	getAccessGroupRuleOptions := &iamaccessgroupsv2.GetAccessGroupRuleOptions{
		AccessGroupID: &grpID,
		RuleID:        &ruleID,
	}
	rule, detailResponse, err := iamAccessGroupsClient.GetAccessGroupRule(getAccessGroupRuleOptions)

	if err != nil || rule == nil {
		if detailResponse != nil && detailResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		} else {
			return fmt.Errorf("Error retrieving access group Rules: %s. API Response: %s", err, detailResponse)
		}
	}

	d.Set("access_group_id", grpID)
	d.Set("name", rule.Name)
	d.Set("expiration", rule.Expiration)
	d.Set("identity_provider", rule.RealmName)
	d.Set("conditions", flattenConditions(rule.Conditions))
	d.Set("rule_id", rule.ID)

	return nil
}

func resourceIBMIAMDynamicRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return err
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	grpID := parts[0]
	ruleID := parts[1]
	getAccessGroupRuleOptions := iamAccessGroupsClient.NewGetAccessGroupRuleOptions(grpID, ruleID)
	_, detailedResponse, err := iamAccessGroupsClient.GetAccessGroupRule(getAccessGroupRuleOptions)
	if err != nil {
		return fmt.Errorf("Error retrieving access group Rules: %s. API Response: %s", err, detailedResponse)
	}

	etag := detailedResponse.GetHeaders().Get("etag")
	realm := d.Get("identity_provider").(string)
	expiration := int64(d.Get("expiration").(int))

	var cond []interface{}
	condition := []iamaccessgroupsv2.RuleConditions{}
	if res, ok := d.GetOk("conditions"); ok {
		cond = res.([]interface{})
		for _, e := range cond {
			r, _ := e.(map[string]interface{})
			value := fmt.Sprintf("\"%s\"", r["value"].(string))
			claim := r["claim"].(string)
			operator := r["operator"].(string)
			conditionParam := iamaccessgroupsv2.RuleConditions{
				Claim:    &claim,
				Operator: &operator,
				Value:    &value,
			}
			condition = append(condition, conditionParam)
		}
	}

	replaceAccessGroupRuleOption := iamAccessGroupsClient.NewReplaceAccessGroupRuleOptions(grpID, ruleID, etag, expiration, realm, condition)
	rule, detailedResponse, err := iamAccessGroupsClient.ReplaceAccessGroupRule(replaceAccessGroupRuleOption)
	if err != nil || rule == nil {
		return fmt.Errorf("Error replacing group(%s) rule(%s). API response: %s", grpID, ruleID, detailedResponse)
	}

	return resourceIBMIAMDynamicRuleRead(d, meta)

}

func resourceIBMIAMDynamicRuleDelete(d *schema.ResourceData, meta interface{}) error {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return err
	}

	grpID := parts[0]
	ruleID := parts[1]
	removeAccessGroupRuleOptions := iamAccessGroupsClient.NewRemoveAccessGroupRuleOptions(grpID, ruleID)
	detailedResponse, err := iamAccessGroupsClient.RemoveAccessGroupRule(removeAccessGroupRuleOptions)
	if err != nil {
		if detailedResponse != nil && detailedResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error getting group(%s) rule(%s). API Response: %s", grpID, ruleID, detailedResponse)
	}

	return nil
}

func resourceIBMIAMDynamicRuleExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	iamAccessGroupsClient, err := meta.(ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return false, err
	}

	parts, err := idParts(d.Id())
	if err != nil {
		return false, err
	}
	if len(parts) < 2 {
		return false, fmt.Errorf("Incorrect ID %s: Id should be a combination of accessGroupID/RuleID", d.Id())
	}
	grpID := parts[0]
	ruleID := parts[1]

	getAccessGroupRuleOptions := iamAccessGroupsClient.NewGetAccessGroupRuleOptions(grpID, ruleID)
	rule, detailResponse, err := iamAccessGroupsClient.GetAccessGroupRule(getAccessGroupRuleOptions)

	if detailResponse != nil && detailResponse.StatusCode == 404 {
		d.SetId("")
		return false, nil
	}
	if err != nil || rule == nil {
		return false, fmt.Errorf("Error getting group(%s) rule(%s). API response: %s", grpID, ruleID, detailResponse)
	}
	return *rule.AccessGroupID == grpID, nil
}
