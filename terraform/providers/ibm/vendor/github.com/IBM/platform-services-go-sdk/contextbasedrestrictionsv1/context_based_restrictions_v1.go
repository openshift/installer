/**
 * (C) Copyright IBM Corp. 2024.
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
 * IBM OpenAPI SDK Code Generator Version: 3.92.1-44330004-20240620-143510
 */

// Package contextbasedrestrictionsv1 : Operations and models for the ContextBasedRestrictionsV1 service
package contextbasedrestrictionsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/platform-services-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// ContextBasedRestrictionsV1 : With the Context Based Restrictions API, you can:
// * Create, list, get, replace, and delete network zones
// * Create, list, get, replace, and delete context-based restriction rules
// * Get account settings
//
// API Version: 1.0.1
type ContextBasedRestrictionsV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://cbr.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "context_based_restrictions"

// ContextBasedRestrictionsV1Options : Service options
type ContextBasedRestrictionsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewContextBasedRestrictionsV1UsingExternalConfig : constructs an instance of ContextBasedRestrictionsV1 with passed in options and external configuration.
func NewContextBasedRestrictionsV1UsingExternalConfig(options *ContextBasedRestrictionsV1Options) (contextBasedRestrictions *ContextBasedRestrictionsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			err = core.SDKErrorf(err, "", "env-auth-error", common.GetComponentInfo())
			return
		}
	}

	contextBasedRestrictions, err = NewContextBasedRestrictionsV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = contextBasedRestrictions.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = contextBasedRestrictions.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewContextBasedRestrictionsV1 : constructs an instance of ContextBasedRestrictionsV1 with passed in options.
func NewContextBasedRestrictionsV1(options *ContextBasedRestrictionsV1Options) (service *ContextBasedRestrictionsV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	baseService, err := core.NewBaseService(serviceOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "new-base-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
			err = core.SDKErrorf(err, "", "set-url-error", common.GetComponentInfo())
			return
		}
	}

	service = &ContextBasedRestrictionsV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "contextBasedRestrictions" suitable for processing requests.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) Clone() *ContextBasedRestrictionsV1 {
	if core.IsNil(contextBasedRestrictions) {
		return nil
	}
	clone := *contextBasedRestrictions
	clone.Service = contextBasedRestrictions.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (contextBasedRestrictions *ContextBasedRestrictionsV1) SetServiceURL(url string) error {
	err := contextBasedRestrictions.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetServiceURL() string {
	return contextBasedRestrictions.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (contextBasedRestrictions *ContextBasedRestrictionsV1) SetDefaultHeaders(headers http.Header) {
	contextBasedRestrictions.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (contextBasedRestrictions *ContextBasedRestrictionsV1) SetEnableGzipCompression(enableGzip bool) {
	contextBasedRestrictions.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetEnableGzipCompression() bool {
	return contextBasedRestrictions.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	contextBasedRestrictions.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) DisableRetries() {
	contextBasedRestrictions.Service.DisableRetries()
}

// CreateZone : Create a network zone
// This operation creates a network zone for the specified account.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) CreateZone(createZoneOptions *CreateZoneOptions) (result *Zone, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.CreateZoneWithContext(context.Background(), createZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateZoneWithContext is an alternate form of the CreateZone method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) CreateZoneWithContext(ctx context.Context, createZoneOptions *CreateZoneOptions) (result *Zone, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createZoneOptions, "createZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "CreateZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*createZoneOptions.XCorrelationID))
	}
	if createZoneOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createZoneOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createZoneOptions.Name != nil {
		body["name"] = createZoneOptions.Name
	}
	if createZoneOptions.AccountID != nil {
		body["account_id"] = createZoneOptions.AccountID
	}
	if createZoneOptions.Description != nil {
		body["description"] = createZoneOptions.Description
	}
	if createZoneOptions.Addresses != nil {
		body["addresses"] = createZoneOptions.Addresses
	}
	if createZoneOptions.Excluded != nil {
		body["excluded"] = createZoneOptions.Excluded
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListZones : List network zones
// This operation lists network zones in the specified account.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListZones(listZonesOptions *ListZonesOptions) (result *ZoneList, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.ListZonesWithContext(context.Background(), listZonesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListZonesWithContext is an alternate form of the ListZones method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListZonesWithContext(ctx context.Context, listZonesOptions *ListZonesOptions) (result *ZoneList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listZonesOptions, "listZonesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listZonesOptions, "listZonesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listZonesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "ListZones")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listZonesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*listZonesOptions.XCorrelationID))
	}
	if listZonesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listZonesOptions.TransactionID))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listZonesOptions.AccountID))
	if listZonesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listZonesOptions.Name))
	}
	if listZonesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listZonesOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_zones", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZoneList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetZone : Get a network zone
// This operation retrieves the network zone identified by the specified zone ID.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetZone(getZoneOptions *GetZoneOptions) (result *Zone, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.GetZoneWithContext(context.Background(), getZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetZoneWithContext is an alternate form of the GetZone method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetZoneWithContext(ctx context.Context, getZoneOptions *GetZoneOptions) (result *Zone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getZoneOptions, "getZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getZoneOptions, "getZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"zone_id": *getZoneOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones/{zone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "GetZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*getZoneOptions.XCorrelationID))
	}
	if getZoneOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getZoneOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceZone : Replace a network zone
// This operation replaces the network zone identified by the specified zone ID. Partial updates are not supported. The
// entire network zone object must be replaced.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ReplaceZone(replaceZoneOptions *ReplaceZoneOptions) (result *Zone, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.ReplaceZoneWithContext(context.Background(), replaceZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceZoneWithContext is an alternate form of the ReplaceZone method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ReplaceZoneWithContext(ctx context.Context, replaceZoneOptions *ReplaceZoneOptions) (result *Zone, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceZoneOptions, "replaceZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceZoneOptions, "replaceZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"zone_id": *replaceZoneOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones/{zone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "ReplaceZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceZoneOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceZoneOptions.IfMatch))
	}
	if replaceZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*replaceZoneOptions.XCorrelationID))
	}
	if replaceZoneOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*replaceZoneOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if replaceZoneOptions.Name != nil {
		body["name"] = replaceZoneOptions.Name
	}
	if replaceZoneOptions.AccountID != nil {
		body["account_id"] = replaceZoneOptions.AccountID
	}
	if replaceZoneOptions.Description != nil {
		body["description"] = replaceZoneOptions.Description
	}
	if replaceZoneOptions.Addresses != nil {
		body["addresses"] = replaceZoneOptions.Addresses
	}
	if replaceZoneOptions.Excluded != nil {
		body["excluded"] = replaceZoneOptions.Excluded
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalZone)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteZone : Delete a network zone
// This operation deletes the network zone identified by the specified zone ID.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) DeleteZone(deleteZoneOptions *DeleteZoneOptions) (response *core.DetailedResponse, err error) {
	response, err = contextBasedRestrictions.DeleteZoneWithContext(context.Background(), deleteZoneOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteZoneWithContext is an alternate form of the DeleteZone method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) DeleteZoneWithContext(ctx context.Context, deleteZoneOptions *DeleteZoneOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteZoneOptions, "deleteZoneOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteZoneOptions, "deleteZoneOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"zone_id": *deleteZoneOptions.ZoneID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones/{zone_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteZoneOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "DeleteZone")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteZoneOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*deleteZoneOptions.XCorrelationID))
	}
	if deleteZoneOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteZoneOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = contextBasedRestrictions.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_zone", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListAvailableServicerefTargets : List available service reference targets
// This operation lists all available service reference targets.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListAvailableServicerefTargets(listAvailableServicerefTargetsOptions *ListAvailableServicerefTargetsOptions) (result *ServiceRefTargetList, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.ListAvailableServicerefTargetsWithContext(context.Background(), listAvailableServicerefTargetsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAvailableServicerefTargetsWithContext is an alternate form of the ListAvailableServicerefTargets method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListAvailableServicerefTargetsWithContext(ctx context.Context, listAvailableServicerefTargetsOptions *ListAvailableServicerefTargetsOptions) (result *ServiceRefTargetList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAvailableServicerefTargetsOptions, "listAvailableServicerefTargetsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones/serviceref_targets`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAvailableServicerefTargetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "ListAvailableServicerefTargets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAvailableServicerefTargetsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*listAvailableServicerefTargetsOptions.XCorrelationID))
	}
	if listAvailableServicerefTargetsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listAvailableServicerefTargetsOptions.TransactionID))
	}

	if listAvailableServicerefTargetsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listAvailableServicerefTargetsOptions.Type))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_available_serviceref_targets", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceRefTargetList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetServicerefTarget : Get service reference target for a specified service name
