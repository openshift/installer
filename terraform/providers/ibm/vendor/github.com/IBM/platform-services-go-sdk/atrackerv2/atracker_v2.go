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
 * IBM OpenAPI SDK Code Generator Version: 3.94.1-71478489-20240820-161623
 */

// Package atrackerv2 : Operations and models for the AtrackerV2 service
package atrackerv2

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

// AtrackerV2 : IBM Cloud Activity Tracker allows you to configure how auditing events are collected and stored in each
// region in your account. Events can be sent to Cloud Object Storage bucket, Logdna or Event Streams.
//
// API Version: 2.0.0
type AtrackerV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.atracker.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "atracker"

// AtrackerV2Options : Service options
type AtrackerV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewAtrackerV2UsingExternalConfig : constructs an instance of AtrackerV2 with passed in options and external configuration.
func NewAtrackerV2UsingExternalConfig(options *AtrackerV2Options) (atracker *AtrackerV2, err error) {
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

	atracker, err = NewAtrackerV2(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = atracker.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = atracker.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewAtrackerV2 : constructs an instance of AtrackerV2 with passed in options.
func NewAtrackerV2(options *AtrackerV2Options) (service *AtrackerV2, err error) {
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

	service = &AtrackerV2{
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
		"eu-de": "https://eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the eu-de region.
		"private.eu-de": "https://private.eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the eu-de region.
		"eu-gb": "https://eu-gb.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the eu-gb region.
		"private.eu-gb": "https://private.eu-gb.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the eu-gb region.
		"eu-es": "https://eu-es.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the eu-es region.
		"private.eu-es": "https://private.eu-es.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the eu-es region.
		"au-syd": "https://au-syd.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the au-syd region.
		"private.au-syd": "https://private.au-syd.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service in the au-syd region.
		"ca-tor": "https://us-east.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for ca-tor points to the us-east region.
		"private.ca-tor": "https://private.us-east.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for ca-tor points to the us-east region.
		"br-sao": "https://us-south.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for br-sao points to the us-south region.
		"private.br-sao": "https://private.us-south.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for br-sao points to the us-south region.
		"eu-fr2": "https://eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for eu-fr2 points to the eu-de region.
		"private.eu-fr2": "https://private.eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for eu-fr2 points to the eu-de region.
		"jp-tok": "https://eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for jp-tok points to the eu-de region.
		"private.jp-tok": "https://private.eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for jp-tok points to the eu-de region.
		"jp-osa": "https://eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for jp-osa points to the eu-de region.
		"private.jp-osa": "https://private.eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for jp-osa points to the eu-de region.
		"in-che": "https://eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for in-che points to the eu-de region.
		"private.in-che": "https://private.eu-de.atracker.cloud.ibm.com", // The server for IBM Cloud Activity Tracker Service for in-che points to the eu-de region.
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", core.SDKErrorf(nil, fmt.Sprintf("service URL for region '%s' not found", region), "invalid-region", common.GetComponentInfo())
}

// Clone makes a copy of "atracker" suitable for processing requests.
func (atracker *AtrackerV2) Clone() *AtrackerV2 {
	if core.IsNil(atracker) {
		return nil
	}
	clone := *atracker
	clone.Service = atracker.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (atracker *AtrackerV2) SetServiceURL(url string) error {
	err := atracker.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (atracker *AtrackerV2) GetServiceURL() string {
	return atracker.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (atracker *AtrackerV2) SetDefaultHeaders(headers http.Header) {
	atracker.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (atracker *AtrackerV2) SetEnableGzipCompression(enableGzip bool) {
	atracker.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (atracker *AtrackerV2) GetEnableGzipCompression() bool {
	return atracker.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (atracker *AtrackerV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	atracker.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (atracker *AtrackerV2) DisableRetries() {
	atracker.Service.DisableRetries()
}

// CreateTarget : Create a target
// Creates a target that includes information about the endpoint and the credentials required to write to that target.
// You can send your logs from all regions to a single target, different targets or multiple targets. One target per
// region is not required. You can define up to 16 targets per account.
func (atracker *AtrackerV2) CreateTarget(createTargetOptions *CreateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	result, response, err = atracker.CreateTargetWithContext(context.Background(), createTargetOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTargetWithContext is an alternate form of the CreateTarget method which supports a Context parameter
func (atracker *AtrackerV2) CreateTargetWithContext(ctx context.Context, createTargetOptions *CreateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTargetOptions, "createTargetOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTargetOptions, "createTargetOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/targets`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "CreateTarget")
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
	if createTargetOptions.LogdnaEndpoint != nil {
		body["logdna_endpoint"] = createTargetOptions.LogdnaEndpoint
	}
	if createTargetOptions.EventstreamsEndpoint != nil {
		body["eventstreams_endpoint"] = createTargetOptions.EventstreamsEndpoint
	}
	if createTargetOptions.CloudlogsEndpoint != nil {
		body["cloudlogs_endpoint"] = createTargetOptions.CloudlogsEndpoint
	}
	if createTargetOptions.Region != nil {
		body["region"] = createTargetOptions.Region
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_target", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListTargets : List targets
// List all targets that are defined for your account.
func (atracker *AtrackerV2) ListTargets(listTargetsOptions *ListTargetsOptions) (result *TargetList, response *core.DetailedResponse, err error) {
	result, response, err = atracker.ListTargetsWithContext(context.Background(), listTargetsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTargetsWithContext is an alternate form of the ListTargets method which supports a Context parameter
func (atracker *AtrackerV2) ListTargetsWithContext(ctx context.Context, listTargetsOptions *ListTargetsOptions) (result *TargetList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listTargetsOptions, "listTargetsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/targets`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTargetsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "ListTargets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listTargetsOptions.Region != nil {
		builder.AddQuery("region", fmt.Sprint(*listTargetsOptions.Region))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_targets", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTargetList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTarget : Get details of a target
// Retrieve the configuration details of a target.
func (atracker *AtrackerV2) GetTarget(getTargetOptions *GetTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	result, response, err = atracker.GetTargetWithContext(context.Background(), getTargetOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTargetWithContext is an alternate form of the GetTarget method which supports a Context parameter
func (atracker *AtrackerV2) GetTargetWithContext(ctx context.Context, getTargetOptions *GetTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTargetOptions, "getTargetOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTargetOptions, "getTargetOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/targets/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "GetTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_target", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTarget : Update a target
// Update the configuration details of a target.
func (atracker *AtrackerV2) ReplaceTarget(replaceTargetOptions *ReplaceTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	result, response, err = atracker.ReplaceTargetWithContext(context.Background(), replaceTargetOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceTargetWithContext is an alternate form of the ReplaceTarget method which supports a Context parameter
func (atracker *AtrackerV2) ReplaceTargetWithContext(ctx context.Context, replaceTargetOptions *ReplaceTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTargetOptions, "replaceTargetOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceTargetOptions, "replaceTargetOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *replaceTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/targets/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "ReplaceTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceTargetOptions.Name != nil {
		body["name"] = replaceTargetOptions.Name
	}
	if replaceTargetOptions.CosEndpoint != nil {
		body["cos_endpoint"] = replaceTargetOptions.CosEndpoint
	}
	if replaceTargetOptions.LogdnaEndpoint != nil {
		body["logdna_endpoint"] = replaceTargetOptions.LogdnaEndpoint
	}
	if replaceTargetOptions.EventstreamsEndpoint != nil {
		body["eventstreams_endpoint"] = replaceTargetOptions.EventstreamsEndpoint
	}
	if replaceTargetOptions.CloudlogsEndpoint != nil {
		body["cloudlogs_endpoint"] = replaceTargetOptions.CloudlogsEndpoint
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_target", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTarget : Delete a target
// Delete a target.
func (atracker *AtrackerV2) DeleteTarget(deleteTargetOptions *DeleteTargetOptions) (result *WarningReport, response *core.DetailedResponse, err error) {
	result, response, err = atracker.DeleteTargetWithContext(context.Background(), deleteTargetOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTargetWithContext is an alternate form of the DeleteTarget method which supports a Context parameter
func (atracker *AtrackerV2) DeleteTargetWithContext(ctx context.Context, deleteTargetOptions *DeleteTargetOptions) (result *WarningReport, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTargetOptions, "deleteTargetOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTargetOptions, "deleteTargetOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/targets/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "DeleteTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_target", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalWarningReport)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ValidateTarget : Validate a target
// Validate a target by checking the credentials to write to the target. The result is included as additional data of
// the target in the section "write_status".
func (atracker *AtrackerV2) ValidateTarget(validateTargetOptions *ValidateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	result, response, err = atracker.ValidateTargetWithContext(context.Background(), validateTargetOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ValidateTargetWithContext is an alternate form of the ValidateTarget method which supports a Context parameter
func (atracker *AtrackerV2) ValidateTargetWithContext(ctx context.Context, validateTargetOptions *ValidateTargetOptions) (result *Target, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(validateTargetOptions, "validateTargetOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(validateTargetOptions, "validateTargetOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *validateTargetOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/targets/{id}/validate`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range validateTargetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "ValidateTarget")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "validate_target", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTarget)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateRoute : Create a route
// Create a route to define the rule that specifies how to manage auditing events.
func (atracker *AtrackerV2) CreateRoute(createRouteOptions *CreateRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	result, response, err = atracker.CreateRouteWithContext(context.Background(), createRouteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateRouteWithContext is an alternate form of the CreateRoute method which supports a Context parameter
func (atracker *AtrackerV2) CreateRouteWithContext(ctx context.Context, createRouteOptions *CreateRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRouteOptions, "createRouteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createRouteOptions, "createRouteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/routes`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "CreateRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRouteOptions.Name != nil {
		body["name"] = createRouteOptions.Name
	}
	if createRouteOptions.Rules != nil {
		body["rules"] = createRouteOptions.Rules
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_route", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoute)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListRoutes : List routes
// List the route that is configured for an account.
func (atracker *AtrackerV2) ListRoutes(listRoutesOptions *ListRoutesOptions) (result *RouteList, response *core.DetailedResponse, err error) {
	result, response, err = atracker.ListRoutesWithContext(context.Background(), listRoutesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListRoutesWithContext is an alternate form of the ListRoutes method which supports a Context parameter
func (atracker *AtrackerV2) ListRoutesWithContext(ctx context.Context, listRoutesOptions *ListRoutesOptions) (result *RouteList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRoutesOptions, "listRoutesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/routes`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listRoutesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "ListRoutes")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_routes", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRouteList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetRoute : Get details of a route
// Get the configuration details of a route.
func (atracker *AtrackerV2) GetRoute(getRouteOptions *GetRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	result, response, err = atracker.GetRouteWithContext(context.Background(), getRouteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetRouteWithContext is an alternate form of the GetRoute method which supports a Context parameter
func (atracker *AtrackerV2) GetRouteWithContext(ctx context.Context, getRouteOptions *GetRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRouteOptions, "getRouteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getRouteOptions, "getRouteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getRouteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/routes/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "GetRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_route", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoute)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceRoute : Update a route
// Update the configuration details of a route.
func (atracker *AtrackerV2) ReplaceRoute(replaceRouteOptions *ReplaceRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	result, response, err = atracker.ReplaceRouteWithContext(context.Background(), replaceRouteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceRouteWithContext is an alternate form of the ReplaceRoute method which supports a Context parameter
func (atracker *AtrackerV2) ReplaceRouteWithContext(ctx context.Context, replaceRouteOptions *ReplaceRouteOptions) (result *Route, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceRouteOptions, "replaceRouteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceRouteOptions, "replaceRouteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *replaceRouteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/routes/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "ReplaceRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceRouteOptions.Name != nil {
		body["name"] = replaceRouteOptions.Name
	}
	if replaceRouteOptions.Rules != nil {
		body["rules"] = replaceRouteOptions.Rules
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_route", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRoute)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteRoute : Delete a route
// Deletes a route.
func (atracker *AtrackerV2) DeleteRoute(deleteRouteOptions *DeleteRouteOptions) (response *core.DetailedResponse, err error) {
	response, err = atracker.DeleteRouteWithContext(context.Background(), deleteRouteOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteRouteWithContext is an alternate form of the DeleteRoute method which supports a Context parameter
func (atracker *AtrackerV2) DeleteRouteWithContext(ctx context.Context, deleteRouteOptions *DeleteRouteOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRouteOptions, "deleteRouteOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteRouteOptions, "deleteRouteOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteRouteOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/routes/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteRouteOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "DeleteRoute")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = atracker.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_route", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetSettings : Get settings
// Get information about the current settings including default targets.
func (atracker *AtrackerV2) GetSettings(getSettingsOptions *GetSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	result, response, err = atracker.GetSettingsWithContext(context.Background(), getSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSettingsWithContext is an alternate form of the GetSettings method which supports a Context parameter
func (atracker *AtrackerV2) GetSettingsWithContext(ctx context.Context, getSettingsOptions *GetSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getSettingsOptions, "getSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/settings`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "GetSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettings)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PutSettings : Modify settings
// Modify the current settings such as default targets, permitted target regions, metadata region primary and secondary.
func (atracker *AtrackerV2) PutSettings(putSettingsOptions *PutSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	result, response, err = atracker.PutSettingsWithContext(context.Background(), putSettingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PutSettingsWithContext is an alternate form of the PutSettings method which supports a Context parameter
func (atracker *AtrackerV2) PutSettingsWithContext(ctx context.Context, putSettingsOptions *PutSettingsOptions) (result *Settings, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(putSettingsOptions, "putSettingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(putSettingsOptions, "putSettingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = atracker.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(atracker.Service.Options.URL, `/api/v2/settings`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range putSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("atracker", "V2", "PutSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if putSettingsOptions.MetadataRegionPrimary != nil {
		body["metadata_region_primary"] = putSettingsOptions.MetadataRegionPrimary
	}
	if putSettingsOptions.PrivateAPIEndpointOnly != nil {
		body["private_api_endpoint_only"] = putSettingsOptions.PrivateAPIEndpointOnly
	}
	if putSettingsOptions.DefaultTargets != nil {
		body["default_targets"] = putSettingsOptions.DefaultTargets
	}
	if putSettingsOptions.PermittedTargetRegions != nil {
		body["permitted_target_regions"] = putSettingsOptions.PermittedTargetRegions
	}
	if putSettingsOptions.MetadataRegionBackup != nil {
		body["metadata_region_backup"] = putSettingsOptions.MetadataRegionBackup
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
	response, err = atracker.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "put_settings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSettings)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "2.0.0")
}

// CloudLogsEndpoint : Property values for the IBM Cloud Logs endpoint in responses.
type CloudLogsEndpoint struct {
	// The CRN of the IBM Cloud Logs instance.
	TargetCRN *string `json:"target_crn" validate:"required"`
}

// UnmarshalCloudLogsEndpoint unmarshals an instance of CloudLogsEndpoint from the specified map of raw messages.
func UnmarshalCloudLogsEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudLogsEndpoint)
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CloudLogsEndpointPrototype : Property values for an IBM Cloud Logs endpoint in requests.
type CloudLogsEndpointPrototype struct {
	// The CRN of the IBM Cloud Logs instance.
	TargetCRN *string `json:"target_crn" validate:"required"`
}

// NewCloudLogsEndpointPrototype : Instantiate CloudLogsEndpointPrototype (Generic Model Constructor)
func (*AtrackerV2) NewCloudLogsEndpointPrototype(targetCRN string) (_model *CloudLogsEndpointPrototype, err error) {
	_model = &CloudLogsEndpointPrototype{
		TargetCRN: core.StringPtr(targetCRN),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCloudLogsEndpointPrototype unmarshals an instance of CloudLogsEndpointPrototype from the specified map of raw messages.
func UnmarshalCloudLogsEndpointPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CloudLogsEndpointPrototype)
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CosEndpoint : Property values for a Cloud Object Storage Endpoint in responses.
type CosEndpoint struct {
	// The host name of the Cloud Object Storage endpoint.
	Endpoint *string `json:"endpoint" validate:"required"`

	// The CRN of the Cloud Object Storage instance.
	TargetCRN *string `json:"target_crn" validate:"required"`

	// The bucket name under the Cloud Object Storage instance.
	Bucket *string `json:"bucket" validate:"required"`

	// Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag
	// to true if service to service is enabled and do not supply an apikey.
	ServiceToServiceEnabled *bool `json:"service_to_service_enabled" validate:"required"`
}

// UnmarshalCosEndpoint unmarshals an instance of CosEndpoint from the specified map of raw messages.
func UnmarshalCosEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CosEndpoint)
	err = core.UnmarshalPrimitive(m, "endpoint", &obj.Endpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		err = core.SDKErrorf(err, "", "bucket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_to_service_enabled", &obj.ServiceToServiceEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_to_service_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CosEndpointPrototype : Property values for a Cloud Object Storage Endpoint in requests.
type CosEndpointPrototype struct {
	// The host name of the Cloud Object Storage endpoint.
	Endpoint *string `json:"endpoint" validate:"required"`

	// The CRN of the Cloud Object Storage instance.
	TargetCRN *string `json:"target_crn" validate:"required"`

	// The bucket name under the Cloud Object Storage instance.
	Bucket *string `json:"bucket" validate:"required"`

	// The IAM API key that has writer access to the Cloud Object Storage instance. This credential is masked in the
	// response. This is required if service_to_service is not enabled.
	APIKey *string `json:"api_key,omitempty"`

	// Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag
	// to true if service to service is enabled and do not supply an apikey.
	ServiceToServiceEnabled *bool `json:"service_to_service_enabled,omitempty"`
}

// NewCosEndpointPrototype : Instantiate CosEndpointPrototype (Generic Model Constructor)
func (*AtrackerV2) NewCosEndpointPrototype(endpoint string, targetCRN string, bucket string) (_model *CosEndpointPrototype, err error) {
	_model = &CosEndpointPrototype{
		Endpoint: core.StringPtr(endpoint),
		TargetCRN: core.StringPtr(targetCRN),
		Bucket: core.StringPtr(bucket),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCosEndpointPrototype unmarshals an instance of CosEndpointPrototype from the specified map of raw messages.
func UnmarshalCosEndpointPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CosEndpointPrototype)
	err = core.UnmarshalPrimitive(m, "endpoint", &obj.Endpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "bucket", &obj.Bucket)
	if err != nil {
		err = core.SDKErrorf(err, "", "bucket-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_to_service_enabled", &obj.ServiceToServiceEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_to_service_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateRouteOptions : The CreateRoute options.
type CreateRouteOptions struct {
	// The name of the route. The name must be 1000 characters or less and cannot include any special characters other than
	// `(space) - . _ :`. Do not include any personal identifying information (PII) in any resource names.
	Name *string `json:"name" validate:"required"`

	// Routing rules that will be evaluated in their order of the array.
	Rules []RulePrototype `json:"rules" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateRouteOptions : Instantiate CreateRouteOptions
func (*AtrackerV2) NewCreateRouteOptions(name string, rules []RulePrototype) *CreateRouteOptions {
	return &CreateRouteOptions{
		Name: core.StringPtr(name),
		Rules: rules,
	}
}

// SetName : Allow user to set Name
func (_options *CreateRouteOptions) SetName(name string) *CreateRouteOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *CreateRouteOptions) SetRules(rules []RulePrototype) *CreateRouteOptions {
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
	// than `(space) - . _ :`. Do not include any personal identifying information (PII) in any resource names.
	Name *string `json:"name" validate:"required"`

	// The type of the target. It can be cloud_object_storage, logdna, event_streams, or cloud_logs. Based on this type you
	// must include cos_endpoint, logdna_endpoint, eventstreams_endpoint or cloudlogs_endpoint.
	TargetType *string `json:"target_type" validate:"required"`

	// Property values for a Cloud Object Storage Endpoint in requests.
	CosEndpoint *CosEndpointPrototype `json:"cos_endpoint,omitempty"`

	// Property values for a LogDNA Endpoint in requests.
	LogdnaEndpoint *LogdnaEndpointPrototype `json:"logdna_endpoint,omitempty"`

	// Property values for an Event Streams Endpoint in requests.
	EventstreamsEndpoint *EventstreamsEndpointPrototype `json:"eventstreams_endpoint,omitempty"`

	// Property values for an IBM Cloud Logs endpoint in requests.
	CloudlogsEndpoint *CloudLogsEndpointPrototype `json:"cloudlogs_endpoint,omitempty"`

	// Include this optional field if you want to create a target in a different region other than the one you are
	// connected.
	Region *string `json:"region,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateTargetOptions.TargetType property.
// The type of the target. It can be cloud_object_storage, logdna, event_streams, or cloud_logs. Based on this type you
// must include cos_endpoint, logdna_endpoint, eventstreams_endpoint or cloudlogs_endpoint.
const (
	CreateTargetOptionsTargetTypeCloudLogsConst = "cloud_logs"
	CreateTargetOptionsTargetTypeCloudObjectStorageConst = "cloud_object_storage"
	CreateTargetOptionsTargetTypeEventStreamsConst = "event_streams"
	CreateTargetOptionsTargetTypeLogdnaConst = "logdna"
)

// NewCreateTargetOptions : Instantiate CreateTargetOptions
func (*AtrackerV2) NewCreateTargetOptions(name string, targetType string) *CreateTargetOptions {
	return &CreateTargetOptions{
		Name: core.StringPtr(name),
		TargetType: core.StringPtr(targetType),
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
func (_options *CreateTargetOptions) SetCosEndpoint(cosEndpoint *CosEndpointPrototype) *CreateTargetOptions {
	_options.CosEndpoint = cosEndpoint
	return _options
}

// SetLogdnaEndpoint : Allow user to set LogdnaEndpoint
func (_options *CreateTargetOptions) SetLogdnaEndpoint(logdnaEndpoint *LogdnaEndpointPrototype) *CreateTargetOptions {
	_options.LogdnaEndpoint = logdnaEndpoint
	return _options
}

// SetEventstreamsEndpoint : Allow user to set EventstreamsEndpoint
func (_options *CreateTargetOptions) SetEventstreamsEndpoint(eventstreamsEndpoint *EventstreamsEndpointPrototype) *CreateTargetOptions {
	_options.EventstreamsEndpoint = eventstreamsEndpoint
	return _options
}

// SetCloudlogsEndpoint : Allow user to set CloudlogsEndpoint
func (_options *CreateTargetOptions) SetCloudlogsEndpoint(cloudlogsEndpoint *CloudLogsEndpointPrototype) *CreateTargetOptions {
	_options.CloudlogsEndpoint = cloudlogsEndpoint
	return _options
}

// SetRegion : Allow user to set Region
func (_options *CreateTargetOptions) SetRegion(region string) *CreateTargetOptions {
	_options.Region = core.StringPtr(region)
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
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteRouteOptions : Instantiate DeleteRouteOptions
func (*AtrackerV2) NewDeleteRouteOptions(id string) *DeleteRouteOptions {
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
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTargetOptions : Instantiate DeleteTargetOptions
func (*AtrackerV2) NewDeleteTargetOptions(id string) *DeleteTargetOptions {
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

// EventstreamsEndpoint : Property values for the Event Streams Endpoint in responses.
type EventstreamsEndpoint struct {
	// The CRN of the Event Streams instance.
	TargetCRN *string `json:"target_crn" validate:"required"`

	// List of broker endpoints.
	Brokers []string `json:"brokers" validate:"required"`

	// The messsage hub topic defined in the Event Streams instance.
	Topic *string `json:"topic" validate:"required"`

	// The user password (api key) for the message hub topic in the Event Streams instance.
	APIKey *string `json:"api_key,omitempty"`

	// Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag
	// to true if service to service is enabled and do not supply an apikey.
	ServiceToServiceEnabled *bool `json:"service_to_service_enabled,omitempty"`
}

// UnmarshalEventstreamsEndpoint unmarshals an instance of EventstreamsEndpoint from the specified map of raw messages.
func UnmarshalEventstreamsEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventstreamsEndpoint)
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "brokers", &obj.Brokers)
	if err != nil {
		err = core.SDKErrorf(err, "", "brokers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "topic", &obj.Topic)
	if err != nil {
		err = core.SDKErrorf(err, "", "topic-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_to_service_enabled", &obj.ServiceToServiceEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_to_service_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EventstreamsEndpointPrototype : Property values for an Event Streams Endpoint in requests.
type EventstreamsEndpointPrototype struct {
	// The CRN of the Event Streams instance.
	TargetCRN *string `json:"target_crn" validate:"required"`

	// List of broker endpoints.
	Brokers []string `json:"brokers" validate:"required"`

	// The messsage hub topic defined in the Event Streams instance.
	Topic *string `json:"topic" validate:"required"`

	// The user password (api key) for the message hub topic in the Event Streams instance.
	APIKey *string `json:"api_key,omitempty"`

	// Determines if IBM Cloud Activity Tracker Event Routing has service to service authentication enabled. Set this flag
	// to true if service to service is enabled and do not supply an apikey.
	ServiceToServiceEnabled *bool `json:"service_to_service_enabled,omitempty"`
}

// NewEventstreamsEndpointPrototype : Instantiate EventstreamsEndpointPrototype (Generic Model Constructor)
func (*AtrackerV2) NewEventstreamsEndpointPrototype(targetCRN string, brokers []string, topic string) (_model *EventstreamsEndpointPrototype, err error) {
	_model = &EventstreamsEndpointPrototype{
		TargetCRN: core.StringPtr(targetCRN),
		Brokers: brokers,
		Topic: core.StringPtr(topic),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalEventstreamsEndpointPrototype unmarshals an instance of EventstreamsEndpointPrototype from the specified map of raw messages.
func UnmarshalEventstreamsEndpointPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EventstreamsEndpointPrototype)
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "brokers", &obj.Brokers)
	if err != nil {
		err = core.SDKErrorf(err, "", "brokers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "topic", &obj.Topic)
	if err != nil {
		err = core.SDKErrorf(err, "", "topic-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.APIKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_to_service_enabled", &obj.ServiceToServiceEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_to_service_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetRouteOptions : The GetRoute options.
type GetRouteOptions struct {
	// The v4 UUID that uniquely identifies the route.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetRouteOptions : Instantiate GetRouteOptions
func (*AtrackerV2) NewGetRouteOptions(id string) *GetRouteOptions {
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

// GetSettingsOptions : The GetSettings options.
type GetSettingsOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetSettingsOptions : Instantiate GetSettingsOptions
func (*AtrackerV2) NewGetSettingsOptions() *GetSettingsOptions {
	return &GetSettingsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetSettingsOptions) SetHeaders(param map[string]string) *GetSettingsOptions {
	options.Headers = param
	return options
}

// GetTargetOptions : The GetTarget options.
type GetTargetOptions struct {
	// The v4 UUID that uniquely identifies the target.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTargetOptions : Instantiate GetTargetOptions
func (*AtrackerV2) NewGetTargetOptions(id string) *GetTargetOptions {
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

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListRoutesOptions : Instantiate ListRoutesOptions
func (*AtrackerV2) NewListRoutesOptions() *ListRoutesOptions {
	return &ListRoutesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListRoutesOptions) SetHeaders(param map[string]string) *ListRoutesOptions {
	options.Headers = param
	return options
}

// ListTargetsOptions : The ListTargets options.
type ListTargetsOptions struct {
	// Limit the query to the specified region.
	Region *string `json:"region,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListTargetsOptions : Instantiate ListTargetsOptions
func (*AtrackerV2) NewListTargetsOptions() *ListTargetsOptions {
	return &ListTargetsOptions{}
}

// SetRegion : Allow user to set Region
func (_options *ListTargetsOptions) SetRegion(region string) *ListTargetsOptions {
	_options.Region = core.StringPtr(region)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTargetsOptions) SetHeaders(param map[string]string) *ListTargetsOptions {
	options.Headers = param
	return options
}

// LogdnaEndpoint : Property values for a LogDNA Endpoint in responses.
type LogdnaEndpoint struct {
	// The CRN of the LogDNA instance.
	TargetCRN *string `json:"target_crn" validate:"required"`
}

// UnmarshalLogdnaEndpoint unmarshals an instance of LogdnaEndpoint from the specified map of raw messages.
func UnmarshalLogdnaEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogdnaEndpoint)
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogdnaEndpointPrototype : Property values for a LogDNA Endpoint in requests.
type LogdnaEndpointPrototype struct {
	// The CRN of the LogDNA instance.
	TargetCRN *string `json:"target_crn" validate:"required"`

	// The LogDNA ingestion key is used for routing logs to a specific LogDNA instance.
	IngestionKey *string `json:"ingestion_key" validate:"required"`
}

// NewLogdnaEndpointPrototype : Instantiate LogdnaEndpointPrototype (Generic Model Constructor)
func (*AtrackerV2) NewLogdnaEndpointPrototype(targetCRN string, ingestionKey string) (_model *LogdnaEndpointPrototype, err error) {
	_model = &LogdnaEndpointPrototype{
		TargetCRN: core.StringPtr(targetCRN),
		IngestionKey: core.StringPtr(ingestionKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalLogdnaEndpointPrototype unmarshals an instance of LogdnaEndpointPrototype from the specified map of raw messages.
func UnmarshalLogdnaEndpointPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogdnaEndpointPrototype)
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ingestion_key", &obj.IngestionKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "ingestion_key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PutSettingsOptions : The PutSettings options.
type PutSettingsOptions struct {
	// To store all your meta data in a single region.
	MetadataRegionPrimary *string `json:"metadata_region_primary" validate:"required"`

	// If you set this true then you cannot access api through public network.
	PrivateAPIEndpointOnly *bool `json:"private_api_endpoint_only" validate:"required"`

	// The target ID List. In the event that no routing rule causes the event to be sent to a target, these targets will
	// receive the event.
	DefaultTargets []string `json:"default_targets,omitempty"`

	// If present then only these regions may be used to define a target.
	PermittedTargetRegions []string `json:"permitted_target_regions,omitempty"`

	// To store all your meta data in a backup region.
	MetadataRegionBackup *string `json:"metadata_region_backup,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewPutSettingsOptions : Instantiate PutSettingsOptions
func (*AtrackerV2) NewPutSettingsOptions(metadataRegionPrimary string, privateAPIEndpointOnly bool) *PutSettingsOptions {
	return &PutSettingsOptions{
		MetadataRegionPrimary: core.StringPtr(metadataRegionPrimary),
		PrivateAPIEndpointOnly: core.BoolPtr(privateAPIEndpointOnly),
	}
}

// SetMetadataRegionPrimary : Allow user to set MetadataRegionPrimary
func (_options *PutSettingsOptions) SetMetadataRegionPrimary(metadataRegionPrimary string) *PutSettingsOptions {
	_options.MetadataRegionPrimary = core.StringPtr(metadataRegionPrimary)
	return _options
}

// SetPrivateAPIEndpointOnly : Allow user to set PrivateAPIEndpointOnly
func (_options *PutSettingsOptions) SetPrivateAPIEndpointOnly(privateAPIEndpointOnly bool) *PutSettingsOptions {
	_options.PrivateAPIEndpointOnly = core.BoolPtr(privateAPIEndpointOnly)
	return _options
}

// SetDefaultTargets : Allow user to set DefaultTargets
func (_options *PutSettingsOptions) SetDefaultTargets(defaultTargets []string) *PutSettingsOptions {
	_options.DefaultTargets = defaultTargets
	return _options
}

// SetPermittedTargetRegions : Allow user to set PermittedTargetRegions
func (_options *PutSettingsOptions) SetPermittedTargetRegions(permittedTargetRegions []string) *PutSettingsOptions {
	_options.PermittedTargetRegions = permittedTargetRegions
	return _options
}

// SetMetadataRegionBackup : Allow user to set MetadataRegionBackup
func (_options *PutSettingsOptions) SetMetadataRegionBackup(metadataRegionBackup string) *PutSettingsOptions {
	_options.MetadataRegionBackup = core.StringPtr(metadataRegionBackup)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PutSettingsOptions) SetHeaders(param map[string]string) *PutSettingsOptions {
	options.Headers = param
	return options
}

// ReplaceRouteOptions : The ReplaceRoute options.
type ReplaceRouteOptions struct {
	// The v4 UUID that uniquely identifies the route.
	ID *string `json:"id" validate:"required,ne="`

	// The name of the route. The name must be 1000 characters or less and cannot include any special characters other than
	// `(space) - . _ :`. Do not include any personal identifying information (PII) in any resource names.
	Name *string `json:"name" validate:"required"`

	// Routing rules that will be evaluated in their order of the array.
	Rules []RulePrototype `json:"rules" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewReplaceRouteOptions : Instantiate ReplaceRouteOptions
func (*AtrackerV2) NewReplaceRouteOptions(id string, name string, rules []RulePrototype) *ReplaceRouteOptions {
	return &ReplaceRouteOptions{
		ID: core.StringPtr(id),
		Name: core.StringPtr(name),
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

// SetRules : Allow user to set Rules
func (_options *ReplaceRouteOptions) SetRules(rules []RulePrototype) *ReplaceRouteOptions {
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
	ID *string `json:"id" validate:"required,ne="`

	// The name of the target. The name must be 1000 characters or less, and cannot include any special characters other
	// than `(space) - . _ :`. Do not include any personal identifying information (PII) in any resource names.
	Name *string `json:"name,omitempty"`

	// Property values for a Cloud Object Storage Endpoint in requests.
	CosEndpoint *CosEndpointPrototype `json:"cos_endpoint,omitempty"`

	// Property values for a LogDNA Endpoint in requests.
	LogdnaEndpoint *LogdnaEndpointPrototype `json:"logdna_endpoint,omitempty"`

	// Property values for an Event Streams Endpoint in requests.
	EventstreamsEndpoint *EventstreamsEndpointPrototype `json:"eventstreams_endpoint,omitempty"`

	// Property values for an IBM Cloud Logs endpoint in requests.
	CloudlogsEndpoint *CloudLogsEndpointPrototype `json:"cloudlogs_endpoint,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewReplaceTargetOptions : Instantiate ReplaceTargetOptions
func (*AtrackerV2) NewReplaceTargetOptions(id string) *ReplaceTargetOptions {
	return &ReplaceTargetOptions{
		ID: core.StringPtr(id),
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

// SetCosEndpoint : Allow user to set CosEndpoint
func (_options *ReplaceTargetOptions) SetCosEndpoint(cosEndpoint *CosEndpointPrototype) *ReplaceTargetOptions {
	_options.CosEndpoint = cosEndpoint
	return _options
}

// SetLogdnaEndpoint : Allow user to set LogdnaEndpoint
func (_options *ReplaceTargetOptions) SetLogdnaEndpoint(logdnaEndpoint *LogdnaEndpointPrototype) *ReplaceTargetOptions {
	_options.LogdnaEndpoint = logdnaEndpoint
	return _options
}

// SetEventstreamsEndpoint : Allow user to set EventstreamsEndpoint
func (_options *ReplaceTargetOptions) SetEventstreamsEndpoint(eventstreamsEndpoint *EventstreamsEndpointPrototype) *ReplaceTargetOptions {
	_options.EventstreamsEndpoint = eventstreamsEndpoint
	return _options
}

// SetCloudlogsEndpoint : Allow user to set CloudlogsEndpoint
func (_options *ReplaceTargetOptions) SetCloudlogsEndpoint(cloudlogsEndpoint *CloudLogsEndpointPrototype) *ReplaceTargetOptions {
	_options.CloudlogsEndpoint = cloudlogsEndpoint
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTargetOptions) SetHeaders(param map[string]string) *ReplaceTargetOptions {
	options.Headers = param
	return options
}

// Route : The route resource. The scope of the route is account wide. That means all the routes are evaluated in all regions,
// except the ones limited by region.
type Route struct {
	// The uuid of the route resource.
	ID *string `json:"id" validate:"required"`

	// The name of the route.
	Name *string `json:"name" validate:"required"`

	// The crn of the route resource.
	CRN *string `json:"crn" validate:"required"`

	// The version of the route.
	Version *int64 `json:"version,omitempty"`

	// The routing rules that will be evaluated in their order of the array. Once a rule is matched, the remaining rules in
	// the route definition will be skipped.
	Rules []Rule `json:"rules" validate:"required"`

	// The timestamp of the route creation time.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The timestamp of the route last updated time.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// The API version of the route.
	APIVersion *int64 `json:"api_version" validate:"required"`

	// An optional message containing information about the route.
	Message *string `json:"message,omitempty"`
}

// UnmarshalRoute unmarshals an instance of Route from the specified map of raw messages.
func UnmarshalRoute(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Route)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_version", &obj.APIVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "routes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Rule : A configuration to route events to pre-defined target.
type Rule struct {
	// The target ID List. All the events will be send to all targets listed in the rule. You can include targets from
	// other regions.
	TargetIds []string `json:"target_ids" validate:"required"`

	// Logs from these locations will be sent to the targets specified. Locations is a superset of regions including global
	// and *.
	Locations []string `json:"locations" validate:"required"`
}

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "target_ids", &obj.TargetIds)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_ids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locations", &obj.Locations)
	if err != nil {
		err = core.SDKErrorf(err, "", "locations-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RulePrototype : A configuration to route events to pre-defined target.
type RulePrototype struct {
	// The target ID List. All the events will be send to all targets listed in the rule. You can include targets from
	// other regions.
	TargetIds []string `json:"target_ids" validate:"required"`

	// Logs from these locations will be sent to the targets specified. Locations is a superset of regions including global
	// and *.
	Locations []string `json:"locations,omitempty"`
}

// NewRulePrototype : Instantiate RulePrototype (Generic Model Constructor)
func (*AtrackerV2) NewRulePrototype(targetIds []string) (_model *RulePrototype, err error) {
	_model = &RulePrototype{
		TargetIds: targetIds,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRulePrototype unmarshals an instance of RulePrototype from the specified map of raw messages.
func UnmarshalRulePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RulePrototype)
	err = core.UnmarshalPrimitive(m, "target_ids", &obj.TargetIds)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_ids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locations", &obj.Locations)
	if err != nil {
		err = core.SDKErrorf(err, "", "locations-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Settings : Activity Tracker Event Routing settings response.
type Settings struct {
	// The target ID List. In the event that no routing rule causes the event to be sent to a target, these targets will
	// receive the event.
	DefaultTargets []string `json:"default_targets" validate:"required"`

	// If present then only these regions may be used to define a target.
	PermittedTargetRegions []string `json:"permitted_target_regions" validate:"required"`

	// To store all your meta data in a single region.
	MetadataRegionPrimary *string `json:"metadata_region_primary" validate:"required"`

	// To store all your meta data in a backup region.
	MetadataRegionBackup *string `json:"metadata_region_backup,omitempty"`

	// If you set this true then you cannot access api through public network.
	PrivateAPIEndpointOnly *bool `json:"private_api_endpoint_only" validate:"required"`

	// API version used for configuring IBM Cloud Activity Tracker Event Routing resources in the account.
	APIVersion *int64 `json:"api_version" validate:"required"`

	// An optional message containing information about the audit log locations.
	Message *string `json:"message,omitempty"`
}

// UnmarshalSettings unmarshals an instance of Settings from the specified map of raw messages.
func UnmarshalSettings(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Settings)
	err = core.UnmarshalPrimitive(m, "default_targets", &obj.DefaultTargets)
	if err != nil {
		err = core.SDKErrorf(err, "", "default_targets-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "permitted_target_regions", &obj.PermittedTargetRegions)
	if err != nil {
		err = core.SDKErrorf(err, "", "permitted_target_regions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata_region_primary", &obj.MetadataRegionPrimary)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata_region_primary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "metadata_region_backup", &obj.MetadataRegionBackup)
	if err != nil {
		err = core.SDKErrorf(err, "", "metadata_region_backup-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "private_api_endpoint_only", &obj.PrivateAPIEndpointOnly)
	if err != nil {
		err = core.SDKErrorf(err, "", "private_api_endpoint_only-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_version", &obj.APIVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Target : Property values for a target in responses.
type Target struct {
	// The uuid of the target resource.
	ID *string `json:"id" validate:"required"`

	// The name of the target resource.
	Name *string `json:"name" validate:"required"`

	// The crn of the target resource.
	CRN *string `json:"crn" validate:"required"`

	// The type of the target.
	TargetType *string `json:"target_type" validate:"required"`

	// Included this optional field if you used it to create a target in a different region other than the one you are
	// connected.
	Region *string `json:"region,omitempty"`

	// Property values for a Cloud Object Storage Endpoint in responses.
	CosEndpoint *CosEndpoint `json:"cos_endpoint,omitempty"`

	// Property values for a LogDNA Endpoint in responses.
	LogdnaEndpoint *LogdnaEndpoint `json:"logdna_endpoint,omitempty"`

	// Property values for the Event Streams Endpoint in responses.
	EventstreamsEndpoint *EventstreamsEndpoint `json:"eventstreams_endpoint,omitempty"`

	// Property values for the IBM Cloud Logs endpoint in responses.
	CloudlogsEndpoint *CloudLogsEndpoint `json:"cloudlogs_endpoint,omitempty"`

	// The status of the write attempt to the target with the provided endpoint parameters.
	WriteStatus *WriteStatus `json:"write_status" validate:"required"`

	// The timestamp of the target creation time.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// The timestamp of the target last updated time.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// An optional message containing information about the target.
	Message *string `json:"message,omitempty"`

	// The API version of the target.
	APIVersion *int64 `json:"api_version" validate:"required"`
}

// Constants associated with the Target.TargetType property.
// The type of the target.
const (
	TargetTargetTypeCloudLogsConst = "cloud_logs"
	TargetTargetTypeCloudObjectStorageConst = "cloud_object_storage"
	TargetTargetTypeEventStreamsConst = "event_streams"
	TargetTargetTypeLogdnaConst = "logdna"
)

// UnmarshalTarget unmarshals an instance of Target from the specified map of raw messages.
func UnmarshalTarget(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Target)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_type", &obj.TargetType)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "cos_endpoint", &obj.CosEndpoint, UnmarshalCosEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos_endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "logdna_endpoint", &obj.LogdnaEndpoint, UnmarshalLogdnaEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "logdna_endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "eventstreams_endpoint", &obj.EventstreamsEndpoint, UnmarshalEventstreamsEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "eventstreams_endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "cloudlogs_endpoint", &obj.CloudlogsEndpoint, UnmarshalCloudLogsEndpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "cloudlogs_endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "write_status", &obj.WriteStatus, UnmarshalWriteStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "write_status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_version", &obj.APIVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_version-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "targets-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ValidateTargetOptions : The ValidateTarget options.
type ValidateTargetOptions struct {
	// The v4 UUID that uniquely identifies the target.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewValidateTargetOptions : Instantiate ValidateTargetOptions
func (*AtrackerV2) NewValidateTargetOptions(id string) *ValidateTargetOptions {
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
	Code *string `json:"code" validate:"required"`

	// The warning message.
	Message *string `json:"message" validate:"required"`
}

// UnmarshalWarning unmarshals an instance of Warning from the specified map of raw messages.
func UnmarshalWarning(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Warning)
	err = core.UnmarshalPrimitive(m, "code", &obj.Code)
	if err != nil {
		err = core.SDKErrorf(err, "", "code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		err = core.SDKErrorf(err, "", "message-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "status_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trace", &obj.Trace)
	if err != nil {
		err = core.SDKErrorf(err, "", "trace-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "warnings", &obj.Warnings, UnmarshalWarning)
	if err != nil {
		err = core.SDKErrorf(err, "", "warnings-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// WriteStatus : The status of the write attempt to the target with the provided endpoint parameters.
type WriteStatus struct {
	// The status such as failed or success.
	Status *string `json:"status" validate:"required"`

	// The timestamp of the failure.
	LastFailure *strfmt.DateTime `json:"last_failure,omitempty"`

	// Detailed description of the cause of the failure.
	ReasonForLastFailure *string `json:"reason_for_last_failure,omitempty"`
}

// UnmarshalWriteStatus unmarshals an instance of WriteStatus from the specified map of raw messages.
func UnmarshalWriteStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WriteStatus)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_failure", &obj.LastFailure)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_failure-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reason_for_last_failure", &obj.ReasonForLastFailure)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason_for_last_failure-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
