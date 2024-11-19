// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineConfigMap() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineConfigMapRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your configmap.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"data": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The key-value pair for the config map. Values must be specified in `KEY=VALUE` format.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
		},
	}
}

func dataSourceIbmCodeEngineConfigMapRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_config_map", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getConfigMapOptions := &codeenginev2.GetConfigMapOptions{}

	getConfigMapOptions.SetProjectID(d.Get("project_id").(string))
	getConfigMapOptions.SetName(d.Get("name").(string))

	configMap, _, err := codeEngineClient.GetConfigMapWithContext(context, getConfigMapOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetConfigMapWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_config_map", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getConfigMapOptions.ProjectID, *getConfigMapOptions.Name))

	if err = d.Set("created_at", configMap.CreatedAt); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting created_at: %s", err), "(Data) ibm_code_engine_config_map", "read")
		return tfErr.GetDiag()
	}

	if configMap.Data != nil {
		if err = d.Set("data", configMap.Data); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting data: %s", err), "(Data) ibm_code_engine_config_map", "read")
			return tfErr.GetDiag()
		}
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting data: %s", err), "(Data) ibm_code_engine_config_map", "read")
			return tfErr.GetDiag()
		}
	}

	if err = d.Set("entity_tag", configMap.EntityTag); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting entity_tag: %s", err), "(Data) ibm_code_engine_config_map", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("href", configMap.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting href: %s", err), "(Data) ibm_code_engine_config_map", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("region", configMap.Region); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting region: %s", err), "(Data) ibm_code_engine_config_map", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("config_map_id", configMap.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting config_map_id: %s", err), "(Data) ibm_code_engine_config_map", "read")
		return tfErr.GetDiag()
	}

	if err = d.Set("resource_type", configMap.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("error setting resource_type: %s", err), "(Data) ibm_code_engine_config_map", "read")
		return tfErr.GetDiag()
	}

	return nil
}
