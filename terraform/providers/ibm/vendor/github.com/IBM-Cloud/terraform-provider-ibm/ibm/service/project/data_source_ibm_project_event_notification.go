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
	"github.com/IBM/project-go-sdk/projectv1"
)

func DataSourceIbmProjectEventNotification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmProjectEventNotificationRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique project ID.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project name.",
			},
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A project descriptive text.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An IBM Cloud resource name, which uniquely identifies a resource.",
			},
			"configs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the configuration. If this parameter is empty, an ID is automatically created for the configuration.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The configuration name.",
						},
						"labels": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A collection of configuration labels.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project configuration description.",
						},
						"locator_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A dotted value of catalogID.versionID.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of a project configuration manual property.",
						},
						"input": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The outputs of a Schematics template property.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The variable name.",
									},
									"type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The variable type.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Can be any value - a string, number, boolean, array, or object.",
									},
									"required": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the variable is required or not.",
									},
								},
							},
						},
						"output": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The outputs of a Schematics template property.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The variable name.",
									},
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A short explanation of the output value.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Can be any value - a string, number, boolean, array, or object.",
									},
								},
							},
						},
						"setting": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Schematics environment variables to use to deploy the configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the configuration setting.",
									},
									"value": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The value of the configuration setting.",
									},
								},
							},
						},
					},
				},
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

func dataSourceIbmProjectEventNotificationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		return diag.FromErr(err)
	}

	getProjectOptions := &projectv1.GetProjectOptions{}

	getProjectOptions.SetID(d.Get("project_id").(string))

	project, response, err := projectClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		log.Printf("[DEBUG] GetProjectWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetProjectWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getProjectOptions.ID))

	if err = d.Set("name", project.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("description", project.Description); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting description: %s", err))
	}

	if err = d.Set("crn", project.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	configs := []map[string]interface{}{}
	if project.Configs != nil {
		for _, modelItem := range project.Configs {
			modelMap, err := dataSourceIbmProjectEventNotificationProjectConfigToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			configs = append(configs, modelMap)
		}
	}
	if err = d.Set("configs", configs); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting configs %s", err))
	}

	metadata := []map[string]interface{}{}
	if project.Metadata != nil {
		modelMap, err := dataSourceIbmProjectEventNotificationProjectMetadataToMap(project.Metadata)
		if err != nil {
			return diag.FromErr(err)
		}
		metadata = append(metadata, modelMap)
	}
	if err = d.Set("metadata", metadata); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting metadata %s", err))
	}

	return nil
}

func dataSourceIbmProjectEventNotificationProjectConfigToMap(model *projectv1.ProjectConfig) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["project_id"] = model.ID
	}
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Labels != nil {
		modelMap["labels"] = model.Labels
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.LocatorID != nil {
		modelMap["locator_id"] = model.LocatorID
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Input != nil {
		input := []map[string]interface{}{}
		for _, inputItem := range model.Input {
			inputItemMap, err := dataSourceIbmProjectEventNotificationInputVariableToMap(&inputItem)
			if err != nil {
				return modelMap, err
			}
			input = append(input, inputItemMap)
		}
		modelMap["input"] = input
	}
	if model.Output != nil {
		output := []map[string]interface{}{}
		for _, outputItem := range model.Output {
			outputItemMap, err := dataSourceIbmProjectEventNotificationOutputValueToMap(&outputItem)
			if err != nil {
				return modelMap, err
			}
			output = append(output, outputItemMap)
		}
		modelMap["output"] = output
	}
	if model.Setting != nil {
		setting := []map[string]interface{}{}
		for _, settingItem := range model.Setting {
			settingItemMap, err := dataSourceIbmProjectEventNotificationProjectConfigSettingCollectionToMap(&settingItem)
			if err != nil {
				return modelMap, err
			}
			setting = append(setting, settingItemMap)
		}
		modelMap["setting"] = setting
	}
	return modelMap, nil
}

func dataSourceIbmProjectEventNotificationInputVariableToMap(model *projectv1.InputVariable) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Type != nil {
		modelMap["type"] = model.Type
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	if model.Required != nil {
		modelMap["required"] = model.Required
	}
	return modelMap, nil
}

func dataSourceIbmProjectEventNotificationOutputValueToMap(model *projectv1.OutputValue) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func dataSourceIbmProjectEventNotificationProjectConfigSettingCollectionToMap(model *projectv1.ProjectConfigSettingCollection) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Value != nil {
		modelMap["value"] = model.Value
	}
	return modelMap, nil
}

func dataSourceIbmProjectEventNotificationProjectMetadataToMap(model *projectv1.ProjectMetadata) (map[string]interface{}, error) {
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
			cumulativeNeedsAttentionViewItemMap, err := dataSourceIbmProjectEventNotificationCumulativeNeedsAttentionToMap(&cumulativeNeedsAttentionViewItem)
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

func dataSourceIbmProjectEventNotificationCumulativeNeedsAttentionToMap(model *projectv1.CumulativeNeedsAttention) (map[string]interface{}, error) {
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
