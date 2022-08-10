// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/scc-go-sdk/v3/posturemanagementv2"
)

func DataSourceIBMSccPostureScope() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSccPostureScopeRead,

		Schema: map[string]*schema.Schema{
			"scope_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id for the given API.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_name .",
			},
			"uuid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_uuid .Will be displayed only when value exists.",
			},
			"partner_uuid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of partner_uuid .Will be displayed only when value exists.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_description .Will be displayed only when value exists.",
			},
			"org_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Stores the value of scope_org_id .Will be displayed only when value exists.",
			},
			"cloud_type_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Stores the value of scope_cloud_type_id .Will be displayed only when value exists.",
			},
			"tld_credential_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Stores the value of scope_tld_credential_id .Will be displayed only when value exists.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_status .Will be displayed only when value exists.",
			},
			"status_msg": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_status_msg .Will be displayed only when value exists.",
			},
			"subset_selected": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Stores the value of scope_subset_selected .Will be displayed only when value exists.",
			},
			"enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Stores the value of scope_enabled .Will be displayed only when value exists.",
			},
			"last_discover_start_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_last_discover_start_time .Will be displayed only when value exists.",
			},
			"last_discover_completed_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_last_discover_completed_time .Will be displayed only when value exists.",
			},
			"last_successful_discover_start_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_last_successful_discover_start_time .Will be displayed only when value exists.",
			},
			"last_successful_discover_completed_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_last_successful_discover_completed_time .Will be displayed only when value exists.",
			},
			"task_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_task_type .Will be displayed only when value exists.",
			},
			"tasks": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Stores the value of scope_tasks .Will be displayed only when value exists.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_logs": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Stores the value of task_logs .",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"task_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of task_id .",
						},
						"task_gateway_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of task_gateway_id .",
						},
						"task_gateway_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_gateway_name .",
						},
						"task_task_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_task_type .",
						},
						"task_gateway_schema_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of task_gateway_schema_id .",
						},
						"task_schema_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_schema_name .",
						},
						"task_discover_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of task_discover_id .",
						},
						"task_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_status .",
						},
						"task_status_msg": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_status_msg .",
						},
						"task_start_time": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of task_start_time .",
						},
						"task_updated_time": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of task_updated_time .",
						},
						"task_derived_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_derived_status .",
						},
						"task_created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of task_created_by .",
						},
					},
				},
			},
			"status_updated_time": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_status_updated_time .Will be displayed only when value exists.",
			},
			"collectors_by_type": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Stores the value of collectors_by_type .Will be displayed only when value exists.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"credentials_by_type": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Stores the value of scope_credentials_by_type .Will be displayed only when value exists.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"credentials_by_sub_categeory_type": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Stores the value of scope_credentials_by_sub_categeory_type .Will be displayed only when value exists.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"sub_categories_by_type": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Stores the value of scope_sub_categories_by_type .Will be displayed only when value exists.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"resource_groups": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_resource_groups .Will be displayed only when value exists.",
			},
			"region_names": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_region_names .Will be displayed only when value exists.",
			},
			"cloud_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_cloud_type .Will be displayed only when value exists.",
			},
			"env_sub_category": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_env_sub_category .Will be displayed only when value exists.",
			},
			"tld_credentail": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Stores the value of ScopeDetailsCredential .",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"credential_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_id .",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_name .",
						},
						"uuid": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_uuid .",
						},
						"credential_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_type .",
						},
						"data": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of credential_data .",
						},
						"display_fields": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Details the fields on the credential. This will change as per credential type selected.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ibm_api_key": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IBM Cloud API Key. This is mandatory for IBM Credential Type.",
									},
									"aws_client_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS client Id.This is mandatory for AWS Cloud.",
									},
									"aws_client_secret": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS client secret.This is mandatory for AWS Cloud.",
									},
									"aws_region": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS region.",
									},
									"aws_arn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "AWS arn value.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "username of the user.This is mandatory for DataBase, Kerbros,OpenStack Credentials.",
									},
									"password": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "password of the user.This is mandatory for DataBase, Kerbros,OpenStack Credentials.",
									},
									"azure_client_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Azure client Id. This is mandatory for Azure Credential type.",
									},
									"azure_client_secret": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Azure client secret.This is mandatory for Azure Credential type.",
									},
									"azure_subscription_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Azure subscription Id.This is mandatory for Azure Credential type.",
									},
									"azure_resource_group": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Azure resource group.",
									},
									"database_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Database name.This is mandatory for Database Credential type.",
									},
									"winrm_authtype": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kerberos windows auth type.This is mandatory for Windows Kerberos Credential type.",
									},
									"winrm_usessl": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kerberos windows ssl.This is mandatory for Windows Kerberos Credential type.",
									},
									"winrm_port": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Kerberos windows port.This is mandatory for Windows Kerberos Credential type.",
									},
									"ms_365_client_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The MS365 client Id.This is mandatory for Windows MS365 Credential type.",
									},
									"ms_365_client_secret": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The MS365 client secret.This is mandatory for Windows MS365 Credential type.",
									},
									"ms_365_tenant_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The MS365 tenantId.This is mandatory for Windows MS365 Credential type.",
									},
									"auth_url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "auth url of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
									},
									"project_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Project name of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
									},
									"user_domain_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "user domain name of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
									},
									"project_domain_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "project domain name of the Open Stack cloud.This is mandatory for Open Stack Credential type.",
									},
									"pem_file_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the PEM file.",
									},
									"pem_data": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The base64 encoded data to associate with the PEM file.",
									},
								},
							},
						},
						"version_timestamp": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of credential_version_timestamp .",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_description .",
						},
						"is_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Stores the value of credential_is_enabled .",
						},
						"gateway_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_gateway_key .",
						},
						"credential_group": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of credential_credential_group .",
						},
						"enabled_credential_group": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Stores the value of credential_enabled_credential_group .",
						},
						"groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Stores the value of credential_groups .",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"credential_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "credential group id.",
									},
									"passphrase": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "passphase of the credential.",
									},
								},
							},
						},
						"purpose": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of credential_purpose .",
						},
					},
				},
			},
			"collectors": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Stores the value of collectors .Will be displayed only when value exists.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"collector_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the collector.",
						},
						"display_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-friendly name of the collector.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the collector.",
						},
						"public_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public key of the collector.Will be used for ssl communciation between collector and orchestrator .This will be populated when collector is installed.",
						},
						"last_heartbeat": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the heartbeat time of a controller . This value exists when collector is installed and running.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of collector.",
						},
						"collector_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector version. This field is populated when collector is installed.",
						},
						"image_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The image version of the collector. This field is populated when collector is installed. \".",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the collector.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the user that created the collector.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO Date/Time the collector was created.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the user that modified the collector.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO Date/Time the collector was modified.",
						},
						"enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Identifies whether the collector is enabled or not(deleted).",
						},
						"registration_code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The registration code of the collector.This is will be used for initial authentication during installation of collector.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the collector.",
						},
						"credential_public_key": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The credential public key.",
						},
						"failure_count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of times the collector has failed.",
						},
						"approved_local_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The approved local gateway ip of the collector. This field will be populated only when collector is installed.",
						},
						"approved_internet_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The approved internet gateway ip of the collector. This field will be populated only when collector is installed.",
						},
						"last_failed_local_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failed local gateway ip. This field will be populated only when collector is installed.",
						},
						"reset_reason": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason for the collector reset .User resets the collector with a reason for reset. The reason entered by the user is saved in this field .",
						},
						"hostname": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector host name. This field will be populated when collector is installed.This will have fully qualified domain name.",
						},
						"install_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The installation path of the collector. This field will be populated when collector is installed.The value will be folder path.",
						},
						"use_private_endpoint": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the collector should use a public or private endpoint. This value is generated based on is_public field value during collector creation. If is_public is set to true, this value will be false.",
						},
						"managed_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The entity that manages the collector.",
						},
						"trial_expiry": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trial expiry. This holds the expiry date of registration_code. This field will be populated when collector is installed.",
						},
						"last_failed_internet_gateway_ip": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failed internet gateway ip of the collector.",
						},
						"status_description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The collector status.",
						},
						"reset_time": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ISO Date/Time of the collector reset. This value will be populated when a collector is reset. The data-time when the reset event is occured is captured in this field.",
						},
						"is_public": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether the collector endpoint is accessible on a public network.If set to `true`, the collector connects to resources in your account over a public network. If set to `false`, the collector connects to resources by using a private IP that is accessible only through the IBM Cloud private network.",
						},
						"is_ubi_image": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Determines whether the collector has a Ubi image.",
						},
					},
				},
			},
			"first_level_scoped_data": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Stores the value of scope_first_level_scoped_data .Will be displayed only when value exists.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scope_object": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of  scope_object .",
						},
						"scope_init_scope": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_init_scope .",
						},
						"scope": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope .",
						},
						"scope_changed": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Stores the value of  scope_changed .",
						},
						"scope_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_id .",
						},
						"scope_properties": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of  scope_properties .",
						},
						"scope_overlay": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_overlay .",
						},
						"scope_new_found": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Stores the value of scope_new_found .",
						},
						"scope_discovery_status": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of scope_discovery_status .",
						},
						"scope_fact_status": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of scope_fact_status .",
						},
						"scope_facts": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_facts .",
						},
						"scope_list_members": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of scope_list_members .",
						},
						"scope_children": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of scope_children .",
						},
						"scope_resource_category": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_resource_category .",
						},
						"scope_resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_resource_type .",
						},
						"scope_resource": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_resource .",
						},
						"scope_resource_attributes": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of scope_resource_attributes .",
						},
						"scope_drift": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of  scope_drift .",
						},
						"scope_parse_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Stores the value of scope_parse_status .",
						},
						"scope_transformed_facts": &schema.Schema{
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Stores the value of scope_transformed_facts .",
						},
						"scope_collector_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Stores the value of scope_collector_id .",
						},
					},
				},
			},
			"discovery_methods": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Stores the value of scope_discovery_methods .Will be displayed only when value exists.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"discovery_method": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_discovery_method .Will be displayed only when value exists.",
			},
			"file_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_file_type .Will be displayed only when value exists.",
			},
			"file_format": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_file_format .Will be displayed only when value exists.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_created_by .Will be displayed only when value exists.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_created_on .Will be displayed only when value exists.",
			},
			"modified_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_modified_by .Will be displayed only when value exists.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_modified_on .Will be displayed only when value exists.",
			},
			"is_discovery_scheduled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Stores the value of scope_is_discovery_scheduled .Will be displayed only when value exists.",
			},
			"interval": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Stores the value of scope_freq .Will be displayed only when value exists.",
			},
			"discovery_setting_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Stores the value of scope_discovery_setting_id .Will be displayed only when value exists.",
			},
			"include_new_eagerly": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Stores the value of scope_include_new_eagerly .Will be displayed only when value exists.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_type .Will be displayed only when value exists.",
			},
			"correlation_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A correlation_Id is created when a scope is created and discovery task is triggered or when a validation is triggered on a Scope. This is used to get the status of the task(discovery or validation).",
			},
			"credential_attributes": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Stores the value of scope_credential_attributes .Will be displayed only when value exists.",
			},
		},
	}
}

func dataSourceIBMSccPostureScopeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	postureManagementClient, err := meta.(conns.ClientSession).PostureManagementV2()
	if err != nil {
		return diag.FromErr(err)
	}

	ScopeDetailsOptions := &posturemanagementv2.GetScopeDetailsOptions{}
	userDetails, err := meta.(conns.ClientSession).BluemixUserDetails()
	if err != nil {
		return diag.FromErr(fmt.Errorf("Error getting userDetails %s", err))
	}

	accountID := userDetails.UserAccount
	ScopeDetailsOptions.SetAccountID(accountID)
	ScopeDetailsOptions.SetID(d.Get("scope_id").(string))

	scope, response, err := postureManagementClient.GetScopeDetailsWithContext(context, ScopeDetailsOptions)
	if err != nil {
		log.Printf("[DEBUG] GetScopeDetailsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetScopeDetailsWithContext failed %s\n%s", err, response))
	}

	d.SetId(*scope.ID)
	if err = d.Set("name", scope.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("uuid", scope.UUID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting uuid: %s", err))
	}
	if err = d.Set("partner_uuid", scope.PartnerUUID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting partner_uuid: %s", err))
	}
	if err = d.Set("description", scope.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}
	if err = d.Set("org_id", flex.IntValue(scope.OrgID)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting org_id: %s", err))
	}
	if err = d.Set("cloud_type_id", flex.IntValue(scope.CloudTypeID)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cloud_type_id: %s", err))
	}
	if err = d.Set("tld_credential_id", flex.IntValue(scope.TldCredentialID)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tld_credential_id: %s", err))
	}
	if err = d.Set("status", scope.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}
	if err = d.Set("status_msg", scope.StatusMsg); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_msg: %s", err))
	}
	if err = d.Set("subset_selected", scope.SubsetSelected); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting subset_selected: %s", err))
	}
	if err = d.Set("enabled", scope.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting enabled: %s", err))
	}
	if err = d.Set("last_discover_start_time", scope.LastDiscoverStartTime); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_discover_start_time: %s", err))
	}
	if err = d.Set("last_discover_completed_time", scope.LastDiscoverCompletedTime); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_discover_completed_time: %s", err))
	}
	if err = d.Set("last_successful_discover_start_time", scope.LastSuccessfulDiscoverStartTime); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_successful_discover_start_time: %s", err))
	}
	if err = d.Set("last_successful_discover_completed_time", scope.LastSuccessfulDiscoverCompletedTime); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting last_successful_discover_completed_time: %s", err))
	}
	if err = d.Set("task_type", scope.TaskType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting task_type: %s", err))
	}

	if scope.Tasks != nil {
		err = d.Set("tasks", dataSourceScopeFlattenTasks(scope.Tasks))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting tasks %s", err))
		}
	}
	if err = d.Set("status_updated_time", scope.StatusUpdatedTime); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status_updated_time: %s", err))
	}

	if scope.CollectorsByType != nil {
		convertedMap := make(map[string]interface{}, len(scope.CollectorsByType))
		for k, v := range scope.CollectorsByType {
			convertedMap[k] = v
		}

		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting collectors_by_type %s", err))
		}
	}

	if scope.CredentialsByType != nil {
		convertedMap := make(map[string]interface{}, len(scope.CredentialsByType))
		for k, v := range scope.CredentialsByType {
			convertedMap[k] = v
		}

		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting credentials_by_type %s", err))
		}
	}

	if scope.CredentialsBySubCategeoryType != nil {
		convertedMap := make(map[string]interface{}, len(scope.CredentialsBySubCategeoryType))
		for k, v := range scope.CredentialsBySubCategeoryType {
			convertedMap[k] = v
		}

		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting credentials_by_sub_categeory_type %s", err))
		}
	}
	if err = d.Set("resource_groups", scope.ResourceGroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_groups: %s", err))
	}
	if err = d.Set("region_names", scope.RegionNames); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting region_names: %s", err))
	}
	if err = d.Set("cloud_type", scope.CloudType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cloud_type: %s", err))
	}
	if err = d.Set("env_sub_category", scope.EnvSubCategory); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting env_sub_category: %s", err))
	}

	if scope.TldCredentail != nil {
		err = d.Set("tld_credentail", dataSourceScopeFlattenTldCredentail(*scope.TldCredentail))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting tld_credentail %s", err))
		}
	}

	if scope.Collectors != nil {
		err = d.Set("collectors", dataSourceScopeFlattenCollectors(scope.Collectors))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting collectors %s", err))
		}
	}

	if scope.FirstLevelScopedData != nil {
		err = d.Set("first_level_scoped_data", dataSourceScopeFlattenFirstLevelScopedData(scope.FirstLevelScopedData))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting first_level_scoped_data %s", err))
		}
	}
	if err = d.Set("discovery_method", scope.DiscoveryMethod); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting discovery_method: %s", err))
	}
	if err = d.Set("file_type", scope.FileType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting file_type: %s", err))
	}
	if err = d.Set("file_format", scope.FileFormat); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting file_format: %s", err))
	}
	if err = d.Set("created_by", scope.CreatedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_by: %s", err))
	}
	if err = d.Set("created_at", scope.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("modified_by", scope.ModifiedBy); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_by: %s", err))
	}
	if err = d.Set("modified_at", scope.ModifiedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}
	if err = d.Set("is_discovery_scheduled", scope.IsDiscoveryScheduled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_discovery_scheduled: %s", err))
	}
	if err = d.Set("interval", flex.IntValue(scope.Interval)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting interval: %s", err))
	}
	if err = d.Set("discovery_setting_id", flex.IntValue(scope.DiscoverySettingID)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting discovery_setting_id: %s", err))
	}
	if err = d.Set("include_new_eagerly", scope.IncludeNewEagerly); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting include_new_eagerly: %s", err))
	}
	if err = d.Set("type", scope.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if err = d.Set("correlation_id", scope.CorrelationID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting correlation_id: %s", err))
	}
	if err = d.Set("credential_attributes", scope.CredentialAttributes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting credential_attributes: %s", err))
	}

	return nil
}