// This operation gets the service reference target for a specified service name.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetServicerefTarget(getServicerefTargetOptions *GetServicerefTargetOptions) (result *ServiceRefTarget, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.GetServicerefTargetWithContext(context.Background(), getServicerefTargetOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetServicerefTargetWithContext is an alternate form of the GetServicerefTarget method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetServicerefTargetWithContext(ctx context.Context, getServicerefTargetOptions *GetServicerefTargetOptions) (result *ServiceRefTarget, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getServicerefTargetOptions, "getServicerefTargetOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getServicerefTargetOptions, "getServicerefTargetOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_name": *getServicerefTargetOptions.ServiceName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/zones/serviceref_targets/{service_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getServicerefTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "GetServicerefTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getServicerefTargetOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*getServicerefTargetOptions.XCorrelationID))
	}
	if getServicerefTargetOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getServicerefTargetOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_serviceref_target", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalServiceRefTarget)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateRule : Create a rule
// This operation creates a rule for the specified account.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) CreateRule(createRuleOptions *CreateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.CreateRuleWithContext(context.Background(), createRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateRuleWithContext is an alternate form of the CreateRule method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) CreateRuleWithContext(ctx context.Context, createRuleOptions *CreateRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createRuleOptions, "createRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/rules`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "CreateRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*createRuleOptions.XCorrelationID))
	}
	if createRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*createRuleOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if createRuleOptions.Description != nil {
		body["description"] = createRuleOptions.Description
	}
	if createRuleOptions.Contexts != nil {
		body["contexts"] = createRuleOptions.Contexts
	}
	if createRuleOptions.Resources != nil {
		body["resources"] = createRuleOptions.Resources
	}
	if createRuleOptions.Operations != nil {
		body["operations"] = createRuleOptions.Operations
	}
	if createRuleOptions.EnforcementMode != nil {
		body["enforcement_mode"] = createRuleOptions.EnforcementMode
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListRules : List rules
// This operation lists rules in the specified account.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListRules(listRulesOptions *ListRulesOptions) (result *RuleList, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.ListRulesWithContext(context.Background(), listRulesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListRulesWithContext is an alternate form of the ListRules method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListRulesWithContext(ctx context.Context, listRulesOptions *ListRulesOptions) (result *RuleList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listRulesOptions, "listRulesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listRulesOptions, "listRulesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/rules`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listRulesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "ListRules")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listRulesOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*listRulesOptions.XCorrelationID))
	}
	if listRulesOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listRulesOptions.TransactionID))
	}

	builder.AddQuery("account_id", fmt.Sprint(*listRulesOptions.AccountID))
	if listRulesOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*listRulesOptions.Region))
	}
	if listRulesOptions.Resource != nil {
		builder.AddQuery("resource", fmt.Sprint(*listRulesOptions.Resource))
	}
	if listRulesOptions.ResourceType != nil {
		builder.AddQuery("resource_type", fmt.Sprint(*listRulesOptions.ResourceType))
	}
	if listRulesOptions.ServiceInstance != nil {
		builder.AddQuery("service_instance", fmt.Sprint(*listRulesOptions.ServiceInstance))
	}
	if listRulesOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listRulesOptions.ServiceName))
	}
	if listRulesOptions.ServiceType != nil {
		builder.AddQuery("service_type", fmt.Sprint(*listRulesOptions.ServiceType))
	}
	if listRulesOptions.ServiceGroupID != nil {
		builder.AddQuery("service_group_id", fmt.Sprint(*listRulesOptions.ServiceGroupID))
	}
	if listRulesOptions.ZoneID != nil {
		builder.AddQuery("zone_id", fmt.Sprint(*listRulesOptions.ZoneID))
	}
	if listRulesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listRulesOptions.Sort))
	}
	if listRulesOptions.EnforcementMode != nil {
		builder.AddQuery("enforcement_mode", fmt.Sprint(*listRulesOptions.EnforcementMode))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_rules", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRuleList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetRule : Get a rule
// This operation retrieves the rule identified by the specified rule ID.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetRule(getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.GetRuleWithContext(context.Background(), getRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetRuleWithContext is an alternate form of the GetRule method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetRuleWithContext(ctx context.Context, getRuleOptions *GetRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRuleOptions, "getRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getRuleOptions, "getRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *getRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "GetRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*getRuleOptions.XCorrelationID))
	}
	if getRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getRuleOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceRule : Replace a rule
// This operation replaces the rule identified by the specified rule ID. Partial updates are not supported. The entire
// rule object must be replaced.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ReplaceRule(replaceRuleOptions *ReplaceRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.ReplaceRuleWithContext(context.Background(), replaceRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceRuleWithContext is an alternate form of the ReplaceRule method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ReplaceRuleWithContext(ctx context.Context, replaceRuleOptions *ReplaceRuleOptions) (result *Rule, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRuleOptions, "replaceRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceRuleOptions, "replaceRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *replaceRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "ReplaceRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceRuleOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceRuleOptions.IfMatch))
	}
	if replaceRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*replaceRuleOptions.XCorrelationID))
	}
	if replaceRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*replaceRuleOptions.TransactionID))
	}

	body := make(map[string]interface{})
	if replaceRuleOptions.Description != nil {
		body["description"] = replaceRuleOptions.Description
	}
	if replaceRuleOptions.Contexts != nil {
		body["contexts"] = replaceRuleOptions.Contexts
	}
	if replaceRuleOptions.Resources != nil {
		body["resources"] = replaceRuleOptions.Resources
	}
	if replaceRuleOptions.Operations != nil {
		body["operations"] = replaceRuleOptions.Operations
	}
	if replaceRuleOptions.EnforcementMode != nil {
		body["enforcement_mode"] = replaceRuleOptions.EnforcementMode
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRule)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteRule : Delete a rule
// This operation deletes the rule identified by the specified rule ID.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) DeleteRule(deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	response, err = contextBasedRestrictions.DeleteRuleWithContext(context.Background(), deleteRuleOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteRuleWithContext is an alternate form of the DeleteRule method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) DeleteRuleWithContext(ctx context.Context, deleteRuleOptions *DeleteRuleOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRuleOptions, "deleteRuleOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteRuleOptions, "deleteRuleOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"rule_id": *deleteRuleOptions.RuleID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/rules/{rule_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteRuleOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "DeleteRule")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if deleteRuleOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*deleteRuleOptions.XCorrelationID))
	}
	if deleteRuleOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*deleteRuleOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = contextBasedRestrictions.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_rule", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetAccountSettings : Get account settings
// This operation retrieves the settings for the specified account ID.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetAccountSettings(getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.GetAccountSettingsWithContext(context.Background(), getAccountSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAccountSettingsWithContext is an alternate form of the GetAccountSettings method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) GetAccountSettingsWithContext(ctx context.Context, getAccountSettingsOptions *GetAccountSettingsOptions) (result *AccountSettings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAccountSettingsOptions, "getAccountSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAccountSettingsOptions, "getAccountSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"account_id": *getAccountSettingsOptions.AccountID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/account_settings/{account_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAccountSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "GetAccountSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getAccountSettingsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*getAccountSettingsOptions.XCorrelationID))
	}
	if getAccountSettingsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*getAccountSettingsOptions.TransactionID))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_account_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAccountSettings)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListAvailableServiceOperations : List available service operations
// This operation lists all available service operations.
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListAvailableServiceOperations(listAvailableServiceOperationsOptions *ListAvailableServiceOperationsOptions) (result *OperationsList, response *core.DetailedResponse, err error) {
	result, response, err = contextBasedRestrictions.ListAvailableServiceOperationsWithContext(context.Background(), listAvailableServiceOperationsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAvailableServiceOperationsWithContext is an alternate form of the ListAvailableServiceOperations method which supports a Context parameter
func (contextBasedRestrictions *ContextBasedRestrictionsV1) ListAvailableServiceOperationsWithContext(ctx context.Context, listAvailableServiceOperationsOptions *ListAvailableServiceOperationsOptions) (result *OperationsList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAvailableServiceOperationsOptions, "listAvailableServiceOperationsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = contextBasedRestrictions.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(contextBasedRestrictions.Service.Options.URL, `/v1/operations`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAvailableServiceOperationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("context_based_restrictions", "V1", "ListAvailableServiceOperations")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listAvailableServiceOperationsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-Id", fmt.Sprint(*listAvailableServiceOperationsOptions.XCorrelationID))
	}
	if listAvailableServiceOperationsOptions.TransactionID != nil {
		builder.AddHeader("Transaction-Id", fmt.Sprint(*listAvailableServiceOperationsOptions.TransactionID))
	}

	if listAvailableServiceOperationsOptions.ServiceName != nil {
		builder.AddQuery("service_name", fmt.Sprint(*listAvailableServiceOperationsOptions.ServiceName))
	}
	if listAvailableServiceOperationsOptions.ServiceGroupID != nil {
		builder.AddQuery("service_group_id", fmt.Sprint(*listAvailableServiceOperationsOptions.ServiceGroupID))
	}
	if listAvailableServiceOperationsOptions.ResourceType != nil {
		builder.AddQuery("resource_type", fmt.Sprint(*listAvailableServiceOperationsOptions.ResourceType))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = contextBasedRestrictions.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_available_service_operations", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOperationsList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.1")
}

