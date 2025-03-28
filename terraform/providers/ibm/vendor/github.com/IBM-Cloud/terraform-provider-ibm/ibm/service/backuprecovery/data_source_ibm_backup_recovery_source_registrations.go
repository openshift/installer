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

func DataSourceIbmBackupRecoverySourceRegistrations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoverySourceRegistrationsRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Ids specifies the list of source registration ids to return. If left empty, every source registration will be returned by default.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"include_source_credentials": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the encrypted crednetial for the registered sources will be included. Credential is first encrypted with internal key and then reencrypted with user supplied encryption key.",
			},
			"encryption_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"use_cached_data": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether we can serve the GET request from the read replica cache. There is a lag of 15 seconds between the read replica and primary data source.",
			},
			"include_external_metadata": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, the external entity metadata like maintenance mode config for the registered sources will be included.",
			},
			"ignore_tenant_migration_in_progress_check": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If true, tenant migration check will be ignored.",
			},
			"registrations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of Protection Source Registrations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Source Registration ID. This can be used to retrieve, edit or delete the source registration.",
						},
						"source_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "ID of top level source object discovered after the registration.",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies information about an object.",
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
									"protection_stats": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the count and size of protected and unprotected objects for the size.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment of the object.",
												},
												"protected_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of the protected leaf objects.",
												},
												"unprotected_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of the unprotected leaf objects.",
												},
												"deleted_protected_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the count of protected leaf objects which were deleted from the source after being protected.",
												},
												"protected_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the protected logical size in bytes.",
												},
												"unprotected_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the unprotected logical size in bytes.",
												},
											},
										},
									},
									"permissions": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the list of users, groups and users that have permissions for a given object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"object_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the id of the object.",
												},
												"users": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of users which has the permissions to the object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the user.",
															},
															"sid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the sid of the user.",
															},
															"domain": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the domain of the user.",
															},
														},
													},
												},
												"groups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of user groups which has permissions to the object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the user group.",
															},
															"sid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the sid of the user group.",
															},
															"domain": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the domain of the user group.",
															},
														},
													},
												},
												"tenant": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a tenant object.",
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
											},
										},
									},
									"mssql_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters for Msssql object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
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
											},
										},
									},
									"physical_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters for Physical object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_system_backup": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if system backup was enabled for the source in a particular run.",
												},
											},
										},
									},
								},
							},
						},
						"environment": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the environment type of the Protection Source.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user specified name for this source.",
						},
						"connection_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the id of the connection from where this source is reachable. This should only be set for a source being registered by a tenant user. This field will be depricated in future. Use connections field.",
						},
						"connections": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specfies the list of connections for the source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connection_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the id of the connection.",
									},
									"entity_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the entity id of the source. The source can a non-root entity.",
									},
									"connector_group_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the connector group id of connector groups.",
									},
									"data_source_connection_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the id of the connection in string format.",
									},
								},
							},
						},
						"connector_group_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the connector group id of connector groups.",
						},
						"data_source_connection_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the id of the connection from where this source is reachable. This should only be set for a source being registered by a tenant user. Also, this is the 'string' of connectionId. This property was added to accommodate for ID values that exceed 2^53 - 1, which is the max value for which JS maintains precision.",
						},
						"advanced_configs": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the advanced configuration for a protection source.",
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
						"authentication_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the status of the authentication during the registration of a Protection Source. 'Pending' indicates the authentication is in progress. 'Scheduled' indicates the authentication is scheduled. 'Finished' indicates the authentication is completed. 'RefreshInProgress' indicates the refresh is in progress.",
						},
						"registration_time_msecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the time when the source was registered in milliseconds.",
						},
						"last_refreshed_time_msecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the time when the source was last refreshed in milliseconds.",
						},
						"external_metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the External metadata of an entity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maintenance_mode_config": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the entity metadata for maintenance mode.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"activation_time_intervals": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the absolute intervals where the maintenance schedule is valid, i.e. maintenance_shedule is considered only for these time ranges. (For example, if there is one time range with [now_usecs, now_usecs + 10 days], the action will be done during the maintenance_schedule for the next 10 days.)The start time must be specified. The end time can be -1 which would denote an indefinite maintenance mode.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"end_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the end time of this time range.",
															},
															"start_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the start time of this time range.",
															},
														},
													},
												},
												"maintenance_schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a schedule for actions to be taken.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"periodic_time_windows": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the time range within the days of the week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_the_week": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the week day.",
																		},
																		"end_time": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the time in hours and minutes.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"hour": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the hour of this time.",
																					},
																					"minute": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the minute of this time.",
																					},
																				},
																			},
																		},
																		"start_time": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the time in hours and minutes.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"hour": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the hour of this time.",
																					},
																					"minute": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the minute of this time.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"schedule_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of schedule for this ScheduleProto.",
															},
															"time_ranges": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the time ranges in usecs.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"end_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the end time of this time range.",
																		},
																		"start_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the start time of this time range.",
																		},
																	},
																},
															},
															"timezone": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the timezone of the user of this ScheduleProto. The timezones have unique names of the form 'Area/Location'.",
															},
														},
													},
												},
												"user_message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "User provided message associated with this maintenance mode.",
												},
												"workflow_intervention_spec_list": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the type of intervention for different workflows when the source goes into maintenance mode.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"intervention": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the intervention type for ongoing tasks.",
															},
															"workflow_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the workflow type for which an intervention would be needed when maintenance mode begins.",
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
						"physical_params": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies parameters to register physical server.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the endpoint IPaddress, URL or hostname of the physical host.",
									},
									"force_register": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The agent running on a physical host will fail the registration if it is already registered as part of another cluster. By setting this option to true, agent can be forced to register with the current cluster.",
									},
									"host_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the type of host.",
									},
									"physical_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the type of physical server.",
									},
									"applications": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the list of applications to be registered with Physical Source.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
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

