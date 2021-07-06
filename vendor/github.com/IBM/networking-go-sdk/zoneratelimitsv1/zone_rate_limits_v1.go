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
 

// Package zoneratelimitsv1 : Operations and models for the ZoneRateLimitsV1 service
package zoneratelimitsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"net/http"
	"reflect"
	"time"
)

// ZoneRateLimitsV1 : Zone Rate Limits
//
// Version: 1.0.1
type ZoneRateLimitsV1 struct {
	Service *core.BaseService

	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string

	// Zone identifier of the zone for which rate limit is to be created.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "zone_rate_limits"

// ZoneRateLimitsV1Options : Service options
type ZoneRateLimitsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `validate:"required"`

	// Zone identifier of the zone for which rate limit is to be created.
	ZoneIdentifier *string `validate:"required"`
}

// NewZoneRateLimitsV1UsingExternalConfig : constructs an instance of ZoneRateLimitsV1 with passed in options and external configuration.
func NewZoneRateLimitsV1UsingExternalConfig(options *ZoneRateLimitsV1Options) (zoneRateLimits *ZoneRateLimitsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	zoneRateLimits, err = NewZoneRateLimitsV1(options)
	if err != nil {
		return
	}

	err = zoneRateLimits.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = zoneRateLimits.Service.SetServiceURL(options.URL)
	}
	return
}