// APIType : Service API Type details.
type APIType struct {
	// The id of the API type.
	APITypeID *string `json:"api_type_id" validate:"required"`

	// The displayed name of the API type.
	DisplayName *string `json:"display_name" validate:"required"`

	// The description of the API type.
	Description *string `json:"description" validate:"required"`

	// The type of the API type.
	Type *string `json:"type" validate:"required"`

	// The actions available for the API type.
	Actions []Action `json:"actions" validate:"required"`

	// The enforcement modes supported by the API type.
	EnforcementModes []string `json:"enforcement_modes,omitempty"`
}

// UnmarshalAPIType unmarshals an instance of APIType from the specified map of raw messages.
func UnmarshalAPIType(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIType)
	err = core.UnmarshalPrimitive(m, "api_type_id", &obj.APITypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_type_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "actions", &obj.Actions, UnmarshalAction)
	if err != nil {
		err = core.SDKErrorf(err, "", "actions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enforcement_modes", &obj.EnforcementModes)
	if err != nil {
		err = core.SDKErrorf(err, "", "enforcement_modes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AccountSettings : An output account settings.
type AccountSettings struct {
	// The globally unique ID of the account settings.
	ID *string `json:"id" validate:"required"`

	// The account settings CRN.
	CRN *string `json:"crn" validate:"required"`

	// the max number of rules allowed for the account.
	RuleCountLimit *int64 `json:"rule_count_limit" validate:"required"`

	// the max number of zones allowed for the account.
	ZoneCountLimit *int64 `json:"zone_count_limit" validate:"required"`

	// the max number of rules with tags allowed for the account.
	TagsRuleCountLimit *int64 `json:"tags_rule_count_limit" validate:"required"`

	// the current number of rules used by the account.
	CurrentRuleCount *int64 `json:"current_rule_count" validate:"required"`

	// the current number of zones used by the account.
	CurrentZoneCount *int64 `json:"current_zone_count" validate:"required"`

	// the current number of rules with tags used by the account.
	CurrentTagsRuleCount *int64 `json:"current_tags_rule_count" validate:"required"`

	// The href link to the resource.
	Href *string `json:"href" validate:"required"`

	// The time the resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at,omitempty"`

	// IAM ID of the user or service which created the resource.
	CreatedByID *string `json:"created_by_id,omitempty"`

	// The last time the resource was modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at,omitempty"`

	// IAM ID of the user or service which modified the resource.
	LastModifiedByID *string `json:"last_modified_by_id,omitempty"`
}

// UnmarshalAccountSettings unmarshals an instance of AccountSettings from the specified map of raw messages.
func UnmarshalAccountSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AccountSettings)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rule_count_limit", &obj.RuleCountLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "rule_count_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "zone_count_limit", &obj.ZoneCountLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "zone_count_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags_rule_count_limit", &obj.TagsRuleCountLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags_rule_count_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "current_rule_count", &obj.CurrentRuleCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "current_rule_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "current_zone_count", &obj.CurrentZoneCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "current_zone_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "current_tags_rule_count", &obj.CurrentTagsRuleCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "current_tags_rule_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Action : Service API Type actions.
type Action struct {
	// The id of the action.
	ActionID *string `json:"action_id" validate:"required"`

	// The description of the action.
	Description *string `json:"description" validate:"required"`
}

// UnmarshalAction unmarshals an instance of Action from the specified map of raw messages.
func UnmarshalAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Action)
	err = core.UnmarshalPrimitive(m, "action_id", &obj.ActionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "action_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Address : A zone address.
// Models which "extend" this model:
// - AddressIPAddress
// - AddressIPAddressRange
// - AddressSubnet
// - AddressVPC
// - AddressServiceRef
type Address struct {
	// The type of address.
	Type *string `json:"type,omitempty"`

	// The IP address.
	Value *string `json:"value,omitempty"`

	// The address id (for use by terraform only).
	ID *string `json:"id,omitempty"`

	// A service reference value.
	Ref *ServiceRefValue `json:"ref,omitempty"`
}

// Constants associated with the Address.Type property.
// The type of address.
const (
	AddressTypeIpaddressConst  = "ipAddress"
	AddressTypeIprangeConst    = "ipRange"
	AddressTypeServicerefConst = "serviceRef"
	AddressTypeSubnetConst     = "subnet"
	AddressTypeVPCConst        = "vpc"
)

func (*Address) isaAddress() bool {
	return true
}

type AddressIntf interface {
	isaAddress() bool
}

// UnmarshalAddress unmarshals an instance of Address from the specified map of raw messages.
func UnmarshalAddress(m map[string]json.RawMessage, result interface{}) (err error) {
	// Retrieve discriminator value to determine correct "subclass".
	var discValue string
	err = core.UnmarshalPrimitive(m, "type", &discValue)
	if err != nil {
		errMsg := fmt.Sprintf("error unmarshalling discriminator property 'type': %s", err.Error())
		err = core.SDKErrorf(err, errMsg, "discriminator-unmarshal-error", common.GetComponentInfo())
		return
	}
	if discValue == "" {
		err = core.SDKErrorf(err, "required discriminator property 'type' not found in JSON object", "missing-discriminator", common.GetComponentInfo())
		return
	}
	if discValue == "ipAddress" {
		err = core.UnmarshalModel(m, "", result, UnmarshalAddressIPAddress)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-AddressIPAddress-error", common.GetComponentInfo())
		}
	} else if discValue == "ipRange" {
		err = core.UnmarshalModel(m, "", result, UnmarshalAddressIPAddressRange)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-AddressIPAddressRange-error", common.GetComponentInfo())
		}
	} else if discValue == "subnet" {
		err = core.UnmarshalModel(m, "", result, UnmarshalAddressSubnet)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-AddressSubnet-error", common.GetComponentInfo())
		}
	} else if discValue == "vpc" {
		err = core.UnmarshalModel(m, "", result, UnmarshalAddressVPC)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-AddressVPC-error", common.GetComponentInfo())
		}
	} else if discValue == "serviceRef" {
		err = core.UnmarshalModel(m, "", result, UnmarshalAddressServiceRef)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-AddressServiceRef-error", common.GetComponentInfo())
		}
	} else {
		errMsg := fmt.Sprintf("unrecognized value for discriminator property 'type': %s", discValue)
		err = core.SDKErrorf(err, errMsg, "invalid-discriminator", common.GetComponentInfo())
	}
	return
}

