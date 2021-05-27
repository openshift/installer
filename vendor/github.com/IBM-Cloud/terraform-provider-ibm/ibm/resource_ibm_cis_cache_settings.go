// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	ibmCISCacheSettings               = "ibm_cis_cache_settings"
	cisCacheSettingsCachingLevel      = "caching_level"
	cisCacheSettingsBrowserExpiration = "browser_expiration"
	cisCacheSettingsDevelopmentMode   = "development_mode"
	cisCacheSettingsQueryStringSort   = "query_string_sort"
	cisCachePurgeAll                  = "purge_all"
	cisCachePurgeByURLs               = "purge_by_urls"
	cisCachePurgeByCacheTags          = "purge_by_tags"
	cisCachePurgeByHosts              = "purge_by_hosts"
	cisCacheSettingsOnOffValidatorID  = "on_off_validator_id"
	cisCacheServeStaleContent         = "serve_stale_content"
)

func resourceIBMCISCacheSettings() *schema.Resource {
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
			cisCacheSettingsCachingLevel: {
				Type:        schema.TypeString,
				Description: "Cache level setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(ibmCISCacheSettings,
					cisCacheSettingsCachingLevel),
			},
			cisCacheServeStaleContent: {
				Type:        schema.TypeString,
				Description: "Serve Stale Content ",
				Default:     "on",
				Optional:    true,
				ValidateFunc: InvokeValidator(ibmCISCacheSettings,
					cisCacheServeStaleContent),
			},
			cisCacheSettingsBrowserExpiration: {
				Type:        schema.TypeInt,
				Description: "Browser Expiration setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(ibmCISCacheSettings,
					cisCacheSettingsBrowserExpiration),
			},
			cisCacheSettingsDevelopmentMode: {
				Type:        schema.TypeString,
				Description: "Development mode setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(ibmCISCacheSettings,
					cisCacheSettingsOnOffValidatorID),
			},
			cisCacheSettingsQueryStringSort: {
				Type:        schema.TypeString,
				Description: "Query String sort setting",
				Optional:    true,
				Computed:    true,
				ValidateFunc: InvokeValidator(ibmCISCacheSettings,
					cisCacheSettingsOnOffValidatorID),
			},
			cisCachePurgeAll: {
				Type:        schema.TypeBool,
				Description: "Purge all setting",
				Optional:    true,
				ConflictsWith: []string{
					cisCachePurgeByURLs,
					cisCachePurgeByCacheTags,
					cisCachePurgeByHosts},
			},
			cisCachePurgeByURLs: {
				Type:        schema.TypeList,
				Description: "Purge by URLs",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{
					cisCachePurgeAll,
					cisCachePurgeByCacheTags,
					cisCachePurgeByHosts},
			},
			cisCachePurgeByCacheTags: {
				Type:        schema.TypeList,
				Description: "Purge by tags",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{
					cisCachePurgeAll,
					cisCachePurgeByURLs,
					cisCachePurgeByHosts},
			},
			cisCachePurgeByHosts: {
				Type:        schema.TypeList,
				Description: "Purge by hosts",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{
					cisCachePurgeAll,
					cisCachePurgeByURLs,
					cisCachePurgeByCacheTags,
				},
			},
		},
		Create:   resourceCISCacheSettingsUpdate,
		Read:     resourceCISCacheSettingsRead,
		Update:   resourceCISCacheSettingsUpdate,
		Delete:   resourceCISCacheSettingsDelete,
		Importer: &schema.ResourceImporter{},
	}
}

func resourceIBMCISCacheSettingsValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	browserCacheTTL := "0, 30, 60, 300, 1200, 1800, 3600, 7200, 10800, 14400," +
		"18000, 28800, 43200, 57600, 72000, 86400, 172800, 259200, 345600, 432000," +
		"691200, 1382400, 2073600, 2678400, 5356800, 16070400, 31536000"
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisCacheSettingsOnOffValidatorID,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "on, off"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisCacheSettingsBrowserExpiration,
			ValidateFunctionIdentifier: ValidateAllowedIntValue,
			Type:                       TypeInt,
			Optional:                   true,
			AllowedValues:              browserCacheTTL})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisCacheServeStaleContent,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Optional:                   true,
			AllowedValues:              "on, off"})
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 cisCacheSettingsCachingLevel,
			ValidateFunctionIdentifier: ValidateAllowedStringValue,
			Type:                       TypeString,
			Required:                   true,
			AllowedValues:              "basic, simplified, aggressive"})
	ibmCISCacheSettingsResourceValidator := ResourceValidator{
		ResourceName: ibmCISCacheSettings,
		Schema:       validateSchema}
	return &ibmCISCacheSettingsResourceValidator
}

