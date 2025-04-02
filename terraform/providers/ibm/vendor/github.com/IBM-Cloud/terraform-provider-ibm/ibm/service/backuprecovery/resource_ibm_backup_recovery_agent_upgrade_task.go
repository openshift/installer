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
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryAgentUpgradeTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryAgentUpgradeTaskCreate,
		ReadContext:   resourceIbmBackupRecoveryAgentUpgradeTaskRead,
		DeleteContext: resourceIbmBackupRecoveryAgentUpgradeTaskDelete,
		UpdateContext: resourceIbmBackupRecoveryAgentUpgradeTaskUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryAgentUpgradeTaskCreate,
		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"retry_task_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies ID of a task which is to be retried.",
			},
			"agent_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the agents upgraded in the task.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the description of the task.",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// ForceNew:    true,
				Description: "Specifies the name of the task.",
			},
			"schedule_end_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the time before which the upgrade task should start execution as a Unix epoch Timestamp (in microseconds). If this is not specified the task will start anytime after scheduleTimeUsecs.",
			},
			"schedule_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the time when the task should start execution as a Unix epoch Timestamp (in microseconds). If no schedule is specified, the task will start immediately.",
			},
			"agents": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the upgrade information for each agent.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the ID of the agent.",
						},
						"info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the upgrade state of the agent.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the time when the upgrade for an agent completed as a Unix epoch Timestamp (in microseconds).",
									},
									"error": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Object that holds the error object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"error_code": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the error code.",
												},
												"message": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the error message.",
												},
												"task_log_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the TaskLogId of the failed task.",
												},
											},
										},
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the source where the agent is installed.",
									},
									"previous_software_version": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the software version of the agent before upgrade.",
									},
									"start_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the time when the upgrade for an agent started as a Unix epoch Timestamp (in microseconds).",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the upgrade status of the agent.<br> 'Scheduled' indicates that upgrade for the agent is yet to start.<br> 'Started' indicates that upgrade for the agent is started.<br> 'Succeeded' indicates that agent was upgraded successfully.<br> 'Failed' indicates that upgrade for the agent has failed.<br> 'Skipped' indicates that upgrade for the agent was skipped.",
									},
								},
							},
						},
					},
				},
			},
			"cluster_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the version to which agents are upgraded.",
			},
			"end_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the time when the upgrade task completed execution as a Unix epoch Timestamp (in microseconds).",
			},
			"error": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Object that holds the error object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_code": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the error code.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the error message.",
						},
						"task_log_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the TaskLogId of the failed task.",
						},
					},
				},
			},
			"is_retryable": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies if a task can be retried.",
			},
			"retried_task_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies ID of a task which was retried if type is 'Retry'.",
			},
			"start_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the time, as a Unix epoch timestamp in microseconds, when the task started execution.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the status of the task.<br> 'Scheduled' indicates that the upgrade task is yet to start.<br> 'Running' indicates that the upgrade task has started execution.<br> 'Succeeded' indicates that the upgrade task completed without an error.<br> 'Failed' indicates that upgrade has failed for all agents. 'PartiallyFailed' indicates that upgrade has failed for some agents.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifes the type of task.<br> 'Auto' indicates an auto agent upgrade task which is started after a cluster upgrade.<br> 'Manual' indicates a schedule based agent upgrade task.<br> 'Retry' indicates an agent upgrade task which was retried.",
			},
		},
	}
}

func checkDiffResourceIbmBackupRecoveryAgentUpgradeTaskCreate(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecoveryAgentUpgradeTask().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_agent_upgrade_task cannot be updated.")
		}
	}
	return nil
}

