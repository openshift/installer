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
 

// Package useragentblockingrulesv1 : Operations and models for the UserAgentBlockingRulesV1 service
package useragentblockingrulesv1

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

// UserAgentBlockingRulesV1 : User-Agent Blocking Rules
//
// Version: 1.0.1
type UserAgentBlockingRulesV1 struct {
	Service *core.BaseService

	// Full crn of the service instance.
	Crn *string

	// Zone identifier (zone id).
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "user_agent_blocking_rules"

// UserAgentBlockingRulesV1Options : Service options
type UserAgentBlockingRulesV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full crn of the service instance.
	Crn *string `validate:"required"`

	// Zone identifier (zone id).
	ZoneIdentifier *string `validate:"required"`
}

// NewUserAgentBlockingRulesV1UsingExternalConfig : constructs an instance of UserAgentBlockingRulesV1 with passed in options and external configuration.
func NewUserAgentBlockingRulesV1UsingExternalConfig(options *UserAgentBlockingRulesV1Options) (userAgentBlockingRules *UserAgentBlockingRulesV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	userAgentBlockingRules, err = NewUserAgentBlockingRulesV1(options)
	if err != nil {
		return
	}

	err = userAgentBlockingRules.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = userAgentBlockingRules.Service.SetServiceURL(options.URL)
	}
	return
}

