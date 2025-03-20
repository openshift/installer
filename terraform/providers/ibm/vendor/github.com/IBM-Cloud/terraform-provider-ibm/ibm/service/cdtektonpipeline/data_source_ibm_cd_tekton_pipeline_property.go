// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.95.2-120e65bc-20240924-152329
 */

package cdtektonpipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtektonpipelinev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func DataSourceIBMCdTektonPipelineProperty() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCdTektonPipelinePropertyRead,

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Tekton pipeline ID.",
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
				Description: "API URL for interacting with the property.",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Options for `single_select` property type. Only needed when using `single_select` property type.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Property type.",
			},
			"locked": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "When true, this property cannot be overridden by a trigger property or at runtime. Attempting to override it will result in run requests being rejected. The default is false.",
			},
			"path": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
			},
		},
	}
}

func dataSourceIBMCdTektonPipelinePropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cd_tekton_pipeline_property", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTektonPipelinePropertyOptions := &cdtektonpipelinev2.GetTektonPipelinePropertyOptions{}

	getTektonPipelinePropertyOptions.SetPipelineID(d.Get("pipeline_id").(string))
	getTektonPipelinePropertyOptions.SetPropertyName(d.Get("property_name").(string))

	property, _, err := cdTektonPipelineClient.GetTektonPipelinePropertyWithContext(context, getTektonPipelinePropertyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTektonPipelinePropertyWithContext failed: %s", err.Error()), "(Data) ibm_cd_tekton_pipeline_property", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getTektonPipelinePropertyOptions.PipelineID, *getTektonPipelinePropertyOptions.PropertyName))

	if err = d.Set("name", property.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-name").GetDiag()
	}

	if !core.IsNil(property.Value) {
		if err = d.Set("value", property.Value); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting value: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-value").GetDiag()
		}
	}

	if !core.IsNil(property.Href) {
		if err = d.Set("href", property.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-href").GetDiag()
		}
	}

	if !core.IsNil(property.Enum) {
		enum := []interface{}{}
		for _, enumItem := range property.Enum {
			enum = append(enum, enumItem)
		}
		if err = d.Set("enum", enum); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting enum: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-enum").GetDiag()
		}
	}

	if err = d.Set("type", property.Type); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting type: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-type").GetDiag()
	}

	if !core.IsNil(property.Locked) {
		if err = d.Set("locked", property.Locked); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting locked: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-locked").GetDiag()
		}
	}

	if !core.IsNil(property.Path) {
		if err = d.Set("path", property.Path); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting path: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-path").GetDiag()
		}
	}

	if !core.IsNil(property.Enum) {
		if err = d.Set("enum", property.Enum); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting enum: %s", err), "(Data) ibm_cd_tekton_pipeline_property", "read", "set-enum").GetDiag()
		}
	}

	return nil
}
