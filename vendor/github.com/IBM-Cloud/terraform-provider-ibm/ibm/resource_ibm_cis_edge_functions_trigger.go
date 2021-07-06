// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	cisEdgeFunctionsTriggerID                   = "trigger_id"
	cisEdgeFunctionsTriggerPattern              = "pattern_url"
	cisEdgeFunctionsTriggerActionName           = "action_name"
	cisEdgeFunctionsTriggerRequestLimitFailOpen = "request_limit_fail_open"
)

func resourceIBMCISEdgeFunctionsTrigger() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISEdgeFunctionsTriggerCreate,
		Read:     resourceIBMCISEdgeFunctionsTriggerRead,
		Update:   resourceIBMCISEdgeFunctionsTriggerUpdate,
		Delete:   resourceIBMCISEdgeFunctionsTriggerDelete,
		Exists:   resourceIBMCISEdgeFunctionsTriggerExists,
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

func resourceIBMCISEdgeFunctionsTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
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
		return fmt.Errorf("Error creating edge function trigger route : %v", err)
	}
	d.SetId(convertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	return resourceIBMCISEdgeFunctionsTriggerRead(d, meta)
}

func resourceIBMCISEdgeFunctionsTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	routeID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
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
			return fmt.Errorf("Error updating edge function trigger route : %v", err)
		}
	}
	return resourceIBMCISEdgeFunctionsTriggerRead(d, meta)
}

func resourceIBMCISEdgeFunctionsTriggerRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return err
	}

	routeID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsTriggerOptions(routeID)
	result, resp, err := cisClient.GetEdgeFunctionsTrigger(opt)
	if err != nil {
		return fmt.Errorf("Error: %v", resp)
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisEdgeFunctionsTriggerID, routeID)
	d.Set(cisEdgeFunctionsTriggerActionName, result.Result.Script)
	d.Set(cisEdgeFunctionsTriggerPattern, result.Result.Pattern)
	d.Set(cisEdgeFunctionsTriggerRequestLimitFailOpen, result.Result.RequestLimitFailOpen)
	return nil
}

func resourceIBMCISEdgeFunctionsTriggerExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return false, err
	}

	routeID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewGetEdgeFunctionsTriggerOptions(routeID)
	_, response, err := cisClient.GetEdgeFunctionsTrigger(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("Edge functions trigger route is not found")
			return false, nil
		}
		return false, fmt.Errorf("Error: %v", response)
	}
	return true, nil
}

func resourceIBMCISEdgeFunctionsTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisEdgeFunctionClientSession()
	if err != nil {
		return fmt.Errorf("Error in creating CIS object")
	}

	routeID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	opt := cisClient.NewDeleteEdgeFunctionsTriggerOptions(routeID)
	_, response, err := cisClient.DeleteEdgeFunctionsTrigger(opt)
	if err != nil {
		return fmt.Errorf("Error in edge function trigger route deletion: %v", response)
	}
	return nil
}
