// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iampolicy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

func DataSourceIBMIAMPolicyAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIAMPolicyAssignmentsRead,

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional template id.",
			},
			"template_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional policy template version.",
			},
			"policy_assignments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of policy assignments.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "policy template id.",
						},
						"template_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "policy template version.",
						},
						"assignment_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Passed in value to correlate with other assignments.",
						},
						"target_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Assignment target type.",
						},
						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "assignment target id.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy assignment ID.",
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The href URL that links to the policies assignments API by policy assignment ID.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UTC timestamp when the policy assignment was created.",
						},
						"created_by_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The iam ID of the entity that created the policy assignment.",
						},
						"last_modified_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UTC timestamp when the policy assignment was last modified.",
						},
						"last_modified_by_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The iam ID of the entity that last modified the policy assignment.",
						},
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Object for each account assigned.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Account ID where resources are assigned.",
									},
									"policy": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Set of properties for the assigned resource.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"resource_created": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "On success, includes the  policy assigned.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "policy id.",
															},
														},
													},
												},
												"error_message": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The error response from API.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"trace": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The unique transaction id for the request.",
															},
															"errors": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The errors encountered during the response.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"code": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The API error code for the error.",
																		},
																		"message": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The error message returned by the API.",
																		},
																		"details": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Additional error details.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"conflicts_with": {
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Details of conflicting resource.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"etag": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The revision number of the resource.",
																								},
																								"role": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The conflicting role id.",
																								},
																								"policy": {
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The conflicting policy id.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"more_info": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Additional info for error.",
																		},
																	},
																},
															},
															"status_code": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The http error code of the response.",
															},
														},
													},
												},
												"status": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The policy assignment status.",
												},
											},
										},
									},
								},
							},
						},
						"options": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of objects with required properties for a policy assignment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"subject_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The policy subject type; either 'iam_id' or 'access_group_id'.",
									},
									"subject_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The policy subject id.",
									},
									"root_requester_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The policy assignment requester id.",
									},
									"root_template_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The template id where this policy is being assigned from.",
									},
									"root_template_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The template version where this policy is being assigned from.",
									},
								},
							},
						},
						"account_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise accountID.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIAMPolicyAssignmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Failed to fetch BluemixUserDetails %s", err))
	}

	accountID := userDetails.UserAccount

	listPolicyAssignmentsOptions := &iampolicymanagementv1.ListPolicyAssignmentsOptions{}

	listPolicyAssignmentsOptions.SetAccountID(accountID)
	if _, ok := d.GetOk("template_id"); ok {
		listPolicyAssignmentsOptions.SetTemplateID(d.Get("template_id").(string))
	}
	if _, ok := d.GetOk("template_version"); ok {
		listPolicyAssignmentsOptions.SetTemplateVersion(d.Get("template_version").(string))
	}

	policyTemplateAssignmentCollection, response, err := iamPolicyManagementClient.ListPolicyAssignmentsWithContext(context, listPolicyAssignmentsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListPolicyAssignmentsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListPolicyAssignmentsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMPolicyAssignmentID(d))

	policyAssignments := []map[string]interface{}{}
	if policyTemplateAssignmentCollection.Assignments != nil {
		for _, modelItem := range policyTemplateAssignmentCollection.Assignments {
			modelMap, err := dataSourceIBMPolicyAssignmentPolicyAssignmentRecordToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			policyAssignments = append(policyAssignments, modelMap)
		}
	}
	if err = d.Set("policy_assignments", policyAssignments); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting policy_assignments %s", err))
	}

	return nil
}

// dataSourceIBMPolicyAssignmentID returns a reasonable ID for the list.
func dataSourceIBMPolicyAssignmentID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMPolicyAssignmentPolicyAssignmentRecordToMap(model *iampolicymanagementv1.PolicyAssignment) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["template_id"] = model.TemplateID
	modelMap["template_version"] = model.TemplateVersion
	modelMap["assignment_id"] = model.AssignmentID
	modelMap["target_type"] = model.TargetType
	modelMap["target"] = model.Target
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Href != nil {
		modelMap["href"] = model.Href
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.CreatedByID != nil {
		modelMap["created_by_id"] = model.CreatedByID
	}
	if model.LastModifiedAt != nil {
		modelMap["last_modified_at"] = model.LastModifiedAt.String()
	}
	if model.LastModifiedByID != nil {
		modelMap["last_modified_by_id"] = model.LastModifiedByID
	}
	if model.AccountID != nil {
		modelMap["account_id"] = model.AccountID
	}

	if model.Options != nil {
		options := []map[string]interface{}{}
		for _, modelItem := range model.Options {
			modelMap, err := dataSourceIBMAssignmentPolicyAssignmentOptionsToMap(&modelItem)
			if err != nil {
				return modelMap, err
			}
			options = append(options, modelMap)
		}
		modelMap["options"] = options
	}
	if model.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range model.Resources {
			resourcesItemMap, err := dataSourceIBMPolicyAssignmentPolicyAssignmentResourcesToMap(&resourcesItem)
			if err != nil {
				return modelMap, err
			}
			resources = append(resources, resourcesItemMap)
		}
		modelMap["resources"] = resources
	}
	return modelMap, nil
}

func dataSourceIBMAssignmentPolicyAssignmentOptionsToMap(model *iampolicymanagementv1.PolicyAssignmentOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["subject_type"] = model.SubjectType
	modelMap["subject_id"] = model.SubjectID
	modelMap["root_requester_id"] = model.RootRequesterID
	if model.RootTemplateID != nil {
		modelMap["root_template_id"] = model.RootTemplateID
	}
	if model.RootTemplateVersion != nil {
		modelMap["root_template_version"] = model.RootTemplateVersion
	}
	return modelMap, nil
}
