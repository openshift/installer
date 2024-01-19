// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func DataSourceIbmSccRule() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		ReadContext: dataSourceIbmSccRuleRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(40 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"rule_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the corresponding rule.",
			},
			"created_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the rule was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the rule.",
			},
			"updated_on": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date when the rule was modified.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who modified the rule.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule ID.",
			},
			"account_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account ID.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The details of a rule's response.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The rule type (allowable values are `user_defined` or `system_defined`).",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version number of a rule.",
			},
			"import": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The collection of import parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parameters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of import parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The import parameter name.",
									},
									"display_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The display name of the property.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The propery description.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The property type.",
									},
								},
							},
						},
					},
				},
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The rule target.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target service name.",
						},
						"service_display_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the target service.",
						},
						"resource_kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target resource kind.",
						},
						"additional_target_attributes": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of targets supported properties.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The additional target attribute name.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The operator.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value.",
									},
								},
							},
						},
					},
				},
			},
			"required_config": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The required configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The required config description.",
						},
						"and": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The `AND` required configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The required config description.",
									},
									"or": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The `OR` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The required config description.",
												},
												"property": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property.",
												},
												"operator": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"and": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The `AND` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The required config description.",
												},
												"property": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property.",
												},
												"operator": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"property": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The property.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The operator.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Schema for any JSON type.",
									},
								},
							},
						},
						"or": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The `OR` required configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The required config description.",
									},
									"or": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The `OR` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The required config description.",
												},
												"property": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property.",
												},
												"operator": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"and": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The `AND` required configurations.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The required config description.",
												},
												"property": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The property.",
												},
												"operator": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The operator.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Schema for any JSON type.",
												},
											},
										},
									},
									"property": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The property.",
									},
									"operator": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The operator.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Schema for any JSON type.",
									},
								},
							},
						},
						"property": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The property.",
						},
						"operator": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operator.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Schema for any JSON type.",
						},
					},
				},
			},
			"labels": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of labels.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	})
}

func dataSourceIbmSccRuleRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	configManagerClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getRuleOptions := &securityandcompliancecenterapiv3.GetRuleOptions{}

	getRuleOptions.SetRuleID(d.Get("rule_id").(string))
	getRuleOptions.SetInstanceID(d.Get("instance_id").(string))

	rule, response, err := configManagerClient.GetRuleWithContext(context, getRuleOptions)
	if err != nil {
		log.Printf("[DEBUG] GetRuleWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetRuleWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getRuleOptions.RuleID))

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

	if err = d.Set("id", rule.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting id: %s", err))
	}

	if err = d.Set("account_id", rule.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}

	if err = d.Set("description", rule.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("type", rule.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("version", rule.Version); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}

	importVar := []map[string]interface{}{}
	if rule.Import != nil {
		modelMap, err := dataSourceIbmSccRuleImportToMap(rule.Import)
		if err != nil {
			return diag.FromErr(err)
		}
		importVar = append(importVar, modelMap)
	}
	if err = d.Set("import", importVar); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting import %s", err))
	}

	target := []map[string]interface{}{}
	if rule.Target != nil {
		modelMap, err := dataSourceIbmSccRuleTargetToMap(rule.Target)
		if err != nil {
			return diag.FromErr(err)
		}
		target = append(target, modelMap)
	}

	if err = d.Set("target", target); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target %s", err))
	}

	if err = d.Set("labels", rule.Labels); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting labels: %s", err))
	}

	requiredConfig := []map[string]interface{}{}
	if rule.RequiredConfig != nil {
		modelMap, err := dataSourceIbmSccRuleRequiredConfigToMap(rule.RequiredConfig)
		if err != nil {
			return diag.FromErr(err)
		}
		requiredConfig = append(requiredConfig, modelMap)
	}
	if err = d.Set("required_config", requiredConfig); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting required_config %s", err))
	}

	return nil
}

func dataSourceIbmSccRuleImportToMap(model *securityandcompliancecenterapiv3.Import) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Parameters != nil {
		parameters := []map[string]interface{}{}
		for _, parametersItem := range model.Parameters {
			parametersItemMap, err := dataSourceIbmSccRuleParameterToMap(&parametersItem)
			if err != nil {
				return modelMap, err
			}
			parameters = append(parameters, parametersItemMap)
		}
		modelMap["parameters"] = parameters
	}
	return modelMap, nil
}

func dataSourceIbmSccRuleParameterToMap(model *securityandcompliancecenterapiv3.Parameter) (map[string]interface{}, error) {
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

func dataSourceIbmSccRuleTargetToMap(model *securityandcompliancecenterapiv3.Target) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["service_name"] = model.ServiceName
	if model.ServiceDisplayName != nil {
		modelMap["service_display_name"] = model.ServiceDisplayName
	}
	modelMap["resource_kind"] = model.ResourceKind
	if model.AdditionalTargetAttributes != nil {
		additionalTargetAttributes := []map[string]interface{}{}
		for _, additionalTargetAttributesItem := range model.AdditionalTargetAttributes {
			additionalTargetAttributesItemMap, err := dataSourceIbmSccRuleAdditionalTargetAttributeToMap(&additionalTargetAttributesItem)
			if err != nil {
				return modelMap, err
			}
			additionalTargetAttributes = append(additionalTargetAttributes, additionalTargetAttributesItemMap)
		}
		modelMap["additional_target_attributes"] = additionalTargetAttributes
	}
	return modelMap, nil
}

