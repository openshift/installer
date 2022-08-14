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

func DataSourceIBMTektonPipelineProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIBMTektonPipelinePropertyRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The tekton pipeline ID.",
			},
			"property_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The property's name.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property name.",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String format property value.",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Options for SINGLE_SELECT property type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"default": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default option for SINGLE_SELECT property type.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property type.",
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "property path for INTEGRATION type properties.",
			},
		},
	}
}

func DataSourceIBMTektonPipelinePropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelinePropertyOptions := &cdtektonpipelinev2.GetTektonPipelinePropertyOptions{}

	getTektonPipelinePropertyOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelinePropertyOptions.SetPropertyName(d.Get("property_name").(string))

	property, response, err := cdTektonPipelineClient.GetTektonPipelinePropertyWithContext(context, getTektonPipelinePropertyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetTektonPipelinePropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelinePropertyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getTektonPipelinePropertyOptions.PipelineID, *getTektonPipelinePropertyOptions.PropertyName))

	if err = d.Set("name", property.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("value", property.Value); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
	}

	if err = d.Set("default", property.Default); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting default: %s", err))
	}

	if err = d.Set("type", property.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("path", property.Path); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting path: %s", err))
	}

	if property.Enum != nil {
		if err = d.Set("enum", property.Enum); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enum: %s", err))
		}
	}

	return nil
}
