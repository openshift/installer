// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package project

import (
	"context"
	"encoding/json"
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the environment.  It is unique within the account across projects and regions.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The description of the environment.",
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
										Description: "The IBM Cloud API Key.",
									},
								},
							},
						},
						"inputs": &schema.Schema{
							Type:        schema.TypeMap,
							Optional:    true,
							Description: "The input variables for configuration definition and environment.",
							Elem:        &schema.Schema{Type: schema.TypeString},
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
										Description: "The unique ID for that compliance profile.",
									},
									"instance_id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A unique ID for an instance of a compliance profile.",
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
				Description: "The project referenced by this resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID.",
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
							Description: "An IBM Cloud resource name, which uniquely identifies a resource.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time format as specified by RFC 3339.",
			},
			"target_account": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target account ID derived from the authentication block values.",
			},
			"modified_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time format as specified by RFC 3339.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL.",
			},
			"project_environment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The environment id as a friendly name.",
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
		return diag.FromErr(err)
	}

	createProjectEnvironmentOptions := &projectv1.CreateProjectEnvironmentOptions{}

	createProjectEnvironmentOptions.SetProjectID(d.Get("project_id").(string))
	definitionModel, err := resourceIbmProjectEnvironmentMapToEnvironmentDefinitionRequiredProperties(d.Get("definition.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createProjectEnvironmentOptions.SetDefinition(definitionModel)

	environment, response, err := projectClient.CreateProjectEnvironmentWithContext(context, createProjectEnvironmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProjectEnvironmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProjectEnvironmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createProjectEnvironmentOptions.ProjectID, *environment.ID))

	return resourceIbmProjectEnvironmentRead(context, d, meta)
}

func resourceIbmProjectEnvironmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getProjectEnvironmentOptions := &projectv1.GetProjectEnvironmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getProjectEnvironmentOptions.SetProjectID(parts[0])
	getProjectEnvironmentOptions.SetID(parts[1])

	environment, response, err := projectClient.GetProjectEnvironmentWithContext(context, getProjectEnvironmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProjectEnvironmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProjectEnvironmentWithContext failed %s\n%s", err, response))
	}

	definitionMap, err := resourceIbmProjectEnvironmentEnvironmentDefinitionRequiredPropertiesToMap(environment.Definition)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("definition", []map[string]interface{}{definitionMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definition: %s", err))
	}
	projectMap, err := resourceIbmProjectEnvironmentProjectReferenceToMap(environment.Project)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("project", []map[string]interface{}{projectMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(environment.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if !core.IsNil(environment.TargetAccount) {
		if err = d.Set("target_account", environment.TargetAccount); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting target_account: %s", err))
		}
	}
	if err = d.Set("modified_at", flex.DateTimeToString(environment.ModifiedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting modified_at: %s", err))
	}
	if err = d.Set("href", environment.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("project_environment_id", environment.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_environment_id: %s", err))
	}

	return nil
}

func resourceIbmProjectEnvironmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateProjectEnvironmentOptions := &projectv1.UpdateProjectEnvironmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateProjectEnvironmentOptions.SetProjectID(parts[0])
	updateProjectEnvironmentOptions.SetID(parts[1])

	hasChange := false

	if d.HasChange("project_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id"))
	}
	if d.HasChange("definition") {
		definition, err := resourceIbmProjectEnvironmentMapToEnvironmentDefinitionProperties(d.Get("definition.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateProjectEnvironmentOptions.SetDefinition(definition)
		hasChange = true
	}

	if hasChange {
		_, response, err := projectClient.UpdateProjectEnvironmentWithContext(context, updateProjectEnvironmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateProjectEnvironmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateProjectEnvironmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmProjectEnvironmentRead(context, d, meta)
}

func resourceIbmProjectEnvironmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProjectEnvironmentOptions := &projectv1.DeleteProjectEnvironmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProjectEnvironmentOptions.SetProjectID(parts[0])
	deleteProjectEnvironmentOptions.SetID(parts[1])

	_, response, err := projectClient.DeleteProjectEnvironmentWithContext(context, deleteProjectEnvironmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProjectEnvironmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteProjectEnvironmentWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmProjectEnvironmentMapToEnvironmentDefinitionRequiredProperties(modelMap map[string]interface{}) (*projectv1.EnvironmentDefinitionRequiredProperties, error) {
	model := &projectv1.EnvironmentDefinitionRequiredProperties{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectEnvironmentMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		bytes, _ := json.Marshal(modelMap["inputs"].(map[string]interface{}))
		newMap := make(map[string]interface{})
		json.Unmarshal(bytes, &newMap)
		if len(newMap) > 0 {
			model.Inputs = newMap
		}
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := resourceIbmProjectEnvironmentMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	return model, nil
}

func resourceIbmProjectEnvironmentMapToProjectConfigAuth(modelMap map[string]interface{}) (*projectv1.ProjectConfigAuth, error) {
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

func resourceIbmProjectEnvironmentMapToProjectComplianceProfile(modelMap map[string]interface{}) (*projectv1.ProjectComplianceProfile, error) {
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

func resourceIbmProjectEnvironmentMapToEnvironmentDefinitionProperties(modelMap map[string]interface{}) (*projectv1.EnvironmentDefinitionProperties, error) {
	model := &projectv1.EnvironmentDefinitionProperties{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectEnvironmentMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		bytes, _ := json.Marshal(modelMap["inputs"].(map[string]interface{}))
		newMap := make(map[string]interface{})
		json.Unmarshal(bytes, &newMap)
		if len(newMap) > 0 {
			model.Inputs = newMap
		}
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := resourceIbmProjectEnvironmentMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	return model, nil
}

func resourceIbmProjectEnvironmentEnvironmentDefinitionRequiredPropertiesToMap(model *projectv1.EnvironmentDefinitionRequiredProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Authorizations != nil {
		authorizationsMap, err := resourceIbmProjectEnvironmentProjectConfigAuthToMap(model.Authorizations)
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
			bytes, err := json.Marshal(v)
			if err != nil {
				return modelMap, err
			}
			inputs[k] = string(bytes)
		}
		if len(inputs) > 0 {
			modelMap["inputs"] = inputs
		}
	}
	if model.ComplianceProfile != nil {
		complianceProfileMap, err := resourceIbmProjectEnvironmentProjectComplianceProfileToMap(model.ComplianceProfile)
		if err != nil {
			return modelMap, err
		}
		if len(complianceProfileMap) > 0 {
			modelMap["compliance_profile"] = []map[string]interface{}{complianceProfileMap}
		}
	}
	return modelMap, nil
}

func resourceIbmProjectEnvironmentProjectConfigAuthToMap(model *projectv1.ProjectConfigAuth) (map[string]interface{}, error) {
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

func resourceIbmProjectEnvironmentProjectComplianceProfileToMap(model *projectv1.ProjectComplianceProfile) (map[string]interface{}, error) {
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

func resourceIbmProjectEnvironmentProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	definitionMap, err := resourceIbmProjectEnvironmentProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = model.Crn
	modelMap["href"] = model.Href
	return modelMap, nil
}

func resourceIbmProjectEnvironmentProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}
