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
 

// Package pageruleapiv1 : Operations and models for the PageRuleApiV1 service
package pageruleapiv1

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

// PageRuleApiV1 : This document describes CIS Pagerule API.
//
// Version: 1.0.0
type PageRuleApiV1 struct {
	Service *core.BaseService

	// instance id.
	Crn *string

	// zone id.
	ZoneID *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "page_rule_api"

// PageRuleApiV1Options : Service options
type PageRuleApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// instance id.
	Crn *string `validate:"required"`

	// zone id.
	ZoneID *string `validate:"required"`
}

// NewPageRuleApiV1UsingExternalConfig : constructs an instance of PageRuleApiV1 with passed in options and external configuration.
func NewPageRuleApiV1UsingExternalConfig(options *PageRuleApiV1Options) (pageRuleApi *PageRuleApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	pageRuleApi, err = NewPageRuleApiV1(options)
	if err != nil {
		return
	}

	err = pageRuleApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = pageRuleApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewPageRuleApiV1 : constructs an instance of PageRuleApiV1 with passed in options.
func NewPageRuleApiV1(options *PageRuleApiV1Options) (service *PageRuleApiV1, err error) {
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

	service = &PageRuleApiV1{
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

// Clone makes a copy of "pageRuleApi" suitable for processing requests.
func (pageRuleApi *PageRuleApiV1) Clone() *PageRuleApiV1 {
	if core.IsNil(pageRuleApi) {
		return nil
	}
	clone := *pageRuleApi
	clone.Service = pageRuleApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (pageRuleApi *PageRuleApiV1) SetServiceURL(url string) error {
	return pageRuleApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (pageRuleApi *PageRuleApiV1) GetServiceURL() string {
	return pageRuleApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (pageRuleApi *PageRuleApiV1) SetDefaultHeaders(headers http.Header) {
	pageRuleApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (pageRuleApi *PageRuleApiV1) SetEnableGzipCompression(enableGzip bool) {
	pageRuleApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (pageRuleApi *PageRuleApiV1) GetEnableGzipCompression() bool {
	return pageRuleApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (pageRuleApi *PageRuleApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	pageRuleApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (pageRuleApi *PageRuleApiV1) DisableRetries() {
	pageRuleApi.Service.DisableRetries()
}

// GetPageRule : Get page rule
// Get a page rule details.
func (pageRuleApi *PageRuleApiV1) GetPageRule(getPageRuleOptions *GetPageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	return pageRuleApi.GetPageRuleWithContext(context.Background(), getPageRuleOptions)
}

// GetPageRuleWithContext is an alternate form of the GetPageRule method which supports a Context parameter
func (pageRuleApi *PageRuleApiV1) GetPageRuleWithContext(ctx context.Context, getPageRuleOptions *GetPageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPageRuleOptions, "getPageRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPageRuleOptions, "getPageRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *pageRuleApi.Crn,
		"zone_id": *pageRuleApi.ZoneID,
		"rule_id": *getPageRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = pageRuleApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(pageRuleApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/pagerules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getPageRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("page_rule_api", "V1", "GetPageRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = pageRuleApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPageRulesResponseWithoutResultInfo)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ChangePageRule : Change page rule
// Change a page rule.
func (pageRuleApi *PageRuleApiV1) ChangePageRule(changePageRuleOptions *ChangePageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	return pageRuleApi.ChangePageRuleWithContext(context.Background(), changePageRuleOptions)
}

// ChangePageRuleWithContext is an alternate form of the ChangePageRule method which supports a Context parameter
func (pageRuleApi *PageRuleApiV1) ChangePageRuleWithContext(ctx context.Context, changePageRuleOptions *ChangePageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(changePageRuleOptions, "changePageRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(changePageRuleOptions, "changePageRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *pageRuleApi.Crn,
		"zone_id": *pageRuleApi.ZoneID,
		"rule_id": *changePageRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = pageRuleApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(pageRuleApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/pagerules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range changePageRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("page_rule_api", "V1", "ChangePageRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if changePageRuleOptions.Targets != nil {
		body["targets"] = changePageRuleOptions.Targets
	}
	if changePageRuleOptions.Actions != nil {
		body["actions"] = changePageRuleOptions.Actions
	}
	if changePageRuleOptions.Priority != nil {
		body["priority"] = changePageRuleOptions.Priority
	}
	if changePageRuleOptions.Status != nil {
		body["status"] = changePageRuleOptions.Status
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
	response, err = pageRuleApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPageRulesResponseWithoutResultInfo)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdatePageRule : Update page rule
// Replace a page rule. The final rule will exactly match the data passed with this request.
func (pageRuleApi *PageRuleApiV1) UpdatePageRule(updatePageRuleOptions *UpdatePageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	return pageRuleApi.UpdatePageRuleWithContext(context.Background(), updatePageRuleOptions)
}

// UpdatePageRuleWithContext is an alternate form of the UpdatePageRule method which supports a Context parameter
func (pageRuleApi *PageRuleApiV1) UpdatePageRuleWithContext(ctx context.Context, updatePageRuleOptions *UpdatePageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePageRuleOptions, "updatePageRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updatePageRuleOptions, "updatePageRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *pageRuleApi.Crn,
		"zone_id": *pageRuleApi.ZoneID,
		"rule_id": *updatePageRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = pageRuleApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(pageRuleApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/pagerules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updatePageRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("page_rule_api", "V1", "UpdatePageRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updatePageRuleOptions.Targets != nil {
		body["targets"] = updatePageRuleOptions.Targets
	}
	if updatePageRuleOptions.Actions != nil {
		body["actions"] = updatePageRuleOptions.Actions
	}
	if updatePageRuleOptions.Priority != nil {
		body["priority"] = updatePageRuleOptions.Priority
	}
	if updatePageRuleOptions.Status != nil {
		body["status"] = updatePageRuleOptions.Status
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
	response, err = pageRuleApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPageRulesResponseWithoutResultInfo)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeletePageRule : Delete page rule
// Delete a page rule.
func (pageRuleApi *PageRuleApiV1) DeletePageRule(deletePageRuleOptions *DeletePageRuleOptions) (result *PageRulesDeleteResponse, response *core.DetailedResponse, err error) {
	return pageRuleApi.DeletePageRuleWithContext(context.Background(), deletePageRuleOptions)
}

// DeletePageRuleWithContext is an alternate form of the DeletePageRule method which supports a Context parameter
func (pageRuleApi *PageRuleApiV1) DeletePageRuleWithContext(ctx context.Context, deletePageRuleOptions *DeletePageRuleOptions) (result *PageRulesDeleteResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePageRuleOptions, "deletePageRuleOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deletePageRuleOptions, "deletePageRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *pageRuleApi.Crn,
		"zone_id": *pageRuleApi.ZoneID,
		"rule_id": *deletePageRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = pageRuleApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(pageRuleApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/pagerules/{rule_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deletePageRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("page_rule_api", "V1", "DeletePageRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = pageRuleApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPageRulesDeleteResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListPageRules : List page rules
// List page rules.
func (pageRuleApi *PageRuleApiV1) ListPageRules(listPageRulesOptions *ListPageRulesOptions) (result *PageRulesResponseListAll, response *core.DetailedResponse, err error) {
	return pageRuleApi.ListPageRulesWithContext(context.Background(), listPageRulesOptions)
}

// ListPageRulesWithContext is an alternate form of the ListPageRules method which supports a Context parameter
func (pageRuleApi *PageRuleApiV1) ListPageRulesWithContext(ctx context.Context, listPageRulesOptions *ListPageRulesOptions) (result *PageRulesResponseListAll, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listPageRulesOptions, "listPageRulesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *pageRuleApi.Crn,
		"zone_id": *pageRuleApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = pageRuleApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(pageRuleApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/pagerules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listPageRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("page_rule_api", "V1", "ListPageRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listPageRulesOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listPageRulesOptions.Status))
	}
	if listPageRulesOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listPageRulesOptions.Order))
	}
	if listPageRulesOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listPageRulesOptions.Direction))
	}
	if listPageRulesOptions.Match != nil {
		builder.AddQuery("match", fmt.Sprint(*listPageRulesOptions.Match))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = pageRuleApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPageRulesResponseListAll)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreatePageRule : Create page rule
// Create a page rule.
func (pageRuleApi *PageRuleApiV1) CreatePageRule(createPageRuleOptions *CreatePageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	return pageRuleApi.CreatePageRuleWithContext(context.Background(), createPageRuleOptions)
}

// CreatePageRuleWithContext is an alternate form of the CreatePageRule method which supports a Context parameter
func (pageRuleApi *PageRuleApiV1) CreatePageRuleWithContext(ctx context.Context, createPageRuleOptions *CreatePageRuleOptions) (result *PageRulesResponseWithoutResultInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createPageRuleOptions, "createPageRuleOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *pageRuleApi.Crn,
		"zone_id": *pageRuleApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = pageRuleApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(pageRuleApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/pagerules`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createPageRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("page_rule_api", "V1", "CreatePageRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createPageRuleOptions.Targets != nil {
		body["targets"] = createPageRuleOptions.Targets
	}
	if createPageRuleOptions.Actions != nil {
		body["actions"] = createPageRuleOptions.Actions
	}
	if createPageRuleOptions.Priority != nil {
		body["priority"] = createPageRuleOptions.Priority
	}
	if createPageRuleOptions.Status != nil {
		body["status"] = createPageRuleOptions.Status
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
	response, err = pageRuleApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPageRulesResponseWithoutResultInfo)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ActionsForwardingUrlValue : value.
type ActionsForwardingUrlValue struct {
	// url.
	URL *string `json:"url,omitempty"`

	// status code.
	StatusCode *int64 `json:"status_code,omitempty"`
}


// UnmarshalActionsForwardingUrlValue unmarshals an instance of ActionsForwardingUrlValue from the specified map of raw messages.
func UnmarshalActionsForwardingUrlValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ActionsForwardingUrlValue)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChangePageRuleOptions : The ChangePageRule options.
type ChangePageRuleOptions struct {
	// rule id.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// targets.
	Targets []TargetsItem `json:"targets,omitempty"`

	// actions.
	Actions []PageRulesBodyActionsItemIntf `json:"actions,omitempty"`

	// priority.
	Priority *int64 `json:"priority,omitempty"`

	// status.
	Status *string `json:"status,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewChangePageRuleOptions : Instantiate ChangePageRuleOptions
func (*PageRuleApiV1) NewChangePageRuleOptions(ruleID string) *ChangePageRuleOptions {
	return &ChangePageRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (options *ChangePageRuleOptions) SetRuleID(ruleID string) *ChangePageRuleOptions {
	options.RuleID = core.StringPtr(ruleID)
	return options
}

// SetTargets : Allow user to set Targets
func (options *ChangePageRuleOptions) SetTargets(targets []TargetsItem) *ChangePageRuleOptions {
	options.Targets = targets
	return options
}

// SetActions : Allow user to set Actions
func (options *ChangePageRuleOptions) SetActions(actions []PageRulesBodyActionsItemIntf) *ChangePageRuleOptions {
	options.Actions = actions
	return options
}

// SetPriority : Allow user to set Priority
func (options *ChangePageRuleOptions) SetPriority(priority int64) *ChangePageRuleOptions {
	options.Priority = core.Int64Ptr(priority)
	return options
}

// SetStatus : Allow user to set Status
func (options *ChangePageRuleOptions) SetStatus(status string) *ChangePageRuleOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ChangePageRuleOptions) SetHeaders(param map[string]string) *ChangePageRuleOptions {
	options.Headers = param
	return options
}

// CreatePageRuleOptions : The CreatePageRule options.
type CreatePageRuleOptions struct {
	// targets.
	Targets []TargetsItem `json:"targets,omitempty"`

	// actions.
	Actions []PageRulesBodyActionsItemIntf `json:"actions,omitempty"`

	// priority.
	Priority *int64 `json:"priority,omitempty"`

	// status.
	Status *string `json:"status,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreatePageRuleOptions : Instantiate CreatePageRuleOptions
func (*PageRuleApiV1) NewCreatePageRuleOptions() *CreatePageRuleOptions {
	return &CreatePageRuleOptions{}
}

// SetTargets : Allow user to set Targets
func (options *CreatePageRuleOptions) SetTargets(targets []TargetsItem) *CreatePageRuleOptions {
	options.Targets = targets
	return options
}

// SetActions : Allow user to set Actions
func (options *CreatePageRuleOptions) SetActions(actions []PageRulesBodyActionsItemIntf) *CreatePageRuleOptions {
	options.Actions = actions
	return options
}

// SetPriority : Allow user to set Priority
func (options *CreatePageRuleOptions) SetPriority(priority int64) *CreatePageRuleOptions {
	options.Priority = core.Int64Ptr(priority)
	return options
}

// SetStatus : Allow user to set Status
func (options *CreatePageRuleOptions) SetStatus(status string) *CreatePageRuleOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePageRuleOptions) SetHeaders(param map[string]string) *CreatePageRuleOptions {
	options.Headers = param
	return options
}

// DeletePageRuleOptions : The DeletePageRule options.
type DeletePageRuleOptions struct {
	// rule id.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeletePageRuleOptions : Instantiate DeletePageRuleOptions
func (*PageRuleApiV1) NewDeletePageRuleOptions(ruleID string) *DeletePageRuleOptions {
	return &DeletePageRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (options *DeletePageRuleOptions) SetRuleID(ruleID string) *DeletePageRuleOptions {
	options.RuleID = core.StringPtr(ruleID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePageRuleOptions) SetHeaders(param map[string]string) *DeletePageRuleOptions {
	options.Headers = param
	return options
}

// GetPageRuleOptions : The GetPageRule options.
type GetPageRuleOptions struct {
	// rule id.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetPageRuleOptions : Instantiate GetPageRuleOptions
func (*PageRuleApiV1) NewGetPageRuleOptions(ruleID string) *GetPageRuleOptions {
	return &GetPageRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (options *GetPageRuleOptions) SetRuleID(ruleID string) *GetPageRuleOptions {
	options.RuleID = core.StringPtr(ruleID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetPageRuleOptions) SetHeaders(param map[string]string) *GetPageRuleOptions {
	options.Headers = param
	return options
}

// ListPageRulesOptions : The ListPageRules options.
type ListPageRulesOptions struct {
	// default value: disabled. valid values: active, disabled.
	Status *string `json:"status,omitempty"`

	// default value: priority. valid values: status, priority.
	Order *string `json:"order,omitempty"`

	// default value: desc. valid values: asc, desc.
	Direction *string `json:"direction,omitempty"`

	// default value: all. valid values: any, all.
	Match *string `json:"match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListPageRulesOptions : Instantiate ListPageRulesOptions
func (*PageRuleApiV1) NewListPageRulesOptions() *ListPageRulesOptions {
	return &ListPageRulesOptions{}
}

// SetStatus : Allow user to set Status
func (options *ListPageRulesOptions) SetStatus(status string) *ListPageRulesOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListPageRulesOptions) SetOrder(order string) *ListPageRulesOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetDirection : Allow user to set Direction
func (options *ListPageRulesOptions) SetDirection(direction string) *ListPageRulesOptions {
	options.Direction = core.StringPtr(direction)
	return options
}

// SetMatch : Allow user to set Match
func (options *ListPageRulesOptions) SetMatch(match string) *ListPageRulesOptions {
	options.Match = core.StringPtr(match)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListPageRulesOptions) SetHeaders(param map[string]string) *ListPageRulesOptions {
	options.Headers = param
	return options
}

// PageRulesBodyActionsItem : PageRulesBodyActionsItem struct
// Models which "extend" this model:
// - PageRulesBodyActionsItemActionsSecurity
// - PageRulesBodyActionsItemActionsSecurityOptions
// - PageRulesBodyActionsItemActionsSsl
// - PageRulesBodyActionsItemActionsTTL
// - PageRulesBodyActionsItemActionsSecurityLevel
// - PageRulesBodyActionsItemActionsCacheLevel
// - PageRulesBodyActionsItemActionsEdgeCacheTTL
// - PageRulesBodyActionsItemActionsForwardingURL
// - PageRulesBodyActionsItemActionsBypassCacheOnCookie
type PageRulesBodyActionsItem struct {
	// " Page rule action field map from UI to API
	//     CF-UI                    map             API,
	// 'Disable Security'           to        'disable_security',
	// 'Browser Integrity Check'    to        'browser_check',
	// 'Server Side Excludes'       to        'server_side_exclude',
	// 'SSL'                        to        'ssl',
	// 'Browser Cache TTL'          to        'browser_cache_ttl',
	// 'Always Online'              to        'always_online',
	// 'Security Level'             to        'security_level',
	// 'Cache Level'                to        'cache_level',
	// 'Edge Cache TTL'             to        'edge_cache_ttl'
	// 'IP Geolocation Header'      to        'ip_geolocation,
	// 'Email Obfuscation'          to        'email_obfuscation',
	// 'Automatic HTTPS Rewrites'   to        'automatic_https_rewrites',
	// 'Opportunistic Encryption'   to        'opportunistic_encryption',
	// 'Forwarding URL'             to        'forwarding_url',
	// 'Always Use HTTPS'           to        'always_use_https',
	// 'Origin Cache Control'       to        'explicit_cache_control',
	// 'Bypass Cache on Cookie'     to        'bypass_cache_on_cookie',
	// 'Cache Deception Armor'      to        'cache_deception_armor',
	// 'WAF'                        to        'waf'
	//
	//                   Page rule conflict list
	// "forwarding_url"             with     all other settings for the rules
	// "always_use_https"           with     all other settings for the rules
	// "disable_security"           with     "email_obfuscation", "server_side_exclude", "waf"
	// ".
	ID *string `json:"id" validate:"required"`

	// value.
	Value interface{} `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItem.ID property.
// " Page rule action field map from UI to API
//     CF-UI                    map             API,
// 'Disable Security'           to        'disable_security',
// 'Browser Integrity Check'    to        'browser_check',
// 'Server Side Excludes'       to        'server_side_exclude',
// 'SSL'                        to        'ssl',
// 'Browser Cache TTL'          to        'browser_cache_ttl',
// 'Always Online'              to        'always_online',
// 'Security Level'             to        'security_level',
// 'Cache Level'                to        'cache_level',
// 'Edge Cache TTL'             to        'edge_cache_ttl'
// 'IP Geolocation Header'      to        'ip_geolocation,
// 'Email Obfuscation'          to        'email_obfuscation',
// 'Automatic HTTPS Rewrites'   to        'automatic_https_rewrites',
// 'Opportunistic Encryption'   to        'opportunistic_encryption',
// 'Forwarding URL'             to        'forwarding_url',
// 'Always Use HTTPS'           to        'always_use_https',
// 'Origin Cache Control'       to        'explicit_cache_control',
// 'Bypass Cache on Cookie'     to        'bypass_cache_on_cookie',
// 'Cache Deception Armor'      to        'cache_deception_armor',
// 'WAF'                        to        'waf'
//
//                   Page rule conflict list
// "forwarding_url"             with     all other settings for the rules
// "always_use_https"           with     all other settings for the rules
// "disable_security"           with     "email_obfuscation", "server_side_exclude", "waf"
// ".
const (
	PageRulesBodyActionsItem_ID_AlwaysOnline = "always_online"
	PageRulesBodyActionsItem_ID_AlwaysUseHttps = "always_use_https"
	PageRulesBodyActionsItem_ID_AutomaticHttpsRewrites = "automatic_https_rewrites"
	PageRulesBodyActionsItem_ID_BrowserCacheTTL = "browser_cache_ttl"
	PageRulesBodyActionsItem_ID_BrowserCheck = "browser_check"
	PageRulesBodyActionsItem_ID_BypassCacheOnCookie = "bypass_cache_on_cookie"
	PageRulesBodyActionsItem_ID_CacheDeceptionArmor = "cache_deception_armor"
	PageRulesBodyActionsItem_ID_CacheLevel = "cache_level"
	PageRulesBodyActionsItem_ID_DisableSecurity = "disable_security"
	PageRulesBodyActionsItem_ID_EdgeCacheTTL = "edge_cache_ttl"
	PageRulesBodyActionsItem_ID_EmailObfuscation = "email_obfuscation"
	PageRulesBodyActionsItem_ID_ExplicitCacheControl = "explicit_cache_control"
	PageRulesBodyActionsItem_ID_ForwardingURL = "forwarding_url"
	PageRulesBodyActionsItem_ID_IpGeolocation = "ip_geolocation"
	PageRulesBodyActionsItem_ID_OpportunisticEncryption = "opportunistic_encryption"
	PageRulesBodyActionsItem_ID_SecurityLevel = "security_level"
	PageRulesBodyActionsItem_ID_ServerSideExclude = "server_side_exclude"
	PageRulesBodyActionsItem_ID_Ssl = "ssl"
	PageRulesBodyActionsItem_ID_Waf = "waf"
)

func (*PageRulesBodyActionsItem) isaPageRulesBodyActionsItem() bool {
	return true
}

type PageRulesBodyActionsItemIntf interface {
	isaPageRulesBodyActionsItem() bool
}

// UnmarshalPageRulesBodyActionsItem unmarshals an instance of PageRulesBodyActionsItem from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesDeleteResponseResult : result.
type PageRulesDeleteResponseResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalPageRulesDeleteResponseResult unmarshals an instance of PageRulesDeleteResponseResult from the specified map of raw messages.
func UnmarshalPageRulesDeleteResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesDeleteResponseResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetsItem : items.
type TargetsItem struct {
	// target.
	Target *string `json:"target" validate:"required"`

	// constraint.
	Constraint *TargetsItemConstraint `json:"constraint" validate:"required"`
}


// NewTargetsItem : Instantiate TargetsItem (Generic Model Constructor)
func (*PageRuleApiV1) NewTargetsItem(target string, constraint *TargetsItemConstraint) (model *TargetsItem, err error) {
	model = &TargetsItem{
		Target: core.StringPtr(target),
		Constraint: constraint,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalTargetsItem unmarshals an instance of TargetsItem from the specified map of raw messages.
func UnmarshalTargetsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetsItem)
	err = core.UnmarshalPrimitive(m, "target", &obj.Target)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "constraint", &obj.Constraint, UnmarshalTargetsItemConstraint)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetsItemConstraint : constraint.
type TargetsItemConstraint struct {
	// operator.
	Operator *string `json:"operator" validate:"required"`

	// value.
	Value *string `json:"value" validate:"required"`
}


// NewTargetsItemConstraint : Instantiate TargetsItemConstraint (Generic Model Constructor)
func (*PageRuleApiV1) NewTargetsItemConstraint(operator string, value string) (model *TargetsItemConstraint, err error) {
	model = &TargetsItemConstraint{
		Operator: core.StringPtr(operator),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalTargetsItemConstraint unmarshals an instance of TargetsItemConstraint from the specified map of raw messages.
func UnmarshalTargetsItemConstraint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetsItemConstraint)
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
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

// UpdatePageRuleOptions : The UpdatePageRule options.
type UpdatePageRuleOptions struct {
	// rule id.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// targets.
	Targets []TargetsItem `json:"targets,omitempty"`

	// actions.
	Actions []PageRulesBodyActionsItemIntf `json:"actions,omitempty"`

	// priority.
	Priority *int64 `json:"priority,omitempty"`

	// status.
	Status *string `json:"status,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdatePageRuleOptions : Instantiate UpdatePageRuleOptions
func (*PageRuleApiV1) NewUpdatePageRuleOptions(ruleID string) *UpdatePageRuleOptions {
	return &UpdatePageRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (options *UpdatePageRuleOptions) SetRuleID(ruleID string) *UpdatePageRuleOptions {
	options.RuleID = core.StringPtr(ruleID)
	return options
}

// SetTargets : Allow user to set Targets
func (options *UpdatePageRuleOptions) SetTargets(targets []TargetsItem) *UpdatePageRuleOptions {
	options.Targets = targets
	return options
}

// SetActions : Allow user to set Actions
func (options *UpdatePageRuleOptions) SetActions(actions []PageRulesBodyActionsItemIntf) *UpdatePageRuleOptions {
	options.Actions = actions
	return options
}

// SetPriority : Allow user to set Priority
func (options *UpdatePageRuleOptions) SetPriority(priority int64) *UpdatePageRuleOptions {
	options.Priority = core.Int64Ptr(priority)
	return options
}

// SetStatus : Allow user to set Status
func (options *UpdatePageRuleOptions) SetStatus(status string) *UpdatePageRuleOptions {
	options.Status = core.StringPtr(status)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePageRuleOptions) SetHeaders(param map[string]string) *UpdatePageRuleOptions {
	options.Headers = param
	return options
}

// PageRulesDeleteResponse : page rules delete response.
type PageRulesDeleteResponse struct {
	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result *PageRulesDeleteResponseResult `json:"result" validate:"required"`
}


// UnmarshalPageRulesDeleteResponse unmarshals an instance of PageRulesDeleteResponse from the specified map of raw messages.
func UnmarshalPageRulesDeleteResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesDeleteResponse)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalPageRulesDeleteResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageRulesResponseListAll : page rule response list all.
type PageRulesResponseListAll struct {
	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result []PageRuleResult `json:"result" validate:"required"`
}


// UnmarshalPageRulesResponseListAll unmarshals an instance of PageRulesResponseListAll from the specified map of raw messages.
func UnmarshalPageRulesResponseListAll(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesResponseListAll)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalPageRuleResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageRulesResponseWithoutResultInfo : page rule response without result information.
type PageRulesResponseWithoutResultInfo struct {
	// success.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// page rule result.
	Result *PageRuleResult `json:"result" validate:"required"`
}


// UnmarshalPageRulesResponseWithoutResultInfo unmarshals an instance of PageRulesResponseWithoutResultInfo from the specified map of raw messages.
func UnmarshalPageRulesResponseWithoutResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesResponseWithoutResultInfo)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalPageRuleResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageRuleResult : page rule result.
type PageRuleResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// targets.
	Targets []TargetsItem `json:"targets" validate:"required"`

	// actions.
	Actions []PageRulesBodyActionsItemIntf `json:"actions" validate:"required"`

	// priority.
	Priority *int64 `json:"priority" validate:"required"`

	// status.
	Status *string `json:"status" validate:"required"`

	// modified date.
	ModifiedOn *string `json:"modified_on" validate:"required"`

	// created date.
	CreatedOn *string `json:"created_on" validate:"required"`
}


// UnmarshalPageRuleResult unmarshals an instance of PageRuleResult from the specified map of raw messages.
func UnmarshalPageRuleResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRuleResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "targets", &obj.Targets, UnmarshalTargetsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalPageRulesBodyActionsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "priority", &obj.Priority)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageRulesBodyActionsItemActionsBypassCacheOnCookie : bypass cache on cookie actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsBypassCacheOnCookie struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsBypassCacheOnCookie.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsBypassCacheOnCookie_ID_BypassCacheOnCookie = "bypass_cache_on_cookie"
)


// NewPageRulesBodyActionsItemActionsBypassCacheOnCookie : Instantiate PageRulesBodyActionsItemActionsBypassCacheOnCookie (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsBypassCacheOnCookie(id string) (model *PageRulesBodyActionsItemActionsBypassCacheOnCookie, err error) {
	model = &PageRulesBodyActionsItemActionsBypassCacheOnCookie{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsBypassCacheOnCookie) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsBypassCacheOnCookie unmarshals an instance of PageRulesBodyActionsItemActionsBypassCacheOnCookie from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsBypassCacheOnCookie(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsBypassCacheOnCookie)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesBodyActionsItemActionsCacheLevel : cache level actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsCacheLevel struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsCacheLevel.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsCacheLevel_ID_CacheLevel = "cache_level"
)

// Constants associated with the PageRulesBodyActionsItemActionsCacheLevel.Value property.
// value.
const (
	PageRulesBodyActionsItemActionsCacheLevel_Value_Aggressive = "aggressive"
	PageRulesBodyActionsItemActionsCacheLevel_Value_Basic = "basic"
	PageRulesBodyActionsItemActionsCacheLevel_Value_Bypass = "bypass"
	PageRulesBodyActionsItemActionsCacheLevel_Value_CacheEverything = "cache_everything"
	PageRulesBodyActionsItemActionsCacheLevel_Value_Simplified = "simplified"
)


// NewPageRulesBodyActionsItemActionsCacheLevel : Instantiate PageRulesBodyActionsItemActionsCacheLevel (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsCacheLevel(id string) (model *PageRulesBodyActionsItemActionsCacheLevel, err error) {
	model = &PageRulesBodyActionsItemActionsCacheLevel{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsCacheLevel) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsCacheLevel unmarshals an instance of PageRulesBodyActionsItemActionsCacheLevel from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsCacheLevel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsCacheLevel)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesBodyActionsItemActionsEdgeCacheTTL : edge cache ttl actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsEdgeCacheTTL struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// ttl value.
	Value *int64 `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsEdgeCacheTTL.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsEdgeCacheTTL_ID_EdgeCacheTTL = "edge_cache_ttl"
)


// NewPageRulesBodyActionsItemActionsEdgeCacheTTL : Instantiate PageRulesBodyActionsItemActionsEdgeCacheTTL (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsEdgeCacheTTL(id string) (model *PageRulesBodyActionsItemActionsEdgeCacheTTL, err error) {
	model = &PageRulesBodyActionsItemActionsEdgeCacheTTL{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsEdgeCacheTTL) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsEdgeCacheTTL unmarshals an instance of PageRulesBodyActionsItemActionsEdgeCacheTTL from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsEdgeCacheTTL(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsEdgeCacheTTL)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesBodyActionsItemActionsForwardingURL : forwarding url actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsForwardingURL struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *ActionsForwardingUrlValue `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsForwardingURL.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsForwardingURL_ID_ForwardingURL = "forwarding_url"
)


// NewPageRulesBodyActionsItemActionsForwardingURL : Instantiate PageRulesBodyActionsItemActionsForwardingURL (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsForwardingURL(id string) (model *PageRulesBodyActionsItemActionsForwardingURL, err error) {
	model = &PageRulesBodyActionsItemActionsForwardingURL{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsForwardingURL) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsForwardingURL unmarshals an instance of PageRulesBodyActionsItemActionsForwardingURL from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsForwardingURL(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsForwardingURL)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "value", &obj.Value, UnmarshalActionsForwardingUrlValue)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageRulesBodyActionsItemActionsSecurity : security actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsSecurity struct {
	// value.
	Value interface{} `json:"value,omitempty"`

	// identifier.
	ID *string `json:"id" validate:"required"`
}

// Constants associated with the PageRulesBodyActionsItemActionsSecurity.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsSecurity_ID_AlwaysUseHttps = "always_use_https"
	PageRulesBodyActionsItemActionsSecurity_ID_DisableSecurity = "disable_security"
)


// NewPageRulesBodyActionsItemActionsSecurity : Instantiate PageRulesBodyActionsItemActionsSecurity (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsSecurity(id string) (model *PageRulesBodyActionsItemActionsSecurity, err error) {
	model = &PageRulesBodyActionsItemActionsSecurity{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsSecurity) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsSecurity unmarshals an instance of PageRulesBodyActionsItemActionsSecurity from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsSecurity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsSecurity)
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PageRulesBodyActionsItemActionsSecurityLevel : security level actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsSecurityLevel struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsSecurityLevel.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsSecurityLevel_ID_SecurityLevel = "security_level"
)

// Constants associated with the PageRulesBodyActionsItemActionsSecurityLevel.Value property.
// value.
const (
	PageRulesBodyActionsItemActionsSecurityLevel_Value_EssentiallyOff = "essentially_off"
	PageRulesBodyActionsItemActionsSecurityLevel_Value_High = "high"
	PageRulesBodyActionsItemActionsSecurityLevel_Value_Low = "low"
	PageRulesBodyActionsItemActionsSecurityLevel_Value_Medium = "medium"
	PageRulesBodyActionsItemActionsSecurityLevel_Value_Off = "off"
	PageRulesBodyActionsItemActionsSecurityLevel_Value_UnderAttack = "under_attack"
)


// NewPageRulesBodyActionsItemActionsSecurityLevel : Instantiate PageRulesBodyActionsItemActionsSecurityLevel (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsSecurityLevel(id string) (model *PageRulesBodyActionsItemActionsSecurityLevel, err error) {
	model = &PageRulesBodyActionsItemActionsSecurityLevel{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsSecurityLevel) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsSecurityLevel unmarshals an instance of PageRulesBodyActionsItemActionsSecurityLevel from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsSecurityLevel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsSecurityLevel)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesBodyActionsItemActionsSecurityOptions : security options.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsSecurityOptions struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsSecurityOptions.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsSecurityOptions_ID_AlwaysOnline = "always_online"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_AutomaticHttpsRewrites = "automatic_https_rewrites"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_BrowserCheck = "browser_check"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_CacheDeceptionArmor = "cache_deception_armor"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_EmailObfuscation = "email_obfuscation"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_ExplicitCacheControl = "explicit_cache_control"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_IpGeolocation = "ip_geolocation"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_OpportunisticEncryption = "opportunistic_encryption"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_ServerSideExclude = "server_side_exclude"
	PageRulesBodyActionsItemActionsSecurityOptions_ID_Waf = "waf"
)

// Constants associated with the PageRulesBodyActionsItemActionsSecurityOptions.Value property.
// value.
const (
	PageRulesBodyActionsItemActionsSecurityOptions_Value_Off = "off"
	PageRulesBodyActionsItemActionsSecurityOptions_Value_On = "on"
)


// NewPageRulesBodyActionsItemActionsSecurityOptions : Instantiate PageRulesBodyActionsItemActionsSecurityOptions (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsSecurityOptions(id string) (model *PageRulesBodyActionsItemActionsSecurityOptions, err error) {
	model = &PageRulesBodyActionsItemActionsSecurityOptions{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsSecurityOptions) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsSecurityOptions unmarshals an instance of PageRulesBodyActionsItemActionsSecurityOptions from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsSecurityOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsSecurityOptions)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesBodyActionsItemActionsSsl : ssl actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsSsl struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsSsl.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsSsl_ID_Ssl = "ssl"
)

// Constants associated with the PageRulesBodyActionsItemActionsSsl.Value property.
// value.
const (
	PageRulesBodyActionsItemActionsSsl_Value_Flexible = "flexible"
	PageRulesBodyActionsItemActionsSsl_Value_Full = "full"
	PageRulesBodyActionsItemActionsSsl_Value_Off = "off"
	PageRulesBodyActionsItemActionsSsl_Value_OriginPull = "origin_pull"
	PageRulesBodyActionsItemActionsSsl_Value_Strict = "strict"
)


// NewPageRulesBodyActionsItemActionsSsl : Instantiate PageRulesBodyActionsItemActionsSsl (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsSsl(id string) (model *PageRulesBodyActionsItemActionsSsl, err error) {
	model = &PageRulesBodyActionsItemActionsSsl{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsSsl) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsSsl unmarshals an instance of PageRulesBodyActionsItemActionsSsl from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsSsl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsSsl)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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

// PageRulesBodyActionsItemActionsTTL : ttl actions.
// This model "extends" PageRulesBodyActionsItem
type PageRulesBodyActionsItemActionsTTL struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// value.
	Value *int64 `json:"value,omitempty"`
}

// Constants associated with the PageRulesBodyActionsItemActionsTTL.ID property.
// identifier.
const (
	PageRulesBodyActionsItemActionsTTL_ID_BrowserCacheTTL = "browser_cache_ttl"
)


// NewPageRulesBodyActionsItemActionsTTL : Instantiate PageRulesBodyActionsItemActionsTTL (Generic Model Constructor)
func (*PageRuleApiV1) NewPageRulesBodyActionsItemActionsTTL(id string) (model *PageRulesBodyActionsItemActionsTTL, err error) {
	model = &PageRulesBodyActionsItemActionsTTL{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

func (*PageRulesBodyActionsItemActionsTTL) isaPageRulesBodyActionsItem() bool {
	return true
}

// UnmarshalPageRulesBodyActionsItemActionsTTL unmarshals an instance of PageRulesBodyActionsItemActionsTTL from the specified map of raw messages.
func UnmarshalPageRulesBodyActionsItemActionsTTL(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageRulesBodyActionsItemActionsTTL)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
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
