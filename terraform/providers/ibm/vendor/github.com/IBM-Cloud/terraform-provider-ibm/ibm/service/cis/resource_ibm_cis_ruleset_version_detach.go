// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"fmt"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMCISRulesetVersionDetach() *schema.Resource {
	return &schema.Resource{
		Create: ResourceIBMCISRulesetVersionDetachCreate,
		Read:   ResourceIBMCISRulesetVersionDetachRead,
		Delete: ResourceIBMCISRulesetVersionDetachDelete,
		Update: ResourceIBMCISRulesetVersionDetachRead,
		// Exists:   ResourceIBMCISRulesetVersionDetachExists,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_ruleset_version_detach",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Optional:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			CISRulesetsId: {
				Type:        schema.TypeString,
				Description: "Associated Ruleset ID",
				Required:    true,
			},
			CISRulesetsVersion: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Ruleset version",
			},
		},
	}
}
func ResourceIBMCISRulesetVersionDetachValidator() *validate.ResourceValidator {
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
		ResourceName: "ibm_cis_ruleset_version_detach",
		Schema:       validateSchema}
	return &ibmCISRulesetValidator
}

func ResourceIBMCISRulesetVersionDetachCreate(d *schema.ResourceData, meta interface{}) error {

	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return fmt.Errorf("[ERROR] Error while getting the CisRulesetsSession %s", err)
	}

	crn := d.Get(cisID).(string)
	sess.Crn = core.StringPtr(crn)

	zoneId := d.Get(cisDomainID).(string)
	rulesetId := d.Get(CISRulesetsId).(string)
	ruleset_version := d.Get(CISRulesetsVersion).(string)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)
		opt := sess.NewDeleteZoneRulesetVersionOptions(rulesetId, ruleset_version)
		res, err := sess.DeleteZoneRulesetVersion(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching the zone ruleset version %s:%s", err, res)
		}
	} else {
		opt := sess.NewDeleteInstanceRulesetVersionOptions(rulesetId, ruleset_version)
		res, err := sess.DeleteInstanceRulesetVersion(opt)
		if err != nil {
			return fmt.Errorf("[ERROR] Error detaching the Instance ruleset version %s:%s", err, res)
		}
	}

	d.SetId(time.Now().UTC().String())
	return ResourceIBMCISRulesetVersionDetachRead(d, meta)
}

func ResourceIBMCISRulesetVersionDetachRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func ResourceIBMCISRulesetVersionDetachDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
