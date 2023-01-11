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
	ibmCISWAFRule          = "ibm_cis_waf_rule"
	cisWAFRuleID           = "rule_id"
	cisWAFRuleDesc         = "description"
	cisWAFRulePriority     = "priority"
	cisWAFRulePackageID    = "package_id"
	cisWAFRuleGroup        = "group"
	cisWAFRuleGroupID      = "id"
	cisWAFRuleGroupName    = "name"
	cisWAFRuleMode         = "mode"
	cisWAFRuleModeOn       = "on"
	cisWAFRuleModeOff      = "off"
	cisWAFRuleAllowedModes = "allowed_modes"
)

func ResourceIBMCISWAFRule() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISWAFRuleUpdate,
		Read:     ResourceIBMCISWAFRuleRead,
		Update:   ResourceIBMCISWAFRuleUpdate,
		Delete:   ResourceIBMCISWAFRuleDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator("ibm_cis_waf_rule",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFRuleID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIS WAF Rule id",
			},
			cisWAFRulePackageID: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CIS WAF Rule package id",
			},
			cisWAFRuleMode: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "CIS WAF Rule mode",
				ValidateFunc: validate.InvokeValidator(ibmCISWAFRule, cisWAFRuleMode),
			},
			cisWAFRuleDesc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIS WAF Rule descriptions",
			},
			cisWAFRulePriority: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "CIS WAF Rule Priority",
			},
			cisWAFRuleGroup: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "CIS WAF Rule group",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisWAFRuleGroupID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "waf rule group id",
						},
						cisWAFRuleGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "waf rule group name",
						},
					},
				},
			},
			cisWAFRuleAllowedModes: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "CIS WAF Rule allowed modes",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func ResourceIBMCISWAFRuleValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	modes := "on, off, default, disable, simulate, block, challenge"
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "cis_id",
			ValidateFunctionIdentifier: validate.ValidateCloudData,
			Type:                       validate.TypeString,
			CloudDataType:              "ResourceInstance",
			CloudDataRange:             []string{"service:internet-svcs"},
			Required:                   true})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisWAFRuleMode,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              modes})
	ibmCISWAFRuleValidator := validate.ResourceValidator{ResourceName: ibmCISWAFRule, Schema: validateSchema}
	return &ibmCISWAFRuleValidator
}

func ResourceIBMCISWAFRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFRuleClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	ruleID := d.Get(cisWAFRuleID).(string)
	packageID, _, _, _ := flex.ConvertTfToCisThreeVar(d.Get(cisWAFRulePackageID).(string))

	if d.HasChange(cisWAFRuleMode) {
		mode := d.Get(cisWAFRuleMode).(string)

		getOpt := cisClient.NewGetWafRuleOptions(packageID, ruleID)
		getResult, getResponse, err := cisClient.GetWafRule(getOpt)
		if err != nil {
			log.Printf("Get WAF rule setting failed: %v", getResponse)
			return err
		}
		getMode := *getResult.Result.Mode
		updateOpt := cisClient.NewUpdateWafRuleOptions(packageID, ruleID)

		// Mode differs based on OWASP and CIS
		if getMode == cisWAFRuleModeOn || getMode == cisWAFRuleModeOff {

			owaspOpt, _ := cisClient.NewWafRuleBodyOwasp(mode)
			updateOpt.SetOwasp(owaspOpt)

		} else {

			cisOpt, _ := cisClient.NewWafRuleBodyCis(mode)
			updateOpt.SetCis(cisOpt)

		}
		_, response, err := cisClient.UpdateWafRule(updateOpt)
		if err != nil {
			log.Printf("Update WAF rule setting failed: %v", response)
			return err
		}
	}

	d.SetId(flex.ConvertCisToTfFourVar(ruleID, packageID, zoneID, crn))
	return ResourceIBMCISWAFRuleRead(d, meta)
}

func ResourceIBMCISWAFRuleRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFRuleClientSession()
	if err != nil {
		return err
	}
	ruleID, packageID, zoneID, crn, _ := flex.ConvertTfToCisFourVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	opt := cisClient.NewGetWafRuleOptions(packageID, ruleID)
	result, response, err := cisClient.GetWafRule(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("WAF Rule is not found!")
			d.SetId("")
			return nil
		}
		log.Printf("Get waf rule setting failed: %v", response)
		return err
	}
	groups := []interface{}{}
	group := map[string]interface{}{}
	group[cisWAFRuleGroupID] = *result.Result.Group.ID
	group[cisWAFRuleGroupName] = *result.Result.Group.Name
	groups = append(groups, group)

	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisWAFRuleID, ruleID)
	d.Set(cisWAFRulePackageID, packageID)
	d.Set(cisWAFRuleDesc, *result.Result.Description)
	d.Set(cisWAFRulePriority, *result.Result.Priority)
	d.Set(cisWAFRuleGroup, groups)
	d.Set(cisWAFRuleMode, *result.Result.Mode)
	d.Set(cisWAFRuleAllowedModes, flex.FlattenStringList(result.Result.AllowedModes))
	return nil
}

func ResourceIBMCISWAFRuleDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS WAF rule resource
	d.SetId("")
	return nil
}
