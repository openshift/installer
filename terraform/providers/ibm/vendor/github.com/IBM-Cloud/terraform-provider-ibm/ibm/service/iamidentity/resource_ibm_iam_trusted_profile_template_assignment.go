// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamidentity

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/iamidentityv1"
)

const (
	InProgress = "in_progress"
	complete   = "complete"
	failed     = "failed"
)

func ResourceIBMTrustedProfileTemplateAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMTrustedProfileTemplateAssignmentCreate,
		ReadContext:   resourceIBMTrustedProfileTemplateAssignmentRead,
		UpdateContext: resourceIBMTrustedProfileTemplateAssignmentUpdate,
		DeleteContext: resourceIBMTrustedProfileTemplateAssignmentDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Template Id.",
			},
			"template_version": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Template version.",
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_trusted_profile_template_assignment", "target_type"),
				Description:  "Assignment target type.",
			},
			"target": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Assignment target.",
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
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Enterprise account Id.",
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
							Required:    true,
							Description: "Target account where the IAM resource is created.",
						},
						"profile": {
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
						"policy_template_refs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Policy resource(s) included only for trusted profile assignments with policy references.",
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
							Elem:        &schema.Schema{Type: schema.TypeString},
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

func ResourceIBMTrustedProfileTemplateAssignmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "target_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "Account, AccountGroup",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_iam_trusted_profile_template_assignment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMTrustedProfileTemplateAssignmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createTrustedProfileAssignmentOptions := &iamidentityv1.CreateTrustedProfileAssignmentOptions{}

	createTrustedProfileAssignmentOptions.SetTemplateID(d.Get("template_id").(string))
	createTrustedProfileAssignmentOptions.SetTemplateVersion(int64(d.Get("template_version").(int)))
	createTrustedProfileAssignmentOptions.SetTargetType(d.Get("target_type").(string))
	createTrustedProfileAssignmentOptions.SetTarget(d.Get("target").(string))

	templateAssignmentResponse, response, err := iamIdentityClient.CreateTrustedProfileAssignmentWithContext(context, createTrustedProfileAssignmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTrustedProfileAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTrustedProfileAssignmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(*templateAssignmentResponse.ID)

	_, err = waitForAssignment(d.Timeout(schema.TimeoutCreate), meta, d, isTrustedProfileTemplateAssigned)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error assigning %s", err))
	}

	return resourceIBMTrustedProfileTemplateAssignmentRead(context, d, meta)
}

func resourceIBMTrustedProfileTemplateAssignmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getTrustedProfileAssignmentOptions := &iamidentityv1.GetTrustedProfileAssignmentOptions{}

	getTrustedProfileAssignmentOptions.SetAssignmentID(d.Id())

	templateAssignmentResponse, response, err := iamIdentityClient.GetTrustedProfileAssignmentWithContext(context, getTrustedProfileAssignmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTrustedProfileAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTrustedProfileAssignmentWithContext failed %s\n%s", err, response))
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
	ctx := []map[string]interface{}{}
	if !core.IsNil(templateAssignmentResponse.Context) {
		contextMap, err := resourceIBMTrustedProfileTemplateAssignmentResponseContextToMap(templateAssignmentResponse.Context)
		if err != nil {
			return diag.FromErr(err)
		}
		ctx = append(ctx, contextMap)
	}
	if err = d.Set("context", ctx); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting context: %s", err))
	}
	if err = d.Set("account_id", templateAssignmentResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting account_id: %s", err))
	}
	if err = d.Set("status", templateAssignmentResponse.Status); err != nil {
		return diag.FromErr(fmt.Errorf("error setting status: %s", err))
	}
	if !core.IsNil(templateAssignmentResponse.Resources) {
		resources := []map[string]interface{}{}
		for _, resourcesItem := range templateAssignmentResponse.Resources {
			resourcesItemMap, err := resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResponseResourceToMap(&resourcesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			resources = append(resources, resourcesItemMap)
		}
		if err = d.Set("resources", resources); err != nil {
			return diag.FromErr(fmt.Errorf("error setting resources: %s", err))
		}
	}
	history := []map[string]interface{}{}
	if !core.IsNil(templateAssignmentResponse.History) {
		for _, historyItem := range templateAssignmentResponse.History {
			historyItemMap, err := resourceIBMTrustedProfileTemplateAssignmentEnityHistoryRecordToMap(&historyItem)
			if err != nil {
				return diag.FromErr(err)
			}
			history = append(history, historyItemMap)
		}
	}
	if err = d.Set("history", history); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting history: %s", err))
	}
	if !core.IsNil(templateAssignmentResponse.Href) {
		if err = d.Set("href", templateAssignmentResponse.Href); err != nil {
			return diag.FromErr(fmt.Errorf("error setting href: %s", err))
		}
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

func resourceIBMTrustedProfileTemplateAssignmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateTrustedProfileAssignmentOptions := &iamidentityv1.UpdateTrustedProfileAssignmentOptions{}
	updateTrustedProfileAssignmentOptions.SetAssignmentID(d.Id())
	updateTrustedProfileAssignmentOptions.SetIfMatch(d.Get("entity_tag").(string))

	hasChange := false

	if d.HasChange("template_version") {
		updateTrustedProfileAssignmentOptions.SetTemplateVersion(int64(d.Get("template_version").(int)))
		hasChange = true
	}

	if hasChange || d.Get("status") == "failed" { // allow the same version to retry failed assignments
		_, response, err := iamIdentityClient.UpdateTrustedProfileAssignmentWithContext(context, updateTrustedProfileAssignmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateTrustedProfileAssignmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateTrustedProfileAssignmentWithContext failed %s\n%s", err, response))
		}

		_, err = waitForAssignment(d.Timeout(schema.TimeoutUpdate), meta, d, isTrustedProfileTemplateAssigned)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error assigning %s", err))
		}
	}

	return resourceIBMTrustedProfileTemplateAssignmentRead(context, d, meta)
}

