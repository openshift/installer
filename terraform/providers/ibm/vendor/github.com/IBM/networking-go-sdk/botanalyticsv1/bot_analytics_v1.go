/**
 * (C) Copyright IBM Corp. 2023.
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
 * IBM OpenAPI SDK Code Generator Version: 3.69.0-370d6400-20230329-174648
 */

// Package botanalyticsv1 : Operations and models for the BotAnalyticsV1 service
package botanalyticsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// BotAnalyticsV1 : Bot Analytics
//
// API Version: 1.0.1
type BotAnalyticsV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string

	// Zone identifier to identifiy the zone.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "bot_analytics"

// BotAnalyticsV1Options : Service options
type BotAnalyticsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`

	// Zone identifier to identifiy the zone.
	ZoneIdentifier *string `validate:"required"`
}

// NewBotAnalyticsV1UsingExternalConfig : constructs an instance of BotAnalyticsV1 with passed in options and external configuration.
func NewBotAnalyticsV1UsingExternalConfig(options *BotAnalyticsV1Options) (botAnalytics *BotAnalyticsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	botAnalytics, err = NewBotAnalyticsV1(options)
	if err != nil {
		return
	}

	err = botAnalytics.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = botAnalytics.Service.SetServiceURL(options.URL)
	}
	return
}

// NewBotAnalyticsV1 : constructs an instance of BotAnalyticsV1 with passed in options.
func NewBotAnalyticsV1(options *BotAnalyticsV1Options) (service *BotAnalyticsV1, err error) {
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

	service = &BotAnalyticsV1{
		Service:        baseService,
		Crn:            options.Crn,
		ZoneIdentifier: options.ZoneIdentifier,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "botAnalytics" suitable for processing requests.
func (botAnalytics *BotAnalyticsV1) Clone() *BotAnalyticsV1 {
	if core.IsNil(botAnalytics) {
		return nil
	}
	clone := *botAnalytics
	clone.Service = botAnalytics.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (botAnalytics *BotAnalyticsV1) SetServiceURL(url string) error {
	return botAnalytics.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (botAnalytics *BotAnalyticsV1) GetServiceURL() string {
	return botAnalytics.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (botAnalytics *BotAnalyticsV1) SetDefaultHeaders(headers http.Header) {
	botAnalytics.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (botAnalytics *BotAnalyticsV1) SetEnableGzipCompression(enableGzip bool) {
	botAnalytics.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (botAnalytics *BotAnalyticsV1) GetEnableGzipCompression() bool {
	return botAnalytics.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (botAnalytics *BotAnalyticsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	botAnalytics.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (botAnalytics *BotAnalyticsV1) DisableRetries() {
	botAnalytics.Service.DisableRetries()
}

// GetBotScore : Get Bot Analytics score source
// Get Bot Analytics score source for a given zone. Use this to identify the most common detection engines used to score
// your traffic.
func (botAnalytics *BotAnalyticsV1) GetBotScore(getBotScoreOptions *GetBotScoreOptions) (result *BotScoreResp, response *core.DetailedResponse, err error) {
	return botAnalytics.GetBotScoreWithContext(context.Background(), getBotScoreOptions)
}

// GetBotScoreWithContext is an alternate form of the GetBotScore method which supports a Context parameter
func (botAnalytics *BotAnalyticsV1) GetBotScoreWithContext(ctx context.Context, getBotScoreOptions *GetBotScoreOptions) (result *BotScoreResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBotScoreOptions, "getBotScoreOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getBotScoreOptions, "getBotScoreOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *botAnalytics.Crn,
		"zone_identifier": *botAnalytics.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = botAnalytics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(botAnalytics.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/bot_analytics/score_source`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getBotScoreOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("bot_analytics", "V1", "GetBotScore")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	sinceVal := strings.Split(fmt.Sprint(*getBotScoreOptions.Since), ".")
	untilVal := strings.Split(fmt.Sprint(*getBotScoreOptions.Until), ".")

	builder.AddQuery("since", sinceVal[0]+"Z")
	builder.AddQuery("until", untilVal[0]+"Z")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = botAnalytics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBotScoreResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetBotTimeseries : Get Bot Analytics timeseries
// Get Bot Analytics timeseries for a given zone.
func (botAnalytics *BotAnalyticsV1) GetBotTimeseries(getBotTimeseriesOptions *GetBotTimeseriesOptions) (result *BotTimeseriesResp, response *core.DetailedResponse, err error) {
	return botAnalytics.GetBotTimeseriesWithContext(context.Background(), getBotTimeseriesOptions)
}

