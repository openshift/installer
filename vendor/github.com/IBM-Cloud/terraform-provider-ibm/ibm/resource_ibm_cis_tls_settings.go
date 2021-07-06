// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISTLSSettings           = "ibm_cis_tls_settings"
	cisTLSSettingsUniversalSSL  = "universal_ssl"
	cisTLSSettingsTLS12Only     = "tls_1_2_only"
	cisTLSSettingsTLS13         = "tls_1_3"
	cisTLSSettingsMinTLSVersion = "min_tls_version"
)

func resourceIBMCISTLSSettings() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
			},
			cisDomainID: {
				Type:             schema.TypeString,
				Description:      "Associated CIS domain",
				Required:         true,
				DiffSuppressFunc: suppressDomainIDDiff,
			},
			cisTLSSettingsUniversalSSL: {
				Type:        schema.TypeBool,
				Description: "Universal SSL setting",
				Optional:    true,
				Computed:    true,
			},
			cisTLSSettingsTLS13: {
				Type:             schema.TypeString,
				Description:      "TLS 1.3 setting",
				Optional:         true,
				Computed:         true,
				ValidateFunc:     InvokeValidator(ibmCISTLSSettings, cisTLSSettingsTLS13),
				DiffSuppressFunc: suppressTLS13Diff,
			},
			cisTLSSettingsMinTLSVersion: {
				Type:         schema.TypeString,
				Description:  "Minimum version of TLS required",
				Optional:     true,
				ValidateFunc: InvokeValidator(ibmCISTLSSettings, cisTLSSettingsMinTLSVersion),
				Default:      "1.1",
			},
		},
		Create:   resourceCISTLSSettingsUpdate,
		Read:     resourceCISTLSSettingsRead,
		Update:   resourceCISTLSSettingsUpdate,
		Delete:   resourceCISTLSSettingsDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func resourceIBMCISTLSSettingsValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisTLSSettingsTLS13,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "on, off, zrt"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisTLSSettingsMinTLSVersion,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "1.1, 1.2, 1.3, 1.4"})
	ibmCISTLSSettingsResourceValidator := ResourceValidator{
		ResourceName: ibmCISTLSSettings,
		Schema:       validateSchema}
	return &ibmCISTLSSettingsResourceValidator
}

func resourceCISTLSSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	if d.HasChange(cisTLSSettingsTLS12Only) ||
		d.HasChange(cisTLSSettingsTLS13) ||
		d.HasChange(cisTLSSettingsUniversalSSL) ||
		d.HasChange(cisTLSSettingsMinTLSVersion) {

		// TLS 1.3 setting
		if tls13, ok := d.GetOk(cisTLSSettingsTLS13); ok {
			opt := cisClient.NewChangeTls13SettingOptions()
			opt.SetValue(tls13.(string))
			_, resp, err := cisClient.ChangeTls13Setting(opt)
			if err != nil {
				log.Printf("Update TLS 1.3 setting Failed : %v\n", resp)
				return err
			}
		}

		// Universal SSL setting
		if universalSSL, ok := d.GetOkExists(cisTLSSettingsUniversalSSL); ok {
			opt := cisClient.NewChangeUniversalCertificateSettingOptions()
			opt.SetEnabled(universalSSL.(bool))
			resp, err := cisClient.ChangeUniversalCertificateSetting(opt)
			if err != nil {
				log.Printf("Update universal ssl setting Failed : %v\n", resp)
				return err
			}
		}

		// Minimum TLS version
		if minTLSVer, ok := d.GetOk(cisTLSSettingsMinTLSVersion); ok {
			cisClient, err := meta.(ClientSession).CisDomainSettingsClientSession()
			if err != nil {
				return err
			}
			cisClient.Crn = core.StringPtr(crn)
			cisClient.ZoneIdentifier = core.StringPtr(zoneID)
			opt := cisClient.NewUpdateMinTlsVersionOptions()
			opt.SetValue(minTLSVer.(string))
			_, resp, err := cisClient.UpdateMinTlsVersion(opt)
			if err != nil {
				log.Printf("Update minimum TLS version setting Failed : %v\n", resp)
				return err
			}
		}
	}
	d.SetId(convertCisToTfTwoVar(zoneID, crn))
	return resourceCISTLSSettingsRead(d, meta)
}

func resourceCISTLSSettingsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	zoneID, crn, _ := convertTftoCisTwoVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	// TLS 1.3 setting
	tls13Result, resp, err := cisClient.GetTls13Setting(cisClient.NewGetTls13SettingOptions())
	if err != nil {
		log.Printf("Get TLS 1.3 setting failed : %v\n", resp)
		return err
	}

	// Universal SSL setting
	universalSSLResult, resp, err := cisClient.GetUniversalCertificateSetting(
		cisClient.NewGetUniversalCertificateSettingOptions())
	if err != nil {
		log.Printf("Update TLS 1.3 setting failed : %v\n", resp)
		return err
	}

	// Minumum TLS version setting
	minTLSClient, err := meta.(ClientSession).CisDomainSettingsClientSession()
	if err != nil {
		return err
	}
	minTLSClient.Crn = core.StringPtr(crn)
	minTLSClient.ZoneIdentifier = core.StringPtr(zoneID)
	minTLSVerResult, resp, err := minTLSClient.GetMinTlsVersion(
		minTLSClient.NewGetMinTlsVersionOptions())
	if err != nil {
		log.Printf("Min TLS Version setting get request failed : %v", resp)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisTLSSettingsTLS13, tls13Result.Result.Value)
	d.Set(cisTLSSettingsUniversalSSL, universalSSLResult.Result.Enabled)
	d.Set(cisTLSSettingsMinTLSVersion, minTLSVerResult.Result.Value)
	return nil
}

func resourceCISTLSSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS resource
	d.SetId("")
	return nil
}

func suppressTLS13Diff(k, old, new string, d *schema.ResourceData) bool {
	// if we enable TLS 1.3, it gives zrt in output.
	if "zrt" == old && new == "on" {
		return true
	}
	return false
}
