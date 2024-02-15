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

func DataSourceIbmProject() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmProjectRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique project ID.",
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
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IBM Cloud location where a resource is deployed.",
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
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group name where the project's data and tools are created.",
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
			"definition": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The definition of the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the project.  It is unique within the account across regions.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A brief explanation of the project's use in the configuration of a deployable architecture. It is possible to create a project without providing a description.",
						},
						"destroy_on_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The policy that indicates whether the resources are destroyed or not when a project is deleted.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmProjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if err = d.Set("crn", project.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(project.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	cumulativeNeedsAttentionView := []map[string]interface{}{}
	if project.CumulativeNeedsAttentionView != nil {
		for _, modelItem := range project.CumulativeNeedsAttentionView {
			modelMap, err := dataSourceIbmProjectCumulativeNeedsAttentionToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			cumulativeNeedsAttentionView = append(cumulativeNeedsAttentionView, modelMap)
		}
	}
	if err = d.Set("cumulative_needs_attention_view", cumulativeNeedsAttentionView); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cumulative_needs_attention_view %s", err))
	}

	if err = d.Set("cumulative_needs_attention_view_error", project.CumulativeNeedsAttentionViewError); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting cumulative_needs_attention_view_error: %s", err))
	}

	if err = d.Set("location", project.Location); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting location: %s", err))
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

	if err = d.Set("resource_group", project.ResourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
	}

	if err = d.Set("event_notifications_crn", project.EventNotificationsCrn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting event_notifications_crn: %s", err))
	}

	configs := []map[string]interface{}{}
	if project.Configs != nil {
		for _, modelItem := range project.Configs {
			modelMap, err := dataSourceIbmProjectProjectConfigSummaryToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			configs = append(configs, modelMap)
		}
	}
	if err = d.Set("configs", configs); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting configs %s", err))
	}

	environments := []map[string]interface{}{}
	if project.Environments != nil {
		for _, modelItem := range project.Environments {
			modelMap, err := dataSourceIbmProjectProjectEnvironmentSummaryToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			environments = append(environments, modelMap)
		}
	}
	if err = d.Set("environments", environments); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting environments %s", err))
	}

	definition := []map[string]interface{}{}
	if project.Definition != nil {
		modelMap, err := dataSourceIbmProjectProjectDefinitionPropertiesToMap(project.Definition)
		if err != nil {
			return diag.FromErr(err)
		}
		definition = append(definition, modelMap)
	}
	if err = d.Set("definition", definition); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definition %s", err))
	}

	return nil
}

func dataSourceIbmProjectCumulativeNeedsAttentionToMap(model *projectv1.CumulativeNeedsAttention) (map[string]interface{}, error) {
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

func dataSourceIbmProjectProjectConfigSummaryToMap(model *projectv1.ProjectConfigSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApprovedVersion != nil {
		approvedVersionMap, err := dataSourceIbmProjectProjectConfigVersionSummaryToMap(model.ApprovedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["approved_version"] = []map[string]interface{}{approvedVersionMap}
	}
	if model.DeployedVersion != nil {
		deployedVersionMap, err := dataSourceIbmProjectProjectConfigVersionSummaryToMap(model.DeployedVersion)
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
	definitionMap, err := dataSourceIbmProjectProjectConfigDefinitionNameDescriptionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	projectMap, err := dataSourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	if model.DeploymentModel != nil {
		modelMap["deployment_model"] = model.DeploymentModel
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectConfigVersionSummaryToMap(model *projectv1.ProjectConfigVersionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["state"] = model.State
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmProjectProjectConfigDefinitionNameDescriptionToMap(model *projectv1.ProjectConfigDefinitionNameDescription) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	definitionMap, err := dataSourceIbmProjectProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = model.Crn
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmProjectProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmProjectProjectEnvironmentSummaryToMap(model *projectv1.ProjectEnvironmentSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	projectMap, err := dataSourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	definitionMap, err := dataSourceIbmProjectEnvironmentDefinitionNameDescriptionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	return modelMap, nil
}

func dataSourceIbmProjectEnvironmentDefinitionNameDescriptionToMap(model *projectv1.EnvironmentDefinitionNameDescription) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = model.Name
	}
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	return modelMap, nil
}

func dataSourceIbmProjectProjectDefinitionPropertiesToMap(model *projectv1.ProjectDefinitionProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	if model.Description != nil {
		modelMap["description"] = model.Description
	}
	modelMap["destroy_on_delete"] = model.DestroyOnDelete
	return modelMap, nil
}
