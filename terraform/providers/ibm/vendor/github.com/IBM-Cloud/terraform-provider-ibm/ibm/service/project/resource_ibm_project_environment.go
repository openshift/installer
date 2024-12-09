// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.92.1-44330004-20240620-143510
 */

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

func ResourceIbmProjectEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmProjectEnvironmentCreate,
		ReadContext:   resourceIbmProjectEnvironmentRead,
		UpdateContext: resourceIbmProjectEnvironmentUpdate,
		DeleteContext: resourceIbmProjectEnvironmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_project_environment", "project_id"),
				Description:  "The unique project ID.",
			},
			"definition": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The environment definition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "",
							Description: "The description of the environment.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the environment. It's unique within the account across projects and regions.",
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
										Sensitive:   true,
										Description: "The IBM Cloud API Key. It can be either raw or pulled from the catalog via a `CRN` or `JSON` blob.",
									},
								},
							},
						},
						"inputs": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The input variables that are used for configuration definition and environment.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"compliance_profile": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The profile that is required for compliance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique ID for the compliance profile.",
									},
									"instance_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A unique ID for the instance of a compliance profile.",
									},
									"instance_location": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The location of the compliance instance.",
									},
									"attachment_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A unique ID for the attachment to a compliance profile.",
									},
									"profile_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the compliance profile.",
									},
								},
							},
						},
					},
				},
			},
			"project": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project that is referenced by this resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The definition of the project reference.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the project.",
									},
								},
							},
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An IBM Cloud resource name that uniquely identifies a resource.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"target_account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target account ID derived from the authentication block values. The target account exists only if the environment currently has an authorization block.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL.",
			},
			"project_environment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The environment ID as a friendly name.",
			},
		},
	}
}

func ResourceIbmProjectEnvironmentValidator() *validate.ResourceValidator {
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_project_environment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmProjectEnvironmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project_environment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createProjectEnvironmentOptions := &projectv1.CreateProjectEnvironmentOptions{}

	createProjectEnvironmentOptions.SetProjectID(d.Get("project_id").(string))
	definitionModel, err := ResourceIbmProjectEnvironmentMapToEnvironmentDefinitionRequiredProperties(d.Get("definition.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "create", "parse-definition").GetDiag()
	}
	createProjectEnvironmentOptions.SetDefinition(definitionModel)

	environment, _, err := projectClient.CreateProjectEnvironmentWithContext(context, createProjectEnvironmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateProjectEnvironmentWithContext failed: %s", err.Error()), "ibm_project_environment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createProjectEnvironmentOptions.ProjectID, *environment.ID))

	return resourceIbmProjectEnvironmentRead(context, d, meta)
}

func resourceIbmProjectEnvironmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project_environment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProjectEnvironmentOptions := &projectv1.GetProjectEnvironmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "sep-id-parts").GetDiag()
	}

	getProjectEnvironmentOptions.SetProjectID(parts[0])
	getProjectEnvironmentOptions.SetID(parts[1])

	environment, response, err := projectClient.GetProjectEnvironmentWithContext(context, getProjectEnvironmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProjectEnvironmentWithContext failed: %s", err.Error()), "ibm_project_environment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	definitionMap, err := ResourceIbmProjectEnvironmentEnvironmentDefinitionRequiredPropertiesResponseToMap(environment.Definition)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "definition-to-map").GetDiag()
	}
	if err = d.Set("definition", []map[string]interface{}{definitionMap}); err != nil {
		err = fmt.Errorf("Error setting definition: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-definition").GetDiag()
	}
	projectMap, err := ResourceIbmProjectEnvironmentProjectReferenceToMap(environment.Project)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "project-to-map").GetDiag()
	}
	if err = d.Set("project", []map[string]interface{}{projectMap}); err != nil {
		err = fmt.Errorf("Error setting project: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-project").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(environment.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-created_at").GetDiag()
	}
	if !core.IsNil(environment.TargetAccount) {
		if err = d.Set("target_account", environment.TargetAccount); err != nil {
			err = fmt.Errorf("Error setting target_account: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-target_account").GetDiag()
		}
	}
	if err = d.Set("modified_at", flex.DateTimeToString(environment.ModifiedAt)); err != nil {
		err = fmt.Errorf("Error setting modified_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-modified_at").GetDiag()
	}
	if err = d.Set("href", environment.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-href").GetDiag()
	}
	if err = d.Set("project_environment_id", environment.ID); err != nil {
		err = fmt.Errorf("Error setting project_environment_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "read", "set-project_environment_id").GetDiag()
	}

	return nil
}

func resourceIbmProjectEnvironmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project_environment", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateProjectEnvironmentOptions := &projectv1.UpdateProjectEnvironmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "update", "sep-id-parts").GetDiag()
	}

	updateProjectEnvironmentOptions.SetProjectID(parts[0])
	updateProjectEnvironmentOptions.SetID(parts[1])

	hasChange := false

	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_project_environment", "update", "project_id-forces-new").GetDiag()
	}
	if d.HasChange("definition") {
		definition, err := ResourceIbmProjectEnvironmentMapToEnvironmentDefinitionPropertiesPatch(d.Get("definition.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "update", "parse-definition").GetDiag()
		}
		updateProjectEnvironmentOptions.SetDefinition(definition)
		hasChange = true
	}

	if hasChange {
		_, _, err = projectClient.UpdateProjectEnvironmentWithContext(context, updateProjectEnvironmentOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateProjectEnvironmentWithContext failed: %s", err.Error()), "ibm_project_environment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmProjectEnvironmentRead(context, d, meta)
}

func resourceIbmProjectEnvironmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project_environment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteProjectEnvironmentOptions := &projectv1.DeleteProjectEnvironmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project_environment", "delete", "sep-id-parts").GetDiag()
	}

	deleteProjectEnvironmentOptions.SetProjectID(parts[0])
	deleteProjectEnvironmentOptions.SetID(parts[1])

	_, _, err = projectClient.DeleteProjectEnvironmentWithContext(context, deleteProjectEnvironmentOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteProjectEnvironmentWithContext failed: %s", err.Error()), "ibm_project_environment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmProjectEnvironmentMapToEnvironmentDefinitionRequiredProperties(modelMap map[string]interface{}) (*projectv1.EnvironmentDefinitionRequiredProperties, error) {
	model := &projectv1.EnvironmentDefinitionRequiredProperties{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := ResourceIbmProjectEnvironmentMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		model.Inputs = modelMap["inputs"].(map[string]interface{})
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := ResourceIbmProjectEnvironmentMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	return model, nil
}

func ResourceIbmProjectEnvironmentMapToProjectConfigAuth(modelMap map[string]interface{}) (*projectv1.ProjectConfigAuth, error) {
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

func ResourceIbmProjectEnvironmentMapToProjectComplianceProfile(modelMap map[string]interface{}) (*projectv1.ProjectComplianceProfile, error) {
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

func ResourceIbmProjectEnvironmentMapToEnvironmentDefinitionPropertiesPatch(modelMap map[string]interface{}) (*projectv1.EnvironmentDefinitionPropertiesPatch, error) {
	model := &projectv1.EnvironmentDefinitionPropertiesPatch{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := ResourceIbmProjectEnvironmentMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		model.Inputs = modelMap["inputs"].(map[string]interface{})
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := ResourceIbmProjectEnvironmentMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	return model, nil
}

func ResourceIbmProjectEnvironmentEnvironmentDefinitionRequiredPropertiesResponseToMap(model *projectv1.EnvironmentDefinitionRequiredPropertiesResponse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["description"] = *model.Description
	modelMap["name"] = *model.Name
	if model.Authorizations != nil {
		authorizationsMap, err := ResourceIbmProjectEnvironmentProjectConfigAuthToMap(model.Authorizations)
		if err != nil {
			return modelMap, err
		}
		if len(authorizationsMap) > 0 {
			modelMap["authorizations"] = []map[string]interface{}{authorizationsMap}
		}
	}
	if model.Inputs != nil {
		inputs := make(map[string]interface{})
		for k, v := range model.Inputs {
			inputs[k] = flex.Stringify(v)
		}
		modelMap["inputs"] = inputs
	}
	if model.ComplianceProfile != nil {
		complianceProfileMap, err := ResourceIbmProjectEnvironmentProjectComplianceProfileToMap(model.ComplianceProfile)
		if err != nil {
			return modelMap, err
		}
		if len(complianceProfileMap) > 0 {
			modelMap["compliance_profile"] = []map[string]interface{}{complianceProfileMap}
		}
	}
	return modelMap, nil
}

func ResourceIbmProjectEnvironmentProjectConfigAuthToMap(model *projectv1.ProjectConfigAuth) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.TrustedProfileID != nil {
		modelMap["trusted_profile_id"] = *model.TrustedProfileID
	}
	if model.Method != nil {
		modelMap["method"] = *model.Method
	}
	if model.ApiKey != nil {
		modelMap["api_key"] = *model.ApiKey
	}
	return modelMap, nil
}

func ResourceIbmProjectEnvironmentProjectComplianceProfileToMap(model *projectv1.ProjectComplianceProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.InstanceID != nil {
		modelMap["instance_id"] = *model.InstanceID
	}
	if model.InstanceLocation != nil {
		modelMap["instance_location"] = *model.InstanceLocation
	}
	if model.AttachmentID != nil {
		modelMap["attachment_id"] = *model.AttachmentID
	}
	if model.ProfileName != nil {
		modelMap["profile_name"] = *model.ProfileName
	}
	return modelMap, nil
}

func ResourceIbmProjectEnvironmentProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["href"] = *model.Href
	definitionMap, err := ResourceIbmProjectEnvironmentProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = *model.Crn
	return modelMap, nil
}

func ResourceIbmProjectEnvironmentProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	return modelMap, nil
}
