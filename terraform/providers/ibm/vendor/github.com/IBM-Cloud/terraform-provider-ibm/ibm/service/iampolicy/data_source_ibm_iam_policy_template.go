// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

func DataSourceIBMIAMPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMPolicyTemplateRead,

		Schema: map[string]*schema.Schema{
			"policy_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of policy templates.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "name of template.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "description of template purpose.",
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "account id where this template will be created.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Template vesrsion.",
						},
						"committed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Template vesrsion committed status.",
						},
						"policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The core set of properties associated with the template's policy objet.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The policy type; either 'access' or 'authorization'.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Allows the customer to use their own words to record the purpose/context related to a policy.",
									},
									"resource": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The resource attributes to which the policy grants access.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"attributes": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of resource attributes to which the policy grants access.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of a resource attribute.",
															},
															"operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The operator of an attribute.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The value of a rule or resource attribute; can be boolean or string for resource attribute. Can be string or an array of strings (e.g., array of days to permit access) for rule attribute.",
															},
														},
													},
												},
												"tags": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Optional list of resource tags to which the policy grants access.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The name of an access management tag.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The value of an access management tag.",
															},
															"operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The operator of an access management tag.",
															},
														},
													},
												},
											},
										},
									},
									"pattern": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Indicates pattern of rule, either 'time-based-conditions:once', 'time-based-conditions:weekly:all-day', or 'time-based-conditions:weekly:custom-hours'.",
									},
									"rule_conditions": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "Rule conditions enforced by the policy",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Key of the condition",
												},
												"operator": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Operator of the condition",
												},
												"value": {
													Type:        schema.TypeList,
													Optional:    true,
													Elem:        &schema.Schema{Type: schema.TypeString},
													Description: "Value of the condition",
												},
												"conditions": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Additional Rule conditions enforced by the policy",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Key of the condition",
															},
															"operator": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Operator of the condition",
															},
															"value": {
																Type:        schema.TypeList,
																Required:    true,
																Elem:        &schema.Schema{Type: schema.TypeString},
																Description: "Value of the condition",
															},
														},
													},
												},
											},
										},
									},

									"rule_operator": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Operator that multiple rule conditions are evaluated over",
									},
									"roles": {
										Type:        schema.TypeList,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Role names of the policy definition",
									},
								},
							},
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy template ID.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMPolicyTemplateRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to fetch BluemixUserDetails %s", err))
	}

	accountID := userDetails.UserAccount

	listPolicyTemplatesOptions := &iampolicymanagementv1.ListPolicyTemplatesOptions{}

	listPolicyTemplatesOptions.SetAccountID(accountID)

	policyTemplateCollection, response, err := iamPolicyManagementClient.ListPolicyTemplatesWithContext(context, listPolicyTemplatesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListPolicyTemplatesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListPolicyTemplatesWithContext failed %s\n%s", err, response))
	}

	d.SetId(accountID)

	policyTemplates := []map[string]interface{}{}
	if policyTemplateCollection.PolicyTemplates != nil {
		for _, modelItem := range policyTemplateCollection.PolicyTemplates {
			modelMap, err := dataSourceIBMPolicyTemplatePolicyTemplateToMap(&modelItem, iamPolicyManagementClient)
			if err != nil {
				return diag.FromErr(err)
			}
			policyTemplates = append(policyTemplates, modelMap)
		}
	}
	if err = d.Set("policy_templates", policyTemplates); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting policy_templates %s", err))
	}

	return nil
}

func dataSourceIBMPolicyTemplatePolicyTemplateToMap(model *iampolicymanagementv1.PolicyTemplate, iamPolicyManagementClient *iampolicymanagementv1.IamPolicyManagementV1) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["account_id"] = model.AccountID
	modelMap["version"] = model.Version
	if model.Committed != nil {
		modelMap["committed"] = model.Committed
	}
	policyMap, err := flattenTemplatePolicy(model.Policy, iamPolicyManagementClient)
	if err != nil {
		return modelMap, err
	}
	modelMap["policy"] = []map[string]interface{}{policyMap}
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}
