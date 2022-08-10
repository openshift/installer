/**
 * (C) Copyright IBM Corp. 2022.
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
 * IBM OpenAPI SDK Code Generator Version: 3.43.4-432d779b-20220119-173927
 */

// Package logpushjobsapiv1 : Operations and models for the LogpushJobsApiV1 service
package logpushjobsapiv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/networking-go-sdk/common"
)

// LogpushJobsApiV1 : CIS Loupush Jobs
type LogpushJobsApiV1 struct {
	Service *core.BaseService

	// Full url-encoded CRN of the service instance.
	Crn *string

	// The domain id.
	ZoneID *string

	// The dataset.
	Dataset *string
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

	// Full url-encoded CRN of the service instance.
	Crn *string `validate:"required"`

	// The domain id.
	ZoneID *string `validate:"required"`

	// The dataset.
	Dataset *string `validate:"required"`
}

// NewLogpushJobsApiV1UsingExternalConfig : constructs an instance of LogpushJobsApiV1 with passed in options and external configuration.
func NewLogpushJobsApiV1UsingExternalConfig(options *LogpushJobsApiV1Options) (logpushJobsApi *LogpushJobsApiV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	logpushJobsApi, err = NewLogpushJobsApiV1(options)
	if err != nil {
		return
	}

	err = logpushJobsApi.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = logpushJobsApi.Service.SetServiceURL(options.URL)
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

	service = &LogpushJobsApiV1{
		Service: baseService,
		Crn: options.Crn,
		ZoneID: options.ZoneID,
		Dataset: options.Dataset,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
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
	return logpushJobsApi.Service.SetServiceURL(url)
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

// GetLogpushJobs : List logpush jobs
// List configured logpush jobs for your domain.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobs(getLogpushJobsOptions *GetLogpushJobsOptions) (result *ListLogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.GetLogpushJobsWithContext(context.Background(), getLogpushJobsOptions)
}

// GetLogpushJobsWithContext is an alternate form of the GetLogpushJobs method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobsWithContext(ctx context.Context, getLogpushJobsOptions *GetLogpushJobsOptions) (result *ListLogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogpushJobsOptions, "getLogpushJobsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/jobs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLogpushJobsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogpushJobs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateLogpushJob : Create a logpush jobs
// Create a new logpush job for the domain.
func (logpushJobsApi *LogpushJobsApiV1) CreateLogpushJob(createLogpushJobOptions *CreateLogpushJobOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.CreateLogpushJobWithContext(context.Background(), createLogpushJobOptions)
}

// CreateLogpushJobWithContext is an alternate form of the CreateLogpushJob method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) CreateLogpushJobWithContext(ctx context.Context, createLogpushJobOptions *CreateLogpushJobOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLogpushJobOptions, "createLogpushJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/jobs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createLogpushJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "CreateLogpushJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLogpushJobOptions.Name != nil {
		body["name"] = createLogpushJobOptions.Name
	}
	if createLogpushJobOptions.Enabled != nil {
		body["enabled"] = createLogpushJobOptions.Enabled
	}
	if createLogpushJobOptions.LogpullOptions != nil {
		body["logpull_options"] = createLogpushJobOptions.LogpullOptions
	}
	if createLogpushJobOptions.DestinationConf != nil {
		body["destination_conf"] = createLogpushJobOptions.DestinationConf
	}
	if createLogpushJobOptions.OwnershipChallenge != nil {
		body["ownership_challenge"] = createLogpushJobOptions.OwnershipChallenge
	}
	if createLogpushJobOptions.Dataset != nil {
		body["dataset"] = createLogpushJobOptions.Dataset
	}
	if createLogpushJobOptions.Frequency != nil {
		body["frequency"] = createLogpushJobOptions.Frequency
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushJob : Get a logpush job
// Get a logpush job  for a given zone.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJob(getLogpushJobOptions *GetLogpushJobOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.GetLogpushJobWithContext(context.Background(), getLogpushJobOptions)
}

// GetLogpushJobWithContext is an alternate form of the GetLogpushJob method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobWithContext(ctx context.Context, getLogpushJobOptions *GetLogpushJobOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLogpushJobOptions, "getLogpushJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLogpushJobOptions, "getLogpushJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id": fmt.Sprint(*getLogpushJobOptions.JobID),
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLogpushJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogpushJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateLogpushJob : Update a logpush job
// Update an existing logpush job for a given zone.
func (logpushJobsApi *LogpushJobsApiV1) UpdateLogpushJob(updateLogpushJobOptions *UpdateLogpushJobOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.UpdateLogpushJobWithContext(context.Background(), updateLogpushJobOptions)
}

