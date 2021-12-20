// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

const (
	iamClaimRuleType     = "type"
	iamClaimRuleOperator = "operator"
)

func resourceIBMIamTrustedProfileClaimRule() *schema.Resource {
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
			},
			"rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier of this claim rule.",
			},
			iamClaimRuleType: &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Type of the calim rule, either 'Profile-SAML' or 'Profile-CR'.",
				ValidateFunc: validateAllowedStringValue([]string{"Profile-SAML", "Profile-CR"}),
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
							Description: "The claim to evaluate against.",
						},
						iamClaimRuleOperator: &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE, NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.",
							ValidateFunc: validateAllowedStringValue([]string{"EQUALS", "NOT_EQUALS", "EQUALS_IGNORE_CASE", "NOT_EQUALS_IGNORE_CASE", "CONTAINS", "IN"}),
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The stringified JSON value that the claim is compared to using the operator.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the claim rule to be created or updated.",
			},
			"realm_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The realm name of the Idp this claim rule applies to. This field is required only if the type is specified as 'Profile-SAML'.",
			},
			"cr_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The compute resource type the rule applies to, required only if type is specified as 'Profile-CR'. Valid values are VSI, IKS_SA, ROKS_SA.",
			},
			"expiration": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Session expiration in seconds, only required if type is 'Profile-SAML'.",
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
		},
	}
}

func resourceIBMIamTrustedProfileClaimRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createClaimRuleOptions := &iamidentityv1.CreateClaimRuleOptions{}
	profile := d.Get("profile_id").(string)
	createClaimRuleOptions.SetProfileID(profile)
	createClaimRuleOptions.SetType(d.Get("type").(string))
	var conditions []iamidentityv1.ProfileClaimRuleConditions
	for _, e := range d.Get("conditions").([]interface{}) {
		value := e.(map[string]interface{})
		conditionsItem := resourceIBMIamTrustedProfileClaimRuleMapToProfileClaimRuleConditions(value)
		conditions = append(conditions, conditionsItem)
	}
	createClaimRuleOptions.SetConditions(conditions)

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

	profileClaimRule, response, err := iamIdentityClient.CreateClaimRule(createClaimRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateClaimRule failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateClaimRule failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", profile, *profileClaimRule.ID))

	return resourceIBMIamTrustedProfileClaimRuleRead(context, d, meta)
}

func resourceIBMIamTrustedProfileClaimRuleMapToProfileClaimRuleConditions(profileClaimRuleConditionsMap map[string]interface{}) iamidentityv1.ProfileClaimRuleConditions {
	profileClaimRuleConditions := iamidentityv1.ProfileClaimRuleConditions{}

	profileClaimRuleConditions.Claim = core.StringPtr(profileClaimRuleConditionsMap["claim"].(string))
	profileClaimRuleConditions.Operator = core.StringPtr(profileClaimRuleConditionsMap["operator"].(string))
	profileClaimRuleConditions.Value = core.StringPtr(profileClaimRuleConditionsMap["value"].(string))

	return profileClaimRuleConditions
}

func resourceIBMIamTrustedProfileClaimRuleMapToResponseContext(responseContextMap map[string]interface{}) iamidentityv1.ResponseContext {
	responseContext := iamidentityv1.ResponseContext{}

	if responseContextMap["transaction_id"] != nil {
		responseContext.TransactionID = core.StringPtr(responseContextMap["transaction_id"].(string))
	}
	if responseContextMap["operation"] != nil {
		responseContext.Operation = core.StringPtr(responseContextMap["operation"].(string))
	}
	if responseContextMap["user_agent"] != nil {
		responseContext.UserAgent = core.StringPtr(responseContextMap["user_agent"].(string))
	}
	if responseContextMap["url"] != nil {
		responseContext.URL = core.StringPtr(responseContextMap["url"].(string))
	}
	if responseContextMap["instance_id"] != nil {
		responseContext.InstanceID = core.StringPtr(responseContextMap["instance_id"].(string))
	}
	if responseContextMap["thread_id"] != nil {
		responseContext.ThreadID = core.StringPtr(responseContextMap["thread_id"].(string))
	}
	if responseContextMap["host"] != nil {
		responseContext.Host = core.StringPtr(responseContextMap["host"].(string))
	}
	if responseContextMap["start_time"] != nil {
		responseContext.StartTime = core.StringPtr(responseContextMap["start_time"].(string))
	}
	if responseContextMap["end_time"] != nil {
		responseContext.EndTime = core.StringPtr(responseContextMap["end_time"].(string))
	}
	if responseContextMap["elapsed_time"] != nil {
		responseContext.ElapsedTime = core.StringPtr(responseContextMap["elapsed_time"].(string))
	}
	if responseContextMap["cluster_name"] != nil {
		responseContext.ClusterName = core.StringPtr(responseContextMap["cluster_name"].(string))
	}

	return responseContext
}

func resourceIBMIamTrustedProfileClaimRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("Invalid ID %s", err))
	}
	getClaimRuleOptions := &iamidentityv1.GetClaimRuleOptions{}

	getClaimRuleOptions.SetRuleID(parts[1])
	getClaimRuleOptions.SetProfileID(parts[0])

	profileClaimRule, response, err := iamIdentityClient.GetClaimRule(getClaimRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetClaimRule failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetClaimRule failed %s\n%s", err, response))
	}

	if err = d.Set("profile_id", getClaimRuleOptions.ProfileID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting profile_id: %s", err))
	}
	if err = d.Set("type", profileClaimRule.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	conditions := []map[string]interface{}{}
	for _, conditionsItem := range profileClaimRule.Conditions {
		conditionsItemMap := resourceIBMIamTrustedProfileClaimRuleProfileClaimRuleConditionsToMap(conditionsItem)
		conditions = append(conditions, conditionsItemMap)
	}
	if err = d.Set("conditions", conditions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting conditions: %s", err))
	}
	if err = d.Set("name", profileClaimRule.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("rule_id", profileClaimRule.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rule_id: %s", err))
	}
	if err = d.Set("realm_name", profileClaimRule.RealmName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting realm_name: %s", err))
	}
	if err = d.Set("cr_type", profileClaimRule.CrType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cr_type: %s", err))
	}
	if err = d.Set("expiration", intValue(profileClaimRule.Expiration)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting expiration: %s", err))
	}
	if err = d.Set("entity_tag", profileClaimRule.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(profileClaimRule.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", dateTimeToString(profileClaimRule.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}

	return nil
}

func resourceIBMIamTrustedProfileClaimRuleProfileClaimRuleConditionsToMap(profileClaimRuleConditions iamidentityv1.ProfileClaimRuleConditions) map[string]interface{} {
	profileClaimRuleConditionsMap := map[string]interface{}{}

	profileClaimRuleConditionsMap["claim"] = profileClaimRuleConditions.Claim
	profileClaimRuleConditionsMap["operator"] = profileClaimRuleConditions.Operator
	profileClaimRuleConditionsMap["value"] = profileClaimRuleConditions.Value

	return profileClaimRuleConditionsMap
}

func resourceIBMIamTrustedProfileClaimRuleResponseContextToMap(responseContext iamidentityv1.ResponseContext) map[string]interface{} {
	responseContextMap := map[string]interface{}{}

	if responseContext.TransactionID != nil {
		responseContextMap["transaction_id"] = responseContext.TransactionID
	}
	if responseContext.Operation != nil {
		responseContextMap["operation"] = responseContext.Operation
	}
	if responseContext.UserAgent != nil {
		responseContextMap["user_agent"] = responseContext.UserAgent
	}
	if responseContext.URL != nil {
		responseContextMap["url"] = responseContext.URL
	}
	if responseContext.InstanceID != nil {
		responseContextMap["instance_id"] = responseContext.InstanceID
	}
	if responseContext.ThreadID != nil {
		responseContextMap["thread_id"] = responseContext.ThreadID
	}
	if responseContext.Host != nil {
		responseContextMap["host"] = responseContext.Host
	}
	if responseContext.StartTime != nil {
		responseContextMap["start_time"] = responseContext.StartTime
	}
	if responseContext.EndTime != nil {
		responseContextMap["end_time"] = responseContext.EndTime
	}
	if responseContext.ElapsedTime != nil {
		responseContextMap["elapsed_time"] = responseContext.ElapsedTime
	}
	if responseContext.ClusterName != nil {
		responseContextMap["cluster_name"] = responseContext.ClusterName
	}

	return responseContextMap
}

func resourceIBMIamTrustedProfileClaimRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("Invalid ID %s", err))
	}
	updateClaimRuleOptions := &iamidentityv1.UpdateClaimRuleOptions{}

	updateClaimRuleOptions.SetIfMatch("*")
	updateClaimRuleOptions.SetRuleID(parts[1])
	updateClaimRuleOptions.SetProfileID(parts[0])
	updateClaimRuleOptions.SetType(d.Get("type").(string))
	var conditions []iamidentityv1.ProfileClaimRuleConditions
	for _, e := range d.Get("conditions").([]interface{}) {
		value := e.(map[string]interface{})
		conditionsItem := resourceIBMIamTrustedProfileClaimRuleMapToProfileClaimRuleConditions(value)
		conditions = append(conditions, conditionsItem)
	}
	updateClaimRuleOptions.SetConditions(conditions)

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

	_, response, err := iamIdentityClient.UpdateClaimRule(updateClaimRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateClaimRule failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateClaimRule failed %s\n%s", err, response))
	}

	return resourceIBMIamTrustedProfileClaimRuleRead(context, d, meta)
}

func resourceIBMIamTrustedProfileClaimRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("Invalid ID %s", err))
	}

	deleteClaimRuleOptions := &iamidentityv1.DeleteClaimRuleOptions{}

	deleteClaimRuleOptions.SetProfileID(parts[0])
	deleteClaimRuleOptions.SetRuleID(parts[1])

	response, err := iamIdentityClient.DeleteClaimRule(deleteClaimRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteClaimRule failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteClaimRule failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
