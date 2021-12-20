// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func dataSourceIBMIamTrustedProfileClaimRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileClaimRuleRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
			},
			"rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the claim rule to get.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "version of the claim rule.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the creation date in ISO format.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If set contains a date time string of the last modification date in ISO format.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The optional claim rule name.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the Calim rule, either 'Profile-SAML' or 'Profile-CR'.",
			},
			"realm_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The realm name of the Idp this claim rule applies to.",
			},
			"expiration": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Session expiration in seconds.",
			},
			"cr_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The compute resource type. Not required if type is Profile-SAML. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"conditions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Conditions of this claim rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"claim": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The claim to evaluate against.",
						},
						"operator": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE, NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The stringified JSON value that the claim is compared to using the operator.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIamTrustedProfileClaimRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getClaimRuleOptions := &iamidentityv1.GetClaimRuleOptions{}

	getClaimRuleOptions.SetProfileID(d.Get("profile_id").(string))
	getClaimRuleOptions.SetRuleID(d.Get("rule_id").(string))

	profileClaimRule, response, err := iamIdentityClient.GetClaimRule(getClaimRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] GetClaimRule failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetClaimRule failed %s\n%s", err, response))
	}
	d.SetId(fmt.Sprintf("%s/%s", *getClaimRuleOptions.ProfileID, *profileClaimRule.ID))
	if err = d.Set("entity_tag", profileClaimRule.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(profileClaimRule.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", dateTimeToString(profileClaimRule.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}
	if err = d.Set("name", profileClaimRule.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("type", profileClaimRule.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("realm_name", profileClaimRule.RealmName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting realm_name: %s", err))
	}
	if err = d.Set("expiration", intValue(profileClaimRule.Expiration)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration: %s", err))
	}
	if err = d.Set("cr_type", profileClaimRule.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cr_type: %s", err))
	}

	if profileClaimRule.Conditions != nil {
		err = d.Set("conditions", dataSourceProfileClaimRuleFlattenConditions(profileClaimRule.Conditions))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting conditions %s", err))
		}
	}

	return nil
}

func dataSourceProfileClaimRuleFlattenConditions(result []iamidentityv1.ProfileClaimRuleConditions) (conditions []map[string]interface{}) {
	for _, conditionsItem := range result {
		conditions = append(conditions, dataSourceProfileClaimRuleConditionsToMap(conditionsItem))
	}

	return conditions
}

func dataSourceProfileClaimRuleConditionsToMap(conditionsItem iamidentityv1.ProfileClaimRuleConditions) (conditionsMap map[string]interface{}) {
	conditionsMap = map[string]interface{}{}

	if conditionsItem.Claim != nil {
		conditionsMap["claim"] = conditionsItem.Claim
	}
	if conditionsItem.Operator != nil {
		conditionsMap["operator"] = conditionsItem.Operator
	}
	if conditionsItem.Value != nil {
		conditionsMap["value"] = conditionsItem.Value
	}

	return conditionsMap
}
