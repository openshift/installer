// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package iamaccessgroup

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/platform-services-go-sdk/iamaccessgroupsv2"
)

const (
	envAGAssignmentTimeoutDurationKey = "IAM_ACCESS_GROUP_ASSIGNMENT_STATE_REFRESH_TIMEOUT_IN_SECONDS"
	InProgress                        = "in_progress"
	complete                          = "complete"
	failed                            = "failed"
)

func ResourceIBMIAMAccessGroupTemplateAssignment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIAMAccessGroupTemplateAssignmentCreate,
		ReadContext:   resourceIBMIAMAccessGroupTemplateAssignmentRead,
		UpdateContext: resourceIBMIAMAccessGroupTemplateAssignmentUpdate,
		DeleteContext: resourceIBMIAMAccessGroupTemplateAssignmentDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"transaction_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template_assignment", "transaction_id"),
				Description:  "An optional transaction id for the request.",
			},
			"template_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template_assignment", "template_id"),
				Description:  "The ID of the template that the assignment is based on.",
			},
			"template_version": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template_assignment", "template_version"),
				Description:  "The version of the template that the assignment is based on.",
			},
			"target_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template_assignment", "target_type"),
				Description:  "The type of the entity that the assignment applies to.",
			},
			"target": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_iam_access_group_template_assignment", "target"),
				Description:  "The ID of the entity that the assignment applies to.",
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the account that the assignment belongs to.",
			},
			"operation": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The operation that the assignment applies to (e.g. 'assign', 'update', 'remove').",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the assignment (e.g. 'accepted', 'in_progress', 'succeeded', 'failed', 'superseded').",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the assignment resource.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the assignment was created.",
			},
			"created_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user or system that created the assignment.",
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the assignment was last updated.",
			},
			"last_modified_by_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user or system that last updated the assignment.",
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIAMAccessGroupTemplateAssignmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "transaction_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z0-9_-]+$`,
			MinValueLength:             1,
			MaxValueLength:             50,
		},
		validate.ValidateSchema{
			Identifier:                 "template_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_-]+$`,
			MinValueLength:             1,
			MaxValueLength:             100,
		},
		validate.ValidateSchema{
			Identifier:                 "template_version",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9]+$`,
			MinValueLength:             1,
			MaxValueLength:             2,
		},
		validate.ValidateSchema{
			Identifier:                 "target_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "Account, AccountGroup",
			Regexp:                     `^[a-zA-Z-]+$`,
			MinValueLength:             7,
			MaxValueLength:             12,
		},
		validate.ValidateSchema{
			Identifier:                 "target",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9_-]+$`,
			MinValueLength:             1,
			MaxValueLength:             50,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_iam_access_group_template_assignment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIAMAccessGroupTemplateAssignmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createAssignmentOptions := &iamaccessgroupsv2.CreateAssignmentOptions{}

	createAssignmentOptions.SetTemplateID(d.Get("template_id").(string))
	createAssignmentOptions.SetTemplateVersion(d.Get("template_version").(string))
	createAssignmentOptions.SetTargetType(d.Get("target_type").(string))
	createAssignmentOptions.SetTarget(d.Get("target").(string))
	if _, ok := d.GetOk("transaction_id"); ok {
		createAssignmentOptions.SetTransactionID(d.Get("transaction_id").(string))
	}

	templateAssignmentResponse, response, err := iamAccessGroupsClient.CreateAssignmentWithContext(context, createAssignmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateAssignmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(*templateAssignmentResponse.ID)

	_, err = waitForAssignment(d.Timeout(schema.TimeoutCreate), meta, d, isAccessGroupTemplateAssigned)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error assigning %s", err))
	}

	return resourceIBMIAMAccessGroupTemplateAssignmentRead(context, d, meta)
}

