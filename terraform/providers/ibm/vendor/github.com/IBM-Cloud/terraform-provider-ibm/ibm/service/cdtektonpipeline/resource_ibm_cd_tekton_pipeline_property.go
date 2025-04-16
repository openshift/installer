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

func ResourceIBMCdTektonPipelineProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdTektonPipelinePropertyCreate,
		ReadContext:   resourceIBMCdTektonPipelinePropertyRead,
		UpdateContext: resourceIBMCdTektonPipelinePropertyUpdate,
		DeleteContext: resourceIBMCdTektonPipelinePropertyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "pipeline_id"),
				Description:  "The Tekton pipeline ID.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "name"),
				Description:  "Property name.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "type"),
				Description:  "Property type.",
			},
			"value": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.SuppressPipelinePropertyRawSecret,
				ValidateFunc:     validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "value"),
				Description:      "Property value. Any string value is valid.",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Options for `single_select` property type. Only needed when using `single_select` property type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "path"),
				Description:  "A dot notation path for `integration` type properties only, to select a value from the tool integration. If left blank the full tool integration data will be used.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the property.",
			},
		},
	}
}

func ResourceIBMCdTektonPipelinePropertyValidator() *validate.ResourceValidator {
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
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-zA-Z_.]{1,253}$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "appconfig, integration, secure, single_select, text",
		},
		validate.ValidateSchema{
			Identifier:                 "value",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             0,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "path",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-zA-Z_.]*$`,
			MinValueLength:             0,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_property", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdTektonPipelinePropertyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelinePropertiesOptions := &cdtektonpipelinev2.CreateTektonPipelinePropertiesOptions{}

	createTektonPipelinePropertiesOptions.SetPipelineID(d.Get("pipeline_id").(string))
	createTektonPipelinePropertiesOptions.SetName(d.Get("name").(string))
	createTektonPipelinePropertiesOptions.SetType(d.Get("type").(string))
	if _, ok := d.GetOk("value"); ok {
		createTektonPipelinePropertiesOptions.SetValue(d.Get("value").(string))
	}
	if _, ok := d.GetOk("enum"); ok {
		var enum []string
		for _, v := range d.Get("enum").([]interface{}) {
			enumItem := v.(string)
			enum = append(enum, enumItem)
		}
		createTektonPipelinePropertiesOptions.SetEnum(enum)
	}
	if _, ok := d.GetOk("path"); ok {
		createTektonPipelinePropertiesOptions.SetPath(d.Get("path").(string))
	}

	property, response, err := cdTektonPipelineClient.CreateTektonPipelinePropertiesWithContext(context, createTektonPipelinePropertiesOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateTektonPipelinePropertiesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateTektonPipelinePropertiesWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelinePropertiesOptions.PipelineID, *property.Name))

	return resourceIBMCdTektonPipelinePropertyRead(context, d, meta)
}

func resourceIBMCdTektonPipelinePropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelinePropertyOptions := &cdtektonpipelinev2.GetTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	getTektonPipelinePropertyOptions.SetPropertyName(parts[1])

	property, response, err := cdTektonPipelineClient.GetTektonPipelinePropertyWithContext(context, getTektonPipelinePropertyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetTektonPipelinePropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetTektonPipelinePropertyWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("pipeline_id", getTektonPipelinePropertyOptions.PipelineID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting pipeline_id: %s", err))
	}
	if err = d.Set("name", property.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("type", property.Type); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
	}
	if !core.IsNil(property.Value) {
		if err = d.Set("value", property.Value); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
		}
	}
	if !core.IsNil(property.Enum) {
		if err = d.Set("enum", property.Enum); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enum: %s", err))
		}
	}
	if !core.IsNil(property.Path) {
		if err = d.Set("path", property.Path); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting path: %s", err))
		}
	}
	if !core.IsNil(property.Href) {
		if err = d.Set("href", property.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}

	return nil
}

func resourceIBMCdTektonPipelinePropertyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelinePropertyOptions := &cdtektonpipelinev2.ReplaceTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	replaceTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	replaceTektonPipelinePropertyOptions.SetPropertyName(parts[1])
	replaceTektonPipelinePropertyOptions.SetName(d.Get("name").(string))
	replaceTektonPipelinePropertyOptions.SetType(d.Get("type").(string))

	hasChange := false

	if d.HasChange("pipeline_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id"))
	}
	if d.HasChange("name") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "name"))
	}
	if d.HasChange("type") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "type"))
	}

	if d.Get("type").(string) == "integration" {
		if d.HasChange("value") || d.HasChange("path") {
			replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
			replaceTektonPipelinePropertyOptions.SetPath(d.Get("path").(string))
			hasChange = true
		}
	} else if d.Get("type").(string) == "single_select" {
		if d.HasChange("enum") || d.HasChange("value") {
			var enum []string
			for _, v := range d.Get("enum").([]interface{}) {
				enumItem := v.(string)
				enum = append(enum, enumItem)
			}
			replaceTektonPipelinePropertyOptions.SetEnum(enum)
			replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
			hasChange = true
		}
	} else {
		if d.HasChange("value") {
			replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
			hasChange = true
		}
	}

	if hasChange {
		_, response, err := cdTektonPipelineClient.ReplaceTektonPipelinePropertyWithContext(context, replaceTektonPipelinePropertyOptions)
		if err != nil {
			log.Printf("[DEBUG] ReplaceTektonPipelinePropertyWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ReplaceTektonPipelinePropertyWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMCdTektonPipelinePropertyRead(context, d, meta)
}

func resourceIBMCdTektonPipelinePropertyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelinePropertyOptions := &cdtektonpipelinev2.DeleteTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	deleteTektonPipelinePropertyOptions.SetPropertyName(parts[1])

	response, err := cdTektonPipelineClient.DeleteTektonPipelinePropertyWithContext(context, deleteTektonPipelinePropertyOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteTektonPipelinePropertyWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteTektonPipelinePropertyWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}