func dataSourceScopeFlattenTasks(result []posturemanagementv2.ScopeDetailsGatewayTask) (tasks []map[string]interface{}) {
	for _, tasksItem := range result {
		tasks = append(tasks, dataSourceScopeTasksToMap(tasksItem))
	}

	return tasks
}

func dataSourceScopeTasksToMap(tasksItem posturemanagementv2.ScopeDetailsGatewayTask) (tasksMap map[string]interface{}) {
	tasksMap = map[string]interface{}{}

	if tasksItem.TaskLogs != nil {
		taskLogsList := []map[string]interface{}{}
		for _, taskLogsItem := range tasksItem.TaskLogs {
			taskLogsList = append(taskLogsList, dataSourceScopeTasksTaskLogsToMap(taskLogsItem))
		}
		tasksMap["task_logs"] = taskLogsList
	}
	if tasksItem.TaskID != nil {
		tasksMap["task_id"] = tasksItem.TaskID
	}
	if tasksItem.TaskGatewayID != nil {
		tasksMap["task_gateway_id"] = tasksItem.TaskGatewayID
	}
	if tasksItem.TaskGatewayName != nil {
		tasksMap["task_gateway_name"] = tasksItem.TaskGatewayName
	}
	if tasksItem.TaskTaskType != nil {
		tasksMap["task_task_type"] = tasksItem.TaskTaskType
	}
	if tasksItem.TaskGatewaySchemaID != nil {
		tasksMap["task_gateway_schema_id"] = tasksItem.TaskGatewaySchemaID
	}
	if tasksItem.TaskSchemaName != nil {
		tasksMap["task_schema_name"] = tasksItem.TaskSchemaName
	}
	if tasksItem.TaskDiscoverID != nil {
		tasksMap["task_discover_id"] = tasksItem.TaskDiscoverID
	}
	if tasksItem.TaskStatus != nil {
		tasksMap["task_status"] = tasksItem.TaskStatus
	}
	if tasksItem.TaskStatusMsg != nil {
		tasksMap["task_status_msg"] = tasksItem.TaskStatusMsg
	}
	if tasksItem.TaskStartTime != nil {
		tasksMap["task_start_time"] = tasksItem.TaskStartTime
	}
	if tasksItem.TaskUpdatedTime != nil {
		tasksMap["task_updated_time"] = tasksItem.TaskUpdatedTime
	}
	if tasksItem.TaskDerivedStatus != nil {
		tasksMap["task_derived_status"] = tasksItem.TaskDerivedStatus
	}
	if tasksItem.TaskCreatedBy != nil {
		tasksMap["task_created_by"] = tasksItem.TaskCreatedBy
	}

	return tasksMap
}