func dataSourceIbmBackupRecoverySourceRegistrationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_source_registrations", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getSourceRegistrationsOptions := &backuprecoveryv1.GetSourceRegistrationsOptions{}

	getSourceRegistrationsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("ids"); ok {
		var ids []int64
		for _, v := range d.Get("ids").([]interface{}) {
			idsItem := int64(v.(int))
			ids = append(ids, idsItem)
		}
		getSourceRegistrationsOptions.SetIds(ids)
	}
	if _, ok := d.GetOk("include_source_credentials"); ok {
		getSourceRegistrationsOptions.SetIncludeSourceCredentials(d.Get("include_source_credentials").(bool))
	}
	if _, ok := d.GetOk("encryption_key"); ok {
		getSourceRegistrationsOptions.SetEncryptionKey(d.Get("encryption_key").(string))
	}
	if _, ok := d.GetOk("use_cached_data"); ok {
		getSourceRegistrationsOptions.SetUseCachedData(d.Get("use_cached_data").(bool))
	}
	if _, ok := d.GetOk("include_external_metadata"); ok {
		getSourceRegistrationsOptions.SetIncludeExternalMetadata(d.Get("include_external_metadata").(bool))
	}
	if _, ok := d.GetOk("ignore_tenant_migration_in_progress_check"); ok {
		getSourceRegistrationsOptions.SetIgnoreTenantMigrationInProgressCheck(d.Get("ignore_tenant_migration_in_progress_check").(bool))
	}

	sourceRegistrations, _, err := backupRecoveryClient.GetSourceRegistrationsWithContext(context, getSourceRegistrationsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSourceRegistrationsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_source_registrations", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoverySourceRegistrationsID(d))

	if !core.IsNil(sourceRegistrations.Registrations) {
		registrations := []map[string]interface{}{}
		for _, registrationsItem := range sourceRegistrations.Registrations {
			registrationsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsSourceRegistrationReponseParamsToMap(&registrationsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_source_registrations", "read", "registrations-to-map").GetDiag()
			}
			registrations = append(registrations, registrationsItemMap)
		}
		if err = d.Set("registrations", registrations); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting registrations: %s", err), "(Data) ibm_backup_recovery_source_registrations", "read", "set-registrations").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoverySourceRegistrationsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoverySourceRegistrationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoverySourceRegistrationsSourceRegistrationReponseParamsToMap(model *backuprecoveryv1.SourceRegistrationReponseParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.SourceID != nil {
		modelMap["source_id"] = flex.IntValue(model.SourceID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ConnectionID != nil {
		modelMap["connection_id"] = flex.IntValue(model.ConnectionID)
	}
	if model.Connections != nil {
		connections := []map[string]interface{}{}
		for _, connectionsItem := range model.Connections {
			connectionsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsConnectionConfigToMap(&connectionsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			connections = append(connections, connectionsItemMap)
		}
		modelMap["connections"] = connections
	}
	if model.ConnectorGroupID != nil {
		modelMap["connector_group_id"] = flex.IntValue(model.ConnectorGroupID)
	}
	if model.DataSourceConnectionID != nil {
		modelMap["data_source_connection_id"] = *model.DataSourceConnectionID
	}
	if model.AdvancedConfigs != nil {
		advancedConfigs := []map[string]interface{}{}
		for _, advancedConfigsItem := range model.AdvancedConfigs {
			advancedConfigsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsKeyValuePairToMap(&advancedConfigsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			advancedConfigs = append(advancedConfigs, advancedConfigsItemMap)
		}
		modelMap["advanced_configs"] = advancedConfigs
	}
	if model.AuthenticationStatus != nil {
		modelMap["authentication_status"] = *model.AuthenticationStatus
	}
	if model.RegistrationTimeMsecs != nil {
		modelMap["registration_time_msecs"] = flex.IntValue(model.RegistrationTimeMsecs)
	}
	if model.LastRefreshedTimeMsecs != nil {
		modelMap["last_refreshed_time_msecs"] = flex.IntValue(model.LastRefreshedTimeMsecs)
	}
	if model.ExternalMetadata != nil {
		externalMetadataMap, err := DataSourceIbmBackupRecoverySourceRegistrationsEntityExternalMetadataToMap(model.ExternalMetadata)
		if err != nil {
			return modelMap, err
		}
		modelMap["external_metadata"] = []map[string]interface{}{externalMetadataMap}
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoverySourceRegistrationsPhysicalSourceRegistrationParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsObjectToMap(model *backuprecoveryv1.Object) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySourceRegistrationsSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	if model.ProtectionStats != nil {
		protectionStats := []map[string]interface{}{}
		for _, protectionStatsItem := range model.ProtectionStats {
			protectionStatsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectProtectionStatsSummaryToMap(&protectionStatsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			protectionStats = append(protectionStats, protectionStatsItemMap)
		}
		modelMap["protection_stats"] = protectionStats
	}
	if model.Permissions != nil {
		permissionsMap, err := DataSourceIbmBackupRecoverySourceRegistrationsPermissionInfoToMap(model.Permissions)
		if err != nil {
			return modelMap, err
		}
		modelMap["permissions"] = []map[string]interface{}{permissionsMap}
	}
	if model.MssqlParams != nil {
		mssqlParamsMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectMssqlParamsToMap(model.MssqlParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mssql_params"] = []map[string]interface{}{mssqlParamsMap}
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectPhysicalParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsSharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySourceRegistrationsSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySourceRegistrationsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsObjectProtectionStatsSummaryToMap(model *backuprecoveryv1.ObjectProtectionStatsSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.ProtectedCount != nil {
		modelMap["protected_count"] = flex.IntValue(model.ProtectedCount)
	}
	if model.UnprotectedCount != nil {
		modelMap["unprotected_count"] = flex.IntValue(model.UnprotectedCount)
	}
	if model.DeletedProtectedCount != nil {
		modelMap["deleted_protected_count"] = flex.IntValue(model.DeletedProtectedCount)
	}
	if model.ProtectedSizeBytes != nil {
		modelMap["protected_size_bytes"] = flex.IntValue(model.ProtectedSizeBytes)
	}
	if model.UnprotectedSizeBytes != nil {
		modelMap["unprotected_size_bytes"] = flex.IntValue(model.UnprotectedSizeBytes)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsPermissionInfoToMap(model *backuprecoveryv1.PermissionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	if model.Users != nil {
		users := []map[string]interface{}{}
		for _, usersItem := range model.Users {
			usersItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsUserToMap(&usersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			users = append(users, usersItemMap)
		}
		modelMap["users"] = users
	}
	if model.Groups != nil {
		groups := []map[string]interface{}{}
		for _, groupsItem := range model.Groups {
			groupsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsGroupToMap(&groupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.Tenant != nil {
		tenantMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTenantToMap(model.Tenant)
		if err != nil {
			return modelMap, err
		}
		modelMap["tenant"] = []map[string]interface{}{tenantMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsUserToMap(model *backuprecoveryv1.User) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Sid != nil {
		modelMap["sid"] = *model.Sid
	}
	if model.Domain != nil {
		modelMap["domain"] = *model.Domain
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsGroupToMap(model *backuprecoveryv1.Group) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Sid != nil {
		modelMap["sid"] = *model.Sid
	}
	if model.Domain != nil {
		modelMap["domain"] = *model.Domain
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
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
		externalVendorMetadataMap, err := DataSourceIbmBackupRecoverySourceRegistrationsExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
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
		networkMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTenantNetworkToMap(model.Network)
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

func DataSourceIbmBackupRecoverySourceRegistrationsExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
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
			customPropertiesItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
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
		metricsConfigMap, err := DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantMetricsConfigToMap(model.MetricsConfig)
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

func DataSourceIbmBackupRecoverySourceRegistrationsExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantMetricsConfigToMap(model *backuprecoveryv1.IbmTenantMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CosResourceConfig != nil {
		cosResourceConfigMap, err := DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantCOSResourceConfigToMap(model.CosResourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_resource_config"] = []map[string]interface{}{cosResourceConfigMap}
	}
	if model.IamMetricsConfig != nil {
		iamMetricsConfigMap, err := DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantIAMMetricsConfigToMap(model.IamMetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["iam_metrics_config"] = []map[string]interface{}{iamMetricsConfigMap}
	}
	if model.MeteringConfig != nil {
		meteringConfigMap, err := DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantMeteringConfigToMap(model.MeteringConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metering_config"] = []map[string]interface{}{meteringConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantCOSResourceConfigToMap(model *backuprecoveryv1.IbmTenantCOSResourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceURL != nil {
		modelMap["resource_url"] = *model.ResourceURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantIAMMetricsConfigToMap(model *backuprecoveryv1.IbmTenantIAMMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IAMURL != nil {
		modelMap["iam_url"] = *model.IAMURL
	}
	if model.BillingApiKeySecretID != nil {
		modelMap["billing_api_key_secret_id"] = *model.BillingApiKeySecretID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsIbmTenantMeteringConfigToMap(model *backuprecoveryv1.IbmTenantMeteringConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySourceRegistrationsTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySourceRegistrationsObjectMssqlParamsToMap(model *backuprecoveryv1.ObjectMssqlParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagInfo != nil {
		aagInfoMap, err := DataSourceIbmBackupRecoverySourceRegistrationsAAGInfoToMap(model.AagInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["aag_info"] = []map[string]interface{}{aagInfoMap}
	}
	if model.HostInfo != nil {
		hostInfoMap, err := DataSourceIbmBackupRecoverySourceRegistrationsHostInformationToMap(model.HostInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["host_info"] = []map[string]interface{}{hostInfoMap}
	}
	if model.IsEncrypted != nil {
		modelMap["is_encrypted"] = *model.IsEncrypted
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsAAGInfoToMap(model *backuprecoveryv1.AAGInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsHostInformationToMap(model *backuprecoveryv1.HostInformation) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySourceRegistrationsObjectPhysicalParamsToMap(model *backuprecoveryv1.ObjectPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsConnectionConfigToMap(model *backuprecoveryv1.ConnectionConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConnectionID != nil {
		modelMap["connection_id"] = flex.IntValue(model.ConnectionID)
	}
	if model.EntityID != nil {
		modelMap["entity_id"] = flex.IntValue(model.EntityID)
	}
	if model.ConnectorGroupID != nil {
		modelMap["connector_group_id"] = flex.IntValue(model.ConnectorGroupID)
	}
	if model.DataSourceConnectionID != nil {
		modelMap["data_source_connection_id"] = *model.DataSourceConnectionID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsKeyValuePairToMap(model *backuprecoveryv1.KeyValuePair) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["key"] = *model.Key
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsEntityExternalMetadataToMap(model *backuprecoveryv1.EntityExternalMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaintenanceModeConfig != nil {
		maintenanceModeConfigMap, err := DataSourceIbmBackupRecoverySourceRegistrationsMaintenanceModeConfigToMap(model.MaintenanceModeConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["maintenance_mode_config"] = []map[string]interface{}{maintenanceModeConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsMaintenanceModeConfigToMap(model *backuprecoveryv1.MaintenanceModeConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ActivationTimeIntervals != nil {
		activationTimeIntervals := []map[string]interface{}{}
		for _, activationTimeIntervalsItem := range model.ActivationTimeIntervals {
			activationTimeIntervalsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTimeRangeUsecsToMap(&activationTimeIntervalsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			activationTimeIntervals = append(activationTimeIntervals, activationTimeIntervalsItemMap)
		}
		modelMap["activation_time_intervals"] = activationTimeIntervals
	}
	if model.MaintenanceSchedule != nil {
		maintenanceScheduleMap, err := DataSourceIbmBackupRecoverySourceRegistrationsScheduleToMap(model.MaintenanceSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["maintenance_schedule"] = []map[string]interface{}{maintenanceScheduleMap}
	}
	if model.UserMessage != nil {
		modelMap["user_message"] = *model.UserMessage
	}
	if model.WorkflowInterventionSpecList != nil {
		workflowInterventionSpecList := []map[string]interface{}{}
		for _, workflowInterventionSpecListItem := range model.WorkflowInterventionSpecList {
			workflowInterventionSpecListItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsWorkflowInterventionSpecToMap(&workflowInterventionSpecListItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			workflowInterventionSpecList = append(workflowInterventionSpecList, workflowInterventionSpecListItemMap)
		}
		modelMap["workflow_intervention_spec_list"] = workflowInterventionSpecList
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsTimeRangeUsecsToMap(model *backuprecoveryv1.TimeRangeUsecs) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["end_time_usecs"] = flex.IntValue(model.EndTimeUsecs)
	modelMap["start_time_usecs"] = flex.IntValue(model.StartTimeUsecs)
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsScheduleToMap(model *backuprecoveryv1.Schedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PeriodicTimeWindows != nil {
		periodicTimeWindows := []map[string]interface{}{}
		for _, periodicTimeWindowsItem := range model.PeriodicTimeWindows {
			periodicTimeWindowsItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTimeWindowToMap(&periodicTimeWindowsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			periodicTimeWindows = append(periodicTimeWindows, periodicTimeWindowsItemMap)
		}
		modelMap["periodic_time_windows"] = periodicTimeWindows
	}
	if model.ScheduleType != nil {
		modelMap["schedule_type"] = *model.ScheduleType
	}
	if model.TimeRanges != nil {
		timeRanges := []map[string]interface{}{}
		for _, timeRangesItem := range model.TimeRanges {
			timeRangesItemMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTimeRangeUsecsToMap(&timeRangesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			timeRanges = append(timeRanges, timeRangesItemMap)
		}
		modelMap["time_ranges"] = timeRanges
	}
	if model.Timezone != nil {
		modelMap["timezone"] = *model.Timezone
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsTimeWindowToMap(model *backuprecoveryv1.TimeWindow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DayOfTheWeek != nil {
		modelMap["day_of_the_week"] = *model.DayOfTheWeek
	}
	if model.EndTime != nil {
		endTimeMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTimeToMap(model.EndTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["end_time"] = []map[string]interface{}{endTimeMap}
	}
	if model.StartTime != nil {
		startTimeMap, err := DataSourceIbmBackupRecoverySourceRegistrationsTimeToMap(model.StartTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["start_time"] = []map[string]interface{}{startTimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsTimeToMap(model *backuprecoveryv1.Time) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hour != nil {
		modelMap["hour"] = flex.IntValue(model.Hour)
	}
	if model.Minute != nil {
		modelMap["minute"] = flex.IntValue(model.Minute)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsWorkflowInterventionSpecToMap(model *backuprecoveryv1.WorkflowInterventionSpec) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["intervention"] = *model.Intervention
	modelMap["workflow_type"] = *model.WorkflowType
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySourceRegistrationsPhysicalSourceRegistrationParamsToMap(model *backuprecoveryv1.PhysicalSourceRegistrationParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["endpoint"] = *model.Endpoint
	if model.ForceRegister != nil {
		modelMap["force_register"] = *model.ForceRegister
	}
	if model.HostType != nil {
		modelMap["host_type"] = *model.HostType
	}
	if model.PhysicalType != nil {
		modelMap["physical_type"] = *model.PhysicalType
	}
	if model.Applications != nil {
		modelMap["applications"] = model.Applications
	}
	return modelMap, nil
}
