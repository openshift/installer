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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/ibm-backup-recovery-sdk-go/backuprecoveryv1"
)

func ResourceIbmBackupRecoveryProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmBackupRecoveryProtectionPolicyCreate,
		ReadContext:   resourceIbmBackupRecoveryProtectionPolicyRead,
		UpdateContext: resourceIbmBackupRecoveryProtectionPolicyUpdate,
		DeleteContext: resourceIbmBackupRecoveryProtectionPolicyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"x_ibm_tenant_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the key to be used to encrypt the source credential. If includeSourceCredentials is set to true this key must be specified.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the Protection Policy.",
			},
			"backup_policy": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Specifies the backup schedule and retentions of a Protection Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"regular": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Specifies the Incremental and Full policy settings and also the common Retention policy settings.\".",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"incremental": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies incremental backup settings for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies settings that defines how frequent backup will be performed for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies how often to start new runs of a Protection Group. <br>'Minutes' specifies that Protection Group run starts periodically after certain number of minutes specified in 'frequency' field. <br>'Hours' specifies that Protection Group run starts periodically after certain number of hours specified in 'frequency' field. <br>'Days' specifies that Protection Group run starts periodically after certain number of days specified in 'frequency' field. <br>'Week' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Month' specifies that new Protection Group runs start monthly on certain day of specific week. This schedule needs 'weekOfMonth' and 'dayOfWeek' fields to be set. <br> Example: To run the Protection Group on Second Sunday of Every Month, following schedule need to be set: <br> unit: 'Month' <br> dayOfWeek: 'Sunday' <br> weekOfMonth: 'Second'.",
															},
															"minute_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of minutes.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"hour_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of hours.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"day_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"week_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"month_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"week_of_month": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																		},
																		"day_of_month": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																		},
																	},
																},
															},
															"year_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_year": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies full backup settings for a Protection Group. Currently, full backup settings can be specified by using either of 'schedule' or 'schdulesAndRetentions' field. Using 'schdulesAndRetentions' is recommended when multiple full backups need to be configured. If full and incremental backup has common retention then only setting 'schedule' is recommended.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that defines how frequent full backup will be performed for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies how often to start new runs of a Protection Group. <br>'Days' specifies that Protection Group run starts periodically on every day. For full backup schedule, currently we only support frequecny of 1 which indicates that full backup will be performed daily. <br>'Weeks' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Months' specifies that new Protection Group runs start monthly on certain day of specific week. This schedule needs 'weekOfMonth' and 'dayOfWeek' fields to be set. <br>'ProtectOnce' specifies that groups using this policy option will run only once and after that group will permanently be disabled. <br> Example: To run the Protection Group on Second Sunday of Every Month, following schedule need to be set: <br> unit: 'Month' <br> dayOfWeek: 'Sunday' <br> weekOfMonth: 'Second'.",
															},
															"day_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"week_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"month_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"week_of_month": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																		},
																		"day_of_month": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																		},
																	},
																},
															},
															"year_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_year": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
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
										Optional:    true,
										Description: "Specifies multiple schedules and retentions for full backup. Specify either of the 'full' or 'fullBackups' values. Its recommended to use 'fullBaackups' value since 'full' will be deprecated after few releases.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies settings that defines how frequent full backup will be performed for a Protection Group.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies how often to start new runs of a Protection Group. <br>'Days' specifies that Protection Group run starts periodically on every day. For full backup schedule, currently we only support frequecny of 1 which indicates that full backup will be performed daily. <br>'Weeks' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Months' specifies that new Protection Group runs start monthly on certain day of specific week. This schedule needs 'weekOfMonth' and 'dayOfWeek' fields to be set. <br>'ProtectOnce' specifies that groups using this policy option will run only once and after that group will permanently be disabled. <br> Example: To run the Protection Group on Second Sunday of Every Month, following schedule need to be set: <br> unit: 'Month' <br> dayOfWeek: 'Sunday' <br> weekOfMonth: 'Second'.",
															},
															"day_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
																		},
																	},
																},
															},
															"week_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Required:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																	},
																},
															},
															"month_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_week": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"week_of_month": &schema.Schema{
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
																		},
																		"day_of_month": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
																		},
																	},
																},
															},
															"year_schedule": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"day_of_year": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
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
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the primary backup target settings for regular backups. If the backup target field is not specified then backup will be taken locally on the Cohesity cluster.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the primary backup location where backups will be stored. If not specified, then default is assumed as local backup on Cohesity cluster.",
												},
												"archival_target_settings": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the primary archival settings. Mainly used for cloud direct archive (CAD) policy where primary backup is stored on archival target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"target_id": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the Archival target id to take primary backup.",
															},
															"target_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies the Archival target name where Snapshots are copied.",
															},
															"tier_settings": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies the settings tier levels configured with each archival target. The tier settings need to be applied in specific order and default tier should always be passed as first entry in tiers array. The following example illustrates how to configure tiering input for AWS tiering. Same type of input structure applied to other cloud platforms also. <br>If user wants to achieve following tiering for backup, <br>User Desired Tiering- <br><t>1.Archive Full back up for 12 Months <br><t>2.Tier Levels <br><t><t>[1,12] [ <br><t><t><t>s3 (1 to 2 months), (default tier) <br><t><t><t>s3 Intelligent tiering (3 to 6 months), <br><t><t><t>s3 One Zone (7 to 9 months) <br><t><t><t>Glacier (10 to 12 months)] <br><t>API Input <br><t><t>1.tiers-[ <br><t><t><t>{'tierType': 'S3','moveAfterUnit':'months', <br><t><t><t>'moveAfter':2 - move from s3 to s3Inte after 2 months}, <br><t><t><t>{'tierType': 'S3Inte','moveAfterUnit':'months', <br><t><t><t>'moveAfter':4 - move from S3Inte to Glacier after 4 months}, <br><t><t><t>{'tierType': 'Glacier', 'moveAfterUnit':'months', <br><t><t><t>'moveAfter': 3 - move from Glacier to S3 One Zone after 3 months }, <br><t><t><t>{'tierType': 'S3 One Zone', 'moveAfterUnit': nil, <br><t><t><t>'moveAfter': nil - For the last record, 'moveAfter' and 'moveAfterUnit' <br><t><t><t>will be ignored since there are no further tier for data movement } <br><t><t><t>}].",
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
																	},
																},
															},
														},
													},
												},
												"use_default_backup_target": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies log backup settings for a Protection Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies settings that defines how frequent log backup will be performed for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies how often to start new Protection Group Runs of a Protection Group. <br>'Minutes' specifies that Protection Group run starts periodically after certain number of minutes specified in 'frequency' field. <br>'Hours' specifies that Protection Group run starts periodically after certain number of hours specified in 'frequency' field.",
												},
												"minute_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of minutes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
															},
														},
													},
												},
												"hour_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of hours.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
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
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies the BMR schedule in case of physical source protection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies settings that defines how frequent bmr backup will be performed for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies how often to start new runs of a Protection Group. <br>'Weeks' specifies that new Protection Group runs start weekly on certain days specified using 'dayOfWeek' field. <br>'Months' specifies that new Protection Group runs start monthly on certain day of specific week.",
												},
												"day_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
															},
														},
													},
												},
												"week_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_of_week": &schema.Schema{
																Type:        schema.TypeList,
																Required:    true,
																Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"month_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_of_week": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"week_of_month": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
															},
															"day_of_month": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
															},
														},
													},
												},
												"year_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_of_year": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
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
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies CDP (Continious Data Protection) backup settings for a Protection Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a CDP backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a CDP backup measured in minutes or hours.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a cdp backup retention.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "Specifies storage snapshot managment backup settings for a Protection Group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies settings that defines how frequent Storage Snapshot Management backup will be performed for a Protection Group.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies how often to start new Protection Group Runs of a Protection Group. <br>'Minutes' specifies that Protection Group run starts periodically after certain number of minutes specified in 'frequency' field. <br>'Hours' specifies that Protection Group run starts periodically after certain number of hours specified in 'frequency' field.",
												},
												"minute_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of minutes.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
															},
														},
													},
												},
												"hour_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of hours.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
															},
														},
													},
												},
												"day_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start after certain number of days.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the backup schedule. <br> Example: If 'frequency' set to 2 and the unit is 'Hours', then Snapshots are backed up every 2 hours. <br> This field is only applicable if unit is 'Minutes', 'Hours' or 'Days'.",
															},
														},
													},
												},
												"week_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to start on certain days of week.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_of_week": &schema.Schema{
																Type:        schema.TypeList,
																Required:    true,
																Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
														},
													},
												},
												"month_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group runs to on specific week and specific days of that week.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_of_week": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "Specifies a list of days of the week when to start Protection Group Runs. <br> Example: To run a Protection Group on every Monday and Tuesday, set the schedule with following values: <br>  unit: 'Weeks' <br>  dayOfWeek: ['Monday','Tuesday'].",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"week_of_month": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Specifies the week of the month (such as 'Third') or nth day of month (such as 'First' or 'Last') in a Monthly Schedule specified by unit field as 'Months'. <br>This field can be used in combination with 'dayOfWeek' to define the day in the month to start the Protection Group Run. <br> Example: if 'weekOfMonth' is set to 'Third' and day is set to 'Monday', a backup is performed on the third Monday of every month. <br> Example: if 'weekOfMonth' is set to 'Last' and dayOfWeek is not set, a backup is performed on the last day of every month.",
															},
															"day_of_month": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the exact date of the month (such as 18) in a Monthly Schedule specified by unit field as 'Years'. <br> Example: if 'dayOfMonth' is set to '18', a backup is performed on the 18th of every month.",
															},
														},
													},
												},
												"year_schedule": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies settings that define a schedule for a Protection Group to run on specific year and specific day of that year.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"day_of_year": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
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
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
							Optional:    true,
							Description: "Specifies the backup timeouts for different type of runs(kFull, kRegular etc.).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"timeout_mins": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies the timeout in mins.",
									},
									"backup_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
				Optional:    true,
				Description: "Specifies the description of the Protection Policy.",
			},
			"blackout_window": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Blackout Windows. If specified, this field defines blackout periods when new Group Runs are not started. If a Group Run has been scheduled but not yet executed and the blackout period starts, the behavior depends on the policy field AbortInBlackoutPeriod.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies a day in the week when no new Protection Group Runs should be started such as 'Sunday'. Specifies a day in a week such as 'Sunday', 'Monday', etc.",
						},
						"start_time": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Specifies the time of day. Used for scheduling purposes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the hour of the day (0-23).",
									},
									"minute": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the minute of the hour (0-59).",
									},
									"time_zone": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "America/Los_Angeles",
										Description: "Specifies the time zone of the user. If not specified, default value is assumed as America/Los_Angeles.",
									},
								},
							},
						},
						"end_time": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Specifies the time of day. Used for scheduling purposes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"hour": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the hour of the day (0-23).",
									},
									"minute": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the minute of the hour (0-59).",
									},
									"time_zone": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "America/Los_Angeles",
										Description: "Specifies the time zone of the user. If not specified, default value is assumed as America/Los_Angeles.",
									},
								},
							},
						},
						"config_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the unique identifier for the target getting added. This field need to be passed olny when policies are updated.",
						},
					},
				},
			},
			"extended_retention": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies additional retention policies that should be applied to the backup snapshots. A backup snapshot will be retained up to a time that is the maximum of all retention policies that are applicable to it.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"schedule": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Specifies a schedule frequency and schedule unit for Extended Retentions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specifies the unit interval for retention of Snapshots. <br>'Runs' means that the Snapshot copy retained after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy retained hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy gets retained daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy is retained weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy is retained monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy is retained yearly at the frequency set in the Frequency.",
									},
									"frequency": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Specifies a factor to multiply the unit by, to determine the retention schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is retained.",
									},
								},
							},
						},
						"retention": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Specifies the retention of a backup.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
									},
									"duration": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
									},
									"data_lock_config": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"mode": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
												},
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
												},
												"enable_worm_on_external_target": &schema.Schema{
													Type:        schema.TypeBool,
													Optional:    true,
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
							Optional:    true,
							Description: "The backup run type to which this extended retention applies to. If this is not set, the extended retention will be applicable to all non-log backup types. Currently, the only value that can be set here is Full.'Regular' indicates a incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a Regular schedule captures all the blocks.'Full' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized.'Log' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time.'System' indicates a system backup. System backups are used to do bare metal recovery of the system to a specific point in time.",
						},
						"config_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the unique identifier for the target getting added. This field need to be passed olny when policies are updated.",
						},
					},
				},
			},
			"remote_target_policy": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Specifies the replication, archival and cloud spin targets of Protection Policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"replication_targets": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
												},
												"frequency": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
												},
											},
										},
									},
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
									},
									"backup_run_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
									},
									"run_timeouts": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout_mins": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the timeout in mins.",
												},
												"backup_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The scheduled backup type(kFull, kRegular etc.).",
												},
											},
										},
									},
									"log_retention": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
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
													Required:    true,
													Description: "Specifies id of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
												},
												"region_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the source id of the AWS protection source registered on IBM cluster.",
												},
											},
										},
									},
									"azure_target_config": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
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
													Optional:    true,
													Description: "Specifies id of the Azure resource group used to filter regions in UI.",
												},
												"resource_group_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies name of the Azure resource group used to filter regions in UI.",
												},
												"source_id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
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
										Required:    true,
										Description: "Specifies the type of target to which replication need to be performed.",
									},
									"remote_target_config": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the configuration for adding remote cluster as repilcation target.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"cluster_id": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
												},
												"frequency": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
												},
											},
										},
									},
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
									},
									"backup_run_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
									},
									"run_timeouts": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout_mins": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the timeout in mins.",
												},
												"backup_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The scheduled backup type(kFull, kRegular etc.).",
												},
											},
										},
									},
									"log_retention": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Required:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the settings tier levels configured with each archival target. The tier settings need to be applied in specific order and default tier should always be passed as first entry in tiers array. The following example illustrates how to configure tiering input for AWS tiering. Same type of input structure applied to other cloud platforms also. <br>If user wants to achieve following tiering for backup, <br>User Desired Tiering- <br><t>1.Archive Full back up for 12 Months <br><t>2.Tier Levels <br><t><t>[1,12] [ <br><t><t><t>s3 (1 to 2 months), (default tier) <br><t><t><t>s3 Intelligent tiering (3 to 6 months), <br><t><t><t>s3 One Zone (7 to 9 months) <br><t><t><t>Glacier (10 to 12 months)] <br><t>API Input <br><t><t>1.tiers-[ <br><t><t><t>{'tierType': 'S3','moveAfterUnit':'months', <br><t><t><t>'moveAfter':2 - move from s3 to s3Inte after 2 months}, <br><t><t><t>{'tierType': 'S3Inte','moveAfterUnit':'months', <br><t><t><t>'moveAfter':4 - move from S3Inte to Glacier after 4 months}, <br><t><t><t>{'tierType': 'Glacier', 'moveAfterUnit':'months', <br><t><t><t>'moveAfter': 3 - move from Glacier to S3 One Zone after 3 months }, <br><t><t><t>{'tierType': 'S3 One Zone', 'moveAfterUnit': nil, <br><t><t><t>'moveAfter': nil - For the last record, 'moveAfter' and 'moveAfterUnit' <br><t><t><t>will be ignored since there are no further tier for data movement } <br><t><t><t>}].",
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
											},
										},
									},
									"extended_retention": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies additional retention policies that should be applied to the archived backup. Archived backup snapshot will be retained up to a time that is the maximum of all retention policies that are applicable to it.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies a schedule frequency and schedule unit for Extended Retentions.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the unit interval for retention of Snapshots. <br>'Runs' means that the Snapshot copy retained after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy retained hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy gets retained daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy is retained weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy is retained monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy is retained yearly at the frequency set in the Frequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the retention schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is retained.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Optional:    true,
													Description: "The backup run type to which this extended retention applies to. If this is not set, the extended retention will be applicable to all non-log backup types. Currently, the only value that can be set here is Full.'Regular' indicates a incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a Regular schedule captures all the blocks.'Full' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized.'Log' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time.'System' indicates a system backup. System backups are used to do bare metal recovery of the system to a specific point in time.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
												},
												"frequency": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
												},
											},
										},
									},
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
									},
									"backup_run_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
									},
									"run_timeouts": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout_mins": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the timeout in mins.",
												},
												"backup_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The scheduled backup type(kFull, kRegular etc.).",
												},
											},
										},
									},
									"log_retention": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the details about Cloud Spin target where backup snapshots may be converted and stored.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"aws_params": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies various resources when converting and deploying a VM to AWS.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"custom_tag_list": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
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
																Required:    true,
																Description: "Specifies id of the AWS region in which to deploy the VM.",
															},
															"subnet_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the subnet within above VPC.",
															},
															"vpc_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the Virtual Private Cloud to chose for the instance type.",
															},
														},
													},
												},
												"azure_params": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies various resources when converting and deploying a VM to Azure.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"availability_set_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the availability set.",
															},
															"network_resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the resource group for the selected virtual network.",
															},
															"resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the Azure resource group. Its value is globally unique within Azure.",
															},
															"storage_account_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
															},
															"storage_container_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the storage container within the above storage account.",
															},
															"storage_resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the resource group for the selected storage account.",
															},
															"temp_vm_resource_group_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the temporary Azure resource group.",
															},
															"temp_vm_storage_account_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the temporary VM storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
															},
															"temp_vm_storage_container_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies id of the temporary VM storage container within the above storage account.",
															},
															"temp_vm_subnet_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies Id of the temporary VM subnet within the above virtual network.",
															},
															"temp_vm_virtual_network_id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies Id of the temporary VM Virtual Network.",
															},
														},
													},
												},
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
												},
												"frequency": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
												},
											},
										},
									},
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
									},
									"backup_run_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
									},
									"run_timeouts": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout_mins": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the timeout in mins.",
												},
												"backup_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The scheduled backup type(kFull, kRegular etc.).",
												},
											},
										},
									},
									"log_retention": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the details about OnpremDeploy target where backup snapshots may be converted and deployed.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
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
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
												},
												"frequency": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
												},
											},
										},
									},
									"retention": &schema.Schema{
										Type:        schema.TypeList,
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Optional:    true,
										Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
									},
									"backup_run_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
									},
									"run_timeouts": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"timeout_mins": &schema.Schema{
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Specifies the timeout in mins.",
												},
												"backup_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The scheduled backup type(kFull, kRegular etc.).",
												},
											},
										},
									},
									"log_retention": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Specifies the retention of a backup.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"unit": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
												},
												"duration": &schema.Schema{
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
												},
												"data_lock_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"mode": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
															},
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
															},
															"enable_worm_on_external_target": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
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
										Required:    true,
										Description: "Specifies the RPaaS target to copy the Snapshots.",
									},
									"target_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Specifies the RPaaS target name where Snapshots are copied.",
									},
									"target_type": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
				Optional:    true,
				Description: "Specifies the configuration for cascaded replications. Using cascaded replication, replication cluster(Rx) can further replicate and archive the snapshot copies to further targets. Its recommended to create cascaded configuration where protection group will be created.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_cluster_id": &schema.Schema{
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Specifies the source cluster id from where the remote operations will be performed to the next set of remote targets.",
						},
						"remote_targets": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Specifies the replication, archival and cloud spin targets of Protection Policy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"replication_targets": &schema.Schema{
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Optional:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"timeout_mins": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the timeout in mins.",
															},
															"backup_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The scheduled backup type(kFull, kRegular etc.).",
															},
														},
													},
												},
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													MaxItems:    1,
													Optional:    true,
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
																Required:    true,
																Description: "Specifies id of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
															},
															"region_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies name of the AWS region in which to replicate the Snapshot to. Applicable if replication target is AWS target.",
															},
															"source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the source id of the AWS protection source registered on IBM cluster.",
															},
														},
													},
												},
												"azure_target_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
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
																Optional:    true,
																Description: "Specifies id of the Azure resource group used to filter regions in UI.",
															},
															"resource_group_name": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Specifies name of the Azure resource group used to filter regions in UI.",
															},
															"source_id": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
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
													Required:    true,
													Description: "Specifies the type of target to which replication need to be performed.",
												},
												"remote_target_config": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the configuration for adding remote cluster as repilcation target.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"cluster_id": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
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
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Optional:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"timeout_mins": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the timeout in mins.",
															},
															"backup_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The scheduled backup type(kFull, kRegular etc.).",
															},
														},
													},
												},
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Required:    true,
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
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the settings tier levels configured with each archival target. The tier settings need to be applied in specific order and default tier should always be passed as first entry in tiers array. The following example illustrates how to configure tiering input for AWS tiering. Same type of input structure applied to other cloud platforms also. <br>If user wants to achieve following tiering for backup, <br>User Desired Tiering- <br><t>1.Archive Full back up for 12 Months <br><t>2.Tier Levels <br><t><t>[1,12] [ <br><t><t><t>s3 (1 to 2 months), (default tier) <br><t><t><t>s3 Intelligent tiering (3 to 6 months), <br><t><t><t>s3 One Zone (7 to 9 months) <br><t><t><t>Glacier (10 to 12 months)] <br><t>API Input <br><t><t>1.tiers-[ <br><t><t><t>{'tierType': 'S3','moveAfterUnit':'months', <br><t><t><t>'moveAfter':2 - move from s3 to s3Inte after 2 months}, <br><t><t><t>{'tierType': 'S3Inte','moveAfterUnit':'months', <br><t><t><t>'moveAfter':4 - move from S3Inte to Glacier after 4 months}, <br><t><t><t>{'tierType': 'Glacier', 'moveAfterUnit':'months', <br><t><t><t>'moveAfter': 3 - move from Glacier to S3 One Zone after 3 months }, <br><t><t><t>{'tierType': 'S3 One Zone', 'moveAfterUnit': nil, <br><t><t><t>'moveAfter': nil - For the last record, 'moveAfter' and 'moveAfterUnit' <br><t><t><t>will be ignored since there are no further tier for data movement } <br><t><t><t>}].",
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
														},
													},
												},
												"extended_retention": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies additional retention policies that should be applied to the archived backup. Archived backup snapshot will be retained up to a time that is the maximum of all retention policies that are applicable to it.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"schedule": &schema.Schema{
																Type:        schema.TypeList,
																MinItems:    1,
																MaxItems:    1,
																Required:    true,
																Description: "Specifies a schedule frequency and schedule unit for Extended Retentions.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the unit interval for retention of Snapshots. <br>'Runs' means that the Snapshot copy retained after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy retained hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy gets retained daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy is retained weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy is retained monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy is retained yearly at the frequency set in the Frequency.",
																		},
																		"frequency": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies a factor to multiply the unit by, to determine the retention schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is retained.",
																		},
																	},
																},
															},
															"retention": &schema.Schema{
																Type:        schema.TypeList,
																MinItems:    1,
																MaxItems:    1,
																Required:    true,
																Description: "Specifies the retention of a backup.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
																		},
																		"data_lock_config": &schema.Schema{
																			Type:        schema.TypeList,
																			MaxItems:    1,
																			Optional:    true,
																			Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																			Elem: &schema.Resource{
																				Schema: map[string]*schema.Schema{
																					"mode": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																					},
																					"unit": &schema.Schema{
																						Type:        schema.TypeString,
																						Required:    true,
																						Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																					},
																					"duration": &schema.Schema{
																						Type:        schema.TypeInt,
																						Required:    true,
																						Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																					},
																					"enable_worm_on_external_target": &schema.Schema{
																						Type:        schema.TypeBool,
																						Optional:    true,
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
																Optional:    true,
																Description: "The backup run type to which this extended retention applies to. If this is not set, the extended retention will be applicable to all non-log backup types. Currently, the only value that can be set here is Full.'Regular' indicates a incremental (CBT) backup. Incremental backups utilizing CBT (if supported) are captured of the target protection objects. The first run of a Regular schedule captures all the blocks.'Full' indicates a full (no CBT) backup. A complete backup (all blocks) of the target protection objects are always captured and Change Block Tracking (CBT) is not utilized.'Log' indicates a Database Log backup. Capture the database transaction logs to allow rolling back to a specific point in time.'System' indicates a system backup. System backups are used to do bare metal recovery of the system to a specific point in time.",
															},
															"config_id": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
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
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Optional:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"timeout_mins": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the timeout in mins.",
															},
															"backup_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The scheduled backup type(kFull, kRegular etc.).",
															},
														},
													},
												},
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the details about Cloud Spin target where backup snapshots may be converted and stored.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"aws_params": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies various resources when converting and deploying a VM to AWS.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"custom_tag_list": &schema.Schema{
																			Type:        schema.TypeList,
																			Optional:    true,
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
																			Required:    true,
																			Description: "Specifies id of the AWS region in which to deploy the VM.",
																		},
																		"subnet_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the subnet within above VPC.",
																		},
																		"vpc_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the Virtual Private Cloud to chose for the instance type.",
																		},
																	},
																},
															},
															"azure_params": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies various resources when converting and deploying a VM to Azure.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"availability_set_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies the availability set.",
																		},
																		"network_resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the resource group for the selected virtual network.",
																		},
																		"resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the Azure resource group. Its value is globally unique within Azure.",
																		},
																		"storage_account_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
																		},
																		"storage_container_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the storage container within the above storage account.",
																		},
																		"storage_resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the resource group for the selected storage account.",
																		},
																		"temp_vm_resource_group_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the temporary Azure resource group.",
																		},
																		"temp_vm_storage_account_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the temporary VM storage account that will contain the storage container within which we will create the blob that will become the VHD disk for the cloned VM.",
																		},
																		"temp_vm_storage_container_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies id of the temporary VM storage container within the above storage account.",
																		},
																		"temp_vm_subnet_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies Id of the temporary VM subnet within the above virtual network.",
																		},
																		"temp_vm_virtual_network_id": &schema.Schema{
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Specifies Id of the temporary VM Virtual Network.",
																		},
																	},
																},
															},
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
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
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Optional:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"timeout_mins": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the timeout in mins.",
															},
															"backup_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The scheduled backup type(kFull, kRegular etc.).",
															},
														},
													},
												},
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the details about OnpremDeploy target where backup snapshots may be converted and deployed.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"id": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
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
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"schedule": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies a schedule fregquency and schedule unit for copying Snapshots to backup targets.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specifies the frequency that Snapshots should be copied to the specified target. Used in combination with multiplier. <br>'Runs' means that the Snapshot copy occurs after the number of Protection Group Runs equals the number specified in the frequency. <br>'Hours' means that the Snapshot copy occurs hourly at the frequency set in the frequency, for example if scheduleFrequency is 2, the copy occurs every 2 hours. <br>'Days' means that the Snapshot copy occurs daily at the frequency set in the frequency. <br>'Weeks' means that the Snapshot copy occurs weekly at the frequency set in the frequency. <br>'Months' means that the Snapshot copy occurs monthly at the frequency set in the Frequency. <br>'Years' means that the Snapshot copy occurs yearly at the frequency set in the scheduleFrequency.",
															},
															"frequency": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies a factor to multiply the unit by, to determine the copy schedule. For example if set to 2 and the unit is hourly, then Snapshots from the first eligible Job Run for every 2 hour period is copied.",
															},
														},
													},
												},
												"retention": &schema.Schema{
													Type:        schema.TypeList,
													MinItems:    1,
													MaxItems:    1,
													Required:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Optional:    true,
													Description: "Specifies if Snapshots are copied from the first completely successful Protection Group Run or the first partially successful Protection Group Run occurring at the start of the replication schedule. <br> If true, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule that was completely successful i.e. Snapshots for all the Objects in the Protection Group were successfully captured. <br> If false, Snapshots are copied from the first Protection Group Run occurring at the start of the replication schedule, even if first Protection Group Run was not completely successful i.e. Snapshots were not captured for all Objects in the Protection Group.",
												},
												"config_id": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies the unique identifier for the target getting added. This field need to be passed only when policies are being updated.",
												},
												"backup_run_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Specifies which type of run should be copied, if not set, all types of runs will be eligible for copying. If set, this will ensure that the first run of given type in the scheduled period will get copied. Currently, this can only be set to Full.",
												},
												"run_timeouts": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Specifies the replication/archival timeouts for different type of runs(kFull, kRegular etc.).",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"timeout_mins": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Specifies the timeout in mins.",
															},
															"backup_type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "The scheduled backup type(kFull, kRegular etc.).",
															},
														},
													},
												},
												"log_retention": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Specifies the retention of a backup.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"unit": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Specificies the Retention Unit of a backup measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Years' then number of retention days will be 365 * 2 = 730 days.",
															},
															"duration": &schema.Schema{
																Type:        schema.TypeInt,
																Required:    true,
																Description: "Specifies the duration for a backup retention. <br> Example. If duration is 7 and unit is Months, the retention of a backup is 7 * 30 = 210 days.",
															},
															"data_lock_config": &schema.Schema{
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Specifies WORM retention type for the snapshots. When a WORM retention type is specified, the snapshots of the Protection Groups using this policy will be kept for the last N days as specified in the duration of the datalock. During that time, the snapshots cannot be deleted.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"mode": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specifies the type of WORM retention type. 'Compliance' implies WORM retention is set for compliance reason. 'Administrative' implies WORM retention is set for administrative purposes.",
																		},
																		"unit": &schema.Schema{
																			Type:        schema.TypeString,
																			Required:    true,
																			Description: "Specificies the Retention Unit of a dataLock measured in days, months or years. <br> If unit is 'Months', then number specified in duration is multiplied to 30. <br> Example: If duration is 4 and unit is 'Months' then number of retention days will be 30 * 4 = 120 days. <br> If unit is 'Years', then number specified in duration is multiplied to 365. <br> If duration is 2 and unit is 'Months' then number of retention days will be 365 * 2 = 730 days.",
																		},
																		"duration": &schema.Schema{
																			Type:        schema.TypeInt,
																			Required:    true,
																			Description: "Specifies the duration for a dataLock. <br> Example. If duration is 7 and unit is Months, the dataLock is enabled for last 7 * 30 = 210 days of the backup.",
																		},
																		"enable_worm_on_external_target": &schema.Schema{
																			Type:        schema.TypeBool,
																			Optional:    true,
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
													Required:    true,
													Description: "Specifies the RPaaS target to copy the Snapshots.",
												},
												"target_name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Specifies the RPaaS target name where Snapshots are copied.",
												},
												"target_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
				MaxItems:    1,
				Optional:    true,
				Description: "Retry Options of a Protection Policy when a Protection Group run fails.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"retries": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the number of times to retry capturing Snapshots before the Protection Group Run fails.",
						},
						"retry_interval_mins": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Specifies the number of minutes before retrying a failed Protection Group.",
						},
					},
				},
			},
			"data_lock": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				// ValidateFunc: validate.InvokeValidator("ibm_backup_recovery_protection_policy", "data_lock"),
				Description: "This field is now deprecated. Please use the DataLockConfig in the backup retention.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the current policy verison. Policy version is incremented for optionally supporting new features and differentialting across releases.",
			},
			"is_cbs_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies true if Calender Based Schedule is supported by client. Default value is assumed as false for this feature.",
			},
			"last_modification_time_usecs": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Specifies the last time this Policy was updated. If this is passed into a PUT request, then the backend will validate that the timestamp passed in matches the time that the policy was actually last modified. If the two timestamps do not match, then the request will be rejected with a stale error.",
			},
			"template_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
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
			"policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "policy ID",
			},
		},
	}
}

