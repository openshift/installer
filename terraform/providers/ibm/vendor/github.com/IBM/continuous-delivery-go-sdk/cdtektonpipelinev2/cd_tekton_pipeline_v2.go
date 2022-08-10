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
 * IBM OpenAPI SDK Code Generator Version: 3.50.0-af9e48c4-20220523-163800
 */

// Package cdtektonpipelinev2 : Operations and models for the CdTektonPipelineV2 service
package cdtektonpipelinev2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	common "github.com/IBM/continuous-delivery-go-sdk/common"
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
const DefaultServiceURL = "https://api.us-south.devops.cloud.ibm.com/v2"

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
			return
		}
	}

	cdTektonPipeline, err = NewCdTektonPipelineV2(options)
	if err != nil {
		return
	}

	err = cdTektonPipeline.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = cdTektonPipeline.Service.SetServiceURL(options.URL)
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
		return
	}

	if options.URL != "" {
		err = baseService.SetServiceURL(options.URL)
		if err != nil {
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
		"us-south": "https://api.us-south.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the us-south region.
		"us-east": "https://api.us-east.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the us-east region.
		"eu-de": "https://api.eu-de.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the eu-de region.
		"eu-gb": "https://api.eu-gb.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the eu-gb region.
		"jp-osa": "https://api.jp-osa.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the jp-osa region.
		"jp-tok": "https://api.jp-tok.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the jp-tok region.
		"au-syd": "https://api.au-syd.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the au-syd region.
		"ca-tor": "https://api.ca-tor.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the ca-tor region.
		"br-sao": "https://api.br-sao.devops.cloud.ibm.com/v2", // The host URL for Tekton Pipeline Service in the br-sao region.
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
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
	return cdTektonPipeline.Service.SetServiceURL(url)
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

// CreateTektonPipeline : Create tekton pipeline
// This request creates a tekton pipeline for a tekton pipeline toolchain integration, user has to use the toolchain
// endpoint to create the tekton pipeline toolchain integration first and then use the generated UUID to create the
// tekton pipeline with or without a specified worker.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipeline(createTektonPipelineOptions *CreateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CreateTektonPipelineWithContext(context.Background(), createTektonPipelineOptions)
}

// CreateTektonPipelineWithContext is an alternate form of the CreateTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineWithContext(ctx context.Context, createTektonPipelineOptions *CreateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createTektonPipelineOptions, "createTektonPipelineOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdTektonPipeline.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdTektonPipeline.Service.Options.URL, `/tekton_pipelines`, nil)
	if err != nil {
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
	if createTektonPipelineOptions.Worker != nil {
		body["worker"] = createTektonPipelineOptions.Worker
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTektonPipeline)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipeline : Get tekton pipeline data
// This request retrieves whole tekton pipeline data.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipeline(getTektonPipelineOptions *GetTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineWithContext(context.Background(), getTektonPipelineOptions)
}

// GetTektonPipelineWithContext is an alternate form of the GetTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineWithContext(ctx context.Context, getTektonPipelineOptions *GetTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineOptions, "getTektonPipelineOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineOptions, "getTektonPipelineOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTektonPipeline)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateTektonPipeline : Update tekton pipeline data
// This request updates tekton pipeline data, but you can only change worker ID in this endpoint. Use other endpoints
// such as /definitions, /triggers, and /properties for detailed updated.
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipeline(updateTektonPipelineOptions *UpdateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.UpdateTektonPipelineWithContext(context.Background(), updateTektonPipelineOptions)
}

// UpdateTektonPipelineWithContext is an alternate form of the UpdateTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipelineWithContext(ctx context.Context, updateTektonPipelineOptions *UpdateTektonPipelineOptions) (result *TektonPipeline, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTektonPipelineOptions, "updateTektonPipelineOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateTektonPipelineOptions, "updateTektonPipelineOptions")
	if err != nil {
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
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateTektonPipelineOptions.Worker != nil {
		body["worker"] = updateTektonPipelineOptions.Worker
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTektonPipeline)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipeline : Delete tekton pipeline instance
// This request deletes tekton pipeline instance that associated with the pipeline toolchain integration.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipeline(deleteTektonPipelineOptions *DeleteTektonPipelineOptions) (response *core.DetailedResponse, err error) {
	return cdTektonPipeline.DeleteTektonPipelineWithContext(context.Background(), deleteTektonPipelineOptions)
}

// DeleteTektonPipelineWithContext is an alternate form of the DeleteTektonPipeline method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineWithContext(ctx context.Context, deleteTektonPipelineOptions *DeleteTektonPipelineOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineOptions, "deleteTektonPipelineOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineOptions, "deleteTektonPipelineOptions")
	if err != nil {
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
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)

	return
}

// ListTektonPipelineRuns : List pipeline run records
// This request list pipeline run records, which has data of that run, such as status, user_info, trigger and other
// information. Default limit is 50.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineRuns(listTektonPipelineRunsOptions *ListTektonPipelineRunsOptions) (result *PipelineRuns, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ListTektonPipelineRunsWithContext(context.Background(), listTektonPipelineRunsOptions)
}

// ListTektonPipelineRunsWithContext is an alternate form of the ListTektonPipelineRuns method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineRunsWithContext(ctx context.Context, listTektonPipelineRunsOptions *ListTektonPipelineRunsOptions) (result *PipelineRuns, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineRunsOptions, "listTektonPipelineRunsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTektonPipelineRunsOptions, "listTektonPipelineRunsOptions")
	if err != nil {
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

	if listTektonPipelineRunsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listTektonPipelineRunsOptions.Limit))
	}
	if listTektonPipelineRunsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listTektonPipelineRunsOptions.Offset))
	}
	if listTektonPipelineRunsOptions.Status != nil {
		builder.AddQuery("status", fmt.Sprint(*listTektonPipelineRunsOptions.Status))
	}
	if listTektonPipelineRunsOptions.TriggerName != nil {
		builder.AddQuery("trigger.name", fmt.Sprint(*listTektonPipelineRunsOptions.TriggerName))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRuns)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineRun : Start a trigger to create a pipelineRun
// This request executes a trigger to create a pipelineRun.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineRun(createTektonPipelineRunOptions *CreateTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CreateTektonPipelineRunWithContext(context.Background(), createTektonPipelineRunOptions)
}

// CreateTektonPipelineRunWithContext is an alternate form of the CreateTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineRunWithContext(ctx context.Context, createTektonPipelineRunOptions *CreateTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineRunOptions, "createTektonPipelineRunOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTektonPipelineRunOptions, "createTektonPipelineRunOptions")
	if err != nil {
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
	if createTektonPipelineRunOptions.TriggerName != nil {
		body["trigger_name"] = createTektonPipelineRunOptions.TriggerName
	}
	if createTektonPipelineRunOptions.TriggerProperties != nil {
		body["trigger_properties"] = createTektonPipelineRunOptions.TriggerProperties
	}
	if createTektonPipelineRunOptions.SecureTriggerProperties != nil {
		body["secure_trigger_properties"] = createTektonPipelineRunOptions.SecureTriggerProperties
	}
	if createTektonPipelineRunOptions.TriggerHeader != nil {
		body["trigger_header"] = createTektonPipelineRunOptions.TriggerHeader
	}
	if createTektonPipelineRunOptions.TriggerBody != nil {
		body["trigger_body"] = createTektonPipelineRunOptions.TriggerBody
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineRun : Get a pipeline run record
// This request retrieves detail of requested pipeline run, to get Kubernetes Resources List of this pipeline run use
// endpoint /tekton_pipelines/{pipeline_id}/tekton_pipelinerun_resource_list/{id}.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRun(getTektonPipelineRunOptions *GetTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineRunWithContext(context.Background(), getTektonPipelineRunOptions)
}

// GetTektonPipelineRunWithContext is an alternate form of the GetTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunWithContext(ctx context.Context, getTektonPipelineRunOptions *GetTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineRunOptions, "getTektonPipelineRunOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineRunOptions, "getTektonPipelineRunOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineRun : Delete a pipeline run record
// This request deletes the requested pipeline run record.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineRun(deleteTektonPipelineRunOptions *DeleteTektonPipelineRunOptions) (response *core.DetailedResponse, err error) {
	return cdTektonPipeline.DeleteTektonPipelineRunWithContext(context.Background(), deleteTektonPipelineRunOptions)
}

// DeleteTektonPipelineRunWithContext is an alternate form of the DeleteTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineRunWithContext(ctx context.Context, deleteTektonPipelineRunOptions *DeleteTektonPipelineRunOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineRunOptions, "deleteTektonPipelineRunOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineRunOptions, "deleteTektonPipelineRunOptions")
	if err != nil {
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
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)

	return
}

// CancelTektonPipelineRun : Cancel a pipeline run
// This request cancels a running pipeline run, use 'force' payload in case you can't cancel a pipeline run normally.
func (cdTektonPipeline *CdTektonPipelineV2) CancelTektonPipelineRun(cancelTektonPipelineRunOptions *CancelTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CancelTektonPipelineRunWithContext(context.Background(), cancelTektonPipelineRunOptions)
}

// CancelTektonPipelineRunWithContext is an alternate form of the CancelTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CancelTektonPipelineRunWithContext(ctx context.Context, cancelTektonPipelineRunOptions *CancelTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(cancelTektonPipelineRunOptions, "cancelTektonPipelineRunOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(cancelTektonPipelineRunOptions, "cancelTektonPipelineRunOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// RerunTektonPipelineRun : Rerun a pipeline run
// This request reruns a past pipeline run with same data. Request body isn't allowed.
func (cdTektonPipeline *CdTektonPipelineV2) RerunTektonPipelineRun(rerunTektonPipelineRunOptions *RerunTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.RerunTektonPipelineRunWithContext(context.Background(), rerunTektonPipelineRunOptions)
}

// RerunTektonPipelineRunWithContext is an alternate form of the RerunTektonPipelineRun method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) RerunTektonPipelineRunWithContext(ctx context.Context, rerunTektonPipelineRunOptions *RerunTektonPipelineRunOptions) (result *PipelineRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(rerunTektonPipelineRunOptions, "rerunTektonPipelineRunOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(rerunTektonPipelineRunOptions, "rerunTektonPipelineRunOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRun)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineRunLogs : Get a list of pipeline run log IDs
// This request fetches list of log IDs of a pipeline run.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogs(getTektonPipelineRunLogsOptions *GetTektonPipelineRunLogsOptions) (result *PipelineRunLogs, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineRunLogsWithContext(context.Background(), getTektonPipelineRunLogsOptions)
}

// GetTektonPipelineRunLogsWithContext is an alternate form of the GetTektonPipelineRunLogs method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogsWithContext(ctx context.Context, getTektonPipelineRunLogsOptions *GetTektonPipelineRunLogsOptions) (result *PipelineRunLogs, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineRunLogsOptions, "getTektonPipelineRunLogsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineRunLogsOptions, "getTektonPipelineRunLogsOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPipelineRunLogs)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineRunLogContent : Get the log content of a pipeline run step by using the step log ID
// This request retrieves log content of a pipeline run step, to get the log ID use endpoint
// /tekton_pipelines/{pipeline_id}/pipeline_runs/{id}/logs.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogContent(getTektonPipelineRunLogContentOptions *GetTektonPipelineRunLogContentOptions) (result *StepLog, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineRunLogContentWithContext(context.Background(), getTektonPipelineRunLogContentOptions)
}

// GetTektonPipelineRunLogContentWithContext is an alternate form of the GetTektonPipelineRunLogContent method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineRunLogContentWithContext(ctx context.Context, getTektonPipelineRunLogContentOptions *GetTektonPipelineRunLogContentOptions) (result *StepLog, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineRunLogContentOptions, "getTektonPipelineRunLogContentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineRunLogContentOptions, "getTektonPipelineRunLogContentOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalStepLog)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ListTektonPipelineDefinitions : List pipeline definitions
// This request fetches pipeline definitions, which is the a collection of individual definition entries, each entry is
// a composition of a repo url, a repo branch and a repo path.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineDefinitions(listTektonPipelineDefinitionsOptions *ListTektonPipelineDefinitionsOptions) (result *Definitions, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ListTektonPipelineDefinitionsWithContext(context.Background(), listTektonPipelineDefinitionsOptions)
}

