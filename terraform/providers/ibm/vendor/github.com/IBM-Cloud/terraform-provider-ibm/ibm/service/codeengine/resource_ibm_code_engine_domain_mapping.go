// Copyright IBM Corp. 2024 All Rights Reserved.
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
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_domain_mapping", "project_id"),
				Description:  "The ID of the project.",
			},
			"component": {
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "A reference to another component.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name of the referenced component.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the referenced resource.",
						},
					},
				},
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_domain_mapping", "name"),
				Description:  "The name of the domain mapping.",
			},
			"tls_secret": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_code_engine_domain_mapping", "tls_secret"),
				Description:  "The name of the TLS secret that includes the certificate and private key of this domain mapping.",
			},
			"cname_target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The value of the CNAME record that must be configured in the DNS settings of the domain, to route traffic properly to the target Code Engine region.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The timestamp when the resource was created.",
			},
			"entity_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The version of the domain mapping instance, which is used to achieve optimistic locking.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "When you provision a new domain mapping, a URL is created identifying the location of the instance.",
			},
			"domain_mapping_id": {
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
				Description: "The type of the Code Engine resource.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the domain mapping.",
			},
			"status_details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The detailed status of the domain mapping.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Optional information to provide more context in case of a 'failed' or 'warning' status.",
						},
					},
				},
			},
			"user_managed": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Specifies whether the domain mapping is managed by the user or by Code Engine.",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies whether the domain mapping is reachable through the public internet, or private IBM network, or only through other components within the same Code Engine project.",
			},
			"etag": {
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
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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

	domainMapping, _, err := codeEngineClient.CreateDomainMappingWithContext(context, createDomainMappingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateDomainMappingWithContext failed: %s", err.Error()), "ibm_code_engine_domain_mapping", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createDomainMappingOptions.ProjectID, *domainMapping.Name))

	_, err = waitForIbmCodeEngineDomainMappingCreate(d, meta)
	if err != nil {
		errMsg := fmt.Sprintf("Error waiting for resource IbmCodeEngineDomainMapping (%s) to be created: %s", d.Id(), err)
		tfErr := flex.TerraformErrorf(err, errMsg, "ibm_code_engine_domain_mapping", "create")
		return tfErr.GetDiag()
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
				if sdkErr, ok := err.(*core.SDKProblem); ok && response.GetStatusCode() == 404 {
					sdkErr.Summary = fmt.Sprintf("The instance %s does not exist anymore: %s", "getDomainMappingOptions", err)
					return nil, "", sdkErr
				}
				return nil, "", err
			}
			failStates := map[string]bool{"failure": true, "failed": true}
			if failStates[*stateObj.Status] {
				return stateObj, *stateObj.Status, fmt.Errorf("the instance %s failed: %s", "getDomainMappingOptions", err)
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
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getDomainMappingOptions := &codeenginev2.GetDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "read")
		return tfErr.GetDiag()
	}

	getDomainMappingOptions.SetProjectID(parts[0])
	getDomainMappingOptions.SetName(parts[1])

	domainMapping, response, err := codeEngineClient.GetDomainMappingWithContext(context, getDomainMappingOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetDomainMappingWithContext failed: %s", err.Error()), "ibm_code_engine_domain_mapping", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("project_id", domainMapping.ProjectID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting project_id: %s", err))
	}
	componentMap, err := resourceIbmCodeEngineDomainMappingComponentRefToMap(domainMapping.Component)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("component", []map[string]interface{}{componentMap}); err != nil {
		return diag.FromErr(fmt.Errorf("error setting component: %s", err))
	}
	if err = d.Set("name", domainMapping.Name); err != nil {
		return diag.FromErr(fmt.Errorf("error setting name: %s", err))
	}
	if err = d.Set("tls_secret", domainMapping.TlsSecret); err != nil {
		return diag.FromErr(fmt.Errorf("error setting tls_secret: %s", err))
	}
	if !core.IsNil(domainMapping.CnameTarget) {
		if err = d.Set("cname_target", domainMapping.CnameTarget); err != nil {
			return diag.FromErr(fmt.Errorf("error setting cname_target: %s", err))
		}
	}
	if !core.IsNil(domainMapping.CreatedAt) {
		if err = d.Set("created_at", domainMapping.CreatedAt); err != nil {
			return diag.FromErr(fmt.Errorf("error setting created_at: %s", err))
		}
	}
	if err = d.Set("entity_tag", domainMapping.EntityTag); err != nil {
		return diag.FromErr(fmt.Errorf("error setting entity_tag: %s", err))
	}
	if !core.IsNil(domainMapping.Href) {
		if err = d.Set("href", domainMapping.Href); err != nil {
			return diag.FromErr(fmt.Errorf("error setting href: %s", err))
		}
	}
	if !core.IsNil(domainMapping.ID) {
		if err = d.Set("domain_mapping_id", domainMapping.ID); err != nil {
			return diag.FromErr(fmt.Errorf("error setting domain_mapping_id: %s", err))
		}
	}
	if !core.IsNil(domainMapping.Region) {
		if err = d.Set("region", domainMapping.Region); err != nil {
			return diag.FromErr(fmt.Errorf("error setting region: %s", err))
		}
	}
	if !core.IsNil(domainMapping.ResourceType) {
		if err = d.Set("resource_type", domainMapping.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("error setting resource_type: %s", err))
		}
	}
	if !core.IsNil(domainMapping.Status) {
		if err = d.Set("status", domainMapping.Status); err != nil {
			return diag.FromErr(fmt.Errorf("error setting status: %s", err))
		}
	}
	if !core.IsNil(domainMapping.StatusDetails) {
		statusDetailsMap, err := resourceIbmCodeEngineDomainMappingDomainMappingStatusToMap(domainMapping.StatusDetails)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("status_details", []map[string]interface{}{statusDetailsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("error setting status_details: %s", err))
		}
	}
	if !core.IsNil(domainMapping.UserManaged) {
		if err = d.Set("user_managed", domainMapping.UserManaged); err != nil {
			return diag.FromErr(fmt.Errorf("error setting user_managed: %s", err))
		}
	}
	if !core.IsNil(domainMapping.Visibility) {
		if err = d.Set("visibility", domainMapping.Visibility); err != nil {
			return diag.FromErr(fmt.Errorf("error setting visibility: %s", err))
		}
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting etag: %s", err), "ibm_code_engine_domain_mapping", "read")
		return tfErr.GetDiag()
	}

	return nil
}

func resourceIbmCodeEngineDomainMappingUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateDomainMappingOptions := &codeenginev2.UpdateDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "update")
		return tfErr.GetDiag()
	}

	updateDomainMappingOptions.SetProjectID(parts[0])
	updateDomainMappingOptions.SetName(parts[1])

	hasChange := false

	patchVals := &codeenginev2.DomainMappingPatch{}
	if d.HasChange("project_id") {
		errMsg := fmt.Sprintf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "project_id")
		tfErr := flex.TerraformErrorf(err, errMsg, "ibm_code_engine_domain_mapping", "update")
		return tfErr.GetDiag()
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
		_, _, err = codeEngineClient.UpdateDomainMappingWithContext(context, updateDomainMappingOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateDomainMappingWithContext failed: %s", err.Error()), "ibm_code_engine_domain_mapping", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIbmCodeEngineDomainMappingRead(context, d, meta)
}

func resourceIbmCodeEngineDomainMappingDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	codeEngineClient, err := meta.(conns.ClientSession).CodeEngineV2()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteDomainMappingOptions := &codeenginev2.DeleteDomainMappingOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_code_engine_domain_mapping", "delete")
		return tfErr.GetDiag()
	}

	deleteDomainMappingOptions.SetProjectID(parts[0])
	deleteDomainMappingOptions.SetName(parts[1])

	_, err = codeEngineClient.DeleteDomainMappingWithContext(context, deleteDomainMappingOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteDomainMappingWithContext failed: %s", err.Error()), "ibm_code_engine_domain_mapping", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
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
