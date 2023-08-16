// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileClaimRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileClaimRuleListRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_claim_rules",
					"profile_id"),
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of claim rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "the unique identifier of the claim rule.",
						},
						"entity_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "version of the claim rule.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If set contains a date time string of the creation date in ISO format.",
						},
						"modified_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "If set contains a date time string of the last modification date in ISO format.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The optional claim rule name.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the Calim rule, either 'Profile-SAML' or 'Profile-CR'.",
						},
						"realm_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The realm name of the Idp this claim rule applies to.",
						},
						"expiration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session expiration in seconds.",
						},
						"cr_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The compute resource type. Not required if type is Profile-SAML. Valid values are VSI, IKS_SA, ROKS_SA.",
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Conditions of this claim rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"claim": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The claim to evaluate against.",
									},
									"operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE, NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.",
									},
									"value": {
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
func DataSourceIBMIamTrustedProfileClaimRulesValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIamTrustedProfileClaimRulesValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_claim_rules", Schema: validateSchema}
	return &iBMIamTrustedProfileClaimRulesValidator
}

func dataSourceIBMIamTrustedProfileClaimRuleListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
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
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting rules %s", err))
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
