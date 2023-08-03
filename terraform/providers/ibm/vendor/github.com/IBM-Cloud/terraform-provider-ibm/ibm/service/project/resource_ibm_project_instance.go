// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package project

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/project-go-sdk/projectv1"
)

func ResourceIbmProjectInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmProjectInstanceCreate,
		ReadContext:   resourceIbmProjectInstanceRead,
		UpdateContext: resourceIbmProjectInstanceUpdate,
		DeleteContext: resourceIbmProjectInstanceDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The project name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A project's descriptive text.",
			},
			"configs": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The project configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The configuration name.",
						},
						"labels": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A collection of configuration labels.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The project configuration description.",
						},
						"locator_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "A dotted value of catalogID.versionID.",
						},
						"input": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "The input values to use to deploy the configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The variable name.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Can be any value - a string, number, boolean, array, or object.",
									},
								},
							},
						},
						"setting": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Schematics environment variables to use to deploy the configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the configuration setting.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The value of the configuration setting.",
									},
								},
							},
						},
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource group where the project's data and tools are created.",
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The location where the project's data and tools are created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An IBM Cloud resource name, which uniquely identifies a resource.",
			},
			"metadata": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The metadata of the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The cumulative list of needs attention items for a project.",
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
										Description: "The unique ID of a project.",
									},
									"config_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique ID of a project.",
									},
									"config_version": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the configuration.",
									},
								},
							},
						},
						"cumulative_needs_attention_view_err": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "True indicates that the fetch of the needs attention items failed.",
						},
						"location": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IBM Cloud location where a resource is deployed.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource group where the project's data and tools are created.",
						},
						"state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project status value.",
						},
						"event_notifications_crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the event notifications instance if one is connected to this project.",
						},
					},
				},
			},
		},
	}
}

func ResourceIbmProjectInstanceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:     "name",
			Type:           validate.TypeString,
			Required:       true,
			Regexp:         `^(?!\s)(?!.*\s$)[^'"<>{}\x00-\x1F]+$`,
			MinValueLength: 1,
			MaxValueLength: 64,
		},
		validate.ValidateSchema{
			Identifier:     "description",
			Type:           validate.TypeString,
			Optional:       true,
			Regexp:         `^$|^(?!\s).*\S$`,
			MinValueLength: 0,
			MaxValueLength: 1024,
		},
		validate.ValidateSchema{
			Identifier:     "resource_group",
			Type:           validate.TypeString,
			Required:       true,
			Regexp:         `^$|^(?!\s)(?!.*\s$)[^'"<>{}\x00-\x1F]*$`,
			MinValueLength: 0,
			MaxValueLength: 40,
		},
		validate.ValidateSchema{
			Identifier:     "location",
			Type:           validate.TypeString,
			Required:       true,
			Regexp:         `^$|^(us-south|us-east|eu-gb|eu-de)$`,
			MinValueLength: 0,
			MaxValueLength: 12,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_project_instance", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmProjectInstanceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	createProjectOptions := &projectv1.CreateProjectOptions{}

	createProjectOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("description"); ok {
		createProjectOptions.SetDescription(d.Get("description").(string))
	}
	if _, ok := d.GetOk("configs"); ok {
		var configs []projectv1.ProjectConfigPrototype
		for _, v := range d.Get("configs").([]interface{}) {
			value := v.(map[string]interface{})
			configsItem, err := resourceIbmProjectInstanceMapToProjectConfigPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			configs = append(configs, *configsItem)
		}
		createProjectOptions.SetConfigs(configs)
	}
	if _, ok := d.GetOk("resource_group"); ok {
		createProjectOptions.SetResourceGroup(d.Get("resource_group").(string))
	}
	if _, ok := d.GetOk("location"); ok {
		createProjectOptions.SetLocation(d.Get("location").(string))
	}

	project, response, err := projectClient.CreateProjectWithContext(context, createProjectOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateProjectWithContext failed %s\n%s", err, response))
	}

	d.SetId(*project.ID)

	_, err = waitForProjectInstanceCreate(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for create project instance (%s) to be succeeded: %s", d.Id(), err))
	}

	return resourceIbmProjectInstanceRead(context, d, meta)
}

func waitForProjectInstanceCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	getProjectOptions := &projectv1.GetProjectOptions{}
	getProjectOptions.SetID(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"not_exists"},
		Target:  []string{"exists"},
		Refresh: func() (interface{}, string, error) {
			_, resp, err := projectClient.GetProject(getProjectOptions)
			if err == nil {
				if resp != nil && resp.StatusCode == 200 {
					return resp, "exists", nil
				} else {
					return resp, "not_exists", nil
				}
			} else {
				return nil, "", fmt.Errorf("[ERROR] Get the project instance %s failed with resp code: %d, err: %v", d.Id(), resp.StatusCode, err)
			}
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      2 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmProjectInstanceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if !core.IsNil(project.Crn) {
		if err = d.Set("crn", project.Crn); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
		}
	}
	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if !core.IsNil(project.Description) {
		if err = d.Set("description", project.Description); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
		}
	}
	if !core.IsNil(project.Configs) {
		configs := []map[string]interface{}{}
		if err = d.Set("configs", configs); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting configs: %s", err))
		}
	}
	if !core.IsNil(project.Metadata) {
		metadataMap, err := resourceIbmProjectInstanceProjectMetadataToMap(project.Metadata)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("metadata", []map[string]interface{}{metadataMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting metadata: %s", err))
		}
	}

	return nil
}

func resourceIbmProjectInstanceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	updateProjectOptions := &projectv1.UpdateProjectOptions{}

	updateProjectOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("name") {
		updateProjectOptions.SetName(d.Get("name").(string))
		hasChange = true
	}
	if d.HasChange("description") {
		updateProjectOptions.SetDescription(d.Get("description").(string))
		hasChange = true
	}
	if d.HasChange("configs") {
		var configs []projectv1.ProjectConfigPrototype
		for _, v := range d.Get("configs").([]interface{}) {
			value := v.(map[string]interface{})
			configsItem, err := resourceIbmProjectInstanceMapToProjectConfigPrototype(value)
			if err != nil {
				return diag.FromErr(err)
			}
			configs = append(configs, *configsItem)
		}
		hasChange = true
	}

	if hasChange {
		_, response, err := projectClient.UpdateProjectWithContext(context, updateProjectOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateProjectWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateProjectWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmProjectInstanceRead(context, d, meta)
}

func resourceIbmProjectInstanceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, err = waitForProjectInstanceDelete(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error waiting for delete project instance (%s) to be succeeded: %s", d.Id(), err))
	}

	d.SetId("")

	return nil
}

