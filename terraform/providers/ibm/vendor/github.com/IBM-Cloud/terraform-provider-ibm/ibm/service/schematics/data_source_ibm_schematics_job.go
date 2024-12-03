// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIBMSchematicsJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSchematicsJobRead,

		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Job Id. Use `GET /v2/jobs` API to look up the Job Ids in your IBM Cloud account.",
			},
			"command_object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the Schematics automation resource.",
			},
			"command_object_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job command object id (workspace-id, action-id).",
			},
			"command_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schematics job command name.",
			},
			"command_parameter": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Schematics job command parameter (playbook-name).",
			},
			"command_options": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Command line options for the command.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"job_inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job inputs used by Action or Workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"job_env_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Environment variables used by the Job while performing Action or Workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"aliases": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of aliases for the variable name.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Description of the meta data.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"min_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Regex for the variable value.",
									},
									"position": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "User defined tags, while running the job.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job name, uniquely derived from the related Workspace or Action.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of your job is derived from the related action or workspace.  The description can be up to 2048 characters long in size.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource-group name derived from the related Workspace or Action.",
			},
			"submitted_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job submission time.",
			},
			"submitted_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Email address of user who submitted the job.",
			},
			"start_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job start time.",
			},
			"end_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job end time.",
			},
			"duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Duration of job execution; example 40 sec.",
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job Status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_job_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Workspace Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace name.",
									},
									"status_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Jobs.",
									},
									"status_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace job status message (eg. App1_Setup_Pending, for a 'Setup' flow in the 'App1' Workspace).",
									},
									"flow_status": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Environment Flow JOB Status.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"flow_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "flow id.",
												},
												"flow_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "flow name.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Flow Job status message - to be displayed along with the status_code;.",
												},
												"workitems": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Environment's individual workItem status details;.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"workspace_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Workspace id.",
															},
															"workspace_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "workspace name.",
															},
															"job_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "workspace job id.",
															},
															"status_code": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of Jobs.",
															},
															"status_message": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "workitem job status message;.",
															},
															"updated_at": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "workitem job status updation timestamp.",
															},
														},
													},
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"template_status": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Workspace Flow Template job status.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"template_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Template Id.",
												},
												"template_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Template name.",
												},
												"flow_index": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Index of the template in the Flow.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Template job status message (eg. VPCt1_Apply_Pending, for a 'VPCt1' Template).",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"action_job_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Action Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action name.",
									},
									"status_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Jobs.",
									},
									"status_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Action Job status message - to be displayed along with the action_status_code.",
									},
									"bastion_status_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Resources.",
									},
									"bastion_status_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Bastion status message - to be displayed along with the bastion_status_code;.",
									},
									"targets_status_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Resources.",
									},
									"targets_status_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Aggregated status message for all target resources,  to be displayed along with the targets_status_code;.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"system_job_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"system_status_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "System job message.",
									},
									"system_status_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Jobs.",
									},
									"schematics_resource_status": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "job staus for each schematics resource.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "system job status message.",
												},
												"schematics_resource_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "id for each resource which is targeted as a part of system job.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"flow_job_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Environment Flow JOB Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flow_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "flow id.",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "flow name.",
									},
									"status_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of Jobs.",
									},
									"status_message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow Job status message - to be displayed along with the status_code;.",
									},
									"workitems": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Environment's individual workItem status details;.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Workspace id.",
												},
												"workspace_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "workspace name.",
												},
												"job_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "workspace job id.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "workitem job status message;.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "workitem job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
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
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Job.",
						},
						"workspace_job_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Workspace Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Workspace name.",
									},
									"flow_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow Id.",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow name.",
									},
									"inputs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Input variables data used by the Workspace Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"outputs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Output variables data from the Workspace Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"settings": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Environment variables used by all the templates in the Workspace.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"template_data": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Input / output data of the Template in the Workspace Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"template_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Template Id.",
												},
												"template_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Template name.",
												},
												"flow_index": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Index of the template in the Flow.",
												},
												"inputs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Job inputs used by the Templates.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of aliases for the variable name.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"outputs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Job output from the Templates.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of aliases for the variable name.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"settings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Environment variables used by the template.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of aliases for the variable name.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"action_job_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Action Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow name.",
									},
									"inputs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Input variables data used by the Action Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"outputs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Output variables data from the Action Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"settings": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Environment variables used by all the templates in the Action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of aliases for the variable name.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Computed:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
									"inventory_record": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Complete inventory resource details with user inputs and system generated data.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique name of your Inventory.  The name can be up to 128 characters long and can include alphanumeric  characters, spaces, dashes, and underscores.",
												},
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Inventory id.",
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The description of your Inventory.  The description can be up to 2048 characters long in size.",
												},
												"location": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
												},
												"resource_group": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Resource-group name for the Inventory definition.  By default, Inventory will be created in Default Resource Group.",
												},
												"created_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Inventory creation time.",
												},
												"created_by": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Email address of user who created the Inventory.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Inventory updation time.",
												},
												"updated_by": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Email address of user who updated the Inventory.",
												},
												"inventories_ini": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Input inventory of host and host group for the playbook,  in the .ini file format.",
												},
												"resource_queries": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Input resource queries that is used to dynamically generate  the inventory of host and host group for the playbook.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"materialized_inventory": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Materialized inventory details used by the Action Job, in .ini format.",
									},
								},
							},
						},
						"system_job_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Controls Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Key ID for which key event is generated.",
									},
									"schematics_resource_id": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "List of the schematics resource id.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"flow_job_data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flow Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flow_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow ID.",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Flow Name.",
									},
									"workitems": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Job data used by each workitem Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"command_object_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "command object id.",
												},
												"command_object_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "command object name.",
												},
												"layers": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "layer name.",
												},
												"source_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Type of source for the Template.",
												},
												"source": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Source of templates, playbooks, or controls.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Type of source for the Template.",
															},
															"git": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Connection details to Git source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"computed_git_repo_url": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "The Complete URL which is computed by git_repo_url, git_repo_folder and branch.",
																		},
																		"git_repo_url": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "URL to the GIT Repo that can be used to clone the template.",
																		},
																		"git_token": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Personal Access Token to connect to Git URLs.",
																		},
																		"git_repo_folder": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Name of the folder in the Git Repo, that contains the template.",
																		},
																		"git_release": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Name of the release tag, used to fetch the Git Repo.",
																		},
																		"git_branch": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Name of the branch, used to fetch the Git Repo.",
																		},
																	},
																},
															},
															"catalog": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Connection details to IBM Cloud Catalog source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"catalog_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "name of the private catalog.",
																		},
																		"offering_name": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Name of the offering in the IBM Catalog.",
																		},
																		"offering_version": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Version string of the offering in the IBM Catalog.",
																		},
																		"offering_kind": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the offering, in the IBM Catalog.",
																		},
																		"offering_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Id of the offering the IBM Catalog.",
																		},
																		"offering_version_id": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Id of the offering version the IBM Catalog.",
																		},
																		"offering_repo_url": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Repo Url of the offering, in the IBM Catalog.",
																		},
																	},
																},
															},
															// "cos_bucket": {
															// 	Type:        schema.TypeList,
															// 	Computed:    true,
															// 	Description: "Connection details to a IBM Cloud Object Storage bucket.",
															// 	Elem: &schema.Resource{
															// 		Schema: map[string]*schema.Schema{
															// 			"cos_bucket_url": {
															// 				Type:        schema.TypeString,
															// 				Computed:    true,
															// 				Description: "COS Bucket Url.",
															// 			},
															// 		},
															// 	},
															// },
														},
													},
												},
												"inputs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Input variables data for the workItem used in FlowJob.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of aliases for the variable name.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"outputs": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Output variables for the workItem.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of aliases for the variable name.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"settings": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Environment variables for the workItem.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of aliases for the variable name.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Computed:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Computed:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem: &schema.Schema{
																				Type: schema.TypeString,
																			},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Computed:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"last_job": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Status of the last job executed by the workitem.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"command_object": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Name of the Schematics automation resource.",
															},
															"command_object_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "command object name (workspace_name/action_name).",
															},
															"command_object_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Workitem command object id, maps to workspace_id or action_id.",
															},
															"command_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Schematics job command name.",
															},
															"job_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Workspace job id.",
															},
															"job_status": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Status of Jobs.",
															},
														},
													},
												},
												"updated_at": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
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
			"bastion": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Describes a bastion resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Bastion Name(Unique).",
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference to the Inventory resource definition.",
						},
					},
				},
			},
			"log_summary": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job log summary record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Workspace Id.",
						},
						"job_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of Job.",
						},
						"log_start_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job log start timestamp.",
						},
						"log_analyzed_till": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Job log update timestamp.",
						},
						"elapsed_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Job log elapsed time (log_analyzed_till - log_start_at).",
						},
						"log_errors": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Job log errors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Error code in the Log.",
									},
									"error_msg": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Summary error message in the log.",
									},
									"error_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of occurrence.",
									},
								},
							},
						},
						"repo_download_job": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Repo download Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scanned_file_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of files scanned.",
									},
									"quarantined_file_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of files quarantined.",
									},
									"detected_filetype": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Detected template or data file type.",
									},
									"inputs_count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of inputs detected.",
									},
									"outputs_count": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Number of outputs detected.",
									},
								},
							},
						},
						"workspace_job": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Workspace Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resources_add": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of resources add.",
									},
									"resources_modify": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of resources modify.",
									},
									"resources_destroy": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of resources destroy.",
									},
								},
							},
						},
						"flow_job": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flow Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workitems_completed": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of workitems completed successfully.",
									},
									"workitems_pending": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of workitems pending in the flow.",
									},
									"workitems_failed": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of workitems failed.",
									},
									"workitems": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "workspace ID.",
												},
												"job_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "workspace JOB ID.",
												},
												"resources_add": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of resources add.",
												},
												"resources_modify": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of resources modify.",
												},
												"resources_destroy": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of resources destroy.",
												},
												"log_url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Log url for job.",
												},
											},
										},
									},
								},
							},
						},
						"action_job": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Flow Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of targets or hosts.",
									},
									"task_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of tasks in playbook.",
									},
									"play_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of plays in playbook.",
									},
									"recap": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Recap records.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "List of target or host name.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
												"ok": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of OK.",
												},
												"changed": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of changed.",
												},
												"failed": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of failed.",
												},
												"skipped": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Number of skipped.",
												},
												"unreachable": {
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
						"system_job": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "System Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_count": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "number of targets or hosts.",
									},
									"success": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of passed.",
									},
									"failed": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Number of failed.",
									},
								},
							},
						},
					},
				},
			},
			"log_store_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job log store URL.",
			},
			"state_store_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job state store URL.",
			},
			"results_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job results store URL.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Job status updation timestamp.",
			},
		},
	}
}

func dataSourceIBMSchematicsJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}
	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}
	getJobOptions := &schematicsv1.GetJobOptions{}

	getJobOptions.SetJobID(d.Get("job_id").(string))

	job, response, err := schematicsClient.GetJobWithContext(context, getJobOptions)
	if err != nil {
		log.Printf("[DEBUG] GetJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(*getJobOptions.JobID)
	if err = d.Set("command_object", job.CommandObject); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_object: %s", err))
	}
	if err = d.Set("command_object_id", job.CommandObjectID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_object_id: %s", err))
	}
	if err = d.Set("command_name", job.CommandName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_name: %s", err))
	}
	if err = d.Set("command_parameter", job.CommandParameter); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_parameter: %s", err))
	}

	if job.Inputs != nil {
		err = d.Set("job_inputs", dataSourceJobFlattenInputs(job.Inputs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting job_inputs %s", err))
		}
	}

	if job.Settings != nil {
		err = d.Set("job_env_settings", dataSourceJobFlattenSettings(job.Settings))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting job_env_settings %s", err))
		}
	}
	if err = d.Set("id", job.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting id: %s", err))
	}
	if err = d.Set("name", job.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("description", job.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
	}
	if err = d.Set("location", job.Location); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting location: %s", err))
	}
	if err = d.Set("resource_group", job.ResourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
	}
	if err = d.Set("submitted_at", flex.DateTimeToString(job.SubmittedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting submitted_at: %s", err))
	}
	if err = d.Set("submitted_by", job.SubmittedBy); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting submitted_by: %s", err))
	}
	if err = d.Set("start_at", flex.DateTimeToString(job.StartAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting start_at: %s", err))
	}
	if err = d.Set("end_at", flex.DateTimeToString(job.EndAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting end_at: %s", err))
	}
	if err = d.Set("duration", job.Duration); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting duration: %s", err))
	}

	if job.Status != nil {
		err = d.Set("status", dataSourceJobFlattenStatus(*job.Status))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting status %s", err))
		}
	}

	if job.Data != nil {
		err = d.Set("data", dataSourceJobFlattenData(*job.Data))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting data %s", err))
		}
	}

	if job.Bastion != nil {
		err = d.Set("bastion", dataSourceJobFlattenBastion(*job.Bastion))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting bastion %s", err))
		}
	}

	if job.LogSummary != nil {
		err = d.Set("log_summary", dataSourceJobFlattenLogSummary(*job.LogSummary))
		if err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting log_summary %s", err))
		}
	}
	if err = d.Set("log_store_url", job.LogStoreURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting log_store_url: %s", err))
	}
	if err = d.Set("state_store_url", job.StateStoreURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting state_store_url: %s", err))
	}
	if err = d.Set("results_url", job.ResultsURL); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting results_url: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(job.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting updated_at: %s", err))
	}

	return nil
}

