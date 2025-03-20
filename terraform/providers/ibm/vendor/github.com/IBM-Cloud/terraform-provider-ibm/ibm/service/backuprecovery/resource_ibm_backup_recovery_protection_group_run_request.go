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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryProtectionGroupRunRequest() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryProtectionGroupRunRequestCreate,
		ReadContext:   resourceIbmBackupRecoveryProtectionGroupRunRequestRead,
		DeleteContext: resourceIbmBackupRecoveryProtectionGroupRunRequestDelete,
		UpdateContext: resourceIbmBackupRecoveryProtectionGroupRunRequestUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryProtectionGroupRun,
		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc: validate.InvokeValidator("ibm_create_protection_group_run_request", "run_type"),
				Description: "Protection group id",
			},
			"run_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_protection_group_run_request", "run_type"),
				Description: "Type of protection run. 'kRegular' indicates an incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a kRegular schedule captures all the blocks. 'kFull' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized. 'kLog' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time. 'kSystem' indicates system volume backup. It produces an image for bare metal recovery.",
			},
			"objects": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the list of objects to be protected by this Protection Group run. These can be leaf objects or non-leaf objects in the protection hierarchy. This must be specified only if a subset of objects from the Protection Groups needs to be protected.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the id of object.",
						},
						"app_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of ids of applications.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"physical_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies physical parameters for this run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metadata_file_path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies metadata file path during run-now requests for physical file based backups for some specific source. If specified, it will override any default metadata/directive file path set at the object level for the source. Also note that if the job default does not specify a metadata/directive file path for the source, then specifying this field for that source during run-now request will be rejected.",
									},
								},
							},
						},
					},
				},
			},
			"targets_config": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Specifies the replication and archival targets.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"use_policy_defaults": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Specifies whether to use default policy settings or not. If specified as true then 'replications' and 'arcihvals' should not be specified. In case of true value, replicatioan targets congfigured in the policy will be added internally.",
						},
						"replications": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of replication targets configurations.",
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
						"archivals": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of archival targets configurations.",
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
						"cloud_replications": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of cloud replication targets configurations.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"aws_target": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the configuration for adding AWS as repilcation target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the AWS Replication target.",
												},
												"region": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies id of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
												},
												"region_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the source id of the AWS protection source registered on IBM cluster.",
												},
											},
										},
									},
									"azure_target": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the configuration for adding Azure as replication target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the Azure Replication target.",
												},
												"resource_group": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies id of the Azure resource group used to filter regions in UI.",
												},
												"resource_group_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the Azure resource group used to filter regions in UI.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the source id of the Azure protection source registered on IBM cluster.",
												},
												"storage_account": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies id of the storage account of Azure replication target which will contain storage container.",
												},
												"storage_account_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the storage account of Azure replication target which will contain storage container.",
												},
												"storage_container": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies id of the storage container of Azure Replication target.",
												},
												"storage_container_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the storage container of Azure Replication target.",
												},
												"storage_resource_group": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies id of the storage resource group of Azure Replication target.",
												},
												"storage_resource_group_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the storage resource group of Azure Replication target.",
												},
											},
										},
									},
									"target_type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the type of target to which replication need to be performed.",
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
					},
				},
			},
		},
	}
}

const (
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
)

