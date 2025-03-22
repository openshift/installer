// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.0-fa797aec-20240814-142622
 */

package backuprecovery

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestCreate,
		ReadContext:   resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestRead,
		DeleteContext: resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestDelete,
		UpdateContext: resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryPerformActionOnProtectionGroupRun,
		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_perform_action_on_protection_group_run_request", "action"),
				Description: "Specifies the type of the action which will be performed on protection runs.",
			},
			"pause_params": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the pause action params for a protection run.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies a unique run id of the Protection Group run.",
						},
					},
				},
			},
			"resume_params": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the resume action params for a protection run.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies a unique run id of the Protection Group run.",
						},
					},
				},
			},
			"cancel_params": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the cancel action params for a protection run.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies a unique run id of the Protection Group run.",
						},
						"local_task_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the task id of the local run.",
						},
						"object_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of entity ids for which we need to cancel the backup tasks. If this is provided it will not cancel the complete run but will cancel only subset of backup tasks (if backup tasks are cancelled correspoding copy task will also get cancelled). If the backup tasks are completed successfully it will not cancel those backup tasks.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"replication_task_id": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the task id of the replication run.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"archival_task_id": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the task id of the archival run.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"cloud_spin_task_id": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the task id of the cloudSpin run.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc: validate.InvokeValidator("ibm_create_protection_group_run_request", "run_type"),
				Description: "Protection group id",
			},
		},
	}
}

func checkDiffResourceIbmBackupRecoveryPerformActionOnProtectionGroupRun(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequest().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_perform_action_on_protection_group_run_request cannot be updated. %s cannot be updated", fieldName)
		}
	}
	return nil
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "action",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "Cancel, Pause, Resume",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_backup_recovery_perform_action_on_protection_group_run_request", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_perform_action_on_protection_group_run_request", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	performActionOnProtectionGroupRunOptions := &backuprecoveryv1.PerformActionOnProtectionGroupRunOptions{}
	performActionOnProtectionGroupRunOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	performActionOnProtectionGroupRunOptions.SetID(d.Get("group_id").(string))
	performActionOnProtectionGroupRunOptions.SetAction(d.Get("action").(string))
	if _, ok := d.GetOk("pause_params"); ok {
		var newPauseParams []backuprecoveryv1.PauseProtectionRunActionParams
		for _, v := range d.Get("pause_params").([]interface{}) {
			value := v.(map[string]interface{})
			newPauseParamsItem, err := ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestMapToPauseProtectionRunActionParams(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_perform_action_on_protection_group_run_request", "create", "parse-pause_params").GetDiag()
			}
			newPauseParams = append(newPauseParams, *newPauseParamsItem)
		}
		performActionOnProtectionGroupRunOptions.SetPauseParams(newPauseParams)
	}
	if _, ok := d.GetOk("resume_params"); ok {
		var newResumeParams []backuprecoveryv1.ResumeProtectionRunActionParams
		for _, v := range d.Get("resume_params").([]interface{}) {
			value := v.(map[string]interface{})
			newResumeParamsItem, err := ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestMapToResumeProtectionRunActionParams(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_perform_action_on_protection_group_run_request", "create", "parse-resume_params").GetDiag()
			}
			newResumeParams = append(newResumeParams, *newResumeParamsItem)
		}
		performActionOnProtectionGroupRunOptions.SetResumeParams(newResumeParams)
	}
	if _, ok := d.GetOk("cancel_params"); ok {
		var newCancelParams []backuprecoveryv1.CancelProtectionGroupRunRequest
		for _, v := range d.Get("cancel_params").([]interface{}) {
			value := v.(map[string]interface{})
			newCancelParamsItem, err := ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestMapToCancelProtectionGroupRunRequest(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_perform_action_on_protection_group_run_request", "create", "parse-cancel_params").GetDiag()
			}
			newCancelParams = append(newCancelParams, *newCancelParamsItem)
		}
		performActionOnProtectionGroupRunOptions.SetCancelParams(newCancelParams)
	}

	performRunActionResponse, _, err := backupRecoveryClient.PerformActionOnProtectionGroupRunWithContext(context, performActionOnProtectionGroupRunOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("PerformActionOnProtectionGroupRunWithContext failed: %s", err.Error()), "ibm_perform_action_on_protection_group_run_request", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(resourceIbmBackupRecoveryProtectionRunActionID(d))

	d.Set("action", performRunActionResponse.Action)

	if !core.IsNil(performRunActionResponse.PauseParams) {
		pauseParams := []map[string]interface{}{}
		for _, pauseParamsItem := range performRunActionResponse.PauseParams {
			pauseParamsItemMap, err := ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestPauseProtectionRunActionParamsToMap(&pauseParamsItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_perform_action_on_protection_group_run_request", "read", "pause_params-to-map").GetDiag()
			}
			pauseParams = append(pauseParams, pauseParamsItemMap)
		}
		if err = d.Set("pause_params", pauseParams); err != nil {
			err = fmt.Errorf("Error setting pause_params: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_perform_action_on_protection_group_run_request", "read", "set-pause_params").GetDiag()
		}
	}
	if !core.IsNil(performRunActionResponse.ResumeParams) {
		resumeParams := []map[string]interface{}{}
		for _, resumeParamsItem := range performRunActionResponse.ResumeParams {
			resumeParamsItemMap, err := ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestResumeProtectionRunActionParamsToMap(&resumeParamsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_perform_action_on_protection_group_run_request", "read", "resume_params-to-map").GetDiag()
			}
			resumeParams = append(resumeParams, resumeParamsItemMap)
		}
		_ = d.Set("resume_params", []map[string]interface{}{})
		if err = d.Set("resume_params", resumeParams); err != nil {
			err = fmt.Errorf("Error setting resume_params: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_perform_action_on_protection_group_run_request", "read", "set-resume_params").GetDiag()
		}
	}
	return resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestRead(context, d, meta)
}