// CreateRuleOptions : The CreateRule options.
type CreateRuleOptions struct {
	// The description of the rule.
	Description *string `json:"description,omitempty"`

	// The contexts this rule applies to.
	Contexts []RuleContext `json:"contexts,omitempty"`

	// The resources this rule apply to.
	Resources []Resource `json:"resources,omitempty"`

	// The operations this rule applies to.
	Operations *NewRuleOperations `json:"operations,omitempty"`

	// The rule enforcement mode:
	//  * `enabled` - The restrictions are enforced and reported. This is the default.
	//  * `disabled` - The restrictions are disabled. Nothing is enforced or reported.
	//  * `report` - The restrictions are evaluated and reported, but not enforced.
	EnforcementMode *string `json:"enforcement_mode,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateRuleOptions.EnforcementMode property.
// The rule enforcement mode:
//   - `enabled` - The restrictions are enforced and reported. This is the default.
//   - `disabled` - The restrictions are disabled. Nothing is enforced or reported.
//   - `report` - The restrictions are evaluated and reported, but not enforced.
const (
	CreateRuleOptionsEnforcementModeDisabledConst = "disabled"
	CreateRuleOptionsEnforcementModeEnabledConst  = "enabled"
	CreateRuleOptionsEnforcementModeReportConst   = "report"
)

// NewCreateRuleOptions : Instantiate CreateRuleOptions
func (*ContextBasedRestrictionsV1) NewCreateRuleOptions() *CreateRuleOptions {
	return &CreateRuleOptions{}
}

// SetDescription : Allow user to set Description
func (_options *CreateRuleOptions) SetDescription(description string) *CreateRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetContexts : Allow user to set Contexts
func (_options *CreateRuleOptions) SetContexts(contexts []RuleContext) *CreateRuleOptions {
	_options.Contexts = contexts
	return _options
}

// SetResources : Allow user to set Resources
func (_options *CreateRuleOptions) SetResources(resources []Resource) *CreateRuleOptions {
	_options.Resources = resources
	return _options
}

// SetOperations : Allow user to set Operations
func (_options *CreateRuleOptions) SetOperations(operations *NewRuleOperations) *CreateRuleOptions {
	_options.Operations = operations
	return _options
}

// SetEnforcementMode : Allow user to set EnforcementMode
func (_options *CreateRuleOptions) SetEnforcementMode(enforcementMode string) *CreateRuleOptions {
	_options.EnforcementMode = core.StringPtr(enforcementMode)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateRuleOptions) SetXCorrelationID(xCorrelationID string) *CreateRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *CreateRuleOptions) SetTransactionID(transactionID string) *CreateRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRuleOptions) SetHeaders(param map[string]string) *CreateRuleOptions {
	options.Headers = param
	return options
}

// CreateZoneOptions : The CreateZone options.
type CreateZoneOptions struct {
	// The name of the zone.
	Name *string `json:"name,omitempty"`

	// The id of the account owning this zone.
	AccountID *string `json:"account_id,omitempty"`

	// The description of the zone.
	Description *string `json:"description,omitempty"`

	// The list of addresses in the zone.
	Addresses []AddressIntf `json:"addresses,omitempty"`

	// The list of excluded addresses in the zone. Only addresses of type `ipAddress`, `ipRange`, and `subnet` can be
	// excluded.
	Excluded []AddressIntf `json:"excluded,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateZoneOptions : Instantiate CreateZoneOptions
func (*ContextBasedRestrictionsV1) NewCreateZoneOptions() *CreateZoneOptions {
	return &CreateZoneOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateZoneOptions) SetName(name string) *CreateZoneOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *CreateZoneOptions) SetAccountID(accountID string) *CreateZoneOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateZoneOptions) SetDescription(description string) *CreateZoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAddresses : Allow user to set Addresses
func (_options *CreateZoneOptions) SetAddresses(addresses []AddressIntf) *CreateZoneOptions {
	_options.Addresses = addresses
	return _options
}

// SetExcluded : Allow user to set Excluded
func (_options *CreateZoneOptions) SetExcluded(excluded []AddressIntf) *CreateZoneOptions {
	_options.Excluded = excluded
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *CreateZoneOptions) SetXCorrelationID(xCorrelationID string) *CreateZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *CreateZoneOptions) SetTransactionID(transactionID string) *CreateZoneOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateZoneOptions) SetHeaders(param map[string]string) *CreateZoneOptions {
	options.Headers = param
	return options
}

// DeleteRuleOptions : The DeleteRule options.
type DeleteRuleOptions struct {
	// The ID of a rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRuleOptions : Instantiate DeleteRuleOptions
func (*ContextBasedRestrictionsV1) NewDeleteRuleOptions(ruleID string) *DeleteRuleOptions {
	return &DeleteRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *DeleteRuleOptions) SetRuleID(ruleID string) *DeleteRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteRuleOptions) SetXCorrelationID(xCorrelationID string) *DeleteRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *DeleteRuleOptions) SetTransactionID(transactionID string) *DeleteRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRuleOptions) SetHeaders(param map[string]string) *DeleteRuleOptions {
	options.Headers = param
	return options
}

// DeleteZoneOptions : The DeleteZone options.
type DeleteZoneOptions struct {
	// The ID of a zone.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteZoneOptions : Instantiate DeleteZoneOptions
func (*ContextBasedRestrictionsV1) NewDeleteZoneOptions(zoneID string) *DeleteZoneOptions {
	return &DeleteZoneOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *DeleteZoneOptions) SetZoneID(zoneID string) *DeleteZoneOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *DeleteZoneOptions) SetXCorrelationID(xCorrelationID string) *DeleteZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *DeleteZoneOptions) SetTransactionID(transactionID string) *DeleteZoneOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteZoneOptions) SetHeaders(param map[string]string) *DeleteZoneOptions {
	options.Headers = param
	return options
}

// GetAccountSettingsOptions : The GetAccountSettings options.
type GetAccountSettingsOptions struct {
	// The ID of the account the settings are for.
	AccountID *string `json:"account_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetAccountSettingsOptions : Instantiate GetAccountSettingsOptions
func (*ContextBasedRestrictionsV1) NewGetAccountSettingsOptions(accountID string) *GetAccountSettingsOptions {
	return &GetAccountSettingsOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *GetAccountSettingsOptions) SetAccountID(accountID string) *GetAccountSettingsOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetAccountSettingsOptions) SetXCorrelationID(xCorrelationID string) *GetAccountSettingsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *GetAccountSettingsOptions) SetTransactionID(transactionID string) *GetAccountSettingsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAccountSettingsOptions) SetHeaders(param map[string]string) *GetAccountSettingsOptions {
	options.Headers = param
	return options
}

// GetRuleOptions : The GetRule options.
type GetRuleOptions struct {
	// The ID of a rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRuleOptions : Instantiate GetRuleOptions
func (*ContextBasedRestrictionsV1) NewGetRuleOptions(ruleID string) *GetRuleOptions {
	return &GetRuleOptions{
		RuleID: core.StringPtr(ruleID),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *GetRuleOptions) SetRuleID(ruleID string) *GetRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetRuleOptions) SetXCorrelationID(xCorrelationID string) *GetRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *GetRuleOptions) SetTransactionID(transactionID string) *GetRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRuleOptions) SetHeaders(param map[string]string) *GetRuleOptions {
	options.Headers = param
	return options
}

// GetServicerefTargetOptions : The GetServicerefTarget options.
type GetServicerefTargetOptions struct {
	// The name of a service.
	ServiceName *string `json:"service_name" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetServicerefTargetOptions : Instantiate GetServicerefTargetOptions
func (*ContextBasedRestrictionsV1) NewGetServicerefTargetOptions(serviceName string) *GetServicerefTargetOptions {
	return &GetServicerefTargetOptions{
		ServiceName: core.StringPtr(serviceName),
	}
}

// SetServiceName : Allow user to set ServiceName
func (_options *GetServicerefTargetOptions) SetServiceName(serviceName string) *GetServicerefTargetOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetServicerefTargetOptions) SetXCorrelationID(xCorrelationID string) *GetServicerefTargetOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *GetServicerefTargetOptions) SetTransactionID(transactionID string) *GetServicerefTargetOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetServicerefTargetOptions) SetHeaders(param map[string]string) *GetServicerefTargetOptions {
	options.Headers = param
	return options
}

// GetZoneOptions : The GetZone options.
type GetZoneOptions struct {
	// The ID of a zone.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetZoneOptions : Instantiate GetZoneOptions
func (*ContextBasedRestrictionsV1) NewGetZoneOptions(zoneID string) *GetZoneOptions {
	return &GetZoneOptions{
		ZoneID: core.StringPtr(zoneID),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *GetZoneOptions) SetZoneID(zoneID string) *GetZoneOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *GetZoneOptions) SetXCorrelationID(xCorrelationID string) *GetZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *GetZoneOptions) SetTransactionID(transactionID string) *GetZoneOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetZoneOptions) SetHeaders(param map[string]string) *GetZoneOptions {
	options.Headers = param
	return options
}

