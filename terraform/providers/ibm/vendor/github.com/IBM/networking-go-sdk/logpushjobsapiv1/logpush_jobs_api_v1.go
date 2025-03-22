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
 * IBM OpenAPI SDK Code Generator Version: 3.98.0-8be2046a-20241205-162752
 */

// Package logpushjobsapiv1 : Operations and models for the LogpushJobsApiV1 service
package logpushjobsapiv1

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// LogpushJobsApiV1 : CIS Logpush Jobs
//
// API Version: 1.0.0
type LogpushJobsApiV1 struct {
	Service *core.BaseService

	// Full URL-encoded CRN of the service instance.
	Crn *string

	// The dataset.
	Dataset *string

	// Zone identifier.
	ZoneID *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "logpush_jobs_api"

// LogpushJobsApiV1Options : Service options
type LogpushJobsApiV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full URL-encoded CRN of the service instance.
	Crn *string `validate:"required"`

	// The dataset.
	Dataset *string `validate:"required"`

	// Zone identifier.
	ZoneID *string `validate:"required"`
}

// NewLogpushJobsApiV1UsingExternalConfig : constructs an instance of LogpushJobsApiV1 with passed in options and external configuration.
func NewLogpushJobsApiV1UsingExternalConfig(options *LogpushJobsApiV1Options) (logpushJobsApi *LogpushJobsApiV1, err error) {
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

	logpushJobsApi, err = NewLogpushJobsApiV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = logpushJobsApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = logpushJobsApi.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewLogpushJobsApiV1 : constructs an instance of LogpushJobsApiV1 with passed in options.
func NewLogpushJobsApiV1(options *LogpushJobsApiV1Options) (service *LogpushJobsApiV1, err error) {
	serviceOptions := &core.ServiceOptions{
		URL:           DefaultServiceURL,
		Authenticator: options.Authenticator,
	}

	err = core.ValidateStruct(options, "options")
	if err != nil {
		err = core.SDKErrorf(err, "", "invalid-global-options", common.GetComponentInfo())
		return
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

	service = &LogpushJobsApiV1{
		Service: baseService,
		Crn:     options.Crn,
		Dataset: options.Dataset,
		ZoneID:  options.ZoneID,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "logpushJobsApi" suitable for processing requests.
func (logpushJobsApi *LogpushJobsApiV1) Clone() *LogpushJobsApiV1 {
	if core.IsNil(logpushJobsApi) {
		return nil
	}
	clone := *logpushJobsApi
	clone.Service = logpushJobsApi.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (logpushJobsApi *LogpushJobsApiV1) SetServiceURL(url string) error {
	err := logpushJobsApi.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (logpushJobsApi *LogpushJobsApiV1) GetServiceURL() string {
	return logpushJobsApi.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (logpushJobsApi *LogpushJobsApiV1) SetDefaultHeaders(headers http.Header) {
	logpushJobsApi.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (logpushJobsApi *LogpushJobsApiV1) SetEnableGzipCompression(enableGzip bool) {
	logpushJobsApi.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (logpushJobsApi *LogpushJobsApiV1) GetEnableGzipCompression() bool {
	return logpushJobsApi.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (logpushJobsApi *LogpushJobsApiV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	logpushJobsApi.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (logpushJobsApi *LogpushJobsApiV1) DisableRetries() {
	logpushJobsApi.Service.DisableRetries()
}

// GetLogpushJobsV2 : List logpush jobs
// List configured logpush jobs for your domain.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobsV2(getLogpushJobsV2Options *GetLogpushJobsV2Options) (result *ListLogpushJobsResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.GetLogpushJobsV2WithContext(context.Background(), getLogpushJobsV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLogpushJobsV2WithContext is an alternate form of the GetLogpushJobsV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobsV2WithContext(ctx context.Context, getLogpushJobsV2Options *GetLogpushJobsV2Options) (result *ListLogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogpushJobsV2Options, "getLogpushJobsV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLogpushJobsV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogpushJobsV2")
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_logpush_jobs_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLogpushJobsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateLogpushJobV2 : Create a logpush jobs
// Create a new logpush job for the domain.
func (logpushJobsApi *LogpushJobsApiV1) CreateLogpushJobV2(createLogpushJobV2Options *CreateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.CreateLogpushJobV2WithContext(context.Background(), createLogpushJobV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateLogpushJobV2WithContext is an alternate form of the CreateLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) CreateLogpushJobV2WithContext(ctx context.Context, createLogpushJobV2Options *CreateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLogpushJobV2Options, "createLogpushJobV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createLogpushJobV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "CreateLogpushJobV2")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if createLogpushJobV2Options.CreateLogpushJobV2Request != nil {
		_, err = builder.SetBodyContentJSON(createLogpushJobV2Options.CreateLogpushJobV2Request)
		if err != nil {
			err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_logpush_job_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushJobV2 : Get a logpush job
// Get a logpush job  for a given zone.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobV2(getLogpushJobV2Options *GetLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.GetLogpushJobV2WithContext(context.Background(), getLogpushJobV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLogpushJobV2WithContext is an alternate form of the GetLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobV2WithContext(ctx context.Context, getLogpushJobV2Options *GetLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLogpushJobV2Options, "getLogpushJobV2Options cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getLogpushJobV2Options, "getLogpushJobV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id":  *getLogpushJobV2Options.JobID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLogpushJobV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogpushJobV2")
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_logpush_job_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateLogpushJobV2 : Update a logpush job
// Update an existing logpush job for a given zone.
func (logpushJobsApi *LogpushJobsApiV1) UpdateLogpushJobV2(updateLogpushJobV2Options *UpdateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.UpdateLogpushJobV2WithContext(context.Background(), updateLogpushJobV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateLogpushJobV2WithContext is an alternate form of the UpdateLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) UpdateLogpushJobV2WithContext(ctx context.Context, updateLogpushJobV2Options *UpdateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLogpushJobV2Options, "updateLogpushJobV2Options cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateLogpushJobV2Options, "updateLogpushJobV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id":  *updateLogpushJobV2Options.JobID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateLogpushJobV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "UpdateLogpushJobV2")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if updateLogpushJobV2Options.UpdateLogpushJobV2Request != nil {
		_, err = builder.SetBodyContentJSON(updateLogpushJobV2Options.UpdateLogpushJobV2Request)
		if err != nil {
			err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_logpush_job_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteLogpushJobV2 : Delete a logpush job
// Delete a logpush job for a zone.
func (logpushJobsApi *LogpushJobsApiV1) DeleteLogpushJobV2(deleteLogpushJobV2Options *DeleteLogpushJobV2Options) (result *DeleteLogpushJobResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.DeleteLogpushJobV2WithContext(context.Background(), deleteLogpushJobV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteLogpushJobV2WithContext is an alternate form of the DeleteLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) DeleteLogpushJobV2WithContext(ctx context.Context, deleteLogpushJobV2Options *DeleteLogpushJobV2Options) (result *DeleteLogpushJobResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLogpushJobV2Options, "deleteLogpushJobV2Options cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteLogpushJobV2Options, "deleteLogpushJobV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id":  *deleteLogpushJobV2Options.JobID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteLogpushJobV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "DeleteLogpushJobV2")
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_logpush_job_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteLogpushJobResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushOwnershipV2 : Get a new ownership challenge sent to your destination
// Get a new ownership challenge.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushOwnershipV2(getLogpushOwnershipV2Options *GetLogpushOwnershipV2Options) (result *OwnershipChallengeResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.GetLogpushOwnershipV2WithContext(context.Background(), getLogpushOwnershipV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLogpushOwnershipV2WithContext is an alternate form of the GetLogpushOwnershipV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushOwnershipV2WithContext(ctx context.Context, getLogpushOwnershipV2Options *GetLogpushOwnershipV2Options) (result *OwnershipChallengeResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogpushOwnershipV2Options, "getLogpushOwnershipV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/ownership`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLogpushOwnershipV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogpushOwnershipV2")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if getLogpushOwnershipV2Options.Cos != nil {
		body["cos"] = getLogpushOwnershipV2Options.Cos
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_logpush_ownership_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOwnershipChallengeResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ValidateLogpushOwnershipChallengeV2 : Validate ownership challenge of the destination
// Validate ownership challenge of the destination.
func (logpushJobsApi *LogpushJobsApiV1) ValidateLogpushOwnershipChallengeV2(validateLogpushOwnershipChallengeV2Options *ValidateLogpushOwnershipChallengeV2Options) (result *OwnershipChallengeValidateResult, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.ValidateLogpushOwnershipChallengeV2WithContext(context.Background(), validateLogpushOwnershipChallengeV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ValidateLogpushOwnershipChallengeV2WithContext is an alternate form of the ValidateLogpushOwnershipChallengeV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ValidateLogpushOwnershipChallengeV2WithContext(ctx context.Context, validateLogpushOwnershipChallengeV2Options *ValidateLogpushOwnershipChallengeV2Options) (result *OwnershipChallengeValidateResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(validateLogpushOwnershipChallengeV2Options, "validateLogpushOwnershipChallengeV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/ownership/validate`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range validateLogpushOwnershipChallengeV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "ValidateLogpushOwnershipChallengeV2")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if validateLogpushOwnershipChallengeV2Options.Cos != nil {
		body["cos"] = validateLogpushOwnershipChallengeV2Options.Cos
	}
	if validateLogpushOwnershipChallengeV2Options.OwnershipChallenge != nil {
		body["ownership_challenge"] = validateLogpushOwnershipChallengeV2Options.OwnershipChallenge
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "validate_logpush_ownership_challenge_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOwnershipChallengeValidateResult)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListFieldsForDatasetV2 : The list of all fields available for a dataset
// The list of all fields available for a dataset.
func (logpushJobsApi *LogpushJobsApiV1) ListFieldsForDatasetV2(listFieldsForDatasetV2Options *ListFieldsForDatasetV2Options) (result *ListFieldsResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.ListFieldsForDatasetV2WithContext(context.Background(), listFieldsForDatasetV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListFieldsForDatasetV2WithContext is an alternate form of the ListFieldsForDatasetV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ListFieldsForDatasetV2WithContext(ctx context.Context, listFieldsForDatasetV2Options *ListFieldsForDatasetV2Options) (result *ListFieldsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listFieldsForDatasetV2Options, "listFieldsForDatasetV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"dataset": *logpushJobsApi.Dataset,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/datasets/{dataset}/fields`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listFieldsForDatasetV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "ListFieldsForDatasetV2")
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_fields_for_dataset_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListFieldsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListLogpushJobsForDatasetV2 : List logpush jobs for dataset
// List configured logpush jobs for a dataset.
func (logpushJobsApi *LogpushJobsApiV1) ListLogpushJobsForDatasetV2(listLogpushJobsForDatasetV2Options *ListLogpushJobsForDatasetV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.ListLogpushJobsForDatasetV2WithContext(context.Background(), listLogpushJobsForDatasetV2Options)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListLogpushJobsForDatasetV2WithContext is an alternate form of the ListLogpushJobsForDatasetV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ListLogpushJobsForDatasetV2WithContext(ctx context.Context, listLogpushJobsForDatasetV2Options *ListLogpushJobsForDatasetV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listLogpushJobsForDatasetV2Options, "listLogpushJobsForDatasetV2Options")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"dataset": *logpushJobsApi.Dataset,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/datasets/{dataset}/jobs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listLogpushJobsForDatasetV2Options.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "ListLogpushJobsForDatasetV2")
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_logpush_jobs_for_dataset_v2", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetLogsRetention : Get log retention
// Get log retention setting for Logpull/Logpush on your domain.
func (logpushJobsApi *LogpushJobsApiV1) GetLogsRetention(getLogsRetentionOptions *GetLogsRetentionOptions) (result *LogRetentionResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.GetLogsRetentionWithContext(context.Background(), getLogsRetentionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetLogsRetentionWithContext is an alternate form of the GetLogsRetention method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogsRetentionWithContext(ctx context.Context, getLogsRetentionOptions *GetLogsRetentionOptions) (result *LogRetentionResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogsRetentionOptions, "getLogsRetentionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logs/retention`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getLogsRetentionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogsRetention")
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_logs_retention", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogRetentionResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateLogRetention : Update log retention
// Update log retention flag for Logpull/Logpush.
func (logpushJobsApi *LogpushJobsApiV1) CreateLogRetention(createLogRetentionOptions *CreateLogRetentionOptions) (result *LogRetentionResp, response *core.DetailedResponse, err error) {
	result, response, err = logpushJobsApi.CreateLogRetentionWithContext(context.Background(), createLogRetentionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateLogRetentionWithContext is an alternate form of the CreateLogRetention method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) CreateLogRetentionWithContext(ctx context.Context, createLogRetentionOptions *CreateLogRetentionOptions) (result *LogRetentionResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLogRetentionOptions, "createLogRetentionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"crn":     *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logs/retention`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createLogRetentionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "CreateLogRetention")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLogRetentionOptions.Flag != nil {
		body["flag"] = createLogRetentionOptions.Flag
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_log_retention", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogRetentionResp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0.0")
}

// CreateLogRetentionOptions : The CreateLogRetention options.
type CreateLogRetentionOptions struct {
	Flag *bool `json:"flag,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateLogRetentionOptions : Instantiate CreateLogRetentionOptions
func (*LogpushJobsApiV1) NewCreateLogRetentionOptions() *CreateLogRetentionOptions {
	return &CreateLogRetentionOptions{}
}

// SetFlag : Allow user to set Flag
func (_options *CreateLogRetentionOptions) SetFlag(flag bool) *CreateLogRetentionOptions {
	_options.Flag = core.BoolPtr(flag)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLogRetentionOptions) SetHeaders(param map[string]string) *CreateLogRetentionOptions {
	options.Headers = param
	return options
}

// CreateLogpushJobV2Options : The CreateLogpushJobV2 options.
type CreateLogpushJobV2Options struct {
	// Create logpush job body.
	CreateLogpushJobV2Request CreateLogpushJobV2RequestIntf `json:"CreateLogpushJobV2Request,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateLogpushJobV2Options : Instantiate CreateLogpushJobV2Options
func (*LogpushJobsApiV1) NewCreateLogpushJobV2Options() *CreateLogpushJobV2Options {
	return &CreateLogpushJobV2Options{}
}

// SetCreateLogpushJobV2Request : Allow user to set CreateLogpushJobV2Request
func (_options *CreateLogpushJobV2Options) SetCreateLogpushJobV2Request(createLogpushJobV2Request CreateLogpushJobV2RequestIntf) *CreateLogpushJobV2Options {
	_options.CreateLogpushJobV2Request = createLogpushJobV2Request
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLogpushJobV2Options) SetHeaders(param map[string]string) *CreateLogpushJobV2Options {
	options.Headers = param
	return options
}

// CreateLogpushJobV2Request : CreateLogpushJobV2Request struct
// Models which "extend" this model:
// - CreateLogpushJobV2RequestLogpushJobCosReq
// - CreateLogpushJobV2RequestLogpushJobLogdnaReq
// - CreateLogpushJobV2RequestLogpushJobIbmclReq
// - CreateLogpushJobV2RequestLogpushJobGenericReq
type CreateLogpushJobV2Request struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the COS bucket where the data will be pushed.
	Cos map[string]interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`

	// Information to identify the LogDNA instance the data will be pushed.
	Logdna map[string]interface{} `json:"logdna,omitempty"`

	// Required information to push logs to your Cloud Logs instance.
	Ibmcl *LogpushJobIbmclReqIbmcl `json:"ibmcl,omitempty"`

	// Uniquely identifies a resource where data will be pushed. Additional configuration parameters supported by the
	// destination may be included.
	DestinationConf *string `json:"destination_conf,omitempty"`
}

// Constants associated with the CreateLogpushJobV2Request.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2Request_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2Request_Dataset_HttpRequests   = "http_requests"
	CreateLogpushJobV2Request_Dataset_RangeEvents    = "range_events"
)

// Constants associated with the CreateLogpushJobV2Request.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2Request_Frequency_High = "high"
	CreateLogpushJobV2Request_Frequency_Low  = "low"
)

func (*CreateLogpushJobV2Request) isaCreateLogpushJobV2Request() bool {
	return true
}

type CreateLogpushJobV2RequestIntf interface {
	isaCreateLogpushJobV2Request() bool
}

// UnmarshalCreateLogpushJobV2Request unmarshals an instance of CreateLogpushJobV2Request from the specified map of raw messages.
func UnmarshalCreateLogpushJobV2Request(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateLogpushJobV2Request)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		err = core.SDKErrorf(err, "", "ownership_challenge-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		err = core.SDKErrorf(err, "", "logdna-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ibmcl", &obj.Ibmcl, UnmarshalLogpushJobIbmclReqIbmcl)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmcl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destination_conf", &obj.DestinationConf)
	if err != nil {
		err = core.SDKErrorf(err, "", "destination_conf-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteLogpushJobV2Options : The DeleteLogpushJobV2 options.
type DeleteLogpushJobV2Options struct {
	// logpush job identifier.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteLogpushJobV2Options : Instantiate DeleteLogpushJobV2Options
func (*LogpushJobsApiV1) NewDeleteLogpushJobV2Options(jobID string) *DeleteLogpushJobV2Options {
	return &DeleteLogpushJobV2Options{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *DeleteLogpushJobV2Options) SetJobID(jobID string) *DeleteLogpushJobV2Options {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLogpushJobV2Options) SetHeaders(param map[string]string) *DeleteLogpushJobV2Options {
	options.Headers = param
	return options
}

// GetLogpushJobV2Options : The GetLogpushJobV2 options.
type GetLogpushJobV2Options struct {
	// logpush job identifier.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLogpushJobV2Options : Instantiate GetLogpushJobV2Options
func (*LogpushJobsApiV1) NewGetLogpushJobV2Options(jobID string) *GetLogpushJobV2Options {
	return &GetLogpushJobV2Options{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *GetLogpushJobV2Options) SetJobID(jobID string) *GetLogpushJobV2Options {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushJobV2Options) SetHeaders(param map[string]string) *GetLogpushJobV2Options {
	options.Headers = param
	return options
}

// GetLogpushJobsV2Options : The GetLogpushJobsV2 options.
type GetLogpushJobsV2Options struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLogpushJobsV2Options : Instantiate GetLogpushJobsV2Options
func (*LogpushJobsApiV1) NewGetLogpushJobsV2Options() *GetLogpushJobsV2Options {
	return &GetLogpushJobsV2Options{}
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushJobsV2Options) SetHeaders(param map[string]string) *GetLogpushJobsV2Options {
	options.Headers = param
	return options
}

// GetLogpushOwnershipV2Options : The GetLogpushOwnershipV2 options.
type GetLogpushOwnershipV2Options struct {
	// Information to identify the COS bucket where the data will be pushed.
	Cos map[string]interface{} `json:"cos,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLogpushOwnershipV2Options : Instantiate GetLogpushOwnershipV2Options
func (*LogpushJobsApiV1) NewGetLogpushOwnershipV2Options() *GetLogpushOwnershipV2Options {
	return &GetLogpushOwnershipV2Options{}
}

// SetCos : Allow user to set Cos
func (_options *GetLogpushOwnershipV2Options) SetCos(cos map[string]interface{}) *GetLogpushOwnershipV2Options {
	_options.Cos = cos
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushOwnershipV2Options) SetHeaders(param map[string]string) *GetLogpushOwnershipV2Options {
	options.Headers = param
	return options
}

// GetLogsRetentionOptions : The GetLogsRetention options.
type GetLogsRetentionOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetLogsRetentionOptions : Instantiate GetLogsRetentionOptions
func (*LogpushJobsApiV1) NewGetLogsRetentionOptions() *GetLogsRetentionOptions {
	return &GetLogsRetentionOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetLogsRetentionOptions) SetHeaders(param map[string]string) *GetLogsRetentionOptions {
	options.Headers = param
	return options
}

// ListFieldsForDatasetV2Options : The ListFieldsForDatasetV2 options.
type ListFieldsForDatasetV2Options struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListFieldsForDatasetV2Options : Instantiate ListFieldsForDatasetV2Options
func (*LogpushJobsApiV1) NewListFieldsForDatasetV2Options() *ListFieldsForDatasetV2Options {
	return &ListFieldsForDatasetV2Options{}
}

// SetHeaders : Allow user to set Headers
func (options *ListFieldsForDatasetV2Options) SetHeaders(param map[string]string) *ListFieldsForDatasetV2Options {
	options.Headers = param
	return options
}

// ListLogpushJobsForDatasetV2Options : The ListLogpushJobsForDatasetV2 options.
type ListLogpushJobsForDatasetV2Options struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListLogpushJobsForDatasetV2Options : Instantiate ListLogpushJobsForDatasetV2Options
func (*LogpushJobsApiV1) NewListLogpushJobsForDatasetV2Options() *ListLogpushJobsForDatasetV2Options {
	return &ListLogpushJobsForDatasetV2Options{}
}

// SetHeaders : Allow user to set Headers
func (options *ListLogpushJobsForDatasetV2Options) SetHeaders(param map[string]string) *ListLogpushJobsForDatasetV2Options {
	options.Headers = param
	return options
}

// LogRetentionRespResult : LogRetentionRespResult struct
type LogRetentionRespResult struct {
	Flag *bool `json:"flag,omitempty"`
}

// UnmarshalLogRetentionRespResult unmarshals an instance of LogRetentionRespResult from the specified map of raw messages.
func UnmarshalLogRetentionRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogRetentionRespResult)
	err = core.UnmarshalPrimitive(m, "flag", &obj.Flag)
	if err != nil {
		err = core.SDKErrorf(err, "", "flag-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogpushJobIbmclReqIbmcl : Required information to push logs to your Cloud Logs instance.
type LogpushJobIbmclReqIbmcl struct {
	// GUID of the IBM Cloud Logs instance where you want to send logs.
	InstanceID *string `json:"instance_id" validate:"required"`

	// Region where the IBM Cloud Logs instance is located.
	Region *string `json:"region" validate:"required"`

	// IBM Cloud API key used to generate a token for pushing to your Cloud Logs instance.
	ApiKey *string `json:"api_key" validate:"required"`
}

// NewLogpushJobIbmclReqIbmcl : Instantiate LogpushJobIbmclReqIbmcl (Generic Model Constructor)
func (*LogpushJobsApiV1) NewLogpushJobIbmclReqIbmcl(instanceID string, region string, apiKey string) (_model *LogpushJobIbmclReqIbmcl, err error) {
	_model = &LogpushJobIbmclReqIbmcl{
		InstanceID: core.StringPtr(instanceID),
		Region:     core.StringPtr(region),
		ApiKey:     core.StringPtr(apiKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalLogpushJobIbmclReqIbmcl unmarshals an instance of LogpushJobIbmclReqIbmcl from the specified map of raw messages.
func UnmarshalLogpushJobIbmclReqIbmcl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogpushJobIbmclReqIbmcl)
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogpushJobsUpdateIbmclReqIbmcl : Required information to push logs to your Cloud Logs instance.
type LogpushJobsUpdateIbmclReqIbmcl struct {
	// GUID of the IBM Cloud Logs instance where you want to send logs.
	InstanceID *string `json:"instance_id,omitempty"`

	// Region where the IBM Cloud Logs instance is located.
	Region *string `json:"region,omitempty"`

	// IBM Cloud API key used to generate a token for pushing to your Cloud Logs instance.
	ApiKey *string `json:"api_key,omitempty"`
}

// UnmarshalLogpushJobsUpdateIbmclReqIbmcl unmarshals an instance of LogpushJobsUpdateIbmclReqIbmcl from the specified map of raw messages.
func UnmarshalLogpushJobsUpdateIbmclReqIbmcl(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogpushJobsUpdateIbmclReqIbmcl)
	err = core.UnmarshalPrimitive(m, "instance_id", &obj.InstanceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "instance_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key", &obj.ApiKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateLogpushJobV2Options : The UpdateLogpushJobV2 options.
type UpdateLogpushJobV2Options struct {
	// logpush job identifier.
	JobID *string `json:"job_id" validate:"required,ne="`

	// Update logpush job.
	UpdateLogpushJobV2Request UpdateLogpushJobV2RequestIntf `json:"UpdateLogpushJobV2Request,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateLogpushJobV2Options : Instantiate UpdateLogpushJobV2Options
func (*LogpushJobsApiV1) NewUpdateLogpushJobV2Options(jobID string) *UpdateLogpushJobV2Options {
	return &UpdateLogpushJobV2Options{
		JobID: core.StringPtr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *UpdateLogpushJobV2Options) SetJobID(jobID string) *UpdateLogpushJobV2Options {
	_options.JobID = core.StringPtr(jobID)
	return _options
}

// SetUpdateLogpushJobV2Request : Allow user to set UpdateLogpushJobV2Request
func (_options *UpdateLogpushJobV2Options) SetUpdateLogpushJobV2Request(updateLogpushJobV2Request UpdateLogpushJobV2RequestIntf) *UpdateLogpushJobV2Options {
	_options.UpdateLogpushJobV2Request = updateLogpushJobV2Request
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLogpushJobV2Options) SetHeaders(param map[string]string) *UpdateLogpushJobV2Options {
	options.Headers = param
	return options
}

// UpdateLogpushJobV2Request : UpdateLogpushJobV2Request struct
// Models which "extend" this model:
// - UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq
// - UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq
// - UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq
// - UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq
type UpdateLogpushJobV2Request struct {
	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the COS bucket where the data will be pushed.
	Cos map[string]interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`

	// Information to identify the LogDNA instance the data will be pushed.
	Logdna map[string]interface{} `json:"logdna,omitempty"`

	// Required information to push logs to your Cloud Logs instance.
	Ibmcl *LogpushJobsUpdateIbmclReqIbmcl `json:"ibmcl,omitempty"`

	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Uniquely identifies a resource where data will be pushed. Additional configuration parameters supported by the
	// destination may be included.
	DestinationConf *string `json:"destination_conf,omitempty"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2Request.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2Request_Frequency_High = "high"
	UpdateLogpushJobV2Request_Frequency_Low  = "low"
)

// Constants associated with the UpdateLogpushJobV2Request.Dataset property.
// Dataset to be pulled.
const (
	UpdateLogpushJobV2Request_Dataset_FirewallEvents = "firewall_events"
	UpdateLogpushJobV2Request_Dataset_HttpRequests   = "http_requests"
	UpdateLogpushJobV2Request_Dataset_RangeEvents    = "range_events"
)

func (*UpdateLogpushJobV2Request) isaUpdateLogpushJobV2Request() bool {
	return true
}

type UpdateLogpushJobV2RequestIntf interface {
	isaUpdateLogpushJobV2Request() bool
}

// UnmarshalUpdateLogpushJobV2Request unmarshals an instance of UpdateLogpushJobV2Request from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2Request(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2Request)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		err = core.SDKErrorf(err, "", "ownership_challenge-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		err = core.SDKErrorf(err, "", "logdna-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ibmcl", &obj.Ibmcl, UnmarshalLogpushJobsUpdateIbmclReqIbmcl)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmcl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destination_conf", &obj.DestinationConf)
	if err != nil {
		err = core.SDKErrorf(err, "", "destination_conf-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ValidateLogpushOwnershipChallengeV2Options : The ValidateLogpushOwnershipChallengeV2 options.
type ValidateLogpushOwnershipChallengeV2Options struct {
	// Information to identify the COS bucket where the data will be pushed.
	Cos map[string]interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewValidateLogpushOwnershipChallengeV2Options : Instantiate ValidateLogpushOwnershipChallengeV2Options
func (*LogpushJobsApiV1) NewValidateLogpushOwnershipChallengeV2Options() *ValidateLogpushOwnershipChallengeV2Options {
	return &ValidateLogpushOwnershipChallengeV2Options{}
}

// SetCos : Allow user to set Cos
func (_options *ValidateLogpushOwnershipChallengeV2Options) SetCos(cos map[string]interface{}) *ValidateLogpushOwnershipChallengeV2Options {
	_options.Cos = cos
	return _options
}

// SetOwnershipChallenge : Allow user to set OwnershipChallenge
func (_options *ValidateLogpushOwnershipChallengeV2Options) SetOwnershipChallenge(ownershipChallenge string) *ValidateLogpushOwnershipChallengeV2Options {
	_options.OwnershipChallenge = core.StringPtr(ownershipChallenge)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ValidateLogpushOwnershipChallengeV2Options) SetHeaders(param map[string]string) *ValidateLogpushOwnershipChallengeV2Options {
	options.Headers = param
	return options
}

// DeleteLogpushJobResp : delete logpush job response.
type DeleteLogpushJobResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result map[string]interface{} `json:"result" validate:"required"`
}

// UnmarshalDeleteLogpushJobResp unmarshals an instance of DeleteLogpushJobResp from the specified map of raw messages.
func UnmarshalDeleteLogpushJobResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLogpushJobResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListFieldsResp : list fields response.
type ListFieldsResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result map[string]interface{} `json:"result,omitempty"`
}

// UnmarshalListFieldsResp unmarshals an instance of ListFieldsResp from the specified map of raw messages.
func UnmarshalListFieldsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListFieldsResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListLogpushJobsResp : List Logpush Jobs Response.
type ListLogpushJobsResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result []LogpushJobPack `json:"result" validate:"required"`
}

// UnmarshalListLogpushJobsResp unmarshals an instance of ListLogpushJobsResp from the specified map of raw messages.
func UnmarshalListLogpushJobsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLogpushJobsResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLogpushJobPack)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogRetentionResp : log retention result.
type LogRetentionResp struct {
	Result *LogRetentionRespResult `json:"result,omitempty"`

	// success response.
	Success *bool `json:"success,omitempty"`

	// errors.
	Errors [][]string `json:"errors,omitempty"`

	// messages.
	Messages [][]string `json:"messages,omitempty"`
}

// UnmarshalLogRetentionResp unmarshals an instance of LogRetentionResp from the specified map of raw messages.
func UnmarshalLogRetentionResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogRetentionResp)
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLogRetentionRespResult)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogpushJobPack : logpush job pack.
type LogpushJobPack struct {
	// Logpush Job ID.
	ID *int64 `json:"id" validate:"required"`

	// Logpush Job Name.
	Name *string `json:"name" validate:"required"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled" validate:"required"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset" validate:"required"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency" validate:"required"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options" validate:"required"`

	// Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.
	DestinationConf *string `json:"destination_conf" validate:"required"`

	// Records the last time for which logs have been successfully pushed.
	LastComplete *string `json:"last_complete" validate:"required"`

	// Records the last time the job failed.
	LastError *string `json:"last_error" validate:"required"`

	// The last failure.
	ErrorMessage *string `json:"error_message" validate:"required"`
}

// UnmarshalLogpushJobPack unmarshals an instance of LogpushJobPack from the specified map of raw messages.
func UnmarshalLogpushJobPack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogpushJobPack)
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
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destination_conf", &obj.DestinationConf)
	if err != nil {
		err = core.SDKErrorf(err, "", "destination_conf-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_complete", &obj.LastComplete)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_complete-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_error", &obj.LastError)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_error-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "error_message", &obj.ErrorMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "error_message-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogpushJobsResp : logpush job response.
type LogpushJobsResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// logpush job pack.
	Result *LogpushJobPack `json:"result" validate:"required"`
}

// UnmarshalLogpushJobsResp unmarshals an instance of LogpushJobsResp from the specified map of raw messages.
func UnmarshalLogpushJobsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogpushJobsResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLogpushJobPack)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OwnershipChallengeResp : Get Logpush Ownership Challenge Response.
type OwnershipChallengeResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// ownership challenge result.
	Result *OwnershipChallengeResult `json:"result" validate:"required"`
}

// UnmarshalOwnershipChallengeResp unmarshals an instance of OwnershipChallengeResp from the specified map of raw messages.
func UnmarshalOwnershipChallengeResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OwnershipChallengeResp)
	err = core.UnmarshalPrimitive(m, "success", &obj.Success)
	if err != nil {
		err = core.SDKErrorf(err, "", "success-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "errors", &obj.Errors)
	if err != nil {
		err = core.SDKErrorf(err, "", "errors-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalOwnershipChallengeResult)
	if err != nil {
		err = core.SDKErrorf(err, "", "result-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OwnershipChallengeResult : ownership challenge result.
type OwnershipChallengeResult struct {
	// file name.
	Filename *string `json:"filename" validate:"required"`

	// valid.
	Valid *bool `json:"valid" validate:"required"`

	// message.
	Messages *string `json:"messages,omitempty"`
}

// UnmarshalOwnershipChallengeResult unmarshals an instance of OwnershipChallengeResult from the specified map of raw messages.
func UnmarshalOwnershipChallengeResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OwnershipChallengeResult)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		err = core.SDKErrorf(err, "", "filename-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "valid", &obj.Valid)
	if err != nil {
		err = core.SDKErrorf(err, "", "valid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "messages", &obj.Messages)
	if err != nil {
		err = core.SDKErrorf(err, "", "messages-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OwnershipChallengeValidateResult : ownership challenge validate result.
type OwnershipChallengeValidateResult struct {
	// valid.
	Valid *bool `json:"valid" validate:"required"`
}

// UnmarshalOwnershipChallengeValidateResult unmarshals an instance of OwnershipChallengeValidateResult from the specified map of raw messages.
func UnmarshalOwnershipChallengeValidateResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OwnershipChallengeValidateResult)
	err = core.UnmarshalPrimitive(m, "valid", &obj.Valid)
	if err != nil {
		err = core.SDKErrorf(err, "", "valid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateLogpushJobV2RequestLogpushJobCosReq : Create COS logpush job input.
// This model "extends" CreateLogpushJobV2Request
type CreateLogpushJobV2RequestLogpushJobCosReq struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the COS bucket where the data will be pushed.
	Cos map[string]interface{} `json:"cos" validate:"required"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge" validate:"required"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the CreateLogpushJobV2RequestLogpushJobCosReq.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2RequestLogpushJobCosReq_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2RequestLogpushJobCosReq_Dataset_HttpRequests   = "http_requests"
	CreateLogpushJobV2RequestLogpushJobCosReq_Dataset_RangeEvents    = "range_events"
)

// Constants associated with the CreateLogpushJobV2RequestLogpushJobCosReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2RequestLogpushJobCosReq_Frequency_High = "high"
	CreateLogpushJobV2RequestLogpushJobCosReq_Frequency_Low  = "low"
)

// NewCreateLogpushJobV2RequestLogpushJobCosReq : Instantiate CreateLogpushJobV2RequestLogpushJobCosReq (Generic Model Constructor)
func (*LogpushJobsApiV1) NewCreateLogpushJobV2RequestLogpushJobCosReq(cos map[string]interface{}, ownershipChallenge string) (_model *CreateLogpushJobV2RequestLogpushJobCosReq, err error) {
	_model = &CreateLogpushJobV2RequestLogpushJobCosReq{
		Cos:                cos,
		OwnershipChallenge: core.StringPtr(ownershipChallenge),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*CreateLogpushJobV2RequestLogpushJobCosReq) isaCreateLogpushJobV2Request() bool {
	return true
}

// UnmarshalCreateLogpushJobV2RequestLogpushJobCosReq unmarshals an instance of CreateLogpushJobV2RequestLogpushJobCosReq from the specified map of raw messages.
func UnmarshalCreateLogpushJobV2RequestLogpushJobCosReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateLogpushJobV2RequestLogpushJobCosReq)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		err = core.SDKErrorf(err, "", "ownership_challenge-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateLogpushJobV2RequestLogpushJobGenericReq : Create logpush job for a generic destination.
// This model "extends" CreateLogpushJobV2Request
type CreateLogpushJobV2RequestLogpushJobGenericReq struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job is enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Uniquely identifies a resource where data will be pushed. Additional configuration parameters supported by the
	// destination may be included.
	DestinationConf *string `json:"destination_conf" validate:"required"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the CreateLogpushJobV2RequestLogpushJobGenericReq.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2RequestLogpushJobGenericReq_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2RequestLogpushJobGenericReq_Dataset_HttpRequests   = "http_requests"
	CreateLogpushJobV2RequestLogpushJobGenericReq_Dataset_RangeEvents    = "range_events"
)

// Constants associated with the CreateLogpushJobV2RequestLogpushJobGenericReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2RequestLogpushJobGenericReq_Frequency_High = "high"
	CreateLogpushJobV2RequestLogpushJobGenericReq_Frequency_Low  = "low"
)

// NewCreateLogpushJobV2RequestLogpushJobGenericReq : Instantiate CreateLogpushJobV2RequestLogpushJobGenericReq (Generic Model Constructor)
func (*LogpushJobsApiV1) NewCreateLogpushJobV2RequestLogpushJobGenericReq(destinationConf string) (_model *CreateLogpushJobV2RequestLogpushJobGenericReq, err error) {
	_model = &CreateLogpushJobV2RequestLogpushJobGenericReq{
		DestinationConf: core.StringPtr(destinationConf),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*CreateLogpushJobV2RequestLogpushJobGenericReq) isaCreateLogpushJobV2Request() bool {
	return true
}

// UnmarshalCreateLogpushJobV2RequestLogpushJobGenericReq unmarshals an instance of CreateLogpushJobV2RequestLogpushJobGenericReq from the specified map of raw messages.
func UnmarshalCreateLogpushJobV2RequestLogpushJobGenericReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateLogpushJobV2RequestLogpushJobGenericReq)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destination_conf", &obj.DestinationConf)
	if err != nil {
		err = core.SDKErrorf(err, "", "destination_conf-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateLogpushJobV2RequestLogpushJobIbmclReq : Create IBM Cloud Logs logpush job input.
// This model "extends" CreateLogpushJobV2Request
type CreateLogpushJobV2RequestLogpushJobIbmclReq struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job is enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Required information to push logs to your Cloud Logs instance.
	Ibmcl *LogpushJobIbmclReqIbmcl `json:"ibmcl" validate:"required"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the CreateLogpushJobV2RequestLogpushJobIbmclReq.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2RequestLogpushJobIbmclReq_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2RequestLogpushJobIbmclReq_Dataset_HttpRequests   = "http_requests"
	CreateLogpushJobV2RequestLogpushJobIbmclReq_Dataset_RangeEvents    = "range_events"
)

// Constants associated with the CreateLogpushJobV2RequestLogpushJobIbmclReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2RequestLogpushJobIbmclReq_Frequency_High = "high"
	CreateLogpushJobV2RequestLogpushJobIbmclReq_Frequency_Low  = "low"
)

// NewCreateLogpushJobV2RequestLogpushJobIbmclReq : Instantiate CreateLogpushJobV2RequestLogpushJobIbmclReq (Generic Model Constructor)
func (*LogpushJobsApiV1) NewCreateLogpushJobV2RequestLogpushJobIbmclReq(ibmcl *LogpushJobIbmclReqIbmcl) (_model *CreateLogpushJobV2RequestLogpushJobIbmclReq, err error) {
	_model = &CreateLogpushJobV2RequestLogpushJobIbmclReq{
		Ibmcl: ibmcl,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*CreateLogpushJobV2RequestLogpushJobIbmclReq) isaCreateLogpushJobV2Request() bool {
	return true
}

// UnmarshalCreateLogpushJobV2RequestLogpushJobIbmclReq unmarshals an instance of CreateLogpushJobV2RequestLogpushJobIbmclReq from the specified map of raw messages.
func UnmarshalCreateLogpushJobV2RequestLogpushJobIbmclReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateLogpushJobV2RequestLogpushJobIbmclReq)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ibmcl", &obj.Ibmcl, UnmarshalLogpushJobIbmclReqIbmcl)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmcl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateLogpushJobV2RequestLogpushJobLogdnaReq : Create LogDNA logpush job input.
// This model "extends" CreateLogpushJobV2Request
type CreateLogpushJobV2RequestLogpushJobLogdnaReq struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the LogDNA instance the data will be pushed.
	Logdna map[string]interface{} `json:"logdna" validate:"required"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the CreateLogpushJobV2RequestLogpushJobLogdnaReq.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Dataset_HttpRequests   = "http_requests"
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Dataset_RangeEvents    = "range_events"
)

// Constants associated with the CreateLogpushJobV2RequestLogpushJobLogdnaReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Frequency_High = "high"
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Frequency_Low  = "low"
)

// NewCreateLogpushJobV2RequestLogpushJobLogdnaReq : Instantiate CreateLogpushJobV2RequestLogpushJobLogdnaReq (Generic Model Constructor)
func (*LogpushJobsApiV1) NewCreateLogpushJobV2RequestLogpushJobLogdnaReq(logdna map[string]interface{}) (_model *CreateLogpushJobV2RequestLogpushJobLogdnaReq, err error) {
	_model = &CreateLogpushJobV2RequestLogpushJobLogdnaReq{
		Logdna: logdna,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*CreateLogpushJobV2RequestLogpushJobLogdnaReq) isaCreateLogpushJobV2Request() bool {
	return true
}

// UnmarshalCreateLogpushJobV2RequestLogpushJobLogdnaReq unmarshals an instance of CreateLogpushJobV2RequestLogpushJobLogdnaReq from the specified map of raw messages.
func UnmarshalCreateLogpushJobV2RequestLogpushJobLogdnaReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateLogpushJobV2RequestLogpushJobLogdnaReq)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		err = core.SDKErrorf(err, "", "logdna-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq : Update COS logpush job input.
// This model "extends" UpdateLogpushJobV2Request
type UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq struct {
	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the COS bucket where the data will be pushed.
	Cos map[string]interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq_Frequency_High = "high"
	UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq_Frequency_Low  = "low"
)

func (*UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq) isaUpdateLogpushJobV2Request() bool {
	return true
}

// UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateCosReq unmarshals an instance of UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateCosReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		err = core.SDKErrorf(err, "", "cos-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		err = core.SDKErrorf(err, "", "ownership_challenge-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq : Create logpush job for a generic destination.
// This model "extends" UpdateLogpushJobV2Request
type UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job is enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Uniquely identifies a resource where data will be pushed. Additional configuration parameters supported by the
	// destination may be included.
	DestinationConf *string `json:"destination_conf,omitempty"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq.Dataset property.
// Dataset to be pulled.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq_Dataset_FirewallEvents = "firewall_events"
	UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq_Dataset_HttpRequests   = "http_requests"
	UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq_Dataset_RangeEvents    = "range_events"
)

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq_Frequency_High = "high"
	UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq_Frequency_Low  = "low"
)

func (*UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq) isaUpdateLogpushJobV2Request() bool {
	return true
}

// UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq unmarshals an instance of UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2RequestLogpushJobsUpdateGenericReq)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "destination_conf", &obj.DestinationConf)
	if err != nil {
		err = core.SDKErrorf(err, "", "destination_conf-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		err = core.SDKErrorf(err, "", "dataset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq : Update IBM Cloud Logs logpush job input.
// This model "extends" UpdateLogpushJobV2Request
type UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq struct {
	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Required information to push logs to your Cloud Logs instance.
	Ibmcl *LogpushJobsUpdateIbmclReqIbmcl `json:"ibmcl,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq_Frequency_High = "high"
	UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq_Frequency_Low  = "low"
)

func (*UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq) isaUpdateLogpushJobV2Request() bool {
	return true
}

// UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq unmarshals an instance of UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2RequestLogpushJobsUpdateIbmclReq)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "ibmcl", &obj.Ibmcl, UnmarshalLogpushJobsUpdateIbmclReqIbmcl)
	if err != nil {
		err = core.SDKErrorf(err, "", "ibmcl-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq : Update LogDNA logpush job input.
// This model "extends" UpdateLogpushJobV2Request
type UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq struct {
	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the LogDNA instance the data will be pushed.
	Logdna map[string]interface{} `json:"logdna,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq_Frequency_High = "high"
	UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq_Frequency_Low  = "low"
)

func (*UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq) isaUpdateLogpushJobV2Request() bool {
	return true
}

// UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq unmarshals an instance of UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		err = core.SDKErrorf(err, "", "logpull_options-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		err = core.SDKErrorf(err, "", "logdna-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		err = core.SDKErrorf(err, "", "frequency-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