func dataSourceScopeTasksTaskLogsToMap(taskLogsItem posturemanagementv2.TaskLogs) (taskLogsMap map[string]interface{}) {
	taskLogsMap = map[string]interface{}{}

	return taskLogsMap
}

func dataSourceScopeFlattenTldCredentail(result posturemanagementv2.ScopeDetailsCredential) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceScopeTldCredentailToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceScopeTldCredentailToMap(tldCredentailItem posturemanagementv2.ScopeDetailsCredential) (tldCredentailMap map[string]interface{}) {
	tldCredentailMap = map[string]interface{}{}

	if tldCredentailItem.ID != nil {
		tldCredentailMap["credential_id"] = tldCredentailItem.ID
	}
	if tldCredentailItem.Name != nil {
		tldCredentailMap["name"] = tldCredentailItem.Name
	}
	if tldCredentailItem.UUID != nil {
		tldCredentailMap["uuid"] = tldCredentailItem.UUID
	}
	if tldCredentailItem.Type != nil {
		tldCredentailMap["credential_type"] = tldCredentailItem.Type
	}
	if tldCredentailItem.Data != nil {
		tldCredentailMap["data"] = tldCredentailItem.Data
	}
	if tldCredentailItem.DisplayFields != nil {
		displayFieldsList := []map[string]interface{}{}
		displayFieldsMap := dataSourceScopeTldCredentailDisplayFieldsToMap(*tldCredentailItem.DisplayFields)
		displayFieldsList = append(displayFieldsList, displayFieldsMap)
		tldCredentailMap["display_fields"] = displayFieldsList
	}
	if tldCredentailItem.VersionTimestamp != nil {
		tldCredentailMap["version_timestamp"] = tldCredentailItem.VersionTimestamp
	}
	if tldCredentailItem.Description != nil {
		tldCredentailMap["description"] = tldCredentailItem.Description
	}
	if tldCredentailItem.IsEnabled != nil {
		tldCredentailMap["is_enabled"] = tldCredentailItem.IsEnabled
	}
	if tldCredentailItem.GatewayKey != nil {
		tldCredentailMap["gateway_key"] = tldCredentailItem.GatewayKey
	}
	if tldCredentailItem.CredentialGroup != nil {
		tldCredentailMap["credential_group"] = tldCredentailItem.CredentialGroup
	}
	if tldCredentailItem.EnabledCredentialGroup != nil {
		tldCredentailMap["enabled_credential_group"] = tldCredentailItem.EnabledCredentialGroup
	}
	if tldCredentailItem.Groups != nil {
		groupsList := []map[string]interface{}{}
		for _, groupsItem := range tldCredentailItem.Groups {
			groupsList = append(groupsList, dataSourceScopeTldCredentailGroupsToMap(groupsItem))
		}
		tldCredentailMap["groups"] = groupsList
	}
	if tldCredentailItem.Purpose != nil {
		tldCredentailMap["purpose"] = tldCredentailItem.Purpose
	}

	return tldCredentailMap
}

