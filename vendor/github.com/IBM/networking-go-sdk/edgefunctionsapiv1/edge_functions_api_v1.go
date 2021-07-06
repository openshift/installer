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
 

// Package edgefunctionsapiv1 : Operations and models for the EdgeFunctionsApiV1 service
package edgefunctionsapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

// EdgeFunctionsApiV1 : Edge Functions
//
// Version: 1.0.0
type EdgeFunctionsApiV1 struct {
	Service *core.BaseService

	// cloud resource name.
	Crn *string

	// zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "edge_functions_api"

// EdgeFunctionsApiV1Options : Service options
type EdgeFunctionsApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// cloud resource name.
	Crn *string `validate:"required"`

	// zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewEdgeFunctionsApiV1UsingExternalConfig : constructs an instance of EdgeFunctionsApiV1 with passed in options and external configuration.
func NewEdgeFunctionsApiV1UsingExternalConfig(options *EdgeFunctionsApiV1Options) (edgeFunctionsApi *EdgeFunctionsApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	edgeFunctionsApi, err = NewEdgeFunctionsApiV1(options)
	if err != nil {
		return
	}

	err = edgeFunctionsApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = edgeFunctionsApi.Service.SetServiceURL(options.URL)
	}
	return
}

// NewEdgeFunctionsApiV1 : constructs an instance of EdgeFunctionsApiV1 with passed in options.
func NewEdgeFunctionsApiV1(options *EdgeFunctionsApiV1Options) (service *EdgeFunctionsApiV1, err error) {
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

	service = &EdgeFunctionsApiV1{
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

// Clone makes a copy of "edgeFunctionsApi" suitable for processing requests.
func (edgeFunctionsApi *EdgeFunctionsApiV1) Clone() *EdgeFunctionsApiV1 {
	if core.IsNil(edgeFunctionsApi) {
		return nil
	}
	clone := *edgeFunctionsApi
	clone.Service = edgeFunctionsApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (edgeFunctionsApi *EdgeFunctionsApiV1) SetServiceURL(url string) error {
	return edgeFunctionsApi.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (edgeFunctionsApi *EdgeFunctionsApiV1) GetServiceURL() string {
	return edgeFunctionsApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (edgeFunctionsApi *EdgeFunctionsApiV1) SetDefaultHeaders(headers http.Header) {
	edgeFunctionsApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (edgeFunctionsApi *EdgeFunctionsApiV1) SetEnableGzipCompression(enableGzip bool) {
	edgeFunctionsApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (edgeFunctionsApi *EdgeFunctionsApiV1) GetEnableGzipCompression() bool {
	return edgeFunctionsApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (edgeFunctionsApi *EdgeFunctionsApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	edgeFunctionsApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (edgeFunctionsApi *EdgeFunctionsApiV1) DisableRetries() {
	edgeFunctionsApi.Service.DisableRetries()
}

// ListEdgeFunctionsActions : Get all edge functions scripts for a given instance
// Get all edge functions scripts for a given instance.
func (edgeFunctionsApi *EdgeFunctionsApiV1) ListEdgeFunctionsActions(listEdgeFunctionsActionsOptions *ListEdgeFunctionsActionsOptions) (result *ListEdgeFunctionsActionsResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.ListEdgeFunctionsActionsWithContext(context.Background(), listEdgeFunctionsActionsOptions)
}

// ListEdgeFunctionsActionsWithContext is an alternate form of the ListEdgeFunctionsActions method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) ListEdgeFunctionsActionsWithContext(ctx context.Context, listEdgeFunctionsActionsOptions *ListEdgeFunctionsActionsOptions) (result *ListEdgeFunctionsActionsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listEdgeFunctionsActionsOptions, "listEdgeFunctionsActionsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/workers/scripts`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listEdgeFunctionsActionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "ListEdgeFunctionsActions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listEdgeFunctionsActionsOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listEdgeFunctionsActionsOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListEdgeFunctionsActionsResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateEdgeFunctionsAction : Upload or replace an edge functions action for a given instance
// Upload or replace an exitsing edge functions action for a given instance.
func (edgeFunctionsApi *EdgeFunctionsApiV1) UpdateEdgeFunctionsAction(updateEdgeFunctionsActionOptions *UpdateEdgeFunctionsActionOptions) (result *GetEdgeFunctionsActionResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.UpdateEdgeFunctionsActionWithContext(context.Background(), updateEdgeFunctionsActionOptions)
}

// UpdateEdgeFunctionsActionWithContext is an alternate form of the UpdateEdgeFunctionsAction method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) UpdateEdgeFunctionsActionWithContext(ctx context.Context, updateEdgeFunctionsActionOptions *UpdateEdgeFunctionsActionOptions) (result *GetEdgeFunctionsActionResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEdgeFunctionsActionOptions, "updateEdgeFunctionsActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateEdgeFunctionsActionOptions, "updateEdgeFunctionsActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"script_name": *updateEdgeFunctionsActionOptions.ScriptName,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/workers/scripts/{script_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateEdgeFunctionsActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "UpdateEdgeFunctionsAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/javascript")
	if updateEdgeFunctionsActionOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateEdgeFunctionsActionOptions.XCorrelationID))
	}

	_, err = builder.SetBodyContent("application/javascript", nil, nil, updateEdgeFunctionsActionOptions.EdgeFunctionsAction)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetEdgeFunctionsActionResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetEdgeFunctionsAction : Download a edge functions action for a given instance
// Fetch raw script content for your worker. Note this is the original script content, not JSON encoded.
func (edgeFunctionsApi *EdgeFunctionsApiV1) GetEdgeFunctionsAction(getEdgeFunctionsActionOptions *GetEdgeFunctionsActionOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.GetEdgeFunctionsActionWithContext(context.Background(), getEdgeFunctionsActionOptions)
}

// GetEdgeFunctionsActionWithContext is an alternate form of the GetEdgeFunctionsAction method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) GetEdgeFunctionsActionWithContext(ctx context.Context, getEdgeFunctionsActionOptions *GetEdgeFunctionsActionOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEdgeFunctionsActionOptions, "getEdgeFunctionsActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEdgeFunctionsActionOptions, "getEdgeFunctionsActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"script_name": *getEdgeFunctionsActionOptions.ScriptName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/workers/scripts/{script_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEdgeFunctionsActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "GetEdgeFunctionsAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/javascript")
	if getEdgeFunctionsActionOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getEdgeFunctionsActionOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = edgeFunctionsApi.Service.Request(request, &result)

	return
}

// DeleteEdgeFunctionsAction : Delete a edge functions action for a given instance
// Delete an edge functions action for a given instance.
func (edgeFunctionsApi *EdgeFunctionsApiV1) DeleteEdgeFunctionsAction(deleteEdgeFunctionsActionOptions *DeleteEdgeFunctionsActionOptions) (result *DeleteEdgeFunctionsActionResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.DeleteEdgeFunctionsActionWithContext(context.Background(), deleteEdgeFunctionsActionOptions)
}

// DeleteEdgeFunctionsActionWithContext is an alternate form of the DeleteEdgeFunctionsAction method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) DeleteEdgeFunctionsActionWithContext(ctx context.Context, deleteEdgeFunctionsActionOptions *DeleteEdgeFunctionsActionOptions) (result *DeleteEdgeFunctionsActionResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEdgeFunctionsActionOptions, "deleteEdgeFunctionsActionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEdgeFunctionsActionOptions, "deleteEdgeFunctionsActionOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"script_name": *deleteEdgeFunctionsActionOptions.ScriptName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/workers/scripts/{script_name}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteEdgeFunctionsActionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "DeleteEdgeFunctionsAction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteEdgeFunctionsActionOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteEdgeFunctionsActionOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteEdgeFunctionsActionResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateEdgeFunctionsTrigger : Create an edge functions trigger on a given zone
// Create an edge functions trigger on a given zone.
func (edgeFunctionsApi *EdgeFunctionsApiV1) CreateEdgeFunctionsTrigger(createEdgeFunctionsTriggerOptions *CreateEdgeFunctionsTriggerOptions) (result *CreateEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.CreateEdgeFunctionsTriggerWithContext(context.Background(), createEdgeFunctionsTriggerOptions)
}

// CreateEdgeFunctionsTriggerWithContext is an alternate form of the CreateEdgeFunctionsTrigger method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) CreateEdgeFunctionsTriggerWithContext(ctx context.Context, createEdgeFunctionsTriggerOptions *CreateEdgeFunctionsTriggerOptions) (result *CreateEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createEdgeFunctionsTriggerOptions, "createEdgeFunctionsTriggerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"zone_identifier": *edgeFunctionsApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/workers/routes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createEdgeFunctionsTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "CreateEdgeFunctionsTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if createEdgeFunctionsTriggerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*createEdgeFunctionsTriggerOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if createEdgeFunctionsTriggerOptions.Pattern != nil {
		body["pattern"] = createEdgeFunctionsTriggerOptions.Pattern
	}
	if createEdgeFunctionsTriggerOptions.Script != nil {
		body["script"] = createEdgeFunctionsTriggerOptions.Script
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
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListEdgeFunctionsTriggers : List all edge functions triggers on a given zone
// List all edge functions triggers on a given zone.
func (edgeFunctionsApi *EdgeFunctionsApiV1) ListEdgeFunctionsTriggers(listEdgeFunctionsTriggersOptions *ListEdgeFunctionsTriggersOptions) (result *ListEdgeFunctionsTriggersResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.ListEdgeFunctionsTriggersWithContext(context.Background(), listEdgeFunctionsTriggersOptions)
}

// ListEdgeFunctionsTriggersWithContext is an alternate form of the ListEdgeFunctionsTriggers method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) ListEdgeFunctionsTriggersWithContext(ctx context.Context, listEdgeFunctionsTriggersOptions *ListEdgeFunctionsTriggersOptions) (result *ListEdgeFunctionsTriggersResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listEdgeFunctionsTriggersOptions, "listEdgeFunctionsTriggersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"zone_identifier": *edgeFunctionsApi.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/workers/routes`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listEdgeFunctionsTriggersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "ListEdgeFunctionsTriggers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if listEdgeFunctionsTriggersOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*listEdgeFunctionsTriggersOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListEdgeFunctionsTriggersResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetEdgeFunctionsTrigger : Get an edge functions trigger on a given zone
// Get an edge functions trigger on a given zone.
func (edgeFunctionsApi *EdgeFunctionsApiV1) GetEdgeFunctionsTrigger(getEdgeFunctionsTriggerOptions *GetEdgeFunctionsTriggerOptions) (result *GetEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.GetEdgeFunctionsTriggerWithContext(context.Background(), getEdgeFunctionsTriggerOptions)
}

// GetEdgeFunctionsTriggerWithContext is an alternate form of the GetEdgeFunctionsTrigger method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) GetEdgeFunctionsTriggerWithContext(ctx context.Context, getEdgeFunctionsTriggerOptions *GetEdgeFunctionsTriggerOptions) (result *GetEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEdgeFunctionsTriggerOptions, "getEdgeFunctionsTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEdgeFunctionsTriggerOptions, "getEdgeFunctionsTriggerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"zone_identifier": *edgeFunctionsApi.ZoneIdentifier,
		"route_id": *getEdgeFunctionsTriggerOptions.RouteID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/workers/routes/{route_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEdgeFunctionsTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "GetEdgeFunctionsTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if getEdgeFunctionsTriggerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*getEdgeFunctionsTriggerOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateEdgeFunctionsTrigger : Update an edge functions trigger on a given zone
// Update an edge functions trigger on a given zone.
func (edgeFunctionsApi *EdgeFunctionsApiV1) UpdateEdgeFunctionsTrigger(updateEdgeFunctionsTriggerOptions *UpdateEdgeFunctionsTriggerOptions) (result *GetEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.UpdateEdgeFunctionsTriggerWithContext(context.Background(), updateEdgeFunctionsTriggerOptions)
}

// UpdateEdgeFunctionsTriggerWithContext is an alternate form of the UpdateEdgeFunctionsTrigger method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) UpdateEdgeFunctionsTriggerWithContext(ctx context.Context, updateEdgeFunctionsTriggerOptions *UpdateEdgeFunctionsTriggerOptions) (result *GetEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEdgeFunctionsTriggerOptions, "updateEdgeFunctionsTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateEdgeFunctionsTriggerOptions, "updateEdgeFunctionsTriggerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"zone_identifier": *edgeFunctionsApi.ZoneIdentifier,
		"route_id": *updateEdgeFunctionsTriggerOptions.RouteID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/workers/routes/{route_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateEdgeFunctionsTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "UpdateEdgeFunctionsTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if updateEdgeFunctionsTriggerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*updateEdgeFunctionsTriggerOptions.XCorrelationID))
	}

	body := make(map[string]interface{})
	if updateEdgeFunctionsTriggerOptions.Pattern != nil {
		body["pattern"] = updateEdgeFunctionsTriggerOptions.Pattern
	}
	if updateEdgeFunctionsTriggerOptions.Script != nil {
		body["script"] = updateEdgeFunctionsTriggerOptions.Script
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
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteEdgeFunctionsTrigger : Delete an edge functions trigger on a given zone
// Delete an edge functions trigger on a given zone.
func (edgeFunctionsApi *EdgeFunctionsApiV1) DeleteEdgeFunctionsTrigger(deleteEdgeFunctionsTriggerOptions *DeleteEdgeFunctionsTriggerOptions) (result *CreateEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	return edgeFunctionsApi.DeleteEdgeFunctionsTriggerWithContext(context.Background(), deleteEdgeFunctionsTriggerOptions)
}

// DeleteEdgeFunctionsTriggerWithContext is an alternate form of the DeleteEdgeFunctionsTrigger method which supports a Context parameter
func (edgeFunctionsApi *EdgeFunctionsApiV1) DeleteEdgeFunctionsTriggerWithContext(ctx context.Context, deleteEdgeFunctionsTriggerOptions *DeleteEdgeFunctionsTriggerOptions) (result *CreateEdgeFunctionsTriggerResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEdgeFunctionsTriggerOptions, "deleteEdgeFunctionsTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEdgeFunctionsTriggerOptions, "deleteEdgeFunctionsTriggerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *edgeFunctionsApi.Crn,
		"zone_identifier": *edgeFunctionsApi.ZoneIdentifier,
		"route_id": *deleteEdgeFunctionsTriggerOptions.RouteID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = edgeFunctionsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(edgeFunctionsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/workers/routes/{route_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteEdgeFunctionsTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("edge_functions_api", "V1", "DeleteEdgeFunctionsTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if deleteEdgeFunctionsTriggerOptions.XCorrelationID != nil {
		builder.AddHeader("X-Correlation-ID", fmt.Sprint(*deleteEdgeFunctionsTriggerOptions.XCorrelationID))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = edgeFunctionsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateEdgeFunctionsTriggerOptions : The CreateEdgeFunctionsTrigger options.
type CreateEdgeFunctionsTriggerOptions struct {
	// a string pattern.
	Pattern *string `json:"pattern,omitempty"`

	// Name of the script to apply when the route is matched. The route is skipped when this is blank/missing.
	Script *string `json:"script,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateEdgeFunctionsTriggerOptions : Instantiate CreateEdgeFunctionsTriggerOptions
func (*EdgeFunctionsApiV1) NewCreateEdgeFunctionsTriggerOptions() *CreateEdgeFunctionsTriggerOptions {
	return &CreateEdgeFunctionsTriggerOptions{}
}

// SetPattern : Allow user to set Pattern
func (options *CreateEdgeFunctionsTriggerOptions) SetPattern(pattern string) *CreateEdgeFunctionsTriggerOptions {
	options.Pattern = core.StringPtr(pattern)
	return options
}

// SetScript : Allow user to set Script
func (options *CreateEdgeFunctionsTriggerOptions) SetScript(script string) *CreateEdgeFunctionsTriggerOptions {
	options.Script = core.StringPtr(script)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *CreateEdgeFunctionsTriggerOptions) SetXCorrelationID(xCorrelationID string) *CreateEdgeFunctionsTriggerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateEdgeFunctionsTriggerOptions) SetHeaders(param map[string]string) *CreateEdgeFunctionsTriggerOptions {
	options.Headers = param
	return options
}

// DeleteEdgeFunctionsActionOptions : The DeleteEdgeFunctionsAction options.
type DeleteEdgeFunctionsActionOptions struct {
	// the edge function action name.
	ScriptName *string `json:"script_name" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEdgeFunctionsActionOptions : Instantiate DeleteEdgeFunctionsActionOptions
func (*EdgeFunctionsApiV1) NewDeleteEdgeFunctionsActionOptions(scriptName string) *DeleteEdgeFunctionsActionOptions {
	return &DeleteEdgeFunctionsActionOptions{
		ScriptName: core.StringPtr(scriptName),
	}
}

// SetScriptName : Allow user to set ScriptName
func (options *DeleteEdgeFunctionsActionOptions) SetScriptName(scriptName string) *DeleteEdgeFunctionsActionOptions {
	options.ScriptName = core.StringPtr(scriptName)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteEdgeFunctionsActionOptions) SetXCorrelationID(xCorrelationID string) *DeleteEdgeFunctionsActionOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEdgeFunctionsActionOptions) SetHeaders(param map[string]string) *DeleteEdgeFunctionsActionOptions {
	options.Headers = param
	return options
}

// DeleteEdgeFunctionsTriggerOptions : The DeleteEdgeFunctionsTrigger options.
type DeleteEdgeFunctionsTriggerOptions struct {
	// trigger identifier.
	RouteID *string `json:"route_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEdgeFunctionsTriggerOptions : Instantiate DeleteEdgeFunctionsTriggerOptions
func (*EdgeFunctionsApiV1) NewDeleteEdgeFunctionsTriggerOptions(routeID string) *DeleteEdgeFunctionsTriggerOptions {
	return &DeleteEdgeFunctionsTriggerOptions{
		RouteID: core.StringPtr(routeID),
	}
}

// SetRouteID : Allow user to set RouteID
func (options *DeleteEdgeFunctionsTriggerOptions) SetRouteID(routeID string) *DeleteEdgeFunctionsTriggerOptions {
	options.RouteID = core.StringPtr(routeID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *DeleteEdgeFunctionsTriggerOptions) SetXCorrelationID(xCorrelationID string) *DeleteEdgeFunctionsTriggerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEdgeFunctionsTriggerOptions) SetHeaders(param map[string]string) *DeleteEdgeFunctionsTriggerOptions {
	options.Headers = param
	return options
}

// GetEdgeFunctionsActionOptions : The GetEdgeFunctionsAction options.
type GetEdgeFunctionsActionOptions struct {
	// the edge function action name.
	ScriptName *string `json:"script_name" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEdgeFunctionsActionOptions : Instantiate GetEdgeFunctionsActionOptions
func (*EdgeFunctionsApiV1) NewGetEdgeFunctionsActionOptions(scriptName string) *GetEdgeFunctionsActionOptions {
	return &GetEdgeFunctionsActionOptions{
		ScriptName: core.StringPtr(scriptName),
	}
}

// SetScriptName : Allow user to set ScriptName
func (options *GetEdgeFunctionsActionOptions) SetScriptName(scriptName string) *GetEdgeFunctionsActionOptions {
	options.ScriptName = core.StringPtr(scriptName)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetEdgeFunctionsActionOptions) SetXCorrelationID(xCorrelationID string) *GetEdgeFunctionsActionOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetEdgeFunctionsActionOptions) SetHeaders(param map[string]string) *GetEdgeFunctionsActionOptions {
	options.Headers = param
	return options
}

// GetEdgeFunctionsTriggerOptions : The GetEdgeFunctionsTrigger options.
type GetEdgeFunctionsTriggerOptions struct {
	// trigger identifier.
	RouteID *string `json:"route_id" validate:"required,ne="`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEdgeFunctionsTriggerOptions : Instantiate GetEdgeFunctionsTriggerOptions
func (*EdgeFunctionsApiV1) NewGetEdgeFunctionsTriggerOptions(routeID string) *GetEdgeFunctionsTriggerOptions {
	return &GetEdgeFunctionsTriggerOptions{
		RouteID: core.StringPtr(routeID),
	}
}

// SetRouteID : Allow user to set RouteID
func (options *GetEdgeFunctionsTriggerOptions) SetRouteID(routeID string) *GetEdgeFunctionsTriggerOptions {
	options.RouteID = core.StringPtr(routeID)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *GetEdgeFunctionsTriggerOptions) SetXCorrelationID(xCorrelationID string) *GetEdgeFunctionsTriggerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetEdgeFunctionsTriggerOptions) SetHeaders(param map[string]string) *GetEdgeFunctionsTriggerOptions {
	options.Headers = param
	return options
}

// ListEdgeFunctionsActionsOptions : The ListEdgeFunctionsActions options.
type ListEdgeFunctionsActionsOptions struct {
	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListEdgeFunctionsActionsOptions : Instantiate ListEdgeFunctionsActionsOptions
func (*EdgeFunctionsApiV1) NewListEdgeFunctionsActionsOptions() *ListEdgeFunctionsActionsOptions {
	return &ListEdgeFunctionsActionsOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListEdgeFunctionsActionsOptions) SetXCorrelationID(xCorrelationID string) *ListEdgeFunctionsActionsOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListEdgeFunctionsActionsOptions) SetHeaders(param map[string]string) *ListEdgeFunctionsActionsOptions {
	options.Headers = param
	return options
}

// ListEdgeFunctionsTriggersOptions : The ListEdgeFunctionsTriggers options.
type ListEdgeFunctionsTriggersOptions struct {
	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListEdgeFunctionsTriggersOptions : Instantiate ListEdgeFunctionsTriggersOptions
func (*EdgeFunctionsApiV1) NewListEdgeFunctionsTriggersOptions() *ListEdgeFunctionsTriggersOptions {
	return &ListEdgeFunctionsTriggersOptions{}
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *ListEdgeFunctionsTriggersOptions) SetXCorrelationID(xCorrelationID string) *ListEdgeFunctionsTriggersOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListEdgeFunctionsTriggersOptions) SetHeaders(param map[string]string) *ListEdgeFunctionsTriggersOptions {
	options.Headers = param
	return options
}

// UpdateEdgeFunctionsActionOptions : The UpdateEdgeFunctionsAction options.
type UpdateEdgeFunctionsActionOptions struct {
	// the edge function action name.
	ScriptName *string `json:"script_name" validate:"required,ne="`

	// upload or replace an edge functions action.
	EdgeFunctionsAction io.ReadCloser `json:"edge_functions_action,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateEdgeFunctionsActionOptions : Instantiate UpdateEdgeFunctionsActionOptions
func (*EdgeFunctionsApiV1) NewUpdateEdgeFunctionsActionOptions(scriptName string) *UpdateEdgeFunctionsActionOptions {
	return &UpdateEdgeFunctionsActionOptions{
		ScriptName: core.StringPtr(scriptName),
	}
}

// SetScriptName : Allow user to set ScriptName
func (options *UpdateEdgeFunctionsActionOptions) SetScriptName(scriptName string) *UpdateEdgeFunctionsActionOptions {
	options.ScriptName = core.StringPtr(scriptName)
	return options
}

// SetEdgeFunctionsAction : Allow user to set EdgeFunctionsAction
func (options *UpdateEdgeFunctionsActionOptions) SetEdgeFunctionsAction(edgeFunctionsAction io.ReadCloser) *UpdateEdgeFunctionsActionOptions {
	options.EdgeFunctionsAction = edgeFunctionsAction
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdateEdgeFunctionsActionOptions) SetXCorrelationID(xCorrelationID string) *UpdateEdgeFunctionsActionOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEdgeFunctionsActionOptions) SetHeaders(param map[string]string) *UpdateEdgeFunctionsActionOptions {
	options.Headers = param
	return options
}

// UpdateEdgeFunctionsTriggerOptions : The UpdateEdgeFunctionsTrigger options.
type UpdateEdgeFunctionsTriggerOptions struct {
	// trigger identifier.
	RouteID *string `json:"route_id" validate:"required,ne="`

	// a string pattern.
	Pattern *string `json:"pattern,omitempty"`

	// Name of the script to apply when the route is matched. The route is skipped when this is blank/missing.
	Script *string `json:"script,omitempty"`

	// Uniquely identifying a request.
	XCorrelationID *string `json:"X-Correlation-ID,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateEdgeFunctionsTriggerOptions : Instantiate UpdateEdgeFunctionsTriggerOptions
func (*EdgeFunctionsApiV1) NewUpdateEdgeFunctionsTriggerOptions(routeID string) *UpdateEdgeFunctionsTriggerOptions {
	return &UpdateEdgeFunctionsTriggerOptions{
		RouteID: core.StringPtr(routeID),
	}
}

// SetRouteID : Allow user to set RouteID
func (options *UpdateEdgeFunctionsTriggerOptions) SetRouteID(routeID string) *UpdateEdgeFunctionsTriggerOptions {
	options.RouteID = core.StringPtr(routeID)
	return options
}

// SetPattern : Allow user to set Pattern
func (options *UpdateEdgeFunctionsTriggerOptions) SetPattern(pattern string) *UpdateEdgeFunctionsTriggerOptions {
	options.Pattern = core.StringPtr(pattern)
	return options
}

// SetScript : Allow user to set Script
func (options *UpdateEdgeFunctionsTriggerOptions) SetScript(script string) *UpdateEdgeFunctionsTriggerOptions {
	options.Script = core.StringPtr(script)
	return options
}

// SetXCorrelationID : Allow user to set XCorrelationID
func (options *UpdateEdgeFunctionsTriggerOptions) SetXCorrelationID(xCorrelationID string) *UpdateEdgeFunctionsTriggerOptions {
	options.XCorrelationID = core.StringPtr(xCorrelationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEdgeFunctionsTriggerOptions) SetHeaders(param map[string]string) *UpdateEdgeFunctionsTriggerOptions {
	options.Headers = param
	return options
}

// CreateEdgeFunctionsTriggerResp : create an edge funtions trigger response.
type CreateEdgeFunctionsTriggerResp struct {
	// edge function trigger id.
	Result *EdgeFunctionsTriggerID `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	// An array with errors.
	Errors []string `json:"errors,omitempty"`

	// An array with messages.
	Messages []string `json:"messages,omitempty"`
}


// UnmarshalCreateEdgeFunctionsTriggerResp unmarshals an instance of CreateEdgeFunctionsTriggerResp from the specified map of raw messages.
func UnmarshalCreateEdgeFunctionsTriggerResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateEdgeFunctionsTriggerResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEdgeFunctionsTriggerID)
	if err != nil {
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteEdgeFunctionsActionResp : create an edge funtions trigger response.
type DeleteEdgeFunctionsActionResp struct {
	// edge function action id.
	Result *EdgeFunctionsActionID `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	// An array with errors.
	Errors []string `json:"errors,omitempty"`

	// An array with messages.
	Messages []string `json:"messages,omitempty"`
}


// UnmarshalDeleteEdgeFunctionsActionResp unmarshals an instance of DeleteEdgeFunctionsActionResp from the specified map of raw messages.
func UnmarshalDeleteEdgeFunctionsActionResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteEdgeFunctionsActionResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEdgeFunctionsActionID)
	if err != nil {
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EdgeFunctionsActionID : edge function action id.
type EdgeFunctionsActionID struct {
	// edge functions action identifier tag.
	ID *string `json:"id,omitempty"`
}


// UnmarshalEdgeFunctionsActionID unmarshals an instance of EdgeFunctionsActionID from the specified map of raw messages.
func UnmarshalEdgeFunctionsActionID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EdgeFunctionsActionID)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EdgeFunctionsActionResp : edge function script.
type EdgeFunctionsActionResp struct {
	// Raw script content, as a string.
	Script *string `json:"script,omitempty"`

	// Hashed script content, can be used in a If-None-Match header when updating.
	Etag *string `json:"etag,omitempty"`

	// handlers.
	Handlers []string `json:"handlers,omitempty"`

	// The time when the script was last modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`

	// The time when the script was last created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// An array with items in the list response.
	Routes []EdgeFunctionsTriggerResp `json:"routes,omitempty"`
}


// UnmarshalEdgeFunctionsActionResp unmarshals an instance of EdgeFunctionsActionResp from the specified map of raw messages.
func UnmarshalEdgeFunctionsActionResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EdgeFunctionsActionResp)
	err = core.UnmarshalPrimitive(m, "script", &obj.Script)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "etag", &obj.Etag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "handlers", &obj.Handlers)
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
	err = core.UnmarshalModel(m, "routes", &obj.Routes, UnmarshalEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EdgeFunctionsTriggerID : edge function trigger id.
type EdgeFunctionsTriggerID struct {
	// edge functions trigger identifier tag.
	ID *string `json:"id,omitempty"`
}


// UnmarshalEdgeFunctionsTriggerID unmarshals an instance of EdgeFunctionsTriggerID from the specified map of raw messages.
func UnmarshalEdgeFunctionsTriggerID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EdgeFunctionsTriggerID)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EdgeFunctionsTriggerResp : edge function trigger id.
type EdgeFunctionsTriggerResp struct {
	// edge functions trigger identifier tag.
	ID *string `json:"id,omitempty"`

	// a string pattern.
	Pattern *string `json:"pattern,omitempty"`

	// Name of the script to apply when the route is matched. The route is skipped when this is blank/missing.
	Script *string `json:"script,omitempty"`

	// request limit fail open or not.
	RequestLimitFailOpen *bool `json:"request_limit_fail_open,omitempty"`
}


// UnmarshalEdgeFunctionsTriggerResp unmarshals an instance of EdgeFunctionsTriggerResp from the specified map of raw messages.
func UnmarshalEdgeFunctionsTriggerResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EdgeFunctionsTriggerResp)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "script", &obj.Script)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "request_limit_fail_open", &obj.RequestLimitFailOpen)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetEdgeFunctionsActionResp : edge funtions action response.
type GetEdgeFunctionsActionResp struct {
	// edge function script.
	Result *EdgeFunctionsActionResp `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	// An array with errors.
	Errors []string `json:"errors,omitempty"`

	// An array with messages.
	Messages []string `json:"messages,omitempty"`
}


// UnmarshalGetEdgeFunctionsActionResp unmarshals an instance of GetEdgeFunctionsActionResp from the specified map of raw messages.
func UnmarshalGetEdgeFunctionsActionResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetEdgeFunctionsActionResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEdgeFunctionsActionResp)
	if err != nil {
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetEdgeFunctionsTriggerResp : edge funtions trigger response.
type GetEdgeFunctionsTriggerResp struct {
	// edge function trigger id.
	Result *EdgeFunctionsTriggerResp `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	// An array with errors.
	Errors []string `json:"errors,omitempty"`

	// An array with messages.
	Messages []string `json:"messages,omitempty"`
}


// UnmarshalGetEdgeFunctionsTriggerResp unmarshals an instance of GetEdgeFunctionsTriggerResp from the specified map of raw messages.
func UnmarshalGetEdgeFunctionsTriggerResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetEdgeFunctionsTriggerResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListEdgeFunctionsActionsResp : edge funtions actions response.
type ListEdgeFunctionsActionsResp struct {
	// An array with items in the list response.
	Result []EdgeFunctionsActionResp `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	// An array with errors.
	Errors []string `json:"errors,omitempty"`

	// An array with messages.
	Messages []string `json:"messages,omitempty"`
}


// UnmarshalListEdgeFunctionsActionsResp unmarshals an instance of ListEdgeFunctionsActionsResp from the specified map of raw messages.
func UnmarshalListEdgeFunctionsActionsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListEdgeFunctionsActionsResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEdgeFunctionsActionResp)
	if err != nil {
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListEdgeFunctionsTriggersResp : edge funtions triggers response.
type ListEdgeFunctionsTriggersResp struct {
	// An array with items in the list response.
	Result []EdgeFunctionsTriggerResp `json:"result,omitempty"`

	// success.
	Success *bool `json:"success,omitempty"`

	// An array with errors.
	Errors []string `json:"errors,omitempty"`

	// An array with messages.
	Messages []string `json:"messages,omitempty"`
}


// UnmarshalListEdgeFunctionsTriggersResp unmarshals an instance of ListEdgeFunctionsTriggersResp from the specified map of raw messages.
func UnmarshalListEdgeFunctionsTriggersResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListEdgeFunctionsTriggersResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalEdgeFunctionsTriggerResp)
	if err != nil {
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
