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

const (
	iamClaimRuleType     = "type"
	iamClaimRuleOperator = "operator"
)

func ResourceIBMIAMTrustedProfileClaimRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIamTrustedProfileClaimRuleCreate,
		ReadContext:   resourceIBMIamTrustedProfileClaimRuleRead,
		UpdateContext: resourceIBMIamTrustedProfileClaimRuleUpdate,
		DeleteContext: resourceIBMIamTrustedProfileClaimRuleDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"profile_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the trusted profile to create a claim rule.",
				ValidateFunc: validate.InvokeValidator("ibm_iam_trusted_profile_claim_rule",
					"profile_id"),
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the claim rule to be created or updated.",
			},
			iamClaimRuleType: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Type of the calim rule, either 'Profile-SAML' or 'Profile-CR'.",
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"Profile-SAML", "Profile-CR"}),
			},
			"realm_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The realm name of the Idp this claim rule applies to.",
			},
			"expiration": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Session expiration in seconds.",
			},
			"cr_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The compute resource type. Not required if type is Profile-SAML. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"conditions": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "Conditions of this claim rule.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"claim": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The claim to evaluate against. [Learn more](/docs/account?topic=account-iam-condition-properties&interface=ui#cr-attribute-names).",
						},
						iamClaimRuleOperator: {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE, NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.",
							ValidateFunc: validate.ValidateAllowedStringValues([]string{"EQUALS", "NOT_EQUALS", "EQUALS_IGNORE_CASE", "NOT_EQUALS_IGNORE_CASE", "CONTAINS", "IN"}),
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The stringified JSON value that the claim is compared to using the operator.",
						},
					},
				},
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
			"rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the unique identifier of the claim rule.",
			},
		},
	}
}
func ResourceIBMIAMTrustedProfileClaimRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "profile_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "iam",
			CloudDataRange:             []string{"service:trusted_profile", "resolved_to:id"},
			Required:                   true})

	iBMIAMTrustedProfileClaimRuleValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_claim_rule", Schema: validateSchema}
	return &iBMIAMTrustedProfileClaimRuleValidator
}

func resourceIBMIamTrustedProfileClaimRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createClaimRuleOptions := &iamidentityv1.CreateClaimRuleOptions{}

	createClaimRuleOptions.SetProfileID(d.Get("profile_id").(string))
	createClaimRuleOptions.SetType(d.Get("type").(string))
	var conditions []iamidentityv1.ProfileClaimRuleConditions
	for _, v := range d.Get("conditions").([]interface{}) {
		value := v.(map[string]interface{})
		conditionsItem, err := ResourceIBMIamTrustedProfileClaimRuleMapToProfileClaimRuleConditions(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "create", "parse-conditions").GetDiag()
		}
		conditions = append(conditions, *conditionsItem)
	}
	createClaimRuleOptions.SetConditions(conditions)
	if _, ok := d.GetOk("context"); ok {
		contextModel, err := ResourceIBMIamTrustedProfileClaimRuleMapToResponseContext(d.Get("context.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "create", "parse-context").GetDiag()
		}
		createClaimRuleOptions.SetContext(contextModel)
	}
	if _, ok := d.GetOk("name"); ok {
		createClaimRuleOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("realm_name"); ok {
		createClaimRuleOptions.SetRealmName(d.Get("realm_name").(string))
	}
	if _, ok := d.GetOk("cr_type"); ok {
		createClaimRuleOptions.SetCrType(d.Get("cr_type").(string))
	}
	if _, ok := d.GetOk("expiration"); ok {
		createClaimRuleOptions.SetExpiration(int64(d.Get("expiration").(int)))
	}

	profileClaimRule, _, err := iamIdentityClient.CreateClaimRuleWithContext(context, createClaimRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateClaimRuleWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_claim_rule", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createClaimRuleOptions.ProfileID, *profileClaimRule.ID))

	return resourceIBMIamTrustedProfileClaimRuleRead(context, d, meta)
}

func resourceIBMIamTrustedProfileClaimRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getClaimRuleOptions := &iamidentityv1.GetClaimRuleOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "sep-id-parts").GetDiag()
	}

	getClaimRuleOptions.SetProfileID(parts[0])
	getClaimRuleOptions.SetRuleID(parts[1])

	profileClaimRule, response, err := iamIdentityClient.GetClaimRuleWithContext(context, getClaimRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetClaimRuleWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_claim_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("profile_id", getClaimRuleOptions.ProfileID); err != nil {
		err = fmt.Errorf("Error setting profile_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-profile_id").GetDiag()
	}

	if !core.IsNil(profileClaimRule.Name) {
		if err = d.Set("name", profileClaimRule.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("type", profileClaimRule.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-type").GetDiag()
	}
	if !core.IsNil(profileClaimRule.RealmName) {
		if err = d.Set("realm_name", profileClaimRule.RealmName); err != nil {
			err = fmt.Errorf("Error setting realm_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-realm_name").GetDiag()
		}
	}
	if !core.IsNil(profileClaimRule.Expiration) {
		if err = d.Set("expiration", flex.IntValue(profileClaimRule.Expiration)); err != nil {
			err = fmt.Errorf("Error setting expiration: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-expiration").GetDiag()
		}
	}
	if !core.IsNil(profileClaimRule.CrType) {
		if err = d.Set("cr_type", profileClaimRule.CrType); err != nil {
			err = fmt.Errorf("Error setting cr_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-cr_type").GetDiag()
		}
	}
	conditions := []map[string]interface{}{}
	for _, conditionsItem := range profileClaimRule.Conditions {
		conditionsItemMap, err := ResourceIBMIamTrustedProfileClaimRuleProfileClaimRuleConditionsToMap(&conditionsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "conditions-to-map").GetDiag()
		}
		conditions = append(conditions, conditionsItemMap)
	}
	if err = d.Set("conditions", conditions); err != nil {
		err = fmt.Errorf("Error setting conditions: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-conditions").GetDiag()
	}
	if err = d.Set("entity_tag", profileClaimRule.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-entity_tag").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(profileClaimRule.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-created_at").GetDiag()
	}
	if !core.IsNil(profileClaimRule.ModifiedAt) {
		if err = d.Set("modified_at", flex.DateTimeToString(profileClaimRule.ModifiedAt)); err != nil {
			err = fmt.Errorf("Error setting modified_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-modified_at").GetDiag()
		}
	}
	if err = d.Set("rule_id", profileClaimRule.ID); err != nil {
		err = fmt.Errorf("Error setting rule_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "read", "set-rule_id").GetDiag()
	}

	return nil
}

func resourceIBMIamTrustedProfileClaimRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateClaimRuleOptions := &iamidentityv1.UpdateClaimRuleOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "update", "sep-id-parts").GetDiag()
	}

	updateClaimRuleOptions.SetIfMatch("*")
	updateClaimRuleOptions.SetRuleID(parts[1])
	updateClaimRuleOptions.SetProfileID(d.Get("profile_id").(string))
	updateClaimRuleOptions.SetType(d.Get("type").(string))
	var conditions []iamidentityv1.ProfileClaimRuleConditions
	for _, v := range d.Get("conditions").([]interface{}) {
		value := v.(map[string]interface{})
		conditionsItem, err := ResourceIBMIamTrustedProfileClaimRuleMapToProfileClaimRuleConditions(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "update", "parse-conditions").GetDiag()
		}
		conditions = append(conditions, *conditionsItem)
	}
	updateClaimRuleOptions.SetConditions(conditions)
	if _, ok := d.GetOk("context"); ok {
		context, err := ResourceIBMIamTrustedProfileClaimRuleMapToResponseContext(d.Get("context.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "update", "parse-context").GetDiag()
		}
		updateClaimRuleOptions.SetContext(context)
	}
	if _, ok := d.GetOk("name"); ok {
		updateClaimRuleOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("realm_name"); ok {
		updateClaimRuleOptions.SetRealmName(d.Get("realm_name").(string))
	}
	if _, ok := d.GetOk("cr_type"); ok {
		updateClaimRuleOptions.SetCrType(d.Get("cr_type").(string))
	}
	if _, ok := d.GetOk("expiration"); ok {
		updateClaimRuleOptions.SetExpiration(int64(d.Get("expiration").(int)))
	}

	_, _, err = iamIdentityClient.UpdateClaimRuleWithContext(context, updateClaimRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateClaimRuleWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_claim_rule", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIBMIamTrustedProfileClaimRuleRead(context, d, meta)
}

func resourceIBMIamTrustedProfileClaimRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteClaimRuleOptions := &iamidentityv1.DeleteClaimRuleOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_iam_trusted_profile_claim_rule", "delete", "sep-id-parts").GetDiag()
	}

	deleteClaimRuleOptions.SetProfileID(parts[0])
	deleteClaimRuleOptions.SetRuleID(parts[1])

	_, err = iamIdentityClient.DeleteClaimRuleWithContext(context, deleteClaimRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteClaimRuleWithContext failed: %s", err.Error()), "ibm_iam_trusted_profile_claim_rule", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIamTrustedProfileClaimRuleMapToProfileClaimRuleConditions(modelMap map[string]interface{}) (*iamidentityv1.ProfileClaimRuleConditions, error) {
	model := &iamidentityv1.ProfileClaimRuleConditions{}
	model.Claim = core.StringPtr(modelMap["claim"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIBMIamTrustedProfileClaimRuleMapToResponseContext(modelMap map[string]interface{}) (*iamidentityv1.ResponseContext, error) {
	model := &iamidentityv1.ResponseContext{}
	if modelMap["transaction_id"] != nil && modelMap["transaction_id"].(string) != "" {
		model.TransactionID = core.StringPtr(modelMap["transaction_id"].(string))
	}
	if modelMap["operation"] != nil && modelMap["operation"].(string) != "" {
		model.Operation = core.StringPtr(modelMap["operation"].(string))
	}
	if modelMap["user_agent"] != nil && modelMap["user_agent"].(string) != "" {
		model.UserAgent = core.StringPtr(modelMap["user_agent"].(string))
	}
	if modelMap["url"] != nil && modelMap["url"].(string) != "" {
		model.URL = core.StringPtr(modelMap["url"].(string))
	}
	if modelMap["instance_id"] != nil && modelMap["instance_id"].(string) != "" {
		model.InstanceID = core.StringPtr(modelMap["instance_id"].(string))
	}
	if modelMap["thread_id"] != nil && modelMap["thread_id"].(string) != "" {
		model.ThreadID = core.StringPtr(modelMap["thread_id"].(string))
	}
	if modelMap["host"] != nil && modelMap["host"].(string) != "" {
		model.Host = core.StringPtr(modelMap["host"].(string))
	}
	if modelMap["start_time"] != nil && modelMap["start_time"].(string) != "" {
		model.StartTime = core.StringPtr(modelMap["start_time"].(string))
	}
	if modelMap["end_time"] != nil && modelMap["end_time"].(string) != "" {
		model.EndTime = core.StringPtr(modelMap["end_time"].(string))
	}
	if modelMap["elapsed_time"] != nil && modelMap["elapsed_time"].(string) != "" {
		model.ElapsedTime = core.StringPtr(modelMap["elapsed_time"].(string))
	}
	if modelMap["cluster_name"] != nil && modelMap["cluster_name"].(string) != "" {
		model.ClusterName = core.StringPtr(modelMap["cluster_name"].(string))
	}
	return model, nil
}

func ResourceIBMIamTrustedProfileClaimRuleProfileClaimRuleConditionsToMap(model *iamidentityv1.ProfileClaimRuleConditions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["claim"] = *model.Claim
	modelMap["operator"] = *model.Operator
	modelMap["value"] = *model.Value
	return modelMap, nil
}
