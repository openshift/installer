// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.94.1-71478489-20240820-161623
 */

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIbmCodeEngineConfigMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineConfigMapCreate,
		ReadContext:   resourceIbmCodeEngineConfigMapRead,
		UpdateContext: resourceIbmCodeEngineConfigMapUpdate,
		DeleteContext: resourceIbmCodeEngineConfigMapDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_config_map", "project_id"),
				Description:  "The ID of the project.",
			},
			"data": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "The key-value pair for the config map. Values must be specified in `KEY=VALUE` format.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_config_map", "name"),
				Description:  "The name of the config map.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the config map instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new config map,  a URL is created identifying the location of the instance.",
			},
			"config_map_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the config map.",
			},
			"etag": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineConfigMapValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "project_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$`,
			MinValueLength:             36,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_config_map", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineConfigMapCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createConfigMapOptions := &codeenginev2.CreateConfigMapOptions{}

	createConfigMapOptions.SetProjectID(d.Get("project_id").(string))
	createConfigMapOptions.SetName(d.Get("name").(string))
	if _, ok := d.GetOk("data"); ok {
		data := make(map[string]string)
		for k, v := range d.Get("data").(map[string]interface{}) {
			data[k] = v.(string)
		}
		createConfigMapOptions.SetData(data)
	}

	configMap, _, err := codeEngineClient.CreateConfigMapWithContext(context, createConfigMapOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateConfigMapWithContext failed: %s", err.Error()), "ibm_code_engine_config_map", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createConfigMapOptions.ProjectID, *configMap.Name))

	return resourceIbmCodeEngineConfigMapRead(context, d, meta)
}

func resourceIbmCodeEngineConfigMapRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getConfigMapOptions := &codeenginev2.GetConfigMapOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "sep-id-parts").GetDiag()
	}

	getConfigMapOptions.SetProjectID(parts[0])
	getConfigMapOptions.SetName(parts[1])

	configMap, response, err := codeEngineClient.GetConfigMapWithContext(context, getConfigMapOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigMapWithContext failed: %s", err.Error()), "ibm_code_engine_config_map", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", configMap.ProjectID); err != nil {
		err = fmt.Errorf("Error setting project_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-project_id").GetDiag()
	}
	if !core.IsNil(configMap.Data) {
		data := make(map[string]string)
		for k, v := range configMap.Data {
			data[k] = string(v)
		}
		if err = d.Set("data", data); err != nil {
			err = fmt.Errorf("Error setting data: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-data").GetDiag()
		}
	}
	if err = d.Set("name", configMap.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-name").GetDiag()
	}
	if !core.IsNil(configMap.CreatedAt) {
		if err = d.Set("created_at", configMap.CreatedAt); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-created_at").GetDiag()
		}
	}
	if err = d.Set("entity_tag", configMap.EntityTag); err != nil {
		err = fmt.Errorf("Error setting entity_tag: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-entity_tag").GetDiag()
	}
	if !core.IsNil(configMap.Href) {
		if err = d.Set("href", configMap.Href); err != nil {
			err = fmt.Errorf("Error setting href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-href").GetDiag()
		}
	}
	if !core.IsNil(configMap.ID) {
		if err = d.Set("config_map_id", configMap.ID); err != nil {
			err = fmt.Errorf("Error setting config_map_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-config_map_id").GetDiag()
		}
	}
	if !core.IsNil(configMap.Region) {
		if err = d.Set("region", configMap.Region); err != nil {
			err = fmt.Errorf("Error setting region: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-region").GetDiag()
		}
	}
	if !core.IsNil(configMap.ResourceType) {
		if err = d.Set("resource_type", configMap.ResourceType); err != nil {
			err = fmt.Errorf("Error setting resource_type: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "read", "set-resource_type").GetDiag()
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_config_map", "read", "set-etag").GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineConfigMapUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	replaceConfigMapOptions := &codeenginev2.ReplaceConfigMapOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "update", "sep-id-parts").GetDiag()
	}

	replaceConfigMapOptions.SetProjectID(parts[0])
	replaceConfigMapOptions.SetName(parts[1])

	hasChange := false

	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		return flex.DiscriminatedTerraformErrorf(nil, errMsg, "ibm_code_engine_config_map", "update", "project_id-forces-new").GetDiag()
	}
	if d.HasChange("data") {
		data := make(map[string]string)
		for k, v := range d.Get("data").(map[string]interface{}) {
			data[k] = v.(string)
		}
		replaceConfigMapOptions.SetData(data)
		hasChange = true
	}
	replaceConfigMapOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		_, _, err = codeEngineClient.ReplaceConfigMapWithContext(context, replaceConfigMapOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReplaceConfigMapWithContext failed: %s", err.Error()), "ibm_code_engine_config_map", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineConfigMapRead(context, d, meta)
}

func resourceIbmCodeEngineConfigMapDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteConfigMapOptions := &codeenginev2.DeleteConfigMapOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_code_engine_config_map", "delete", "sep-id-parts").GetDiag()
	}

	deleteConfigMapOptions.SetProjectID(parts[0])
	deleteConfigMapOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteConfigMapWithContext(context, deleteConfigMapOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteConfigMapWithContext failed: %s", err.Error()), "ibm_code_engine_config_map", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