func checkDiffResourceIbmBackupRecoveryProtectionGroupRun(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// skip if it's a new resource
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	// display a warning in the plan if resource is updated
	for fieldName := range ResourceIbmBackupRecoveryProtectionGroupRunRequest().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_protection_group_run_request cannot be updated.")
		}
	}
	return nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "run_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "kFull, kHydrateCDP, kLog, kRegular, kStorageArraySnapshot, kSystem",
		},
		validate.ValidateSchema{
			Identifier:                 "group_id",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_backup_recovery_protection_group_run_request", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmBackupRecoveryProtectionGroupRunRequestCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_group_run_request", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createProtectionGroupRunOptions := &backuprecoveryv1.CreateProtectionGroupRunOptions{}

	createProtectionGroupRunOptions.SetID(d.Get("group_id").(string))
	createProtectionGroupRunOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	createProtectionGroupRunOptions.SetRunType(d.Get("run_type").(string))
	if _, ok := d.GetOk("objects"); ok {
		var newObjects []backuprecoveryv1.RunObject
		for _, v := range d.Get("objects").([]interface{}) {
			value := v.(map[string]interface{})
			newObjectsItem, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunObject(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_group_run_request", "create", "parse-objects").GetDiag()
			}
			newObjects = append(newObjects, *newObjectsItem)
		}
		createProtectionGroupRunOptions.SetObjects(newObjects)
	}
	if _, ok := d.GetOk("targets_config"); ok {
		newTargetsConfigModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunTargetsConfiguration(d.Get("targets_config.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_group_run_request", "create", "parse-targets_config").GetDiag()
		}
		createProtectionGroupRunOptions.SetTargetsConfig(newTargetsConfigModel)
	}

	createProtectionGroupRunResponse, _, err := backupRecoveryClient.CreateProtectionGroupRunWithContext(context, createProtectionGroupRunOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateProtectionGroupRunWithContext failed: %s", err.Error()), "ibm_backup_recovery_protection_group_run_request", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*createProtectionGroupRunResponse.ProtectionGroupID)
	if err = d.Set("group_id", *createProtectionGroupRunResponse.ProtectionGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting group_id: %s", err))
	}
	return nil
}

func resourceIbmBackupRecoveryProtectionGroupRunRequestRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIbmBackupRecoveryProtectionGroupRunRequestDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmBackupRecoveryProtectionGroupRunRequestUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "update" operation.
	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource update will only affect terraform state and not the actual backend resource",
		Detail:   "Update operation for this resource is not supported and will only affect the terraform statefile. No changes will be made to the backend resource. Please use ibm_backup_recovery_update_protection_group_run_request resource for updates.",
	}
	// d.SetId("")
	diags = append(diags, warning)
	return diags
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunObject(modelMap map[string]interface{}) (*backuprecoveryv1.RunObject, error) {
	model := &backuprecoveryv1.RunObject{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["app_ids"] != nil {
		appIds := []int64{}
		for _, appIdsItem := range modelMap["app_ids"].([]interface{}) {
			appIds = append(appIds, int64(appIdsItem.(int)))
		}
		model.AppIds = appIds
	}
	if modelMap["physical_params"] != nil && len(modelMap["physical_params"].([]interface{})) > 0 {
		PhysicalParamsModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunObjectPhysicalParams(modelMap["physical_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PhysicalParams = PhysicalParamsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunObjectPhysicalParams(modelMap map[string]interface{}) (*backuprecoveryv1.RunObjectPhysicalParams, error) {
	model := &backuprecoveryv1.RunObjectPhysicalParams{}
	if modelMap["metadata_file_path"] != nil && modelMap["metadata_file_path"].(string) != "" {
		model.MetadataFilePath = core.StringPtr(modelMap["metadata_file_path"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunTargetsConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.RunTargetsConfiguration, error) {
	model := &backuprecoveryv1.RunTargetsConfiguration{}
	if modelMap["use_policy_defaults"] != nil {
		model.UsePolicyDefaults = core.BoolPtr(modelMap["use_policy_defaults"].(bool))
	}
	if modelMap["replications"] != nil {
		replications := []backuprecoveryv1.RunReplicationConfig{}
		for _, replicationsItem := range modelMap["replications"].([]interface{}) {
			replicationsItemModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunReplicationConfig(replicationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			replications = append(replications, *replicationsItemModel)
		}
		model.Replications = replications
	}
	if modelMap["archivals"] != nil {
		archivals := []backuprecoveryv1.RunArchivalConfig{}
		for _, archivalsItem := range modelMap["archivals"].([]interface{}) {
			archivalsItemModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunArchivalConfig(archivalsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			archivals = append(archivals, *archivalsItemModel)
		}
		model.Archivals = archivals
	}
	if modelMap["cloud_replications"] != nil {
		cloudReplications := []backuprecoveryv1.RunCloudReplicationConfig{}
		for _, cloudReplicationsItem := range modelMap["cloud_replications"].([]interface{}) {
			cloudReplicationsItemModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunCloudReplicationConfig(cloudReplicationsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			cloudReplications = append(cloudReplications, *cloudReplicationsItemModel)
		}
		model.CloudReplications = cloudReplications
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunReplicationConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RunReplicationConfig, error) {
	model := &backuprecoveryv1.RunReplicationConfig{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["retention"] != nil && len(modelMap["retention"].([]interface{})) > 0 {
		RetentionModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Retention = RetentionModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRetention(modelMap map[string]interface{}) (*backuprecoveryv1.Retention, error) {
	model := &backuprecoveryv1.Retention{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["data_lock_config"] != nil && len(modelMap["data_lock_config"].([]interface{})) > 0 {
		DataLockConfigModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToDataLockConfig(modelMap["data_lock_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DataLockConfig = DataLockConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToDataLockConfig(modelMap map[string]interface{}) (*backuprecoveryv1.DataLockConfig, error) {
	model := &backuprecoveryv1.DataLockConfig{}
	model.Mode = core.StringPtr(modelMap["mode"].(string))
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["enable_worm_on_external_target"] != nil {
		model.EnableWormOnExternalTarget = core.BoolPtr(modelMap["enable_worm_on_external_target"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunArchivalConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RunArchivalConfig, error) {
	model := &backuprecoveryv1.RunArchivalConfig{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	model.ArchivalTargetType = core.StringPtr(modelMap["archival_target_type"].(string))
	if modelMap["retention"] != nil && len(modelMap["retention"].([]interface{})) > 0 {
		RetentionModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
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

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRunCloudReplicationConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RunCloudReplicationConfig, error) {
	model := &backuprecoveryv1.RunCloudReplicationConfig{}
	if modelMap["aws_target"] != nil && len(modelMap["aws_target"].([]interface{})) > 0 {
		AwsTargetModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToAWSTargetConfig(modelMap["aws_target"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AwsTarget = AwsTargetModel
	}
	if modelMap["azure_target"] != nil && len(modelMap["azure_target"].([]interface{})) > 0 {
		AzureTargetModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToAzureTargetConfig(modelMap["azure_target"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AzureTarget = AzureTargetModel
	}
	model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	if modelMap["retention"] != nil && len(modelMap["retention"].([]interface{})) > 0 {
		RetentionModel, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Retention = RetentionModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToAWSTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTargetConfig, error) {
	model := &backuprecoveryv1.AWSTargetConfig{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Region = core.Int64Ptr(int64(modelMap["region"].(int)))
	if modelMap["region_name"] != nil && modelMap["region_name"].(string) != "" {
		model.RegionName = core.StringPtr(modelMap["region_name"].(string))
	}
	model.SourceID = core.Int64Ptr(int64(modelMap["source_id"].(int)))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestMapToAzureTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTargetConfig, error) {
	model := &backuprecoveryv1.AzureTargetConfig{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["resource_group"] != nil {
		model.ResourceGroup = core.Int64Ptr(int64(modelMap["resource_group"].(int)))
	}
	if modelMap["resource_group_name"] != nil && modelMap["resource_group_name"].(string) != "" {
		model.ResourceGroupName = core.StringPtr(modelMap["resource_group_name"].(string))
	}
	model.SourceID = core.Int64Ptr(int64(modelMap["source_id"].(int)))
	if modelMap["storage_account"] != nil {
		model.StorageAccount = core.Int64Ptr(int64(modelMap["storage_account"].(int)))
	}
	if modelMap["storage_account_name"] != nil && modelMap["storage_account_name"].(string) != "" {
		model.StorageAccountName = core.StringPtr(modelMap["storage_account_name"].(string))
	}
	if modelMap["storage_container"] != nil {
		model.StorageContainer = core.Int64Ptr(int64(modelMap["storage_container"].(int)))
	}
	if modelMap["storage_container_name"] != nil && modelMap["storage_container_name"].(string) != "" {
		model.StorageContainerName = core.StringPtr(modelMap["storage_container_name"].(string))
	}
	if modelMap["storage_resource_group"] != nil {
		model.StorageResourceGroup = core.Int64Ptr(int64(modelMap["storage_resource_group"].(int)))
	}
	if modelMap["storage_resource_group_name"] != nil && modelMap["storage_resource_group_name"].(string) != "" {
		model.StorageResourceGroupName = core.StringPtr(modelMap["storage_resource_group_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRunObjectToMap(model *backuprecoveryv1.RunObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.AppIds != nil {
		modelMap["app_ids"] = model.AppIds
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRunObjectPhysicalParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRunObjectPhysicalParamsToMap(model *backuprecoveryv1.RunObjectPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MetadataFilePath != nil {
		modelMap["metadata_file_path"] = *model.MetadataFilePath
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRunTargetsConfigurationToMap(model *backuprecoveryv1.RunTargetsConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UsePolicyDefaults != nil {
		modelMap["use_policy_defaults"] = *model.UsePolicyDefaults
	}
	if model.Replications != nil {
		replications := []map[string]interface{}{}
		for _, replicationsItem := range model.Replications {
			replicationsItemMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRunReplicationConfigToMap(&replicationsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			replications = append(replications, replicationsItemMap)
		}
		modelMap["replications"] = replications
	}
	if model.Archivals != nil {
		archivals := []map[string]interface{}{}
		for _, archivalsItem := range model.Archivals {
			archivalsItemMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRunArchivalConfigToMap(&archivalsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			archivals = append(archivals, archivalsItemMap)
		}
		modelMap["archivals"] = archivals
	}
	if model.CloudReplications != nil {
		cloudReplications := []map[string]interface{}{}
		for _, cloudReplicationsItem := range model.CloudReplications {
			cloudReplicationsItemMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRunCloudReplicationConfigToMap(&cloudReplicationsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cloudReplications = append(cloudReplications, cloudReplicationsItemMap)
		}
		modelMap["cloud_replications"] = cloudReplications
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRunReplicationConfigToMap(model *backuprecoveryv1.RunReplicationConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Retention != nil {
		retentionMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRetentionToMap(model.Retention)
		if err != nil {
			return modelMap, err
		}
		modelMap["retention"] = []map[string]interface{}{retentionMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRetentionToMap(model *backuprecoveryv1.Retention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestDataLockConfigToMap(model *backuprecoveryv1.DataLockConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = *model.Mode
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.EnableWormOnExternalTarget != nil {
		modelMap["enable_worm_on_external_target"] = *model.EnableWormOnExternalTarget
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRunArchivalConfigToMap(model *backuprecoveryv1.RunArchivalConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	modelMap["archival_target_type"] = *model.ArchivalTargetType
	if model.Retention != nil {
		retentionMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRetentionToMap(model.Retention)
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

func ResourceIbmBackupRecoveryProtectionGroupRunRequestRunCloudReplicationConfigToMap(model *backuprecoveryv1.RunCloudReplicationConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTarget != nil {
		awsTargetMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestAWSTargetConfigToMap(model.AwsTarget)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_target"] = []map[string]interface{}{awsTargetMap}
	}
	if model.AzureTarget != nil {
		azureTargetMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestAzureTargetConfigToMap(model.AzureTarget)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_target"] = []map[string]interface{}{azureTargetMap}
	}
	modelMap["target_type"] = *model.TargetType
	if model.Retention != nil {
		retentionMap, err := ResourceIbmBackupRecoveryProtectionGroupRunRequestRetentionToMap(model.Retention)
		if err != nil {
			return modelMap, err
		}
		modelMap["retention"] = []map[string]interface{}{retentionMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestAWSTargetConfigToMap(model *backuprecoveryv1.AWSTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	modelMap["region"] = flex.IntValue(model.Region)
	if model.RegionName != nil {
		modelMap["region_name"] = *model.RegionName
	}
	modelMap["source_id"] = flex.IntValue(model.SourceID)
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionGroupRunRequestAzureTargetConfigToMap(model *backuprecoveryv1.AzureTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceGroup != nil {
		modelMap["resource_group"] = flex.IntValue(model.ResourceGroup)
	}
	if model.ResourceGroupName != nil {
		modelMap["resource_group_name"] = *model.ResourceGroupName
	}
	modelMap["source_id"] = flex.IntValue(model.SourceID)
	if model.StorageAccount != nil {
		modelMap["storage_account"] = flex.IntValue(model.StorageAccount)
	}
	if model.StorageAccountName != nil {
		modelMap["storage_account_name"] = *model.StorageAccountName
	}
	if model.StorageContainer != nil {
		modelMap["storage_container"] = flex.IntValue(model.StorageContainer)
	}
	if model.StorageContainerName != nil {
		modelMap["storage_container_name"] = *model.StorageContainerName
	}
	if model.StorageResourceGroup != nil {
		modelMap["storage_resource_group"] = flex.IntValue(model.StorageResourceGroup)
	}
	if model.StorageResourceGroupName != nil {
		modelMap["storage_resource_group_name"] = *model.StorageResourceGroupName
	}
	return modelMap, nil
}
