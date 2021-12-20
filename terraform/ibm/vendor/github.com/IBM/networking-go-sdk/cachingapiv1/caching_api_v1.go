/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * IBM OpenAPI SDK Code Generator Version: 3.26.0-4b317b0c-20210127-171701
 */
 

// Package cachingapiv1 : Operations and models for the CachingApiV1 service
package cachingapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
	"net/http"
	"reflect"
	"time"
)

// CachingApiV1 : This document describes CIS caching  API.
//
// Version: 1.0.0
type CachingApiV1 struct {
	Service *core.BaseService

	// cloud resource name.
	Crn *string

	// zone id.
	ZoneID *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "caching_api"

// CachingApiV1Options : Service options
type CachingApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// cloud resource name.
	Crn *string `validate:"required"`

	// zone id.
	ZoneID *string `validate:"required"`
}

// NewCachingApiV1UsingExternalConfig : constructs an instance of CachingApiV1 with passed in options and external configuration.
func NewCachingApiV1UsingExternalConfig(options *CachingApiV1Options) (cachingApi *CachingApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	cachingApi, err = NewCachingApiV1(options)
	if err != nil {
		return
	}

	err = cachingApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = cachingApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewCachingApiV1 : constructs an instance of CachingApiV1 with passed in options.
func NewCachingApiV1(options *CachingApiV1Options) (service *CachingApiV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	err = core.ValidateStruct(options, "options")
	if err != nil {
		return
	}

	baseService, err := core.NewBaseService(serviceOptions)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
			return
		}
	}

	service = &CachingApiV1{
		Service: baseService,
		Crn: options.Crn,
		ZoneID: options.ZoneID,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "cachingApi" suitable for processing requests.
func (cachingApi *CachingApiV1) Clone() *CachingApiV1 {
	if core.IsNil(cachingApi) {
		return nil
	}
	clone := *cachingApi
	clone.Service = cachingApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (cachingApi *CachingApiV1) SetServiceURL(url string) error {
	return cachingApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (cachingApi *CachingApiV1) GetServiceURL() string {
	return cachingApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (cachingApi *CachingApiV1) SetDefaultHeaders(headers http.Header) {
	cachingApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (cachingApi *CachingApiV1) SetEnableGzipCompression(enableGzip bool) {
	cachingApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (cachingApi *CachingApiV1) GetEnableGzipCompression() bool {
	return cachingApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (cachingApi *CachingApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	cachingApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (cachingApi *CachingApiV1) DisableRetries() {
	cachingApi.Service.DisableRetries()
}

// PurgeAll : Purge all
// All resources in CDN edge servers' cache should be removed. This may have dramatic affects on your origin server load
// after performing this action.
func (cachingApi *CachingApiV1) PurgeAll(purgeAllOptions *PurgeAllOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	return cachingApi.PurgeAllWithContext(context.Background(), purgeAllOptions)
}

// PurgeAllWithContext is an alternate form of the PurgeAll method which supports a Context parameter
func (cachingApi *CachingApiV1) PurgeAllWithContext(ctx context.Context, purgeAllOptions *PurgeAllOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(purgeAllOptions, "purgeAllOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/purge_cache/purge_all`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range purgeAllOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "PurgeAll")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPurgeAllResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// PurgeByUrls : Purge URLs
// Granularly remove one or more files from CDN edge servers' cache either by specifying URLs.
func (cachingApi *CachingApiV1) PurgeByUrls(purgeByUrlsOptions *PurgeByUrlsOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	return cachingApi.PurgeByUrlsWithContext(context.Background(), purgeByUrlsOptions)
}

// PurgeByUrlsWithContext is an alternate form of the PurgeByUrls method which supports a Context parameter
func (cachingApi *CachingApiV1) PurgeByUrlsWithContext(ctx context.Context, purgeByUrlsOptions *PurgeByUrlsOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(purgeByUrlsOptions, "purgeByUrlsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/purge_cache/purge_by_urls`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range purgeByUrlsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "PurgeByUrls")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if purgeByUrlsOptions.Files != nil {
		body["files"] = purgeByUrlsOptions.Files
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPurgeAllResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// PurgeByCacheTags : Purge Cache-Tags
// Granularly remove one or more files from CDN edge servers' cache either by specifying the associated Cache-Tags.
func (cachingApi *CachingApiV1) PurgeByCacheTags(purgeByCacheTagsOptions *PurgeByCacheTagsOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	return cachingApi.PurgeByCacheTagsWithContext(context.Background(), purgeByCacheTagsOptions)
}

// PurgeByCacheTagsWithContext is an alternate form of the PurgeByCacheTags method which supports a Context parameter
func (cachingApi *CachingApiV1) PurgeByCacheTagsWithContext(ctx context.Context, purgeByCacheTagsOptions *PurgeByCacheTagsOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(purgeByCacheTagsOptions, "purgeByCacheTagsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/purge_cache/purge_by_cache_tags`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range purgeByCacheTagsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "PurgeByCacheTags")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if purgeByCacheTagsOptions.Tags != nil {
		body["tags"] = purgeByCacheTagsOptions.Tags
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPurgeAllResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// PurgeByHosts : Purge host names
// Granularly remove one or more files from CDN edge servers' cache either by specifying the hostnames.
func (cachingApi *CachingApiV1) PurgeByHosts(purgeByHostsOptions *PurgeByHostsOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	return cachingApi.PurgeByHostsWithContext(context.Background(), purgeByHostsOptions)
}

// PurgeByHostsWithContext is an alternate form of the PurgeByHosts method which supports a Context parameter
func (cachingApi *CachingApiV1) PurgeByHostsWithContext(ctx context.Context, purgeByHostsOptions *PurgeByHostsOptions) (result *PurgeAllResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(purgeByHostsOptions, "purgeByHostsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/purge_cache/purge_by_hosts`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range purgeByHostsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "PurgeByHosts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if purgeByHostsOptions.Hosts != nil {
		body["hosts"] = purgeByHostsOptions.Hosts
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPurgeAllResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetBrowserCacheTTL : Get browser cache TTL setting
// Browser Cache TTL (in seconds) specifies how long CDN edge servers cached resources will remain on your visitors'
// computers.
func (cachingApi *CachingApiV1) GetBrowserCacheTTL(getBrowserCacheTtlOptions *GetBrowserCacheTtlOptions) (result *BrowserTTLResponse, response *core.DetailedResponse, err error) {
	return cachingApi.GetBrowserCacheTTLWithContext(context.Background(), getBrowserCacheTtlOptions)
}

// GetBrowserCacheTTLWithContext is an alternate form of the GetBrowserCacheTTL method which supports a Context parameter
func (cachingApi *CachingApiV1) GetBrowserCacheTTLWithContext(ctx context.Context, getBrowserCacheTtlOptions *GetBrowserCacheTtlOptions) (result *BrowserTTLResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getBrowserCacheTtlOptions, "getBrowserCacheTtlOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/browser_cache_ttl`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getBrowserCacheTtlOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "GetBrowserCacheTTL")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBrowserTTLResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateBrowserCacheTTL : Change browser cache TTL setting
// Browser Cache TTL (in seconds) specifies how long CDN edge servers cached resources will remain on your visitors'
// computers.
func (cachingApi *CachingApiV1) UpdateBrowserCacheTTL(updateBrowserCacheTtlOptions *UpdateBrowserCacheTtlOptions) (result *BrowserTTLResponse, response *core.DetailedResponse, err error) {
	return cachingApi.UpdateBrowserCacheTTLWithContext(context.Background(), updateBrowserCacheTtlOptions)
}

// UpdateBrowserCacheTTLWithContext is an alternate form of the UpdateBrowserCacheTTL method which supports a Context parameter
func (cachingApi *CachingApiV1) UpdateBrowserCacheTTLWithContext(ctx context.Context, updateBrowserCacheTtlOptions *UpdateBrowserCacheTtlOptions) (result *BrowserTTLResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateBrowserCacheTtlOptions, "updateBrowserCacheTtlOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/browser_cache_ttl`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateBrowserCacheTtlOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "UpdateBrowserCacheTTL")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateBrowserCacheTtlOptions.Value != nil {
		body["value"] = updateBrowserCacheTtlOptions.Value
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBrowserTTLResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetServeStaleContent : Get Serve Stale Content setting
// When enabled, Serve Stale Content will serve pages from CDN edge servers' cache if your server is offline.
func (cachingApi *CachingApiV1) GetServeStaleContent(getServeStaleContentOptions *GetServeStaleContentOptions) (result *ServeStaleContentResponse, response *core.DetailedResponse, err error) {
	return cachingApi.GetServeStaleContentWithContext(context.Background(), getServeStaleContentOptions)
}

// GetServeStaleContentWithContext is an alternate form of the GetServeStaleContent method which supports a Context parameter
func (cachingApi *CachingApiV1) GetServeStaleContentWithContext(ctx context.Context, getServeStaleContentOptions *GetServeStaleContentOptions) (result *ServeStaleContentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getServeStaleContentOptions, "getServeStaleContentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/always_online`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getServeStaleContentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "GetServeStaleContent")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServeStaleContentResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateServeStaleContent : Change Serve Stale Content setting
// When enabled, Serve Stale Content will serve pages from CDN edge servers' cache if your server is offline.
func (cachingApi *CachingApiV1) UpdateServeStaleContent(updateServeStaleContentOptions *UpdateServeStaleContentOptions) (result *ServeStaleContentResponse, response *core.DetailedResponse, err error) {
	return cachingApi.UpdateServeStaleContentWithContext(context.Background(), updateServeStaleContentOptions)
}

// UpdateServeStaleContentWithContext is an alternate form of the UpdateServeStaleContent method which supports a Context parameter
func (cachingApi *CachingApiV1) UpdateServeStaleContentWithContext(ctx context.Context, updateServeStaleContentOptions *UpdateServeStaleContentOptions) (result *ServeStaleContentResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateServeStaleContentOptions, "updateServeStaleContentOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/always_online`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateServeStaleContentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "UpdateServeStaleContent")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateServeStaleContentOptions.Value != nil {
		body["value"] = updateServeStaleContentOptions.Value
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServeStaleContentResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetDevelopmentMode : Get development mode setting
// Get development mode setting.
func (cachingApi *CachingApiV1) GetDevelopmentMode(getDevelopmentModeOptions *GetDevelopmentModeOptions) (result *DeveopmentModeResponse, response *core.DetailedResponse, err error) {
	return cachingApi.GetDevelopmentModeWithContext(context.Background(), getDevelopmentModeOptions)
}

// GetDevelopmentModeWithContext is an alternate form of the GetDevelopmentMode method which supports a Context parameter
func (cachingApi *CachingApiV1) GetDevelopmentModeWithContext(ctx context.Context, getDevelopmentModeOptions *GetDevelopmentModeOptions) (result *DeveopmentModeResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getDevelopmentModeOptions, "getDevelopmentModeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/development_mode`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getDevelopmentModeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "GetDevelopmentMode")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeveopmentModeResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateDevelopmentMode : Change development mode setting
// Change development mode setting.
func (cachingApi *CachingApiV1) UpdateDevelopmentMode(updateDevelopmentModeOptions *UpdateDevelopmentModeOptions) (result *DeveopmentModeResponse, response *core.DetailedResponse, err error) {
	return cachingApi.UpdateDevelopmentModeWithContext(context.Background(), updateDevelopmentModeOptions)
}

// UpdateDevelopmentModeWithContext is an alternate form of the UpdateDevelopmentMode method which supports a Context parameter
func (cachingApi *CachingApiV1) UpdateDevelopmentModeWithContext(ctx context.Context, updateDevelopmentModeOptions *UpdateDevelopmentModeOptions) (result *DeveopmentModeResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateDevelopmentModeOptions, "updateDevelopmentModeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/development_mode`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateDevelopmentModeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "UpdateDevelopmentMode")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateDevelopmentModeOptions.Value != nil {
		body["value"] = updateDevelopmentModeOptions.Value
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeveopmentModeResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetQueryStringSort : Get Enable Query String Sort setting
// Get Enable Query String Sort setting.
func (cachingApi *CachingApiV1) GetQueryStringSort(getQueryStringSortOptions *GetQueryStringSortOptions) (result *EnableQueryStringSortResponse, response *core.DetailedResponse, err error) {
	return cachingApi.GetQueryStringSortWithContext(context.Background(), getQueryStringSortOptions)
}

// GetQueryStringSortWithContext is an alternate form of the GetQueryStringSort method which supports a Context parameter
func (cachingApi *CachingApiV1) GetQueryStringSortWithContext(ctx context.Context, getQueryStringSortOptions *GetQueryStringSortOptions) (result *EnableQueryStringSortResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getQueryStringSortOptions, "getQueryStringSortOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/sort_query_string_for_cache`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getQueryStringSortOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "GetQueryStringSort")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnableQueryStringSortResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateQueryStringSort : Change Enable Query String Sort setting
// Change Enable Query String Sort setting.
func (cachingApi *CachingApiV1) UpdateQueryStringSort(updateQueryStringSortOptions *UpdateQueryStringSortOptions) (result *EnableQueryStringSortResponse, response *core.DetailedResponse, err error) {
	return cachingApi.UpdateQueryStringSortWithContext(context.Background(), updateQueryStringSortOptions)
}

// UpdateQueryStringSortWithContext is an alternate form of the UpdateQueryStringSort method which supports a Context parameter
func (cachingApi *CachingApiV1) UpdateQueryStringSortWithContext(ctx context.Context, updateQueryStringSortOptions *UpdateQueryStringSortOptions) (result *EnableQueryStringSortResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateQueryStringSortOptions, "updateQueryStringSortOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/sort_query_string_for_cache`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateQueryStringSortOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "UpdateQueryStringSort")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateQueryStringSortOptions.Value != nil {
		body["value"] = updateQueryStringSortOptions.Value
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnableQueryStringSortResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetCacheLevel : Get cache level setting
// Get cache level setting of a specific zone.
func (cachingApi *CachingApiV1) GetCacheLevel(getCacheLevelOptions *GetCacheLevelOptions) (result *CacheLevelResponse, response *core.DetailedResponse, err error) {
	return cachingApi.GetCacheLevelWithContext(context.Background(), getCacheLevelOptions)
}

// GetCacheLevelWithContext is an alternate form of the GetCacheLevel method which supports a Context parameter
func (cachingApi *CachingApiV1) GetCacheLevelWithContext(ctx context.Context, getCacheLevelOptions *GetCacheLevelOptions) (result *CacheLevelResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getCacheLevelOptions, "getCacheLevelOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/cache_level`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getCacheLevelOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "GetCacheLevel")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCacheLevelResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateCacheLevel : Set cache level setting
// The `basic` setting will cache most static resources (i.e., css, images, and JavaScript). The `simplified` setting
// will ignore the query string when delivering a cached resource. The `aggressive` setting will cache all static
// resources, including ones with a query string.
func (cachingApi *CachingApiV1) UpdateCacheLevel(updateCacheLevelOptions *UpdateCacheLevelOptions) (result *CacheLevelResponse, response *core.DetailedResponse, err error) {
	return cachingApi.UpdateCacheLevelWithContext(context.Background(), updateCacheLevelOptions)
}

// UpdateCacheLevelWithContext is an alternate form of the UpdateCacheLevel method which supports a Context parameter
func (cachingApi *CachingApiV1) UpdateCacheLevelWithContext(ctx context.Context, updateCacheLevelOptions *UpdateCacheLevelOptions) (result *CacheLevelResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateCacheLevelOptions, "updateCacheLevelOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *cachingApi.Crn,
		"zone_id": *cachingApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cachingApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cachingApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/settings/cache_level`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateCacheLevelOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("caching_api", "V1", "UpdateCacheLevel")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateCacheLevelOptions.Value != nil {
		body["value"] = updateCacheLevelOptions.Value
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cachingApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCacheLevelResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// BrowserTTLResponseResult : result object.
type BrowserTTLResponseResult struct {
	// ttl type.
	ID *string `json:"id,omitempty"`

	// ttl value.
	Value *int64 `json:"value,omitempty"`

	// editable.
	Editable *bool `json:"editable,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`
}


// UnmarshalBrowserTTLResponseResult unmarshals an instance of BrowserTTLResponseResult from the specified map of raw messages.
func UnmarshalBrowserTTLResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BrowserTTLResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CacheLevelResponseResult : result.
type CacheLevelResponseResult struct {
	// cache level.
	ID *string `json:"id,omitempty"`

	// cache level.
	Value *string `json:"value,omitempty"`

	// editable value.
	Editable *bool `json:"editable,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`
}


// UnmarshalCacheLevelResponseResult unmarshals an instance of CacheLevelResponseResult from the specified map of raw messages.
func UnmarshalCacheLevelResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CacheLevelResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeveopmentModeResponseResult : result object.
type DeveopmentModeResponseResult struct {
	// object id.
	ID *string `json:"id,omitempty"`

	// on/off value.
	Value *string `json:"value,omitempty"`

	// editable value.
	Editable *bool `json:"editable,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`
}


// UnmarshalDeveopmentModeResponseResult unmarshals an instance of DeveopmentModeResponseResult from the specified map of raw messages.
func UnmarshalDeveopmentModeResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeveopmentModeResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnableQueryStringSortResponseResult : result of sort query string.
type EnableQueryStringSortResponseResult struct {
	// cache id.
	ID *string `json:"id,omitempty"`

	// on/off value.
	Value *string `json:"value,omitempty"`

	// editable propery.
	Editable *bool `json:"editable,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`
}


// UnmarshalEnableQueryStringSortResponseResult unmarshals an instance of EnableQueryStringSortResponseResult from the specified map of raw messages.
func UnmarshalEnableQueryStringSortResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnableQueryStringSortResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetBrowserCacheTtlOptions : The GetBrowserCacheTTL options.
type GetBrowserCacheTtlOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetBrowserCacheTtlOptions : Instantiate GetBrowserCacheTtlOptions
func (*CachingApiV1) NewGetBrowserCacheTtlOptions() *GetBrowserCacheTtlOptions {
	return &GetBrowserCacheTtlOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetBrowserCacheTtlOptions) SetHeaders(param map[string]string) *GetBrowserCacheTtlOptions {
	options.Headers = param
	return options
}

// GetCacheLevelOptions : The GetCacheLevel options.
type GetCacheLevelOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetCacheLevelOptions : Instantiate GetCacheLevelOptions
func (*CachingApiV1) NewGetCacheLevelOptions() *GetCacheLevelOptions {
	return &GetCacheLevelOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetCacheLevelOptions) SetHeaders(param map[string]string) *GetCacheLevelOptions {
	options.Headers = param
	return options
}

// GetDevelopmentModeOptions : The GetDevelopmentMode options.
type GetDevelopmentModeOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetDevelopmentModeOptions : Instantiate GetDevelopmentModeOptions
func (*CachingApiV1) NewGetDevelopmentModeOptions() *GetDevelopmentModeOptions {
	return &GetDevelopmentModeOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetDevelopmentModeOptions) SetHeaders(param map[string]string) *GetDevelopmentModeOptions {
	options.Headers = param
	return options
}

// GetQueryStringSortOptions : The GetQueryStringSort options.
type GetQueryStringSortOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetQueryStringSortOptions : Instantiate GetQueryStringSortOptions
func (*CachingApiV1) NewGetQueryStringSortOptions() *GetQueryStringSortOptions {
	return &GetQueryStringSortOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetQueryStringSortOptions) SetHeaders(param map[string]string) *GetQueryStringSortOptions {
	options.Headers = param
	return options
}

// GetServeStaleContentOptions : The GetServeStaleContent options.
type GetServeStaleContentOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetServeStaleContentOptions : Instantiate GetServeStaleContentOptions
func (*CachingApiV1) NewGetServeStaleContentOptions() *GetServeStaleContentOptions {
	return &GetServeStaleContentOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetServeStaleContentOptions) SetHeaders(param map[string]string) *GetServeStaleContentOptions {
	options.Headers = param
	return options
}

// PurgeAllOptions : The PurgeAll options.
type PurgeAllOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPurgeAllOptions : Instantiate PurgeAllOptions
func (*CachingApiV1) NewPurgeAllOptions() *PurgeAllOptions {
	return &PurgeAllOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *PurgeAllOptions) SetHeaders(param map[string]string) *PurgeAllOptions {
	options.Headers = param
	return options
}

// PurgeAllResponseResult : purge object.
type PurgeAllResponseResult struct {
	// purge id.
	ID *string `json:"id,omitempty"`
}


// UnmarshalPurgeAllResponseResult unmarshals an instance of PurgeAllResponseResult from the specified map of raw messages.
func UnmarshalPurgeAllResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PurgeAllResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PurgeByCacheTagsOptions : The PurgeByCacheTags options.
type PurgeByCacheTagsOptions struct {
	// array of tags.
	Tags []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPurgeByCacheTagsOptions : Instantiate PurgeByCacheTagsOptions
func (*CachingApiV1) NewPurgeByCacheTagsOptions() *PurgeByCacheTagsOptions {
	return &PurgeByCacheTagsOptions{}
}

// SetTags : Allow user to set Tags
func (options *PurgeByCacheTagsOptions) SetTags(tags []string) *PurgeByCacheTagsOptions {
	options.Tags = tags
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PurgeByCacheTagsOptions) SetHeaders(param map[string]string) *PurgeByCacheTagsOptions {
	options.Headers = param
	return options
}

// PurgeByHostsOptions : The PurgeByHosts options.
type PurgeByHostsOptions struct {
	// hosts name.
	Hosts []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPurgeByHostsOptions : Instantiate PurgeByHostsOptions
func (*CachingApiV1) NewPurgeByHostsOptions() *PurgeByHostsOptions {
	return &PurgeByHostsOptions{}
}

// SetHosts : Allow user to set Hosts
func (options *PurgeByHostsOptions) SetHosts(hosts []string) *PurgeByHostsOptions {
	options.Hosts = hosts
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PurgeByHostsOptions) SetHeaders(param map[string]string) *PurgeByHostsOptions {
	options.Headers = param
	return options
}

// PurgeByUrlsOptions : The PurgeByUrls options.
type PurgeByUrlsOptions struct {
	// purge url array.
	Files []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPurgeByUrlsOptions : Instantiate PurgeByUrlsOptions
func (*CachingApiV1) NewPurgeByUrlsOptions() *PurgeByUrlsOptions {
	return &PurgeByUrlsOptions{}
}

// SetFiles : Allow user to set Files
func (options *PurgeByUrlsOptions) SetFiles(files []string) *PurgeByUrlsOptions {
	options.Files = files
	return options
}

// SetHeaders : Allow user to set Headers
func (options *PurgeByUrlsOptions) SetHeaders(param map[string]string) *PurgeByUrlsOptions {
	options.Headers = param
	return options
}

// ServeStaleContentResponseResult : result object.
type ServeStaleContentResponseResult struct {
	// serve stale content cache id.
	ID *string `json:"id,omitempty"`

	// on/off value.
	Value *string `json:"value,omitempty"`

	// editable value.
	Editable *bool `json:"editable,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`
}


// UnmarshalServeStaleContentResponseResult unmarshals an instance of ServeStaleContentResponseResult from the specified map of raw messages.
func UnmarshalServeStaleContentResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServeStaleContentResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "editable", &obj.Editable)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateBrowserCacheTtlOptions : The UpdateBrowserCacheTTL options.
type UpdateBrowserCacheTtlOptions struct {
	// ttl value.
	Value *int64

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateBrowserCacheTtlOptions : Instantiate UpdateBrowserCacheTtlOptions
func (*CachingApiV1) NewUpdateBrowserCacheTtlOptions() *UpdateBrowserCacheTtlOptions {
	return &UpdateBrowserCacheTtlOptions{}
}

// SetValue : Allow user to set Value
func (options *UpdateBrowserCacheTtlOptions) SetValue(value int64) *UpdateBrowserCacheTtlOptions {
	options.Value = core.Int64Ptr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBrowserCacheTtlOptions) SetHeaders(param map[string]string) *UpdateBrowserCacheTtlOptions {
	options.Headers = param
	return options
}

// UpdateCacheLevelOptions : The UpdateCacheLevel options.
type UpdateCacheLevelOptions struct {
	// cache level.
	Value *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateCacheLevelOptions.Value property.
// cache level.
const (
	UpdateCacheLevelOptions_Value_Aggressive = "aggressive"
	UpdateCacheLevelOptions_Value_Basic = "basic"
	UpdateCacheLevelOptions_Value_Simplified = "simplified"
)

// NewUpdateCacheLevelOptions : Instantiate UpdateCacheLevelOptions
func (*CachingApiV1) NewUpdateCacheLevelOptions() *UpdateCacheLevelOptions {
	return &UpdateCacheLevelOptions{}
}

// SetValue : Allow user to set Value
func (options *UpdateCacheLevelOptions) SetValue(value string) *UpdateCacheLevelOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCacheLevelOptions) SetHeaders(param map[string]string) *UpdateCacheLevelOptions {
	options.Headers = param
	return options
}

// UpdateDevelopmentModeOptions : The UpdateDevelopmentMode options.
type UpdateDevelopmentModeOptions struct {
	// on/off value.
	Value *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateDevelopmentModeOptions.Value property.
// on/off value.
const (
	UpdateDevelopmentModeOptions_Value_Off = "off"
	UpdateDevelopmentModeOptions_Value_On = "on"
)

// NewUpdateDevelopmentModeOptions : Instantiate UpdateDevelopmentModeOptions
func (*CachingApiV1) NewUpdateDevelopmentModeOptions() *UpdateDevelopmentModeOptions {
	return &UpdateDevelopmentModeOptions{}
}

// SetValue : Allow user to set Value
func (options *UpdateDevelopmentModeOptions) SetValue(value string) *UpdateDevelopmentModeOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDevelopmentModeOptions) SetHeaders(param map[string]string) *UpdateDevelopmentModeOptions {
	options.Headers = param
	return options
}

// UpdateQueryStringSortOptions : The UpdateQueryStringSort options.
type UpdateQueryStringSortOptions struct {
	// on/off property value.
	Value *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateQueryStringSortOptions.Value property.
// on/off property value.
const (
	UpdateQueryStringSortOptions_Value_Off = "off"
	UpdateQueryStringSortOptions_Value_On = "on"
)

// NewUpdateQueryStringSortOptions : Instantiate UpdateQueryStringSortOptions
func (*CachingApiV1) NewUpdateQueryStringSortOptions() *UpdateQueryStringSortOptions {
	return &UpdateQueryStringSortOptions{}
}

// SetValue : Allow user to set Value
func (options *UpdateQueryStringSortOptions) SetValue(value string) *UpdateQueryStringSortOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateQueryStringSortOptions) SetHeaders(param map[string]string) *UpdateQueryStringSortOptions {
	options.Headers = param
	return options
}

// UpdateServeStaleContentOptions : The UpdateServeStaleContent options.
type UpdateServeStaleContentOptions struct {
	// on/off value.
	Value *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateServeStaleContentOptions.Value property.
// on/off value.
const (
	UpdateServeStaleContentOptions_Value_Off = "off"
	UpdateServeStaleContentOptions_Value_On = "on"
)

// NewUpdateServeStaleContentOptions : Instantiate UpdateServeStaleContentOptions
func (*CachingApiV1) NewUpdateServeStaleContentOptions() *UpdateServeStaleContentOptions {
	return &UpdateServeStaleContentOptions{}
}

// SetValue : Allow user to set Value
func (options *UpdateServeStaleContentOptions) SetValue(value string) *UpdateServeStaleContentOptions {
	options.Value = core.StringPtr(value)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateServeStaleContentOptions) SetHeaders(param map[string]string) *UpdateServeStaleContentOptions {
	options.Headers = param
	return options
}

// BrowserTTLResponse : browser ttl response.
type BrowserTTLResponse struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result object.
	Result *BrowserTTLResponseResult `json:"result" validate:"required"`
}


// UnmarshalBrowserTTLResponse unmarshals an instance of BrowserTTLResponse from the specified map of raw messages.
func UnmarshalBrowserTTLResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BrowserTTLResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalBrowserTTLResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CacheLevelResponse : cache level response.
type CacheLevelResponse struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result *CacheLevelResponseResult `json:"result" validate:"required"`
}


// UnmarshalCacheLevelResponse unmarshals an instance of CacheLevelResponse from the specified map of raw messages.
func UnmarshalCacheLevelResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CacheLevelResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCacheLevelResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeveopmentModeResponse : development mode response.
type DeveopmentModeResponse struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result object.
	Result *DeveopmentModeResponseResult `json:"result" validate:"required"`
}


// UnmarshalDeveopmentModeResponse unmarshals an instance of DeveopmentModeResponse from the specified map of raw messages.
func UnmarshalDeveopmentModeResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeveopmentModeResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeveopmentModeResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnableQueryStringSortResponse : sort query string response.
type EnableQueryStringSortResponse struct {
	// success response true/false.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result of sort query string.
	Result *EnableQueryStringSortResponseResult `json:"result" validate:"required"`
}


// UnmarshalEnableQueryStringSortResponse unmarshals an instance of EnableQueryStringSortResponse from the specified map of raw messages.
func UnmarshalEnableQueryStringSortResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnableQueryStringSortResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEnableQueryStringSortResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PurgeAllResponse : purge all response.
type PurgeAllResponse struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// purge object.
	Result *PurgeAllResponseResult `json:"result" validate:"required"`
}


// UnmarshalPurgeAllResponse unmarshals an instance of PurgeAllResponse from the specified map of raw messages.
func UnmarshalPurgeAllResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PurgeAllResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalPurgeAllResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServeStaleContentResponse : serve stale conent response.
type ServeStaleContentResponse struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result object.
	Result *ServeStaleContentResponseResult `json:"result" validate:"required"`
}


// UnmarshalServeStaleContentResponse unmarshals an instance of ServeStaleContentResponse from the specified map of raw messages.
func UnmarshalServeStaleContentResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServeStaleContentResponse)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalServeStaleContentResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
