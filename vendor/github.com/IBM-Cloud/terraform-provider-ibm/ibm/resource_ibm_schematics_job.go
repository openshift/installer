// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func resourceIBMSchematicsJob() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSchematicsJobCreate,
		Read:     resourceIBMSchematicsJobRead,
		Update:   resourceIBMSchematicsJobUpdate,
		Delete:   resourceIBMSchematicsJobDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"command_object": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_schematics_job", "command_object"),
				Description:  "Name of the Schematics automation resource.",
			},
			"command_object_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Job command object ID (`workspace-id, action-id or control-id`).",
			},
			"command_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_schematics_job", "command_name"),
				Description:  "Schematics job command name.",
			},
			"command_parameter": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Schematics job command parameter (`playbook-name, capsule-name or flow-name`).",
			},
			"command_options": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Command line options for the command.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"job_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Job inputs used by an action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of aliases for the variable name.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"job_env_settings": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Environment variables used by the job while performing an action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of aliases for the variable name.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User defined tags, while running the job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"location": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_schematics_job", "location"),
				Description:  "List of action locations supported by IBM Cloud Schematics service.  **Note** this does not limit the location of the resources provisioned using Schematics.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job Status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_job_status": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Action Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action name.",
									},
									"status_code": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of the jobs.",
									},
									"status_message": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action job status message to be displayed along with the `action_status_code`.",
									},
									"bastion_status_code": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of the resources.",
									},
									"bastion_status_message": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Bastion status message to be displayed along with the `bastion_status_code`.",
									},
									"targets_status_code": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of the resources.",
									},
									"targets_status_message": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Aggregated status message for all target resources, to be displayed along with the `targets_status_code`.",
									},
									"updated_at": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
					},
				},
			},
			"data": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Job data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the job.",
						},
						"action_job_data": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Action Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow name.",
									},
									"inputs": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Input variables data used by an action job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"outputs": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Output variables data from an action job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"settings": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Environment variables used by all the templates in an action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": &schema.Schema{
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": &schema.Schema{
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": &schema.Schema{
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"updated_at": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
					},
				},
			},
			"bastion": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Complete target details with the user inputs and the system generated data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target name.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target type (`cluster`, `vsi`, `icd`, `vpc`).",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Target description.",
						},
						"resource_query": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resource selection query string.",
						},
						"credential_ref": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Override credential for each resource.  Reference to credentials values, used by all the resources.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Target ID.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Targets creation time.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "E-mail address of the user who created the targets.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Targets updation time.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "E-mail address of user who updated the targets.",
						},
						"sys_lock": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "System lock status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sys_locked": &schema.Schema{
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the Workspace locked by the Schematic action ?.",
									},
									"sys_locked_by": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the user who performed the action, that lead to lock the Workspace.",
									},
									"sys_locked_at": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "When the user performed the action that lead to lock the Workspace ?.",
									},
								},
							},
						},
						"resource_ids": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Array of the resource IDs.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"job_log_summary": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Job log summary record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Workspace ID.",
						},
						"job_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of Job.",
						},
						"log_start_at": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Job log start timestamp.",
						},
						"log_analyzed_till": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Job log update timestamp.",
						},
						"elapsed_time": &schema.Schema{
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "Job log elapsed time (`log_analyzed_till - log_start_at`).",
						},
						"log_errors": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Job log errors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Error code in the Log.",
									},
									"error_msg": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Summary error message in the log.",
									},
									"error_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Number of occurrence.",
									},
								},
							},
						},
						"repo_download_job": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Repo download Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scanned_file_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of files scanned.",
									},
									"quarantined_file_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of files quarantined.",
									},
									"detected_filetype": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Detected template or data file type.",
									},
									"inputs_count": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Number of inputs detected.",
									},
									"outputs_count": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Number of outputs detected.",
									},
								},
							},
						},
						"action_job": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Flow Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of targets or hosts.",
									},
									"task_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of tasks in playbook.",
									},
									"play_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of plays in playbook.",
									},
									"recap": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Recap records.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target": &schema.Schema{
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of target or host name.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"ok": &schema.Schema{
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of OK.",
												},
												"changed": &schema.Schema{
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of changed.",
												},
												"failed": &schema.Schema{
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of failed.",
												},
												"skipped": &schema.Schema{
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of skipped.",
												},
												"unreachable": &schema.Schema{
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of unreachable.",
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
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job name, uniquely derived from the related action.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job description derived from the related action.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group name derived from the related action.",
			},
			"submitted_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job submission time.",
			},
			"submitted_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who submitted the job.",
			},
			"start_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job start time.",
			},
			"end_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job end time.",
			},
			"duration": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Duration of job execution, for example, `40 sec`.",
			},
			"targets_ini": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inventory of host and host group for the playbook in `INI` file format. For example, `\"targets_ini\": \"[webserverhost]  172.22.192.6  [dbhost]  172.22.192.5\"`. For more information, about an inventory host group syntax, see [Inventory host groups](/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).",
			},
			"log_store_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job log store URL.",
			},
			"state_store_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job state store URL.",
			},
			"results_url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job results store URL.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job status updation timestamp.",
			},
		},
	}
}

func resourceIBMSchematicsJobValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "command_object",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              "action, workspace",
		},
		ValidateSchema{
			Identifier:                 "command_name",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              "ansible_playbook_check, ansible_playbook_run, helm_install, helm_list, helm_show, opa_evaluate, terraform_init, terrform_apply, terrform_destroy, terrform_plan, terrform_refresh, terrform_show, terrform_taint, workspace_apply_flow, workspace_custom_flow, workspace_destroy_flow, workspace_init_flow, workspace_plan_flow, workspace_refresh_flow, workspace_show_flow",
		},
		ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              "eu-de, eu-gb, us-east, us-south",
		})

	resourceValidator := ResourceValidator{ResourceName: "ibm_schematics_job", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSchematicsJobCreate(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	session, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	iamRefreshToken := session.Config.IAMRefreshToken

	createJobOptions := &schematicsv1.CreateJobOptions{}
	createJobOptions.SetRefreshToken(iamRefreshToken)

	if _, ok := d.GetOk("command_object"); ok {
		createJobOptions.SetCommandObject(d.Get("command_object").(string))
	}
	if _, ok := d.GetOk("command_object_id"); ok {
		createJobOptions.SetCommandObjectID(d.Get("command_object_id").(string))
	}
	if _, ok := d.GetOk("command_name"); ok {
		createJobOptions.SetCommandName(d.Get("command_name").(string))
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		createJobOptions.SetCommandParameter(d.Get("command_parameter").(string))
	}
	if _, ok := d.GetOk("command_options"); ok {
		createJobOptions.SetCommandOptions(d.Get("command_options").([]string))
	}
	if _, ok := d.GetOk("job_inputs"); ok {
		var inputs []schematicsv1.VariableData
		for _, e := range d.Get("job_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			inputsItem := resourceIBMSchematicsJobMapToVariableData(value)
			inputs = append(inputs, inputsItem)
		}
		createJobOptions.SetInputs(inputs)
	}
	if _, ok := d.GetOk("job_env_settings"); ok {
		var settings []schematicsv1.VariableData
		for _, e := range d.Get("job_env_settings").([]interface{}) {
			value := e.(map[string]interface{})
			settingsItem := resourceIBMSchematicsJobMapToVariableData(value)
			settings = append(settings, settingsItem)
		}
		createJobOptions.SetSettings(settings)
	}
	if _, ok := d.GetOk("tags"); ok {
		createJobOptions.SetTags(expandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("location"); ok {
		createJobOptions.SetLocation(d.Get("location").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		statusAttr := d.Get("status").([]interface{})
		if len(statusAttr) > 0 {
			status := resourceIBMSchematicsJobMapToJobStatus(d.Get("status.0").(map[string]interface{}))
			createJobOptions.SetStatus(&status)
		}
	}
	if _, ok := d.GetOk("data"); ok {
		dataAttr := d.Get("data").([]interface{})
		if len(dataAttr) > 0 {
			data := resourceIBMSchematicsJobMapToJobData(d.Get("data.0").(map[string]interface{}))
			createJobOptions.SetData(&data)
		}
	}
	if _, ok := d.GetOk("bastion"); ok {
		bastionAttr := d.Get("bastion").([]interface{})
		if len(bastionAttr) > 0 {
			bastion := resourceIBMSchematicsJobMapToTargetResourceset(d.Get("bastion.0").(map[string]interface{}))
			createJobOptions.SetBastion(&bastion)
		}
	}
	if _, ok := d.GetOk("job_log_summary"); ok {
		jobLogSummaryAttr := d.Get("job_log_summary").([]interface{})
		if len(jobLogSummaryAttr) > 0 {
			logSummary := resourceIBMSchematicsJobMapToJobLogSummary(d.Get("job_log_summary.0").(map[string]interface{}))
			createJobOptions.SetLogSummary(&logSummary)
		}
	}

	job, response, err := schematicsClient.CreateJobWithContext(context.TODO(), createJobOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateJobWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*job.ID)

	return resourceIBMSchematicsJobRead(d, meta)
}

func resourceIBMSchematicsJobMapToVariableData(variableDataMap map[string]interface{}) schematicsv1.VariableData {
	variableData := schematicsv1.VariableData{}

	if variableDataMap["name"] != nil {
		variableData.Name = core.StringPtr(variableDataMap["name"].(string))
	}
	if variableDataMap["value"] != nil {
		variableData.Value = core.StringPtr(variableDataMap["value"].(string))
	}
	if variableDataMap["metadata"] != nil && len(variableDataMap["metadata"].([]interface{})) != 0 {
		variableMetaData := resourceIBMSchematicsJobMapToVariableMetadata(variableDataMap["metadata"].([]interface{})[0].(map[string]interface{}))
		variableData.Metadata = &variableMetaData
	}
	if variableDataMap["link"] != nil {
		variableData.Link = core.StringPtr(variableDataMap["link"].(string))
	}

	return variableData
}

func resourceIBMSchematicsJobMapToVariableMetadata(variableMetadataMap map[string]interface{}) schematicsv1.VariableMetadata {
	variableMetadata := schematicsv1.VariableMetadata{}

	if variableMetadataMap["type"] != nil {
		variableMetadata.Type = core.StringPtr(variableMetadataMap["type"].(string))
	}
	if variableMetadataMap["aliases"] != nil {
		aliases := []string{}
		for _, aliasesItem := range variableMetadataMap["aliases"].([]interface{}) {
			aliases = append(aliases, aliasesItem.(string))
		}
		variableMetadata.Aliases = aliases
	}
	if variableMetadataMap["description"] != nil {
		variableMetadata.Description = core.StringPtr(variableMetadataMap["description"].(string))
	}
	if variableMetadataMap["default_value"] != nil {
		variableMetadata.DefaultValue = core.StringPtr(variableMetadataMap["default_value"].(string))
	}
	if variableMetadataMap["secure"] != nil {
		variableMetadata.Secure = core.BoolPtr(variableMetadataMap["secure"].(bool))
	}
	if variableMetadataMap["immutable"] != nil {
		variableMetadata.Immutable = core.BoolPtr(variableMetadataMap["immutable"].(bool))
	}
	if variableMetadataMap["hidden"] != nil {
		variableMetadata.Hidden = core.BoolPtr(variableMetadataMap["hidden"].(bool))
	}
	if variableMetadataMap["options"] != nil {
		options := []string{}
		for _, optionsItem := range variableMetadataMap["options"].([]interface{}) {
			options = append(options, optionsItem.(string))
		}
		variableMetadata.Options = options
	}
	if variableMetadataMap["min_value"] != nil {
		variableMetadata.MinValue = core.Int64Ptr(int64(variableMetadataMap["min_value"].(int)))
	}
	if variableMetadataMap["max_value"] != nil {
		variableMetadata.MaxValue = core.Int64Ptr(int64(variableMetadataMap["max_value"].(int)))
	}
	if variableMetadataMap["min_length"] != nil {
		variableMetadata.MinLength = core.Int64Ptr(int64(variableMetadataMap["min_length"].(int)))
	}
	if variableMetadataMap["max_length"] != nil {
		variableMetadata.MaxLength = core.Int64Ptr(int64(variableMetadataMap["max_length"].(int)))
	}
	if variableMetadataMap["matches"] != nil {
		variableMetadata.Matches = core.StringPtr(variableMetadataMap["matches"].(string))
	}
	if variableMetadataMap["position"] != nil {
		variableMetadata.Position = core.Int64Ptr(int64(variableMetadataMap["position"].(int)))
	}
	if variableMetadataMap["group_by"] != nil {
		variableMetadata.GroupBy = core.StringPtr(variableMetadataMap["group_by"].(string))
	}
	if variableMetadataMap["source"] != nil {
		variableMetadata.Source = core.StringPtr(variableMetadataMap["source"].(string))
	}

	return variableMetadata
}

func resourceIBMSchematicsJobMapToJobStatus(jobStatusMap map[string]interface{}) schematicsv1.JobStatus {
	jobStatus := schematicsv1.JobStatus{}

	if jobStatusMap["action_job_status"] != nil {
		actionJobStatus := resourceIBMSchematicsJobMapToJobStatusAction(jobStatusMap["action_job_status"].([]interface{})[0].(map[string]interface{}))
		jobStatus.ActionJobStatus = &actionJobStatus
	}

	return jobStatus
}

func resourceIBMSchematicsJobMapToJobStatusAction(jobStatusActionMap map[string]interface{}) schematicsv1.JobStatusAction {
	jobStatusAction := schematicsv1.JobStatusAction{}

	if jobStatusActionMap["action_name"] != nil {
		jobStatusAction.ActionName = core.StringPtr(jobStatusActionMap["action_name"].(string))
	}
	if jobStatusActionMap["status_code"] != nil {
		jobStatusAction.StatusCode = core.StringPtr(jobStatusActionMap["status_code"].(string))
	}
	if jobStatusActionMap["status_message"] != nil {
		jobStatusAction.StatusMessage = core.StringPtr(jobStatusActionMap["status_message"].(string))
	}
	if jobStatusActionMap["bastion_status_code"] != nil {
		jobStatusAction.BastionStatusCode = core.StringPtr(jobStatusActionMap["bastion_status_code"].(string))
	}
	if jobStatusActionMap["bastion_status_message"] != nil {
		jobStatusAction.BastionStatusMessage = core.StringPtr(jobStatusActionMap["bastion_status_message"].(string))
	}
	if jobStatusActionMap["targets_status_code"] != nil {
		jobStatusAction.TargetsStatusCode = core.StringPtr(jobStatusActionMap["targets_status_code"].(string))
	}
	if jobStatusActionMap["targets_status_message"] != nil {
		jobStatusAction.TargetsStatusMessage = core.StringPtr(jobStatusActionMap["targets_status_message"].(string))
	}
	if jobStatusActionMap["updated_at"] != nil {
		updatedAt, err := strfmt.ParseDateTime(jobStatusActionMap["updated_at"].(string))
		if err != nil {
			jobStatusAction.UpdatedAt = &updatedAt
		}
	}

	return jobStatusAction
}

func resourceIBMSchematicsJobMapToJobData(jobDataMap map[string]interface{}) schematicsv1.JobData {
	jobData := schematicsv1.JobData{}

	jobData.JobType = core.StringPtr(jobDataMap["job_type"].(string))
	if jobDataMap["action_job_data"] != nil {
		actionJobData := resourceIBMSchematicsJobMapToJobDataAction(jobDataMap["action_job_data"].([]interface{})[0].(map[string]interface{}))
		jobData.ActionJobData = &actionJobData
	}

	return jobData
}

func resourceIBMSchematicsJobMapToJobDataAction(jobDataActionMap map[string]interface{}) schematicsv1.JobDataAction {
	jobDataAction := schematicsv1.JobDataAction{}

	if jobDataActionMap["action_name"] != nil {
		jobDataAction.ActionName = core.StringPtr(jobDataActionMap["action_name"].(string))
	}
	if jobDataActionMap["inputs"] != nil {
		inputs := []schematicsv1.VariableData{}
		for _, inputsItem := range jobDataActionMap["inputs"].([]interface{}) {
			inputsItemModel := resourceIBMSchematicsJobMapToVariableData(inputsItem.(map[string]interface{}))
			inputs = append(inputs, inputsItemModel)
		}
		jobDataAction.Inputs = inputs
	}
	if jobDataActionMap["outputs"] != nil {
		outputs := []schematicsv1.VariableData{}
		for _, outputsItem := range jobDataActionMap["outputs"].([]interface{}) {
			outputsItemModel := resourceIBMSchematicsJobMapToVariableData(outputsItem.(map[string]interface{}))
			outputs = append(outputs, outputsItemModel)
		}
		jobDataAction.Outputs = outputs
	}
	if jobDataActionMap["settings"] != nil {
		settings := []schematicsv1.VariableData{}
		for _, settingsItem := range jobDataActionMap["settings"].([]interface{}) {
			settingsItemModel := resourceIBMSchematicsJobMapToVariableData(settingsItem.(map[string]interface{}))
			settings = append(settings, settingsItemModel)
		}
		jobDataAction.Settings = settings
	}
	if jobDataActionMap["updated_at"] != nil {

	}

	return jobDataAction
}

func resourceIBMSchematicsJobMapToTargetResourceset(targetResourcesetMap map[string]interface{}) schematicsv1.TargetResourceset {
	targetResourceset := schematicsv1.TargetResourceset{}

	if targetResourcesetMap["name"] != nil {
		targetResourceset.Name = core.StringPtr(targetResourcesetMap["name"].(string))
	}
	if targetResourcesetMap["type"] != nil {
		targetResourceset.Type = core.StringPtr(targetResourcesetMap["type"].(string))
	}
	if targetResourcesetMap["description"] != nil {
		targetResourceset.Description = core.StringPtr(targetResourcesetMap["description"].(string))
	}
	if targetResourcesetMap["resource_query"] != nil {
		targetResourceset.ResourceQuery = core.StringPtr(targetResourcesetMap["resource_query"].(string))
	}
	if targetResourcesetMap["credential_ref"] != nil {
		targetResourceset.CredentialRef = core.StringPtr(targetResourcesetMap["credential_ref"].(string))
	}
	if targetResourcesetMap["id"] != nil {
		targetResourceset.ID = core.StringPtr(targetResourcesetMap["id"].(string))
	}
	if targetResourcesetMap["created_at"] != nil {

	}
	if targetResourcesetMap["created_by"] != nil {
		targetResourceset.CreatedBy = core.StringPtr(targetResourcesetMap["created_by"].(string))
	}
	if targetResourcesetMap["updated_at"] != nil {

	}
	if targetResourcesetMap["updated_by"] != nil {
		targetResourceset.UpdatedBy = core.StringPtr(targetResourcesetMap["updated_by"].(string))
	}
	if targetResourcesetMap["sys_lock"] != nil {
		sysLock := resourceIBMSchematicsJobMapToSystemLock(targetResourcesetMap["sys_lock"].(map[string]interface{}))
		targetResourceset.SysLock = &sysLock
	}
	if targetResourcesetMap["resource_ids"] != nil {
		resourceIds := []string{}
		for _, resourceIdsItem := range targetResourcesetMap["resource_ids"].([]interface{}) {
			resourceIds = append(resourceIds, resourceIdsItem.(string))
		}
		targetResourceset.ResourceIds = resourceIds
	}

	return targetResourceset
}

func resourceIBMSchematicsJobMapToSystemLock(systemLockMap map[string]interface{}) schematicsv1.SystemLock {
	systemLock := schematicsv1.SystemLock{}

	if systemLockMap["sys_locked"] != nil {
		systemLock.SysLocked = core.BoolPtr(systemLockMap["sys_locked"].(bool))
	}
	if systemLockMap["sys_locked_by"] != nil {
		systemLock.SysLockedBy = core.StringPtr(systemLockMap["sys_locked_by"].(string))
	}
	if systemLockMap["sys_locked_at"] != nil {

	}

	return systemLock
}

func resourceIBMSchematicsJobMapToJobLogSummary(jobLogSummaryMap map[string]interface{}) schematicsv1.JobLogSummary {
	jobLogSummary := schematicsv1.JobLogSummary{}

	if jobLogSummaryMap["job_id"] != nil {
		jobLogSummary.JobID = core.StringPtr(jobLogSummaryMap["job_id"].(string))
	}
	if jobLogSummaryMap["job_type"] != nil {
		jobLogSummary.JobType = core.StringPtr(jobLogSummaryMap["job_type"].(string))
	}
	if jobLogSummaryMap["log_start_at"] != nil {

	}
	if jobLogSummaryMap["log_analyzed_till"] != nil {

	}
	if jobLogSummaryMap["elapsed_time"] != nil {
		jobLogSummary.ElapsedTime = core.Float64Ptr(jobLogSummaryMap["elapsed_time"].(float64))
	}
	if jobLogSummaryMap["log_errors"] != nil {
		logErrors := []schematicsv1.JobLogSummaryLogErrorsItem{}
		for _, logErrorsItem := range jobLogSummaryMap["log_errors"].([]interface{}) {
			logErrorsItemModel := resourceIBMSchematicsJobMapToJobLogSummaryLogErrorsItem(logErrorsItem.(map[string]interface{}))
			logErrors = append(logErrors, logErrorsItemModel)
		}
		jobLogSummary.LogErrors = logErrors
	}
	if jobLogSummaryMap["repo_download_job"] != nil {
		repoDownloadJob := resourceIBMSchematicsJobMapToJobLogSummaryRepoDownloadJob(jobLogSummaryMap["repo_download_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.RepoDownloadJob = &repoDownloadJob
	}
	if jobLogSummaryMap["action_job"] != nil {
		actionJob := resourceIBMSchematicsJobMapToJobLogSummaryActionJob(jobLogSummaryMap["action_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.ActionJob = &actionJob
	}

	return jobLogSummary
}

func resourceIBMSchematicsJobMapToJobLogSummaryLogErrorsItem(jobLogSummaryLogErrorsItemMap map[string]interface{}) schematicsv1.JobLogSummaryLogErrorsItem {
	jobLogSummaryLogErrorsItem := schematicsv1.JobLogSummaryLogErrorsItem{}

	if jobLogSummaryLogErrorsItemMap["error_code"] != nil {
		jobLogSummaryLogErrorsItem.ErrorCode = core.StringPtr(jobLogSummaryLogErrorsItemMap["error_code"].(string))
	}
	if jobLogSummaryLogErrorsItemMap["error_msg"] != nil {
		jobLogSummaryLogErrorsItem.ErrorMsg = core.StringPtr(jobLogSummaryLogErrorsItemMap["error_msg"].(string))
	}
	if jobLogSummaryLogErrorsItemMap["error_count"] != nil {
		jobLogSummaryLogErrorsItem.ErrorCount = core.Float64Ptr(jobLogSummaryLogErrorsItemMap["error_count"].(float64))
	}

	return jobLogSummaryLogErrorsItem
}

func resourceIBMSchematicsJobMapToJobLogSummaryRepoDownloadJob(jobLogSummaryRepoDownloadJobMap map[string]interface{}) schematicsv1.JobLogSummaryRepoDownloadJob {
	jobLogSummaryRepoDownloadJob := schematicsv1.JobLogSummaryRepoDownloadJob{}

	if jobLogSummaryRepoDownloadJobMap["scanned_file_count"] != nil {
		jobLogSummaryRepoDownloadJob.ScannedFileCount = core.Float64Ptr(jobLogSummaryRepoDownloadJobMap["scanned_file_count"].(float64))
	}
	if jobLogSummaryRepoDownloadJobMap["quarantined_file_count"] != nil {
		jobLogSummaryRepoDownloadJob.QuarantinedFileCount = core.Float64Ptr(jobLogSummaryRepoDownloadJobMap["quarantined_file_count"].(float64))
	}
	if jobLogSummaryRepoDownloadJobMap["detected_filetype"] != nil {
		jobLogSummaryRepoDownloadJob.DetectedFiletype = core.StringPtr(jobLogSummaryRepoDownloadJobMap["detected_filetype"].(string))
	}
	if jobLogSummaryRepoDownloadJobMap["inputs_count"] != nil {
		jobLogSummaryRepoDownloadJob.InputsCount = core.StringPtr(jobLogSummaryRepoDownloadJobMap["inputs_count"].(string))
	}
	if jobLogSummaryRepoDownloadJobMap["outputs_count"] != nil {
		jobLogSummaryRepoDownloadJob.OutputsCount = core.StringPtr(jobLogSummaryRepoDownloadJobMap["outputs_count"].(string))
	}

	return jobLogSummaryRepoDownloadJob
}

func resourceIBMSchematicsJobMapToJobLogSummaryActionJob(jobLogSummaryActionJobMap map[string]interface{}) schematicsv1.JobLogSummaryActionJob {
	jobLogSummaryActionJob := schematicsv1.JobLogSummaryActionJob{}

	if jobLogSummaryActionJobMap["target_count"] != nil {
		jobLogSummaryActionJob.TargetCount = core.Float64Ptr(jobLogSummaryActionJobMap["target_count"].(float64))
	}
	if jobLogSummaryActionJobMap["task_count"] != nil {
		jobLogSummaryActionJob.TaskCount = core.Float64Ptr(jobLogSummaryActionJobMap["task_count"].(float64))
	}
	if jobLogSummaryActionJobMap["play_count"] != nil {
		jobLogSummaryActionJob.PlayCount = core.Float64Ptr(jobLogSummaryActionJobMap["play_count"].(float64))
	}
	if jobLogSummaryActionJobMap["recap"] != nil {
		recap := resourceIBMSchematicsJobMapToJobLogSummaryActionJobRecap(jobLogSummaryActionJobMap["recap"].([]interface{})[0].(map[string]interface{}))
		jobLogSummaryActionJob.Recap = &recap
	}

	return jobLogSummaryActionJob
}

func resourceIBMSchematicsJobMapToJobLogSummaryActionJobRecap(jobLogSummaryActionJobRecapMap map[string]interface{}) schematicsv1.JobLogSummaryActionJobRecap {
	jobLogSummaryActionJobRecap := schematicsv1.JobLogSummaryActionJobRecap{}

	if jobLogSummaryActionJobRecapMap["target"] != nil {
		target := []string{}
		for _, targetItem := range jobLogSummaryActionJobRecapMap["target"].([]interface{}) {
			target = append(target, targetItem.(string))
		}
		jobLogSummaryActionJobRecap.Target = target
	}
	if jobLogSummaryActionJobRecapMap["ok"] != nil {
		jobLogSummaryActionJobRecap.Ok = core.Float64Ptr(jobLogSummaryActionJobRecapMap["ok"].(float64))
	}
	if jobLogSummaryActionJobRecapMap["changed"] != nil {
		jobLogSummaryActionJobRecap.Changed = core.Float64Ptr(jobLogSummaryActionJobRecapMap["changed"].(float64))
	}
	if jobLogSummaryActionJobRecapMap["failed"] != nil {
		jobLogSummaryActionJobRecap.Failed = core.Float64Ptr(jobLogSummaryActionJobRecapMap["failed"].(float64))
	}
	if jobLogSummaryActionJobRecapMap["skipped"] != nil {
		jobLogSummaryActionJobRecap.Skipped = core.Float64Ptr(jobLogSummaryActionJobRecapMap["skipped"].(float64))
	}
	if jobLogSummaryActionJobRecapMap["unreachable"] != nil {
		jobLogSummaryActionJobRecap.Unreachable = core.Float64Ptr(jobLogSummaryActionJobRecapMap["unreachable"].(float64))
	}

	return jobLogSummaryActionJobRecap
}

func resourceIBMSchematicsJobRead(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	getJobOptions := &schematicsv1.GetJobOptions{}

	getJobOptions.SetJobID(d.Id())

	job, response, err := schematicsClient.GetJobWithContext(context.TODO(), getJobOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetJobWithContext failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("command_object", job.CommandObject); err != nil {
		return fmt.Errorf("Error setting command_object: %s", err)
	}
	if err = d.Set("command_object_id", job.CommandObjectID); err != nil {
		return fmt.Errorf("Error setting command_object_id: %s", err)
	}
	if err = d.Set("command_name", job.CommandName); err != nil {
		return fmt.Errorf("Error setting command_name: %s", err)
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		if err = d.Set("command_parameter", d.Get("command_parameter").(string)); err != nil {
			return fmt.Errorf("Error setting command_parameter: %s", err)
		}
	}
	if job.CommandOptions != nil {
		if err = d.Set("command_options", job.CommandOptions); err != nil {
			return fmt.Errorf("Error setting command_options: %s", err)
		}
	}
	if job.Inputs != nil {
		inputs := []map[string]interface{}{}
		for _, inputsItem := range job.Inputs {
			inputsItemMap := resourceIBMSchematicsJobVariableDataToMap(inputsItem)
			inputs = append(inputs, inputsItemMap)
		}
		if err = d.Set("job_inputs", inputs); err != nil {
			return fmt.Errorf("Error setting job_inputs: %s", err)
		}
	}
	if job.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range job.Settings {
			settingsItemMap := resourceIBMSchematicsJobVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
		}
		if err = d.Set("job_env_settings", settings); err != nil {
			return fmt.Errorf("Error setting job_env_settings: %s", err)
		}
	}
	if job.Tags != nil {
		if err = d.Set("tags", job.Tags); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}
	if err = d.Set("location", job.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if job.Status != nil {
		statusMap := resourceIBMSchematicsJobJobStatusToMap(*job.Status)
		if err = d.Set("status", []map[string]interface{}{statusMap}); err != nil {
			return fmt.Errorf("Error setting status: %s", err)
		}
	}
	if job.Data != nil {
		dataMap := resourceIBMSchematicsJobJobDataToMap(*job.Data)
		if err = d.Set("data", []map[string]interface{}{dataMap}); err != nil {
			return fmt.Errorf("Error setting data: %s", err)
		}
	}
	if job.Bastion != nil {
		bastionMap := resourceIBMSchematicsJobTargetResourcesetToMap(*job.Bastion)
		if err = d.Set("bastion", []map[string]interface{}{bastionMap}); err != nil {
			return fmt.Errorf("Error setting bastion: %s", err)
		}
	}
	if job.LogSummary != nil {
		logSummaryMap := resourceIBMSchematicsJobJobLogSummaryToMap(*job.LogSummary)
		if err = d.Set("job_log_summary", []map[string]interface{}{logSummaryMap}); err != nil {
			return fmt.Errorf("Error setting job_log_summary: %s", err)
		}
	}
	if err = d.Set("name", job.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("description", job.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("resource_group", job.ResourceGroup); err != nil {
		return fmt.Errorf("Error setting resource_group: %s", err)
	}
	if err = d.Set("submitted_at", job.SubmittedAt.String()); err != nil {
		return fmt.Errorf("Error setting submitted_at: %s", err)
	}
	if err = d.Set("submitted_by", job.SubmittedBy); err != nil {
		return fmt.Errorf("Error setting submitted_by: %s", err)
	}
	if err = d.Set("start_at", job.StartAt.String()); err != nil {
		return fmt.Errorf("Error setting start_at: %s", err)
	}
	if err = d.Set("end_at", job.EndAt.String()); err != nil {
		return fmt.Errorf("Error setting end_at: %s", err)
	}
	if err = d.Set("duration", job.Duration); err != nil {
		return fmt.Errorf("Error setting duration: %s", err)
	}
	if err = d.Set("targets_ini", job.TargetsIni); err != nil {
		return fmt.Errorf("Error setting targets_ini: %s", err)
	}
	if err = d.Set("log_store_url", job.LogStoreURL); err != nil {
		return fmt.Errorf("Error setting log_store_url: %s", err)
	}
	if err = d.Set("state_store_url", job.StateStoreURL); err != nil {
		return fmt.Errorf("Error setting state_store_url: %s", err)
	}
	if err = d.Set("results_url", job.ResultsURL); err != nil {
		return fmt.Errorf("Error setting results_url: %s", err)
	}
	if err = d.Set("updated_at", job.UpdatedAt.String()); err != nil {
		return fmt.Errorf("Error setting updated_at: %s", err)
	}

	return nil
}

func resourceIBMSchematicsJobVariableDataToMap(variableData schematicsv1.VariableData) map[string]interface{} {
	variableDataMap := map[string]interface{}{}

	variableDataMap["name"] = variableData.Name
	variableDataMap["value"] = variableData.Value
	if variableData.Metadata != nil {
		MetadataMap := resourceIBMSchematicsJobVariableMetadataToMap(*variableData.Metadata)
		variableDataMap["metadata"] = []map[string]interface{}{MetadataMap}
	}
	variableDataMap["link"] = variableData.Link

	return variableDataMap
}

func resourceIBMSchematicsJobVariableMetadataToMap(variableMetadata schematicsv1.VariableMetadata) map[string]interface{} {
	variableMetadataMap := map[string]interface{}{}

	variableMetadataMap["type"] = variableMetadata.Type
	if variableMetadata.Aliases != nil {
		variableMetadataMap["aliases"] = variableMetadata.Aliases
	}
	variableMetadataMap["description"] = variableMetadata.Description
	variableMetadataMap["default_value"] = variableMetadata.DefaultValue
	variableMetadataMap["secure"] = variableMetadata.Secure
	variableMetadataMap["immutable"] = variableMetadata.Immutable
	variableMetadataMap["hidden"] = variableMetadata.Hidden
	if variableMetadata.Options != nil {
		variableMetadataMap["options"] = variableMetadata.Options
	}
	variableMetadataMap["min_value"] = intValue(variableMetadata.MinValue)
	variableMetadataMap["max_value"] = intValue(variableMetadata.MaxValue)
	variableMetadataMap["min_length"] = intValue(variableMetadata.MinLength)
	variableMetadataMap["max_length"] = intValue(variableMetadata.MaxLength)
	variableMetadataMap["matches"] = variableMetadata.Matches
	variableMetadataMap["position"] = intValue(variableMetadata.Position)
	variableMetadataMap["group_by"] = variableMetadata.GroupBy
	variableMetadataMap["source"] = variableMetadata.Source

	return variableMetadataMap
}

func resourceIBMSchematicsJobJobStatusToMap(jobStatus schematicsv1.JobStatus) map[string]interface{} {
	jobStatusMap := map[string]interface{}{}

	if jobStatus.ActionJobStatus != nil {
		ActionJobStatusMap := resourceIBMSchematicsJobJobStatusActionToMap(*jobStatus.ActionJobStatus)
		jobStatusMap["action_job_status"] = []map[string]interface{}{ActionJobStatusMap}
	}

	return jobStatusMap
}

func resourceIBMSchematicsJobJobStatusActionToMap(jobStatusAction schematicsv1.JobStatusAction) map[string]interface{} {
	jobStatusActionMap := map[string]interface{}{}

	jobStatusActionMap["action_name"] = jobStatusAction.ActionName
	jobStatusActionMap["status_code"] = jobStatusAction.StatusCode
	jobStatusActionMap["status_message"] = jobStatusAction.StatusMessage
	jobStatusActionMap["bastion_status_code"] = jobStatusAction.BastionStatusCode
	jobStatusActionMap["bastion_status_message"] = jobStatusAction.BastionStatusMessage
	jobStatusActionMap["targets_status_code"] = jobStatusAction.TargetsStatusCode
	jobStatusActionMap["targets_status_message"] = jobStatusAction.TargetsStatusMessage
	jobStatusActionMap["updated_at"] = jobStatusAction.UpdatedAt.String()

	return jobStatusActionMap
}

func resourceIBMSchematicsJobJobDataToMap(jobData schematicsv1.JobData) map[string]interface{} {
	jobDataMap := map[string]interface{}{}

	jobDataMap["job_type"] = jobData.JobType
	if jobData.ActionJobData != nil {
		ActionJobDataMap := resourceIBMSchematicsJobJobDataActionToMap(*jobData.ActionJobData)
		jobDataMap["action_job_data"] = []map[string]interface{}{ActionJobDataMap}
	}

	return jobDataMap
}

func resourceIBMSchematicsJobJobDataActionToMap(jobDataAction schematicsv1.JobDataAction) map[string]interface{} {
	jobDataActionMap := map[string]interface{}{}

	jobDataActionMap["action_name"] = jobDataAction.ActionName
	if jobDataAction.Inputs != nil {
		inputs := []map[string]interface{}{}
		for _, inputsItem := range jobDataAction.Inputs {
			inputsItemMap := resourceIBMSchematicsJobVariableDataToMap(inputsItem)
			inputs = append(inputs, inputsItemMap)
			// TODO: handle Inputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataActionMap["inputs"] = inputs
	}
	if jobDataAction.Outputs != nil {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range jobDataAction.Outputs {
			outputsItemMap := resourceIBMSchematicsJobVariableDataToMap(outputsItem)
			outputs = append(outputs, outputsItemMap)
			// TODO: handle Outputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataActionMap["outputs"] = outputs
	}
	if jobDataAction.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range jobDataAction.Settings {
			settingsItemMap := resourceIBMSchematicsJobVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
			// TODO: handle Settings of type TypeList -- list of non-primitive, not model items
		}
		jobDataActionMap["settings"] = settings
	}
	jobDataActionMap["updated_at"] = jobDataAction.UpdatedAt.String()

	return jobDataActionMap
}

func resourceIBMSchematicsJobTargetResourcesetToMap(targetResourceset schematicsv1.TargetResourceset) map[string]interface{} {
	targetResourcesetMap := map[string]interface{}{}

	targetResourcesetMap["name"] = targetResourceset.Name
	targetResourcesetMap["type"] = targetResourceset.Type
	targetResourcesetMap["description"] = targetResourceset.Description
	targetResourcesetMap["resource_query"] = targetResourceset.ResourceQuery
	targetResourcesetMap["credential_ref"] = targetResourceset.CredentialRef
	targetResourcesetMap["id"] = targetResourceset.ID
	targetResourcesetMap["created_at"] = targetResourceset.CreatedAt.String()
	targetResourcesetMap["created_by"] = targetResourceset.CreatedBy
	targetResourcesetMap["updated_at"] = targetResourceset.UpdatedAt.String()
	targetResourcesetMap["updated_by"] = targetResourceset.UpdatedBy
	if targetResourceset.SysLock != nil {
		SysLockMap := resourceIBMSchematicsJobSystemLockToMap(*targetResourceset.SysLock)
		targetResourcesetMap["sys_lock"] = []map[string]interface{}{SysLockMap}
	}
	if targetResourceset.ResourceIds != nil {
		targetResourcesetMap["resource_ids"] = targetResourceset.ResourceIds
	}

	return targetResourcesetMap
}

func resourceIBMSchematicsJobSystemLockToMap(systemLock schematicsv1.SystemLock) map[string]interface{} {
	systemLockMap := map[string]interface{}{}

	systemLockMap["sys_locked"] = systemLock.SysLocked
	systemLockMap["sys_locked_by"] = systemLock.SysLockedBy
	systemLockMap["sys_locked_at"] = systemLock.SysLockedAt.String()

	return systemLockMap
}

func resourceIBMSchematicsJobJobLogSummaryToMap(jobLogSummary schematicsv1.JobLogSummary) map[string]interface{} {
	jobLogSummaryMap := map[string]interface{}{}

	jobLogSummaryMap["job_id"] = jobLogSummary.JobID
	jobLogSummaryMap["job_type"] = jobLogSummary.JobType
	jobLogSummaryMap["log_start_at"] = jobLogSummary.LogStartAt.String()
	jobLogSummaryMap["log_analyzed_till"] = jobLogSummary.LogAnalyzedTill.String()
	jobLogSummaryMap["elapsed_time"] = jobLogSummary.ElapsedTime
	if jobLogSummary.LogErrors != nil {
		logErrors := []map[string]interface{}{}
		for _, logErrorsItem := range jobLogSummary.LogErrors {
			logErrorsItemMap := resourceIBMSchematicsJobJobLogSummaryLogErrorsItemToMap(logErrorsItem)
			logErrors = append(logErrors, logErrorsItemMap)
			// TODO: handle LogErrors of type TypeList -- list of non-primitive, not model items
		}
		jobLogSummaryMap["log_errors"] = logErrors
	}
	if jobLogSummary.RepoDownloadJob != nil {
		RepoDownloadJobMap := resourceIBMSchematicsJobJobLogSummaryRepoDownloadJobToMap(*jobLogSummary.RepoDownloadJob)
		jobLogSummaryMap["repo_download_job"] = []map[string]interface{}{RepoDownloadJobMap}
	}
	if jobLogSummary.ActionJob != nil {
		ActionJobMap := resourceIBMSchematicsJobJobLogSummaryActionJobToMap(*jobLogSummary.ActionJob)
		jobLogSummaryMap["action_job"] = []map[string]interface{}{ActionJobMap}
	}

	return jobLogSummaryMap
}

func resourceIBMSchematicsJobJobLogSummaryLogErrorsItemToMap(jobLogSummaryLogErrorsItem schematicsv1.JobLogSummaryLogErrorsItem) map[string]interface{} {
	jobLogSummaryLogErrorsItemMap := map[string]interface{}{}

	jobLogSummaryLogErrorsItemMap["error_code"] = jobLogSummaryLogErrorsItem.ErrorCode
	jobLogSummaryLogErrorsItemMap["error_msg"] = jobLogSummaryLogErrorsItem.ErrorMsg
	jobLogSummaryLogErrorsItemMap["error_count"] = jobLogSummaryLogErrorsItem.ErrorCount

	return jobLogSummaryLogErrorsItemMap
}

func resourceIBMSchematicsJobJobLogSummaryRepoDownloadJobToMap(jobLogSummaryRepoDownloadJob schematicsv1.JobLogSummaryRepoDownloadJob) map[string]interface{} {
	jobLogSummaryRepoDownloadJobMap := map[string]interface{}{}

	jobLogSummaryRepoDownloadJobMap["scanned_file_count"] = jobLogSummaryRepoDownloadJob.ScannedFileCount
	jobLogSummaryRepoDownloadJobMap["quarantined_file_count"] = jobLogSummaryRepoDownloadJob.QuarantinedFileCount
	jobLogSummaryRepoDownloadJobMap["detected_filetype"] = jobLogSummaryRepoDownloadJob.DetectedFiletype
	jobLogSummaryRepoDownloadJobMap["inputs_count"] = jobLogSummaryRepoDownloadJob.InputsCount
	jobLogSummaryRepoDownloadJobMap["outputs_count"] = jobLogSummaryRepoDownloadJob.OutputsCount

	return jobLogSummaryRepoDownloadJobMap
}

func resourceIBMSchematicsJobJobLogSummaryActionJobToMap(jobLogSummaryActionJob schematicsv1.JobLogSummaryActionJob) map[string]interface{} {
	jobLogSummaryActionJobMap := map[string]interface{}{}

	jobLogSummaryActionJobMap["target_count"] = jobLogSummaryActionJob.TargetCount
	jobLogSummaryActionJobMap["task_count"] = jobLogSummaryActionJob.TaskCount
	jobLogSummaryActionJobMap["play_count"] = jobLogSummaryActionJob.PlayCount
	if jobLogSummaryActionJob.Recap != nil {
		RecapMap := resourceIBMSchematicsJobJobLogSummaryActionJobRecapToMap(*jobLogSummaryActionJob.Recap)
		jobLogSummaryActionJobMap["recap"] = []map[string]interface{}{RecapMap}
	}

	return jobLogSummaryActionJobMap
}

func resourceIBMSchematicsJobJobLogSummaryActionJobRecapToMap(jobLogSummaryActionJobRecap schematicsv1.JobLogSummaryActionJobRecap) map[string]interface{} {
	jobLogSummaryActionJobRecapMap := map[string]interface{}{}

	if jobLogSummaryActionJobRecap.Target != nil {
		jobLogSummaryActionJobRecapMap["target"] = jobLogSummaryActionJobRecap.Target
	}
	jobLogSummaryActionJobRecapMap["ok"] = jobLogSummaryActionJobRecap.Ok
	jobLogSummaryActionJobRecapMap["changed"] = jobLogSummaryActionJobRecap.Changed
	jobLogSummaryActionJobRecapMap["failed"] = jobLogSummaryActionJobRecap.Failed
	jobLogSummaryActionJobRecapMap["skipped"] = jobLogSummaryActionJobRecap.Skipped
	jobLogSummaryActionJobRecapMap["unreachable"] = jobLogSummaryActionJobRecap.Unreachable

	return jobLogSummaryActionJobRecapMap
}

func resourceIBMSchematicsJobUpdate(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	session, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	iamRefreshToken := session.Config.IAMRefreshToken

	replaceJobOptions := &schematicsv1.ReplaceJobOptions{}

	replaceJobOptions.SetJobID(d.Id())
	replaceJobOptions.SetRefreshToken(iamRefreshToken)

	if _, ok := d.GetOk("command_object"); ok {
		replaceJobOptions.SetCommandObject(d.Get("command_object").(string))
	}
	if _, ok := d.GetOk("command_object_id"); ok {
		replaceJobOptions.SetCommandObjectID(d.Get("command_object_id").(string))
	}
	if _, ok := d.GetOk("command_name"); ok {
		replaceJobOptions.SetCommandName(d.Get("command_name").(string))
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		replaceJobOptions.SetCommandParameter(d.Get("command_parameter").(string))
	}
	if _, ok := d.GetOk("command_options"); ok {
		replaceJobOptions.SetCommandOptions(expandStringList(d.Get("command_options").([]interface{})))
	}
	if _, ok := d.GetOk("job_inputs"); ok {
		var inputs []schematicsv1.VariableData
		for _, e := range d.Get("job_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			inputsItem := resourceIBMSchematicsJobMapToVariableData(value)
			inputs = append(inputs, inputsItem)
		}
		replaceJobOptions.SetInputs(inputs)
	}
	if _, ok := d.GetOk("job_env_settings"); ok {
		var settings []schematicsv1.VariableData
		for _, e := range d.Get("job_env_settings").([]interface{}) {
			value := e.(map[string]interface{})
			settingsItem := resourceIBMSchematicsJobMapToVariableData(value)
			settings = append(settings, settingsItem)
		}
		replaceJobOptions.SetSettings(settings)
	}
	if _, ok := d.GetOk("tags"); ok {
		replaceJobOptions.SetTags(expandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("location"); ok {
		replaceJobOptions.SetLocation(d.Get("location").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		statusAttr := d.Get("status").([]interface{})
		if len(statusAttr) > 0 {
			status := resourceIBMSchematicsJobMapToJobStatus(d.Get("status.0").(map[string]interface{}))
			replaceJobOptions.SetStatus(&status)
		}
	}
	if _, ok := d.GetOk("data"); ok {
		dataAttr := d.Get("data").([]interface{})
		if len(dataAttr) > 0 {
			data := resourceIBMSchematicsJobMapToJobData(d.Get("data.0").(map[string]interface{}))
			replaceJobOptions.SetData(&data)
		}
	}
	if _, ok := d.GetOk("bastion"); ok {
		bastionAttr := d.Get("bastion").([]interface{})
		if len(bastionAttr) > 0 {
			bastion := resourceIBMSchematicsJobMapToTargetResourceset(d.Get("bastion.0").(map[string]interface{}))
			replaceJobOptions.SetBastion(&bastion)
		}
	}
	if _, ok := d.GetOk("job_log_summary"); ok {
		jobLogSummaryAttr := d.Get("job_log_summary").([]interface{})
		if len(jobLogSummaryAttr) > 0 {
			logSummary := resourceIBMSchematicsJobMapToJobLogSummary(d.Get("job_log_summary.0").(map[string]interface{}))
			replaceJobOptions.SetLogSummary(&logSummary)
		}
	}

	_, response, err := schematicsClient.ReplaceJobWithContext(context.TODO(), replaceJobOptions)
	if err != nil {
		log.Printf("[DEBUG] ReplaceJobWithContext failed %s\n%s", err, response)
		return err
	}

	return resourceIBMSchematicsJobRead(d, meta)
}

func resourceIBMSchematicsJobDelete(d *schema.ResourceData, meta interface{}) error {

	session, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	deleteJobOptions := &schematicsv1.DeleteJobOptions{}

	iamRefreshToken := session.Config.IAMRefreshToken
	deleteJobOptions.SetRefreshToken(iamRefreshToken)

	deleteJobOptions.SetJobID(d.Id())

	response, err := schematicsClient.DeleteJobWithContext(context.TODO(), deleteJobOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteJobWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}
