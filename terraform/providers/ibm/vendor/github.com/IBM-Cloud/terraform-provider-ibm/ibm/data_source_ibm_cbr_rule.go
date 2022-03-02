// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func dataSourceIBMCbrRule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCbrRuleRead,

		Schema: map[string]*schema.Schema{
			"rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of a rule.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule CRN.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the rule.",
			},
			"contexts": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The contexts this rule applies to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attribute name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
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
				Computed:    true,
				Description: "The resources this rule apply to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attributes": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource attributes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attribute name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attribute value.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attribute operator.",
									},
								},
							},
						},
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The optional resource tags.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tag attribute name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tag attribute value.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The attribute operator.",
									},
								},
							},
						},
					},
				},
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
		},
	}
}

func dataSourceIBMCbrRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	contextBasedRestrictionsClient, err := meta.(ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleOptions := &contextbasedrestrictionsv1.GetRuleOptions{}

	getRuleOptions.SetRuleID(d.Get("rule_id").(string))

	rule, response, err := contextBasedRestrictionsClient.GetRuleWithContext(context, getRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] GetRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getRuleOptions.RuleID))
	if err = d.Set("crn", rule.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("description", rule.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if rule.Contexts != nil {
		err = d.Set("contexts", dataSourceRuleFlattenContexts(rule.Contexts))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting contexts %s", err))
		}
	}

	if rule.Resources != nil {
		err = d.Set("resources", dataSourceRuleFlattenResources(rule.Resources))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resources %s", err))
		}
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

	return nil
}

func dataSourceRuleFlattenContexts(result []contextbasedrestrictionsv1.RuleContext) (contexts []map[string]interface{}) {
	for _, contextsItem := range result {
		contexts = append(contexts, dataSourceRuleContextsToMap(contextsItem))
	}

	return contexts
}

func dataSourceRuleContextsToMap(contextsItem contextbasedrestrictionsv1.RuleContext) (contextsMap map[string]interface{}) {
	contextsMap = map[string]interface{}{}

	if contextsItem.Attributes != nil {
		attributesList := []map[string]interface{}{}
		for _, attributesItem := range contextsItem.Attributes {
			attributesList = append(attributesList, dataSourceRuleContextsAttributesToMap(attributesItem))
		}
		contextsMap["attributes"] = attributesList
	}

	return contextsMap
}

func dataSourceRuleContextsAttributesToMap(attributesItem contextbasedrestrictionsv1.RuleContextAttribute) (attributesMap map[string]interface{}) {
	attributesMap = map[string]interface{}{}

	if attributesItem.Name != nil {
		attributesMap["name"] = attributesItem.Name
	}
	if attributesItem.Value != nil {
		attributesMap["value"] = attributesItem.Value
	}

	return attributesMap
}

func dataSourceRuleFlattenResources(result []contextbasedrestrictionsv1.Resource) (resources []map[string]interface{}) {
	for _, resourcesItem := range result {
		resources = append(resources, dataSourceRuleResourcesToMap(resourcesItem))
	}

	return resources
}

func dataSourceRuleResourcesToMap(resourcesItem contextbasedrestrictionsv1.Resource) (resourcesMap map[string]interface{}) {
	resourcesMap = map[string]interface{}{}

	if resourcesItem.Attributes != nil {
		attributesList := []map[string]interface{}{}
		for _, attributesItem := range resourcesItem.Attributes {
			attributesList = append(attributesList, dataSourceRuleResourcesAttributesToMap(attributesItem))
		}
		resourcesMap["attributes"] = attributesList
	}
	if resourcesItem.Tags != nil {
		tagsList := []map[string]interface{}{}
		for _, tagsItem := range resourcesItem.Tags {
			tagsList = append(tagsList, dataSourceRuleResourcesTagsToMap(tagsItem))
		}
		resourcesMap["tags"] = tagsList
	}

	return resourcesMap
}

func dataSourceRuleResourcesAttributesToMap(attributesItem contextbasedrestrictionsv1.ResourceAttribute) (attributesMap map[string]interface{}) {
	attributesMap = map[string]interface{}{}

	if attributesItem.Name != nil {
		attributesMap["name"] = attributesItem.Name
	}
	if attributesItem.Value != nil {
		attributesMap["value"] = attributesItem.Value
	}
	if attributesItem.Operator != nil {
		attributesMap["operator"] = attributesItem.Operator
	}

	return attributesMap
}

func dataSourceRuleResourcesTagsToMap(tagsItem contextbasedrestrictionsv1.ResourceTagAttribute) (tagsMap map[string]interface{}) {
	tagsMap = map[string]interface{}{}

	if tagsItem.Name != nil {
		tagsMap["name"] = tagsItem.Name
	}
	if tagsItem.Value != nil {
		tagsMap["value"] = tagsItem.Value
	}
	if tagsItem.Operator != nil {
		tagsMap["operator"] = tagsItem.Operator
	}

	return tagsMap
}
