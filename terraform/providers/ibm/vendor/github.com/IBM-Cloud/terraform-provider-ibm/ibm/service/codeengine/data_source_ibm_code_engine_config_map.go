// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
)

func DataSourceIbmCodeEngineConfigMap() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineConfigMapRead,

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of your configmap.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"data": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The key-value pair for the config map. Values must be specified in `KEY=VALUE` format.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the config map instance, which is used to achieve optimistic locking.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new config map,  a URL is created identifying the location of the instance.",
			},
			"config_map_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"resource_type": &schema.Schema{
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
		return diag.FromErr(err)
	}

	getConfigMapOptions := &codeenginev2.GetConfigMapOptions{}

	getConfigMapOptions.SetProjectID(d.Get("project_id").(string))
	getConfigMapOptions.SetName(d.Get("name").(string))

	configMap, response, err := codeEngineClient.GetConfigMapWithContext(context, getConfigMapOptions)
	if err != nil {
		log.Printf("[DEBUG] GetConfigMapWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetConfigMapWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getConfigMapOptions.ProjectID, *getConfigMapOptions.Name))

	if err = d.Set("created_at", configMap.CreatedAt); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if configMap.Data != nil {
		if err = d.Set("data", configMap.Data); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data: %s", err))
		}
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting data %s", err))
		}
	}

	if err = d.Set("entity_tag", configMap.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}

	if err = d.Set("href", configMap.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("config_map_id", configMap.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting config_map_id: %s", err))
	}

	if err = d.Set("resource_type", configMap.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}
