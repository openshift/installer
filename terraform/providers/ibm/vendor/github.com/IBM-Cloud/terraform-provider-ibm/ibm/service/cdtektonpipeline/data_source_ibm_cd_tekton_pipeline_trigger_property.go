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

func DataSourceIBMCdTektonPipelineTriggerProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdTektonPipelineTriggerPropertyRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Tekton pipeline ID.",
			},
			"trigger_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The trigger ID.",
			},
			"property_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The property name.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property name.",
			},
			"value": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property value. Any string value is valid.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the trigger property.",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Options for `single_select` property type. Only needed for `single_select` property type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property type.",
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
			},
		},
	}
}

func dataSourceIBMCdTektonPipelineTriggerPropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelineTriggerPropertyOptions := &cdtektonpipelinev2.GetTektonPipelineTriggerPropertyOptions{}

	getTektonPipelineTriggerPropertyOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelineTriggerPropertyOptions.SetTriggerID(d.Get("trigger_id").(string))
	getTektonPipelineTriggerPropertyOptions.SetPropertyName(d.Get("property_name").(string))

	triggerProperty, response, err := cdTektonPipelineClient.GetTektonPipelineTriggerPropertyWithContext(context, getTektonPipelineTriggerPropertyOptions)
	if err != nil {
		log.Printf("[DEBUG] GetTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelineTriggerPropertyWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *getTektonPipelineTriggerPropertyOptions.PipelineID, *getTektonPipelineTriggerPropertyOptions.TriggerID, *getTektonPipelineTriggerPropertyOptions.PropertyName))

	if err = d.Set("name", triggerProperty.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("value", triggerProperty.Value); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
	}

	if err = d.Set("href", triggerProperty.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("type", triggerProperty.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}

	if err = d.Set("path", triggerProperty.Path); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting path: %s", err))
	}

	if triggerProperty.Enum != nil {
		if err = d.Set("enum", triggerProperty.Enum); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enum: %s", err))
		}
	}

	return nil
}
