// Copyright IBM Corp. 2022 All Rights Reserved.
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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtoolchainv2"
)

func ResourceIBMCdToolchainToolArtifactory() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMCdToolchainToolArtifactoryCreate,
		ReadContext:   ResourceIBMCdToolchainToolArtifactoryRead,
		UpdateContext: ResourceIBMCdToolchainToolArtifactoryUpdate,
		DeleteContext: ResourceIBMCdToolchainToolArtifactoryDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_artifactory", "toolchain_id"),
				Description:  "ID of the toolchain to bind tool to.",
			},
			"parameters": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Parameters to be used to create the tool.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type a name for this tool integration, for example: my-artifactory. This name displays on your toolchain.",
						},
						"dashboard_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL that you want to navigate to when you click the Artifactory integration tile.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "Choose the type of repository for your Artifactory integration.",
						},
						"user_id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the User ID or email for your Artifactory repository.",
						},
						"token": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: flex.SuppressHashedRawSecret,
							Sensitive:        true,
							Description:      "Type the API key for your Artifactory repository.",
						},
						"release_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL for your Artifactory release repository.",
						},
						"mirror_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL for your Artifactory virtual repository, which is a repository that can see your private repositories and a cache of the public repositories.",
						},
						"snapshot_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL for your Artifactory snapshot repository.",
						},
						"repository_name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the name of your artifactory repository where your docker images are located.",
						},
						"repository_url": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Type the URL of your artifactory repository where your docker images are located.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_artifactory", "name"),
				Description:  "Name of tool.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Resource group where tool can be found.",
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
							Description: "URI representing the this resource through the UI.",
						},
						"api_href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "URI representing the this resource through an API.",
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

func ResourceIBMCdToolchainToolArtifactoryValidator() *validate.ResourceValidator {
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
			Regexp:                     `^([^\\x00-\\x7F]|[a-zA-Z0-9-._ ])+$`,
			MinValueLength:             0,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain_tool_artifactory", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMCdToolchainToolArtifactoryCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createToolOptions := &cdtoolchainv2.CreateToolOptions{}

	createToolOptions.SetToolchainID(d.Get("toolchain_id").(string))
	createToolOptions.SetToolTypeID("artifactory")
	parametersModel := GetParametersForCreate(d, ResourceIBMCdToolchainToolArtifactory(), nil)
	createToolOptions.SetParameters(parametersModel)
	if _, ok := d.GetOk("name"); ok {
		createToolOptions.SetName(d.Get("name").(string))
	}

	postToolResponse, response, err := cdToolchainClient.CreateToolWithContext(context, createToolOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateToolWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateToolWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createToolOptions.ToolchainID, *postToolResponse.ID))

	return ResourceIBMCdToolchainToolArtifactoryRead(context, d, meta)
}

func ResourceIBMCdToolchainToolArtifactoryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	getToolByIDResponse, response, err := cdToolchainClient.GetToolByIDWithContext(context, getToolByIDOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetToolByIDWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetToolByIDWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("toolchain_id", getToolByIDResponse.ToolchainID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_id: %s", err))
	}
	parametersMap := GetParametersFromRead(getToolByIDResponse.Parameters, ResourceIBMCdToolchainToolArtifactory(), nil)
	if err = d.Set("parameters", []map[string]interface{}{parametersMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting parameters: %s", err))
	}
	if err = d.Set("name", getToolByIDResponse.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("resource_group_id", getToolByIDResponse.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("crn", getToolByIDResponse.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("toolchain_crn", getToolByIDResponse.ToolchainCRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting toolchain_crn: %s", err))
	}
	if err = d.Set("href", getToolByIDResponse.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	referentMap, err := ResourceIBMCdToolchainToolArtifactoryToolReferentToMap(getToolByIDResponse.Referent)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("referent", []map[string]interface{}{referentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting referent: %s", err))
	}
	if err = d.Set("updated_at", flex.DateTimeToString(getToolByIDResponse.UpdatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting updated_at: %s", err))
	}
	if err = d.Set("state", getToolByIDResponse.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("tool_id", getToolByIDResponse.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tool_id: %s", err))
	}

	return nil
}

func ResourceIBMCdToolchainToolArtifactoryUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	updateToolOptions.SetToolTypeID("artifactory")

	hasChange := false

	if d.HasChange("toolchain_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "toolchain_id"))
	}
	if d.HasChange("parameters") {
		parameters := GetParametersForUpdate(d, ResourceIBMCdToolchainToolArtifactory(), nil)
		updateToolOptions.SetParameters(parameters)
		hasChange = true
	}
	if d.HasChange("name") {
		updateToolOptions.SetName(d.Get("name").(string))
		hasChange = true
	}

	if hasChange {
		response, err := cdToolchainClient.UpdateToolWithContext(context, updateToolOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateToolWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateToolWithContext failed %s\n%s", err, response))
		}
	}

	return ResourceIBMCdToolchainToolArtifactoryRead(context, d, meta)
}

func ResourceIBMCdToolchainToolArtifactoryDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func ResourceIBMCdToolchainToolArtifactoryToolReferentToMap(model *cdtoolchainv2.ToolReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = model.APIHref
	}
	return modelMap, nil
}