func dataSourceIbmSccRuleAdditionalTargetAttributeToMap(model *securityandcompliancecenterapiv3.AdditionalTargetAttribute) (map[string]interface{}, error) {
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

func dataSourceIbmSccRuleRequiredConfigToMap(model securityandcompliancecenterapiv3.RequiredConfigIntf) (map[string]interface{}, error) {
	if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd); ok {
		return dataSourceIbmSccRuleRequiredConfigAndToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr); ok {
		return dataSourceIbmSccRuleRequiredConfigOrToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase); ok {
		return dataSourceIbmSccRuleRequiredConfigRequiredConfigBaseToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfig); ok {
		modelMap := make(map[string]interface{})
		model := model.(*securityandcompliancecenterapiv3.RequiredConfig)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.And != nil {
			and := []map[string]interface{}{}
			for _, andItem := range model.And {
				andItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(andItem)
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
				orItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(orItem)
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
		if model.Value != nil {
			modelMap["value"] = model.Value
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized securityandcompliancecenterapiv3.RequiredConfigIntf subtype encountered")
	}
}

func dataSourceIbmSccRuleRequiredConfigItemsToMap(model securityandcompliancecenterapiv3.RequiredConfigItemsIntf) (map[string]interface{}, error) {
	if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr); ok {
		return dataSourceIbmSccRuleRequiredConfigItemsRequiredConfigOrToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd); ok {
		return dataSourceIbmSccRuleRequiredConfigItemsRequiredConfigAndToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase); ok {
		return dataSourceIbmSccRuleRequiredConfigItemsRequiredConfigBaseToMap(model.(*securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase))
	} else if _, ok := model.(*securityandcompliancecenterapiv3.RequiredConfigItems); ok {
		modelMap := make(map[string]interface{})
		model := model.(*securityandcompliancecenterapiv3.RequiredConfigItems)
		if model.Description != nil {
			modelMap["description"] = model.Description
		}
		if model.Or != nil {
			or := []map[string]interface{}{}
			for _, orItem := range model.Or {
				orItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(orItem)
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
				andItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(andItem)
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
		if model.Value != nil {
			modelMap["value"] = model.Value
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized securityandcompliancecenterapiv3.RequiredConfigItemsIntf subtype encountered")
	}
}

func dataSourceIbmSccRuleRequiredConfigBaseToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase) (map[string]interface{}, error) {
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

func dataSourceIbmSccRuleRequiredConfigItemsRequiredConfigOrToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigOr) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Or != nil {
		or := []map[string]interface{}{}
		for _, orItem := range model.Or {
			orItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(orItem.(*securityandcompliancecenterapiv3.RequiredConfigItems))
			if err != nil {
				return modelMap, err
			}
			or = append(or, orItemMap)
		}
		modelMap["or"] = or
	}
	return modelMap, nil
}

func dataSourceIbmSccRuleRequiredConfigItemsRequiredConfigAndToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigAnd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.And != nil {
		and := []map[string]interface{}{}
		for _, andItem := range model.And {
			andItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(andItem.(*securityandcompliancecenterapiv3.RequiredConfigItems))
			if err != nil {
				return modelMap, err
			}
			and = append(and, andItemMap)
		}
		modelMap["and"] = and
	}
	return modelMap, nil
}

func dataSourceIbmSccRuleRequiredConfigItemsRequiredConfigBaseToMap(model *securityandcompliancecenterapiv3.RequiredConfigItemsRequiredConfigBase) (map[string]interface{}, error) {
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

func dataSourceIbmSccRuleRequiredConfigAndToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigAnd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.And != nil {
		and := []map[string]interface{}{}
		for _, andItem := range model.And {
			andItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(andItem)
			if err != nil {
				return modelMap, err
			}
			and = append(and, andItemMap)
		}
		modelMap["and"] = and
	}
	return modelMap, nil
}

func dataSourceIbmSccRuleRequiredConfigOrToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigOr) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Or != nil {
		or := []map[string]interface{}{}
		for _, orItem := range model.Or {
			orItemMap, err := dataSourceIbmSccRuleRequiredConfigItemsToMap(orItem)
			if err != nil {
				return modelMap, err
			}
			or = append(or, orItemMap)
		}
		modelMap["or"] = or
	}
	return modelMap, nil
}

func dataSourceIbmSccRuleRequiredConfigRequiredConfigBaseToMap(model *securityandcompliancecenterapiv3.RequiredConfigRequiredConfigBase) (map[string]interface{}, error) {
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
