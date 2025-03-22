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
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func DataSourceIbmBackupRecovery() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryRead,

		Schema: map[string]*schema.Schema{
			"recovery_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the id of a Recovery.",
			},
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the name of the Recovery.",
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
			"snapshot_environment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies the type of snapshot environment for which the Recovery was performed.",
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
			"creation_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the information about the creation of the protection group or recovery.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_name": &schema.Schema{
							Type:        schema.TypeString,
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
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
							Computed:    true,
							Description: "Specifies the globally unique id for this retrieval of an archive task.",
						},
						"uptier_expiry_times": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies how much time the retrieved entity is present in the hot-tiers.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
					},
				},
			},
			"is_multi_stage_restore": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the current recovery operation is a multi-stage restore operation. This is currently used by VMware recoveres for the migration/hot-standby use case.",
			},
			"physical_params": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the recovery options specific to Physical environment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"objects": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the list of Recover Object parameters. For recovering files, specifies the object contains the file to recover.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshot_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the snapshot id.",
									},
									"point_in_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the timestamp (in microseconds. from epoch) for recovering to a point-in-time in the past.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection group id of the object snapshot.",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
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
									"progress_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task id for Recovery of VM.",
									},
									"recover_from_standby": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
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
							Computed:    true,
							Description: "Specifies the type of recover action to be performed.",
						},
						"recover_volume_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the parameters to recover Physical Volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
									"physical_target_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the params for recovering to a physical target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_target": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the target entity where the volumes are being mounted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
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
													Computed:    true,
													Description: "Specifies the mapping from source volumes to destination volumes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source_volume_guid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the guid of the source volume.",
															},
															"destination_volume_guid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the guid of the destination volume.",
															},
														},
													},
												},
												"force_unmount_volume": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether volume would be dismounted first during LockVolume failure. If not specified, default is false.",
												},
												"vlan_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on Cohesity, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on Cohesity, then the partition hostname or VIPs will be used for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
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
							Computed:    true,
							Description: "Specifies the parameters to mount Physical Volumes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
									"physical_target_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the params for recovering to a physical target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mount_to_original_target": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to mount to the original target. If true, originalTargetConfig must be specified. If false, newTargetConfig must be specified.",
												},
												"original_target_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the configuration for mounting to the original target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"server_credentials": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies credentials to access the target server. This is required if the server is of Linux OS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the username to access target entity.",
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
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
													Computed:    true,
													Description: "Specifies the configuration for mounting to a new target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mount_target": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the target entity to recover to.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
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
																Computed:    true,
																Description: "Specifies credentials to access the target server. This is required if the server is of Linux OS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the username to access target entity.",
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
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
													Computed:    true,
													Description: "Specifies whether to perform a read-only mount. Default is false.",
												},
												"volume_names": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the names of volumes that need to be mounted. If this is not specified then all volumes that are part of the source VM will be mounted on the target VM.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"mounted_volume_mapping": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the mapping of original volumes and mounted volumes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"original_volume": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the original volume.",
															},
															"mounted_volume": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the point where the volume is mounted.",
															},
															"file_system_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of the file system of the volume.",
															},
														},
													},
												},
												"vlan_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on Cohesity, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on Cohesity, then the partition hostname or VIPs will be used for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
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
							Computed:    true,
							Description: "Specifies the parameters to perform a file and folder recovery.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"files_and_folders": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the information about the files and folders to be recovered.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"absolute_path": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the absolute path to the file or folder.",
												},
												"destination_dir": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the destination directory where the file/directory was copied.",
												},
												"is_directory": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
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
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"is_view_file_recovery": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specify if the recovery is of type view file/folder.",
												},
											},
										},
									},
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
									"physical_target_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters to recover to a Physical target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"recover_target": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the target entity where the volumes are being mounted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
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
													Computed:    true,
													Description: "If this is true, then files will be restored to original paths.",
												},
												"overwrite_existing": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to overwrite existing file/folder during recovery.",
												},
												"alternate_restore_directory": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the directory path where restore should happen if restore_to_original_paths is set to false.",
												},
												"preserve_attributes": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to preserve file/folder attributes during recovery.",
												},
												"preserve_timestamps": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to preserve the original time stamps.",
												},
												"preserve_acls": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Whether to preserve the ACLs of the original file.",
												},
												"continue_on_error": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to continue recovering other volumes if one of the volumes fails to recover. Default value is false.",
												},
												"save_success_files": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to save success files or not. Default value is false.",
												},
												"vlan_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on Cohesity, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on Cohesity, then the partition hostname or VIPs will be used for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
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
													Computed:    true,
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
							Computed:    true,
							Description: "Specifies the parameters to download files and folders.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"expiry_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the time upto which the download link is available.",
									},
									"files_and_folders": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the info about the files and folders to be recovered.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"absolute_path": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the absolute path to the file or folder.",
												},
												"destination_dir": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the destination directory where the file/directory was copied.",
												},
												"is_directory": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
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
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"is_view_file_recovery": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specify if the recovery is of type view file/folder.",
												},
											},
										},
									},
									"download_file_path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the path location to download the files and folders.",
									},
								},
							},
						},
						"system_recovery_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the parameters to perform a system recovery.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"full_nas_path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the path to the recovery view.",
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
				Description: "Specifies the recovery options specific to Sql environment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recover_app_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the parameters to recover Sql databases.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"snapshot_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the snapshot id.",
									},
									"point_in_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the timestamp (in microseconds. from epoch) for recovering to a point-in-time in the past.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection group id of the object snapshot.",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
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
									"progress_task_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Progress monitor task id for Recovery of VM.",
									},
									"recover_from_standby": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
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
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"bytes_restored": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specify the total bytes restored.",
									},
									"aag_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Object details for Mssql.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the AAG name.",
												},
												"object_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the AAG object Id.",
												},
											},
										},
									},
									"host_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the host information for a objects. This is mainly populated in case of App objects where app object is hosted by another object such as VM or physical server.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the id of the host object.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the host object.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment of the object.",
												},
											},
										},
									},
									"is_encrypted": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether the database is TDE enabled.",
									},
									"sql_target_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the params for recovering to a sql host. Specifiy seperate settings for each db object that need to be recovered. Provided sql backup should be recovered to same type of target host. For Example: If you have sql backup taken from a physical host then that should be recovered to physical host only.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"new_source_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the destination Source configuration parameters where the databases will be recovered. This is mandatory if recoverToNewSource is set to true.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keep_cdc": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether to keep CDC (Change Data Capture) on recovered databases or not. If not passed, this is assumed to be true. If withNoRecovery is passed as true, then this field must not be set to true. Passing this field as true in this scenario will be a invalid request.",
															},
															"multi_stage_restore_options": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the parameters related to multi stage Sql restore.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enable_auto_sync": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Set this to true if you want to enable auto sync for multi stage restore.",
																		},
																		"enable_multi_stage_restore": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Set this to true if you are creating a multi-stage Sql restore task needed for features such as Hot-Standby.",
																		},
																	},
																},
															},
															"native_log_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the WITH clause to be used in native sql log restore command. This is only applicable for native log restore.",
															},
															"native_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "'with_clause' contains 'with clause' to be used in native sql restore command. This is only applicable for database restore of native sql backup. Here user can specify multiple restore options. Example: 'WITH BUFFERCOUNT = 575, MAXTRANSFERSIZE = 2097152'.",
															},
															"overwriting_policy": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a policy to be used while recovering existing databases.",
															},
															"replay_entire_last_log": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the option to set replay last log bit while creating the sql restore task and doing restore to latest point-in-time. If this is set to true, we will replay the entire last log without STOPAT.",
															},
															"restore_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time in the past to which the Sql database needs to be restored. This allows for granular recovery of Sql databases. If this is not set, the Sql database will be restored from the full/incremental snapshot.",
															},
															"secondary_data_files_dir_list": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the secondary data filename pattern and corresponding direcories of the DB. Secondary data files are optional and are user defined. The recommended file extention for secondary files is \".ndf\". If this option is specified and the destination folders do not exist they will be automatically created.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"directory": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the directory where to keep the files matching the pattern.",
																		},
																		"filename_pattern": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a pattern to be matched with filenames. This can be a regex expression.",
																		},
																	},
																},
															},
															"with_no_recovery": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the flag to bring DBs online or not after successful recovery. If this is passed as true, then it means DBs won't be brought online.",
															},
															"data_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the directory where to put the database data files. Missing directory will be automatically created.",
															},
															"database_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a new name for the restored database. If this field is not specified, then the original database will be overwritten after recovery.",
															},
															"host": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the source id of target host where databases will be recovered. This source id can be a physical host or virtual machine.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
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
																Computed:    true,
																Description: "Specifies an instance name of the Sql Server that should be used for restoring databases to.",
															},
															"log_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the directory where to put the database log files. Missing directory will be automatically created.",
															},
														},
													},
												},
												"original_source_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the Source configuration if databases are being recovered to Original Source. If not specified, all the configuration parameters will be retained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"keep_cdc": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether to keep CDC (Change Data Capture) on recovered databases or not. If not passed, this is assumed to be true. If withNoRecovery is passed as true, then this field must not be set to true. Passing this field as true in this scenario will be a invalid request.",
															},
															"multi_stage_restore_options": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the parameters related to multi stage Sql restore.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enable_auto_sync": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Set this to true if you want to enable auto sync for multi stage restore.",
																		},
																		"enable_multi_stage_restore": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Set this to true if you are creating a multi-stage Sql restore task needed for features such as Hot-Standby.",
																		},
																	},
																},
															},
															"native_log_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the WITH clause to be used in native sql log restore command. This is only applicable for native log restore.",
															},
															"native_recovery_with_clause": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "'with_clause' contains 'with clause' to be used in native sql restore command. This is only applicable for database restore of native sql backup. Here user can specify multiple restore options. Example: 'WITH BUFFERCOUNT = 575, MAXTRANSFERSIZE = 2097152'.",
															},
															"overwriting_policy": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a policy to be used while recovering existing databases.",
															},
															"replay_entire_last_log": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the option to set replay last log bit while creating the sql restore task and doing restore to latest point-in-time. If this is set to true, we will replay the entire last log without STOPAT.",
															},
															"restore_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the time in the past to which the Sql database needs to be restored. This allows for granular recovery of Sql databases. If this is not set, the Sql database will be restored from the full/incremental snapshot.",
															},
															"secondary_data_files_dir_list": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the secondary data filename pattern and corresponding direcories of the DB. Secondary data files are optional and are user defined. The recommended file extention for secondary files is \".ndf\". If this option is specified and the destination folders do not exist they will be automatically created.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"directory": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the directory where to keep the files matching the pattern.",
																		},
																		"filename_pattern": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a pattern to be matched with filenames. This can be a regex expression.",
																		},
																	},
																},
															},
															"with_no_recovery": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the flag to bring DBs online or not after successful recovery. If this is passed as true, then it means DBs won't be brought online.",
															},
															"capture_tail_logs": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Set this to true if tail logs are to be captured before the recovery operation. This is only applicable if database is not being renamed.",
															},
															"data_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the directory where to put the database data files. Missing directory will be automatically created. If you are overwriting the existing database then this field will be ignored.",
															},
															"log_file_directory_location": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the directory where to put the database log files. Missing directory will be automatically created. If you are overwriting the existing database then this field will be ignored.",
															},
															"new_database_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a new name for the restored database. If this field is not specified, then the original database will be overwritten after recovery.",
															},
														},
													},
												},
												"recover_to_new_source": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies the parameter whether the recovery should be performed to a new sources or an original Source Target.",
												},
											},
										},
									},
									"target_environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment of the recovery target. The corresponding params below must be filled out.",
									},
								},
							},
						},
						"recovery_action": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the type of recover action to be performed.",
						},
						"vlan_config": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies VLAN Params associated with the recovered. If this is not specified, then the VLAN settings will be automatically selected from one of the below options: a. If VLANs are configured on IBM, then the VLAN host/VIP will be automatically based on the client's (e.g. ESXI host) IP address. b. If VLANs are not configured on IBM, then the partition hostname or VIPs will be used for Recovery.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "If this is set, then the Cohesity host name or the IP address associated with this vlan is used for mounting Cohesity's view on the remote host.",
									},
									"disable_vlan": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
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
	}
}

func dataSourceIbmBackupRecoveryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_recovery", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getRecoveryByIdOptions := &backuprecoveryv1.GetRecoveryByIdOptions{}
	tenantId := d.Get("x_ibm_tenant_id").(string)
	getRecoveryByIdOptions.SetID(d.Get("recovery_id").(string))
	getRecoveryByIdOptions.SetXIBMTenantID(tenantId)

	recovery, _, err := backupRecoveryClient.GetRecoveryByIDWithContext(context, getRecoveryByIdOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRecoveryByIDWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_recovery", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	recoveryId := fmt.Sprintf("%s::%s", tenantId, d.Get("recovery_id").(string))
	d.SetId(recoveryId)

	if !core.IsNil(recovery.Name) {
		if err = d.Set("name", recovery.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-name").GetDiag()
		}
	}

	if !core.IsNil(recovery.StartTimeUsecs) {
		if err = d.Set("start_time_usecs", flex.IntValue(recovery.StartTimeUsecs)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting start_time_usecs: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-start_time_usecs").GetDiag()
		}
	}

	if !core.IsNil(recovery.EndTimeUsecs) {
		if err = d.Set("end_time_usecs", flex.IntValue(recovery.EndTimeUsecs)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting end_time_usecs: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-end_time_usecs").GetDiag()
		}
	}

	if !core.IsNil(recovery.Status) {
		if err = d.Set("status", recovery.Status); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-status").GetDiag()
		}
	}

	if !core.IsNil(recovery.ProgressTaskID) {
		if err = d.Set("progress_task_id", recovery.ProgressTaskID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting progress_task_id: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-progress_task_id").GetDiag()
		}
	}

	if !core.IsNil(recovery.SnapshotEnvironment) {
		if err = d.Set("snapshot_environment", recovery.SnapshotEnvironment); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting snapshot_environment: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-snapshot_environment").GetDiag()
		}
	}

	if !core.IsNil(recovery.RecoveryAction) {
		if err = d.Set("recovery_action", recovery.RecoveryAction); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting recovery_action: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-recovery_action").GetDiag()
		}
	}

	if !core.IsNil(recovery.Permissions) {
		permissions := []map[string]interface{}{}
		for _, permissionsItem := range recovery.Permissions {
			permissionsItemMap, err := DataSourceIbmBackupRecoveryTenantToMap(&permissionsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_recovery", "read", "permissions-to-map").GetDiag()
			}
			permissions = append(permissions, permissionsItemMap)
		}
		if err = d.Set("permissions", permissions); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting permissions: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-permissions").GetDiag()
		}
	}

	if !core.IsNil(recovery.CreationInfo) {
		creationInfo := []map[string]interface{}{}
		creationInfoMap, err := DataSourceIbmBackupRecoveryCreationInfoToMap(recovery.CreationInfo)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_recovery", "read", "creation_info-to-map").GetDiag()
		}
		creationInfo = append(creationInfo, creationInfoMap)
		if err = d.Set("creation_info", creationInfo); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting creation_info: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-creation_info").GetDiag()
		}
	}

	if !core.IsNil(recovery.CanTearDown) {
		if err = d.Set("can_tear_down", recovery.CanTearDown); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting can_tear_down: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-can_tear_down").GetDiag()
		}
	}

	if !core.IsNil(recovery.TearDownStatus) {
		if err = d.Set("tear_down_status", recovery.TearDownStatus); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tear_down_status: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-tear_down_status").GetDiag()
		}
	}

	if !core.IsNil(recovery.TearDownMessage) {
		if err = d.Set("tear_down_message", recovery.TearDownMessage); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting tear_down_message: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-tear_down_message").GetDiag()
		}
	}

	if !core.IsNil(recovery.Messages) {
		messages := []interface{}{}
		for _, messagesItem := range recovery.Messages {
			messages = append(messages, messagesItem)
		}
		if err = d.Set("messages", messages); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting messages: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-messages").GetDiag()
		}
	}

	if !core.IsNil(recovery.IsParentRecovery) {
		if err = d.Set("is_parent_recovery", recovery.IsParentRecovery); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting is_parent_recovery: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-is_parent_recovery").GetDiag()
		}
	}

	if !core.IsNil(recovery.ParentRecoveryID) {
		if err = d.Set("parent_recovery_id", recovery.ParentRecoveryID); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting parent_recovery_id: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-parent_recovery_id").GetDiag()
		}
	}

	if !core.IsNil(recovery.RetrieveArchiveTasks) {
		retrieveArchiveTasks := []map[string]interface{}{}
		for _, retrieveArchiveTasksItem := range recovery.RetrieveArchiveTasks {
			retrieveArchiveTasksItemMap, err := DataSourceIbmBackupRecoveryRetrieveArchiveTaskToMap(&retrieveArchiveTasksItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_recovery", "read", "retrieve_archive_tasks-to-map").GetDiag()
			}
			retrieveArchiveTasks = append(retrieveArchiveTasks, retrieveArchiveTasksItemMap)
		}
		if err = d.Set("retrieve_archive_tasks", retrieveArchiveTasks); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting retrieve_archive_tasks: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-retrieve_archive_tasks").GetDiag()
		}
	}

	if !core.IsNil(recovery.IsMultiStageRestore) {
		if err = d.Set("is_multi_stage_restore", recovery.IsMultiStageRestore); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting is_multi_stage_restore: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-is_multi_stage_restore").GetDiag()
		}
	}

	if !core.IsNil(recovery.PhysicalParams) {
		physicalParams := []map[string]interface{}{}
		physicalParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsToMap(recovery.PhysicalParams)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_recovery", "read", "physical_params-to-map").GetDiag()
		}
		physicalParams = append(physicalParams, physicalParamsMap)
		if err = d.Set("physical_params", physicalParams); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting physical_params: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-physical_params").GetDiag()
		}
	}

	if !core.IsNil(recovery.MssqlParams) {
		mssqlParams := []map[string]interface{}{}
		mssqlParamsMap, err := DataSourceIbmBackupRecoveryRecoverSqlParamsToMap(recovery.MssqlParams)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_recovery", "read", "mssql_params-to-map").GetDiag()
		}
		mssqlParams = append(mssqlParams, mssqlParamsMap)
		if err = d.Set("mssql_params", mssqlParams); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting mssql_params: %s", err), "(Data) ibm_backup_recovery_recovery", "read", "set-mssql_params").GetDiag()
		}
	}

	return nil
}

func DataSourceIbmBackupRecoveryTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
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
		externalVendorMetadataMap, err := DataSourceIbmBackupRecoveryExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
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
		networkMap, err := DataSourceIbmBackupRecoveryTenantNetworkToMap(model.Network)
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

func DataSourceIbmBackupRecoveryExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := DataSourceIbmBackupRecoveryIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
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
			customPropertiesItemMap, err := DataSourceIbmBackupRecoveryExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
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
		metricsConfigMap, err := DataSourceIbmBackupRecoveryIbmTenantMetricsConfigToMap(model.MetricsConfig)
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

func DataSourceIbmBackupRecoveryExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryIbmTenantMetricsConfigToMap(model *backuprecoveryv1.IbmTenantMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CosResourceConfig != nil {
		cosResourceConfigMap, err := DataSourceIbmBackupRecoveryIbmTenantCOSResourceConfigToMap(model.CosResourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_resource_config"] = []map[string]interface{}{cosResourceConfigMap}
	}
	if model.IamMetricsConfig != nil {
		iamMetricsConfigMap, err := DataSourceIbmBackupRecoveryIbmTenantIAMMetricsConfigToMap(model.IamMetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["iam_metrics_config"] = []map[string]interface{}{iamMetricsConfigMap}
	}
	if model.MeteringConfig != nil {
		meteringConfigMap, err := DataSourceIbmBackupRecoveryIbmTenantMeteringConfigToMap(model.MeteringConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metering_config"] = []map[string]interface{}{meteringConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryIbmTenantCOSResourceConfigToMap(model *backuprecoveryv1.IbmTenantCOSResourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceURL != nil {
		modelMap["resource_url"] = *model.ResourceURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryIbmTenantIAMMetricsConfigToMap(model *backuprecoveryv1.IbmTenantIAMMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IAMURL != nil {
		modelMap["iam_url"] = *model.IAMURL
	}
	if model.BillingApiKeySecretID != nil {
		modelMap["billing_api_key_secret_id"] = *model.BillingApiKeySecretID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryIbmTenantMeteringConfigToMap(model *backuprecoveryv1.IbmTenantMeteringConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryCreationInfoToMap(model *backuprecoveryv1.CreationInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRetrieveArchiveTaskToMap(model *backuprecoveryv1.RetrieveArchiveTask) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TaskUid != nil {
		modelMap["task_uid"] = *model.TaskUid
	}
	if model.UptierExpiryTimes != nil {
		modelMap["uptier_expiry_times"] = model.UptierExpiryTimes
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverPhysicalParamsToMap(model *backuprecoveryv1.RecoverPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	objects := []map[string]interface{}{}
	for _, objectsItem := range model.Objects {
		objectsItemMap, err := DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsToMap(&objectsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		objects = append(objects, objectsItemMap)
	}
	modelMap["objects"] = objects
	modelMap["recovery_action"] = *model.RecoveryAction
	if model.RecoverVolumeParams != nil {
		recoverVolumeParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsRecoverVolumeParamsToMap(model.RecoverVolumeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["recover_volume_params"] = []map[string]interface{}{recoverVolumeParamsMap}
	}
	if model.MountVolumeParams != nil {
		mountVolumeParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsMountVolumeParamsToMap(model.MountVolumeParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mount_volume_params"] = []map[string]interface{}{mountVolumeParamsMap}
	}
	if model.RecoverFileAndFolderParams != nil {
		recoverFileAndFolderParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsRecoverFileAndFolderParamsToMap(model.RecoverFileAndFolderParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["recover_file_and_folder_params"] = []map[string]interface{}{recoverFileAndFolderParamsMap}
	}
	if model.DownloadFileAndFolderParams != nil {
		downloadFileAndFolderParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsDownloadFileAndFolderParamsToMap(model.DownloadFileAndFolderParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["download_file_and_folder_params"] = []map[string]interface{}{downloadFileAndFolderParamsMap}
	}
	if model.SystemRecoveryParams != nil {
		systemRecoveryParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalParamsSystemRecoveryParamsToMap(model.SystemRecoveryParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["system_recovery_params"] = []map[string]interface{}{systemRecoveryParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParams) (map[string]interface{}, error) {
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
		objectInfoMap, err := DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsObjectInfoToMap(model.ObjectInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["object_info"] = []map[string]interface{}{objectInfoMap}
	}
	if model.SnapshotTargetType != nil {
		modelMap["snapshot_target_type"] = *model.SnapshotTargetType
	}
	if model.ArchivalTargetInfo != nil {
		archivalTargetInfoMap, err := DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model.ArchivalTargetInfo)
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

func DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsObjectInfoToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParamsObjectInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoveryObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoveryObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoveryObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoveryObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoveryObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoveryObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model *backuprecoveryv1.CommonRecoverObjectSnapshotParamsArchivalTargetInfo) (map[string]interface{}, error) {
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
		tierSettingsMap, err := DataSourceIbmBackupRecoveryArchivalTargetTierInfoToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryArchivalTargetTierInfoToMap(model *backuprecoveryv1.ArchivalTargetTierInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := DataSourceIbmBackupRecoveryAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := DataSourceIbmBackupRecoveryAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}
	if model.GoogleTiering != nil {
		googleTieringMap, err := DataSourceIbmBackupRecoveryGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := DataSourceIbmBackupRecoveryOracleTiersToMap(model.OracleTiering)
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

func DataSourceIbmBackupRecoveryAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := DataSourceIbmBackupRecoveryAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryRecoverPhysicalParamsRecoverVolumeParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsRecoverVolumeParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_environment"] = *model.TargetEnvironment
	if model.PhysicalTargetParams != nil {
		physicalTargetParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalVolumeParamsPhysicalTargetParamsToMap(model.PhysicalTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_target_params"] = []map[string]interface{}{physicalTargetParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverPhysicalVolumeParamsPhysicalTargetParamsToMap(model *backuprecoveryv1.RecoverPhysicalVolumeParamsPhysicalTargetParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	mountTargetMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeMountTargetToMap(model.MountTarget)
	if err != nil {
		return modelMap, err
	}
	modelMap["mount_target"] = []map[string]interface{}{mountTargetMap}
	volumeMapping := []map[string]interface{}{}
	for _, volumeMappingItem := range model.VolumeMapping {
		volumeMappingItemMap, err := DataSourceIbmBackupRecoveryRecoverVolumeMappingToMap(&volumeMappingItem) // #nosec G601
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
		vlanConfigMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeMountTargetToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeMountTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverVolumeMappingToMap(model *backuprecoveryv1.RecoverVolumeMapping) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["source_volume_guid"] = *model.SourceVolumeGuid
	modelMap["destination_volume_guid"] = *model.DestinationVolumeGuid
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverVolumeVlanConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverVolumeVlanConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryRecoverPhysicalParamsMountVolumeParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsMountVolumeParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_environment"] = *model.TargetEnvironment
	if model.PhysicalTargetParams != nil {
		physicalTargetParamsMap, err := DataSourceIbmBackupRecoveryMountPhysicalVolumeParamsPhysicalTargetParamsToMap(model.PhysicalTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_target_params"] = []map[string]interface{}{physicalTargetParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryMountPhysicalVolumeParamsPhysicalTargetParamsToMap(model *backuprecoveryv1.MountPhysicalVolumeParamsPhysicalTargetParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_to_original_target"] = *model.MountToOriginalTarget
	if model.OriginalTargetConfig != nil {
		originalTargetConfigMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeOriginalTargetConfigToMap(model.OriginalTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_target_config"] = []map[string]interface{}{originalTargetConfigMap}
	}
	if model.NewTargetConfig != nil {
		newTargetConfigMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeNewTargetConfigToMap(model.NewTargetConfig)
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
			mountedVolumeMappingItemMap, err := DataSourceIbmBackupRecoveryMountedVolumeMappingToMap(&mountedVolumeMappingItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			mountedVolumeMapping = append(mountedVolumeMapping, mountedVolumeMappingItemMap)
		}
		modelMap["mounted_volume_mapping"] = mountedVolumeMapping
	}
	if model.VlanConfig != nil {
		vlanConfigMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeOriginalTargetConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForMountVolumeOriginalTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ServerCredentials != nil {
		serverCredentialsMap, err := DataSourceIbmBackupRecoveryPhysicalMountVolumesOriginalTargetConfigServerCredentialsToMap(model.ServerCredentials)
		if err != nil {
			return modelMap, err
		}
		modelMap["server_credentials"] = []map[string]interface{}{serverCredentialsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryPhysicalMountVolumesOriginalTargetConfigServerCredentialsToMap(model *backuprecoveryv1.PhysicalMountVolumesOriginalTargetConfigServerCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["username"] = *model.Username
	modelMap["password"] = *model.Password
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeNewTargetConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForMountVolumeNewTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	mountTargetMap, err := DataSourceIbmBackupRecoveryRecoverTargetToMap(model.MountTarget)
	if err != nil {
		return modelMap, err
	}
	modelMap["mount_target"] = []map[string]interface{}{mountTargetMap}
	if model.ServerCredentials != nil {
		serverCredentialsMap, err := DataSourceIbmBackupRecoveryPhysicalMountVolumesNewTargetConfigServerCredentialsToMap(model.ServerCredentials)
		if err != nil {
			return modelMap, err
		}
		modelMap["server_credentials"] = []map[string]interface{}{serverCredentialsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverTargetToMap(model *backuprecoveryv1.RecoverTarget) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryPhysicalMountVolumesNewTargetConfigServerCredentialsToMap(model *backuprecoveryv1.PhysicalMountVolumesNewTargetConfigServerCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["username"] = *model.Username
	modelMap["password"] = *model.Password
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryMountedVolumeMappingToMap(model *backuprecoveryv1.MountedVolumeMapping) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForMountVolumeVlanConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForMountVolumeVlanConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryRecoverPhysicalParamsRecoverFileAndFolderParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsRecoverFileAndFolderParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	filesAndFolders := []map[string]interface{}{}
	for _, filesAndFoldersItem := range model.FilesAndFolders {
		filesAndFoldersItemMap, err := DataSourceIbmBackupRecoveryCommonRecoverFileAndFolderInfoToMap(&filesAndFoldersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		filesAndFolders = append(filesAndFolders, filesAndFoldersItemMap)
	}
	modelMap["files_and_folders"] = filesAndFolders
	modelMap["target_environment"] = *model.TargetEnvironment
	if model.PhysicalTargetParams != nil {
		physicalTargetParamsMap, err := DataSourceIbmBackupRecoveryRecoverPhysicalFileAndFolderParamsPhysicalTargetParamsToMap(model.PhysicalTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_target_params"] = []map[string]interface{}{physicalTargetParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryCommonRecoverFileAndFolderInfoToMap(model *backuprecoveryv1.CommonRecoverFileAndFolderInfo) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryRecoverPhysicalFileAndFolderParamsPhysicalTargetParamsToMap(model *backuprecoveryv1.RecoverPhysicalFileAndFolderParamsPhysicalTargetParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	recoverTargetMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderRecoverTargetToMap(model.RecoverTarget)
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
		vlanConfigMap, err := DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderVlanConfigToMap(model.VlanConfig)
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

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderRecoverTargetToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderRecoverTarget) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryPhysicalTargetParamsForRecoverFileAndFolderVlanConfigToMap(model *backuprecoveryv1.PhysicalTargetParamsForRecoverFileAndFolderVlanConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryRecoverPhysicalParamsDownloadFileAndFolderParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsDownloadFileAndFolderParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ExpiryTimeUsecs != nil {
		modelMap["expiry_time_usecs"] = flex.IntValue(model.ExpiryTimeUsecs)
	}
	if model.FilesAndFolders != nil {
		filesAndFolders := []map[string]interface{}{}
		for _, filesAndFoldersItem := range model.FilesAndFolders {
			filesAndFoldersItemMap, err := DataSourceIbmBackupRecoveryCommonRecoverFileAndFolderInfoToMap(&filesAndFoldersItem) // #nosec G601
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

func DataSourceIbmBackupRecoveryRecoverPhysicalParamsSystemRecoveryParamsToMap(model *backuprecoveryv1.RecoverPhysicalParamsSystemRecoveryParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FullNasPath != nil {
		modelMap["full_nas_path"] = *model.FullNasPath
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverSqlParamsToMap(model *backuprecoveryv1.RecoverSqlParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.RecoverAppParams != nil {
		recoverAppParams := []map[string]interface{}{}
		for _, recoverAppParamsItem := range model.RecoverAppParams {
			recoverAppParamsItemMap, err := DataSourceIbmBackupRecoveryRecoverSqlAppParamsToMap(&recoverAppParamsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			recoverAppParams = append(recoverAppParams, recoverAppParamsItemMap)
		}
		modelMap["recover_app_params"] = recoverAppParams
	}
	modelMap["recovery_action"] = *model.RecoveryAction
	if model.VlanConfig != nil {
		vlanConfigMap, err := DataSourceIbmBackupRecoveryVlanConfigToMap(model.VlanConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_config"] = []map[string]interface{}{vlanConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverSqlAppParamsToMap(model *backuprecoveryv1.RecoverSqlAppParams) (map[string]interface{}, error) {
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
		objectInfoMap, err := DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsObjectInfoToMap(model.ObjectInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["object_info"] = []map[string]interface{}{objectInfoMap}
	}
	if model.SnapshotTargetType != nil {
		modelMap["snapshot_target_type"] = *model.SnapshotTargetType
	}
	if model.ArchivalTargetInfo != nil {
		archivalTargetInfoMap, err := DataSourceIbmBackupRecoveryCommonRecoverObjectSnapshotParamsArchivalTargetInfoToMap(model.ArchivalTargetInfo)
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
		aagInfoMap, err := DataSourceIbmBackupRecoveryAAGInfoToMap(model.AagInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["aag_info"] = []map[string]interface{}{aagInfoMap}
	}
	if model.HostInfo != nil {
		hostInfoMap, err := DataSourceIbmBackupRecoveryHostInformationToMap(model.HostInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["host_info"] = []map[string]interface{}{hostInfoMap}
	}
	if model.IsEncrypted != nil {
		modelMap["is_encrypted"] = *model.IsEncrypted
	}
	if model.SqlTargetParams != nil {
		sqlTargetParamsMap, err := DataSourceIbmBackupRecoverySqlTargetParamsForRecoverSqlAppToMap(model.SqlTargetParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["sql_target_params"] = []map[string]interface{}{sqlTargetParamsMap}
	}
	modelMap["target_environment"] = *model.TargetEnvironment
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryAAGInfoToMap(model *backuprecoveryv1.AAGInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryHostInformationToMap(model *backuprecoveryv1.HostInformation) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySqlTargetParamsForRecoverSqlAppToMap(model *backuprecoveryv1.SqlTargetParamsForRecoverSqlApp) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.NewSourceConfig != nil {
		newSourceConfigMap, err := DataSourceIbmBackupRecoveryRecoverSqlAppNewSourceConfigToMap(model.NewSourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["new_source_config"] = []map[string]interface{}{newSourceConfigMap}
	}
	if model.OriginalSourceConfig != nil {
		originalSourceConfigMap, err := DataSourceIbmBackupRecoveryRecoverSqlAppOriginalSourceConfigToMap(model.OriginalSourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["original_source_config"] = []map[string]interface{}{originalSourceConfigMap}
	}
	modelMap["recover_to_new_source"] = *model.RecoverToNewSource
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverSqlAppNewSourceConfigToMap(model *backuprecoveryv1.RecoverSqlAppNewSourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.KeepCdc != nil {
		modelMap["keep_cdc"] = *model.KeepCdc
	}
	if model.MultiStageRestoreOptions != nil {
		multiStageRestoreOptionsMap, err := DataSourceIbmBackupRecoveryMultiStageRestoreOptionsToMap(model.MultiStageRestoreOptions)
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
			secondaryDataFilesDirListItemMap, err := DataSourceIbmBackupRecoveryFilenamePatternToDirectoryToMap(&secondaryDataFilesDirListItem) // #nosec G601
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
	hostMap, err := DataSourceIbmBackupRecoveryObjectIdentifierToMap(model.Host)
	if err != nil {
		return modelMap, err
	}
	modelMap["host"] = []map[string]interface{}{hostMap}
	modelMap["instance_name"] = *model.InstanceName
	modelMap["log_file_directory_location"] = *model.LogFileDirectoryLocation
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryMultiStageRestoreOptionsToMap(model *backuprecoveryv1.MultiStageRestoreOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableAutoSync != nil {
		modelMap["enable_auto_sync"] = *model.EnableAutoSync
	}
	if model.EnableMultiStageRestore != nil {
		modelMap["enable_multi_stage_restore"] = *model.EnableMultiStageRestore
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryFilenamePatternToDirectoryToMap(model *backuprecoveryv1.FilenamePatternToDirectory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Directory != nil {
		modelMap["directory"] = *model.Directory
	}
	if model.FilenamePattern != nil {
		modelMap["filename_pattern"] = *model.FilenamePattern
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryObjectIdentifierToMap(model *backuprecoveryv1.RecoveryObjectIdentifier) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = flex.IntValue(model.ID)
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryRecoverSqlAppOriginalSourceConfigToMap(model *backuprecoveryv1.RecoverSqlAppOriginalSourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.KeepCdc != nil {
		modelMap["keep_cdc"] = *model.KeepCdc
	}
	if model.MultiStageRestoreOptions != nil {
		multiStageRestoreOptionsMap, err := DataSourceIbmBackupRecoveryMultiStageRestoreOptionsToMap(model.MultiStageRestoreOptions)
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
			secondaryDataFilesDirListItemMap, err := DataSourceIbmBackupRecoveryFilenamePatternToDirectoryToMap(&secondaryDataFilesDirListItem) // #nosec G601
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

func DataSourceIbmBackupRecoveryVlanConfigToMap(model *backuprecoveryv1.RecoveryVlanConfig) (map[string]interface{}, error) {
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
