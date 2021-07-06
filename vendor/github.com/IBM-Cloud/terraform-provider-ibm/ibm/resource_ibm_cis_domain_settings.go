// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISDomainSettings                             = "ibm_cis_domain_settings"
	cisDomainSettingsDNSSEC                          = "dnssec"
	cisDomainSettingsWAF                             = "waf"
	cisDomainSettingsSSL                             = "ssl"
	cisDomainSettingsCertificateStatus               = "certificate_status"
	cisDomainSettingsMinTLSVersion                   = "min_tls_version"
	cisDomainSettingsCNAMEFlattening                 = "cname_flattening"
	cisDomainSettingsOpportunisticEncryption         = "opportunistic_encryption"
	cisDomainSettingsAutomaticHTPSRewrites           = "automatic_https_rewrites"
	cisDomainSettingsAlwaysUseHTTPS                  = "always_use_https"
	cisDomainSettingsIPv6                            = "ipv6"
	cisDomainSettingsBrowserCheck                    = "browser_check"
	cisDomainSettingsHotlinkProtection               = "hotlink_protection"
	cisDomainSettingsHTTP2                           = "http2"
	cisDomainSettingsImageLoadOptimization           = "image_load_optimization"
	cisDomainSettingsImageSizeOptimization           = "image_size_optimization"
	cisDomainSettingsIPGeoLocation                   = "ip_geolocation"
	cisDomainSettingsOriginErrorPagePassThru         = "origin_error_page_pass_thru"
	cisDomainSettingsBrotli                          = "brotli"
	cisDomainSettingsPseudoIPv4                      = "pseudo_ipv4"
	cisDomainSettingsPrefetchPreload                 = "prefetch_preload"
	cisDomainSettingsResponseBuffering               = "response_buffering"
	cisDomainSettingsScriptLoadOptimisation          = "script_load_optimization"
	cisDomainSettingsServerSideExclude               = "server_side_exclude"
	cisDomainSettingsTLSClientAuth                   = "tls_client_auth"
	cisDomainSettingsTrueClientIPHeader              = "true_client_ip_header"
	cisDomainSettingsWebSockets                      = "websockets"
	cisDomainSettingsChallengeTTL                    = "challenge_ttl"
	cisDomainSettingsMinify                          = "minify"
	cisDomainSettingsMinifyCSS                       = "css"
	cisDomainSettingsMinifyHTML                      = "html"
	cisDomainSettingsMinifyJS                        = "js"
	cisDomainSettingsSecurityHeader                  = "security_header"
	cisDomainSettingsSecurityHeaderEnabled           = "enabled"
	cisDomainSettingsSecurityHeaderMaxAge            = "max_age"
	cisDomainSettingsSecurityHeaderIncludeSubdomains = "include_subdomains"
	cisDomainSettingsSecurityHeaderNoSniff           = "nosniff"
	cisDomainSettingsMobileRedirect                  = "mobile_redirect"
	cisDomainSettingsMobileRedirectStatus            = "status"
	cisDomainSettingsMobileRedirectMobileSubdomain   = "mobile_subdomain"
	cisDomainSettingsMobileRedirectStripURI          = "strip_uri"
	cisDomainSettingsMaxUpload                       = "max_upload"
	cisDomainSettingsCipher                          = "cipher"
	cisDomainSettingsONOFFValidatorID                = "on_off"
	cisDomainSettingsActiveDisableValidatorID        = "active_disable"
	cisDomainSettingsSSLSettingValidatorID           = "ssl_setting"
	cisDomainSettingsTLSVersionValidatorID           = "tls_version"
	cisDomainSettingsCNAMEFlattenValidatorID         = "cname_flatten"
	cisDomainSettingsImgSizeOptimizeValidatorID      = "img_size_optimize"
	cisDomainSettingsPseudoIPv4ValidatorID           = "psuedo_ipv4"
	cisDomainSettingsChallengeTTLValidatorID         = "challenge_ttl"
	cisDomainSettingsMaxUploadValidatorID            = "max_upload"
	cisDomainSettingsCipherValidatorID               = "cipher"
)