// UpdateLogpushJobWithContext is an alternate form of the UpdateLogpushJob method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) UpdateLogpushJobWithContext(ctx context.Context, updateLogpushJobOptions *UpdateLogpushJobOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLogpushJobOptions, "updateLogpushJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateLogpushJobOptions, "updateLogpushJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id": fmt.Sprint(*updateLogpushJobOptions.JobID),
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateLogpushJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "UpdateLogpushJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateLogpushJobOptions.Enabled != nil {
		body["enabled"] = updateLogpushJobOptions.Enabled
	}
	if updateLogpushJobOptions.LogpullOptions != nil {
		body["logpull_options"] = updateLogpushJobOptions.LogpullOptions
	}
	if updateLogpushJobOptions.DestinationConf != nil {
		body["destination_conf"] = updateLogpushJobOptions.DestinationConf
	}
	if updateLogpushJobOptions.OwnershipChallenge != nil {
		body["ownership_challenge"] = updateLogpushJobOptions.OwnershipChallenge
	}
	if updateLogpushJobOptions.Frequency != nil {
		body["frequency"] = updateLogpushJobOptions.Frequency
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteLogpushJob : Delete a logpush job
// Delete a logpush job for a zone.
func (logpushJobsApi *LogpushJobsApiV1) DeleteLogpushJob(deleteLogpushJobOptions *DeleteLogpushJobOptions) (result *DeleteLogpushJobResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.DeleteLogpushJobWithContext(context.Background(), deleteLogpushJobOptions)
}

// DeleteLogpushJobWithContext is an alternate form of the DeleteLogpushJob method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) DeleteLogpushJobWithContext(ctx context.Context, deleteLogpushJobOptions *DeleteLogpushJobOptions) (result *DeleteLogpushJobResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLogpushJobOptions, "deleteLogpushJobOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLogpushJobOptions, "deleteLogpushJobOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id": fmt.Sprint(*deleteLogpushJobOptions.JobID),
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteLogpushJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "DeleteLogpushJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteLogpushJobResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListFieldsForDataset : List all fields available for a dataset
// The list of all fields available for a given dataset.
func (logpushJobsApi *LogpushJobsApiV1) ListFieldsForDataset(listFieldsForDatasetOptions *ListFieldsForDatasetOptions) (result *ListFieldsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.ListFieldsForDatasetWithContext(context.Background(), listFieldsForDatasetOptions)
}

// ListFieldsForDatasetWithContext is an alternate form of the ListFieldsForDataset method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ListFieldsForDatasetWithContext(ctx context.Context, listFieldsForDatasetOptions *ListFieldsForDatasetOptions) (result *ListFieldsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listFieldsForDatasetOptions, "listFieldsForDatasetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"dataset": *logpushJobsApi.Dataset,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/datasets/{dataset}/fields`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listFieldsForDatasetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "ListFieldsForDataset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListFieldsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListLogpushJobsForDataset : List logpush jobs for a dataset
// List configured logpush jobs for a dataset.
func (logpushJobsApi *LogpushJobsApiV1) ListLogpushJobsForDataset(listLogpushJobsForDatasetOptions *ListLogpushJobsForDatasetOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.ListLogpushJobsForDatasetWithContext(context.Background(), listLogpushJobsForDatasetOptions)
}

// ListLogpushJobsForDatasetWithContext is an alternate form of the ListLogpushJobsForDataset method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ListLogpushJobsForDatasetWithContext(ctx context.Context, listLogpushJobsForDatasetOptions *ListLogpushJobsForDatasetOptions) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listLogpushJobsForDatasetOptions, "listLogpushJobsForDatasetOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"dataset": *logpushJobsApi.Dataset,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/datasets/{dataset}/jobs`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listLogpushJobsForDatasetOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "ListLogpushJobsForDataset")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushOwnership : Get a new ownership challenge sent to your destination
// Get a new ownership challenge.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushOwnership(getLogpushOwnershipOptions *GetLogpushOwnershipOptions) (result *OwnershipChallengeResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.GetLogpushOwnershipWithContext(context.Background(), getLogpushOwnershipOptions)
}

// GetLogpushOwnershipWithContext is an alternate form of the GetLogpushOwnership method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushOwnershipWithContext(ctx context.Context, getLogpushOwnershipOptions *GetLogpushOwnershipOptions) (result *OwnershipChallengeResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogpushOwnershipOptions, "getLogpushOwnershipOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/ownership`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLogpushOwnershipOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "GetLogpushOwnership")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if getLogpushOwnershipOptions.DestinationConf != nil {
		body["destination_conf"] = getLogpushOwnershipOptions.DestinationConf
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOwnershipChallengeResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ValidateLogpushOwnershipChallenge : Validate ownership challenge of the destination
// Validate ownership challenge of the destination.
func (logpushJobsApi *LogpushJobsApiV1) ValidateLogpushOwnershipChallenge(validateLogpushOwnershipChallengeOptions *ValidateLogpushOwnershipChallengeOptions) (result *OwnershipChallengeValidateResult, response *core.DetailedResponse, err error) {
	return logpushJobsApi.ValidateLogpushOwnershipChallengeWithContext(context.Background(), validateLogpushOwnershipChallengeOptions)
}

