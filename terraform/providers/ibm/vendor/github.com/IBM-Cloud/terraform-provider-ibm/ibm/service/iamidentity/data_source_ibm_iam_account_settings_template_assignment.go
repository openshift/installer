// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

func DataSourceIBMAccountSettingsTemplateAssignment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAccountSettingsTemplateAssignmentRead,

		Schema: map[string]*schema.Schema{
			"assignment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the Assignment Record.",
			},
			"include_history": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defines if the entity history is included in the response.",
			},
			"context": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Context with key properties for problem determination.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"transaction_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The transaction ID of the inbound REST request.",
						},
						"operation": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation of the inbound REST request.",
						},
						"user_agent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user agent of the inbound REST request.",
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of that cluster.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID of the server instance processing the request.",
						},
						"thread_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The thread ID of the server instance processing the request.",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The host of the server instance processing the request.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start time of the request.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The finish time of the request.",
						},
						"elapsed_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The elapsed time in msec.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster name.",
						},
					},
				},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assignment record Id.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enterprise account Id.",
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Template Id.",
			},
			"template_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Template version.",
			},
			"target_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Assignment target type.",
			},
			"target": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Assignment target.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assignment status.",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Status breakdown per target account of IAM resources created or errors encountered in attempting to create those IAM resources. IAM resources are only included in the response providing the assignment is not in progress. IAM resources are also only included when getting a single assignment, and excluded by list APIs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target account where the IAM resource is created.",
						},
						"account_settings": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy Template Id, only returned for a profile assignment with policy references.",
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Policy version, only returned for a profile assignment with policy references.",
									},
									"resource_created": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Body parameters for created resource.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Id of the created resource.",
												},
											},
										},
									},
									"error_message": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Body parameters for assignment error.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the error.",
												},
												"error_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Internal error code.",
												},
												"message": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Error message detailing the nature of the error.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Internal status code for the error.",
												},
											},
										},
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status for the target account's assignment.",
									},
								},
							},
						},
					},
				},
			},
			"history": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Assignment history.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestamp when the action was triggered.",
						},
						"iam_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IAM ID of the identity which triggered the action.",
						},
						"iam_id_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account of the identity which triggered the action.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Action of the history entry.",
						},
						"params": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Params of the history entry.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Message which summarizes the executed action.",
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Href.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assignment created at.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the identity that created the assignment.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Assignment modified at.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IAMid of the identity that last modified the assignment.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Entity tag for this assignment record.",
			},
		},
	}
}

func dataSourceIBMAccountSettingsTemplateAssignmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getAccountSettingsAssignmentOptions := &iamidentityv1.GetAccountSettingsAssignmentOptions{}

	getAccountSettingsAssignmentOptions.SetAssignmentID(d.Get("assignment_id").(string))
	if _, ok := d.GetOk("include_history"); ok {
		getAccountSettingsAssignmentOptions.SetIncludeHistory(d.Get("include_history").(bool))
	}

	templateAssignmentResponse, response, err := iamIdentityClient.GetAccountSettingsAssignmentWithContext(context, getAccountSettingsAssignmentOptions)
	if err != nil {
		log.Printf("[DEBUG] GetAccountSettingsAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAccountSettingsAssignmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getAccountSettingsAssignmentOptions.AssignmentID))

	var respContext []map[string]interface{}
	if templateAssignmentResponse.Context != nil {
		modelMap, err := dataSourceIBMAccountSettingsTemplateAssignmentResponseContextToMap(templateAssignmentResponse.Context)
		if err != nil {
			return diag.FromErr(err)
		}
		respContext = append(respContext, modelMap)
	}
	if err = d.Set("context", respContext); err != nil {
		return diag.FromErr(fmt.Errorf("error setting context %s", err))
	}

	if err = d.Set("id", templateAssignmentResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting id: %s", err))
	}

	if err = d.Set("account_id", templateAssignmentResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting account_id: %s", err))
	}

	if err = d.Set("template_id", templateAssignmentResponse.TemplateID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting template_id: %s", err))
	}

	if err = d.Set("template_version", flex.IntValue(templateAssignmentResponse.TemplateVersion)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting template_version: %s", err))
	}

	if err = d.Set("target_type", templateAssignmentResponse.TargetType); err != nil {
		return diag.FromErr(fmt.Errorf("error setting target_type: %s", err))
	}

	if err = d.Set("target", templateAssignmentResponse.Target); err != nil {
		return diag.FromErr(fmt.Errorf("error setting target: %s", err))
	}

	if err = d.Set("status", templateAssignmentResponse.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting status: %s", err))
	}

	var resources []map[string]interface{}
	if templateAssignmentResponse.Resources != nil {
		for _, modelItem := range templateAssignmentResponse.Resources {
			modelMap, err := dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResponseResourceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, modelMap)
		}
	}
	if err = d.Set("resources", resources); err != nil {
		return diag.FromErr(fmt.Errorf("error setting resources %s", err))
	}

	var history []map[string]interface{}
	if templateAssignmentResponse.History != nil {
		for _, modelItem := range templateAssignmentResponse.History {
			modelMap, err := dataSourceIBMAccountSettingsTemplateAssignmentEnityHistoryRecordToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, modelMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("error setting history %s", err))
	}

	if err = d.Set("href", templateAssignmentResponse.Href); err != nil {
		return diag.FromErr(fmt.Errorf("error setting href: %s", err))
	}

	if err = d.Set("created_at", templateAssignmentResponse.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
	}

	if err = d.Set("created_by_id", templateAssignmentResponse.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting created_by_id: %s", err))
	}

	if err = d.Set("last_modified_at", templateAssignmentResponse.LastModifiedAt); err != nil {
		return diag.FromErr(fmt.Errorf("error setting last_modified_at: %s", err))
	}

	if err = d.Set("last_modified_by_id", templateAssignmentResponse.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting last_modified_by_id: %s", err))
	}

	if err = d.Set("entity_tag", templateAssignmentResponse.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("error setting entity_tag: %s", err))
	}

	return nil
}