// GetBotTimeseriesWithContext is an alternate form of the GetBotTimeseries method which supports a Context parameter
func (botAnalytics *BotAnalyticsV1) GetBotTimeseriesWithContext(ctx context.Context, getBotTimeseriesOptions *GetBotTimeseriesOptions) (result *BotTimeseriesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBotTimeseriesOptions, "getBotTimeseriesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getBotTimeseriesOptions, "getBotTimeseriesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *botAnalytics.Crn,
		"zone_identifier": *botAnalytics.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = botAnalytics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(botAnalytics.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/bot_analytics/timeseries`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getBotTimeseriesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("bot_analytics", "V1", "GetBotTimeseries")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	sinceVal := strings.Split(fmt.Sprint(*getBotTimeseriesOptions.Since), ".")
	untilVal := strings.Split(fmt.Sprint(*getBotTimeseriesOptions.Until), ".")

	builder.AddQuery("since", sinceVal[0]+"Z")
	builder.AddQuery("until", untilVal[0]+"Z")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = botAnalytics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBotTimeseriesResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetBotTopns : Get Bot Analytics top attributes
// Get Bot Analytics top attributes for a given zone. Use this to view more detailed information on specific IP
// addresses and other characteristics.
func (botAnalytics *BotAnalyticsV1) GetBotTopns(getBotTopnsOptions *GetBotTopnsOptions) (result *BotTopnsResp, response *core.DetailedResponse, err error) {
	return botAnalytics.GetBotTopnsWithContext(context.Background(), getBotTopnsOptions)
}

// GetBotTopnsWithContext is an alternate form of the GetBotTopns method which supports a Context parameter
func (botAnalytics *BotAnalyticsV1) GetBotTopnsWithContext(ctx context.Context, getBotTopnsOptions *GetBotTopnsOptions) (result *BotTopnsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBotTopnsOptions, "getBotTopnsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getBotTopnsOptions, "getBotTopnsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *botAnalytics.Crn,
		"zone_identifier": *botAnalytics.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = botAnalytics.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(botAnalytics.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/bot_analytics/top_ns`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getBotTopnsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("bot_analytics", "V1", "GetBotTopns")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	sinceVal := strings.Split(fmt.Sprint(*getBotTopnsOptions.Since), ".")
	untilVal := strings.Split(fmt.Sprint(*getBotTopnsOptions.Until), ".")

	builder.AddQuery("since", sinceVal[0]+"Z")
	builder.AddQuery("until", untilVal[0]+"Z")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = botAnalytics.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBotTopnsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// BotScoreRespResultItem : BotScoreRespResultItem struct
type BotScoreRespResultItem struct {
	BotScore []BotScoreRespResultItemBotScoreItem `json:"botScore,omitempty"`
}

// UnmarshalBotScoreRespResultItem unmarshals an instance of BotScoreRespResultItem from the specified map of raw messages.
func UnmarshalBotScoreRespResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotScoreRespResultItem)
	err = core.UnmarshalModel(m, "botScore", &obj.BotScore, UnmarshalBotScoreRespResultItemBotScoreItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BotScoreRespResultItemBotScoreItem : BotScoreRespResultItemBotScoreItem struct
type BotScoreRespResultItemBotScoreItem struct {
	Avg *BotScoreRespResultItemBotScoreItemAvg `json:"avg,omitempty"`

	Count *float64 `json:"count,omitempty"`

	Dimensions *BotScoreRespResultItemBotScoreItemDimensions `json:"dimensions,omitempty"`
}

// UnmarshalBotScoreRespResultItemBotScoreItem unmarshals an instance of BotScoreRespResultItemBotScoreItem from the specified map of raw messages.
func UnmarshalBotScoreRespResultItemBotScoreItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotScoreRespResultItemBotScoreItem)
	err = core.UnmarshalModel(m, "avg", &obj.Avg, UnmarshalBotScoreRespResultItemBotScoreItemAvg)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "dimensions", &obj.Dimensions, UnmarshalBotScoreRespResultItemBotScoreItemDimensions)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BotScoreRespResultItemBotScoreItemAvg : BotScoreRespResultItemBotScoreItemAvg struct
type BotScoreRespResultItemBotScoreItemAvg struct {
	SampleInterval *float64 `json:"sampleInterval,omitempty"`
}

// UnmarshalBotScoreRespResultItemBotScoreItemAvg unmarshals an instance of BotScoreRespResultItemBotScoreItemAvg from the specified map of raw messages.
func UnmarshalBotScoreRespResultItemBotScoreItemAvg(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotScoreRespResultItemBotScoreItemAvg)
	err = core.UnmarshalPrimitive(m, "sampleInterval", &obj.SampleInterval)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BotScoreRespResultItemBotScoreItemDimensions : BotScoreRespResultItemBotScoreItemDimensions struct
type BotScoreRespResultItemBotScoreItemDimensions struct {
	BotScoreSrcName *string `json:"botScoreSrcName,omitempty"`
}

// UnmarshalBotScoreRespResultItemBotScoreItemDimensions unmarshals an instance of BotScoreRespResultItemBotScoreItemDimensions from the specified map of raw messages.
func UnmarshalBotScoreRespResultItemBotScoreItemDimensions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotScoreRespResultItemBotScoreItemDimensions)
	err = core.UnmarshalPrimitive(m, "botScoreSrcName", &obj.BotScoreSrcName)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BotTimeseriesRespResultItem : BotTimeseriesRespResultItem struct