// ValidateLogpushOwnershipChallengeWithContext is an alternate form of the ValidateLogpushOwnershipChallenge method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ValidateLogpushOwnershipChallengeWithContext(ctx context.Context, validateLogpushOwnershipChallengeOptions *ValidateLogpushOwnershipChallengeOptions) (result *OwnershipChallengeValidateResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(validateLogpushOwnershipChallengeOptions, "validateLogpushOwnershipChallengeOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v1/{crn}/zones/{zone_id}/logpush/ownership/validate`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range validateLogpushOwnershipChallengeOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("logpush_jobs_api", "V1", "ValidateLogpushOwnershipChallenge")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if validateLogpushOwnershipChallengeOptions.DestinationConf != nil {
		body["destination_conf"] = validateLogpushOwnershipChallengeOptions.DestinationConf
	}
	if validateLogpushOwnershipChallengeOptions.OwnershipChallenge != nil {
		body["ownership_challenge"] = validateLogpushOwnershipChallengeOptions.OwnershipChallenge
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
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOwnershipChallengeValidateResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushJobsV2 : List logpush jobs
// List configured logpush jobs for your domain.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobsV2(getLogpushJobsV2Options *GetLogpushJobsV2Options) (result *ListLogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.GetLogpushJobsV2WithContext(context.Background(), getLogpushJobsV2Options)
}

// GetLogpushJobsV2WithContext is an alternate form of the GetLogpushJobsV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobsV2WithContext(ctx context.Context, getLogpushJobsV2Options *GetLogpushJobsV2Options) (result *ListLogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogpushJobsV2Options, "getLogpushJobsV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateLogpushJobV2 : Create a logpush jobs
// Create a new logpush job for the domain.
func (logpushJobsApi *LogpushJobsApiV1) CreateLogpushJobV2(createLogpushJobV2Options *CreateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.CreateLogpushJobV2WithContext(context.Background(), createLogpushJobV2Options)
}

// CreateLogpushJobV2WithContext is an alternate form of the CreateLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) CreateLogpushJobV2WithContext(ctx context.Context, createLogpushJobV2Options *CreateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLogpushJobV2Options, "createLogpushJobV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs`, pathParamsMap)
	if err != nil {
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
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushJobV2 : Get a logpush job
// Get a logpush job  for a given zone.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobV2(getLogpushJobV2Options *GetLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.GetLogpushJobV2WithContext(context.Background(), getLogpushJobV2Options)
}

// GetLogpushJobV2WithContext is an alternate form of the GetLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushJobV2WithContext(ctx context.Context, getLogpushJobV2Options *GetLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLogpushJobV2Options, "getLogpushJobV2Options cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLogpushJobV2Options, "getLogpushJobV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id": fmt.Sprint(*getLogpushJobV2Options.JobID),
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateLogpushJobV2 : Update a logpush job
// Update an existing logpush job for a given zone.
func (logpushJobsApi *LogpushJobsApiV1) UpdateLogpushJobV2(updateLogpushJobV2Options *UpdateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.UpdateLogpushJobV2WithContext(context.Background(), updateLogpushJobV2Options)
}

// UpdateLogpushJobV2WithContext is an alternate form of the UpdateLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) UpdateLogpushJobV2WithContext(ctx context.Context, updateLogpushJobV2Options *UpdateLogpushJobV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLogpushJobV2Options, "updateLogpushJobV2Options cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateLogpushJobV2Options, "updateLogpushJobV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id": fmt.Sprint(*updateLogpushJobV2Options.JobID),
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
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
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteLogpushJobV2 : Delete a logpush job
// Delete a logpush job for a zone.
func (logpushJobsApi *LogpushJobsApiV1) DeleteLogpushJobV2(deleteLogpushJobV2Options *DeleteLogpushJobV2Options) (result *DeleteLogpushJobResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.DeleteLogpushJobV2WithContext(context.Background(), deleteLogpushJobV2Options)
}

// DeleteLogpushJobV2WithContext is an alternate form of the DeleteLogpushJobV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) DeleteLogpushJobV2WithContext(ctx context.Context, deleteLogpushJobV2Options *DeleteLogpushJobV2Options) (result *DeleteLogpushJobResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLogpushJobV2Options, "deleteLogpushJobV2Options cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLogpushJobV2Options, "deleteLogpushJobV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"job_id": fmt.Sprint(*deleteLogpushJobV2Options.JobID),
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/jobs/{job_id}`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteLogpushJobResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetLogpushOwnershipV2 : Get a new ownership challenge sent to your destination
// Get a new ownership challenge.
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushOwnershipV2(getLogpushOwnershipV2Options *GetLogpushOwnershipV2Options) (result *OwnershipChallengeResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.GetLogpushOwnershipV2WithContext(context.Background(), getLogpushOwnershipV2Options)
}

// GetLogpushOwnershipV2WithContext is an alternate form of the GetLogpushOwnershipV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) GetLogpushOwnershipV2WithContext(ctx context.Context, getLogpushOwnershipV2Options *GetLogpushOwnershipV2Options) (result *OwnershipChallengeResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(getLogpushOwnershipV2Options, "getLogpushOwnershipV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/ownership`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOwnershipChallengeResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ValidateLogpushOwnershipChallengeV2 : Validate ownership challenge of the destination
// Validate ownership challenge of the destination.
func (logpushJobsApi *LogpushJobsApiV1) ValidateLogpushOwnershipChallengeV2(validateLogpushOwnershipChallengeV2Options *ValidateLogpushOwnershipChallengeV2Options) (result *OwnershipChallengeValidateResult, response *core.DetailedResponse, err error) {
	return logpushJobsApi.ValidateLogpushOwnershipChallengeV2WithContext(context.Background(), validateLogpushOwnershipChallengeV2Options)
}

