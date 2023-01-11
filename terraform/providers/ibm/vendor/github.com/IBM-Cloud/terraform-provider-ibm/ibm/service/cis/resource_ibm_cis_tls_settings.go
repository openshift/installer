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
	ibmCISTLSSettings           = "ibm_cis_tls_settings"
	cisTLSSettingsUniversalSSL  = "universal_ssl"
	cisTLSSettingsTLS12Only     = "tls_1_2_only"
	cisTLSSettingsTLS13         = "tls_1_3"
	cisTLSSettingsMinTLSVersion = "min_tls_version"
)

func ResourceIBMCISTLSSettings() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			cisID: {
				Type:        schema.TypeString,
				Description: "CIS instance crn",
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_cis_tls_settings",
					"cis_id"),
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
				ValidateFunc:     validate.InvokeValidator(ibmCISTLSSettings, cisTLSSettingsTLS13),
				DiffSuppressFunc: suppressTLS13Diff,
			},
			cisTLSSettingsMinTLSVersion: {
				Type:         schema.TypeString,
				Description:  "Minimum version of TLS required",
				Optional:     true,
				ValidateFunc: validate.InvokeValidator(ibmCISTLSSettings, cisTLSSettingsMinTLSVersion),
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

func ResourceIBMCISTLSSettingsValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)

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
			Identifier:                 cisTLSSettingsTLS13,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "on, off, zrt"})
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 cisTLSSettingsMinTLSVersion,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "1.1, 1.2, 1.3, 1.4"})
	ibmCISTLSSettingsResourceValidator := validate.ResourceValidator{
		ResourceName: ibmCISTLSSettings,
		Schema:       validateSchema}
	return &ibmCISTLSSettingsResourceValidator
}

func resourceCISTLSSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, _ := flex.ConvertTftoCisTwoVar(d.Get(cisDomainID).(string))
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
			cisClient, err := meta.(conns.ClientSession).CisDomainSettingsClientSession()
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
	d.SetId(flex.ConvertCisToTfTwoVar(zoneID, crn))
	return resourceCISTLSSettingsRead(d, meta)
}

func resourceCISTLSSettingsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(conns.ClientSession).CisSSLClientSession()
	if err != nil {
		return err
	}
	zoneID, crn, _ := flex.ConvertTftoCisTwoVar(d.Id())
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
	minTLSClient, err := meta.(conns.ClientSession).CisDomainSettingsClientSession()
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
	if old == "zrt" && new == "on" {
		return true
	}
	return false
}
