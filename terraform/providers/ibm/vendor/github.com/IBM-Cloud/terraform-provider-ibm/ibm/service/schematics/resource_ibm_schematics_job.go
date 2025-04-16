// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func ResourceIBMSchematicsJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSchematicsJobCreate,
		ReadContext:   resourceIBMSchematicsJobRead,
		UpdateContext: resourceIBMSchematicsJobUpdate,
		DeleteContext: resourceIBMSchematicsJobDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"command_object": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_schematics_job", "command_object"),
				Description:  "Name of the Schematics automation resource.",
			},
			"command_object_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Job command object id (workspace-id, action-id).",
			},
			"command_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_schematics_job", "command_name"),
				Description:  "Schematics job command name.",
			},
			"command_parameter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Schematics job command parameter (playbook-name).",
			},
			"command_options": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Command line options for the command.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"job_inputs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Job inputs used by Action or Workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the variable.",
									},
									"aliases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of aliases for the variable name.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the meta data.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"min_value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Regex for the variable value.",
									},
									"position": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"job_env_settings": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Environment variables used by the Job while performing Action or Workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value for the variable or reference to the value.",
						},
						"metadata": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "User editable metadata for the variables.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the variable.",
									},
									"aliases": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of aliases for the variable name.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Description of the meta data.",
									},
									"default_value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Default value for the variable, if the override value is not specified.",
									},
									"secure": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable secure or sensitive ?.",
									},
									"immutable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Is the variable readonly ?.",
									},
									"hidden": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "If true, the variable will not be displayed on UI or CLI.",
									},
									"options": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"min_value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum value of the variable. Applicable for integer type.",
									},
									"max_value": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum value of the variable. Applicable for integer type.",
									},
									"min_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Minimum length of the variable value. Applicable for string type.",
									},
									"max_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Maximum length of the variable value. Applicable for string type.",
									},
									"matches": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Regex for the variable value.",
									},
									"position": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Relative position of this variable in a list.",
									},
									"group_by": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Display name of the group this variable belongs to.",
									},
									"source": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Source of this meta-data.",
									},
								},
							},
						},
						"link": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Reference link to the variable value By default the expression will point to self.value.",
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "User defined tags, while running the job.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"location": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_schematics_job", "location"),
				Description:  "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Job Status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_job_status": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Workspace Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Workspace name.",
									},
									"status_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of Jobs.",
									},
									"status_message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Workspace job status message (eg. App1_Setup_Pending, for a 'Setup' flow in the 'App1' Workspace).",
									},
									"flow_status": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Environment Flow JOB Status.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"flow_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "flow id.",
												},
												"flow_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "flow name.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Flow Job status message - to be displayed along with the status_code;.",
												},
												"workitems": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Environment's individual workItem status details;.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"workspace_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Workspace id.",
															},
															"workspace_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "workspace name.",
															},
															"job_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "workspace job id.",
															},
															"status_code": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Status of Jobs.",
															},
															"status_message": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "workitem job status message;.",
															},
															"updated_at": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "workitem job status updation timestamp.",
															},
														},
													},
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"template_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Workspace Flow Template job status.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"template_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template Id.",
												},
												"template_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template name.",
												},
												"flow_index": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Index of the template in the Flow.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template job status message (eg. VPCt1_Apply_Pending, for a 'VPCt1' Template).",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"action_job_status": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Action Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action name.",
									},
									"status_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of Jobs.",
									},
									"status_message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Action Job status message - to be displayed along with the action_status_code.",
									},
									"bastion_status_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of Resources.",
									},
									"bastion_status_message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Bastion status message - to be displayed along with the bastion_status_code;.",
									},
									"targets_status_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of Resources.",
									},
									"targets_status_message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Aggregated status message for all target resources,  to be displayed along with the targets_status_code;.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"system_job_status": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "System Job Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"system_status_message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "System job message.",
									},
									"system_status_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of Jobs.",
									},
									"schematics_resource_status": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "job staus for each schematics resource.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"status_code": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "system job status message.",
												},
												"schematics_resource_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "id for each resource which is targeted as a part of system job.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"flow_job_status": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Environment Flow JOB Status.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flow_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "flow id.",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "flow name.",
									},
									"status_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Status of Jobs.",
									},
									"status_message": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow Job status message - to be displayed along with the status_code;.",
									},
									"workitems": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Environment's individual workItem status details;.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Workspace id.",
												},
												"workspace_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "workspace name.",
												},
												"job_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "workspace job id.",
												},
												"status_code": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Status of Jobs.",
												},
												"status_message": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "workitem job status message;.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "workitem job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
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
			"data": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Job data.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of Job.",
						},
						"workspace_job_data": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Workspace Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workspace_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Workspace name.",
									},
									"flow_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow Id.",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow name.",
									},
									"inputs": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Input variables data used by the Workspace Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"outputs": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Output variables data from the Workspace Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"settings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Environment variables used by all the templates in the Workspace.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"template_data": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Input / output data of the Template in the Workspace Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"template_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template Id.",
												},
												"template_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Template name.",
												},
												"flow_index": {
													Type:        schema.TypeInt,
													Optional:    true,
													Description: "Index of the template in the Flow.",
												},
												"inputs": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Job inputs used by the Templates.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of aliases for the variable name.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"outputs": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Job output from the Templates.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of aliases for the variable name.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"settings": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Environment variables used by the template.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of aliases for the variable name.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"action_job_data": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Action Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow name.",
									},
									"inputs": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Input variables data used by the Action Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"outputs": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Output variables data from the Action Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"settings": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Environment variables used by all the templates in the Action.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Name of the variable.",
												},
												"value": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Value for the variable or reference to the value.",
												},
												"metadata": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "User editable metadata for the variables.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Type of the variable.",
															},
															"aliases": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of aliases for the variable name.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"description": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Description of the meta data.",
															},
															"default_value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Default value for the variable, if the override value is not specified.",
															},
															"secure": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable secure or sensitive ?.",
															},
															"immutable": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "Is the variable readonly ?.",
															},
															"hidden": {
																Type:        schema.TypeBool,
																Optional:    true,
																Description: "If true, the variable will not be displayed on UI or CLI.",
															},
															"options": {
																Type:        schema.TypeList,
																Optional:    true,
																Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																Elem:        &schema.Schema{Type: schema.TypeString},
															},
															"min_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum value of the variable. Applicable for integer type.",
															},
															"max_value": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum value of the variable. Applicable for integer type.",
															},
															"min_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Minimum length of the variable value. Applicable for string type.",
															},
															"max_length": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Maximum length of the variable value. Applicable for string type.",
															},
															"matches": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Regex for the variable value.",
															},
															"position": {
																Type:        schema.TypeInt,
																Optional:    true,
																Description: "Relative position of this variable in a list.",
															},
															"group_by": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Display name of the group this variable belongs to.",
															},
															"source": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Source of this meta-data.",
															},
														},
													},
												},
												"link": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Reference link to the variable value By default the expression will point to self.value.",
												},
											},
										},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
									"inventory_record": {
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "Complete inventory resource details with user inputs and system generated data.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The unique name of your Inventory.  The name can be up to 128 characters long and can include alphanumeric  characters, spaces, dashes, and underscores.",
												},
												"id": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Inventory id.",
												},
												"description": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The description of your Inventory.  The description can be up to 2048 characters long in size.",
												},
												"location": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
												},
												"resource_group": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Resource-group name for the Inventory definition.  By default, Inventory will be created in Default Resource Group.",
												},
												"created_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Inventory creation time.",
												},
												"created_by": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Email address of user who created the Inventory.",
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Inventory updation time.",
												},
												"updated_by": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Email address of user who updated the Inventory.",
												},
												"inventories_ini": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Input inventory of host and host group for the playbook,  in the .ini file format.",
												},
												"resource_queries": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Input resource queries that is used to dynamically generate  the inventory of host and host group for the playbook.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"materialized_inventory": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Materialized inventory details used by the Action Job, in .ini format.",
									},
								},
							},
						},
						"system_job_data": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Controls Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Key ID for which key event is generated.",
									},
									"schematics_resource_id": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "List of the schematics resource id.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"updated_at": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Job status updation timestamp.",
									},
								},
							},
						},
						"flow_job_data": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Flow Job data.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"flow_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow ID.",
									},
									"flow_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Flow Name.",
									},
									"workitems": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Job data used by each workitem Job.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"command_object_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "command object id.",
												},
												"command_object_name": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "command object name.",
												},
												"layers": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "layer name.",
												},
												"source_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Type of source for the Template.",
												},
												"source": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Source of templates, playbooks, or controls.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"source_type": {
																Type:        schema.TypeString,
																Required:    true,
																Description: "Type of source for the Template.",
															},
															"git": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Connection details to Git source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"computed_git_repo_url": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "The Complete URL which is computed by git_repo_url, git_repo_folder and branch.",
																		},
																		"git_repo_url": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "URL to the GIT Repo that can be used to clone the template.",
																		},
																		"git_token": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Personal Access Token to connect to Git URLs.",
																		},
																		"git_repo_folder": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Name of the folder in the Git Repo, that contains the template.",
																		},
																		"git_release": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Name of the release tag, used to fetch the Git Repo.",
																		},
																		"git_branch": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Name of the branch, used to fetch the Git Repo.",
																		},
																	},
																},
															},
															"catalog": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "Connection details to IBM Cloud Catalog source.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"catalog_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "name of the private catalog.",
																		},
																		"offering_name": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Name of the offering in the IBM Catalog.",
																		},
																		"offering_version": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Version string of the offering in the IBM Catalog.",
																		},
																		"offering_kind": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the offering, in the IBM Catalog.",
																		},
																		"offering_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Id of the offering the IBM Catalog.",
																		},
																		"offering_version_id": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Id of the offering version the IBM Catalog.",
																		},
																		"offering_repo_url": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Repo Url of the offering, in the IBM Catalog.",
																		},
																	},
																},
															},
															// "cos_bucket": {
															// 	Type:        schema.TypeList,
															// 	MaxItems:    1,
															// 	Optional:    true,
															// 	Description: "Connection details to a IBM Cloud Object Storage bucket.",
															// 	Elem: &schema.Resource{
															// 		Schema: map[string]*schema.Schema{
															// 			"cos_bucket_url": {
															// 				Type:        schema.TypeString,
															// 				Optional:    true,
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
													Optional:    true,
													Description: "Input variables data for the workItem used in FlowJob.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of aliases for the variable name.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"outputs": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Output variables for the workItem.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of aliases for the variable name.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"settings": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "Environment variables for the workItem.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the variable.",
															},
															"value": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Value for the variable or reference to the value.",
															},
															"metadata": {
																Type:        schema.TypeList,
																MaxItems:    1,
																Optional:    true,
																Description: "User editable metadata for the variables.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Type of the variable.",
																		},
																		"aliases": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of aliases for the variable name.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"description": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Description of the meta data.",
																		},
																		"default_value": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Default value for the variable, if the override value is not specified.",
																		},
																		"secure": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable secure or sensitive ?.",
																		},
																		"immutable": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "Is the variable readonly ?.",
																		},
																		"hidden": {
																			Type:        schema.TypeBool,
																			Optional:    true,
																			Description: "If true, the variable will not be displayed on UI or CLI.",
																		},
																		"options": {
																			Type:        schema.TypeList,
																			Optional:    true,
																			Description: "List of possible values for this variable.  If type is integer or date, then the array of string will be  converted to array of integers or date during runtime.",
																			Elem:        &schema.Schema{Type: schema.TypeString},
																		},
																		"min_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum value of the variable. Applicable for integer type.",
																		},
																		"max_value": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum value of the variable. Applicable for integer type.",
																		},
																		"min_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Minimum length of the variable value. Applicable for string type.",
																		},
																		"max_length": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Maximum length of the variable value. Applicable for string type.",
																		},
																		"matches": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Regex for the variable value.",
																		},
																		"position": {
																			Type:        schema.TypeInt,
																			Optional:    true,
																			Description: "Relative position of this variable in a list.",
																		},
																		"group_by": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Display name of the group this variable belongs to.",
																		},
																		"source": {
																			Type:        schema.TypeString,
																			Optional:    true,
																			Description: "Source of this meta-data.",
																		},
																	},
																},
															},
															"link": {
																Type:        schema.TypeString,
																Optional:    true,
																Computed:    true,
																Description: "Reference link to the variable value By default the expression will point to self.value.",
															},
														},
													},
												},
												"last_job": {
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "Status of the last job executed by the workitem.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"command_object": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Name of the Schematics automation resource.",
															},
															"command_object_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "command object name (workspace_name/action_name).",
															},
															"command_object_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Workitem command object id, maps to workspace_id or action_id.",
															},
															"command_name": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Schematics job command name.",
															},
															"job_id": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Workspace job id.",
															},
															"job_status": {
																Type:        schema.TypeString,
																Optional:    true,
																Description: "Status of Jobs.",
															},
														},
													},
												},
												"updated_at": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "Job status updation timestamp.",
												},
											},
										},
									},
									"updated_at": {
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
			"bastion": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Describes a bastion resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Bastion Name(Unique).",
						},
						"host": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Reference to the Inventory resource definition.",
						},
					},
				},
			},
			"log_summary": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Job log summary record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Workspace Id.",
						},
						"job_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type of Job.",
						},
						"log_start_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Job log start timestamp.",
						},
						"log_analyzed_till": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Job log update timestamp.",
						},
						"elapsed_time": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "Job log elapsed time (log_analyzed_till - log_start_at).",
						},
						"log_errors": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Job log errors.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Error code in the Log.",
									},
									"error_msg": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Summary error message in the log.",
									},
									"error_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Number of occurrence.",
									},
								},
							},
						},
						"repo_download_job": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Repo download Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"scanned_file_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of files scanned.",
									},
									"quarantined_file_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of files quarantined.",
									},
									"detected_filetype": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Detected template or data file type.",
									},
									"inputs_count": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Number of inputs detected.",
									},
									"outputs_count": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Number of outputs detected.",
									},
								},
							},
						},
						"workspace_job": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Workspace Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resources_add": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of resources add.",
									},
									"resources_modify": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of resources modify.",
									},
									"resources_destroy": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of resources destroy.",
									},
								},
							},
						},
						"flow_job": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Flow Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workitems_completed": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of workitems completed successfully.",
									},
									"workitems_pending": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of workitems pending in the flow.",
									},
									"workitems_failed": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "Number of workitems failed.",
									},
									"workitems": {
										Type:     schema.TypeList,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"workspace_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "workspace ID.",
												},
												"job_id": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: "workspace JOB ID.",
												},
												"resources_add": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Computed:    true,
													Description: "Number of resources add.",
												},
												"resources_modify": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Computed:    true,
													Description: "Number of resources modify.",
												},
												"resources_destroy": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Computed:    true,
													Description: "Number of resources destroy.",
												},
												"log_url": {
													Type:        schema.TypeString,
													Optional:    true,
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
							Optional:    true,
							Description: "Flow Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of targets or hosts.",
									},
									"task_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of tasks in playbook.",
									},
									"play_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of plays in playbook.",
									},
									"recap": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Recap records.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"target": {
													Type:        schema.TypeList,
													Optional:    true,
													Description: "List of target or host name.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"ok": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of OK.",
												},
												"changed": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of changed.",
												},
												"failed": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of failed.",
												},
												"skipped": {
													Type:        schema.TypeFloat,
													Optional:    true,
													Description: "Number of skipped.",
												},
												"unreachable": {
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
						"system_job": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "System Job log summary.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"target_count": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Computed:    true,
										Description: "number of targets or hosts.",
									},
									"success": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Number of passed.",
									},
									"failed": {
										Type:        schema.TypeFloat,
										Optional:    true,
										Description: "Number of failed.",
									},
								},
							},
						},
					},
				},
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

func ResourceIBMSchematicsJobValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "command_object",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "action, environment, system, workspace",
		},
		validate.ValidateSchema{
			Identifier:                 "command_name",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "ansible_playbook_check, ansible_playbook_run, create_action, create_cart, create_environment, create_workspace, delete_action, delete_environment, delete_workspace, environment_init, environment_install, environment_uninstall, patch_action, patch_workspace, put_action, put_environment, put_workspace, repository_process, system_key_delete, system_key_disable, system_key_enable, system_key_restore, system_key_rotate, workspace_apply, workspace_destroy, workspace_plan, workspace_refresh",
		},
		validate.ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "eu-de, eu-gb, us-east, us-south",
		})

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_schematics_job", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSchematicsJobCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}

	iamRefreshToken := session.Config.IAMRefreshToken

	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}
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
		var jobInputs []schematicsv1.VariableData
		for _, e := range d.Get("job_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			jobInputsItem := resourceIBMSchematicsJobMapToVariableData(value)
			jobInputs = append(jobInputs, jobInputsItem)
		}
		createJobOptions.SetInputs(jobInputs)
	}
	if _, ok := d.GetOk("job_env_settings"); ok {
		var jobEnvSettings []schematicsv1.VariableData
		for _, e := range d.Get("job_env_settings").([]interface{}) {
			value := e.(map[string]interface{})
			jobEnvSettingsItem := resourceIBMSchematicsJobMapToVariableData(value)
			jobEnvSettings = append(jobEnvSettings, jobEnvSettingsItem)
		}
		createJobOptions.SetSettings(jobEnvSettings)
	}
	if _, ok := d.GetOk("tags"); ok {
		createJobOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
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
			bastion := resourceIBMSchematicsJobMapToBastionResourceDefinition(d.Get("bastion.0").(map[string]interface{}))
			createJobOptions.SetBastion(&bastion)
		}
	}
	if _, ok := d.GetOk("log_summary"); ok {
		logSummaryAttr := d.Get("log_summary").([]interface{})
		if len(logSummaryAttr) > 0 {
			logSummary := resourceIBMSchematicsJobMapToJobLogSummary(d.Get("log_summary.0").(map[string]interface{}))
			createJobOptions.SetLogSummary(&logSummary)
		}
	}

	job, response, err := schematicsClient.CreateJobWithContext(context, createJobOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(*job.ID)

	return resourceIBMSchematicsJobRead(context, d, meta)
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
func resourceIBMSchematicsJobMapToCredentialVariableMetadata(variableMetadataMap map[string]interface{}) schematicsv1.CredentialVariableMetadata {
	variableMetadata := schematicsv1.CredentialVariableMetadata{}

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
	if variableMetadataMap["immutable"] != nil {
		variableMetadata.Immutable = core.BoolPtr(variableMetadataMap["immutable"].(bool))
	}
	if variableMetadataMap["hidden"] != nil {
		variableMetadata.Hidden = core.BoolPtr(variableMetadataMap["hidden"].(bool))
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

	if jobStatusMap["workspace_job_status"] != nil {
		// TODO: handle WorkspaceJobStatus of type JobStatusWorkspace -- not primitive type, not list
	}
	if jobStatusMap["action_job_status"] != nil {
		actionJobStatus := resourceIBMSchematicsJobMapToJobStatusAction(jobStatusMap["action_job_status"].([]interface{})[0].(map[string]interface{}))
		jobStatus.ActionJobStatus = &actionJobStatus
	}

	return jobStatus
}

func resourceIBMSchematicsJobMapToJobStatusWorkspace(jobStatusWorkspaceMap map[string]interface{}) schematicsv1.JobStatusWorkspace {
	jobStatusWorkspace := schematicsv1.JobStatusWorkspace{}

	if jobStatusWorkspaceMap["workspace_name"] != nil {
		jobStatusWorkspace.WorkspaceName = core.StringPtr(jobStatusWorkspaceMap["workspace_name"].(string))
	}
	if jobStatusWorkspaceMap["status_code"] != nil {
		jobStatusWorkspace.StatusCode = core.StringPtr(jobStatusWorkspaceMap["status_code"].(string))
	}
	if jobStatusWorkspaceMap["status_message"] != nil {
		jobStatusWorkspace.StatusMessage = core.StringPtr(jobStatusWorkspaceMap["status_message"].(string))
	}
	if jobStatusWorkspaceMap["flow_status"] != nil {
		// TODO: handle FlowStatus of type JobStatusFlow -- not primitive type, not list
	}
	if jobStatusWorkspaceMap["template_status"] != nil {
		templateStatus := []schematicsv1.JobStatusTemplate{}
		for _, templateStatusItem := range jobStatusWorkspaceMap["template_status"].([]interface{}) {
			templateStatusItemModel := resourceIBMSchematicsJobMapToJobStatusTemplate(templateStatusItem.(map[string]interface{}))
			templateStatus = append(templateStatus, templateStatusItemModel)
		}
		jobStatusWorkspace.TemplateStatus = templateStatus
	}
	if jobStatusWorkspaceMap["updated_at"] != nil {

	}

	return jobStatusWorkspace
}

func resourceIBMSchematicsJobMapToJobStatusFlow(jobStatusFlowMap map[string]interface{}) schematicsv1.JobStatusFlow {
	jobStatusFlow := schematicsv1.JobStatusFlow{}

	if jobStatusFlowMap["flow_id"] != nil {
		jobStatusFlow.FlowID = core.StringPtr(jobStatusFlowMap["flow_id"].(string))
	}
	if jobStatusFlowMap["flow_name"] != nil {
		jobStatusFlow.FlowName = core.StringPtr(jobStatusFlowMap["flow_name"].(string))
	}
	if jobStatusFlowMap["status_code"] != nil {
		jobStatusFlow.StatusCode = core.StringPtr(jobStatusFlowMap["status_code"].(string))
	}
	if jobStatusFlowMap["status_message"] != nil {
		jobStatusFlow.StatusMessage = core.StringPtr(jobStatusFlowMap["status_message"].(string))
	}
	if jobStatusFlowMap["workitems"] != nil {
		workitems := []schematicsv1.JobStatusWorkitem{}
		for _, workitemsItem := range jobStatusFlowMap["workitems"].([]interface{}) {
			workitemsItemModel := resourceIBMSchematicsJobMapToJobStatusWorkitem(workitemsItem.(map[string]interface{}))
			workitems = append(workitems, workitemsItemModel)
		}
		jobStatusFlow.Workitems = workitems
	}
	if jobStatusFlowMap["updated_at"] != nil {

	}

	return jobStatusFlow
}

func resourceIBMSchematicsJobMapToJobStatusWorkitem(jobStatusWorkitemMap map[string]interface{}) schematicsv1.JobStatusWorkitem {
	jobStatusWorkitem := schematicsv1.JobStatusWorkitem{}

	if jobStatusWorkitemMap["workspace_id"] != nil {
		jobStatusWorkitem.WorkspaceID = core.StringPtr(jobStatusWorkitemMap["workspace_id"].(string))
	}
	if jobStatusWorkitemMap["workspace_name"] != nil {
		jobStatusWorkitem.WorkspaceName = core.StringPtr(jobStatusWorkitemMap["workspace_name"].(string))
	}
	if jobStatusWorkitemMap["job_id"] != nil {
		jobStatusWorkitem.JobID = core.StringPtr(jobStatusWorkitemMap["job_id"].(string))
	}
	if jobStatusWorkitemMap["status_code"] != nil {
		jobStatusWorkitem.StatusCode = core.StringPtr(jobStatusWorkitemMap["status_code"].(string))
	}
	if jobStatusWorkitemMap["status_message"] != nil {
		jobStatusWorkitem.StatusMessage = core.StringPtr(jobStatusWorkitemMap["status_message"].(string))
	}
	if jobStatusWorkitemMap["updated_at"] != nil {

	}

	return jobStatusWorkitem
}

func resourceIBMSchematicsJobMapToJobStatusTemplate(jobStatusTemplateMap map[string]interface{}) schematicsv1.JobStatusTemplate {
	jobStatusTemplate := schematicsv1.JobStatusTemplate{}

	if jobStatusTemplateMap["template_id"] != nil {
		jobStatusTemplate.TemplateID = core.StringPtr(jobStatusTemplateMap["template_id"].(string))
	}
	if jobStatusTemplateMap["template_name"] != nil {
		jobStatusTemplate.TemplateName = core.StringPtr(jobStatusTemplateMap["template_name"].(string))
	}
	if jobStatusTemplateMap["flow_index"] != nil {
		jobStatusTemplate.FlowIndex = core.Int64Ptr(int64(jobStatusTemplateMap["flow_index"].(int)))
	}
	if jobStatusTemplateMap["status_code"] != nil {
		jobStatusTemplate.StatusCode = core.StringPtr(jobStatusTemplateMap["status_code"].(string))
	}
	if jobStatusTemplateMap["status_message"] != nil {
		jobStatusTemplate.StatusMessage = core.StringPtr(jobStatusTemplateMap["status_message"].(string))
	}
	if jobStatusTemplateMap["updated_at"] != nil {

	}

	return jobStatusTemplate
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

func resourceIBMSchematicsJobMapToJobStatusSystem(jobStatusSystemMap map[string]interface{}) schematicsv1.JobStatusSystem {
	jobStatusSystem := schematicsv1.JobStatusSystem{}

	if jobStatusSystemMap["system_status_message"] != nil {
		jobStatusSystem.SystemStatusMessage = core.StringPtr(jobStatusSystemMap["system_status_message"].(string))
	}
	if jobStatusSystemMap["system_status_code"] != nil {
		jobStatusSystem.SystemStatusCode = core.StringPtr(jobStatusSystemMap["system_status_code"].(string))
	}
	if jobStatusSystemMap["schematics_resource_status"] != nil {
		schematicsResourceStatus := []schematicsv1.JobStatusSchematicsResources{}
		for _, schematicsResourceStatusItem := range jobStatusSystemMap["schematics_resource_status"].([]interface{}) {
			schematicsResourceStatusItemModel := resourceIBMSchematicsJobMapToJobStatusSchematicsResources(schematicsResourceStatusItem.(map[string]interface{}))
			schematicsResourceStatus = append(schematicsResourceStatus, schematicsResourceStatusItemModel)
		}
		jobStatusSystem.SchematicsResourceStatus = schematicsResourceStatus
	}
	if jobStatusSystemMap["updated_at"] != nil {

	}

	return jobStatusSystem
}

func resourceIBMSchematicsJobMapToJobStatusSchematicsResources(jobStatusSchematicsResourcesMap map[string]interface{}) schematicsv1.JobStatusSchematicsResources {
	jobStatusSchematicsResources := schematicsv1.JobStatusSchematicsResources{}

	if jobStatusSchematicsResourcesMap["status_code"] != nil {
		jobStatusSchematicsResources.StatusCode = core.StringPtr(jobStatusSchematicsResourcesMap["status_code"].(string))
	}
	if jobStatusSchematicsResourcesMap["status_message"] != nil {
		jobStatusSchematicsResources.StatusMessage = core.StringPtr(jobStatusSchematicsResourcesMap["status_message"].(string))
	}
	if jobStatusSchematicsResourcesMap["schematics_resource_id"] != nil {
		jobStatusSchematicsResources.SchematicsResourceID = core.StringPtr(jobStatusSchematicsResourcesMap["schematics_resource_id"].(string))
	}
	if jobStatusSchematicsResourcesMap["updated_at"] != nil {

	}

	return jobStatusSchematicsResources
}

func resourceIBMSchematicsJobMapToJobData(jobDataMap map[string]interface{}) schematicsv1.JobData {
	jobData := schematicsv1.JobData{}

	jobData.JobType = core.StringPtr(jobDataMap["job_type"].(string))
	if jobDataMap["workspace_job_data"] != nil {
		workspaceJobData := resourceIBMSchematicsJobMapToJobDataWorkspace(jobDataMap["workspace_job_data"].([]interface{})[0].(map[string]interface{}))
		jobData.WorkspaceJobData = &workspaceJobData
	}
	if jobDataMap["action_job_data"] != nil {
		actionJobData := resourceIBMSchematicsJobMapToJobDataAction(jobDataMap["action_job_data"].([]interface{})[0].(map[string]interface{}))
		jobData.ActionJobData = &actionJobData
	}
	if jobDataMap["system_job_data"] != nil {
		systemJobData := resourceIBMSchematicsJobMapToJobDataSystem(jobDataMap["system_job_data"].([]interface{})[0].(map[string]interface{}))
		jobData.SystemJobData = &systemJobData
	}
	if jobDataMap["flow_job_data"] != nil {
		flowJobData := resourceIBMSchematicsJobMapToJobDataFlow(jobDataMap["flow_job_data"].([]interface{})[0].(map[string]interface{}))
		jobData.FlowJobData = &flowJobData
	}

	return jobData
}

func resourceIBMSchematicsJobMapToJobDataWorkspace(jobDataWorkspaceMap map[string]interface{}) schematicsv1.JobDataWorkspace {
	jobDataWorkspace := schematicsv1.JobDataWorkspace{}

	if jobDataWorkspaceMap["workspace_name"] != nil {
		jobDataWorkspace.WorkspaceName = core.StringPtr(jobDataWorkspaceMap["workspace_name"].(string))
	}
	if jobDataWorkspaceMap["flow_id"] != nil {
		jobDataWorkspace.FlowID = core.StringPtr(jobDataWorkspaceMap["flow_id"].(string))
	}
	if jobDataWorkspaceMap["flow_name"] != nil {
		jobDataWorkspace.FlowName = core.StringPtr(jobDataWorkspaceMap["flow_name"].(string))
	}
	if jobDataWorkspaceMap["inputs"] != nil {
		inputs := []schematicsv1.VariableData{}
		for _, inputsItem := range jobDataWorkspaceMap["inputs"].([]interface{}) {
			inputsItemModel := resourceIBMSchematicsJobMapToVariableData(inputsItem.(map[string]interface{}))
			inputs = append(inputs, inputsItemModel)
		}
		jobDataWorkspace.Inputs = inputs
	}
	if jobDataWorkspaceMap["outputs"] != nil {
		outputs := []schematicsv1.VariableData{}
		for _, outputsItem := range jobDataWorkspaceMap["outputs"].([]interface{}) {
			outputsItemModel := resourceIBMSchematicsJobMapToVariableData(outputsItem.(map[string]interface{}))
			outputs = append(outputs, outputsItemModel)
		}
		jobDataWorkspace.Outputs = outputs
	}
	if jobDataWorkspaceMap["settings"] != nil {
		settings := []schematicsv1.VariableData{}
		for _, settingsItem := range jobDataWorkspaceMap["settings"].([]interface{}) {
			settingsItemModel := resourceIBMSchematicsJobMapToVariableData(settingsItem.(map[string]interface{}))
			settings = append(settings, settingsItemModel)
		}
		jobDataWorkspace.Settings = settings
	}
	if jobDataWorkspaceMap["template_data"] != nil {
		templateData := []schematicsv1.JobDataTemplate{}
		for _, templateDataItem := range jobDataWorkspaceMap["template_data"].([]interface{}) {
			templateDataItemModel := resourceIBMSchematicsJobMapToJobDataTemplate(templateDataItem.(map[string]interface{}))
			templateData = append(templateData, templateDataItemModel)
		}
		jobDataWorkspace.TemplateData = templateData
	}
	if jobDataWorkspaceMap["updated_at"] != nil {

	}

	return jobDataWorkspace
}

func resourceIBMSchematicsJobMapToJobDataTemplate(jobDataTemplateMap map[string]interface{}) schematicsv1.JobDataTemplate {
	jobDataTemplate := schematicsv1.JobDataTemplate{}

	if jobDataTemplateMap["template_id"] != nil {
		jobDataTemplate.TemplateID = core.StringPtr(jobDataTemplateMap["template_id"].(string))
	}
	if jobDataTemplateMap["template_name"] != nil {
		jobDataTemplate.TemplateName = core.StringPtr(jobDataTemplateMap["template_name"].(string))
	}
	if jobDataTemplateMap["flow_index"] != nil {
		jobDataTemplate.FlowIndex = core.Int64Ptr(int64(jobDataTemplateMap["flow_index"].(int)))
	}
	if jobDataTemplateMap["inputs"] != nil {
		inputs := []schematicsv1.VariableData{}
		for _, inputsItem := range jobDataTemplateMap["inputs"].([]interface{}) {
			inputsItemModel := resourceIBMSchematicsJobMapToVariableData(inputsItem.(map[string]interface{}))
			inputs = append(inputs, inputsItemModel)
		}
		jobDataTemplate.Inputs = inputs
	}
	if jobDataTemplateMap["outputs"] != nil {
		outputs := []schematicsv1.VariableData{}
		for _, outputsItem := range jobDataTemplateMap["outputs"].([]interface{}) {
			outputsItemModel := resourceIBMSchematicsJobMapToVariableData(outputsItem.(map[string]interface{}))
			outputs = append(outputs, outputsItemModel)
		}
		jobDataTemplate.Outputs = outputs
	}
	if jobDataTemplateMap["settings"] != nil {
		settings := []schematicsv1.VariableData{}
		for _, settingsItem := range jobDataTemplateMap["settings"].([]interface{}) {
			settingsItemModel := resourceIBMSchematicsJobMapToVariableData(settingsItem.(map[string]interface{}))
			settings = append(settings, settingsItemModel)
		}
		jobDataTemplate.Settings = settings
	}
	if jobDataTemplateMap["updated_at"] != nil {

	}

	return jobDataTemplate
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
	if jobDataActionMap["inventory_record"] != nil {
		// TODO: handle InventoryRecord of type InventoryResourceRecord -- not primitive type, not list
	}
	if jobDataActionMap["materialized_inventory"] != nil {
		jobDataAction.MaterializedInventory = core.StringPtr(jobDataActionMap["materialized_inventory"].(string))
	}

	return jobDataAction
}

func resourceIBMSchematicsJobMapToInventoryResourceRecord(inventoryResourceRecordMap map[string]interface{}) schematicsv1.InventoryResourceRecord {
	inventoryResourceRecord := schematicsv1.InventoryResourceRecord{}

	if inventoryResourceRecordMap["name"] != nil {
		inventoryResourceRecord.Name = core.StringPtr(inventoryResourceRecordMap["name"].(string))
	}
	if inventoryResourceRecordMap["id"] != nil {
		inventoryResourceRecord.ID = core.StringPtr(inventoryResourceRecordMap["id"].(string))
	}
	if inventoryResourceRecordMap["description"] != nil {
		inventoryResourceRecord.Description = core.StringPtr(inventoryResourceRecordMap["description"].(string))
	}
	if inventoryResourceRecordMap["location"] != nil {
		inventoryResourceRecord.Location = core.StringPtr(inventoryResourceRecordMap["location"].(string))
	}
	if inventoryResourceRecordMap["resource_group"] != nil {
		inventoryResourceRecord.ResourceGroup = core.StringPtr(inventoryResourceRecordMap["resource_group"].(string))
	}
	if inventoryResourceRecordMap["created_at"] != nil {

	}
	if inventoryResourceRecordMap["created_by"] != nil {
		inventoryResourceRecord.CreatedBy = core.StringPtr(inventoryResourceRecordMap["created_by"].(string))
	}
	if inventoryResourceRecordMap["updated_at"] != nil {

	}
	if inventoryResourceRecordMap["updated_by"] != nil {
		inventoryResourceRecord.UpdatedBy = core.StringPtr(inventoryResourceRecordMap["updated_by"].(string))
	}
	if inventoryResourceRecordMap["inventories_ini"] != nil {
		inventoryResourceRecord.InventoriesIni = core.StringPtr(inventoryResourceRecordMap["inventories_ini"].(string))
	}
	if inventoryResourceRecordMap["resource_queries"] != nil {
		resourceQueries := []string{}
		for _, resourceQueriesItem := range inventoryResourceRecordMap["resource_queries"].([]interface{}) {
			resourceQueries = append(resourceQueries, resourceQueriesItem.(string))
		}
		inventoryResourceRecord.ResourceQueries = resourceQueries
	}

	return inventoryResourceRecord
}

func resourceIBMSchematicsJobMapToJobDataSystem(jobDataSystemMap map[string]interface{}) schematicsv1.JobDataSystem {
	jobDataSystem := schematicsv1.JobDataSystem{}

	if jobDataSystemMap["key_id"] != nil {
		jobDataSystem.KeyID = core.StringPtr(jobDataSystemMap["key_id"].(string))
	}
	if jobDataSystemMap["schematics_resource_id"] != nil {
		schematicsResourceID := []string{}
		for _, schematicsResourceIDItem := range jobDataSystemMap["schematics_resource_id"].([]interface{}) {
			schematicsResourceID = append(schematicsResourceID, schematicsResourceIDItem.(string))
		}
		jobDataSystem.SchematicsResourceID = schematicsResourceID
	}
	if jobDataSystemMap["updated_at"] != nil {

	}

	return jobDataSystem
}

func resourceIBMSchematicsJobMapToJobDataFlow(jobDataFlowMap map[string]interface{}) schematicsv1.JobDataFlow {
	jobDataFlow := schematicsv1.JobDataFlow{}

	if jobDataFlowMap["flow_id"] != nil {
		jobDataFlow.FlowID = core.StringPtr(jobDataFlowMap["flow_id"].(string))
	}
	if jobDataFlowMap["flow_name"] != nil {
		jobDataFlow.FlowName = core.StringPtr(jobDataFlowMap["flow_name"].(string))
	}
	if jobDataFlowMap["workitems"] != nil {
		workitems := []schematicsv1.JobDataWorkItem{}
		for _, workitemsItem := range jobDataFlowMap["workitems"].([]interface{}) {
			workitemsItemModel := resourceIBMSchematicsJobMapToJobDataWorkItem(workitemsItem.(map[string]interface{}))
			workitems = append(workitems, workitemsItemModel)
		}
		jobDataFlow.Workitems = workitems
	}
	if jobDataFlowMap["updated_at"] != nil {

	}

	return jobDataFlow
}

func resourceIBMSchematicsJobMapToJobDataWorkItem(jobDataWorkItemMap map[string]interface{}) schematicsv1.JobDataWorkItem {
	jobDataWorkItem := schematicsv1.JobDataWorkItem{}

	if jobDataWorkItemMap["command_object_id"] != nil {
		jobDataWorkItem.CommandObjectID = core.StringPtr(jobDataWorkItemMap["command_object_id"].(string))
	}
	if jobDataWorkItemMap["command_object_name"] != nil {
		jobDataWorkItem.CommandObjectName = core.StringPtr(jobDataWorkItemMap["command_object_name"].(string))
	}
	if jobDataWorkItemMap["layers"] != nil {
		jobDataWorkItem.Layers = core.StringPtr(jobDataWorkItemMap["layers"].(string))
	}
	if jobDataWorkItemMap["source_type"] != nil {
		jobDataWorkItem.SourceType = core.StringPtr(jobDataWorkItemMap["source_type"].(string))
	}
	if jobDataWorkItemMap["source"] != nil {
		// TODO: handle Source of type ExternalSource -- not primitive type, not list
	}
	if jobDataWorkItemMap["inputs"] != nil {
		inputs := []schematicsv1.VariableData{}
		for _, inputsItem := range jobDataWorkItemMap["inputs"].([]interface{}) {
			inputsItemModel := resourceIBMSchematicsJobMapToVariableData(inputsItem.(map[string]interface{}))
			inputs = append(inputs, inputsItemModel)
		}
		jobDataWorkItem.Inputs = inputs
	}
	if jobDataWorkItemMap["outputs"] != nil {
		outputs := []schematicsv1.VariableData{}
		for _, outputsItem := range jobDataWorkItemMap["outputs"].([]interface{}) {
			outputsItemModel := resourceIBMSchematicsJobMapToVariableData(outputsItem.(map[string]interface{}))
			outputs = append(outputs, outputsItemModel)
		}
		jobDataWorkItem.Outputs = outputs
	}
	if jobDataWorkItemMap["settings"] != nil {
		settings := []schematicsv1.VariableData{}
		for _, settingsItem := range jobDataWorkItemMap["settings"].([]interface{}) {
			settingsItemModel := resourceIBMSchematicsJobMapToVariableData(settingsItem.(map[string]interface{}))
			settings = append(settings, settingsItemModel)
		}
		jobDataWorkItem.Settings = settings
	}
	if jobDataWorkItemMap["updated_at"] != nil {

	}

	return jobDataWorkItem
}

func resourceIBMSchematicsJobMapToExternalSource(externalSourceMap map[string]interface{}) schematicsv1.ExternalSource {
	externalSource := schematicsv1.ExternalSource{}

	externalSource.SourceType = core.StringPtr(externalSourceMap["source_type"].(string))
	if externalSourceMap["git"] != nil {
		// TODO: handle Git of type ExternalSourceGit -- not primitive type, not list
	}
	if externalSourceMap["catalog"] != nil {
		// TODO: handle Catalog of type ExternalSourceCatalog -- not primitive type, not list
	}

	return externalSource
}

func resourceIBMSchematicsJobMapToExternalSourceGit(externalSourceGitMap map[string]interface{}) schematicsv1.GitSource {
	externalSourceGit := schematicsv1.GitSource{}

	if externalSourceGitMap["computed_git_repo_url"] != nil {
		externalSourceGit.ComputedGitRepoURL = core.StringPtr(externalSourceGitMap["computed_git_repo_url"].(string))
	}
	if externalSourceGitMap["git_repo_url"] != nil {
		externalSourceGit.GitRepoURL = core.StringPtr(externalSourceGitMap["git_repo_url"].(string))
	}
	if externalSourceGitMap["git_token"] != nil {
		externalSourceGit.GitToken = core.StringPtr(externalSourceGitMap["git_token"].(string))
	}
	if externalSourceGitMap["git_repo_folder"] != nil {
		externalSourceGit.GitRepoFolder = core.StringPtr(externalSourceGitMap["git_repo_folder"].(string))
	}
	if externalSourceGitMap["git_release"] != nil {
		externalSourceGit.GitRelease = core.StringPtr(externalSourceGitMap["git_release"].(string))
	}
	if externalSourceGitMap["git_branch"] != nil {
		externalSourceGit.GitBranch = core.StringPtr(externalSourceGitMap["git_branch"].(string))
	}

	return externalSourceGit
}

func resourceIBMSchematicsJobMapToExternalSourceCatalog(externalSourceCatalogMap map[string]interface{}) schematicsv1.CatalogSource {
	externalSourceCatalog := schematicsv1.CatalogSource{}

	if externalSourceCatalogMap["catalog_name"] != nil {
		externalSourceCatalog.CatalogName = core.StringPtr(externalSourceCatalogMap["catalog_name"].(string))
	}
	if externalSourceCatalogMap["offering_name"] != nil {
		externalSourceCatalog.OfferingName = core.StringPtr(externalSourceCatalogMap["offering_name"].(string))
	}
	if externalSourceCatalogMap["offering_version"] != nil {
		externalSourceCatalog.OfferingVersion = core.StringPtr(externalSourceCatalogMap["offering_version"].(string))
	}
	if externalSourceCatalogMap["offering_kind"] != nil {
		externalSourceCatalog.OfferingKind = core.StringPtr(externalSourceCatalogMap["offering_kind"].(string))
	}
	if externalSourceCatalogMap["offering_id"] != nil {
		externalSourceCatalog.OfferingID = core.StringPtr(externalSourceCatalogMap["offering_id"].(string))
	}
	if externalSourceCatalogMap["offering_version_id"] != nil {
		externalSourceCatalog.OfferingVersionID = core.StringPtr(externalSourceCatalogMap["offering_version_id"].(string))
	}
	if externalSourceCatalogMap["offering_repo_url"] != nil {
		externalSourceCatalog.OfferingRepoURL = core.StringPtr(externalSourceCatalogMap["offering_repo_url"].(string))
	}

	return externalSourceCatalog
}

// func resourceIBMSchematicsJobMapToExternalSourceCosBucket(externalSourceCosBucketMap map[string]interface{}) schematicsv1.ExternalSourceCosBucket {
// 	externalSourceCosBucket := schematicsv1.ExternalSourceCosBucket{}

// 	if externalSourceCosBucketMap["cos_bucket_url"] != nil {
// 		externalSourceCosBucket.CosBucketURL = core.StringPtr(externalSourceCosBucketMap["cos_bucket_url"].(string))
// 	}

// 	return externalSourceCosBucket
// }

func resourceIBMSchematicsJobMapToJobDataWorkItemLastJob(jobDataWorkItemLastJobMap map[string]interface{}) schematicsv1.JobDataWorkItemLastJob {
	jobDataWorkItemLastJob := schematicsv1.JobDataWorkItemLastJob{}

	if jobDataWorkItemLastJobMap["command_object"] != nil {
		jobDataWorkItemLastJob.CommandObject = core.StringPtr(jobDataWorkItemLastJobMap["command_object"].(string))
	}
	if jobDataWorkItemLastJobMap["command_object_name"] != nil {
		jobDataWorkItemLastJob.CommandObjectName = core.StringPtr(jobDataWorkItemLastJobMap["command_object_name"].(string))
	}
	if jobDataWorkItemLastJobMap["command_object_id"] != nil {
		jobDataWorkItemLastJob.CommandObjectID = core.StringPtr(jobDataWorkItemLastJobMap["command_object_id"].(string))
	}
	if jobDataWorkItemLastJobMap["command_name"] != nil {
		jobDataWorkItemLastJob.CommandName = core.StringPtr(jobDataWorkItemLastJobMap["command_name"].(string))
	}
	if jobDataWorkItemLastJobMap["job_id"] != nil {
		jobDataWorkItemLastJob.JobID = core.StringPtr(jobDataWorkItemLastJobMap["job_id"].(string))
	}
	if jobDataWorkItemLastJobMap["job_status"] != nil {
		jobDataWorkItemLastJob.JobStatus = core.StringPtr(jobDataWorkItemLastJobMap["job_status"].(string))
	}

	return jobDataWorkItemLastJob
}

func resourceIBMSchematicsJobMapToBastionResourceDefinition(bastionResourceDefinitionMap map[string]interface{}) schematicsv1.BastionResourceDefinition {
	bastionResourceDefinition := schematicsv1.BastionResourceDefinition{}

	if bastionResourceDefinitionMap["name"] != nil {
		bastionResourceDefinition.Name = core.StringPtr(bastionResourceDefinitionMap["name"].(string))
	}
	if bastionResourceDefinitionMap["host"] != nil {
		bastionResourceDefinition.Host = core.StringPtr(bastionResourceDefinitionMap["host"].(string))
	}

	return bastionResourceDefinition
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
		logErrors := []schematicsv1.JobLogSummaryLogErrors{}
		for _, logErrorsItem := range jobLogSummaryMap["log_errors"].([]interface{}) {
			logErrorsItemModel := resourceIBMSchematicsJobMapToJobLogSummaryLogErrors(logErrorsItem.(map[string]interface{}))
			logErrors = append(logErrors, logErrorsItemModel)
		}
		jobLogSummary.LogErrors = logErrors
	}
	if jobLogSummaryMap["repo_download_job"] != nil && jobLogSummaryMap["repo_download_job"].([]interface{})[0] != nil {
		repoDownloadJob := resourceIBMSchematicsJobMapToJobLogSummaryRepoDownloadJob(jobLogSummaryMap["repo_download_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.RepoDownloadJob = &repoDownloadJob
	}
	if jobLogSummaryMap["workspace_job"] != nil && jobLogSummaryMap["workspace_job"].([]interface{})[0] != nil {
		workspaceJob := resourceIBMSchematicsJobMapToJobLogSummaryWorkspaceJob(jobLogSummaryMap["workspace_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.WorkspaceJob = &workspaceJob
	}
	if jobLogSummaryMap["flow_job"] != nil && jobLogSummaryMap["flow_job"].([]interface{})[0] != nil {
		flowJob := resourceIBMSchematicsJobMapToJobLogSummaryFlowJob(jobLogSummaryMap["flow_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.FlowJob = &flowJob
	}
	if jobLogSummaryMap["action_job"] != nil && jobLogSummaryMap["action_job"].([]interface{})[0] != nil {
		actionJob := resourceIBMSchematicsJobMapToJobLogSummaryActionJob(jobLogSummaryMap["action_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.ActionJob = &actionJob
	}
	if jobLogSummaryMap["system_job"] != nil && jobLogSummaryMap["system_job"].([]interface{})[0] != nil {
		systemJob := resourceIBMSchematicsJobMapToJobLogSummarySystemJob(jobLogSummaryMap["system_job"].([]interface{})[0].(map[string]interface{}))
		jobLogSummary.SystemJob = &systemJob
	}

	return jobLogSummary
}

func resourceIBMSchematicsJobMapToJobLogSummaryLogErrors(jobLogSummaryLogErrorsMap map[string]interface{}) schematicsv1.JobLogSummaryLogErrors {
	jobLogSummaryLogErrors := schematicsv1.JobLogSummaryLogErrors{}

	if jobLogSummaryLogErrorsMap["error_code"] != nil {
		jobLogSummaryLogErrors.ErrorCode = core.StringPtr(jobLogSummaryLogErrorsMap["error_code"].(string))
	}
	if jobLogSummaryLogErrorsMap["error_msg"] != nil {
		jobLogSummaryLogErrors.ErrorMsg = core.StringPtr(jobLogSummaryLogErrorsMap["error_msg"].(string))
	}
	if jobLogSummaryLogErrorsMap["error_count"] != nil {
		jobLogSummaryLogErrors.ErrorCount = core.Float64Ptr(jobLogSummaryLogErrorsMap["error_count"].(float64))
	}

	return jobLogSummaryLogErrors
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

func resourceIBMSchematicsJobMapToJobLogSummaryWorkspaceJob(jobLogSummaryWorkspaceJobMap map[string]interface{}) schematicsv1.JobLogSummaryWorkspaceJob {
	jobLogSummaryWorkspaceJob := schematicsv1.JobLogSummaryWorkspaceJob{}

	if jobLogSummaryWorkspaceJobMap["resources_add"] != nil {
		jobLogSummaryWorkspaceJob.ResourcesAdd = core.Float64Ptr(jobLogSummaryWorkspaceJobMap["resources_add"].(float64))
	}
	if jobLogSummaryWorkspaceJobMap["resources_modify"] != nil {
		jobLogSummaryWorkspaceJob.ResourcesModify = core.Float64Ptr(jobLogSummaryWorkspaceJobMap["resources_modify"].(float64))
	}
	if jobLogSummaryWorkspaceJobMap["resources_destroy"] != nil {
		jobLogSummaryWorkspaceJob.ResourcesDestroy = core.Float64Ptr(jobLogSummaryWorkspaceJobMap["resources_destroy"].(float64))
	}

	return jobLogSummaryWorkspaceJob
}

func resourceIBMSchematicsJobMapToJobLogSummaryFlowJob(jobLogSummaryFlowJobMap map[string]interface{}) schematicsv1.JobLogSummaryFlowJob {
	jobLogSummaryFlowJob := schematicsv1.JobLogSummaryFlowJob{}

	if jobLogSummaryFlowJobMap["workitems_completed"] != nil {
		jobLogSummaryFlowJob.WorkitemsCompleted = core.Float64Ptr(jobLogSummaryFlowJobMap["workitems_completed"].(float64))
	}
	if jobLogSummaryFlowJobMap["workitems_pending"] != nil {
		jobLogSummaryFlowJob.WorkitemsPending = core.Float64Ptr(jobLogSummaryFlowJobMap["workitems_pending"].(float64))
	}
	if jobLogSummaryFlowJobMap["workitems_failed"] != nil {
		jobLogSummaryFlowJob.WorkitemsFailed = core.Float64Ptr(jobLogSummaryFlowJobMap["workitems_failed"].(float64))
	}
	if jobLogSummaryFlowJobMap["workitems"] != nil {
		workitems := []schematicsv1.JobLogSummaryWorkitems{}
		for _, workitemsItem := range jobLogSummaryFlowJobMap["workitems"].([]interface{}) {
			workitemsItemModel := resourceIBMSchematicsJobMapToJobLogSummaryWorkitems(workitemsItem.(map[string]interface{}))
			workitems = append(workitems, workitemsItemModel)
		}
		jobLogSummaryFlowJob.Workitems = workitems
	}

	return jobLogSummaryFlowJob
}

func resourceIBMSchematicsJobMapToJobLogSummaryWorkitems(jobLogSummaryWorkitemsMap map[string]interface{}) schematicsv1.JobLogSummaryWorkitems {
	jobLogSummaryWorkitems := schematicsv1.JobLogSummaryWorkitems{}

	if jobLogSummaryWorkitemsMap["workspace_id"] != nil {
		jobLogSummaryWorkitems.WorkspaceID = core.StringPtr(jobLogSummaryWorkitemsMap["workspace_id"].(string))
	}
	if jobLogSummaryWorkitemsMap["job_id"] != nil {
		jobLogSummaryWorkitems.JobID = core.StringPtr(jobLogSummaryWorkitemsMap["job_id"].(string))
	}
	if jobLogSummaryWorkitemsMap["resources_add"] != nil {
		jobLogSummaryWorkitems.ResourcesAdd = core.Float64Ptr(jobLogSummaryWorkitemsMap["resources_add"].(float64))
	}
	if jobLogSummaryWorkitemsMap["resources_modify"] != nil {
		jobLogSummaryWorkitems.ResourcesModify = core.Float64Ptr(jobLogSummaryWorkitemsMap["resources_modify"].(float64))
	}
	if jobLogSummaryWorkitemsMap["resources_destroy"] != nil {
		jobLogSummaryWorkitems.ResourcesDestroy = core.Float64Ptr(jobLogSummaryWorkitemsMap["resources_destroy"].(float64))
	}
	if jobLogSummaryWorkitemsMap["log_url"] != nil {
		jobLogSummaryWorkitems.LogURL = core.StringPtr(jobLogSummaryWorkitemsMap["log_url"].(string))
	}

	return jobLogSummaryWorkitems
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
	if jobLogSummaryActionJobMap["recap"] != nil && jobLogSummaryActionJobMap["recap"].([]interface{})[0] != nil {
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

func resourceIBMSchematicsJobMapToJobLogSummarySystemJob(jobLogSummarySystemJobMap map[string]interface{}) schematicsv1.JobLogSummarySystemJob {
	jobLogSummarySystemJob := schematicsv1.JobLogSummarySystemJob{}

	if jobLogSummarySystemJobMap["target_count"] != nil {
		jobLogSummarySystemJob.TargetCount = core.Float64Ptr(jobLogSummarySystemJobMap["target_count"].(float64))
	}
	if jobLogSummarySystemJobMap["success"] != nil {
		jobLogSummarySystemJob.Success = core.Float64Ptr(jobLogSummarySystemJobMap["success"].(float64))
	}
	if jobLogSummarySystemJobMap["failed"] != nil {
		jobLogSummarySystemJob.Failed = core.Float64Ptr(jobLogSummarySystemJobMap["failed"].(float64))
	}

	return jobLogSummarySystemJob
}

func resourceIBMSchematicsJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}
	jobIDSplit := strings.Split(d.Id(), ".")
	region := jobIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	getJobOptions := &schematicsv1.GetJobOptions{}

	getJobOptions.SetJobID(d.Id())

	job, response, err := schematicsClient.GetJobWithContext(context, getJobOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetJobWithContext failed %s\n%s", err, response))
	}
	if err = d.Set("command_object", job.CommandObject); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_object: %s", err))
	}
	if err = d.Set("command_object_id", job.CommandObjectID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_object_id: %s", err))
	}
	if err = d.Set("command_name", job.CommandName); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_name: %s", err))
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		if err = d.Set("command_parameter", d.Get("command_parameter").(string)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_parameter: %s", err))
		}
	}
	if job.CommandOptions != nil {
		if err = d.Set("command_options", job.CommandOptions); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting command_options: %s", err))
		}
	}
	if job.Inputs != nil {
		jobInputs := []map[string]interface{}{}
		for _, jobInputsItem := range job.Inputs {
			jobInputsItemMap := resourceIBMSchematicsJobVariableDataToMap(jobInputsItem)
			jobInputs = append(jobInputs, jobInputsItemMap)
		}
		if err = d.Set("job_inputs", jobInputs); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting job_inputs: %s", err))
		}
	}
	if job.Settings != nil {
		jobEnvSettings := []map[string]interface{}{}
		for _, jobEnvSettingsItem := range job.Settings {
			jobEnvSettingsItemMap := resourceIBMSchematicsJobVariableDataToMap(jobEnvSettingsItem)
			jobEnvSettings = append(jobEnvSettings, jobEnvSettingsItemMap)
		}
		if err = d.Set("job_env_settings", jobEnvSettings); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting job_env_settings: %s", err))
		}
	}
	if job.Tags != nil {
		if err = d.Set("tags", job.Tags); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting tags: %s", err))
		}
	}
	if err = d.Set("location", job.Location); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting location: %s", err))
	}
	if job.Status != nil {
		statusMap := resourceIBMSchematicsJobJobStatusToMap(*job.Status)
		if err = d.Set("status", []map[string]interface{}{statusMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
		}
	}
	if job.Data != nil {
		dataMap := resourceIBMSchematicsJobJobDataToMap(*job.Data)
		if err = d.Set("data", []map[string]interface{}{dataMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting data: %s", err))
		}
	}
	if job.Bastion != nil {
		bastionMap := resourceIBMSchematicsJobBastionResourceDefinitionToMap(*job.Bastion)
		if err = d.Set("bastion", []map[string]interface{}{bastionMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting bastion: %s", err))
		}
	}
	if job.LogSummary != nil {
		logSummaryMap := resourceIBMSchematicsJobJobLogSummaryToMap(*job.LogSummary)
		if err = d.Set("log_summary", []map[string]interface{}{logSummaryMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting log_summary: %s", err))
		}
	}
	if err = d.Set("name", job.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
	}
	if err = d.Set("description", job.Description); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting description: %s", err))
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

func resourceIBMSchematicsJobVariableDataToMap(variableData schematicsv1.VariableData) map[string]interface{} {
	variableDataMap := map[string]interface{}{}

	if variableData.Name != nil {
		variableDataMap["name"] = variableData.Name
	}
	if variableData.Value != nil {
		variableDataMap["value"] = variableData.Value
	}
	if variableData.Metadata != nil {
		MetadataMap := resourceIBMSchematicsJobVariableMetadataToMap(*variableData.Metadata)
		variableDataMap["metadata"] = []map[string]interface{}{MetadataMap}
	}
	if variableData.Link != nil {
		variableDataMap["link"] = variableData.Link
	}

	return variableDataMap
}

func resourceIBMSchematicsJobVariableMetadataToMap(variableMetadata schematicsv1.VariableMetadata) map[string]interface{} {
	variableMetadataMap := map[string]interface{}{}

	if variableMetadata.Type != nil {
		variableMetadataMap["type"] = variableMetadata.Type
	}
	if variableMetadata.Aliases != nil {
		variableMetadataMap["aliases"] = variableMetadata.Aliases
	}
	if variableMetadata.Description != nil {
		variableMetadataMap["description"] = variableMetadata.Description
	}
	if variableMetadata.DefaultValue != nil {
		variableMetadataMap["default_value"] = variableMetadata.DefaultValue
	}
	if variableMetadata.Secure != nil {
		variableMetadataMap["secure"] = variableMetadata.Secure
	}
	if variableMetadata.Immutable != nil {
		variableMetadataMap["immutable"] = variableMetadata.Immutable
	}
	if variableMetadata.Hidden != nil {
		variableMetadataMap["hidden"] = variableMetadata.Hidden
	}
	if variableMetadata.Options != nil {
		variableMetadataMap["options"] = variableMetadata.Options
	}
	if variableMetadata.MinValue != nil {
		variableMetadataMap["min_value"] = flex.IntValue(variableMetadata.MinValue)
	}
	if variableMetadata.MaxValue != nil {
		variableMetadataMap["max_value"] = flex.IntValue(variableMetadata.MaxValue)
	}
	if variableMetadata.MinLength != nil {
		variableMetadataMap["min_length"] = flex.IntValue(variableMetadata.MinLength)
	}
	if variableMetadata.MaxLength != nil {
		variableMetadataMap["max_length"] = flex.IntValue(variableMetadata.MaxLength)
	}
	if variableMetadata.Matches != nil {
		variableMetadataMap["matches"] = variableMetadata.Matches
	}
	if variableMetadata.Position != nil {
		variableMetadataMap["position"] = flex.IntValue(variableMetadata.Position)
	}
	if variableMetadata.GroupBy != nil {
		variableMetadataMap["group_by"] = variableMetadata.GroupBy
	}
	if variableMetadata.Source != nil {
		variableMetadataMap["source"] = variableMetadata.Source
	}

	return variableMetadataMap
}

func resourceIBMSchematicsJobJobStatusToMap(jobStatus schematicsv1.JobStatus) map[string]interface{} {
	jobStatusMap := map[string]interface{}{}

	if jobStatus.WorkspaceJobStatus != nil {
		WorkspaceJobStatusMap := resourceIBMSchematicsJobJobStatusWorkspaceToMap(*jobStatus.WorkspaceJobStatus)
		jobStatusMap["workspace_job_status"] = []map[string]interface{}{WorkspaceJobStatusMap}
	}
	if jobStatus.ActionJobStatus != nil {
		ActionJobStatusMap := resourceIBMSchematicsJobJobStatusActionToMap(*jobStatus.ActionJobStatus)
		jobStatusMap["action_job_status"] = []map[string]interface{}{ActionJobStatusMap}
	}
	if jobStatus.SystemJobStatus != nil {
		SystemJobStatusMap := resourceIBMSchematicsJobJobStatusSystemToMap(*jobStatus.SystemJobStatus)
		jobStatusMap["system_job_status"] = []map[string]interface{}{SystemJobStatusMap}
	}
	if jobStatus.FlowJobStatus != nil {
		FlowJobStatusMap := resourceIBMSchematicsJobJobStatusFlowToMap(*jobStatus.FlowJobStatus)
		jobStatusMap["flow_job_status"] = []map[string]interface{}{FlowJobStatusMap}
	}

	return jobStatusMap
}

func resourceIBMSchematicsJobJobStatusWorkspaceToMap(jobStatusWorkspace schematicsv1.JobStatusWorkspace) map[string]interface{} {
	jobStatusWorkspaceMap := map[string]interface{}{}

	if jobStatusWorkspace.WorkspaceName != nil {
		jobStatusWorkspaceMap["workspace_name"] = jobStatusWorkspace.WorkspaceName
	}
	if jobStatusWorkspace.StatusCode != nil {
		jobStatusWorkspaceMap["status_code"] = jobStatusWorkspace.StatusCode
	}
	if jobStatusWorkspace.StatusMessage != nil {
		jobStatusWorkspaceMap["status_message"] = jobStatusWorkspace.StatusMessage
	}
	if jobStatusWorkspace.FlowStatus != nil {
		FlowStatusMap := resourceIBMSchematicsJobJobStatusFlowToMap(*jobStatusWorkspace.FlowStatus)
		jobStatusWorkspaceMap["flow_status"] = []map[string]interface{}{FlowStatusMap}
	}
	if jobStatusWorkspace.TemplateStatus != nil {
		templateStatus := []map[string]interface{}{}
		for _, templateStatusItem := range jobStatusWorkspace.TemplateStatus {
			templateStatusItemMap := resourceIBMSchematicsJobJobStatusTemplateToMap(templateStatusItem)
			templateStatus = append(templateStatus, templateStatusItemMap)
			// TODO: handle TemplateStatus of type TypeList -- list of non-primitive, not model items
		}
		jobStatusWorkspaceMap["template_status"] = templateStatus
	}
	if jobStatusWorkspace.UpdatedAt != nil {
		jobStatusWorkspaceMap["updated_at"] = jobStatusWorkspace.UpdatedAt.String()
	}

	return jobStatusWorkspaceMap
}

func resourceIBMSchematicsJobJobStatusFlowToMap(jobStatusFlow schematicsv1.JobStatusFlow) map[string]interface{} {
	jobStatusFlowMap := map[string]interface{}{}

	if jobStatusFlow.FlowID != nil {
		jobStatusFlowMap["flow_id"] = jobStatusFlow.FlowID
	}
	if jobStatusFlow.FlowName != nil {
		jobStatusFlowMap["flow_name"] = jobStatusFlow.FlowName
	}
	if jobStatusFlow.StatusCode != nil {
		jobStatusFlowMap["status_code"] = jobStatusFlow.StatusCode
	}
	if jobStatusFlow.StatusMessage != nil {
		jobStatusFlowMap["status_message"] = jobStatusFlow.StatusMessage
	}
	if jobStatusFlow.Workitems != nil {
		workitems := []map[string]interface{}{}
		for _, workitemsItem := range jobStatusFlow.Workitems {
			workitemsItemMap := resourceIBMSchematicsJobJobStatusWorkitemToMap(workitemsItem)
			workitems = append(workitems, workitemsItemMap)
			// TODO: handle Workitems of type TypeList -- list of non-primitive, not model items
		}
		jobStatusFlowMap["workitems"] = workitems
	}
	if jobStatusFlow.UpdatedAt != nil {
		jobStatusFlowMap["updated_at"] = jobStatusFlow.UpdatedAt.String()
	}

	return jobStatusFlowMap
}

func resourceIBMSchematicsJobJobStatusWorkitemToMap(jobStatusWorkitem schematicsv1.JobStatusWorkitem) map[string]interface{} {
	jobStatusWorkitemMap := map[string]interface{}{}

	if jobStatusWorkitem.WorkspaceID != nil {
		jobStatusWorkitemMap["workspace_id"] = jobStatusWorkitem.WorkspaceID
	}
	if jobStatusWorkitem.WorkspaceName != nil {
		jobStatusWorkitemMap["workspace_name"] = jobStatusWorkitem.WorkspaceName
	}
	if jobStatusWorkitem.JobID != nil {
		jobStatusWorkitemMap["job_id"] = jobStatusWorkitem.JobID
	}
	if jobStatusWorkitem.StatusCode != nil {
		jobStatusWorkitemMap["status_code"] = jobStatusWorkitem.StatusCode
	}
	if jobStatusWorkitem.StatusMessage != nil {
		jobStatusWorkitemMap["status_message"] = jobStatusWorkitem.StatusMessage
	}
	if jobStatusWorkitem.UpdatedAt != nil {
		jobStatusWorkitemMap["updated_at"] = jobStatusWorkitem.UpdatedAt.String()
	}

	return jobStatusWorkitemMap
}

func resourceIBMSchematicsJobJobStatusTemplateToMap(jobStatusTemplate schematicsv1.JobStatusTemplate) map[string]interface{} {
	jobStatusTemplateMap := map[string]interface{}{}

	if jobStatusTemplate.TemplateID != nil {
		jobStatusTemplateMap["template_id"] = jobStatusTemplate.TemplateID
	}
	if jobStatusTemplate.TemplateName != nil {
		jobStatusTemplateMap["template_name"] = jobStatusTemplate.TemplateName
	}
	if jobStatusTemplate.FlowIndex != nil {
		jobStatusTemplateMap["flow_index"] = flex.IntValue(jobStatusTemplate.FlowIndex)
	}
	if jobStatusTemplate.StatusCode != nil {
		jobStatusTemplateMap["status_code"] = jobStatusTemplate.StatusCode
	}
	if jobStatusTemplate.StatusMessage != nil {
		jobStatusTemplateMap["status_message"] = jobStatusTemplate.StatusMessage
	}
	if jobStatusTemplate.UpdatedAt != nil {
		jobStatusTemplateMap["updated_at"] = jobStatusTemplate.UpdatedAt.String()
	}

	return jobStatusTemplateMap
}

func resourceIBMSchematicsJobJobStatusActionToMap(jobStatusAction schematicsv1.JobStatusAction) map[string]interface{} {
	jobStatusActionMap := map[string]interface{}{}

	if jobStatusAction.ActionName != nil {
		jobStatusActionMap["action_name"] = jobStatusAction.ActionName
	}
	if jobStatusAction.StatusCode != nil {
		jobStatusActionMap["status_code"] = jobStatusAction.StatusCode
	}
	if jobStatusAction.StatusMessage != nil {
		jobStatusActionMap["status_message"] = jobStatusAction.StatusMessage
	}
	if jobStatusAction.BastionStatusCode != nil {
		jobStatusActionMap["bastion_status_code"] = jobStatusAction.BastionStatusCode
	}
	if jobStatusAction.BastionStatusMessage != nil {
		jobStatusActionMap["bastion_status_message"] = jobStatusAction.BastionStatusMessage
	}
	if jobStatusAction.TargetsStatusCode != nil {
		jobStatusActionMap["targets_status_code"] = jobStatusAction.TargetsStatusCode
	}
	if jobStatusAction.TargetsStatusMessage != nil {
		jobStatusActionMap["targets_status_message"] = jobStatusAction.TargetsStatusMessage
	}
	if jobStatusAction.UpdatedAt != nil {
		jobStatusActionMap["updated_at"] = jobStatusAction.UpdatedAt.String()
	}

	return jobStatusActionMap
}

func resourceIBMSchematicsJobJobStatusSystemToMap(jobStatusSystem schematicsv1.JobStatusSystem) map[string]interface{} {
	jobStatusSystemMap := map[string]interface{}{}

	if jobStatusSystem.SystemStatusMessage != nil {
		jobStatusSystemMap["system_status_message"] = jobStatusSystem.SystemStatusMessage
	}
	if jobStatusSystem.SystemStatusCode != nil {
		jobStatusSystemMap["system_status_code"] = jobStatusSystem.SystemStatusCode
	}
	if jobStatusSystem.SchematicsResourceStatus != nil {
		schematicsResourceStatus := []map[string]interface{}{}
		for _, schematicsResourceStatusItem := range jobStatusSystem.SchematicsResourceStatus {
			schematicsResourceStatusItemMap := resourceIBMSchematicsJobJobStatusSchematicsResourcesToMap(schematicsResourceStatusItem)
			schematicsResourceStatus = append(schematicsResourceStatus, schematicsResourceStatusItemMap)
			// TODO: handle SchematicsResourceStatus of type TypeList -- list of non-primitive, not model items
		}
		jobStatusSystemMap["schematics_resource_status"] = schematicsResourceStatus
	}
	if jobStatusSystem.UpdatedAt != nil {
		jobStatusSystemMap["updated_at"] = jobStatusSystem.UpdatedAt.String()
	}

	return jobStatusSystemMap
}

func resourceIBMSchematicsJobJobStatusSchematicsResourcesToMap(jobStatusSchematicsResources schematicsv1.JobStatusSchematicsResources) map[string]interface{} {
	jobStatusSchematicsResourcesMap := map[string]interface{}{}

	if jobStatusSchematicsResources.StatusCode != nil {
		jobStatusSchematicsResourcesMap["status_code"] = jobStatusSchematicsResources.StatusCode
	}
	if jobStatusSchematicsResources.StatusMessage != nil {
		jobStatusSchematicsResourcesMap["status_message"] = jobStatusSchematicsResources.StatusMessage
	}
	if jobStatusSchematicsResources.SchematicsResourceID != nil {
		jobStatusSchematicsResourcesMap["schematics_resource_id"] = jobStatusSchematicsResources.SchematicsResourceID
	}
	if jobStatusSchematicsResources.UpdatedAt != nil {
		jobStatusSchematicsResourcesMap["updated_at"] = jobStatusSchematicsResources.UpdatedAt.String()
	}

	return jobStatusSchematicsResourcesMap
}

func resourceIBMSchematicsJobJobDataToMap(jobData schematicsv1.JobData) map[string]interface{} {
	jobDataMap := map[string]interface{}{}

	jobDataMap["job_type"] = jobData.JobType
	if jobData.WorkspaceJobData != nil {
		WorkspaceJobDataMap := resourceIBMSchematicsJobJobDataWorkspaceToMap(*jobData.WorkspaceJobData)
		jobDataMap["workspace_job_data"] = []map[string]interface{}{WorkspaceJobDataMap}
	}
	if jobData.ActionJobData != nil {
		ActionJobDataMap := resourceIBMSchematicsJobJobDataActionToMap(*jobData.ActionJobData)
		jobDataMap["action_job_data"] = []map[string]interface{}{ActionJobDataMap}
	}
	if jobData.SystemJobData != nil {
		SystemJobDataMap := resourceIBMSchematicsJobJobDataSystemToMap(*jobData.SystemJobData)
		jobDataMap["system_job_data"] = []map[string]interface{}{SystemJobDataMap}
	}
	if jobData.FlowJobData != nil {
		FlowJobDataMap := resourceIBMSchematicsJobJobDataFlowToMap(*jobData.FlowJobData)
		jobDataMap["flow_job_data"] = []map[string]interface{}{FlowJobDataMap}
	}

	return jobDataMap
}

func resourceIBMSchematicsJobJobDataWorkspaceToMap(jobDataWorkspace schematicsv1.JobDataWorkspace) map[string]interface{} {
	jobDataWorkspaceMap := map[string]interface{}{}

	if jobDataWorkspace.WorkspaceName != nil {
		jobDataWorkspaceMap["workspace_name"] = jobDataWorkspace.WorkspaceName
	}
	if jobDataWorkspace.FlowID != nil {
		jobDataWorkspaceMap["flow_id"] = jobDataWorkspace.FlowID
	}
	if jobDataWorkspace.FlowName != nil {
		jobDataWorkspaceMap["flow_name"] = jobDataWorkspace.FlowName
	}
	if jobDataWorkspace.Inputs != nil {
		inputs := []map[string]interface{}{}
		for _, inputsItem := range jobDataWorkspace.Inputs {
			inputsItemMap := resourceIBMSchematicsJobVariableDataToMap(inputsItem)
			inputs = append(inputs, inputsItemMap)
			// TODO: handle Inputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkspaceMap["inputs"] = inputs
	}
	if jobDataWorkspace.Outputs != nil {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range jobDataWorkspace.Outputs {
			outputsItemMap := resourceIBMSchematicsJobVariableDataToMap(outputsItem)
			outputs = append(outputs, outputsItemMap)
			// TODO: handle Outputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkspaceMap["outputs"] = outputs
	}
	if jobDataWorkspace.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range jobDataWorkspace.Settings {
			settingsItemMap := resourceIBMSchematicsJobVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
			// TODO: handle Settings of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkspaceMap["settings"] = settings
	}
	if jobDataWorkspace.TemplateData != nil {
		templateData := []map[string]interface{}{}
		for _, templateDataItem := range jobDataWorkspace.TemplateData {
			templateDataItemMap := resourceIBMSchematicsJobJobDataTemplateToMap(templateDataItem)
			templateData = append(templateData, templateDataItemMap)
			// TODO: handle TemplateData of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkspaceMap["template_data"] = templateData
	}
	if jobDataWorkspace.UpdatedAt != nil {
		jobDataWorkspaceMap["updated_at"] = jobDataWorkspace.UpdatedAt.String()
	}

	return jobDataWorkspaceMap
}

func resourceIBMSchematicsJobJobDataTemplateToMap(jobDataTemplate schematicsv1.JobDataTemplate) map[string]interface{} {
	jobDataTemplateMap := map[string]interface{}{}

	if jobDataTemplate.TemplateID != nil {
		jobDataTemplateMap["template_id"] = jobDataTemplate.TemplateID
	}
	if jobDataTemplate.TemplateName != nil {
		jobDataTemplateMap["template_name"] = jobDataTemplate.TemplateName
	}
	if jobDataTemplate.FlowIndex != nil {
		jobDataTemplateMap["flow_index"] = flex.IntValue(jobDataTemplate.FlowIndex)
	}
	if jobDataTemplate.Inputs != nil {
		inputs := []map[string]interface{}{}
		for _, inputsItem := range jobDataTemplate.Inputs {
			inputsItemMap := resourceIBMSchematicsJobVariableDataToMap(inputsItem)
			inputs = append(inputs, inputsItemMap)
			// TODO: handle Inputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataTemplateMap["inputs"] = inputs
	}
	if jobDataTemplate.Outputs != nil {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range jobDataTemplate.Outputs {
			outputsItemMap := resourceIBMSchematicsJobVariableDataToMap(outputsItem)
			outputs = append(outputs, outputsItemMap)
			// TODO: handle Outputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataTemplateMap["outputs"] = outputs
	}
	if jobDataTemplate.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range jobDataTemplate.Settings {
			settingsItemMap := resourceIBMSchematicsJobVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
			// TODO: handle Settings of type TypeList -- list of non-primitive, not model items
		}
		jobDataTemplateMap["settings"] = settings
	}
	if jobDataTemplate.UpdatedAt != nil {
		jobDataTemplateMap["updated_at"] = jobDataTemplate.UpdatedAt.String()
	}

	return jobDataTemplateMap
}

func resourceIBMSchematicsJobJobDataActionToMap(jobDataAction schematicsv1.JobDataAction) map[string]interface{} {
	jobDataActionMap := map[string]interface{}{}

	if jobDataAction.ActionName != nil {
		jobDataActionMap["action_name"] = jobDataAction.ActionName
	}
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
	if jobDataAction.UpdatedAt != nil {
		jobDataActionMap["updated_at"] = jobDataAction.UpdatedAt.String()
	}
	if jobDataAction.InventoryRecord != nil {
		InventoryRecordMap := resourceIBMSchematicsJobInventoryResourceRecordToMap(*jobDataAction.InventoryRecord)
		jobDataActionMap["inventory_record"] = []map[string]interface{}{InventoryRecordMap}
	}
	if jobDataAction.MaterializedInventory != nil {
		jobDataActionMap["materialized_inventory"] = jobDataAction.MaterializedInventory
	}

	return jobDataActionMap
}

func resourceIBMSchematicsJobInventoryResourceRecordToMap(inventoryResourceRecord schematicsv1.InventoryResourceRecord) map[string]interface{} {
	inventoryResourceRecordMap := map[string]interface{}{}

	if inventoryResourceRecord.Name != nil {
		inventoryResourceRecordMap["name"] = inventoryResourceRecord.Name
	}
	if inventoryResourceRecord.ID != nil {
		inventoryResourceRecordMap["id"] = inventoryResourceRecord.ID
	}
	if inventoryResourceRecord.Description != nil {
		inventoryResourceRecordMap["description"] = inventoryResourceRecord.Description
	}
	if inventoryResourceRecord.Location != nil {
		inventoryResourceRecordMap["location"] = inventoryResourceRecord.Location
	}
	if inventoryResourceRecord.ResourceGroup != nil {
		inventoryResourceRecordMap["resource_group"] = inventoryResourceRecord.ResourceGroup
	}
	if inventoryResourceRecord.CreatedAt != nil {
		inventoryResourceRecordMap["created_at"] = inventoryResourceRecord.CreatedAt.String()
	}
	if inventoryResourceRecord.CreatedBy != nil {
		inventoryResourceRecordMap["created_by"] = inventoryResourceRecord.CreatedBy
	}
	if inventoryResourceRecord.UpdatedAt != nil {
		inventoryResourceRecordMap["updated_at"] = inventoryResourceRecord.UpdatedAt.String()
	}
	if inventoryResourceRecord.UpdatedBy != nil {
		inventoryResourceRecordMap["updated_by"] = inventoryResourceRecord.UpdatedBy
	}
	if inventoryResourceRecord.InventoriesIni != nil {
		inventoryResourceRecordMap["inventories_ini"] = inventoryResourceRecord.InventoriesIni
	}
	if inventoryResourceRecord.ResourceQueries != nil {
		inventoryResourceRecordMap["resource_queries"] = inventoryResourceRecord.ResourceQueries
	}

	return inventoryResourceRecordMap
}

func resourceIBMSchematicsJobJobDataSystemToMap(jobDataSystem schematicsv1.JobDataSystem) map[string]interface{} {
	jobDataSystemMap := map[string]interface{}{}

	if jobDataSystem.KeyID != nil {
		jobDataSystemMap["key_id"] = jobDataSystem.KeyID
	}
	if jobDataSystem.SchematicsResourceID != nil {
		jobDataSystemMap["schematics_resource_id"] = jobDataSystem.SchematicsResourceID
	}
	if jobDataSystem.UpdatedAt != nil {
		jobDataSystemMap["updated_at"] = jobDataSystem.UpdatedAt.String()
	}

	return jobDataSystemMap
}

func resourceIBMSchematicsJobJobDataFlowToMap(jobDataFlow schematicsv1.JobDataFlow) map[string]interface{} {
	jobDataFlowMap := map[string]interface{}{}

	if jobDataFlow.FlowID != nil {
		jobDataFlowMap["flow_id"] = jobDataFlow.FlowID
	}
	if jobDataFlow.FlowName != nil {
		jobDataFlowMap["flow_name"] = jobDataFlow.FlowName
	}
	if jobDataFlow.Workitems != nil {
		workitems := []map[string]interface{}{}
		for _, workitemsItem := range jobDataFlow.Workitems {
			workitemsItemMap := resourceIBMSchematicsJobJobDataWorkItemToMap(workitemsItem)
			workitems = append(workitems, workitemsItemMap)
			// TODO: handle Workitems of type TypeList -- list of non-primitive, not model items
		}
		jobDataFlowMap["workitems"] = workitems
	}
	if jobDataFlow.UpdatedAt != nil {
		jobDataFlowMap["updated_at"] = jobDataFlow.UpdatedAt.String()
	}

	return jobDataFlowMap
}

func resourceIBMSchematicsJobJobDataWorkItemToMap(jobDataWorkItem schematicsv1.JobDataWorkItem) map[string]interface{} {
	jobDataWorkItemMap := map[string]interface{}{}

	if jobDataWorkItem.CommandObjectID != nil {
		jobDataWorkItemMap["command_object_id"] = jobDataWorkItem.CommandObjectID
	}
	if jobDataWorkItem.CommandObjectName != nil {
		jobDataWorkItemMap["command_object_name"] = jobDataWorkItem.CommandObjectName
	}
	if jobDataWorkItem.Layers != nil {
		jobDataWorkItemMap["layers"] = jobDataWorkItem.Layers
	}
	if jobDataWorkItem.SourceType != nil {
		jobDataWorkItemMap["source_type"] = jobDataWorkItem.SourceType
	}
	if jobDataWorkItem.Source != nil {
		SourceMap := resourceIBMSchematicsJobExternalSourceToMap(*jobDataWorkItem.Source)
		jobDataWorkItemMap["source"] = []map[string]interface{}{SourceMap}
	}
	if jobDataWorkItem.Inputs != nil {
		inputs := []map[string]interface{}{}
		for _, inputsItem := range jobDataWorkItem.Inputs {
			inputsItemMap := resourceIBMSchematicsJobVariableDataToMap(inputsItem)
			inputs = append(inputs, inputsItemMap)
			// TODO: handle Inputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkItemMap["inputs"] = inputs
	}
	if jobDataWorkItem.Outputs != nil {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range jobDataWorkItem.Outputs {
			outputsItemMap := resourceIBMSchematicsJobVariableDataToMap(outputsItem)
			outputs = append(outputs, outputsItemMap)
			// TODO: handle Outputs of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkItemMap["outputs"] = outputs
	}
	if jobDataWorkItem.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range jobDataWorkItem.Settings {
			settingsItemMap := resourceIBMSchematicsJobVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
			// TODO: handle Settings of type TypeList -- list of non-primitive, not model items
		}
		jobDataWorkItemMap["settings"] = settings
	}
	if jobDataWorkItem.LastJob != nil {
		LastJobMap := resourceIBMSchematicsJobJobDataWorkItemLastJobToMap(*jobDataWorkItem.LastJob)
		jobDataWorkItemMap["last_job"] = []map[string]interface{}{LastJobMap}
	}
	if jobDataWorkItem.UpdatedAt != nil {
		jobDataWorkItemMap["updated_at"] = jobDataWorkItem.UpdatedAt.String()
	}

	return jobDataWorkItemMap
}

func resourceIBMSchematicsJobExternalSourceToMap(externalSource schematicsv1.ExternalSource) map[string]interface{} {
	externalSourceMap := map[string]interface{}{}

	externalSourceMap["source_type"] = externalSource.SourceType
	if externalSource.Git != nil {
		GitMap := resourceIBMSchematicsJobExternalSourceGitToMap(*externalSource.Git)
		externalSourceMap["git"] = []map[string]interface{}{GitMap}
	}
	if externalSource.Catalog != nil {
		CatalogMap := resourceIBMSchematicsJobExternalSourceCatalogToMap(*externalSource.Catalog)
		externalSourceMap["catalog"] = []map[string]interface{}{CatalogMap}
	}
	// if externalSource.CosBucket != nil {
	// 	CosBucketMap := resourceIBMSchematicsJobExternalSourceCosBucketToMap(*externalSource.CosBucket)
	// 	externalSourceMap["cos_bucket"] = []map[string]interface{}{CosBucketMap}
	// }

	return externalSourceMap
}

func resourceIBMSchematicsJobExternalSourceGitToMap(externalSourceGit schematicsv1.GitSource) map[string]interface{} {
	externalSourceGitMap := map[string]interface{}{}

	if externalSourceGit.ComputedGitRepoURL != nil {
		externalSourceGitMap["computed_git_repo_url"] = externalSourceGit.ComputedGitRepoURL
	}
	if externalSourceGit.GitRepoURL != nil {
		externalSourceGitMap["git_repo_url"] = externalSourceGit.GitRepoURL
	}
	if externalSourceGit.GitToken != nil {
		externalSourceGitMap["git_token"] = externalSourceGit.GitToken
	}
	if externalSourceGit.GitRepoFolder != nil {
		externalSourceGitMap["git_repo_folder"] = externalSourceGit.GitRepoFolder
	}
	if externalSourceGit.GitRelease != nil {
		externalSourceGitMap["git_release"] = externalSourceGit.GitRelease
	}
	if externalSourceGit.GitBranch != nil {
		externalSourceGitMap["git_branch"] = externalSourceGit.GitBranch
	}

	return externalSourceGitMap
}

func resourceIBMSchematicsJobExternalSourceCatalogToMap(externalSourceCatalog schematicsv1.CatalogSource) map[string]interface{} {
	externalSourceCatalogMap := map[string]interface{}{}

	if externalSourceCatalog.CatalogName != nil {
		externalSourceCatalogMap["catalog_name"] = externalSourceCatalog.CatalogName
	}
	if externalSourceCatalog.OfferingName != nil {
		externalSourceCatalogMap["offering_name"] = externalSourceCatalog.OfferingName
	}
	if externalSourceCatalog.OfferingVersion != nil {
		externalSourceCatalogMap["offering_version"] = externalSourceCatalog.OfferingVersion
	}
	if externalSourceCatalog.OfferingKind != nil {
		externalSourceCatalogMap["offering_kind"] = externalSourceCatalog.OfferingKind
	}
	if externalSourceCatalog.OfferingID != nil {
		externalSourceCatalogMap["offering_id"] = externalSourceCatalog.OfferingID
	}
	if externalSourceCatalog.OfferingVersionID != nil {
		externalSourceCatalogMap["offering_version_id"] = externalSourceCatalog.OfferingVersionID
	}
	if externalSourceCatalog.OfferingRepoURL != nil {
		externalSourceCatalogMap["offering_repo_url"] = externalSourceCatalog.OfferingRepoURL
	}

	return externalSourceCatalogMap
}

// func resourceIBMSchematicsJobExternalSourceCosBucketToMap(externalSourceCosBucket schematicsv1.ExternalSourceCosBucket) map[string]interface{} {
// 	externalSourceCosBucketMap := map[string]interface{}{}

// 	if externalSourceCosBucket.CosBucketURL != nil {
// 		externalSourceCosBucketMap["cos_bucket_url"] = externalSourceCosBucket.CosBucketURL
// 	}

// 	return externalSourceCosBucketMap
// }

func resourceIBMSchematicsJobJobDataWorkItemLastJobToMap(jobDataWorkItemLastJob schematicsv1.JobDataWorkItemLastJob) map[string]interface{} {
	jobDataWorkItemLastJobMap := map[string]interface{}{}

	if jobDataWorkItemLastJob.CommandObject != nil {
		jobDataWorkItemLastJobMap["command_object"] = jobDataWorkItemLastJob.CommandObject
	}
	if jobDataWorkItemLastJob.CommandObjectName != nil {
		jobDataWorkItemLastJobMap["command_object_name"] = jobDataWorkItemLastJob.CommandObjectName
	}
	if jobDataWorkItemLastJob.CommandObjectID != nil {
		jobDataWorkItemLastJobMap["command_object_id"] = jobDataWorkItemLastJob.CommandObjectID
	}
	if jobDataWorkItemLastJob.CommandName != nil {
		jobDataWorkItemLastJobMap["command_name"] = jobDataWorkItemLastJob.CommandName
	}
	if jobDataWorkItemLastJob.JobID != nil {
		jobDataWorkItemLastJobMap["job_id"] = jobDataWorkItemLastJob.JobID
	}
	if jobDataWorkItemLastJob.JobStatus != nil {
		jobDataWorkItemLastJobMap["job_status"] = jobDataWorkItemLastJob.JobStatus
	}

	return jobDataWorkItemLastJobMap
}

func resourceIBMSchematicsJobBastionResourceDefinitionToMap(bastionResourceDefinition schematicsv1.BastionResourceDefinition) map[string]interface{} {
	bastionResourceDefinitionMap := map[string]interface{}{}

	if bastionResourceDefinition.Name != nil {
		bastionResourceDefinitionMap["name"] = bastionResourceDefinition.Name
	}
	if bastionResourceDefinition.Host != nil {
		bastionResourceDefinitionMap["host"] = bastionResourceDefinition.Host
	}

	return bastionResourceDefinitionMap
}

func resourceIBMSchematicsJobJobLogSummaryToMap(jobLogSummary schematicsv1.JobLogSummary) map[string]interface{} {
	jobLogSummaryMap := map[string]interface{}{}

	if jobLogSummary.JobID != nil {
		jobLogSummaryMap["job_id"] = jobLogSummary.JobID
	}
	if jobLogSummary.JobType != nil {
		jobLogSummaryMap["job_type"] = jobLogSummary.JobType
	}
	if jobLogSummary.LogStartAt != nil {
		jobLogSummaryMap["log_start_at"] = jobLogSummary.LogStartAt.String()
	}
	if jobLogSummary.LogAnalyzedTill != nil {
		jobLogSummaryMap["log_analyzed_till"] = jobLogSummary.LogAnalyzedTill.String()
	}
	if jobLogSummary.ElapsedTime != nil {
		jobLogSummaryMap["elapsed_time"] = jobLogSummary.ElapsedTime
	}
	if jobLogSummary.LogErrors != nil {
		logErrors := []map[string]interface{}{}
		for _, logErrorsItem := range jobLogSummary.LogErrors {
			logErrorsItemMap := resourceIBMSchematicsJobJobLogSummaryLogErrorsToMap(logErrorsItem)
			logErrors = append(logErrors, logErrorsItemMap)
			// TODO: handle LogErrors of type TypeList -- list of non-primitive, not model items
		}
		jobLogSummaryMap["log_errors"] = logErrors
	}
	if jobLogSummary.RepoDownloadJob != nil {
		RepoDownloadJobMap := resourceIBMSchematicsJobJobLogSummaryRepoDownloadJobToMap(*jobLogSummary.RepoDownloadJob)
		jobLogSummaryMap["repo_download_job"] = []map[string]interface{}{RepoDownloadJobMap}
	}
	if jobLogSummary.WorkspaceJob != nil {
		WorkspaceJobMap := resourceIBMSchematicsJobJobLogSummaryWorkspaceJobToMap(*jobLogSummary.WorkspaceJob)
		jobLogSummaryMap["workspace_job"] = []map[string]interface{}{WorkspaceJobMap}
	}
	if jobLogSummary.FlowJob != nil {
		FlowJobMap := resourceIBMSchematicsJobJobLogSummaryFlowJobToMap(*jobLogSummary.FlowJob)
		jobLogSummaryMap["flow_job"] = []map[string]interface{}{FlowJobMap}
	}
	if jobLogSummary.ActionJob != nil {
		ActionJobMap := resourceIBMSchematicsJobJobLogSummaryActionJobToMap(*jobLogSummary.ActionJob)
		jobLogSummaryMap["action_job"] = []map[string]interface{}{ActionJobMap}
	}
	if jobLogSummary.SystemJob != nil {
		SystemJobMap := resourceIBMSchematicsJobJobLogSummarySystemJobToMap(*jobLogSummary.SystemJob)
		jobLogSummaryMap["system_job"] = []map[string]interface{}{SystemJobMap}
	}

	return jobLogSummaryMap
}

func resourceIBMSchematicsJobJobLogSummaryLogErrorsToMap(jobLogSummaryLogErrors schematicsv1.JobLogSummaryLogErrors) map[string]interface{} {
	jobLogSummaryLogErrorsMap := map[string]interface{}{}

	if jobLogSummaryLogErrors.ErrorCode != nil {
		jobLogSummaryLogErrorsMap["error_code"] = jobLogSummaryLogErrors.ErrorCode
	}
	if jobLogSummaryLogErrors.ErrorMsg != nil {
		jobLogSummaryLogErrorsMap["error_msg"] = jobLogSummaryLogErrors.ErrorMsg
	}
	if jobLogSummaryLogErrors.ErrorCount != nil {
		jobLogSummaryLogErrorsMap["error_count"] = jobLogSummaryLogErrors.ErrorCount
	}

	return jobLogSummaryLogErrorsMap
}

func resourceIBMSchematicsJobJobLogSummaryRepoDownloadJobToMap(jobLogSummaryRepoDownloadJob schematicsv1.JobLogSummaryRepoDownloadJob) map[string]interface{} {
	jobLogSummaryRepoDownloadJobMap := map[string]interface{}{}

	if jobLogSummaryRepoDownloadJob.ScannedFileCount != nil {
		jobLogSummaryRepoDownloadJobMap["scanned_file_count"] = jobLogSummaryRepoDownloadJob.ScannedFileCount
	}
	if jobLogSummaryRepoDownloadJob.QuarantinedFileCount != nil {
		jobLogSummaryRepoDownloadJobMap["quarantined_file_count"] = jobLogSummaryRepoDownloadJob.QuarantinedFileCount
	}
	if jobLogSummaryRepoDownloadJob.DetectedFiletype != nil {
		jobLogSummaryRepoDownloadJobMap["detected_filetype"] = jobLogSummaryRepoDownloadJob.DetectedFiletype
	}
	if jobLogSummaryRepoDownloadJob.InputsCount != nil {
		jobLogSummaryRepoDownloadJobMap["inputs_count"] = jobLogSummaryRepoDownloadJob.InputsCount
	}
	if jobLogSummaryRepoDownloadJob.OutputsCount != nil {
		jobLogSummaryRepoDownloadJobMap["outputs_count"] = jobLogSummaryRepoDownloadJob.OutputsCount
	}

	return jobLogSummaryRepoDownloadJobMap
}

func resourceIBMSchematicsJobJobLogSummaryWorkspaceJobToMap(jobLogSummaryWorkspaceJob schematicsv1.JobLogSummaryWorkspaceJob) map[string]interface{} {
	jobLogSummaryWorkspaceJobMap := map[string]interface{}{}

	if jobLogSummaryWorkspaceJob.ResourcesAdd != nil {
		jobLogSummaryWorkspaceJobMap["resources_add"] = jobLogSummaryWorkspaceJob.ResourcesAdd
	}
	if jobLogSummaryWorkspaceJob.ResourcesModify != nil {
		jobLogSummaryWorkspaceJobMap["resources_modify"] = jobLogSummaryWorkspaceJob.ResourcesModify
	}
	if jobLogSummaryWorkspaceJob.ResourcesDestroy != nil {
		jobLogSummaryWorkspaceJobMap["resources_destroy"] = jobLogSummaryWorkspaceJob.ResourcesDestroy
	}

	return jobLogSummaryWorkspaceJobMap
}

func resourceIBMSchematicsJobJobLogSummaryFlowJobToMap(jobLogSummaryFlowJob schematicsv1.JobLogSummaryFlowJob) map[string]interface{} {
	jobLogSummaryFlowJobMap := map[string]interface{}{}

	if jobLogSummaryFlowJob.WorkitemsCompleted != nil {
		jobLogSummaryFlowJobMap["workitems_completed"] = jobLogSummaryFlowJob.WorkitemsCompleted
	}
	if jobLogSummaryFlowJob.WorkitemsPending != nil {
		jobLogSummaryFlowJobMap["workitems_pending"] = jobLogSummaryFlowJob.WorkitemsPending
	}
	if jobLogSummaryFlowJob.WorkitemsFailed != nil {
		jobLogSummaryFlowJobMap["workitems_failed"] = jobLogSummaryFlowJob.WorkitemsFailed
	}
	if jobLogSummaryFlowJob.Workitems != nil {
		workitems := []map[string]interface{}{}
		for _, workitemsItem := range jobLogSummaryFlowJob.Workitems {
			workitemsItemMap := resourceIBMSchematicsJobJobLogSummaryWorkitemsToMap(workitemsItem)
			workitems = append(workitems, workitemsItemMap)
			// TODO: handle Workitems of type TypeList -- list of non-primitive, not model items
		}
		jobLogSummaryFlowJobMap["workitems"] = workitems
	}

	return jobLogSummaryFlowJobMap
}

func resourceIBMSchematicsJobJobLogSummaryWorkitemsToMap(jobLogSummaryWorkitems schematicsv1.JobLogSummaryWorkitems) map[string]interface{} {
	jobLogSummaryWorkitemsMap := map[string]interface{}{}

	if jobLogSummaryWorkitems.WorkspaceID != nil {
		jobLogSummaryWorkitemsMap["workspace_id"] = jobLogSummaryWorkitems.WorkspaceID
	}
	if jobLogSummaryWorkitems.JobID != nil {
		jobLogSummaryWorkitemsMap["job_id"] = jobLogSummaryWorkitems.JobID
	}
	if jobLogSummaryWorkitems.ResourcesAdd != nil {
		jobLogSummaryWorkitemsMap["resources_add"] = jobLogSummaryWorkitems.ResourcesAdd
	}
	if jobLogSummaryWorkitems.ResourcesModify != nil {
		jobLogSummaryWorkitemsMap["resources_modify"] = jobLogSummaryWorkitems.ResourcesModify
	}
	if jobLogSummaryWorkitems.ResourcesDestroy != nil {
		jobLogSummaryWorkitemsMap["resources_destroy"] = jobLogSummaryWorkitems.ResourcesDestroy
	}
	if jobLogSummaryWorkitems.LogURL != nil {
		jobLogSummaryWorkitemsMap["log_url"] = jobLogSummaryWorkitems.LogURL
	}

	return jobLogSummaryWorkitemsMap
}

func resourceIBMSchematicsJobJobLogSummaryActionJobToMap(jobLogSummaryActionJob schematicsv1.JobLogSummaryActionJob) map[string]interface{} {
	jobLogSummaryActionJobMap := map[string]interface{}{}

	if jobLogSummaryActionJob.TargetCount != nil {
		jobLogSummaryActionJobMap["target_count"] = jobLogSummaryActionJob.TargetCount
	}
	if jobLogSummaryActionJob.TaskCount != nil {
		jobLogSummaryActionJobMap["task_count"] = jobLogSummaryActionJob.TaskCount
	}
	if jobLogSummaryActionJob.PlayCount != nil {
		jobLogSummaryActionJobMap["play_count"] = jobLogSummaryActionJob.PlayCount
	}
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
	if jobLogSummaryActionJobRecap.Ok != nil {
		jobLogSummaryActionJobRecapMap["ok"] = jobLogSummaryActionJobRecap.Ok
	}
	if jobLogSummaryActionJobRecap.Changed != nil {
		jobLogSummaryActionJobRecapMap["changed"] = jobLogSummaryActionJobRecap.Changed
	}
	if jobLogSummaryActionJobRecap.Failed != nil {
		jobLogSummaryActionJobRecapMap["failed"] = jobLogSummaryActionJobRecap.Failed
	}
	if jobLogSummaryActionJobRecap.Skipped != nil {
		jobLogSummaryActionJobRecapMap["skipped"] = jobLogSummaryActionJobRecap.Skipped
	}
	if jobLogSummaryActionJobRecap.Unreachable != nil {
		jobLogSummaryActionJobRecapMap["unreachable"] = jobLogSummaryActionJobRecap.Unreachable
	}

	return jobLogSummaryActionJobRecapMap
}

func resourceIBMSchematicsJobJobLogSummarySystemJobToMap(jobLogSummarySystemJob schematicsv1.JobLogSummarySystemJob) map[string]interface{} {
	jobLogSummarySystemJobMap := map[string]interface{}{}

	if jobLogSummarySystemJob.TargetCount != nil {
		jobLogSummarySystemJobMap["target_count"] = jobLogSummarySystemJob.TargetCount
	}
	if jobLogSummarySystemJob.Success != nil {
		jobLogSummarySystemJobMap["success"] = jobLogSummarySystemJob.Success
	}
	if jobLogSummarySystemJob.Failed != nil {
		jobLogSummarySystemJobMap["failed"] = jobLogSummarySystemJob.Failed
	}

	return jobLogSummarySystemJobMap
}

func resourceIBMSchematicsJobUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}

	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}

	iamRefreshToken := session.Config.IAMRefreshToken
	jobIDSplit := strings.Split(d.Id(), ".")
	region := jobIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	updateJobOptions := &schematicsv1.UpdateJobOptions{}

	updateJobOptions.SetJobID(d.Id())
	updateJobOptions.SetRefreshToken(iamRefreshToken)

	if _, ok := d.GetOk("command_object"); ok {
		updateJobOptions.SetCommandObject(d.Get("command_object").(string))
	}
	if _, ok := d.GetOk("command_object_id"); ok {
		updateJobOptions.SetCommandObjectID(d.Get("command_object_id").(string))
	}
	if _, ok := d.GetOk("command_name"); ok {
		updateJobOptions.SetCommandName(d.Get("command_name").(string))
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		updateJobOptions.SetCommandParameter(d.Get("command_parameter").(string))
	}
	if _, ok := d.GetOk("command_options"); ok {
		updateJobOptions.SetCommandOptions(flex.ExpandStringList(d.Get("command_options").([]interface{})))
	}
	if _, ok := d.GetOk("job_inputs"); ok {
		var jobInputs []schematicsv1.VariableData
		for _, e := range d.Get("job_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			jobInputsItem := resourceIBMSchematicsJobMapToVariableData(value)
			jobInputs = append(jobInputs, jobInputsItem)
		}
		updateJobOptions.SetInputs(jobInputs)
	}
	if _, ok := d.GetOk("job_env_settings"); ok {
		var jobEnvSettings []schematicsv1.VariableData
		for _, e := range d.Get("job_env_settings").([]interface{}) {
			value := e.(map[string]interface{})
			jobEnvSettingsItem := resourceIBMSchematicsJobMapToVariableData(value)
			jobEnvSettings = append(jobEnvSettings, jobEnvSettingsItem)
		}
		updateJobOptions.SetSettings(jobEnvSettings)
	}
	if _, ok := d.GetOk("tags"); ok {
		updateJobOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("location"); ok {
		updateJobOptions.SetLocation(d.Get("location").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		statusAttr := d.Get("status").([]interface{})
		if len(statusAttr) > 0 {
			status := resourceIBMSchematicsJobMapToJobStatus(d.Get("status.0").(map[string]interface{}))
			updateJobOptions.SetStatus(&status)
		}
	}
	if _, ok := d.GetOk("data"); ok {
		dataAttr := d.Get("data").([]interface{})
		if len(dataAttr) > 0 {
			data := resourceIBMSchematicsJobMapToJobData(d.Get("data.0").(map[string]interface{}))
			updateJobOptions.SetData(&data)
		}
	}
	if _, ok := d.GetOk("bastion"); ok {
		bastionAttr := d.Get("bastion").([]interface{})
		if len(bastionAttr) > 0 {
			bastion := resourceIBMSchematicsJobMapToBastionResourceDefinition(d.Get("bastion.0").(map[string]interface{}))
			updateJobOptions.SetBastion(&bastion)
		}
	}
	if _, ok := d.GetOk("log_summary"); ok {
		jobLogSummaryAttr := d.Get("log_summary").([]interface{})
		if len(jobLogSummaryAttr) > 0 {
			logSummary := resourceIBMSchematicsJobMapToJobLogSummary(d.Get("log_summary.0").(map[string]interface{}))
			updateJobOptions.SetLogSummary(&logSummary)
		}
	}

	_, response, err := schematicsClient.UpdateJobWithContext(context, updateJobOptions)
	if err != nil {
		log.Printf("[DEBUG] UpdateJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("UpdateJobWithContext failed %s\n%s", err, response))
	}

	return resourceIBMSchematicsJobRead(context, d, meta)
}

func resourceIBMSchematicsJobDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	session, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return diag.FromErr(err)
	}

	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		return diag.FromErr(err)
	}
	jobIDSplit := strings.Split(d.Id(), ".")
	region := jobIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	deleteJobOptions := &schematicsv1.DeleteJobOptions{}

	iamRefreshToken := session.Config.IAMRefreshToken
	deleteJobOptions.SetRefreshToken(iamRefreshToken)

	deleteJobOptions.SetJobID(d.Id())

	response, err := schematicsClient.DeleteJobWithContext(context, deleteJobOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteJobWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
