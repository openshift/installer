// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

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
	"github.com/IBM/scc-go-sdk/v3/configurationgovernancev1"
)

func ResourceIBMSccTemplateAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccTemplateAttachmentCreate,
		ReadContext:   resourceIBMSccTemplateAttachmentRead,
		UpdateContext: resourceIBMSccTemplateAttachmentUpdate,
		DeleteContext: resourceIBMSccTemplateAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"attachment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID that uniquely identifies the template.",
			},
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID that uniquely identifies the template.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Your IBM Cloud account ID.",
			},
			"included_scope": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The extent at which the template can be attached across your accounts.",
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
							Description: "The ID of the scope, such as an enterprise, account, or account group, where you want to apply the customized defaults that are associated with a template.",
						},
						"scope_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The type of scope.",
							ValidateFunc: validate.InvokeValidator("ibm_scc_template_attachment", "scope_type"),
						},
					},
				},
			},
			"excluded_scopes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
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
							Description: "The ID of the scope, such as an enterprise, account, or account group, where you want to apply the customized defaults that are associated with a template.",
						},
						"scope_type": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The type of scope.",
							ValidateFunc: validate.InvokeValidator("ibm_scc_template_attachment", "scope_type"),
						},
					},
				},
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIBMSccTemplateAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createTemplateAttachmentsOptions := &configurationgovernancev1.CreateTemplateAttachmentsOptions{}

	createTemplateAttachmentsOptions.SetTemplateID(d.Get("template_id").(string))
	var attachment []configurationgovernancev1.TemplateAttachmentRequest
	attachmentItem, err := resourceIBMSccTemplateAttachmentMapToTemplateAttachmentRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}
	attachment = append(attachment, *attachmentItem)
	createTemplateAttachmentsOptions.SetAttachments(attachment)

	createTemplateAttachmentsResponse, response, err := configurationGovernanceClient.CreateTemplateAttachmentsWithContext(context, createTemplateAttachmentsOptions)
	if err != nil || response.StatusCode > 300 {
		log.Printf("[DEBUG] CreateTemplateAttachmentsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTemplateAttachmentsWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createTemplateAttachmentsOptions.TemplateID, *createTemplateAttachmentsResponse.Attachments[0].AttachmentID))

	return resourceIBMSccTemplateAttachmentRead(context, d, meta)
}

func resourceIBMSccTemplateAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getTemplateAttachmentOptions := &configurationgovernancev1.GetTemplateAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTemplateAttachmentOptions.SetTemplateID(parts[0])
	getTemplateAttachmentOptions.SetAttachmentID(parts[1])

	templateAttachment, response, err := configurationGovernanceClient.GetTemplateAttachmentWithContext(context, getTemplateAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTemplateAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTemplateAttachmentWithContext failed %s\n%s", err, response))
	}

	// TODO: handle argument of type []interface{}
	if err = d.Set("template_id", templateAttachment.TemplateID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting template_id: %s", err))
	}
	if err = d.Set("account_id", templateAttachment.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	includedScopeMap, err := resourceIBMSccTemplateAttachmentTemplateScopeToMap(templateAttachment.IncludedScope)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("included_scope", []map[string]interface{}{includedScopeMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting included_scope: %s", err))
	}

	excludedScope := []map[string]interface{}{}
	if templateAttachment.ExcludedScopes != nil {
		for _, excludedScopeItem := range templateAttachment.ExcludedScopes {
			excludedScopeItemMap, err := resourceIBMSccTemplateAttachmentTemplateScopeToMap(&excludedScopeItem)
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

func resourceIBMSccTemplateAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateTemplateAttachmentOptions := &configurationgovernancev1.UpdateTemplateAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateTemplateAttachmentOptions.SetTemplateID(parts[0])
	updateTemplateAttachmentOptions.SetAttachmentID(parts[1])

	if d.HasChange("template_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "template_id"))
	}

	hasChange := d.HasChange("included_scope") || d.HasChange("excluded_scopes")

	updateTemplateAttachmentOptions.SetIfMatch(d.Get("version").(string))

	if hasChange {
		updateTemplateAttachmentOptions.SetIfMatch(d.Get("version").(string))
		updateTemplateAttachmentOptions.SetTemplateID(d.Get("template_id").(string))
		updateTemplateAttachmentOptions.SetAccountID(d.Get("account_id").(string))

		includedScope, err := resourceIBMSccTemplateAttachmentMapToTemplateScope(d.Get("included_scope.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTemplateAttachmentOptions.SetIncludedScope(includedScope)

		excludedScopes := []configurationgovernancev1.TemplateScope{}
		if d.Get("excluded_scopes") != nil {
			for _, scopeItem := range d.Get("excluded_scopes").([]interface{}) {
				excludedScope, err := resourceIBMSccTemplateAttachmentMapToTemplateScope(scopeItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(err)
				}
				excludedScopes = append(excludedScopes, *excludedScope)
			}
		}
		updateTemplateAttachmentOptions.SetExcludedScopes(excludedScopes)

		_, response, err := configurationGovernanceClient.UpdateTemplateAttachmentWithContext(context, updateTemplateAttachmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTemplateAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTemplateAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMSccTemplateAttachmentRead(context, d, meta)
}

func resourceIBMSccTemplateAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTemplateAttachmentOptions := &configurationgovernancev1.DeleteTemplateAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTemplateAttachmentOptions.SetTemplateID(parts[0])
	deleteTemplateAttachmentOptions.SetAttachmentID(parts[1])

	response, err := configurationGovernanceClient.DeleteTemplateAttachmentWithContext(context, deleteTemplateAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTemplateAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTemplateAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMSccTemplateAttachmentMapToTemplateAttachmentRequest(d *schema.ResourceData) (*configurationgovernancev1.TemplateAttachmentRequest, error) {
	model := &configurationgovernancev1.TemplateAttachmentRequest{}
	model.AccountID = core.StringPtr(d.Get("account_id").(string))
	IncludedScopeModel, err := resourceIBMSccTemplateAttachmentMapToTemplateScope(d.Get("included_scope.0").(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.IncludedScope = IncludedScopeModel
	if d.Get("excluded_scopes") != nil {
		excludedScopes := []configurationgovernancev1.TemplateScope{}
		for _, excludedScopesItem := range d.Get("excluded_scopes").([]interface{}) {
			excludedScopesItemModel, err := resourceIBMSccTemplateAttachmentMapToTemplateScope(excludedScopesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			excludedScopes = append(excludedScopes, *excludedScopesItemModel)
		}
		model.ExcludedScopes = excludedScopes
	}
	return model, nil
}

func resourceIBMSccTemplateAttachmentMapToTemplateScope(modelMap map[string]interface{}) (*configurationgovernancev1.TemplateScope, error) {
	model := &configurationgovernancev1.TemplateScope{}
	if modelMap["note"] != nil {
		model.Note = core.StringPtr(modelMap["note"].(string))
	}
	model.ScopeID = core.StringPtr(modelMap["scope_id"].(string))
	model.ScopeType = core.StringPtr(modelMap["scope_type"].(string))
	return model, nil
}

func resourceIBMSccTemplateAttachmentTemplateAttachmentRequestToMap(model *configurationgovernancev1.TemplateAttachmentRequest) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["account_id"] = model.AccountID
	includedScopeMap, err := resourceIBMSccTemplateAttachmentTemplateScopeToMap(model.IncludedScope)
	if err != nil {
		return modelMap, err
	}
	modelMap["included_scope"] = []map[string]interface{}{includedScopeMap}
	if model.ExcludedScopes != nil {
		excludedScopes := []map[string]interface{}{}
		for _, excludedScopesItem := range model.ExcludedScopes {
			excludedScopesItemMap, err := resourceIBMSccTemplateAttachmentTemplateScopeToMap(&excludedScopesItem)
			if err != nil {
				return modelMap, err
			}
			excludedScopes = append(excludedScopes, excludedScopesItemMap)
		}
		modelMap["excluded_scopes"] = excludedScopes
	}
	return modelMap, nil
}

func resourceIBMSccTemplateAttachmentTemplateScopeToMap(model *configurationgovernancev1.TemplateScope) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Note != nil {
		modelMap["note"] = model.Note
	}
	modelMap["scope_id"] = model.ScopeID
	modelMap["scope_type"] = model.ScopeType
	return modelMap, nil
}