// ValidateLogpushOwnershipChallengeV2WithContext is an alternate form of the ValidateLogpushOwnershipChallengeV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ValidateLogpushOwnershipChallengeV2WithContext(ctx context.Context, validateLogpushOwnershipChallengeV2Options *ValidateLogpushOwnershipChallengeV2Options) (result *OwnershipChallengeValidateResult, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(validateLogpushOwnershipChallengeV2Options, "validateLogpushOwnershipChallengeV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/ownership/validate`, pathParamsMap)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOwnershipChallengeValidateResult)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListFieldsForDatasetV2 : The list of all fields available for a dataset
// The list of all fields available for a dataset.
func (logpushJobsApi *LogpushJobsApiV1) ListFieldsForDatasetV2(listFieldsForDatasetV2Options *ListFieldsForDatasetV2Options) (result *ListFieldsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.ListFieldsForDatasetV2WithContext(context.Background(), listFieldsForDatasetV2Options)
}

// ListFieldsForDatasetV2WithContext is an alternate form of the ListFieldsForDatasetV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ListFieldsForDatasetV2WithContext(ctx context.Context, listFieldsForDatasetV2Options *ListFieldsForDatasetV2Options) (result *ListFieldsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listFieldsForDatasetV2Options, "listFieldsForDatasetV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"dataset": *logpushJobsApi.Dataset,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/datasets/{dataset}/fields`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListFieldsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListLogpushJobsForDatasetV2 : List logpush jobs for dataset
// List configured logpush jobs for a dataset.
func (logpushJobsApi *LogpushJobsApiV1) ListLogpushJobsForDatasetV2(listLogpushJobsForDatasetV2Options *ListLogpushJobsForDatasetV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	return logpushJobsApi.ListLogpushJobsForDatasetV2WithContext(context.Background(), listLogpushJobsForDatasetV2Options)
}