// ListTektonPipelineDefinitionsWithContext is an alternate form of the ListTektonPipelineDefinitions method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineDefinitionsWithContext(ctx context.Context, listTektonPipelineDefinitionsOptions *ListTektonPipelineDefinitionsOptions) (result *Definitions, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineDefinitionsOptions, "listTektonPipelineDefinitionsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTektonPipelineDefinitionsOptions, "listTektonPipelineDefinitionsOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinitions)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineDefinition : Create a single definition
// This request adds a single definition.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineDefinition(createTektonPipelineDefinitionOptions *CreateTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CreateTektonPipelineDefinitionWithContext(context.Background(), createTektonPipelineDefinitionOptions)
}

// CreateTektonPipelineDefinitionWithContext is an alternate form of the CreateTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineDefinitionWithContext(ctx context.Context, createTektonPipelineDefinitionOptions *CreateTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineDefinitionOptions, "createTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTektonPipelineDefinitionOptions, "createTektonPipelineDefinitionOptions")
	if err != nil {
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
	if createTektonPipelineDefinitionOptions.ScmSource != nil {
		body["scm_source"] = createTektonPipelineDefinitionOptions.ScmSource
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinition)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineDefinition : Retrieve a single definition entry
// This request fetches a single definition entry, which is a composition of the definition repo's url, branch and
// directory path.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineDefinition(getTektonPipelineDefinitionOptions *GetTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineDefinitionWithContext(context.Background(), getTektonPipelineDefinitionOptions)
}

// GetTektonPipelineDefinitionWithContext is an alternate form of the GetTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineDefinitionWithContext(ctx context.Context, getTektonPipelineDefinitionOptions *GetTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineDefinitionOptions, "getTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineDefinitionOptions, "getTektonPipelineDefinitionOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinition)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTektonPipelineDefinition : Edit a single definition entry
// This request replaces a single definition's repo directory path or repo branch. Its scm_source.url and
// service_instance_id are immutable.
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineDefinition(replaceTektonPipelineDefinitionOptions *ReplaceTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ReplaceTektonPipelineDefinitionWithContext(context.Background(), replaceTektonPipelineDefinitionOptions)
}

// ReplaceTektonPipelineDefinitionWithContext is an alternate form of the ReplaceTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineDefinitionWithContext(ctx context.Context, replaceTektonPipelineDefinitionOptions *ReplaceTektonPipelineDefinitionOptions) (result *Definition, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTektonPipelineDefinitionOptions, "replaceTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceTektonPipelineDefinitionOptions, "replaceTektonPipelineDefinitionOptions")
	if err != nil {
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
	if replaceTektonPipelineDefinitionOptions.ScmSource != nil {
		body["scm_source"] = replaceTektonPipelineDefinitionOptions.ScmSource
	}
	if replaceTektonPipelineDefinitionOptions.ServiceInstanceID != nil {
		body["service_instance_id"] = replaceTektonPipelineDefinitionOptions.ServiceInstanceID
	}
	if replaceTektonPipelineDefinitionOptions.ID != nil {
		body["id"] = replaceTektonPipelineDefinitionOptions.ID
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDefinition)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineDefinition : Delete a single definition entry
// This request deletes a single definition from definition list.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineDefinition(deleteTektonPipelineDefinitionOptions *DeleteTektonPipelineDefinitionOptions) (response *core.DetailedResponse, err error) {
	return cdTektonPipeline.DeleteTektonPipelineDefinitionWithContext(context.Background(), deleteTektonPipelineDefinitionOptions)
}

// DeleteTektonPipelineDefinitionWithContext is an alternate form of the DeleteTektonPipelineDefinition method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineDefinitionWithContext(ctx context.Context, deleteTektonPipelineDefinitionOptions *DeleteTektonPipelineDefinitionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineDefinitionOptions, "deleteTektonPipelineDefinitionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineDefinitionOptions, "deleteTektonPipelineDefinitionOptions")
	if err != nil {
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
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)

	return
}

// ListTektonPipelineProperties : List pipeline environment properties
// This request lists the pipeline environment properties which every pipelineRun execution has access to.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineProperties(listTektonPipelinePropertiesOptions *ListTektonPipelinePropertiesOptions) (result *EnvProperties, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ListTektonPipelinePropertiesWithContext(context.Background(), listTektonPipelinePropertiesOptions)
}

// ListTektonPipelinePropertiesWithContext is an alternate form of the ListTektonPipelineProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelinePropertiesWithContext(ctx context.Context, listTektonPipelinePropertiesOptions *ListTektonPipelinePropertiesOptions) (result *EnvProperties, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelinePropertiesOptions, "listTektonPipelinePropertiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTektonPipelinePropertiesOptions, "listTektonPipelinePropertiesOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvProperties)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineProperties : Create pipeline environment property
// This request creates a single pipeline environment property.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineProperties(createTektonPipelinePropertiesOptions *CreateTektonPipelinePropertiesOptions) (result *Property, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CreateTektonPipelinePropertiesWithContext(context.Background(), createTektonPipelinePropertiesOptions)
}

// CreateTektonPipelinePropertiesWithContext is an alternate form of the CreateTektonPipelineProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelinePropertiesWithContext(ctx context.Context, createTektonPipelinePropertiesOptions *CreateTektonPipelinePropertiesOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelinePropertiesOptions, "createTektonPipelinePropertiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTektonPipelinePropertiesOptions, "createTektonPipelinePropertiesOptions")
	if err != nil {
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
	if createTektonPipelinePropertiesOptions.Value != nil {
		body["value"] = createTektonPipelinePropertiesOptions.Value
	}
	if createTektonPipelinePropertiesOptions.Enum != nil {
		body["enum"] = createTektonPipelinePropertiesOptions.Enum
	}
	if createTektonPipelinePropertiesOptions.Default != nil {
		body["default"] = createTektonPipelinePropertiesOptions.Default
	}
	if createTektonPipelinePropertiesOptions.Type != nil {
		body["type"] = createTektonPipelinePropertiesOptions.Type
	}
	if createTektonPipelinePropertiesOptions.Path != nil {
		body["path"] = createTektonPipelinePropertiesOptions.Path
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineProperty : Get a single pipeline environment property
// This request gets a single pipeline environment property data.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineProperty(getTektonPipelinePropertyOptions *GetTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelinePropertyWithContext(context.Background(), getTektonPipelinePropertyOptions)
}

// GetTektonPipelinePropertyWithContext is an alternate form of the GetTektonPipelineProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelinePropertyWithContext(ctx context.Context, getTektonPipelinePropertyOptions *GetTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelinePropertyOptions, "getTektonPipelinePropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelinePropertyOptions, "getTektonPipelinePropertyOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTektonPipelineProperty : Replace a single pipeline environment property value
// This request updates a single pipeline environment property value, its type or name are immutable. For any property
// type, the entered value replaces the existing value.
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineProperty(replaceTektonPipelinePropertyOptions *ReplaceTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ReplaceTektonPipelinePropertyWithContext(context.Background(), replaceTektonPipelinePropertyOptions)
}

// ReplaceTektonPipelinePropertyWithContext is an alternate form of the ReplaceTektonPipelineProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelinePropertyWithContext(ctx context.Context, replaceTektonPipelinePropertyOptions *ReplaceTektonPipelinePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTektonPipelinePropertyOptions, "replaceTektonPipelinePropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceTektonPipelinePropertyOptions, "replaceTektonPipelinePropertyOptions")
	if err != nil {
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
	if replaceTektonPipelinePropertyOptions.Value != nil {
		body["value"] = replaceTektonPipelinePropertyOptions.Value
	}
	if replaceTektonPipelinePropertyOptions.Enum != nil {
		body["enum"] = replaceTektonPipelinePropertyOptions.Enum
	}
	if replaceTektonPipelinePropertyOptions.Default != nil {
		body["default"] = replaceTektonPipelinePropertyOptions.Default
	}
	if replaceTektonPipelinePropertyOptions.Type != nil {
		body["type"] = replaceTektonPipelinePropertyOptions.Type
	}
	if replaceTektonPipelinePropertyOptions.Path != nil {
		body["path"] = replaceTektonPipelinePropertyOptions.Path
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineProperty : Delete a single pipeline environment property
// This request deletes a single pipeline environment property.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineProperty(deleteTektonPipelinePropertyOptions *DeleteTektonPipelinePropertyOptions) (response *core.DetailedResponse, err error) {
	return cdTektonPipeline.DeleteTektonPipelinePropertyWithContext(context.Background(), deleteTektonPipelinePropertyOptions)
}

// DeleteTektonPipelinePropertyWithContext is an alternate form of the DeleteTektonPipelineProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelinePropertyWithContext(ctx context.Context, deleteTektonPipelinePropertyOptions *DeleteTektonPipelinePropertyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelinePropertyOptions, "deleteTektonPipelinePropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTektonPipelinePropertyOptions, "deleteTektonPipelinePropertyOptions")
	if err != nil {
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
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)

	return
}

// ListTektonPipelineTriggers : List pipeline triggers
// This request lists pipeline triggers.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggers(listTektonPipelineTriggersOptions *ListTektonPipelineTriggersOptions) (result *Triggers, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ListTektonPipelineTriggersWithContext(context.Background(), listTektonPipelineTriggersOptions)
}

// ListTektonPipelineTriggersWithContext is an alternate form of the ListTektonPipelineTriggers method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggersWithContext(ctx context.Context, listTektonPipelineTriggersOptions *ListTektonPipelineTriggersOptions) (result *Triggers, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineTriggersOptions, "listTektonPipelineTriggersOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTektonPipelineTriggersOptions, "listTektonPipelineTriggersOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggers)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineTrigger : Create a trigger or duplicate a trigger
// This request creates a trigger or duplicate a trigger from an existing trigger.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTrigger(createTektonPipelineTriggerOptions *CreateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CreateTektonPipelineTriggerWithContext(context.Background(), createTektonPipelineTriggerOptions)
}

// CreateTektonPipelineTriggerWithContext is an alternate form of the CreateTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTriggerWithContext(ctx context.Context, createTektonPipelineTriggerOptions *CreateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineTriggerOptions, "createTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTektonPipelineTriggerOptions, "createTektonPipelineTriggerOptions")
	if err != nil {
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

	if createTektonPipelineTriggerOptions.Trigger != nil {
		_, err = builder.SetBodyContentJSON(createTektonPipelineTriggerOptions.Trigger)
		if err != nil {
			return
		}
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineTrigger : Get a single trigger
// This request retrieves a single trigger detail.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTrigger(getTektonPipelineTriggerOptions *GetTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineTriggerWithContext(context.Background(), getTektonPipelineTriggerOptions)
}

// GetTektonPipelineTriggerWithContext is an alternate form of the GetTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTriggerWithContext(ctx context.Context, getTektonPipelineTriggerOptions *GetTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineTriggerOptions, "getTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineTriggerOptions, "getTektonPipelineTriggerOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// UpdateTektonPipelineTrigger : Edit a single trigger entry
// This request changes a single field or many fields of a trigger, note that some fields are immutable, and use
// /properties to update trigger properties.
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipelineTrigger(updateTektonPipelineTriggerOptions *UpdateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.UpdateTektonPipelineTriggerWithContext(context.Background(), updateTektonPipelineTriggerOptions)
}

// UpdateTektonPipelineTriggerWithContext is an alternate form of the UpdateTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) UpdateTektonPipelineTriggerWithContext(ctx context.Context, updateTektonPipelineTriggerOptions *UpdateTektonPipelineTriggerOptions) (result TriggerIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateTektonPipelineTriggerOptions, "updateTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateTektonPipelineTriggerOptions, "updateTektonPipelineTriggerOptions")
	if err != nil {
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
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateTektonPipelineTriggerOptions.Type != nil {
		body["type"] = updateTektonPipelineTriggerOptions.Type
	}
	if updateTektonPipelineTriggerOptions.Name != nil {
		body["name"] = updateTektonPipelineTriggerOptions.Name
	}
	if updateTektonPipelineTriggerOptions.EventListener != nil {
		body["event_listener"] = updateTektonPipelineTriggerOptions.EventListener
	}
	if updateTektonPipelineTriggerOptions.Tags != nil {
		body["tags"] = updateTektonPipelineTriggerOptions.Tags
	}
	if updateTektonPipelineTriggerOptions.Worker != nil {
		body["worker"] = updateTektonPipelineTriggerOptions.Worker
	}
	if updateTektonPipelineTriggerOptions.Concurrency != nil {
		body["concurrency"] = updateTektonPipelineTriggerOptions.Concurrency
	}
	if updateTektonPipelineTriggerOptions.Disabled != nil {
		body["disabled"] = updateTektonPipelineTriggerOptions.Disabled
	}
	if updateTektonPipelineTriggerOptions.Secret != nil {
		body["secret"] = updateTektonPipelineTriggerOptions.Secret
	}
	if updateTektonPipelineTriggerOptions.Cron != nil {
		body["cron"] = updateTektonPipelineTriggerOptions.Cron
	}
	if updateTektonPipelineTriggerOptions.Timezone != nil {
		body["timezone"] = updateTektonPipelineTriggerOptions.Timezone
	}
	if updateTektonPipelineTriggerOptions.ScmSource != nil {
		body["scm_source"] = updateTektonPipelineTriggerOptions.ScmSource
	}
	if updateTektonPipelineTriggerOptions.Events != nil {
		body["events"] = updateTektonPipelineTriggerOptions.Events
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrigger)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineTrigger : Delete a single trigger
// This request deletes a single trigger.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTrigger(deleteTektonPipelineTriggerOptions *DeleteTektonPipelineTriggerOptions) (response *core.DetailedResponse, err error) {
	return cdTektonPipeline.DeleteTektonPipelineTriggerWithContext(context.Background(), deleteTektonPipelineTriggerOptions)
}

// DeleteTektonPipelineTriggerWithContext is an alternate form of the DeleteTektonPipelineTrigger method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTriggerWithContext(ctx context.Context, deleteTektonPipelineTriggerOptions *DeleteTektonPipelineTriggerOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineTriggerOptions, "deleteTektonPipelineTriggerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineTriggerOptions, "deleteTektonPipelineTriggerOptions")
	if err != nil {
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
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)

	return
}