func resourceIBMTrustedProfileTemplateAssignmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTrustedProfileAssignmentOptions := &iamidentityv1.DeleteTrustedProfileAssignmentOptions{}

	deleteTrustedProfileAssignmentOptions.SetAssignmentID(d.Id())

	_, response, err := iamIdentityClient.DeleteTrustedProfileAssignmentWithContext(context, deleteTrustedProfileAssignmentOptions)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[DEBUG] DeleteTrustedProfileAssignmentWithContext failed %s\n%s", err, response))
	}

	_, err = waitForAssignment(d.Timeout(schema.TimeoutDelete), meta, d, isTrustedProfileAssignmentRemoved)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error removing assignment %s", err))
	}

	d.SetId("")

	return nil
}

func waitForAssignment(timeout time.Duration, meta interface{}, d *schema.ResourceData, refreshFn func(string, interface{}) resource.StateRefreshFunc) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{InProgress},
		Target:       []string{complete},
		Refresh:      refreshFn(d.Id(), meta),
		Delay:        30 * time.Second,
		PollInterval: time.Minute,
		Timeout:      timeout,
	}

	return stateConf.WaitForState()
}

func isTrustedProfileAssignmentRemoved(id string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()

		getOptions := &iamidentityv1.GetTrustedProfileAssignmentOptions{
			AssignmentID: &id,
		}
		assignment, response, err := iamIdentityClient.GetTrustedProfileAssignment(getOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return assignment, complete, nil
			}

			return nil, failed, fmt.Errorf("[ERROR] The assignment %s failed to delete or deletion was not completed within specific timeout period: %s\n%s", id, err, response)
		} else {
			log.Printf("Assignment removal still in progress\n")
		}
		return assignment, InProgress, nil
	}
}

func isTrustedProfileTemplateAssigned(id string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		iamIdentityClient, err := meta.(conns.ClientSession).IAMIdentityV1API()

		getOptions := &iamidentityv1.GetTrustedProfileAssignmentOptions{
			AssignmentID: &id,
		}
		assignment, response, err := iamIdentityClient.GetTrustedProfileAssignment(getOptions)
		if err != nil {
			return nil, failed, fmt.Errorf("[ERROR] The assignment %s failed or did not complete within specific timeout period: %s\n%s", id, err, response)
		}

		if assignment != nil {
			if *assignment.Status == "accepted" || *assignment.Status == "in_progress" {
				log.Printf("Assignment still in progress\n")
				return assignment, InProgress, nil
			}

			if *assignment.Status == "succeeded" {
				return assignment, complete, nil
			}

			if *assignment.Status == "failed" {
				return assignment, failed, fmt.Errorf("[ERROR] The assignment %s did complete but with a 'failed' status. Please check assignment resource for detailed errors: %s\n", id, response)
			}
		}

		return assignment, failed, fmt.Errorf("[ERROR] Unexpected status reached for assignment %s.: %s\n", id, response)
	}
}

func resourceIBMTrustedProfileTemplateAssignmentResponseContextToMap(model *iamidentityv1.ResponseContext) (map[string]interface{}, error) {
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

func resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResponseResourceToMap(model *iamidentityv1.TemplateAssignmentResponseResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target"] = model.Target
	if model.Profile != nil {
		profileMap, err := resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(model.Profile)
		if err != nil {
			return modelMap, err
		}
		modelMap["profile"] = []map[string]interface{}{profileMap}
	}
	if model.PolicyTemplateRefs != nil {
		policyTemplateRefs := []map[string]interface{}{}
		for _, policyTemplateRefsItem := range model.PolicyTemplateRefs {
			policyTemplateRefsItemMap, err := resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(&policyTemplateRefsItem)
			if err != nil {
				return modelMap, err
			}
			policyTemplateRefs = append(policyTemplateRefs, policyTemplateRefsItemMap)
		}
		modelMap["policy_template_refs"] = policyTemplateRefs
	}
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResponseResourceDetailToMap(model *iamidentityv1.TemplateAssignmentResponseResourceDetail) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.Version != nil {
		modelMap["version"] = model.Version
	}
	if model.ResourceCreated != nil {
		resourceCreatedMap, err := resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResourceToMap(model.ResourceCreated)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_created"] = []map[string]interface{}{resourceCreatedMap}
	}
	if model.ErrorMessage != nil {
		errorMessageMap, err := resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResourceErrorToMap(model.ErrorMessage)
		if err != nil {
			return modelMap, err
		}
		modelMap["error_message"] = []map[string]interface{}{errorMessageMap}
	}
	modelMap["status"] = model.Status
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResourceToMap(model *iamidentityv1.TemplateAssignmentResource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	return modelMap, nil
}

func resourceIBMTrustedProfileTemplateAssignmentTemplateAssignmentResourceErrorToMap(model *iamidentityv1.TemplateAssignmentResourceError) (map[string]interface{}, error) {
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

func resourceIBMTrustedProfileTemplateAssignmentEnityHistoryRecordToMap(model *iamidentityv1.EnityHistoryRecord) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["timestamp"] = model.Timestamp
	modelMap["iam_id"] = model.IamID
	modelMap["iam_id_account"] = model.IamIDAccount
	modelMap["action"] = model.Action
	modelMap["params"] = model.Params
	modelMap["message"] = model.Message
	return modelMap, nil
}