func resourceIBMCISSettings() *schema.Resource {
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
			cisDomainSettingsDNSSEC: {
				Type:        schema.TypeString,
				Description: "DNS Sec setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsActiveDisableValidatorID),
			},
			cisDomainSettingsWAF: {
				Type:        schema.TypeString,
				Description: "WAF setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsSSL: {
				Type:        schema.TypeString,
				Description: "SSL/TLS setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsSSLSettingValidatorID),
			},
			cisDomainSettingsCertificateStatus: {
				Type:        schema.TypeString,
				Description: "Certificate status",
				Computed:    true,
				Deprecated:  "This field is deprecated",
			},
			cisDomainSettingsMinTLSVersion: {
				Type:        schema.TypeString,
				Description: "Minimum version of TLS required",
				Optional:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsTLSVersionValidatorID),
				Default: "1.1",
			},
			cisDomainSettingsCNAMEFlattening: {
				Type:        schema.TypeString,
				Description: "cname_flattening setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsCNAMEFlattenValidatorID),
			},
			cisDomainSettingsOpportunisticEncryption: {
				Type:        schema.TypeString,
				Description: "opportunistic_encryption setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsAutomaticHTPSRewrites: {
				Type:        schema.TypeString,
				Description: "automatic_https_rewrites setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsAlwaysUseHTTPS: {
				Type:        schema.TypeString,
				Description: "always_use_https setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsIPv6: {
				Type:        schema.TypeString,
				Description: "ipv6 setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsBrowserCheck: {
				Type:        schema.TypeString,
				Description: "browser_check setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsHotlinkProtection: {
				Type:        schema.TypeString,
				Description: "hotlink_protection setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsHTTP2: {
				Type:        schema.TypeString,
				Description: "http2 setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsImageLoadOptimization: {
				Type:        schema.TypeString,
				Description: "image_load_optimization setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsImageSizeOptimization: {
				Type:        schema.TypeString,
				Description: "image_size_optimization setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsImgSizeOptimizeValidatorID),
			},
			cisDomainSettingsIPGeoLocation: {
				Type:        schema.TypeString,
				Description: "ip_geolocation setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsOriginErrorPagePassThru: {
				Type:        schema.TypeString,
				Description: "origin_error_page_pass_thru setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsBrotli: {
				Type:        schema.TypeString,
				Description: "brotli setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsPseudoIPv4: {
				Type:        schema.TypeString,
				Description: "pseudo_ipv4 setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsPseudoIPv4ValidatorID),
			},
			cisDomainSettingsPrefetchPreload: {
				Type:        schema.TypeString,
				Description: "prefetch_preload setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsResponseBuffering: {
				Type:        schema.TypeString,
				Description: "response_buffering setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsScriptLoadOptimisation: {
				Type:        schema.TypeString,
				Description: "script_load_optimization setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsServerSideExclude: {
				Type:        schema.TypeString,
				Description: "server_side_exclude setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsTLSClientAuth: {
				Type:        schema.TypeString,
				Description: "tls_client_auth setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsTrueClientIPHeader: {
				Type:        schema.TypeString,
				Description: "true_client_ip_header setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsWebSockets: {
				Type:        schema.TypeString,
				Description: "websockets setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsONOFFValidatorID),
			},
			cisDomainSettingsChallengeTTL: {
				Type:        schema.TypeInt,
				Description: "Challenge TTL setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsChallengeTTLValidatorID),
			},
			cisDomainSettingsMaxUpload: {
				Type:        schema.TypeInt,
				Description: "Maximum upload",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(
					ibmCISDomainSettings,
					cisDomainSettingsMaxUploadValidatorID),
			},
			cisDomainSettingsCipher: {
				Type:        schema.TypeSet,
				Description: "Cipher settings",
				Optional:    true,
				Computed:    true,
				Set:         schema.HashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: InvokeValidator(
						ibmCISDomainSettings,
						cisDomainSettingsCipherValidatorID),
				},
			},
			cisDomainSettingsMinify: {
				Type:        schema.TypeList,
				Description: "Minify setting",
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisDomainSettingsMinifyCSS: {
							Type:        schema.TypeString,
							Description: "Minify CSS setting",
							Required:    true,
							ValidateFunc: InvokeValidator(
								ibmCISDomainSettings,
								cisDomainSettingsONOFFValidatorID),
						},
						cisDomainSettingsMinifyHTML: {
							Type:        schema.TypeString,
							Description: "Minify HTML setting",
							Required:    true,
							ValidateFunc: InvokeValidator(
								ibmCISDomainSettings,
								cisDomainSettingsONOFFValidatorID),
						},
						cisDomainSettingsMinifyJS: {
							Type:        schema.TypeString,
							Description: "Minify JS setting",
							Required:    true,
							ValidateFunc: InvokeValidator(
								ibmCISDomainSettings,
								cisDomainSettingsONOFFValidatorID),
						},
					},
				},
			},
			cisDomainSettingsSecurityHeader: {
				Type:        schema.TypeList,
				Description: "Security Header Setting",
				Optional:    true,
				Computed:    true,
				MinItems:    1,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisDomainSettingsSecurityHeaderEnabled: {
							Type:        schema.TypeBool,
							Description: "security header enabled/disabled",
							Required:    true,
						},
						cisDomainSettingsSecurityHeaderIncludeSubdomains: {
							Type:        schema.TypeBool,
							Description: "security header subdomain included or not",
							Required:    true,
						},
						cisDomainSettingsSecurityHeaderMaxAge: {
							Type:        schema.TypeInt,
							Description: "security header max age",
							Required:    true,
						},
						cisDomainSettingsSecurityHeaderNoSniff: {
							Type:        schema.TypeBool,
							Description: "security header no sniff",
							Required:    true,
						},
					},
				},
			},
			cisDomainSettingsMobileRedirect: {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						cisDomainSettingsMobileRedirectStatus: {
							Type:        schema.TypeString,
							Description: "mobile redirect status",
							Required:    true,
							ValidateFunc: InvokeValidator(
								ibmCISDomainSettings,
								cisDomainSettingsONOFFValidatorID),
						},
						cisDomainSettingsMobileRedirectMobileSubdomain: {
							Type:        schema.TypeString,
							Description: "Mobile redirect subdomain",
							Optional:    true,
							Computed:    true,
						},
						cisDomainSettingsMobileRedirectStripURI: {
							Type:        schema.TypeBool,
							Description: "mobile redirect strip URI",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},

		Create:   resourceCISSettingsUpdate,
		Read:     resourceCISSettingsRead,
		Update:   resourceCISSettingsUpdate,
		Delete:   resourceCISSettingsDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func resourceIBMCISDomainSettingValidator() *ResourceValidator {

	sslSetting := "off, flexible, full, strict, origin_pull"
	tlsVersion := "1.1, 1.2, 1.3, 1.4"
	cnameFlatten := "flatten_at_root, flatten_all, flatten_none"
	imgSizeOptimize := "lossless, off, lossy"
	pseudoIPv4 := "overwrite_header, off, add_header"
	challengeTTL := "300, 900, 1800, 2700, 3600, 7200, 10800, 14400, 28800, 57600, 86400, 604800, 2592000, 31536000"
	maxUpload := "100, 125, 150, 175, 200, 225, 250, 275, 300, 325, 350, 375, 400, 425, 450, 475, 500"
	cipher := "ECDHE-ECDSA-AES128-GCM-SHA256,ECDHE-ECDSA-CHACHA20-POLY1305, ECDHE-RSA-AES128-GCM-SHA256,ECDHE-RSA-CHACHA20-POLY1305, ECDHE-ECDSA-AES128-SHA256, ECDHE-ECDSA-AES128-SHA, ECDHE-RSA-AES128-SHA256, ECDHE-RSA-AES128-SHA, AES128-GCM-SHA256, AES128-SHA256, AES128-SHA, ECDHE-ECDSA-AES256-GCM-SHA384, ECDHE-ECDSA-AES256-SHA384, ECDHE-RSA-AES256-GCM-SHA384, ECDHE-RSA-AES256-SHA384, ECDHE-RSA-AES256-SHA, AES256-GCM-SHA384, AES256-SHA256, AES256-SHA, DES-CBC3-SHA"

	validateSchema := make([]ValidateSchema, 1)

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsONOFFValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "on, off"})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsActiveDisableValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "active, disabled"})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsSSLSettingValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              sslSetting})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsTLSVersionValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              tlsVersion})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsCNAMEFlattenValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              cnameFlatten})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsImgSizeOptimizeValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              imgSizeOptimize})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsPseudoIPv4ValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              pseudoIPv4})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsChallengeTTLValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedIntValue,
			Type:                       TypeInt,
			Optional:                   true,
			AllowedValues:              challengeTTL})

	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsMaxUploadValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedIntValue,
			Type:                       TypeInt,
			Optional:                   true,
			AllowedValues:              maxUpload})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisDomainSettingsCipherValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              cipher})
	ibmCISDomainSettingResourceValidator := ResourceValidator{
		ResourceName: ibmCISDomainSettings,
		Schema:       validateSchema}
	return &ibmCISDomainSettingResourceValidator
}

