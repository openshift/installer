// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func resourceIBMCISWAFPackage() *schema.Resource {
	return &schema.Resource{
		Create:   resourceIBMCISWAFPackageUpdate,
		Read:     resourceIBMCISWAFPackageRead,
		Update:   resourceIBMCISWAFPackageUpdate,
		Delete:   resourceIBMCISWAFPackageDelete,
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
				ValidateFunc: InvokeValidator(
					ibmCISWAFPackage, cisWAFPackageSensitivity),
			},
			cisWAFPackageActionMode: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "WAF pakcage action mode",
				ValidateFunc: InvokeValidator(
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

func resourceIBMCISWAFPackageValidator() *ResourceValidator {

	validateSchema := make([]ValidateSchema, 1)
	sesitivity := "high, medium, low, off"
	actionMode := "simulate, block, challenge"

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisWAFPackageSensitivity,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              sesitivity})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisWAFPackageActionMode,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              actionMode})
	ibmCISWAFPackageValidator := ResourceValidator{ResourceName: ibmCISWAFPackage, Schema: validateSchema}
	return &ibmCISWAFPackageValidator
}

func resourceIBMCISWAFPackageUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisWAFPackageClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	packageID, _, _, err := convertTfToCisThreeVar(d.Get(cisWAFPackageID).(string))
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
		d.SetId(convertCisToTfThreeVar(*result.Result.ID, zoneID, crn))
	}

	return resourceIBMCISWAFPackageRead(d, meta)
}

func resourceIBMCISWAFPackageRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisWAFPackageClientSession()
	if err != nil {
		return err
	}
	packageID, zoneID, crn, err := convertTfToCisThreeVar(d.Id())
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

func resourceIBMCISWAFPackageDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS WAF Package resource
	d.SetId("")
	return nil
}
