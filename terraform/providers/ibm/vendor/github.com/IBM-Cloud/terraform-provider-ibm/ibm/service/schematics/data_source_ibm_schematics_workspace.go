// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package schematics

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/schematics-go-sdk/schematicsv1"
)

func DataSourceIBMSchematicsWorkspace() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMSchematicsWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the workspace.  To find the workspace ID, use the `GET /v1/workspaces` API.",
			},
			"applied_shareddata_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of applied shared dataset ID.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"catalog_ref": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the software template that you chose from the IBM Cloud catalog. This information is returned for IBM Cloud catalog offerings only.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dry_run": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Dry run.",
						},
						"owning_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Owning account ID of the catalog.",
						},
						"item_icon_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to the icon of the software template in the IBM Cloud catalog.",
						},
						"item_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the software template that you chose to install from the IBM Cloud catalog. This software is provisioned with Schematics.",
						},
						"item_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the software that you chose to install from the IBM Cloud catalog.",
						},
						"item_readme_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to the readme file of the software template in the IBM Cloud catalog.",
						},
						"item_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to the software template in the IBM Cloud catalog.",
						},
						"launch_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to the dashboard to access your software.",
						},
						"offering_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the software template that you chose to install from the IBM Cloud catalog.",
						},
						"service_extensions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Service extensions defined as string of json",
						}}}},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the workspace was created.",
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID that created the workspace.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The workspace CRN.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the workspace.",
			},
			"last_health_check_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the last health check was performed by Schematics.",
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The IBM Cloud location where your workspace was provisioned.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the workspace.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group the workspace was provisioned in.",
			},
			"runtime_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the provisioning engine, state file, and runtime logs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_cmd": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The command that was used to apply the Terraform template or IBM Cloud catalog software template.",
						},
						"engine_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provisioning engine that was used to apply the Terraform template or IBM Cloud catalog software template.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the provisioning engine that was used.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID that was assigned to your Terraform template or IBM Cloud catalog software template.",
						},
						"log_store_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL to access the logs that were created during the creation, update, or deletion of your IBM Cloud resources.",
						},
						"output_values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of Output values.",
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
						},
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of resources.",
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
						},
						"state_store_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL where the Terraform statefile (`terraform.tfstate`) is stored. You can use the statefile to find an overview of IBM Cloud resources that were created by Schematics. Schematics uses the statefile as an inventory list to determine future create, update, or deletion jobs.",
						},
					},
				},
			},
			"shared_data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the Target used by the templates originating from IBM Cloud catalog offerings. This information is not relevant when you create a workspace from your own Terraform template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the cluster where you want to provision the resources of all IBM Cloud catalog templates that are included in the catalog offering.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Target cluster name.",
						},
						"entitlement_keys": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The entitlement key that you want to use to install IBM Cloud entitled software.",
							Elem: &schema.Schema{
								Type: schema.TypeMap,
							},
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Kubernetes namespace or OpenShift project where the resources of all IBM Cloud catalog templates that are included in the catalog offering are deployed into.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IBM Cloud region that you want to use for the resources of all IBM Cloud catalog templates that are included in the catalog offering.",
						},
						"resource_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the resource group that you want to use for the resources of all IBM Cloud catalog templates that are included in the catalog offering.",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the workspace.   **Active**: After you successfully ran your infrastructure code by applying your Terraform execution plan, the state of your workspace changes to `Active`.   **Connecting**: Schematics tries to connect to the template in your source repo. If successfully connected, the template is downloaded and metadata, such as input parameters, is extracted. After the template is downloaded, the state of the workspace changes to `Scanning`.   **Draft**: The workspace is created without a reference to a GitHub or GitLab repository.   **Failed**: If errors occur during the execution of your infrastructure code in IBM Cloud Schematics, your workspace status is set to `Failed`.   **Inactive**: The Terraform template was scanned successfully and the workspace creation is complete. You can now start running Schematics plan and apply jobs to provision the IBM Cloud resources that you specified in your template. If you have an `Active` workspace and decide to remove all your resources, your workspace is set to `Inactive` after all your resources are removed.   **In progress**: When you instruct IBM Cloud Schematics to run your infrastructure code by applying your Terraform execution plan, the status of our workspace changes to `In progress`.   **Scanning**: The download of the Terraform template is complete and vulnerability scanning started. If the scan is successful, the workspace state changes to `Inactive`. If errors in your template are found, the state changes to `Template Error`.   **Stopped**: The Schematics plan, apply, or destroy job was cancelled manually.   **Template Error**: The Schematics template contains errors and cannot be processed.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tags that are associated with the workspace.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template_env_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of environment values.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the variable.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enter the value as a string for the primitive types such as `bool`, `number`, `string`, and `HCL` format for the complex variables, as you provide in a `.tfvars` file. **You need to enter escaped string of `HCL` format for the complex variable value**. For more information, about how to declare variables in a terraform configuration file and provide value to schematics, see [Providing values for the declared variables](/docs/schematics?topic=schematics-create-tf-config#declare-variable).",
						},
						"secure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to `true`, the value of your input variable is protected and not returned in your API response.",
						},
						"hidden": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to `true`, the value of your input variable is protected and not returned in your API response.",
						},
					},
				},
			},
			"template_git_folder": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subfolder in your GitHub or GitLab repository where your Terraform template is stored. If your template is stored in the root directory, `.` is returned.",
			},
			"template_init_state_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Init state file.",
			},
			"template_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The Terraform version that was used to run your Terraform code.",
			},
			"template_uninstall_script_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Uninstall script name.",
			},
			"template_values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A list of variable values that you want to apply during the Helm chart installation. The list must be provided in JSON format, such as `\"\"autoscaling:  enabled: true  minReplicas: 2\"`. The values that you define here override the default Helm chart values. This field is supported only for IBM Cloud catalog offerings that are provisioned by using the Terraform Helm provider.",
			},
			"template_values_metadata": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of input variables that are associated with the workspace.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the variable.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the variable.",
						},
						"aliases": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of aliases for the variable name.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the meta data.",
						},
						"cloud_data_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cloud data type of the variable. eg. resource_group_id, region, vpc_id.",
						},
						"default": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Default value for the variable only if the override value is not specified.",
						},
						"link_status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the link.",
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
							Description: "If **true**, the variable is not displayed on UI or Command line.",
						},
						"required": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If the variable required?.",
						},
						"options": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of possible values for this variable.  If type is **integer** or **date**, then the array of string is  converted to array of integers or date during the runtime.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"min_value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum value of the variable. Applicable for the integer type.",
						},
						"max_value": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum value of the variable. Applicable for the integer type.",
						},
						"min_length": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The minimum length of the variable value. Applicable for the string type.",
						},
						"max_length": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum length of the variable value. Applicable for the string type.",
						},
						"matches": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The regex for the variable value.",
						},
						"position": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The relative position of this variable in a list.",
						},
						"group_by": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The display name of the group this variable belongs to.",
						},
						"source": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source of this meta-data.",
						},
						"metadata": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A list of input variables that are associated with the workspace.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Default value for the variable only if the override value is not specified.",
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the meta data.",
									},
									"hidden": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If **true**, the variable is not displayed on UI or Command line.",
									},
									"required": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If the variable required?.",
									},
									"options": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The list of possible values for this variable.  If type is **integer** or **date**, then the array of string is  converted to array of integers or date during the runtime.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Type of the variable.",
									},
									"secure": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "If set to `true`, the value of your input variable is protected and not returned in your API response.",
									},
								},
							},
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the variable. Applicable for the integer type.",
						},
					},
				},
			},
			"template_inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the input variables that your template uses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of your input variable.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the variable.",
						},
						"secure": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "If set to `true`, the value of your input variable is protected and not returned in your API response.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "`Terraform v0.11` supports `string`, `list`, `map` data type. For more information, about the syntax, see [Configuring input variables](https://www.terraform.io/docs/configuration-0-11/variables.html).<br> `Terraform v0.12` additionally, supports `bool`, `number` and complex data types such as `list(type)`, `map(type)`,`object({attribute name=type,..})`, `set(type)`, `tuple([type])`. For more information, about the syntax to use the complex data type, see [Configuring variables](https://www.terraform.io/docs/configuration/variables.html#type-constraints).",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enter the value as a string for the primitive types such as `bool`, `number`, `string`, and `HCL` format for the complex variables, as you provide in a `.tfvars` file. **You need to enter escaped string of `HCL` format for the complex variable value**. For more information, about how to declare variables in a terraform configuration file and provide value to schematics, see [Providing values for the declared variables](https://cloud.ibm.com/docs/schematics?topic=schematics-create-tf-config#declare-variable).",
						},
					},
				},
			},
			"template_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Workspace template ref.",
			},
			"template_git_branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The repository branch.",
			},
			"template_git_full_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Full repository URL.",
			},
			"template_git_has_uploadedgitrepotar": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Has uploaded Git repository tar.",
			},
			"template_git_release": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The repository release.",
			},
			"template_git_repo_sha_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The repository SHA value.",
			},
			"template_git_repo_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The repository URL.",
			},
			"template_git_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The source URL.",
			},

			/*"template_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Workspace type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},*/
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the workspace was last updated.",
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID that updated the workspace.",
			},
			"is_frozen": {
				Type:       schema.TypeBool,
				Computed:   true,
				Deprecated: "use frozen instead",
			},
			"frozen": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, the workspace is frozen and changes to the workspace are disabled.",
			},
			"frozen_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the workspace was frozen.",
			},
			"frozen_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID that froze the workspace.",
			},
			"is_locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, the workspace is locked and disabled for changes.",
				Deprecated:  "Use locked instead",
			},
			"locked": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, the workspace is locked and disabled for changes.",
			},
			"locked_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID that initiated a resource-related job, such as applying or destroying resources, that locked the workspace.",
			},
			"locked_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the workspace was locked.",
			},
			"status_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The success or error code that was returned for the last plan, apply, or destroy job that ran against your workspace.",
			},
			"status_msg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The success or error message that was returned for the last plan, apply, or destroy job that ran against your workspace.",
			},
			flex.ResourceControllerURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the IBM Cloud dashboard that can be used to explore and view details about this workspace",
			},
		},
	}
}