func waitForProjectInstanceDelete(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return false, err
	}
	instanceID := d.Id()
	getProjectOptions := &projectv1.GetProjectOptions{}
	getProjectOptions.SetID(instanceID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"exists"},
		Target:  []string{"not_exists"},
		Refresh: func() (interface{}, string, error) {
			_, resp, err := projectClient.GetProject(getProjectOptions)
			if err != nil {
				if resp != nil && resp.StatusCode == 404 {
					return resp, "not_exists", nil
				} else {
					return resp, "exists", nil
				}
			} else {
				return resp, "exists", nil
			}
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      2 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmProjectInstanceMapToProjectConfigPrototype(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototype, error) {
	model := &projectv1.ProjectConfigPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["labels"] != nil {
		labels := []string{}
		for _, labelsItem := range modelMap["labels"].([]interface{}) {
			labels = append(labels, labelsItem.(string))
		}
		model.Labels = labels
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	if modelMap["input"] != nil {
		input := []projectv1.ProjectConfigInputVariable{}
		for _, inputItem := range modelMap["input"].([]interface{}) {
			inputItemModel, err := resourceIbmProjectInstanceMapToProjectConfigInputVariable(inputItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			input = append(input, *inputItemModel)
		}
		model.Input = input
	}
	if modelMap["setting"] != nil {
		setting := []projectv1.ProjectConfigSettingCollection{}
		for _, settingItem := range modelMap["setting"].([]interface{}) {
			settingItemModel, err := resourceIbmProjectInstanceMapToProjectConfigSettingCollection(settingItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			setting = append(setting, *settingItemModel)
		}
		model.Setting = setting
	}
	return model, nil
}

func resourceIbmProjectInstanceMapToProjectConfigInputVariable(modelMap map[string]interface{}) (*projectv1.ProjectConfigInputVariable, error) {
	model := &projectv1.ProjectConfigInputVariable{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["value"] != nil {
		model.Value = modelMap["value"].(string)
	}
	return model, nil
}

func resourceIbmProjectInstanceMapToProjectConfigSettingCollection(modelMap map[string]interface{}) (*projectv1.ProjectConfigSettingCollection, error) {
	model := &projectv1.ProjectConfigSettingCollection{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.Value = core.StringPtr(modelMap["value"].(string))
	return model, nil
}

func resourceIbmProjectInstanceProjectConfigPrototypeToMap(model *projectv1.ProjectConfigPrototype) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = model.ID
	}
	modelMap["name"] = model.Name
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["locator_id"] = model.LocatorID
	if model.Input != nil {
		input := []map[string]interface{}{}
		for _, inputItem := range model.Input {
			inputItemMap, err := resourceIbmProjectInstanceProjectConfigInputVariableToMap(&inputItem)
			if err != nil {
				return modelMap, err
			}
			input = append(input, inputItemMap)
		}
		modelMap["input"] = input
	}
	if model.Setting != nil {
		setting := []map[string]interface{}{}
		for _, settingItem := range model.Setting {
			settingItemMap, err := resourceIbmProjectInstanceProjectConfigSettingCollectionToMap(&settingItem)
			if err != nil {
				return modelMap, err
			}
			setting = append(setting, settingItemMap)
		}
		modelMap["setting"] = setting
	}
	return modelMap, nil
}

func resourceIbmProjectInstanceProjectConfigInputVariableToMap(model *projectv1.ProjectConfigInputVariable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func resourceIbmProjectInstanceProjectConfigSettingCollectionToMap(model *projectv1.ProjectConfigSettingCollection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["value"] = model.Value
	return modelMap, nil
}

func resourceIbmProjectInstanceProjectMetadataToMap(model *projectv1.ProjectMetadata) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Crn != nil {
		modelMap["crn"] = model.Crn
	}
	if model.CreatedAt != nil {
		modelMap["created_at"] = model.CreatedAt.String()
	}
	if model.CumulativeNeedsAttentionView != nil {
		cumulativeNeedsAttentionView := []map[string]interface{}{}
		for _, cumulativeNeedsAttentionViewItem := range model.CumulativeNeedsAttentionView {
			cumulativeNeedsAttentionViewItemMap, err := resourceIbmProjectInstanceCumulativeNeedsAttentionToMap(&cumulativeNeedsAttentionViewItem)
			if err != nil {
				return modelMap, err
			}
			cumulativeNeedsAttentionView = append(cumulativeNeedsAttentionView, cumulativeNeedsAttentionViewItemMap)
		}
		modelMap["cumulative_needs_attention_view"] = cumulativeNeedsAttentionView
	}
	if model.CumulativeNeedsAttentionViewErr != nil {
		modelMap["cumulative_needs_attention_view_err"] = model.CumulativeNeedsAttentionViewErr
	}
	if model.Location != nil {
		modelMap["location"] = model.Location
	}
	if model.ResourceGroup != nil {
		modelMap["resource_group"] = model.ResourceGroup
	}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	if model.EventNotificationsCrn != nil {
		modelMap["event_notifications_crn"] = model.EventNotificationsCrn
	}
	return modelMap, nil
}

func resourceIbmProjectInstanceCumulativeNeedsAttentionToMap(model *projectv1.CumulativeNeedsAttention) (map[string]interface{}, error) {
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
	return modelMap, nil
}
