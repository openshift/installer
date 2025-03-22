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
 * IBM OpenAPI SDK Code Generator Version: 3.95.2-120e65bc-20240924-152329
 */

// Package cdtektonpipelinev2 : Operations and models for the CdTektonPipelineV2 service
package cdtektonpipelinev2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	common "github.com/IBM/continuous-delivery-go-sdk/v2/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

// CdTektonPipelineV2 : Continuous Delivery Tekton pipeline API definition <br><br> Maximum request payload size is 512
// KB <br><br> All calls require an <strong>Authorization</strong> HTTP header. <br><br> The following header is the
// accepted authentication mechanism and required credentials for each <ul><li><b>Bearer:</b> an IBM Cloud IAM token
// (authorized for all endpoints)</li>
//
// API Version: 2.0.0
type CdTektonPipelineV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.us-south.devops.cloud.ibm.com/pipeline/v2"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "cd_tekton_pipeline"

// CdTektonPipelineV2Options : Service options
type CdTektonPipelineV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewCdTektonPipelineV2UsingExternalConfig : constructs an instance of CdTektonPipelineV2 with passed in options and external configuration.
func NewCdTektonPipelineV2UsingExternalConfig(options *CdTektonPipelineV2Options) (cdTektonPipeline *CdTektonPipelineV2, err error) {
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

	cdTektonPipeline, err = NewCdTektonPipelineV2(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = cdTektonPipeline.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = cdTektonPipeline.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewCdTektonPipelineV2 : constructs an instance of CdTektonPipelineV2 with passed in options.
func NewCdTektonPipelineV2(options *CdTektonPipelineV2Options) (service *CdTektonPipelineV2, err error) {
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

	service = &CdTektonPipelineV2{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://api.us-south.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the us-south region.
		"us-east": "https://api.us-east.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the us-east region.
		"eu-de": "https://api.eu-de.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the eu-de region.
		"eu-gb": "https://api.eu-gb.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the eu-gb region.
		"eu-es": "https://api.eu-es.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the eu-es region.
		"jp-osa": "https://api.jp-osa.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the jp-osa region.
		"jp-tok": "https://api.jp-tok.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the jp-tok region.
		"au-syd": "https://api.au-syd.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the au-syd region.
		"ca-tor": "https://api.ca-tor.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the ca-tor region.
		"br-sao": "https://api.br-sao.devops.cloud.ibm.com/pipeline/v2", // The host URL for Tekton Pipeline Service in the br-sao region.
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", core.SDKErrorf(nil, fmt.Sprintf("service URL for region '%s' not found", region), "invalid-region", common.GetComponentInfo())
}

// Clone makes a copy of "cdTektonPipeline" suitable for processing requests.
func (cdTektonPipeline *CdTektonPipelineV2) Clone() *CdTektonPipelineV2 {
	if core.IsNil(cdTektonPipeline) {
		return nil
	}
	clone := *cdTektonPipeline
	clone.Service = cdTektonPipeline.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (cdTektonPipeline *CdTektonPipelineV2) SetServiceURL(url string) error {
	err := cdTektonPipeline.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (cdTektonPipeline *CdTektonPipelineV2) GetServiceURL() string {
	return cdTektonPipeline.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (cdTektonPipeline *CdTektonPipelineV2) SetDefaultHeaders(headers http.Header) {
	cdTektonPipeline.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (cdTektonPipeline *CdTektonPipelineV2) SetEnableGzipCompression(enableGzip bool) {
	cdTektonPipeline.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (cdTektonPipeline *CdTektonPipelineV2) GetEnableGzipCompression() bool {
	return cdTektonPipeline.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (cdTektonPipeline *CdTektonPipelineV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	cdTektonPipeline.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (cdTektonPipeline *CdTektonPipelineV2) DisableRetries() {
	cdTektonPipeline.Service.DisableRetries()
}

// CreateTektonPipeline : Create Tekton pipeline
// This request creates a Tekton pipeline. Requires a pipeline tool already created in the toolchain using the toolchain
// API https://cloud.ibm.com/apidocs/toolchain#create-tool, and use the tool ID to create the Tekton pipeline.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipeline(createTektonPipelineOptions *CreateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CreateTektonPipelineWithContext(context.Background(), createTektonPipelineOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTektonPipelineWithContext is an alternate form of the CreateTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineWithContext(ctx context.Context, createTektonPipelineOptions *CreateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineOptions, "createTektonPipelineOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTektonPipelineOptions, "createTektonPipelineOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTektonPipelineOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CreateTektonPipeline")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTektonPipelineOptions.ID != nil {
		body["id"] = createTektonPipelineOptions.ID
	}
	if createTektonPipelineOptions.NextBuildNumber != nil {
		body["next_build_number"] = createTektonPipelineOptions.NextBuildNumber
	}
	if createTektonPipelineOptions.EnableNotifications != nil {
		body["enable_notifications"] = createTektonPipelineOptions.EnableNotifications
	}
	if createTektonPipelineOptions.EnablePartialCloning != nil {
		body["enable_partial_cloning"] = createTektonPipelineOptions.EnablePartialCloning
	}
	if createTektonPipelineOptions.Worker != nil {
		body["worker"] = createTektonPipelineOptions.Worker
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tekton_pipeline", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTektonPipeline)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipeline : Get Tekton pipeline data
// This request retrieves the Tekton pipeline data for the pipeline identified by `{id}`.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipeline(getTektonPipelineOptions *GetTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineWithContext(context.Background(), getTektonPipelineOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineWithContext is an alternate form of the GetTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineWithContext(ctx context.Context, getTektonPipelineOptions *GetTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineOptions, "getTektonPipelineOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineOptions, "getTektonPipelineOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getTektonPipelineOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipeline")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTektonPipeline)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateTektonPipeline : Update Tekton pipeline data
// This request updates Tekton pipeline data, but you can only change worker ID in this endpoint. Use other endpoints
// such as /definitions, /triggers, and /properties for other configuration updates.
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipeline(updateTektonPipelineOptions *UpdateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.UpdateTektonPipelineWithContext(context.Background(), updateTektonPipelineOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateTektonPipelineWithContext is an alternate form of the UpdateTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipelineWithContext(ctx context.Context, updateTektonPipelineOptions *UpdateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTektonPipelineOptions, "updateTektonPipelineOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateTektonPipelineOptions, "updateTektonPipelineOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *updateTektonPipelineOptions.ID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateTektonPipelineOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "UpdateTektonPipeline")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateTektonPipelineOptions.TektonPipelinePatch != nil {
		_, err = builder.SetBodyContentJSON(updateTektonPipelineOptions.TektonPipelinePatch)
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_tekton_pipeline", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTektonPipeline)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipeline : Delete Tekton pipeline instance
// This request deletes Tekton pipeline instance that is associated with the pipeline toolchain integration.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipeline(deleteTektonPipelineOptions *DeleteTektonPipelineOptions) (response *core.DetailedResponse, err error) {
	response, err = cdTektonPipeline.DeleteTektonPipelineWithContext(context.Background(), deleteTektonPipelineOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTektonPipelineWithContext is an alternate form of the DeleteTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineWithContext(ctx context.Context, deleteTektonPipelineOptions *DeleteTektonPipelineOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineOptions, "deleteTektonPipelineOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineOptions, "deleteTektonPipelineOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteTektonPipelineOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTektonPipelineOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DeleteTektonPipeline")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tekton_pipeline", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListTektonPipelineRuns : List pipeline run records
// This request lists pipeline run records, which has data about the runs, such as status, user_info, trigger and other
// information. Default limit is 50.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineRuns(listTektonPipelineRunsOptions *ListTektonPipelineRunsOptions) (result *PipelineRunsCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ListTektonPipelineRunsWithContext(context.Background(), listTektonPipelineRunsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTektonPipelineRunsWithContext is an alternate form of the ListTektonPipelineRuns method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineRunsWithContext(ctx context.Context, listTektonPipelineRunsOptions *ListTektonPipelineRunsOptions) (result *PipelineRunsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineRunsOptions, "listTektonPipelineRunsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listTektonPipelineRunsOptions, "listTektonPipelineRunsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *listTektonPipelineRunsOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTektonPipelineRunsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ListTektonPipelineRuns")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listTektonPipelineRunsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listTektonPipelineRunsOptions.Start))
	}
	if listTektonPipelineRunsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listTektonPipelineRunsOptions.Limit))
	}
	if listTektonPipelineRunsOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listTektonPipelineRunsOptions.Status))
	}
	if listTektonPipelineRunsOptions.TriggerName != nil {
		builder.AddQuery("trigger.name", fmt.Sprint(*listTektonPipelineRunsOptions.TriggerName))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tekton_pipeline_runs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRunsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineRun : Trigger a pipeline run
// Trigger a new pipeline run using either the manual or the timed trigger, specifying the additional properties or
// overriding existing ones as needed.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineRun(createTektonPipelineRunOptions *CreateTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CreateTektonPipelineRunWithContext(context.Background(), createTektonPipelineRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTektonPipelineRunWithContext is an alternate form of the CreateTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineRunWithContext(ctx context.Context, createTektonPipelineRunOptions *CreateTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineRunOptions, "createTektonPipelineRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTektonPipelineRunOptions, "createTektonPipelineRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *createTektonPipelineRunOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTektonPipelineRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CreateTektonPipelineRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTektonPipelineRunOptions.Description != nil {
		body["description"] = createTektonPipelineRunOptions.Description
	}
	if createTektonPipelineRunOptions.TriggerName != nil {
		body["trigger_name"] = createTektonPipelineRunOptions.TriggerName
	}
	if createTektonPipelineRunOptions.TriggerProperties != nil {
		body["trigger_properties"] = createTektonPipelineRunOptions.TriggerProperties
	}
	if createTektonPipelineRunOptions.SecureTriggerProperties != nil {
		body["secure_trigger_properties"] = createTektonPipelineRunOptions.SecureTriggerProperties
	}
	if createTektonPipelineRunOptions.TriggerHeaders != nil {
		body["trigger_headers"] = createTektonPipelineRunOptions.TriggerHeaders
	}
	if createTektonPipelineRunOptions.TriggerBody != nil {
		body["trigger_body"] = createTektonPipelineRunOptions.TriggerBody
	}
	if createTektonPipelineRunOptions.Trigger != nil {
		body["trigger"] = createTektonPipelineRunOptions.Trigger
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tekton_pipeline_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineRun : Get a pipeline run record
// This request retrieves details of the pipeline run identified by `{id}`.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRun(getTektonPipelineRunOptions *GetTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineRunWithContext(context.Background(), getTektonPipelineRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineRunWithContext is an alternate form of the GetTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunWithContext(ctx context.Context, getTektonPipelineRunOptions *GetTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineRunOptions, "getTektonPipelineRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineRunOptions, "getTektonPipelineRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelineRunOptions.PipelineID,
		"id": *getTektonPipelineRunOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getTektonPipelineRunOptions.Includes != nil {
		builder.AddQuery("includes", fmt.Sprint(*getTektonPipelineRunOptions.Includes))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineRun : Delete a pipeline run record
// This request deletes the pipeline run record identified by `{id}`.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineRun(deleteTektonPipelineRunOptions *DeleteTektonPipelineRunOptions) (response *core.DetailedResponse, err error) {
	response, err = cdTektonPipeline.DeleteTektonPipelineRunWithContext(context.Background(), deleteTektonPipelineRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTektonPipelineRunWithContext is an alternate form of the DeleteTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineRunWithContext(ctx context.Context, deleteTektonPipelineRunOptions *DeleteTektonPipelineRunOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineRunOptions, "deleteTektonPipelineRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineRunOptions, "deleteTektonPipelineRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *deleteTektonPipelineRunOptions.PipelineID,
		"id": *deleteTektonPipelineRunOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTektonPipelineRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DeleteTektonPipelineRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tekton_pipeline_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CancelTektonPipelineRun : Cancel a pipeline run
// This request cancels a running pipeline run identified by `{id}`. Use `force: true` in the body if the pipeline run
// can't be cancelled normally.
func (cdTektonPipeline *CdTektonPipelineV2) CancelTektonPipelineRun(cancelTektonPipelineRunOptions *CancelTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CancelTektonPipelineRunWithContext(context.Background(), cancelTektonPipelineRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CancelTektonPipelineRunWithContext is an alternate form of the CancelTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CancelTektonPipelineRunWithContext(ctx context.Context, cancelTektonPipelineRunOptions *CancelTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(cancelTektonPipelineRunOptions, "cancelTektonPipelineRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(cancelTektonPipelineRunOptions, "cancelTektonPipelineRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *cancelTektonPipelineRunOptions.PipelineID,
		"id": *cancelTektonPipelineRunOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs/{id}/cancel`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range cancelTektonPipelineRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CancelTektonPipelineRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if cancelTektonPipelineRunOptions.Force != nil {
		body["force"] = cancelTektonPipelineRunOptions.Force
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "cancel_tekton_pipeline_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// RerunTektonPipelineRun : Rerun a pipeline run
// This request reruns a past pipeline run, which is identified by `{id}`, with the same data. Request body isn't
// allowed.
func (cdTektonPipeline *CdTektonPipelineV2) RerunTektonPipelineRun(rerunTektonPipelineRunOptions *RerunTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.RerunTektonPipelineRunWithContext(context.Background(), rerunTektonPipelineRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// RerunTektonPipelineRunWithContext is an alternate form of the RerunTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) RerunTektonPipelineRunWithContext(ctx context.Context, rerunTektonPipelineRunOptions *RerunTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(rerunTektonPipelineRunOptions, "rerunTektonPipelineRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(rerunTektonPipelineRunOptions, "rerunTektonPipelineRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *rerunTektonPipelineRunOptions.PipelineID,
		"id": *rerunTektonPipelineRunOptions.ID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs/{id}/rerun`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range rerunTektonPipelineRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "RerunTektonPipelineRun")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "rerun_tekton_pipeline_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineRunLogs : Get a list of pipeline run log objects
// This request fetches a list of log data for a pipeline run identified by `{id}`. The `href` in each log entry can be
// used to fetch that individual log.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogs(getTektonPipelineRunLogsOptions *GetTektonPipelineRunLogsOptions) (result *LogsCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineRunLogsWithContext(context.Background(), getTektonPipelineRunLogsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineRunLogsWithContext is an alternate form of the GetTektonPipelineRunLogs method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogsWithContext(ctx context.Context, getTektonPipelineRunLogsOptions *GetTektonPipelineRunLogsOptions) (result *LogsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineRunLogsOptions, "getTektonPipelineRunLogsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineRunLogsOptions, "getTektonPipelineRunLogsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelineRunLogsOptions.PipelineID,
		"id": *getTektonPipelineRunLogsOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs/{id}/logs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineRunLogsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineRunLogs")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_run_logs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLogsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineRunLogContent : Get the log content of a pipeline run step
// This request retrieves the log content of a pipeline run step, where the step is identified by `{id}`. To get the log
// ID use the endpoint `/tekton_pipelines/{pipeline_id}/pipeline_runs/{id}/logs`.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogContent(getTektonPipelineRunLogContentOptions *GetTektonPipelineRunLogContentOptions) (result *StepLog, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineRunLogContentWithContext(context.Background(), getTektonPipelineRunLogContentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineRunLogContentWithContext is an alternate form of the GetTektonPipelineRunLogContent method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogContentWithContext(ctx context.Context, getTektonPipelineRunLogContentOptions *GetTektonPipelineRunLogContentOptions) (result *StepLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineRunLogContentOptions, "getTektonPipelineRunLogContentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineRunLogContentOptions, "getTektonPipelineRunLogContentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelineRunLogContentOptions.PipelineID,
		"pipeline_run_id": *getTektonPipelineRunLogContentOptions.PipelineRunID,
		"id": *getTektonPipelineRunLogContentOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/pipeline_runs/{pipeline_run_id}/logs/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineRunLogContentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineRunLogContent")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_run_log_content", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStepLog)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListTektonPipelineDefinitions : List pipeline definitions
// This request fetches pipeline definitions, which is a collection of individual definition entries. Each entry
// consists of a repository url, a repository path and a branch or tag. The referenced repository URL must match the URL
// of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
// https://cloud.ibm.com/apidocs/toolchain#list-tools. The branch or tag of the definition must match against a
// corresponding branch or tag in the chosen repository, and the path must match a subfolder in the repository.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineDefinitions(listTektonPipelineDefinitionsOptions *ListTektonPipelineDefinitionsOptions) (result *DefinitionsCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ListTektonPipelineDefinitionsWithContext(context.Background(), listTektonPipelineDefinitionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTektonPipelineDefinitionsWithContext is an alternate form of the ListTektonPipelineDefinitions method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineDefinitionsWithContext(ctx context.Context, listTektonPipelineDefinitionsOptions *ListTektonPipelineDefinitionsOptions) (result *DefinitionsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineDefinitionsOptions, "listTektonPipelineDefinitionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listTektonPipelineDefinitionsOptions, "listTektonPipelineDefinitionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *listTektonPipelineDefinitionsOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/definitions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTektonPipelineDefinitionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ListTektonPipelineDefinitions")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tekton_pipeline_definitions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinitionsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineDefinition : Create a single definition
// This request adds a single definition. The source properties should consist of a repository url, a repository path
// and a branch or tag. The referenced repository URL must match the URL of a repository tool integration in the parent
// toolchain. Obtain the list of integrations from the toolchain API https://cloud.ibm.com/apidocs/toolchain#list-tools.
// The branch or tag of the definition must match against a corresponding branch or tag in the chosen repository, and
// the path must match a subfolder in the repository.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineDefinition(createTektonPipelineDefinitionOptions *CreateTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CreateTektonPipelineDefinitionWithContext(context.Background(), createTektonPipelineDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTektonPipelineDefinitionWithContext is an alternate form of the CreateTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineDefinitionWithContext(ctx context.Context, createTektonPipelineDefinitionOptions *CreateTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineDefinitionOptions, "createTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTektonPipelineDefinitionOptions, "createTektonPipelineDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *createTektonPipelineDefinitionOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/definitions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTektonPipelineDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CreateTektonPipelineDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTektonPipelineDefinitionOptions.Source != nil {
		body["source"] = createTektonPipelineDefinitionOptions.Source
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tekton_pipeline_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinition)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineDefinition : Retrieve a single definition entry
// This request fetches a single definition entry, which consists of the definition repository URL, a repository path,
// and a branch or tag.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineDefinition(getTektonPipelineDefinitionOptions *GetTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineDefinitionWithContext(context.Background(), getTektonPipelineDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineDefinitionWithContext is an alternate form of the GetTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineDefinitionWithContext(ctx context.Context, getTektonPipelineDefinitionOptions *GetTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineDefinitionOptions, "getTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineDefinitionOptions, "getTektonPipelineDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelineDefinitionOptions.PipelineID,
		"definition_id": *getTektonPipelineDefinitionOptions.DefinitionID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/definitions/{definition_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineDefinition")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinition)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTektonPipelineDefinition : Edit a single definition entry
// This request updates a definition entry identified by `{definition_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineDefinition(replaceTektonPipelineDefinitionOptions *ReplaceTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ReplaceTektonPipelineDefinitionWithContext(context.Background(), replaceTektonPipelineDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceTektonPipelineDefinitionWithContext is an alternate form of the ReplaceTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineDefinitionWithContext(ctx context.Context, replaceTektonPipelineDefinitionOptions *ReplaceTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTektonPipelineDefinitionOptions, "replaceTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceTektonPipelineDefinitionOptions, "replaceTektonPipelineDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *replaceTektonPipelineDefinitionOptions.PipelineID,
		"definition_id": *replaceTektonPipelineDefinitionOptions.DefinitionID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/definitions/{definition_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceTektonPipelineDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ReplaceTektonPipelineDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceTektonPipelineDefinitionOptions.Source != nil {
		body["source"] = replaceTektonPipelineDefinitionOptions.Source
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_tekton_pipeline_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinition)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineDefinition : Delete a single definition entry
// This request deletes a single definition from the definition list.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineDefinition(deleteTektonPipelineDefinitionOptions *DeleteTektonPipelineDefinitionOptions) (response *core.DetailedResponse, err error) {
	response, err = cdTektonPipeline.DeleteTektonPipelineDefinitionWithContext(context.Background(), deleteTektonPipelineDefinitionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTektonPipelineDefinitionWithContext is an alternate form of the DeleteTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineDefinitionWithContext(ctx context.Context, deleteTektonPipelineDefinitionOptions *DeleteTektonPipelineDefinitionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineDefinitionOptions, "deleteTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineDefinitionOptions, "deleteTektonPipelineDefinitionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *deleteTektonPipelineDefinitionOptions.PipelineID,
		"definition_id": *deleteTektonPipelineDefinitionOptions.DefinitionID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/definitions/{definition_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTektonPipelineDefinitionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DeleteTektonPipelineDefinition")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tekton_pipeline_definition", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListTektonPipelineProperties : List the pipeline's environment properties
// This request lists the environment properties of the pipeline identified by  `{pipeline_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineProperties(listTektonPipelinePropertiesOptions *ListTektonPipelinePropertiesOptions) (result *PropertiesCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ListTektonPipelinePropertiesWithContext(context.Background(), listTektonPipelinePropertiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTektonPipelinePropertiesWithContext is an alternate form of the ListTektonPipelineProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelinePropertiesWithContext(ctx context.Context, listTektonPipelinePropertiesOptions *ListTektonPipelinePropertiesOptions) (result *PropertiesCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelinePropertiesOptions, "listTektonPipelinePropertiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listTektonPipelinePropertiesOptions, "listTektonPipelinePropertiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *listTektonPipelinePropertiesOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/properties`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTektonPipelinePropertiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ListTektonPipelineProperties")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listTektonPipelinePropertiesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listTektonPipelinePropertiesOptions.Name))
	}
	if listTektonPipelinePropertiesOptions.Type != nil {
		builder.AddQuery("type", strings.Join(listTektonPipelinePropertiesOptions.Type, ","))
	}
	if listTektonPipelinePropertiesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listTektonPipelinePropertiesOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tekton_pipeline_properties", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPropertiesCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineProperties : Create a pipeline environment property
// This request creates an environment property.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineProperties(createTektonPipelinePropertiesOptions *CreateTektonPipelinePropertiesOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CreateTektonPipelinePropertiesWithContext(context.Background(), createTektonPipelinePropertiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTektonPipelinePropertiesWithContext is an alternate form of the CreateTektonPipelineProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelinePropertiesWithContext(ctx context.Context, createTektonPipelinePropertiesOptions *CreateTektonPipelinePropertiesOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelinePropertiesOptions, "createTektonPipelinePropertiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTektonPipelinePropertiesOptions, "createTektonPipelinePropertiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *createTektonPipelinePropertiesOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/properties`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTektonPipelinePropertiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CreateTektonPipelineProperties")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTektonPipelinePropertiesOptions.Name != nil {
		body["name"] = createTektonPipelinePropertiesOptions.Name
	}
	if createTektonPipelinePropertiesOptions.Type != nil {
		body["type"] = createTektonPipelinePropertiesOptions.Type
	}
	if createTektonPipelinePropertiesOptions.Value != nil {
		body["value"] = createTektonPipelinePropertiesOptions.Value
	}
	if createTektonPipelinePropertiesOptions.Enum != nil {
		body["enum"] = createTektonPipelinePropertiesOptions.Enum
	}
	if createTektonPipelinePropertiesOptions.Locked != nil {
		body["locked"] = createTektonPipelinePropertiesOptions.Locked
	}
	if createTektonPipelinePropertiesOptions.Path != nil {
		body["path"] = createTektonPipelinePropertiesOptions.Path
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tekton_pipeline_properties", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineProperty : Get a pipeline environment property
// This request gets the data of an environment property identified by `{property_name}`.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineProperty(getTektonPipelinePropertyOptions *GetTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelinePropertyWithContext(context.Background(), getTektonPipelinePropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelinePropertyWithContext is an alternate form of the GetTektonPipelineProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelinePropertyWithContext(ctx context.Context, getTektonPipelinePropertyOptions *GetTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelinePropertyOptions, "getTektonPipelinePropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelinePropertyOptions, "getTektonPipelinePropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelinePropertyOptions.PipelineID,
		"property_name": *getTektonPipelinePropertyOptions.PropertyName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/properties/{property_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelinePropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineProperty")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTektonPipelineProperty : Replace the value of an environment property
// This request updates the value of an environment property identified by `{property_name}`, its type and name are
// immutable.
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineProperty(replaceTektonPipelinePropertyOptions *ReplaceTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ReplaceTektonPipelinePropertyWithContext(context.Background(), replaceTektonPipelinePropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceTektonPipelinePropertyWithContext is an alternate form of the ReplaceTektonPipelineProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelinePropertyWithContext(ctx context.Context, replaceTektonPipelinePropertyOptions *ReplaceTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTektonPipelinePropertyOptions, "replaceTektonPipelinePropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceTektonPipelinePropertyOptions, "replaceTektonPipelinePropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *replaceTektonPipelinePropertyOptions.PipelineID,
		"property_name": *replaceTektonPipelinePropertyOptions.PropertyName,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/properties/{property_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceTektonPipelinePropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ReplaceTektonPipelineProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceTektonPipelinePropertyOptions.Name != nil {
		body["name"] = replaceTektonPipelinePropertyOptions.Name
	}
	if replaceTektonPipelinePropertyOptions.Type != nil {
		body["type"] = replaceTektonPipelinePropertyOptions.Type
	}
	if replaceTektonPipelinePropertyOptions.Value != nil {
		body["value"] = replaceTektonPipelinePropertyOptions.Value
	}
	if replaceTektonPipelinePropertyOptions.Enum != nil {
		body["enum"] = replaceTektonPipelinePropertyOptions.Enum
	}
	if replaceTektonPipelinePropertyOptions.Locked != nil {
		body["locked"] = replaceTektonPipelinePropertyOptions.Locked
	}
	if replaceTektonPipelinePropertyOptions.Path != nil {
		body["path"] = replaceTektonPipelinePropertyOptions.Path
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_tekton_pipeline_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineProperty : Delete a single pipeline environment property
// This request deletes a single pipeline environment property.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineProperty(deleteTektonPipelinePropertyOptions *DeleteTektonPipelinePropertyOptions) (response *core.DetailedResponse, err error) {
	response, err = cdTektonPipeline.DeleteTektonPipelinePropertyWithContext(context.Background(), deleteTektonPipelinePropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTektonPipelinePropertyWithContext is an alternate form of the DeleteTektonPipelineProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelinePropertyWithContext(ctx context.Context, deleteTektonPipelinePropertyOptions *DeleteTektonPipelinePropertyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelinePropertyOptions, "deleteTektonPipelinePropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTektonPipelinePropertyOptions, "deleteTektonPipelinePropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *deleteTektonPipelinePropertyOptions.PipelineID,
		"property_name": *deleteTektonPipelinePropertyOptions.PropertyName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/properties/{property_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTektonPipelinePropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DeleteTektonPipelineProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tekton_pipeline_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListTektonPipelineTriggers : List pipeline triggers
// This request lists pipeline triggers for the pipeline identified by `{pipeline_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggers(listTektonPipelineTriggersOptions *ListTektonPipelineTriggersOptions) (result *TriggersCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ListTektonPipelineTriggersWithContext(context.Background(), listTektonPipelineTriggersOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTektonPipelineTriggersWithContext is an alternate form of the ListTektonPipelineTriggers method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggersWithContext(ctx context.Context, listTektonPipelineTriggersOptions *ListTektonPipelineTriggersOptions) (result *TriggersCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineTriggersOptions, "listTektonPipelineTriggersOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listTektonPipelineTriggersOptions, "listTektonPipelineTriggersOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *listTektonPipelineTriggersOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTektonPipelineTriggersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ListTektonPipelineTriggers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listTektonPipelineTriggersOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listTektonPipelineTriggersOptions.Type))
	}
	if listTektonPipelineTriggersOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listTektonPipelineTriggersOptions.Name))
	}
	if listTektonPipelineTriggersOptions.EventListener != nil {
		builder.AddQuery("event_listener", fmt.Sprint(*listTektonPipelineTriggersOptions.EventListener))
	}
	if listTektonPipelineTriggersOptions.WorkerID != nil {
		builder.AddQuery("worker.id", fmt.Sprint(*listTektonPipelineTriggersOptions.WorkerID))
	}
	if listTektonPipelineTriggersOptions.WorkerName != nil {
		builder.AddQuery("worker.name", fmt.Sprint(*listTektonPipelineTriggersOptions.WorkerName))
	}
	if listTektonPipelineTriggersOptions.Disabled != nil {
		builder.AddQuery("disabled", fmt.Sprint(*listTektonPipelineTriggersOptions.Disabled))
	}
	if listTektonPipelineTriggersOptions.Tags != nil {
		builder.AddQuery("tags", fmt.Sprint(*listTektonPipelineTriggersOptions.Tags))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tekton_pipeline_triggers", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggersCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineTrigger : Create a trigger
// This request creates a trigger.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTrigger(createTektonPipelineTriggerOptions *CreateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CreateTektonPipelineTriggerWithContext(context.Background(), createTektonPipelineTriggerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTektonPipelineTriggerWithContext is an alternate form of the CreateTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTriggerWithContext(ctx context.Context, createTektonPipelineTriggerOptions *CreateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineTriggerOptions, "createTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTektonPipelineTriggerOptions, "createTektonPipelineTriggerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *createTektonPipelineTriggerOptions.PipelineID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTektonPipelineTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CreateTektonPipelineTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTektonPipelineTriggerOptions.Type != nil {
		body["type"] = createTektonPipelineTriggerOptions.Type
	}
	if createTektonPipelineTriggerOptions.Name != nil {
		body["name"] = createTektonPipelineTriggerOptions.Name
	}
	if createTektonPipelineTriggerOptions.EventListener != nil {
		body["event_listener"] = createTektonPipelineTriggerOptions.EventListener
	}
	if createTektonPipelineTriggerOptions.Tags != nil {
		body["tags"] = createTektonPipelineTriggerOptions.Tags
	}
	if createTektonPipelineTriggerOptions.Worker != nil {
		body["worker"] = createTektonPipelineTriggerOptions.Worker
	}
	if createTektonPipelineTriggerOptions.MaxConcurrentRuns != nil {
		body["max_concurrent_runs"] = createTektonPipelineTriggerOptions.MaxConcurrentRuns
	}
	if createTektonPipelineTriggerOptions.Enabled != nil {
		body["enabled"] = createTektonPipelineTriggerOptions.Enabled
	}
	if createTektonPipelineTriggerOptions.Secret != nil {
		body["secret"] = createTektonPipelineTriggerOptions.Secret
	}
	if createTektonPipelineTriggerOptions.Cron != nil {
		body["cron"] = createTektonPipelineTriggerOptions.Cron
	}
	if createTektonPipelineTriggerOptions.Timezone != nil {
		body["timezone"] = createTektonPipelineTriggerOptions.Timezone
	}
	if createTektonPipelineTriggerOptions.Source != nil {
		body["source"] = createTektonPipelineTriggerOptions.Source
	}
	if createTektonPipelineTriggerOptions.Events != nil {
		body["events"] = createTektonPipelineTriggerOptions.Events
	}
	if createTektonPipelineTriggerOptions.Filter != nil {
		body["filter"] = createTektonPipelineTriggerOptions.Filter
	}
	if createTektonPipelineTriggerOptions.Favorite != nil {
		body["favorite"] = createTektonPipelineTriggerOptions.Favorite
	}
	if createTektonPipelineTriggerOptions.EnableEventsFromForks != nil {
		body["enable_events_from_forks"] = createTektonPipelineTriggerOptions.EnableEventsFromForks
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tekton_pipeline_trigger", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineTrigger : Get a single trigger
// This request retrieves a single trigger identified by `{trigger_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTrigger(getTektonPipelineTriggerOptions *GetTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineTriggerWithContext(context.Background(), getTektonPipelineTriggerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineTriggerWithContext is an alternate form of the GetTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTriggerWithContext(ctx context.Context, getTektonPipelineTriggerOptions *GetTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineTriggerOptions, "getTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineTriggerOptions, "getTektonPipelineTriggerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelineTriggerOptions.PipelineID,
		"trigger_id": *getTektonPipelineTriggerOptions.TriggerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineTrigger")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_trigger", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateTektonPipelineTrigger : Edit a trigger
// This request changes a single field or many fields of the trigger identified by `{trigger_id}`. Note that some fields
// are immutable, and use `/properties` endpoint to update trigger properties.
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipelineTrigger(updateTektonPipelineTriggerOptions *UpdateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.UpdateTektonPipelineTriggerWithContext(context.Background(), updateTektonPipelineTriggerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateTektonPipelineTriggerWithContext is an alternate form of the UpdateTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipelineTriggerWithContext(ctx context.Context, updateTektonPipelineTriggerOptions *UpdateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTektonPipelineTriggerOptions, "updateTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateTektonPipelineTriggerOptions, "updateTektonPipelineTriggerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *updateTektonPipelineTriggerOptions.PipelineID,
		"trigger_id": *updateTektonPipelineTriggerOptions.TriggerID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateTektonPipelineTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "UpdateTektonPipelineTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	if updateTektonPipelineTriggerOptions.TriggerPatch != nil {
		_, err = builder.SetBodyContentJSON(updateTektonPipelineTriggerOptions.TriggerPatch)
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_tekton_pipeline_trigger", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineTrigger : Delete a single trigger
// This request deletes the trigger identified by `{trigger_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTrigger(deleteTektonPipelineTriggerOptions *DeleteTektonPipelineTriggerOptions) (response *core.DetailedResponse, err error) {
	response, err = cdTektonPipeline.DeleteTektonPipelineTriggerWithContext(context.Background(), deleteTektonPipelineTriggerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTektonPipelineTriggerWithContext is an alternate form of the DeleteTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTriggerWithContext(ctx context.Context, deleteTektonPipelineTriggerOptions *DeleteTektonPipelineTriggerOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineTriggerOptions, "deleteTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineTriggerOptions, "deleteTektonPipelineTriggerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *deleteTektonPipelineTriggerOptions.PipelineID,
		"trigger_id": *deleteTektonPipelineTriggerOptions.TriggerID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTektonPipelineTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DeleteTektonPipelineTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tekton_pipeline_trigger", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DuplicateTektonPipelineTrigger : Duplicate a trigger
// This request duplicates a trigger from an existing trigger identified by `{source_trigger_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) DuplicateTektonPipelineTrigger(duplicateTektonPipelineTriggerOptions *DuplicateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.DuplicateTektonPipelineTriggerWithContext(context.Background(), duplicateTektonPipelineTriggerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DuplicateTektonPipelineTriggerWithContext is an alternate form of the DuplicateTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DuplicateTektonPipelineTriggerWithContext(ctx context.Context, duplicateTektonPipelineTriggerOptions *DuplicateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(duplicateTektonPipelineTriggerOptions, "duplicateTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(duplicateTektonPipelineTriggerOptions, "duplicateTektonPipelineTriggerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *duplicateTektonPipelineTriggerOptions.PipelineID,
		"source_trigger_id": *duplicateTektonPipelineTriggerOptions.SourceTriggerID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{source_trigger_id}/duplicate`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range duplicateTektonPipelineTriggerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DuplicateTektonPipelineTrigger")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if duplicateTektonPipelineTriggerOptions.Name != nil {
		body["name"] = duplicateTektonPipelineTriggerOptions.Name
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "duplicate_tekton_pipeline_trigger", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListTektonPipelineTriggerProperties : List trigger properties
// This request lists trigger properties for the trigger identified by `{trigger_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggerProperties(listTektonPipelineTriggerPropertiesOptions *ListTektonPipelineTriggerPropertiesOptions) (result *TriggerPropertiesCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ListTektonPipelineTriggerPropertiesWithContext(context.Background(), listTektonPipelineTriggerPropertiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTektonPipelineTriggerPropertiesWithContext is an alternate form of the ListTektonPipelineTriggerProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggerPropertiesWithContext(ctx context.Context, listTektonPipelineTriggerPropertiesOptions *ListTektonPipelineTriggerPropertiesOptions) (result *TriggerPropertiesCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineTriggerPropertiesOptions, "listTektonPipelineTriggerPropertiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listTektonPipelineTriggerPropertiesOptions, "listTektonPipelineTriggerPropertiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *listTektonPipelineTriggerPropertiesOptions.PipelineID,
		"trigger_id": *listTektonPipelineTriggerPropertiesOptions.TriggerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}/properties`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTektonPipelineTriggerPropertiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ListTektonPipelineTriggerProperties")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listTektonPipelineTriggerPropertiesOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listTektonPipelineTriggerPropertiesOptions.Name))
	}
	if listTektonPipelineTriggerPropertiesOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listTektonPipelineTriggerPropertiesOptions.Type))
	}
	if listTektonPipelineTriggerPropertiesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listTektonPipelineTriggerPropertiesOptions.Sort))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tekton_pipeline_trigger_properties", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerPropertiesCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineTriggerProperties : Create a trigger property
// This request creates a property in the trigger identified by `{trigger_id}`.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTriggerProperties(createTektonPipelineTriggerPropertiesOptions *CreateTektonPipelineTriggerPropertiesOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.CreateTektonPipelineTriggerPropertiesWithContext(context.Background(), createTektonPipelineTriggerPropertiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTektonPipelineTriggerPropertiesWithContext is an alternate form of the CreateTektonPipelineTriggerProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTriggerPropertiesWithContext(ctx context.Context, createTektonPipelineTriggerPropertiesOptions *CreateTektonPipelineTriggerPropertiesOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineTriggerPropertiesOptions, "createTektonPipelineTriggerPropertiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTektonPipelineTriggerPropertiesOptions, "createTektonPipelineTriggerPropertiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *createTektonPipelineTriggerPropertiesOptions.PipelineID,
		"trigger_id": *createTektonPipelineTriggerPropertiesOptions.TriggerID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}/properties`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTektonPipelineTriggerPropertiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "CreateTektonPipelineTriggerProperties")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createTektonPipelineTriggerPropertiesOptions.Name != nil {
		body["name"] = createTektonPipelineTriggerPropertiesOptions.Name
	}
	if createTektonPipelineTriggerPropertiesOptions.Type != nil {
		body["type"] = createTektonPipelineTriggerPropertiesOptions.Type
	}
	if createTektonPipelineTriggerPropertiesOptions.Value != nil {
		body["value"] = createTektonPipelineTriggerPropertiesOptions.Value
	}
	if createTektonPipelineTriggerPropertiesOptions.Enum != nil {
		body["enum"] = createTektonPipelineTriggerPropertiesOptions.Enum
	}
	if createTektonPipelineTriggerPropertiesOptions.Path != nil {
		body["path"] = createTektonPipelineTriggerPropertiesOptions.Path
	}
	if createTektonPipelineTriggerPropertiesOptions.Locked != nil {
		body["locked"] = createTektonPipelineTriggerPropertiesOptions.Locked
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tekton_pipeline_trigger_properties", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperty)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineTriggerProperty : Get a trigger property
// This request retrieves a trigger property.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTriggerProperty(getTektonPipelineTriggerPropertyOptions *GetTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.GetTektonPipelineTriggerPropertyWithContext(context.Background(), getTektonPipelineTriggerPropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTektonPipelineTriggerPropertyWithContext is an alternate form of the GetTektonPipelineTriggerProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTriggerPropertyWithContext(ctx context.Context, getTektonPipelineTriggerPropertyOptions *GetTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineTriggerPropertyOptions, "getTektonPipelineTriggerPropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTektonPipelineTriggerPropertyOptions, "getTektonPipelineTriggerPropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *getTektonPipelineTriggerPropertyOptions.PipelineID,
		"trigger_id": *getTektonPipelineTriggerPropertyOptions.TriggerID,
		"property_name": *getTektonPipelineTriggerPropertyOptions.PropertyName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}/properties/{property_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTektonPipelineTriggerPropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "GetTektonPipelineTriggerProperty")
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tekton_pipeline_trigger_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperty)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTektonPipelineTriggerProperty : Replace a trigger property value
// This request updates a trigger property value, type and name are immutable.
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineTriggerProperty(replaceTektonPipelineTriggerPropertyOptions *ReplaceTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	result, response, err = cdTektonPipeline.ReplaceTektonPipelineTriggerPropertyWithContext(context.Background(), replaceTektonPipelineTriggerPropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceTektonPipelineTriggerPropertyWithContext is an alternate form of the ReplaceTektonPipelineTriggerProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineTriggerPropertyWithContext(ctx context.Context, replaceTektonPipelineTriggerPropertyOptions *ReplaceTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTektonPipelineTriggerPropertyOptions, "replaceTektonPipelineTriggerPropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceTektonPipelineTriggerPropertyOptions, "replaceTektonPipelineTriggerPropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *replaceTektonPipelineTriggerPropertyOptions.PipelineID,
		"trigger_id": *replaceTektonPipelineTriggerPropertyOptions.TriggerID,
		"property_name": *replaceTektonPipelineTriggerPropertyOptions.PropertyName,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}/properties/{property_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceTektonPipelineTriggerPropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "ReplaceTektonPipelineTriggerProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if replaceTektonPipelineTriggerPropertyOptions.Name != nil {
		body["name"] = replaceTektonPipelineTriggerPropertyOptions.Name
	}
	if replaceTektonPipelineTriggerPropertyOptions.Type != nil {
		body["type"] = replaceTektonPipelineTriggerPropertyOptions.Type
	}
	if replaceTektonPipelineTriggerPropertyOptions.Value != nil {
		body["value"] = replaceTektonPipelineTriggerPropertyOptions.Value
	}
	if replaceTektonPipelineTriggerPropertyOptions.Enum != nil {
		body["enum"] = replaceTektonPipelineTriggerPropertyOptions.Enum
	}
	if replaceTektonPipelineTriggerPropertyOptions.Path != nil {
		body["path"] = replaceTektonPipelineTriggerPropertyOptions.Path
	}
	if replaceTektonPipelineTriggerPropertyOptions.Locked != nil {
		body["locked"] = replaceTektonPipelineTriggerPropertyOptions.Locked
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_tekton_pipeline_trigger_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperty)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineTriggerProperty : Delete a trigger property
// This request deletes a trigger property.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTriggerProperty(deleteTektonPipelineTriggerPropertyOptions *DeleteTektonPipelineTriggerPropertyOptions) (response *core.DetailedResponse, err error) {
	response, err = cdTektonPipeline.DeleteTektonPipelineTriggerPropertyWithContext(context.Background(), deleteTektonPipelineTriggerPropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTektonPipelineTriggerPropertyWithContext is an alternate form of the DeleteTektonPipelineTriggerProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTriggerPropertyWithContext(ctx context.Context, deleteTektonPipelineTriggerPropertyOptions *DeleteTektonPipelineTriggerPropertyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineTriggerPropertyOptions, "deleteTektonPipelineTriggerPropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineTriggerPropertyOptions, "deleteTektonPipelineTriggerPropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"pipeline_id": *deleteTektonPipelineTriggerPropertyOptions.PipelineID,
		"trigger_id": *deleteTektonPipelineTriggerPropertyOptions.TriggerID,
		"property_name": *deleteTektonPipelineTriggerPropertyOptions.PropertyName,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines/{pipeline_id}/triggers/{trigger_id}/properties/{property_name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTektonPipelineTriggerPropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_tekton_pipeline", "V2", "DeleteTektonPipelineTriggerProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tekton_pipeline_trigger_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "2.0.0")
}

// CancelTektonPipelineRunOptions : The CancelTektonPipelineRun options.
type CancelTektonPipelineRunOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Flag indicating whether the pipeline cancellation action is forced or not.
	Force *bool `json:"force,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCancelTektonPipelineRunOptions : Instantiate CancelTektonPipelineRunOptions
func (*CdTektonPipelineV2) NewCancelTektonPipelineRunOptions(pipelineID string, id string) *CancelTektonPipelineRunOptions {
	return &CancelTektonPipelineRunOptions{
		PipelineID: core.StringPtr(pipelineID),
		ID: core.StringPtr(id),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CancelTektonPipelineRunOptions) SetPipelineID(pipelineID string) *CancelTektonPipelineRunOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetID : Allow user to set ID
func (_options *CancelTektonPipelineRunOptions) SetID(id string) *CancelTektonPipelineRunOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetForce : Allow user to set Force
func (_options *CancelTektonPipelineRunOptions) SetForce(force bool) *CancelTektonPipelineRunOptions {
	_options.Force = core.BoolPtr(force)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CancelTektonPipelineRunOptions) SetHeaders(param map[string]string) *CancelTektonPipelineRunOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineDefinitionOptions : The CreateTektonPipelineDefinition options.
type CreateTektonPipelineDefinitionOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Source repository containing the Tekton pipeline definition.
	Source *DefinitionSource `json:"source" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateTektonPipelineDefinitionOptions : Instantiate CreateTektonPipelineDefinitionOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineDefinitionOptions(pipelineID string, source *DefinitionSource) *CreateTektonPipelineDefinitionOptions {
	return &CreateTektonPipelineDefinitionOptions{
		PipelineID: core.StringPtr(pipelineID),
		Source: source,
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelineDefinitionOptions) SetPipelineID(pipelineID string) *CreateTektonPipelineDefinitionOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetSource : Allow user to set Source
func (_options *CreateTektonPipelineDefinitionOptions) SetSource(source *DefinitionSource) *CreateTektonPipelineDefinitionOptions {
	_options.Source = source
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineDefinitionOptions) SetHeaders(param map[string]string) *CreateTektonPipelineDefinitionOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineOptions : The CreateTektonPipeline options.
type CreateTektonPipelineOptions struct {
	// The ID for the associated pipeline tool, which was already created in the target toolchain. To get the pipeline ID
	// call the toolchain API https://cloud.ibm.com/apidocs/toolchain#list-tools and find the pipeline tool.
	ID *string `json:"id" validate:"required"`

	// Specify the build number that will be used for the next pipeline run. Build numbers can be any positive whole number
	// between 0 and 100000000000000.
	NextBuildNumber *int64 `json:"next_build_number,omitempty"`

	// Flag to enable notifications for this pipeline. If enabled, the Tekton pipeline run events will be published to all
	// the destinations specified by the Slack and Event Notifications integrations in the parent toolchain.
	EnableNotifications *bool `json:"enable_notifications,omitempty"`

	// Flag to enable partial cloning for this pipeline. When partial clone is enabled, only the files contained within the
	// paths specified in definition repositories are read and cloned, this means that symbolic links might not work.
	EnablePartialCloning *bool `json:"enable_partial_cloning,omitempty"`

	// Specify the worker that is to be used to run the trigger, indicated by a worker object with only the worker ID. If
	// not specified or set as `worker: { id: 'public' }`, the IBM Managed shared workers are used.
	Worker *WorkerIdentity `json:"worker,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateTektonPipelineOptions : Instantiate CreateTektonPipelineOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineOptions(id string) *CreateTektonPipelineOptions {
	return &CreateTektonPipelineOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *CreateTektonPipelineOptions) SetID(id string) *CreateTektonPipelineOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetNextBuildNumber : Allow user to set NextBuildNumber
func (_options *CreateTektonPipelineOptions) SetNextBuildNumber(nextBuildNumber int64) *CreateTektonPipelineOptions {
	_options.NextBuildNumber = core.Int64Ptr(nextBuildNumber)
	return _options
}

// SetEnableNotifications : Allow user to set EnableNotifications
func (_options *CreateTektonPipelineOptions) SetEnableNotifications(enableNotifications bool) *CreateTektonPipelineOptions {
	_options.EnableNotifications = core.BoolPtr(enableNotifications)
	return _options
}

// SetEnablePartialCloning : Allow user to set EnablePartialCloning
func (_options *CreateTektonPipelineOptions) SetEnablePartialCloning(enablePartialCloning bool) *CreateTektonPipelineOptions {
	_options.EnablePartialCloning = core.BoolPtr(enablePartialCloning)
	return _options
}

// SetWorker : Allow user to set Worker
func (_options *CreateTektonPipelineOptions) SetWorker(worker *WorkerIdentity) *CreateTektonPipelineOptions {
	_options.Worker = worker
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineOptions) SetHeaders(param map[string]string) *CreateTektonPipelineOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelinePropertiesOptions : The CreateTektonPipelineProperties options.
type CreateTektonPipelinePropertiesOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// Property value. Any string value is valid.
	Value *string `json:"value,omitempty"`

	// Options for `single_select` property type. Only needed when using `single_select` property type.
	Enum []string `json:"enum,omitempty"`

	// When true, this property cannot be overridden by a trigger property or at runtime. Attempting to override it will
	// result in run requests being rejected. The default is false.
	Locked *bool `json:"locked,omitempty"`

	// A dot notation path for `integration` type properties only, to select a value from the tool integration. If left
	// blank the full tool integration data will be used.
	Path *string `json:"path,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateTektonPipelinePropertiesOptions.Type property.
// Property type.
const (
	CreateTektonPipelinePropertiesOptionsTypeAppconfigConst = "appconfig"
	CreateTektonPipelinePropertiesOptionsTypeIntegrationConst = "integration"
	CreateTektonPipelinePropertiesOptionsTypeSecureConst = "secure"
	CreateTektonPipelinePropertiesOptionsTypeSingleSelectConst = "single_select"
	CreateTektonPipelinePropertiesOptionsTypeTextConst = "text"
)

// NewCreateTektonPipelinePropertiesOptions : Instantiate CreateTektonPipelinePropertiesOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelinePropertiesOptions(pipelineID string, name string, typeVar string) *CreateTektonPipelinePropertiesOptions {
	return &CreateTektonPipelinePropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelinePropertiesOptions) SetPipelineID(pipelineID string) *CreateTektonPipelinePropertiesOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateTektonPipelinePropertiesOptions) SetName(name string) *CreateTektonPipelinePropertiesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateTektonPipelinePropertiesOptions) SetType(typeVar string) *CreateTektonPipelinePropertiesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetValue : Allow user to set Value
func (_options *CreateTektonPipelinePropertiesOptions) SetValue(value string) *CreateTektonPipelinePropertiesOptions {
	_options.Value = core.StringPtr(value)
	return _options
}

// SetEnum : Allow user to set Enum
func (_options *CreateTektonPipelinePropertiesOptions) SetEnum(enum []string) *CreateTektonPipelinePropertiesOptions {
	_options.Enum = enum
	return _options
}

// SetLocked : Allow user to set Locked
func (_options *CreateTektonPipelinePropertiesOptions) SetLocked(locked bool) *CreateTektonPipelinePropertiesOptions {
	_options.Locked = core.BoolPtr(locked)
	return _options
}

// SetPath : Allow user to set Path
func (_options *CreateTektonPipelinePropertiesOptions) SetPath(path string) *CreateTektonPipelinePropertiesOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelinePropertiesOptions) SetHeaders(param map[string]string) *CreateTektonPipelinePropertiesOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineRunOptions : The CreateTektonPipelineRun options.
type CreateTektonPipelineRunOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Optional description for the created PipelineRun.
	Description *string `json:"description,omitempty"`

	// Trigger name.
	TriggerName *string `json:"trigger_name,omitempty"`

	// An object containing string values only. It provides additional 'text' properties or overrides existing
	// pipeline/trigger properties that can be used in the created run.
	TriggerProperties map[string]interface{} `json:"trigger_properties,omitempty"`

	// An object containing string values only. It provides additional `secure` properties or overrides existing `secure`
	// pipeline/trigger properties that can be used in the created run.
	SecureTriggerProperties map[string]interface{} `json:"secure_trigger_properties,omitempty"`

	// An object containing string values only that provides the request headers. Use `$(header.header_key_name)` to access
	// it in a TriggerBinding. Most commonly used as part of a Generic Webhook to provide a verification token or signature
	// in the request headers.
	TriggerHeaders map[string]interface{} `json:"trigger_headers,omitempty"`

	// An object that provides the request body. Use `$(body.body_key_name)` to access it in a TriggerBinding. Most
	// commonly used to pass in additional properties or override properties for the pipeline run that is created.
	TriggerBody map[string]interface{} `json:"trigger_body,omitempty"`

	// Trigger details passed when triggering a Tekton pipeline run.
	Trigger *PipelineRunTrigger `json:"trigger,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateTektonPipelineRunOptions : Instantiate CreateTektonPipelineRunOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineRunOptions(pipelineID string) *CreateTektonPipelineRunOptions {
	return &CreateTektonPipelineRunOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelineRunOptions) SetPipelineID(pipelineID string) *CreateTektonPipelineRunOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateTektonPipelineRunOptions) SetDescription(description string) *CreateTektonPipelineRunOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTriggerName : Allow user to set TriggerName
func (_options *CreateTektonPipelineRunOptions) SetTriggerName(triggerName string) *CreateTektonPipelineRunOptions {
	_options.TriggerName = core.StringPtr(triggerName)
	return _options
}

// SetTriggerProperties : Allow user to set TriggerProperties
func (_options *CreateTektonPipelineRunOptions) SetTriggerProperties(triggerProperties map[string]interface{}) *CreateTektonPipelineRunOptions {
	_options.TriggerProperties = triggerProperties
	return _options
}

// SetSecureTriggerProperties : Allow user to set SecureTriggerProperties
func (_options *CreateTektonPipelineRunOptions) SetSecureTriggerProperties(secureTriggerProperties map[string]interface{}) *CreateTektonPipelineRunOptions {
	_options.SecureTriggerProperties = secureTriggerProperties
	return _options
}

// SetTriggerHeaders : Allow user to set TriggerHeaders
func (_options *CreateTektonPipelineRunOptions) SetTriggerHeaders(triggerHeaders map[string]interface{}) *CreateTektonPipelineRunOptions {
	_options.TriggerHeaders = triggerHeaders
	return _options
}

// SetTriggerBody : Allow user to set TriggerBody
func (_options *CreateTektonPipelineRunOptions) SetTriggerBody(triggerBody map[string]interface{}) *CreateTektonPipelineRunOptions {
	_options.TriggerBody = triggerBody
	return _options
}

// SetTrigger : Allow user to set Trigger
func (_options *CreateTektonPipelineRunOptions) SetTrigger(trigger *PipelineRunTrigger) *CreateTektonPipelineRunOptions {
	_options.Trigger = trigger
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineRunOptions) SetHeaders(param map[string]string) *CreateTektonPipelineRunOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineTriggerOptions : The CreateTektonPipelineTrigger options.
type CreateTektonPipelineTriggerOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener" validate:"required"`

	// Trigger tags array.
	Tags []string `json:"tags"`

	// Specify the worker used to run the trigger. Use `worker: { id: 'public' }` to use the IBM Managed workers. The
	// default is to inherit the worker set in the pipeline settings, which can also be explicitly set using `worker: { id:
	// 'inherit' }`.
	Worker *WorkerIdentity `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled
	// for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Flag to check if the trigger is enabled. If omitted the trigger is enabled by default.
	Enabled *bool `json:"enabled,omitempty"`

	// Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.
	Secret *GenericSecret `json:"secret,omitempty"`

	// Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is
	// every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week.
	// Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.
	Cron *string `json:"cron,omitempty"`

	// Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates
	// this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC.
	// Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.
	Timezone *string `json:"timezone,omitempty"`

	// Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the
	// URL of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
	// https://cloud.ibm.com/apidocs/toolchain#list-tools.
	Source *TriggerSourcePrototype `json:"source,omitempty"`

	// Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger
	// listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the
	// 'merge request' term, they correspond to the generic term i.e. 'pull request'.
	Events []string `json:"events"`

	// Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used
	// for event filtering against the Git webhook payloads.
	Filter *string `json:"filter,omitempty"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`

	// Only used for SCM triggers. When enabled, pull request events from forks of the selected repository will trigger a
	// pipeline run.
	EnableEventsFromForks *bool `json:"enable_events_from_forks,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateTektonPipelineTriggerOptions.Type property.
// Trigger type.
const (
	CreateTektonPipelineTriggerOptionsTypeGenericConst = "generic"
	CreateTektonPipelineTriggerOptionsTypeManualConst = "manual"
	CreateTektonPipelineTriggerOptionsTypeScmConst = "scm"
	CreateTektonPipelineTriggerOptionsTypeTimerConst = "timer"
)

// Constants associated with the CreateTektonPipelineTriggerOptions.Events property.
// List of events. Supported options are 'push' Git webhook events, 'pull_request_closed' Git webhook events and
// 'pull_request' for 'open pull request' or 'update pull request' Git webhook events.
const (
	CreateTektonPipelineTriggerOptionsEventsPullRequestConst = "pull_request"
	CreateTektonPipelineTriggerOptionsEventsPullRequestClosedConst = "pull_request_closed"
	CreateTektonPipelineTriggerOptionsEventsPushConst = "push"
)

// NewCreateTektonPipelineTriggerOptions : Instantiate CreateTektonPipelineTriggerOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineTriggerOptions(pipelineID string, typeVar string, name string, eventListener string) *CreateTektonPipelineTriggerOptions {
	return &CreateTektonPipelineTriggerOptions{
		PipelineID: core.StringPtr(pipelineID),
		Type: core.StringPtr(typeVar),
		Name: core.StringPtr(name),
		EventListener: core.StringPtr(eventListener),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelineTriggerOptions) SetPipelineID(pipelineID string) *CreateTektonPipelineTriggerOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateTektonPipelineTriggerOptions) SetType(typeVar string) *CreateTektonPipelineTriggerOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateTektonPipelineTriggerOptions) SetName(name string) *CreateTektonPipelineTriggerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetEventListener : Allow user to set EventListener
func (_options *CreateTektonPipelineTriggerOptions) SetEventListener(eventListener string) *CreateTektonPipelineTriggerOptions {
	_options.EventListener = core.StringPtr(eventListener)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateTektonPipelineTriggerOptions) SetTags(tags []string) *CreateTektonPipelineTriggerOptions {
	_options.Tags = tags
	return _options
}

// SetWorker : Allow user to set Worker
func (_options *CreateTektonPipelineTriggerOptions) SetWorker(worker *WorkerIdentity) *CreateTektonPipelineTriggerOptions {
	_options.Worker = worker
	return _options
}

// SetMaxConcurrentRuns : Allow user to set MaxConcurrentRuns
func (_options *CreateTektonPipelineTriggerOptions) SetMaxConcurrentRuns(maxConcurrentRuns int64) *CreateTektonPipelineTriggerOptions {
	_options.MaxConcurrentRuns = core.Int64Ptr(maxConcurrentRuns)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateTektonPipelineTriggerOptions) SetEnabled(enabled bool) *CreateTektonPipelineTriggerOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetSecret : Allow user to set Secret
func (_options *CreateTektonPipelineTriggerOptions) SetSecret(secret *GenericSecret) *CreateTektonPipelineTriggerOptions {
	_options.Secret = secret
	return _options
}

// SetCron : Allow user to set Cron
func (_options *CreateTektonPipelineTriggerOptions) SetCron(cron string) *CreateTektonPipelineTriggerOptions {
	_options.Cron = core.StringPtr(cron)
	return _options
}

// SetTimezone : Allow user to set Timezone
func (_options *CreateTektonPipelineTriggerOptions) SetTimezone(timezone string) *CreateTektonPipelineTriggerOptions {
	_options.Timezone = core.StringPtr(timezone)
	return _options
}

// SetSource : Allow user to set Source
func (_options *CreateTektonPipelineTriggerOptions) SetSource(source *TriggerSourcePrototype) *CreateTektonPipelineTriggerOptions {
	_options.Source = source
	return _options
}

// SetEvents : Allow user to set Events
func (_options *CreateTektonPipelineTriggerOptions) SetEvents(events []string) *CreateTektonPipelineTriggerOptions {
	_options.Events = events
	return _options
}

// SetFilter : Allow user to set Filter
func (_options *CreateTektonPipelineTriggerOptions) SetFilter(filter string) *CreateTektonPipelineTriggerOptions {
	_options.Filter = core.StringPtr(filter)
	return _options
}

// SetFavorite : Allow user to set Favorite
func (_options *CreateTektonPipelineTriggerOptions) SetFavorite(favorite bool) *CreateTektonPipelineTriggerOptions {
	_options.Favorite = core.BoolPtr(favorite)
	return _options
}

// SetEnableEventsFromForks : Allow user to set EnableEventsFromForks
func (_options *CreateTektonPipelineTriggerOptions) SetEnableEventsFromForks(enableEventsFromForks bool) *CreateTektonPipelineTriggerOptions {
	_options.EnableEventsFromForks = core.BoolPtr(enableEventsFromForks)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *CreateTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineTriggerPropertiesOptions : The CreateTektonPipelineTriggerProperties options.
type CreateTektonPipelineTriggerPropertiesOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// Property value. Any string value is valid.
	Value *string `json:"value,omitempty"`

	// Options for `single_select` property type. Only needed for `single_select` property type.
	Enum []string `json:"enum,omitempty"`

	// A dot notation path for `integration` type properties only, to select a value from the tool integration. If left
	// blank the full tool integration data will be used.
	Path *string `json:"path,omitempty"`

	// When true, this property cannot be overridden at runtime. Attempting to override it will result in run requests
	// being rejected. The default is false.
	Locked *bool `json:"locked,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateTektonPipelineTriggerPropertiesOptions.Type property.
// Property type.
const (
	CreateTektonPipelineTriggerPropertiesOptionsTypeAppconfigConst = "appconfig"
	CreateTektonPipelineTriggerPropertiesOptionsTypeIntegrationConst = "integration"
	CreateTektonPipelineTriggerPropertiesOptionsTypeSecureConst = "secure"
	CreateTektonPipelineTriggerPropertiesOptionsTypeSingleSelectConst = "single_select"
	CreateTektonPipelineTriggerPropertiesOptionsTypeTextConst = "text"
)

// NewCreateTektonPipelineTriggerPropertiesOptions : Instantiate CreateTektonPipelineTriggerPropertiesOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineTriggerPropertiesOptions(pipelineID string, triggerID string, name string, typeVar string) *CreateTektonPipelineTriggerPropertiesOptions {
	return &CreateTektonPipelineTriggerPropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetPipelineID(pipelineID string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetTriggerID(triggerID string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetName(name string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetType(typeVar string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetValue : Allow user to set Value
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetValue(value string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Value = core.StringPtr(value)
	return _options
}

// SetEnum : Allow user to set Enum
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetEnum(enum []string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Enum = enum
	return _options
}

// SetPath : Allow user to set Path
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetPath(path string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetLocked : Allow user to set Locked
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetLocked(locked bool) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Locked = core.BoolPtr(locked)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineTriggerPropertiesOptions) SetHeaders(param map[string]string) *CreateTektonPipelineTriggerPropertiesOptions {
	options.Headers = param
	return options
}

// Definition : Tekton pipeline definition entry object, consisting of a repository url, a repository path and a branch or tag. The
// referenced repository URL must match the URL of a repository tool integration in the parent toolchain. Obtain the
// list of integrations from the toolchain API https://cloud.ibm.com/apidocs/toolchain#list-tools. The branch or tag of
// the definition must match against a corresponding branch or tag in the chosen repository, and the path must match a
// subfolder in the repository.
type Definition struct {
	// Source repository containing the Tekton pipeline definition.
	Source *DefinitionSource `json:"source" validate:"required"`

	// API URL for interacting with the definition.
	Href *string `json:"href,omitempty"`

	// The aggregated definition ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalDefinition unmarshals an instance of Definition from the specified map of raw messages.
func UnmarshalDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Definition)
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalDefinitionSource)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
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

// DefinitionSource : Source repository containing the Tekton pipeline definition.
type DefinitionSource struct {
	// The only supported source type is "git", indicating that the source is a git repository.
	Type *string `json:"type" validate:"required"`

	// Properties of the source, which define the URL of the repository and a branch or tag.
	Properties *DefinitionSourceProperties `json:"properties" validate:"required"`
}

// NewDefinitionSource : Instantiate DefinitionSource (Generic Model Constructor)
func (*CdTektonPipelineV2) NewDefinitionSource(typeVar string, properties *DefinitionSourceProperties) (_model *DefinitionSource, err error) {
	_model = &DefinitionSource{
		Type: core.StringPtr(typeVar),
		Properties: properties,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalDefinitionSource unmarshals an instance of DefinitionSource from the specified map of raw messages.
func UnmarshalDefinitionSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefinitionSource)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalDefinitionSourceProperties)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefinitionSourceProperties : Properties of the source, which define the URL of the repository and a branch or tag.
type DefinitionSourceProperties struct {
	// URL of the definition repository.
	URL *string `json:"url" validate:"required"`

	// A branch from the repo, specify one of branch or tag only.
	Branch *string `json:"branch,omitempty"`

	// A tag from the repo, specify one of branch or tag only.
	Tag *string `json:"tag,omitempty"`

	// The path to the definition's YAML files.
	Path *string `json:"path" validate:"required"`

	// Reference to the repository tool in the parent toolchain.
	Tool *Tool `json:"tool,omitempty"`
}

// NewDefinitionSourceProperties : Instantiate DefinitionSourceProperties (Generic Model Constructor)
func (*CdTektonPipelineV2) NewDefinitionSourceProperties(url string, path string) (_model *DefinitionSourceProperties, err error) {
	_model = &DefinitionSourceProperties{
		URL: core.StringPtr(url),
		Path: core.StringPtr(path),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalDefinitionSourceProperties unmarshals an instance of DefinitionSourceProperties from the specified map of raw messages.
func UnmarshalDefinitionSourceProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefinitionSourceProperties)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		err = core.SDKErrorf(err, "", "branch-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tag", &obj.Tag)
	if err != nil {
		err = core.SDKErrorf(err, "", "tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "tool", &obj.Tool, UnmarshalTool)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefinitionsCollection : Pipeline definitions is a collection of individual definition entries, each entry consists of a repository URL, a
// repository path, and a branch or tag.
type DefinitionsCollection struct {
	// The list of all definitions in the pipeline.
	Definitions []Definition `json:"definitions" validate:"required"`
}

// UnmarshalDefinitionsCollection unmarshals an instance of DefinitionsCollection from the specified map of raw messages.
func UnmarshalDefinitionsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefinitionsCollection)
	err = core.UnmarshalModel(m, "definitions", &obj.Definitions, UnmarshalDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "definitions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTektonPipelineDefinitionOptions : The DeleteTektonPipelineDefinition options.
type DeleteTektonPipelineDefinitionOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The definition ID.
	DefinitionID *string `json:"definition_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTektonPipelineDefinitionOptions : Instantiate DeleteTektonPipelineDefinitionOptions
func (*CdTektonPipelineV2) NewDeleteTektonPipelineDefinitionOptions(pipelineID string, definitionID string) *DeleteTektonPipelineDefinitionOptions {
	return &DeleteTektonPipelineDefinitionOptions{
		PipelineID: core.StringPtr(pipelineID),
		DefinitionID: core.StringPtr(definitionID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *DeleteTektonPipelineDefinitionOptions) SetPipelineID(pipelineID string) *DeleteTektonPipelineDefinitionOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetDefinitionID : Allow user to set DefinitionID
func (_options *DeleteTektonPipelineDefinitionOptions) SetDefinitionID(definitionID string) *DeleteTektonPipelineDefinitionOptions {
	_options.DefinitionID = core.StringPtr(definitionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTektonPipelineDefinitionOptions) SetHeaders(param map[string]string) *DeleteTektonPipelineDefinitionOptions {
	options.Headers = param
	return options
}

// DeleteTektonPipelineOptions : The DeleteTektonPipeline options.
type DeleteTektonPipelineOptions struct {
	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTektonPipelineOptions : Instantiate DeleteTektonPipelineOptions
func (*CdTektonPipelineV2) NewDeleteTektonPipelineOptions(id string) *DeleteTektonPipelineOptions {
	return &DeleteTektonPipelineOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteTektonPipelineOptions) SetID(id string) *DeleteTektonPipelineOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTektonPipelineOptions) SetHeaders(param map[string]string) *DeleteTektonPipelineOptions {
	options.Headers = param
	return options
}

// DeleteTektonPipelinePropertyOptions : The DeleteTektonPipelineProperty options.
type DeleteTektonPipelinePropertyOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The property name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTektonPipelinePropertyOptions : Instantiate DeleteTektonPipelinePropertyOptions
func (*CdTektonPipelineV2) NewDeleteTektonPipelinePropertyOptions(pipelineID string, propertyName string) *DeleteTektonPipelinePropertyOptions {
	return &DeleteTektonPipelinePropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		PropertyName: core.StringPtr(propertyName),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *DeleteTektonPipelinePropertyOptions) SetPipelineID(pipelineID string) *DeleteTektonPipelinePropertyOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetPropertyName : Allow user to set PropertyName
func (_options *DeleteTektonPipelinePropertyOptions) SetPropertyName(propertyName string) *DeleteTektonPipelinePropertyOptions {
	_options.PropertyName = core.StringPtr(propertyName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTektonPipelinePropertyOptions) SetHeaders(param map[string]string) *DeleteTektonPipelinePropertyOptions {
	options.Headers = param
	return options
}

// DeleteTektonPipelineRunOptions : The DeleteTektonPipelineRun options.
type DeleteTektonPipelineRunOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTektonPipelineRunOptions : Instantiate DeleteTektonPipelineRunOptions
func (*CdTektonPipelineV2) NewDeleteTektonPipelineRunOptions(pipelineID string, id string) *DeleteTektonPipelineRunOptions {
	return &DeleteTektonPipelineRunOptions{
		PipelineID: core.StringPtr(pipelineID),
		ID: core.StringPtr(id),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *DeleteTektonPipelineRunOptions) SetPipelineID(pipelineID string) *DeleteTektonPipelineRunOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteTektonPipelineRunOptions) SetID(id string) *DeleteTektonPipelineRunOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTektonPipelineRunOptions) SetHeaders(param map[string]string) *DeleteTektonPipelineRunOptions {
	options.Headers = param
	return options
}

// DeleteTektonPipelineTriggerOptions : The DeleteTektonPipelineTrigger options.
type DeleteTektonPipelineTriggerOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTektonPipelineTriggerOptions : Instantiate DeleteTektonPipelineTriggerOptions
func (*CdTektonPipelineV2) NewDeleteTektonPipelineTriggerOptions(pipelineID string, triggerID string) *DeleteTektonPipelineTriggerOptions {
	return &DeleteTektonPipelineTriggerOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *DeleteTektonPipelineTriggerOptions) SetPipelineID(pipelineID string) *DeleteTektonPipelineTriggerOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *DeleteTektonPipelineTriggerOptions) SetTriggerID(triggerID string) *DeleteTektonPipelineTriggerOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *DeleteTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// DeleteTektonPipelineTriggerPropertyOptions : The DeleteTektonPipelineTriggerProperty options.
type DeleteTektonPipelineTriggerPropertyOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// The property name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTektonPipelineTriggerPropertyOptions : Instantiate DeleteTektonPipelineTriggerPropertyOptions
func (*CdTektonPipelineV2) NewDeleteTektonPipelineTriggerPropertyOptions(pipelineID string, triggerID string, propertyName string) *DeleteTektonPipelineTriggerPropertyOptions {
	return &DeleteTektonPipelineTriggerPropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
		PropertyName: core.StringPtr(propertyName),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *DeleteTektonPipelineTriggerPropertyOptions) SetPipelineID(pipelineID string) *DeleteTektonPipelineTriggerPropertyOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *DeleteTektonPipelineTriggerPropertyOptions) SetTriggerID(triggerID string) *DeleteTektonPipelineTriggerPropertyOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetPropertyName : Allow user to set PropertyName
func (_options *DeleteTektonPipelineTriggerPropertyOptions) SetPropertyName(propertyName string) *DeleteTektonPipelineTriggerPropertyOptions {
	_options.PropertyName = core.StringPtr(propertyName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTektonPipelineTriggerPropertyOptions) SetHeaders(param map[string]string) *DeleteTektonPipelineTriggerPropertyOptions {
	options.Headers = param
	return options
}

// DuplicateTektonPipelineTriggerOptions : The DuplicateTektonPipelineTrigger options.
type DuplicateTektonPipelineTriggerOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The ID of the trigger to duplicate.
	SourceTriggerID *string `json:"source_trigger_id" validate:"required,ne="`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDuplicateTektonPipelineTriggerOptions : Instantiate DuplicateTektonPipelineTriggerOptions
func (*CdTektonPipelineV2) NewDuplicateTektonPipelineTriggerOptions(pipelineID string, sourceTriggerID string, name string) *DuplicateTektonPipelineTriggerOptions {
	return &DuplicateTektonPipelineTriggerOptions{
		PipelineID: core.StringPtr(pipelineID),
		SourceTriggerID: core.StringPtr(sourceTriggerID),
		Name: core.StringPtr(name),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *DuplicateTektonPipelineTriggerOptions) SetPipelineID(pipelineID string) *DuplicateTektonPipelineTriggerOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetSourceTriggerID : Allow user to set SourceTriggerID
func (_options *DuplicateTektonPipelineTriggerOptions) SetSourceTriggerID(sourceTriggerID string) *DuplicateTektonPipelineTriggerOptions {
	_options.SourceTriggerID = core.StringPtr(sourceTriggerID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DuplicateTektonPipelineTriggerOptions) SetName(name string) *DuplicateTektonPipelineTriggerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DuplicateTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *DuplicateTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// GenericSecret : Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.
type GenericSecret struct {
	// Secret type.
	Type *string `json:"type,omitempty"`

	// Secret value, not needed if secret type is `internal_validation`.
	Value *string `json:"value,omitempty"`

	// Secret location, not needed if secret type is `internal_validation`.
	Source *string `json:"source,omitempty"`

	// Secret name, not needed if type is `internal_validation`.
	KeyName *string `json:"key_name,omitempty"`

	// Algorithm used for `digest_matches` secret type. Only needed for `digest_matches` secret type.
	Algorithm *string `json:"algorithm,omitempty"`
}

// Constants associated with the GenericSecret.Type property.
// Secret type.
const (
	GenericSecretTypeDigestMatchesConst = "digest_matches"
	GenericSecretTypeInternalValidationConst = "internal_validation"
	GenericSecretTypeTokenMatchesConst = "token_matches"
)

// Constants associated with the GenericSecret.Source property.
// Secret location, not needed if secret type is `internal_validation`.
const (
	GenericSecretSourceHeaderConst = "header"
	GenericSecretSourcePayloadConst = "payload"
	GenericSecretSourceQueryConst = "query"
)

// Constants associated with the GenericSecret.Algorithm property.
// Algorithm used for `digest_matches` secret type. Only needed for `digest_matches` secret type.
const (
	GenericSecretAlgorithmMd4Const = "md4"
	GenericSecretAlgorithmMd5Const = "md5"
	GenericSecretAlgorithmRipemd160Const = "ripemd160"
	GenericSecretAlgorithmSha1Const = "sha1"
	GenericSecretAlgorithmSha256Const = "sha256"
	GenericSecretAlgorithmSha384Const = "sha384"
	GenericSecretAlgorithmSha512Const = "sha512"
	GenericSecretAlgorithmSha512224Const = "sha512_224"
	GenericSecretAlgorithmSha512256Const = "sha512_256"
)

// UnmarshalGenericSecret unmarshals an instance of GenericSecret from the specified map of raw messages.
func UnmarshalGenericSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GenericSecret)
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
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "key_name", &obj.KeyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "key_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		err = core.SDKErrorf(err, "", "algorithm-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the GenericSecret
func (genericSecret *GenericSecret) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(genericSecret.Type) {
		_patch["type"] = genericSecret.Type
	}
	if !core.IsNil(genericSecret.Value) {
		_patch["value"] = genericSecret.Value
	}
	if !core.IsNil(genericSecret.Source) {
		_patch["source"] = genericSecret.Source
	}
	if !core.IsNil(genericSecret.KeyName) {
		_patch["key_name"] = genericSecret.KeyName
	}
	if !core.IsNil(genericSecret.Algorithm) {
		_patch["algorithm"] = genericSecret.Algorithm
	}

	return
}

// GetTektonPipelineDefinitionOptions : The GetTektonPipelineDefinition options.
type GetTektonPipelineDefinitionOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The definition ID.
	DefinitionID *string `json:"definition_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelineDefinitionOptions : Instantiate GetTektonPipelineDefinitionOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineDefinitionOptions(pipelineID string, definitionID string) *GetTektonPipelineDefinitionOptions {
	return &GetTektonPipelineDefinitionOptions{
		PipelineID: core.StringPtr(pipelineID),
		DefinitionID: core.StringPtr(definitionID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelineDefinitionOptions) SetPipelineID(pipelineID string) *GetTektonPipelineDefinitionOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetDefinitionID : Allow user to set DefinitionID
func (_options *GetTektonPipelineDefinitionOptions) SetDefinitionID(definitionID string) *GetTektonPipelineDefinitionOptions {
	_options.DefinitionID = core.StringPtr(definitionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineDefinitionOptions) SetHeaders(param map[string]string) *GetTektonPipelineDefinitionOptions {
	options.Headers = param
	return options
}

// GetTektonPipelineOptions : The GetTektonPipeline options.
type GetTektonPipelineOptions struct {
	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelineOptions : Instantiate GetTektonPipelineOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineOptions(id string) *GetTektonPipelineOptions {
	return &GetTektonPipelineOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetTektonPipelineOptions) SetID(id string) *GetTektonPipelineOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineOptions) SetHeaders(param map[string]string) *GetTektonPipelineOptions {
	options.Headers = param
	return options
}

// GetTektonPipelinePropertyOptions : The GetTektonPipelineProperty options.
type GetTektonPipelinePropertyOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The property name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelinePropertyOptions : Instantiate GetTektonPipelinePropertyOptions
func (*CdTektonPipelineV2) NewGetTektonPipelinePropertyOptions(pipelineID string, propertyName string) *GetTektonPipelinePropertyOptions {
	return &GetTektonPipelinePropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		PropertyName: core.StringPtr(propertyName),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelinePropertyOptions) SetPipelineID(pipelineID string) *GetTektonPipelinePropertyOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetPropertyName : Allow user to set PropertyName
func (_options *GetTektonPipelinePropertyOptions) SetPropertyName(propertyName string) *GetTektonPipelinePropertyOptions {
	_options.PropertyName = core.StringPtr(propertyName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelinePropertyOptions) SetHeaders(param map[string]string) *GetTektonPipelinePropertyOptions {
	options.Headers = param
	return options
}

// GetTektonPipelineRunLogContentOptions : The GetTektonPipelineRunLogContent options.
type GetTektonPipelineRunLogContentOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The Tekton pipeline run ID.
	PipelineRunID *string `json:"pipeline_run_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelineRunLogContentOptions : Instantiate GetTektonPipelineRunLogContentOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineRunLogContentOptions(pipelineID string, pipelineRunID string, id string) *GetTektonPipelineRunLogContentOptions {
	return &GetTektonPipelineRunLogContentOptions{
		PipelineID: core.StringPtr(pipelineID),
		PipelineRunID: core.StringPtr(pipelineRunID),
		ID: core.StringPtr(id),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelineRunLogContentOptions) SetPipelineID(pipelineID string) *GetTektonPipelineRunLogContentOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetPipelineRunID : Allow user to set PipelineRunID
func (_options *GetTektonPipelineRunLogContentOptions) SetPipelineRunID(pipelineRunID string) *GetTektonPipelineRunLogContentOptions {
	_options.PipelineRunID = core.StringPtr(pipelineRunID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetTektonPipelineRunLogContentOptions) SetID(id string) *GetTektonPipelineRunLogContentOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineRunLogContentOptions) SetHeaders(param map[string]string) *GetTektonPipelineRunLogContentOptions {
	options.Headers = param
	return options
}

// GetTektonPipelineRunLogsOptions : The GetTektonPipelineRunLogs options.
type GetTektonPipelineRunLogsOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelineRunLogsOptions : Instantiate GetTektonPipelineRunLogsOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineRunLogsOptions(pipelineID string, id string) *GetTektonPipelineRunLogsOptions {
	return &GetTektonPipelineRunLogsOptions{
		PipelineID: core.StringPtr(pipelineID),
		ID: core.StringPtr(id),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelineRunLogsOptions) SetPipelineID(pipelineID string) *GetTektonPipelineRunLogsOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetTektonPipelineRunLogsOptions) SetID(id string) *GetTektonPipelineRunLogsOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineRunLogsOptions) SetHeaders(param map[string]string) *GetTektonPipelineRunLogsOptions {
	options.Headers = param
	return options
}

// GetTektonPipelineRunOptions : The GetTektonPipelineRun options.
type GetTektonPipelineRunOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Defines if response includes definition metadata.
	Includes *string `json:"includes,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetTektonPipelineRunOptions.Includes property.
// Defines if response includes definition metadata.
const (
	GetTektonPipelineRunOptionsIncludesDefinitionsConst = "definitions"
)

// NewGetTektonPipelineRunOptions : Instantiate GetTektonPipelineRunOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineRunOptions(pipelineID string, id string) *GetTektonPipelineRunOptions {
	return &GetTektonPipelineRunOptions{
		PipelineID: core.StringPtr(pipelineID),
		ID: core.StringPtr(id),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelineRunOptions) SetPipelineID(pipelineID string) *GetTektonPipelineRunOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetTektonPipelineRunOptions) SetID(id string) *GetTektonPipelineRunOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetIncludes : Allow user to set Includes
func (_options *GetTektonPipelineRunOptions) SetIncludes(includes string) *GetTektonPipelineRunOptions {
	_options.Includes = core.StringPtr(includes)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineRunOptions) SetHeaders(param map[string]string) *GetTektonPipelineRunOptions {
	options.Headers = param
	return options
}

// GetTektonPipelineTriggerOptions : The GetTektonPipelineTrigger options.
type GetTektonPipelineTriggerOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelineTriggerOptions : Instantiate GetTektonPipelineTriggerOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineTriggerOptions(pipelineID string, triggerID string) *GetTektonPipelineTriggerOptions {
	return &GetTektonPipelineTriggerOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelineTriggerOptions) SetPipelineID(pipelineID string) *GetTektonPipelineTriggerOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *GetTektonPipelineTriggerOptions) SetTriggerID(triggerID string) *GetTektonPipelineTriggerOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *GetTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// GetTektonPipelineTriggerPropertyOptions : The GetTektonPipelineTriggerProperty options.
type GetTektonPipelineTriggerPropertyOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// The property name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTektonPipelineTriggerPropertyOptions : Instantiate GetTektonPipelineTriggerPropertyOptions
func (*CdTektonPipelineV2) NewGetTektonPipelineTriggerPropertyOptions(pipelineID string, triggerID string, propertyName string) *GetTektonPipelineTriggerPropertyOptions {
	return &GetTektonPipelineTriggerPropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
		PropertyName: core.StringPtr(propertyName),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *GetTektonPipelineTriggerPropertyOptions) SetPipelineID(pipelineID string) *GetTektonPipelineTriggerPropertyOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *GetTektonPipelineTriggerPropertyOptions) SetTriggerID(triggerID string) *GetTektonPipelineTriggerPropertyOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetPropertyName : Allow user to set PropertyName
func (_options *GetTektonPipelineTriggerPropertyOptions) SetPropertyName(propertyName string) *GetTektonPipelineTriggerPropertyOptions {
	_options.PropertyName = core.StringPtr(propertyName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTektonPipelineTriggerPropertyOptions) SetHeaders(param map[string]string) *GetTektonPipelineTriggerPropertyOptions {
	options.Headers = param
	return options
}

// ListTektonPipelineDefinitionsOptions : The ListTektonPipelineDefinitions options.
type ListTektonPipelineDefinitionsOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListTektonPipelineDefinitionsOptions : Instantiate ListTektonPipelineDefinitionsOptions
func (*CdTektonPipelineV2) NewListTektonPipelineDefinitionsOptions(pipelineID string) *ListTektonPipelineDefinitionsOptions {
	return &ListTektonPipelineDefinitionsOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ListTektonPipelineDefinitionsOptions) SetPipelineID(pipelineID string) *ListTektonPipelineDefinitionsOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTektonPipelineDefinitionsOptions) SetHeaders(param map[string]string) *ListTektonPipelineDefinitionsOptions {
	options.Headers = param
	return options
}

// ListTektonPipelinePropertiesOptions : The ListTektonPipelineProperties options.
type ListTektonPipelinePropertiesOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Filters the collection to resources with the specified property name.
	Name *string `json:"name,omitempty"`

	// Filters the collection to resources with the specified property type.
	Type []string `json:"type,omitempty"`

	// Sorts the returned properties by name, in ascending order using `name` or in descending order using `-name`.
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListTektonPipelinePropertiesOptions.Type property.
// Type of property.
const (
	ListTektonPipelinePropertiesOptionsTypeAppconfigConst = "appconfig"
	ListTektonPipelinePropertiesOptionsTypeIntegrationConst = "integration"
	ListTektonPipelinePropertiesOptionsTypeSecureConst = "secure"
	ListTektonPipelinePropertiesOptionsTypeSingleSelectConst = "single_select"
	ListTektonPipelinePropertiesOptionsTypeTextConst = "text"
)

// NewListTektonPipelinePropertiesOptions : Instantiate ListTektonPipelinePropertiesOptions
func (*CdTektonPipelineV2) NewListTektonPipelinePropertiesOptions(pipelineID string) *ListTektonPipelinePropertiesOptions {
	return &ListTektonPipelinePropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ListTektonPipelinePropertiesOptions) SetPipelineID(pipelineID string) *ListTektonPipelinePropertiesOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListTektonPipelinePropertiesOptions) SetName(name string) *ListTektonPipelinePropertiesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListTektonPipelinePropertiesOptions) SetType(typeVar []string) *ListTektonPipelinePropertiesOptions {
	_options.Type = typeVar
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListTektonPipelinePropertiesOptions) SetSort(sort string) *ListTektonPipelinePropertiesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTektonPipelinePropertiesOptions) SetHeaders(param map[string]string) *ListTektonPipelinePropertiesOptions {
	options.Headers = param
	return options
}

// ListTektonPipelineRunsOptions : The ListTektonPipelineRuns options.
type ListTektonPipelineRunsOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// A page token that identifies the start point of the list of pipeline runs. This value is included in the response
	// body of each request to fetch pipeline runs.
	Start *string `json:"start,omitempty"`

	// The number of pipeline runs to return, sorted by creation time, most recent first.
	Limit *int64 `json:"limit,omitempty"`

	// Filters the collection to resources with the specified status.
	Status *string `json:"status,omitempty"`

	// Filters the collection to resources with the specified trigger name.
	TriggerName *string `json:"trigger.name,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListTektonPipelineRunsOptions.Status property.
// Filters the collection to resources with the specified status.
const (
	ListTektonPipelineRunsOptionsStatusCancelledConst = "cancelled"
	ListTektonPipelineRunsOptionsStatusErrorConst = "error"
	ListTektonPipelineRunsOptionsStatusFailedConst = "failed"
	ListTektonPipelineRunsOptionsStatusPendingConst = "pending"
	ListTektonPipelineRunsOptionsStatusQueuedConst = "queued"
	ListTektonPipelineRunsOptionsStatusRunningConst = "running"
	ListTektonPipelineRunsOptionsStatusSucceededConst = "succeeded"
	ListTektonPipelineRunsOptionsStatusWaitingConst = "waiting"
)

// NewListTektonPipelineRunsOptions : Instantiate ListTektonPipelineRunsOptions
func (*CdTektonPipelineV2) NewListTektonPipelineRunsOptions(pipelineID string) *ListTektonPipelineRunsOptions {
	return &ListTektonPipelineRunsOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ListTektonPipelineRunsOptions) SetPipelineID(pipelineID string) *ListTektonPipelineRunsOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListTektonPipelineRunsOptions) SetStart(start string) *ListTektonPipelineRunsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListTektonPipelineRunsOptions) SetLimit(limit int64) *ListTektonPipelineRunsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStatus : Allow user to set Status
func (_options *ListTektonPipelineRunsOptions) SetStatus(status string) *ListTektonPipelineRunsOptions {
	_options.Status = core.StringPtr(status)
	return _options
}

// SetTriggerName : Allow user to set TriggerName
func (_options *ListTektonPipelineRunsOptions) SetTriggerName(triggerName string) *ListTektonPipelineRunsOptions {
	_options.TriggerName = core.StringPtr(triggerName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTektonPipelineRunsOptions) SetHeaders(param map[string]string) *ListTektonPipelineRunsOptions {
	options.Headers = param
	return options
}

// ListTektonPipelineTriggerPropertiesOptions : The ListTektonPipelineTriggerProperties options.
type ListTektonPipelineTriggerPropertiesOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Filter properties by `name`.
	Name *string `json:"name,omitempty"`

	// Filter properties by `type`. Valid types are `secure`, `text`, `integration`, `single_select`, `appconfig`.
	Type *string `json:"type,omitempty"`

	// Sort properties by name. They can be sorted in ascending order using `name` or in descending order using `-name`.
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListTektonPipelineTriggerPropertiesOptions : Instantiate ListTektonPipelineTriggerPropertiesOptions
func (*CdTektonPipelineV2) NewListTektonPipelineTriggerPropertiesOptions(pipelineID string, triggerID string) *ListTektonPipelineTriggerPropertiesOptions {
	return &ListTektonPipelineTriggerPropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ListTektonPipelineTriggerPropertiesOptions) SetPipelineID(pipelineID string) *ListTektonPipelineTriggerPropertiesOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *ListTektonPipelineTriggerPropertiesOptions) SetTriggerID(triggerID string) *ListTektonPipelineTriggerPropertiesOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListTektonPipelineTriggerPropertiesOptions) SetName(name string) *ListTektonPipelineTriggerPropertiesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListTektonPipelineTriggerPropertiesOptions) SetType(typeVar string) *ListTektonPipelineTriggerPropertiesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListTektonPipelineTriggerPropertiesOptions) SetSort(sort string) *ListTektonPipelineTriggerPropertiesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTektonPipelineTriggerPropertiesOptions) SetHeaders(param map[string]string) *ListTektonPipelineTriggerPropertiesOptions {
	options.Headers = param
	return options
}

// ListTektonPipelineTriggersOptions : The ListTektonPipelineTriggers options.
type ListTektonPipelineTriggersOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Optional filter by "type", accepts a comma separated list of types. Valid types are "manual", "scm", "generic", and
	// "timer".
	Type *string `json:"type,omitempty"`

	// Optional filter by "name", accepts a single string value.
	Name *string `json:"name,omitempty"`

	// Optional filter by "event_listener", accepts a single string value.
	EventListener *string `json:"event_listener,omitempty"`

	// Optional filter by "worker.id", accepts a single string value.
	WorkerID *string `json:"worker.id,omitempty"`

	// Optional filter by "worker.name", accepts a single string value.
	WorkerName *string `json:"worker.name,omitempty"`

	// Optional filter by "disabled" state, possible values are "true" or "false".
	Disabled *string `json:"disabled,omitempty"`

	// Optional filter by "tags", accepts a comma separated list of tags. The response lists triggers having at least one
	// matching tag.
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListTektonPipelineTriggersOptions : Instantiate ListTektonPipelineTriggersOptions
func (*CdTektonPipelineV2) NewListTektonPipelineTriggersOptions(pipelineID string) *ListTektonPipelineTriggersOptions {
	return &ListTektonPipelineTriggersOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ListTektonPipelineTriggersOptions) SetPipelineID(pipelineID string) *ListTektonPipelineTriggersOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetType : Allow user to set Type
func (_options *ListTektonPipelineTriggersOptions) SetType(typeVar string) *ListTektonPipelineTriggersOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListTektonPipelineTriggersOptions) SetName(name string) *ListTektonPipelineTriggersOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetEventListener : Allow user to set EventListener
func (_options *ListTektonPipelineTriggersOptions) SetEventListener(eventListener string) *ListTektonPipelineTriggersOptions {
	_options.EventListener = core.StringPtr(eventListener)
	return _options
}

// SetWorkerID : Allow user to set WorkerID
func (_options *ListTektonPipelineTriggersOptions) SetWorkerID(workerID string) *ListTektonPipelineTriggersOptions {
	_options.WorkerID = core.StringPtr(workerID)
	return _options
}

// SetWorkerName : Allow user to set WorkerName
func (_options *ListTektonPipelineTriggersOptions) SetWorkerName(workerName string) *ListTektonPipelineTriggersOptions {
	_options.WorkerName = core.StringPtr(workerName)
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *ListTektonPipelineTriggersOptions) SetDisabled(disabled string) *ListTektonPipelineTriggersOptions {
	_options.Disabled = core.StringPtr(disabled)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ListTektonPipelineTriggersOptions) SetTags(tags string) *ListTektonPipelineTriggersOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTektonPipelineTriggersOptions) SetHeaders(param map[string]string) *ListTektonPipelineTriggersOptions {
	options.Headers = param
	return options
}

// Log : Log data for Tekton pipeline run steps.
type Log struct {
	// API for getting log content.
	Href *string `json:"href,omitempty"`

	// Step log ID.
	ID *string `json:"id" validate:"required"`

	// <podName>/<containerName> of this log.
	Name *string `json:"name" validate:"required"`
}

// UnmarshalLog unmarshals an instance of Log from the specified map of raw messages.
func UnmarshalLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Log)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LogsCollection : List of pipeline run log objects.
type LogsCollection struct {
	// The list of pipeline run log objects.
	Logs []Log `json:"logs" validate:"required"`
}

// UnmarshalLogsCollection unmarshals an instance of LogsCollection from the specified map of raw messages.
func UnmarshalLogsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LogsCollection)
	err = core.UnmarshalModel(m, "logs", &obj.Logs, UnmarshalLog)
	if err != nil {
		err = core.SDKErrorf(err, "", "logs-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRun : Single Tekton pipeline run object.
type PipelineRun struct {
	// Universally Unique Identifier.
	ID *string `json:"id" validate:"required"`

	// General href URL.
	Href *string `json:"href,omitempty"`

	// Information about the user that triggered a pipeline run. Only included for pipeline runs that were manually
	// triggered.
	UserInfo *UserInfo `json:"user_info,omitempty"`

	// Status of the pipeline run.
	Status *string `json:"status" validate:"required"`

	// The aggregated definition ID.
	DefinitionID *string `json:"definition_id" validate:"required"`

	// Reference to the pipeline definition of a pipeline run.
	Definition *RunDefinition `json:"definition,omitempty"`

	// Optional description for the created PipelineRun.
	Description *string `json:"description,omitempty"`

	// Worker details used in this pipeline run.
	Worker *PipelineRunWorker `json:"worker" validate:"required"`

	// The ID of the pipeline to which this pipeline run belongs.
	PipelineID *string `json:"pipeline_id" validate:"required"`

	// Reference to the pipeline to which a pipeline run belongs.
	Pipeline *RunPipeline `json:"pipeline,omitempty"`

	// Listener name used to start the run.
	ListenerName *string `json:"listener_name" validate:"required"`

	// Tekton pipeline trigger.
	Trigger TriggerIntf `json:"trigger" validate:"required"`

	// Event parameters object in String format that was passed in upon creation of this pipeline run, the contents depends
	// on the type of trigger. For example, the Git event payload is included for Git triggers, or in the case of a manual
	// trigger the override and added properties are included.
	EventParamsBlob *string `json:"event_params_blob" validate:"required"`

	// Trigger headers object in String format that was passed in upon creation of this pipeline run. Omitted if no
	// trigger_headers object was provided when creating the pipeline run.
	TriggerHeaders *string `json:"trigger_headers,omitempty"`

	// Properties used in this Tekton pipeline run. Not included when fetching the list of pipeline runs.
	Properties []Property `json:"properties,omitempty"`

	// Standard RFC 3339 Date Time String.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Standard RFC 3339 Date Time String. Only included if the run has been updated since it was created.
	UpdatedAt *strfmt.DateTime `json:"updated_at,omitempty"`

	// URL for the details page of this pipeline run.
	RunURL *string `json:"run_url" validate:"required"`

	// Error message that provides details when a pipeline run encounters an error.
	ErrorMessage *string `json:"error_message,omitempty"`
}

// Constants associated with the PipelineRun.Status property.
// Status of the pipeline run.
const (
	PipelineRunStatusCancelledConst = "cancelled"
	PipelineRunStatusErrorConst = "error"
	PipelineRunStatusFailedConst = "failed"
	PipelineRunStatusPendingConst = "pending"
	PipelineRunStatusQueuedConst = "queued"
	PipelineRunStatusRunningConst = "running"
	PipelineRunStatusSucceededConst = "succeeded"
	PipelineRunStatusWaitingConst = "waiting"
)

// UnmarshalPipelineRun unmarshals an instance of PipelineRun from the specified map of raw messages.
func UnmarshalPipelineRun(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRun)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_info", &obj.UserInfo, UnmarshalUserInfo)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_info-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "definition_id", &obj.DefinitionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definition", &obj.Definition, UnmarshalRunDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "definition-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalPipelineRunWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pipeline_id", &obj.PipelineID)
	if err != nil {
		err = core.SDKErrorf(err, "", "pipeline_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "pipeline", &obj.Pipeline, UnmarshalRunPipeline)
	if err != nil {
		err = core.SDKErrorf(err, "", "pipeline-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "listener_name", &obj.ListenerName)
	if err != nil {
		err = core.SDKErrorf(err, "", "listener_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "trigger", &obj.Trigger, UnmarshalTrigger)
	if err != nil {
		err = core.SDKErrorf(err, "", "trigger-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_params_blob", &obj.EventParamsBlob)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_params_blob-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trigger_headers", &obj.TriggerHeaders)
	if err != nil {
		err = core.SDKErrorf(err, "", "trigger_headers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "run_url", &obj.RunURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_url-error", common.GetComponentInfo())
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

// PipelineRunTrigger : Trigger details passed when triggering a Tekton pipeline run.
type PipelineRunTrigger struct {
	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// An object containing string values only. It provides additional 'text' properties or overrides existing
	// pipeline/trigger properties that can be used in the created run.
	Properties map[string]interface{} `json:"properties,omitempty"`

	// An object containing string values only. It provides additional `secure` properties or overrides existing `secure`
	// pipeline/trigger properties that can be used in the created run.
	SecureProperties map[string]interface{} `json:"secure_properties,omitempty"`

	// An object containing string values only that provides the request headers. Use `$(header.header_key_name)` to access
	// it in a TriggerBinding. Most commonly used as part of a Generic Webhook to provide a verification token or signature
	// in the request headers.
	HeadersVar map[string]interface{} `json:"headers,omitempty"`

	// An object that provides the request body. Use `$(body.body_key_name)` to access it in a TriggerBinding. Most
	// commonly used to pass in additional properties or override properties for the pipeline run that is created.
	Body map[string]interface{} `json:"body,omitempty"`
}

// NewPipelineRunTrigger : Instantiate PipelineRunTrigger (Generic Model Constructor)
func (*CdTektonPipelineV2) NewPipelineRunTrigger(name string) (_model *PipelineRunTrigger, err error) {
	_model = &PipelineRunTrigger{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPipelineRunTrigger unmarshals an instance of PipelineRunTrigger from the specified map of raw messages.
func UnmarshalPipelineRunTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunTrigger)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "properties", &obj.Properties)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secure_properties", &obj.SecureProperties)
	if err != nil {
		err = core.SDKErrorf(err, "", "secure_properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "headers", &obj.HeadersVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "headers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "body", &obj.Body)
	if err != nil {
		err = core.SDKErrorf(err, "", "body-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRunWorker : Worker details used in this pipeline run.
type PipelineRunWorker struct {
	// Name of the worker. Computed based on the worker ID.
	Name *string `json:"name,omitempty"`

	// The agent ID of the corresponding private worker integration used for this pipeline run.
	AgentID *string `json:"agent_id,omitempty"`

	// The Service ID of the corresponding private worker integration used for this pipeline run.
	ServiceID *string `json:"service_id,omitempty"`

	// Universally Unique Identifier.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalPipelineRunWorker unmarshals an instance of PipelineRunWorker from the specified map of raw messages.
func UnmarshalPipelineRunWorker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunWorker)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "agent_id", &obj.AgentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "agent_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_id-error", common.GetComponentInfo())
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

// PipelineRunsCollection : Tekton pipeline runs object.
type PipelineRunsCollection struct {
	// Tekton pipeline runs list.
	PipelineRuns []PipelineRun `json:"pipeline_runs" validate:"required"`

	// The number of pipeline runs to return, sorted by creation time, most recent first.
	Limit *int64 `json:"limit" validate:"required"`

	// First page of pipeline runs.
	First *RunsFirstPage `json:"first" validate:"required"`

	// Next page of pipeline runs relative to the `start` and `limit` params. Only included when there are more pages
	// available.
	Next *RunsNextPage `json:"next,omitempty"`

	// Last page of pipeline runs relative to the `start` and `limit` params. Only included when the last page has been
	// reached.
	Last *RunsLastPage `json:"last,omitempty"`
}

// UnmarshalPipelineRunsCollection unmarshals an instance of PipelineRunsCollection from the specified map of raw messages.
func UnmarshalPipelineRunsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunsCollection)
	err = core.UnmarshalModel(m, "pipeline_runs", &obj.PipelineRuns, UnmarshalPipelineRun)
	if err != nil {
		err = core.SDKErrorf(err, "", "pipeline_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalRunsFirstPage)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalRunsNextPage)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalRunsLastPage)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PipelineRunsCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	start, err := core.GetQueryParam(resp.Next.Href, "start")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if start == nil {
		return nil, nil
	}
	return start, nil
}

// PropertiesCollection : Pipeline properties object.
type PropertiesCollection struct {
	// Pipeline properties list.
	Properties []Property `json:"properties" validate:"required"`
}

// UnmarshalPropertiesCollection unmarshals an instance of PropertiesCollection from the specified map of raw messages.
func UnmarshalPropertiesCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PropertiesCollection)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Property : Property object.
type Property struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property value. Any string value is valid.
	Value *string `json:"value,omitempty"`

	// API URL for interacting with the property.
	Href *string `json:"href,omitempty"`

	// Options for `single_select` property type. Only needed when using `single_select` property type.
	Enum []string `json:"enum,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// When true, this property cannot be overridden by a trigger property or at runtime. Attempting to override it will
	// result in run requests being rejected. The default is false.
	Locked *bool `json:"locked,omitempty"`

	// A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left
	// blank the full tool integration data will be used.
	Path *string `json:"path,omitempty"`
}

// Constants associated with the Property.Type property.
// Property type.
const (
	PropertyTypeAppconfigConst = "appconfig"
	PropertyTypeIntegrationConst = "integration"
	PropertyTypeSecureConst = "secure"
	PropertyTypeSingleSelectConst = "single_select"
	PropertyTypeTextConst = "text"
)

// UnmarshalProperty unmarshals an instance of Property from the specified map of raw messages.
func UnmarshalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Property)
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
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		err = core.SDKErrorf(err, "", "enum-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceTektonPipelineDefinitionOptions : The ReplaceTektonPipelineDefinition options.
type ReplaceTektonPipelineDefinitionOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The definition ID.
	DefinitionID *string `json:"definition_id" validate:"required,ne="`

	// Source repository containing the Tekton pipeline definition.
	Source *DefinitionSource `json:"source" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewReplaceTektonPipelineDefinitionOptions : Instantiate ReplaceTektonPipelineDefinitionOptions
func (*CdTektonPipelineV2) NewReplaceTektonPipelineDefinitionOptions(pipelineID string, definitionID string, source *DefinitionSource) *ReplaceTektonPipelineDefinitionOptions {
	return &ReplaceTektonPipelineDefinitionOptions{
		PipelineID: core.StringPtr(pipelineID),
		DefinitionID: core.StringPtr(definitionID),
		Source: source,
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ReplaceTektonPipelineDefinitionOptions) SetPipelineID(pipelineID string) *ReplaceTektonPipelineDefinitionOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetDefinitionID : Allow user to set DefinitionID
func (_options *ReplaceTektonPipelineDefinitionOptions) SetDefinitionID(definitionID string) *ReplaceTektonPipelineDefinitionOptions {
	_options.DefinitionID = core.StringPtr(definitionID)
	return _options
}

// SetSource : Allow user to set Source
func (_options *ReplaceTektonPipelineDefinitionOptions) SetSource(source *DefinitionSource) *ReplaceTektonPipelineDefinitionOptions {
	_options.Source = source
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTektonPipelineDefinitionOptions) SetHeaders(param map[string]string) *ReplaceTektonPipelineDefinitionOptions {
	options.Headers = param
	return options
}

// ReplaceTektonPipelinePropertyOptions : The ReplaceTektonPipelineProperty options.
type ReplaceTektonPipelinePropertyOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The property name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// Property value. Any string value is valid.
	Value *string `json:"value,omitempty"`

	// Options for `single_select` property type. Only needed when using `single_select` property type.
	Enum []string `json:"enum,omitempty"`

	// When true, this property cannot be overridden by a trigger property or at runtime. Attempting to override it will
	// result in run requests being rejected. The default is false.
	Locked *bool `json:"locked,omitempty"`

	// A dot notation path for `integration` type properties only, to select a value from the tool integration. If left
	// blank the full tool integration data will be used.
	Path *string `json:"path,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ReplaceTektonPipelinePropertyOptions.Type property.
// Property type.
const (
	ReplaceTektonPipelinePropertyOptionsTypeAppconfigConst = "appconfig"
	ReplaceTektonPipelinePropertyOptionsTypeIntegrationConst = "integration"
	ReplaceTektonPipelinePropertyOptionsTypeSecureConst = "secure"
	ReplaceTektonPipelinePropertyOptionsTypeSingleSelectConst = "single_select"
	ReplaceTektonPipelinePropertyOptionsTypeTextConst = "text"
)

// NewReplaceTektonPipelinePropertyOptions : Instantiate ReplaceTektonPipelinePropertyOptions
func (*CdTektonPipelineV2) NewReplaceTektonPipelinePropertyOptions(pipelineID string, propertyName string, name string, typeVar string) *ReplaceTektonPipelinePropertyOptions {
	return &ReplaceTektonPipelinePropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		PropertyName: core.StringPtr(propertyName),
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ReplaceTektonPipelinePropertyOptions) SetPipelineID(pipelineID string) *ReplaceTektonPipelinePropertyOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetPropertyName : Allow user to set PropertyName
func (_options *ReplaceTektonPipelinePropertyOptions) SetPropertyName(propertyName string) *ReplaceTektonPipelinePropertyOptions {
	_options.PropertyName = core.StringPtr(propertyName)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceTektonPipelinePropertyOptions) SetName(name string) *ReplaceTektonPipelinePropertyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceTektonPipelinePropertyOptions) SetType(typeVar string) *ReplaceTektonPipelinePropertyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetValue : Allow user to set Value
func (_options *ReplaceTektonPipelinePropertyOptions) SetValue(value string) *ReplaceTektonPipelinePropertyOptions {
	_options.Value = core.StringPtr(value)
	return _options
}

// SetEnum : Allow user to set Enum
func (_options *ReplaceTektonPipelinePropertyOptions) SetEnum(enum []string) *ReplaceTektonPipelinePropertyOptions {
	_options.Enum = enum
	return _options
}

// SetLocked : Allow user to set Locked
func (_options *ReplaceTektonPipelinePropertyOptions) SetLocked(locked bool) *ReplaceTektonPipelinePropertyOptions {
	_options.Locked = core.BoolPtr(locked)
	return _options
}

// SetPath : Allow user to set Path
func (_options *ReplaceTektonPipelinePropertyOptions) SetPath(path string) *ReplaceTektonPipelinePropertyOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTektonPipelinePropertyOptions) SetHeaders(param map[string]string) *ReplaceTektonPipelinePropertyOptions {
	options.Headers = param
	return options
}

// ReplaceTektonPipelineTriggerPropertyOptions : The ReplaceTektonPipelineTriggerProperty options.
type ReplaceTektonPipelineTriggerPropertyOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// The property name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// Property value. Any string value is valid.
	Value *string `json:"value,omitempty"`

	// Options for `single_select` property type. Only needed for `single_select` property type.
	Enum []string `json:"enum,omitempty"`

	// A dot notation path for `integration` type properties only, to select a value from the tool integration. If left
	// blank the full tool integration data will be used.
	Path *string `json:"path,omitempty"`

	// When true, this property cannot be overridden at runtime. Attempting to override it will result in run requests
	// being rejected. The default is false.
	Locked *bool `json:"locked,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ReplaceTektonPipelineTriggerPropertyOptions.Type property.
// Property type.
const (
	ReplaceTektonPipelineTriggerPropertyOptionsTypeAppconfigConst = "appconfig"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeIntegrationConst = "integration"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeSecureConst = "secure"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeSingleSelectConst = "single_select"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeTextConst = "text"
)

// NewReplaceTektonPipelineTriggerPropertyOptions : Instantiate ReplaceTektonPipelineTriggerPropertyOptions
func (*CdTektonPipelineV2) NewReplaceTektonPipelineTriggerPropertyOptions(pipelineID string, triggerID string, propertyName string, name string, typeVar string) *ReplaceTektonPipelineTriggerPropertyOptions {
	return &ReplaceTektonPipelineTriggerPropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
		PropertyName: core.StringPtr(propertyName),
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetPipelineID(pipelineID string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetTriggerID(triggerID string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetPropertyName : Allow user to set PropertyName
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetPropertyName(propertyName string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.PropertyName = core.StringPtr(propertyName)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetName(name string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetType(typeVar string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetValue : Allow user to set Value
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetValue(value string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Value = core.StringPtr(value)
	return _options
}

// SetEnum : Allow user to set Enum
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetEnum(enum []string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Enum = enum
	return _options
}

// SetPath : Allow user to set Path
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetPath(path string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetLocked : Allow user to set Locked
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetLocked(locked bool) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Locked = core.BoolPtr(locked)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTektonPipelineTriggerPropertyOptions) SetHeaders(param map[string]string) *ReplaceTektonPipelineTriggerPropertyOptions {
	options.Headers = param
	return options
}

// RerunTektonPipelineRunOptions : The RerunTektonPipelineRun options.
type RerunTektonPipelineRunOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewRerunTektonPipelineRunOptions : Instantiate RerunTektonPipelineRunOptions
func (*CdTektonPipelineV2) NewRerunTektonPipelineRunOptions(pipelineID string, id string) *RerunTektonPipelineRunOptions {
	return &RerunTektonPipelineRunOptions{
		PipelineID: core.StringPtr(pipelineID),
		ID: core.StringPtr(id),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *RerunTektonPipelineRunOptions) SetPipelineID(pipelineID string) *RerunTektonPipelineRunOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetID : Allow user to set ID
func (_options *RerunTektonPipelineRunOptions) SetID(id string) *RerunTektonPipelineRunOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RerunTektonPipelineRunOptions) SetHeaders(param map[string]string) *RerunTektonPipelineRunOptions {
	options.Headers = param
	return options
}

// ResourceGroupReference : The resource group in which the pipeline was created.
type ResourceGroupReference struct {
	// ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalResourceGroupReference unmarshals an instance of ResourceGroupReference from the specified map of raw messages.
func UnmarshalResourceGroupReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceGroupReference)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RunDefinition : Reference to the pipeline definition of a pipeline run.
type RunDefinition struct {
	// The ID of the definition used for a pipeline run.
	ID *string `json:"id,omitempty"`
}

// UnmarshalRunDefinition unmarshals an instance of RunDefinition from the specified map of raw messages.
func UnmarshalRunDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RunDefinition)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RunPipeline : Reference to the pipeline to which a pipeline run belongs.
type RunPipeline struct {
	// The ID of the pipeline to which a pipeline run belongs.
	ID *string `json:"id,omitempty"`
}

// UnmarshalRunPipeline unmarshals an instance of RunPipeline from the specified map of raw messages.
func UnmarshalRunPipeline(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RunPipeline)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RunsFirstPage : First page of pipeline runs.
type RunsFirstPage struct {
	// General href URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalRunsFirstPage unmarshals an instance of RunsFirstPage from the specified map of raw messages.
func UnmarshalRunsFirstPage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RunsFirstPage)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RunsLastPage : Last page of pipeline runs relative to the `start` and `limit` params. Only included when the last page has been
// reached.
type RunsLastPage struct {
	// General href URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalRunsLastPage unmarshals an instance of RunsLastPage from the specified map of raw messages.
func UnmarshalRunsLastPage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RunsLastPage)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RunsNextPage : Next page of pipeline runs relative to the `start` and `limit` params. Only included when there are more pages
// available.
type RunsNextPage struct {
	// General href URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalRunsNextPage unmarshals an instance of RunsNextPage from the specified map of raw messages.
func UnmarshalRunsNextPage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RunsNextPage)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// StepLog : Logs for a Tekton pipeline run step.
type StepLog struct {
	// The raw log content of the step. Only included when fetching an individual log object.
	Data *string `json:"data" validate:"required"`

	// Step log ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalStepLog unmarshals an instance of StepLog from the specified map of raw messages.
func UnmarshalStepLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StepLog)
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
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

// TektonPipeline : Tekton pipeline object.
type TektonPipeline struct {
	// String.
	Name *string `json:"name" validate:"required"`

	// Pipeline status.
	Status *string `json:"status" validate:"required"`

	// The resource group in which the pipeline was created.
	ResourceGroup *ResourceGroupReference `json:"resource_group" validate:"required"`

	// Toolchain object containing references to the parent toolchain.
	Toolchain *ToolchainReference `json:"toolchain" validate:"required"`

	// Universally Unique Identifier.
	ID *string `json:"id" validate:"required"`

	// Definition list.
	Definitions []Definition `json:"definitions" validate:"required"`

	// Tekton pipeline's environment properties.
	Properties []Property `json:"properties" validate:"required"`

	// Standard RFC 3339 Date Time String.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Standard RFC 3339 Date Time String.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Tekton pipeline triggers list.
	Triggers []TriggerIntf `json:"triggers" validate:"required"`

	// Details of the worker used to run the pipeline.
	Worker *Worker `json:"worker" validate:"required"`

	// URL for this pipeline showing the list of pipeline runs.
	RunsURL *string `json:"runs_url" validate:"required"`

	// API URL for interacting with the pipeline.
	Href *string `json:"href,omitempty"`

	// The latest pipeline run build number. If this property is absent, the pipeline hasn't had any pipeline runs.
	BuildNumber *int64 `json:"build_number" validate:"required"`

	// The build number that will be used for the next pipeline run.
	NextBuildNumber *int64 `json:"next_build_number,omitempty"`

	// Flag to enable notifications for this pipeline. If enabled, the Tekton pipeline run events will be published to all
	// the destinations specified by the Slack and Event Notifications integrations in the parent toolchain. If omitted,
	// this feature is disabled by default.
	EnableNotifications *bool `json:"enable_notifications" validate:"required"`

	// Flag to enable partial cloning for this pipeline. When partial clone is enabled, only the files contained within the
	// paths specified in definition repositories are read and cloned, this means that symbolic links might not work. If
	// omitted, this feature is disabled by default.
	EnablePartialCloning *bool `json:"enable_partial_cloning" validate:"required"`

	// Flag to check if the trigger is enabled.
	Enabled *bool `json:"enabled" validate:"required"`
}

// Constants associated with the TektonPipeline.Status property.
// Pipeline status.
const (
	TektonPipelineStatusConfiguredConst = "configured"
	TektonPipelineStatusConfiguringConst = "configuring"
)

// UnmarshalTektonPipeline unmarshals an instance of TektonPipeline from the specified map of raw messages.
func UnmarshalTektonPipeline(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TektonPipeline)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "resource_group", &obj.ResourceGroup, UnmarshalResourceGroupReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "toolchain", &obj.Toolchain, UnmarshalToolchainReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "definitions", &obj.Definitions, UnmarshalDefinition)
	if err != nil {
		err = core.SDKErrorf(err, "", "definitions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "triggers", &obj.Triggers, UnmarshalTrigger)
	if err != nil {
		err = core.SDKErrorf(err, "", "triggers-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "runs_url", &obj.RunsURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "runs_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "build_number", &obj.BuildNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "build_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "next_build_number", &obj.NextBuildNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "next_build_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_notifications", &obj.EnableNotifications)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_notifications-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_partial_cloning", &obj.EnablePartialCloning)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_partial_cloning-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TektonPipelinePatch : Request body used to update this pipeline.
type TektonPipelinePatch struct {
	// Specify the build number that will be used for the next pipeline run. Build numbers can be any positive whole number
	// between 0 and 100000000000000.
	NextBuildNumber *int64 `json:"next_build_number,omitempty"`

	// Flag to enable notifications for this pipeline. If enabled, the Tekton pipeline run events will be published to all
	// the destinations specified by the Slack and Event Notifications integrations in the parent toolchain.
	EnableNotifications *bool `json:"enable_notifications,omitempty"`

	// Flag to enable partial cloning for this pipeline. When partial clone is enabled, only the files contained within the
	// paths specified in definition repositories are read and cloned, this means that symbolic links might not work.
	EnablePartialCloning *bool `json:"enable_partial_cloning,omitempty"`

	// Specify the worker that is to be used to run the trigger, indicated by a worker object with only the worker ID. If
	// not specified or set as `worker: { id: 'public' }`, the IBM Managed shared workers are used.
	Worker *WorkerIdentity `json:"worker,omitempty"`
}

// UnmarshalTektonPipelinePatch unmarshals an instance of TektonPipelinePatch from the specified map of raw messages.
func UnmarshalTektonPipelinePatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TektonPipelinePatch)
	err = core.UnmarshalPrimitive(m, "next_build_number", &obj.NextBuildNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "next_build_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_notifications", &obj.EnableNotifications)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_notifications-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_partial_cloning", &obj.EnablePartialCloning)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_partial_cloning-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorkerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the TektonPipelinePatch
func (tektonPipelinePatch *TektonPipelinePatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(tektonPipelinePatch.NextBuildNumber) {
		_patch["next_build_number"] = tektonPipelinePatch.NextBuildNumber
	}
	if !core.IsNil(tektonPipelinePatch.EnableNotifications) {
		_patch["enable_notifications"] = tektonPipelinePatch.EnableNotifications
	}
	if !core.IsNil(tektonPipelinePatch.EnablePartialCloning) {
		_patch["enable_partial_cloning"] = tektonPipelinePatch.EnablePartialCloning
	}
	if !core.IsNil(tektonPipelinePatch.Worker) {
		_patch["worker"] = tektonPipelinePatch.Worker.asPatch()
	}

	return
}

// Tool : Reference to the repository tool in the parent toolchain.
type Tool struct {
	// ID of the repository tool instance in the parent toolchain.
	ID *string `json:"id" validate:"required"`
}

// NewTool : Instantiate Tool (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTool(id string) (_model *Tool, err error) {
	_model = &Tool{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTool unmarshals an instance of Tool from the specified map of raw messages.
func UnmarshalTool(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tool)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainReference : Toolchain object containing references to the parent toolchain.
type ToolchainReference struct {
	// Universally Unique Identifier.
	ID *string `json:"id" validate:"required"`

	// The CRN for the toolchain that contains the Tekton pipeline.
	CRN *string `json:"crn" validate:"required"`
}

// UnmarshalToolchainReference unmarshals an instance of ToolchainReference from the specified map of raw messages.
func UnmarshalToolchainReference(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainReference)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Trigger : Tekton pipeline trigger.
// Models which "extend" this model:
// - TriggerManualTrigger
// - TriggerScmTrigger
// - TriggerTimerTrigger
// - TriggerGenericTrigger
type Trigger struct {
	// Trigger type.
	Type *string `json:"type,omitempty"`

	// Trigger name.
	Name *string `json:"name,omitempty"`

	// API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.
	Href *string `json:"href,omitempty"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener,omitempty"`

	// The Trigger ID.
	ID *string `json:"id,omitempty"`

	// Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline
	// run.
	Properties []TriggerProperty `json:"properties,omitempty"`

	// Optional trigger tags array.
	Tags []string `json:"tags"`

	// Details of the worker used to run the trigger.
	Worker *Worker `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled
	// for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Flag to check if the trigger is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`

	// When enabled, pull request events from forks of the selected repository will trigger a pipeline run.
	EnableEventsFromForks *bool `json:"enable_events_from_forks,omitempty"`

	// Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the
	// URL of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
	// https://cloud.ibm.com/apidocs/toolchain#list-tools.
	Source *TriggerSource `json:"source,omitempty"`

	// Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger
	// listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the
	// 'merge request' term, they correspond to the generic term i.e. 'pull request'.
	Events []string `json:"events"`

	// Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used
	// for event filtering against the Git webhook payloads.
	Filter *string `json:"filter,omitempty"`

	// Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is
	// every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week.
	// Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.
	Cron *string `json:"cron,omitempty"`

	// Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates
	// this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC.
	// Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.
	Timezone *string `json:"timezone,omitempty"`

	// Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.
	Secret *GenericSecret `json:"secret,omitempty"`

	// Webhook URL that can be used to trigger pipeline runs.
	WebhookURL *string `json:"webhook_url,omitempty"`
}

// Constants associated with the Trigger.Events property.
// List of events. Supported options are 'push' Git webhook events, 'pull_request_closed' Git webhook events and
// 'pull_request' for 'open pull request' or 'update pull request' Git webhook events.
const (
	TriggerEventsPullRequestConst = "pull_request"
	TriggerEventsPullRequestClosedConst = "pull_request_closed"
	TriggerEventsPushConst = "push"
)
func (*Trigger) isaTrigger() bool {
	return true
}

type TriggerIntf interface {
	isaTrigger() bool
}

// UnmarshalTrigger unmarshals an instance of Trigger from the specified map of raw messages.
func UnmarshalTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Trigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_listener-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_concurrent_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "favorite", &obj.Favorite)
	if err != nil {
		err = core.SDKErrorf(err, "", "favorite-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_events_from_forks", &obj.EnableEventsFromForks)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_events_from_forks-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalTriggerSource)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "events", &obj.Events)
	if err != nil {
		err = core.SDKErrorf(err, "", "events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "filter", &obj.Filter)
	if err != nil {
		err = core.SDKErrorf(err, "", "filter-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		err = core.SDKErrorf(err, "", "cron-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		err = core.SDKErrorf(err, "", "timezone-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "webhook_url", &obj.WebhookURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "webhook_url-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerPatch : Tekton pipeline trigger object used for updating the trigger.
type TriggerPatch struct {
	// Trigger type.
	Type *string `json:"type,omitempty"`

	// Trigger name.
	Name *string `json:"name,omitempty"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener,omitempty"`

	// Trigger tags array. Optional tags for the trigger.
	Tags []string `json:"tags"`

	// Specify the worker used to run the trigger. Use `worker: { id: 'public' }` to use the IBM Managed workers. Use
	// `worker: { id: 'inherit' }` to inherit the worker used by the pipeline.
	Worker *WorkerIdentity `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If set to 0 then the custom concurrency limit is
	// disabled for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Defines if this trigger is enabled.
	Enabled *bool `json:"enabled,omitempty"`

	// Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.
	Secret *GenericSecret `json:"secret,omitempty"`

	// Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is
	// every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week.
	// Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.
	Cron *string `json:"cron,omitempty"`

	// Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates
	// this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC.
	// Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.
	Timezone *string `json:"timezone,omitempty"`

	// Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the
	// URL of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
	// https://cloud.ibm.com/apidocs/toolchain#list-tools.
	Source *TriggerSourcePrototype `json:"source,omitempty"`

	// Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger
	// listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the
	// 'merge request' term, they correspond to the generic term i.e. 'pull request'.
	Events []string `json:"events"`

	// Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used
	// for event filtering against the Git webhook payloads.
	Filter *string `json:"filter,omitempty"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`

	// Only used for SCM triggers. When enabled, pull request events from forks of the selected repository will trigger a
	// pipeline run.
	EnableEventsFromForks *bool `json:"enable_events_from_forks,omitempty"`
}

// Constants associated with the TriggerPatch.Type property.
// Trigger type.
const (
	TriggerPatchTypeGenericConst = "generic"
	TriggerPatchTypeManualConst = "manual"
	TriggerPatchTypeScmConst = "scm"
	TriggerPatchTypeTimerConst = "timer"
)

// Constants associated with the TriggerPatch.Events property.
// List of events. Supported options are 'push' Git webhook events, 'pull_request_closed' Git webhook events and
// 'pull_request' for 'open pull request' or 'update pull request' Git webhook events.
const (
	TriggerPatchEventsPullRequestConst = "pull_request"
	TriggerPatchEventsPullRequestClosedConst = "pull_request_closed"
	TriggerPatchEventsPushConst = "push"
)

// UnmarshalTriggerPatch unmarshals an instance of TriggerPatch from the specified map of raw messages.
func UnmarshalTriggerPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerPatch)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_listener-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorkerIdentity)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_concurrent_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		err = core.SDKErrorf(err, "", "cron-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		err = core.SDKErrorf(err, "", "timezone-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalTriggerSourcePrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "events", &obj.Events)
	if err != nil {
		err = core.SDKErrorf(err, "", "events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "filter", &obj.Filter)
	if err != nil {
		err = core.SDKErrorf(err, "", "filter-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "favorite", &obj.Favorite)
	if err != nil {
		err = core.SDKErrorf(err, "", "favorite-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_events_from_forks", &obj.EnableEventsFromForks)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_events_from_forks-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the TriggerPatch
func (triggerPatch *TriggerPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(triggerPatch.Type) {
		_patch["type"] = triggerPatch.Type
	}
	if !core.IsNil(triggerPatch.Name) {
		_patch["name"] = triggerPatch.Name
	}
	if !core.IsNil(triggerPatch.EventListener) {
		_patch["event_listener"] = triggerPatch.EventListener
	}
	if !core.IsNil(triggerPatch.Tags) {
		_patch["tags"] = triggerPatch.Tags
	}
	if !core.IsNil(triggerPatch.Worker) {
		_patch["worker"] = triggerPatch.Worker.asPatch()
	}
	if !core.IsNil(triggerPatch.MaxConcurrentRuns) {
		_patch["max_concurrent_runs"] = triggerPatch.MaxConcurrentRuns
	}
	if !core.IsNil(triggerPatch.Enabled) {
		_patch["enabled"] = triggerPatch.Enabled
	}
	if !core.IsNil(triggerPatch.Secret) {
		_patch["secret"] = triggerPatch.Secret.asPatch()
	}
	if !core.IsNil(triggerPatch.Cron) {
		_patch["cron"] = triggerPatch.Cron
	}
	if !core.IsNil(triggerPatch.Timezone) {
		_patch["timezone"] = triggerPatch.Timezone
	}
	if !core.IsNil(triggerPatch.Source) {
		_patch["source"] = triggerPatch.Source.asPatch()
	}
	if !core.IsNil(triggerPatch.Events) {
		_patch["events"] = triggerPatch.Events
	}
	if !core.IsNil(triggerPatch.Filter) {
		_patch["filter"] = triggerPatch.Filter
	}
	if !core.IsNil(triggerPatch.Favorite) {
		_patch["favorite"] = triggerPatch.Favorite
	}
	if !core.IsNil(triggerPatch.EnableEventsFromForks) {
		_patch["enable_events_from_forks"] = triggerPatch.EnableEventsFromForks
	}

	return
}

// TriggerPropertiesCollection : Trigger properties object.
type TriggerPropertiesCollection struct {
	// Trigger properties list.
	Properties []TriggerProperty `json:"properties" validate:"required"`
}

// UnmarshalTriggerPropertiesCollection unmarshals an instance of TriggerPropertiesCollection from the specified map of raw messages.
func UnmarshalTriggerPropertiesCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerPropertiesCollection)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerProperty : Trigger property object.
type TriggerProperty struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property value. Any string value is valid.
	Value *string `json:"value,omitempty"`

	// API URL for interacting with the trigger property.
	Href *string `json:"href,omitempty"`

	// Options for `single_select` property type. Only needed for `single_select` property type.
	Enum []string `json:"enum,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// A dot notation path for `integration` type properties only, that selects a value from the tool integration. If left
	// blank the full tool integration data will be used.
	Path *string `json:"path,omitempty"`

	// When true, this property cannot be overridden at runtime. Attempting to override it will result in run requests
	// being rejected. The default is false.
	Locked *bool `json:"locked,omitempty"`
}

// Constants associated with the TriggerProperty.Type property.
// Property type.
const (
	TriggerPropertyTypeAppconfigConst = "appconfig"
	TriggerPropertyTypeIntegrationConst = "integration"
	TriggerPropertyTypeSecureConst = "secure"
	TriggerPropertyTypeSingleSelectConst = "single_select"
	TriggerPropertyTypeTextConst = "text"
)

// UnmarshalTriggerProperty unmarshals an instance of TriggerProperty from the specified map of raw messages.
func UnmarshalTriggerProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerProperty)
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
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		err = core.SDKErrorf(err, "", "enum-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "locked", &obj.Locked)
	if err != nil {
		err = core.SDKErrorf(err, "", "locked-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerSource : Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the URL
// of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
// https://cloud.ibm.com/apidocs/toolchain#list-tools.
type TriggerSource struct {
	// The only supported source type is "git", indicating that the source is a git repository.
	Type *string `json:"type" validate:"required"`

	// Properties of the source, which define the URL of the repository and a branch or pattern.
	Properties *TriggerSourceProperties `json:"properties" validate:"required"`
}

// UnmarshalTriggerSource unmarshals an instance of TriggerSource from the specified map of raw messages.
func UnmarshalTriggerSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerSource)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerSourceProperties)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerSourceProperties : Properties of the source, which define the URL of the repository and a branch or pattern.
type TriggerSourceProperties struct {
	// URL of the repository to which the trigger is listening.
	URL *string `json:"url" validate:"required"`

	// Name of a branch from the repo. Only one of branch, pattern, or filter should be specified.
	Branch *string `json:"branch,omitempty"`

	// The pattern of Git branch or tag. You can specify a glob pattern such as '!test' or '*master' to match against
	// multiple tags or branches in the repository.The glob pattern used must conform to Bash 4.3 specifications, see bash
	// documentation for more info: https://www.gnu.org/software/bash/manual/bash.html#Pattern-Matching. Only one of
	// branch, pattern, or filter should be specified.
	Pattern *string `json:"pattern,omitempty"`

	// True if the repository server is not addressable on the public internet. IBM Cloud will not be able to validate the
	// connection details you provide.
	BlindConnection *bool `json:"blind_connection" validate:"required"`

	// Repository webhook ID. It is generated upon trigger creation.
	HookID *string `json:"hook_id,omitempty"`

	// Reference to the repository tool in the parent toolchain.
	Tool *Tool `json:"tool" validate:"required"`
}

// UnmarshalTriggerSourceProperties unmarshals an instance of TriggerSourceProperties from the specified map of raw messages.
func UnmarshalTriggerSourceProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerSourceProperties)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		err = core.SDKErrorf(err, "", "branch-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		err = core.SDKErrorf(err, "", "pattern-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "blind_connection", &obj.BlindConnection)
	if err != nil {
		err = core.SDKErrorf(err, "", "blind_connection-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "hook_id", &obj.HookID)
	if err != nil {
		err = core.SDKErrorf(err, "", "hook_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "tool", &obj.Tool, UnmarshalTool)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerSourcePropertiesPrototype : Properties of the source, which define the URL of the repository and a branch or pattern.
type TriggerSourcePropertiesPrototype struct {
	// URL of the repository to which the trigger is listening.
	URL *string `json:"url" validate:"required"`

	// Name of a branch from the repo. Only one of branch, pattern, or filter should be specified.
	Branch *string `json:"branch,omitempty"`

	// The pattern of Git branch or tag. You can specify a glob pattern such as '!test' or '*master' to match against
	// multiple tags or branches in the repository.The glob pattern used must conform to Bash 4.3 specifications, see bash
	// documentation for more info: https://www.gnu.org/software/bash/manual/bash.html#Pattern-Matching. Only one of
	// branch, pattern, or filter should be specified.
	Pattern *string `json:"pattern,omitempty"`
}

// NewTriggerSourcePropertiesPrototype : Instantiate TriggerSourcePropertiesPrototype (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerSourcePropertiesPrototype(url string) (_model *TriggerSourcePropertiesPrototype, err error) {
	_model = &TriggerSourcePropertiesPrototype{
		URL: core.StringPtr(url),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTriggerSourcePropertiesPrototype unmarshals an instance of TriggerSourcePropertiesPrototype from the specified map of raw messages.
func UnmarshalTriggerSourcePropertiesPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerSourcePropertiesPrototype)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		err = core.SDKErrorf(err, "", "branch-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		err = core.SDKErrorf(err, "", "pattern-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the TriggerSourcePropertiesPrototype
func (triggerSourcePropertiesPrototype *TriggerSourcePropertiesPrototype) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(triggerSourcePropertiesPrototype.URL) {
		_patch["url"] = triggerSourcePropertiesPrototype.URL
	}
	if !core.IsNil(triggerSourcePropertiesPrototype.Branch) {
		_patch["branch"] = triggerSourcePropertiesPrototype.Branch
	}
	if !core.IsNil(triggerSourcePropertiesPrototype.Pattern) {
		_patch["pattern"] = triggerSourcePropertiesPrototype.Pattern
	}

	return
}

// TriggerSourcePrototype : Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the URL
// of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
// https://cloud.ibm.com/apidocs/toolchain#list-tools.
type TriggerSourcePrototype struct {
	// The only supported source type is "git", indicating that the source is a git repository.
	Type *string `json:"type" validate:"required"`

	// Properties of the source, which define the URL of the repository and a branch or pattern.
	Properties *TriggerSourcePropertiesPrototype `json:"properties" validate:"required"`
}

// NewTriggerSourcePrototype : Instantiate TriggerSourcePrototype (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerSourcePrototype(typeVar string, properties *TriggerSourcePropertiesPrototype) (_model *TriggerSourcePrototype, err error) {
	_model = &TriggerSourcePrototype{
		Type: core.StringPtr(typeVar),
		Properties: properties,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTriggerSourcePrototype unmarshals an instance of TriggerSourcePrototype from the specified map of raw messages.
func UnmarshalTriggerSourcePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerSourcePrototype)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerSourcePropertiesPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the TriggerSourcePrototype
func (triggerSourcePrototype *TriggerSourcePrototype) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(triggerSourcePrototype.Type) {
		_patch["type"] = triggerSourcePrototype.Type
	}
	if !core.IsNil(triggerSourcePrototype.Properties) {
		_patch["properties"] = triggerSourcePrototype.Properties.asPatch()
	}

	return
}

// TriggersCollection : Tekton pipeline triggers object.
type TriggersCollection struct {
	// Tekton pipeline triggers list.
	Triggers []TriggerIntf `json:"triggers" validate:"required"`
}

// UnmarshalTriggersCollection unmarshals an instance of TriggersCollection from the specified map of raw messages.
func UnmarshalTriggersCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersCollection)
	err = core.UnmarshalModel(m, "triggers", &obj.Triggers, UnmarshalTrigger)
	if err != nil {
		err = core.SDKErrorf(err, "", "triggers-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateTektonPipelineOptions : The UpdateTektonPipeline options.
type UpdateTektonPipelineOptions struct {
	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// JSON Merge-Patch content for update_tekton_pipeline.
	TektonPipelinePatch map[string]interface{} `json:"TektonPipelinePatch,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateTektonPipelineOptions : Instantiate UpdateTektonPipelineOptions
func (*CdTektonPipelineV2) NewUpdateTektonPipelineOptions(id string) *UpdateTektonPipelineOptions {
	return &UpdateTektonPipelineOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *UpdateTektonPipelineOptions) SetID(id string) *UpdateTektonPipelineOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetTektonPipelinePatch : Allow user to set TektonPipelinePatch
func (_options *UpdateTektonPipelineOptions) SetTektonPipelinePatch(tektonPipelinePatch map[string]interface{}) *UpdateTektonPipelineOptions {
	_options.TektonPipelinePatch = tektonPipelinePatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTektonPipelineOptions) SetHeaders(param map[string]string) *UpdateTektonPipelineOptions {
	options.Headers = param
	return options
}

// UpdateTektonPipelineTriggerOptions : The UpdateTektonPipelineTrigger options.
type UpdateTektonPipelineTriggerOptions struct {
	// The Tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_tekton_pipeline_trigger.
	TriggerPatch map[string]interface{} `json:"TriggerPatch,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateTektonPipelineTriggerOptions : Instantiate UpdateTektonPipelineTriggerOptions
func (*CdTektonPipelineV2) NewUpdateTektonPipelineTriggerOptions(pipelineID string, triggerID string) *UpdateTektonPipelineTriggerOptions {
	return &UpdateTektonPipelineTriggerOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *UpdateTektonPipelineTriggerOptions) SetPipelineID(pipelineID string) *UpdateTektonPipelineTriggerOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTriggerID : Allow user to set TriggerID
func (_options *UpdateTektonPipelineTriggerOptions) SetTriggerID(triggerID string) *UpdateTektonPipelineTriggerOptions {
	_options.TriggerID = core.StringPtr(triggerID)
	return _options
}

// SetTriggerPatch : Allow user to set TriggerPatch
func (_options *UpdateTektonPipelineTriggerOptions) SetTriggerPatch(triggerPatch map[string]interface{}) *UpdateTektonPipelineTriggerOptions {
	_options.TriggerPatch = triggerPatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *UpdateTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// UserInfo : Information about the user that triggered a pipeline run. Only included for pipeline runs that were manually
// triggered.
type UserInfo struct {
	// IBM Cloud IAM ID.
	IamID *string `json:"iam_id" validate:"required"`

	// User email address.
	Sub *string `json:"sub" validate:"required"`
}

// UnmarshalUserInfo unmarshals an instance of UserInfo from the specified map of raw messages.
func UnmarshalUserInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserInfo)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		err = core.SDKErrorf(err, "", "iam_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sub", &obj.Sub)
	if err != nil {
		err = core.SDKErrorf(err, "", "sub-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Worker : Details of the worker used to run the pipeline.
type Worker struct {
	// Name of the worker. Computed based on the worker ID.
	Name *string `json:"name,omitempty"`

	// Type of the worker. Computed based on the worker ID.
	Type *string `json:"type,omitempty"`

	// ID of the worker.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalWorker unmarshals an instance of Worker from the specified map of raw messages.
func UnmarshalWorker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Worker)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
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

// WorkerIdentity : Specify the worker that is to be used to run the trigger, indicated by a worker object with only the worker ID. If
// not specified or set as `worker: { id: 'public' }`, the IBM Managed shared workers are used.
type WorkerIdentity struct {
	// ID of the worker.
	ID *string `json:"id" validate:"required"`
}

// NewWorkerIdentity : Instantiate WorkerIdentity (Generic Model Constructor)
func (*CdTektonPipelineV2) NewWorkerIdentity(id string) (_model *WorkerIdentity, err error) {
	_model = &WorkerIdentity{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalWorkerIdentity unmarshals an instance of WorkerIdentity from the specified map of raw messages.
func UnmarshalWorkerIdentity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkerIdentity)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the WorkerIdentity
func (workerIdentity *WorkerIdentity) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(workerIdentity.ID) {
		_patch["id"] = workerIdentity.ID
	}

	return
}

// TriggerGenericTrigger : Generic webhook trigger, which triggers a pipeline run when the Tekton Pipeline Service receives a POST event with
// secrets.
// This model "extends" Trigger
type TriggerGenericTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.
	Href *string `json:"href,omitempty"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener" validate:"required"`

	// The Trigger ID.
	ID *string `json:"id" validate:"required"`

	// Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline
	// run.
	Properties []TriggerProperty `json:"properties,omitempty"`

	// Optional trigger tags array.
	Tags []string `json:"tags"`

	// Details of the worker used to run the trigger.
	Worker *Worker `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled
	// for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Flag to check if the trigger is enabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`

	// Only needed for Generic Webhook trigger type. The secret is used to start the Generic Webhook trigger.
	Secret *GenericSecret `json:"secret,omitempty"`

	// Webhook URL that can be used to trigger pipeline runs.
	WebhookURL *string `json:"webhook_url,omitempty"`

	// Stores the CEL (Common Expression Language) expression value which is used for event filtering against the webhook
	// payloads.
	Filter *string `json:"filter,omitempty"`
}

func (*TriggerGenericTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerGenericTrigger unmarshals an instance of TriggerGenericTrigger from the specified map of raw messages.
func UnmarshalTriggerGenericTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerGenericTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_listener-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_concurrent_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "favorite", &obj.Favorite)
	if err != nil {
		err = core.SDKErrorf(err, "", "favorite-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "webhook_url", &obj.WebhookURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "webhook_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "filter", &obj.Filter)
	if err != nil {
		err = core.SDKErrorf(err, "", "filter-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerManualTrigger : Manual trigger.
// This model "extends" Trigger
type TriggerManualTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.
	Href *string `json:"href,omitempty"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener" validate:"required"`

	// The Trigger ID.
	ID *string `json:"id" validate:"required"`

	// Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline
	// run.
	Properties []TriggerProperty `json:"properties,omitempty"`

	// Optional trigger tags array.
	Tags []string `json:"tags"`

	// Details of the worker used to run the trigger.
	Worker *Worker `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled
	// for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Flag to check if the trigger is enabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`
}

func (*TriggerManualTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerManualTrigger unmarshals an instance of TriggerManualTrigger from the specified map of raw messages.
func UnmarshalTriggerManualTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerManualTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_listener-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_concurrent_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "favorite", &obj.Favorite)
	if err != nil {
		err = core.SDKErrorf(err, "", "favorite-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerScmTrigger : Git trigger type. It automatically triggers a pipeline run when the Tekton Pipeline Service receives a corresponding
// Git webhook event.
// This model "extends" Trigger
type TriggerScmTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.
	Href *string `json:"href,omitempty"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener" validate:"required"`

	// The Trigger ID.
	ID *string `json:"id" validate:"required"`

	// Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline
	// run.
	Properties []TriggerProperty `json:"properties,omitempty"`

	// Optional trigger tags array.
	Tags []string `json:"tags"`

	// Details of the worker used to run the trigger.
	Worker *Worker `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled
	// for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Flag to check if the trigger is enabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`

	// When enabled, pull request events from forks of the selected repository will trigger a pipeline run.
	EnableEventsFromForks *bool `json:"enable_events_from_forks,omitempty"`

	// Source repository for a Git trigger. Only required for Git triggers. The referenced repository URL must match the
	// URL of a repository tool integration in the parent toolchain. Obtain the list of integrations from the toolchain API
	// https://cloud.ibm.com/apidocs/toolchain#list-tools.
	Source *TriggerSource `json:"source,omitempty"`

	// Either 'events' or 'filter' is required specifically for Git triggers. Stores a list of events that a Git trigger
	// listens to. Choose one or more from 'push', 'pull_request', and 'pull_request_closed'. If SCM repositories use the
	// 'merge request' term, they correspond to the generic term i.e. 'pull request'.
	Events []string `json:"events"`

	// Either 'events' or 'filter' can be used. Stores the CEL (Common Expression Language) expression value which is used
	// for event filtering against the Git webhook payloads.
	Filter *string `json:"filter,omitempty"`
}

// Constants associated with the TriggerScmTrigger.Events property.
// List of events. Supported options are 'push' Git webhook events, 'pull_request_closed' Git webhook events and
// 'pull_request' for 'open pull request' or 'update pull request' Git webhook events.
const (
	TriggerScmTriggerEventsPullRequestConst = "pull_request"
	TriggerScmTriggerEventsPullRequestClosedConst = "pull_request_closed"
	TriggerScmTriggerEventsPushConst = "push"
)

func (*TriggerScmTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerScmTrigger unmarshals an instance of TriggerScmTrigger from the specified map of raw messages.
func UnmarshalTriggerScmTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerScmTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_listener-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_concurrent_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "favorite", &obj.Favorite)
	if err != nil {
		err = core.SDKErrorf(err, "", "favorite-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enable_events_from_forks", &obj.EnableEventsFromForks)
	if err != nil {
		err = core.SDKErrorf(err, "", "enable_events_from_forks-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "source", &obj.Source, UnmarshalTriggerSource)
	if err != nil {
		err = core.SDKErrorf(err, "", "source-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "events", &obj.Events)
	if err != nil {
		err = core.SDKErrorf(err, "", "events-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "filter", &obj.Filter)
	if err != nil {
		err = core.SDKErrorf(err, "", "filter-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerTimerTrigger : Timer trigger, which triggers pipeline runs according to the provided CRON value and timezone.
// This model "extends" Trigger
type TriggerTimerTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// API URL for interacting with the trigger. Only included when fetching the list of pipeline triggers.
	Href *string `json:"href,omitempty"`

	// Event listener name. The name of the event listener to which the trigger is associated. The event listeners are
	// defined in the definition repositories of the Tekton pipeline.
	EventListener *string `json:"event_listener" validate:"required"`

	// The Trigger ID.
	ID *string `json:"id" validate:"required"`

	// Optional trigger properties are used to override or supplement the pipeline properties when triggering a pipeline
	// run.
	Properties []TriggerProperty `json:"properties,omitempty"`

	// Optional trigger tags array.
	Tags []string `json:"tags"`

	// Details of the worker used to run the trigger.
	Worker *Worker `json:"worker,omitempty"`

	// Defines the maximum number of concurrent runs for this trigger. If omitted then the concurrency limit is disabled
	// for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`

	// Flag to check if the trigger is enabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// Mark the trigger as a favorite.
	Favorite *bool `json:"favorite,omitempty"`

	// Only needed for timer triggers. CRON expression that indicates when this trigger will activate. Maximum frequency is
	// every 5 minutes. The string is based on UNIX crontab syntax: minute, hour, day of month, month, day of week.
	// Example: The CRON expression 0 *_/2 * * * - translates to - every 2 hours.
	Cron *string `json:"cron,omitempty"`

	// Only used for timer triggers. Specify the timezone used for this timer trigger, which will ensure the CRON activates
	// this trigger relative to the specified timezone. If no timezone is specified, the default timezone used is UTC.
	// Valid timezones are those listed in the IANA timezone database, https://www.iana.org/time-zones.
	Timezone *string `json:"timezone,omitempty"`
}

func (*TriggerTimerTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerTimerTrigger unmarshals an instance of TriggerTimerTrigger from the specified map of raw messages.
func UnmarshalTriggerTimerTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerTimerTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		err = core.SDKErrorf(err, "", "event_listener-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		err = core.SDKErrorf(err, "", "worker-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		err = core.SDKErrorf(err, "", "max_concurrent_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "favorite", &obj.Favorite)
	if err != nil {
		err = core.SDKErrorf(err, "", "favorite-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		err = core.SDKErrorf(err, "", "cron-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		err = core.SDKErrorf(err, "", "timezone-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

//
// TektonPipelineRunsPager can be used to simplify the use of the "ListTektonPipelineRuns" method.
//
type TektonPipelineRunsPager struct {
	hasNext bool
	options *ListTektonPipelineRunsOptions
	client  *CdTektonPipelineV2
	pageContext struct {
		next *string
	}
}

// NewTektonPipelineRunsPager returns a new TektonPipelineRunsPager instance.
func (cdTektonPipeline *CdTektonPipelineV2) NewTektonPipelineRunsPager(options *ListTektonPipelineRunsOptions) (pager *TektonPipelineRunsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListTektonPipelineRunsOptions = *options
	pager = &TektonPipelineRunsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  cdTektonPipeline,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *TektonPipelineRunsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *TektonPipelineRunsPager) GetNextWithContext(ctx context.Context) (page []PipelineRun, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListTektonPipelineRunsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		var start *string
		start, err = core.GetQueryParam(result.Next.Href, "start")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'start' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.PipelineRuns

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *TektonPipelineRunsPager) GetAllWithContext(ctx context.Context) (allItems []PipelineRun, err error) {
	for pager.HasNext() {
		var nextPage []PipelineRun
		nextPage, err = pager.GetNextWithContext(ctx)
		if err != nil {
			err = core.RepurposeSDKProblem(err, "error-getting-next-page")
			return
		}
		allItems = append(allItems, nextPage...)
	}
	return
}

// GetNext invokes GetNextWithContext() using context.Background() as the Context parameter.
func (pager *TektonPipelineRunsPager) GetNext() (page []PipelineRun, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *TektonPipelineRunsPager) GetAll() (allItems []PipelineRun, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
