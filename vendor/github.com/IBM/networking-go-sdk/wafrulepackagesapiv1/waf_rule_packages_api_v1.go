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
 

// Package wafrulepackagesapiv1 : Operations and models for the WafRulePackagesApiV1 service
package wafrulepackagesapiv1

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

// WafRulePackagesApiV1 : This document describes CIS WAF Rule Packages API.
//
// Version: 1.0.1
type WafRulePackagesApiV1 struct {
	Service *core.BaseService

	// Cloud resource name.
	Crn *string

	// Zone ID.
	ZoneID *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "waf_rule_packages_api"

// WafRulePackagesApiV1Options : Service options
type WafRulePackagesApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Cloud resource name.
	Crn *string `validate:"required"`

	// Zone ID.
	ZoneID *string `validate:"required"`
}

// NewWafRulePackagesApiV1UsingExternalConfig : constructs an instance of WafRulePackagesApiV1 with passed in options and external configuration.
func NewWafRulePackagesApiV1UsingExternalConfig(options *WafRulePackagesApiV1Options) (wafRulePackagesApi *WafRulePackagesApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	wafRulePackagesApi, err = NewWafRulePackagesApiV1(options)
	if err != nil {
		return
	}

	err = wafRulePackagesApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = wafRulePackagesApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewWafRulePackagesApiV1 : constructs an instance of WafRulePackagesApiV1 with passed in options.
func NewWafRulePackagesApiV1(options *WafRulePackagesApiV1Options) (service *WafRulePackagesApiV1, err error) {
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

	service = &WafRulePackagesApiV1{
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

// Clone makes a copy of "wafRulePackagesApi" suitable for processing requests.
func (wafRulePackagesApi *WafRulePackagesApiV1) Clone() *WafRulePackagesApiV1 {
	if core.IsNil(wafRulePackagesApi) {
		return nil
	}
	clone := *wafRulePackagesApi
	clone.Service = wafRulePackagesApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (wafRulePackagesApi *WafRulePackagesApiV1) SetServiceURL(url string) error {
	return wafRulePackagesApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (wafRulePackagesApi *WafRulePackagesApiV1) GetServiceURL() string {
	return wafRulePackagesApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (wafRulePackagesApi *WafRulePackagesApiV1) SetDefaultHeaders(headers http.Header) {
	wafRulePackagesApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (wafRulePackagesApi *WafRulePackagesApiV1) SetEnableGzipCompression(enableGzip bool) {
	wafRulePackagesApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (wafRulePackagesApi *WafRulePackagesApiV1) GetEnableGzipCompression() bool {
	return wafRulePackagesApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (wafRulePackagesApi *WafRulePackagesApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	wafRulePackagesApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (wafRulePackagesApi *WafRulePackagesApiV1) DisableRetries() {
	wafRulePackagesApi.Service.DisableRetries()
}

// ListWafPackages : List all WAF rule packages
// Get firewall packages for a zone.
func (wafRulePackagesApi *WafRulePackagesApiV1) ListWafPackages(listWafPackagesOptions *ListWafPackagesOptions) (result *WafPackagesResponse, response *core.DetailedResponse, err error) {
	return wafRulePackagesApi.ListWafPackagesWithContext(context.Background(), listWafPackagesOptions)
}

// ListWafPackagesWithContext is an alternate form of the ListWafPackages method which supports a Context parameter
func (wafRulePackagesApi *WafRulePackagesApiV1) ListWafPackagesWithContext(ctx context.Context, listWafPackagesOptions *ListWafPackagesOptions) (result *WafPackagesResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listWafPackagesOptions, "listWafPackagesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *wafRulePackagesApi.Crn,
		"zone_id": *wafRulePackagesApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = wafRulePackagesApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(wafRulePackagesApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/firewall/waf/packages`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listWafPackagesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("waf_rule_packages_api", "V1", "ListWafPackages")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listWafPackagesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listWafPackagesOptions.Name))
	}
	if listWafPackagesOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listWafPackagesOptions.Page))
	}
	if listWafPackagesOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listWafPackagesOptions.PerPage))
	}
	if listWafPackagesOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listWafPackagesOptions.Order))
	}
	if listWafPackagesOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listWafPackagesOptions.Direction))
	}
	if listWafPackagesOptions.Match != nil {
		builder.AddQuery("match", fmt.Sprint(*listWafPackagesOptions.Match))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = wafRulePackagesApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWafPackagesResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWafPackage : Get WAF rule package
// Get information about a single firewall package.
func (wafRulePackagesApi *WafRulePackagesApiV1) GetWafPackage(getWafPackageOptions *GetWafPackageOptions) (result *WafPackageResponse, response *core.DetailedResponse, err error) {
	return wafRulePackagesApi.GetWafPackageWithContext(context.Background(), getWafPackageOptions)
}

// GetWafPackageWithContext is an alternate form of the GetWafPackage method which supports a Context parameter
func (wafRulePackagesApi *WafRulePackagesApiV1) GetWafPackageWithContext(ctx context.Context, getWafPackageOptions *GetWafPackageOptions) (result *WafPackageResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getWafPackageOptions, "getWafPackageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getWafPackageOptions, "getWafPackageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *wafRulePackagesApi.Crn,
		"zone_id": *wafRulePackagesApi.ZoneID,
		"package_id": *getWafPackageOptions.PackageID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = wafRulePackagesApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(wafRulePackagesApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/firewall/waf/packages/{package_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getWafPackageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("waf_rule_packages_api", "V1", "GetWafPackage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = wafRulePackagesApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWafPackageResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateWafPackage : Change WAF rule package package
// Change the sensitivity and action for an anomaly detection type WAF rule package.
func (wafRulePackagesApi *WafRulePackagesApiV1) UpdateWafPackage(updateWafPackageOptions *UpdateWafPackageOptions) (result *WafPackageResponse, response *core.DetailedResponse, err error) {
	return wafRulePackagesApi.UpdateWafPackageWithContext(context.Background(), updateWafPackageOptions)
}

// UpdateWafPackageWithContext is an alternate form of the UpdateWafPackage method which supports a Context parameter
func (wafRulePackagesApi *WafRulePackagesApiV1) UpdateWafPackageWithContext(ctx context.Context, updateWafPackageOptions *UpdateWafPackageOptions) (result *WafPackageResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateWafPackageOptions, "updateWafPackageOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateWafPackageOptions, "updateWafPackageOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *wafRulePackagesApi.Crn,
		"zone_id": *wafRulePackagesApi.ZoneID,
		"package_id": *updateWafPackageOptions.PackageID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = wafRulePackagesApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(wafRulePackagesApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/firewall/waf/packages/{package_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateWafPackageOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("waf_rule_packages_api", "V1", "UpdateWafPackage")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateWafPackageOptions.Sensitivity != nil {
		body["sensitivity"] = updateWafPackageOptions.Sensitivity
	}
	if updateWafPackageOptions.ActionMode != nil {
		body["action_mode"] = updateWafPackageOptions.ActionMode
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
	response, err = wafRulePackagesApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWafPackageResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetWafPackageOptions : The GetWafPackage options.
type GetWafPackageOptions struct {
	// Package ID.
	PackageID *string `json:"package_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetWafPackageOptions : Instantiate GetWafPackageOptions
func (*WafRulePackagesApiV1) NewGetWafPackageOptions(packageID string) *GetWafPackageOptions {
	return &GetWafPackageOptions{
		PackageID: core.StringPtr(packageID),
	}
}

// SetPackageID : Allow user to set PackageID
func (options *GetWafPackageOptions) SetPackageID(packageID string) *GetWafPackageOptions {
	options.PackageID = core.StringPtr(packageID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetWafPackageOptions) SetHeaders(param map[string]string) *GetWafPackageOptions {
	options.Headers = param
	return options
}

// ListWafPackagesOptions : The ListWafPackages options.
type ListWafPackagesOptions struct {
	// Name of the firewall package.
	Name *string `json:"name,omitempty"`

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

// Constants associated with the ListWafPackagesOptions.Direction property.
// Direction to order packages.
const (
	ListWafPackagesOptions_Direction_Asc = "asc"
	ListWafPackagesOptions_Direction_Desc = "desc"
)

// Constants associated with the ListWafPackagesOptions.Match property.
// Whether to match all search requirements or at least one (any).
const (
	ListWafPackagesOptions_Match_All = "all"
	ListWafPackagesOptions_Match_Any = "any"
)

// NewListWafPackagesOptions : Instantiate ListWafPackagesOptions
func (*WafRulePackagesApiV1) NewListWafPackagesOptions() *ListWafPackagesOptions {
	return &ListWafPackagesOptions{}
}

// SetName : Allow user to set Name
func (options *ListWafPackagesOptions) SetName(name string) *ListWafPackagesOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetPage : Allow user to set Page
func (options *ListWafPackagesOptions) SetPage(page int64) *ListWafPackagesOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListWafPackagesOptions) SetPerPage(perPage int64) *ListWafPackagesOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListWafPackagesOptions) SetOrder(order string) *ListWafPackagesOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetDirection : Allow user to set Direction
func (options *ListWafPackagesOptions) SetDirection(direction string) *ListWafPackagesOptions {
	options.Direction = core.StringPtr(direction)
	return options
}

// SetMatch : Allow user to set Match
func (options *ListWafPackagesOptions) SetMatch(match string) *ListWafPackagesOptions {
	options.Match = core.StringPtr(match)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListWafPackagesOptions) SetHeaders(param map[string]string) *ListWafPackagesOptions {
	options.Headers = param
	return options
}

// UpdateWafPackageOptions : The UpdateWafPackage options.
type UpdateWafPackageOptions struct {
	// Package ID.
	PackageID *string `json:"package_id" validate:"required,ne="`

	// The sensitivity of the firewall package.
	Sensitivity *string `json:"sensitivity,omitempty"`

	// The default action that will be taken for rules under the firewall package.
	ActionMode *string `json:"action_mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateWafPackageOptions.Sensitivity property.
// The sensitivity of the firewall package.
const (
	UpdateWafPackageOptions_Sensitivity_High = "high"
	UpdateWafPackageOptions_Sensitivity_Low = "low"
	UpdateWafPackageOptions_Sensitivity_Medium = "medium"
	UpdateWafPackageOptions_Sensitivity_Off = "off"
)

// Constants associated with the UpdateWafPackageOptions.ActionMode property.
// The default action that will be taken for rules under the firewall package.
const (
	UpdateWafPackageOptions_ActionMode_Block = "block"
	UpdateWafPackageOptions_ActionMode_Challenge = "challenge"
	UpdateWafPackageOptions_ActionMode_Simulate = "simulate"
)

// NewUpdateWafPackageOptions : Instantiate UpdateWafPackageOptions
func (*WafRulePackagesApiV1) NewUpdateWafPackageOptions(packageID string) *UpdateWafPackageOptions {
	return &UpdateWafPackageOptions{
		PackageID: core.StringPtr(packageID),
	}
}

// SetPackageID : Allow user to set PackageID
func (options *UpdateWafPackageOptions) SetPackageID(packageID string) *UpdateWafPackageOptions {
	options.PackageID = core.StringPtr(packageID)
	return options
}

// SetSensitivity : Allow user to set Sensitivity
func (options *UpdateWafPackageOptions) SetSensitivity(sensitivity string) *UpdateWafPackageOptions {
	options.Sensitivity = core.StringPtr(sensitivity)
	return options
}

// SetActionMode : Allow user to set ActionMode
func (options *UpdateWafPackageOptions) SetActionMode(actionMode string) *UpdateWafPackageOptions {
	options.ActionMode = core.StringPtr(actionMode)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateWafPackageOptions) SetHeaders(param map[string]string) *UpdateWafPackageOptions {
	options.Headers = param
	return options
}

// WafPackageResponseResult : Container for response information.
type WafPackageResponseResult struct {
	// ID.
	ID *string `json:"id,omitempty"`

	// Name.
	Name *string `json:"name,omitempty"`

	// Description.
	Description *string `json:"description,omitempty"`

	// Detection mode.
	DetectionMode *string `json:"detection_mode,omitempty"`

	// Value.
	ZoneID *string `json:"zone_id,omitempty"`

	// Value.
	Status *string `json:"status,omitempty"`

	// Value.
	Sensitivity *string `json:"sensitivity,omitempty"`

	// Value.
	ActionMode *string `json:"action_mode,omitempty"`
}


// UnmarshalWafPackageResponseResult unmarshals an instance of WafPackageResponseResult from the specified map of raw messages.
func UnmarshalWafPackageResponseResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafPackageResponseResult)
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
	err = core.UnmarshalPrimitive(m, "detection_mode", &obj.DetectionMode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_id", &obj.ZoneID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sensitivity", &obj.Sensitivity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "action_mode", &obj.ActionMode)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WafPackagesResponseResultInfo : Statistics of results.
type WafPackagesResponseResultInfo struct {
	// Page number.
	Page *int64 `json:"page" validate:"required"`

	// Number of results per page.
	PerPage *int64 `json:"per_page" validate:"required"`

	// Number of results.
	Count *int64 `json:"count" validate:"required"`

	// Total number of results.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalWafPackagesResponseResultInfo unmarshals an instance of WafPackagesResponseResultInfo from the specified map of raw messages.
func UnmarshalWafPackagesResponseResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafPackagesResponseResultInfo)
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

// WafPackagesResponseResultItem : WafPackagesResponseResultItem struct
type WafPackagesResponseResultItem struct {
	// ID.
	ID *string `json:"id,omitempty"`

	// Name.
	Name *string `json:"name,omitempty"`

	// Description.
	Description *string `json:"description,omitempty"`

	// Detection mode.
	DetectionMode *string `json:"detection_mode,omitempty"`

	// Value.
	ZoneID *string `json:"zone_id,omitempty"`

	// Value.
	Status *string `json:"status,omitempty"`
}


// UnmarshalWafPackagesResponseResultItem unmarshals an instance of WafPackagesResponseResultItem from the specified map of raw messages.
func UnmarshalWafPackagesResponseResultItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafPackagesResponseResultItem)
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
	err = core.UnmarshalPrimitive(m, "detection_mode", &obj.DetectionMode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_id", &obj.ZoneID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WafPackageResponse : waf package response.
type WafPackageResponse struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result *WafPackageResponseResult `json:"result" validate:"required"`
}


// UnmarshalWafPackageResponse unmarshals an instance of WafPackageResponse from the specified map of raw messages.
func UnmarshalWafPackageResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafPackageResponse)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalWafPackageResponseResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WafPackagesResponse : waf packages response.
type WafPackagesResponse struct {
	// Was operation successful.
	Success *bool `json:"success" validate:"required"`

	// Array of errors encountered.
	Errors [][]string `json:"errors" validate:"required"`

	// Array of messages returned.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for response information.
	Result []WafPackagesResponseResultItem `json:"result" validate:"required"`

	// Statistics of results.
	ResultInfo *WafPackagesResponseResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalWafPackagesResponse unmarshals an instance of WafPackagesResponse from the specified map of raw messages.
func UnmarshalWafPackagesResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WafPackagesResponse)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalWafPackagesResponseResultItem)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalWafPackagesResponseResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