// ListLogpushJobsForDatasetV2WithContext is an alternate form of the ListLogpushJobsForDatasetV2 method which supports a Context parameter
func (logpushJobsApi *LogpushJobsApiV1) ListLogpushJobsForDatasetV2WithContext(ctx context.Context, listLogpushJobsForDatasetV2Options *ListLogpushJobsForDatasetV2Options) (result *LogpushJobsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listLogpushJobsForDatasetV2Options, "listLogpushJobsForDatasetV2Options")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *logpushJobsApi.Crn,
		"zone_id": *logpushJobsApi.ZoneID,
		"dataset": *logpushJobsApi.Dataset,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = logpushJobsApi.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(logpushJobsApi.Service.Options.URL, `/v2/{crn}/zones/{zone_id}/logpush/datasets/{dataset}/jobs`, pathParamsMap)
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = logpushJobsApi.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogpushJobsResp)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateLogpushJobOptions : The CreateLogpushJob options.
type CreateLogpushJobOptions struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.
	DestinationConf *string `json:"destination_conf,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateLogpushJobOptions.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobOptions_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobOptions_Dataset_HttpRequests = "http_requests"
	CreateLogpushJobOptions_Dataset_RangeEvents = "range_events"
)

// Constants associated with the CreateLogpushJobOptions.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobOptions_Frequency_High = "high"
	CreateLogpushJobOptions_Frequency_Low = "low"
)

// NewCreateLogpushJobOptions : Instantiate CreateLogpushJobOptions
func (*LogpushJobsApiV1) NewCreateLogpushJobOptions() *CreateLogpushJobOptions {
	return &CreateLogpushJobOptions{}
}

// SetName : Allow user to set Name
func (_options *CreateLogpushJobOptions) SetName(name string) *CreateLogpushJobOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateLogpushJobOptions) SetEnabled(enabled bool) *CreateLogpushJobOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetLogpullOptions : Allow user to set LogpullOptions
func (_options *CreateLogpushJobOptions) SetLogpullOptions(logpullOptions string) *CreateLogpushJobOptions {
	_options.LogpullOptions = core.StringPtr(logpullOptions)
	return _options
}

// SetDestinationConf : Allow user to set DestinationConf
func (_options *CreateLogpushJobOptions) SetDestinationConf(destinationConf string) *CreateLogpushJobOptions {
	_options.DestinationConf = core.StringPtr(destinationConf)
	return _options
}

// SetOwnershipChallenge : Allow user to set OwnershipChallenge
func (_options *CreateLogpushJobOptions) SetOwnershipChallenge(ownershipChallenge string) *CreateLogpushJobOptions {
	_options.OwnershipChallenge = core.StringPtr(ownershipChallenge)
	return _options
}

// SetDataset : Allow user to set Dataset
func (_options *CreateLogpushJobOptions) SetDataset(dataset string) *CreateLogpushJobOptions {
	_options.Dataset = core.StringPtr(dataset)
	return _options
}

// SetFrequency : Allow user to set Frequency
func (_options *CreateLogpushJobOptions) SetFrequency(frequency string) *CreateLogpushJobOptions {
	_options.Frequency = core.StringPtr(frequency)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLogpushJobOptions) SetHeaders(param map[string]string) *CreateLogpushJobOptions {
	options.Headers = param
	return options
}

// CreateLogpushJobV2Options : The CreateLogpushJobV2 options.
type CreateLogpushJobV2Options struct {
	// Create logpush job body.
	CreateLogpushJobV2Request CreateLogpushJobV2RequestIntf `json:"CreateLogpushJobV2Request,omitempty"`

	// Allows users to set headers on API requests
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
type CreateLogpushJobV2Request struct {
	// Logpush Job Name.
	Name *string `json:"name,omitempty"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the COS bucket where the data will be pushed.
	Cos interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`

	// Information to identify the LogDNA instance the data will be pushed.
	Logdna interface{} `json:"logdna,omitempty"`
}

// Constants associated with the CreateLogpushJobV2Request.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2Request_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2Request_Dataset_HttpRequests = "http_requests"
	CreateLogpushJobV2Request_Dataset_RangeEvents = "range_events"
)

// Constants associated with the CreateLogpushJobV2Request.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2Request_Frequency_High = "high"
	CreateLogpushJobV2Request_Frequency_Low = "low"
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
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteLogpushJobOptions : The DeleteLogpushJob options.
type DeleteLogpushJobOptions struct {
	// logpush job identifier.
	JobID *int64 `json:"job_id" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLogpushJobOptions : Instantiate DeleteLogpushJobOptions