// NewUserAgentBlockingRulesV1 : constructs an instance of UserAgentBlockingRulesV1 with passed in options.
func NewUserAgentBlockingRulesV1(options *UserAgentBlockingRulesV1Options) (service *UserAgentBlockingRulesV1, err error) {
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

	service = &UserAgentBlockingRulesV1{
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

// Clone makes a copy of "userAgentBlockingRules" suitable for processing requests.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) Clone() *UserAgentBlockingRulesV1 {
	if core.IsNil(userAgentBlockingRules) {
		return nil
	}
	clone := *userAgentBlockingRules
	clone.Service = userAgentBlockingRules.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (userAgentBlockingRules *UserAgentBlockingRulesV1) SetServiceURL(url string) error {
	return userAgentBlockingRules.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (userAgentBlockingRules *UserAgentBlockingRulesV1) GetServiceURL() string {
	return userAgentBlockingRules.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (userAgentBlockingRules *UserAgentBlockingRulesV1) SetDefaultHeaders(headers http.Header) {
	userAgentBlockingRules.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (userAgentBlockingRules *UserAgentBlockingRulesV1) SetEnableGzipCompression(enableGzip bool) {
	userAgentBlockingRules.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (userAgentBlockingRules *UserAgentBlockingRulesV1) GetEnableGzipCompression() bool {
	return userAgentBlockingRules.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	userAgentBlockingRules.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) DisableRetries() {
	userAgentBlockingRules.Service.DisableRetries()
}

// ListAllZoneUserAgentRules : List all user-agent blocking rules
// List all user agent blocking rules.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) ListAllZoneUserAgentRules(listAllZoneUserAgentRulesOptions *ListAllZoneUserAgentRulesOptions) (result *ListUseragentRulesResp, response *core.DetailedResponse, err error) {
	return userAgentBlockingRules.ListAllZoneUserAgentRulesWithContext(context.Background(), listAllZoneUserAgentRulesOptions)
}

// ListAllZoneUserAgentRulesWithContext is an alternate form of the ListAllZoneUserAgentRules method which supports a Context parameter
func (userAgentBlockingRules *UserAgentBlockingRulesV1) ListAllZoneUserAgentRulesWithContext(ctx context.Context, listAllZoneUserAgentRulesOptions *ListAllZoneUserAgentRulesOptions) (result *ListUseragentRulesResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllZoneUserAgentRulesOptions, "listAllZoneUserAgentRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *userAgentBlockingRules.Crn,
		"zone_identifier": *userAgentBlockingRules.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = userAgentBlockingRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(userAgentBlockingRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/ua_rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllZoneUserAgentRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("user_agent_blocking_rules", "V1", "ListAllZoneUserAgentRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllZoneUserAgentRulesOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listAllZoneUserAgentRulesOptions.Page))
	}
	if listAllZoneUserAgentRulesOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listAllZoneUserAgentRulesOptions.PerPage))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = userAgentBlockingRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListUseragentRulesResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneUserAgentRule : Create user-agent blocking rule
// Create a new user-agent blocking rule for a given zone under a service instance.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) CreateZoneUserAgentRule(createZoneUserAgentRuleOptions *CreateZoneUserAgentRuleOptions) (result *UseragentRuleResp, response *core.DetailedResponse, err error) {
	return userAgentBlockingRules.CreateZoneUserAgentRuleWithContext(context.Background(), createZoneUserAgentRuleOptions)
}

// CreateZoneUserAgentRuleWithContext is an alternate form of the CreateZoneUserAgentRule method which supports a Context parameter
func (userAgentBlockingRules *UserAgentBlockingRulesV1) CreateZoneUserAgentRuleWithContext(ctx context.Context, createZoneUserAgentRuleOptions *CreateZoneUserAgentRuleOptions) (result *UseragentRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createZoneUserAgentRuleOptions, "createZoneUserAgentRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *userAgentBlockingRules.Crn,
		"zone_identifier": *userAgentBlockingRules.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = userAgentBlockingRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(userAgentBlockingRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/ua_rules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createZoneUserAgentRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("user_agent_blocking_rules", "V1", "CreateZoneUserAgentRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createZoneUserAgentRuleOptions.Paused != nil {
		body["paused"] = createZoneUserAgentRuleOptions.Paused
	}
	if createZoneUserAgentRuleOptions.Description != nil {
		body["description"] = createZoneUserAgentRuleOptions.Description
	}
	if createZoneUserAgentRuleOptions.Mode != nil {
		body["mode"] = createZoneUserAgentRuleOptions.Mode
	}
	if createZoneUserAgentRuleOptions.Configuration != nil {
		body["configuration"] = createZoneUserAgentRuleOptions.Configuration
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
	response, err = userAgentBlockingRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUseragentRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteZoneUserAgentRule : Delete user-agent blocking rule
// Delete a user-agent blocking rule for a particular zone, given its id.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) DeleteZoneUserAgentRule(deleteZoneUserAgentRuleOptions *DeleteZoneUserAgentRuleOptions) (result *DeleteUseragentRuleResp, response *core.DetailedResponse, err error) {
	return userAgentBlockingRules.DeleteZoneUserAgentRuleWithContext(context.Background(), deleteZoneUserAgentRuleOptions)
}

// DeleteZoneUserAgentRuleWithContext is an alternate form of the DeleteZoneUserAgentRule method which supports a Context parameter
func (userAgentBlockingRules *UserAgentBlockingRulesV1) DeleteZoneUserAgentRuleWithContext(ctx context.Context, deleteZoneUserAgentRuleOptions *DeleteZoneUserAgentRuleOptions) (result *DeleteUseragentRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneUserAgentRuleOptions, "deleteZoneUserAgentRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneUserAgentRuleOptions, "deleteZoneUserAgentRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *userAgentBlockingRules.Crn,
		"zone_identifier": *userAgentBlockingRules.ZoneIdentifier,
		"useragent_rule_identifier": *deleteZoneUserAgentRuleOptions.UseragentRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = userAgentBlockingRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(userAgentBlockingRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/ua_rules/{useragent_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneUserAgentRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("user_agent_blocking_rules", "V1", "DeleteZoneUserAgentRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = userAgentBlockingRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteUseragentRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetUserAgentRule : Get user-agent blocking rule
// For a given service instance, zone id and user-agent rule id, get the user-agent blocking rule details.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) GetUserAgentRule(getUserAgentRuleOptions *GetUserAgentRuleOptions) (result *UseragentRuleResp, response *core.DetailedResponse, err error) {
	return userAgentBlockingRules.GetUserAgentRuleWithContext(context.Background(), getUserAgentRuleOptions)
}

// GetUserAgentRuleWithContext is an alternate form of the GetUserAgentRule method which supports a Context parameter
func (userAgentBlockingRules *UserAgentBlockingRulesV1) GetUserAgentRuleWithContext(ctx context.Context, getUserAgentRuleOptions *GetUserAgentRuleOptions) (result *UseragentRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getUserAgentRuleOptions, "getUserAgentRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getUserAgentRuleOptions, "getUserAgentRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *userAgentBlockingRules.Crn,
		"zone_identifier": *userAgentBlockingRules.ZoneIdentifier,
		"useragent_rule_identifier": *getUserAgentRuleOptions.UseragentRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = userAgentBlockingRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(userAgentBlockingRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/ua_rules/{useragent_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getUserAgentRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("user_agent_blocking_rules", "V1", "GetUserAgentRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = userAgentBlockingRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUseragentRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateUserAgentRule : Update user-agent blocking rule
// Update an existing user-agent blocking rule for a given zone under a given service instance.
func (userAgentBlockingRules *UserAgentBlockingRulesV1) UpdateUserAgentRule(updateUserAgentRuleOptions *UpdateUserAgentRuleOptions) (result *UseragentRuleResp, response *core.DetailedResponse, err error) {
	return userAgentBlockingRules.UpdateUserAgentRuleWithContext(context.Background(), updateUserAgentRuleOptions)
}

// UpdateUserAgentRuleWithContext is an alternate form of the UpdateUserAgentRule method which supports a Context parameter
func (userAgentBlockingRules *UserAgentBlockingRulesV1) UpdateUserAgentRuleWithContext(ctx context.Context, updateUserAgentRuleOptions *UpdateUserAgentRuleOptions) (result *UseragentRuleResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateUserAgentRuleOptions, "updateUserAgentRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateUserAgentRuleOptions, "updateUserAgentRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *userAgentBlockingRules.Crn,
		"zone_identifier": *userAgentBlockingRules.ZoneIdentifier,
		"useragent_rule_identifier": *updateUserAgentRuleOptions.UseragentRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = userAgentBlockingRules.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(userAgentBlockingRules.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/ua_rules/{useragent_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateUserAgentRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("user_agent_blocking_rules", "V1", "UpdateUserAgentRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateUserAgentRuleOptions.Paused != nil {
		body["paused"] = updateUserAgentRuleOptions.Paused
	}
	if updateUserAgentRuleOptions.Description != nil {
		body["description"] = updateUserAgentRuleOptions.Description
	}
	if updateUserAgentRuleOptions.Mode != nil {
		body["mode"] = updateUserAgentRuleOptions.Mode
	}
	if updateUserAgentRuleOptions.Configuration != nil {
		body["configuration"] = updateUserAgentRuleOptions.Configuration
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
	response, err = userAgentBlockingRules.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUseragentRuleResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneUserAgentRuleOptions : The CreateZoneUserAgentRule options.
type CreateZoneUserAgentRuleOptions struct {
	// Whether this user-agent rule is currently disabled.
	Paused *bool `json:"paused,omitempty"`

	// Some useful information about this rule to help identify the purpose of it.
	Description *string `json:"description,omitempty"`

	// The type of action to perform.
	Mode *string `json:"mode,omitempty"`

	// Target/Value pair to use for this rule. The value is the exact UserAgent to match.
	Configuration *UseragentRuleInputConfiguration `json:"configuration,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateZoneUserAgentRuleOptions.Mode property.
// The type of action to perform.
const (
	CreateZoneUserAgentRuleOptions_Mode_Block = "block"
	CreateZoneUserAgentRuleOptions_Mode_Challenge = "challenge"
	CreateZoneUserAgentRuleOptions_Mode_JsChallenge = "js_challenge"
)

// NewCreateZoneUserAgentRuleOptions : Instantiate CreateZoneUserAgentRuleOptions
func (*UserAgentBlockingRulesV1) NewCreateZoneUserAgentRuleOptions() *CreateZoneUserAgentRuleOptions {
	return &CreateZoneUserAgentRuleOptions{}
}

// SetPaused : Allow user to set Paused
func (options *CreateZoneUserAgentRuleOptions) SetPaused(paused bool) *CreateZoneUserAgentRuleOptions {
	options.Paused = core.BoolPtr(paused)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateZoneUserAgentRuleOptions) SetDescription(description string) *CreateZoneUserAgentRuleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetMode : Allow user to set Mode
func (options *CreateZoneUserAgentRuleOptions) SetMode(mode string) *CreateZoneUserAgentRuleOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetConfiguration : Allow user to set Configuration
func (options *CreateZoneUserAgentRuleOptions) SetConfiguration(configuration *UseragentRuleInputConfiguration) *CreateZoneUserAgentRuleOptions {
	options.Configuration = configuration
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateZoneUserAgentRuleOptions) SetHeaders(param map[string]string) *CreateZoneUserAgentRuleOptions {
	options.Headers = param
	return options
}

// DeleteUseragentRuleRespResult : Container for response information.
type DeleteUseragentRuleRespResult struct {
	// ID.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteUseragentRuleRespResult unmarshals an instance of DeleteUseragentRuleRespResult from the specified map of raw messages.
func UnmarshalDeleteUseragentRuleRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteUseragentRuleRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteZoneUserAgentRuleOptions : The DeleteZoneUserAgentRule options.
type DeleteZoneUserAgentRuleOptions struct {
	// Identifier of the user-agent rule to be deleted.
	UseragentRuleIdentifier *string `json:"useragent_rule_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneUserAgentRuleOptions : Instantiate DeleteZoneUserAgentRuleOptions
func (*UserAgentBlockingRulesV1) NewDeleteZoneUserAgentRuleOptions(useragentRuleIdentifier string) *DeleteZoneUserAgentRuleOptions {
	return &DeleteZoneUserAgentRuleOptions{
		UseragentRuleIdentifier: core.StringPtr(useragentRuleIdentifier),
	}
}

// SetUseragentRuleIdentifier : Allow user to set UseragentRuleIdentifier
func (options *DeleteZoneUserAgentRuleOptions) SetUseragentRuleIdentifier(useragentRuleIdentifier string) *DeleteZoneUserAgentRuleOptions {
	options.UseragentRuleIdentifier = core.StringPtr(useragentRuleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneUserAgentRuleOptions) SetHeaders(param map[string]string) *DeleteZoneUserAgentRuleOptions {
	options.Headers = param
	return options
}

// GetUserAgentRuleOptions : The GetUserAgentRule options.
type GetUserAgentRuleOptions struct {
	// Identifier of user-agent blocking rule for the given zone.
	UseragentRuleIdentifier *string `json:"useragent_rule_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetUserAgentRuleOptions : Instantiate GetUserAgentRuleOptions
func (*UserAgentBlockingRulesV1) NewGetUserAgentRuleOptions(useragentRuleIdentifier string) *GetUserAgentRuleOptions {
	return &GetUserAgentRuleOptions{
		UseragentRuleIdentifier: core.StringPtr(useragentRuleIdentifier),
	}
}

// SetUseragentRuleIdentifier : Allow user to set UseragentRuleIdentifier
func (options *GetUserAgentRuleOptions) SetUseragentRuleIdentifier(useragentRuleIdentifier string) *GetUserAgentRuleOptions {
	options.UseragentRuleIdentifier = core.StringPtr(useragentRuleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetUserAgentRuleOptions) SetHeaders(param map[string]string) *GetUserAgentRuleOptions {
	options.Headers = param
	return options
}

// ListAllZoneUserAgentRulesOptions : The ListAllZoneUserAgentRules options.
type ListAllZoneUserAgentRulesOptions struct {
	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Maximum number of user-agent rules per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllZoneUserAgentRulesOptions : Instantiate ListAllZoneUserAgentRulesOptions
func (*UserAgentBlockingRulesV1) NewListAllZoneUserAgentRulesOptions() *ListAllZoneUserAgentRulesOptions {
	return &ListAllZoneUserAgentRulesOptions{}
}

// SetPage : Allow user to set Page
func (options *ListAllZoneUserAgentRulesOptions) SetPage(page int64) *ListAllZoneUserAgentRulesOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListAllZoneUserAgentRulesOptions) SetPerPage(perPage int64) *ListAllZoneUserAgentRulesOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllZoneUserAgentRulesOptions) SetHeaders(param map[string]string) *ListAllZoneUserAgentRulesOptions {
	options.Headers = param
	return options
}

// ListUseragentRulesRespResultInfo : Statistics of results.
type ListUseragentRulesRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalListUseragentRulesRespResultInfo unmarshals an instance of ListUseragentRulesRespResultInfo from the specified map of raw messages.
func UnmarshalListUseragentRulesRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListUseragentRulesRespResultInfo)
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

// UpdateUserAgentRuleOptions : The UpdateUserAgentRule options.
type UpdateUserAgentRuleOptions struct {
	// Identifier of user-agent rule.
	UseragentRuleIdentifier *string `json:"useragent_rule_identifier" validate:"required,ne="`

	// Whether this user-agent rule is currently disabled.
	Paused *bool `json:"paused,omitempty"`

	// Some useful information about this rule to help identify the purpose of it.
	Description *string `json:"description,omitempty"`

	// The type of action to perform.
	Mode *string `json:"mode,omitempty"`

	// Target/Value pair to use for this rule. The value is the exact UserAgent to match.
	Configuration *UseragentRuleInputConfiguration `json:"configuration,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateUserAgentRuleOptions.Mode property.
// The type of action to perform.
const (
	UpdateUserAgentRuleOptions_Mode_Block = "block"
	UpdateUserAgentRuleOptions_Mode_Challenge = "challenge"
	UpdateUserAgentRuleOptions_Mode_JsChallenge = "js_challenge"
)

// NewUpdateUserAgentRuleOptions : Instantiate UpdateUserAgentRuleOptions
func (*UserAgentBlockingRulesV1) NewUpdateUserAgentRuleOptions(useragentRuleIdentifier string) *UpdateUserAgentRuleOptions {
	return &UpdateUserAgentRuleOptions{
		UseragentRuleIdentifier: core.StringPtr(useragentRuleIdentifier),
	}
}

// SetUseragentRuleIdentifier : Allow user to set UseragentRuleIdentifier
func (options *UpdateUserAgentRuleOptions) SetUseragentRuleIdentifier(useragentRuleIdentifier string) *UpdateUserAgentRuleOptions {
	options.UseragentRuleIdentifier = core.StringPtr(useragentRuleIdentifier)
	return options
}

// SetPaused : Allow user to set Paused
func (options *UpdateUserAgentRuleOptions) SetPaused(paused bool) *UpdateUserAgentRuleOptions {
	options.Paused = core.BoolPtr(paused)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateUserAgentRuleOptions) SetDescription(description string) *UpdateUserAgentRuleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetMode : Allow user to set Mode
func (options *UpdateUserAgentRuleOptions) SetMode(mode string) *UpdateUserAgentRuleOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetConfiguration : Allow user to set Configuration
func (options *UpdateUserAgentRuleOptions) SetConfiguration(configuration *UseragentRuleInputConfiguration) *UpdateUserAgentRuleOptions {
	options.Configuration = configuration
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateUserAgentRuleOptions) SetHeaders(param map[string]string) *UpdateUserAgentRuleOptions {
	options.Headers = param
	return options
}

// UseragentRuleInputConfiguration : Target/Value pair to use for this rule. The value is the exact UserAgent to match.
type UseragentRuleInputConfiguration struct {
	// properties.
	Target *string `json:"target" validate:"required"`

	// The exact UserAgent string to match with this rule.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the UseragentRuleInputConfiguration.Target property.
// properties.
const (
	UseragentRuleInputConfiguration_Target_Ua = "ua"
)


// NewUseragentRuleInputConfiguration : Instantiate UseragentRuleInputConfiguration (Generic Model Constructor)
func (*UserAgentBlockingRulesV1) NewUseragentRuleInputConfiguration(target string, value string) (model *UseragentRuleInputConfiguration, err error) {
	model = &UseragentRuleInputConfiguration{
		Target: core.StringPtr(target),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalUseragentRuleInputConfiguration unmarshals an instance of UseragentRuleInputConfiguration from the specified map of raw messages.
func UnmarshalUseragentRuleInputConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UseragentRuleInputConfiguration)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
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

// UseragentRuleObjectConfiguration : Target/Value pair to use for this rule. The value is the exact UserAgent to match.
type UseragentRuleObjectConfiguration struct {
	// properties.
	Target *string `json:"target" validate:"required"`

	// The exact UserAgent string to match with this rule.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the UseragentRuleObjectConfiguration.Target property.
// properties.
const (
	UseragentRuleObjectConfiguration_Target_Ua = "ua"
)


// UnmarshalUseragentRuleObjectConfiguration unmarshals an instance of UseragentRuleObjectConfiguration from the specified map of raw messages.
func UnmarshalUseragentRuleObjectConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UseragentRuleObjectConfiguration)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
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

// DeleteUseragentRuleResp : user agent delete response.
type DeleteUseragentRuleResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *DeleteUseragentRuleRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteUseragentRuleResp unmarshals an instance of DeleteUseragentRuleResp from the specified map of raw messages.
func UnmarshalDeleteUseragentRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteUseragentRuleResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteUseragentRuleRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListUseragentRulesResp : user agent rules response.
type ListUseragentRulesResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []UseragentRuleObject `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListUseragentRulesRespResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListUseragentRulesResp unmarshals an instance of ListUseragentRulesResp from the specified map of raw messages.
func UnmarshalListUseragentRulesResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListUseragentRulesResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalUseragentRuleObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListUseragentRulesRespResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UseragentRuleObject : user agent rule object.
type UseragentRuleObject struct {
	// Identifier of the user-agent blocking rule.
	ID *string `json:"id" validate:"required"`

	// Whether this user-agent rule is currently disabled.
	Paused *bool `json:"paused" validate:"required"`

	// Some useful information about this rule to help identify the purpose of it.
	Description *string `json:"description" validate:"required"`

	// The type of action to perform.
	Mode *string `json:"mode" validate:"required"`

	// Target/Value pair to use for this rule. The value is the exact UserAgent to match.
	Configuration *UseragentRuleObjectConfiguration `json:"configuration" validate:"required"`
}

// Constants associated with the UseragentRuleObject.Mode property.
// The type of action to perform.
const (
	UseragentRuleObject_Mode_Block = "block"
	UseragentRuleObject_Mode_Challenge = "challenge"
	UseragentRuleObject_Mode_JsChallenge = "js_challenge"
)


// UnmarshalUseragentRuleObject unmarshals an instance of UseragentRuleObject from the specified map of raw messages.
func UnmarshalUseragentRuleObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UseragentRuleObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "paused", &obj.Paused)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "configuration", &obj.Configuration, UnmarshalUseragentRuleObjectConfiguration)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UseragentRuleResp : user agent rule response.
type UseragentRuleResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// user agent rule object.
	Result *UseragentRuleObject `json:"result" validate:"required"`
}


// UnmarshalUseragentRuleResp unmarshals an instance of UseragentRuleResp from the specified map of raw messages.
func UnmarshalUseragentRuleResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UseragentRuleResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalUseragentRuleObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
