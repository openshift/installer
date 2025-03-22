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

func DataSourceIbmBackupRecoveryProtectionGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryProtectionGroupsRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"request_initiator_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of request from UI, which is used for services like magneto to determine the priority of requests.",
			},
			"ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by a list of Protection Group ids.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"names": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by a list of Protection Group names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by Policy ids that are associated with Protection Groups. Only Protection Groups associated with the specified Policy ids, are returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_groups_with_datalock_only": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to only return Protection Groups with a datalock.",
			},
			"environments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by environment types such as 'kVMware', 'kView', etc. Only Protection Groups protecting the specified environment types are returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_active": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter by Inactive or Active Protection Groups. If not set, all Inactive and Active Protection Groups are returned. If true, only Active Protection Groups are returned. If false, only Inactive Protection Groups are returned. When you create a Protection Group on a Primary Cluster with a replication schedule, the Cluster creates an Inactive copy of the Protection Group on the Remote Cluster. In addition, when an Active and running Protection Group is deactivated, the Protection Group becomes Inactive.",
			},
			"is_deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, return only Protection Groups that have been deleted but still have Snapshots associated with them. If false, return all Protection Groups except those Protection Groups that have been deleted and still have Snapshots associated with them. A Protection Group that is deleted with all its Snapshots is not returned for either of these cases.",
			},
			"is_paused": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter by paused or non paused Protection Groups, If not set, all paused and non paused Protection Groups are returned. If true, only paused Protection Groups are returned. If false, only non paused Protection Groups are returned.",
			},
			"last_run_local_backup_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by last local backup run status.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_run_replication_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by last remote replication run status.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_run_archival_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by last cloud archival run status.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_run_cloud_spin_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by last cloud spin run status.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"last_run_any_status": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by last any run status.<br> 'Running' indicates that the run is still running.<br> 'Canceled' indicates that the run has been canceled.<br> 'Canceling' indicates that the run is in the process of being canceled.<br> 'Failed' indicates that the run has failed.<br> 'Missed' indicates that the run was unable to take place at the scheduled time because the previous run was still happening.<br> 'Succeeded' indicates that the run has finished successfully.<br> 'SucceededWithWarning' indicates that the run finished successfully, but there were some warning messages.<br> 'Paused' indicates that the ongoing run has been paused.<br> 'Skipped' indicates that the run was skipped.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_last_run_sla_violated": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, return Protection Groups for which last run SLA was violated.",
			},
			"include_last_run_info": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the response will include last run info. If it is false or not specified, the last run info won't be returned.",
			},
			"prune_excluded_source_ids": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the response will not include the list of excluded source IDs in groups that contain this field. This can be set to true in order to improve performance if excluded source IDs are not needed by the user.",
			},
			"prune_source_ids": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the response will exclude the list of source IDs within the group specified.",
			},
			"use_cached_data": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether we can serve the GET request from the read replica cache. There is a lag of 15 seconds between the read replica and primary data source.",
			},
			"source_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by Source ids that are associated with Protection Groups. Only Protection Groups associated with the specified Source ids, are returned.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"protection_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of Protection Groups which were returned by the request.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the ID of the Protection Group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the name of the Protection Group.",
						},
						"cluster_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the cluster ID.",
						},
						"region_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the region ID.",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the unique id of the Protection Policy associated with the Protection Group. The Policy provides retry settings Protection Schedules, Priority, SLA, etc.",
						},
						"priority": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the priority of the Protection Group.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies a description of the Protection Group.",
						},
						"start_time": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the time of day. Used for scheduling purposes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the hour of the day (0-23).",
									},
									"minute": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the minute of the hour (0-59).",
									},
									"time_zone": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the time zone of the user. If not specified, default value is assumed as America/Los_Angeles.",
									},
								},
							},
						},
						"end_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the end time in micro seconds for this Protection Group. If this is not specified, the Protection Group won't be ended.",
						},
						"last_modified_timestamp_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the last time this protection group was updated. If this is passed into a PUT request, then the backend will validate that the timestamp passed in matches the time that the protection group was actually last modified. If the two timestamps do not match, then the request will be rejected with a stale error.",
						},
						"alert_policy": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies a policy for alerting users of the status of a Protection Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backup_run_status": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the run status for which the user would like to receive alerts.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"alert_targets": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of targets to receive the alerts.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"email_address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies an email address to receive an alert.",
												},
												"language": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the language of the delivery target. Default value is 'en-us'.",
												},
												"recipient_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the recipient type of email recipient. Default value is 'kTo'.",
												},
											},
										},
									},
									"raise_object_level_failure_alert": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether object level alerts are raised for backup failures after the backup run.",
									},
									"raise_object_level_failure_alert_after_last_attempt": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether object level alerts are raised for backup failures after last backup attempt.",
									},
									"raise_object_level_failure_alert_after_each_attempt": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether object level alerts are raised for backup failures after each backup attempt.",
									},
								},
							},
						},
						"sla": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the SLA parameters for this Protection Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"backup_run_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the type of run this rule should apply to.",
									},
									"sla_minutes": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of minutes allotted to a run of the specified type before SLA is considered violated.",
									},
								},
							},
						},
						"qos_policy": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies whether the Protection Group will be written to HDD or SSD.",
						},
						"abort_in_blackouts": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether currently executing jobs should abort if a blackout period specified by a policy starts. Available only if the selected policy has at least one blackout period. Default value is false.",
						},
						"pause_in_blackouts": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies whether currently executing jobs should be paused if a blackout period specified by a policy starts. Available only if the selected policy has at least one blackout period. Default value is false. This field should not be set to true if 'abortInBlackouts' is sent as true.",
						},
						"is_active": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if the Protection Group is active or not.",
						},
						"is_deleted": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if the Protection Group has been deleted.",
						},
						"is_paused": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if the the Protection Group is paused. New runs are not scheduled for the paused Protection Groups. Active run if any is not impacted.",
						},
						"environment": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the environment of the Protection Group.",
						},
						"last_run": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the parameters which are common between Protection Group runs of all Protection Groups.",
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
						"permissions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the list of tenants that have permissions for this protection group.",
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
						"is_protect_once": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies if the the Protection Group is using a protect once type of policy. This field is helpful to identify run happen for this group.",
						},
						"missing_entities": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the Information about missing entities.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the ID of the object.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the name of the object.",
									},
									"parent_source_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the id of the parent source of the object.",
									},
									"parent_source_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the name of the parent source of the object.",
									},
								},
							},
						},
						"invalid_entities": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the Information about invalid entities. An entity will be considered invalid if it is part of an active protection group but has lost compatibility for the given backup type.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the ID of the object.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the name of the object.",
									},
									"parent_source_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the id of the parent source of the object.",
									},
									"parent_source_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the name of the parent source of the object.",
									},
								},
							},
						},
						"num_protected_objects": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of protected objects of the Protection Group.",
						},
						"advanced_configs": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the advanced configuration for a protection job.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "key.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "value.",
									},
								},
							},
						},
						"physical_params": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the Physical Protection Group type.",
									},
									"volume_protection_type_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters which are specific to Volume based physical Protection Groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"objects": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the ID of the object protected.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object protected.",
															},
															"volume_guids": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of GUIDs of volumes protected. If empty, then all volumes will be protected by default.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"enable_system_backup": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether or not to take a system backup. Applicable only for windows sources.",
															},
															"excluded_vss_writers": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies writer names which should be excluded from physical volume based backups for a given source.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"indexing_policy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings for indexing files found in an Object (such as a VM) so these files can be searched and recovered. This also specifies inclusion and exclusion rules that determine the directories to index.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enable_indexing": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the files found in an Object (such as a VM) should be indexed. If true (the default), files are indexed.",
															},
															"include_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Indexed Directories. Specifies a list of directories to index. Regular expression can also be specified to provide the directory paths. Example: /Users/<wildcard>/AppData.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"exclude_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Excluded Directories. Specifies a list of directories to exclude from indexing.Regular expression can also be specified to provide the directory paths. Example: /Users/<wildcard>/AppData.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"perform_source_side_deduplication": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not to perform source side deduplication on this Protection Group.",
												},
												"quiesce": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies Whether to take app-consistent snapshots by quiescing apps and the filesystem before taking a backup.",
												},
												"continue_on_quiesce_failure": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to continue backing up on quiesce failure.",
												},
												"incremental_backup_after_restart": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not to perform an incremental backup after the server restarts. This is applicable to windows environments.",
												},
												"pre_post_script": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the params for pre and post scripts.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"pre_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PreBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																		"continue_on_error": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the script needs to continue even if there is an occurence of an error. If this flag is set to true, then Backup Run will start even if the pre backup script fails. If not specified or false, then backup run will not start when script fails.",
																		},
																	},
																},
															},
															"post_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PostBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																	},
																},
															},
														},
													},
												},
												"dedup_exclusion_source_ids": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies ids of sources for which deduplication has to be disabled.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"excluded_vss_writers": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies writer names which should be excluded from physical volume based backups.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"cobmr_backup": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to take a CoBMR backup.",
												},
											},
										},
									},
									"file_protection_type_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters which are specific to Physical related Protection Groups.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"excluded_vss_writers": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies writer names which should be excluded from physical file based backups.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"objects": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of objects protected by this Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"excluded_vss_writers": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies writer names which should be excluded from physical file based backups.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the ID of the object protected.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object protected.",
															},
															"file_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a list of file paths to be protected by this Protection Group.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"included_path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a path to be included on the source. All paths under this path will be included unless they are specifically mentioned in excluded paths.",
																		},
																		"excluded_paths": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a set of paths nested under the include path which should be excluded from the Protection Group.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"skip_nested_volumes": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether to skip any nested volumes (both local and network) that are mounted under include path. Applicable only for windows sources.",
																		},
																	},
																},
															},
															"uses_path_level_skip_nested_volume_setting": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether path level or object level skip nested volume setting will be used.",
															},
															"nested_volume_types_to_skip": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies mount types of nested volumes to be skipped.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"follow_nas_symlink_target": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether to follow NAS target pointed by symlink for windows sources.",
															},
															"metadata_file_path": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the path of metadatafile on source. This file contains absolute paths of files that needs to be backed up on the same source.",
															},
														},
													},
												},
												"indexing_policy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings for indexing files found in an Object (such as a VM) so these files can be searched and recovered. This also specifies inclusion and exclusion rules that determine the directories to index.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enable_indexing": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the files found in an Object (such as a VM) should be indexed. If true (the default), files are indexed.",
															},
															"include_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Indexed Directories. Specifies a list of directories to index. Regular expression can also be specified to provide the directory paths. Example: /Users/<wildcard>/AppData.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"exclude_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Excluded Directories. Specifies a list of directories to exclude from indexing.Regular expression can also be specified to provide the directory paths. Example: /Users/<wildcard>/AppData.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"perform_source_side_deduplication": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not to perform source side deduplication on this Protection Group.",
												},
												"perform_brick_based_deduplication": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not to perform brick based deduplication on this Protection Group.",
												},
												"task_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the timeouts for all the objects inside this Protection Group, for both full and incremental backups.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"timeout_mins": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the timeout in mins.",
															},
															"backup_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The scheduled backup type(kFull, kRegular etc.).",
															},
														},
													},
												},
												"quiesce": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies Whether to take app-consistent snapshots by quiescing apps and the filesystem before taking a backup.",
												},
												"continue_on_quiesce_failure": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to continue backing up on quiesce failure.",
												},
												"cobmr_backup": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to take CoBMR backup.",
												},
												"pre_post_script": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the params for pre and post scripts.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"pre_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PreBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																		"continue_on_error": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the script needs to continue even if there is an occurence of an error. If this flag is set to true, then Backup Run will start even if the pre backup script fails. If not specified or false, then backup run will not start when script fails.",
																		},
																	},
																},
															},
															"post_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PostBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																	},
																},
															},
														},
													},
												},
												"dedup_exclusion_source_ids": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies ids of sources for which deduplication has to be disabled.",
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
												},
												"global_exclude_paths": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies global exclude filters which are applied to all sources in a job.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"global_exclude_fs": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies global exclude filesystems which are applied to all sources in a job.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ignorable_errors": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the Errors to be ignored in error db.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"allow_parallel_runs": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not this job can have parallel runs.",
												},
											},
										},
									},
								},
							},
						},
						"mssql_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the parameters specific to MSSQL Protection Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_protection_type_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the params to create a File based MSSQL Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aag_backup_preference_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the preference type for backing up databases that are part of an AAG. If not specified, then default preferences of the AAG server are applied. This field wont be applicable if user DB preference is set to skip AAG databases.",
												},
												"advanced_settings": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "This is used to regulate certain gflag values from the UI. The values passed by the user from the UI will be used for the respective gflags.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cloned_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error if SQL database is cloned.",
															},
															"db_backup_if_not_online_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error if SQL database is not online.",
															},
															"missing_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Fail the backup job when the database is missing. The database may be missing if it is deleted or corrupted.",
															},
															"offline_restoring_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Fail the backup job when database is offline or restoring.",
															},
															"read_only_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to skip backup for read-only SQL databases.",
															},
															"report_all_non_autoprotect_db_errors": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error for all dbs in non-autoprotect jobs.",
															},
														},
													},
												},
												"backup_system_dbs": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to backup system databases. If not specified then parameter is set to true.",
												},
												"exclude_filters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of exclusion filters applied during the group creation or edit. These exclusion filters can be wildcard supported strings or regular expressions. Objects satisfying the will filters will be excluded during backup and also auto protected objects will be ignored if filtered by any of the filters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"filter_string": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the filter string using wildcard supported strings or regular expressions.",
															},
															"is_regular_expression": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the provided filter string is a regular expression or not. This needs to be explicitly set to true if user is trying to filter by regular expressions. Not providing this value in case of regular expression can result in unintended results. The default value is assumed to be false.",
															},
														},
													},
												},
												"full_backups_copy_only": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether full backups should be copy-only.",
												},
												"log_backup_num_streams": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of streams to be used for log backups.",
												},
												"log_backup_with_clause": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the WithClause to be used for log backups.",
												},
												"pre_post_script": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the params for pre and post scripts.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"pre_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PreBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																		"continue_on_error": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the script needs to continue even if there is an occurence of an error. If this flag is set to true, then Backup Run will start even if the pre backup script fails. If not specified or false, then backup run will not start when script fails.",
																		},
																	},
																},
															},
															"post_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PostBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																	},
																},
															},
														},
													},
												},
												"use_aag_preferences_from_server": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not the AAG backup preferences specified on the SQL Server host should be used.",
												},
												"user_db_backup_preference_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the preference type for backing up user databases on the host.",
												},
												"additional_host_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings which are to be applied to specific host containers in this protection group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"disable_source_side_deduplication": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether or not to disable source side deduplication on this source. The default behavior is false unless the user has set 'performSourceSideDeduplication' to true.",
															},
															"host_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the id of the host container on which databases are hosted.",
															},
															"host_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the host container on which databases are hosted.",
															},
														},
													},
												},
												"objects": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of object params to be protected.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the ID of the object being protected. If this is a non leaf level object, then the object will be auto-protected unless leaf objects are specified for exclusion.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object being protected.",
															},
															"source_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of source being protected.",
															},
														},
													},
												},
												"perform_source_side_deduplication": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not to perform source side deduplication on this Protection Group.",
												},
											},
										},
									},
									"native_protection_type_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the params to create a Native based MSSQL Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aag_backup_preference_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the preference type for backing up databases that are part of an AAG. If not specified, then default preferences of the AAG server are applied. This field wont be applicable if user DB preference is set to skip AAG databases.",
												},
												"advanced_settings": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "This is used to regulate certain gflag values from the UI. The values passed by the user from the UI will be used for the respective gflags.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cloned_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error if SQL database is cloned.",
															},
															"db_backup_if_not_online_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error if SQL database is not online.",
															},
															"missing_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Fail the backup job when the database is missing. The database may be missing if it is deleted or corrupted.",
															},
															"offline_restoring_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Fail the backup job when database is offline or restoring.",
															},
															"read_only_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to skip backup for read-only SQL databases.",
															},
															"report_all_non_autoprotect_db_errors": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error for all dbs in non-autoprotect jobs.",
															},
														},
													},
												},
												"backup_system_dbs": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to backup system databases. If not specified then parameter is set to true.",
												},
												"exclude_filters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of exclusion filters applied during the group creation or edit. These exclusion filters can be wildcard supported strings or regular expressions. Objects satisfying the will filters will be excluded during backup and also auto protected objects will be ignored if filtered by any of the filters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"filter_string": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the filter string using wildcard supported strings or regular expressions.",
															},
															"is_regular_expression": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the provided filter string is a regular expression or not. This needs to be explicitly set to true if user is trying to filter by regular expressions. Not providing this value in case of regular expression can result in unintended results. The default value is assumed to be false.",
															},
														},
													},
												},
												"full_backups_copy_only": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether full backups should be copy-only.",
												},
												"log_backup_num_streams": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of streams to be used for log backups.",
												},
												"log_backup_with_clause": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the WithClause to be used for log backups.",
												},
												"pre_post_script": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the params for pre and post scripts.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"pre_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PreBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																		"continue_on_error": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the script needs to continue even if there is an occurence of an error. If this flag is set to true, then Backup Run will start even if the pre backup script fails. If not specified or false, then backup run will not start when script fails.",
																		},
																	},
																},
															},
															"post_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PostBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																	},
																},
															},
														},
													},
												},
												"use_aag_preferences_from_server": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not the AAG backup preferences specified on the SQL Server host should be used.",
												},
												"user_db_backup_preference_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the preference type for backing up user databases on the host.",
												},
												"num_streams": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of streams to be used.",
												},
												"objects": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of object params to be protected.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the ID of the object being protected. If this is a non leaf level object, then the object will be auto-protected unless leaf objects are specified for exclusion.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object being protected.",
															},
															"source_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of source being protected.",
															},
														},
													},
												},
												"with_clause": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the WithClause to be used.",
												},
											},
										},
									},
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the MSSQL Protection Group type.",
									},
									"volume_protection_type_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the params to create a Volume based MSSQL Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aag_backup_preference_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the preference type for backing up databases that are part of an AAG. If not specified, then default preferences of the AAG server are applied. This field wont be applicable if user DB preference is set to skip AAG databases.",
												},
												"advanced_settings": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "This is used to regulate certain gflag values from the UI. The values passed by the user from the UI will be used for the respective gflags.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cloned_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error if SQL database is cloned.",
															},
															"db_backup_if_not_online_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error if SQL database is not online.",
															},
															"missing_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Fail the backup job when the database is missing. The database may be missing if it is deleted or corrupted.",
															},
															"offline_restoring_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Fail the backup job when database is offline or restoring.",
															},
															"read_only_db_backup_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to skip backup for read-only SQL databases.",
															},
															"report_all_non_autoprotect_db_errors": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Whether to report error for all dbs in non-autoprotect jobs.",
															},
														},
													},
												},
												"backup_system_dbs": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to backup system databases. If not specified then parameter is set to true.",
												},
												"exclude_filters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of exclusion filters applied during the group creation or edit. These exclusion filters can be wildcard supported strings or regular expressions. Objects satisfying the will filters will be excluded during backup and also auto protected objects will be ignored if filtered by any of the filters.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"filter_string": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the filter string using wildcard supported strings or regular expressions.",
															},
															"is_regular_expression": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the provided filter string is a regular expression or not. This needs to be explicitly set to true if user is trying to filter by regular expressions. Not providing this value in case of regular expression can result in unintended results. The default value is assumed to be false.",
															},
														},
													},
												},
												"full_backups_copy_only": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether full backups should be copy-only.",
												},
												"log_backup_num_streams": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of streams to be used for log backups.",
												},
												"log_backup_with_clause": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the WithClause to be used for log backups.",
												},
												"pre_post_script": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the params for pre and post scripts.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"pre_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PreBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																		"continue_on_error": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the script needs to continue even if there is an occurence of an error. If this flag is set to true, then Backup Run will start even if the pre backup script fails. If not specified or false, then backup run will not start when script fails.",
																		},
																	},
																},
															},
															"post_script": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the common params for PostBackup scripts.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the absolute path to the script on the remote host.",
																		},
																		"params": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the arguments or parameters and values to pass into the remote script. For example if the script expects values for the 'database' and 'user' parameters, specify the parameters and values using the following string: \"database=myDatabase user=me\".",
																		},
																		"timeout_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the timeout of the script in seconds. The script will be killed if it exceeds this value. By default, no timeout will occur if left empty.",
																		},
																		"is_active": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the script should be enabled, default value set to true.",
																		},
																	},
																},
															},
														},
													},
												},
												"use_aag_preferences_from_server": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or not the AAG backup preferences specified on the SQL Server host should be used.",
												},
												"user_db_backup_preference_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the preference type for backing up user databases on the host.",
												},
												"additional_host_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings which are to be applied to specific host containers in this protection group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enable_system_backup": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether to enable system/bmr backup using 3rd party tools installed on agent host.",
															},
															"host_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the id of the host container on which databases are hosted.",
															},
															"host_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the host container on which databases are hosted.",
															},
															"volume_guids": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of volume GUIDs to be protected. If not specified, all the volumes of the host will be protected. Note that volumes of host on which databases are hosted are protected even if its not mentioned in this list.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"backup_db_volumes_only": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to only backup volumes on which the specified databases reside. If not specified (default), all the volumes of the host will be protected.",
												},
												"incremental_backup_after_restart": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether or to perform incremental backups the first time after a server restarts. By default, a full backup will be performed.",
												},
												"indexing_policy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings for indexing files found in an Object (such as a VM) so these files can be searched and recovered. This also specifies inclusion and exclusion rules that determine the directories to index.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enable_indexing": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the files found in an Object (such as a VM) should be indexed. If true (the default), files are indexed.",
															},
															"include_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Indexed Directories. Specifies a list of directories to index. Regular expression can also be specified to provide the directory paths. Example: /Users/<wildcard>/AppData.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"exclude_paths": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Excluded Directories. Specifies a list of directories to exclude from indexing.Regular expression can also be specified to provide the directory paths. Example: /Users/<wildcard>/AppData.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"objects": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of object ids to be protected.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the ID of the object being protected. If this is a non leaf level object, then the object will be auto-protected unless leaf objects are specified for exclusion.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object being protected.",
															},
															"source_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of source being protected.",
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