func (*LogpushJobsApiV1) NewDeleteLogpushJobOptions(jobID int64) *DeleteLogpushJobOptions {
	return &DeleteLogpushJobOptions{
		JobID: core.Int64Ptr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *DeleteLogpushJobOptions) SetJobID(jobID int64) *DeleteLogpushJobOptions {
	_options.JobID = core.Int64Ptr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLogpushJobOptions) SetHeaders(param map[string]string) *DeleteLogpushJobOptions {
	options.Headers = param
	return options
}

// DeleteLogpushJobV2Options : The DeleteLogpushJobV2 options.
type DeleteLogpushJobV2Options struct {
	// logpush job identifier.
	JobID *int64 `json:"job_id" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLogpushJobV2Options : Instantiate DeleteLogpushJobV2Options
func (*LogpushJobsApiV1) NewDeleteLogpushJobV2Options(jobID int64) *DeleteLogpushJobV2Options {
	return &DeleteLogpushJobV2Options{
		JobID: core.Int64Ptr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *DeleteLogpushJobV2Options) SetJobID(jobID int64) *DeleteLogpushJobV2Options {
	_options.JobID = core.Int64Ptr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLogpushJobV2Options) SetHeaders(param map[string]string) *DeleteLogpushJobV2Options {
	options.Headers = param
	return options
}

// GetLogpushJobOptions : The GetLogpushJob options.
type GetLogpushJobOptions struct {
	// logpush job identifier.
	JobID *int64 `json:"job_id" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLogpushJobOptions : Instantiate GetLogpushJobOptions
func (*LogpushJobsApiV1) NewGetLogpushJobOptions(jobID int64) *GetLogpushJobOptions {
	return &GetLogpushJobOptions{
		JobID: core.Int64Ptr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *GetLogpushJobOptions) SetJobID(jobID int64) *GetLogpushJobOptions {
	_options.JobID = core.Int64Ptr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushJobOptions) SetHeaders(param map[string]string) *GetLogpushJobOptions {
	options.Headers = param
	return options
}

// GetLogpushJobV2Options : The GetLogpushJobV2 options.
type GetLogpushJobV2Options struct {
	// logpush job identifier.
	JobID *int64 `json:"job_id" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLogpushJobV2Options : Instantiate GetLogpushJobV2Options
func (*LogpushJobsApiV1) NewGetLogpushJobV2Options(jobID int64) *GetLogpushJobV2Options {
	return &GetLogpushJobV2Options{
		JobID: core.Int64Ptr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *GetLogpushJobV2Options) SetJobID(jobID int64) *GetLogpushJobV2Options {
	_options.JobID = core.Int64Ptr(jobID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushJobV2Options) SetHeaders(param map[string]string) *GetLogpushJobV2Options {
	options.Headers = param
	return options
}

// GetLogpushJobsOptions : The GetLogpushJobs options.
type GetLogpushJobsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLogpushJobsOptions : Instantiate GetLogpushJobsOptions
func (*LogpushJobsApiV1) NewGetLogpushJobsOptions() *GetLogpushJobsOptions {
	return &GetLogpushJobsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushJobsOptions) SetHeaders(param map[string]string) *GetLogpushJobsOptions {
	options.Headers = param
	return options
}

// GetLogpushJobsV2Options : The GetLogpushJobsV2 options.
type GetLogpushJobsV2Options struct {

	// Allows users to set headers on API requests
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

// GetLogpushOwnershipOptions : The GetLogpushOwnership options.
type GetLogpushOwnershipOptions struct {
	// Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.
	DestinationConf *string `json:"destination_conf,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLogpushOwnershipOptions : Instantiate GetLogpushOwnershipOptions
func (*LogpushJobsApiV1) NewGetLogpushOwnershipOptions() *GetLogpushOwnershipOptions {
	return &GetLogpushOwnershipOptions{}
}

// SetDestinationConf : Allow user to set DestinationConf
func (_options *GetLogpushOwnershipOptions) SetDestinationConf(destinationConf string) *GetLogpushOwnershipOptions {
	_options.DestinationConf = core.StringPtr(destinationConf)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushOwnershipOptions) SetHeaders(param map[string]string) *GetLogpushOwnershipOptions {
	options.Headers = param
	return options
}

// GetLogpushOwnershipV2Options : The GetLogpushOwnershipV2 options.
type GetLogpushOwnershipV2Options struct {
	// Information to identify the COS bucket where the data will be pushed.
	Cos interface{} `json:"cos,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLogpushOwnershipV2Options : Instantiate GetLogpushOwnershipV2Options
func (*LogpushJobsApiV1) NewGetLogpushOwnershipV2Options() *GetLogpushOwnershipV2Options {
	return &GetLogpushOwnershipV2Options{}
}

// SetCos : Allow user to set Cos
func (_options *GetLogpushOwnershipV2Options) SetCos(cos interface{}) *GetLogpushOwnershipV2Options {
	_options.Cos = cos
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetLogpushOwnershipV2Options) SetHeaders(param map[string]string) *GetLogpushOwnershipV2Options {
	options.Headers = param
	return options
}

// ListFieldsForDatasetOptions : The ListFieldsForDataset options.
type ListFieldsForDatasetOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListFieldsForDatasetOptions : Instantiate ListFieldsForDatasetOptions
func (*LogpushJobsApiV1) NewListFieldsForDatasetOptions() *ListFieldsForDatasetOptions {
	return &ListFieldsForDatasetOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListFieldsForDatasetOptions) SetHeaders(param map[string]string) *ListFieldsForDatasetOptions {
	options.Headers = param
	return options
}

// ListFieldsForDatasetV2Options : The ListFieldsForDatasetV2 options.
type ListFieldsForDatasetV2Options struct {

	// Allows users to set headers on API requests
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

// ListLogpushJobsForDatasetOptions : The ListLogpushJobsForDataset options.
type ListLogpushJobsForDatasetOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListLogpushJobsForDatasetOptions : Instantiate ListLogpushJobsForDatasetOptions
func (*LogpushJobsApiV1) NewListLogpushJobsForDatasetOptions() *ListLogpushJobsForDatasetOptions {
	return &ListLogpushJobsForDatasetOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListLogpushJobsForDatasetOptions) SetHeaders(param map[string]string) *ListLogpushJobsForDatasetOptions {
	options.Headers = param
	return options
}

// ListLogpushJobsForDatasetV2Options : The ListLogpushJobsForDatasetV2 options.
type ListLogpushJobsForDatasetV2Options struct {

	// Allows users to set headers on API requests
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

// UpdateLogpushJobOptions : The UpdateLogpushJob options.
type UpdateLogpushJobOptions struct {
	// logpush job identifier.
	JobID *int64 `json:"job_id" validate:"required"`

	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.
	DestinationConf *string `json:"destination_conf,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateLogpushJobOptions.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobOptions_Frequency_High = "high"
	UpdateLogpushJobOptions_Frequency_Low = "low"
)

// NewUpdateLogpushJobOptions : Instantiate UpdateLogpushJobOptions
func (*LogpushJobsApiV1) NewUpdateLogpushJobOptions(jobID int64) *UpdateLogpushJobOptions {
	return &UpdateLogpushJobOptions{
		JobID: core.Int64Ptr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *UpdateLogpushJobOptions) SetJobID(jobID int64) *UpdateLogpushJobOptions {
	_options.JobID = core.Int64Ptr(jobID)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateLogpushJobOptions) SetEnabled(enabled bool) *UpdateLogpushJobOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetLogpullOptions : Allow user to set LogpullOptions
func (_options *UpdateLogpushJobOptions) SetLogpullOptions(logpullOptions string) *UpdateLogpushJobOptions {
	_options.LogpullOptions = core.StringPtr(logpullOptions)
	return _options
}

// SetDestinationConf : Allow user to set DestinationConf
func (_options *UpdateLogpushJobOptions) SetDestinationConf(destinationConf string) *UpdateLogpushJobOptions {
	_options.DestinationConf = core.StringPtr(destinationConf)
	return _options
}

// SetOwnershipChallenge : Allow user to set OwnershipChallenge
func (_options *UpdateLogpushJobOptions) SetOwnershipChallenge(ownershipChallenge string) *UpdateLogpushJobOptions {
	_options.OwnershipChallenge = core.StringPtr(ownershipChallenge)
	return _options
}

// SetFrequency : Allow user to set Frequency
func (_options *UpdateLogpushJobOptions) SetFrequency(frequency string) *UpdateLogpushJobOptions {
	_options.Frequency = core.StringPtr(frequency)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLogpushJobOptions) SetHeaders(param map[string]string) *UpdateLogpushJobOptions {
	options.Headers = param
	return options
}

// UpdateLogpushJobV2Options : The UpdateLogpushJobV2 options.
type UpdateLogpushJobV2Options struct {
	// logpush job identifier.
	JobID *int64 `json:"job_id" validate:"required"`

	// Update logpush job.
	UpdateLogpushJobV2Request UpdateLogpushJobV2RequestIntf `json:"UpdateLogpushJobV2Request,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateLogpushJobV2Options : Instantiate UpdateLogpushJobV2Options
func (*LogpushJobsApiV1) NewUpdateLogpushJobV2Options(jobID int64) *UpdateLogpushJobV2Options {
	return &UpdateLogpushJobV2Options{
		JobID: core.Int64Ptr(jobID),
	}
}

// SetJobID : Allow user to set JobID
func (_options *UpdateLogpushJobV2Options) SetJobID(jobID int64) *UpdateLogpushJobV2Options {
	_options.JobID = core.Int64Ptr(jobID)
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
type UpdateLogpushJobV2Request struct {
	// Whether the logpush job enabled or not.
	Enabled *bool `json:"enabled,omitempty"`

	// Configuration string.
	LogpullOptions *string `json:"logpull_options,omitempty"`

	// Information to identify the COS bucket where the data will be pushed.
	Cos interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`

	// Information to identify the LogDNA instance the data will be pushed.
	Logdna interface{} `json:"logdna,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2Request.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2Request_Frequency_High = "high"
	UpdateLogpushJobV2Request_Frequency_Low = "low"
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
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ValidateLogpushOwnershipChallengeOptions : The ValidateLogpushOwnershipChallenge options.
type ValidateLogpushOwnershipChallengeOptions struct {
	// Uniquely identifies a resource (such as an s3 bucket) where data will be pushed.
	DestinationConf *string `json:"destination_conf,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewValidateLogpushOwnershipChallengeOptions : Instantiate ValidateLogpushOwnershipChallengeOptions
func (*LogpushJobsApiV1) NewValidateLogpushOwnershipChallengeOptions() *ValidateLogpushOwnershipChallengeOptions {
	return &ValidateLogpushOwnershipChallengeOptions{}
}

// SetDestinationConf : Allow user to set DestinationConf
func (_options *ValidateLogpushOwnershipChallengeOptions) SetDestinationConf(destinationConf string) *ValidateLogpushOwnershipChallengeOptions {
	_options.DestinationConf = core.StringPtr(destinationConf)
	return _options
}

// SetOwnershipChallenge : Allow user to set OwnershipChallenge
func (_options *ValidateLogpushOwnershipChallengeOptions) SetOwnershipChallenge(ownershipChallenge string) *ValidateLogpushOwnershipChallengeOptions {
	_options.OwnershipChallenge = core.StringPtr(ownershipChallenge)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ValidateLogpushOwnershipChallengeOptions) SetHeaders(param map[string]string) *ValidateLogpushOwnershipChallengeOptions {
	options.Headers = param
	return options
}

// ValidateLogpushOwnershipChallengeV2Options : The ValidateLogpushOwnershipChallengeV2 options.
type ValidateLogpushOwnershipChallengeV2Options struct {
	// Information to identify the COS bucket where the data will be pushed.
	Cos interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewValidateLogpushOwnershipChallengeV2Options : Instantiate ValidateLogpushOwnershipChallengeV2Options
func (*LogpushJobsApiV1) NewValidateLogpushOwnershipChallengeV2Options() *ValidateLogpushOwnershipChallengeV2Options {
	return &ValidateLogpushOwnershipChallengeV2Options{}
}

// SetCos : Allow user to set Cos
func (_options *ValidateLogpushOwnershipChallengeV2Options) SetCos(cos interface{}) *ValidateLogpushOwnershipChallengeV2Options {
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
	// success respose.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result interface{} `json:"result" validate:"required"`
}

// UnmarshalDeleteLogpushJobResp unmarshals an instance of DeleteLogpushJobResp from the specified map of raw messages.
func UnmarshalDeleteLogpushJobResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLogpushJobResp)
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
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
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
	Result interface{} `json:"result,omitempty"`
}

// UnmarshalListFieldsResp unmarshals an instance of ListFieldsResp from the specified map of raw messages.
func UnmarshalListFieldsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListFieldsResp)
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
	err = core.UnmarshalPrimitive(m, "result", &obj.Result)
	if err != nil {
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLogpushJobPack)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "destination_conf", &obj.DestinationConf)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_complete", &obj.LastComplete)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_error", &obj.LastError)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "error_message", &obj.ErrorMessage)
	if err != nil {
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLogpushJobPack)
	if err != nil {
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalOwnershipChallengeResult)
	if err != nil {
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
		return
	}
	err = core.UnmarshalPrimitive(m, "valid", &obj.Valid)
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
	Cos interface{} `json:"cos" validate:"required"`

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
	CreateLogpushJobV2RequestLogpushJobCosReq_Dataset_HttpRequests = "http_requests"
	CreateLogpushJobV2RequestLogpushJobCosReq_Dataset_RangeEvents = "range_events"
)

