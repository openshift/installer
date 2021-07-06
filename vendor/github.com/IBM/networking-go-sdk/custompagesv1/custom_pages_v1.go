/**
 * (C) Copyright IBM Corp. 2020.
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
 * IBM OpenAPI SDK Code Generator Version: 3.20.0-debb9f29-20201203-202043
 */
 

// Package custompagesv1 : Operations and models for the CustomPagesV1 service
package custompagesv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
	"net/http"
	"reflect"
	"time"
)

// CustomPagesV1 : Custom Pages
//
// Version: 1.0.0
type CustomPagesV1 struct {
	Service *core.BaseService

	// Full crn of the service instance.
	Crn *string

	// Zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "custom_pages"

// CustomPagesV1Options : Service options
type CustomPagesV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full crn of the service instance.
	Crn *string `validate:"required"`

	// Zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewCustomPagesV1UsingExternalConfig : constructs an instance of CustomPagesV1 with passed in options and external configuration.
func NewCustomPagesV1UsingExternalConfig(options *CustomPagesV1Options) (customPages *CustomPagesV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	customPages, err = NewCustomPagesV1(options)
	if err != nil {
		return
	}

	err = customPages.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = customPages.Service.SetServiceURL(options.URL)
	}
	return
}

// NewCustomPagesV1 : constructs an instance of CustomPagesV1 with passed in options.
func NewCustomPagesV1(options *CustomPagesV1Options) (service *CustomPagesV1, err error) {
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

	service = &CustomPagesV1{
		Service: baseService,
		Crn: options.Crn,
		ZoneIdentifier: options.ZoneIdentifier,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "customPages" suitable for processing requests.
func (customPages *CustomPagesV1) Clone() *CustomPagesV1 {
	if core.IsNil(customPages) {
		return nil
	}
	clone := *customPages
	clone.Service = customPages.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (customPages *CustomPagesV1) SetServiceURL(url string) error {
	return customPages.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (customPages *CustomPagesV1) GetServiceURL() string {
	return customPages.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (customPages *CustomPagesV1) SetDefaultHeaders(headers http.Header) {
	customPages.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (customPages *CustomPagesV1) SetEnableGzipCompression(enableGzip bool) {
	customPages.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (customPages *CustomPagesV1) GetEnableGzipCompression() bool {
	return customPages.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (customPages *CustomPagesV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	customPages.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (customPages *CustomPagesV1) DisableRetries() {
	customPages.Service.DisableRetries()
}

// ListInstanceCustomPages : List all custom pages for a given instance
// List all custom pages for a given instance.
func (customPages *CustomPagesV1) ListInstanceCustomPages(listInstanceCustomPagesOptions *ListInstanceCustomPagesOptions) (result *ListCustomPagesResp, response *core.DetailedResponse, err error) {
	return customPages.ListInstanceCustomPagesWithContext(context.Background(), listInstanceCustomPagesOptions)
}

// ListInstanceCustomPagesWithContext is an alternate form of the ListInstanceCustomPages method which supports a Context parameter
func (customPages *CustomPagesV1) ListInstanceCustomPagesWithContext(ctx context.Context, listInstanceCustomPagesOptions *ListInstanceCustomPagesOptions) (result *ListCustomPagesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listInstanceCustomPagesOptions, "listInstanceCustomPagesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *customPages.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = customPages.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(customPages.Service.Options.URL, `/v1/{crn}/custom_pages`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listInstanceCustomPagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("custom_pages", "V1", "ListInstanceCustomPages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = customPages.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListCustomPagesResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetInstanceCustomPage : Get a custom page for a given instance
// Get a specific custom page for a given instance.
func (customPages *CustomPagesV1) GetInstanceCustomPage(getInstanceCustomPageOptions *GetInstanceCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	return customPages.GetInstanceCustomPageWithContext(context.Background(), getInstanceCustomPageOptions)
}

// GetInstanceCustomPageWithContext is an alternate form of the GetInstanceCustomPage method which supports a Context parameter
func (customPages *CustomPagesV1) GetInstanceCustomPageWithContext(ctx context.Context, getInstanceCustomPageOptions *GetInstanceCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getInstanceCustomPageOptions, "getInstanceCustomPageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getInstanceCustomPageOptions, "getInstanceCustomPageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *customPages.Crn,
		"page_identifier": *getInstanceCustomPageOptions.PageIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = customPages.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(customPages.Service.Options.URL, `/v1/{crn}/custom_pages/{page_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getInstanceCustomPageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("custom_pages", "V1", "GetInstanceCustomPage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = customPages.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomPageSpecificResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateInstanceCustomPage : Update a custom page for a given instance
// Update a specific custom page for a given instance.
func (customPages *CustomPagesV1) UpdateInstanceCustomPage(updateInstanceCustomPageOptions *UpdateInstanceCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	return customPages.UpdateInstanceCustomPageWithContext(context.Background(), updateInstanceCustomPageOptions)
}

// UpdateInstanceCustomPageWithContext is an alternate form of the UpdateInstanceCustomPage method which supports a Context parameter
func (customPages *CustomPagesV1) UpdateInstanceCustomPageWithContext(ctx context.Context, updateInstanceCustomPageOptions *UpdateInstanceCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateInstanceCustomPageOptions, "updateInstanceCustomPageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateInstanceCustomPageOptions, "updateInstanceCustomPageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *customPages.Crn,
		"page_identifier": *updateInstanceCustomPageOptions.PageIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = customPages.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(customPages.Service.Options.URL, `/v1/{crn}/custom_pages/{page_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateInstanceCustomPageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("custom_pages", "V1", "UpdateInstanceCustomPage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateInstanceCustomPageOptions.URL != nil {
		body["url"] = updateInstanceCustomPageOptions.URL
	}
	if updateInstanceCustomPageOptions.State != nil {
		body["state"] = updateInstanceCustomPageOptions.State
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
	response, err = customPages.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomPageSpecificResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListZoneCustomPages : List all custom pages for a given zone
// List all custom pages for a given zone.
func (customPages *CustomPagesV1) ListZoneCustomPages(listZoneCustomPagesOptions *ListZoneCustomPagesOptions) (result *ListCustomPagesResp, response *core.DetailedResponse, err error) {
	return customPages.ListZoneCustomPagesWithContext(context.Background(), listZoneCustomPagesOptions)
}

// ListZoneCustomPagesWithContext is an alternate form of the ListZoneCustomPages method which supports a Context parameter
func (customPages *CustomPagesV1) ListZoneCustomPagesWithContext(ctx context.Context, listZoneCustomPagesOptions *ListZoneCustomPagesOptions) (result *ListCustomPagesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listZoneCustomPagesOptions, "listZoneCustomPagesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *customPages.Crn,
		"zone_identifier": *customPages.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = customPages.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(customPages.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_pages`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listZoneCustomPagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("custom_pages", "V1", "ListZoneCustomPages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = customPages.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListCustomPagesResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetZoneCustomPage : Get a custom page for a given zone
// Get a specific custom page for a given zone.
func (customPages *CustomPagesV1) GetZoneCustomPage(getZoneCustomPageOptions *GetZoneCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	return customPages.GetZoneCustomPageWithContext(context.Background(), getZoneCustomPageOptions)
}

// GetZoneCustomPageWithContext is an alternate form of the GetZoneCustomPage method which supports a Context parameter
func (customPages *CustomPagesV1) GetZoneCustomPageWithContext(ctx context.Context, getZoneCustomPageOptions *GetZoneCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneCustomPageOptions, "getZoneCustomPageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getZoneCustomPageOptions, "getZoneCustomPageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *customPages.Crn,
		"zone_identifier": *customPages.ZoneIdentifier,
		"page_identifier": *getZoneCustomPageOptions.PageIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = customPages.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(customPages.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_pages/{page_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getZoneCustomPageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("custom_pages", "V1", "GetZoneCustomPage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = customPages.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomPageSpecificResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateZoneCustomPage : Update a custom page for a given zone
// Update a specific custom page for a given zone.
func (customPages *CustomPagesV1) UpdateZoneCustomPage(updateZoneCustomPageOptions *UpdateZoneCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	return customPages.UpdateZoneCustomPageWithContext(context.Background(), updateZoneCustomPageOptions)
}

// UpdateZoneCustomPageWithContext is an alternate form of the UpdateZoneCustomPage method which supports a Context parameter
func (customPages *CustomPagesV1) UpdateZoneCustomPageWithContext(ctx context.Context, updateZoneCustomPageOptions *UpdateZoneCustomPageOptions) (result *CustomPageSpecificResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateZoneCustomPageOptions, "updateZoneCustomPageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateZoneCustomPageOptions, "updateZoneCustomPageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *customPages.Crn,
		"zone_identifier": *customPages.ZoneIdentifier,
		"page_identifier": *updateZoneCustomPageOptions.PageIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = customPages.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(customPages.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/custom_pages/{page_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateZoneCustomPageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("custom_pages", "V1", "UpdateZoneCustomPage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateZoneCustomPageOptions.URL != nil {
		body["url"] = updateZoneCustomPageOptions.URL
	}
	if updateZoneCustomPageOptions.State != nil {
		body["state"] = updateZoneCustomPageOptions.State
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
	response, err = customPages.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCustomPageSpecificResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetInstanceCustomPageOptions : The GetInstanceCustomPage options.
type GetInstanceCustomPageOptions struct {
	// Custom page identifier.
	PageIdentifier *string `json:"page_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetInstanceCustomPageOptions.PageIdentifier property.
// Custom page identifier.
const (
	GetInstanceCustomPageOptions_PageIdentifier_AlwaysOnline = "always_online"
	GetInstanceCustomPageOptions_PageIdentifier_BasicChallenge = "basic_challenge"
	GetInstanceCustomPageOptions_PageIdentifier_CountryChallenge = "country_challenge"
	GetInstanceCustomPageOptions_PageIdentifier_IpBlock = "ip_block"
	GetInstanceCustomPageOptions_PageIdentifier_RatelimitBlock = "ratelimit_block"
	GetInstanceCustomPageOptions_PageIdentifier_UnderAttack = "under_attack"
	GetInstanceCustomPageOptions_PageIdentifier_WafBlock = "waf_block"
	GetInstanceCustomPageOptions_PageIdentifier_WafChallenge = "waf_challenge"
)

// NewGetInstanceCustomPageOptions : Instantiate GetInstanceCustomPageOptions
func (*CustomPagesV1) NewGetInstanceCustomPageOptions(pageIdentifier string) *GetInstanceCustomPageOptions {
	return &GetInstanceCustomPageOptions{
		PageIdentifier: core.StringPtr(pageIdentifier),
	}
}

// SetPageIdentifier : Allow user to set PageIdentifier
func (options *GetInstanceCustomPageOptions) SetPageIdentifier(pageIdentifier string) *GetInstanceCustomPageOptions {
	options.PageIdentifier = core.StringPtr(pageIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetInstanceCustomPageOptions) SetHeaders(param map[string]string) *GetInstanceCustomPageOptions {
	options.Headers = param
	return options
}

// GetZoneCustomPageOptions : The GetZoneCustomPage options.
type GetZoneCustomPageOptions struct {
	// Custom page identifier.
	PageIdentifier *string `json:"page_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetZoneCustomPageOptions.PageIdentifier property.
// Custom page identifier.
const (
	GetZoneCustomPageOptions_PageIdentifier_AlwaysOnline = "always_online"
	GetZoneCustomPageOptions_PageIdentifier_BasicChallenge = "basic_challenge"
	GetZoneCustomPageOptions_PageIdentifier_CountryChallenge = "country_challenge"
	GetZoneCustomPageOptions_PageIdentifier_IpBlock = "ip_block"
	GetZoneCustomPageOptions_PageIdentifier_RatelimitBlock = "ratelimit_block"
	GetZoneCustomPageOptions_PageIdentifier_UnderAttack = "under_attack"
	GetZoneCustomPageOptions_PageIdentifier_WafBlock = "waf_block"
	GetZoneCustomPageOptions_PageIdentifier_WafChallenge = "waf_challenge"
)

// NewGetZoneCustomPageOptions : Instantiate GetZoneCustomPageOptions
func (*CustomPagesV1) NewGetZoneCustomPageOptions(pageIdentifier string) *GetZoneCustomPageOptions {
	return &GetZoneCustomPageOptions{
		PageIdentifier: core.StringPtr(pageIdentifier),
	}
}

// SetPageIdentifier : Allow user to set PageIdentifier
func (options *GetZoneCustomPageOptions) SetPageIdentifier(pageIdentifier string) *GetZoneCustomPageOptions {
	options.PageIdentifier = core.StringPtr(pageIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneCustomPageOptions) SetHeaders(param map[string]string) *GetZoneCustomPageOptions {
	options.Headers = param
	return options
}

// ListCustomPagesRespResultInfo : Statistics of results.
type ListCustomPagesRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of total pages.
	TotalPages *int64 `json:"total_pages" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalListCustomPagesRespResultInfo unmarshals an instance of ListCustomPagesRespResultInfo from the specified map of raw messages.
func UnmarshalListCustomPagesRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListCustomPagesRespResultInfo)
	err = core.UnmarshalPrimitive(m, "page", &obj.Page)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "per_page", &obj.PerPage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_pages", &obj.TotalPages)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListInstanceCustomPagesOptions : The ListInstanceCustomPages options.
type ListInstanceCustomPagesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListInstanceCustomPagesOptions : Instantiate ListInstanceCustomPagesOptions
func (*CustomPagesV1) NewListInstanceCustomPagesOptions() *ListInstanceCustomPagesOptions {
	return &ListInstanceCustomPagesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListInstanceCustomPagesOptions) SetHeaders(param map[string]string) *ListInstanceCustomPagesOptions {
	options.Headers = param
	return options
}

// ListZoneCustomPagesOptions : The ListZoneCustomPages options.
type ListZoneCustomPagesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListZoneCustomPagesOptions : Instantiate ListZoneCustomPagesOptions
func (*CustomPagesV1) NewListZoneCustomPagesOptions() *ListZoneCustomPagesOptions {
	return &ListZoneCustomPagesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListZoneCustomPagesOptions) SetHeaders(param map[string]string) *ListZoneCustomPagesOptions {
	options.Headers = param
	return options
}

// UpdateInstanceCustomPageOptions : The UpdateInstanceCustomPage options.
type UpdateInstanceCustomPageOptions struct {
	// Custom page identifier.
	PageIdentifier *string `json:"page_identifier" validate:"required,ne="`

	// A URL that is associated with the Custom Page.
	URL *string `json:"url,omitempty"`

	// The Custom Page state.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateInstanceCustomPageOptions.PageIdentifier property.
// Custom page identifier.
const (
	UpdateInstanceCustomPageOptions_PageIdentifier_AlwaysOnline = "always_online"
	UpdateInstanceCustomPageOptions_PageIdentifier_BasicChallenge = "basic_challenge"
	UpdateInstanceCustomPageOptions_PageIdentifier_CountryChallenge = "country_challenge"
	UpdateInstanceCustomPageOptions_PageIdentifier_IpBlock = "ip_block"
	UpdateInstanceCustomPageOptions_PageIdentifier_RatelimitBlock = "ratelimit_block"
	UpdateInstanceCustomPageOptions_PageIdentifier_UnderAttack = "under_attack"
	UpdateInstanceCustomPageOptions_PageIdentifier_WafBlock = "waf_block"
	UpdateInstanceCustomPageOptions_PageIdentifier_WafChallenge = "waf_challenge"
)

// Constants associated with the UpdateInstanceCustomPageOptions.State property.
// The Custom Page state.
const (
	UpdateInstanceCustomPageOptions_State_Customized = "customized"
	UpdateInstanceCustomPageOptions_State_Default = "default"
)

// NewUpdateInstanceCustomPageOptions : Instantiate UpdateInstanceCustomPageOptions
func (*CustomPagesV1) NewUpdateInstanceCustomPageOptions(pageIdentifier string) *UpdateInstanceCustomPageOptions {
	return &UpdateInstanceCustomPageOptions{
		PageIdentifier: core.StringPtr(pageIdentifier),
	}
}

// SetPageIdentifier : Allow user to set PageIdentifier
func (options *UpdateInstanceCustomPageOptions) SetPageIdentifier(pageIdentifier string) *UpdateInstanceCustomPageOptions {
	options.PageIdentifier = core.StringPtr(pageIdentifier)
	return options
}

// SetURL : Allow user to set URL
func (options *UpdateInstanceCustomPageOptions) SetURL(url string) *UpdateInstanceCustomPageOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetState : Allow user to set State
func (options *UpdateInstanceCustomPageOptions) SetState(state string) *UpdateInstanceCustomPageOptions {
	options.State = core.StringPtr(state)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateInstanceCustomPageOptions) SetHeaders(param map[string]string) *UpdateInstanceCustomPageOptions {
	options.Headers = param
	return options
}

// UpdateZoneCustomPageOptions : The UpdateZoneCustomPage options.
type UpdateZoneCustomPageOptions struct {
	// Custom page identifier.
	PageIdentifier *string `json:"page_identifier" validate:"required,ne="`

	// A URL that is associated with the Custom Page.
	URL *string `json:"url,omitempty"`

	// The Custom Page state.
	State *string `json:"state,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateZoneCustomPageOptions.PageIdentifier property.
// Custom page identifier.
const (
	UpdateZoneCustomPageOptions_PageIdentifier_AlwaysOnline = "always_online"
	UpdateZoneCustomPageOptions_PageIdentifier_BasicChallenge = "basic_challenge"
	UpdateZoneCustomPageOptions_PageIdentifier_CountryChallenge = "country_challenge"
	UpdateZoneCustomPageOptions_PageIdentifier_IpBlock = "ip_block"
	UpdateZoneCustomPageOptions_PageIdentifier_RatelimitBlock = "ratelimit_block"
	UpdateZoneCustomPageOptions_PageIdentifier_UnderAttack = "under_attack"
	UpdateZoneCustomPageOptions_PageIdentifier_WafBlock = "waf_block"
	UpdateZoneCustomPageOptions_PageIdentifier_WafChallenge = "waf_challenge"
)

// Constants associated with the UpdateZoneCustomPageOptions.State property.
// The Custom Page state.
const (
	UpdateZoneCustomPageOptions_State_Customized = "customized"
	UpdateZoneCustomPageOptions_State_Default = "default"
)

// NewUpdateZoneCustomPageOptions : Instantiate UpdateZoneCustomPageOptions
func (*CustomPagesV1) NewUpdateZoneCustomPageOptions(pageIdentifier string) *UpdateZoneCustomPageOptions {
	return &UpdateZoneCustomPageOptions{
		PageIdentifier: core.StringPtr(pageIdentifier),
	}
}

// SetPageIdentifier : Allow user to set PageIdentifier
func (options *UpdateZoneCustomPageOptions) SetPageIdentifier(pageIdentifier string) *UpdateZoneCustomPageOptions {
	options.PageIdentifier = core.StringPtr(pageIdentifier)
	return options
}

// SetURL : Allow user to set URL
func (options *UpdateZoneCustomPageOptions) SetURL(url string) *UpdateZoneCustomPageOptions {
	options.URL = core.StringPtr(url)
	return options
}

// SetState : Allow user to set State
func (options *UpdateZoneCustomPageOptions) SetState(state string) *UpdateZoneCustomPageOptions {
	options.State = core.StringPtr(state)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateZoneCustomPageOptions) SetHeaders(param map[string]string) *UpdateZoneCustomPageOptions {
	options.Headers = param
	return options
}

// CustomPageObject : custom page object.
type CustomPageObject struct {
	// Custom page identifier.
	ID *string `json:"id" validate:"required"`

	// Description of custom page.
	Description *string `json:"description" validate:"required"`

	// array of page tokens.
	RequiredTokens []string `json:"required_tokens" validate:"required"`

	// Preview target.
	PreviewTarget *string `json:"preview_target" validate:"required"`

	// Created date.
	CreatedOn *strfmt.DateTime `json:"created_on" validate:"required"`

	// Modified date.
	ModifiedOn *strfmt.DateTime `json:"modified_on" validate:"required"`

	// A URL that is associated with the Custom Page.
	URL *string `json:"url" validate:"required"`

	// The Custom Page state.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the CustomPageObject.ID property.
// Custom page identifier.
const (
	CustomPageObject_ID_AlwaysOnline = "always_online"
	CustomPageObject_ID_BasicChallenge = "basic_challenge"
	CustomPageObject_ID_CountryChallenge = "country_challenge"
	CustomPageObject_ID_IpBlock = "ip_block"
	CustomPageObject_ID_RatelimitBlock = "ratelimit_block"
	CustomPageObject_ID_UnderAttack = "under_attack"
	CustomPageObject_ID_WafBlock = "waf_block"
	CustomPageObject_ID_WafChallenge = "waf_challenge"
)

// Constants associated with the CustomPageObject.RequiredTokens property.
const (
	CustomPageObject_RequiredTokens_AlwaysOnlineNoCopyBox = "::ALWAYS_ONLINE_NO_COPY_BOX::"
	CustomPageObject_RequiredTokens_CaptchaBox = "::CAPTCHA_BOX::"
	CustomPageObject_RequiredTokens_CloudflareError1000sBox = "::CLOUDFLARE_ERROR_1000S_BOX::"
	CustomPageObject_RequiredTokens_CloudflareError500sBox = "::CLOUDFLARE_ERROR_500S_BOX::"
	CustomPageObject_RequiredTokens_ImUnderAttackBox = "::IM_UNDER_ATTACK_BOX::"
)

// Constants associated with the CustomPageObject.State property.
// The Custom Page state.
const (
	CustomPageObject_State_Customized = "customized"
	CustomPageObject_State_Default = "default"
)


// UnmarshalCustomPageObject unmarshals an instance of CustomPageObject from the specified map of raw messages.
func UnmarshalCustomPageObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomPageObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "required_tokens", &obj.RequiredTokens)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "preview_target", &obj.PreviewTarget)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CustomPageSpecificResp : custom page specific response.
type CustomPageSpecificResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// custom page object.
	Result *CustomPageObject `json:"result" validate:"required"`
}


// UnmarshalCustomPageSpecificResp unmarshals an instance of CustomPageSpecificResp from the specified map of raw messages.
func UnmarshalCustomPageSpecificResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CustomPageSpecificResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCustomPageObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListCustomPagesResp : list of custom pages response.
type ListCustomPagesResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// custom pages array.
	Result []CustomPageObject `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListCustomPagesRespResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListCustomPagesResp unmarshals an instance of ListCustomPagesResp from the specified map of raw messages.
func UnmarshalListCustomPagesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListCustomPagesResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalCustomPageObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListCustomPagesRespResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