// ListTektonPipelineTriggerProperties : List trigger environment properties
// This request lists trigger environment properties.
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggerProperties(listTektonPipelineTriggerPropertiesOptions *ListTektonPipelineTriggerPropertiesOptions) (result *TriggerProperties, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ListTektonPipelineTriggerPropertiesWithContext(context.Background(), listTektonPipelineTriggerPropertiesOptions)
}

// ListTektonPipelineTriggerPropertiesWithContext is an alternate form of the ListTektonPipelineTriggerProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ListTektonPipelineTriggerPropertiesWithContext(ctx context.Context, listTektonPipelineTriggerPropertiesOptions *ListTektonPipelineTriggerPropertiesOptions) (result *TriggerProperties, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTektonPipelineTriggerPropertiesOptions, "listTektonPipelineTriggerPropertiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listTektonPipelineTriggerPropertiesOptions, "listTektonPipelineTriggerPropertiesOptions")
	if err != nil {
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

	builder.AddQuery("name", fmt.Sprint(*listTektonPipelineTriggerPropertiesOptions.Name))
	builder.AddQuery("type", fmt.Sprint(*listTektonPipelineTriggerPropertiesOptions.Type))
	builder.AddQuery("sort", fmt.Sprint(*listTektonPipelineTriggerPropertiesOptions.Sort))

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperties)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateTektonPipelineTriggerProperties : Create trigger's environment property
// This request creates a trigger's property.
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTriggerProperties(createTektonPipelineTriggerPropertiesOptions *CreateTektonPipelineTriggerPropertiesOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.CreateTektonPipelineTriggerPropertiesWithContext(context.Background(), createTektonPipelineTriggerPropertiesOptions)
}

