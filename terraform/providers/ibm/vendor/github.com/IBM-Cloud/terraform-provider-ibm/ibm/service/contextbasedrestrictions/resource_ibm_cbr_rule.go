// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package contextbasedrestrictions

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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"x_correlation_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "x_correlation_id"),
				Description:  "The supplied or generated value of this header is logged for a request and repeated in a response header for the corresponding response. The same value is used for downstream requests and retries of those requests. If a value of this headers is not supplied in a request, the service generates a random (version 4) UUID.",
			},
			"transaction_id": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Deprecated:   "This argument is deprecated and may be removed in a future release",
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "transaction_id"),
				Description:  "The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.",
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "description"),
				Description:  "The description of the rule.",
			},
			"contexts": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
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
							Type:        schema.TypeSet,
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
				Computed:    true,
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
									"display_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
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
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cbr_rule", "enforcement_mode"),
				Description:  "The rule enforcement mode: * `enabled` - The restrictions are enforced and reported. This is the default. * `disabled` - The restrictions are disabled. Nothing is enforced or reported. * `report` - The restrictions are evaluated and reported, but not enforced.",
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
			"etag": &schema.Schema{
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cbr_rule", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCbrRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createRuleOptions := &contextbasedrestrictionsv1.CreateRuleOptions{}

	if _, ok := d.GetOk("description"); ok {
		createRuleOptions.SetDescription(d.Get("description").(string))
	}
	contexts := []contextbasedrestrictionsv1.RuleContext{}
	if _, ok := d.GetOk("contexts"); ok {
		for _, v := range d.Get("contexts").([]interface{}) {
			value := v.(map[string]interface{})
			contextsItem, err := ResourceIBMCbrRuleMapToRuleContext(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "create", "ResourceIBMCbrRuleMapToRuleContext").GetDiag()
			}
			contexts = append(contexts, *contextsItem)
		}
	}
	createRuleOptions.SetContexts(contexts)
	if _, ok := d.GetOk("resources"); ok {
		var resources []contextbasedrestrictionsv1.Resource
		for _, v := range d.Get("resources").([]interface{}) {
			value := v.(map[string]interface{})
			resourcesItem, err := ResourceIBMCbrRuleMapToResource(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "create", "ResourceIBMCbrRuleMapToResource").GetDiag()
			}
			resources = append(resources, *resourcesItem)
		}
		createRuleOptions.SetResources(resources)
	}
	if _, ok := d.GetOk("operations"); ok {
		operationsModel, err := ResourceIBMCbrRuleMapToNewRuleOperations(d.Get("operations.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "create", "ResourceIBMCbrRuleMapToNewRuleOperations").GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateRuleWithContext failed: %s", err.Error()), "ibm_cbr_rule", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*rule.ID)

	if err := ResourceIBMCbrRuleSetData(rule, response, d); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "create", "resourceIBMCbrRuleSetData").GetDiag()
	}

	return nil
}

func resourceIBMCbrRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getRuleOptions := &contextbasedrestrictionsv1.GetRuleOptions{}

	getRuleOptions.SetRuleID(d.Id())

	rule, response, err := contextBasedRestrictionsClient.GetRuleWithContext(context, getRuleOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRuleWithContext failed: %s", err.Error()), "ibm_cbr_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err := ResourceIBMCbrRuleSetData(rule, response, d); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "read", "resourceIBMCbrRuleSetData").GetDiag()
	}

	return nil
}

func resourceIBMCbrRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	replaceRuleOptions := &contextbasedrestrictionsv1.ReplaceRuleOptions{}

	replaceRuleOptions.SetRuleID(d.Id())
	if _, ok := d.GetOk("x_correlation_id"); ok {
		replaceRuleOptions.SetXCorrelationID(d.Get("x_correlation_id").(string))
	}
	if _, ok := d.GetOk("transaction_id"); ok {
		replaceRuleOptions.SetTransactionID(d.Get("transaction_id").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		replaceRuleOptions.SetDescription(d.Get("description").(string))
	}
	contexts := []contextbasedrestrictionsv1.RuleContext{}
	if _, ok := d.GetOk("contexts"); ok {
		for _, v := range d.Get("contexts").([]interface{}) {
			value := v.(map[string]interface{})
			contextsItem, err := ResourceIBMCbrRuleMapToRuleContext(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "update", "ResourceIBMCbrRuleMapToRuleContext").GetDiag()
			}
			contexts = append(contexts, *contextsItem)
		}
	}
	replaceRuleOptions.SetContexts(contexts)
	if _, ok := d.GetOk("resources"); ok {
		var resources []contextbasedrestrictionsv1.Resource
		for _, v := range d.Get("resources").([]interface{}) {
			value := v.(map[string]interface{})
			resourcesItem, err := ResourceIBMCbrRuleMapToResource(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "update", "ResourceIBMCbrRuleMapToResource").GetDiag()
			}
			resources = append(resources, *resourcesItem)
		}
		replaceRuleOptions.SetResources(resources)
	}
	if _, ok := d.GetOk("operations"); ok {
		operations, err := ResourceIBMCbrRuleMapToNewRuleOperations(d.Get("operations.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "update", "ResourceIBMCbrRuleMapToNewRuleOperations").GetDiag()
		}
		replaceRuleOptions.SetOperations(operations)
	}
	if _, ok := d.GetOk("enforcement_mode"); ok {
		replaceRuleOptions.SetEnforcementMode(d.Get("enforcement_mode").(string))
	}
	replaceRuleOptions.SetIfMatch(d.Get("etag").(string))

	rule, response, err := contextBasedRestrictionsClient.ReplaceRuleWithContext(context, replaceRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceRuleWithContext failed: %s", err.Error()), "ibm_cbr_rule", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err := ResourceIBMCbrRuleSetData(rule, response, d); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "update", "resourceIBMCbrRuleSetData").GetDiag()
	}

	return nil
}

func resourceIBMCbrRuleDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cbr_rule", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteRuleOptions := &contextbasedrestrictionsv1.DeleteRuleOptions{}

	deleteRuleOptions.SetRuleID(d.Id())

	_, err = contextBasedRestrictionsClient.DeleteRuleWithContext(context, deleteRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteRuleWithContext failed: %s", err.Error()), "ibm_cbr_rule", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMCbrRuleSetData(rule *contextbasedrestrictionsv1.Rule, response *core.DetailedResponse, d *schema.ResourceData) error {
	if !core.IsNil(rule.Description) {
		if err := d.Set("description", rule.Description); err != nil {
			return fmt.Errorf("Error setting description: %s", err)
		}
	}
	contexts := []map[string]interface{}{}
	if !core.IsNil(rule.Contexts) {
		for _, contextsItem := range rule.Contexts {
			contextsItemMap, err := ResourceIBMCbrRuleRuleContextToMap(&contextsItem) // #nosec G601
			if err != nil {
				return fmt.Errorf("Error map rule context: %s", err)
			}
			contexts = append(contexts, contextsItemMap)
		}
	}
	if err := d.Set("contexts", contexts); err != nil {
		return fmt.Errorf("Error setting contexts: %s", err)
	}
	if !core.IsNil(rule.Resources) {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range rule.Resources {
			resourcesItemMap, err := ResourceIBMCbrRuleResourceToMap(&resourcesItem) // #nosec G601
			if err != nil {
				return fmt.Errorf("Error map rule resource: %s", err)
			}
			resources = append(resources, resourcesItemMap)
		}
		if err := d.Set("resources", resources); err != nil {
			return fmt.Errorf("Error setting resources: %s", err)
		}
	}
	if !core.IsNil(rule.Operations) {
		operationsMap, err := ResourceIBMCbrRuleNewRuleOperationsToMap(rule.Operations)
		if err != nil {
			return fmt.Errorf("Error map rule operations: %s", err)
		}
		if err = d.Set("operations", []map[string]interface{}{operationsMap}); err != nil {
			return fmt.Errorf("Error setting operations: %s", err)
		}
	}
	if !core.IsNil(rule.EnforcementMode) {
		if err := d.Set("enforcement_mode", rule.EnforcementMode); err != nil {
			return fmt.Errorf("Error setting enforcement_mode: %s", err)
		}
	}
	if err := d.Set("crn", rule.CRN); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err := d.Set("href", rule.Href); err != nil {
		return fmt.Errorf("Error setting href: %s", err)
	}
	if err := d.Set("created_at", flex.DateTimeToString(rule.CreatedAt)); err != nil {
		return fmt.Errorf("Error setting created_at: %s", err)
	}
	if err := d.Set("created_by_id", rule.CreatedByID); err != nil {
		return fmt.Errorf("Error setting created_by_id: %s", err)
	}
	if err := d.Set("last_modified_at", flex.DateTimeToString(rule.LastModifiedAt)); err != nil {
		return fmt.Errorf("Error setting last_modified_at: %s", err)
	}
	if err := d.Set("last_modified_by_id", rule.LastModifiedByID); err != nil {
		return fmt.Errorf("Error setting last_modified_by_id: %s", err)
	}
	if err := d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return fmt.Errorf("Error setting etag: %s", err)
	}
	return nil
}