var settingsList = []string{
	cisDomainSettingsDNSSEC,
	cisDomainSettingsWAF,
	cisDomainSettingsSSL,
	cisDomainSettingsMinTLSVersion,
	cisDomainSettingsCNAMEFlattening,
	cisDomainSettingsOpportunisticEncryption,
	cisDomainSettingsAutomaticHTPSRewrites,
	cisDomainSettingsAlwaysUseHTTPS,
	cisDomainSettingsIPv6,
	cisDomainSettingsBrowserCheck,
	cisDomainSettingsHotlinkProtection,
	cisDomainSettingsHTTP2,
	cisDomainSettingsImageLoadOptimization,
	cisDomainSettingsImageSizeOptimization,
	cisDomainSettingsIPGeoLocation,
	cisDomainSettingsOriginErrorPagePassThru,
	cisDomainSettingsBrotli,
	cisDomainSettingsPseudoIPv4,
	cisDomainSettingsPrefetchPreload,
	cisDomainSettingsResponseBuffering,
	cisDomainSettingsScriptLoadOptimisation,
	cisDomainSettingsServerSideExclude,
	cisDomainSettingsTLSClientAuth,
	cisDomainSettingsTrueClientIPHeader,
	cisDomainSettingsWebSockets,
	cisDomainSettingsChallengeTTL,
	cisDomainSettingsMinify,
	cisDomainSettingsSecurityHeader,
	cisDomainSettingsMobileRedirect,
	cisDomainSettingsMaxUpload,
	cisDomainSettingsCipher,
}

func resourceCISSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisDomainSettingsClientSession()
	if err != nil {
		return err
	}

	cisID := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(cisID)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	for _, item := range settingsList {
		var err error
		var resp *core.DetailedResponse

		switch item {
		case cisDomainSettingsDNSSEC:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateZoneDnssecOptions()
					opt.SetStatus(v.(string))
					_, resp, err = cisClient.UpdateZoneDnssec(opt)
				}
			}
		case cisDomainSettingsWAF:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateWebApplicationFirewallOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateWebApplicationFirewall(opt)
				}
			}
		case cisDomainSettingsSSL:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					cisClient, err := meta.(ClientSession).CisSSLClientSession()
					if err != nil {
						return err
					}
					cisClient.Crn = core.StringPtr(cisID)
					cisClient.ZoneIdentifier = core.StringPtr(zoneID)
					opt := cisClient.NewChangeSslSettingOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.ChangeSslSetting(opt)
				}
			}

		case cisDomainSettingsMinTLSVersion:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateMinTlsVersionOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateMinTlsVersion(opt)
				}
			}
		case cisDomainSettingsCNAMEFlattening:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateZoneCnameFlatteningOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateZoneCnameFlattening(opt)
				}
			}
		case cisDomainSettingsOpportunisticEncryption:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateOpportunisticEncryptionOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateOpportunisticEncryption(opt)
				}
			}
		case cisDomainSettingsAutomaticHTPSRewrites:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateAutomaticHttpsRewritesOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateAutomaticHttpsRewrites(opt)
				}
			}
		case cisDomainSettingsAlwaysUseHTTPS:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateAlwaysUseHttpsOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateAlwaysUseHttps(opt)
				}
			}
		case cisDomainSettingsIPv6:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateIpv6Options()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateIpv6(opt)
				}
			}
		case cisDomainSettingsBrowserCheck:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateBrowserCheckOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateBrowserCheck(opt)
				}
			}
		case cisDomainSettingsHotlinkProtection:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateHotlinkProtectionOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateHotlinkProtection(opt)
				}
			}
		case cisDomainSettingsHTTP2:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateHttp2Options()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateHttp2(opt)
				}
			}
		case cisDomainSettingsImageLoadOptimization:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateImageLoadOptimizationOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateImageLoadOptimization(opt)
				}
			}
		case cisDomainSettingsImageSizeOptimization:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateImageSizeOptimizationOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateImageSizeOptimization(opt)
				}
			}
		case cisDomainSettingsIPGeoLocation:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateIpGeolocationOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateIpGeolocation(opt)
				}
			}
		case cisDomainSettingsOriginErrorPagePassThru:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateEnableErrorPagesOnOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateEnableErrorPagesOn(opt)
				}
			}
		case cisDomainSettingsPseudoIPv4:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdatePseudoIpv4Options()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdatePseudoIpv4(opt)
				}
			}
		case cisDomainSettingsPrefetchPreload:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdatePrefetchPreloadOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdatePrefetchPreload(opt)
				}
			}
		case cisDomainSettingsResponseBuffering:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateResponseBufferingOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateResponseBuffering(opt)
				}
			}
		case cisDomainSettingsScriptLoadOptimisation:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateScriptLoadOptimizationOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateScriptLoadOptimization(opt)
				}
			}
		case cisDomainSettingsServerSideExclude:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateServerSideExcludeOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateServerSideExclude(opt)
				}
			}
		case cisDomainSettingsTLSClientAuth:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateTlsClientAuthOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateTlsClientAuth(opt)
				}
			}
		case cisDomainSettingsTrueClientIPHeader:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateTrueClientIpOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateTrueClientIp(opt)
				}
			}
		case cisDomainSettingsWebSockets:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateWebSocketsOptions()
					opt.SetValue(v.(string))
					_, resp, err = cisClient.UpdateWebSockets(opt)
				}
			}
		case cisDomainSettingsChallengeTTL:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateChallengeTtlOptions()
					opt.SetValue(int64(v.(int)))
					_, resp, err = cisClient.UpdateChallengeTTL(opt)
				}
			}
		case cisDomainSettingsMaxUpload:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					opt := cisClient.NewUpdateMaxUploadOptions()
					opt.SetValue(int64(v.(int)))
					_, resp, err = cisClient.UpdateMaxUpload(opt)
				}
			}
		case cisDomainSettingsCipher:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					cipherValue := expandStringList(v.(*schema.Set).List())
					opt := cisClient.NewUpdateCiphersOptions()
					opt.SetValue(cipherValue)
					_, resp, err = cisClient.UpdateCiphers(opt)
				}
			}
		case cisDomainSettingsMinify:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					dataMap := v.([]interface{})[0].(map[string]interface{})
					css := dataMap[cisDomainSettingsMinifyCSS].(string)
					html := dataMap[cisDomainSettingsMinifyHTML].(string)
					js := dataMap[cisDomainSettingsMinifyJS].(string)
					minifyVal, err := cisClient.NewMinifySettingValue(css, html, js)
					if err != nil {
						log.Println("Invalid minfiy setting values")
						return err
					}
					opt := cisClient.NewUpdateMinifyOptions()
					opt.SetValue(minifyVal)
					_, resp, err = cisClient.UpdateMinify(opt)
				}
			}
		case cisDomainSettingsSecurityHeader:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					dataMap := v.([]interface{})[0].(map[string]interface{})
					enabled := dataMap[cisDomainSettingsSecurityHeaderEnabled].(bool)
					nosniff := dataMap[cisDomainSettingsSecurityHeaderNoSniff].(bool)
					includeSubdomain := dataMap[cisDomainSettingsSecurityHeaderIncludeSubdomains].(bool)
					maxAge := int64(dataMap[cisDomainSettingsSecurityHeaderMaxAge].(int))
					securityVal, err := cisClient.NewSecurityHeaderSettingValueStrictTransportSecurity(
						enabled, maxAge, includeSubdomain, nosniff)
					if err != nil {
						log.Println("Invalid security header setting values")
						return err
					}
					securityOpt, err := cisClient.NewSecurityHeaderSettingValue(securityVal)
					if err != nil {
						log.Println("Invalid security header setting options")
						return err
					}
					opt := cisClient.NewUpdateSecurityHeaderOptions()
					opt.SetValue(securityOpt)
					_, resp, err = cisClient.UpdateSecurityHeader(opt)
				}
			}
		case cisDomainSettingsMobileRedirect:
			if d.HasChange(item) {
				if v, ok := d.GetOk(item); ok {
					dataMap := v.([]interface{})[0].(map[string]interface{})
					status := dataMap[cisDomainSettingsMobileRedirectStatus].(string)
					mobileSubdomain := dataMap[cisDomainSettingsMobileRedirectMobileSubdomain].(string)
					stripURI := dataMap[cisDomainSettingsMobileRedirectStripURI].(bool)
					mobileOpt, err := cisClient.NewMobileRedirecSettingValue(status, mobileSubdomain, stripURI)
					if err != nil {
						log.Println("Invalid mobile redirect options")
						return err
					}
					opt := cisClient.NewUpdateMobileRedirectOptions()
					opt.SetValue(mobileOpt)
					_, resp, err = cisClient.UpdateMobileRedirect(opt)
				}
			}
		}
		if err != nil {
			if resp != nil && resp.StatusCode == 405 {
				log.Printf("[WARN] Update %s : %s", item, err)
				continue
			}
			log.Printf("Update settings Failed on %s, %v\n", item, resp)
			return err
		}
	}
	d.SetId(convertCisToTfTwoVar(zoneID, cisID))
	return resourceCISSettingsRead(d, meta)
}

func resourceCISSettingsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisDomainSettingsClientSession()
	if err != nil {
		return err
	}

	zoneID, crn, _ := convertTftoCisTwoVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneIdentifier = core.StringPtr(zoneID)

	for _, item := range settingsList {
		var settingErr error
		var settingResponse *core.DetailedResponse
		switch item {
		case cisDomainSettingsDNSSEC:
			opt := cisClient.NewGetZoneDnssecOptions()
			result, resp, err := cisClient.GetZoneDnssec(opt)
			if err == nil {
				d.Set(cisDomainSettingsDNSSEC, result.Result.Status)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsWAF:
			opt := cisClient.NewGetWebApplicationFirewallOptions()
			result, resp, err := cisClient.GetWebApplicationFirewall(opt)
			if err == nil {
				d.Set(cisDomainSettingsWAF, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsSSL:
			cisClient, err := meta.(ClientSession).CisSSLClientSession()
			if err != nil {
				return err
			}
			cisClient.Crn = core.StringPtr(crn)
			cisClient.ZoneIdentifier = core.StringPtr(zoneID)
			opt := cisClient.NewGetSslSettingOptions()
			result, resp, err := cisClient.GetSslSetting(opt)
			if err == nil {
				d.Set(cisDomainSettingsSSL, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsBrotli:
			cisClient, err := meta.(ClientSession).CisAPI()
			if err != nil {
				return err
			}
			settingsResult, err := cisClient.Settings().GetSetting(crn, zoneID, item)
			if err == nil {
				settingsObj := *settingsResult
				d.Set(item, settingsObj.Value)
			}
			settingErr = err

		case cisDomainSettingsMinTLSVersion:
			opt := cisClient.NewGetMinTlsVersionOptions()
			result, resp, err := cisClient.GetMinTlsVersion(opt)
			if err == nil {
				d.Set(cisDomainSettingsMinTLSVersion, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsCNAMEFlattening:
			opt := cisClient.NewGetZoneCnameFlatteningOptions()
			result, resp, err := cisClient.GetZoneCnameFlattening(opt)
			if err == nil {
				d.Set(cisDomainSettingsCNAMEFlattening, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsOpportunisticEncryption:
			opt := cisClient.NewGetOpportunisticEncryptionOptions()
			result, resp, err := cisClient.GetOpportunisticEncryption(opt)
			if err == nil {
				d.Set(cisDomainSettingsOpportunisticEncryption, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsAutomaticHTPSRewrites:
			opt := cisClient.NewGetAutomaticHttpsRewritesOptions()
			result, resp, err := cisClient.GetAutomaticHttpsRewrites(opt)
			if err == nil {
				d.Set(cisDomainSettingsAutomaticHTPSRewrites, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsAlwaysUseHTTPS:
			opt := cisClient.NewGetAlwaysUseHttpsOptions()
			result, resp, err := cisClient.GetAlwaysUseHttps(opt)
			if err == nil {
				d.Set(cisDomainSettingsAlwaysUseHTTPS, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsIPv6:
			opt := cisClient.NewGetIpv6Options()
			result, resp, err := cisClient.GetIpv6(opt)
			if err == nil {
				d.Set(cisDomainSettingsIPv6, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsBrowserCheck:
			opt := cisClient.NewGetBrowserCheckOptions()
			result, resp, err := cisClient.GetBrowserCheck(opt)
			if err == nil {
				d.Set(cisDomainSettingsBrowserCheck, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsHotlinkProtection:
			opt := cisClient.NewGetHotlinkProtectionOptions()
			result, resp, err := cisClient.GetHotlinkProtection(opt)
			if err == nil {
				d.Set(cisDomainSettingsHotlinkProtection, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsHTTP2:
			opt := cisClient.NewGetHttp2Options()
			result, resp, err := cisClient.GetHttp2(opt)
			if err == nil {
				d.Set(cisDomainSettingsHTTP2, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsImageLoadOptimization:
			opt := cisClient.NewGetImageLoadOptimizationOptions()
			result, resp, err := cisClient.GetImageLoadOptimization(opt)
			if err == nil {
				d.Set(cisDomainSettingsImageLoadOptimization, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsImageSizeOptimization:
			opt := cisClient.NewGetImageSizeOptimizationOptions()
			result, resp, err := cisClient.GetImageSizeOptimization(opt)
			if err == nil {
				d.Set(cisDomainSettingsImageSizeOptimization, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsIPGeoLocation:
			opt := cisClient.NewGetIpGeolocationOptions()
			result, resp, err := cisClient.GetIpGeolocation(opt)
			if err == nil {
				d.Set(cisDomainSettingsIPGeoLocation, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsOriginErrorPagePassThru:
			opt := cisClient.NewGetEnableErrorPagesOnOptions()
			result, resp, err := cisClient.GetEnableErrorPagesOn(opt)
			if err == nil {
				d.Set(cisDomainSettingsOriginErrorPagePassThru, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsPseudoIPv4:
			opt := cisClient.NewGetPseudoIpv4Options()
			result, resp, err := cisClient.GetPseudoIpv4(opt)
			if err == nil {
				d.Set(cisDomainSettingsPseudoIPv4, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsPrefetchPreload:
			opt := cisClient.NewGetPrefetchPreloadOptions()
			result, resp, err := cisClient.GetPrefetchPreload(opt)
			if err == nil {
				d.Set(cisDomainSettingsPrefetchPreload, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsResponseBuffering:
			opt := cisClient.NewGetResponseBufferingOptions()
			result, resp, err := cisClient.GetResponseBuffering(opt)
			if err == nil {
				d.Set(cisDomainSettingsResponseBuffering, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsScriptLoadOptimisation:
			opt := cisClient.NewGetScriptLoadOptimizationOptions()
			result, resp, err := cisClient.GetScriptLoadOptimization(opt)
			if err == nil {
				d.Set(cisDomainSettingsScriptLoadOptimisation, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsServerSideExclude:
			opt := cisClient.NewGetServerSideExcludeOptions()
			result, resp, err := cisClient.GetServerSideExclude(opt)
			if err == nil {
				d.Set(cisDomainSettingsServerSideExclude, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsTLSClientAuth:
			opt := cisClient.NewGetTlsClientAuthOptions()
			result, resp, err := cisClient.GetTlsClientAuth(opt)
			if err == nil {
				d.Set(cisDomainSettingsTLSClientAuth, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsTrueClientIPHeader:
			opt := cisClient.NewGetTrueClientIpOptions()
			result, resp, err := cisClient.GetTrueClientIp(opt)
			if err == nil {
				d.Set(cisDomainSettingsTrueClientIPHeader, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsWebSockets:
			opt := cisClient.NewGetWebSocketsOptions()
			result, resp, err := cisClient.GetWebSockets(opt)
			if err == nil {
				d.Set(cisDomainSettingsWebSockets, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsChallengeTTL:
			opt := cisClient.NewGetChallengeTtlOptions()
			result, resp, err := cisClient.GetChallengeTTL(opt)
			if err == nil {
				d.Set(cisDomainSettingsChallengeTTL, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsMaxUpload:
			opt := cisClient.NewGetMaxUploadOptions()
			result, resp, err := cisClient.GetMaxUpload(opt)
			if err == nil {
				d.Set(cisDomainSettingsMaxUpload, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsCipher:
			opt := cisClient.NewGetCiphersOptions()
			result, resp, err := cisClient.GetCiphers(opt)
			if err == nil {
				d.Set(cisDomainSettingsCipher, result.Result.Value)
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsMinify:
			opt := cisClient.NewGetMinifyOptions()
			result, resp, err := cisClient.GetMinify(opt)
			if err == nil {
				minify := result.Result.Value
				value := map[string]string{
					cisDomainSettingsMinifyCSS:  *minify.Css,
					cisDomainSettingsMinifyHTML: *minify.HTML,
					cisDomainSettingsMinifyJS:   *minify.Js,
				}
				d.Set(cisDomainSettingsMinify, []interface{}{value})
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsSecurityHeader:
			opt := cisClient.NewGetSecurityHeaderOptions()
			result, resp, err := cisClient.GetSecurityHeader(opt)
			if err == nil {

				if result.Result.Value != nil && result.Result.Value.StrictTransportSecurity != nil {

					securityHeader := result.Result.Value.StrictTransportSecurity
					value := map[string]interface{}{}
					if securityHeader.Enabled != nil {
						value[cisDomainSettingsSecurityHeaderEnabled] = *securityHeader.Enabled
					}
					if securityHeader.Nosniff != nil {
						value[cisDomainSettingsSecurityHeaderNoSniff] = *securityHeader.Nosniff
					}
					if securityHeader.IncludeSubdomains != nil {
						value[cisDomainSettingsSecurityHeaderIncludeSubdomains] = *securityHeader.IncludeSubdomains
					}
					if securityHeader.MaxAge != nil {
						value[cisDomainSettingsSecurityHeaderMaxAge] = *securityHeader.MaxAge
					}
					d.Set(cisDomainSettingsSecurityHeader, []interface{}{value})
				}
			}
			settingResponse = resp
			settingErr = err

		case cisDomainSettingsMobileRedirect:
			opt := cisClient.NewGetMobileRedirectOptions()
			result, resp, err := cisClient.GetMobileRedirect(opt)
			if err == nil {
				if result.Result.Value != nil {

					value := result.Result.Value

					uri := map[string]interface{}{}
					if value.MobileSubdomain != nil {
						uri[cisDomainSettingsMobileRedirectMobileSubdomain] = *value.MobileSubdomain
					}
					if value.Status != nil {
						uri[cisDomainSettingsMobileRedirectStatus] = *value.Status
					}
					if value.StripURI != nil {
						uri[cisDomainSettingsMobileRedirectStripURI] = *value.StripURI
					}
					d.Set(cisDomainSettingsMobileRedirect, []interface{}{uri})
				}
			}
			settingResponse = resp
			settingErr = err
		}

		if settingErr != nil {
			if settingResponse != nil && settingResponse.StatusCode == 405 {
				log.Printf("[WARN] Get %s. : %s", item, settingErr)
				continue
			}
			log.Printf("Get settings failed on %s, %v\n", item, settingErr)
			return settingErr
		}
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	return nil
}

func resourceCISSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS resource
	d.SetId("")
	return nil
}
