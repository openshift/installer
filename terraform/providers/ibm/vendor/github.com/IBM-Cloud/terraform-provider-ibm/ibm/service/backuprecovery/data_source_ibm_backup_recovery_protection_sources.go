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

func DataSourceIbmBackupRecoveryProtectionSources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryProtectionSourcesRead,

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"exclude_office365_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the Object types to be filtered out for Office 365 that match the passed in types such as 'kDomain', 'kOutlook', 'kMailbox', etc. For example, set this parameter to 'kMailbox' to exclude Mailbox Objects from being returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"get_teams_channels": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Filter policies by a list of policy ids.",
			},
			"after_cursor_entity_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the entity id starting from which the items are to be returned.",
			},
			"before_cursor_entity_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the entity id upto which the items are to be returned.",
			},
			"node_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the entity id for the Node at any level within the Source entity hierarchy whose children are to be paginated.",
			},
			"page_size": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the maximum number of entities to be returned within the page.",
			},
			"has_valid_mailbox": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, users with valid mailbox will be returned.",
			},
			"has_valid_onedrive": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, users with valid onedrive will be returned.",
			},
			"is_security_group": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, Groups which are security enabled will be returned.",
			},
			"backup_recovery_protection_source_nodes_id": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Return the Object subtree for the passed in Protection Source id.",
			},
			"num_levels": &schema.Schema{
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "Specifies the expected number of levels from the root node to be returned in the entity hierarchy response.",
			},
			"exclude_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter out the Object types (and their subtrees) that match the passed in types such as 'kVCenter', 'kFolder', 'kDatacenter', 'kComputeResource', 'kResourcePool', 'kDatastore', 'kHostSystem', 'kVirtualMachine', etc. For example, set this parameter to 'kResourcePool' to exclude Resource Pool Objects from being returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_aws_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the Object types to be filtered out for AWS that match the passed in types such as 'kEC2Instance', 'kRDSInstance', 'kAuroraCluster', 'kTag', 'kAuroraTag', 'kRDSTag', kS3Bucket, kS3Tag. For example, set this parameter to 'kEC2Instance' to exclude ec2 instance from being returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_kubernetes_types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the Object types to be filtered out for Kubernetes that match the passed in types such as 'kService'. For example, set this parameter to 'kService' to exclude services from being returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_datastores": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this parameter to true to also return kDatastore object types found in the Source in addition to their Object subtrees. By default, datastores are not returned.",
			},
			"include_networks": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this parameter to true to also return kNetwork object types found in the Source in addition to their Object subtrees. By default, network objects are not returned.",
			},
			"include_vm_folders": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this parameter to true to also return kVMFolder object types found in the Source in addition to their Object subtrees. By default, VM folder objects are not returned.",
			},
			"include_sfdc_fields": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this parameter to true to also return fields of the object found in the Source in addition to their Object subtrees. By default, Sfdc object fields are not returned.",
			},
			"include_system_v_apps": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set this parameter to true to also return system VApp object types found in the Source in addition to their Object subtrees. By default, VM folder objects are not returned.",
			},
			"environments": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Return only Protection Sources that match the passed in environment type such as 'kVMware', 'kSQL', 'kView' 'kPhysical', 'kPuppeteer', 'kPure', 'kNetapp', 'kGenericNas', 'kHyperV', 'kAcropolis', or 'kAzure'. For example, set this parameter to 'kVMware' to only return the Sources (and their Object subtrees) found in the 'kVMware' (VMware vCenter Server) environment.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"environment": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This field is deprecated. Use environments instead.",
			},
			"include_entity_permission_info": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If specified, then a list of entites with permissions assigned to them are returned.",
			},
			"sids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter the object subtree for the sids given in the list.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"include_source_credentials": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If specified, then crednetial for the registered sources will be included. Credential is first encrypted with internal key and then reencrypted with user supplied 'encryption_key'.",
			},
			"encryption_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key to be used to encrypt the source credential. If include_source_credentials is set to true this key must be specified.",
			},
			"include_object_protection_info": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If specified, the object protection of entities(if any) will be returned.",
			},
			"prune_non_critical_info": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to prune non critical info within entities. Incase of VMs, virtual disk information will be pruned. Incase of Office365, metadata about user entities will be pruned. This can be used to limit the size of the response by caller.",
			},
			"prune_aggregation_info": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether to prune the aggregation information about the number of entities protected/unprotected.",
			},
			"request_initiator_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the type of the request. Possible values are UIUser and UIAuto, which means the request is triggered by user or is an auto refresh request. Services like magneto will use this to determine the priority of the requests, so that it can more intelligently handle overload situations by prioritizing higher priority requests.",
			},
			"use_cached_data": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether we can serve the GET request to the read replica cache. setting this to true ensures that the API request is served to the read replica. setting this to false will serve the request to the master.",
			},
			"all_under_hierarchy": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "AllUnderHierarchy specifies if objects of all the tenants under the hierarchy of the logged in user's organization should be returned.",
			},
			"protection_sources": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies list of protection sources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_nodes": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the child subtree used to store additional application-level Objects. Different environments use the subtree to store application-level information. For example for SQL Server, this subtree stores the SQL Server instances running on a VM.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"nodes": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies children of the current node in the Protection Sources hierarchy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"application_nodes": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the child subtree used to store additional application-level Objects. Different environments use the subtree to store application-level information. For example for SQL Server, this subtree stores the SQL Server instances running on a VM.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{},
													},
												},
												"entity_pagination_parameters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the cursor based pagination parameters for Protection Source and its children. Pagination is supported at a given level within the Protection Source Hierarchy with the help of before or after cursors. A Cursor will always refer to a specific source within the source dataset but will be invalidated if the item is removed.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"after_cursor_entity_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the entity id starting from which the items are to be returned.",
															},
															"before_cursor_entity_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the entity id upto which the items are to be returned.",
															},
															"node_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the entity id for the Node at any level within the Source entity hierarchy whose children are to be paginated.",
															},
															"page_size": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the maximum number of entities to be returned within the page.",
															},
														},
													},
												},
												"entity_permission_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the permission information of entities.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"entity_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the entity id.",
															},
															"groups": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies struct with basic group details.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"domain": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies domain name of the user.",
																		},
																		"group_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies group name of the group.",
																		},
																		"sid": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies unique Security ID (SID) of the user.",
																		},
																		"tenant_ids": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the tenants to which the group belongs to.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"is_inferred": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the Entity Permission Information is inferred or not. For example, SQL application hosted over vCenter will have inferred entity permission information.",
															},
															"is_registered_by_sp": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether this entity is registered by the SP or not. This will be populated only if the entity is a root entity. Refer to magneto/base/permissions.proto for details.",
															},
															"registering_tenant_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the tenant id that registered this entity. This will be populated only if the entity is a root entity.",
															},
															"tenant": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies struct with basic tenant details.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"bifrost_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if this tenant is bifrost enabled or not.",
																		},
																		"is_managed_on_helios": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether this tenant is manged on helios.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies name of the tenant.",
																		},
																		"tenant_id": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies the unique id of the tenant.",
																		},
																	},
																},
															},
															"users": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies struct with basic user details.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"domain": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies domain name of the user.",
																		},
																		"sid": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies unique Security ID (SID) of the user.",
																		},
																		"tenant_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the tenant to which the user belongs to.",
																		},
																		"user_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies user name of the user.",
																		},
																	},
																},
															},
														},
													},
												},
												"logical_size": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the logical size of the data in bytes for the Object on this node. Presence of this field indicates this node is a leaf node.",
												},
												"nodes": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies children of the current node in the Protection Sources hierarchy. When representing Objects in memory, the entire Object subtree hierarchy is represented. You can use this subtree to navigate down the Object hierarchy.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{},
													},
												},
												"object_protection_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the Object Protection Info of the Protection Source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"auto_protect_parent_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the auto protect parent id if this entity is protected based on auto protection. This is only specified for leaf entities.",
															},
															"entity_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the entity id.",
															},
															"has_active_object_protection_spec": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies if the entity is under object protection.",
															},
														},
													},
												},
												"protected_sources_summary": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Array of Protected Objects. Specifies aggregated information about all the child Objects of this node that are currently protected by a Protection Job. There is one entry for each environment that is being backed up. The aggregated information for the Object hierarchy's environment will be available at the 0th index of the vector.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
															},
															"leaves_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the number of leaf nodes under the subtree of this node.",
															},
															"total_logical_size": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total logical size of the data under the subtree of this node.",
															},
														},
													},
												},
												"protection_source": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies details about an Acropolis Protection Source when the environment is set to 'kAcropolis'.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"connection_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the connection id of the tenant.",
															},
															"connector_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the connector group id of the connector groups.",
															},
															"custom_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the user provided custom name of the Protection Source.",
															},
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the environment (such as 'kVMware' or 'kSQL') where the Protection Source exists. Depending on the environment, one of the following Protection Sources are initialized.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies an id of the Protection Source.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a name of the Protection Source.",
															},
															"parent_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies an id of the parent of the Protection Source.",
															},
															"physical_protection_source": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a Protection Source in a Physical environment.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"agents": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifiles the agents running on the Physical Protection Source and the status information.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"cbmr_version": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the version if Cristie BMR product is installed on the host.",
																					},
																					"file_cbt_info": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "CBT version and service state info.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"file_version": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"build_ver": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																											"major_ver": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																											"minor_ver": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																											"revision_num": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																										},
																									},
																								},
																								"is_installed": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Indicates whether the cbt driver is installed.",
																								},
																								"reboot_status": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Indicates whether host is rebooted post VolCBT installation.",
																								},
																								"service_state": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Structure to Hold Service Status.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": &schema.Schema{
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																											"state": &schema.Schema{
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"host_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the host type where the agent is running. This is only set for persistent agents.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the agent's id.",
																					},
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the agent's name.",
																					},
																					"oracle_multi_node_channel_supported": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether oracle multi node multi channel is supported or not.",
																					},
																					"registration_info": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies information about a registered Source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"access_info": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the parameters required to establish a connection with a particular environment.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"connection_id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
																											},
																											"connector_group_id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
																											},
																											"endpoint": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
																											},
																											"environment": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																											},
																											"id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
																											},
																											"version": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
																											},
																										},
																									},
																								},
																								"allowed_ip_addresses": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"authentication_error_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
																								},
																								"authentication_status": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
																								},
																								"blacklisted_ip_addresses": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "This field is deprecated. Use DeniedIpAddresses instead.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"denied_ip_addresses": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"environments": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies a list of applications environment that are registered with this Protection Source such as 'kSQL'. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"is_db_authenticated": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if application entity dbAuthenticated or not.",
																								},
																								"is_storage_array_snapshot_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if this source entity has enabled storage array snapshot or not.",
																								},
																								"link_vms_across_vcenter": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
																								},
																								"minimum_free_space_gb": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																								},
																								"minimum_free_space_percent": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																								},
																								"password": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies password of the username to access the target source.",
																								},
																								"physical_params": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"applications": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"password": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies password of the username to access the target source.",
																											},
																											"throttling_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the source side throttling configuration.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"cpu_throttling_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Throttling Configuration Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"fixed_threshold": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																																	},
																																	"pattern_type": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																																	},
																																	"throttling_windows": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Throttling Window Parameters Definition.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day_time_window": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Day Time Window Parameters.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"end_time": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies the Day Time Parameters.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"day": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
																																											Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																										},
																																										"time": &schema.Schema{
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
																																							"start_time": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies the Day Time Parameters.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"day": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
																																											Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																										},
																																										"time": &schema.Schema{
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
																																						},
																																					},
																																				},
																																				"threshold": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "Throttling threshold applicable in the window.",
																																				},
																																			},
																																		},
																																	},
																																},
																															},
																														},
																														"network_throttling_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Throttling Configuration Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"fixed_threshold": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																																	},
																																	"pattern_type": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																																	},
																																	"throttling_windows": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Throttling Window Parameters Definition.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day_time_window": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Day Time Window Parameters.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"end_time": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies the Day Time Parameters.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"day": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
																																											Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																										},
																																										"time": &schema.Schema{
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
																																							"start_time": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies the Day Time Parameters.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"day": &schema.Schema{
																																											Type:        schema.TypeString,
																																											Computed:    true,
																																											Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																										},
																																										"time": &schema.Schema{
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
																																						},
																																					},
																																				},
																																				"threshold": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "Throttling threshold applicable in the window.",
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
																											"username": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies username to access the target source.",
																											},
																										},
																									},
																								},
																								"progress_monitor_path": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
																								},
																								"refresh_error_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
																								},
																								"refresh_time_usecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
																								},
																								"registered_apps_info": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies information of the applications registered on this protection source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"authentication_error_message": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
																											},
																											"authentication_status": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
																											},
																											"environment": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																											},
																											"host_settings_check_results": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"check_type": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																														},
																														"result_type": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																														},
																														"user_message": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies a descriptive message for failed/warning types.",
																														},
																													},
																												},
																											},
																											"refresh_error_message": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
																											},
																										},
																									},
																								},
																								"registration_time_usecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
																								},
																								"subnets": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"component": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Component that has reserved the subnet.",
																											},
																											"description": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Description of the subnet.",
																											},
																											"id": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "ID of the subnet.",
																											},
																											"ip": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies either an IPv6 address or an IPv4 address.",
																											},
																											"netmask_bits": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "netmaskBits.",
																											},
																											"netmask_ip4": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
																											},
																											"nfs_access": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Component that has reserved the subnet.",
																											},
																											"nfs_all_squash": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
																											},
																											"nfs_root_squash": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether clients from this subnet can mount as root on NFS.",
																											},
																											"s3_access": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																											},
																											"smb_access": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																											},
																											"tenant_id": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the unique id of the tenant.",
																											},
																										},
																									},
																								},
																								"throttling_policy": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the throttling policy for a registered Protection Source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"enforce_max_streams": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																											},
																											"enforce_registered_source_max_backups": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																											},
																											"is_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																											},
																											"latency_thresholds": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"active_task_msecs": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																														},
																														"new_task_msecs": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																														},
																													},
																												},
																											},
																											"max_concurrent_streams": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																											},
																											"nas_source_params": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																														},
																														"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																														},
																														"max_parallel_read_write_full_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																														},
																														"max_parallel_read_write_incremental_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																														},
																													},
																												},
																											},
																											"registered_source_max_concurrent_backups": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																											},
																											"storage_array_snapshot_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Configuration.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"is_max_snapshots_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																														},
																														"is_max_space_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																														},
																														"storage_array_snapshot_max_space_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Space Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshot_space_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
																																	},
																																},
																															},
																														},
																														"storage_array_snapshot_throttling_policies": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies throttling policies configured for individual volume/lun.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"id": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "Specifies the volume id of the storage array snapshot config.",
																																	},
																																	"is_max_snapshots_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																	},
																																	"is_max_space_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																	},
																																	"max_snapshot_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_snapshots": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Max number of storage snapshots allowed per volume/lun.",
																																				},
																																			},
																																		},
																																	},
																																	"max_space_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Max Space Config.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_snapshot_space_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Max number of storage snapshots allowed per volume/lun.",
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
																								"throttling_policy_overrides": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies throttling policy override for a Datastore in a registered entity.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"datastore_id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the Protection Source id of the Datastore.",
																											},
																											"datastore_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the display name of the Datastore.",
																											},
																											"throttling_policy": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the throttling policy for a registered Protection Source.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"enforce_max_streams": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																														},
																														"enforce_registered_source_max_backups": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																														},
																														"is_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																														},
																														"latency_thresholds": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"active_task_msecs": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																																	},
																																	"new_task_msecs": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																																	},
																																},
																															},
																														},
																														"max_concurrent_streams": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																														},
																														"nas_source_params": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																																	},
																																	"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																																	},
																																	"max_parallel_read_write_full_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																																	},
																																	"max_parallel_read_write_incremental_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																																	},
																																},
																															},
																														},
																														"registered_source_max_concurrent_backups": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																														},
																														"storage_array_snapshot_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Configuration.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"is_max_snapshots_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																	},
																																	"is_max_space_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																	},
																																	"storage_array_snapshot_max_space_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Max Space Config.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_snapshot_space_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Max number of storage snapshots allowed per volume/lun.",
																																				},
																																			},
																																		},
																																	},
																																	"storage_array_snapshot_throttling_policies": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies throttling policies configured for individual volume/lun.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"id": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "Specifies the volume id of the storage array snapshot config.",
																																				},
																																				"is_max_snapshots_config_enabled": &schema.Schema{
																																					Type:        schema.TypeBool,
																																					Computed:    true,
																																					Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																				},
																																				"is_max_space_config_enabled": &schema.Schema{
																																					Type:        schema.TypeBool,
																																					Computed:    true,
																																					Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																				},
																																				"max_snapshot_config": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"max_snapshots": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Computed:    true,
																																								Description: "Max number of storage snapshots allowed per volume/lun.",
																																							},
																																						},
																																					},
																																				},
																																				"max_space_config": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies Storage Array Snapshot Max Space Config.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"max_snapshot_space_percentage": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Computed:    true,
																																								Description: "Max number of storage snapshots allowed per volume/lun.",
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
																									},
																								},
																								"use_o_auth_for_exchange_online": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
																								},
																								"use_vm_bios_uuid": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
																								},
																								"user_messages": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"username": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies username to access the target source.",
																								},
																								"vlan_params": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the VLAN configuration for Recovery.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"vlan": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																											},
																											"disable_vlan": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
																											},
																											"interface_name": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																											},
																										},
																									},
																								},
																								"warning_messages": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																							},
																						},
																					},
																					"source_side_dedup_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether source side dedup is enabled or not.",
																					},
																					"status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the agent status. Specifies the status of the agent running on a physical source.",
																					},
																					"status_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies additional details about the agent status.",
																					},
																					"upgradability": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the upgradability of the agent running on the physical server. Specifies the upgradability of the agent running on the physical server.",
																					},
																					"upgrade_status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the status of the upgrade of the agent on a physical server. Specifies the status of the upgrade of the agent on a physical server.",
																					},
																					"upgrade_status_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies detailed message about the agent upgrade failure. This field is not set for successful upgrade.",
																					},
																					"version": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the version of the Agent software.",
																					},
																					"vol_cbt_info": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "CBT version and service state info.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"file_version": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"build_ver": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																											"major_ver": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																											"minor_ver": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																											"revision_num": &schema.Schema{
																												Type:     schema.TypeFloat,
																												Computed: true,
																											},
																										},
																									},
																								},
																								"is_installed": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Indicates whether the cbt driver is installed.",
																								},
																								"reboot_status": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Indicates whether host is rebooted post VolCBT installation.",
																								},
																								"service_state": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Structure to Hold Service Status.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"name": &schema.Schema{
																												Type:     schema.TypeString,
																												Computed: true,
																											},
																											"state": &schema.Schema{
																												Type:     schema.TypeString,
																												Computed: true,
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
																		"cluster_source_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of cluster resource this source represents.",
																		},
																		"host_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the hostname.",
																		},
																		"host_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the environment type for the host.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies an id for an object that is unique across Cohesity Clusters. The id is composite of all the ids listed below.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"cluster_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Cohesity Cluster id where the object was created.",
																					},
																					"cluster_incarnation_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies an id for the Cohesity Cluster that is generated when a Cohesity Cluster is initially created.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a unique id assigned to an object (such as a Job) by the Cohesity Cluster.",
																					},
																				},
																			},
																		},
																		"is_proxy_host": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the physical host is a proxy host.",
																		},
																		"memory_size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total memory on the host in bytes.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a human readable name of the Protection Source.",
																		},
																		"networking_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the struct containing information about network addresses configured on the given box. This is needed for dealing with Windows/Oracle Cluster resources that we discover and protect automatically.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"resource_vec": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The list of resources on the system that are accessible by an IP address.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"endpoints": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The endpoints by which the resource is accessible.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"fqdn": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The Fully Qualified Domain Name.",
																											},
																											"ipv4_addr": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The IPv4 address.",
																											},
																											"ipv6_addr": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The IPv6 address.",
																											},
																										},
																									},
																								},
																								"type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The type of the resource.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"num_processors": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the number of processors on the host.",
																		},
																		"os_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a human readable name of the OS of the Protection Source.",
																		},
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of managed Object in a Physical Protection Source. 'kGroup' indicates the EH container.",
																		},
																		"vcs_version": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies cluster version for VCS host.",
																		},
																		"volumes": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Array of Physical Volumes. Specifies the volumes available on the physical host. These fields are populated only for the kPhysicalHost type.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"device_path": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the path to the device that hosts the volume locally.",
																					},
																					"guid": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies an id for the Physical Volume.",
																					},
																					"is_boot_volume": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether the volume is boot volume.",
																					},
																					"is_extended_attributes_supported": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether this volume supports extended attributes (like ACLs) when performing file backups.",
																					},
																					"is_protected": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if a volume is protected by a Job.",
																					},
																					"is_shared_volume": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether the volume is shared volume.",
																					},
																					"label": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a volume label that can be used for displaying additional identifying information about a volume.",
																					},
																					"logical_size_bytes": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the logical size of the volume in bytes that is not reduced by change-block tracking, compression and deduplication.",
																					},
																					"mount_points": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the mount points where the volume is mounted, for example- 'C:', '/mnt/foo' etc.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"mount_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies mount type of volume e.g. nfs, autofs, ext4 etc.",
																					},
																					"network_path": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the full path to connect to the network attached volume. For example, (IP or hostname):/path/to/share for NFS volumes).",
																					},
																					"used_size_bytes": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the size used by the volume in bytes.",
																					},
																				},
																			},
																		},
																		"vsswriters": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies vss writer information about a Physical Protection Source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"is_writer_excluded": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "If true, the writer will be excluded by default.",
																					},
																					"writer_name": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies the name of the writer.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"sql_protection_source": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies an Object representing one SQL Server instance or database.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_available_for_vss_backup": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the database is marked as available for backup according to the SQL Server VSS writer. This may be false if either the state of the databases is not online, or if the VSS writer is not online. This field is set only for type 'kDatabase'.",
																		},
																		"created_timestamp": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the time when the database was created. It is displayed in the timezone of the SQL server on which this database is running.",
																		},
																		"database_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the database name of the SQL Protection Source, if the type is database.",
																		},
																		"db_aag_entity_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the AAG entity id if the database is part of an AAG. This field is set only for type 'kDatabase'.",
																		},
																		"db_aag_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the name of the AAG if the database is part of an AAG. This field is set only for type 'kDatabase'.",
																		},
																		"db_compatibility_level": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the versions of SQL server that the database is compatible with.",
																		},
																		"db_file_groups": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the information about the set of file groups for this db on the host. This is only set if the type is kDatabase.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"db_files": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the last known information about the set of database files on the host. This field is set only for type 'kDatabase'.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"file_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the format type of the file that SQL database stores the data. Specifies the format type of the file that SQL database stores the data. 'kRows' refers to a data file 'kLog' refers to a log file 'kFileStream' refers to a directory containing FILESTREAM data 'kNotSupportedType' is for information purposes only. Not supported. 'kFullText' refers to a full-text catalog.",
																					},
																					"full_path": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the full path of the database file on the SQL host machine.",
																					},
																					"size_bytes": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the last known size of the database file.",
																					},
																				},
																			},
																		},
																		"db_owner_username": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the name of the database owner.",
																		},
																		"default_database_location": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the default path for data files for DBs in an instance.",
																		},
																		"default_log_location": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the default path for log files for DBs in an instance.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a unique id for a SQL Protection Source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"created_date_msecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a unique identifier generated from the date the database is created or renamed. Cohesity uses this identifier in combination with the databaseId to uniquely identify a database.",
																					},
																					"database_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a unique id of the database but only for the life of the database. SQL Server may reuse database ids. Cohesity uses the createDateMsecs in combination with this databaseId to uniquely identify a database.",
																					},
																					"instance_id": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies unique id for the SQL Server instance. This id does not change during the life of the instance.",
																					},
																				},
																			},
																		},
																		"is_encrypted": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the database is TDE enabled.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the instance name of the SQL Protection Source.",
																		},
																		"owner_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the id of the container VM for the SQL Protection Source.",
																		},
																		"recovery_model": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the Recovery Model for the database in SQL environment. Only meaningful for the 'kDatabase' SQL Protection Source. Specifies the Recovery Model set for the Microsoft SQL Server. 'kSimpleRecoveryModel' indicates the Simple SQL Recovery Model which does not utilize log backups. 'kFullRecoveryModel' indicates the Full SQL Recovery Model which requires log backups and allows recovery to a single point in time. 'kBulkLoggedRecoveryModel' indicates the Bulk Logged SQL Recovery Model which requires log backups and allows high-performance bulk copy operations.",
																		},
																		"sql_server_db_state": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The state of the database as returned by SQL Server. Indicates the state of the database. The values correspond to the 'state' field in the system table sys.databases. See https://goo.gl/P66XqM. 'kOnline' indicates that database is in online state. 'kRestoring' indicates that database is in restore state. 'kRecovering' indicates that database is in recovery state. 'kRecoveryPending' indicates that database recovery is in pending state. 'kSuspect' indicates that primary filegroup is suspect and may be damaged. 'kEmergency' indicates that manually forced emergency state. 'kOffline' indicates that database is in offline state. 'kCopying' indicates that database is in copying state. 'kOfflineSecondary' indicates that secondary database is in offline state.",
																		},
																		"sql_server_instance_version": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the Server Instance Version.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"build": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the build.",
																					},
																					"major_version": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the major version.",
																					},
																					"minor_version": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the minor version.",
																					},
																					"revision": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the revision.",
																					},
																					"version_string": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the version string.",
																					},
																				},
																			},
																		},
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of the managed Object in a SQL Protection Source. Examples of SQL Objects include 'kInstance' and 'kDatabase'. 'kInstance' indicates that SQL server instance is being protected. 'kDatabase' indicates that SQL server database is being protected. 'kAAG' indicates that SQL AAG (AlwaysOn Availability Group) is being protected. 'kAAGRootContainer' indicates that SQL AAG's root container is being protected. 'kRootContainer' indicates root container for SQL sources.",
																		},
																	},
																},
															},
														},
													},
												},
												"registration_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies information about a registered Source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"access_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the parameters required to establish a connection with a particular environment.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"connection_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
																		},
																		"connector_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
																		},
																		"endpoint": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
																		},
																		"environment": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
																		},
																		"version": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
																		},
																	},
																},
															},
															"allowed_ip_addresses": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"authentication_error_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
															},
															"authentication_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
															},
															"blacklisted_ip_addresses": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "This field is deprecated. Use DeniedIpAddresses instead.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"denied_ip_addresses": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"environments": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a list of applications environment that are registered with this Protection Source such as 'kSQL'. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"is_db_authenticated": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if application entity dbAuthenticated or not.",
															},
															"is_storage_array_snapshot_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if this source entity has enabled storage array snapshot or not.",
															},
															"link_vms_across_vcenter": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
															},
															"minimum_free_space_gb": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
															},
															"minimum_free_space_percent": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
															},
															"password": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies password of the username to access the target source.",
															},
															"physical_params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"applications": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies password of the username to access the target source.",
																		},
																		"throttling_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the source side throttling configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"cpu_throttling_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Throttling Configuration Parameters.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"fixed_threshold": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																								},
																								"pattern_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																								},
																								"throttling_windows": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Throttling Window Parameters Definition.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"day_time_window": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Day Time Window Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"end_time": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"day": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																	},
																																	"time": &schema.Schema{
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
																														"start_time": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"day": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																	},
																																	"time": &schema.Schema{
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
																													},
																												},
																											},
																											"threshold": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Throttling threshold applicable in the window.",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"network_throttling_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Throttling Configuration Parameters.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"fixed_threshold": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																								},
																								"pattern_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																								},
																								"throttling_windows": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Throttling Window Parameters Definition.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"day_time_window": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Day Time Window Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"end_time": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"day": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																	},
																																	"time": &schema.Schema{
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
																														"start_time": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"day": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																	},
																																	"time": &schema.Schema{
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
																													},
																												},
																											},
																											"threshold": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Throttling threshold applicable in the window.",
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
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies username to access the target source.",
																		},
																	},
																},
															},
															"progress_monitor_path": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
															},
															"refresh_error_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
															},
															"refresh_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
															},
															"registered_apps_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies information of the applications registered on this protection source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"authentication_error_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
																		},
																		"authentication_status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
																		},
																		"environment": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																		},
																		"host_settings_check_results": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"check_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																					},
																					"result_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																					},
																					"user_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a descriptive message for failed/warning types.",
																					},
																				},
																			},
																		},
																		"refresh_error_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
																		},
																	},
																},
															},
															"registration_time_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
															},
															"subnets": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"component": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Component that has reserved the subnet.",
																		},
																		"description": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the subnet.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "ID of the subnet.",
																		},
																		"ip": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies either an IPv6 address or an IPv4 address.",
																		},
																		"netmask_bits": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "netmaskBits.",
																		},
																		"netmask_ip4": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
																		},
																		"nfs_access": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Component that has reserved the subnet.",
																		},
																		"nfs_all_squash": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
																		},
																		"nfs_root_squash": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether clients from this subnet can mount as root on NFS.",
																		},
																		"s3_access": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																		},
																		"smb_access": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																		},
																		"tenant_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unique id of the tenant.",
																		},
																	},
																},
															},
															"throttling_policy": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the throttling policy for a registered Protection Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enforce_max_streams": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																		},
																		"enforce_registered_source_max_backups": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																		},
																		"is_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																		},
																		"latency_thresholds": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"active_task_msecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																					},
																					"new_task_msecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																					},
																				},
																			},
																		},
																		"max_concurrent_streams": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																		},
																		"nas_source_params": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																					},
																					"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																					},
																					"max_parallel_read_write_full_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																					},
																					"max_parallel_read_write_incremental_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																					},
																				},
																			},
																		},
																		"registered_source_max_concurrent_backups": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																		},
																		"storage_array_snapshot_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies Storage Array Snapshot Configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"is_max_snapshots_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																					},
																					"is_max_space_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																					},
																					"storage_array_snapshot_max_space_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Max Space Config.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_snapshot_space_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Max number of storage snapshots allowed per volume/lun.",
																								},
																							},
																						},
																					},
																					"storage_array_snapshot_throttling_policies": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies throttling policies configured for individual volume/lun.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the volume id of the storage array snapshot config.",
																								},
																								"is_max_snapshots_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																								},
																								"is_max_space_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																								},
																								"max_snapshot_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshots": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
																											},
																										},
																									},
																								},
																								"max_space_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Space Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshot_space_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
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
															"throttling_policy_overrides": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies throttling policy override for a Datastore in a registered entity.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"datastore_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Protection Source id of the Datastore.",
																		},
																		"datastore_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the display name of the Datastore.",
																		},
																		"throttling_policy": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the throttling policy for a registered Protection Source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"enforce_max_streams": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																					},
																					"enforce_registered_source_max_backups": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																					},
																					"is_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																					},
																					"latency_thresholds": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"active_task_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																								},
																								"new_task_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																								},
																							},
																						},
																					},
																					"max_concurrent_streams": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																					},
																					"nas_source_params": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																								},
																								"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																								},
																								"max_parallel_read_write_full_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																								},
																								"max_parallel_read_write_incremental_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																								},
																							},
																						},
																					},
																					"registered_source_max_concurrent_backups": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																					},
																					"storage_array_snapshot_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_max_snapshots_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																								},
																								"is_max_space_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																								},
																								"storage_array_snapshot_max_space_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Space Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshot_space_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
																											},
																										},
																									},
																								},
																								"storage_array_snapshot_throttling_policies": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies throttling policies configured for individual volume/lun.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the volume id of the storage array snapshot config.",
																											},
																											"is_max_snapshots_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																											},
																											"is_max_space_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																											},
																											"max_snapshot_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshots": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
																														},
																													},
																												},
																											},
																											"max_space_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Space Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshot_space_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
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
																},
															},
															"use_o_auth_for_exchange_online": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
															},
															"use_vm_bios_uuid": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
															},
															"user_messages": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"username": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies username to access the target source.",
															},
															"vlan_params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the VLAN configuration for Recovery.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"vlan": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																		},
																		"disable_vlan": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
																		},
																		"interface_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																		},
																	},
																},
															},
															"warning_messages": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"total_downtiered_size_in_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the total bytes downtiered from the source so far.",
												},
												"total_uptiered_size_in_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the total bytes uptiered to the source so far.",
												},
												"unprotected_sources_summary": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Aggregated information about a node subtree.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
															},
															"leaves_count": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the number of leaf nodes under the subtree of this node.",
															},
															"total_logical_size": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total logical size of the data under the subtree of this node.",
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
						"entity_pagination_parameters": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the cursor based pagination parameters for Protection Source and its children. Pagination is supported at a given level within the Protection Source Hierarchy with the help of before or after cursors. A Cursor will always refer to a specific source within the source dataset but will be invalidated if the item is removed.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"after_cursor_entity_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the entity id starting from which the items are to be returned.",
									},
									"before_cursor_entity_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the entity id upto which the items are to be returned.",
									},
									"node_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the entity id for the Node at any level within the Source entity hierarchy whose children are to be paginated.",
									},
									"page_size": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the maximum number of entities to be returned within the page.",
									},
								},
							},
						},
						"entity_permission_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the permission information of entities.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"entity_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the entity id.",
									},
									"groups": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies struct with basic group details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"domain": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies domain name of the user.",
												},
												"group_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies group name of the group.",
												},
												"sid": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies unique Security ID (SID) of the user.",
												},
												"tenant_ids": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the tenants to which the group belongs to.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"is_inferred": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether the Entity Permission Information is inferred or not. For example, SQL application hosted over vCenter will have inferred entity permission information.",
									},
									"is_registered_by_sp": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether this entity is registered by the SP or not. This will be populated only if the entity is a root entity. Refer to magneto/base/permissions.proto for details.",
									},
									"registering_tenant_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the tenant id that registered this entity. This will be populated only if the entity is a root entity.",
									},
									"tenant": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies struct with basic tenant details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"bifrost_enabled": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if this tenant is bifrost enabled or not.",
												},
												"is_managed_on_helios": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether this tenant is manged on helios.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies name of the tenant.",
												},
												"tenant_id": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies the unique id of the tenant.",
												},
											},
										},
									},
									"users": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies struct with basic user details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"domain": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies domain name of the user.",
												},
												"sid": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies unique Security ID (SID) of the user.",
												},
												"tenant_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the tenant to which the user belongs to.",
												},
												"user_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies user name of the user.",
												},
											},
										},
									},
								},
							},
						},
						"logical_size": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the logical size of the data in bytes for the Object on this node. Presence of this field indicates this node is a leaf node.",
						},
						"nodes": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies children of the current node in the Protection Sources hierarchy. When representing Objects in memory, the entire Object subtree hierarchy is represented. You can use this subtree to navigate down the Object hierarchy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"application_nodes": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the child subtree used to store additional application-level Objects. Different environments use the subtree to store application-level information. For example for SQL Server, this subtree stores the SQL Server instances running on a VM.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"nodes": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies children of the current node in the Protection Sources hierarchy.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"application_nodes": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the child subtree used to store additional application-level Objects. Different environments use the subtree to store application-level information. For example for SQL Server, this subtree stores the SQL Server instances running on a VM.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
															"entity_pagination_parameters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the cursor based pagination parameters for Protection Source and its children. Pagination is supported at a given level within the Protection Source Hierarchy with the help of before or after cursors. A Cursor will always refer to a specific source within the source dataset but will be invalidated if the item is removed.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"after_cursor_entity_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the entity id starting from which the items are to be returned.",
																		},
																		"before_cursor_entity_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the entity id upto which the items are to be returned.",
																		},
																		"node_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the entity id for the Node at any level within the Source entity hierarchy whose children are to be paginated.",
																		},
																		"page_size": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the maximum number of entities to be returned within the page.",
																		},
																	},
																},
															},
															"entity_permission_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the permission information of entities.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"entity_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the entity id.",
																		},
																		"groups": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies struct with basic group details.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"domain": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies domain name of the user.",
																					},
																					"group_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies group name of the group.",
																					},
																					"sid": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies unique Security ID (SID) of the user.",
																					},
																					"tenant_ids": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the tenants to which the group belongs to.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																				},
																			},
																		},
																		"is_inferred": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the Entity Permission Information is inferred or not. For example, SQL application hosted over vCenter will have inferred entity permission information.",
																		},
																		"is_registered_by_sp": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether this entity is registered by the SP or not. This will be populated only if the entity is a root entity. Refer to magneto/base/permissions.proto for details.",
																		},
																		"registering_tenant_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the tenant id that registered this entity. This will be populated only if the entity is a root entity.",
																		},
																		"tenant": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies struct with basic tenant details.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"bifrost_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if this tenant is bifrost enabled or not.",
																					},
																					"is_managed_on_helios": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether this tenant is manged on helios.",
																					},
																					"name": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies name of the tenant.",
																					},
																					"tenant_id": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies the unique id of the tenant.",
																					},
																				},
																			},
																		},
																		"users": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies struct with basic user details.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"domain": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies domain name of the user.",
																					},
																					"sid": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies unique Security ID (SID) of the user.",
																					},
																					"tenant_id": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the tenant to which the user belongs to.",
																					},
																					"user_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies user name of the user.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"logical_size": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the logical size of the data in bytes for the Object on this node. Presence of this field indicates this node is a leaf node.",
															},
															"nodes": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies children of the current node in the Protection Sources hierarchy. When representing Objects in memory, the entire Object subtree hierarchy is represented. You can use this subtree to navigate down the Object hierarchy.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{},
																},
															},
															"object_protection_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the Object Protection Info of the Protection Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"auto_protect_parent_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the auto protect parent id if this entity is protected based on auto protection. This is only specified for leaf entities.",
																		},
																		"entity_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the entity id.",
																		},
																		"has_active_object_protection_spec": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies if the entity is under object protection.",
																		},
																	},
																},
															},
															"protected_sources_summary": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Protected Objects. Specifies aggregated information about all the child Objects of this node that are currently protected by a Protection Job. There is one entry for each environment that is being backed up. The aggregated information for the Object hierarchy's environment will be available at the 0th index of the vector.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"environment": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
																		},
																		"leaves_count": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the number of leaf nodes under the subtree of this node.",
																		},
																		"total_logical_size": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total logical size of the data under the subtree of this node.",
																		},
																	},
																},
															},
															"protection_source": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies details about an Acropolis Protection Source when the environment is set to 'kAcropolis'.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"connection_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the connection id of the tenant.",
																		},
																		"connector_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the connector group id of the connector groups.",
																		},
																		"custom_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the user provided custom name of the Protection Source.",
																		},
																		"environment": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the environment (such as 'kVMware' or 'kSQL') where the Protection Source exists. Depending on the environment, one of the following Protection Sources are initialized.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies an id of the Protection Source.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a name of the Protection Source.",
																		},
																		"parent_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies an id of the parent of the Protection Source.",
																		},
																		"physical_protection_source": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a Protection Source in a Physical environment.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"agents": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifiles the agents running on the Physical Protection Source and the status information.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"cbmr_version": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the version if Cristie BMR product is installed on the host.",
																								},
																								"file_cbt_info": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "CBT version and service state info.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"file_version": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"build_ver": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																														"major_ver": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																														"minor_ver": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																														"revision_num": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																													},
																												},
																											},
																											"is_installed": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Indicates whether the cbt driver is installed.",
																											},
																											"reboot_status": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Indicates whether host is rebooted post VolCBT installation.",
																											},
																											"service_state": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Structure to Hold Service Status.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"name": &schema.Schema{
																															Type:     schema.TypeString,
																															Computed: true,
																														},
																														"state": &schema.Schema{
																															Type:     schema.TypeString,
																															Computed: true,
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"host_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the host type where the agent is running. This is only set for persistent agents.",
																								},
																								"id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the agent's id.",
																								},
																								"name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the agent's name.",
																								},
																								"oracle_multi_node_channel_supported": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether oracle multi node multi channel is supported or not.",
																								},
																								"registration_info": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies information about a registered Source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"access_info": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the parameters required to establish a connection with a particular environment.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"connection_id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
																														},
																														"connector_group_id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
																														},
																														"endpoint": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
																														},
																														"environment": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																														},
																														"id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
																														},
																														"version": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
																														},
																													},
																												},
																											},
																											"allowed_ip_addresses": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"authentication_error_message": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
																											},
																											"authentication_status": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
																											},
																											"blacklisted_ip_addresses": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "This field is deprecated. Use DeniedIpAddresses instead.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"denied_ip_addresses": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"environments": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies a list of applications environment that are registered with this Protection Source such as 'kSQL'. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"is_db_authenticated": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if application entity dbAuthenticated or not.",
																											},
																											"is_storage_array_snapshot_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if this source entity has enabled storage array snapshot or not.",
																											},
																											"link_vms_across_vcenter": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
																											},
																											"minimum_free_space_gb": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																											},
																											"minimum_free_space_percent": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																											},
																											"password": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies password of the username to access the target source.",
																											},
																											"physical_params": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"applications": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																															Elem: &schema.Schema{
																																Type: schema.TypeString,
																															},
																														},
																														"password": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies password of the username to access the target source.",
																														},
																														"throttling_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the source side throttling configuration.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"cpu_throttling_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Throttling Configuration Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"fixed_threshold": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																																				},
																																				"pattern_type": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																																				},
																																				"throttling_windows": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Throttling Window Parameters Definition.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"day_time_window": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies the Day Time Window Parameters.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"end_time": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Specifies the Day Time Parameters.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"day": &schema.Schema{
																																														Type:        schema.TypeString,
																																														Computed:    true,
																																														Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																													},
																																													"time": &schema.Schema{
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
																																										"start_time": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Specifies the Day Time Parameters.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"day": &schema.Schema{
																																														Type:        schema.TypeString,
																																														Computed:    true,
																																														Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																													},
																																													"time": &schema.Schema{
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
																																									},
																																								},
																																							},
																																							"threshold": &schema.Schema{
																																								Type:        schema.TypeInt,
																																								Computed:    true,
																																								Description: "Throttling threshold applicable in the window.",
																																							},
																																						},
																																					},
																																				},
																																			},
																																		},
																																	},
																																	"network_throttling_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Throttling Configuration Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"fixed_threshold": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																																				},
																																				"pattern_type": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																																				},
																																				"throttling_windows": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Throttling Window Parameters Definition.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"day_time_window": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies the Day Time Window Parameters.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"end_time": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Specifies the Day Time Parameters.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"day": &schema.Schema{
																																														Type:        schema.TypeString,
																																														Computed:    true,
																																														Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																													},
																																													"time": &schema.Schema{
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
																																										"start_time": &schema.Schema{
																																											Type:        schema.TypeList,
																																											Computed:    true,
																																											Description: "Specifies the Day Time Parameters.",
																																											Elem: &schema.Resource{
																																												Schema: map[string]*schema.Schema{
																																													"day": &schema.Schema{
																																														Type:        schema.TypeString,
																																														Computed:    true,
																																														Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																													},
																																													"time": &schema.Schema{
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
																																									},
																																								},
																																							},
																																							"threshold": &schema.Schema{
																																								Type:        schema.TypeInt,
																																								Computed:    true,
																																								Description: "Throttling threshold applicable in the window.",
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
																														"username": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies username to access the target source.",
																														},
																													},
																												},
																											},
																											"progress_monitor_path": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
																											},
																											"refresh_error_message": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
																											},
																											"refresh_time_usecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
																											},
																											"registered_apps_info": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies information of the applications registered on this protection source.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"authentication_error_message": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
																														},
																														"authentication_status": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
																														},
																														"environment": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																														},
																														"host_settings_check_results": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"check_type": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																																	},
																																	"result_type": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																																	},
																																	"user_message": &schema.Schema{
																																		Type:        schema.TypeString,
																																		Computed:    true,
																																		Description: "Specifies a descriptive message for failed/warning types.",
																																	},
																																},
																															},
																														},
																														"refresh_error_message": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
																														},
																													},
																												},
																											},
																											"registration_time_usecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
																											},
																											"subnets": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"component": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Component that has reserved the subnet.",
																														},
																														"description": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Description of the subnet.",
																														},
																														"id": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "ID of the subnet.",
																														},
																														"ip": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies either an IPv6 address or an IPv4 address.",
																														},
																														"netmask_bits": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "netmaskBits.",
																														},
																														"netmask_ip4": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
																														},
																														"nfs_access": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Component that has reserved the subnet.",
																														},
																														"nfs_all_squash": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
																														},
																														"nfs_root_squash": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether clients from this subnet can mount as root on NFS.",
																														},
																														"s3_access": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																														},
																														"smb_access": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																														},
																														"tenant_id": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the unique id of the tenant.",
																														},
																													},
																												},
																											},
																											"throttling_policy": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the throttling policy for a registered Protection Source.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"enforce_max_streams": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																														},
																														"enforce_registered_source_max_backups": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																														},
																														"is_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																														},
																														"latency_thresholds": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"active_task_msecs": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																																	},
																																	"new_task_msecs": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																																	},
																																},
																															},
																														},
																														"max_concurrent_streams": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																														},
																														"nas_source_params": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																																	},
																																	"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																																	},
																																	"max_parallel_read_write_full_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																																	},
																																	"max_parallel_read_write_incremental_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																																	},
																																},
																															},
																														},
																														"registered_source_max_concurrent_backups": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																														},
																														"storage_array_snapshot_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Configuration.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"is_max_snapshots_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																	},
																																	"is_max_space_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																	},
																																	"storage_array_snapshot_max_space_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Max Space Config.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_snapshot_space_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Max number of storage snapshots allowed per volume/lun.",
																																				},
																																			},
																																		},
																																	},
																																	"storage_array_snapshot_throttling_policies": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies throttling policies configured for individual volume/lun.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"id": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "Specifies the volume id of the storage array snapshot config.",
																																				},
																																				"is_max_snapshots_config_enabled": &schema.Schema{
																																					Type:        schema.TypeBool,
																																					Computed:    true,
																																					Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																				},
																																				"is_max_space_config_enabled": &schema.Schema{
																																					Type:        schema.TypeBool,
																																					Computed:    true,
																																					Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																				},
																																				"max_snapshot_config": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"max_snapshots": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Computed:    true,
																																								Description: "Max number of storage snapshots allowed per volume/lun.",
																																							},
																																						},
																																					},
																																				},
																																				"max_space_config": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies Storage Array Snapshot Max Space Config.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"max_snapshot_space_percentage": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Computed:    true,
																																								Description: "Max number of storage snapshots allowed per volume/lun.",
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
																											"throttling_policy_overrides": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies throttling policy override for a Datastore in a registered entity.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"datastore_id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Specifies the Protection Source id of the Datastore.",
																														},
																														"datastore_name": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the display name of the Datastore.",
																														},
																														"throttling_policy": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the throttling policy for a registered Protection Source.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"enforce_max_streams": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																																	},
																																	"enforce_registered_source_max_backups": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																																	},
																																	"is_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																																	},
																																	"latency_thresholds": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"active_task_msecs": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																																				},
																																				"new_task_msecs": &schema.Schema{
																																					Type:        schema.TypeInt,
																																					Computed:    true,
																																					Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																																				},
																																			},
																																		},
																																	},
																																	"max_concurrent_streams": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																																	},
																																	"nas_source_params": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																																				},
																																				"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																																				},
																																				"max_parallel_read_write_full_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																																				},
																																				"max_parallel_read_write_incremental_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																																				},
																																			},
																																		},
																																	},
																																	"registered_source_max_concurrent_backups": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																																	},
																																	"storage_array_snapshot_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Configuration.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"is_max_snapshots_config_enabled": &schema.Schema{
																																					Type:        schema.TypeBool,
																																					Computed:    true,
																																					Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																				},
																																				"is_max_space_config_enabled": &schema.Schema{
																																					Type:        schema.TypeBool,
																																					Computed:    true,
																																					Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																				},
																																				"storage_array_snapshot_max_space_config": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies Storage Array Snapshot Max Space Config.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"max_snapshot_space_percentage": &schema.Schema{
																																								Type:        schema.TypeFloat,
																																								Computed:    true,
																																								Description: "Max number of storage snapshots allowed per volume/lun.",
																																							},
																																						},
																																					},
																																				},
																																				"storage_array_snapshot_throttling_policies": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies throttling policies configured for individual volume/lun.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"id": &schema.Schema{
																																								Type:        schema.TypeInt,
																																								Computed:    true,
																																								Description: "Specifies the volume id of the storage array snapshot config.",
																																							},
																																							"is_max_snapshots_config_enabled": &schema.Schema{
																																								Type:        schema.TypeBool,
																																								Computed:    true,
																																								Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																							},
																																							"is_max_space_config_enabled": &schema.Schema{
																																								Type:        schema.TypeBool,
																																								Computed:    true,
																																								Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																							},
																																							"max_snapshot_config": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"max_snapshots": &schema.Schema{
																																											Type:        schema.TypeFloat,
																																											Computed:    true,
																																											Description: "Max number of storage snapshots allowed per volume/lun.",
																																										},
																																									},
																																								},
																																							},
																																							"max_space_config": &schema.Schema{
																																								Type:        schema.TypeList,
																																								Computed:    true,
																																								Description: "Specifies Storage Array Snapshot Max Space Config.",
																																								Elem: &schema.Resource{
																																									Schema: map[string]*schema.Schema{
																																										"max_snapshot_space_percentage": &schema.Schema{
																																											Type:        schema.TypeFloat,
																																											Computed:    true,
																																											Description: "Max number of storage snapshots allowed per volume/lun.",
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
																												},
																											},
																											"use_o_auth_for_exchange_online": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
																											},
																											"use_vm_bios_uuid": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
																											},
																											"user_messages": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																											"username": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies username to access the target source.",
																											},
																											"vlan_params": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the VLAN configuration for Recovery.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"vlan": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																														},
																														"disable_vlan": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
																														},
																														"interface_name": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																														},
																													},
																												},
																											},
																											"warning_messages": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
																												Elem: &schema.Schema{
																													Type: schema.TypeString,
																												},
																											},
																										},
																									},
																								},
																								"source_side_dedup_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether source side dedup is enabled or not.",
																								},
																								"status": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the agent status. Specifies the status of the agent running on a physical source.",
																								},
																								"status_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies additional details about the agent status.",
																								},
																								"upgradability": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the upgradability of the agent running on the physical server. Specifies the upgradability of the agent running on the physical server.",
																								},
																								"upgrade_status": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the status of the upgrade of the agent on a physical server. Specifies the status of the upgrade of the agent on a physical server.",
																								},
																								"upgrade_status_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies detailed message about the agent upgrade failure. This field is not set for successful upgrade.",
																								},
																								"version": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the version of the Agent software.",
																								},
																								"vol_cbt_info": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "CBT version and service state info.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"file_version": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"build_ver": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																														"major_ver": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																														"minor_ver": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																														"revision_num": &schema.Schema{
																															Type:     schema.TypeFloat,
																															Computed: true,
																														},
																													},
																												},
																											},
																											"is_installed": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Indicates whether the cbt driver is installed.",
																											},
																											"reboot_status": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Indicates whether host is rebooted post VolCBT installation.",
																											},
																											"service_state": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Structure to Hold Service Status.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"name": &schema.Schema{
																															Type:     schema.TypeString,
																															Computed: true,
																														},
																														"state": &schema.Schema{
																															Type:     schema.TypeString,
																															Computed: true,
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
																					"cluster_source_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of cluster resource this source represents.",
																					},
																					"host_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the hostname.",
																					},
																					"host_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the environment type for the host.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies an id for an object that is unique across Cohesity Clusters. The id is composite of all the ids listed below.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"cluster_id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the Cohesity Cluster id where the object was created.",
																								},
																								"cluster_incarnation_id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies an id for the Cohesity Cluster that is generated when a Cohesity Cluster is initially created.",
																								},
																								"id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies a unique id assigned to an object (such as a Job) by the Cohesity Cluster.",
																								},
																							},
																						},
																					},
																					"is_proxy_host": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the physical host is a proxy host.",
																					},
																					"memory_size_bytes": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the total memory on the host in bytes.",
																					},
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a human readable name of the Protection Source.",
																					},
																					"networking_info": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the struct containing information about network addresses configured on the given box. This is needed for dealing with Windows/Oracle Cluster resources that we discover and protect automatically.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"resource_vec": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "The list of resources on the system that are accessible by an IP address.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"endpoints": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "The endpoints by which the resource is accessible.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"fqdn": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The Fully Qualified Domain Name.",
																														},
																														"ipv4_addr": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The IPv4 address.",
																														},
																														"ipv6_addr": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "The IPv6 address.",
																														},
																													},
																												},
																											},
																											"type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "The type of the resource.",
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																					"num_processors": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the number of processors on the host.",
																					},
																					"os_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a human readable name of the OS of the Protection Source.",
																					},
																					"type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of managed Object in a Physical Protection Source. 'kGroup' indicates the EH container.",
																					},
																					"vcs_version": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies cluster version for VCS host.",
																					},
																					"volumes": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Array of Physical Volumes. Specifies the volumes available on the physical host. These fields are populated only for the kPhysicalHost type.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"device_path": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the path to the device that hosts the volume locally.",
																								},
																								"guid": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies an id for the Physical Volume.",
																								},
																								"is_boot_volume": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether the volume is boot volume.",
																								},
																								"is_extended_attributes_supported": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether this volume supports extended attributes (like ACLs) when performing file backups.",
																								},
																								"is_protected": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if a volume is protected by a Job.",
																								},
																								"is_shared_volume": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether the volume is shared volume.",
																								},
																								"label": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies a volume label that can be used for displaying additional identifying information about a volume.",
																								},
																								"logical_size_bytes": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the logical size of the volume in bytes that is not reduced by change-block tracking, compression and deduplication.",
																								},
																								"mount_points": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the mount points where the volume is mounted, for example- 'C:', '/mnt/foo' etc.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"mount_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies mount type of volume e.g. nfs, autofs, ext4 etc.",
																								},
																								"network_path": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the full path to connect to the network attached volume. For example, (IP or hostname):/path/to/share for NFS volumes).",
																								},
																								"used_size_bytes": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the size used by the volume in bytes.",
																								},
																							},
																						},
																					},
																					"vsswriters": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies vss writer information about a Physical Protection Source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_writer_excluded": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "If true, the writer will be excluded by default.",
																								},
																								"writer_name": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies the name of the writer.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"sql_protection_source": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies an Object representing one SQL Server instance or database.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"is_available_for_vss_backup": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether the database is marked as available for backup according to the SQL Server VSS writer. This may be false if either the state of the databases is not online, or if the VSS writer is not online. This field is set only for type 'kDatabase'.",
																					},
																					"created_timestamp": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the time when the database was created. It is displayed in the timezone of the SQL server on which this database is running.",
																					},
																					"database_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the database name of the SQL Protection Source, if the type is database.",
																					},
																					"db_aag_entity_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the AAG entity id if the database is part of an AAG. This field is set only for type 'kDatabase'.",
																					},
																					"db_aag_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the name of the AAG if the database is part of an AAG. This field is set only for type 'kDatabase'.",
																					},
																					"db_compatibility_level": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the versions of SQL server that the database is compatible with.",
																					},
																					"db_file_groups": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the information about the set of file groups for this db on the host. This is only set if the type is kDatabase.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"db_files": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the last known information about the set of database files on the host. This field is set only for type 'kDatabase'.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"file_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the format type of the file that SQL database stores the data. Specifies the format type of the file that SQL database stores the data. 'kRows' refers to a data file 'kLog' refers to a log file 'kFileStream' refers to a directory containing FILESTREAM data 'kNotSupportedType' is for information purposes only. Not supported. 'kFullText' refers to a full-text catalog.",
																								},
																								"full_path": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the full path of the database file on the SQL host machine.",
																								},
																								"size_bytes": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the last known size of the database file.",
																								},
																							},
																						},
																					},
																					"db_owner_username": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the name of the database owner.",
																					},
																					"default_database_location": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the default path for data files for DBs in an instance.",
																					},
																					"default_log_location": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the default path for log files for DBs in an instance.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a unique id for a SQL Protection Source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"created_date_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies a unique identifier generated from the date the database is created or renamed. Cohesity uses this identifier in combination with the databaseId to uniquely identify a database.",
																								},
																								"database_id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies a unique id of the database but only for the life of the database. SQL Server may reuse database ids. Cohesity uses the createDateMsecs in combination with this databaseId to uniquely identify a database.",
																								},
																								"instance_id": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies unique id for the SQL Server instance. This id does not change during the life of the instance.",
																								},
																							},
																						},
																					},
																					"is_encrypted": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether the database is TDE enabled.",
																					},
																					"name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the instance name of the SQL Protection Source.",
																					},
																					"owner_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the id of the container VM for the SQL Protection Source.",
																					},
																					"recovery_model": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the Recovery Model for the database in SQL environment. Only meaningful for the 'kDatabase' SQL Protection Source. Specifies the Recovery Model set for the Microsoft SQL Server. 'kSimpleRecoveryModel' indicates the Simple SQL Recovery Model which does not utilize log backups. 'kFullRecoveryModel' indicates the Full SQL Recovery Model which requires log backups and allows recovery to a single point in time. 'kBulkLoggedRecoveryModel' indicates the Bulk Logged SQL Recovery Model which requires log backups and allows high-performance bulk copy operations.",
																					},
																					"sql_server_db_state": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The state of the database as returned by SQL Server. Indicates the state of the database. The values correspond to the 'state' field in the system table sys.databases. See https://goo.gl/P66XqM. 'kOnline' indicates that database is in online state. 'kRestoring' indicates that database is in restore state. 'kRecovering' indicates that database is in recovery state. 'kRecoveryPending' indicates that database recovery is in pending state. 'kSuspect' indicates that primary filegroup is suspect and may be damaged. 'kEmergency' indicates that manually forced emergency state. 'kOffline' indicates that database is in offline state. 'kCopying' indicates that database is in copying state. 'kOfflineSecondary' indicates that secondary database is in offline state.",
																					},
																					"sql_server_instance_version": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Server Instance Version.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"build": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the build.",
																								},
																								"major_version": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the major version.",
																								},
																								"minor_version": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the minor version.",
																								},
																								"revision": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the revision.",
																								},
																								"version_string": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the version string.",
																								},
																							},
																						},
																					},
																					"type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of the managed Object in a SQL Protection Source. Examples of SQL Objects include 'kInstance' and 'kDatabase'. 'kInstance' indicates that SQL server instance is being protected. 'kDatabase' indicates that SQL server database is being protected. 'kAAG' indicates that SQL AAG (AlwaysOn Availability Group) is being protected. 'kAAGRootContainer' indicates that SQL AAG's root container is being protected. 'kRootContainer' indicates root container for SQL sources.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"registration_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies information about a registered Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"access_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the parameters required to establish a connection with a particular environment.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"connection_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
																					},
																					"connector_group_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
																					},
																					"endpoint": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
																					},
																					"environment": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
																					},
																					"version": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
																					},
																				},
																			},
																		},
																		"allowed_ip_addresses": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"authentication_error_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
																		},
																		"authentication_status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
																		},
																		"blacklisted_ip_addresses": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "This field is deprecated. Use DeniedIpAddresses instead.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"denied_ip_addresses": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"environments": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of applications environment that are registered with this Protection Source such as 'kSQL'. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"is_db_authenticated": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if application entity dbAuthenticated or not.",
																		},
																		"is_storage_array_snapshot_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if this source entity has enabled storage array snapshot or not.",
																		},
																		"link_vms_across_vcenter": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
																		},
																		"minimum_free_space_gb": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																		},
																		"minimum_free_space_percent": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies password of the username to access the target source.",
																		},
																		"physical_params": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"applications": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"password": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies password of the username to access the target source.",
																					},
																					"throttling_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the source side throttling configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"cpu_throttling_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Throttling Configuration Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"fixed_threshold": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																											},
																											"pattern_type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																											},
																											"throttling_windows": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Throttling Window Parameters Definition.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day_time_window": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Window Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"end_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																	"start_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																},
																															},
																														},
																														"threshold": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Throttling threshold applicable in the window.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"network_throttling_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Throttling Configuration Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"fixed_threshold": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																											},
																											"pattern_type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																											},
																											"throttling_windows": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Throttling Window Parameters Definition.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day_time_window": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Window Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"end_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																	"start_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																},
																															},
																														},
																														"threshold": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Throttling threshold applicable in the window.",
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
																					"username": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies username to access the target source.",
																					},
																				},
																			},
																		},
																		"progress_monitor_path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
																		},
																		"refresh_error_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
																		},
																		"refresh_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
																		},
																		"registered_apps_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies information of the applications registered on this protection source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"authentication_error_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
																					},
																					"authentication_status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
																					},
																					"environment": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																					},
																					"host_settings_check_results": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"check_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																								},
																								"result_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																								},
																								"user_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies a descriptive message for failed/warning types.",
																								},
																							},
																						},
																					},
																					"refresh_error_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
																					},
																				},
																			},
																		},
																		"registration_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
																		},
																		"subnets": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"component": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Component that has reserved the subnet.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Description of the subnet.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "ID of the subnet.",
																					},
																					"ip": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies either an IPv6 address or an IPv4 address.",
																					},
																					"netmask_bits": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "netmaskBits.",
																					},
																					"netmask_ip4": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
																					},
																					"nfs_access": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Component that has reserved the subnet.",
																					},
																					"nfs_all_squash": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
																					},
																					"nfs_root_squash": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether clients from this subnet can mount as root on NFS.",
																					},
																					"s3_access": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																					},
																					"smb_access": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																					},
																					"tenant_id": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the unique id of the tenant.",
																					},
																				},
																			},
																		},
																		"throttling_policy": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the throttling policy for a registered Protection Source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"enforce_max_streams": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																					},
																					"enforce_registered_source_max_backups": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																					},
																					"is_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																					},
																					"latency_thresholds": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"active_task_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																								},
																								"new_task_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																								},
																							},
																						},
																					},
																					"max_concurrent_streams": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																					},
																					"nas_source_params": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																								},
																								"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																								},
																								"max_parallel_read_write_full_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																								},
																								"max_parallel_read_write_incremental_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																								},
																							},
																						},
																					},
																					"registered_source_max_concurrent_backups": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																					},
																					"storage_array_snapshot_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_max_snapshots_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																								},
																								"is_max_space_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																								},
																								"storage_array_snapshot_max_space_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Space Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshot_space_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
																											},
																										},
																									},
																								},
																								"storage_array_snapshot_throttling_policies": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies throttling policies configured for individual volume/lun.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the volume id of the storage array snapshot config.",
																											},
																											"is_max_snapshots_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																											},
																											"is_max_space_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																											},
																											"max_snapshot_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshots": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
																														},
																													},
																												},
																											},
																											"max_space_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Space Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshot_space_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
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
																		"throttling_policy_overrides": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies throttling policy override for a Datastore in a registered entity.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"datastore_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Protection Source id of the Datastore.",
																					},
																					"datastore_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the display name of the Datastore.",
																					},
																					"throttling_policy": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the throttling policy for a registered Protection Source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"enforce_max_streams": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																								},
																								"enforce_registered_source_max_backups": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																								},
																								"is_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																								},
																								"latency_thresholds": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"active_task_msecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																											},
																											"new_task_msecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																											},
																										},
																									},
																								},
																								"max_concurrent_streams": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																								},
																								"nas_source_params": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																											},
																											"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																											},
																											"max_parallel_read_write_full_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																											},
																											"max_parallel_read_write_incremental_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																											},
																										},
																									},
																								},
																								"registered_source_max_concurrent_backups": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																								},
																								"storage_array_snapshot_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Configuration.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"is_max_snapshots_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																											},
																											"is_max_space_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																											},
																											"storage_array_snapshot_max_space_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Space Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshot_space_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
																														},
																													},
																												},
																											},
																											"storage_array_snapshot_throttling_policies": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies throttling policies configured for individual volume/lun.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Specifies the volume id of the storage array snapshot config.",
																														},
																														"is_max_snapshots_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																														},
																														"is_max_space_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																														},
																														"max_snapshot_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshots": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
																																	},
																																},
																															},
																														},
																														"max_space_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Space Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshot_space_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
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
																			},
																		},
																		"use_o_auth_for_exchange_online": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
																		},
																		"use_vm_bios_uuid": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
																		},
																		"user_messages": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies username to access the target source.",
																		},
																		"vlan_params": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the VLAN configuration for Recovery.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"vlan": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																					},
																					"disable_vlan": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
																					},
																					"interface_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																					},
																				},
																			},
																		},
																		"warning_messages": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"total_downtiered_size_in_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total bytes downtiered from the source so far.",
															},
															"total_uptiered_size_in_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total bytes uptiered to the source so far.",
															},
															"unprotected_sources_summary": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Aggregated information about a node subtree.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"environment": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
																		},
																		"leaves_count": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the number of leaf nodes under the subtree of this node.",
																		},
																		"total_logical_size": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the total logical size of the data under the subtree of this node.",
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
									"entity_pagination_parameters": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the cursor based pagination parameters for Protection Source and its children. Pagination is supported at a given level within the Protection Source Hierarchy with the help of before or after cursors. A Cursor will always refer to a specific source within the source dataset but will be invalidated if the item is removed.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"after_cursor_entity_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the entity id starting from which the items are to be returned.",
												},
												"before_cursor_entity_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the entity id upto which the items are to be returned.",
												},
												"node_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the entity id for the Node at any level within the Source entity hierarchy whose children are to be paginated.",
												},
												"page_size": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the maximum number of entities to be returned within the page.",
												},
											},
										},
									},
									"entity_permission_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the permission information of entities.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"entity_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the entity id.",
												},
												"groups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies struct with basic group details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"domain": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies domain name of the user.",
															},
															"group_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies group name of the group.",
															},
															"sid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies unique Security ID (SID) of the user.",
															},
															"tenant_ids": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the tenants to which the group belongs to.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
												"is_inferred": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether the Entity Permission Information is inferred or not. For example, SQL application hosted over vCenter will have inferred entity permission information.",
												},
												"is_registered_by_sp": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether this entity is registered by the SP or not. This will be populated only if the entity is a root entity. Refer to magneto/base/permissions.proto for details.",
												},
												"registering_tenant_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the tenant id that registered this entity. This will be populated only if the entity is a root entity.",
												},
												"tenant": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies struct with basic tenant details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"bifrost_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if this tenant is bifrost enabled or not.",
															},
															"is_managed_on_helios": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether this tenant is manged on helios.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies name of the tenant.",
															},
															"tenant_id": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the unique id of the tenant.",
															},
														},
													},
												},
												"users": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies struct with basic user details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"domain": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies domain name of the user.",
															},
															"sid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies unique Security ID (SID) of the user.",
															},
															"tenant_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the tenant to which the user belongs to.",
															},
															"user_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies user name of the user.",
															},
														},
													},
												},
											},
										},
									},
									"logical_size": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the logical size of the data in bytes for the Object on this node. Presence of this field indicates this node is a leaf node.",
									},
									"nodes": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies children of the current node in the Protection Sources hierarchy. When representing Objects in memory, the entire Object subtree hierarchy is represented. You can use this subtree to navigate down the Object hierarchy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{},
										},
									},
									"object_protection_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the Object Protection Info of the Protection Source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"auto_protect_parent_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the auto protect parent id if this entity is protected based on auto protection. This is only specified for leaf entities.",
												},
												"entity_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the entity id.",
												},
												"has_active_object_protection_spec": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies if the entity is under object protection.",
												},
											},
										},
									},
									"protected_sources_summary": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Array of Protected Objects. Specifies aggregated information about all the child Objects of this node that are currently protected by a Protection Job. There is one entry for each environment that is being backed up. The aggregated information for the Object hierarchy's environment will be available at the 0th index of the vector.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
												},
												"leaves_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of leaf nodes under the subtree of this node.",
												},
												"total_logical_size": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the total logical size of the data under the subtree of this node.",
												},
											},
										},
									},
									"protection_source": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies details about an Acropolis Protection Source when the environment is set to 'kAcropolis'.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"connection_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the connection id of the tenant.",
												},
												"connector_group_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the connector group id of the connector groups.",
												},
												"custom_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the user provided custom name of the Protection Source.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment (such as 'kVMware' or 'kSQL') where the Protection Source exists. Depending on the environment, one of the following Protection Sources are initialized.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies an id of the Protection Source.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies a name of the Protection Source.",
												},
												"parent_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies an id of the parent of the Protection Source.",
												},
												"physical_protection_source": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a Protection Source in a Physical environment.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"agents": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifiles the agents running on the Physical Protection Source and the status information.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cbmr_version": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the version if Cristie BMR product is installed on the host.",
																		},
																		"file_cbt_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "CBT version and service state info.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"file_version": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"build_ver": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																								"major_ver": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																								"minor_ver": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																								"revision_num": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"is_installed": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Indicates whether the cbt driver is installed.",
																					},
																					"reboot_status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Indicates whether host is rebooted post VolCBT installation.",
																					},
																					"service_state": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Structure to Hold Service Status.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"name": &schema.Schema{
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																								"state": &schema.Schema{
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"host_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the host type where the agent is running. This is only set for persistent agents.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the agent's id.",
																		},
																		"name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the agent's name.",
																		},
																		"oracle_multi_node_channel_supported": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether oracle multi node multi channel is supported or not.",
																		},
																		"registration_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies information about a registered Source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"access_info": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the parameters required to establish a connection with a particular environment.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"connection_id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
																								},
																								"connector_group_id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
																								},
																								"endpoint": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
																								},
																								"environment": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																								},
																								"id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
																								},
																								"version": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
																								},
																							},
																						},
																					},
																					"allowed_ip_addresses": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"authentication_error_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
																					},
																					"authentication_status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
																					},
																					"blacklisted_ip_addresses": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "This field is deprecated. Use DeniedIpAddresses instead.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"denied_ip_addresses": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"environments": &schema.Schema{
																						Type:        schema.TypeList,
																						Optional:    true,
																						Description: "Return only Protection Sources that match the passed in environment type such as 'kVMware', 'kSQL', 'kView' 'kPhysical', 'kPuppeteer', 'kPure', 'kNetapp', 'kGenericNas', 'kHyperV', 'kAcropolis', or 'kAzure'. For example, set this parameter to 'kVMware' to only return the Sources (and their Object subtrees) found in the 'kVMware' (VMware vCenter Server) environment.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"is_db_authenticated": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if application entity dbAuthenticated or not.",
																					},
																					"is_storage_array_snapshot_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if this source entity has enabled storage array snapshot or not.",
																					},
																					"link_vms_across_vcenter": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
																					},
																					"minimum_free_space_gb": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																					},
																					"minimum_free_space_percent": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																					},
																					"password": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies password of the username to access the target source.",
																					},
																					"physical_params": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"applications": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																									Elem: &schema.Schema{
																										Type: schema.TypeString,
																									},
																								},
																								"password": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies password of the username to access the target source.",
																								},
																								"throttling_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the source side throttling configuration.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"cpu_throttling_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Throttling Configuration Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"fixed_threshold": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																														},
																														"pattern_type": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																														},
																														"throttling_windows": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Throttling Window Parameters Definition.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"day_time_window": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Window Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"end_time": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Day Time Parameters.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"day": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
																																								Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																							},
																																							"time": &schema.Schema{
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
																																				"start_time": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Day Time Parameters.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"day": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
																																								Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																							},
																																							"time": &schema.Schema{
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
																																			},
																																		},
																																	},
																																	"threshold": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "Throttling threshold applicable in the window.",
																																	},
																																},
																															},
																														},
																													},
																												},
																											},
																											"network_throttling_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Throttling Configuration Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"fixed_threshold": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																														},
																														"pattern_type": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																														},
																														"throttling_windows": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Throttling Window Parameters Definition.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"day_time_window": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Window Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"end_time": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Day Time Parameters.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"day": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
																																								Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																							},
																																							"time": &schema.Schema{
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
																																				"start_time": &schema.Schema{
																																					Type:        schema.TypeList,
																																					Computed:    true,
																																					Description: "Specifies the Day Time Parameters.",
																																					Elem: &schema.Resource{
																																						Schema: map[string]*schema.Schema{
																																							"day": &schema.Schema{
																																								Type:        schema.TypeString,
																																								Computed:    true,
																																								Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																							},
																																							"time": &schema.Schema{
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
																																			},
																																		},
																																	},
																																	"threshold": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "Throttling threshold applicable in the window.",
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
																								"username": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies username to access the target source.",
																								},
																							},
																						},
																					},
																					"progress_monitor_path": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
																					},
																					"refresh_error_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
																					},
																					"refresh_time_usecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
																					},
																					"registered_apps_info": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies information of the applications registered on this protection source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"authentication_error_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
																								},
																								"authentication_status": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
																								},
																								"environment": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																								},
																								"host_settings_check_results": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"check_type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																											},
																											"result_type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																											},
																											"user_message": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies a descriptive message for failed/warning types.",
																											},
																										},
																									},
																								},
																								"refresh_error_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
																								},
																							},
																						},
																					},
																					"registration_time_usecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
																					},
																					"subnets": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"component": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Component that has reserved the subnet.",
																								},
																								"description": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Description of the subnet.",
																								},
																								"id": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "ID of the subnet.",
																								},
																								"ip": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies either an IPv6 address or an IPv4 address.",
																								},
																								"netmask_bits": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "netmaskBits.",
																								},
																								"netmask_ip4": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
																								},
																								"nfs_access": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Component that has reserved the subnet.",
																								},
																								"nfs_all_squash": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
																								},
																								"nfs_root_squash": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether clients from this subnet can mount as root on NFS.",
																								},
																								"s3_access": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																								},
																								"smb_access": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																								},
																								"tenant_id": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the unique id of the tenant.",
																								},
																							},
																						},
																					},
																					"throttling_policy": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the throttling policy for a registered Protection Source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"enforce_max_streams": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																								},
																								"enforce_registered_source_max_backups": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																								},
																								"is_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																								},
																								"latency_thresholds": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"active_task_msecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																											},
																											"new_task_msecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																											},
																										},
																									},
																								},
																								"max_concurrent_streams": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																								},
																								"nas_source_params": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																											},
																											"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																											},
																											"max_parallel_read_write_full_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																											},
																											"max_parallel_read_write_incremental_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																											},
																										},
																									},
																								},
																								"registered_source_max_concurrent_backups": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																								},
																								"storage_array_snapshot_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Configuration.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"is_max_snapshots_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																											},
																											"is_max_space_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																											},
																											"storage_array_snapshot_max_space_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Space Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshot_space_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
																														},
																													},
																												},
																											},
																											"storage_array_snapshot_throttling_policies": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies throttling policies configured for individual volume/lun.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Specifies the volume id of the storage array snapshot config.",
																														},
																														"is_max_snapshots_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																														},
																														"is_max_space_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																														},
																														"max_snapshot_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshots": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
																																	},
																																},
																															},
																														},
																														"max_space_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Space Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshot_space_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
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
																					"throttling_policy_overrides": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies throttling policy override for a Datastore in a registered entity.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"datastore_id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the Protection Source id of the Datastore.",
																								},
																								"datastore_name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the display name of the Datastore.",
																								},
																								"throttling_policy": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the throttling policy for a registered Protection Source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"enforce_max_streams": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																											},
																											"enforce_registered_source_max_backups": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																											},
																											"is_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																											},
																											"latency_thresholds": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"active_task_msecs": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																														},
																														"new_task_msecs": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																														},
																													},
																												},
																											},
																											"max_concurrent_streams": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																											},
																											"nas_source_params": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																														},
																														"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																														},
																														"max_parallel_read_write_full_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																														},
																														"max_parallel_read_write_incremental_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																														},
																													},
																												},
																											},
																											"registered_source_max_concurrent_backups": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																											},
																											"storage_array_snapshot_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Configuration.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"is_max_snapshots_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																														},
																														"is_max_space_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																														},
																														"storage_array_snapshot_max_space_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Space Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshot_space_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
																																	},
																																},
																															},
																														},
																														"storage_array_snapshot_throttling_policies": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies throttling policies configured for individual volume/lun.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"id": &schema.Schema{
																																		Type:        schema.TypeInt,
																																		Computed:    true,
																																		Description: "Specifies the volume id of the storage array snapshot config.",
																																	},
																																	"is_max_snapshots_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																																	},
																																	"is_max_space_config_enabled": &schema.Schema{
																																		Type:        schema.TypeBool,
																																		Computed:    true,
																																		Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																																	},
																																	"max_snapshot_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_snapshots": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Max number of storage snapshots allowed per volume/lun.",
																																				},
																																			},
																																		},
																																	},
																																	"max_space_config": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies Storage Array Snapshot Max Space Config.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"max_snapshot_space_percentage": &schema.Schema{
																																					Type:        schema.TypeFloat,
																																					Computed:    true,
																																					Description: "Max number of storage snapshots allowed per volume/lun.",
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
																						},
																					},
																					"use_o_auth_for_exchange_online": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
																					},
																					"use_vm_bios_uuid": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
																					},
																					"user_messages": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"username": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies username to access the target source.",
																					},
																					"vlan_params": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the VLAN configuration for Recovery.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"vlan": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																								},
																								"disable_vlan": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
																								},
																								"interface_name": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																								},
																							},
																						},
																					},
																					"warning_messages": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																				},
																			},
																		},
																		"source_side_dedup_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether source side dedup is enabled or not.",
																		},
																		"status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the agent status. Specifies the status of the agent running on a physical source.",
																		},
																		"status_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies additional details about the agent status.",
																		},
																		"upgradability": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the upgradability of the agent running on the physical server. Specifies the upgradability of the agent running on the physical server.",
																		},
																		"upgrade_status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the status of the upgrade of the agent on a physical server. Specifies the status of the upgrade of the agent on a physical server.",
																		},
																		"upgrade_status_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies detailed message about the agent upgrade failure. This field is not set for successful upgrade.",
																		},
																		"version": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the version of the Agent software.",
																		},
																		"vol_cbt_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "CBT version and service state info.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"file_version": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"build_ver": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																								"major_ver": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																								"minor_ver": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																								"revision_num": &schema.Schema{
																									Type:     schema.TypeFloat,
																									Computed: true,
																								},
																							},
																						},
																					},
																					"is_installed": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Indicates whether the cbt driver is installed.",
																					},
																					"reboot_status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Indicates whether host is rebooted post VolCBT installation.",
																					},
																					"service_state": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Structure to Hold Service Status.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"name": &schema.Schema{
																									Type:     schema.TypeString,
																									Computed: true,
																								},
																								"state": &schema.Schema{
																									Type:     schema.TypeString,
																									Computed: true,
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
															"cluster_source_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of cluster resource this source represents.",
															},
															"host_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the hostname.",
															},
															"host_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the environment type for the host.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies an id for an object that is unique across Cohesity Clusters. The id is composite of all the ids listed below.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cluster_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Cohesity Cluster id where the object was created.",
																		},
																		"cluster_incarnation_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies an id for the Cohesity Cluster that is generated when a Cohesity Cluster is initially created.",
																		},
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a unique id assigned to an object (such as a Job) by the Cohesity Cluster.",
																		},
																	},
																},
															},
															"is_proxy_host": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the physical host is a proxy host.",
															},
															"memory_size_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the total memory on the host in bytes.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a human readable name of the Protection Source.",
															},
															"networking_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the struct containing information about network addresses configured on the given box. This is needed for dealing with Windows/Oracle Cluster resources that we discover and protect automatically.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"resource_vec": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The list of resources on the system that are accessible by an IP address.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"endpoints": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "The endpoints by which the resource is accessible.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"fqdn": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The Fully Qualified Domain Name.",
																								},
																								"ipv4_addr": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The IPv4 address.",
																								},
																								"ipv6_addr": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "The IPv6 address.",
																								},
																							},
																						},
																					},
																					"type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The type of the resource.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"num_processors": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the number of processors on the host.",
															},
															"os_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a human readable name of the OS of the Protection Source.",
															},
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of managed Object in a Physical Protection Source. 'kGroup' indicates the EH container.",
															},
															"vcs_version": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies cluster version for VCS host.",
															},
															"volumes": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Array of Physical Volumes. Specifies the volumes available on the physical host. These fields are populated only for the kPhysicalHost type.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"device_path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the path to the device that hosts the volume locally.",
																		},
																		"guid": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies an id for the Physical Volume.",
																		},
																		"is_boot_volume": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the volume is boot volume.",
																		},
																		"is_extended_attributes_supported": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether this volume supports extended attributes (like ACLs) when performing file backups.",
																		},
																		"is_protected": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if a volume is protected by a Job.",
																		},
																		"is_shared_volume": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether the volume is shared volume.",
																		},
																		"label": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a volume label that can be used for displaying additional identifying information about a volume.",
																		},
																		"logical_size_bytes": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the logical size of the volume in bytes that is not reduced by change-block tracking, compression and deduplication.",
																		},
																		"mount_points": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the mount points where the volume is mounted, for example- 'C:', '/mnt/foo' etc.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"mount_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies mount type of volume e.g. nfs, autofs, ext4 etc.",
																		},
																		"network_path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the full path to connect to the network attached volume. For example, (IP or hostname):/path/to/share for NFS volumes).",
																		},
																		"used_size_bytes": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the size used by the volume in bytes.",
																		},
																	},
																},
															},
															"vsswriters": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies vss writer information about a Physical Protection Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_writer_excluded": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the writer will be excluded by default.",
																		},
																		"writer_name": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies the name of the writer.",
																		},
																	},
																},
															},
														},
													},
												},
												"sql_protection_source": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies an Object representing one SQL Server instance or database.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_available_for_vss_backup": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the database is marked as available for backup according to the SQL Server VSS writer. This may be false if either the state of the databases is not online, or if the VSS writer is not online. This field is set only for type 'kDatabase'.",
															},
															"created_timestamp": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the time when the database was created. It is displayed in the timezone of the SQL server on which this database is running.",
															},
															"database_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the database name of the SQL Protection Source, if the type is database.",
															},
															"db_aag_entity_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the AAG entity id if the database is part of an AAG. This field is set only for type 'kDatabase'.",
															},
															"db_aag_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the AAG if the database is part of an AAG. This field is set only for type 'kDatabase'.",
															},
															"db_compatibility_level": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the versions of SQL server that the database is compatible with.",
															},
															"db_file_groups": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the information about the set of file groups for this db on the host. This is only set if the type is kDatabase.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"db_files": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the last known information about the set of database files on the host. This field is set only for type 'kDatabase'.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"file_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the format type of the file that SQL database stores the data. Specifies the format type of the file that SQL database stores the data. 'kRows' refers to a data file 'kLog' refers to a log file 'kFileStream' refers to a directory containing FILESTREAM data 'kNotSupportedType' is for information purposes only. Not supported. 'kFullText' refers to a full-text catalog.",
																		},
																		"full_path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the full path of the database file on the SQL host machine.",
																		},
																		"size_bytes": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the last known size of the database file.",
																		},
																	},
																},
															},
															"db_owner_username": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the name of the database owner.",
															},
															"default_database_location": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the default path for data files for DBs in an instance.",
															},
															"default_log_location": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the default path for log files for DBs in an instance.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a unique id for a SQL Protection Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"created_date_msecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a unique identifier generated from the date the database is created or renamed. Cohesity uses this identifier in combination with the databaseId to uniquely identify a database.",
																		},
																		"database_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a unique id of the database but only for the life of the database. SQL Server may reuse database ids. Cohesity uses the createDateMsecs in combination with this databaseId to uniquely identify a database.",
																		},
																		"instance_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies unique id for the SQL Server instance. This id does not change during the life of the instance.",
																		},
																	},
																},
															},
															"is_encrypted": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the database is TDE enabled.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the instance name of the SQL Protection Source.",
															},
															"owner_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the id of the container VM for the SQL Protection Source.",
															},
															"recovery_model": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Recovery Model for the database in SQL environment. Only meaningful for the 'kDatabase' SQL Protection Source. Specifies the Recovery Model set for the Microsoft SQL Server. 'kSimpleRecoveryModel' indicates the Simple SQL Recovery Model which does not utilize log backups. 'kFullRecoveryModel' indicates the Full SQL Recovery Model which requires log backups and allows recovery to a single point in time. 'kBulkLoggedRecoveryModel' indicates the Bulk Logged SQL Recovery Model which requires log backups and allows high-performance bulk copy operations.",
															},
															"sql_server_db_state": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The state of the database as returned by SQL Server. Indicates the state of the database. The values correspond to the 'state' field in the system table sys.databases. See https://goo.gl/P66XqM. 'kOnline' indicates that database is in online state. 'kRestoring' indicates that database is in restore state. 'kRecovering' indicates that database is in recovery state. 'kRecoveryPending' indicates that database recovery is in pending state. 'kSuspect' indicates that primary filegroup is suspect and may be damaged. 'kEmergency' indicates that manually forced emergency state. 'kOffline' indicates that database is in offline state. 'kCopying' indicates that database is in copying state. 'kOfflineSecondary' indicates that secondary database is in offline state.",
															},
															"sql_server_instance_version": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the Server Instance Version.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"build": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the build.",
																		},
																		"major_version": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the major version.",
																		},
																		"minor_version": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the minor version.",
																		},
																		"revision": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the revision.",
																		},
																		"version_string": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the version string.",
																		},
																	},
																},
															},
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of the managed Object in a SQL Protection Source. Examples of SQL Objects include 'kInstance' and 'kDatabase'. 'kInstance' indicates that SQL server instance is being protected. 'kDatabase' indicates that SQL server database is being protected. 'kAAG' indicates that SQL AAG (AlwaysOn Availability Group) is being protected. 'kAAGRootContainer' indicates that SQL AAG's root container is being protected. 'kRootContainer' indicates root container for SQL sources.",
															},
														},
													},
												},
											},
										},
									},
									"registration_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information about a registered Source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"access_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the parameters required to establish a connection with a particular environment.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"connection_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
															},
															"connector_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
															},
															"endpoint": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
															},
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
															},
															"version": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
															},
														},
													},
												},
												"allowed_ip_addresses": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"authentication_error_message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
												},
												"authentication_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
												},
												"blacklisted_ip_addresses": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "This field is deprecated. Use DeniedIpAddresses instead.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"denied_ip_addresses": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"environments": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Return only Protection Sources that match the passed in environment type such as 'kVMware', 'kSQL', 'kView' 'kPhysical', 'kPuppeteer', 'kPure', 'kNetapp', 'kGenericNas', 'kHyperV', 'kAcropolis', or 'kAzure'. For example, set this parameter to 'kVMware' to only return the Sources (and their Object subtrees) found in the 'kVMware' (VMware vCenter Server) environment.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"is_db_authenticated": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if application entity dbAuthenticated or not.",
												},
												"is_storage_array_snapshot_enabled": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if this source entity has enabled storage array snapshot or not.",
												},
												"link_vms_across_vcenter": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
												},
												"minimum_free_space_gb": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
												},
												"minimum_free_space_percent": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
												},
												"password": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies password of the username to access the target source.",
												},
												"physical_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"applications": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"password": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies password of the username to access the target source.",
															},
															"throttling_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the source side throttling configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cpu_throttling_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the Throttling Configuration Parameters.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"fixed_threshold": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																					},
																					"pattern_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																					},
																					"throttling_windows": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Throttling Window Parameters Definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"day_time_window": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Day Time Window Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"end_time": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Day Time Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																														},
																														"time": &schema.Schema{
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
																											"start_time": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Day Time Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																														},
																														"time": &schema.Schema{
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
																										},
																									},
																								},
																								"threshold": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Throttling threshold applicable in the window.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"network_throttling_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the Throttling Configuration Parameters.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"fixed_threshold": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																					},
																					"pattern_type": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																					},
																					"throttling_windows": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Throttling Window Parameters Definition.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"day_time_window": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Day Time Window Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"end_time": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Day Time Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																														},
																														"time": &schema.Schema{
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
																											"start_time": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Day Time Parameters.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day": &schema.Schema{
																															Type:        schema.TypeString,
																															Computed:    true,
																															Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																														},
																														"time": &schema.Schema{
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
																										},
																									},
																								},
																								"threshold": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Throttling threshold applicable in the window.",
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
															"username": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies username to access the target source.",
															},
														},
													},
												},
												"progress_monitor_path": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
												},
												"refresh_error_message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
												},
												"refresh_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
												},
												"registered_apps_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies information of the applications registered on this protection source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"authentication_error_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
															},
															"authentication_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
															},
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
															},
															"host_settings_check_results": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"check_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																		},
																		"result_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																		},
																		"user_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a descriptive message for failed/warning types.",
																		},
																	},
																},
															},
															"refresh_error_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
															},
														},
													},
												},
												"registration_time_usecs": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
												},
												"subnets": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"component": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Component that has reserved the subnet.",
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the subnet.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "ID of the subnet.",
															},
															"ip": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies either an IPv6 address or an IPv4 address.",
															},
															"netmask_bits": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "netmaskBits.",
															},
															"netmask_ip4": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
															},
															"nfs_access": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Component that has reserved the subnet.",
															},
															"nfs_all_squash": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
															},
															"nfs_root_squash": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether clients from this subnet can mount as root on NFS.",
															},
															"s3_access": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
															},
															"smb_access": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
															},
															"tenant_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique id of the tenant.",
															},
														},
													},
												},
												"throttling_policy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the throttling policy for a registered Protection Source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enforce_max_streams": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
															},
															"enforce_registered_source_max_backups": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
															},
															"is_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
															},
															"latency_thresholds": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"active_task_msecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																		},
																		"new_task_msecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																		},
																	},
																},
															},
															"max_concurrent_streams": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
															},
															"nas_source_params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																		},
																		"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																		},
																		"max_parallel_read_write_full_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																		},
																		"max_parallel_read_write_incremental_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																		},
																	},
																},
															},
															"registered_source_max_concurrent_backups": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
															},
															"storage_array_snapshot_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies Storage Array Snapshot Configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_max_snapshots_config_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																		},
																		"is_max_space_config_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																		},
																		"storage_array_snapshot_max_space_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies Storage Array Snapshot Max Space Config.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_snapshot_space_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Max number of storage snapshots allowed per volume/lun.",
																					},
																				},
																			},
																		},
																		"storage_array_snapshot_throttling_policies": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies throttling policies configured for individual volume/lun.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the volume id of the storage array snapshot config.",
																					},
																					"is_max_snapshots_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																					},
																					"is_max_space_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																					},
																					"max_snapshot_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_snapshots": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Max number of storage snapshots allowed per volume/lun.",
																								},
																							},
																						},
																					},
																					"max_space_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Max Space Config.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_snapshot_space_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Max number of storage snapshots allowed per volume/lun.",
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
												"throttling_policy_overrides": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies throttling policy override for a Datastore in a registered entity.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"datastore_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the Protection Source id of the Datastore.",
															},
															"datastore_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the display name of the Datastore.",
															},
															"throttling_policy": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the throttling policy for a registered Protection Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"enforce_max_streams": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																		},
																		"enforce_registered_source_max_backups": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																		},
																		"is_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																		},
																		"latency_thresholds": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"active_task_msecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																					},
																					"new_task_msecs": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																					},
																				},
																			},
																		},
																		"max_concurrent_streams": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																		},
																		"nas_source_params": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																					},
																					"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																					},
																					"max_parallel_read_write_full_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																					},
																					"max_parallel_read_write_incremental_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																					},
																				},
																			},
																		},
																		"registered_source_max_concurrent_backups": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																		},
																		"storage_array_snapshot_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies Storage Array Snapshot Configuration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"is_max_snapshots_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																					},
																					"is_max_space_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																					},
																					"storage_array_snapshot_max_space_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Max Space Config.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_snapshot_space_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Max number of storage snapshots allowed per volume/lun.",
																								},
																							},
																						},
																					},
																					"storage_array_snapshot_throttling_policies": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies throttling policies configured for individual volume/lun.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"id": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the volume id of the storage array snapshot config.",
																								},
																								"is_max_snapshots_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																								},
																								"is_max_space_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																								},
																								"max_snapshot_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshots": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
																											},
																										},
																									},
																								},
																								"max_space_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Space Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshot_space_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
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
													},
												},
												"use_o_auth_for_exchange_online": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
												},
												"use_vm_bios_uuid": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
												},
												"user_messages": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"username": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies username to access the target source.",
												},
												"vlan_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the VLAN configuration for Recovery.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"vlan": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
															},
															"disable_vlan": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
															},
															"interface_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
															},
														},
													},
												},
												"warning_messages": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"total_downtiered_size_in_bytes": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the total bytes downtiered from the source so far.",
									},
									"total_uptiered_size_in_bytes": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the total bytes uptiered to the source so far.",
									},
									"unprotected_sources_summary": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Aggregated information about a node subtree.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
												},
												"leaves_count": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of leaf nodes under the subtree of this node.",
												},
												"total_logical_size": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the total logical size of the data under the subtree of this node.",
												},
											},
										},
									},
								},
							},
						},
						"object_protection_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the Object Protection Info of the Protection Source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"auto_protect_parent_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the auto protect parent id if this entity is protected based on auto protection. This is only specified for leaf entities.",
									},
									"entity_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the entity id.",
									},
									"has_active_object_protection_spec": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies if the entity is under object protection.",
									},
								},
							},
						},
						"protected_sources_summary": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of Protected Objects. Specifies aggregated information about all the child Objects of this node that are currently protected by a Protection Job. There is one entry for each environment that is being backed up. The aggregated information for the Object hierarchy's environment will be available at the 0th index of the vector.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
									},
									"leaves_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of leaf nodes under the subtree of this node.",
									},
									"total_logical_size": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the total logical size of the data under the subtree of this node.",
									},
								},
							},
						},
						"protection_source": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies details about an Acropolis Protection Source when the environment is set to 'kAcropolis'.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"connection_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the connection id of the tenant.",
									},
									"connector_group_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the connector group id of the connector groups.",
									},
									"custom_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the user provided custom name of the Protection Source.",
									},
									"environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment (such as 'kVMware' or 'kSQL') where the Protection Source exists. Depending on the environment, one of the following Protection Sources are initialized.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies an id of the Protection Source.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies a name of the Protection Source.",
									},
									"parent_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies an id of the parent of the Protection Source.",
									},
									"physical_protection_source": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a Protection Source in a Physical environment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"agents": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifiles the agents running on the Physical Protection Source and the status information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cbmr_version": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the version if Cristie BMR product is installed on the host.",
															},
															"file_cbt_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "CBT version and service state info.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"file_version": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"build_ver": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																					"major_ver": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																					"minor_ver": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																					"revision_num": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"is_installed": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Indicates whether the cbt driver is installed.",
																		},
																		"reboot_status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Indicates whether host is rebooted post VolCBT installation.",
																		},
																		"service_state": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Structure to Hold Service Status.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"state": &schema.Schema{
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																				},
																			},
																		},
																	},
																},
															},
															"host_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the host type where the agent is running. This is only set for persistent agents.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the agent's id.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the agent's name.",
															},
															"oracle_multi_node_channel_supported": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether oracle multi node multi channel is supported or not.",
															},
															"registration_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies information about a registered Source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"access_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the parameters required to establish a connection with a particular environment.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"connection_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
																					},
																					"connector_group_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
																					},
																					"endpoint": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
																					},
																					"environment": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
																					},
																					"version": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
																					},
																				},
																			},
																		},
																		"allowed_ip_addresses": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"authentication_error_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
																		},
																		"authentication_status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
																		},
																		"blacklisted_ip_addresses": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "This field is deprecated. Use DeniedIpAddresses instead.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"denied_ip_addresses": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"environments": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Return only Protection Sources that match the passed in environment type such as 'kVMware', 'kSQL', 'kView' 'kPhysical', 'kPuppeteer', 'kPure', 'kNetapp', 'kGenericNas', 'kHyperV', 'kAcropolis', or 'kAzure'. For example, set this parameter to 'kVMware' to only return the Sources (and their Object subtrees) found in the 'kVMware' (VMware vCenter Server) environment.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"is_db_authenticated": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if application entity dbAuthenticated or not.",
																		},
																		"is_storage_array_snapshot_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if this source entity has enabled storage array snapshot or not.",
																		},
																		"link_vms_across_vcenter": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
																		},
																		"minimum_free_space_gb": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																		},
																		"minimum_free_space_percent": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
																		},
																		"password": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies password of the username to access the target source.",
																		},
																		"physical_params": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"applications": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"password": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies password of the username to access the target source.",
																					},
																					"throttling_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the source side throttling configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"cpu_throttling_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Throttling Configuration Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"fixed_threshold": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																											},
																											"pattern_type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																											},
																											"throttling_windows": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Throttling Window Parameters Definition.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day_time_window": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Window Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"end_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																	"start_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																},
																															},
																														},
																														"threshold": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Throttling threshold applicable in the window.",
																														},
																													},
																												},
																											},
																										},
																									},
																								},
																								"network_throttling_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Throttling Configuration Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"fixed_threshold": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																											},
																											"pattern_type": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																											},
																											"throttling_windows": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies the Throttling Window Parameters Definition.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"day_time_window": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies the Day Time Window Parameters.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"end_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																	"start_time": &schema.Schema{
																																		Type:        schema.TypeList,
																																		Computed:    true,
																																		Description: "Specifies the Day Time Parameters.",
																																		Elem: &schema.Resource{
																																			Schema: map[string]*schema.Schema{
																																				"day": &schema.Schema{
																																					Type:        schema.TypeString,
																																					Computed:    true,
																																					Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																																				},
																																				"time": &schema.Schema{
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
																																},
																															},
																														},
																														"threshold": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Throttling threshold applicable in the window.",
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
																					"username": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies username to access the target source.",
																					},
																				},
																			},
																		},
																		"progress_monitor_path": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
																		},
																		"refresh_error_message": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
																		},
																		"refresh_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
																		},
																		"registered_apps_info": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies information of the applications registered on this protection source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"authentication_error_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
																					},
																					"authentication_status": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
																					},
																					"environment": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
																					},
																					"host_settings_check_results": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"check_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
																								},
																								"result_type": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
																								},
																								"user_message": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies a descriptive message for failed/warning types.",
																								},
																							},
																						},
																					},
																					"refresh_error_message": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
																					},
																				},
																			},
																		},
																		"registration_time_usecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
																		},
																		"subnets": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"component": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Component that has reserved the subnet.",
																					},
																					"description": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Description of the subnet.",
																					},
																					"id": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "ID of the subnet.",
																					},
																					"ip": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies either an IPv6 address or an IPv4 address.",
																					},
																					"netmask_bits": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "netmaskBits.",
																					},
																					"netmask_ip4": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
																					},
																					"nfs_access": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Component that has reserved the subnet.",
																					},
																					"nfs_all_squash": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
																					},
																					"nfs_root_squash": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether clients from this subnet can mount as root on NFS.",
																					},
																					"s3_access": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																					},
																					"smb_access": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
																					},
																					"tenant_id": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the unique id of the tenant.",
																					},
																				},
																			},
																		},
																		"throttling_policy": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the throttling policy for a registered Protection Source.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"enforce_max_streams": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																					},
																					"enforce_registered_source_max_backups": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																					},
																					"is_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																					},
																					"latency_thresholds": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"active_task_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																								},
																								"new_task_msecs": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																								},
																							},
																						},
																					},
																					"max_concurrent_streams": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																					},
																					"nas_source_params": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																								},
																								"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																								},
																								"max_parallel_read_write_full_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																								},
																								"max_parallel_read_write_incremental_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																								},
																							},
																						},
																					},
																					"registered_source_max_concurrent_backups": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																					},
																					"storage_array_snapshot_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Configuration.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"is_max_snapshots_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																								},
																								"is_max_space_config_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																								},
																								"storage_array_snapshot_max_space_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Max Space Config.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_snapshot_space_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Max number of storage snapshots allowed per volume/lun.",
																											},
																										},
																									},
																								},
																								"storage_array_snapshot_throttling_policies": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies throttling policies configured for individual volume/lun.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"id": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "Specifies the volume id of the storage array snapshot config.",
																											},
																											"is_max_snapshots_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																											},
																											"is_max_space_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																											},
																											"max_snapshot_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshots": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
																														},
																													},
																												},
																											},
																											"max_space_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Space Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshot_space_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
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
																		"throttling_policy_overrides": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies throttling policy override for a Datastore in a registered entity.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"datastore_id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the Protection Source id of the Datastore.",
																					},
																					"datastore_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the display name of the Datastore.",
																					},
																					"throttling_policy": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the throttling policy for a registered Protection Source.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"enforce_max_streams": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
																								},
																								"enforce_registered_source_max_backups": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
																								},
																								"is_enabled": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
																								},
																								"latency_thresholds": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"active_task_msecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																											},
																											"new_task_msecs": &schema.Schema{
																												Type:        schema.TypeInt,
																												Computed:    true,
																												Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																											},
																										},
																									},
																								},
																								"max_concurrent_streams": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
																								},
																								"nas_source_params": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																											},
																											"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																											},
																											"max_parallel_read_write_full_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																											},
																											"max_parallel_read_write_incremental_percentage": &schema.Schema{
																												Type:        schema.TypeFloat,
																												Computed:    true,
																												Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																											},
																										},
																									},
																								},
																								"registered_source_max_concurrent_backups": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
																								},
																								"storage_array_snapshot_config": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies Storage Array Snapshot Configuration.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"is_max_snapshots_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																											},
																											"is_max_space_config_enabled": &schema.Schema{
																												Type:        schema.TypeBool,
																												Computed:    true,
																												Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																											},
																											"storage_array_snapshot_max_space_config": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies Storage Array Snapshot Max Space Config.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"max_snapshot_space_percentage": &schema.Schema{
																															Type:        schema.TypeFloat,
																															Computed:    true,
																															Description: "Max number of storage snapshots allowed per volume/lun.",
																														},
																													},
																												},
																											},
																											"storage_array_snapshot_throttling_policies": &schema.Schema{
																												Type:        schema.TypeList,
																												Computed:    true,
																												Description: "Specifies throttling policies configured for individual volume/lun.",
																												Elem: &schema.Resource{
																													Schema: map[string]*schema.Schema{
																														"id": &schema.Schema{
																															Type:        schema.TypeInt,
																															Computed:    true,
																															Description: "Specifies the volume id of the storage array snapshot config.",
																														},
																														"is_max_snapshots_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																														},
																														"is_max_space_config_enabled": &schema.Schema{
																															Type:        schema.TypeBool,
																															Computed:    true,
																															Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																														},
																														"max_snapshot_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshots": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
																																	},
																																},
																															},
																														},
																														"max_space_config": &schema.Schema{
																															Type:        schema.TypeList,
																															Computed:    true,
																															Description: "Specifies Storage Array Snapshot Max Space Config.",
																															Elem: &schema.Resource{
																																Schema: map[string]*schema.Schema{
																																	"max_snapshot_space_percentage": &schema.Schema{
																																		Type:        schema.TypeFloat,
																																		Computed:    true,
																																		Description: "Max number of storage snapshots allowed per volume/lun.",
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
																			},
																		},
																		"use_o_auth_for_exchange_online": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
																		},
																		"use_vm_bios_uuid": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
																		},
																		"user_messages": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"username": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies username to access the target source.",
																		},
																		"vlan_params": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the VLAN configuration for Recovery.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"vlan": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																					},
																					"disable_vlan": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
																					},
																					"interface_name": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
																					},
																				},
																			},
																		},
																		"warning_messages": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"source_side_dedup_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether source side dedup is enabled or not.",
															},
															"status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the agent status. Specifies the status of the agent running on a physical source.",
															},
															"status_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies additional details about the agent status.",
															},
															"upgradability": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the upgradability of the agent running on the physical server. Specifies the upgradability of the agent running on the physical server.",
															},
															"upgrade_status": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the status of the upgrade of the agent on a physical server. Specifies the status of the upgrade of the agent on a physical server.",
															},
															"upgrade_status_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies detailed message about the agent upgrade failure. This field is not set for successful upgrade.",
															},
															"version": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the version of the Agent software.",
															},
															"vol_cbt_info": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "CBT version and service state info.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"file_version": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Subcomponent version. The interpretation of the version is based on operating system.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"build_ver": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																					"major_ver": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																					"minor_ver": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																					"revision_num": &schema.Schema{
																						Type:     schema.TypeFloat,
																						Computed: true,
																					},
																				},
																			},
																		},
																		"is_installed": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Indicates whether the cbt driver is installed.",
																		},
																		"reboot_status": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Indicates whether host is rebooted post VolCBT installation.",
																		},
																		"service_state": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Structure to Hold Service Status.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"name": &schema.Schema{
																						Type:     schema.TypeString,
																						Computed: true,
																					},
																					"state": &schema.Schema{
																						Type:     schema.TypeString,
																						Computed: true,
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
												"cluster_source_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of cluster resource this source represents.",
												},
												"host_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the hostname.",
												},
												"host_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment type for the host.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies an id for an object that is unique across Cohesity Clusters. The id is composite of all the ids listed below.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the Cohesity Cluster id where the object was created.",
															},
															"cluster_incarnation_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies an id for the Cohesity Cluster that is generated when a Cohesity Cluster is initially created.",
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a unique id assigned to an object (such as a Job) by the Cohesity Cluster.",
															},
														},
													},
												},
												"is_proxy_host": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if the physical host is a proxy host.",
												},
												"memory_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the total memory on the host in bytes.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies a human readable name of the Protection Source.",
												},
												"networking_info": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the struct containing information about network addresses configured on the given box. This is needed for dealing with Windows/Oracle Cluster resources that we discover and protect automatically.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"resource_vec": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The list of resources on the system that are accessible by an IP address.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"endpoints": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "The endpoints by which the resource is accessible.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"fqdn": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The Fully Qualified Domain Name.",
																					},
																					"ipv4_addr": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The IPv4 address.",
																					},
																					"ipv6_addr": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "The IPv6 address.",
																					},
																				},
																			},
																		},
																		"type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The type of the resource.",
																		},
																	},
																},
															},
														},
													},
												},
												"num_processors": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the number of processors on the host.",
												},
												"os_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies a human readable name of the OS of the Protection Source.",
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of managed Object in a Physical Protection Source. 'kGroup' indicates the EH container.",
												},
												"vcs_version": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies cluster version for VCS host.",
												},
												"volumes": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Array of Physical Volumes. Specifies the volumes available on the physical host. These fields are populated only for the kPhysicalHost type.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"device_path": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the path to the device that hosts the volume locally.",
															},
															"guid": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies an id for the Physical Volume.",
															},
															"is_boot_volume": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the volume is boot volume.",
															},
															"is_extended_attributes_supported": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether this volume supports extended attributes (like ACLs) when performing file backups.",
															},
															"is_protected": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if a volume is protected by a Job.",
															},
															"is_shared_volume": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether the volume is shared volume.",
															},
															"label": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a volume label that can be used for displaying additional identifying information about a volume.",
															},
															"logical_size_bytes": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the logical size of the volume in bytes that is not reduced by change-block tracking, compression and deduplication.",
															},
															"mount_points": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the mount points where the volume is mounted, for example- 'C:', '/mnt/foo' etc.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"mount_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies mount type of volume e.g. nfs, autofs, ext4 etc.",
															},
															"network_path": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the full path to connect to the network attached volume. For example, (IP or hostname):/path/to/share for NFS volumes).",
															},
															"used_size_bytes": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the size used by the volume in bytes.",
															},
														},
													},
												},
												"vsswriters": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies vss writer information about a Physical Protection Source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_writer_excluded": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the writer will be excluded by default.",
															},
															"writer_name": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies the name of the writer.",
															},
														},
													},
												},
											},
										},
									},
									"sql_protection_source": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies an Object representing one SQL Server instance or database.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_available_for_vss_backup": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether the database is marked as available for backup according to the SQL Server VSS writer. This may be false if either the state of the databases is not online, or if the VSS writer is not online. This field is set only for type 'kDatabase'.",
												},
												"created_timestamp": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the time when the database was created. It is displayed in the timezone of the SQL server on which this database is running.",
												},
												"database_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the database name of the SQL Protection Source, if the type is database.",
												},
												"db_aag_entity_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the AAG entity id if the database is part of an AAG. This field is set only for type 'kDatabase'.",
												},
												"db_aag_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the AAG if the database is part of an AAG. This field is set only for type 'kDatabase'.",
												},
												"db_compatibility_level": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the versions of SQL server that the database is compatible with.",
												},
												"db_file_groups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the information about the set of file groups for this db on the host. This is only set if the type is kDatabase.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"db_files": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the last known information about the set of database files on the host. This field is set only for type 'kDatabase'.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"file_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the format type of the file that SQL database stores the data. Specifies the format type of the file that SQL database stores the data. 'kRows' refers to a data file 'kLog' refers to a log file 'kFileStream' refers to a directory containing FILESTREAM data 'kNotSupportedType' is for information purposes only. Not supported. 'kFullText' refers to a full-text catalog.",
															},
															"full_path": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the full path of the database file on the SQL host machine.",
															},
															"size_bytes": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the last known size of the database file.",
															},
														},
													},
												},
												"db_owner_username": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the name of the database owner.",
												},
												"default_database_location": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the default path for data files for DBs in an instance.",
												},
												"default_log_location": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the default path for log files for DBs in an instance.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a unique id for a SQL Protection Source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"created_date_msecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a unique identifier generated from the date the database is created or renamed. Cohesity uses this identifier in combination with the databaseId to uniquely identify a database.",
															},
															"database_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a unique id of the database but only for the life of the database. SQL Server may reuse database ids. Cohesity uses the createDateMsecs in combination with this databaseId to uniquely identify a database.",
															},
															"instance_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies unique id for the SQL Server instance. This id does not change during the life of the instance.",
															},
														},
													},
												},
												"is_encrypted": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether the database is TDE enabled.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the instance name of the SQL Protection Source.",
												},
												"owner_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the id of the container VM for the SQL Protection Source.",
												},
												"recovery_model": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the Recovery Model for the database in SQL environment. Only meaningful for the 'kDatabase' SQL Protection Source. Specifies the Recovery Model set for the Microsoft SQL Server. 'kSimpleRecoveryModel' indicates the Simple SQL Recovery Model which does not utilize log backups. 'kFullRecoveryModel' indicates the Full SQL Recovery Model which requires log backups and allows recovery to a single point in time. 'kBulkLoggedRecoveryModel' indicates the Bulk Logged SQL Recovery Model which requires log backups and allows high-performance bulk copy operations.",
												},
												"sql_server_db_state": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The state of the database as returned by SQL Server. Indicates the state of the database. The values correspond to the 'state' field in the system table sys.databases. See https://goo.gl/P66XqM. 'kOnline' indicates that database is in online state. 'kRestoring' indicates that database is in restore state. 'kRecovering' indicates that database is in recovery state. 'kRecoveryPending' indicates that database recovery is in pending state. 'kSuspect' indicates that primary filegroup is suspect and may be damaged. 'kEmergency' indicates that manually forced emergency state. 'kOffline' indicates that database is in offline state. 'kCopying' indicates that database is in copying state. 'kOfflineSecondary' indicates that secondary database is in offline state.",
												},
												"sql_server_instance_version": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the Server Instance Version.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"build": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the build.",
															},
															"major_version": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the major version.",
															},
															"minor_version": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the minor version.",
															},
															"revision": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the revision.",
															},
															"version_string": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the version string.",
															},
														},
													},
												},
												"type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of the managed Object in a SQL Protection Source. Examples of SQL Objects include 'kInstance' and 'kDatabase'. 'kInstance' indicates that SQL server instance is being protected. 'kDatabase' indicates that SQL server database is being protected. 'kAAG' indicates that SQL AAG (AlwaysOn Availability Group) is being protected. 'kAAGRootContainer' indicates that SQL AAG's root container is being protected. 'kRootContainer' indicates root container for SQL sources.",
												},
											},
										},
									},
								},
							},
						},
						"registration_info": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies information about a registered Source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters required to establish a connection with a particular environment.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"connection_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "ID of the Bifrost (HyX or Rigel) network realm (i.e. a connection) associated with the source.",
												},
												"connector_group_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the Id of the connector group. Each connector group is collection of Rigel/hyx. Each entity will be tagged with connector group id.",
												},
												"endpoint": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specify an IP address or URL of the environment. (such as the IP address of the vCenter Server for a VMware environment).",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the environment like VMware, SQL, where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies a Unique id that is generated when the Source is registered. This is a convenience field that is used to maintain an index to different connection params.",
												},
												"version": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Version is updated each time the connector parameters are updated. This is used to discard older connector parameters.",
												},
											},
										},
									},
									"allowed_ip_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the list of IP Addresses on the registered source to be exclusively allowed for doing any type of IO operations.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"authentication_error_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies an authentication error message. This indicates the given credentials are rejected and the registration of the source is not successful.",
									},
									"authentication_status": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the status of the authenticating to the Protection Source when registering it with Cohesity Cluster.",
									},
									"blacklisted_ip_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "This field is deprecated. Use DeniedIpAddresses instead.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"denied_ip_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the list of IP Addresses on the registered source to be denied for doing any type of IO operations.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"environments": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Return only Protection Sources that match the passed in environment type such as 'kVMware', 'kSQL', 'kView' 'kPhysical', 'kPuppeteer', 'kPure', 'kNetapp', 'kGenericNas', 'kHyperV', 'kAcropolis', or 'kAzure'. For example, set this parameter to 'kVMware' to only return the Sources (and their Object subtrees) found in the 'kVMware' (VMware vCenter Server) environment.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"is_db_authenticated": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies if application entity dbAuthenticated or not.",
									},
									"is_storage_array_snapshot_enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies if this source entity has enabled storage array snapshot or not.",
									},
									"link_vms_across_vcenter": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies if the VM linking feature is enabled for this VCenter This means that VMs present in this VCenter which earlier belonged to some other VCenter(also registerd on same cluster) and were migrated, will be linked during EH refresh. This will enable preserving snapshot chains for migrated VMs.",
									},
									"minimum_free_space_gb": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the minimum free space in GiB of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in GiB) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
									},
									"minimum_free_space_percent": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the minimum free space in percentage of the space expected to be available on the datastore where the virtual disks of the VM being backed up. If the amount of free space(in percentage) is lower than the value given by this field, backup will be aborted. Note that this field is applicable only to 'kVMware' type of environments.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies password of the username to access the target source.",
									},
									"physical_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the parameters required to register Application Servers running in a Protection Source specific to a physical adapter.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"applications": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the types of applications such as 'kSQL', 'kExchange', 'kAD' running on the Protection Source. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"password": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies password of the username to access the target source.",
												},
												"throttling_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the source side throttling configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cpu_throttling_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the Throttling Configuration Parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"fixed_threshold": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																		},
																		"pattern_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																		},
																		"throttling_windows": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the Throttling Window Parameters Definition.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_time_window": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Day Time Window Parameters.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"end_time": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Day Time Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"day": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																											},
																											"time": &schema.Schema{
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
																								"start_time": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Day Time Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"day": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																											},
																											"time": &schema.Schema{
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
																							},
																						},
																					},
																					"threshold": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Throttling threshold applicable in the window.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"network_throttling_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the Throttling Configuration Parameters.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"fixed_threshold": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Fixed baseline threshold for throttling. This is mandatory for any other throttling type than kNoThrottling.",
																		},
																		"pattern_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the throttling pattern. 'kNoThrottling' indicates that throttling is not in force. 'kBaseThrottling' indicates indicates a constant base level throttling. 'kFixed' indicates a constant base level throttling.",
																		},
																		"throttling_windows": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the Throttling Window Parameters Definition.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_time_window": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies the Day Time Window Parameters.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"end_time": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Day Time Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"day": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																											},
																											"time": &schema.Schema{
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
																								"start_time": &schema.Schema{
																									Type:        schema.TypeList,
																									Computed:    true,
																									Description: "Specifies the Day Time Parameters.",
																									Elem: &schema.Resource{
																										Schema: map[string]*schema.Schema{
																											"day": &schema.Schema{
																												Type:        schema.TypeString,
																												Computed:    true,
																												Description: "Specifies the day of the week (such as 'kMonday') for scheduling throttling. Specifies a day in a week such as 'kSunday', 'kMonday', etc.",
																											},
																											"time": &schema.Schema{
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
																							},
																						},
																					},
																					"threshold": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Throttling threshold applicable in the window.",
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
												"username": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies username to access the target source.",
												},
											},
										},
									},
									"progress_monitor_path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Captures the current progress and pulse details w.r.t to either the registration or refresh.",
									},
									"refresh_error_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies a message if there was any error encountered during the last rebuild of the Protection Source tree. If there was no error during the last rebuild, this field is reset.",
									},
									"refresh_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source tree was most recently fetched and built.",
									},
									"registered_apps_info": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies information of the applications registered on this protection source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"authentication_error_message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "pecifies an authentication error message. This indicates the given credentials are rejected and the registration of the application is not successful.",
												},
												"authentication_status": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the status of authenticating to the Protection Source when registering this application with Cohesity Cluster. If the status is 'kFinished' and there is no error, registration is successful. Specifies the status of the authentication during the registration of a Protection Source. 'kPending' indicates the authentication is in progress. 'kScheduled' indicates the authentication is scheduled. 'kFinished' indicates the authentication is completed. 'kRefreshInProgress' indicates the refresh is in progress.",
												},
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the application environment. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc.",
												},
												"host_settings_check_results": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the list of check results internally performed to verify status of various services such as 'AgnetRunning', 'SQLWriterRunning' etc.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"check_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of the check internally performed. Specifies the type of the host check performed internally. 'kIsAgentPortAccessible' indicates the check for agent port access. 'kIsAgentRunning' indicates the status for the Cohesity agent service. 'kIsSQLWriterRunning' indicates the status for SQLWriter service. 'kAreSQLInstancesRunning' indicates the run status for all the SQL instances in the host. 'kCheckServiceLoginsConfig' checks the privileges and sysadmin status of the logins used by the SQL instance services, Cohesity agent service and the SQLWriter service. 'kCheckSQLFCIVIP' checks whether the SQL FCI is registered with a valid VIP or FQDN. 'kCheckSQLDiskSpace' checks whether volumes containing SQL DBs have at least 10% free space.",
															},
															"result_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of the result returned after performing the internal host check. Specifies the type of the host check result performed internally. 'kPass' indicates that the respective check was successful. 'kFail' indicates that the respective check failed as some mandatory setting is not met 'kWarning' indicates that the respective check has warning as certain non-mandatory setting is not met.",
															},
															"user_message": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies a descriptive message for failed/warning types.",
															},
														},
													},
												},
												"refresh_error_message": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies a message if there was any error encountered during the last rebuild of the application tree. If there was no error during the last rebuild, this field is reset.",
												},
											},
										},
									},
									"registration_time_usecs": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the Unix epoch time (in microseconds) when the Protection Source was registered.",
									},
									"subnets": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the list of subnets added during creation or updation of vmare source. Currently, this field will only be populated in case of VMware registration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"component": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Component that has reserved the subnet.",
												},
												"description": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Description of the subnet.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "ID of the subnet.",
												},
												"ip": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies either an IPv6 address or an IPv4 address.",
												},
												"netmask_bits": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "netmaskBits.",
												},
												"netmask_ip4": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the netmask using an IP4 address. The netmask can only be set using netmaskIp4 if the IP address is an IPv4 address.",
												},
												"nfs_access": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Component that has reserved the subnet.",
												},
												"nfs_all_squash": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether all clients from this subnet can map view with view_all_squash_uid/view_all_squash_gid configured in the view.",
												},
												"nfs_root_squash": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether clients from this subnet can mount as root on NFS.",
												},
												"s3_access": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies whether clients from this subnet can access using S3 protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
												},
												"smb_access": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies whether clients from this subnet can mount using SMB protocol. Protocol access level. 'kDisabled' indicates Protocol access level 'Disabled' 'kReadOnly' indicates Protocol access level 'ReadOnly' 'kReadWrite' indicates Protocol access level 'ReadWrite'.",
												},
												"tenant_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unique id of the tenant.",
												},
											},
										},
									},
									"throttling_policy": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the throttling policy for a registered Protection Source.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enforce_max_streams": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
												},
												"enforce_registered_source_max_backups": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
												},
												"is_enabled": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
												},
												"latency_thresholds": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"active_task_msecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
															},
															"new_task_msecs": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
															},
														},
													},
												},
												"max_concurrent_streams": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
												},
												"nas_source_params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
															},
															"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
															},
															"max_parallel_read_write_full_percentage": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
															},
															"max_parallel_read_write_incremental_percentage": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
															},
														},
													},
												},
												"registered_source_max_concurrent_backups": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
												},
												"storage_array_snapshot_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies Storage Array Snapshot Configuration.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_max_snapshots_config_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
															},
															"is_max_space_config_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the storage array snapshot max space config is enabled or not.",
															},
															"storage_array_snapshot_max_space_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies Storage Array Snapshot Max Space Config.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"max_snapshot_space_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Max number of storage snapshots allowed per volume/lun.",
																		},
																	},
																},
															},
															"storage_array_snapshot_throttling_policies": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies throttling policies configured for individual volume/lun.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the volume id of the storage array snapshot config.",
																		},
																		"is_max_snapshots_config_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																		},
																		"is_max_space_config_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																		},
																		"max_snapshot_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_snapshots": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Max number of storage snapshots allowed per volume/lun.",
																					},
																				},
																			},
																		},
																		"max_space_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies Storage Array Snapshot Max Space Config.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_snapshot_space_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Max number of storage snapshots allowed per volume/lun.",
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
									"throttling_policy_overrides": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies throttling policy override for a Datastore in a registered entity.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"datastore_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the Protection Source id of the Datastore.",
												},
												"datastore_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the display name of the Datastore.",
												},
												"throttling_policy": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the throttling policy for a registered Protection Source.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"enforce_max_streams": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether datastore streams are configured for all datastores that are part of the registered entity. If set to true, number of streams from Cohesity cluster to the registered entity will be limited to the value set for maxConcurrentStreams. If not set or set to false, there is no max limit for the number of concurrent streams.",
															},
															"enforce_registered_source_max_backups": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether no. of backups are configured for the registered entity. If set to true, number of backups made by Cohesity cluster in the registered entity will be limited to the value set for RegisteredSourceMaxConcurrentBackups. If not set or set to false, there is no max limit for the number of concurrent backups.",
															},
															"is_enabled": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Indicates whether read operations to the datastores, which are part of the registered Protection Source, are throttled.",
															},
															"latency_thresholds": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies latency thresholds that trigger throttling for all datastores found in the registered Protection Source or specific to one datastore.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"active_task_msecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "If the latency of a datastore is above this value, existing backup tasks using the datastore are throttled.",
																		},
																		"new_task_msecs": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "If the latency of a datastore is above this value, then new backup tasks using the datastore will not be started.",
																		},
																	},
																},
															},
															"max_concurrent_streams": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the limit on the number of streams Cohesity cluster will make concurrently to the datastores of the registered entity. This limit is enforced only when the flag enforceMaxStreams is set to true.",
															},
															"nas_source_params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the NAS specific source throttling parameters during source registration or during backup of the source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"max_parallel_metadata_fetch_full_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during full backup of the source.",
																		},
																		"max_parallel_metadata_fetch_incremental_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent metadata to be fetched during incremental backup of the source.",
																		},
																		"max_parallel_read_write_full_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent IO during full backup of the source.",
																		},
																		"max_parallel_read_write_incremental_percentage": &schema.Schema{
																			Type:        schema.TypeFloat,
																			Computed:    true,
																			Description: "Specifies the percentage value of maximum concurrent IO during incremental backup of the source.",
																		},
																	},
																},
															},
															"registered_source_max_concurrent_backups": &schema.Schema{
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Specifies the limit on the number of backups Cohesity cluster will make concurrently to the registered entity. This limit is enforced only when the flag enforceRegisteredSourceMaxBackups is set to true.",
															},
															"storage_array_snapshot_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies Storage Array Snapshot Configuration.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"is_max_snapshots_config_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																		},
																		"is_max_space_config_enabled": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																		},
																		"storage_array_snapshot_max_space_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies Storage Array Snapshot Max Space Config.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"max_snapshot_space_percentage": &schema.Schema{
																						Type:        schema.TypeFloat,
																						Computed:    true,
																						Description: "Max number of storage snapshots allowed per volume/lun.",
																					},
																				},
																			},
																		},
																		"storage_array_snapshot_throttling_policies": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies throttling policies configured for individual volume/lun.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"id": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the volume id of the storage array snapshot config.",
																					},
																					"is_max_snapshots_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max snapshots config is enabled or not.",
																					},
																					"is_max_space_config_enabled": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies if the storage array snapshot max space config is enabled or not.",
																					},
																					"max_snapshot_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Max Snapshots Config.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_snapshots": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Max number of storage snapshots allowed per volume/lun.",
																								},
																							},
																						},
																					},
																					"max_space_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies Storage Array Snapshot Max Space Config.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"max_snapshot_space_percentage": &schema.Schema{
																									Type:        schema.TypeFloat,
																									Computed:    true,
																									Description: "Max number of storage snapshots allowed per volume/lun.",
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
										},
									},
									"use_o_auth_for_exchange_online": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies whether OAuth should be used for authentication in case of Exchange Online.",
									},
									"use_vm_bios_uuid": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Specifies if registered vCenter is using BIOS UUID to track virtual machines.",
									},
									"user_messages": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the additional details encountered during registration. Though the registration may succeed, user messages imply the host environment requires some cleanup or fixing.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies username to access the target source.",
									},
									"vlan_params": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the VLAN configuration for Recovery.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"vlan": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Specifies the VLAN to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
												},
												"disable_vlan": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies whether to use the VIPs even when VLANs are configured on the Cluster. If configured, VLAN IP addresses are used by default. If VLANs are not configured, this flag is ignored. Set this flag to true to force using the partition VIPs when VLANs are configured on the Cluster.",
												},
												"interface_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the physical interface group name to use for mounting Cohesity's view on the remote host. If specified, Cohesity hostname or the IP address on this VLAN is used.",
												},
											},
										},
									},
									"warning_messages": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a list of warnings encountered during registration. Though the registration may succeed, warning messages imply the host environment requires some cleanup or fixing.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
						"total_downtiered_size_in_bytes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the total bytes downtiered from the source so far.",
						},
						"total_uptiered_size_in_bytes": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the total bytes uptiered to the source so far.",
						},
						"unprotected_sources_summary": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Aggregated information about a node subtree.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"environment": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the environment such as 'kSQL' or 'kVMware', where the Protection Source exists. Supported environment types such as 'kView', 'kSQL', 'kVMware', etc. NOTE 'kPuppeteer' refers to Cohesity's Remote Adapter. 'kVMware' indicates the VMware Protection Source environment. 'kHyperV' indicates the HyperV Protection Source environment. 'kSQL' indicates the SQL Protection Source environment. 'kView' indicates the View Protection Source environment. 'kPuppeteer' indicates the Cohesity's Remote Adapter. 'kPhysical' indicates the physical Protection Source environment. 'kPure' indicates the Pure Storage Protection Source environment. 'kNimble' indicates the Nimble Storage Protection Source environment. 'kAzure' indicates the Microsoft's Azure Protection Source environment. 'kNetapp' indicates the Netapp Protection Source environment. 'kAgent' indicates the Agent Protection Source environment. 'kGenericNas' indicates the Generic Network Attached Storage Protection Source environment. 'kAcropolis' indicates the Acropolis Protection Source environment. 'kPhysicalFiles' indicates the Physical Files Protection Source environment. 'kIbmFlashSystem' indicates the IBM Flash System Protection Source environment. 'kIsilon' indicates the Dell EMC's Isilon Protection Source environment. 'kGPFS' indicates IBM's GPFS Protection Source environment. 'kKVM' indicates the KVM Protection Source environment. 'kAWS' indicates the AWS Protection Source environment. 'kExchange' indicates the Exchange Protection Source environment. 'kHyperVVSS' indicates the HyperV VSS Protection Source environment. 'kOracle' indicates the Oracle Protection Source environment. 'kGCP' indicates the Google Cloud Platform Protection Source environment. 'kFlashBlade' indicates the Flash Blade Protection Source environment. 'kAWSNative' indicates the AWS Native Protection Source environment. 'kO365' indicates the Office 365 Protection Source environment. 'kO365Outlook' indicates Office 365 outlook Protection Source environment. 'kHyperFlex' indicates the Hyper Flex Protection Source environment. 'kGCPNative' indicates the GCP Native Protection Source environment. 'kAzureNative' indicates the Azure Native Protection Source environment. 'kKubernetes' indicates a Kubernetes Protection Source environment. 'kElastifile' indicates Elastifile Protection Source environment. 'kAD' indicates Active Directory Protection Source environment. 'kRDSSnapshotManager' indicates AWS RDS Protection Source environment. 'kCassandra' indicates Cassandra Protection Source environment. 'kMongoDB' indicates MongoDB Protection Source environment. 'kCouchbase' indicates Couchbase Protection Source environment. 'kHdfs' indicates Hdfs Protection Source environment. 'kHive' indicates Hive Protection Source environment. 'kHBase' indicates HBase Protection Source environment. 'kUDA' indicates Universal Data Adapter Protection Source environment. 'kO365Teams' indicates the Office365 Teams Protection Source environment. 'kO365Group' indicates the Office365 Groups Protection Source environment. 'kO365Exchange' indicates the Office365 Mailbox Protection Source environment. 'kO365OneDrive' indicates the Office365 OneDrive Protection Source environment. 'kO365Sharepoint' indicates the Office365 SharePoint Protection Source environment. 'kO365PublicFolders' indicates the Office365 PublicFolders Protection Source environment.",
									},
									"leaves_count": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of leaf nodes under the subtree of this node.",
									},
									"total_logical_size": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the total logical size of the data under the subtree of this node.",
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

