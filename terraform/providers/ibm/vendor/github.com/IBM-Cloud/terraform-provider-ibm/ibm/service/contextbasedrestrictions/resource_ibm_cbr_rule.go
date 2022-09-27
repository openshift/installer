// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package contextbasedrestrictions

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
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func ResourceIBMCbrRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCbrRuleCreate,
		ReadContext:   resourceIBMCbrRuleRead,
		UpdateContext: resourceIBMCbrRuleUpdate,
		DeleteContext: resourceIBMCbrRuleDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "description"),
				Description:  "The description of the rule.",
			},
			"contexts": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The contexts this rule applies to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The attribute name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The attribute value.",
									},
								},
							},
						},
					},
				},
			},
			"resources": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The resources this rule apply to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The resource attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The attribute name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The attribute value.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The attribute operator.",
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The optional resource tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The tag attribute name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The tag attribute value.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The attribute operator.",
									},
								},
							},
						},
					},
				},
			},
			"operations": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The operations this rule applies to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The API types this rule applies to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_type_id": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"enforcement_mode": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "enforcement_mode"),
				Description:  "The rule enforcement mode: * `enabled` - The restrictions are enforced and reported. This is the default. * `disabled` - The restrictions are disabled. Nothing is enforced or reported. * `report` - The restrictions are evaluated and reported, but not enforced.",
			},
			"x_correlation_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "x_correlation_id"),
				Description:  "The supplied or generated value of this header is logged for a request and repeated in a response header for the corresponding response. The same value is used for downstream requests and retries of those requests. If a value of this headers is not supplied in a request, the service generates a random (version 4) UUID.",
			},
			"transaction_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "transaction_id"),
				Description:  "The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule CRN.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The href link to the resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time the resource was created.",
			},
			"created_by_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which created the resource.",
			},
			"last_modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last time the resource was modified.",
			},
			"last_modified_by_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAM ID of the user or service which modified the resource.",
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMCbrRuleValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[\x20-\xFE]*$`,
			MinValueLength:             0,
			MaxValueLength:             300,
		},
		validate.ValidateSchema{
			Identifier:                 "enforcement_mode",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "disabled, enabled, report",
		},
		validate.ValidateSchema{
			Identifier:                 "x_correlation_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 ,\-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             1024,
		},
		validate.ValidateSchema{
			Identifier:                 "transaction_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9 ,\-_]+$`,
			MinValueLength:             1,
			MaxValueLength:             1024,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cbr_rule", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCbrRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createRuleOptions := &contextbasedrestrictionsv1.CreateRuleOptions{}

	if _, ok := d.GetOk("description"); ok {
		createRuleOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("contexts"); ok {
		var contexts []contextbasedrestrictionsv1.RuleContext
		for _, e := range d.Get("contexts").([]interface{}) {
			value := e.(map[string]interface{})
			contextsItem, err := resourceIBMCbrRuleMapToRuleContext(value)
			if err != nil {
				return diag.FromErr(err)
			}
			contexts = append(contexts, *contextsItem)
		}
		createRuleOptions.SetContexts(contexts)
	}
	if _, ok := d.GetOk("resources"); ok {
		var resources []contextbasedrestrictionsv1.Resource
		for _, e := range d.Get("resources").([]interface{}) {
			value := e.(map[string]interface{})
			resourcesItem, err := resourceIBMCbrRuleMapToResource(value)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, *resourcesItem)
		}
		createRuleOptions.SetResources(resources)
	}
	if _, ok := d.GetOk("operations"); ok {
		operationsModel, err := resourceIBMCbrRuleMapToNewRuleOperations(d.Get("operations.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createRuleOptions.SetOperations(operationsModel)
	}
	if _, ok := d.GetOk("enforcement_mode"); ok {
		createRuleOptions.SetEnforcementMode(d.Get("enforcement_mode").(string))
	}
	if _, ok := d.GetOk("x_correlation_id"); ok {
		createRuleOptions.SetXCorrelationID(d.Get("x_correlation_id").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		createRuleOptions.SetTransactionID(d.Get("transaction_id").(string))
	}

	rule, response, err := contextBasedRestrictionsClient.CreateRuleWithContext(context, createRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId(*rule.ID)

	return resourceIBMCbrRuleRead(context, d, meta)
}

func resourceIBMCbrRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleOptions := &contextbasedrestrictionsv1.GetRuleOptions{}

	getRuleOptions.SetRuleID(d.Id())

	rule, response, err := contextBasedRestrictionsClient.GetRuleWithContext(context, getRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRuleWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("x_correlation_id", getRuleOptions.XCorrelationID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting x_correlation_id: %s", err))
	}
	if err = d.Set("transaction_id", getRuleOptions.TransactionID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting transaction_id: %s", err))
	}
	if err = d.Set("description", rule.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	contexts := []map[string]interface{}{}
	if rule.Contexts != nil {
		for _, contextsItem := range rule.Contexts {
			contextsItemMap, err := resourceIBMCbrRuleRuleContextToMap(&contextsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			contexts = append(contexts, contextsItemMap)
		}
	}
	if err = d.Set("contexts", contexts); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting contexts: %s", err))
	}
	resources := []map[string]interface{}{}
	if rule.Resources != nil {
		for _, resourcesItem := range rule.Resources {
			resourcesItemMap, err := resourceIBMCbrRuleResourceToMap(&resourcesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, resourcesItemMap)
		}
	}
	if err = d.Set("resources", resources); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resources: %s", err))
	}
	if rule.Operations != nil {
		operationsMap, err := resourceIBMCbrRuleNewRuleOperationsToMap(rule.Operations)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("operations", []map[string]interface{}{operationsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting operations: %s", err))
		}
	}
	if err = d.Set("enforcement_mode", rule.EnforcementMode); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enforcement_mode: %s", err))
	}
	if err = d.Set("crn", rule.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", rule.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(rule.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by_id", rule.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}
	if err = d.Set("last_modified_at", flex.DateTimeToString(rule.LastModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_at: %s", err))
	}
	if err = d.Set("last_modified_by_id", rule.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_by_id: %s", err))
	}
	if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	return nil
}

func resourceIBMCbrRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceRuleOptions := &contextbasedrestrictionsv1.ReplaceRuleOptions{}

	replaceRuleOptions.SetRuleID(d.Id())
	if _, ok := d.GetOk("description"); ok {
		replaceRuleOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("contexts"); ok {
		var contexts []contextbasedrestrictionsv1.RuleContext
		for _, e := range d.Get("contexts").([]interface{}) {
			value := e.(map[string]interface{})
			contextsItem, err := resourceIBMCbrRuleMapToRuleContext(value)
			if err != nil {
				return diag.FromErr(err)
			}
			contexts = append(contexts, *contextsItem)
		}
		replaceRuleOptions.SetContexts(contexts)
	}
	if _, ok := d.GetOk("resources"); ok {
		var resources []contextbasedrestrictionsv1.Resource
		for _, e := range d.Get("resources").([]interface{}) {
			value := e.(map[string]interface{})
			resourcesItem, err := resourceIBMCbrRuleMapToResource(value)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, *resourcesItem)
		}
		replaceRuleOptions.SetResources(resources)
	}
	if _, ok := d.GetOk("operations"); ok {
		operations, err := resourceIBMCbrRuleMapToNewRuleOperations(d.Get("operations.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceRuleOptions.SetOperations(operations)
	}
	if _, ok := d.GetOk("enforcement_mode"); ok {
		replaceRuleOptions.SetEnforcementMode(d.Get("enforcement_mode").(string))
	}
	if _, ok := d.GetOk("x_correlation_id"); ok {
		replaceRuleOptions.SetXCorrelationID(d.Get("x_correlation_id").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		replaceRuleOptions.SetTransactionID(d.Get("transaction_id").(string))
	}
	replaceRuleOptions.SetIfMatch(d.Get("version").(string))

	_, response, err := contextBasedRestrictionsClient.ReplaceRuleWithContext(context, replaceRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] ReplaceRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ReplaceRuleWithContext failed %s\n%s", err, response))
	}

	return resourceIBMCbrRuleRead(context, d, meta)
}

func resourceIBMCbrRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRuleOptions := &contextbasedrestrictionsv1.DeleteRuleOptions{}

	deleteRuleOptions.SetRuleID(d.Id())

	response, err := contextBasedRestrictionsClient.DeleteRuleWithContext(context, deleteRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCbrRuleMapToRuleContext(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.RuleContext, error) {
	model := &contextbasedrestrictionsv1.RuleContext{}
	attributes := []contextbasedrestrictionsv1.RuleContextAttribute{}
	for _, attributesItem := range modelMap["attributes"].([]interface{}) {
		attributesItemModel, err := resourceIBMCbrRuleMapToRuleContextAttribute(attributesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		attributes = append(attributes, *attributesItemModel)
	}
	model.Attributes = attributes
	return model, nil
}

func resourceIBMCbrRuleMapToRuleContextAttribute(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.RuleContextAttribute, error) {
	model := &contextbasedrestrictionsv1.RuleContextAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIBMCbrRuleMapToResource(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.Resource, error) {
	model := &contextbasedrestrictionsv1.Resource{}
	attributes := []contextbasedrestrictionsv1.ResourceAttribute{}
	for _, attributesItem := range modelMap["attributes"].([]interface{}) {
		attributesItemModel, err := resourceIBMCbrRuleMapToResourceAttribute(attributesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		attributes = append(attributes, *attributesItemModel)
	}
	model.Attributes = attributes
	if modelMap["tags"] != nil {
		tags := []contextbasedrestrictionsv1.ResourceTagAttribute{}
		for _, tagsItem := range modelMap["tags"].([]interface{}) {
			tagsItemModel, err := resourceIBMCbrRuleMapToResourceTagAttribute(tagsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			tags = append(tags, *tagsItemModel)
		}
		model.Tags = tags
	}
	return model, nil
}

func resourceIBMCbrRuleMapToResourceAttribute(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.ResourceAttribute, error) {
	model := &contextbasedrestrictionsv1.ResourceAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	return model, nil
}

func resourceIBMCbrRuleMapToResourceTagAttribute(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.ResourceTagAttribute, error) {
	model := &contextbasedrestrictionsv1.ResourceTagAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	return model, nil
}

func resourceIBMCbrRuleMapToNewRuleOperations(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.NewRuleOperations, error) {
	model := &contextbasedrestrictionsv1.NewRuleOperations{}
	apiTypes := []contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem{}
	for _, apiTypesItem := range modelMap["api_types"].([]interface{}) {
		apiTypesItemModel, err := resourceIBMCbrRuleMapToNewRuleOperationsAPITypesItem(apiTypesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		apiTypes = append(apiTypes, *apiTypesItemModel)
	}
	model.APITypes = apiTypes
	return model, nil
}

func resourceIBMCbrRuleMapToNewRuleOperationsAPITypesItem(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem, error) {
	model := &contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem{}
	model.APITypeID = core.StringPtr(modelMap["api_type_id"].(string))
	return model, nil
}

func resourceIBMCbrRuleRuleContextToMap(model *contextbasedrestrictionsv1.RuleContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	attributes := []map[string]interface{}{}
	for _, attributesItem := range model.Attributes {
		attributesItemMap, err := resourceIBMCbrRuleRuleContextAttributeToMap(&attributesItem)
		if err != nil {
			return modelMap, err
		}
		attributes = append(attributes, attributesItemMap)
	}
	modelMap["attributes"] = attributes
	return modelMap, nil
}

func resourceIBMCbrRuleRuleContextAttributeToMap(model *contextbasedrestrictionsv1.RuleContextAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIBMCbrRuleResourceToMap(model *contextbasedrestrictionsv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	attributes := []map[string]interface{}{}
	for _, attributesItem := range model.Attributes {
		attributesItemMap, err := resourceIBMCbrRuleResourceAttributeToMap(&attributesItem)
		if err != nil {
			return modelMap, err
		}
		attributes = append(attributes, attributesItemMap)
	}
	modelMap["attributes"] = attributes
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := resourceIBMCbrRuleResourceTagAttributeToMap(&tagsItem)
			if err != nil {
				return modelMap, err
			}
			tags = append(tags, tagsItemMap)
		}
		modelMap["tags"] = tags
	}
	return modelMap, nil
}

func resourceIBMCbrRuleResourceAttributeToMap(model *contextbasedrestrictionsv1.ResourceAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	return modelMap, nil
}

func resourceIBMCbrRuleResourceTagAttributeToMap(model *contextbasedrestrictionsv1.ResourceTagAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	return modelMap, nil
}

func resourceIBMCbrRuleNewRuleOperationsToMap(model *contextbasedrestrictionsv1.NewRuleOperations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	apiTypes := []map[string]interface{}{}
	for _, apiTypesItem := range model.APITypes {
		apiTypesItemMap, err := resourceIBMCbrRuleNewRuleOperationsAPITypesItemToMap(&apiTypesItem)
		if err != nil {
			return modelMap, err
		}
		apiTypes = append(apiTypes, apiTypesItemMap)
	}
	modelMap["api_types"] = apiTypes
	return modelMap, nil
}

func resourceIBMCbrRuleNewRuleOperationsAPITypesItemToMap(model *contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["api_type_id"] = model.APITypeID
	return modelMap, nil
}
