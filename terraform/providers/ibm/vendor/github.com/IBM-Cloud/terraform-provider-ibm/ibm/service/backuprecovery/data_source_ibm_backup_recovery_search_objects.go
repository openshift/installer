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

func DataSourceIbmBackupRecoverySearchObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoverySearchObjectsRead,

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
			"search_string": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the search string to filter the objects. This search string will be applicable for objectnames. User can specify a wildcard character '*' as a suffix to a string where all object names are matched with the prefix string. For example, if vm1 and vm2 are the names of objects, user can specify vm* to list the objects. If not specified, then all the objects will be returned which will match other filtering criteria.",
			},
			"environments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the environment type to filter objects.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"protection_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the protection type to filter objects.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"protection_group_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of Protection Group ids to filter the objects. If specified, the objects protected by specified Protection Group ids will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"object_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of Object ids to filter.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"os_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the operating system types to filter objects on.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of Protection Source object ids to filter the objects. If specified, the object which are present in those Sources will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},
			"source_uuids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of Protection Source object uuids to filter the objects. If specified, the object which are present in those Sources will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_protected": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies the protection status of objects. If set to true, only protected objects will be returned. If set to false, only unprotected objects will be returned. If not specified, all objects will be returned.",
			},
			"is_deleted": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, then objects which are deleted on atleast one cluster will be returned. If not set or set to false then objects which are registered on atleast one cluster are returned.",
			},
			"last_run_status_list": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies a list of status of the object's last protection run. Only objects with last run status of these will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_identifiers": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the list of cluster identifiers. Format is clusterId:clusterIncarnationId. Only records from clusters having these identifiers will be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_deleted_objects": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to include deleted objects in response. These objects can't be protected but can be recovered. This field is deprecated.",
			},
			"pagination_cookie": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the pagination cookie with which subsequent parts of the response can be fetched.",
			},
			"object_count": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of objects to be fetched for the specified pagination cookie.",
			},
			"must_have_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies tags which must be all present in the document.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"might_have_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies list of tags, one or more of which might be present in the document. These are OR'ed together and the resulting criteria AND'ed with the rest of the query.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"must_have_snapshot_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies snapshot tags which must be all present in the document.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"might_have_snapshot_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies list of snapshot tags, one or more of which might be present in the document. These are OR'ed together and the resulting criteria AND'ed with the rest of the query.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tag_search_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the tag name to filter the tagged objects and snapshots. User can specify a wildcard character '*' as a suffix to a string where all object's tag names are matched with the prefix string.",
			},
			"tag_names": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the tag names to filter the tagged objects and snapshots.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tag_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the tag names to filter the tagged objects and snapshots.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tag_categories": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the tag category to filter the objects and snapshots.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tag_sub_categories": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the tag subcategory to filter the objects and snapshots.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_helios_tag_info_for_objects": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "pecifies whether to include helios tags information for objects in response. Default value is false.",
			},
			"external_filters": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the key-value pairs to filtering the results for the search. Each filter is of the form 'key:value'. The filter 'externalFilters:k1:v1&externalFilters:k2:v2&externalFilters:k2:v3' returns the documents where each document will match the query (k1=v1) AND (k2=v2 OR k2 = v3). Allowed keys: - vmBiosUuid - graphUuid - arn - instanceId - bucketName - azureId.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of Objects.",
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
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies tag applied to the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies Id of tag applied to the object.",
									},
								},
							},
						},
						"snapshot_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies snapshot tags applied to the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies Id of tag applied to the object.",
									},
									"run_ids": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"helios_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the helios tag information for the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies category of tag applied to the object.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies name of tag applied to the object.",
									},
									"sub_category": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies subCategory of tag applied to the object.",
									},
									"third_party_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies thirdPartyName of tag applied to the object.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the type (ex custom, thirdparty, system) of tag applied to the object.",
									},
									"ui_color": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the color of tag applied to the object.",
									},
									"updated_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies update time of tag applied to the object.",
									},
									"uuid": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies Uuid of tag applied to the object.",
									},
								},
							},
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"object_protection_infos": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the object info on each cluster.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"object_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the object id.",
									},
									"source_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the source id.",
									},
									"view_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the view id for the object.",
									},
									"region_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the region id where this object belongs to.",
									},
									"cluster_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the cluster id where this object belongs to.",
									},
									"cluster_incarnation_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the cluster incarnation id where this object belongs to.",
									},
									"tenant_ids": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of Tenants the object belongs to.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"is_deleted": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether the object is deleted. Deleted objects can't be protected but can be recovered or unprotected.",
									},
									"protection_groups": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of protection groups protecting this object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the protection group name.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the protection group id.",
												},
												"protection_env_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the protection type of the job if any.",
												},
												"policy_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the policy name for this group.",
												},
												"policy_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the policy id for this group.",
												},
												"last_backup_run_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of last local back up run.",
												},
												"last_archival_run_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of last archival run.",
												},
												"last_replication_run_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of last replication run.",
												},
												"last_run_sla_violated": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if the sla is violated in last run.",
												},
											},
										},
									},
									"object_backup_configuration": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of object protections.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"policy_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the policy name for this group.",
												},
												"policy_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the policy id for this protection.",
												},
												"last_backup_run_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of last local back up run.",
												},
												"last_archival_run_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of last archival run.",
												},
												"last_replication_run_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of last replication run.",
												},
												"last_run_sla_violated": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if the sla is violated in last run.",
												},
											},
										},
									},
									"last_run_status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the status of the object's last protection run.",
									},
								},
							},
						},
						"secondary_ids": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies secondary IDs associated to the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies name of the secondary ID for an object.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies value of the secondary ID for an object.",
									},
								},
							},
						},
						"tagged_snapshots": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the helios tagged snapshots (snapshots which are tagged by user or thirdparty in control plane) for the object.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cluster_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the cluster Id of the tagged snapshot.",
									},
									"cluster_incarnation_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the clusterIncarnationId of the tagged snapshot.",
									},
									"job_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the jobId of the tagged snapshot.",
									},
									"object_uuid": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the object uuid of the tagged snapshot.",
									},
									"run_start_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the runStartTimeUsecs of the tagged snapshot.",
									},
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies tag applied to the object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"category": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies category of tag applied to the object.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of tag applied to the object.",
												},
												"sub_category": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies subCategory of tag applied to the object.",
												},
												"third_party_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies thirdPartyName of tag applied to the object.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type (ex custom, thirdparty, system) of tag applied to the object.",
												},
												"ui_color": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the color of tag applied to the object.",
												},
												"updated_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies update time of tag applied to the object.",
												},
												"uuid": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies Uuid of tag applied to the object.",
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