func ResourceIbmBackupRecoveryProtectionPolicyValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "data_lock",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "Administrative, Compliance",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_backup_recovery_protection_policy", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmBackupRecoveryProtectionPolicyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tenantId := d.Get("x_ibm_tenant_id").(string)

	createProtectionPolicyOptions := &backuprecoveryv1.CreateProtectionPolicyOptions{}

	createProtectionPolicyOptions.SetXIBMTenantID(tenantId)
	createProtectionPolicyOptions.SetName(d.Get("name").(string))
	backupPolicyModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToBackupPolicy(d.Get("backup_policy.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "parse-backup_policy").GetDiag()
	}
	createProtectionPolicyOptions.SetBackupPolicy(backupPolicyModel)
	if _, ok := d.GetOk("description"); ok {
		createProtectionPolicyOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("blackout_window"); ok {
		var blackoutWindow []backuprecoveryv1.BlackoutWindow
		for _, v := range d.Get("blackout_window").([]interface{}) {
			value := v.(map[string]interface{})
			blackoutWindowItem, err := ResourceIbmBackupRecoveryProtectionPolicyMapToBlackoutWindow(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "parse-blackout_window").GetDiag()
			}
			blackoutWindow = append(blackoutWindow, *blackoutWindowItem)
		}
		createProtectionPolicyOptions.SetBlackoutWindow(blackoutWindow)
	}
	if _, ok := d.GetOk("extended_retention"); ok {
		var extendedRetention []backuprecoveryv1.ExtendedRetentionPolicy
		for _, v := range d.Get("extended_retention").([]interface{}) {
			value := v.(map[string]interface{})
			extendedRetentionItem, err := ResourceIbmBackupRecoveryProtectionPolicyMapToExtendedRetentionPolicy(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "parse-extended_retention").GetDiag()
			}
			extendedRetention = append(extendedRetention, *extendedRetentionItem)
		}
		createProtectionPolicyOptions.SetExtendedRetention(extendedRetention)
	}
	if _, ok := d.GetOk("remote_target_policy"); ok {
		remoteTargetPolicyModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetsConfiguration(d.Get("remote_target_policy.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "parse-remote_target_policy").GetDiag()
		}
		createProtectionPolicyOptions.SetRemoteTargetPolicy(remoteTargetPolicyModel)
	}
	if _, ok := d.GetOk("cascaded_targets_config"); ok {
		var cascadedTargetsConfig []backuprecoveryv1.CascadedTargetConfiguration
		for _, v := range d.Get("cascaded_targets_config").([]interface{}) {
			value := v.(map[string]interface{})
			cascadedTargetsConfigItem, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCascadedTargetConfiguration(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "parse-cascaded_targets_config").GetDiag()
			}
			cascadedTargetsConfig = append(cascadedTargetsConfig, *cascadedTargetsConfigItem)
		}
		createProtectionPolicyOptions.SetCascadedTargetsConfig(cascadedTargetsConfig)
	}
	if _, ok := d.GetOk("retry_options"); ok {
		retryOptionsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetryOptions(d.Get("retry_options.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "create", "parse-retry_options").GetDiag()
		}
		createProtectionPolicyOptions.SetRetryOptions(retryOptionsModel)
	}
	if _, ok := d.GetOk("data_lock"); ok {
		createProtectionPolicyOptions.SetDataLock(d.Get("data_lock").(string))
	}
	if _, ok := d.GetOk("version"); ok {
		createProtectionPolicyOptions.SetVersion(int64(d.Get("version").(int)))
	}
	if _, ok := d.GetOk("is_cbs_enabled"); ok {
		createProtectionPolicyOptions.SetIsCBSEnabled(d.Get("is_cbs_enabled").(bool))
	}
	// if _, ok := d.GetOk("last_modification_time_usecs"); ok {
	// 	createProtectionPolicyOptions.SetLastModificationTimeUsecs(int64(d.Get("last_modification_time_usecs").(int)))
	// }
	if _, ok := d.GetOk("template_id"); ok {
		createProtectionPolicyOptions.SetTemplateID(d.Get("template_id").(string))
	}

	protectionPolicyResponse, _, err := backupRecoveryClient.CreateProtectionPolicyWithContext(context, createProtectionPolicyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateProtectionPolicyWithContext failed: %s", err.Error()), "ibm_backup_recovery_protection_policy", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	policyId := fmt.Sprintf("%s::%s", tenantId, *protectionPolicyResponse.ID)
	d.SetId(policyId)

	return resourceIbmBackupRecoveryProtectionPolicyRead(context, d, meta)
}

func resourceIbmBackupRecoveryProtectionPolicyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tenantId := d.Get("x_ibm_tenant_id").(string)
	policyId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		policyId = ParseId(d.Id(), "id")
	}
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProtectionPolicyByIdOptions := &backuprecoveryv1.GetProtectionPolicyByIdOptions{}

	getProtectionPolicyByIdOptions.SetID(policyId)
	getProtectionPolicyByIdOptions.SetXIBMTenantID(tenantId)

	protectionPolicyResponse, response, err := backupRecoveryClient.GetProtectionPolicyByIDWithContext(context, getProtectionPolicyByIdOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProtectionPolicyByIDWithContext failed: %s", err.Error()), "ibm_backup_recovery_protection_policy", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("x_ibm_tenant_id", tenantId); err != nil {
		err = fmt.Errorf("Error setting x_ibm_tenant_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-x_ibm_tenant_id").GetDiag()
	}

	if err = d.Set("name", protectionPolicyResponse.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-name").GetDiag()
	}

	if err = d.Set("policy_id", policyId); err != nil {
		err = fmt.Errorf("Error setting policy_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-policy_id").GetDiag()
	}

	backupPolicyMap, err := ResourceIbmBackupRecoveryProtectionPolicyBackupPolicyToMap(protectionPolicyResponse.BackupPolicy)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "backup_policy-to-map").GetDiag()
	}
	if err = d.Set("backup_policy", []map[string]interface{}{backupPolicyMap}); err != nil {
		err = fmt.Errorf("Error setting backup_policy: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-backup_policy").GetDiag()
	}
	if !core.IsNil(protectionPolicyResponse.Description) {
		if err = d.Set("description", protectionPolicyResponse.Description); err != nil {
			err = fmt.Errorf("Error setting description: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-description").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.BlackoutWindow) {
		blackoutWindow := []map[string]interface{}{}
		for _, blackoutWindowItem := range protectionPolicyResponse.BlackoutWindow {
			blackoutWindowItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyBlackoutWindowToMap(&blackoutWindowItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "blackout_window-to-map").GetDiag()
			}
			blackoutWindow = append(blackoutWindow, blackoutWindowItemMap)
		}
		if err = d.Set("blackout_window", blackoutWindow); err != nil {
			err = fmt.Errorf("Error setting blackout_window: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-blackout_window").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.ExtendedRetention) {
		extendedRetention := []map[string]interface{}{}
		for _, extendedRetentionItem := range protectionPolicyResponse.ExtendedRetention {
			extendedRetentionItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyExtendedRetentionPolicyToMap(&extendedRetentionItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "extended_retention-to-map").GetDiag()
			}
			extendedRetention = append(extendedRetention, extendedRetentionItemMap)
		}
		if err = d.Set("extended_retention", extendedRetention); err != nil {
			err = fmt.Errorf("Error setting extended_retention: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-extended_retention").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.RemoteTargetPolicy) {
		remoteTargetPolicyMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetsConfigurationToMap(protectionPolicyResponse.RemoteTargetPolicy)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "remote_target_policy-to-map").GetDiag()
		}
		if err = d.Set("remote_target_policy", []map[string]interface{}{remoteTargetPolicyMap}); err != nil {
			err = fmt.Errorf("Error setting remote_target_policy: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-remote_target_policy").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.CascadedTargetsConfig) {
		cascadedTargetsConfig := []map[string]interface{}{}
		for _, cascadedTargetsConfigItem := range protectionPolicyResponse.CascadedTargetsConfig {
			cascadedTargetsConfigItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCascadedTargetConfigurationToMap(&cascadedTargetsConfigItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "cascaded_targets_config-to-map").GetDiag()
			}
			cascadedTargetsConfig = append(cascadedTargetsConfig, cascadedTargetsConfigItemMap)
		}
		if err = d.Set("cascaded_targets_config", cascadedTargetsConfig); err != nil {
			err = fmt.Errorf("Error setting cascaded_targets_config: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-cascaded_targets_config").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.RetryOptions) {
		retryOptionsMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetryOptionsToMap(protectionPolicyResponse.RetryOptions)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "retry_options-to-map").GetDiag()
		}
		if err = d.Set("retry_options", []map[string]interface{}{retryOptionsMap}); err != nil {
			err = fmt.Errorf("Error setting retry_options: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-retry_options").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.DataLock) {
		if err = d.Set("data_lock", protectionPolicyResponse.DataLock); err != nil {
			err = fmt.Errorf("Error setting data_lock: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-data_lock").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.Version) {
		if err = d.Set("version", flex.IntValue(protectionPolicyResponse.Version)); err != nil {
			err = fmt.Errorf("Error setting version: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-version").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.IsCBSEnabled) {
		if err = d.Set("is_cbs_enabled", protectionPolicyResponse.IsCBSEnabled); err != nil {
			err = fmt.Errorf("Error setting is_cbs_enabled: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-is_cbs_enabled").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.LastModificationTimeUsecs) {
		if err = d.Set("last_modification_time_usecs", flex.IntValue(protectionPolicyResponse.LastModificationTimeUsecs)); err != nil {
			err = fmt.Errorf("Error setting last_modification_time_usecs: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-last_modification_time_usecs").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.TemplateID) {
		if err = d.Set("template_id", protectionPolicyResponse.TemplateID); err != nil {
			err = fmt.Errorf("Error setting template_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-template_id").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.IsUsable) {
		if err = d.Set("is_usable", protectionPolicyResponse.IsUsable); err != nil {
			err = fmt.Errorf("Error setting is_usable: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-is_usable").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.IsReplicated) {
		if err = d.Set("is_replicated", protectionPolicyResponse.IsReplicated); err != nil {
			err = fmt.Errorf("Error setting is_replicated: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-is_replicated").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.NumProtectionGroups) {
		if err = d.Set("num_protection_groups", flex.IntValue(protectionPolicyResponse.NumProtectionGroups)); err != nil {
			err = fmt.Errorf("Error setting num_protection_groups: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-num_protection_groups").GetDiag()
		}
	}
	if !core.IsNil(protectionPolicyResponse.NumProtectedObjects) {
		if err = d.Set("num_protected_objects", flex.IntValue(protectionPolicyResponse.NumProtectedObjects)); err != nil {
			err = fmt.Errorf("Error setting num_protected_objects: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "read", "set-num_protected_objects").GetDiag()
		}
	}

	return nil
}

func resourceIbmBackupRecoveryProtectionPolicyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tenantId := d.Get("x_ibm_tenant_id").(string)
	policyId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		policyId = ParseId(d.Id(), "id")
	}

	updateProtectionPolicyOptions := &backuprecoveryv1.UpdateProtectionPolicyOptions{}

	updateProtectionPolicyOptions.SetID(policyId)
	updateProtectionPolicyOptions.SetXIBMTenantID(tenantId)
	updateProtectionPolicyOptions.SetName(d.Get("name").(string))
	backupPolicy, err := ResourceIbmBackupRecoveryProtectionPolicyMapToBackupPolicy(d.Get("backup_policy.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "parse-backup_policy").GetDiag()
	}
	updateProtectionPolicyOptions.SetBackupPolicy(backupPolicy)
	if _, ok := d.GetOk("description"); ok {
		updateProtectionPolicyOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("blackout_window"); ok {
		var blackoutWindow []backuprecoveryv1.BlackoutWindow
		for _, v := range d.Get("blackout_window").([]interface{}) {
			value := v.(map[string]interface{})
			blackoutWindowItem, err := ResourceIbmBackupRecoveryProtectionPolicyMapToBlackoutWindow(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "parse-blackout_window").GetDiag()
			}
			blackoutWindow = append(blackoutWindow, *blackoutWindowItem)
		}
		updateProtectionPolicyOptions.SetBlackoutWindow(blackoutWindow)
	}
	if _, ok := d.GetOk("extended_retention"); ok {
		var extendedRetention []backuprecoveryv1.ExtendedRetentionPolicy
		for _, v := range d.Get("extended_retention").([]interface{}) {
			value := v.(map[string]interface{})
			extendedRetentionItem, err := ResourceIbmBackupRecoveryProtectionPolicyMapToExtendedRetentionPolicy(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "parse-extended_retention").GetDiag()
			}
			extendedRetention = append(extendedRetention, *extendedRetentionItem)
		}
		updateProtectionPolicyOptions.SetExtendedRetention(extendedRetention)
	}
	if _, ok := d.GetOk("remote_target_policy"); ok {
		remoteTargetPolicy, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetsConfiguration(d.Get("remote_target_policy.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "parse-remote_target_policy").GetDiag()
		}
		updateProtectionPolicyOptions.SetRemoteTargetPolicy(remoteTargetPolicy)
	}
	if _, ok := d.GetOk("cascaded_targets_config"); ok {
		var cascadedTargetsConfig []backuprecoveryv1.CascadedTargetConfiguration
		for _, v := range d.Get("cascaded_targets_config").([]interface{}) {
			value := v.(map[string]interface{})
			cascadedTargetsConfigItem, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCascadedTargetConfiguration(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "parse-cascaded_targets_config").GetDiag()
			}
			cascadedTargetsConfig = append(cascadedTargetsConfig, *cascadedTargetsConfigItem)
		}
		updateProtectionPolicyOptions.SetCascadedTargetsConfig(cascadedTargetsConfig)
	}
	if _, ok := d.GetOk("retry_options"); ok {
		retryOptions, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetryOptions(d.Get("retry_options.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "update", "parse-retry_options").GetDiag()
		}
		updateProtectionPolicyOptions.SetRetryOptions(retryOptions)
	}
	if _, ok := d.GetOk("data_lock"); ok {
		updateProtectionPolicyOptions.SetDataLock(d.Get("data_lock").(string))
	}
	if _, ok := d.GetOk("version"); ok {
		updateProtectionPolicyOptions.SetVersion(int64(d.Get("version").(int)))
	}
	if _, ok := d.GetOk("is_cbs_enabled"); ok {
		updateProtectionPolicyOptions.SetIsCBSEnabled(d.Get("is_cbs_enabled").(bool))
	}
	if _, ok := d.GetOk("last_modification_time_usecs"); ok {
		updateProtectionPolicyOptions.SetLastModificationTimeUsecs(int64(d.Get("last_modification_time_usecs").(int)))
	}
	if _, ok := d.GetOk("template_id"); ok {
		updateProtectionPolicyOptions.SetTemplateID(d.Get("template_id").(string))
	}

	_, _, err = backupRecoveryClient.UpdateProtectionPolicyWithContext(context, updateProtectionPolicyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateProtectionPolicyWithContext failed: %s", err.Error()), "ibm_backup_recovery_protection_policy", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	return resourceIbmBackupRecoveryProtectionPolicyRead(context, d, meta)
}

func resourceIbmBackupRecoveryProtectionPolicyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	backupRecoveryClient, err := meta.(conns.ClientSession).BackupRecoveryV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_backup_recovery_protection_policy", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tenantId := d.Get("x_ibm_tenant_id").(string)
	policyId := d.Id()
	if strings.Contains(d.Id(), "::") {
		tenantId = ParseId(d.Id(), "tenantId")
		policyId = ParseId(d.Id(), "id")
	}

	deleteProtectionPolicyOptions := &backuprecoveryv1.DeleteProtectionPolicyOptions{}

	deleteProtectionPolicyOptions.SetID(policyId)
	deleteProtectionPolicyOptions.SetXIBMTenantID(tenantId)

	_, err = backupRecoveryClient.DeleteProtectionPolicyWithContext(context, deleteProtectionPolicyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteProtectionPolicyWithContext failed: %s", err.Error()), "ibm_backup_recovery_protection_policy", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.BackupPolicy, error) {
	model := &backuprecoveryv1.BackupPolicy{}
	RegularModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRegularBackupPolicy(modelMap["regular"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Regular = RegularModel
	if modelMap["log"] != nil && len(modelMap["log"].([]interface{})) > 0 {
		LogModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogBackupPolicy(modelMap["log"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Log = LogModel
	}
	if modelMap["bmr"] != nil && len(modelMap["bmr"].([]interface{})) > 0 {
		BmrModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToBmrBackupPolicy(modelMap["bmr"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Bmr = BmrModel
	}
	if modelMap["cdp"] != nil && len(modelMap["cdp"].([]interface{})) > 0 {
		CdpModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCdpBackupPolicy(modelMap["cdp"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Cdp = CdpModel
	}
	if modelMap["storage_array_snapshot"] != nil && len(modelMap["storage_array_snapshot"].([]interface{})) > 0 {
		StorageArraySnapshotModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToStorageArraySnapshotBackupPolicy(modelMap["storage_array_snapshot"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.StorageArraySnapshot = StorageArraySnapshotModel
	}
	if modelMap["run_timeouts"] != nil {
		runTimeouts := []backuprecoveryv1.CancellationTimeoutParams{}
		for _, runTimeoutsItem := range modelMap["run_timeouts"].([]interface{}) {
			runTimeoutsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(runTimeoutsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			runTimeouts = append(runTimeouts, *runTimeoutsItemModel)
		}
		model.RunTimeouts = runTimeouts
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToRegularBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.RegularBackupPolicy, error) {
	model := &backuprecoveryv1.RegularBackupPolicy{}
	if modelMap["incremental"] != nil && len(modelMap["incremental"].([]interface{})) > 0 {
		IncrementalModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToIncrementalBackupPolicy(modelMap["incremental"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Incremental = IncrementalModel
	}
	if modelMap["full"] != nil && len(modelMap["full"].([]interface{})) > 0 {
		FullModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToFullBackupPolicy(modelMap["full"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Full = FullModel
	}
	if modelMap["full_backups"] != nil {
		fullBackups := []backuprecoveryv1.FullScheduleAndRetention{}
		for _, fullBackupsItem := range modelMap["full_backups"].([]interface{}) {
			fullBackupsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToFullScheduleAndRetention(fullBackupsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			fullBackups = append(fullBackups, *fullBackupsItemModel)
		}
		model.FullBackups = fullBackups
	}
	if modelMap["retention"] != nil && len(modelMap["retention"].([]interface{})) > 0 {
		RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Retention = RetentionModel
	}
	if modelMap["primary_backup_target"] != nil && len(modelMap["primary_backup_target"].([]interface{})) > 0 {
		PrimaryBackupTargetModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToPrimaryBackupTarget(modelMap["primary_backup_target"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryBackupTarget = PrimaryBackupTargetModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToIncrementalBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.IncrementalBackupPolicy, error) {
	model := &backuprecoveryv1.IncrementalBackupPolicy{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToIncrementalSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToIncrementalSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.IncrementalSchedule, error) {
	model := &backuprecoveryv1.IncrementalSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["minute_schedule"] != nil && len(modelMap["minute_schedule"].([]interface{})) > 0 {
		MinuteScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMinuteSchedule(modelMap["minute_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MinuteSchedule = MinuteScheduleModel
	}
	if modelMap["hour_schedule"] != nil && len(modelMap["hour_schedule"].([]interface{})) > 0 {
		HourScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToHourSchedule(modelMap["hour_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.HourSchedule = HourScheduleModel
	}
	if modelMap["day_schedule"] != nil && len(modelMap["day_schedule"].([]interface{})) > 0 {
		DayScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDaySchedule(modelMap["day_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DaySchedule = DayScheduleModel
	}
	if modelMap["week_schedule"] != nil && len(modelMap["week_schedule"].([]interface{})) > 0 {
		WeekScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToWeekSchedule(modelMap["week_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WeekSchedule = WeekScheduleModel
	}
	if modelMap["month_schedule"] != nil && len(modelMap["month_schedule"].([]interface{})) > 0 {
		MonthScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMonthSchedule(modelMap["month_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MonthSchedule = MonthScheduleModel
	}
	if modelMap["year_schedule"] != nil && len(modelMap["year_schedule"].([]interface{})) > 0 {
		YearScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToYearSchedule(modelMap["year_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.YearSchedule = YearScheduleModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToMinuteSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.MinuteSchedule, error) {
	model := &backuprecoveryv1.MinuteSchedule{}
	model.Frequency = core.Int64Ptr(int64(modelMap["frequency"].(int)))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToHourSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.HourSchedule, error) {
	model := &backuprecoveryv1.HourSchedule{}
	model.Frequency = core.Int64Ptr(int64(modelMap["frequency"].(int)))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToDaySchedule(modelMap map[string]interface{}) (*backuprecoveryv1.DaySchedule, error) {
	model := &backuprecoveryv1.DaySchedule{}
	model.Frequency = core.Int64Ptr(int64(modelMap["frequency"].(int)))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToWeekSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.WeekSchedule, error) {
	model := &backuprecoveryv1.WeekSchedule{}
	dayOfWeek := []string{}
	for _, dayOfWeekItem := range modelMap["day_of_week"].([]interface{}) {
		dayOfWeek = append(dayOfWeek, dayOfWeekItem.(string))
	}
	model.DayOfWeek = dayOfWeek
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToMonthSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.MonthSchedule, error) {
	model := &backuprecoveryv1.MonthSchedule{}
	if modelMap["day_of_week"] != nil {
		dayOfWeek := []string{}
		for _, dayOfWeekItem := range modelMap["day_of_week"].([]interface{}) {
			dayOfWeek = append(dayOfWeek, dayOfWeekItem.(string))
		}
		model.DayOfWeek = dayOfWeek
	}
	if modelMap["week_of_month"] != nil && modelMap["week_of_month"].(string) != "" {
		model.WeekOfMonth = core.StringPtr(modelMap["week_of_month"].(string))
	}
	if modelMap["day_of_month"] != nil {
		model.DayOfMonth = core.Int64Ptr(int64(modelMap["day_of_month"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToYearSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.YearSchedule, error) {
	model := &backuprecoveryv1.YearSchedule{}
	model.DayOfYear = core.StringPtr(modelMap["day_of_year"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToFullBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.FullBackupPolicy, error) {
	model := &backuprecoveryv1.FullBackupPolicy{}
	if modelMap["schedule"] != nil && len(modelMap["schedule"].([]interface{})) > 0 {
		ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToFullSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Schedule = ScheduleModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToFullSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.FullSchedule, error) {
	model := &backuprecoveryv1.FullSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["day_schedule"] != nil && len(modelMap["day_schedule"].([]interface{})) > 0 {
		DayScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDaySchedule(modelMap["day_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DaySchedule = DayScheduleModel
	}
	if modelMap["week_schedule"] != nil && len(modelMap["week_schedule"].([]interface{})) > 0 {
		WeekScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToWeekSchedule(modelMap["week_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WeekSchedule = WeekScheduleModel
	}
	if modelMap["month_schedule"] != nil && len(modelMap["month_schedule"].([]interface{})) > 0 {
		MonthScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMonthSchedule(modelMap["month_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MonthSchedule = MonthScheduleModel
	}
	if modelMap["year_schedule"] != nil && len(modelMap["year_schedule"].([]interface{})) > 0 {
		YearScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToYearSchedule(modelMap["year_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.YearSchedule = YearScheduleModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToFullScheduleAndRetention(modelMap map[string]interface{}) (*backuprecoveryv1.FullScheduleAndRetention, error) {
	model := &backuprecoveryv1.FullScheduleAndRetention{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToFullSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap map[string]interface{}) (*backuprecoveryv1.Retention, error) {
	model := &backuprecoveryv1.Retention{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["data_lock_config"] != nil && len(modelMap["data_lock_config"].([]interface{})) > 0 {
		DataLockConfigModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDataLockConfig(modelMap["data_lock_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DataLockConfig = DataLockConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToDataLockConfig(modelMap map[string]interface{}) (*backuprecoveryv1.DataLockConfig, error) {
	model := &backuprecoveryv1.DataLockConfig{}
	model.Mode = core.StringPtr(modelMap["mode"].(string))
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["enable_worm_on_external_target"] != nil {
		model.EnableWormOnExternalTarget = core.BoolPtr(modelMap["enable_worm_on_external_target"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToPrimaryBackupTarget(modelMap map[string]interface{}) (*backuprecoveryv1.PrimaryBackupTarget, error) {
	model := &backuprecoveryv1.PrimaryBackupTarget{}
	if modelMap["target_type"] != nil && modelMap["target_type"].(string) != "" {
		model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	}
	if modelMap["archival_target_settings"] != nil && len(modelMap["archival_target_settings"].([]interface{})) > 0 {
		ArchivalTargetSettingsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToPrimaryArchivalTarget(modelMap["archival_target_settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ArchivalTargetSettings = ArchivalTargetSettingsModel
	}
	if modelMap["use_default_backup_target"] != nil {
		model.UseDefaultBackupTarget = core.BoolPtr(modelMap["use_default_backup_target"].(bool))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToPrimaryArchivalTarget(modelMap map[string]interface{}) (*backuprecoveryv1.PrimaryArchivalTarget, error) {
	model := &backuprecoveryv1.PrimaryArchivalTarget{}
	model.TargetID = core.Int64Ptr(int64(modelMap["target_id"].(int)))
	if modelMap["target_name"] != nil && modelMap["target_name"].(string) != "" {
		model.TargetName = core.StringPtr(modelMap["target_name"].(string))
	}
	if modelMap["tier_settings"] != nil && len(modelMap["tier_settings"].([]interface{})) > 0 {
		TierSettingsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTierLevelSettings(modelMap["tier_settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.TierSettings = TierSettingsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToTierLevelSettings(modelMap map[string]interface{}) (*backuprecoveryv1.TierLevelSettings, error) {
	model := &backuprecoveryv1.TierLevelSettings{}
	if modelMap["aws_tiering"] != nil && len(modelMap["aws_tiering"].([]interface{})) > 0 {
		AwsTieringModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAWSTiers(modelMap["aws_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AwsTiering = AwsTieringModel
	}
	if modelMap["azure_tiering"] != nil && len(modelMap["azure_tiering"].([]interface{})) > 0 {
		AzureTieringModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAzureTiers(modelMap["azure_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AzureTiering = AzureTieringModel
	}
	if modelMap["cloud_platform"] != nil && modelMap["cloud_platform"].(string) != "" {
		model.CloudPlatform = core.StringPtr(modelMap["cloud_platform"].(string))
	}
	if modelMap["google_tiering"] != nil && len(modelMap["google_tiering"].([]interface{})) > 0 {
		GoogleTieringModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToGoogleTiers(modelMap["google_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.GoogleTiering = GoogleTieringModel
	}
	if modelMap["oracle_tiering"] != nil && len(modelMap["oracle_tiering"].([]interface{})) > 0 {
		OracleTieringModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToOracleTiers(modelMap["oracle_tiering"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.OracleTiering = OracleTieringModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAWSTiers(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTiers, error) {
	model := &backuprecoveryv1.AWSTiers{}
	tiers := []backuprecoveryv1.AWSTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAWSTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAWSTier(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTier, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyMapToAzureTiers(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTiers, error) {
	model := &backuprecoveryv1.AzureTiers{}
	if modelMap["tiers"] != nil {
		tiers := []backuprecoveryv1.AzureTier{}
		for _, tiersItem := range modelMap["tiers"].([]interface{}) {
			tiersItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAzureTier(tiersItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			tiers = append(tiers, *tiersItemModel)
		}
		model.Tiers = tiers
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAzureTier(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTier, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyMapToGoogleTiers(modelMap map[string]interface{}) (*backuprecoveryv1.GoogleTiers, error) {
	model := &backuprecoveryv1.GoogleTiers{}
	tiers := []backuprecoveryv1.GoogleTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToGoogleTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToGoogleTier(modelMap map[string]interface{}) (*backuprecoveryv1.GoogleTier, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyMapToOracleTiers(modelMap map[string]interface{}) (*backuprecoveryv1.OracleTiers, error) {
	model := &backuprecoveryv1.OracleTiers{}
	tiers := []backuprecoveryv1.OracleTier{}
	for _, tiersItem := range modelMap["tiers"].([]interface{}) {
		tiersItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToOracleTier(tiersItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		tiers = append(tiers, *tiersItemModel)
	}
	model.Tiers = tiers
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToOracleTier(modelMap map[string]interface{}) (*backuprecoveryv1.OracleTier, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyMapToLogBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.LogBackupPolicy, error) {
	model := &backuprecoveryv1.LogBackupPolicy{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToLogSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.LogSchedule, error) {
	model := &backuprecoveryv1.LogSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["minute_schedule"] != nil && len(modelMap["minute_schedule"].([]interface{})) > 0 {
		MinuteScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMinuteSchedule(modelMap["minute_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MinuteSchedule = MinuteScheduleModel
	}
	if modelMap["hour_schedule"] != nil && len(modelMap["hour_schedule"].([]interface{})) > 0 {
		HourScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToHourSchedule(modelMap["hour_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.HourSchedule = HourScheduleModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToBmrBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.BmrBackupPolicy, error) {
	model := &backuprecoveryv1.BmrBackupPolicy{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToBmrSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToBmrSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.BmrSchedule, error) {
	model := &backuprecoveryv1.BmrSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["day_schedule"] != nil && len(modelMap["day_schedule"].([]interface{})) > 0 {
		DayScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDaySchedule(modelMap["day_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DaySchedule = DayScheduleModel
	}
	if modelMap["week_schedule"] != nil && len(modelMap["week_schedule"].([]interface{})) > 0 {
		WeekScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToWeekSchedule(modelMap["week_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WeekSchedule = WeekScheduleModel
	}
	if modelMap["month_schedule"] != nil && len(modelMap["month_schedule"].([]interface{})) > 0 {
		MonthScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMonthSchedule(modelMap["month_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MonthSchedule = MonthScheduleModel
	}
	if modelMap["year_schedule"] != nil && len(modelMap["year_schedule"].([]interface{})) > 0 {
		YearScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToYearSchedule(modelMap["year_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.YearSchedule = YearScheduleModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCdpBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.CdpBackupPolicy, error) {
	model := &backuprecoveryv1.CdpBackupPolicy{}
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCdpRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCdpRetention(modelMap map[string]interface{}) (*backuprecoveryv1.CdpRetention, error) {
	model := &backuprecoveryv1.CdpRetention{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["data_lock_config"] != nil && len(modelMap["data_lock_config"].([]interface{})) > 0 {
		DataLockConfigModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDataLockConfig(modelMap["data_lock_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DataLockConfig = DataLockConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToStorageArraySnapshotBackupPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.StorageArraySnapshotBackupPolicy, error) {
	model := &backuprecoveryv1.StorageArraySnapshotBackupPolicy{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToStorageArraySnapshotSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToStorageArraySnapshotSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.StorageArraySnapshotSchedule, error) {
	model := &backuprecoveryv1.StorageArraySnapshotSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["minute_schedule"] != nil && len(modelMap["minute_schedule"].([]interface{})) > 0 {
		MinuteScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMinuteSchedule(modelMap["minute_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MinuteSchedule = MinuteScheduleModel
	}
	if modelMap["hour_schedule"] != nil && len(modelMap["hour_schedule"].([]interface{})) > 0 {
		HourScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToHourSchedule(modelMap["hour_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.HourSchedule = HourScheduleModel
	}
	if modelMap["day_schedule"] != nil && len(modelMap["day_schedule"].([]interface{})) > 0 {
		DayScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDaySchedule(modelMap["day_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DaySchedule = DayScheduleModel
	}
	if modelMap["week_schedule"] != nil && len(modelMap["week_schedule"].([]interface{})) > 0 {
		WeekScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToWeekSchedule(modelMap["week_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.WeekSchedule = WeekScheduleModel
	}
	if modelMap["month_schedule"] != nil && len(modelMap["month_schedule"].([]interface{})) > 0 {
		MonthScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToMonthSchedule(modelMap["month_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.MonthSchedule = MonthScheduleModel
	}
	if modelMap["year_schedule"] != nil && len(modelMap["year_schedule"].([]interface{})) > 0 {
		YearScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToYearSchedule(modelMap["year_schedule"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.YearSchedule = YearScheduleModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(modelMap map[string]interface{}) (*backuprecoveryv1.CancellationTimeoutParams, error) {
	model := &backuprecoveryv1.CancellationTimeoutParams{}
	if modelMap["timeout_mins"] != nil {
		model.TimeoutMins = core.Int64Ptr(int64(modelMap["timeout_mins"].(int)))
	}
	if modelMap["backup_type"] != nil && modelMap["backup_type"].(string) != "" {
		model.BackupType = core.StringPtr(modelMap["backup_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToBlackoutWindow(modelMap map[string]interface{}) (*backuprecoveryv1.BlackoutWindow, error) {
	model := &backuprecoveryv1.BlackoutWindow{}
	model.Day = core.StringPtr(modelMap["day"].(string))
	StartTimeModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTimeOfDay(modelMap["start_time"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.StartTime = StartTimeModel
	EndTimeModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTimeOfDay(modelMap["end_time"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.EndTime = EndTimeModel
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToTimeOfDay(modelMap map[string]interface{}) (*backuprecoveryv1.TimeOfDay, error) {
	model := &backuprecoveryv1.TimeOfDay{}
	model.Hour = core.Int64Ptr(int64(modelMap["hour"].(int)))
	model.Minute = core.Int64Ptr(int64(modelMap["minute"].(int)))
	if modelMap["time_zone"] != nil && modelMap["time_zone"].(string) != "" {
		model.TimeZone = core.StringPtr(modelMap["time_zone"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToExtendedRetentionPolicy(modelMap map[string]interface{}) (*backuprecoveryv1.ExtendedRetentionPolicy, error) {
	model := &backuprecoveryv1.ExtendedRetentionPolicy{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToExtendedRetentionSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	if modelMap["run_type"] != nil && modelMap["run_type"].(string) != "" {
		model.RunType = core.StringPtr(modelMap["run_type"].(string))
	}
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToExtendedRetentionSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.ExtendedRetentionSchedule, error) {
	model := &backuprecoveryv1.ExtendedRetentionSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["frequency"] != nil {
		model.Frequency = core.Int64Ptr(int64(modelMap["frequency"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToTargetsConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.TargetsConfiguration, error) {
	model := &backuprecoveryv1.TargetsConfiguration{}
	if modelMap["replication_targets"] != nil {
		replicationTargets := []backuprecoveryv1.ReplicationTargetConfiguration{}
		for _, replicationTargetsItem := range modelMap["replication_targets"].([]interface{}) {
			replicationTargetsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToReplicationTargetConfiguration(replicationTargetsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			replicationTargets = append(replicationTargets, *replicationTargetsItemModel)
		}
		model.ReplicationTargets = replicationTargets
	}
	if modelMap["archival_targets"] != nil {
		archivalTargets := []backuprecoveryv1.ArchivalTargetConfiguration{}
		for _, archivalTargetsItem := range modelMap["archival_targets"].([]interface{}) {
			archivalTargetsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToArchivalTargetConfiguration(archivalTargetsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			archivalTargets = append(archivalTargets, *archivalTargetsItemModel)
		}
		model.ArchivalTargets = archivalTargets
	}
	if modelMap["cloud_spin_targets"] != nil {
		cloudSpinTargets := []backuprecoveryv1.CloudSpinTargetConfiguration{}
		for _, cloudSpinTargetsItem := range modelMap["cloud_spin_targets"].([]interface{}) {
			cloudSpinTargetsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCloudSpinTargetConfiguration(cloudSpinTargetsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			cloudSpinTargets = append(cloudSpinTargets, *cloudSpinTargetsItemModel)
		}
		model.CloudSpinTargets = cloudSpinTargets
	}
	if modelMap["onprem_deploy_targets"] != nil {
		onpremDeployTargets := []backuprecoveryv1.OnpremDeployTargetConfiguration{}
		for _, onpremDeployTargetsItem := range modelMap["onprem_deploy_targets"].([]interface{}) {
			onpremDeployTargetsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToOnpremDeployTargetConfiguration(onpremDeployTargetsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			onpremDeployTargets = append(onpremDeployTargets, *onpremDeployTargetsItemModel)
		}
		model.OnpremDeployTargets = onpremDeployTargets
	}
	if modelMap["rpaas_targets"] != nil {
		rpaasTargets := []backuprecoveryv1.RpaasTargetConfiguration{}
		for _, rpaasTargetsItem := range modelMap["rpaas_targets"].([]interface{}) {
			rpaasTargetsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRpaasTargetConfiguration(rpaasTargetsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			rpaasTargets = append(rpaasTargets, *rpaasTargetsItemModel)
		}
		model.RpaasTargets = rpaasTargets
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToReplicationTargetConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.ReplicationTargetConfiguration, error) {
	model := &backuprecoveryv1.ReplicationTargetConfiguration{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	if modelMap["copy_on_run_success"] != nil {
		model.CopyOnRunSuccess = core.BoolPtr(modelMap["copy_on_run_success"].(bool))
	}
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	if modelMap["backup_run_type"] != nil && modelMap["backup_run_type"].(string) != "" {
		model.BackupRunType = core.StringPtr(modelMap["backup_run_type"].(string))
	}
	if modelMap["run_timeouts"] != nil {
		runTimeouts := []backuprecoveryv1.CancellationTimeoutParams{}
		for _, runTimeoutsItem := range modelMap["run_timeouts"].([]interface{}) {
			runTimeoutsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(runTimeoutsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			runTimeouts = append(runTimeouts, *runTimeoutsItemModel)
		}
		model.RunTimeouts = runTimeouts
	}
	if modelMap["log_retention"] != nil && len(modelMap["log_retention"].([]interface{})) > 0 {
		LogRetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogRetention(modelMap["log_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRetention = LogRetentionModel
	}
	if modelMap["aws_target_config"] != nil && len(modelMap["aws_target_config"].([]interface{})) > 0 {
		AwsTargetConfigModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAWSTargetConfig(modelMap["aws_target_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AwsTargetConfig = AwsTargetConfigModel
	}
	if modelMap["azure_target_config"] != nil && len(modelMap["azure_target_config"].([]interface{})) > 0 {
		AzureTargetConfigModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAzureTargetConfig(modelMap["azure_target_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AzureTargetConfig = AzureTargetConfigModel
	}
	model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	if modelMap["remote_target_config"] != nil && len(modelMap["remote_target_config"].([]interface{})) > 0 {
		RemoteTargetConfigModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRemoteTargetConfig(modelMap["remote_target_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.RemoteTargetConfig = RemoteTargetConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToTargetSchedule(modelMap map[string]interface{}) (*backuprecoveryv1.TargetSchedule, error) {
	model := &backuprecoveryv1.TargetSchedule{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	if modelMap["frequency"] != nil {
		model.Frequency = core.Int64Ptr(int64(modelMap["frequency"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToLogRetention(modelMap map[string]interface{}) (*backuprecoveryv1.LogRetention, error) {
	model := &backuprecoveryv1.LogRetention{}
	model.Unit = core.StringPtr(modelMap["unit"].(string))
	model.Duration = core.Int64Ptr(int64(modelMap["duration"].(int)))
	if modelMap["data_lock_config"] != nil && len(modelMap["data_lock_config"].([]interface{})) > 0 {
		DataLockConfigModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToDataLockConfig(modelMap["data_lock_config"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DataLockConfig = DataLockConfigModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAWSTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.AWSTargetConfig, error) {
	model := &backuprecoveryv1.AWSTargetConfig{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Region = core.Int64Ptr(int64(modelMap["region"].(int)))
	if modelMap["region_name"] != nil && modelMap["region_name"].(string) != "" {
		model.RegionName = core.StringPtr(modelMap["region_name"].(string))
	}
	model.SourceID = core.Int64Ptr(int64(modelMap["source_id"].(int)))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAzureTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.AzureTargetConfig, error) {
	model := &backuprecoveryv1.AzureTargetConfig{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["resource_group"] != nil {
		model.ResourceGroup = core.Int64Ptr(int64(modelMap["resource_group"].(int)))
	}
	if modelMap["resource_group_name"] != nil && modelMap["resource_group_name"].(string) != "" {
		model.ResourceGroupName = core.StringPtr(modelMap["resource_group_name"].(string))
	}
	model.SourceID = core.Int64Ptr(int64(modelMap["source_id"].(int)))
	if modelMap["storage_account"] != nil && modelMap["storage_account"].(int) != 0 {
		model.StorageAccount = core.Int64Ptr(int64(modelMap["storage_account"].(int)))
	}
	if modelMap["storage_account_name"] != nil && modelMap["storage_account_name"].(string) != "" {
		model.StorageAccountName = core.StringPtr(modelMap["storage_account_name"].(string))
	}
	if modelMap["storage_container"] != nil && modelMap["storage_container"].(int) != 0 {
		model.StorageContainer = core.Int64Ptr(int64(modelMap["storage_container"].(int)))
	}
	if modelMap["storage_container_name"] != nil && modelMap["storage_container_name"].(string) != "" {
		model.StorageContainerName = core.StringPtr(modelMap["storage_container_name"].(string))
	}
	if modelMap["storage_resource_group"] != nil && modelMap["storage_resource_group"].(int) != 0 {
		model.StorageResourceGroup = core.Int64Ptr(int64(modelMap["storage_resource_group"].(int)))
	}
	if modelMap["storage_resource_group_name"] != nil && modelMap["storage_resource_group_name"].(string) != "" {
		model.StorageResourceGroupName = core.StringPtr(modelMap["storage_resource_group_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToRemoteTargetConfig(modelMap map[string]interface{}) (*backuprecoveryv1.RemoteTargetConfig, error) {
	model := &backuprecoveryv1.RemoteTargetConfig{}
	model.ClusterID = core.Int64Ptr(int64(modelMap["cluster_id"].(int)))
	if modelMap["cluster_name"] != nil && modelMap["cluster_name"].(string) != "" {
		model.ClusterName = core.StringPtr(modelMap["cluster_name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToArchivalTargetConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.ArchivalTargetConfiguration, error) {
	model := &backuprecoveryv1.ArchivalTargetConfiguration{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	if modelMap["copy_on_run_success"] != nil {
		model.CopyOnRunSuccess = core.BoolPtr(modelMap["copy_on_run_success"].(bool))
	}
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	if modelMap["backup_run_type"] != nil && modelMap["backup_run_type"].(string) != "" {
		model.BackupRunType = core.StringPtr(modelMap["backup_run_type"].(string))
	}
	if modelMap["run_timeouts"] != nil {
		runTimeouts := []backuprecoveryv1.CancellationTimeoutParams{}
		for _, runTimeoutsItem := range modelMap["run_timeouts"].([]interface{}) {
			runTimeoutsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(runTimeoutsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			runTimeouts = append(runTimeouts, *runTimeoutsItemModel)
		}
		model.RunTimeouts = runTimeouts
	}
	if modelMap["log_retention"] != nil && len(modelMap["log_retention"].([]interface{})) > 0 {
		LogRetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogRetention(modelMap["log_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRetention = LogRetentionModel
	}
	model.TargetID = core.Int64Ptr(int64(modelMap["target_id"].(int)))
	if modelMap["target_name"] != nil && modelMap["target_name"].(string) != "" {
		model.TargetName = core.StringPtr(modelMap["target_name"].(string))
	}
	if modelMap["target_type"] != nil && modelMap["target_type"].(string) != "" {
		model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	}
	if modelMap["tier_settings"] != nil && len(modelMap["tier_settings"].([]interface{})) > 0 {
		TierSettingsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTierLevelSettings(modelMap["tier_settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.TierSettings = TierSettingsModel
	}
	if modelMap["extended_retention"] != nil {
		extendedRetention := []backuprecoveryv1.ExtendedRetentionPolicy{}
		for _, extendedRetentionItem := range modelMap["extended_retention"].([]interface{}) {
			extendedRetentionItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToExtendedRetentionPolicy(extendedRetentionItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			extendedRetention = append(extendedRetention, *extendedRetentionItemModel)
		}
		model.ExtendedRetention = extendedRetention
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCloudSpinTargetConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.CloudSpinTargetConfiguration, error) {
	model := &backuprecoveryv1.CloudSpinTargetConfiguration{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	if modelMap["copy_on_run_success"] != nil {
		model.CopyOnRunSuccess = core.BoolPtr(modelMap["copy_on_run_success"].(bool))
	}
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	if modelMap["backup_run_type"] != nil && modelMap["backup_run_type"].(string) != "" {
		model.BackupRunType = core.StringPtr(modelMap["backup_run_type"].(string))
	}
	if modelMap["run_timeouts"] != nil {
		runTimeouts := []backuprecoveryv1.CancellationTimeoutParams{}
		for _, runTimeoutsItem := range modelMap["run_timeouts"].([]interface{}) {
			runTimeoutsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(runTimeoutsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			runTimeouts = append(runTimeouts, *runTimeoutsItemModel)
		}
		model.RunTimeouts = runTimeouts
	}
	if modelMap["log_retention"] != nil && len(modelMap["log_retention"].([]interface{})) > 0 {
		LogRetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogRetention(modelMap["log_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRetention = LogRetentionModel
	}
	TargetModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCloudSpinTarget(modelMap["target"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Target = TargetModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCloudSpinTarget(modelMap map[string]interface{}) (*backuprecoveryv1.CloudSpinTarget, error) {
	model := &backuprecoveryv1.CloudSpinTarget{}
	if modelMap["aws_params"] != nil && len(modelMap["aws_params"].([]interface{})) > 0 {
		AwsParamsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAwsCloudSpinParams(modelMap["aws_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AwsParams = AwsParamsModel
	}
	if modelMap["azure_params"] != nil && len(modelMap["azure_params"].([]interface{})) > 0 {
		AzureParamsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToAzureCloudSpinParams(modelMap["azure_params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AzureParams = AzureParamsModel
	}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAwsCloudSpinParams(modelMap map[string]interface{}) (*backuprecoveryv1.AwsCloudSpinParams, error) {
	model := &backuprecoveryv1.AwsCloudSpinParams{}
	if modelMap["custom_tag_list"] != nil {
		customTagList := []backuprecoveryv1.CustomTagParams{}
		for _, customTagListItem := range modelMap["custom_tag_list"].([]interface{}) {
			customTagListItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCustomTagParams(customTagListItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			customTagList = append(customTagList, *customTagListItemModel)
		}
		model.CustomTagList = customTagList
	}
	model.Region = core.Int64Ptr(int64(modelMap["region"].(int)))
	if modelMap["subnet_id"] != nil {
		model.SubnetID = core.Int64Ptr(int64(modelMap["subnet_id"].(int)))
	}
	if modelMap["vpc_id"] != nil {
		model.VpcID = core.Int64Ptr(int64(modelMap["vpc_id"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCustomTagParams(modelMap map[string]interface{}) (*backuprecoveryv1.CustomTagParams, error) {
	model := &backuprecoveryv1.CustomTagParams{}
	model.Key = core.StringPtr(modelMap["key"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToAzureCloudSpinParams(modelMap map[string]interface{}) (*backuprecoveryv1.AzureCloudSpinParams, error) {
	model := &backuprecoveryv1.AzureCloudSpinParams{}
	if modelMap["availability_set_id"] != nil {
		model.AvailabilitySetID = core.Int64Ptr(int64(modelMap["availability_set_id"].(int)))
	}
	if modelMap["network_resource_group_id"] != nil {
		model.NetworkResourceGroupID = core.Int64Ptr(int64(modelMap["network_resource_group_id"].(int)))
	}
	if modelMap["resource_group_id"] != nil {
		model.ResourceGroupID = core.Int64Ptr(int64(modelMap["resource_group_id"].(int)))
	}
	if modelMap["storage_account_id"] != nil {
		model.StorageAccountID = core.Int64Ptr(int64(modelMap["storage_account_id"].(int)))
	}
	if modelMap["storage_container_id"] != nil {
		model.StorageContainerID = core.Int64Ptr(int64(modelMap["storage_container_id"].(int)))
	}
	if modelMap["storage_resource_group_id"] != nil {
		model.StorageResourceGroupID = core.Int64Ptr(int64(modelMap["storage_resource_group_id"].(int)))
	}
	if modelMap["temp_vm_resource_group_id"] != nil {
		model.TempVmResourceGroupID = core.Int64Ptr(int64(modelMap["temp_vm_resource_group_id"].(int)))
	}
	if modelMap["temp_vm_storage_account_id"] != nil {
		model.TempVmStorageAccountID = core.Int64Ptr(int64(modelMap["temp_vm_storage_account_id"].(int)))
	}
	if modelMap["temp_vm_storage_container_id"] != nil {
		model.TempVmStorageContainerID = core.Int64Ptr(int64(modelMap["temp_vm_storage_container_id"].(int)))
	}
	if modelMap["temp_vm_subnet_id"] != nil {
		model.TempVmSubnetID = core.Int64Ptr(int64(modelMap["temp_vm_subnet_id"].(int)))
	}
	if modelMap["temp_vm_virtual_network_id"] != nil {
		model.TempVmVirtualNetworkID = core.Int64Ptr(int64(modelMap["temp_vm_virtual_network_id"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToOnpremDeployTargetConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.OnpremDeployTargetConfiguration, error) {
	model := &backuprecoveryv1.OnpremDeployTargetConfiguration{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	if modelMap["copy_on_run_success"] != nil {
		model.CopyOnRunSuccess = core.BoolPtr(modelMap["copy_on_run_success"].(bool))
	}
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	if modelMap["backup_run_type"] != nil && modelMap["backup_run_type"].(string) != "" {
		model.BackupRunType = core.StringPtr(modelMap["backup_run_type"].(string))
	}
	if modelMap["run_timeouts"] != nil {
		runTimeouts := []backuprecoveryv1.CancellationTimeoutParams{}
		for _, runTimeoutsItem := range modelMap["run_timeouts"].([]interface{}) {
			runTimeoutsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(runTimeoutsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			runTimeouts = append(runTimeouts, *runTimeoutsItemModel)
		}
		model.RunTimeouts = runTimeouts
	}
	if modelMap["log_retention"] != nil && len(modelMap["log_retention"].([]interface{})) > 0 {
		LogRetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogRetention(modelMap["log_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRetention = LogRetentionModel
	}
	if modelMap["params"] != nil && len(modelMap["params"].([]interface{})) > 0 {
		ParamsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToOnpremDeployParams(modelMap["params"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Params = ParamsModel
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToOnpremDeployParams(modelMap map[string]interface{}) (*backuprecoveryv1.OnpremDeployParams, error) {
	model := &backuprecoveryv1.OnpremDeployParams{}
	if modelMap["id"] != nil {
		model.ID = core.Int64Ptr(int64(modelMap["id"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToRpaasTargetConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.RpaasTargetConfiguration, error) {
	model := &backuprecoveryv1.RpaasTargetConfiguration{}
	ScheduleModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetSchedule(modelMap["schedule"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Schedule = ScheduleModel
	RetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToRetention(modelMap["retention"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Retention = RetentionModel
	if modelMap["copy_on_run_success"] != nil {
		model.CopyOnRunSuccess = core.BoolPtr(modelMap["copy_on_run_success"].(bool))
	}
	if modelMap["config_id"] != nil && modelMap["config_id"].(string) != "" {
		model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	}
	if modelMap["backup_run_type"] != nil && modelMap["backup_run_type"].(string) != "" {
		model.BackupRunType = core.StringPtr(modelMap["backup_run_type"].(string))
	}
	if modelMap["run_timeouts"] != nil {
		runTimeouts := []backuprecoveryv1.CancellationTimeoutParams{}
		for _, runTimeoutsItem := range modelMap["run_timeouts"].([]interface{}) {
			runTimeoutsItemModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToCancellationTimeoutParams(runTimeoutsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			runTimeouts = append(runTimeouts, *runTimeoutsItemModel)
		}
		model.RunTimeouts = runTimeouts
	}
	if modelMap["log_retention"] != nil && len(modelMap["log_retention"].([]interface{})) > 0 {
		LogRetentionModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToLogRetention(modelMap["log_retention"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.LogRetention = LogRetentionModel
	}
	model.TargetID = core.Int64Ptr(int64(modelMap["target_id"].(int)))
	if modelMap["target_name"] != nil && modelMap["target_name"].(string) != "" {
		model.TargetName = core.StringPtr(modelMap["target_name"].(string))
	}
	if modelMap["target_type"] != nil && modelMap["target_type"].(string) != "" {
		model.TargetType = core.StringPtr(modelMap["target_type"].(string))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToCascadedTargetConfiguration(modelMap map[string]interface{}) (*backuprecoveryv1.CascadedTargetConfiguration, error) {
	model := &backuprecoveryv1.CascadedTargetConfiguration{}
	model.SourceClusterID = core.Int64Ptr(int64(modelMap["source_cluster_id"].(int)))
	RemoteTargetsModel, err := ResourceIbmBackupRecoveryProtectionPolicyMapToTargetsConfiguration(modelMap["remote_targets"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.RemoteTargets = RemoteTargetsModel
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMapToRetryOptions(modelMap map[string]interface{}) (*backuprecoveryv1.RetryOptions, error) {
	model := &backuprecoveryv1.RetryOptions{}
	if modelMap["retries"] != nil {
		model.Retries = core.Int64Ptr(int64(modelMap["retries"].(int)))
	}
	if modelMap["retry_interval_mins"] != nil {
		model.RetryIntervalMins = core.Int64Ptr(int64(modelMap["retry_interval_mins"].(int)))
	}
	return model, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyBackupPolicyToMap(model *backuprecoveryv1.BackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	regularMap, err := ResourceIbmBackupRecoveryProtectionPolicyRegularBackupPolicyToMap(model.Regular)
	if err != nil {
		return modelMap, err
	}
	modelMap["regular"] = []map[string]interface{}{regularMap}
	if model.Log != nil {
		logMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogBackupPolicyToMap(model.Log)
		if err != nil {
			return modelMap, err
		}
		modelMap["log"] = []map[string]interface{}{logMap}
	}
	if model.Bmr != nil {
		bmrMap, err := ResourceIbmBackupRecoveryProtectionPolicyBmrBackupPolicyToMap(model.Bmr)
		if err != nil {
			return modelMap, err
		}
		modelMap["bmr"] = []map[string]interface{}{bmrMap}
	}
	if model.Cdp != nil {
		cdpMap, err := ResourceIbmBackupRecoveryProtectionPolicyCdpBackupPolicyToMap(model.Cdp)
		if err != nil {
			return modelMap, err
		}
		modelMap["cdp"] = []map[string]interface{}{cdpMap}
	}
	if model.StorageArraySnapshot != nil {
		storageArraySnapshotMap, err := ResourceIbmBackupRecoveryProtectionPolicyStorageArraySnapshotBackupPolicyToMap(model.StorageArraySnapshot)
		if err != nil {
			return modelMap, err
		}
		modelMap["storage_array_snapshot"] = []map[string]interface{}{storageArraySnapshotMap}
	}
	if model.RunTimeouts != nil {
		runTimeouts := []map[string]interface{}{}
		for _, runTimeoutsItem := range model.RunTimeouts {
			runTimeoutsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyRegularBackupPolicyToMap(model *backuprecoveryv1.RegularBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Incremental != nil {
		incrementalMap, err := ResourceIbmBackupRecoveryProtectionPolicyIncrementalBackupPolicyToMap(model.Incremental)
		if err != nil {
			return modelMap, err
		}
		modelMap["incremental"] = []map[string]interface{}{incrementalMap}
	}
	if model.Full != nil {
		fullMap, err := ResourceIbmBackupRecoveryProtectionPolicyFullBackupPolicyToMap(model.Full)
		if err != nil {
			return modelMap, err
		}
		modelMap["full"] = []map[string]interface{}{fullMap}
	}
	if model.FullBackups != nil {
		fullBackups := []map[string]interface{}{}
		for _, fullBackupsItem := range model.FullBackups {
			fullBackupsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyFullScheduleAndRetentionToMap(&fullBackupsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			fullBackups = append(fullBackups, fullBackupsItemMap)
		}
		modelMap["full_backups"] = fullBackups
	}
	if model.Retention != nil {
		retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
		if err != nil {
			return modelMap, err
		}
		modelMap["retention"] = []map[string]interface{}{retentionMap}
	}
	if model.PrimaryBackupTarget != nil {
		primaryBackupTargetMap, err := ResourceIbmBackupRecoveryProtectionPolicyPrimaryBackupTargetToMap(model.PrimaryBackupTarget)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_backup_target"] = []map[string]interface{}{primaryBackupTargetMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyIncrementalBackupPolicyToMap(model *backuprecoveryv1.IncrementalBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyIncrementalScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyIncrementalScheduleToMap(model *backuprecoveryv1.IncrementalSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.MinuteSchedule != nil {
		minuteScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMinuteScheduleToMap(model.MinuteSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["minute_schedule"] = []map[string]interface{}{minuteScheduleMap}
	}
	if model.HourSchedule != nil {
		hourScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyHourScheduleToMap(model.HourSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["hour_schedule"] = []map[string]interface{}{hourScheduleMap}
	}
	if model.DaySchedule != nil {
		dayScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMinuteScheduleToMap(model *backuprecoveryv1.MinuteSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["frequency"] = flex.IntValue(model.Frequency)
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyHourScheduleToMap(model *backuprecoveryv1.HourSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["frequency"] = flex.IntValue(model.Frequency)
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyDayScheduleToMap(model *backuprecoveryv1.DaySchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["frequency"] = flex.IntValue(model.Frequency)
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyWeekScheduleToMap(model *backuprecoveryv1.WeekSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["day_of_week"] = model.DayOfWeek
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyMonthScheduleToMap(model *backuprecoveryv1.MonthSchedule) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyYearScheduleToMap(model *backuprecoveryv1.YearSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["day_of_year"] = *model.DayOfYear
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyFullBackupPolicyToMap(model *backuprecoveryv1.FullBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Schedule != nil {
		scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyFullScheduleToMap(model.Schedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyFullScheduleToMap(model *backuprecoveryv1.FullSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.DaySchedule != nil {
		dayScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyFullScheduleAndRetentionToMap(model *backuprecoveryv1.FullScheduleAndRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyFullScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model *backuprecoveryv1.Retention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := ResourceIbmBackupRecoveryProtectionPolicyDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyDataLockConfigToMap(model *backuprecoveryv1.DataLockConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = *model.Mode
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.EnableWormOnExternalTarget != nil {
		modelMap["enable_worm_on_external_target"] = *model.EnableWormOnExternalTarget
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyPrimaryBackupTargetToMap(model *backuprecoveryv1.PrimaryBackupTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TargetType != nil {
		modelMap["target_type"] = *model.TargetType
	}
	if model.ArchivalTargetSettings != nil {
		archivalTargetSettingsMap, err := ResourceIbmBackupRecoveryProtectionPolicyPrimaryArchivalTargetToMap(model.ArchivalTargetSettings)
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

func ResourceIbmBackupRecoveryProtectionPolicyPrimaryArchivalTargetToMap(model *backuprecoveryv1.PrimaryArchivalTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["target_id"] = flex.IntValue(model.TargetID)
	if model.TargetName != nil {
		modelMap["target_name"] = *model.TargetName
	}
	if model.TierSettings != nil {
		tierSettingsMap, err := ResourceIbmBackupRecoveryProtectionPolicyTierLevelSettingsToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyTierLevelSettingsToMap(model *backuprecoveryv1.TierLevelSettings) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsTiering != nil {
		awsTieringMap, err := ResourceIbmBackupRecoveryProtectionPolicyAWSTiersToMap(model.AwsTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_tiering"] = []map[string]interface{}{awsTieringMap}
	}
	if model.AzureTiering != nil {
		azureTieringMap, err := ResourceIbmBackupRecoveryProtectionPolicyAzureTiersToMap(model.AzureTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_tiering"] = []map[string]interface{}{azureTieringMap}
	}
	if model.CloudPlatform != nil {
		modelMap["cloud_platform"] = *model.CloudPlatform
	}
	if model.GoogleTiering != nil {
		googleTieringMap, err := ResourceIbmBackupRecoveryProtectionPolicyGoogleTiersToMap(model.GoogleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["google_tiering"] = []map[string]interface{}{googleTieringMap}
	}
	if model.OracleTiering != nil {
		oracleTieringMap, err := ResourceIbmBackupRecoveryProtectionPolicyOracleTiersToMap(model.OracleTiering)
		if err != nil {
			return modelMap, err
		}
		modelMap["oracle_tiering"] = []map[string]interface{}{oracleTieringMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyAWSTiersToMap(model *backuprecoveryv1.AWSTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyAWSTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyAWSTierToMap(model *backuprecoveryv1.AWSTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyAzureTiersToMap(model *backuprecoveryv1.AzureTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Tiers != nil {
		tiers := []map[string]interface{}{}
		for _, tiersItem := range model.Tiers {
			tiersItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyAzureTierToMap(&tiersItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			tiers = append(tiers, tiersItemMap)
		}
		modelMap["tiers"] = tiers
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyAzureTierToMap(model *backuprecoveryv1.AzureTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyGoogleTiersToMap(model *backuprecoveryv1.GoogleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyGoogleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyGoogleTierToMap(model *backuprecoveryv1.GoogleTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyOracleTiersToMap(model *backuprecoveryv1.OracleTiers) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	tiers := []map[string]interface{}{}
	for _, tiersItem := range model.Tiers {
		tiersItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyOracleTierToMap(&tiersItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		tiers = append(tiers, tiersItemMap)
	}
	modelMap["tiers"] = tiers
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyOracleTierToMap(model *backuprecoveryv1.OracleTier) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyLogBackupPolicyToMap(model *backuprecoveryv1.LogBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyLogScheduleToMap(model *backuprecoveryv1.LogSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.MinuteSchedule != nil {
		minuteScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMinuteScheduleToMap(model.MinuteSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["minute_schedule"] = []map[string]interface{}{minuteScheduleMap}
	}
	if model.HourSchedule != nil {
		hourScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyHourScheduleToMap(model.HourSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["hour_schedule"] = []map[string]interface{}{hourScheduleMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyBmrBackupPolicyToMap(model *backuprecoveryv1.BmrBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyBmrScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyBmrScheduleToMap(model *backuprecoveryv1.BmrSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.DaySchedule != nil {
		dayScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyCdpBackupPolicyToMap(model *backuprecoveryv1.CdpBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyCdpRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyCdpRetentionToMap(model *backuprecoveryv1.CdpRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := ResourceIbmBackupRecoveryProtectionPolicyDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyStorageArraySnapshotBackupPolicyToMap(model *backuprecoveryv1.StorageArraySnapshotBackupPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyStorageArraySnapshotScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
	if err != nil {
		return modelMap, err
	}
	modelMap["retention"] = []map[string]interface{}{retentionMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyStorageArraySnapshotScheduleToMap(model *backuprecoveryv1.StorageArraySnapshotSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.MinuteSchedule != nil {
		minuteScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMinuteScheduleToMap(model.MinuteSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["minute_schedule"] = []map[string]interface{}{minuteScheduleMap}
	}
	if model.HourSchedule != nil {
		hourScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyHourScheduleToMap(model.HourSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["hour_schedule"] = []map[string]interface{}{hourScheduleMap}
	}
	if model.DaySchedule != nil {
		dayScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyDayScheduleToMap(model.DaySchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["day_schedule"] = []map[string]interface{}{dayScheduleMap}
	}
	if model.WeekSchedule != nil {
		weekScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyWeekScheduleToMap(model.WeekSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["week_schedule"] = []map[string]interface{}{weekScheduleMap}
	}
	if model.MonthSchedule != nil {
		monthScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyMonthScheduleToMap(model.MonthSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["month_schedule"] = []map[string]interface{}{monthScheduleMap}
	}
	if model.YearSchedule != nil {
		yearScheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyYearScheduleToMap(model.YearSchedule)
		if err != nil {
			return modelMap, err
		}
		modelMap["year_schedule"] = []map[string]interface{}{yearScheduleMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(model *backuprecoveryv1.CancellationTimeoutParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TimeoutMins != nil {
		modelMap["timeout_mins"] = flex.IntValue(model.TimeoutMins)
	}
	if model.BackupType != nil {
		modelMap["backup_type"] = *model.BackupType
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyBlackoutWindowToMap(model *backuprecoveryv1.BlackoutWindow) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["day"] = *model.Day
	startTimeMap, err := ResourceIbmBackupRecoveryProtectionPolicyTimeOfDayToMap(model.StartTime)
	if err != nil {
		return modelMap, err
	}
	modelMap["start_time"] = []map[string]interface{}{startTimeMap}
	endTimeMap, err := ResourceIbmBackupRecoveryProtectionPolicyTimeOfDayToMap(model.EndTime)
	if err != nil {
		return modelMap, err
	}
	modelMap["end_time"] = []map[string]interface{}{endTimeMap}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyTimeOfDayToMap(model *backuprecoveryv1.TimeOfDay) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hour"] = flex.IntValue(model.Hour)
	modelMap["minute"] = flex.IntValue(model.Minute)
	if model.TimeZone != nil {
		modelMap["time_zone"] = *model.TimeZone
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyExtendedRetentionPolicyToMap(model *backuprecoveryv1.ExtendedRetentionPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyExtendedRetentionScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
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

func ResourceIbmBackupRecoveryProtectionPolicyExtendedRetentionScheduleToMap(model *backuprecoveryv1.ExtendedRetentionSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.Frequency != nil {
		modelMap["frequency"] = flex.IntValue(model.Frequency)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyTargetsConfigurationToMap(model *backuprecoveryv1.TargetsConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ReplicationTargets != nil {
		replicationTargets := []map[string]interface{}{}
		for _, replicationTargetsItem := range model.ReplicationTargets {
			replicationTargetsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyReplicationTargetConfigurationToMap(&replicationTargetsItem) // #nosec G601
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
			archivalTargetsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyArchivalTargetConfigurationToMap(&archivalTargetsItem) // #nosec G601
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
			cloudSpinTargetsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCloudSpinTargetConfigurationToMap(&cloudSpinTargetsItem) // #nosec G601
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
			onpremDeployTargetsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyOnpremDeployTargetConfigurationToMap(&onpremDeployTargetsItem) // #nosec G601
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
			rpaasTargetsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyRpaasTargetConfigurationToMap(&rpaasTargetsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			rpaasTargets = append(rpaasTargets, rpaasTargetsItemMap)
		}
		modelMap["rpaas_targets"] = rpaasTargets
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyReplicationTargetConfigurationToMap(model *backuprecoveryv1.ReplicationTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
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
			runTimeoutsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	if model.AwsTargetConfig != nil {
		awsTargetConfigMap, err := ResourceIbmBackupRecoveryProtectionPolicyAWSTargetConfigToMap(model.AwsTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_target_config"] = []map[string]interface{}{awsTargetConfigMap}
	}
	if model.AzureTargetConfig != nil {
		azureTargetConfigMap, err := ResourceIbmBackupRecoveryProtectionPolicyAzureTargetConfigToMap(model.AzureTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["azure_target_config"] = []map[string]interface{}{azureTargetConfigMap}
	}
	modelMap["target_type"] = *model.TargetType
	if model.RemoteTargetConfig != nil {
		remoteTargetConfigMap, err := ResourceIbmBackupRecoveryProtectionPolicyRemoteTargetConfigToMap(model.RemoteTargetConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote_target_config"] = []map[string]interface{}{remoteTargetConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyTargetScheduleToMap(model *backuprecoveryv1.TargetSchedule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	if model.Frequency != nil {
		modelMap["frequency"] = flex.IntValue(model.Frequency)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyLogRetentionToMap(model *backuprecoveryv1.LogRetention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["unit"] = *model.Unit
	modelMap["duration"] = flex.IntValue(model.Duration)
	if model.DataLockConfig != nil {
		dataLockConfigMap, err := ResourceIbmBackupRecoveryProtectionPolicyDataLockConfigToMap(model.DataLockConfig)
		if err != nil {
			return modelMap, err
		}
		modelMap["data_lock_config"] = []map[string]interface{}{dataLockConfigMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyAWSTargetConfigToMap(model *backuprecoveryv1.AWSTargetConfig) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyAzureTargetConfigToMap(model *backuprecoveryv1.AzureTargetConfig) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyRemoteTargetConfigToMap(model *backuprecoveryv1.RemoteTargetConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cluster_id"] = flex.IntValue(model.ClusterID)
	if model.ClusterName != nil {
		modelMap["cluster_name"] = *model.ClusterName
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyArchivalTargetConfigurationToMap(model *backuprecoveryv1.ArchivalTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
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
			runTimeoutsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogRetentionToMap(model.LogRetention)
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
		tierSettingsMap, err := ResourceIbmBackupRecoveryProtectionPolicyTierLevelSettingsToMap(model.TierSettings)
		if err != nil {
			return modelMap, err
		}
		modelMap["tier_settings"] = []map[string]interface{}{tierSettingsMap}
	}
	if model.ExtendedRetention != nil {
		extendedRetention := []map[string]interface{}{}
		for _, extendedRetentionItem := range model.ExtendedRetention {
			extendedRetentionItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyExtendedRetentionPolicyToMap(&extendedRetentionItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			extendedRetention = append(extendedRetention, extendedRetentionItemMap)
		}
		modelMap["extended_retention"] = extendedRetention
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyCloudSpinTargetConfigurationToMap(model *backuprecoveryv1.CloudSpinTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
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
			runTimeoutsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	targetMap, err := ResourceIbmBackupRecoveryProtectionPolicyCloudSpinTargetToMap(model.Target)
	if err != nil {
		return modelMap, err
	}
	modelMap["target"] = []map[string]interface{}{targetMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyCloudSpinTargetToMap(model *backuprecoveryv1.CloudSpinTarget) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.AwsParams != nil {
		awsParamsMap, err := ResourceIbmBackupRecoveryProtectionPolicyAwsCloudSpinParamsToMap(model.AwsParams)
		if err != nil {
			return modelMap, err
		}
		modelMap["aws_params"] = []map[string]interface{}{awsParamsMap}
	}
	if model.AzureParams != nil {
		azureParamsMap, err := ResourceIbmBackupRecoveryProtectionPolicyAzureCloudSpinParamsToMap(model.AzureParams)
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

func ResourceIbmBackupRecoveryProtectionPolicyAwsCloudSpinParamsToMap(model *backuprecoveryv1.AwsCloudSpinParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CustomTagList != nil {
		customTagList := []map[string]interface{}{}
		for _, customTagListItem := range model.CustomTagList {
			customTagListItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCustomTagParamsToMap(&customTagListItem) // #nosec G601
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

func ResourceIbmBackupRecoveryProtectionPolicyCustomTagParamsToMap(model *backuprecoveryv1.CustomTagParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["key"] = *model.Key
	modelMap["value"] = *model.Value
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyAzureCloudSpinParamsToMap(model *backuprecoveryv1.AzureCloudSpinParams) (map[string]interface{}, error) {
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

func ResourceIbmBackupRecoveryProtectionPolicyOnpremDeployTargetConfigurationToMap(model *backuprecoveryv1.OnpremDeployTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
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
			runTimeoutsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogRetentionToMap(model.LogRetention)
		if err != nil {
			return modelMap, err
		}
		modelMap["log_retention"] = []map[string]interface{}{logRetentionMap}
	}
	if model.Params != nil {
		paramsMap, err := ResourceIbmBackupRecoveryProtectionPolicyOnpremDeployParamsToMap(model.Params)
		if err != nil {
			return modelMap, err
		}
		modelMap["params"] = []map[string]interface{}{paramsMap}
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyOnpremDeployParamsToMap(model *backuprecoveryv1.OnpremDeployParams) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = flex.IntValue(model.ID)
	}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyRpaasTargetConfigurationToMap(model *backuprecoveryv1.RpaasTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	scheduleMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetScheduleToMap(model.Schedule)
	if err != nil {
		return modelMap, err
	}
	modelMap["schedule"] = []map[string]interface{}{scheduleMap}
	retentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyRetentionToMap(model.Retention)
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
			runTimeoutsItemMap, err := ResourceIbmBackupRecoveryProtectionPolicyCancellationTimeoutParamsToMap(&runTimeoutsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			runTimeouts = append(runTimeouts, runTimeoutsItemMap)
		}
		modelMap["run_timeouts"] = runTimeouts
	}
	if model.LogRetention != nil {
		logRetentionMap, err := ResourceIbmBackupRecoveryProtectionPolicyLogRetentionToMap(model.LogRetention)
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

func ResourceIbmBackupRecoveryProtectionPolicyCascadedTargetConfigurationToMap(model *backuprecoveryv1.CascadedTargetConfiguration) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["source_cluster_id"] = flex.IntValue(model.SourceClusterID)
	remoteTargetsMap, err := ResourceIbmBackupRecoveryProtectionPolicyTargetsConfigurationToMap(model.RemoteTargets)
	if err != nil {
		return modelMap, err
	}
	modelMap["remote_targets"] = []map[string]interface{}{remoteTargetsMap}
	return modelMap, nil
}

func ResourceIbmBackupRecoveryProtectionPolicyRetryOptionsToMap(model *backuprecoveryv1.RetryOptions) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Retries != nil {
		modelMap["retries"] = flex.IntValue(model.Retries)
	}
	if model.RetryIntervalMins != nil {
		modelMap["retry_interval_mins"] = flex.IntValue(model.RetryIntervalMins)
	}
	return modelMap, nil
}
