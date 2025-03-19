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

func ResourceIBMCdTektonPipelineTriggerProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCdTektonPipelineTriggerPropertyCreate,
		ReadContext:   resourceIBMCdTektonPipelineTriggerPropertyRead,
		UpdateContext: resourceIBMCdTektonPipelineTriggerPropertyUpdate,
		DeleteContext: resourceIBMCdTektonPipelineTriggerPropertyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger_property", "pipeline_id"),
				Description:  "The Tekton pipeline ID.",
			},
			"trigger_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger_property", "trigger_id"),
				Description:  "The trigger ID.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger_property", "name"),
				Description:  "Property name.",
			},
			"value": &schema.Schema{
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: flex.SuppressTriggerPropertyRawSecret,
				ValidateFunc:     validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger_property", "value"),
				Description:      "Property value. Any string value is valid.",
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Options for `single_select` property type. Only needed for `single_select` property type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger_property", "type"),
				Description:  "Property type.",
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_trigger_property", "path"),
				Description:  "A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left blank the full tool integration data will be used.",
			},
			"locked": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "When true, this property cannot be overridden at runtime. Attempting to override it will result in run requests being rejected. The default is false.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "API URL for interacting with the trigger property.",
			},
		},
	}
}

func ResourceIBMCdTektonPipelineTriggerPropertyValidator() *validate.ResourceValidator {
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
			Identifier:                 "trigger_id",
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_trigger_property", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMCdTektonPipelineTriggerPropertyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createTektonPipelineTriggerPropertiesOptions := &cdtektonpipelinev2.CreateTektonPipelineTriggerPropertiesOptions{}

	createTektonPipelineTriggerPropertiesOptions.SetPipelineID(d.Get("pipeline_id").(string))
	createTektonPipelineTriggerPropertiesOptions.SetTriggerID(d.Get("trigger_id").(string))
	createTektonPipelineTriggerPropertiesOptions.SetName(d.Get("name").(string))
	createTektonPipelineTriggerPropertiesOptions.SetType(d.Get("type").(string))
	if _, ok := d.GetOk("value"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetValue(d.Get("value").(string))
	}
	if _, ok := d.GetOk("enum"); ok {
		var enum []string
		for _, v := range d.Get("enum").([]interface{}) {
			enumItem := v.(string)
			enum = append(enum, enumItem)
		}
		createTektonPipelineTriggerPropertiesOptions.SetEnum(enum)
	}
	if _, ok := d.GetOk("path"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetPath(d.Get("path").(string))
	}
	if _, ok := d.GetOk("locked"); ok {
		createTektonPipelineTriggerPropertiesOptions.SetLocked(d.Get("locked").(bool))
	}

	triggerProperty, _, err := cdTektonPipelineClient.CreateTektonPipelineTriggerPropertiesWithContext(context, createTektonPipelineTriggerPropertiesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateTektonPipelineTriggerPropertiesWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger_property", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createTektonPipelineTriggerPropertiesOptions.PipelineID, *createTektonPipelineTriggerPropertiesOptions.TriggerID, *triggerProperty.Name))

	return resourceIBMCdTektonPipelineTriggerPropertyRead(context, d, meta)
}

func resourceIBMCdTektonPipelineTriggerPropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getTektonPipelineTriggerPropertyOptions := &cdtektonpipelinev2.GetTektonPipelineTriggerPropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "sep-id-parts").GetDiag()
	}

	getTektonPipelineTriggerPropertyOptions.SetPipelineID(parts[0])
	getTektonPipelineTriggerPropertyOptions.SetTriggerID(parts[1])
	getTektonPipelineTriggerPropertyOptions.SetPropertyName(parts[2])

	triggerProperty, response, err := cdTektonPipelineClient.GetTektonPipelineTriggerPropertyWithContext(context, getTektonPipelineTriggerPropertyOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetTektonPipelineTriggerPropertyWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger_property", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("name", triggerProperty.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-name").GetDiag()
	}
	if !core.IsNil(triggerProperty.Value) {
		if err = d.Set("value", triggerProperty.Value); err != nil {
			err = fmt.Errorf("Error setting value: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-value").GetDiag()
		}
	}
	if !core.IsNil(triggerProperty.Enum) {
		if err = d.Set("enum", triggerProperty.Enum); err != nil {
			err = fmt.Errorf("Error setting enum: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-enum").GetDiag()
		}
	}
	if err = d.Set("type", triggerProperty.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-type").GetDiag()
	}
	if !core.IsNil(triggerProperty.Path) {
		if err = d.Set("path", triggerProperty.Path); err != nil {
			err = fmt.Errorf("Error setting path: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-path").GetDiag()
		}
	}
	if !core.IsNil(triggerProperty.Locked) {
		if err = d.Set("locked", triggerProperty.Locked); err != nil {
			err = fmt.Errorf("Error setting locked: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-locked").GetDiag()
		}
	}
	if !core.IsNil(triggerProperty.Href) {
		if err = d.Set("href", triggerProperty.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "read", "set-href").GetDiag()
		}
	}

	return nil
}

