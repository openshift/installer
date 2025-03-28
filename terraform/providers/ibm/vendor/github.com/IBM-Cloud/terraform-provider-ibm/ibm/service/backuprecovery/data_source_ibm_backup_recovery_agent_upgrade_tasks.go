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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func DataSourceIbmBackupRecoveryAgentUpgradeTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryAgentUpgradeTasksRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies IDs of tasks to be fetched.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"tasks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of agent upgrade tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_i_ds": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the agents upgraded in the task.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"agents": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the upgrade information for each agent.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the ID of the agent.",
									},
									"info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the upgrade state of the agent.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"end_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the time when the upgrade for an agent completed as a Unix epoch Timestamp (in microseconds).",
												},
												"error": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Object that holds the error object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"error_code": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the error code.",
															},
															"message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the error message.",
															},
															"task_log_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the TaskLogId of the failed task.",
															},
														},
													},
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the source where the agent is installed.",
												},
												"previous_software_version": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the software version of the agent before upgrade.",
												},
												"start_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the time when the upgrade for an agent started as a Unix epoch Timestamp (in microseconds).",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
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
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the description of the task.",
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
										Computed:    true,
										Description: "Specifies the error code.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the error message.",
									},
									"task_log_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the TaskLogId of the failed task.",
									},
								},
							},
						},
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the ID of the task.",
						},
						"is_retryable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if a task can be retried.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the name of the task.",
						},
						"retried_task_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies ID of a task which was retried if type is 'Retry'.",
						},
						"schedule_end_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the time before which the upgrade task should start execution as a Unix epoch Timestamp (in microseconds). If this is not specified the task will start anytime after scheduleTimeUsecs.",
						},
						"schedule_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the time when the task should start execution as a Unix epoch Timestamp (in microseconds). If no schedule is specified, the task will start immediately.",
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
				},
			},
		},
	}
}

func dataSourceIbmBackupRecoveryAgentUpgradeTasksRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_agent_upgrade_tasks", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getUpgradeTasksOptions := &backuprecoveryv1.GetUpgradeTasksOptions{}

	getUpgradeTasksOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("ids"); ok {
		var ids []int64
		for _, v := range d.Get("ids").([]interface{}) {
			idsItem := int64(v.(int))
			ids = append(ids, idsItem)
		}
		getUpgradeTasksOptions.SetIds(ids)
	}

	agentUpgradeTaskStates, _, err := backupRecoveryClient.GetUpgradeTasksWithContext(context, getUpgradeTasksOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetUpgradeTasksWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_agent_upgrade_tasks", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryAgentUpgradeTasksID(d))

	if !core.IsNil(agentUpgradeTaskStates.Tasks) {
		tasks := []map[string]interface{}{}
		for _, tasksItem := range agentUpgradeTaskStates.Tasks {
			tasksItemMap, err := DataSourceIbmBackupRecoveryAgentUpgradeTasksAgentUpgradeTaskStateToMap(&tasksItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_agent_upgrade_tasks", "read", "tasks-to-map").GetDiag()
			}
			tasks = append(tasks, tasksItemMap)
		}
		if err = d.Set("tasks", tasks); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tasks: %s", err), "(Data) ibm_backup_recovery_agent_upgrade_tasks", "read", "set-tasks").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryAgentUpgradeTasksID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryAgentUpgradeTasksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryAgentUpgradeTasksAgentUpgradeTaskStateToMap(model *backuprecoveryv1.AgentUpgradeTaskState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AgentIDs != nil {
		modelMap["agent_i_ds"] = model.AgentIDs
	}
	if model.Agents != nil {
		agents := []map[string]interface{}{}
		for _, agentsItem := range model.Agents {
			agentsItemMap, err := DataSourceIbmBackupRecoveryAgentUpgradeTasksAgentUpgradeInfoObjectToMap(&agentsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			agents = append(agents, agentsItemMap)
		}
		modelMap["agents"] = agents
	}
	if model.ClusterVersion != nil {
		modelMap["cluster_version"] = *model.ClusterVersion
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Error != nil {
		errorMap, err := DataSourceIbmBackupRecoveryAgentUpgradeTasksErrorToMap(model.Error)
		if err != nil {
			return modelMap, err
		}
		modelMap["error"] = []map[string]interface{}{errorMap}
	}
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.IsRetryable != nil {
		modelMap["is_retryable"] = *model.IsRetryable
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.RetriedTaskID != nil {
		modelMap["retried_task_id"] = flex.IntValue(model.RetriedTaskID)
	}
	if model.ScheduleEndTimeUsecs != nil {
		modelMap["schedule_end_time_usecs"] = flex.IntValue(model.ScheduleEndTimeUsecs)
	}
	if model.ScheduleTimeUsecs != nil {
		modelMap["schedule_time_usecs"] = flex.IntValue(model.ScheduleTimeUsecs)
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryAgentUpgradeTasksAgentUpgradeInfoObjectToMap(model *backuprecoveryv1.AgentUpgradeInfoObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Info != nil {
		infoMap, err := DataSourceIbmBackupRecoveryAgentUpgradeTasksAgentInfoObjectToMap(model.Info)
		if err != nil {
			return modelMap, err
		}
		modelMap["info"] = []map[string]interface{}{infoMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryAgentUpgradeTasksAgentInfoObjectToMap(model *backuprecoveryv1.AgentInfoObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Error != nil {
		errorMap, err := DataSourceIbmBackupRecoveryAgentUpgradeTasksErrorToMap(model.Error)
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

func DataSourceIbmBackupRecoveryAgentUpgradeTasksErrorToMap(model *backuprecoveryv1.Error) (map[string]interface{}, error) {
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
