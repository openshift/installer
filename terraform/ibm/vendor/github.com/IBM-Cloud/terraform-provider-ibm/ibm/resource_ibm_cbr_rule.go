// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func resourceIBMCbrRule() *schema.Resource {
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
				ValidateFunc: InvokeValidator("ibm_cbr_rule", "description"),
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

func resourceIBMCbrRuleValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "description",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^[\x20-\xFE]*$`,
			MinValueLength:             0,
			MaxValueLength:             300,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_cbr_rule", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCbrRuleCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createRuleOptions := &contextbasedrestrictionsv1.CreateRuleOptions{}

	accountID, err := getIBMCbrAccountId(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk("description"); ok {
		createRuleOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("contexts"); ok {
		var contexts []contextbasedrestrictionsv1.RuleContext
		for _, e := range d.Get("contexts").([]interface{}) {
			value := e.(map[string]interface{})
			contextsItem := resourceIBMCbrRuleMapToRuleContext(value)
			contexts = append(contexts, contextsItem)
		}
		createRuleOptions.SetContexts(contexts)
	}
	if _, ok := d.GetOk("resources"); ok {
		var resources []contextbasedrestrictionsv1.Resource
		for _, e := range d.Get("resources").([]interface{}) {
			value := e.(map[string]interface{})
			resourcesItem := resourceIBMCbrRuleMapToResource(value, accountID)
			resources = append(resources, resourcesItem)
		}

		createRuleOptions.SetResources(resources)
	}

	rule, response, err := contextBasedRestrictionsClient.CreateRuleWithContext(context, createRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId(*rule.ID)

	return resourceIBMCbrRuleRead(context, d, meta)
}

func resourceIBMCbrRuleMapToRuleContext(ruleContextMap map[string]interface{}) contextbasedrestrictionsv1.RuleContext {
	ruleContext := contextbasedrestrictionsv1.RuleContext{}

	attributes := []contextbasedrestrictionsv1.RuleContextAttribute{}
	for _, attributesItem := range ruleContextMap["attributes"].([]interface{}) {
		attributesItemModel := resourceIBMCbrRuleMapToRuleContextAttribute(attributesItem.(map[string]interface{}))
		attributes = append(attributes, attributesItemModel)
	}
	ruleContext.Attributes = attributes

	return ruleContext
}

func resourceIBMCbrRuleMapToRuleContextAttribute(ruleContextAttributeMap map[string]interface{}) contextbasedrestrictionsv1.RuleContextAttribute {
	ruleContextAttribute := contextbasedrestrictionsv1.RuleContextAttribute{}

	ruleContextAttribute.Name = core.StringPtr(ruleContextAttributeMap["name"].(string))
	ruleContextAttribute.Value = core.StringPtr(ruleContextAttributeMap["value"].(string))
	return ruleContextAttribute
}

func resourceIBMCbrRuleAccountIdAttribute(accountID string) contextbasedrestrictionsv1.ResourceAttribute {
	accountIdAttribute := contextbasedrestrictionsv1.ResourceAttribute{}

	accountIdAttribute.Name = core.StringPtr("accountId")
	accountIdAttribute.Value = core.StringPtr(accountID)
	//accountIdAttribute.Operator = core.StringPtr("")

	return accountIdAttribute

}

func resourceIBMCbrRuleMapToResource(resourceMap map[string]interface{}, accountID string) contextbasedrestrictionsv1.Resource {
	resource := contextbasedrestrictionsv1.Resource{}

	attributes := []contextbasedrestrictionsv1.ResourceAttribute{}

	attributes = append(attributes, resourceIBMCbrRuleAccountIdAttribute(accountID))

	for _, attributesItem := range resourceMap["attributes"].([]interface{}) {
		attributesItemModel := resourceIBMCbrRuleMapToResourceAttribute(attributesItem.(map[string]interface{}))

		if *attributesItemModel.Name != "accountId" {
			attributes = append(attributes, attributesItemModel)
		}
	}

	resource.Attributes = attributes

	if resourceMap["tags"] != nil {
		tags := []contextbasedrestrictionsv1.ResourceTagAttribute{}
		for _, tagsItem := range resourceMap["tags"].([]interface{}) {
			tagsItemModel := resourceIBMCbrRuleMapToResourceTagAttribute(tagsItem.(map[string]interface{}))
			tags = append(tags, tagsItemModel)
		}
		resource.Tags = tags
	}

	return resource
}

func resourceIBMCbrRuleMapToResourceAttribute(resourceAttributeMap map[string]interface{}) contextbasedrestrictionsv1.ResourceAttribute {
	resourceAttribute := contextbasedrestrictionsv1.ResourceAttribute{}

	resourceAttribute.Name = core.StringPtr(resourceAttributeMap["name"].(string))
	resourceAttribute.Value = core.StringPtr(resourceAttributeMap["value"].(string))
	if resourceAttributeMap["operator"] != nil && resourceAttributeMap["operator"] != "" {
		resourceAttribute.Operator = core.StringPtr(resourceAttributeMap["operator"].(string))
	}

	return resourceAttribute
}

func resourceIBMCbrRuleMapToResourceTagAttribute(resourceTagAttributeMap map[string]interface{}) contextbasedrestrictionsv1.ResourceTagAttribute {
	resourceTagAttribute := contextbasedrestrictionsv1.ResourceTagAttribute{}

	resourceTagAttribute.Name = core.StringPtr(resourceTagAttributeMap["name"].(string))
	resourceTagAttribute.Value = core.StringPtr(resourceTagAttributeMap["value"].(string))
	if resourceTagAttributeMap["operator"] != nil && resourceTagAttributeMap["operator"] != "" {
		resourceTagAttribute.Operator = core.StringPtr(resourceTagAttributeMap["operator"].(string))
	}

	return resourceTagAttribute
}

func resourceIBMCbrRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
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

	if err = d.Set("description", rule.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if rule.Contexts != nil {
		contexts := []map[string]interface{}{}
		for _, contextsItem := range rule.Contexts {
			contextsItemMap := resourceIBMCbrRuleRuleContextToMap(contextsItem)
			contexts = append(contexts, contextsItemMap)
		}
		if err = d.Set("contexts", contexts); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting contexts: %s", err))
		}
	}
	if rule.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range rule.Resources {
			resourcesItemMap := resourceIBMCbrRuleResourceToMap(resourcesItem)
			resources = append(resources, resourcesItemMap)
		}
		if err = d.Set("resources", resources); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resources: %s", err))
		}
	}
	if err = d.Set("crn", rule.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("href", rule.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("created_at", dateTimeToString(rule.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by_id", rule.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}
	if err = d.Set("last_modified_at", dateTimeToString(rule.LastModifiedAt)); err != nil {
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

func resourceIBMCbrRuleRuleContextToMap(ruleContext contextbasedrestrictionsv1.RuleContext) map[string]interface{} {
	ruleContextMap := map[string]interface{}{}

	attributes := []map[string]interface{}{}
	for _, attributesItem := range ruleContext.Attributes {
		attributesItemMap := resourceIBMCbrRuleRuleContextAttributeToMap(attributesItem)
		attributes = append(attributes, attributesItemMap)
		// TODO: handle Attributes of type TypeList -- list of non-primitive, not model items
	}
	ruleContextMap["attributes"] = attributes

	return ruleContextMap
}

func resourceIBMCbrRuleRuleContextAttributeToMap(ruleContextAttribute contextbasedrestrictionsv1.RuleContextAttribute) map[string]interface{} {
	ruleContextAttributeMap := map[string]interface{}{}

	ruleContextAttributeMap["name"] = ruleContextAttribute.Name
	ruleContextAttributeMap["value"] = ruleContextAttribute.Value

	return ruleContextAttributeMap
}

func resourceIBMCbrRuleResourceToMap(resource contextbasedrestrictionsv1.Resource) map[string]interface{} {
	resourceMap := map[string]interface{}{}

	attributes := []map[string]interface{}{}
	for _, attributesItem := range resource.Attributes {

		if *attributesItem.Name != "accountId" {
			attributesItemMap := resourceIBMCbrRuleResourceAttributeToMap(attributesItem)
			attributes = append(attributes, attributesItemMap)
		}
		// TODO: handle Attributes of type TypeList -- list of non-primitive, not model items
	}
	resourceMap["attributes"] = attributes
	if resource.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range resource.Tags {
			tagsItemMap := resourceIBMCbrRuleResourceTagAttributeToMap(tagsItem)
			tags = append(tags, tagsItemMap)
			// TODO: handle Tags of type TypeList -- list of non-primitive, not model items
		}
		resourceMap["tags"] = tags
	}

	return resourceMap
}

func resourceIBMCbrRuleResourceAttributeToMap(resourceAttribute contextbasedrestrictionsv1.ResourceAttribute) map[string]interface{} {
	resourceAttributeMap := map[string]interface{}{}

	resourceAttributeMap["name"] = resourceAttribute.Name
	resourceAttributeMap["value"] = resourceAttribute.Value
	if resourceAttribute.Operator != nil {
		resourceAttributeMap["operator"] = resourceAttribute.Operator
	}

	return resourceAttributeMap
}

func resourceIBMCbrRuleResourceTagAttributeToMap(resourceTagAttribute contextbasedrestrictionsv1.ResourceTagAttribute) map[string]interface{} {
	resourceTagAttributeMap := map[string]interface{}{}

	resourceTagAttributeMap["name"] = resourceTagAttribute.Name
	resourceTagAttributeMap["value"] = resourceTagAttribute.Value
	if resourceTagAttribute.Operator != nil {
		resourceTagAttributeMap["operator"] = resourceTagAttribute.Operator
	}

	return resourceTagAttributeMap
}

func resourceIBMCbrRuleUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceRuleOptions := &contextbasedrestrictionsv1.ReplaceRuleOptions{}

	accountID, err := getIBMCbrAccountId(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	replaceRuleOptions.SetRuleID(d.Id())
	if _, ok := d.GetOk("description"); ok {
		replaceRuleOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("contexts"); ok {
		var contexts []contextbasedrestrictionsv1.RuleContext
		for _, e := range d.Get("contexts").([]interface{}) {
			value := e.(map[string]interface{})
			contextsItem := resourceIBMCbrRuleMapToRuleContext(value)
			contexts = append(contexts, contextsItem)
		}
		replaceRuleOptions.SetContexts(contexts)
	}
	if _, ok := d.GetOk("resources"); ok {
		var resources []contextbasedrestrictionsv1.Resource
		for _, e := range d.Get("resources").([]interface{}) {
			value := e.(map[string]interface{})
			resourcesItem := resourceIBMCbrRuleMapToResource(value, accountID)
			resources = append(resources, resourcesItem)
		}
		replaceRuleOptions.SetResources(resources)
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
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteRuleOptions := &contextbasedrestrictionsv1.DeleteRuleOptions{}

	deleteRuleOptions.SetRuleID(d.Id())

	// deleteRuleOptions.SetIfMatch(d.Get("version").(string))

	response, err := contextBasedRestrictionsClient.DeleteRuleWithContext(context, deleteRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
