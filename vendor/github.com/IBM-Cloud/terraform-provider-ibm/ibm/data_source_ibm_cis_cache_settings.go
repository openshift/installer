// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"log"
	"time"

	"github.com/IBM/go-sdk-core/v4/core"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceIBMCISCacheSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCISCacheSettingsRead,
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
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Cache Level Setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cache level id",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cache level value",
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "cache level editable",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cache level modified on",
						},
					},
				},
			},
			cisCacheServeStaleContent: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Serve Stale Content ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "serve stale content id ",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "serve stale content value ",
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "serve stale content editable ",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "serve stale content modified on ",
						},
					},
				},
			},
			cisCacheSettingsBrowserExpiration: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Browser Expiration setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "browser expiration id",
						},
						"value": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "browser expiration value",
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "browser expiration editable",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "browser expiration modified on",
						},
					},
				},
			},
			cisCacheSettingsDevelopmentMode: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Development mode setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "development mode id",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "development mode value",
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "development mode editable",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "development mode modified on",
						},
					},
				},
			},
			cisCacheSettingsQueryStringSort: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Query String sort setting",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "query string sort id",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "query qtring sort value",
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "query string sort editable",
						},
						"modified_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "query string sort modified on",
						},
					},
				},
			},
		},
	}
}
func dataSourceCISCacheSettingsRead(d *schema.ResourceData, meta interface{}) error {
	cisClient, err := meta.(ClientSession).CisCacheClientSession()
	if err != nil {
		return err
	}

	crn := d.Get(cisID).(string)
	zoneID, _, _ := convertTftoCisTwoVar(d.Get(cisDomainID).(string))
	cisClient.Crn = core.StringPtr(crn)
	cisClient.ZoneID = core.StringPtr(zoneID)

	// Cache Level Setting
	cacheLevel_result, resp, err := cisClient.GetCacheLevel(cisClient.NewGetCacheLevelOptions())

	if err != nil {
		log.Printf("Get Cache Level  setting failed : %v\n", resp)
		return err
	}
	if cacheLevel_result != nil || cacheLevel_result.Result != nil {

		cacheLevels := make([]map[string]interface{}, 0)
		cacheLevel := make(map[string]interface{})

		if cacheLevel_result.Result.ID != nil {
			cacheLevel["id"] = cacheLevel_result.Result.ID
		}
		if cacheLevel_result.Result.Value != nil {
			cacheLevel["value"] = cacheLevel_result.Result.Value
		}
		if cacheLevel_result.Result.Editable != nil {
			cacheLevel["editable"] = cacheLevel_result.Result.Editable
		}
		if cacheLevel_result.Result.ModifiedOn != nil {
			cacheLevel["modified_on"] = cacheLevel_result.Result.ModifiedOn
		}
		cacheLevels = append(cacheLevels, cacheLevel)
		d.Set(cisCacheSettingsCachingLevel, cacheLevels)

	}
	// Serve Stale Content setting
	servestaleContent_result, resp, err := cisClient.GetServeStaleContent(cisClient.NewGetServeStaleContentOptions())

	if err != nil {
		log.Printf("Get Serve Stale Content setting failed : %v\n", resp)
		return err
	}
	if servestaleContent_result != nil || servestaleContent_result.Result != nil {

		servestalecontents := make([]map[string]interface{}, 0)
		servestalecontent := make(map[string]interface{})

		if servestaleContent_result.Result.ID != nil {
			servestalecontent["id"] = servestaleContent_result.Result.ID
		}
		if servestaleContent_result.Result.Value != nil {
			servestalecontent["value"] = servestaleContent_result.Result.Value
		}
		if servestaleContent_result.Result.Editable != nil {
			servestalecontent["editable"] = servestaleContent_result.Result.Editable
		}
		if servestaleContent_result.Result.ModifiedOn != nil {
			servestalecontent["modified_on"] = servestaleContent_result.Result.ModifiedOn
		}
		servestalecontents = append(servestalecontents, servestalecontent)
		d.Set(cisCacheServeStaleContent, servestalecontents)

	}

	// Browser Expiration setting
	browserCacheTTL_result, resp, err := cisClient.GetBrowserCacheTTL(cisClient.NewGetBrowserCacheTtlOptions())

	if err != nil {
		log.Printf("Get browser expiration setting failed : %v\n", resp)
		return err
	}
	if browserCacheTTL_result != nil || browserCacheTTL_result.Result != nil {

		browserCacheTTLs := make([]map[string]interface{}, 0)
		browserCacheTTL := make(map[string]interface{})

		if browserCacheTTL_result.Result.ID != nil {
			browserCacheTTL["id"] = browserCacheTTL_result.Result.ID
		}
		if browserCacheTTL_result.Result.Value != nil {
			browserCacheTTL["value"] = browserCacheTTL_result.Result.Value
		}
		if browserCacheTTL_result.Result.Editable != nil {
			browserCacheTTL["editable"] = browserCacheTTL_result.Result.Editable
		}
		if browserCacheTTL_result.Result.ModifiedOn != nil {
			browserCacheTTL["modified_on"] = browserCacheTTL_result.Result.ModifiedOn
		}
		browserCacheTTLs = append(browserCacheTTLs, browserCacheTTL)
		d.Set(cisCacheSettingsBrowserExpiration, browserCacheTTLs)

	}
	// development mode setting
	devMode_result, resp, err := cisClient.GetDevelopmentMode(cisClient.NewGetDevelopmentModeOptions())

	if err != nil {
		log.Printf("Get development mode setting failed : %v", resp)
		return err
	}
	if devMode_result != nil || devMode_result.Result != nil {

		devModes := make([]map[string]interface{}, 0)
		devMode := make(map[string]interface{})

		if devMode_result.Result.ID != nil {
			devMode["id"] = devMode_result.Result.ID
		}
		if devMode_result.Result.Value != nil {
			devMode["value"] = devMode_result.Result.Value
		}
		if devMode_result.Result.Editable != nil {
			devMode["editable"] = devMode_result.Result.Editable
		}
		if devMode_result.Result.ModifiedOn != nil {
			devMode["modified_on"] = devMode_result.Result.ModifiedOn
		}
		devModes = append(devModes, devMode)
		d.Set(cisCacheSettingsDevelopmentMode, devModes)

	}

	// Query string sort setting
	queryStringSort_result, resp, err := cisClient.GetQueryStringSort(cisClient.NewGetQueryStringSortOptions())

	if err != nil {
		log.Printf("Get query string sort setting failed : %v", resp)
		return err
	}
	if queryStringSort_result != nil || queryStringSort_result.Result != nil {

		queryStringSorts := make([]map[string]interface{}, 0)
		queryStringSort := make(map[string]interface{})

		if queryStringSort_result.Result.ID != nil {
			queryStringSort["id"] = queryStringSort_result.Result.ID
		}
		if queryStringSort_result.Result.Value != nil {
			queryStringSort["value"] = queryStringSort_result.Result.Value
		}
		if queryStringSort_result.Result.Editable != nil {
			queryStringSort["editable"] = queryStringSort_result.Result.Editable
		}
		if queryStringSort_result.Result.ModifiedOn != nil {
			queryStringSort["modified_on"] = queryStringSort_result.Result.ModifiedOn
		}
		queryStringSorts = append(queryStringSorts, queryStringSort)
		d.Set(cisCacheSettingsQueryStringSort, queryStringSorts)

	}
	d.SetId(dataSourceIBMCISCacheSettingID(d))
	d.Set(cisID, crn)
	d.Set(cisDomainID, zoneID)
	return nil
}
func dataSourceIBMCISCacheSettingID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
