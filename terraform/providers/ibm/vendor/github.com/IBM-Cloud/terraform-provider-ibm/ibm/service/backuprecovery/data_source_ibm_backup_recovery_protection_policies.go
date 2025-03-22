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

func DataSourceIbmBackupRecoveryProtectionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmBackupRecoveryProtectionPoliciesRead,

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
				Description: "Filter policies by a list of policy ids.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"policy_names": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter policies by a list of policy names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"types": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Types specifies the policy type of policies to be returned.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude_linked_policies": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If excludeLinkedPolicies is set to true then only local policies created on cluster will be returned. The result will exclude all linked policies created from policy templates.",
			},
			"include_replicated_policies": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If includeReplicatedPolicies is set to true, then response will also contain replicated policies. By default, replication policies are not included in the response.",
			},
			"include_stats": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If includeStats is set to true, then response will return number of protection groups and objects. By default, the protection stats are not included in the response.",
			},
			"policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Specifies a list of protection policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the name of the Protection Policy.",
						},
						"backup_policy": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the backup schedule and retentions of a Protection Policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"regular": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the Incremental and Full policy settings and also the common Retention policy settings.\".",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"incremental": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies incremental backup settings for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that defines how frequent backup will be performed for a Protection Group.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies how often to start new runs of a Protection Group. <br>'Minutes' specifies that Protection Group run starts periodically after certain number of minutes specified in 'frequency' field. <br>'Hours' specifies that Protection Group run starts periodically after certain number of hours specified in 'frequency' field. <br>'Days' specifies that Protection Group run starts periodically after certain number of days specified in 'frequency' field. <br>'Week' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Month' specifies that new Protection Group runs start monthly on certain day of specific week. This schedule needs 'weekOfMonth' and 'dayOfWeek' fields to be set. <br> Example: To run the Protection Group on Second Sunday of Every Month, following schedule need to be set: <br> unit: 'Month' <br> dayOfWeek: 'Sunday' <br> weekOfMonth: 'Second'.",
																		},
																		"minute_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of minutes.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"frequency": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																					},
																				},
																			},
																		},
																		"hour_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of hours.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"frequency": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																					},
																				},
																			},
																		},
																		"day_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"frequency": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																					},
																				},
																			},
																		},
																		"week_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_week": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																				},
																			},
																		},
																		"month_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_week": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"week_of_month": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																					},
																					"day_of_month": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																					},
																				},
																			},
																		},
																		"year_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_year": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the day of the Year (such as 'First' or 'Last') in a Yearly Schedule. <br>This field is used to define the day in the year to start the Protection Group Run. <br> Example: if 'dayOfYear' is set to 'First', a backup is performed on the first day of every year. <br> Example: if 'dayOfYear' is set to 'Last', a backup is performed on the last day of every year.",
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
												"full": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies full backup settings for a Protection Group. Currently, full backup settings can be specified by using either of 'schedule' or 'schdulesAndRetentions' field. Using 'schdulesAndRetentions' is recommended when multiple full backups need to be configured. If full and incremental backup has common retention then only setting 'schedule' is recommended.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that defines how frequent full backup will be performed for a Protection Group.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies how often to start new runs of a Protection Group. <br>'Days' specifies that Protection Group run starts periodically on every day. For full backup schedule, currently we only support frequecny of 1 which indicates that full backup will be performed daily. <br>'Weeks' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Months' specifies that new Protection Group runs start monthly on certain day of specific week. This schedule needs 'weekOfMonth' and 'dayOfWeek' fields to be set. <br>'ProtectOnce' specifies that groups using this policy option will run only once and after that group will permanently be disabled. <br> Example: To run the Protection Group on Second Sunday of Every Month, following schedule need to be set: <br> unit: 'Month' <br> dayOfWeek: 'Sunday' <br> weekOfMonth: 'Second'.",
																		},
																		"day_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"frequency": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																					},
																				},
																			},
																		},
																		"week_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_week": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																				},
																			},
																		},
																		"month_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_week": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"week_of_month": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																					},
																					"day_of_month": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																					},
																				},
																			},
																		},
																		"year_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_year": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the day of the Year (such as 'First' or 'Last') in a Yearly Schedule. <br>This field is used to define the day in the year to start the Protection Group Run. <br> Example: if 'dayOfYear' is set to 'First', a backup is performed on the first day of every year. <br> Example: if 'dayOfYear' is set to 'Last', a backup is performed on the last day of every year.",
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
												"full_backups": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies multiple schedules and retentions for full backup. Specify either of the 'full' or 'fullBackups' values. Its recommended to use 'fullBaackups' value since 'full' will be deprecated after few releases.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that defines how frequent full backup will be performed for a Protection Group.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies how often to start new runs of a Protection Group. <br>'Days' specifies that Protection Group run starts periodically on every day. For full backup schedule, currently we only support frequecny of 1 which indicates that full backup will be performed daily. <br>'Weeks' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Months' specifies that new Protection Group runs start monthly on certain day of specific week. This schedule needs 'weekOfMonth' and 'dayOfWeek' fields to be set. <br>'ProtectOnce' specifies that groups using this policy option will run only once and after that group will permanently be disabled. <br> Example: To run the Protection Group on Second Sunday of Every Month, following schedule need to be set: <br> unit: 'Month' <br> dayOfWeek: 'Sunday' <br> weekOfMonth: 'Second'.",
																		},
																		"day_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"frequency": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																					},
																				},
																			},
																		},
																		"week_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_week": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																				},
																			},
																		},
																		"month_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_week": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																						Elem: &schema.Schema{
																							Type: schema.TypeString,
																						},
																					},
																					"week_of_month": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																					},
																					"day_of_month": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																					},
																				},
																			},
																		},
																		"year_schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"day_of_year": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the day of the Year (such as 'First' or 'Last') in a Yearly Schedule. <br>This field is used to define the day in the year to start the Protection Group Run. <br> Example: if 'dayOfYear' is set to 'First', a backup is performed on the first day of every year. <br> Example: if 'dayOfYear' is set to 'Last', a backup is performed on the last day of every year.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
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
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"primary_backup_target": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the primary backup target settings for regular backups. If the backup target field is not specified then backup will be taken locally on the Cohesity cluster.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"target_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the primary backup location where backups will be stored. If not specified, then default is assumed as local backup on Cohesity cluster.",
															},
															"archival_target_settings": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the primary archival settings. Mainly used for cloud direct archive (CAD) policy where primary backup is stored on archival target.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"target_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the Archival target id to take primary backup.",
																		},
																		"target_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the Archival target name where Snapshots are copied.",
																		},
																		"tier_settings": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the settings tier levels configured with each archival target. The tier settings need to be applied in specific order and default tier should always be passed as first entry in tiers array. The following example illustrates how to configure tiering input for AWS tiering. Same type of input structure applied to other cloud platforms also. <br>If user wants to achieve following tiering for backup, <br>User Desired Tiering- <br><t>1.Archive Full back up for 12 Months <br><t>2.Tier Levels <br><t><t>[1,12] [ <br><t><t><t>s3 (1 to 2 months), (default tier) <br><t><t><t>s3 Intelligent tiering (3 to 6 months), <br><t><t><t>s3 One Zone (7 to 9 months) <br><t><t><t>Glacier (10 to 12 months)] <br><t>API Input <br><t><t>1.tiers-[ <br><t><t><t>{'tierType': 'S3','moveAfterUnit':'months', <br><t><t><t>'moveAfter':2 - move from s3 to s3Inte after 2 months}, <br><t><t><t>{'tierType': 'S3Inte','moveAfterUnit':'months', <br><t><t><t>'moveAfter':4 - move from S3Inte to Glacier after 4 months}, <br><t><t><t>{'tierType': 'Glacier', 'moveAfterUnit':'months', <br><t><t><t>'moveAfter': 3 - move from Glacier to S3 One Zone after 3 months }, <br><t><t><t>{'tierType': 'S3 One Zone', 'moveAfterUnit': nil, <br><t><t><t>'moveAfter': nil - For the last record, 'moveAfter' and 'moveAfterUnit' <br><t><t><t>will be ignored since there are no further tier for data movement } <br><t><t><t>}].",
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
																				},
																			},
																		},
																	},
																},
															},
															"use_default_backup_target": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if the default primary backup target must be used for backups. If this is not specified or set to false, then targets specified in 'archivalTargetSettings' will be used for backups. If the value is specified as true, then default backup target is used internally. This field should only be set in the environment where tenant policy management is enabled and external targets are assigned to tenant when provisioning tenants.",
															},
														},
													},
												},
											},
										},
									},
									"log": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies log backup settings for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings that defines how frequent log backup will be performed for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies how often to start new Protection Group Runs of a Protection Group. <br>'Minutes' specifies that Protection Group run starts periodically after certain number of minutes specified in 'frequency' field. <br>'Hours' specifies that Protection Group run starts periodically after certain number of hours specified in 'frequency' field.",
															},
															"minute_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of minutes.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"hour_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of hours.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
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
									"bmr": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the BMR schedule in case of physical source protection.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings that defines how frequent bmr backup will be performed for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies how often to start new runs of a Protection Group. <br>'Weeks' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Months' specifies that new Protection Group runs start monthly on certain day of specific week.",
															},
															"day_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"week_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"month_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"week_of_month": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																		},
																		"day_of_month": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																		},
																	},
																},
															},
															"year_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_year": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the day of the Year (such as 'First' or 'Last') in a Yearly Schedule. <br>This field is used to define the day in the year to start the Protection Group Run. <br> Example: if 'dayOfYear' is set to 'First', a backup is performed on the first day of every year. <br> Example: if 'dayOfYear' is set to 'Last', a backup is performed on the last day of every year.",
																		},
																	},
																},
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
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
									"cdp": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies CDP (Continious Data Protection) backup settings for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a CDP backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a CDP backup measured in minutes or hours.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a cdp backup retention.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
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
									"storage_array_snapshot": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies storage snapshot managment backup settings for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies settings that defines how frequent Storage Snapshot Management backup will be performed for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies how often to start new Protection Group Runs of a Protection Group. <br>'Minutes' specifies that Protection Group run starts periodically after certain number of minutes specified in 'frequency' field. <br>'Hours' specifies that Protection Group run starts periodically after certain number of hours specified in 'frequency' field.",
															},
															"minute_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of minutes.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"hour_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of hours.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"day_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"week_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																	},
																},
															},
															"month_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"week_of_month": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																		},
																		"day_of_month": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																		},
																	},
																},
															},
															"year_schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_year": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the day of the Year (such as 'First' or 'Last') in a Yearly Schedule. <br>This field is used to define the day in the year to start the Protection Group Run. <br> Example: if 'dayOfYear' is set to 'First', a backup is performed on the first day of every year. <br> Example: if 'dayOfYear' is set to 'Last', a backup is performed on the last day of every year.",
																		},
																	},
																},
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
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
									"run_timeouts": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the backup timeouts for different type of runs(kFull, kRegular etc.).",
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
								},
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the description of the Protection Policy.",
						},
						"blackout_window": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of Blackout Windows. If specified, this field defines blackout periods when new Group Runs are not started. If a Group Run has been scheduled but not yet executed and the blackout period starts, the behavior depends on the policy field AbortInBlackoutPeriod.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"day": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies a day in the week when no new Protection Group Runs should be started such as 'Sunday'. Specifies a day in a week such as 'Sunday', 'Monday', etc.",
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
									"end_time": &schema.Schema{
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
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed olny when policies are updated.",
									},
								},
							},
						},
						"extended_retention": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies additional retention policies that should be applied to the backup snapshots. A backup snapshot will be retained up to a time that is the maximum of all retention policies that are applicable to it.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies a schedule frequency and schedule unit for Extended Retentions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unit interval for retention of Snapshots. <br>'Runs' means that the Snapshot copy retained after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy retained hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy gets retained daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy is retained weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy is retained monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy is retained yearly at the frequency set in the Frequency.",
												},
												"frequency": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies a factor to multiply the unit by, to determine the retention schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is retained.",
												},
											},
										},
									},
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
															},
														},
													},
												},
											},
										},
									},
									"run_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The backup run type to which this extended retention applies to. If this is not set, the extended retention will be applicable to all non-log backup types. Currently, the only value that can be set here is Full.'Regular' indicates a incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a Regular schedule captures all the blocks.'Full' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized.'Log' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time.'System' indicates a system backup. System backups are used to do bare metal recovery of the system to a specific point in time.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed olny when policies are updated.",
									},
								},
							},
						},
						"remote_target_policy": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the replication, archival and cloud spin targets of Protection Policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replication_targets": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"copy_on_run_success": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
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
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the type of target to which replication need to be performed.",
												},
												"remote_target_config": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the configuration for adding remote cluster as repilcation target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the cluster id of the target replication cluster.",
															},
															"cluster_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the cluster name of the target replication cluster.",
															},
														},
													},
												},
											},
										},
									},
									"archival_targets": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"copy_on_run_success": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"target_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the Archival target to copy the Snapshots to.",
												},
												"target_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the Archival target name where Snapshots are copied.",
												},
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the Archival target type where Snapshots are copied.",
												},
												"tier_settings": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the settings tier levels configured with each archival target. The tier settings need to be applied in specific order and default tier should always be passed as first entry in tiers array. The following example illustrates how to configure tiering input for AWS tiering. Same type of input structure applied to other cloud platforms also. <br>If user wants to achieve following tiering for backup, <br>User Desired Tiering- <br><t>1.Archive Full back up for 12 Months <br><t>2.Tier Levels <br><t><t>[1,12] [ <br><t><t><t>s3 (1 to 2 months), (default tier) <br><t><t><t>s3 Intelligent tiering (3 to 6 months), <br><t><t><t>s3 One Zone (7 to 9 months) <br><t><t><t>Glacier (10 to 12 months)] <br><t>API Input <br><t><t>1.tiers-[ <br><t><t><t>{'tierType': 'S3','moveAfterUnit':'months', <br><t><t><t>'moveAfter':2 - move from s3 to s3Inte after 2 months}, <br><t><t><t>{'tierType': 'S3Inte','moveAfterUnit':'months', <br><t><t><t>'moveAfter':4 - move from S3Inte to Glacier after 4 months}, <br><t><t><t>{'tierType': 'Glacier', 'moveAfterUnit':'months', <br><t><t><t>'moveAfter': 3 - move from Glacier to S3 One Zone after 3 months }, <br><t><t><t>{'tierType': 'S3 One Zone', 'moveAfterUnit': nil, <br><t><t><t>'moveAfter': nil - For the last record, 'moveAfter' and 'moveAfterUnit' <br><t><t><t>will be ignored since there are no further tier for data movement } <br><t><t><t>}].",
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
														},
													},
												},
												"extended_retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies additional retention policies that should be applied to the archived backup. Archived backup snapshot will be retained up to a time that is the maximum of all retention policies that are applicable to it.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a schedule frequency and schedule unit for Extended Retentions.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unit interval for retention of Snapshots. <br>'Runs' means that the Snapshot copy retained after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy retained hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy gets retained daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy is retained weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy is retained monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy is retained yearly at the frequency set in the Frequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the retention schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is retained.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The backup run type to which this extended retention applies to. If this is not set, the extended retention will be applicable to all non-log backup types. Currently, the only value that can be set here is Full.'Regular' indicates a incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a Regular schedule captures all the blocks.'Full' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized.'Log' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time.'System' indicates a system backup. System backups are used to do bare metal recovery of the system to a specific point in time.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier for the target getting added. This field need to be passed olny when policies are updated.",
															},
														},
													},
												},
											},
										},
									},
									"cloud_spin_targets": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"copy_on_run_success": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"target": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the details about Cloud Spin target where backup snapshots may be converted and stored.",
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
														},
													},
												},
											},
										},
									},
									"onprem_deploy_targets": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"copy_on_run_success": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"params": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the details about OnpremDeploy target where backup snapshots may be converted and deployed.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the unique id of the onprem entity.",
															},
														},
													},
												},
											},
										},
									},
									"rpaas_targets": &schema.Schema{
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"copy_on_run_success": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																		},
																	},
																},
															},
														},
													},
												},
												"target_id": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Specifies the RPaaS target to copy the Snapshots.",
												},
												"target_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the RPaaS target name where Snapshots are copied.",
												},
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the RPaaS target type where Snapshots are copied.",
												},
											},
										},
									},
								},
							},
						},
						"cascaded_targets_config": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Specifies the configuration for cascaded replications. Using cascaded replication, replication cluster(Rx) can further replicate and archive the snapshot copies to further targets. Its recommended to create cascaded configuration where protection group will be created.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_cluster_id": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the source cluster id from where the remote operations will be performed to the next set of remote targets.",
									},
									"remote_targets": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Specifies the replication, archival and cloud spin targets of Protection Policy.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"replication_targets": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"copy_on_run_success": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
															},
															"backup_run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
															},
															"run_timeouts": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
															"log_retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
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
															"target_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the type of target to which replication need to be performed.",
															},
															"remote_target_config": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the configuration for adding remote cluster as repilcation target.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"cluster_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the cluster id of the target replication cluster.",
																		},
																		"cluster_name": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the cluster name of the target replication cluster.",
																		},
																	},
																},
															},
														},
													},
												},
												"archival_targets": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"copy_on_run_success": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
															},
															"backup_run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
															},
															"run_timeouts": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
															"log_retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"target_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the Archival target to copy the Snapshots to.",
															},
															"target_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Archival target name where Snapshots are copied.",
															},
															"target_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Archival target type where Snapshots are copied.",
															},
															"tier_settings": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the settings tier levels configured with each archival target. The tier settings need to be applied in specific order and default tier should always be passed as first entry in tiers array. The following example illustrates how to configure tiering input for AWS tiering. Same type of input structure applied to other cloud platforms also. <br>If user wants to achieve following tiering for backup, <br>User Desired Tiering- <br><t>1.Archive Full back up for 12 Months <br><t>2.Tier Levels <br><t><t>[1,12] [ <br><t><t><t>s3 (1 to 2 months), (default tier) <br><t><t><t>s3 Intelligent tiering (3 to 6 months), <br><t><t><t>s3 One Zone (7 to 9 months) <br><t><t><t>Glacier (10 to 12 months)] <br><t>API Input <br><t><t>1.tiers-[ <br><t><t><t>{'tierType': 'S3','moveAfterUnit':'months', <br><t><t><t>'moveAfter':2 - move from s3 to s3Inte after 2 months}, <br><t><t><t>{'tierType': 'S3Inte','moveAfterUnit':'months', <br><t><t><t>'moveAfter':4 - move from S3Inte to Glacier after 4 months}, <br><t><t><t>{'tierType': 'Glacier', 'moveAfterUnit':'months', <br><t><t><t>'moveAfter': 3 - move from Glacier to S3 One Zone after 3 months }, <br><t><t><t>{'tierType': 'S3 One Zone', 'moveAfterUnit': nil, <br><t><t><t>'moveAfter': nil - For the last record, 'moveAfter' and 'moveAfterUnit' <br><t><t><t>will be ignored since there are no further tier for data movement } <br><t><t><t>}].",
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
																	},
																},
															},
															"extended_retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies additional retention policies that should be applied to the archived backup. Archived backup snapshot will be retained up to a time that is the maximum of all retention policies that are applicable to it.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"schedule": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies a schedule frequency and schedule unit for Extended Retentions.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the unit interval for retention of Snapshots. <br>'Runs' means that the Snapshot copy retained after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy retained hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy gets retained daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy is retained weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy is retained monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy is retained yearly at the frequency set in the Frequency.",
																					},
																					"frequency": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies a factor to multiply the unit by, to determine the retention schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is retained.",
																					},
																				},
																			},
																		},
																		"retention": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies the retention of a backup.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																					},
																					"data_lock_config": &schema.Schema{
																						Type:        schema.TypeList,
																						Computed:    true,
																						Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																						Elem: &schema.Resource{
																							Schema: map[string]*schema.Schema{
																								"mode": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																								},
																								"unit": &schema.Schema{
																									Type:        schema.TypeString,
																									Computed:    true,
																									Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																								},
																								"duration": &schema.Schema{
																									Type:        schema.TypeInt,
																									Computed:    true,
																									Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																								},
																								"enable_worm_on_external_target": &schema.Schema{
																									Type:        schema.TypeBool,
																									Computed:    true,
																									Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																		"run_type": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The backup run type to which this extended retention applies to. If this is not set, the extended retention will be applicable to all non-log backup types. Currently, the only value that can be set here is Full.'Regular' indicates a incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a Regular schedule captures all the blocks.'Full' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized.'Log' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time.'System' indicates a system backup. System backups are used to do bare metal recovery of the system to a specific point in time.",
																		},
																		"config_id": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the unique identifier for the target getting added. This field need to be passed olny when policies are updated.",
																		},
																	},
																},
															},
														},
													},
												},
												"cloud_spin_targets": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"copy_on_run_success": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
															},
															"backup_run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
															},
															"run_timeouts": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
															"log_retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"target": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the details about Cloud Spin target where backup snapshots may be converted and stored.",
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
																	},
																},
															},
														},
													},
												},
												"onprem_deploy_targets": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"copy_on_run_success": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
															},
															"backup_run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
															},
															"run_timeouts": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
															"log_retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"params": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the details about OnpremDeploy target where backup snapshots may be converted and deployed.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the unique id of the onprem entity.",
																		},
																	},
																},
															},
														},
													},
												},
												"rpaas_targets": &schema.Schema{
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"copy_on_run_success": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
															},
															"backup_run_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
															},
															"run_timeouts": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
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
															"log_retention": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Computed:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Computed:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Computed:    true,
																						Description: "Specifies whether objects in the external target associated with this policy need to be made immutable.",
																					},
																				},
																			},
																		},
																	},
																},
															},
															"target_id": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Specifies the RPaaS target to copy the Snapshots.",
															},
															"target_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the RPaaS target name where Snapshots are copied.",
															},
															"target_type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the RPaaS target type where Snapshots are copied.",
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
						"retry_options": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Retry Options of a Protection Policy when a Protection Group run fails.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retries": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of times to retry capturing Snapshots before the Protection Group Run fails.",
									},
									"retry_interval_mins": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Specifies the number of minutes before retrying a failed Protection Group.",
									},
								},
							},
						},
						"data_lock": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "This field is now deprecated. Please use the DataLockConfig in the backup retention.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the current policy verison. Policy version is incremented for optionally supporting new features and differentialting across releases.",
						},
						"is_cbs_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Specifies true if Calender Based Schedule is supported by client. Default value is assumed as false for this feature.",
						},
						"last_modification_time_usecs": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the last time this Policy was updated. If this is passed into a PUT request, then the backend will validate that the timestamp passed in matches the time that the policy was actually last modified. If the two timestamps do not match, then the request will be rejected with a stale error.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies a unique Policy id assigned by the Cohesity Cluster.",
						},
						"template_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the parent policy template id to which the policy is linked to. This field is set only when policy is created from template.",
						},
						"is_usable": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field is set to true if the linked policy which is internally created from a policy templates qualifies as usable to create more policies on the cluster. If the linked policy is partially filled and can not create a working policy then this field will be set to false. In case of normal policy created on the cluster, this field wont be populated.",
						},
						"is_replicated": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This field is set to true when policy is the replicated policy.",
						},
						"num_protection_groups": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of protection groups using the protection policy.",
						},
						"num_protected_objects": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the number of protected objects using the protection policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmBackupRecoveryProtectionPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_policies", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProtectionPoliciesOptions := &backuprecoveryv1.GetProtectionPoliciesOptions{}

	getProtectionPoliciesOptions.SetXIBMTenantID(d.Get("x_ibm_tenant_id").(string))
	if _, ok := d.GetOk("request_initiator_type"); ok {
		getProtectionPoliciesOptions.SetRequestInitiatorType(d.Get("request_initiator_type").(string))
	}
	if _, ok := d.GetOk("ids"); ok {
		var ids []string
		for _, v := range d.Get("ids").([]interface{}) {
			idsItem := v.(string)
			ids = append(ids, idsItem)
		}
		getProtectionPoliciesOptions.SetIds(ids)
	}
	if _, ok := d.GetOk("policy_names"); ok {
		var policyNames []string
		for _, v := range d.Get("policy_names").([]interface{}) {
			policyNamesItem := v.(string)
			policyNames = append(policyNames, policyNamesItem)
		}
		getProtectionPoliciesOptions.SetPolicyNames(policyNames)
	}
	if _, ok := d.GetOk("types"); ok {
		var types []string
		for _, v := range d.Get("types").([]interface{}) {
			typesItem := v.(string)
			types = append(types, typesItem)
		}
		getProtectionPoliciesOptions.SetTypes(types)
	}
	if _, ok := d.GetOk("exclude_linked_policies"); ok {
		getProtectionPoliciesOptions.SetExcludeLinkedPolicies(d.Get("exclude_linked_policies").(bool))
	}
	if _, ok := d.GetOk("include_replicated_policies"); ok {
		getProtectionPoliciesOptions.SetIncludeReplicatedPolicies(d.Get("include_replicated_policies").(bool))
	}
	if _, ok := d.GetOk("include_stats"); ok {
		getProtectionPoliciesOptions.SetIncludeStats(d.Get("include_stats").(bool))
	}

	protectionPoliciesResponse, _, err := backupRecoveryClient.GetProtectionPoliciesWithContext(context, getProtectionPoliciesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProtectionPoliciesWithContext failed: %s", err.Error()), "(Data) ibm_backup_recovery_protection_policies", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIbmBackupRecoveryProtectionPoliciesID(d))

	if !core.IsNil(protectionPoliciesResponse.Policies) {
		policies := []map[string]interface{}{}
		for _, policiesItem := range protectionPoliciesResponse.Policies {
			policiesItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesProtectionPolicyResponseToMap(&policiesItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_backup_recovery_protection_policies", "read", "policies-to-map").GetDiag()
			}
			policies = append(policies, policiesItemMap)
		}
		if err = d.Set("policies", policies); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting policies: %s", err), "(Data) ibm_backup_recovery_protection_policies", "read", "set-policies").GetDiag()
		}
	}

	return nil
}

