// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.98.0-8be2046a-20241205-162752
 */

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileClaimRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileClaimRuleRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_claim_rule",
					"profile_id"),
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
				Description: "Type of the claim rule, either 'Profile-SAML' or 'Profile-CR'.",
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
							Description: "The claim to evaluate against. [Learn more](/docs/account?topic=account-iam-condition-properties&interface=ui#cr-attribute-names).",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClaimRuleOptions := &iamidentityv1.GetClaimRuleOptions{}

	getClaimRuleOptions.SetProfileID(d.Get("profile_id").(string))
	getClaimRuleOptions.SetRuleID(d.Get("rule_id").(string))

	profileClaimRule, _, err := iamIdentityClient.GetClaimRuleWithContext(context, getClaimRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClaimRuleWithContext failed: %s", err.Error()), "(Data) ibm_iam_trusted_profile_claim_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(fmt.Sprintf("%s/%s", *getClaimRuleOptions.ProfileID, *profileClaimRule.ID))
	if err = d.Set("entity_tag", profileClaimRule.EntityTag); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting entity_tag: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-entity_tag").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(profileClaimRule.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-created_at").GetDiag()
	}

	if !core.IsNil(profileClaimRule.ModifiedAt) {
		if err = d.Set("modified_at", flex.DateTimeToString(profileClaimRule.ModifiedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting modified_at: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-modified_at").GetDiag()
		}
	}

	if !core.IsNil(profileClaimRule.Name) {
		if err = d.Set("name", profileClaimRule.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-name").GetDiag()
		}
	}

	if err = d.Set("type", profileClaimRule.Type); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-type").GetDiag()
	}

	if !core.IsNil(profileClaimRule.RealmName) {
		if err = d.Set("realm_name", profileClaimRule.RealmName); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting realm_name: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-realm_name").GetDiag()
		}
	}

	if err = d.Set("expiration", flex.IntValue(profileClaimRule.Expiration)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting expiration: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-expiration").GetDiag()
	}

	if !core.IsNil(profileClaimRule.CrType) {
		if err = d.Set("cr_type", profileClaimRule.CrType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cr_type: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-cr_type").GetDiag()
		}
	}

	conditions := []map[string]interface{}{}
	for _, conditionsItem := range profileClaimRule.Conditions {
		conditionsItemMap, err := DataSourceIBMIamTrustedProfileClaimRuleProfileClaimRuleConditionsToMap(&conditionsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "conditions-to-map").GetDiag()
		}
		conditions = append(conditions, conditionsItemMap)
	}
	if err = d.Set("conditions", conditions); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting conditions: %s", err), "(Data) ibm_iam_trusted_profile_claim_rule", "read", "set-conditions").GetDiag()
	}

	return nil
}

func DataSourceIBMIamTrustedProfileClaimRuleProfileClaimRuleConditionsToMap(model *iamidentityv1.ProfileClaimRuleConditions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["claim"] = *model.Claim
	modelMap["operator"] = *model.Operator
	modelMap["value"] = *model.Value
	return modelMap, nil
}
