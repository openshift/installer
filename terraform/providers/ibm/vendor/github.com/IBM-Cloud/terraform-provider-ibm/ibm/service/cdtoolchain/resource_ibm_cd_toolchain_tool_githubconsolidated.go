// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

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
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtoolchainv2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMCdToolchainToolGithubconsolidated() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdToolchainToolGithubconsolidatedCreate,
		ReadContext:   resourceIBMCdToolchainToolGithubconsolidatedRead,
		UpdateContext: resourceIBMCdToolchainToolGithubconsolidatedUpdate,
		DeleteContext: resourceIBMCdToolchainToolGithubconsolidatedDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_githubconsolidated", "toolchain_id"),
				Description:  "ID of the toolchain to bind the tool to.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_githubconsolidated", "name"),
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
							Optional:    true,
							Computed:    true,
							Description: "Set this value to 'github' for github.com, or 'githubcustom' for a custom GitHub Enterprise server.",
						},
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The title of the server. e.g. My GitHub Enterprise Server.",
						},
						"api_root_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The API root URL for the GitHub server.",
						},
						"default_branch": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The default branch of the git repository.",
						},
						"root_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The Root URL of the server. e.g. https://github.example.com.",
						},
						"blind_connection": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Setting this value to true means the server is not addressable on the public internet. IBM Cloud will not be able to validate the connection details you provide. Certain functionality that requires API access to the git server will be disabled. Delivery pipeline will only work using a private worker that has network access to the git server.",
						},
						"owner_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The GitHub user or organization that owns the repository.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name of the new GitHub repository to create.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL of the GitHub repository for this tool integration.  This parameter is required when linking to an existing repository.  The value will be computed when creating a new repository, cloning, or forking a repository.",
						},
						"source_repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The URL of the repository that you are forking or cloning.  This parameter is required when forking or cloning a repository.  It is not used when creating a new repository or linking to an existing repository.",
						},
						"token_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The token URL used for authorizing with the GitHub server.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The operation that should be performed to initialize the new tool integration. Use 'new' or 'new_if_not_exists' to create a new git repository, 'clone' or 'clone_if_not_exists' to clone an existing repository into a new git repository, 'fork' or 'fork_if_not_exists' to fork an existing git repository, or 'link' to link to an existing git repository. If you attempt to apply a resource with type 'new', 'clone', or 'fork' when the target repo already exists, the attempt will fail. If you apply a resource with type 'new_if_not_exists`, 'clone_if_not_exists', or 'fork_if_not_exists' when the target repo already exists, the existing repo will be used as-is.",
						},
						"private_repo": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Set this value to 'true' to make the repository private when creating a new repository or when cloning or forking a repository.  This parameter is not used when linking to an existing repository.",
						},
						"auto_init": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Setting this value to true will initialize this repository with a README.  This parameter is only used when creating a new repository.",
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
							Optional:    true,
							Computed:    true,
							Description: "The ID of the GitHub repository.",
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
							Description: "Setting this value to true will enable issues on the GitHub repository and add an issues tool card to the toolchain.  Setting the value to false will remove the tool card from the toolchain, but will not impact whether or not issues are enabled on the GitHub repository itself.",
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
							Description: "Set this value to 'github' for github.com, or 'githubcustom' for a custom GitHub Enterprise server.",
						},
						"title": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The title of the server. e.g. My GitHub Enterprise Server.",
						},
						"root_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The Root URL of the server. e.g. https://github.example.com.",
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
							Description: "The GitHub user or organization that owns the repository.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The name of the new GitHub repository to create.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The URL of the GitHub repository for this tool integration.  This parameter is required when linking to an existing repository.  The value will be computed when creating a new repository, cloning, or forking a repository.",
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
							Default:     false,
							ForceNew:    true,
							Description: "Set this value to 'true' to make the repository private when creating a new repository or when cloning or forking a repository.  This parameter is not used when linking to an existing repository.",
						},
						"auto_init": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							ForceNew:    true,
							Description: "Setting this value to true will initialize this repository with a README.  This parameter is only used when creating a new repository.",
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
							Computed:    true,
							Description: "URI representing this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
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

