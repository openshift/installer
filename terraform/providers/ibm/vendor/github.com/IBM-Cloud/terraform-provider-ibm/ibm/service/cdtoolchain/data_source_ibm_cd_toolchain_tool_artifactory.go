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

func DataSourceIBMCdToolchainToolArtifactory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdToolchainToolArtifactoryRead,

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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this tool integration.",
						},
						"dashboard_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the Artifactory server dashboard for this integration. In the graphical UI, this is the dashboard that the browser will navigate to when you click the Artifactory integration tile.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of repository for your Artifactory integration.",
						},
						"user_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The User ID or email for your Artifactory repository.",
						},
						"token": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
							Description: "The Access token for your Artifactory repository. You can use a toolchain secret reference for this parameter. For more information, see [Protecting your sensitive data in Continuous Delivery](https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-cd_data_security#cd_secure_credentials).",
						},
						"release_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for your Artifactory release repository.",
						},
						"mirror_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for your Artifactory virtual repository, which is a repository that can see your private repositories and a cache of the public repositories.",
						},
						"snapshot_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for your Artifactory snapshot repository.",
						},
						"repository_name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of your Artifactory repository where your docker images are located.",
						},
						"repository_url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of your Artifactory repository where your docker images are located.",
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

func dataSourceIBMCdToolchainToolArtifactoryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	getToolByIDOptions.SetToolchainID(d.Get("toolchain_id").(string))
	getToolByIDOptions.SetToolID(d.Get("tool_id").(string))

	toolchainTool, _, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetToolByIDWithContext failed: %s", err.Error()), "(Data) ibm_cd_toolchain_tool_artifactory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if *toolchainTool.ToolTypeID != "artifactory" {
		return flex.TerraformErrorf(err, fmt.Sprintf("Retrieved tool is not the correct type: %s", err), "(Data) ibm_cd_toolchain_tool", "read").GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getToolByIDOptions.ToolchainID, *getToolByIDOptions.ToolID))

	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_group_id: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-resource_group_id").GetDiag()
	}

	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-crn").GetDiag()
	}

	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting toolchain_crn: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-toolchain_crn").GetDiag()
	}

	if err = d.Set("href", toolchainTool.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-href").GetDiag()
	}

	referent := []map[string]interface{}{}
	referentMap, err := DataSourceIBMCdToolchainToolArtifactoryToolModelReferentToMap(toolchainTool.Referent)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "referent-to-map").GetDiag()
	}
	referent = append(referent, referentMap)
	if err = d.Set("referent", referent); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting referent: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-referent").GetDiag()
	}

	if !core.IsNil(toolchainTool.Name) {
		if err = d.Set("name", toolchainTool.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-name").GetDiag()
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-updated_at").GetDiag()
	}

	parameters := []map[string]interface{}{}
	parametersMap := GetParametersFromRead(toolchainTool.Parameters, DataSourceIBMCdToolchainToolArtifactory(), nil)
	parameters = append(parameters, parametersMap)
	if err = d.Set("parameters", parameters); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting parameters: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-parameters").GetDiag()
	}

	if err = d.Set("state", toolchainTool.State); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting state: %s", err), "(Data) ibm_cd_toolchain_tool_artifactory", "read", "set-state").GetDiag()
	}

	return nil
}

func DataSourceIBMCdToolchainToolArtifactoryToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = *model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = *model.APIHref
	}
	return modelMap, nil
}