func resourceCISCacheSettingsUpdate(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisCacheClientSession()
	if err != nil {
		return err
	}
	crn := d.Get(cisID).(string)
	zoneID, _, err := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	if d.HasChange(cisCacheSettingsCachingLevel) ||
		d.HasChange(cisCacheSettingsBrowserExpiration) ||
		d.HasChange(cisCacheSettingsDevelopmentMode) ||
		d.HasChange(cisCacheSettingsQueryStringSort) ||
		d.HasChange(cisCachePurgeAll) ||
		d.HasChange(cisCachePurgeByURLs) ||
		d.HasChange(cisCachePurgeByCacheTags) ||
		d.HasChange(cisCachePurgeByHosts) ||
		d.HasChange(cisCacheServeStaleContent) {

		// Caching Level Setting
		if value, ok := d.GetOk(cisCacheSettingsCachingLevel); ok {
			opt := cisClient.NewUpdateCacheLevelOptions()
			opt.SetValue(value.(string))
			_, resp, err := cisClient.UpdateCacheLevel(opt)
			if err != nil {
				log.Printf("Update caching level failed : %v\n", resp)
				return err
			}
		}
		// Serve Stale Content Setting
		if value, ok := d.GetOk(cisCacheServeStaleContent); ok {
			opt := cisClient.NewUpdateServeStaleContentOptions()
			opt.SetValue(value.(string))
			_, resp, err := cisClient.UpdateServeStaleContent(opt)
			if err != nil {
				log.Printf("Update Serve Stale Content Setting failed : %v\n", resp)
				return err
			}
		}

		// Browser Expiration setting
		if value, ok := d.GetOk(cisCacheSettingsBrowserExpiration); ok {
			opt := cisClient.NewUpdateBrowserCacheTtlOptions()
			opt.SetValue(int64(value.(int)))
			_, resp, err := cisClient.UpdateBrowserCacheTTL(opt)
			if err != nil {
				log.Printf("Update browser expiration setting failed : %v\n", resp)
				return err
			}
		}

		// development mode setting
		if value, ok := d.GetOk(cisCacheSettingsDevelopmentMode); ok {
			opt := cisClient.NewUpdateDevelopmentModeOptions()
			opt.SetValue(value.(string))
			_, resp, err := cisClient.UpdateDevelopmentMode(opt)
			if err != nil {
				log.Printf("Update development mode setting failed : %v\n", resp)
				return err
			}
		}
		// Query string sort setting
		if value, ok := d.GetOk(cisCacheSettingsQueryStringSort); ok {
			opt := cisClient.NewUpdateQueryStringSortOptions()
			opt.SetValue(value.(string))
			_, resp, err := cisClient.UpdateQueryStringSort(opt)
			if err != nil {
				log.Printf("Update query string sort setting failed : %v\n", resp)
				return err
			}
		}

		if value, ok := d.GetOkExists(cisCachePurgeAll); ok {
			if value.(bool) == true {
				opt := cisClient.NewPurgeAllOptions()
				result, response, err := cisClient.PurgeAll(opt)
				if err != nil {
					log.Printf("Purge all failed : %v", response)
					return err
				}
				log.Printf("Purge all successful : %s", *result.Result.ID)
			}
		}
		if value, ok := d.GetOk(cisCachePurgeByURLs); ok {
			urls := expandStringList(value.([]interface{}))
			opt := cisClient.NewPurgeByUrlsOptions()
			opt.SetFiles(urls)
			_, response, err := cisClient.PurgeByUrls(opt)
			if err != nil {
				log.Printf("Purge by urls failed : %v", response)
				return err
			}
		}
		if value, ok := d.GetOk(cisCachePurgeByCacheTags); ok {
			cacheTags := expandStringList(value.([]interface{}))
			opt := cisClient.NewPurgeByCacheTagsOptions()
			opt.SetTags(cacheTags)
			result, response, err := cisClient.PurgeByCacheTags(opt)
			if err != nil {
				log.Printf("Purge by cache tags failed : %v", response)
				return err
			}
			log.Printf("Purge by tags successful : %s", *result.Result.ID)

		}
		if value, ok := d.GetOk(cisCachePurgeByHosts); ok {
			hosts := expandStringList(value.([]interface{}))
			opt := cisClient.NewPurgeByHostsOptions()
			opt.SetHosts(hosts)
			result, response, err := cisClient.PurgeByHosts(opt)
			if err != nil {
				log.Printf("Purge by hosts failed : %v", response)
				return err
			}
			log.Printf("Purge by hosts successful : %s", *result.Result.ID)
		}
	}
	d.SetId(convertCisToTfTwoVar(zoneID, crn))
	return resourceCISCacheSettingsRead(d, meta)
}

func resourceCISCacheSettingsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisCacheClientSession()
	if err != nil {
		return err
	}
	zoneID, crn, _ := convertTftoCisTwoVar(d.Id())
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	// Caching Level Setting
	cacheLevel, resp, err := cisClient.GetCacheLevel(cisClient.NewGetCacheLevelOptions())
	if err != nil {
		log.Printf("Get caching leve setting failed : %v\n", resp)
		return err
	}

	// Serve Stale Content setting
	servestaleContent, resp, err := cisClient.GetServeStaleContent(cisClient.NewGetServeStaleContentOptions())
	if err != nil {
		log.Printf("Get Serve Stale Content setting failed : %v\n", resp)
		return err
	}

	// Browser Expiration setting
	browserCacheTTL, resp, err := cisClient.GetBrowserCacheTTL(
		cisClient.NewGetBrowserCacheTtlOptions())
	if err != nil {
		log.Printf("Get browser expiration setting failed : %v\n", resp)
		return err
	}

	// development mode setting
	devMode, resp, err := cisClient.GetDevelopmentMode(
		cisClient.NewGetDevelopmentModeOptions())
	if err != nil {
		log.Printf("Get development mode setting failed : %v", resp)
		return err
	}

	// Query string sort setting
	queryStringSort, resp, err := cisClient.GetQueryStringSort(
		cisClient.NewGetQueryStringSortOptions())
	if err != nil {
		log.Printf("Get query string sort setting failed : %v", resp)
		return err
	}
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	d.Set(cisCacheSettingsBrowserExpiration, *browserCacheTTL.Result.Value)
	d.Set(cisCacheSettingsCachingLevel, *cacheLevel.Result.Value)
	d.Set(cisCacheSettingsDevelopmentMode, *devMode.Result.Value)
	d.Set(cisCacheSettingsQueryStringSort, *queryStringSort.Result.Value)
	d.Set(cisCacheServeStaleContent, *servestaleContent.Result.Value)
	return nil
}

func resourceCISCacheSettingsDelete(d *schema.ResourceData, meta interface{}) error {
	// Nothing to delete on CIS resource
	d.SetId("")
	return nil
}
