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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iampolicymanagementv1"
)

func DataSourceIBMIAMPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIAMPolicyAssignmentRead,

		Schema: map[string]*schema.Schema{
			"assignment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy template assignment ID.",
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "1.0",
				Description: "The policy template assignment new format",
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
			"subject": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "assignment access type subject details",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"target": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "assignment target details",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "policy template details",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Object for each account assigned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "assignment target details",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
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
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enterprise accountID.",
			},
		},
	}
}

func dataSourceIBMIAMPolicyAssignmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamPolicyManagementClient, err := meta.(conns.ClientSession).IAMPolicyManagementV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getPolicyAssignmentOptions := &iampolicymanagementv1.GetPolicyAssignmentOptions{}

	getPolicyAssignmentOptions.SetAssignmentID(d.Get("assignment_id").(string))
	getPolicyAssignmentOptions.SetVersion("1.0")

	assignmentResponse, response, err := iamPolicyManagementClient.GetPolicyAssignmentWithContext(context, getPolicyAssignmentOptions)
	if err != nil {
		log.Printf("[DEBUG] GetPolicyAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetPolicyAssignmentWithContext failed %s\n%s", err, response))
	}

	policyAssignmentRecord := assignmentResponse.(*iampolicymanagementv1.PolicyTemplateAssignmentItems)
	d.SetId(*policyAssignmentRecord.ID)

	targetMap, err := ResourceIBMPolicyAssignmentAssignmentTargetDetailsToMap(policyAssignmentRecord.Target)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("id", policyAssignmentRecord.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %s", err))
	}
	if err = d.Set("target", targetMap); err != nil {
		return diag.FromErr(fmt.Errorf("error setting target: %s", err))
	}
	if policyAssignmentRecord.Template != nil {
		templateMap, err := DataSourceIBMPolicyAssignmentAssignmentTemplateDetailsToMap(policyAssignmentRecord.Template)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("template", templateMap); err != nil {
			return diag.FromErr(fmt.Errorf("error setting template: %s", err))
		}
	}

	if err = d.Set("href", policyAssignmentRecord.Href); err != nil {
		return diag.FromErr(fmt.Errorf("error setting href: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(policyAssignmentRecord.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
	}

	if err = d.Set("last_modified_at", flex.DateTimeToString(policyAssignmentRecord.LastModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting last_modified_at: %s", err))
	}

	if err = d.Set("account_id", policyAssignmentRecord.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting account_id: %s", err))
	}

	resources := []map[string]interface{}{}
	if policyAssignmentRecord.Resources != nil {
		for _, modelItem := range policyAssignmentRecord.Resources {
			modelMap, err := ResourceIBMPolicyAssignmentPolicyAssignmentV1ResourcesToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, modelMap)
		}
	}
	if err = d.Set("resources", resources); err != nil {
		return diag.FromErr(fmt.Errorf("error setting resources %s", err))
	}

	if policyAssignmentRecord.Subject != nil {
		modelMap, err := DataSourceIBMPolicyAssignmentPolicyAssignmentV1Subject(policyAssignmentRecord.Subject)
		if err != nil {
			return diag.FromErr(err)
		}

		if err = d.Set("subject", modelMap); err != nil {
			return diag.FromErr(fmt.Errorf("error setting subject %s", err))
		}
	}

	return nil
}

func ResourceIBMPolicyAssignmentConflictsWithToMap(model *iampolicymanagementv1.ConflictsWith) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Etag != nil {
		modelMap["etag"] = *model.Etag
	}
	if model.Role != nil {
		modelMap["role"] = *model.Role
	}
	if model.Policy != nil {
		modelMap["policy"] = *model.Policy
	}
	return modelMap, nil
}

func ResourceIBMPolicyAssignmentErrorDetailsToMap(model *iampolicymanagementv1.ErrorDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConflictsWith != nil {
		conflictsWithMap, err := ResourceIBMPolicyAssignmentConflictsWithToMap(model.ConflictsWith)
		if err != nil {
			return modelMap, err
		}
		modelMap["conflicts_with"] = []map[string]interface{}{conflictsWithMap}
	}
	return modelMap, nil
}

func ResourceIBMPolicyAssignmentAssignmentResourceCreatedToMap(model *iampolicymanagementv1.AssignmentResourceCreated) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	return modelMap, nil
}

func ResourceIBMPolicyAssignmentErrorObjectToMap(model *iampolicymanagementv1.ErrorObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.Details != nil {
		detailsMap, err := ResourceIBMPolicyAssignmentErrorDetailsToMap(model.Details)
		if err != nil {
			return modelMap, err
		}
		modelMap["details"] = []map[string]interface{}{detailsMap}
	}
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMPolicyAssignmentErrorResponseToMap(model *iampolicymanagementv1.ErrorResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Trace != nil {
		modelMap["trace"] = *model.Trace
	}
	if model.Errors != nil {
		errors := []map[string]interface{}{}
		for _, errorsItem := range model.Errors {
			errorsItemMap, err := ResourceIBMPolicyAssignmentErrorObjectToMap(&errorsItem)
			if err != nil {
				return modelMap, err
			}
			errors = append(errors, errorsItemMap)
		}
		modelMap["errors"] = errors
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = flex.IntValue(model.StatusCode)
	}
	return modelMap, nil
}

func ResourceIBMPolicyAssignmentPolicyAssignmentResourcePolicyToMap(model *iampolicymanagementv1.PolicyAssignmentResourcePolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceCreated != nil {
		resourceCreatedMap, err := ResourceIBMPolicyAssignmentAssignmentResourceCreatedToMap(model.ResourceCreated)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_created"] = []map[string]interface{}{resourceCreatedMap}
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.ErrorMessage != nil {
		errorMessageMap, err := ResourceIBMPolicyAssignmentErrorResponseToMap(model.ErrorMessage)
		if err != nil {
			return modelMap, err
		}
		modelMap["error_message"] = []map[string]interface{}{errorMessageMap}
	}
	return modelMap, nil
}

func ResourceIBMPolicyAssignmentPolicyAssignmentV1ResourcesToMap(model *iampolicymanagementv1.PolicyAssignmentV1Resources) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Target != nil {
		targetMap, err := ResourceIBMPolicyAssignmentResourceTargetDetailsToMap(model.Target)
		if err != nil {
			return modelMap, err
		}
		modelMap["target"] = targetMap
	}
	if model.Policy != nil {
		policyMap, err := ResourceIBMPolicyAssignmentPolicyAssignmentResourcePolicyToMap(model.Policy)
		if err != nil {
			return modelMap, err
		}
		modelMap["policy"] = []map[string]interface{}{policyMap}
	}
	return modelMap, nil
}
func DataSourceIBMPolicyAssignmentPolicyAssignmentV1Subject(model *iampolicymanagementv1.PolicyAssignmentV1Subject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}
