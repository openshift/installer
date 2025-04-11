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
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Description: "The value that is set as a prefix in the component that is bound.",
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
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_binding", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getBindingOptions := &codeenginev2.GetBindingOptions{}

	getBindingOptions.SetProjectID(d.Get("project_id").(string))
	getBindingOptions.SetID(d.Get("binding_id").(string))

	binding, _, err := codeEngineClient.GetBindingWithContext(context, getBindingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBindingWithContext failed: %s", err.Error()), "(Data) ibm_code_engine_binding", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getBindingOptions.ProjectID, *getBindingOptions.ID))

	component := []map[string]interface{}{}
	componentMap, err := DataSourceIbmCodeEngineBindingComponentRefToMap(binding.Component)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_code_engine_binding", "read", "component-to-map").GetDiag()
	}
	component = append(component, componentMap)
	if err = d.Set("component", component); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting component: %s", err), "(Data) ibm_code_engine_binding", "read", "set-component").GetDiag()
	}

	if !core.IsNil(binding.Href) {
		if err = d.Set("href", binding.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_code_engine_binding", "read", "set-href").GetDiag()
		}
	}

	if err = d.Set("prefix", binding.Prefix); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting prefix: %s", err), "(Data) ibm_code_engine_binding", "read", "set-prefix").GetDiag()
	}

	if !core.IsNil(binding.ResourceType) {
		if err = d.Set("resource_type", binding.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_code_engine_binding", "read", "set-resource_type").GetDiag()
		}
	}

	if err = d.Set("secret_name", binding.SecretName); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting secret_name: %s", err), "(Data) ibm_code_engine_binding", "read", "set-secret_name").GetDiag()
	}

	if !core.IsNil(binding.Status) {
		if err = d.Set("status", binding.Status); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_code_engine_binding", "read", "set-status").GetDiag()
		}
	}

	return nil
}

func DataSourceIbmCodeEngineBindingComponentRefToMap(model *codeenginev2.ComponentRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}