func ResourceIBMCdToolchainToolGithubconsolidatedValidator() *validate.ResourceValidator {
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain_tool_githubconsolidated", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdToolchainToolGithubconsolidatedCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createToolOptions := &cdtoolchainv2.CreateToolOptions{}

	createToolOptions.SetToolchainID(d.Get("toolchain_id").(string))
	createToolOptions.SetToolTypeID("githubconsolidated")
	remapFields := map[string]string{
		"toolchain_issues_enabled": "has_issues",
	}
	parametersModel := GetParametersForCreate(d, ResourceIBMCdToolchainToolGithubconsolidated(), remapFields)
	createToolOptions.SetParameters(parametersModel)
	if _, ok := d.GetOk("name"); ok {
		createToolOptions.SetName(d.Get("name").(string))
	}

	toolchainToolPost, _, err := cdToolchainClient.CreateToolWithContext(context, createToolOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateToolWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_githubconsolidated", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createToolOptions.ToolchainID, *toolchainToolPost.ID))

	return resourceIBMCdToolchainToolGithubconsolidatedRead(context, d, meta)
}

func resourceIBMCdToolchainToolGithubconsolidatedRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "sep-id-parts").GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetToolByIDWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_githubconsolidated", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("toolchain_id", toolchainTool.ToolchainID); err != nil {
		err = fmt.Errorf("Error setting toolchain_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-toolchain_id").GetDiag()
	}
	if !core.IsNil(toolchainTool.Name) {
		if err = d.Set("name", toolchainTool.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-name").GetDiag()
		}
	}
	remapFields := map[string]string{
		"toolchain_issues_enabled": "has_issues",
	}
	parametersMap := GetParametersFromRead(toolchainTool.Parameters, ResourceIBMCdToolchainToolGithubconsolidated(), remapFields)
	if err = d.Set("parameters", []map[string]interface{}{parametersMap}); err != nil {
		err = fmt.Errorf("Error setting parameters: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-parameters").GetDiag()
	}
	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		err = fmt.Errorf("Error setting resource_group_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-resource_group_id").GetDiag()
	}
	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-crn").GetDiag()
	}
	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		err = fmt.Errorf("Error setting toolchain_crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-toolchain_crn").GetDiag()
	}
	if err = d.Set("href", toolchainTool.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-href").GetDiag()
	}
	referentMap, err := ResourceIBMCdToolchainToolGithubconsolidatedToolModelReferentToMap(toolchainTool.Referent)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "referent-to-map").GetDiag()
	}
	if err = d.Set("referent", []map[string]interface{}{referentMap}); err != nil {
		err = fmt.Errorf("Error setting referent: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-referent").GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		err = fmt.Errorf("Error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("state", toolchainTool.State); err != nil {
		err = fmt.Errorf("Error setting state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-state").GetDiag()
	}
	if err = d.Set("tool_id", toolchainTool.ID); err != nil {
		err = fmt.Errorf("Error setting tool_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "read", "set-tool_id").GetDiag()
	}

	return nil
}

func resourceIBMCdToolchainToolGithubconsolidatedUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateToolOptions := &cdtoolchainv2.UpdateToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "update", "sep-id-parts").GetDiag()
	}

	updateToolOptions.SetToolchainID(parts[0])
	updateToolOptions.SetToolID(parts[1])

	hasChange := false

	patchVals := &cdtoolchainv2.ToolchainToolPrototypePatch{}
	if d.HasChange("toolchain_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "toolchain_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_toolchain_tool_githubconsolidated", "update", "toolchain_id-forces-new").GetDiag()
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
		parameters := GetParametersForUpdate(d, ResourceIBMCdToolchainToolGithubconsolidated(), remapFields)
		patchVals.Parameters = parameters
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateToolOptions.ToolchainToolPrototypePatch = ResourceIBMCdToolchainToolGithubconsolidatedToolchainToolPrototypePatchAsPatch(patchVals, d)

		_, _, err = cdToolchainClient.UpdateToolWithContext(context, updateToolOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateToolWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_githubconsolidated", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdToolchainToolGithubconsolidatedRead(context, d, meta)
}

func resourceIBMCdToolchainToolGithubconsolidatedDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteToolOptions := &cdtoolchainv2.DeleteToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_githubconsolidated", "delete", "sep-id-parts").GetDiag()
	}

	deleteToolOptions.SetToolchainID(parts[0])
	deleteToolOptions.SetToolID(parts[1])

	_, err = cdToolchainClient.DeleteToolWithContext(context, deleteToolOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteToolWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_githubconsolidated", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMCdToolchainToolGithubconsolidatedToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = *model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = *model.APIHref
	}
	return modelMap, nil
}

func ResourceIBMCdToolchainToolGithubconsolidatedToolchainToolPrototypePatchAsPatch(patchVals *cdtoolchainv2.ToolchainToolPrototypePatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	}
	path = "tool_type_id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["tool_type_id"] = nil
	}
	path = "parameters"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["parameters"] = nil
	} else if exists && patch["parameters"] != nil {
		ResourceIBMCdToolchainToolGithubconsolidatedToolModelParametersAsPatch(patch["parameters"].(map[string]interface{}), d)
	}

	return patch
}

func ResourceIBMCdToolchainToolGithubconsolidatedToolModelParametersAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "parameters.0.git_id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["git_id"] = nil
	}
	path = "parameters.0.title"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["title"] = nil
	}
	path = "parameters.0.api_root_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["api_root_url"] = nil
	}
	path = "parameters.0.default_branch"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["default_branch"] = nil
	}
	path = "parameters.0.root_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["root_url"] = nil
	}
	path = "parameters.0.blind_connection"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["blind_connection"] = nil
	}
	path = "parameters.0.owner_id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["owner_id"] = nil
	}
	path = "parameters.0.repo_name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["repo_name"] = nil
	}
	path = "parameters.0.repo_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["repo_url"] = nil
	}
	path = "parameters.0.source_repo_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["source_repo_url"] = nil
	}
	path = "parameters.0.token_url"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["token_url"] = nil
	}
	path = "parameters.0.type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["type"] = nil
	}
	path = "parameters.0.private_repo"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["private_repo"] = nil
	}
	path = "parameters.0.auto_init"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auto_init"] = nil
	}
	path = "parameters.0.enable_traceability"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["enable_traceability"] = nil
	}
	path = "parameters.0.integration_owner"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["integration_owner"] = nil
	}
	path = "parameters.0.repo_id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["repo_id"] = nil
	}
	path = "parameters.0.auth_type"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["auth_type"] = nil
	}
	path = "parameters.0.api_token"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["api_token"] = nil
	}
	path = "parameters.0.toolchain_issues_enabled"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["toolchain_issues_enabled"] = nil
	}
}
