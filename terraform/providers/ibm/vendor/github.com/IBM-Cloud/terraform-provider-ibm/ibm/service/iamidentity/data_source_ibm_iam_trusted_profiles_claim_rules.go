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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMIamTrustedProfileClaimRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIamTrustedProfileClaimRuleListRead,

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the trusted profile.",
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_iam_trusted_profile_claim_rules",
					"profile_id"),
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profiles_claim_rules", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listClaimRulesOptions := &iamidentityv1.ListClaimRulesOptions{}

	listClaimRulesOptions.SetProfileID(d.Get("profile_id").(string))

	profileClaimRuleList, _, err := iamIdentityClient.ListClaimRulesWithContext(context, listClaimRulesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListClaimRulesWithContext failed: %s", err.Error()), "(Data) ibm_iam_trusted_profiles_claim_rules", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIamTrustedProfilesClaimRulesID(d))

	rules := []map[string]interface{}{}
	for _, rulesItem := range profileClaimRuleList.Rules {
		rulesItemMap, err := DataSourceIBMIamTrustedProfilesClaimRulesProfileClaimRuleToMap(&rulesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_iam_trusted_profiles_claim_rules", "read", "rules-to-map").GetDiag()
		}
		rules = append(rules, rulesItemMap)
	}
	if err = d.Set("rules", rules); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting rules: %s", err), "(Data) ibm_iam_trusted_profiles_claim_rules", "read", "set-rules").GetDiag()
	}

	return nil
}

// dataSourceIBMIamTrustedProfilesClaimRulesID returns a reasonable ID for the list.
func dataSourceIBMIamTrustedProfilesClaimRulesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIamTrustedProfilesClaimRulesProfileClaimRuleToMap(model *iamidentityv1.ProfileClaimRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["entity_tag"] = *model.EntityTag
	modelMap["created_at"] = model.CreatedAt.String()
	if model.ModifiedAt != nil {
		modelMap["modified_at"] = model.ModifiedAt.String()
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	modelMap["type"] = *model.Type
	if model.RealmName != nil {
		modelMap["realm_name"] = *model.RealmName
	}
	modelMap["expiration"] = flex.IntValue(model.Expiration)
	if model.CrType != nil {
		modelMap["cr_type"] = *model.CrType
	}
	conditions := []map[string]interface{}{}
	for _, conditionsItem := range model.Conditions {
		conditionsItemMap, err := DataSourceIBMIamTrustedProfilesClaimRulesProfileClaimRuleConditionsToMap(&conditionsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		conditions = append(conditions, conditionsItemMap)
	}
	modelMap["conditions"] = conditions
	return modelMap, nil
}

func DataSourceIBMIamTrustedProfilesClaimRulesProfileClaimRuleConditionsToMap(model *iamidentityv1.ProfileClaimRuleConditions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["claim"] = *model.Claim
	modelMap["operator"] = *model.Operator
	modelMap["value"] = *model.Value
	return modelMap, nil
}
