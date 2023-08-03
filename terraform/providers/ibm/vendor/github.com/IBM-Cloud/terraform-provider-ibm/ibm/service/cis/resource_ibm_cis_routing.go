// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	ibmCISRouting          = "ibm_cis_routing"
	cisRoutingSmartRouting = "smart_routing"
)

func ResourceIBMCISRouting() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISRoutingUpdate,
		Read:     ResourceIBMCISRoutingRead,
		Update:   ResourceIBMCISRoutingUpdate,
		Delete:   ResourceIBMCISRoutingDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator("ibm_cis_routing",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisRoutingSmartRouting: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Smart Routing value",
				ValidateFunc: validate.InvokeValidator(ibmCISRouting, cisRoutingSmartRouting),
			},
		},
	}
}

func ResourceIBMCISRoutingValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	smartRoutingValues := "on, off"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisRoutingSmartRouting,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              smartRoutingValues})
	ibmCISRoutingValidator := validate.ResourceValidator{ResourceName: ibmCISRouting, Schema: validateSchema}
	return &ibmCISRoutingValidator
}

func ResourceIBMCISRoutingUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisRoutingClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	if d.HasChange(cisRoutingSmartRouting) {
		smartRoutingValue := d.Get(cisRoutingSmartRouting).(string)
		opt := cisClient.NewUpdateSmartRoutingOptions()
		opt.SetValue(smartRoutingValue)
		_, response, err := cisClient.UpdateSmartRouting(opt)
		if err != nil {
			log.Printf("Update smart route setting failed: %v", response)
			return err
		}
	}

	d.SetId(flex.ConvertCisToTfTwoVar(zoneID, crn))
	return ResourceIBMCISRoutingRead(d, meta)
}

func ResourceIBMCISRoutingRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisRoutingClientSession()
	if err != nil {
		return err
	}
	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewGetSmartRoutingOptions()
	result, response, err := cisClient.GetSmartRouting(opt)
	if err != nil {
		log.Printf("Get smart route setting failed: %v", response)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisRoutingSmartRouting, *result.Result.Value)
	return nil
}

func ResourceIBMCISRoutingDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS resource
	d.SetId("")
	return nil
}