func dataSourceIbmBackupRecoverySearchObjectsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_search_objects", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	searchObjectsOptions := &backuprecoveryv1.SearchObjectsOptions{}

	searchObjectsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("request_initiator_type"); ok {
		searchObjectsOptions.SetRequestInitiatorType(d.Get("request_initiator_type").(string))
	}
	if _, ok := d.GetOk("search_string"); ok {
		searchObjectsOptions.SetSearchString(d.Get("search_string").(string))
	}
	if _, ok := d.GetOk("environments"); ok {
		var environments []string
		for _, v := range d.Get("environments").([]interface{}) {
			environmentsItem := v.(string)
			environments = append(environments, environmentsItem)
		}
		searchObjectsOptions.SetEnvironments(environments)
	}
	if _, ok := d.GetOk("protection_types"); ok {
		var protectionTypes []string
		for _, v := range d.Get("protection_types").([]interface{}) {
			protectionTypesItem := v.(string)
			protectionTypes = append(protectionTypes, protectionTypesItem)
		}
		searchObjectsOptions.SetProtectionTypes(protectionTypes)
	}
	if _, ok := d.GetOk("protection_group_ids"); ok {
		var protectionGroupIds []string
		for _, v := range d.Get("protection_group_ids").([]interface{}) {
			protectionGroupIdsItem := v.(string)
			protectionGroupIds = append(protectionGroupIds, protectionGroupIdsItem)
		}
		searchObjectsOptions.SetProtectionGroupIds(protectionGroupIds)
	}
	if _, ok := d.GetOk("object_ids"); ok {
		var objectIds []int64
		for _, v := range d.Get("object_ids").([]interface{}) {
			objectIdsItem := int64(v.(int))
			objectIds = append(objectIds, objectIdsItem)
		}
		searchObjectsOptions.SetObjectIds(objectIds)
	}
	if _, ok := d.GetOk("os_types"); ok {
		var osTypes []string
		for _, v := range d.Get("os_types").([]interface{}) {
			osTypesItem := v.(string)
			osTypes = append(osTypes, osTypesItem)
		}
		searchObjectsOptions.SetOsTypes(osTypes)
	}
	if _, ok := d.GetOk("source_ids"); ok {
		var sourceIds []int64
		for _, v := range d.Get("source_ids").([]interface{}) {
			sourceIdsItem := int64(v.(int))
			sourceIds = append(sourceIds, sourceIdsItem)
		}
		searchObjectsOptions.SetSourceIds(sourceIds)
	}
	if _, ok := d.GetOk("source_uuids"); ok {
		var sourceUUIDs []string
		for _, v := range d.Get("source_uuids").([]interface{}) {
			sourceUUIDsItem := v.(string)
			sourceUUIDs = append(sourceUUIDs, sourceUUIDsItem)
		}
		searchObjectsOptions.SetSourceUUIDs(sourceUUIDs)
	}
	if _, ok := d.GetOk("is_protected"); ok {
		searchObjectsOptions.SetIsProtected(d.Get("is_protected").(bool))
	}
	if _, ok := d.GetOk("is_deleted"); ok {
		searchObjectsOptions.SetIsDeleted(d.Get("is_deleted").(bool))
	}
	if _, ok := d.GetOk("last_run_status_list"); ok {
		var lastRunStatusList []string
		for _, v := range d.Get("last_run_status_list").([]interface{}) {
			lastRunStatusListItem := v.(string)
			lastRunStatusList = append(lastRunStatusList, lastRunStatusListItem)
		}
		searchObjectsOptions.SetLastRunStatusList(lastRunStatusList)
	}
	if _, ok := d.GetOk("cluster_identifiers"); ok {
		var clusterIdentifiers []string
		for _, v := range d.Get("cluster_identifiers").([]interface{}) {
			clusterIdentifiersItem := v.(string)
			clusterIdentifiers = append(clusterIdentifiers, clusterIdentifiersItem)
		}
		searchObjectsOptions.SetClusterIdentifiers(clusterIdentifiers)
	}
	if _, ok := d.GetOk("include_deleted_objects"); ok {
		searchObjectsOptions.SetIncludeDeletedObjects(d.Get("include_deleted_objects").(bool))
	}
	if _, ok := d.GetOk("pagination_cookie"); ok {
		searchObjectsOptions.SetPaginationCookie(d.Get("pagination_cookie").(string))
	}
	if _, ok := d.GetOk("object_count"); ok {
		searchObjectsOptions.SetCount(int64(d.Get("object_count").(int)))
	}
	if _, ok := d.GetOk("must_have_tag_ids"); ok {
		var mustHaveTagIds []string
		for _, v := range d.Get("must_have_tag_ids").([]interface{}) {
			mustHaveTagIdsItem := v.(string)
			mustHaveTagIds = append(mustHaveTagIds, mustHaveTagIdsItem)
		}
		searchObjectsOptions.SetMustHaveTagIds(mustHaveTagIds)
	}
	if _, ok := d.GetOk("might_have_tag_ids"); ok {
		var mightHaveTagIds []string
		for _, v := range d.Get("might_have_tag_ids").([]interface{}) {
			mightHaveTagIdsItem := v.(string)
			mightHaveTagIds = append(mightHaveTagIds, mightHaveTagIdsItem)
		}
		searchObjectsOptions.SetMightHaveTagIds(mightHaveTagIds)
	}
	if _, ok := d.GetOk("must_have_snapshot_tag_ids"); ok {
		var mustHaveSnapshotTagIds []string
		for _, v := range d.Get("must_have_snapshot_tag_ids").([]interface{}) {
			mustHaveSnapshotTagIdsItem := v.(string)
			mustHaveSnapshotTagIds = append(mustHaveSnapshotTagIds, mustHaveSnapshotTagIdsItem)
		}
		searchObjectsOptions.SetMustHaveSnapshotTagIds(mustHaveSnapshotTagIds)
	}
	if _, ok := d.GetOk("might_have_snapshot_tag_ids"); ok {
		var mightHaveSnapshotTagIds []string
		for _, v := range d.Get("might_have_snapshot_tag_ids").([]interface{}) {
			mightHaveSnapshotTagIdsItem := v.(string)
			mightHaveSnapshotTagIds = append(mightHaveSnapshotTagIds, mightHaveSnapshotTagIdsItem)
		}
		searchObjectsOptions.SetMightHaveSnapshotTagIds(mightHaveSnapshotTagIds)
	}
	if _, ok := d.GetOk("tag_search_name"); ok {
		searchObjectsOptions.SetTagSearchName(d.Get("tag_search_name").(string))
	}
	if _, ok := d.GetOk("tag_names"); ok {
		var tagNames []string
		for _, v := range d.Get("tag_names").([]interface{}) {
			tagNamesItem := v.(string)
			tagNames = append(tagNames, tagNamesItem)
		}
		searchObjectsOptions.SetTagNames(tagNames)
	}
	if _, ok := d.GetOk("tag_types"); ok {
		var tagTypes []string
		for _, v := range d.Get("tag_types").([]interface{}) {
			tagTypesItem := v.(string)
			tagTypes = append(tagTypes, tagTypesItem)
		}
		searchObjectsOptions.SetTagTypes(tagTypes)
	}
	if _, ok := d.GetOk("tag_categories"); ok {
		var tagCategories []string
		for _, v := range d.Get("tag_categories").([]interface{}) {
			tagCategoriesItem := v.(string)
			tagCategories = append(tagCategories, tagCategoriesItem)
		}
		searchObjectsOptions.SetTagCategories(tagCategories)
	}
	if _, ok := d.GetOk("tag_sub_categories"); ok {
		var tagSubCategories []string
		for _, v := range d.Get("tag_sub_categories").([]interface{}) {
			tagSubCategoriesItem := v.(string)
			tagSubCategories = append(tagSubCategories, tagSubCategoriesItem)
		}
		searchObjectsOptions.SetTagSubCategories(tagSubCategories)
	}
	if _, ok := d.GetOk("include_helios_tag_info_for_objects"); ok {
		searchObjectsOptions.SetIncludeHeliosTagInfoForObjects(d.Get("include_helios_tag_info_for_objects").(bool))
	}
	if _, ok := d.GetOk("external_filters"); ok {
		var externalFilters []string
		for _, v := range d.Get("external_filters").([]interface{}) {
			externalFiltersItem := v.(string)
			externalFilters = append(externalFilters, externalFiltersItem)
		}
		searchObjectsOptions.SetExternalFilters(externalFilters)
	}

	objectsSearchResponseBody, _, err := backupRecoveryClient.SearchObjectsWithContext(context, searchObjectsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SearchObjectsWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_search_objects", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoverySearchObjectsID(d))

	if !core.IsNil(objectsSearchResponseBody.Objects) {
		objects := []map[string]interface{}{}
		for _, objectsItem := range objectsSearchResponseBody.Objects {
			objectsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsSearchObjectToMap(&objectsItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_search_objects", "read", "objects-to-map").GetDiag()
			}
			objects = append(objects, objectsItemMap)
		}
		if err = d.Set("objects", objects); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting objects: %s", err), "(Data) ibm_backup_recovery_search_objects", "read", "set-objects").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoverySearchObjectsID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoverySearchObjectsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoverySearchObjectsSearchObjectToMap(model *backuprecoveryv1.SearchObject) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	if model.ProtectionStats != nil {
		protectionStats := []map[string]interface{}{}
		for _, protectionStatsItem := range model.ProtectionStats {
			protectionStatsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectProtectionStatsSummaryToMap(&protectionStatsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			protectionStats = append(protectionStats, protectionStatsItemMap)
		}
		modelMap["protection_stats"] = protectionStats
	}
	if model.Permissions != nil {
		permissionsMap, err := DataSourceIbmBackupRecoverySearchObjectsPermissionInfoToMap(model.Permissions)
		if err != nil {
			return modelMap, err
		}
		modelMap["permissions"] = []map[string]interface{}{permissionsMap}
	}
	if model.MssqlParams != nil {
		mssqlParamsMap, err := DataSourceIbmBackupRecoverySearchObjectsSearchObjectMssqlParamsToMap(model.MssqlParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mssql_params"] = []map[string]interface{}{mssqlParamsMap}
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoverySearchObjectsSearchObjectPhysicalParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsTagInfoToMap(&tagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tags = append(tags, tagsItemMap)
		}
		modelMap["tags"] = tags
	}
	if model.SnapshotTags != nil {
		snapshotTags := []map[string]interface{}{}
		for _, snapshotTagsItem := range model.SnapshotTags {
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.HeliosTags != nil {
		heliosTags := []map[string]interface{}{}
		for _, heliosTagsItem := range model.HeliosTags {
			heliosTagsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsHeliosTagInfoToMap(&heliosTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			heliosTags = append(heliosTags, heliosTagsItemMap)
		}
		modelMap["helios_tags"] = heliosTags
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchObjectsSearchObjectSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ObjectProtectionInfos != nil {
		objectProtectionInfos := []map[string]interface{}{}
		for _, objectProtectionInfosItem := range model.ObjectProtectionInfos {
			objectProtectionInfosItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectProtectionInfoToMap(&objectProtectionInfosItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			objectProtectionInfos = append(objectProtectionInfos, objectProtectionInfosItemMap)
		}
		modelMap["object_protection_infos"] = objectProtectionInfos
	}
	if model.SecondaryIds != nil {
		secondaryIds := []map[string]interface{}{}
		for _, secondaryIdsItem := range model.SecondaryIds {
			secondaryIdsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsSecondaryIDToMap(&secondaryIdsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			secondaryIds = append(secondaryIds, secondaryIdsItemMap)
		}
		modelMap["secondary_ids"] = secondaryIds
	}
	if model.TaggedSnapshots != nil {
		taggedSnapshots := []map[string]interface{}{}
		for _, taggedSnapshotsItem := range model.TaggedSnapshots {
			taggedSnapshotsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsTaggedSnapshotInfoToMap(&taggedSnapshotsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			taggedSnapshots = append(taggedSnapshots, taggedSnapshotsItemMap)
		}
		modelMap["tagged_snapshots"] = taggedSnapshots
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsSharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsObjectProtectionStatsSummaryToMap(model *backuprecoveryv1.ObjectProtectionStatsSummary) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchObjectsPermissionInfoToMap(model *backuprecoveryv1.PermissionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	if model.Users != nil {
		users := []map[string]interface{}{}
		for _, usersItem := range model.Users {
			usersItemMap, err := DataSourceIbmBackupRecoverySearchObjectsUserToMap(&usersItem) // #nosec G601
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
			groupsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsGroupToMap(&groupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.Tenant != nil {
		tenantMap, err := DataSourceIbmBackupRecoverySearchObjectsTenantToMap(model.Tenant)
		if err != nil {
			return modelMap, err
		}
		modelMap["tenant"] = []map[string]interface{}{tenantMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsUserToMap(model *backuprecoveryv1.User) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchObjectsGroupToMap(model *backuprecoveryv1.Group) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchObjectsTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
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
		externalVendorMetadataMap, err := DataSourceIbmBackupRecoverySearchObjectsExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
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
		networkMap, err := DataSourceIbmBackupRecoverySearchObjectsTenantNetworkToMap(model.Network)
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

func DataSourceIbmBackupRecoverySearchObjectsExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := DataSourceIbmBackupRecoverySearchObjectsIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
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
			customPropertiesItemMap, err := DataSourceIbmBackupRecoverySearchObjectsExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
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
		metricsConfigMap, err := DataSourceIbmBackupRecoverySearchObjectsIbmTenantMetricsConfigToMap(model.MetricsConfig)
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

func DataSourceIbmBackupRecoverySearchObjectsExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsIbmTenantMetricsConfigToMap(model *backuprecoveryv1.IbmTenantMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CosResourceConfig != nil {
		cosResourceConfigMap, err := DataSourceIbmBackupRecoverySearchObjectsIbmTenantCOSResourceConfigToMap(model.CosResourceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cos_resource_config"] = []map[string]interface{}{cosResourceConfigMap}
	}
	if model.IamMetricsConfig != nil {
		iamMetricsConfigMap, err := DataSourceIbmBackupRecoverySearchObjectsIbmTenantIAMMetricsConfigToMap(model.IamMetricsConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["iam_metrics_config"] = []map[string]interface{}{iamMetricsConfigMap}
	}
	if model.MeteringConfig != nil {
		meteringConfigMap, err := DataSourceIbmBackupRecoverySearchObjectsIbmTenantMeteringConfigToMap(model.MeteringConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["metering_config"] = []map[string]interface{}{meteringConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsIbmTenantCOSResourceConfigToMap(model *backuprecoveryv1.IbmTenantCOSResourceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceURL != nil {
		modelMap["resource_url"] = *model.ResourceURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsIbmTenantIAMMetricsConfigToMap(model *backuprecoveryv1.IbmTenantIAMMetricsConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IAMURL != nil {
		modelMap["iam_url"] = *model.IAMURL
	}
	if model.BillingApiKeySecretID != nil {
		modelMap["billing_api_key_secret_id"] = *model.BillingApiKeySecretID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsIbmTenantMeteringConfigToMap(model *backuprecoveryv1.IbmTenantMeteringConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchObjectsTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchObjectsSearchObjectMssqlParamsToMap(model *backuprecoveryv1.SearchObjectMssqlParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagInfo != nil {
		aagInfoMap, err := DataSourceIbmBackupRecoverySearchObjectsAAGInfoToMap(model.AagInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["aag_info"] = []map[string]interface{}{aagInfoMap}
	}
	if model.HostInfo != nil {
		hostInfoMap, err := DataSourceIbmBackupRecoverySearchObjectsHostInformationToMap(model.HostInfo)
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

func DataSourceIbmBackupRecoverySearchObjectsAAGInfoToMap(model *backuprecoveryv1.AAGInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsHostInformationToMap(model *backuprecoveryv1.HostInformation) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchObjectsSearchObjectPhysicalParamsToMap(model *backuprecoveryv1.SearchObjectPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsTagInfoToMap(model *backuprecoveryv1.TagInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["tag_id"] = *model.TagID
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsSnapshotTagInfoToMap(model *backuprecoveryv1.SnapshotTagInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["tag_id"] = *model.TagID
	if model.RunIds != nil {
		modelMap["run_ids"] = model.RunIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsHeliosTagInfoToMap(model *backuprecoveryv1.HeliosTagInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Category != nil {
		modelMap["category"] = *model.Category
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.SubCategory != nil {
		modelMap["sub_category"] = *model.SubCategory
	}
	if model.ThirdPartyName != nil {
		modelMap["third_party_name"] = *model.ThirdPartyName
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.UiColor != nil {
		modelMap["ui_color"] = *model.UiColor
	}
	if model.UpdatedTimeUsecs != nil {
		modelMap["updated_time_usecs"] = flex.IntValue(model.UpdatedTimeUsecs)
	}
	modelMap["uuid"] = *model.UUID
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsSearchObjectSourceInfoToMap(model *backuprecoveryv1.SearchObjectSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	if model.ProtectionStats != nil {
		protectionStats := []map[string]interface{}{}
		for _, protectionStatsItem := range model.ProtectionStats {
			protectionStatsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectProtectionStatsSummaryToMap(&protectionStatsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			protectionStats = append(protectionStats, protectionStatsItemMap)
		}
		modelMap["protection_stats"] = protectionStats
	}
	if model.Permissions != nil {
		permissionsMap, err := DataSourceIbmBackupRecoverySearchObjectsPermissionInfoToMap(model.Permissions)
		if err != nil {
			return modelMap, err
		}
		modelMap["permissions"] = []map[string]interface{}{permissionsMap}
	}
	if model.MssqlParams != nil {
		mssqlParamsMap, err := DataSourceIbmBackupRecoverySearchObjectsSearchObjectSourceInfoMssqlParamsToMap(model.MssqlParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mssql_params"] = []map[string]interface{}{mssqlParamsMap}
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoverySearchObjectsSearchObjectSourceInfoPhysicalParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsSearchObjectSourceInfoMssqlParamsToMap(model *backuprecoveryv1.SearchObjectSourceInfoMssqlParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagInfo != nil {
		aagInfoMap, err := DataSourceIbmBackupRecoverySearchObjectsAAGInfoToMap(model.AagInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["aag_info"] = []map[string]interface{}{aagInfoMap}
	}
	if model.HostInfo != nil {
		hostInfoMap, err := DataSourceIbmBackupRecoverySearchObjectsHostInformationToMap(model.HostInfo)
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

func DataSourceIbmBackupRecoverySearchObjectsSearchObjectSourceInfoPhysicalParamsToMap(model *backuprecoveryv1.SearchObjectSourceInfoPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsObjectProtectionInfoToMap(model *backuprecoveryv1.ObjectProtectionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	if model.SourceID != nil {
		modelMap["source_id"] = flex.IntValue(model.SourceID)
	}
	if model.ViewID != nil {
		modelMap["view_id"] = flex.IntValue(model.ViewID)
	}
	if model.RegionID != nil {
		modelMap["region_id"] = *model.RegionID
	}
	if model.ClusterID != nil {
		modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	}
	if model.ClusterIncarnationID != nil {
		modelMap["cluster_incarnation_id"] = flex.IntValue(model.ClusterIncarnationID)
	}
	if model.TenantIds != nil {
		modelMap["tenant_ids"] = model.TenantIds
	}
	if model.IsDeleted != nil {
		modelMap["is_deleted"] = *model.IsDeleted
	}
	if model.ProtectionGroups != nil {
		protectionGroups := []map[string]interface{}{}
		for _, protectionGroupsItem := range model.ProtectionGroups {
			protectionGroupsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsObjectProtectionGroupSummaryToMap(&protectionGroupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			protectionGroups = append(protectionGroups, protectionGroupsItemMap)
		}
		modelMap["protection_groups"] = protectionGroups
	}
	if model.ObjectBackupConfiguration != nil {
		objectBackupConfiguration := []map[string]interface{}{}
		for _, objectBackupConfigurationItem := range model.ObjectBackupConfiguration {
			objectBackupConfigurationItemMap, err := DataSourceIbmBackupRecoverySearchObjectsProtectionSummaryToMap(&objectBackupConfigurationItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			objectBackupConfiguration = append(objectBackupConfiguration, objectBackupConfigurationItemMap)
		}
		modelMap["object_backup_configuration"] = objectBackupConfiguration
	}
	if model.LastRunStatus != nil {
		modelMap["last_run_status"] = *model.LastRunStatus
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsObjectProtectionGroupSummaryToMap(model *backuprecoveryv1.ObjectProtectionGroupSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ProtectionEnvType != nil {
		modelMap["protection_env_type"] = *model.ProtectionEnvType
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.LastBackupRunStatus != nil {
		modelMap["last_backup_run_status"] = *model.LastBackupRunStatus
	}
	if model.LastArchivalRunStatus != nil {
		modelMap["last_archival_run_status"] = *model.LastArchivalRunStatus
	}
	if model.LastReplicationRunStatus != nil {
		modelMap["last_replication_run_status"] = *model.LastReplicationRunStatus
	}
	if model.LastRunSlaViolated != nil {
		modelMap["last_run_sla_violated"] = *model.LastRunSlaViolated
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsProtectionSummaryToMap(model *backuprecoveryv1.ProtectionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.LastBackupRunStatus != nil {
		modelMap["last_backup_run_status"] = *model.LastBackupRunStatus
	}
	if model.LastArchivalRunStatus != nil {
		modelMap["last_archival_run_status"] = *model.LastArchivalRunStatus
	}
	if model.LastReplicationRunStatus != nil {
		modelMap["last_replication_run_status"] = *model.LastReplicationRunStatus
	}
	if model.LastRunSlaViolated != nil {
		modelMap["last_run_sla_violated"] = *model.LastRunSlaViolated
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsSecondaryIDToMap(model *backuprecoveryv1.SecondaryID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchObjectsTaggedSnapshotInfoToMap(model *backuprecoveryv1.TaggedSnapshotInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterID != nil {
		modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	}
	if model.ClusterIncarnationID != nil {
		modelMap["cluster_incarnation_id"] = flex.IntValue(model.ClusterIncarnationID)
	}
	if model.JobID != nil {
		modelMap["job_id"] = flex.IntValue(model.JobID)
	}
	if model.ObjectUUID != nil {
		modelMap["object_uuid"] = *model.ObjectUUID
	}
	if model.RunStartTimeUsecs != nil {
		modelMap["run_start_time_usecs"] = flex.IntValue(model.RunStartTimeUsecs)
	}
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchObjectsHeliosTagInfoToMap(&tagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tags = append(tags, tagsItemMap)
		}
		modelMap["tags"] = tags
	}
	return modelMap, nil
}
