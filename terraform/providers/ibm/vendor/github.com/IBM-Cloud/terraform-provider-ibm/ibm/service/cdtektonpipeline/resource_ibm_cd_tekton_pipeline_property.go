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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/v2/cdtektonpipelinev2"
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
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "type"),
				Description:  "Property type.",
			},
			"locked": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this property cannot be overridden by a trigger property or at runtime. Attempting to override it will result in run requests being rejected. The default is false.",
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "path"),
				Description:  "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
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
			Identifier:                 "value",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^.*$`,
			MinValueLength:             0,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "appconfig, integration, secure, single_select, text",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
	if _, ok := d.GetOk("locked"); ok {
		createTektonPipelinePropertiesOptions.SetLocked(d.Get("locked").(bool))
	}
	if _, ok := d.GetOk("path"); ok {
		createTektonPipelinePropertiesOptions.SetPath(d.Get("path").(string))
	}

	property, _, err := cdTektonPipelineClient.CreateTektonPipelinePropertiesWithContext(context, createTektonPipelinePropertiesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTektonPipelinePropertiesWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_property", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createTektonPipelinePropertiesOptions.PipelineID, *property.Name))

	return resourceIBMCdTektonPipelinePropertyRead(context, d, meta)
}

func resourceIBMCdTektonPipelinePropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTektonPipelinePropertyOptions := &cdtektonpipelinev2.GetTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "sep-id-parts").GetDiag()
	}

	getTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	getTektonPipelinePropertyOptions.SetPropertyName(parts[1])

	property, response, err := cdTektonPipelineClient.GetTektonPipelinePropertyWithContext(context, getTektonPipelinePropertyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTektonPipelinePropertyWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_property", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", property.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-name").GetDiag()
	}
	if !core.IsNil(property.Value) {
		if err = d.Set("value", property.Value); err != nil {
			err = fmt.Errorf("Error setting value: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-value").GetDiag()
		}
	}
	if !core.IsNil(property.Enum) {
		if err = d.Set("enum", property.Enum); err != nil {
			err = fmt.Errorf("Error setting enum: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-enum").GetDiag()
		}
	}
	if err = d.Set("type", property.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-type").GetDiag()
	}
	if !core.IsNil(property.Locked) {
		if err = d.Set("locked", property.Locked); err != nil {
			err = fmt.Errorf("Error setting locked: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-locked").GetDiag()
		}
	}
	if !core.IsNil(property.Path) {
		if err = d.Set("path", property.Path); err != nil {
			err = fmt.Errorf("Error setting path: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-path").GetDiag()
		}
	}
	if !core.IsNil(property.Href) {
		if err = d.Set("href", property.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "read", "set-href").GetDiag()
		}
	}

	return nil
}

func resourceIBMCdTektonPipelinePropertyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	replaceTektonPipelinePropertyOptions := &cdtektonpipelinev2.ReplaceTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "update", "sep-id-parts").GetDiag()
	}

	replaceTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	replaceTektonPipelinePropertyOptions.SetPropertyName(parts[1])
	replaceTektonPipelinePropertyOptions.SetName(d.Get("name").(string))
	replaceTektonPipelinePropertyOptions.SetType(d.Get("type").(string))

	hasChange := false

	if d.HasChange("pipeline_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_property", "update", "pipeline_id-forces-new").GetDiag()
	}
	if d.HasChange("name") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "name")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_property", "update", "name-forces-new").GetDiag()
	}
	if d.HasChange("type") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "type")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_property", "update", "type-forces-new").GetDiag()
	}
	if d.HasChange("locked") {
		replaceTektonPipelinePropertyOptions.SetLocked(d.Get("locked").(bool))
		hasChange = true
	}
	if d.Get("type").(string) == "integration" {
		if d.HasChange("value") || d.HasChange("path") || d.HasChange("locked") {
			replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
			replaceTektonPipelinePropertyOptions.SetPath(d.Get("path").(string))
			hasChange = true
		}
	} else if d.Get("type").(string) == "single_select" {
		if d.HasChange("enum") || d.HasChange("value") || d.HasChange("locked") {
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
		if d.HasChange("value") || d.HasChange("locked") {
			replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
			hasChange = true
		}
	}

	if hasChange {
		_, _, err = cdTektonPipelineClient.ReplaceTektonPipelinePropertyWithContext(context, replaceTektonPipelinePropertyOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceTektonPipelinePropertyWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_property", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdTektonPipelinePropertyRead(context, d, meta)
}

func resourceIBMCdTektonPipelinePropertyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteTektonPipelinePropertyOptions := &cdtektonpipelinev2.DeleteTektonPipelinePropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_property", "delete", "sep-id-parts").GetDiag()
	}

	deleteTektonPipelinePropertyOptions.SetPipelineID(parts[0])
	deleteTektonPipelinePropertyOptions.SetPropertyName(parts[1])

	_, err = cdTektonPipelineClient.DeleteTektonPipelinePropertyWithContext(context, deleteTektonPipelinePropertyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTektonPipelinePropertyWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_property", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
