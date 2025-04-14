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
 * IBM OpenAPI SDK Code Generator Version: 3.64.1-cee95189-20230124-211647
 */

// Package botmanagementv1 : Operations and models for the BotManagementV1 service
package botmanagementv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// BotManagementV1 : Bot Management
//
// API Version: 1.0.1
type BotManagementV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string

	// Identifier of zone.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "bot_management"

// BotManagementV1Options : Service options
type BotManagementV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`

	// Identifier of zone.
	ZoneIdentifier *string `validate:"required"`
}

// NewBotManagementV1UsingExternalConfig : constructs an instance of BotManagementV1 with passed in options and external configuration.
func NewBotManagementV1UsingExternalConfig(options *BotManagementV1Options) (botManagement *BotManagementV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	botManagement, err = NewBotManagementV1(options)
	if err != nil {
		return
	}

	err = botManagement.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = botManagement.Service.SetServiceURL(options.URL)
	}
	return
}

// NewBotManagementV1 : constructs an instance of BotManagementV1 with passed in options.
func NewBotManagementV1(options *BotManagementV1Options) (service *BotManagementV1, err error) {
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

	service = &BotManagementV1{
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

// Clone makes a copy of "botManagement" suitable for processing requests.
func (botManagement *BotManagementV1) Clone() *BotManagementV1 {
	if core.IsNil(botManagement) {
		return nil
	}
	clone := *botManagement
	clone.Service = botManagement.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (botManagement *BotManagementV1) SetServiceURL(url string) error {
	return botManagement.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (botManagement *BotManagementV1) GetServiceURL() string {
	return botManagement.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (botManagement *BotManagementV1) SetDefaultHeaders(headers http.Header) {
	botManagement.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (botManagement *BotManagementV1) SetEnableGzipCompression(enableGzip bool) {
	botManagement.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (botManagement *BotManagementV1) GetEnableGzipCompression() bool {
	return botManagement.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (botManagement *BotManagementV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	botManagement.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (botManagement *BotManagementV1) DisableRetries() {
	botManagement.Service.DisableRetries()
}

// GetBotManagement : Get Bot management setting
// Get Bot management setting for a given zone.
func (botManagement *BotManagementV1) GetBotManagement(getBotManagementOptions *GetBotManagementOptions) (result *BotMgtResp, response *core.DetailedResponse, err error) {
	return botManagement.GetBotManagementWithContext(context.Background(), getBotManagementOptions)
}

// GetBotManagementWithContext is an alternate form of the GetBotManagement method which supports a Context parameter
func (botManagement *BotManagementV1) GetBotManagementWithContext(ctx context.Context, getBotManagementOptions *GetBotManagementOptions) (result *BotMgtResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getBotManagementOptions, "getBotManagementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *botManagement.Crn,
		"zone_identifier": *botManagement.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = botManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(botManagement.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/bot_management`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getBotManagementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("bot_management", "V1", "GetBotManagement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = botManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBotMgtResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateBotManagement : Update Bot management setting
// Update Bot management setting for given zone.
func (botManagement *BotManagementV1) UpdateBotManagement(updateBotManagementOptions *UpdateBotManagementOptions) (result *BotMgtResp, response *core.DetailedResponse, err error) {
	return botManagement.UpdateBotManagementWithContext(context.Background(), updateBotManagementOptions)
}

// UpdateBotManagementWithContext is an alternate form of the UpdateBotManagement method which supports a Context parameter
func (botManagement *BotManagementV1) UpdateBotManagementWithContext(ctx context.Context, updateBotManagementOptions *UpdateBotManagementOptions) (result *BotMgtResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(updateBotManagementOptions, "updateBotManagementOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn":             *botManagement.Crn,
		"zone_identifier": *botManagement.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = botManagement.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(botManagement.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/bot_management`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateBotManagementOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("bot_management", "V1", "UpdateBotManagement")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateBotManagementOptions.FightMode != nil {
		body["fight_mode"] = updateBotManagementOptions.FightMode
	}
	if updateBotManagementOptions.SessionScore != nil {
		body["session_score"] = updateBotManagementOptions.SessionScore
	}
	if updateBotManagementOptions.EnableJs != nil {
		body["enable_js"] = updateBotManagementOptions.EnableJs
	}
	if updateBotManagementOptions.AuthIdLogging != nil {
		body["auth_id_logging"] = updateBotManagementOptions.AuthIdLogging
	}
	if updateBotManagementOptions.UseLatestModel != nil {
		body["use_latest_model"] = updateBotManagementOptions.UseLatestModel
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
	response, err = botManagement.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBotMgtResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// BotMgtRespResult : Container for response information.
type BotMgtRespResult struct {
	FightMode *bool `json:"fight_mode,omitempty"`

	SessionScore *bool `json:"session_score,omitempty"`

	EnableJs *bool `json:"enable_js,omitempty"`

	AuthIdLogging *bool `json:"auth_id_logging,omitempty"`

	UseLatestModel *bool `json:"use_latest_model,omitempty"`
}

// UnmarshalBotMgtRespResult unmarshals an instance of BotMgtRespResult from the specified map of raw messages.
func UnmarshalBotMgtRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotMgtRespResult)
	err = core.UnmarshalPrimitive(m, "fight_mode", &obj.FightMode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "session_score", &obj.SessionScore)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_js", &obj.EnableJs)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "auth_id_logging", &obj.AuthIdLogging)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "use_latest_model", &obj.UseLatestModel)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetBotManagementOptions : The GetBotManagement options.
type GetBotManagementOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetBotManagementOptions : Instantiate GetBotManagementOptions
func (*BotManagementV1) NewGetBotManagementOptions() *GetBotManagementOptions {
	return &GetBotManagementOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetBotManagementOptions) SetHeaders(param map[string]string) *GetBotManagementOptions {
	options.Headers = param
	return options
}

// UpdateBotManagementOptions : The UpdateBotManagement options.
type UpdateBotManagementOptions struct {
	FightMode *bool `json:"fight_mode,omitempty"`

	SessionScore *bool `json:"session_score,omitempty"`

	EnableJs *bool `json:"enable_js,omitempty"`

	AuthIdLogging *bool `json:"auth_id_logging,omitempty"`

	UseLatestModel *bool `json:"use_latest_model,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateBotManagementOptions : Instantiate UpdateBotManagementOptions
func (*BotManagementV1) NewUpdateBotManagementOptions() *UpdateBotManagementOptions {
	return &UpdateBotManagementOptions{}
}

// SetFightMode : Allow user to set FightMode
func (_options *UpdateBotManagementOptions) SetFightMode(fightMode bool) *UpdateBotManagementOptions {
	_options.FightMode = core.BoolPtr(fightMode)
	return _options
}

// SetSessionScore : Allow user to set SessionScore
func (_options *UpdateBotManagementOptions) SetSessionScore(sessionScore bool) *UpdateBotManagementOptions {
	_options.SessionScore = core.BoolPtr(sessionScore)
	return _options
}

// SetEnableJs : Allow user to set EnableJs
func (_options *UpdateBotManagementOptions) SetEnableJs(enableJs bool) *UpdateBotManagementOptions {
	_options.EnableJs = core.BoolPtr(enableJs)
	return _options
}

// SetAuthIdLogging : Allow user to set AuthIdLogging
func (_options *UpdateBotManagementOptions) SetAuthIdLogging(authIdLogging bool) *UpdateBotManagementOptions {
	_options.AuthIdLogging = core.BoolPtr(authIdLogging)
	return _options
}

// SetUseLatestModel : Allow user to set UseLatestModel
func (_options *UpdateBotManagementOptions) SetUseLatestModel(useLatestModel bool) *UpdateBotManagementOptions {
	_options.UseLatestModel = core.BoolPtr(useLatestModel)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBotManagementOptions) SetHeaders(param map[string]string) *UpdateBotManagementOptions {
	options.Headers = param
	return options
}

// BotMgtResp : Bot Management Response.
type BotMgtResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *BotMgtRespResult `json:"result" validate:"required"`
}

// UnmarshalBotMgtResp unmarshals an instance of BotMgtResp from the specified map of raw messages.
func UnmarshalBotMgtResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BotMgtResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalBotMgtRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
