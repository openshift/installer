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

func ResourceIBMCdToolchainToolSonarqube() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdToolchainToolSonarqubeCreate,
		ReadContext:   resourceIBMCdToolchainToolSonarqubeRead,
		UpdateContext: resourceIBMCdToolchainToolSonarqubeUpdate,
		DeleteContext: resourceIBMCdToolchainToolSonarqubeDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"toolchain_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_sonarqube", "toolchain_id"),
				Description:  "ID of the toolchain to bind the tool to.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_toolchain_tool_sonarqube", "name"),
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name for this tool integration.",
						},
						"user_login": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user id for authenticating to the SonarQube server.",
						},
						"user_password": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: flex.SuppressHashedRawSecret,
							Sensitive:        true,
							Description:      "The password or token for authenticating to the SonarQube server. You can use a toolchain secret reference for this parameter. For more information, see [Protecting your sensitive data in Continuous Delivery](https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-cd_data_security#cd_secure_credentials).",
						},
						"blind_connection": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "When set to true, instructs IBM Cloud Continuous Delivery to not validate the configuration of this integration. Set this to true if the SonarQube server is not addressable on the public internet.",
						},
						"server_url": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL of the SonarQube server.",
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

func ResourceIBMCdToolchainToolSonarqubeValidator() *validate.ResourceValidator {
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_toolchain_tool_sonarqube", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdToolchainToolSonarqubeCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createToolOptions := &cdtoolchainv2.CreateToolOptions{}

	createToolOptions.SetToolchainID(d.Get("toolchain_id").(string))
	createToolOptions.SetToolTypeID("sonarqube")
	remapFields := map[string]string{
		"server_url": "dashboard_url",
	}
	parametersModel := GetParametersForCreate(d, ResourceIBMCdToolchainToolSonarqube(), remapFields)
	createToolOptions.SetParameters(parametersModel)
	if _, ok := d.GetOk("name"); ok {
		createToolOptions.SetName(d.Get("name").(string))
	}

	toolchainToolPost, _, err := cdToolchainClient.CreateToolWithContext(context, createToolOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateToolWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_sonarqube", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createToolOptions.ToolchainID, *toolchainToolPost.ID))

	return resourceIBMCdToolchainToolSonarqubeRead(context, d, meta)
}

func resourceIBMCdToolchainToolSonarqubeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getToolByIDOptions := &cdtoolchainv2.GetToolByIDOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "sep-id-parts").GetDiag()
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetToolByIDWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_sonarqube", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("toolchain_id", toolchainTool.ToolchainID); err != nil {
		err = fmt.Errorf("Error setting toolchain_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-toolchain_id").GetDiag()
	}
	if !core.IsNil(toolchainTool.Name) {
		if err = d.Set("name", toolchainTool.Name); err != nil {
			err = fmt.Errorf("Error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-name").GetDiag()
		}
	}
	remapFields := map[string]string{
		"server_url": "dashboard_url",
	}
	parametersMap := GetParametersFromRead(toolchainTool.Parameters, ResourceIBMCdToolchainToolSonarqube(), remapFields)
	if err = d.Set("parameters", []map[string]interface{}{parametersMap}); err != nil {
		err = fmt.Errorf("Error setting parameters: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-parameters").GetDiag()
	}
	if err = d.Set("resource_group_id", toolchainTool.ResourceGroupID); err != nil {
		err = fmt.Errorf("Error setting resource_group_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-resource_group_id").GetDiag()
	}
	if err = d.Set("crn", toolchainTool.CRN); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-crn").GetDiag()
	}
	if err = d.Set("toolchain_crn", toolchainTool.ToolchainCRN); err != nil {
		err = fmt.Errorf("Error setting toolchain_crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-toolchain_crn").GetDiag()
	}
	if err = d.Set("href", toolchainTool.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-href").GetDiag()
	}
	referentMap, err := ResourceIBMCdToolchainToolSonarqubeToolModelReferentToMap(toolchainTool.Referent)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "referent-to-map").GetDiag()
	}
	if err = d.Set("referent", []map[string]interface{}{referentMap}); err != nil {
		err = fmt.Errorf("Error setting referent: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-referent").GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(toolchainTool.UpdatedAt)); err != nil {
		err = fmt.Errorf("Error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("state", toolchainTool.State); err != nil {
		err = fmt.Errorf("Error setting state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-state").GetDiag()
	}
	if err = d.Set("tool_id", toolchainTool.ID); err != nil {
		err = fmt.Errorf("Error setting tool_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "read", "set-tool_id").GetDiag()
	}

	return nil
}

func resourceIBMCdToolchainToolSonarqubeUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateToolOptions := &cdtoolchainv2.UpdateToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "update", "sep-id-parts").GetDiag()
	}

	updateToolOptions.SetToolchainID(parts[0])
	updateToolOptions.SetToolID(parts[1])

	hasChange := false

	patchVals := &cdtoolchainv2.ToolchainToolPrototypePatch{}
	if d.HasChange("toolchain_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "toolchain_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_toolchain_tool_sonarqube", "update", "toolchain_id-forces-new").GetDiag()
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("parameters") {
		remapFields := map[string]string{
			"server_url": "dashboard_url",
		}
		parameters := GetParametersForUpdate(d, ResourceIBMCdToolchainToolSonarqube(), remapFields)
		patchVals.Parameters = parameters
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateToolOptions.ToolchainToolPrototypePatch = ResourceIBMCdToolchainToolSonarqubeToolchainToolPrototypePatchAsPatch(patchVals, d)

		_, _, err = cdToolchainClient.UpdateToolWithContext(context, updateToolOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateToolWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_sonarqube", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdToolchainToolSonarqubeRead(context, d, meta)
}

func resourceIBMCdToolchainToolSonarqubeDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdToolchainClient, err := meta.(conns.ClientSession).CdToolchainV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteToolOptions := &cdtoolchainv2.DeleteToolOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_toolchain_tool_sonarqube", "delete", "sep-id-parts").GetDiag()
	}

	deleteToolOptions.SetToolchainID(parts[0])
	deleteToolOptions.SetToolID(parts[1])

	_, err = cdToolchainClient.DeleteToolWithContext(context, deleteToolOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteToolWithContext failed: %s", err.Error()), "ibm_cd_toolchain_tool_sonarqube", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMCdToolchainToolSonarqubeToolModelReferentToMap(model *cdtoolchainv2.ToolModelReferent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.UIHref != nil {
		modelMap["ui_href"] = *model.UIHref
	}
	if model.APIHref != nil {
		modelMap["api_href"] = *model.APIHref
	}
	return modelMap, nil
}

func ResourceIBMCdToolchainToolSonarqubeToolchainToolPrototypePatchAsPatch(patchVals *cdtoolchainv2.ToolchainToolPrototypePatch, d *schema.ResourceData) map[string]interface{} {
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
		ResourceIBMCdToolchainToolSonarqubeToolModelParametersAsPatch(patch["parameters"].(map[string]interface{}), d)
	}

	return patch
}

func ResourceIBMCdToolchainToolSonarqubeToolModelParametersAsPatch(patch map[string]interface{}, d *schema.ResourceData) {
	var path string

	path = "parameters.0.user_login"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["user_login"] = nil
	}
	path = "parameters.0.user_password"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["user_password"] = nil
	}
	path = "parameters.0.blind_connection"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["blind_connection"] = nil
	}
}
