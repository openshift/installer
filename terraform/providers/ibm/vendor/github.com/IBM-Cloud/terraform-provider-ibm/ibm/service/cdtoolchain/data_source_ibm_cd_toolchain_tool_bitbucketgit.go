// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtoolchain

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func DataSourceIBMCdToolchainToolBitbucketgit() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdToolchainToolBitbucketgitRead,

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
							Description: "Set this value to 'bitbucketgit' for bitbucket.org, or to the GUID of a custom Bitbucket server.",
						},
						"api_root_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The API root URL for the Bitbucket Server.",
						},
						"default_branch": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The default branch of the git repository.",
						},
						"owner_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Bitbucket user or group that owns the repository.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the new Bitbucket repository to create.  This parameter is required when creating a new repository, cloning, or forking a repository.  The value will be computed when linking to an existing repository.",
						},
						"repo_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the bitbucket repository for this tool integration.  This parameter is required when linking to an existing repository.  The value will be computed when creating a new repository, cloning, or forking a repository.",
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
							Description: "The ID of the Bitbucket repository.",
						},
						"toolchain_issues_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Setting this value to true will enable issues on the Bitbucket repository and add an issues tool card to the toolchain.  Setting the value to false will remove the tool card from the toolchain, but will not impact whether or not issues are enabled on the Bitbucket repository itself.",
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

func dataSourceIBMCdToolchainToolBitbucketgitRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	getToolByIDOptions.SetToolchainID(d.Get("toolchain_id").(string))
	getToolByIDOptions.SetToolID(d.Get("tool_id").(string))

	toolchainTool, response, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		log.Printf("[DEBUG] GetToolByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolByIDWithContext failed %s\n%s", err, response))
	}

	if *toolchainTool.ToolTypeID != "bitbucketgit" {
		return diag.FromErr(fmt.Errorf("Retrieved tool is not the correct type: %s", *toolchainTool.ToolTypeID))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getToolByIDOptions.ToolchainID, *getToolByIDOptions.ToolID))

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

	referent := []map[string]interface{}{}
	if toolchainTool.Referent != nil {
		modelMap, err := dataSourceIBMCdToolchainToolBitbucketgitToolModelReferentToMap(toolchainTool.Referent)
		if err != nil {
			return diag.FromErr(err)
		}
		referent = append(referent, modelMap)
	}
	if err = d.Set("referent", referent); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent %s", err))
	}

	if err = d.Set("name", toolchainTool.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}

	parameters := []map[string]interface{}{}
	if toolchainTool.Parameters != nil {
		remapFields := map[string]string{
			"toolchain_issues_enabled": "has_issues",
		}
		modelMap := GetParametersFromRead(toolchainTool.Parameters, DataSourceIBMCdToolchainToolBitbucketgit(), remapFields)
		parameters = append(parameters, modelMap)
	}
	if err = d.Set("parameters", parameters); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parameters %s", err))
	}

	if err = d.Set("state", toolchainTool.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}

	return nil
}

func dataSourceIBMCdToolchainToolBitbucketgitToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = model.APIHref
	}
	return modelMap, nil
}
