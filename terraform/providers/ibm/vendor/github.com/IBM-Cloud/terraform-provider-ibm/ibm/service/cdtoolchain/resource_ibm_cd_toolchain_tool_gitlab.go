// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMCdToolchainToolGitlab() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdToolchainToolGitlabCreate,
		ReadContext:   resourceIBMCdToolchainToolGitlabRead,
		UpdateContext: resourceIBMCdToolchainToolGitlabUpdate,
		DeleteContext: resourceIBMCdToolchainToolGitlabDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_gitlab", "toolchain_id"),
				Description:  "ID of the toolchain to bind the tool to.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_gitlab", "name"),
				Description:  "Name of the tool.",
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool integration can be found in the <a href=\"https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations\">Configuring tool integrations page</a>.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Set this value to 'gitlab' for gitlab.com, or 'gitlabcustom' for a custom GitLab server.",
						},
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The title of the server. e.g. My GitLab Enterprise Server.",
						},
						"api_root_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API root URL for the GitLab Server.",
						},
						"default_branch": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default branch of the git repository.",
						},
						"root_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Root URL of the server. e.g. https://gitlab.example.com.",
						},
						"blind_connection": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Setting this value to true means the server is not addressable on the public internet. IBM Cloud will not be able to validate the connection details you provide. Certain functionality that requires API access to the git server will be disabled. Delivery pipeline will only work using a private worker that has network access to the git server.",
						},
						"owner_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The GitLab user or group that owns the repository.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the new GitLab repository to create.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the GitLab repository for this tool integration.  This parameter is required when linking to an existing repository.  The value will be computed when creating a new repository, cloning, or forking a repository.",
						},
						"source_repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the repository that you are forking or cloning.  This parameter is required when forking or cloning a repository.  It is not used when creating a new repository or linking to an existing repository.",
						},
						"token_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The token URL used for authorizing with the GitLab server.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The operation that should be performed to initialize the new tool integration. Use 'new' or 'new_if_not_exists' to create a new git repository, 'clone' or 'clone_if_not_exists' to clone an existing repository into a new git repository, 'fork' or 'fork_if_not_exists' to fork an existing git repository, or 'link' to link to an existing git repository. If you attempt to apply a resource with type 'new', 'clone', or 'fork' when the target repo already exists, the attempt will fail. If you apply a resource with type 'new_if_not_exists`, 'clone_if_not_exists', or 'fork_if_not_exists' when the target repo already exists, the existing repo will be used as-is.",
						},
						"private_repo": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Set this value to 'true' to make the repository private when creating a new repository or when cloning or forking a repository.  This parameter is not used when linking to an existing repository.",
						},
						"enable_traceability": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Set this value to 'true' to track the deployment of code changes by creating tags, labels and comments on commits, pull requests and referenced issues.",
						},
						"integration_owner": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: flex.SuppressAllowBlank,
							Description:      "Select the user which git operations will be performed as.",
						},
						"repo_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the GitLab project.",
						},
						"auth_type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Select the method of authentication that will be used to access the git provider. The default value is 'oauth'.",
						},
						"api_token": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: flex.SuppressHashedRawSecret,
							Sensitive:        true,
							Description:      "Personal Access Token. Required if ‘auth_type’ is set to ‘pat’, ignored otherwise.",
						},
						"toolchain_issues_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							Description: "Setting this value to true will enable issues on the GitLab repository and add an issues tool card to the toolchain.  Setting the value to false will remove the tool card from the toolchain, but will not impact whether or not issues are enabled on the GitLab repository itself.",
						},
					},
				},
			},
			"initialization": &schema.Schema{
				Type:     schema.TypeList,
				MinItems: 1,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "Set this value to 'gitlab' for gitlab.com, or 'gitlabcustom' for a custom GitLab server.",
						},
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The title of the server. e.g. My GitLab Enterprise Server.",
						},
						"root_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The Root URL of the server. e.g. https://gitlab.example.com.",
						},
						"blind_connection": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Setting this value to true means the server is not addressable on the public internet. IBM Cloud will not be able to validate the connection details you provide. Certain functionality that requires API access to the git server will be disabled. Delivery pipeline will only work using a private worker that has network access to the git server.",
						},
						"owner_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The GitLab user or group that owns the repository.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The name of the new GitLab repository to create.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The URL of the GitLab repository for this tool integration.  This parameter is required when linking to an existing repository.  The value will be computed when creating a new repository, cloning, or forking a repository.",
						},
						"source_repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The URL of the repository that you are forking or cloning.  This parameter is required when forking or cloning a repository.  It is not used when creating a new repository or linking to an existing repository.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "The operation that should be performed to initialize the new tool integration. Use 'new' or 'new_if_not_exists' to create a new git repository, 'clone' or 'clone_if_not_exists' to clone an existing repository into a new git repository, 'fork' or 'fork_if_not_exists' to fork an existing git repository, or 'link' to link to an existing git repository. If you attempt to apply a resource with type 'new', 'clone', or 'fork' when the target repo already exists, the attempt will fail. If you apply a resource with type 'new_if_not_exists`, 'clone_if_not_exists', or 'fork_if_not_exists' when the target repo already exists, the existing repo will be used as-is.",
						},
						"private_repo": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Set this value to 'true' to make the repository private when creating a new repository or when cloning or forking a repository.  This parameter is not used when linking to an existing repository.",
						},
					},
				},
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group where the tool is located.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool CRN.",
			},
			"toolchain_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CRN of toolchain which the tool is bound to.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URI representing the tool.",
			},
			"referent": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information on URIs to access this resource through the UI or API.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ui_href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URI representing this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URI representing this resource through an API.",
						},
					},
				},
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest tool update timestamp.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current configuration state of the tool.",
			},
			"tool_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Tool ID.",
			},
		},
	}
}

func ResourceIBMCdToolchainToolGitlabValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "toolchain_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([^\x00-\x7F]|[a-zA-Z0-9-._ ])+$`,
			MinValueLength:             0,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain_tool_gitlab", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdToolchainToolGitlabCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createToolOptions := &cdtoolchainv2.CreateToolOptions{}

	createToolOptions.SetToolchainID(d.Get("toolchain_id").(string))
	createToolOptions.SetToolTypeID("gitlab")
	remapFields := map[string]string{
		"toolchain_issues_enabled": "has_issues",
	}
	parametersModel := GetParametersForCreate(d, ResourceIBMCdToolchainToolGitlab(), remapFields)
	createToolOptions.SetParameters(parametersModel)
	if _, ok := d.GetOk("name"); ok {
		createToolOptions.SetName(d.Get("name").(string))
	}

	toolchainToolPost, response, err := cdToolchainClient.CreateToolWithContext(context, createToolOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateToolWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateToolWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createToolOptions.ToolchainID, *toolchainToolPost.ID))

	return resourceIBMCdToolchainToolGitlabRead(context, d, meta)
}

func resourceIBMCdToolchainToolGitlabRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions.SetToolchainID(parts[0])
	getToolByIDOptions.SetToolID(parts[1])

	var toolchainTool *cdtoolchainv2.ToolchainTool
	var response *core.DetailedResponse
	err = resource.RetryContext(context, 10*time.Second, func() *resource.RetryError {
		toolchainTool, response, err = cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
		if err != nil || toolchainTool == nil {
			if response != nil && response.StatusCode == 404 {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if conns.IsResourceTimeoutError(err) {
		toolchainTool, response, err = cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	}
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetToolByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolByIDWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("toolchain_id", toolchainTool.ToolchainID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_id: %s", err))
	}
	if !core.IsNil(toolchainTool.Name) {
		if err = d.Set("name", toolchainTool.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	remapFields := map[string]string{
		"toolchain_issues_enabled": "has_issues",
	}
	parametersMap := GetParametersFromRead(toolchainTool.Parameters, ResourceIBMCdToolchainToolGitlab(), remapFields)
	if err = d.Set("parameters", []map[string]interface{}{parametersMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parameters: %s", err))
	}
	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_crn: %s", err))
	}
	if err = d.Set("href", toolchainTool.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	referentMap, err := resourceIBMCdToolchainToolGitlabToolModelReferentToMap(toolchainTool.Referent)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("referent", []map[string]interface{}{referentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("state", toolchainTool.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("tool_id", toolchainTool.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tool_id: %s", err))
	}

	return nil
}

func resourceIBMCdToolchainToolGitlabUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateToolOptions := &cdtoolchainv2.UpdateToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateToolOptions.SetToolchainID(parts[0])
	updateToolOptions.SetToolID(parts[1])

	hasChange := false

	patchVals := &cdtoolchainv2.ToolchainToolPrototypePatch{}
	if d.HasChange("toolchain_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "toolchain_id"))
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("parameters") {
		remapFields := map[string]string{
			"toolchain_issues_enabled": "has_issues",
		}
		parameters := GetParametersForUpdate(d, ResourceIBMCdToolchainToolGitlab(), remapFields)
		patchVals.Parameters = parameters
		hasChange = true
	}

	if hasChange {
		updateToolOptions.ToolchainToolPrototypePatch, _ = patchVals.AsPatch()
		_, response, err := cdToolchainClient.UpdateToolWithContext(context, updateToolOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateToolWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateToolWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdToolchainToolGitlabRead(context, d, meta)
}

func resourceIBMCdToolchainToolGitlabDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolOptions := &cdtoolchainv2.DeleteToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteToolOptions.SetToolchainID(parts[0])
	deleteToolOptions.SetToolID(parts[1])

	response, err := cdToolchainClient.DeleteToolWithContext(context, deleteToolOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteToolWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteToolWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCdToolchainToolGitlabToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = model.APIHref
	}
	return modelMap, nil
}
