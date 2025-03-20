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
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryDownloadFilesFolders() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryDownloadFilesFoldersCreate,
		ReadContext:   resourceIbmBackupRecoveryDownloadFilesFoldersRead,
		DeleteContext: resourceIbmBackupRecoveryDownloadFilesFoldersDelete,
		UpdateContext: resourceIbmBackupRecoveryDownloadFilesFoldersUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecoveryDownloadFilesFolders,
		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"documents": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the list of documents to download using item ids. Only one of filesAndFolders or documents should be used. Currently only files are supported by documents.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_directory": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies whether the document is a directory. Since currently only files are supported this should always be false.",
						},
						"item_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the item id of the document.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the name of the recovery task. This field must be set and must be a unique name.",
			},
			"object": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the common snapshot parameters for a protected object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"snapshot_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the snapshot id.",
						},
						"point_in_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the timestamp (in microseconds. from epoch) for recovering to a point-in-time in the past.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the protection group id of the object snapshot.",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the protection group name of the object snapshot.",
						},
						"snapshot_creation_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the time when the snapshot is created in Unix timestamp epoch in microseconds.",
						},
						"object_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the information about the object for which the snapshot is taken.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies object id.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the object.",
									},
									"source_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies registered source id to which object belongs.",
									},
									"source_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies registered source name to which object belongs.",
									},
									"environment": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the environment of the object.",
									},
									"object_hash": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the hash identifier of the object.",
									},
									"object_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the type of the object.",
									},
									"logical_size_bytes": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the logical size of object in bytes.",
									},
									"uuid": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the uuid which is a unique identifier of the object.",
									},
									"global_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the global id which is a unique identifier of the object.",
									},
									"protection_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the protection type of the object if any.",
									},
									"sharepoint_site_summary": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the common parameters for Sharepoint site objects.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"site_web_url": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the web url for the Sharepoint site.",
												},
											},
										},
									},
									"os_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the operating system type of the object.",
									},
									"child_objects": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies child object details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies object id.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the name of the object.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies registered source id to which object belongs.",
												},
												"source_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies registered source name to which object belongs.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the environment of the object.",
												},
												"object_hash": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the hash identifier of the object.",
												},
												"object_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the type of the object.",
												},
												"logical_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the logical size of object in bytes.",
												},
												"uuid": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the uuid which is a unique identifier of the object.",
												},
												"global_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the global id which is a unique identifier of the object.",
												},
												"protection_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the protection type of the object if any.",
												},
												"sharepoint_site_summary": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the common parameters for Sharepoint site objects.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"site_web_url": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the web url for the Sharepoint site.",
															},
														},
													},
												},
												"os_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the operating system type of the object.",
												},
												"child_objects": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies child object details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{},
													},
												},
												"v_center_summary": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_cloud_env": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Computed:    true,
																Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
															},
														},
													},
												},
												"windows_cluster_summary": &schema.Schema{
													Type:     schema.TypeList,
													Optional: true,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_source_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the type of cluster resource this source represents.",
															},
														},
													},
												},
											},
										},
									},
									"v_center_summary": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_cloud_env": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
												},
											},
										},
									},
									"windows_cluster_summary": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_source_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the type of cluster resource this source represents.",
												},
											},
										},
									},
								},
							},
						},
						"snapshot_target_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the snapshot target type.",
						},
						"archival_target_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the archival target information if the snapshot is an archival snapshot.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the archival target ID.",
									},
									"archival_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the archival task id. This is a protection group UID which only applies when archival type is 'Tape'.",
									},
									"target_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the archival target name.",
									},
									"target_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the archival target type.",
									},
									"usage_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the usage type for the target.",
									},
									"ownership_context": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the ownership context for the target.",
									},
									"tier_settings": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the tier info for archival.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aws_tiering": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies aws tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Required:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies Azure tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies the cloud platform to enable tiering.",
												},
												"google_tiering": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies Google tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Required:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies Oracle tiers.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"tiers": &schema.Schema{
																Type:        schema.TypeList,
																Required:    true,
																Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"move_after_unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																		},
																		"move_after": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies the type of the current tier where the snapshot resides. This will be specified if the run is a CAD run.",
												},
											},
										},
									},
								},
							},
						},
						"progress_task_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Progress monitor task id for Recovery of VM.",
						},
						"recover_from_standby": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies that user wants to perform standby restore if it is enabled for this object.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of the Recovery. 'Running' indicates that the Recovery is still running. 'Canceled' indicates that the Recovery has been cancelled. 'Canceling' indicates that the Recovery is in the process of being cancelled. 'Failed' indicates that the Recovery has failed. 'Succeeded' indicates that the Recovery has finished successfully. 'SucceededWithWarning' indicates that the Recovery finished successfully, but there were some warning messages. 'Skipped' indicates that the Recovery task was skipped.",
						},
						"start_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the start time of the Recovery in Unix timestamp epoch in microseconds.",
						},
						"end_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the end time of the Recovery in Unix timestamp epoch in microseconds. This field will be populated only after Recovery is finished.",
						},
						"messages": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specify error messages about the object.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"bytes_restored": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specify the total bytes restored.",
						},
					},
				},
			},
			"parent_recovery_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_recovery_download_files_folders", "parent_recovery_id"),
				Description: "If current recovery is child task triggered through another parent recovery operation, then this field will specify the id of the parent recovery.",
			},
			"files_and_folders": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the list of files and folders to download.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"absolute_path": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the absolute path of the file or folder.",
						},
						"is_directory": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Specifies whether the file or folder object is a directory.",
						},
					},
				},
			},
			"glacier_retrieval_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_recovery_download_files_folders", "glacier_retrieval_type"),
				Description: "Specifies the glacier retrieval type when restoring or downloding files or folders from a Glacier-based cloud snapshot.",
			},
			"recovery_request_initiator_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_recovery", "request_initiator_type"),
				Description: "Specifies the type of request from UI, which is used for services like magneto to determine the priority of requests.",
			},
			"recovery_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ForceNew:    true,
				Description: "Specifies the name of the Recovery.",
			},
			"recovery_snapshot_environment": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				// ForceNew:     true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_recovery", "snapshot_environment"),
				Description: "Specifies the type of snapshot environment for which the Recovery was performed.",
			},
			"recovery_physical_params": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				// ForceNew:    true,
				Description: "Specifies the recovery options specific to Physical environment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"objects": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies the list of Recover Object parameters. For recovering files, specifies the object contains the file to recover.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshot_id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the snapshot id.",
									},
									"point_in_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the timestamp (in microseconds. from epoch) for recovering to a point-in-time in the past.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the protection group id of the object snapshot.",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the protection group name of the object snapshot.",
									},
									"snapshot_creation_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the time when the snapshot is created in Unix timestamp epoch in microseconds.",
									},
									"object_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the information about the object for which the snapshot is taken.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies object id.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the object.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies registered source id to which object belongs.",
												},
												"source_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies registered source name to which object belongs.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the environment of the object.",
												},
												"object_hash": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the hash identifier of the object.",
												},
												"object_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the type of the object.",
												},
												"logical_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the logical size of object in bytes.",
												},
												"uuid": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the uuid which is a unique identifier of the object.",
												},
												"global_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the global id which is a unique identifier of the object.",
												},
												"protection_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the protection type of the object if any.",
												},
												"sharepoint_site_summary": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the common parameters for Sharepoint site objects.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"site_web_url": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the web url for the Sharepoint site.",
															},
														},
													},
												},
												"os_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the operating system type of the object.",
												},
												"child_objects": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies child object details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies object id.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the name of the object.",
															},
															"source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies registered source id to which object belongs.",
															},
															"source_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies registered source name to which object belongs.",
															},
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the environment of the object.",
															},
															"object_hash": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the hash identifier of the object.",
															},
															"object_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the type of the object.",
															},
															"logical_size_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the logical size of object in bytes.",
															},
															"uuid": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the uuid which is a unique identifier of the object.",
															},
															"global_id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the global id which is a unique identifier of the object.",
															},
															"protection_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the protection type of the object if any.",
															},
															"sharepoint_site_summary": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies the common parameters for Sharepoint site objects.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"site_web_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the web url for the Sharepoint site.",
																		},
																	},
																},
															},
															"os_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the operating system type of the object.",
															},
															"child_objects": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Specifies child object details.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
															"v_center_summary": &schema.Schema{
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_cloud_env": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
																		},
																	},
																},
															},
															"windows_cluster_summary": &schema.Schema{
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cluster_source_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the type of cluster resource this source represents.",
																		},
																	},
																},
															},
														},
													},
												},
												"v_center_summary": &schema.Schema{
													Type:     schema.TypeList,
													MaxItems: 1,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_cloud_env": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
															},
														},
													},
												},
												"windows_cluster_summary": &schema.Schema{
													Type:     schema.TypeList,
													MaxItems: 1,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_source_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the type of cluster resource this source represents.",
															},
														},
													},
												},
											},
										},
									},
									"snapshot_target_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the snapshot target type.",
									},
									"storage_domain_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the ID of the Storage Domain where this snapshot is stored.",
									},
									"archival_target_info": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the archival target information if the snapshot is an archival snapshot.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the archival target ID.",
												},
												"archival_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the archival task id. This is a protection group UID which only applies when archival type is 'Tape'.",
												},
												"target_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the archival target name.",
												},
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the archival target type.",
												},
												"usage_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the usage type for the target.",
												},
												"ownership_context": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the ownership context for the target.",
												},
												"tier_settings": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the tier info for archival.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"aws_tiering": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies aws tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies Azure tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																Required:    true,
																Description: "Specifies the cloud platform to enable tiering.",
															},
															"google_tiering": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies Google tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies Oracle tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																Optional:    true,
																Description: "Specifies the type of the current tier where the snapshot resides. This will be specified if the run is a CAD run.",
															},
														},
													},
												},
											},
										},
									},
									"progress_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task id for Recovery of VM.",
									},
									"recover_from_standby": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies that user wants to perform standby restore if it is enabled for this object.",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the Recovery. 'Running' indicates that the Recovery is still running. 'Canceled' indicates that the Recovery has been cancelled. 'Canceling' indicates that the Recovery is in the process of being cancelled. 'Failed' indicates that the Recovery has failed. 'Succeeded' indicates that the Recovery has finished successfully. 'SucceededWithWarning' indicates that the Recovery finished successfully, but there were some warning messages. 'Skipped' indicates that the Recovery task was skipped.",
									},
									"start_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the start time of the Recovery in Unix timestamp epoch in microseconds.",
									},
									"end_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the end time of the Recovery in Unix timestamp epoch in microseconds. This field will be populated only after Recovery is finished.",
									},
									"messages": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specify error messages about the object.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"bytes_restored": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specify the total bytes restored.",
									},
								},
							},
						},
						"recovery_action": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the type of recover action to be performed.",
						},
						"recover_volume_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the parameters to recover Physical Volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
									"physical_target_params": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the params for recovering to a physical target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_target": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the target entity where the volumes are being mounted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the id of the object.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object.",
															},
														},
													},
												},
												"volume_mapping": &schema.Schema{
													Type:        schema.TypeList,
													Required:    true,
													Description: "Specifies the mapping from source volumes to destination volumes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source_volume_guid": &schema.Schema{
																Type:         schema.TypeString,
																ValidateFunc: validation.StringIsNotEmpty,
																Required:     true,
																Description:  "Specifies the guid of the source volume.",
															},
															"destination_volume_guid": &schema.Schema{
																Type:         schema.TypeString,
																ValidateFunc: validation.StringIsNotEmpty,
																Required:     true,
																Description:  "Specifies the guid of the destination volume.",
															},
														},
													},
												},
												"force_unmount_volume": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether volume would be dismounted first during LockVolume failure. If not specified, default is false.",
												},
												"vlan_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on Cohesity, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on Cohesity, then the partition hostname or VIPs will be used for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If this is set to true, then even if VLANs are configured on the system, the partition VIPs will be used for the Recovery.",
															},
															"interface_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Interface group to use for Recovery.",
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
						"mount_volume_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the parameters to mount Physical Volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
									"physical_target_params": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the params for recovering to a physical target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_to_original_target": &schema.Schema{
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Specifies whether to mount to the original target. If true, originalTargetConfig must be specified. If false, newTargetConfig must be specified.",
												},
												"original_target_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the configuration for mounting to the original target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"server_credentials": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies credentials to access the target server. This is required if the server is of Linux OS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the username to access target entity.",
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the password to access target entity.",
																		},
																	},
																},
															},
														},
													},
												},
												"new_target_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the configuration for mounting to a new target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mount_target": &schema.Schema{
																Type:        schema.TypeList,
																MinItems:    1,
																MaxItems:    1,
																Required:    true,
																Description: "Specifies the target entity to recover to.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the id of the object.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the name of the object.",
																		},
																		"parent_source_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the id of the parent source of the target.",
																		},
																		"parent_source_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the name of the parent source of the target.",
																		},
																	},
																},
															},
															"server_credentials": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies credentials to access the target server. This is required if the server is of Linux OS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the username to access target entity.",
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the password to access target entity.",
																		},
																	},
																},
															},
														},
													},
												},
												"read_only_mount": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to perform a read-only mount. Default is false.",
												},
												"volume_names": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the names of volumes that need to be mounted. If this is not specified then all volumes that are part of the source VM will be mounted on the target VM.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"mounted_volume_mapping": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the mapping of original volumes and mounted volumes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"original_volume": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the name of the original volume.",
															},
															"mounted_volume": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the name of the point where the volume is mounted.",
															},
															"file_system_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the type of the file system of the volume.",
															},
														},
													},
												},
												"vlan_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on Cohesity, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on Cohesity, then the partition hostname or VIPs will be used for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If this is set to true, then even if VLANs are configured on the system, the partition VIPs will be used for the Recovery.",
															},
															"interface_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Interface group to use for Recovery.",
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
						"recover_file_and_folder_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the parameters to perform a file and folder recovery.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"files_and_folders": &schema.Schema{
										Type:        schema.TypeList,
										Required:    true,
										Description: "Specifies the information about the files and folders to be recovered.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"absolute_path": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the absolute path to the file or folder.",
												},
												"destination_dir": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the destination directory where the file/directory was copied.",
												},
												"is_directory": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether this is a directory or not.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the recovery status for this file or folder.",
												},
												"messages": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specify error messages about the file during recovery.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"is_view_file_recovery": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specify if the recovery is of type view file/folder.",
												},
											},
										},
									},
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
									"physical_target_params": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the parameters to recover to a Physical target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"recover_target": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the target entity where the volumes are being mounted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the id of the object.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the object.",
															},
															"parent_source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the id of the parent source of the target.",
															},
															"parent_source_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the parent source of the target.",
															},
														},
													},
												},
												"restore_to_original_paths": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "If this is true, then files will be restored to original paths.",
												},
												"overwrite_existing": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to overwrite existing file/folder during recovery.",
												},
												"alternate_restore_directory": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the directory path where restore should happen if restore_to_original_paths is set to false.",
												},
												"preserve_attributes": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to preserve file/folder attributes during recovery.",
												},
												"preserve_timestamps": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to preserve the original time stamps.",
												},
												"preserve_acls": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Whether to preserve the ACLs of the original file.",
												},
												"continue_on_error": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to continue recovering other volumes if one of the volumes fails to recover. Default value is false.",
												},
												"save_success_files": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether to save success files or not. Default value is false.",
												},
												"vlan_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on Cohesity, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on Cohesity, then the partition hostname or VIPs will be used for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If this is set to true, then even if VLANs are configured on the system, the partition VIPs will be used for the Recovery.",
															},
															"interface_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Interface group to use for Recovery.",
															},
														},
													},
												},
												"restore_entity_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the restore type (restore everything or ACLs only) when restoring or downloading files or folders from a Physical file based or block based backup snapshot.",
												},
											},
										},
									},
								},
							},
						},
						"download_file_and_folder_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the parameters to download files and folders.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expiry_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the time upto which the download link is available.",
									},
									"files_and_folders": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the info about the files and folders to be recovered.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"absolute_path": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the absolute path to the file or folder.",
												},
												"destination_dir": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the destination directory where the file/directory was copied.",
												},
												"is_directory": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specifies whether this is a directory or not.",
												},
												"status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the recovery status for this file or folder.",
												},
												"messages": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specify error messages about the file during recovery.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"is_view_file_recovery": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Description: "Specify if the recovery is of type view file/folder.",
												},
											},
										},
									},
									"download_file_path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the path location to download the files and folders.",
									},
								},
							},
						},
						"system_recovery_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the parameters to perform a system recovery.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"full_nas_path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the path to the recovery view.",
									},
								},
							},
						},
					},
				},
			},
			"recovery_mssql_params": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				// ForceNew:    true,
				Description: "Specifies the recovery options specific to Sql environment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recover_app_params": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the parameters to recover Sql databases.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshot_id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the snapshot id.",
									},
									"point_in_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the timestamp (in microseconds. from epoch) for recovering to a point-in-time in the past.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the protection group id of the object snapshot.",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the protection group name of the object snapshot.",
									},
									"snapshot_creation_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the time when the snapshot is created in Unix timestamp epoch in microseconds.",
									},
									"object_info": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the information about the object for which the snapshot is taken.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies object id.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the object.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies registered source id to which object belongs.",
												},
												"source_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies registered source name to which object belongs.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the environment of the object.",
												},
												"object_hash": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the hash identifier of the object.",
												},
												"object_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the type of the object.",
												},
												"logical_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the logical size of object in bytes.",
												},
												"uuid": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the uuid which is a unique identifier of the object.",
												},
												"global_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the global id which is a unique identifier of the object.",
												},
												"protection_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the protection type of the object if any.",
												},
												"sharepoint_site_summary": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the common parameters for Sharepoint site objects.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"site_web_url": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the web url for the Sharepoint site.",
															},
														},
													},
												},
												"os_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the operating system type of the object.",
												},
												"child_objects": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies child object details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies object id.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the name of the object.",
															},
															"source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies registered source id to which object belongs.",
															},
															"source_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies registered source name to which object belongs.",
															},
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the environment of the object.",
															},
															"object_hash": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the hash identifier of the object.",
															},
															"object_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the type of the object.",
															},
															"logical_size_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the logical size of object in bytes.",
															},
															"uuid": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the uuid which is a unique identifier of the object.",
															},
															"global_id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the global id which is a unique identifier of the object.",
															},
															"protection_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the protection type of the object if any.",
															},
															"sharepoint_site_summary": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies the common parameters for Sharepoint site objects.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"site_web_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the web url for the Sharepoint site.",
																		},
																	},
																},
															},
															"os_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the operating system type of the object.",
															},
															"child_objects": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Specifies child object details.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
															"v_center_summary": &schema.Schema{
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_cloud_env": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
																		},
																	},
																},
															},
															"windows_cluster_summary": &schema.Schema{
																Type:     schema.TypeList,
																MaxItems: 1,
																Optional: true,
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cluster_source_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the type of cluster resource this source represents.",
																		},
																	},
																},
															},
														},
													},
												},
												"v_center_summary": &schema.Schema{
													Type:     schema.TypeList,
													MaxItems: 1,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_cloud_env": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
															},
														},
													},
												},
												"windows_cluster_summary": &schema.Schema{
													Type:     schema.TypeList,
													MaxItems: 1,
													Optional: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_source_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the type of cluster resource this source represents.",
															},
														},
													},
												},
											},
										},
									},
									"snapshot_target_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the snapshot target type.",
									},
									"storage_domain_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the ID of the Storage Domain where this snapshot is stored.",
									},
									"archival_target_info": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the archival target information if the snapshot is an archival snapshot.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the archival target ID.",
												},
												"archival_task_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the archival task id. This is a protection group UID which only applies when archival type is 'Tape'.",
												},
												"target_name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the archival target name.",
												},
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the archival target type.",
												},
												"usage_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the usage type for the target.",
												},
												"ownership_context": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the ownership context for the target.",
												},
												"tier_settings": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the tier info for archival.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"aws_tiering": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies aws tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies Azure tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																Required:    true,
																Description: "Specifies the cloud platform to enable tiering.",
															},
															"google_tiering": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies Google tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies Oracle tiers.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"tiers": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies the tiers that are used to move the archived backup from current tier to next tier. The order of the tiers determines which tier will be used next for moving the archived backup. The first tier input should always be default tier where backup will be acrhived. Each tier specifies how much time after the backup will be moved to next tier from the current tier.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"move_after_unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Optional:    true,
																						Computed:    true,
																						Description: "Specifies the unit for moving the data from current tier to next tier. This unit will be a base unit for the 'moveAfter' field specified below.",
																					},
																					"move_after": &schema.Schema{
																						Type:        schema.TypeInt,
																						Optional:    true,
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
																Optional:    true,
																Description: "Specifies the type of the current tier where the snapshot resides. This will be specified if the run is a CAD run.",
															},
														},
													},
												},
											},
										},
									},
									"progress_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task id for Recovery of VM.",
									},
									"recover_from_standby": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies that user wants to perform standby restore if it is enabled for this object.",
									},
									"status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the Recovery. 'Running' indicates that the Recovery is still running. 'Canceled' indicates that the Recovery has been cancelled. 'Canceling' indicates that the Recovery is in the process of being cancelled. 'Failed' indicates that the Recovery has failed. 'Succeeded' indicates that the Recovery has finished successfully. 'SucceededWithWarning' indicates that the Recovery finished successfully, but there were some warning messages. 'Skipped' indicates that the Recovery task was skipped.",
									},
									"start_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the start time of the Recovery in Unix timestamp epoch in microseconds.",
									},
									"end_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the end time of the Recovery in Unix timestamp epoch in microseconds. This field will be populated only after Recovery is finished.",
									},
									"messages": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specify error messages about the object.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"bytes_restored": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specify the total bytes restored.",
									},
									"aag_info": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Object details for Mssql.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the AAG name.",
												},
												"object_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the AAG object Id.",
												},
											},
										},
									},
									"host_info": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the host information for a objects. This is mainly populated in case of App objects where app object is hosted by another object such as VM or physical server.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the id of the host object.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the name of the host object.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the environment of the object.",
												},
											},
										},
									},
									"is_encrypted": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Specifies whether the database is TDE enabled.",
									},
									"sql_target_params": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the params for recovering to a sql host. Specifiy seperate settings for each db object that need to be recovered. Provided sql backup should be recovered to same type of target host. For Example: If you have sql backup taken from a physical host then that should be recovered to physical host only.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"new_source_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the destination Source configuration parameters where the databases will be recovered. This is mandatory if recoverToNewSource is set to true.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keep_cdc": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies whether to keep CDC (Change Data Capture) on recovered databases or not. If not passed, this is assumed to be true. If withNoRecovery is passed as true, then this field must not be set to true. Passing this field as true in this scenario will be a invalid request.",
															},
															"multi_stage_restore_options": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies the parameters related to multi stage Sql restore.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enable_auto_sync": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Set this to true if you want to enable auto sync for multi stage restore.",
																		},
																		"enable_multi_stage_restore": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Set this to true if you are creating a multi-stage Sql restore task needed for features such as Hot-Standby.",
																		},
																	},
																},
															},
															"native_log_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the WITH clause to be used in native sql log restore command. This is only applicable for native log restore.",
															},
															"native_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "'with_clause' contains 'with clause' to be used in native sql restore command. This is only applicable for database restore of native sql backup. Here user can specify multiple restore options. Example: 'WITH BUFFERCOUNT = 575, MAXTRANSFERSIZE = 2097152'.",
															},
															"overwriting_policy": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies a policy to be used while recovering existing databases.",
															},
															"replay_entire_last_log": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies the option to set replay last log bit while creating the sql restore task and doing restore to latest point-in-time. If this is set to true, we will replay the entire last log without STOPAT.",
															},
															"restore_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the time in the past to which the Sql database needs to be restored. This allows for granular recovery of Sql databases. If this is not set, the Sql database will be restored from the full/incremental snapshot.",
															},
															"secondary_data_files_dir_list": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Specifies the secondary data filename pattern and corresponding direcories of the DB. Secondary data files are optional and are user defined. The recommended file extention for secondary files is \".ndf\". If this option is specified and the destination folders do not exist they will be automatically created.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"directory": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the directory where to keep the files matching the pattern.",
																		},
																		"filename_pattern": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies a pattern to be matched with filenames. This can be a regex expression.",
																		},
																	},
																},
															},
															"with_no_recovery": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies the flag to bring DBs online or not after successful recovery. If this is passed as true, then it means DBs won't be brought online.",
															},
															"data_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the directory where to put the database data files. Missing directory will be automatically created.",
															},
															"database_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies a new name for the restored database. If this field is not specified, then the original database will be overwritten after recovery.",
															},
															"host": &schema.Schema{
																Type:        schema.TypeList,
																MinItems:    1,
																MaxItems:    1,
																Required:    true,
																Description: "Specifies the source id of target host where databases will be recovered. This source id can be a physical host or virtual machine.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the id of the object.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the name of the object.",
																		},
																	},
																},
															},
															"instance_name": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies an instance name of the Sql Server that should be used for restoring databases to.",
															},
															"log_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the directory where to put the database log files. Missing directory will be automatically created.",
															},
														},
													},
												},
												"original_source_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the Source configuration if databases are being recovered to Original Source. If not specified, all the configuration parameters will be retained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keep_cdc": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies whether to keep CDC (Change Data Capture) on recovered databases or not. If not passed, this is assumed to be true. If withNoRecovery is passed as true, then this field must not be set to true. Passing this field as true in this scenario will be a invalid request.",
															},
															"multi_stage_restore_options": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies the parameters related to multi stage Sql restore.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enable_auto_sync": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Set this to true if you want to enable auto sync for multi stage restore.",
																		},
																		"enable_multi_stage_restore": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Set this to true if you are creating a multi-stage Sql restore task needed for features such as Hot-Standby.",
																		},
																	},
																},
															},
															"native_log_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the WITH clause to be used in native sql log restore command. This is only applicable for native log restore.",
															},
															"native_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "'with_clause' contains 'with clause' to be used in native sql restore command. This is only applicable for database restore of native sql backup. Here user can specify multiple restore options. Example: 'WITH BUFFERCOUNT = 575, MAXTRANSFERSIZE = 2097152'.",
															},
															"overwriting_policy": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies a policy to be used while recovering existing databases.",
															},
															"replay_entire_last_log": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies the option to set replay last log bit while creating the sql restore task and doing restore to latest point-in-time. If this is set to true, we will replay the entire last log without STOPAT.",
															},
															"restore_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the time in the past to which the Sql database needs to be restored. This allows for granular recovery of Sql databases. If this is not set, the Sql database will be restored from the full/incremental snapshot.",
															},
															"secondary_data_files_dir_list": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Specifies the secondary data filename pattern and corresponding direcories of the DB. Secondary data files are optional and are user defined. The recommended file extention for secondary files is \".ndf\". If this option is specified and the destination folders do not exist they will be automatically created.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"directory": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the directory where to keep the files matching the pattern.",
																		},
																		"filename_pattern": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies a pattern to be matched with filenames. This can be a regex expression.",
																		},
																	},
																},
															},
															"with_no_recovery": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Specifies the flag to bring DBs online or not after successful recovery. If this is passed as true, then it means DBs won't be brought online.",
															},
															"capture_tail_logs": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Set this to true if tail logs are to be captured before the recovery operation. This is only applicable if database is not being renamed.",
															},
															"data_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the directory where to put the database data files. Missing directory will be automatically created. If you are overwriting the existing database then this field will be ignored.",
															},
															"log_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the directory where to put the database log files. Missing directory will be automatically created. If you are overwriting the existing database then this field will be ignored.",
															},
															"new_database_name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies a new name for the restored database. If this field is not specified, then the original database will be overwritten after recovery.",
															},
														},
													},
												},
												"recover_to_new_source": &schema.Schema{
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Specifies the parameter whether the recovery should be performed to a new sources or an original Source Target.",
												},
											},
										},
									},
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
								},
							},
						},
						"recovery_action": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the type of recover action to be performed.",
						},
						"vlan_config": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on IBM, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on IBM, then the partition hostname or VIPs will be used for Recovery.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
									},
									"disable_vlan": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If this is set to true, then even if VLANs are configured on the system, the partition VIPs will be used for the Recovery.",
									},
									"interface_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Interface group to use for Recovery.",
									},
								},
							},
						},
					},
				},
			},
			"recovery_start_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the start time of the Recovery in Unix timestamp epoch in microseconds.",
			},
			"recovery_end_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the end time of the Recovery in Unix timestamp epoch in microseconds. This field will be populated only after Recovery is finished.",
			},
			"recovery_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Recovery. 'Running' indicates that the Recovery is still running. 'Canceled' indicates that the Recovery has been cancelled. 'Canceling' indicates that the Recovery is in the process of being cancelled. 'Failed' indicates that the Recovery has failed. 'Succeeded' indicates that the Recovery has finished successfully. 'SucceededWithWarning' indicates that the Recovery finished successfully, but there were some warning messages. 'Skipped' indicates that the Recovery task was skipped.",
			},
			"recovery_progress_task_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Progress monitor task id for Recovery.",
			},
			"recovery_action": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of recover action.",
			},
			"recovery_permissions": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of tenants that have permissions for this recovery.",
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
							Optional:    true,
							Computed:    true,
							Description: "Description about the tenant.",
						},
						"external_vendor_metadata": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the additional metadata for the tenant that is specifically set by the external vendors who are responsible for managing tenants. This field will only applicable if tenant creation is happening for a specially provisioned clusters for external vendors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ibm_tenant_metadata_params": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the additional metadata for the tenant that is specifically set by the external vendor of type 'IBM'.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"account_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the unique identifier of the IBM's account ID.",
												},
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the unique CRN associated with the tenant.",
												},
												"custom_properties": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the list of custom properties associated with the tenant. External vendors can choose to set any properties inside following list. Note that the fields set inside the following will not be available for direct filtering. API callers should make sure that no sensitive information such as passwords is sent in these fields.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"key": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the unique key for custom property.",
															},
															"value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the value for the above custom key.",
															},
														},
													},
												},
												"liveness_mode": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the current liveness mode of the tenant. This mode may change based on AZ failures when vendor chooses to failover or failback the tenants to other AZs.",
												},
												"ownership_mode": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the current ownership mode for the tenant. The ownership of the tenant represents the active role for functioning of the tenant.",
												},
												"resource_group_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the Resource Group ID associated with the tenant.",
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
							Optional:    true,
							Computed:    true,
							Description: "The tenant id.",
						},
						"is_managed_on_helios": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "Name of the Tenant.",
						},
						"network": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "The hostname for Cohesity cluster as seen by tenants and as is routable from the tenant's network. Tenant's VLAN's hostname, if available can be used instead but it is mandatory to provide this value if there's no VLAN hostname to use. Also, when set, this field would take precedence over VLAN hostname.",
									},
									"cluster_ips": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Set of IPs as seen from the tenant's network for the Cohesity cluster. Only one from 'clusterHostname' and 'clusterIps' is needed.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Current Status of the Tenant.",
						},
					},
				},
			},
			"recovery_creation_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the information about the creation of the protection group or recovery.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the user who created the protection group or recovery.",
						},
					},
				},
			},
			"recovery_can_tear_down": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether it's possible to tear down the objects created by the recovery.",
			},
			"recovery_tear_down_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the status of the tear down operation. This is only set when the canTearDown is set to true. 'DestroyScheduled' indicates that the tear down is ready to schedule. 'Destroying' indicates that the tear down is still running. 'Destroyed' indicates that the tear down succeeded. 'DestroyError' indicates that the tear down failed.",
			},
			"recovery_tear_down_message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the error message about the tear down operation if it fails.",
			},
			"recovery_messages": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies messages about the recovery.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"is_parent_recovery": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the current recovery operation has created child recoveries. This is currently used in SQL recovery where multiple child recoveries can be tracked under a common/parent recovery.",
			},
			"recovery_retrieve_archive_tasks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of persistent state of a retrieve of an archive task.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_uid": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the globally unique id for this retrieval of an archive task.",
						},
						"uptier_expiry_times": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies how much time the retrieved entity is present in the hot-tiers.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"recovery_is_multi_stage_restore": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the current recovery operation is a multi-stage restore operation. This is currently used by VMware recoveres for the migration/hot-standby use case.",
			},
		},
	}
}