// ListAvailableServiceOperationsOptions : The ListAvailableServiceOperations options.
type ListAvailableServiceOperationsOptions struct {
	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// The name of the service.
	ServiceName *string `json:"service_name,omitempty"`

	// The id of the service group.
	ServiceGroupID *string `json:"service_group_id,omitempty"`

	// The type of resource.
	ResourceType *string `json:"resource_type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAvailableServiceOperationsOptions : Instantiate ListAvailableServiceOperationsOptions
func (*ContextBasedRestrictionsV1) NewListAvailableServiceOperationsOptions() *ListAvailableServiceOperationsOptions {
	return &ListAvailableServiceOperationsOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListAvailableServiceOperationsOptions) SetXCorrelationID(xCorrelationID string) *ListAvailableServiceOperationsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ListAvailableServiceOperationsOptions) SetTransactionID(transactionID string) *ListAvailableServiceOperationsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListAvailableServiceOperationsOptions) SetServiceName(serviceName string) *ListAvailableServiceOperationsOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetServiceGroupID : Allow user to set ServiceGroupID
func (_options *ListAvailableServiceOperationsOptions) SetServiceGroupID(serviceGroupID string) *ListAvailableServiceOperationsOptions {
	_options.ServiceGroupID = core.StringPtr(serviceGroupID)
	return _options
}

// SetResourceType : Allow user to set ResourceType
func (_options *ListAvailableServiceOperationsOptions) SetResourceType(resourceType string) *ListAvailableServiceOperationsOptions {
	_options.ResourceType = core.StringPtr(resourceType)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAvailableServiceOperationsOptions) SetHeaders(param map[string]string) *ListAvailableServiceOperationsOptions {
	options.Headers = param
	return options
}

// ListAvailableServicerefTargetsOptions : The ListAvailableServicerefTargets options.
type ListAvailableServicerefTargetsOptions struct {
	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Specifies the types of services to retrieve.
	Type *string `json:"type,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListAvailableServicerefTargetsOptions.Type property.
// Specifies the types of services to retrieve.
const (
	ListAvailableServicerefTargetsOptionsTypeAllConst             = "all"
	ListAvailableServicerefTargetsOptionsTypePlatformServiceConst = "platform_service"
)

// NewListAvailableServicerefTargetsOptions : Instantiate ListAvailableServicerefTargetsOptions
func (*ContextBasedRestrictionsV1) NewListAvailableServicerefTargetsOptions() *ListAvailableServicerefTargetsOptions {
	return &ListAvailableServicerefTargetsOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListAvailableServicerefTargetsOptions) SetXCorrelationID(xCorrelationID string) *ListAvailableServicerefTargetsOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ListAvailableServicerefTargetsOptions) SetTransactionID(transactionID string) *ListAvailableServicerefTargetsOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListAvailableServicerefTargetsOptions) SetType(typeVar string) *ListAvailableServicerefTargetsOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAvailableServicerefTargetsOptions) SetHeaders(param map[string]string) *ListAvailableServicerefTargetsOptions {
	options.Headers = param
	return options
}

// ListRulesOptions : The ListRules options.
type ListRulesOptions struct {
	// The ID of the managing account.
	AccountID *string `json:"account_id" validate:"required"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// The `region` resource attribute.
	Region *string `json:"region,omitempty"`

	// The `resource` resource attribute.
	Resource *string `json:"resource,omitempty"`

	// The `resourceType` resource attribute.
	ResourceType *string `json:"resource_type,omitempty"`

	// The `serviceInstance` resource attribute.
	ServiceInstance *string `json:"service_instance,omitempty"`

	// The `serviceName` resource attribute.
	ServiceName *string `json:"service_name,omitempty"`

	// The rule's `serviceType` resource attribute.
	ServiceType *string `json:"service_type,omitempty"`

	// The rule's `service_group_id` resource attribute.
	ServiceGroupID *string `json:"service_group_id,omitempty"`

	// The globally unique ID of the zone.
	ZoneID *string `json:"zone_id,omitempty"`

	// Sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// The rule's `enforcement_mode` attribute.
	EnforcementMode *string `json:"enforcement_mode,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListRulesOptions.EnforcementMode property.
// The rule's `enforcement_mode` attribute.
const (
	ListRulesOptionsEnforcementModeDisabledConst = "disabled"
	ListRulesOptionsEnforcementModeEnabledConst  = "enabled"
	ListRulesOptionsEnforcementModeReportConst   = "report"
)

// NewListRulesOptions : Instantiate ListRulesOptions
func (*ContextBasedRestrictionsV1) NewListRulesOptions(accountID string) *ListRulesOptions {
	return &ListRulesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListRulesOptions) SetAccountID(accountID string) *ListRulesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListRulesOptions) SetXCorrelationID(xCorrelationID string) *ListRulesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ListRulesOptions) SetTransactionID(transactionID string) *ListRulesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetRegion : Allow user to set Region
func (_options *ListRulesOptions) SetRegion(region string) *ListRulesOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetResource : Allow user to set Resource
func (_options *ListRulesOptions) SetResource(resource string) *ListRulesOptions {
	_options.Resource = core.StringPtr(resource)
	return _options
}

// SetResourceType : Allow user to set ResourceType
func (_options *ListRulesOptions) SetResourceType(resourceType string) *ListRulesOptions {
	_options.ResourceType = core.StringPtr(resourceType)
	return _options
}

// SetServiceInstance : Allow user to set ServiceInstance
func (_options *ListRulesOptions) SetServiceInstance(serviceInstance string) *ListRulesOptions {
	_options.ServiceInstance = core.StringPtr(serviceInstance)
	return _options
}

// SetServiceName : Allow user to set ServiceName
func (_options *ListRulesOptions) SetServiceName(serviceName string) *ListRulesOptions {
	_options.ServiceName = core.StringPtr(serviceName)
	return _options
}

// SetServiceType : Allow user to set ServiceType
func (_options *ListRulesOptions) SetServiceType(serviceType string) *ListRulesOptions {
	_options.ServiceType = core.StringPtr(serviceType)
	return _options
}

// SetServiceGroupID : Allow user to set ServiceGroupID
func (_options *ListRulesOptions) SetServiceGroupID(serviceGroupID string) *ListRulesOptions {
	_options.ServiceGroupID = core.StringPtr(serviceGroupID)
	return _options
}

