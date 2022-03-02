/**
 * (C) Copyright IBM Corp. 2021.
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
 * IBM OpenAPI SDK Code Generator Version: 3.36.1-694fc13e-20210723-211159
 */

// Package atrackerv1 : Operations and models for the AtrackerV1 service
package atrackerv1

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

// AtrackerV1 : Activity Tracker is a platform service that you can configure in each region in your account to define
// how auditing events are collected and stored. Events are stored in a Cloud Object Storage bucket that is also
// available in the account.
//
// Version: 1.1.0
type AtrackerV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.atracker.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "atracker"

// AtrackerV1Options : Service options
type AtrackerV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewAtrackerV1UsingExternalConfig : constructs an instance of AtrackerV1 with passed in options and external configuration.
func NewAtrackerV1UsingExternalConfig(options *AtrackerV1Options) (atracker *AtrackerV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	atracker, err = NewAtrackerV1(options)
	if err != nil {
		return
	}

	err = atracker.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = atracker.Service.SetServiceURL(options.URL)
	}
	return
}

// NewAtrackerV1 : constructs an instance of AtrackerV1 with passed in options.
func NewAtrackerV1(options *AtrackerV1Options) (service *AtrackerV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
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

	service = &AtrackerV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://us-south.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the us-south region.
		"private.us-south": "https://private.us-south.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the us-south region.
		"us-east": "https://us-east.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the us-east region.
		"private.us-east": "https://private.us-east.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the us-east region.
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "atracker" suitable for processing requests.
func (atracker *AtrackerV1) Clone() *AtrackerV1 {
	if core.IsNil(atracker) {
		return nil
	}
	clone := *atracker
	clone.Service = atracker.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (atracker *AtrackerV1) SetServiceURL(url string) error {
	return atracker.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (atracker *AtrackerV1) GetServiceURL() string {
	return atracker.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (atracker *AtrackerV1) SetDefaultHeaders(headers http.Header) {
	atracker.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (atracker *AtrackerV1) SetEnableGzipCompression(enableGzip bool) {
	atracker.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (atracker *AtrackerV1) GetEnableGzipCompression() bool {
	return atracker.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (atracker *AtrackerV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	atracker.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (atracker *AtrackerV1) DisableRetries() {
	atracker.Service.DisableRetries()
}

// CreateTarget : Create a target
// Creates a Cloud Object Storage (COS) target that includes information about the COS endpoint and the credentials to
// access the bucket. You must define a COS target per region.  Notice that although you can use the same COS bucket for
// collecting auditing events in your account across multiple regions, you should consider defining a bucket in each
// region to reduce performance and network latency issues. You can define up to 16 targets per region.
func (atracker *AtrackerV1) CreateTarget(createTargetOptions *CreateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return atracker.CreateTargetWithContext(context.Background(), createTargetOptions)
}

// CreateTargetWithContext is an alternate form of the CreateTarget method which supports a Context parameter
func (atracker *AtrackerV1) CreateTargetWithContext(ctx context.Context, createTargetOptions *CreateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTargetOptions, "createTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTargetOptions, "createTargetOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/targets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "CreateTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTargetOptions.Name != nil {
		body["name"] = createTargetOptions.Name
	}
	if createTargetOptions.TargetType != nil {
		body["target_type"] = createTargetOptions.TargetType
	}
	if createTargetOptions.CosEndpoint != nil {
		body["cos_endpoint"] = createTargetOptions.CosEndpoint
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListTargets : List targets
// List all Cloud Object Storage (COS) targets that are defined in a region.
func (atracker *AtrackerV1) ListTargets(listTargetsOptions *ListTargetsOptions) (result *TargetList, response *core.DetailedResponse, err error) {
	return atracker.ListTargetsWithContext(context.Background(), listTargetsOptions)
}

// ListTargetsWithContext is an alternate form of the ListTargets method which supports a Context parameter
func (atracker *AtrackerV1) ListTargetsWithContext(ctx context.Context, listTargetsOptions *ListTargetsOptions) (result *TargetList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listTargetsOptions, "listTargetsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/targets`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listTargetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "ListTargets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTargetList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTarget : Get details of a target
// Retrieve the configuration details of a target.
func (atracker *AtrackerV1) GetTarget(getTargetOptions *GetTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return atracker.GetTargetWithContext(context.Background(), getTargetOptions)
}

// GetTargetWithContext is an alternate form of the GetTarget method which supports a Context parameter
func (atracker *AtrackerV1) GetTargetWithContext(ctx context.Context, getTargetOptions *GetTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTargetOptions, "getTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTargetOptions, "getTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/targets/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "GetTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTarget : Update a target
// Update the configuration details of a target.
func (atracker *AtrackerV1) ReplaceTarget(replaceTargetOptions *ReplaceTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return atracker.ReplaceTargetWithContext(context.Background(), replaceTargetOptions)
}

// ReplaceTargetWithContext is an alternate form of the ReplaceTarget method which supports a Context parameter
func (atracker *AtrackerV1) ReplaceTargetWithContext(ctx context.Context, replaceTargetOptions *ReplaceTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTargetOptions, "replaceTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceTargetOptions, "replaceTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *replaceTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/targets/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "ReplaceTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceTargetOptions.Name != nil {
		body["name"] = replaceTargetOptions.Name
	}
	if replaceTargetOptions.TargetType != nil {
		body["target_type"] = replaceTargetOptions.TargetType
	}
	if replaceTargetOptions.CosEndpoint != nil {
		body["cos_endpoint"] = replaceTargetOptions.CosEndpoint
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTarget : Delete a target
// Delete a target.
func (atracker *AtrackerV1) DeleteTarget(deleteTargetOptions *DeleteTargetOptions) (result *WarningReport, response *core.DetailedResponse, err error) {
	return atracker.DeleteTargetWithContext(context.Background(), deleteTargetOptions)
}

// DeleteTargetWithContext is an alternate form of the DeleteTarget method which supports a Context parameter
func (atracker *AtrackerV1) DeleteTargetWithContext(ctx context.Context, deleteTargetOptions *DeleteTargetOptions) (result *WarningReport, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTargetOptions, "deleteTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTargetOptions, "deleteTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/targets/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "DeleteTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWarningReport)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ValidateTarget : Validate a target
// Validate a target by checking the credentials to write to the bucket. The result is included as additional data of
// the target in the section "cos_write_status".
func (atracker *AtrackerV1) ValidateTarget(validateTargetOptions *ValidateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	return atracker.ValidateTargetWithContext(context.Background(), validateTargetOptions)
}

// ValidateTargetWithContext is an alternate form of the ValidateTarget method which supports a Context parameter
func (atracker *AtrackerV1) ValidateTargetWithContext(ctx context.Context, validateTargetOptions *ValidateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(validateTargetOptions, "validateTargetOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(validateTargetOptions, "validateTargetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *validateTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/targets/{id}/validate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range validateTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "ValidateTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateRoute : Create a route
// Create a route to define the rule that specifies how to manage auditing events in a region.  You can define 1 route
// only per region. You can configure 1 target only per route. To define how to manage global events, that is, auditing
// events in your account that are not region specific, you must configure 1 route in your account to collect and route
// global events. You must set the receive_global_events field to true.
func (atracker *AtrackerV1) CreateRoute(createRouteOptions *CreateRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	return atracker.CreateRouteWithContext(context.Background(), createRouteOptions)
}

// CreateRouteWithContext is an alternate form of the CreateRoute method which supports a Context parameter
func (atracker *AtrackerV1) CreateRouteWithContext(ctx context.Context, createRouteOptions *CreateRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRouteOptions, "createRouteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRouteOptions, "createRouteOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/routes`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "CreateRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRouteOptions.Name != nil {
		body["name"] = createRouteOptions.Name
	}
	if createRouteOptions.ReceiveGlobalEvents != nil {
		body["receive_global_events"] = createRouteOptions.ReceiveGlobalEvents
	}
	if createRouteOptions.Rules != nil {
		body["rules"] = createRouteOptions.Rules
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoute)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListRoutes : List routes
// List the route that is configured in a region.
func (atracker *AtrackerV1) ListRoutes(listRoutesOptions *ListRoutesOptions) (result *RouteList, response *core.DetailedResponse, err error) {
	return atracker.ListRoutesWithContext(context.Background(), listRoutesOptions)
}

// ListRoutesWithContext is an alternate form of the ListRoutes method which supports a Context parameter
func (atracker *AtrackerV1) ListRoutesWithContext(ctx context.Context, listRoutesOptions *ListRoutesOptions) (result *RouteList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRoutesOptions, "listRoutesOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/routes`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRoutesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "ListRoutes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRouteList)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetRoute : Get details of a route
// Get the configuration details of a route.
func (atracker *AtrackerV1) GetRoute(getRouteOptions *GetRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	return atracker.GetRouteWithContext(context.Background(), getRouteOptions)
}

// GetRouteWithContext is an alternate form of the GetRoute method which supports a Context parameter
func (atracker *AtrackerV1) GetRouteWithContext(ctx context.Context, getRouteOptions *GetRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRouteOptions, "getRouteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRouteOptions, "getRouteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *getRouteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/routes/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "GetRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoute)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceRoute : Update a route
// Update the configuration details of a route.
func (atracker *AtrackerV1) ReplaceRoute(replaceRouteOptions *ReplaceRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	return atracker.ReplaceRouteWithContext(context.Background(), replaceRouteOptions)
}

// ReplaceRouteWithContext is an alternate form of the ReplaceRoute method which supports a Context parameter
func (atracker *AtrackerV1) ReplaceRouteWithContext(ctx context.Context, replaceRouteOptions *ReplaceRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRouteOptions, "replaceRouteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceRouteOptions, "replaceRouteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *replaceRouteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/routes/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range replaceRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "ReplaceRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceRouteOptions.Name != nil {
		body["name"] = replaceRouteOptions.Name
	}
	if replaceRouteOptions.ReceiveGlobalEvents != nil {
		body["receive_global_events"] = replaceRouteOptions.ReceiveGlobalEvents
	}
	if replaceRouteOptions.Rules != nil {
		body["rules"] = replaceRouteOptions.Rules
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoute)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteRoute : Delete a route
// Deletes a route.
func (atracker *AtrackerV1) DeleteRoute(deleteRouteOptions *DeleteRouteOptions) (response *core.DetailedResponse, err error) {
	return atracker.DeleteRouteWithContext(context.Background(), deleteRouteOptions)
}

// DeleteRouteWithContext is an alternate form of the DeleteRoute method which supports a Context parameter
func (atracker *AtrackerV1) DeleteRouteWithContext(ctx context.Context, deleteRouteOptions *DeleteRouteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRouteOptions, "deleteRouteOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRouteOptions, "deleteRouteOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteRouteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/routes/{id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "DeleteRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = atracker.Service.Request(request, nil)

	return
}

// GetEndpoints : Get endpoints
// Get information about the public and private endpoints that are enabled in a region when you use the Activity Tracker
// API.
func (atracker *AtrackerV1) GetEndpoints(getEndpointsOptions *GetEndpointsOptions) (result *Endpoints, response *core.DetailedResponse, err error) {
	return atracker.GetEndpointsWithContext(context.Background(), getEndpointsOptions)
}

// GetEndpointsWithContext is an alternate form of the GetEndpoints method which supports a Context parameter
func (atracker *AtrackerV1) GetEndpointsWithContext(ctx context.Context, getEndpointsOptions *GetEndpointsOptions) (result *Endpoints, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getEndpointsOptions, "getEndpointsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/endpoints`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "GetEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoints)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// PatchEndpoints : Modify endpoints
// Configure the public endpoint availability in a region to use the Activity Tracker API. By default, the private
// endpoint is enabled and cannot be disabled.
func (atracker *AtrackerV1) PatchEndpoints(patchEndpointsOptions *PatchEndpointsOptions) (result *Endpoints, response *core.DetailedResponse, err error) {
	return atracker.PatchEndpointsWithContext(context.Background(), patchEndpointsOptions)
}

// PatchEndpointsWithContext is an alternate form of the PatchEndpoints method which supports a Context parameter
func (atracker *AtrackerV1) PatchEndpointsWithContext(ctx context.Context, patchEndpointsOptions *PatchEndpointsOptions) (result *Endpoints, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(patchEndpointsOptions, "patchEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(patchEndpointsOptions, "patchEndpointsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v1/endpoints`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range patchEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V1", "PatchEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if patchEndpointsOptions.APIEndpoint != nil {
		body["api_endpoint"] = patchEndpointsOptions.APIEndpoint
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoints)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// APIEndpoint : Activity Tracker API endpoint.
type APIEndpoint struct {
	// The public URL of Activity Tracker in a region.
	PublicURL *string `json:"public_url" validate:"required"`

	// Indicates whether or not the public endpoint is enabled in the account.
	PublicEnabled *bool `json:"public_enabled" validate:"required"`

	// The private URL of Activity Tracker. This URL cannot be disabled.
	PrivateURL *string `json:"private_url" validate:"required"`

	// The private endpoint is always enabled.
	PrivateEnabled *bool `json:"private_enabled,omitempty"`
}

// UnmarshalAPIEndpoint unmarshals an instance of APIEndpoint from the specified map of raw messages.
func UnmarshalAPIEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(APIEndpoint)
	err = core.UnmarshalPrimitive(m, "public_url", &obj.PublicURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "public_enabled", &obj.PublicEnabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_url", &obj.PrivateURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "private_enabled", &obj.PrivateEnabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRouteOptions : The CreateRoute options.
type CreateRouteOptions struct {
	// The name of the route. The name must be 1000 characters or less and cannot include any special characters other than
	// `(space) - . _ :`.
	Name *string `json:"name" validate:"required"`

	// Indicates whether or not all global events should be forwarded to this region.
	ReceiveGlobalEvents *bool `json:"receive_global_events" validate:"required"`

	// Routing rules that will be evaluated in their order of the array.
	Rules []Rule `json:"rules" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateRouteOptions : Instantiate CreateRouteOptions
func (*AtrackerV1) NewCreateRouteOptions(name string, receiveGlobalEvents bool, rules []Rule) *CreateRouteOptions {
	return &CreateRouteOptions{
		Name: core.StringPtr(name),
		ReceiveGlobalEvents: core.BoolPtr(receiveGlobalEvents),
		Rules: rules,
	}
}

// SetName : Allow user to set Name
func (_options *CreateRouteOptions) SetName(name string) *CreateRouteOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetReceiveGlobalEvents : Allow user to set ReceiveGlobalEvents
func (_options *CreateRouteOptions) SetReceiveGlobalEvents(receiveGlobalEvents bool) *CreateRouteOptions {
	_options.ReceiveGlobalEvents = core.BoolPtr(receiveGlobalEvents)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *CreateRouteOptions) SetRules(rules []Rule) *CreateRouteOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRouteOptions) SetHeaders(param map[string]string) *CreateRouteOptions {
	options.Headers = param
	return options
}

// CreateTargetOptions : The CreateTarget options.
type CreateTargetOptions struct {
	// The name of the target. The name must be 1000 characters or less, and cannot include any special characters other
	// than `(space) - . _ :`.
	Name *string `json:"name" validate:"required"`

	// The type of the target.
	TargetType *string `json:"target_type" validate:"required"`

	// Property values for a Cloud Object Storage Endpoint.
	CosEndpoint *CosEndpoint `json:"cos_endpoint" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateTargetOptions.TargetType property.
// The type of the target.
const (
	CreateTargetOptionsTargetTypeCloudObjectStorageConst = "cloud_object_storage"
)

// NewCreateTargetOptions : Instantiate CreateTargetOptions
func (*AtrackerV1) NewCreateTargetOptions(name string, targetType string, cosEndpoint *CosEndpoint) *CreateTargetOptions {
	return &CreateTargetOptions{
		Name: core.StringPtr(name),
		TargetType: core.StringPtr(targetType),
		CosEndpoint: cosEndpoint,
	}
}

// SetName : Allow user to set Name
func (_options *CreateTargetOptions) SetName(name string) *CreateTargetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetTargetType : Allow user to set TargetType
func (_options *CreateTargetOptions) SetTargetType(targetType string) *CreateTargetOptions {
	_options.TargetType = core.StringPtr(targetType)
	return _options
}

// SetCosEndpoint : Allow user to set CosEndpoint
func (_options *CreateTargetOptions) SetCosEndpoint(cosEndpoint *CosEndpoint) *CreateTargetOptions {
	_options.CosEndpoint = cosEndpoint
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTargetOptions) SetHeaders(param map[string]string) *CreateTargetOptions {
	options.Headers = param
	return options
}

// DeleteRouteOptions : The DeleteRoute options.
type DeleteRouteOptions struct {
	// The v4 UUID that uniquely identifies the route.
	ID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRouteOptions : Instantiate DeleteRouteOptions
func (*AtrackerV1) NewDeleteRouteOptions(id string) *DeleteRouteOptions {
	return &DeleteRouteOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteRouteOptions) SetID(id string) *DeleteRouteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRouteOptions) SetHeaders(param map[string]string) *DeleteRouteOptions {
	options.Headers = param
	return options
}

// DeleteTargetOptions : The DeleteTarget options.
type DeleteTargetOptions struct {
	// The v4 UUID that uniquely identifies the target.
	ID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteTargetOptions : Instantiate DeleteTargetOptions
func (*AtrackerV1) NewDeleteTargetOptions(id string) *DeleteTargetOptions {
	return &DeleteTargetOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteTargetOptions) SetID(id string) *DeleteTargetOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTargetOptions) SetHeaders(param map[string]string) *DeleteTargetOptions {
	options.Headers = param
	return options
}

// Endpoints : Activity Tracker endpoints.
type Endpoints struct {
	// Activity Tracker API endpoint.
	APIEndpoint *APIEndpoint `json:"api_endpoint" validate:"required"`
}

// UnmarshalEndpoints unmarshals an instance of Endpoints from the specified map of raw messages.
func UnmarshalEndpoints(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Endpoints)
	err = core.UnmarshalModel(m, "api_endpoint", &obj.APIEndpoint, UnmarshalAPIEndpoint)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointsRequestAPIEndpoint : Activity Tracker service API endpoint.
type EndpointsRequestAPIEndpoint struct {
	// Indicate whether or not the public endpoint is enabled in an account.
	PublicEnabled *bool `json:"public_enabled,omitempty"`
}

// UnmarshalEndpointsRequestAPIEndpoint unmarshals an instance of EndpointsRequestAPIEndpoint from the specified map of raw messages.
func UnmarshalEndpointsRequestAPIEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointsRequestAPIEndpoint)
	err = core.UnmarshalPrimitive(m, "public_enabled", &obj.PublicEnabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetEndpointsOptions : The GetEndpoints options.
type GetEndpointsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEndpointsOptions : Instantiate GetEndpointsOptions
func (*AtrackerV1) NewGetEndpointsOptions() *GetEndpointsOptions {
	return &GetEndpointsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetEndpointsOptions) SetHeaders(param map[string]string) *GetEndpointsOptions {
	options.Headers = param
	return options
}

// GetRouteOptions : The GetRoute options.
type GetRouteOptions struct {
	// The v4 UUID that uniquely identifies the route.
	ID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRouteOptions : Instantiate GetRouteOptions
func (*AtrackerV1) NewGetRouteOptions(id string) *GetRouteOptions {
	return &GetRouteOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetRouteOptions) SetID(id string) *GetRouteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetRouteOptions) SetHeaders(param map[string]string) *GetRouteOptions {
	options.Headers = param
	return options
}

// GetTargetOptions : The GetTarget options.
type GetTargetOptions struct {
	// The v4 UUID that uniquely identifies the target.
	ID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetTargetOptions : Instantiate GetTargetOptions
func (*AtrackerV1) NewGetTargetOptions(id string) *GetTargetOptions {
	return &GetTargetOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetTargetOptions) SetID(id string) *GetTargetOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTargetOptions) SetHeaders(param map[string]string) *GetTargetOptions {
	options.Headers = param
	return options
}

// ListRoutesOptions : The ListRoutes options.
type ListRoutesOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListRoutesOptions : Instantiate ListRoutesOptions
func (*AtrackerV1) NewListRoutesOptions() *ListRoutesOptions {
	return &ListRoutesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListRoutesOptions) SetHeaders(param map[string]string) *ListRoutesOptions {
	options.Headers = param
	return options
}

// ListTargetsOptions : The ListTargets options.
type ListTargetsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListTargetsOptions : Instantiate ListTargetsOptions
func (*AtrackerV1) NewListTargetsOptions() *ListTargetsOptions {
	return &ListTargetsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListTargetsOptions) SetHeaders(param map[string]string) *ListTargetsOptions {
	options.Headers = param
	return options
}

// PatchEndpointsOptions : The PatchEndpoints options.
type PatchEndpointsOptions struct {
	// Activity Tracker service API endpoint.
	APIEndpoint *EndpointsRequestAPIEndpoint `json:"api_endpoint,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewPatchEndpointsOptions : Instantiate PatchEndpointsOptions
func (*AtrackerV1) NewPatchEndpointsOptions() *PatchEndpointsOptions {
	return &PatchEndpointsOptions{}
}

// SetAPIEndpoint : Allow user to set APIEndpoint
func (_options *PatchEndpointsOptions) SetAPIEndpoint(apiEndpoint *EndpointsRequestAPIEndpoint) *PatchEndpointsOptions {
	_options.APIEndpoint = apiEndpoint
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PatchEndpointsOptions) SetHeaders(param map[string]string) *PatchEndpointsOptions {
	options.Headers = param
	return options
}

// ReplaceRouteOptions : The ReplaceRoute options.
type ReplaceRouteOptions struct {
	// The v4 UUID that uniquely identifies the route.
	ID *string `json:"-" validate:"required,ne="`

	// The name of the route. The name must be 1000 characters or less and cannot include any special characters other than
	// `(space) - . _ :`.
	Name *string `json:"name" validate:"required"`

	// Indicates whether or not all global events should be forwarded to this region.
	ReceiveGlobalEvents *bool `json:"receive_global_events" validate:"required"`

	// Routing rules that will be evaluated in their order of the array.
	Rules []Rule `json:"rules" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceRouteOptions : Instantiate ReplaceRouteOptions
func (*AtrackerV1) NewReplaceRouteOptions(id string, name string, receiveGlobalEvents bool, rules []Rule) *ReplaceRouteOptions {
	return &ReplaceRouteOptions{
		ID: core.StringPtr(id),
		Name: core.StringPtr(name),
		ReceiveGlobalEvents: core.BoolPtr(receiveGlobalEvents),
		Rules: rules,
	}
}

// SetID : Allow user to set ID
func (_options *ReplaceRouteOptions) SetID(id string) *ReplaceRouteOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceRouteOptions) SetName(name string) *ReplaceRouteOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetReceiveGlobalEvents : Allow user to set ReceiveGlobalEvents
func (_options *ReplaceRouteOptions) SetReceiveGlobalEvents(receiveGlobalEvents bool) *ReplaceRouteOptions {
	_options.ReceiveGlobalEvents = core.BoolPtr(receiveGlobalEvents)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *ReplaceRouteOptions) SetRules(rules []Rule) *ReplaceRouteOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceRouteOptions) SetHeaders(param map[string]string) *ReplaceRouteOptions {
	options.Headers = param
	return options
}

// ReplaceTargetOptions : The ReplaceTarget options.
type ReplaceTargetOptions struct {
	// The v4 UUID that uniquely identifies the target.
	ID *string `json:"-" validate:"required,ne="`

	// The name of the target. The name must be 1000 characters or less, and cannot include any special characters other
	// than `(space) - . _ :`.
	Name *string `json:"name" validate:"required"`

	// The type of the target.
	TargetType *string `json:"target_type" validate:"required"`

	// Property values for a Cloud Object Storage Endpoint.
	CosEndpoint *CosEndpoint `json:"cos_endpoint" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceTargetOptions.TargetType property.
// The type of the target.
const (
	ReplaceTargetOptionsTargetTypeCloudObjectStorageConst = "cloud_object_storage"
)

// NewReplaceTargetOptions : Instantiate ReplaceTargetOptions
func (*AtrackerV1) NewReplaceTargetOptions(id string, name string, targetType string, cosEndpoint *CosEndpoint) *ReplaceTargetOptions {
	return &ReplaceTargetOptions{
		ID: core.StringPtr(id),
		Name: core.StringPtr(name),
		TargetType: core.StringPtr(targetType),
		CosEndpoint: cosEndpoint,
	}
}

// SetID : Allow user to set ID
func (_options *ReplaceTargetOptions) SetID(id string) *ReplaceTargetOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceTargetOptions) SetName(name string) *ReplaceTargetOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetTargetType : Allow user to set TargetType
func (_options *ReplaceTargetOptions) SetTargetType(targetType string) *ReplaceTargetOptions {
	_options.TargetType = core.StringPtr(targetType)
	return _options
}

// SetCosEndpoint : Allow user to set CosEndpoint
func (_options *ReplaceTargetOptions) SetCosEndpoint(cosEndpoint *CosEndpoint) *ReplaceTargetOptions {
	_options.CosEndpoint = cosEndpoint
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTargetOptions) SetHeaders(param map[string]string) *ReplaceTargetOptions {
	options.Headers = param
	return options
}

// Route : The route resource.
type Route struct {
	// The uuid of the route resource.
	ID *string `json:"id" validate:"required"`

	// The name of the route.
	Name *string `json:"name" validate:"required"`

	// The crn of the route resource.
	CRN *string `json:"crn" validate:"required"`

	// The version of the route.
	Version *int64 `json:"version,omitempty"`

	// Indicates whether or not all global events should be forwarded to this region.
	ReceiveGlobalEvents *bool `json:"receive_global_events" validate:"required"`

	// The routing rules that will be evaluated in their order of the array.
	Rules []Rule `json:"rules" validate:"required"`

	// The timestamp of the route creation time.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The timestamp of the route last updated time.
	Updated *strfmt.DateTime `json:"updated,omitempty"`
}

// UnmarshalRoute unmarshals an instance of Route from the specified map of raw messages.
func UnmarshalRoute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Route)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "receive_global_events", &obj.ReceiveGlobalEvents)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RouteList : A list of route resources.
type RouteList struct {
	// A list of route resources.
	Routes []Route `json:"routes" validate:"required"`
}

// UnmarshalRouteList unmarshals an instance of RouteList from the specified map of raw messages.
func UnmarshalRouteList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RouteList)
	err = core.UnmarshalModel(m, "routes", &obj.Routes, UnmarshalRoute)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Rule : The request payload to create a regional route.
type Rule struct {
	// The target ID List. Only 1 target id is supported.
	TargetIds []string `json:"target_ids" validate:"required"`
}

// NewRule : Instantiate Rule (Generic Model Constructor)
func (*AtrackerV1) NewRule(targetIds []string) (_model *Rule, err error) {
	_model = &Rule{
		TargetIds: targetIds,
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "target_ids", &obj.TargetIds)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Target : Property values for a target in the response. Credentials associated with the target are encrypted and masked as
// REDACTED in the response.
type Target struct {
	// The uuid of the target resource.
	ID *string `json:"id" validate:"required"`

	// The name of the target resource.
	Name *string `json:"name" validate:"required"`

	// The crn of the target resource.
	CRN *string `json:"crn" validate:"required"`

	// The type of the target.
	TargetType *string `json:"target_type" validate:"required"`

	// The encryption key that is used to encrypt events before Activity Tracker services buffer them on storage. This
	// credential is masked in the response.
	EncryptKey *string `json:"encrypt_key,omitempty"`

	// Property values for a Cloud Object Storage Endpoint.
	CosEndpoint *CosEndpoint `json:"cos_endpoint,omitempty"`

	// The status of the write attempt with the provided cos_endpoint parameters.
	CosWriteStatus *CosWriteStatus `json:"cos_write_status,omitempty"`

	// The timestamp of the target creation time.
	Created *strfmt.DateTime `json:"created,omitempty"`

	// The timestamp of the target last updated time.
	Updated *strfmt.DateTime `json:"updated,omitempty"`
}

// Constants associated with the Target.TargetType property.
// The type of the target.
const (
	TargetTargetTypeCloudObjectStorageConst = "cloud_object_storage"
)

// UnmarshalTarget unmarshals an instance of Target from the specified map of raw messages.
func UnmarshalTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Target)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_type", &obj.TargetType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "encrypt_key", &obj.EncryptKey)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cos_endpoint", &obj.CosEndpoint, UnmarshalCosEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "cos_write_status", &obj.CosWriteStatus, UnmarshalCosWriteStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated", &obj.Updated)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TargetList : A list of target resources.
type TargetList struct {
	// A list of target resources.
	Targets []Target `json:"targets" validate:"required"`
}

// UnmarshalTargetList unmarshals an instance of TargetList from the specified map of raw messages.
func UnmarshalTargetList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetList)
	err = core.UnmarshalModel(m, "targets", &obj.Targets, UnmarshalTarget)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ValidateTargetOptions : The ValidateTarget options.
type ValidateTargetOptions struct {
	// The v4 UUID that uniquely identifies the target.
	ID *string `json:"-" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewValidateTargetOptions : Instantiate ValidateTargetOptions
func (*AtrackerV1) NewValidateTargetOptions(id string) *ValidateTargetOptions {
	return &ValidateTargetOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *ValidateTargetOptions) SetID(id string) *ValidateTargetOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ValidateTargetOptions) SetHeaders(param map[string]string) *ValidateTargetOptions {
	options.Headers = param
	return options
}

// Warning : The warning object.
type Warning struct {
	// The warning code.
	Code *string `json:"code,omitempty"`

	// The warning message.
	Message *string `json:"message,omitempty"`
}

// UnmarshalWarning unmarshals an instance of Warning from the specified map of raw messages.
func UnmarshalWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Warning)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WarningReport : Description of a warning that occurred in a service request.
type WarningReport struct {
	// The status code.
	StatusCode *int64 `json:"status_code,omitempty"`

	// The transaction-id of the API request.
	Trace *string `json:"trace,omitempty"`

	// The warning array triggered by the API request.
	Warnings []Warning `json:"warnings,omitempty"`
}

// UnmarshalWarningReport unmarshals an instance of WarningReport from the specified map of raw messages.
func UnmarshalWarningReport(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WarningReport)
	err = core.UnmarshalPrimitive(m, "status_code", &obj.StatusCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "warnings", &obj.Warnings, UnmarshalWarning)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CosEndpoint : Property values for a Cloud Object Storage Endpoint.
type CosEndpoint struct {
	// The host name of the Cloud Object Storage endpoint.
	Endpoint *string `json:"endpoint" validate:"required"`

	// The CRN of the Cloud Object Storage instance.
	TargetCRN *string `json:"target_crn" validate:"required"`

	// The bucket name under the Cloud Object Storage instance.
	Bucket *string `json:"bucket" validate:"required"`

	// The IAM API key that has writer access to the Cloud Object Storage instance. This credential is masked in the
	// response.
	APIKey *string `json:"api_key" validate:"required"`
}

// NewCosEndpoint : Instantiate CosEndpoint (Generic Model Constructor)
func (*AtrackerV1) NewCosEndpoint(endpoint string, targetCRN string, bucket string, apiKey string) (_model *CosEndpoint, err error) {
	_model = &CosEndpoint{
		Endpoint: core.StringPtr(endpoint),
		TargetCRN: core.StringPtr(targetCRN),
		Bucket: core.StringPtr(bucket),
		APIKey: core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalCosEndpoint unmarshals an instance of CosEndpoint from the specified map of raw messages.
func UnmarshalCosEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CosEndpoint)
	err = core.UnmarshalPrimitive(m, "endpoint", &obj.Endpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CosWriteStatus : The status of the write attempt with the provided cos_endpoint parameters.
type CosWriteStatus struct {
	// The status such as failed or success.
	Status *string `json:"status,omitempty"`

	// The timestamp of the failure.
	LastFailure *strfmt.DateTime `json:"last_failure,omitempty"`

	// Detailed description of the cause of the failure.
	ReasonForLastFailure *string `json:"reason_for_last_failure,omitempty"`
}

// UnmarshalCosWriteStatus unmarshals an instance of CosWriteStatus from the specified map of raw messages.
func UnmarshalCosWriteStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CosWriteStatus)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_failure", &obj.LastFailure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reason_for_last_failure", &obj.ReasonForLastFailure)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