func dataSourceScopeTldCredentailDisplayFieldsToMap(displayFieldsItem posturemanagementv2.ScopeDetailsCredentialDisplayFields) (displayFieldsMap map[string]interface{}) {
	displayFieldsMap = map[string]interface{}{}

	if displayFieldsItem.IBMAPIKey != nil {
		displayFieldsMap["ibm_api_key"] = displayFieldsItem.IBMAPIKey
	}
	if displayFieldsItem.AwsClientID != nil {
		displayFieldsMap["aws_client_id"] = displayFieldsItem.AwsClientID
	}
	if displayFieldsItem.AwsClientSecret != nil {
		displayFieldsMap["aws_client_secret"] = displayFieldsItem.AwsClientSecret
	}
	if displayFieldsItem.AwsRegion != nil {
		displayFieldsMap["aws_region"] = displayFieldsItem.AwsRegion
	}
	if displayFieldsItem.AwsArn != nil {
		displayFieldsMap["aws_arn"] = displayFieldsItem.AwsArn
	}
	if displayFieldsItem.Username != nil {
		displayFieldsMap["username"] = displayFieldsItem.Username
	}
	if displayFieldsItem.Password != nil {
		displayFieldsMap["password"] = displayFieldsItem.Password
	}
	if displayFieldsItem.AzureClientID != nil {
		displayFieldsMap["azure_client_id"] = displayFieldsItem.AzureClientID
	}
	if displayFieldsItem.AzureClientSecret != nil {
		displayFieldsMap["azure_client_secret"] = displayFieldsItem.AzureClientSecret
	}
	if displayFieldsItem.AzureSubscriptionID != nil {
		displayFieldsMap["azure_subscription_id"] = displayFieldsItem.AzureSubscriptionID
	}
	if displayFieldsItem.AzureResourceGroup != nil {
		displayFieldsMap["azure_resource_group"] = displayFieldsItem.AzureResourceGroup
	}
	if displayFieldsItem.DatabaseName != nil {
		displayFieldsMap["database_name"] = displayFieldsItem.DatabaseName
	}
	if displayFieldsItem.WinrmAuthtype != nil {
		displayFieldsMap["winrm_authtype"] = displayFieldsItem.WinrmAuthtype
	}
	if displayFieldsItem.WinrmUsessl != nil {
		displayFieldsMap["winrm_usessl"] = displayFieldsItem.WinrmUsessl
	}
	if displayFieldsItem.WinrmPort != nil {
		displayFieldsMap["winrm_port"] = displayFieldsItem.WinrmPort
	}
	if displayFieldsItem.Ms365ClientID != nil {
		displayFieldsMap["ms_365_client_id"] = displayFieldsItem.Ms365ClientID
	}
	if displayFieldsItem.Ms365ClientSecret != nil {
		displayFieldsMap["ms_365_client_secret"] = displayFieldsItem.Ms365ClientSecret
	}
	if displayFieldsItem.Ms365TenantID != nil {
		displayFieldsMap["ms_365_tenant_id"] = displayFieldsItem.Ms365TenantID
	}
	if displayFieldsItem.AuthURL != nil {
		displayFieldsMap["auth_url"] = displayFieldsItem.AuthURL
	}
	if displayFieldsItem.ProjectName != nil {
		displayFieldsMap["project_name"] = displayFieldsItem.ProjectName
	}
	if displayFieldsItem.UserDomainName != nil {
		displayFieldsMap["user_domain_name"] = displayFieldsItem.UserDomainName
	}
	if displayFieldsItem.ProjectDomainName != nil {
		displayFieldsMap["project_domain_name"] = displayFieldsItem.ProjectDomainName
	}

	return displayFieldsMap
}

