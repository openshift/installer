// Copyright IBM Corp. 2022 All Rights Reserved.
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

func DataSourceIBMTektonPipelineDefinition() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMTektonPipelineDefinitionRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tekton pipeline ID.",
			},
			"definition_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The definition ID.",
			},
			"scm_source": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Scm source for tekton pipeline defintion.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "General href URL.",
						},
						"branch": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A branch of the repo, branch field doesn't coexist with tag field.",
						},
						"tag": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A tag of the repo.",
						},
						"path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path to the definitions yaml files.",
						},
					},
				},
			},
			"service_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID.",
			},
		},
	}
}

func DataSourceIBMTektonPipelineDefinitionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	scmSource := []map[string]interface{}{}
	if definition.ScmSource != nil {
		modelMap, err := DataSourceIBMTektonPipelineDefinitionDefinitionScmSourceToMap(definition.ScmSource)
		if err != nil {
			return diag.FromErr(err)
		}
		scmSource = append(scmSource, modelMap)
	}
	if err = d.Set("scm_source", scmSource); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting scm_source %s", err))
	}

	if err = d.Set("service_instance_id", definition.ServiceInstanceID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting service_instance_id: %s", err))
	}

	return nil
}

func DataSourceIBMTektonPipelineDefinitionDefinitionScmSourceToMap(model *cdtektonpipelinev2.DefinitionScmSource) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.URL != nil {
		modelMap["url"] = *model.URL
	}
	if model.Branch != nil {
		modelMap["branch"] = *model.Branch
	}
	if model.Tag != nil {
		modelMap["tag"] = *model.Tag
	}
	if model.Path != nil {
		modelMap["path"] = *model.Path
	}
	return modelMap, nil
}