func SchematicsEndpointURL(region, meta interface{}) (string, bool, error) {
	sess, err := meta.(conns.ClientSession).BluemixSession()
	if err != nil {
		return "", false, err
	}
	if region != sess.Config.Region {
		// overide provider region with resource/datasource region argument
		// update client with updated region endpoint
		schematicsEndpoint := fmt.Sprintf("https://%s.%s", fmt.Sprintf("%s.schematics", region), "cloud.ibm.com")
		visibility := sess.Config.Visibility
		if visibility == "private" || visibility == "public-and-private" {
			schematicsEndpoint = fmt.Sprintf("https://%s.%s", fmt.Sprintf("private-%s.schematics", region), "cloud.ibm.com")
		}
		schematicsEndpointURL := conns.EnvFallBack([]string{"IBMCLOUD_SCHEMATICS_API_ENDPOINT"}, schematicsEndpoint)
		return schematicsEndpointURL, true, nil
	}
	return "", false, nil
}

func dataSourceIBMSchematicsWorkspaceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	schematicsClient, err := meta.(conns.ClientSession).SchematicsV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead schematicsClient initialization failed: %s", err.Error()), "ibm_schematics_workspace", "read")
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

	getWorkspaceOptions := &schematicsv1.GetWorkspaceOptions{}

	getWorkspaceOptions.SetWID(d.Get("workspace_id").(string))

	workspaceResponse, response, err := schematicsClient.GetWorkspaceWithContext(context, getWorkspaceOptions)
	if err != nil {

		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead GetWorkspaceWithContext failed with error: %s and response:\n%s", err, response), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getWorkspaceOptions.WID)

	if workspaceResponse.CatalogRef != nil {
		err = d.Set("catalog_ref", dataSourceWorkspaceResponseFlattenCatalogRef(*workspaceResponse.CatalogRef))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(workspaceResponse.CreatedAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("created_by", workspaceResponse.CreatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("crn", workspaceResponse.Crn); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("description", workspaceResponse.Description); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("last_health_check_at", flex.DateTimeToString(workspaceResponse.LastHealthCheckAt)); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("location", workspaceResponse.Location); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("name", workspaceResponse.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set("resource_group", workspaceResponse.ResourceGroup); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if workspaceResponse.RuntimeData != nil {
		err = d.Set("runtime_data", dataSourceWorkspaceResponseFlattenRuntimeData(workspaceResponse.RuntimeData))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if workspaceResponse.SharedData != nil {
		err = d.Set("shared_data", dataSourceWorkspaceResponseFlattenSharedData(*workspaceResponse.SharedData))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("status", workspaceResponse.Status); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if workspaceResponse.TemplateData != nil {
		templateData := dataSourceWorkspaceResponseFlattenTemplateData(workspaceResponse.TemplateData)

		if err = d.Set("template_env_settings", templateData[0]["env_values"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_git_folder", templateData[0]["folder"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_init_state_file", templateData[0]["init_state_file"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_type", templateData[0]["type"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_uninstall_script_name", templateData[0]["uninstall_script_name"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_values", templateData[0]["values"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		err = d.Set("template_values_metadata", dataSourceWorkspaceResponseFlattenValuesMetadata(templateData[0]["values_metadata"]))
		if err != nil {
			fmt.Println(fmt.Errorf("[ERROR] Error reading template_values_metadata %s", err))
		}
		if err = d.Set("template_inputs", templateData[0]["variablestore"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("template_ref", workspaceResponse.TemplateRef); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if workspaceResponse.TemplateRepo != nil {
		templateRepoMap := dataSourceWorkspaceResponseFlattenTemplateRepo(*workspaceResponse.TemplateRepo)
		if err = d.Set("template_git_branch", templateRepoMap[0]["branch"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_git_release", templateRepoMap[0]["release"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_git_repo_sha_value", templateRepoMap[0]["repo_sha_value"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_git_repo_url", templateRepoMap[0]["repo_url"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_git_url", templateRepoMap[0]["url"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("template_git_has_uploadedgitrepotar", templateRepoMap[0]["has_uploadedgitrepotar"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	/*if err = d.Set("type", workspaceResponse.Type); err != nil {
		return fmt.Errorf("[ERROR] Error setting type: %s", err)
	}*/
	if workspaceResponse.UpdatedAt != nil {
		if err = d.Set("updated_at", workspaceResponse.UpdatedAt.String()); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if err = d.Set("updated_by", workspaceResponse.UpdatedBy); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if workspaceResponse.WorkspaceStatus != nil {
		workspaceStatusMap := dataSourceWorkspaceResponseFlattenWorkspaceStatus(*workspaceResponse.WorkspaceStatus)
		if err = d.Set("is_frozen", workspaceStatusMap[0]["frozen"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("frozen", workspaceStatusMap[0]["frozen"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("frozen_at", workspaceStatusMap[0]["frozen_at"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("frozen_by", workspaceStatusMap[0]["frozen_by"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("is_locked", workspaceStatusMap[0]["locked"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("locked", workspaceStatusMap[0]["locked"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("locked_by", workspaceStatusMap[0]["locked_by"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("locked_time", workspaceStatusMap[0]["locked_time"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	if workspaceResponse.WorkspaceStatusMsg != nil {
		workspaceStatusMsgMap := dataSourceWorkspaceResponseFlattenWorkspaceStatusMsg(*workspaceResponse.WorkspaceStatusMsg)
		if err = d.Set("status_code", workspaceStatusMsgMap[0]["status_code"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set("status_msg", workspaceStatusMsgMap[0]["status_msg"]); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed with error: %s", err), "ibm_schematics_workspace", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	controller, err := flex.GetBaseController(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("dataSourceIBMSchematicsWorkspaceRead failed: %s", err.Error()), "ibm_schematics_workspace", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.Set(flex.ResourceControllerURL, controller+"/schematics")

	return nil
}

func dataSourceWorkspaceResponseFlattenCatalogRef(result schematicsv1.CatalogRef) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceWorkspaceResponseCatalogRefToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceWorkspaceResponseCatalogRefToMap(catalogRefItem schematicsv1.CatalogRef) (catalogRefMap map[string]interface{}) {
	catalogRefMap = map[string]interface{}{}

	if catalogRefItem.DryRun != nil {
		catalogRefMap["dry_run"] = catalogRefItem.DryRun
	}
	if catalogRefItem.OwningAccount != nil {
		catalogRefMap["owning_account"] = catalogRefItem.OwningAccount
	}
	if catalogRefItem.ItemIconURL != nil {
		catalogRefMap["item_icon_url"] = catalogRefItem.ItemIconURL
	}
	if catalogRefItem.ItemID != nil {
		catalogRefMap["item_id"] = catalogRefItem.ItemID
	}
	if catalogRefItem.ItemName != nil {
		catalogRefMap["item_name"] = catalogRefItem.ItemName
	}
	if catalogRefItem.ItemReadmeURL != nil {
		catalogRefMap["item_readme_url"] = catalogRefItem.ItemReadmeURL
	}
	if catalogRefItem.ItemURL != nil {
		catalogRefMap["item_url"] = catalogRefItem.ItemURL
	}
	if catalogRefItem.LaunchURL != nil {
		catalogRefMap["launch_url"] = catalogRefItem.LaunchURL
	}
	if catalogRefItem.OfferingVersion != nil {
		catalogRefMap["offering_version"] = catalogRefItem.OfferingVersion
	}
	if catalogRefItem.ServiceExtensions != nil {
		serviceExtensionsByte, err := json.MarshalIndent(catalogRefItem.ServiceExtensions, "", "")
		if err != nil {

		}
		serviceExtensionsJSON := string(serviceExtensionsByte[:])
		catalogRefMap["service_extensions"] = serviceExtensionsJSON
	}
	return catalogRefMap
}

func dataSourceWorkspaceResponseFlattenRuntimeData(result []schematicsv1.TemplateRunTimeDataResponse) (runtimeData []map[string]interface{}) {
	for _, runtimeDataItem := range result {
		runtimeData = append(runtimeData, dataSourceWorkspaceResponseRuntimeDataToMap(runtimeDataItem))
	}

	return runtimeData
}

func dataSourceWorkspaceResponseRuntimeDataToMap(runtimeDataItem schematicsv1.TemplateRunTimeDataResponse) (runtimeDataMap map[string]interface{}) {
	runtimeDataMap = map[string]interface{}{}

	if runtimeDataItem.EngineCmd != nil {
		runtimeDataMap["engine_cmd"] = runtimeDataItem.EngineCmd
	}
	if runtimeDataItem.EngineName != nil {
		runtimeDataMap["engine_name"] = runtimeDataItem.EngineName
	}
	if runtimeDataItem.EngineVersion != nil {
		runtimeDataMap["engine_version"] = runtimeDataItem.EngineVersion
	}
	if runtimeDataItem.ID != nil {
		runtimeDataMap["id"] = runtimeDataItem.ID
	}
	if runtimeDataItem.LogStoreURL != nil {
		runtimeDataMap["log_store_url"] = runtimeDataItem.LogStoreURL
	}
	if runtimeDataItem.OutputValues != nil {
		runtimeDataMap["output_values"] = runtimeDataItem.OutputValues
	}
	if runtimeDataItem.Resources != nil {
		runtimeDataMap["resources"] = runtimeDataItem.Resources
	}
	if runtimeDataItem.StateStoreURL != nil {
		runtimeDataMap["state_store_url"] = runtimeDataItem.StateStoreURL
	}

	return runtimeDataMap
}

func dataSourceWorkspaceResponseFlattenSharedData(result schematicsv1.SharedTargetDataResponse) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceWorkspaceResponseSharedDataToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceWorkspaceResponseSharedDataToMap(sharedDataItem schematicsv1.SharedTargetDataResponse) (sharedDataMap map[string]interface{}) {
	sharedDataMap = map[string]interface{}{}

	if sharedDataItem.ClusterID != nil {
		sharedDataMap["cluster_id"] = sharedDataItem.ClusterID
	}
	if sharedDataItem.ClusterName != nil {
		sharedDataMap["cluster_name"] = sharedDataItem.ClusterName
	}
	if sharedDataItem.EntitlementKeys != nil {
		sharedDataMap["entitlement_keys"] = sharedDataItem.EntitlementKeys
	}
	if sharedDataItem.Namespace != nil {
		sharedDataMap["namespace"] = sharedDataItem.Namespace
	}
	if sharedDataItem.Region != nil {
		sharedDataMap["region"] = sharedDataItem.Region
	}
	if sharedDataItem.ResourceGroupID != nil {
		sharedDataMap["resource_group_id"] = sharedDataItem.ResourceGroupID
	}

	return sharedDataMap
}

func dataSourceWorkspaceResponseFlattenTemplateData(result []schematicsv1.TemplateSourceDataResponse) (templateData []map[string]interface{}) {
	for _, templateDataItem := range result {
		templateData = append(templateData, dataSourceWorkspaceResponseTemplateDataToMap(templateDataItem))
	}

	return templateData
}

func dataSourceWorkspaceResponseTemplateDataToMap(templateDataItem schematicsv1.TemplateSourceDataResponse) (templateDataMap map[string]interface{}) {
	templateDataMap = map[string]interface{}{}

	if templateDataItem.EnvValues != nil {
		envValuesList := []map[string]interface{}{}
		for _, envValuesItem := range templateDataItem.EnvValues {
			envValuesList = append(envValuesList, dataSourceWorkspaceResponseTemplateDataEnvValuesToMap(envValuesItem))
		}
		templateDataMap["env_values"] = envValuesList
	}
	if templateDataItem.Folder != nil {
		templateDataMap["folder"] = templateDataItem.Folder
	}
	if templateDataItem.Compact != nil {
		templateDataMap["compact"] = templateDataItem.Compact
	}
	if templateDataItem.HasGithubtoken != nil {
		templateDataMap["has_githubtoken"] = templateDataItem.HasGithubtoken
	}
	if templateDataItem.ID != nil {
		templateDataMap["id"] = templateDataItem.ID
	}
	if templateDataItem.Type != nil {
		templateDataMap["type"] = templateDataItem.Type
	}
	if templateDataItem.UninstallScriptName != nil {
		templateDataMap["uninstall_script_name"] = templateDataItem.UninstallScriptName
	}
	if templateDataItem.Values != nil {
		templateDataMap["values"] = templateDataItem.Values
	}
	if templateDataItem.ValuesMetadata != nil {
		valuesMetadataList := []interface{}{}
		for _, valuesMetadataItem := range templateDataItem.ValuesMetadata {
			valuesMetadataList = append(valuesMetadataList, valuesMetadataItem)
		}
		templateDataMap["values_metadata"] = valuesMetadataList
	}
	if templateDataItem.ValuesURL != nil {
		templateDataMap["values_url"] = templateDataItem.ValuesURL
	}
	if templateDataItem.Variablestore != nil {
		variablestoreList := []map[string]interface{}{}
		for _, variablestoreItem := range templateDataItem.Variablestore {
			variablestoreList = append(variablestoreList, dataSourceWorkspaceResponseTemplateDataVariablestoreToMap(variablestoreItem))
		}
		templateDataMap["variablestore"] = variablestoreList
	}

	return templateDataMap
}

func dataSourceWorkspaceResponseFlattenValuesMetadata(result interface{}) (valuesMetadata []map[string]interface{}) {
	if result != nil {
		for _, res := range result.([]interface{}) {
			valuesMetadataMap := dataSourceWorkspaceResponseValuesMetadataToMap(res.(map[string]interface{}))
			valuesMetadata = append(valuesMetadata, valuesMetadataMap)
		}
	}
	return valuesMetadata
}

func dataSourceWorkspaceResponseValuesMetadataToMap(valuesMetadataItem map[string]interface{}) map[string]interface{} {
	valuesMetadataMap := map[string]interface{}{}

	if valuesMetadataItem["name"] != nil {
		valuesMetadataMap["name"] = valuesMetadataItem["name"].(string)
	}
	if valuesMetadataItem["type"] != nil {
		valuesMetadataMap["type"] = valuesMetadataItem["type"].(string)
	}

	if valuesMetadataItem["aliases"] != nil {
		valuesMetadataMap["aliases"] = valuesMetadataItem["aliases"]
	}

	if valuesMetadataItem["description"] != nil {
		valuesMetadataMap["description"] = valuesMetadataItem["description"].(string)
	}

	if valuesMetadataItem["cloud_data_type"] != nil {
		valuesMetadataMap["cloud_data_type"] = valuesMetadataItem["cloud_data_type"].(string)
	}

	if valuesMetadataItem["default"] != nil {
		valuesMetadataMap["default"] = valuesMetadataItem["default"].(string)
	}

	if valuesMetadataItem["link_status"] != nil {
		valuesMetadataMap["link_status"] = valuesMetadataItem["link_status"].(string)
	}

	if valuesMetadataItem["secure"] != nil {
		valuesMetadataMap["secure"] = valuesMetadataItem["secure"]
	}

	if valuesMetadataItem["immutable"] != nil {
		valuesMetadataMap["immutable"] = valuesMetadataItem["immutable"]
	}

	if valuesMetadataItem["hidden"] != nil {
		valuesMetadataMap["hidden"] = valuesMetadataItem["hidden"]
	}

	if valuesMetadataItem["required"] != nil {
		valuesMetadataMap["required"] = valuesMetadataItem["required"]
	}

	if valuesMetadataItem["options"] != nil {
		valuesMetadataMap["options"] = valuesMetadataItem["options"]
	}

	if valuesMetadataItem["min_value"] != nil {
		valuesMetadataMap["min_value"] = valuesMetadataItem["min_value"]
	}

	if valuesMetadataItem["max_value"] != nil {
		valuesMetadataMap["max_value"] = valuesMetadataItem["max_value"]
	}

	if valuesMetadataItem["min_length"] != nil {
		valuesMetadataMap["min_length"] = valuesMetadataItem["min_length"]
	}

	if valuesMetadataItem["max_length"] != nil {
		valuesMetadataMap["max_length"] = valuesMetadataItem["max_length"]
	}

	if valuesMetadataItem["matches"] != nil {
		valuesMetadataMap["matches"] = valuesMetadataItem["matches"].(string)
	}

	if valuesMetadataItem["position"] != nil {
		valuesMetadataMap["position"] = valuesMetadataItem["position"]
	}

	if valuesMetadataItem["group_by"] != nil {
		valuesMetadataMap["group_by"] = valuesMetadataItem["group_by"].(string)
	}

	if valuesMetadataItem["source"] != nil {
		valuesMetadataMap["source"] = valuesMetadataItem["source"].(string)
	}

	if valuesMetadataItem["metadata"] != nil {
		metadataList := []map[string]interface{}{}

		valuesMetadataMap["metadata"] = append(metadataList, valuesMetadataItem["metadata"].(map[string]interface{}))

	}
	if valuesMetadataItem["value"] != nil {
		valuesMetadataMap["value"] = valuesMetadataItem["value"].(string)
	}

	return valuesMetadataMap
}
func dataSourceIbmSchematicsWorkspaceVariableMetadataToMap(model *schematicsv1.VariableMetadata) map[string]interface{} {
	modelMap := make(map[string]interface{})
	if model.Type != nil {
		modelMap["type"] = *model.Type
	}
	if model.Aliases != nil {
		modelMap["aliases"] = model.Aliases
	}
	if model.Description != nil {
		modelMap["description"] = *model.Description
	}
	if model.CloudDataType != nil {
		modelMap["cloud_data_type"] = *model.CloudDataType
	}
	if model.DefaultValue != nil {
		modelMap["default_value"] = *model.DefaultValue
	}
	if model.LinkStatus != nil {
		modelMap["link_status"] = *model.LinkStatus
	}
	if model.Secure != nil {
		modelMap["secure"] = *model.Secure
	}
	if model.Immutable != nil {
		modelMap["immutable"] = *model.Immutable
	}
	if model.Hidden != nil {
		modelMap["hidden"] = *model.Hidden
	}
	if model.Required != nil {
		modelMap["required"] = *model.Required
	}
	if model.Options != nil {
		modelMap["options"] = model.Options
	}
	if model.MinValue != nil {
		modelMap["min_value"] = *model.MinValue
	}
	if model.MaxValue != nil {
		modelMap["max_value"] = *model.MaxValue
	}
	if model.MinLength != nil {
		modelMap["min_length"] = *model.MinLength
	}
	if model.MaxLength != nil {
		modelMap["max_length"] = *model.MaxLength
	}
	if model.Matches != nil {
		modelMap["matches"] = *model.Matches
	}
	if model.Position != nil {
		modelMap["position"] = *model.Position
	}
	if model.GroupBy != nil {
		modelMap["group_by"] = *model.GroupBy
	}
	if model.Source != nil {
		modelMap["source"] = *model.Source
	}
	return modelMap
}
func dataSourceWorkspaceResponseTemplateDataEnvValuesToMap(envValuesItem schematicsv1.EnvVariableResponse) (envValuesMap map[string]interface{}) {
	envValuesMap = map[string]interface{}{}

	if envValuesItem.Hidden != nil {
		envValuesMap["hidden"] = *envValuesItem.Hidden
	}
	if envValuesItem.Name != nil {
		envValuesMap["name"] = envValuesItem.Name
	}
	if envValuesItem.Secure != nil {
		envValuesMap["secure"] = *envValuesItem.Secure
	}
	if envValuesItem.Value != nil {
		envValuesMap["value"] = envValuesItem.Value
	}

	return envValuesMap
}

func dataSourceWorkspaceResponseTemplateDataVariablestoreToMap(variablestoreItem schematicsv1.WorkspaceVariableResponse) (variablestoreMap map[string]interface{}) {
	variablestoreMap = map[string]interface{}{}

	if variablestoreItem.Description != nil {
		variablestoreMap["description"] = variablestoreItem.Description
	}
	if variablestoreItem.Name != nil {
		variablestoreMap["name"] = variablestoreItem.Name
	}
	if variablestoreItem.Secure != nil {
		variablestoreMap["secure"] = variablestoreItem.Secure
	}
	if variablestoreItem.Type != nil {
		variablestoreMap["type"] = variablestoreItem.Type
	}
	if variablestoreItem.Value != nil {
		variablestoreMap["value"] = variablestoreItem.Value
	}

	return variablestoreMap
}

func dataSourceWorkspaceResponseFlattenTemplateRepo(result schematicsv1.TemplateRepoResponse) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceWorkspaceResponseTemplateRepoToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceWorkspaceResponseTemplateRepoToMap(templateRepoItem schematicsv1.TemplateRepoResponse) (templateRepoMap map[string]interface{}) {
	templateRepoMap = map[string]interface{}{}

	if templateRepoItem.Branch != nil {
		templateRepoMap["branch"] = templateRepoItem.Branch
	}
	if templateRepoItem.FullURL != nil {
		templateRepoMap["full_url"] = templateRepoItem.FullURL
	}
	if templateRepoItem.HasUploadedgitrepotar != nil {
		templateRepoMap["has_uploadedgitrepotar"] = templateRepoItem.HasUploadedgitrepotar
	}
	if templateRepoItem.Release != nil {
		templateRepoMap["release"] = templateRepoItem.Release
	}
	if templateRepoItem.RepoShaValue != nil {
		templateRepoMap["repo_sha_value"] = templateRepoItem.RepoShaValue
	}
	if templateRepoItem.RepoURL != nil {
		templateRepoMap["repo_url"] = templateRepoItem.RepoURL
	}
	if templateRepoItem.URL != nil {
		templateRepoMap["url"] = templateRepoItem.URL
	}

	return templateRepoMap
}

func dataSourceWorkspaceResponseFlattenWorkspaceStatus(result schematicsv1.WorkspaceStatusResponse) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceWorkspaceResponseWorkspaceStatusToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceWorkspaceResponseWorkspaceStatusToMap(workspaceStatusItem schematicsv1.WorkspaceStatusResponse) (workspaceStatusMap map[string]interface{}) {
	workspaceStatusMap = map[string]interface{}{}

	if workspaceStatusItem.Frozen != nil {
		workspaceStatusMap["frozen"] = workspaceStatusItem.Frozen
	}
	if workspaceStatusItem.FrozenAt != nil {
		workspaceStatusMap["frozen_at"] = workspaceStatusItem.FrozenAt.String()
	}
	if workspaceStatusItem.FrozenBy != nil {
		workspaceStatusMap["frozen_by"] = workspaceStatusItem.FrozenBy
	}
	if workspaceStatusItem.Locked != nil {
		workspaceStatusMap["locked"] = workspaceStatusItem.Locked
	}
	if workspaceStatusItem.LockedBy != nil {
		workspaceStatusMap["locked_by"] = workspaceStatusItem.LockedBy
	}
	if workspaceStatusItem.LockedTime != nil {
		workspaceStatusMap["locked_time"] = workspaceStatusItem.LockedTime.String()
	}

	return workspaceStatusMap
}

func dataSourceWorkspaceResponseFlattenWorkspaceStatusMsg(result schematicsv1.WorkspaceStatusMessage) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceWorkspaceResponseWorkspaceStatusMsgToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceWorkspaceResponseWorkspaceStatusMsgToMap(workspaceStatusMsgItem schematicsv1.WorkspaceStatusMessage) (workspaceStatusMsgMap map[string]interface{}) {
	workspaceStatusMsgMap = map[string]interface{}{}

	if workspaceStatusMsgItem.StatusCode != nil {
		workspaceStatusMsgMap["status_code"] = workspaceStatusMsgItem.StatusCode
	}
	if workspaceStatusMsgItem.StatusMsg != nil {
		workspaceStatusMsgMap["status_msg"] = workspaceStatusMsgItem.StatusMsg
	}

	return workspaceStatusMsgMap
}
