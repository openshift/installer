// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package codeengine

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/bluemix-go/bmxerror"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/code-engine-go-sdk/codeenginev2"
	"github.com/IBM/go-sdk-core/v5/core"
)

func ResourceIbmCodeEngineDomainMapping() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmCodeEngineDomainMappingCreate,
		ReadContext:   resourceIbmCodeEngineDomainMappingRead,
		UpdateContext: resourceIbmCodeEngineDomainMappingUpdate,
		DeleteContext: resourceIbmCodeEngineDomainMappingDelete,
		Importer:      &schema.ResourceImporter{},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_domain_mapping", "project_id"),
				Description:  "The ID of the project.",
			},
			"component": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_domain_mapping", "name"),
				Description:  "The name of the domain mapping.",
			},
			"tls_secret": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_domain_mapping", "tls_secret"),
				Description:  "The name of the TLS secret that holds the certificate and private key of this domain mapping.",
			},
			"cname_target": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exposes the value of the CNAME record that needs to be configured in the DNS settings of the domain, to route traffic properly to the target Code Engine region.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"domain_mapping_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The identifier of the resource.",
			},
			"entity_tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the domain mapping instance, which is used to achieve optimistic locking.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new domain mapping, a URL is created identifying the location of the instance.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the CE Resource.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the domain mapping.",
			},
			"status_details": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the domain mapping.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"user_managed": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Exposes whether the domain mapping is managed by the user or by Code Engine.",
			},
			"visibility": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Exposes whether the domain mapping is reachable through the public internet, or private IBM network, or only through other components within the same Code Engine project.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmCodeEngineDomainMappingValidator() *validate.ResourceValidator {
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
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)+$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
		validate.ValidateSchema{
			Identifier:                 "tls_secret",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-z0-9]([\-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([\-a-z0-9]*[a-z0-9])?)*$`,
			MinValueLength:             1,
			MaxValueLength:             253,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_code_engine_domain_mapping", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmCodeEngineDomainMappingCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	createDomainMappingOptions := &codeenginev2.CreateDomainMappingOptions{}

	createDomainMappingOptions.SetProjectID(d.Get("project_id").(string))
	componentModel, err := resourceIbmCodeEngineDomainMappingMapToComponentRef(d.Get("component.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createDomainMappingOptions.SetComponent(componentModel)
	createDomainMappingOptions.SetName(d.Get("name").(string))
	createDomainMappingOptions.SetTlsSecret(d.Get("tls_secret").(string))

	domainMapping, response, err := codeEngineClient.CreateDomainMappingWithContext(context, createDomainMappingOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateDomainMappingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateDomainMappingWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createDomainMappingOptions.ProjectID, *domainMapping.Name))

	_, err = waitForIbmCodeEngineDomainMappingCreate(d, meta)
	if err != nil {
		return diag.FromErr(fmt.Errorf(
			"Error waiting for resource IbmCodeEngineDomainMapping (%s) to be created: %s", d.Id(), err))
	}

	return resourceIbmCodeEngineDomainMappingRead(context, d, meta)
}

func waitForIbmCodeEngineDomainMappingCreate(d *schema.ResourceData, meta interface{}) (interface{}, error) {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return false, err
	}
	getDomainMappingOptions := &codeenginev2.GetDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return false, err
	}

	getDomainMappingOptions.SetProjectID(parts[0])
	getDomainMappingOptions.SetName(parts[1])

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deploying"},
		Target:  []string{"ready", "failed"},
		Refresh: func() (interface{}, string, error) {
			stateObj, response, err := codeEngineClient.GetDomainMapping(getDomainMappingOptions)
			if err != nil {
				if apiErr, ok := err.(bmxerror.RequestFailure); ok && apiErr.StatusCode() == 404 {
					return nil, "", fmt.Errorf("The instance %s does not exist anymore: %s\n%s", "getDomainMappingOptions", err, response)
				}
				return nil, "", err
			}
			failStates := map[string]bool{"failure": true, "failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("The instance %s failed: %s\n%s", "getDomainMappingOptions", err, response)
			}
			return stateObj, *stateObj.Status, nil
		},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      60 * time.Second,
		MinTimeout: 60 * time.Second,
	}

	return stateConf.WaitForState()
}

func resourceIbmCodeEngineDomainMappingRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	getDomainMappingOptions := &codeenginev2.GetDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getDomainMappingOptions.SetProjectID(parts[0])
	getDomainMappingOptions.SetName(parts[1])

	domainMapping, response, err := codeEngineClient.GetDomainMappingWithContext(context, getDomainMappingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDomainMappingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDomainMappingWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("project_id", domainMapping.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting project_id: %s", err))
	}
	componentMap, err := resourceIbmCodeEngineDomainMappingComponentRefToMap(domainMapping.Component)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("component", []map[string]interface{}{componentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting component: %s", err))
	}
	if err = d.Set("name", domainMapping.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if err = d.Set("tls_secret", domainMapping.TlsSecret); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting tls_secret: %s", err))
	}
	if !core.IsNil(domainMapping.CnameTarget) {
		if err = d.Set("cname_target", domainMapping.CnameTarget); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cname_target: %s", err))
		}
	}
	if !core.IsNil(domainMapping.CreatedAt) {
		if err = d.Set("created_at", domainMapping.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
		}
	}
	if err = d.Set("entity_tag", domainMapping.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting entity_tag: %s", err))
	}
	if !core.IsNil(domainMapping.Href) {
		if err = d.Set("href", domainMapping.Href); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
		}
	}
	if !core.IsNil(domainMapping.ID) {
		if err = d.Set("domain_mapping_id", domainMapping.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting domain_mapping_id: %s", err))
		}
	}
	if !core.IsNil(domainMapping.ResourceType) {
		if err = d.Set("resource_type", domainMapping.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(domainMapping.Status) {
		if err = d.Set("status", domainMapping.Status); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
		}
	}
	if !core.IsNil(domainMapping.StatusDetails) {
		statusDetailsMap, err := resourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(domainMapping.StatusDetails)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("status_details", []map[string]interface{}{statusDetailsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_details: %s", err))
		}
	}
	if !core.IsNil(domainMapping.UserManaged) {
		if err = d.Set("user_managed", domainMapping.UserManaged); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting user_managed: %s", err))
		}
	}
	if !core.IsNil(domainMapping.Visibility) {
		if err = d.Set("visibility", domainMapping.Visibility); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting visibility: %s", err))
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIbmCodeEngineDomainMappingUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	updateDomainMappingOptions := &codeenginev2.UpdateDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateDomainMappingOptions.SetProjectID(parts[0])
	updateDomainMappingOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.DomainMappingPatch{}
	if d.HasChange("project_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id"))
	}
	if d.HasChange("name") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "name"))
	}
	if d.HasChange("component") {
		component, err := resourceIbmCodeEngineDomainMappingMapToComponentRef(d.Get("component.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchVals.Component = component
		hasChange = true
	}
	if d.HasChange("tls_secret") {
		newTlsSecret := d.Get("tls_secret").(string)
		patchVals.TlsSecret = &newTlsSecret
		hasChange = true
	}
	updateDomainMappingOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateDomainMappingOptions.DomainMapping, _ = patchVals.AsPatch()
		_, response, err := codeEngineClient.UpdateDomainMappingWithContext(context, updateDomainMappingOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateDomainMappingWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateDomainMappingWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmCodeEngineDomainMappingRead(context, d, meta)
}

func resourceIbmCodeEngineDomainMappingDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDomainMappingOptions := &codeenginev2.DeleteDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDomainMappingOptions.SetProjectID(parts[0])
	deleteDomainMappingOptions.SetName(parts[1])

	response, err := codeEngineClient.DeleteDomainMappingWithContext(context, deleteDomainMappingOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDomainMappingWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteDomainMappingWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmCodeEngineDomainMappingMapToComponentRef(modelMap map[string]interface{}) (*codeenginev2.ComponentRef, error) {
	model := &codeenginev2.ComponentRef{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	model.ResourceType = core.StringPtr(modelMap["resource_type"].(string))
	return model, nil
}

func resourceIbmCodeEngineDomainMappingComponentRefToMap(model *codeenginev2.ComponentRef) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(model *codeenginev2.DomainMappingStatus) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Reason != nil {
		modelMap["reason"] = model.Reason
	}
	return modelMap, nil
}