func dataSourceIbmBackupRecoveryProtectionGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_groups", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProtectionGroupsOptions := &backuprecoveryv1.GetProtectionGroupsOptions{}

	getProtectionGroupsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("request_initiator_type"); ok {
		getProtectionGroupsOptions.SetRequestInitiatorType(d.Get("request_initiator_type").(string))
	}
	if _, ok := d.GetOk("ids"); ok {
		var ids []string
		for _, v := range d.Get("ids").([]interface{}) {
			idsItem := v.(string)
			ids = append(ids, idsItem)
		}
		getProtectionGroupsOptions.SetIds(ids)
	}
	if _, ok := d.GetOk("names"); ok {
		var names []string
		for _, v := range d.Get("names").([]interface{}) {
			namesItem := v.(string)
			names = append(names, namesItem)
		}
		getProtectionGroupsOptions.SetNames(names)
	}
	if _, ok := d.GetOk("policy_ids"); ok {
		var policyIds []string
		for _, v := range d.Get("policy_ids").([]interface{}) {
			policyIdsItem := v.(string)
			policyIds = append(policyIds, policyIdsItem)
		}
		getProtectionGroupsOptions.SetPolicyIds(policyIds)
	}
	if _, ok := d.GetOk("include_groups_with_datalock_only"); ok {
		getProtectionGroupsOptions.SetIncludeGroupsWithDatalockOnly(d.Get("include_groups_with_datalock_only").(bool))
	}
	if _, ok := d.GetOk("environments"); ok {
		var environments []string
		for _, v := range d.Get("environments").([]interface{}) {
			environmentsItem := v.(string)
			environments = append(environments, environmentsItem)
		}
		getProtectionGroupsOptions.SetEnvironments(environments)
	}
	if _, ok := d.GetOk("is_active"); ok {
		getProtectionGroupsOptions.SetIsActive(d.Get("is_active").(bool))
	}
	if _, ok := d.GetOk("is_deleted"); ok {
		getProtectionGroupsOptions.SetIsDeleted(d.Get("is_deleted").(bool))
	}
	if _, ok := d.GetOk("is_paused"); ok {
		getProtectionGroupsOptions.SetIsPaused(d.Get("is_paused").(bool))
	}
	if _, ok := d.GetOk("last_run_local_backup_status"); ok {
		var lastRunLocalBackupStatus []string
		for _, v := range d.Get("last_run_local_backup_status").([]interface{}) {
			lastRunLocalBackupStatusItem := v.(string)
			lastRunLocalBackupStatus = append(lastRunLocalBackupStatus, lastRunLocalBackupStatusItem)
		}
		getProtectionGroupsOptions.SetLastRunLocalBackupStatus(lastRunLocalBackupStatus)
	}
	if _, ok := d.GetOk("last_run_replication_status"); ok {
		var lastRunReplicationStatus []string
		for _, v := range d.Get("last_run_replication_status").([]interface{}) {
			lastRunReplicationStatusItem := v.(string)
			lastRunReplicationStatus = append(lastRunReplicationStatus, lastRunReplicationStatusItem)
		}
		getProtectionGroupsOptions.SetLastRunReplicationStatus(lastRunReplicationStatus)
	}
	if _, ok := d.GetOk("last_run_archival_status"); ok {
		var lastRunArchivalStatus []string
		for _, v := range d.Get("last_run_archival_status").([]interface{}) {
			lastRunArchivalStatusItem := v.(string)
			lastRunArchivalStatus = append(lastRunArchivalStatus, lastRunArchivalStatusItem)
		}
		getProtectionGroupsOptions.SetLastRunArchivalStatus(lastRunArchivalStatus)
	}
	if _, ok := d.GetOk("last_run_cloud_spin_status"); ok {
		var lastRunCloudSpinStatus []string
		for _, v := range d.Get("last_run_cloud_spin_status").([]interface{}) {
			lastRunCloudSpinStatusItem := v.(string)
			lastRunCloudSpinStatus = append(lastRunCloudSpinStatus, lastRunCloudSpinStatusItem)
		}
		getProtectionGroupsOptions.SetLastRunCloudSpinStatus(lastRunCloudSpinStatus)
	}
	if _, ok := d.GetOk("last_run_any_status"); ok {
		var lastRunAnyStatus []string
		for _, v := range d.Get("last_run_any_status").([]interface{}) {
			lastRunAnyStatusItem := v.(string)
			lastRunAnyStatus = append(lastRunAnyStatus, lastRunAnyStatusItem)
		}
		getProtectionGroupsOptions.SetLastRunAnyStatus(lastRunAnyStatus)
	}
	if _, ok := d.GetOk("is_last_run_sla_violated"); ok {
		getProtectionGroupsOptions.SetIsLastRunSlaViolated(d.Get("is_last_run_sla_violated").(bool))
	}
	if _, ok := d.GetOk("include_last_run_info"); ok {
		getProtectionGroupsOptions.SetIncludeLastRunInfo(d.Get("include_last_run_info").(bool))
	}
	if _, ok := d.GetOk("prune_excluded_source_ids"); ok {
		getProtectionGroupsOptions.SetPruneExcludedSourceIds(d.Get("prune_excluded_source_ids").(bool))
	}
	if _, ok := d.GetOk("prune_source_ids"); ok {
		getProtectionGroupsOptions.SetPruneSourceIds(d.Get("prune_source_ids").(bool))
	}
	if _, ok := d.GetOk("use_cached_data"); ok {
		getProtectionGroupsOptions.SetUseCachedData(d.Get("use_cached_data").(bool))
	}
	if _, ok := d.GetOk("source_ids"); ok {
		var sourceIds []int64
		for _, v := range d.Get("source_ids").([]interface{}) {
			sourceIdsItem := int64(v.(int))
			sourceIds = append(sourceIds, sourceIdsItem)
		}
		getProtectionGroupsOptions.SetSourceIds(sourceIds)
	}

	protectionGroupsResponse, _, err := backupRecoveryClient.GetProtectionGroupsWithContext(context, getProtectionGroupsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProtectionGroupsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_protection_groups", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryProtectionGroupsID(d))

	if !core.IsNil(protectionGroupsResponse.ProtectionGroups) {
		protectionGroups := []map[string]interface{}{}
		for _, protectionGroupsItem := range protectionGroupsResponse.ProtectionGroups {
			protectionGroupsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsProtectionGroupResponseToMap(&protectionGroupsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_groups", "read", "protection_groups-to-map").GetDiag()
			}
			protectionGroups = append(protectionGroups, protectionGroupsItemMap)
		}
		if err = d.Set("protection_groups", protectionGroups); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting protection_groups: %s", err), "(Data) ibm_backup_recovery_protection_groups", "read", "set-protection_groups").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryProtectionGroupsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryProtectionGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryProtectionGroupsProtectionGroupResponseToMap(model *backuprecoveryv1.ProtectionGroupResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ClusterID != nil {
		modelMap["cluster_id"] = *model.ClusterID
	}
	if model.RegionID != nil {
		modelMap["region_id"] = *model.RegionID
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.Priority != nil {
		modelMap["priority"] = *model.Priority
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.StartTime != nil {
		startTimeMap, err := DataSourceIbmBackupRecoveryProtectionGroupsTimeOfDayToMap(model.StartTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["start_time"] = []map[string]interface{}{startTimeMap}
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.LastModifiedTimestampUsecs != nil {
		modelMap["last_modified_timestamp_usecs"] = flex.IntValue(model.LastModifiedTimestampUsecs)
	}
	if model.AlertPolicy != nil {
		alertPolicyMap, err := DataSourceIbmBackupRecoveryProtectionGroupsProtectionGroupAlertingPolicyToMap(model.AlertPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["alert_policy"] = []map[string]interface{}{alertPolicyMap}
	}
	if model.Sla != nil {
		sla := []map[string]interface{}{}
		for _, slaItem := range model.Sla {
			slaItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsSlaRuleToMap(&slaItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			sla = append(sla, slaItemMap)
		}
		modelMap["sla"] = sla
	}
	if model.QosPolicy != nil {
		modelMap["qos_policy"] = *model.QosPolicy
	}
	if model.AbortInBlackouts != nil {
		modelMap["abort_in_blackouts"] = *model.AbortInBlackouts
	}
	if model.PauseInBlackouts != nil {
		modelMap["pause_in_blackouts"] = *model.PauseInBlackouts
	}
	if model.IsActive != nil {
		modelMap["is_active"] = *model.IsActive
	}
	if model.IsDeleted != nil {
		modelMap["is_deleted"] = *model.IsDeleted
	}
	if model.IsPaused != nil {
		modelMap["is_paused"] = *model.IsPaused
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.LastRun != nil {
		lastRunMap, err := DataSourceIbmBackupRecoveryProtectionGroupsProtectionGroupRunToMap(model.LastRun)
		if err != nil {
			return modelMap, err
		}
		modelMap["last_run"] = []map[string]interface{}{lastRunMap}
	}
	if model.Permissions != nil {
		permissions := []map[string]interface{}{}
		for _, permissionsItem := range model.Permissions {
			permissionsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsTenantToMap(&permissionsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			permissions = append(permissions, permissionsItemMap)
		}
		modelMap["permissions"] = permissions
	}
	if model.IsProtectOnce != nil {
		modelMap["is_protect_once"] = *model.IsProtectOnce
	}
	if model.MissingEntities != nil {
		missingEntities := []map[string]interface{}{}
		for _, missingEntitiesItem := range model.MissingEntities {
			missingEntitiesItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMissingEntityParamsToMap(&missingEntitiesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			missingEntities = append(missingEntities, missingEntitiesItemMap)
		}
		modelMap["missing_entities"] = missingEntities
	}
	if model.InvalidEntities != nil {
		invalidEntities := []map[string]interface{}{}
		for _, invalidEntitiesItem := range model.InvalidEntities {
			invalidEntitiesItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMissingEntityParamsToMap(&invalidEntitiesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			invalidEntities = append(invalidEntities, invalidEntitiesItemMap)
		}
		modelMap["invalid_entities"] = invalidEntities
	}
	if model.NumProtectedObjects != nil {
		modelMap["num_protected_objects"] = flex.IntValue(model.NumProtectedObjects)
	}
	if model.AdvancedConfigs != nil {
		advancedConfigs := []map[string]interface{}{}
		for _, advancedConfigsItem := range model.AdvancedConfigs {
			advancedConfigsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsKeyValuePairToMap(&advancedConfigsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			advancedConfigs = append(advancedConfigs, advancedConfigsItemMap)
		}
		modelMap["advanced_configs"] = advancedConfigs
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPhysicalProtectionGroupParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	if model.MssqlParams != nil {
		mssqlParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLProtectionGroupParamsToMap(model.MssqlParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mssql_params"] = []map[string]interface{}{mssqlParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsTimeOfDayToMap(model *backuprecoveryv1.TimeOfDay) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hour"] = flex.IntValue(model.Hour)
	modelMap["minute"] = flex.IntValue(model.Minute)
	if model.TimeZone != nil {
		modelMap["time_zone"] = *model.TimeZone
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsProtectionGroupAlertingPolicyToMap(model *backuprecoveryv1.ProtectionGroupAlertingPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["backup_run_status"] = model.BackupRunStatus
	if model.AlertTargets != nil {
		alertTargets := []map[string]interface{}{}
		for _, alertTargetsItem := range model.AlertTargets {
			alertTargetsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAlertTargetToMap(&alertTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			alertTargets = append(alertTargets, alertTargetsItemMap)
		}
		modelMap["alert_targets"] = alertTargets
	}
	if model.RaiseObjectLevelFailureAlert != nil {
		modelMap["raise_object_level_failure_alert"] = *model.RaiseObjectLevelFailureAlert
	}
	if model.RaiseObjectLevelFailureAlertAfterLastAttempt != nil {
		modelMap["raise_object_level_failure_alert_after_last_attempt"] = *model.RaiseObjectLevelFailureAlertAfterLastAttempt
	}
	if model.RaiseObjectLevelFailureAlertAfterEachAttempt != nil {
		modelMap["raise_object_level_failure_alert_after_each_attempt"] = *model.RaiseObjectLevelFailureAlertAfterEachAttempt
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsAlertTargetToMap(model *backuprecoveryv1.AlertTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["email_address"] = *model.EmailAddress
	if model.Language != nil {
		modelMap["language"] = *model.Language
	}
	if model.RecipientType != nil {
		modelMap["recipient_type"] = *model.RecipientType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsSlaRuleToMap(model *backuprecoveryv1.SlaRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BackupRunType != nil {
		modelMap["backup_run_type"] = *model.BackupRunType
	}
	if model.SlaMinutes != nil {
		modelMap["sla_minutes"] = flex.IntValue(model.SlaMinutes)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsProtectionGroupRunToMap(model *backuprecoveryv1.ProtectionGroupRun) (map[string]interface{}, error) {
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
		originClusterIdentifierMap, err := DataSourceIbmBackupRecoveryProtectionGroupsClusterIdentifierToMap(model.OriginClusterIdentifier)
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
			objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsObjectRunResultToMap(&objectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			objects = append(objects, objectsItemMap)
		}
		modelMap["objects"] = objects
	}
	if model.LocalBackupInfo != nil {
		localBackupInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupRunSummaryToMap(model.LocalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_backup_info"] = []map[string]interface{}{localBackupInfoMap}
	}
	if model.OriginalBackupInfo != nil {
		originalBackupInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupRunSummaryToMap(model.OriginalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_backup_info"] = []map[string]interface{}{originalBackupInfoMap}
	}
	if model.ReplicationInfo != nil {
		replicationInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsReplicationRunSummaryToMap(model.ReplicationInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["replication_info"] = []map[string]interface{}{replicationInfoMap}
	}
	if model.ArchivalInfo != nil {
		archivalInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsArchivalRunSummaryToMap(model.ArchivalInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_info"] = []map[string]interface{}{archivalInfoMap}
	}
	if model.CloudSpinInfo != nil {
		cloudSpinInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinRunSummaryToMap(model.CloudSpinInfo)
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
			permissionsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsTenantToMap(&permissionsItem) // #nosec G601
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

func DataSourceIbmBackupRecoveryProtectionGroupsClusterIdentifierToMap(model *backuprecoveryv1.ClusterIdentifier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsObjectRunResultToMap(model *backuprecoveryv1.ObjectRunResult) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Object != nil {
		objectMap, err := DataSourceIbmBackupRecoveryProtectionGroupsObjectSummaryToMap(model.Object)
		if err != nil {
			return modelMap, err
		}
		modelMap["object"] = []map[string]interface{}{objectMap}
	}
	if model.LocalSnapshotInfo != nil {
		localSnapshotInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupRunToMap(model.LocalSnapshotInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_snapshot_info"] = []map[string]interface{}{localSnapshotInfoMap}
	}
	if model.OriginalBackupInfo != nil {
		originalBackupInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupRunToMap(model.OriginalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_backup_info"] = []map[string]interface{}{originalBackupInfoMap}
	}
	if model.ReplicationInfo != nil {
		replicationInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsReplicationRunToMap(model.ReplicationInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["replication_info"] = []map[string]interface{}{replicationInfoMap}
	}
	if model.ArchivalInfo != nil {
		archivalInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsArchivalRunToMap(model.ArchivalInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_info"] = []map[string]interface{}{archivalInfoMap}
	}
	if model.CloudSpinInfo != nil {
		cloudSpinInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinRunToMap(model.CloudSpinInfo)
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

func DataSourceIbmBackupRecoveryProtectionGroupsObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoveryProtectionGroupsSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoveryProtectionGroupsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoveryProtectionGroupsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsSharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsBackupRunToMap(model *backuprecoveryv1.BackupRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SnapshotInfo != nil {
		snapshotInfoMap, err := DataSourceIbmBackupRecoveryProtectionGroupsSnapshotInfoToMap(model.SnapshotInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["snapshot_info"] = []map[string]interface{}{snapshotInfoMap}
	}
	if model.FailedAttempts != nil {
		failedAttempts := []map[string]interface{}{}
		for _, failedAttemptsItem := range model.FailedAttempts {
			failedAttemptsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupAttemptToMap(&failedAttemptsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			failedAttempts = append(failedAttempts, failedAttemptsItemMap)
		}
		modelMap["failed_attempts"] = failedAttempts
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsSnapshotInfoToMap(model *backuprecoveryv1.SnapshotInfo) (map[string]interface{}, error) {
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
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupDataStatsToMap(model.Stats)
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
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsBackupDataStatsToMap(model *backuprecoveryv1.BackupDataStats) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsDataLockConstraintsToMap(model *backuprecoveryv1.DataLockConstraints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Mode != nil {
		modelMap["mode"] = *model.Mode
	}
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsBackupAttemptToMap(model *backuprecoveryv1.BackupAttempt) (map[string]interface{}, error) {
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
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupDataStatsToMap(model.Stats)
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

func DataSourceIbmBackupRecoveryProtectionGroupsReplicationRunToMap(model *backuprecoveryv1.ReplicationRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplicationTargetResults != nil {
		replicationTargetResults := []map[string]interface{}{}
		for _, replicationTargetResultsItem := range model.ReplicationTargetResults {
			replicationTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsReplicationTargetResultToMap(&replicationTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			replicationTargetResults = append(replicationTargetResults, replicationTargetResultsItemMap)
		}
		modelMap["replication_target_results"] = replicationTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsReplicationTargetResultToMap(model *backuprecoveryv1.ReplicationTargetResult) (map[string]interface{}, error) {
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
		awsTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAWSTargetConfigToMap(model.AwsTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_target_config"] = []map[string]interface{}{awsTargetConfigMap}
	}
	if model.AzureTargetConfig != nil {
		azureTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAzureTargetConfigToMap(model.AzureTargetConfig)
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
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsReplicationDataStatsToMap(model.Stats)
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
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsDataLockConstraintsToMap(model.DataLockConstraints)
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

func DataSourceIbmBackupRecoveryProtectionGroupsAWSTargetConfigToMap(model *backuprecoveryv1.AWSTargetConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsAzureTargetConfigToMap(model *backuprecoveryv1.AzureTargetConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsReplicationDataStatsToMap(model *backuprecoveryv1.ReplicationDataStats) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsArchivalRunToMap(model *backuprecoveryv1.ArchivalRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchivalTargetResults != nil {
		archivalTargetResults := []map[string]interface{}{}
		for _, archivalTargetResultsItem := range model.ArchivalTargetResults {
			archivalTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsArchivalTargetResultToMap(&archivalTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			archivalTargetResults = append(archivalTargetResults, archivalTargetResultsItemMap)
		}
		modelMap["archival_target_results"] = archivalTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsArchivalTargetResultToMap(model *backuprecoveryv1.ArchivalTargetResult) (map[string]interface{}, error) {
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
		tierSettingsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsArchivalTargetTierInfoToMap(model.TierSettings)
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
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsArchivalDataStatsToMap(model.Stats)
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
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	if model.OnLegalHold != nil {
		modelMap["on_legal_hold"] = *model.OnLegalHold
	}
	if model.WormProperties != nil {
		wormPropertiesMap, err := DataSourceIbmBackupRecoveryProtectionGroupsWormPropertiesToMap(model.WormProperties)
		if err != nil {
			return modelMap, err
		}
		modelMap["worm_properties"] = []map[string]interface{}{wormPropertiesMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsArchivalTargetTierInfoToMap(model *backuprecoveryv1.ArchivalTargetTierInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}
	if model.GoogleTiering != nil {
		googleTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupsGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := DataSourceIbmBackupRecoveryProtectionGroupsOracleTiersToMap(model.OracleTiering)
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

func DataSourceIbmBackupRecoveryProtectionGroupsAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsArchivalDataStatsToMap(model *backuprecoveryv1.ArchivalDataStats) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsWormPropertiesToMap(model *backuprecoveryv1.WormProperties) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinRunToMap(model *backuprecoveryv1.CloudSpinRun) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CloudSpinTargetResults != nil {
		cloudSpinTargetResults := []map[string]interface{}{}
		for _, cloudSpinTargetResultsItem := range model.CloudSpinTargetResults {
			cloudSpinTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinTargetResultToMap(&cloudSpinTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cloudSpinTargetResults = append(cloudSpinTargetResults, cloudSpinTargetResultsItemMap)
		}
		modelMap["cloud_spin_target_results"] = cloudSpinTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinTargetResultToMap(model *backuprecoveryv1.CloudSpinTargetResult) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsParams != nil {
		awsParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAwsCloudSpinParamsToMap(model.AwsParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_params"] = []map[string]interface{}{awsParamsMap}
	}
	if model.AzureParams != nil {
		azureParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAzureCloudSpinParamsToMap(model.AzureParams)
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
		statsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinDataStatsToMap(model.Stats)
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
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsDataLockConstraintsToMap(model.DataLockConstraints)
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

func DataSourceIbmBackupRecoveryProtectionGroupsAwsCloudSpinParamsToMap(model *backuprecoveryv1.AwsCloudSpinParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomTagList != nil {
		customTagList := []map[string]interface{}{}
		for _, customTagListItem := range model.CustomTagList {
			customTagListItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCustomTagParamsToMap(&customTagListItem) // #nosec G601
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

func DataSourceIbmBackupRecoveryProtectionGroupsCustomTagParamsToMap(model *backuprecoveryv1.CustomTagParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["key"] = *model.Key
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsAzureCloudSpinParamsToMap(model *backuprecoveryv1.AzureCloudSpinParams) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinDataStatsToMap(model *backuprecoveryv1.CloudSpinDataStats) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PhysicalBytesTransferred != nil {
		modelMap["physical_bytes_transferred"] = flex.IntValue(model.PhysicalBytesTransferred)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsBackupRunSummaryToMap(model *backuprecoveryv1.BackupRunSummary) (map[string]interface{}, error) {
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
		localSnapshotStatsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsBackupDataStatsToMap(model.LocalSnapshotStats)
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
		dataLockConstraintsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsDataLockConstraintsToMap(model.DataLockConstraints)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_constraints"] = []map[string]interface{}{dataLockConstraintsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsReplicationRunSummaryToMap(model *backuprecoveryv1.ReplicationRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplicationTargetResults != nil {
		replicationTargetResults := []map[string]interface{}{}
		for _, replicationTargetResultsItem := range model.ReplicationTargetResults {
			replicationTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsReplicationTargetResultToMap(&replicationTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			replicationTargetResults = append(replicationTargetResults, replicationTargetResultsItemMap)
		}
		modelMap["replication_target_results"] = replicationTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsArchivalRunSummaryToMap(model *backuprecoveryv1.ArchivalRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ArchivalTargetResults != nil {
		archivalTargetResults := []map[string]interface{}{}
		for _, archivalTargetResultsItem := range model.ArchivalTargetResults {
			archivalTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsArchivalTargetResultToMap(&archivalTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			archivalTargetResults = append(archivalTargetResults, archivalTargetResultsItemMap)
		}
		modelMap["archival_target_results"] = archivalTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinRunSummaryToMap(model *backuprecoveryv1.CloudSpinRunSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CloudSpinTargetResults != nil {
		cloudSpinTargetResults := []map[string]interface{}{}
		for _, cloudSpinTargetResultsItem := range model.CloudSpinTargetResults {
			cloudSpinTargetResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCloudSpinTargetResultToMap(&cloudSpinTargetResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cloudSpinTargetResults = append(cloudSpinTargetResults, cloudSpinTargetResultsItemMap)
		}
		modelMap["cloud_spin_target_results"] = cloudSpinTargetResults
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
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
		externalVendorMetadataMap, err := DataSourceIbmBackupRecoveryProtectionGroupsExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
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
		networkMap, err := DataSourceIbmBackupRecoveryProtectionGroupsTenantNetworkToMap(model.Network)
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

func DataSourceIbmBackupRecoveryProtectionGroupsExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
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
			customPropertiesItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
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
		metricsConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantMetricsConfigToMap(model.MetricsConfig)
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

func DataSourceIbmBackupRecoveryProtectionGroupsExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}
func DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantMetricsConfigToMap(model *backuprecoveryv1.IbmTenantMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CosResourceConfig != nil {
		cosResourceConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantCOSResourceConfigToMap(model.CosResourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_resource_config"] = []map[string]interface{}{cosResourceConfigMap}
	}
	if model.IamMetricsConfig != nil {
		iamMetricsConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantIAMMetricsConfigToMap(model.IamMetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["iam_metrics_config"] = []map[string]interface{}{iamMetricsConfigMap}
	}
	if model.MeteringConfig != nil {
		meteringConfigMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantMeteringConfigToMap(model.MeteringConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metering_config"] = []map[string]interface{}{meteringConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantCOSResourceConfigToMap(model *backuprecoveryv1.IbmTenantCOSResourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceURL != nil {
		modelMap["resource_url"] = *model.ResourceURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantIAMMetricsConfigToMap(model *backuprecoveryv1.IbmTenantIAMMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IAMURL != nil {
		modelMap["iam_url"] = *model.IAMURL
	}
	if model.BillingApiKeySecretID != nil {
		modelMap["billing_api_key_secret_id"] = *model.BillingApiKeySecretID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsIbmTenantMeteringConfigToMap(model *backuprecoveryv1.IbmTenantMeteringConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionGroupsMissingEntityParamsToMap(model *backuprecoveryv1.MissingEntityParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ParentSourceID != nil {
		modelMap["parent_source_id"] = flex.IntValue(model.ParentSourceID)
	}
	if model.ParentSourceName != nil {
		modelMap["parent_source_name"] = *model.ParentSourceName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsKeyValuePairToMap(model *backuprecoveryv1.KeyValuePair) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["key"] = *model.Key
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPhysicalProtectionGroupParamsToMap(model *backuprecoveryv1.PhysicalProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["protection_type"] = *model.ProtectionType
	if model.VolumeProtectionTypeParams != nil {
		volumeProtectionTypeParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPhysicalVolumeProtectionGroupParamsToMap(model.VolumeProtectionTypeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["volume_protection_type_params"] = []map[string]interface{}{volumeProtectionTypeParamsMap}
	}
	if model.FileProtectionTypeParams != nil {
		fileProtectionTypeParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPhysicalFileProtectionGroupParamsToMap(model.FileProtectionTypeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["file_protection_type_params"] = []map[string]interface{}{fileProtectionTypeParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPhysicalVolumeProtectionGroupParamsToMap(model *backuprecoveryv1.PhysicalVolumeProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPhysicalVolumeProtectionGroupObjectParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	if model.IndexingPolicy != nil {
		indexingPolicyMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIndexingPolicyToMap(model.IndexingPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["indexing_policy"] = []map[string]interface{}{indexingPolicyMap}
	}
	if model.PerformSourceSideDeduplication != nil {
		modelMap["perform_source_side_deduplication"] = *model.PerformSourceSideDeduplication
	}
	if model.Quiesce != nil {
		modelMap["quiesce"] = *model.Quiesce
	}
	if model.ContinueOnQuiesceFailure != nil {
		modelMap["continue_on_quiesce_failure"] = *model.ContinueOnQuiesceFailure
	}
	if model.IncrementalBackupAfterRestart != nil {
		modelMap["incremental_backup_after_restart"] = *model.IncrementalBackupAfterRestart
	}
	if model.PrePostScript != nil {
		prePostScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPrePostScriptParamsToMap(model.PrePostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["pre_post_script"] = []map[string]interface{}{prePostScriptMap}
	}
	if model.DedupExclusionSourceIds != nil {
		modelMap["dedup_exclusion_source_ids"] = model.DedupExclusionSourceIds
	}
	if model.ExcludedVssWriters != nil {
		modelMap["excluded_vss_writers"] = model.ExcludedVssWriters
	}
	if model.CobmrBackup != nil {
		modelMap["cobmr_backup"] = *model.CobmrBackup
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPhysicalVolumeProtectionGroupObjectParamsToMap(model *backuprecoveryv1.PhysicalVolumeProtectionGroupObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.VolumeGuids != nil {
		modelMap["volume_guids"] = model.VolumeGuids
	}
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	if model.ExcludedVssWriters != nil {
		modelMap["excluded_vss_writers"] = model.ExcludedVssWriters
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsIndexingPolicyToMap(model *backuprecoveryv1.IndexingPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["enable_indexing"] = *model.EnableIndexing
	if model.IncludePaths != nil {
		modelMap["include_paths"] = model.IncludePaths
	}
	if model.ExcludePaths != nil {
		modelMap["exclude_paths"] = model.ExcludePaths
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPrePostScriptParamsToMap(model *backuprecoveryv1.PrePostScriptParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PreScript != nil {
		preScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCommonPreBackupScriptParamsToMap(model.PreScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["pre_script"] = []map[string]interface{}{preScriptMap}
	}
	if model.PostScript != nil {
		postScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCommonPostBackupScriptParamsToMap(model.PostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["post_script"] = []map[string]interface{}{postScriptMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsCommonPreBackupScriptParamsToMap(model *backuprecoveryv1.CommonPreBackupScriptParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["path"] = *model.Path
	if model.Params != nil {
		modelMap["params"] = *model.Params
	}
	if model.TimeoutSecs != nil {
		modelMap["timeout_secs"] = flex.IntValue(model.TimeoutSecs)
	}
	if model.IsActive != nil {
		modelMap["is_active"] = *model.IsActive
	}
	if model.ContinueOnError != nil {
		modelMap["continue_on_error"] = *model.ContinueOnError
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsCommonPostBackupScriptParamsToMap(model *backuprecoveryv1.CommonPostBackupScriptParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["path"] = *model.Path
	if model.Params != nil {
		modelMap["params"] = *model.Params
	}
	if model.TimeoutSecs != nil {
		modelMap["timeout_secs"] = flex.IntValue(model.TimeoutSecs)
	}
	if model.IsActive != nil {
		modelMap["is_active"] = *model.IsActive
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPhysicalFileProtectionGroupParamsToMap(model *backuprecoveryv1.PhysicalFileProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExcludedVssWriters != nil {
		modelMap["excluded_vss_writers"] = model.ExcludedVssWriters
	}
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPhysicalFileProtectionGroupObjectParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	if model.IndexingPolicy != nil {
		indexingPolicyMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIndexingPolicyToMap(model.IndexingPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["indexing_policy"] = []map[string]interface{}{indexingPolicyMap}
	}
	if model.PerformSourceSideDeduplication != nil {
		modelMap["perform_source_side_deduplication"] = *model.PerformSourceSideDeduplication
	}
	if model.PerformBrickBasedDeduplication != nil {
		modelMap["perform_brick_based_deduplication"] = *model.PerformBrickBasedDeduplication
	}
	if model.TaskTimeouts != nil {
		taskTimeouts := []map[string]interface{}{}
		for _, taskTimeoutsItem := range model.TaskTimeouts {
			taskTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsCancellationTimeoutParamsToMap(&taskTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			taskTimeouts = append(taskTimeouts, taskTimeoutsItemMap)
		}
		modelMap["task_timeouts"] = taskTimeouts
	}
	if model.Quiesce != nil {
		modelMap["quiesce"] = *model.Quiesce
	}
	if model.ContinueOnQuiesceFailure != nil {
		modelMap["continue_on_quiesce_failure"] = *model.ContinueOnQuiesceFailure
	}
	if model.CobmrBackup != nil {
		modelMap["cobmr_backup"] = *model.CobmrBackup
	}
	if model.PrePostScript != nil {
		prePostScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPrePostScriptParamsToMap(model.PrePostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["pre_post_script"] = []map[string]interface{}{prePostScriptMap}
	}
	if model.DedupExclusionSourceIds != nil {
		modelMap["dedup_exclusion_source_ids"] = model.DedupExclusionSourceIds
	}
	if model.GlobalExcludePaths != nil {
		modelMap["global_exclude_paths"] = model.GlobalExcludePaths
	}
	if model.GlobalExcludeFS != nil {
		modelMap["global_exclude_fs"] = model.GlobalExcludeFS
	}
	if model.IgnorableErrors != nil {
		modelMap["ignorable_errors"] = model.IgnorableErrors
	}
	if model.AllowParallelRuns != nil {
		modelMap["allow_parallel_runs"] = *model.AllowParallelRuns
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPhysicalFileProtectionGroupObjectParamsToMap(model *backuprecoveryv1.PhysicalFileProtectionGroupObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExcludedVssWriters != nil {
		modelMap["excluded_vss_writers"] = model.ExcludedVssWriters
	}
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.FilePaths != nil {
		filePaths := []map[string]interface{}{}
		for _, filePathsItem := range model.FilePaths {
			filePathsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPhysicalFileBackupPathParamsToMap(&filePathsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			filePaths = append(filePaths, filePathsItemMap)
		}
		modelMap["file_paths"] = filePaths
	}
	if model.UsesPathLevelSkipNestedVolumeSetting != nil {
		modelMap["uses_path_level_skip_nested_volume_setting"] = *model.UsesPathLevelSkipNestedVolumeSetting
	}
	if model.NestedVolumeTypesToSkip != nil {
		modelMap["nested_volume_types_to_skip"] = model.NestedVolumeTypesToSkip
	}
	if model.FollowNasSymlinkTarget != nil {
		modelMap["follow_nas_symlink_target"] = *model.FollowNasSymlinkTarget
	}
	if model.MetadataFilePath != nil {
		modelMap["metadata_file_path"] = *model.MetadataFilePath
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsPhysicalFileBackupPathParamsToMap(model *backuprecoveryv1.PhysicalFileBackupPathParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["included_path"] = *model.IncludedPath
	if model.ExcludedPaths != nil {
		modelMap["excluded_paths"] = model.ExcludedPaths
	}
	if model.SkipNestedVolumes != nil {
		modelMap["skip_nested_volumes"] = *model.SkipNestedVolumes
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsCancellationTimeoutParamsToMap(model *backuprecoveryv1.CancellationTimeoutParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TimeoutMins != nil {
		modelMap["timeout_mins"] = flex.IntValue(model.TimeoutMins)
	}
	if model.BackupType != nil {
		modelMap["backup_type"] = *model.BackupType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLProtectionGroupParamsToMap(model *backuprecoveryv1.MSSQLProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FileProtectionTypeParams != nil {
		fileProtectionTypeParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLFileProtectionGroupParamsToMap(model.FileProtectionTypeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["file_protection_type_params"] = []map[string]interface{}{fileProtectionTypeParamsMap}
	}
	if model.NativeProtectionTypeParams != nil {
		nativeProtectionTypeParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLNativeProtectionGroupParamsToMap(model.NativeProtectionTypeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["native_protection_type_params"] = []map[string]interface{}{nativeProtectionTypeParamsMap}
	}
	modelMap["protection_type"] = *model.ProtectionType
	if model.VolumeProtectionTypeParams != nil {
		volumeProtectionTypeParamsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLVolumeProtectionGroupParamsToMap(model.VolumeProtectionTypeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["volume_protection_type_params"] = []map[string]interface{}{volumeProtectionTypeParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLFileProtectionGroupParamsToMap(model *backuprecoveryv1.MSSQLFileProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagBackupPreferenceType != nil {
		modelMap["aag_backup_preference_type"] = *model.AagBackupPreferenceType
	}
	if model.AdvancedSettings != nil {
		advancedSettingsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAdvancedSettingsToMap(model.AdvancedSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["advanced_settings"] = []map[string]interface{}{advancedSettingsMap}
	}
	if model.BackupSystemDbs != nil {
		modelMap["backup_system_dbs"] = *model.BackupSystemDbs
	}
	if model.ExcludeFilters != nil {
		excludeFilters := []map[string]interface{}{}
		for _, excludeFiltersItem := range model.ExcludeFilters {
			excludeFiltersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsFilterToMap(&excludeFiltersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			excludeFilters = append(excludeFilters, excludeFiltersItemMap)
		}
		modelMap["exclude_filters"] = excludeFilters
	}
	if model.FullBackupsCopyOnly != nil {
		modelMap["full_backups_copy_only"] = *model.FullBackupsCopyOnly
	}
	if model.LogBackupNumStreams != nil {
		modelMap["log_backup_num_streams"] = flex.IntValue(model.LogBackupNumStreams)
	}
	if model.LogBackupWithClause != nil {
		modelMap["log_backup_with_clause"] = *model.LogBackupWithClause
	}
	if model.PrePostScript != nil {
		prePostScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPrePostScriptParamsToMap(model.PrePostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["pre_post_script"] = []map[string]interface{}{prePostScriptMap}
	}
	if model.UseAagPreferencesFromServer != nil {
		modelMap["use_aag_preferences_from_server"] = *model.UseAagPreferencesFromServer
	}
	if model.UserDbBackupPreferenceType != nil {
		modelMap["user_db_backup_preference_type"] = *model.UserDbBackupPreferenceType
	}
	if model.AdditionalHostParams != nil {
		additionalHostParams := []map[string]interface{}{}
		for _, additionalHostParamsItem := range model.AdditionalHostParams {
			additionalHostParamsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLFileProtectionGroupHostParamsToMap(&additionalHostParamsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			additionalHostParams = append(additionalHostParams, additionalHostParamsItemMap)
		}
		modelMap["additional_host_params"] = additionalHostParams
	}
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLFileProtectionGroupObjectParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	if model.PerformSourceSideDeduplication != nil {
		modelMap["perform_source_side_deduplication"] = *model.PerformSourceSideDeduplication
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsAdvancedSettingsToMap(model *backuprecoveryv1.AdvancedSettings) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClonedDbBackupStatus != nil {
		modelMap["cloned_db_backup_status"] = *model.ClonedDbBackupStatus
	}
	if model.DbBackupIfNotOnlineStatus != nil {
		modelMap["db_backup_if_not_online_status"] = *model.DbBackupIfNotOnlineStatus
	}
	if model.MissingDbBackupStatus != nil {
		modelMap["missing_db_backup_status"] = *model.MissingDbBackupStatus
	}
	if model.OfflineRestoringDbBackupStatus != nil {
		modelMap["offline_restoring_db_backup_status"] = *model.OfflineRestoringDbBackupStatus
	}
	if model.ReadOnlyDbBackupStatus != nil {
		modelMap["read_only_db_backup_status"] = *model.ReadOnlyDbBackupStatus
	}
	if model.ReportAllNonAutoprotectDbErrors != nil {
		modelMap["report_all_non_autoprotect_db_errors"] = *model.ReportAllNonAutoprotectDbErrors
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsFilterToMap(model *backuprecoveryv1.Filter) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FilterString != nil {
		modelMap["filter_string"] = *model.FilterString
	}
	if model.IsRegularExpression != nil {
		modelMap["is_regular_expression"] = *model.IsRegularExpression
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLFileProtectionGroupHostParamsToMap(model *backuprecoveryv1.MSSQLFileProtectionGroupHostParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DisableSourceSideDeduplication != nil {
		modelMap["disable_source_side_deduplication"] = *model.DisableSourceSideDeduplication
	}
	modelMap["host_id"] = flex.IntValue(model.HostID)
	if model.HostName != nil {
		modelMap["host_name"] = *model.HostName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLFileProtectionGroupObjectParamsToMap(model *backuprecoveryv1.MSSQLFileProtectionGroupObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SourceType != nil {
		modelMap["source_type"] = *model.SourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLNativeProtectionGroupParamsToMap(model *backuprecoveryv1.MSSQLNativeProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagBackupPreferenceType != nil {
		modelMap["aag_backup_preference_type"] = *model.AagBackupPreferenceType
	}
	if model.AdvancedSettings != nil {
		advancedSettingsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAdvancedSettingsToMap(model.AdvancedSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["advanced_settings"] = []map[string]interface{}{advancedSettingsMap}
	}
	if model.BackupSystemDbs != nil {
		modelMap["backup_system_dbs"] = *model.BackupSystemDbs
	}
	if model.ExcludeFilters != nil {
		excludeFilters := []map[string]interface{}{}
		for _, excludeFiltersItem := range model.ExcludeFilters {
			excludeFiltersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsFilterToMap(&excludeFiltersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			excludeFilters = append(excludeFilters, excludeFiltersItemMap)
		}
		modelMap["exclude_filters"] = excludeFilters
	}
	if model.FullBackupsCopyOnly != nil {
		modelMap["full_backups_copy_only"] = *model.FullBackupsCopyOnly
	}
	if model.LogBackupNumStreams != nil {
		modelMap["log_backup_num_streams"] = flex.IntValue(model.LogBackupNumStreams)
	}
	if model.LogBackupWithClause != nil {
		modelMap["log_backup_with_clause"] = *model.LogBackupWithClause
	}
	if model.PrePostScript != nil {
		prePostScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPrePostScriptParamsToMap(model.PrePostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["pre_post_script"] = []map[string]interface{}{prePostScriptMap}
	}
	if model.UseAagPreferencesFromServer != nil {
		modelMap["use_aag_preferences_from_server"] = *model.UseAagPreferencesFromServer
	}
	if model.UserDbBackupPreferenceType != nil {
		modelMap["user_db_backup_preference_type"] = *model.UserDbBackupPreferenceType
	}
	if model.NumStreams != nil {
		modelMap["num_streams"] = flex.IntValue(model.NumStreams)
	}
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLNativeProtectionGroupObjectParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	if model.WithClause != nil {
		modelMap["with_clause"] = *model.WithClause
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLNativeProtectionGroupObjectParamsToMap(model *backuprecoveryv1.MSSQLNativeProtectionGroupObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SourceType != nil {
		modelMap["source_type"] = *model.SourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLVolumeProtectionGroupParamsToMap(model *backuprecoveryv1.MSSQLVolumeProtectionGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagBackupPreferenceType != nil {
		modelMap["aag_backup_preference_type"] = *model.AagBackupPreferenceType
	}
	if model.AdvancedSettings != nil {
		advancedSettingsMap, err := DataSourceIbmBackupRecoveryProtectionGroupsAdvancedSettingsToMap(model.AdvancedSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["advanced_settings"] = []map[string]interface{}{advancedSettingsMap}
	}
	if model.BackupSystemDbs != nil {
		modelMap["backup_system_dbs"] = *model.BackupSystemDbs
	}
	if model.ExcludeFilters != nil {
		excludeFilters := []map[string]interface{}{}
		for _, excludeFiltersItem := range model.ExcludeFilters {
			excludeFiltersItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsFilterToMap(&excludeFiltersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			excludeFilters = append(excludeFilters, excludeFiltersItemMap)
		}
		modelMap["exclude_filters"] = excludeFilters
	}
	if model.FullBackupsCopyOnly != nil {
		modelMap["full_backups_copy_only"] = *model.FullBackupsCopyOnly
	}
	if model.LogBackupNumStreams != nil {
		modelMap["log_backup_num_streams"] = flex.IntValue(model.LogBackupNumStreams)
	}
	if model.LogBackupWithClause != nil {
		modelMap["log_backup_with_clause"] = *model.LogBackupWithClause
	}
	if model.PrePostScript != nil {
		prePostScriptMap, err := DataSourceIbmBackupRecoveryProtectionGroupsPrePostScriptParamsToMap(model.PrePostScript)
		if err != nil {
			return modelMap, err
		}
		modelMap["pre_post_script"] = []map[string]interface{}{prePostScriptMap}
	}
	if model.UseAagPreferencesFromServer != nil {
		modelMap["use_aag_preferences_from_server"] = *model.UseAagPreferencesFromServer
	}
	if model.UserDbBackupPreferenceType != nil {
		modelMap["user_db_backup_preference_type"] = *model.UserDbBackupPreferenceType
	}
	if model.AdditionalHostParams != nil {
		additionalHostParams := []map[string]interface{}{}
		for _, additionalHostParamsItem := range model.AdditionalHostParams {
			additionalHostParamsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLVolumeProtectionGroupHostParamsToMap(&additionalHostParamsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			additionalHostParams = append(additionalHostParams, additionalHostParamsItemMap)
		}
		modelMap["additional_host_params"] = additionalHostParams
	}
	if model.BackupDbVolumesOnly != nil {
		modelMap["backup_db_volumes_only"] = *model.BackupDbVolumesOnly
	}
	if model.IncrementalBackupAfterRestart != nil {
		modelMap["incremental_backup_after_restart"] = *model.IncrementalBackupAfterRestart
	}
	if model.IndexingPolicy != nil {
		indexingPolicyMap, err := DataSourceIbmBackupRecoveryProtectionGroupsIndexingPolicyToMap(model.IndexingPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["indexing_policy"] = []map[string]interface{}{indexingPolicyMap}
	}
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := DataSourceIbmBackupRecoveryProtectionGroupsMSSQLVolumeProtectionGroupObjectParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLVolumeProtectionGroupHostParamsToMap(model *backuprecoveryv1.MSSQLVolumeProtectionGroupHostParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	modelMap["host_id"] = flex.IntValue(model.HostID)
	if model.HostName != nil {
		modelMap["host_name"] = *model.HostName
	}
	if model.VolumeGuids != nil {
		modelMap["volume_guids"] = model.VolumeGuids
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionGroupsMSSQLVolumeProtectionGroupObjectParamsToMap(model *backuprecoveryv1.MSSQLVolumeProtectionGroupObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SourceType != nil {
		modelMap["source_type"] = *model.SourceType
	}
	return modelMap, nil
}