func ResourceIBMCbrRuleMapToRuleContext(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.RuleContext, error) {
	model := &contextbasedrestrictionsv1.RuleContext{}
	attributes := []contextbasedrestrictionsv1.RuleContextAttribute{}
	for _, attributesItem := range modelMap["attributes"].([]interface{}) {
		attributesItemModel, err := ResourceIBMCbrRuleMapToRuleContextAttribute(attributesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		attributes = append(attributes, *attributesItemModel)
	}
	model.Attributes = attributes
	return model, nil
}

func ResourceIBMCbrRuleMapToRuleContextAttribute(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.RuleContextAttribute, error) {
	model := &contextbasedrestrictionsv1.RuleContextAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIBMCbrRuleMapToResource(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.Resource, error) {
	model := &contextbasedrestrictionsv1.Resource{}
	attributes := []contextbasedrestrictionsv1.ResourceAttribute{}
	attributeList := modelMap["attributes"].(*schema.Set).List()
	// for _, attributesItem := range modelMap["attributes"].([]interface{}) {
	for _, attributesItem := range attributeList {
		attributesItemModel, err := ResourceIBMCbrRuleMapToResourceAttribute(attributesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		attributes = append(attributes, *attributesItemModel)
	}
	model.Attributes = attributes
	if modelMap["tags"] != nil {
		tags := []contextbasedrestrictionsv1.ResourceTagAttribute{}
		for _, tagsItem := range modelMap["tags"].([]interface{}) {
			tagsItemModel, err := ResourceIBMCbrRuleMapToResourceTagAttribute(tagsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			tags = append(tags, *tagsItemModel)
		}
		model.Tags = tags
	}
	return model, nil
}

func ResourceIBMCbrRuleMapToResourceAttribute(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.ResourceAttribute, error) {
	model := &contextbasedrestrictionsv1.ResourceAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	return model, nil
}

func ResourceIBMCbrRuleMapToResourceTagAttribute(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.ResourceTagAttribute, error) {
	model := &contextbasedrestrictionsv1.ResourceTagAttribute{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	if modelMap["operator"] != nil && modelMap["operator"].(string) != "" {
		model.Operator = core.StringPtr(modelMap["operator"].(string))
	}
	return model, nil
}

func ResourceIBMCbrRuleMapToNewRuleOperations(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.NewRuleOperations, error) {
	model := &contextbasedrestrictionsv1.NewRuleOperations{}
	apiTypes := []contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem{}
	for _, apiTypesItem := range modelMap["api_types"].([]interface{}) {
		apiTypesItemModel, err := ResourceIBMCbrRuleMapToNewRuleOperationsAPITypesItem(apiTypesItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		apiTypes = append(apiTypes, *apiTypesItemModel)
	}
	model.APITypes = apiTypes
	return model, nil
}

func ResourceIBMCbrRuleMapToNewRuleOperationsAPITypesItem(modelMap map[string]interface{}) (*contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem, error) {
	model := &contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem{}
	model.APITypeID = core.StringPtr(modelMap["api_type_id"].(string))
	return model, nil
}

func ResourceIBMCbrRuleRuleContextToMap(model *contextbasedrestrictionsv1.RuleContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	attributes := []map[string]interface{}{}
	for _, attributesItem := range model.Attributes {
		attributesItemMap, err := ResourceIBMCbrRuleRuleContextAttributeToMap(&attributesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		attributes = append(attributes, attributesItemMap)
	}
	modelMap["attributes"] = attributes
	return modelMap, nil
}

func ResourceIBMCbrRuleRuleContextAttributeToMap(model *contextbasedrestrictionsv1.RuleContextAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	return modelMap, nil
}

func compareResAttrSetFunc(v interface{}) int {
	m := v.(map[string]interface{})
	name := (m["name"]).(*string)
	return schema.HashString(*name)
}

func ResourceIBMCbrRuleResourceToMap(model *contextbasedrestrictionsv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	attributes := []interface{}{}
	//attributes := []map[string]interface{}{}
	for _, attributesItem := range model.Attributes {
		attributesItemMap, err := ResourceIBMCbrRuleResourceAttributeToMap(&attributesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		attributes = append(attributes, attributesItemMap)
	}
	attributeSet := schema.NewSet(compareResAttrSetFunc, attributes)
	modelMap["attributes"] = attributeSet
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := ResourceIBMCbrRuleResourceTagAttributeToMap(&tagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tags = append(tags, tagsItemMap)
		}
		modelMap["tags"] = tags
	}
	return modelMap, nil
}

func ResourceIBMCbrRuleResourceAttributeToMap(model *contextbasedrestrictionsv1.ResourceAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	return modelMap, nil
}

func ResourceIBMCbrRuleResourceTagAttributeToMap(model *contextbasedrestrictionsv1.ResourceTagAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	if model.Operator != nil {
		modelMap["operator"] = model.Operator
	}
	return modelMap, nil
}

func ResourceIBMCbrRuleNewRuleOperationsToMap(model *contextbasedrestrictionsv1.NewRuleOperations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	apiTypes := []map[string]interface{}{}
	for _, apiTypesItem := range model.APITypes {
		apiTypesItemMap, err := ResourceIBMCbrRuleNewRuleOperationsAPITypesItemToMap(&apiTypesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		apiTypes = append(apiTypes, apiTypesItemMap)
	}
	modelMap["api_types"] = apiTypes
	return modelMap, nil
}

func ResourceIBMCbrRuleNewRuleOperationsAPITypesItemToMap(model *contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["api_type_id"] = model.APITypeID
	return modelMap, nil
}
