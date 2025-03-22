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

func DataSourceIbmBackupRecoveryObjectSnapshots() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryObjectSnapshotsRead,

		Schema: map[string]*schema.Schema{
			"object_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Specifies the id of the Object.",
			},
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"from_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the timestamp in Unix time epoch in microseconds to filter Object's snapshots which were taken after this value.",
			},
			"to_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the timestamp in Unix time epoch in microseconds to filter Object's snapshots which were taken before this value.",
			},
			"run_start_from_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the timestamp in Unix time epoch in microseconds to filter Object's snapshots which were run after this value.",
			},
			"run_start_to_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the timestamp in Unix time epoch in microseconds to filter Object's snapshots which were run before this value.",
			},
			"snapshot_actions": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of recovery actions. Only snapshots that apply to these actions will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by run type. Only protection runs matching the specified types will be returned. By default, CDP hydration snapshots are not included unless explicitly queried using this field.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"protection_group_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "If specified, this returns only the snapshots of the specified object ID, which belong to the provided protection group IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"run_instance_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by a list of run instance IDs. If specified, only snapshots created by these protection runs will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"region_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by a list of region IDs.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"object_action_keys": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by ObjectActionKey, which uniquely represents the protection of an object. An object can be protected in multiple ways but at most once for a given combination of ObjectActionKey. When specified, only snapshots matching the given action keys are returned for the corresponding object.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"snapshots": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of snapshots.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aws_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies parameters of AWS type snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection type of AWS snapshots.",
									},
								},
							},
						},
						"azure_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies parameters of Azure type snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection type of Azure snapshots.",
									},
								},
							},
						},
						"cluster_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the cluster id where this snapshot belongs to.",
						},
						"cluster_incarnation_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the cluster incarnation id where this snapshot belongs to.",
						},
						"elastifile_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the common parameters for NAS objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_nas_mount_protocols": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of NAS mount protocols supported by this object.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"environment": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the snapshot environment.",
						},
						"expiry_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the expiry time of the snapshot in Unix timestamp epoch in microseconds. If the snapshot has no expiry, this property will not be set.",
						},
						"external_target_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies archival target summary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the archival target ID.",
									},
									"archival_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the archival task id. This is a protection group UID which only applies when archival type is 'Tape'.",
									},
									"target_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the archival target name.",
									},
									"target_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the archival target type.",
									},
									"usage_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the usage type for the target.",
									},
									"ownership_context": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the ownership context for the target.",
									},
									"tier_settings": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the tier info for archival.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aws_tiering": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies aws tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the time period after which the backup will be moved from current tier to next tier.",
																		},
																		"tier_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the AWS tier types.",
																		},
																	},
																},
															},
														},
													},
												},
												"azure_tiering": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies Azure tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the time period after which the backup will be moved from current tier to next tier.",
																		},
																		"tier_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the Azure tier types.",
																		},
																	},
																},
															},
														},
													},
												},
												"cloud_platform": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the cloud platform to enable tiering.",
												},
												"google_tiering": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies Google tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the time period after which the backup will be moved from current tier to next tier.",
																		},
																		"tier_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the Google tier types.",
																		},
																	},
																},
															},
														},
													},
												},
												"oracle_tiering": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies Oracle tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the time period after which the backup will be moved from current tier to next tier.",
																		},
																		"tier_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the Oracle tier types.",
																		},
																	},
																},
															},
														},
													},
												},
												"current_tier_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of the current tier where the snapshot resides. This will be specified if the run is a CAD run.",
												},
											},
										},
									},
								},
							},
						},
						"flashblade_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the common parameters for Flashblade objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_nas_mount_protocols": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of NAS mount protocols supported by this object.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"generic_nas_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the common parameters for NAS objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_nas_mount_protocols": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of NAS mount protocols supported by this object.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"gpfs_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the common parameters for NAS objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_nas_mount_protocols": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of NAS mount protocols supported by this object.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"has_data_lock": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if this snapshot has datalock.",
						},
						"hyperv_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies parameters of HyperV type snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection type of HyperV snapshots.",
									},
								},
							},
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the id of the snapshot.",
						},
						"indexing_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the indexing status of objects in this snapshot.<br> 'InProgress' indicates the indexing is in progress.<br> 'Done' indicates indexing is done.<br> 'NoIndex' indicates indexing is not applicable.<br> 'Error' indicates indexing failed with error.",
						},
						"isilon_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the common parameters for Isilon objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_nas_mount_protocols": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of NAS mount protocols supported by this object.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"netapp_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the common parameters for Netapp objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"supported_nas_mount_protocols": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of NAS mount protocols supported by this object.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"volume_extended_style": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the extended style of a NetApp volume.",
									},
									"volume_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the Netapp volume type.",
									},
								},
							},
						},
						"object_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the object id which the snapshot is taken from.",
						},
						"object_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the object name which the snapshot is taken from.",
						},
						"on_legal_hold": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if this snapshot is on legalhold.",
						},
						"ownership_context": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the ownership context for the target.",
						},
						"physical_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies parameters of Physical type snapshots.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_system_backup": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies if system backup was enabled for the source in that particular run.",
									},
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection type of Physical snapshots.",
									},
								},
							},
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies id of the Protection Group.",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies name of the Protection Group.",
						},
						"protection_group_run_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies id of the Protection Group Run.",
						},
						"region_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the region id where this snapshot belongs to.",
						},
						"run_instance_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the instance id of the protection run which create the snapshot.",
						},
						"run_start_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the start time of the run in micro seconds.",
						},
						"run_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the type of protection run created this snapshot.",
						},
						"sfdc_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the Salesforce objects mutation parameters.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"records_added": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of records added for the Object.",
									},
									"records_modified": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of records updated for the Object.",
									},
									"records_removed": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of records removed from the Object.",
									},
								},
							},
						},
						"snapshot_target_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the target type where the Object's snapshot resides.",
						},
						"snapshot_timestamp_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the timestamp in Unix time epoch in microseconds when the snapshot is taken for the specified Object.",
						},
						"source_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the source protection group id in case of replication.",
						},
						"source_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the object source id which the snapshot is taken from.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the Storage Domain id where the snapshot of object is present.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmBackupRecoveryObjectSnapshotsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_object_snapshots", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getObjectSnapshotsOptions := &backuprecoveryv1.GetObjectSnapshotsOptions{}

	getObjectSnapshotsOptions.SetID(int64(d.Get("object_id").(int)))
	getObjectSnapshotsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("from_time_usecs"); ok {
		getObjectSnapshotsOptions.SetFromTimeUsecs(int64(d.Get("from_time_usecs").(int)))
	}
	if _, ok := d.GetOk("to_time_usecs"); ok {
		getObjectSnapshotsOptions.SetToTimeUsecs(int64(d.Get("to_time_usecs").(int)))
	}
	if _, ok := d.GetOk("run_start_from_time_usecs"); ok {
		getObjectSnapshotsOptions.SetRunStartFromTimeUsecs(int64(d.Get("run_start_from_time_usecs").(int)))
	}
	if _, ok := d.GetOk("run_start_to_time_usecs"); ok {
		getObjectSnapshotsOptions.SetRunStartToTimeUsecs(int64(d.Get("run_start_to_time_usecs").(int)))
	}
	if _, ok := d.GetOk("snapshot_actions"); ok {
		var snapshotActions []string
		for _, v := range d.Get("snapshot_actions").([]interface{}) {
			snapshotActionsItem := v.(string)
			snapshotActions = append(snapshotActions, snapshotActionsItem)
		}
		getObjectSnapshotsOptions.SetSnapshotActions(snapshotActions)
	}
	if _, ok := d.GetOk("run_types"); ok {
		var runTypes []string
		for _, v := range d.Get("run_types").([]interface{}) {
			runTypesItem := v.(string)
			runTypes = append(runTypes, runTypesItem)
		}
		getObjectSnapshotsOptions.SetRunTypes(runTypes)
	}
	if _, ok := d.GetOk("protection_group_ids"); ok {
		var protectionGroupIds []string
		for _, v := range d.Get("protection_group_ids").([]interface{}) {
			protectionGroupIdsItem := v.(string)
			protectionGroupIds = append(protectionGroupIds, protectionGroupIdsItem)
		}
		getObjectSnapshotsOptions.SetProtectionGroupIds(protectionGroupIds)
	}
	if _, ok := d.GetOk("run_instance_ids"); ok {
		var runInstanceIds []int64
		for _, v := range d.Get("run_instance_ids").([]interface{}) {
			runInstanceIdsItem := int64(v.(int))
			runInstanceIds = append(runInstanceIds, runInstanceIdsItem)
		}
		getObjectSnapshotsOptions.SetRunInstanceIds(runInstanceIds)
	}
	if _, ok := d.GetOk("region_ids"); ok {
		var regionIds []string
		for _, v := range d.Get("region_ids").([]interface{}) {
			regionIdsItem := v.(string)
			regionIds = append(regionIds, regionIdsItem)
		}
		getObjectSnapshotsOptions.SetRegionIds(regionIds)
	}
	if _, ok := d.GetOk("object_action_keys"); ok {
		var objectActionKeys []string
		for _, v := range d.Get("object_action_keys").([]interface{}) {
			objectActionKeysItem := v.(string)
			objectActionKeys = append(objectActionKeys, objectActionKeysItem)
		}
		getObjectSnapshotsOptions.SetObjectActionKeys(objectActionKeys)
	}

	getObjectSnapshotsResponse, _, err := backupRecoveryClient.GetObjectSnapshotsWithContext(context, getObjectSnapshotsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetObjectSnapshotsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_object_snapshots", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryObjectSnapshotsID(d))

	if !core.IsNil(getObjectSnapshotsResponse.Snapshots) {
		snapshots := []map[string]interface{}{}
		for _, snapshotsItem := range getObjectSnapshotsResponse.Snapshots {
			snapshotsItemMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsObjectSnapshotToMap(&snapshotsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_object_snapshots", "read", "snapshots-to-map").GetDiag()
			}
			snapshots = append(snapshots, snapshotsItemMap)
		}
		if err = d.Set("snapshots", snapshots); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting snapshots: %s", err), "(Data) ibm_backup_recovery_object_snapshots", "read", "set-snapshots").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryObjectSnapshotsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryObjectSnapshotsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryObjectSnapshotsObjectSnapshotToMap(model *backuprecoveryv1.ObjectSnapshot) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsParams != nil {
		awsParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsAwsSnapshotParamsToMap(model.AwsParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_params"] = []map[string]interface{}{awsParamsMap}
	}
	if model.AzureParams != nil {
		azureParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsAzureSnapshotParamsToMap(model.AzureParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_params"] = []map[string]interface{}{azureParamsMap}
	}
	if model.ClusterID != nil {
		modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	}
	if model.ClusterIncarnationID != nil {
		modelMap["cluster_incarnation_id"] = flex.IntValue(model.ClusterIncarnationID)
	}
	if model.ElastifileParams != nil {
		elastifileParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsCommonNasObjectParamsToMap(model.ElastifileParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["elastifile_params"] = []map[string]interface{}{elastifileParamsMap}
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.ExternalTargetInfo != nil {
		externalTargetInfoMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsArchivalTargetSummaryInfoToMap(model.ExternalTargetInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["external_target_info"] = []map[string]interface{}{externalTargetInfoMap}
	}
	if model.FlashbladeParams != nil {
		flashbladeParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsFlashbladeObjectParamsToMap(model.FlashbladeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["flashblade_params"] = []map[string]interface{}{flashbladeParamsMap}
	}
	if model.GenericNasParams != nil {
		genericNasParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsCommonNasObjectParamsToMap(model.GenericNasParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["generic_nas_params"] = []map[string]interface{}{genericNasParamsMap}
	}
	if model.GpfsParams != nil {
		gpfsParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsCommonNasObjectParamsToMap(model.GpfsParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["gpfs_params"] = []map[string]interface{}{gpfsParamsMap}
	}
	if model.HasDataLock != nil {
		modelMap["has_data_lock"] = *model.HasDataLock
	}
	if model.HypervParams != nil {
		hypervParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsHypervSnapshotParamsToMap(model.HypervParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["hyperv_params"] = []map[string]interface{}{hypervParamsMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.IndexingStatus != nil {
		modelMap["indexing_status"] = *model.IndexingStatus
	}
	if model.IsilonParams != nil {
		isilonParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsIsilonObjectParamsToMap(model.IsilonParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["isilon_params"] = []map[string]interface{}{isilonParamsMap}
	}
	if model.NetappParams != nil {
		netappParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsNetappObjectParamsToMap(model.NetappParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["netapp_params"] = []map[string]interface{}{netappParamsMap}
	}
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	if model.ObjectName != nil {
		modelMap["object_name"] = *model.ObjectName
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	if model.OwnershipContext != nil {
		modelMap["ownership_context"] = *model.OwnershipContext
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsPhysicalSnapshotParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.ProtectionGroupRunID != nil {
		modelMap["protection_group_run_id"] = *model.ProtectionGroupRunID
	}
	if model.RegionID != nil {
		modelMap["region_id"] = *model.RegionID
	}
	if model.RunInstanceID != nil {
		modelMap["run_instance_id"] = flex.IntValue(model.RunInstanceID)
	}
	if model.RunStartTimeUsecs != nil {
		modelMap["run_start_time_usecs"] = flex.IntValue(model.RunStartTimeUsecs)
	}
	if model.RunType != nil {
		modelMap["run_type"] = *model.RunType
	}
	if model.SfdcParams != nil {
		sfdcParamsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsSfdcObjectParamsToMap(model.SfdcParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["sfdc_params"] = []map[string]interface{}{sfdcParamsMap}
	}
	if model.SnapshotTargetType != nil {
		modelMap["snapshot_target_type"] = *model.SnapshotTargetType
	}
	if model.SnapshotTimestampUsecs != nil {
		modelMap["snapshot_timestamp_usecs"] = flex.IntValue(model.SnapshotTimestampUsecs)
	}
	if model.SourceGroupID != nil {
		modelMap["source_group_id"] = *model.SourceGroupID
	}
	if model.SourceID != nil {
		modelMap["source_id"] = flex.IntValue(model.SourceID)
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsAwsSnapshotParamsToMap(model *backuprecoveryv1.AwsSnapshotParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ProtectionType != nil {
		modelMap["protection_type"] = *model.ProtectionType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsAzureSnapshotParamsToMap(model *backuprecoveryv1.AzureSnapshotParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ProtectionType != nil {
		modelMap["protection_type"] = *model.ProtectionType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsCommonNasObjectParamsToMap(model *backuprecoveryv1.CommonNasObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SupportedNasMountProtocols != nil {
		modelMap["supported_nas_mount_protocols"] = model.SupportedNasMountProtocols
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsArchivalTargetSummaryInfoToMap(model *backuprecoveryv1.ArchivalTargetSummaryInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetID != nil {
		modelMap["target_id"] = flex.IntValue(model.TargetID)
	}
	if model.ArchivalTaskID != nil {
		modelMap["archival_task_id"] = *model.ArchivalTaskID
	}
	if model.TargetName != nil {
		modelMap["target_name"] = *model.TargetName
	}
	if model.TargetType != nil {
		modelMap["target_type"] = *model.TargetType
	}
	if model.UsageType != nil {
		modelMap["usage_type"] = *model.UsageType
	}
	if model.OwnershipContext != nil {
		modelMap["ownership_context"] = *model.OwnershipContext
	}
	if model.TierSettings != nil {
		tierSettingsMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsArchivalTargetTierInfoToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsArchivalTargetTierInfoToMap(model *backuprecoveryv1.ArchivalTargetTierInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}
	if model.GoogleTiering != nil {
		googleTieringMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsOracleTiersToMap(model.OracleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["oracle_tiering"] = []map[string]interface{}{oracleTieringMap}
	}
	if model.CurrentTierType != nil {
		modelMap["current_tier_type"] = *model.CurrentTierType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoveAfterUnit != nil {
		modelMap["move_after_unit"] = *model.MoveAfterUnit
	}
	if model.MoveAfter != nil {
		modelMap["move_after"] = flex.IntValue(model.MoveAfter)
	}
	modelMap["tier_type"] = *model.TierType
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoveAfterUnit != nil {
		modelMap["move_after_unit"] = *model.MoveAfterUnit
	}
	if model.MoveAfter != nil {
		modelMap["move_after"] = flex.IntValue(model.MoveAfter)
	}
	modelMap["tier_type"] = *model.TierType
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoveAfterUnit != nil {
		modelMap["move_after_unit"] = *model.MoveAfterUnit
	}
	if model.MoveAfter != nil {
		modelMap["move_after"] = flex.IntValue(model.MoveAfter)
	}
	modelMap["tier_type"] = *model.TierType
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryObjectSnapshotsOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MoveAfterUnit != nil {
		modelMap["move_after_unit"] = *model.MoveAfterUnit
	}
	if model.MoveAfter != nil {
		modelMap["move_after"] = flex.IntValue(model.MoveAfter)
	}
	modelMap["tier_type"] = *model.TierType
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsFlashbladeObjectParamsToMap(model *backuprecoveryv1.FlashbladeObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SupportedNasMountProtocols != nil {
		modelMap["supported_nas_mount_protocols"] = model.SupportedNasMountProtocols
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsHypervSnapshotParamsToMap(model *backuprecoveryv1.HypervSnapshotParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ProtectionType != nil {
		modelMap["protection_type"] = *model.ProtectionType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsIsilonObjectParamsToMap(model *backuprecoveryv1.IsilonObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SupportedNasMountProtocols != nil {
		modelMap["supported_nas_mount_protocols"] = model.SupportedNasMountProtocols
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsNetappObjectParamsToMap(model *backuprecoveryv1.NetappObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SupportedNasMountProtocols != nil {
		modelMap["supported_nas_mount_protocols"] = model.SupportedNasMountProtocols
	}
	if model.VolumeExtendedStyle != nil {
		modelMap["volume_extended_style"] = *model.VolumeExtendedStyle
	}
	if model.VolumeType != nil {
		modelMap["volume_type"] = *model.VolumeType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsPhysicalSnapshotParamsToMap(model *backuprecoveryv1.PhysicalSnapshotParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	if model.ProtectionType != nil {
		modelMap["protection_type"] = *model.ProtectionType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSnapshotsSfdcObjectParamsToMap(model *backuprecoveryv1.SfdcObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RecordsAdded != nil {
		modelMap["records_added"] = flex.IntValue(model.RecordsAdded)
	}
	if model.RecordsModified != nil {
		modelMap["records_modified"] = flex.IntValue(model.RecordsModified)
	}
	if model.RecordsRemoved != nil {
		modelMap["records_removed"] = flex.IntValue(model.RecordsRemoved)
	}
	return modelMap, nil
}
