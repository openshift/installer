// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

const (
	actionName = "name"
)

func resourceIBMSchematicsAction() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSchematicsActionCreate,
		Read:     resourceIBMSchematicsActionRead,
		Update:   resourceIBMSchematicsActionUpdate,
		Delete:   resourceIBMSchematicsActionDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Action name (unique for an account).",
				ValidateFunc: InvokeValidator("ibm_schematics_action", actionName),
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Action description.",
			},
			"location": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_schematics_action", "location"),
				Description:  "List of action locations supported by IBM Cloud Schematics service.  **Note** this does not limit the location of the resources provisioned using Schematics.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Resource-group name for an action.  By default, action is created in default resource group.",
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Action tags.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_state": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "User defined status of the Schematics object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "User defined states  * `draft` Object can be modified, and can be used by jobs run by an author, during execution  * `live` Object can be modified, and can be used by jobs during execution  * `locked` Object cannot be modified, and can be used by jobs during execution  * `disable` Object can be modified, and cannot be used by Jobs during execution.",
						},
						"set_by": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name of the user who set the state of an Object.",
						},
						"set_at": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "When the user who set the state of an Object.",
						},
					},
				},
			},
			"source_readme_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the `README` file, for the source.",
			},
			"source": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Source of templates, playbooks, or controls.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of source for the Template.",
						},
						"git": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Connection details to Git source.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"git_repo_url": &schema.Schema{
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "URL to the GIT Repo that can be used to clone the template.",
										ValidateFunc: validation.IsURLWithHTTPorHTTPS,
									},
									"git_token": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Personal Access Token to connect to Git URLs.",
									},
									"git_repo_folder": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the folder in the Git Repo, that contains the template.",
									},
									"git_release": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the release tag, used to fetch the Git Repo.",
									},
									"git_branch": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Name of the branch, used to fetch the Git Repo.",
									},
								},
							},
						},
					},
				},
			},
			"source_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_schematics_action", "source_type"),
				Description:  "Type of source for the Template.",
			},
			"command_parameter": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Schematics job command parameter (playbook-name, capsule-name or flow-name).",
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
			"targets_ini": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Inventory of host and host group for the playbook in `INI` file format. For example, `\"targets_ini\": \"[webserverhost]  172.22.192.6  [dbhost]  172.22.192.5\"`. For more information, about an inventory host group syntax, see [Inventory host groups](/docs/schematics?topic=schematics-schematics-cli-reference#schematics-inventory-host-grps).",
			},
			"credentials": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "credentials of the Action.",
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
			"action_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Input variables for an action.",
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
			"action_outputs": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Output variables for an action.",
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
				Description: "Environment variables for an action.",
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
			"trigger_record_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID to the trigger.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Computed state of an action.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status_code": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Status of automation (workspace or action).",
						},
						"status_job_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Job id reference for this status.",
						},
						"status_message": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Automation status message - to be displayed along with the status_code.",
						},
					},
				},
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
			"x_github_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform template.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action Cloud Resource Name.",
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action account ID.",
			},
			"source_created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action Playbook Source creation time.",
			},
			"source_created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of user who created the Action Playbook Source.",
			},
			"source_updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The action playbook updation time.",
			},
			"source_updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of user who updated the action playbook source.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action creation time.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who created an action.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Action updation time.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "E-mail address of the user who updated an action.",
			},
			"namespace": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the namespace.",
			},
			"playbook_names": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Playbook names retrieved from the respository.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceIBMSchematicsActionValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              "eu-de, eu-gb, us-east, us-south",
		},
		ValidateSchema{
			Identifier:                 "source_type",
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              "external_scm, git_hub, git_hub_enterprise, git_lab, ibm_cloud_catalog, ibm_git_lab, local",
		},
		ValidateSchema{
			Identifier:                 actionName,
			ValidateFunctionIdentifier: StringLenBetween,
			Type:                       TypeString,
			MinValueLength:             1,
			MaxValueLength:             65,
			Optional:                   true,
		})

	resourceValidator := ResourceValidator{ResourceName: "ibm_schematics_action", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMSchematicsActionCreate(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
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
		createActionOptions.SetTags(expandStringList(d.Get("tags").([]interface{})))
	}
	if _, ok := d.GetOk("user_state"); ok {
		userStateAttr := d.Get("user_state").([]interface{})
		if len(userStateAttr) > 0 {
			userState := resourceIBMSchematicsActionMapToUserState(d.Get("user_state.0").(map[string]interface{}))
			createActionOptions.SetUserState(&userState)
		}
	}
	if _, ok := d.GetOk("source_readme_url"); ok {
		createActionOptions.SetSourceReadmeURL(d.Get("source_readme_url").(string))
	}
	if _, ok := d.GetOk("source"); ok {
		sourceAttr := d.Get("source").([]interface{})
		if len(sourceAttr) > 0 {
			source := resourceIBMSchematicsActionMapToExternalSource(d.Get("source.0").(map[string]interface{}))
			createActionOptions.SetSource(&source)
		}
	}
	if _, ok := d.GetOk("source_type"); ok {
		createActionOptions.SetSourceType(d.Get("source_type").(string))
	}
	if _, ok := d.GetOk("command_parameter"); ok {
		createActionOptions.SetCommandParameter(d.Get("command_parameter").(string))
	}
	if _, ok := d.GetOk("bastion"); ok {
		bastionAttr := d.Get("bastion").([]interface{})
		if len(bastionAttr) > 0 {
			bastion := resourceIBMSchematicsActionMapToTargetResourceset(d.Get("bastion.0").(map[string]interface{}))
			createActionOptions.SetBastion(&bastion)
		}
	}
	if _, ok := d.GetOk("targets_ini"); ok {
		createActionOptions.SetTargetsIni(d.Get("targets_ini").(string))
	}
	if _, ok := d.GetOk("credentials"); ok {
		var credentials []schematicsv1.VariableData
		for _, e := range d.Get("credentials").([]interface{}) {
			value := e.(map[string]interface{})
			credentialsItem := resourceIBMSchematicsActionMapToVariableData(value)
			credentials = append(credentials, credentialsItem)
		}
		createActionOptions.SetCredentials(credentials)
	}
	if _, ok := d.GetOk("action_inputs"); ok {
		var inputs []schematicsv1.VariableData
		for _, e := range d.Get("action_inputs").([]interface{}) {
			value := e.(map[string]interface{})
			inputsItem := resourceIBMSchematicsActionMapToVariableData(value)
			inputs = append(inputs, inputsItem)
		}
		createActionOptions.SetInputs(inputs)
	}
	if _, ok := d.GetOk("action_outputs"); ok {
		var outputs []schematicsv1.VariableData
		for _, e := range d.Get("action_outputs").([]interface{}) {
			value := e.(map[string]interface{})
			outputsItem := resourceIBMSchematicsActionMapToVariableData(value)
			outputs = append(outputs, outputsItem)
		}
		createActionOptions.SetOutputs(outputs)
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
	if _, ok := d.GetOk("trigger_record_id"); ok {
		createActionOptions.SetTriggerRecordID(d.Get("trigger_record_id").(string))
	}
	if _, ok := d.GetOk("state"); ok {
		stateAttr := d.Get("state").([]interface{})
		if len(stateAttr) > 0 {
			state := resourceIBMSchematicsActionMapToActionState(d.Get("state.0").(map[string]interface{}))
			createActionOptions.SetState(&state)
		}
	}
	if _, ok := d.GetOk("sys_lock"); ok {
		sysLockAttr := d.Get("sys_lock").([]interface{})
		if len(sysLockAttr) > 0 {
			sysLock := resourceIBMSchematicsActionMapToSystemLock(d.Get("sys_lock.0").(map[string]interface{}))
			createActionOptions.SetSysLock(&sysLock)
		}
	}
	if _, ok := d.GetOk("x_github_token"); ok {
		createActionOptions.SetXGithubToken(d.Get("x_github_token").(string))
	}

	action, response, err := schematicsClient.CreateActionWithContext(context.TODO(), createActionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateActionWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*action.ID)

	return resourceIBMSchematicsActionRead(d, meta)
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

	return externalSource
}

func resourceIBMSchematicsActionMapToExternalSourceGit(externalSourceGitMap map[string]interface{}) schematicsv1.ExternalSourceGit {
	externalSourceGit := schematicsv1.ExternalSourceGit{}

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

func resourceIBMSchematicsActionMapToTargetResourceset(targetResourcesetMap map[string]interface{}) schematicsv1.TargetResourceset {
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
		createdAt, err := strfmt.ParseDateTime(targetResourcesetMap["created_at"].(string))
		if err != nil {
			targetResourceset.CreatedAt = &createdAt
		}
	}
	if targetResourcesetMap["created_by"] != nil {
		targetResourceset.CreatedBy = core.StringPtr(targetResourcesetMap["created_by"].(string))
	}
	if targetResourcesetMap["updated_at"] != nil {
		updatedAt, err := strfmt.ParseDateTime(targetResourcesetMap["updated_at"].(string))
		if err != nil {
			targetResourceset.CreatedAt = &updatedAt
		}
	}
	if targetResourcesetMap["updated_by"] != nil {
		targetResourceset.UpdatedBy = core.StringPtr(targetResourcesetMap["updated_by"].(string))
	}
	if targetResourcesetMap["sys_lock"] != nil && len(targetResourcesetMap["sys_lock"].([]interface{})) != 0 {
		sysLock := resourceIBMSchematicsActionMapToSystemLock(targetResourcesetMap["sys_lock"].([]interface{})[0].(map[string]interface{}))
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

func resourceIBMSchematicsActionMapToSystemLock(systemLockMap map[string]interface{}) schematicsv1.SystemLock {
	systemLock := schematicsv1.SystemLock{}

	if systemLockMap["sys_locked"] != nil {
		systemLock.SysLocked = core.BoolPtr(systemLockMap["sys_locked"].(bool))
	}
	if systemLockMap["sys_locked_by"] != nil {
		systemLock.SysLockedBy = core.StringPtr(systemLockMap["sys_locked_by"].(string))
	}
	if systemLockMap["sys_locked_at"] != nil {
		sysLockedAt, err := strfmt.ParseDateTime(systemLockMap["sys_locked_at"].(string))
		if err != nil {
			systemLock.SysLockedAt = &sysLockedAt
		}
	}

	return systemLock
}

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

func resourceIBMSchematicsActionRead(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	getActionOptions := &schematicsv1.GetActionOptions{}

	getActionOptions.SetActionID(d.Id())

	action, response, err := schematicsClient.GetActionWithContext(context.TODO(), getActionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetActionWithContext failed %s\n%s", err, response)
		return err
	}

	if err = d.Set("name", action.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("description", action.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("location", action.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err = d.Set("resource_group", action.ResourceGroup); err != nil {
		return fmt.Errorf("Error setting resource_group: %s", err)
	}
	if action.Tags != nil {
		if err = d.Set("tags", action.Tags); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}
	if action.UserState != nil {
		userStateMap := resourceIBMSchematicsActionUserStateToMap(*action.UserState)
		if err = d.Set("user_state", []map[string]interface{}{userStateMap}); err != nil {
			return fmt.Errorf("Error setting user_state: %s", err)
		}
	}
	if err = d.Set("source_readme_url", action.SourceReadmeURL); err != nil {
		return fmt.Errorf("Error setting source_readme_url: %s", err)
	}
	if _, ok := d.GetOk("source"); ok {
		if action.Source != nil {
			sourceMap := resourceIBMSchematicsActionExternalSourceToMap(*action.Source)
			if err = d.Set("source", []map[string]interface{}{sourceMap}); err != nil {
				return fmt.Errorf("Error setting source: %s", err)
			}
		}
	}
	if err = d.Set("source_type", action.SourceType); err != nil {
		return fmt.Errorf("Error setting source_type: %s", err)
	}
	if err = d.Set("command_parameter", action.CommandParameter); err != nil {
		return fmt.Errorf("Error setting command_parameter: %s", err)
	}
	if _, ok := d.GetOk("bastion"); ok {
		if action.Bastion != nil {
			bastionMap := resourceIBMSchematicsActionTargetResourcesetToMap(*action.Bastion)
			if err = d.Set("bastion", []map[string]interface{}{bastionMap}); err != nil {
				return fmt.Errorf("Error setting bastion: %s", err)
			}
		}
	}
	if err = d.Set("targets_ini", action.TargetsIni); err != nil {
		return fmt.Errorf("Error setting targets_ini: %s", err)
	}
	if action.Credentials != nil {
		credentials := []map[string]interface{}{}
		for _, credentialsItem := range action.Credentials {
			credentialsItemMap := resourceIBMSchematicsActionVariableDataToMap(credentialsItem)
			credentials = append(credentials, credentialsItemMap)
		}
		if err = d.Set("credentials", credentials); err != nil {
			return fmt.Errorf("Error setting credentials: %s", err)
		}
	}
	if action.Inputs != nil {
		inputs := []map[string]interface{}{}
		for _, inputsItem := range action.Inputs {
			inputsItemMap := resourceIBMSchematicsActionVariableDataToMap(inputsItem)
			inputs = append(inputs, inputsItemMap)
		}
		if err = d.Set("action_inputs", inputs); err != nil {
			return fmt.Errorf("Error setting action_inputs: %s", err)
		}
	}
	if action.Outputs != nil {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range action.Outputs {
			outputsItemMap := resourceIBMSchematicsActionVariableDataToMap(outputsItem)
			outputs = append(outputs, outputsItemMap)
		}
		if err = d.Set("action_outputs", outputs); err != nil {
			return fmt.Errorf("Error setting action_outputs: %s", err)
		}
	}
	if action.Settings != nil {
		settings := []map[string]interface{}{}
		for _, settingsItem := range action.Settings {
			settingsItemMap := resourceIBMSchematicsActionVariableDataToMap(settingsItem)
			settings = append(settings, settingsItemMap)
		}
		if err = d.Set("settings", settings); err != nil {
			return fmt.Errorf("Error setting settings: %s", err)
		}
	}
	if err = d.Set("trigger_record_id", action.TriggerRecordID); err != nil {
		return fmt.Errorf("Error setting trigger_record_id: %s", err)
	}
	if action.State != nil {
		stateMap := resourceIBMSchematicsActionActionStateToMap(*action.State)
		if err = d.Set("state", []map[string]interface{}{stateMap}); err != nil {
			return fmt.Errorf("Error setting state: %s", err)
		}
	}
	if action.SysLock != nil {
		sysLockMap := resourceIBMSchematicsActionSystemLockToMap(*action.SysLock)
		if err = d.Set("sys_lock", []map[string]interface{}{sysLockMap}); err != nil {
			return fmt.Errorf("Error setting sys_lock: %s", err)
		}
	}
	if err = d.Set("crn", action.Crn); err != nil {
		return fmt.Errorf("Error setting crn: %s", err)
	}
	if err = d.Set("account", action.Account); err != nil {
		return fmt.Errorf("Error setting account: %s", err)
	}
	if action.SourceCreatedAt != nil {
		if err = d.Set("source_created_at", action.SourceCreatedAt.String()); err != nil {
			return fmt.Errorf("Error setting source_created_at: %s", err)
		}
	}
	if err = d.Set("source_created_by", action.SourceCreatedBy); err != nil {
		return fmt.Errorf("Error setting source_created_by: %s", err)
	}
	if action.SourceUpdatedAt != nil {
		if err = d.Set("source_updated_at", action.SourceUpdatedAt.String()); err != nil {
			return fmt.Errorf("Error setting source_updated_at: %s", err)
		}
	}
	if err = d.Set("source_updated_by", action.SourceUpdatedBy); err != nil {
		return fmt.Errorf("Error setting source_updated_by: %s", err)
	}
	if action.CreatedAt != nil {
		if err = d.Set("created_at", action.CreatedAt.String()); err != nil {
			return fmt.Errorf("Error setting created_at: %s", err)
		}
	}
	if err = d.Set("created_by", action.CreatedBy); err != nil {
		return fmt.Errorf("Error setting created_by: %s", err)
	}
	if action.UpdatedAt != nil {
		if err = d.Set("updated_at", action.UpdatedAt.String()); err != nil {
			return fmt.Errorf("Error setting updated_at: %s", err)
		}
	}
	if err = d.Set("updated_by", action.UpdatedBy); err != nil {
		return fmt.Errorf("Error setting updated_by: %s", err)
	}
	if err = d.Set("namespace", action.Namespace); err != nil {
		return fmt.Errorf("Error setting namespace: %s", err)
	}
	if action.PlaybookNames != nil && len(action.PlaybookNames) > 0 {
		if err = d.Set("playbook_names", action.PlaybookNames); err != nil {
			return fmt.Errorf("Error setting playbook_names: %s", err)
		}
	} else {
		d.Set("playbook_names", []string{})
	}

	return nil
}

func resourceIBMSchematicsActionUserStateToMap(userState schematicsv1.UserState) map[string]interface{} {
	userStateMap := map[string]interface{}{}

	userStateMap["state"] = userState.State
	userStateMap["set_by"] = userState.SetBy
	userStateMap["set_at"] = userState.SetAt.String()

	return userStateMap
}

func resourceIBMSchematicsActionExternalSourceToMap(externalSource schematicsv1.ExternalSource) map[string]interface{} {
	externalSourceMap := map[string]interface{}{}

	externalSourceMap["source_type"] = externalSource.SourceType
	if externalSource.Git != nil {
		GitMap := resourceIBMSchematicsActionExternalSourceGitToMap(*externalSource.Git)
		externalSourceMap["git"] = []map[string]interface{}{GitMap}
	}

	return externalSourceMap
}

func resourceIBMSchematicsActionExternalSourceGitToMap(externalSourceGit schematicsv1.ExternalSourceGit) map[string]interface{} {
	externalSourceGitMap := map[string]interface{}{}

	externalSourceGitMap["git_repo_url"] = externalSourceGit.GitRepoURL
	externalSourceGitMap["git_token"] = externalSourceGit.GitToken
	externalSourceGitMap["git_repo_folder"] = externalSourceGit.GitRepoFolder
	externalSourceGitMap["git_release"] = externalSourceGit.GitRelease
	externalSourceGitMap["git_branch"] = externalSourceGit.GitBranch

	return externalSourceGitMap
}

func resourceIBMSchematicsActionTargetResourcesetToMap(targetResourceset schematicsv1.TargetResourceset) map[string]interface{} {
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
		SysLockMap := resourceIBMSchematicsActionSystemLockToMap(*targetResourceset.SysLock)
		targetResourcesetMap["sys_lock"] = []map[string]interface{}{SysLockMap}
	}
	if targetResourceset.ResourceIds != nil {
		targetResourcesetMap["resource_ids"] = targetResourceset.ResourceIds
	}

	return targetResourcesetMap
}

func resourceIBMSchematicsActionSystemLockToMap(systemLock schematicsv1.SystemLock) map[string]interface{} {
	systemLockMap := map[string]interface{}{}

	systemLockMap["sys_locked"] = systemLock.SysLocked
	systemLockMap["sys_locked_by"] = systemLock.SysLockedBy
	systemLockMap["sys_locked_at"] = systemLock.SysLockedAt.String()

	return systemLockMap
}

func resourceIBMSchematicsActionVariableDataToMap(variableData schematicsv1.VariableData) map[string]interface{} {
	variableDataMap := map[string]interface{}{}

	variableDataMap["name"] = variableData.Name
	variableDataMap["value"] = variableData.Value
	if variableData.Metadata != nil {
		MetadataMap := resourceIBMSchematicsActionVariableMetadataToMap(*variableData.Metadata)
		variableDataMap["metadata"] = []map[string]interface{}{MetadataMap}
	}
	variableDataMap["link"] = variableData.Link

	return variableDataMap
}

func resourceIBMSchematicsActionVariableMetadataToMap(variableMetadata schematicsv1.VariableMetadata) map[string]interface{} {
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

func resourceIBMSchematicsActionActionStateToMap(actionState schematicsv1.ActionState) map[string]interface{} {
	actionStateMap := map[string]interface{}{}

	actionStateMap["status_code"] = actionState.StatusCode
	actionStateMap["status_job_id"] = actionState.StatusJobID
	actionStateMap["status_message"] = actionState.StatusMessage

	return actionStateMap
}

func resourceIBMSchematicsActionUpdate(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
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
		updateActionOptions.SetTags(expandStringList(d.Get("tags").([]interface{})))
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
	if d.HasChange("bastion") {
		bastionAttr := d.Get("bastion").([]interface{})
		if len(bastionAttr) > 0 {
			bastion := resourceIBMSchematicsActionMapToTargetResourceset(d.Get("bastion.0").(map[string]interface{}))
			updateActionOptions.SetBastion(&bastion)
			hasChange = true
		}
	}
	if d.HasChange("targets_ini") {
		updateActionOptions.SetTargetsIni(d.Get("targets_ini").(string))
		hasChange = true
	}
	if d.HasChange("credentials") {
		var credentials []schematicsv1.VariableData
		for _, e := range d.Get("credentials").([]interface{}) {
			value := e.(map[string]interface{})
			credentialsItem := resourceIBMSchematicsActionMapToVariableData(value)
			credentials = append(credentials, credentialsItem)
		}
		updateActionOptions.SetCredentials(credentials)
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
	if d.HasChange("trigger_record_id") {
		updateActionOptions.SetTriggerRecordID(d.Get("trigger_record_id").(string))
		hasChange = true
	}
	if d.HasChange("state") {
		stateAttr := d.Get("state").([]interface{})
		if len(stateAttr) > 0 {
			state := resourceIBMSchematicsActionMapToActionState(d.Get("state.0").(map[string]interface{}))
			updateActionOptions.SetState(&state)
			hasChange = true
		}
	}
	if d.HasChange("sys_lock") {
		sysLockAttr := d.Get("sys_lock").([]interface{})
		if len(sysLockAttr) > 0 {
			sysLock := resourceIBMSchematicsActionMapToSystemLock(d.Get("sys_lock.0").(map[string]interface{}))
			updateActionOptions.SetSysLock(&sysLock)
			hasChange = true
		}
	}

	if hasChange {
		_, response, err := schematicsClient.UpdateActionWithContext(context.TODO(), updateActionOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateActionWithContext failed %s\n%s", err, response)
			return err
		}
	}

	return resourceIBMSchematicsActionRead(d, meta)
}

func resourceIBMSchematicsActionDelete(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	deleteActionOptions := &schematicsv1.DeleteActionOptions{}

	deleteActionOptions.SetActionID(d.Id())

	response, err := schematicsClient.DeleteActionWithContext(context.TODO(), deleteActionOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteActionWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}