func resourceIbmBackupRecoveryAgentUpgradeTaskCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createUpgradeTaskOptions := &backuprecoveryv1.CreateUpgradeTaskOptions{}

	createUpgradeTaskOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("agent_ids"); ok {
		var agentIDs []int64
		for _, v := range d.Get("agent_ids").([]interface{}) {
			agentIDsItem := int64(v.(int))
			agentIDs = append(agentIDs, agentIDsItem)
		}
		createUpgradeTaskOptions.SetAgentIDs(agentIDs)
	}

	if _, ok := d.GetOk("description"); ok {
		createUpgradeTaskOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		createUpgradeTaskOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("retry_task_id"); ok {
		createUpgradeTaskOptions.SetRetryTaskID(int64(d.Get("retry_task_id").(int)))
	}
	if _, ok := d.GetOk("schedule_end_time_usecs"); ok {
		createUpgradeTaskOptions.SetScheduleEndTimeUsecs(int64(d.Get("schedule_end_time_usecs").(int)))
	}
	if _, ok := d.GetOk("schedule_time_usecs"); ok {
		createUpgradeTaskOptions.SetScheduleTimeUsecs(int64(d.Get("schedule_time_usecs").(int)))
	}

	agentUpgradeTaskState, _, err := backupRecoveryClient.CreateUpgradeTaskWithContext(context, createUpgradeTaskOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateUpgradeTaskWithContext failed: %s", err.Error()), "ibm_backup_recovery_agent_upgrade_task", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(strconv.Itoa(int(*agentUpgradeTaskState.ID)))

	return resourceIbmBackupRecoveryAgentUpgradeTaskRead(context, d, meta)
}

func resourceIbmBackupRecoveryAgentUpgradeTaskRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getUpgradeTasksOptions := &backuprecoveryv1.GetUpgradeTasksOptions{}

	getUpgradeTasksOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getUpgradeTasksOptions.SetIds([]int64{int64(id)})

	agentUpgradeTaskStates, response, err := backupRecoveryClient.GetUpgradeTasksWithContext(context, getUpgradeTasksOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetUpgradeTasksWithContext failed: %s", err.Error()), "ibm_backup_recovery_agent_upgrade_task", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].AgentIDs) {
		agentIDs := []interface{}{}
		for _, agentIDsItem := range agentUpgradeTaskStates.Tasks[0].AgentIDs {
			agentIDs = append(agentIDs, int64(agentIDsItem))
		}
		if err = d.Set("agent_ids", agentIDs); err != nil {
			err = fmt.Errorf("Error setting agent_ids: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-agent_ids").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].Description) {
		if err = d.Set("description", agentUpgradeTaskStates.Tasks[0].Description); err != nil {
			err = fmt.Errorf("Error setting description: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-description").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].Name) {
		if err = d.Set("name", agentUpgradeTaskStates.Tasks[0].Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].ScheduleEndTimeUsecs) {
		if err = d.Set("schedule_end_time_usecs", flex.IntValue(agentUpgradeTaskStates.Tasks[0].ScheduleEndTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting schedule_end_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-schedule_end_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].ScheduleTimeUsecs) {
		if err = d.Set("schedule_time_usecs", flex.IntValue(agentUpgradeTaskStates.Tasks[0].ScheduleTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting schedule_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-schedule_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].Agents) {
		agents := []map[string]interface{}{}
		for _, agentsItem := range agentUpgradeTaskStates.Tasks[0].Agents {
			agentsItemMap, err := ResourceIbmBackupRecoveryAgentUpgradeTaskAgentUpgradeInfoObjectToMap(&agentsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "agents-to-map").GetDiag()
			}
			agents = append(agents, agentsItemMap)
		}
		if err = d.Set("agents", agents); err != nil {
			err = fmt.Errorf("Error setting agents: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-agents").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].ClusterVersion) {
		if err = d.Set("cluster_version", agentUpgradeTaskStates.Tasks[0].ClusterVersion); err != nil {
			err = fmt.Errorf("Error setting cluster_version: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-cluster_version").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].EndTimeUsecs) {
		if err = d.Set("end_time_usecs", flex.IntValue(agentUpgradeTaskStates.Tasks[0].EndTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting end_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-end_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].Error) {
		errorMap, err := ResourceIbmBackupRecoveryAgentUpgradeTaskErrorToMap(agentUpgradeTaskStates.Tasks[0].Error)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "error-to-map").GetDiag()
		}
		if err = d.Set("error", []map[string]interface{}{errorMap}); err != nil {
			err = fmt.Errorf("Error setting error: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-error").GetDiag()
		}
	} else {
		if err = d.Set("error", []interface{}{}); err != nil {
			err = fmt.Errorf("Error setting error: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-error").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].IsRetryable) {
		if err = d.Set("is_retryable", agentUpgradeTaskStates.Tasks[0].IsRetryable); err != nil {
			err = fmt.Errorf("Error setting is_retryable: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-is_retryable").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].RetriedTaskID) {
		if err = d.Set("retried_task_id", flex.IntValue(agentUpgradeTaskStates.Tasks[0].RetriedTaskID)); err != nil {
			err = fmt.Errorf("Error setting retried_task_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-retried_task_id").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].StartTimeUsecs) {
		if err = d.Set("start_time_usecs", flex.IntValue(agentUpgradeTaskStates.Tasks[0].StartTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting start_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-start_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].Status) {
		if err = d.Set("status", agentUpgradeTaskStates.Tasks[0].Status); err != nil {
			err = fmt.Errorf("Error setting status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-status").GetDiag()
		}
	}
	if !core.IsNil(agentUpgradeTaskStates.Tasks[0].Type) {
		if err = d.Set("type", agentUpgradeTaskStates.Tasks[0].Type); err != nil {
			err = fmt.Errorf("Error setting type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_agent_upgrade_task", "read", "set-type").GetDiag()
		}
	}

	return nil
}

func resourceIbmBackupRecoveryAgentUpgradeTaskDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmBackupRecoveryAgentUpgradeTaskUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.
	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource Update Will Only Affect Terraform State",
		Detail:   "Update operation for this resource is not supported and will only affect the terraform statefile. No changes will be made to the backend resource. ",
	}
	diags = append(diags, warning)
	return diags
}

func ResourceIbmBackupRecoveryAgentUpgradeTaskAgentUpgradeInfoObjectToMap(model *backuprecoveryv1.AgentUpgradeInfoObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Info != nil {
		infoMap, err := ResourceIbmBackupRecoveryAgentUpgradeTaskAgentInfoObjectToMap(model.Info)
		if err != nil {
			return modelMap, err
		}
		modelMap["info"] = []map[string]interface{}{infoMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryAgentUpgradeTaskAgentInfoObjectToMap(model *backuprecoveryv1.AgentInfoObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Error != nil {
		errorMap, err := ResourceIbmBackupRecoveryAgentUpgradeTaskErrorToMap(model.Error)
		if err != nil {
			return modelMap, err
		}
		modelMap["error"] = []map[string]interface{}{errorMap}
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.PreviousSoftwareVersion != nil {
		modelMap["previous_software_version"] = *model.PreviousSoftwareVersion
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryAgentUpgradeTaskErrorToMap(model *backuprecoveryv1.Error) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ErrorCode != nil {
		modelMap["error_code"] = *model.ErrorCode
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.TaskLogID != nil {
		modelMap["task_log_id"] = *model.TaskLogID
	}
	return modelMap, nil
}
