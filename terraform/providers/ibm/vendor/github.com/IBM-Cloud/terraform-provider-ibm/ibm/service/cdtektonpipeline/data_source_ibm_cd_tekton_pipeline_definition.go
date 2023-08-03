// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtektonpipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
)

func DataSourceIBMCdTektonPipelineDefinition() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdTektonPipelineDefinitionRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Tekton pipeline ID.",
			},
			"definition_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The definition ID.",
			},
			"source": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Source repository containing the Tekton pipeline definition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The only supported source type is \"git\", indicating that the source is a git repository.",
						},
						"properties": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Properties of the source, which define the URL of the repository and a branch or tag.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "URL of the definition repository.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A branch from the repo, specify one of branch or tag only.",
									},
									"tag": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A tag from the repo, specify one of branch or tag only.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The path to the definition's YAML files.",
									},
									"tool": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Reference to the repository tool in the parent toolchain.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "ID of the repository tool instance in the parent toolchain.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the definition.",
			},
		},
	}
}

func dataSourceIBMCdTektonPipelineDefinitionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineDefinitionOptions := &cdtektonpipelinev2.GetTektonPipelineDefinitionOptions{}

	getTektonPipelineDefinitionOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelineDefinitionOptions.SetDefinitionID(d.Get("definition_id").(string))

	definition, response, err := cdTektonPipelineClient.GetTektonPipelineDefinitionWithContext(context, getTektonPipelineDefinitionOptions)
	if err != nil {
		log.Printf("[DEBUG] GetTektonPipelineDefinitionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineDefinitionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getTektonPipelineDefinitionOptions.PipelineID, *getTektonPipelineDefinitionOptions.DefinitionID))

	source := []map[string]interface{}{}
	if definition.Source != nil {
		modelMap, err := dataSourceIBMCdTektonPipelineDefinitionDefinitionSourceToMap(definition.Source)
		if err != nil {
			return diag.FromErr(err)
		}
		source = append(source, modelMap)
	}
	if err = d.Set("source", source); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source %s", err))
	}

	if err = d.Set("href", definition.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	return nil
}

func dataSourceIBMCdTektonPipelineDefinitionDefinitionSourceToMap(model *cdtektonpipelinev2.DefinitionSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	propertiesMap, err := dataSourceIBMCdTektonPipelineDefinitionDefinitionSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func dataSourceIBMCdTektonPipelineDefinitionDefinitionSourcePropertiesToMap(model *cdtektonpipelinev2.DefinitionSourceProperties) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["url"] = model.URL
	if model.Branch != nil {
		modelMap["branch"] = model.Branch
	}
	if model.Tag != nil {
		modelMap["tag"] = model.Tag
	}
	modelMap["path"] = model.Path
	if model.Tool != nil {
		toolMap, err := dataSourceIBMCdTektonPipelineDefinitionToolToMap(model.Tool)
		if err != nil {
			return modelMap, err
		}
		modelMap["tool"] = []map[string]interface{}{toolMap}
	}
	return modelMap, nil
}

func dataSourceIBMCdTektonPipelineDefinitionToolToMap(model *cdtektonpipelinev2.Tool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}
