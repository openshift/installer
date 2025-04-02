// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIbmSchematicsPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmSchematicsPolicyRead,

		Schema: map[string]*schema.Schema{
			"policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID to get the details of policy.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of Schematics customization policy.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of Schematics customization policy.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name for the policy.  By default, Policy will be created in `default` Resource Group.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tags for the Schematics customization policy.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User defined status of the Schematics object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User-defined states  * `draft` Object can be modified; can be used by Jobs run by the author, during execution  * `live` Object can be modified; can be used by Jobs during execution  * `locked` Object cannot be modified; can be used by Jobs during execution  * `disable` Object can be modified. cannot be used by Jobs during execution.",
						},
						"set_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the User who set the state of the Object.",
						},
						"set_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "When the User who set the state of the Object.",
						},
					},
				},
			},
			"kind": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Policy kind or categories for managing and deriving policy decision  * `agent_assignment_policy` Agent assignment policy for job execution.",
			},
			"target": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The objects for the Schematics policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"selector_kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Types of schematics object selector.",
						},
						"selector_ids": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Static selectors of schematics object ids (agent, workspace, action or blueprint) for the Schematics policy.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"selector_scope": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Selectors to dynamically list of schematics object ids (agent, workspace, action or blueprint) for the Schematics policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"kind": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the Schematics automation resource.",
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The tag based selector.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"resource_groups": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The resource group based selector.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"locations": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The location based selector.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
			"parameter": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The parameter to tune the Schematics policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_assignment_policy_parameter": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Parameters for the `agent_assignment_policy`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"selector_kind": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Types of schematics object selector.",
									},
									"selector_ids": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The static selectors of schematics object ids (workspace, action or blueprint) for the Schematics policy.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"selector_scope": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The selectors to dynamically list of schematics object ids (workspace, action or blueprint) for the Schematics policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"kind": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the Schematics automation resource.",
												},
												"tags": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The tag based selector.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"resource_groups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The resource group based selector.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"locations": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The location based selector.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The system generated policy Id.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy CRN.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Account id.",
			},
			"scoped_resources": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of scoped Schematics resources targeted by the policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"kind": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Schematics automation resource.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Schematics resource Id.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy creation time.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user who created the policy.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy updation time.",
			},
		},
	}
}

func dataSourceIbmSchematicsPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getPolicyOptions := &schematicsv1.GetPolicyOptions{}

	getPolicyOptions.SetPolicyID(d.Get("policy_id").(string))

	policy, response, err := schematicsClient.GetPolicyWithContext(context, getPolicyOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead GetPolicyWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s", *getPolicyOptions.PolicyID))

	if err = d.Set("name", policy.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("description", policy.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_group", policy.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("location", policy.Location); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	state := []map[string]interface{}{}
	if policy.State != nil {
		modelMap, err := dataSourceIbmSchematicsPolicyUserStateToMap(policy.State)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		state = append(state, modelMap)
	}
	if err = d.Set("state", state); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("kind", policy.Kind); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if policy.Tags != nil {
		if err = d.Set("tags", policy.Tags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	target := []map[string]interface{}{}
	if policy.Target != nil {
		modelMap, err := dataSourceIbmSchematicsPolicyPolicyObjectsToMap(policy.Target)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		target = append(target, modelMap)
	}
	if err = d.Set("target", target); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	parameter := []map[string]interface{}{}
	if policy.Parameter != nil {
		modelMap, err := dataSourceIbmSchematicsPolicyPolicyParameterToMap(policy.Parameter)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		parameter = append(parameter, modelMap)
	}
	if err = d.Set("parameter", parameter); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("id", policy.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("crn", policy.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("account", policy.Account); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	scopedResources := []map[string]interface{}{}
	if policy.ScopedResources != nil {
		for _, modelItem := range policy.ScopedResources {
			modelMap, err := dataSourceIbmSchematicsPolicyScopedResourceToMap(&modelItem)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed: %s", err.Error()), "ibm_schematics_policy", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			scopedResources = append(scopedResources, modelMap)
		}
	}
	if err = d.Set("scoped_resources", scopedResources); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(policy.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created_by", policy.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(policy.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIbmSchematicsPolicyRead failed with error: %s", err), "ibm_schematics_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return nil
}

func dataSourceIbmSchematicsPolicyUserStateToMap(model *schematicsv1.UserState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.State != nil {
		modelMap["state"] = *model.State
	}
	if model.SetBy != nil {
		modelMap["set_by"] = *model.SetBy
	}
	if model.SetAt != nil {
		modelMap["set_at"] = model.SetAt.String()
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsPolicyPolicyObjectsToMap(model *schematicsv1.PolicyObjects) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SelectorKind != nil {
		modelMap["selector_kind"] = *model.SelectorKind
	}
	if model.SelectorIds != nil {
		modelMap["selector_ids"] = model.SelectorIds
	}
	if model.SelectorScope != nil {
		selectorScope := []map[string]interface{}{}
		for _, selectorScopeItem := range model.SelectorScope {
			selectorScopeItemMap, err := dataSourceIbmSchematicsPolicyPolicyObjectSelectorToMap(&selectorScopeItem)
			if err != nil {
				return modelMap, err
			}
			selectorScope = append(selectorScope, selectorScopeItemMap)
		}
		modelMap["selector_scope"] = selectorScope
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsPolicyPolicyObjectSelectorToMap(model *schematicsv1.PolicyObjectSelector) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Kind != nil {
		modelMap["kind"] = *model.Kind
	}
	if model.Tags != nil {
		modelMap["tags"] = model.Tags
	}
	if model.ResourceGroups != nil {
		modelMap["resource_groups"] = model.ResourceGroups
	}
	if model.Locations != nil {
		modelMap["locations"] = model.Locations
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsPolicyPolicyParameterToMap(model *schematicsv1.PolicyParameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentAssignmentPolicyParameter != nil {
		agentAssignmentPolicyParameterMap, err := dataSourceIbmSchematicsPolicyAgentAssignmentPolicyParameterToMap(model.AgentAssignmentPolicyParameter)
		if err != nil {
			return modelMap, err
		}
		modelMap["agent_assignment_policy_parameter"] = []map[string]interface{}{agentAssignmentPolicyParameterMap}
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsPolicyAgentAssignmentPolicyParameterToMap(model *schematicsv1.AgentAssignmentPolicyParameter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SelectorKind != nil {
		modelMap["selector_kind"] = *model.SelectorKind
	}
	if model.SelectorIds != nil {
		modelMap["selector_ids"] = model.SelectorIds
	}
	if model.SelectorScope != nil {
		selectorScope := []map[string]interface{}{}
		for _, selectorScopeItem := range model.SelectorScope {
			selectorScopeItemMap, err := dataSourceIbmSchematicsPolicyPolicyObjectSelectorToMap(&selectorScopeItem)
			if err != nil {
				return modelMap, err
			}
			selectorScope = append(selectorScope, selectorScopeItemMap)
		}
		modelMap["selector_scope"] = selectorScope
	}
	return modelMap, nil
}

func dataSourceIbmSchematicsPolicyScopedResourceToMap(model *schematicsv1.ScopedResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Kind != nil {
		modelMap["kind"] = *model.Kind
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}
