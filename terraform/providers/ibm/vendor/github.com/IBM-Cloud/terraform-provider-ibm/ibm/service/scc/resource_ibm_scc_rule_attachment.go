// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/configurationgovernancev1"
)

func ResourceIBMSccRuleAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccRuleAttachmentCreate,
		ReadContext:   resourceIBMSccRuleAttachmentRead,
		UpdateContext: resourceIBMSccRuleAttachmentUpdate,
		DeleteContext: resourceIBMSccRuleAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"attachment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID that uniquely identifies the attachment.",
			},
			"rule_id": &schema.Schema{
				ForceNew:    true,
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID that uniquely identifies the rule.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Your IBM Cloud account ID.",
			},
			"included_scope": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The extent at which the rule can be attached across your accounts.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"note": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A short description or alias to assign to the scope.",
						},
						"scope_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the scope, such as an enterprise, account, or account group, that you want to evaluate.",
						},
						"scope_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The type of scope that you want to evaluate.",
							ValidateFunc: validate.InvokeValidator("ibm_scc_rule_attachment", "scope_type"),
						},
					},
				},
				MaxItems: 1,
			},
			"excluded_scopes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The extent at which the rule can be excluded from the included scope.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"note": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A short description or alias to assign to the scope.",
						},
						"scope_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the scope, such as an enterprise, account, or account group, that you want to evaluate.",
						},
						"scope_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The type of scope that you want to evaluate.",
							ValidateFunc: validate.InvokeValidator("ibm_scc_rule_attachment", "scope_type"),
						},
					},
				},
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: customdiff.All(
			// update the version number via API GET if any of the fields are true
			customdiff.ComputedIf("version", func(_ context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				return diff.HasChange("included_scope") || diff.HasChange("excluded_scopes")
			}),
		),
	}
}

func resourceIBMSccRuleAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createRuleAttachmentsOptions := &configurationgovernancev1.CreateRuleAttachmentsOptions{}

	createRuleAttachmentsOptions.SetRuleID(d.Get("rule_id").(string))
	var attachment []configurationgovernancev1.RuleAttachmentRequest
	attachmentItem, err := resourceIBMSccRuleAttachmentMapToRuleAttachmentRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	attachment = append(attachment, *attachmentItem)
	createRuleAttachmentsOptions.SetAttachments(attachment)

	createRuleAttachmentsResponse, response, err := configurationGovernanceClient.CreateRuleAttachmentsWithContext(context, createRuleAttachmentsOptions)
	if err != nil || response.StatusCode > 300 {
		log.Printf("[DEBUG] CreateRuleAttachmentsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRuleAttachmentsWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createRuleAttachmentsOptions.RuleID, *createRuleAttachmentsResponse.Attachments[0].AttachmentID))

	return resourceIBMSccRuleAttachmentRead(context, d, meta)
}

func resourceIBMSccRuleAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleAttachmentOptions := &configurationgovernancev1.GetRuleAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleAttachmentOptions.SetRuleID(parts[0])
	getRuleAttachmentOptions.SetAttachmentID(parts[1])

	ruleAttachment, response, err := configurationGovernanceClient.GetRuleAttachmentWithContext(context, getRuleAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetRuleAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRuleAttachmentWithContext failed %s\n%s", err, response))
	}

	// TODO: handle argument of type []interface{}
	if err = d.Set("rule_id", ruleAttachment.RuleID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting rule_id: %s", err))
	}
	if err = d.Set("account_id", ruleAttachment.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	includedScopeMap, err := resourceIBMSccRuleAttachmentRuleScopeToMap(ruleAttachment.IncludedScope)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("included_scope", []map[string]interface{}{includedScopeMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting included_scope: %s", err))
	}

	excludedScope := []map[string]interface{}{}
	if ruleAttachment.ExcludedScopes != nil {
		for _, excludedScopeItem := range ruleAttachment.ExcludedScopes {
			excludedScopeItemMap, err := resourceIBMSccRuleAttachmentRuleScopeToMap(&excludedScopeItem)
			if err != nil {
				return diag.FromErr(err)
			}
			excludedScope = append(excludedScope, excludedScopeItemMap)
		}
	}
	if err = d.Set("excluded_scopes", excludedScope); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting excluded_scopes: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	return nil
}

func resourceIBMSccRuleAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateRuleAttachmentOptions := &configurationgovernancev1.UpdateRuleAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateRuleAttachmentOptions.SetRuleID(parts[0])
	updateRuleAttachmentOptions.SetAttachmentID(parts[1])

	// This code is never going to work since the schema has ForceNew in property rule_id
	// if d.HasChange("rule_id") {
	// 	return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
	// 		" The resource must be re-created to update this property.", "rule_id"))
	// }

	hasChange := d.HasChange("included_scope") || d.HasChange("excluded_scopes")

	if hasChange {
		updateRuleAttachmentOptions.SetIfMatch(d.Get("version").(string))
		updateRuleAttachmentOptions.SetRuleID(d.Get("rule_id").(string))
		updateRuleAttachmentOptions.SetAccountID(d.Get("account_id").(string))

		includedScope, err := resourceIBMSccRuleAttachmentMapToRuleScope(d.Get("included_scope.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateRuleAttachmentOptions.SetIncludedScope(includedScope)

		excludedScopes := []configurationgovernancev1.RuleScope{}
		if d.Get("excluded_scopes") != nil {
			for _, scopeItem := range d.Get("excluded_scopes").([]interface{}) {
				excludedScope, err := resourceIBMSccRuleAttachmentMapToRuleScope(scopeItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(err)
				}
				excludedScopes = append(excludedScopes, *excludedScope)
			}
		}
		updateRuleAttachmentOptions.SetExcludedScopes(excludedScopes)

		_, response, err := configurationGovernanceClient.UpdateRuleAttachmentWithContext(context, updateRuleAttachmentOptions)
		if err != nil || response.StatusCode > 300 {
			log.Printf("[DEBUG] UpdateRuleAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateRuleAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMSccRuleAttachmentRead(context, d, meta)
}

func resourceIBMSccRuleAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRuleAttachmentOptions := &configurationgovernancev1.DeleteRuleAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRuleAttachmentOptions.SetRuleID(parts[0])
	deleteRuleAttachmentOptions.SetAttachmentID(parts[1])

	response, err := configurationGovernanceClient.DeleteRuleAttachmentWithContext(context, deleteRuleAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRuleAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteRuleAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMSccRuleAttachmentMapToRuleAttachmentRequest(d *schema.ResourceData) (*configurationgovernancev1.RuleAttachmentRequest, error) {
	model := &configurationgovernancev1.RuleAttachmentRequest{}
	model.AccountID = core.StringPtr(d.Get("account_id").(string))
	includedScopeList := d.Get("included_scope").([]interface{})
	IncludedScopeModel, err := resourceIBMSccRuleAttachmentMapToRuleScope(includedScopeList[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.IncludedScope = IncludedScopeModel
	if d.Get("excluded_scopes") != nil {
		excludedScopes := []configurationgovernancev1.RuleScope{}
		for _, excludedScopesItem := range d.Get("excluded_scopes").([]interface{}) {
			excludedScopesItemModel, err := resourceIBMSccRuleAttachmentMapToRuleScope(excludedScopesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			excludedScopes = append(excludedScopes, *excludedScopesItemModel)
		}
		model.ExcludedScopes = excludedScopes
	}
	return model, nil
}

func resourceIBMSccRuleAttachmentMapToRuleScope(modelMap map[string]interface{}) (*configurationgovernancev1.RuleScope, error) {
	model := &configurationgovernancev1.RuleScope{}
	if modelMap["note"] != nil {
		model.Note = core.StringPtr(modelMap["note"].(string))
	}
	model.ScopeID = core.StringPtr(modelMap["scope_id"].(string))
	model.ScopeType = core.StringPtr(modelMap["scope_type"].(string))
	return model, nil
}

func resourceIBMSccRuleAttachmentRuleAttachmentRequestToMap(model *configurationgovernancev1.RuleAttachmentRequest) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["account_id"] = model.AccountID
	includedScopeMap, err := resourceIBMSccRuleAttachmentRuleScopeToMap(model.IncludedScope)
	if err != nil {
		return modelMap, err
	}
	modelMap["included_scope"] = []map[string]interface{}{includedScopeMap}
	if model.ExcludedScopes != nil {
		excludedScopes := []map[string]interface{}{}
		for _, excludedScopesItem := range model.ExcludedScopes {
			excludedScopesItemMap, err := resourceIBMSccRuleAttachmentRuleScopeToMap(&excludedScopesItem)
			if err != nil {
				return modelMap, err
			}
			excludedScopes = append(excludedScopes, excludedScopesItemMap)
		}
		modelMap["excluded_scopes"] = excludedScopes
	}
	return modelMap, nil
}

func resourceIBMSccRuleAttachmentRuleScopeToMap(model *configurationgovernancev1.RuleScope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Note != nil {
		modelMap["note"] = model.Note
	}
	modelMap["scope_id"] = model.ScopeID
	modelMap["scope_type"] = model.ScopeType
	return modelMap, nil
}
