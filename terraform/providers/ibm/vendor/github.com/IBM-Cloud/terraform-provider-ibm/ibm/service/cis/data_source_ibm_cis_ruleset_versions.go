// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package cis

import (
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	CISRulesetsVersionOutput = "ruleset_versions"
)

func DataSourceIBMCISRulesetVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataIBMCISRulesetVersionsRead,
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeDataSourceValidator(
					"ibm_cis_ruleset_versions",
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
				Required:    true,
				Description: "ID of the Ruleset",
			},
			CISRulesetsVersion: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ruleset version",
			},
			CISRulesetsVersionOutput: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for response information.",
				Elem:        CISResponseObject,
			},
			CISRulesetsObjectOutput: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Container for response information.",
				Elem:        CISResponseObject,
			},
		},
	}
}

func DataSourceIBMCISRulesetVersionsValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)

	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "resource_instance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})

	IBMCISRulesetValidator := validate.ResourceValidator{
		ResourceName: "ibm_cis_ruleset_versions",
		Schema:       validateSchema}
	return &IBMCISRulesetValidator
}

func dataIBMCISRulesetVersionsRead(d *schema.ResourceData, meta interface{}) error {
	sess, err := meta.(conns.ClientSession).CisRulesetsSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneId := d.Get(cisDomainID).(string)
	ruleset_version := d.Get(CISRulesetsVersion).(string)
	rulesetId := d.Get(CISRulesetsId).(string)

	sess.Crn = core.StringPtr(crn)

	if zoneId != "" {
		sess.ZoneIdentifier = core.StringPtr(zoneId)

		if ruleset_version != "" {
			opt := sess.NewGetZoneRulesetVersionOptions(rulesetId, ruleset_version)
			result, resp, err := sess.GetZoneRulesetVersion(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}
			rulesetObj := flattenCISRulesets(*result.Result)

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsObjectOutput, rulesetObj)
			d.Set(cisID, crn)

		} else {
			opt := sess.NewGetZoneRulesetVersionsOptions(rulesetId)
			result, resp, err := sess.GetZoneRulesetVersions(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}

			rulesetList := make([]map[string]interface{}, 0)
			for _, rulesetObj := range result.Result {
				rulesetOutput := map[string]interface{}{}
				rulesetOutput[CISRulesetsDescription] = *rulesetObj.Description
				rulesetOutput[CISRulesetsKind] = *rulesetObj.Kind
				rulesetOutput[CISRulesetsName] = *rulesetObj.Name
				rulesetOutput[CISRulesetsPhase] = *rulesetObj.Phase
				rulesetOutput[CISRulesetsLastUpdatedAt] = *rulesetObj.LastUpdated
				rulesetOutput[CISRulesetsVersion] = *rulesetObj.Version

				if rulesetOutput[CISRulesetsPhase] == "http_request_firewall_managed" {
					rulesetList = append(rulesetList, rulesetOutput)
				}
			}

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsVersionOutput, rulesetList)
			d.Set(cisID, crn)
		}

	} else {

		if ruleset_version != "" {
			opt := sess.NewGetInstanceRulesetVersionOptions(rulesetId, ruleset_version)
			result, resp, err := sess.GetInstanceRulesetVersion(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}

			rulesetObj := flattenCISRulesets(*result.Result)

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsObjectOutput, rulesetObj)
			d.Set(cisID, crn)

		} else {
			opt := sess.NewGetInstanceRulesetVersionsOptions(rulesetId)
			result, resp, err := sess.GetInstanceRulesetVersions(opt)
			if err != nil {
				log.Printf("[WARN] List all Instance rulesets failed: %v\n", resp)
				return err
			}

			rulesetList := make([]map[string]interface{}, 0)
			for _, rulesetObj := range result.Result {
				rulesetOutput := map[string]interface{}{}
				rulesetOutput[CISRulesetsDescription] = *rulesetObj.Description
				rulesetOutput[CISRulesetsKind] = *rulesetObj.Kind
				rulesetOutput[CISRulesetsName] = *rulesetObj.Name
				rulesetOutput[CISRulesetsPhase] = *rulesetObj.Phase
				rulesetOutput[CISRulesetsLastUpdatedAt] = *rulesetObj.LastUpdated
				rulesetOutput[CISRulesetsVersion] = *rulesetObj.Version
				rulesetOutput[CISRulesetsId] = *rulesetObj.ID

				if rulesetOutput[CISRulesetsPhase] == "http_request_firewall_managed" {
					rulesetList = append(rulesetList, rulesetOutput)
				}
			}

			d.SetId(dataSourceCISRulesetsCheckID(d))
			d.Set(CISRulesetsVersionOutput, rulesetList)
			d.Set(cisID, crn)
		}
	}

	return nil
}