type BotTimeseriesRespResultItem struct {
	BotScore []map[string]interface{} `json:"botScore,omitempty"`
}

// UnmarshalBotTimeseriesRespResultItem unmarshals an instance of BotTimeseriesRespResultItem from the specified map of raw messages.
func UnmarshalBotTimeseriesRespResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotTimeseriesRespResultItem)
	err = core.UnmarshalPrimitive(m, "botScore", &obj.BotScore)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetBotScoreOptions : The GetBotScore options.
type GetBotScoreOptions struct {
	// UTC datetime for start of query.
	Since *strfmt.DateTime `json:"since" validate:"required"`

	// UTC datetime for end of query.
	Until *strfmt.DateTime `json:"until" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetBotScoreOptions : Instantiate GetBotScoreOptions
func (*BotAnalyticsV1) NewGetBotScoreOptions(since *strfmt.DateTime, until *strfmt.DateTime) *GetBotScoreOptions {
	return &GetBotScoreOptions{
		Since: since,
		Until: until,
	}
}

// SetSince : Allow user to set Since
func (_options *GetBotScoreOptions) SetSince(since *strfmt.DateTime) *GetBotScoreOptions {
	_options.Since = since
	return _options
}

// SetUntil : Allow user to set Until
func (_options *GetBotScoreOptions) SetUntil(until *strfmt.DateTime) *GetBotScoreOptions {
	_options.Until = until
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBotScoreOptions) SetHeaders(param map[string]string) *GetBotScoreOptions {
	options.Headers = param
	return options
}

// GetBotTimeseriesOptions : The GetBotTimeseries options.
type GetBotTimeseriesOptions struct {
	// UTC datetime for start of query.
	Since *strfmt.DateTime `json:"since" validate:"required"`

	// UTC datetime for end of query.
	Until *strfmt.DateTime `json:"until" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetBotTimeseriesOptions : Instantiate GetBotTimeseriesOptions
func (*BotAnalyticsV1) NewGetBotTimeseriesOptions(since *strfmt.DateTime, until *strfmt.DateTime) *GetBotTimeseriesOptions {
	return &GetBotTimeseriesOptions{
		Since: since,
		Until: until,
	}
}

// SetSince : Allow user to set Since
func (_options *GetBotTimeseriesOptions) SetSince(since *strfmt.DateTime) *GetBotTimeseriesOptions {
	_options.Since = since
	return _options
}

// SetUntil : Allow user to set Until
func (_options *GetBotTimeseriesOptions) SetUntil(until *strfmt.DateTime) *GetBotTimeseriesOptions {
	_options.Until = until
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBotTimeseriesOptions) SetHeaders(param map[string]string) *GetBotTimeseriesOptions {
	options.Headers = param
	return options
}

// GetBotTopnsOptions : The GetBotTopns options.
type GetBotTopnsOptions struct {
	// UTC datetime for start of query.
	Since *strfmt.DateTime `json:"since" validate:"required"`

	// UTC datetime for end of query.
	Until *strfmt.DateTime `json:"until" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetBotTopnsOptions : Instantiate GetBotTopnsOptions
func (*BotAnalyticsV1) NewGetBotTopnsOptions(since *strfmt.DateTime, until *strfmt.DateTime) *GetBotTopnsOptions {
	return &GetBotTopnsOptions{
		Since: since,
		Until: until,
	}
}

// SetSince : Allow user to set Since
func (_options *GetBotTopnsOptions) SetSince(since *strfmt.DateTime) *GetBotTopnsOptions {
	_options.Since = since
	return _options
}

// SetUntil : Allow user to set Until
func (_options *GetBotTopnsOptions) SetUntil(until *strfmt.DateTime) *GetBotTopnsOptions {
	_options.Until = until
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBotTopnsOptions) SetHeaders(param map[string]string) *GetBotTopnsOptions {
	options.Headers = param
	return options
}

// BotScoreResp : Bot Score Source Response.
type BotScoreResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []BotScoreRespResultItem `json:"result" validate:"required"`
}

// UnmarshalBotScoreResp unmarshals an instance of BotScoreResp from the specified map of raw messages.
func UnmarshalBotScoreResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotScoreResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalBotScoreRespResultItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BotTimeseriesResp : Bot Timeseries Response.
type BotTimeseriesResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []map[string]interface{} `json:"result" validate:"required"`
}

// UnmarshalBotTimeseriesResp unmarshals an instance of BotTimeseriesResp from the specified map of raw messages.
func UnmarshalBotTimeseriesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotTimeseriesResp)
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
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BotTopnsResp : Bot top attributes response.
type BotTopnsResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []map[string]interface{} `json:"result" validate:"required"`
}

// UnmarshalBotTopnsResp unmarshals an instance of BotTopnsResp from the specified map of raw messages.
func UnmarshalBotTopnsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotTopnsResp)
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
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
