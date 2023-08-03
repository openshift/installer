// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileClaimRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileClaimRuleRead,

		Schema: map[string]*schema.Schema{
			"profile_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_claim_rule",
					"profile_id"),
			},
			"rule_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the claim rule to get.",
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
	}
}

func DataSourceIBMIamTrustedProfileClaimRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIamTrustedProfileClaimRuleValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_claim_rule", Schema: validateSchema}
	return &iBMIamTrustedProfileClaimRuleValidator
}

func dataSourceIBMIamTrustedProfileClaimRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
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
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(profileClaimRule.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", flex.DateTimeToString(profileClaimRule.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting modified_at: %s", err))
	}
	if err = d.Set("name", profileClaimRule.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("type", profileClaimRule.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}
	if err = d.Set("realm_name", profileClaimRule.RealmName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting realm_name: %s", err))
	}
	if err = d.Set("expiration", flex.IntValue(profileClaimRule.Expiration)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting expiration: %s", err))
	}
	if err = d.Set("cr_type", profileClaimRule.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting cr_type: %s", err))
	}

	if profileClaimRule.Conditions != nil {
		err = d.Set("conditions", dataSourceProfileClaimRuleFlattenConditions(profileClaimRule.Conditions))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting conditions %s", err))
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
