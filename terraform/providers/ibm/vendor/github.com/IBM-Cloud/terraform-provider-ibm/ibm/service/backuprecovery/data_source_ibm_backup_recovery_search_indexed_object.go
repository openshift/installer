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

func DataSourceIbmBackupRecoverySearchIndexedObject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoverySearchIndexedObjectRead,
		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"protection_group_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies a list of Protection Group ids to filter the indexed objects. If specified, the objects indexed by specified Protection Group ids will be returned.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"storage_domain_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the Storage Domain ids to filter indexed objects for which Protection Groups are writing data to Cohesity Views on the specified Storage Domains.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "TenantId contains id of the tenant for which objects are to be returned.",
			},
			"include_tenants": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If true, the response will include objects which belongs to all tenants which the current user has permission to see. Default value is false.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "\"This field is deprecated. Please use mightHaveTagIds.\".",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"snapshot_tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "\"This field is deprecated. Please use mightHaveSnapshotTagIds.\".",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"must_have_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies tags which must be all present in the document.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"might_have_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies list of tags, one or more of which might be present in the document. These are OR'ed together and the resulting criteria AND'ed with the rest of the query.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"must_have_snapshot_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies snapshot tags which must be all present in the document.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"might_have_snapshot_tag_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies list of snapshot tags, one or more of which might be present in the document. These are OR'ed together and the resulting criteria AND'ed with the rest of the query.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"pagination_cookie": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the pagination cookie with which subsequent parts of the response can be fetched.",
			},
			"object_count": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the number of indexed objects to be fetched for the specified pagination cookie.",
			},
			"object_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				// ForceNew:     true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_search_indexed_object", "object_type"),
				Description: "Specifies the object type to be searched for.",
			},
			"use_cached_data": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies whether we can serve the GET request from the read replica cache. There is a lag of 15 seconds between the read replica and primary data source.",
			},
			"cassandra_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search Cassandra on a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cassandra_object_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies one or more Cassandra object types to be searched.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the Cassandra Objects.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"couchbase_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search CouchBase on a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"couchbase_object_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies Couchbase object types be searched. For Couchbase it can only be set to 'CouchbaseBuckets'.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the Couchbase Objects.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"email_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request parameters to search for emails and email folders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attendees_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the calendar items which have specified email addresses as attendees.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"bcc_recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the emails which are sent to specified email addresses in BCC.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"cc_recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the emails which are sent to specified email addresses in CC.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"created_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds where the created time of the email/item is less than specified value.",
						},
						"created_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds where the created time of the email/item is more than specified value.",
						},
						"due_date_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds where the last modification time of the email/item is less than specified value.",
						},
						"due_date_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds where the last modification time of the email/item is more than specified value.",
						},
						"email_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the contact items which have specified text in email address.",
						},
						"email_subject": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the emails which have the specified text in its subject.",
						},
						"first_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the contacts with specified text in first name.",
						},
						"folder_names": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the emails which are categorized to specified folders.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"has_attachment": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Filters the emails which have attachment.",
						},
						"last_modified_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds where the last modification time of the email/item is less than specified value.",
						},
						"last_modified_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds where the last modification time of the email/item is more than specified value.",
						},
						"last_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the contacts with specified text in last name.",
						},
						"middle_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the contacts with specified text in middle name.",
						},
						"organizer_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the calendar items which are organized by specified User's email address.",
						},
						"received_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds where the received time of the email is less than specified value.",
						},
						"received_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds where the received time of the email is more than specified value.",
						},
						"recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the emails which are sent to specified email addresses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"sender_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the emails which are received from specified User's email address.",
						},
						"source_environment": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the source environment.",
						},
						"task_status_types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of task item status types. Task items having status within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of mailbox item types. Only items within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"o365_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies email search request params specific to O365 environment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the domain Ids in which mailboxes are registered.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"mailbox_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the mailbox Ids which contains the emails/folders.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
					},
				},
			},
			"exchange_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the parameters which are specific for searching Exchange mailboxes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the Exchange Objects.",
						},
					},
				},
			},
			"file_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request parameters to search for files and file folders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the search string to filter the files. User can specify a wildcard character '*' as a suffix to a string where all files name are matched with the prefix string.",
						},
						"types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of file types. Only files within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"source_environments": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of the source environments. Only files from these types of source will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
						"object_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of object ids. Only files found in these objects will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"hbase_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search Hbase on a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hbase_object_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies one or more Hbase object types be searched.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the Hbase Objects.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"hdfs_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search HDFS on a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hdfs_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies types as Folders or Files or both to be searched.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the HDFS Folders and Files.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"hive_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search Hive on a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hive_object_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies one or more Hive object types be searched.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the Hive Objects.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"mongodb_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search Mongo DB on a cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mongo_db_object_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies one or more MongoDB object types be searched.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the MongoDB Objects.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"ms_groups_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request params to search for Groups items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mailbox_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the request parameters to search for mailbox items and folders.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"attendees_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Filters the calendar items which have specified email addresses as attendees.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"bcc_recipient_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Filters the emails which are sent to specified email addresses in BCC.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"cc_recipient_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Filters the emails which are sent to specified email addresses in CC.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"created_end_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the end time in Unix timestamp epoch in seconds where the created time of the email/item is less than specified value.",
									},
									"created_start_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the start time in Unix timestamp epoch in seconds where the created time of the email/item is more than specified value.",
									},
									"due_date_end_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the end time in Unix timestamp epoch in seconds where the last modification time of the email/item is less than specified value.",
									},
									"due_date_start_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the start time in Unix timestamp epoch in seconds where the last modification time of the email/item is more than specified value.",
									},
									"email_address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the contact items which have specified text in email address.",
									},
									"email_subject": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the emails which have the specified text in its subject.",
									},
									"first_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the contacts with specified text in first name.",
									},
									"folder_names": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Filters the emails which are categorized to specified folders.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"has_attachment": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Filters the emails which have attachment.",
									},
									"last_modified_end_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the end time in Unix timestamp epoch in seconds where the last modification time of the email/item is less than specified value.",
									},
									"last_modified_start_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the start time in Unix timestamp epoch in seconds where the last modification time of the email/item is more than specified value.",
									},
									"last_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the contacts with specified text in last name.",
									},
									"middle_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the contacts with specified text in middle name.",
									},
									"organizer_address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the calendar items which are organized by specified User's email address.",
									},
									"received_end_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the end time in Unix timestamp epoch in seconds where the received time of the email is less than specified value.",
									},
									"received_start_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the start time in Unix timestamp epoch in seconds where the received time of the email is more than specified value.",
									},
									"recipient_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Filters the emails which are sent to specified email addresses.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"sender_address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Filters the emails which are received from specified User's email address.",
									},
									"source_environment": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the source environment.",
									},
									"task_status_types": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of task item status types. Task items having status within the given types will be returned.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"types": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of mailbox item types. Only items within the given types will be returned.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"o365_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies O365 specific params search request params to search for indexed items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the domain Ids in which indexed items are searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"group_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Group ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"site_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Sharepoint site ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"teams_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Teams ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"user_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the user ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"site_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the request parameters to search for files/folders in document libraries.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category_types": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies a list of document library types. Only items within the given types will be returned.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"creation_end_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the end time in Unix timestamp epoch in seconds when the file/folder is created.",
									},
									"creation_start_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the start time in Unix timestamp epoch in seconds when the file/folder is created.",
									},
									"include_files": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Specifies whether to include files in the response. Default is true.",
									},
									"include_folders": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Specifies whether to include folders in the response. Default is true.",
									},
									"o365_params": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies O365 specific params search request params to search for indexed items.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"domain_ids": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the domain Ids in which indexed items are searched.",
													Elem:        &schema.Schema{Type: schema.TypeInt},
												},
												"group_ids": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the Group ids across which the indexed items needs to be searched.",
													Elem:        &schema.Schema{Type: schema.TypeInt},
												},
												"site_ids": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the Sharepoint site ids across which the indexed items needs to be searched.",
													Elem:        &schema.Schema{Type: schema.TypeInt},
												},
												"teams_ids": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the Teams ids across which the indexed items needs to be searched.",
													Elem:        &schema.Schema{Type: schema.TypeInt},
												},
												"user_ids": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the user ids across which the indexed items needs to be searched.",
													Elem:        &schema.Schema{Type: schema.TypeInt},
												},
											},
										},
									},
									"owner_names": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the list of owner names to filter on owner of the file/folder.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"search_string": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the search string to filter the files/folders. User can specify a wildcard character '*' as a suffix to a string where all item names are matched with the prefix string.",
									},
									"size_bytes_lower_limit": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the minimum size of the file in bytes.",
									},
									"size_bytes_upper_limit": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the maximum size of the file in bytes.",
									},
								},
							},
						},
					},
				},
			},
			"ms_teams_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request params to search for Teams items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of teams files types. Only items within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"channel_names": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the list of channel names to filter while doing search for files.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"channel_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the request parameters related to channels for Microsoft365 teams.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel_email": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the email id of the channel.",
									},
									"channel_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the unique id of the channel.",
									},
									"channel_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the name of the channel. Only items within the specified channel will be returned.",
									},
									"include_private_channels": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Specifies whether to include private channels in the response. Default is true.",
									},
									"include_public_channels": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Default:     true,
										Description: "Specifies whether to include public channels in the response. Default is true.",
									},
								},
							},
						},
						"creation_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds when the item is created.",
						},
						"creation_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds when the item is created.",
						},
						"o365_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies O365 specific params search request params to search for indexed items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the domain Ids in which indexed items are searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"group_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Group ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"site_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Sharepoint site ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"teams_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Teams ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"user_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the user ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"owner_names": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the list of owner email ids to filter on owner of the item.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the search string to filter the items. User can specify a wildcard character '*' as a suffix to a string where all item names are matched with the prefix string.",
						},
						"size_bytes_lower_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the minimum size of the item in bytes.",
						},
						"size_bytes_upper_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum size of the item in bytes.",
						},
						"types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of Teams item types. Only items within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"one_drive_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request parameters to search for files/folders in document libraries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of document library types. Only items within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"creation_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds when the file/folder is created.",
						},
						"creation_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds when the file/folder is created.",
						},
						"include_files": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Specifies whether to include files in the response. Default is true.",
						},
						"include_folders": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Specifies whether to include folders in the response. Default is true.",
						},
						"o365_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies O365 specific params search request params to search for indexed items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the domain Ids in which indexed items are searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"group_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Group ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"site_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Sharepoint site ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"teams_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Teams ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"user_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the user ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"owner_names": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the list of owner names to filter on owner of the file/folder.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the search string to filter the files/folders. User can specify a wildcard character '*' as a suffix to a string where all item names are matched with the prefix string.",
						},
						"size_bytes_lower_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the minimum size of the file in bytes.",
						},
						"size_bytes_upper_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum size of the file in bytes.",
						},
					},
				},
			},
			"public_folder_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request parameters to search for Public Folder items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the search string to filter the items. User can specify a wildcard character '*' as a suffix to a string where all item names are matched with the prefix string.",
						},
						"types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of public folder item types. Only items within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"has_attachment": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Filters the public folder items which have attachment.",
						},
						"sender_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filters the public folder items which are received from specified user's email address.",
						},
						"recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the public folder items which are sent to specified email addresses.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"cc_recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the public folder items which are sent to specified email addresses in CC.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"bcc_recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Filters the public folder items which are sent to specified email addresses in BCC.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"received_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds where the received time of the public folder item is more than specified value.",
						},
						"received_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds where the received time of the public folder items is less than specified value.",
						},
					},
				},
			},
			"sfdc_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the parameters which are specific for searching Salesforce records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mutation_types": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "Specifies a list of mutuation types for an object.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"object_name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the name of the object.",
						},
						"query_string": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the query string to search records. Query string can be one or multiples clauses joined together by 'AND' or 'OR' claused.",
						},
						"snapshot_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the id of the snapshot for the object.",
						},
					},
				},
			},
			"sharepoint_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the request parameters to search for files/folders in document libraries.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category_types": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of document library types. Only items within the given types will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"creation_end_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the end time in Unix timestamp epoch in seconds when the file/folder is created.",
						},
						"creation_start_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the start time in Unix timestamp epoch in seconds when the file/folder is created.",
						},
						"include_files": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Specifies whether to include files in the response. Default is true.",
						},
						"include_folders": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Specifies whether to include folders in the response. Default is true.",
						},
						"o365_params": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies O365 specific params search request params to search for indexed items.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the domain Ids in which indexed items are searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"group_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Group ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"site_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Sharepoint site ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"teams_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the Teams ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
									"user_ids": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the user ids across which the indexed items needs to be searched.",
										Elem:        &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
						"owner_names": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies the list of owner names to filter on owner of the file/folder.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the search string to filter the files/folders. User can specify a wildcard character '*' as a suffix to a string where all item names are matched with the prefix string.",
						},
						"size_bytes_lower_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the minimum size of the file in bytes.",
						},
						"size_bytes_upper_limit": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the maximum size of the file in bytes.",
						},
					},
				},
			},
			"uda_params": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				ForceNew:    true,
				Description: "Parameters required to search Universal Data Adapter objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"search_string": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the search string to search the Universal Data Adapter Objects.",
						},
						"source_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Specifies a list of source ids. Only files found in these sources will be returned.",
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"cassandra_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Cassandra objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
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
										Type: schema.TypeInt,

										Computed:    true,
										Description: "Specifies registered source id to which object belongs.",
									},
									"source_name": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies registered source name to which object belongs.",
									},
									"environment": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the environment of the object.",
									},
									"object_hash": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the hash identifier of the object.",
									},
									"object_type": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the type of the object.",
									},
									"logical_size_bytes": &schema.Schema{
										Type: schema.TypeInt,

										Computed:    true,
										Description: "Specifies the logical size of object in bytes.",
									},
									"uuid": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the uuid which is a unique identifier of the object.",
									},
									"global_id": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the global id which is a unique identifier of the object.",
									},
									"protection_type": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the protection type of the object if any.",
									},
									"sharepoint_site_summary": &schema.Schema{
										Type: schema.TypeList,

										Computed:    true,
										Description: "Specifies the common parameters for Sharepoint site objects.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"site_web_url": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the web url for the Sharepoint site.",
												},
											},
										},
									},
									"os_type": &schema.Schema{
										Type: schema.TypeString,

										Computed:    true,
										Description: "Specifies the operating system type of the object.",
									},
									"child_objects": &schema.Schema{
										Type: schema.TypeList,

										Computed:    true,
										Description: "Specifies child object details.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type: schema.TypeInt,

													Computed:    true,
													Description: "Specifies object id.",
												},
												"name": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the name of the object.",
												},
												"source_id": &schema.Schema{
													Type: schema.TypeInt,

													Computed:    true,
													Description: "Specifies registered source id to which object belongs.",
												},
												"source_name": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies registered source name to which object belongs.",
												},
												"environment": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the environment of the object.",
												},
												"object_hash": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the hash identifier of the object.",
												},
												"object_type": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the type of the object.",
												},
												"logical_size_bytes": &schema.Schema{
													Type: schema.TypeInt,

													Computed:    true,
													Description: "Specifies the logical size of object in bytes.",
												},
												"uuid": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the uuid which is a unique identifier of the object.",
												},
												"global_id": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the global id which is a unique identifier of the object.",
												},
												"protection_type": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the protection type of the object if any.",
												},
												"sharepoint_site_summary": &schema.Schema{
													Type: schema.TypeList,

													Computed:    true,
													Description: "Specifies the common parameters for Sharepoint site objects.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"site_web_url": &schema.Schema{
																Type: schema.TypeString,

																Computed:    true,
																Description: "Specifies the web url for the Sharepoint site.",
															},
														},
													},
												},
												"os_type": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the operating system type of the object.",
												},
												"child_objects": &schema.Schema{
													Type: schema.TypeList,

													Computed:    true,
													Description: "Specifies child object details.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{},
													},
												},
												"v_center_summary": &schema.Schema{
													Type: schema.TypeList,

													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"is_cloud_env": &schema.Schema{
																Type: schema.TypeBool,

																Computed:    true,
																Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
															},
														},
													},
												},
												"windows_cluster_summary": &schema.Schema{
													Type: schema.TypeList,

													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_source_type": &schema.Schema{
																Type: schema.TypeString,

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
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_cloud_env": &schema.Schema{
													Type: schema.TypeBool,

													Computed:    true,
													Description: "Specifies that registered vCenter source is a VMC (VMware Cloud) environment or not.",
												},
											},
										},
									},
									"windows_cluster_summary": &schema.Schema{
										Type: schema.TypeList,

										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_source_type": &schema.Schema{
													Type: schema.TypeString,

													Computed:    true,
													Description: "Specifies the type of cluster resource this source represents.",
												},
											},
										},
									},
								},
							},
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"keyspace_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies type of Keyspace.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the Cassandra Object type.",
						},
					},
				},
			},
			"couchbase_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Couchbase objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Couchbase Object Type. For Couchbase this is alywas set to Bucket.",
						},
					},
				},
			},
			"emails": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed emails and email folders.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"bcc_recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the email addresses of all the BCC receipients of this email.\".",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"cc_recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the email addresses of all the CC receipients of this email.\".",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"created_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Unix timestamp epoch in seconds at which this item is created.\".",
						},
						"directory_path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the directory path to this mailbox item.",
						},
						"email_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the email addresses of a contact.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"email_subject": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the subject of this email.",
						},
						"first_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the contact's first name.",
						},
						"folder_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specify the name of the email folder.",
						},
						"has_attachment": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Specifies whether email has an attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the email object.",
						},
						"last_modification_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the name of the person who modified this item.\".",
						},
						"last_modification_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Unix timestamp epoch in seconds at which this item was modified.\".",
						},
						"last_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the contact's last name.",
						},
						"optional_attendees_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the email addresses of all the optional attendees of this calendar item.\".",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"organizer_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the calendar item organizer's email address.\".",
						},
						"parent_folder_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of parent folder the mailbox item.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path to this mailbox item.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Protection Group id protecting the mailbox.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Protection Group name protecting the mailbox item.\".",
						},
						"received_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Unix timestamp epoch in seconds at which this email is received.\".",
						},
						"recipient_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the email addresses of all receipients of this email.\".",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"required_attendees_addresses": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the email addresses of all required attendees of this calendar item.\".",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"sender_address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the sender's email address.",
						},
						"sent_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Unix timestamp epoch in seconds at which this email is sent.\".",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"task_completion_date_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Unix timestamp epoch in seconds at which this task item was completed.\".",
						},
						"task_due_date_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Unix timestamp epoch in seconds at which this task item is due.\".",
						},
						"task_status": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the task item status type.",
						},
						"tenant_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specify the tenant id to which this email belongs to.\".",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Mailbox item type.",
						},
						"user_object_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Object Summary.",
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
					},
				},
			},
			"exchange_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed HDFS objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"database_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the Exchange database corresponding to the mailbox.",
						},
						"email": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the email corresponding to the mailbox.",
						},
						"object_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the Exchange mailbox.",
						},
					},
				},
			},
			"files": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed files.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the file name.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path to this file.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the file type.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this file.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this file.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
									"protection_stats": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the count and size of protected and unprotected objects for the size.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the environment of the object.",
												},
												"protected_count": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the count of the protected leaf objects.",
												},
												"unprotected_count": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the count of the unprotected leaf objects.",
												},
												"deleted_protected_count": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the count of protected leaf objects which were deleted from the source after being protected.",
												},
												"protected_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the protected logical size in bytes.",
												},
												"unprotected_size_bytes": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the unprotected logical size in bytes.",
												},
											},
										},
									},
									"permissions": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the list of users, groups and users that have permissions for a given object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"object_id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the id of the object.",
												},
												"users": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the list of users which has the permissions to the object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the name of the user.",
															},
															"sid": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the sid of the user.",
															},
															"domain": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the domain of the user.",
															},
														},
													},
												},
												"groups": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the list of user groups which has permissions to the object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the name of the user group.",
															},
															"sid": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the sid of the user group.",
															},
															"domain": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the domain of the user group.",
															},
														},
													},
												},
												"tenant": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
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
											},
										},
									},
									"mssql_params": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the parameters for Msssql object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aag_info": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Object details for Mssql.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the AAG name.",
															},
															"object_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the AAG object Id.",
															},
														},
													},
												},
												"host_info": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the host information for a objects. This is mainly populated in case of App objects where app object is hosted by another object such as VM or physical server.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the id of the host object.",
															},
															"name": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the name of the host object.",
															},
															"environment": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the environment of the object.",
															},
														},
													},
												},
												"is_encrypted": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Specifies whether the database is TDE enabled.",
												},
											},
										},
									},
									"physical_params": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the parameters for Physical object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable_system_backup": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
													Computed:    true,
													Description: "Specifies if system backup was enabled for the source in a particular run.",
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
			"hbase_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Hbase objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Hbase Object Type.",
						},
					},
				},
			},
			"hdfs_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed HDFS objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the HDFS Object Type.",
						},
					},
				},
			},
			"hive_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Hive objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Hive Object Type.",
						},
					},
				},
			},
			"mongo_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Mongo objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"cdp_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the CDP related information for a given object. This field will only be populated when protection group is configured with policy having CDP retention settings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"allow_re_enable_cdp": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Specifies if re-enabling CDP is allowed or not through UI without any job or policy update through API.",
									},
									"cdp_enabled": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Specifies whether CDP is currently active or not. CDP might have been active on this object before, but it might not be anymore.",
									},
									"last_run_info": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the last backup information for a given CDP object.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"local_backup_info": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Computed:    true,
													Description: "Specifies the last local backup information for a given CDP object.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"end_time_in_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the end time of the last local backup taken.",
															},
															"start_time_in_usecs": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Computed:    true,
																Description: "Specifies the start time of the last local backup taken.",
															},
														},
													},
												},
											},
										},
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the protection group id to which this CDP object belongs.",
									},
								},
							},
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Mongo Object Type.",
						},
					},
				},
			},
			"ms_group_items": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed M365 Groups items like group mail items, files etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"mailbox_item": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies an email or an email folder.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
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
										Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies runs the tags are applied to.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"bcc_recipient_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the email addresses of all the BCC receipients of this email.\".",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"cc_recipient_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the email addresses of all the CC receipients of this email.\".",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"created_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Unix timestamp epoch in seconds at which this item is created.\".",
									},
									"directory_path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the directory path to this mailbox item.",
									},
									"email_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the email addresses of a contact.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"email_subject": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the subject of this email.",
									},
									"first_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the contact's first name.",
									},
									"folder_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specify the name of the email folder.",
									},
									"has_attachment": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: "Specifies whether email has an attachment.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the id of the email object.",
									},
									"last_modification_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the name of the person who modified this item.\".",
									},
									"last_modification_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Unix timestamp epoch in seconds at which this item was modified.\".",
									},
									"last_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the contact's last name.",
									},
									"optional_attendees_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the email addresses of all the optional attendees of this calendar item.\".",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"organizer_address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the calendar item organizer's email address.\".",
									},
									"parent_folder_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the id of parent folder the mailbox item.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the path to this mailbox item.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Protection Group id protecting the mailbox.\".",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Protection Group name protecting the mailbox item.\".",
									},
									"received_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Unix timestamp epoch in seconds at which this email is received.\".",
									},
									"recipient_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the email addresses of all receipients of this email.\".",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"required_attendees_addresses": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the email addresses of all required attendees of this calendar item.\".",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"sender_address": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the sender's email address.",
									},
									"sent_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Unix timestamp epoch in seconds at which this email is sent.\".",
									},
									"storage_domain_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
									},
									"task_completion_date_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Unix timestamp epoch in seconds at which this task item was completed.\".",
									},
									"task_due_date_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Unix timestamp epoch in seconds at which this task item is due.\".",
									},
									"task_status": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the task item status type.",
									},
									"tenant_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specify the tenant id to which this email belongs to.\".",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Mailbox item type.",
									},
									"user_object_info": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Object Summary.",
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
								},
							},
						},
						"site_item": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies a Document Library indexed item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
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
										Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies runs the tags are applied to.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the object.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the path of the object.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the protection group id which contains this object.\".",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the protection group name which contains this object.\".",
									},
									"policy_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the protection policy id for this file.",
									},
									"policy_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the protection policy name for this file.",
									},
									"storage_domain_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
									},
									"source_info": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Source Object information.",
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
									"creation_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Unix timestamp epoch in seconds at which this item is created.",
									},
									"file_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the file type.",
									},
									"item_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the id of the document library item.",
									},
									"item_size": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the size in bytes for the indexed item.",
									},
									"owner_email": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the email of the owner of the document library item.",
									},
									"owner_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the owner of the document library item.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the M365 Group item type.",
						},
					},
				},
			},
			"one_drive_items": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed one drive items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"creation_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Unix timestamp epoch in seconds at which this item is created.",
						},
						"file_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the file type.",
						},
						"item_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the document library item.",
						},
						"item_size": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the size in bytes for the indexed item.",
						},
						"owner_email": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the email of the owner of the document library item.",
						},
						"owner_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the owner of the document library item.",
						},
					},
				},
			},
			"public_folder_items": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Public folder items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Public folder item type.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed item.",
						},
						"subject": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the subject of the indexed item.",
						},
						"has_attachments": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Specifies whether the item has any attachments.",
						},
						"item_class": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the item class of the indexed item.",
						},
						"received_time_secs": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Unix timestamp epoch in seconds at which this item is received.",
						},
						"item_size": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the size in bytes for the indexed item.",
						},
						"parent_folder_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of parent folder the indexed item.",
						},
					},
				},
			},
			"sfdc_records": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the list of salesforce records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"column_names": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the column names for the records.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"records": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Each record is represented by an array of strings having the same order as the 'columnNames'.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"sharepoint_items": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Sharepoint items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"document_library_item": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies a Document Library indexed item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tags": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
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
										Optional:    true,
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
													Optional:    true,
													Computed:    true,
													Description: "Specifies runs the tags are applied to.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the object.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the path of the object.",
									},
									"protection_group_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the protection group id which contains this object.\".",
									},
									"protection_group_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the protection group name which contains this object.\".",
									},
									"policy_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the protection policy id for this file.",
									},
									"policy_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the protection policy name for this file.",
									},
									"storage_domain_id": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
									},
									"source_info": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Source Object information.",
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
									"creation_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Unix timestamp epoch in seconds at which this item is created.",
									},
									"file_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the file type.",
									},
									"item_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the id of the document library item.",
									},
									"item_size": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the size in bytes for the indexed item.",
									},
									"owner_email": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the email of the owner of the document library item.",
									},
									"owner_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the owner of the document library item.",
									},
								},
							},
						},
					},
				},
			},
			"teams_items": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed M365 Teams items like channels, files etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"channel_item": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies a M365 Teams channel item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel_email": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the email of this channel.",
									},
									"channel_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the id of this channel.",
									},
									"channel_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the channel name.",
									},
									"channel_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the channel type.",
									},
									"creation_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Unix timestamp epoch in seconds at which this channel is created.",
									},
									"owner_names": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the names of owners of this channel.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"file_item": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies a M365 Teams channel file item.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"creation_time_secs": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the Unix timestamp epoch in seconds at which this item is created.",
									},
									"drive_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the name of the drive location for this file.",
									},
									"file_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the file type.",
									},
									"item_size": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Specifies the size in bytes for the indexed item.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the M365 Teams item type.",
						},
					},
				},
			},
			"uda_objects": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies the indexed Universal Data Adapter objects.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tags": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
										Computed:    true,
										Description: "Specifies runs the tags are applied to.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the name of the object.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the path of the object.",
						},
						"protection_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group id which contains this object.\".",
						},
						"protection_group_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the protection group name which contains this object.\".",
						},
						"policy_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy id for this file.",
						},
						"policy_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the protection policy name for this file.",
						},
						"storage_domain_id": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "\"Specifies the Storage Domain id where the backup data of Object is present.\".",
						},
						"source_info": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the Source Object information.",
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
						"full_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the full name of the indexed object.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the id of the indexed object.",
						},
						"object_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the type of the indexed object.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmBackupRecoverySearchIndexedObjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	searchIndexedObjectsOptions := &backuprecoveryv1.SearchIndexedObjectsOptions{}

	searchIndexedObjectsOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	searchIndexedObjectsOptions.SetObjectType(d.Get("object_type").(string))
	if _, ok := d.GetOk("protection_group_ids"); ok {
		var protectionGroupIds []string
		for _, v := range d.Get("protection_group_ids").([]interface{}) {
			protectionGroupIdsItem := v.(string)
			protectionGroupIds = append(protectionGroupIds, protectionGroupIdsItem)
		}
		searchIndexedObjectsOptions.SetProtectionGroupIds(protectionGroupIds)
	}
	if _, ok := d.GetOk("storage_domain_ids"); ok {
		var storageDomainIds []int64
		for _, v := range d.Get("storage_domain_ids").([]interface{}) {
			storageDomainIdsItem := int64(v.(int))
			storageDomainIds = append(storageDomainIds, storageDomainIdsItem)
		}
		searchIndexedObjectsOptions.SetStorageDomainIds(storageDomainIds)
	}
	if _, ok := d.GetOk("tenant_id"); ok {
		searchIndexedObjectsOptions.SetTenantID(d.Get("tenant_id").(string))
	}
	if _, ok := d.GetOk("include_tenants"); ok {
		searchIndexedObjectsOptions.SetIncludeTenants(d.Get("include_tenants").(bool))
	}
	if _, ok := d.GetOk("tags"); ok {
		var tags []string
		for _, v := range d.Get("tags").([]interface{}) {
			tagsItem := v.(string)
			tags = append(tags, tagsItem)
		}
		searchIndexedObjectsOptions.SetTags(tags)
	}
	if _, ok := d.GetOk("snapshot_tags"); ok {
		var snapshotTags []string
		for _, v := range d.Get("snapshot_tags").([]interface{}) {
			snapshotTagsItem := v.(string)
			snapshotTags = append(snapshotTags, snapshotTagsItem)
		}
		searchIndexedObjectsOptions.SetSnapshotTags(snapshotTags)
	}
	if _, ok := d.GetOk("must_have_tag_ids"); ok {
		var mustHaveTagIds []string
		for _, v := range d.Get("must_have_tag_ids").([]interface{}) {
			mustHaveTagIdsItem := v.(string)
			mustHaveTagIds = append(mustHaveTagIds, mustHaveTagIdsItem)
		}
		searchIndexedObjectsOptions.SetMustHaveTagIds(mustHaveTagIds)
	}
	if _, ok := d.GetOk("might_have_tag_ids"); ok {
		var mightHaveTagIds []string
		for _, v := range d.Get("might_have_tag_ids").([]interface{}) {
			mightHaveTagIdsItem := v.(string)
			mightHaveTagIds = append(mightHaveTagIds, mightHaveTagIdsItem)
		}
		searchIndexedObjectsOptions.SetMightHaveTagIds(mightHaveTagIds)
	}
	if _, ok := d.GetOk("must_have_snapshot_tag_ids"); ok {
		var mustHaveSnapshotTagIds []string
		for _, v := range d.Get("must_have_snapshot_tag_ids").([]interface{}) {
			mustHaveSnapshotTagIdsItem := v.(string)
			mustHaveSnapshotTagIds = append(mustHaveSnapshotTagIds, mustHaveSnapshotTagIdsItem)
		}
		searchIndexedObjectsOptions.SetMustHaveSnapshotTagIds(mustHaveSnapshotTagIds)
	}
	if _, ok := d.GetOk("might_have_snapshot_tag_ids"); ok {
		var mightHaveSnapshotTagIds []string
		for _, v := range d.Get("might_have_snapshot_tag_ids").([]interface{}) {
			mightHaveSnapshotTagIdsItem := v.(string)
			mightHaveSnapshotTagIds = append(mightHaveSnapshotTagIds, mightHaveSnapshotTagIdsItem)
		}
		searchIndexedObjectsOptions.SetMightHaveSnapshotTagIds(mightHaveSnapshotTagIds)
	}
	if _, ok := d.GetOk("pagination_cookie"); ok {
		searchIndexedObjectsOptions.SetPaginationCookie(d.Get("pagination_cookie").(string))
	}
	if _, ok := d.GetOk("object_count"); ok {
		searchIndexedObjectsOptions.SetCount(int64(d.Get("object_count").(int)))
	}
	if _, ok := d.GetOk("use_cached_data"); ok {
		searchIndexedObjectsOptions.SetUseCachedData(d.Get("use_cached_data").(bool))
	}
	if _, ok := d.GetOk("cassandra_params"); ok {
		cassandraParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToCassandraOnPremSearchParams(d.Get("cassandra_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-cassandra_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetCassandraParams(cassandraParamsModel)
	}
	if _, ok := d.GetOk("couchbase_params"); ok {
		couchbaseParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToCouchBaseOnPremSearchParams(d.Get("couchbase_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-couchbase_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetCouchbaseParams(couchbaseParamsModel)
	}
	if _, ok := d.GetOk("email_params"); ok {
		emailParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchEmailRequestParams(d.Get("email_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-email_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetEmailParams(emailParamsModel)
	}
	if _, ok := d.GetOk("exchange_params"); ok {
		exchangeParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchExchangeObjectsRequestParams(d.Get("exchange_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-exchange_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetExchangeParams(exchangeParamsModel)
	}
	if _, ok := d.GetOk("file_params"); ok {
		fileParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchFileRequestParams(d.Get("file_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-file_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetFileParams(fileParamsModel)
	}
	if _, ok := d.GetOk("hbase_params"); ok {
		hbaseParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToHbaseOnPremSearchParams(d.Get("hbase_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-hbase_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetHbaseParams(hbaseParamsModel)
	}
	if _, ok := d.GetOk("hdfs_params"); ok {
		hdfsParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToHDFSOnPremSearchParams(d.Get("hdfs_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-hdfs_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetHdfsParams(hdfsParamsModel)
	}
	if _, ok := d.GetOk("hive_params"); ok {
		hiveParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToHiveOnPremSearchParams(d.Get("hive_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-hive_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetHiveParams(hiveParamsModel)
	}
	if _, ok := d.GetOk("mongodb_params"); ok {
		mongodbParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToMongoDbOnPremSearchParams(d.Get("mongodb_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-mongodb_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetMongodbParams(mongodbParamsModel)
	}
	if _, ok := d.GetOk("ms_groups_params"); ok {
		msGroupsParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchMsGroupsRequestParams(d.Get("ms_groups_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-ms_groups_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetMsGroupsParams(msGroupsParamsModel)
	}
	if _, ok := d.GetOk("ms_teams_params"); ok {
		msTeamsParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchMsTeamsRequestParams(d.Get("ms_teams_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-ms_teams_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetMsTeamsParams(msTeamsParamsModel)
	}
	if _, ok := d.GetOk("one_drive_params"); ok {
		oneDriveParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchDocumentLibraryRequestParams(d.Get("one_drive_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-one_drive_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetOneDriveParams(oneDriveParamsModel)
	}
	if _, ok := d.GetOk("public_folder_params"); ok {
		publicFolderParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchPublicFolderRequestParams(d.Get("public_folder_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-public_folder_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetPublicFolderParams(publicFolderParamsModel)
	}
	if _, ok := d.GetOk("sfdc_params"); ok {
		sfdcParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchSfdcRecordsRequestParams(d.Get("sfdc_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-sfdc_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetSfdcParams(sfdcParamsModel)
	}
	if _, ok := d.GetOk("sharepoint_params"); ok {
		sharepointParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchDocumentLibraryRequestParams(d.Get("sharepoint_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-sharepoint_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetSharepointParams(sharepointParamsModel)
	}
	if _, ok := d.GetOk("uda_params"); ok {
		udaParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToUdaOnPremSearchParams(d.Get("uda_params.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_search_indexed_object", "create", "parse-uda_params").GetDiag()
		}
		searchIndexedObjectsOptions.SetUdaParams(udaParamsModel)
	}

	searchIndexedObjectsResponse, _, err := backupRecoveryClient.SearchIndexedObjectsWithContext(context, searchIndexedObjectsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("SearchIndexedObjectsWithContext failed: %s", err.Error()), "ibm_search_indexed_object", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(DataSourceIbmBackupRecoverySearchIndexedObjectID(d))

	if searchIndexedObjectsResponse != nil {
		if err = d.Set("object_type", searchIndexedObjectsResponse.ObjectType); err != nil {
			err = fmt.Errorf("Error setting object_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-object_type").GetDiag()
		}
		if !core.IsNil(searchIndexedObjectsResponse.Count) {
			if err = d.Set("object_count", flex.IntValue(searchIndexedObjectsResponse.Count)); err != nil {
				err = fmt.Errorf("Error setting count: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-count").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.PaginationCookie) {
			if err = d.Set("pagination_cookie", searchIndexedObjectsResponse.PaginationCookie); err != nil {
				err = fmt.Errorf("Error setting pagination_cookie: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-pagination_cookie").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.CassandraObjects) {
			cassandraObjects := []map[string]interface{}{}
			for _, cassandraObjectsItem := range searchIndexedObjectsResponse.CassandraObjects {
				cassandraObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCassandraIndexedObjectToMap(&cassandraObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "cassandra_objects-to-map").GetDiag()
				}
				cassandraObjects = append(cassandraObjects, cassandraObjectsItemMap)
			}
			if err = d.Set("cassandra_objects", cassandraObjects); err != nil {
				err = fmt.Errorf("Error setting cassandra_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-cassandra_objects").GetDiag()
			}
		} else {
			if err = d.Set("cassandra_objects", []map[string]interface{}{}); err != nil {
				err = fmt.Errorf("Error setting cassandra_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-cassandra_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.CouchbaseObjects) {
			couchbaseObjects := []map[string]interface{}{}
			for _, couchbaseObjectsItem := range searchIndexedObjectsResponse.CouchbaseObjects {
				couchbaseObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCouchbaseIndexedObjectToMap(&couchbaseObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "couchbase_objects-to-map").GetDiag()
				}
				couchbaseObjects = append(couchbaseObjects, couchbaseObjectsItemMap)
			}
			if err = d.Set("couchbase_objects", couchbaseObjects); err != nil {
				err = fmt.Errorf("Error setting couchbase_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-couchbase_objects").GetDiag()
			}
		} else {
			if err = d.Set("couchbase_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting couchbase_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-couchbase_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.Emails) {
			emails := []map[string]interface{}{}
			for _, emailsItem := range searchIndexedObjectsResponse.Emails {
				emailsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectEmailToMap(&emailsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "emails-to-map").GetDiag()
				}
				emails = append(emails, emailsItemMap)
			}
			if err = d.Set("emails", emails); err != nil {
				err = fmt.Errorf("Error setting emails: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-emails").GetDiag()
			}
		} else {
			if err = d.Set("emails", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting emails: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-emails").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.ExchangeObjects) {
			exchangeObjects := []map[string]interface{}{}
			for _, exchangeObjectsItem := range searchIndexedObjectsResponse.ExchangeObjects {
				exchangeObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectExchangeIndexedObjectToMap(&exchangeObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "exchange_objects-to-map").GetDiag()
				}
				exchangeObjects = append(exchangeObjects, exchangeObjectsItemMap)
			}
			if err = d.Set("exchange_objects", exchangeObjects); err != nil {
				err = fmt.Errorf("Error setting exchange_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-exchange_objects").GetDiag()
			}
		} else {
			if err = d.Set("exchange_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting exchange_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-exchange_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.Files) {
			files := []map[string]interface{}{}
			for _, filesItem := range searchIndexedObjectsResponse.Files {
				filesItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectFileToMap(&filesItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "files-to-map").GetDiag()
				}
				files = append(files, filesItemMap)
			}
			if err = d.Set("files", files); err != nil {
				err = fmt.Errorf("Error setting files: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-files").GetDiag()
			}
		} else {
			if err = d.Set("files", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting files: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-files").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.HbaseObjects) {
			hbaseObjects := []map[string]interface{}{}
			for _, hbaseObjectsItem := range searchIndexedObjectsResponse.HbaseObjects {
				hbaseObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectHbaseIndexedObjectToMap(&hbaseObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "hbase_objects-to-map").GetDiag()
				}
				hbaseObjects = append(hbaseObjects, hbaseObjectsItemMap)
			}
			if err = d.Set("hbase_objects", hbaseObjects); err != nil {
				err = fmt.Errorf("Error setting hbase_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-hbase_objects").GetDiag()
			}
		} else {
			if err = d.Set("hbase_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting hbase_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-hbase_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.HdfsObjects) {
			hdfsObjects := []map[string]interface{}{}
			for _, hdfsObjectsItem := range searchIndexedObjectsResponse.HdfsObjects {
				hdfsObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectHDFSIndexedObjectToMap(&hdfsObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "hdfs_objects-to-map").GetDiag()
				}
				hdfsObjects = append(hdfsObjects, hdfsObjectsItemMap)
			}
			if err = d.Set("hdfs_objects", hdfsObjects); err != nil {
				err = fmt.Errorf("Error setting hdfs_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-hdfs_objects").GetDiag()
			}
		} else {
			if err = d.Set("hdfs_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting hdfs_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-hdfs_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.HiveObjects) {
			hiveObjects := []map[string]interface{}{}
			for _, hiveObjectsItem := range searchIndexedObjectsResponse.HiveObjects {
				hiveObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectHiveIndexedObjectToMap(&hiveObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "hive_objects-to-map").GetDiag()
				}
				hiveObjects = append(hiveObjects, hiveObjectsItemMap)
			}
			if err = d.Set("hive_objects", hiveObjects); err != nil {
				err = fmt.Errorf("Error setting hive_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-hive_objects").GetDiag()
			}
		} else {
			if err = d.Set("hive_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting hive_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-hive_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.MongoObjects) {
			mongoObjects := []map[string]interface{}{}
			for _, mongoObjectsItem := range searchIndexedObjectsResponse.MongoObjects {
				mongoObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectMongoIndexedObjectToMap(&mongoObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "mongo_objects-to-map").GetDiag()
				}
				mongoObjects = append(mongoObjects, mongoObjectsItemMap)
			}
			if err = d.Set("mongo_objects", mongoObjects); err != nil {
				err = fmt.Errorf("Error setting mongo_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-mongo_objects").GetDiag()
			}
		} else {
			if err = d.Set("mongo_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting mongo_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-mongo_objects").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.MsGroupItems) {
			msGroupItems := []map[string]interface{}{}
			for _, msGroupItemsItem := range searchIndexedObjectsResponse.MsGroupItems {
				msGroupItemsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectMsGroupItemToMap(&msGroupItemsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "ms_group_items-to-map").GetDiag()
				}
				msGroupItems = append(msGroupItems, msGroupItemsItemMap)
			}
			if err = d.Set("ms_group_items", msGroupItems); err != nil {
				err = fmt.Errorf("Error setting ms_group_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-ms_group_items").GetDiag()
			}
		} else {
			if err = d.Set("ms_group_items", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting ms_group_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-ms_group_items").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.OneDriveItems) {
			oneDriveItems := []map[string]interface{}{}
			for _, oneDriveItemsItem := range searchIndexedObjectsResponse.OneDriveItems {
				oneDriveItemsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectDocumentLibraryItemToMap(&oneDriveItemsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "one_drive_items-to-map").GetDiag()
				}
				oneDriveItems = append(oneDriveItems, oneDriveItemsItemMap)
			}
			if err = d.Set("one_drive_items", oneDriveItems); err != nil {
				err = fmt.Errorf("Error setting one_drive_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-one_drive_items").GetDiag()
			}
		} else {
			if err = d.Set("one_drive_items", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting one_drive_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-one_drive_items").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.PublicFolderItems) {
			publicFolderItems := []map[string]interface{}{}
			for _, publicFolderItemsItem := range searchIndexedObjectsResponse.PublicFolderItems {
				publicFolderItemsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectPublicFolderItemToMap(&publicFolderItemsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "public_folder_items-to-map").GetDiag()
				}
				publicFolderItems = append(publicFolderItems, publicFolderItemsItemMap)
			}
			if err = d.Set("public_folder_items", publicFolderItems); err != nil {
				err = fmt.Errorf("Error setting public_folder_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-public_folder_items").GetDiag()
			}
		} else {
			if err = d.Set("public_folder_items", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting public_folder_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-public_folder_items").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.SfdcRecords) {
			sfdcRecordsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSfdcRecordsToMap(searchIndexedObjectsResponse.SfdcRecords)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "sfdc_records-to-map").GetDiag()
			}
			if err = d.Set("sfdc_records", []map[string]interface{}{sfdcRecordsMap}); err != nil {
				err = fmt.Errorf("Error setting sfdc_records: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-sfdc_records").GetDiag()
			}
		} else {
			if err = d.Set("sfdc_records", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting sfdc_records: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-sfdc_records").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.SharepointItems) {
			sharepointItems := []map[string]interface{}{}
			for _, sharepointItemsItem := range searchIndexedObjectsResponse.SharepointItems {
				sharepointItemsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointItemToMap(&sharepointItemsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "sharepoint_items-to-map").GetDiag()
				}
				sharepointItems = append(sharepointItems, sharepointItemsItemMap)
			}
			if err = d.Set("sharepoint_items", sharepointItems); err != nil {
				err = fmt.Errorf("Error setting sharepoint_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-sharepoint_items").GetDiag()
			}
		} else {
			if err = d.Set("sharepoint_items", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting sharepoint_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-sharepoint_items").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.TeamsItems) {
			teamsItems := []map[string]interface{}{}
			for _, teamsItemsItem := range searchIndexedObjectsResponse.TeamsItems {
				teamsItemsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTeamsItemToMap(&teamsItemsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "teams_items-to-map").GetDiag()
				}
				teamsItems = append(teamsItems, teamsItemsItemMap)
			}
			if err = d.Set("teams_items", teamsItems); err != nil {
				err = fmt.Errorf("Error setting teams_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-teams_items").GetDiag()
			}
		} else {
			if err = d.Set("teams_items", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting teams_items: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-teams_items").GetDiag()
			}
		}
		if !core.IsNil(searchIndexedObjectsResponse.UdaObjects) {
			udaObjects := []map[string]interface{}{}
			for _, udaObjectsItem := range searchIndexedObjectsResponse.UdaObjects {
				udaObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectUdaIndexedObjectToMap(&udaObjectsItem) // #nosec G601
				if err != nil {
					return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "uda_objects-to-map").GetDiag()
				}
				udaObjects = append(udaObjects, udaObjectsItemMap)
			}
			if err = d.Set("uda_objects", udaObjects); err != nil {
				err = fmt.Errorf("Error setting uda_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-uda_objects").GetDiag()
			}
		} else {
			if err = d.Set("uda_objects", []interface{}{}); err != nil {
				err = fmt.Errorf("Error setting uda_objects: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_search_indexed_object", "read", "set-uda_objects").GetDiag()
			}
		}
	}
	return nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToCassandraOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.CassandraOnPremSearchParams, error) {
	model := &backuprecoveryv1.CassandraOnPremSearchParams{}
	cassandraObjectTypes := []string{}
	for _, cassandraObjectTypesItem := range modelMap["cassandra_object_types"].([]interface{}) {
		cassandraObjectTypes = append(cassandraObjectTypes, cassandraObjectTypesItem.(string))
	}
	model.CassandraObjectTypes = cassandraObjectTypes
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToCouchBaseOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.CouchBaseOnPremSearchParams, error) {
	model := &backuprecoveryv1.CouchBaseOnPremSearchParams{}
	couchbaseObjectTypes := []string{}
	for _, couchbaseObjectTypesItem := range modelMap["couchbase_object_types"].([]interface{}) {
		couchbaseObjectTypes = append(couchbaseObjectTypes, couchbaseObjectTypesItem.(string))
	}
	model.CouchbaseObjectTypes = couchbaseObjectTypes
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchEmailRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchEmailRequestParams, error) {
	model := &backuprecoveryv1.SearchEmailRequestParams{}
	if modelMap["attendees_addresses"] != nil {
		attendeesAddresses := []string{}
		for _, attendeesAddressesItem := range modelMap["attendees_addresses"].([]interface{}) {
			attendeesAddresses = append(attendeesAddresses, attendeesAddressesItem.(string))
		}
		model.AttendeesAddresses = attendeesAddresses
	}
	if modelMap["bcc_recipient_addresses"] != nil {
		bccRecipientAddresses := []string{}
		for _, bccRecipientAddressesItem := range modelMap["bcc_recipient_addresses"].([]interface{}) {
			bccRecipientAddresses = append(bccRecipientAddresses, bccRecipientAddressesItem.(string))
		}
		model.BccRecipientAddresses = bccRecipientAddresses
	}
	if modelMap["cc_recipient_addresses"] != nil {
		ccRecipientAddresses := []string{}
		for _, ccRecipientAddressesItem := range modelMap["cc_recipient_addresses"].([]interface{}) {
			ccRecipientAddresses = append(ccRecipientAddresses, ccRecipientAddressesItem.(string))
		}
		model.CcRecipientAddresses = ccRecipientAddresses
	}
	if modelMap["created_end_time_secs"] != nil {
		model.CreatedEndTimeSecs = core.Int64Ptr(int64(modelMap["created_end_time_secs"].(int)))
	}
	if modelMap["created_start_time_secs"] != nil {
		model.CreatedStartTimeSecs = core.Int64Ptr(int64(modelMap["created_start_time_secs"].(int)))
	}
	if modelMap["due_date_end_time_secs"] != nil {
		model.DueDateEndTimeSecs = core.Int64Ptr(int64(modelMap["due_date_end_time_secs"].(int)))
	}
	if modelMap["due_date_start_time_secs"] != nil {
		model.DueDateStartTimeSecs = core.Int64Ptr(int64(modelMap["due_date_start_time_secs"].(int)))
	}
	if modelMap["email_address"] != nil && modelMap["email_address"].(string) != "" {
		model.EmailAddress = core.StringPtr(modelMap["email_address"].(string))
	}
	if modelMap["email_subject"] != nil && modelMap["email_subject"].(string) != "" {
		model.EmailSubject = core.StringPtr(modelMap["email_subject"].(string))
	}
	if modelMap["first_name"] != nil && modelMap["first_name"].(string) != "" {
		model.FirstName = core.StringPtr(modelMap["first_name"].(string))
	}
	if modelMap["folder_names"] != nil {
		folderNames := []string{}
		for _, folderNamesItem := range modelMap["folder_names"].([]interface{}) {
			folderNames = append(folderNames, folderNamesItem.(string))
		}
		model.FolderNames = folderNames
	}
	if modelMap["has_attachment"] != nil {
		model.HasAttachment = core.BoolPtr(modelMap["has_attachment"].(bool))
	}
	if modelMap["last_modified_end_time_secs"] != nil {
		model.LastModifiedEndTimeSecs = core.Int64Ptr(int64(modelMap["last_modified_end_time_secs"].(int)))
	}
	if modelMap["last_modified_start_time_secs"] != nil {
		model.LastModifiedStartTimeSecs = core.Int64Ptr(int64(modelMap["last_modified_start_time_secs"].(int)))
	}
	if modelMap["last_name"] != nil && modelMap["last_name"].(string) != "" {
		model.LastName = core.StringPtr(modelMap["last_name"].(string))
	}
	if modelMap["middle_name"] != nil && modelMap["middle_name"].(string) != "" {
		model.MiddleName = core.StringPtr(modelMap["middle_name"].(string))
	}
	if modelMap["organizer_address"] != nil && modelMap["organizer_address"].(string) != "" {
		model.OrganizerAddress = core.StringPtr(modelMap["organizer_address"].(string))
	}
	if modelMap["received_end_time_secs"] != nil {
		model.ReceivedEndTimeSecs = core.Int64Ptr(int64(modelMap["received_end_time_secs"].(int)))
	}
	if modelMap["received_start_time_secs"] != nil {
		model.ReceivedStartTimeSecs = core.Int64Ptr(int64(modelMap["received_start_time_secs"].(int)))
	}
	if modelMap["recipient_addresses"] != nil {
		recipientAddresses := []string{}
		for _, recipientAddressesItem := range modelMap["recipient_addresses"].([]interface{}) {
			recipientAddresses = append(recipientAddresses, recipientAddressesItem.(string))
		}
		model.RecipientAddresses = recipientAddresses
	}
	if modelMap["sender_address"] != nil && modelMap["sender_address"].(string) != "" {
		model.SenderAddress = core.StringPtr(modelMap["sender_address"].(string))
	}
	if modelMap["source_environment"] != nil && modelMap["source_environment"].(string) != "" {
		model.SourceEnvironment = core.StringPtr(modelMap["source_environment"].(string))
	}
	if modelMap["task_status_types"] != nil {
		taskStatusTypes := []string{}
		for _, taskStatusTypesItem := range modelMap["task_status_types"].([]interface{}) {
			taskStatusTypes = append(taskStatusTypes, taskStatusTypesItem.(string))
		}
		model.TaskStatusTypes = taskStatusTypes
	}
	if modelMap["types"] != nil {
		types := []string{}
		for _, typesItem := range modelMap["types"].([]interface{}) {
			types = append(types, typesItem.(string))
		}
		model.Types = types
	}
	if modelMap["o365_params"] != nil && len(modelMap["o365_params"].([]interface{})) > 0 {
		O365ParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365SearchEmailsRequestParams(modelMap["o365_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.O365Params = O365ParamsModel
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365SearchEmailsRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.O365SearchEmailsRequestParams, error) {
	model := &backuprecoveryv1.O365SearchEmailsRequestParams{}
	if modelMap["domain_ids"] != nil {
		domainIds := []int64{}
		for _, domainIdsItem := range modelMap["domain_ids"].([]interface{}) {
			domainIds = append(domainIds, int64(domainIdsItem.(int)))
		}
		model.DomainIds = domainIds
	}
	if modelMap["mailbox_ids"] != nil {
		mailboxIds := []int64{}
		for _, mailboxIdsItem := range modelMap["mailbox_ids"].([]interface{}) {
			mailboxIds = append(mailboxIds, int64(mailboxIdsItem.(int)))
		}
		model.MailboxIds = mailboxIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchExchangeObjectsRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchExchangeObjectsRequestParams, error) {
	model := &backuprecoveryv1.SearchExchangeObjectsRequestParams{}
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchFileRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchFileRequestParams, error) {
	model := &backuprecoveryv1.SearchFileRequestParams{}
	if modelMap["search_string"] != nil && modelMap["search_string"].(string) != "" {
		model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	}
	if modelMap["types"] != nil {
		types := []string{}
		for _, typesItem := range modelMap["types"].([]interface{}) {
			types = append(types, typesItem.(string))
		}
		model.Types = types
	}
	if modelMap["source_environments"] != nil {
		sourceEnvironments := []string{}
		for _, sourceEnvironmentsItem := range modelMap["source_environments"].([]interface{}) {
			sourceEnvironments = append(sourceEnvironments, sourceEnvironmentsItem.(string))
		}
		model.SourceEnvironments = sourceEnvironments
	}
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	if modelMap["object_ids"] != nil {
		objectIds := []int64{}
		for _, objectIdsItem := range modelMap["object_ids"].([]interface{}) {
			objectIds = append(objectIds, int64(objectIdsItem.(int)))
		}
		model.ObjectIds = objectIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToHbaseOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.HbaseOnPremSearchParams, error) {
	model := &backuprecoveryv1.HbaseOnPremSearchParams{}
	hbaseObjectTypes := []string{}
	for _, hbaseObjectTypesItem := range modelMap["hbase_object_types"].([]interface{}) {
		hbaseObjectTypes = append(hbaseObjectTypes, hbaseObjectTypesItem.(string))
	}
	model.HbaseObjectTypes = hbaseObjectTypes
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToHDFSOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.HDFSOnPremSearchParams, error) {
	model := &backuprecoveryv1.HDFSOnPremSearchParams{}
	hdfsTypes := []string{}
	for _, hdfsTypesItem := range modelMap["hdfs_types"].([]interface{}) {
		hdfsTypes = append(hdfsTypes, hdfsTypesItem.(string))
	}
	model.HdfsTypes = hdfsTypes
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToHiveOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.HiveOnPremSearchParams, error) {
	model := &backuprecoveryv1.HiveOnPremSearchParams{}
	hiveObjectTypes := []string{}
	for _, hiveObjectTypesItem := range modelMap["hive_object_types"].([]interface{}) {
		hiveObjectTypes = append(hiveObjectTypes, hiveObjectTypesItem.(string))
	}
	model.HiveObjectTypes = hiveObjectTypes
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToMongoDbOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.MongoDbOnPremSearchParams, error) {
	model := &backuprecoveryv1.MongoDbOnPremSearchParams{}
	mongoDbObjectTypes := []string{}
	for _, mongoDbObjectTypesItem := range modelMap["mongo_db_object_types"].([]interface{}) {
		mongoDbObjectTypes = append(mongoDbObjectTypes, mongoDbObjectTypesItem.(string))
	}
	model.MongoDBObjectTypes = mongoDbObjectTypes
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchMsGroupsRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchMsGroupsRequestParams, error) {
	model := &backuprecoveryv1.SearchMsGroupsRequestParams{}
	if modelMap["mailbox_params"] != nil && len(modelMap["mailbox_params"].([]interface{})) > 0 {
		MailboxParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchEmailRequestParamsBase(modelMap["mailbox_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MailboxParams = MailboxParamsModel
	}
	if modelMap["o365_params"] != nil && len(modelMap["o365_params"].([]interface{})) > 0 {
		O365ParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365SearchRequestParams(modelMap["o365_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.O365Params = O365ParamsModel
	}
	if modelMap["site_params"] != nil && len(modelMap["site_params"].([]interface{})) > 0 {
		SiteParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchDocumentLibraryRequestParams(modelMap["site_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.SiteParams = SiteParamsModel
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchEmailRequestParamsBase(modelMap map[string]interface{}) (*backuprecoveryv1.SearchEmailRequestParamsBase, error) {
	model := &backuprecoveryv1.SearchEmailRequestParamsBase{}
	if modelMap["attendees_addresses"] != nil {
		attendeesAddresses := []string{}
		for _, attendeesAddressesItem := range modelMap["attendees_addresses"].([]interface{}) {
			attendeesAddresses = append(attendeesAddresses, attendeesAddressesItem.(string))
		}
		model.AttendeesAddresses = attendeesAddresses
	}
	if modelMap["bcc_recipient_addresses"] != nil {
		bccRecipientAddresses := []string{}
		for _, bccRecipientAddressesItem := range modelMap["bcc_recipient_addresses"].([]interface{}) {
			bccRecipientAddresses = append(bccRecipientAddresses, bccRecipientAddressesItem.(string))
		}
		model.BccRecipientAddresses = bccRecipientAddresses
	}
	if modelMap["cc_recipient_addresses"] != nil {
		ccRecipientAddresses := []string{}
		for _, ccRecipientAddressesItem := range modelMap["cc_recipient_addresses"].([]interface{}) {
			ccRecipientAddresses = append(ccRecipientAddresses, ccRecipientAddressesItem.(string))
		}
		model.CcRecipientAddresses = ccRecipientAddresses
	}
	if modelMap["created_end_time_secs"] != nil {
		model.CreatedEndTimeSecs = core.Int64Ptr(int64(modelMap["created_end_time_secs"].(int)))
	}
	if modelMap["created_start_time_secs"] != nil {
		model.CreatedStartTimeSecs = core.Int64Ptr(int64(modelMap["created_start_time_secs"].(int)))
	}
	if modelMap["due_date_end_time_secs"] != nil {
		model.DueDateEndTimeSecs = core.Int64Ptr(int64(modelMap["due_date_end_time_secs"].(int)))
	}
	if modelMap["due_date_start_time_secs"] != nil {
		model.DueDateStartTimeSecs = core.Int64Ptr(int64(modelMap["due_date_start_time_secs"].(int)))
	}
	if modelMap["email_address"] != nil && modelMap["email_address"].(string) != "" {
		model.EmailAddress = core.StringPtr(modelMap["email_address"].(string))
	}
	if modelMap["email_subject"] != nil && modelMap["email_subject"].(string) != "" {
		model.EmailSubject = core.StringPtr(modelMap["email_subject"].(string))
	}
	if modelMap["first_name"] != nil && modelMap["first_name"].(string) != "" {
		model.FirstName = core.StringPtr(modelMap["first_name"].(string))
	}
	if modelMap["folder_names"] != nil {
		folderNames := []string{}
		for _, folderNamesItem := range modelMap["folder_names"].([]interface{}) {
			folderNames = append(folderNames, folderNamesItem.(string))
		}
		model.FolderNames = folderNames
	}
	if modelMap["has_attachment"] != nil {
		model.HasAttachment = core.BoolPtr(modelMap["has_attachment"].(bool))
	}
	if modelMap["last_modified_end_time_secs"] != nil {
		model.LastModifiedEndTimeSecs = core.Int64Ptr(int64(modelMap["last_modified_end_time_secs"].(int)))
	}
	if modelMap["last_modified_start_time_secs"] != nil {
		model.LastModifiedStartTimeSecs = core.Int64Ptr(int64(modelMap["last_modified_start_time_secs"].(int)))
	}
	if modelMap["last_name"] != nil && modelMap["last_name"].(string) != "" {
		model.LastName = core.StringPtr(modelMap["last_name"].(string))
	}
	if modelMap["middle_name"] != nil && modelMap["middle_name"].(string) != "" {
		model.MiddleName = core.StringPtr(modelMap["middle_name"].(string))
	}
	if modelMap["organizer_address"] != nil && modelMap["organizer_address"].(string) != "" {
		model.OrganizerAddress = core.StringPtr(modelMap["organizer_address"].(string))
	}
	if modelMap["received_end_time_secs"] != nil {
		model.ReceivedEndTimeSecs = core.Int64Ptr(int64(modelMap["received_end_time_secs"].(int)))
	}
	if modelMap["received_start_time_secs"] != nil {
		model.ReceivedStartTimeSecs = core.Int64Ptr(int64(modelMap["received_start_time_secs"].(int)))
	}
	if modelMap["recipient_addresses"] != nil {
		recipientAddresses := []string{}
		for _, recipientAddressesItem := range modelMap["recipient_addresses"].([]interface{}) {
			recipientAddresses = append(recipientAddresses, recipientAddressesItem.(string))
		}
		model.RecipientAddresses = recipientAddresses
	}
	if modelMap["sender_address"] != nil && modelMap["sender_address"].(string) != "" {
		model.SenderAddress = core.StringPtr(modelMap["sender_address"].(string))
	}
	if modelMap["source_environment"] != nil && modelMap["source_environment"].(string) != "" {
		model.SourceEnvironment = core.StringPtr(modelMap["source_environment"].(string))
	}
	if modelMap["task_status_types"] != nil {
		taskStatusTypes := []string{}
		for _, taskStatusTypesItem := range modelMap["task_status_types"].([]interface{}) {
			taskStatusTypes = append(taskStatusTypes, taskStatusTypesItem.(string))
		}
		model.TaskStatusTypes = taskStatusTypes
	}
	if modelMap["types"] != nil {
		types := []string{}
		for _, typesItem := range modelMap["types"].([]interface{}) {
			types = append(types, typesItem.(string))
		}
		model.Types = types
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365SearchRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.O365SearchRequestParams, error) {
	model := &backuprecoveryv1.O365SearchRequestParams{}
	if modelMap["domain_ids"] != nil {
		domainIds := []int64{}
		for _, domainIdsItem := range modelMap["domain_ids"].([]interface{}) {
			domainIds = append(domainIds, int64(domainIdsItem.(int)))
		}
		model.DomainIds = domainIds
	}
	if modelMap["group_ids"] != nil {
		groupIds := []int64{}
		for _, groupIdsItem := range modelMap["group_ids"].([]interface{}) {
			groupIds = append(groupIds, int64(groupIdsItem.(int)))
		}
		model.GroupIds = groupIds
	}
	if modelMap["site_ids"] != nil {
		siteIds := []int64{}
		for _, siteIdsItem := range modelMap["site_ids"].([]interface{}) {
			siteIds = append(siteIds, int64(siteIdsItem.(int)))
		}
		model.SiteIds = siteIds
	}
	if modelMap["teams_ids"] != nil {
		teamsIds := []int64{}
		for _, teamsIdsItem := range modelMap["teams_ids"].([]interface{}) {
			teamsIds = append(teamsIds, int64(teamsIdsItem.(int)))
		}
		model.TeamsIds = teamsIds
	}
	if modelMap["user_ids"] != nil {
		userIds := []int64{}
		for _, userIdsItem := range modelMap["user_ids"].([]interface{}) {
			userIds = append(userIds, int64(userIdsItem.(int)))
		}
		model.UserIds = userIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchDocumentLibraryRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchDocumentLibraryRequestParams, error) {
	model := &backuprecoveryv1.SearchDocumentLibraryRequestParams{}
	if modelMap["category_types"] != nil {
		categoryTypes := []string{}
		for _, categoryTypesItem := range modelMap["category_types"].([]interface{}) {
			categoryTypes = append(categoryTypes, categoryTypesItem.(string))
		}
		model.CategoryTypes = categoryTypes
	}
	if modelMap["creation_end_time_secs"] != nil {
		model.CreationEndTimeSecs = core.Int64Ptr(int64(modelMap["creation_end_time_secs"].(int)))
	}
	if modelMap["creation_start_time_secs"] != nil {
		model.CreationStartTimeSecs = core.Int64Ptr(int64(modelMap["creation_start_time_secs"].(int)))
	}
	if modelMap["include_files"] != nil {
		model.IncludeFiles = core.BoolPtr(modelMap["include_files"].(bool))
	}
	if modelMap["include_folders"] != nil {
		model.IncludeFolders = core.BoolPtr(modelMap["include_folders"].(bool))
	}
	if modelMap["o365_params"] != nil && len(modelMap["o365_params"].([]interface{})) > 0 {
		O365ParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365SearchRequestParams(modelMap["o365_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.O365Params = O365ParamsModel
	}
	if modelMap["owner_names"] != nil {
		ownerNames := []string{}
		for _, ownerNamesItem := range modelMap["owner_names"].([]interface{}) {
			ownerNames = append(ownerNames, ownerNamesItem.(string))
		}
		model.OwnerNames = ownerNames
	}
	if modelMap["search_string"] != nil && modelMap["search_string"].(string) != "" {
		model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	}
	if modelMap["size_bytes_lower_limit"] != nil {
		model.SizeBytesLowerLimit = core.Int64Ptr(int64(modelMap["size_bytes_lower_limit"].(int)))
	}
	if modelMap["size_bytes_upper_limit"] != nil {
		model.SizeBytesUpperLimit = core.Int64Ptr(int64(modelMap["size_bytes_upper_limit"].(int)))
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchMsTeamsRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchMsTeamsRequestParams, error) {
	model := &backuprecoveryv1.SearchMsTeamsRequestParams{}
	if modelMap["category_types"] != nil {
		categoryTypes := []string{}
		for _, categoryTypesItem := range modelMap["category_types"].([]interface{}) {
			categoryTypes = append(categoryTypes, categoryTypesItem.(string))
		}
		model.CategoryTypes = categoryTypes
	}
	if modelMap["channel_names"] != nil {
		channelNames := []string{}
		for _, channelNamesItem := range modelMap["channel_names"].([]interface{}) {
			channelNames = append(channelNames, channelNamesItem.(string))
		}
		model.ChannelNames = channelNames
	}
	if modelMap["channel_params"] != nil && len(modelMap["channel_params"].([]interface{})) > 0 {
		ChannelParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365TeamsChannelsSearchRequestParams(modelMap["channel_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ChannelParams = ChannelParamsModel
	}
	if modelMap["creation_end_time_secs"] != nil {
		model.CreationEndTimeSecs = core.Int64Ptr(int64(modelMap["creation_end_time_secs"].(int)))
	}
	if modelMap["creation_start_time_secs"] != nil {
		model.CreationStartTimeSecs = core.Int64Ptr(int64(modelMap["creation_start_time_secs"].(int)))
	}
	if modelMap["o365_params"] != nil && len(modelMap["o365_params"].([]interface{})) > 0 {
		O365ParamsModel, err := DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365SearchRequestParams(modelMap["o365_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.O365Params = O365ParamsModel
	}
	if modelMap["owner_names"] != nil {
		ownerNames := []string{}
		for _, ownerNamesItem := range modelMap["owner_names"].([]interface{}) {
			ownerNames = append(ownerNames, ownerNamesItem.(string))
		}
		model.OwnerNames = ownerNames
	}
	if modelMap["search_string"] != nil && modelMap["search_string"].(string) != "" {
		model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	}
	if modelMap["size_bytes_lower_limit"] != nil {
		model.SizeBytesLowerLimit = core.Int64Ptr(int64(modelMap["size_bytes_lower_limit"].(int)))
	}
	if modelMap["size_bytes_upper_limit"] != nil {
		model.SizeBytesUpperLimit = core.Int64Ptr(int64(modelMap["size_bytes_upper_limit"].(int)))
	}
	if modelMap["types"] != nil {
		types := []string{}
		for _, typesItem := range modelMap["types"].([]interface{}) {
			types = append(types, typesItem.(string))
		}
		model.Types = types
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToO365TeamsChannelsSearchRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.O365TeamsChannelsSearchRequestParams, error) {
	model := &backuprecoveryv1.O365TeamsChannelsSearchRequestParams{}
	if modelMap["channel_email"] != nil && modelMap["channel_email"].(string) != "" {
		model.ChannelEmail = core.StringPtr(modelMap["channel_email"].(string))
	}
	if modelMap["channel_id"] != nil && modelMap["channel_id"].(string) != "" {
		model.ChannelID = core.StringPtr(modelMap["channel_id"].(string))
	}
	if modelMap["channel_name"] != nil && modelMap["channel_name"].(string) != "" {
		model.ChannelName = core.StringPtr(modelMap["channel_name"].(string))
	}
	if modelMap["include_private_channels"] != nil {
		model.IncludePrivateChannels = core.BoolPtr(modelMap["include_private_channels"].(bool))
	}
	if modelMap["include_public_channels"] != nil {
		model.IncludePublicChannels = core.BoolPtr(modelMap["include_public_channels"].(bool))
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchPublicFolderRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchPublicFolderRequestParams, error) {
	model := &backuprecoveryv1.SearchPublicFolderRequestParams{}
	if modelMap["search_string"] != nil && modelMap["search_string"].(string) != "" {
		model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	}
	if modelMap["types"] != nil {
		types := []string{}
		for _, typesItem := range modelMap["types"].([]interface{}) {
			types = append(types, typesItem.(string))
		}
		model.Types = types
	}
	if modelMap["has_attachment"] != nil {
		model.HasAttachment = core.BoolPtr(modelMap["has_attachment"].(bool))
	}
	if modelMap["sender_address"] != nil && modelMap["sender_address"].(string) != "" {
		model.SenderAddress = core.StringPtr(modelMap["sender_address"].(string))
	}
	if modelMap["recipient_addresses"] != nil {
		recipientAddresses := []string{}
		for _, recipientAddressesItem := range modelMap["recipient_addresses"].([]interface{}) {
			recipientAddresses = append(recipientAddresses, recipientAddressesItem.(string))
		}
		model.RecipientAddresses = recipientAddresses
	}
	if modelMap["cc_recipient_addresses"] != nil {
		ccRecipientAddresses := []string{}
		for _, ccRecipientAddressesItem := range modelMap["cc_recipient_addresses"].([]interface{}) {
			ccRecipientAddresses = append(ccRecipientAddresses, ccRecipientAddressesItem.(string))
		}
		model.CcRecipientAddresses = ccRecipientAddresses
	}
	if modelMap["bcc_recipient_addresses"] != nil {
		bccRecipientAddresses := []string{}
		for _, bccRecipientAddressesItem := range modelMap["bcc_recipient_addresses"].([]interface{}) {
			bccRecipientAddresses = append(bccRecipientAddresses, bccRecipientAddressesItem.(string))
		}
		model.BccRecipientAddresses = bccRecipientAddresses
	}
	if modelMap["received_start_time_secs"] != nil {
		model.ReceivedStartTimeSecs = core.Int64Ptr(int64(modelMap["received_start_time_secs"].(int)))
	}
	if modelMap["received_end_time_secs"] != nil {
		model.ReceivedEndTimeSecs = core.Int64Ptr(int64(modelMap["received_end_time_secs"].(int)))
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToSearchSfdcRecordsRequestParams(modelMap map[string]interface{}) (*backuprecoveryv1.SearchSfdcRecordsRequestParams, error) {
	model := &backuprecoveryv1.SearchSfdcRecordsRequestParams{}
	mutationTypes := []string{}
	for _, mutationTypesItem := range modelMap["mutation_types"].([]interface{}) {
		mutationTypes = append(mutationTypes, mutationTypesItem.(string))
	}
	model.MutationTypes = mutationTypes
	model.ObjectName = core.StringPtr(modelMap["object_name"].(string))
	if modelMap["query_string"] != nil && modelMap["query_string"].(string) != "" {
		model.QueryString = core.StringPtr(modelMap["query_string"].(string))
	}
	model.SnapshotID = core.StringPtr(modelMap["snapshot_id"].(string))
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMapToUdaOnPremSearchParams(modelMap map[string]interface{}) (*backuprecoveryv1.UdaOnPremSearchParams, error) {
	model := &backuprecoveryv1.UdaOnPremSearchParams{}
	model.SearchString = core.StringPtr(modelMap["search_string"].(string))
	if modelMap["source_ids"] != nil {
		sourceIds := []int64{}
		for _, sourceIdsItem := range modelMap["source_ids"].([]interface{}) {
			sourceIds = append(sourceIds, int64(sourceIdsItem.(int)))
		}
		model.SourceIds = sourceIds
	}
	return model, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCassandraOnPremSearchParamsToMap(model *backuprecoveryv1.CassandraOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cassandra_object_types"] = model.CassandraObjectTypes
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCouchBaseOnPremSearchParamsToMap(model *backuprecoveryv1.CouchBaseOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["couchbase_object_types"] = model.CouchbaseObjectTypes
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchEmailRequestParamsToMap(model *backuprecoveryv1.SearchEmailRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AttendeesAddresses != nil {
		modelMap["attendees_addresses"] = model.AttendeesAddresses
	}
	if model.BccRecipientAddresses != nil {
		modelMap["bcc_recipient_addresses"] = model.BccRecipientAddresses
	}
	if model.CcRecipientAddresses != nil {
		modelMap["cc_recipient_addresses"] = model.CcRecipientAddresses
	}
	if model.CreatedEndTimeSecs != nil {
		modelMap["created_end_time_secs"] = flex.IntValue(model.CreatedEndTimeSecs)
	}
	if model.CreatedStartTimeSecs != nil {
		modelMap["created_start_time_secs"] = flex.IntValue(model.CreatedStartTimeSecs)
	}
	if model.DueDateEndTimeSecs != nil {
		modelMap["due_date_end_time_secs"] = flex.IntValue(model.DueDateEndTimeSecs)
	}
	if model.DueDateStartTimeSecs != nil {
		modelMap["due_date_start_time_secs"] = flex.IntValue(model.DueDateStartTimeSecs)
	}
	if model.EmailAddress != nil {
		modelMap["email_address"] = *model.EmailAddress
	}
	if model.EmailSubject != nil {
		modelMap["email_subject"] = *model.EmailSubject
	}
	if model.FirstName != nil {
		modelMap["first_name"] = *model.FirstName
	}
	if model.FolderNames != nil {
		modelMap["folder_names"] = model.FolderNames
	}
	if model.HasAttachment != nil {
		modelMap["has_attachment"] = *model.HasAttachment
	}
	if model.LastModifiedEndTimeSecs != nil {
		modelMap["last_modified_end_time_secs"] = flex.IntValue(model.LastModifiedEndTimeSecs)
	}
	if model.LastModifiedStartTimeSecs != nil {
		modelMap["last_modified_start_time_secs"] = flex.IntValue(model.LastModifiedStartTimeSecs)
	}
	if model.LastName != nil {
		modelMap["last_name"] = *model.LastName
	}
	if model.MiddleName != nil {
		modelMap["middle_name"] = *model.MiddleName
	}
	if model.OrganizerAddress != nil {
		modelMap["organizer_address"] = *model.OrganizerAddress
	}
	if model.ReceivedEndTimeSecs != nil {
		modelMap["received_end_time_secs"] = flex.IntValue(model.ReceivedEndTimeSecs)
	}
	if model.ReceivedStartTimeSecs != nil {
		modelMap["received_start_time_secs"] = flex.IntValue(model.ReceivedStartTimeSecs)
	}
	if model.RecipientAddresses != nil {
		modelMap["recipient_addresses"] = model.RecipientAddresses
	}
	if model.SenderAddress != nil {
		modelMap["sender_address"] = *model.SenderAddress
	}
	if model.SourceEnvironment != nil {
		modelMap["source_environment"] = *model.SourceEnvironment
	}
	if model.TaskStatusTypes != nil {
		modelMap["task_status_types"] = model.TaskStatusTypes
	}
	if model.Types != nil {
		modelMap["types"] = model.Types
	}
	if model.O365Params != nil {
		o365ParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectO365SearchEmailsRequestParamsToMap(model.O365Params)
		if err != nil {
			return modelMap, err
		}
		modelMap["o365_params"] = []map[string]interface{}{o365ParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectO365SearchEmailsRequestParamsToMap(model *backuprecoveryv1.O365SearchEmailsRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DomainIds != nil {
		modelMap["domain_ids"] = model.DomainIds
	}
	if model.MailboxIds != nil {
		modelMap["mailbox_ids"] = model.MailboxIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchExchangeObjectsRequestParamsToMap(model *backuprecoveryv1.SearchExchangeObjectsRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["search_string"] = *model.SearchString
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchFileRequestParamsToMap(model *backuprecoveryv1.SearchFileRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SearchString != nil {
		modelMap["search_string"] = *model.SearchString
	}
	if model.Types != nil {
		modelMap["types"] = model.Types
	}
	if model.SourceEnvironments != nil {
		modelMap["source_environments"] = model.SourceEnvironments
	}
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	if model.ObjectIds != nil {
		modelMap["object_ids"] = model.ObjectIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHbaseOnPremSearchParamsToMap(model *backuprecoveryv1.HbaseOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hbase_object_types"] = model.HbaseObjectTypes
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHDFSOnPremSearchParamsToMap(model *backuprecoveryv1.HDFSOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hdfs_types"] = model.HdfsTypes
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHiveOnPremSearchParamsToMap(model *backuprecoveryv1.HiveOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hive_object_types"] = model.HiveObjectTypes
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMongoDbOnPremSearchParamsToMap(model *backuprecoveryv1.MongoDbOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mongo_db_object_types"] = model.MongoDBObjectTypes
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchMsGroupsRequestParamsToMap(model *backuprecoveryv1.SearchMsGroupsRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.MailboxParams != nil {
		mailboxParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSearchEmailRequestParamsBaseToMap(model.MailboxParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mailbox_params"] = []map[string]interface{}{mailboxParamsMap}
	}
	if model.O365Params != nil {
		o365ParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectO365SearchRequestParamsToMap(model.O365Params)
		if err != nil {
			return modelMap, err
		}
		modelMap["o365_params"] = []map[string]interface{}{o365ParamsMap}
	}
	if model.SiteParams != nil {
		siteParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSearchDocumentLibraryRequestParamsToMap(model.SiteParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["site_params"] = []map[string]interface{}{siteParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchEmailRequestParamsBaseToMap(model *backuprecoveryv1.SearchEmailRequestParamsBase) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AttendeesAddresses != nil {
		modelMap["attendees_addresses"] = model.AttendeesAddresses
	}
	if model.BccRecipientAddresses != nil {
		modelMap["bcc_recipient_addresses"] = model.BccRecipientAddresses
	}
	if model.CcRecipientAddresses != nil {
		modelMap["cc_recipient_addresses"] = model.CcRecipientAddresses
	}
	if model.CreatedEndTimeSecs != nil {
		modelMap["created_end_time_secs"] = flex.IntValue(model.CreatedEndTimeSecs)
	}
	if model.CreatedStartTimeSecs != nil {
		modelMap["created_start_time_secs"] = flex.IntValue(model.CreatedStartTimeSecs)
	}
	if model.DueDateEndTimeSecs != nil {
		modelMap["due_date_end_time_secs"] = flex.IntValue(model.DueDateEndTimeSecs)
	}
	if model.DueDateStartTimeSecs != nil {
		modelMap["due_date_start_time_secs"] = flex.IntValue(model.DueDateStartTimeSecs)
	}
	if model.EmailAddress != nil {
		modelMap["email_address"] = *model.EmailAddress
	}
	if model.EmailSubject != nil {
		modelMap["email_subject"] = *model.EmailSubject
	}
	if model.FirstName != nil {
		modelMap["first_name"] = *model.FirstName
	}
	if model.FolderNames != nil {
		modelMap["folder_names"] = model.FolderNames
	}
	if model.HasAttachment != nil {
		modelMap["has_attachment"] = *model.HasAttachment
	}
	if model.LastModifiedEndTimeSecs != nil {
		modelMap["last_modified_end_time_secs"] = flex.IntValue(model.LastModifiedEndTimeSecs)
	}
	if model.LastModifiedStartTimeSecs != nil {
		modelMap["last_modified_start_time_secs"] = flex.IntValue(model.LastModifiedStartTimeSecs)
	}
	if model.LastName != nil {
		modelMap["last_name"] = *model.LastName
	}
	if model.MiddleName != nil {
		modelMap["middle_name"] = *model.MiddleName
	}
	if model.OrganizerAddress != nil {
		modelMap["organizer_address"] = *model.OrganizerAddress
	}
	if model.ReceivedEndTimeSecs != nil {
		modelMap["received_end_time_secs"] = flex.IntValue(model.ReceivedEndTimeSecs)
	}
	if model.ReceivedStartTimeSecs != nil {
		modelMap["received_start_time_secs"] = flex.IntValue(model.ReceivedStartTimeSecs)
	}
	if model.RecipientAddresses != nil {
		modelMap["recipient_addresses"] = model.RecipientAddresses
	}
	if model.SenderAddress != nil {
		modelMap["sender_address"] = *model.SenderAddress
	}
	if model.SourceEnvironment != nil {
		modelMap["source_environment"] = *model.SourceEnvironment
	}
	if model.TaskStatusTypes != nil {
		modelMap["task_status_types"] = model.TaskStatusTypes
	}
	if model.Types != nil {
		modelMap["types"] = model.Types
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectO365SearchRequestParamsToMap(model *backuprecoveryv1.O365SearchRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DomainIds != nil {
		modelMap["domain_ids"] = model.DomainIds
	}
	if model.GroupIds != nil {
		modelMap["group_ids"] = model.GroupIds
	}
	if model.SiteIds != nil {
		modelMap["site_ids"] = model.SiteIds
	}
	if model.TeamsIds != nil {
		modelMap["teams_ids"] = model.TeamsIds
	}
	if model.UserIds != nil {
		modelMap["user_ids"] = model.UserIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchDocumentLibraryRequestParamsToMap(model *backuprecoveryv1.SearchDocumentLibraryRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CategoryTypes != nil {
		modelMap["category_types"] = model.CategoryTypes
	}
	if model.CreationEndTimeSecs != nil {
		modelMap["creation_end_time_secs"] = flex.IntValue(model.CreationEndTimeSecs)
	}
	if model.CreationStartTimeSecs != nil {
		modelMap["creation_start_time_secs"] = flex.IntValue(model.CreationStartTimeSecs)
	}
	if model.IncludeFiles != nil {
		modelMap["include_files"] = *model.IncludeFiles
	}
	if model.IncludeFolders != nil {
		modelMap["include_folders"] = *model.IncludeFolders
	}
	if model.O365Params != nil {
		o365ParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectO365SearchRequestParamsToMap(model.O365Params)
		if err != nil {
			return modelMap, err
		}
		modelMap["o365_params"] = []map[string]interface{}{o365ParamsMap}
	}
	if model.OwnerNames != nil {
		modelMap["owner_names"] = model.OwnerNames
	}
	if model.SearchString != nil {
		modelMap["search_string"] = *model.SearchString
	}
	if model.SizeBytesLowerLimit != nil {
		modelMap["size_bytes_lower_limit"] = flex.IntValue(model.SizeBytesLowerLimit)
	}
	if model.SizeBytesUpperLimit != nil {
		modelMap["size_bytes_upper_limit"] = flex.IntValue(model.SizeBytesUpperLimit)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchMsTeamsRequestParamsToMap(model *backuprecoveryv1.SearchMsTeamsRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CategoryTypes != nil {
		modelMap["category_types"] = model.CategoryTypes
	}
	if model.ChannelNames != nil {
		modelMap["channel_names"] = model.ChannelNames
	}
	if model.ChannelParams != nil {
		channelParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectO365TeamsChannelsSearchRequestParamsToMap(model.ChannelParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["channel_params"] = []map[string]interface{}{channelParamsMap}
	}
	if model.CreationEndTimeSecs != nil {
		modelMap["creation_end_time_secs"] = flex.IntValue(model.CreationEndTimeSecs)
	}
	if model.CreationStartTimeSecs != nil {
		modelMap["creation_start_time_secs"] = flex.IntValue(model.CreationStartTimeSecs)
	}
	if model.O365Params != nil {
		o365ParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectO365SearchRequestParamsToMap(model.O365Params)
		if err != nil {
			return modelMap, err
		}
		modelMap["o365_params"] = []map[string]interface{}{o365ParamsMap}
	}
	if model.OwnerNames != nil {
		modelMap["owner_names"] = model.OwnerNames
	}
	if model.SearchString != nil {
		modelMap["search_string"] = *model.SearchString
	}
	if model.SizeBytesLowerLimit != nil {
		modelMap["size_bytes_lower_limit"] = flex.IntValue(model.SizeBytesLowerLimit)
	}
	if model.SizeBytesUpperLimit != nil {
		modelMap["size_bytes_upper_limit"] = flex.IntValue(model.SizeBytesUpperLimit)
	}
	if model.Types != nil {
		modelMap["types"] = model.Types
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectO365TeamsChannelsSearchRequestParamsToMap(model *backuprecoveryv1.O365TeamsChannelsSearchRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ChannelEmail != nil {
		modelMap["channel_email"] = *model.ChannelEmail
	}
	if model.ChannelID != nil {
		modelMap["channel_id"] = *model.ChannelID
	}
	if model.ChannelName != nil {
		modelMap["channel_name"] = *model.ChannelName
	}
	if model.IncludePrivateChannels != nil {
		modelMap["include_private_channels"] = *model.IncludePrivateChannels
	}
	if model.IncludePublicChannels != nil {
		modelMap["include_public_channels"] = *model.IncludePublicChannels
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchPublicFolderRequestParamsToMap(model *backuprecoveryv1.SearchPublicFolderRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SearchString != nil {
		modelMap["search_string"] = *model.SearchString
	}
	if model.Types != nil {
		modelMap["types"] = model.Types
	}
	if model.HasAttachment != nil {
		modelMap["has_attachment"] = *model.HasAttachment
	}
	if model.SenderAddress != nil {
		modelMap["sender_address"] = *model.SenderAddress
	}
	if model.RecipientAddresses != nil {
		modelMap["recipient_addresses"] = model.RecipientAddresses
	}
	if model.CcRecipientAddresses != nil {
		modelMap["cc_recipient_addresses"] = model.CcRecipientAddresses
	}
	if model.BccRecipientAddresses != nil {
		modelMap["bcc_recipient_addresses"] = model.BccRecipientAddresses
	}
	if model.ReceivedStartTimeSecs != nil {
		modelMap["received_start_time_secs"] = flex.IntValue(model.ReceivedStartTimeSecs)
	}
	if model.ReceivedEndTimeSecs != nil {
		modelMap["received_end_time_secs"] = flex.IntValue(model.ReceivedEndTimeSecs)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSearchSfdcRecordsRequestParamsToMap(model *backuprecoveryv1.SearchSfdcRecordsRequestParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mutation_types"] = model.MutationTypes
	modelMap["object_name"] = *model.ObjectName
	if model.QueryString != nil {
		modelMap["query_string"] = *model.QueryString
	}
	modelMap["snapshot_id"] = *model.SnapshotID
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectUdaOnPremSearchParamsToMap(model *backuprecoveryv1.UdaOnPremSearchParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["search_string"] = *model.SearchString
	if model.SourceIds != nil {
		modelMap["source_ids"] = model.SourceIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCassandraIndexedObjectToMap(model *backuprecoveryv1.CassandraIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCassandraIndexedObjectSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.KeyspaceType != nil {
		modelMap["keyspace_type"] = *model.KeyspaceType
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(model *backuprecoveryv1.TagInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["tag_id"] = *model.TagID
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(model *backuprecoveryv1.SnapshotTagInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["tag_id"] = *model.TagID
	if model.RunIds != nil {
		modelMap["run_ids"] = model.RunIds
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCassandraIndexedObjectSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model *backuprecoveryv1.SharepointObjectParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.SiteWebURL != nil {
		modelMap["site_web_url"] = *model.SiteWebURL
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(model *backuprecoveryv1.ObjectSummary) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model *backuprecoveryv1.ObjectTypeVCenterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IsCloudEnv != nil {
		modelMap["is_cloud_env"] = *model.IsCloudEnv
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model *backuprecoveryv1.ObjectTypeWindowsClusterParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ClusterSourceType != nil {
		modelMap["cluster_source_type"] = *model.ClusterSourceType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCouchbaseIndexedObjectToMap(model *backuprecoveryv1.CouchbaseIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCouchbaseIndexedObjectSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCouchbaseIndexedObjectSourceInfoToMap(model *backuprecoveryv1.CouchbaseIndexedObjectSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectEmailToMap(model *backuprecoveryv1.Email) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.BccRecipientAddresses != nil {
		modelMap["bcc_recipient_addresses"] = model.BccRecipientAddresses
	}
	if model.CcRecipientAddresses != nil {
		modelMap["cc_recipient_addresses"] = model.CcRecipientAddresses
	}
	if model.CreatedTimeSecs != nil {
		modelMap["created_time_secs"] = flex.IntValue(model.CreatedTimeSecs)
	}
	if model.DirectoryPath != nil {
		modelMap["directory_path"] = *model.DirectoryPath
	}
	if model.EmailAddresses != nil {
		modelMap["email_addresses"] = model.EmailAddresses
	}
	if model.EmailSubject != nil {
		modelMap["email_subject"] = *model.EmailSubject
	}
	if model.FirstName != nil {
		modelMap["first_name"] = *model.FirstName
	}
	if model.FolderName != nil {
		modelMap["folder_name"] = *model.FolderName
	}
	if model.HasAttachment != nil {
		modelMap["has_attachment"] = *model.HasAttachment
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.LastModificationName != nil {
		modelMap["last_modification_name"] = *model.LastModificationName
	}
	if model.LastModificationTimeSecs != nil {
		modelMap["last_modification_time_secs"] = flex.IntValue(model.LastModificationTimeSecs)
	}
	if model.LastName != nil {
		modelMap["last_name"] = *model.LastName
	}
	if model.OptionalAttendeesAddresses != nil {
		modelMap["optional_attendees_addresses"] = model.OptionalAttendeesAddresses
	}
	if model.OrganizerAddress != nil {
		modelMap["organizer_address"] = *model.OrganizerAddress
	}
	if model.ParentFolderID != nil {
		modelMap["parent_folder_id"] = flex.IntValue(model.ParentFolderID)
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.ReceivedTimeSecs != nil {
		modelMap["received_time_secs"] = flex.IntValue(model.ReceivedTimeSecs)
	}
	if model.RecipientAddresses != nil {
		modelMap["recipient_addresses"] = model.RecipientAddresses
	}
	if model.RequiredAttendeesAddresses != nil {
		modelMap["required_attendees_addresses"] = model.RequiredAttendeesAddresses
	}
	if model.SenderAddress != nil {
		modelMap["sender_address"] = *model.SenderAddress
	}
	if model.SentTimeSecs != nil {
		modelMap["sent_time_secs"] = flex.IntValue(model.SentTimeSecs)
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.TaskCompletionDateTimeSecs != nil {
		modelMap["task_completion_date_time_secs"] = flex.IntValue(model.TaskCompletionDateTimeSecs)
	}
	if model.TaskDueDateTimeSecs != nil {
		modelMap["task_due_date_time_secs"] = flex.IntValue(model.TaskDueDateTimeSecs)
	}
	if model.TaskStatus != nil {
		modelMap["task_status"] = *model.TaskStatus
	}
	if model.TenantID != nil {
		modelMap["tenant_id"] = *model.TenantID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.UserObjectInfo != nil {
		userObjectInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(model.UserObjectInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["user_object_info"] = []map[string]interface{}{userObjectInfoMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectExchangeIndexedObjectToMap(model *backuprecoveryv1.ExchangeIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectExchangeIndexedObjectSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.DatabaseName != nil {
		modelMap["database_name"] = *model.DatabaseName
	}
	if model.Email != nil {
		modelMap["email"] = *model.Email
	}
	if model.ObjectName != nil {
		modelMap["object_name"] = *model.ObjectName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectExchangeIndexedObjectSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectFileToMap(model *backuprecoveryv1.File) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectFileSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectFileSourceInfoToMap(model *backuprecoveryv1.FileSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	if model.ProtectionStats != nil {
		protectionStats := []map[string]interface{}{}
		for _, protectionStatsItem := range model.ProtectionStats {
			protectionStatsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectProtectionStatsSummaryToMap(&protectionStatsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			protectionStats = append(protectionStats, protectionStatsItemMap)
		}
		modelMap["protection_stats"] = protectionStats
	}
	if model.Permissions != nil {
		permissionsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectPermissionInfoToMap(model.Permissions)
		if err != nil {
			return modelMap, err
		}
		modelMap["permissions"] = []map[string]interface{}{permissionsMap}
	}
	if model.MssqlParams != nil {
		mssqlParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectMssqlParamsToMap(model.MssqlParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["mssql_params"] = []map[string]interface{}{mssqlParamsMap}
	}
	if model.PhysicalParams != nil {
		physicalParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectPhysicalParamsToMap(model.PhysicalParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["physical_params"] = []map[string]interface{}{physicalParamsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectObjectProtectionStatsSummaryToMap(model *backuprecoveryv1.ObjectProtectionStatsSummary) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchIndexedObjectPermissionInfoToMap(model *backuprecoveryv1.PermissionInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	if model.Users != nil {
		users := []map[string]interface{}{}
		for _, usersItem := range model.Users {
			usersItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectUserToMap(&usersItem) // #nosec G601
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
			groupsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectGroupToMap(&groupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			groups = append(groups, groupsItemMap)
		}
		modelMap["groups"] = groups
	}
	if model.Tenant != nil {
		tenantMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTenantToMap(model.Tenant)
		if err != nil {
			return modelMap, err
		}
		modelMap["tenant"] = []map[string]interface{}{tenantMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectUserToMap(model *backuprecoveryv1.User) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchIndexedObjectGroupToMap(model *backuprecoveryv1.Group) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchIndexedObjectTenantToMap(model *backuprecoveryv1.Tenant) (map[string]interface{}, error) {
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
		externalVendorMetadataMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectExternalVendorTenantMetadataToMap(model.ExternalVendorMetadata)
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
		networkMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTenantNetworkToMap(model.Network)
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

func DataSourceIbmBackupRecoverySearchIndexedObjectExternalVendorTenantMetadataToMap(model *backuprecoveryv1.ExternalVendorTenantMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.IbmTenantMetadataParams != nil {
		ibmTenantMetadataParamsMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectIbmTenantMetadataParamsToMap(model.IbmTenantMetadataParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["ibm_tenant_metadata_params"] = []map[string]interface{}{ibmTenantMetadataParamsMap}
	}
	modelMap["type"] = *model.Type
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectIbmTenantMetadataParamsToMap(model *backuprecoveryv1.IbmTenantMetadataParams) (map[string]interface{}, error) {
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
			customPropertiesItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectExternalVendorCustomPropertiesToMap(&customPropertiesItem) // #nosec G601
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
	if model.OwnershipMode != nil {
		modelMap["ownership_mode"] = *model.OwnershipMode
	}
	if model.ResourceGroupID != nil {
		modelMap["resource_group_id"] = *model.ResourceGroupID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectExternalVendorCustomPropertiesToMap(model *backuprecoveryv1.ExternalVendorCustomProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Key != nil {
		modelMap["key"] = *model.Key
	}
	if model.Value != nil {
		modelMap["value"] = *model.Value
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectTenantNetworkToMap(model *backuprecoveryv1.TenantNetwork) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchIndexedObjectObjectMssqlParamsToMap(model *backuprecoveryv1.FileSourceInfoMssqlParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AagInfo != nil {
		aagInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectAAGInfoToMap(model.AagInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["aag_info"] = []map[string]interface{}{aagInfoMap}
	}
	if model.HostInfo != nil {
		hostInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectHostInformationToMap(model.HostInfo)
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

func DataSourceIbmBackupRecoverySearchIndexedObjectAAGInfoToMap(model *backuprecoveryv1.AAGInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ObjectID != nil {
		modelMap["object_id"] = flex.IntValue(model.ObjectID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHostInformationToMap(model *backuprecoveryv1.HostInformation) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoverySearchIndexedObjectObjectPhysicalParamsToMap(model *backuprecoveryv1.FileSourceInfoPhysicalParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnableSystemBackup != nil {
		modelMap["enable_system_backup"] = *model.EnableSystemBackup
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHbaseIndexedObjectToMap(model *backuprecoveryv1.HbaseIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectHbaseIndexedObjectSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHbaseIndexedObjectSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHDFSIndexedObjectToMap(model *backuprecoveryv1.HDFSIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectHDFSIndexedObjectSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHDFSIndexedObjectSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectHiveIndexedObjectToMap(model *backuprecoveryv1.HiveIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCommonIndexedObjectParamsSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCommonIndexedObjectParamsSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMongoIndexedObjectToMap(model *backuprecoveryv1.MongoIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCommonIndexedObjectParamsSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.CdpInfo != nil {
		cdpInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCdpObjectInfoToMap(model.CdpInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["cdp_info"] = []map[string]interface{}{cdpInfoMap}
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCdpObjectInfoToMap(model *backuprecoveryv1.CdpObjectInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AllowReEnableCdp != nil {
		modelMap["allow_re_enable_cdp"] = *model.AllowReEnableCdp
	}
	if model.CdpEnabled != nil {
		modelMap["cdp_enabled"] = *model.CdpEnabled
	}
	if model.LastRunInfo != nil {
		lastRunInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCdpObjectLastRunInfoToMap(model.LastRunInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["last_run_info"] = []map[string]interface{}{lastRunInfoMap}
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCdpObjectLastRunInfoToMap(model *backuprecoveryv1.CdpObjectLastRunInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.LocalBackupInfo != nil {
		localBackupInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCdpLocalBackupInfoToMap(model.LocalBackupInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["local_backup_info"] = []map[string]interface{}{localBackupInfoMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectCdpLocalBackupInfoToMap(model *backuprecoveryv1.CdpLocalBackupInfo) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EndTimeInUsecs != nil {
		modelMap["end_time_in_usecs"] = flex.IntValue(model.EndTimeInUsecs)
	}
	if model.StartTimeInUsecs != nil {
		modelMap["start_time_in_usecs"] = flex.IntValue(model.StartTimeInUsecs)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectMsGroupItemToMap(model *backuprecoveryv1.MsGroupItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCommonIndexedObjectParamsSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.MailboxItem != nil {
		mailboxItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectEmailToMap(model.MailboxItem)
		if err != nil {
			return modelMap, err
		}
		modelMap["mailbox_item"] = []map[string]interface{}{mailboxItemMap}
	}
	if model.SiteItem != nil {
		siteItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectDocumentLibraryItemToMap(model.SiteItem)
		if err != nil {
			return modelMap, err
		}
		modelMap["site_item"] = []map[string]interface{}{siteItemMap}
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectDocumentLibraryItemToMap(model *backuprecoveryv1.DocumentLibraryItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectDocumentLibraryItemSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.CreationTimeSecs != nil {
		modelMap["creation_time_secs"] = flex.IntValue(model.CreationTimeSecs)
	}
	if model.FileType != nil {
		modelMap["file_type"] = *model.FileType
	}
	if model.ItemID != nil {
		modelMap["item_id"] = *model.ItemID
	}
	if model.ItemSize != nil {
		modelMap["item_size"] = flex.IntValue(model.ItemSize)
	}
	if model.OwnerEmail != nil {
		modelMap["owner_email"] = *model.OwnerEmail
	}
	if model.OwnerName != nil {
		modelMap["owner_name"] = *model.OwnerName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectDocumentLibraryItemSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectPublicFolderItemToMap(model *backuprecoveryv1.PublicFolderItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectPublicFolderItemSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Subject != nil {
		modelMap["subject"] = *model.Subject
	}
	if model.HasAttachments != nil {
		modelMap["has_attachments"] = *model.HasAttachments
	}
	if model.ItemClass != nil {
		modelMap["item_class"] = *model.ItemClass
	}
	if model.ReceivedTimeSecs != nil {
		modelMap["received_time_secs"] = flex.IntValue(model.ReceivedTimeSecs)
	}
	if model.ItemSize != nil {
		modelMap["item_size"] = flex.IntValue(model.ItemSize)
	}
	if model.ParentFolderID != nil {
		modelMap["parent_folder_id"] = *model.ParentFolderID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectPublicFolderItemSourceInfoToMap(model *backuprecoveryv1.CommonIndexedObjectParamsSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSfdcRecordsToMap(model *backuprecoveryv1.SfdcRecords) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ColumnNames != nil {
		modelMap["column_names"] = model.ColumnNames
	}
	if model.Records != nil {
		modelMap["records"] = model.Records
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectSharepointItemToMap(model *backuprecoveryv1.SharepointItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DocumentLibraryItem != nil {
		documentLibraryItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectDocumentLibraryItemToMap(model.DocumentLibraryItem)
		if err != nil {
			return modelMap, err
		}
		modelMap["document_library_item"] = []map[string]interface{}{documentLibraryItemMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectTeamsItemToMap(model *backuprecoveryv1.TeamsItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTeamsItemSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.ChannelItem != nil {
		channelItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectChannelItemToMap(model.ChannelItem)
		if err != nil {
			return modelMap, err
		}
		modelMap["channel_item"] = []map[string]interface{}{channelItemMap}
	}
	if model.FileItem != nil {
		fileItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTeamsFileItemToMap(model.FileItem)
		if err != nil {
			return modelMap, err
		}
		modelMap["file_item"] = []map[string]interface{}{fileItemMap}
	}
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectTeamsItemSourceInfoToMap(model *backuprecoveryv1.TeamsItemSourceInfo) (map[string]interface{}, error) {
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
		sharepointSiteSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSharepointObjectParamsToMap(model.SharepointSiteSummary)
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
			childObjectsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectSummaryToMap(&childObjectsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			childObjects = append(childObjects, childObjectsItemMap)
		}
		modelMap["child_objects"] = childObjects
	}
	if model.VCenterSummary != nil {
		vCenterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeVCenterParamsToMap(model.VCenterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["v_center_summary"] = []map[string]interface{}{vCenterSummaryMap}
	}
	if model.WindowsClusterSummary != nil {
		windowsClusterSummaryMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectObjectTypeWindowsClusterParamsToMap(model.WindowsClusterSummary)
		if err != nil {
			return modelMap, err
		}
		modelMap["windows_cluster_summary"] = []map[string]interface{}{windowsClusterSummaryMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectChannelItemToMap(model *backuprecoveryv1.ChannelItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ChannelEmail != nil {
		modelMap["channel_email"] = *model.ChannelEmail
	}
	if model.ChannelID != nil {
		modelMap["channel_id"] = *model.ChannelID
	}
	if model.ChannelName != nil {
		modelMap["channel_name"] = *model.ChannelName
	}
	if model.ChannelType != nil {
		modelMap["channel_type"] = *model.ChannelType
	}
	if model.CreationTimeSecs != nil {
		modelMap["creation_time_secs"] = flex.IntValue(model.CreationTimeSecs)
	}
	if model.OwnerNames != nil {
		modelMap["owner_names"] = model.OwnerNames
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectTeamsFileItemToMap(model *backuprecoveryv1.TeamsFileItem) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CreationTimeSecs != nil {
		modelMap["creation_time_secs"] = flex.IntValue(model.CreationTimeSecs)
	}
	if model.DriveName != nil {
		modelMap["drive_name"] = *model.DriveName
	}
	if model.FileType != nil {
		modelMap["file_type"] = *model.FileType
	}
	if model.ItemSize != nil {
		modelMap["item_size"] = flex.IntValue(model.ItemSize)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoverySearchIndexedObjectUdaIndexedObjectToMap(model *backuprecoveryv1.UdaIndexedObject) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tags != nil {
		tags := []map[string]interface{}{}
		for _, tagsItem := range model.Tags {
			tagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectTagInfoToMap(&tagsItem) // #nosec G601
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
			snapshotTagsItemMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectSnapshotTagInfoToMap(&snapshotTagsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			snapshotTags = append(snapshotTags, snapshotTagsItemMap)
		}
		modelMap["snapshot_tags"] = snapshotTags
	}
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	if model.ProtectionGroupID != nil {
		modelMap["protection_group_id"] = *model.ProtectionGroupID
	}
	if model.ProtectionGroupName != nil {
		modelMap["protection_group_name"] = *model.ProtectionGroupName
	}
	if model.PolicyID != nil {
		modelMap["policy_id"] = *model.PolicyID
	}
	if model.PolicyName != nil {
		modelMap["policy_name"] = *model.PolicyName
	}
	if model.StorageDomainID != nil {
		modelMap["storage_domain_id"] = flex.IntValue(model.StorageDomainID)
	}
	if model.SourceInfo != nil {
		sourceInfoMap, err := DataSourceIbmBackupRecoverySearchIndexedObjectCommonIndexedObjectParamsSourceInfoToMap(model.SourceInfo)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_info"] = []map[string]interface{}{sourceInfoMap}
	}
	if model.FullName != nil {
		modelMap["full_name"] = *model.FullName
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.ObjectType != nil {
		modelMap["object_type"] = *model.ObjectType
	}
	return modelMap, nil
}
