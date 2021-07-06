// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISRouting          = "ibm_cis_routing"
	cisRoutingSmartRouting = "smart_routing"
)

func resourceIBMCISRouting() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISRoutingUpdate,
		Read:     resourceIBMCISRoutingRead,
		Update:   resourceIBMCISRoutingUpdate,
		Delete:   resourceIBMCISRoutingDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
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
				ValidateFunc: InvokeValidator(ibmCISRouting, cisRoutingSmartRouting),
			},
		},
	}
}

func resourceIBMCISRoutingValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	smartRoutingValues := "on, off"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisRoutingSmartRouting,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              smartRoutingValues})
	ibmCISRoutingValidator := ResourceValidator{ResourceName: ibmCISRouting, Schema: validateSchema}
	return &ibmCISRoutingValidator
}

func resourceIBMCISRoutingUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRoutingClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
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

	d.SetId(convertCisToTfTwoVar(zoneID, crn))
	return resourceIBMCISRoutingRead(d, meta)
}

func resourceIBMCISRoutingRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisRoutingClientSession()
	if err != nil {
		return err
	}
	zoneID, crn, err := convertTftoCisTwoVar(d.Id())
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

func resourceIBMCISRoutingDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS resource
	d.SetId("")
	return nil
}
