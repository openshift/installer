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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtoolchainv2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func DataSourceIBMCdToolchainToolHostedgit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdToolchainToolHostedgitRead,

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the toolchain.",
			},
			"tool_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the tool bound to the toolchain.",
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
							Computed:    true,
							Description: "URI representing this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "URI representing this resource through an API.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the tool.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Latest tool update timestamp.",
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool integration can be found in the <a href=\"https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations\">Configuring tool integrations page</a>.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Set this value to 'hostedgit' to target Git Repos and Issue Tracking.",
						},
						"api_root_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API root URL for the GitLab server.",
						},
						"default_branch": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default branch of the git repository.",
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
							Description: "The token URL used for authorizing with the Bitbucket server.",
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
							Computed:    true,
							Description: "Set this value to 'true' to track the deployment of code changes by creating tags, labels and comments on commits, pull requests and referenced issues.",
						},
						"integration_owner": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Select the user which git operations will be performed as.",
						},
						"repo_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the Git Repos and Issue Tracking project.",
						},
						"auth_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Select the method of authentication that will be used to access the git provider. The default value is 'oauth'.",
						},
						"api_token": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "Personal Access Token. Required if 'auth_type' is set to 'pat', ignored otherwise.",
						},
						"toolchain_issues_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Setting this value to true will enable issues on the GitLab repository and add an issues tool card to the toolchain.  Setting the value to false will remove the tool card from the toolchain, but will not impact whether or not issues are enabled on the GitLab repository itself.",
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current configuration state of the tool.",
			},
		},
	}
}

func dataSourceIBMCdToolchainToolHostedgitRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	getToolByIDOptions.SetToolchainID(d.Get("toolchain_id").(string))
	getToolByIDOptions.SetToolID(d.Get("tool_id").(string))

	toolchainTool, _, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetToolByIDWithContext failed: %s", err.Error()), "(Data) ibm_cd_toolchain_tool_hostedgit", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if *toolchainTool.ToolTypeID != "hostedgit" {
		return flex.TerraformErrorf(err, fmt.Sprintf("Retrieved tool is not the correct type: %s", err), "(Data) ibm_cd_toolchain_tool", "read").GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getToolByIDOptions.ToolchainID, *getToolByIDOptions.ToolID))

	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-resource_group_id").GetDiag()
	}

	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-crn").GetDiag()
	}

	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting toolchain_crn: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-toolchain_crn").GetDiag()
	}

	if err = d.Set("href", toolchainTool.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-href").GetDiag()
	}

	referent := []map[string]interface{}{}
	referentMap, err := DataSourceIBMCdToolchainToolHostedgitToolModelReferentToMap(toolchainTool.Referent)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "referent-to-map").GetDiag()
	}
	referent = append(referent, referentMap)
	if err = d.Set("referent", referent); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting referent: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-referent").GetDiag()
	}

	if !core.IsNil(toolchainTool.Name) {
		if err = d.Set("name", toolchainTool.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-name").GetDiag()
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-updated_at").GetDiag()
	}

	parameters := []map[string]interface{}{}
	remapFields := map[string]string{
		"toolchain_issues_enabled": "has_issues",
	}
	parametersMap := GetParametersFromRead(toolchainTool.Parameters, DataSourceIBMCdToolchainToolHostedgit(), remapFields)
	parameters = append(parameters, parametersMap)
	if err = d.Set("parameters", parameters); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting parameters: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-parameters").GetDiag()
	}

	if err = d.Set("state", toolchainTool.State); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "(Data) ibm_cd_toolchain_tool_hostedgit", "read", "set-state").GetDiag()
	}

	return nil
}

func DataSourceIBMCdToolchainToolHostedgitToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = *model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = *model.APIHref
	}
	return modelMap, nil
}