func dataSourceScopeTldCredentailGroupsToMap(groupsItem posturemanagementv2.CredentialGroup) (groupsMap map[string]interface{}) {
	groupsMap = map[string]interface{}{}

	if groupsItem.ID != nil {
		groupsMap["credential_group_id"] = groupsItem.ID
	}
	if groupsItem.Passphrase != nil {
		groupsMap["passphrase"] = groupsItem.Passphrase
	}

	return groupsMap
}

func dataSourceScopeFlattenCollectors(result []posturemanagementv2.Collector) (collectors []map[string]interface{}) {
	for _, collectorsItem := range result {
		collectors = append(collectors, dataSourceScopeCollectorsToMap(collectorsItem))
	}

	return collectors
}

func dataSourceScopeCollectorsToMap(collectorsItem posturemanagementv2.Collector) (collectorsMap map[string]interface{}) {
	collectorsMap = map[string]interface{}{}

	if collectorsItem.ID != nil {
		collectorsMap["collector_id"] = collectorsItem.ID
	}
	if collectorsItem.DisplayName != nil {
		collectorsMap["display_name"] = collectorsItem.DisplayName
	}
	if collectorsItem.Name != nil {
		collectorsMap["name"] = collectorsItem.Name
	}
	if collectorsItem.PublicKey != nil {
		collectorsMap["public_key"] = collectorsItem.PublicKey
	}
	if collectorsItem.LastHeartbeat != nil {
		collectorsMap["last_heartbeat"] = collectorsItem.LastHeartbeat.String()
	}
	if collectorsItem.Status != nil {
		collectorsMap["status"] = collectorsItem.Status
	}
	if collectorsItem.CollectorVersion != nil {
		collectorsMap["collector_version"] = collectorsItem.CollectorVersion
	}
	if collectorsItem.ImageVersion != nil {
		collectorsMap["image_version"] = collectorsItem.ImageVersion
	}
	if collectorsItem.Description != nil {
		collectorsMap["description"] = collectorsItem.Description
	}
	if collectorsItem.CreatedBy != nil {
		collectorsMap["created_by"] = collectorsItem.CreatedBy
	}
	if collectorsItem.CreatedAt != nil {
		collectorsMap["created_at"] = collectorsItem.CreatedAt.String()
	}
	if collectorsItem.UpdatedBy != nil {
		collectorsMap["updated_by"] = collectorsItem.UpdatedBy
	}
	if collectorsItem.UpdatedAt != nil {
		collectorsMap["updated_at"] = collectorsItem.UpdatedAt.String()
	}
	if collectorsItem.Enabled != nil {
		collectorsMap["enabled"] = collectorsItem.Enabled
	}
	if collectorsItem.RegistrationCode != nil {
		collectorsMap["registration_code"] = collectorsItem.RegistrationCode
	}
	if collectorsItem.Type != nil {
		collectorsMap["type"] = collectorsItem.Type
	}
	if collectorsItem.CredentialPublicKey != nil {
		collectorsMap["credential_public_key"] = collectorsItem.CredentialPublicKey
	}
	if collectorsItem.FailureCount != nil {
		collectorsMap["failure_count"] = collectorsItem.FailureCount
	}
	if collectorsItem.ApprovedLocalGatewayIP != nil {
		collectorsMap["approved_local_gateway_ip"] = collectorsItem.ApprovedLocalGatewayIP
	}
	if collectorsItem.ApprovedInternetGatewayIP != nil {
		collectorsMap["approved_internet_gateway_ip"] = collectorsItem.ApprovedInternetGatewayIP
	}
	if collectorsItem.LastFailedLocalGatewayIP != nil {
		collectorsMap["last_failed_local_gateway_ip"] = collectorsItem.LastFailedLocalGatewayIP
	}
	if collectorsItem.ResetReason != nil {
		collectorsMap["reset_reason"] = collectorsItem.ResetReason
	}
	if collectorsItem.Hostname != nil {
		collectorsMap["hostname"] = collectorsItem.Hostname
	}
	if collectorsItem.InstallPath != nil {
		collectorsMap["install_path"] = collectorsItem.InstallPath
	}
	if collectorsItem.UsePrivateEndpoint != nil {
		collectorsMap["use_private_endpoint"] = collectorsItem.UsePrivateEndpoint
	}
	if collectorsItem.ManagedBy != nil {
		collectorsMap["managed_by"] = collectorsItem.ManagedBy
	}
	if collectorsItem.TrialExpiry != nil {
		collectorsMap["trial_expiry"] = collectorsItem.TrialExpiry.String()
	}
	if collectorsItem.LastFailedInternetGatewayIP != nil {
		collectorsMap["last_failed_internet_gateway_ip"] = collectorsItem.LastFailedInternetGatewayIP
	}
	if collectorsItem.StatusDescription != nil {
		collectorsMap["status_description"] = collectorsItem.StatusDescription
	}
	if collectorsItem.ResetTime != nil {
		collectorsMap["reset_time"] = collectorsItem.ResetTime.String()
	}
	if collectorsItem.IsPublic != nil {
		collectorsMap["is_public"] = collectorsItem.IsPublic
	}
	if collectorsItem.IsUbiImage != nil {
		collectorsMap["is_ubi_image"] = collectorsItem.IsUbiImage
	}

	return collectorsMap
}

