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

func DataSourceIbmBackupRecoveryProtectionGroupRuns() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryProtectionGroupRunsRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"protection_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies a unique id of the Protection Group.",
			},
			"request_initiator_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of request from UI, which is used for services like magneto to determine the priority of requests.",
			},
			"run_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the protection run id.",
			},
			"start_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Start time for time range filter. Specify the start time as a Unix epoch Timestamp (in microseconds), only runs executing after this time will be returned. By default it is endTimeUsecs minus an hour.",
			},
			"end_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "End time for time range filter. Specify the end time as a Unix epoch Timestamp (in microseconds), only runs executing before this time will be returned. By default it is current time.",
			},
			"run_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by run type. Only protection run matching the specified types will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_object_details": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies if the result includes the object details for each protection run. If set to true, details of the protected object will be returned. If set to false or not specified, details will not be returned.",
			},
			"local_backup_run_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of local backup status, runs matching the status will be returned.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"replication_run_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of replication status, runs matching the status will be returned.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"archival_run_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of archival status, runs matching the status will be returned.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cloud_spin_run_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of cloud spin status, runs matching the status will be returned.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"num_runs": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the max number of runs. If not specified, at most 100 runs will be returned.",
			},
			"exclude_non_restorable_runs": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Specifies whether to exclude non restorable runs. Run is treated restorable only if there is atleast one object snapshot (which may be either a local or an archival snapshot) which is not deleted or expired. Default value is false.",
			},
			"run_tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of tags for protection runs. If this is specified, only the runs which match these tags will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"use_cached_data": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether we can serve the GET request from the read replica cache. There is a lag of 15 seconds between the read replica and primary data source.",
			},
			"filter_by_end_time": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the runs with backup end time within the specified time range will be returned. Otherwise, the runs with start time in the time range are returned.",
			},
			"snapshot_target_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the snapshot's target type which should be filtered.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"only_return_successful_copy_run": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "only successful copyruns are returned.",
			},
			"filter_by_copy_task_end_time": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, then the details of the runs for which any copyTask completed in the given timerange will be returned. Only one of filterByEndTime and filterByCopyTaskEndTime can be set.",
			},
			"runs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of Protection Group runs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the ID of the Protection Group run.",
						},
						"protection_group_instance_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Protection Group instance Id. This field will be removed later.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ProtectionGroupId to which this run belongs.",
						},
						"is_replication_run": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if this protection run is a replication run.",
						},
						"origin_cluster_identifier": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the information about a cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the id of the cluster.",
									},
									"cluster_incarnation_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the incarnation id of the cluster.",
									},
									"cluster_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the name of the cluster.",
									},
								},
							},
						},
						"origin_protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ProtectionGroupId to which this run belongs on the primary cluster if this run is a replication run.",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the Protection Group to which this run belongs.",
						},
						"is_local_snapshots_deleted": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if snapshots for this run has been deleted.",
						},
						"objects": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Snapahot, replication, archival results for each object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the Object Summary.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies object id.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the object.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies registered source id to which object belongs.",
												},
												"source_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies registered source name to which object belongs.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment of the object.",
												},
												"object_hash": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the hash identifier of the object.",
												},
												"object_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of the object.",
												},
												"logical_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the logical size of object in bytes.",
												},
												"uuid": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the uuid which is a unique identifier of the object.",
												},
												"global_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the global id which is a unique identifier of the object.",
												},
												"protection_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the protection type of the object if any.",
												},
												"sharepoint_site_summary": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the common parameters for Sharepoint site objects.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"site_web_url": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the web url for the Sharepoint site.",
															},
														},
													},
												},
												"os_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the operating system type of the object.",
												},
												"child_objects": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies child object details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{},
													},
												},
												"v_center_summary": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_cloud_env": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
															},
														},
													},
												},
												"windows_cluster_summary": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_source_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of cluster resource this source represents.",
															},
														},
													},
												},
											},
										},
									},
									"local_snapshot_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information about backup run for an object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"snapshot_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Snapshot info for an object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"snapshot_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Snapshot id for a successful snapshot. This field will not be set if the Protection Group Run has no successful attempt.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of snapshot.",
															},
															"status_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "A message decribing the status. This will be populated currently only for kWaitingForOlderBackupRun status.",
															},
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"admitted_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the backup task was admitted to run in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"permit_grant_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated to the time when permit is granted again.",
															},
															"queue_duration_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration between the startTime and when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated considering the time when permit is granted again. Queue duration = PermitGrantTimeUsecs - StartTimeUsecs.",
															},
															"snapshot_creation_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the source snapshot was taken in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about local snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical size of object(s) in bytes.",
																		},
																		"bytes_written": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total size of data in bytes written after taking backup.",
																		},
																		"bytes_read": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical bytes read for creating the snapshot.",
																		},
																	},
																},
															},
															"progress_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for backup of the object.",
															},
															"indexing_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for the indexing of documents in an object.",
															},
															"stats_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Stats task for an object.",
															},
															"warnings": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a list of warning messages.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"is_manually_deleted": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the snapshot is deleted manually.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds) for an object.",
															},
															"total_file_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The total number of file and directory entities visited in this backup. Only applicable to file based backups.",
															},
															"backup_file_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The total number of file and directory entities that are backed up in this run. Only applicable to file based backups.",
															},
															"data_lock_constraints": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the dataLock constraints for local or target snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"expiry_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
																		},
																	},
																},
															},
														},
													},
												},
												"failed_attempts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Failed backup attempts for an object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"admitted_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the backup task was admitted to run in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"permit_grant_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated to the time when permit is granted again.",
															},
															"queue_duration_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration between the startTime and when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated considering the time when permit is granted again. Queue duration = PermitGrantTimeUsecs - StartTimeUsecs.",
															},
															"snapshot_creation_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the source snapshot was taken in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of the attempt for an object. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Pausing' indicates that the ongoing run is in the process of being paused. 'Resuming' indicates that the already paused run is in the process of being running again. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about local snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical size of object(s) in bytes.",
																		},
																		"bytes_written": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total size of data in bytes written after taking backup.",
																		},
																		"bytes_read": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical bytes read for creating the snapshot.",
																		},
																	},
																},
															},
															"progress_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for an object.",
															},
															"message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "A message about the error if encountered while performing backup.",
															},
														},
													},
												},
											},
										},
									},
									"original_backup_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information about backup run for an object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"snapshot_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Snapshot info for an object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"snapshot_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Snapshot id for a successful snapshot. This field will not be set if the Protection Group Run has no successful attempt.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of snapshot.",
															},
															"status_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "A message decribing the status. This will be populated currently only for kWaitingForOlderBackupRun status.",
															},
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"admitted_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the backup task was admitted to run in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"permit_grant_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated to the time when permit is granted again.",
															},
															"queue_duration_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration between the startTime and when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated considering the time when permit is granted again. Queue duration = PermitGrantTimeUsecs - StartTimeUsecs.",
															},
															"snapshot_creation_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the source snapshot was taken in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about local snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical size of object(s) in bytes.",
																		},
																		"bytes_written": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total size of data in bytes written after taking backup.",
																		},
																		"bytes_read": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical bytes read for creating the snapshot.",
																		},
																	},
																},
															},
															"progress_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for backup of the object.",
															},
															"indexing_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for the indexing of documents in an object.",
															},
															"stats_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Stats task for an object.",
															},
															"warnings": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a list of warning messages.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"is_manually_deleted": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the snapshot is deleted manually.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds) for an object.",
															},
															"total_file_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The total number of file and directory entities visited in this backup. Only applicable to file based backups.",
															},
															"backup_file_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "The total number of file and directory entities that are backed up in this run. Only applicable to file based backups.",
															},
															"data_lock_constraints": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the dataLock constraints for local or target snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"expiry_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
																		},
																	},
																},
															},
														},
													},
												},
												"failed_attempts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Failed backup attempts for an object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of attempt in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"admitted_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the backup task was admitted to run in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"permit_grant_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated to the time when permit is granted again.",
															},
															"queue_duration_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration between the startTime and when gatekeeper permit is granted to the backup task. If the backup task is rescheduled due to errors, the field is updated considering the time when permit is granted again. Queue duration = PermitGrantTimeUsecs - StartTimeUsecs.",
															},
															"snapshot_creation_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the source snapshot was taken in Unix epoch Timestamp(in microseconds) for an object.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of the attempt for an object. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Pausing' indicates that the ongoing run is in the process of being paused. 'Resuming' indicates that the already paused run is in the process of being running again. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about local snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical size of object(s) in bytes.",
																		},
																		"bytes_written": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total size of data in bytes written after taking backup.",
																		},
																		"bytes_read": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical bytes read for creating the snapshot.",
																		},
																	},
																},
															},
															"progress_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for an object.",
															},
															"message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "A message about the error if encountered while performing backup.",
															},
														},
													},
												},
											},
										},
									},
									"replication_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information about replication run for an object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"replication_target_results": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Replication result for a target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the id of the cluster.",
															},
															"cluster_incarnation_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the incarnation id of the cluster.",
															},
															"cluster_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the cluster.",
															},
															"aws_target_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
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
																			Computed:    true,
																			Description: "Specifies id of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
																		},
																		"region_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies name of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
																		},
																		"source_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the source id of the AWS protection source registered on IBM cluster.",
																		},
																	},
																},
															},
															"azure_target_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
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
																			Computed:    true,
																			Description: "Specifies id of the Azure resource group used to filter regions in UI.",
																		},
																		"resource_group_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies name of the Azure resource group used to filter regions in UI.",
																		},
																		"source_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
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
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of replication in Unix epoch Timestamp(in microseconds) for a target.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of replication in Unix epoch Timestamp(in microseconds) for a target.",
															},
															"queued_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time when the replication is queued for schedule in Unix epoch Timestamp(in microseconds) for a target.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of the replication for a target. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
															},
															"message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Message about the replication run.",
															},
															"percentage_completed": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the progress in percentage.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about replication data.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total logical size in bytes.",
																		},
																		"logical_bytes_transferred": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total logical bytes transferred.",
																		},
																		"physical_bytes_transferred": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total physical bytes transferred.",
																		},
																	},
																},
															},
															"is_manually_deleted": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the snapshot is deleted manually.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds) for an object.",
															},
															"replication_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Task UID for a replication protection run. This is for tasks that are replicated from another cluster.",
															},
															"entries_changed": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the number of metadata actions completed during the protection run.",
															},
															"is_in_bound": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the direction of the replication. If the snapshot is replicated to this cluster, then isInBound is true. If the snapshot is replicated from this cluster to another cluster, then isInBound is false.",
															},
															"data_lock_constraints": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the dataLock constraints for local or target snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"expiry_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
																		},
																	},
																},
															},
															"on_legal_hold": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the legal hold status for a replication target.",
															},
															"multi_object_replication": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether view based replication was used. In this case, the view containing all objects is replicated as a whole instead of replicating on a per object basis.",
															},
														},
													},
												},
											},
										},
									},
									"archival_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information about archival run for an object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"archival_target_results": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Archival result for an archival target.",
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
															"run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of Protection Group run. 'kRegular' indicates an incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a kRegular schedule captures all the blocks. 'kFull' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized. 'kLog' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time. 'kSystem' indicates system volume backup. It produces an image for bare metal recovery.",
															},
															"is_sla_violated": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicated if SLA has been violated for this run.",
															},
															"snapshot_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Snapshot id for a successful snapshot. This field will not be set if the archival Run fails to take the snapshot.",
															},
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of replication run in Unix epoch Timestamp(in microseconds) for an archival target.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of replication run in Unix epoch Timestamp(in microseconds) for an archival target.",
															},
															"queued_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time when the archival is queued for schedule in Unix epoch Timestamp(in microseconds) for a target.",
															},
															"is_incremental": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether this is an incremental archive. If set to true, this is an incremental archive, otherwise this is a full archive.",
															},
															"is_forever_incremental": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether this is forever incremental or not.",
															},
															"is_cad_archive": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Whether this is CAD archive or not.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of the replication run for an archival target. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
															},
															"message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Message about the archival run.",
															},
															"progress_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task id for archival.",
															},
															"stats_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Run Stats task id for archival.",
															},
															"indexing_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task for indexing.",
															},
															"successful_objects_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the count of objects for which backup was successful.",
															},
															"failed_objects_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the count of objects for which backup failed.",
															},
															"cancelled_objects_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the count of objects for which backup was cancelled.",
															},
															"successful_app_objects_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the count of app objects for which backup was successful.",
															},
															"failed_app_objects_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the count of app objects for which backup failed.",
															},
															"cancelled_app_objects_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the count of app objects for which backup was cancelled.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about archival data.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the logicalSizeBytes.",
																		},
																		"bytes_read": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies total logical bytes read for creating the snapshot.",
																		},
																		"logical_bytes_transferred": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the logical bytes transferred.",
																		},
																		"physical_bytes_transferred": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the physical bytes transferred.",
																		},
																		"avg_logical_transfer_rate_bps": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the average rate of transfer in bytes per second.",
																		},
																		"file_walk_done": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the file system walk is done. Only applicable to file based backups.",
																		},
																		"total_file_count": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total number of file and directory entities visited in this backup. Only applicable to file based backups.",
																		},
																		"backup_file_count": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total number of file and directory entities that are backed up in this run. Only applicable to file based backups.",
																		},
																	},
																},
															},
															"is_manually_deleted": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the snapshot is deleted manually.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
															},
															"data_lock_constraints": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the dataLock constraints for local or target snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"expiry_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
																		},
																	},
																},
															},
															"on_legal_hold": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the legal hold status for a archival target.",
															},
															"worm_properties": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the WORM related properties for this archive.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_archive_worm_compliant": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether this archive run is WORM compliant.",
																		},
																		"worm_non_compliance_reason": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies reason of archive not being worm compliant.",
																		},
																		"worm_expiry_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the time at which the WORM protection expires.",
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
									"cloud_spin_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information about Cloud Spin run for an object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cloud_spin_target_results": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Cloud Spin result for a target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"aws_params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies various resources when converting and deploying a VM to AWS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"custom_tag_list": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies tags of various resources when converting and deploying a VM to AWS.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"key": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies key of the custom tag.",
																					},
																					"value": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies value of the custom tag.",
																					},
																				},
																			},
																		},
																		"region": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the AWS region in which to deploy the VM.",
																		},
																		"subnet_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the subnet within above VPC.",
																		},
																		"vpc_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the Virtual Private Cloud to chose for the instance type.",
																		},
																	},
																},
															},
															"azure_params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies various resources when converting and deploying a VM to Azure.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"availability_set_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the availability set.",
																		},
																		"network_resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the resource group for the selected virtual network.",
																		},
																		"resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the Azure resource group. Its value is globally unique within Azure.",
																		},
																		"storage_account_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
																		},
																		"storage_container_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the storage container within the above storage account.",
																		},
																		"storage_resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the resource group for the selected storage account.",
																		},
																		"temp_vm_resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the temporary Azure resource group.",
																		},
																		"temp_vm_storage_account_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the temporary VM storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
																		},
																		"temp_vm_storage_container_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies id of the temporary VM storage container within the above storage account.",
																		},
																		"temp_vm_subnet_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies Id of the temporary VM subnet within the above virtual network.",
																		},
																		"temp_vm_virtual_network_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies Id of the temporary VM Virtual Network.",
																		},
																	},
																},
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the unique id of the cloud spin entity.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the already added cloud spin target.",
															},
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of Cloud Spin in Unix epoch Timestamp(in microseconds) for a target.",
															},
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of Cloud Spin in Unix epoch Timestamp(in microseconds) for a target.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of the Cloud Spin for a target. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
															},
															"message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Message about the Cloud Spin run.",
															},
															"stats": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies statistics about Cloud Spin data.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"physical_bytes_transferred": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the physical bytes transferred.",
																		},
																	},
																},
															},
															"is_manually_deleted": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the snapshot is deleted manually.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds) for an object.",
															},
															"cloudspin_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Task ID for a CloudSpin protection run.",
															},
															"progress_task_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Progress monitor task id for Cloud Spin run.",
															},
															"data_lock_constraints": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the dataLock constraints for local or target snapshot.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"expiry_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
																		},
																	},
																},
															},
															"on_legal_hold": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the legal hold status for a cloud spin target.",
															},
														},
													},
												},
											},
										},
									},
									"on_legal_hold": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies if object's snapshot is on legal hold.",
									},
								},
							},
						},
						"local_backup_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies summary information about local snapshot run across all objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"run_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of Protection Group run. 'kRegular' indicates an incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a kRegular schedule captures all the blocks. 'kFull' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized. 'kLog' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time. 'kSystem' indicates system volume backup. It produces an image for bare metal recovery. 'kStorageArraySnapshot' indicates storage array snapshot backup.",
									},
									"is_sla_violated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicated if SLA has been violated for this run.",
									},
									"start_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the start time of backup run in Unix epoch Timestamp(in microseconds).",
									},
									"end_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the end time of backup run in Unix epoch Timestamp(in microseconds).",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the backup run. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
									},
									"messages": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Message about the backup run.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"successful_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup was successful.",
									},
									"skipped_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup was skipped.",
									},
									"failed_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup failed.",
									},
									"cancelled_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup was cancelled.",
									},
									"successful_app_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of app objects for which backup was successful.",
									},
									"failed_app_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of app objects for which backup failed.",
									},
									"cancelled_app_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of app objects for which backup was cancelled.",
									},
									"local_snapshot_stats": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies statistics about local snapshot.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"logical_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies total logical size of object(s) in bytes.",
												},
												"bytes_written": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies total size of data in bytes written after taking backup.",
												},
												"bytes_read": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies total logical bytes read for creating the snapshot.",
												},
											},
										},
									},
									"indexing_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task for indexing.",
									},
									"progress_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task id for local backup run.",
									},
									"stats_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Stats task id for local backup run.",
									},
									"data_lock": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "This field is deprecated. Use DataLockConstraints field instead.",
									},
									"local_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task ID for a local protection run.",
									},
									"data_lock_constraints": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the dataLock constraints for local or target snapshot.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
												},
												"expiry_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
												},
											},
										},
									},
								},
							},
						},
						"original_backup_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies summary information about local snapshot run across all objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"run_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of Protection Group run. 'kRegular' indicates an incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a kRegular schedule captures all the blocks. 'kFull' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized. 'kLog' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time. 'kSystem' indicates system volume backup. It produces an image for bare metal recovery. 'kStorageArraySnapshot' indicates storage array snapshot backup.",
									},
									"is_sla_violated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Indicated if SLA has been violated for this run.",
									},
									"start_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the start time of backup run in Unix epoch Timestamp(in microseconds).",
									},
									"end_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the end time of backup run in Unix epoch Timestamp(in microseconds).",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the backup run. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
									},
									"messages": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Message about the backup run.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"successful_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup was successful.",
									},
									"skipped_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup was skipped.",
									},
									"failed_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup failed.",
									},
									"cancelled_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of objects for which backup was cancelled.",
									},
									"successful_app_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of app objects for which backup was successful.",
									},
									"failed_app_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of app objects for which backup failed.",
									},
									"cancelled_app_objects_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the count of app objects for which backup was cancelled.",
									},
									"local_snapshot_stats": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies statistics about local snapshot.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"logical_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies total logical size of object(s) in bytes.",
												},
												"bytes_written": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies total size of data in bytes written after taking backup.",
												},
												"bytes_read": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies total logical bytes read for creating the snapshot.",
												},
											},
										},
									},
									"indexing_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task for indexing.",
									},
									"progress_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task id for local backup run.",
									},
									"stats_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Stats task id for local backup run.",
									},
									"data_lock": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "This field is deprecated. Use DataLockConstraints field instead.",
									},
									"local_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Task ID for a local protection run.",
									},
									"data_lock_constraints": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the dataLock constraints for local or target snapshot.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
												},
												"expiry_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
												},
											},
										},
									},
								},
							},
						},
						"replication_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies summary information about replication run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replication_target_results": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Replication results for each replication target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the id of the cluster.",
												},
												"cluster_incarnation_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the incarnation id of the cluster.",
												},
												"cluster_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the cluster.",
												},
												"aws_target_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
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
																Computed:    true,
																Description: "Specifies id of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
															},
															"region_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies name of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
															},
															"source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the source id of the AWS protection source registered on IBM cluster.",
															},
														},
													},
												},
												"azure_target_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
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
																Computed:    true,
																Description: "Specifies id of the Azure resource group used to filter regions in UI.",
															},
															"resource_group_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies name of the Azure resource group used to filter regions in UI.",
															},
															"source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
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
												"start_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the start time of replication in Unix epoch Timestamp(in microseconds) for a target.",
												},
												"end_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the end time of replication in Unix epoch Timestamp(in microseconds) for a target.",
												},
												"queued_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the time when the replication is queued for schedule in Unix epoch Timestamp(in microseconds) for a target.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of the replication for a target. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
												},
												"message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Message about the replication run.",
												},
												"percentage_completed": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the progress in percentage.",
												},
												"stats": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies statistics about replication data.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"logical_size_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total logical size in bytes.",
															},
															"logical_bytes_transferred": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total logical bytes transferred.",
															},
															"physical_bytes_transferred": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total physical bytes transferred.",
															},
														},
													},
												},
												"is_manually_deleted": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether the snapshot is deleted manually.",
												},
												"expiry_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds) for an object.",
												},
												"replication_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Task UID for a replication protection run. This is for tasks that are replicated from another cluster.",
												},
												"entries_changed": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of metadata actions completed during the protection run.",
												},
												"is_in_bound": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies the direction of the replication. If the snapshot is replicated to this cluster, then isInBound is true. If the snapshot is replicated from this cluster to another cluster, then isInBound is false.",
												},
												"data_lock_constraints": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the dataLock constraints for local or target snapshot.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
															},
														},
													},
												},
												"on_legal_hold": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies the legal hold status for a replication target.",
												},
												"multi_object_replication": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether view based replication was used. In this case, the view containing all objects is replicated as a whole instead of replicating on a per object basis.",
												},
											},
										},
									},
								},
							},
						},
						"archival_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies summary information about archival run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"archival_target_results": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Archival results for each archival target.",
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
												"run_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of Protection Group run. 'kRegular' indicates an incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a kRegular schedule captures all the blocks. 'kFull' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized. 'kLog' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time. 'kSystem' indicates system volume backup. It produces an image for bare metal recovery.",
												},
												"is_sla_violated": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicated if SLA has been violated for this run.",
												},
												"snapshot_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Snapshot id for a successful snapshot. This field will not be set if the archival Run fails to take the snapshot.",
												},
												"start_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the start time of replication run in Unix epoch Timestamp(in microseconds) for an archival target.",
												},
												"end_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the end time of replication run in Unix epoch Timestamp(in microseconds) for an archival target.",
												},
												"queued_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the time when the archival is queued for schedule in Unix epoch Timestamp(in microseconds) for a target.",
												},
												"is_incremental": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether this is an incremental archive. If set to true, this is an incremental archive, otherwise this is a full archive.",
												},
												"is_forever_incremental": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether this is forever incremental or not.",
												},
												"is_cad_archive": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether this is CAD archive or not.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of the replication run for an archival target. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
												},
												"message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Message about the archival run.",
												},
												"progress_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Progress monitor task id for archival.",
												},
												"stats_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Run Stats task id for archival.",
												},
												"indexing_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Progress monitor task for indexing.",
												},
												"successful_objects_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of objects for which backup was successful.",
												},
												"failed_objects_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of objects for which backup failed.",
												},
												"cancelled_objects_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of objects for which backup was cancelled.",
												},
												"successful_app_objects_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of app objects for which backup was successful.",
												},
												"failed_app_objects_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of app objects for which backup failed.",
												},
												"cancelled_app_objects_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of app objects for which backup was cancelled.",
												},
												"stats": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies statistics about archival data.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"logical_size_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the logicalSizeBytes.",
															},
															"bytes_read": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies total logical bytes read for creating the snapshot.",
															},
															"logical_bytes_transferred": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the logical bytes transferred.",
															},
															"physical_bytes_transferred": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the physical bytes transferred.",
															},
															"avg_logical_transfer_rate_bps": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the average rate of transfer in bytes per second.",
															},
															"file_walk_done": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the file system walk is done. Only applicable to file based backups.",
															},
															"total_file_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total number of file and directory entities visited in this backup. Only applicable to file based backups.",
															},
															"backup_file_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total number of file and directory entities that are backed up in this run. Only applicable to file based backups.",
															},
														},
													},
												},
												"is_manually_deleted": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether the snapshot is deleted manually.",
												},
												"expiry_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
												},
												"data_lock_constraints": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the dataLock constraints for local or target snapshot.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
															},
														},
													},
												},
												"on_legal_hold": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies the legal hold status for a archival target.",
												},
												"worm_properties": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the WORM related properties for this archive.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_archive_worm_compliant": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether this archive run is WORM compliant.",
															},
															"worm_non_compliance_reason": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies reason of archive not being worm compliant.",
															},
															"worm_expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time at which the WORM protection expires.",
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
						"cloud_spin_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies summary information about cloud spin run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cloud_spin_target_results": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Cloud Spin results for each Cloud Spin target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aws_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies various resources when converting and deploying a VM to AWS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"custom_tag_list": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies tags of various resources when converting and deploying a VM to AWS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies key of the custom tag.",
																		},
																		"value": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies value of the custom tag.",
																		},
																	},
																},
															},
															"region": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the AWS region in which to deploy the VM.",
															},
															"subnet_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the subnet within above VPC.",
															},
															"vpc_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the Virtual Private Cloud to chose for the instance type.",
															},
														},
													},
												},
												"azure_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies various resources when converting and deploying a VM to Azure.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"availability_set_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the availability set.",
															},
															"network_resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the resource group for the selected virtual network.",
															},
															"resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the Azure resource group. Its value is globally unique within Azure.",
															},
															"storage_account_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
															},
															"storage_container_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the storage container within the above storage account.",
															},
															"storage_resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the resource group for the selected storage account.",
															},
															"temp_vm_resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the temporary Azure resource group.",
															},
															"temp_vm_storage_account_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the temporary VM storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
															},
															"temp_vm_storage_container_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies id of the temporary VM storage container within the above storage account.",
															},
															"temp_vm_subnet_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies Id of the temporary VM subnet within the above virtual network.",
															},
															"temp_vm_virtual_network_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies Id of the temporary VM Virtual Network.",
															},
														},
													},
												},
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the unique id of the cloud spin entity.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the already added cloud spin target.",
												},
												"start_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the start time of Cloud Spin in Unix epoch Timestamp(in microseconds) for a target.",
												},
												"end_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the end time of Cloud Spin in Unix epoch Timestamp(in microseconds) for a target.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of the Cloud Spin for a target. 'Running' indicates that the run is still running. 'Canceled' indicates that the run has been canceled. 'Canceling' indicates that the run is in the process of being canceled. 'Paused' indicates that the ongoing run has been paused. 'Failed' indicates that the run has failed. 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening. 'Succeeded' indicates that the run has finished successfully. 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages. 'Skipped' indicates that the run was skipped.",
												},
												"message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Message about the Cloud Spin run.",
												},
												"stats": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies statistics about Cloud Spin data.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"physical_bytes_transferred": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the physical bytes transferred.",
															},
														},
													},
												},
												"is_manually_deleted": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether the snapshot is deleted manually.",
												},
												"expiry_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds) for an object.",
												},
												"cloudspin_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Task ID for a CloudSpin protection run.",
												},
												"progress_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Progress monitor task id for Cloud Spin run.",
												},
												"data_lock_constraints": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the dataLock constraints for local or target snapshot.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of WORM retention type. <br>'Compliance' implies WORM retention is set for compliance reason. <br>'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"expiry_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the expiry time of attempt in Unix epoch Timestamp (in microseconds).",
															},
														},
													},
												},
												"on_legal_hold": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies the legal hold status for a cloud spin target.",
												},
											},
										},
									},
								},
							},
						},
						"on_legal_hold": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if the Protection Run is on legal hold.",
						},
						"permissions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the list of tenants that have permissions for this protection group run.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"created_at_time_msecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Epoch time when tenant was created.",
									},
									"deleted_at_time_msecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Epoch time when tenant was last updated.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description about the tenant.",
									},
									"external_vendor_metadata": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the additional metadata for the tenant that is specifically set by the external vendors who are responsible for managing tenants. This field will only applicable if tenant creation is happening for a specially provisioned clusters for external vendors.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ibm_tenant_metadata_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the additional metadata for the tenant that is specifically set by the external vendor of type 'IBM'.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"account_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier of the IBM's account ID.",
															},
															"crn": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique CRN associated with the tenant.",
															},
															"custom_properties": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of custom properties associated with the tenant. External vendors can choose to set any properties inside following list. Note that the fields set inside the following will not be available for direct filtering. API callers should make sure that no sensitive information such as passwords is sent in these fields.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"key": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unique key for custom property.",
																		},
																		"value": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the value for the above custom key.",
																		},
																	},
																},
															},
															"liveness_mode": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the current liveness mode of the tenant. This mode may change based on AZ failures when vendor chooses to failover or failback the tenants to other AZs.",
															},
															"metrics_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the metadata for metrics configuration. The metadata defined here will be used by cluster to send the usgae metrics to IBM cloud metering service for calculating the tenant billing.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cos_resource_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the details of COS resource configuration required for posting metrics and trackinb billing information for IBM tenants.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"resource_url": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the resource COS resource configuration endpoint that will be used for fetching bucket usage for a given tenant.",
																					},
																				},
																			},
																		},
																		"iam_metrics_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the IAM configuration that will be used for accessing the billing service in IBM cloud.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"iam_url": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the IAM URL needed to fetch the operator token from IBM. The operator token is needed to make service API calls to IBM billing service.",
																					},
																					"billing_api_key_secret_id": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies Id of the secret that contains the API key.",
																					},
																				},
																			},
																		},
																		"metering_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the metering configuration that will be used for IBM cluster to send the billing details to IBM billing service.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"part_ids": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the list of part identifiers used for metrics identification.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"submission_interval_in_secs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the frequency in seconds at which the metrics will be pushed to IBM billing service from cluster.",
																					},
																					"url": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the base metering URL that will be used by cluster to send the billing information.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"ownership_mode": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the current ownership mode for the tenant. The ownership of the tenant represents the active role for functioning of the tenant.",
															},
															"plan_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Plan Id associated with the tenant. This field is introduced for tracking purposes inside IBM enviournment.",
															},
															"resource_group_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Resource Group ID associated with the tenant.",
															},
															"resource_instance_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Resource Instance ID associated with the tenant. This field is introduced for tracking purposes inside IBM enviournment.",
															},
														},
													},
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of the external vendor. The type specific parameters must be specified the provided type.",
												},
											},
										},
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The tenant id.",
									},
									"is_managed_on_helios": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Flag to indicate if tenant is managed on helios.",
									},
									"last_updated_at_time_msecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Epoch time when tenant was last updated.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the Tenant.",
									},
									"network": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Networking information about a Tenant on a Cluster.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"connector_enabled": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether connector (hybrid extender) is enabled.",
												},
												"cluster_hostname": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The hostname for Cohesity cluster as seen by tenants and as is routable from the tenant's network. Tenant's VLAN's hostname, if available can be used instead but it is mandatory to provide this value if there's no VLAN hostname to use. Also, when set, this field would take precedence over VLAN hostname.",
												},
												"cluster_ips": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Set of IPs as seen from the tenant's network for the Cohesity cluster. Only one from 'clusterHostname' and 'clusterIps' is needed.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Current Status of the Tenant.",
									},
								},
							},
						},
						"is_cloud_archival_direct": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether the run is a CAD run if cloud archive direct feature is enabled. If this field is true, the primary backup copy will only be available at the given archived location.",
						},
						"has_local_snapshot": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether the run has a local snapshot. For cloud retrieved runs there may not be local snapshots.",
						},
						"environment": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the environment of the Protection Group.",
						},
						"externally_triggered_backup_tag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The tag of externally triggered backup job.",
						},
					},
				},
			},
			"total_runs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the count of total runs exist for the given set of filters. The number of runs in single API call are limited and this count can be used to estimate query filter values to get next set of remaining runs. Please note that this field will only be populated if startTimeUsecs or endTimeUsecs or both are specified in query parameters.",
			},
		},
	}
}