func checkDiffResourceIbmBackupRecoveryDownloadFilesFolders(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecoveryDownloadFilesFolders().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_recovery_download_files_folders cannot be updated.")
		}
	}
	return nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "parent_recovery_id",
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^\d+:\d+:\d+$`,
		},
		validate.ValidateSchema{
			Identifier:                 "glacier_retrieval_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "kExpeditedNoPCU, kExpeditedWithPCU, kStandard",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_backup_recovery_recovery_download_files_folders", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmBackupRecoveryDownloadFilesFoldersCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery_download_files_folders", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createDownloadFilesAndFoldersRecoveryOptions := &backuprecoveryv1.CreateDownloadFilesAndFoldersRecoveryOptions{}

	createDownloadFilesAndFoldersRecoveryOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	createDownloadFilesAndFoldersRecoveryOptions.SetName(d.Get("name").(string))
	objectModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToCommonRecoverObjectSnapshotParams(d.Get("object.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery_download_files_folders", "create", "parse-object").GetDiag()
	}
	createDownloadFilesAndFoldersRecoveryOptions.SetObject(objectModel)
	var filesAndFolders []backuprecoveryv1.FilesAndFoldersObject
	for _, v := range d.Get("files_and_folders").([]interface{}) {
		value := v.(map[string]interface{})
		filesAndFoldersItem, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToFilesAndFoldersObject(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery_download_files_folders", "create", "parse-files_and_folders").GetDiag()
		}
		filesAndFolders = append(filesAndFolders, *filesAndFoldersItem)
	}
	createDownloadFilesAndFoldersRecoveryOptions.SetFilesAndFolders(filesAndFolders)
	if _, ok := d.GetOk("documents"); ok {
		var documents []backuprecoveryv1.DocumentObject
		for _, v := range d.Get("documents").([]interface{}) {
			value := v.(map[string]interface{})
			documentsItem, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToDocumentObject(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery_download_files_folders", "create", "parse-documents").GetDiag()
			}
			documents = append(documents, *documentsItem)
		}
		createDownloadFilesAndFoldersRecoveryOptions.SetDocuments(documents)
	}
	if _, ok := d.GetOk("parent_recovery_id"); ok {
		createDownloadFilesAndFoldersRecoveryOptions.SetParentRecoveryID(d.Get("parent_recovery_id").(string))
	}
	if _, ok := d.GetOk("glacier_retrieval_type"); ok {
		createDownloadFilesAndFoldersRecoveryOptions.SetGlacierRetrievalType(d.Get("glacier_retrieval_type").(string))
	}

	recovery, _, err := backupRecoveryClient.CreateDownloadFilesAndFoldersRecoveryWithContext(context, createDownloadFilesAndFoldersRecoveryOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDownloadFilesAndFoldersRecoveryWithContext failed: %s", err.Error()), "ibm_backup_recovery_recovery_download_files_folders", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*recovery.ID)

	return resourceIbmBackupRecoveryDownloadFilesFoldersRead(context, d, meta)
}

func resourceIbmBackupRecoveryDownloadFilesFoldersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery_download_files_folders", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getRecoveryByIdOptions := &backuprecoveryv1.GetRecoveryByIdOptions{}

	getRecoveryByIdOptions.SetID(d.Id())
	getRecoveryByIdOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))

	recovery, response, err := backupRecoveryClient.GetRecoveryByIDWithContext(context, getRecoveryByIdOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRecoveryByIDWithContext failed: %s", err.Error()), "ibm_backup_recovery_recovery_download_files_folders", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getRecoveryByIdOptions.ID)

	if !core.IsNil(recovery.Name) {
		if err = d.Set("recovery_name", recovery.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_recovery", "read", "set-name").GetDiag()
		}
	}

	if !core.IsNil(recovery.StartTimeUsecs) {
		if err = d.Set("recovery_start_time_usecs", flex.IntValue(recovery.StartTimeUsecs)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting start_time_usecs: %s", err), "(Data) ibm_recovery", "read", "set-start_time_usecs").GetDiag()
		}
	}

	if !core.IsNil(recovery.EndTimeUsecs) {
		if err = d.Set("recovery_end_time_usecs", flex.IntValue(recovery.EndTimeUsecs)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting end_time_usecs: %s", err), "(Data) ibm_recovery", "read", "set-end_time_usecs").GetDiag()
		}
	}

	if !core.IsNil(recovery.Status) {
		if err = d.Set("recovery_status", recovery.Status); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_recovery", "read", "set-status").GetDiag()
		}
	}

	if !core.IsNil(recovery.ProgressTaskID) {
		if err = d.Set("recovery_progress_task_id", recovery.ProgressTaskID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting progress_task_id: %s", err), "(Data) ibm_recovery", "read", "set-progress_task_id").GetDiag()
		}
	}

	if !core.IsNil(recovery.SnapshotEnvironment) {
		if err = d.Set("recovery_snapshot_environment", recovery.SnapshotEnvironment); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting snapshot_environment: %s", err), "(Data) ibm_recovery", "read", "set-snapshot_environment").GetDiag()
		}
	}

	if !core.IsNil(recovery.RecoveryAction) {
		if err = d.Set("recovery_action", recovery.RecoveryAction); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting recovery_action: %s", err), "(Data) ibm_recovery", "read", "set-recovery_action").GetDiag()
		}
	}

	if !core.IsNil(recovery.Permissions) {
		permissions := []map[string]interface{}{}
		for _, permissionsItem := range recovery.Permissions {
			permissionsItemMap, err := DataSourceIbmBackupRecoveryTenantToMap(&permissionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_recovery", "read", "permissions-to-map").GetDiag()
			}
			permissions = append(permissions, permissionsItemMap)
		}
		if err = d.Set("recovery_permissions", permissions); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting permissions: %s", err), "(Data) ibm_recovery", "read", "set-permissions").GetDiag()
		}
	}

	if !core.IsNil(recovery.CreationInfo) {
		creationInfo := []map[string]interface{}{}
		creationInfoMap, err := DataSourceIbmBackupRecoveryCreationInfoToMap(recovery.CreationInfo)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_recovery", "read", "creation_info-to-map").GetDiag()
		}
		creationInfo = append(creationInfo, creationInfoMap)
		if err = d.Set("recovery_creation_info", creationInfo); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting creation_info: %s", err), "(Data) ibm_recovery", "read", "set-creation_info").GetDiag()
		}
	}

	if !core.IsNil(recovery.CanTearDown) {
		if err = d.Set("recovery_can_tear_down", recovery.CanTearDown); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting can_tear_down: %s", err), "(Data) ibm_recovery", "read", "set-can_tear_down").GetDiag()
		}
	}

	if !core.IsNil(recovery.TearDownStatus) {
		if err = d.Set("recovery_tear_down_status", recovery.TearDownStatus); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tear_down_status: %s", err), "(Data) ibm_recovery", "read", "set-tear_down_status").GetDiag()
		}
	}

	if !core.IsNil(recovery.TearDownMessage) {
		if err = d.Set("recovery_tear_down_message", recovery.TearDownMessage); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tear_down_message: %s", err), "(Data) ibm_recovery", "read", "set-tear_down_message").GetDiag()
		}
	}

	if !core.IsNil(recovery.Messages) {
		messages := []interface{}{}
		for _, messagesItem := range recovery.Messages {
			messages = append(messages, messagesItem)
		}
		if err = d.Set("recovery_messages", messages); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting messages: %s", err), "(Data) ibm_recovery", "read", "set-messages").GetDiag()
		}
	} else {
		if err = d.Set("recovery_messages", []interface{}{}); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting messages: %s", err), "(Data) ibm_recovery", "read", "set-messages").GetDiag()
		}
	}

	if !core.IsNil(recovery.IsParentRecovery) {
		if err = d.Set("is_parent_recovery", recovery.IsParentRecovery); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting is_parent_recovery: %s", err), "(Data) ibm_recovery", "read", "set-is_parent_recovery").GetDiag()
		}
	}

	if !core.IsNil(recovery.ParentRecoveryID) {
		if err = d.Set("parent_recovery_id", recovery.ParentRecoveryID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting parent_recovery_id: %s", err), "(Data) ibm_recovery", "read", "set-parent_recovery_id").GetDiag()
		}
	}

	if !core.IsNil(recovery.RetrieveArchiveTasks) {
		retrieveArchiveTasks := []map[string]interface{}{}
		for _, retrieveArchiveTasksItem := range recovery.RetrieveArchiveTasks {
			retrieveArchiveTasksItemMap, err := DataSourceIbmBackupRecoveryRetrieveArchiveTaskToMap(&retrieveArchiveTasksItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_recovery", "read", "retrieve_archive_tasks-to-map").GetDiag()
			}
			retrieveArchiveTasks = append(retrieveArchiveTasks, retrieveArchiveTasksItemMap)
		}
		if err = d.Set("recovery_retrieve_archive_tasks", retrieveArchiveTasks); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting retrieve_archive_tasks: %s", err), "(Data) ibm_recovery", "read", "set-retrieve_archive_tasks").GetDiag()
		}
	} else {
		if err = d.Set("recovery_retrieve_archive_tasks", []interface{}{}); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting mssql_params: %s", err), "(Data) ibm_recovery", "read", "set-mssql_params").GetDiag()
		}
	}

	if !core.IsNil(recovery.IsMultiStageRestore) {
		if err = d.Set("recovery_is_multi_stage_restore", recovery.IsMultiStageRestore); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting is_multi_stage_restore: %s", err), "(Data) ibm_recovery", "read", "set-is_multi_stage_restore").GetDiag()
		}
	}

	if !core.IsNil(recovery.PhysicalParams) {
		physicalParams := []map[string]interface{}{}
		physicalParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsToMap(recovery.PhysicalParams)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_recovery", "read", "physical_params-to-map").GetDiag()
		}
		physicalParams = append(physicalParams, physicalParamsMap)
		if err = d.Set("recovery_physical_params", physicalParams); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting physical_params: %s", err), "(Data) ibm_recovery", "read", "set-physical_params").GetDiag()
		}
	}

	if !core.IsNil(recovery.MssqlParams) {
		mssqlParams := []map[string]interface{}{}
		mssqlParamsMap, err := DataSourceIbmBackupRecoveryRecoverSqlParamsToMap(recovery.MssqlParams)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_recovery", "read", "mssql_params-to-map").GetDiag()
		}
		mssqlParams = append(mssqlParams, mssqlParamsMap)
		if err = d.Set("recovery_mssql_params", mssqlParams); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting mssql_params: %s", err), "(Data) ibm_recovery", "read", "set-mssql_params").GetDiag()
		}
	} else {
		if err = d.Set("recovery_mssql_params", []interface{}{}); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting mssql_params: %s", err), "(Data) ibm_recovery", "read", "set-mssql_params").GetDiag()
		}
	}

	return nil
}

func resourceIbmBackupRecoveryDownloadFilesFoldersDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmBackupRecoveryDownloadFilesFoldersUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This resource does not support a "update" operation.
	var diags diag.Diagnostics
	warning := diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource update will only affect terraform state and not the actual backend resource",
		Detail:   "Update operation for this resource is not supported and will only affect the terraform statefile. No changes will be made to the backend resource.",
	}
	diags = append(diags, warning)
	return diags
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToCommonRecoverObjectSnapshotParams(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverObjectSnapshotParams, error) {
	model := &backuprecoveryv1.CommonRecoverObjectSnapshotParams{}
	model.SnapshotID = core.StringPtr(modelMap["snapshot_id"].(string))
	if modelMap["point_in_time_usecs"] != nil {
		model.PointInTimeUsecs = core.Int64Ptr(int64(modelMap["point_in_time_usecs"].(int)))
	}
	if modelMap["protection_group_id"] != nil && modelMap["protection_group_id"].(string) != "" {
		model.ProtectionGroupID = core.StringPtr(modelMap["protection_group_id"].(string))
	}
	if modelMap["protection_group_name"] != nil && modelMap["protection_group_name"].(string) != "" {
		model.ProtectionGroupName = core.StringPtr(modelMap["protection_group_name"].(string))
	}
	if modelMap["snapshot_creation_time_usecs"] != nil && modelMap["snapshot_creation_time_usecs"].(int) != 0 {
		model.SnapshotCreationTimeUsecs = core.Int64Ptr(int64(modelMap["snapshot_creation_time_usecs"].(int)))
	}
	if modelMap["object_info"] != nil && len(modelMap["object_info"].([]interface{})) > 0 {
		ObjectInfoModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToCommonRecoverObjectSnapshotParamsObjectInfo(modelMap["object_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ObjectInfo = ObjectInfoModel
	}
	if modelMap["snapshot_target_type"] != nil && modelMap["snapshot_target_type"].(string) != "" {
		model.SnapshotTargetType = core.StringPtr(modelMap["snapshot_target_type"].(string))
	}
	if modelMap["archival_target_info"] != nil && len(modelMap["archival_target_info"].([]interface{})) > 0 {
		ArchivalTargetInfoModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToCommonRecoverObjectSnapshotParamsArchivalTargetInfo(modelMap["archival_target_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ArchivalTargetInfo = ArchivalTargetInfoModel
	}
	if modelMap["progress_task_id"] != nil && modelMap["progress_task_id"].(string) != "" {
		model.ProgressTaskID = core.StringPtr(modelMap["progress_task_id"].(string))
	}
	if modelMap["recover_from_standby"] != nil {
		model.RecoverFromStandby = core.BoolPtr(modelMap["recover_from_standby"].(bool))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["start_time_usecs"] != nil && modelMap["start_time_usecs"].(int) != 0 {
		model.StartTimeUsecs = core.Int64Ptr(int64(modelMap["start_time_usecs"].(int)))
	}
	if modelMap["end_time_usecs"] != nil && modelMap["end_time_usecs"].(int) != 0 {
		model.EndTimeUsecs = core.Int64Ptr(int64(modelMap["end_time_usecs"].(int)))
	}
	if modelMap["messages"] != nil {
		messages := []string{}
		for _, messagesItem := range modelMap["messages"].([]interface{}) {
			messages = append(messages, messagesItem.(string))
		}
		model.Messages = messages
	}
	if modelMap["bytes_restored"] != nil && modelMap["bytes_restored"].(int) != 0 {
		model.BytesRestored = core.Int64Ptr(int64(modelMap["bytes_restored"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToCommonRecoverObjectSnapshotParamsObjectInfo(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverObjectSnapshotParamsObjectInfo, error) {
	model := &backuprecoveryv1.CommonRecoverObjectSnapshotParamsObjectInfo{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["source_id"] != nil {
		model.SourceID = core.Int64Ptr(int64(modelMap["source_id"].(int)))
	}
	if modelMap["source_name"] != nil && modelMap["source_name"].(string) != "" {
		model.SourceName = core.StringPtr(modelMap["source_name"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	if modelMap["object_hash"] != nil && modelMap["object_hash"].(string) != "" {
		model.ObjectHash = core.StringPtr(modelMap["object_hash"].(string))
	}
	if modelMap["object_type"] != nil && modelMap["object_type"].(string) != "" {
		model.ObjectType = core.StringPtr(modelMap["object_type"].(string))
	}
	if modelMap["logical_size_bytes"] != nil {
		model.LogicalSizeBytes = core.Int64Ptr(int64(modelMap["logical_size_bytes"].(int)))
	}
	if modelMap["uuid"] != nil && modelMap["uuid"].(string) != "" {
		model.UUID = core.StringPtr(modelMap["uuid"].(string))
	}
	if modelMap["global_id"] != nil && modelMap["global_id"].(string) != "" {
		model.GlobalID = core.StringPtr(modelMap["global_id"].(string))
	}
	if modelMap["protection_type"] != nil && modelMap["protection_type"].(string) != "" {
		model.ProtectionType = core.StringPtr(modelMap["protection_type"].(string))
	}
	if modelMap["sharepoint_site_summary"] != nil && len(modelMap["sharepoint_site_summary"].([]interface{})) > 0 {
		SharepointSiteSummaryModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToSharepointObjectParams(modelMap["sharepoint_site_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SharepointSiteSummary = SharepointSiteSummaryModel
	}
	if modelMap["os_type"] != nil && modelMap["os_type"].(string) != "" {
		model.OsType = core.StringPtr(modelMap["os_type"].(string))
	}
	if modelMap["child_objects"] != nil {
		childObjects := []backuprecoveryv1.ObjectSummary{}
		for _, childObjectsItem := range modelMap["child_objects"].([]interface{}) {
			childObjectsItemModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectSummary(childObjectsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			childObjects = append(childObjects, *childObjectsItemModel)
		}
		model.ChildObjects = childObjects
	}
	if modelMap["v_center_summary"] != nil && len(modelMap["v_center_summary"].([]interface{})) > 0 {
		VCenterSummaryModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectTypeVCenterParams(modelMap["v_center_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VCenterSummary = VCenterSummaryModel
	}
	if modelMap["windows_cluster_summary"] != nil && len(modelMap["windows_cluster_summary"].([]interface{})) > 0 {
		WindowsClusterSummaryModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectTypeWindowsClusterParams(modelMap["windows_cluster_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WindowsClusterSummary = WindowsClusterSummaryModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToSharepointObjectParams(modelMap map[string]interface{}) (*backuprecoveryv1.SharepointObjectParams, error) {
	model := &backuprecoveryv1.SharepointObjectParams{}
	if modelMap["site_web_url"] != nil && modelMap["site_web_url"].(string) != "" {
		model.SiteWebURL = core.StringPtr(modelMap["site_web_url"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectSummary(modelMap map[string]interface{}) (*backuprecoveryv1.ObjectSummary, error) {
	model := &backuprecoveryv1.ObjectSummary{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["source_id"] != nil {
		model.SourceID = core.Int64Ptr(int64(modelMap["source_id"].(int)))
	}
	if modelMap["source_name"] != nil && modelMap["source_name"].(string) != "" {
		model.SourceName = core.StringPtr(modelMap["source_name"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	if modelMap["object_hash"] != nil && modelMap["object_hash"].(string) != "" {
		model.ObjectHash = core.StringPtr(modelMap["object_hash"].(string))
	}
	if modelMap["object_type"] != nil && modelMap["object_type"].(string) != "" {
		model.ObjectType = core.StringPtr(modelMap["object_type"].(string))
	}
	if modelMap["logical_size_bytes"] != nil {
		model.LogicalSizeBytes = core.Int64Ptr(int64(modelMap["logical_size_bytes"].(int)))
	}
	if modelMap["uuid"] != nil && modelMap["uuid"].(string) != "" {
		model.UUID = core.StringPtr(modelMap["uuid"].(string))
	}
	if modelMap["global_id"] != nil && modelMap["global_id"].(string) != "" {
		model.GlobalID = core.StringPtr(modelMap["global_id"].(string))
	}
	if modelMap["protection_type"] != nil && modelMap["protection_type"].(string) != "" {
		model.ProtectionType = core.StringPtr(modelMap["protection_type"].(string))
	}
	if modelMap["sharepoint_site_summary"] != nil && len(modelMap["sharepoint_site_summary"].([]interface{})) > 0 {
		SharepointSiteSummaryModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToSharepointObjectParams(modelMap["sharepoint_site_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SharepointSiteSummary = SharepointSiteSummaryModel
	}
	if modelMap["os_type"] != nil && modelMap["os_type"].(string) != "" {
		model.OsType = core.StringPtr(modelMap["os_type"].(string))
	}
	if modelMap["child_objects"] != nil {
		childObjects := []backuprecoveryv1.ObjectSummary{}
		for _, childObjectsItem := range modelMap["child_objects"].([]interface{}) {
			childObjectsItemModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectSummary(childObjectsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			childObjects = append(childObjects, *childObjectsItemModel)
		}
		model.ChildObjects = childObjects
	}
	if modelMap["v_center_summary"] != nil && len(modelMap["v_center_summary"].([]interface{})) > 0 {
		VCenterSummaryModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectTypeVCenterParams(modelMap["v_center_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VCenterSummary = VCenterSummaryModel
	}
	if modelMap["windows_cluster_summary"] != nil && len(modelMap["windows_cluster_summary"].([]interface{})) > 0 {
		WindowsClusterSummaryModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectTypeWindowsClusterParams(modelMap["windows_cluster_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WindowsClusterSummary = WindowsClusterSummaryModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectTypeVCenterParams(modelMap map[string]interface{}) (*backuprecoveryv1.ObjectTypeVCenterParams, error) {
	model := &backuprecoveryv1.ObjectTypeVCenterParams{}
	if modelMap["is_cloud_env"] != nil {
		model.IsCloudEnv = core.BoolPtr(modelMap["is_cloud_env"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToObjectTypeWindowsClusterParams(modelMap map[string]interface{}) (*backuprecoveryv1.ObjectTypeWindowsClusterParams, error) {
	model := &backuprecoveryv1.ObjectTypeWindowsClusterParams{}
	if modelMap["cluster_source_type"] != nil && modelMap["cluster_source_type"].(string) != "" {
		model.ClusterSourceType = core.StringPtr(modelMap["cluster_source_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToCommonRecoverObjectSnapshotParamsArchivalTargetInfo(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverObjectSnapshotParamsArchivalTargetInfo, error) {
	model := &backuprecoveryv1.CommonRecoverObjectSnapshotParamsArchivalTargetInfo{}
	if modelMap["target_id"] != nil {
		model.TargetID = core.Int64Ptr(int64(modelMap["target_id"].(int)))
	}
	if modelMap["archival_task_id"] != nil && modelMap["archival_task_id"].(string) != "" {
		model.ArchivalTaskID = core.StringPtr(modelMap["archival_task_id"].(string))
	}
	if modelMap["target_name"] != nil && modelMap["target_name"].(string) != "" {
		model.TargetName = core.StringPtr(modelMap["target_name"].(string))
	}
	if modelMap["target_type"] != nil && modelMap["target_type"].(string) != "" {
		model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	}
	if modelMap["usage_type"] != nil && modelMap["usage_type"].(string) != "" {
		model.UsageType = core.StringPtr(modelMap["usage_type"].(string))
	}
	if modelMap["ownership_context"] != nil && modelMap["ownership_context"].(string) != "" {
		model.OwnershipContext = core.StringPtr(modelMap["ownership_context"].(string))
	}
	if modelMap["tier_settings"] != nil && len(modelMap["tier_settings"].([]interface{})) > 0 {
		TierSettingsModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToArchivalTargetTierInfo(modelMap["tier_settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.TierSettings = TierSettingsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToArchivalTargetTierInfo(modelMap map[string]interface{}) (*backuprecoveryv1.ArchivalTargetTierInfo, error) {
	model := &backuprecoveryv1.ArchivalTargetTierInfo{}
	if modelMap["aws_tiering"] != nil && len(modelMap["aws_tiering"].([]interface{})) > 0 {
		AwsTieringModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAWSTiers(modelMap["aws_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AwsTiering = AwsTieringModel
	}
	if modelMap["azure_tiering"] != nil && len(modelMap["azure_tiering"].([]interface{})) > 0 {
		AzureTieringModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAzureTiers(modelMap["azure_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AzureTiering = AzureTieringModel
	}
	if modelMap["cloud_platform"] != nil && modelMap["cloud_platform"].(string) != "" {
		model.CloudPlatform = core.StringPtr(modelMap["cloud_platform"].(string))
	}
	if modelMap["google_tiering"] != nil && len(modelMap["google_tiering"].([]interface{})) > 0 {
		GoogleTieringModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToGoogleTiers(modelMap["google_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.GoogleTiering = GoogleTieringModel
	}
	if modelMap["oracle_tiering"] != nil && len(modelMap["oracle_tiering"].([]interface{})) > 0 {
		OracleTieringModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToOracleTiers(modelMap["oracle_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.OracleTiering = OracleTieringModel
	}
	if modelMap["current_tier_type"] != nil && modelMap["current_tier_type"].(string) != "" {
		model.CurrentTierType = core.StringPtr(modelMap["current_tier_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAWSTiers(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTiers, error) {
	model := &backuprecoveryv1.AWSTiers{}
	tiers := []backuprecoveryv1.AWSTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAWSTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAWSTier(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTier, error) {
	model := &backuprecoveryv1.AWSTier{}
	if modelMap["move_after_unit"] != nil && modelMap["move_after_unit"].(string) != "" {
		model.MoveAfterUnit = core.StringPtr(modelMap["move_after_unit"].(string))
	}
	if modelMap["move_after"] != nil {
		model.MoveAfter = core.Int64Ptr(int64(modelMap["move_after"].(int)))
	}
	model.TierType = core.StringPtr(modelMap["tier_type"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAzureTiers(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTiers, error) {
	model := &backuprecoveryv1.AzureTiers{}
	if modelMap["tiers"] != nil {
		tiers := []backuprecoveryv1.AzureTier{}
		for _, tiersItem := range modelMap["tiers"].([]interface{}) {
			tiersItemModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAzureTier(tiersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			tiers = append(tiers, *tiersItemModel)
		}
		model.Tiers = tiers
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToAzureTier(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTier, error) {
	model := &backuprecoveryv1.AzureTier{}
	if modelMap["move_after_unit"] != nil && modelMap["move_after_unit"].(string) != "" {
		model.MoveAfterUnit = core.StringPtr(modelMap["move_after_unit"].(string))
	}
	if modelMap["move_after"] != nil {
		model.MoveAfter = core.Int64Ptr(int64(modelMap["move_after"].(int)))
	}
	model.TierType = core.StringPtr(modelMap["tier_type"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToGoogleTiers(modelMap map[string]interface{}) (*backuprecoveryv1.GoogleTiers, error) {
	model := &backuprecoveryv1.GoogleTiers{}
	tiers := []backuprecoveryv1.GoogleTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToGoogleTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToGoogleTier(modelMap map[string]interface{}) (*backuprecoveryv1.GoogleTier, error) {
	model := &backuprecoveryv1.GoogleTier{}
	if modelMap["move_after_unit"] != nil && modelMap["move_after_unit"].(string) != "" {
		model.MoveAfterUnit = core.StringPtr(modelMap["move_after_unit"].(string))
	}
	if modelMap["move_after"] != nil {
		model.MoveAfter = core.Int64Ptr(int64(modelMap["move_after"].(int)))
	}
	model.TierType = core.StringPtr(modelMap["tier_type"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToOracleTiers(modelMap map[string]interface{}) (*backuprecoveryv1.OracleTiers, error) {
	model := &backuprecoveryv1.OracleTiers{}
	tiers := []backuprecoveryv1.OracleTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryDownloadFilesFoldersMapToOracleTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToOracleTier(modelMap map[string]interface{}) (*backuprecoveryv1.OracleTier, error) {
	model := &backuprecoveryv1.OracleTier{}
	if modelMap["move_after_unit"] != nil && modelMap["move_after_unit"].(string) != "" {
		model.MoveAfterUnit = core.StringPtr(modelMap["move_after_unit"].(string))
	}
	if modelMap["move_after"] != nil {
		model.MoveAfter = core.Int64Ptr(int64(modelMap["move_after"].(int)))
	}
	model.TierType = core.StringPtr(modelMap["tier_type"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToFilesAndFoldersObject(modelMap map[string]interface{}) (*backuprecoveryv1.FilesAndFoldersObject, error) {
	model := &backuprecoveryv1.FilesAndFoldersObject{}
	model.AbsolutePath = core.StringPtr(modelMap["absolute_path"].(string))
	if modelMap["is_directory"] != nil {
		model.IsDirectory = core.BoolPtr(modelMap["is_directory"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersMapToDocumentObject(modelMap map[string]interface{}) (*backuprecoveryv1.DocumentObject, error) {
	model := &backuprecoveryv1.DocumentObject{}
	if modelMap["is_directory"] != nil {
		model.IsDirectory = core.BoolPtr(modelMap["is_directory"].(bool))
	}
	model.ItemID = core.StringPtr(modelMap["item_id"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersDocumentObjectToMap(model *backuprecoveryv1.DocumentObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsDirectory != nil {
		modelMap["is_directory"] = *model.IsDirectory
	}
	modelMap["item_id"] = *model.ItemID
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersCommonRecoverObjectSnapshotParamsToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["snapshot_id"] = *model.SnapshotID
	if model.PointInTimeUsecs != nil {
		modelMap["point_in_time_usecs"] = flex.IntValue(model.PointInTimeUsecs)
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.SnapshotCreationTimeUsecs != nil && *(model.SnapshotCreationTimeUsecs) != 0 {
		modelMap["snapshot_creation_time_usecs"] = flex.IntValue(model.SnapshotCreationTimeUsecs)
	}
	if model.ObjectInfo != nil {
		objectInfoMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersCommonRecoverObjectSnapshotParamsObjectInfoToMap(model.ObjectInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["object_info"] = []map[string]interface{}{objectInfoMap}
	}
	if model.SnapshotTargetType != nil {
		modelMap["snapshot_target_type"] = *model.SnapshotTargetType
	}
	if model.ArchivalTargetInfo != nil {
		archivalTargetInfoMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model.ArchivalTargetInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_target_info"] = []map[string]interface{}{archivalTargetInfoMap}
	}
	if model.ProgressTaskID != nil {
		modelMap["progress_task_id"] = *model.ProgressTaskID
	}
	if model.RecoverFromStandby != nil {
		modelMap["recover_from_standby"] = *model.RecoverFromStandby
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.StartTimeUsecs != nil && *(model.StartTimeUsecs) != 0 {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil && *(model.EndTimeUsecs) != 0 {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Messages != nil {
		modelMap["messages"] = model.Messages
	}
	if model.BytesRestored != nil && *(model.BytesRestored) != 0 {
		modelMap["bytes_restored"] = flex.IntValue(model.BytesRestored)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersCommonRecoverObjectSnapshotParamsObjectInfoToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParamsObjectInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersSharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParamsArchivalTargetInfo) (map[string]interface{}, error) {
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
		tierSettingsMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersArchivalTargetTierInfoToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersArchivalTargetTierInfoToMap(model *backuprecoveryv1.ArchivalTargetTierInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}
	if model.GoogleTiering != nil {
		googleTieringMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersOracleTiersToMap(model.OracleTiering)
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

func ResourceIbmBackupRecoveryDownloadFilesFoldersAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryDownloadFilesFoldersAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryDownloadFilesFoldersGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryDownloadFilesFoldersOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryDownloadFilesFoldersOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryDownloadFilesFoldersOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryDownloadFilesFoldersFilesAndFoldersObjectToMap(model *backuprecoveryv1.FilesAndFoldersObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["absolute_path"] = *model.AbsolutePath
	if model.IsDirectory != nil {
		modelMap["is_directory"] = *model.IsDirectory
	}
	return modelMap, nil
}
