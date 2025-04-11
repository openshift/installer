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
				Description: "The resource group name where the project's data and tools are created.",
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
							Description: "The name of the project.  It's unique within the account across regions.",
						},
						"destroy_on_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Required:    true,
							Description: "The policy that indicates whether the resources are destroyed or not when a project is deleted.",
						},
						"description": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "A brief explanation of the project's use in the configuration of a deployable architecture. You can create a project without providing a description.",
						},
						"auto_deploy": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "A boolean flag to enable auto deploy.",
						},
						"monitoring_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "A boolean flag to enable automatic drift detection. Use this field to run a daily check to compare your configurations to your deployed resources to detect any difference.",
						},
					},
				},
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An IBM Cloud resource name that uniquely identifies a resource.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
			},
			"cumulative_needs_attention_view": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The cumulative list of needs attention items for a project. If the view is successfully retrieved, an empty or nonempty array is returned.",
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
							Description: "A unique ID for this individual event.",
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
				Description: "A value of `true` indicates that the fetch of the needs attention items failed. This property only exists if there was an error when you retrieved the cumulative needs attention view.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group ID where the project's data and tools are created.",
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
				Description: "The CRN of the Event Notifications instance if one is connected to this project.",
			},
			"configs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The project configurations. These configurations are only included in the response of creating a project if a configuration array is specified in the request payload.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"approved_version": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A summary of a project configuration version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A summary of the definition in a project configuration version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the project environment.",
												},
												"locator_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													ForceNew:    true,
													Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
												},
											},
										},
									},
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The state of the configuration.",
									},
									"state_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Computed state code clarifying the prerequisites for validation for the configuration.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the configuration.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A URL.",
									},
								},
							},
						},
						"deployed_version": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A summary of a project configuration version.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"definition": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "A summary of the definition in a project configuration version.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"environment_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The ID of the project environment.",
												},
												"locator_id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													ForceNew:    true,
													Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
												},
											},
										},
									},
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The state of the configuration.",
									},
									"state_code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Computed state code clarifying the prerequisites for validation for the configuration.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The version number of the configuration.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A URL.",
									},
								},
							},
						},
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
							Description: "A date and time value in the format YYYY-MM-DDTHH:mm:ssZ or YYYY-MM-DDTHH:mm:ss.sssZ to match the date and time format as specified by RFC 3339.",
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
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The description of a project configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A project configuration description.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The configuration name. It's unique within the account across projects and regions.",
									},
									"locator_id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										ForceNew:    true,
										Description: "A unique concatenation of the catalog ID and the version ID that identify the deployable architecture in the catalog. I you're importing from an existing Schematics workspace that is not backed by cart, a `locator_id` is required. If you're using a Schematics workspace that is backed by cart, a `locator_id` is not necessary because the Schematics workspace has one.> There are 3 scenarios:> 1. If only a `locator_id` is specified, a new Schematics workspace is instantiated with that `locator_id`.> 2. If only a schematics `workspace_crn` is specified, a `400` is returned if a `locator_id` is not found in the existing schematics workspace.> 3. If both a Schematics `workspace_crn` and a `locator_id` are specified, a `400` message is returned if the specified `locator_id` does not agree with the `locator_id` in the existing Schematics workspace.> For more information of creating a Schematics workspace, see [Creating workspaces and importing your Terraform template](/docs/schematics?topic=schematics-sch-create-wks).",
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
				Description: "The project environment. These environments are only included in the response if project environments were created on the project.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The environment ID as a friendly name.",
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
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A URL.",
						},
						"definition": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The environment definition that is used in the project collection.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"description": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The description of the environment.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the environment. It's unique within the account across projects and regions.",
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
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "ca-tor, eu-de, eu-gb, us-east, us-south",
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
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createProjectOptions := &projectv1.CreateProjectOptions{}

	definitionModel, err := ResourceIbmProjectMapToProjectPrototypeDefinition(d.Get("definition.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "create", "parse-definition").GetDiag()
	}
	createProjectOptions.SetDefinition(definitionModel)
	createProjectOptions.SetLocation(d.Get("location").(string))
	createProjectOptions.SetResourceGroup(d.Get("resource_group").(string))
	if _, ok := d.GetOk("configs"); ok {
		var configs []projectv1.ProjectConfigPrototype
		for _, v := range d.Get("configs").([]interface{}) {
			value := v.(map[string]interface{})
			configsItem, err := ResourceIbmProjectMapToProjectConfigPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "create", "parse-configs").GetDiag()
			}
			configs = append(configs, *configsItem)
		}
		createProjectOptions.SetConfigs(configs)
	}
	if _, ok := d.GetOk("environments"); ok {
		var environments []projectv1.EnvironmentPrototype
		for _, v := range d.Get("environments").([]interface{}) {
			value := v.(map[string]interface{})
			environmentsItem, err := ResourceIbmProjectMapToEnvironmentPrototype(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "create", "parse-environments").GetDiag()
			}
			environments = append(environments, *environmentsItem)
		}
		createProjectOptions.SetEnvironments(environments)
	}

	project, _, err := projectClient.CreateProjectWithContext(context, createProjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateProjectWithContext failed: %s", err.Error()), "ibm_project", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*project.ID)

	return resourceIbmProjectRead(context, d, meta)
}

func resourceIbmProjectRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getProjectOptions := &projectv1.GetProjectOptions{}

	getProjectOptions.SetID(d.Id())

	project, response, err := projectClient.GetProjectWithContext(context, getProjectOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetProjectWithContext failed: %s", err.Error()), "ibm_project", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("location", project.Location); err != nil {
		err = fmt.Errorf("Error setting location: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-location").GetDiag()
	}
	if err = d.Set("resource_group", project.ResourceGroup); err != nil {
		err = fmt.Errorf("Error setting resource_group: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-resource_group").GetDiag()
	}
	definitionMap, err := ResourceIbmProjectProjectDefinitionPropertiesToMap(project.Definition)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "definition-to-map").GetDiag()
	}
	if err = d.Set("definition", []map[string]interface{}{definitionMap}); err != nil {
		err = fmt.Errorf("Error setting definition: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-definition").GetDiag()
	}
	if err = d.Set("crn", project.Crn); err != nil {
		err = fmt.Errorf("Error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-crn").GetDiag()
	}
	if err = d.Set("created_at", flex.DateTimeToString(project.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-created_at").GetDiag()
	}
	cumulativeNeedsAttentionView := []map[string]interface{}{}
	for _, cumulativeNeedsAttentionViewItem := range project.CumulativeNeedsAttentionView {
		cumulativeNeedsAttentionViewItemMap, err := ResourceIbmProjectCumulativeNeedsAttentionToMap(&cumulativeNeedsAttentionViewItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "cumulative_needs_attention_view-to-map").GetDiag()
		}
		cumulativeNeedsAttentionView = append(cumulativeNeedsAttentionView, cumulativeNeedsAttentionViewItemMap)
	}
	if err = d.Set("cumulative_needs_attention_view", cumulativeNeedsAttentionView); err != nil {
		err = fmt.Errorf("Error setting cumulative_needs_attention_view: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-cumulative_needs_attention_view").GetDiag()
	}
	if !core.IsNil(project.CumulativeNeedsAttentionViewError) {
		if err = d.Set("cumulative_needs_attention_view_error", project.CumulativeNeedsAttentionViewError); err != nil {
			err = fmt.Errorf("Error setting cumulative_needs_attention_view_error: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-cumulative_needs_attention_view_error").GetDiag()
		}
	}
	if err = d.Set("resource_group_id", project.ResourceGroupID); err != nil {
		err = fmt.Errorf("Error setting resource_group_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-resource_group_id").GetDiag()
	}
	if err = d.Set("state", project.State); err != nil {
		err = fmt.Errorf("Error setting state: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-state").GetDiag()
	}
	if err = d.Set("href", project.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-href").GetDiag()
	}
	if !core.IsNil(project.EventNotificationsCrn) {
		if err = d.Set("event_notifications_crn", project.EventNotificationsCrn); err != nil {
			err = fmt.Errorf("Error setting event_notifications_crn: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-event_notifications_crn").GetDiag()
		}
	}
	configs := []map[string]interface{}{}
	for _, configsItem := range project.Configs {
		configsItemMap, err := ResourceIbmProjectProjectConfigSummaryToMap(&configsItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "configs-to-map").GetDiag()
		}
		configs = append(configs, configsItemMap)
	}
	if err = d.Set("configs", configs); err != nil {
		err = fmt.Errorf("Error setting configs: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-configs").GetDiag()
	}
	environments := []map[string]interface{}{}
	for _, environmentsItem := range project.Environments {
		environmentsItemMap, err := ResourceIbmProjectProjectEnvironmentSummaryToMap(&environmentsItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "environments-to-map").GetDiag()
		}
		environments = append(environments, environmentsItemMap)
	}
	if err = d.Set("environments", environments); err != nil {
		err = fmt.Errorf("Error setting environments: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "read", "set-environments").GetDiag()
	}

	return nil
}

func resourceIbmProjectUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateProjectOptions := &projectv1.UpdateProjectOptions{}

	updateProjectOptions.SetID(d.Id())

	hasChange := false

	if d.HasChange("definition") {
		definition, err := ResourceIbmProjectMapToProjectPatchDefinitionBlock(d.Get("definition.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_project", "update", "parse-definition").GetDiag()
		}
		updateProjectOptions.SetDefinition(definition)
		hasChange = true
	}

	if hasChange {
		_, _, err = projectClient.UpdateProjectWithContext(context, updateProjectOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateProjectWithContext failed: %s", err.Error()), "ibm_project", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmProjectRead(context, d, meta)
}

func resourceIbmProjectDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectClient, err := meta.(conns.ClientSession).ProjectV1()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_project", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteProjectOptions := &projectv1.DeleteProjectOptions{}

	deleteProjectOptions.SetID(d.Id())

	_, _, err = projectClient.DeleteProjectWithContext(context, deleteProjectOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteProjectWithContext failed: %s", err.Error()), "ibm_project", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIbmProjectMapToProjectPrototypeDefinition(modelMap map[string]interface{}) (*projectv1.ProjectPrototypeDefinition, error) {
	model := &projectv1.ProjectPrototypeDefinition{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["destroy_on_delete"] != nil {
		model.DestroyOnDelete = core.BoolPtr(modelMap["destroy_on_delete"].(bool))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["auto_deploy"] != nil {
		model.AutoDeploy = core.BoolPtr(modelMap["auto_deploy"].(bool))
	}
	if modelMap["monitoring_enabled"] != nil {
		model.MonitoringEnabled = core.BoolPtr(modelMap["monitoring_enabled"].(bool))
	}
	return model, nil
}

func ResourceIbmProjectMapToProjectConfigPrototype(modelMap map[string]interface{}) (*projectv1.ProjectConfigPrototype, error) {
	model := &projectv1.ProjectConfigPrototype{}
	DefinitionModel, err := ResourceIbmProjectMapToProjectConfigDefinitionPrototype(modelMap["definition"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Definition = DefinitionModel
	if modelMap["schematics"] != nil && len(modelMap["schematics"].([]interface{})) > 0 {
		SchematicsModel, err := ResourceIbmProjectMapToSchematicsWorkspace(modelMap["schematics"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Schematics = SchematicsModel
	}
	return model, nil
}

func ResourceIbmProjectMapToProjectConfigDefinitionPrototype(modelMap map[string]interface{}) (projectv1.ProjectConfigDefinitionPrototypeIntf, error) {
	model := &projectv1.ProjectConfigDefinitionPrototype{}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := ResourceIbmProjectMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	if modelMap["locator_id"] != nil && modelMap["locator_id"].(string) != "" {
		model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["environment_id"] != nil && modelMap["environment_id"].(string) != "" {
		model.EnvironmentID = core.StringPtr(modelMap["environment_id"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := ResourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		model.Inputs = modelMap["inputs"].(map[string]interface{})
	}
	if modelMap["settings"] != nil {
		model.Settings = modelMap["settings"].(map[string]interface{})
	}
	if modelMap["resource_crns"] != nil {
		resourceCrns := []string{}
		for _, resourceCrnsItem := range modelMap["resource_crns"].([]interface{}) {
			resourceCrns = append(resourceCrns, resourceCrnsItem.(string))
		}
		model.ResourceCrns = resourceCrns
	}
	return model, nil
}

func ResourceIbmProjectMapToProjectComplianceProfile(modelMap map[string]interface{}) (*projectv1.ProjectComplianceProfile, error) {
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

func ResourceIbmProjectMapToProjectConfigAuth(modelMap map[string]interface{}) (*projectv1.ProjectConfigAuth, error) {
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

func ResourceIbmProjectMapToStackConfigMember(modelMap map[string]interface{}) (*projectv1.StackConfigMember, error) {
	model := &projectv1.StackConfigMember{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.ConfigID = core.StringPtr(modelMap["config_id"].(string))
	return model, nil
}

func ResourceIbmProjectMapToProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype(modelMap map[string]interface{}) (*projectv1.ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype, error) {
	model := &projectv1.ProjectConfigDefinitionPrototypeDAConfigDefinitionPropertiesPrototype{}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := ResourceIbmProjectMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	if modelMap["locator_id"] != nil && modelMap["locator_id"].(string) != "" {
		model.LocatorID = core.StringPtr(modelMap["locator_id"].(string))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["environment_id"] != nil && modelMap["environment_id"].(string) != "" {
		model.EnvironmentID = core.StringPtr(modelMap["environment_id"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := ResourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		model.Inputs = modelMap["inputs"].(map[string]interface{})
	}
	if modelMap["settings"] != nil {
		model.Settings = modelMap["settings"].(map[string]interface{})
	}
	return model, nil
}

func ResourceIbmProjectMapToProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype(modelMap map[string]interface{}) (*projectv1.ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype, error) {
	model := &projectv1.ProjectConfigDefinitionPrototypeResourceConfigDefinitionPropertiesPrototype{}
	if modelMap["resource_crns"] != nil {
		resourceCrns := []string{}
		for _, resourceCrnsItem := range modelMap["resource_crns"].([]interface{}) {
			resourceCrns = append(resourceCrns, resourceCrnsItem.(string))
		}
		model.ResourceCrns = resourceCrns
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["environment_id"] != nil && modelMap["environment_id"].(string) != "" {
		model.EnvironmentID = core.StringPtr(modelMap["environment_id"].(string))
	}
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := ResourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		model.Inputs = modelMap["inputs"].(map[string]interface{})
	}
	if modelMap["settings"] != nil {
		model.Settings = modelMap["settings"].(map[string]interface{})
	}
	return model, nil
}

func ResourceIbmProjectMapToSchematicsWorkspace(modelMap map[string]interface{}) (*projectv1.SchematicsWorkspace, error) {
	model := &projectv1.SchematicsWorkspace{}
	if modelMap["workspace_crn"] != nil && modelMap["workspace_crn"].(string) != "" {
		model.WorkspaceCrn = core.StringPtr(modelMap["workspace_crn"].(string))
	}
	return model, nil
}

func ResourceIbmProjectMapToEnvironmentPrototype(modelMap map[string]interface{}) (*projectv1.EnvironmentPrototype, error) {
	model := &projectv1.EnvironmentPrototype{}
	DefinitionModel, err := ResourceIbmProjectMapToEnvironmentDefinitionRequiredProperties(modelMap["definition"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Definition = DefinitionModel
	return model, nil
}

func ResourceIbmProjectMapToEnvironmentDefinitionRequiredProperties(modelMap map[string]interface{}) (*projectv1.EnvironmentDefinitionRequiredProperties, error) {
	model := &projectv1.EnvironmentDefinitionRequiredProperties{}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	model.Name = core.StringPtr(modelMap["name"].(string))
	if modelMap["authorizations"] != nil && len(modelMap["authorizations"].([]interface{})) > 0 {
		AuthorizationsModel, err := ResourceIbmProjectMapToProjectConfigAuth(modelMap["authorizations"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Authorizations = AuthorizationsModel
	}
	if modelMap["inputs"] != nil {
		model.Inputs = modelMap["inputs"].(map[string]interface{})
	}
	if modelMap["compliance_profile"] != nil && len(modelMap["compliance_profile"].([]interface{})) > 0 {
		ComplianceProfileModel, err := ResourceIbmProjectMapToProjectComplianceProfile(modelMap["compliance_profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ComplianceProfile = ComplianceProfileModel
	}
	return model, nil
}

func ResourceIbmProjectMapToProjectPatchDefinitionBlock(modelMap map[string]interface{}) (*projectv1.ProjectPatchDefinitionBlock, error) {
	model := &projectv1.ProjectPatchDefinitionBlock{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["destroy_on_delete"] != nil {
		model.DestroyOnDelete = core.BoolPtr(modelMap["destroy_on_delete"].(bool))
	}
	if modelMap["auto_deploy"] != nil {
		model.AutoDeploy = core.BoolPtr(modelMap["auto_deploy"].(bool))
	}
	if modelMap["description"] != nil && modelMap["description"].(string) != "" {
		model.Description = core.StringPtr(modelMap["description"].(string))
	}
	if modelMap["monitoring_enabled"] != nil {
		model.MonitoringEnabled = core.BoolPtr(modelMap["monitoring_enabled"].(bool))
	}
	return model, nil
}

func ResourceIbmProjectProjectDefinitionPropertiesToMap(model *projectv1.ProjectDefinitionProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["destroy_on_delete"] = *model.DestroyOnDelete
	modelMap["description"] = *model.Description
	if model.AutoDeploy != nil {
		modelMap["auto_deploy"] = *model.AutoDeploy
	}
	if model.MonitoringEnabled != nil {
		modelMap["monitoring_enabled"] = *model.MonitoringEnabled
	}
	return modelMap, nil
}

func ResourceIbmProjectCumulativeNeedsAttentionToMap(model *projectv1.CumulativeNeedsAttention) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Event != nil {
		modelMap["event"] = *model.Event
	}
	if model.EventID != nil {
		modelMap["event_id"] = *model.EventID
	}
	if model.ConfigID != nil {
		modelMap["config_id"] = *model.ConfigID
	}
	if model.ConfigVersion != nil {
		modelMap["config_version"] = flex.IntValue(model.ConfigVersion)
	}
	return modelMap, nil
}

func ResourceIbmProjectProjectConfigSummaryToMap(model *projectv1.ProjectConfigSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ApprovedVersion != nil {
		approvedVersionMap, err := ResourceIbmProjectProjectConfigVersionSummaryToMap(model.ApprovedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["approved_version"] = []map[string]interface{}{approvedVersionMap}
	}
	if model.DeployedVersion != nil {
		deployedVersionMap, err := ResourceIbmProjectProjectConfigVersionSummaryToMap(model.DeployedVersion)
		if err != nil {
			return modelMap, err
		}
		modelMap["deployed_version"] = []map[string]interface{}{deployedVersionMap}
	}
	modelMap["id"] = *model.ID
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["state"] = *model.State
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["modified_at"] = model.ModifiedAt.String()
	modelMap["href"] = *model.Href
	definitionMap, err := ResourceIbmProjectProjectConfigSummaryDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	projectMap, err := ResourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	if model.DeploymentModel != nil {
		modelMap["deployment_model"] = *model.DeploymentModel
	}
	return modelMap, nil
}

func ResourceIbmProjectProjectConfigVersionSummaryToMap(model *projectv1.ProjectConfigVersionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	definitionMap, err := ResourceIbmProjectProjectConfigVersionDefinitionSummaryToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["state"] = *model.State
	if model.StateCode != nil {
		modelMap["state_code"] = *model.StateCode
	}
	modelMap["version"] = flex.IntValue(model.Version)
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func ResourceIbmProjectProjectConfigVersionDefinitionSummaryToMap(model *projectv1.ProjectConfigVersionDefinitionSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EnvironmentID != nil {
		modelMap["environment_id"] = *model.EnvironmentID
	}
	if model.LocatorID != nil {
		modelMap["locator_id"] = *model.LocatorID
	}
	return modelMap, nil
}

func ResourceIbmProjectProjectConfigSummaryDefinitionToMap(model *projectv1.ProjectConfigSummaryDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["description"] = *model.Description
	modelMap["name"] = *model.Name
	if model.LocatorID != nil {
		modelMap["locator_id"] = *model.LocatorID
	}
	return modelMap, nil
}

func ResourceIbmProjectProjectReferenceToMap(model *projectv1.ProjectReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["href"] = *model.Href
	definitionMap, err := ResourceIbmProjectProjectDefinitionReferenceToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	modelMap["crn"] = *model.Crn
	return modelMap, nil
}

func ResourceIbmProjectProjectDefinitionReferenceToMap(model *projectv1.ProjectDefinitionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIbmProjectProjectEnvironmentSummaryToMap(model *projectv1.ProjectEnvironmentSummary) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	projectMap, err := ResourceIbmProjectProjectReferenceToMap(model.Project)
	if err != nil {
		return modelMap, err
	}
	modelMap["project"] = []map[string]interface{}{projectMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = *model.Href
	definitionMap, err := ResourceIbmProjectProjectEnvironmentSummaryDefinitionToMap(model.Definition)
	if err != nil {
		return modelMap, err
	}
	modelMap["definition"] = []map[string]interface{}{definitionMap}
	return modelMap, nil
}

func ResourceIbmProjectProjectEnvironmentSummaryDefinitionToMap(model *projectv1.ProjectEnvironmentSummaryDefinition) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["description"] = *model.Description
	modelMap["name"] = *model.Name
	return modelMap, nil
}