func dataSourceIbmBackupRecoveryProtectionGroupRunsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_group_runs", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProtectionGroupRunsOptions := &backuprecoveryv1.GetProtectionGroupRunsOptions{}

	getProtectionGroupRunsOptions.SetID(d.Get("protection_group_id").(string))
	getProtectionGroupRunsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("request_initiator_type"); ok {
		getProtectionGroupRunsOptions.SetRequestInitiatorType(d.Get("request_initiator_type").(string))
	}
	if _, ok := d.GetOk("run_id"); ok {
		getProtectionGroupRunsOptions.SetRunID(d.Get("run_id").(string))
	}
	if _, ok := d.GetOk("start_time_usecs"); ok {
		getProtectionGroupRunsOptions.SetStartTimeUsecs(int64(d.Get("start_time_usecs").(int)))
	}
	if _, ok := d.GetOk("end_time_usecs"); ok {
		getProtectionGroupRunsOptions.SetEndTimeUsecs(int64(d.Get("end_time_usecs").(int)))
	}
	if _, ok := d.GetOk("run_types"); ok {
		var runTypes []string
		for _, v := range d.Get("run_types").([]interface{}) {
			runTypesItem := v.(string)
			runTypes = append(runTypes, runTypesItem)
		}
		getProtectionGroupRunsOptions.SetRunTypes(runTypes)
	}
	if _, ok := d.GetOk("include_object_details"); ok {
		getProtectionGroupRunsOptions.SetIncludeObjectDetails(d.Get("include_object_details").(bool))
	}
	if _, ok := d.GetOk("local_backup_run_status"); ok {
		var localBackupRunStatus []string
		for _, v := range d.Get("local_backup_run_status").([]interface{}) {
			localBackupRunStatusItem := v.(string)
			localBackupRunStatus = append(localBackupRunStatus, localBackupRunStatusItem)
		}
		getProtectionGroupRunsOptions.SetLocalBackupRunStatus(localBackupRunStatus)
	}
	if _, ok := d.GetOk("replication_run_status"); ok {
		var replicationRunStatus []string
		for _, v := range d.Get("replication_run_status").([]interface{}) {
			replicationRunStatusItem := v.(string)
			replicationRunStatus = append(replicationRunStatus, replicationRunStatusItem)
		}
		getProtectionGroupRunsOptions.SetReplicationRunStatus(replicationRunStatus)
	}
	if _, ok := d.GetOk("archival_run_status"); ok {
		var archivalRunStatus []string
		for _, v := range d.Get("archival_run_status").([]interface{}) {
			archivalRunStatusItem := v.(string)
			archivalRunStatus = append(archivalRunStatus, archivalRunStatusItem)
		}
		getProtectionGroupRunsOptions.SetArchivalRunStatus(archivalRunStatus)
	}
	if _, ok := d.GetOk("cloud_spin_run_status"); ok {
		var cloudSpinRunStatus []string
		for _, v := range d.Get("cloud_spin_run_status").([]interface{}) {
			cloudSpinRunStatusItem := v.(string)
			cloudSpinRunStatus = append(cloudSpinRunStatus, cloudSpinRunStatusItem)
		}
		getProtectionGroupRunsOptions.SetCloudSpinRunStatus(cloudSpinRunStatus)
	}
	if _, ok := d.GetOk("num_runs"); ok {
		getProtectionGroupRunsOptions.SetNumRuns(int64(d.Get("num_runs").(int)))
	}
	if _, ok := d.GetOk("exclude_non_restorable_runs"); ok {
		getProtectionGroupRunsOptions.SetExcludeNonRestorableRuns(d.Get("exclude_non_restorable_runs").(bool))
	}
	if _, ok := d.GetOk("run_tags"); ok {
		var runTags []string
		for _, v := range d.Get("run_tags").([]interface{}) {
			runTagsItem := v.(string)
			runTags = append(runTags, runTagsItem)
		}
		getProtectionGroupRunsOptions.SetRunTags(runTags)
	}
	if _, ok := d.GetOk("use_cached_data"); ok {
		getProtectionGroupRunsOptions.SetUseCachedData(d.Get("use_cached_data").(bool))
	}
	if _, ok := d.GetOk("filter_by_end_time"); ok {
		getProtectionGroupRunsOptions.SetFilterByEndTime(d.Get("filter_by_end_time").(bool))
	}
	if _, ok := d.GetOk("snapshot_target_types"); ok {
		var snapshotTargetTypes []string
		for _, v := range d.Get("snapshot_target_types").([]interface{}) {
			snapshotTargetTypesItem := v.(string)
			snapshotTargetTypes = append(snapshotTargetTypes, snapshotTargetTypesItem)
		}
		getProtectionGroupRunsOptions.SetSnapshotTargetTypes(snapshotTargetTypes)
	}
	if _, ok := d.GetOk("only_return_successful_copy_run"); ok {
		getProtectionGroupRunsOptions.SetOnlyReturnSuccessfulCopyRun(d.Get("only_return_successful_copy_run").(bool))
	}
	if _, ok := d.GetOk("filter_by_copy_task_end_time"); ok {
		getProtectionGroupRunsOptions.SetFilterByCopyTaskEndTime(d.Get("filter_by_copy_task_end_time").(bool))
	}

	protectionGroupRunsResponse, _, err := backupRecoveryClient.GetProtectionGroupRunsWithContext(context, getProtectionGroupRunsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProtectionGroupRunsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_protection_group_runs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryProtectionGroupRunsID(d))

	if !core.IsNil(protectionGroupRunsResponse.Runs) {
		runs := []map[string]interface{}{}
		for _, runsItem := range protectionGroupRunsResponse.Runs {
			runsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsProtectionGroupRunToMap(&runsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_group_runs", "read", "runs-to-map").GetDiag()
			}
			runs = append(runs, runsItemMap)
		}
		if err = d.Set("runs", runs); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting runs: %s", err), "(Data) ibm_backup_recovery_protection_group_runs", "read", "set-runs").GetDiag()
		}
	}

	if !core.IsNil(protectionGroupRunsResponse.TotalRuns) {
		if err = d.Set("total_runs", flex.IntValue(protectionGroupRunsResponse.TotalRuns)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_runs: %s", err), "(Data) ibm_backup_recovery_protection_group_runs", "read", "set-total_runs").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryProtectionGroupRunsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryProtectionGroupRunsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsProtectionGroupRunToMap(model *backuprecoveryv1.ProtectionGroupRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ProtectionGroupInstanceID != nil {
		modelMap["protection_group_instance_id"] = flex.IntValue(model.ProtectionGroupInstanceID)
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.IsReplicationRun != nil {
		modelMap["is_replication_run"] = *model.IsReplicationRun
	}
	if model.OriginClusterIdentifier != nil {
		originClusterIdentifierMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsClusterIdentifierToMap(model.OriginClusterIdentifier)
		if err != nil {
			return modelMap, err
		}
		modelMap["origin_cluster_identifier"] = []map[string]interface{}{originClusterIdentifierMap}
	}
	if model.OriginProtectionGroupID != nil {
		modelMap["origin_protection_group_id"] = *model.OriginProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.IsLocalSnapshotsDeleted != nil {
		modelMap["is_local_snapshots_deleted"] = *model.IsLocalSnapshotsDeleted
	}
	if model.Objects != nil {
		objects := []map[string]interface{}{}
		for _, objectsItem := range model.Objects {
			objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsObjectRunResultToMap(&objectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			objects = append(objects, objectsItemMap)
		}
		modelMap["objects"] = objects
	}
	if model.LocalBackupInfo != nil {
		localBackupInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupRunSummaryToMap(model.LocalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_backup_info"] = []map[string]interface{}{localBackupInfoMap}
	}
	if model.OriginalBackupInfo != nil {
		originalBackupInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupRunSummaryToMap(model.OriginalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_backup_info"] = []map[string]interface{}{originalBackupInfoMap}
	}
	if model.ReplicationInfo != nil {
		replicationInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationRunSummaryToMap(model.ReplicationInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["replication_info"] = []map[string]interface{}{replicationInfoMap}
	}
	if model.ArchivalInfo != nil {
		archivalInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalRunSummaryToMap(model.ArchivalInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_info"] = []map[string]interface{}{archivalInfoMap}
	}
	if model.CloudSpinInfo != nil {
		cloudSpinInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinRunSummaryToMap(model.CloudSpinInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["cloud_spin_info"] = []map[string]interface{}{cloudSpinInfoMap}
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	if model.Permissions != nil {
		permissions := []map[string]interface{}{}
		for _, permissionsItem := range model.Permissions {
			permissionsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsTenantToMap(&permissionsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			permissions = append(permissions, permissionsItemMap)
		}
		modelMap["permissions"] = permissions
	}
	if model.IsCloudArchivalDirect != nil {
		modelMap["is_cloud_archival_direct"] = *model.IsCloudArchivalDirect
	}
	if model.HasLocalSnapshot != nil {
		modelMap["has_local_snapshot"] = *model.HasLocalSnapshot
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.ExternallyTriggeredBackupTag != nil {
		modelMap["externally_triggered_backup_tag"] = *model.ExternallyTriggeredBackupTag
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsClusterIdentifierToMap(model *backuprecoveryv1.ClusterIdentifier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterID != nil {
		modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	}
	if model.ClusterIncarnationID != nil {
		modelMap["cluster_incarnation_id"] = flex.IntValue(model.ClusterIncarnationID)
	}
	if model.ClusterName != nil {
		modelMap["cluster_name"] = *model.ClusterName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsObjectRunResultToMap(model *backuprecoveryv1.ObjectRunResult) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Object != nil {
		objectMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsObjectSummaryToMap(model.Object)
		if err != nil {
			return modelMap, err
		}
		modelMap["object"] = []map[string]interface{}{objectMap}
	}
	if model.LocalSnapshotInfo != nil {
		localSnapshotInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupRunToMap(model.LocalSnapshotInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_snapshot_info"] = []map[string]interface{}{localSnapshotInfoMap}
	}
	if model.OriginalBackupInfo != nil {
		originalBackupInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupRunToMap(model.OriginalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_backup_info"] = []map[string]interface{}{originalBackupInfoMap}
	}
	if model.ReplicationInfo != nil {
		replicationInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationRunToMap(model.ReplicationInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["replication_info"] = []map[string]interface{}{replicationInfoMap}
	}
	if model.ArchivalInfo != nil {
		archivalInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalRunToMap(model.ArchivalInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_info"] = []map[string]interface{}{archivalInfoMap}
	}
	if model.CloudSpinInfo != nil {
		cloudSpinInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinRunToMap(model.CloudSpinInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["cloud_spin_info"] = []map[string]interface{}{cloudSpinInfoMap}
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SourceID != nil {
		modelMap["source_id"] = flex.IntValue(model.SourceID)
	}
	if model.SourceName != nil {
		modelMap["source_name"] = *model.SourceName
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.ObjectHash != nil {
		modelMap["object_hash"] = *model.ObjectHash
	}
	if model.ObjectType != nil {
		modelMap["object_type"] = *model.ObjectType
	}
	if model.LogicalSizeBytes != nil {
		modelMap["logical_size_bytes"] = flex.IntValue(model.LogicalSizeBytes)
	}
	if model.UUID != nil {
		modelMap["uuid"] = *model.UUID
	}
	if model.GlobalID != nil {
		modelMap["global_id"] = *model.GlobalID
	}
	if model.ProtectionType != nil {
		modelMap["protection_type"] = *model.ProtectionType
	}
	if model.SharepointSiteSummary != nil {
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsSharepointObjectParamsToMap(model.SharepointSiteSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["sharepoint_site_summary"] = []map[string]interface{}{sharepointSiteSummaryMap}
	}
	if model.OsType != nil {
		modelMap["os_type"] = *model.OsType
	}
	if model.ChildObjects != nil {
		childObjects := []map[string]interface{}{}
		for _, childObjectsItem := range model.ChildObjects {
			childObjectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsSharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsBackupRunToMap(model *backuprecoveryv1.BackupRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SnapshotInfo != nil {
		snapshotInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsSnapshotInfoToMap(model.SnapshotInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["snapshot_info"] = []map[string]interface{}{snapshotInfoMap}
	}
	if model.FailedAttempts != nil {
		failedAttempts := []map[string]interface{}{}
		for _, failedAttemptsItem := range model.FailedAttempts {
			failedAttemptsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupAttemptToMap(&failedAttemptsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			failedAttempts = append(failedAttempts, failedAttemptsItemMap)
		}
		modelMap["failed_attempts"] = failedAttempts
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsSnapshotInfoToMap(model *backuprecoveryv1.SnapshotInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SnapshotID != nil {
		modelMap["snapshot_id"] = *model.SnapshotID
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = *model.StatusMessage
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.AdmittedTimeUsecs != nil {
		modelMap["admitted_time_usecs"] = flex.IntValue(model.AdmittedTimeUsecs)
	}
	if model.PermitGrantTimeUsecs != nil {
		modelMap["permit_grant_time_usecs"] = flex.IntValue(model.PermitGrantTimeUsecs)
	}
	if model.QueueDurationUsecs != nil {
		modelMap["queue_duration_usecs"] = flex.IntValue(model.QueueDurationUsecs)
	}
	if model.SnapshotCreationTimeUsecs != nil {
		modelMap["snapshot_creation_time_usecs"] = flex.IntValue(model.SnapshotCreationTimeUsecs)
	}
	if model.Stats != nil {
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupDataStatsToMap(model.Stats)
		if err != nil {
			return modelMap, err
		}
		modelMap["stats"] = []map[string]interface{}{statsMap}
	}
	if model.ProgressTaskID != nil {
		modelMap["progress_task_id"] = *model.ProgressTaskID
	}
	if model.IndexingTaskID != nil {
		modelMap["indexing_task_id"] = *model.IndexingTaskID
	}
	if model.StatsTaskID != nil {
		modelMap["stats_task_id"] = *model.StatsTaskID
	}
	if model.Warnings != nil {
		modelMap["warnings"] = model.Warnings
	}
	if model.IsManuallyDeleted != nil {
		modelMap["is_manually_deleted"] = *model.IsManuallyDeleted
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.TotalFileCount != nil {
		modelMap["total_file_count"] = flex.IntValue(model.TotalFileCount)
	}
	if model.BackupFileCount != nil {
		modelMap["backup_file_count"] = flex.IntValue(model.BackupFileCount)
	}
	if model.DataLockConstraints != nil {
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsBackupDataStatsToMap(model *backuprecoveryv1.BackupDataStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LogicalSizeBytes != nil {
		modelMap["logical_size_bytes"] = flex.IntValue(model.LogicalSizeBytes)
	}
	if model.BytesWritten != nil {
		modelMap["bytes_written"] = flex.IntValue(model.BytesWritten)
	}
	if model.BytesRead != nil {
		modelMap["bytes_read"] = flex.IntValue(model.BytesRead)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsDataLockConstraintsToMap(model *backuprecoveryv1.DataLockConstraints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Mode != nil {
		modelMap["mode"] = *model.Mode
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsBackupAttemptToMap(model *backuprecoveryv1.BackupAttempt) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.AdmittedTimeUsecs != nil {
		modelMap["admitted_time_usecs"] = flex.IntValue(model.AdmittedTimeUsecs)
	}
	if model.PermitGrantTimeUsecs != nil {
		modelMap["permit_grant_time_usecs"] = flex.IntValue(model.PermitGrantTimeUsecs)
	}
	if model.QueueDurationUsecs != nil {
		modelMap["queue_duration_usecs"] = flex.IntValue(model.QueueDurationUsecs)
	}
	if model.SnapshotCreationTimeUsecs != nil {
		modelMap["snapshot_creation_time_usecs"] = flex.IntValue(model.SnapshotCreationTimeUsecs)
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Stats != nil {
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupDataStatsToMap(model.Stats)
		if err != nil {
			return modelMap, err
		}
		modelMap["stats"] = []map[string]interface{}{statsMap}
	}
	if model.ProgressTaskID != nil {
		modelMap["progress_task_id"] = *model.ProgressTaskID
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationRunToMap(model *backuprecoveryv1.ReplicationRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplicationTargetResults != nil {
		replicationTargetResults := []map[string]interface{}{}
		for _, replicationTargetResultsItem := range model.ReplicationTargetResults {
			replicationTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationTargetResultToMap(&replicationTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			replicationTargetResults = append(replicationTargetResults, replicationTargetResultsItemMap)
		}
		modelMap["replication_target_results"] = replicationTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationTargetResultToMap(model *backuprecoveryv1.ReplicationTargetResult) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterID != nil {
		modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	}
	if model.ClusterIncarnationID != nil {
		modelMap["cluster_incarnation_id"] = flex.IntValue(model.ClusterIncarnationID)
	}
	if model.ClusterName != nil {
		modelMap["cluster_name"] = *model.ClusterName
	}
	if model.AwsTargetConfig != nil {
		awsTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAWSTargetConfigToMap(model.AwsTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_target_config"] = []map[string]interface{}{awsTargetConfigMap}
	}
	if model.AzureTargetConfig != nil {
		azureTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAzureTargetConfigToMap(model.AzureTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_target_config"] = []map[string]interface{}{azureTargetConfigMap}
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.QueuedTimeUsecs != nil {
		modelMap["queued_time_usecs"] = flex.IntValue(model.QueuedTimeUsecs)
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.PercentageCompleted != nil {
		modelMap["percentage_completed"] = flex.IntValue(model.PercentageCompleted)
	}
	if model.Stats != nil {
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationDataStatsToMap(model.Stats)
		if err != nil {
			return modelMap, err
		}
		modelMap["stats"] = []map[string]interface{}{statsMap}
	}
	if model.IsManuallyDeleted != nil {
		modelMap["is_manually_deleted"] = *model.IsManuallyDeleted
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.ReplicationTaskID != nil {
		modelMap["replication_task_id"] = *model.ReplicationTaskID
	}
	if model.EntriesChanged != nil {
		modelMap["entries_changed"] = flex.IntValue(model.EntriesChanged)
	}
	if model.IsInBound != nil {
		modelMap["is_in_bound"] = *model.IsInBound
	}
	if model.DataLockConstraints != nil {
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	if model.MultiObjectReplication != nil {
		modelMap["multi_object_replication"] = *model.MultiObjectReplication
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsAWSTargetConfigToMap(model *backuprecoveryv1.AWSTargetConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsAzureTargetConfigToMap(model *backuprecoveryv1.AzureTargetConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationDataStatsToMap(model *backuprecoveryv1.ReplicationDataStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LogicalSizeBytes != nil {
		modelMap["logical_size_bytes"] = flex.IntValue(model.LogicalSizeBytes)
	}
	if model.LogicalBytesTransferred != nil {
		modelMap["logical_bytes_transferred"] = flex.IntValue(model.LogicalBytesTransferred)
	}
	if model.PhysicalBytesTransferred != nil {
		modelMap["physical_bytes_transferred"] = flex.IntValue(model.PhysicalBytesTransferred)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalRunToMap(model *backuprecoveryv1.ArchivalRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchivalTargetResults != nil {
		archivalTargetResults := []map[string]interface{}{}
		for _, archivalTargetResultsItem := range model.ArchivalTargetResults {
			archivalTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalTargetResultToMap(&archivalTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			archivalTargetResults = append(archivalTargetResults, archivalTargetResultsItemMap)
		}
		modelMap["archival_target_results"] = archivalTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalTargetResultToMap(model *backuprecoveryv1.ArchivalTargetResult) (map[string]interface{}, error) {
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
		tierSettingsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalTargetTierInfoToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	if model.RunType != nil {
		modelMap["run_type"] = *model.RunType
	}
	if model.IsSlaViolated != nil {
		modelMap["is_sla_violated"] = *model.IsSlaViolated
	}
	if model.SnapshotID != nil {
		modelMap["snapshot_id"] = *model.SnapshotID
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.QueuedTimeUsecs != nil {
		modelMap["queued_time_usecs"] = flex.IntValue(model.QueuedTimeUsecs)
	}
	if model.IsIncremental != nil {
		modelMap["is_incremental"] = *model.IsIncremental
	}
	if model.IsForeverIncremental != nil {
		modelMap["is_forever_incremental"] = *model.IsForeverIncremental
	}
	if model.IsCadArchive != nil {
		modelMap["is_cad_archive"] = *model.IsCadArchive
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.ProgressTaskID != nil {
		modelMap["progress_task_id"] = *model.ProgressTaskID
	}
	if model.StatsTaskID != nil {
		modelMap["stats_task_id"] = *model.StatsTaskID
	}
	if model.IndexingTaskID != nil {
		modelMap["indexing_task_id"] = *model.IndexingTaskID
	}
	if model.SuccessfulObjectsCount != nil {
		modelMap["successful_objects_count"] = flex.IntValue(model.SuccessfulObjectsCount)
	}
	if model.FailedObjectsCount != nil {
		modelMap["failed_objects_count"] = flex.IntValue(model.FailedObjectsCount)
	}
	if model.CancelledObjectsCount != nil {
		modelMap["cancelled_objects_count"] = flex.IntValue(model.CancelledObjectsCount)
	}
	if model.SuccessfulAppObjectsCount != nil {
		modelMap["successful_app_objects_count"] = flex.IntValue(model.SuccessfulAppObjectsCount)
	}
	if model.FailedAppObjectsCount != nil {
		modelMap["failed_app_objects_count"] = flex.IntValue(model.FailedAppObjectsCount)
	}
	if model.CancelledAppObjectsCount != nil {
		modelMap["cancelled_app_objects_count"] = flex.IntValue(model.CancelledAppObjectsCount)
	}
	if model.Stats != nil {
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalDataStatsToMap(model.Stats)
		if err != nil {
			return modelMap, err
		}
		modelMap["stats"] = []map[string]interface{}{statsMap}
	}
	if model.IsManuallyDeleted != nil {
		modelMap["is_manually_deleted"] = *model.IsManuallyDeleted
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.DataLockConstraints != nil {
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	if model.WormProperties != nil {
		wormPropertiesMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsWormPropertiesToMap(model.WormProperties)
		if err != nil {
			return modelMap, err
		}
		modelMap["worm_properties"] = []map[string]interface{}{wormPropertiesMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalTargetTierInfoToMap(model *backuprecoveryv1.ArchivalTargetTierInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}

	if model.GoogleTiering != nil {
		googleTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsOracleTiersToMap(model.OracleTiering)
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalDataStatsToMap(model *backuprecoveryv1.ArchivalDataStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LogicalSizeBytes != nil {
		modelMap["logical_size_bytes"] = flex.IntValue(model.LogicalSizeBytes)
	}
	if model.BytesRead != nil {
		modelMap["bytes_read"] = flex.IntValue(model.BytesRead)
	}
	if model.LogicalBytesTransferred != nil {
		modelMap["logical_bytes_transferred"] = flex.IntValue(model.LogicalBytesTransferred)
	}
	if model.PhysicalBytesTransferred != nil {
		modelMap["physical_bytes_transferred"] = flex.IntValue(model.PhysicalBytesTransferred)
	}
	if model.AvgLogicalTransferRateBps != nil {
		modelMap["avg_logical_transfer_rate_bps"] = flex.IntValue(model.AvgLogicalTransferRateBps)
	}
	if model.FileWalkDone != nil {
		modelMap["file_walk_done"] = *model.FileWalkDone
	}
	if model.TotalFileCount != nil {
		modelMap["total_file_count"] = flex.IntValue(model.TotalFileCount)
	}
	if model.BackupFileCount != nil {
		modelMap["backup_file_count"] = flex.IntValue(model.BackupFileCount)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsWormPropertiesToMap(model *backuprecoveryv1.WormProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsArchiveWormCompliant != nil {
		modelMap["is_archive_worm_compliant"] = *model.IsArchiveWormCompliant
	}
	if model.WormNonComplianceReason != nil {
		modelMap["worm_non_compliance_reason"] = *model.WormNonComplianceReason
	}
	if model.WormExpiryTimeUsecs != nil {
		modelMap["worm_expiry_time_usecs"] = flex.IntValue(model.WormExpiryTimeUsecs)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinRunToMap(model *backuprecoveryv1.CloudSpinRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CloudSpinTargetResults != nil {
		cloudSpinTargetResults := []map[string]interface{}{}
		for _, cloudSpinTargetResultsItem := range model.CloudSpinTargetResults {
			cloudSpinTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinTargetResultToMap(&cloudSpinTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cloudSpinTargetResults = append(cloudSpinTargetResults, cloudSpinTargetResultsItemMap)
		}
		modelMap["cloud_spin_target_results"] = cloudSpinTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinTargetResultToMap(model *backuprecoveryv1.CloudSpinTargetResult) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsParams != nil {
		awsParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAwsCloudSpinParamsToMap(model.AwsParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_params"] = []map[string]interface{}{awsParamsMap}
	}
	if model.AzureParams != nil {
		azureParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsAzureCloudSpinParamsToMap(model.AzureParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_params"] = []map[string]interface{}{azureParamsMap}
	}
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Message != nil {
		modelMap["message"] = *model.Message
	}
	if model.Stats != nil {
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinDataStatsToMap(model.Stats)
		if err != nil {
			return modelMap, err
		}
		modelMap["stats"] = []map[string]interface{}{statsMap}
	}
	if model.IsManuallyDeleted != nil {
		modelMap["is_manually_deleted"] = *model.IsManuallyDeleted
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.CloudspinTaskID != nil {
		modelMap["cloudspin_task_id"] = *model.CloudspinTaskID
	}
	if model.ProgressTaskID != nil {
		modelMap["progress_task_id"] = *model.ProgressTaskID
	}
	if model.DataLockConstraints != nil {
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsAwsCloudSpinParamsToMap(model *backuprecoveryv1.AwsCloudSpinParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomTagList != nil {
		customTagList := []map[string]interface{}{}
		for _, customTagListItem := range model.CustomTagList {
			customTagListItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsCustomTagParamsToMap(&customTagListItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			customTagList = append(customTagList, customTagListItemMap)
		}
		modelMap["custom_tag_list"] = customTagList
	}
	modelMap["region"] = flex.IntValue(model.Region)
	if model.SubnetID != nil {
		modelMap["subnet_id"] = flex.IntValue(model.SubnetID)
	}
	if model.VpcID != nil {
		modelMap["vpc_id"] = flex.IntValue(model.VpcID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsCustomTagParamsToMap(model *backuprecoveryv1.CustomTagParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["key"] = *model.Key
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsAzureCloudSpinParamsToMap(model *backuprecoveryv1.AzureCloudSpinParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AvailabilitySetID != nil {
		modelMap["availability_set_id"] = flex.IntValue(model.AvailabilitySetID)
	}
	if model.NetworkResourceGroupID != nil {
		modelMap["network_resource_group_id"] = flex.IntValue(model.NetworkResourceGroupID)
	}
	if model.ResourceGroupID != nil {
		modelMap["resource_group_id"] = flex.IntValue(model.ResourceGroupID)
	}
	if model.StorageAccountID != nil {
		modelMap["storage_account_id"] = flex.IntValue(model.StorageAccountID)
	}
	if model.StorageContainerID != nil {
		modelMap["storage_container_id"] = flex.IntValue(model.StorageContainerID)
	}
	if model.StorageResourceGroupID != nil {
		modelMap["storage_resource_group_id"] = flex.IntValue(model.StorageResourceGroupID)
	}
	if model.TempVmResourceGroupID != nil {
		modelMap["temp_vm_resource_group_id"] = flex.IntValue(model.TempVmResourceGroupID)
	}
	if model.TempVmStorageAccountID != nil {
		modelMap["temp_vm_storage_account_id"] = flex.IntValue(model.TempVmStorageAccountID)
	}
	if model.TempVmStorageContainerID != nil {
		modelMap["temp_vm_storage_container_id"] = flex.IntValue(model.TempVmStorageContainerID)
	}
	if model.TempVmSubnetID != nil {
		modelMap["temp_vm_subnet_id"] = flex.IntValue(model.TempVmSubnetID)
	}
	if model.TempVmVirtualNetworkID != nil {
		modelMap["temp_vm_virtual_network_id"] = flex.IntValue(model.TempVmVirtualNetworkID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinDataStatsToMap(model *backuprecoveryv1.CloudSpinDataStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PhysicalBytesTransferred != nil {
		modelMap["physical_bytes_transferred"] = flex.IntValue(model.PhysicalBytesTransferred)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsBackupRunSummaryToMap(model *backuprecoveryv1.BackupRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RunType != nil {
		modelMap["run_type"] = *model.RunType
	}
	if model.IsSlaViolated != nil {
		modelMap["is_sla_violated"] = *model.IsSlaViolated
	}
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Messages != nil {
		modelMap["messages"] = model.Messages
	}
	if model.SuccessfulObjectsCount != nil {
		modelMap["successful_objects_count"] = flex.IntValue(model.SuccessfulObjectsCount)
	}
	if model.SkippedObjectsCount != nil {
		modelMap["skipped_objects_count"] = flex.IntValue(model.SkippedObjectsCount)
	}
	if model.FailedObjectsCount != nil {
		modelMap["failed_objects_count"] = flex.IntValue(model.FailedObjectsCount)
	}
	if model.CancelledObjectsCount != nil {
		modelMap["cancelled_objects_count"] = flex.IntValue(model.CancelledObjectsCount)
	}
	if model.SuccessfulAppObjectsCount != nil {
		modelMap["successful_app_objects_count"] = flex.IntValue(model.SuccessfulAppObjectsCount)
	}
	if model.FailedAppObjectsCount != nil {
		modelMap["failed_app_objects_count"] = flex.IntValue(model.FailedAppObjectsCount)
	}
	if model.CancelledAppObjectsCount != nil {
		modelMap["cancelled_app_objects_count"] = flex.IntValue(model.CancelledAppObjectsCount)
	}
	if model.LocalSnapshotStats != nil {
		localSnapshotStatsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsBackupDataStatsToMap(model.LocalSnapshotStats)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_snapshot_stats"] = []map[string]interface{}{localSnapshotStatsMap}
	}
	if model.IndexingTaskID != nil {
		modelMap["indexing_task_id"] = *model.IndexingTaskID
	}
	if model.ProgressTaskID != nil {
		modelMap["progress_task_id"] = *model.ProgressTaskID
	}
	if model.StatsTaskID != nil {
		modelMap["stats_task_id"] = *model.StatsTaskID
	}
	if model.DataLock != nil {
		modelMap["data_lock"] = *model.DataLock
	}
	if model.LocalTaskID != nil {
		modelMap["local_task_id"] = *model.LocalTaskID
	}
	if model.DataLockConstraints != nil {
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationRunSummaryToMap(model *backuprecoveryv1.ReplicationRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplicationTargetResults != nil {
		replicationTargetResults := []map[string]interface{}{}
		for _, replicationTargetResultsItem := range model.ReplicationTargetResults {
			replicationTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsReplicationTargetResultToMap(&replicationTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			replicationTargetResults = append(replicationTargetResults, replicationTargetResultsItemMap)
		}
		modelMap["replication_target_results"] = replicationTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalRunSummaryToMap(model *backuprecoveryv1.ArchivalRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchivalTargetResults != nil {
		archivalTargetResults := []map[string]interface{}{}
		for _, archivalTargetResultsItem := range model.ArchivalTargetResults {
			archivalTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsArchivalTargetResultToMap(&archivalTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			archivalTargetResults = append(archivalTargetResults, archivalTargetResultsItemMap)
		}
		modelMap["archival_target_results"] = archivalTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinRunSummaryToMap(model *backuprecoveryv1.CloudSpinRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CloudSpinTargetResults != nil {
		cloudSpinTargetResults := []map[string]interface{}{}
		for _, cloudSpinTargetResultsItem := range model.CloudSpinTargetResults {
			cloudSpinTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsCloudSpinTargetResultToMap(&cloudSpinTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cloudSpinTargetResults = append(cloudSpinTargetResults, cloudSpinTargetResultsItemMap)
		}
		modelMap["cloud_spin_target_results"] = cloudSpinTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedAtTimeMsecs != nil && *(model.CreatedAtTimeMsecs) != 0 {
		modelMap["created_at_time_msecs"] = flex.IntValue(model.CreatedAtTimeMsecs)
	}
	if model.DeletedAtTimeMsecs != nil && *(model.DeletedAtTimeMsecs) != 0 {
		modelMap["deleted_at_time_msecs"] = flex.IntValue(model.DeletedAtTimeMsecs)
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.ExternalVendorMetadata != nil {
		externalVendorMetadataMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
		if err != nil {
			return modelMap, err
		}
		modelMap["external_vendor_metadata"] = []map[string]interface{}{externalVendorMetadataMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.IsManagedOnHelios != nil {
		modelMap["is_managed_on_helios"] = *model.IsManagedOnHelios
	}
	if model.LastUpdatedAtTimeMsecs != nil && *(model.LastUpdatedAtTimeMsecs) != 0 {
		modelMap["last_updated_at_time_msecs"] = flex.IntValue(model.LastUpdatedAtTimeMsecs)
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Network != nil {
		networkMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsTenantNetworkToMap(model.Network)
		if err != nil {
			return modelMap, err
		}
		modelMap["network"] = []map[string]interface{}{networkMap}
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccountID != nil {
		modelMap["account_id"] = *model.AccountID
	}
	if model.Crn != nil {
		modelMap["crn"] = *model.Crn
	}
	if model.CustomProperties != nil {
		customProperties := []map[string]interface{}{}
		for _, customPropertiesItem := range model.CustomProperties {
			customPropertiesItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			customProperties = append(customProperties, customPropertiesItemMap)
		}
		modelMap["custom_properties"] = customProperties
	}
	if model.LivenessMode != nil {
		modelMap["liveness_mode"] = *model.LivenessMode
	}
	if model.MetricsConfig != nil {
		metricsConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantMetricsConfigToMap(model.MetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metrics_config"] = []map[string]interface{}{metricsConfigMap}
	}
	if model.OwnershipMode != nil {
		modelMap["ownership_mode"] = *model.OwnershipMode
	}
	if model.PlanID != nil {
		modelMap["plan_id"] = *model.PlanID
	}
	if model.ResourceGroupID != nil {
		modelMap["resource_group_id"] = *model.ResourceGroupID
	}
	if model.ResourceInstanceID != nil {
		modelMap["resource_instance_id"] = *model.ResourceInstanceID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}
func DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantMetricsConfigToMap(model *backuprecoveryv1.IbmTenantMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CosResourceConfig != nil {
		cosResourceConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantCOSResourceConfigToMap(model.CosResourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_resource_config"] = []map[string]interface{}{cosResourceConfigMap}
	}
	if model.IamMetricsConfig != nil {
		iamMetricsConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantIAMMetricsConfigToMap(model.IamMetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["iam_metrics_config"] = []map[string]interface{}{iamMetricsConfigMap}
	}
	if model.MeteringConfig != nil {
		meteringConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantMeteringConfigToMap(model.MeteringConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metering_config"] = []map[string]interface{}{meteringConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantCOSResourceConfigToMap(model *backuprecoveryv1.IbmTenantCOSResourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceURL != nil {
		modelMap["resource_url"] = *model.ResourceURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantIAMMetricsConfigToMap(model *backuprecoveryv1.IbmTenantIAMMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IAMURL != nil {
		modelMap["iam_url"] = *model.IAMURL
	}
	if model.BillingApiKeySecretID != nil {
		modelMap["billing_api_key_secret_id"] = *model.BillingApiKeySecretID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsIbmTenantMeteringConfigToMap(model *backuprecoveryv1.IbmTenantMeteringConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PartIds != nil {
		modelMap["part_ids"] = model.PartIds
	}
	if model.SubmissionIntervalInSecs != nil {
		modelMap["submission_interval_in_secs"] = flex.IntValue(model.SubmissionIntervalInSecs)
	}
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupRunsTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["connector_enabled"] = *model.ConnectorEnabled
	if model.ClusterHostname != nil {
		modelMap["cluster_hostname"] = *model.ClusterHostname
	}
	if model.ClusterIps != nil {
		modelMap["cluster_ips"] = model.ClusterIps
	}
	return modelMap, nil
}
