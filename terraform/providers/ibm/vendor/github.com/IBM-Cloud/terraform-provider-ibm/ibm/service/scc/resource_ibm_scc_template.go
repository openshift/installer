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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v3/configurationgovernancev1"
)

func ResourceIBMSccTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSccTemplateCreate,
		ReadContext:   resourceIBMSccTemplateRead,
		UpdateContext: resourceIBMSccTemplateUpdate,
		DeleteContext: resourceIBMSccTemplateDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Your IBM Cloud account ID.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A human-readablse alias to assign to your template.",
				ValidateFunc: validate.InvokeValidator("ibm_scc_template", "name"),
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "An extended description of your template.",
				ValidateFunc: validate.InvokeValidator("ibm_scc_template", "description"),
			},
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID that uniquely identifies the template.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The properties that describe the resource that you want to targetwith the rule or template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": &schema.Schema{
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The programmatic name of the IBM Cloud service that you want to target with the rule or template.",
							ValidateFunc: validate.InvokeValidator("ibm_scc_template", "service_name"),
						},
						"resource_kind": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of resource that you want to target.",
						},
						"additional_target_attributes": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "An extra qualifier for the resource kind. When you include additional attributes, only the resources that match the definition are included in the rule or template.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the additional attribute that you want to use to further qualify the target.Options differ depending on the service or resource that you are targeting with a rule or template. For more information, refer to the service documentation.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The value that you want to apply to `name` field.Options differ depending on the rule or template that you configure. For more information, refer to the service documentation.",
									},
								},
							},
						},
					},
				},
			},
			"customized_defaults": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "A list of default property values to apply to your template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"property": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the resource property that you want to configure.Property options differ depending on the service or resource that you are targeting with a template. To view a list of properties that are compatible with templates, refer to the service documentation.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The custom value that you want to apply as the default for the resource property in the `name` field.This value is used to to override the default value that is provided by IBM when a resource is created. Value options differ depending on the resource that you are configuring. To learn more about your options, refer to the service documentation.",
						},
					},
				},
			},
		},
	}
}

func resourceIBMSccTemplateCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createTemplatesOptions := &configurationgovernancev1.CreateTemplatesOptions{}

	var template []configurationgovernancev1.CreateTemplateRequest
	templateItem, err := resourceIBMSccTemplateMapToCreateTemplateRequest(d)
	template = append(template, *templateItem)
	createTemplatesOptions.SetTemplates(template)

	createTemplatesResponse, response, err := configurationGovernanceClient.CreateTemplatesWithContext(context, createTemplatesOptions)
	if err != nil || response.GetStatusCode() == 207 || response.StatusCode > 300 {
		log.Printf("[DEBUG] CreateTemplatesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTemplatesWithContext failed %s\n%s", err, response))
	}

	d.SetId(*createTemplatesResponse.Templates[0].Template.TemplateID)

	return resourceIBMSccTemplateRead(context, d, meta)
}

func resourceIBMSccTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getTemplateOptions := &configurationgovernancev1.GetTemplateOptions{}

	getTemplateOptions.SetTemplateID(d.Id())

	templateResponse, response, err := configurationGovernanceClient.GetTemplateWithContext(context, getTemplateOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTemplateWithContext failed %s\n%s", err, response))
	}

	// TODO: handle argument of type []interface{}
	if err = d.Set("account_id", templateResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("name", templateResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("description", templateResponse.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	targetMap, err := resourceIBMSccTemplateSimpleTargetResourceToMap(templateResponse.Target)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("target", []map[string]interface{}{targetMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target: %s", err))
	}

	customizedDefaults := []map[string]interface{}{}
	for _, customizedDefaultsItem := range templateResponse.CustomizedDefaults {
		customizedDefaultsItemMap, err := resourceIBMSccTemplateTemplateCustomizedDefaultPropertyToMap(&customizedDefaultsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		customizedDefaults = append(customizedDefaults, customizedDefaultsItemMap)
	}
	if err = d.Set("customized_defaults", customizedDefaults); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting customized_defaults: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	return nil
}

func resourceIBMSccTemplateUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateTemplateOptions := &configurationgovernancev1.UpdateTemplateOptions{}

	updateTemplateOptions.SetTemplateID(d.Id())

	hasChange := d.HasChange("name") || d.HasChange("description") ||
		d.HasChange("target") || d.HasChange("customized_defaults")

	if hasChange {
		updateTemplateOptions.SetIfMatch(d.Get("version").(string))
		updateTemplateOptions.SetName(d.Get("name").(string))
		updateTemplateOptions.SetAccountID(d.Get("account_id").(string))
		updateTemplateOptions.SetDescription(d.Get("description").(string))

		target, err := resourceIBMSccTemplateMapToSimpleTargetResource(d.Get("target.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateTemplateOptions.SetTarget(target)

		customizedDefaults := []configurationgovernancev1.TemplateCustomizedDefaultProperty{}
		for _, customizedDefaultsItem := range d.Get("customized_defaults").([]interface{}) {
			if customizedDefaultsItem != nil {
				customizedDefaultsItemModel, err := resourceIBMSccTemplateMapToTemplateCustomizedDefaultProperty(customizedDefaultsItem.(map[string]interface{}))
				if err != nil {
					return diag.FromErr(err)
				}
				customizedDefaults = append(customizedDefaults, *customizedDefaultsItemModel)
			}
		}
		updateTemplateOptions.SetCustomizedDefaults(customizedDefaults)

		_, response, err := configurationGovernanceClient.UpdateTemplateWithContext(context, updateTemplateOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTemplateWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTemplateWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMSccTemplateRead(context, d, meta)
}

func resourceIBMSccTemplateDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configurationGovernanceClient, err := meta.(conns.ClientSession).ConfigurationGovernanceV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTemplateOptions := &configurationgovernancev1.DeleteTemplateOptions{}

	deleteTemplateOptions.SetTemplateID(d.Id())

	response, err := configurationGovernanceClient.DeleteTemplateWithContext(context, deleteTemplateOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTemplateWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTemplateWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMSccTemplateMapToCreateTemplateRequest(d *schema.ResourceData) (*configurationgovernancev1.CreateTemplateRequest, error) {
	model := &configurationgovernancev1.CreateTemplateRequest{}
	if d.Get("request_id") != nil {
		model.RequestID = core.StringPtr(d.Get("request_id").(string))
	}
	TemplateModel, err := resourceIBMSccTemplateMapToTemplate(d)
	if err != nil {
		return model, err
	}
	model.Template = TemplateModel
	return model, nil
}

func resourceIBMSccTemplateMapToTemplate(d *schema.ResourceData) (*configurationgovernancev1.Template, error) {
	model := &configurationgovernancev1.Template{}
	model.AccountID = core.StringPtr(d.Get("account_id").(string))
	model.Name = core.StringPtr(d.Get("name").(string))
	model.Description = core.StringPtr(d.Get("description").(string))
	if d.Get("template_id") != nil {
		model.TemplateID = core.StringPtr(d.Get("template_id").(string))
	}
	targetList := d.Get("target").([]interface{})
	TargetModel, err := resourceIBMSccTemplateMapToSimpleTargetResource(targetList[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Target = TargetModel
	customizedDefaults := []configurationgovernancev1.TemplateCustomizedDefaultProperty{}
	for _, customizedDefaultsItem := range d.Get("customized_defaults").([]interface{}) {
		customizedDefaultsItemModel, err := resourceIBMSccTemplateMapToTemplateCustomizedDefaultProperty(customizedDefaultsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		customizedDefaults = append(customizedDefaults, *customizedDefaultsItemModel)
	}
	model.CustomizedDefaults = customizedDefaults
	return model, nil
}

func resourceIBMSccTemplateMapToSimpleTargetResource(modelMap map[string]interface{}) (*configurationgovernancev1.SimpleTargetResource, error) {
	model := &configurationgovernancev1.SimpleTargetResource{}
	model.ServiceName = core.StringPtr(modelMap["service_name"].(string))
	model.ResourceKind = core.StringPtr(modelMap["resource_kind"].(string))
	if modelMap["additional_target_attributes"] != nil {
		additionalTargetAttributes := []configurationgovernancev1.BaseTargetAttribute{}
		for _, additionalTargetAttributesItem := range modelMap["additional_target_attributes"].([]interface{}) {
			additionalTargetAttributesItemModel, err := resourceIBMSccTemplateMapToBaseTargetAttribute(additionalTargetAttributesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, *additionalTargetAttributesItemModel)
		}
		model.AdditionalTargetAttributes = additionalTargetAttributes
	}
	return model, nil
}

func resourceIBMSccTemplateMapToBaseTargetAttribute(modelMap map[string]interface{}) (*configurationgovernancev1.BaseTargetAttribute, error) {
	model := &configurationgovernancev1.BaseTargetAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMSccTemplateMapToTemplateCustomizedDefaultProperty(modelMap map[string]interface{}) (*configurationgovernancev1.TemplateCustomizedDefaultProperty, error) {
	model := &configurationgovernancev1.TemplateCustomizedDefaultProperty{}
	model.Property = core.StringPtr(modelMap["property"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMSccTemplateCreateTemplateRequestToMap(model *configurationgovernancev1.CreateTemplateRequest) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RequestID != nil {
		modelMap["request_id"] = model.RequestID
	}
	templateMap, err := resourceIBMSccTemplateTemplateToMap(model.Template)
	if err != nil {
		return modelMap, err
	}
	modelMap["template"] = []map[string]interface{}{templateMap}
	return modelMap, nil
}

func resourceIBMSccTemplateTemplateToMap(model *configurationgovernancev1.Template) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["account_id"] = model.AccountID
	modelMap["name"] = model.Name
	modelMap["description"] = model.Description
	if model.TemplateID != nil {
		modelMap["template_id"] = model.TemplateID
	}
	targetMap, err := resourceIBMSccTemplateSimpleTargetResourceToMap(model.Target)
	if err != nil {
		return modelMap, err
	}
	modelMap["target"] = []map[string]interface{}{targetMap}
	customizedDefaults := []map[string]interface{}{}
	for _, customizedDefaultsItem := range model.CustomizedDefaults {
		customizedDefaultsItemMap, err := resourceIBMSccTemplateTemplateCustomizedDefaultPropertyToMap(&customizedDefaultsItem)
		if err != nil {
			return modelMap, err
		}
		customizedDefaults = append(customizedDefaults, customizedDefaultsItemMap)
	}
	modelMap["customized_defaults"] = customizedDefaults
	return modelMap, nil
}

func resourceIBMSccTemplateSimpleTargetResourceToMap(model *configurationgovernancev1.SimpleTargetResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["service_name"] = model.ServiceName
	modelMap["resource_kind"] = model.ResourceKind
	if model.AdditionalTargetAttributes != nil {
		additionalTargetAttributes := []map[string]interface{}{}
		for _, additionalTargetAttributesItem := range model.AdditionalTargetAttributes {
			additionalTargetAttributesItemMap, err := resourceIBMSccTemplateBaseTargetAttributeToMap(&additionalTargetAttributesItem)
			if err != nil {
				return modelMap, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, additionalTargetAttributesItemMap)
		}
		modelMap["additional_target_attributes"] = additionalTargetAttributes
	}
	return modelMap, nil
}

func resourceIBMSccTemplateBaseTargetAttributeToMap(model *configurationgovernancev1.BaseTargetAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIBMSccTemplateTemplateCustomizedDefaultPropertyToMap(model *configurationgovernancev1.TemplateCustomizedDefaultProperty) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["property"] = model.Property
	modelMap["value"] = model.Value
	return modelMap, nil
}