func resourceIBMCdTektonPipelineTriggerPropertyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	replaceTektonPipelineTriggerPropertyOptions := &cdtektonpipelinev2.ReplaceTektonPipelineTriggerPropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "update", "sep-id-parts").GetDiag()
	}

	replaceTektonPipelineTriggerPropertyOptions.SetPipelineID(parts[0])
	replaceTektonPipelineTriggerPropertyOptions.SetTriggerID(parts[1])
	replaceTektonPipelineTriggerPropertyOptions.SetPropertyName(parts[2])
	replaceTektonPipelineTriggerPropertyOptions.SetName(d.Get("name").(string))
	replaceTektonPipelineTriggerPropertyOptions.SetType(d.Get("type").(string))

	hasChange := false

	if d.HasChange("pipeline_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "pipeline_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_trigger_property", "update", "pipeline_id-forces-new").GetDiag()
	}
	if d.HasChange("trigger_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "trigger_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_trigger_property", "update", "trigger_id-forces-new").GetDiag()
	}
	if d.HasChange("name") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "name")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_trigger_property", "update", "name-forces-new").GetDiag()
	}
	if d.HasChange("type") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "type")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_cd_tekton_pipeline_trigger_property", "update", "type-forces-new").GetDiag()
	}
	if d.HasChange("locked") {
		replaceTektonPipelineTriggerPropertyOptions.SetLocked(d.Get("locked").(bool))
		hasChange = true
	}
	if d.Get("type").(string) == "integration" {
		if d.HasChange("value") || d.HasChange("path") || d.HasChange("locked") {
			replaceTektonPipelineTriggerPropertyOptions.SetValue(d.Get("value").(string))
			replaceTektonPipelineTriggerPropertyOptions.SetPath(d.Get("path").(string))
			hasChange = true
		}
	} else if d.Get("type").(string) == "single_select" {
		if d.HasChange("enum") || d.HasChange("value") || d.HasChange("locked") {
			var enum []string
			for _, v := range d.Get("enum").([]interface{}) {
				enumItem := v.(string)
				enum = append(enum, enumItem)
			}
			replaceTektonPipelineTriggerPropertyOptions.SetEnum(enum)
			replaceTektonPipelineTriggerPropertyOptions.SetValue(d.Get("value").(string))
			hasChange = true
		}
	} else {
		if d.HasChange("value") || d.HasChange("locked") {
			replaceTektonPipelineTriggerPropertyOptions.SetValue(d.Get("value").(string))
			hasChange = true
		}
	}

	if hasChange {
		_, _, err = cdTektonPipelineClient.ReplaceTektonPipelineTriggerPropertyWithContext(context, replaceTektonPipelineTriggerPropertyOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceTektonPipelineTriggerPropertyWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger_property", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMCdTektonPipelineTriggerPropertyRead(context, d, meta)
}

func resourceIBMCdTektonPipelineTriggerPropertyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteTektonPipelineTriggerPropertyOptions := &cdtektonpipelinev2.DeleteTektonPipelineTriggerPropertyOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cd_tekton_pipeline_trigger_property", "delete", "sep-id-parts").GetDiag()
	}

	deleteTektonPipelineTriggerPropertyOptions.SetPipelineID(parts[0])
	deleteTektonPipelineTriggerPropertyOptions.SetTriggerID(parts[1])
	deleteTektonPipelineTriggerPropertyOptions.SetPropertyName(parts[2])

	_, err = cdTektonPipelineClient.DeleteTektonPipelineTriggerPropertyWithContext(context, deleteTektonPipelineTriggerPropertyOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteTektonPipelineTriggerPropertyWithContext failed: %s", err.Error()), "ibm_cd_tekton_pipeline_trigger_property", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
