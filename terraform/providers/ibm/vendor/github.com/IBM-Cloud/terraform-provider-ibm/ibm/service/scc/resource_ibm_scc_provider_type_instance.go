// ng Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package scc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/scc-go-sdk/v5/securityandcompliancecenterapiv3"
)

func ResourceIbmSccProviderTypeInstance() *schema.Resource {
	return AddSchemaData(&schema.Resource{
		CreateContext: resourceIbmSccProviderTypeInstanceCreate,
		ReadContext:   resourceIbmSccProviderTypeInstanceRead,
		UpdateContext: resourceIbmSccProviderTypeInstanceUpdate,
		DeleteContext: resourceIbmSccProviderTypeInstanceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"provider_type_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_provider_type_instance", "provider_type_id"),
				Description:  "The provider type ID.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_scc_provider_type_instance", "name"),
				Description:  "The name of the provider type instance.",
			},
			"attributes": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the provider type.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which resource was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Time at which resource was updated.",
			},
			"provider_type_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the provider type instance.",
			},
		},
	})
}

func ResourceIbmSccProviderTypeInstanceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "provider_type_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[a-zA-Z0-9 ,\-_]+$`,
			MinValueLength:             32,
			MaxValueLength:             36,
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `[A-Za-z0-9]+`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_scc_provider_type_instance", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmSccProviderTypeInstanceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	createProviderTypeInstanceOptions := &securityandcompliancecenterapiv3.CreateProviderTypeInstanceOptions{}
	instanceID := d.Get("instance_id").(string)

	createProviderTypeInstanceOptions.SetProviderTypeID(d.Get("provider_type_id").(string))
	createProviderTypeInstanceOptions.SetName(d.Get("name").(string))
	createProviderTypeInstanceOptions.SetInstanceID(instanceID)
	attributesModel, err := resourceIbmSccProviderTypeInstanceMapToProviderTypeInstanceAttributes(d.Get("attributes").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createProviderTypeInstanceOptions.SetAttributes(attributesModel)

	providerTypeInstanceItem, response, err := securityAndComplianceCenterApIsClient.CreateProviderTypeInstanceWithContext(context, createProviderTypeInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateProviderTypeInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("CreateProviderTypeInstanceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceID, *createProviderTypeInstanceOptions.ProviderTypeID, *providerTypeInstanceItem.ID))

	return resourceIbmSccProviderTypeInstanceRead(context, d, meta)
}

func resourceIbmSccProviderTypeInstanceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	getProviderTypeInstanceOptions := &securityandcompliancecenterapiv3.GetProviderTypeInstanceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getProviderTypeInstanceOptions.SetInstanceID(parts[0])
	getProviderTypeInstanceOptions.SetProviderTypeID(parts[1])
	getProviderTypeInstanceOptions.SetProviderTypeInstanceID(parts[2])

	providerTypeInstanceItem, response, err := securityAndComplianceCenterApIsClient.GetProviderTypeInstanceWithContext(context, getProviderTypeInstanceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetProviderTypeInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("GetProviderTypeInstanceWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("instance_id", parts[0]); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting instance_id: %s", err))
	}

	if err = d.Set("name", providerTypeInstanceItem.Name); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting name: %s", err))
	}

	attributesMap, err := resourceIbmSccProviderTypeInstanceProviderTypeInstanceAttributesToMap(providerTypeInstanceItem.Attributes)
	if err != nil {
		return diag.FromErr(err)
	}

	if err = d.Set("attributes", attributesMap); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting attributes: %s", err))
	}

	if !core.IsNil(providerTypeInstanceItem.Type) {
		if err = d.Set("type", providerTypeInstanceItem.Type); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting type: %s", err))
		}
	}
	if !core.IsNil(providerTypeInstanceItem.CreatedAt) {
		if err = d.Set("created_at", flex.DateTimeToString(providerTypeInstanceItem.CreatedAt)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting created_at: %s", err))
		}
	}
	if !core.IsNil(providerTypeInstanceItem.UpdatedAt) {
		if err = d.Set("updated_at", flex.DateTimeToString(providerTypeInstanceItem.UpdatedAt)); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting updated_at: %s", err))
		}
	}
	if !core.IsNil(providerTypeInstanceItem.ID) {
		if err = d.Set("provider_type_instance_id", providerTypeInstanceItem.ID); err != nil {
			return diag.FromErr(flex.FmtErrorf("Error setting provider_type_instance_id: %s", err))
		}
	}
	if err = d.Set("provider_type_id", parts[1]); err != nil {
		return diag.FromErr(flex.FmtErrorf("Error setting provider_type_id: %s", err))
	}

	return nil
}

func resourceIbmSccProviderTypeInstanceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	updateProviderTypeInstanceOptions := &securityandcompliancecenterapiv3.UpdateProviderTypeInstanceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: add updateProviderTypeInstanceOptions.SetInstanceID to scc-go-sdk
	updateProviderTypeInstanceOptions.SetInstanceID(parts[0])
	updateProviderTypeInstanceOptions.SetProviderTypeID(parts[1])
	updateProviderTypeInstanceOptions.SetProviderTypeInstanceID(parts[2])

	hasChange := false

	if d.HasChange("provider_type_id") {
		return diag.FromErr(flex.FmtErrorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "provider_type_id"))
	}
	if d.HasChange("attributes") {
		updateProviderTypeInstanceOptions.SetAttributes(d.Get("attributes").(map[string]interface{}))
		hasChange = true
	}
	if d.HasChange("name") {
		updateProviderTypeInstanceOptions.SetName(d.Get("name").(string))
		hasChange = true
	}

	if hasChange {
		_, response, err := securityAndComplianceCenterApIsClient.UpdateProviderTypeInstanceWithContext(context, updateProviderTypeInstanceOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateProviderTypeInstanceWithContext failed %s\n%s", err, response)
			return diag.FromErr(flex.FmtErrorf("UpdateProviderTypeInstanceWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmSccProviderTypeInstanceRead(context, d, meta)
}

func resourceIbmSccProviderTypeInstanceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	securityAndComplianceCenterApIsClient, err := meta.(conns.ClientSession).SecurityAndComplianceCenterV3()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProviderTypeInstanceOptions := &securityandcompliancecenterapiv3.DeleteProviderTypeInstanceOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteProviderTypeInstanceOptions.SetInstanceID(parts[0])
	deleteProviderTypeInstanceOptions.SetProviderTypeID(parts[1])
	deleteProviderTypeInstanceOptions.SetProviderTypeInstanceID(parts[2])

	response, err := securityAndComplianceCenterApIsClient.DeleteProviderTypeInstanceWithContext(context, deleteProviderTypeInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteProviderTypeInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(flex.FmtErrorf("DeleteProviderTypeInstanceWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmSccProviderTypeInstanceMapToProviderTypeInstanceAttributes(modelMap map[string]interface{}) (map[string]interface{}, error) {
	return modelMap, nil
}

func resourceIbmSccProviderTypeInstanceProviderTypeInstanceAttributesToMap(modelMap map[string]interface{}) (map[string]interface{}, error) {
	return modelMap, nil
}
