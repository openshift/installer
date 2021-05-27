// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func dataSourceIBMSchematicsJob() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIBMSchematicsJobRead,

		Schema: map[string]*schema.Schema{
			"job_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Use GET jobs API to look up the Job IDs in your IBM Cloud account.",
			},
			"command_object": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Schematics automation resource.",
			},
			"command_object_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job command object ID (`workspace-id, action-id or control-id`).",
			},
			"command_name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schematics job command name.",
			},
			"command_options": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Command line options for the command.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"job_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job inputs used by an action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"job_env_settings": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Environment variables used by the job while performing an action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User defined tags, while running the job.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job ID.",
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
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "List of action locations supported by IBM Cloud Schematics service.  **Note** this does not limit the location of the resources provisioned using Schematics.",
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
			"status": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job Status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_job_status": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Action Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action name.",
									},
									"status_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the jobs.",
									},
									"status_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action job status message to be displayed along with the `action_status_code`.",
									},
									"bastion_status_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the resources.",
									},
									"bastion_status_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bastion status message to be displayed along with the `bastion_status_code`.",
									},
									"targets_status_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of the resources.",
									},
									"targets_status_message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Aggregated status message for all target resources, to be displayed along with the `targets_status_code`.",
									},
									"updated_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
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
				Computed:    true,
				Description: "Job data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the job.",
						},
						"action_job_data": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Action Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow name.",
									},
									"inputs": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Input variables data used by an action job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"outputs": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Output variables data from an action job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"settings": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Environment variables used by all the templates in an action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": &schema.Schema{
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": &schema.Schema{
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": &schema.Schema{
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"updated_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
					},
				},
			},
			"targets_ini": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Inventory of host and host group for the playbook in `INI` file format. For example, `\"targets_ini\": \"[webserverhost]  172.22.192.6  [dbhost]  172.22.192.5\"`. For more information, about an inventory host group syntax, see [Inventory host groups](/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).",
			},
			"bastion": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Complete target details with the user inputs and the system generated data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target name.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target type (`cluster`, `vsi`, `icd`, `vpc`).",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target description.",
						},
						"resource_query": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource selection query string.",
						},
						"credential_ref": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Override credential for each resource.  Reference to credentials values, used by all the resources.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target ID.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Targets creation time.",
						},
						"created_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "E-mail address of the user who created the targets.",
						},
						"updated_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Targets updation time.",
						},
						"updated_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "E-mail address of user who updated the targets.",
						},
						"sys_lock": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System lock status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sys_locked": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the Workspace locked by the Schematic action ?.",
									},
									"sys_locked_by": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Name of the user who performed the action, that lead to lock the Workspace.",
									},
									"sys_locked_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "When the user performed the action that lead to lock the Workspace ?.",
									},
								},
							},
						},
						"resource_ids": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Array of the resource IDs.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"job_log_summary": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job log summary record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace ID.",
						},
						"job_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Job.",
						},
						"log_start_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job log start timestamp.",
						},
						"log_analyzed_till": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job log update timestamp.",
						},
						"elapsed_time": &schema.Schema{
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Job log elapsed time (`log_analyzed_till - log_start_at`).",
						},
						"log_errors": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Job log errors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error code in the Log.",
									},
									"error_msg": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Summary error message in the log.",
									},
									"error_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of occurrence.",
									},
								},
							},
						},
						"repo_download_job": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Repo download Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scanned_file_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of files scanned.",
									},
									"quarantined_file_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of files quarantined.",
									},
									"detected_filetype": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detected template or data file type.",
									},
									"inputs_count": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of inputs detected.",
									},
									"outputs_count": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of outputs detected.",
									},
								},
							},
						},
						"action_job": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flow Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of targets or hosts.",
									},
									"task_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of tasks in playbook.",
									},
									"play_count": &schema.Schema{
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of plays in playbook.",
									},
									"recap": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Recap records.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of target or host name.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ok": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of OK.",
												},
												"changed": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of changed.",
												},
												"failed": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of failed.",
												},
												"skipped": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of skipped.",
												},
												"unreachable": &schema.Schema{
													Type:        schema.TypeFloat,
													Computed:    true,
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

func dataSourceIBMSchematicsJobRead(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	getJobOptions := &schematicsv1.GetJobOptions{}

	getJobOptions.SetJobID(d.Get("job_id").(string))

	job, response, err := schematicsClient.GetJobWithContext(context.TODO(), getJobOptions)
	if err != nil {
		log.Printf("[DEBUG] GetJobWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*job.ID)
	if err = d.Set("command_object", job.CommandObject); err != nil {
		return fmt.Errorf("Error setting command_object: %s", err)
	}
	if err = d.Set("command_object_id", job.CommandObjectID); err != nil {
		return fmt.Errorf("Error setting command_object_id: %s", err)
	}
	if err = d.Set("command_name", job.CommandName); err != nil {
		return fmt.Errorf("Error setting command_name: %s", err)
	}
	if err = d.Set("command_options", job.CommandOptions); err != nil {
		return fmt.Errorf("Error setting command_options: %s", err)
	}

	if job.Inputs != nil {
		err = d.Set("job_inputs", dataSourceJobFlattenInputs(job.Inputs))
		if err != nil {
			return fmt.Errorf("Error setting job_inputs %s", err)
		}
	}

	if job.Settings != nil {
		err = d.Set("job_env_settings", dataSourceJobFlattenSettings(job.Settings))
		if err != nil {
			return fmt.Errorf("Error setting job_env_settings %s", err)
		}
	}
	if err = d.Set("tags", job.Tags); err != nil {
		return fmt.Errorf("Error setting tags: %s", err)
	}
	if err = d.Set("id", job.ID); err != nil {
		return fmt.Errorf("Error setting id: %s", err)
	}
	if err = d.Set("name", job.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("description", job.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("location", job.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err = d.Set("resource_group", job.ResourceGroup); err != nil {
		return fmt.Errorf("Error setting resource_group: %s", err)
	}
	if job.SubmittedAt != nil {
		if err = d.Set("submitted_at", job.SubmittedAt.String()); err != nil {
			return fmt.Errorf("Error setting submitted_at: %s", err)
		}
	}
	if err = d.Set("submitted_by", job.SubmittedBy); err != nil {
		return fmt.Errorf("Error setting submitted_by: %s", err)
	}
	if job.StartAt != nil {
		if err = d.Set("start_at", job.StartAt.String()); err != nil {
			return fmt.Errorf("Error setting start_at: %s", err)
		}
	}
	if job.EndAt != nil {
		if err = d.Set("end_at", job.EndAt.String()); err != nil {
			return fmt.Errorf("Error setting end_at: %s", err)
		}
	}
	if err = d.Set("duration", job.Duration); err != nil {
		return fmt.Errorf("Error setting duration: %s", err)
	}

	if job.Status != nil {
		err = d.Set("status", dataSourceJobFlattenStatus(*job.Status))
		if err != nil {
			return fmt.Errorf("Error setting status %s", err)
		}
	}

	if job.Data != nil {
		err = d.Set("data", dataSourceJobFlattenData(*job.Data))
		if err != nil {
			return fmt.Errorf("Error setting data %s", err)
		}
	}
	if err = d.Set("targets_ini", job.TargetsIni); err != nil {
		return fmt.Errorf("Error setting targets_ini: %s", err)
	}

	if job.Bastion != nil {
		err = d.Set("bastion", dataSourceJobFlattenBastion(*job.Bastion))
		if err != nil {
			return fmt.Errorf("Error setting bastion %s", err)
		}
	}

	if job.LogSummary != nil {
		err = d.Set("job_log_summary", dataSourceJobFlattenLogSummary(*job.LogSummary))
		if err != nil {
			return fmt.Errorf("Error setting job_log_summary %s", err)
		}
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
	if job.UpdatedAt != nil {
		if err = d.Set("updated_at", job.UpdatedAt.String()); err != nil {
			return fmt.Errorf("Error setting updated_at: %s", err)
		}
	}

	return nil
}

func dataSourceJobFlattenInputs(result []schematicsv1.VariableData) (inputs []map[string]interface{}) {
	for _, inputsItem := range result {
		inputs = append(inputs, dataSourceJobInputsToMap(inputsItem))
	}

	return inputs
}

func dataSourceJobInputsToMap(inputsItem schematicsv1.VariableData) (inputsMap map[string]interface{}) {
	inputsMap = map[string]interface{}{}

	if inputsItem.Name != nil {
		inputsMap["name"] = inputsItem.Name
	}
	if inputsItem.Value != nil {
		inputsMap["value"] = inputsItem.Value
	}
	if inputsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceJobInputsMetadataToMap(*inputsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		inputsMap["metadata"] = metadataList
	}
	if inputsItem.Link != nil {
		inputsMap["link"] = inputsItem.Link
	}

	return inputsMap
}

func dataSourceJobInputsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceJobOutputsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceJobFlattenSettings(result []schematicsv1.VariableData) (settings []map[string]interface{}) {
	for _, settingsItem := range result {
		settings = append(settings, dataSourceJobSettingsToMap(settingsItem))
	}

	return settings
}

func dataSourceJobSettingsToMap(settingsItem schematicsv1.VariableData) (settingsMap map[string]interface{}) {
	settingsMap = map[string]interface{}{}

	if settingsItem.Name != nil {
		settingsMap["name"] = settingsItem.Name
	}
	if settingsItem.Value != nil {
		settingsMap["value"] = settingsItem.Value
	}
	if settingsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceJobSettingsMetadataToMap(*settingsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		settingsMap["metadata"] = metadataList
	}
	if settingsItem.Link != nil {
		settingsMap["link"] = settingsItem.Link
	}

	return settingsMap
}

func dataSourceJobSettingsMetadataToMap(metadataItem schematicsv1.VariableMetadata) (metadataMap map[string]interface{}) {
	metadataMap = map[string]interface{}{}

	if metadataItem.Type != nil {
		metadataMap["type"] = metadataItem.Type
	}
	if metadataItem.Aliases != nil {
		metadataMap["aliases"] = metadataItem.Aliases
	}
	if metadataItem.Description != nil {
		metadataMap["description"] = metadataItem.Description
	}
	if metadataItem.DefaultValue != nil {
		metadataMap["default_value"] = metadataItem.DefaultValue
	}
	if metadataItem.Secure != nil {
		metadataMap["secure"] = metadataItem.Secure
	}
	if metadataItem.Immutable != nil {
		metadataMap["immutable"] = metadataItem.Immutable
	}
	if metadataItem.Hidden != nil {
		metadataMap["hidden"] = metadataItem.Hidden
	}
	if metadataItem.Options != nil {
		metadataMap["options"] = metadataItem.Options
	}
	if metadataItem.MinValue != nil {
		metadataMap["min_value"] = metadataItem.MinValue
	}
	if metadataItem.MaxValue != nil {
		metadataMap["max_value"] = metadataItem.MaxValue
	}
	if metadataItem.MinLength != nil {
		metadataMap["min_length"] = metadataItem.MinLength
	}
	if metadataItem.MaxLength != nil {
		metadataMap["max_length"] = metadataItem.MaxLength
	}
	if metadataItem.Matches != nil {
		metadataMap["matches"] = metadataItem.Matches
	}
	if metadataItem.Position != nil {
		metadataMap["position"] = metadataItem.Position
	}
	if metadataItem.GroupBy != nil {
		metadataMap["group_by"] = metadataItem.GroupBy
	}
	if metadataItem.Source != nil {
		metadataMap["source"] = metadataItem.Source
	}

	return metadataMap
}

func dataSourceJobFlattenStatus(result schematicsv1.JobStatus) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceJobStatusToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceJobStatusToMap(statusItem schematicsv1.JobStatus) (statusMap map[string]interface{}) {
	statusMap = map[string]interface{}{}

	if statusItem.ActionJobStatus != nil {
		actionJobStatusList := []map[string]interface{}{}
		actionJobStatusMap := dataSourceJobStatusActionJobStatusToMap(*statusItem.ActionJobStatus)
		actionJobStatusList = append(actionJobStatusList, actionJobStatusMap)
		statusMap["action_job_status"] = actionJobStatusList
	}

	return statusMap
}

func dataSourceJobStatusActionJobStatusToMap(actionJobStatusItem schematicsv1.JobStatusAction) (actionJobStatusMap map[string]interface{}) {
	actionJobStatusMap = map[string]interface{}{}

	if actionJobStatusItem.ActionName != nil {
		actionJobStatusMap["action_name"] = actionJobStatusItem.ActionName
	}
	if actionJobStatusItem.StatusCode != nil {
		actionJobStatusMap["status_code"] = actionJobStatusItem.StatusCode
	}
	if actionJobStatusItem.StatusMessage != nil {
		actionJobStatusMap["status_message"] = actionJobStatusItem.StatusMessage
	}
	if actionJobStatusItem.BastionStatusCode != nil {
		actionJobStatusMap["bastion_status_code"] = actionJobStatusItem.BastionStatusCode
	}
	if actionJobStatusItem.BastionStatusMessage != nil {
		actionJobStatusMap["bastion_status_message"] = actionJobStatusItem.BastionStatusMessage
	}
	if actionJobStatusItem.TargetsStatusCode != nil {
		actionJobStatusMap["targets_status_code"] = actionJobStatusItem.TargetsStatusCode
	}
	if actionJobStatusItem.TargetsStatusMessage != nil {
		actionJobStatusMap["targets_status_message"] = actionJobStatusItem.TargetsStatusMessage
	}
	if actionJobStatusItem.UpdatedAt != nil {
		actionJobStatusMap["updated_at"] = actionJobStatusItem.UpdatedAt.String()
	}

	return actionJobStatusMap
}

func dataSourceJobFlattenData(result schematicsv1.JobData) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceJobDataToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceJobDataToMap(dataItem schematicsv1.JobData) (dataMap map[string]interface{}) {
	dataMap = map[string]interface{}{}

	if dataItem.JobType != nil {
		dataMap["job_type"] = dataItem.JobType
	}
	if dataItem.ActionJobData != nil {
		actionJobDataList := []map[string]interface{}{}
		actionJobDataMap := dataSourceJobDataActionJobDataToMap(*dataItem.ActionJobData)
		actionJobDataList = append(actionJobDataList, actionJobDataMap)
		dataMap["action_job_data"] = actionJobDataList
	}

	return dataMap
}

func dataSourceJobDataActionJobDataToMap(actionJobDataItem schematicsv1.JobDataAction) (actionJobDataMap map[string]interface{}) {
	actionJobDataMap = map[string]interface{}{}

	if actionJobDataItem.ActionName != nil {
		actionJobDataMap["action_name"] = actionJobDataItem.ActionName
	}
	if actionJobDataItem.Inputs != nil {
		inputsList := []map[string]interface{}{}
		for _, inputsItem := range actionJobDataItem.Inputs {
			inputsList = append(inputsList, dataSourceJobActionJobDataInputsToMap(inputsItem))
		}
		actionJobDataMap["inputs"] = inputsList
	}
	if actionJobDataItem.Outputs != nil {
		outputsList := []map[string]interface{}{}
		for _, outputsItem := range actionJobDataItem.Outputs {
			outputsList = append(outputsList, dataSourceJobActionJobDataOutputsToMap(outputsItem))
		}
		actionJobDataMap["outputs"] = outputsList
	}
	if actionJobDataItem.Settings != nil {
		settingsList := []map[string]interface{}{}
		for _, settingsItem := range actionJobDataItem.Settings {
			settingsList = append(settingsList, dataSourceJobActionJobDataSettingsToMap(settingsItem))
		}
		actionJobDataMap["settings"] = settingsList
	}
	if actionJobDataItem.UpdatedAt != nil {
		actionJobDataMap["updated_at"] = actionJobDataItem.UpdatedAt.String()
	}

	return actionJobDataMap
}

func dataSourceJobActionJobDataInputsToMap(inputsItem schematicsv1.VariableData) (inputsMap map[string]interface{}) {
	inputsMap = map[string]interface{}{}

	if inputsItem.Name != nil {
		inputsMap["name"] = inputsItem.Name
	}
	if inputsItem.Value != nil {
		inputsMap["value"] = inputsItem.Value
	}
	if inputsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceJobInputsMetadataToMap(*inputsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		inputsMap["metadata"] = metadataList
	}
	if inputsItem.Link != nil {
		inputsMap["link"] = inputsItem.Link
	}

	return inputsMap
}

func dataSourceJobActionJobDataOutputsToMap(outputsItem schematicsv1.VariableData) (outputsMap map[string]interface{}) {
	outputsMap = map[string]interface{}{}

	if outputsItem.Name != nil {
		outputsMap["name"] = outputsItem.Name
	}
	if outputsItem.Value != nil {
		outputsMap["value"] = outputsItem.Value
	}
	if outputsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceJobOutputsMetadataToMap(*outputsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		outputsMap["metadata"] = metadataList
	}
	if outputsItem.Link != nil {
		outputsMap["link"] = outputsItem.Link
	}

	return outputsMap
}

func dataSourceJobActionJobDataSettingsToMap(settingsItem schematicsv1.VariableData) (settingsMap map[string]interface{}) {
	settingsMap = map[string]interface{}{}

	if settingsItem.Name != nil {
		settingsMap["name"] = settingsItem.Name
	}
	if settingsItem.Value != nil {
		settingsMap["value"] = settingsItem.Value
	}
	if settingsItem.Metadata != nil {
		metadataList := []map[string]interface{}{}
		metadataMap := dataSourceJobSettingsMetadataToMap(*settingsItem.Metadata)
		metadataList = append(metadataList, metadataMap)
		settingsMap["metadata"] = metadataList
	}
	if settingsItem.Link != nil {
		settingsMap["link"] = settingsItem.Link
	}

	return settingsMap
}

func dataSourceJobFlattenBastion(result schematicsv1.TargetResourceset) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceJobBastionToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceJobBastionToMap(bastionItem schematicsv1.TargetResourceset) (bastionMap map[string]interface{}) {
	bastionMap = map[string]interface{}{}

	if bastionItem.Name != nil {
		bastionMap["name"] = bastionItem.Name
	}
	if bastionItem.Type != nil {
		bastionMap["type"] = bastionItem.Type
	}
	if bastionItem.Description != nil {
		bastionMap["description"] = bastionItem.Description
	}
	if bastionItem.ResourceQuery != nil {
		bastionMap["resource_query"] = bastionItem.ResourceQuery
	}
	if bastionItem.CredentialRef != nil {
		bastionMap["credential_ref"] = bastionItem.CredentialRef
	}
	if bastionItem.ID != nil {
		bastionMap["id"] = bastionItem.ID
	}
	if bastionItem.CreatedAt != nil {
		bastionMap["created_at"] = bastionItem.CreatedAt.String()
	}
	if bastionItem.CreatedBy != nil {
		bastionMap["created_by"] = bastionItem.CreatedBy
	}
	if bastionItem.UpdatedAt != nil {
		bastionMap["updated_at"] = bastionItem.UpdatedAt.String()
	}
	if bastionItem.UpdatedBy != nil {
		bastionMap["updated_by"] = bastionItem.UpdatedBy
	}
	if bastionItem.SysLock != nil {
		sysLockList := []map[string]interface{}{}
		sysLockMap := dataSourceJobBastionSysLockToMap(*bastionItem.SysLock)
		sysLockList = append(sysLockList, sysLockMap)
		bastionMap["sys_lock"] = sysLockList
	}
	if bastionItem.ResourceIds != nil {
		bastionMap["resource_ids"] = bastionItem.ResourceIds
	}

	return bastionMap
}

func dataSourceJobBastionSysLockToMap(sysLockItem schematicsv1.SystemLock) (sysLockMap map[string]interface{}) {
	sysLockMap = map[string]interface{}{}

	if sysLockItem.SysLocked != nil {
		sysLockMap["sys_locked"] = sysLockItem.SysLocked
	}
	if sysLockItem.SysLockedBy != nil {
		sysLockMap["sys_locked_by"] = sysLockItem.SysLockedBy
	}
	if sysLockItem.SysLockedAt != nil {
		sysLockMap["sys_locked_at"] = sysLockItem.SysLockedAt.String()
	}

	return sysLockMap
}

func dataSourceJobFlattenLogSummary(result schematicsv1.JobLogSummary) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceJobLogSummaryToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceJobLogSummaryToMap(logSummaryItem schematicsv1.JobLogSummary) (logSummaryMap map[string]interface{}) {
	logSummaryMap = map[string]interface{}{}

	if logSummaryItem.JobID != nil {
		logSummaryMap["job_id"] = logSummaryItem.JobID
	}
	if logSummaryItem.JobType != nil {
		logSummaryMap["job_type"] = logSummaryItem.JobType
	}
	if logSummaryItem.LogStartAt != nil {
		logSummaryMap["log_start_at"] = logSummaryItem.LogStartAt.String()
	}
	if logSummaryItem.LogAnalyzedTill != nil {
		logSummaryMap["log_analyzed_till"] = logSummaryItem.LogAnalyzedTill.String()
	}
	if logSummaryItem.ElapsedTime != nil {
		logSummaryMap["elapsed_time"] = logSummaryItem.ElapsedTime
	}
	if logSummaryItem.LogErrors != nil {
		logErrorsList := []map[string]interface{}{}
		for _, logErrorsItem := range logSummaryItem.LogErrors {
			logErrorsList = append(logErrorsList, dataSourceJobLogSummaryLogErrorsToMap(logErrorsItem))
		}
		logSummaryMap["log_errors"] = logErrorsList
	}
	if logSummaryItem.RepoDownloadJob != nil {
		repoDownloadJobList := []map[string]interface{}{}
		repoDownloadJobMap := dataSourceJobLogSummaryRepoDownloadJobToMap(*logSummaryItem.RepoDownloadJob)
		repoDownloadJobList = append(repoDownloadJobList, repoDownloadJobMap)
		logSummaryMap["repo_download_job"] = repoDownloadJobList
	}
	if logSummaryItem.ActionJob != nil {
		actionJobList := []map[string]interface{}{}
		actionJobMap := dataSourceJobLogSummaryActionJobToMap(*logSummaryItem.ActionJob)
		actionJobList = append(actionJobList, actionJobMap)
		logSummaryMap["action_job"] = actionJobList
	}

	return logSummaryMap
}

func dataSourceJobLogSummaryLogErrorsToMap(logErrorsItem schematicsv1.JobLogSummaryLogErrorsItem) (logErrorsMap map[string]interface{}) {
	logErrorsMap = map[string]interface{}{}

	if logErrorsItem.ErrorCode != nil {
		logErrorsMap["error_code"] = logErrorsItem.ErrorCode
	}
	if logErrorsItem.ErrorMsg != nil {
		logErrorsMap["error_msg"] = logErrorsItem.ErrorMsg
	}
	if logErrorsItem.ErrorCount != nil {
		logErrorsMap["error_count"] = logErrorsItem.ErrorCount
	}

	return logErrorsMap
}

func dataSourceJobLogSummaryRepoDownloadJobToMap(repoDownloadJobItem schematicsv1.JobLogSummaryRepoDownloadJob) (repoDownloadJobMap map[string]interface{}) {
	repoDownloadJobMap = map[string]interface{}{}

	if repoDownloadJobItem.ScannedFileCount != nil {
		repoDownloadJobMap["scanned_file_count"] = repoDownloadJobItem.ScannedFileCount
	}
	if repoDownloadJobItem.QuarantinedFileCount != nil {
		repoDownloadJobMap["quarantined_file_count"] = repoDownloadJobItem.QuarantinedFileCount
	}
	if repoDownloadJobItem.DetectedFiletype != nil {
		repoDownloadJobMap["detected_filetype"] = repoDownloadJobItem.DetectedFiletype
	}
	if repoDownloadJobItem.InputsCount != nil {
		repoDownloadJobMap["inputs_count"] = repoDownloadJobItem.InputsCount
	}
	if repoDownloadJobItem.OutputsCount != nil {
		repoDownloadJobMap["outputs_count"] = repoDownloadJobItem.OutputsCount
	}

	return repoDownloadJobMap
}

func dataSourceJobLogSummaryActionJobToMap(actionJobItem schematicsv1.JobLogSummaryActionJob) (actionJobMap map[string]interface{}) {
	actionJobMap = map[string]interface{}{}

	if actionJobItem.TargetCount != nil {
		actionJobMap["target_count"] = actionJobItem.TargetCount
	}
	if actionJobItem.TaskCount != nil {
		actionJobMap["task_count"] = actionJobItem.TaskCount
	}
	if actionJobItem.PlayCount != nil {
		actionJobMap["play_count"] = actionJobItem.PlayCount
	}
	if actionJobItem.Recap != nil {
		recapList := []map[string]interface{}{}
		recapMap := dataSourceJobActionJobRecapToMap(*actionJobItem.Recap)
		recapList = append(recapList, recapMap)
		actionJobMap["recap"] = recapList
	}

	return actionJobMap
}

func dataSourceJobActionJobRecapToMap(recapItem schematicsv1.JobLogSummaryActionJobRecap) (recapMap map[string]interface{}) {
	recapMap = map[string]interface{}{}

	if recapItem.Target != nil {
		recapMap["target"] = recapItem.Target
	}
	if recapItem.Ok != nil {
		recapMap["ok"] = recapItem.Ok
	}
	if recapItem.Changed != nil {
		recapMap["changed"] = recapItem.Changed
	}
	if recapItem.Failed != nil {
		recapMap["failed"] = recapItem.Failed
	}
	if recapItem.Skipped != nil {
		recapMap["skipped"] = recapItem.Skipped
	}
	if recapItem.Unreachable != nil {
		recapMap["unreachable"] = recapItem.Unreachable
	}

	return recapMap
}
