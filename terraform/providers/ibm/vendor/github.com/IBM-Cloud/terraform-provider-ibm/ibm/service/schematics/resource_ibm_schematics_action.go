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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

const (
	actionName = "name"
)

func ResourceIBMSchematicsAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMSchematicsActionCreate,
		ReadContext:   resourceIBMSchematicsActionRead,
		UpdateContext: resourceIBMSchematicsActionUpdate,
		DeleteContext: resourceIBMSchematicsActionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The unique name of your action. The name can be up to 128 characters long and can include alphanumeric characters, spaces, dashes, and underscores. **Example** you can use the name to stop action.",
				ValidateFunc: validate.InvokeValidator("ibm_schematics_action", "name"),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action description.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "List of locations supported by IBM Cloud Schematics service.  While creating your workspace or action, choose the right region, since it cannot be changed.  Note, this does not limit the location of the IBM Cloud resources, provisioned using Schematics.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource-group name for an action.  By default, action is created in default resource group.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Action tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_state": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "User defined status of the Schematics object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User-defined states  * `draft` Object can be modified; can be used by Jobs run by the author, during execution  * `live` Object can be modified; can be used by Jobs during execution  * `locked` Object cannot be modified; can be used by Jobs during execution  * `disable` Object can be modified. cannot be used by Jobs during execution.",
						},
						"set_by": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the User who set the state of the Object.",
						},
						"set_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "When the User who set the state of the Object.",
						},
					},
				},
			},
			"source_readme_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the `README` file, for the source URL.",
			},
			"source": {
				Type:        schema.TypeList,
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
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "URL to the GIT Repo that can be used to clone the template.",
										ValidateFunc: validation.IsURLWithHTTPorHTTPS,
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
			"source_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_schematics_action", "source_type"),
				Description:  "Type of source for the Template.",
			},
			"command_parameter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Schematics job command parameter (playbook-name).",
			},
			"inventory": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Target inventory record ID, used by the action or ansible playbook.",
			},
			"credentials": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "credentials of the Action.",
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
			"bastion_credential": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "User editable variable data & system generated reference to value.",
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
			"targets_ini": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory of host and host group for the playbook in `INI` file format. For example, `\"targets_ini\": \"[webserverhost]  172.22.192.6  [dbhost]  172.22.192.5\"`. For more information, about an inventory host group syntax, see [Inventory host groups](https://cloud.ibm.com/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).",
			},
			"action_inputs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Input variables for the Action.",
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
			"action_outputs": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Output variables for the Action.",
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
				Description: "Environment variables for the Action.",
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
			"state": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed state of the Action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status of automation (workspace or action).",
						},
						"status_job_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Job id reference for this status.",
						},
						"status_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automation status message - to be displayed along with the status_code.",
						},
					},
				},
			},
			"sys_lock": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "System lock status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sys_locked": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Is the automation locked by a Schematic job ?.",
						},
						"sys_locked_by": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the User who performed the job, that lead to the locking of the automation.",
						},
						"sys_locked_at": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "When the User performed the job that lead to locking of the automation ?.",
						},
					},
				},
			},
			"x_github_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform template.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action Cloud Resource Name.",
			},
			"account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action account ID.",
			},
			"source_created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action Playbook Source creation time.",
			},
			"source_created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of user who created the Action Playbook Source.",
			},
			"source_updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action playbook updation time.",
			},
			"source_updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of user who updated the action playbook source.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action creation time.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who created an action.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action updation time.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who updated an action.",
			},
			"playbook_names": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Playbook names retrieved from the respository.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ResourceIBMSchematicsActionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "source_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "cos_bucket, external_scm, git_hub, git_hub_enterprise, git_lab, ibm_schematics_action_catalog, ibm_git_lab, local",
		},
		validate.ValidateSchema{
			Identifier:                 actionName,
			ValidateFunctionIdentifier: validate.StringLenBetween,
			Type:                       validate.TypeString,
			MinValueLength:             1,
			MaxValueLength:             65,
			Optional:                   true,
		})

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_schematics_action", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSchematicsActionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionCreate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_action", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if r, ok := d.GetOk("location"); ok {
		region := r.(string)
		schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
		if updatedURL {
			schematicsClient.Service.Options.URL = schematicsURL
		}
	}
	createActionOptions := &schematicsv1.CreateActionOptions{}

	if _, ok := d.GetOk("name"); ok {
		createActionOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("description"); ok {
		createActionOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("location"); ok {
		createActionOptions.SetLocation(d.Get("location").(string))
	}
	if _, ok := d.GetOk("resource_group"); ok {
		createActionOptions.SetResourceGroup(d.Get("resource_group").(string))
	}
	if _, ok := d.GetOk("tags"); ok {
		createActionOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("user_state"); ok {
		userState := resourceIBMSchematicsActionMapToUserState(d.Get("user_state.0").(map[string]interface{}))
		createActionOptions.SetUserState(&userState)
	}
	if _, ok := d.GetOk("source_readme_url"); ok {
		createActionOptions.SetSourceReadmeURL(d.Get("source_readme_url").(string))
	}
	if _, ok := d.GetOk("source"); ok {
		source := resourceIBMSchematicsActionMapToExternalSource(d.Get("source.0").(map[string]interface{}))
		createActionOptions.SetSource(&source)
	}
	if _, ok := d.GetOk("source_type"); ok {
		createActionOptions.SetSourceType(d.Get("source_type").(string))
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		createActionOptions.SetCommandParameter(d.Get("command_parameter").(string))
	}
	if _, ok := d.GetOk("inventory"); ok {
		createActionOptions.SetInventory(d.Get("inventory").(string))
	}
	if _, ok := d.GetOk("credentials"); ok {
		var credentials []schematicsv1.CredentialVariableData
		for _, e := range d.Get("credentials").([]interface{}) {
			value := e.(map[string]interface{})
			credentialsItem := resourceIBMSchematicsActionMapToCredentialsVariableData(value)
			credentials = append(credentials, credentialsItem)
		}
		createActionOptions.SetCredentials(credentials)
	}
	if _, ok := d.GetOk("bastion"); ok {
		bastion := resourceIBMSchematicsActionMapToBastionResourceDefinition(d.Get("bastion.0").(map[string]interface{}))
		createActionOptions.SetBastion(&bastion)
	}
	if _, ok := d.GetOk("bastion_credential"); ok {
		bastionCredential := resourceIBMSchematicsActionMapToCredentialsVariableData(d.Get("bastion_credential.0").(map[string]interface{}))
		createActionOptions.SetBastionCredential(&bastionCredential)
	}
	if _, ok := d.GetOk("targets_ini"); ok {
		createActionOptions.SetTargetsIni(d.Get("targets_ini").(string))
	}
	if _, ok := d.GetOk("action_inputs"); ok {
		var actionInputs []schematicsv1.VariableData
		for _, e := range d.Get("action_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			actionInputsItem := resourceIBMSchematicsActionMapToVariableData(value)
			actionInputs = append(actionInputs, actionInputsItem)
		}
		createActionOptions.SetInputs(actionInputs)
	}
	if _, ok := d.GetOk("action_outputs"); ok {
		var actionOutputs []schematicsv1.VariableData
		for _, e := range d.Get("action_outputs").([]interface{}) {
			value := e.(map[string]interface{})
			actionOutputsItem := resourceIBMSchematicsActionMapToVariableData(value)
			actionOutputs = append(actionOutputs, actionOutputsItem)
		}
		createActionOptions.SetOutputs(actionOutputs)
	}
	if _, ok := d.GetOk("settings"); ok {
		var settings []schematicsv1.VariableData
		for _, e := range d.Get("settings").([]interface{}) {
			value := e.(map[string]interface{})
			settingsItem := resourceIBMSchematicsActionMapToVariableData(value)
			settings = append(settings, settingsItem)
		}
		createActionOptions.SetSettings(settings)
	}
	if _, ok := d.GetOk("x_github_token"); ok {
		createActionOptions.SetXGithubToken(d.Get("x_github_token").(string))
	}

	action, response, err := schematicsClient.CreateActionWithContext(context, createActionOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionCreate CreateActionWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*action.ID)

	return resourceIBMSchematicsActionRead(context, d, meta)
}

func resourceIBMSchematicsActionMapToUserState(userStateMap map[string]interface{}) schematicsv1.UserState {
	userState := schematicsv1.UserState{}

	if userStateMap["state"] != nil {
		userState.State = core.StringPtr(userStateMap["state"].(string))
	}
	if userStateMap["set_by"] != nil {
		userState.SetBy = core.StringPtr(userStateMap["set_by"].(string))
	}
	if userStateMap["set_at"] != nil {
		setAt, err := strfmt.ParseDateTime(userStateMap["set_at"].(string))
		if err != nil {
			userState.SetAt = &setAt
		}
	}

	return userState
}

func resourceIBMSchematicsActionMapToExternalSource(externalSourceMap map[string]interface{}) schematicsv1.ExternalSource {
	externalSource := schematicsv1.ExternalSource{}

	externalSource.SourceType = core.StringPtr(externalSourceMap["source_type"].(string))
	if externalSourceMap["git"] != nil {
		externalSourceGit := resourceIBMSchematicsActionMapToExternalSourceGit(externalSourceMap["git"].([]interface{})[0].(map[string]interface{}))
		externalSource.Git = &externalSourceGit
	}
	if externalSourceMap["catalog"] != nil && len(externalSourceMap["catalog"].([]interface{})) > 0 {
		externalSourceCatalog := resourceIBMSchematicsActionMapToExternalSourceCatalog(externalSourceMap["catalog"].([]interface{})[0].(map[string]interface{}))
		externalSource.Catalog = &externalSourceCatalog
	}

	return externalSource
}

func resourceIBMSchematicsActionMapToExternalSourceGit(externalSourceGitMap map[string]interface{}) schematicsv1.GitSource {
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

func resourceIBMSchematicsActionMapToExternalSourceCatalog(externalSourceCatalogMap map[string]interface{}) schematicsv1.CatalogSource {
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

// func resourceIBMSchematicsActionMapToExternalSourceCosBucket(externalSourceCosBucketMap map[string]interface{}) schematicsv1.ExternalSourceCosBucket {
// 	externalSourceCosBucket := schematicsv1.ExternalSourceCosBucket{}

// 	if externalSourceCosBucketMap["cos_bucket_url"] != nil {
// 		externalSourceCosBucket.CosBucketURL = core.StringPtr(externalSourceCosBucketMap["cos_bucket_url"].(string))
// 	}

// 	return externalSourceCosBucket
// }

func resourceIBMSchematicsActionMapToVariableData(variableDataMap map[string]interface{}) schematicsv1.VariableData {
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
func resourceIBMSchematicsActionMapToCredentialsVariableData(variableDataMap map[string]interface{}) schematicsv1.CredentialVariableData {
	variableData := schematicsv1.CredentialVariableData{}

	if variableDataMap["name"] != nil {
		variableData.Name = core.StringPtr(variableDataMap["name"].(string))
	}
	if variableDataMap["value"] != nil {
		variableData.Value = core.StringPtr(variableDataMap["value"].(string))
	}
	if variableDataMap["metadata"] != nil && len(variableDataMap["metadata"].([]interface{})) != 0 {
		variableMetaData := resourceIBMSchematicsJobMapToCredentialVariableMetadata(variableDataMap["metadata"].([]interface{})[0].(map[string]interface{}))
		variableData.Metadata = &variableMetaData
	}
	if variableDataMap["link"] != nil {
		variableData.Link = core.StringPtr(variableDataMap["link"].(string))
	}

	return variableData
}

func resourceIBMSchematicsActionMapToVariableMetadata(variableMetadataMap map[string]interface{}) schematicsv1.VariableMetadata {
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

func resourceIBMSchematicsActionMapToBastionResourceDefinition(bastionResourceDefinitionMap map[string]interface{}) schematicsv1.BastionResourceDefinition {
	bastionResourceDefinition := schematicsv1.BastionResourceDefinition{}

	if bastionResourceDefinitionMap["name"] != nil {
		bastionResourceDefinition.Name = core.StringPtr(bastionResourceDefinitionMap["name"].(string))
	}
	if bastionResourceDefinitionMap["host"] != nil {
		bastionResourceDefinition.Host = core.StringPtr(bastionResourceDefinitionMap["host"].(string))
	}

	return bastionResourceDefinition
}

func resourceIBMSchematicsActionMapToActionState(actionStateMap map[string]interface{}) schematicsv1.ActionState {
	actionState := schematicsv1.ActionState{}

	if actionStateMap["status_code"] != nil {
		actionState.StatusCode = core.StringPtr(actionStateMap["status_code"].(string))
	}
	if actionStateMap["status_job_id"] != nil {
		actionState.StatusJobID = core.StringPtr(actionStateMap["status_job_id"].(string))
	}
	if actionStateMap["status_message"] != nil {
		actionState.StatusMessage = core.StringPtr(actionStateMap["status_message"].(string))
	}

	return actionState
}

func resourceIBMSchematicsActionMapToSystemLock(systemLockMap map[string]interface{}) schematicsv1.SystemLock {
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

func resourceIBMSchematicsActionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	actionIDSplit := strings.Split(d.Id(), ".")
	region := actionIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	getActionOptions := &schematicsv1.GetActionOptions{}

	getActionOptions.SetActionID(d.Id())

	action, response, err := schematicsClient.GetActionWithContext(context, getActionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead GetActionWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("name", action.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("description", action.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("location", action.Location); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("resource_group", action.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if action.Tags != nil {
		if err = d.Set("tags", action.Tags); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if action.UserState != nil {
		userStateMap := resourceIBMSchematicsActionUserStateToMap(*action.UserState)
		if err = d.Set("user_state", []map[string]interface{}{userStateMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("source_readme_url", action.SourceReadmeURL); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if _, ok := d.GetOk("source"); ok {
		if action.Source != nil {
			sourceMap := resourceIBMSchematicsActionExternalSourceToMap(*action.Source)
			if err = d.Set("source", []map[string]interface{}{sourceMap}); err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
	}
	if err = d.Set("source_type", action.SourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("command_parameter", action.CommandParameter); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("inventory", action.Inventory); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if action.Credentials != nil {
		credentials := []map[string]interface{}{}
		for _, credentialsItem := range action.Credentials {
			credentialsItemMap := resourceIBMSchematicsActionCredentialVariableDataToMap(credentialsItem)
			credentials = append(credentials, credentialsItemMap)
		}
		if err = d.Set("credentials", credentials); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if _, ok := d.GetOk("bastion"); ok {
		if action.Bastion != nil {
			bastionMap := resourceIBMSchematicsActionBastionResourceDefinitionToMap(*action.Bastion)
			if err = d.Set("bastion", []map[string]interface{}{bastionMap}); err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
	}
	if action.BastionCredential != nil {
		bastionCredentialMap := resourceIBMSchematicsActionCredentialVariableDataToMap(*action.BastionCredential)
		if err = d.Set("bastion_credential", []map[string]interface{}{bastionCredentialMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("targets_ini", action.TargetsIni); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if action.Inputs != nil {
		actionInputs := []map[string]interface{}{}
		for _, actionInputsItem := range action.Inputs {
			actionInputsItemMap := resourceIBMSchematicsActionVariableDataToMap(actionInputsItem)
			actionInputs = append(actionInputs, actionInputsItemMap)
		}
		if err = d.Set("action_inputs", actionInputs); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if action.Outputs != nil {
		actionOutputs := []map[string]interface{}{}
		for _, actionOutputsItem := range action.Outputs {
			actionOutputsItemMap := resourceIBMSchematicsActionVariableDataToMap(actionOutputsItem)
			actionOutputs = append(actionOutputs, actionOutputsItemMap)
		}
		if err = d.Set("action_outputs", actionOutputs); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if action.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range action.Settings {
			settingsItemMap := resourceIBMSchematicsActionVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
		}
		if err = d.Set("settings", settings); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if action.State != nil {
		stateMap := resourceIBMSchematicsActionActionStateToMap(*action.State)
		if err = d.Set("state", []map[string]interface{}{stateMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if action.SysLock != nil {
		sysLockMap := resourceIBMSchematicsActionSystemLockToMap(*action.SysLock)
		if err = d.Set("sys_lock", []map[string]interface{}{sysLockMap}); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("crn", action.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("account", action.Account); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("source_created_at", flex.DateTimeToString(action.SourceCreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("source_created_by", action.SourceCreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("source_updated_at", flex.DateTimeToString(action.SourceUpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("source_updated_by", action.SourceUpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(action.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", action.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(action.UpdatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("updated_by", action.UpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if action.PlaybookNames != nil && len(action.PlaybookNames) > 0 {
		if err = d.Set("playbook_names", action.PlaybookNames); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionRead failed with error: %s", err), "ibm_schematics_action", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	} else {
		d.Set("playbook_names", []string{})
	}

	return nil
}

func resourceIBMSchematicsActionUserStateToMap(userState schematicsv1.UserState) map[string]interface{} {
	userStateMap := map[string]interface{}{}

	if userState.State != nil {
		userStateMap["state"] = userState.State
	}
	if userState.SetBy != nil {
		userStateMap["set_by"] = userState.SetBy
	}
	if userState.SetAt != nil {
		userStateMap["set_at"] = userState.SetAt.String()
	}

	return userStateMap
}

func resourceIBMSchematicsActionExternalSourceToMap(externalSource schematicsv1.ExternalSource) map[string]interface{} {
	externalSourceMap := map[string]interface{}{}

	externalSourceMap["source_type"] = externalSource.SourceType
	if externalSource.Git != nil {
		GitMap := resourceIBMSchematicsActionExternalSourceGitToMap(*externalSource.Git)
		externalSourceMap["git"] = []map[string]interface{}{GitMap}
	}
	if externalSource.Catalog != nil {
		CatalogMap := resourceIBMSchematicsActionExternalSourceCatalogToMap(*externalSource.Catalog)
		externalSourceMap["catalog"] = []map[string]interface{}{CatalogMap}
	}
	// if externalSource.CosBucket != nil {
	// 	CosBucketMap := resourceIBMSchematicsActionExternalSourceCosBucketToMap(*externalSource.CosBucket)
	// 	externalSourceMap["cos_bucket"] = []map[string]interface{}{CosBucketMap}
	// }

	return externalSourceMap
}

func resourceIBMSchematicsActionExternalSourceGitToMap(externalSourceGit schematicsv1.GitSource) map[string]interface{} {
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

func resourceIBMSchematicsActionExternalSourceCatalogToMap(externalSourceCatalog schematicsv1.CatalogSource) map[string]interface{} {
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

// func resourceIBMSchematicsActionExternalSourceCosBucketToMap(externalSourceCosBucket schematicsv1.ExternalSourceCosBucket) map[string]interface{} {
// 	externalSourceCosBucketMap := map[string]interface{}{}

// 	if externalSourceCosBucket.CosBucketURL != nil {
// 		externalSourceCosBucketMap["cos_bucket_url"] = externalSourceCosBucket.CosBucketURL
// 	}

// 	return externalSourceCosBucketMap
// }

func resourceIBMSchematicsActionVariableDataToMap(variableData schematicsv1.VariableData) map[string]interface{} {
	variableDataMap := map[string]interface{}{}

	if variableData.Name != nil {
		variableDataMap["name"] = variableData.Name
	}
	if variableData.Value != nil {
		variableDataMap["value"] = variableData.Value
	}
	if variableData.Metadata != nil {
		MetadataMap := resourceIBMSchematicsActionVariableMetadataToMap(*variableData.Metadata)
		variableDataMap["metadata"] = []map[string]interface{}{MetadataMap}
	}
	if variableData.Link != nil {
		variableDataMap["link"] = variableData.Link
	}

	return variableDataMap
}
func resourceIBMSchematicsActionCredentialVariableDataToMap(variableData schematicsv1.CredentialVariableData) map[string]interface{} {
	variableDataMap := map[string]interface{}{}

	if variableData.Name != nil {
		variableDataMap["name"] = variableData.Name
	}
	if variableData.Value != nil {
		variableDataMap["value"] = variableData.Value
	}
	if variableData.Metadata != nil {
		MetadataMap := resourceIBMSchematicsActionCredentialVariableMetadataToMap(*variableData.Metadata)
		variableDataMap["metadata"] = []map[string]interface{}{MetadataMap}
	}
	if variableData.Link != nil {
		variableDataMap["link"] = variableData.Link
	}

	return variableDataMap
}

func resourceIBMSchematicsActionCredentialVariableMetadataToMap(variableMetadata schematicsv1.CredentialVariableMetadata) map[string]interface{} {
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
	if variableMetadata.Immutable != nil {
		variableMetadataMap["immutable"] = variableMetadata.Immutable
	}
	if variableMetadata.Hidden != nil {
		variableMetadataMap["hidden"] = variableMetadata.Hidden
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
func resourceIBMSchematicsActionVariableMetadataToMap(variableMetadata schematicsv1.VariableMetadata) map[string]interface{} {
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

func resourceIBMSchematicsActionBastionResourceDefinitionToMap(bastionResourceDefinition schematicsv1.BastionResourceDefinition) map[string]interface{} {
	bastionResourceDefinitionMap := map[string]interface{}{}

	if bastionResourceDefinition.Name != nil {
		bastionResourceDefinitionMap["name"] = bastionResourceDefinition.Name
	}
	if bastionResourceDefinition.Host != nil {
		bastionResourceDefinitionMap["host"] = bastionResourceDefinition.Host
	}

	return bastionResourceDefinitionMap
}

func resourceIBMSchematicsActionActionStateToMap(actionState schematicsv1.ActionState) map[string]interface{} {
	actionStateMap := map[string]interface{}{}

	if actionState.StatusCode != nil {
		actionStateMap["status_code"] = actionState.StatusCode
	}
	if actionState.StatusJobID != nil {
		actionStateMap["status_job_id"] = actionState.StatusJobID
	}
	if actionState.StatusMessage != nil {
		actionStateMap["status_message"] = actionState.StatusMessage
	}

	return actionStateMap
}

func resourceIBMSchematicsActionSystemLockToMap(systemLock schematicsv1.SystemLock) map[string]interface{} {
	systemLockMap := map[string]interface{}{}

	if systemLock.SysLocked != nil {
		systemLockMap["sys_locked"] = systemLock.SysLocked
	}
	if systemLock.SysLockedBy != nil {
		systemLockMap["sys_locked_by"] = systemLock.SysLockedBy
	}
	if systemLock.SysLockedAt != nil {
		systemLockMap["sys_locked_at"] = systemLock.SysLockedAt.String()
	}

	return systemLockMap
}

func resourceIBMSchematicsActionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionUpdate schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_action", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	actionIDSplit := strings.Split(d.Id(), ".")
	region := actionIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}

	updateActionOptions := &schematicsv1.UpdateActionOptions{}

	updateActionOptions.SetActionID(d.Id())

	hasChange := false

	if d.HasChange("name") {
		updateActionOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateActionOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("location") {
		updateActionOptions.SetLocation(d.Get("location").(string))
		hasChange = true
	}
	if d.HasChange("resource_group") {
		updateActionOptions.SetResourceGroup(d.Get("resource_group").(string))
		hasChange = true
	}
	if d.HasChange("tags") {
		updateActionOptions.SetTags(flex.ExpandStringList(d.Get("tags").([]interface{})))
		hasChange = true
	}
	if d.HasChange("user_state") {
		userStateAttr := d.Get("user_state").([]interface{})
		if len(userStateAttr) > 0 {
			userState := resourceIBMSchematicsActionMapToUserState(d.Get("user_state.0").(map[string]interface{}))
			updateActionOptions.SetUserState(&userState)
			hasChange = true
		}
	}
	if d.HasChange("source_readme_url") {
		updateActionOptions.SetSourceReadmeURL(d.Get("source_readme_url").(string))
		hasChange = true
	}
	if d.HasChange("source") {
		sourceAttr := d.Get("source").([]interface{})
		if len(sourceAttr) > 0 {
			source := resourceIBMSchematicsActionMapToExternalSource(d.Get("source.0").(map[string]interface{}))
			updateActionOptions.SetSource(&source)
			hasChange = true
		}
	}
	if d.HasChange("source_type") {
		updateActionOptions.SetSourceType(d.Get("source_type").(string))
		hasChange = true
	}
	if d.HasChange("command_parameter") {
		updateActionOptions.SetCommandParameter(d.Get("command_parameter").(string))
		hasChange = true
	}
	if d.HasChange("inventory") {
		updateActionOptions.SetInventory(d.Get("inventory").(string))
		hasChange = true
	}
	if d.HasChange("credentials") {
		// TODO: handle Credentials of type TypeList -- not primitive, not model
		hasChange = true
	}
	if d.HasChange("bastion") {
		bastionAttr := d.Get("bastion").([]interface{})
		if len(bastionAttr) > 0 {
			bastion := resourceIBMSchematicsActionMapToBastionResourceDefinition(d.Get("bastion.0").(map[string]interface{}))
			updateActionOptions.SetBastion(&bastion)
			hasChange = true
		}
	}
	if d.HasChange("inventory") {
		updateActionOptions.SetInventory(d.Get("inventory").(string))
		hasChange = true
	}
	if d.HasChange("bastion_credential") {
		bastionCredential := resourceIBMSchematicsActionMapToCredentialsVariableData(d.Get("bastion_credential.0").(map[string]interface{}))
		updateActionOptions.SetBastionCredential(&bastionCredential)
		hasChange = true
	}
	if d.HasChange("targets_ini") {
		updateActionOptions.SetTargetsIni(d.Get("targets_ini").(string))
		hasChange = true
	}
	if d.HasChange("action_inputs") {
		var inputs []schematicsv1.VariableData
		for _, e := range d.Get("action_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			inputsItem := resourceIBMSchematicsActionMapToVariableData(value)
			inputs = append(inputs, inputsItem)
		}
		updateActionOptions.SetInputs(inputs)
		hasChange = true
	}
	if d.HasChange("action_outputs") {
		var outputs []schematicsv1.VariableData
		for _, e := range d.Get("action_outputs").([]interface{}) {
			value := e.(map[string]interface{})
			outputsItem := resourceIBMSchematicsActionMapToVariableData(value)
			outputs = append(outputs, outputsItem)
		}
		updateActionOptions.SetOutputs(outputs)
		hasChange = true
	}
	if d.HasChange("settings") {
		var settings []schematicsv1.VariableData
		for _, e := range d.Get("settings").([]interface{}) {
			value := e.(map[string]interface{})
			settingsItem := resourceIBMSchematicsActionMapToVariableData(value)
			settings = append(settings, settingsItem)
		}
		updateActionOptions.SetSettings(settings)
		hasChange = true
	}
	if hasChange {
		_, response, err := schematicsClient.UpdateActionWithContext(context, updateActionOptions)
		if err != nil {

			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionUpdate failed with error: %s and response:\n%s", err, response), "ibm_schematics_action", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMSchematicsActionRead(context, d, meta)
}

func resourceIBMSchematicsActionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionDelete schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_action", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	actionIDSplit := strings.Split(d.Id(), ".")
	region := actionIDSplit[0]
	schematicsURL, updatedURL, _ := SchematicsEndpointURL(region, meta)
	if updatedURL {
		schematicsClient.Service.Options.URL = schematicsURL
	}
	deleteActionOptions := &schematicsv1.DeleteActionOptions{}

	deleteActionOptions.SetActionID(d.Id())

	response, err := schematicsClient.DeleteActionWithContext(context, deleteActionOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("resourceIBMSchematicsActionDelete DeleteActionWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_action", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
