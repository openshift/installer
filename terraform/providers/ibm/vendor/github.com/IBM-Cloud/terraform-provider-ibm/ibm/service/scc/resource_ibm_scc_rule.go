// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func ResourceIbmSccRule() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		CreateContext: resourceIbmSccRuleCreate,
		ReadContext:   resourceIbmSccRuleRead,
		UpdateContext: resourceIbmSccRuleUpdate,
		DeleteContext: resourceIbmSccRuleDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Deprecation list
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A human-readable alias to assign to your rule.",
				Deprecated:  "name is now deprecated",
			},
			"rule_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of rule. Rules that you create are `user_defined`.",
				Deprecated:  "use type instead",
			},
			"creation_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the resource was created.",
				Deprecated:  "use created_on instead",
			},
			"modification_date": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the resource was last modified.",
				Deprecated:  "use updated_on instead",
			},
			"modified_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for the user or application that last modified the resource.",
				Deprecated:  "use updated_by",
			},
			"enforcement_actions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The actions that the service must run on your behalf when a request to create or modify the target resource does not comply with your conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "To block a request from completing, use `disallow`.",
						},
					},
				},
				MaxItems:   1,
				Deprecated: "enforcement_actions is now deprecated",
			},
			// End of Deprecation list
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the rule was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the rule.",
			},
			"description": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_rule", "description"),
				Description:  "The details of a rule's response.",
			},
			// Manual Intervention
			"etag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The etag of the rule.",
			},
			// End Manual Intervention
			"import": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The collection of import parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameters": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of import parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The import parameter name.",
									},
									"display_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The display name of the property.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The propery description.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The property type.",
									},
								},
							},
						},
					},
				},
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The list of labels.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"required_config": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The required configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The required config description.",
						},
						"and": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The `AND` required configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The required config description.",
									},
									"or": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The `OR` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The required config description.",
												},
												"property": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The property.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"and": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The `AND` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The required config description.",
												},
												"property": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The property.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"property": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The property.",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The operator.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema for any JSON type.",
									},
								},
							},
						},
						"or": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The `OR` required configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The required config description.",
									},
									"or": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The `OR` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The required config description.",
												},
												"property": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The property.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"and": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The `AND` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The required config description.",
												},
												"property": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The property.",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The operator.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"property": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The property.",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The operator.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Schema for any JSON type.",
									},
								},
							},
						},
						"property": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The property.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The operator.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Schema for any JSON type.",
						},
					},
				},
			},
			"target": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The rule target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target service name.",
						},
						"service_display_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The display name of the target service.",
							// Manual Intervention
							DiffSuppressFunc: func(_, oldVal, newVal string, d *schema.ResourceData) bool {
								if newVal == "" {
									return true
								}
								if strings.ToLower(oldVal) == strings.ToLower(newVal) {
									return true
								}
								return false
							},
							// End Manual Intervention
						},
						"resource_kind": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target resource kind.",
						},
						"additional_target_attributes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The list of targets supported properties.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The additional target attribute name.",
									},
									"operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The operator.",
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The value.",
									},
								},
							},
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule type (allowable values are `user_defined` or `system_defined`).",
			},
			"updated_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the rule was modified.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who modified the rule.",
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_rule", "version"),
				Description:  "The version number of a rule.",
			},
		},
	})
}

func ResourceIbmSccRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[A-Za-z0-9]+`,
			MinValueLength:             0,
			MaxValueLength:             512,
		},
		validate.ValidateSchema{
			Identifier:                 "version",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[0-9][0-9.]*$`,
			MinValueLength:             5,
			MaxValueLength:             10,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_rule", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configManagerClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	createRuleOptions := &securityandcompliancecenterapiv3.CreateRuleOptions{}

	createRuleOptions.SetDescription(d.Get("description").(string))
	// Manual Intervention
	targetModel, err := resourceIbmSccRuleMapToTarget(d.Get("target.0").(map[string]interface{}))
	// End Manual Intervention
	if err != nil {
		return diag.FromErr(err)
	}
	createRuleOptions.SetTarget(targetModel)
	requiredConfigModel, err := resourceIbmSccRuleMapToRequiredConfig(d.Get("required_config.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createRuleOptions.SetRequiredConfig(requiredConfigModel)
	if _, ok := d.GetOk("version"); ok {
		createRuleOptions.SetVersion(d.Get("version").(string))
	}
	if _, ok := d.GetOk("import"); ok {
		importVarModel, err := resourceIbmSccRuleMapToImport(d.Get("import.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createRuleOptions.SetImport(importVarModel)
	}
	if _, ok := d.GetOk("labels"); ok {
		labels := make([]string, 0)
		for _, v := range d.Get("labels").([]interface{}) {
			labelsItem := v.(string)
			labels = append(labels, labelsItem)
		}
		createRuleOptions.SetLabels(labels)
	}

	instance_id := d.Get("instance_id").(string)
	createRuleOptions.SetInstanceID(instance_id)
	rule, response, err := configManagerClient.CreateRuleWithContext(context, createRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId(instance_id + "/" + *rule.ID)

	return resourceIbmSccRuleRead(context, d, meta)
}

func resourceIbmSccRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configManagerClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleOptions := &securityandcompliancecenterapiv3.GetRuleOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	getRuleOptions.SetInstanceID(parts[0])
	getRuleOptions.SetRuleID(parts[1])

	rule, response, err := configManagerClient.GetRuleWithContext(context, getRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRuleWithContext failed %s\n%s", err, response))
	}
	// Manual Intervention
	if err = d.Set("instance_id", parts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("rule_id", parts[1]); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting instance_id: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("ETag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}
	// End Manual Intervention

	if err = d.Set("description", rule.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if !core.IsNil(rule.Version) {
		if err = d.Set("version", rule.Version); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
		}
	}

	if !core.IsNil(rule.Import) {
		importVarMap, err := resourceIbmSccRuleImportToMap(rule.Import)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("import", []map[string]interface{}{importVarMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting import: %s", err))
		}
	}

	targetMap, err := resourceIbmSccRuleTargetToMap(rule.Target)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target: %s", err))
	}
	requiredConfigMap, err := resourceIbmSccRuleRequiredConfigToMap(rule.RequiredConfig)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("required_config", []map[string]interface{}{requiredConfigMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting required_config: %s", err))
	}
	if !core.IsNil(rule.Labels) {
		log.Printf("[INFO] rule.Labels = %v\n", rule.Labels)
		if err = d.Set("labels", rule.Labels); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting labels: %s", err))
		}
	}
	if err = d.Set("created_on", flex.DateTimeToString(rule.CreatedOn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_on: %s", err))
	}
	if err = d.Set("created_by", rule.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("updated_on", flex.DateTimeToString(rule.UpdatedOn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_on: %s", err))
	}
	if err = d.Set("updated_by", rule.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_by: %s", err))
	}
	if err = d.Set("account_id", rule.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("type", rule.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	return nil
}

func resourceIbmSccRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configManagerClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceRuleOptions := &securityandcompliancecenterapiv3.ReplaceRuleOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	replaceRuleOptions.SetInstanceID(parts[0])
	replaceRuleOptions.SetRuleID(parts[1])
	replaceRuleOptions.SetIfMatch(d.Get("etag").(string))

	hasChange := false

	if d.HasChange("description") || d.HasChange("target") || d.HasChange("required_config") {
		replaceRuleOptions.SetDescription(d.Get("description").(string))
		target, err := resourceIbmSccRuleMapToTarget(d.Get("target.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceRuleOptions.SetTarget(target)
		requiredConfig, err := resourceIbmSccRuleMapToRequiredConfig(d.Get("required_config.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceRuleOptions.SetRequiredConfig(requiredConfig)
		hasChange = true
	}
	if d.HasChange("version") {
		replaceRuleOptions.SetVersion(d.Get("version").(string))
		hasChange = true
	}
	if d.HasChange("import") {
		importVar, err := resourceIbmSccRuleMapToImport(d.Get("import.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceRuleOptions.SetImport(importVar)
		hasChange = true
	}
	if d.HasChange("labels") {
		hasChange = true
	}

	if hasChange {
		if _, ok := d.GetOk("labels"); ok {
			labels := make([]string, 0)
			for _, v := range d.Get("labels").([]interface{}) {
				labelsItem := v.(string)
				labels = append(labels, labelsItem)
			}
			replaceRuleOptions.SetLabels(labels)
		}
		if _, ok := d.GetOk("import"); ok {
			importVar, err := resourceIbmSccRuleMapToImport(d.Get("import.0").(map[string]interface{}))
			if err != nil {
				return diag.FromErr(err)
			}
			replaceRuleOptions.SetImport(importVar)
		}
		_, response, err := configManagerClient.ReplaceRuleWithContext(context, replaceRuleOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceRuleWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceRuleWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSccRuleRead(context, d, meta)
}

func resourceIbmSccRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configManagerClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRuleOptions := &securityandcompliancecenterapiv3.DeleteRuleOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}
	deleteRuleOptions.SetInstanceID(parts[0])
	deleteRuleOptions.SetRuleID(parts[1])

	response, err := configManagerClient.DeleteRuleWithContext(context, deleteRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmSccRuleMapToTarget(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Target, error) {
	model := &securityandcompliancecenterapiv3.Target{}
	model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	if modelMap["service_display_name"] != nil && modelMap["service_display_name"].(string) != "" {
		model.ServiceDisplayName = core.StringPtr(modelMap["service_display_name"].(string))
	}
	model.ResourceKind = core.StringPtr(modelMap["resource_kind"].(string))
	if modelMap["additional_target_attributes"] != nil {
		additionalTargetAttributes := []securityandcompliancecenterapiv3.AdditionalTargetAttribute{}
		for _, additionalTargetAttributesItem := range modelMap["additional_target_attributes"].([]interface{}) {
			additionalTargetAttributesItemModel, err := resourceIbmSccRuleMapToAdditionalTargetAttribute(additionalTargetAttributesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, *additionalTargetAttributesItemModel)
		}
		model.AdditionalTargetAttributes = additionalTargetAttributes
	}
	return model, nil
}

func resourceIbmSccRuleMapToAdditionalTargetAttribute(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.AdditionalTargetAttribute, error) {
	model := &securityandcompliancecenterapiv3.AdditionalTargetAttribute{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	if modelMap["value"] != nil && modelMap["value"].(string) != "" {
		model.Value = core.StringPtr(modelMap["value"].(string))
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfig(modelMap map[string]interface{}) (securityandcompliancecenterapiv3.RequiredConfigIntf, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfig{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["and"] != nil {
		and := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			andItemModel, err := resourceIbmSccRuleMapToRequiredConfigItems(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	if modelMap["or"] != nil {
		or := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			orItemModel, err := resourceIbmSccRuleMapToRequiredConfigItems(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	if modelMap["property"] != nil && modelMap["property"].(string) != "" {
		model.Property = core.StringPtr(modelMap["property"].(string))
	}
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	// Manual Intervention
	if modelMap["value"] != nil {
		// model.Value = modelMap["value"].(string)
		sLit := strings.Trim(modelMap["value"].(string), "[]")
		sList := strings.Split(sLit, ",")
		if len(sList) == 1 {
			model.Value = modelMap["value"].(string)
		} else {
			model.Value = sList
		}

	}
	// End Manual Intervention
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigItems(modelMap map[string]interface{}) (securityandcompliancecenterapiv3.RequiredConfigItemsIntf, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigItems{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["or"] != nil {
		or := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			orItemModel, err := resourceIbmSccRuleMapToRequiredConfigItems(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	if modelMap["and"] != nil {
		and := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			andItemModel, err := resourceIbmSccRuleMapToRequiredConfigItems(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	if modelMap["property"] != nil && modelMap["property"].(string) != "" {
		model.Property = core.StringPtr(modelMap["property"].(string))
	}
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	// Manual Intervention
	if modelMap["value"] != nil && len(modelMap["value"].(string)) > 0 {
		// model.Value = modelMap["value"].(string)
		sLit := strings.Trim(modelMap["value"].(string), "[]")
		sList := strings.Split(sLit, ",")
		if len(sList) == 1 {
			model.Value = modelMap["value"].(string)
		} else {
			model.Value = sList
		}
	}
	// Manual Intervention
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigBase(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	if modelMap["value"] != nil {
		sLit := strings.Trim(modelMap["value"].(string), "[]")
		sList := strings.Split(sLit, ",")
		if len(sList) == 1 {
			model.Value = modelMap["value"].(string)
		} else {
			model.Value = sList
		}
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigItemsRequiredConfigOrDepth1(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["or"] != nil {
		or := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			orItemModel, err := resourceIbmSccRuleMapToRequiredConfigItemsRequiredConfigBase(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	if modelMap["and"] != nil {
		or := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, orItem := range modelMap["and"].([]interface{}) {
			orItemModel, err := resourceIbmSccRuleMapToRequiredConfigItemsRequiredConfigBase(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigItemsRequiredConfigAndDepth1(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["and"] != nil {
		and := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			andItemModel, err := resourceIbmSccRuleMapToRequiredConfigItemsRequiredConfigBase(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigItemsRequiredConfigBase(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	if modelMap["value"] != nil {
		model.Value = modelMap["value"].(string)
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigAnd(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["and"] != nil {
		and := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, andItem := range modelMap["and"].([]interface{}) {
			andItemModel, err := resourceIbmSccRuleMapToRequiredConfigItems(andItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			and = append(and, andItemModel)
		}
		model.And = and
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigOr(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["or"] != nil {
		or := []securityandcompliancecenterapiv3.RequiredConfigItemsIntf{}
		for _, orItem := range modelMap["or"].([]interface{}) {
			orItemModel, err := resourceIbmSccRuleMapToRequiredConfigItems(orItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			or = append(or, orItemModel)
		}
		model.Or = or
	}
	return model, nil
}

func resourceIbmSccRuleMapToRequiredConfigRequiredConfigBase(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase, error) {
	model := &securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Operator = core.StringPtr(modelMap["operator"].(string))
	if modelMap["value"] != nil {
		model.Value = modelMap["value"].(string)
	}
	return model, nil
}

func resourceIbmSccRuleMapToImport(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Import, error) {
	model := &securityandcompliancecenterapiv3.Import{}
	if modelMap["parameters"] != nil {
		parameters := []securityandcompliancecenterapiv3.Parameter{}
		for _, parametersItem := range modelMap["parameters"].([]interface{}) {
			parametersItemModel, err := resourceIbmSccRuleMapToParameter(parametersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			parameters = append(parameters, *parametersItemModel)
		}
		model.Parameters = parameters
	}
	return model, nil
}

func resourceIbmSccRuleMapToParameter(modelMap map[string]interface{}) (*securityandcompliancecenterapiv3.Parameter, error) {
	model := &securityandcompliancecenterapiv3.Parameter{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["display_name"] != nil && modelMap["display_name"].(string) != "" {
		model.DisplayName = core.StringPtr(modelMap["display_name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIbmSccRuleImportToMap(model *securityandcompliancecenterapiv3.Import) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := resourceIbmSccRuleParameterToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func resourceIbmSccRuleParameterToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.DisplayName != nil {
		modelMap["display_name"] = model.DisplayName
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIbmSccRuleTargetToMap(model *securityandcompliancecenterapiv3.Target) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["service_name"] = model.ServiceName
	if model.ServiceDisplayName != nil {
		modelMap["service_display_name"] = model.ServiceDisplayName
	}
	modelMap["resource_kind"] = model.ResourceKind
	if model.AdditionalTargetAttributes != nil {
		additionalTargetAttributes := []map[string]interface{}{}
		for _, additionalTargetAttributesItem := range model.AdditionalTargetAttributes {
			additionalTargetAttributesItemMap, err := resourceIbmSccRuleAdditionalTargetAttributeToMap(&additionalTargetAttributesItem)
			if err != nil {
				return modelMap, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, additionalTargetAttributesItemMap)
		}
		modelMap["additional_target_attributes"] = additionalTargetAttributes
	}
	return modelMap, nil
}

func resourceIbmSccRuleAdditionalTargetAttributeToMap(model *securityandcompliancecenterapiv3.AdditionalTargetAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigToMap(model securityandcompliancecenterapiv3.RequiredConfigIntf) (map[string]interface{}, error) {
	if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd); ok {
		return resourceIbmSccRuleRequiredConfigAndToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr); ok {
		return resourceIbmSccRuleRequiredConfigOrToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase); ok {
		return resourceIbmSccRuleRequiredConfigRequiredConfigBaseToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfig); ok {
		modelMap := make(map[string]interface{})
		model := model.(*securityandcompliancecenterapiv3.RequiredConfig)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.And != nil {
			and := []map[string]interface{}{}
			for _, andItem := range model.And {
				andItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(andItem)
				if err != nil {
					return modelMap, err
				}
				and = append(and, andItemMap)
			}
			modelMap["and"] = and
		}
		if model.Or != nil {
			or := []map[string]interface{}{}
			for _, orItem := range model.Or {
				orItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(orItem)
				if err != nil {
					return modelMap, err
				}
				or = append(or, orItemMap)
			}
			modelMap["or"] = or
		}
		if model.Property != nil {
			modelMap["property"] = model.Property
		}
		if model.Operator != nil {
			modelMap["operator"] = model.Operator
		}
		// Manual Intervention
		if model.Value != nil {
			// model.Value is a schema.TypeString, so it needs to converted to String Type
			switch v := model.Value.(type) {
			case string:
				modelMap["value"] = v
			case []interface{}:
				lst := []string{}
				for _, val := range v {
					vStr := "'" + val.(string) + "'"
					lst = append(lst, vStr)
				}
				modelMap["value"] = "[" + strings.Join(lst, ",") + "]"
			}
		}
		// End Manual Intervention
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized securityandcompliancecenterapiv3.RequiredConfigIntf subtype encountered")
	}
}

func resourceIbmSccRuleRequiredConfigItemsToMap(model securityandcompliancecenterapiv3.RequiredConfigItemsIntf) (map[string]interface{}, error) {
	if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr); ok {
		return resourceIbmSccRuleRequiredConfigItemsRequiredConfigOrDepth1ToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd); ok {
		return resourceIbmSccRuleRequiredConfigItemsRequiredConfigAndDepth1ToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase); ok {
		return resourceIbmSccRuleRequiredConfigItemsRequiredConfigBaseToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItems); ok {
		modelMap := make(map[string]interface{})
		model := model.(*securityandcompliancecenterapiv3.RequiredConfigItems)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Or != nil {
			or := []map[string]interface{}{}
			for _, orItem := range model.Or {
				orItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(orItem)
				if err != nil {
					return modelMap, err
				}
				or = append(or, orItemMap)
			}
			modelMap["or"] = or
		}
		if model.And != nil {
			and := []map[string]interface{}{}
			for _, andItem := range model.And {
				andItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(andItem)
				if err != nil {
					return modelMap, err
				}
				and = append(and, andItemMap)
			}
			modelMap["and"] = and
		}
		if model.Property != nil {
			modelMap["property"] = model.Property
		}
		if model.Operator != nil {
			modelMap["operator"] = model.Operator
		}
		// Manual Intervention
		if model.Value != nil {
			// modelMap["value"] = model.Value
			switch v := model.Value.(type) {
			case string:
				modelMap["value"] = v
			case []interface{}:
				lst := []string{}
				for _, val := range v {
					vStr := val.(string)
					lst = append(lst, vStr)
				}
				modelMap["value"] = strings.Join(lst, ",")
			}
		}
		// End Manual Intervention
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized securityandcompliancecenterapiv3.RequiredConfigItemsIntf subtype encountered")
	}
}

func resourceIbmSccRuleRequiredConfigBaseToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["property"] = model.Property
	modelMap["operator"] = model.Operator
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigItemsRequiredConfigOrDepth1ToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Or != nil {
		or := []map[string]interface{}{}
		for _, orItem := range model.Or {
			orItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(orItem)
			if err != nil {
				return modelMap, err
			}
			or = append(or, orItemMap)
		}
		modelMap["or"] = or
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigItemsRequiredConfigAndDepth1ToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.And != nil {
		and := []map[string]interface{}{}
		for _, andItem := range model.And {
			andItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(andItem)
			if err != nil {
				return modelMap, err
			}
			and = append(and, andItemMap)
		}
		modelMap["and"] = and
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigItemsRequiredConfigBaseToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["property"] = model.Property
	modelMap["operator"] = model.Operator
	if model.Value != nil {
		// modelMap["value"] = model.Value
		switch v := model.Value.(type) {
		case string:
			modelMap["value"] = v
		case []string:
			s := strings.Join(v, ",")
			modelMap["value"] = s
		default:
			fmt.Printf("******** the type is %v\n", reflect.TypeOf(v))
		}
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigAndToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.And != nil {
		and := []map[string]interface{}{}
		for _, andItem := range model.And {
			andItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(andItem)
			if err != nil {
				return modelMap, err
			}
			and = append(and, andItemMap)
		}
		modelMap["and"] = and
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigOrToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Or != nil {
		or := []map[string]interface{}{}
		for _, orItem := range model.Or {
			orItemMap, err := resourceIbmSccRuleRequiredConfigItemsToMap(orItem)
			if err != nil {
				return modelMap, err
			}
			or = append(or, orItemMap)
		}
		modelMap["or"] = or
	}
	return modelMap, nil
}

func resourceIbmSccRuleRequiredConfigRequiredConfigBaseToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["property"] = model.Property
	modelMap["operator"] = model.Operator
	// Manual Intervention
	if model.Value != nil {
		// model.Value is a schema.TypeString, so it needs to converted to String Type
		switch v := model.Value.(type) {
		case string:
			modelMap["value"] = v
		case []interface{}:
			lst := []string{}
			for _, val := range v {
				vStr := "'" + val.(string) + "'"
				lst = append(lst, vStr)
			}
			modelMap["value"] = "[" + "'" + strings.Join(lst, ",") + "'" + "]"
		}
	}
	// End Manual Intervention

	return modelMap, nil
}