func dataSourceJobFlattenInputs(result []schematicsv1.VariableData) (jobInputs []map[string]interface{}) {
	for _, jobInputsItem := range result {
		jobInputs = append(jobInputs, dataSourceJobInputsToMap(jobInputsItem))
	}

	return jobInputs
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

func dataSourceJobFlattenSettings(result []schematicsv1.VariableData) (jobEnvSettings []map[string]interface{}) {
	for _, jobEnvSettingsItem := range result {
		jobEnvSettings = append(jobEnvSettings, dataSourceJobSettingsToMap(jobEnvSettingsItem))
	}

	return jobEnvSettings
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

	if statusItem.WorkspaceJobStatus != nil {
		workspaceJobStatusList := []map[string]interface{}{}
		workspaceJobStatusMap := dataSourceJobStatusWorkspaceJobStatusToMap(*statusItem.WorkspaceJobStatus)
		workspaceJobStatusList = append(workspaceJobStatusList, workspaceJobStatusMap)
		statusMap["workspace_job_status"] = workspaceJobStatusList
	}
	if statusItem.ActionJobStatus != nil {
		actionJobStatusList := []map[string]interface{}{}
		actionJobStatusMap := dataSourceJobStatusActionJobStatusToMap(*statusItem.ActionJobStatus)
		actionJobStatusList = append(actionJobStatusList, actionJobStatusMap)
		statusMap["action_job_status"] = actionJobStatusList
	}
	if statusItem.SystemJobStatus != nil {
		systemJobStatusList := []map[string]interface{}{}
		systemJobStatusMap := dataSourceJobStatusSystemJobStatusToMap(*statusItem.SystemJobStatus)
		systemJobStatusList = append(systemJobStatusList, systemJobStatusMap)
		statusMap["system_job_status"] = systemJobStatusList
	}
	if statusItem.FlowJobStatus != nil {
		flowJobStatusList := []map[string]interface{}{}
		flowJobStatusMap := dataSourceJobStatusFlowJobStatusToMap(*statusItem.FlowJobStatus)
		flowJobStatusList = append(flowJobStatusList, flowJobStatusMap)
		statusMap["flow_job_status"] = flowJobStatusList
	}

	return statusMap
}

func dataSourceJobStatusWorkspaceJobStatusToMap(workspaceJobStatusItem schematicsv1.JobStatusWorkspace) (workspaceJobStatusMap map[string]interface{}) {
	workspaceJobStatusMap = map[string]interface{}{}

	if workspaceJobStatusItem.WorkspaceName != nil {
		workspaceJobStatusMap["workspace_name"] = workspaceJobStatusItem.WorkspaceName
	}
	if workspaceJobStatusItem.StatusCode != nil {
		workspaceJobStatusMap["status_code"] = workspaceJobStatusItem.StatusCode
	}
	if workspaceJobStatusItem.StatusMessage != nil {
		workspaceJobStatusMap["status_message"] = workspaceJobStatusItem.StatusMessage
	}
	if workspaceJobStatusItem.FlowStatus != nil {
		flowStatusList := []map[string]interface{}{}
		flowStatusMap := dataSourceJobWorkspaceJobStatusFlowStatusToMap(*workspaceJobStatusItem.FlowStatus)
		flowStatusList = append(flowStatusList, flowStatusMap)
		workspaceJobStatusMap["flow_status"] = flowStatusList
	}
	if workspaceJobStatusItem.TemplateStatus != nil {
		templateStatusList := []map[string]interface{}{}
		for _, templateStatusItem := range workspaceJobStatusItem.TemplateStatus {
			templateStatusList = append(templateStatusList, dataSourceJobWorkspaceJobStatusTemplateStatusToMap(templateStatusItem))
		}
		workspaceJobStatusMap["template_status"] = templateStatusList
	}
	if workspaceJobStatusItem.UpdatedAt != nil {
		workspaceJobStatusMap["updated_at"] = workspaceJobStatusItem.UpdatedAt.String()
	}

	return workspaceJobStatusMap
}

func dataSourceJobWorkspaceJobStatusFlowStatusToMap(flowStatusItem schematicsv1.JobStatusFlow) (flowStatusMap map[string]interface{}) {
	flowStatusMap = map[string]interface{}{}

	if flowStatusItem.FlowID != nil {
		flowStatusMap["flow_id"] = flowStatusItem.FlowID
	}
	if flowStatusItem.FlowName != nil {
		flowStatusMap["flow_name"] = flowStatusItem.FlowName
	}
	if flowStatusItem.StatusCode != nil {
		flowStatusMap["status_code"] = flowStatusItem.StatusCode
	}
	if flowStatusItem.StatusMessage != nil {
		flowStatusMap["status_message"] = flowStatusItem.StatusMessage
	}
	if flowStatusItem.Workitems != nil {
		workitemsList := []map[string]interface{}{}
		for _, workitemsItem := range flowStatusItem.Workitems {
			workitemsList = append(workitemsList, dataSourceJobFlowStatusWorkitemsToMap(workitemsItem))
		}
		flowStatusMap["workitems"] = workitemsList
	}
	if flowStatusItem.UpdatedAt != nil {
		flowStatusMap["updated_at"] = flowStatusItem.UpdatedAt.String()
	}

	return flowStatusMap
}

func dataSourceJobFlowStatusWorkitemsToMap(workitemsItem schematicsv1.JobStatusWorkitem) (workitemsMap map[string]interface{}) {
	workitemsMap = map[string]interface{}{}

	if workitemsItem.WorkspaceID != nil {
		workitemsMap["workspace_id"] = workitemsItem.WorkspaceID
	}
	if workitemsItem.WorkspaceName != nil {
		workitemsMap["workspace_name"] = workitemsItem.WorkspaceName
	}
	if workitemsItem.JobID != nil {
		workitemsMap["job_id"] = workitemsItem.JobID
	}
	if workitemsItem.StatusCode != nil {
		workitemsMap["status_code"] = workitemsItem.StatusCode
	}
	if workitemsItem.StatusMessage != nil {
		workitemsMap["status_message"] = workitemsItem.StatusMessage
	}
	if workitemsItem.UpdatedAt != nil {
		workitemsMap["updated_at"] = workitemsItem.UpdatedAt.String()
	}

	return workitemsMap
}

func dataSourceJobWorkspaceJobStatusTemplateStatusToMap(templateStatusItem schematicsv1.JobStatusTemplate) (templateStatusMap map[string]interface{}) {
	templateStatusMap = map[string]interface{}{}

	if templateStatusItem.TemplateID != nil {
		templateStatusMap["template_id"] = templateStatusItem.TemplateID
	}
	if templateStatusItem.TemplateName != nil {
		templateStatusMap["template_name"] = templateStatusItem.TemplateName
	}
	if templateStatusItem.FlowIndex != nil {
		templateStatusMap["flow_index"] = templateStatusItem.FlowIndex
	}
	if templateStatusItem.StatusCode != nil {
		templateStatusMap["status_code"] = templateStatusItem.StatusCode
	}
	if templateStatusItem.StatusMessage != nil {
		templateStatusMap["status_message"] = templateStatusItem.StatusMessage
	}
	if templateStatusItem.UpdatedAt != nil {
		templateStatusMap["updated_at"] = templateStatusItem.UpdatedAt.String()
	}

	return templateStatusMap
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

func dataSourceJobStatusSystemJobStatusToMap(systemJobStatusItem schematicsv1.JobStatusSystem) (systemJobStatusMap map[string]interface{}) {
	systemJobStatusMap = map[string]interface{}{}

	if systemJobStatusItem.SystemStatusMessage != nil {
		systemJobStatusMap["system_status_message"] = systemJobStatusItem.SystemStatusMessage
	}
	if systemJobStatusItem.SystemStatusCode != nil {
		systemJobStatusMap["system_status_code"] = systemJobStatusItem.SystemStatusCode
	}
	if systemJobStatusItem.SchematicsResourceStatus != nil {
		schematicsResourceStatusList := []map[string]interface{}{}
		for _, schematicsResourceStatusItem := range systemJobStatusItem.SchematicsResourceStatus {
			schematicsResourceStatusList = append(schematicsResourceStatusList, dataSourceJobSystemJobStatusSchematicsResourceStatusToMap(schematicsResourceStatusItem))
		}
		systemJobStatusMap["schematics_resource_status"] = schematicsResourceStatusList
	}
	if systemJobStatusItem.UpdatedAt != nil {
		systemJobStatusMap["updated_at"] = systemJobStatusItem.UpdatedAt.String()
	}

	return systemJobStatusMap
}

func dataSourceJobSystemJobStatusSchematicsResourceStatusToMap(schematicsResourceStatusItem schematicsv1.JobStatusSchematicsResources) (schematicsResourceStatusMap map[string]interface{}) {
	schematicsResourceStatusMap = map[string]interface{}{}

	if schematicsResourceStatusItem.StatusCode != nil {
		schematicsResourceStatusMap["status_code"] = schematicsResourceStatusItem.StatusCode
	}
	if schematicsResourceStatusItem.StatusMessage != nil {
		schematicsResourceStatusMap["status_message"] = schematicsResourceStatusItem.StatusMessage
	}
	if schematicsResourceStatusItem.SchematicsResourceID != nil {
		schematicsResourceStatusMap["schematics_resource_id"] = schematicsResourceStatusItem.SchematicsResourceID
	}
	if schematicsResourceStatusItem.UpdatedAt != nil {
		schematicsResourceStatusMap["updated_at"] = schematicsResourceStatusItem.UpdatedAt.String()
	}

	return schematicsResourceStatusMap
}

func dataSourceJobStatusFlowJobStatusToMap(flowJobStatusItem schematicsv1.JobStatusFlow) (flowJobStatusMap map[string]interface{}) {
	flowJobStatusMap = map[string]interface{}{}

	if flowJobStatusItem.FlowID != nil {
		flowJobStatusMap["flow_id"] = flowJobStatusItem.FlowID
	}
	if flowJobStatusItem.FlowName != nil {
		flowJobStatusMap["flow_name"] = flowJobStatusItem.FlowName
	}
	if flowJobStatusItem.StatusCode != nil {
		flowJobStatusMap["status_code"] = flowJobStatusItem.StatusCode
	}
	if flowJobStatusItem.StatusMessage != nil {
		flowJobStatusMap["status_message"] = flowJobStatusItem.StatusMessage
	}
	if flowJobStatusItem.Workitems != nil {
		workitemsList := []map[string]interface{}{}
		for _, workitemsItem := range flowJobStatusItem.Workitems {
			workitemsList = append(workitemsList, dataSourceJobFlowJobStatusWorkitemsToMap(workitemsItem))
		}
		flowJobStatusMap["workitems"] = workitemsList
	}
	if flowJobStatusItem.UpdatedAt != nil {
		flowJobStatusMap["updated_at"] = flowJobStatusItem.UpdatedAt.String()
	}

	return flowJobStatusMap
}

func dataSourceJobFlowJobStatusWorkitemsToMap(workitemsItem schematicsv1.JobStatusWorkitem) (workitemsMap map[string]interface{}) {
	workitemsMap = map[string]interface{}{}

	if workitemsItem.WorkspaceID != nil {
		workitemsMap["workspace_id"] = workitemsItem.WorkspaceID
	}
	if workitemsItem.WorkspaceName != nil {
		workitemsMap["workspace_name"] = workitemsItem.WorkspaceName
	}
	if workitemsItem.JobID != nil {
		workitemsMap["job_id"] = workitemsItem.JobID
	}
	if workitemsItem.StatusCode != nil {
		workitemsMap["status_code"] = workitemsItem.StatusCode
	}
	if workitemsItem.StatusMessage != nil {
		workitemsMap["status_message"] = workitemsItem.StatusMessage
	}
	if workitemsItem.UpdatedAt != nil {
		workitemsMap["updated_at"] = workitemsItem.UpdatedAt.String()
	}

	return workitemsMap
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
	if dataItem.WorkspaceJobData != nil {
		workspaceJobDataList := []map[string]interface{}{}
		workspaceJobDataMap := dataSourceJobDataWorkspaceJobDataToMap(*dataItem.WorkspaceJobData)
		workspaceJobDataList = append(workspaceJobDataList, workspaceJobDataMap)
		dataMap["workspace_job_data"] = workspaceJobDataList
	}
	if dataItem.ActionJobData != nil {
		actionJobDataList := []map[string]interface{}{}
		actionJobDataMap := dataSourceJobDataActionJobDataToMap(*dataItem.ActionJobData)
		actionJobDataList = append(actionJobDataList, actionJobDataMap)
		dataMap["action_job_data"] = actionJobDataList
	}
	if dataItem.SystemJobData != nil {
		systemJobDataList := []map[string]interface{}{}
		systemJobDataMap := dataSourceJobDataSystemJobDataToMap(*dataItem.SystemJobData)
		systemJobDataList = append(systemJobDataList, systemJobDataMap)
		dataMap["system_job_data"] = systemJobDataList
	}
	if dataItem.FlowJobData != nil {
		flowJobDataList := []map[string]interface{}{}
		flowJobDataMap := dataSourceJobDataFlowJobDataToMap(*dataItem.FlowJobData)
		flowJobDataList = append(flowJobDataList, flowJobDataMap)
		dataMap["flow_job_data"] = flowJobDataList
	}

	return dataMap
}

func dataSourceJobDataWorkspaceJobDataToMap(workspaceJobDataItem schematicsv1.JobDataWorkspace) (workspaceJobDataMap map[string]interface{}) {
	workspaceJobDataMap = map[string]interface{}{}

	if workspaceJobDataItem.WorkspaceName != nil {
		workspaceJobDataMap["workspace_name"] = workspaceJobDataItem.WorkspaceName
	}
	if workspaceJobDataItem.FlowID != nil {
		workspaceJobDataMap["flow_id"] = workspaceJobDataItem.FlowID
	}
	if workspaceJobDataItem.FlowName != nil {
		workspaceJobDataMap["flow_name"] = workspaceJobDataItem.FlowName
	}
	if workspaceJobDataItem.Inputs != nil {
		inputsList := []map[string]interface{}{}
		for _, inputsItem := range workspaceJobDataItem.Inputs {
			inputsList = append(inputsList, dataSourceJobWorkspaceJobDataInputsToMap(inputsItem))
		}
		workspaceJobDataMap["inputs"] = inputsList
	}
	if workspaceJobDataItem.Outputs != nil {
		outputsList := []map[string]interface{}{}
		for _, outputsItem := range workspaceJobDataItem.Outputs {
			outputsList = append(outputsList, dataSourceJobWorkspaceJobDataOutputsToMap(outputsItem))
		}
		workspaceJobDataMap["outputs"] = outputsList
	}
	if workspaceJobDataItem.Settings != nil {
		settingsList := []map[string]interface{}{}
		for _, settingsItem := range workspaceJobDataItem.Settings {
			settingsList = append(settingsList, dataSourceJobWorkspaceJobDataSettingsToMap(settingsItem))
		}
		workspaceJobDataMap["settings"] = settingsList
	}
	if workspaceJobDataItem.TemplateData != nil {
		templateDataList := []map[string]interface{}{}
		for _, templateDataItem := range workspaceJobDataItem.TemplateData {
			templateDataList = append(templateDataList, dataSourceJobWorkspaceJobDataTemplateDataToMap(templateDataItem))
		}
		workspaceJobDataMap["template_data"] = templateDataList
	}
	if workspaceJobDataItem.UpdatedAt != nil {
		workspaceJobDataMap["updated_at"] = workspaceJobDataItem.UpdatedAt.String()
	}

	return workspaceJobDataMap
}

func dataSourceJobWorkspaceJobDataInputsToMap(inputsItem schematicsv1.VariableData) (inputsMap map[string]interface{}) {
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

func dataSourceJobWorkspaceJobDataOutputsToMap(outputsItem schematicsv1.VariableData) (outputsMap map[string]interface{}) {
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

func dataSourceJobWorkspaceJobDataSettingsToMap(settingsItem schematicsv1.VariableData) (settingsMap map[string]interface{}) {
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

func dataSourceJobWorkspaceJobDataTemplateDataToMap(templateDataItem schematicsv1.JobDataTemplate) (templateDataMap map[string]interface{}) {
	templateDataMap = map[string]interface{}{}

	if templateDataItem.TemplateID != nil {
		templateDataMap["template_id"] = templateDataItem.TemplateID
	}
	if templateDataItem.TemplateName != nil {
		templateDataMap["template_name"] = templateDataItem.TemplateName
	}
	if templateDataItem.FlowIndex != nil {
		templateDataMap["flow_index"] = templateDataItem.FlowIndex
	}
	if templateDataItem.Inputs != nil {
		inputsList := []map[string]interface{}{}
		for _, inputsItem := range templateDataItem.Inputs {
			inputsList = append(inputsList, dataSourceJobTemplateDataInputsToMap(inputsItem))
		}
		templateDataMap["inputs"] = inputsList
	}
	if templateDataItem.Outputs != nil {
		outputsList := []map[string]interface{}{}
		for _, outputsItem := range templateDataItem.Outputs {
			outputsList = append(outputsList, dataSourceJobTemplateDataOutputsToMap(outputsItem))
		}
		templateDataMap["outputs"] = outputsList
	}
	if templateDataItem.Settings != nil {
		settingsList := []map[string]interface{}{}
		for _, settingsItem := range templateDataItem.Settings {
			settingsList = append(settingsList, dataSourceJobTemplateDataSettingsToMap(settingsItem))
		}
		templateDataMap["settings"] = settingsList
	}
	if templateDataItem.UpdatedAt != nil {
		templateDataMap["updated_at"] = templateDataItem.UpdatedAt.String()
	}

	return templateDataMap
}

func dataSourceJobTemplateDataInputsToMap(inputsItem schematicsv1.VariableData) (inputsMap map[string]interface{}) {
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

func dataSourceJobTemplateDataOutputsToMap(outputsItem schematicsv1.VariableData) (outputsMap map[string]interface{}) {
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

func dataSourceJobTemplateDataSettingsToMap(settingsItem schematicsv1.VariableData) (settingsMap map[string]interface{}) {
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
	if actionJobDataItem.InventoryRecord != nil {
		inventoryRecordList := []map[string]interface{}{}
		inventoryRecordMap := dataSourceJobActionJobDataInventoryRecordToMap(*actionJobDataItem.InventoryRecord)
		inventoryRecordList = append(inventoryRecordList, inventoryRecordMap)
		actionJobDataMap["inventory_record"] = inventoryRecordList
	}
	if actionJobDataItem.MaterializedInventory != nil {
		actionJobDataMap["materialized_inventory"] = actionJobDataItem.MaterializedInventory
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

func dataSourceJobActionJobDataInventoryRecordToMap(inventoryRecordItem schematicsv1.InventoryResourceRecord) (inventoryRecordMap map[string]interface{}) {
	inventoryRecordMap = map[string]interface{}{}

	if inventoryRecordItem.Name != nil {
		inventoryRecordMap["name"] = inventoryRecordItem.Name
	}
	if inventoryRecordItem.ID != nil {
		inventoryRecordMap["id"] = inventoryRecordItem.ID
	}
	if inventoryRecordItem.Description != nil {
		inventoryRecordMap["description"] = inventoryRecordItem.Description
	}
	if inventoryRecordItem.Location != nil {
		inventoryRecordMap["location"] = inventoryRecordItem.Location
	}
	if inventoryRecordItem.ResourceGroup != nil {
		inventoryRecordMap["resource_group"] = inventoryRecordItem.ResourceGroup
	}
	if inventoryRecordItem.CreatedAt != nil {
		inventoryRecordMap["created_at"] = inventoryRecordItem.CreatedAt.String()
	}
	if inventoryRecordItem.CreatedBy != nil {
		inventoryRecordMap["created_by"] = inventoryRecordItem.CreatedBy
	}
	if inventoryRecordItem.UpdatedAt != nil {
		inventoryRecordMap["updated_at"] = inventoryRecordItem.UpdatedAt.String()
	}
	if inventoryRecordItem.UpdatedBy != nil {
		inventoryRecordMap["updated_by"] = inventoryRecordItem.UpdatedBy
	}
	if inventoryRecordItem.InventoriesIni != nil {
		inventoryRecordMap["inventories_ini"] = inventoryRecordItem.InventoriesIni
	}
	if inventoryRecordItem.ResourceQueries != nil {
		inventoryRecordMap["resource_queries"] = inventoryRecordItem.ResourceQueries
	}

	return inventoryRecordMap
}

func dataSourceJobDataSystemJobDataToMap(systemJobDataItem schematicsv1.JobDataSystem) (systemJobDataMap map[string]interface{}) {
	systemJobDataMap = map[string]interface{}{}

	if systemJobDataItem.KeyID != nil {
		systemJobDataMap["key_id"] = systemJobDataItem.KeyID
	}
	if systemJobDataItem.SchematicsResourceID != nil {
		systemJobDataMap["schematics_resource_id"] = systemJobDataItem.SchematicsResourceID
	}
	if systemJobDataItem.UpdatedAt != nil {
		systemJobDataMap["updated_at"] = systemJobDataItem.UpdatedAt.String()
	}

	return systemJobDataMap
}

func dataSourceJobDataFlowJobDataToMap(flowJobDataItem schematicsv1.JobDataFlow) (flowJobDataMap map[string]interface{}) {
	flowJobDataMap = map[string]interface{}{}

	if flowJobDataItem.FlowID != nil {
		flowJobDataMap["flow_id"] = flowJobDataItem.FlowID
	}
	if flowJobDataItem.FlowName != nil {
		flowJobDataMap["flow_name"] = flowJobDataItem.FlowName
	}
	if flowJobDataItem.Workitems != nil {
		workitemsList := []map[string]interface{}{}
		for _, workitemsItem := range flowJobDataItem.Workitems {
			workitemsList = append(workitemsList, dataSourceJobFlowJobDataWorkitemsToMap(workitemsItem))
		}
		flowJobDataMap["workitems"] = workitemsList
	}
	if flowJobDataItem.UpdatedAt != nil {
		flowJobDataMap["updated_at"] = flowJobDataItem.UpdatedAt.String()
	}

	return flowJobDataMap
}

func dataSourceJobFlowJobDataWorkitemsToMap(workitemsItem schematicsv1.JobDataWorkItem) (workitemsMap map[string]interface{}) {
	workitemsMap = map[string]interface{}{}

	if workitemsItem.CommandObjectID != nil {
		workitemsMap["command_object_id"] = workitemsItem.CommandObjectID
	}
	if workitemsItem.CommandObjectName != nil {
		workitemsMap["command_object_name"] = workitemsItem.CommandObjectName
	}
	if workitemsItem.Layers != nil {
		workitemsMap["layers"] = workitemsItem.Layers
	}
	if workitemsItem.SourceType != nil {
		workitemsMap["source_type"] = workitemsItem.SourceType
	}
	if workitemsItem.Source != nil {
		sourceList := []map[string]interface{}{}
		sourceMap := dataSourceJobWorkitemsSourceToMap(*workitemsItem.Source)
		sourceList = append(sourceList, sourceMap)
		workitemsMap["source"] = sourceList
	}
	if workitemsItem.Inputs != nil {
		inputsList := []map[string]interface{}{}
		for _, inputsItem := range workitemsItem.Inputs {
			inputsList = append(inputsList, dataSourceJobWorkitemsInputsToMap(inputsItem))
		}
		workitemsMap["inputs"] = inputsList
	}
	if workitemsItem.Outputs != nil {
		outputsList := []map[string]interface{}{}
		for _, outputsItem := range workitemsItem.Outputs {
			outputsList = append(outputsList, dataSourceJobWorkitemsOutputsToMap(outputsItem))
		}
		workitemsMap["outputs"] = outputsList
	}
	if workitemsItem.Settings != nil {
		settingsList := []map[string]interface{}{}
		for _, settingsItem := range workitemsItem.Settings {
			settingsList = append(settingsList, dataSourceJobWorkitemsSettingsToMap(settingsItem))
		}
		workitemsMap["settings"] = settingsList
	}
	if workitemsItem.LastJob != nil {
		lastJobList := []map[string]interface{}{}
		lastJobMap := dataSourceJobWorkitemsLastJobToMap(*workitemsItem.LastJob)
		lastJobList = append(lastJobList, lastJobMap)
		workitemsMap["last_job"] = lastJobList
	}
	if workitemsItem.UpdatedAt != nil {
		workitemsMap["updated_at"] = workitemsItem.UpdatedAt.String()
	}

	return workitemsMap
}

func dataSourceJobWorkitemsSourceToMap(sourceItem schematicsv1.ExternalSource) (sourceMap map[string]interface{}) {
	sourceMap = map[string]interface{}{}

	if sourceItem.SourceType != nil {
		sourceMap["source_type"] = sourceItem.SourceType
	}
	if sourceItem.Git != nil {
		gitList := []map[string]interface{}{}
		gitMap := dataSourceJobSourceGitToMap(*sourceItem.Git)
		gitList = append(gitList, gitMap)
		sourceMap["git"] = gitList
	}
	if sourceItem.Catalog != nil {
		catalogList := []map[string]interface{}{}
		catalogMap := dataSourceJobSourceCatalogToMap(*sourceItem.Catalog)
		catalogList = append(catalogList, catalogMap)
		sourceMap["catalog"] = catalogList
	}
	// if sourceItem.CosBucket != nil {
	// 	cosBucketList := []map[string]interface{}{}
	// 	cosBucketMap := dataSourceJobSourceCosBucketToMap(*sourceItem.CosBucket)
	// 	cosBucketList = append(cosBucketList, cosBucketMap)
	// 	sourceMap["cos_bucket"] = cosBucketList
	// }

	return sourceMap
}

func dataSourceJobSourceGitToMap(gitItem schematicsv1.GitSource) (gitMap map[string]interface{}) {
	gitMap = map[string]interface{}{}

	if gitItem.ComputedGitRepoURL != nil {
		gitMap["computed_git_repo_url"] = gitItem.ComputedGitRepoURL
	}
	if gitItem.GitRepoURL != nil {
		gitMap["git_repo_url"] = gitItem.GitRepoURL
	}
	if gitItem.GitToken != nil {
		gitMap["git_token"] = gitItem.GitToken
	}
	if gitItem.GitRepoFolder != nil {
		gitMap["git_repo_folder"] = gitItem.GitRepoFolder
	}
	if gitItem.GitRelease != nil {
		gitMap["git_release"] = gitItem.GitRelease
	}
	if gitItem.GitBranch != nil {
		gitMap["git_branch"] = gitItem.GitBranch
	}

	return gitMap
}

func dataSourceJobSourceCatalogToMap(catalogItem schematicsv1.CatalogSource) (catalogMap map[string]interface{}) {
	catalogMap = map[string]interface{}{}

	if catalogItem.CatalogName != nil {
		catalogMap["catalog_name"] = catalogItem.CatalogName
	}
	if catalogItem.OfferingName != nil {
		catalogMap["offering_name"] = catalogItem.OfferingName
	}
	if catalogItem.OfferingVersion != nil {
		catalogMap["offering_version"] = catalogItem.OfferingVersion
	}
	if catalogItem.OfferingKind != nil {
		catalogMap["offering_kind"] = catalogItem.OfferingKind
	}
	if catalogItem.OfferingID != nil {
		catalogMap["offering_id"] = catalogItem.OfferingID
	}
	if catalogItem.OfferingVersionID != nil {
		catalogMap["offering_version_id"] = catalogItem.OfferingVersionID
	}
	if catalogItem.OfferingRepoURL != nil {
		catalogMap["offering_repo_url"] = catalogItem.OfferingRepoURL
	}

	return catalogMap
}

// func dataSourceJobSourceCosBucketToMap(cosBucketItem schematicsv1.ExternalSourceCosBucket) (cosBucketMap map[string]interface{}) {
// 	cosBucketMap = map[string]interface{}{}

// 	if cosBucketItem.CosBucketURL != nil {
// 		cosBucketMap["cos_bucket_url"] = cosBucketItem.CosBucketURL
// 	}

// 	return cosBucketMap
// }

func dataSourceJobWorkitemsInputsToMap(inputsItem schematicsv1.VariableData) (inputsMap map[string]interface{}) {
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

func dataSourceJobWorkitemsOutputsToMap(outputsItem schematicsv1.VariableData) (outputsMap map[string]interface{}) {
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

func dataSourceJobWorkitemsSettingsToMap(settingsItem schematicsv1.VariableData) (settingsMap map[string]interface{}) {
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

func dataSourceJobWorkitemsLastJobToMap(lastJobItem schematicsv1.JobDataWorkItemLastJob) (lastJobMap map[string]interface{}) {
	lastJobMap = map[string]interface{}{}

	if lastJobItem.CommandObject != nil {
		lastJobMap["command_object"] = lastJobItem.CommandObject
	}
	if lastJobItem.CommandObjectName != nil {
		lastJobMap["command_object_name"] = lastJobItem.CommandObjectName
	}
	if lastJobItem.CommandObjectID != nil {
		lastJobMap["command_object_id"] = lastJobItem.CommandObjectID
	}
	if lastJobItem.CommandName != nil {
		lastJobMap["command_name"] = lastJobItem.CommandName
	}
	if lastJobItem.JobID != nil {
		lastJobMap["job_id"] = lastJobItem.JobID
	}
	if lastJobItem.JobStatus != nil {
		lastJobMap["job_status"] = lastJobItem.JobStatus
	}

	return lastJobMap
}

func dataSourceJobFlattenBastion(result schematicsv1.BastionResourceDefinition) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceJobBastionToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceJobBastionToMap(bastionItem schematicsv1.BastionResourceDefinition) (bastionMap map[string]interface{}) {
	bastionMap = map[string]interface{}{}

	if bastionItem.Name != nil {
		bastionMap["name"] = bastionItem.Name
	}
	if bastionItem.Host != nil {
		bastionMap["host"] = bastionItem.Host
	}

	return bastionMap
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
	if logSummaryItem.WorkspaceJob != nil {
		workspaceJobList := []map[string]interface{}{}
		workspaceJobMap := dataSourceJobLogSummaryWorkspaceJobToMap(*logSummaryItem.WorkspaceJob)
		workspaceJobList = append(workspaceJobList, workspaceJobMap)
		logSummaryMap["workspace_job"] = workspaceJobList
	}
	if logSummaryItem.FlowJob != nil {
		flowJobList := []map[string]interface{}{}
		flowJobMap := dataSourceJobLogSummaryFlowJobToMap(*logSummaryItem.FlowJob)
		flowJobList = append(flowJobList, flowJobMap)
		logSummaryMap["flow_job"] = flowJobList
	}
	if logSummaryItem.ActionJob != nil {
		actionJobList := []map[string]interface{}{}
		actionJobMap := dataSourceJobLogSummaryActionJobToMap(*logSummaryItem.ActionJob)
		actionJobList = append(actionJobList, actionJobMap)
		logSummaryMap["action_job"] = actionJobList
	}
	if logSummaryItem.SystemJob != nil {
		systemJobList := []map[string]interface{}{}
		systemJobMap := dataSourceJobLogSummarySystemJobToMap(*logSummaryItem.SystemJob)
		systemJobList = append(systemJobList, systemJobMap)
		logSummaryMap["system_job"] = systemJobList
	}

	return logSummaryMap
}

func dataSourceJobLogSummaryLogErrorsToMap(logErrorsItem schematicsv1.JobLogSummaryLogErrors) (logErrorsMap map[string]interface{}) {
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

func dataSourceJobLogSummaryWorkspaceJobToMap(workspaceJobItem schematicsv1.JobLogSummaryWorkspaceJob) (workspaceJobMap map[string]interface{}) {
	workspaceJobMap = map[string]interface{}{}

	if workspaceJobItem.ResourcesAdd != nil {
		workspaceJobMap["resources_add"] = workspaceJobItem.ResourcesAdd
	}
	if workspaceJobItem.ResourcesModify != nil {
		workspaceJobMap["resources_modify"] = workspaceJobItem.ResourcesModify
	}
	if workspaceJobItem.ResourcesDestroy != nil {
		workspaceJobMap["resources_destroy"] = workspaceJobItem.ResourcesDestroy
	}

	return workspaceJobMap
}

func dataSourceJobLogSummaryFlowJobToMap(flowJobItem schematicsv1.JobLogSummaryFlowJob) (flowJobMap map[string]interface{}) {
	flowJobMap = map[string]interface{}{}

	if flowJobItem.WorkitemsCompleted != nil {
		flowJobMap["workitems_completed"] = flowJobItem.WorkitemsCompleted
	}
	if flowJobItem.WorkitemsPending != nil {
		flowJobMap["workitems_pending"] = flowJobItem.WorkitemsPending
	}
	if flowJobItem.WorkitemsFailed != nil {
		flowJobMap["workitems_failed"] = flowJobItem.WorkitemsFailed
	}
	if flowJobItem.Workitems != nil {
		workitemsList := []map[string]interface{}{}
		for _, workitemsItem := range flowJobItem.Workitems {
			workitemsList = append(workitemsList, dataSourceJobFlowJobWorkitemsToMap(workitemsItem))
		}
		flowJobMap["workitems"] = workitemsList
	}

	return flowJobMap
}

func dataSourceJobFlowJobWorkitemsToMap(workitemsItem schematicsv1.JobLogSummaryWorkitems) (workitemsMap map[string]interface{}) {
	workitemsMap = map[string]interface{}{}

	if workitemsItem.WorkspaceID != nil {
		workitemsMap["workspace_id"] = workitemsItem.WorkspaceID
	}
	if workitemsItem.JobID != nil {
		workitemsMap["job_id"] = workitemsItem.JobID
	}
	if workitemsItem.ResourcesAdd != nil {
		workitemsMap["resources_add"] = workitemsItem.ResourcesAdd
	}
	if workitemsItem.ResourcesModify != nil {
		workitemsMap["resources_modify"] = workitemsItem.ResourcesModify
	}
	if workitemsItem.ResourcesDestroy != nil {
		workitemsMap["resources_destroy"] = workitemsItem.ResourcesDestroy
	}
	if workitemsItem.LogURL != nil {
		workitemsMap["log_url"] = workitemsItem.LogURL
	}

	return workitemsMap
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

func dataSourceJobLogSummarySystemJobToMap(systemJobItem schematicsv1.JobLogSummarySystemJob) (systemJobMap map[string]interface{}) {
	systemJobMap = map[string]interface{}{}

	if systemJobItem.TargetCount != nil {
		systemJobMap["target_count"] = systemJobItem.TargetCount
	}
	if systemJobItem.Success != nil {
		systemJobMap["success"] = systemJobItem.Success
	}
	if systemJobItem.Failed != nil {
		systemJobMap["failed"] = systemJobItem.Failed
	}

	return systemJobMap
}
