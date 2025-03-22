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
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIbmBackupRecovery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryCreate,
		ReadContext:   resourceIbmBackupRecoveryRead,
		DeleteContext: resourceIbmBackupRecoveryDelete,
		UpdateContext: resourceIbmBackupRecoveryUpdate,
		Importer:      &schema.ResourceImporter{},
		CustomizeDiff: checkDiffResourceIbmBackupRecovery,
		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"request_initiator_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_recovery", "request_initiator_type"),
				Description: "Specifies the type of request from UI, which is used for services like magneto to determine the priority of requests.",
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:    true,
				Description: "Specifies the name of the Recovery.",
			},
			"recovery_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Recovery ID",
			},
			"snapshot_environment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:     true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_recovery", "snapshot_environment"),
				Description: "Specifies the type of snapshot environment for which the Recovery was performed.",
			},
			"physical_params": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
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
										Computed:    true,
										Description: "Specifies the protection group id of the object snapshot.",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
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
			"mssql_params": &schema.Schema{
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
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
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the Recovery. 'Running' indicates that the Recovery is still running. 'Canceled' indicates that the Recovery has been cancelled. 'Canceling' indicates that the Recovery is in the process of being cancelled. 'Failed' indicates that the Recovery has failed. 'Succeeded' indicates that the Recovery has finished successfully. 'SucceededWithWarning' indicates that the Recovery finished successfully, but there were some warning messages. 'Skipped' indicates that the Recovery task was skipped.",
			},
			"progress_task_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Progress monitor task id for Recovery.",
			},
			"recovery_action": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of recover action.",
			},
			"permissions": &schema.Schema{
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
												"metrics_config": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the metadata for metrics configuration. The metadata defined here will be used by cluster to send the usgae metrics to IBM cloud metering service for calculating the tenant billing.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cos_resource_config": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the details of COS resource configuration required for posting metrics and trackinb billing information for IBM tenants.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"resource_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the resource COS resource configuration endpoint that will be used for fetching bucket usage for a given tenant.",
																		},
																	},
																},
															},
															"iam_metrics_config": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the IAM configuration that will be used for accessing the billing service in IBM cloud.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"iam_url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the IAM URL needed to fetch the operator token from IBM. The operator token is needed to make service API calls to IBM billing service.",
																		},
																		"billing_api_key_secret_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies Id of the secret that contains the API key.",
																		},
																	},
																},
															},
															"metering_config": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the metering configuration that will be used for IBM cluster to send the billing details to IBM billing service.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"part_ids": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the list of part identifiers used for metrics identification.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"submission_interval_in_secs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Computed:    true,
																			Description: "Specifies the frequency in seconds at which the metrics will be pushed to IBM billing service from cluster.",
																		},
																		"url": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies the current ownership mode for the tenant. The ownership of the tenant represents the active role for functioning of the tenant.",
												},
												"plan_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the Plan Id associated with the tenant. This field is introduced for tracking purposes inside IBM enviournment.",
												},
												"resource_group_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the Resource Group ID associated with the tenant.",
												},
												"resource_instance_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
			"creation_info": &schema.Schema{
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
			"can_tear_down": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether it's possible to tear down the objects created by the recovery.",
			},
			"tear_down_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the status of the tear down operation. This is only set when the canTearDown is set to true. 'DestroyScheduled' indicates that the tear down is ready to schedule. 'Destroying' indicates that the tear down is still running. 'Destroyed' indicates that the tear down succeeded. 'DestroyError' indicates that the tear down failed.",
			},
			"tear_down_message": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the error message about the tear down operation if it fails.",
			},
			"messages": &schema.Schema{
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
			"parent_recovery_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If current recovery is child recovery triggered by another parent recovery operation, then this field willt specify the id of the parent recovery.",
			},
			"retrieve_archive_tasks": &schema.Schema{
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
			"is_multi_stage_restore": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the current recovery operation is a multi-stage restore operation. This is currently used by VMware recoveres for the migration/hot-standby use case.",
			},
		},
	}
}

func checkDiffResourceIbmBackupRecovery(context context.Context, d *schema.ResourceDiff, meta interface{}) error {
	// oldId, _ := d.GetChange("x_ibm_tenant_id")
	// if oldId == "" {
	// 	return nil
	// }

	// return if it's a new resource
	if d.Id() == "" {
		return nil
	}

	for fieldName := range ResourceIbmBackupRecovery().Schema {
		if d.HasChange(fieldName) {
			return fmt.Errorf("[ERROR] Resource ibm_backup_recovery_recovery cannot be updated. Field: %s", fieldName)
		}
	}
	return nil
}

func resourceIbmBackupRecoveryCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createRecoveryOptions := &backuprecoveryv1.CreateRecoveryOptions{}
	tenantId := d.Get("x_ibm_tenant_id").(string)
	createRecoveryOptions.SetXIBMTenantID(tenantId)
	createRecoveryOptions.SetName(d.Get("name").(string))
	createRecoveryOptions.SetSnapshotEnvironment(d.Get("snapshot_environment").(string))
	if _, ok := d.GetOk("physical_params"); ok {
		physicalParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalParams(d.Get("physical_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "create", "parse-physical_params").GetDiag()
		}
		createRecoveryOptions.SetPhysicalParams(physicalParamsModel)
	}
	if _, ok := d.GetOk("mssql_params"); ok {
		mssqlParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverSqlParams(d.Get("mssql_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "create", "parse-mssql_params").GetDiag()
		}
		createRecoveryOptions.SetMssqlParams(mssqlParamsModel)
	}
	if _, ok := d.GetOk("request_initiator_type"); ok {
		createRecoveryOptions.SetRequestInitiatorType(d.Get("request_initiator_type").(string))
	}

	recovery, _, err := backupRecoveryClient.CreateRecoveryWithContext(context, createRecoveryOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateRecoveryWithContext failed: %s", err.Error()), "ibm_backup_recovery_recovery", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	recoveryId := fmt.Sprintf("%s::%s", tenantId, *recovery.ID)
	d.SetId(recoveryId)

	return resourceIbmBackupRecoveryRead(context, d, meta)
}

func resourceIbmBackupRecoveryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getRecoveryByIdOptions := &backuprecoveryv1.GetRecoveryByIdOptions{}
	tenantId := d.Get("x_ibm_tenant_id").(string)
	recoveryId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		recoveryId = ParseId(d.Id(), "id")
	}
	getRecoveryByIdOptions.SetID(recoveryId)
	getRecoveryByIdOptions.SetXIBMTenantID(tenantId)

	recovery, response, err := backupRecoveryClient.GetRecoveryByIDWithContext(context, getRecoveryByIdOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRecoveryByIDWithContext failed: %s", err.Error()), "ibm_backup_recovery_recovery", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", recovery.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-name").GetDiag()
	}
	if err = d.Set("x_ibm_tenant_id", tenantId); err != nil {
		err = fmt.Errorf("Error setting x_ibm_tenant_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-x_ibm_tenant_id").GetDiag()
	}
	if err = d.Set("recovery_id", recoveryId); err != nil {
		err = fmt.Errorf("Error setting recovery_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-recovery_id").GetDiag()
	}
	if err = d.Set("snapshot_environment", recovery.SnapshotEnvironment); err != nil {
		err = fmt.Errorf("Error setting snapshot_environment: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-snapshot_environment").GetDiag()
	}
	if !core.IsNil(recovery.PhysicalParams) {
		physicalParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalParamsToMap(recovery.PhysicalParams)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "physical_params-to-map").GetDiag()
		}
		if err = d.Set("physical_params", []map[string]interface{}{physicalParamsMap}); err != nil {
			err = fmt.Errorf("Error setting physical_params: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-physical_params").GetDiag()
		}
	}
	if !core.IsNil(recovery.MssqlParams) {
		mssqlParamsMap, err := ResourceIbmBackupRecoveryRecoverSqlParamsToMap(recovery.MssqlParams)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "mssql_params-to-map").GetDiag()
		}
		if err = d.Set("mssql_params", []map[string]interface{}{mssqlParamsMap}); err != nil {
			err = fmt.Errorf("Error setting mssql_params: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-mssql_params").GetDiag()
		}
	}
	if !core.IsNil(recovery.StartTimeUsecs) {
		if err = d.Set("start_time_usecs", flex.IntValue(recovery.StartTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting start_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-start_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(recovery.EndTimeUsecs) {
		if err = d.Set("end_time_usecs", flex.IntValue(recovery.EndTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting end_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-end_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(recovery.Status) {
		if err = d.Set("status", recovery.Status); err != nil {
			err = fmt.Errorf("Error setting status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-status").GetDiag()
		}
	}
	if !core.IsNil(recovery.ProgressTaskID) {
		if err = d.Set("progress_task_id", recovery.ProgressTaskID); err != nil {
			err = fmt.Errorf("Error setting progress_task_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-progress_task_id").GetDiag()
		}
	}
	if !core.IsNil(recovery.RecoveryAction) {
		if err = d.Set("recovery_action", recovery.RecoveryAction); err != nil {
			err = fmt.Errorf("Error setting recovery_action: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-recovery_action").GetDiag()
		}
	}
	if !core.IsNil(recovery.Permissions) {
		permissions := []map[string]interface{}{}
		for _, permissionsItem := range recovery.Permissions {
			permissionsItemMap, err := ResourceIbmBackupRecoveryTenantToMap(&permissionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "permissions-to-map").GetDiag()
			}
			permissions = append(permissions, permissionsItemMap)
		}
		if err = d.Set("permissions", permissions); err != nil {
			err = fmt.Errorf("Error setting permissions: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-permissions").GetDiag()
		}
	}
	if !core.IsNil(recovery.CreationInfo) {
		creationInfoMap, err := ResourceIbmBackupRecoveryCreationInfoToMap(recovery.CreationInfo)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "creation_info-to-map").GetDiag()
		}
		if err = d.Set("creation_info", []map[string]interface{}{creationInfoMap}); err != nil {
			err = fmt.Errorf("Error setting creation_info: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-creation_info").GetDiag()
		}
	}
	if !core.IsNil(recovery.CanTearDown) {
		if err = d.Set("can_tear_down", recovery.CanTearDown); err != nil {
			err = fmt.Errorf("Error setting can_tear_down: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-can_tear_down").GetDiag()
		}
	}
	if !core.IsNil(recovery.TearDownStatus) {
		if err = d.Set("tear_down_status", recovery.TearDownStatus); err != nil {
			err = fmt.Errorf("Error setting tear_down_status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-tear_down_status").GetDiag()
		}
	}
	if !core.IsNil(recovery.TearDownMessage) {
		if err = d.Set("tear_down_message", recovery.TearDownMessage); err != nil {
			err = fmt.Errorf("Error setting tear_down_message: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-tear_down_message").GetDiag()
		}
	}
	if !core.IsNil(recovery.Messages) {
		if err = d.Set("messages", recovery.Messages); err != nil {
			err = fmt.Errorf("Error setting messages: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-messages").GetDiag()
		}
	} else {
		if err = d.Set("messages", []interface{}{}); err != nil {
			err = fmt.Errorf("Error setting messages: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-messages").GetDiag()
		}
	}
	if !core.IsNil(recovery.IsParentRecovery) {
		if err = d.Set("is_parent_recovery", recovery.IsParentRecovery); err != nil {
			err = fmt.Errorf("Error setting is_parent_recovery: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-is_parent_recovery").GetDiag()
		}
	}
	if !core.IsNil(recovery.ParentRecoveryID) {
		if err = d.Set("parent_recovery_id", recovery.ParentRecoveryID); err != nil {
			err = fmt.Errorf("Error setting parent_recovery_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-parent_recovery_id").GetDiag()
		}
	}
	if !core.IsNil(recovery.RetrieveArchiveTasks) {
		retrieveArchiveTasks := []map[string]interface{}{}
		for _, retrieveArchiveTasksItem := range recovery.RetrieveArchiveTasks {
			retrieveArchiveTasksItemMap, err := ResourceIbmBackupRecoveryRetrieveArchiveTaskToMap(&retrieveArchiveTasksItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "retrieve_archive_tasks-to-map").GetDiag()
			}
			retrieveArchiveTasks = append(retrieveArchiveTasks, retrieveArchiveTasksItemMap)
		}
		if err = d.Set("retrieve_archive_tasks", retrieveArchiveTasks); err != nil {
			err = fmt.Errorf("Error setting retrieve_archive_tasks: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-retrieve_archive_tasks").GetDiag()
		}
	} else {
		if err = d.Set("retrieve_archive_tasks", []interface{}{}); err != nil {
			err = fmt.Errorf("Error setting retrieve_archive_tasks: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-retrieve_archive_tasks").GetDiag()
		}
	}
	if !core.IsNil(recovery.IsMultiStageRestore) {
		if err = d.Set("is_multi_stage_restore", recovery.IsMultiStageRestore); err != nil {
			err = fmt.Errorf("Error setting is_multi_stage_restore: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_recovery", "read", "set-is_multi_stage_restore").GetDiag()
		}
	}

	return nil
}

func resourceIbmBackupRecoveryDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceIbmBackupRecoveryUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func ResourceIbmBackupRecoveryMapToRecoverPhysicalParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalParams{}
	objects := []backuprecoveryv1.CommonRecoverObjectSnapshotParams{}
	for _, objectsItem := range modelMap["objects"].([]interface{}) {
		objectsItemModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParams(objectsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		objects = append(objects, *objectsItemModel)
	}
	model.Objects = objects
	model.RecoveryAction = core.StringPtr(modelMap["recovery_action"].(string))
	if modelMap["recover_volume_params"] != nil && len(modelMap["recover_volume_params"].([]interface{})) > 0 {
		RecoverVolumeParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsRecoverVolumeParams(modelMap["recover_volume_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.RecoverVolumeParams = RecoverVolumeParamsModel
	}
	if modelMap["mount_volume_params"] != nil && len(modelMap["mount_volume_params"].([]interface{})) > 0 {
		MountVolumeParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsMountVolumeParams(modelMap["mount_volume_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MountVolumeParams = MountVolumeParamsModel
	}
	if modelMap["recover_file_and_folder_params"] != nil && len(modelMap["recover_file_and_folder_params"].([]interface{})) > 0 {
		RecoverFileAndFolderParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsRecoverFileAndFolderParams(modelMap["recover_file_and_folder_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.RecoverFileAndFolderParams = RecoverFileAndFolderParamsModel
	}
	if modelMap["download_file_and_folder_params"] != nil && len(modelMap["download_file_and_folder_params"].([]interface{})) > 0 {
		DownloadFileAndFolderParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsDownloadFileAndFolderParams(modelMap["download_file_and_folder_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DownloadFileAndFolderParams = DownloadFileAndFolderParamsModel
	}
	if modelMap["system_recovery_params"] != nil && len(modelMap["system_recovery_params"].([]interface{})) > 0 {
		SystemRecoveryParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsSystemRecoveryParams(modelMap["system_recovery_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SystemRecoveryParams = SystemRecoveryParamsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParams(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverObjectSnapshotParams, error) {
	model := &backuprecoveryv1.CommonRecoverObjectSnapshotParams{}
	model.SnapshotID = core.StringPtr(modelMap["snapshot_id"].(string))
	if modelMap["point_in_time_usecs"] != nil && modelMap["point_in_time_usecs"].(int) != 0 {
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
		ObjectInfoModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParamsObjectInfo(modelMap["object_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ObjectInfo = ObjectInfoModel
	}
	if modelMap["snapshot_target_type"] != nil && modelMap["snapshot_target_type"].(string) != "" {
		model.SnapshotTargetType = core.StringPtr(modelMap["snapshot_target_type"].(string))
	}
	if modelMap["archival_target_info"] != nil && len(modelMap["archival_target_info"].([]interface{})) > 0 {
		ArchivalTargetInfoModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParamsArchivalTargetInfo(modelMap["archival_target_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ArchivalTargetInfo = ArchivalTargetInfoModel
	}
	if modelMap["progress_task_id"] != nil && modelMap["progress_task_id"].(string) != "" {
		model.ProgressTaskID = core.StringPtr(modelMap["progress_task_id"].(string))
	}
	if modelMap["recover_from_standby"] != nil && modelMap["recover_from_standby"].(bool) != false {
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

func ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParamsObjectInfo(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverObjectSnapshotParamsObjectInfo, error) {
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
		SharepointSiteSummaryModel, err := ResourceIbmBackupRecoveryMapToSharepointObjectParams(modelMap["sharepoint_site_summary"].([]interface{})[0].(map[string]interface{}))
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
			childObjectsItemModel, err := ResourceIbmBackupRecoveryMapToObjectSummary(childObjectsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			childObjects = append(childObjects, *childObjectsItemModel)
		}
		model.ChildObjects = childObjects
	}
	if modelMap["v_center_summary"] != nil && len(modelMap["v_center_summary"].([]interface{})) > 0 {
		VCenterSummaryModel, err := ResourceIbmBackupRecoveryMapToObjectTypeVCenterParams(modelMap["v_center_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VCenterSummary = VCenterSummaryModel
	}
	if modelMap["windows_cluster_summary"] != nil && len(modelMap["windows_cluster_summary"].([]interface{})) > 0 {
		WindowsClusterSummaryModel, err := ResourceIbmBackupRecoveryMapToObjectTypeWindowsClusterParams(modelMap["windows_cluster_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WindowsClusterSummary = WindowsClusterSummaryModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToSharepointObjectParams(modelMap map[string]interface{}) (*backuprecoveryv1.SharepointObjectParams, error) {
	model := &backuprecoveryv1.SharepointObjectParams{}
	if modelMap["site_web_url"] != nil && modelMap["site_web_url"].(string) != "" {
		model.SiteWebURL = core.StringPtr(modelMap["site_web_url"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToObjectSummary(modelMap map[string]interface{}) (*backuprecoveryv1.ObjectSummary, error) {
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
		SharepointSiteSummaryModel, err := ResourceIbmBackupRecoveryMapToSharepointObjectParams(modelMap["sharepoint_site_summary"].([]interface{})[0].(map[string]interface{}))
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
			childObjectsItemModel, err := ResourceIbmBackupRecoveryMapToObjectSummary(childObjectsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			childObjects = append(childObjects, *childObjectsItemModel)
		}
		model.ChildObjects = childObjects
	}
	if modelMap["v_center_summary"] != nil && len(modelMap["v_center_summary"].([]interface{})) > 0 {
		VCenterSummaryModel, err := ResourceIbmBackupRecoveryMapToObjectTypeVCenterParams(modelMap["v_center_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VCenterSummary = VCenterSummaryModel
	}
	if modelMap["windows_cluster_summary"] != nil && len(modelMap["windows_cluster_summary"].([]interface{})) > 0 {
		WindowsClusterSummaryModel, err := ResourceIbmBackupRecoveryMapToObjectTypeWindowsClusterParams(modelMap["windows_cluster_summary"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WindowsClusterSummary = WindowsClusterSummaryModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToObjectTypeVCenterParams(modelMap map[string]interface{}) (*backuprecoveryv1.ObjectTypeVCenterParams, error) {
	model := &backuprecoveryv1.ObjectTypeVCenterParams{}
	if modelMap["is_cloud_env"] != nil {
		model.IsCloudEnv = core.BoolPtr(modelMap["is_cloud_env"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToObjectTypeWindowsClusterParams(modelMap map[string]interface{}) (*backuprecoveryv1.ObjectTypeWindowsClusterParams, error) {
	model := &backuprecoveryv1.ObjectTypeWindowsClusterParams{}
	if modelMap["cluster_source_type"] != nil && modelMap["cluster_source_type"].(string) != "" {
		model.ClusterSourceType = core.StringPtr(modelMap["cluster_source_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParamsArchivalTargetInfo(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverObjectSnapshotParamsArchivalTargetInfo, error) {
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
		TierSettingsModel, err := ResourceIbmBackupRecoveryMapToArchivalTargetTierInfo(modelMap["tier_settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.TierSettings = TierSettingsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToArchivalTargetTierInfo(modelMap map[string]interface{}) (*backuprecoveryv1.ArchivalTargetTierInfo, error) {
	model := &backuprecoveryv1.ArchivalTargetTierInfo{}
	if modelMap["aws_tiering"] != nil && len(modelMap["aws_tiering"].([]interface{})) > 0 {
		AwsTieringModel, err := ResourceIbmBackupRecoveryMapToAWSTiers(modelMap["aws_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AwsTiering = AwsTieringModel
	}
	if modelMap["azure_tiering"] != nil && len(modelMap["azure_tiering"].([]interface{})) > 0 {
		AzureTieringModel, err := ResourceIbmBackupRecoveryMapToAzureTiers(modelMap["azure_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AzureTiering = AzureTieringModel
	}
	if modelMap["cloud_platform"] != nil && modelMap["cloud_platform"].(string) != "" {
		model.CloudPlatform = core.StringPtr(modelMap["cloud_platform"].(string))
	}
	if modelMap["google_tiering"] != nil && len(modelMap["google_tiering"].([]interface{})) > 0 {
		GoogleTieringModel, err := ResourceIbmBackupRecoveryMapToGoogleTiers(modelMap["google_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.GoogleTiering = GoogleTieringModel
	}
	if modelMap["oracle_tiering"] != nil && len(modelMap["oracle_tiering"].([]interface{})) > 0 {
		OracleTieringModel, err := ResourceIbmBackupRecoveryMapToOracleTiers(modelMap["oracle_tiering"].([]interface{})[0].(map[string]interface{}))
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

func ResourceIbmBackupRecoveryMapToAWSTiers(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTiers, error) {
	model := &backuprecoveryv1.AWSTiers{}
	tiers := []backuprecoveryv1.AWSTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryMapToAWSTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryMapToAWSTier(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTier, error) {
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

func ResourceIbmBackupRecoveryMapToAzureTiers(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTiers, error) {
	model := &backuprecoveryv1.AzureTiers{}
	if modelMap["tiers"] != nil {
		tiers := []backuprecoveryv1.AzureTier{}
		for _, tiersItem := range modelMap["tiers"].([]interface{}) {
			tiersItemModel, err := ResourceIbmBackupRecoveryMapToAzureTier(tiersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			tiers = append(tiers, *tiersItemModel)
		}
		model.Tiers = tiers
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToAzureTier(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTier, error) {
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

func ResourceIbmBackupRecoveryMapToGoogleTiers(modelMap map[string]interface{}) (*backuprecoveryv1.GoogleTiers, error) {
	model := &backuprecoveryv1.GoogleTiers{}
	tiers := []backuprecoveryv1.GoogleTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryMapToGoogleTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryMapToGoogleTier(modelMap map[string]interface{}) (*backuprecoveryv1.GoogleTier, error) {
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

func ResourceIbmBackupRecoveryMapToOracleTiers(modelMap map[string]interface{}) (*backuprecoveryv1.OracleTiers, error) {
	model := &backuprecoveryv1.OracleTiers{}
	tiers := []backuprecoveryv1.OracleTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryMapToOracleTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryMapToOracleTier(modelMap map[string]interface{}) (*backuprecoveryv1.OracleTier, error) {
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

func ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsRecoverVolumeParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalParamsRecoverVolumeParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalParamsRecoverVolumeParams{}
	model.TargetEnvironment = core.StringPtr(modelMap["target_environment"].(string))
	if modelMap["physical_target_params"] != nil && len(modelMap["physical_target_params"].([]interface{})) > 0 {
		PhysicalTargetParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalVolumeParamsPhysicalTargetParams(modelMap["physical_target_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PhysicalTargetParams = PhysicalTargetParamsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverPhysicalVolumeParamsPhysicalTargetParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalVolumeParamsPhysicalTargetParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalVolumeParamsPhysicalTargetParams{}
	MountTargetModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverVolumeMountTarget(modelMap["mount_target"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.MountTarget = MountTargetModel
	volumeMapping := []backuprecoveryv1.RecoverVolumeMapping{}
	for _, volumeMappingItem := range modelMap["volume_mapping"].([]interface{}) {
		volumeMappingItemModel, err := ResourceIbmBackupRecoveryMapToRecoverVolumeMapping(volumeMappingItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		volumeMapping = append(volumeMapping, *volumeMappingItemModel)
	}
	model.VolumeMapping = volumeMapping
	if modelMap["force_unmount_volume"] != nil {
		model.ForceUnmountVolume = core.BoolPtr(modelMap["force_unmount_volume"].(bool))
	}
	if modelMap["vlan_config"] != nil && len(modelMap["vlan_config"].([]interface{})) > 0 {
		VlanConfigModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverVolumeVlanConfig(modelMap["vlan_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VlanConfig = VlanConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverVolumeMountTarget(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeMountTarget, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeMountTarget{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverVolumeMapping(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverVolumeMapping, error) {
	model := &backuprecoveryv1.RecoverVolumeMapping{}
	model.SourceVolumeGuid = core.StringPtr(modelMap["source_volume_guid"].(string))
	model.DestinationVolumeGuid = core.StringPtr(modelMap["destination_volume_guid"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverVolumeVlanConfig(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeVlanConfig, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeVlanConfig{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["disable_vlan"] != nil {
		model.DisableVlan = core.BoolPtr(modelMap["disable_vlan"].(bool))
	}
	if modelMap["interface_name"] != nil && modelMap["interface_name"].(string) != "" {
		model.InterfaceName = core.StringPtr(modelMap["interface_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsMountVolumeParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalParamsMountVolumeParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalParamsMountVolumeParams{}
	model.TargetEnvironment = core.StringPtr(modelMap["target_environment"].(string))
	if modelMap["physical_target_params"] != nil && len(modelMap["physical_target_params"].([]interface{})) > 0 {
		PhysicalTargetParamsModel, err := ResourceIbmBackupRecoveryMapToMountPhysicalVolumeParamsPhysicalTargetParams(modelMap["physical_target_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PhysicalTargetParams = PhysicalTargetParamsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToMountPhysicalVolumeParamsPhysicalTargetParams(modelMap map[string]interface{}) (*backuprecoveryv1.MountPhysicalVolumeParamsPhysicalTargetParams, error) {
	model := &backuprecoveryv1.MountPhysicalVolumeParamsPhysicalTargetParams{}
	model.MountToOriginalTarget = core.BoolPtr(modelMap["mount_to_original_target"].(bool))
	if modelMap["original_target_config"] != nil && len(modelMap["original_target_config"].([]interface{})) > 0 {
		OriginalTargetConfigModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForMountVolumeOriginalTargetConfig(modelMap["original_target_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.OriginalTargetConfig = OriginalTargetConfigModel
	}
	if modelMap["new_target_config"] != nil && len(modelMap["new_target_config"].([]interface{})) > 0 {
		NewTargetConfigModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForMountVolumeNewTargetConfig(modelMap["new_target_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.NewTargetConfig = NewTargetConfigModel
	}
	if modelMap["read_only_mount"] != nil {
		model.ReadOnlyMount = core.BoolPtr(modelMap["read_only_mount"].(bool))
	}
	if modelMap["volume_names"] != nil {
		volumeNames := []string{}
		for _, volumeNamesItem := range modelMap["volume_names"].([]interface{}) {
			volumeNames = append(volumeNames, volumeNamesItem.(string))
		}
		model.VolumeNames = volumeNames
	}
	if modelMap["mounted_volume_mapping"] != nil {
		mountedVolumeMapping := []backuprecoveryv1.MountedVolumeMapping{}
		for _, mountedVolumeMappingItem := range modelMap["mounted_volume_mapping"].([]interface{}) {
			mountedVolumeMappingItemModel, err := ResourceIbmBackupRecoveryMapToMountedVolumeMapping(mountedVolumeMappingItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			mountedVolumeMapping = append(mountedVolumeMapping, *mountedVolumeMappingItemModel)
		}
		model.MountedVolumeMapping = mountedVolumeMapping
	}
	if modelMap["vlan_config"] != nil && len(modelMap["vlan_config"].([]interface{})) > 0 {
		VlanConfigModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForMountVolumeVlanConfig(modelMap["vlan_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VlanConfig = VlanConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForMountVolumeOriginalTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForMountVolumeOriginalTargetConfig, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForMountVolumeOriginalTargetConfig{}
	if modelMap["server_credentials"] != nil && len(modelMap["server_credentials"].([]interface{})) > 0 {
		ServerCredentialsModel, err := ResourceIbmBackupRecoveryMapToPhysicalMountVolumesOriginalTargetConfigServerCredentials(modelMap["server_credentials"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ServerCredentials = ServerCredentialsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalMountVolumesOriginalTargetConfigServerCredentials(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalMountVolumesOriginalTargetConfigServerCredentials, error) {
	model := &backuprecoveryv1.PhysicalMountVolumesOriginalTargetConfigServerCredentials{}
	model.Username = core.StringPtr(modelMap["username"].(string))
	model.Password = core.StringPtr(modelMap["password"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForMountVolumeNewTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForMountVolumeNewTargetConfig, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForMountVolumeNewTargetConfig{}
	MountTargetModel, err := ResourceIbmBackupRecoveryMapToRecoverTarget(modelMap["mount_target"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.MountTarget = MountTargetModel
	if modelMap["server_credentials"] != nil && len(modelMap["server_credentials"].([]interface{})) > 0 {
		ServerCredentialsModel, err := ResourceIbmBackupRecoveryMapToPhysicalMountVolumesNewTargetConfigServerCredentials(modelMap["server_credentials"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ServerCredentials = ServerCredentialsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverTarget(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverTarget, error) {
	model := &backuprecoveryv1.RecoverTarget{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["parent_source_id"] != nil {
		model.ParentSourceID = core.Int64Ptr(int64(modelMap["parent_source_id"].(int)))
	}
	if modelMap["parent_source_name"] != nil && modelMap["parent_source_name"].(string) != "" {
		model.ParentSourceName = core.StringPtr(modelMap["parent_source_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalMountVolumesNewTargetConfigServerCredentials(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalMountVolumesNewTargetConfigServerCredentials, error) {
	model := &backuprecoveryv1.PhysicalMountVolumesNewTargetConfigServerCredentials{}
	model.Username = core.StringPtr(modelMap["username"].(string))
	model.Password = core.StringPtr(modelMap["password"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryMapToMountedVolumeMapping(modelMap map[string]interface{}) (*backuprecoveryv1.MountedVolumeMapping, error) {
	model := &backuprecoveryv1.MountedVolumeMapping{}
	if modelMap["original_volume"] != nil && modelMap["original_volume"].(string) != "" {
		model.OriginalVolume = core.StringPtr(modelMap["original_volume"].(string))
	}
	if modelMap["mounted_volume"] != nil && modelMap["mounted_volume"].(string) != "" {
		model.MountedVolume = core.StringPtr(modelMap["mounted_volume"].(string))
	}
	if modelMap["file_system_type"] != nil && modelMap["file_system_type"].(string) != "" {
		model.FileSystemType = core.StringPtr(modelMap["file_system_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForMountVolumeVlanConfig(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForMountVolumeVlanConfig, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForMountVolumeVlanConfig{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["disable_vlan"] != nil {
		model.DisableVlan = core.BoolPtr(modelMap["disable_vlan"].(bool))
	}
	if modelMap["interface_name"] != nil && modelMap["interface_name"].(string) != "" {
		model.InterfaceName = core.StringPtr(modelMap["interface_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsRecoverFileAndFolderParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalParamsRecoverFileAndFolderParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalParamsRecoverFileAndFolderParams{}
	filesAndFolders := []backuprecoveryv1.CommonRecoverFileAndFolderInfo{}
	for _, filesAndFoldersItem := range modelMap["files_and_folders"].([]interface{}) {
		filesAndFoldersItemModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverFileAndFolderInfo(filesAndFoldersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		filesAndFolders = append(filesAndFolders, *filesAndFoldersItemModel)
	}
	model.FilesAndFolders = filesAndFolders
	model.TargetEnvironment = core.StringPtr(modelMap["target_environment"].(string))
	if modelMap["physical_target_params"] != nil && len(modelMap["physical_target_params"].([]interface{})) > 0 {
		PhysicalTargetParamsModel, err := ResourceIbmBackupRecoveryMapToRecoverPhysicalFileAndFolderParamsPhysicalTargetParams(modelMap["physical_target_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PhysicalTargetParams = PhysicalTargetParamsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToCommonRecoverFileAndFolderInfo(modelMap map[string]interface{}) (*backuprecoveryv1.CommonRecoverFileAndFolderInfo, error) {
	model := &backuprecoveryv1.CommonRecoverFileAndFolderInfo{}
	model.AbsolutePath = core.StringPtr(modelMap["absolute_path"].(string))
	if modelMap["destination_dir"] != nil && modelMap["destination_dir"].(string) != "" {
		model.DestinationDir = core.StringPtr(modelMap["destination_dir"].(string))
	}
	if modelMap["is_directory"] != nil {
		model.IsDirectory = core.BoolPtr(modelMap["is_directory"].(bool))
	}
	if modelMap["status"] != nil && modelMap["status"].(string) != "" {
		model.Status = core.StringPtr(modelMap["status"].(string))
	}
	if modelMap["messages"] != nil {
		messages := []string{}
		for _, messagesItem := range modelMap["messages"].([]interface{}) {
			messages = append(messages, messagesItem.(string))
		}
		model.Messages = messages
	}
	if modelMap["is_view_file_recovery"] != nil {
		model.IsViewFileRecovery = core.BoolPtr(modelMap["is_view_file_recovery"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverPhysicalFileAndFolderParamsPhysicalTargetParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalFileAndFolderParamsPhysicalTargetParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalFileAndFolderParamsPhysicalTargetParams{}
	RecoverTargetModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverFileAndFolderRecoverTarget(modelMap["recover_target"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.RecoverTarget = RecoverTargetModel
	if modelMap["restore_to_original_paths"] != nil {
		model.RestoreToOriginalPaths = core.BoolPtr(modelMap["restore_to_original_paths"].(bool))
	}
	if modelMap["overwrite_existing"] != nil {
		model.OverwriteExisting = core.BoolPtr(modelMap["overwrite_existing"].(bool))
	}
	if modelMap["alternate_restore_directory"] != nil && modelMap["alternate_restore_directory"].(string) != "" {
		model.AlternateRestoreDirectory = core.StringPtr(modelMap["alternate_restore_directory"].(string))
	}
	if modelMap["preserve_attributes"] != nil {
		model.PreserveAttributes = core.BoolPtr(modelMap["preserve_attributes"].(bool))
	}
	if modelMap["preserve_timestamps"] != nil {
		model.PreserveTimestamps = core.BoolPtr(modelMap["preserve_timestamps"].(bool))
	}
	if modelMap["preserve_acls"] != nil {
		model.PreserveAcls = core.BoolPtr(modelMap["preserve_acls"].(bool))
	}
	if modelMap["continue_on_error"] != nil {
		model.ContinueOnError = core.BoolPtr(modelMap["continue_on_error"].(bool))
	}
	if modelMap["save_success_files"] != nil {
		model.SaveSuccessFiles = core.BoolPtr(modelMap["save_success_files"].(bool))
	}
	if modelMap["vlan_config"] != nil && len(modelMap["vlan_config"].([]interface{})) > 0 {
		VlanConfigModel, err := ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverFileAndFolderVlanConfig(modelMap["vlan_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VlanConfig = VlanConfigModel
	}
	if modelMap["restore_entity_type"] != nil && modelMap["restore_entity_type"].(string) != "" {
		model.RestoreEntityType = core.StringPtr(modelMap["restore_entity_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverFileAndFolderRecoverTarget(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderRecoverTarget, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderRecoverTarget{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["parent_source_id"] != nil && modelMap["parent_source_id"].(int) != 0 {
		model.ParentSourceID = core.Int64Ptr(int64(modelMap["parent_source_id"].(int)))
	}
	if modelMap["parent_source_name"] != nil && modelMap["parent_source_name"].(string) != "" {
		model.ParentSourceName = core.StringPtr(modelMap["parent_source_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToPhysicalTargetParamsForRecoverFileAndFolderVlanConfig(modelMap map[string]interface{}) (*backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderVlanConfig, error) {
	model := &backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderVlanConfig{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["disable_vlan"] != nil {
		model.DisableVlan = core.BoolPtr(modelMap["disable_vlan"].(bool))
	}
	if modelMap["interface_name"] != nil && modelMap["interface_name"].(string) != "" {
		model.InterfaceName = core.StringPtr(modelMap["interface_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsDownloadFileAndFolderParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalParamsDownloadFileAndFolderParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalParamsDownloadFileAndFolderParams{}
	if modelMap["expiry_time_usecs"] != nil {
		model.ExpiryTimeUsecs = core.Int64Ptr(int64(modelMap["expiry_time_usecs"].(int)))
	}
	if modelMap["files_and_folders"] != nil {
		filesAndFolders := []backuprecoveryv1.CommonRecoverFileAndFolderInfo{}
		for _, filesAndFoldersItem := range modelMap["files_and_folders"].([]interface{}) {
			filesAndFoldersItemModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverFileAndFolderInfo(filesAndFoldersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			filesAndFolders = append(filesAndFolders, *filesAndFoldersItemModel)
		}
		model.FilesAndFolders = filesAndFolders
	}
	if modelMap["download_file_path"] != nil && modelMap["download_file_path"].(string) != "" {
		model.DownloadFilePath = core.StringPtr(modelMap["download_file_path"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverPhysicalParamsSystemRecoveryParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverPhysicalParamsSystemRecoveryParams, error) {
	model := &backuprecoveryv1.RecoverPhysicalParamsSystemRecoveryParams{}
	if modelMap["full_nas_path"] != nil && modelMap["full_nas_path"].(string) != "" {
		model.FullNasPath = core.StringPtr(modelMap["full_nas_path"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverSqlParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverSqlParams, error) {
	model := &backuprecoveryv1.RecoverSqlParams{}
	if modelMap["recover_app_params"] != nil {
		recoverAppParams := []backuprecoveryv1.RecoverSqlAppParams{}
		for _, recoverAppParamsItem := range modelMap["recover_app_params"].([]interface{}) {
			recoverAppParamsItemModel, err := ResourceIbmBackupRecoveryMapToRecoverSqlAppParams(recoverAppParamsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			recoverAppParams = append(recoverAppParams, *recoverAppParamsItemModel)
		}
		model.RecoverAppParams = recoverAppParams
	}
	model.RecoveryAction = core.StringPtr(modelMap["recovery_action"].(string))
	if modelMap["vlan_config"] != nil && len(modelMap["vlan_config"].([]interface{})) > 0 {
		VlanConfigModel, err := ResourceIbmBackupRecoveryMapToRecoveryVlanConfig(modelMap["vlan_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.VlanConfig = VlanConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverSqlAppParams(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverSqlAppParams, error) {
	model := &backuprecoveryv1.RecoverSqlAppParams{}
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
	if modelMap["snapshot_creation_time_usecs"] != nil {
		model.SnapshotCreationTimeUsecs = core.Int64Ptr(int64(modelMap["snapshot_creation_time_usecs"].(int)))
	}
	if modelMap["object_info"] != nil && len(modelMap["object_info"].([]interface{})) > 0 {
		ObjectInfoModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParamsObjectInfo(modelMap["object_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ObjectInfo = ObjectInfoModel
	}
	if modelMap["snapshot_target_type"] != nil && modelMap["snapshot_target_type"].(string) != "" {
		model.SnapshotTargetType = core.StringPtr(modelMap["snapshot_target_type"].(string))
	}
	if modelMap["archival_target_info"] != nil && len(modelMap["archival_target_info"].([]interface{})) > 0 {
		ArchivalTargetInfoModel, err := ResourceIbmBackupRecoveryMapToCommonRecoverObjectSnapshotParamsArchivalTargetInfo(modelMap["archival_target_info"].([]interface{})[0].(map[string]interface{}))
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
	if modelMap["start_time_usecs"] != nil {
		model.StartTimeUsecs = core.Int64Ptr(int64(modelMap["start_time_usecs"].(int)))
	}
	if modelMap["end_time_usecs"] != nil {
		model.EndTimeUsecs = core.Int64Ptr(int64(modelMap["end_time_usecs"].(int)))
	}
	if modelMap["messages"] != nil {
		messages := []string{}
		for _, messagesItem := range modelMap["messages"].([]interface{}) {
			messages = append(messages, messagesItem.(string))
		}
		model.Messages = messages
	}
	if modelMap["bytes_restored"] != nil {
		model.BytesRestored = core.Int64Ptr(int64(modelMap["bytes_restored"].(int)))
	}
	if modelMap["aag_info"] != nil && len(modelMap["aag_info"].([]interface{})) > 0 {
		AagInfoModel, err := ResourceIbmBackupRecoveryMapToAAGInfo(modelMap["aag_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AagInfo = AagInfoModel
	}
	if modelMap["host_info"] != nil && len(modelMap["host_info"].([]interface{})) > 0 {
		HostInfoModel, err := ResourceIbmBackupRecoveryMapToHostInformation(modelMap["host_info"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.HostInfo = HostInfoModel
	}
	if modelMap["is_encrypted"] != nil {
		model.IsEncrypted = core.BoolPtr(modelMap["is_encrypted"].(bool))
	}
	if modelMap["sql_target_params"] != nil && len(modelMap["sql_target_params"].([]interface{})) > 0 {
		SqlTargetParamsModel, err := ResourceIbmBackupRecoveryMapToSqlTargetParamsForRecoverSqlApp(modelMap["sql_target_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SqlTargetParams = SqlTargetParamsModel
	}
	model.TargetEnvironment = core.StringPtr(modelMap["target_environment"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryMapToAAGInfo(modelMap map[string]interface{}) (*backuprecoveryv1.AAGInfo, error) {
	model := &backuprecoveryv1.AAGInfo{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["object_id"] != nil {
		model.ObjectID = core.Int64Ptr(int64(modelMap["object_id"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToHostInformation(modelMap map[string]interface{}) (*backuprecoveryv1.HostInformation, error) {
	model := &backuprecoveryv1.HostInformation{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToSqlTargetParamsForRecoverSqlApp(modelMap map[string]interface{}) (*backuprecoveryv1.SqlTargetParamsForRecoverSqlApp, error) {
	model := &backuprecoveryv1.SqlTargetParamsForRecoverSqlApp{}
	if modelMap["new_source_config"] != nil && len(modelMap["new_source_config"].([]interface{})) > 0 {
		NewSourceConfigModel, err := ResourceIbmBackupRecoveryMapToRecoverSqlAppNewSourceConfig(modelMap["new_source_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.NewSourceConfig = NewSourceConfigModel
	}
	if modelMap["original_source_config"] != nil && len(modelMap["original_source_config"].([]interface{})) > 0 {
		OriginalSourceConfigModel, err := ResourceIbmBackupRecoveryMapToRecoverSqlAppOriginalSourceConfig(modelMap["original_source_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.OriginalSourceConfig = OriginalSourceConfigModel
	}
	model.RecoverToNewSource = core.BoolPtr(modelMap["recover_to_new_source"].(bool))
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverSqlAppNewSourceConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverSqlAppNewSourceConfig, error) {
	model := &backuprecoveryv1.RecoverSqlAppNewSourceConfig{}
	if modelMap["keep_cdc"] != nil {
		model.KeepCdc = core.BoolPtr(modelMap["keep_cdc"].(bool))
	}
	if modelMap["multi_stage_restore_options"] != nil && len(modelMap["multi_stage_restore_options"].([]interface{})) > 0 {
		MultiStageRestoreOptionsModel, err := ResourceIbmBackupRecoveryMapToMultiStageRestoreOptions(modelMap["multi_stage_restore_options"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MultiStageRestoreOptions = MultiStageRestoreOptionsModel
	}
	if modelMap["native_log_recovery_with_clause"] != nil && modelMap["native_log_recovery_with_clause"].(string) != "" {
		model.NativeLogRecoveryWithClause = core.StringPtr(modelMap["native_log_recovery_with_clause"].(string))
	}
	if modelMap["native_recovery_with_clause"] != nil && modelMap["native_recovery_with_clause"].(string) != "" {
		model.NativeRecoveryWithClause = core.StringPtr(modelMap["native_recovery_with_clause"].(string))
	}
	if modelMap["overwriting_policy"] != nil && modelMap["overwriting_policy"].(string) != "" {
		model.OverwritingPolicy = core.StringPtr(modelMap["overwriting_policy"].(string))
	}
	if modelMap["replay_entire_last_log"] != nil {
		model.ReplayEntireLastLog = core.BoolPtr(modelMap["replay_entire_last_log"].(bool))
	}
	if modelMap["restore_time_usecs"] != nil {
		model.RestoreTimeUsecs = core.Int64Ptr(int64(modelMap["restore_time_usecs"].(int)))
	}
	if modelMap["secondary_data_files_dir_list"] != nil {
		secondaryDataFilesDirList := []backuprecoveryv1.FilenamePatternToDirectory{}
		for _, secondaryDataFilesDirListItem := range modelMap["secondary_data_files_dir_list"].([]interface{}) {
			secondaryDataFilesDirListItemModel, err := ResourceIbmBackupRecoveryMapToFilenamePatternToDirectory(secondaryDataFilesDirListItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			secondaryDataFilesDirList = append(secondaryDataFilesDirList, *secondaryDataFilesDirListItemModel)
		}
		model.SecondaryDataFilesDirList = secondaryDataFilesDirList
	}
	if modelMap["with_no_recovery"] != nil {
		model.WithNoRecovery = core.BoolPtr(modelMap["with_no_recovery"].(bool))
	}
	model.DataFileDirectoryLocation = core.StringPtr(modelMap["data_file_directory_location"].(string))
	if modelMap["database_name"] != nil && modelMap["database_name"].(string) != "" {
		model.DatabaseName = core.StringPtr(modelMap["database_name"].(string))
	}
	HostModel, err := ResourceIbmBackupRecoveryMapToRecoveryObjectIdentifier(modelMap["host"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Host = HostModel
	model.InstanceName = core.StringPtr(modelMap["instance_name"].(string))
	model.LogFileDirectoryLocation = core.StringPtr(modelMap["log_file_directory_location"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryMapToMultiStageRestoreOptions(modelMap map[string]interface{}) (*backuprecoveryv1.MultiStageRestoreOptions, error) {
	model := &backuprecoveryv1.MultiStageRestoreOptions{}
	if modelMap["enable_auto_sync"] != nil {
		model.EnableAutoSync = core.BoolPtr(modelMap["enable_auto_sync"].(bool))
	}
	if modelMap["enable_multi_stage_restore"] != nil {
		model.EnableMultiStageRestore = core.BoolPtr(modelMap["enable_multi_stage_restore"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToFilenamePatternToDirectory(modelMap map[string]interface{}) (*backuprecoveryv1.FilenamePatternToDirectory, error) {
	model := &backuprecoveryv1.FilenamePatternToDirectory{}
	if modelMap["directory"] != nil && modelMap["directory"].(string) != "" {
		model.Directory = core.StringPtr(modelMap["directory"].(string))
	}
	if modelMap["filename_pattern"] != nil && modelMap["filename_pattern"].(string) != "" {
		model.FilenamePattern = core.StringPtr(modelMap["filename_pattern"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoveryObjectIdentifier(modelMap map[string]interface{}) (*backuprecoveryv1.RecoveryObjectIdentifier, error) {
	model := &backuprecoveryv1.RecoveryObjectIdentifier{}
	model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoverSqlAppOriginalSourceConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RecoverSqlAppOriginalSourceConfig, error) {
	model := &backuprecoveryv1.RecoverSqlAppOriginalSourceConfig{}
	if modelMap["keep_cdc"] != nil {
		model.KeepCdc = core.BoolPtr(modelMap["keep_cdc"].(bool))
	}
	if modelMap["multi_stage_restore_options"] != nil && len(modelMap["multi_stage_restore_options"].([]interface{})) > 0 {
		MultiStageRestoreOptionsModel, err := ResourceIbmBackupRecoveryMapToMultiStageRestoreOptions(modelMap["multi_stage_restore_options"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MultiStageRestoreOptions = MultiStageRestoreOptionsModel
	}
	if modelMap["native_log_recovery_with_clause"] != nil && modelMap["native_log_recovery_with_clause"].(string) != "" {
		model.NativeLogRecoveryWithClause = core.StringPtr(modelMap["native_log_recovery_with_clause"].(string))
	}
	if modelMap["native_recovery_with_clause"] != nil && modelMap["native_recovery_with_clause"].(string) != "" {
		model.NativeRecoveryWithClause = core.StringPtr(modelMap["native_recovery_with_clause"].(string))
	}
	if modelMap["overwriting_policy"] != nil && modelMap["overwriting_policy"].(string) != "" {
		model.OverwritingPolicy = core.StringPtr(modelMap["overwriting_policy"].(string))
	}
	if modelMap["replay_entire_last_log"] != nil {
		model.ReplayEntireLastLog = core.BoolPtr(modelMap["replay_entire_last_log"].(bool))
	}
	if modelMap["restore_time_usecs"] != nil {
		model.RestoreTimeUsecs = core.Int64Ptr(int64(modelMap["restore_time_usecs"].(int)))
	}
	if modelMap["secondary_data_files_dir_list"] != nil {
		secondaryDataFilesDirList := []backuprecoveryv1.FilenamePatternToDirectory{}
		for _, secondaryDataFilesDirListItem := range modelMap["secondary_data_files_dir_list"].([]interface{}) {
			secondaryDataFilesDirListItemModel, err := ResourceIbmBackupRecoveryMapToFilenamePatternToDirectory(secondaryDataFilesDirListItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			secondaryDataFilesDirList = append(secondaryDataFilesDirList, *secondaryDataFilesDirListItemModel)
		}
		model.SecondaryDataFilesDirList = secondaryDataFilesDirList
	}
	if modelMap["with_no_recovery"] != nil {
		model.WithNoRecovery = core.BoolPtr(modelMap["with_no_recovery"].(bool))
	}
	if modelMap["capture_tail_logs"] != nil {
		model.CaptureTailLogs = core.BoolPtr(modelMap["capture_tail_logs"].(bool))
	}
	if modelMap["data_file_directory_location"] != nil && modelMap["data_file_directory_location"].(string) != "" {
		model.DataFileDirectoryLocation = core.StringPtr(modelMap["data_file_directory_location"].(string))
	}
	if modelMap["log_file_directory_location"] != nil && modelMap["log_file_directory_location"].(string) != "" {
		model.LogFileDirectoryLocation = core.StringPtr(modelMap["log_file_directory_location"].(string))
	}
	if modelMap["new_database_name"] != nil && modelMap["new_database_name"].(string) != "" {
		model.NewDatabaseName = core.StringPtr(modelMap["new_database_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryMapToRecoveryVlanConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RecoveryVlanConfig, error) {
	model := &backuprecoveryv1.RecoveryVlanConfig{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["disable_vlan"] != nil {
		model.DisableVlan = core.BoolPtr(modelMap["disable_vlan"].(bool))
	}
	if modelMap["interface_name"] != nil && modelMap["interface_name"].(string) != "" {
		model.InterfaceName = core.StringPtr(modelMap["interface_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalParamsToMap(model *backuprecoveryv1.RecoverPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	modelMap["recovery_action"] = *model.RecoveryAction
	if model.RecoverVolumeParams != nil {
		recoverVolumeParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalParamsRecoverVolumeParamsToMap(model.RecoverVolumeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["recover_volume_params"] = []map[string]interface{}{recoverVolumeParamsMap}
	}
	if model.MountVolumeParams != nil {
		mountVolumeParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalParamsMountVolumeParamsToMap(model.MountVolumeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mount_volume_params"] = []map[string]interface{}{mountVolumeParamsMap}
	}
	if model.RecoverFileAndFolderParams != nil {
		recoverFileAndFolderParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalParamsRecoverFileAndFolderParamsToMap(model.RecoverFileAndFolderParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["recover_file_and_folder_params"] = []map[string]interface{}{recoverFileAndFolderParamsMap}
	}
	if model.DownloadFileAndFolderParams != nil {
		downloadFileAndFolderParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalParamsDownloadFileAndFolderParamsToMap(model.DownloadFileAndFolderParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["download_file_and_folder_params"] = []map[string]interface{}{downloadFileAndFolderParamsMap}
	}
	if model.SystemRecoveryParams != nil {
		systemRecoveryParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalParamsSystemRecoveryParamsToMap(model.SystemRecoveryParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["system_recovery_params"] = []map[string]interface{}{systemRecoveryParamsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParams) (map[string]interface{}, error) {
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
		objectInfoMap, err := ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsObjectInfoToMap(model.ObjectInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["object_info"] = []map[string]interface{}{objectInfoMap}
	}
	if model.SnapshotTargetType != nil {
		modelMap["snapshot_target_type"] = *model.SnapshotTargetType
	}
	if model.ArchivalTargetInfo != nil {
		archivalTargetInfoMap, err := ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model.ArchivalTargetInfo)
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
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Messages != nil {
		modelMap["messages"] = model.Messages
	}
	if model.BytesRestored != nil {
		modelMap["bytes_restored"] = flex.IntValue(model.BytesRestored)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsObjectInfoToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParamsObjectInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := ResourceIbmBackupRecoverySharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := ResourceIbmBackupRecoveryObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := ResourceIbmBackupRecoveryObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := ResourceIbmBackupRecoveryObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoverySharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := ResourceIbmBackupRecoverySharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := ResourceIbmBackupRecoveryObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := ResourceIbmBackupRecoveryObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := ResourceIbmBackupRecoveryObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParamsArchivalTargetInfo) (map[string]interface{}, error) {
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
		tierSettingsMap, err := ResourceIbmBackupRecoveryArchivalTargetTierInfoToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryArchivalTargetTierInfoToMap(model *backuprecoveryv1.ArchivalTargetTierInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := ResourceIbmBackupRecoveryAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := ResourceIbmBackupRecoveryAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}
	if model.GoogleTiering != nil {
		googleTieringMap, err := ResourceIbmBackupRecoveryGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := ResourceIbmBackupRecoveryOracleTiersToMap(model.OracleTiering)
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

func ResourceIbmBackupRecoveryAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := ResourceIbmBackupRecoveryAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryRecoverPhysicalParamsRecoverVolumeParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsRecoverVolumeParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_environment"] = *model.TargetEnvironment
	if model.PhysicalTargetParams != nil {
		physicalTargetParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalVolumeParamsPhysicalTargetParamsToMap(model.PhysicalTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_target_params"] = []map[string]interface{}{physicalTargetParamsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalVolumeParamsPhysicalTargetParamsToMap(model *backuprecoveryv1.RecoverPhysicalVolumeParamsPhysicalTargetParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	mountTargetMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeMountTargetToMap(model.MountTarget)
	if err != nil {
		return modelMap, err
	}
	modelMap["mount_target"] = []map[string]interface{}{mountTargetMap}
	volumeMapping := []map[string]interface{}{}
	for _, volumeMappingItem := range model.VolumeMapping {
		volumeMappingItemMap, err := ResourceIbmBackupRecoveryRecoverVolumeMappingToMap(&volumeMappingItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		volumeMapping = append(volumeMapping, volumeMappingItemMap)
	}
	modelMap["volume_mapping"] = volumeMapping
	if model.ForceUnmountVolume != nil {
		modelMap["force_unmount_volume"] = *model.ForceUnmountVolume
	}
	if model.VlanConfig != nil {
		vlanConfigMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeMountTargetToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeMountTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverVolumeMappingToMap(model *backuprecoveryv1.RecoverVolumeMapping) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["source_volume_guid"] = *model.SourceVolumeGuid
	modelMap["destination_volume_guid"] = *model.DestinationVolumeGuid
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeVlanConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeVlanConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.DisableVlan != nil {
		modelMap["disable_vlan"] = *model.DisableVlan
	}
	if model.InterfaceName != nil {
		modelMap["interface_name"] = *model.InterfaceName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalParamsMountVolumeParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsMountVolumeParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_environment"] = *model.TargetEnvironment
	if model.PhysicalTargetParams != nil {
		physicalTargetParamsMap, err := ResourceIbmBackupRecoveryMountPhysicalVolumeParamsPhysicalTargetParamsToMap(model.PhysicalTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_target_params"] = []map[string]interface{}{physicalTargetParamsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryMountPhysicalVolumeParamsPhysicalTargetParamsToMap(model *backuprecoveryv1.MountPhysicalVolumeParamsPhysicalTargetParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_to_original_target"] = *model.MountToOriginalTarget
	if model.OriginalTargetConfig != nil {
		originalTargetConfigMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeOriginalTargetConfigToMap(model.OriginalTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_target_config"] = []map[string]interface{}{originalTargetConfigMap}
	}
	if model.NewTargetConfig != nil {
		newTargetConfigMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeNewTargetConfigToMap(model.NewTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["new_target_config"] = []map[string]interface{}{newTargetConfigMap}
	}
	if model.ReadOnlyMount != nil {
		modelMap["read_only_mount"] = *model.ReadOnlyMount
	}
	if model.VolumeNames != nil {
		modelMap["volume_names"] = model.VolumeNames
	}
	if model.MountedVolumeMapping != nil {
		mountedVolumeMapping := []map[string]interface{}{}
		for _, mountedVolumeMappingItem := range model.MountedVolumeMapping {
			mountedVolumeMappingItemMap, err := ResourceIbmBackupRecoveryMountedVolumeMappingToMap(&mountedVolumeMappingItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			mountedVolumeMapping = append(mountedVolumeMapping, mountedVolumeMappingItemMap)
		}
		modelMap["mounted_volume_mapping"] = mountedVolumeMapping
	}
	if model.VlanConfig != nil {
		vlanConfigMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeOriginalTargetConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForMountVolumeOriginalTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ServerCredentials != nil {
		serverCredentialsMap, err := ResourceIbmBackupRecoveryPhysicalMountVolumesOriginalTargetConfigServerCredentialsToMap(model.ServerCredentials)
		if err != nil {
			return modelMap, err
		}
		modelMap["server_credentials"] = []map[string]interface{}{serverCredentialsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalMountVolumesOriginalTargetConfigServerCredentialsToMap(model *backuprecoveryv1.PhysicalMountVolumesOriginalTargetConfigServerCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["username"] = *model.Username
	modelMap["password"] = *model.Password
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeNewTargetConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForMountVolumeNewTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	mountTargetMap, err := ResourceIbmBackupRecoveryRecoverTargetToMap(model.MountTarget)
	if err != nil {
		return modelMap, err
	}
	modelMap["mount_target"] = []map[string]interface{}{mountTargetMap}
	if model.ServerCredentials != nil {
		serverCredentialsMap, err := ResourceIbmBackupRecoveryPhysicalMountVolumesNewTargetConfigServerCredentialsToMap(model.ServerCredentials)
		if err != nil {
			return modelMap, err
		}
		modelMap["server_credentials"] = []map[string]interface{}{serverCredentialsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverTargetToMap(model *backuprecoveryv1.RecoverTarget) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryPhysicalMountVolumesNewTargetConfigServerCredentialsToMap(model *backuprecoveryv1.PhysicalMountVolumesNewTargetConfigServerCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["username"] = *model.Username
	modelMap["password"] = *model.Password
	return modelMap, nil
}

func ResourceIbmBackupRecoveryMountedVolumeMappingToMap(model *backuprecoveryv1.MountedVolumeMapping) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.OriginalVolume != nil {
		modelMap["original_volume"] = *model.OriginalVolume
	}
	if model.MountedVolume != nil {
		modelMap["mounted_volume"] = *model.MountedVolume
	}
	if model.FileSystemType != nil {
		modelMap["file_system_type"] = *model.FileSystemType
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeVlanConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForMountVolumeVlanConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.DisableVlan != nil {
		modelMap["disable_vlan"] = *model.DisableVlan
	}
	if model.InterfaceName != nil {
		modelMap["interface_name"] = *model.InterfaceName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalParamsRecoverFileAndFolderParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsRecoverFileAndFolderParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	filesAndFolders := []map[string]interface{}{}
	for _, filesAndFoldersItem := range model.FilesAndFolders {
		filesAndFoldersItemMap, err := ResourceIbmBackupRecoveryCommonRecoverFileAndFolderInfoToMap(&filesAndFoldersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		filesAndFolders = append(filesAndFolders, filesAndFoldersItemMap)
	}
	modelMap["files_and_folders"] = filesAndFolders
	modelMap["target_environment"] = *model.TargetEnvironment
	if model.PhysicalTargetParams != nil {
		physicalTargetParamsMap, err := ResourceIbmBackupRecoveryRecoverPhysicalFileAndFolderParamsPhysicalTargetParamsToMap(model.PhysicalTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_target_params"] = []map[string]interface{}{physicalTargetParamsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryCommonRecoverFileAndFolderInfoToMap(model *backuprecoveryv1.CommonRecoverFileAndFolderInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["absolute_path"] = *model.AbsolutePath
	if model.DestinationDir != nil {
		modelMap["destination_dir"] = *model.DestinationDir
	}
	if model.IsDirectory != nil {
		modelMap["is_directory"] = *model.IsDirectory
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.Messages != nil {
		modelMap["messages"] = model.Messages
	}
	if model.IsViewFileRecovery != nil {
		modelMap["is_view_file_recovery"] = *model.IsViewFileRecovery
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalFileAndFolderParamsPhysicalTargetParamsToMap(model *backuprecoveryv1.RecoverPhysicalFileAndFolderParamsPhysicalTargetParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	recoverTargetMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderRecoverTargetToMap(model.RecoverTarget)
	if err != nil {
		return modelMap, err
	}
	modelMap["recover_target"] = []map[string]interface{}{recoverTargetMap}
	if model.RestoreToOriginalPaths != nil {
		modelMap["restore_to_original_paths"] = *model.RestoreToOriginalPaths
	}
	if model.OverwriteExisting != nil {
		modelMap["overwrite_existing"] = *model.OverwriteExisting
	}
	if model.AlternateRestoreDirectory != nil {
		modelMap["alternate_restore_directory"] = *model.AlternateRestoreDirectory
	}
	if model.PreserveAttributes != nil {
		modelMap["preserve_attributes"] = *model.PreserveAttributes
	}
	if model.PreserveTimestamps != nil {
		modelMap["preserve_timestamps"] = *model.PreserveTimestamps
	}
	if model.PreserveAcls != nil {
		modelMap["preserve_acls"] = *model.PreserveAcls
	}
	if model.ContinueOnError != nil {
		modelMap["continue_on_error"] = *model.ContinueOnError
	}
	if model.SaveSuccessFiles != nil {
		modelMap["save_success_files"] = *model.SaveSuccessFiles
	}
	if model.VlanConfig != nil {
		vlanConfigMap, err := ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	if model.RestoreEntityType != nil {
		modelMap["restore_entity_type"] = *model.RestoreEntityType
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderRecoverTargetToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderRecoverTarget) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderVlanConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderVlanConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.DisableVlan != nil {
		modelMap["disable_vlan"] = *model.DisableVlan
	}
	if model.InterfaceName != nil {
		modelMap["interface_name"] = *model.InterfaceName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalParamsDownloadFileAndFolderParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsDownloadFileAndFolderParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.FilesAndFolders != nil {
		filesAndFolders := []map[string]interface{}{}
		for _, filesAndFoldersItem := range model.FilesAndFolders {
			filesAndFoldersItemMap, err := ResourceIbmBackupRecoveryCommonRecoverFileAndFolderInfoToMap(&filesAndFoldersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			filesAndFolders = append(filesAndFolders, filesAndFoldersItemMap)
		}
		modelMap["files_and_folders"] = filesAndFolders
	}
	if model.DownloadFilePath != nil {
		modelMap["download_file_path"] = *model.DownloadFilePath
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverPhysicalParamsSystemRecoveryParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsSystemRecoveryParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FullNasPath != nil {
		modelMap["full_nas_path"] = *model.FullNasPath
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverSqlParamsToMap(model *backuprecoveryv1.RecoverSqlParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RecoverAppParams != nil {
		recoverAppParams := []map[string]interface{}{}
		for _, recoverAppParamsItem := range model.RecoverAppParams {
			recoverAppParamsItemMap, err := ResourceIbmBackupRecoveryRecoverSqlAppParamsToMap(&recoverAppParamsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			recoverAppParams = append(recoverAppParams, recoverAppParamsItemMap)
		}
		modelMap["recover_app_params"] = recoverAppParams
	}
	modelMap["recovery_action"] = *model.RecoveryAction
	if model.VlanConfig != nil {
		vlanConfigMap, err := ResourceIbmBackupRecoveryVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverSqlAppParamsToMap(model *backuprecoveryv1.RecoverSqlAppParams) (map[string]interface{}, error) {
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
	if model.SnapshotCreationTimeUsecs != nil {
		modelMap["snapshot_creation_time_usecs"] = flex.IntValue(model.SnapshotCreationTimeUsecs)
	}
	if model.ObjectInfo != nil {
		objectInfoMap, err := ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsObjectInfoToMap(model.ObjectInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["object_info"] = []map[string]interface{}{objectInfoMap}
	}
	if model.SnapshotTargetType != nil {
		modelMap["snapshot_target_type"] = *model.SnapshotTargetType
	}
	if model.ArchivalTargetInfo != nil {
		archivalTargetInfoMap, err := ResourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model.ArchivalTargetInfo)
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
	if model.StartTimeUsecs != nil {
		modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	}
	if model.EndTimeUsecs != nil {
		modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	}
	if model.Messages != nil {
		modelMap["messages"] = model.Messages
	}
	if model.BytesRestored != nil {
		modelMap["bytes_restored"] = flex.IntValue(model.BytesRestored)
	}
	if model.AagInfo != nil {
		aagInfoMap, err := ResourceIbmBackupRecoveryAAGInfoToMap(model.AagInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["aag_info"] = []map[string]interface{}{aagInfoMap}
	}
	if model.HostInfo != nil {
		hostInfoMap, err := ResourceIbmBackupRecoveryHostInformationToMap(model.HostInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["host_info"] = []map[string]interface{}{hostInfoMap}
	}
	if model.IsEncrypted != nil {
		modelMap["is_encrypted"] = *model.IsEncrypted
	}
	if model.SqlTargetParams != nil {
		sqlTargetParamsMap, err := ResourceIbmBackupRecoverySqlTargetParamsForRecoverSqlAppToMap(model.SqlTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["sql_target_params"] = []map[string]interface{}{sqlTargetParamsMap}
	}
	modelMap["target_environment"] = *model.TargetEnvironment
	return modelMap, nil
}

func ResourceIbmBackupRecoveryAAGInfoToMap(model *backuprecoveryv1.AAGInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryHostInformationToMap(model *backuprecoveryv1.HostInformation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoverySqlTargetParamsForRecoverSqlAppToMap(model *backuprecoveryv1.SqlTargetParamsForRecoverSqlApp) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewSourceConfig != nil {
		newSourceConfigMap, err := ResourceIbmBackupRecoveryRecoverSqlAppNewSourceConfigToMap(model.NewSourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["new_source_config"] = []map[string]interface{}{newSourceConfigMap}
	}
	if model.OriginalSourceConfig != nil {
		originalSourceConfigMap, err := ResourceIbmBackupRecoveryRecoverSqlAppOriginalSourceConfigToMap(model.OriginalSourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_source_config"] = []map[string]interface{}{originalSourceConfigMap}
	}
	modelMap["recover_to_new_source"] = *model.RecoverToNewSource
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverSqlAppNewSourceConfigToMap(model *backuprecoveryv1.RecoverSqlAppNewSourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.KeepCdc != nil {
		modelMap["keep_cdc"] = *model.KeepCdc
	}
	if model.MultiStageRestoreOptions != nil {
		multiStageRestoreOptionsMap, err := ResourceIbmBackupRecoveryMultiStageRestoreOptionsToMap(model.MultiStageRestoreOptions)
		if err != nil {
			return modelMap, err
		}
		modelMap["multi_stage_restore_options"] = []map[string]interface{}{multiStageRestoreOptionsMap}
	}
	if model.NativeLogRecoveryWithClause != nil {
		modelMap["native_log_recovery_with_clause"] = *model.NativeLogRecoveryWithClause
	}
	if model.NativeRecoveryWithClause != nil {
		modelMap["native_recovery_with_clause"] = *model.NativeRecoveryWithClause
	}
	if model.OverwritingPolicy != nil {
		modelMap["overwriting_policy"] = *model.OverwritingPolicy
	}
	if model.ReplayEntireLastLog != nil {
		modelMap["replay_entire_last_log"] = *model.ReplayEntireLastLog
	}
	if model.RestoreTimeUsecs != nil {
		modelMap["restore_time_usecs"] = flex.IntValue(model.RestoreTimeUsecs)
	}
	if model.SecondaryDataFilesDirList != nil {
		secondaryDataFilesDirList := []map[string]interface{}{}
		for _, secondaryDataFilesDirListItem := range model.SecondaryDataFilesDirList {
			secondaryDataFilesDirListItemMap, err := ResourceIbmBackupRecoveryFilenamePatternToDirectoryToMap(&secondaryDataFilesDirListItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			secondaryDataFilesDirList = append(secondaryDataFilesDirList, secondaryDataFilesDirListItemMap)
		}
		modelMap["secondary_data_files_dir_list"] = secondaryDataFilesDirList
	}
	if model.WithNoRecovery != nil {
		modelMap["with_no_recovery"] = *model.WithNoRecovery
	}
	modelMap["data_file_directory_location"] = *model.DataFileDirectoryLocation
	if model.DatabaseName != nil {
		modelMap["database_name"] = *model.DatabaseName
	}
	hostMap, err := ResourceIbmBackupRecoveryObjectIdentifierToMap(model.Host)
	if err != nil {
		return modelMap, err
	}
	modelMap["host"] = []map[string]interface{}{hostMap}
	modelMap["instance_name"] = *model.InstanceName
	modelMap["log_file_directory_location"] = *model.LogFileDirectoryLocation
	return modelMap, nil
}

func ResourceIbmBackupRecoveryMultiStageRestoreOptionsToMap(model *backuprecoveryv1.MultiStageRestoreOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableAutoSync != nil {
		modelMap["enable_auto_sync"] = *model.EnableAutoSync
	}
	if model.EnableMultiStageRestore != nil {
		modelMap["enable_multi_stage_restore"] = *model.EnableMultiStageRestore
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryFilenamePatternToDirectoryToMap(model *backuprecoveryv1.FilenamePatternToDirectory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Directory != nil {
		modelMap["directory"] = *model.Directory
	}
	if model.FilenamePattern != nil {
		modelMap["filename_pattern"] = *model.FilenamePattern
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryObjectIdentifierToMap(model *backuprecoveryv1.RecoveryObjectIdentifier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRecoverSqlAppOriginalSourceConfigToMap(model *backuprecoveryv1.RecoverSqlAppOriginalSourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.KeepCdc != nil {
		modelMap["keep_cdc"] = *model.KeepCdc
	}
	if model.MultiStageRestoreOptions != nil {
		multiStageRestoreOptionsMap, err := ResourceIbmBackupRecoveryMultiStageRestoreOptionsToMap(model.MultiStageRestoreOptions)
		if err != nil {
			return modelMap, err
		}
		modelMap["multi_stage_restore_options"] = []map[string]interface{}{multiStageRestoreOptionsMap}
	}
	if model.NativeLogRecoveryWithClause != nil {
		modelMap["native_log_recovery_with_clause"] = *model.NativeLogRecoveryWithClause
	}
	if model.NativeRecoveryWithClause != nil {
		modelMap["native_recovery_with_clause"] = *model.NativeRecoveryWithClause
	}
	if model.OverwritingPolicy != nil {
		modelMap["overwriting_policy"] = *model.OverwritingPolicy
	}
	if model.ReplayEntireLastLog != nil {
		modelMap["replay_entire_last_log"] = *model.ReplayEntireLastLog
	}
	if model.RestoreTimeUsecs != nil {
		modelMap["restore_time_usecs"] = flex.IntValue(model.RestoreTimeUsecs)
	}
	if model.SecondaryDataFilesDirList != nil {
		secondaryDataFilesDirList := []map[string]interface{}{}
		for _, secondaryDataFilesDirListItem := range model.SecondaryDataFilesDirList {
			secondaryDataFilesDirListItemMap, err := ResourceIbmBackupRecoveryFilenamePatternToDirectoryToMap(&secondaryDataFilesDirListItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			secondaryDataFilesDirList = append(secondaryDataFilesDirList, secondaryDataFilesDirListItemMap)
		}
		modelMap["secondary_data_files_dir_list"] = secondaryDataFilesDirList
	}
	if model.WithNoRecovery != nil {
		modelMap["with_no_recovery"] = *model.WithNoRecovery
	}
	if model.CaptureTailLogs != nil {
		modelMap["capture_tail_logs"] = *model.CaptureTailLogs
	}
	if model.DataFileDirectoryLocation != nil {
		modelMap["data_file_directory_location"] = *model.DataFileDirectoryLocation
	}
	if model.LogFileDirectoryLocation != nil {
		modelMap["log_file_directory_location"] = *model.LogFileDirectoryLocation
	}
	if model.NewDatabaseName != nil {
		modelMap["new_database_name"] = *model.NewDatabaseName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryVlanConfigToMap(model *backuprecoveryv1.RecoveryVlanConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.DisableVlan != nil {
		modelMap["disable_vlan"] = *model.DisableVlan
	}
	if model.InterfaceName != nil {
		modelMap["interface_name"] = *model.InterfaceName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
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
		externalVendorMetadataMap, err := ResourceIbmBackupRecoveryExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
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
		networkMap, err := ResourceIbmBackupRecoveryTenantNetworkToMap(model.Network)
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

func ResourceIbmBackupRecoveryExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := ResourceIbmBackupRecoveryIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func ResourceIbmBackupRecoveryIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
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
			customPropertiesItemMap, err := ResourceIbmBackupRecoveryExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
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
		metricsConfigMap, err := ResourceIbmBackupRecoveryIbmTenantMetricsConfigToMap(model.MetricsConfig)
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

func ResourceIbmBackupRecoveryExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryIbmTenantMetricsConfigToMap(model *backuprecoveryv1.IbmTenantMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CosResourceConfig != nil {
		cosResourceConfigMap, err := ResourceIbmBackupRecoveryIbmTenantCOSResourceConfigToMap(model.CosResourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_resource_config"] = []map[string]interface{}{cosResourceConfigMap}
	}
	if model.IamMetricsConfig != nil {
		iamMetricsConfigMap, err := ResourceIbmBackupRecoveryIbmTenantIAMMetricsConfigToMap(model.IamMetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["iam_metrics_config"] = []map[string]interface{}{iamMetricsConfigMap}
	}
	if model.MeteringConfig != nil {
		meteringConfigMap, err := ResourceIbmBackupRecoveryIbmTenantMeteringConfigToMap(model.MeteringConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metering_config"] = []map[string]interface{}{meteringConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryIbmTenantCOSResourceConfigToMap(model *backuprecoveryv1.IbmTenantCOSResourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceURL != nil {
		modelMap["resource_url"] = *model.ResourceURL
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryIbmTenantIAMMetricsConfigToMap(model *backuprecoveryv1.IbmTenantIAMMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IAMURL != nil {
		modelMap["iam_url"] = *model.IAMURL
	}
	if model.BillingApiKeySecretID != nil {
		modelMap["billing_api_key_secret_id"] = *model.BillingApiKeySecretID
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryIbmTenantMeteringConfigToMap(model *backuprecoveryv1.IbmTenantMeteringConfig) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryCreationInfoToMap(model *backuprecoveryv1.CreationInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryRetrieveArchiveTaskToMap(model *backuprecoveryv1.RetrieveArchiveTask) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TaskUid != nil {
		modelMap["task_uid"] = *model.TaskUid
	}
	if model.UptierExpiryTimes != nil {
		modelMap["uptier_expiry_times"] = model.UptierExpiryTimes
	}
	return modelMap, nil
}
