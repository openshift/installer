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

func ResourceIbmProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmProjectCreate,
		ReadContext:   resourceIbmProjectRead,
		UpdateContext: resourceIbmProjectUpdate,
		DeleteContext: resourceIbmProjectDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"location": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_project", "location"),
				Description:  "The IBM Cloud location where a resource is deployed.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The resource group where the project's data and tools are created.",
			},
			"definition": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The definition of the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the project.  It is unique within the account across regions.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A brief explanation of the project's use in the configuration of a deployable architecture. It is possible to create a project without providing a description.",
						},
						"destroy_on_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Required:    true,
							Description: "The policy that indicates whether the resources are destroyed or not when a project is deleted.",
						},
					},
				},
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An IBM Cloud resource name, which uniquely identifies a resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ, matching the date and time format as specified by RFC 3339.",
			},
			"cumulative_needs_attention_view": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cumulative list of needs attention items for a project. If the view is successfully retrieved, an array which could be empty is returned.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The event name.",
						},
						"event_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for that individual event.",
						},
						"config_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A unique ID for the configuration.",
						},
						"config_version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number of the configuration.",
						},
					},
				},
			},
			"cumulative_needs_attention_view_error": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "True indicates that the fetch of the needs attention items failed. It only exists if there was an error while retrieving the cumulative needs attention view.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group id where the project's data and tools are created.",
			},
			"state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project status value.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A URL.",
			},
			"event_notifications_crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the event notifications instance if one is connected to this project.",
			},
			"configs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project configurations. These configurations are only included in the response of creating a project if a configs array is specified in the request payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version of the configuration.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The state of the configuration.",
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
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The name and description of a project configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configuration name. It is unique within the account across projects and regions.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A project configuration description.",
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
						"deployment_model": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration type.",
						},
					},
				},
			},
			"environments": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project environments. These environments are only included in the response if project environments were created on the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment id as a friendly name.",
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
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The environment definition used in the project collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the environment.  It is unique within the account across projects and regions.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the environment.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceIbmProjectValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "location",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^$|^(us-south|us-east|eu-gb|eu-de)$`,
			MinValueLength:             0,
			MaxValueLength:             12,
		},
		validate.ValidateSchema{
			Identifier:                 "resource_group",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^(?!\s)(?!.*\s$)[^'"` + "`" + `<>{}\x00-\x1F]*$`,
			MinValueLength:             0,
			MaxValueLength:             64,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_project", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmProjectCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createProjectOptions := &projectv1.CreateProjectOptions{}

	definitionModel, err := resourceIbmProjectMapToProjectPrototypeDefinition(d.Get("definition.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createProjectOptions.SetDefinition(definitionModel)
	createProjectOptions.SetLocation(d.Get("location").(string))
	createProjectOptions.SetResourceGroup(d.Get("resource_group").(string))
	if _, ok := d.GetOk("configs"); ok {
		var configs []projectv1.ProjectConfigPrototype
		for _, v := range d.Get("configs").([]interface{}) {
			value := v.(map[string]interface{})
			configsItem, err := resourceIbmProjectMapToProjectConfigPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			configs = append(configs, *configsItem)
		}
		createProjectOptions.SetConfigs(configs)
	}
	if _, ok := d.GetOk("environments"); ok {
		var environments []projectv1.EnvironmentPrototype
		for _, v := range d.Get("environments").([]interface{}) {
			value := v.(map[string]interface{})
			environmentsItem, err := resourceIbmProjectMapToEnvironmentPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			environments = append(environments, *environmentsItem)
		}
		createProjectOptions.SetEnvironments(environments)
	}

	project, response, err := projectClient.CreateProjectWithContext(context, createProjectOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProjectWithContext failed %s\n%s", err, response))
	}

	d.SetId(*project.ID)

	return resourceIbmProjectRead(context, d, meta)
}

func resourceIbmProjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getProjectOptions := &projectv1.GetProjectOptions{}

	getProjectOptions.SetID(d.Id())

	project, response, err := projectClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProjectWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("location", project.Location); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
	}
	if err = d.Set("resource_group", project.ResourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
	}
	definitionMap, err := resourceIbmProjectProjectDefinitionPropertiesToMap(project.Definition)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("definition", []map[string]interface{}{definitionMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definition: %s", err))
	}
	if err = d.Set("crn", project.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(project.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if !core.IsNil(project.CumulativeNeedsAttentionView) {
		cumulativeNeedsAttentionView := []map[string]interface{}{}
		for _, cumulativeNeedsAttentionViewItem := range project.CumulativeNeedsAttentionView {
			cumulativeNeedsAttentionViewItemMap, err := resourceIbmProjectCumulativeNeedsAttentionToMap(&cumulativeNeedsAttentionViewItem)
			if err != nil {
				return diag.FromErr(err)
			}
			cumulativeNeedsAttentionView = append(cumulativeNeedsAttentionView, cumulativeNeedsAttentionViewItemMap)
		}
		if err = d.Set("cumulative_needs_attention_view", cumulativeNeedsAttentionView); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cumulative_needs_attention_view: %s", err))
		}
	}
	if !core.IsNil(project.CumulativeNeedsAttentionViewError) {
		if err = d.Set("cumulative_needs_attention_view_error", project.CumulativeNeedsAttentionViewError); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cumulative_needs_attention_view_error: %s", err))
		}
	}
	if err = d.Set("resource_group_id", project.ResourceGroupID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group_id: %s", err))
	}
	if err = d.Set("state", project.State); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting state: %s", err))
	}
	if err = d.Set("href", project.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if !core.IsNil(project.EventNotificationsCrn) {
		if err = d.Set("event_notifications_crn", project.EventNotificationsCrn); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting event_notifications_crn: %s", err))
		}
	}
	if !core.IsNil(project.Configs) {
		configs := []map[string]interface{}{}
		for _, configsItem := range project.Configs {
			configsItemMap, err := resourceIbmProjectProjectConfigSummaryToMap(&configsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			configs = append(configs, configsItemMap)
		}
		if err = d.Set("configs", configs); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting configs: %s", err))
		}
	}
	if !core.IsNil(project.Environments) {
		environments := []map[string]interface{}{}
		for _, environmentsItem := range project.Environments {
			environmentsItemMap, err := resourceIbmProjectProjectEnvironmentSummaryToMap(&environmentsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			environments = append(environments, environmentsItemMap)
		}
		if err = d.Set("environments", environments); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting environments: %s", err))
		}
	}

	return nil
}

func resourceIbmProjectUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateProjectOptions := &projectv1.UpdateProjectOptions{}

	updateProjectOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("definition") {
		definition, err := resourceIbmProjectMapToProjectPatchDefinitionBlock(d.Get("definition.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		updateProjectOptions.SetDefinition(definition)
		hasChange = true
	}

	if hasChange {
		_, response, err := projectClient.UpdateProjectWithContext(context, updateProjectOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateProjectWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateProjectWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmProjectRead(context, d, meta)
}

func resourceIbmProjectDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProjectOptions := &projectv1.DeleteProjectOptions{}

	deleteProjectOptions.SetID(d.Id())

	response, err := projectClient.DeleteProjectWithContext(context, deleteProjectOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteProjectWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmProjectMapToProjectPrototypeDefinition(modelMap map[string]interface{}) (*projectv1.ProjectPrototypeDefinition, error) {
	model := &projectv1.ProjectPrototypeDefinition{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["destroy_on_delete"] != nil {
		model.DestroyOnDelete = core.BoolPtr(modelMap["destroy_on_delete"].(bool))
	}
	return model, nil
}

func resourceIbmProjectMapToProjectConfigPrototype(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototype, error) {
	model := &projectv1.ProjectConfigPrototype{}
	DefinitionModel, err := resourceIbmProjectMapToProjectConfigPrototypeDefinitionBlock(modelMap["definition"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Definition = DefinitionModel
	if modelMap["schematics"] != nil && len(modelMap["schematics"].([]interface{})) > 0 {
		SchematicsModel, err := resourceIbmProjectMapToSchematicsWorkspace(modelMap["schematics"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Schematics = SchematicsModel
	}
	return model, nil
}

func resourceIbmProjectMapToProjectConfigPrototypeDefinitionBlock(modelMap map[string]interface{}) (projectv1.ProjectConfigPrototypeDefinitionBlockIntf, error) {
	model := &projectv1.ProjectConfigPrototypeDefinitionBlock{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["environment_id"] != nil && modelMap["environment_id"].(string) != "" {
		model.EnvironmentID = core.StringPtr(modelMap["environment_id"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
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
	if modelMap["settings"] != nil {
		bytes, _ := json.Marshal(modelMap["settings"].(map[string]interface{}))
		newMap := make(map[string]interface{})
		json.Unmarshal(bytes, &newMap)
		if len(newMap) > 0 {
			model.Settings = newMap
		}
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := resourceIbmProjectMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	if modelMap["locator_id"] != nil && modelMap["locator_id"].(string) != "" {
		model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	}
	if modelMap["resource_crns"] != nil && len(modelMap["resource_crns"].([]interface{})) > 0 {
		resourceCrns := []string{}
		for _, resourceCrnsItem := range modelMap["resource_crns"].([]interface{}) {
			resourceCrns = append(resourceCrns, resourceCrnsItem.(string))
		}
		model.ResourceCrns = resourceCrns
	}
	return model, nil
}

func resourceIbmProjectMapToProjectConfigAuth(modelMap map[string]interface{}) (*projectv1.ProjectConfigAuth, error) {
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

func resourceIbmProjectMapToProjectComplianceProfile(modelMap map[string]interface{}) (*projectv1.ProjectComplianceProfile, error) {
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

func resourceIbmProjectMapToProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties, error) {
	model := &projectv1.ProjectConfigPrototypeDefinitionBlockDAConfigDefinitionProperties{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["environment_id"] != nil && modelMap["environment_id"].(string) != "" {
		model.EnvironmentID = core.StringPtr(modelMap["environment_id"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
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
	if modelMap["settings"] != nil {
		bytes, _ := json.Marshal(modelMap["settings"].(map[string]interface{}))
		newMap := make(map[string]interface{})
		json.Unmarshal(bytes, &newMap)
		if len(newMap) > 0 {
			model.Settings = newMap
		}
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := resourceIbmProjectMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	if modelMap["locator_id"] != nil && modelMap["locator_id"].(string) != "" {
		model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	}
	return model, nil
}

func resourceIbmProjectMapToProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties, error) {
	model := &projectv1.ProjectConfigPrototypeDefinitionBlockResourceConfigDefinitionProperties{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["environment_id"] != nil && modelMap["environment_id"].(string) != "" {
		model.EnvironmentID = core.StringPtr(modelMap["environment_id"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
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
	if modelMap["settings"] != nil {
		bytes, _ := json.Marshal(modelMap["settings"].(map[string]interface{}))
		newMap := make(map[string]interface{})
		json.Unmarshal(bytes, &newMap)
		if len(newMap) > 0 {
			model.Settings = newMap
		}
	}
	if modelMap["resource_crns"] != nil && len(modelMap["resource_crns"].([]interface{})) > 0 {
		resourceCrns := []string{}
		for _, resourceCrnsItem := range modelMap["resource_crns"].([]interface{}) {
			resourceCrns = append(resourceCrns, resourceCrnsItem.(string))
		}
		model.ResourceCrns = resourceCrns
	}
	return model, nil
}

func resourceIbmProjectMapToSchematicsWorkspace(modelMap map[string]interface{}) (*projectv1.SchematicsWorkspace, error) {
	model := &projectv1.SchematicsWorkspace{}
	if modelMap["workspace_crn"] != nil && modelMap["workspace_crn"].(string) != "" {
		model.WorkspaceCrn = core.StringPtr(modelMap["workspace_crn"].(string))
	}
	return model, nil
}

func resourceIbmProjectMapToEnvironmentPrototype(modelMap map[string]interface{}) (*projectv1.EnvironmentPrototype, error) {
	model := &projectv1.EnvironmentPrototype{}
	DefinitionModel, err := resourceIbmProjectMapToEnvironmentDefinitionRequiredProperties(modelMap["definition"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Definition = DefinitionModel
	return model, nil
}

func resourceIbmProjectMapToEnvironmentDefinitionRequiredProperties(modelMap map[string]interface{}) (*projectv1.EnvironmentDefinitionRequiredProperties, error) {
	model := &projectv1.EnvironmentDefinitionRequiredProperties{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := resourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
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
		ComplianceProfileModel, err := resourceIbmProjectMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	return model, nil
}

func resourceIbmProjectMapToProjectPatchDefinitionBlock(modelMap map[string]interface{}) (*projectv1.ProjectPatchDefinitionBlock, error) {
	model := &projectv1.ProjectPatchDefinitionBlock{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["destroy_on_delete"] != nil {
		model.DestroyOnDelete = core.BoolPtr(modelMap["destroy_on_delete"].(bool))
	}
	return model, nil
}

func resourceIbmProjectProjectDefinitionPropertiesToMap(model *projectv1.ProjectDefinitionProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["destroy_on_delete"] = model.DestroyOnDelete
	return modelMap, nil
}

func resourceIbmProjectCumulativeNeedsAttentionToMap(model *projectv1.CumulativeNeedsAttention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Event != nil {
		modelMap["event"] = model.Event
	}
	if model.EventID != nil {
		modelMap["event_id"] = model.EventID
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = model.ConfigID
	}
	if model.ConfigVersion != nil {
		modelMap["config_version"] = flex.IntValue(model.ConfigVersion)
	}
	return modelMap, nil
}

func resourceIbmProjectProjectConfigSummaryToMap(model *projectv1.ProjectConfigSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApprovedVersion != nil {
		approvedVersionMap, err := resourceIbmProjectProjectConfigVersionSummaryToMap(model.ApprovedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["approved_version"] = []map[string]interface{}{approvedVersionMap}
	}
	if model.DeployedVersion != nil {
		deployedVersionMap, err := resourceIbmProjectProjectConfigVersionSummaryToMap(model.DeployedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["deployed_version"] = []map[string]interface{}{deployedVersionMap}
	}
	modelMap["id"] = model.ID
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["state"] = model.State
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["modified_at"] = model.ModifiedAt.String()
	modelMap["href"] = model.Href
	definitionMap, err := resourceIbmProjectProjectConfigDefinitionNameDescriptionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	projectMap, err := resourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	if model.DeploymentModel != nil {
		modelMap["deployment_model"] = model.DeploymentModel
	}
	return modelMap, nil
}

func resourceIbmProjectProjectConfigVersionSummaryToMap(model *projectv1.ProjectConfigVersionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["state"] = model.State
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["href"] = model.Href
	return modelMap, nil
}

func resourceIbmProjectProjectConfigDefinitionNameDescriptionToMap(model *projectv1.ProjectConfigDefinitionNameDescription) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func resourceIbmProjectProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	definitionMap, err := resourceIbmProjectProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = model.Crn
	modelMap["href"] = model.Href
	return modelMap, nil
}

func resourceIbmProjectProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIbmProjectProjectEnvironmentSummaryToMap(model *projectv1.ProjectEnvironmentSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	projectMap, err := resourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	definitionMap, err := resourceIbmProjectEnvironmentDefinitionNameDescriptionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	return modelMap, nil
}

func resourceIbmProjectEnvironmentDefinitionNameDescriptionToMap(model *projectv1.EnvironmentDefinitionNameDescription) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}
