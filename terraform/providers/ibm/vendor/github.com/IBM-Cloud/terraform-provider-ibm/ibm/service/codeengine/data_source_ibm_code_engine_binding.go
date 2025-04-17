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

func DataSourceIbmCodeEngineBinding() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmCodeEngineBindingRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project.",
			},
			"component": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A reference to another component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the referenced component.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the referenced resource.",
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new binding,  a URL is created identifying the location of the instance.",
			},
			"binding_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the binding.",
			},
			"prefix": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value that is set as prefix in the component that is bound.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the binding.",
			},
			"secret_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The service access secret that is bound to a component.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the binding.",
			},
		},
	}
}

func dataSourceIbmCodeEngineBindingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getBindingOptions := &codeenginev2.GetBindingOptions{}

	getBindingOptions.SetProjectID(d.Get("project_id").(string))
	getBindingOptions.SetID(d.Get("binding_id").(string))

	binding, response, err := codeEngineClient.GetBindingWithContext(context, getBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] GetBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBindingWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getBindingOptions.ProjectID, *getBindingOptions.ID))

	component := []map[string]interface{}{}
	if binding.Component != nil {
		modelMap, err := dataSourceIbmCodeEngineBindingComponentRefToMap(binding.Component)
		if err != nil {
			return diag.FromErr(err)
		}
		component = append(component, modelMap)
	}

	errString := "Error setting %s %s"
	if err = d.Set("component", component); err != nil {
		return diag.FromErr(fmt.Errorf(errString, "component", err))
	}

	if err = d.Set("href", binding.Href); err != nil {
		return diag.FromErr(fmt.Errorf(errString, "href", err))
	}

	if err = d.Set("prefix", binding.Prefix); err != nil {
		return diag.FromErr(fmt.Errorf(errString, "prefix", err))
	}

	if err = d.Set("resource_type", binding.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf(errString, "resource_type", err))
	}

	if err = d.Set("secret_name", binding.SecretName); err != nil {
		return diag.FromErr(fmt.Errorf(errString, "secret_name", err))
	}

	if err = d.Set("status", binding.Status); err != nil {
		return diag.FromErr(fmt.Errorf(errString, "status", err))
	}

	return nil
}

func dataSourceIbmCodeEngineBindingComponentRefToMap(model *codeenginev2.ComponentRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