// Constants associated with the CreateLogpushJobV2RequestLogpushJobCosReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2RequestLogpushJobCosReq_Frequency_High = "high"
	CreateLogpushJobV2RequestLogpushJobCosReq_Frequency_Low = "low"
)

// NewCreateLogpushJobV2RequestLogpushJobCosReq : Instantiate CreateLogpushJobV2RequestLogpushJobCosReq (Generic Model Constructor)
func (*LogpushJobsApiV1) NewCreateLogpushJobV2RequestLogpushJobCosReq(cos interface{}, ownershipChallenge string) (_model *CreateLogpushJobV2RequestLogpushJobCosReq, err error) {
	_model = &CreateLogpushJobV2RequestLogpushJobCosReq{
		Cos: cos,
		OwnershipChallenge: core.StringPtr(ownershipChallenge),
	}
	err = core.ValidateStruct(_model, "required parameters")
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
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
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
	Logdna interface{} `json:"logdna" validate:"required"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// Dataset to be pulled.
	Dataset *string `json:"dataset,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the CreateLogpushJobV2RequestLogpushJobLogdnaReq.Dataset property.
// Dataset to be pulled.
const (
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Dataset_FirewallEvents = "firewall_events"
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Dataset_HttpRequests = "http_requests"
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Dataset_RangeEvents = "range_events"
)

// Constants associated with the CreateLogpushJobV2RequestLogpushJobLogdnaReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Frequency_High = "high"
	CreateLogpushJobV2RequestLogpushJobLogdnaReq_Frequency_Low = "low"
)

