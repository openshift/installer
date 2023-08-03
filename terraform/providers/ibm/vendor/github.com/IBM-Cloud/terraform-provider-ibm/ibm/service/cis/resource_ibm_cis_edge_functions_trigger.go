// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	cisEdgeFunctionsTriggerID                   = "trigger_id"
	cisEdgeFunctionsTriggerPattern              = "pattern_url"
	cisEdgeFunctionsTriggerActionName           = "action_name"
	cisEdgeFunctionsTriggerRequestLimitFailOpen = "request_limit_fail_open"
)

func ResourceIBMCISEdgeFunctionsTrigger() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISEdgeFunctionsTriggerCreate,
		Read:     ResourceIBMCISEdgeFunctionsTriggerRead,
		Update:   ResourceIBMCISEdgeFunctionsTriggerUpdate,
		Delete:   ResourceIBMCISEdgeFunctionsTriggerDelete,
		Exists:   ResourceIBMCISEdgeFunctionsTriggerExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator("ibm_cis_edge_functions_trigger",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDataDiff,
			},
			cisEdgeFunctionsTriggerID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIS Edge Functions trigger route ID",
			},
			cisEdgeFunctionsTriggerPattern: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Edge function trigger pattern",
			},
			cisEdgeFunctionsTriggerActionName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Edge function trigger action name",
			},
			cisEdgeFunctionsTriggerRequestLimitFailOpen: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Edge function trigger request limit fail open",
			},
		},
	}
}
func ResourceIBMCISEdgeFunctionsTriggerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISEdgeFunctionsTriggerValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_edge_functions_trigger",
		Schema:       validateSchema}
	return &ibmCISEdgeFunctionsTriggerValidator
}

func ResourceIBMCISEdgeFunctionsTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)
	opt := cisClient.NewCreateEdgeFunctionsTriggerOptions()
	if action, ok := d.GetOk(cisEdgeFunctionsTriggerActionName); ok {
		opt.SetScript(action.(string))
	}
	pattern := d.Get(cisEdgeFunctionsTriggerPattern).(string)
	opt.SetPattern(pattern)

	result, _, err := cisClient.CreateEdgeFunctionsTrigger(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error creating edge function trigger route : %v", err)
	}
	d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return ResourceIBMCISEdgeFunctionsTriggerRead(d, meta)
}

func ResourceIBMCISEdgeFunctionsTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	routeID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	if d.HasChange(cisEdgeFunctionsTriggerActionName) ||
		d.HasChange(cisEdgeFunctionsTriggerPattern) {
		opt := cisClient.NewUpdateEdgeFunctionsTriggerOptions(routeID)

		if action, ok := d.GetOk(cisEdgeFunctionsTriggerActionName); ok {
			opt.SetScript(action.(string))
		}
		pattern := d.Get(cisEdgeFunctionsTriggerPattern).(string)
		opt.SetPattern(pattern)

		_, _, err := cisClient.UpdateEdgeFunctionsTrigger(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error updating edge function trigger route : %v", err)
		}
	}
	return ResourceIBMCISEdgeFunctionsTriggerRead(d, meta)
}

func ResourceIBMCISEdgeFunctionsTriggerRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	routeID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsTriggerOptions(routeID)
	result, resp, err := cisClient.GetEdgeFunctionsTrigger(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error: %v", resp)
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisEdgeFunctionsTriggerID, routeID)
	d.Set(cisEdgeFunctionsTriggerActionName, result.Result.Script)
	d.Set(cisEdgeFunctionsTriggerPattern, result.Result.Pattern)
	d.Set(cisEdgeFunctionsTriggerRequestLimitFailOpen, result.Result.RequestLimitFailOpen)
	return nil
}

func ResourceIBMCISEdgeFunctionsTriggerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return false, err
	}

	routeID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsTriggerOptions(routeID)
	_, response, err := cisClient.GetEdgeFunctionsTrigger(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Edge functions trigger route is not found")
			return false, nil
		}
		return false, fmt.Errorf("[ERROR] Error: %v", response)
	}
	return true, nil
}

func ResourceIBMCISEdgeFunctionsTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error in creating CIS object")
	}

	routeID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewDeleteEdgeFunctionsTriggerOptions(routeID)
	_, response, err := cisClient.DeleteEdgeFunctionsTrigger(opt)
	if err != nil {
		return fmt.Errorf("[ERROR] Error in edge function trigger route deletion: %v", response)
	}
	return nil
}