func resourceIbmBackupRecoveryProtectionRunActionID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.

	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Delete Not Supported",
		Detail:   "The resource definition will be only be removed from the terraform statefile. This resource cannot be deleted from the backend. ",
	}
	diags = append(diags, warning)
	d.SetId("")
	return diags
}

func resourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "update" operation.
	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource update will only affect terraform state and not the actual backend resource",
		Detail:   "Update operation for this resource is not supported and will only affect the terraform statefile. No changes will be made to the backend resource.",
	}
	diags = append(diags, warning)
	// d.SetId("")
	return diags
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestMapToPauseProtectionRunActionParams(modelMap map[string]interface{}) (*backuprecoveryv1.PauseProtectionRunActionParams, error) {
	model := &backuprecoveryv1.PauseProtectionRunActionParams{}
	model.RunID = core.StringPtr(modelMap["run_id"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestMapToResumeProtectionRunActionParams(modelMap map[string]interface{}) (*backuprecoveryv1.ResumeProtectionRunActionParams, error) {
	model := &backuprecoveryv1.ResumeProtectionRunActionParams{}
	model.RunID = core.StringPtr(modelMap["run_id"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestMapToCancelProtectionGroupRunRequest(modelMap map[string]interface{}) (*backuprecoveryv1.CancelProtectionGroupRunRequest, error) {
	model := &backuprecoveryv1.CancelProtectionGroupRunRequest{}
	model.RunID = core.StringPtr(modelMap["run_id"].(string))
	if modelMap["local_task_id"] != nil && modelMap["local_task_id"].(string) != "" {
		model.LocalTaskID = core.StringPtr(modelMap["local_task_id"].(string))
	}
	if modelMap["object_ids"] != nil {
		objectIds := []int64{}
		for _, objectIdsItem := range modelMap["object_ids"].([]interface{}) {
			objectIds = append(objectIds, int64(objectIdsItem.(int)))
		}
		model.ObjectIds = objectIds
	}
	if modelMap["replication_task_id"] != nil {
		replicationTaskID := []string{}
		for _, replicationTaskIDItem := range modelMap["replication_task_id"].([]interface{}) {
			replicationTaskID = append(replicationTaskID, replicationTaskIDItem.(string))
		}
		model.ReplicationTaskID = replicationTaskID
	}
	if modelMap["archival_task_id"] != nil {
		archivalTaskID := []string{}
		for _, archivalTaskIDItem := range modelMap["archival_task_id"].([]interface{}) {
			archivalTaskID = append(archivalTaskID, archivalTaskIDItem.(string))
		}
		model.ArchivalTaskID = archivalTaskID
	}
	if modelMap["cloud_spin_task_id"] != nil {
		cloudSpinTaskID := []string{}
		for _, cloudSpinTaskIDItem := range modelMap["cloud_spin_task_id"].([]interface{}) {
			cloudSpinTaskID = append(cloudSpinTaskID, cloudSpinTaskIDItem.(string))
		}
		model.CloudSpinTaskID = cloudSpinTaskID
	}
	return model, nil
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestPauseProtectionRunActionParamsToMap(model *backuprecoveryv1.PauseProtectionRunActionResponseParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["run_id"] = *model.RunID
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestResumeProtectionRunActionParamsToMap(model *backuprecoveryv1.ResumeProtectionRunActionResponseParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["run_id"] = *model.RunID
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPerformActionOnProtectionGroupRunRequestCancelProtectionGroupRunRequestToMap(model *backuprecoveryv1.CancelProtectionGroupRunResponseParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["run_id"] = *model.RunID
	return modelMap, nil
}
