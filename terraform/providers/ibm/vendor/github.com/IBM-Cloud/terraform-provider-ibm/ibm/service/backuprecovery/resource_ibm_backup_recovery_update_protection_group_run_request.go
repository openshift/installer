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
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestCreate,
		ReadContext:   resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRead,
		DeleteContext: resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestDelete,
		UpdateContext: resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdate,
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryUpdateProtectionGroupRunRequest,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"update_protection_group_run_params": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				// ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies a unique Protection Group Run id.",
						},
						"local_snapshot_config": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the params to perform actions on local snapshot taken by a Protection Group Run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_legal_hold": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies whether to retain the snapshot for legal purpose. If set to true, the snapshots cannot be deleted until the retention period. Note that using this option may cause the Cluster to run out of space. If set to false explicitly, the hold is removed, and the snapshots will expire as specified in the policy of the Protection Group. If this field is not specified, there is no change to the hold of the run. This field can be set only by a User having Data Security Role.",
									},
									"delete_snapshot": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies whether to delete the snapshot. When this is set to true, all other params will be ignored.",
									},
									"data_lock": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept until the maximum of the snapshot retention time. During that time, the snapshots cannot be deleted. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
									},
									"days_to_keep": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies number of days to retain the snapshots. If positive, then this value is added to exisiting expiry time thereby increasing  the retention period of the snapshot. Conversly, if this value is negative, then value is subtracted to existing expiry time thereby decreasing the retention period of the snaphot. Here, by this operation if expiry time goes below current time then snapshot is immediately deleted.",
									},
								},
							},
						},
						"replication_snapshot_config": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the params to perform actions on replication snapshots taken by a Protection Group Run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"new_snapshot_config": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the new configuration about adding Replication Snapshot to existing Protection Group Run.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies id of Remote Cluster to copy the Snapshots to.",
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
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
									"update_existing_snapshot_config": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the configuration about updating an existing Replication Snapshot Run.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the cluster id of the replication cluster.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the cluster name of the replication cluster.",
												},
												"enable_legal_hold": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to retain the snapshot for legal purpose. If set to true, the snapshots cannot be deleted until the retention period. Note that using this option may cause the Cluster to run out of space. If set to false explicitly, the hold is removed, and the snapshots will expire as specified in the policy of the Protection Group. If this field is not specified, there is no change to the hold of the run. This field can be set only by a User having Data Security Role.",
												},
												"delete_snapshot": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to delete the snapshot. When this is set to true, all other params will be ignored.",
												},
												"resync": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to retry the replication operation in case if earlier attempt failed. If not specified or set to false, replication is not retried.",
												},
												"data_lock": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept until the maximum of the snapshot retention time. During that time, the snapshots cannot be deleted. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
												},
												"days_to_keep": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies number of days to retain the snapshots. If positive, then this value is added to exisiting expiry time thereby increasing  the retention period of the snapshot. Conversly, if this value is negative, then value is subtracted to existing expiry time thereby decreasing the retention period of the snaphot. Here, by this operation if expiry time goes below current time then snapshot is immediately deleted.",
												},
											},
										},
									},
								},
							},
						},
						"archival_snapshot_config": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the params to perform actions on archival snapshots taken by a Protection Group Run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"new_snapshot_config": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the new configuration about adding Archival Snapshot to existing Protection Group Run.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the Archival target to copy the Snapshots to.",
												},
												"archival_target_type": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the snapshot's archival target type from which recovery has been performed.",
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"copy_only_fully_successful": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies if Snapshots are copied from a fully successful Protection Group Run or a partially successful Protection Group Run. If false, Snapshots are copied the Protection Group Run, even if the Run was not fully successful i.e. Snapshots were not captured for all Objects in the Protection Group. If true, Snapshots are copied only when the run is fully successful.",
												},
											},
										},
									},
									"update_existing_snapshot_config": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the configuration about updating an existing Archival Snapshot Run.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the id of the archival target.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the archival target.",
												},
												"archival_target_type": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the snapshot's archival target type from which recovery has been performed.",
												},
												"enable_legal_hold": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to retain the snapshot for legal purpose. If set to true, the snapshots cannot be deleted until the retention period. Note that using this option may cause the Cluster to run out of space. If set to false explicitly, the hold is removed, and the snapshots will expire as specified in the policy of the Protection Group. If this field is not specified, there is no change to the hold of the run. This field can be set only by a User having Data Security Role.",
												},
												"delete_snapshot": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to delete the snapshot. When this is set to true, all other params will be ignored.",
												},
												"resync": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to retry the archival operation in case if earlier attempt failed. If not specified or set to false, archival is not retried.",
												},
												"data_lock": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept until the maximum of the snapshot retention time. During that time, the snapshots cannot be deleted. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
												},
												"days_to_keep": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies number of days to retain the snapshots. If positive, then this value is added to exisiting expiry time thereby increasing  the retention period of the snapshot. Conversly, if this value is negative, then value is subtracted to existing expiry time thereby decreasing the retention period of the snaphot. Here, by this operation if expiry time goes below current time then snapshot is immediately deleted.",
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
			"run_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID.",
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc: validate.InvokeValidator("ibm_create_protection_group_run_request", "run_type"),
				Description: "Protection group id",
			},
			"successful_run_ids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ForceNew:    true,
				Description: "Specifies a list of Protection Group ids for which the state should change.",
			},
			"failed_runs": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				// ForceNew:    true,
				Description: "Specfies the list of connections for the source.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"run_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the id of the connection.",
						},
						"error_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the entity id of the source. The source can a non-root entity.",
						},
					},
				},
			},
		},
	}
}

