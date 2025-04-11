// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMCISRulesetEntryPointVersion() *schema.Resource {
	return &schema.Resource{
		Read:     ResourceIBMCISRulesetEntryPointVersionRead,
		Create:   ResourceIBMCISRulesetEntryPointVersionUpdate,
		Update:   ResourceIBMCISRulesetEntryPointVersionUpdate,
		Delete:   ResourceIBMCISRulesetEntryPointVersionDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Optional:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			CISRulesetPhase: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Ruleset phase",
			},
			CISRulesetsEntryPointOutput: {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Container for response information.",
				Elem:        CISResourceResponseObject,
			},
		},
	}
}
func ResourceIBMCISRulesetEntryPointVersionValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	ibmCISRulesetValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_ruleset_entrypoint_version",
		Schema:       validateSchema}
	return &ibmCISRulesetValidator
}

func ResourceIBMCISRulesetEntryPointVersionRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisRulesetsSession %s", err)
	}

	ruleset_phase, zoneId, crn, err := flex.ConvertTfToCisThreeVar(d.Id())
	if err != nil {
		return fmt.Errorf("[ERROR] Error while ConvertTftoCisThreeVar %s", err)
	}
	sess.Crn = core.StringPtr(crn)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)
		opt := sess.Clone().NewGetZoneEntrypointRulesetOptions(ruleset_phase)
		result, resp, err := sess.GetZoneEntrypointRuleset(opt)
		if err != nil {
			return fmt.Errorf("[WARN] Get zone ruleset failed: %v", resp)
		}
		rulesetObj := flattenCISRulesets(*result.Result)

		d.Set(CISRulesetsEntryPointOutput, rulesetObj)
		d.Set(cisDomainID, zoneId)
		d.Set(cisID, crn)
		d.Set(CISRulesetPhase, ruleset_phase)

	} else {
		opt := sess.NewGetInstanceEntrypointRulesetOptions(ruleset_phase)
		result, resp, err := sess.GetInstanceEntrypointRuleset(opt)
		if err != nil {
			return fmt.Errorf("[WARN] Get zone ruleset failed: %v", resp)
		}
		rulesetObj := flattenCISRulesets(*result.Result)

		d.Set(CISRulesetsEntryPointOutput, rulesetObj)
		d.Set(cisDomainID, zoneId)
		d.Set(cisID, crn)
		d.Set(CISRulesetPhase, ruleset_phase)

	}

	return nil
}

func ResourceIBMCISRulesetEntryPointVersionUpdate(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisRulesetsSession %s", err)
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	zoneId := d.Get(cisDomainID).(string)
	ruleset_phase := d.Get(CISRulesetPhase).(string)

	if zoneId != "" {
		sess.ZoneIdentifier = &zoneId

		opt := sess.NewUpdateZoneEntrypointRulesetOptions(ruleset_phase)

		cis_ruleset_object := d.Get(CISRulesetsObjectOutput)

		rulesetsObject := cis_ruleset_object.(*schema.Set).List()[0].(map[string]interface{})
		opt.SetDescription(rulesetsObject[CISRulesetsDescription].(string))
		opt.SetName(rulesetsObject[CISRulesetsName].(string))

		rulesObj := expandCISRules(rulesetsObject[CISRulesetsRules])
		opt.SetRules(rulesObj)

		result, resp, err := sess.UpdateZoneEntrypointRuleset(opt)
		if err != nil || result == nil {
			return fmt.Errorf("[ERROR] Error while Update Zone Entrypoint Rulesets %s %s", err, resp)
		}

	} else {
		opt := sess.NewUpdateInstanceEntrypointRulesetOptions(ruleset_phase)

		rulesetsObject := d.Get(CISRulesetsObjectOutput).([]interface{})[0].(map[string]interface{})
		opt.SetDescription(rulesetsObject[CISRulesetsDescription].(string))
		opt.SetName(rulesetsObject[CISRulesetsName].(string))

		rulesObj := expandCISRules(rulesetsObject[CISRulesetsRules])
		opt.SetRules(rulesObj)

		result, resp, err := sess.UpdateInstanceEntrypointRuleset(opt)
		if err != nil || result == nil {
			return fmt.Errorf("[ERROR] Error while Update Entrypoint Rulesets %s %s", err, resp)
		}

	}
	d.SetId(dataSourceCISRulesetsEPCheckID(d))
	return ResourceIBMCISRulesetEntryPointVersionRead(d, meta)
}

func ResourceIBMCISRulesetEntryPointVersionDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
