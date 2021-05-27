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
 

// Package wafrulegroupsapiv1 : Operations and models for the WafRuleGroupsApiV1 service
package wafrulegroupsapiv1

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

// WafRuleGroupsApiV1 : This document describes CIS WAF Rule Groups API.
//
// Version: 1.0.0
type WafRuleGroupsApiV1 struct {
	Service *core.BaseService

	// cloud resource name.
	Crn *string

	// Zone ID.
	ZoneID *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "waf_rule_groups_api"

// WafRuleGroupsApiV1Options : Service options
type WafRuleGroupsApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// cloud resource name.
	Crn *string `validate:"required"`

	// Zone ID.
	ZoneID *string `validate:"required"`
}

// NewWafRuleGroupsApiV1UsingExternalConfig : constructs an instance of WafRuleGroupsApiV1 with passed in options and external configuration.
func NewWafRuleGroupsApiV1UsingExternalConfig(options *WafRuleGroupsApiV1Options) (wafRuleGroupsApi *WafRuleGroupsApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	wafRuleGroupsApi, err = NewWafRuleGroupsApiV1(options)
	if err != nil {
		return
	}

	err = wafRuleGroupsApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = wafRuleGroupsApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewWafRuleGroupsApiV1 : constructs an instance of WafRuleGroupsApiV1 with passed in options.
func NewWafRuleGroupsApiV1(options *WafRuleGroupsApiV1Options) (service *WafRuleGroupsApiV1, err error) {
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

	service = &WafRuleGroupsApiV1{
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

// Clone makes a copy of "wafRuleGroupsApi" suitable for processing requests.
func (wafRuleGroupsApi *WafRuleGroupsApiV1) Clone() *WafRuleGroupsApiV1 {
	if core.IsNil(wafRuleGroupsApi) {
		return nil
	}
	clone := *wafRuleGroupsApi
	clone.Service = wafRuleGroupsApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (wafRuleGroupsApi *WafRuleGroupsApiV1) SetServiceURL(url string) error {
	return wafRuleGroupsApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (wafRuleGroupsApi *WafRuleGroupsApiV1) GetServiceURL() string {
	return wafRuleGroupsApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (wafRuleGroupsApi *WafRuleGroupsApiV1) SetDefaultHeaders(headers http.Header) {
	wafRuleGroupsApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (wafRuleGroupsApi *WafRuleGroupsApiV1) SetEnableGzipCompression(enableGzip bool) {
	wafRuleGroupsApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (wafRuleGroupsApi *WafRuleGroupsApiV1) GetEnableGzipCompression() bool {
	return wafRuleGroupsApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (wafRuleGroupsApi *WafRuleGroupsApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	wafRuleGroupsApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (wafRuleGroupsApi *WafRuleGroupsApiV1) DisableRetries() {
	wafRuleGroupsApi.Service.DisableRetries()
}

// ListWafRuleGroups : List all WAF rule groups
// List all WAF rule groups contained within a package.
func (wafRuleGroupsApi *WafRuleGroupsApiV1) ListWafRuleGroups(listWafRuleGroupsOptions *ListWafRuleGroupsOptions) (result *WafGroupsResponse, response *core.DetailedResponse, err error) {
	return wafRuleGroupsApi.ListWafRuleGroupsWithContext(context.Background(), listWafRuleGroupsOptions)
}

// ListWafRuleGroupsWithContext is an alternate form of the ListWafRuleGroups method which supports a Context parameter
func (wafRuleGroupsApi *WafRuleGroupsApiV1) ListWafRuleGroupsWithContext(ctx context.Context, listWafRuleGroupsOptions *ListWafRuleGroupsOptions) (result *WafGroupsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listWafRuleGroupsOptions, "listWafRuleGroupsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listWafRuleGroupsOptions, "listWafRuleGroupsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *wafRuleGroupsApi.Crn,
		"zone_id": *wafRuleGroupsApi.ZoneID,
		"pkg_id": *listWafRuleGroupsOptions.PkgID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = wafRuleGroupsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(wafRuleGroupsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/firewall/waf/packages/{pkg_id}/groups`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listWafRuleGroupsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("waf_rule_groups_api", "V1", "ListWafRuleGroups")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listWafRuleGroupsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listWafRuleGroupsOptions.Name))
	}
	if listWafRuleGroupsOptions.Mode != nil {
		builder.AddQuery("mode", fmt.Sprint(*listWafRuleGroupsOptions.Mode))
	}
	if listWafRuleGroupsOptions.RulesCount != nil {
		builder.AddQuery("rules_count", fmt.Sprint(*listWafRuleGroupsOptions.RulesCount))
	}
	if listWafRuleGroupsOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listWafRuleGroupsOptions.Page))
	}
	if listWafRuleGroupsOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listWafRuleGroupsOptions.PerPage))
	}
	if listWafRuleGroupsOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listWafRuleGroupsOptions.Order))
	}
	if listWafRuleGroupsOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listWafRuleGroupsOptions.Direction))
	}
	if listWafRuleGroupsOptions.Match != nil {
		builder.AddQuery("match", fmt.Sprint(*listWafRuleGroupsOptions.Match))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = wafRuleGroupsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWafGroupsResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWafRuleGroup : Get WAF rule group
// Get a single WAF rule group.
func (wafRuleGroupsApi *WafRuleGroupsApiV1) GetWafRuleGroup(getWafRuleGroupOptions *GetWafRuleGroupOptions) (result *WafGroupResponse, response *core.DetailedResponse, err error) {
	return wafRuleGroupsApi.GetWafRuleGroupWithContext(context.Background(), getWafRuleGroupOptions)
}

// GetWafRuleGroupWithContext is an alternate form of the GetWafRuleGroup method which supports a Context parameter
func (wafRuleGroupsApi *WafRuleGroupsApiV1) GetWafRuleGroupWithContext(ctx context.Context, getWafRuleGroupOptions *GetWafRuleGroupOptions) (result *WafGroupResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWafRuleGroupOptions, "getWafRuleGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWafRuleGroupOptions, "getWafRuleGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *wafRuleGroupsApi.Crn,
		"zone_id": *wafRuleGroupsApi.ZoneID,
		"pkg_id": *getWafRuleGroupOptions.PkgID,
		"group_id": *getWafRuleGroupOptions.GroupID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = wafRuleGroupsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(wafRuleGroupsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/firewall/waf/packages/{pkg_id}/groups/{group_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWafRuleGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("waf_rule_groups_api", "V1", "GetWafRuleGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = wafRuleGroupsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWafGroupResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateWafRuleGroup : Update WAF rule group
// Update the state of a WAF rule group.
func (wafRuleGroupsApi *WafRuleGroupsApiV1) UpdateWafRuleGroup(updateWafRuleGroupOptions *UpdateWafRuleGroupOptions) (result *WafGroupResponse, response *core.DetailedResponse, err error) {
	return wafRuleGroupsApi.UpdateWafRuleGroupWithContext(context.Background(), updateWafRuleGroupOptions)
}

// UpdateWafRuleGroupWithContext is an alternate form of the UpdateWafRuleGroup method which supports a Context parameter
func (wafRuleGroupsApi *WafRuleGroupsApiV1) UpdateWafRuleGroupWithContext(ctx context.Context, updateWafRuleGroupOptions *UpdateWafRuleGroupOptions) (result *WafGroupResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateWafRuleGroupOptions, "updateWafRuleGroupOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateWafRuleGroupOptions, "updateWafRuleGroupOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *wafRuleGroupsApi.Crn,
		"zone_id": *wafRuleGroupsApi.ZoneID,
		"pkg_id": *updateWafRuleGroupOptions.PkgID,
		"group_id": *updateWafRuleGroupOptions.GroupID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = wafRuleGroupsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(wafRuleGroupsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/firewall/waf/packages/{pkg_id}/groups/{group_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateWafRuleGroupOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("waf_rule_groups_api", "V1", "UpdateWafRuleGroup")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateWafRuleGroupOptions.Mode != nil {
		body["mode"] = updateWafRuleGroupOptions.Mode
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
	response, err = wafRuleGroupsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWafGroupResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWafRuleGroupOptions : The GetWafRuleGroup options.
type GetWafRuleGroupOptions struct {
	// Package ID.
	PkgID *string `json:"pkg_id" validate:"required,ne="`

	// Group ID.
	GroupID *string `json:"group_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWafRuleGroupOptions : Instantiate GetWafRuleGroupOptions
func (*WafRuleGroupsApiV1) NewGetWafRuleGroupOptions(pkgID string, groupID string) *GetWafRuleGroupOptions {
	return &GetWafRuleGroupOptions{
		PkgID: core.StringPtr(pkgID),
		GroupID: core.StringPtr(groupID),
	}
}

// SetPkgID : Allow user to set PkgID
func (options *GetWafRuleGroupOptions) SetPkgID(pkgID string) *GetWafRuleGroupOptions {
	options.PkgID = core.StringPtr(pkgID)
	return options
}

// SetGroupID : Allow user to set GroupID
func (options *GetWafRuleGroupOptions) SetGroupID(groupID string) *GetWafRuleGroupOptions {
	options.GroupID = core.StringPtr(groupID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWafRuleGroupOptions) SetHeaders(param map[string]string) *GetWafRuleGroupOptions {
	options.Headers = param
	return options
}

// ListWafRuleGroupsOptions : The ListWafRuleGroups options.
type ListWafRuleGroupsOptions struct {
	// Package ID.
	PkgID *string `json:"pkg_id" validate:"required,ne="`

	// Name of the firewall package.
	Name *string `json:"name,omitempty"`

	// Whether or not the rules contained within this group are configurable/usable.
	Mode *string `json:"mode,omitempty"`

	// How many rules are contained within this group.
	RulesCount *string `json:"rules_count,omitempty"`

	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Number of packages per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Field to order packages by.
	Order *string `json:"order,omitempty"`

	// Direction to order packages.
	Direction *string `json:"direction,omitempty"`

	// Whether to match all search requirements or at least one (any).
	Match *string `json:"match,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListWafRuleGroupsOptions.Mode property.
// Whether or not the rules contained within this group are configurable/usable.
const (
	ListWafRuleGroupsOptions_Mode_Off = "off"
	ListWafRuleGroupsOptions_Mode_On = "on"
)

// Constants associated with the ListWafRuleGroupsOptions.Direction property.
// Direction to order packages.
const (
	ListWafRuleGroupsOptions_Direction_Asc = "asc"
	ListWafRuleGroupsOptions_Direction_Desc = "desc"
)

// Constants associated with the ListWafRuleGroupsOptions.Match property.
// Whether to match all search requirements or at least one (any).
const (
	ListWafRuleGroupsOptions_Match_All = "all"
	ListWafRuleGroupsOptions_Match_Any = "any"
)

// NewListWafRuleGroupsOptions : Instantiate ListWafRuleGroupsOptions
func (*WafRuleGroupsApiV1) NewListWafRuleGroupsOptions(pkgID string) *ListWafRuleGroupsOptions {
	return &ListWafRuleGroupsOptions{
		PkgID: core.StringPtr(pkgID),
	}
}

// SetPkgID : Allow user to set PkgID
func (options *ListWafRuleGroupsOptions) SetPkgID(pkgID string) *ListWafRuleGroupsOptions {
	options.PkgID = core.StringPtr(pkgID)
	return options
}

// SetName : Allow user to set Name
func (options *ListWafRuleGroupsOptions) SetName(name string) *ListWafRuleGroupsOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetMode : Allow user to set Mode
func (options *ListWafRuleGroupsOptions) SetMode(mode string) *ListWafRuleGroupsOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetRulesCount : Allow user to set RulesCount
func (options *ListWafRuleGroupsOptions) SetRulesCount(rulesCount string) *ListWafRuleGroupsOptions {
	options.RulesCount = core.StringPtr(rulesCount)
	return options
}

// SetPage : Allow user to set Page
func (options *ListWafRuleGroupsOptions) SetPage(page int64) *ListWafRuleGroupsOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListWafRuleGroupsOptions) SetPerPage(perPage int64) *ListWafRuleGroupsOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListWafRuleGroupsOptions) SetOrder(order string) *ListWafRuleGroupsOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetDirection : Allow user to set Direction
func (options *ListWafRuleGroupsOptions) SetDirection(direction string) *ListWafRuleGroupsOptions {
	options.Direction = core.StringPtr(direction)
	return options
}

// SetMatch : Allow user to set Match
func (options *ListWafRuleGroupsOptions) SetMatch(match string) *ListWafRuleGroupsOptions {
	options.Match = core.StringPtr(match)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListWafRuleGroupsOptions) SetHeaders(param map[string]string) *ListWafRuleGroupsOptions {
	options.Headers = param
	return options
}

// UpdateWafRuleGroupOptions : The UpdateWafRuleGroup options.
type UpdateWafRuleGroupOptions struct {
	// Package ID.
	PkgID *string `json:"pkg_id" validate:"required,ne="`

	// Group ID.
	GroupID *string `json:"group_id" validate:"required,ne="`

	// Whether or not the rules contained within this group are configurable/usable.
	Mode *string `json:"mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateWafRuleGroupOptions.Mode property.
// Whether or not the rules contained within this group are configurable/usable.
const (
	UpdateWafRuleGroupOptions_Mode_Off = "off"
	UpdateWafRuleGroupOptions_Mode_On = "on"
)

// NewUpdateWafRuleGroupOptions : Instantiate UpdateWafRuleGroupOptions
func (*WafRuleGroupsApiV1) NewUpdateWafRuleGroupOptions(pkgID string, groupID string) *UpdateWafRuleGroupOptions {
	return &UpdateWafRuleGroupOptions{
		PkgID: core.StringPtr(pkgID),
		GroupID: core.StringPtr(groupID),
	}
}

// SetPkgID : Allow user to set PkgID
func (options *UpdateWafRuleGroupOptions) SetPkgID(pkgID string) *UpdateWafRuleGroupOptions {
	options.PkgID = core.StringPtr(pkgID)
	return options
}

// SetGroupID : Allow user to set GroupID
func (options *UpdateWafRuleGroupOptions) SetGroupID(groupID string) *UpdateWafRuleGroupOptions {
	options.GroupID = core.StringPtr(groupID)
	return options
}

// SetMode : Allow user to set Mode
func (options *UpdateWafRuleGroupOptions) SetMode(mode string) *UpdateWafRuleGroupOptions {
	options.Mode = core.StringPtr(mode)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateWafRuleGroupOptions) SetHeaders(param map[string]string) *UpdateWafRuleGroupOptions {
	options.Headers = param
	return options
}

// WafGroupResponseResultInfo : Statistics of results.
type WafGroupResponseResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalWafGroupResponseResultInfo unmarshals an instance of WafGroupResponseResultInfo from the specified map of raw messages.
func UnmarshalWafGroupResponseResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafGroupResponseResultInfo)
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

// WafGroupsResponseResultInfo : Statistics of results.
type WafGroupsResponseResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalWafGroupsResponseResultInfo unmarshals an instance of WafGroupsResponseResultInfo from the specified map of raw messages.
func UnmarshalWafGroupsResponseResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafGroupsResponseResultInfo)
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

// WafGroupResponse : waf group response.
type WafGroupResponse struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// waf rule properties.
	Result *WafRuleProperties `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *WafGroupResponseResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalWafGroupResponse unmarshals an instance of WafGroupResponse from the specified map of raw messages.
func UnmarshalWafGroupResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafGroupResponse)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalWafRuleProperties)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalWafGroupResponseResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WafGroupsResponse : waf groups response.
type WafGroupsResponse struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []WafRuleProperties `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *WafGroupsResponseResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalWafGroupsResponse unmarshals an instance of WafGroupsResponse from the specified map of raw messages.
func UnmarshalWafGroupsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafGroupsResponse)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalWafRuleProperties)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalWafGroupsResponseResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WafRuleProperties : waf rule properties.
type WafRuleProperties struct {
	// ID.
	ID *string `json:"id,omitempty"`

	// Name.
	Name *string `json:"name,omitempty"`

	// Description.
	Description *string `json:"description,omitempty"`

	// Number of rules.
	RulesCount *int64 `json:"rules_count,omitempty"`

	// Number of modified rules.
	ModifiedRulesCount *int64 `json:"modified_rules_count,omitempty"`

	// Package ID.
	PackageID *string `json:"package_id,omitempty"`

	// Mode.
	Mode *string `json:"mode,omitempty"`

	// Allowed Modes.
	AllowedModes []string `json:"allowed_modes,omitempty"`
}


// UnmarshalWafRuleProperties unmarshals an instance of WafRuleProperties from the specified map of raw messages.
func UnmarshalWafRuleProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafRuleProperties)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rules_count", &obj.RulesCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_rules_count", &obj.ModifiedRulesCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "package_id", &obj.PackageID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "mode", &obj.Mode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allowed_modes", &obj.AllowedModes)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