// NewCreateLogpushJobV2RequestLogpushJobLogdnaReq : Instantiate CreateLogpushJobV2RequestLogpushJobLogdnaReq (Generic Model Constructor)
func (*LogpushJobsApiV1) NewCreateLogpushJobV2RequestLogpushJobLogdnaReq(logdna interface{}) (_model *CreateLogpushJobV2RequestLogpushJobLogdnaReq, err error) {
	_model = &CreateLogpushJobV2RequestLogpushJobLogdnaReq{
		Logdna: logdna,
	}
	err = core.ValidateStruct(_model, "required parameters")
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
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "dataset", &obj.Dataset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
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
	Cos interface{} `json:"cos,omitempty"`

	// Ownership challenge token to prove destination ownership.
	OwnershipChallenge *string `json:"ownership_challenge,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq_Frequency_High = "high"
	UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq_Frequency_Low = "low"
)

func (*UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq) isaUpdateLogpushJobV2Request() bool {
	return true
}

// UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateCosReq unmarshals an instance of UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateCosReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2RequestLogpushJobsUpdateCosReq)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cos", &obj.Cos)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ownership_challenge", &obj.OwnershipChallenge)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
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
	Logdna interface{} `json:"logdna,omitempty"`

	// The frequency at which CIS sends batches of logs to your destination.
	Frequency *string `json:"frequency,omitempty"`
}

// Constants associated with the UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq.Frequency property.
// The frequency at which CIS sends batches of logs to your destination.
const (
	UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq_Frequency_High = "high"
	UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq_Frequency_Low = "low"
)

func (*UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq) isaUpdateLogpushJobV2Request() bool {
	return true
}

// UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq unmarshals an instance of UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq from the specified map of raw messages.
func UnmarshalUpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateLogpushJobV2RequestLogpushJobsUpdateLogdnaReq)
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logpull_options", &obj.LogpullOptions)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "logdna", &obj.Logdna)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "frequency", &obj.Frequency)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
