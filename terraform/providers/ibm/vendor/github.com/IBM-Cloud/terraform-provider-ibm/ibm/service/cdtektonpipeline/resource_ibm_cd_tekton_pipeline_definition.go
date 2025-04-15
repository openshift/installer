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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIBMCdTektonPipelineDefinition() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdTektonPipelineDefinitionCreate,
		ReadContext:   resourceIBMCdTektonPipelineDefinitionRead,
		UpdateContext: resourceIBMCdTektonPipelineDefinitionUpdate,
		DeleteContext: resourceIBMCdTektonPipelineDefinitionDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_definition", "pipeline_id"),
				Description:  "The Tekton pipeline ID.",
			},
			"source": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "Source repository containing the Tekton pipeline definition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The only supported source type is \"git\", indicating that the source is a git repository.",
						},
						"properties": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "Properties of the source, which define the URL of the repository and a branch or tag.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "URL of the definition repository.",
									},
									"branch": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A branch from the repo, specify one of branch or tag only.",
									},
									"tag": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "A tag from the repo, specify one of branch or tag only.",
									},
									"path": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
			"definition_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The aggregated definition ID.",
			},
		},
	}
}

func ResourceIBMCdTektonPipelineDefinitionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "pipeline_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z]+$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_definition", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdTektonPipelineDefinitionCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelineDefinitionOptions := &cdtektonpipelinev2.CreateTektonPipelineDefinitionOptions{}

	createTektonPipelineDefinitionOptions.SetPipelineID(d.Get("pipeline_id").(string))
	sourceModel, err := resourceIBMCdTektonPipelineDefinitionMapToDefinitionSource(d.Get("source.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createTektonPipelineDefinitionOptions.SetSource(sourceModel)

	definition, response, err := cdTektonPipelineClient.CreateTektonPipelineDefinitionWithContext(context, createTektonPipelineDefinitionOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelineDefinitionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelineDefinitionWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelineDefinitionOptions.PipelineID, *definition.ID))

	return resourceIBMCdTektonPipelineDefinitionRead(context, d, meta)
}

func resourceIBMCdTektonPipelineDefinitionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineDefinitionOptions := &cdtektonpipelinev2.GetTektonPipelineDefinitionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineDefinitionOptions.SetPipelineID(parts[0])
	getTektonPipelineDefinitionOptions.SetDefinitionID(parts[1])

	definition, response, err := cdTektonPipelineClient.GetTektonPipelineDefinitionWithContext(context, getTektonPipelineDefinitionOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelineDefinitionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineDefinitionWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("pipeline_id", getTektonPipelineDefinitionOptions.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}
	sourceMap, err := resourceIBMCdTektonPipelineDefinitionDefinitionSourceToMap(definition.Source)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("source", []map[string]interface{}{sourceMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting source: %s", err))
	}
	if !core.IsNil(definition.Href) {
		if err = d.Set("href", definition.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if err = d.Set("definition_id", definition.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting definition_id: %s", err))
	}

	return nil
}

func resourceIBMCdTektonPipelineDefinitionUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelineDefinitionOptions := &cdtektonpipelinev2.ReplaceTektonPipelineDefinitionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelineDefinitionOptions.SetPipelineID(parts[0])
	replaceTektonPipelineDefinitionOptions.SetDefinitionID(parts[1])

	hasChange := false

	if d.HasChange("pipeline_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id"))
	}
	if d.HasChange("source") {
		source, err := resourceIBMCdTektonPipelineDefinitionMapToDefinitionSource(d.Get("source.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		replaceTektonPipelineDefinitionOptions.SetSource(source)
		hasChange = true
	}

	if hasChange {
		_, response, err := cdTektonPipelineClient.ReplaceTektonPipelineDefinitionWithContext(context, replaceTektonPipelineDefinitionOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceTektonPipelineDefinitionWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceTektonPipelineDefinitionWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdTektonPipelineDefinitionRead(context, d, meta)
}

func resourceIBMCdTektonPipelineDefinitionDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineDefinitionOptions := &cdtektonpipelinev2.DeleteTektonPipelineDefinitionOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelineDefinitionOptions.SetPipelineID(parts[0])
	deleteTektonPipelineDefinitionOptions.SetDefinitionID(parts[1])

	response, err := cdTektonPipelineClient.DeleteTektonPipelineDefinitionWithContext(context, deleteTektonPipelineDefinitionOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelineDefinitionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelineDefinitionWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMCdTektonPipelineDefinitionMapToDefinitionSource(modelMap map[string]interface{}) (*cdtektonpipelinev2.DefinitionSource, error) {
	model := &cdtektonpipelinev2.DefinitionSource{}
	model.Type = core.StringPtr(modelMap["type"].(string))
	PropertiesModel, err := resourceIBMCdTektonPipelineDefinitionMapToDefinitionSourceProperties(modelMap["properties"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Properties = PropertiesModel
	return model, nil
}

func resourceIBMCdTektonPipelineDefinitionMapToDefinitionSourceProperties(modelMap map[string]interface{}) (*cdtektonpipelinev2.DefinitionSourceProperties, error) {
	model := &cdtektonpipelinev2.DefinitionSourceProperties{}
	model.URL = core.StringPtr(modelMap["url"].(string))
	if modelMap["branch"] != nil && modelMap["branch"].(string) != "" {
		model.Branch = core.StringPtr(modelMap["branch"].(string))
	}
	if modelMap["tag"] != nil && modelMap["tag"].(string) != "" {
		model.Tag = core.StringPtr(modelMap["tag"].(string))
	}
	model.Path = core.StringPtr(modelMap["path"].(string))
	if modelMap["tool"] != nil && len(modelMap["tool"].([]interface{})) > 0 {
		ToolModel, err := resourceIBMCdTektonPipelineDefinitionMapToTool(modelMap["tool"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Tool = ToolModel
	}
	return model, nil
}

func resourceIBMCdTektonPipelineDefinitionMapToTool(modelMap map[string]interface{}) (*cdtektonpipelinev2.Tool, error) {
	model := &cdtektonpipelinev2.Tool{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMCdTektonPipelineDefinitionDefinitionSourceToMap(model *cdtektonpipelinev2.DefinitionSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["type"] = model.Type
	propertiesMap, err := resourceIBMCdTektonPipelineDefinitionDefinitionSourcePropertiesToMap(model.Properties)
	if err != nil {
		return modelMap, err
	}
	modelMap["properties"] = []map[string]interface{}{propertiesMap}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineDefinitionDefinitionSourcePropertiesToMap(model *cdtektonpipelinev2.DefinitionSourceProperties) (map[string]interface{}, error) {
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
		toolMap, err := resourceIBMCdTektonPipelineDefinitionToolToMap(model.Tool)
		if err != nil {
			return modelMap, err
		}
		modelMap["tool"] = []map[string]interface{}{toolMap}
	}
	return modelMap, nil
}

func resourceIBMCdTektonPipelineDefinitionToolToMap(model *cdtektonpipelinev2.Tool) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	return modelMap, nil
}