func checkDiffResourceIbmBackupRecoveryUpdateProtectionGroupRunRequest(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequest().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_update_protection_group_run_request cannot be updated.")
		}
	}
	return nil
}

func resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_update_protection_group_run_request", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateProtectionGroupRunOptions := &backuprecoveryv1.UpdateProtectionGroupRunOptions{}

	updateProtectionGroupRunOptions.SetID(d.Get("group_id").(string))
	updateProtectionGroupRunOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	var newUpdateProtectionGroupRunParams []backuprecoveryv1.UpdateProtectionGroupRunParams
	for _, v := range d.Get("update_protection_group_run_params").([]interface{}) {
		value := v.(map[string]interface{})
		newUpdateProtectionGroupRunParamsItem, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateProtectionGroupRunParams(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_update_protection_group_run_request", "create", "parse-update_protection_group_run_params").GetDiag()
		}
		newUpdateProtectionGroupRunParams = append(newUpdateProtectionGroupRunParams, *newUpdateProtectionGroupRunParamsItem)
	}
	updateProtectionGroupRunOptions.SetUpdateProtectionGroupRunParams(newUpdateProtectionGroupRunParams)

	updateProtectionGroupRunResponse, _, err := backupRecoveryClient.UpdateProtectionGroupRunWithContext(context, updateProtectionGroupRunOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateProtectionGroupRunWithContext failed: %s", err.Error()), "ibm_backup_recovery_update_protection_group_run_request", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.Set("successful_run_ids", strings.Join(updateProtectionGroupRunResponse.SuccessfulRunIds[:], ","))

	if !core.IsNil(updateProtectionGroupRunResponse.FailedRuns) {
		failedRuns := []map[string]interface{}{}
		for _, failedRun := range updateProtectionGroupRunResponse.FailedRuns {
			failedRunsMap, err := resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateProtectionGroupRunFailedRuns(&failedRun)
			if err != nil {
				return diag.FromErr(err)
			}
			failedRuns = append(failedRuns, failedRunsMap)
		}
		if err = d.Set("failed_runs", failedRuns); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting failedRuns: %s", err))
		}
	}

	d.SetId(*updateProtectionGroupRunOptions.ID)

	return resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRead(context, d, meta)
}

func resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateProtectionGroupRunFailedRuns(model *backuprecoveryv1.FailedRunDetails) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RunID != nil {
		modelMap["run_id"] = model.RunID
	}
	if model.ErrorMessage != nil {
		modelMap["error_message"] = model.ErrorMessage
	}
	return modelMap, nil
}

func resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "delete" operation.
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

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateProtectionGroupRunParams(modelMap map[string]interface{}) (*backuprecoveryv1.UpdateProtectionGroupRunParams, error) {
	model := &backuprecoveryv1.UpdateProtectionGroupRunParams{}
	model.RunID = core.StringPtr(modelMap["run_id"].(string))
	if modelMap["local_snapshot_config"] != nil && len(modelMap["local_snapshot_config"].([]interface{})) > 0 {
		LocalSnapshotConfigModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateLocalSnapshotConfig(modelMap["local_snapshot_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LocalSnapshotConfig = LocalSnapshotConfigModel
	}
	if modelMap["replication_snapshot_config"] != nil && len(modelMap["replication_snapshot_config"].([]interface{})) > 0 {
		ReplicationSnapshotConfigModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateReplicationSnapshotConfig(modelMap["replication_snapshot_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ReplicationSnapshotConfig = ReplicationSnapshotConfigModel
	}
	if modelMap["archival_snapshot_config"] != nil && len(modelMap["archival_snapshot_config"].([]interface{})) > 0 {
		ArchivalSnapshotConfigModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateArchivalSnapshotConfig(modelMap["archival_snapshot_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ArchivalSnapshotConfig = ArchivalSnapshotConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateLocalSnapshotConfig(modelMap map[string]interface{}) (*backuprecoveryv1.UpdateLocalSnapshotConfig, error) {
	model := &backuprecoveryv1.UpdateLocalSnapshotConfig{}
	if modelMap["enable_legal_hold"] != nil {
		model.EnableLegalHold = core.BoolPtr(modelMap["enable_legal_hold"].(bool))
	}
	if modelMap["delete_snapshot"] != nil {
		model.DeleteSnapshot = core.BoolPtr(modelMap["delete_snapshot"].(bool))
	}
	if modelMap["data_lock"] != nil && modelMap["data_lock"].(string) != "" {
		model.DataLock = core.StringPtr(modelMap["data_lock"].(string))
	}
	if modelMap["days_to_keep"] != nil {
		model.DaysToKeep = core.Int64Ptr(int64(modelMap["days_to_keep"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateReplicationSnapshotConfig(modelMap map[string]interface{}) (*backuprecoveryv1.UpdateReplicationSnapshotConfig, error) {
	model := &backuprecoveryv1.UpdateReplicationSnapshotConfig{}
	if modelMap["new_snapshot_config"] != nil {
		newSnapshotConfig := []backuprecoveryv1.RunReplicationConfig{}
		for _, newSnapshotConfigItem := range modelMap["new_snapshot_config"].([]interface{}) {
			newSnapshotConfigItemModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRunReplicationConfig(newSnapshotConfigItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			newSnapshotConfig = append(newSnapshotConfig, *newSnapshotConfigItemModel)
		}
		model.NewSnapshotConfig = newSnapshotConfig
	}
	if modelMap["update_existing_snapshot_config"] != nil {
		updateExistingSnapshotConfig := []backuprecoveryv1.UpdateExistingReplicationSnapshotConfig{}
		for _, updateExistingSnapshotConfigItem := range modelMap["update_existing_snapshot_config"].([]interface{}) {
			updateExistingSnapshotConfigItemModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateExistingReplicationSnapshotConfig(updateExistingSnapshotConfigItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			updateExistingSnapshotConfig = append(updateExistingSnapshotConfig, *updateExistingSnapshotConfigItemModel)
		}
		model.UpdateExistingSnapshotConfig = updateExistingSnapshotConfig
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRunReplicationConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RunReplicationConfig, error) {
	model := &backuprecoveryv1.RunReplicationConfig{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["retention"] != nil && len(modelMap["retention"].([]interface{})) > 0 {
		RetentionModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Retention = RetentionModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRetention(modelMap map[string]interface{}) (*backuprecoveryv1.Retention, error) {
	model := &backuprecoveryv1.Retention{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["data_lock_config"] != nil && len(modelMap["data_lock_config"].([]interface{})) > 0 {
		DataLockConfigModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToDataLockConfig(modelMap["data_lock_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DataLockConfig = DataLockConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToDataLockConfig(modelMap map[string]interface{}) (*backuprecoveryv1.DataLockConfig, error) {
	model := &backuprecoveryv1.DataLockConfig{}
	model.Mode = core.StringPtr(modelMap["mode"].(string))
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["enable_worm_on_external_target"] != nil {
		model.EnableWormOnExternalTarget = core.BoolPtr(modelMap["enable_worm_on_external_target"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateExistingReplicationSnapshotConfig(modelMap map[string]interface{}) (*backuprecoveryv1.UpdateExistingReplicationSnapshotConfig, error) {
	model := &backuprecoveryv1.UpdateExistingReplicationSnapshotConfig{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["enable_legal_hold"] != nil {
		model.EnableLegalHold = core.BoolPtr(modelMap["enable_legal_hold"].(bool))
	}
	if modelMap["delete_snapshot"] != nil {
		model.DeleteSnapshot = core.BoolPtr(modelMap["delete_snapshot"].(bool))
	}
	if modelMap["resync"] != nil {
		model.Resync = core.BoolPtr(modelMap["resync"].(bool))
	}
	if modelMap["data_lock"] != nil && modelMap["data_lock"].(string) != "" {
		model.DataLock = core.StringPtr(modelMap["data_lock"].(string))
	}
	if modelMap["days_to_keep"] != nil {
		model.DaysToKeep = core.Int64Ptr(int64(modelMap["days_to_keep"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateArchivalSnapshotConfig(modelMap map[string]interface{}) (*backuprecoveryv1.UpdateArchivalSnapshotConfig, error) {
	model := &backuprecoveryv1.UpdateArchivalSnapshotConfig{}
	if modelMap["new_snapshot_config"] != nil {
		newSnapshotConfig := []backuprecoveryv1.RunArchivalConfig{}
		for _, newSnapshotConfigItem := range modelMap["new_snapshot_config"].([]interface{}) {
			newSnapshotConfigItemModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRunArchivalConfig(newSnapshotConfigItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			newSnapshotConfig = append(newSnapshotConfig, *newSnapshotConfigItemModel)
		}
		model.NewSnapshotConfig = newSnapshotConfig
	}
	if modelMap["update_existing_snapshot_config"] != nil {
		updateExistingSnapshotConfig := []backuprecoveryv1.UpdateExistingArchivalSnapshotConfig{}
		for _, updateExistingSnapshotConfigItem := range modelMap["update_existing_snapshot_config"].([]interface{}) {
			updateExistingSnapshotConfigItemModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateExistingArchivalSnapshotConfig(updateExistingSnapshotConfigItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			updateExistingSnapshotConfig = append(updateExistingSnapshotConfig, *updateExistingSnapshotConfigItemModel)
		}
		model.UpdateExistingSnapshotConfig = updateExistingSnapshotConfig
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRunArchivalConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RunArchivalConfig, error) {
	model := &backuprecoveryv1.RunArchivalConfig{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	model.ArchivalTargetType = core.StringPtr(modelMap["archival_target_type"].(string))
	if modelMap["retention"] != nil && len(modelMap["retention"].([]interface{})) > 0 {
		RetentionModel, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Retention = RetentionModel
	}
	if modelMap["copy_only_fully_successful"] != nil {
		model.CopyOnlyFullySuccessful = core.BoolPtr(modelMap["copy_only_fully_successful"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestMapToUpdateExistingArchivalSnapshotConfig(modelMap map[string]interface{}) (*backuprecoveryv1.UpdateExistingArchivalSnapshotConfig, error) {
	model := &backuprecoveryv1.UpdateExistingArchivalSnapshotConfig{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.ArchivalTargetType = core.StringPtr(modelMap["archival_target_type"].(string))
	if modelMap["enable_legal_hold"] != nil {
		model.EnableLegalHold = core.BoolPtr(modelMap["enable_legal_hold"].(bool))
	}
	if modelMap["delete_snapshot"] != nil {
		model.DeleteSnapshot = core.BoolPtr(modelMap["delete_snapshot"].(bool))
	}
	if modelMap["resync"] != nil {
		model.Resync = core.BoolPtr(modelMap["resync"].(bool))
	}
	if modelMap["data_lock"] != nil && modelMap["data_lock"].(string) != "" {
		model.DataLock = core.StringPtr(modelMap["data_lock"].(string))
	}
	if modelMap["days_to_keep"] != nil {
		model.DaysToKeep = core.Int64Ptr(int64(modelMap["days_to_keep"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateProtectionGroupRunParamsToMap(model *backuprecoveryv1.UpdateProtectionGroupRunParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["run_id"] = *model.RunID
	if model.LocalSnapshotConfig != nil {
		localSnapshotConfigMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateLocalSnapshotConfigToMap(model.LocalSnapshotConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_snapshot_config"] = []map[string]interface{}{localSnapshotConfigMap}
	}
	if model.ReplicationSnapshotConfig != nil {
		replicationSnapshotConfigMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateReplicationSnapshotConfigToMap(model.ReplicationSnapshotConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["replication_snapshot_config"] = []map[string]interface{}{replicationSnapshotConfigMap}
	}
	if model.ArchivalSnapshotConfig != nil {
		archivalSnapshotConfigMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateArchivalSnapshotConfigToMap(model.ArchivalSnapshotConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_snapshot_config"] = []map[string]interface{}{archivalSnapshotConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateLocalSnapshotConfigToMap(model *backuprecoveryv1.UpdateLocalSnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableLegalHold != nil {
		modelMap["enable_legal_hold"] = *model.EnableLegalHold
	}
	if model.DeleteSnapshot != nil {
		modelMap["delete_snapshot"] = *model.DeleteSnapshot
	}
	if model.DataLock != nil {
		modelMap["data_lock"] = *model.DataLock
	}
	if model.DaysToKeep != nil {
		modelMap["days_to_keep"] = flex.IntValue(model.DaysToKeep)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateReplicationSnapshotConfigToMap(model *backuprecoveryv1.UpdateReplicationSnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewSnapshotConfig != nil {
		newSnapshotConfig := []map[string]interface{}{}
		for _, newSnapshotConfigItem := range model.NewSnapshotConfig {
			newSnapshotConfigItemMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRunReplicationConfigToMap(&newSnapshotConfigItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			newSnapshotConfig = append(newSnapshotConfig, newSnapshotConfigItemMap)
		}
		modelMap["new_snapshot_config"] = newSnapshotConfig
	}
	if model.UpdateExistingSnapshotConfig != nil {
		updateExistingSnapshotConfig := []map[string]interface{}{}
		for _, updateExistingSnapshotConfigItem := range model.UpdateExistingSnapshotConfig {
			updateExistingSnapshotConfigItemMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateExistingReplicationSnapshotConfigToMap(&updateExistingSnapshotConfigItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			updateExistingSnapshotConfig = append(updateExistingSnapshotConfig, updateExistingSnapshotConfigItemMap)
		}
		modelMap["update_existing_snapshot_config"] = updateExistingSnapshotConfig
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRunReplicationConfigToMap(model *backuprecoveryv1.RunReplicationConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Retention != nil {
		retentionMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRetentionToMap(model.Retention)
		if err != nil {
			return modelMap, err
		}
		modelMap["retention"] = []map[string]interface{}{retentionMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRetentionToMap(model *backuprecoveryv1.Retention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestDataLockConfigToMap(model *backuprecoveryv1.DataLockConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = *model.Mode
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.EnableWormOnExternalTarget != nil {
		modelMap["enable_worm_on_external_target"] = *model.EnableWormOnExternalTarget
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateExistingReplicationSnapshotConfigToMap(model *backuprecoveryv1.UpdateExistingReplicationSnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.EnableLegalHold != nil {
		modelMap["enable_legal_hold"] = *model.EnableLegalHold
	}
	if model.DeleteSnapshot != nil {
		modelMap["delete_snapshot"] = *model.DeleteSnapshot
	}
	if model.Resync != nil {
		modelMap["resync"] = *model.Resync
	}
	if model.DataLock != nil {
		modelMap["data_lock"] = *model.DataLock
	}
	if model.DaysToKeep != nil {
		modelMap["days_to_keep"] = flex.IntValue(model.DaysToKeep)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateArchivalSnapshotConfigToMap(model *backuprecoveryv1.UpdateArchivalSnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewSnapshotConfig != nil {
		newSnapshotConfig := []map[string]interface{}{}
		for _, newSnapshotConfigItem := range model.NewSnapshotConfig {
			newSnapshotConfigItemMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRunArchivalConfigToMap(&newSnapshotConfigItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			newSnapshotConfig = append(newSnapshotConfig, newSnapshotConfigItemMap)
		}
		modelMap["new_snapshot_config"] = newSnapshotConfig
	}
	if model.UpdateExistingSnapshotConfig != nil {
		updateExistingSnapshotConfig := []map[string]interface{}{}
		for _, updateExistingSnapshotConfigItem := range model.UpdateExistingSnapshotConfig {
			updateExistingSnapshotConfigItemMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateExistingArchivalSnapshotConfigToMap(&updateExistingSnapshotConfigItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			updateExistingSnapshotConfig = append(updateExistingSnapshotConfig, updateExistingSnapshotConfigItemMap)
		}
		modelMap["update_existing_snapshot_config"] = updateExistingSnapshotConfig
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRunArchivalConfigToMap(model *backuprecoveryv1.RunArchivalConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	modelMap["archival_target_type"] = *model.ArchivalTargetType
	if model.Retention != nil {
		retentionMap, err := ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestRetentionToMap(model.Retention)
		if err != nil {
			return modelMap, err
		}
		modelMap["retention"] = []map[string]interface{}{retentionMap}
	}
	if model.CopyOnlyFullySuccessful != nil {
		modelMap["copy_only_fully_successful"] = *model.CopyOnlyFullySuccessful
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryUpdateProtectionGroupRunRequestUpdateExistingArchivalSnapshotConfigToMap(model *backuprecoveryv1.UpdateExistingArchivalSnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	modelMap["archival_target_type"] = *model.ArchivalTargetType
	if model.EnableLegalHold != nil {
		modelMap["enable_legal_hold"] = *model.EnableLegalHold
	}
	if model.DeleteSnapshot != nil {
		modelMap["delete_snapshot"] = *model.DeleteSnapshot
	}
	if model.Resync != nil {
		modelMap["resync"] = *model.Resync
	}
	if model.DataLock != nil {
		modelMap["data_lock"] = *model.DataLock
	}
	if model.DaysToKeep != nil {
		modelMap["days_to_keep"] = flex.IntValue(model.DaysToKeep)
	}
	return modelMap, nil
}
