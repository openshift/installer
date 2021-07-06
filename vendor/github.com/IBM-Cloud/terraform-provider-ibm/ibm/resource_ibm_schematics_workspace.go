// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/schematics-go-sdk/schematicsv1"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	schematicsWorkspaceName         = "name"
	schematicsWorkspaceDescription  = "description"
	schematicsWorkspaceTemplateType = "template_type"
)

func resourceIBMSchematicsWorkspace() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMSchematicsWorkspaceCreate,
		Read:     resourceIBMSchematicsWorkspaceRead,
		Update:   resourceIBMSchematicsWorkspaceUpdate,
		Delete:   resourceIBMSchematicsWorkspaceDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"applied_shareddata_ids": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of applied shared dataset id.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"catalog_ref": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Information about the software template that you chose from the IBM Cloud catalog. This information is returned for IBM Cloud catalog offerings only.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dry_run": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Dry run.",
						},
						"item_icon_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL to the icon of the software template in the IBM Cloud catalog.",
						},
						"item_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the software template that you chose to install from the IBM Cloud catalog. This software is provisioned with Schematics.",
						},
						"item_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the software that you chose to install from the IBM Cloud catalog.",
						},
						"item_readme_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL to the readme file of the software template in the IBM Cloud catalog.",
						},
						"item_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL to the software template in the IBM Cloud catalog.",
						},
						"launch_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL to the dashboard to access your software.",
						},
						"offering_version": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The version of the software template that you chose to install from the IBM Cloud catalog.",
						},
					},
				},
			},
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The description of the workspace.",
				ValidateFunc: InvokeValidator("ibm_schematics_workspace", schematicsWorkspaceDescription),
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location where you want to create your Schematics workspace and run Schematics actions. The location that you enter must match the API endpoint that you use. For example, if you use the Frankfurt API endpoint, you must specify `eu-de` as your location. If you use an API endpoint for a geography and you do not specify a location, Schematics determines the location based on availability.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The name of your workspace. The name can be up to 128 characters long and can include alphanumeric characters, spaces, dashes, and underscores. When you create a workspace for your own Terraform template, consider including the microservice component that you set up with your Terraform template and the IBM Cloud environment where you want to deploy your resources in your name.",
				ValidateFunc: InvokeValidator("ibm_schematics_workspace", schematicsWorkspaceName),
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the resource group where you want to provision the workspace.",
			},
			"shared_data": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "Information that is shared across templates in IBM Cloud catalog offerings. This information is not provided when you create a workspace from your own Terraform template.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_created_on": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster created on.",
						},
						"cluster_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the cluster where you want to provision the resources of all IBM Cloud catalog templates that are included in the catalog offering.",
						},
						"cluster_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster name.",
						},
						"cluster_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster type.",
						},
						"entitlement_keys": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The entitlement key that you want to use to install IBM Cloud entitled software.",
							Elem:        &schema.Schema{Type: schema.TypeMap},
						},
						"namespace": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The Kubernetes namespace or OpenShift project where the resources of all IBM Cloud catalog templates that are included in the catalog offering are deployed into.",
						},
						"region": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The IBM Cloud region that you want to use for the resources of all IBM Cloud catalog templates that are included in the catalog offering.",
						},
						"resource_group_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the resource group that you want to use for the resources of all IBM Cloud catalog templates that are included in the catalog offering.",
						},
						"worker_count": &schema.Schema{
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cluster worker count.",
						},
						"worker_machine_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster worker type.",
						},
					},
				},
			},
			"tags": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of tags that are associated with the workspace.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"template_env_settings": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of environment variables that you want to apply during the execution of a bash script or Terraform action. This field must be provided as a list of key-value pairs, for example, **TF_LOG=debug**. Each entry will be a map with one entry where `key is the environment variable name and value is value`. You can define environment variables for IBM Cloud catalog offerings that are provisioned by using a bash script.",
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
			"template_git_folder": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The subfolder in your GitHub or GitLab repository where your Terraform template is stored.",
			},
			"template_init_state_file": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The content of an existing Terraform statefile that you want to import in to your workspace. To get the content of a Terraform statefile for a specific Terraform template in an existing workspace, run `ibmcloud terraform state pull --id <workspace_id> --template <template_id>`.",
			},
			"template_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The Terraform version that you want to use to run your Terraform code. Enter `terraform_v0.12` to use Terraform version 0.12, and `terraform_v0.11` to use Terraform version 0.11. Make sure that your Terraform config files are compatible with the Terraform version that you select.",
				ValidateFunc: InvokeValidator("ibm_schematics_workspace", schematicsWorkspaceTemplateType),
			},
			"template_uninstall_script_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Uninstall script name.",
			},
			"template_values": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A list of variable values that you want to apply during the Helm chart installation. The list must be provided in JSON format, such as `\"autoscaling:  enabled: true  minReplicas: 2\"`. The values that you define here override the default Helm chart values. This field is supported only for IBM Cloud catalog offerings that are provisioned by using the Terraform Helm provider.",
			},
			"template_values_metadata": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of values metadata.",
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
			"template_inputs": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "VariablesRequest -.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of your input variable.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the variable.",
						},
						"secure": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "If set to `true`, the value of your input variable is protected and not returned in your API response.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "`Terraform v0.11` supports `string`, `list`, `map` data type. For more information, about the syntax, see [Configuring input variables](https://www.terraform.io/docs/configuration-0-11/variables.html). <br> `Terraform v0.12` additionally, supports `bool`, `number` and complex data types such as `list(type)`, `map(type)`, `object({attribute name=type,..})`, `set(type)`, `tuple([type])`. For more information, about the syntax to use the complex data type, see [Configuring variables](https://www.terraform.io/docs/configuration/variables.html#type-constraints).",
						},
						"use_default": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Variable uses default value; and is not over-ridden.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Enter the value as a string for the primitive types such as `bool`, `number`, `string`, and `HCL` format for the complex variables, as you provide in a `.tfvars` file. **You need to enter escaped string of `HCL` format for the complex variable value**. For more information, about how to declare variables in a terraform configuration file and provide value to schematics, see [Providing values for the declared variables](/docs/schematics?topic=schematics-create-tf-config#declare-variable).",
						},
					},
				},
			},
			"template_ref": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workspace template ref.",
			},
			"template_git_branch": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The branch in GitHub where your Terraform template is stored.",
			},
			"template_git_release": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The release tag in GitHub of your Terraform template.",
			},
			"template_git_repo_sha_value": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Repo SHA value.",
			},
			"template_git_repo_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URL to the repository where the IBM Cloud catalog software template is stored.",
			},
			"template_git_url": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The URL to the GitHub or GitLab repository where your Terraform and public bit bucket template is stored. For more information of the environment variable syntax, see [Create workspace new](/docs/schematics?topic=schematics-schematics-cli-reference#schematics-workspace-new).",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"template_git_has_uploadedgitrepotar": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Has uploaded git repo tar",
			},
			/*"template_type": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "List of Workspace type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},*/
			"frozen": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If set to true, the workspace is frozen and changes to the workspace are disabled.",
			},
			"frozen_at": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The timestamp when the workspace was frozen.",
			},
			"frozen_by": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The user ID that froze the workspace.",
			},
			"locked": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If set to true, the workspace is locked and disabled for changes.",
			},
			"locked_by": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user ID that initiated a resource-related action, such as applying or destroying resources, that locked the workspace.",
			},
			"locked_time": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The timestamp when the workspace was locked.",
			},
			"x_github_token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The personal access token to authenticate with your private GitHub or GitLab repository and access your Terraform template.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the workspace was created.",
			},
			"created_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID that created the workspace.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Workspace CRN.",
			},
			"last_health_check_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the last health check was performed by Schematics.",
			},
			"runtime_data": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information about the provisioning engine, state file, and runtime logs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_cmd": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The command that was used to apply the Terraform template or IBM Cloud catalog software template.",
						},
						"engine_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The provisioning engine that was used to apply the Terraform template or IBM Cloud catalog software template.",
						},
						"engine_version": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The version of the provisioning engine that was used.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID that was assigned to your Terraform template or IBM Cloud catalog software template.",
						},
						"log_store_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL to access the logs that were created during the creation, update, or deletion of your IBM Cloud resources.",
						},
						"output_values": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of Output values.",
							Elem:        &schema.Schema{Type: schema.TypeMap},
						},
						"resources": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of resources.",
							Elem:        &schema.Schema{Type: schema.TypeMap},
						},
						"state_store_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL where the Terraform statefile (`terraform.tfstate`) is stored. You can use the statefile to find an overview of IBM Cloud resources that were created by Schematics. Schematics uses the statefile as an inventory list to determine future create, update, or deletion actions.",
						},
					},
				},
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the workspace.  **Active**: After you successfully ran your infrastructure code by applying your Terraform execution plan, the state of your workspace changes to `Active`.  **Connecting**: Schematics tries to connect to the template in your source repo. If successfully connected, the template is downloaded and metadata, such as input parameters, is extracted. After the template is downloaded, the state of the workspace changes to `Scanning`.  **Draft**: The workspace is created without a reference to a GitHub or GitLab repository.  **Failed**: If errors occur during the execution of your infrastructure code in IBM Cloud Schematics, your workspace status is set to `Failed`.  **Inactive**: The Terraform template was scanned successfully and the workspace creation is complete. You can now start running Schematics plan and apply actions to provision the IBM Cloud resources that you specified in your template. If you have an `Active` workspace and decide to remove all your resources, your workspace is set to `Inactive` after all your resources are removed.  **In progress**: When you instruct IBM Cloud Schematics to run your infrastructure code by applying your Terraform execution plan, the status of our workspace changes to `In progress`.  **Scanning**: The download of the Terraform template is complete and vulnerability scanning started. If the scan is successful, the workspace state changes to `Inactive`. If errors in your template are found, the state changes to `Template Error`.  **Stopped**: The Schematics plan, apply, or destroy action was cancelled manually.  **Template Error**: The Schematics template contains errors and cannot be processed.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the workspace was last updated.",
			},
			"updated_by": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID that updated the workspace.",
			},
			"status_code": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The success or error code that was returned for the last plan, apply, or destroy action that ran against your workspace.",
			},
			"status_msg": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The success or error message that was returned for the last plan, apply, or destroy action that ran against your workspace.",
			},
		},
	}
}

func resourceIBMSchematicsWorkspaceValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 schematicsWorkspaceName,
			ValidateFunctionIdentifier: ValidateRegexp,
			Type:                       TypeString,
			Regexp:                     `^[a-zA-Z0-9][a-zA-Z0-9-_ ]*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
			Required:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 schematicsWorkspaceDescription,
			ValidateFunctionIdentifier: StringLenBetween,
			Type:                       TypeString,
			MinValueLength:             0,
			MaxValueLength:             2048,
			Optional:                   true})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 schematicsWorkspaceTemplateType,
			ValidateFunctionIdentifier: ValidateRegexp,
			Type:                       TypeString,
			Regexp:                     `^terraform_v0\.(?:11|12|13)(?:\.\d+)?$`,
			Default:                    "[]",
			Optional:                   true})

	ibmSchematicsWorkspaceResourceValidator := ResourceValidator{ResourceName: "ibm_schematics_workspace", Schema: validateSchema}
	return &ibmSchematicsWorkspaceResourceValidator
}

func resourceIBMSchematicsWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	createWorkspaceOptions := &schematicsv1.CreateWorkspaceOptions{}

	if _, ok := d.GetOk("applied_shareddata_ids"); ok {
		createWorkspaceOptions.SetAppliedShareddataIds(expandStringList(d.Get("applied_shareddata_ids").([]interface{})))
	}
	if _, ok := d.GetOk("catalog_ref"); ok {
		catalogRefAttr := d.Get("catalog_ref").([]interface{})
		if len(catalogRefAttr) > 0 {
			catalogRef := resourceIBMSchematicsWorkspaceMapToCatalogRef(d.Get("catalog_ref.0").(map[string]interface{}))
			createWorkspaceOptions.SetCatalogRef(&catalogRef)
		}
	}
	if _, ok := d.GetOk("description"); ok {
		createWorkspaceOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("location"); ok {
		createWorkspaceOptions.SetLocation(d.Get("location").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		createWorkspaceOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("resource_group"); ok {
		createWorkspaceOptions.SetResourceGroup(d.Get("resource_group").(string))
	}
	if _, ok := d.GetOk("shared_data"); ok {
		sharedDataAttr := d.Get("shared_data").([]interface{})
		if len(sharedDataAttr) > 0 {
			sharedData := resourceIBMSchematicsWorkspaceMapToSharedTargetData(d.Get("shared_data.0").(map[string]interface{}))
			createWorkspaceOptions.SetSharedData(&sharedData)
		}
	}
	if _, ok := d.GetOk("tags"); ok {
		createWorkspaceOptions.SetTags(expandStringList(d.Get("tags").([]interface{})))
	}

	var templateData []schematicsv1.TemplateSourceDataRequest

	templateSourceDataRequestMap := map[string]interface{}{}
	hasTemplateData := false

	if _, ok := d.GetOk("template_env_settings"); ok {
		templateSourceDataRequestMap["env_values"] = d.Get("template_env_settings").([]interface{})
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_git_folder"); ok {
		templateSourceDataRequestMap["folder"] = d.Get("template_git_folder").(string)
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_init_state_file"); ok {
		templateSourceDataRequestMap["init_state_file"] = d.Get("template_init_state_file").(string)
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_type"); ok {
		templateSourceDataRequestMap["type"] = d.Get("template_type").(string)
		createWorkspaceOptions.SetType([]string{d.Get("template_type").(string)})
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_uninstall_script_name"); ok {
		templateSourceDataRequestMap["uninstall_script_name"] = d.Get("template_uninstall_script_name").(string)
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_values"); ok {
		templateSourceDataRequestMap["values"] = d.Get("template_values").(string)
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_values_metadata"); ok {
		templateSourceDataRequestMap["values_metadata"] = d.Get("template_values_metadata").([]interface{})
		hasTemplateData = true
	}
	if _, ok := d.GetOk("template_inputs"); ok {
		templateSourceDataRequestMap["variablestore"] = d.Get("template_inputs").([]interface{})
		hasTemplateData = true
	}
	if hasTemplateData {
		templateDataItem := resourceIBMSchematicsWorkspaceMapToTemplateSourceDataRequest(templateSourceDataRequestMap)
		templateData = append(templateData, templateDataItem)
		createWorkspaceOptions.SetTemplateData(templateData)
	}
	if _, ok := d.GetOk("template_ref"); ok {
		createWorkspaceOptions.SetTemplateRef(d.Get("template_ref").(string))
	}

	templateRepoRequestMap := map[string]interface{}{}
	hasTemplateRepo := false
	if _, ok := d.GetOk("template_git_branch"); ok {
		templateRepoRequestMap["branch"] = d.Get("template_git_branch").(string)
		hasTemplateRepo = true
	}
	if _, ok := d.GetOk("template_git_release"); ok {
		templateRepoRequestMap["release"] = d.Get("template_git_release").(string)
		hasTemplateRepo = true
	}
	if _, ok := d.GetOk("template_git_repo_sha_value"); ok {
		templateRepoRequestMap["repo_sha_value"] = d.Get("template_git_repo_sha_value").(string)
		hasTemplateRepo = true
	}
	if _, ok := d.GetOk("template_git_repo_url"); ok {
		templateRepoRequestMap["repo_url"] = d.Get("template_git_repo_url").(string)
		hasTemplateRepo = true
	}
	if _, ok := d.GetOk("template_git_url"); ok {
		templateRepoRequestMap["url"] = d.Get("template_git_url").(string)
		hasTemplateRepo = true
	}
	if _, ok := d.GetOk("template_git_has_uploadedgitrepotar"); ok {
		templateRepoRequestMap["has_uploadedgitrepotar"] = d.Get("template_git_has_uploadedgitrepotar").(string)
		hasTemplateRepo = true
	}
	if hasTemplateRepo {
		templateRepo := resourceIBMSchematicsWorkspaceMapToTemplateRepoRequest(templateRepoRequestMap)
		createWorkspaceOptions.SetTemplateRepo(&templateRepo)
	}

	/*if _, ok := d.GetOk("template_type"); ok {
		createWorkspaceOptions.SetType(expandStringList(d.Get("template_type").([]interface{})))
	}*/
	workspaceStatusRequestMap := map[string]interface{}{}
	hasWorkspaceStatus := false
	if _, ok := d.GetOk("frozen"); ok {
		workspaceStatusRequestMap["frozen"] = d.Get("frozen").(bool)
		hasWorkspaceStatus = true
	}
	if _, ok := d.GetOk("frozen_at"); ok {
		workspaceStatusRequestMap["frozen_at"] = d.Get("frozen_at").(string)
		hasWorkspaceStatus = true
	}
	if _, ok := d.GetOk("frozen_by"); ok {
		workspaceStatusRequestMap["frozen_by"] = d.Get("frozen_by").(string)
		hasWorkspaceStatus = true
	}
	if _, ok := d.GetOk("locked"); ok {
		workspaceStatusRequestMap["locked"] = d.Get("locked").(bool)
		hasWorkspaceStatus = true
	}
	if _, ok := d.GetOk("locked_by"); ok {
		workspaceStatusRequestMap["locked_by"] = d.Get("locked_by").(string)
		hasWorkspaceStatus = true
	}
	if _, ok := d.GetOk("locked_time"); ok {
		workspaceStatusRequestMap["locked_time"] = d.Get("locked_time").(string)
		hasWorkspaceStatus = true
	}
	if hasWorkspaceStatus {
		workspaceStatus := resourceIBMSchematicsWorkspaceMapToWorkspaceStatusRequest(workspaceStatusRequestMap)
		createWorkspaceOptions.SetWorkspaceStatus(&workspaceStatus)
	}

	if _, ok := d.GetOk("x_github_token"); ok {
		createWorkspaceOptions.SetXGithubToken(d.Get("x_github_token").(string))
	}

	workspaceResponse, response, err := schematicsClient.CreateWorkspaceWithContext(context.TODO(), createWorkspaceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateWorkspaceWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId(*workspaceResponse.ID)

	return resourceIBMSchematicsWorkspaceRead(d, meta)
}

func resourceIBMSchematicsWorkspaceMapToCatalogRef(catalogRefMap map[string]interface{}) schematicsv1.CatalogRef {
	catalogRef := schematicsv1.CatalogRef{}

	if catalogRefMap["dry_run"] != nil {
		catalogRef.DryRun = core.BoolPtr(catalogRefMap["dry_run"].(bool))
	}
	if catalogRefMap["item_icon_url"] != nil {
		catalogRef.ItemIconURL = core.StringPtr(catalogRefMap["item_icon_url"].(string))
	}
	if catalogRefMap["item_id"] != nil {
		catalogRef.ItemID = core.StringPtr(catalogRefMap["item_id"].(string))
	}
	if catalogRefMap["item_name"] != nil {
		catalogRef.ItemName = core.StringPtr(catalogRefMap["item_name"].(string))
	}
	if catalogRefMap["item_readme_url"] != nil {
		catalogRef.ItemReadmeURL = core.StringPtr(catalogRefMap["item_readme_url"].(string))
	}
	if catalogRefMap["item_url"] != nil {
		catalogRef.ItemURL = core.StringPtr(catalogRefMap["item_url"].(string))
	}
	if catalogRefMap["launch_url"] != nil {
		catalogRef.LaunchURL = core.StringPtr(catalogRefMap["launch_url"].(string))
	}
	if catalogRefMap["offering_version"] != nil {
		catalogRef.OfferingVersion = core.StringPtr(catalogRefMap["offering_version"].(string))
	}

	return catalogRef
}

func resourceIBMSchematicsWorkspaceMapToSharedTargetData(sharedTargetDataMap map[string]interface{}) schematicsv1.SharedTargetData {
	sharedTargetData := schematicsv1.SharedTargetData{}

	if sharedTargetDataMap["cluster_created_on"] != nil {
		sharedTargetData.ClusterCreatedOn = core.StringPtr(sharedTargetDataMap["cluster_created_on"].(string))
	}
	if sharedTargetDataMap["cluster_id"] != nil {
		sharedTargetData.ClusterID = core.StringPtr(sharedTargetDataMap["cluster_id"].(string))
	}
	if sharedTargetDataMap["cluster_name"] != nil {
		sharedTargetData.ClusterName = core.StringPtr(sharedTargetDataMap["cluster_name"].(string))
	}
	if sharedTargetDataMap["cluster_type"] != nil {
		sharedTargetData.ClusterType = core.StringPtr(sharedTargetDataMap["cluster_type"].(string))
	}
	if sharedTargetDataMap["entitlement_keys"] != nil {
		entitlementKeys := []interface{}{}
		for _, entitlementKeysItem := range sharedTargetDataMap["entitlement_keys"].([]interface{}) {
			entitlementKeys = append(entitlementKeys, entitlementKeysItem.(interface{}))
		}
		sharedTargetData.EntitlementKeys = entitlementKeys
	}
	if sharedTargetDataMap["namespace"] != nil {
		sharedTargetData.Namespace = core.StringPtr(sharedTargetDataMap["namespace"].(string))
	}
	if sharedTargetDataMap["region"] != nil {
		sharedTargetData.Region = core.StringPtr(sharedTargetDataMap["region"].(string))
	}
	if sharedTargetDataMap["resource_group_id"] != nil {
		sharedTargetData.ResourceGroupID = core.StringPtr(sharedTargetDataMap["resource_group_id"].(string))
	}
	if sharedTargetDataMap["worker_count"] != nil {
		sharedTargetData.WorkerCount = core.Int64Ptr(int64(sharedTargetDataMap["worker_count"].(int)))
	}
	if sharedTargetDataMap["worker_machine_type"] != nil {
		sharedTargetData.WorkerMachineType = core.StringPtr(sharedTargetDataMap["worker_machine_type"].(string))
	}

	return sharedTargetData
}

func resourceIBMSchematicsWorkspaceMapToTemplateSourceDataRequest(templateSourceDataRequestMap map[string]interface{}) schematicsv1.TemplateSourceDataRequest {
	templateSourceDataRequest := schematicsv1.TemplateSourceDataRequest{}

	if templateSourceDataRequestMap["env_values"] != nil {
		envValues := []interface{}{}
		for _, envValuesItem := range templateSourceDataRequestMap["env_values"].([]interface{}) {
			envValues = append(envValues, envValuesItem.(interface{}))
		}
		templateSourceDataRequest.EnvValues = envValues
	}
	if templateSourceDataRequestMap["folder"] != nil {
		templateSourceDataRequest.Folder = core.StringPtr(templateSourceDataRequestMap["folder"].(string))
	}
	if templateSourceDataRequestMap["init_state_file"] != nil {
		templateSourceDataRequest.InitStateFile = core.StringPtr(templateSourceDataRequestMap["init_state_file"].(string))
	}
	if templateSourceDataRequestMap["type"] != nil {
		templateSourceDataRequest.Type = core.StringPtr(templateSourceDataRequestMap["type"].(string))
	}
	if templateSourceDataRequestMap["uninstall_script_name"] != nil {
		templateSourceDataRequest.UninstallScriptName = core.StringPtr(templateSourceDataRequestMap["uninstall_script_name"].(string))
	}
	if templateSourceDataRequestMap["values"] != nil {
		templateSourceDataRequest.Values = core.StringPtr(templateSourceDataRequestMap["values"].(string))
	}
	if templateSourceDataRequestMap["values_metadata"] != nil {
		valuesMetadata := []interface{}{}
		for _, valuesMetadataItem := range templateSourceDataRequestMap["values_metadata"].([]interface{}) {
			valuesMetadata = append(valuesMetadata, valuesMetadataItem.(interface{}))
		}
		templateSourceDataRequest.ValuesMetadata = valuesMetadata
	}
	if templateSourceDataRequestMap["variablestore"] != nil {
		variablestore := []schematicsv1.WorkspaceVariableRequest{}
		for _, variablestoreItem := range templateSourceDataRequestMap["variablestore"].([]interface{}) {
			variablestoreItemModel := resourceIBMSchematicsWorkspaceMapToWorkspaceVariableRequest(variablestoreItem.(map[string]interface{}))
			variablestore = append(variablestore, variablestoreItemModel)
		}
		templateSourceDataRequest.Variablestore = variablestore
	}

	return templateSourceDataRequest
}

func resourceIBMSchematicsWorkspaceMapToWorkspaceVariableRequest(workspaceVariableRequestMap map[string]interface{}) schematicsv1.WorkspaceVariableRequest {
	workspaceVariableRequest := schematicsv1.WorkspaceVariableRequest{}

	if workspaceVariableRequestMap["description"] != nil {
		workspaceVariableRequest.Description = core.StringPtr(workspaceVariableRequestMap["description"].(string))
	}
	if workspaceVariableRequestMap["name"] != nil {
		workspaceVariableRequest.Name = core.StringPtr(workspaceVariableRequestMap["name"].(string))
	}
	if workspaceVariableRequestMap["secure"] != nil {
		workspaceVariableRequest.Secure = core.BoolPtr(workspaceVariableRequestMap["secure"].(bool))
	}
	if workspaceVariableRequestMap["type"] != nil {
		workspaceVariableRequest.Type = core.StringPtr(workspaceVariableRequestMap["type"].(string))
	}
	if workspaceVariableRequestMap["use_default"] != nil {
		workspaceVariableRequest.UseDefault = core.BoolPtr(workspaceVariableRequestMap["use_default"].(bool))
	}
	if workspaceVariableRequestMap["value"] != nil {
		workspaceVariableRequest.Value = core.StringPtr(workspaceVariableRequestMap["value"].(string))
	}

	return workspaceVariableRequest
}

func resourceIBMSchematicsWorkspaceMapToTemplateRepoRequest(templateRepoRequestMap map[string]interface{}) schematicsv1.TemplateRepoRequest {
	templateRepoRequest := schematicsv1.TemplateRepoRequest{}

	if templateRepoRequestMap["branch"] != nil {
		templateRepoRequest.Branch = core.StringPtr(templateRepoRequestMap["branch"].(string))
	}
	if templateRepoRequestMap["release"] != nil {
		templateRepoRequest.Release = core.StringPtr(templateRepoRequestMap["release"].(string))
	}
	if templateRepoRequestMap["repo_sha_value"] != nil {
		templateRepoRequest.RepoShaValue = core.StringPtr(templateRepoRequestMap["repo_sha_value"].(string))
	}
	if templateRepoRequestMap["repo_url"] != nil {
		templateRepoRequest.RepoURL = core.StringPtr(templateRepoRequestMap["repo_url"].(string))
	}
	if templateRepoRequestMap["url"] != nil {
		templateRepoRequest.URL = core.StringPtr(templateRepoRequestMap["url"].(string))
	}

	return templateRepoRequest
}

func resourceIBMSchematicsWorkspaceMapToTemplateRepoUpdateRequest(templateRepoUpdateRequestMap map[string]interface{}) schematicsv1.TemplateRepoUpdateRequest {
	templateRepoUpdateRequest := schematicsv1.TemplateRepoUpdateRequest{}

	if templateRepoUpdateRequestMap["branch"] != nil {
		templateRepoUpdateRequest.Branch = core.StringPtr(templateRepoUpdateRequestMap["branch"].(string))
	}
	if templateRepoUpdateRequestMap["release"] != nil {
		templateRepoUpdateRequest.Release = core.StringPtr(templateRepoUpdateRequestMap["release"].(string))
	}
	if templateRepoUpdateRequestMap["repo_sha_value"] != nil {
		templateRepoUpdateRequest.RepoShaValue = core.StringPtr(templateRepoUpdateRequestMap["repo_sha_value"].(string))
	}
	if templateRepoUpdateRequestMap["repo_url"] != nil {
		templateRepoUpdateRequest.RepoURL = core.StringPtr(templateRepoUpdateRequestMap["repo_url"].(string))
	}
	if templateRepoUpdateRequestMap["url"] != nil {
		templateRepoUpdateRequest.URL = core.StringPtr(templateRepoUpdateRequestMap["url"].(string))
	}

	return templateRepoUpdateRequest
}

func resourceIBMSchematicsWorkspaceMapToWorkspaceStatusRequest(workspaceStatusRequestMap map[string]interface{}) schematicsv1.WorkspaceStatusRequest {
	workspaceStatusRequest := schematicsv1.WorkspaceStatusRequest{}

	if workspaceStatusRequestMap["frozen"] != nil {
		workspaceStatusRequest.Frozen = core.BoolPtr(workspaceStatusRequestMap["frozen"].(bool))
	}
	if workspaceStatusRequestMap["frozen_at"] != nil {
		frozenAt, err := strfmt.ParseDateTime(workspaceStatusRequestMap["frozen_at"].(string))
		if err != nil {
			workspaceStatusRequest.FrozenAt = &frozenAt
		}
	}
	if workspaceStatusRequestMap["frozen_by"] != nil {
		workspaceStatusRequest.FrozenBy = core.StringPtr(workspaceStatusRequestMap["frozen_by"].(string))
	}
	if workspaceStatusRequestMap["locked"] != nil {
		workspaceStatusRequest.Locked = core.BoolPtr(workspaceStatusRequestMap["locked"].(bool))
	}
	if workspaceStatusRequestMap["locked_by"] != nil {
		workspaceStatusRequest.LockedBy = core.StringPtr(workspaceStatusRequestMap["locked_by"].(string))
	}
	if workspaceStatusRequestMap["locked_time"] != nil {
		lockedTime, err := strfmt.ParseDateTime(workspaceStatusRequestMap["locked_time"].(string))
		if err != nil {
			workspaceStatusRequest.LockedTime = &lockedTime
		}
	}

	return workspaceStatusRequest
}

func resourceIBMSchematicsWorkspaceMapToWorkspaceStatusUpdateRequest(workspaceStatusUpdateRequestMap map[string]interface{}) schematicsv1.WorkspaceStatusUpdateRequest {
	workspaceStatusUpdateRequest := schematicsv1.WorkspaceStatusUpdateRequest{}

	if workspaceStatusUpdateRequestMap["frozen"] != nil {
		workspaceStatusUpdateRequest.Frozen = core.BoolPtr(workspaceStatusUpdateRequestMap["frozen"].(bool))
	}
	if workspaceStatusUpdateRequestMap["frozen_at"] != nil {
		frozenAt := workspaceStatusUpdateRequestMap["frozen_at"].(strfmt.DateTime)
		workspaceStatusUpdateRequest.FrozenAt = &frozenAt
	}
	if workspaceStatusUpdateRequestMap["frozen_by"] != nil {
		workspaceStatusUpdateRequest.FrozenBy = core.StringPtr(workspaceStatusUpdateRequestMap["frozen_by"].(string))
	}
	if workspaceStatusUpdateRequestMap["locked"] != nil {
		workspaceStatusUpdateRequest.Locked = core.BoolPtr(workspaceStatusUpdateRequestMap["locked"].(bool))
	}
	if workspaceStatusUpdateRequestMap["locked_by"] != nil {
		workspaceStatusUpdateRequest.LockedBy = core.StringPtr(workspaceStatusUpdateRequestMap["locked_by"].(string))
	}
	if workspaceStatusUpdateRequestMap["locked_time"] != nil {
		lockedTime := workspaceStatusUpdateRequestMap["locked_time"].(strfmt.DateTime)
		workspaceStatusUpdateRequest.LockedTime = &lockedTime
	}

	return workspaceStatusUpdateRequest
}

func resourceIBMSchematicsWorkspaceRead(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	getWorkspaceOptions := &schematicsv1.GetWorkspaceOptions{}

	getWorkspaceOptions.SetWID(d.Id())

	workspaceResponse, response, err := schematicsClient.GetWorkspaceWithContext(context.TODO(), getWorkspaceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetWorkspaceWithContext failed %s\n%s", err, response)
		return err
	}

	if workspaceResponse.AppliedShareddataIds != nil {
		if err = d.Set("applied_shareddata_ids", workspaceResponse.AppliedShareddataIds); err != nil {
			return fmt.Errorf("Error setting applied_shareddata_ids: %s", err)
		}
	}
	if workspaceResponse.CatalogRef != nil {
		catalogRefMap := resourceIBMSchematicsWorkspaceCatalogRefToMap(*workspaceResponse.CatalogRef)
		if err = d.Set("catalog_ref", []map[string]interface{}{catalogRefMap}); err != nil {
			return fmt.Errorf("Error setting catalog_ref: %s", err)
		}
	}
	if err = d.Set("description", workspaceResponse.Description); err != nil {
		return fmt.Errorf("Error setting description: %s", err)
	}
	if err = d.Set("location", workspaceResponse.Location); err != nil {
		return fmt.Errorf("Error setting location: %s", err)
	}
	if err = d.Set("name", workspaceResponse.Name); err != nil {
		return fmt.Errorf("Error setting name: %s", err)
	}
	if err = d.Set("resource_group", workspaceResponse.ResourceGroup); err != nil {
		return fmt.Errorf("Error setting resource_group: %s", err)
	}
	if _, ok := d.GetOk("shared_data"); ok {
		if workspaceResponse.SharedData != nil {
			sharedDataMap := resourceIBMSchematicsWorkspaceSharedTargetDataResponseToMap(*workspaceResponse.SharedData)
			if err = d.Set("shared_data", []map[string]interface{}{sharedDataMap}); err != nil {
				return fmt.Errorf("Error reading shared_data: %s", err)
			}
		}
	}
	if workspaceResponse.Tags != nil {
		if err = d.Set("tags", workspaceResponse.Tags); err != nil {
			return fmt.Errorf("Error setting tags: %s", err)
		}
	}
	if workspaceResponse.TemplateData != nil {
		templateData := []map[string]interface{}{}
		for _, templateDataItem := range workspaceResponse.TemplateData {
			templateDataItemMap := resourceIBMSchematicsWorkspaceTemplateSourceDataResponseToMap(templateDataItem)
			templateData = append(templateData, templateDataItemMap)
		}
		if err = d.Set("template_env_settings", templateData[0]["env_values"]); err != nil {
			return fmt.Errorf("Error reading env_values: %s", err)
		}
		if err = d.Set("template_git_folder", templateData[0]["folder"]); err != nil {
			return fmt.Errorf("Error reading folder: %s", err)
		}
		if err = d.Set("template_init_state_file", templateData[0]["init_state_file"]); err != nil {
			return fmt.Errorf("Error reading init_state_file: %s", err)
		}
		if err = d.Set("template_type", templateData[0]["type"]); err != nil {
			return fmt.Errorf("Error reading type: %s", err)
		}
		if err = d.Set("template_uninstall_script_name", templateData[0]["uninstall_script_name"]); err != nil {
			return fmt.Errorf("Error reading uninstall_script_name: %s", err)
		}
		if err = d.Set("template_values", templateData[0]["values"]); err != nil {
			return fmt.Errorf("Error reading values: %s", err)
		}
		if err = d.Set("template_values_metadata", templateData[0]["values_metadata"]); err != nil {
			return fmt.Errorf("Error reading values_metadata: %s", err)
		}
		if err = d.Set("template_inputs", templateData[0]["variablestore"]); err != nil {
			return fmt.Errorf("Error reading variablestore: %s", err)
		}

	}
	if err = d.Set("template_ref", workspaceResponse.TemplateRef); err != nil {
		return fmt.Errorf("Error setting template_ref: %s", err)
	}
	if workspaceResponse.TemplateRepo != nil {
		templateRepoMap := resourceIBMSchematicsWorkspaceTemplateRepoResponseToMap(*workspaceResponse.TemplateRepo)
		if err = d.Set("template_git_branch", templateRepoMap["branch"]); err != nil {
			return fmt.Errorf("Error reading branch: %s", err)
		}
		if err = d.Set("template_git_release", templateRepoMap["release"]); err != nil {
			return fmt.Errorf("Error reading release: %s", err)
		}
		if err = d.Set("template_git_repo_sha_value", templateRepoMap["repo_sha_value"]); err != nil {
			return fmt.Errorf("Error reading repo_sha_value: %s", err)
		}
		if err = d.Set("template_git_repo_url", templateRepoMap["repo_url"]); err != nil {
			return fmt.Errorf("Error reading repo_url: %s", err)
		}
		if err = d.Set("template_git_url", templateRepoMap["url"]); err != nil {
			return fmt.Errorf("Error reading url: %s", err)
		}
		if err = d.Set("template_git_has_uploadedgitrepotar", templateRepoMap["has_uploadedgitrepotar"]); err != nil {
			return fmt.Errorf("Error reading has_uploadedgitrepotar: %s", err)
		}
	}
	/*if workspaceResponse.Type != nil {
		if err = d.Set("template_type", workspaceResponse.Type); err != nil {
			return fmt.Errorf("Error reading type: %s", err)
		}
	}*/
	if workspaceResponse.WorkspaceStatus != nil {
		workspaceStatusMap := resourceIBMSchematicsWorkspaceWorkspaceStatusResponseToMap(*workspaceResponse.WorkspaceStatus)
		if err = d.Set("frozen", workspaceStatusMap["frozen"]); err != nil {
			return fmt.Errorf("Error reading frozen: %s", err)
		}
		if err = d.Set("frozen_at", workspaceStatusMap["frozen_at"]); err != nil {
			return fmt.Errorf("Error reading frozen_at: %s", err)
		}
		if err = d.Set("frozen_by", workspaceStatusMap["frozen_by"]); err != nil {
			return fmt.Errorf("Error reading frozen_by: %s", err)
		}
		if err = d.Set("locked", workspaceStatusMap["locked"]); err != nil {
			return fmt.Errorf("Error reading locked: %s", err)
		}
		if err = d.Set("locked_by", workspaceStatusMap["locked_by"]); err != nil {
			return fmt.Errorf("Error reading locked_by: %s", err)
		}
		if err = d.Set("locked_time", workspaceStatusMap["locked_time"]); err != nil {
			return fmt.Errorf("Error reading locked_time: %s", err)
		}
	}
	if workspaceResponse.CreatedAt != nil {
		if err = d.Set("created_at", workspaceResponse.CreatedAt.String()); err != nil {
			return fmt.Errorf("Error reading created_at: %s", err)
		}
	}
	if err = d.Set("created_by", workspaceResponse.CreatedBy); err != nil {
		return fmt.Errorf("Error setting created_by: %s", err)
	}
	if err = d.Set("crn", workspaceResponse.Crn); err != nil {
		return fmt.Errorf("Error reading crn: %s", err)
	}
	if workspaceResponse.LastHealthCheckAt != nil {
		if err = d.Set("last_health_check_at", workspaceResponse.LastHealthCheckAt.String()); err != nil {
			return fmt.Errorf("Error reading last_health_check_at: %s", err)
		}
	}
	if workspaceResponse.RuntimeData != nil {
		runtimeData := []map[string]interface{}{}
		for _, runtimeDataItem := range workspaceResponse.RuntimeData {
			runtimeDataItemMap := resourceIBMSchematicsWorkspaceTemplateRunTimeDataResponseToMap(runtimeDataItem)
			runtimeData = append(runtimeData, runtimeDataItemMap)
		}
		if err = d.Set("runtime_data", runtimeData); err != nil {
			return fmt.Errorf("Error setting runtime_data: %s", err)
		}
	}
	if err = d.Set("status", workspaceResponse.Status); err != nil {
		return fmt.Errorf("Error setting status: %s", err)
	}
	if workspaceResponse.UpdatedAt != nil {
		if err = d.Set("updated_at", workspaceResponse.UpdatedAt.String()); err != nil {
			return fmt.Errorf("Error reading updated_at: %s", err)
		}
	}
	if err = d.Set("updated_by", workspaceResponse.UpdatedBy); err != nil {
		return fmt.Errorf("Error setting updated_by: %s", err)
	}
	if workspaceResponse.WorkspaceStatusMsg != nil {
		workspaceStatusMsgMap := resourceIBMSchematicsWorkspaceWorkspaceStatusMessageToMap(*workspaceResponse.WorkspaceStatusMsg)
		if err = d.Set("status_code", workspaceStatusMsgMap["status_code"]); err != nil {
			return fmt.Errorf("Error reading status_code: %s", err)
		}
		if err = d.Set("status_msg", workspaceStatusMsgMap["status_msg"]); err != nil {
			return fmt.Errorf("Error reading status_msg: %s", err)
		}
	}

	return nil
}

func resourceIBMSchematicsWorkspaceCatalogRefToMap(catalogRef schematicsv1.CatalogRef) map[string]interface{} {
	catalogRefMap := map[string]interface{}{}

	catalogRefMap["dry_run"] = catalogRef.DryRun
	catalogRefMap["item_icon_url"] = catalogRef.ItemIconURL
	catalogRefMap["item_id"] = catalogRef.ItemID
	catalogRefMap["item_name"] = catalogRef.ItemName
	catalogRefMap["item_readme_url"] = catalogRef.ItemReadmeURL
	catalogRefMap["item_url"] = catalogRef.ItemURL
	catalogRefMap["launch_url"] = catalogRef.LaunchURL
	catalogRefMap["offering_version"] = catalogRef.OfferingVersion

	return catalogRefMap
}

func resourceIBMSchematicsWorkspaceSharedTargetDataToMap(sharedTargetData schematicsv1.SharedTargetData) map[string]interface{} {
	sharedTargetDataMap := map[string]interface{}{}

	sharedTargetDataMap["cluster_created_on"] = sharedTargetData.ClusterCreatedOn
	sharedTargetDataMap["cluster_id"] = sharedTargetData.ClusterID
	sharedTargetDataMap["cluster_name"] = sharedTargetData.ClusterName
	sharedTargetDataMap["cluster_type"] = sharedTargetData.ClusterType
	if sharedTargetData.EntitlementKeys != nil {
		entitlementKeys := []interface{}{}
		for _, entitlementKeysItem := range sharedTargetData.EntitlementKeys {
			entitlementKeys = append(entitlementKeys, entitlementKeysItem)
		}
		sharedTargetDataMap["entitlement_keys"] = entitlementKeys
	}
	sharedTargetDataMap["namespace"] = sharedTargetData.Namespace
	sharedTargetDataMap["region"] = sharedTargetData.Region
	sharedTargetDataMap["resource_group_id"] = sharedTargetData.ResourceGroupID
	sharedTargetDataMap["worker_count"] = intValue(sharedTargetData.WorkerCount)
	sharedTargetDataMap["worker_machine_type"] = sharedTargetData.WorkerMachineType

	return sharedTargetDataMap
}

func resourceIBMSchematicsWorkspaceSharedTargetDataResponseToMap(sharedTargetData schematicsv1.SharedTargetDataResponse) map[string]interface{} {
	sharedTargetDataResponseMap := map[string]interface{}{}

	sharedTargetDataResponseMap["cluster_id"] = sharedTargetData.ClusterID
	sharedTargetDataResponseMap["cluster_name"] = sharedTargetData.ClusterName
	if sharedTargetData.EntitlementKeys != nil {
		entitlementKeys := []interface{}{}
		for _, entitlementKeysItem := range sharedTargetData.EntitlementKeys {
			entitlementKeys = append(entitlementKeys, entitlementKeysItem)
		}
		sharedTargetDataResponseMap["entitlement_keys"] = entitlementKeys
	}
	sharedTargetDataResponseMap["namespace"] = sharedTargetData.Namespace
	sharedTargetDataResponseMap["region"] = sharedTargetData.Region
	sharedTargetDataResponseMap["resource_group_id"] = sharedTargetData.ResourceGroupID

	return sharedTargetDataResponseMap
}

func resourceIBMSchematicsWorkspaceTemplateSourceDataRequestToMap(templateSourceDataRequest schematicsv1.TemplateSourceDataRequest) map[string]interface{} {
	templateSourceDataRequestMap := map[string]interface{}{}

	if templateSourceDataRequest.EnvValues != nil {
		envValues := []interface{}{}
		for _, envValuesItem := range templateSourceDataRequest.EnvValues {
			envValues = append(envValues, envValuesItem)
		}
		templateSourceDataRequestMap["env_values"] = envValues
	}
	templateSourceDataRequestMap["folder"] = templateSourceDataRequest.Folder
	templateSourceDataRequestMap["init_state_file"] = templateSourceDataRequest.InitStateFile
	templateSourceDataRequestMap["type"] = templateSourceDataRequest.Type
	templateSourceDataRequestMap["uninstall_script_name"] = templateSourceDataRequest.UninstallScriptName
	templateSourceDataRequestMap["values"] = templateSourceDataRequest.Values
	if templateSourceDataRequest.ValuesMetadata != nil {
		valuesMetadata := []interface{}{}
		for _, valuesMetadataItem := range templateSourceDataRequest.ValuesMetadata {
			valuesMetadata = append(valuesMetadata, valuesMetadataItem)
		}
		templateSourceDataRequestMap["values_metadata"] = valuesMetadata
	}
	if templateSourceDataRequest.Variablestore != nil {
		variablestore := []map[string]interface{}{}
		for _, variablestoreItem := range templateSourceDataRequest.Variablestore {
			variablestoreItemMap := resourceIBMSchematicsWorkspaceWorkspaceVariableRequestToMap(variablestoreItem)
			variablestore = append(variablestore, variablestoreItemMap)
			// TODO: handle Variablestore of type TypeList -- list of non-primitive, not model items
		}
		templateSourceDataRequestMap["variablestore"] = variablestore
	}

	return templateSourceDataRequestMap
}

func resourceIBMSchematicsWorkspaceTemplateSourceDataResponseToMap(templateSourceDataResponse schematicsv1.TemplateSourceDataResponse) map[string]interface{} {
	templateSourceDataResponseMap := map[string]interface{}{}

	if templateSourceDataResponse.EnvValues != nil {
		envValues := []map[string]interface{}{}
		for _, envValuesItem := range templateSourceDataResponse.EnvValues {
			flattenedEnvVals := map[string]interface{}{}
			if envValuesItem.Name != nil {
				flattenedEnvVals[*envValuesItem.Name] = envValuesItem.Value
			}

			envValues = append(envValues, flattenedEnvVals)
		}
		templateSourceDataResponseMap["env_values"] = envValues
	}
	if templateSourceDataResponse.Type != nil {
		templateSourceDataResponseMap["type"] = templateSourceDataResponse.Type
	}
	templateSourceDataResponseMap["folder"] = templateSourceDataResponse.Folder
	templateSourceDataResponseMap["uninstall_script_name"] = templateSourceDataResponse.UninstallScriptName
	templateSourceDataResponseMap["values"] = templateSourceDataResponse.Values
	if templateSourceDataResponse.ValuesMetadata != nil {
		valuesMetadata := []interface{}{}
		for _, valuesMetadataItem := range templateSourceDataResponse.ValuesMetadata {
			valuesMetadata = append(valuesMetadata, valuesMetadataItem)
		}
		templateSourceDataResponseMap["values_metadata"] = valuesMetadata
	}
	if templateSourceDataResponse.Variablestore != nil {
		variablestore := []map[string]interface{}{}
		for _, variablestoreItem := range templateSourceDataResponse.Variablestore {
			variablestoreItemMap := resourceIBMSchematicsWorkspaceWorkspaceVariableResponseToMap(variablestoreItem)
			variablestore = append(variablestore, variablestoreItemMap)
		}
		templateSourceDataResponseMap["variablestore"] = variablestore
	}

	return templateSourceDataResponseMap
}

func resourceIBMSchematicsWorkspaceWorkspaceVariableRequestToMap(workspaceVariableRequest schematicsv1.WorkspaceVariableRequest) map[string]interface{} {
	workspaceVariableRequestMap := map[string]interface{}{}

	workspaceVariableRequestMap["description"] = workspaceVariableRequest.Description
	workspaceVariableRequestMap["name"] = workspaceVariableRequest.Name
	workspaceVariableRequestMap["secure"] = workspaceVariableRequest.Secure
	workspaceVariableRequestMap["type"] = workspaceVariableRequest.Type
	workspaceVariableRequestMap["use_default"] = workspaceVariableRequest.UseDefault
	workspaceVariableRequestMap["value"] = workspaceVariableRequest.Value

	return workspaceVariableRequestMap
}

func resourceIBMSchematicsWorkspaceWorkspaceVariableResponseToMap(workspaceVariableResponse schematicsv1.WorkspaceVariableResponse) map[string]interface{} {
	workspaceVariableRequestMap := map[string]interface{}{}

	workspaceVariableRequestMap["description"] = workspaceVariableResponse.Description
	workspaceVariableRequestMap["name"] = workspaceVariableResponse.Name
	workspaceVariableRequestMap["secure"] = workspaceVariableResponse.Secure
	workspaceVariableRequestMap["type"] = workspaceVariableResponse.Type
	workspaceVariableRequestMap["value"] = workspaceVariableResponse.Value

	return workspaceVariableRequestMap
}

func resourceIBMSchematicsWorkspaceTemplateRepoRequestToMap(templateRepoRequest schematicsv1.TemplateRepoRequest) map[string]interface{} {
	templateRepoRequestMap := map[string]interface{}{}

	templateRepoRequestMap["branch"] = templateRepoRequest.Branch
	templateRepoRequestMap["release"] = templateRepoRequest.Release
	templateRepoRequestMap["repo_sha_value"] = templateRepoRequest.RepoShaValue
	templateRepoRequestMap["repo_url"] = templateRepoRequest.RepoURL
	templateRepoRequestMap["url"] = templateRepoRequest.URL

	return templateRepoRequestMap
}

func resourceIBMSchematicsWorkspaceTemplateRepoResponseToMap(templateRepoResponse schematicsv1.TemplateRepoResponse) map[string]interface{} {
	templateRepoResponseMap := map[string]interface{}{}

	templateRepoResponseMap["branch"] = templateRepoResponse.Branch
	templateRepoResponseMap["release"] = templateRepoResponse.Release
	templateRepoResponseMap["repo_sha_value"] = templateRepoResponse.RepoShaValue
	templateRepoResponseMap["repo_url"] = templateRepoResponse.RepoURL
	templateRepoResponseMap["url"] = templateRepoResponse.URL
	templateRepoResponseMap["has_uploadedgitrepotar"] = templateRepoResponse.HasUploadedgitrepotar

	return templateRepoResponseMap
}

func resourceIBMSchematicsWorkspaceWorkspaceStatusRequestToMap(workspaceStatusRequest schematicsv1.WorkspaceStatusRequest) map[string]interface{} {
	workspaceStatusRequestMap := map[string]interface{}{}

	workspaceStatusRequestMap["frozen"] = workspaceStatusRequest.Frozen
	workspaceStatusRequestMap["frozen_at"] = workspaceStatusRequest.FrozenAt.String()
	workspaceStatusRequestMap["frozen_by"] = workspaceStatusRequest.FrozenBy
	workspaceStatusRequestMap["locked"] = workspaceStatusRequest.Locked
	workspaceStatusRequestMap["locked_by"] = workspaceStatusRequest.LockedBy
	workspaceStatusRequestMap["locked_time"] = workspaceStatusRequest.LockedTime.String()

	return workspaceStatusRequestMap
}

func resourceIBMSchematicsWorkspaceWorkspaceStatusResponseToMap(workspaceStatusResponse schematicsv1.WorkspaceStatusResponse) map[string]interface{} {
	workspaceStatusResponseMap := map[string]interface{}{}

	workspaceStatusResponseMap["frozen"] = workspaceStatusResponse.Frozen
	if workspaceStatusResponse.FrozenAt != nil {
		workspaceStatusResponseMap["frozen_at"] = workspaceStatusResponse.FrozenAt.String()
	}
	workspaceStatusResponseMap["frozen_by"] = workspaceStatusResponse.FrozenBy
	workspaceStatusResponseMap["locked"] = workspaceStatusResponse.Locked
	workspaceStatusResponseMap["locked_by"] = workspaceStatusResponse.LockedBy
	if workspaceStatusResponse.LockedTime != nil {
		workspaceStatusResponseMap["locked_time"] = workspaceStatusResponse.LockedTime.String()
	}

	return workspaceStatusResponseMap
}

func resourceIBMSchematicsWorkspaceTemplateRunTimeDataResponseToMap(templateRunTimeDataResponse schematicsv1.TemplateRunTimeDataResponse) map[string]interface{} {
	templateRunTimeDataResponseMap := map[string]interface{}{}

	templateRunTimeDataResponseMap["engine_cmd"] = templateRunTimeDataResponse.EngineCmd
	templateRunTimeDataResponseMap["engine_name"] = templateRunTimeDataResponse.EngineName
	templateRunTimeDataResponseMap["engine_version"] = templateRunTimeDataResponse.EngineVersion
	templateRunTimeDataResponseMap["id"] = templateRunTimeDataResponse.ID
	templateRunTimeDataResponseMap["log_store_url"] = templateRunTimeDataResponse.LogStoreURL
	if templateRunTimeDataResponse.OutputValues != nil {
		outputValues := []interface{}{}
		for _, outputValuesItem := range templateRunTimeDataResponse.OutputValues {
			outputValues = append(outputValues, outputValuesItem)
		}
		templateRunTimeDataResponseMap["output_values"] = outputValues
	}
	if templateRunTimeDataResponse.Resources != nil {
		resources := []interface{}{}
		for _, resourcesItem := range templateRunTimeDataResponse.Resources {
			resources = append(resources, resourcesItem)
		}
		templateRunTimeDataResponseMap["resources"] = resources
	}
	templateRunTimeDataResponseMap["state_store_url"] = templateRunTimeDataResponse.StateStoreURL

	return templateRunTimeDataResponseMap
}

func resourceIBMSchematicsWorkspaceWorkspaceStatusMessageToMap(workspaceStatusMessage schematicsv1.WorkspaceStatusMessage) map[string]interface{} {
	workspaceStatusMessageMap := map[string]interface{}{}

	workspaceStatusMessageMap["status_code"] = workspaceStatusMessage.StatusCode
	workspaceStatusMessageMap["status_msg"] = workspaceStatusMessage.StatusMsg

	return workspaceStatusMessageMap
}

func resourceIBMSchematicsWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	updateWorkspaceOptions := &schematicsv1.UpdateWorkspaceOptions{}

	updateWorkspaceOptions.SetWID(d.Id())

	hasChange := false

	if d.HasChange("catalog_ref") {
		catalogRefAttr := d.Get("catalog_ref").([]interface{})
		if len(catalogRefAttr) > 0 {
			catalogRef := resourceIBMSchematicsWorkspaceMapToCatalogRef(d.Get("catalog_ref.0").(map[string]interface{}))
			updateWorkspaceOptions.SetCatalogRef(&catalogRef)
			hasChange = true
		}
	}
	if d.HasChange("description") {
		updateWorkspaceOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("name") {
		updateWorkspaceOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("shared_data") {
		sharedDataAttr := d.Get("shared_data").([]interface{})
		if len(sharedDataAttr) > 0 {
			sharedData := resourceIBMSchematicsWorkspaceMapToSharedTargetData(d.Get("shared_data.0").(map[string]interface{}))
			updateWorkspaceOptions.SetSharedData(&sharedData)
			hasChange = true
		}
	}
	if d.HasChange("tags") {
		updateWorkspaceOptions.SetTags(expandStringList(d.Get("tags").([]interface{})))
		hasChange = true
	}

	var templateData []schematicsv1.TemplateSourceDataRequest

	templateSourceDataRequestMap := map[string]interface{}{}
	hasTemplateData := false

	if d.HasChange("template_env_settings") {
		templateSourceDataRequestMap["env_values"] = d.Get("template_env_settings").([]interface{})
		hasTemplateData = true
	}
	if d.HasChange("template_git_folder") {
		templateSourceDataRequestMap["folder"] = d.Get("template_git_folder").(string)
		hasTemplateData = true
	}
	if d.HasChange("template_init_state_file") {
		templateSourceDataRequestMap["init_state_file"] = d.Get("template_init_state_file").(string)
		hasTemplateData = true
	}
	if d.HasChange("template_type") {
		templateSourceDataRequestMap["type"] = d.Get("template_type").(string)
		updateWorkspaceOptions.SetType([]string{d.Get("template_type").(string)})
		hasTemplateData = true
	}
	if d.HasChange("template_uninstall_script_name") {
		templateSourceDataRequestMap["uninstall_script_name"] = d.Get("template_uninstall_script_name").(string)
		hasTemplateData = true
	}
	if d.HasChange("template_values") {
		templateSourceDataRequestMap["values"] = d.Get("template_values").(string)
		hasTemplateData = true
	}
	if d.HasChange("template_values_metadata") {
		templateSourceDataRequestMap["values_metadata"] = d.Get("template_values_metadata").([]interface{})
		hasTemplateData = true
	}
	if d.HasChange("template_inputs") {
		templateSourceDataRequestMap["variablestore"] = d.Get("template_inputs").([]interface{})
		hasTemplateData = true
	}
	if hasTemplateData {
		templateDataItem := resourceIBMSchematicsWorkspaceMapToTemplateSourceDataRequest(templateSourceDataRequestMap)
		templateData = append(templateData, templateDataItem)
		updateWorkspaceOptions.SetTemplateData(templateData)
		hasChange = true
	}

	templateRepoRequestMap := map[string]interface{}{}
	hasTemplateRepo := false
	if d.HasChange("template_git_branch") {
		templateRepoRequestMap["branch"] = d.Get("template_git_branch").(bool)
		hasTemplateRepo = true
	}
	if d.HasChange("template_git_release") {
		templateRepoRequestMap["release"] = d.Get("template_git_release").(string)
		hasTemplateRepo = true
	}
	if d.HasChange("template_git_repo_sha_value") {
		templateRepoRequestMap["repo_sha_value"] = d.Get("template_git_repo_sha_value").(string)
		hasTemplateRepo = true
	}
	if d.HasChange("template_git_repo_url") {
		templateRepoRequestMap["repo_url"] = d.Get("template_git_repo_url").(string)
		hasTemplateRepo = true
	}
	if d.HasChange("template_git_url") {
		templateRepoRequestMap["url"] = d.Get("template_git_url").(string)
		hasTemplateRepo = true
	}
	if d.HasChange("template_git_has_uploadedgitrepotar") {
		templateRepoRequestMap["has_uploadedgitrepotar"] = d.Get("template_git_has_uploadedgitrepotar").(string)
		hasTemplateRepo = true
	}
	if hasTemplateRepo {
		templateRepo := resourceIBMSchematicsWorkspaceMapToTemplateRepoUpdateRequest(templateRepoRequestMap)
		updateWorkspaceOptions.SetTemplateRepo(&templateRepo)
		hasChange = true
	}

	if d.HasChange("template_type") {
		updateWorkspaceOptions.SetType([]string{d.Get("template_type").(string)})
		hasChange = true
	}

	workspaceStatusRequestMap := map[string]interface{}{}
	workspaceStatus := false
	if d.HasChange("frozen") {
		workspaceStatusRequestMap["frozen"] = d.Get("frozen").(bool)
		workspaceStatus = true
	}
	if d.HasChange("frozen_at") {
		workspaceStatusRequestMap["frozen_at"] = d.Get("frozen_at").(string)
		workspaceStatus = true
	}
	if d.HasChange("frozen_by") {
		workspaceStatusRequestMap["frozen_by"] = d.Get("frozen_by").(string)
		workspaceStatus = true
	}
	if d.HasChange("locked") {
		workspaceStatusRequestMap["locked"] = d.Get("locked").(bool)
		workspaceStatus = true
	}
	if d.HasChange("locked_by") {
		workspaceStatusRequestMap["locked_by"] = d.Get("locked_by").(string)
		workspaceStatus = true
	}
	if d.HasChange("locked_time") {
		workspaceStatusRequestMap["locked_time"] = d.Get("locked_time").(string)
		workspaceStatus = true
	}
	if workspaceStatus {
		workspaceStatus := resourceIBMSchematicsWorkspaceMapToWorkspaceStatusUpdateRequest(workspaceStatusRequestMap)
		updateWorkspaceOptions.SetWorkspaceStatus(&workspaceStatus)
		hasChange = true
	}

	if hasChange {
		_, response, err := schematicsClient.UpdateWorkspaceWithContext(context.TODO(), updateWorkspaceOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateWorkspaceWithContext failed %s\n%s", err, response)
			return err
		}
	}

	return resourceIBMSchematicsWorkspaceRead(d, meta)
}

func resourceIBMSchematicsWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	schematicsClient, err := meta.(ClientSession).SchematicsV1()
	if err != nil {
		return err
	}

	session, err := meta.(ClientSession).BluemixSession()
	if err != nil {
		return err
	}

	deleteWorkspaceOptions := &schematicsv1.DeleteWorkspaceOptions{}

	deleteWorkspaceOptions.SetWID(d.Id())

	iamRefreshToken := session.Config.IAMRefreshToken
	deleteWorkspaceOptions.SetRefreshToken(iamRefreshToken)

	_, response, err := schematicsClient.DeleteWorkspaceWithContext(context.TODO(), deleteWorkspaceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteWorkspaceWithContext failed %s\n%s", err, response)
		return err
	}

	d.SetId("")

	return nil
}