// SetZoneID : Allow user to set ZoneID
func (_options *ListRulesOptions) SetZoneID(zoneID string) *ListRulesOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListRulesOptions) SetSort(sort string) *ListRulesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetEnforcementMode : Allow user to set EnforcementMode
func (_options *ListRulesOptions) SetEnforcementMode(enforcementMode string) *ListRulesOptions {
	_options.EnforcementMode = core.StringPtr(enforcementMode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListRulesOptions) SetHeaders(param map[string]string) *ListRulesOptions {
	options.Headers = param
	return options
}

// ListZonesOptions : The ListZones options.
type ListZonesOptions struct {
	// The ID of the managing account.
	AccountID *string `json:"account_id" validate:"required"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// The name of the zone.
	Name *string `json:"name,omitempty"`

	// Sorts results by using a valid sort field. To learn more, see
	// [Sorting](https://cloud.ibm.com/docs/api-handbook?topic=api-handbook-sorting).
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListZonesOptions : Instantiate ListZonesOptions
func (*ContextBasedRestrictionsV1) NewListZonesOptions(accountID string) *ListZonesOptions {
	return &ListZonesOptions{
		AccountID: core.StringPtr(accountID),
	}
}

// SetAccountID : Allow user to set AccountID
func (_options *ListZonesOptions) SetAccountID(accountID string) *ListZonesOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ListZonesOptions) SetXCorrelationID(xCorrelationID string) *ListZonesOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ListZonesOptions) SetTransactionID(transactionID string) *ListZonesOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListZonesOptions) SetName(name string) *ListZonesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListZonesOptions) SetSort(sort string) *ListZonesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListZonesOptions) SetHeaders(param map[string]string) *ListZonesOptions {
	options.Headers = param
	return options
}

// NewRuleOperations : The operations this rule applies to.
type NewRuleOperations struct {
	// The API types this rule applies to.
	APITypes []NewRuleOperationsAPITypesItem `json:"api_types" validate:"required"`
}

// NewNewRuleOperations : Instantiate NewRuleOperations (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewNewRuleOperations(apiTypes []NewRuleOperationsAPITypesItem) (_model *NewRuleOperations, err error) {
	_model = &NewRuleOperations{
		APITypes: apiTypes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalNewRuleOperations unmarshals an instance of NewRuleOperations from the specified map of raw messages.
func UnmarshalNewRuleOperations(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NewRuleOperations)
	err = core.UnmarshalModel(m, "api_types", &obj.APITypes, UnmarshalNewRuleOperationsAPITypesItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_types-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// NewRuleOperationsAPITypesItem : NewRuleOperationsAPITypesItem struct
type NewRuleOperationsAPITypesItem struct {
	APITypeID *string `json:"api_type_id" validate:"required"`
}

// NewNewRuleOperationsAPITypesItem : Instantiate NewRuleOperationsAPITypesItem (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewNewRuleOperationsAPITypesItem(apiTypeID string) (_model *NewRuleOperationsAPITypesItem, err error) {
	_model = &NewRuleOperationsAPITypesItem{
		APITypeID: core.StringPtr(apiTypeID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalNewRuleOperationsAPITypesItem unmarshals an instance of NewRuleOperationsAPITypesItem from the specified map of raw messages.
func UnmarshalNewRuleOperationsAPITypesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(NewRuleOperationsAPITypesItem)
	err = core.UnmarshalPrimitive(m, "api_type_id", &obj.APITypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_type_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OperationsList : The response object of the `list_available_service_operations` operation.
type OperationsList struct {
	// The returned API types.
	APITypes []APIType `json:"api_types" validate:"required"`
}

// UnmarshalOperationsList unmarshals an instance of OperationsList from the specified map of raw messages.
func UnmarshalOperationsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OperationsList)
	err = core.UnmarshalModel(m, "api_types", &obj.APITypes, UnmarshalAPIType)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_types-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceRuleOptions : The ReplaceRule options.
type ReplaceRuleOptions struct {
	// The ID of a rule.
	RuleID *string `json:"rule_id" validate:"required,ne="`

	// The current revision of the resource being updated. This can be found in the Create/Get/Update resource response
	// ETag header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The description of the rule.
	Description *string `json:"description,omitempty"`

	// The contexts this rule applies to.
	Contexts []RuleContext `json:"contexts,omitempty"`

	// The resources this rule apply to.
	Resources []Resource `json:"resources,omitempty"`

	// The operations this rule applies to.
	Operations *NewRuleOperations `json:"operations,omitempty"`

	// The rule enforcement mode:
	//  * `enabled` - The restrictions are enforced and reported. This is the default.
	//  * `disabled` - The restrictions are disabled. Nothing is enforced or reported.
	//  * `report` - The restrictions are evaluated and reported, but not enforced.
	EnforcementMode *string `json:"enforcement_mode,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceRuleOptions.EnforcementMode property.
// The rule enforcement mode:
//   - `enabled` - The restrictions are enforced and reported. This is the default.
//   - `disabled` - The restrictions are disabled. Nothing is enforced or reported.
//   - `report` - The restrictions are evaluated and reported, but not enforced.
const (
	ReplaceRuleOptionsEnforcementModeDisabledConst = "disabled"
	ReplaceRuleOptionsEnforcementModeEnabledConst  = "enabled"
	ReplaceRuleOptionsEnforcementModeReportConst   = "report"
)

// NewReplaceRuleOptions : Instantiate ReplaceRuleOptions
func (*ContextBasedRestrictionsV1) NewReplaceRuleOptions(ruleID string, ifMatch string) *ReplaceRuleOptions {
	return &ReplaceRuleOptions{
		RuleID:  core.StringPtr(ruleID),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetRuleID : Allow user to set RuleID
func (_options *ReplaceRuleOptions) SetRuleID(ruleID string) *ReplaceRuleOptions {
	_options.RuleID = core.StringPtr(ruleID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceRuleOptions) SetIfMatch(ifMatch string) *ReplaceRuleOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceRuleOptions) SetDescription(description string) *ReplaceRuleOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetContexts : Allow user to set Contexts
func (_options *ReplaceRuleOptions) SetContexts(contexts []RuleContext) *ReplaceRuleOptions {
	_options.Contexts = contexts
	return _options
}

// SetResources : Allow user to set Resources
func (_options *ReplaceRuleOptions) SetResources(resources []Resource) *ReplaceRuleOptions {
	_options.Resources = resources
	return _options
}

// SetOperations : Allow user to set Operations
func (_options *ReplaceRuleOptions) SetOperations(operations *NewRuleOperations) *ReplaceRuleOptions {
	_options.Operations = operations
	return _options
}

// SetEnforcementMode : Allow user to set EnforcementMode
func (_options *ReplaceRuleOptions) SetEnforcementMode(enforcementMode string) *ReplaceRuleOptions {
	_options.EnforcementMode = core.StringPtr(enforcementMode)
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ReplaceRuleOptions) SetXCorrelationID(xCorrelationID string) *ReplaceRuleOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ReplaceRuleOptions) SetTransactionID(transactionID string) *ReplaceRuleOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceRuleOptions) SetHeaders(param map[string]string) *ReplaceRuleOptions {
	options.Headers = param
	return options
}

// ReplaceZoneOptions : The ReplaceZone options.
type ReplaceZoneOptions struct {
	// The ID of a zone.
	ZoneID *string `json:"zone_id" validate:"required,ne="`

	// The current revision of the resource being updated. This can be found in the Create/Get/Update resource response
	// ETag header.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The name of the zone.
	Name *string `json:"name,omitempty"`

	// The id of the account owning this zone.
	AccountID *string `json:"account_id,omitempty"`

	// The description of the zone.
	Description *string `json:"description,omitempty"`

	// The list of addresses in the zone.
	Addresses []AddressIntf `json:"addresses,omitempty"`

	// The list of excluded addresses in the zone. Only addresses of type `ipAddress`, `ipRange`, and `subnet` can be
	// excluded.
	Excluded []AddressIntf `json:"excluded,omitempty"`

	// The supplied or generated value of this header is logged for a request and repeated in a response header for the
	// corresponding response. The same value is used for downstream requests and retries of those requests. If a value of
	// this headers is not supplied in a request, the service generates a random (version 4) UUID.
	XCorrelationID *string `json:"X-Correlation-Id,omitempty"`

	// The `Transaction-Id` header behaves as the `X-Correlation-Id` header. It is supported for backward compatibility
	// with other IBM platform services that support the `Transaction-Id` header only. If both `X-Correlation-Id` and
	// `Transaction-Id` are provided, `X-Correlation-Id` has the precedence over `Transaction-Id`.
	// Deprecated: this field is deprecated and may be removed in a future release.
	TransactionID *string `json:"Transaction-Id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceZoneOptions : Instantiate ReplaceZoneOptions
func (*ContextBasedRestrictionsV1) NewReplaceZoneOptions(zoneID string, ifMatch string) *ReplaceZoneOptions {
	return &ReplaceZoneOptions{
		ZoneID:  core.StringPtr(zoneID),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetZoneID : Allow user to set ZoneID
func (_options *ReplaceZoneOptions) SetZoneID(zoneID string) *ReplaceZoneOptions {
	_options.ZoneID = core.StringPtr(zoneID)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceZoneOptions) SetIfMatch(ifMatch string) *ReplaceZoneOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceZoneOptions) SetName(name string) *ReplaceZoneOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetAccountID : Allow user to set AccountID
func (_options *ReplaceZoneOptions) SetAccountID(accountID string) *ReplaceZoneOptions {
	_options.AccountID = core.StringPtr(accountID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *ReplaceZoneOptions) SetDescription(description string) *ReplaceZoneOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetAddresses : Allow user to set Addresses
func (_options *ReplaceZoneOptions) SetAddresses(addresses []AddressIntf) *ReplaceZoneOptions {
	_options.Addresses = addresses
	return _options
}

// SetExcluded : Allow user to set Excluded
func (_options *ReplaceZoneOptions) SetExcluded(excluded []AddressIntf) *ReplaceZoneOptions {
	_options.Excluded = excluded
	return _options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (_options *ReplaceZoneOptions) SetXCorrelationID(xCorrelationID string) *ReplaceZoneOptions {
	_options.XCorrelationID = core.StringPtr(xCorrelationID)
	return _options
}

// SetTransactionID : Allow user to set TransactionID
// Deprecated: this method is deprecated and may be removed in a future release.
func (_options *ReplaceZoneOptions) SetTransactionID(transactionID string) *ReplaceZoneOptions {
	_options.TransactionID = core.StringPtr(transactionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceZoneOptions) SetHeaders(param map[string]string) *ReplaceZoneOptions {
	options.Headers = param
	return options
}

// Resource : An rule resource.
type Resource struct {
	// The resource attributes.
	Attributes []ResourceAttribute `json:"attributes" validate:"required"`

	// The optional resource tags.
	Tags []ResourceTagAttribute `json:"tags,omitempty"`
}

// NewResource : Instantiate Resource (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewResource(attributes []ResourceAttribute) (_model *Resource, err error) {
	_model = &Resource{
		Attributes: attributes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResource unmarshals an instance of Resource from the specified map of raw messages.
func UnmarshalResource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Resource)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalResourceAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "tags", &obj.Tags, UnmarshalResourceTagAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceAttribute : A rule resource attribute.
type ResourceAttribute struct {
	// The attribute name.
	Name *string `json:"name" validate:"required"`

	// The attribute value.
	Value *string `json:"value" validate:"required"`

	// The attribute operator.
	Operator *string `json:"operator,omitempty"`
}

// NewResourceAttribute : Instantiate ResourceAttribute (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewResourceAttribute(name string, value string) (_model *ResourceAttribute, err error) {
	_model = &ResourceAttribute{
		Name:  core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResourceAttribute unmarshals an instance of ResourceAttribute from the specified map of raw messages.
func UnmarshalResourceAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceAttribute)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResourceTagAttribute : A rule resource tag attribute.
type ResourceTagAttribute struct {
	// The tag attribute name.
	Name *string `json:"name" validate:"required"`

	// The tag attribute value.
	Value *string `json:"value" validate:"required"`

	// The attribute operator.
	Operator *string `json:"operator,omitempty"`
}

// NewResourceTagAttribute : Instantiate ResourceTagAttribute (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewResourceTagAttribute(name string, value string) (_model *ResourceTagAttribute, err error) {
	_model = &ResourceTagAttribute{
		Name:  core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalResourceTagAttribute unmarshals an instance of ResourceTagAttribute from the specified map of raw messages.
func UnmarshalResourceTagAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceTagAttribute)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Rule : An output rule.
type Rule struct {
	// The globally unique ID of the rule.
	ID *string `json:"id" validate:"required"`

	// The rule CRN.
	CRN *string `json:"crn" validate:"required"`

	// The description of the rule.
	Description *string `json:"description" validate:"required"`

	// The contexts this rule applies to.
	Contexts []RuleContext `json:"contexts" validate:"required"`

	// The resources this rule apply to.
	Resources []Resource `json:"resources" validate:"required"`

	// The operations this rule applies to.
	Operations *NewRuleOperations `json:"operations,omitempty"`

	// The rule enforcement mode:
	//  * `enabled` - The restrictions are enforced and reported. This is the default.
	//  * `disabled` - The restrictions are disabled. Nothing is enforced or reported.
	//  * `report` - The restrictions are evaluated and reported, but not enforced.
	EnforcementMode *string `json:"enforcement_mode,omitempty"`

	// The href link to the resource.
	Href *string `json:"href" validate:"required"`

	// The time the resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// IAM ID of the user or service which created the resource.
	CreatedByID *string `json:"created_by_id" validate:"required"`

	// The last time the resource was modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at" validate:"required"`

	// IAM ID of the user or service which modified the resource.
	LastModifiedByID *string `json:"last_modified_by_id" validate:"required"`
}

// Constants associated with the Rule.EnforcementMode property.
// The rule enforcement mode:
//   - `enabled` - The restrictions are enforced and reported. This is the default.
//   - `disabled` - The restrictions are disabled. Nothing is enforced or reported.
//   - `report` - The restrictions are evaluated and reported, but not enforced.
const (
	RuleEnforcementModeDisabledConst = "disabled"
	RuleEnforcementModeEnabledConst  = "enabled"
	RuleEnforcementModeReportConst   = "report"
)

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "contexts", &obj.Contexts, UnmarshalRuleContext)
	if err != nil {
		err = core.SDKErrorf(err, "", "contexts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resources", &obj.Resources, UnmarshalResource)
	if err != nil {
		err = core.SDKErrorf(err, "", "resources-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "operations", &obj.Operations, UnmarshalNewRuleOperations)
	if err != nil {
		err = core.SDKErrorf(err, "", "operations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enforcement_mode", &obj.EnforcementMode)
	if err != nil {
		err = core.SDKErrorf(err, "", "enforcement_mode-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleContext : A rule context.
type RuleContext struct {
	// The attributes.
	Attributes []RuleContextAttribute `json:"attributes" validate:"required"`
}

// NewRuleContext : Instantiate RuleContext (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewRuleContext(attributes []RuleContextAttribute) (_model *RuleContext, err error) {
	_model = &RuleContext{
		Attributes: attributes,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRuleContext unmarshals an instance of RuleContext from the specified map of raw messages.
func UnmarshalRuleContext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleContext)
	err = core.UnmarshalModel(m, "attributes", &obj.Attributes, UnmarshalRuleContextAttribute)
	if err != nil {
		err = core.SDKErrorf(err, "", "attributes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleContextAttribute : An rule context attribute.
type RuleContextAttribute struct {
	// The attribute name.
	Name *string `json:"name" validate:"required"`

	// The attribute value.
	Value *string `json:"value" validate:"required"`
}

// NewRuleContextAttribute : Instantiate RuleContextAttribute (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewRuleContextAttribute(name string, value string) (_model *RuleContextAttribute, err error) {
	_model = &RuleContextAttribute{
		Name:  core.StringPtr(name),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRuleContextAttribute unmarshals an instance of RuleContextAttribute from the specified map of raw messages.
func UnmarshalRuleContextAttribute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleContextAttribute)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RuleList : The response object of the ListRules operation.
type RuleList struct {
	// The number of returned results.
	Count *int64 `json:"count" validate:"required"`

	// The returned rules.
	Rules []Rule `json:"rules" validate:"required"`
}

// UnmarshalRuleList unmarshals an instance of RuleList from the specified map of raw messages.
func UnmarshalRuleList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RuleList)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceRefTarget : Summary information about a service reference target.
type ServiceRefTarget struct {
	// The name of the service.
	ServiceName *string `json:"service_name" validate:"required"`

	// The type of the service.
	ServiceType *string `json:"service_type,omitempty"`

	// The locations the service is available.
	Locations []ServiceRefTargetLocationsItem `json:"locations,omitempty"`
}

// UnmarshalServiceRefTarget unmarshals an instance of ServiceRefTarget from the specified map of raw messages.
func UnmarshalServiceRefTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceRefTarget)
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_type", &obj.ServiceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "locations", &obj.Locations, UnmarshalServiceRefTargetLocationsItem)
	if err != nil {
		err = core.SDKErrorf(err, "", "locations-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceRefTargetList : A list of service reference targets.
type ServiceRefTargetList struct {
	// The number of returned results.
	Count *int64 `json:"count" validate:"required"`

	// The list of service reference targets.
	Targets []ServiceRefTarget `json:"targets" validate:"required"`
}

// UnmarshalServiceRefTargetList unmarshals an instance of ServiceRefTargetList from the specified map of raw messages.
func UnmarshalServiceRefTargetList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceRefTargetList)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "targets", &obj.Targets, UnmarshalServiceRefTarget)
	if err != nil {
		err = core.SDKErrorf(err, "", "targets-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceRefTargetLocationsItem : ServiceRefTargetLocationsItem struct
type ServiceRefTargetLocationsItem struct {
	// The location display name.
	DisplayName *string `json:"display_name,omitempty"`

	// The location kind.
	Kind *string `json:"kind,omitempty"`

	// The location name.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalServiceRefTargetLocationsItem unmarshals an instance of ServiceRefTargetLocationsItem from the specified map of raw messages.
func UnmarshalServiceRefTargetLocationsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceRefTargetLocationsItem)
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "kind", &obj.Kind)
	if err != nil {
		err = core.SDKErrorf(err, "", "kind-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceRefValue : A service reference value.
type ServiceRefValue struct {
	// The id of the account owning the service.
	AccountID *string `json:"account_id" validate:"required"`

	// The service type.
	ServiceType *string `json:"service_type,omitempty"`

	// The service name.
	ServiceName *string `json:"service_name,omitempty"`

	// The service instance.
	ServiceInstance *string `json:"service_instance,omitempty"`

	// The location.
	Location *string `json:"location,omitempty"`
}

// NewServiceRefValue : Instantiate ServiceRefValue (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewServiceRefValue(accountID string) (_model *ServiceRefValue, err error) {
	_model = &ServiceRefValue{
		AccountID: core.StringPtr(accountID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalServiceRefValue unmarshals an instance of ServiceRefValue from the specified map of raw messages.
func UnmarshalServiceRefValue(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceRefValue)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_type", &obj.ServiceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance", &obj.ServiceInstance)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Zone : An output zone.
type Zone struct {
	// The globally unique ID of the zone.
	ID *string `json:"id" validate:"required"`

	// The zone CRN.
	CRN *string `json:"crn" validate:"required"`

	// The number of addresses in the zone.
	AddressCount *int64 `json:"address_count" validate:"required"`

	// The number of excluded addresses in the zone.
	ExcludedCount *int64 `json:"excluded_count" validate:"required"`

	// The name of the zone.
	Name *string `json:"name" validate:"required"`

	// The id of the account owning this zone.
	AccountID *string `json:"account_id" validate:"required"`

	// The description of the zone.
	Description *string `json:"description" validate:"required"`

	// The list of addresses in the zone.
	Addresses []AddressIntf `json:"addresses" validate:"required"`

	// The list of excluded addresses in the zone. Only addresses of type `ipAddress`, `ipRange`, and `subnet` can be
	// excluded.
	Excluded []AddressIntf `json:"excluded" validate:"required"`

	// The href link to the resource.
	Href *string `json:"href" validate:"required"`

	// The time the resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// IAM ID of the user or service which created the resource.
	CreatedByID *string `json:"created_by_id" validate:"required"`

	// The last time the resource was modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at" validate:"required"`

	// IAM ID of the user or service which modified the resource.
	LastModifiedByID *string `json:"last_modified_by_id" validate:"required"`
}

// UnmarshalZone unmarshals an instance of Zone from the specified map of raw messages.
func UnmarshalZone(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Zone)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "address_count", &obj.AddressCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "address_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "excluded_count", &obj.ExcludedCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "excluded_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "addresses", &obj.Addresses, UnmarshalAddress)
	if err != nil {
		err = core.SDKErrorf(err, "", "addresses-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "excluded", &obj.Excluded, UnmarshalAddress)
	if err != nil {
		err = core.SDKErrorf(err, "", "excluded-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ZoneList : The response object of the ListZones operation.
type ZoneList struct {
	// The number of returned results.
	Count *int64 `json:"count" validate:"required"`

	// The returned zones.
	Zones []ZoneSummary `json:"zones" validate:"required"`
}

// UnmarshalZoneList unmarshals an instance of ZoneList from the specified map of raw messages.
func UnmarshalZoneList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneList)
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		err = core.SDKErrorf(err, "", "count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "zones", &obj.Zones, UnmarshalZoneSummary)
	if err != nil {
		err = core.SDKErrorf(err, "", "zones-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ZoneSummary : An output zone summary.
type ZoneSummary struct {
	// The globally unique ID of the zone.
	ID *string `json:"id" validate:"required"`

	// The zone CRN.
	CRN *string `json:"crn" validate:"required"`

	// The name of the zone.
	Name *string `json:"name" validate:"required"`

	// The description of the zone.
	Description *string `json:"description,omitempty"`

	// A preview of addresses in the zone (3 addresses maximum).
	AddressesPreview []AddressIntf `json:"addresses_preview" validate:"required"`

	// The number of addresses in the zone.
	AddressCount *int64 `json:"address_count" validate:"required"`

	// The number of excluded addresses in the zone.
	ExcludedCount *int64 `json:"excluded_count" validate:"required"`

	// The href link to the resource.
	Href *string `json:"href" validate:"required"`

	// The time the resource was created.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// IAM ID of the user or service which created the resource.
	CreatedByID *string `json:"created_by_id" validate:"required"`

	// The last time the resource was modified.
	LastModifiedAt *strfmt.DateTime `json:"last_modified_at" validate:"required"`

	// IAM ID of the user or service which modified the resource.
	LastModifiedByID *string `json:"last_modified_by_id" validate:"required"`
}

// UnmarshalZoneSummary unmarshals an instance of ZoneSummary from the specified map of raw messages.
func UnmarshalZoneSummary(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ZoneSummary)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "addresses_preview", &obj.AddressesPreview, UnmarshalAddress)
	if err != nil {
		err = core.SDKErrorf(err, "", "addresses_preview-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "address_count", &obj.AddressCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "address_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "excluded_count", &obj.ExcludedCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "excluded_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by_id", &obj.CreatedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_at", &obj.LastModifiedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_modified_by_id", &obj.LastModifiedByID)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_modified_by_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AddressIPAddress : A single IP address. IPv4 and IPv6 are supported.
// This model "extends" Address
type AddressIPAddress struct {
	// The type of address.
	Type *string `json:"type" validate:"required"`

	// The IP address.
	Value *string `json:"value" validate:"required"`

	// The address id (for use by terraform only).
	ID *string `json:"id,omitempty"`
}

// Constants associated with the AddressIPAddress.Type property.
// The type of address.
const (
	AddressIPAddressTypeIpaddressConst = "ipAddress"
)

// NewAddressIPAddress : Instantiate AddressIPAddress (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewAddressIPAddress(typeVar string, value string) (_model *AddressIPAddress, err error) {
	_model = &AddressIPAddress{
		Type:  core.StringPtr(typeVar),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*AddressIPAddress) isaAddress() bool {
	return true
}

// UnmarshalAddressIPAddress unmarshals an instance of AddressIPAddress from the specified map of raw messages.
func UnmarshalAddressIPAddress(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AddressIPAddress)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AddressIPAddressRange : An IP address range. IPv4 and IPv6 are supported.
// This model "extends" Address
type AddressIPAddressRange struct {
	// The type of address.
	Type *string `json:"type" validate:"required"`

	// The ip range in <first-ip>-<last-ip> format.
	Value *string `json:"value" validate:"required"`

	// The address id (for use by terraform only).
	ID *string `json:"id,omitempty"`
}

// Constants associated with the AddressIPAddressRange.Type property.
// The type of address.
const (
	AddressIPAddressRangeTypeIprangeConst = "ipRange"
)

// NewAddressIPAddressRange : Instantiate AddressIPAddressRange (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewAddressIPAddressRange(typeVar string, value string) (_model *AddressIPAddressRange, err error) {
	_model = &AddressIPAddressRange{
		Type:  core.StringPtr(typeVar),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*AddressIPAddressRange) isaAddress() bool {
	return true
}

// UnmarshalAddressIPAddressRange unmarshals an instance of AddressIPAddressRange from the specified map of raw messages.
func UnmarshalAddressIPAddressRange(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AddressIPAddressRange)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AddressServiceRef : A service reference.
// This model "extends" Address
type AddressServiceRef struct {
	// The type of address.
	Type *string `json:"type" validate:"required"`

	// A service reference value.
	Ref *ServiceRefValue `json:"ref" validate:"required"`

	// The address id (for use by terraform only).
	ID *string `json:"id,omitempty"`
}

// Constants associated with the AddressServiceRef.Type property.
// The type of address.
const (
	AddressServiceRefTypeServicerefConst = "serviceRef"
)

// NewAddressServiceRef : Instantiate AddressServiceRef (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewAddressServiceRef(typeVar string, ref *ServiceRefValue) (_model *AddressServiceRef, err error) {
	_model = &AddressServiceRef{
		Type: core.StringPtr(typeVar),
		Ref:  ref,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*AddressServiceRef) isaAddress() bool {
	return true
}

// UnmarshalAddressServiceRef unmarshals an instance of AddressServiceRef from the specified map of raw messages.
func UnmarshalAddressServiceRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AddressServiceRef)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ref", &obj.Ref, UnmarshalServiceRefValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "ref-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AddressSubnet : A subnet in CIDR format.
// This model "extends" Address
type AddressSubnet struct {
	// The type of address.
	Type *string `json:"type" validate:"required"`

	// The subnet in CIDR format.
	Value *string `json:"value" validate:"required"`

	// The address id (for use by terraform only).
	ID *string `json:"id,omitempty"`
}

// Constants associated with the AddressSubnet.Type property.
// The type of address.
const (
	AddressSubnetTypeSubnetConst = "subnet"
)

// NewAddressSubnet : Instantiate AddressSubnet (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewAddressSubnet(typeVar string, value string) (_model *AddressSubnet, err error) {
	_model = &AddressSubnet{
		Type:  core.StringPtr(typeVar),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*AddressSubnet) isaAddress() bool {
	return true
}

// UnmarshalAddressSubnet unmarshals an instance of AddressSubnet from the specified map of raw messages.
func UnmarshalAddressSubnet(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AddressSubnet)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AddressVPC : A single VPC address.
// This model "extends" Address
type AddressVPC struct {
	// The type of address.
	Type *string `json:"type" validate:"required"`

	// The VPC CRN.
	Value *string `json:"value" validate:"required"`

	// The address id (for use by terraform only).
	ID *string `json:"id,omitempty"`
}

// Constants associated with the AddressVPC.Type property.
// The type of address.
const (
	AddressVPCTypeVPCConst = "vpc"
)

// NewAddressVPC : Instantiate AddressVPC (Generic Model Constructor)
func (*ContextBasedRestrictionsV1) NewAddressVPC(typeVar string, value string) (_model *AddressVPC, err error) {
	_model = &AddressVPC{
		Type:  core.StringPtr(typeVar),
		Value: core.StringPtr(value),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*AddressVPC) isaAddress() bool {
	return true
}

// UnmarshalAddressVPC unmarshals an instance of AddressVPC from the specified map of raw messages.
func UnmarshalAddressVPC(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AddressVPC)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
