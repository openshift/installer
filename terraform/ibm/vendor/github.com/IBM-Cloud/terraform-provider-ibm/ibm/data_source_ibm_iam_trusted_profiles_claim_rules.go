// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func dataSourceIBMIamTrustedProfileClaimRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileClaimRuleListRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
			},
			"rules": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of claim rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the unique identifier of the claim rule.",
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
				},
			},
		},
	}
}

func dataSourceIBMIamTrustedProfileClaimRuleListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listClaimRulesOptions := &iamidentityv1.ListClaimRulesOptions{}

	listClaimRulesOptions.SetProfileID(d.Get("profile_id").(string))

	profileClaimRuleList, response, err := iamIdentityClient.ListClaimRules(listClaimRulesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListClaimRules failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListClaimRules failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIamTrustedProfileClaimRuleListID(d))

	if profileClaimRuleList.Rules != nil {
		err = d.Set("rules", dataSourceProfileClaimRuleListFlattenRules(profileClaimRuleList.Rules))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting rules %s", err))
		}
	}

	return nil
}

// dataSourceIBMIamTrustedProfileClaimRuleListID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfileClaimRuleListID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceProfileClaimRuleListFlattenRules(result []iamidentityv1.ProfileClaimRule) (rules []map[string]interface{}) {
	for _, rulesItem := range result {
		rules = append(rules, dataSourceProfileClaimRuleListRuleToMap(rulesItem))
	}

	return rules
}

func dataSourceProfileClaimRuleListRuleToMap(rulesItem iamidentityv1.ProfileClaimRule) (rulesMap map[string]interface{}) {
	rulesMap = map[string]interface{}{}

	if rulesItem.ID != nil {
		rulesMap["id"] = rulesItem.ID
	}
	if rulesItem.EntityTag != nil {
		rulesMap["entity_tag"] = rulesItem.EntityTag
	}
	if rulesItem.CreatedAt != nil {
		rulesMap["created_at"] = rulesItem.CreatedAt.String()
	}
	if rulesItem.ModifiedAt != nil {
		rulesMap["modified_at"] = rulesItem.ModifiedAt.String()
	}
	if rulesItem.Name != nil {
		rulesMap["name"] = rulesItem.Name
	}
	if rulesItem.Type != nil {
		rulesMap["type"] = rulesItem.Type
	}
	if rulesItem.RealmName != nil {
		rulesMap["realm_name"] = rulesItem.RealmName
	}
	if rulesItem.Expiration != nil {
		rulesMap["expiration"] = rulesItem.Expiration
	}
	if rulesItem.CrType != nil {
		rulesMap["cr_type"] = rulesItem.CrType
	}
	if rulesItem.Conditions != nil {
		conditionsList := []map[string]interface{}{}
		for _, conditionsItem := range rulesItem.Conditions {
			conditionsList = append(conditionsList, dataSourceProfileClaimRuleListRulesConditionsToMap(conditionsItem))
		}
		rulesMap["conditions"] = conditionsList
	}

	return rulesMap
}

func dataSourceProfileClaimRuleListRulesConditionsToMap(conditionsItem iamidentityv1.ProfileClaimRuleConditions) (conditionsMap map[string]interface{}) {
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