func dataSourceIBMAccountSettingsTemplateAssignmentResponseContextToMap(model *iamidentityv1.ResponseContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TransactionID != nil {
		modelMap["transaction_id"] = model.TransactionID
	}
	if model.Operation != nil {
		modelMap["operation"] = model.Operation
	}
	if model.UserAgent != nil {
		modelMap["user_agent"] = model.UserAgent
	}
	if model.URL != nil {
		modelMap["url"] = model.URL
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = model.InstanceID
	}
	if model.ThreadID != nil {
		modelMap["thread_id"] = model.ThreadID
	}
	if model.Host != nil {
		modelMap["host"] = model.Host
	}
	if model.StartTime != nil {
		modelMap["start_time"] = model.StartTime
	}
	if model.EndTime != nil {
		modelMap["end_time"] = model.EndTime
	}
	if model.ElapsedTime != nil {
		modelMap["elapsed_time"] = model.ElapsedTime
	}
	if model.ClusterName != nil {
		modelMap["cluster_name"] = model.ClusterName
	}
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResponseResourceToMap(model *iamidentityv1.TemplateAssignmentResponseResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target"] = model.Target
	if model.Profile != nil {
		profileMap, err := dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(model.Profile)
		if err != nil {
			return modelMap, err
		}
		modelMap["profile"] = []map[string]interface{}{profileMap}
	}
	if model.AccountSettings != nil {
		accountSettingsMap, err := dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(model.AccountSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["account_settings"] = []map[string]interface{}{accountSettingsMap}
	}
	if model.PolicyTemplateRefs != nil {
		var policyTemplateRefs []map[string]interface{}
		for _, policyTemplateRefsItem := range model.PolicyTemplateRefs {
			policyTemplateRefsItemMap, err := dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(&policyTemplateRefsItem)
			if err != nil {
				return modelMap, err
			}
			policyTemplateRefs = append(policyTemplateRefs, policyTemplateRefsItemMap)
		}
		modelMap["policy_template_refs"] = policyTemplateRefs
	}
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(model *iamidentityv1.TemplateAssignmentResponseResourceDetail) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.ResourceCreated != nil {
		resourceCreatedMap, err := dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResourceToMap(model.ResourceCreated)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_created"] = []map[string]interface{}{resourceCreatedMap}
	}
	if model.ErrorMessage != nil {
		errorMessageMap, err := dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResourceErrorToMap(model.ErrorMessage)
		if err != nil {
			return modelMap, err
		}
		modelMap["error_message"] = []map[string]interface{}{errorMessageMap}
	}
	modelMap["status"] = model.Status
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResourceToMap(model *iamidentityv1.TemplateAssignmentResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateAssignmentTemplateAssignmentResourceErrorToMap(model *iamidentityv1.TemplateAssignmentResourceError) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.ErrorCode != nil {
		modelMap["error_code"] = model.ErrorCode
	}
	if model.Message != nil {
		modelMap["message"] = model.Message
	}
	if model.StatusCode != nil {
		modelMap["status_code"] = model.StatusCode
	}
	return modelMap, nil
}

func dataSourceIBMAccountSettingsTemplateAssignmentEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = model.Timestamp
	modelMap["iam_id"] = model.IamID
	modelMap["iam_id_account"] = model.IamIDAccount
	modelMap["action"] = model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = model.Message
	return modelMap, nil
}
