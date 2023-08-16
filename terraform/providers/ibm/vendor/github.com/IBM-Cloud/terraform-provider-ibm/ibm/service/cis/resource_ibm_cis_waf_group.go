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
	ibmCISWAFGroup                = "ibm_cis_waf_group"
	cisWAFGroupID                 = "group_id"
	cisWAFGroupPackageID          = "package_id"
	cisWAFGroupMode               = "mode"
	cisWAFGroupName               = "name"
	cisWAFGroupRulesCount         = "rules_count"
	cisWAFGroupModifiedRulesCount = "modified_rules_count"
	cisWAFGroupDesc               = "description"
	cisWAFCheckMode               = "check_mode"
)

func ResourceIBMCISWAFGroup() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISWAFGroupUpdate,
		Read:     ResourceIBMCISWAFGroupRead,
		Update:   ResourceIBMCISWAFGroupUpdate,
		Delete:   ResourceIBMCISWAFGroupDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator("ibm_cis_waf_group",
					"cis_id"),
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
				ValidateFunc: validate.InvokeValidator(ibmCISWAFGroup, cisWAFGroupMode),
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
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "WAF Rule group rules count",
			},
			cisWAFGroupModifiedRulesCount: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "WAF Rule group modified rules count",
			},
			cisWAFCheckMode: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Check Mode before making a create/update request",
				Default:     false,
			},
		},
	}
}

func ResourceIBMCISWAFGroupValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	mode := "on, off"

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
			Identifier:                 cisWAFGroupMode,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              mode})
	ibmCISWAFGroupValidator := validate.ResourceValidator{ResourceName: ibmCISWAFGroup, Schema: validateSchema}
	return &ibmCISWAFGroupValidator
}

func ResourceIBMCISWAFGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFGroupClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	packageID, _, _, _ := flex.ConvertTfToCisThreeVar(d.Get(cisWAFGroupPackageID).(string))
	groupID := d.Get(cisWAFGroupID).(string)
	mode := d.Get(cisWAFGroupMode).(string)

	checkMode := d.Get(cisWAFCheckMode)

	if checkMode == true {
		opt := cisClient.NewGetWafRuleGroupOptions(packageID, groupID)
		result, _, error := cisClient.GetWafRuleGroup(opt)
		if err != nil {
			log.Printf("Get waf rule group setting failed: %v", error)
			return err
		}

		actualMode := *result.Result.Mode
		if actualMode == mode {
			d.SetId(flex.ConvertCisToTfFourVar(groupID, packageID, zoneID, crn))
			return ResourceIBMCISWAFGroupRead(d, meta)
		}
	}

	if d.HasChange(cisWAFGroupMode) {
		opt := cisClient.NewUpdateWafRuleGroupOptions(packageID, groupID)
		opt.SetMode(mode)
		_, response, err := cisClient.UpdateWafRuleGroup(opt)
		if err != nil {
			log.Printf("Update waf rule group mode failed: %v", response)
			return err
		}
	}
	d.SetId(flex.ConvertCisToTfFourVar(groupID, packageID, zoneID, crn))
	return ResourceIBMCISWAFGroupRead(d, meta)
}

func ResourceIBMCISWAFGroupRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFGroupClientSession()
	if err != nil {
		return err
	}
	groupID, packageID, zoneID, crn, _ := flex.ConvertTfToCisFourVar(d.Id())
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

func ResourceIBMCISWAFGroupDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS WAF Group resource
	d.SetId("")
	return nil
}
