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
	ibmCISWAFPackage           = "ibm_cis_waf_package"
	cisWAFPackageID            = "package_id"
	cisWAFPackageName          = "name"
	cisWAFPackageDescription   = "description"
	cisWAFPackageDetectionMode = "detection_mode"
	cisWAFPackageSensitivity   = "sensitivity"
	cisWAFPackageActionMode    = "action_mode"
)

func ResourceIBMCISWAFPackage() *schema.Resource {
	return &schema.Resource{
		Create:   ResourceIBMCISWAFPackageUpdate,
		Read:     ResourceIBMCISWAFPackageRead,
		Update:   ResourceIBMCISWAFPackageUpdate,
		Delete:   ResourceIBMCISWAFPackageDelete,
		Importer: &schema.ResourceImporter{},
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CIS Intance CRN",
				ValidateFunc: validate.InvokeValidator("ibm_cis_waf_package",
					"cis_id"),
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "CIS Domain ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFPackageID: {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "WAF pakcage ID",
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisWAFPackageName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF pakcage name",
			},
			cisWAFPackageDetectionMode: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF pakcage detection mode",
			},
			cisWAFPackageSensitivity: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "WAF pakcage sensitivity",
				ValidateFunc: validate.InvokeValidator(
					ibmCISWAFPackage, cisWAFPackageSensitivity),
			},
			cisWAFPackageActionMode: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "WAF pakcage action mode",
				ValidateFunc: validate.InvokeValidator(
					ibmCISWAFPackage, cisWAFPackageActionMode),
			},
			cisWAFPackageDescription: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "WAF package description",
			},
		},
	}
}

func ResourceIBMCISWAFPackageValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	sesitivity := "high, medium, low, off"
	actionMode := "simulate, block, challenge"

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
			Identifier:                 cisWAFPackageSensitivity,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              sesitivity})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisWAFPackageActionMode,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              actionMode})
	ibmCISWAFPackageValidator := validate.ResourceValidator{ResourceName: ibmCISWAFPackage, Schema: validateSchema}
	return &ibmCISWAFPackageValidator
}

func ResourceIBMCISWAFPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFPackageClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
	packageID, _, _, _ := flex.ConvertTfToCisThreeVar(d.Get(cisWAFPackageID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	if d.HasChange(cisWAFPackageSensitivity) ||
		d.HasChange(cisWAFPackageActionMode) {
		opt := cisClient.NewUpdateWafPackageOptions(packageID)
		if v, ok := d.GetOk(cisWAFPackageSensitivity); ok {
			opt.SetSensitivity(v.(string))
		}
		if v, ok := d.GetOk(cisWAFPackageActionMode); ok {
			opt.SetActionMode(v.(string))
		}
		result, response, err := cisClient.UpdateWafPackage(opt)
		if err != nil {
			log.Printf("Update waf package setting failed: %v", response)
			return err
		}
		d.SetId(flex.ConvertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	}

	return ResourceIBMCISWAFPackageRead(d, meta)
}

func ResourceIBMCISWAFPackageRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisWAFPackageClientSession()
	if err != nil {
		return err
	}
	packageID, zoneID, crn, _ := flex.ConvertTfToCisThreeVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)
	opt := cisClient.NewGetWafPackageOptions(packageID)
	result, response, err := cisClient.GetWafPackage(opt)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			log.Printf("WAF package is not found!")
			d.SetId("")
			return nil
		}
		log.Printf("Get waf package setting failed: %v", response)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisWAFPackageID, result.Result.ID)
	d.Set(cisWAFPackageName, result.Result.Name)
	d.Set(cisWAFPackageDetectionMode, result.Result.DetectionMode)
	d.Set(cisWAFPackageActionMode, result.Result.ActionMode)
	d.Set(cisWAFPackageSensitivity, result.Result.Sensitivity)
	d.Set(cisWAFPackageDescription, result.Result.Description)
	return nil
}

func ResourceIBMCISWAFPackageDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS WAF Package resource
	d.SetId("")
	return nil
}
