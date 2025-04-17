// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmCodeEngineBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineBindingCreate,
		ReadContext:   resourceIbmCodeEngineBindingRead,
		DeleteContext: resourceIbmCodeEngineBindingDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_binding", "project_id"),
				Description:  "The ID of the project.",
			},
			"component": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "A reference to another component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the referenced component.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the referenced resource.",
						},
					},
				},
			},
			"prefix": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_binding", "prefix"),
				Description:  "Optional value that is set as prefix in the component that is bound. Will be generated if not provided.",
			},
			"secret_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_binding", "secret_name"),
				Description:  "The service access secret that is binding to a component.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new binding,  a URL is created identifying the location of the instance.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the binding.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the binding.",
			},
			"binding_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the binding.",
			},
		},
	}
}

func ResourceIbmCodeEngineBindingValidator() *validate.ResourceValidator {
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
			Identifier:                 "prefix",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[A-Z]([_A-Z0-9]*[A-Z0-9])*$`,
			MinValueLength:             0,
			MaxValueLength:             31,
		},
		validate.ValidateSchema{
			Identifier:                 "secret_name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_binding", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineBindingCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createBindingOptions := &codeenginev2.CreateBindingOptions{}

	createBindingOptions.SetProjectID(d.Get("project_id").(string))
	componentModel, err := resourceIbmCodeEngineBindingMapToComponentRef(d.Get("component.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createBindingOptions.SetComponent(componentModel)
	createBindingOptions.SetPrefix(d.Get("prefix").(string))
	createBindingOptions.SetSecretName(d.Get("secret_name").(string))

	binding, response, err := codeEngineClient.CreateBindingWithContext(context, createBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateBindingWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createBindingOptions.ProjectID, *binding.ID))

	return resourceIbmCodeEngineBindingRead(context, d, meta)
}

func resourceIbmCodeEngineBindingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getBindingOptions := &codeenginev2.GetBindingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getBindingOptions.SetProjectID(parts[0])
	getBindingOptions.SetID(parts[1])

	binding, response, err := codeEngineClient.GetBindingWithContext(context, getBindingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBindingWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("project_id", binding.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_id: %s", err))
	}
	componentMap, err := resourceIbmCodeEngineBindingComponentRefToMap(binding.Component)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("component", []map[string]interface{}{componentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting component: %s", err))
	}
	if err = d.Set("prefix", binding.Prefix); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting prefix: %s", err))
	}
	if err = d.Set("secret_name", binding.SecretName); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secret_name: %s", err))
	}
	if !core.IsNil(binding.Href) {
		if err = d.Set("href", binding.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(binding.ResourceType) {
		if err = d.Set("resource_type", binding.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(binding.Status) {
		if err = d.Set("status", binding.Status); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
		}
	}
	if !core.IsNil(binding.ID) {
		if err = d.Set("binding_id", binding.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting binding_id: %s", err))
		}
	}

	return nil
}

func resourceIbmCodeEngineBindingDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBindingOptions := &codeenginev2.DeleteBindingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBindingOptions.SetProjectID(parts[0])
	deleteBindingOptions.SetID(parts[1])

	response, err := codeEngineClient.DeleteBindingWithContext(context, deleteBindingOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBindingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteBindingWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmCodeEngineBindingMapToComponentRef(modelMap map[string]interface{}) (*codeenginev2.ComponentRef, error) {
	model := &codeenginev2.ComponentRef{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.ResourceType = core.StringPtr(modelMap["resource_type"].(string))
	return model, nil
}

func resourceIbmCodeEngineBindingComponentRefToMap(model *codeenginev2.ComponentRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}