// NewZoneRateLimitsV1 : constructs an instance of ZoneRateLimitsV1 with passed in options.
func NewZoneRateLimitsV1(options *ZoneRateLimitsV1Options) (service *ZoneRateLimitsV1, err error) {
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

	service = &ZoneRateLimitsV1{
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

// Clone makes a copy of "zoneRateLimits" suitable for processing requests.
func (zoneRateLimits *ZoneRateLimitsV1) Clone() *ZoneRateLimitsV1 {
	if core.IsNil(zoneRateLimits) {
		return nil
	}
	clone := *zoneRateLimits
	clone.Service = zoneRateLimits.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (zoneRateLimits *ZoneRateLimitsV1) SetServiceURL(url string) error {
	return zoneRateLimits.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (zoneRateLimits *ZoneRateLimitsV1) GetServiceURL() string {
	return zoneRateLimits.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (zoneRateLimits *ZoneRateLimitsV1) SetDefaultHeaders(headers http.Header) {
	zoneRateLimits.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (zoneRateLimits *ZoneRateLimitsV1) SetEnableGzipCompression(enableGzip bool) {
	zoneRateLimits.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (zoneRateLimits *ZoneRateLimitsV1) GetEnableGzipCompression() bool {
	return zoneRateLimits.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (zoneRateLimits *ZoneRateLimitsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	zoneRateLimits.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (zoneRateLimits *ZoneRateLimitsV1) DisableRetries() {
	zoneRateLimits.Service.DisableRetries()
}

// ListAllZoneRateLimits : List all rate limits
// The details of Rate Limit for a given zone under a given service instance.
func (zoneRateLimits *ZoneRateLimitsV1) ListAllZoneRateLimits(listAllZoneRateLimitsOptions *ListAllZoneRateLimitsOptions) (result *ListRatelimitResp, response *core.DetailedResponse, err error) {
	return zoneRateLimits.ListAllZoneRateLimitsWithContext(context.Background(), listAllZoneRateLimitsOptions)
}

// ListAllZoneRateLimitsWithContext is an alternate form of the ListAllZoneRateLimits method which supports a Context parameter
func (zoneRateLimits *ZoneRateLimitsV1) ListAllZoneRateLimitsWithContext(ctx context.Context, listAllZoneRateLimitsOptions *ListAllZoneRateLimitsOptions) (result *ListRatelimitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllZoneRateLimitsOptions, "listAllZoneRateLimitsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneRateLimits.Crn,
		"zone_identifier": *zoneRateLimits.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneRateLimits.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneRateLimits.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rate_limits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllZoneRateLimitsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_rate_limits", "V1", "ListAllZoneRateLimits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllZoneRateLimitsOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listAllZoneRateLimitsOptions.Page))
	}
	if listAllZoneRateLimitsOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listAllZoneRateLimitsOptions.PerPage))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneRateLimits.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListRatelimitResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneRateLimits : Create rate limit
// Create a new rate limit for a given zone under a service instance.
func (zoneRateLimits *ZoneRateLimitsV1) CreateZoneRateLimits(createZoneRateLimitsOptions *CreateZoneRateLimitsOptions) (result *RatelimitResp, response *core.DetailedResponse, err error) {
	return zoneRateLimits.CreateZoneRateLimitsWithContext(context.Background(), createZoneRateLimitsOptions)
}

// CreateZoneRateLimitsWithContext is an alternate form of the CreateZoneRateLimits method which supports a Context parameter
func (zoneRateLimits *ZoneRateLimitsV1) CreateZoneRateLimitsWithContext(ctx context.Context, createZoneRateLimitsOptions *CreateZoneRateLimitsOptions) (result *RatelimitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createZoneRateLimitsOptions, "createZoneRateLimitsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneRateLimits.Crn,
		"zone_identifier": *zoneRateLimits.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneRateLimits.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneRateLimits.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rate_limits`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createZoneRateLimitsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_rate_limits", "V1", "CreateZoneRateLimits")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createZoneRateLimitsOptions.Disabled != nil {
		body["disabled"] = createZoneRateLimitsOptions.Disabled
	}
	if createZoneRateLimitsOptions.Description != nil {
		body["description"] = createZoneRateLimitsOptions.Description
	}
	if createZoneRateLimitsOptions.Bypass != nil {
		body["bypass"] = createZoneRateLimitsOptions.Bypass
	}
	if createZoneRateLimitsOptions.Threshold != nil {
		body["threshold"] = createZoneRateLimitsOptions.Threshold
	}
	if createZoneRateLimitsOptions.Period != nil {
		body["period"] = createZoneRateLimitsOptions.Period
	}
	if createZoneRateLimitsOptions.Action != nil {
		body["action"] = createZoneRateLimitsOptions.Action
	}
	if createZoneRateLimitsOptions.Correlate != nil {
		body["correlate"] = createZoneRateLimitsOptions.Correlate
	}
	if createZoneRateLimitsOptions.Match != nil {
		body["match"] = createZoneRateLimitsOptions.Match
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
	response, err = zoneRateLimits.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRatelimitResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteZoneRateLimit : Delete rate limit
// Delete a rate limit given its id.
func (zoneRateLimits *ZoneRateLimitsV1) DeleteZoneRateLimit(deleteZoneRateLimitOptions *DeleteZoneRateLimitOptions) (result *DeleteRateLimitResp, response *core.DetailedResponse, err error) {
	return zoneRateLimits.DeleteZoneRateLimitWithContext(context.Background(), deleteZoneRateLimitOptions)
}

// DeleteZoneRateLimitWithContext is an alternate form of the DeleteZoneRateLimit method which supports a Context parameter
func (zoneRateLimits *ZoneRateLimitsV1) DeleteZoneRateLimitWithContext(ctx context.Context, deleteZoneRateLimitOptions *DeleteZoneRateLimitOptions) (result *DeleteRateLimitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneRateLimitOptions, "deleteZoneRateLimitOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneRateLimitOptions, "deleteZoneRateLimitOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneRateLimits.Crn,
		"zone_identifier": *zoneRateLimits.ZoneIdentifier,
		"rate_limit_identifier": *deleteZoneRateLimitOptions.RateLimitIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneRateLimits.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneRateLimits.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rate_limits/{rate_limit_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneRateLimitOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_rate_limits", "V1", "DeleteZoneRateLimit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneRateLimits.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteRateLimitResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetRateLimit : Get a rate limit
// Get the details of a rate limit for a given zone under a given service instance.
func (zoneRateLimits *ZoneRateLimitsV1) GetRateLimit(getRateLimitOptions *GetRateLimitOptions) (result *RatelimitResp, response *core.DetailedResponse, err error) {
	return zoneRateLimits.GetRateLimitWithContext(context.Background(), getRateLimitOptions)
}

// GetRateLimitWithContext is an alternate form of the GetRateLimit method which supports a Context parameter
func (zoneRateLimits *ZoneRateLimitsV1) GetRateLimitWithContext(ctx context.Context, getRateLimitOptions *GetRateLimitOptions) (result *RatelimitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRateLimitOptions, "getRateLimitOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRateLimitOptions, "getRateLimitOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneRateLimits.Crn,
		"zone_identifier": *zoneRateLimits.ZoneIdentifier,
		"rate_limit_identifier": *getRateLimitOptions.RateLimitIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneRateLimits.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneRateLimits.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rate_limits/{rate_limit_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRateLimitOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_rate_limits", "V1", "GetRateLimit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneRateLimits.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRatelimitResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateRateLimit : Update rate limit
// Update an existing rate limit for a given zone under a service instance.
func (zoneRateLimits *ZoneRateLimitsV1) UpdateRateLimit(updateRateLimitOptions *UpdateRateLimitOptions) (result *RatelimitResp, response *core.DetailedResponse, err error) {
	return zoneRateLimits.UpdateRateLimitWithContext(context.Background(), updateRateLimitOptions)
}

// UpdateRateLimitWithContext is an alternate form of the UpdateRateLimit method which supports a Context parameter
func (zoneRateLimits *ZoneRateLimitsV1) UpdateRateLimitWithContext(ctx context.Context, updateRateLimitOptions *UpdateRateLimitOptions) (result *RatelimitResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRateLimitOptions, "updateRateLimitOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRateLimitOptions, "updateRateLimitOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneRateLimits.Crn,
		"zone_identifier": *zoneRateLimits.ZoneIdentifier,
		"rate_limit_identifier": *updateRateLimitOptions.RateLimitIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneRateLimits.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneRateLimits.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/rate_limits/{rate_limit_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRateLimitOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_rate_limits", "V1", "UpdateRateLimit")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateRateLimitOptions.Disabled != nil {
		body["disabled"] = updateRateLimitOptions.Disabled
	}
	if updateRateLimitOptions.Description != nil {
		body["description"] = updateRateLimitOptions.Description
	}
	if updateRateLimitOptions.Bypass != nil {
		body["bypass"] = updateRateLimitOptions.Bypass
	}
	if updateRateLimitOptions.Threshold != nil {
		body["threshold"] = updateRateLimitOptions.Threshold
	}
	if updateRateLimitOptions.Period != nil {
		body["period"] = updateRateLimitOptions.Period
	}
	if updateRateLimitOptions.Action != nil {
		body["action"] = updateRateLimitOptions.Action
	}
	if updateRateLimitOptions.Correlate != nil {
		body["correlate"] = updateRateLimitOptions.Correlate
	}
	if updateRateLimitOptions.Match != nil {
		body["match"] = updateRateLimitOptions.Match
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
	response, err = zoneRateLimits.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRatelimitResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneRateLimitsOptions : The CreateZoneRateLimits options.
type CreateZoneRateLimitsOptions struct {
	// Whether this ratelimit is currently disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// A note that you can use to describe the reason for a rate limit.
	Description *string `json:"description,omitempty"`

	// Criteria that would allow the rate limit to be bypassed, for example to express that you shouldn't apply a rate
	// limit to a given set of URLs.
	Bypass []RatelimitInputBypassItem `json:"bypass,omitempty"`

	// The threshold that triggers the rate limit mitigations, combine with period. i.e. threshold per period.
	Threshold *int64 `json:"threshold,omitempty"`

	// The time in seconds to count matching traffic. If the count exceeds threshold within this period the action will be
	// performed.
	Period *int64 `json:"period,omitempty"`

	// action.
	Action *RatelimitInputAction `json:"action,omitempty"`

	// Enable NAT based rate limits.
	Correlate *RatelimitInputCorrelate `json:"correlate,omitempty"`

	// Determines which traffic the rate limit counts towards the threshold. Needs to be one of "request" or "response"
	// objects.
	Match *RatelimitInputMatch `json:"match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateZoneRateLimitsOptions : Instantiate CreateZoneRateLimitsOptions
func (*ZoneRateLimitsV1) NewCreateZoneRateLimitsOptions() *CreateZoneRateLimitsOptions {
	return &CreateZoneRateLimitsOptions{}
}

// SetDisabled : Allow user to set Disabled
func (options *CreateZoneRateLimitsOptions) SetDisabled(disabled bool) *CreateZoneRateLimitsOptions {
	options.Disabled = core.BoolPtr(disabled)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateZoneRateLimitsOptions) SetDescription(description string) *CreateZoneRateLimitsOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetBypass : Allow user to set Bypass
func (options *CreateZoneRateLimitsOptions) SetBypass(bypass []RatelimitInputBypassItem) *CreateZoneRateLimitsOptions {
	options.Bypass = bypass
	return options
}

// SetThreshold : Allow user to set Threshold
func (options *CreateZoneRateLimitsOptions) SetThreshold(threshold int64) *CreateZoneRateLimitsOptions {
	options.Threshold = core.Int64Ptr(threshold)
	return options
}

// SetPeriod : Allow user to set Period
func (options *CreateZoneRateLimitsOptions) SetPeriod(period int64) *CreateZoneRateLimitsOptions {
	options.Period = core.Int64Ptr(period)
	return options
}

// SetAction : Allow user to set Action
func (options *CreateZoneRateLimitsOptions) SetAction(action *RatelimitInputAction) *CreateZoneRateLimitsOptions {
	options.Action = action
	return options
}

// SetCorrelate : Allow user to set Correlate
func (options *CreateZoneRateLimitsOptions) SetCorrelate(correlate *RatelimitInputCorrelate) *CreateZoneRateLimitsOptions {
	options.Correlate = correlate
	return options
}

// SetMatch : Allow user to set Match
func (options *CreateZoneRateLimitsOptions) SetMatch(match *RatelimitInputMatch) *CreateZoneRateLimitsOptions {
	options.Match = match
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateZoneRateLimitsOptions) SetHeaders(param map[string]string) *CreateZoneRateLimitsOptions {
	options.Headers = param
	return options
}

// DeleteRateLimitRespResult : Container for response information.
type DeleteRateLimitRespResult struct {
	// ID.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteRateLimitRespResult unmarshals an instance of DeleteRateLimitRespResult from the specified map of raw messages.
func UnmarshalDeleteRateLimitRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteRateLimitRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteZoneRateLimitOptions : The DeleteZoneRateLimit options.
type DeleteZoneRateLimitOptions struct {
	// Identifier of the rate limit to be deleted.
	RateLimitIdentifier *string `json:"rate_limit_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneRateLimitOptions : Instantiate DeleteZoneRateLimitOptions
func (*ZoneRateLimitsV1) NewDeleteZoneRateLimitOptions(rateLimitIdentifier string) *DeleteZoneRateLimitOptions {
	return &DeleteZoneRateLimitOptions{
		RateLimitIdentifier: core.StringPtr(rateLimitIdentifier),
	}
}

// SetRateLimitIdentifier : Allow user to set RateLimitIdentifier
func (options *DeleteZoneRateLimitOptions) SetRateLimitIdentifier(rateLimitIdentifier string) *DeleteZoneRateLimitOptions {
	options.RateLimitIdentifier = core.StringPtr(rateLimitIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneRateLimitOptions) SetHeaders(param map[string]string) *DeleteZoneRateLimitOptions {
	options.Headers = param
	return options
}

// GetRateLimitOptions : The GetRateLimit options.
type GetRateLimitOptions struct {
	// Identifier of rate limit for the given zone.
	RateLimitIdentifier *string `json:"rate_limit_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRateLimitOptions : Instantiate GetRateLimitOptions
func (*ZoneRateLimitsV1) NewGetRateLimitOptions(rateLimitIdentifier string) *GetRateLimitOptions {
	return &GetRateLimitOptions{
		RateLimitIdentifier: core.StringPtr(rateLimitIdentifier),
	}
}

// SetRateLimitIdentifier : Allow user to set RateLimitIdentifier
func (options *GetRateLimitOptions) SetRateLimitIdentifier(rateLimitIdentifier string) *GetRateLimitOptions {
	options.RateLimitIdentifier = core.StringPtr(rateLimitIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetRateLimitOptions) SetHeaders(param map[string]string) *GetRateLimitOptions {
	options.Headers = param
	return options
}

// ListAllZoneRateLimitsOptions : The ListAllZoneRateLimits options.
type ListAllZoneRateLimitsOptions struct {
	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Maximum number of rate limits per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllZoneRateLimitsOptions : Instantiate ListAllZoneRateLimitsOptions
func (*ZoneRateLimitsV1) NewListAllZoneRateLimitsOptions() *ListAllZoneRateLimitsOptions {
	return &ListAllZoneRateLimitsOptions{}
}

// SetPage : Allow user to set Page
func (options *ListAllZoneRateLimitsOptions) SetPage(page int64) *ListAllZoneRateLimitsOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListAllZoneRateLimitsOptions) SetPerPage(perPage int64) *ListAllZoneRateLimitsOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllZoneRateLimitsOptions) SetHeaders(param map[string]string) *ListAllZoneRateLimitsOptions {
	options.Headers = param
	return options
}

// ListRatelimitRespResultInfo : Statistics of results.
type ListRatelimitRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalListRatelimitRespResultInfo unmarshals an instance of ListRatelimitRespResultInfo from the specified map of raw messages.
func UnmarshalListRatelimitRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListRatelimitRespResultInfo)
	err = core.UnmarshalPrimitive(m, "page", &obj.Page)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "per_page", &obj.PerPage)
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

// RatelimitInputAction : action.
type RatelimitInputAction struct {
	// The type of action to perform.
	Mode *string `json:"mode" validate:"required"`

	// The time in seconds as an integer to perform the mitigation action. Must be the same or greater than the period.
	// This field is valid only when mode is "simulate" or "ban".
	Timeout *int64 `json:"timeout,omitempty"`

	// Custom content-type and body to return, this overrides the custom error for the zone. This field is not required.
	// Omission will result in default HTML error page.This field is valid only when mode is "simulate" or "ban".
	Response *RatelimitInputActionResponse `json:"response,omitempty"`
}

// Constants associated with the RatelimitInputAction.Mode property.
// The type of action to perform.
const (
	RatelimitInputAction_Mode_Ban = "ban"
	RatelimitInputAction_Mode_Challenge = "challenge"
	RatelimitInputAction_Mode_JsChallenge = "js_challenge"
	RatelimitInputAction_Mode_Simulate = "simulate"
)


// NewRatelimitInputAction : Instantiate RatelimitInputAction (Generic Model Constructor)
func (*ZoneRateLimitsV1) NewRatelimitInputAction(mode string) (model *RatelimitInputAction, err error) {
	model = &RatelimitInputAction{
		Mode: core.StringPtr(mode),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRatelimitInputAction unmarshals an instance of RatelimitInputAction from the specified map of raw messages.
func UnmarshalRatelimitInputAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputAction)
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response", &obj.Response, UnmarshalRatelimitInputActionResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputActionResponse : Custom content-type and body to return, this overrides the custom error for the zone. This field is not required.
// Omission will result in default HTML error page.This field is valid only when mode is "simulate" or "ban".
type RatelimitInputActionResponse struct {
	// The content type of the body.
	ContentType *string `json:"content_type,omitempty"`

	// The body to return, the content here should conform to the content_type.
	Body *string `json:"body,omitempty"`
}

// Constants associated with the RatelimitInputActionResponse.ContentType property.
// The content type of the body.
const (
	RatelimitInputActionResponse_ContentType_ApplicationJSON = "application/json"
	RatelimitInputActionResponse_ContentType_TextPlain = "text/plain"
	RatelimitInputActionResponse_ContentType_TextXml = "text/xml"
)


// UnmarshalRatelimitInputActionResponse unmarshals an instance of RatelimitInputActionResponse from the specified map of raw messages.
func UnmarshalRatelimitInputActionResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputActionResponse)
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "body", &obj.Body)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputBypassItem : RatelimitInputBypassItem struct
type RatelimitInputBypassItem struct {
	// Rate limit name.
	Name *string `json:"name" validate:"required"`

	// The url to bypass.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the RatelimitInputBypassItem.Name property.
// Rate limit name.
const (
	RatelimitInputBypassItem_Name_URL = "url"
)


// NewRatelimitInputBypassItem : Instantiate RatelimitInputBypassItem (Generic Model Constructor)
func (*ZoneRateLimitsV1) NewRatelimitInputBypassItem(name string, value string) (model *RatelimitInputBypassItem, err error) {
	model = &RatelimitInputBypassItem{
		Name: core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRatelimitInputBypassItem unmarshals an instance of RatelimitInputBypassItem from the specified map of raw messages.
func UnmarshalRatelimitInputBypassItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputBypassItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputCorrelate : Enable NAT based rate limits.
type RatelimitInputCorrelate struct {
	// NAT rate limits by.
	By *string `json:"by" validate:"required"`
}

// Constants associated with the RatelimitInputCorrelate.By property.
// NAT rate limits by.
const (
	RatelimitInputCorrelate_By_Nat = "nat"
)


// NewRatelimitInputCorrelate : Instantiate RatelimitInputCorrelate (Generic Model Constructor)
func (*ZoneRateLimitsV1) NewRatelimitInputCorrelate(by string) (model *RatelimitInputCorrelate, err error) {
	model = &RatelimitInputCorrelate{
		By: core.StringPtr(by),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRatelimitInputCorrelate unmarshals an instance of RatelimitInputCorrelate from the specified map of raw messages.
func UnmarshalRatelimitInputCorrelate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputCorrelate)
	err = core.UnmarshalPrimitive(m, "by", &obj.By)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputMatch : Determines which traffic the rate limit counts towards the threshold. Needs to be one of "request" or "response"
// objects.
type RatelimitInputMatch struct {
	// request.
	Request *RatelimitInputMatchRequest `json:"request,omitempty"`

	// response.
	Response *RatelimitInputMatchResponse `json:"response,omitempty"`
}


// UnmarshalRatelimitInputMatch unmarshals an instance of RatelimitInputMatch from the specified map of raw messages.
func UnmarshalRatelimitInputMatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputMatch)
	err = core.UnmarshalModel(m, "request", &obj.Request, UnmarshalRatelimitInputMatchRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response", &obj.Response, UnmarshalRatelimitInputMatchResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputMatchRequest : request.
type RatelimitInputMatchRequest struct {
	// A subset of the list HTTP methods, or ["_ALL_"] for selecting all methods.
	Methods []string `json:"methods,omitempty"`

	// HTTP schemes list, or ["_ALL_"] for selecting all schemes.
	Schemes []string `json:"schemes,omitempty"`

	// The URL pattern to match comprised of the host and path, i.e. example.org/path. Wildcard are expanded to match
	// applicable traffic, query strings are not matched. Use * for all traffic to your zone.
	URL *string `json:"url" validate:"required"`
}

// Constants associated with the RatelimitInputMatchRequest.Methods property.
const (
	RatelimitInputMatchRequest_Methods_All = "_ALL_"
	RatelimitInputMatchRequest_Methods_Delete = "DELETE"
	RatelimitInputMatchRequest_Methods_Get = "GET"
	RatelimitInputMatchRequest_Methods_Head = "HEAD"
	RatelimitInputMatchRequest_Methods_Patch = "PATCH"
	RatelimitInputMatchRequest_Methods_Post = "POST"
	RatelimitInputMatchRequest_Methods_Put = "PUT"
)

// Constants associated with the RatelimitInputMatchRequest.Schemes property.
const (
	RatelimitInputMatchRequest_Schemes_All = "_ALL_"
	RatelimitInputMatchRequest_Schemes_Http = "HTTP"
	RatelimitInputMatchRequest_Schemes_Https = "HTTPS"
)


// NewRatelimitInputMatchRequest : Instantiate RatelimitInputMatchRequest (Generic Model Constructor)
func (*ZoneRateLimitsV1) NewRatelimitInputMatchRequest(url string) (model *RatelimitInputMatchRequest, err error) {
	model = &RatelimitInputMatchRequest{
		URL: core.StringPtr(url),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRatelimitInputMatchRequest unmarshals an instance of RatelimitInputMatchRequest from the specified map of raw messages.
func UnmarshalRatelimitInputMatchRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputMatchRequest)
	err = core.UnmarshalPrimitive(m, "methods", &obj.Methods)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schemes", &obj.Schemes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputMatchResponse : response.
type RatelimitInputMatchResponse struct {
	// HTTP Status codes, can be one [403], many [401,403] or indicate all by not providing this value. This field is not
	// required.
	Status []int64 `json:"status,omitempty"`

	// Array of response headers to match. If a response does not meet the header criteria then the request will not be
	// counted towards the rate limit.
	HeadersVar []RatelimitInputMatchResponseHeadersItem `json:"headers,omitempty"`

	// Deprecated, please use response headers instead and also provide "origin_traffic:false" to avoid legacy behaviour
	// interacting with the response.headers property.
	OriginTraffic *bool `json:"origin_traffic,omitempty"`
}


// UnmarshalRatelimitInputMatchResponse unmarshals an instance of RatelimitInputMatchResponse from the specified map of raw messages.
func UnmarshalRatelimitInputMatchResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputMatchResponse)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "headers", &obj.HeadersVar, UnmarshalRatelimitInputMatchResponseHeadersItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "origin_traffic", &obj.OriginTraffic)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitInputMatchResponseHeadersItem : RatelimitInputMatchResponseHeadersItem struct
type RatelimitInputMatchResponseHeadersItem struct {
	// The name of the response header to match.
	Name *string `json:"name" validate:"required"`

	// The operator when matchin, eq means equals, ne means not equals.
	Op *string `json:"op" validate:"required"`

	// The value of the header, which will be exactly matched.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the RatelimitInputMatchResponseHeadersItem.Op property.
// The operator when matchin, eq means equals, ne means not equals.
const (
	RatelimitInputMatchResponseHeadersItem_Op_Eq = "eq"
	RatelimitInputMatchResponseHeadersItem_Op_Ne = "ne"
)

// Constants associated with the RatelimitInputMatchResponseHeadersItem.Value property.
// The value of the header, which will be exactly matched.
const (
	RatelimitInputMatchResponseHeadersItem_Value_Hit = "HIT"
)


// NewRatelimitInputMatchResponseHeadersItem : Instantiate RatelimitInputMatchResponseHeadersItem (Generic Model Constructor)
func (*ZoneRateLimitsV1) NewRatelimitInputMatchResponseHeadersItem(name string, op string, value string) (model *RatelimitInputMatchResponseHeadersItem, err error) {
	model = &RatelimitInputMatchResponseHeadersItem{
		Name: core.StringPtr(name),
		Op: core.StringPtr(op),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRatelimitInputMatchResponseHeadersItem unmarshals an instance of RatelimitInputMatchResponseHeadersItem from the specified map of raw messages.
func UnmarshalRatelimitInputMatchResponseHeadersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitInputMatchResponseHeadersItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "op", &obj.Op)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectAction : action.
type RatelimitObjectAction struct {
	// The type of action to perform.
	Mode *string `json:"mode" validate:"required"`

	// The time in seconds as an integer to perform the mitigation action. Must be the same or greater than the period.
	// This field is valid only when mode is "simulate" or "ban".
	Timeout *int64 `json:"timeout,omitempty"`

	// Custom content-type and body to return, this overrides the custom error for the zone. This field is not required.
	// Omission will result in default HTML error page.This field is valid only when mode is "simulate" or "ban".
	Response *RatelimitObjectActionResponse `json:"response,omitempty"`
}

// Constants associated with the RatelimitObjectAction.Mode property.
// The type of action to perform.
const (
	RatelimitObjectAction_Mode_Ban = "ban"
	RatelimitObjectAction_Mode_Challenge = "challenge"
	RatelimitObjectAction_Mode_JsChallenge = "js_challenge"
	RatelimitObjectAction_Mode_Simulate = "simulate"
)


// UnmarshalRatelimitObjectAction unmarshals an instance of RatelimitObjectAction from the specified map of raw messages.
func UnmarshalRatelimitObjectAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectAction)
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response", &obj.Response, UnmarshalRatelimitObjectActionResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectActionResponse : Custom content-type and body to return, this overrides the custom error for the zone. This field is not required.
// Omission will result in default HTML error page.This field is valid only when mode is "simulate" or "ban".
type RatelimitObjectActionResponse struct {
	// The content type of the body.
	ContentType *string `json:"content_type" validate:"required"`

	// The body to return, the content here should conform to the content_type.
	Body *string `json:"body" validate:"required"`
}

// Constants associated with the RatelimitObjectActionResponse.ContentType property.
// The content type of the body.
const (
	RatelimitObjectActionResponse_ContentType_ApplicationJSON = "application/json"
	RatelimitObjectActionResponse_ContentType_TextPlain = "text/plain"
	RatelimitObjectActionResponse_ContentType_TextXml = "text/xml"
)


// UnmarshalRatelimitObjectActionResponse unmarshals an instance of RatelimitObjectActionResponse from the specified map of raw messages.
func UnmarshalRatelimitObjectActionResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectActionResponse)
	err = core.UnmarshalPrimitive(m, "content_type", &obj.ContentType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "body", &obj.Body)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectBypassItem : RatelimitObjectBypassItem struct
type RatelimitObjectBypassItem struct {
	// rate limit name.
	Name *string `json:"name" validate:"required"`

	// The url to bypass.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the RatelimitObjectBypassItem.Name property.
// rate limit name.
const (
	RatelimitObjectBypassItem_Name_URL = "url"
)


// UnmarshalRatelimitObjectBypassItem unmarshals an instance of RatelimitObjectBypassItem from the specified map of raw messages.
func UnmarshalRatelimitObjectBypassItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectBypassItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectCorrelate : Enable NAT based rate limits.
type RatelimitObjectCorrelate struct {
	// rate limit enabled by.
	By *string `json:"by" validate:"required"`
}

// Constants associated with the RatelimitObjectCorrelate.By property.
// rate limit enabled by.
const (
	RatelimitObjectCorrelate_By_Nat = "nat"
)


// UnmarshalRatelimitObjectCorrelate unmarshals an instance of RatelimitObjectCorrelate from the specified map of raw messages.
func UnmarshalRatelimitObjectCorrelate(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectCorrelate)
	err = core.UnmarshalPrimitive(m, "by", &obj.By)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectMatch : Determines which traffic the rate limit counts towards the threshold. Needs to be one of "request" or "response"
// objects.
type RatelimitObjectMatch struct {
	// request.
	Request *RatelimitObjectMatchRequest `json:"request,omitempty"`

	// response.
	Response *RatelimitObjectMatchResponse `json:"response,omitempty"`
}


// UnmarshalRatelimitObjectMatch unmarshals an instance of RatelimitObjectMatch from the specified map of raw messages.
func UnmarshalRatelimitObjectMatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectMatch)
	err = core.UnmarshalModel(m, "request", &obj.Request, UnmarshalRatelimitObjectMatchRequest)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "response", &obj.Response, UnmarshalRatelimitObjectMatchResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectMatchRequest : request.
type RatelimitObjectMatchRequest struct {
	// A subset of the list HTTP methods, or ["_ALL_"] for selecting all methods.
	Methods []string `json:"methods,omitempty"`

	// HTTP schemes list, or ["_ALL_"] for selecting all schemes.
	Schemes []string `json:"schemes,omitempty"`

	// The URL pattern to match comprised of the host and path, i.e. example.org/path. Wildcard are expanded to match
	// applicable traffic, query strings are not matched. Use * for all traffic to your zone.
	URL *string `json:"url" validate:"required"`
}

// Constants associated with the RatelimitObjectMatchRequest.Methods property.
const (
	RatelimitObjectMatchRequest_Methods_All = "_ALL_"
	RatelimitObjectMatchRequest_Methods_Delete = "DELETE"
	RatelimitObjectMatchRequest_Methods_Get = "GET"
	RatelimitObjectMatchRequest_Methods_Head = "HEAD"
	RatelimitObjectMatchRequest_Methods_Patch = "PATCH"
	RatelimitObjectMatchRequest_Methods_Post = "POST"
	RatelimitObjectMatchRequest_Methods_Put = "PUT"
)

// Constants associated with the RatelimitObjectMatchRequest.Schemes property.
const (
	RatelimitObjectMatchRequest_Schemes_All = "_ALL_"
	RatelimitObjectMatchRequest_Schemes_Http = "HTTP"
	RatelimitObjectMatchRequest_Schemes_Https = "HTTPS"
)


// UnmarshalRatelimitObjectMatchRequest unmarshals an instance of RatelimitObjectMatchRequest from the specified map of raw messages.
func UnmarshalRatelimitObjectMatchRequest(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectMatchRequest)
	err = core.UnmarshalPrimitive(m, "methods", &obj.Methods)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "schemes", &obj.Schemes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectMatchResponse : response.
type RatelimitObjectMatchResponse struct {
	// HTTP Status codes, can be one [403], many [401,403] or indicate all by not providing this value. This field is not
	// required.
	Status []int64 `json:"status,omitempty"`

	// Array of response headers to match. If a response does not meet the header criteria then the request will not be
	// counted towards the rate limit.
	HeadersVar []RatelimitObjectMatchResponseHeadersItem `json:"headers,omitempty"`

	// Deprecated, please use response headers instead and also provide "origin_traffic:false" to avoid legacy behaviour
	// interacting with the response.headers property.
	OriginTraffic *bool `json:"origin_traffic,omitempty"`
}


// UnmarshalRatelimitObjectMatchResponse unmarshals an instance of RatelimitObjectMatchResponse from the specified map of raw messages.
func UnmarshalRatelimitObjectMatchResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectMatchResponse)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "headers", &obj.HeadersVar, UnmarshalRatelimitObjectMatchResponseHeadersItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "origin_traffic", &obj.OriginTraffic)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObjectMatchResponseHeadersItem : RatelimitObjectMatchResponseHeadersItem struct
type RatelimitObjectMatchResponseHeadersItem struct {
	// The name of the response header to match.
	Name *string `json:"name" validate:"required"`

	// The operator when matchin, eq means equals, ne means not equals.
	Op *string `json:"op" validate:"required"`

	// The value of the header, which will be exactly matched.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the RatelimitObjectMatchResponseHeadersItem.Op property.
// The operator when matchin, eq means equals, ne means not equals.
const (
	RatelimitObjectMatchResponseHeadersItem_Op_Eq = "eq"
	RatelimitObjectMatchResponseHeadersItem_Op_Ne = "ne"
)

// Constants associated with the RatelimitObjectMatchResponseHeadersItem.Value property.
// The value of the header, which will be exactly matched.
const (
	RatelimitObjectMatchResponseHeadersItem_Value_Hit = "HIT"
)


// UnmarshalRatelimitObjectMatchResponseHeadersItem unmarshals an instance of RatelimitObjectMatchResponseHeadersItem from the specified map of raw messages.
func UnmarshalRatelimitObjectMatchResponseHeadersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObjectMatchResponseHeadersItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "op", &obj.Op)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateRateLimitOptions : The UpdateRateLimit options.
type UpdateRateLimitOptions struct {
	// Identifier of rate limit.
	RateLimitIdentifier *string `json:"rate_limit_identifier" validate:"required,ne="`

	// Whether this ratelimit is currently disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// A note that you can use to describe the reason for a rate limit.
	Description *string `json:"description,omitempty"`

	// Criteria that would allow the rate limit to be bypassed, for example to express that you shouldn't apply a rate
	// limit to a given set of URLs.
	Bypass []RatelimitInputBypassItem `json:"bypass,omitempty"`

	// The threshold that triggers the rate limit mitigations, combine with period. i.e. threshold per period.
	Threshold *int64 `json:"threshold,omitempty"`

	// The time in seconds to count matching traffic. If the count exceeds threshold within this period the action will be
	// performed.
	Period *int64 `json:"period,omitempty"`

	// action.
	Action *RatelimitInputAction `json:"action,omitempty"`

	// Enable NAT based rate limits.
	Correlate *RatelimitInputCorrelate `json:"correlate,omitempty"`

	// Determines which traffic the rate limit counts towards the threshold. Needs to be one of "request" or "response"
	// objects.
	Match *RatelimitInputMatch `json:"match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateRateLimitOptions : Instantiate UpdateRateLimitOptions
func (*ZoneRateLimitsV1) NewUpdateRateLimitOptions(rateLimitIdentifier string) *UpdateRateLimitOptions {
	return &UpdateRateLimitOptions{
		RateLimitIdentifier: core.StringPtr(rateLimitIdentifier),
	}
}

// SetRateLimitIdentifier : Allow user to set RateLimitIdentifier
func (options *UpdateRateLimitOptions) SetRateLimitIdentifier(rateLimitIdentifier string) *UpdateRateLimitOptions {
	options.RateLimitIdentifier = core.StringPtr(rateLimitIdentifier)
	return options
}

// SetDisabled : Allow user to set Disabled
func (options *UpdateRateLimitOptions) SetDisabled(disabled bool) *UpdateRateLimitOptions {
	options.Disabled = core.BoolPtr(disabled)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateRateLimitOptions) SetDescription(description string) *UpdateRateLimitOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetBypass : Allow user to set Bypass
func (options *UpdateRateLimitOptions) SetBypass(bypass []RatelimitInputBypassItem) *UpdateRateLimitOptions {
	options.Bypass = bypass
	return options
}

// SetThreshold : Allow user to set Threshold
func (options *UpdateRateLimitOptions) SetThreshold(threshold int64) *UpdateRateLimitOptions {
	options.Threshold = core.Int64Ptr(threshold)
	return options
}

// SetPeriod : Allow user to set Period
func (options *UpdateRateLimitOptions) SetPeriod(period int64) *UpdateRateLimitOptions {
	options.Period = core.Int64Ptr(period)
	return options
}

// SetAction : Allow user to set Action
func (options *UpdateRateLimitOptions) SetAction(action *RatelimitInputAction) *UpdateRateLimitOptions {
	options.Action = action
	return options
}

// SetCorrelate : Allow user to set Correlate
func (options *UpdateRateLimitOptions) SetCorrelate(correlate *RatelimitInputCorrelate) *UpdateRateLimitOptions {
	options.Correlate = correlate
	return options
}

// SetMatch : Allow user to set Match
func (options *UpdateRateLimitOptions) SetMatch(match *RatelimitInputMatch) *UpdateRateLimitOptions {
	options.Match = match
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRateLimitOptions) SetHeaders(param map[string]string) *UpdateRateLimitOptions {
	options.Headers = param
	return options
}

// DeleteRateLimitResp : rate limit delete response.
type DeleteRateLimitResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *DeleteRateLimitRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteRateLimitResp unmarshals an instance of DeleteRateLimitResp from the specified map of raw messages.
func UnmarshalDeleteRateLimitResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteRateLimitResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteRateLimitRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListRatelimitResp : rate limit list response.
type ListRatelimitResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []RatelimitObject `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListRatelimitRespResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListRatelimitResp unmarshals an instance of ListRatelimitResp from the specified map of raw messages.
func UnmarshalListRatelimitResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListRatelimitResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRatelimitObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListRatelimitRespResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitObject : rate limit object.
type RatelimitObject struct {
	// Identifier of the rate limit.
	ID *string `json:"id" validate:"required"`

	// Whether this ratelimit is currently disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// A note that you can use to describe the reason for a rate limit.
	Description *string `json:"description" validate:"required"`

	// Criteria that would allow the rate limit to be bypassed, for example to express that you shouldn't apply a rate
	// limit to a given set of URLs.
	Bypass []RatelimitObjectBypassItem `json:"bypass" validate:"required"`

	// The threshold that triggers the rate limit mitigations, combine with period. i.e. threshold per period.
	Threshold *int64 `json:"threshold" validate:"required"`

	// The time in seconds to count matching traffic. If the count exceeds threshold within this period the action will be
	// performed.
	Period *int64 `json:"period" validate:"required"`

	// Enable NAT based rate limits.
	Correlate *RatelimitObjectCorrelate `json:"correlate,omitempty"`

	// action.
	Action *RatelimitObjectAction `json:"action" validate:"required"`

	// Determines which traffic the rate limit counts towards the threshold. Needs to be one of "request" or "response"
	// objects.
	Match *RatelimitObjectMatch `json:"match" validate:"required"`
}


// UnmarshalRatelimitObject unmarshals an instance of RatelimitObject from the specified map of raw messages.
func UnmarshalRatelimitObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "bypass", &obj.Bypass, UnmarshalRatelimitObjectBypassItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "threshold", &obj.Threshold)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "period", &obj.Period)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "correlate", &obj.Correlate, UnmarshalRatelimitObjectCorrelate)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "action", &obj.Action, UnmarshalRatelimitObjectAction)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "match", &obj.Match, UnmarshalRatelimitObjectMatch)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RatelimitResp : rate limit response.
type RatelimitResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// rate limit object.
	Result *RatelimitObject `json:"result" validate:"required"`
}


// UnmarshalRatelimitResp unmarshals an instance of RatelimitResp from the specified map of raw messages.
func UnmarshalRatelimitResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RatelimitResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRatelimitObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
