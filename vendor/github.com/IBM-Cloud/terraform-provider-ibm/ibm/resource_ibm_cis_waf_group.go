// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISWAFGroup                = "ibm_cis_waf_group"
	cisWAFGroupID                 = "group_id"
	cisWAFGroupPackageID          = "package_id"
	cisWAFGroupMode               = "mode"
	cisWAFGroupName               = "name"
	cisWAFGroupRulesCount         = "rules_count"
	cisWAFGroupModifiedRulesCount = "modified_rules_count"
	cisWAFGroupDesc               = "description"
)

func resourceIBMCISWAFGroup() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISWAFGroupUpdate,
		Read:     resourceIBMCISWAFGroupRead,
		Update:   resourceIBMCISWAFGroupUpdate,
		Delete:   resourceIBMCISWAFGroupDelete,
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
			cisWAFGroupPackageID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "WAF Rule package id",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFGroupID: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "WAF Rule group id",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFGroupMode: {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "WAF Rule group mode on/off",
				ValidateFunc: InvokeValidator(ibmCISWAFGroup, cisWAFGroupMode),
			},
			cisWAFGroupName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF Rule group name",
			},
			cisWAFGroupDesc: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF Rule group description",
			},
			cisWAFGroupRulesCount: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF Rule group rules count",
			},
			cisWAFGroupModifiedRulesCount: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF Rule group modified rules count",
			},
		},
	}
}

func resourceIBMCISWAFGroupValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	mode := "on, off"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisWAFGroupMode,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              mode})
	ibmCISWAFGroupValidator := ResourceValidator{ResourceName: ibmCISWAFGroup, Schema: validateSchema}
	return &ibmCISWAFGroupValidator
}

func resourceIBMCISWAFGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisWAFGroupClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	packageID, _, _, _ := convertTfToCisThreeVar(d.Get(cisWAFGroupPackageID).(string))
	groupID := d.Get(cisWAFGroupID).(string)

	if d.HasChange(cisWAFGroupMode) {
		mode := d.Get(cisWAFGroupMode).(string)
		opt := cisClient.NewUpdateWafRuleGroupOptions(packageID, groupID)
		opt.SetMode(mode)
		_, response, err := cisClient.UpdateWafRuleGroup(opt)
		if err != nil {
			log.Printf("Update waf rule group mode failed: %v", response)
			return err
		}
	}
	d.SetId(convertCisToTfFourVar(groupID, packageID, zoneID, crn))
	return resourceIBMCISWAFGroupRead(d, meta)
}

func resourceIBMCISWAFGroupRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisWAFGroupClientSession()
	if err != nil {
		return err
	}
	groupID, packageID, zoneID, crn, err := convertTfToCisFourVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	opt := cisClient.NewGetWafRuleGroupOptions(packageID, groupID)
	result, response, err := cisClient.GetWafRuleGroup(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("WAF group is not found!")
			d.SetId("")
			return nil
		}
		log.Printf("Get waf rule group setting failed: %v", response)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisWAFGroupID, groupID)
	d.Set(cisWAFGroupPackageID, result.Result.PackageID)
	d.Set(cisWAFGroupMode, result.Result.Mode)
	d.Set(cisWAFGroupName, result.Result.Name)
	d.Set(cisWAFGroupDesc, result.Result.Description)
	d.Set(cisWAFGroupModifiedRulesCount, result.Result.ModifiedRulesCount)
	d.Set(cisWAFGroupRulesCount, result.Result.RulesCount)
	return nil
}

func resourceIBMCISWAFGroupDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS WAF Group resource
	d.SetId("")
	return nil
}