func resourceIBMIAMAccessGroupTemplateAssignmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getAssignmentOptions := &iamaccessgroupsv2.GetAssignmentOptions{}

	getAssignmentOptions.SetAssignmentID(d.Id())

	templateAssignmentVerboseResponse, response, err := iamAccessGroupsClient.GetAssignmentWithContext(context, getAssignmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAssignmentWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("template_id", templateAssignmentVerboseResponse.TemplateID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting template_id: %s", err))
	}
	if err = d.Set("template_version", templateAssignmentVerboseResponse.TemplateVersion); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting template_version: %s", err))
	}
	if err = d.Set("target_type", templateAssignmentVerboseResponse.TargetType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target_type: %s", err))
	}
	if err = d.Set("target", templateAssignmentVerboseResponse.Target); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting target: %s", err))
	}
	if err = d.Set("account_id", templateAssignmentVerboseResponse.AccountID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting account_id: %s", err))
	}
	if err = d.Set("operation", templateAssignmentVerboseResponse.Operation); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting operation: %s", err))
	}
	if err = d.Set("status", templateAssignmentVerboseResponse.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("href", templateAssignmentVerboseResponse.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(templateAssignmentVerboseResponse.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("created_by_id", templateAssignmentVerboseResponse.CreatedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by_id: %s", err))
	}
	if err = d.Set("last_modified_at", flex.DateTimeToString(templateAssignmentVerboseResponse.LastModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_at: %s", err))
	}
	if err = d.Set("last_modified_by_id", templateAssignmentVerboseResponse.LastModifiedByID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_modified_by_id: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIBMIAMAccessGroupTemplateAssignmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateAssignmentOptions := &iamaccessgroupsv2.UpdateAssignmentOptions{}

	updateAssignmentOptions.SetAssignmentID(d.Id())

	getAssignmentOptions := &iamaccessgroupsv2.GetAssignmentOptions{}

	getAssignmentOptions.SetAssignmentID(d.Id())

	templateAssignmentVerboseResponse, response, err := iamAccessGroupsClient.GetAssignmentWithContext(context, getAssignmentOptions)

	if err != nil && templateAssignmentVerboseResponse == nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetAssignmentWithContext failed %s\n%s", err, response))
	}

	updateAssignmentOptions.SetIfMatch(response.Headers.Get("ETag"))

	hasChange := false

	if d.HasChange("template_version") {
		updateAssignmentOptions.SetTemplateVersion(d.Get("template_version").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := iamAccessGroupsClient.UpdateAssignmentWithContext(context, updateAssignmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateAssignmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateAssignmentWithContext failed %s\n%s", err, response))
		}
		waitForAssignment(d.Timeout(schema.TimeoutUpdate), meta, d, isAccessGroupTemplateAssigned)
	}

	return resourceIBMIAMAccessGroupTemplateAssignmentRead(context, d, meta)
}

func resourceIBMIAMAccessGroupTemplateAssignmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteAssignmentOptions := &iamaccessgroupsv2.DeleteAssignmentOptions{}

	deleteAssignmentOptions.SetAssignmentID(d.Id())

	response, err := iamAccessGroupsClient.DeleteAssignmentWithContext(context, deleteAssignmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteAssignmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteAssignmentWithContext failed %s\n%s", err, response))
	}

	waitForAssignment(d.Timeout(schema.TimeoutDelete), meta, d, isAccessGroupTemplateAssignmentDeleted)

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

func isAccessGroupTemplateAssigned(id string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
		if err != nil {
			return nil, failed, err
		}

		getAssignmentOptions := &iamaccessgroupsv2.GetAssignmentOptions{}

		getAssignmentOptions.SetAssignmentID(id)

		assignment, response, err := iamAccessGroupsClient.GetAssignment(getAssignmentOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				return nil, failed, err
			}
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

func isAccessGroupTemplateAssignmentDeleted(id string, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		iamAccessGroupsClient, err := meta.(conns.ClientSession).IAMAccessGroupsV2()
		if err != nil {
			return nil, failed, err
		}

		getAssignmentOptions := &iamaccessgroupsv2.GetAssignmentOptions{}

		getAssignmentOptions.SetAssignmentID(id)

		assignment, response, err := iamAccessGroupsClient.GetAssignment(getAssignmentOptions)

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
