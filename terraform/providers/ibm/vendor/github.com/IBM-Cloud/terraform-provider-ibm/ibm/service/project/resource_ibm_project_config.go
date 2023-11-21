// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package project

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/project-go-sdk/projectv1"
)

func ResourceIbmProjectConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmProjectConfigCreate,
		ReadContext:   resourceIbmProjectConfigRead,
		UpdateContext: resourceIbmProjectConfigUpdate,
		DeleteContext: resourceIbmProjectConfigDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_project_config", "project_id"),
				Description:  "The unique project ID.",
			},
			"schematics": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "A schematics workspace associated to a project configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"workspace_crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "An existing schematics workspace CRN.",
						},
					},
				},
			},
			"definition": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The type and output of a project configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The configuration name.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A project configuration description.",
						},
						"environment": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the project environment.",
						},
						"authorizations": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The authorization details. You can authorize by using a trusted profile or an API key in Secrets Manager.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"trusted_profile_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The trusted profile ID.",
									},
									"method": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The authorization method. You can authorize by using a trusted profile or an API key in Secrets Manager.",
									},
									"api_key": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The IBM Cloud API Key.",
									},
								},
							},
						},
						"compliance_profile": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The profile required for compliance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique ID.",
									},
									"instance_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique ID.",
									},
									"instance_location": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The location of the compliance instance.",
									},
									"attachment_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique ID.",
									},
									"profile_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the compliance profile.",
									},
								},
							},
						},
						"locator_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: "A dotted value of catalogID.versionID.",
						},
						"inputs": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The input variables for configuration definition and environment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"settings": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Schematics environment variables to use to deploy the configuration.Settings are only available if they were specified when the configuration was initially created.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of a project configuration manual property.",
						},
					},
				},
			},
			"version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the configuration.",
			},
			"is_draft": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag that indicates whether the version of the configuration is draft, or active.",
			},
			"needs_attention_state": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The needs attention state of a configuration.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time format as specified by RFC 3339.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time format as specified by RFC 3339.",
			},
			"last_saved_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time format as specified by RFC 3339.",
			},
			"outputs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The outputs of a Schematics template property.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The variable name.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A short explanation of the output value.",
						},
						"value": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Can be any value - a string, number, boolean, array, or object.",
						},
					},
				},
			},
			"project": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project referenced by this resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique ID.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The definition of the project reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the project.",
									},
								},
							},
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "An IBM Cloud resource name, which uniquely identifies a resource.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "A URL.",
						},
					},
				},
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state of the configuration.",
			},
			"update_available": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The flag that indicates whether a configuration update is available.",
			},
			"project_config_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.",
			},
		},
	}
}

func ResourceIbmProjectConfigValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "project_id",
			ValidateFunctionIdentifier: validate.ValidateRegexp,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[\.\-0-9a-zA-Z]+$`,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_project_config", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmProjectConfigCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createConfigOptions := &projectv1.CreateConfigOptions{}

	createConfigOptions.SetProjectID(d.Get("project_id").(string))
	definitionModel, err := resourceIbmProjectConfigMapToProjectConfigPrototypeDefinitionBlock(d.Get("definition.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createConfigOptions.SetDefinition(definitionModel)
	if _, ok := d.GetOk("schematics"); ok {
		schematicsModel, err := resourceIbmProjectConfigMapToSchematicsWorkspace(d.Get("schematics.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createConfigOptions.SetSchematics(schematicsModel)
	}

	projectConfig, response, err := projectClient.CreateConfigWithContext(context, createConfigOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateConfigWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateConfigWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createConfigOptions.ProjectID, *projectConfig.ID))

	return resourceIbmProjectConfigRead(context, d, meta)
}

func resourceIbmProjectConfigRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigOptions := &projectv1.GetConfigOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getConfigOptions.SetProjectID(parts[0])
	getConfigOptions.SetID(parts[1])

	projectConfig, response, err := projectClient.GetConfigWithContext(context, getConfigOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetConfigWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigWithContext failed %s\n%s", err, response))
	}

	definitionMap, err := resourceIbmProjectConfigProjectConfigResponseDefinitionToMap(projectConfig.Definition)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("definition", []map[string]interface{}{definitionMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definition: %s", err))
	}
	if err = d.Set("version", flex.IntValue(projectConfig.Version)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting version: %s", err))
	}
	if err = d.Set("is_draft", projectConfig.IsDraft); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting is_draft: %s", err))
	}
	if !core.IsNil(projectConfig.NeedsAttentionState) {
		needsAttentionState := []interface{}{}
		for _, needsAttentionStateItem := range projectConfig.NeedsAttentionState {
			needsAttentionState = append(needsAttentionState, needsAttentionStateItem)
		}
		if err = d.Set("needs_attention_state", needsAttentionState); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting needs_attention_state: %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(projectConfig.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("modified_at", flex.DateTimeToString(projectConfig.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}
	if !core.IsNil(projectConfig.LastSavedAt) {
		if err = d.Set("last_saved_at", flex.DateTimeToString(projectConfig.LastSavedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting last_saved_at: %s", err))
		}
	}
	if !core.IsNil(projectConfig.Outputs) {
		outputs := []map[string]interface{}{}
		for _, outputsItem := range projectConfig.Outputs {
			outputsItemMap, err := resourceIbmProjectConfigOutputValueToMap(&outputsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			outputs = append(outputs, outputsItemMap)
		}
		if err = d.Set("outputs", outputs); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting outputs: %s", err))
		}
	}
	projectMap, err := resourceIbmProjectConfigProjectReferenceToMap(projectConfig.Project)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("project", []map[string]interface{}{projectMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project: %s", err))
	}
	if err = d.Set("state", projectConfig.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("update_available", projectConfig.UpdateAvailable); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting update_available: %s", err))
	}
	if err = d.Set("project_config_id", projectConfig.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_config_id: %s", err))
	}

	return nil
}

func resourceIbmProjectConfigUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateConfigOptions := &projectv1.UpdateConfigOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateConfigOptions.SetProjectID(parts[0])
	updateConfigOptions.SetID(parts[1])

	hasChange := false

	if d.HasChange("definition") {
		definition, err := resourceIbmProjectConfigMapToProjectConfigPrototypePatchDefinitionBlock(d.Get("definition.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateConfigOptions.SetDefinition(definition)
		hasChange = true
	}

	if hasChange {
		_, response, err := projectClient.UpdateConfigWithContext(context, updateConfigOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateConfigWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateConfigWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmProjectConfigRead(context, d, meta)
}

func resourceIbmProjectConfigDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteConfigOptions := &projectv1.DeleteConfigOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteConfigOptions.SetProjectID(parts[0])
	deleteConfigOptions.SetID(parts[1])

	_, response, err := projectClient.DeleteConfigWithContext(context, deleteConfigOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteConfigWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteConfigWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmProjectConfigMapToProjectConfigPrototypeDefinitionBlock(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototypeDefinitionBlock, error) {
	model := &projectv1.ProjectConfigPrototypeDefinitionBlock{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectConfigMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := resourceIbmProjectConfigMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	if modelMap["locator_id"] != nil && modelMap["locator_id"].(string) != "" {
		model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	}
	if modelMap["inputs"] != nil && len(modelMap["inputs"].([]interface{})) > 0 {
		InputsModel, err := resourceIbmProjectConfigMapToInputVariable(modelMap["inputs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Inputs = InputsModel
	}
	if modelMap["settings"] != nil && len(modelMap["settings"].([]interface{})) > 0 {
		SettingsModel, err := resourceIbmProjectConfigMapToProjectConfigSetting(modelMap["settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Settings = SettingsModel
	}
	return model, nil
}

func resourceIbmProjectConfigMapToProjectConfigAuth(modelMap map[string]interface{}) (*projectv1.ProjectConfigAuth, error) {
	model := &projectv1.ProjectConfigAuth{}
	if modelMap["trusted_profile_id"] != nil && modelMap["trusted_profile_id"].(string) != "" {
		model.TrustedProfileID = core.StringPtr(modelMap["trusted_profile_id"].(string))
	}
	if modelMap["method"] != nil && modelMap["method"].(string) != "" {
		model.Method = core.StringPtr(modelMap["method"].(string))
	}
	if modelMap["api_key"] != nil && modelMap["api_key"].(string) != "" {
		model.ApiKey = core.StringPtr(modelMap["api_key"].(string))
	}
	return model, nil
}

func resourceIbmProjectConfigMapToProjectComplianceProfile(modelMap map[string]interface{}) (*projectv1.ProjectComplianceProfile, error) {
	model := &projectv1.ProjectComplianceProfile{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["instance_id"] != nil && modelMap["instance_id"].(string) != "" {
		model.InstanceID = core.StringPtr(modelMap["instance_id"].(string))
	}
	if modelMap["instance_location"] != nil && modelMap["instance_location"].(string) != "" {
		model.InstanceLocation = core.StringPtr(modelMap["instance_location"].(string))
	}
	if modelMap["attachment_id"] != nil && modelMap["attachment_id"].(string) != "" {
		model.AttachmentID = core.StringPtr(modelMap["attachment_id"].(string))
	}
	if modelMap["profile_name"] != nil && modelMap["profile_name"].(string) != "" {
		model.ProfileName = core.StringPtr(modelMap["profile_name"].(string))
	}
	return model, nil
}

func resourceIbmProjectConfigMapToInputVariable(modelMap map[string]interface{}) (*projectv1.InputVariable, error) {
	model := &projectv1.InputVariable{}
	return model, nil
}

func resourceIbmProjectConfigMapToProjectConfigSetting(modelMap map[string]interface{}) (*projectv1.ProjectConfigSetting, error) {
	model := &projectv1.ProjectConfigSetting{}
	return model, nil
}

func resourceIbmProjectConfigMapToSchematicsWorkspace(modelMap map[string]interface{}) (*projectv1.SchematicsWorkspace, error) {
	model := &projectv1.SchematicsWorkspace{}
	if modelMap["workspace_crn"] != nil && modelMap["workspace_crn"].(string) != "" {
		model.WorkspaceCrn = core.StringPtr(modelMap["workspace_crn"].(string))
	}
	return model, nil
}

func resourceIbmProjectConfigMapToProjectConfigPrototypePatchDefinitionBlock(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototypePatchDefinitionBlock, error) {
	model := &projectv1.ProjectConfigPrototypePatchDefinitionBlock{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["environment"] != nil && modelMap["environment"].(string) != "" {
		model.Environment = core.StringPtr(modelMap["environment"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectConfigMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := resourceIbmProjectConfigMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	if modelMap["locator_id"] != nil && modelMap["locator_id"].(string) != "" {
		model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	}
	if modelMap["inputs"] != nil && len(modelMap["inputs"].([]interface{})) > 0 {
		InputsModel, err := resourceIbmProjectConfigMapToInputVariable(modelMap["inputs"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Inputs = InputsModel
	}
	if modelMap["settings"] != nil && len(modelMap["settings"].([]interface{})) > 0 {
		SettingsModel, err := resourceIbmProjectConfigMapToProjectConfigSetting(modelMap["settings"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Settings = SettingsModel
	}
	return model, nil
}

func resourceIbmProjectConfigSchematicsWorkspaceToMap(model *projectv1.SchematicsWorkspace) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.WorkspaceCrn != nil {
		modelMap["workspace_crn"] = model.WorkspaceCrn
	}
	return modelMap, nil
}

func resourceIbmProjectConfigProjectConfigResponseDefinitionToMap(model *projectv1.ProjectConfigResponseDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Environment != nil {
		modelMap["environment"] = model.Environment
	}
	if model.Authorizations != nil {
		authorizationsMap, err := resourceIbmProjectConfigProjectConfigAuthToMap(model.Authorizations)
		if err != nil {
			return modelMap, err
		}
		modelMap["authorizations"] = []map[string]interface{}{authorizationsMap}
	}
	if model.ComplianceProfile != nil {
		complianceProfileMap, err := resourceIbmProjectConfigProjectComplianceProfileToMap(model.ComplianceProfile)
		if err != nil {
			return modelMap, err
		}
		if len(complianceProfileMap) > 0 {
			modelMap["compliance_profile"] = []map[string]interface{}{complianceProfileMap}
		}
	}
	modelMap["locator_id"] = model.LocatorID
	if model.Inputs != nil {
		inputsMap, err := resourceIbmProjectConfigInputVariableToMap(model.Inputs)
		if err != nil {
			return modelMap, err
		}
		if len(inputsMap) > 0 {
			modelMap["inputs"] = []map[string]interface{}{inputsMap}
		}
	}
	if model.Settings != nil {
		settingsMap, err := resourceIbmProjectConfigProjectConfigSettingToMap(model.Settings)
		if err != nil {
			return modelMap, err
		}
		modelMap["settings"] = []map[string]interface{}{settingsMap}
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	return modelMap, nil
}

func resourceIbmProjectConfigProjectConfigAuthToMap(model *projectv1.ProjectConfigAuth) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = model.TrustedProfileID
	}
	if model.Method != nil {
		modelMap["method"] = model.Method
	}
	if model.ApiKey != nil {
		modelMap["api_key"] = model.ApiKey
	}
	return modelMap, nil
}

func resourceIbmProjectConfigProjectComplianceProfileToMap(model *projectv1.ProjectComplianceProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = model.InstanceID
	}
	if model.InstanceLocation != nil {
		modelMap["instance_location"] = model.InstanceLocation
	}
	if model.AttachmentID != nil {
		modelMap["attachment_id"] = model.AttachmentID
	}
	if model.ProfileName != nil {
		modelMap["profile_name"] = model.ProfileName
	}
	return modelMap, nil
}

func resourceIbmProjectConfigInputVariableToMap(model *projectv1.InputVariable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func resourceIbmProjectConfigProjectConfigSettingToMap(model *projectv1.ProjectConfigSetting) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	return modelMap, nil
}

func resourceIbmProjectConfigOutputValueToMap(model *projectv1.OutputValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIbmProjectConfigProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	definitionMap, err := resourceIbmProjectConfigProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = model.Crn
	modelMap["href"] = model.Href
	return modelMap, nil
}

func resourceIbmProjectConfigProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}