// dataSourceIbmBackupRecoveryProtectionPoliciesID returns a reasonable ID for the list.
func dataSourceIbmBackupRecoveryProtectionPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIbmBackupRecoveryProtectionPoliciesProtectionPolicyResponseToMap(model *backuprecoveryv1.ProtectionPolicyResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	backupPolicyMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesBackupPolicyToMap(model.BackupPolicy)
	if err != nil {
		return modelMap, err
	}
	modelMap["backup_policy"] = []map[string]interface{}{backupPolicyMap}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.BlackoutWindow != nil {
		blackoutWindow := []map[string]interface{}{}
		for _, blackoutWindowItem := range model.BlackoutWindow {
			blackoutWindowItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesBlackoutWindowToMap(&blackoutWindowItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			blackoutWindow = append(blackoutWindow, blackoutWindowItemMap)
		}
		modelMap["blackout_window"] = blackoutWindow
	}
	if model.ExtendedRetention != nil {
		extendedRetention := []map[string]interface{}{}
		for _, extendedRetentionItem := range model.ExtendedRetention {
			extendedRetentionItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesExtendedRetentionPolicyToMap(&extendedRetentionItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			extendedRetention = append(extendedRetention, extendedRetentionItemMap)
		}
		modelMap["extended_retention"] = extendedRetention
	}
	if model.RemoteTargetPolicy != nil {
		remoteTargetPolicyMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetsConfigurationToMap(model.RemoteTargetPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote_target_policy"] = []map[string]interface{}{remoteTargetPolicyMap}
	}
	if model.CascadedTargetsConfig != nil {
		cascadedTargetsConfig := []map[string]interface{}{}
		for _, cascadedTargetsConfigItem := range model.CascadedTargetsConfig {
			cascadedTargetsConfigItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCascadedTargetConfigurationToMap(&cascadedTargetsConfigItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cascadedTargetsConfig = append(cascadedTargetsConfig, cascadedTargetsConfigItemMap)
		}
		modelMap["cascaded_targets_config"] = cascadedTargetsConfig
	}
	if model.RetryOptions != nil {
		retryOptionsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetryOptionsToMap(model.RetryOptions)
		if err != nil {
			return modelMap, err
		}
		modelMap["retry_options"] = []map[string]interface{}{retryOptionsMap}
	}
	if model.DataLock != nil {
		modelMap["data_lock"] = *model.DataLock
	}
	if model.Version != nil {
		modelMap["version"] = flex.IntValue(model.Version)
	}
	if model.IsCBSEnabled != nil {
		modelMap["is_cbs_enabled"] = *model.IsCBSEnabled
	}
	if model.LastModificationTimeUsecs != nil {
		modelMap["last_modification_time_usecs"] = flex.IntValue(model.LastModificationTimeUsecs)
	}
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.TemplateID != nil {
		modelMap["template_id"] = *model.TemplateID
	}
	if model.IsUsable != nil {
		modelMap["is_usable"] = *model.IsUsable
	}
	if model.IsReplicated != nil {
		modelMap["is_replicated"] = *model.IsReplicated
	}
	if model.NumProtectionGroups != nil {
		modelMap["num_protection_groups"] = flex.IntValue(model.NumProtectionGroups)
	}
	if model.NumProtectedObjects != nil {
		modelMap["num_protected_objects"] = flex.IntValue(model.NumProtectedObjects)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesBackupPolicyToMap(model *backuprecoveryv1.BackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	regularMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRegularBackupPolicyToMap(model.Regular)
	if err != nil {
		return modelMap, err
	}
	modelMap["regular"] = []map[string]interface{}{regularMap}
	if model.Log != nil {
		logMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogBackupPolicyToMap(model.Log)
		if err != nil {
			return modelMap, err
		}
		modelMap["log"] = []map[string]interface{}{logMap}
	}
	if model.Bmr != nil {
		bmrMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesBmrBackupPolicyToMap(model.Bmr)
		if err != nil {
			return modelMap, err
		}
		modelMap["bmr"] = []map[string]interface{}{bmrMap}
	}
	if model.Cdp != nil {
		cdpMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCdpBackupPolicyToMap(model.Cdp)
		if err != nil {
			return modelMap, err
		}
		modelMap["cdp"] = []map[string]interface{}{cdpMap}
	}
	if model.StorageArraySnapshot != nil {
		storageArraySnapshotMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesStorageArraySnapshotBackupPolicyToMap(model.StorageArraySnapshot)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_array_snapshot"] = []map[string]interface{}{storageArraySnapshotMap}
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesRegularBackupPolicyToMap(model *backuprecoveryv1.RegularBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Incremental != nil {
		incrementalMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesIncrementalBackupPolicyToMap(model.Incremental)
		if err != nil {
			return modelMap, err
		}
		modelMap["incremental"] = []map[string]interface{}{incrementalMap}
	}
	if model.Full != nil {
		fullMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesFullBackupPolicyToMap(model.Full)
		if err != nil {
			return modelMap, err
		}
		modelMap["full"] = []map[string]interface{}{fullMap}
	}
	if model.FullBackups != nil {
		fullBackups := []map[string]interface{}{}
		for _, fullBackupsItem := range model.FullBackups {
			fullBackupsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesFullScheduleAndRetentionToMap(&fullBackupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			fullBackups = append(fullBackups, fullBackupsItemMap)
		}
		modelMap["full_backups"] = fullBackups
	}
	if model.Retention != nil {
		retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
		if err != nil {
			return modelMap, err
		}
		modelMap["retention"] = []map[string]interface{}{retentionMap}
	}
	if model.PrimaryBackupTarget != nil {
		primaryBackupTargetMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesPrimaryBackupTargetToMap(model.PrimaryBackupTarget)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_backup_target"] = []map[string]interface{}{primaryBackupTargetMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesIncrementalBackupPolicyToMap(model *backuprecoveryv1.IncrementalBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesIncrementalScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesIncrementalScheduleToMap(model *backuprecoveryv1.IncrementalSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.MinuteSchedule != nil {
		minuteScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMinuteScheduleToMap(model.MinuteSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["minute_schedule"] = []map[string]interface{}{minuteScheduleMap}
	}
	if model.HourSchedule != nil {
		hourScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesHourScheduleToMap(model.HourSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["hour_schedule"] = []map[string]interface{}{hourScheduleMap}
	}
	if model.DaySchedule != nil {
		dayScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesMinuteScheduleToMap(model *backuprecoveryv1.MinuteSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["frequency"] = flex.IntValue(model.Frequency)
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesHourScheduleToMap(model *backuprecoveryv1.HourSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["frequency"] = flex.IntValue(model.Frequency)
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesDayScheduleToMap(model *backuprecoveryv1.DaySchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["frequency"] = flex.IntValue(model.Frequency)
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesWeekScheduleToMap(model *backuprecoveryv1.WeekSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["day_of_week"] = model.DayOfWeek
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesMonthScheduleToMap(model *backuprecoveryv1.MonthSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.DayOfWeek != nil {
		modelMap["day_of_week"] = model.DayOfWeek
	}
	if model.WeekOfMonth != nil {
		modelMap["week_of_month"] = *model.WeekOfMonth
	}
	if model.DayOfMonth != nil {
		modelMap["day_of_month"] = flex.IntValue(model.DayOfMonth)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesYearScheduleToMap(model *backuprecoveryv1.YearSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["day_of_year"] = *model.DayOfYear
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesFullBackupPolicyToMap(model *backuprecoveryv1.FullBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Schedule != nil {
		scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesFullScheduleToMap(model.Schedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesFullScheduleToMap(model *backuprecoveryv1.FullSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.DaySchedule != nil {
		dayScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesFullScheduleAndRetentionToMap(model *backuprecoveryv1.FullScheduleAndRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesFullScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model *backuprecoveryv1.Retention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesDataLockConfigToMap(model *backuprecoveryv1.DataLockConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = *model.Mode
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.EnableWormOnExternalTarget != nil {
		modelMap["enable_worm_on_external_target"] = *model.EnableWormOnExternalTarget
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesPrimaryBackupTargetToMap(model *backuprecoveryv1.PrimaryBackupTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetType != nil {
		modelMap["target_type"] = *model.TargetType
	}
	if model.ArchivalTargetSettings != nil {
		archivalTargetSettingsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesPrimaryArchivalTargetToMap(model.ArchivalTargetSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["archival_target_settings"] = []map[string]interface{}{archivalTargetSettingsMap}
	}
	if model.UseDefaultBackupTarget != nil {
		modelMap["use_default_backup_target"] = *model.UseDefaultBackupTarget
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesPrimaryArchivalTargetToMap(model *backuprecoveryv1.PrimaryArchivalTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_id"] = flex.IntValue(model.TargetID)
	if model.TargetName != nil {
		modelMap["target_name"] = *model.TargetName
	}
	if model.TierSettings != nil {
		tierSettingsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTierLevelSettingsToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesTierLevelSettingsToMap(model *backuprecoveryv1.TierLevelSettings) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	modelMap["cloud_platform"] = *model.CloudPlatform
	if model.GoogleTiering != nil {
		googleTieringMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesOracleTiersToMap(model.OracleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["oracle_tiering"] = []map[string]interface{}{oracleTieringMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesLogBackupPolicyToMap(model *backuprecoveryv1.LogBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesLogScheduleToMap(model *backuprecoveryv1.LogSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.MinuteSchedule != nil {
		minuteScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMinuteScheduleToMap(model.MinuteSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["minute_schedule"] = []map[string]interface{}{minuteScheduleMap}
	}
	if model.HourSchedule != nil {
		hourScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesHourScheduleToMap(model.HourSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["hour_schedule"] = []map[string]interface{}{hourScheduleMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesBmrBackupPolicyToMap(model *backuprecoveryv1.BmrBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesBmrScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesBmrScheduleToMap(model *backuprecoveryv1.BmrSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.DaySchedule != nil {
		dayScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesCdpBackupPolicyToMap(model *backuprecoveryv1.CdpBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCdpRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesCdpRetentionToMap(model *backuprecoveryv1.CdpRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesStorageArraySnapshotBackupPolicyToMap(model *backuprecoveryv1.StorageArraySnapshotBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesStorageArraySnapshotScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesStorageArraySnapshotScheduleToMap(model *backuprecoveryv1.StorageArraySnapshotSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.MinuteSchedule != nil {
		minuteScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMinuteScheduleToMap(model.MinuteSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["minute_schedule"] = []map[string]interface{}{minuteScheduleMap}
	}
	if model.HourSchedule != nil {
		hourScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesHourScheduleToMap(model.HourSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["hour_schedule"] = []map[string]interface{}{hourScheduleMap}
	}
	if model.DaySchedule != nil {
		dayScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(model *backuprecoveryv1.CancellationTimeoutParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TimeoutMins != nil {
		modelMap["timeout_mins"] = flex.IntValue(model.TimeoutMins)
	}
	if model.BackupType != nil {
		modelMap["backup_type"] = *model.BackupType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesBlackoutWindowToMap(model *backuprecoveryv1.BlackoutWindow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["day"] = *model.Day
	startTimeMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTimeOfDayToMap(model.StartTime)
	if err != nil {
		return modelMap, err
	}
	modelMap["start_time"] = []map[string]interface{}{startTimeMap}
	endTimeMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTimeOfDayToMap(model.EndTime)
	if err != nil {
		return modelMap, err
	}
	modelMap["end_time"] = []map[string]interface{}{endTimeMap}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesTimeOfDayToMap(model *backuprecoveryv1.TimeOfDay) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hour"] = flex.IntValue(model.Hour)
	modelMap["minute"] = flex.IntValue(model.Minute)
	if model.TimeZone != nil {
		modelMap["time_zone"] = *model.TimeZone
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesExtendedRetentionPolicyToMap(model *backuprecoveryv1.ExtendedRetentionPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesExtendedRetentionScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	if model.RunType != nil {
		modelMap["run_type"] = *model.RunType
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesExtendedRetentionScheduleToMap(model *backuprecoveryv1.ExtendedRetentionSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.Frequency != nil {
		modelMap["frequency"] = flex.IntValue(model.Frequency)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesTargetsConfigurationToMap(model *backuprecoveryv1.TargetsConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplicationTargets != nil {
		replicationTargets := []map[string]interface{}{}
		for _, replicationTargetsItem := range model.ReplicationTargets {
			replicationTargetsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesReplicationTargetConfigurationToMap(&replicationTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			replicationTargets = append(replicationTargets, replicationTargetsItemMap)
		}
		modelMap["replication_targets"] = replicationTargets
	}
	if model.ArchivalTargets != nil {
		archivalTargets := []map[string]interface{}{}
		for _, archivalTargetsItem := range model.ArchivalTargets {
			archivalTargetsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesArchivalTargetConfigurationToMap(&archivalTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			archivalTargets = append(archivalTargets, archivalTargetsItemMap)
		}
		modelMap["archival_targets"] = archivalTargets
	}
	if model.CloudSpinTargets != nil {
		cloudSpinTargets := []map[string]interface{}{}
		for _, cloudSpinTargetsItem := range model.CloudSpinTargets {
			cloudSpinTargetsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCloudSpinTargetConfigurationToMap(&cloudSpinTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			cloudSpinTargets = append(cloudSpinTargets, cloudSpinTargetsItemMap)
		}
		modelMap["cloud_spin_targets"] = cloudSpinTargets
	}
	if model.OnpremDeployTargets != nil {
		onpremDeployTargets := []map[string]interface{}{}
		for _, onpremDeployTargetsItem := range model.OnpremDeployTargets {
			onpremDeployTargetsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesOnpremDeployTargetConfigurationToMap(&onpremDeployTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			onpremDeployTargets = append(onpremDeployTargets, onpremDeployTargetsItemMap)
		}
		modelMap["onprem_deploy_targets"] = onpremDeployTargets
	}
	if model.RpaasTargets != nil {
		rpaasTargets := []map[string]interface{}{}
		for _, rpaasTargetsItem := range model.RpaasTargets {
			rpaasTargetsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRpaasTargetConfigurationToMap(&rpaasTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			rpaasTargets = append(rpaasTargets, rpaasTargetsItemMap)
		}
		modelMap["rpaas_targets"] = rpaasTargets
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesReplicationTargetConfigurationToMap(model *backuprecoveryv1.ReplicationTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	if model.CopyOnRunSuccess != nil {
		modelMap["copy_on_run_success"] = *model.CopyOnRunSuccess
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	if model.BackupRunType != nil {
		modelMap["backup_run_type"] = *model.BackupRunType
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	if model.AwsTargetConfig != nil {
		awsTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAWSTargetConfigToMap(model.AwsTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_target_config"] = []map[string]interface{}{awsTargetConfigMap}
	}
	if model.AzureTargetConfig != nil {
		azureTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAzureTargetConfigToMap(model.AzureTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_target_config"] = []map[string]interface{}{azureTargetConfigMap}
	}
	modelMap["target_type"] = *model.TargetType
	if model.RemoteTargetConfig != nil {
		remoteTargetConfigMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRemoteTargetConfigToMap(model.RemoteTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote_target_config"] = []map[string]interface{}{remoteTargetConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesTargetScheduleToMap(model *backuprecoveryv1.TargetSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.Frequency != nil {
		modelMap["frequency"] = flex.IntValue(model.Frequency)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesLogRetentionToMap(model *backuprecoveryv1.LogRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesAWSTargetConfigToMap(model *backuprecoveryv1.AWSTargetConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesAzureTargetConfigToMap(model *backuprecoveryv1.AzureTargetConfig) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesRemoteTargetConfigToMap(model *backuprecoveryv1.RemoteTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	if model.ClusterName != nil {
		modelMap["cluster_name"] = *model.ClusterName
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesArchivalTargetConfigurationToMap(model *backuprecoveryv1.ArchivalTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	if model.CopyOnRunSuccess != nil {
		modelMap["copy_on_run_success"] = *model.CopyOnRunSuccess
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	if model.BackupRunType != nil {
		modelMap["backup_run_type"] = *model.BackupRunType
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	modelMap["target_id"] = flex.IntValue(model.TargetID)
	if model.TargetName != nil {
		modelMap["target_name"] = *model.TargetName
	}
	if model.TargetType != nil {
		modelMap["target_type"] = *model.TargetType
	}
	if model.TierSettings != nil {
		tierSettingsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTierLevelSettingsToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	if model.ExtendedRetention != nil {
		extendedRetention := []map[string]interface{}{}
		for _, extendedRetentionItem := range model.ExtendedRetention {
			extendedRetentionItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesExtendedRetentionPolicyToMap(&extendedRetentionItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			extendedRetention = append(extendedRetention, extendedRetentionItemMap)
		}
		modelMap["extended_retention"] = extendedRetention
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesCloudSpinTargetConfigurationToMap(model *backuprecoveryv1.CloudSpinTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	if model.CopyOnRunSuccess != nil {
		modelMap["copy_on_run_success"] = *model.CopyOnRunSuccess
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	if model.BackupRunType != nil {
		modelMap["backup_run_type"] = *model.BackupRunType
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	targetMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCloudSpinTargetToMap(model.Target)
	if err != nil {
		return modelMap, err
	}
	modelMap["target"] = []map[string]interface{}{targetMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesCloudSpinTargetToMap(model *backuprecoveryv1.CloudSpinTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsParams != nil {
		awsParamsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAwsCloudSpinParamsToMap(model.AwsParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_params"] = []map[string]interface{}{awsParamsMap}
	}
	if model.AzureParams != nil {
		azureParamsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesAzureCloudSpinParamsToMap(model.AzureParams)
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
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesAwsCloudSpinParamsToMap(model *backuprecoveryv1.AwsCloudSpinParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomTagList != nil {
		customTagList := []map[string]interface{}{}
		for _, customTagListItem := range model.CustomTagList {
			customTagListItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCustomTagParamsToMap(&customTagListItem) // #nosec G601
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

func DataSourceIbmBackupRecoveryProtectionPoliciesCustomTagParamsToMap(model *backuprecoveryv1.CustomTagParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["key"] = *model.Key
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesAzureCloudSpinParamsToMap(model *backuprecoveryv1.AzureCloudSpinParams) (map[string]interface{}, error) {
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

func DataSourceIbmBackupRecoveryProtectionPoliciesOnpremDeployTargetConfigurationToMap(model *backuprecoveryv1.OnpremDeployTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	if model.CopyOnRunSuccess != nil {
		modelMap["copy_on_run_success"] = *model.CopyOnRunSuccess
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	if model.BackupRunType != nil {
		modelMap["backup_run_type"] = *model.BackupRunType
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	if model.Params != nil {
		paramsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesOnpremDeployParamsToMap(model.Params)
		if err != nil {
			return modelMap, err
		}
		modelMap["params"] = []map[string]interface{}{paramsMap}
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesOnpremDeployParamsToMap(model *backuprecoveryv1.OnpremDeployParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesRpaasTargetConfigurationToMap(model *backuprecoveryv1.RpaasTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	if model.CopyOnRunSuccess != nil {
		modelMap["copy_on_run_success"] = *model.CopyOnRunSuccess
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	if model.BackupRunType != nil {
		modelMap["backup_run_type"] = *model.BackupRunType
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	modelMap["target_id"] = flex.IntValue(model.TargetID)
	if model.TargetName != nil {
		modelMap["target_name"] = *model.TargetName
	}
	if model.TargetType != nil {
		modelMap["target_type"] = *model.TargetType
	}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesCascadedTargetConfigurationToMap(model *backuprecoveryv1.CascadedTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["source_cluster_id"] = flex.IntValue(model.SourceClusterID)
	remoteTargetsMap, err := DataSourceIbmBackupRecoveryProtectionPoliciesTargetsConfigurationToMap(model.RemoteTargets)
	if err != nil {
		return modelMap, err
	}
	modelMap["remote_targets"] = []map[string]interface{}{remoteTargetsMap}
	return modelMap, nil
}

func DataSourceIbmBackupRecoveryProtectionPoliciesRetryOptionsToMap(model *backuprecoveryv1.RetryOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Retries != nil {
		modelMap["retries"] = flex.IntValue(model.Retries)
	}
	if model.RetryIntervalMins != nil {
		modelMap["retry_interval_mins"] = flex.IntValue(model.RetryIntervalMins)
	}
	return modelMap, nil
}
