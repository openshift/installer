// Copyright IBM Corp. 2022 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cdtektonpipeline

import (
	"context"
	"crypto/hmac"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/sha3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/continuous-delivery-go-sdk/cdtektonpipelinev2"
	"github.com/google/go-cmp/cmp"
)

func ResourceIBMTektonPipelineProperty() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceIBMTektonPipelinePropertyCreate,
		ReadContext:   ResourceIBMTektonPipelinePropertyRead,
		UpdateContext: ResourceIBMTektonPipelinePropertyUpdate,
		DeleteContext: ResourceIBMTektonPipelinePropertyDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"pipeline_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "pipeline_id"),
				Description:  "The tekton pipeline ID.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "name"),
				Description:  "Property name.",
			},
			"value": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "value"),
				Description:  "String format property value.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Get("type").(string) == "SECURE" {
						segs := []string{d.Get("pipeline_id").(string), d.Get("name").(string)}
						secret := strings.Join(segs, ".")
						mac := hmac.New(sha3.New512, []byte(secret))
						mac.Write([]byte(new))
						secureHmac := hex.EncodeToString(mac.Sum(nil))
						hasEnvChange := !cmp.Equal(strings.Join([]string{"hash", "SHA3-512", secureHmac}, ":"), old)
						if hasEnvChange {
							return false
						}
						return true
					} else {
						if old == new {
							return true
						}
						return false
					}
				},
			},
			"enum": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Options for SINGLE_SELECT property type.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"default": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "default"),
				Description:  "Default option for SINGLE_SELECT property type.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "type"),
				Description:  "Property type.",
			},
			"path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_cd_tekton_pipeline_property", "path"),
				Description:  "property path for INTEGRATION type properties.",
			},
		},
	}
}

func ResourceIBMTektonPipelinePropertyValidator() *validate.ResourceValidator {
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
			Regexp:                     `^[-0-9a-zA-Z_.]{1,234}$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "value",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `.`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
		validate.ValidateSchema{
			Identifier:                 "default",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[-0-9a-zA-Z_.]{1,235}$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "APPCONFIG, INTEGRATION, SECURE, SINGLE_SELECT, TEXT",
		},
		validate.ValidateSchema{
			Identifier:                 "path",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `.`,
			MinValueLength:             1,
			MaxValueLength:             4096,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_cd_tekton_pipeline_property", Schema: validateSchema}
	return &resourceValidator
}

func ResourceIBMTektonPipelinePropertyCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cdTektonPipelineClient, err := meta.(conns.ClientSession).CdTektonPipelineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createTektonPipelinePropertiesOptions := &cdtektonpipelinev2.CreateTektonPipelinePropertiesOptions{}

	createTektonPipelinePropertiesOptions.SetPipelineID(d.Get("pipeline_id").(string))
	if _, ok := d.GetOk("name"); ok {
		createTektonPipelinePropertiesOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("value"); ok {
		createTektonPipelinePropertiesOptions.SetValue(d.Get("value").(string))
	}
	if _, ok := d.GetOk("enum"); ok {
		enumInterface := d.Get("enum").([]interface{})
		enum := make([]string, len(enumInterface))
		for i, v := range enumInterface {
			enum[i] = fmt.Sprint(v)
		}
		createTektonPipelinePropertiesOptions.SetEnum(enum)
	}
	if _, ok := d.GetOk("default"); ok {
		createTektonPipelinePropertiesOptions.SetDefault(d.Get("default").(string))
	}
	if _, ok := d.GetOk("type"); ok {
		createTektonPipelinePropertiesOptions.SetType(d.Get("type").(string))
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

	return ResourceIBMTektonPipelinePropertyRead(context, d, meta)
}

func ResourceIBMTektonPipelinePropertyRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err = d.Set("value", property.Value); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting value: %s", err))
	}
	if property.Enum != nil {
		if err = d.Set("enum", property.Enum); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enum: %s", err))
		}
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

	return nil
}

func ResourceIBMTektonPipelinePropertyUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChange("name") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "name"))
	}

	if d.Get("type").(string) == "INTEGRATION" {
		if d.HasChange("value") || d.HasChange("path") {
			replaceTektonPipelinePropertyOptions.SetValue(d.Get("value").(string))
			replaceTektonPipelinePropertyOptions.SetPath(d.Get("path").(string))
			hasChange = true
		}
	} else if d.Get("type").(string) == "SINGLE_SELECT" {
		if d.HasChange("enum") || d.HasChange("default") {
			enumInterface := d.Get("enum").([]interface{})
			enum := make([]string, len(enumInterface))
			for i, v := range enumInterface {
				enum[i] = fmt.Sprint(v)
			}
			replaceTektonPipelinePropertyOptions.SetEnum(enum)
			replaceTektonPipelinePropertyOptions.SetDefault(d.Get("default").(string))
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

	return ResourceIBMTektonPipelinePropertyRead(context, d, meta)
}

func ResourceIBMTektonPipelinePropertyDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