func dataSourceIbmBackupRecoveryProtectionSourcesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_sources", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listProtectionSourcesOptions := &backuprecoveryv1.ListProtectionSourcesOptions{}

	listProtectionSourcesOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("exclude_office365_types"); ok {
		var excludeOffice365Types []string
		for _, v := range d.Get("exclude_office365_types").([]interface{}) {
			excludeOffice365TypesItem := v.(string)
			excludeOffice365Types = append(excludeOffice365Types, excludeOffice365TypesItem)
		}
		listProtectionSourcesOptions.SetExcludeOffice365Types(excludeOffice365Types)
	}
	if _, ok := d.GetOk("get_teams_channels"); ok {
		listProtectionSourcesOptions.SetGetTeamsChannels(d.Get("get_teams_channels").(bool))
	}
	if _, ok := d.GetOk("after_cursor_entity_id"); ok {
		listProtectionSourcesOptions.SetAfterCursorEntityID(int64(d.Get("after_cursor_entity_id").(int)))
	}
	if _, ok := d.GetOk("before_cursor_entity_id"); ok {
		listProtectionSourcesOptions.SetBeforeCursorEntityID(int64(d.Get("before_cursor_entity_id").(int)))
	}
	if _, ok := d.GetOk("node_id"); ok {
		listProtectionSourcesOptions.SetNodeID(int64(d.Get("node_id").(int)))
	}
	if _, ok := d.GetOk("page_size"); ok {
		listProtectionSourcesOptions.SetPageSize(int64(d.Get("page_size").(int)))
	}
	if _, ok := d.GetOk("has_valid_mailbox"); ok {
		listProtectionSourcesOptions.SetHasValidMailbox(d.Get("has_valid_mailbox").(bool))
	}
	if _, ok := d.GetOk("has_valid_onedrive"); ok {
		listProtectionSourcesOptions.SetHasValidOnedrive(d.Get("has_valid_onedrive").(bool))
	}
	if _, ok := d.GetOk("is_security_group"); ok {
		listProtectionSourcesOptions.SetIsSecurityGroup(d.Get("is_security_group").(bool))
	}
	if _, ok := d.GetOk("backup_recovery_protection_sources_id"); ok {
		listProtectionSourcesOptions.SetID(int64(d.Get("backup_recovery_protection_sources_id").(int)))
	}
	if _, ok := d.GetOk("num_levels"); ok {
		listProtectionSourcesOptions.SetNumLevels(d.Get("num_levels").(float64))
	}
	if _, ok := d.GetOk("exclude_types"); ok {
		var excludeTypes []string
		for _, v := range d.Get("exclude_types").([]interface{}) {
			excludeTypesItem := v.(string)
			excludeTypes = append(excludeTypes, excludeTypesItem)
		}
		listProtectionSourcesOptions.SetExcludeTypes(excludeTypes)
	}
	if _, ok := d.GetOk("exclude_aws_types"); ok {
		var excludeAwsTypes []string
		for _, v := range d.Get("exclude_aws_types").([]interface{}) {
			excludeAwsTypesItem := v.(string)
			excludeAwsTypes = append(excludeAwsTypes, excludeAwsTypesItem)
		}
		listProtectionSourcesOptions.SetExcludeAwsTypes(excludeAwsTypes)
	}
	if _, ok := d.GetOk("exclude_kubernetes_types"); ok {
		var excludeKubernetesTypes []string
		for _, v := range d.Get("exclude_kubernetes_types").([]interface{}) {
			excludeKubernetesTypesItem := v.(string)
			excludeKubernetesTypes = append(excludeKubernetesTypes, excludeKubernetesTypesItem)
		}
		listProtectionSourcesOptions.SetExcludeKubernetesTypes(excludeKubernetesTypes)
	}
	if _, ok := d.GetOk("include_datastores"); ok {
		listProtectionSourcesOptions.SetIncludeDatastores(d.Get("include_datastores").(bool))
	}
	if _, ok := d.GetOk("include_networks"); ok {
		listProtectionSourcesOptions.SetIncludeNetworks(d.Get("include_networks").(bool))
	}
	if _, ok := d.GetOk("include_vm_folders"); ok {
		listProtectionSourcesOptions.SetIncludeVMFolders(d.Get("include_vm_folders").(bool))
	}
	if _, ok := d.GetOk("include_sfdc_fields"); ok {
		listProtectionSourcesOptions.SetIncludeSfdcFields(d.Get("include_sfdc_fields").(bool))
	}
	if _, ok := d.GetOk("include_system_v_apps"); ok {
		listProtectionSourcesOptions.SetIncludeSystemVApps(d.Get("include_system_v_apps").(bool))
	}
	if _, ok := d.GetOk("environments"); ok {
		var environments []string
		for _, v := range d.Get("environments").([]interface{}) {
			environmentsItem := v.(string)
			environments = append(environments, environmentsItem)
		}
		listProtectionSourcesOptions.SetEnvironments(environments)
	}
	if _, ok := d.GetOk("environment"); ok {
		listProtectionSourcesOptions.SetEnvironment(d.Get("environment").(string))
	}
	if _, ok := d.GetOk("include_entity_permission_info"); ok {
		listProtectionSourcesOptions.SetIncludeEntityPermissionInfo(d.Get("include_entity_permission_info").(bool))
	}
	if _, ok := d.GetOk("sids"); ok {
		var sids []string
		for _, v := range d.Get("sids").([]interface{}) {
			sidsItem := v.(string)
			sids = append(sids, sidsItem)
		}
		listProtectionSourcesOptions.SetSids(sids)
	}
	if _, ok := d.GetOk("include_source_credentials"); ok {
		listProtectionSourcesOptions.SetIncludeSourceCredentials(d.Get("include_source_credentials").(bool))
	}
	if _, ok := d.GetOk("encryption_key"); ok {
		listProtectionSourcesOptions.SetEncryptionKey(d.Get("encryption_key").(string))
	}
	if _, ok := d.GetOk("include_object_protection_info"); ok {
		listProtectionSourcesOptions.SetIncludeObjectProtectionInfo(d.Get("include_object_protection_info").(bool))
	}
	if _, ok := d.GetOk("prune_non_critical_info"); ok {
		listProtectionSourcesOptions.SetPruneNonCriticalInfo(d.Get("prune_non_critical_info").(bool))
	}
	if _, ok := d.GetOk("prune_aggregation_info"); ok {
		listProtectionSourcesOptions.SetPruneAggregationInfo(d.Get("prune_aggregation_info").(bool))
	}
	if _, ok := d.GetOk("request_initiator_type"); ok {
		listProtectionSourcesOptions.SetRequestInitiatorType(d.Get("request_initiator_type").(string))
	}
	if _, ok := d.GetOk("use_cached_data"); ok {
		listProtectionSourcesOptions.SetUseCachedData(d.Get("use_cached_data").(bool))
	}
	if _, ok := d.GetOk("all_under_hierarchy"); ok {
		listProtectionSourcesOptions.SetAllUnderHierarchy(d.Get("all_under_hierarchy").(bool))
	}

	protectionSourcesResponse, _, err := backupRecoveryClient.ListProtectionSourcesWithContext(context, listProtectionSourcesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListProtectionSourcesWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_protection_sources", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryProtectionSourcesID(d))

	if !core.IsNil(protectionSourcesResponse) {
		protectionSources := []map[string]interface{}{}
		for _, protectionSourcesItem := range protectionSourcesResponse {
			protectionSourcesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceNodesToMap(&protectionSourcesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_sources", "read", "protection_sources-to-map").GetDiag()
			}
			protectionSources = append(protectionSources, protectionSourcesItemMap)
		}
		if err = d.Set("protection_sources", protectionSources); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting protection_sources: %s", err), "(Data) ibm_backup_recovery_protection_sources", "read", "set-protection_sources").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryProtectionSourcesID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryProtectionSourcesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceNodesToMap(model *backuprecoveryv1.ProtectionSourceNodes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApplicationNodes != nil {
		applicationNodesMap, err := DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceApplicationNodesToMap(model.ApplicationNodes)
		if err != nil {
			return modelMap, err
		}
		modelMap["application_nodes"] = []map[string]interface{}{applicationNodesMap}
	}
	if model.EntityPaginationParameters != nil {
		entityPaginationParametersMap, err := DataSourceIbmBackupRecoveryProtectionSourcesPaginationParametersToMap(model.EntityPaginationParameters)
		if err != nil {
			return modelMap, err
		}
		modelMap["entity_pagination_parameters"] = []map[string]interface{}{entityPaginationParametersMap}
	}
	if model.EntityPermissionInfo != nil {
		entityPermissionInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesEntityPermissionInfoToMap(model.EntityPermissionInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["entity_permission_info"] = []map[string]interface{}{entityPermissionInfoMap}
	}
	if model.LogicalSize != nil {
		modelMap["logical_size"] = flex.IntValue(model.LogicalSize)
	}
	if model.Nodes != nil {
		nodes := []map[string]interface{}{}
		for _, nodesItem := range model.Nodes {
			nodesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceNodesToMap(&nodesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			nodes = append(nodes, nodesItemMap)
		}
		modelMap["nodes"] = nodes
	}
	if model.ObjectProtectionInfo != nil {
		objectProtectionInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceObjectProtectionInfoToMap(model.ObjectProtectionInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["object_protection_info"] = []map[string]interface{}{objectProtectionInfoMap}
	}
	if model.ProtectedSourcesSummary != nil {
		protectedSourcesSummary := []map[string]interface{}{}
		for _, protectedSourcesSummaryItem := range model.ProtectedSourcesSummary {
			protectedSourcesSummaryItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesSubtreeInfoToMap(&protectedSourcesSummaryItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			protectedSourcesSummary = append(protectedSourcesSummary, protectedSourcesSummaryItemMap)
		}
		modelMap["protected_sources_summary"] = protectedSourcesSummary
	}
	if model.ProtectionSource != nil {
		protectionSourceMap, err := DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceNodeToMap(model.ProtectionSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["protection_source"] = []map[string]interface{}{protectionSourceMap}
	}
	if model.RegistrationInfo != nil {
		registrationInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesAgentRegistrationInfoToMap(model.RegistrationInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["registration_info"] = []map[string]interface{}{registrationInfoMap}
	}
	if model.TotalDowntieredSizeInBytes != nil {
		modelMap["total_downtiered_size_in_bytes"] = flex.IntValue(model.TotalDowntieredSizeInBytes)
	}
	if model.TotalUptieredSizeInBytes != nil {
		modelMap["total_uptiered_size_in_bytes"] = flex.IntValue(model.TotalUptieredSizeInBytes)
	}
	if model.UnprotectedSourcesSummary != nil {
		unprotectedSourcesSummary := []map[string]interface{}{}
		for _, unprotectedSourcesSummaryItem := range model.UnprotectedSourcesSummary {
			unprotectedSourcesSummaryItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesSubtreeInfoToMap(&unprotectedSourcesSummaryItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			unprotectedSourcesSummary = append(unprotectedSourcesSummary, unprotectedSourcesSummaryItemMap)
		}
		modelMap["unprotected_sources_summary"] = unprotectedSourcesSummary
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceApplicationNodesToMap(model *backuprecoveryv1.ProtectionSourceApplicationNodes) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Nodes != nil {
		nodes := []map[string]interface{}{}
		for _, nodesItem := range model.Nodes {
			nodesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceNodesToMap(&nodesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			nodes = append(nodes, nodesItemMap)
		}
		modelMap["nodes"] = nodes
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesPaginationParametersToMap(model *backuprecoveryv1.PaginationParameters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AfterCursorEntityID != nil {
		modelMap["after_cursor_entity_id"] = flex.IntValue(model.AfterCursorEntityID)
	}
	if model.BeforeCursorEntityID != nil {
		modelMap["before_cursor_entity_id"] = flex.IntValue(model.BeforeCursorEntityID)
	}
	if model.NodeID != nil {
		modelMap["node_id"] = flex.IntValue(model.NodeID)
	}
	if model.PageSize != nil {
		modelMap["page_size"] = flex.IntValue(model.PageSize)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesEntityPermissionInfoToMap(model *backuprecoveryv1.EntityPermissionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EntityID != nil {
		modelMap["entity_id"] = flex.IntValue(model.EntityID)
	}
	if model.Groups != nil {
		groups := []map[string]interface{}{}
		for _, groupsItem := range model.Groups {
			groupsItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesEntityGroupParamsToMap(&groupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.IsInferred != nil {
		modelMap["is_inferred"] = *model.IsInferred
	}
	if model.IsRegisteredBySp != nil {
		modelMap["is_registered_by_sp"] = *model.IsRegisteredBySp
	}
	if model.RegisteringTenantID != nil {
		modelMap["registering_tenant_id"] = *model.RegisteringTenantID
	}
	if model.Tenant != nil {
		tenantMap, err := DataSourceIbmBackupRecoveryProtectionSourcesEntityTenantInfoToMap(model.Tenant)
		if err != nil {
			return modelMap, err
		}
		modelMap["tenant"] = []map[string]interface{}{tenantMap}
	}
	if model.Users != nil {
		users := []map[string]interface{}{}
		for _, usersItem := range model.Users {
			usersItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesEntityUserInfoToMap(&usersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			users = append(users, usersItemMap)
		}
		modelMap["users"] = users
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesEntityGroupParamsToMap(model *backuprecoveryv1.EntityGroupParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Domain != nil {
		modelMap["domain"] = *model.Domain
	}
	if model.GroupName != nil {
		modelMap["group_name"] = *model.GroupName
	}
	if model.Sid != nil {
		modelMap["sid"] = *model.Sid
	}
	if model.TenantIds != nil {
		modelMap["tenant_ids"] = model.TenantIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesEntityTenantInfoToMap(model *backuprecoveryv1.EntityTenantInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BifrostEnabled != nil {
		modelMap["bifrost_enabled"] = *model.BifrostEnabled
	}
	if model.IsManagedOnHelios != nil {
		modelMap["is_managed_on_helios"] = *model.IsManagedOnHelios
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.TenantID != nil {
		modelMap["tenant_id"] = *model.TenantID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesEntityUserInfoToMap(model *backuprecoveryv1.EntityUserInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Domain != nil {
		modelMap["domain"] = *model.Domain
	}
	if model.Sid != nil {
		modelMap["sid"] = *model.Sid
	}
	if model.TenantID != nil {
		modelMap["tenant_id"] = *model.TenantID
	}
	if model.UserName != nil {
		modelMap["user_name"] = *model.UserName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceObjectProtectionInfoToMap(model *backuprecoveryv1.ProtectionSourceObjectProtectionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AutoProtectParentID != nil {
		modelMap["auto_protect_parent_id"] = flex.IntValue(model.AutoProtectParentID)
	}
	if model.EntityID != nil {
		modelMap["entity_id"] = flex.IntValue(model.EntityID)
	}
	if model.HasActiveObjectProtectionSpec != nil {
		modelMap["has_active_object_protection_spec"] = flex.IntValue(model.HasActiveObjectProtectionSpec)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesSubtreeInfoToMap(model *backuprecoveryv1.SubtreeInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.LeavesCount != nil {
		modelMap["leaves_count"] = flex.IntValue(model.LeavesCount)
	}
	if model.TotalLogicalSize != nil {
		modelMap["total_logical_size"] = flex.IntValue(model.TotalLogicalSize)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesProtectionSourceNodeToMap(model *backuprecoveryv1.ProtectionSourceNode) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConnectionID != nil {
		modelMap["connection_id"] = flex.IntValue(model.ConnectionID)
	}
	if model.ConnectorGroupID != nil {
		modelMap["connector_group_id"] = flex.IntValue(model.ConnectorGroupID)
	}
	if model.CustomName != nil {
		modelMap["custom_name"] = *model.CustomName
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ParentID != nil {
		modelMap["parent_id"] = flex.IntValue(model.ParentID)
	}
	if model.PhysicalProtectionSource != nil {
		physicalProtectionSourceMap, err := DataSourceIbmBackupRecoveryProtectionSourcesPhysicalProtectionSourceToMap(model.PhysicalProtectionSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_protection_source"] = []map[string]interface{}{physicalProtectionSourceMap}
	}
	if model.SqlProtectionSource != nil {
		sqlProtectionSourceMap, err := DataSourceIbmBackupRecoveryProtectionSourcesSqlProtectionSourceToMap(model.SqlProtectionSource)
		if err != nil {
			return modelMap, err
		}
		modelMap["sql_protection_source"] = []map[string]interface{}{sqlProtectionSourceMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesPhysicalProtectionSourceToMap(model *backuprecoveryv1.PhysicalProtectionSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Agents != nil {
		agents := []map[string]interface{}{}
		for _, agentsItem := range model.Agents {
			agentsItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesAgentInformationToMap(&agentsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			agents = append(agents, agentsItemMap)
		}
		modelMap["agents"] = agents
	}
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	if model.HostName != nil {
		modelMap["host_name"] = *model.HostName
	}
	if model.HostType != nil {
		modelMap["host_type"] = *model.HostType
	}
	if model.ID != nil {
		idMap, err := DataSourceIbmBackupRecoveryProtectionSourcesUniqueGlobalIDToMap(model.ID)
		if err != nil {
			return modelMap, err
		}
		modelMap["id"] = []map[string]interface{}{idMap}
	}
	if model.IsProxyHost != nil {
		modelMap["is_proxy_host"] = *model.IsProxyHost
	}
	if model.MemorySizeBytes != nil {
		modelMap["memory_size_bytes"] = flex.IntValue(model.MemorySizeBytes)
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.NetworkingInfo != nil {
		networkingInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesNetworkingInformationToMap(model.NetworkingInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["networking_info"] = []map[string]interface{}{networkingInfoMap}
	}
	if model.NumProcessors != nil {
		modelMap["num_processors"] = flex.IntValue(model.NumProcessors)
	}
	if model.OsName != nil {
		modelMap["os_name"] = *model.OsName
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.VcsVersion != nil {
		modelMap["vcs_version"] = *model.VcsVersion
	}
	if model.Volumes != nil {
		volumes := []map[string]interface{}{}
		for _, volumesItem := range model.Volumes {
			volumesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesPhysicalVolumeToMap(&volumesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			volumes = append(volumes, volumesItemMap)
		}
		modelMap["volumes"] = volumes
	}
	if model.Vsswriters != nil {
		vsswriters := []map[string]interface{}{}
		for _, vsswritersItem := range model.Vsswriters {
			vsswritersItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesVssWritersToMap(&vsswritersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			vsswriters = append(vsswriters, vsswritersItemMap)
		}
		modelMap["vsswriters"] = vsswriters
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesAgentInformationToMap(model *backuprecoveryv1.AgentInformation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CbmrVersion != nil {
		modelMap["cbmr_version"] = *model.CbmrVersion
	}
	if model.FileCbtInfo != nil {
		fileCbtInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesCbtInfoToMap(model.FileCbtInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["file_cbt_info"] = []map[string]interface{}{fileCbtInfoMap}
	}
	if model.HostType != nil {
		modelMap["host_type"] = *model.HostType
	}
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.OracleMultiNodeChannelSupported != nil {
		modelMap["oracle_multi_node_channel_supported"] = *model.OracleMultiNodeChannelSupported
	}
	if model.RegistrationInfo != nil {
		registrationInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesAgentRegistrationInfoToMap(model.RegistrationInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["registration_info"] = []map[string]interface{}{registrationInfoMap}
	}
	if model.SourceSideDedupEnabled != nil {
		modelMap["source_side_dedup_enabled"] = *model.SourceSideDedupEnabled
	}
	if model.Status != nil {
		modelMap["status"] = *model.Status
	}
	if model.StatusMessage != nil {
		modelMap["status_message"] = *model.StatusMessage
	}
	if model.Upgradability != nil {
		modelMap["upgradability"] = *model.Upgradability
	}
	if model.UpgradeStatus != nil {
		modelMap["upgrade_status"] = *model.UpgradeStatus
	}
	if model.UpgradeStatusMessage != nil {
		modelMap["upgrade_status_message"] = *model.UpgradeStatusMessage
	}
	if model.Version != nil {
		modelMap["version"] = *model.Version
	}
	if model.VolCbtInfo != nil {
		volCbtInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesCbtInfoToMap(model.VolCbtInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["vol_cbt_info"] = []map[string]interface{}{volCbtInfoMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesCbtInfoToMap(model *backuprecoveryv1.CbtInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FileVersion != nil {
		fileVersionMap, err := DataSourceIbmBackupRecoveryProtectionSourcesCbtFileVersionToMap(model.FileVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["file_version"] = []map[string]interface{}{fileVersionMap}
	}
	if model.IsInstalled != nil {
		modelMap["is_installed"] = *model.IsInstalled
	}
	if model.RebootStatus != nil {
		modelMap["reboot_status"] = *model.RebootStatus
	}
	if model.ServiceState != nil {
		serviceStateMap, err := DataSourceIbmBackupRecoveryProtectionSourcesCbtServiceStateToMap(model.ServiceState)
		if err != nil {
			return modelMap, err
		}
		modelMap["service_state"] = []map[string]interface{}{serviceStateMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesCbtFileVersionToMap(model *backuprecoveryv1.CbtFileVersion) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.BuildVer != nil {
		modelMap["build_ver"] = *model.BuildVer
	}
	if model.MajorVer != nil {
		modelMap["major_ver"] = *model.MajorVer
	}
	if model.MinorVer != nil {
		modelMap["minor_ver"] = *model.MinorVer
	}
	if model.RevisionNum != nil {
		modelMap["revision_num"] = *model.RevisionNum
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesCbtServiceStateToMap(model *backuprecoveryv1.CbtServiceState) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["state"] = *model.State
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesAgentRegistrationInfoToMap(model *backuprecoveryv1.AgentRegistrationInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AccessInfo != nil {
		accessInfoMap, err := DataSourceIbmBackupRecoveryProtectionSourcesAgentAccessInfoToMap(model.AccessInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["access_info"] = []map[string]interface{}{accessInfoMap}
	}
	if model.AllowedIpAddresses != nil {
		modelMap["allowed_ip_addresses"] = model.AllowedIpAddresses
	}
	if model.AuthenticationErrorMessage != nil {
		modelMap["authentication_error_message"] = *model.AuthenticationErrorMessage
	}
	if model.AuthenticationStatus != nil {
		modelMap["authentication_status"] = *model.AuthenticationStatus
	}
	if model.BlacklistedIpAddresses != nil {
		modelMap["blacklisted_ip_addresses"] = model.BlacklistedIpAddresses
	}
	if model.DeniedIpAddresses != nil {
		modelMap["denied_ip_addresses"] = model.DeniedIpAddresses
	}
	if model.Environments != nil {
		modelMap["environments"] = model.Environments
	}
	if model.IsDbAuthenticated != nil {
		modelMap["is_db_authenticated"] = *model.IsDbAuthenticated
	}
	if model.IsStorageArraySnapshotEnabled != nil {
		modelMap["is_storage_array_snapshot_enabled"] = *model.IsStorageArraySnapshotEnabled
	}
	if model.LinkVmsAcrossVcenter != nil {
		modelMap["link_vms_across_vcenter"] = *model.LinkVmsAcrossVcenter
	}
	if model.MinimumFreeSpaceGB != nil {
		modelMap["minimum_free_space_gb"] = flex.IntValue(model.MinimumFreeSpaceGB)
	}
	if model.MinimumFreeSpacePercent != nil {
		modelMap["minimum_free_space_percent"] = flex.IntValue(model.MinimumFreeSpacePercent)
	}
	if model.Password != nil {
		modelMap["password"] = *model.Password
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoveryProtectionSourcesAgentPhysicalParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	if model.ProgressMonitorPath != nil {
		modelMap["progress_monitor_path"] = *model.ProgressMonitorPath
	}
	if model.RefreshErrorMessage != nil {
		modelMap["refresh_error_message"] = *model.RefreshErrorMessage
	}
	if model.RefreshTimeUsecs != nil {
		modelMap["refresh_time_usecs"] = flex.IntValue(model.RefreshTimeUsecs)
	}
	if model.RegisteredAppsInfo != nil {
		registeredAppsInfo := []map[string]interface{}{}
		for _, registeredAppsInfoItem := range model.RegisteredAppsInfo {
			registeredAppsInfoItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesRegisteredAppInfoToMap(&registeredAppsInfoItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			registeredAppsInfo = append(registeredAppsInfo, registeredAppsInfoItemMap)
		}
		modelMap["registered_apps_info"] = registeredAppsInfo
	}
	if model.RegistrationTimeUsecs != nil {
		modelMap["registration_time_usecs"] = flex.IntValue(model.RegistrationTimeUsecs)
	}
	if model.Subnets != nil {
		subnets := []map[string]interface{}{}
		for _, subnetsItem := range model.Subnets {
			subnetsItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesSubnetToMap(&subnetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			subnets = append(subnets, subnetsItemMap)
		}
		modelMap["subnets"] = subnets
	}
	if model.ThrottlingPolicy != nil {
		throttlingPolicyMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingPolicyToMap(model.ThrottlingPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["throttling_policy"] = []map[string]interface{}{throttlingPolicyMap}
	}
	if model.ThrottlingPolicyOverrides != nil {
		throttlingPolicyOverrides := []map[string]interface{}{}
		for _, throttlingPolicyOverridesItem := range model.ThrottlingPolicyOverrides {
			throttlingPolicyOverridesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingPolicyOverridesToMap(&throttlingPolicyOverridesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			throttlingPolicyOverrides = append(throttlingPolicyOverrides, throttlingPolicyOverridesItemMap)
		}
		modelMap["throttling_policy_overrides"] = throttlingPolicyOverrides
	}
	if model.UseOAuthForExchangeOnline != nil {
		modelMap["use_o_auth_for_exchange_online"] = *model.UseOAuthForExchangeOnline
	}
	if model.UseVmBiosUUID != nil {
		modelMap["use_vm_bios_uuid"] = *model.UseVmBiosUUID
	}
	if model.UserMessages != nil {
		modelMap["user_messages"] = model.UserMessages
	}
	if model.Username != nil {
		modelMap["username"] = *model.Username
	}
	if model.VlanParams != nil {
		vlanParamsMap, err := DataSourceIbmBackupRecoveryProtectionSourcesRegisteredSourceVlanConfigToMap(model.VlanParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["vlan_params"] = []map[string]interface{}{vlanParamsMap}
	}
	if model.WarningMessages != nil {
		modelMap["warning_messages"] = model.WarningMessages
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesAgentAccessInfoToMap(model *backuprecoveryv1.AgentAccessInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ConnectionID != nil {
		modelMap["connection_id"] = flex.IntValue(model.ConnectionID)
	}
	if model.ConnectorGroupID != nil {
		modelMap["connector_group_id"] = flex.IntValue(model.ConnectorGroupID)
	}
	if model.Endpoint != nil {
		modelMap["endpoint"] = *model.Endpoint
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.Version != nil {
		modelMap["version"] = flex.IntValue(model.Version)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesAgentPhysicalParamsToMap(model *backuprecoveryv1.AgentPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Applications != nil {
		modelMap["applications"] = model.Applications
	}
	if model.Password != nil {
		modelMap["password"] = *model.Password
	}
	if model.ThrottlingConfig != nil {
		throttlingConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingConfigToMap(model.ThrottlingConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["throttling_config"] = []map[string]interface{}{throttlingConfigMap}
	}
	if model.Username != nil {
		modelMap["username"] = *model.Username
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesThrottlingConfigToMap(model *backuprecoveryv1.ThrottlingConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CpuThrottlingConfig != nil {
		cpuThrottlingConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingConfigurationParamsToMap(model.CpuThrottlingConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["cpu_throttling_config"] = []map[string]interface{}{cpuThrottlingConfigMap}
	}
	if model.NetworkThrottlingConfig != nil {
		networkThrottlingConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingConfigurationParamsToMap(model.NetworkThrottlingConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["network_throttling_config"] = []map[string]interface{}{networkThrottlingConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesThrottlingConfigurationParamsToMap(model *backuprecoveryv1.ThrottlingConfigurationParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FixedThreshold != nil {
		modelMap["fixed_threshold"] = flex.IntValue(model.FixedThreshold)
	}
	if model.PatternType != nil {
		modelMap["pattern_type"] = *model.PatternType
	}
	if model.ThrottlingWindows != nil {
		throttlingWindows := []map[string]interface{}{}
		for _, throttlingWindowsItem := range model.ThrottlingWindows {
			throttlingWindowsItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingWindowToMap(&throttlingWindowsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			throttlingWindows = append(throttlingWindows, throttlingWindowsItemMap)
		}
		modelMap["throttling_windows"] = throttlingWindows
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesThrottlingWindowToMap(model *backuprecoveryv1.ThrottlingWindow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DayTimeWindow != nil {
		dayTimeWindowMap, err := DataSourceIbmBackupRecoveryProtectionSourcesDayTimeWindowToMap(model.DayTimeWindow)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_time_window"] = []map[string]interface{}{dayTimeWindowMap}
	}
	if model.Threshold != nil {
		modelMap["threshold"] = flex.IntValue(model.Threshold)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesDayTimeWindowToMap(model *backuprecoveryv1.DayTimeWindow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EndTime != nil {
		endTimeMap, err := DataSourceIbmBackupRecoveryProtectionSourcesDayTimeParamsToMap(model.EndTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["end_time"] = []map[string]interface{}{endTimeMap}
	}
	if model.StartTime != nil {
		startTimeMap, err := DataSourceIbmBackupRecoveryProtectionSourcesDayTimeParamsToMap(model.StartTime)
		if err != nil {
			return modelMap, err
		}
		modelMap["start_time"] = []map[string]interface{}{startTimeMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesDayTimeParamsToMap(model *backuprecoveryv1.DayTimeParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Day != nil {
		modelMap["day"] = *model.Day
	}
	if model.Time != nil {
		timeMap, err := DataSourceIbmBackupRecoveryProtectionSourcesTimeToMap(model.Time)
		if err != nil {
			return modelMap, err
		}
		modelMap["time"] = []map[string]interface{}{timeMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesTimeToMap(model *backuprecoveryv1.Time) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Hour != nil {
		modelMap["hour"] = flex.IntValue(model.Hour)
	}
	if model.Minute != nil {
		modelMap["minute"] = flex.IntValue(model.Minute)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesRegisteredAppInfoToMap(model *backuprecoveryv1.RegisteredAppInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AuthenticationErrorMessage != nil {
		modelMap["authentication_error_message"] = *model.AuthenticationErrorMessage
	}
	if model.AuthenticationStatus != nil {
		modelMap["authentication_status"] = *model.AuthenticationStatus
	}
	if model.Environment != nil {
		modelMap["environment"] = *model.Environment
	}
	if model.HostSettingsCheckResults != nil {
		hostSettingsCheckResults := []map[string]interface{}{}
		for _, hostSettingsCheckResultsItem := range model.HostSettingsCheckResults {
			hostSettingsCheckResultsItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesHostSettingsCheckResultToMap(&hostSettingsCheckResultsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			hostSettingsCheckResults = append(hostSettingsCheckResults, hostSettingsCheckResultsItemMap)
		}
		modelMap["host_settings_check_results"] = hostSettingsCheckResults
	}
	if model.RefreshErrorMessage != nil {
		modelMap["refresh_error_message"] = *model.RefreshErrorMessage
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesHostSettingsCheckResultToMap(model *backuprecoveryv1.HostSettingsCheckResult) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CheckType != nil {
		modelMap["check_type"] = *model.CheckType
	}
	if model.ResultType != nil {
		modelMap["result_type"] = *model.ResultType
	}
	if model.UserMessage != nil {
		modelMap["user_message"] = *model.UserMessage
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesSubnetToMap(model *backuprecoveryv1.Subnet) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Component != nil {
		modelMap["component"] = *model.Component
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Ip != nil {
		modelMap["ip"] = *model.Ip
	}
	if model.NetmaskBits != nil {
		modelMap["netmask_bits"] = *model.NetmaskBits
	}
	if model.NetmaskIp4 != nil {
		modelMap["netmask_ip4"] = *model.NetmaskIp4
	}
	if model.NfsAccess != nil {
		modelMap["nfs_access"] = *model.NfsAccess
	}
	if model.NfsAllSquash != nil {
		modelMap["nfs_all_squash"] = *model.NfsAllSquash
	}
	if model.NfsRootSquash != nil {
		modelMap["nfs_root_squash"] = *model.NfsRootSquash
	}
	if model.S3Access != nil {
		modelMap["s3_access"] = *model.S3Access
	}
	if model.SmbAccess != nil {
		modelMap["smb_access"] = *model.SmbAccess
	}
	if model.TenantID != nil {
		modelMap["tenant_id"] = *model.TenantID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesThrottlingPolicyToMap(model *backuprecoveryv1.ThrottlingPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnforceMaxStreams != nil {
		modelMap["enforce_max_streams"] = *model.EnforceMaxStreams
	}
	if model.EnforceRegisteredSourceMaxBackups != nil {
		modelMap["enforce_registered_source_max_backups"] = *model.EnforceRegisteredSourceMaxBackups
	}
	if model.IsEnabled != nil {
		modelMap["is_enabled"] = *model.IsEnabled
	}
	if model.LatencyThresholds != nil {
		latencyThresholdsMap, err := DataSourceIbmBackupRecoveryProtectionSourcesLatencyThresholdsToMap(model.LatencyThresholds)
		if err != nil {
			return modelMap, err
		}
		modelMap["latency_thresholds"] = []map[string]interface{}{latencyThresholdsMap}
	}
	if model.MaxConcurrentStreams != nil {
		modelMap["max_concurrent_streams"] = *model.MaxConcurrentStreams
	}
	if model.NasSourceParams != nil {
		nasSourceParamsMap, err := DataSourceIbmBackupRecoveryProtectionSourcesNasSourceParamsToMap(model.NasSourceParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["nas_source_params"] = []map[string]interface{}{nasSourceParamsMap}
	}
	if model.RegisteredSourceMaxConcurrentBackups != nil {
		modelMap["registered_source_max_concurrent_backups"] = *model.RegisteredSourceMaxConcurrentBackups
	}
	if model.StorageArraySnapshotConfig != nil {
		storageArraySnapshotConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesStorageArraySnapshotConfigToMap(model.StorageArraySnapshotConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_array_snapshot_config"] = []map[string]interface{}{storageArraySnapshotConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesLatencyThresholdsToMap(model *backuprecoveryv1.LatencyThresholds) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ActiveTaskMsecs != nil {
		modelMap["active_task_msecs"] = flex.IntValue(model.ActiveTaskMsecs)
	}
	if model.NewTaskMsecs != nil {
		modelMap["new_task_msecs"] = flex.IntValue(model.NewTaskMsecs)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesNasSourceParamsToMap(model *backuprecoveryv1.NasSourceParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxParallelMetadataFetchFullPercentage != nil {
		modelMap["max_parallel_metadata_fetch_full_percentage"] = *model.MaxParallelMetadataFetchFullPercentage
	}
	if model.MaxParallelMetadataFetchIncrementalPercentage != nil {
		modelMap["max_parallel_metadata_fetch_incremental_percentage"] = *model.MaxParallelMetadataFetchIncrementalPercentage
	}
	if model.MaxParallelReadWriteFullPercentage != nil {
		modelMap["max_parallel_read_write_full_percentage"] = *model.MaxParallelReadWriteFullPercentage
	}
	if model.MaxParallelReadWriteIncrementalPercentage != nil {
		modelMap["max_parallel_read_write_incremental_percentage"] = *model.MaxParallelReadWriteIncrementalPercentage
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesStorageArraySnapshotConfigToMap(model *backuprecoveryv1.StorageArraySnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsMaxSnapshotsConfigEnabled != nil {
		modelMap["is_max_snapshots_config_enabled"] = *model.IsMaxSnapshotsConfigEnabled
	}
	if model.IsMaxSpaceConfigEnabled != nil {
		modelMap["is_max_space_config_enabled"] = *model.IsMaxSpaceConfigEnabled
	}
	if model.StorageArraySnapshotMaxSpaceConfig != nil {
		storageArraySnapshotMaxSpaceConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesStorageArraySnapshotMaxSpaceConfigToMap(model.StorageArraySnapshotMaxSpaceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_array_snapshot_max_space_config"] = []map[string]interface{}{storageArraySnapshotMaxSpaceConfigMap}
	}
	if model.StorageArraySnapshotThrottlingPolicies != nil {
		storageArraySnapshotThrottlingPolicies := []map[string]interface{}{}
		for _, storageArraySnapshotThrottlingPoliciesItem := range model.StorageArraySnapshotThrottlingPolicies {
			storageArraySnapshotThrottlingPoliciesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesStorageArraySnapshotThrottlingPoliciesToMap(&storageArraySnapshotThrottlingPoliciesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			storageArraySnapshotThrottlingPolicies = append(storageArraySnapshotThrottlingPolicies, storageArraySnapshotThrottlingPoliciesItemMap)
		}
		modelMap["storage_array_snapshot_throttling_policies"] = storageArraySnapshotThrottlingPolicies
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesStorageArraySnapshotMaxSpaceConfigToMap(model *backuprecoveryv1.StorageArraySnapshotMaxSpaceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSnapshotSpacePercentage != nil {
		modelMap["max_snapshot_space_percentage"] = *model.MaxSnapshotSpacePercentage
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesStorageArraySnapshotThrottlingPoliciesToMap(model *backuprecoveryv1.StorageArraySnapshotThrottlingPolicies) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	if model.IsMaxSnapshotsConfigEnabled != nil {
		modelMap["is_max_snapshots_config_enabled"] = *model.IsMaxSnapshotsConfigEnabled
	}
	if model.IsMaxSpaceConfigEnabled != nil {
		modelMap["is_max_space_config_enabled"] = *model.IsMaxSpaceConfigEnabled
	}
	if model.MaxSnapshotConfig != nil {
		maxSnapshotConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesMaxSnapshotConfigToMap(model.MaxSnapshotConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["max_snapshot_config"] = []map[string]interface{}{maxSnapshotConfigMap}
	}
	if model.MaxSpaceConfig != nil {
		maxSpaceConfigMap, err := DataSourceIbmBackupRecoveryProtectionSourcesMaxSpaceConfigToMap(model.MaxSpaceConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["max_space_config"] = []map[string]interface{}{maxSpaceConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesMaxSnapshotConfigToMap(model *backuprecoveryv1.MaxSnapshotConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSnapshots != nil {
		modelMap["max_snapshots"] = *model.MaxSnapshots
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesMaxSpaceConfigToMap(model *backuprecoveryv1.MaxSpaceConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MaxSnapshotSpacePercentage != nil {
		modelMap["max_snapshot_space_percentage"] = *model.MaxSnapshotSpacePercentage
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesThrottlingPolicyOverridesToMap(model *backuprecoveryv1.ThrottlingPolicyOverrides) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DatastoreID != nil {
		modelMap["datastore_id"] = flex.IntValue(model.DatastoreID)
	}
	if model.DatastoreName != nil {
		modelMap["datastore_name"] = *model.DatastoreName
	}
	if model.ThrottlingPolicy != nil {
		throttlingPolicyMap, err := DataSourceIbmBackupRecoveryProtectionSourcesThrottlingPolicyToMap(model.ThrottlingPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["throttling_policy"] = []map[string]interface{}{throttlingPolicyMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesRegisteredSourceVlanConfigToMap(model *backuprecoveryv1.RegisteredSourceVlanConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Vlan != nil {
		modelMap["vlan"] = *model.Vlan
	}
	if model.DisableVlan != nil {
		modelMap["disable_vlan"] = *model.DisableVlan
	}
	if model.InterfaceName != nil {
		modelMap["interface_name"] = *model.InterfaceName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesUniqueGlobalIDToMap(model *backuprecoveryv1.UniqueGlobalID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterID != nil {
		modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	}
	if model.ClusterIncarnationID != nil {
		modelMap["cluster_incarnation_id"] = flex.IntValue(model.ClusterIncarnationID)
	}
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesNetworkingInformationToMap(model *backuprecoveryv1.NetworkingInformation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ResourceVec != nil {
		resourceVec := []map[string]interface{}{}
		for _, resourceVecItem := range model.ResourceVec {
			resourceVecItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesClusterNetworkResourceInformationToMap(&resourceVecItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			resourceVec = append(resourceVec, resourceVecItemMap)
		}
		modelMap["resource_vec"] = resourceVec
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesClusterNetworkResourceInformationToMap(model *backuprecoveryv1.ClusterNetworkResourceInformation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Endpoints != nil {
		endpoints := []map[string]interface{}{}
		for _, endpointsItem := range model.Endpoints {
			endpointsItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesClusterNetworkingEndpointToMap(&endpointsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			endpoints = append(endpoints, endpointsItemMap)
		}
		modelMap["endpoints"] = endpoints
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesClusterNetworkingEndpointToMap(model *backuprecoveryv1.ClusterNetworkingEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Fqdn != nil {
		modelMap["fqdn"] = *model.Fqdn
	}
	if model.Ipv4Addr != nil {
		modelMap["ipv4_addr"] = *model.Ipv4Addr
	}
	if model.Ipv6Addr != nil {
		modelMap["ipv6_addr"] = *model.Ipv6Addr
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesPhysicalVolumeToMap(model *backuprecoveryv1.PhysicalVolume) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DevicePath != nil {
		modelMap["device_path"] = *model.DevicePath
	}
	if model.Guid != nil {
		modelMap["guid"] = *model.Guid
	}
	if model.IsBootVolume != nil {
		modelMap["is_boot_volume"] = *model.IsBootVolume
	}
	if model.IsExtendedAttributesSupported != nil {
		modelMap["is_extended_attributes_supported"] = *model.IsExtendedAttributesSupported
	}
	if model.IsProtected != nil {
		modelMap["is_protected"] = *model.IsProtected
	}
	if model.IsSharedVolume != nil {
		modelMap["is_shared_volume"] = *model.IsSharedVolume
	}
	if model.Label != nil {
		modelMap["label"] = *model.Label
	}
	if model.LogicalSizeBytes != nil {
		modelMap["logical_size_bytes"] = *model.LogicalSizeBytes
	}
	if model.MountPoints != nil {
		modelMap["mount_points"] = model.MountPoints
	}
	if model.MountType != nil {
		modelMap["mount_type"] = *model.MountType
	}
	if model.NetworkPath != nil {
		modelMap["network_path"] = *model.NetworkPath
	}
	if model.UsedSizeBytes != nil {
		modelMap["used_size_bytes"] = *model.UsedSizeBytes
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesVssWritersToMap(model *backuprecoveryv1.VssWriters) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsWriterExcluded != nil {
		modelMap["is_writer_excluded"] = *model.IsWriterExcluded
	}
	if model.WriterName != nil {
		modelMap["writer_name"] = *model.WriterName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesSqlProtectionSourceToMap(model *backuprecoveryv1.SqlProtectionSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsAvailableForVssBackup != nil {
		modelMap["is_available_for_vss_backup"] = *model.IsAvailableForVssBackup
	}
	if model.CreatedTimestamp != nil {
		modelMap["created_timestamp"] = *model.CreatedTimestamp
	}
	if model.DatabaseName != nil {
		modelMap["database_name"] = *model.DatabaseName
	}
	if model.DbAagEntityID != nil {
		modelMap["db_aag_entity_id"] = flex.IntValue(model.DbAagEntityID)
	}
	if model.DbAagName != nil {
		modelMap["db_aag_name"] = *model.DbAagName
	}
	if model.DbCompatibilityLevel != nil {
		modelMap["db_compatibility_level"] = flex.IntValue(model.DbCompatibilityLevel)
	}
	if model.DbFileGroups != nil {
		modelMap["db_file_groups"] = model.DbFileGroups
	}
	if model.DbFiles != nil {
		dbFiles := []map[string]interface{}{}
		for _, dbFilesItem := range model.DbFiles {
			dbFilesItemMap, err := DataSourceIbmBackupRecoveryProtectionSourcesDatabaseFileInformationToMap(&dbFilesItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			dbFiles = append(dbFiles, dbFilesItemMap)
		}
		modelMap["db_files"] = dbFiles
	}
	if model.DbOwnerUsername != nil {
		modelMap["db_owner_username"] = *model.DbOwnerUsername
	}
	if model.DefaultDatabaseLocation != nil {
		modelMap["default_database_location"] = *model.DefaultDatabaseLocation
	}
	if model.DefaultLogLocation != nil {
		modelMap["default_log_location"] = *model.DefaultLogLocation
	}
	if model.ID != nil {
		idMap, err := DataSourceIbmBackupRecoveryProtectionSourcesSQLSourceIDToMap(model.ID)
		if err != nil {
			return modelMap, err
		}
		modelMap["id"] = []map[string]interface{}{idMap}
	}
	if model.IsEncrypted != nil {
		modelMap["is_encrypted"] = *model.IsEncrypted
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.OwnerID != nil {
		modelMap["owner_id"] = flex.IntValue(model.OwnerID)
	}
	if model.RecoveryModel != nil {
		modelMap["recovery_model"] = *model.RecoveryModel
	}
	if model.SqlServerDbState != nil {
		modelMap["sql_server_db_state"] = *model.SqlServerDbState
	}
	if model.SqlServerInstanceVersion != nil {
		sqlServerInstanceVersionMap, err := DataSourceIbmBackupRecoveryProtectionSourcesSQLServerInstanceVersionToMap(model.SqlServerInstanceVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["sql_server_instance_version"] = []map[string]interface{}{sqlServerInstanceVersionMap}
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesDatabaseFileInformationToMap(model *backuprecoveryv1.DatabaseFileInformation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.FileType != nil {
		modelMap["file_type"] = *model.FileType
	}
	if model.FullPath != nil {
		modelMap["full_path"] = *model.FullPath
	}
	if model.SizeBytes != nil {
		modelMap["size_bytes"] = flex.IntValue(model.SizeBytes)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesSQLSourceIDToMap(model *backuprecoveryv1.SQLSourceID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreatedDateMsecs != nil {
		modelMap["created_date_msecs"] = flex.IntValue(model.CreatedDateMsecs)
	}
	if model.DatabaseID != nil {
		modelMap["database_id"] = flex.IntValue(model.DatabaseID)
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = *model.InstanceID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionSourcesSQLServerInstanceVersionToMap(model *backuprecoveryv1.SQLServerInstanceVersion) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Build != nil {
		modelMap["build"] = *model.Build
	}
	if model.MajorVersion != nil {
		modelMap["major_version"] = *model.MajorVersion
	}
	if model.MinorVersion != nil {
		modelMap["minor_version"] = *model.MinorVersion
	}
	if model.Revision != nil {
		modelMap["revision"] = *model.Revision
	}
	if model.VersionString != nil {
		modelMap["version_string"] = *model.VersionString
	}
	return modelMap, nil
}
