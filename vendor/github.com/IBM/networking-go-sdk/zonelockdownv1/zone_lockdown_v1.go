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
 

// Package zonelockdownv1 : Operations and models for the ZoneLockdownV1 service
package zonelockdownv1

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

// ZoneLockdownV1 : Zone Lockdown
//
// Version: 1.0.1
type ZoneLockdownV1 struct {
	Service *core.BaseService

	// Full crn of the service instance.
	Crn *string

	// Zone identifier (zone id).
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "zone_lockdown"

// ZoneLockdownV1Options : Service options
type ZoneLockdownV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full crn of the service instance.
	Crn *string `validate:"required"`

	// Zone identifier (zone id).
	ZoneIdentifier *string `validate:"required"`
}

// NewZoneLockdownV1UsingExternalConfig : constructs an instance of ZoneLockdownV1 with passed in options and external configuration.
func NewZoneLockdownV1UsingExternalConfig(options *ZoneLockdownV1Options) (zoneLockdown *ZoneLockdownV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	zoneLockdown, err = NewZoneLockdownV1(options)
	if err != nil {
		return
	}

	err = zoneLockdown.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = zoneLockdown.Service.SetServiceURL(options.URL)
	}
	return
}

// NewZoneLockdownV1 : constructs an instance of ZoneLockdownV1 with passed in options.
func NewZoneLockdownV1(options *ZoneLockdownV1Options) (service *ZoneLockdownV1, err error) {
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

	service = &ZoneLockdownV1{
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

// Clone makes a copy of "zoneLockdown" suitable for processing requests.
func (zoneLockdown *ZoneLockdownV1) Clone() *ZoneLockdownV1 {
	if core.IsNil(zoneLockdown) {
		return nil
	}
	clone := *zoneLockdown
	clone.Service = zoneLockdown.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (zoneLockdown *ZoneLockdownV1) SetServiceURL(url string) error {
	return zoneLockdown.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (zoneLockdown *ZoneLockdownV1) GetServiceURL() string {
	return zoneLockdown.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (zoneLockdown *ZoneLockdownV1) SetDefaultHeaders(headers http.Header) {
	zoneLockdown.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (zoneLockdown *ZoneLockdownV1) SetEnableGzipCompression(enableGzip bool) {
	zoneLockdown.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (zoneLockdown *ZoneLockdownV1) GetEnableGzipCompression() bool {
	return zoneLockdown.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (zoneLockdown *ZoneLockdownV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	zoneLockdown.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (zoneLockdown *ZoneLockdownV1) DisableRetries() {
	zoneLockdown.Service.DisableRetries()
}

// ListAllZoneLockownRules : List all lockdown rules
// List all lockdown rules for a zone.
func (zoneLockdown *ZoneLockdownV1) ListAllZoneLockownRules(listAllZoneLockownRulesOptions *ListAllZoneLockownRulesOptions) (result *ListLockdownResp, response *core.DetailedResponse, err error) {
	return zoneLockdown.ListAllZoneLockownRulesWithContext(context.Background(), listAllZoneLockownRulesOptions)
}

// ListAllZoneLockownRulesWithContext is an alternate form of the ListAllZoneLockownRules method which supports a Context parameter
func (zoneLockdown *ZoneLockdownV1) ListAllZoneLockownRulesWithContext(ctx context.Context, listAllZoneLockownRulesOptions *ListAllZoneLockownRulesOptions) (result *ListLockdownResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllZoneLockownRulesOptions, "listAllZoneLockownRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneLockdown.Crn,
		"zone_identifier": *zoneLockdown.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneLockdown.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneLockdown.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/lockdowns`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllZoneLockownRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_lockdown", "V1", "ListAllZoneLockownRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllZoneLockownRulesOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listAllZoneLockownRulesOptions.Page))
	}
	if listAllZoneLockownRulesOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listAllZoneLockownRulesOptions.PerPage))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneLockdown.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLockdownResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneLockdownRule : Create lockdown rule
// Create a new lockdown rule for a given zone under a service instance.
func (zoneLockdown *ZoneLockdownV1) CreateZoneLockdownRule(createZoneLockdownRuleOptions *CreateZoneLockdownRuleOptions) (result *LockdownResp, response *core.DetailedResponse, err error) {
	return zoneLockdown.CreateZoneLockdownRuleWithContext(context.Background(), createZoneLockdownRuleOptions)
}

// CreateZoneLockdownRuleWithContext is an alternate form of the CreateZoneLockdownRule method which supports a Context parameter
func (zoneLockdown *ZoneLockdownV1) CreateZoneLockdownRuleWithContext(ctx context.Context, createZoneLockdownRuleOptions *CreateZoneLockdownRuleOptions) (result *LockdownResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createZoneLockdownRuleOptions, "createZoneLockdownRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneLockdown.Crn,
		"zone_identifier": *zoneLockdown.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneLockdown.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneLockdown.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/lockdowns`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createZoneLockdownRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_lockdown", "V1", "CreateZoneLockdownRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createZoneLockdownRuleOptions.ID != nil {
		body["id"] = createZoneLockdownRuleOptions.ID
	}
	if createZoneLockdownRuleOptions.Paused != nil {
		body["paused"] = createZoneLockdownRuleOptions.Paused
	}
	if createZoneLockdownRuleOptions.Description != nil {
		body["description"] = createZoneLockdownRuleOptions.Description
	}
	if createZoneLockdownRuleOptions.Urls != nil {
		body["urls"] = createZoneLockdownRuleOptions.Urls
	}
	if createZoneLockdownRuleOptions.Configurations != nil {
		body["configurations"] = createZoneLockdownRuleOptions.Configurations
	}
	if createZoneLockdownRuleOptions.Priority != nil {
		body["priority"] = createZoneLockdownRuleOptions.Priority
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
	response, err = zoneLockdown.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLockdownResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteZoneLockdownRule : Delete lockdown rule
// Delete a lockdown rule for a particular zone, given its id.
func (zoneLockdown *ZoneLockdownV1) DeleteZoneLockdownRule(deleteZoneLockdownRuleOptions *DeleteZoneLockdownRuleOptions) (result *DeleteLockdownResp, response *core.DetailedResponse, err error) {
	return zoneLockdown.DeleteZoneLockdownRuleWithContext(context.Background(), deleteZoneLockdownRuleOptions)
}

// DeleteZoneLockdownRuleWithContext is an alternate form of the DeleteZoneLockdownRule method which supports a Context parameter
func (zoneLockdown *ZoneLockdownV1) DeleteZoneLockdownRuleWithContext(ctx context.Context, deleteZoneLockdownRuleOptions *DeleteZoneLockdownRuleOptions) (result *DeleteLockdownResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneLockdownRuleOptions, "deleteZoneLockdownRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteZoneLockdownRuleOptions, "deleteZoneLockdownRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneLockdown.Crn,
		"zone_identifier": *zoneLockdown.ZoneIdentifier,
		"lockdown_rule_identifier": *deleteZoneLockdownRuleOptions.LockdownRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneLockdown.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneLockdown.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/lockdowns/{lockdown_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteZoneLockdownRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_lockdown", "V1", "DeleteZoneLockdownRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneLockdown.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteLockdownResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetLockdown : Get lockdown rule
// For a given service instance, zone id and lockdown rule id, get the lockdown rule details.
func (zoneLockdown *ZoneLockdownV1) GetLockdown(getLockdownOptions *GetLockdownOptions) (result *LockdownResp, response *core.DetailedResponse, err error) {
	return zoneLockdown.GetLockdownWithContext(context.Background(), getLockdownOptions)
}

// GetLockdownWithContext is an alternate form of the GetLockdown method which supports a Context parameter
func (zoneLockdown *ZoneLockdownV1) GetLockdownWithContext(ctx context.Context, getLockdownOptions *GetLockdownOptions) (result *LockdownResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLockdownOptions, "getLockdownOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLockdownOptions, "getLockdownOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneLockdown.Crn,
		"zone_identifier": *zoneLockdown.ZoneIdentifier,
		"lockdown_rule_identifier": *getLockdownOptions.LockdownRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneLockdown.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneLockdown.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/lockdowns/{lockdown_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLockdownOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_lockdown", "V1", "GetLockdown")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = zoneLockdown.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLockdownResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateLockdownRule : Update lockdown rule
// Update an existing lockdown rule for a given zone under a given service instance.
func (zoneLockdown *ZoneLockdownV1) UpdateLockdownRule(updateLockdownRuleOptions *UpdateLockdownRuleOptions) (result *LockdownResp, response *core.DetailedResponse, err error) {
	return zoneLockdown.UpdateLockdownRuleWithContext(context.Background(), updateLockdownRuleOptions)
}

// UpdateLockdownRuleWithContext is an alternate form of the UpdateLockdownRule method which supports a Context parameter
func (zoneLockdown *ZoneLockdownV1) UpdateLockdownRuleWithContext(ctx context.Context, updateLockdownRuleOptions *UpdateLockdownRuleOptions) (result *LockdownResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLockdownRuleOptions, "updateLockdownRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateLockdownRuleOptions, "updateLockdownRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *zoneLockdown.Crn,
		"zone_identifier": *zoneLockdown.ZoneIdentifier,
		"lockdown_rule_identifier": *updateLockdownRuleOptions.LockdownRuleIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = zoneLockdown.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(zoneLockdown.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/firewall/lockdowns/{lockdown_rule_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateLockdownRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("zone_lockdown", "V1", "UpdateLockdownRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateLockdownRuleOptions.ID != nil {
		body["id"] = updateLockdownRuleOptions.ID
	}
	if updateLockdownRuleOptions.Paused != nil {
		body["paused"] = updateLockdownRuleOptions.Paused
	}
	if updateLockdownRuleOptions.Description != nil {
		body["description"] = updateLockdownRuleOptions.Description
	}
	if updateLockdownRuleOptions.Urls != nil {
		body["urls"] = updateLockdownRuleOptions.Urls
	}
	if updateLockdownRuleOptions.Configurations != nil {
		body["configurations"] = updateLockdownRuleOptions.Configurations
	}
	if updateLockdownRuleOptions.Priority != nil {
		body["priority"] = updateLockdownRuleOptions.Priority
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
	response, err = zoneLockdown.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLockdownResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateZoneLockdownRuleOptions : The CreateZoneLockdownRule options.
type CreateZoneLockdownRuleOptions struct {
	// Lockdown rule identifier.
	ID *string `json:"id,omitempty"`

	// Whether this zone lockdown is currently paused.
	Paused *bool `json:"paused,omitempty"`

	// A note that you can use to describe the reason for a Lockdown rule.
	Description *string `json:"description,omitempty"`

	// URLs to be included in this rule definition. Wildcards are permitted. The URL pattern entered here will be escaped
	// before use. This limits the URL to just simple wildcard patterns.
	Urls []string `json:"urls,omitempty"`

	// List of IP addresses or CIDR ranges to use for this rule. This can include any number of ip or ip_range
	// configurations that can access the provided URLs.
	Configurations []LockdownInputConfigurationsItem `json:"configurations,omitempty"`

	// firewall priority.
	Priority *int64 `json:"priority,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateZoneLockdownRuleOptions : Instantiate CreateZoneLockdownRuleOptions
func (*ZoneLockdownV1) NewCreateZoneLockdownRuleOptions() *CreateZoneLockdownRuleOptions {
	return &CreateZoneLockdownRuleOptions{}
}

// SetID : Allow user to set ID
func (options *CreateZoneLockdownRuleOptions) SetID(id string) *CreateZoneLockdownRuleOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetPaused : Allow user to set Paused
func (options *CreateZoneLockdownRuleOptions) SetPaused(paused bool) *CreateZoneLockdownRuleOptions {
	options.Paused = core.BoolPtr(paused)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateZoneLockdownRuleOptions) SetDescription(description string) *CreateZoneLockdownRuleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetUrls : Allow user to set Urls
func (options *CreateZoneLockdownRuleOptions) SetUrls(urls []string) *CreateZoneLockdownRuleOptions {
	options.Urls = urls
	return options
}

// SetConfigurations : Allow user to set Configurations
func (options *CreateZoneLockdownRuleOptions) SetConfigurations(configurations []LockdownInputConfigurationsItem) *CreateZoneLockdownRuleOptions {
	options.Configurations = configurations
	return options
}

// SetPriority : Allow user to set Priority
func (options *CreateZoneLockdownRuleOptions) SetPriority(priority int64) *CreateZoneLockdownRuleOptions {
	options.Priority = core.Int64Ptr(priority)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateZoneLockdownRuleOptions) SetHeaders(param map[string]string) *CreateZoneLockdownRuleOptions {
	options.Headers = param
	return options
}

// DeleteLockdownRespResult : Container for response information.
type DeleteLockdownRespResult struct {
	// ID.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteLockdownRespResult unmarshals an instance of DeleteLockdownRespResult from the specified map of raw messages.
func UnmarshalDeleteLockdownRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLockdownRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteZoneLockdownRuleOptions : The DeleteZoneLockdownRule options.
type DeleteZoneLockdownRuleOptions struct {
	// Identifier of the lockdown rule to be deleted.
	LockdownRuleIdentifier *string `json:"lockdown_rule_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneLockdownRuleOptions : Instantiate DeleteZoneLockdownRuleOptions
func (*ZoneLockdownV1) NewDeleteZoneLockdownRuleOptions(lockdownRuleIdentifier string) *DeleteZoneLockdownRuleOptions {
	return &DeleteZoneLockdownRuleOptions{
		LockdownRuleIdentifier: core.StringPtr(lockdownRuleIdentifier),
	}
}

// SetLockdownRuleIdentifier : Allow user to set LockdownRuleIdentifier
func (options *DeleteZoneLockdownRuleOptions) SetLockdownRuleIdentifier(lockdownRuleIdentifier string) *DeleteZoneLockdownRuleOptions {
	options.LockdownRuleIdentifier = core.StringPtr(lockdownRuleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneLockdownRuleOptions) SetHeaders(param map[string]string) *DeleteZoneLockdownRuleOptions {
	options.Headers = param
	return options
}

// GetLockdownOptions : The GetLockdown options.
type GetLockdownOptions struct {
	// Identifier of lockdown rule for the given zone.
	LockdownRuleIdentifier *string `json:"lockdown_rule_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLockdownOptions : Instantiate GetLockdownOptions
func (*ZoneLockdownV1) NewGetLockdownOptions(lockdownRuleIdentifier string) *GetLockdownOptions {
	return &GetLockdownOptions{
		LockdownRuleIdentifier: core.StringPtr(lockdownRuleIdentifier),
	}
}

// SetLockdownRuleIdentifier : Allow user to set LockdownRuleIdentifier
func (options *GetLockdownOptions) SetLockdownRuleIdentifier(lockdownRuleIdentifier string) *GetLockdownOptions {
	options.LockdownRuleIdentifier = core.StringPtr(lockdownRuleIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetLockdownOptions) SetHeaders(param map[string]string) *GetLockdownOptions {
	options.Headers = param
	return options
}

// ListAllZoneLockownRulesOptions : The ListAllZoneLockownRules options.
type ListAllZoneLockownRulesOptions struct {
	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Maximum number of lockdown rules per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllZoneLockownRulesOptions : Instantiate ListAllZoneLockownRulesOptions
func (*ZoneLockdownV1) NewListAllZoneLockownRulesOptions() *ListAllZoneLockownRulesOptions {
	return &ListAllZoneLockownRulesOptions{}
}

// SetPage : Allow user to set Page
func (options *ListAllZoneLockownRulesOptions) SetPage(page int64) *ListAllZoneLockownRulesOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListAllZoneLockownRulesOptions) SetPerPage(perPage int64) *ListAllZoneLockownRulesOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllZoneLockownRulesOptions) SetHeaders(param map[string]string) *ListAllZoneLockownRulesOptions {
	options.Headers = param
	return options
}

// ListLockdownRespResultInfo : Statistics of results.
type ListLockdownRespResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalListLockdownRespResultInfo unmarshals an instance of ListLockdownRespResultInfo from the specified map of raw messages.
func UnmarshalListLockdownRespResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLockdownRespResultInfo)
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

// LockdownInputConfigurationsItem : LockdownInputConfigurationsItem struct
type LockdownInputConfigurationsItem struct {
	// properties.
	Target *string `json:"target" validate:"required"`

	// IP addresses or CIDR, if target is "ip", then value should be an IP addresses, otherwise CIDR.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the LockdownInputConfigurationsItem.Target property.
// properties.
const (
	LockdownInputConfigurationsItem_Target_Ip = "ip"
	LockdownInputConfigurationsItem_Target_IpRange = "ip_range"
)


// NewLockdownInputConfigurationsItem : Instantiate LockdownInputConfigurationsItem (Generic Model Constructor)
func (*ZoneLockdownV1) NewLockdownInputConfigurationsItem(target string, value string) (model *LockdownInputConfigurationsItem, err error) {
	model = &LockdownInputConfigurationsItem{
		Target: core.StringPtr(target),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalLockdownInputConfigurationsItem unmarshals an instance of LockdownInputConfigurationsItem from the specified map of raw messages.
func UnmarshalLockdownInputConfigurationsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LockdownInputConfigurationsItem)
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

// LockdownObjectConfigurationsItem : LockdownObjectConfigurationsItem struct
type LockdownObjectConfigurationsItem struct {
	// target.
	Target *string `json:"target" validate:"required"`

	// IP addresses or CIDR, if target is "ip", then value should be an IP addresses, otherwise CIDR.
	Value *string `json:"value" validate:"required"`
}

// Constants associated with the LockdownObjectConfigurationsItem.Target property.
// target.
const (
	LockdownObjectConfigurationsItem_Target_Ip = "ip"
	LockdownObjectConfigurationsItem_Target_IpRange = "ip_range"
)


// UnmarshalLockdownObjectConfigurationsItem unmarshals an instance of LockdownObjectConfigurationsItem from the specified map of raw messages.
func UnmarshalLockdownObjectConfigurationsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LockdownObjectConfigurationsItem)
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

// UpdateLockdownRuleOptions : The UpdateLockdownRule options.
type UpdateLockdownRuleOptions struct {
	// Identifier of lockdown rule.
	LockdownRuleIdentifier *string `json:"lockdown_rule_identifier" validate:"required,ne="`

	// Lockdown rule identifier.
	ID *string `json:"id,omitempty"`

	// Whether this zone lockdown is currently paused.
	Paused *bool `json:"paused,omitempty"`

	// A note that you can use to describe the reason for a Lockdown rule.
	Description *string `json:"description,omitempty"`

	// URLs to be included in this rule definition. Wildcards are permitted. The URL pattern entered here will be escaped
	// before use. This limits the URL to just simple wildcard patterns.
	Urls []string `json:"urls,omitempty"`

	// List of IP addresses or CIDR ranges to use for this rule. This can include any number of ip or ip_range
	// configurations that can access the provided URLs.
	Configurations []LockdownInputConfigurationsItem `json:"configurations,omitempty"`

	// firewall priority.
	Priority *int64 `json:"priority,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateLockdownRuleOptions : Instantiate UpdateLockdownRuleOptions
func (*ZoneLockdownV1) NewUpdateLockdownRuleOptions(lockdownRuleIdentifier string) *UpdateLockdownRuleOptions {
	return &UpdateLockdownRuleOptions{
		LockdownRuleIdentifier: core.StringPtr(lockdownRuleIdentifier),
	}
}

// SetLockdownRuleIdentifier : Allow user to set LockdownRuleIdentifier
func (options *UpdateLockdownRuleOptions) SetLockdownRuleIdentifier(lockdownRuleIdentifier string) *UpdateLockdownRuleOptions {
	options.LockdownRuleIdentifier = core.StringPtr(lockdownRuleIdentifier)
	return options
}

// SetID : Allow user to set ID
func (options *UpdateLockdownRuleOptions) SetID(id string) *UpdateLockdownRuleOptions {
	options.ID = core.StringPtr(id)
	return options
}

// SetPaused : Allow user to set Paused
func (options *UpdateLockdownRuleOptions) SetPaused(paused bool) *UpdateLockdownRuleOptions {
	options.Paused = core.BoolPtr(paused)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateLockdownRuleOptions) SetDescription(description string) *UpdateLockdownRuleOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetUrls : Allow user to set Urls
func (options *UpdateLockdownRuleOptions) SetUrls(urls []string) *UpdateLockdownRuleOptions {
	options.Urls = urls
	return options
}

// SetConfigurations : Allow user to set Configurations
func (options *UpdateLockdownRuleOptions) SetConfigurations(configurations []LockdownInputConfigurationsItem) *UpdateLockdownRuleOptions {
	options.Configurations = configurations
	return options
}

// SetPriority : Allow user to set Priority
func (options *UpdateLockdownRuleOptions) SetPriority(priority int64) *UpdateLockdownRuleOptions {
	options.Priority = core.Int64Ptr(priority)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLockdownRuleOptions) SetHeaders(param map[string]string) *UpdateLockdownRuleOptions {
	options.Headers = param
	return options
}

// DeleteLockdownResp : delete lockdown response.
type DeleteLockdownResp struct {
	// Operation success flag.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *DeleteLockdownRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteLockdownResp unmarshals an instance of DeleteLockdownResp from the specified map of raw messages.
func UnmarshalDeleteLockdownResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLockdownResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteLockdownRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListLockdownResp : list lockdown response.
type ListLockdownResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []LockdownObject `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *ListLockdownRespResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListLockdownResp unmarshals an instance of ListLockdownResp from the specified map of raw messages.
func UnmarshalListLockdownResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLockdownResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLockdownObject)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalListLockdownRespResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LockdownObject : lockdown object.
type LockdownObject struct {
	// Lockdown rule identifier.
	ID *string `json:"id" validate:"required"`

	// firewall priority.
	Priority *int64 `json:"priority,omitempty"`

	// Whether this zone lockdown is currently paused.
	Paused *bool `json:"paused" validate:"required"`

	// A note that you can use to describe the reason for a Lockdown rule.
	Description *string `json:"description" validate:"required"`

	// URLs to be included in this rule definition. Wildcards are permitted. The URL pattern entered here will be escaped
	// before use. This limits the URL to just simple wildcard patterns.
	Urls []string `json:"urls" validate:"required"`

	// List of IP addresses or CIDR ranges to use for this rule. This can include any number of ip or ip_range
	// configurations that can access the provided URLs.
	Configurations []LockdownObjectConfigurationsItem `json:"configurations" validate:"required"`
}


// UnmarshalLockdownObject unmarshals an instance of LockdownObject from the specified map of raw messages.
func UnmarshalLockdownObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LockdownObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
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
	err = core.UnmarshalPrimitive(m, "urls", &obj.Urls)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "configurations", &obj.Configurations, UnmarshalLockdownObjectConfigurationsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LockdownResp : lockdown response.
type LockdownResp struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// lockdown object.
	Result *LockdownObject `json:"result" validate:"required"`
}


// UnmarshalLockdownResp unmarshals an instance of LockdownResp from the specified map of raw messages.
func UnmarshalLockdownResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LockdownResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLockdownObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