// CreateTektonPipelineTriggerPropertiesWithContext is an alternate form of the CreateTektonPipelineTriggerProperties method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) CreateTektonPipelineTriggerPropertiesWithContext(ctx context.Context, createTektonPipelineTriggerPropertiesOptions *CreateTektonPipelineTriggerPropertiesOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTektonPipelineTriggerPropertiesOptions, "createTektonPipelineTriggerPropertiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createTektonPipelineTriggerPropertiesOptions, "createTektonPipelineTriggerPropertiesOptions")
	if err != nil {
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
	if createTektonPipelineTriggerPropertiesOptions.Value != nil {
		body["value"] = createTektonPipelineTriggerPropertiesOptions.Value
	}
	if createTektonPipelineTriggerPropertiesOptions.Enum != nil {
		body["enum"] = createTektonPipelineTriggerPropertiesOptions.Enum
	}
	if createTektonPipelineTriggerPropertiesOptions.Default != nil {
		body["default"] = createTektonPipelineTriggerPropertiesOptions.Default
	}
	if createTektonPipelineTriggerPropertiesOptions.Type != nil {
		body["type"] = createTektonPipelineTriggerPropertiesOptions.Type
	}
	if createTektonPipelineTriggerPropertiesOptions.Path != nil {
		body["path"] = createTektonPipelineTriggerPropertiesOptions.Path
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperty)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetTektonPipelineTriggerProperty : Get a trigger's environment property
// This request retrieves a trigger's environment property.
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTriggerProperty(getTektonPipelineTriggerPropertyOptions *GetTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.GetTektonPipelineTriggerPropertyWithContext(context.Background(), getTektonPipelineTriggerPropertyOptions)
}

// GetTektonPipelineTriggerPropertyWithContext is an alternate form of the GetTektonPipelineTriggerProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) GetTektonPipelineTriggerPropertyWithContext(ctx context.Context, getTektonPipelineTriggerPropertyOptions *GetTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTektonPipelineTriggerPropertyOptions, "getTektonPipelineTriggerPropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getTektonPipelineTriggerPropertyOptions, "getTektonPipelineTriggerPropertyOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperty)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// ReplaceTektonPipelineTriggerProperty : Replace a trigger's environment property value
// This request updates a trigger's environment property value, its type or name are immutable. For any property type,
// the entered value replaces the existing value.
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineTriggerProperty(replaceTektonPipelineTriggerPropertyOptions *ReplaceTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	return cdTektonPipeline.ReplaceTektonPipelineTriggerPropertyWithContext(context.Background(), replaceTektonPipelineTriggerPropertyOptions)
}

// ReplaceTektonPipelineTriggerPropertyWithContext is an alternate form of the ReplaceTektonPipelineTriggerProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) ReplaceTektonPipelineTriggerPropertyWithContext(ctx context.Context, replaceTektonPipelineTriggerPropertyOptions *ReplaceTektonPipelineTriggerPropertyOptions) (result *TriggerProperty, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceTektonPipelineTriggerPropertyOptions, "replaceTektonPipelineTriggerPropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(replaceTektonPipelineTriggerPropertyOptions, "replaceTektonPipelineTriggerPropertyOptions")
	if err != nil {
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
	if replaceTektonPipelineTriggerPropertyOptions.Value != nil {
		body["value"] = replaceTektonPipelineTriggerPropertyOptions.Value
	}
	if replaceTektonPipelineTriggerPropertyOptions.Enum != nil {
		body["enum"] = replaceTektonPipelineTriggerPropertyOptions.Enum
	}
	if replaceTektonPipelineTriggerPropertyOptions.Default != nil {
		body["default"] = replaceTektonPipelineTriggerPropertyOptions.Default
	}
	if replaceTektonPipelineTriggerPropertyOptions.Type != nil {
		body["type"] = replaceTektonPipelineTriggerPropertyOptions.Type
	}
	if replaceTektonPipelineTriggerPropertyOptions.Path != nil {
		body["path"] = replaceTektonPipelineTriggerPropertyOptions.Path
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
	response, err = cdTektonPipeline.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTriggerProperty)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTektonPipelineTriggerProperty : Delete a trigger's property
// this request deletes a trigger's property.
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTriggerProperty(deleteTektonPipelineTriggerPropertyOptions *DeleteTektonPipelineTriggerPropertyOptions) (response *core.DetailedResponse, err error) {
	return cdTektonPipeline.DeleteTektonPipelineTriggerPropertyWithContext(context.Background(), deleteTektonPipelineTriggerPropertyOptions)
}

// DeleteTektonPipelineTriggerPropertyWithContext is an alternate form of the DeleteTektonPipelineTriggerProperty method which supports a Context parameter
func (cdTektonPipeline *CdTektonPipelineV2) DeleteTektonPipelineTriggerPropertyWithContext(ctx context.Context, deleteTektonPipelineTriggerPropertyOptions *DeleteTektonPipelineTriggerPropertyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTektonPipelineTriggerPropertyOptions, "deleteTektonPipelineTriggerPropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteTektonPipelineTriggerPropertyOptions, "deleteTektonPipelineTriggerPropertyOptions")
	if err != nil {
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
		return
	}

	response, err = cdTektonPipeline.Service.Request(request, nil)

	return
}

// CancelTektonPipelineRunOptions : The CancelTektonPipelineRun options.
type CancelTektonPipelineRunOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Flag whether force cancel.
	Force *bool `json:"force,omitempty"`

	// Allows users to set headers on API requests
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

// Concurrency : Concurrency object.
type Concurrency struct {
	// Defines the maximum number of concurrent runs for this trigger.
	MaxConcurrentRuns *int64 `json:"max_concurrent_runs,omitempty"`
}

// UnmarshalConcurrency unmarshals an instance of Concurrency from the specified map of raw messages.
func UnmarshalConcurrency(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Concurrency)
	err = core.UnmarshalPrimitive(m, "max_concurrent_runs", &obj.MaxConcurrentRuns)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateTektonPipelineDefinitionOptions : The CreateTektonPipelineDefinition options.
type CreateTektonPipelineDefinitionOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Scm source for tekton pipeline defintion.
	ScmSource *DefinitionScmSource `json:"scm_source,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateTektonPipelineDefinitionOptions : Instantiate CreateTektonPipelineDefinitionOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineDefinitionOptions(pipelineID string) *CreateTektonPipelineDefinitionOptions {
	return &CreateTektonPipelineDefinitionOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelineDefinitionOptions) SetPipelineID(pipelineID string) *CreateTektonPipelineDefinitionOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetScmSource : Allow user to set ScmSource
func (_options *CreateTektonPipelineDefinitionOptions) SetScmSource(scmSource *DefinitionScmSource) *CreateTektonPipelineDefinitionOptions {
	_options.ScmSource = scmSource
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineDefinitionOptions) SetHeaders(param map[string]string) *CreateTektonPipelineDefinitionOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineOptions : The CreateTektonPipeline options.
type CreateTektonPipelineOptions struct {
	// UUID.
	ID *string `json:"id,omitempty"`

	// Worker object with worker ID only.
	Worker *WorkerWithID `json:"worker,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateTektonPipelineOptions : Instantiate CreateTektonPipelineOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineOptions() *CreateTektonPipelineOptions {
	return &CreateTektonPipelineOptions{}
}

// SetID : Allow user to set ID
func (_options *CreateTektonPipelineOptions) SetID(id string) *CreateTektonPipelineOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetWorker : Allow user to set Worker
func (_options *CreateTektonPipelineOptions) SetWorker(worker *WorkerWithID) *CreateTektonPipelineOptions {
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Property name.
	Name *string `json:"name,omitempty"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type,omitempty"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateTektonPipelinePropertiesOptions.Type property.
// Property type.
const (
	CreateTektonPipelinePropertiesOptionsTypeAppconfigConst = "APPCONFIG"
	CreateTektonPipelinePropertiesOptionsTypeIntegrationConst = "INTEGRATION"
	CreateTektonPipelinePropertiesOptionsTypeSecureConst = "SECURE"
	CreateTektonPipelinePropertiesOptionsTypeSingleSelectConst = "SINGLE_SELECT"
	CreateTektonPipelinePropertiesOptionsTypeTextConst = "TEXT"
)

// NewCreateTektonPipelinePropertiesOptions : Instantiate CreateTektonPipelinePropertiesOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelinePropertiesOptions(pipelineID string) *CreateTektonPipelinePropertiesOptions {
	return &CreateTektonPipelinePropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
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

// SetDefault : Allow user to set Default
func (_options *CreateTektonPipelinePropertiesOptions) SetDefault(defaultVar string) *CreateTektonPipelinePropertiesOptions {
	_options.Default = core.StringPtr(defaultVar)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateTektonPipelinePropertiesOptions) SetType(typeVar string) *CreateTektonPipelinePropertiesOptions {
	_options.Type = core.StringPtr(typeVar)
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Trigger name.
	TriggerName *string `json:"trigger_name,omitempty"`

	// A valid single dictionary object containing string values only to provide extra TEXT trigger properties or override
	// defined TEXT type trigger properties.
	TriggerProperties map[string]interface{} `json:"trigger_properties,omitempty"`

	// A valid single dictionary object containing string values only to provide extra SECURE type trigger properties or
	// override defined SECURE type trigger properties.
	SecureTriggerProperties map[string]interface{} `json:"secure_trigger_properties,omitempty"`

	// A valid single dictionary object with only string values to provide header, use $(header.header_key_name) to access
	// it in triggerBinding.
	TriggerHeader map[string]interface{} `json:"trigger_header,omitempty"`

	// A valid object to provide body, use $(body.body_key_name) to access it in triggerBinding.
	TriggerBody map[string]interface{} `json:"trigger_body,omitempty"`

	// Allows users to set headers on API requests
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

// SetTriggerHeader : Allow user to set TriggerHeader
func (_options *CreateTektonPipelineRunOptions) SetTriggerHeader(triggerHeader map[string]interface{}) *CreateTektonPipelineRunOptions {
	_options.TriggerHeader = triggerHeader
	return _options
}

// SetTriggerBody : Allow user to set TriggerBody
func (_options *CreateTektonPipelineRunOptions) SetTriggerBody(triggerBody map[string]interface{}) *CreateTektonPipelineRunOptions {
	_options.TriggerBody = triggerBody
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineRunOptions) SetHeaders(param map[string]string) *CreateTektonPipelineRunOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineTriggerOptions : The CreateTektonPipelineTrigger options.
type CreateTektonPipelineTriggerOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Tekton pipeline trigger object.
	Trigger TriggerIntf `json:"Trigger,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateTektonPipelineTriggerOptions : Instantiate CreateTektonPipelineTriggerOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineTriggerOptions(pipelineID string) *CreateTektonPipelineTriggerOptions {
	return &CreateTektonPipelineTriggerOptions{
		PipelineID: core.StringPtr(pipelineID),
	}
}

// SetPipelineID : Allow user to set PipelineID
func (_options *CreateTektonPipelineTriggerOptions) SetPipelineID(pipelineID string) *CreateTektonPipelineTriggerOptions {
	_options.PipelineID = core.StringPtr(pipelineID)
	return _options
}

// SetTrigger : Allow user to set Trigger
func (_options *CreateTektonPipelineTriggerOptions) SetTrigger(trigger TriggerIntf) *CreateTektonPipelineTriggerOptions {
	_options.Trigger = trigger
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *CreateTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// CreateTektonPipelineTriggerPropertiesOptions : The CreateTektonPipelineTriggerProperties options.
type CreateTektonPipelineTriggerPropertiesOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Property name.
	Name *string `json:"name,omitempty"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type,omitempty"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateTektonPipelineTriggerPropertiesOptions.Type property.
// Property type.
const (
	CreateTektonPipelineTriggerPropertiesOptionsTypeAppconfigConst = "APPCONFIG"
	CreateTektonPipelineTriggerPropertiesOptionsTypeIntegrationConst = "INTEGRATION"
	CreateTektonPipelineTriggerPropertiesOptionsTypeSecureConst = "SECURE"
	CreateTektonPipelineTriggerPropertiesOptionsTypeSingleSelectConst = "SINGLE_SELECT"
	CreateTektonPipelineTriggerPropertiesOptionsTypeTextConst = "TEXT"
)

// NewCreateTektonPipelineTriggerPropertiesOptions : Instantiate CreateTektonPipelineTriggerPropertiesOptions
func (*CdTektonPipelineV2) NewCreateTektonPipelineTriggerPropertiesOptions(pipelineID string, triggerID string) *CreateTektonPipelineTriggerPropertiesOptions {
	return &CreateTektonPipelineTriggerPropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
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

// SetDefault : Allow user to set Default
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetDefault(defaultVar string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Default = core.StringPtr(defaultVar)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetType(typeVar string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetPath : Allow user to set Path
func (_options *CreateTektonPipelineTriggerPropertiesOptions) SetPath(path string) *CreateTektonPipelineTriggerPropertiesOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTektonPipelineTriggerPropertiesOptions) SetHeaders(param map[string]string) *CreateTektonPipelineTriggerPropertiesOptions {
	options.Headers = param
	return options
}

// Definition : Tekton pipeline definition entry object.
type Definition struct {
	// Scm source for tekton pipeline defintion.
	ScmSource *DefinitionScmSource `json:"scm_source" validate:"required"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id" validate:"required"`

	// UUID.
	ID *string `json:"id" validate:"required"`
}

// NewDefinition : Instantiate Definition (Generic Model Constructor)
func (*CdTektonPipelineV2) NewDefinition(scmSource *DefinitionScmSource, serviceInstanceID string, id string) (_model *Definition, err error) {
	_model = &Definition{
		ScmSource: scmSource,
		ServiceInstanceID: core.StringPtr(serviceInstanceID),
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalDefinition unmarshals an instance of Definition from the specified map of raw messages.
func UnmarshalDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Definition)
	err = core.UnmarshalModel(m, "scm_source", &obj.ScmSource, UnmarshalDefinitionScmSource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
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

// DefinitionScmSource : Scm source for tekton pipeline defintion.
type DefinitionScmSource struct {
	// General href URL.
	URL *string `json:"url" validate:"required"`

	// A branch of the repo, branch field doesn't coexist with tag field.
	Branch *string `json:"branch,omitempty"`

	// A tag of the repo.
	Tag *string `json:"tag,omitempty"`

	// The path to the definitions yaml files.
	Path *string `json:"path" validate:"required"`
}

// NewDefinitionScmSource : Instantiate DefinitionScmSource (Generic Model Constructor)
func (*CdTektonPipelineV2) NewDefinitionScmSource(url string, path string) (_model *DefinitionScmSource, err error) {
	_model = &DefinitionScmSource{
		URL: core.StringPtr(url),
		Path: core.StringPtr(path),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalDefinitionScmSource unmarshals an instance of DefinitionScmSource from the specified map of raw messages.
func UnmarshalDefinitionScmSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefinitionScmSource)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tag", &obj.Tag)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Definitions : Pipeline definitions is a collection of individual definition entries, each entry is a composition of a repo url,
// repo branch/tag and a certain directory path.
type Definitions struct {
	// Definition list.
	Definitions []DefinitionsDefinitionsItem `json:"definitions" validate:"required"`
}

// UnmarshalDefinitions unmarshals an instance of Definitions from the specified map of raw messages.
func UnmarshalDefinitions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Definitions)
	err = core.UnmarshalModel(m, "definitions", &obj.Definitions, UnmarshalDefinitionsDefinitionsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DefinitionsDefinitionsItem : Tekton pipeline definition entry object.
type DefinitionsDefinitionsItem struct {
	// Scm source for tekton pipeline defintion.
	ScmSource *DefinitionScmSource `json:"scm_source" validate:"required"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id" validate:"required"`

	// UUID.
	ID *string `json:"id" validate:"required"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// UnmarshalDefinitionsDefinitionsItem unmarshals an instance of DefinitionsDefinitionsItem from the specified map of raw messages.
func UnmarshalDefinitionsDefinitionsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DefinitionsDefinitionsItem)
	err = core.UnmarshalModel(m, "scm_source", &obj.ScmSource, UnmarshalDefinitionScmSource)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteTektonPipelineDefinitionOptions : The DeleteTektonPipelineDefinition options.
type DeleteTektonPipelineDefinitionOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The definition ID.
	DefinitionID *string `json:"definition_id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The property's name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// The property's name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests
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

// EnvProperties : Pipeline properties object.
type EnvProperties struct {
	// Pipeline properties list.
	Properties []Property `json:"properties" validate:"required"`
}

// UnmarshalEnvProperties unmarshals an instance of EnvProperties from the specified map of raw messages.
func UnmarshalEnvProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvProperties)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Events : Needed only for git trigger type. Events object defines the events this git trigger listening to.
type Events struct {
	// If true, the trigger starts when tekton pipeline service receive a repo's 'push' git webhook event.
	Push *bool `json:"push,omitempty"`

	// If true, the trigger starts when tekton pipeline service receive a repo pull reqeust's 'close' git webhook event.
	PullRequestClosed *bool `json:"pull_request_closed,omitempty"`

	// If true, the trigger starts when tekton pipeline service receive a repo pull reqeust's 'open' or 'update' git
	// webhook event.
	PullRequest *bool `json:"pull_request,omitempty"`
}

// UnmarshalEvents unmarshals an instance of Events from the specified map of raw messages.
func UnmarshalEvents(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Events)
	err = core.UnmarshalPrimitive(m, "push", &obj.Push)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pull_request_closed", &obj.PullRequestClosed)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pull_request", &obj.PullRequest)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GenericSecret : Needed only for generic trigger type. Secret used to start generic trigger.
type GenericSecret struct {
	// Secret type.
	Type *string `json:"type,omitempty"`

	// Secret value, not needed if secret type is "internalValidation".
	Value *string `json:"value,omitempty"`

	// Secret location, not needed if secret type is "internalValidation".
	Source *string `json:"source,omitempty"`

	// Secret name, not needed if type is "internalValidation".
	KeyName *string `json:"key_name,omitempty"`

	// Algorithm used for "digestMatches" secret type.
	Algorithm *string `json:"algorithm,omitempty"`
}

// Constants associated with the GenericSecret.Type property.
// Secret type.
const (
	GenericSecretTypeDigestmatchesConst = "digestMatches"
	GenericSecretTypeInternalvalidationConst = "internalValidation"
	GenericSecretTypeTokenmatchesConst = "tokenMatches"
)

// Constants associated with the GenericSecret.Source property.
// Secret location, not needed if secret type is "internalValidation".
const (
	GenericSecretSourceHeaderConst = "header"
	GenericSecretSourcePayloadConst = "payload"
	GenericSecretSourceQueryConst = "query"
)

// Constants associated with the GenericSecret.Algorithm property.
// Algorithm used for "digestMatches" secret type.
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
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source", &obj.Source)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "key_name", &obj.KeyName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "algorithm", &obj.Algorithm)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetTektonPipelineDefinitionOptions : The GetTektonPipelineDefinition options.
type GetTektonPipelineDefinitionOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The definition ID.
	DefinitionID *string `json:"definition_id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The property's name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The tekton pipeline run ID.
	PipelineRunID *string `json:"pipeline_run_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Defines if response includes definition metadata.
	Includes *string `json:"includes,omitempty"`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// The property's name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// Filters the collection to resources with the specified property name.
	Name *string `json:"name,omitempty"`

	// Filters the collection to resources with the specified property type.
	Type []string `json:"type,omitempty"`

	// Sorts the returned properties by a property field, for properties you can sort them alphabetically by "name" in
	// ascending order or with "-name" in descending order.
	Sort *string `json:"sort,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListTektonPipelinePropertiesOptions.Type property.
// Query in URL.
const (
	ListTektonPipelinePropertiesOptionsTypeAppconfigConst = "APPCONFIG"
	ListTektonPipelinePropertiesOptionsTypeIntegrationConst = "INTEGRATION"
	ListTektonPipelinePropertiesOptionsTypeSecureConst = "SECURE"
	ListTektonPipelinePropertiesOptionsTypeSingleSelectConst = "SINGLE_SELECT"
	ListTektonPipelinePropertiesOptionsTypeTextConst = "TEXT"
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The number of pipeline runs to return, sorted by creation time, most recent first.
	Limit *int64 `json:"limit,omitempty"`

	// Skip the specified number of pipeline runs.
	Offset *int64 `json:"offset,omitempty"`

	// Filters the collection to resources with the specified status.
	Status *string `json:"status,omitempty"`

	// Filters the collection to resources with the specified trigger name.
	TriggerName *string `json:"trigger.name,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListTektonPipelineRunsOptions.Status property.
// Filters the collection to resources with the specified status.
const (
	ListTektonPipelineRunsOptionsStatusCancelledConst = "cancelled"
	ListTektonPipelineRunsOptionsStatusCancellingConst = "cancelling"
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

// SetLimit : Allow user to set Limit
func (_options *ListTektonPipelineRunsOptions) SetLimit(limit int64) *ListTektonPipelineRunsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListTektonPipelineRunsOptions) SetOffset(offset int64) *ListTektonPipelineRunsOptions {
	_options.Offset = core.Int64Ptr(offset)
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// filter properties by the name of property.
	Name *string `json:"name" validate:"required"`

	// filter properties by the type of property, avaialble types are
	// "SECURE","TEXT","INTEGRATION","SINGLE_SELECT","APPCONFIG".
	Type *string `json:"type" validate:"required"`

	// sort properties by the a property field, for properties you can only sort them by "name" or "-name".
	Sort *string `json:"sort" validate:"required"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListTektonPipelineTriggerPropertiesOptions : Instantiate ListTektonPipelineTriggerPropertiesOptions
func (*CdTektonPipelineV2) NewListTektonPipelineTriggerPropertiesOptions(pipelineID string, triggerID string, name string, typeVar string, sort string) *ListTektonPipelineTriggerPropertiesOptions {
	return &ListTektonPipelineTriggerPropertiesOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
		Sort: core.StringPtr(sort),
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// filter the triggers by trigger "type", possible values are "manual", "scm", "generic", "timer".
	Type *string `json:"type,omitempty"`

	// filter the triggers by trigger "name", accept single string value.
	Name *string `json:"name,omitempty"`

	// filter the triggers by trigger "event_listener", accept single string value.
	EventListener *string `json:"event_listener,omitempty"`

	// filter the triggers by trigger "worker.id", accept single string value.
	WorkerID *string `json:"worker.id,omitempty"`

	// filter the triggers by trigger "worker.name", accept single string value.
	WorkerName *string `json:"worker.name,omitempty"`

	// filter the triggers by trigger "disabled" flag, possbile state are "true" and "false".
	Disabled *string `json:"disabled,omitempty"`

	// filter the triggers by trigger "tags", the response lists triggers if it has one matching tag.
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests
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

// PipelineRun : Single tekton pipeline run object.
type PipelineRun struct {
	// UUID.
	ID *string `json:"id" validate:"required"`

	// User information.
	UserInfo *UserInfo `json:"user_info,omitempty"`

	// Status of the pipeline run.
	Status *string `json:"status" validate:"required"`

	// The aggregated definition ID used for this pipelineRun.
	DefinitionID *string `json:"definition_id" validate:"required"`

	// worker details used in this pipelineRun.
	Worker *PipelineRunWorker `json:"worker" validate:"required"`

	// UUID.
	PipelineID *string `json:"pipeline_id" validate:"required"`

	// Listener name used to start the run.
	ListenerName *string `json:"listener_name" validate:"required"`

	// Tekton pipeline trigger object.
	Trigger TriggerIntf `json:"trigger" validate:"required"`

	// Event parameters object passed to this pipeline run in String format, the contents depends on the type of trigger,
	// for example, for git trigger it includes the git event payload.
	EventParamsBlob *string `json:"event_params_blob" validate:"required"`

	// Heads passed to this pipeline run in String format, tekton pipeline service can't control the shape of the contents.
	EventHeaderParamsBlob *string `json:"event_header_params_blob,omitempty"`

	// Properties used in this tekton pipeline run.
	Properties []Property `json:"properties,omitempty"`

	// Standard RFC 3339 Date Time String.
	Created *strfmt.DateTime `json:"created" validate:"required"`

	// Standard RFC 3339 Date Time String.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Dashboard URL for this pipeline run.
	HTMLURL *string `json:"html_url" validate:"required"`
}

// Constants associated with the PipelineRun.Status property.
// Status of the pipeline run.
const (
	PipelineRunStatusCancelledConst = "cancelled"
	PipelineRunStatusCancellingConst = "cancelling"
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
		return
	}
	err = core.UnmarshalModel(m, "user_info", &obj.UserInfo, UnmarshalUserInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "definition_id", &obj.DefinitionID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalPipelineRunWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pipeline_id", &obj.PipelineID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "listener_name", &obj.ListenerName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "trigger", &obj.Trigger, UnmarshalTrigger)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_params_blob", &obj.EventParamsBlob)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_header_params_blob", &obj.EventHeaderParamsBlob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
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
	err = core.UnmarshalPrimitive(m, "html_url", &obj.HTMLURL)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRunLog : Pipeline run logId object.
type PipelineRunLog struct {
	// <podName>/<containerName> of this log.
	Name *string `json:"name,omitempty"`

	// Generated log ID.
	ID *string `json:"id,omitempty"`

	// API for getting log content.
	Href *string `json:"href,omitempty"`
}

// UnmarshalPipelineRunLog unmarshals an instance of PipelineRunLog from the specified map of raw messages.
func UnmarshalPipelineRunLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunLog)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRunLogs : List of pipeline run log ID object.
type PipelineRunLogs struct {
	Logs []PipelineRunLog `json:"logs,omitempty"`
}

// UnmarshalPipelineRunLogs unmarshals an instance of PipelineRunLogs from the specified map of raw messages.
func UnmarshalPipelineRunLogs(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunLogs)
	err = core.UnmarshalModel(m, "logs", &obj.Logs, UnmarshalPipelineRunLog)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRunWorker : worker details used in this pipelineRun.
type PipelineRunWorker struct {
	// Worker name.
	Name *string `json:"name,omitempty"`

	// The agent ID of the corresponding private worker integration used for this pipelineRun.
	Agent *string `json:"agent,omitempty"`

	// The Service ID of the corresponding private worker integration used for this pipelineRun.
	ServiceID *string `json:"service_id,omitempty"`

	// UUID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalPipelineRunWorker unmarshals an instance of PipelineRunWorker from the specified map of raw messages.
func UnmarshalPipelineRunWorker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunWorker)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "agent", &obj.Agent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_id", &obj.ServiceID)
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

// PipelineRuns : Tekton pipeline runs object.
type PipelineRuns struct {
	// Tekton pipeline runs list.
	PipelineRuns []PipelineRunsPipelineRunsItem `json:"pipeline_runs" validate:"required"`

	// Skip a specified number of pipeline runs.
	Offset *int64 `json:"offset" validate:"required"`

	// The number of pipeline runs to return, sorted by creation time, most recent first.
	Limit *int64 `json:"limit" validate:"required"`

	// First page of pipeline runs.
	First *PipelineRunsFirst `json:"first" validate:"required"`

	// Next page of pipeline runs relative to the offset and limit.
	Next *PipelineRunsNext `json:"next,omitempty"`
}

// UnmarshalPipelineRuns unmarshals an instance of PipelineRuns from the specified map of raw messages.
func UnmarshalPipelineRuns(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRuns)
	err = core.UnmarshalModel(m, "pipeline_runs", &obj.PipelineRuns, UnmarshalPipelineRunsPipelineRunsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPipelineRunsFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPipelineRunsNext)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PipelineRuns) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil || offset == nil {
		return nil, err
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// PipelineRunsFirst : First page of pipeline runs.
type PipelineRunsFirst struct {
	// General href URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPipelineRunsFirst unmarshals an instance of PipelineRunsFirst from the specified map of raw messages.
func UnmarshalPipelineRunsFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunsFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRunsNext : Next page of pipeline runs relative to the offset and limit.
type PipelineRunsNext struct {
	// General href URL.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPipelineRunsNext unmarshals an instance of PipelineRunsNext from the specified map of raw messages.
func UnmarshalPipelineRunsNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunsNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PipelineRunsPipelineRunsItem : Single tekton pipeline run object.
type PipelineRunsPipelineRunsItem struct {
	// UUID.
	ID *string `json:"id" validate:"required"`

	// User information.
	UserInfo *UserInfo `json:"user_info,omitempty"`

	// Status of the pipeline run.
	Status *string `json:"status" validate:"required"`

	// The aggregated definition ID used for this pipelineRun.
	DefinitionID *string `json:"definition_id" validate:"required"`

	// worker details used in this pipelineRun.
	Worker *PipelineRunWorker `json:"worker" validate:"required"`

	// UUID.
	PipelineID *string `json:"pipeline_id" validate:"required"`

	// Listener name used to start the run.
	ListenerName *string `json:"listener_name" validate:"required"`

	// Tekton pipeline trigger object.
	Trigger TriggerIntf `json:"trigger" validate:"required"`

	// Event parameters object passed to this pipeline run in String format, the contents depends on the type of trigger,
	// for example, for git trigger it includes the git event payload.
	EventParamsBlob *string `json:"event_params_blob" validate:"required"`

	// Heads passed to this pipeline run in String format, tekton pipeline service can't control the shape of the contents.
	EventHeaderParamsBlob *string `json:"event_header_params_blob,omitempty"`

	// Properties used in this tekton pipeline run.
	Properties []Property `json:"properties,omitempty"`

	// Standard RFC 3339 Date Time String.
	Created *strfmt.DateTime `json:"created" validate:"required"`

	// Standard RFC 3339 Date Time String.
	Updated *strfmt.DateTime `json:"updated,omitempty"`

	// Dashboard URL for this pipeline run.
	HTMLURL *string `json:"html_url" validate:"required"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the PipelineRunsPipelineRunsItem.Status property.
// Status of the pipeline run.
const (
	PipelineRunsPipelineRunsItemStatusCancelledConst = "cancelled"
	PipelineRunsPipelineRunsItemStatusCancellingConst = "cancelling"
	PipelineRunsPipelineRunsItemStatusErrorConst = "error"
	PipelineRunsPipelineRunsItemStatusFailedConst = "failed"
	PipelineRunsPipelineRunsItemStatusPendingConst = "pending"
	PipelineRunsPipelineRunsItemStatusQueuedConst = "queued"
	PipelineRunsPipelineRunsItemStatusRunningConst = "running"
	PipelineRunsPipelineRunsItemStatusSucceededConst = "succeeded"
	PipelineRunsPipelineRunsItemStatusWaitingConst = "waiting"
)

// UnmarshalPipelineRunsPipelineRunsItem unmarshals an instance of PipelineRunsPipelineRunsItem from the specified map of raw messages.
func UnmarshalPipelineRunsPipelineRunsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PipelineRunsPipelineRunsItem)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "user_info", &obj.UserInfo, UnmarshalUserInfo)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "definition_id", &obj.DefinitionID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalPipelineRunWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pipeline_id", &obj.PipelineID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "listener_name", &obj.ListenerName)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "trigger", &obj.Trigger, UnmarshalTrigger)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_params_blob", &obj.EventParamsBlob)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_header_params_blob", &obj.EventHeaderParamsBlob)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
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
	err = core.UnmarshalPrimitive(m, "html_url", &obj.HTMLURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Property : Property object.
type Property struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`
}

// Constants associated with the Property.Type property.
// Property type.
const (
	PropertyTypeAppconfigConst = "APPCONFIG"
	PropertyTypeIntegrationConst = "INTEGRATION"
	PropertyTypeSecureConst = "SECURE"
	PropertyTypeSingleSelectConst = "SINGLE_SELECT"
	PropertyTypeTextConst = "TEXT"
)

// NewProperty : Instantiate Property (Generic Model Constructor)
func (*CdTektonPipelineV2) NewProperty(name string, typeVar string) (_model *Property, err error) {
	_model = &Property{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalProperty unmarshals an instance of Property from the specified map of raw messages.
func UnmarshalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Property)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceTektonPipelineDefinitionOptions : The ReplaceTektonPipelineDefinition options.
type ReplaceTektonPipelineDefinitionOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The definition ID.
	DefinitionID *string `json:"definition_id" validate:"required,ne="`

	// Scm source for tekton pipeline defintion.
	ScmSource *DefinitionScmSource `json:"scm_source,omitempty"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewReplaceTektonPipelineDefinitionOptions : Instantiate ReplaceTektonPipelineDefinitionOptions
func (*CdTektonPipelineV2) NewReplaceTektonPipelineDefinitionOptions(pipelineID string, definitionID string) *ReplaceTektonPipelineDefinitionOptions {
	return &ReplaceTektonPipelineDefinitionOptions{
		PipelineID: core.StringPtr(pipelineID),
		DefinitionID: core.StringPtr(definitionID),
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

// SetScmSource : Allow user to set ScmSource
func (_options *ReplaceTektonPipelineDefinitionOptions) SetScmSource(scmSource *DefinitionScmSource) *ReplaceTektonPipelineDefinitionOptions {
	_options.ScmSource = scmSource
	return _options
}

// SetServiceInstanceID : Allow user to set ServiceInstanceID
func (_options *ReplaceTektonPipelineDefinitionOptions) SetServiceInstanceID(serviceInstanceID string) *ReplaceTektonPipelineDefinitionOptions {
	_options.ServiceInstanceID = core.StringPtr(serviceInstanceID)
	return _options
}

// SetID : Allow user to set ID
func (_options *ReplaceTektonPipelineDefinitionOptions) SetID(id string) *ReplaceTektonPipelineDefinitionOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTektonPipelineDefinitionOptions) SetHeaders(param map[string]string) *ReplaceTektonPipelineDefinitionOptions {
	options.Headers = param
	return options
}

// ReplaceTektonPipelinePropertyOptions : The ReplaceTektonPipelineProperty options.
type ReplaceTektonPipelinePropertyOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The property's name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Property name.
	Name *string `json:"name,omitempty"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type,omitempty"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceTektonPipelinePropertyOptions.Type property.
// Property type.
const (
	ReplaceTektonPipelinePropertyOptionsTypeAppconfigConst = "APPCONFIG"
	ReplaceTektonPipelinePropertyOptionsTypeIntegrationConst = "INTEGRATION"
	ReplaceTektonPipelinePropertyOptionsTypeSecureConst = "SECURE"
	ReplaceTektonPipelinePropertyOptionsTypeSingleSelectConst = "SINGLE_SELECT"
	ReplaceTektonPipelinePropertyOptionsTypeTextConst = "TEXT"
)

// NewReplaceTektonPipelinePropertyOptions : Instantiate ReplaceTektonPipelinePropertyOptions
func (*CdTektonPipelineV2) NewReplaceTektonPipelinePropertyOptions(pipelineID string, propertyName string) *ReplaceTektonPipelinePropertyOptions {
	return &ReplaceTektonPipelinePropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		PropertyName: core.StringPtr(propertyName),
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

// SetDefault : Allow user to set Default
func (_options *ReplaceTektonPipelinePropertyOptions) SetDefault(defaultVar string) *ReplaceTektonPipelinePropertyOptions {
	_options.Default = core.StringPtr(defaultVar)
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceTektonPipelinePropertyOptions) SetType(typeVar string) *ReplaceTektonPipelinePropertyOptions {
	_options.Type = core.StringPtr(typeVar)
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
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// The property's name.
	PropertyName *string `json:"property_name" validate:"required,ne="`

	// Property name.
	Name *string `json:"name,omitempty"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type,omitempty"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ReplaceTektonPipelineTriggerPropertyOptions.Type property.
// Property type.
const (
	ReplaceTektonPipelineTriggerPropertyOptionsTypeAppconfigConst = "APPCONFIG"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeIntegrationConst = "INTEGRATION"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeSecureConst = "SECURE"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeSingleSelectConst = "SINGLE_SELECT"
	ReplaceTektonPipelineTriggerPropertyOptionsTypeTextConst = "TEXT"
)

// NewReplaceTektonPipelineTriggerPropertyOptions : Instantiate ReplaceTektonPipelineTriggerPropertyOptions
func (*CdTektonPipelineV2) NewReplaceTektonPipelineTriggerPropertyOptions(pipelineID string, triggerID string, propertyName string) *ReplaceTektonPipelineTriggerPropertyOptions {
	return &ReplaceTektonPipelineTriggerPropertyOptions{
		PipelineID: core.StringPtr(pipelineID),
		TriggerID: core.StringPtr(triggerID),
		PropertyName: core.StringPtr(propertyName),
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

// SetDefault : Allow user to set Default
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetDefault(defaultVar string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Default = core.StringPtr(defaultVar)
	return _options
}

// SetType : Allow user to set Type
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetType(typeVar string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetPath : Allow user to set Path
func (_options *ReplaceTektonPipelineTriggerPropertyOptions) SetPath(path string) *ReplaceTektonPipelineTriggerPropertyOptions {
	_options.Path = core.StringPtr(path)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceTektonPipelineTriggerPropertyOptions) SetHeaders(param map[string]string) *ReplaceTektonPipelineTriggerPropertyOptions {
	options.Headers = param
	return options
}

// RerunTektonPipelineRunOptions : The RerunTektonPipelineRun options.
type RerunTektonPipelineRunOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests
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

// StepLog : Log object for tekton pipeline run step.
type StepLog struct {
	// Step log ID.
	ID *string `json:"id" validate:"required"`

	// The raw log content of step.
	Data *string `json:"data" validate:"required"`
}

// UnmarshalStepLog unmarshals an instance of StepLog from the specified map of raw messages.
func UnmarshalStepLog(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(StepLog)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
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

	// ID.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain object.
	Toolchain *Toolchain `json:"toolchain" validate:"required"`

	// UUID.
	ID *string `json:"id" validate:"required"`

	// Definition list.
	Definitions []Definition `json:"definitions" validate:"required"`

	// Tekton pipeline's environment properties.
	Properties []Property `json:"properties" validate:"required"`

	// Standard RFC 3339 Date Time String.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Standard RFC 3339 Date Time String.
	Created *strfmt.DateTime `json:"created" validate:"required"`

	// Tekton pipeline definition document detail object. If this property is absent, the pipeline has no definitions
	// added.
	PipelineDefinition *TektonPipelinePipelineDefinition `json:"pipeline_definition,omitempty"`

	// Tekton pipeline triggers list.
	Triggers []TriggerIntf `json:"triggers" validate:"required"`

	// Default pipeline worker used to run the pipeline.
	Worker *Worker `json:"worker" validate:"required"`

	// Dashboard URL of this pipeline.
	HTMLURL *string `json:"html_url" validate:"required"`

	// The latest pipeline run build number. If this property is absent, the pipeline hasn't had any pipelineRuns.
	BuildNumber *int64 `json:"build_number,omitempty"`

	// Flag whether this pipeline enabled.
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
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "toolchain", &obj.Toolchain, UnmarshalToolchain)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "definitions", &obj.Definitions, UnmarshalDefinition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created", &obj.Created)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "pipeline_definition", &obj.PipelineDefinition, UnmarshalTektonPipelinePipelineDefinition)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "triggers", &obj.Triggers, UnmarshalTrigger)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "html_url", &obj.HTMLURL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "build_number", &obj.BuildNumber)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TektonPipelinePipelineDefinition : Tekton pipeline definition document detail object. If this property is absent, the pipeline has no definitions added.
type TektonPipelinePipelineDefinition struct {
	// The state of pipeline definition status.
	Status *string `json:"status,omitempty"`

	// UUID.
	ID *string `json:"id,omitempty"`
}

// Constants associated with the TektonPipelinePipelineDefinition.Status property.
// The state of pipeline definition status.
const (
	TektonPipelinePipelineDefinitionStatusFailedConst = "failed"
	TektonPipelinePipelineDefinitionStatusOutdatedConst = "outdated"
	TektonPipelinePipelineDefinitionStatusUpdatedConst = "updated"
	TektonPipelinePipelineDefinitionStatusUpdatingConst = "updating"
)

// UnmarshalTektonPipelinePipelineDefinition unmarshals an instance of TektonPipelinePipelineDefinition from the specified map of raw messages.
func UnmarshalTektonPipelinePipelineDefinition(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TektonPipelinePipelineDefinition)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
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

// Toolchain : Toolchain object.
type Toolchain struct {
	// UUID.
	ID *string `json:"id" validate:"required"`

	// The CRN for the toolchain that containing the tekton pipeline.
	CRN *string `json:"crn" validate:"required"`
}

// UnmarshalToolchain unmarshals an instance of Toolchain from the specified map of raw messages.
func UnmarshalToolchain(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Toolchain)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Trigger : Tekton pipeline trigger object.
// Models which "extend" this model:
// - TriggerDuplicateTrigger
// - TriggerManualTrigger
// - TriggerScmTrigger
// - TriggerTimerTrigger
// - TriggerGenericTrigger
type Trigger struct {
	// source trigger ID to clone from.
	SourceTriggerID *string `json:"source_trigger_id,omitempty"`

	// name of the duplicated trigger.
	Name *string `json:"name,omitempty"`

	// Trigger type.
	Type *string `json:"type,omitempty"`

	// Event listener name.
	EventListener *string `json:"event_listener,omitempty"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// Scm source for git type tekton pipeline trigger.
	ScmSource *TriggerScmSource `json:"scm_source,omitempty"`

	// Needed only for git trigger type. Events object defines the events this git trigger listening to.
	Events *Events `json:"events,omitempty"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`

	// Needed only for timer trigger type. Cron expression for timer trigger. Maximum frequency is every 5 minutes.
	Cron *string `json:"cron,omitempty"`

	// Needed only for timer trigger type. Timezones for timer trigger.
	Timezone *string `json:"timezone,omitempty"`

	// Needed only for generic trigger type. Secret used to start generic trigger.
	Secret *GenericSecret `json:"secret,omitempty"`
}
func (*Trigger) isaTrigger() bool {
	return true
}

type TriggerIntf interface {
	isaTrigger() bool
}

// UnmarshalTrigger unmarshals an instance of Trigger from the specified map of raw messages.
func UnmarshalTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Trigger)
	err = core.UnmarshalPrimitive(m, "source_trigger_id", &obj.SourceTriggerID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scm_source", &obj.ScmSource, UnmarshalTriggerScmSource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "events", &obj.Events, UnmarshalEvents)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerGenericTriggerPropertiesItem : Trigger Property object.
type TriggerGenericTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggerGenericTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggerGenericTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggerGenericTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggerGenericTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggerGenericTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerGenericTriggerPropertiesItemTypeTextConst = "TEXT"
)

// NewTriggerGenericTriggerPropertiesItem : Instantiate TriggerGenericTriggerPropertiesItem (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerGenericTriggerPropertiesItem(name string, typeVar string) (_model *TriggerGenericTriggerPropertiesItem, err error) {
	_model = &TriggerGenericTriggerPropertiesItem{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTriggerGenericTriggerPropertiesItem unmarshals an instance of TriggerGenericTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggerGenericTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerGenericTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerManualTriggerPropertiesItem : Trigger Property object.
type TriggerManualTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggerManualTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggerManualTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggerManualTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggerManualTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggerManualTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerManualTriggerPropertiesItemTypeTextConst = "TEXT"
)

// NewTriggerManualTriggerPropertiesItem : Instantiate TriggerManualTriggerPropertiesItem (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerManualTriggerPropertiesItem(name string, typeVar string) (_model *TriggerManualTriggerPropertiesItem, err error) {
	_model = &TriggerManualTriggerPropertiesItem{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTriggerManualTriggerPropertiesItem unmarshals an instance of TriggerManualTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggerManualTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerManualTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerProperties : Trigger properties object.
type TriggerProperties struct {
	// Trigger properties list.
	Properties []TriggerPropertiesPropertiesItem `json:"properties" validate:"required"`
}

// UnmarshalTriggerProperties unmarshals an instance of TriggerProperties from the specified map of raw messages.
func UnmarshalTriggerProperties(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerProperties)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerPropertiesPropertiesItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerPropertiesItem : Trigger Property object.
type TriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggerPropertiesItem.Type property.
// Property type.
const (
	TriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerPropertiesItemTypeTextConst = "TEXT"
)

// NewTriggerPropertiesItem : Instantiate TriggerPropertiesItem (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerPropertiesItem(name string, typeVar string) (_model *TriggerPropertiesItem, err error) {
	_model = &TriggerPropertiesItem{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTriggerPropertiesItem unmarshals an instance of TriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerPropertiesPropertiesItem : Trigger Property object.
type TriggerPropertiesPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggerPropertiesPropertiesItem.Type property.
// Property type.
const (
	TriggerPropertiesPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggerPropertiesPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggerPropertiesPropertiesItemTypeSecureConst = "SECURE"
	TriggerPropertiesPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerPropertiesPropertiesItemTypeTextConst = "TEXT"
)

// UnmarshalTriggerPropertiesPropertiesItem unmarshals an instance of TriggerPropertiesPropertiesItem from the specified map of raw messages.
func UnmarshalTriggerPropertiesPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerPropertiesPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerProperty : Trigger Property object.
type TriggerProperty struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`
}

// Constants associated with the TriggerProperty.Type property.
// Property type.
const (
	TriggerPropertyTypeAppconfigConst = "APPCONFIG"
	TriggerPropertyTypeIntegrationConst = "INTEGRATION"
	TriggerPropertyTypeSecureConst = "SECURE"
	TriggerPropertyTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerPropertyTypeTextConst = "TEXT"
)

// NewTriggerProperty : Instantiate TriggerProperty (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerProperty(name string, typeVar string) (_model *TriggerProperty, err error) {
	_model = &TriggerProperty{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTriggerProperty unmarshals an instance of TriggerProperty from the specified map of raw messages.
func UnmarshalTriggerProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerProperty)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerScmSource : Scm source for git type tekton pipeline trigger.
type TriggerScmSource struct {
	// Needed only for git trigger type. Repo URL that listening to.
	URL *string `json:"url,omitempty"`

	// Needed only for git trigger type. Branch name of the repo. Branch field doesn't coexist with pattern field.
	Branch *string `json:"branch,omitempty"`

	// Needed only for git trigger type. Git branch or tag pattern to listen to. Please refer to
	// https://github.com/micromatch/micromatch for pattern syntax.
	Pattern *string `json:"pattern,omitempty"`

	// Needed only for git trigger type. Branch name of the repo.
	BlindConnection *bool `json:"blind_connection,omitempty"`

	// Webhook ID.
	HookID *string `json:"hook_id,omitempty"`
}

// UnmarshalTriggerScmSource unmarshals an instance of TriggerScmSource from the specified map of raw messages.
func UnmarshalTriggerScmSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerScmSource)
	err = core.UnmarshalPrimitive(m, "url", &obj.URL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "branch", &obj.Branch)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pattern", &obj.Pattern)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "blind_connection", &obj.BlindConnection)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "hook_id", &obj.HookID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerScmTriggerPropertiesItem : Trigger Property object.
type TriggerScmTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggerScmTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggerScmTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggerScmTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggerScmTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggerScmTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerScmTriggerPropertiesItemTypeTextConst = "TEXT"
)

// NewTriggerScmTriggerPropertiesItem : Instantiate TriggerScmTriggerPropertiesItem (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerScmTriggerPropertiesItem(name string, typeVar string) (_model *TriggerScmTriggerPropertiesItem, err error) {
	_model = &TriggerScmTriggerPropertiesItem{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTriggerScmTriggerPropertiesItem unmarshals an instance of TriggerScmTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggerScmTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerScmTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerTimerTriggerPropertiesItem : Trigger Property object.
type TriggerTimerTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggerTimerTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggerTimerTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggerTimerTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggerTimerTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggerTimerTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggerTimerTriggerPropertiesItemTypeTextConst = "TEXT"
)

// NewTriggerTimerTriggerPropertiesItem : Instantiate TriggerTimerTriggerPropertiesItem (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerTimerTriggerPropertiesItem(name string, typeVar string) (_model *TriggerTimerTriggerPropertiesItem, err error) {
	_model = &TriggerTimerTriggerPropertiesItem{
		Name: core.StringPtr(name),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalTriggerTimerTriggerPropertiesItem unmarshals an instance of TriggerTimerTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggerTimerTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerTimerTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Triggers : Tekton pipeline triggers object.
type Triggers struct {
	// Tekton pipeline triggers list.
	Triggers []TriggersTriggersItemIntf `json:"triggers" validate:"required"`
}

// UnmarshalTriggers unmarshals an instance of Triggers from the specified map of raw messages.
func UnmarshalTriggers(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Triggers)
	err = core.UnmarshalModel(m, "triggers", &obj.Triggers, UnmarshalTriggersTriggersItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItem : Tekton pipeline trigger object.
// Models which "extend" this model:
// - TriggersTriggersItemTriggerDuplicateTrigger
// - TriggersTriggersItemTriggerManualTrigger
// - TriggersTriggersItemTriggerScmTrigger
// - TriggersTriggersItemTriggerTimerTrigger
// - TriggersTriggersItemTriggerGenericTrigger
type TriggersTriggersItem struct {
	// source trigger ID to clone from.
	SourceTriggerID *string `json:"source_trigger_id,omitempty"`

	// name of the duplicated trigger.
	Name *string `json:"name,omitempty"`

	// Trigger type.
	Type *string `json:"type,omitempty"`

	// Event listener name.
	EventListener *string `json:"event_listener,omitempty"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// Scm source for git type tekton pipeline trigger.
	ScmSource *TriggerScmSource `json:"scm_source,omitempty"`

	// Needed only for git trigger type. Events object defines the events this git trigger listening to.
	Events *Events `json:"events,omitempty"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`

	// Needed only for timer trigger type. Cron expression for timer trigger. Maximum frequency is every 5 minutes.
	Cron *string `json:"cron,omitempty"`

	// Needed only for timer trigger type. Timezones for timer trigger.
	Timezone *string `json:"timezone,omitempty"`

	// Needed only for generic trigger type. Secret used to start generic trigger.
	Secret *GenericSecret `json:"secret,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}
func (*TriggersTriggersItem) isaTriggersTriggersItem() bool {
	return true
}

type TriggersTriggersItemIntf interface {
	isaTriggersTriggersItem() bool
}

// UnmarshalTriggersTriggersItem unmarshals an instance of TriggersTriggersItem from the specified map of raw messages.
func UnmarshalTriggersTriggersItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItem)
	err = core.UnmarshalPrimitive(m, "source_trigger_id", &obj.SourceTriggerID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scm_source", &obj.ScmSource, UnmarshalTriggerScmSource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "events", &obj.Events, UnmarshalEvents)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerGenericTriggerPropertiesItem : Trigger Property object.
type TriggersTriggersItemTriggerGenericTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggersTriggersItemTriggerGenericTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggersTriggersItemTriggerGenericTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggersTriggersItemTriggerGenericTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggersTriggersItemTriggerGenericTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggersTriggersItemTriggerGenericTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggersTriggersItemTriggerGenericTriggerPropertiesItemTypeTextConst = "TEXT"
)

// UnmarshalTriggersTriggersItemTriggerGenericTriggerPropertiesItem unmarshals an instance of TriggersTriggersItemTriggerGenericTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerGenericTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerGenericTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerManualTriggerPropertiesItem : Trigger Property object.
type TriggersTriggersItemTriggerManualTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggersTriggersItemTriggerManualTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggersTriggersItemTriggerManualTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggersTriggersItemTriggerManualTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggersTriggersItemTriggerManualTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggersTriggersItemTriggerManualTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggersTriggersItemTriggerManualTriggerPropertiesItemTypeTextConst = "TEXT"
)

// UnmarshalTriggersTriggersItemTriggerManualTriggerPropertiesItem unmarshals an instance of TriggersTriggersItemTriggerManualTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerManualTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerManualTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerScmTriggerPropertiesItem : Trigger Property object.
type TriggersTriggersItemTriggerScmTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggersTriggersItemTriggerScmTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggersTriggersItemTriggerScmTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggersTriggersItemTriggerScmTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggersTriggersItemTriggerScmTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggersTriggersItemTriggerScmTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggersTriggersItemTriggerScmTriggerPropertiesItemTypeTextConst = "TEXT"
)

// UnmarshalTriggersTriggersItemTriggerScmTriggerPropertiesItem unmarshals an instance of TriggersTriggersItemTriggerScmTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerScmTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerScmTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerTimerTriggerPropertiesItem : Trigger Property object.
type TriggersTriggersItemTriggerTimerTriggerPropertiesItem struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// String format property value.
	Value *string `json:"value,omitempty"`

	// Options for SINGLE_SELECT property type.
	Enum []string `json:"enum,omitempty"`

	// Default option for SINGLE_SELECT property type.
	Default *string `json:"default,omitempty"`

	// Property type.
	Type *string `json:"type" validate:"required"`

	// property path for INTEGRATION type properties.
	Path *string `json:"path,omitempty"`

	// General href URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the TriggersTriggersItemTriggerTimerTriggerPropertiesItem.Type property.
// Property type.
const (
	TriggersTriggersItemTriggerTimerTriggerPropertiesItemTypeAppconfigConst = "APPCONFIG"
	TriggersTriggersItemTriggerTimerTriggerPropertiesItemTypeIntegrationConst = "INTEGRATION"
	TriggersTriggersItemTriggerTimerTriggerPropertiesItemTypeSecureConst = "SECURE"
	TriggersTriggersItemTriggerTimerTriggerPropertiesItemTypeSingleSelectConst = "SINGLE_SELECT"
	TriggersTriggersItemTriggerTimerTriggerPropertiesItemTypeTextConst = "TEXT"
)

// UnmarshalTriggersTriggersItemTriggerTimerTriggerPropertiesItem unmarshals an instance of TriggersTriggersItemTriggerTimerTriggerPropertiesItem from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerTimerTriggerPropertiesItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerTimerTriggerPropertiesItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enum", &obj.Enum)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateTektonPipelineOptions : The UpdateTektonPipeline options.
type UpdateTektonPipelineOptions struct {
	// ID of current instance.
	ID *string `json:"id" validate:"required,ne="`

	// Worker object with worker ID only.
	Worker *WorkerWithID `json:"worker,omitempty"`

	// Allows users to set headers on API requests
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

// SetWorker : Allow user to set Worker
func (_options *UpdateTektonPipelineOptions) SetWorker(worker *WorkerWithID) *UpdateTektonPipelineOptions {
	_options.Worker = worker
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTektonPipelineOptions) SetHeaders(param map[string]string) *UpdateTektonPipelineOptions {
	options.Headers = param
	return options
}

// UpdateTektonPipelineTriggerOptions : The UpdateTektonPipelineTrigger options.
type UpdateTektonPipelineTriggerOptions struct {
	// The tekton pipeline ID.
	PipelineID *string `json:"pipeline_id" validate:"required,ne="`

	// The trigger ID.
	TriggerID *string `json:"trigger_id" validate:"required,ne="`

	// Trigger type.
	Type *string `json:"type,omitempty"`

	// Trigger name.
	Name *string `json:"name,omitempty"`

	// Event listener name.
	EventListener *string `json:"event_listener,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// Defines if this trigger is disabled.
	Disabled *bool `json:"disabled,omitempty"`

	// Needed only for generic trigger type. Secret used to start generic trigger.
	Secret *GenericSecret `json:"secret,omitempty"`

	// Needed only for timer trigger type. Cron expression for timer trigger.
	Cron *string `json:"cron,omitempty"`

	// Needed only for timer trigger type. Timezones for timer trigger.
	Timezone *string `json:"timezone,omitempty"`

	// Scm source for git type tekton pipeline trigger.
	ScmSource *TriggerScmSource `json:"scm_source,omitempty"`

	// Needed only for git trigger type. Events object defines the events this git trigger listening to.
	Events *Events `json:"events,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateTektonPipelineTriggerOptions.Type property.
// Trigger type.
const (
	UpdateTektonPipelineTriggerOptionsTypeGenericConst = "generic"
	UpdateTektonPipelineTriggerOptionsTypeManualConst = "manual"
	UpdateTektonPipelineTriggerOptionsTypeScmConst = "scm"
	UpdateTektonPipelineTriggerOptionsTypeTimerConst = "timer"
)

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

// SetType : Allow user to set Type
func (_options *UpdateTektonPipelineTriggerOptions) SetType(typeVar string) *UpdateTektonPipelineTriggerOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateTektonPipelineTriggerOptions) SetName(name string) *UpdateTektonPipelineTriggerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetEventListener : Allow user to set EventListener
func (_options *UpdateTektonPipelineTriggerOptions) SetEventListener(eventListener string) *UpdateTektonPipelineTriggerOptions {
	_options.EventListener = core.StringPtr(eventListener)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateTektonPipelineTriggerOptions) SetTags(tags []string) *UpdateTektonPipelineTriggerOptions {
	_options.Tags = tags
	return _options
}

// SetWorker : Allow user to set Worker
func (_options *UpdateTektonPipelineTriggerOptions) SetWorker(worker *Worker) *UpdateTektonPipelineTriggerOptions {
	_options.Worker = worker
	return _options
}

// SetConcurrency : Allow user to set Concurrency
func (_options *UpdateTektonPipelineTriggerOptions) SetConcurrency(concurrency *Concurrency) *UpdateTektonPipelineTriggerOptions {
	_options.Concurrency = concurrency
	return _options
}

// SetDisabled : Allow user to set Disabled
func (_options *UpdateTektonPipelineTriggerOptions) SetDisabled(disabled bool) *UpdateTektonPipelineTriggerOptions {
	_options.Disabled = core.BoolPtr(disabled)
	return _options
}

// SetSecret : Allow user to set Secret
func (_options *UpdateTektonPipelineTriggerOptions) SetSecret(secret *GenericSecret) *UpdateTektonPipelineTriggerOptions {
	_options.Secret = secret
	return _options
}

// SetCron : Allow user to set Cron
func (_options *UpdateTektonPipelineTriggerOptions) SetCron(cron string) *UpdateTektonPipelineTriggerOptions {
	_options.Cron = core.StringPtr(cron)
	return _options
}

// SetTimezone : Allow user to set Timezone
func (_options *UpdateTektonPipelineTriggerOptions) SetTimezone(timezone string) *UpdateTektonPipelineTriggerOptions {
	_options.Timezone = core.StringPtr(timezone)
	return _options
}

// SetScmSource : Allow user to set ScmSource
func (_options *UpdateTektonPipelineTriggerOptions) SetScmSource(scmSource *TriggerScmSource) *UpdateTektonPipelineTriggerOptions {
	_options.ScmSource = scmSource
	return _options
}

// SetEvents : Allow user to set Events
func (_options *UpdateTektonPipelineTriggerOptions) SetEvents(events *Events) *UpdateTektonPipelineTriggerOptions {
	_options.Events = events
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateTektonPipelineTriggerOptions) SetHeaders(param map[string]string) *UpdateTektonPipelineTriggerOptions {
	options.Headers = param
	return options
}

// UserInfo : User information.
type UserInfo struct {
	// IBM Cloud IAM ID.
	IamID *string `json:"iam_id" validate:"required"`

	// User Email address.
	Sub *string `json:"sub" validate:"required"`
}

// UnmarshalUserInfo unmarshals an instance of UserInfo from the specified map of raw messages.
func UnmarshalUserInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserInfo)
	err = core.UnmarshalPrimitive(m, "iam_id", &obj.IamID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sub", &obj.Sub)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Worker : Default pipeline worker used to run the pipeline.
type Worker struct {
	// worker name.
	Name *string `json:"name,omitempty"`

	// worker type.
	Type *string `json:"type,omitempty"`

	ID *string `json:"id" validate:"required"`
}

// Constants associated with the Worker.Type property.
// worker type.
const (
	WorkerTypePrivateConst = "private"
	WorkerTypePublicConst = "public"
)

// NewWorker : Instantiate Worker (Generic Model Constructor)
func (*CdTektonPipelineV2) NewWorker(id string) (_model *Worker, err error) {
	_model = &Worker{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalWorker unmarshals an instance of Worker from the specified map of raw messages.
func UnmarshalWorker(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Worker)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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

// WorkerWithID : Worker object with worker ID only.
type WorkerWithID struct {
	ID *string `json:"id" validate:"required"`
}

// NewWorkerWithID : Instantiate WorkerWithID (Generic Model Constructor)
func (*CdTektonPipelineV2) NewWorkerWithID(id string) (_model *WorkerWithID, err error) {
	_model = &WorkerWithID{
		ID: core.StringPtr(id),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

// UnmarshalWorkerWithID unmarshals an instance of WorkerWithID from the specified map of raw messages.
func UnmarshalWorkerWithID(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(WorkerWithID)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerDuplicateTrigger : request body to duplicate trigger.
// This model "extends" Trigger
type TriggerDuplicateTrigger struct {
	// source trigger ID to clone from.
	SourceTriggerID *string `json:"source_trigger_id" validate:"required"`

	// name of the duplicated trigger.
	Name *string `json:"name" validate:"required"`
}

// NewTriggerDuplicateTrigger : Instantiate TriggerDuplicateTrigger (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerDuplicateTrigger(sourceTriggerID string, name string) (_model *TriggerDuplicateTrigger, err error) {
	_model = &TriggerDuplicateTrigger{
		SourceTriggerID: core.StringPtr(sourceTriggerID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*TriggerDuplicateTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerDuplicateTrigger unmarshals an instance of TriggerDuplicateTrigger from the specified map of raw messages.
func UnmarshalTriggerDuplicateTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerDuplicateTrigger)
	err = core.UnmarshalPrimitive(m, "source_trigger_id", &obj.SourceTriggerID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerGenericTrigger : Generic trigger, which triggers pipeline when tekton pipeline service receive a valie POST event with secrets.
// This model "extends" Trigger
type TriggerGenericTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggerGenericTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// Needed only for generic trigger type. Secret used to start generic trigger.
	Secret *GenericSecret `json:"secret,omitempty"`
}

// NewTriggerGenericTrigger : Instantiate TriggerGenericTrigger (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerGenericTrigger(typeVar string, name string, eventListener string, disabled bool) (_model *TriggerGenericTrigger, err error) {
	_model = &TriggerGenericTrigger{
		Type: core.StringPtr(typeVar),
		Name: core.StringPtr(name),
		EventListener: core.StringPtr(eventListener),
		Disabled: core.BoolPtr(disabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*TriggerGenericTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerGenericTrigger unmarshals an instance of TriggerGenericTrigger from the specified map of raw messages.
func UnmarshalTriggerGenericTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerGenericTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerGenericTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
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

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggerManualTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`
}

// NewTriggerManualTrigger : Instantiate TriggerManualTrigger (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerManualTrigger(typeVar string, name string, eventListener string, disabled bool) (_model *TriggerManualTrigger, err error) {
	_model = &TriggerManualTrigger{
		Type: core.StringPtr(typeVar),
		Name: core.StringPtr(name),
		EventListener: core.StringPtr(eventListener),
		Disabled: core.BoolPtr(disabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*TriggerManualTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerManualTrigger unmarshals an instance of TriggerManualTrigger from the specified map of raw messages.
func UnmarshalTriggerManualTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerManualTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerManualTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerScmTrigger : Git type trigger, which automatically triggers pipelineRun when tekton pipeline service receive a valid corresponding
// git event.
// This model "extends" Trigger
type TriggerScmTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggerScmTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// Scm source for git type tekton pipeline trigger.
	ScmSource *TriggerScmSource `json:"scm_source,omitempty"`

	// Needed only for git trigger type. Events object defines the events this git trigger listening to.
	Events *Events `json:"events,omitempty"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`
}

// NewTriggerScmTrigger : Instantiate TriggerScmTrigger (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerScmTrigger(typeVar string, name string, eventListener string, disabled bool) (_model *TriggerScmTrigger, err error) {
	_model = &TriggerScmTrigger{
		Type: core.StringPtr(typeVar),
		Name: core.StringPtr(name),
		EventListener: core.StringPtr(eventListener),
		Disabled: core.BoolPtr(disabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*TriggerScmTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerScmTrigger unmarshals an instance of TriggerScmTrigger from the specified map of raw messages.
func UnmarshalTriggerScmTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerScmTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerScmTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scm_source", &obj.ScmSource, UnmarshalTriggerScmSource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "events", &obj.Events, UnmarshalEvents)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggerTimerTrigger : Timer trigger, which triggers pipelineRun according to the cron value and time zone.
// This model "extends" Trigger
type TriggerTimerTrigger struct {
	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggerTimerTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// Needed only for timer trigger type. Cron expression for timer trigger. Maximum frequency is every 5 minutes.
	Cron *string `json:"cron,omitempty"`

	// Needed only for timer trigger type. Timezones for timer trigger.
	Timezone *string `json:"timezone,omitempty"`
}

// NewTriggerTimerTrigger : Instantiate TriggerTimerTrigger (Generic Model Constructor)
func (*CdTektonPipelineV2) NewTriggerTimerTrigger(typeVar string, name string, eventListener string, disabled bool) (_model *TriggerTimerTrigger, err error) {
	_model = &TriggerTimerTrigger{
		Type: core.StringPtr(typeVar),
		Name: core.StringPtr(name),
		EventListener: core.StringPtr(eventListener),
		Disabled: core.BoolPtr(disabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	return
}

func (*TriggerTimerTrigger) isaTrigger() bool {
	return true
}

// UnmarshalTriggerTimerTrigger unmarshals an instance of TriggerTimerTrigger from the specified map of raw messages.
func UnmarshalTriggerTimerTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggerTimerTrigger)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggerTimerTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerDuplicateTrigger : request body to duplicate trigger.
// This model "extends" TriggersTriggersItem
type TriggersTriggersItemTriggerDuplicateTrigger struct {
	// General href URL.
	Href *string `json:"href,omitempty"`

	// source trigger ID to clone from.
	SourceTriggerID *string `json:"source_trigger_id" validate:"required"`

	// name of the duplicated trigger.
	Name *string `json:"name" validate:"required"`
}

func (*TriggersTriggersItemTriggerDuplicateTrigger) isaTriggersTriggersItem() bool {
	return true
}

// UnmarshalTriggersTriggersItemTriggerDuplicateTrigger unmarshals an instance of TriggersTriggersItemTriggerDuplicateTrigger from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerDuplicateTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerDuplicateTrigger)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_trigger_id", &obj.SourceTriggerID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerGenericTrigger : Generic trigger, which triggers pipeline when tekton pipeline service receive a valie POST event with secrets.
// This model "extends" TriggersTriggersItem
type TriggersTriggersItemTriggerGenericTrigger struct {
	// General href URL.
	Href *string `json:"href,omitempty"`

	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggersTriggersItemTriggerGenericTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// Needed only for generic trigger type. Secret used to start generic trigger.
	Secret *GenericSecret `json:"secret,omitempty"`
}

func (*TriggersTriggersItemTriggerGenericTrigger) isaTriggersTriggersItem() bool {
	return true
}

// UnmarshalTriggersTriggersItemTriggerGenericTrigger unmarshals an instance of TriggersTriggersItemTriggerGenericTrigger from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerGenericTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerGenericTrigger)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggersTriggersItemTriggerGenericTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "secret", &obj.Secret, UnmarshalGenericSecret)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerManualTrigger : Manual trigger.
// This model "extends" TriggersTriggersItem
type TriggersTriggersItemTriggerManualTrigger struct {
	// General href URL.
	Href *string `json:"href,omitempty"`

	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggersTriggersItemTriggerManualTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`
}

func (*TriggersTriggersItemTriggerManualTrigger) isaTriggersTriggersItem() bool {
	return true
}

// UnmarshalTriggersTriggersItemTriggerManualTrigger unmarshals an instance of TriggersTriggersItemTriggerManualTrigger from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerManualTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerManualTrigger)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggersTriggersItemTriggerManualTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerScmTrigger : Git type trigger, which automatically triggers pipelineRun when tekton pipeline service receive a valid corresponding
// git event.
// This model "extends" TriggersTriggersItem
type TriggersTriggersItemTriggerScmTrigger struct {
	// General href URL.
	Href *string `json:"href,omitempty"`

	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggersTriggersItemTriggerScmTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// Scm source for git type tekton pipeline trigger.
	ScmSource *TriggerScmSource `json:"scm_source,omitempty"`

	// Needed only for git trigger type. Events object defines the events this git trigger listening to.
	Events *Events `json:"events,omitempty"`

	// UUID.
	ServiceInstanceID *string `json:"service_instance_id,omitempty"`
}

func (*TriggersTriggersItemTriggerScmTrigger) isaTriggersTriggersItem() bool {
	return true
}

// UnmarshalTriggersTriggersItemTriggerScmTrigger unmarshals an instance of TriggersTriggersItemTriggerScmTrigger from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerScmTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerScmTrigger)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggersTriggersItemTriggerScmTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "scm_source", &obj.ScmSource, UnmarshalTriggerScmSource)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "events", &obj.Events, UnmarshalEvents)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_instance_id", &obj.ServiceInstanceID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TriggersTriggersItemTriggerTimerTrigger : Timer trigger, which triggers pipelineRun according to the cron value and time zone.
// This model "extends" TriggersTriggersItem
type TriggersTriggersItemTriggerTimerTrigger struct {
	// General href URL.
	Href *string `json:"href,omitempty"`

	// Trigger type.
	Type *string `json:"type" validate:"required"`

	// Trigger name.
	Name *string `json:"name" validate:"required"`

	// Event listener name.
	EventListener *string `json:"event_listener" validate:"required"`

	// UUID.
	ID *string `json:"id,omitempty"`

	// Trigger properties.
	Properties []TriggersTriggersItemTriggerTimerTriggerPropertiesItem `json:"properties,omitempty"`

	// Trigger tags array.
	Tags []string `json:"tags,omitempty"`

	// Trigger worker used to run the trigger, the trigger worker overrides the default pipeline worker.If not exist, this
	// trigger uses default pipeline worker.
	Worker *Worker `json:"worker,omitempty"`

	// Concurrency object.
	Concurrency *Concurrency `json:"concurrency,omitempty"`

	// flag whether the trigger is disabled.
	Disabled *bool `json:"disabled" validate:"required"`

	// Needed only for timer trigger type. Cron expression for timer trigger. Maximum frequency is every 5 minutes.
	Cron *string `json:"cron,omitempty"`

	// Needed only for timer trigger type. Timezones for timer trigger.
	Timezone *string `json:"timezone,omitempty"`
}

func (*TriggersTriggersItemTriggerTimerTrigger) isaTriggersTriggersItem() bool {
	return true
}

// UnmarshalTriggersTriggersItemTriggerTimerTrigger unmarshals an instance of TriggersTriggersItemTriggerTimerTrigger from the specified map of raw messages.
func UnmarshalTriggersTriggersItemTriggerTimerTrigger(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TriggersTriggersItemTriggerTimerTrigger)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "event_listener", &obj.EventListener)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalTriggersTriggersItemTriggerTimerTriggerPropertiesItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "worker", &obj.Worker, UnmarshalWorker)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "concurrency", &obj.Concurrency, UnmarshalConcurrency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled", &obj.Disabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "cron", &obj.Cron)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timezone", &obj.Timezone)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
