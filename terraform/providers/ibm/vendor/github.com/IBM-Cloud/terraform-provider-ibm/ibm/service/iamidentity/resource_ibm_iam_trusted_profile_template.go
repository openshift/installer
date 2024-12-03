// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func ResourceIBMTrustedProfileTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMTrustedProfileTemplateCreate,
		ReadContext:   resourceIBMTrustedProfileTemplateRead,
		UpdateContext: resourceIBMTrustedProfileTemplateUpdate,
		DeleteContext: resourceIBMTrustedProfileTemplateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the account where the template resides.",
			},
			"name": {
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"name", "description", "profile"},
				Optional:     true,
				Description:  "The name of the trusted profile template. This is visible only in the enterprise account.",
			},
			"description": {
				Type:         schema.TypeString,
				AtLeastOneOf: []string{"name", "description", "profile"},
				Optional:     true,
				Description:  "The description of the trusted profile template. Describe the template for enterprise account users.",
			},
			"profile": {
				Type:         schema.TypeList,
				AtLeastOneOf: []string{"name", "description", "profile"},
				MaxItems:     1,
				Optional:     true,
				Description:  "Input body parameters for the TemplateProfileComponent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the Profile.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description of the Profile.",
						},
						"rules": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Rules for the Profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the claim rule to be created or updated.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the claim rule.",
									},
									"realm_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The realm name of the Idp this claim rule applies to. This field is required only if the type is specified as 'Profile-SAML'.",
									},
									"expiration": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Session expiration in seconds, only required if type is 'Profile-SAML'.",
									},
									"conditions": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Conditions of this claim rule.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"claim": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The claim to evaluate against.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operation to perform on the claim. valid values are EQUALS, NOT_EQUALS, EQUALS_IGNORE_CASE, NOT_EQUALS_IGNORE_CASE, CONTAINS, IN.",
												},
												"value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The stringified JSON value that the claim is compared to using the operator.",
												},
											},
										},
									},
								},
							},
						},
						"identities": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Identities for the Profile.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"iam_id": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "IAM ID of the identity.",
									},
									"identifier": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Identifier of the identity that can assume the trusted profiles. This can be a user identifier (IAM id), serviceid or crn. Internally it uses account id of the service id for the identifier 'serviceid' and for the identifier 'crn' it uses account id contained in the CRN.",
									},
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the identity.",
									},
									"accounts": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Only valid for the type user. Accounts from which a user can assume the trusted profile.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the identity that can assume the trusted profile. This is optional field for all the types of identities. When this field is not set for the identity type 'serviceid' then the description of the service id is used. Description is recommended for the identity type 'crn' E.g. 'Instance 1234 of IBM Cloud Service project'.",
									},
								},
							},
						},
					},
				},
			},
			"policy_template_references": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Existing policy templates that you can reference to assign access in the trusted profile component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of Access Policy Template.",
						},
						"version": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Version of Access Policy Template.",
						},
					},
				},
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the the template.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the the template.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of the the template.",
			}, "committed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Committed flag determines if the template is ready for assignment.",
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "History of the trusted profile template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the history entry.",
						},
						"params": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Params of the history entry.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entity tag for this templateId-version combination.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cloud resource name.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of when the template was created.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the creator.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of when the template was last modified.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the identity that made the latest modification.",
			},
		},
	}
}

func resourceIBMTrustedProfileTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk("template_id"); ok { // if template_id is present then we need to create a new version of this template instead
		return resourceIBMTrustedProfileTemplateCreateVersion(context, d, meta)
	}

	createProfileTemplateOptions := &iamidentityv1.CreateProfileTemplateOptions{}

	if _, ok := d.GetOk("name"); ok {
		createProfileTemplateOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createProfileTemplateOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("profile"); ok {
		profileModel, err := resourceIBMTrustedProfileTemplateMapToTemplateProfileComponentRequest(d.Get("profile.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createProfileTemplateOptions.SetProfile(profileModel)
	}
	if _, ok := d.GetOk("policy_template_references"); ok {
		var policyTemplateReferences []iamidentityv1.PolicyTemplateReference
		for _, v := range d.Get("policy_template_references").(*schema.Set).List() {
			value := v.(map[string]interface{})
			policyTemplateReferencesItem, err := resourceIBMTrustedProfileTemplateMapToPolicyTemplateReference(value)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplateReferences = append(policyTemplateReferences, *policyTemplateReferencesItem)
		}
		createProfileTemplateOptions.SetPolicyTemplateReferences(policyTemplateReferences)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	accountID := userDetails.UserAccount
	createProfileTemplateOptions.SetAccountID(accountID)

	trustedProfileTemplateResponse, response, err := iamIdentityClient.CreateProfileTemplateWithContext(context, createProfileTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProfileTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProfileTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId(buildResourceIdFromTemplateVersion(*trustedProfileTemplateResponse.ID, *trustedProfileTemplateResponse.Version))

	if d.Get("committed").(bool) {
		err := resourceIBMTrustedProfileTemplateCommit(context, d, meta)
		if err != nil {
			log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateCommit failed %s", err)
			return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateCommit failed %s", err))
		}
	}

	return resourceIBMTrustedProfileTemplateRead(context, d, meta)
}

func resourceIBMTrustedProfileTemplateCreateVersion(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createProfileTemplateVersionOptions := &iamidentityv1.CreateProfileTemplateVersionOptions{}

	id, _, err := parseResourceId(d.Get("template_id").(string))
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateRead failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateRead failed %s", err))
	}

	createProfileTemplateVersionOptions.SetTemplateID(id)

	if _, ok := d.GetOk("name"); ok {
		createProfileTemplateVersionOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createProfileTemplateVersionOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("profile"); ok {
		profileModel, err := resourceIBMTrustedProfileTemplateMapToTemplateProfileComponentRequest(d.Get("profile.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createProfileTemplateVersionOptions.SetProfile(profileModel)
	}
	if _, ok := d.GetOk("policy_template_references"); ok {
		var policyTemplateReferences []iamidentityv1.PolicyTemplateReference
		for _, v := range d.Get("policy_template_references").(*schema.Set).List() {
			value := v.(map[string]interface{})
			policyTemplateReferencesItem, err := resourceIBMTrustedProfileTemplateMapToPolicyTemplateReference(value)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplateReferences = append(policyTemplateReferences, *policyTemplateReferencesItem)
		}
		createProfileTemplateVersionOptions.SetPolicyTemplateReferences(policyTemplateReferences)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	accountID := userDetails.UserAccount
	createProfileTemplateVersionOptions.SetAccountID(accountID)

	trustedProfileTemplateVersionResponse, response, err := iamIdentityClient.CreateProfileTemplateVersionWithContext(context, createProfileTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProfileTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProfileTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId(buildResourceIdFromTemplateVersion(*trustedProfileTemplateVersionResponse.ID, *trustedProfileTemplateVersionResponse.Version))

	if d.Get("committed").(bool) {
		err := resourceIBMTrustedProfileTemplateCommit(context, d, meta)
		if err != nil {
			log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateCommit failed %s", err)
			return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateCommit failed %s", err))
		}
	}

	return resourceIBMTrustedProfileTemplateRead(context, d, meta)
}

func resourceIBMTrustedProfileTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getProfileTemplateVersionOptions := &iamidentityv1.GetProfileTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateRead failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateRead failed %s", err))
	}

	getProfileTemplateVersionOptions.SetTemplateID(id)
	getProfileTemplateVersionOptions.SetVersion(version)

	trustedProfileTemplateResponse, response, err := iamIdentityClient.GetProfileTemplateVersionWithContext(context, getProfileTemplateVersionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProfileTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProfileTemplateVersionWithContext failed %s\n%s", err, response))
	}

	if !core.IsNil(trustedProfileTemplateResponse.Version) {
		if err = d.Set("version", trustedProfileTemplateResponse.Version); err != nil {
			return diag.FromErr(fmt.Errorf("error setting version: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.AccountID) {
		if err = d.Set("account_id", trustedProfileTemplateResponse.AccountID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting account_id: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.Name) {
		if err = d.Set("name", trustedProfileTemplateResponse.Name); err != nil {
			return diag.FromErr(fmt.Errorf("error setting name: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.Description) {
		if err = d.Set("description", trustedProfileTemplateResponse.Description); err != nil {
			return diag.FromErr(fmt.Errorf("error setting description: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.Profile) {
		profileMap, err := resourceIBMTrustedProfileTemplateTemplateProfileComponentResponseToMap(trustedProfileTemplateResponse.Profile)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("profile", []map[string]interface{}{profileMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting profile: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.PolicyTemplateReferences) {
		var policyTemplateReferences []map[string]interface{}
		for _, policyTemplateReferencesItem := range trustedProfileTemplateResponse.PolicyTemplateReferences {
			policyTemplateReferencesItemMap, err := resourceIBMTrustedProfileTemplatePolicyTemplateReferenceToMap(&policyTemplateReferencesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplateReferences = append(policyTemplateReferences, policyTemplateReferencesItemMap)
		}
		if err = d.Set("policy_template_references", policyTemplateReferences); err != nil {
			return diag.FromErr(fmt.Errorf("error setting policy_template_references: %s", err))
		}
	}
	if err = d.Set("id", trustedProfileTemplateResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %s", err))
	}
	if !core.IsNil(trustedProfileTemplateResponse.Committed) {
		if err = d.Set("committed", trustedProfileTemplateResponse.Committed); err != nil {
			return diag.FromErr(fmt.Errorf("error setting committed: %s", err))
		}
	}
	var history []map[string]interface{}
	if !core.IsNil(trustedProfileTemplateResponse.History) {
		for _, historyItem := range trustedProfileTemplateResponse.History {
			historyItemMap, err := resourceIBMTrustedProfileTemplateEntityHistoryRecordToMap(&historyItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, historyItemMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("error setting history: %s", err))
	}

	if !core.IsNil(trustedProfileTemplateResponse.EntityTag) {
		if err = d.Set("entity_tag", trustedProfileTemplateResponse.EntityTag); err != nil {
			return diag.FromErr(fmt.Errorf("error setting entity_tag: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.CRN) {
		if err = d.Set("crn", trustedProfileTemplateResponse.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("error setting crn: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.CreatedAt) {
		if err = d.Set("created_at", trustedProfileTemplateResponse.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.CreatedByID) {
		if err = d.Set("created_by_id", trustedProfileTemplateResponse.CreatedByID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_by_id: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.LastModifiedAt) {
		if err = d.Set("last_modified_at", trustedProfileTemplateResponse.LastModifiedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting last_modified_at: %s", err))
		}
	}
	if !core.IsNil(trustedProfileTemplateResponse.LastModifiedByID) {
		if err = d.Set("last_modified_by_id", trustedProfileTemplateResponse.LastModifiedByID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting last_modified_by_id: %s", err))
		}
	}

	return nil
}

func resourceIBMTrustedProfileTemplateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateProfileTemplateVersionOptions := &iamidentityv1.UpdateProfileTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateUpdate failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateUpdate failed %s", err))
	}

	updateProfileTemplateVersionOptions.SetTemplateID(id)
	updateProfileTemplateVersionOptions.SetVersion(version)
	updateProfileTemplateVersionOptions.SetIfMatch(d.Get("entity_tag").(string))

	hasChange := false

	if d.HasChange("name") {
		updateProfileTemplateVersionOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateProfileTemplateVersionOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("profile") {
		profile, err := resourceIBMTrustedProfileTemplateMapToTemplateProfileComponentRequest(d.Get("profile.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateProfileTemplateVersionOptions.SetProfile(profile)
		hasChange = true
	}
	if d.HasChange("policy_template_references") {
		var policyTemplateReferences []iamidentityv1.PolicyTemplateReference
		for _, v := range d.Get("policy_template_references").(*schema.Set).List() {
			value := v.(map[string]interface{})
			policyTemplateReferencesItem, err := resourceIBMTrustedProfileTemplateMapToPolicyTemplateReference(value)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplateReferences = append(policyTemplateReferences, *policyTemplateReferencesItem)
		}
		updateProfileTemplateVersionOptions.SetPolicyTemplateReferences(policyTemplateReferences)
		hasChange = true
	}

	if hasChange {
		_, response, err := iamIdentityClient.UpdateProfileTemplateVersionWithContext(context, updateProfileTemplateVersionOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateProfileTemplateVersionWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateProfileTemplateVersionWithContext failed %s\n%s", err, response))
		}
	}

	if d.HasChange("committed") {
		if d.Get("committed").(bool) {
			err := resourceIBMTrustedProfileTemplateCommit(context, d, meta)
			if err != nil {
				log.Printf("[DEBUG] resourceIBMTrustedProfileTemplateCommit failed %s", err)
				return diag.FromErr(fmt.Errorf("resourceIBMTrustedProfileTemplateCommit failed %s", err))
			}
		} else {
			return diag.FromErr(fmt.Errorf("A committed template cannot be uncommitted"))
		}
	}

	return resourceIBMTrustedProfileTemplateRead(context, d, meta)
}

func resourceIBMTrustedProfileTemplateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProfileTemplateVersionOptions := &iamidentityv1.DeleteProfileTemplateVersionOptions{}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		log.Printf("[DEBUG] resourceIBMAccountSettingsTemplateDelete failed %s", err)
		return diag.FromErr(fmt.Errorf("resourceIBMAccountSettingsTemplateDelete failed %s", err))
	}

	deleteProfileTemplateVersionOptions.SetTemplateID(id)
	deleteProfileTemplateVersionOptions.SetVersion(version)

	response, err := iamIdentityClient.DeleteProfileTemplateVersionWithContext(context, deleteProfileTemplateVersionOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProfileTemplateVersionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteProfileTemplateVersionWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMTrustedProfileTemplateCommit(context context.Context, d *schema.ResourceData, meta interface{}) error {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return err
	}

	id, version, err := parseResourceId(d.Id())
	if err != nil {
		return err
	}

	commitTrustedProfileTemplateVersionOptions := iamIdentityClient.NewCommitProfileTemplateOptions(id, version)
	_, err = iamIdentityClient.CommitProfileTemplateWithContext(context, commitTrustedProfileTemplateVersionOptions)
	if err != nil {
		return err
	}

	return nil
}

func resourceIBMTrustedProfileTemplateMapToTemplateProfileComponentRequest(modelMap map[string]interface{}) (*iamidentityv1.TemplateProfileComponentRequest, error) {
	model := &iamidentityv1.TemplateProfileComponentRequest{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["rules"] != nil {
		var rules []iamidentityv1.TrustedProfileTemplateClaimRule
		for _, rulesItem := range modelMap["rules"].([]interface{}) {
			rulesItemModel, err := resourceIBMTrustedProfileTemplateMapToTrustedProfileTemplateClaimRule(rulesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			rules = append(rules, *rulesItemModel)
		}
		model.Rules = rules
	}
	if modelMap["identities"] != nil {
		var identities []iamidentityv1.ProfileIdentityRequest
		for _, identitiesItem := range modelMap["identities"].([]interface{}) {
			identitiesItemModel, err := resourceIBMTrustedProfileTemplateMapToProfileIdentityRequest(identitiesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			identities = append(identities, *identitiesItemModel)
		}
		model.Identities = identities
	}
	return model, nil
}

func resourceIBMTrustedProfileTemplateMapToTrustedProfileTemplateClaimRule(modelMap map[string]interface{}) (*iamidentityv1.TrustedProfileTemplateClaimRule, error) {
	model := &iamidentityv1.TrustedProfileTemplateClaimRule{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Type = core.StringPtr(modelMap["type"].(string))
	if modelMap["realm_name"] != nil && modelMap["realm_name"].(string) != "" {
		model.RealmName = core.StringPtr(modelMap["realm_name"].(string))
	}
	if modelMap["expiration"] != nil {
		model.Expiration = core.Int64Ptr(int64(modelMap["expiration"].(int)))
	}
	var conditions []iamidentityv1.ProfileClaimRuleConditions
	for _, conditionsItem := range modelMap["conditions"].([]interface{}) {
		conditionsItemModel, err := resourceIBMTrustedProfileTemplateMapToProfileClaimRuleConditions(conditionsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		conditions = append(conditions, *conditionsItemModel)
	}
	model.Conditions = conditions
	return model, nil
}

func resourceIBMTrustedProfileTemplateMapToProfileClaimRuleConditions(modelMap map[string]interface{}) (*iamidentityv1.ProfileClaimRuleConditions, error) {
	model := &iamidentityv1.ProfileClaimRuleConditions{}
	model.Claim = core.StringPtr(modelMap["claim"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMTrustedProfileTemplateMapToProfileIdentityRequest(modelMap map[string]interface{}) (*iamidentityv1.ProfileIdentityRequest, error) {
	model := &iamidentityv1.ProfileIdentityRequest{}
	model.Identifier = core.StringPtr(modelMap["identifier"].(string))
	model.Type = core.StringPtr(modelMap["type"].(string))
	if modelMap["accounts"] != nil {
		var accounts []string
		for _, accountsItem := range modelMap["accounts"].([]interface{}) {
			accounts = append(accounts, accountsItem.(string))
		}
		model.Accounts = accounts
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	return model, nil
}

func resourceIBMTrustedProfileTemplateMapToPolicyTemplateReference(modelMap map[string]interface{}) (*iamidentityv1.PolicyTemplateReference, error) {
	model := &iamidentityv1.PolicyTemplateReference{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	model.Version = core.StringPtr(modelMap["version"].(string))
	return model, nil
}

func resourceIBMTrustedProfileTemplateTemplateProfileComponentResponseToMap(model *iamidentityv1.TemplateProfileComponentResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Rules != nil {
		var rules []map[string]interface{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := resourceIBMTrustedProfileTemplateTrustedProfileTemplateClaimRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	if model.Identities != nil {
		var identities []map[string]interface{}
		for _, identitiesItem := range model.Identities {
			identitiesItemMap, err := resourceIBMTrustedProfileTemplateProfileIdentityResponseToMap(&identitiesItem)
			if err != nil {
				return modelMap, err
			}
			identities = append(identities, identitiesItemMap)
		}
		modelMap["identities"] = identities
	}
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateTrustedProfileTemplateClaimRuleToMap(model *iamidentityv1.TrustedProfileTemplateClaimRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	modelMap["type"] = model.Type
	if model.RealmName != nil {
		modelMap["realm_name"] = model.RealmName
	}
	if model.Expiration != nil {
		modelMap["expiration"] = flex.IntValue(model.Expiration)
	}
	var conditions []map[string]interface{}
	for _, conditionsItem := range model.Conditions {
		conditionsItemMap, err := resourceIBMTrustedProfileTemplateProfileClaimRuleConditionsToMap(&conditionsItem)
		if err != nil {
			return modelMap, err
		}
		conditions = append(conditions, conditionsItemMap)
	}
	modelMap["conditions"] = conditions
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateProfileClaimRuleConditionsToMap(model *iamidentityv1.ProfileClaimRuleConditions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["claim"] = model.Claim
	modelMap["operator"] = model.Operator
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateProfileIdentityResponseToMap(model *iamidentityv1.ProfileIdentityResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["iam_id"] = model.IamID
	modelMap["identifier"] = model.Identifier
	modelMap["type"] = model.Type
	if model.Accounts != nil {
		modelMap["accounts"] = model.Accounts
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplatePolicyTemplateReferenceToMap(model *iamidentityv1.PolicyTemplateReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["version"] = model.Version
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateEntityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = model.Timestamp
	modelMap["iam_id"] = model.IamID
	modelMap["iam_id_account"] = model.IamIDAccount
	modelMap["action"] = model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = model.Message
	return modelMap, nil
}
