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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/contextbasedrestrictionsv1"
)

func DataSourceIBMCbrRule() *schema.Resource {
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
			"operations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The operations this rule applies to.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_types": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The API types this rule applies to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_type_id": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"enforcement_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule enforcement mode: * `enabled` - The restrictions are enforced and reported. This is the default. * `disabled` - The restrictions are disabled. Nothing is enforced or reported. * `report` - The restrictions are evaluated and reported, but not enforced.",
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
	contextBasedRestrictionsClient, err := meta.(conns.ClientSession).ContextBasedRestrictionsV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_rule", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getRuleOptions := &contextbasedrestrictionsv1.GetRuleOptions{}

	getRuleOptions.SetRuleID(d.Get("rule_id").(string))

	rule, _, err := contextBasedRestrictionsClient.GetRuleWithContext(context, getRuleOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRuleWithContext failed: %s", err.Error()), "(Data) ibm_cbr_rule", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getRuleOptions.RuleID))

	if err = d.Set("crn", rule.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cbr_rule", "read", "set-crn").GetDiag()
	}

	if err = d.Set("description", rule.Description); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting description: %s", err), "(Data) ibm_cbr_rule", "read", "set-description").GetDiag()
	}

	contexts := []map[string]interface{}{}
	if rule.Contexts != nil {
		for _, modelItem := range rule.Contexts {
			modelMap, err := DataSourceIBMCbrRuleRuleContextToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error mapping rule context: %s", err), "(Data) ibm_cbr_rule", "read", "DataSourceIBMCbrRuleRuleContextToMap").GetDiag()
			}
			contexts = append(contexts, modelMap)
		}
	}
	if err = d.Set("contexts", contexts); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting contexts: %s", err), "(Data) ibm_cbr_rule", "read", "set-contexts").GetDiag()
	}

	resources := []map[string]interface{}{}
	if rule.Resources != nil {
		for _, modelItem := range rule.Resources {
			modelMap, err := DataSourceIBMCbrRuleResourceToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error mapping rule resource: %s", err), "(Data) ibm_cbr_rule", "read", "DataSourceIBMCbrRuleResourceToMap").GetDiag()
			}
			resources = append(resources, modelMap)
		}
	}
	if err = d.Set("resources", resources); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resources: %s", err), "(Data) ibm_cbr_rule", "read", "set-resources").GetDiag()
	}

	operations := []map[string]interface{}{}
	if !core.IsNil(rule.Operations) {
		modelMap, err := DataSourceIBMCbrRuleNewRuleOperationsToMap(rule.Operations)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cbr_rule", "read", "DataSourceIBMCbrRuleNewRuleOperationsToMap").GetDiag()
		}
		operations = append(operations, modelMap)
	}
	if err = d.Set("operations", operations); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting operations: %s", err), "(Data) ibm_cbr_rule", "read", "set-operations").GetDiag()
	}

	if !core.IsNil(rule.EnforcementMode) {
		if err = d.Set("enforcement_mode", rule.EnforcementMode); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting enforcement_mode: %s", err), "(Data) ibm_cbr_rule", "read", "set-enforcement_mode").GetDiag()
		}
	}

	if err = d.Set("href", rule.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cbr_rule", "read", "set-href").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(rule.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_cbr_rule", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("created_by_id", rule.CreatedByID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_by_id: %s", err), "(Data) ibm_cbr_rule", "read", "set-created_by_id").GetDiag()
	}

	if err = d.Set("last_modified_at", flex.DateTimeToString(rule.LastModifiedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting last_modified_at: %s", err), "(Data) ibm_cbr_rule", "read", "set-last_modified_at").GetDiag()
	}

	if err = d.Set("last_modified_by_id", rule.LastModifiedByID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting last_modified_by_id: %s", err), "(Data) ibm_cbr_rule", "read", "set-last_modified_by_id").GetDiag()
	}

	return nil
}

func DataSourceIBMCbrRuleRuleContextToMap(model *contextbasedrestrictionsv1.RuleContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Attributes != nil {
		attributes := []map[string]interface{}{}
		for _, attributesItem := range model.Attributes {
			attributesItemMap, err := DataSourceIBMCbrRuleRuleContextAttributeToMap(&attributesItem)
			if err != nil {
				return modelMap, err
			}
			attributes = append(attributes, attributesItemMap)
		}
		modelMap["attributes"] = attributes
	}
	return modelMap, nil
}

func DataSourceIBMCbrRuleRuleContextAttributeToMap(model *contextbasedrestrictionsv1.RuleContextAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIBMCbrRuleResourceToMap(model *contextbasedrestrictionsv1.Resource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Attributes != nil {
		attributes := []map[string]interface{}{}
		for _, attributesItem := range model.Attributes {
			attributesItemMap, err := DataSourceIBMCbrRuleResourceAttributeToMap(&attributesItem)
			if err != nil {
				return modelMap, err
			}
			attributes = append(attributes, attributesItemMap)
		}
		modelMap["attributes"] = attributes
	}
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIBMCbrRuleResourceTagAttributeToMap(&tagsItem)
			if err != nil {
				return modelMap, err
			}
			tags = append(tags, tagsItemMap)
		}
		modelMap["tags"] = tags
	}
	return modelMap, nil
}

func DataSourceIBMCbrRuleResourceAttributeToMap(model *contextbasedrestrictionsv1.ResourceAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Operator != nil {
		modelMap["operator"] = *model.Operator
	}
	return modelMap, nil
}

func DataSourceIBMCbrRuleResourceTagAttributeToMap(model *contextbasedrestrictionsv1.ResourceTagAttribute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	if model.Operator != nil {
		modelMap["operator"] = *model.Operator
	}
	return modelMap, nil
}

func DataSourceIBMCbrRuleNewRuleOperationsToMap(model *contextbasedrestrictionsv1.NewRuleOperations) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.APITypes != nil {
		apiTypes := []map[string]interface{}{}
		for _, apiTypesItem := range model.APITypes {
			apiTypesItemMap, err := DataSourceIBMCbrRuleNewRuleOperationsAPITypesItemToMap(&apiTypesItem)
			if err != nil {
				return modelMap, err
			}
			apiTypes = append(apiTypes, apiTypesItemMap)
		}
		modelMap["api_types"] = apiTypes
	}
	return modelMap, nil
}

func DataSourceIBMCbrRuleNewRuleOperationsAPITypesItemToMap(model *contextbasedrestrictionsv1.NewRuleOperationsAPITypesItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.APITypeID != nil {
		modelMap["api_type_id"] = *model.APITypeID
	}
	return modelMap, nil
}