func dataSourceScopeFlattenFirstLevelScopedData(result []posturemanagementv2.ScopeDetailsAssetData) (firstLevelScopedData []map[string]interface{}) {
	for _, firstLevelScopedDataItem := range result {
		firstLevelScopedData = append(firstLevelScopedData, dataSourceScopeFirstLevelScopedDataToMap(firstLevelScopedDataItem))
	}

	return firstLevelScopedData
}

func dataSourceScopeFirstLevelScopedDataToMap(firstLevelScopedDataItem posturemanagementv2.ScopeDetailsAssetData) (firstLevelScopedDataMap map[string]interface{}) {
	firstLevelScopedDataMap = map[string]interface{}{}

	if firstLevelScopedDataItem.ScopeObject != nil {
		firstLevelScopedDataMap["scope_object"] = firstLevelScopedDataItem.ScopeObject
	}
	if firstLevelScopedDataItem.ScopeInitScope != nil {
		firstLevelScopedDataMap["scope_init_scope"] = firstLevelScopedDataItem.ScopeInitScope
	}
	if firstLevelScopedDataItem.Scope != nil {
		firstLevelScopedDataMap["scope"] = firstLevelScopedDataItem.Scope
	}
	if firstLevelScopedDataItem.ScopeChanged != nil {
		firstLevelScopedDataMap["scope_changed"] = firstLevelScopedDataItem.ScopeChanged
	}
	if firstLevelScopedDataItem.ScopeID != nil {
		firstLevelScopedDataMap["scope_id"] = firstLevelScopedDataItem.ScopeID
	}
	if firstLevelScopedDataItem.ScopeProperties != nil {
		firstLevelScopedDataMap["scope_properties"] = firstLevelScopedDataItem.ScopeProperties
	}
	if firstLevelScopedDataItem.ScopeOverlay != nil {
		firstLevelScopedDataMap["scope_overlay"] = firstLevelScopedDataItem.ScopeOverlay
	}
	if firstLevelScopedDataItem.ScopeNewFound != nil {
		firstLevelScopedDataMap["scope_new_found"] = firstLevelScopedDataItem.ScopeNewFound
	}
	if firstLevelScopedDataItem.ScopeDiscoveryStatus != nil {
		firstLevelScopedDataMap["scope_discovery_status"] = firstLevelScopedDataItem.ScopeDiscoveryStatus
	}
	if firstLevelScopedDataItem.ScopeFactStatus != nil {
		firstLevelScopedDataMap["scope_fact_status"] = firstLevelScopedDataItem.ScopeFactStatus
	}
	if firstLevelScopedDataItem.ScopeFacts != nil {
		firstLevelScopedDataMap["scope_facts"] = firstLevelScopedDataItem.ScopeFacts
	}
	if firstLevelScopedDataItem.ScopeListMembers != nil {
		firstLevelScopedDataMap["scope_list_members"] = firstLevelScopedDataItem.ScopeListMembers
	}
	if firstLevelScopedDataItem.ScopeChildren != nil {
		firstLevelScopedDataMap["scope_children"] = firstLevelScopedDataItem.ScopeChildren
	}
	if firstLevelScopedDataItem.ScopeResourceCategory != nil {
		firstLevelScopedDataMap["scope_resource_category"] = firstLevelScopedDataItem.ScopeResourceCategory
	}
	if firstLevelScopedDataItem.ScopeResourceType != nil {
		firstLevelScopedDataMap["scope_resource_type"] = firstLevelScopedDataItem.ScopeResourceType
	}
	if firstLevelScopedDataItem.ScopeResource != nil {
		firstLevelScopedDataMap["scope_resource"] = firstLevelScopedDataItem.ScopeResource
	}
	if firstLevelScopedDataItem.ScopeResourceAttributes != nil {
		firstLevelScopedDataMap["scope_resource_attributes"] = firstLevelScopedDataItem.ScopeResourceAttributes
	}
	if firstLevelScopedDataItem.ScopeDrift != nil {
		firstLevelScopedDataMap["scope_drift"] = firstLevelScopedDataItem.ScopeDrift
	}
	if firstLevelScopedDataItem.ScopeParseStatus != nil {
		firstLevelScopedDataMap["scope_parse_status"] = firstLevelScopedDataItem.ScopeParseStatus
	}
	if firstLevelScopedDataItem.ScopeTransformedFacts != nil {
		firstLevelScopedDataMap["scope_transformed_facts"] = firstLevelScopedDataItem.ScopeTransformedFacts
	}
	if firstLevelScopedDataItem.ScopeCollectorID != nil {
		firstLevelScopedDataMap["scope_collector_id"] = firstLevelScopedDataItem.ScopeCollectorID
	}

	return firstLevelScopedDataMap
}
