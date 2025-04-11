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

// Package codeenginev2 : Operations and models for the CodeEngineV2 service
package codeenginev2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM/code-engine-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
)

// CodeEngineV2 : REST API for Code Engine
//
// API Version: 2.0.0
type CodeEngineV2 struct {
	Service *core.BaseService

	// The API version, in format `YYYY-MM-DD`. For the API behavior documented here, specify any date between `2021-03-31`
	// and `2024-12-10`.
	Version *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.au-syd.codeengine.cloud.ibm.com/v2"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "code_engine"

// CodeEngineV2Options : Service options
type CodeEngineV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// The API version, in format `YYYY-MM-DD`. For the API behavior documented here, specify any date between `2021-03-31`
	// and `2024-12-10`.
	Version *string 
}

// NewCodeEngineV2UsingExternalConfig : constructs an instance of CodeEngineV2 with passed in options and external configuration.
func NewCodeEngineV2UsingExternalConfig(options *CodeEngineV2Options) (codeEngine *CodeEngineV2, err error) {
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

	codeEngine, err = NewCodeEngineV2(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = codeEngine.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = codeEngine.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewCodeEngineV2 : constructs an instance of CodeEngineV2 with passed in options.
func NewCodeEngineV2(options *CodeEngineV2Options) (service *CodeEngineV2, err error) {
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

	service = &CodeEngineV2{
		Service: baseService,
		Version: options.Version,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "codeEngine" suitable for processing requests.
func (codeEngine *CodeEngineV2) Clone() *CodeEngineV2 {
	if core.IsNil(codeEngine) {
		return nil
	}
	clone := *codeEngine
	clone.Service = codeEngine.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (codeEngine *CodeEngineV2) SetServiceURL(url string) error {
	err := codeEngine.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (codeEngine *CodeEngineV2) GetServiceURL() string {
	return codeEngine.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (codeEngine *CodeEngineV2) SetDefaultHeaders(headers http.Header) {
	codeEngine.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (codeEngine *CodeEngineV2) SetEnableGzipCompression(enableGzip bool) {
	codeEngine.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (codeEngine *CodeEngineV2) GetEnableGzipCompression() bool {
	return codeEngine.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (codeEngine *CodeEngineV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	codeEngine.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (codeEngine *CodeEngineV2) DisableRetries() {
	codeEngine.Service.DisableRetries()
}

// ListProjects : List all projects
// List all projects in the current account.
func (codeEngine *CodeEngineV2) ListProjects(listProjectsOptions *ListProjectsOptions) (result *ProjectList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListProjectsWithContext(context.Background(), listProjectsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListProjectsWithContext is an alternate form of the ListProjects method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListProjectsWithContext(ctx context.Context, listProjectsOptions *ListProjectsOptions) (result *ProjectList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listProjectsOptions, "listProjectsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listProjectsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListProjects")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listProjectsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listProjectsOptions.Limit))
	}
	if listProjectsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listProjectsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_projects", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateProject : Create a project
// Create a Code Engine project on IBM Cloud. The project will be created in the region that corresponds to the API
// endpoint that is being called.
func (codeEngine *CodeEngineV2) CreateProject(createProjectOptions *CreateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateProjectWithContext(context.Background(), createProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateProjectWithContext is an alternate form of the CreateProject method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateProjectWithContext(ctx context.Context, createProjectOptions *CreateProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createProjectOptions, "createProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createProjectOptions, "createProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateProject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createProjectOptions.Name != nil {
		body["name"] = createProjectOptions.Name
	}
	if createProjectOptions.ResourceGroupID != nil {
		body["resource_group_id"] = createProjectOptions.ResourceGroupID
	}
	if createProjectOptions.Tags != nil {
		body["tags"] = createProjectOptions.Tags
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProject : Get a project
// Display the details of a single project.
func (codeEngine *CodeEngineV2) GetProject(getProjectOptions *GetProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetProjectWithContext(context.Background(), getProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProjectWithContext is an alternate form of the GetProject method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetProjectWithContext(ctx context.Context, getProjectOptions *GetProjectOptions) (result *Project, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectOptions, "getProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProjectOptions, "getProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *getProjectOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetProject")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProject)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteProject : Delete a project
// Delete a project.
func (codeEngine *CodeEngineV2) DeleteProject(deleteProjectOptions *DeleteProjectOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteProjectWithContext(context.Background(), deleteProjectOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteProjectWithContext is an alternate form of the DeleteProject method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteProjectWithContext(ctx context.Context, deleteProjectOptions *DeleteProjectOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteProjectOptions, "deleteProjectOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteProjectOptions, "deleteProjectOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"id": *deleteProjectOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteProjectOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteProject")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_project", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListAllowedOutboundDestination : List allowed outbound destinations
// List all allowed outbound destinations in a project.
func (codeEngine *CodeEngineV2) ListAllowedOutboundDestination(listAllowedOutboundDestinationOptions *ListAllowedOutboundDestinationOptions) (result *AllowedOutboundDestinationList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListAllowedOutboundDestinationWithContext(context.Background(), listAllowedOutboundDestinationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAllowedOutboundDestinationWithContext is an alternate form of the ListAllowedOutboundDestination method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListAllowedOutboundDestinationWithContext(ctx context.Context, listAllowedOutboundDestinationOptions *ListAllowedOutboundDestinationOptions) (result *AllowedOutboundDestinationList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAllowedOutboundDestinationOptions, "listAllowedOutboundDestinationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listAllowedOutboundDestinationOptions, "listAllowedOutboundDestinationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listAllowedOutboundDestinationOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/allowed_outbound_destinations`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAllowedOutboundDestinationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListAllowedOutboundDestination")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAllowedOutboundDestinationOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAllowedOutboundDestinationOptions.Limit))
	}
	if listAllowedOutboundDestinationOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listAllowedOutboundDestinationOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_allowed_outbound_destination", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllowedOutboundDestinationList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateAllowedOutboundDestination : Create an allowed outbound destination
// Create an allowed outbound destination.
func (codeEngine *CodeEngineV2) CreateAllowedOutboundDestination(createAllowedOutboundDestinationOptions *CreateAllowedOutboundDestinationOptions) (result AllowedOutboundDestinationIntf, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateAllowedOutboundDestinationWithContext(context.Background(), createAllowedOutboundDestinationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAllowedOutboundDestinationWithContext is an alternate form of the CreateAllowedOutboundDestination method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateAllowedOutboundDestinationWithContext(ctx context.Context, createAllowedOutboundDestinationOptions *CreateAllowedOutboundDestinationOptions) (result AllowedOutboundDestinationIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAllowedOutboundDestinationOptions, "createAllowedOutboundDestinationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAllowedOutboundDestinationOptions, "createAllowedOutboundDestinationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createAllowedOutboundDestinationOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/allowed_outbound_destinations`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAllowedOutboundDestinationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateAllowedOutboundDestination")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	_, err = builder.SetBodyContentJSON(createAllowedOutboundDestinationOptions.AllowedOutboundDestination)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_allowed_outbound_destination", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllowedOutboundDestination)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAllowedOutboundDestination : Get an allowed outbound destination
// Display the details of an allowed outbound destination.
func (codeEngine *CodeEngineV2) GetAllowedOutboundDestination(getAllowedOutboundDestinationOptions *GetAllowedOutboundDestinationOptions) (result AllowedOutboundDestinationIntf, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetAllowedOutboundDestinationWithContext(context.Background(), getAllowedOutboundDestinationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAllowedOutboundDestinationWithContext is an alternate form of the GetAllowedOutboundDestination method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetAllowedOutboundDestinationWithContext(ctx context.Context, getAllowedOutboundDestinationOptions *GetAllowedOutboundDestinationOptions) (result AllowedOutboundDestinationIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAllowedOutboundDestinationOptions, "getAllowedOutboundDestinationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAllowedOutboundDestinationOptions, "getAllowedOutboundDestinationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getAllowedOutboundDestinationOptions.ProjectID,
		"name": *getAllowedOutboundDestinationOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/allowed_outbound_destinations/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAllowedOutboundDestinationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetAllowedOutboundDestination")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_allowed_outbound_destination", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllowedOutboundDestination)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAllowedOutboundDestination : Delete an allowed outbound destination
// Delete an allowed outbound destination.
func (codeEngine *CodeEngineV2) DeleteAllowedOutboundDestination(deleteAllowedOutboundDestinationOptions *DeleteAllowedOutboundDestinationOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteAllowedOutboundDestinationWithContext(context.Background(), deleteAllowedOutboundDestinationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAllowedOutboundDestinationWithContext is an alternate form of the DeleteAllowedOutboundDestination method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteAllowedOutboundDestinationWithContext(ctx context.Context, deleteAllowedOutboundDestinationOptions *DeleteAllowedOutboundDestinationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAllowedOutboundDestinationOptions, "deleteAllowedOutboundDestinationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAllowedOutboundDestinationOptions, "deleteAllowedOutboundDestinationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteAllowedOutboundDestinationOptions.ProjectID,
		"name": *deleteAllowedOutboundDestinationOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/allowed_outbound_destinations/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAllowedOutboundDestinationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteAllowedOutboundDestination")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_allowed_outbound_destination", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateAllowedOutboundDestination : Update an allowed outbound destination
// Update an allowed outbound destination.
func (codeEngine *CodeEngineV2) UpdateAllowedOutboundDestination(updateAllowedOutboundDestinationOptions *UpdateAllowedOutboundDestinationOptions) (result AllowedOutboundDestinationIntf, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.UpdateAllowedOutboundDestinationWithContext(context.Background(), updateAllowedOutboundDestinationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAllowedOutboundDestinationWithContext is an alternate form of the UpdateAllowedOutboundDestination method which supports a Context parameter
func (codeEngine *CodeEngineV2) UpdateAllowedOutboundDestinationWithContext(ctx context.Context, updateAllowedOutboundDestinationOptions *UpdateAllowedOutboundDestinationOptions) (result AllowedOutboundDestinationIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAllowedOutboundDestinationOptions, "updateAllowedOutboundDestinationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAllowedOutboundDestinationOptions, "updateAllowedOutboundDestinationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateAllowedOutboundDestinationOptions.ProjectID,
		"name": *updateAllowedOutboundDestinationOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/allowed_outbound_destinations/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAllowedOutboundDestinationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "UpdateAllowedOutboundDestination")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateAllowedOutboundDestinationOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateAllowedOutboundDestinationOptions.IfMatch))
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	_, err = builder.SetBodyContentJSON(updateAllowedOutboundDestinationOptions.AllowedOutboundDestination)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_allowed_outbound_destination", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAllowedOutboundDestination)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProjectEgressIps : List egress IP addresses
// Lists all egress IP addresses (public and private) that are used by components running in this project. For
// information about using egress IP addresses, see [Code Engine public and private IP
// addresses](https://cloud.ibm.com/docs/codeengine?topic=codeengine-network-addresses).
func (codeEngine *CodeEngineV2) GetProjectEgressIps(getProjectEgressIpsOptions *GetProjectEgressIpsOptions) (result *ProjectEgressIPAddresses, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetProjectEgressIpsWithContext(context.Background(), getProjectEgressIpsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProjectEgressIpsWithContext is an alternate form of the GetProjectEgressIps method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetProjectEgressIpsWithContext(ctx context.Context, getProjectEgressIpsOptions *GetProjectEgressIpsOptions) (result *ProjectEgressIPAddresses, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectEgressIpsOptions, "getProjectEgressIpsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProjectEgressIpsOptions, "getProjectEgressIpsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getProjectEgressIpsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/egress_ips`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProjectEgressIpsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetProjectEgressIps")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_project_egress_ips", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectEgressIPAddresses)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetProjectStatusDetails : Get the status details for a project
// Retrieves status details about the given project.
func (codeEngine *CodeEngineV2) GetProjectStatusDetails(getProjectStatusDetailsOptions *GetProjectStatusDetailsOptions) (result *ProjectStatusDetails, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetProjectStatusDetailsWithContext(context.Background(), getProjectStatusDetailsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetProjectStatusDetailsWithContext is an alternate form of the GetProjectStatusDetails method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetProjectStatusDetailsWithContext(ctx context.Context, getProjectStatusDetailsOptions *GetProjectStatusDetailsOptions) (result *ProjectStatusDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getProjectStatusDetailsOptions, "getProjectStatusDetailsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getProjectStatusDetailsOptions, "getProjectStatusDetailsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getProjectStatusDetailsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/status_details`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getProjectStatusDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetProjectStatusDetails")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_project_status_details", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProjectStatusDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListApps : List applications
// List all applications in a project.
func (codeEngine *CodeEngineV2) ListApps(listAppsOptions *ListAppsOptions) (result *AppList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListAppsWithContext(context.Background(), listAppsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAppsWithContext is an alternate form of the ListApps method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListAppsWithContext(ctx context.Context, listAppsOptions *ListAppsOptions) (result *AppList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAppsOptions, "listAppsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listAppsOptions, "listAppsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listAppsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAppsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListApps")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}
	if listAppsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAppsOptions.Limit))
	}
	if listAppsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listAppsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_apps", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAppList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateApp : Create an application
// Create an application.
func (codeEngine *CodeEngineV2) CreateApp(createAppOptions *CreateAppOptions) (result *App, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateAppWithContext(context.Background(), createAppOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateAppWithContext is an alternate form of the CreateApp method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateAppWithContext(ctx context.Context, createAppOptions *CreateAppOptions) (result *App, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createAppOptions, "createAppOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createAppOptions, "createAppOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createAppOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	body := make(map[string]interface{})
	if createAppOptions.ImageReference != nil {
		body["image_reference"] = createAppOptions.ImageReference
	}
	if createAppOptions.Name != nil {
		body["name"] = createAppOptions.Name
	}
	if createAppOptions.ImagePort != nil {
		body["image_port"] = createAppOptions.ImagePort
	}
	if createAppOptions.ImageSecret != nil {
		body["image_secret"] = createAppOptions.ImageSecret
	}
	if createAppOptions.ManagedDomainMappings != nil {
		body["managed_domain_mappings"] = createAppOptions.ManagedDomainMappings
	}
	if createAppOptions.ProbeLiveness != nil {
		body["probe_liveness"] = createAppOptions.ProbeLiveness
	}
	if createAppOptions.ProbeReadiness != nil {
		body["probe_readiness"] = createAppOptions.ProbeReadiness
	}
	if createAppOptions.RunArguments != nil {
		body["run_arguments"] = createAppOptions.RunArguments
	}
	if createAppOptions.RunAsUser != nil {
		body["run_as_user"] = createAppOptions.RunAsUser
	}
	if createAppOptions.RunCommands != nil {
		body["run_commands"] = createAppOptions.RunCommands
	}
	if createAppOptions.RunEnvVariables != nil {
		body["run_env_variables"] = createAppOptions.RunEnvVariables
	}
	if createAppOptions.RunServiceAccount != nil {
		body["run_service_account"] = createAppOptions.RunServiceAccount
	}
	if createAppOptions.RunVolumeMounts != nil {
		body["run_volume_mounts"] = createAppOptions.RunVolumeMounts
	}
	if createAppOptions.ScaleConcurrency != nil {
		body["scale_concurrency"] = createAppOptions.ScaleConcurrency
	}
	if createAppOptions.ScaleConcurrencyTarget != nil {
		body["scale_concurrency_target"] = createAppOptions.ScaleConcurrencyTarget
	}
	if createAppOptions.ScaleCpuLimit != nil {
		body["scale_cpu_limit"] = createAppOptions.ScaleCpuLimit
	}
	if createAppOptions.ScaleDownDelay != nil {
		body["scale_down_delay"] = createAppOptions.ScaleDownDelay
	}
	if createAppOptions.ScaleEphemeralStorageLimit != nil {
		body["scale_ephemeral_storage_limit"] = createAppOptions.ScaleEphemeralStorageLimit
	}
	if createAppOptions.ScaleInitialInstances != nil {
		body["scale_initial_instances"] = createAppOptions.ScaleInitialInstances
	}
	if createAppOptions.ScaleMaxInstances != nil {
		body["scale_max_instances"] = createAppOptions.ScaleMaxInstances
	}
	if createAppOptions.ScaleMemoryLimit != nil {
		body["scale_memory_limit"] = createAppOptions.ScaleMemoryLimit
	}
	if createAppOptions.ScaleMinInstances != nil {
		body["scale_min_instances"] = createAppOptions.ScaleMinInstances
	}
	if createAppOptions.ScaleRequestTimeout != nil {
		body["scale_request_timeout"] = createAppOptions.ScaleRequestTimeout
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_app", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetApp : Get an application
// Display the details of an application.
func (codeEngine *CodeEngineV2) GetApp(getAppOptions *GetAppOptions) (result *App, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetAppWithContext(context.Background(), getAppOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAppWithContext is an alternate form of the GetApp method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetAppWithContext(ctx context.Context, getAppOptions *GetAppOptions) (result *App, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAppOptions, "getAppOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAppOptions, "getAppOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getAppOptions.ProjectID,
		"name": *getAppOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_app", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteApp : Delete an application
// Delete an application.
func (codeEngine *CodeEngineV2) DeleteApp(deleteAppOptions *DeleteAppOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteAppWithContext(context.Background(), deleteAppOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAppWithContext is an alternate form of the DeleteApp method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteAppWithContext(ctx context.Context, deleteAppOptions *DeleteAppOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAppOptions, "deleteAppOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAppOptions, "deleteAppOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteAppOptions.ProjectID,
		"name": *deleteAppOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_app", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateApp : Update an application
// An application contains one or more revisions. A revision represents an immutable version of the configuration
// properties of the application. Each update of an application configuration property creates a new revision of the
// application. [Learn more](https://cloud.ibm.com/docs/codeengine?topic=codeengine-update-app).
func (codeEngine *CodeEngineV2) UpdateApp(updateAppOptions *UpdateAppOptions) (result *App, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.UpdateAppWithContext(context.Background(), updateAppOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateAppWithContext is an alternate form of the UpdateApp method which supports a Context parameter
func (codeEngine *CodeEngineV2) UpdateAppWithContext(ctx context.Context, updateAppOptions *UpdateAppOptions) (result *App, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateAppOptions, "updateAppOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateAppOptions, "updateAppOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateAppOptions.ProjectID,
		"name": *updateAppOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "UpdateApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateAppOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateAppOptions.IfMatch))
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	_, err = builder.SetBodyContentJSON(updateAppOptions.App)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_app", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApp)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListAppRevisions : List application revisions
// List all application revisions in a particular application.
func (codeEngine *CodeEngineV2) ListAppRevisions(listAppRevisionsOptions *ListAppRevisionsOptions) (result *AppRevisionList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListAppRevisionsWithContext(context.Background(), listAppRevisionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAppRevisionsWithContext is an alternate form of the ListAppRevisions method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListAppRevisionsWithContext(ctx context.Context, listAppRevisionsOptions *ListAppRevisionsOptions) (result *AppRevisionList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAppRevisionsOptions, "listAppRevisionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listAppRevisionsOptions, "listAppRevisionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listAppRevisionsOptions.ProjectID,
		"app_name": *listAppRevisionsOptions.AppName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{app_name}/revisions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAppRevisionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListAppRevisions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAppRevisionsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAppRevisionsOptions.Limit))
	}
	if listAppRevisionsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listAppRevisionsOptions.Start))
	}
	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_app_revisions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAppRevisionList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetAppRevision : Get an application revision
// Display the details of an application revision.
func (codeEngine *CodeEngineV2) GetAppRevision(getAppRevisionOptions *GetAppRevisionOptions) (result *AppRevision, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetAppRevisionWithContext(context.Background(), getAppRevisionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAppRevisionWithContext is an alternate form of the GetAppRevision method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetAppRevisionWithContext(ctx context.Context, getAppRevisionOptions *GetAppRevisionOptions) (result *AppRevision, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getAppRevisionOptions, "getAppRevisionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getAppRevisionOptions, "getAppRevisionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getAppRevisionOptions.ProjectID,
		"app_name": *getAppRevisionOptions.AppName,
		"name": *getAppRevisionOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{app_name}/revisions/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getAppRevisionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetAppRevision")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_app_revision", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAppRevision)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteAppRevision : Delete an application revision
// Delete an application revision.
func (codeEngine *CodeEngineV2) DeleteAppRevision(deleteAppRevisionOptions *DeleteAppRevisionOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteAppRevisionWithContext(context.Background(), deleteAppRevisionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteAppRevisionWithContext is an alternate form of the DeleteAppRevision method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteAppRevisionWithContext(ctx context.Context, deleteAppRevisionOptions *DeleteAppRevisionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteAppRevisionOptions, "deleteAppRevisionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteAppRevisionOptions, "deleteAppRevisionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteAppRevisionOptions.ProjectID,
		"app_name": *deleteAppRevisionOptions.AppName,
		"name": *deleteAppRevisionOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{app_name}/revisions/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteAppRevisionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteAppRevision")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_app_revision", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListAppInstances : List application instances
// List all instances of an application.
func (codeEngine *CodeEngineV2) ListAppInstances(listAppInstancesOptions *ListAppInstancesOptions) (result *AppInstanceList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListAppInstancesWithContext(context.Background(), listAppInstancesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListAppInstancesWithContext is an alternate form of the ListAppInstances method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListAppInstancesWithContext(ctx context.Context, listAppInstancesOptions *ListAppInstancesOptions) (result *AppInstanceList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listAppInstancesOptions, "listAppInstancesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listAppInstancesOptions, "listAppInstancesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listAppInstancesOptions.ProjectID,
		"app_name": *listAppInstancesOptions.AppName,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/apps/{app_name}/instances`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listAppInstancesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListAppInstances")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listAppInstancesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listAppInstancesOptions.Limit))
	}
	if listAppInstancesOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listAppInstancesOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_app_instances", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalAppInstanceList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListJobs : List jobs
// List all jobs in a project.
func (codeEngine *CodeEngineV2) ListJobs(listJobsOptions *ListJobsOptions) (result *JobList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListJobsWithContext(context.Background(), listJobsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListJobsWithContext is an alternate form of the ListJobs method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListJobsWithContext(ctx context.Context, listJobsOptions *ListJobsOptions) (result *JobList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listJobsOptions, "listJobsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listJobsOptions, "listJobsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listJobsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/jobs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listJobsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListJobs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}
	if listJobsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listJobsOptions.Limit))
	}
	if listJobsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listJobsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_jobs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateJob : Create a job
// Create a job.
func (codeEngine *CodeEngineV2) CreateJob(createJobOptions *CreateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateJobWithContext(context.Background(), createJobOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateJobWithContext is an alternate form of the CreateJob method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateJobWithContext(ctx context.Context, createJobOptions *CreateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createJobOptions, "createJobOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createJobOptions, "createJobOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createJobOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/jobs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	body := make(map[string]interface{})
	if createJobOptions.ImageReference != nil {
		body["image_reference"] = createJobOptions.ImageReference
	}
	if createJobOptions.Name != nil {
		body["name"] = createJobOptions.Name
	}
	if createJobOptions.ImageSecret != nil {
		body["image_secret"] = createJobOptions.ImageSecret
	}
	if createJobOptions.RunArguments != nil {
		body["run_arguments"] = createJobOptions.RunArguments
	}
	if createJobOptions.RunAsUser != nil {
		body["run_as_user"] = createJobOptions.RunAsUser
	}
	if createJobOptions.RunCommands != nil {
		body["run_commands"] = createJobOptions.RunCommands
	}
	if createJobOptions.RunEnvVariables != nil {
		body["run_env_variables"] = createJobOptions.RunEnvVariables
	}
	if createJobOptions.RunMode != nil {
		body["run_mode"] = createJobOptions.RunMode
	}
	if createJobOptions.RunServiceAccount != nil {
		body["run_service_account"] = createJobOptions.RunServiceAccount
	}
	if createJobOptions.RunVolumeMounts != nil {
		body["run_volume_mounts"] = createJobOptions.RunVolumeMounts
	}
	if createJobOptions.ScaleArraySpec != nil {
		body["scale_array_spec"] = createJobOptions.ScaleArraySpec
	}
	if createJobOptions.ScaleCpuLimit != nil {
		body["scale_cpu_limit"] = createJobOptions.ScaleCpuLimit
	}
	if createJobOptions.ScaleEphemeralStorageLimit != nil {
		body["scale_ephemeral_storage_limit"] = createJobOptions.ScaleEphemeralStorageLimit
	}
	if createJobOptions.ScaleMaxExecutionTime != nil {
		body["scale_max_execution_time"] = createJobOptions.ScaleMaxExecutionTime
	}
	if createJobOptions.ScaleMemoryLimit != nil {
		body["scale_memory_limit"] = createJobOptions.ScaleMemoryLimit
	}
	if createJobOptions.ScaleRetryLimit != nil {
		body["scale_retry_limit"] = createJobOptions.ScaleRetryLimit
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_job", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetJob : Get a job
// Display the details of a job.
func (codeEngine *CodeEngineV2) GetJob(getJobOptions *GetJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetJobWithContext(context.Background(), getJobOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetJobWithContext is an alternate form of the GetJob method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetJobWithContext(ctx context.Context, getJobOptions *GetJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getJobOptions, "getJobOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getJobOptions, "getJobOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getJobOptions.ProjectID,
		"name": *getJobOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/jobs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_job", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteJob : Delete a job
// Delete a job.
func (codeEngine *CodeEngineV2) DeleteJob(deleteJobOptions *DeleteJobOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteJobWithContext(context.Background(), deleteJobOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteJobWithContext is an alternate form of the DeleteJob method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteJobWithContext(ctx context.Context, deleteJobOptions *DeleteJobOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteJobOptions, "deleteJobOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteJobOptions, "deleteJobOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteJobOptions.ProjectID,
		"name": *deleteJobOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/jobs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_job", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateJob : Update a job
// Update the given job.
func (codeEngine *CodeEngineV2) UpdateJob(updateJobOptions *UpdateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.UpdateJobWithContext(context.Background(), updateJobOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateJobWithContext is an alternate form of the UpdateJob method which supports a Context parameter
func (codeEngine *CodeEngineV2) UpdateJobWithContext(ctx context.Context, updateJobOptions *UpdateJobOptions) (result *Job, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateJobOptions, "updateJobOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateJobOptions, "updateJobOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateJobOptions.ProjectID,
		"name": *updateJobOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/jobs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateJobOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "UpdateJob")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateJobOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateJobOptions.IfMatch))
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	_, err = builder.SetBodyContentJSON(updateJobOptions.Job)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_job", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJob)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListJobRuns : List job runs
// List all job runs in a project.
func (codeEngine *CodeEngineV2) ListJobRuns(listJobRunsOptions *ListJobRunsOptions) (result *JobRunList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListJobRunsWithContext(context.Background(), listJobRunsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListJobRunsWithContext is an alternate form of the ListJobRuns method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListJobRunsWithContext(ctx context.Context, listJobRunsOptions *ListJobRunsOptions) (result *JobRunList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listJobRunsOptions, "listJobRunsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listJobRunsOptions, "listJobRunsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listJobRunsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/job_runs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listJobRunsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListJobRuns")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}
	if listJobRunsOptions.JobName != nil {
		builder.AddQuery("job_name", fmt.Sprint(*listJobRunsOptions.JobName))
	}
	if listJobRunsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listJobRunsOptions.Limit))
	}
	if listJobRunsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listJobRunsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_job_runs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobRunList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateJobRun : Create a job run
// Create an job run.
func (codeEngine *CodeEngineV2) CreateJobRun(createJobRunOptions *CreateJobRunOptions) (result *JobRun, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateJobRunWithContext(context.Background(), createJobRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateJobRunWithContext is an alternate form of the CreateJobRun method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateJobRunWithContext(ctx context.Context, createJobRunOptions *CreateJobRunOptions) (result *JobRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createJobRunOptions, "createJobRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createJobRunOptions, "createJobRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createJobRunOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/job_runs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createJobRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateJobRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	body := make(map[string]interface{})
	if createJobRunOptions.ImageReference != nil {
		body["image_reference"] = createJobRunOptions.ImageReference
	}
	if createJobRunOptions.ImageSecret != nil {
		body["image_secret"] = createJobRunOptions.ImageSecret
	}
	if createJobRunOptions.JobName != nil {
		body["job_name"] = createJobRunOptions.JobName
	}
	if createJobRunOptions.Name != nil {
		body["name"] = createJobRunOptions.Name
	}
	if createJobRunOptions.RunArguments != nil {
		body["run_arguments"] = createJobRunOptions.RunArguments
	}
	if createJobRunOptions.RunAsUser != nil {
		body["run_as_user"] = createJobRunOptions.RunAsUser
	}
	if createJobRunOptions.RunCommands != nil {
		body["run_commands"] = createJobRunOptions.RunCommands
	}
	if createJobRunOptions.RunEnvVariables != nil {
		body["run_env_variables"] = createJobRunOptions.RunEnvVariables
	}
	if createJobRunOptions.RunMode != nil {
		body["run_mode"] = createJobRunOptions.RunMode
	}
	if createJobRunOptions.RunServiceAccount != nil {
		body["run_service_account"] = createJobRunOptions.RunServiceAccount
	}
	if createJobRunOptions.RunVolumeMounts != nil {
		body["run_volume_mounts"] = createJobRunOptions.RunVolumeMounts
	}
	if createJobRunOptions.ScaleArraySizeVariableOverride != nil {
		body["scale_array_size_variable_override"] = createJobRunOptions.ScaleArraySizeVariableOverride
	}
	if createJobRunOptions.ScaleArraySpec != nil {
		body["scale_array_spec"] = createJobRunOptions.ScaleArraySpec
	}
	if createJobRunOptions.ScaleCpuLimit != nil {
		body["scale_cpu_limit"] = createJobRunOptions.ScaleCpuLimit
	}
	if createJobRunOptions.ScaleEphemeralStorageLimit != nil {
		body["scale_ephemeral_storage_limit"] = createJobRunOptions.ScaleEphemeralStorageLimit
	}
	if createJobRunOptions.ScaleMaxExecutionTime != nil {
		body["scale_max_execution_time"] = createJobRunOptions.ScaleMaxExecutionTime
	}
	if createJobRunOptions.ScaleMemoryLimit != nil {
		body["scale_memory_limit"] = createJobRunOptions.ScaleMemoryLimit
	}
	if createJobRunOptions.ScaleRetryLimit != nil {
		body["scale_retry_limit"] = createJobRunOptions.ScaleRetryLimit
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_job_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetJobRun : Get a job run
// Display the details of a job run.
func (codeEngine *CodeEngineV2) GetJobRun(getJobRunOptions *GetJobRunOptions) (result *JobRun, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetJobRunWithContext(context.Background(), getJobRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetJobRunWithContext is an alternate form of the GetJobRun method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetJobRunWithContext(ctx context.Context, getJobRunOptions *GetJobRunOptions) (result *JobRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getJobRunOptions, "getJobRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getJobRunOptions, "getJobRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getJobRunOptions.ProjectID,
		"name": *getJobRunOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/job_runs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getJobRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetJobRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_job_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalJobRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteJobRun : Delete a job run
// Delete a job run.
func (codeEngine *CodeEngineV2) DeleteJobRun(deleteJobRunOptions *DeleteJobRunOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteJobRunWithContext(context.Background(), deleteJobRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteJobRunWithContext is an alternate form of the DeleteJobRun method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteJobRunWithContext(ctx context.Context, deleteJobRunOptions *DeleteJobRunOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteJobRunOptions, "deleteJobRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteJobRunOptions, "deleteJobRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteJobRunOptions.ProjectID,
		"name": *deleteJobRunOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/job_runs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteJobRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteJobRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_job_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListFunctionRuntimes : List the function runtimes
// List all valid function runtimes.
func (codeEngine *CodeEngineV2) ListFunctionRuntimes(listFunctionRuntimesOptions *ListFunctionRuntimesOptions) (result *FunctionRuntimeList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListFunctionRuntimesWithContext(context.Background(), listFunctionRuntimesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListFunctionRuntimesWithContext is an alternate form of the ListFunctionRuntimes method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListFunctionRuntimesWithContext(ctx context.Context, listFunctionRuntimesOptions *ListFunctionRuntimesOptions) (result *FunctionRuntimeList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listFunctionRuntimesOptions, "listFunctionRuntimesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/function_runtimes`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listFunctionRuntimesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListFunctionRuntimes")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_function_runtimes", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFunctionRuntimeList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListFunctions : List functions
// List all functions in a project.
func (codeEngine *CodeEngineV2) ListFunctions(listFunctionsOptions *ListFunctionsOptions) (result *FunctionList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListFunctionsWithContext(context.Background(), listFunctionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListFunctionsWithContext is an alternate form of the ListFunctions method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListFunctionsWithContext(ctx context.Context, listFunctionsOptions *ListFunctionsOptions) (result *FunctionList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listFunctionsOptions, "listFunctionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listFunctionsOptions, "listFunctionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listFunctionsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/functions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listFunctionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListFunctions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}
	if listFunctionsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listFunctionsOptions.Limit))
	}
	if listFunctionsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listFunctionsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_functions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFunctionList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateFunction : Create a function
// Create a function.
func (codeEngine *CodeEngineV2) CreateFunction(createFunctionOptions *CreateFunctionOptions) (result *Function, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateFunctionWithContext(context.Background(), createFunctionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateFunctionWithContext is an alternate form of the CreateFunction method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateFunctionWithContext(ctx context.Context, createFunctionOptions *CreateFunctionOptions) (result *Function, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createFunctionOptions, "createFunctionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createFunctionOptions, "createFunctionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createFunctionOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/functions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createFunctionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateFunction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	body := make(map[string]interface{})
	if createFunctionOptions.CodeReference != nil {
		body["code_reference"] = createFunctionOptions.CodeReference
	}
	if createFunctionOptions.Name != nil {
		body["name"] = createFunctionOptions.Name
	}
	if createFunctionOptions.Runtime != nil {
		body["runtime"] = createFunctionOptions.Runtime
	}
	if createFunctionOptions.CodeBinary != nil {
		body["code_binary"] = createFunctionOptions.CodeBinary
	}
	if createFunctionOptions.CodeMain != nil {
		body["code_main"] = createFunctionOptions.CodeMain
	}
	if createFunctionOptions.CodeSecret != nil {
		body["code_secret"] = createFunctionOptions.CodeSecret
	}
	if createFunctionOptions.ManagedDomainMappings != nil {
		body["managed_domain_mappings"] = createFunctionOptions.ManagedDomainMappings
	}
	if createFunctionOptions.RunEnvVariables != nil {
		body["run_env_variables"] = createFunctionOptions.RunEnvVariables
	}
	if createFunctionOptions.ScaleConcurrency != nil {
		body["scale_concurrency"] = createFunctionOptions.ScaleConcurrency
	}
	if createFunctionOptions.ScaleCpuLimit != nil {
		body["scale_cpu_limit"] = createFunctionOptions.ScaleCpuLimit
	}
	if createFunctionOptions.ScaleDownDelay != nil {
		body["scale_down_delay"] = createFunctionOptions.ScaleDownDelay
	}
	if createFunctionOptions.ScaleMaxExecutionTime != nil {
		body["scale_max_execution_time"] = createFunctionOptions.ScaleMaxExecutionTime
	}
	if createFunctionOptions.ScaleMemoryLimit != nil {
		body["scale_memory_limit"] = createFunctionOptions.ScaleMemoryLimit
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_function", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFunction)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetFunction : Get a function
// Display the details of a function.
func (codeEngine *CodeEngineV2) GetFunction(getFunctionOptions *GetFunctionOptions) (result *Function, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetFunctionWithContext(context.Background(), getFunctionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetFunctionWithContext is an alternate form of the GetFunction method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetFunctionWithContext(ctx context.Context, getFunctionOptions *GetFunctionOptions) (result *Function, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getFunctionOptions, "getFunctionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getFunctionOptions, "getFunctionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getFunctionOptions.ProjectID,
		"name": *getFunctionOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/functions/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getFunctionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetFunction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_function", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFunction)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteFunction : Delete a function
// Delete a function.
func (codeEngine *CodeEngineV2) DeleteFunction(deleteFunctionOptions *DeleteFunctionOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteFunctionWithContext(context.Background(), deleteFunctionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteFunctionWithContext is an alternate form of the DeleteFunction method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteFunctionWithContext(ctx context.Context, deleteFunctionOptions *DeleteFunctionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteFunctionOptions, "deleteFunctionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteFunctionOptions, "deleteFunctionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteFunctionOptions.ProjectID,
		"name": *deleteFunctionOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/functions/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteFunctionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteFunction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_function", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateFunction : Update a function
// Update the given function.
func (codeEngine *CodeEngineV2) UpdateFunction(updateFunctionOptions *UpdateFunctionOptions) (result *Function, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.UpdateFunctionWithContext(context.Background(), updateFunctionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateFunctionWithContext is an alternate form of the UpdateFunction method which supports a Context parameter
func (codeEngine *CodeEngineV2) UpdateFunctionWithContext(ctx context.Context, updateFunctionOptions *UpdateFunctionOptions) (result *Function, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFunctionOptions, "updateFunctionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateFunctionOptions, "updateFunctionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateFunctionOptions.ProjectID,
		"name": *updateFunctionOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/functions/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateFunctionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "UpdateFunction")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateFunctionOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateFunctionOptions.IfMatch))
	}

	if codeEngine.Version != nil {
		builder.AddQuery("version", fmt.Sprint(*codeEngine.Version))
	}

	_, err = builder.SetBodyContentJSON(updateFunctionOptions.Function)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_function", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFunction)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListBindings : List bindings
// List all bindings in a project.
func (codeEngine *CodeEngineV2) ListBindings(listBindingsOptions *ListBindingsOptions) (result *BindingList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListBindingsWithContext(context.Background(), listBindingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListBindingsWithContext is an alternate form of the ListBindings method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListBindingsWithContext(ctx context.Context, listBindingsOptions *ListBindingsOptions) (result *BindingList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listBindingsOptions, "listBindingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listBindingsOptions, "listBindingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listBindingsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/bindings`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listBindingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListBindings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listBindingsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listBindingsOptions.Limit))
	}
	if listBindingsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listBindingsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_bindings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBindingList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateBinding : Create a binding
// Create a binding. Creating a service binding with a Code Engine app will update the app, creating a new revision. For
// more information see the [documentaion](https://cloud.ibm.com/docs/codeengine?topic=codeengine-service-binding).
func (codeEngine *CodeEngineV2) CreateBinding(createBindingOptions *CreateBindingOptions) (result *Binding, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateBindingWithContext(context.Background(), createBindingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateBindingWithContext is an alternate form of the CreateBinding method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateBindingWithContext(ctx context.Context, createBindingOptions *CreateBindingOptions) (result *Binding, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createBindingOptions, "createBindingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createBindingOptions, "createBindingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createBindingOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/bindings`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createBindingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateBinding")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createBindingOptions.Component != nil {
		body["component"] = createBindingOptions.Component
	}
	if createBindingOptions.Prefix != nil {
		body["prefix"] = createBindingOptions.Prefix
	}
	if createBindingOptions.SecretName != nil {
		body["secret_name"] = createBindingOptions.SecretName
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_binding", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBinding)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetBinding : Get a binding
// Display the details of a binding.
func (codeEngine *CodeEngineV2) GetBinding(getBindingOptions *GetBindingOptions) (result *Binding, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetBindingWithContext(context.Background(), getBindingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBindingWithContext is an alternate form of the GetBinding method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetBindingWithContext(ctx context.Context, getBindingOptions *GetBindingOptions) (result *Binding, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBindingOptions, "getBindingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBindingOptions, "getBindingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getBindingOptions.ProjectID,
		"id": *getBindingOptions.ID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/bindings/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getBindingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetBinding")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_binding", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBinding)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteBinding : Delete a binding
// Delete a binding.
func (codeEngine *CodeEngineV2) DeleteBinding(deleteBindingOptions *DeleteBindingOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteBindingWithContext(context.Background(), deleteBindingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteBindingWithContext is an alternate form of the DeleteBinding method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteBindingWithContext(ctx context.Context, deleteBindingOptions *DeleteBindingOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteBindingOptions, "deleteBindingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteBindingOptions, "deleteBindingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteBindingOptions.ProjectID,
		"id": *deleteBindingOptions.ID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/bindings/{id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteBindingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteBinding")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_binding", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListBuilds : List builds
// List all builds in a project.
func (codeEngine *CodeEngineV2) ListBuilds(listBuildsOptions *ListBuildsOptions) (result *BuildList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListBuildsWithContext(context.Background(), listBuildsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListBuildsWithContext is an alternate form of the ListBuilds method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListBuildsWithContext(ctx context.Context, listBuildsOptions *ListBuildsOptions) (result *BuildList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listBuildsOptions, "listBuildsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listBuildsOptions, "listBuildsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listBuildsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/builds`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listBuildsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListBuilds")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listBuildsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listBuildsOptions.Limit))
	}
	if listBuildsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listBuildsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_builds", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuildList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateBuild : Create a build
// Create a build.
func (codeEngine *CodeEngineV2) CreateBuild(createBuildOptions *CreateBuildOptions) (result *Build, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateBuildWithContext(context.Background(), createBuildOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateBuildWithContext is an alternate form of the CreateBuild method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateBuildWithContext(ctx context.Context, createBuildOptions *CreateBuildOptions) (result *Build, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createBuildOptions, "createBuildOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createBuildOptions, "createBuildOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createBuildOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/builds`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createBuildOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateBuild")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createBuildOptions.Name != nil {
		body["name"] = createBuildOptions.Name
	}
	if createBuildOptions.OutputImage != nil {
		body["output_image"] = createBuildOptions.OutputImage
	}
	if createBuildOptions.OutputSecret != nil {
		body["output_secret"] = createBuildOptions.OutputSecret
	}
	if createBuildOptions.StrategyType != nil {
		body["strategy_type"] = createBuildOptions.StrategyType
	}
	if createBuildOptions.SourceContextDir != nil {
		body["source_context_dir"] = createBuildOptions.SourceContextDir
	}
	if createBuildOptions.SourceRevision != nil {
		body["source_revision"] = createBuildOptions.SourceRevision
	}
	if createBuildOptions.SourceSecret != nil {
		body["source_secret"] = createBuildOptions.SourceSecret
	}
	if createBuildOptions.SourceType != nil {
		body["source_type"] = createBuildOptions.SourceType
	}
	if createBuildOptions.SourceURL != nil {
		body["source_url"] = createBuildOptions.SourceURL
	}
	if createBuildOptions.StrategySize != nil {
		body["strategy_size"] = createBuildOptions.StrategySize
	}
	if createBuildOptions.StrategySpecFile != nil {
		body["strategy_spec_file"] = createBuildOptions.StrategySpecFile
	}
	if createBuildOptions.Timeout != nil {
		body["timeout"] = createBuildOptions.Timeout
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_build", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuild)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetBuild : Get a build
// Display the details of a build.
func (codeEngine *CodeEngineV2) GetBuild(getBuildOptions *GetBuildOptions) (result *Build, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetBuildWithContext(context.Background(), getBuildOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBuildWithContext is an alternate form of the GetBuild method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetBuildWithContext(ctx context.Context, getBuildOptions *GetBuildOptions) (result *Build, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBuildOptions, "getBuildOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBuildOptions, "getBuildOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getBuildOptions.ProjectID,
		"name": *getBuildOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/builds/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getBuildOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetBuild")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_build", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuild)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteBuild : Delete a build
// Delete a build.
func (codeEngine *CodeEngineV2) DeleteBuild(deleteBuildOptions *DeleteBuildOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteBuildWithContext(context.Background(), deleteBuildOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteBuildWithContext is an alternate form of the DeleteBuild method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteBuildWithContext(ctx context.Context, deleteBuildOptions *DeleteBuildOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteBuildOptions, "deleteBuildOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteBuildOptions, "deleteBuildOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteBuildOptions.ProjectID,
		"name": *deleteBuildOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/builds/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteBuildOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteBuild")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_build", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateBuild : Update a build
// Update a build.
func (codeEngine *CodeEngineV2) UpdateBuild(updateBuildOptions *UpdateBuildOptions) (result *Build, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.UpdateBuildWithContext(context.Background(), updateBuildOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateBuildWithContext is an alternate form of the UpdateBuild method which supports a Context parameter
func (codeEngine *CodeEngineV2) UpdateBuildWithContext(ctx context.Context, updateBuildOptions *UpdateBuildOptions) (result *Build, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateBuildOptions, "updateBuildOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateBuildOptions, "updateBuildOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateBuildOptions.ProjectID,
		"name": *updateBuildOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/builds/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateBuildOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "UpdateBuild")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateBuildOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateBuildOptions.IfMatch))
	}

	_, err = builder.SetBodyContentJSON(updateBuildOptions.Build)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_build", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuild)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListBuildRuns : List build runs
// List all build runs in a project.
func (codeEngine *CodeEngineV2) ListBuildRuns(listBuildRunsOptions *ListBuildRunsOptions) (result *BuildRunList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListBuildRunsWithContext(context.Background(), listBuildRunsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListBuildRunsWithContext is an alternate form of the ListBuildRuns method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListBuildRunsWithContext(ctx context.Context, listBuildRunsOptions *ListBuildRunsOptions) (result *BuildRunList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listBuildRunsOptions, "listBuildRunsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listBuildRunsOptions, "listBuildRunsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listBuildRunsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/build_runs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listBuildRunsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListBuildRuns")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listBuildRunsOptions.BuildName != nil {
		builder.AddQuery("build_name", fmt.Sprint(*listBuildRunsOptions.BuildName))
	}
	if listBuildRunsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listBuildRunsOptions.Limit))
	}
	if listBuildRunsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listBuildRunsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_build_runs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuildRunList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateBuildRun : Create a build run
// Create a build run.
func (codeEngine *CodeEngineV2) CreateBuildRun(createBuildRunOptions *CreateBuildRunOptions) (result *BuildRun, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateBuildRunWithContext(context.Background(), createBuildRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateBuildRunWithContext is an alternate form of the CreateBuildRun method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateBuildRunWithContext(ctx context.Context, createBuildRunOptions *CreateBuildRunOptions) (result *BuildRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createBuildRunOptions, "createBuildRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createBuildRunOptions, "createBuildRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createBuildRunOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/build_runs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createBuildRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateBuildRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createBuildRunOptions.BuildName != nil {
		body["build_name"] = createBuildRunOptions.BuildName
	}
	if createBuildRunOptions.Name != nil {
		body["name"] = createBuildRunOptions.Name
	}
	if createBuildRunOptions.OutputImage != nil {
		body["output_image"] = createBuildRunOptions.OutputImage
	}
	if createBuildRunOptions.OutputSecret != nil {
		body["output_secret"] = createBuildRunOptions.OutputSecret
	}
	if createBuildRunOptions.ServiceAccount != nil {
		body["service_account"] = createBuildRunOptions.ServiceAccount
	}
	if createBuildRunOptions.SourceContextDir != nil {
		body["source_context_dir"] = createBuildRunOptions.SourceContextDir
	}
	if createBuildRunOptions.SourceRevision != nil {
		body["source_revision"] = createBuildRunOptions.SourceRevision
	}
	if createBuildRunOptions.SourceSecret != nil {
		body["source_secret"] = createBuildRunOptions.SourceSecret
	}
	if createBuildRunOptions.SourceType != nil {
		body["source_type"] = createBuildRunOptions.SourceType
	}
	if createBuildRunOptions.SourceURL != nil {
		body["source_url"] = createBuildRunOptions.SourceURL
	}
	if createBuildRunOptions.StrategySize != nil {
		body["strategy_size"] = createBuildRunOptions.StrategySize
	}
	if createBuildRunOptions.StrategySpecFile != nil {
		body["strategy_spec_file"] = createBuildRunOptions.StrategySpecFile
	}
	if createBuildRunOptions.StrategyType != nil {
		body["strategy_type"] = createBuildRunOptions.StrategyType
	}
	if createBuildRunOptions.Timeout != nil {
		body["timeout"] = createBuildRunOptions.Timeout
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_build_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuildRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetBuildRun : Get a build run
// Display the details of a build run.
func (codeEngine *CodeEngineV2) GetBuildRun(getBuildRunOptions *GetBuildRunOptions) (result *BuildRun, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetBuildRunWithContext(context.Background(), getBuildRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetBuildRunWithContext is an alternate form of the GetBuildRun method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetBuildRunWithContext(ctx context.Context, getBuildRunOptions *GetBuildRunOptions) (result *BuildRun, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getBuildRunOptions, "getBuildRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getBuildRunOptions, "getBuildRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getBuildRunOptions.ProjectID,
		"name": *getBuildRunOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/build_runs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getBuildRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetBuildRun")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_build_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalBuildRun)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteBuildRun : Delete a build run
// Delete a build run.
func (codeEngine *CodeEngineV2) DeleteBuildRun(deleteBuildRunOptions *DeleteBuildRunOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteBuildRunWithContext(context.Background(), deleteBuildRunOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteBuildRunWithContext is an alternate form of the DeleteBuildRun method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteBuildRunWithContext(ctx context.Context, deleteBuildRunOptions *DeleteBuildRunOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteBuildRunOptions, "deleteBuildRunOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteBuildRunOptions, "deleteBuildRunOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteBuildRunOptions.ProjectID,
		"name": *deleteBuildRunOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/build_runs/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteBuildRunOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteBuildRun")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_build_run", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListDomainMappings : List domain mappings
// List all domain mappings in a project.
func (codeEngine *CodeEngineV2) ListDomainMappings(listDomainMappingsOptions *ListDomainMappingsOptions) (result *DomainMappingList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListDomainMappingsWithContext(context.Background(), listDomainMappingsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListDomainMappingsWithContext is an alternate form of the ListDomainMappings method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListDomainMappingsWithContext(ctx context.Context, listDomainMappingsOptions *ListDomainMappingsOptions) (result *DomainMappingList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listDomainMappingsOptions, "listDomainMappingsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listDomainMappingsOptions, "listDomainMappingsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listDomainMappingsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/domain_mappings`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listDomainMappingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListDomainMappings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listDomainMappingsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listDomainMappingsOptions.Limit))
	}
	if listDomainMappingsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listDomainMappingsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_domain_mappings", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDomainMappingList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateDomainMapping : Create a domain mapping
// Create a domain mapping.
func (codeEngine *CodeEngineV2) CreateDomainMapping(createDomainMappingOptions *CreateDomainMappingOptions) (result *DomainMapping, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateDomainMappingWithContext(context.Background(), createDomainMappingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateDomainMappingWithContext is an alternate form of the CreateDomainMapping method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateDomainMappingWithContext(ctx context.Context, createDomainMappingOptions *CreateDomainMappingOptions) (result *DomainMapping, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createDomainMappingOptions, "createDomainMappingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createDomainMappingOptions, "createDomainMappingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createDomainMappingOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/domain_mappings`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createDomainMappingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateDomainMapping")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createDomainMappingOptions.Component != nil {
		body["component"] = createDomainMappingOptions.Component
	}
	if createDomainMappingOptions.Name != nil {
		body["name"] = createDomainMappingOptions.Name
	}
	if createDomainMappingOptions.TlsSecret != nil {
		body["tls_secret"] = createDomainMappingOptions.TlsSecret
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_domain_mapping", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDomainMapping)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetDomainMapping : Get a domain mapping
// Get domain mapping.
func (codeEngine *CodeEngineV2) GetDomainMapping(getDomainMappingOptions *GetDomainMappingOptions) (result *DomainMapping, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetDomainMappingWithContext(context.Background(), getDomainMappingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetDomainMappingWithContext is an alternate form of the GetDomainMapping method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetDomainMappingWithContext(ctx context.Context, getDomainMappingOptions *GetDomainMappingOptions) (result *DomainMapping, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getDomainMappingOptions, "getDomainMappingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getDomainMappingOptions, "getDomainMappingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getDomainMappingOptions.ProjectID,
		"name": *getDomainMappingOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/domain_mappings/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getDomainMappingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetDomainMapping")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_domain_mapping", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDomainMapping)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteDomainMapping : Delete a domain mapping
// Delete a domain mapping.
func (codeEngine *CodeEngineV2) DeleteDomainMapping(deleteDomainMappingOptions *DeleteDomainMappingOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteDomainMappingWithContext(context.Background(), deleteDomainMappingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteDomainMappingWithContext is an alternate form of the DeleteDomainMapping method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteDomainMappingWithContext(ctx context.Context, deleteDomainMappingOptions *DeleteDomainMappingOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteDomainMappingOptions, "deleteDomainMappingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteDomainMappingOptions, "deleteDomainMappingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteDomainMappingOptions.ProjectID,
		"name": *deleteDomainMappingOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/domain_mappings/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteDomainMappingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteDomainMapping")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_domain_mapping", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateDomainMapping : Update a domain mapping
// Update a domain mapping.
func (codeEngine *CodeEngineV2) UpdateDomainMapping(updateDomainMappingOptions *UpdateDomainMappingOptions) (result *DomainMapping, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.UpdateDomainMappingWithContext(context.Background(), updateDomainMappingOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateDomainMappingWithContext is an alternate form of the UpdateDomainMapping method which supports a Context parameter
func (codeEngine *CodeEngineV2) UpdateDomainMappingWithContext(ctx context.Context, updateDomainMappingOptions *UpdateDomainMappingOptions) (result *DomainMapping, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateDomainMappingOptions, "updateDomainMappingOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateDomainMappingOptions, "updateDomainMappingOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *updateDomainMappingOptions.ProjectID,
		"name": *updateDomainMappingOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/domain_mappings/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateDomainMappingOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "UpdateDomainMapping")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")
	if updateDomainMappingOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*updateDomainMappingOptions.IfMatch))
	}

	_, err = builder.SetBodyContentJSON(updateDomainMappingOptions.DomainMapping)
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_domain_mapping", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDomainMapping)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListConfigMaps : List config maps
// List all config maps in a project.
func (codeEngine *CodeEngineV2) ListConfigMaps(listConfigMapsOptions *ListConfigMapsOptions) (result *ConfigMapList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListConfigMapsWithContext(context.Background(), listConfigMapsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListConfigMapsWithContext is an alternate form of the ListConfigMaps method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListConfigMapsWithContext(ctx context.Context, listConfigMapsOptions *ListConfigMapsOptions) (result *ConfigMapList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listConfigMapsOptions, "listConfigMapsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listConfigMapsOptions, "listConfigMapsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listConfigMapsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/config_maps`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listConfigMapsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListConfigMaps")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listConfigMapsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listConfigMapsOptions.Limit))
	}
	if listConfigMapsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listConfigMapsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_config_maps", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigMapList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateConfigMap : Create a config map
// Create a config map.
func (codeEngine *CodeEngineV2) CreateConfigMap(createConfigMapOptions *CreateConfigMapOptions) (result *ConfigMap, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateConfigMapWithContext(context.Background(), createConfigMapOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateConfigMapWithContext is an alternate form of the CreateConfigMap method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateConfigMapWithContext(ctx context.Context, createConfigMapOptions *CreateConfigMapOptions) (result *ConfigMap, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createConfigMapOptions, "createConfigMapOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createConfigMapOptions, "createConfigMapOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createConfigMapOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/config_maps`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createConfigMapOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateConfigMap")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createConfigMapOptions.Name != nil {
		body["name"] = createConfigMapOptions.Name
	}
	if createConfigMapOptions.Data != nil {
		body["data"] = createConfigMapOptions.Data
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_config_map", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigMap)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetConfigMap : Get a config map
// Display the details of a config map.
func (codeEngine *CodeEngineV2) GetConfigMap(getConfigMapOptions *GetConfigMapOptions) (result *ConfigMap, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetConfigMapWithContext(context.Background(), getConfigMapOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetConfigMapWithContext is an alternate form of the GetConfigMap method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetConfigMapWithContext(ctx context.Context, getConfigMapOptions *GetConfigMapOptions) (result *ConfigMap, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getConfigMapOptions, "getConfigMapOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getConfigMapOptions, "getConfigMapOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getConfigMapOptions.ProjectID,
		"name": *getConfigMapOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/config_maps/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getConfigMapOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetConfigMap")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_config_map", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigMap)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceConfigMap : Update a config map
// Update a config map.
func (codeEngine *CodeEngineV2) ReplaceConfigMap(replaceConfigMapOptions *ReplaceConfigMapOptions) (result *ConfigMap, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ReplaceConfigMapWithContext(context.Background(), replaceConfigMapOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceConfigMapWithContext is an alternate form of the ReplaceConfigMap method which supports a Context parameter
func (codeEngine *CodeEngineV2) ReplaceConfigMapWithContext(ctx context.Context, replaceConfigMapOptions *ReplaceConfigMapOptions) (result *ConfigMap, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceConfigMapOptions, "replaceConfigMapOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceConfigMapOptions, "replaceConfigMapOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *replaceConfigMapOptions.ProjectID,
		"name": *replaceConfigMapOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/config_maps/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceConfigMapOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ReplaceConfigMap")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceConfigMapOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceConfigMapOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replaceConfigMapOptions.Data != nil {
		body["data"] = replaceConfigMapOptions.Data
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_config_map", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigMap)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteConfigMap : Delete a config map
// Delete a config map.
func (codeEngine *CodeEngineV2) DeleteConfigMap(deleteConfigMapOptions *DeleteConfigMapOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteConfigMapWithContext(context.Background(), deleteConfigMapOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteConfigMapWithContext is an alternate form of the DeleteConfigMap method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteConfigMapWithContext(ctx context.Context, deleteConfigMapOptions *DeleteConfigMapOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteConfigMapOptions, "deleteConfigMapOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteConfigMapOptions, "deleteConfigMapOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteConfigMapOptions.ProjectID,
		"name": *deleteConfigMapOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/config_maps/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteConfigMapOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteConfigMap")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_config_map", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListSecrets : List secrets
// List all secrets in a project.
func (codeEngine *CodeEngineV2) ListSecrets(listSecretsOptions *ListSecretsOptions) (result *SecretList, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ListSecretsWithContext(context.Background(), listSecretsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListSecretsWithContext is an alternate form of the ListSecrets method which supports a Context parameter
func (codeEngine *CodeEngineV2) ListSecretsWithContext(ctx context.Context, listSecretsOptions *ListSecretsOptions) (result *SecretList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSecretsOptions, "listSecretsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listSecretsOptions, "listSecretsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *listSecretsOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/secrets`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listSecretsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ListSecrets")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSecretsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSecretsOptions.Limit))
	}
	if listSecretsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listSecretsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_secrets", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecretList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateSecret : Create a secret
// Create a secret.
func (codeEngine *CodeEngineV2) CreateSecret(createSecretOptions *CreateSecretOptions) (result *Secret, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.CreateSecretWithContext(context.Background(), createSecretOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateSecretWithContext is an alternate form of the CreateSecret method which supports a Context parameter
func (codeEngine *CodeEngineV2) CreateSecretWithContext(ctx context.Context, createSecretOptions *CreateSecretOptions) (result *Secret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSecretOptions, "createSecretOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createSecretOptions, "createSecretOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *createSecretOptions.ProjectID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/secrets`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "CreateSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSecretOptions.Format != nil {
		body["format"] = createSecretOptions.Format
	}
	if createSecretOptions.Name != nil {
		body["name"] = createSecretOptions.Name
	}
	if createSecretOptions.Data != nil {
		body["data"] = createSecretOptions.Data
	}
	if createSecretOptions.ServiceAccess != nil {
		body["service_access"] = createSecretOptions.ServiceAccess
	}
	if createSecretOptions.ServiceOperator != nil {
		body["service_operator"] = createSecretOptions.ServiceOperator
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_secret", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecret)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetSecret : Get a secret
// Get a secret.
func (codeEngine *CodeEngineV2) GetSecret(getSecretOptions *GetSecretOptions) (result *Secret, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.GetSecretWithContext(context.Background(), getSecretOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSecretWithContext is an alternate form of the GetSecret method which supports a Context parameter
func (codeEngine *CodeEngineV2) GetSecretWithContext(ctx context.Context, getSecretOptions *GetSecretOptions) (result *Secret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSecretOptions, "getSecretOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getSecretOptions, "getSecretOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *getSecretOptions.ProjectID,
		"name": *getSecretOptions.Name,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/secrets/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "GetSecret")
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_secret", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecret)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ReplaceSecret : Update a secret
// Update a secret.
func (codeEngine *CodeEngineV2) ReplaceSecret(replaceSecretOptions *ReplaceSecretOptions) (result *Secret, response *core.DetailedResponse, err error) {
	result, response, err = codeEngine.ReplaceSecretWithContext(context.Background(), replaceSecretOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ReplaceSecretWithContext is an alternate form of the ReplaceSecret method which supports a Context parameter
func (codeEngine *CodeEngineV2) ReplaceSecretWithContext(ctx context.Context, replaceSecretOptions *ReplaceSecretOptions) (result *Secret, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(replaceSecretOptions, "replaceSecretOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(replaceSecretOptions, "replaceSecretOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *replaceSecretOptions.ProjectID,
		"name": *replaceSecretOptions.Name,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/secrets/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range replaceSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "ReplaceSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if replaceSecretOptions.IfMatch != nil {
		builder.AddHeader("If-Match", fmt.Sprint(*replaceSecretOptions.IfMatch))
	}

	body := make(map[string]interface{})
	if replaceSecretOptions.Format != nil {
		body["format"] = replaceSecretOptions.Format
	}
	if replaceSecretOptions.Data != nil {
		body["data"] = replaceSecretOptions.Data
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
	response, err = codeEngine.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "replace_secret", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSecret)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteSecret : Delete a secret
// Delete a secret.
func (codeEngine *CodeEngineV2) DeleteSecret(deleteSecretOptions *DeleteSecretOptions) (response *core.DetailedResponse, err error) {
	response, err = codeEngine.DeleteSecretWithContext(context.Background(), deleteSecretOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteSecretWithContext is an alternate form of the DeleteSecret method which supports a Context parameter
func (codeEngine *CodeEngineV2) DeleteSecretWithContext(ctx context.Context, deleteSecretOptions *DeleteSecretOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSecretOptions, "deleteSecretOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteSecretOptions, "deleteSecretOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"project_id": *deleteSecretOptions.ProjectID,
		"name": *deleteSecretOptions.Name,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = codeEngine.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(codeEngine.Service.Options.URL, `/projects/{project_id}/secrets/{name}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteSecretOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("code_engine", "V2", "DeleteSecret")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = codeEngine.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_secret", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "2.0.0")
}

// AllowedOutboundDestination : AllowedOutboundDestination Describes the model of an allowed outbound destination.
// Models which "extend" this model:
// - AllowedOutboundDestinationCidrBlockData
type AllowedOutboundDestination struct {
	// The version of the allowed outbound destination, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
	Type *string `json:"type" validate:"required"`

	// The IPv4 address range.
	CidrBlock *string `json:"cidr_block,omitempty"`

	// The name of the CIDR block.
	Name *string `json:"name,omitempty"`
}

// Constants associated with the AllowedOutboundDestination.Type property.
// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
const (
	AllowedOutboundDestination_Type_CidrBlock = "cidr_block"
)
func (*AllowedOutboundDestination) isaAllowedOutboundDestination() bool {
	return true
}

type AllowedOutboundDestinationIntf interface {
	isaAllowedOutboundDestination() bool
}

// UnmarshalAllowedOutboundDestination unmarshals an instance of AllowedOutboundDestination from the specified map of raw messages.
func UnmarshalAllowedOutboundDestination(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestination)
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cidr_block", &obj.CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "cidr_block-error", common.GetComponentInfo())
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

// AllowedOutboundDestinationList : Contains a list of allowed outbound destinations and pagination information.
type AllowedOutboundDestinationList struct {
	// List of all allowed outbound destinations.
	AllowedOutboundDestinations []AllowedOutboundDestinationIntf `json:"allowed_outbound_destinations" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalAllowedOutboundDestinationList unmarshals an instance of AllowedOutboundDestinationList from the specified map of raw messages.
func UnmarshalAllowedOutboundDestinationList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestinationList)
	err = core.UnmarshalModel(m, "allowed_outbound_destinations", &obj.AllowedOutboundDestinations, UnmarshalAllowedOutboundDestination)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_outbound_destinations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AllowedOutboundDestinationList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// AllowedOutboundDestinationPatch : AllowedOutboundDestinationPatch is the request model for allowed outbound destination update operations.
// Models which "extend" this model:
// - AllowedOutboundDestinationPatchCidrBlockDataPatch
type AllowedOutboundDestinationPatch struct {
	// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
	Type *string `json:"type,omitempty"`

	// The IPv4 address range.
	CidrBlock *string `json:"cidr_block,omitempty"`
}

// Constants associated with the AllowedOutboundDestinationPatch.Type property.
// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
const (
	AllowedOutboundDestinationPatch_Type_CidrBlock = "cidr_block"
)
func (*AllowedOutboundDestinationPatch) isaAllowedOutboundDestinationPatch() bool {
	return true
}

type AllowedOutboundDestinationPatchIntf interface {
	isaAllowedOutboundDestinationPatch() bool
}

// UnmarshalAllowedOutboundDestinationPatch unmarshals an instance of AllowedOutboundDestinationPatch from the specified map of raw messages.
func UnmarshalAllowedOutboundDestinationPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestinationPatch)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cidr_block", &obj.CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "cidr_block-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the AllowedOutboundDestinationPatch
func (allowedOutboundDestinationPatch *AllowedOutboundDestinationPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(allowedOutboundDestinationPatch.Type) {
		_patch["type"] = allowedOutboundDestinationPatch.Type
	}
	if !core.IsNil(allowedOutboundDestinationPatch.CidrBlock) {
		_patch["cidr_block"] = allowedOutboundDestinationPatch.CidrBlock
	}

	return
}

// AllowedOutboundDestinationPrototype : AllowedOutboundDestinationPrototype is the request model for allowed outbound destination create operations.
// Models which "extend" this model:
// - AllowedOutboundDestinationPrototypeCidrBlockDataPrototype
type AllowedOutboundDestinationPrototype struct {
	// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
	Type *string `json:"type" validate:"required"`

	// The IPv4 address range.
	CidrBlock *string `json:"cidr_block,omitempty"`

	// The name of the CIDR block.
	Name *string `json:"name,omitempty"`
}

// Constants associated with the AllowedOutboundDestinationPrototype.Type property.
// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
const (
	AllowedOutboundDestinationPrototype_Type_CidrBlock = "cidr_block"
)
func (*AllowedOutboundDestinationPrototype) isaAllowedOutboundDestinationPrototype() bool {
	return true
}

type AllowedOutboundDestinationPrototypeIntf interface {
	isaAllowedOutboundDestinationPrototype() bool
}

// UnmarshalAllowedOutboundDestinationPrototype unmarshals an instance of AllowedOutboundDestinationPrototype from the specified map of raw messages.
func UnmarshalAllowedOutboundDestinationPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestinationPrototype)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cidr_block", &obj.CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "cidr_block-error", common.GetComponentInfo())
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

// App : App is the response model for app resources.
type App struct {
	// Reference to a build that is associated with the application.
	Build *string `json:"build,omitempty"`

	// Reference to a build run that is associated with the application.
	BuildRun *string `json:"build_run,omitempty"`

	// References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as
	// environment variables in the application.
	ComputedEnvVariables []EnvVar `json:"computed_env_variables,omitempty"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// Optional URL to invoke the app. Depending on visibility,  this is accessible publicly or in the private network
	// only. Empty in case 'managed_domain_mappings' is set to 'local'.
	Endpoint *string `json:"endpoint,omitempty"`

	// The URL to the app that is only visible within the project.
	EndpointInternal *string `json:"endpoint_internal,omitempty"`

	// The version of the app instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// When you provision a new app,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is
	// used to connect to the port that is exposed by the container image.
	ImagePort *int64 `json:"image_port,omitempty"`

	// The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference" validate:"required"`

	// Optional name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the app will be created but cannot reach the ready status, until this property is provided, too.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Optional value controlling which of the system managed domain mappings will be setup for the application. Valid
	// values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project
	// supports application private visibility.
	ManagedDomainMappings *string `json:"managed_domain_mappings" validate:"required"`

	// The name of the app.
	Name *string `json:"name" validate:"required"`

	// Response model for probes.
	ProbeLiveness *Probe `json:"probe_liveness,omitempty"`

	// Response model for probes.
	ProbeReadiness *Probe `json:"probe_readiness,omitempty"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the app.
	ResourceType *string `json:"resource_type,omitempty"`

	// Optional arguments for the app that are passed to start the container. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments" validate:"required"`

	// Optional user ID (UID) to run the app.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Optional commands for the app that are passed to start the container. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands" validate:"required"`

	// References to config maps, secrets or literal values, which are defined by the app owner and are exposed as
	// environment variables in the application.
	RunEnvVariables []EnvVar `json:"run_env_variables" validate:"required"`

	// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
	// `none`, `reader`, and `writer`.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Mounts of config maps or secrets.
	RunVolumeMounts []VolumeMount `json:"run_volume_mounts" validate:"required"`

	// Optional maximum number of requests that can be processed concurrently per instance.
	ScaleConcurrency *int64 `json:"scale_concurrency,omitempty"`

	// Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use
	// this value to scale up instances based on concurrent number of requests. This option defaults to the value of the
	// `scale_concurrency` option, if not specified.
	ScaleConcurrencyTarget *int64 `json:"scale_concurrency_target,omitempty"`

	// Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit" validate:"required"`

	// Optional amount of time in seconds that delays the scale-down behavior for an app instance.
	ScaleDownDelay *int64 `json:"scale_down_delay,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit" validate:"required"`

	// Optional initial number of instances that are created upon app creation or app update.
	ScaleInitialInstances *int64 `json:"scale_initial_instances,omitempty"`

	// Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper
	// scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas
	// for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).
	ScaleMaxInstances *int64 `json:"scale_max_instances" validate:"required"`

	// Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit" validate:"required"`

	// Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if
	// not hit by any request for some time.
	ScaleMinInstances *int64 `json:"scale_min_instances" validate:"required"`

	// Optional amount of time in seconds that is allowed for a running app to respond to a request.
	ScaleRequestTimeout *int64 `json:"scale_request_timeout" validate:"required"`

	// The current status of the app.
	Status *string `json:"status,omitempty"`

	// The detailed status of the application.
	StatusDetails *AppStatus `json:"status_details,omitempty"`
}

// Constants associated with the App.ManagedDomainMappings property.
// Optional value controlling which of the system managed domain mappings will be setup for the application. Valid
// values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project
// supports application private visibility.
const (
	App_ManagedDomainMappings_Local = "local"
	App_ManagedDomainMappings_LocalPrivate = "local_private"
	App_ManagedDomainMappings_LocalPublic = "local_public"
)

// Constants associated with the App.ResourceType property.
// The type of the app.
const (
	App_ResourceType_AppV2 = "app_v2"
)

// Constants associated with the App.RunServiceAccount property.
// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
// `none`, `reader`, and `writer`.
const (
	App_RunServiceAccount_Default = "default"
	App_RunServiceAccount_Manager = "manager"
	App_RunServiceAccount_None = "none"
	App_RunServiceAccount_Reader = "reader"
	App_RunServiceAccount_Writer = "writer"
)

// Constants associated with the App.Status property.
// The current status of the app.
const (
	App_Status_Deploying = "deploying"
	App_Status_Failed = "failed"
	App_Status_Ready = "ready"
	App_Status_Warning = "warning"
)

// UnmarshalApp unmarshals an instance of App from the specified map of raw messages.
func UnmarshalApp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(App)
	err = core.UnmarshalPrimitive(m, "build", &obj.Build)
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "build_run", &obj.BuildRun)
	if err != nil {
		err = core.SDKErrorf(err, "", "build_run-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "computed_env_variables", &obj.ComputedEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "computed_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoint", &obj.Endpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoint_internal", &obj.EndpointInternal)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint_internal-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "image_port", &obj.ImagePort)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_reference", &obj.ImageReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_secret", &obj.ImageSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "managed_domain_mappings", &obj.ManagedDomainMappings)
	if err != nil {
		err = core.SDKErrorf(err, "", "managed_domain_mappings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "probe_liveness", &obj.ProbeLiveness, UnmarshalProbe)
	if err != nil {
		err = core.SDKErrorf(err, "", "probe_liveness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "probe_readiness", &obj.ProbeReadiness, UnmarshalProbe)
	if err != nil {
		err = core.SDKErrorf(err, "", "probe_readiness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_arguments", &obj.RunArguments)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_arguments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_as_user", &obj.RunAsUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_as_user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_commands", &obj.RunCommands)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_commands-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_service_account", &obj.RunServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_volume_mounts", &obj.RunVolumeMounts, UnmarshalVolumeMount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_volume_mounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency", &obj.ScaleConcurrency)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency_target", &obj.ScaleConcurrencyTarget)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency_target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_down_delay", &obj.ScaleDownDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_down_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_initial_instances", &obj.ScaleInitialInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_initial_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_instances", &obj.ScaleMaxInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_min_instances", &obj.ScaleMinInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_min_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_request_timeout", &obj.ScaleRequestTimeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_request_timeout-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalAppStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AppInstance : AppInstance is the response model for app instance resources.
type AppInstance struct {
	// The name of the application that is associated with this instance.
	AppName *string `json:"app_name" validate:"required"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// When you provision a new app instance, a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the app instance.
	Name *string `json:"name,omitempty"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the app instance.
	ResourceType *string `json:"resource_type,omitempty"`

	// The number of restarts of the app instance.
	Restarts *int64 `json:"restarts,omitempty"`

	// The name of the revision that is associated with this instance.
	RevisionName *string `json:"revision_name" validate:"required"`

	// The number of CPU set for the instance. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit" validate:"required"`

	// The amount of ephemeral storage set for the instance. The amount specified as ephemeral storage, must not exceed the
	// amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or Gigabyte (G), whereas
	// G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit" validate:"required"`

	// The amount of memory set for the instance. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit" validate:"required"`

	// The current status of the instance.
	Status *string `json:"status,omitempty"`

	// The status of a container.
	SystemContainer *ContainerStatus `json:"system_container,omitempty"`

	// The status of a container.
	UserContainer *ContainerStatus `json:"user_container,omitempty"`
}

// Constants associated with the AppInstance.ResourceType property.
// The type of the app instance.
const (
	AppInstance_ResourceType_AppInstanceV2 = "app_instance_v2"
)

// Constants associated with the AppInstance.Status property.
// The current status of the instance.
const (
	AppInstance_Status_Failed = "failed"
	AppInstance_Status_Pending = "pending"
	AppInstance_Status_Running = "running"
	AppInstance_Status_Succeeded = "succeeded"
)

// UnmarshalAppInstance unmarshals an instance of AppInstance from the specified map of raw messages.
func UnmarshalAppInstance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppInstance)
	err = core.UnmarshalPrimitive(m, "app_name", &obj.AppName)
	if err != nil {
		err = core.SDKErrorf(err, "", "app_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "restarts", &obj.Restarts)
	if err != nil {
		err = core.SDKErrorf(err, "", "restarts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "revision_name", &obj.RevisionName)
	if err != nil {
		err = core.SDKErrorf(err, "", "revision_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "system_container", &obj.SystemContainer, UnmarshalContainerStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "system_container-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "user_container", &obj.UserContainer, UnmarshalContainerStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_container-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AppInstanceList : Contains a list of app instances and pagination information.
type AppInstanceList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// List of all app instances.
	Instances []AppInstance `json:"instances" validate:"required"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalAppInstanceList unmarshals an instance of AppInstanceList from the specified map of raw messages.
func UnmarshalAppInstanceList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppInstanceList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "instances", &obj.Instances, UnmarshalAppInstance)
	if err != nil {
		err = core.SDKErrorf(err, "", "instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AppInstanceList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// AppList : Contains a list of apps and pagination information.
type AppList struct {
	// List of all apps.
	Apps []App `json:"apps" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalAppList unmarshals an instance of AppList from the specified map of raw messages.
func UnmarshalAppList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppList)
	err = core.UnmarshalModel(m, "apps", &obj.Apps, UnmarshalApp)
	if err != nil {
		err = core.SDKErrorf(err, "", "apps-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AppList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// AppPatch : App is the request model for app update operations.
type AppPatch struct {
	// Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is
	// used to connect to the port that is exposed by the container image.
	ImagePort *int64 `json:"image_port,omitempty"`

	// The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference,omitempty"`

	// Optional name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the app will be created but cannot reach the ready status, until this property is provided, too.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Optional value controlling which of the system managed domain mappings will be setup for the application. Valid
	// values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project
	// supports application private visibility.
	ManagedDomainMappings *string `json:"managed_domain_mappings,omitempty"`

	// Request model for probes.
	ProbeLiveness *ProbePrototype `json:"probe_liveness,omitempty"`

	// Request model for probes.
	ProbeReadiness *ProbePrototype `json:"probe_readiness,omitempty"`

	// Optional arguments for the app that are passed to start the container. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments,omitempty"`

	// Optional user ID (UID) to run the app.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Optional commands for the app that are passed to start the container. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands,omitempty"`

	// Optional references to config maps, secrets or literal values.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
	// `none`, `reader`, and `writer`.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or a secrets. In case this is provided, existing `run_volume_mounts` will be
	// overwritten.
	RunVolumeMounts []VolumeMountPrototype `json:"run_volume_mounts,omitempty"`

	// Optional maximum number of requests that can be processed concurrently per instance.
	ScaleConcurrency *int64 `json:"scale_concurrency,omitempty"`

	// Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use
	// this value to scale up instances based on concurrent number of requests. This option defaults to the value of the
	// `scale_concurrency` option, if not specified.
	ScaleConcurrencyTarget *int64 `json:"scale_concurrency_target,omitempty"`

	// Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of time in seconds that delays the scale-down behavior for an app instance.
	ScaleDownDelay *int64 `json:"scale_down_delay,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit,omitempty"`

	// Optional initial number of instances that are created upon app creation or app update.
	ScaleInitialInstances *int64 `json:"scale_initial_instances,omitempty"`

	// Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper
	// scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas
	// for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).
	ScaleMaxInstances *int64 `json:"scale_max_instances,omitempty"`

	// Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if
	// not hit by any request for some time.
	ScaleMinInstances *int64 `json:"scale_min_instances,omitempty"`

	// Optional amount of time in seconds that is allowed for a running app to respond to a request.
	ScaleRequestTimeout *int64 `json:"scale_request_timeout,omitempty"`
}

// Constants associated with the AppPatch.ManagedDomainMappings property.
// Optional value controlling which of the system managed domain mappings will be setup for the application. Valid
// values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project
// supports application private visibility.
const (
	AppPatch_ManagedDomainMappings_Local = "local"
	AppPatch_ManagedDomainMappings_LocalPrivate = "local_private"
	AppPatch_ManagedDomainMappings_LocalPublic = "local_public"
)

// Constants associated with the AppPatch.RunServiceAccount property.
// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
// `none`, `reader`, and `writer`.
const (
	AppPatch_RunServiceAccount_Default = "default"
	AppPatch_RunServiceAccount_Manager = "manager"
	AppPatch_RunServiceAccount_None = "none"
	AppPatch_RunServiceAccount_Reader = "reader"
	AppPatch_RunServiceAccount_Writer = "writer"
)

// UnmarshalAppPatch unmarshals an instance of AppPatch from the specified map of raw messages.
func UnmarshalAppPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppPatch)
	err = core.UnmarshalPrimitive(m, "image_port", &obj.ImagePort)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_reference", &obj.ImageReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_secret", &obj.ImageSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "managed_domain_mappings", &obj.ManagedDomainMappings)
	if err != nil {
		err = core.SDKErrorf(err, "", "managed_domain_mappings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "probe_liveness", &obj.ProbeLiveness, UnmarshalProbePrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "probe_liveness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "probe_readiness", &obj.ProbeReadiness, UnmarshalProbePrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "probe_readiness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_arguments", &obj.RunArguments)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_arguments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_as_user", &obj.RunAsUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_as_user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_commands", &obj.RunCommands)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_commands-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVarPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_service_account", &obj.RunServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_volume_mounts", &obj.RunVolumeMounts, UnmarshalVolumeMountPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_volume_mounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency", &obj.ScaleConcurrency)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency_target", &obj.ScaleConcurrencyTarget)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency_target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_down_delay", &obj.ScaleDownDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_down_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_initial_instances", &obj.ScaleInitialInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_initial_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_instances", &obj.ScaleMaxInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_min_instances", &obj.ScaleMinInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_min_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_request_timeout", &obj.ScaleRequestTimeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_request_timeout-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the AppPatch
func (appPatch *AppPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(appPatch.ImagePort) {
		_patch["image_port"] = appPatch.ImagePort
	}
	if !core.IsNil(appPatch.ImageReference) {
		_patch["image_reference"] = appPatch.ImageReference
	}
	if !core.IsNil(appPatch.ImageSecret) {
		_patch["image_secret"] = appPatch.ImageSecret
	}
	if !core.IsNil(appPatch.ManagedDomainMappings) {
		_patch["managed_domain_mappings"] = appPatch.ManagedDomainMappings
	}
	if !core.IsNil(appPatch.ProbeLiveness) {
		_patch["probe_liveness"] = appPatch.ProbeLiveness.asPatch()
	}
	if !core.IsNil(appPatch.ProbeReadiness) {
		_patch["probe_readiness"] = appPatch.ProbeReadiness.asPatch()
	}
	if !core.IsNil(appPatch.RunArguments) {
		_patch["run_arguments"] = appPatch.RunArguments
	}
	if !core.IsNil(appPatch.RunAsUser) {
		_patch["run_as_user"] = appPatch.RunAsUser
	}
	if !core.IsNil(appPatch.RunCommands) {
		_patch["run_commands"] = appPatch.RunCommands
	}
	if !core.IsNil(appPatch.RunEnvVariables) {
		var runEnvVariablesPatches []map[string]interface{}
		for _, runEnvVariables := range appPatch.RunEnvVariables {
			runEnvVariablesPatches = append(runEnvVariablesPatches, runEnvVariables.asPatch())
		}
		_patch["run_env_variables"] = runEnvVariablesPatches
	}
	if !core.IsNil(appPatch.RunServiceAccount) {
		_patch["run_service_account"] = appPatch.RunServiceAccount
	}
	if !core.IsNil(appPatch.RunVolumeMounts) {
		var runVolumeMountsPatches []map[string]interface{}
		for _, runVolumeMounts := range appPatch.RunVolumeMounts {
			runVolumeMountsPatches = append(runVolumeMountsPatches, runVolumeMounts.asPatch())
		}
		_patch["run_volume_mounts"] = runVolumeMountsPatches
	}
	if !core.IsNil(appPatch.ScaleConcurrency) {
		_patch["scale_concurrency"] = appPatch.ScaleConcurrency
	}
	if !core.IsNil(appPatch.ScaleConcurrencyTarget) {
		_patch["scale_concurrency_target"] = appPatch.ScaleConcurrencyTarget
	}
	if !core.IsNil(appPatch.ScaleCpuLimit) {
		_patch["scale_cpu_limit"] = appPatch.ScaleCpuLimit
	}
	if !core.IsNil(appPatch.ScaleDownDelay) {
		_patch["scale_down_delay"] = appPatch.ScaleDownDelay
	}
	if !core.IsNil(appPatch.ScaleEphemeralStorageLimit) {
		_patch["scale_ephemeral_storage_limit"] = appPatch.ScaleEphemeralStorageLimit
	}
	if !core.IsNil(appPatch.ScaleInitialInstances) {
		_patch["scale_initial_instances"] = appPatch.ScaleInitialInstances
	}
	if !core.IsNil(appPatch.ScaleMaxInstances) {
		_patch["scale_max_instances"] = appPatch.ScaleMaxInstances
	}
	if !core.IsNil(appPatch.ScaleMemoryLimit) {
		_patch["scale_memory_limit"] = appPatch.ScaleMemoryLimit
	}
	if !core.IsNil(appPatch.ScaleMinInstances) {
		_patch["scale_min_instances"] = appPatch.ScaleMinInstances
	}
	if !core.IsNil(appPatch.ScaleRequestTimeout) {
		_patch["scale_request_timeout"] = appPatch.ScaleRequestTimeout
	}

	return
}

// AppRevision : AppRevision is the response model for app revision resources.
type AppRevision struct {
	// Name of the associated app.
	AppName *string `json:"app_name,omitempty"`

	// References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as
	// environment variables in the application.
	ComputedEnvVariables []EnvVar `json:"computed_env_variables,omitempty"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// When you provision a new revision,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is
	// used to connect to the port that is exposed by the container image.
	ImagePort *int64 `json:"image_port,omitempty"`

	// The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference" validate:"required"`

	// Optional name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the app will be created but cannot reach the ready status, until this property is provided, too.
	ImageSecret *string `json:"image_secret,omitempty"`

	// The name of the app revision.
	Name *string `json:"name,omitempty"`

	// Response model for probes.
	ProbeLiveness *Probe `json:"probe_liveness,omitempty"`

	// Response model for probes.
	ProbeReadiness *Probe `json:"probe_readiness,omitempty"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the app revision.
	ResourceType *string `json:"resource_type,omitempty"`

	// Optional arguments for the app that are passed to start the container. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments" validate:"required"`

	// Optional user ID (UID) to run the app.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Optional commands for the app that are passed to start the container. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands" validate:"required"`

	// References to config maps, secrets or literal values, which are defined by the app owner and are exposed as
	// environment variables in the application.
	RunEnvVariables []EnvVar `json:"run_env_variables" validate:"required"`

	// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
	// `none`, `reader`, and `writer`.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Mounts of config maps or secrets.
	RunVolumeMounts []VolumeMount `json:"run_volume_mounts" validate:"required"`

	// Optional maximum number of requests that can be processed concurrently per instance.
	ScaleConcurrency *int64 `json:"scale_concurrency,omitempty"`

	// Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use
	// this value to scale up instances based on concurrent number of requests. This option defaults to the value of the
	// `scale_concurrency` option, if not specified.
	ScaleConcurrencyTarget *int64 `json:"scale_concurrency_target,omitempty"`

	// Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit" validate:"required"`

	// Optional amount of time in seconds that delays the scale-down behavior for an app instance.
	ScaleDownDelay *int64 `json:"scale_down_delay,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit" validate:"required"`

	// Optional initial number of instances that are created upon app creation or app update.
	ScaleInitialInstances *int64 `json:"scale_initial_instances,omitempty"`

	// Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper
	// scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas
	// for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).
	ScaleMaxInstances *int64 `json:"scale_max_instances" validate:"required"`

	// Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit" validate:"required"`

	// Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if
	// not hit by any request for some time.
	ScaleMinInstances *int64 `json:"scale_min_instances" validate:"required"`

	// Optional amount of time in seconds that is allowed for a running app to respond to a request.
	ScaleRequestTimeout *int64 `json:"scale_request_timeout" validate:"required"`

	// The current status of the app revision.
	Status *string `json:"status,omitempty"`

	// The detailed status of the application revision.
	StatusDetails *AppRevisionStatus `json:"status_details,omitempty"`
}

// Constants associated with the AppRevision.ResourceType property.
// The type of the app revision.
const (
	AppRevision_ResourceType_AppRevisionV2 = "app_revision_v2"
)

// Constants associated with the AppRevision.RunServiceAccount property.
// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
// `none`, `reader`, and `writer`.
const (
	AppRevision_RunServiceAccount_Default = "default"
	AppRevision_RunServiceAccount_Manager = "manager"
	AppRevision_RunServiceAccount_None = "none"
	AppRevision_RunServiceAccount_Reader = "reader"
	AppRevision_RunServiceAccount_Writer = "writer"
)

// Constants associated with the AppRevision.Status property.
// The current status of the app revision.
const (
	AppRevision_Status_Failed = "failed"
	AppRevision_Status_Loading = "loading"
	AppRevision_Status_Ready = "ready"
	AppRevision_Status_Warning = "warning"
)

// UnmarshalAppRevision unmarshals an instance of AppRevision from the specified map of raw messages.
func UnmarshalAppRevision(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppRevision)
	err = core.UnmarshalPrimitive(m, "app_name", &obj.AppName)
	if err != nil {
		err = core.SDKErrorf(err, "", "app_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "computed_env_variables", &obj.ComputedEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "computed_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "image_port", &obj.ImagePort)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_reference", &obj.ImageReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_secret", &obj.ImageSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "probe_liveness", &obj.ProbeLiveness, UnmarshalProbe)
	if err != nil {
		err = core.SDKErrorf(err, "", "probe_liveness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "probe_readiness", &obj.ProbeReadiness, UnmarshalProbe)
	if err != nil {
		err = core.SDKErrorf(err, "", "probe_readiness-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_arguments", &obj.RunArguments)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_arguments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_as_user", &obj.RunAsUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_as_user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_commands", &obj.RunCommands)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_commands-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_service_account", &obj.RunServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_volume_mounts", &obj.RunVolumeMounts, UnmarshalVolumeMount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_volume_mounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency", &obj.ScaleConcurrency)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency_target", &obj.ScaleConcurrencyTarget)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency_target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_down_delay", &obj.ScaleDownDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_down_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_initial_instances", &obj.ScaleInitialInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_initial_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_instances", &obj.ScaleMaxInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_min_instances", &obj.ScaleMinInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_min_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_request_timeout", &obj.ScaleRequestTimeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_request_timeout-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalAppRevisionStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AppRevisionList : Contains a list of app revisions and pagination information.
type AppRevisionList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`

	// List of all app revisions.
	Revisions []AppRevision `json:"revisions" validate:"required"`
}

// UnmarshalAppRevisionList unmarshals an instance of AppRevisionList from the specified map of raw messages.
func UnmarshalAppRevisionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppRevisionList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "revisions", &obj.Revisions, UnmarshalAppRevision)
	if err != nil {
		err = core.SDKErrorf(err, "", "revisions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *AppRevisionList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// AppRevisionStatus : The detailed status of the application revision.
type AppRevisionStatus struct {
	// The number of running instances of the revision.
	ActualInstances *int64 `json:"actual_instances,omitempty"`

	// Optional information to provide more context in case of a 'failed' or 'warning' status.
	Reason *string `json:"reason,omitempty"`
}

// Constants associated with the AppRevisionStatus.Reason property.
// Optional information to provide more context in case of a 'failed' or 'warning' status.
const (
	AppRevisionStatus_Reason_ContainerFailedExitCode0 = "container_failed_exit_code_0"
	AppRevisionStatus_Reason_ContainerFailedExitCode1 = "container_failed_exit_code_1"
	AppRevisionStatus_Reason_ContainerFailedExitCode139 = "container_failed_exit_code_139"
	AppRevisionStatus_Reason_ContainerFailedExitCode24 = "container_failed_exit_code_24"
	AppRevisionStatus_Reason_Deploying = "deploying"
	AppRevisionStatus_Reason_DeployingWaitingForResources = "deploying_waiting_for_resources"
	AppRevisionStatus_Reason_FetchImageFailedMissingPullCredentials = "fetch_image_failed_missing_pull_credentials"
	AppRevisionStatus_Reason_FetchImageFailedMissingPullSecret = "fetch_image_failed_missing_pull_secret"
	AppRevisionStatus_Reason_FetchImageFailedRegistryNotFound = "fetch_image_failed_registry_not_found"
	AppRevisionStatus_Reason_FetchImageFailedUnknownManifest = "fetch_image_failed_unknown_manifest"
	AppRevisionStatus_Reason_FetchImageFailedUnknownRepository = "fetch_image_failed_unknown_repository"
	AppRevisionStatus_Reason_FetchImageFailedWrongPullCredentials = "fetch_image_failed_wrong_pull_credentials"
	AppRevisionStatus_Reason_ImagePullBackOff = "image_pull_back_off"
	AppRevisionStatus_Reason_InitialScaleNeverAchieved = "initial_scale_never_achieved"
	AppRevisionStatus_Reason_InvalidTarHeaderImagePullErr = "invalid_tar_header_image_pull_err"
	AppRevisionStatus_Reason_Ready = "ready"
	AppRevisionStatus_Reason_Waiting = "waiting"
)

// UnmarshalAppRevisionStatus unmarshals an instance of AppRevisionStatus from the specified map of raw messages.
func UnmarshalAppRevisionStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppRevisionStatus)
	err = core.UnmarshalPrimitive(m, "actual_instances", &obj.ActualInstances)
	if err != nil {
		err = core.SDKErrorf(err, "", "actual_instances-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AppStatus : The detailed status of the application.
type AppStatus struct {
	// Latest app revision that has been created.
	LatestCreatedRevision *string `json:"latest_created_revision,omitempty"`

	// Latest app revision that reached a ready state.
	LatestReadyRevision *string `json:"latest_ready_revision,omitempty"`

	// Optional information to provide more context in case of a 'failed' or 'warning' status.
	Reason *string `json:"reason,omitempty"`
}

// Constants associated with the AppStatus.Reason property.
// Optional information to provide more context in case of a 'failed' or 'warning' status.
const (
	AppStatus_Reason_Deploying = "deploying"
	AppStatus_Reason_NoRevisionReady = "no_revision_ready"
	AppStatus_Reason_Ready = "ready"
	AppStatus_Reason_ReadyButLatestRevisionFailed = "ready_but_latest_revision_failed"
	AppStatus_Reason_WaitingForResources = "waiting_for_resources"
)

// UnmarshalAppStatus unmarshals an instance of AppStatus from the specified map of raw messages.
func UnmarshalAppStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AppStatus)
	err = core.UnmarshalPrimitive(m, "latest_created_revision", &obj.LatestCreatedRevision)
	if err != nil {
		err = core.SDKErrorf(err, "", "latest_created_revision-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "latest_ready_revision", &obj.LatestReadyRevision)
	if err != nil {
		err = core.SDKErrorf(err, "", "latest_ready_revision-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Binding : Describes the model of a binding.
type Binding struct {
	// A reference to another component.
	Component *ComponentRef `json:"component" validate:"required"`

	// When you provision a new binding,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The ID of the binding.
	ID *string `json:"id,omitempty"`

	// The value that is set as a prefix in the component that is bound.
	Prefix *string `json:"prefix" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The type of the binding.
	ResourceType *string `json:"resource_type,omitempty"`

	// The service access secret that is bound to a component.
	SecretName *string `json:"secret_name" validate:"required"`

	// The current status of the binding.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the Binding.ResourceType property.
// The type of the binding.
const (
	Binding_ResourceType_BindingV2 = "binding_v2"
)

// UnmarshalBinding unmarshals an instance of Binding from the specified map of raw messages.
func UnmarshalBinding(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Binding)
	err = core.UnmarshalModel(m, "component", &obj.Component, UnmarshalComponentRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "component-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "prefix", &obj.Prefix)
	if err != nil {
		err = core.SDKErrorf(err, "", "prefix-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_name", &obj.SecretName)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BindingList : Contains a list of bindings and pagination information.
type BindingList struct {
	// List of all bindings.
	Bindings []Binding `json:"bindings" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalBindingList unmarshals an instance of BindingList from the specified map of raw messages.
func UnmarshalBindingList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BindingList)
	err = core.UnmarshalModel(m, "bindings", &obj.Bindings, UnmarshalBinding)
	if err != nil {
		err = core.SDKErrorf(err, "", "bindings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *BindingList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// Build : Response model for build definitions.
type Build struct {
	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The version of the build instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// When you provision a new build,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the build.
	Name *string `json:"name,omitempty"`

	// The name of the image.
	OutputImage *string `json:"output_image" validate:"required"`

	// The secret that is required to access the image registry. Make sure that the secret is granted with push permissions
	// towards the specified container registry namespace.
	OutputSecret *string `json:"output_secret" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the build.
	ResourceType *string `json:"resource_type,omitempty"`

	// Optional directory in the repository that contains the buildpacks file or the Dockerfile.
	SourceContextDir *string `json:"source_context_dir,omitempty"`

	// Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and
	// uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.
	SourceRevision *string `json:"source_revision,omitempty"`

	// Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`.
	// Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but
	// cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this
	// field must be omitted.
	SourceSecret *string `json:"source_secret,omitempty"`

	// Specifies the type of source to determine if your build source is in a repository or based on local source code.
	// * local - For builds from local source code.
	// * git - For builds from git version controlled source code.
	SourceType *string `json:"source_type" validate:"required"`

	// The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is
	// `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like
	// `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like
	// `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.
	SourceURL *string `json:"source_url,omitempty"`

	// The current status of the build.
	Status *string `json:"status,omitempty"`

	// The detailed status of the build.
	StatusDetails *BuildStatus `json:"status_details,omitempty"`

	// Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`,
	// `large`, `xlarge`, `xxlarge`.
	StrategySize *string `json:"strategy_size" validate:"required"`

	// Optional path to the specification file that is used for build strategies for building an image.
	StrategySpecFile *string `json:"strategy_spec_file,omitempty"`

	// The strategy to use for building the image.
	StrategyType *string `json:"strategy_type" validate:"required"`

	// The maximum amount of time, in seconds, that can pass before the build must succeed or fail.
	Timeout *int64 `json:"timeout,omitempty"`
}

// Constants associated with the Build.ResourceType property.
// The type of the build.
const (
	Build_ResourceType_BuildV2 = "build_v2"
)

// Constants associated with the Build.SourceType property.
// Specifies the type of source to determine if your build source is in a repository or based on local source code.
// * local - For builds from local source code.
// * git - For builds from git version controlled source code.
const (
	Build_SourceType_Git = "git"
	Build_SourceType_Local = "local"
)

// Constants associated with the Build.Status property.
// The current status of the build.
const (
	Build_Status_Failed = "failed"
	Build_Status_Ready = "ready"
)

// UnmarshalBuild unmarshals an instance of Build from the specified map of raw messages.
func UnmarshalBuild(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Build)
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "output_image", &obj.OutputImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_image-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "output_secret", &obj.OutputSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_context_dir", &obj.SourceContextDir)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_context_dir-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_revision", &obj.SourceRevision)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_revision-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_secret", &obj.SourceSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_type", &obj.SourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_url", &obj.SourceURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalBuildStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_size", &obj.StrategySize)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_spec_file", &obj.StrategySpecFile)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_spec_file-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_type", &obj.StrategyType)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeout-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BuildList : Contains a list of builds and pagination information.
type BuildList struct {
	// List of all builds.
	Builds []Build `json:"builds" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalBuildList unmarshals an instance of BuildList from the specified map of raw messages.
func UnmarshalBuildList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BuildList)
	err = core.UnmarshalModel(m, "builds", &obj.Builds, UnmarshalBuild)
	if err != nil {
		err = core.SDKErrorf(err, "", "builds-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *BuildList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// BuildPatch : Patch a build object.
type BuildPatch struct {
	// The name of the image.
	OutputImage *string `json:"output_image,omitempty"`

	// The secret that is required to access the image registry. Make sure that the secret is granted with push permissions
	// towards the specified container registry namespace.
	OutputSecret *string `json:"output_secret,omitempty"`

	// Optional directory in the repository that contains the buildpacks file or the Dockerfile.
	SourceContextDir *string `json:"source_context_dir,omitempty"`

	// Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and
	// uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.
	SourceRevision *string `json:"source_revision,omitempty"`

	// Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`.
	// Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but
	// cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this
	// field must be omitted.
	SourceSecret *string `json:"source_secret,omitempty"`

	// Specifies the type of source to determine if your build source is in a repository or based on local source code.
	// * local - For builds from local source code.
	// * git - For builds from git version controlled source code.
	SourceType *string `json:"source_type,omitempty"`

	// The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is
	// `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like
	// `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like
	// `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.
	SourceURL *string `json:"source_url,omitempty"`

	// Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`,
	// `large`, `xlarge`, `xxlarge`.
	StrategySize *string `json:"strategy_size,omitempty"`

	// Optional path to the specification file that is used for build strategies for building an image.
	StrategySpecFile *string `json:"strategy_spec_file,omitempty"`

	// The strategy to use for building the image.
	StrategyType *string `json:"strategy_type,omitempty"`

	// The maximum amount of time, in seconds, that can pass before the build must succeed or fail.
	Timeout *int64 `json:"timeout,omitempty"`
}

// Constants associated with the BuildPatch.SourceType property.
// Specifies the type of source to determine if your build source is in a repository or based on local source code.
// * local - For builds from local source code.
// * git - For builds from git version controlled source code.
const (
	BuildPatch_SourceType_Git = "git"
	BuildPatch_SourceType_Local = "local"
)

// UnmarshalBuildPatch unmarshals an instance of BuildPatch from the specified map of raw messages.
func UnmarshalBuildPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BuildPatch)
	err = core.UnmarshalPrimitive(m, "output_image", &obj.OutputImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_image-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "output_secret", &obj.OutputSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_context_dir", &obj.SourceContextDir)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_context_dir-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_revision", &obj.SourceRevision)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_revision-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_secret", &obj.SourceSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_type", &obj.SourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_url", &obj.SourceURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_size", &obj.StrategySize)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_spec_file", &obj.StrategySpecFile)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_spec_file-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_type", &obj.StrategyType)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeout-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the BuildPatch
func (buildPatch *BuildPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(buildPatch.OutputImage) {
		_patch["output_image"] = buildPatch.OutputImage
	}
	if !core.IsNil(buildPatch.OutputSecret) {
		_patch["output_secret"] = buildPatch.OutputSecret
	}
	if !core.IsNil(buildPatch.SourceContextDir) {
		_patch["source_context_dir"] = buildPatch.SourceContextDir
	}
	if !core.IsNil(buildPatch.SourceRevision) {
		_patch["source_revision"] = buildPatch.SourceRevision
	}
	if !core.IsNil(buildPatch.SourceSecret) {
		_patch["source_secret"] = buildPatch.SourceSecret
	}
	if !core.IsNil(buildPatch.SourceType) {
		_patch["source_type"] = buildPatch.SourceType
	}
	if !core.IsNil(buildPatch.SourceURL) {
		_patch["source_url"] = buildPatch.SourceURL
	}
	if !core.IsNil(buildPatch.StrategySize) {
		_patch["strategy_size"] = buildPatch.StrategySize
	}
	if !core.IsNil(buildPatch.StrategySpecFile) {
		_patch["strategy_spec_file"] = buildPatch.StrategySpecFile
	}
	if !core.IsNil(buildPatch.StrategyType) {
		_patch["strategy_type"] = buildPatch.StrategyType
	}
	if !core.IsNil(buildPatch.Timeout) {
		_patch["timeout"] = buildPatch.Timeout
	}

	return
}

// BuildRun : Response model for build run objects.
type BuildRun struct {
	// Optional name of the build on which this build run is based on. If specified, the build run will inherit the
	// configuration of the referenced build. If not specified, make sure to specify at least the fields `strategy_type`,
	// `source_url`, `output_image` and `output_secret` to describe the build run.
	BuildName *string `json:"build_name" validate:"required"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// When you trigger a new build run,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the build run.
	Name *string `json:"name" validate:"required"`

	// The name of the image.
	OutputImage *string `json:"output_image,omitempty"`

	// The secret that is required to access the image registry. Make sure that the secret is granted with push permissions
	// towards the specified container registry namespace.
	OutputSecret *string `json:"output_secret,omitempty"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the build run.
	ResourceType *string `json:"resource_type,omitempty"`

	// Optional service account, which is used for resource control.” or “Optional service account that is used for
	// resource control.
	ServiceAccount *string `json:"service_account,omitempty"`

	// Optional directory in the repository that contains the buildpacks file or the Dockerfile.
	SourceContextDir *string `json:"source_context_dir,omitempty"`

	// Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and
	// uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.
	SourceRevision *string `json:"source_revision,omitempty"`

	// Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`.
	// Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but
	// cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this
	// field must be omitted.
	SourceSecret *string `json:"source_secret,omitempty"`

	// Specifies the type of source to determine if your build source is in a repository or based on local source code.
	// * local - For builds from local source code.
	// * git - For builds from git version controlled source code.
	SourceType *string `json:"source_type,omitempty"`

	// The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is
	// `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like
	// `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like
	// `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.
	SourceURL *string `json:"source_url,omitempty"`

	// The current status of the build run.
	Status *string `json:"status,omitempty"`

	// Current status condition of a build run.
	StatusDetails *BuildRunStatus `json:"status_details,omitempty"`

	// Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`,
	// `large`, `xlarge`, `xxlarge`.
	StrategySize *string `json:"strategy_size,omitempty"`

	// Optional path to the specification file that is used for build strategies for building an image.
	StrategySpecFile *string `json:"strategy_spec_file,omitempty"`

	// The strategy to use for building the image.
	StrategyType *string `json:"strategy_type,omitempty"`

	// The maximum amount of time, in seconds, that can pass before the build must succeed or fail.
	Timeout *int64 `json:"timeout,omitempty"`
}

// Constants associated with the BuildRun.ResourceType property.
// The type of the build run.
const (
	BuildRun_ResourceType_BuildRunV2 = "build_run_v2"
)

// Constants associated with the BuildRun.ServiceAccount property.
// Optional service account, which is used for resource control.” or “Optional service account that is used for resource
// control.
const (
	BuildRun_ServiceAccount_Default = "default"
	BuildRun_ServiceAccount_Manager = "manager"
	BuildRun_ServiceAccount_None = "none"
	BuildRun_ServiceAccount_Reader = "reader"
	BuildRun_ServiceAccount_Writer = "writer"
)

// Constants associated with the BuildRun.SourceType property.
// Specifies the type of source to determine if your build source is in a repository or based on local source code.
// * local - For builds from local source code.
// * git - For builds from git version controlled source code.
const (
	BuildRun_SourceType_Git = "git"
	BuildRun_SourceType_Local = "local"
)

// Constants associated with the BuildRun.Status property.
// The current status of the build run.
const (
	BuildRun_Status_Failed = "failed"
	BuildRun_Status_Pending = "pending"
	BuildRun_Status_Running = "running"
	BuildRun_Status_Succeeded = "succeeded"
)

// UnmarshalBuildRun unmarshals an instance of BuildRun from the specified map of raw messages.
func UnmarshalBuildRun(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BuildRun)
	err = core.UnmarshalPrimitive(m, "build_name", &obj.BuildName)
	if err != nil {
		err = core.SDKErrorf(err, "", "build_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "output_image", &obj.OutputImage)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_image-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "output_secret", &obj.OutputSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_account", &obj.ServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_context_dir", &obj.SourceContextDir)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_context_dir-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_revision", &obj.SourceRevision)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_revision-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_secret", &obj.SourceSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_type", &obj.SourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_url", &obj.SourceURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalBuildRunStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_size", &obj.StrategySize)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_spec_file", &obj.StrategySpecFile)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_spec_file-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "strategy_type", &obj.StrategyType)
	if err != nil {
		err = core.SDKErrorf(err, "", "strategy_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeout-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BuildRunList : Contains a list of build runs and pagination information.
type BuildRunList struct {
	// List of all build runs.
	BuildRuns []BuildRun `json:"build_runs" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalBuildRunList unmarshals an instance of BuildRunList from the specified map of raw messages.
func UnmarshalBuildRunList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BuildRunList)
	err = core.UnmarshalModel(m, "build_runs", &obj.BuildRuns, UnmarshalBuildRun)
	if err != nil {
		err = core.SDKErrorf(err, "", "build_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *BuildRunList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// BuildRunStatus : Current status condition of a build run.
type BuildRunStatus struct {
	// Time the build run completed.
	CompletionTime *string `json:"completion_time,omitempty"`

	// The default branch name of the git source.
	GitBranchName *string `json:"git_branch_name,omitempty"`

	// The commit author of a git source.
	GitCommitAuthor *string `json:"git_commit_author,omitempty"`

	// The commit sha of the git source.
	GitCommitSha *string `json:"git_commit_sha,omitempty"`

	// Describes the time the build run completed.
	OutputDigest *string `json:"output_digest,omitempty"`

	// Optional information to provide more context in case of a 'failed' or 'warning' status.
	Reason *string `json:"reason,omitempty"`

	// The timestamp of the source.
	SourceTimestamp *string `json:"source_timestamp,omitempty"`

	// Time the build run started.
	StartTime *string `json:"start_time,omitempty"`
}

// Constants associated with the BuildRunStatus.Reason property.
// Optional information to provide more context in case of a 'failed' or 'warning' status.
const (
	BuildRunStatus_Reason_BuildNotFound = "build_not_found"
	BuildRunStatus_Reason_ExceededEphemeralStorage = "exceeded_ephemeral_storage"
	BuildRunStatus_Reason_Failed = "failed"
	BuildRunStatus_Reason_FailedToExecuteBuildRun = "failed_to_execute_build_run"
	BuildRunStatus_Reason_InvalidBuildConfiguration = "invalid_build_configuration"
	BuildRunStatus_Reason_MissingCodeRepoAccess = "missing_code_repo_access"
	BuildRunStatus_Reason_MissingRegistryAccess = "missing_registry_access"
	BuildRunStatus_Reason_MissingSecrets = "missing_secrets"
	BuildRunStatus_Reason_MissingTaskRun = "missing_task_run"
	BuildRunStatus_Reason_Pending = "pending"
	BuildRunStatus_Reason_PodEvicted = "pod_evicted"
	BuildRunStatus_Reason_PodEvictedBecauseOfStorageQuotaExceeds = "pod_evicted_because_of_storage_quota_exceeds"
	BuildRunStatus_Reason_Running = "running"
	BuildRunStatus_Reason_Succeeded = "succeeded"
	BuildRunStatus_Reason_TaskRunGenerationFailed = "task_run_generation_failed"
	BuildRunStatus_Reason_Timeout = "timeout"
	BuildRunStatus_Reason_UnknownStrategy = "unknown_strategy"
)

// UnmarshalBuildRunStatus unmarshals an instance of BuildRunStatus from the specified map of raw messages.
func UnmarshalBuildRunStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BuildRunStatus)
	err = core.UnmarshalPrimitive(m, "completion_time", &obj.CompletionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "completion_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_branch_name", &obj.GitBranchName)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_branch_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_commit_author", &obj.GitCommitAuthor)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_author-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_commit_sha", &obj.GitCommitSha)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_sha-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "output_digest", &obj.OutputDigest)
	if err != nil {
		err = core.SDKErrorf(err, "", "output_digest-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "source_timestamp", &obj.SourceTimestamp)
	if err != nil {
		err = core.SDKErrorf(err, "", "source_timestamp-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "start_time-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// BuildStatus : The detailed status of the build.
type BuildStatus struct {
	// Optional information to provide more context in case of a 'failed' or 'warning' status.
	Reason *string `json:"reason,omitempty"`
}

// Constants associated with the BuildStatus.Reason property.
// Optional information to provide more context in case of a 'failed' or 'warning' status.
const (
	BuildStatus_Reason_ClusterBuildStrategyNotFound = "cluster_build_strategy_not_found"
	BuildStatus_Reason_Failed = "failed"
	BuildStatus_Reason_MultipleSecretRefNotFound = "multiple_secret_ref_not_found"
	BuildStatus_Reason_Registered = "registered"
	BuildStatus_Reason_RemoteRepositoryUnreachable = "remote_repository_unreachable"
	BuildStatus_Reason_RuntimePathsCanNotBeEmpty = "runtime_paths_can_not_be_empty"
	BuildStatus_Reason_SetOwnerReferenceFailed = "set_owner_reference_failed"
	BuildStatus_Reason_SpecOutputSecretRefNotFound = "spec_output_secret_ref_not_found"
	BuildStatus_Reason_SpecRuntimeSecretRefNotFound = "spec_runtime_secret_ref_not_found"
	BuildStatus_Reason_SpecSourceSecretNotFound = "spec_source_secret_not_found"
	BuildStatus_Reason_StrategyNotFound = "strategy_not_found"
)

// UnmarshalBuildStatus unmarshals an instance of BuildStatus from the specified map of raw messages.
func UnmarshalBuildStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(BuildStatus)
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ComponentRef : A reference to another component.
type ComponentRef struct {
	// The name of the referenced component.
	Name *string `json:"name" validate:"required"`

	// The type of the referenced resource.
	ResourceType *string `json:"resource_type" validate:"required"`
}

// NewComponentRef : Instantiate ComponentRef (Generic Model Constructor)
func (*CodeEngineV2) NewComponentRef(name string, resourceType string) (_model *ComponentRef, err error) {
	_model = &ComponentRef{
		Name: core.StringPtr(name),
		ResourceType: core.StringPtr(resourceType),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalComponentRef unmarshals an instance of ComponentRef from the specified map of raw messages.
func UnmarshalComponentRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ComponentRef)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the ComponentRef
func (componentRef *ComponentRef) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(componentRef.Name) {
		_patch["name"] = componentRef.Name
	}
	if !core.IsNil(componentRef.ResourceType) {
		_patch["resource_type"] = componentRef.ResourceType
	}

	return
}

// ConfigMap : Describes the model of a configmap.
type ConfigMap struct {
	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The key-value pair for the config map. Values must be specified in `KEY=VALUE` format.
	Data map[string]string `json:"data,omitempty"`

	// The version of the config map instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// When you provision a new config map,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the config map.
	Name *string `json:"name" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the config map.
	ResourceType *string `json:"resource_type,omitempty"`
}

// Constants associated with the ConfigMap.ResourceType property.
// The type of the config map.
const (
	ConfigMap_ResourceType_ConfigMapV2 = "config_map_v2"
)

// UnmarshalConfigMap unmarshals an instance of ConfigMap from the specified map of raw messages.
func UnmarshalConfigMap(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigMap)
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigMapList : Contains a list of config maps and pagination information.
type ConfigMapList struct {
	// List of all config maps.
	ConfigMaps []ConfigMap `json:"config_maps" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalConfigMapList unmarshals an instance of ConfigMapList from the specified map of raw messages.
func UnmarshalConfigMapList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigMapList)
	err = core.UnmarshalModel(m, "config_maps", &obj.ConfigMaps, UnmarshalConfigMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "config_maps-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ConfigMapList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ContainerStatus : The status of a container.
type ContainerStatus struct {
	// Details of the observed container status.
	CurrentState *ContainerStatusDetails `json:"current_state,omitempty"`

	// Details of the observed container status.
	LastObservedState *ContainerStatusDetails `json:"last_observed_state,omitempty"`
}

// UnmarshalContainerStatus unmarshals an instance of ContainerStatus from the specified map of raw messages.
func UnmarshalContainerStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ContainerStatus)
	err = core.UnmarshalModel(m, "current_state", &obj.CurrentState, UnmarshalContainerStatusDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "current_state-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last_observed_state", &obj.LastObservedState, UnmarshalContainerStatusDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_observed_state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ContainerStatusDetails : Details of the observed container status.
type ContainerStatusDetails struct {
	// The time the container terminated. Only populated in an observed failure state.
	CompletedAt *string `json:"completed_at,omitempty"`

	// The status of the container.
	ContainerStatus *string `json:"container_status,omitempty"`

	// The exit code of the last termination of the container. Only populated in an observed failure state.
	ExitCode *int64 `json:"exit_code,omitempty"`

	// The reason the container is not yet running or has failed. Only populated in non-running states.
	Reason *string `json:"reason,omitempty"`

	// The time the container started.
	StartedAt *string `json:"started_at,omitempty"`
}

// Constants associated with the ContainerStatusDetails.Reason property.
// The reason the container is not yet running or has failed. Only populated in non-running states.
const (
	ContainerStatusDetails_Reason_ContainerFailedExitCode0 = "container_failed_exit_code_0"
	ContainerStatusDetails_Reason_ContainerFailedExitCode1 = "container_failed_exit_code_1"
	ContainerStatusDetails_Reason_ContainerFailedExitCode139 = "container_failed_exit_code_139"
	ContainerStatusDetails_Reason_ContainerFailedExitCode24 = "container_failed_exit_code_24"
	ContainerStatusDetails_Reason_Deploying = "deploying"
	ContainerStatusDetails_Reason_DeployingWaitingForResources = "deploying_waiting_for_resources"
	ContainerStatusDetails_Reason_FetchImageFailedMissingPullCredentials = "fetch_image_failed_missing_pull_credentials"
	ContainerStatusDetails_Reason_FetchImageFailedMissingPullSecret = "fetch_image_failed_missing_pull_secret"
	ContainerStatusDetails_Reason_FetchImageFailedRegistryNotFound = "fetch_image_failed_registry_not_found"
	ContainerStatusDetails_Reason_FetchImageFailedUnknownManifest = "fetch_image_failed_unknown_manifest"
	ContainerStatusDetails_Reason_FetchImageFailedUnknownRepository = "fetch_image_failed_unknown_repository"
	ContainerStatusDetails_Reason_FetchImageFailedWrongPullCredentials = "fetch_image_failed_wrong_pull_credentials"
	ContainerStatusDetails_Reason_ImagePullBackOff = "image_pull_back_off"
	ContainerStatusDetails_Reason_InitialScaleNeverAchieved = "initial_scale_never_achieved"
	ContainerStatusDetails_Reason_InvalidTarHeaderImagePullErr = "invalid_tar_header_image_pull_err"
	ContainerStatusDetails_Reason_Ready = "ready"
	ContainerStatusDetails_Reason_Waiting = "waiting"
)

// UnmarshalContainerStatusDetails unmarshals an instance of ContainerStatusDetails from the specified map of raw messages.
func UnmarshalContainerStatusDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ContainerStatusDetails)
	err = core.UnmarshalPrimitive(m, "completed_at", &obj.CompletedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "completed_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "container_status", &obj.ContainerStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "container_status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "exit_code", &obj.ExitCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "exit_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "started_at", &obj.StartedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "started_at-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateAllowedOutboundDestinationOptions : The CreateAllowedOutboundDestination options.
type CreateAllowedOutboundDestinationOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// AllowedOutboundDestination prototype.
	AllowedOutboundDestination AllowedOutboundDestinationPrototypeIntf `json:"allowed_outbound_destination" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateAllowedOutboundDestinationOptions : Instantiate CreateAllowedOutboundDestinationOptions
func (*CodeEngineV2) NewCreateAllowedOutboundDestinationOptions(projectID string, allowedOutboundDestination AllowedOutboundDestinationPrototypeIntf) *CreateAllowedOutboundDestinationOptions {
	return &CreateAllowedOutboundDestinationOptions{
		ProjectID: core.StringPtr(projectID),
		AllowedOutboundDestination: allowedOutboundDestination,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateAllowedOutboundDestinationOptions) SetProjectID(projectID string) *CreateAllowedOutboundDestinationOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetAllowedOutboundDestination : Allow user to set AllowedOutboundDestination
func (_options *CreateAllowedOutboundDestinationOptions) SetAllowedOutboundDestination(allowedOutboundDestination AllowedOutboundDestinationPrototypeIntf) *CreateAllowedOutboundDestinationOptions {
	_options.AllowedOutboundDestination = allowedOutboundDestination
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAllowedOutboundDestinationOptions) SetHeaders(param map[string]string) *CreateAllowedOutboundDestinationOptions {
	options.Headers = param
	return options
}

// CreateAppOptions : The CreateApp options.
type CreateAppOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of the image that is used for this app. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference" validate:"required"`

	// The name of the app. Use a name that is unique within the project.
	Name *string `json:"name" validate:"required"`

	// Optional port the app listens on. While the app will always be exposed via port `443` for end users, this port is
	// used to connect to the port that is exposed by the container image.
	ImagePort *int64 `json:"image_port,omitempty"`

	// Optional name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the app will be created but cannot reach the ready status, until this property is provided, too.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Optional value controlling which of the system managed domain mappings will be setup for the application. Valid
	// values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project
	// supports application private visibility.
	ManagedDomainMappings *string `json:"managed_domain_mappings,omitempty"`

	// Request model for probes.
	ProbeLiveness *ProbePrototype `json:"probe_liveness,omitempty"`

	// Request model for probes.
	ProbeReadiness *ProbePrototype `json:"probe_readiness,omitempty"`

	// Optional arguments for the app that are passed to start the container. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments,omitempty"`

	// Optional user ID (UID) to run the app.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Optional commands for the app that are passed to start the container. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands,omitempty"`

	// Optional references to config maps, secrets or literal values that are exposed as environment variables within the
	// running application.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
	// `none`, `reader`, and `writer`.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or a secrets.
	RunVolumeMounts []VolumeMountPrototype `json:"run_volume_mounts,omitempty"`

	// Optional maximum number of requests that can be processed concurrently per instance.
	ScaleConcurrency *int64 `json:"scale_concurrency,omitempty"`

	// Optional threshold of concurrent requests per instance at which one or more additional instances are created. Use
	// this value to scale up instances based on concurrent number of requests. This option defaults to the value of the
	// `scale_concurrency` option, if not specified.
	ScaleConcurrencyTarget *int64 `json:"scale_concurrency_target,omitempty"`

	// Optional number of CPU set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of time in seconds that delays the scale-down behavior for an app instance.
	ScaleDownDelay *int64 `json:"scale_down_delay,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the app. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit,omitempty"`

	// Optional initial number of instances that are created upon app creation or app update.
	ScaleInitialInstances *int64 `json:"scale_initial_instances,omitempty"`

	// Optional maximum number of instances for this app. If you set this value to `0`, this property does not set a upper
	// scaling limit. However, the app scaling is still limited by the project quota for instances. See [Limits and quotas
	// for Code Engine](https://cloud.ibm.com/docs/codeengine?topic=codeengine-limits).
	ScaleMaxInstances *int64 `json:"scale_max_instances,omitempty"`

	// Optional amount of memory set for the instance of the app. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// Optional minimum number of instances for this app. If you set this value to `0`, the app will scale down to zero, if
	// not hit by any request for some time.
	ScaleMinInstances *int64 `json:"scale_min_instances,omitempty"`

	// Optional amount of time in seconds that is allowed for a running app to respond to a request.
	ScaleRequestTimeout *int64 `json:"scale_request_timeout,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateAppOptions.ManagedDomainMappings property.
// Optional value controlling which of the system managed domain mappings will be setup for the application. Valid
// values are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project
// supports application private visibility.
const (
	CreateAppOptions_ManagedDomainMappings_Local = "local"
	CreateAppOptions_ManagedDomainMappings_LocalPrivate = "local_private"
	CreateAppOptions_ManagedDomainMappings_LocalPublic = "local_public"
)

// Constants associated with the CreateAppOptions.RunServiceAccount property.
// Optional name of the service account. For built-in service accounts, you can use the shortened names `manager` ,
// `none`, `reader`, and `writer`.
const (
	CreateAppOptions_RunServiceAccount_Default = "default"
	CreateAppOptions_RunServiceAccount_Manager = "manager"
	CreateAppOptions_RunServiceAccount_None = "none"
	CreateAppOptions_RunServiceAccount_Reader = "reader"
	CreateAppOptions_RunServiceAccount_Writer = "writer"
)

// NewCreateAppOptions : Instantiate CreateAppOptions
func (*CodeEngineV2) NewCreateAppOptions(projectID string, imageReference string, name string) *CreateAppOptions {
	return &CreateAppOptions{
		ProjectID: core.StringPtr(projectID),
		ImageReference: core.StringPtr(imageReference),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateAppOptions) SetProjectID(projectID string) *CreateAppOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetImageReference : Allow user to set ImageReference
func (_options *CreateAppOptions) SetImageReference(imageReference string) *CreateAppOptions {
	_options.ImageReference = core.StringPtr(imageReference)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateAppOptions) SetName(name string) *CreateAppOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetImagePort : Allow user to set ImagePort
func (_options *CreateAppOptions) SetImagePort(imagePort int64) *CreateAppOptions {
	_options.ImagePort = core.Int64Ptr(imagePort)
	return _options
}

// SetImageSecret : Allow user to set ImageSecret
func (_options *CreateAppOptions) SetImageSecret(imageSecret string) *CreateAppOptions {
	_options.ImageSecret = core.StringPtr(imageSecret)
	return _options
}

// SetManagedDomainMappings : Allow user to set ManagedDomainMappings
func (_options *CreateAppOptions) SetManagedDomainMappings(managedDomainMappings string) *CreateAppOptions {
	_options.ManagedDomainMappings = core.StringPtr(managedDomainMappings)
	return _options
}

// SetProbeLiveness : Allow user to set ProbeLiveness
func (_options *CreateAppOptions) SetProbeLiveness(probeLiveness *ProbePrototype) *CreateAppOptions {
	_options.ProbeLiveness = probeLiveness
	return _options
}

// SetProbeReadiness : Allow user to set ProbeReadiness
func (_options *CreateAppOptions) SetProbeReadiness(probeReadiness *ProbePrototype) *CreateAppOptions {
	_options.ProbeReadiness = probeReadiness
	return _options
}

// SetRunArguments : Allow user to set RunArguments
func (_options *CreateAppOptions) SetRunArguments(runArguments []string) *CreateAppOptions {
	_options.RunArguments = runArguments
	return _options
}

// SetRunAsUser : Allow user to set RunAsUser
func (_options *CreateAppOptions) SetRunAsUser(runAsUser int64) *CreateAppOptions {
	_options.RunAsUser = core.Int64Ptr(runAsUser)
	return _options
}

// SetRunCommands : Allow user to set RunCommands
func (_options *CreateAppOptions) SetRunCommands(runCommands []string) *CreateAppOptions {
	_options.RunCommands = runCommands
	return _options
}

// SetRunEnvVariables : Allow user to set RunEnvVariables
func (_options *CreateAppOptions) SetRunEnvVariables(runEnvVariables []EnvVarPrototype) *CreateAppOptions {
	_options.RunEnvVariables = runEnvVariables
	return _options
}

// SetRunServiceAccount : Allow user to set RunServiceAccount
func (_options *CreateAppOptions) SetRunServiceAccount(runServiceAccount string) *CreateAppOptions {
	_options.RunServiceAccount = core.StringPtr(runServiceAccount)
	return _options
}

// SetRunVolumeMounts : Allow user to set RunVolumeMounts
func (_options *CreateAppOptions) SetRunVolumeMounts(runVolumeMounts []VolumeMountPrototype) *CreateAppOptions {
	_options.RunVolumeMounts = runVolumeMounts
	return _options
}

// SetScaleConcurrency : Allow user to set ScaleConcurrency
func (_options *CreateAppOptions) SetScaleConcurrency(scaleConcurrency int64) *CreateAppOptions {
	_options.ScaleConcurrency = core.Int64Ptr(scaleConcurrency)
	return _options
}

// SetScaleConcurrencyTarget : Allow user to set ScaleConcurrencyTarget
func (_options *CreateAppOptions) SetScaleConcurrencyTarget(scaleConcurrencyTarget int64) *CreateAppOptions {
	_options.ScaleConcurrencyTarget = core.Int64Ptr(scaleConcurrencyTarget)
	return _options
}

// SetScaleCpuLimit : Allow user to set ScaleCpuLimit
func (_options *CreateAppOptions) SetScaleCpuLimit(scaleCpuLimit string) *CreateAppOptions {
	_options.ScaleCpuLimit = core.StringPtr(scaleCpuLimit)
	return _options
}

// SetScaleDownDelay : Allow user to set ScaleDownDelay
func (_options *CreateAppOptions) SetScaleDownDelay(scaleDownDelay int64) *CreateAppOptions {
	_options.ScaleDownDelay = core.Int64Ptr(scaleDownDelay)
	return _options
}

// SetScaleEphemeralStorageLimit : Allow user to set ScaleEphemeralStorageLimit
func (_options *CreateAppOptions) SetScaleEphemeralStorageLimit(scaleEphemeralStorageLimit string) *CreateAppOptions {
	_options.ScaleEphemeralStorageLimit = core.StringPtr(scaleEphemeralStorageLimit)
	return _options
}

// SetScaleInitialInstances : Allow user to set ScaleInitialInstances
func (_options *CreateAppOptions) SetScaleInitialInstances(scaleInitialInstances int64) *CreateAppOptions {
	_options.ScaleInitialInstances = core.Int64Ptr(scaleInitialInstances)
	return _options
}

// SetScaleMaxInstances : Allow user to set ScaleMaxInstances
func (_options *CreateAppOptions) SetScaleMaxInstances(scaleMaxInstances int64) *CreateAppOptions {
	_options.ScaleMaxInstances = core.Int64Ptr(scaleMaxInstances)
	return _options
}

// SetScaleMemoryLimit : Allow user to set ScaleMemoryLimit
func (_options *CreateAppOptions) SetScaleMemoryLimit(scaleMemoryLimit string) *CreateAppOptions {
	_options.ScaleMemoryLimit = core.StringPtr(scaleMemoryLimit)
	return _options
}

// SetScaleMinInstances : Allow user to set ScaleMinInstances
func (_options *CreateAppOptions) SetScaleMinInstances(scaleMinInstances int64) *CreateAppOptions {
	_options.ScaleMinInstances = core.Int64Ptr(scaleMinInstances)
	return _options
}

// SetScaleRequestTimeout : Allow user to set ScaleRequestTimeout
func (_options *CreateAppOptions) SetScaleRequestTimeout(scaleRequestTimeout int64) *CreateAppOptions {
	_options.ScaleRequestTimeout = core.Int64Ptr(scaleRequestTimeout)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateAppOptions) SetHeaders(param map[string]string) *CreateAppOptions {
	options.Headers = param
	return options
}

// CreateBindingOptions : The CreateBinding options.
type CreateBindingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// A reference to another component.
	Component *ComponentRef `json:"component" validate:"required"`

	// Optional value that is set as prefix in the component that is bound. Will be generated if not provided.
	Prefix *string `json:"prefix" validate:"required"`

	// The service access secret that is bound to a component.
	SecretName *string `json:"secret_name" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateBindingOptions : Instantiate CreateBindingOptions
func (*CodeEngineV2) NewCreateBindingOptions(projectID string, component *ComponentRef, prefix string, secretName string) *CreateBindingOptions {
	return &CreateBindingOptions{
		ProjectID: core.StringPtr(projectID),
		Component: component,
		Prefix: core.StringPtr(prefix),
		SecretName: core.StringPtr(secretName),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateBindingOptions) SetProjectID(projectID string) *CreateBindingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetComponent : Allow user to set Component
func (_options *CreateBindingOptions) SetComponent(component *ComponentRef) *CreateBindingOptions {
	_options.Component = component
	return _options
}

// SetPrefix : Allow user to set Prefix
func (_options *CreateBindingOptions) SetPrefix(prefix string) *CreateBindingOptions {
	_options.Prefix = core.StringPtr(prefix)
	return _options
}

// SetSecretName : Allow user to set SecretName
func (_options *CreateBindingOptions) SetSecretName(secretName string) *CreateBindingOptions {
	_options.SecretName = core.StringPtr(secretName)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateBindingOptions) SetHeaders(param map[string]string) *CreateBindingOptions {
	options.Headers = param
	return options
}

// CreateBuildOptions : The CreateBuild options.
type CreateBuildOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of the build. Use a name that is unique within the project.
	Name *string `json:"name" validate:"required"`

	// The name of the image.
	OutputImage *string `json:"output_image" validate:"required"`

	// The secret that is required to access the image registry. Make sure that the secret is granted with push permissions
	// towards the specified container registry namespace.
	OutputSecret *string `json:"output_secret" validate:"required"`

	// The strategy to use for building the image.
	StrategyType *string `json:"strategy_type" validate:"required"`

	// Optional directory in the repository that contains the buildpacks file or the Dockerfile.
	SourceContextDir *string `json:"source_context_dir,omitempty"`

	// Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and
	// uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.
	SourceRevision *string `json:"source_revision,omitempty"`

	// Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`.
	// Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but
	// cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this
	// field must be omitted.
	SourceSecret *string `json:"source_secret,omitempty"`

	// Specifies the type of source to determine if your build source is in a repository or based on local source code.
	// * local - For builds from local source code.
	// * git - For builds from git version controlled source code.
	SourceType *string `json:"source_type,omitempty"`

	// The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is
	// `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like
	// `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like
	// `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.
	SourceURL *string `json:"source_url,omitempty"`

	// Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`,
	// `large`, `xlarge`, `xxlarge`.
	StrategySize *string `json:"strategy_size,omitempty"`

	// Optional path to the specification file that is used for build strategies for building an image.
	StrategySpecFile *string `json:"strategy_spec_file,omitempty"`

	// The maximum amount of time, in seconds, that can pass before the build must succeed or fail.
	Timeout *int64 `json:"timeout,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateBuildOptions.SourceType property.
// Specifies the type of source to determine if your build source is in a repository or based on local source code.
// * local - For builds from local source code.
// * git - For builds from git version controlled source code.
const (
	CreateBuildOptions_SourceType_Git = "git"
	CreateBuildOptions_SourceType_Local = "local"
)

// NewCreateBuildOptions : Instantiate CreateBuildOptions
func (*CodeEngineV2) NewCreateBuildOptions(projectID string, name string, outputImage string, outputSecret string, strategyType string) *CreateBuildOptions {
	return &CreateBuildOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		OutputImage: core.StringPtr(outputImage),
		OutputSecret: core.StringPtr(outputSecret),
		StrategyType: core.StringPtr(strategyType),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateBuildOptions) SetProjectID(projectID string) *CreateBuildOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateBuildOptions) SetName(name string) *CreateBuildOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetOutputImage : Allow user to set OutputImage
func (_options *CreateBuildOptions) SetOutputImage(outputImage string) *CreateBuildOptions {
	_options.OutputImage = core.StringPtr(outputImage)
	return _options
}

// SetOutputSecret : Allow user to set OutputSecret
func (_options *CreateBuildOptions) SetOutputSecret(outputSecret string) *CreateBuildOptions {
	_options.OutputSecret = core.StringPtr(outputSecret)
	return _options
}

// SetStrategyType : Allow user to set StrategyType
func (_options *CreateBuildOptions) SetStrategyType(strategyType string) *CreateBuildOptions {
	_options.StrategyType = core.StringPtr(strategyType)
	return _options
}

// SetSourceContextDir : Allow user to set SourceContextDir
func (_options *CreateBuildOptions) SetSourceContextDir(sourceContextDir string) *CreateBuildOptions {
	_options.SourceContextDir = core.StringPtr(sourceContextDir)
	return _options
}

// SetSourceRevision : Allow user to set SourceRevision
func (_options *CreateBuildOptions) SetSourceRevision(sourceRevision string) *CreateBuildOptions {
	_options.SourceRevision = core.StringPtr(sourceRevision)
	return _options
}

// SetSourceSecret : Allow user to set SourceSecret
func (_options *CreateBuildOptions) SetSourceSecret(sourceSecret string) *CreateBuildOptions {
	_options.SourceSecret = core.StringPtr(sourceSecret)
	return _options
}

// SetSourceType : Allow user to set SourceType
func (_options *CreateBuildOptions) SetSourceType(sourceType string) *CreateBuildOptions {
	_options.SourceType = core.StringPtr(sourceType)
	return _options
}

// SetSourceURL : Allow user to set SourceURL
func (_options *CreateBuildOptions) SetSourceURL(sourceURL string) *CreateBuildOptions {
	_options.SourceURL = core.StringPtr(sourceURL)
	return _options
}

// SetStrategySize : Allow user to set StrategySize
func (_options *CreateBuildOptions) SetStrategySize(strategySize string) *CreateBuildOptions {
	_options.StrategySize = core.StringPtr(strategySize)
	return _options
}

// SetStrategySpecFile : Allow user to set StrategySpecFile
func (_options *CreateBuildOptions) SetStrategySpecFile(strategySpecFile string) *CreateBuildOptions {
	_options.StrategySpecFile = core.StringPtr(strategySpecFile)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *CreateBuildOptions) SetTimeout(timeout int64) *CreateBuildOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateBuildOptions) SetHeaders(param map[string]string) *CreateBuildOptions {
	options.Headers = param
	return options
}

// CreateBuildRunOptions : The CreateBuildRun options.
type CreateBuildRunOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional name of the build on which this build run is based on. If specified, the build run will inherit the
	// configuration of the referenced build. If not specified, make sure to specify at least the fields `strategy_type`,
	// `source_url`, `output_image` and `output_secret` to describe the build run.
	BuildName *string `json:"build_name,omitempty"`

	// Name of the build run. This field is optional, if the field `build_name` is specified and its value will be
	// generated like so: `[BUILD_NAME]-run-[timestamp with format: YYMMDD-hhmmss] if not set.`.
	Name *string `json:"name,omitempty"`

	// The name of the image.
	OutputImage *string `json:"output_image,omitempty"`

	// The secret that is required to access the image registry. Make sure that the secret is granted with push permissions
	// towards the specified container registry namespace.
	OutputSecret *string `json:"output_secret,omitempty"`

	// Optional service account, which is used for resource control.” or “Optional service account that is used for
	// resource control.
	ServiceAccount *string `json:"service_account,omitempty"`

	// Optional directory in the repository that contains the buildpacks file or the Dockerfile.
	SourceContextDir *string `json:"source_context_dir,omitempty"`

	// Commit, tag, or branch in the source repository to pull. This field is optional if the `source_type` is `git` and
	// uses the HEAD of default branch if not specified. If the `source_type` value is `local`, this field must be omitted.
	SourceRevision *string `json:"source_revision,omitempty"`

	// Name of the secret that is used access the repository source. This field is optional if the `source_type` is `git`.
	// Additionally, if the `source_url` points to a repository that requires authentication, the build will be created but
	// cannot access any source code, until this property is provided, too. If the `source_type` value is `local`, this
	// field must be omitted.
	SourceSecret *string `json:"source_secret,omitempty"`

	// Specifies the type of source to determine if your build source is in a repository or based on local source code.
	// * local - For builds from local source code.
	// * git - For builds from git version controlled source code.
	SourceType *string `json:"source_type,omitempty"`

	// The URL of the code repository. This field is required if the `source_type` is `git`. If the `source_type` value is
	// `local`, this field must be omitted. If the repository is publicly available you can provide a 'https' URL like
	// `https://github.com/IBM/CodeEngine`. If the repository requires authentication, you need to provide a 'ssh' URL like
	// `git@github.com:IBM/CodeEngine.git` along with a `source_secret` that points to a secret of format `ssh_auth`.
	SourceURL *string `json:"source_url,omitempty"`

	// Optional size for the build, which determines the amount of resources used. Build sizes are `small`, `medium`,
	// `large`, `xlarge`, `xxlarge`.
	StrategySize *string `json:"strategy_size,omitempty"`

	// Optional path to the specification file that is used for build strategies for building an image.
	StrategySpecFile *string `json:"strategy_spec_file,omitempty"`

	// The strategy to use for building the image.
	StrategyType *string `json:"strategy_type,omitempty"`

	// The maximum amount of time, in seconds, that can pass before the build must succeed or fail.
	Timeout *int64 `json:"timeout,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateBuildRunOptions.ServiceAccount property.
// Optional service account, which is used for resource control.” or “Optional service account that is used for resource
// control.
const (
	CreateBuildRunOptions_ServiceAccount_Default = "default"
	CreateBuildRunOptions_ServiceAccount_Manager = "manager"
	CreateBuildRunOptions_ServiceAccount_None = "none"
	CreateBuildRunOptions_ServiceAccount_Reader = "reader"
	CreateBuildRunOptions_ServiceAccount_Writer = "writer"
)

// Constants associated with the CreateBuildRunOptions.SourceType property.
// Specifies the type of source to determine if your build source is in a repository or based on local source code.
// * local - For builds from local source code.
// * git - For builds from git version controlled source code.
const (
	CreateBuildRunOptions_SourceType_Git = "git"
	CreateBuildRunOptions_SourceType_Local = "local"
)

// NewCreateBuildRunOptions : Instantiate CreateBuildRunOptions
func (*CodeEngineV2) NewCreateBuildRunOptions(projectID string) *CreateBuildRunOptions {
	return &CreateBuildRunOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateBuildRunOptions) SetProjectID(projectID string) *CreateBuildRunOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetBuildName : Allow user to set BuildName
func (_options *CreateBuildRunOptions) SetBuildName(buildName string) *CreateBuildRunOptions {
	_options.BuildName = core.StringPtr(buildName)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateBuildRunOptions) SetName(name string) *CreateBuildRunOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetOutputImage : Allow user to set OutputImage
func (_options *CreateBuildRunOptions) SetOutputImage(outputImage string) *CreateBuildRunOptions {
	_options.OutputImage = core.StringPtr(outputImage)
	return _options
}

// SetOutputSecret : Allow user to set OutputSecret
func (_options *CreateBuildRunOptions) SetOutputSecret(outputSecret string) *CreateBuildRunOptions {
	_options.OutputSecret = core.StringPtr(outputSecret)
	return _options
}

// SetServiceAccount : Allow user to set ServiceAccount
func (_options *CreateBuildRunOptions) SetServiceAccount(serviceAccount string) *CreateBuildRunOptions {
	_options.ServiceAccount = core.StringPtr(serviceAccount)
	return _options
}

// SetSourceContextDir : Allow user to set SourceContextDir
func (_options *CreateBuildRunOptions) SetSourceContextDir(sourceContextDir string) *CreateBuildRunOptions {
	_options.SourceContextDir = core.StringPtr(sourceContextDir)
	return _options
}

// SetSourceRevision : Allow user to set SourceRevision
func (_options *CreateBuildRunOptions) SetSourceRevision(sourceRevision string) *CreateBuildRunOptions {
	_options.SourceRevision = core.StringPtr(sourceRevision)
	return _options
}

// SetSourceSecret : Allow user to set SourceSecret
func (_options *CreateBuildRunOptions) SetSourceSecret(sourceSecret string) *CreateBuildRunOptions {
	_options.SourceSecret = core.StringPtr(sourceSecret)
	return _options
}

// SetSourceType : Allow user to set SourceType
func (_options *CreateBuildRunOptions) SetSourceType(sourceType string) *CreateBuildRunOptions {
	_options.SourceType = core.StringPtr(sourceType)
	return _options
}

// SetSourceURL : Allow user to set SourceURL
func (_options *CreateBuildRunOptions) SetSourceURL(sourceURL string) *CreateBuildRunOptions {
	_options.SourceURL = core.StringPtr(sourceURL)
	return _options
}

// SetStrategySize : Allow user to set StrategySize
func (_options *CreateBuildRunOptions) SetStrategySize(strategySize string) *CreateBuildRunOptions {
	_options.StrategySize = core.StringPtr(strategySize)
	return _options
}

// SetStrategySpecFile : Allow user to set StrategySpecFile
func (_options *CreateBuildRunOptions) SetStrategySpecFile(strategySpecFile string) *CreateBuildRunOptions {
	_options.StrategySpecFile = core.StringPtr(strategySpecFile)
	return _options
}

// SetStrategyType : Allow user to set StrategyType
func (_options *CreateBuildRunOptions) SetStrategyType(strategyType string) *CreateBuildRunOptions {
	_options.StrategyType = core.StringPtr(strategyType)
	return _options
}

// SetTimeout : Allow user to set Timeout
func (_options *CreateBuildRunOptions) SetTimeout(timeout int64) *CreateBuildRunOptions {
	_options.Timeout = core.Int64Ptr(timeout)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateBuildRunOptions) SetHeaders(param map[string]string) *CreateBuildRunOptions {
	options.Headers = param
	return options
}

// CreateConfigMapOptions : The CreateConfigMap options.
type CreateConfigMapOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of the configmap. Use a name that is unique within the project.
	Name *string `json:"name" validate:"required"`

	// The key-value pair for the config map. Values must be specified in `KEY=VALUE` format. Each `KEY` field must consist
	// of alphanumeric characters, `-`, `_` or `.` and must not be exceed a max length of 253 characters. Each `VALUE`
	// field can consists of any character and must not be exceed a max length of 1048576 characters.
	Data map[string]string `json:"data,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateConfigMapOptions : Instantiate CreateConfigMapOptions
func (*CodeEngineV2) NewCreateConfigMapOptions(projectID string, name string) *CreateConfigMapOptions {
	return &CreateConfigMapOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateConfigMapOptions) SetProjectID(projectID string) *CreateConfigMapOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateConfigMapOptions) SetName(name string) *CreateConfigMapOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetData : Allow user to set Data
func (_options *CreateConfigMapOptions) SetData(data map[string]string) *CreateConfigMapOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateConfigMapOptions) SetHeaders(param map[string]string) *CreateConfigMapOptions {
	options.Headers = param
	return options
}

// CreateDomainMappingOptions : The CreateDomainMapping options.
type CreateDomainMappingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// A reference to another component.
	Component *ComponentRef `json:"component" validate:"required"`

	// The name of the domain mapping.
	Name *string `json:"name" validate:"required"`

	// The name of the TLS secret that includes the certificate and private key of this domain mapping.
	TlsSecret *string `json:"tls_secret" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateDomainMappingOptions : Instantiate CreateDomainMappingOptions
func (*CodeEngineV2) NewCreateDomainMappingOptions(projectID string, component *ComponentRef, name string, tlsSecret string) *CreateDomainMappingOptions {
	return &CreateDomainMappingOptions{
		ProjectID: core.StringPtr(projectID),
		Component: component,
		Name: core.StringPtr(name),
		TlsSecret: core.StringPtr(tlsSecret),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateDomainMappingOptions) SetProjectID(projectID string) *CreateDomainMappingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetComponent : Allow user to set Component
func (_options *CreateDomainMappingOptions) SetComponent(component *ComponentRef) *CreateDomainMappingOptions {
	_options.Component = component
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateDomainMappingOptions) SetName(name string) *CreateDomainMappingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetTlsSecret : Allow user to set TlsSecret
func (_options *CreateDomainMappingOptions) SetTlsSecret(tlsSecret string) *CreateDomainMappingOptions {
	_options.TlsSecret = core.StringPtr(tlsSecret)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateDomainMappingOptions) SetHeaders(param map[string]string) *CreateDomainMappingOptions {
	options.Headers = param
	return options
}

// CreateFunctionOptions : The CreateFunction options.
type CreateFunctionOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Specifies either a reference to a code bundle or the source code itself. To specify the source code, use the data
	// URL scheme and include the source code as base64 encoded. The data URL scheme is defined in [RFC
	// 2397](https://tools.ietf.org/html/rfc2397).
	CodeReference *string `json:"code_reference" validate:"required"`

	// The name of the function. Use a name that is unique within the project.
	Name *string `json:"name" validate:"required"`

	// The managed runtime used to execute the injected code.
	Runtime *string `json:"runtime" validate:"required"`

	// Specifies whether the code is binary or not. Defaults to false when `code_reference` is set to a data URL. When
	// `code_reference` is set to a code bundle URL, this field is always true.
	CodeBinary *bool `json:"code_binary,omitempty"`

	// Specifies the name of the function that should be invoked.
	CodeMain *string `json:"code_main,omitempty"`

	// The name of the secret that is used to access the specified `code_reference`. The secret is used to authenticate
	// with a non-public endpoint that is specified as`code_reference`.
	CodeSecret *string `json:"code_secret,omitempty"`

	// Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values
	// are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports
	// function private visibility.
	ManagedDomainMappings *string `json:"managed_domain_mappings,omitempty"`

	// Optional references to config maps, secrets or literal values.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// Number of parallel requests handled by a single instance, supported only by Node.js, default is `1`.
	ScaleConcurrency *int64 `json:"scale_concurrency,omitempty"`

	// Optional amount of CPU set for the instance of the function. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of time in seconds that delays the scale down behavior for a function.
	ScaleDownDelay *int64 `json:"scale_down_delay,omitempty"`

	// Timeout in secs after which the function is terminated.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the function. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateFunctionOptions.ManagedDomainMappings property.
// Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values
// are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports
// function private visibility.
const (
	CreateFunctionOptions_ManagedDomainMappings_Local = "local"
	CreateFunctionOptions_ManagedDomainMappings_LocalPrivate = "local_private"
	CreateFunctionOptions_ManagedDomainMappings_LocalPublic = "local_public"
)

// NewCreateFunctionOptions : Instantiate CreateFunctionOptions
func (*CodeEngineV2) NewCreateFunctionOptions(projectID string, codeReference string, name string, runtime string) *CreateFunctionOptions {
	return &CreateFunctionOptions{
		ProjectID: core.StringPtr(projectID),
		CodeReference: core.StringPtr(codeReference),
		Name: core.StringPtr(name),
		Runtime: core.StringPtr(runtime),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateFunctionOptions) SetProjectID(projectID string) *CreateFunctionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetCodeReference : Allow user to set CodeReference
func (_options *CreateFunctionOptions) SetCodeReference(codeReference string) *CreateFunctionOptions {
	_options.CodeReference = core.StringPtr(codeReference)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateFunctionOptions) SetName(name string) *CreateFunctionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRuntime : Allow user to set Runtime
func (_options *CreateFunctionOptions) SetRuntime(runtime string) *CreateFunctionOptions {
	_options.Runtime = core.StringPtr(runtime)
	return _options
}

// SetCodeBinary : Allow user to set CodeBinary
func (_options *CreateFunctionOptions) SetCodeBinary(codeBinary bool) *CreateFunctionOptions {
	_options.CodeBinary = core.BoolPtr(codeBinary)
	return _options
}

// SetCodeMain : Allow user to set CodeMain
func (_options *CreateFunctionOptions) SetCodeMain(codeMain string) *CreateFunctionOptions {
	_options.CodeMain = core.StringPtr(codeMain)
	return _options
}

// SetCodeSecret : Allow user to set CodeSecret
func (_options *CreateFunctionOptions) SetCodeSecret(codeSecret string) *CreateFunctionOptions {
	_options.CodeSecret = core.StringPtr(codeSecret)
	return _options
}

// SetManagedDomainMappings : Allow user to set ManagedDomainMappings
func (_options *CreateFunctionOptions) SetManagedDomainMappings(managedDomainMappings string) *CreateFunctionOptions {
	_options.ManagedDomainMappings = core.StringPtr(managedDomainMappings)
	return _options
}

// SetRunEnvVariables : Allow user to set RunEnvVariables
func (_options *CreateFunctionOptions) SetRunEnvVariables(runEnvVariables []EnvVarPrototype) *CreateFunctionOptions {
	_options.RunEnvVariables = runEnvVariables
	return _options
}

// SetScaleConcurrency : Allow user to set ScaleConcurrency
func (_options *CreateFunctionOptions) SetScaleConcurrency(scaleConcurrency int64) *CreateFunctionOptions {
	_options.ScaleConcurrency = core.Int64Ptr(scaleConcurrency)
	return _options
}

// SetScaleCpuLimit : Allow user to set ScaleCpuLimit
func (_options *CreateFunctionOptions) SetScaleCpuLimit(scaleCpuLimit string) *CreateFunctionOptions {
	_options.ScaleCpuLimit = core.StringPtr(scaleCpuLimit)
	return _options
}

// SetScaleDownDelay : Allow user to set ScaleDownDelay
func (_options *CreateFunctionOptions) SetScaleDownDelay(scaleDownDelay int64) *CreateFunctionOptions {
	_options.ScaleDownDelay = core.Int64Ptr(scaleDownDelay)
	return _options
}

// SetScaleMaxExecutionTime : Allow user to set ScaleMaxExecutionTime
func (_options *CreateFunctionOptions) SetScaleMaxExecutionTime(scaleMaxExecutionTime int64) *CreateFunctionOptions {
	_options.ScaleMaxExecutionTime = core.Int64Ptr(scaleMaxExecutionTime)
	return _options
}

// SetScaleMemoryLimit : Allow user to set ScaleMemoryLimit
func (_options *CreateFunctionOptions) SetScaleMemoryLimit(scaleMemoryLimit string) *CreateFunctionOptions {
	_options.ScaleMemoryLimit = core.StringPtr(scaleMemoryLimit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateFunctionOptions) SetHeaders(param map[string]string) *CreateFunctionOptions {
	options.Headers = param
	return options
}

// CreateJobOptions : The CreateJob options.
type CreateJobOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference" validate:"required"`

	// The name of the job. Use a name that is unique within the project.
	Name *string `json:"name" validate:"required"`

	// The name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the job / job runs will be created but submitted job runs will fail, until this property is
	// provided, too. This property must not be set on a job run, which references a job template.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Set arguments for the job that are passed to start job run containers. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments,omitempty"`

	// The user ID (UID) to run the job.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Set commands for the job that are passed to start job run containers. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands,omitempty"`

	// Optional references to config maps, secrets or literal values.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
	// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
	// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
	RunMode *string `json:"run_mode,omitempty"`

	// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
	// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or a secrets.
	RunVolumeMounts []VolumeMountPrototype `json:"run_volume_mounts,omitempty"`

	// Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges,
	// such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number
	// of unique array indices that you specify with this parameter determines the number of job instances to run.
	ScaleArraySpec *string `json:"scale_array_spec,omitempty"`

	// Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit,omitempty"`

	// The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is
	// `task`.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// The number of times to rerun an instance of the job before the job is marked as failed. This property can only be
	// specified if `run_mode` is `task`.
	ScaleRetryLimit *int64 `json:"scale_retry_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateJobOptions.RunMode property.
// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
const (
	CreateJobOptions_RunMode_Daemon = "daemon"
	CreateJobOptions_RunMode_Task = "task"
)

// Constants associated with the CreateJobOptions.RunServiceAccount property.
// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
const (
	CreateJobOptions_RunServiceAccount_Default = "default"
	CreateJobOptions_RunServiceAccount_Manager = "manager"
	CreateJobOptions_RunServiceAccount_None = "none"
	CreateJobOptions_RunServiceAccount_Reader = "reader"
	CreateJobOptions_RunServiceAccount_Writer = "writer"
)

// NewCreateJobOptions : Instantiate CreateJobOptions
func (*CodeEngineV2) NewCreateJobOptions(projectID string, imageReference string, name string) *CreateJobOptions {
	return &CreateJobOptions{
		ProjectID: core.StringPtr(projectID),
		ImageReference: core.StringPtr(imageReference),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateJobOptions) SetProjectID(projectID string) *CreateJobOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetImageReference : Allow user to set ImageReference
func (_options *CreateJobOptions) SetImageReference(imageReference string) *CreateJobOptions {
	_options.ImageReference = core.StringPtr(imageReference)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateJobOptions) SetName(name string) *CreateJobOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetImageSecret : Allow user to set ImageSecret
func (_options *CreateJobOptions) SetImageSecret(imageSecret string) *CreateJobOptions {
	_options.ImageSecret = core.StringPtr(imageSecret)
	return _options
}

// SetRunArguments : Allow user to set RunArguments
func (_options *CreateJobOptions) SetRunArguments(runArguments []string) *CreateJobOptions {
	_options.RunArguments = runArguments
	return _options
}

// SetRunAsUser : Allow user to set RunAsUser
func (_options *CreateJobOptions) SetRunAsUser(runAsUser int64) *CreateJobOptions {
	_options.RunAsUser = core.Int64Ptr(runAsUser)
	return _options
}

// SetRunCommands : Allow user to set RunCommands
func (_options *CreateJobOptions) SetRunCommands(runCommands []string) *CreateJobOptions {
	_options.RunCommands = runCommands
	return _options
}

// SetRunEnvVariables : Allow user to set RunEnvVariables
func (_options *CreateJobOptions) SetRunEnvVariables(runEnvVariables []EnvVarPrototype) *CreateJobOptions {
	_options.RunEnvVariables = runEnvVariables
	return _options
}

// SetRunMode : Allow user to set RunMode
func (_options *CreateJobOptions) SetRunMode(runMode string) *CreateJobOptions {
	_options.RunMode = core.StringPtr(runMode)
	return _options
}

// SetRunServiceAccount : Allow user to set RunServiceAccount
func (_options *CreateJobOptions) SetRunServiceAccount(runServiceAccount string) *CreateJobOptions {
	_options.RunServiceAccount = core.StringPtr(runServiceAccount)
	return _options
}

// SetRunVolumeMounts : Allow user to set RunVolumeMounts
func (_options *CreateJobOptions) SetRunVolumeMounts(runVolumeMounts []VolumeMountPrototype) *CreateJobOptions {
	_options.RunVolumeMounts = runVolumeMounts
	return _options
}

// SetScaleArraySpec : Allow user to set ScaleArraySpec
func (_options *CreateJobOptions) SetScaleArraySpec(scaleArraySpec string) *CreateJobOptions {
	_options.ScaleArraySpec = core.StringPtr(scaleArraySpec)
	return _options
}

// SetScaleCpuLimit : Allow user to set ScaleCpuLimit
func (_options *CreateJobOptions) SetScaleCpuLimit(scaleCpuLimit string) *CreateJobOptions {
	_options.ScaleCpuLimit = core.StringPtr(scaleCpuLimit)
	return _options
}

// SetScaleEphemeralStorageLimit : Allow user to set ScaleEphemeralStorageLimit
func (_options *CreateJobOptions) SetScaleEphemeralStorageLimit(scaleEphemeralStorageLimit string) *CreateJobOptions {
	_options.ScaleEphemeralStorageLimit = core.StringPtr(scaleEphemeralStorageLimit)
	return _options
}

// SetScaleMaxExecutionTime : Allow user to set ScaleMaxExecutionTime
func (_options *CreateJobOptions) SetScaleMaxExecutionTime(scaleMaxExecutionTime int64) *CreateJobOptions {
	_options.ScaleMaxExecutionTime = core.Int64Ptr(scaleMaxExecutionTime)
	return _options
}

// SetScaleMemoryLimit : Allow user to set ScaleMemoryLimit
func (_options *CreateJobOptions) SetScaleMemoryLimit(scaleMemoryLimit string) *CreateJobOptions {
	_options.ScaleMemoryLimit = core.StringPtr(scaleMemoryLimit)
	return _options
}

// SetScaleRetryLimit : Allow user to set ScaleRetryLimit
func (_options *CreateJobOptions) SetScaleRetryLimit(scaleRetryLimit int64) *CreateJobOptions {
	_options.ScaleRetryLimit = core.Int64Ptr(scaleRetryLimit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateJobOptions) SetHeaders(param map[string]string) *CreateJobOptions {
	options.Headers = param
	return options
}

// CreateJobRunOptions : The CreateJobRun options.
type CreateJobRunOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference,omitempty"`

	// The name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the job / job runs will be created but submitted job runs will fail, until this property is
	// provided, too. This property must not be set on a job run, which references a job template.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Optional name of the job on which this job run is based on. If specified, the job run will inherit the configuration
	// of the referenced job.
	JobName *string `json:"job_name,omitempty"`

	// The name of the job. Use a name that is unique within the project.
	Name *string `json:"name,omitempty"`

	// Set arguments for the job that are passed to start job run containers. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments,omitempty"`

	// The user ID (UID) to run the job.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Set commands for the job that are passed to start job run containers. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands,omitempty"`

	// Optional references to config maps, secrets or literal values.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
	// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
	// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
	RunMode *string `json:"run_mode,omitempty"`

	// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
	// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or a secrets.
	RunVolumeMounts []VolumeMountPrototype `json:"run_volume_mounts,omitempty"`

	// Optional value to override the JOB_ARRAY_SIZE environment variable for a job run.
	ScaleArraySizeVariableOverride *int64 `json:"scale_array_size_variable_override,omitempty"`

	// Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges,
	// such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number
	// of unique array indices that you specify with this parameter determines the number of job instances to run.
	ScaleArraySpec *string `json:"scale_array_spec,omitempty"`

	// Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit,omitempty"`

	// The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is
	// `task`.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// The number of times to rerun an instance of the job before the job is marked as failed. This property can only be
	// specified if `run_mode` is `task`.
	ScaleRetryLimit *int64 `json:"scale_retry_limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateJobRunOptions.RunMode property.
// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
const (
	CreateJobRunOptions_RunMode_Daemon = "daemon"
	CreateJobRunOptions_RunMode_Task = "task"
)

// Constants associated with the CreateJobRunOptions.RunServiceAccount property.
// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
const (
	CreateJobRunOptions_RunServiceAccount_Default = "default"
	CreateJobRunOptions_RunServiceAccount_Manager = "manager"
	CreateJobRunOptions_RunServiceAccount_None = "none"
	CreateJobRunOptions_RunServiceAccount_Reader = "reader"
	CreateJobRunOptions_RunServiceAccount_Writer = "writer"
)

// NewCreateJobRunOptions : Instantiate CreateJobRunOptions
func (*CodeEngineV2) NewCreateJobRunOptions(projectID string) *CreateJobRunOptions {
	return &CreateJobRunOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateJobRunOptions) SetProjectID(projectID string) *CreateJobRunOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetImageReference : Allow user to set ImageReference
func (_options *CreateJobRunOptions) SetImageReference(imageReference string) *CreateJobRunOptions {
	_options.ImageReference = core.StringPtr(imageReference)
	return _options
}

// SetImageSecret : Allow user to set ImageSecret
func (_options *CreateJobRunOptions) SetImageSecret(imageSecret string) *CreateJobRunOptions {
	_options.ImageSecret = core.StringPtr(imageSecret)
	return _options
}

// SetJobName : Allow user to set JobName
func (_options *CreateJobRunOptions) SetJobName(jobName string) *CreateJobRunOptions {
	_options.JobName = core.StringPtr(jobName)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateJobRunOptions) SetName(name string) *CreateJobRunOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetRunArguments : Allow user to set RunArguments
func (_options *CreateJobRunOptions) SetRunArguments(runArguments []string) *CreateJobRunOptions {
	_options.RunArguments = runArguments
	return _options
}

// SetRunAsUser : Allow user to set RunAsUser
func (_options *CreateJobRunOptions) SetRunAsUser(runAsUser int64) *CreateJobRunOptions {
	_options.RunAsUser = core.Int64Ptr(runAsUser)
	return _options
}

// SetRunCommands : Allow user to set RunCommands
func (_options *CreateJobRunOptions) SetRunCommands(runCommands []string) *CreateJobRunOptions {
	_options.RunCommands = runCommands
	return _options
}

// SetRunEnvVariables : Allow user to set RunEnvVariables
func (_options *CreateJobRunOptions) SetRunEnvVariables(runEnvVariables []EnvVarPrototype) *CreateJobRunOptions {
	_options.RunEnvVariables = runEnvVariables
	return _options
}

// SetRunMode : Allow user to set RunMode
func (_options *CreateJobRunOptions) SetRunMode(runMode string) *CreateJobRunOptions {
	_options.RunMode = core.StringPtr(runMode)
	return _options
}

// SetRunServiceAccount : Allow user to set RunServiceAccount
func (_options *CreateJobRunOptions) SetRunServiceAccount(runServiceAccount string) *CreateJobRunOptions {
	_options.RunServiceAccount = core.StringPtr(runServiceAccount)
	return _options
}

// SetRunVolumeMounts : Allow user to set RunVolumeMounts
func (_options *CreateJobRunOptions) SetRunVolumeMounts(runVolumeMounts []VolumeMountPrototype) *CreateJobRunOptions {
	_options.RunVolumeMounts = runVolumeMounts
	return _options
}

// SetScaleArraySizeVariableOverride : Allow user to set ScaleArraySizeVariableOverride
func (_options *CreateJobRunOptions) SetScaleArraySizeVariableOverride(scaleArraySizeVariableOverride int64) *CreateJobRunOptions {
	_options.ScaleArraySizeVariableOverride = core.Int64Ptr(scaleArraySizeVariableOverride)
	return _options
}

// SetScaleArraySpec : Allow user to set ScaleArraySpec
func (_options *CreateJobRunOptions) SetScaleArraySpec(scaleArraySpec string) *CreateJobRunOptions {
	_options.ScaleArraySpec = core.StringPtr(scaleArraySpec)
	return _options
}

// SetScaleCpuLimit : Allow user to set ScaleCpuLimit
func (_options *CreateJobRunOptions) SetScaleCpuLimit(scaleCpuLimit string) *CreateJobRunOptions {
	_options.ScaleCpuLimit = core.StringPtr(scaleCpuLimit)
	return _options
}

// SetScaleEphemeralStorageLimit : Allow user to set ScaleEphemeralStorageLimit
func (_options *CreateJobRunOptions) SetScaleEphemeralStorageLimit(scaleEphemeralStorageLimit string) *CreateJobRunOptions {
	_options.ScaleEphemeralStorageLimit = core.StringPtr(scaleEphemeralStorageLimit)
	return _options
}

// SetScaleMaxExecutionTime : Allow user to set ScaleMaxExecutionTime
func (_options *CreateJobRunOptions) SetScaleMaxExecutionTime(scaleMaxExecutionTime int64) *CreateJobRunOptions {
	_options.ScaleMaxExecutionTime = core.Int64Ptr(scaleMaxExecutionTime)
	return _options
}

// SetScaleMemoryLimit : Allow user to set ScaleMemoryLimit
func (_options *CreateJobRunOptions) SetScaleMemoryLimit(scaleMemoryLimit string) *CreateJobRunOptions {
	_options.ScaleMemoryLimit = core.StringPtr(scaleMemoryLimit)
	return _options
}

// SetScaleRetryLimit : Allow user to set ScaleRetryLimit
func (_options *CreateJobRunOptions) SetScaleRetryLimit(scaleRetryLimit int64) *CreateJobRunOptions {
	_options.ScaleRetryLimit = core.Int64Ptr(scaleRetryLimit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateJobRunOptions) SetHeaders(param map[string]string) *CreateJobRunOptions {
	options.Headers = param
	return options
}

// CreateProjectOptions : The CreateProject options.
type CreateProjectOptions struct {
	// The name of the project.
	Name *string `json:"name" validate:"required"`

	// Optional ID of the resource group for your project deployment. If this field is not defined, the default resource
	// group of the account will be used.
	ResourceGroupID *string `json:"resource_group_id,omitempty"`

	// Optional list of labels to assign to your project. Tags are not part of the project resource that is returned by the
	// server, but can be obtained and managed through the Global Tagging API in IBM Cloud. Find more information on
	// [Global Tagging API docs](https://cloud.ibm.com/apidocs/tagging).
	Tags []string `json:"tags,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateProjectOptions : Instantiate CreateProjectOptions
func (*CodeEngineV2) NewCreateProjectOptions(name string) *CreateProjectOptions {
	return &CreateProjectOptions{
		Name: core.StringPtr(name),
	}
}

// SetName : Allow user to set Name
func (_options *CreateProjectOptions) SetName(name string) *CreateProjectOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *CreateProjectOptions) SetResourceGroupID(resourceGroupID string) *CreateProjectOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateProjectOptions) SetTags(tags []string) *CreateProjectOptions {
	_options.Tags = tags
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateProjectOptions) SetHeaders(param map[string]string) *CreateProjectOptions {
	options.Headers = param
	return options
}

// CreateSecretOptions : The CreateSecret options.
type CreateSecretOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Specify the format of the secret. The format of the secret will determine how the secret is used.
	Format *string `json:"format" validate:"required"`

	// The name of the secret.
	Name *string `json:"name" validate:"required"`

	// Data container that allows to specify config parameters and their values as a key-value map. Each key field must
	// consist of alphanumeric characters, `-`, `_` or `.` and must not exceed a max length of 253 characters. Each value
	// field can consists of any character and must not exceed a max length of 1048576 characters.
	Data SecretDataIntf `json:"data,omitempty"`

	// Properties for Service Access Secrets.
	ServiceAccess *ServiceAccessSecretPrototypeProps `json:"service_access,omitempty"`

	// Properties for the IBM Cloud Operator Secrets.
	ServiceOperator *OperatorSecretPrototypeProps `json:"service_operator,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateSecretOptions.Format property.
// Specify the format of the secret. The format of the secret will determine how the secret is used.
const (
	CreateSecretOptions_Format_BasicAuth = "basic_auth"
	CreateSecretOptions_Format_Generic = "generic"
	CreateSecretOptions_Format_Other = "other"
	CreateSecretOptions_Format_Registry = "registry"
	CreateSecretOptions_Format_ServiceAccess = "service_access"
	CreateSecretOptions_Format_ServiceOperator = "service_operator"
	CreateSecretOptions_Format_SshAuth = "ssh_auth"
	CreateSecretOptions_Format_Tls = "tls"
)

// NewCreateSecretOptions : Instantiate CreateSecretOptions
func (*CodeEngineV2) NewCreateSecretOptions(projectID string, format string, name string) *CreateSecretOptions {
	return &CreateSecretOptions{
		ProjectID: core.StringPtr(projectID),
		Format: core.StringPtr(format),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *CreateSecretOptions) SetProjectID(projectID string) *CreateSecretOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *CreateSecretOptions) SetFormat(format string) *CreateSecretOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateSecretOptions) SetName(name string) *CreateSecretOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetData : Allow user to set Data
func (_options *CreateSecretOptions) SetData(data SecretDataIntf) *CreateSecretOptions {
	_options.Data = data
	return _options
}

// SetServiceAccess : Allow user to set ServiceAccess
func (_options *CreateSecretOptions) SetServiceAccess(serviceAccess *ServiceAccessSecretPrototypeProps) *CreateSecretOptions {
	_options.ServiceAccess = serviceAccess
	return _options
}

// SetServiceOperator : Allow user to set ServiceOperator
func (_options *CreateSecretOptions) SetServiceOperator(serviceOperator *OperatorSecretPrototypeProps) *CreateSecretOptions {
	_options.ServiceOperator = serviceOperator
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSecretOptions) SetHeaders(param map[string]string) *CreateSecretOptions {
	options.Headers = param
	return options
}

// DeleteAllowedOutboundDestinationOptions : The DeleteAllowedOutboundDestination options.
type DeleteAllowedOutboundDestinationOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your allowed outbound destination.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAllowedOutboundDestinationOptions : Instantiate DeleteAllowedOutboundDestinationOptions
func (*CodeEngineV2) NewDeleteAllowedOutboundDestinationOptions(projectID string, name string) *DeleteAllowedOutboundDestinationOptions {
	return &DeleteAllowedOutboundDestinationOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteAllowedOutboundDestinationOptions) SetProjectID(projectID string) *DeleteAllowedOutboundDestinationOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteAllowedOutboundDestinationOptions) SetName(name string) *DeleteAllowedOutboundDestinationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAllowedOutboundDestinationOptions) SetHeaders(param map[string]string) *DeleteAllowedOutboundDestinationOptions {
	options.Headers = param
	return options
}

// DeleteAppOptions : The DeleteApp options.
type DeleteAppOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAppOptions : Instantiate DeleteAppOptions
func (*CodeEngineV2) NewDeleteAppOptions(projectID string, name string) *DeleteAppOptions {
	return &DeleteAppOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteAppOptions) SetProjectID(projectID string) *DeleteAppOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteAppOptions) SetName(name string) *DeleteAppOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAppOptions) SetHeaders(param map[string]string) *DeleteAppOptions {
	options.Headers = param
	return options
}

// DeleteAppRevisionOptions : The DeleteAppRevision options.
type DeleteAppRevisionOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	AppName *string `json:"app_name" validate:"required,ne="`

	// The name of your application revision.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteAppRevisionOptions : Instantiate DeleteAppRevisionOptions
func (*CodeEngineV2) NewDeleteAppRevisionOptions(projectID string, appName string, name string) *DeleteAppRevisionOptions {
	return &DeleteAppRevisionOptions{
		ProjectID: core.StringPtr(projectID),
		AppName: core.StringPtr(appName),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteAppRevisionOptions) SetProjectID(projectID string) *DeleteAppRevisionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetAppName : Allow user to set AppName
func (_options *DeleteAppRevisionOptions) SetAppName(appName string) *DeleteAppRevisionOptions {
	_options.AppName = core.StringPtr(appName)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteAppRevisionOptions) SetName(name string) *DeleteAppRevisionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteAppRevisionOptions) SetHeaders(param map[string]string) *DeleteAppRevisionOptions {
	options.Headers = param
	return options
}

// DeleteBindingOptions : The DeleteBinding options.
type DeleteBindingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The id of your binding.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteBindingOptions : Instantiate DeleteBindingOptions
func (*CodeEngineV2) NewDeleteBindingOptions(projectID string, id string) *DeleteBindingOptions {
	return &DeleteBindingOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteBindingOptions) SetProjectID(projectID string) *DeleteBindingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *DeleteBindingOptions) SetID(id string) *DeleteBindingOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteBindingOptions) SetHeaders(param map[string]string) *DeleteBindingOptions {
	options.Headers = param
	return options
}

// DeleteBuildOptions : The DeleteBuild options.
type DeleteBuildOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your build.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteBuildOptions : Instantiate DeleteBuildOptions
func (*CodeEngineV2) NewDeleteBuildOptions(projectID string, name string) *DeleteBuildOptions {
	return &DeleteBuildOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteBuildOptions) SetProjectID(projectID string) *DeleteBuildOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteBuildOptions) SetName(name string) *DeleteBuildOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteBuildOptions) SetHeaders(param map[string]string) *DeleteBuildOptions {
	options.Headers = param
	return options
}

// DeleteBuildRunOptions : The DeleteBuildRun options.
type DeleteBuildRunOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your build run.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteBuildRunOptions : Instantiate DeleteBuildRunOptions
func (*CodeEngineV2) NewDeleteBuildRunOptions(projectID string, name string) *DeleteBuildRunOptions {
	return &DeleteBuildRunOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteBuildRunOptions) SetProjectID(projectID string) *DeleteBuildRunOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteBuildRunOptions) SetName(name string) *DeleteBuildRunOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteBuildRunOptions) SetHeaders(param map[string]string) *DeleteBuildRunOptions {
	options.Headers = param
	return options
}

// DeleteConfigMapOptions : The DeleteConfigMap options.
type DeleteConfigMapOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your configmap.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteConfigMapOptions : Instantiate DeleteConfigMapOptions
func (*CodeEngineV2) NewDeleteConfigMapOptions(projectID string, name string) *DeleteConfigMapOptions {
	return &DeleteConfigMapOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteConfigMapOptions) SetProjectID(projectID string) *DeleteConfigMapOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteConfigMapOptions) SetName(name string) *DeleteConfigMapOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteConfigMapOptions) SetHeaders(param map[string]string) *DeleteConfigMapOptions {
	options.Headers = param
	return options
}

// DeleteDomainMappingOptions : The DeleteDomainMapping options.
type DeleteDomainMappingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your domain mapping.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteDomainMappingOptions : Instantiate DeleteDomainMappingOptions
func (*CodeEngineV2) NewDeleteDomainMappingOptions(projectID string, name string) *DeleteDomainMappingOptions {
	return &DeleteDomainMappingOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteDomainMappingOptions) SetProjectID(projectID string) *DeleteDomainMappingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteDomainMappingOptions) SetName(name string) *DeleteDomainMappingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteDomainMappingOptions) SetHeaders(param map[string]string) *DeleteDomainMappingOptions {
	options.Headers = param
	return options
}

// DeleteFunctionOptions : The DeleteFunction options.
type DeleteFunctionOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your function.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteFunctionOptions : Instantiate DeleteFunctionOptions
func (*CodeEngineV2) NewDeleteFunctionOptions(projectID string, name string) *DeleteFunctionOptions {
	return &DeleteFunctionOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteFunctionOptions) SetProjectID(projectID string) *DeleteFunctionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteFunctionOptions) SetName(name string) *DeleteFunctionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteFunctionOptions) SetHeaders(param map[string]string) *DeleteFunctionOptions {
	options.Headers = param
	return options
}

// DeleteJobOptions : The DeleteJob options.
type DeleteJobOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your job.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteJobOptions : Instantiate DeleteJobOptions
func (*CodeEngineV2) NewDeleteJobOptions(projectID string, name string) *DeleteJobOptions {
	return &DeleteJobOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteJobOptions) SetProjectID(projectID string) *DeleteJobOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteJobOptions) SetName(name string) *DeleteJobOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteJobOptions) SetHeaders(param map[string]string) *DeleteJobOptions {
	options.Headers = param
	return options
}

// DeleteJobRunOptions : The DeleteJobRun options.
type DeleteJobRunOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your job run.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteJobRunOptions : Instantiate DeleteJobRunOptions
func (*CodeEngineV2) NewDeleteJobRunOptions(projectID string, name string) *DeleteJobRunOptions {
	return &DeleteJobRunOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteJobRunOptions) SetProjectID(projectID string) *DeleteJobRunOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteJobRunOptions) SetName(name string) *DeleteJobRunOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteJobRunOptions) SetHeaders(param map[string]string) *DeleteJobRunOptions {
	options.Headers = param
	return options
}

// DeleteProjectOptions : The DeleteProject options.
type DeleteProjectOptions struct {
	// The ID of the project.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteProjectOptions : Instantiate DeleteProjectOptions
func (*CodeEngineV2) NewDeleteProjectOptions(id string) *DeleteProjectOptions {
	return &DeleteProjectOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *DeleteProjectOptions) SetID(id string) *DeleteProjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteProjectOptions) SetHeaders(param map[string]string) *DeleteProjectOptions {
	options.Headers = param
	return options
}

// DeleteSecretOptions : The DeleteSecret options.
type DeleteSecretOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your secret.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteSecretOptions : Instantiate DeleteSecretOptions
func (*CodeEngineV2) NewDeleteSecretOptions(projectID string, name string) *DeleteSecretOptions {
	return &DeleteSecretOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *DeleteSecretOptions) SetProjectID(projectID string) *DeleteSecretOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *DeleteSecretOptions) SetName(name string) *DeleteSecretOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSecretOptions) SetHeaders(param map[string]string) *DeleteSecretOptions {
	options.Headers = param
	return options
}

// DomainMapping : Response model for domain mapping definitions.
type DomainMapping struct {
	// The value of the CNAME record that must be configured in the DNS settings of the domain, to route traffic properly
	// to the target Code Engine region.
	CnameTarget *string `json:"cname_target,omitempty"`

	// A reference to another component.
	Component *ComponentRef `json:"component" validate:"required"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The version of the domain mapping instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// When you provision a new domain mapping, a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the domain mapping.
	Name *string `json:"name" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the Code Engine resource.
	ResourceType *string `json:"resource_type,omitempty"`

	// The current status of the domain mapping.
	Status *string `json:"status,omitempty"`

	// The detailed status of the domain mapping.
	StatusDetails *DomainMappingStatus `json:"status_details,omitempty"`

	// The name of the TLS secret that includes the certificate and private key of this domain mapping.
	TlsSecret *string `json:"tls_secret" validate:"required"`

	// Specifies whether the domain mapping is managed by the user or by Code Engine.
	UserManaged *bool `json:"user_managed,omitempty"`

	// Specifies whether the domain mapping is reachable through the public internet, or private IBM network, or only
	// through other components within the same Code Engine project.
	Visibility *string `json:"visibility,omitempty"`
}

// Constants associated with the DomainMapping.ResourceType property.
// The type of the Code Engine resource.
const (
	DomainMapping_ResourceType_DomainMappingV2 = "domain_mapping_v2"
)

// Constants associated with the DomainMapping.Status property.
// The current status of the domain mapping.
const (
	DomainMapping_Status_Deploying = "deploying"
	DomainMapping_Status_Failed = "failed"
	DomainMapping_Status_Ready = "ready"
)

// Constants associated with the DomainMapping.Visibility property.
// Specifies whether the domain mapping is reachable through the public internet, or private IBM network, or only
// through other components within the same Code Engine project.
const (
	DomainMapping_Visibility_Custom = "custom"
	DomainMapping_Visibility_Private = "private"
	DomainMapping_Visibility_Project = "project"
	DomainMapping_Visibility_Public = "public"
)

// UnmarshalDomainMapping unmarshals an instance of DomainMapping from the specified map of raw messages.
func UnmarshalDomainMapping(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DomainMapping)
	err = core.UnmarshalPrimitive(m, "cname_target", &obj.CnameTarget)
	if err != nil {
		err = core.SDKErrorf(err, "", "cname_target-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "component", &obj.Component, UnmarshalComponentRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "component-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalDomainMappingStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tls_secret", &obj.TlsSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "tls_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_managed", &obj.UserManaged)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_managed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "visibility", &obj.Visibility)
	if err != nil {
		err = core.SDKErrorf(err, "", "visibility-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DomainMappingList : Contains a list of domain mappings and pagination information.
type DomainMappingList struct {
	// List of all domain mappings.
	DomainMappings []DomainMapping `json:"domain_mappings" validate:"required"`

	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalDomainMappingList unmarshals an instance of DomainMappingList from the specified map of raw messages.
func UnmarshalDomainMappingList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DomainMappingList)
	err = core.UnmarshalModel(m, "domain_mappings", &obj.DomainMappings, UnmarshalDomainMapping)
	if err != nil {
		err = core.SDKErrorf(err, "", "domain_mappings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *DomainMappingList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// DomainMappingPatch : Patch a domain mappings object.
type DomainMappingPatch struct {
	// A reference to another component.
	Component *ComponentRef `json:"component,omitempty"`

	// The name of the TLS secret that includes the certificate and private key of this domain mapping.
	TlsSecret *string `json:"tls_secret,omitempty"`
}

// UnmarshalDomainMappingPatch unmarshals an instance of DomainMappingPatch from the specified map of raw messages.
func UnmarshalDomainMappingPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DomainMappingPatch)
	err = core.UnmarshalModel(m, "component", &obj.Component, UnmarshalComponentRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "component-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tls_secret", &obj.TlsSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "tls_secret-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the DomainMappingPatch
func (domainMappingPatch *DomainMappingPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(domainMappingPatch.Component) {
		_patch["component"] = domainMappingPatch.Component.asPatch()
	}
	if !core.IsNil(domainMappingPatch.TlsSecret) {
		_patch["tls_secret"] = domainMappingPatch.TlsSecret
	}

	return
}

// DomainMappingStatus : The detailed status of the domain mapping.
type DomainMappingStatus struct {
	// Optional information to provide more context in case of a 'failed' or 'warning' status.
	Reason *string `json:"reason,omitempty"`
}

// Constants associated with the DomainMappingStatus.Reason property.
// Optional information to provide more context in case of a 'failed' or 'warning' status.
const (
	DomainMappingStatus_Reason_AppRefFailed = "app_ref_failed"
	DomainMappingStatus_Reason_Deploying = "deploying"
	DomainMappingStatus_Reason_DomainAlreadyClaimed = "domain_already_claimed"
	DomainMappingStatus_Reason_Failed = "failed"
	DomainMappingStatus_Reason_FailedReconcileIngress = "failed_reconcile_ingress"
	DomainMappingStatus_Reason_Ready = "ready"
)

// UnmarshalDomainMappingStatus unmarshals an instance of DomainMappingStatus from the specified map of raw messages.
func UnmarshalDomainMappingStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DomainMappingStatus)
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvVar : Response model for environment variables.
type EnvVar struct {
	// The key to reference as environment variable.
	Key *string `json:"key,omitempty"`

	// The name of the environment variable.
	Name *string `json:"name,omitempty"`

	// A prefix that can be added to all keys of a full secret or config map reference.
	Prefix *string `json:"prefix,omitempty"`

	// The name of the secret or config map.
	Reference *string `json:"reference,omitempty"`

	// Specify the type of the environment variable.
	Type *string `json:"type" validate:"required"`

	// The literal value of the environment variable.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the EnvVar.Type property.
// Specify the type of the environment variable.
const (
	EnvVar_Type_ConfigMapFullReference = "config_map_full_reference"
	EnvVar_Type_ConfigMapKeyReference = "config_map_key_reference"
	EnvVar_Type_Literal = "literal"
	EnvVar_Type_SecretFullReference = "secret_full_reference"
	EnvVar_Type_SecretKeyReference = "secret_key_reference"
)

// UnmarshalEnvVar unmarshals an instance of EnvVar from the specified map of raw messages.
func UnmarshalEnvVar(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvVar)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "prefix", &obj.Prefix)
	if err != nil {
		err = core.SDKErrorf(err, "", "prefix-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvVarPrototype : Prototype model for environment variables.
type EnvVarPrototype struct {
	// The key to reference as environment variable.
	Key *string `json:"key,omitempty"`

	// The name of the environment variable.
	Name *string `json:"name,omitempty"`

	// A prefix that can be added to all keys of a full secret or config map reference.
	Prefix *string `json:"prefix,omitempty"`

	// The name of the secret or config map.
	Reference *string `json:"reference,omitempty"`

	// Specify the type of the environment variable.
	Type *string `json:"type,omitempty"`

	// The literal value of the environment variable.
	Value *string `json:"value,omitempty"`
}

// Constants associated with the EnvVarPrototype.Type property.
// Specify the type of the environment variable.
const (
	EnvVarPrototype_Type_ConfigMapFullReference = "config_map_full_reference"
	EnvVarPrototype_Type_ConfigMapKeyReference = "config_map_key_reference"
	EnvVarPrototype_Type_Literal = "literal"
	EnvVarPrototype_Type_SecretFullReference = "secret_full_reference"
	EnvVarPrototype_Type_SecretKeyReference = "secret_key_reference"
)

// UnmarshalEnvVarPrototype unmarshals an instance of EnvVarPrototype from the specified map of raw messages.
func UnmarshalEnvVarPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvVarPrototype)
	err = core.UnmarshalPrimitive(m, "key", &obj.Key)
	if err != nil {
		err = core.SDKErrorf(err, "", "key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "prefix", &obj.Prefix)
	if err != nil {
		err = core.SDKErrorf(err, "", "prefix-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the EnvVarPrototype
func (envVarPrototype *EnvVarPrototype) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(envVarPrototype.Key) {
		_patch["key"] = envVarPrototype.Key
	}
	if !core.IsNil(envVarPrototype.Name) {
		_patch["name"] = envVarPrototype.Name
	}
	if !core.IsNil(envVarPrototype.Prefix) {
		_patch["prefix"] = envVarPrototype.Prefix
	}
	if !core.IsNil(envVarPrototype.Reference) {
		_patch["reference"] = envVarPrototype.Reference
	}
	if !core.IsNil(envVarPrototype.Type) {
		_patch["type"] = envVarPrototype.Type
	}
	if !core.IsNil(envVarPrototype.Value) {
		_patch["value"] = envVarPrototype.Value
	}

	return
}

// Function : Function is the response model for function resources.
type Function struct {
	// Specifies whether the code is binary or not. Defaults to false when `code_reference` is set to a data URL. When
	// `code_reference` is set to a code bundle URL, this field is always true.
	CodeBinary *bool `json:"code_binary" validate:"required"`

	// Specifies the name of the function that should be invoked.
	CodeMain *string `json:"code_main,omitempty"`

	// Specifies either a reference to a code bundle or the source code itself. To specify the source code, use the data
	// URL scheme and include the source code as base64 encoded. The data URL scheme is defined in [RFC
	// 2397](https://tools.ietf.org/html/rfc2397).
	CodeReference *string `json:"code_reference" validate:"required"`

	// The name of the secret that is used to access the specified `code_reference`. The secret is used to authenticate
	// with a non-public endpoint that is specified as`code_reference`.
	CodeSecret *string `json:"code_secret,omitempty"`

	// References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as
	// environment variables in the function.
	ComputedEnvVariables []EnvVar `json:"computed_env_variables,omitempty"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// URL to invoke the function.
	Endpoint *string `json:"endpoint,omitempty"`

	// URL to function that is only visible within the project.
	EndpointInternal *string `json:"endpoint_internal,omitempty"`

	// The version of the function instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// When you provision a new function, a relative URL path is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values
	// are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports
	// function private visibility.
	ManagedDomainMappings *string `json:"managed_domain_mappings" validate:"required"`

	// The name of the function.
	Name *string `json:"name" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the function.
	ResourceType *string `json:"resource_type,omitempty"`

	// References to config maps, secrets or literal values, which are defined by the function owner and are exposed as
	// environment variables in the function.
	RunEnvVariables []EnvVar `json:"run_env_variables" validate:"required"`

	// The managed runtime used to execute the injected code.
	Runtime *string `json:"runtime" validate:"required"`

	// Number of parallel requests handled by a single instance, supported only by Node.js, default is `1`.
	ScaleConcurrency *int64 `json:"scale_concurrency" validate:"required"`

	// Optional amount of CPU set for the instance of the function. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit" validate:"required"`

	// Optional amount of time in seconds that delays the scale down behavior for a function.
	ScaleDownDelay *int64 `json:"scale_down_delay" validate:"required"`

	// Timeout in secs after which the function is terminated.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time" validate:"required"`

	// Optional amount of memory set for the instance of the function. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit" validate:"required"`

	// The current status of the function.
	Status *string `json:"status,omitempty"`

	// The detailed status of the function.
	StatusDetails *FunctionStatus `json:"status_details" validate:"required"`
}

// Constants associated with the Function.ManagedDomainMappings property.
// Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values
// are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports
// function private visibility.
const (
	Function_ManagedDomainMappings_Local = "local"
	Function_ManagedDomainMappings_LocalPrivate = "local_private"
	Function_ManagedDomainMappings_LocalPublic = "local_public"
)

// Constants associated with the Function.ResourceType property.
// The type of the function.
const (
	Function_ResourceType_FunctionV2 = "function_v2"
)

// Constants associated with the Function.Status property.
// The current status of the function.
const (
	Function_Status_Deploying = "deploying"
	Function_Status_Failed = "failed"
	Function_Status_Offline = "offline"
	Function_Status_Ready = "ready"
)

// UnmarshalFunction unmarshals an instance of Function from the specified map of raw messages.
func UnmarshalFunction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Function)
	err = core.UnmarshalPrimitive(m, "code_binary", &obj.CodeBinary)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_binary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "code_main", &obj.CodeMain)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_main-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "code_reference", &obj.CodeReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "code_secret", &obj.CodeSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "computed_env_variables", &obj.ComputedEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "computed_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoint", &obj.Endpoint)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoint_internal", &obj.EndpointInternal)
	if err != nil {
		err = core.SDKErrorf(err, "", "endpoint_internal-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "managed_domain_mappings", &obj.ManagedDomainMappings)
	if err != nil {
		err = core.SDKErrorf(err, "", "managed_domain_mappings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "runtime", &obj.Runtime)
	if err != nil {
		err = core.SDKErrorf(err, "", "runtime-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency", &obj.ScaleConcurrency)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_down_delay", &obj.ScaleDownDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_down_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_execution_time", &obj.ScaleMaxExecutionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_execution_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalFunctionStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FunctionList : Contains a list of functions and pagination information.
type FunctionList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// List of all functions.
	Functions []Function `json:"functions" validate:"required"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalFunctionList unmarshals an instance of FunctionList from the specified map of raw messages.
func UnmarshalFunctionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FunctionList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "functions", &obj.Functions, UnmarshalFunction)
	if err != nil {
		err = core.SDKErrorf(err, "", "functions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *FunctionList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// FunctionPatch : Request model for function update operations.
type FunctionPatch struct {
	// Specifies whether the code is binary or not. Defaults to false when `code_reference` is set to a data URL. When
	// `code_reference` is set to a code bundle URL, this field is always true.
	CodeBinary *bool `json:"code_binary,omitempty"`

	// Specifies the name of the function that should be invoked.
	CodeMain *string `json:"code_main,omitempty"`

	// Specifies either a reference to a code bundle or the source code itself. To specify the source code, use the data
	// URL scheme and include the source code as base64 encoded. The data URL scheme is defined in [RFC
	// 2397](https://tools.ietf.org/html/rfc2397).
	CodeReference *string `json:"code_reference,omitempty"`

	// The name of the secret that is used to access the specified `code_reference`. The secret is used to authenticate
	// with a non-public endpoint that is specified as`code_reference`.
	CodeSecret *string `json:"code_secret,omitempty"`

	// Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values
	// are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports
	// function private visibility.
	ManagedDomainMappings *string `json:"managed_domain_mappings,omitempty"`

	// Optional references to config maps, secrets or literal values.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// The managed runtime used to execute the injected code.
	Runtime *string `json:"runtime,omitempty"`

	// Number of parallel requests handled by a single instance, supported only by Node.js, default is `1`.
	ScaleConcurrency *int64 `json:"scale_concurrency,omitempty"`

	// Optional amount of CPU set for the instance of the function. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of time in seconds that delays the scale down behavior for a function.
	ScaleDownDelay *int64 `json:"scale_down_delay,omitempty"`

	// Timeout in secs after which the function is terminated.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the function. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`
}

// Constants associated with the FunctionPatch.ManagedDomainMappings property.
// Optional value controlling which of the system managed domain mappings will be setup for the function. Valid values
// are 'local_public', 'local_private' and 'local'. Visibility can only be 'local_private' if the project supports
// function private visibility.
const (
	FunctionPatch_ManagedDomainMappings_Local = "local"
	FunctionPatch_ManagedDomainMappings_LocalPrivate = "local_private"
	FunctionPatch_ManagedDomainMappings_LocalPublic = "local_public"
)

// UnmarshalFunctionPatch unmarshals an instance of FunctionPatch from the specified map of raw messages.
func UnmarshalFunctionPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FunctionPatch)
	err = core.UnmarshalPrimitive(m, "code_binary", &obj.CodeBinary)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_binary-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "code_main", &obj.CodeMain)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_main-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "code_reference", &obj.CodeReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "code_secret", &obj.CodeSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "code_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "managed_domain_mappings", &obj.ManagedDomainMappings)
	if err != nil {
		err = core.SDKErrorf(err, "", "managed_domain_mappings-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVarPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "runtime", &obj.Runtime)
	if err != nil {
		err = core.SDKErrorf(err, "", "runtime-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_concurrency", &obj.ScaleConcurrency)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_concurrency-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_down_delay", &obj.ScaleDownDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_down_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_execution_time", &obj.ScaleMaxExecutionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_execution_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the FunctionPatch
func (functionPatch *FunctionPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(functionPatch.CodeBinary) {
		_patch["code_binary"] = functionPatch.CodeBinary
	}
	if !core.IsNil(functionPatch.CodeMain) {
		_patch["code_main"] = functionPatch.CodeMain
	}
	if !core.IsNil(functionPatch.CodeReference) {
		_patch["code_reference"] = functionPatch.CodeReference
	}
	if !core.IsNil(functionPatch.CodeSecret) {
		_patch["code_secret"] = functionPatch.CodeSecret
	}
	if !core.IsNil(functionPatch.ManagedDomainMappings) {
		_patch["managed_domain_mappings"] = functionPatch.ManagedDomainMappings
	}
	if !core.IsNil(functionPatch.RunEnvVariables) {
		var runEnvVariablesPatches []map[string]interface{}
		for _, runEnvVariables := range functionPatch.RunEnvVariables {
			runEnvVariablesPatches = append(runEnvVariablesPatches, runEnvVariables.asPatch())
		}
		_patch["run_env_variables"] = runEnvVariablesPatches
	}
	if !core.IsNil(functionPatch.Runtime) {
		_patch["runtime"] = functionPatch.Runtime
	}
	if !core.IsNil(functionPatch.ScaleConcurrency) {
		_patch["scale_concurrency"] = functionPatch.ScaleConcurrency
	}
	if !core.IsNil(functionPatch.ScaleCpuLimit) {
		_patch["scale_cpu_limit"] = functionPatch.ScaleCpuLimit
	}
	if !core.IsNil(functionPatch.ScaleDownDelay) {
		_patch["scale_down_delay"] = functionPatch.ScaleDownDelay
	}
	if !core.IsNil(functionPatch.ScaleMaxExecutionTime) {
		_patch["scale_max_execution_time"] = functionPatch.ScaleMaxExecutionTime
	}
	if !core.IsNil(functionPatch.ScaleMemoryLimit) {
		_patch["scale_memory_limit"] = functionPatch.ScaleMemoryLimit
	}

	return
}

// FunctionRuntime : Response model for Function runtime objects.
type FunctionRuntime struct {
	// Whether the function runtime is the default for the code bundle family.
	Default *bool `json:"default,omitempty"`

	// Whether the function runtime is deprecated.
	Deprecated *bool `json:"deprecated,omitempty"`

	// The code bundle family of the function runtime.
	Family *string `json:"family,omitempty"`

	// The ID of the function runtime.
	ID *string `json:"id,omitempty"`

	// The name of the function runtime.
	Name *string `json:"name,omitempty"`

	// Whether the function runtime is optimized.
	Optimized *bool `json:"optimized,omitempty"`
}

// UnmarshalFunctionRuntime unmarshals an instance of FunctionRuntime from the specified map of raw messages.
func UnmarshalFunctionRuntime(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FunctionRuntime)
	err = core.UnmarshalPrimitive(m, "default", &obj.Default)
	if err != nil {
		err = core.SDKErrorf(err, "", "default-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "deprecated", &obj.Deprecated)
	if err != nil {
		err = core.SDKErrorf(err, "", "deprecated-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "family", &obj.Family)
	if err != nil {
		err = core.SDKErrorf(err, "", "family-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "optimized", &obj.Optimized)
	if err != nil {
		err = core.SDKErrorf(err, "", "optimized-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FunctionRuntimeList : Contains a list of Function runtimes.
type FunctionRuntimeList struct {
	// List of all Function runtimes.
	FunctionRuntimes []FunctionRuntime `json:"function_runtimes,omitempty"`
}

// UnmarshalFunctionRuntimeList unmarshals an instance of FunctionRuntimeList from the specified map of raw messages.
func UnmarshalFunctionRuntimeList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FunctionRuntimeList)
	err = core.UnmarshalModel(m, "function_runtimes", &obj.FunctionRuntimes, UnmarshalFunctionRuntime)
	if err != nil {
		err = core.SDKErrorf(err, "", "function_runtimes-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FunctionStatus : The detailed status of the function.
type FunctionStatus struct {
	// Provides additional information about the status of the function.
	Reason *string `json:"reason,omitempty"`
}

// Constants associated with the FunctionStatus.Reason property.
// Provides additional information about the status of the function.
const (
	FunctionStatus_Reason_Deploying = "deploying"
	FunctionStatus_Reason_DeployingConfiguringRoutes = "deploying_configuring_routes"
	FunctionStatus_Reason_NoCodeBundle = "no_code_bundle"
	FunctionStatus_Reason_Offline = "offline"
	FunctionStatus_Reason_Ready = "ready"
	FunctionStatus_Reason_ReadyLastUpdateFailed = "ready_last_update_failed"
	FunctionStatus_Reason_ReadyUpdateInProgress = "ready_update_in_progress"
	FunctionStatus_Reason_UnknownReason = "unknown_reason"
)

// UnmarshalFunctionStatus unmarshals an instance of FunctionStatus from the specified map of raw messages.
func UnmarshalFunctionStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FunctionStatus)
	err = core.UnmarshalPrimitive(m, "reason", &obj.Reason)
	if err != nil {
		err = core.SDKErrorf(err, "", "reason-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetAllowedOutboundDestinationOptions : The GetAllowedOutboundDestination options.
type GetAllowedOutboundDestinationOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your allowed outbound destination.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAllowedOutboundDestinationOptions : Instantiate GetAllowedOutboundDestinationOptions
func (*CodeEngineV2) NewGetAllowedOutboundDestinationOptions(projectID string, name string) *GetAllowedOutboundDestinationOptions {
	return &GetAllowedOutboundDestinationOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetAllowedOutboundDestinationOptions) SetProjectID(projectID string) *GetAllowedOutboundDestinationOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetAllowedOutboundDestinationOptions) SetName(name string) *GetAllowedOutboundDestinationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAllowedOutboundDestinationOptions) SetHeaders(param map[string]string) *GetAllowedOutboundDestinationOptions {
	options.Headers = param
	return options
}

// GetAppOptions : The GetApp options.
type GetAppOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAppOptions : Instantiate GetAppOptions
func (*CodeEngineV2) NewGetAppOptions(projectID string, name string) *GetAppOptions {
	return &GetAppOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetAppOptions) SetProjectID(projectID string) *GetAppOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetAppOptions) SetName(name string) *GetAppOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAppOptions) SetHeaders(param map[string]string) *GetAppOptions {
	options.Headers = param
	return options
}

// GetAppRevisionOptions : The GetAppRevision options.
type GetAppRevisionOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	AppName *string `json:"app_name" validate:"required,ne="`

	// The name of your application revision.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetAppRevisionOptions : Instantiate GetAppRevisionOptions
func (*CodeEngineV2) NewGetAppRevisionOptions(projectID string, appName string, name string) *GetAppRevisionOptions {
	return &GetAppRevisionOptions{
		ProjectID: core.StringPtr(projectID),
		AppName: core.StringPtr(appName),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetAppRevisionOptions) SetProjectID(projectID string) *GetAppRevisionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetAppName : Allow user to set AppName
func (_options *GetAppRevisionOptions) SetAppName(appName string) *GetAppRevisionOptions {
	_options.AppName = core.StringPtr(appName)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetAppRevisionOptions) SetName(name string) *GetAppRevisionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetAppRevisionOptions) SetHeaders(param map[string]string) *GetAppRevisionOptions {
	options.Headers = param
	return options
}

// GetBindingOptions : The GetBinding options.
type GetBindingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The id of your binding.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetBindingOptions : Instantiate GetBindingOptions
func (*CodeEngineV2) NewGetBindingOptions(projectID string, id string) *GetBindingOptions {
	return &GetBindingOptions{
		ProjectID: core.StringPtr(projectID),
		ID: core.StringPtr(id),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetBindingOptions) SetProjectID(projectID string) *GetBindingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetID : Allow user to set ID
func (_options *GetBindingOptions) SetID(id string) *GetBindingOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBindingOptions) SetHeaders(param map[string]string) *GetBindingOptions {
	options.Headers = param
	return options
}

// GetBuildOptions : The GetBuild options.
type GetBuildOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your build.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetBuildOptions : Instantiate GetBuildOptions
func (*CodeEngineV2) NewGetBuildOptions(projectID string, name string) *GetBuildOptions {
	return &GetBuildOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetBuildOptions) SetProjectID(projectID string) *GetBuildOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetBuildOptions) SetName(name string) *GetBuildOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBuildOptions) SetHeaders(param map[string]string) *GetBuildOptions {
	options.Headers = param
	return options
}

// GetBuildRunOptions : The GetBuildRun options.
type GetBuildRunOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your build run.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetBuildRunOptions : Instantiate GetBuildRunOptions
func (*CodeEngineV2) NewGetBuildRunOptions(projectID string, name string) *GetBuildRunOptions {
	return &GetBuildRunOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetBuildRunOptions) SetProjectID(projectID string) *GetBuildRunOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetBuildRunOptions) SetName(name string) *GetBuildRunOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetBuildRunOptions) SetHeaders(param map[string]string) *GetBuildRunOptions {
	options.Headers = param
	return options
}

// GetConfigMapOptions : The GetConfigMap options.
type GetConfigMapOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your configmap.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetConfigMapOptions : Instantiate GetConfigMapOptions
func (*CodeEngineV2) NewGetConfigMapOptions(projectID string, name string) *GetConfigMapOptions {
	return &GetConfigMapOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetConfigMapOptions) SetProjectID(projectID string) *GetConfigMapOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetConfigMapOptions) SetName(name string) *GetConfigMapOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetConfigMapOptions) SetHeaders(param map[string]string) *GetConfigMapOptions {
	options.Headers = param
	return options
}

// GetDomainMappingOptions : The GetDomainMapping options.
type GetDomainMappingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your domain mapping.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetDomainMappingOptions : Instantiate GetDomainMappingOptions
func (*CodeEngineV2) NewGetDomainMappingOptions(projectID string, name string) *GetDomainMappingOptions {
	return &GetDomainMappingOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetDomainMappingOptions) SetProjectID(projectID string) *GetDomainMappingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetDomainMappingOptions) SetName(name string) *GetDomainMappingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetDomainMappingOptions) SetHeaders(param map[string]string) *GetDomainMappingOptions {
	options.Headers = param
	return options
}

// GetFunctionOptions : The GetFunction options.
type GetFunctionOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your function.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetFunctionOptions : Instantiate GetFunctionOptions
func (*CodeEngineV2) NewGetFunctionOptions(projectID string, name string) *GetFunctionOptions {
	return &GetFunctionOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetFunctionOptions) SetProjectID(projectID string) *GetFunctionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetFunctionOptions) SetName(name string) *GetFunctionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetFunctionOptions) SetHeaders(param map[string]string) *GetFunctionOptions {
	options.Headers = param
	return options
}

// GetJobOptions : The GetJob options.
type GetJobOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your job.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetJobOptions : Instantiate GetJobOptions
func (*CodeEngineV2) NewGetJobOptions(projectID string, name string) *GetJobOptions {
	return &GetJobOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetJobOptions) SetProjectID(projectID string) *GetJobOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetJobOptions) SetName(name string) *GetJobOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetJobOptions) SetHeaders(param map[string]string) *GetJobOptions {
	options.Headers = param
	return options
}

// GetJobRunOptions : The GetJobRun options.
type GetJobRunOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your job run.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetJobRunOptions : Instantiate GetJobRunOptions
func (*CodeEngineV2) NewGetJobRunOptions(projectID string, name string) *GetJobRunOptions {
	return &GetJobRunOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetJobRunOptions) SetProjectID(projectID string) *GetJobRunOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetJobRunOptions) SetName(name string) *GetJobRunOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetJobRunOptions) SetHeaders(param map[string]string) *GetJobRunOptions {
	options.Headers = param
	return options
}

// GetProjectEgressIpsOptions : The GetProjectEgressIps options.
type GetProjectEgressIpsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProjectEgressIpsOptions : Instantiate GetProjectEgressIpsOptions
func (*CodeEngineV2) NewGetProjectEgressIpsOptions(projectID string) *GetProjectEgressIpsOptions {
	return &GetProjectEgressIpsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetProjectEgressIpsOptions) SetProjectID(projectID string) *GetProjectEgressIpsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProjectEgressIpsOptions) SetHeaders(param map[string]string) *GetProjectEgressIpsOptions {
	options.Headers = param
	return options
}

// GetProjectOptions : The GetProject options.
type GetProjectOptions struct {
	// The ID of the project.
	ID *string `json:"id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProjectOptions : Instantiate GetProjectOptions
func (*CodeEngineV2) NewGetProjectOptions(id string) *GetProjectOptions {
	return &GetProjectOptions{
		ID: core.StringPtr(id),
	}
}

// SetID : Allow user to set ID
func (_options *GetProjectOptions) SetID(id string) *GetProjectOptions {
	_options.ID = core.StringPtr(id)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProjectOptions) SetHeaders(param map[string]string) *GetProjectOptions {
	options.Headers = param
	return options
}

// GetProjectStatusDetailsOptions : The GetProjectStatusDetails options.
type GetProjectStatusDetailsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetProjectStatusDetailsOptions : Instantiate GetProjectStatusDetailsOptions
func (*CodeEngineV2) NewGetProjectStatusDetailsOptions(projectID string) *GetProjectStatusDetailsOptions {
	return &GetProjectStatusDetailsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetProjectStatusDetailsOptions) SetProjectID(projectID string) *GetProjectStatusDetailsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetProjectStatusDetailsOptions) SetHeaders(param map[string]string) *GetProjectStatusDetailsOptions {
	options.Headers = param
	return options
}

// GetSecretOptions : The GetSecret options.
type GetSecretOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your secret.
	Name *string `json:"name" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetSecretOptions : Instantiate GetSecretOptions
func (*CodeEngineV2) NewGetSecretOptions(projectID string, name string) *GetSecretOptions {
	return &GetSecretOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *GetSecretOptions) SetProjectID(projectID string) *GetSecretOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *GetSecretOptions) SetName(name string) *GetSecretOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSecretOptions) SetHeaders(param map[string]string) *GetSecretOptions {
	options.Headers = param
	return options
}

// Job : Job is the response model for job resources.
type Job struct {
	// Reference to a build that is associated with the job.
	Build *string `json:"build,omitempty"`

	// Reference to a build run that is associated with the job.
	BuildRun *string `json:"build_run,omitempty"`

	// References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as
	// environment variables in the job run.
	ComputedEnvVariables []EnvVar `json:"computed_env_variables,omitempty"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The version of the job instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// When you provision a new job,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference" validate:"required"`

	// The name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the job / job runs will be created but submitted job runs will fail, until this property is
	// provided, too. This property must not be set on a job run, which references a job template.
	ImageSecret *string `json:"image_secret,omitempty"`

	// The name of the job.
	Name *string `json:"name" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the job.
	ResourceType *string `json:"resource_type,omitempty"`

	// Set arguments for the job that are passed to start job run containers. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments" validate:"required"`

	// The user ID (UID) to run the job.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Set commands for the job that are passed to start job run containers. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands" validate:"required"`

	// References to config maps, secrets or literal values, which are defined by the function owner and are exposed as
	// environment variables in the job run.
	RunEnvVariables []EnvVar `json:"run_env_variables" validate:"required"`

	// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
	// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
	// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
	RunMode *string `json:"run_mode" validate:"required"`

	// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
	// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or secrets.
	RunVolumeMounts []VolumeMount `json:"run_volume_mounts" validate:"required"`

	// Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges,
	// such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number
	// of unique array indices that you specify with this parameter determines the number of job instances to run.
	ScaleArraySpec *string `json:"scale_array_spec" validate:"required"`

	// Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit" validate:"required"`

	// Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit" validate:"required"`

	// The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is
	// `task`.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit" validate:"required"`

	// The number of times to rerun an instance of the job before the job is marked as failed. This property can only be
	// specified if `run_mode` is `task`.
	ScaleRetryLimit *int64 `json:"scale_retry_limit,omitempty"`
}

// Constants associated with the Job.ResourceType property.
// The type of the job.
const (
	Job_ResourceType_JobV2 = "job_v2"
)

// Constants associated with the Job.RunMode property.
// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
const (
	Job_RunMode_Daemon = "daemon"
	Job_RunMode_Task = "task"
)

// Constants associated with the Job.RunServiceAccount property.
// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
const (
	Job_RunServiceAccount_Default = "default"
	Job_RunServiceAccount_Manager = "manager"
	Job_RunServiceAccount_None = "none"
	Job_RunServiceAccount_Reader = "reader"
	Job_RunServiceAccount_Writer = "writer"
)

// UnmarshalJob unmarshals an instance of Job from the specified map of raw messages.
func UnmarshalJob(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Job)
	err = core.UnmarshalPrimitive(m, "build", &obj.Build)
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "build_run", &obj.BuildRun)
	if err != nil {
		err = core.SDKErrorf(err, "", "build_run-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "computed_env_variables", &obj.ComputedEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "computed_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "image_reference", &obj.ImageReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_secret", &obj.ImageSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_arguments", &obj.RunArguments)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_arguments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_as_user", &obj.RunAsUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_as_user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_commands", &obj.RunCommands)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_commands-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_mode", &obj.RunMode)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_mode-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_service_account", &obj.RunServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_volume_mounts", &obj.RunVolumeMounts, UnmarshalVolumeMount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_volume_mounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_array_spec", &obj.ScaleArraySpec)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_array_spec-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_execution_time", &obj.ScaleMaxExecutionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_execution_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_retry_limit", &obj.ScaleRetryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_retry_limit-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobList : Contains a list of jobs and pagination information.
type JobList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// List of all jobs.
	Jobs []Job `json:"jobs" validate:"required"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalJobList unmarshals an instance of JobList from the specified map of raw messages.
func UnmarshalJobList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "jobs", &obj.Jobs, UnmarshalJob)
	if err != nil {
		err = core.SDKErrorf(err, "", "jobs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *JobList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// JobPatch : Request model for job update operations.
type JobPatch struct {
	// The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference,omitempty"`

	// The name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the job / job runs will be created but submitted job runs will fail, until this property is
	// provided, too. This property must not be set on a job run, which references a job template.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Set arguments for the job that are passed to start job run containers. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments,omitempty"`

	// The user ID (UID) to run the job.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Set commands for the job that are passed to start job run containers. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands,omitempty"`

	// Optional references to config maps, secrets or literal values.
	RunEnvVariables []EnvVarPrototype `json:"run_env_variables,omitempty"`

	// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
	// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
	// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
	RunMode *string `json:"run_mode,omitempty"`

	// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
	// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or a secrets. In case this is provided, existing `run_volume_mounts` will be
	// overwritten.
	RunVolumeMounts []VolumeMountPrototype `json:"run_volume_mounts,omitempty"`

	// Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges,
	// such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number
	// of unique array indices that you specify with this parameter determines the number of job instances to run.
	ScaleArraySpec *string `json:"scale_array_spec,omitempty"`

	// Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit,omitempty"`

	// The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is
	// `task`.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// The number of times to rerun an instance of the job before the job is marked as failed. This property can only be
	// specified if `run_mode` is `task`.
	ScaleRetryLimit *int64 `json:"scale_retry_limit,omitempty"`
}

// Constants associated with the JobPatch.RunMode property.
// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
const (
	JobPatch_RunMode_Daemon = "daemon"
	JobPatch_RunMode_Task = "task"
)

// Constants associated with the JobPatch.RunServiceAccount property.
// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
const (
	JobPatch_RunServiceAccount_Default = "default"
	JobPatch_RunServiceAccount_Manager = "manager"
	JobPatch_RunServiceAccount_None = "none"
	JobPatch_RunServiceAccount_Reader = "reader"
	JobPatch_RunServiceAccount_Writer = "writer"
)

// UnmarshalJobPatch unmarshals an instance of JobPatch from the specified map of raw messages.
func UnmarshalJobPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobPatch)
	err = core.UnmarshalPrimitive(m, "image_reference", &obj.ImageReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_secret", &obj.ImageSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_arguments", &obj.RunArguments)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_arguments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_as_user", &obj.RunAsUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_as_user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_commands", &obj.RunCommands)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_commands-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVarPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_mode", &obj.RunMode)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_mode-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_service_account", &obj.RunServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_volume_mounts", &obj.RunVolumeMounts, UnmarshalVolumeMountPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_volume_mounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_array_spec", &obj.ScaleArraySpec)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_array_spec-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_execution_time", &obj.ScaleMaxExecutionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_execution_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_retry_limit", &obj.ScaleRetryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_retry_limit-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the JobPatch
func (jobPatch *JobPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(jobPatch.ImageReference) {
		_patch["image_reference"] = jobPatch.ImageReference
	}
	if !core.IsNil(jobPatch.ImageSecret) {
		_patch["image_secret"] = jobPatch.ImageSecret
	}
	if !core.IsNil(jobPatch.RunArguments) {
		_patch["run_arguments"] = jobPatch.RunArguments
	}
	if !core.IsNil(jobPatch.RunAsUser) {
		_patch["run_as_user"] = jobPatch.RunAsUser
	}
	if !core.IsNil(jobPatch.RunCommands) {
		_patch["run_commands"] = jobPatch.RunCommands
	}
	if !core.IsNil(jobPatch.RunEnvVariables) {
		var runEnvVariablesPatches []map[string]interface{}
		for _, runEnvVariables := range jobPatch.RunEnvVariables {
			runEnvVariablesPatches = append(runEnvVariablesPatches, runEnvVariables.asPatch())
		}
		_patch["run_env_variables"] = runEnvVariablesPatches
	}
	if !core.IsNil(jobPatch.RunMode) {
		_patch["run_mode"] = jobPatch.RunMode
	}
	if !core.IsNil(jobPatch.RunServiceAccount) {
		_patch["run_service_account"] = jobPatch.RunServiceAccount
	}
	if !core.IsNil(jobPatch.RunVolumeMounts) {
		var runVolumeMountsPatches []map[string]interface{}
		for _, runVolumeMounts := range jobPatch.RunVolumeMounts {
			runVolumeMountsPatches = append(runVolumeMountsPatches, runVolumeMounts.asPatch())
		}
		_patch["run_volume_mounts"] = runVolumeMountsPatches
	}
	if !core.IsNil(jobPatch.ScaleArraySpec) {
		_patch["scale_array_spec"] = jobPatch.ScaleArraySpec
	}
	if !core.IsNil(jobPatch.ScaleCpuLimit) {
		_patch["scale_cpu_limit"] = jobPatch.ScaleCpuLimit
	}
	if !core.IsNil(jobPatch.ScaleEphemeralStorageLimit) {
		_patch["scale_ephemeral_storage_limit"] = jobPatch.ScaleEphemeralStorageLimit
	}
	if !core.IsNil(jobPatch.ScaleMaxExecutionTime) {
		_patch["scale_max_execution_time"] = jobPatch.ScaleMaxExecutionTime
	}
	if !core.IsNil(jobPatch.ScaleMemoryLimit) {
		_patch["scale_memory_limit"] = jobPatch.ScaleMemoryLimit
	}
	if !core.IsNil(jobPatch.ScaleRetryLimit) {
		_patch["scale_retry_limit"] = jobPatch.ScaleRetryLimit
	}

	return
}

// JobRun : Response model for job run resources.
type JobRun struct {
	// References to config maps, secrets or literal values, which are defined and set by Code Engine and are exposed as
	// environment variables in the job run.
	ComputedEnvVariables []EnvVar `json:"computed_env_variables,omitempty"`

	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// When you provision a new job run,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the image that is used for this job. The format is `REGISTRY/NAMESPACE/REPOSITORY:TAG` where `REGISTRY`
	// and `TAG` are optional. If `REGISTRY` is not specified, the default is `docker.io`. If `TAG` is not specified, the
	// default is `latest`. If the image reference points to a registry that requires authentication, make sure to also
	// specify the property `image_secret`.
	ImageReference *string `json:"image_reference,omitempty"`

	// The name of the image registry access secret. The image registry access secret is used to authenticate with a
	// private registry when you download the container image. If the image reference points to a registry that requires
	// authentication, the job / job runs will be created but submitted job runs will fail, until this property is
	// provided, too. This property must not be set on a job run, which references a job template.
	ImageSecret *string `json:"image_secret,omitempty"`

	// Optional name of the job reference of this job run. If specified, the job run will inherit the configuration of the
	// referenced job.
	JobName *string `json:"job_name,omitempty"`

	// The name of the job run.
	Name *string `json:"name,omitempty"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the job run.
	ResourceType *string `json:"resource_type,omitempty"`

	// Set arguments for the job that are passed to start job run containers. If not specified an empty string array will
	// be applied and the arguments specified by the container image, will be used to start the container.
	RunArguments []string `json:"run_arguments" validate:"required"`

	// The user ID (UID) to run the job.
	RunAsUser *int64 `json:"run_as_user,omitempty"`

	// Set commands for the job that are passed to start job run containers. If not specified an empty string array will be
	// applied and the command specified by the container image, will be used to start the container.
	RunCommands []string `json:"run_commands" validate:"required"`

	// References to config maps, secrets or literal values, which are defined by the function owner and are exposed as
	// environment variables in the job run.
	RunEnvVariables []EnvVar `json:"run_env_variables" validate:"required"`

	// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
	// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
	// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
	RunMode *string `json:"run_mode,omitempty"`

	// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
	// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
	RunServiceAccount *string `json:"run_service_account,omitempty"`

	// Optional mounts of config maps or secrets.
	RunVolumeMounts []VolumeMount `json:"run_volume_mounts" validate:"required"`

	// Optional value to override the JOB_ARRAY_SIZE environment variable for a job run.
	ScaleArraySizeVariableOverride *int64 `json:"scale_array_size_variable_override,omitempty"`

	// Define a custom set of array indices as a comma-separated list containing single values and hyphen-separated ranges,
	// such as  5,12-14,23,27. Each instance gets its array index value from the environment variable JOB_INDEX. The number
	// of unique array indices that you specify with this parameter determines the number of job instances to run.
	ScaleArraySpec *string `json:"scale_array_spec,omitempty"`

	// Optional amount of CPU set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo).
	ScaleCpuLimit *string `json:"scale_cpu_limit,omitempty"`

	// Optional amount of ephemeral storage to set for the instance of the job. The amount specified as ephemeral storage,
	// must not exceed the amount of `scale_memory_limit`. The units for specifying ephemeral storage are Megabyte (M) or
	// Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information see [Units of
	// measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleEphemeralStorageLimit *string `json:"scale_ephemeral_storage_limit,omitempty"`

	// The maximum execution time in seconds for runs of the job. This property can only be specified if `run_mode` is
	// `task`.
	ScaleMaxExecutionTime *int64 `json:"scale_max_execution_time,omitempty"`

	// Optional amount of memory set for the instance of the job. For valid values see [Supported memory and CPU
	// combinations](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo). The units for specifying memory
	// are Megabyte (M) or Gigabyte (G), whereas G and M are the shorthand expressions for GB and MB. For more information
	// see [Units of measurement](https://cloud.ibm.com/docs/codeengine?topic=codeengine-mem-cpu-combo#unit-measurements).
	ScaleMemoryLimit *string `json:"scale_memory_limit,omitempty"`

	// The number of times to rerun an instance of the job before the job is marked as failed. This property can only be
	// specified if `run_mode` is `task`.
	ScaleRetryLimit *int64 `json:"scale_retry_limit,omitempty"`

	// The current status of the job run.
	Status *string `json:"status,omitempty"`

	// The detailed status of the job run.
	StatusDetails *JobRunStatus `json:"status_details,omitempty"`
}

// Constants associated with the JobRun.ResourceType property.
// The type of the job run.
const (
	JobRun_ResourceType_JobRunV2 = "job_run_v2"
)

// Constants associated with the JobRun.RunMode property.
// The mode for runs of the job. Valid values are `task` and `daemon`. In `task` mode, the `max_execution_time` and
// `retry_limit` properties apply. In `daemon` mode, since there is no timeout and failed instances are restarted
// indefinitely, the `max_execution_time` and `retry_limit` properties are not allowed.
const (
	JobRun_RunMode_Daemon = "daemon"
	JobRun_RunMode_Task = "task"
)

// Constants associated with the JobRun.RunServiceAccount property.
// The name of the service account. For built-in service accounts, you can use the shortened names `manager`, `none`,
// `reader`, and `writer`. This property must not be set on a job run, which references a job template.
const (
	JobRun_RunServiceAccount_Default = "default"
	JobRun_RunServiceAccount_Manager = "manager"
	JobRun_RunServiceAccount_None = "none"
	JobRun_RunServiceAccount_Reader = "reader"
	JobRun_RunServiceAccount_Writer = "writer"
)

// Constants associated with the JobRun.Status property.
// The current status of the job run.
const (
	JobRun_Status_Completed = "completed"
	JobRun_Status_Failed = "failed"
	JobRun_Status_Pending = "pending"
	JobRun_Status_Running = "running"
)

// UnmarshalJobRun unmarshals an instance of JobRun from the specified map of raw messages.
func UnmarshalJobRun(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobRun)
	err = core.UnmarshalModel(m, "computed_env_variables", &obj.ComputedEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "computed_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "image_reference", &obj.ImageReference)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "image_secret", &obj.ImageSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "image_secret-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "job_name", &obj.JobName)
	if err != nil {
		err = core.SDKErrorf(err, "", "job_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_arguments", &obj.RunArguments)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_arguments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_as_user", &obj.RunAsUser)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_as_user-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_commands", &obj.RunCommands)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_commands-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_env_variables", &obj.RunEnvVariables, UnmarshalEnvVar)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_env_variables-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_mode", &obj.RunMode)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_mode-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "run_service_account", &obj.RunServiceAccount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_service_account-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "run_volume_mounts", &obj.RunVolumeMounts, UnmarshalVolumeMount)
	if err != nil {
		err = core.SDKErrorf(err, "", "run_volume_mounts-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_array_size_variable_override", &obj.ScaleArraySizeVariableOverride)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_array_size_variable_override-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_array_spec", &obj.ScaleArraySpec)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_array_spec-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_cpu_limit", &obj.ScaleCpuLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_cpu_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_ephemeral_storage_limit", &obj.ScaleEphemeralStorageLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_ephemeral_storage_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_max_execution_time", &obj.ScaleMaxExecutionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_max_execution_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_memory_limit", &obj.ScaleMemoryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_memory_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "scale_retry_limit", &obj.ScaleRetryLimit)
	if err != nil {
		err = core.SDKErrorf(err, "", "scale_retry_limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "status_details", &obj.StatusDetails, UnmarshalJobRunStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_details-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// JobRunList : Contains a list of job runs and pagination information.
type JobRunList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// List of all jobs.
	JobRuns []JobRun `json:"job_runs" validate:"required"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`
}

// UnmarshalJobRunList unmarshals an instance of JobRunList from the specified map of raw messages.
func UnmarshalJobRunList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobRunList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "job_runs", &obj.JobRuns, UnmarshalJobRun)
	if err != nil {
		err = core.SDKErrorf(err, "", "job_runs-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *JobRunList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// JobRunStatus : The detailed status of the job run.
type JobRunStatus struct {
	// Time the job run completed.
	CompletionTime *string `json:"completion_time,omitempty"`

	// Number of failed job run instances.
	Failed *int64 `json:"failed,omitempty"`

	// List of job run indices that failed.
	FailedIndices *string `json:"failed_indices,omitempty"`

	// Number of pending job run instances.
	Pending *int64 `json:"pending,omitempty"`

	// List of job run indices that are pending.
	PendingIndices *string `json:"pending_indices,omitempty"`

	// Number of requested job run instances.
	Requested *int64 `json:"requested,omitempty"`

	// Number of running job run instances.
	Running *int64 `json:"running,omitempty"`

	// List of job run indices that are running.
	RunningIndices *string `json:"running_indices,omitempty"`

	// Time the job run started.
	StartTime *string `json:"start_time,omitempty"`

	// Number of succeeded job run instances.
	Succeeded *int64 `json:"succeeded,omitempty"`

	// List of job run indices that succeeded.
	SucceededIndices *string `json:"succeeded_indices,omitempty"`

	// Number of job run instances with unknown state.
	Unknown *int64 `json:"unknown,omitempty"`
}

// UnmarshalJobRunStatus unmarshals an instance of JobRunStatus from the specified map of raw messages.
func UnmarshalJobRunStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(JobRunStatus)
	err = core.UnmarshalPrimitive(m, "completion_time", &obj.CompletionTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "completion_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed", &obj.Failed)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "failed_indices", &obj.FailedIndices)
	if err != nil {
		err = core.SDKErrorf(err, "", "failed_indices-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pending", &obj.Pending)
	if err != nil {
		err = core.SDKErrorf(err, "", "pending-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "pending_indices", &obj.PendingIndices)
	if err != nil {
		err = core.SDKErrorf(err, "", "pending_indices-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "requested", &obj.Requested)
	if err != nil {
		err = core.SDKErrorf(err, "", "requested-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "running", &obj.Running)
	if err != nil {
		err = core.SDKErrorf(err, "", "running-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "running_indices", &obj.RunningIndices)
	if err != nil {
		err = core.SDKErrorf(err, "", "running_indices-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start_time", &obj.StartTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "start_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "succeeded", &obj.Succeeded)
	if err != nil {
		err = core.SDKErrorf(err, "", "succeeded-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "succeeded_indices", &obj.SucceededIndices)
	if err != nil {
		err = core.SDKErrorf(err, "", "succeeded_indices-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "unknown", &obj.Unknown)
	if err != nil {
		err = core.SDKErrorf(err, "", "unknown-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListAllowedOutboundDestinationOptions : The ListAllowedOutboundDestination options.
type ListAllowedOutboundDestinationOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of allowed outbound destinations per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListAllowedOutboundDestinationOptions : Instantiate ListAllowedOutboundDestinationOptions
func (*CodeEngineV2) NewListAllowedOutboundDestinationOptions(projectID string) *ListAllowedOutboundDestinationOptions {
	return &ListAllowedOutboundDestinationOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListAllowedOutboundDestinationOptions) SetProjectID(projectID string) *ListAllowedOutboundDestinationOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAllowedOutboundDestinationOptions) SetLimit(limit int64) *ListAllowedOutboundDestinationOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListAllowedOutboundDestinationOptions) SetStart(start string) *ListAllowedOutboundDestinationOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAllowedOutboundDestinationOptions) SetHeaders(param map[string]string) *ListAllowedOutboundDestinationOptions {
	options.Headers = param
	return options
}

// ListAppInstancesOptions : The ListAppInstances options.
type ListAppInstancesOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	AppName *string `json:"app_name" validate:"required,ne="`

	// Optional maximum number of apps per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListAppInstancesOptions : Instantiate ListAppInstancesOptions
func (*CodeEngineV2) NewListAppInstancesOptions(projectID string, appName string) *ListAppInstancesOptions {
	return &ListAppInstancesOptions{
		ProjectID: core.StringPtr(projectID),
		AppName: core.StringPtr(appName),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListAppInstancesOptions) SetProjectID(projectID string) *ListAppInstancesOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetAppName : Allow user to set AppName
func (_options *ListAppInstancesOptions) SetAppName(appName string) *ListAppInstancesOptions {
	_options.AppName = core.StringPtr(appName)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAppInstancesOptions) SetLimit(limit int64) *ListAppInstancesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListAppInstancesOptions) SetStart(start string) *ListAppInstancesOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAppInstancesOptions) SetHeaders(param map[string]string) *ListAppInstancesOptions {
	options.Headers = param
	return options
}

// ListAppRevisionsOptions : The ListAppRevisions options.
type ListAppRevisionsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	AppName *string `json:"app_name" validate:"required,ne="`

	// Optional maximum number of apps per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListAppRevisionsOptions : Instantiate ListAppRevisionsOptions
func (*CodeEngineV2) NewListAppRevisionsOptions(projectID string, appName string) *ListAppRevisionsOptions {
	return &ListAppRevisionsOptions{
		ProjectID: core.StringPtr(projectID),
		AppName: core.StringPtr(appName),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListAppRevisionsOptions) SetProjectID(projectID string) *ListAppRevisionsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetAppName : Allow user to set AppName
func (_options *ListAppRevisionsOptions) SetAppName(appName string) *ListAppRevisionsOptions {
	_options.AppName = core.StringPtr(appName)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAppRevisionsOptions) SetLimit(limit int64) *ListAppRevisionsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListAppRevisionsOptions) SetStart(start string) *ListAppRevisionsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAppRevisionsOptions) SetHeaders(param map[string]string) *ListAppRevisionsOptions {
	options.Headers = param
	return options
}

// ListAppsOptions : The ListApps options.
type ListAppsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of apps per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListAppsOptions : Instantiate ListAppsOptions
func (*CodeEngineV2) NewListAppsOptions(projectID string) *ListAppsOptions {
	return &ListAppsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListAppsOptions) SetProjectID(projectID string) *ListAppsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListAppsOptions) SetLimit(limit int64) *ListAppsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListAppsOptions) SetStart(start string) *ListAppsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListAppsOptions) SetHeaders(param map[string]string) *ListAppsOptions {
	options.Headers = param
	return options
}

// ListBindingsOptions : The ListBindings options.
type ListBindingsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of bindings per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListBindingsOptions : Instantiate ListBindingsOptions
func (*CodeEngineV2) NewListBindingsOptions(projectID string) *ListBindingsOptions {
	return &ListBindingsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListBindingsOptions) SetProjectID(projectID string) *ListBindingsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListBindingsOptions) SetLimit(limit int64) *ListBindingsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListBindingsOptions) SetStart(start string) *ListBindingsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListBindingsOptions) SetHeaders(param map[string]string) *ListBindingsOptions {
	options.Headers = param
	return options
}

// ListBuildRunsOptions : The ListBuildRuns options.
type ListBuildRunsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional name of the build that should be filtered for.
	BuildName *string `json:"build_name,omitempty"`

	// Optional maximum number of build runs per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListBuildRunsOptions : Instantiate ListBuildRunsOptions
func (*CodeEngineV2) NewListBuildRunsOptions(projectID string) *ListBuildRunsOptions {
	return &ListBuildRunsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListBuildRunsOptions) SetProjectID(projectID string) *ListBuildRunsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetBuildName : Allow user to set BuildName
func (_options *ListBuildRunsOptions) SetBuildName(buildName string) *ListBuildRunsOptions {
	_options.BuildName = core.StringPtr(buildName)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListBuildRunsOptions) SetLimit(limit int64) *ListBuildRunsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListBuildRunsOptions) SetStart(start string) *ListBuildRunsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListBuildRunsOptions) SetHeaders(param map[string]string) *ListBuildRunsOptions {
	options.Headers = param
	return options
}

// ListBuildsOptions : The ListBuilds options.
type ListBuildsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of builds per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListBuildsOptions : Instantiate ListBuildsOptions
func (*CodeEngineV2) NewListBuildsOptions(projectID string) *ListBuildsOptions {
	return &ListBuildsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListBuildsOptions) SetProjectID(projectID string) *ListBuildsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListBuildsOptions) SetLimit(limit int64) *ListBuildsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListBuildsOptions) SetStart(start string) *ListBuildsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListBuildsOptions) SetHeaders(param map[string]string) *ListBuildsOptions {
	options.Headers = param
	return options
}

// ListConfigMapsOptions : The ListConfigMaps options.
type ListConfigMapsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of config maps per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListConfigMapsOptions : Instantiate ListConfigMapsOptions
func (*CodeEngineV2) NewListConfigMapsOptions(projectID string) *ListConfigMapsOptions {
	return &ListConfigMapsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListConfigMapsOptions) SetProjectID(projectID string) *ListConfigMapsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListConfigMapsOptions) SetLimit(limit int64) *ListConfigMapsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListConfigMapsOptions) SetStart(start string) *ListConfigMapsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListConfigMapsOptions) SetHeaders(param map[string]string) *ListConfigMapsOptions {
	options.Headers = param
	return options
}

// ListDomainMappingsOptions : The ListDomainMappings options.
type ListDomainMappingsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of domain mappings per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListDomainMappingsOptions : Instantiate ListDomainMappingsOptions
func (*CodeEngineV2) NewListDomainMappingsOptions(projectID string) *ListDomainMappingsOptions {
	return &ListDomainMappingsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListDomainMappingsOptions) SetProjectID(projectID string) *ListDomainMappingsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListDomainMappingsOptions) SetLimit(limit int64) *ListDomainMappingsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListDomainMappingsOptions) SetStart(start string) *ListDomainMappingsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListDomainMappingsOptions) SetHeaders(param map[string]string) *ListDomainMappingsOptions {
	options.Headers = param
	return options
}

// ListFirstMetadata : Describes properties needed to retrieve the first page of a result list.
type ListFirstMetadata struct {
	// Href that points to the first page.
	Href *string `json:"href,omitempty"`
}

// UnmarshalListFirstMetadata unmarshals an instance of ListFirstMetadata from the specified map of raw messages.
func UnmarshalListFirstMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListFirstMetadata)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListFunctionRuntimesOptions : The ListFunctionRuntimes options.
type ListFunctionRuntimesOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListFunctionRuntimesOptions : Instantiate ListFunctionRuntimesOptions
func (*CodeEngineV2) NewListFunctionRuntimesOptions() *ListFunctionRuntimesOptions {
	return &ListFunctionRuntimesOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListFunctionRuntimesOptions) SetHeaders(param map[string]string) *ListFunctionRuntimesOptions {
	options.Headers = param
	return options
}

// ListFunctionsOptions : The ListFunctions options.
type ListFunctionsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of functions per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the 'next_url' field of the
	// operation response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListFunctionsOptions : Instantiate ListFunctionsOptions
func (*CodeEngineV2) NewListFunctionsOptions(projectID string) *ListFunctionsOptions {
	return &ListFunctionsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListFunctionsOptions) SetProjectID(projectID string) *ListFunctionsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListFunctionsOptions) SetLimit(limit int64) *ListFunctionsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListFunctionsOptions) SetStart(start string) *ListFunctionsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListFunctionsOptions) SetHeaders(param map[string]string) *ListFunctionsOptions {
	options.Headers = param
	return options
}

// ListJobRunsOptions : The ListJobRuns options.
type ListJobRunsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional name of the job that you want to use to filter.
	JobName *string `json:"job_name,omitempty"`

	// Optional maximum number of job runs per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListJobRunsOptions : Instantiate ListJobRunsOptions
func (*CodeEngineV2) NewListJobRunsOptions(projectID string) *ListJobRunsOptions {
	return &ListJobRunsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListJobRunsOptions) SetProjectID(projectID string) *ListJobRunsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetJobName : Allow user to set JobName
func (_options *ListJobRunsOptions) SetJobName(jobName string) *ListJobRunsOptions {
	_options.JobName = core.StringPtr(jobName)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListJobRunsOptions) SetLimit(limit int64) *ListJobRunsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListJobRunsOptions) SetStart(start string) *ListJobRunsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobRunsOptions) SetHeaders(param map[string]string) *ListJobRunsOptions {
	options.Headers = param
	return options
}

// ListJobsOptions : The ListJobs options.
type ListJobsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of jobs per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListJobsOptions : Instantiate ListJobsOptions
func (*CodeEngineV2) NewListJobsOptions(projectID string) *ListJobsOptions {
	return &ListJobsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListJobsOptions) SetProjectID(projectID string) *ListJobsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListJobsOptions) SetLimit(limit int64) *ListJobsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListJobsOptions) SetStart(start string) *ListJobsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListJobsOptions) SetHeaders(param map[string]string) *ListJobsOptions {
	options.Headers = param
	return options
}

// ListNextMetadata : Describes properties needed to retrieve the next page of a result list.
type ListNextMetadata struct {
	// Href that points to the next page.
	Href *string `json:"href,omitempty"`

	// Token.
	Start *string `json:"start,omitempty"`
}

// UnmarshalListNextMetadata unmarshals an instance of ListNextMetadata from the specified map of raw messages.
func UnmarshalListNextMetadata(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListNextMetadata)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListProjectsOptions : The ListProjects options.
type ListProjectsOptions struct {
	// Optional maximum number of projects per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. Any additional query
	// parameters are ignored if a page token is present. If omitted, the first page of results is returned. This value is
	// obtained from the 'start' query parameter in the `next` object of the operation response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListProjectsOptions : Instantiate ListProjectsOptions
func (*CodeEngineV2) NewListProjectsOptions() *ListProjectsOptions {
	return &ListProjectsOptions{}
}

// SetLimit : Allow user to set Limit
func (_options *ListProjectsOptions) SetLimit(limit int64) *ListProjectsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListProjectsOptions) SetStart(start string) *ListProjectsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListProjectsOptions) SetHeaders(param map[string]string) *ListProjectsOptions {
	options.Headers = param
	return options
}

// ListSecretsOptions : The ListSecrets options.
type ListSecretsOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// Optional maximum number of secrets per page.
	Limit *int64 `json:"limit,omitempty"`

	// An optional token that indicates the beginning of the page of results to be returned. If omitted, the first page of
	// results is returned. This value is obtained from the 'start' query parameter in the `next` object of the operation
	// response.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListSecretsOptions : Instantiate ListSecretsOptions
func (*CodeEngineV2) NewListSecretsOptions(projectID string) *ListSecretsOptions {
	return &ListSecretsOptions{
		ProjectID: core.StringPtr(projectID),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ListSecretsOptions) SetProjectID(projectID string) *ListSecretsOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSecretsOptions) SetLimit(limit int64) *ListSecretsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListSecretsOptions) SetStart(start string) *ListSecretsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSecretsOptions) SetHeaders(param map[string]string) *ListSecretsOptions {
	options.Headers = param
	return options
}

// OperatorSecretProps : Properties for the IBM Cloud Operator Secret.
type OperatorSecretProps struct {
	// The ID of the apikey associated with the operator secret.
	ApikeyID *string `json:"apikey_id" validate:"required"`

	// The list of resource groups (by ID) that the operator secret can bind services in.
	ResourceGroupIds []string `json:"resource_group_ids" validate:"required"`

	// A reference to a Service ID.
	Serviceid *ServiceIDRef `json:"serviceid" validate:"required"`

	// Specifies whether the operator secret is user managed.
	UserManaged *bool `json:"user_managed" validate:"required"`
}

// UnmarshalOperatorSecretProps unmarshals an instance of OperatorSecretProps from the specified map of raw messages.
func UnmarshalOperatorSecretProps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OperatorSecretProps)
	err = core.UnmarshalPrimitive(m, "apikey_id", &obj.ApikeyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "apikey_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_ids", &obj.ResourceGroupIds)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_ids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalServiceIDRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceid-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "user_managed", &obj.UserManaged)
	if err != nil {
		err = core.SDKErrorf(err, "", "user_managed-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OperatorSecretPrototypeProps : Properties for the IBM Cloud Operator Secrets.
type OperatorSecretPrototypeProps struct {
	// The list of resource groups (by ID) that the operator secret can bind services in.
	ResourceGroupIds []string `json:"resource_group_ids,omitempty"`

	// A reference to the Service ID.
	Serviceid *ServiceIDRefPrototype `json:"serviceid,omitempty"`
}

// UnmarshalOperatorSecretPrototypeProps unmarshals an instance of OperatorSecretPrototypeProps from the specified map of raw messages.
func UnmarshalOperatorSecretPrototypeProps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OperatorSecretPrototypeProps)
	err = core.UnmarshalPrimitive(m, "resource_group_ids", &obj.ResourceGroupIds)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_ids-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalServiceIDRefPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Probe : Response model for probes.
type Probe struct {
	// The number of consecutive, unsuccessful checks for the probe to be considered failed.
	FailureThreshold *int64 `json:"failure_threshold,omitempty"`

	// The amount of time in seconds to wait before the first probe check is performed.
	InitialDelay *int64 `json:"initial_delay,omitempty"`

	// The amount of time in seconds between probe checks.
	Interval *int64 `json:"interval,omitempty"`

	// The path of the HTTP request to the resource. A path is only supported for a probe with a `type` of `http`.
	Path *string `json:"path,omitempty"`

	// The port on which to probe the resource.
	Port *int64 `json:"port,omitempty"`

	// The amount of time in seconds that the probe waits for a response from the application before it times out and
	// fails.
	Timeout *int64 `json:"timeout,omitempty"`

	// Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the Probe.Type property.
// Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.
const (
	Probe_Type_Http = "http"
	Probe_Type_Tcp = "tcp"
)

// UnmarshalProbe unmarshals an instance of Probe from the specified map of raw messages.
func UnmarshalProbe(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Probe)
	err = core.UnmarshalPrimitive(m, "failure_threshold", &obj.FailureThreshold)
	if err != nil {
		err = core.SDKErrorf(err, "", "failure_threshold-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "initial_delay", &obj.InitialDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "initial_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		err = core.SDKErrorf(err, "", "interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeout-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProbePrototype : Request model for probes.
type ProbePrototype struct {
	// The number of consecutive, unsuccessful checks for the probe to be considered failed.
	FailureThreshold *int64 `json:"failure_threshold,omitempty"`

	// The amount of time in seconds to wait before the first probe check is performed.
	InitialDelay *int64 `json:"initial_delay,omitempty"`

	// The amount of time in seconds between probe checks.
	Interval *int64 `json:"interval,omitempty"`

	// The path of the HTTP request to the resource. A path is only supported for a probe with a `type` of `http`.
	Path *string `json:"path,omitempty"`

	// The port on which to probe the resource.
	Port *int64 `json:"port,omitempty"`

	// The amount of time in seconds that the probe waits for a response from the application before it times out and
	// fails.
	Timeout *int64 `json:"timeout,omitempty"`

	// Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.
	Type *string `json:"type,omitempty"`
}

// Constants associated with the ProbePrototype.Type property.
// Specifies whether to use HTTP or TCP for the probe checks. The default is TCP.
const (
	ProbePrototype_Type_Http = "http"
	ProbePrototype_Type_Tcp = "tcp"
)

// UnmarshalProbePrototype unmarshals an instance of ProbePrototype from the specified map of raw messages.
func UnmarshalProbePrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProbePrototype)
	err = core.UnmarshalPrimitive(m, "failure_threshold", &obj.FailureThreshold)
	if err != nil {
		err = core.SDKErrorf(err, "", "failure_threshold-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "initial_delay", &obj.InitialDelay)
	if err != nil {
		err = core.SDKErrorf(err, "", "initial_delay-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		err = core.SDKErrorf(err, "", "interval-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		err = core.SDKErrorf(err, "", "path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		err = core.SDKErrorf(err, "", "timeout-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the ProbePrototype
func (probePrototype *ProbePrototype) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(probePrototype.FailureThreshold) {
		_patch["failure_threshold"] = probePrototype.FailureThreshold
	}
	if !core.IsNil(probePrototype.InitialDelay) {
		_patch["initial_delay"] = probePrototype.InitialDelay
	}
	if !core.IsNil(probePrototype.Interval) {
		_patch["interval"] = probePrototype.Interval
	}
	if !core.IsNil(probePrototype.Path) {
		_patch["path"] = probePrototype.Path
	}
	if !core.IsNil(probePrototype.Port) {
		_patch["port"] = probePrototype.Port
	}
	if !core.IsNil(probePrototype.Timeout) {
		_patch["timeout"] = probePrototype.Timeout
	}
	if !core.IsNil(probePrototype.Type) {
		_patch["type"] = probePrototype.Type
	}

	return
}

// Project : Describes the model of a project.
type Project struct {
	// An alphanumeric value identifying the account ID.
	AccountID *string `json:"account_id,omitempty"`

	// The timestamp when the project was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The CRN of the project.
	Crn *string `json:"crn,omitempty"`

	// When you provision a new resource, a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The ID of the project.
	ID *string `json:"id,omitempty"`

	// The name of the project.
	Name *string `json:"name" validate:"required"`

	// The region for your project deployment. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de', 'eu-gb', 'jp-osa',
	// 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The ID of the resource group.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// The type of the project.
	ResourceType *string `json:"resource_type,omitempty"`

	// The current state of the project. For example, when the project is created and is ready for use, the status of the
	// project is active.
	Status *string `json:"status,omitempty"`
}

// Constants associated with the Project.ResourceType property.
// The type of the project.
const (
	Project_ResourceType_ProjectV2 = "project_v2"
)

// Constants associated with the Project.Status property.
// The current state of the project. For example, when the project is created and is ready for use, the status of the
// project is active.
const (
	Project_Status_Active = "active"
	Project_Status_Creating = "creating"
	Project_Status_CreationFailed = "creation_failed"
	Project_Status_Deleting = "deleting"
	Project_Status_DeletionFailed = "deletion_failed"
	Project_Status_HardDeleted = "hard_deleted"
	Project_Status_HardDeleting = "hard_deleting"
	Project_Status_HardDeletionFailed = "hard_deletion_failed"
	Project_Status_Inactive = "inactive"
	Project_Status_PendingRemoval = "pending_removal"
	Project_Status_Preparing = "preparing"
	Project_Status_SoftDeleted = "soft_deleted"
)

// UnmarshalProject unmarshals an instance of Project from the specified map of raw messages.
func UnmarshalProject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Project)
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectEgressIPAddresses : Describes the model of egress IP addresses.
type ProjectEgressIPAddresses struct {
	// List of IBM private network IP addresses.
	Private []string `json:"private,omitempty"`

	// List of public IP addresses.
	Public []string `json:"public,omitempty"`
}

// UnmarshalProjectEgressIPAddresses unmarshals an instance of ProjectEgressIPAddresses from the specified map of raw messages.
func UnmarshalProjectEgressIPAddresses(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectEgressIPAddresses)
	err = core.UnmarshalPrimitive(m, "private", &obj.Private)
	if err != nil {
		err = core.SDKErrorf(err, "", "private-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "public", &obj.Public)
	if err != nil {
		err = core.SDKErrorf(err, "", "public-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ProjectList : Contains a list of projects and pagination information.
type ProjectList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`

	// List of projects.
	Projects []Project `json:"projects" validate:"required"`
}

// UnmarshalProjectList unmarshals an instance of ProjectList from the specified map of raw messages.
func UnmarshalProjectList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "projects", &obj.Projects, UnmarshalProject)
	if err != nil {
		err = core.SDKErrorf(err, "", "projects-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ProjectList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ProjectStatusDetails : Describes the model of a project status details.
type ProjectStatusDetails struct {
	// Status of the domain created for the project.
	Domain *string `json:"domain" validate:"required"`

	// Defines whether a project is enabled for management and consumption.
	Project *string `json:"project" validate:"required"`

	// Return true when project is not VPE enabled.
	VpeNotEnabled *bool `json:"vpe_not_enabled,omitempty"`
}

// Constants associated with the ProjectStatusDetails.Domain property.
// Status of the domain created for the project.
const (
	ProjectStatusDetails_Domain_Ready = "ready"
	ProjectStatusDetails_Domain_Unknown = "unknown"
)

// Constants associated with the ProjectStatusDetails.Project property.
// Defines whether a project is enabled for management and consumption.
const (
	ProjectStatusDetails_Project_Disabled = "disabled"
	ProjectStatusDetails_Project_Enabled = "enabled"
)

// UnmarshalProjectStatusDetails unmarshals an instance of ProjectStatusDetails from the specified map of raw messages.
func UnmarshalProjectStatusDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ProjectStatusDetails)
	err = core.UnmarshalPrimitive(m, "domain", &obj.Domain)
	if err != nil {
		err = core.SDKErrorf(err, "", "domain-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project", &obj.Project)
	if err != nil {
		err = core.SDKErrorf(err, "", "project-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "vpe_not_enabled", &obj.VpeNotEnabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "vpe_not_enabled-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ReplaceConfigMapOptions : The ReplaceConfigMap options.
type ReplaceConfigMapOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your configmap.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the config map settings to be updated. Specify the version that you retrieved as entity_tag (ETag header)
	// when reading the config map. This value helps identifying parallel usage of this API. Pass * to indicate to update
	// any version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// The key-value pair for the config map. Values must be specified in `KEY=VALUE` format. Each `KEY` field must consist
	// of alphanumeric characters, `-`, `_` or `.` and must not be exceed a max length of 253 characters. Each `VALUE`
	// field can consists of any character and must not be exceed a max length of 1048576 characters.
	Data map[string]string `json:"data,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewReplaceConfigMapOptions : Instantiate ReplaceConfigMapOptions
func (*CodeEngineV2) NewReplaceConfigMapOptions(projectID string, name string, ifMatch string) *ReplaceConfigMapOptions {
	return &ReplaceConfigMapOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ReplaceConfigMapOptions) SetProjectID(projectID string) *ReplaceConfigMapOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceConfigMapOptions) SetName(name string) *ReplaceConfigMapOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceConfigMapOptions) SetIfMatch(ifMatch string) *ReplaceConfigMapOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetData : Allow user to set Data
func (_options *ReplaceConfigMapOptions) SetData(data map[string]string) *ReplaceConfigMapOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceConfigMapOptions) SetHeaders(param map[string]string) *ReplaceConfigMapOptions {
	options.Headers = param
	return options
}

// ReplaceSecretOptions : The ReplaceSecret options.
type ReplaceSecretOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your secret.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the secret settings to be updated. Specify the version that you retrieved as entity_tag (ETag header)
	// when reading the secret. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Specify the format of the secret. The format of the secret will determine how the secret is used.
	Format *string `json:"format" validate:"required"`

	// Data container that allows to specify config parameters and their values as a key-value map. Each key field must
	// consist of alphanumeric characters, `-`, `_` or `.` and must not exceed a max length of 253 characters. Each value
	// field can consists of any character and must not exceed a max length of 1048576 characters.
	Data SecretDataIntf `json:"data,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ReplaceSecretOptions.Format property.
// Specify the format of the secret. The format of the secret will determine how the secret is used.
const (
	ReplaceSecretOptions_Format_BasicAuth = "basic_auth"
	ReplaceSecretOptions_Format_Generic = "generic"
	ReplaceSecretOptions_Format_Other = "other"
	ReplaceSecretOptions_Format_Registry = "registry"
	ReplaceSecretOptions_Format_ServiceAccess = "service_access"
	ReplaceSecretOptions_Format_ServiceOperator = "service_operator"
	ReplaceSecretOptions_Format_SshAuth = "ssh_auth"
	ReplaceSecretOptions_Format_Tls = "tls"
)

// NewReplaceSecretOptions : Instantiate ReplaceSecretOptions
func (*CodeEngineV2) NewReplaceSecretOptions(projectID string, name string, ifMatch string, format string) *ReplaceSecretOptions {
	return &ReplaceSecretOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		Format: core.StringPtr(format),
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *ReplaceSecretOptions) SetProjectID(projectID string) *ReplaceSecretOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *ReplaceSecretOptions) SetName(name string) *ReplaceSecretOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *ReplaceSecretOptions) SetIfMatch(ifMatch string) *ReplaceSecretOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *ReplaceSecretOptions) SetFormat(format string) *ReplaceSecretOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetData : Allow user to set Data
func (_options *ReplaceSecretOptions) SetData(data SecretDataIntf) *ReplaceSecretOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ReplaceSecretOptions) SetHeaders(param map[string]string) *ReplaceSecretOptions {
	options.Headers = param
	return options
}

// ResourceKeyRef : The service credential associated with the secret.
type ResourceKeyRef struct {
	// ID of the service credential associated with the secret.
	ID *string `json:"id,omitempty"`

	// Name of the service credential associated with the secret.
	Name *string `json:"name,omitempty"`
}

// UnmarshalResourceKeyRef unmarshals an instance of ResourceKeyRef from the specified map of raw messages.
func UnmarshalResourceKeyRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceKeyRef)
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

// ResourceKeyRefPrototype : The service credential associated with the secret.
type ResourceKeyRefPrototype struct {
	// ID of the service credential associated with the secret.
	ID *string `json:"id,omitempty"`
}

// UnmarshalResourceKeyRefPrototype unmarshals an instance of ResourceKeyRefPrototype from the specified map of raw messages.
func UnmarshalResourceKeyRefPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResourceKeyRefPrototype)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RoleRef : A reference to the Role and Role CRN for service binding.
type RoleRef struct {
	// CRN of the IAM Role for this service access secret.
	Crn *string `json:"crn,omitempty"`

	// Role of the service credential.
	Name *string `json:"name,omitempty"`
}

// UnmarshalRoleRef unmarshals an instance of RoleRef from the specified map of raw messages.
func UnmarshalRoleRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RoleRef)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
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

// RoleRefPrototype : A reference to the Role and Role CRN for service binding.
type RoleRefPrototype struct {
	// CRN of the IAM Role for this service access secret.
	Crn *string `json:"crn,omitempty"`
}

// UnmarshalRoleRefPrototype unmarshals an instance of RoleRefPrototype from the specified map of raw messages.
func UnmarshalRoleRefPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RoleRefPrototype)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Secret : Describes the model of a secret.
type Secret struct {
	// The timestamp when the resource was created.
	CreatedAt *string `json:"created_at,omitempty"`

	// Data container that allows to specify config parameters and their values as a key-value map. Each key field must
	// consist of alphanumeric characters, `-`, `_` or `.` and must not exceed a max length of 253 characters. Each value
	// field can consists of any character and must not exceed a max length of 1048576 characters.
	Data map[string]string `json:"data,omitempty"`

	// The version of the secret instance, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Specify the format of the secret.
	Format *string `json:"format,omitempty"`

	// When you provision a new secret,  a URL is created identifying the location of the instance.
	Href *string `json:"href,omitempty"`

	// The identifier of the resource.
	ID *string `json:"id,omitempty"`

	// The name of the secret.
	Name *string `json:"name" validate:"required"`

	// The ID of the project in which the resource is located.
	ProjectID *string `json:"project_id,omitempty"`

	// The region of the project the resource is located in. Possible values: 'au-syd', 'br-sao', 'ca-tor', 'eu-de',
	// 'eu-gb', 'jp-osa', 'jp-tok', 'us-east', 'us-south'.
	Region *string `json:"region,omitempty"`

	// The type of the secret.
	ResourceType *string `json:"resource_type,omitempty"`

	// Properties for Service Access Secrets.
	ServiceAccess *ServiceAccessSecretProps `json:"service_access,omitempty"`

	// Properties for the IBM Cloud Operator Secret.
	ServiceOperator *OperatorSecretProps `json:"service_operator,omitempty"`
}

// Constants associated with the Secret.Format property.
// Specify the format of the secret.
const (
	Secret_Format_BasicAuth = "basic_auth"
	Secret_Format_Generic = "generic"
	Secret_Format_Other = "other"
	Secret_Format_Registry = "registry"
	Secret_Format_ServiceAccess = "service_access"
	Secret_Format_ServiceOperator = "service_operator"
	Secret_Format_SshAuth = "ssh_auth"
	Secret_Format_Tls = "tls"
)

// UnmarshalSecret unmarshals an instance of Secret from the specified map of raw messages.
func UnmarshalSecret(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Secret)
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "data", &obj.Data)
	if err != nil {
		err = core.SDKErrorf(err, "", "data-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		err = core.SDKErrorf(err, "", "format-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "project_id", &obj.ProjectID)
	if err != nil {
		err = core.SDKErrorf(err, "", "project_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "region", &obj.Region)
	if err != nil {
		err = core.SDKErrorf(err, "", "region-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_type", &obj.ResourceType)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_access", &obj.ServiceAccess, UnmarshalServiceAccessSecretProps)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_access-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_operator", &obj.ServiceOperator, UnmarshalOperatorSecretProps)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_operator-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretData : Data container that allows to specify config parameters and their values as a key-value map. Each key field must
// consist of alphanumeric characters, `-`, `_` or `.` and must not exceed a max length of 253 characters. Each value
// field can consists of any character and must not exceed a max length of 1048576 characters.
// This type supports additional properties of type *string.
// Models which "extend" this model:
// - SecretDataGenericSecretData
// - SecretDataBasicAuthSecretData
// - SecretDataRegistrySecretData
// - SecretDataSSHSecretData
// - SecretDataTLSSecretData
type SecretData struct {
	// Basic auth username.
	Username *string `json:"username,omitempty"`

	// Basic auth password.
	Password *string `json:"password,omitempty"`

	// Registry server.
	Server *string `json:"server,omitempty"`

	// Registry email address.
	Email *string `json:"email,omitempty"`

	// SSH key.
	SshKey *string `json:"ssh_key,omitempty"`

	// Known hosts.
	KnownHosts *string `json:"known_hosts,omitempty"`

	// The TLS certificate used in a TLS secret.
	TlsCert *string `json:"tls_cert,omitempty"`

	// The TLS key used in a TLS secret.
	TlsKey *string `json:"tls_key,omitempty"`

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}
func (*SecretData) isaSecretData() bool {
	return true
}

type SecretDataIntf interface {
	isaSecretData() bool
	SetProperty(key string, value *string)
	SetProperties(m map[string]*string)
	GetProperty(key string) *string
	GetProperties() map[string]*string
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretData.
func (o *SecretData) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretData.
func (o *SecretData) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretData.
func (o *SecretData) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretData.
func (o *SecretData) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretData
func (o *SecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Username != nil {
		m["username"] = o.Username
	}
	if o.Password != nil {
		m["password"] = o.Password
	}
	if o.Server != nil {
		m["server"] = o.Server
	}
	if o.Email != nil {
		m["email"] = o.Email
	}
	if o.SshKey != nil {
		m["ssh_key"] = o.SshKey
	}
	if o.KnownHosts != nil {
		m["known_hosts"] = o.KnownHosts
	}
	if o.TlsCert != nil {
		m["tls_cert"] = o.TlsCert
	}
	if o.TlsKey != nil {
		m["tls_key"] = o.TlsKey
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalSecretData unmarshals an instance of SecretData from the specified map of raw messages.
func UnmarshalSecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretData)
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	delete(m, "username")
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	delete(m, "password")
	err = core.UnmarshalPrimitive(m, "server", &obj.Server)
	if err != nil {
		err = core.SDKErrorf(err, "", "server-error", common.GetComponentInfo())
		return
	}
	delete(m, "server")
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	delete(m, "email")
	err = core.UnmarshalPrimitive(m, "ssh_key", &obj.SshKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "ssh_key-error", common.GetComponentInfo())
		return
	}
	delete(m, "ssh_key")
	err = core.UnmarshalPrimitive(m, "known_hosts", &obj.KnownHosts)
	if err != nil {
		err = core.SDKErrorf(err, "", "known_hosts-error", common.GetComponentInfo())
		return
	}
	delete(m, "known_hosts")
	err = core.UnmarshalPrimitive(m, "tls_cert", &obj.TlsCert)
	if err != nil {
		err = core.SDKErrorf(err, "", "tls_cert-error", common.GetComponentInfo())
		return
	}
	delete(m, "tls_cert")
	err = core.UnmarshalPrimitive(m, "tls_key", &obj.TlsKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "tls_key-error", common.GetComponentInfo())
		return
	}
	delete(m, "tls_key")
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretList : List of secret resources.
type SecretList struct {
	// Describes properties needed to retrieve the first page of a result list.
	First *ListFirstMetadata `json:"first,omitempty"`

	// Maximum number of resources per page.
	Limit *int64 `json:"limit" validate:"required"`

	// Describes properties needed to retrieve the next page of a result list.
	Next *ListNextMetadata `json:"next,omitempty"`

	// List of secrets.
	Secrets []Secret `json:"secrets" validate:"required"`
}

// UnmarshalSecretList unmarshals an instance of SecretList from the specified map of raw messages.
func UnmarshalSecretList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretList)
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalListFirstMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalListNextMetadata)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "secrets", &obj.Secrets, UnmarshalSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "secrets-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SecretList) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ServiceAccessSecretProps : Properties for Service Access Secrets.
type ServiceAccessSecretProps struct {
	// The service credential associated with the secret.
	ResourceKey *ResourceKeyRef `json:"resource_key" validate:"required"`

	// A reference to the Role and Role CRN for service binding.
	Role *RoleRef `json:"role,omitempty"`

	// The IBM Cloud service instance associated with the secret.
	ServiceInstance *ServiceInstanceRef `json:"service_instance" validate:"required"`

	// A reference to a Service ID.
	Serviceid *ServiceIDRef `json:"serviceid,omitempty"`
}

// UnmarshalServiceAccessSecretProps unmarshals an instance of ServiceAccessSecretProps from the specified map of raw messages.
func UnmarshalServiceAccessSecretProps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceAccessSecretProps)
	err = core.UnmarshalModel(m, "resource_key", &obj.ResourceKey, UnmarshalResourceKeyRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "role", &obj.Role, UnmarshalRoleRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_instance", &obj.ServiceInstance, UnmarshalServiceInstanceRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalServiceIDRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceAccessSecretPrototypeProps : Properties for Service Access Secrets.
type ServiceAccessSecretPrototypeProps struct {
	// The service credential associated with the secret.
	ResourceKey *ResourceKeyRefPrototype `json:"resource_key" validate:"required"`

	// A reference to the Role and Role CRN for service binding.
	Role *RoleRefPrototype `json:"role,omitempty"`

	// The IBM Cloud service instance associated with the secret.
	ServiceInstance *ServiceInstanceRefPrototype `json:"service_instance" validate:"required"`

	// A reference to a Service ID.
	Serviceid *ServiceIDRef `json:"serviceid,omitempty"`
}

// NewServiceAccessSecretPrototypeProps : Instantiate ServiceAccessSecretPrototypeProps (Generic Model Constructor)
func (*CodeEngineV2) NewServiceAccessSecretPrototypeProps(resourceKey *ResourceKeyRefPrototype, serviceInstance *ServiceInstanceRefPrototype) (_model *ServiceAccessSecretPrototypeProps, err error) {
	_model = &ServiceAccessSecretPrototypeProps{
		ResourceKey: resourceKey,
		ServiceInstance: serviceInstance,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalServiceAccessSecretPrototypeProps unmarshals an instance of ServiceAccessSecretPrototypeProps from the specified map of raw messages.
func UnmarshalServiceAccessSecretPrototypeProps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceAccessSecretPrototypeProps)
	err = core.UnmarshalModel(m, "resource_key", &obj.ResourceKey, UnmarshalResourceKeyRefPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_key-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "role", &obj.Role, UnmarshalRoleRefPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "role-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "service_instance", &obj.ServiceInstance, UnmarshalServiceInstanceRefPrototype)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_instance-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "serviceid", &obj.Serviceid, UnmarshalServiceIDRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "serviceid-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceIDRef : A reference to a Service ID.
type ServiceIDRef struct {
	// CRN value of a Service ID.
	Crn *string `json:"crn,omitempty"`

	// The ID of the Service ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalServiceIDRef unmarshals an instance of ServiceIDRef from the specified map of raw messages.
func UnmarshalServiceIDRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceIDRef)
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
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

// ServiceIDRefPrototype : A reference to the Service ID.
type ServiceIDRefPrototype struct {
	// The ID of the Service ID.
	ID *string `json:"id,omitempty"`
}

// UnmarshalServiceIDRefPrototype unmarshals an instance of ServiceIDRefPrototype from the specified map of raw messages.
func UnmarshalServiceIDRefPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceIDRefPrototype)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceInstanceRef : The IBM Cloud service instance associated with the secret.
type ServiceInstanceRef struct {
	// ID of the IBM Cloud service instance associated with the secret.
	ID *string `json:"id,omitempty"`

	// Type of IBM Cloud service associated with the secret.
	Type *string `json:"type,omitempty"`
}

// UnmarshalServiceInstanceRef unmarshals an instance of ServiceInstanceRef from the specified map of raw messages.
func UnmarshalServiceInstanceRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceInstanceRef)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ServiceInstanceRefPrototype : The IBM Cloud service instance associated with the secret.
type ServiceInstanceRefPrototype struct {
	// ID of the IBM Cloud service instance associated with the secret.
	ID *string `json:"id,omitempty"`
}

// UnmarshalServiceInstanceRefPrototype unmarshals an instance of ServiceInstanceRefPrototype from the specified map of raw messages.
func UnmarshalServiceInstanceRefPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ServiceInstanceRefPrototype)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateAllowedOutboundDestinationOptions : The UpdateAllowedOutboundDestination options.
type UpdateAllowedOutboundDestinationOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your allowed outbound destination.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the allowed outbound destination to be updated. Specify the version that you retrieved as entity_tag
	// (ETag header) when reading the allowed outbound destination. This value helps identifying parallel usage of this
	// API. Pass * to indicate to update any version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// AllowedOutboundDestination patch.
	AllowedOutboundDestination map[string]interface{} `json:"allowed_outbound_destination" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateAllowedOutboundDestinationOptions : Instantiate UpdateAllowedOutboundDestinationOptions
func (*CodeEngineV2) NewUpdateAllowedOutboundDestinationOptions(projectID string, name string, ifMatch string, allowedOutboundDestination map[string]interface{}) *UpdateAllowedOutboundDestinationOptions {
	return &UpdateAllowedOutboundDestinationOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		AllowedOutboundDestination: allowedOutboundDestination,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateAllowedOutboundDestinationOptions) SetProjectID(projectID string) *UpdateAllowedOutboundDestinationOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAllowedOutboundDestinationOptions) SetName(name string) *UpdateAllowedOutboundDestinationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateAllowedOutboundDestinationOptions) SetIfMatch(ifMatch string) *UpdateAllowedOutboundDestinationOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetAllowedOutboundDestination : Allow user to set AllowedOutboundDestination
func (_options *UpdateAllowedOutboundDestinationOptions) SetAllowedOutboundDestination(allowedOutboundDestination map[string]interface{}) *UpdateAllowedOutboundDestinationOptions {
	_options.AllowedOutboundDestination = allowedOutboundDestination
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAllowedOutboundDestinationOptions) SetHeaders(param map[string]string) *UpdateAllowedOutboundDestinationOptions {
	options.Headers = param
	return options
}

// UpdateAppOptions : The UpdateApp options.
type UpdateAppOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your application.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the app settings to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the app. This value helps identifying parallel usage of this API. Pass * to indicate to update any version
	// available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// App patch.
	App map[string]interface{} `json:"app" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateAppOptions : Instantiate UpdateAppOptions
func (*CodeEngineV2) NewUpdateAppOptions(projectID string, name string, ifMatch string, app map[string]interface{}) *UpdateAppOptions {
	return &UpdateAppOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		App: app,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateAppOptions) SetProjectID(projectID string) *UpdateAppOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateAppOptions) SetName(name string) *UpdateAppOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateAppOptions) SetIfMatch(ifMatch string) *UpdateAppOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetApp : Allow user to set App
func (_options *UpdateAppOptions) SetApp(app map[string]interface{}) *UpdateAppOptions {
	_options.App = app
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateAppOptions) SetHeaders(param map[string]string) *UpdateAppOptions {
	options.Headers = param
	return options
}

// UpdateBuildOptions : The UpdateBuild options.
type UpdateBuildOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your build.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the build settings to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the build. This value helps identifying parallel usage of this API. Pass * to indicate to update any version
	// available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Build patch.
	Build map[string]interface{} `json:"build" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateBuildOptions : Instantiate UpdateBuildOptions
func (*CodeEngineV2) NewUpdateBuildOptions(projectID string, name string, ifMatch string, build map[string]interface{}) *UpdateBuildOptions {
	return &UpdateBuildOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		Build: build,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateBuildOptions) SetProjectID(projectID string) *UpdateBuildOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateBuildOptions) SetName(name string) *UpdateBuildOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateBuildOptions) SetIfMatch(ifMatch string) *UpdateBuildOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetBuild : Allow user to set Build
func (_options *UpdateBuildOptions) SetBuild(build map[string]interface{}) *UpdateBuildOptions {
	_options.Build = build
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateBuildOptions) SetHeaders(param map[string]string) *UpdateBuildOptions {
	options.Headers = param
	return options
}

// UpdateDomainMappingOptions : The UpdateDomainMapping options.
type UpdateDomainMappingOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your domain mapping.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the domain mapping to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the domain mapping. This value helps identify parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// DomainMapping patch.
	DomainMapping map[string]interface{} `json:"domain_mapping" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateDomainMappingOptions : Instantiate UpdateDomainMappingOptions
func (*CodeEngineV2) NewUpdateDomainMappingOptions(projectID string, name string, ifMatch string, domainMapping map[string]interface{}) *UpdateDomainMappingOptions {
	return &UpdateDomainMappingOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		DomainMapping: domainMapping,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateDomainMappingOptions) SetProjectID(projectID string) *UpdateDomainMappingOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateDomainMappingOptions) SetName(name string) *UpdateDomainMappingOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateDomainMappingOptions) SetIfMatch(ifMatch string) *UpdateDomainMappingOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetDomainMapping : Allow user to set DomainMapping
func (_options *UpdateDomainMappingOptions) SetDomainMapping(domainMapping map[string]interface{}) *UpdateDomainMappingOptions {
	_options.DomainMapping = domainMapping
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateDomainMappingOptions) SetHeaders(param map[string]string) *UpdateDomainMappingOptions {
	options.Headers = param
	return options
}

// UpdateFunctionOptions : The UpdateFunction options.
type UpdateFunctionOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your function.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the function settings to be updated. Specify the version that you retrieved as entity_tag (ETag header)
	// when reading the function. This value helps identifying parallel usage of this API. Pass * to indicate to update any
	// version available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Function patch.
	Function map[string]interface{} `json:"function" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateFunctionOptions : Instantiate UpdateFunctionOptions
func (*CodeEngineV2) NewUpdateFunctionOptions(projectID string, name string, ifMatch string, function map[string]interface{}) *UpdateFunctionOptions {
	return &UpdateFunctionOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		Function: function,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateFunctionOptions) SetProjectID(projectID string) *UpdateFunctionOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateFunctionOptions) SetName(name string) *UpdateFunctionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateFunctionOptions) SetIfMatch(ifMatch string) *UpdateFunctionOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetFunction : Allow user to set Function
func (_options *UpdateFunctionOptions) SetFunction(function map[string]interface{}) *UpdateFunctionOptions {
	_options.Function = function
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateFunctionOptions) SetHeaders(param map[string]string) *UpdateFunctionOptions {
	options.Headers = param
	return options
}

// UpdateJobOptions : The UpdateJob options.
type UpdateJobOptions struct {
	// The ID of the project.
	ProjectID *string `json:"project_id" validate:"required,ne="`

	// The name of your job.
	Name *string `json:"name" validate:"required,ne="`

	// Version of the job settings to be updated. Specify the version that you retrieved as entity_tag (ETag header) when
	// reading the job. This value helps identifying parallel usage of this API. Pass * to indicate to update any version
	// available. This might result in stale updates.
	IfMatch *string `json:"If-Match" validate:"required"`

	// Job patch prototype.
	Job map[string]interface{} `json:"job" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateJobOptions : Instantiate UpdateJobOptions
func (*CodeEngineV2) NewUpdateJobOptions(projectID string, name string, ifMatch string, job map[string]interface{}) *UpdateJobOptions {
	return &UpdateJobOptions{
		ProjectID: core.StringPtr(projectID),
		Name: core.StringPtr(name),
		IfMatch: core.StringPtr(ifMatch),
		Job: job,
	}
}

// SetProjectID : Allow user to set ProjectID
func (_options *UpdateJobOptions) SetProjectID(projectID string) *UpdateJobOptions {
	_options.ProjectID = core.StringPtr(projectID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateJobOptions) SetName(name string) *UpdateJobOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetIfMatch : Allow user to set IfMatch
func (_options *UpdateJobOptions) SetIfMatch(ifMatch string) *UpdateJobOptions {
	_options.IfMatch = core.StringPtr(ifMatch)
	return _options
}

// SetJob : Allow user to set Job
func (_options *UpdateJobOptions) SetJob(job map[string]interface{}) *UpdateJobOptions {
	_options.Job = job
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateJobOptions) SetHeaders(param map[string]string) *UpdateJobOptions {
	options.Headers = param
	return options
}

// VolumeMount : Response model of a volume mount.
type VolumeMount struct {
	// The path that should be mounted.
	MountPath *string `json:"mount_path" validate:"required"`

	// The name of the mount.
	Name *string `json:"name" validate:"required"`

	// The name of the referenced secret or config map.
	Reference *string `json:"reference" validate:"required"`

	// Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the VolumeMount.Type property.
// Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.
const (
	VolumeMount_Type_ConfigMap = "config_map"
	VolumeMount_Type_Secret = "secret"
)

// UnmarshalVolumeMount unmarshals an instance of VolumeMount from the specified map of raw messages.
func UnmarshalVolumeMount(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VolumeMount)
	err = core.UnmarshalPrimitive(m, "mount_path", &obj.MountPath)
	if err != nil {
		err = core.SDKErrorf(err, "", "mount_path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// VolumeMountPrototype : Prototype model of a volume mount.
type VolumeMountPrototype struct {
	// The path that should be mounted.
	MountPath *string `json:"mount_path" validate:"required"`

	// Optional name of the mount. If not set, it will be generated based on the `ref` and a random ID. In case the `ref`
	// is longer than 58 characters, it will be cut off.
	Name *string `json:"name,omitempty"`

	// The name of the referenced secret or config map.
	Reference *string `json:"reference" validate:"required"`

	// Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.
	Type *string `json:"type" validate:"required"`
}

// Constants associated with the VolumeMountPrototype.Type property.
// Specify the type of the volume mount. Allowed types are: 'config_map', 'secret'.
const (
	VolumeMountPrototype_Type_ConfigMap = "config_map"
	VolumeMountPrototype_Type_Secret = "secret"
)

// NewVolumeMountPrototype : Instantiate VolumeMountPrototype (Generic Model Constructor)
func (*CodeEngineV2) NewVolumeMountPrototype(mountPath string, reference string, typeVar string) (_model *VolumeMountPrototype, err error) {
	_model = &VolumeMountPrototype{
		MountPath: core.StringPtr(mountPath),
		Reference: core.StringPtr(reference),
		Type: core.StringPtr(typeVar),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalVolumeMountPrototype unmarshals an instance of VolumeMountPrototype from the specified map of raw messages.
func UnmarshalVolumeMountPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VolumeMountPrototype)
	err = core.UnmarshalPrimitive(m, "mount_path", &obj.MountPath)
	if err != nil {
		err = core.SDKErrorf(err, "", "mount_path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "reference", &obj.Reference)
	if err != nil {
		err = core.SDKErrorf(err, "", "reference-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// asPatch returns a generic map representation of the VolumeMountPrototype
func (volumeMountPrototype *VolumeMountPrototype) asPatch() (_patch map[string]interface{}) {
	_patch = map[string]interface{}{}
	if !core.IsNil(volumeMountPrototype.MountPath) {
		_patch["mount_path"] = volumeMountPrototype.MountPath
	}
	if !core.IsNil(volumeMountPrototype.Name) {
		_patch["name"] = volumeMountPrototype.Name
	}
	if !core.IsNil(volumeMountPrototype.Reference) {
		_patch["reference"] = volumeMountPrototype.Reference
	}
	if !core.IsNil(volumeMountPrototype.Type) {
		_patch["type"] = volumeMountPrototype.Type
	}

	return
}

// AllowedOutboundDestinationPatchCidrBlockDataPatch : Update an allowed outbound destination by using a CIDR block.
// This model "extends" AllowedOutboundDestinationPatch
type AllowedOutboundDestinationPatchCidrBlockDataPatch struct {
	// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
	Type *string `json:"type,omitempty"`

	// The IPv4 address range.
	CidrBlock *string `json:"cidr_block,omitempty"`
}

// Constants associated with the AllowedOutboundDestinationPatchCidrBlockDataPatch.Type property.
// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
const (
	AllowedOutboundDestinationPatchCidrBlockDataPatch_Type_CidrBlock = "cidr_block"
)

func (*AllowedOutboundDestinationPatchCidrBlockDataPatch) isaAllowedOutboundDestinationPatch() bool {
	return true
}

// UnmarshalAllowedOutboundDestinationPatchCidrBlockDataPatch unmarshals an instance of AllowedOutboundDestinationPatchCidrBlockDataPatch from the specified map of raw messages.
func UnmarshalAllowedOutboundDestinationPatchCidrBlockDataPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestinationPatchCidrBlockDataPatch)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cidr_block", &obj.CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "cidr_block-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the AllowedOutboundDestinationPatchCidrBlockDataPatch
func (allowedOutboundDestinationPatchCidrBlockDataPatch *AllowedOutboundDestinationPatchCidrBlockDataPatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(allowedOutboundDestinationPatchCidrBlockDataPatch.Type) {
		_patch["type"] = allowedOutboundDestinationPatchCidrBlockDataPatch.Type
	}
	if !core.IsNil(allowedOutboundDestinationPatchCidrBlockDataPatch.CidrBlock) {
		_patch["cidr_block"] = allowedOutboundDestinationPatchCidrBlockDataPatch.CidrBlock
	}

	return
}

// AllowedOutboundDestinationPrototypeCidrBlockDataPrototype : Create an allowed outbound destination by using a CIDR block.
// This model "extends" AllowedOutboundDestinationPrototype
type AllowedOutboundDestinationPrototypeCidrBlockDataPrototype struct {
	// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
	Type *string `json:"type" validate:"required"`

	// The IPv4 address range.
	CidrBlock *string `json:"cidr_block" validate:"required"`

	// The name of the CIDR block.
	Name *string `json:"name" validate:"required"`
}

// Constants associated with the AllowedOutboundDestinationPrototypeCidrBlockDataPrototype.Type property.
// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
const (
	AllowedOutboundDestinationPrototypeCidrBlockDataPrototype_Type_CidrBlock = "cidr_block"
)

// NewAllowedOutboundDestinationPrototypeCidrBlockDataPrototype : Instantiate AllowedOutboundDestinationPrototypeCidrBlockDataPrototype (Generic Model Constructor)
func (*CodeEngineV2) NewAllowedOutboundDestinationPrototypeCidrBlockDataPrototype(typeVar string, cidrBlock string, name string) (_model *AllowedOutboundDestinationPrototypeCidrBlockDataPrototype, err error) {
	_model = &AllowedOutboundDestinationPrototypeCidrBlockDataPrototype{
		Type: core.StringPtr(typeVar),
		CidrBlock: core.StringPtr(cidrBlock),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*AllowedOutboundDestinationPrototypeCidrBlockDataPrototype) isaAllowedOutboundDestinationPrototype() bool {
	return true
}

// UnmarshalAllowedOutboundDestinationPrototypeCidrBlockDataPrototype unmarshals an instance of AllowedOutboundDestinationPrototypeCidrBlockDataPrototype from the specified map of raw messages.
func UnmarshalAllowedOutboundDestinationPrototypeCidrBlockDataPrototype(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestinationPrototypeCidrBlockDataPrototype)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cidr_block", &obj.CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "cidr_block-error", common.GetComponentInfo())
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

// AllowedOutboundDestinationCidrBlockData : Allowed outbound destination CIDR block.
// This model "extends" AllowedOutboundDestination
type AllowedOutboundDestinationCidrBlockData struct {
	// The version of the allowed outbound destination, which is used to achieve optimistic locking.
	EntityTag *string `json:"entity_tag" validate:"required"`

	// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
	Type *string `json:"type" validate:"required"`

	// The IPv4 address range.
	CidrBlock *string `json:"cidr_block" validate:"required"`

	// The name of the CIDR block.
	Name *string `json:"name" validate:"required"`
}

// Constants associated with the AllowedOutboundDestinationCidrBlockData.Type property.
// Specify the type of the allowed outbound destination. Allowed types are: 'cidr_block'.
const (
	AllowedOutboundDestinationCidrBlockData_Type_CidrBlock = "cidr_block"
)

func (*AllowedOutboundDestinationCidrBlockData) isaAllowedOutboundDestination() bool {
	return true
}

// UnmarshalAllowedOutboundDestinationCidrBlockData unmarshals an instance of AllowedOutboundDestinationCidrBlockData from the specified map of raw messages.
func UnmarshalAllowedOutboundDestinationCidrBlockData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AllowedOutboundDestinationCidrBlockData)
	err = core.UnmarshalPrimitive(m, "entity_tag", &obj.EntityTag)
	if err != nil {
		err = core.SDKErrorf(err, "", "entity_tag-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "cidr_block", &obj.CidrBlock)
	if err != nil {
		err = core.SDKErrorf(err, "", "cidr_block-error", common.GetComponentInfo())
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

// SecretDataBasicAuthSecretData : SecretDataBasicAuthSecretData struct
// This type supports additional properties of type *string.
// This model "extends" SecretData
type SecretDataBasicAuthSecretData struct {
	// Basic auth username.
	Username *string `json:"username" validate:"required"`

	// Basic auth password.
	Password *string `json:"password" validate:"required"`

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}

// NewSecretDataBasicAuthSecretData : Instantiate SecretDataBasicAuthSecretData (Generic Model Constructor)
func (*CodeEngineV2) NewSecretDataBasicAuthSecretData(username string, password string) (_model *SecretDataBasicAuthSecretData, err error) {
	_model = &SecretDataBasicAuthSecretData{
		Username: core.StringPtr(username),
		Password: core.StringPtr(password),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*SecretDataBasicAuthSecretData) isaSecretData() bool {
	return true
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretDataBasicAuthSecretData.
func (o *SecretDataBasicAuthSecretData) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretDataBasicAuthSecretData.
func (o *SecretDataBasicAuthSecretData) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretDataBasicAuthSecretData.
func (o *SecretDataBasicAuthSecretData) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretDataBasicAuthSecretData.
func (o *SecretDataBasicAuthSecretData) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretDataBasicAuthSecretData
func (o *SecretDataBasicAuthSecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Username != nil {
		m["username"] = o.Username
	}
	if o.Password != nil {
		m["password"] = o.Password
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalSecretDataBasicAuthSecretData unmarshals an instance of SecretDataBasicAuthSecretData from the specified map of raw messages.
func UnmarshalSecretDataBasicAuthSecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretDataBasicAuthSecretData)
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	delete(m, "username")
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	delete(m, "password")
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretDataGenericSecretData : Data container that allows to specify config parameters and their values as a key-value map. Each key field must
// consist of alphanumeric characters, `-`, `_` or `.` and must not be exceed a max length of 253 characters. Each value
// field can consists of any character and must not be exceed a max length of 1048576 characters.
// This type supports additional properties of type *string.
// This model "extends" SecretData
type SecretDataGenericSecretData struct {

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}

func (*SecretDataGenericSecretData) isaSecretData() bool {
	return true
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretDataGenericSecretData.
func (o *SecretDataGenericSecretData) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretDataGenericSecretData.
func (o *SecretDataGenericSecretData) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretDataGenericSecretData.
func (o *SecretDataGenericSecretData) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretDataGenericSecretData.
func (o *SecretDataGenericSecretData) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretDataGenericSecretData
func (o *SecretDataGenericSecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalSecretDataGenericSecretData unmarshals an instance of SecretDataGenericSecretData from the specified map of raw messages.
func UnmarshalSecretDataGenericSecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretDataGenericSecretData)
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretDataRegistrySecretData : SecretDataRegistrySecretData struct
// This type supports additional properties of type *string.
// This model "extends" SecretData
type SecretDataRegistrySecretData struct {
	// Registry username.
	Username *string `json:"username" validate:"required"`

	// Registry password.
	Password *string `json:"password" validate:"required"`

	// Registry server.
	Server *string `json:"server" validate:"required"`

	// Registry email address.
	Email *string `json:"email,omitempty"`

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}

// NewSecretDataRegistrySecretData : Instantiate SecretDataRegistrySecretData (Generic Model Constructor)
func (*CodeEngineV2) NewSecretDataRegistrySecretData(username string, password string, server string) (_model *SecretDataRegistrySecretData, err error) {
	_model = &SecretDataRegistrySecretData{
		Username: core.StringPtr(username),
		Password: core.StringPtr(password),
		Server: core.StringPtr(server),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*SecretDataRegistrySecretData) isaSecretData() bool {
	return true
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretDataRegistrySecretData.
func (o *SecretDataRegistrySecretData) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretDataRegistrySecretData.
func (o *SecretDataRegistrySecretData) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretDataRegistrySecretData.
func (o *SecretDataRegistrySecretData) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretDataRegistrySecretData.
func (o *SecretDataRegistrySecretData) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretDataRegistrySecretData
func (o *SecretDataRegistrySecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.Username != nil {
		m["username"] = o.Username
	}
	if o.Password != nil {
		m["password"] = o.Password
	}
	if o.Server != nil {
		m["server"] = o.Server
	}
	if o.Email != nil {
		m["email"] = o.Email
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalSecretDataRegistrySecretData unmarshals an instance of SecretDataRegistrySecretData from the specified map of raw messages.
func UnmarshalSecretDataRegistrySecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretDataRegistrySecretData)
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	delete(m, "username")
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	delete(m, "password")
	err = core.UnmarshalPrimitive(m, "server", &obj.Server)
	if err != nil {
		err = core.SDKErrorf(err, "", "server-error", common.GetComponentInfo())
		return
	}
	delete(m, "server")
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
		return
	}
	delete(m, "email")
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretDataSSHSecretData : Secret Data field used by SSH secrets.
// This type supports additional properties of type *string.
// This model "extends" SecretData
type SecretDataSSHSecretData struct {
	// SSH key.
	SshKey *string `json:"ssh_key" validate:"required"`

	// Known hosts.
	KnownHosts *string `json:"known_hosts,omitempty"`

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}

// NewSecretDataSSHSecretData : Instantiate SecretDataSSHSecretData (Generic Model Constructor)
func (*CodeEngineV2) NewSecretDataSSHSecretData(sshKey string) (_model *SecretDataSSHSecretData, err error) {
	_model = &SecretDataSSHSecretData{
		SshKey: core.StringPtr(sshKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*SecretDataSSHSecretData) isaSecretData() bool {
	return true
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretDataSSHSecretData.
func (o *SecretDataSSHSecretData) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretDataSSHSecretData.
func (o *SecretDataSSHSecretData) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretDataSSHSecretData.
func (o *SecretDataSSHSecretData) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretDataSSHSecretData.
func (o *SecretDataSSHSecretData) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretDataSSHSecretData
func (o *SecretDataSSHSecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.SshKey != nil {
		m["ssh_key"] = o.SshKey
	}
	if o.KnownHosts != nil {
		m["known_hosts"] = o.KnownHosts
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalSecretDataSSHSecretData unmarshals an instance of SecretDataSSHSecretData from the specified map of raw messages.
func UnmarshalSecretDataSSHSecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretDataSSHSecretData)
	err = core.UnmarshalPrimitive(m, "ssh_key", &obj.SshKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "ssh_key-error", common.GetComponentInfo())
		return
	}
	delete(m, "ssh_key")
	err = core.UnmarshalPrimitive(m, "known_hosts", &obj.KnownHosts)
	if err != nil {
		err = core.SDKErrorf(err, "", "known_hosts-error", common.GetComponentInfo())
		return
	}
	delete(m, "known_hosts")
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SecretDataTLSSecretData : SecretDataTLSSecretData struct
// This type supports additional properties of type *string.
// This model "extends" SecretData
type SecretDataTLSSecretData struct {
	// The TLS certificate used in a TLS secret.
	TlsCert *string `json:"tls_cert" validate:"required"`

	// The TLS key used in a TLS secret.
	TlsKey *string `json:"tls_key" validate:"required"`

	// Allows users to set arbitrary properties of type *string.
	additionalProperties map[string]*string
}

// NewSecretDataTLSSecretData : Instantiate SecretDataTLSSecretData (Generic Model Constructor)
func (*CodeEngineV2) NewSecretDataTLSSecretData(tlsCert string, tlsKey string) (_model *SecretDataTLSSecretData, err error) {
	_model = &SecretDataTLSSecretData{
		TlsCert: core.StringPtr(tlsCert),
		TlsKey: core.StringPtr(tlsKey),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*SecretDataTLSSecretData) isaSecretData() bool {
	return true
}

// SetProperty allows the user to set an arbitrary property on an instance of SecretDataTLSSecretData.
func (o *SecretDataTLSSecretData) SetProperty(key string, value *string) {
	if o.additionalProperties == nil {
		o.additionalProperties = make(map[string]*string)
	}
	o.additionalProperties[key] = value
}

// SetProperties allows the user to set a map of arbitrary properties on an instance of SecretDataTLSSecretData.
func (o *SecretDataTLSSecretData) SetProperties(m map[string]*string) {
	o.additionalProperties = make(map[string]*string)
	for k, v := range m {
		o.additionalProperties[k] = v
	}
}

// GetProperty allows the user to retrieve an arbitrary property from an instance of SecretDataTLSSecretData.
func (o *SecretDataTLSSecretData) GetProperty(key string) *string {
	return o.additionalProperties[key]
}

// GetProperties allows the user to retrieve the map of arbitrary properties from an instance of SecretDataTLSSecretData.
func (o *SecretDataTLSSecretData) GetProperties() map[string]*string {
	return o.additionalProperties
}

// MarshalJSON performs custom serialization for instances of SecretDataTLSSecretData
func (o *SecretDataTLSSecretData) MarshalJSON() (buffer []byte, err error) {
	m := make(map[string]interface{})
	if len(o.additionalProperties) > 0 {
		for k, v := range o.additionalProperties {
			m[k] = v
		}
	}
	if o.TlsCert != nil {
		m["tls_cert"] = o.TlsCert
	}
	if o.TlsKey != nil {
		m["tls_key"] = o.TlsKey
	}
	buffer, err = json.Marshal(m)
	if err != nil {
		err = core.SDKErrorf(err, "", "model-marshal", common.GetComponentInfo())
	}
	return
}

// UnmarshalSecretDataTLSSecretData unmarshals an instance of SecretDataTLSSecretData from the specified map of raw messages.
func UnmarshalSecretDataTLSSecretData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SecretDataTLSSecretData)
	err = core.UnmarshalPrimitive(m, "tls_cert", &obj.TlsCert)
	if err != nil {
		err = core.SDKErrorf(err, "", "tls_cert-error", common.GetComponentInfo())
		return
	}
	delete(m, "tls_cert")
	err = core.UnmarshalPrimitive(m, "tls_key", &obj.TlsKey)
	if err != nil {
		err = core.SDKErrorf(err, "", "tls_key-error", common.GetComponentInfo())
		return
	}
	delete(m, "tls_key")
	for k := range m {
		var v *string
		e := core.UnmarshalPrimitive(m, k, &v)
		if e != nil {
			err = core.SDKErrorf(e, "", "additional-properties-error", common.GetComponentInfo())
			return
		}
		obj.SetProperty(k, v)
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

//
// ProjectsPager can be used to simplify the use of the "ListProjects" method.
//
type ProjectsPager struct {
	hasNext bool
	options *ListProjectsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewProjectsPager returns a new ProjectsPager instance.
func (codeEngine *CodeEngineV2) NewProjectsPager(options *ListProjectsOptions) (pager *ProjectsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListProjectsOptions = *options
	pager = &ProjectsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ProjectsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ProjectsPager) GetNextWithContext(ctx context.Context) (page []Project, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListProjectsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Projects

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ProjectsPager) GetAllWithContext(ctx context.Context) (allItems []Project, err error) {
	for pager.HasNext() {
		var nextPage []Project
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
func (pager *ProjectsPager) GetNext() (page []Project, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ProjectsPager) GetAll() (allItems []Project, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// AllowedOutboundDestinationPager can be used to simplify the use of the "ListAllowedOutboundDestination" method.
//
type AllowedOutboundDestinationPager struct {
	hasNext bool
	options *ListAllowedOutboundDestinationOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewAllowedOutboundDestinationPager returns a new AllowedOutboundDestinationPager instance.
func (codeEngine *CodeEngineV2) NewAllowedOutboundDestinationPager(options *ListAllowedOutboundDestinationOptions) (pager *AllowedOutboundDestinationPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListAllowedOutboundDestinationOptions = *options
	pager = &AllowedOutboundDestinationPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AllowedOutboundDestinationPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AllowedOutboundDestinationPager) GetNextWithContext(ctx context.Context) (page []AllowedOutboundDestinationIntf, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListAllowedOutboundDestinationWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.AllowedOutboundDestinations

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AllowedOutboundDestinationPager) GetAllWithContext(ctx context.Context) (allItems []AllowedOutboundDestinationIntf, err error) {
	for pager.HasNext() {
		var nextPage []AllowedOutboundDestinationIntf
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
func (pager *AllowedOutboundDestinationPager) GetNext() (page []AllowedOutboundDestinationIntf, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AllowedOutboundDestinationPager) GetAll() (allItems []AllowedOutboundDestinationIntf, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// AppsPager can be used to simplify the use of the "ListApps" method.
//
type AppsPager struct {
	hasNext bool
	options *ListAppsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewAppsPager returns a new AppsPager instance.
func (codeEngine *CodeEngineV2) NewAppsPager(options *ListAppsOptions) (pager *AppsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListAppsOptions = *options
	pager = &AppsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AppsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AppsPager) GetNextWithContext(ctx context.Context) (page []App, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListAppsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Apps

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AppsPager) GetAllWithContext(ctx context.Context) (allItems []App, err error) {
	for pager.HasNext() {
		var nextPage []App
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
func (pager *AppsPager) GetNext() (page []App, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AppsPager) GetAll() (allItems []App, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// AppRevisionsPager can be used to simplify the use of the "ListAppRevisions" method.
//
type AppRevisionsPager struct {
	hasNext bool
	options *ListAppRevisionsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewAppRevisionsPager returns a new AppRevisionsPager instance.
func (codeEngine *CodeEngineV2) NewAppRevisionsPager(options *ListAppRevisionsOptions) (pager *AppRevisionsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListAppRevisionsOptions = *options
	pager = &AppRevisionsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AppRevisionsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AppRevisionsPager) GetNextWithContext(ctx context.Context) (page []AppRevision, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListAppRevisionsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Revisions

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AppRevisionsPager) GetAllWithContext(ctx context.Context) (allItems []AppRevision, err error) {
	for pager.HasNext() {
		var nextPage []AppRevision
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
func (pager *AppRevisionsPager) GetNext() (page []AppRevision, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AppRevisionsPager) GetAll() (allItems []AppRevision, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// AppInstancesPager can be used to simplify the use of the "ListAppInstances" method.
//
type AppInstancesPager struct {
	hasNext bool
	options *ListAppInstancesOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewAppInstancesPager returns a new AppInstancesPager instance.
func (codeEngine *CodeEngineV2) NewAppInstancesPager(options *ListAppInstancesOptions) (pager *AppInstancesPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListAppInstancesOptions = *options
	pager = &AppInstancesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *AppInstancesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *AppInstancesPager) GetNextWithContext(ctx context.Context) (page []AppInstance, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListAppInstancesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Instances

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *AppInstancesPager) GetAllWithContext(ctx context.Context) (allItems []AppInstance, err error) {
	for pager.HasNext() {
		var nextPage []AppInstance
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
func (pager *AppInstancesPager) GetNext() (page []AppInstance, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *AppInstancesPager) GetAll() (allItems []AppInstance, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// JobsPager can be used to simplify the use of the "ListJobs" method.
//
type JobsPager struct {
	hasNext bool
	options *ListJobsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewJobsPager returns a new JobsPager instance.
func (codeEngine *CodeEngineV2) NewJobsPager(options *ListJobsOptions) (pager *JobsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListJobsOptions = *options
	pager = &JobsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *JobsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *JobsPager) GetNextWithContext(ctx context.Context) (page []Job, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListJobsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Jobs

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *JobsPager) GetAllWithContext(ctx context.Context) (allItems []Job, err error) {
	for pager.HasNext() {
		var nextPage []Job
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
func (pager *JobsPager) GetNext() (page []Job, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *JobsPager) GetAll() (allItems []Job, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// JobRunsPager can be used to simplify the use of the "ListJobRuns" method.
//
type JobRunsPager struct {
	hasNext bool
	options *ListJobRunsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewJobRunsPager returns a new JobRunsPager instance.
func (codeEngine *CodeEngineV2) NewJobRunsPager(options *ListJobRunsOptions) (pager *JobRunsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListJobRunsOptions = *options
	pager = &JobRunsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *JobRunsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *JobRunsPager) GetNextWithContext(ctx context.Context) (page []JobRun, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListJobRunsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.JobRuns

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *JobRunsPager) GetAllWithContext(ctx context.Context) (allItems []JobRun, err error) {
	for pager.HasNext() {
		var nextPage []JobRun
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
func (pager *JobRunsPager) GetNext() (page []JobRun, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *JobRunsPager) GetAll() (allItems []JobRun, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// FunctionsPager can be used to simplify the use of the "ListFunctions" method.
//
type FunctionsPager struct {
	hasNext bool
	options *ListFunctionsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewFunctionsPager returns a new FunctionsPager instance.
func (codeEngine *CodeEngineV2) NewFunctionsPager(options *ListFunctionsOptions) (pager *FunctionsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListFunctionsOptions = *options
	pager = &FunctionsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *FunctionsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *FunctionsPager) GetNextWithContext(ctx context.Context) (page []Function, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListFunctionsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Functions

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *FunctionsPager) GetAllWithContext(ctx context.Context) (allItems []Function, err error) {
	for pager.HasNext() {
		var nextPage []Function
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
func (pager *FunctionsPager) GetNext() (page []Function, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *FunctionsPager) GetAll() (allItems []Function, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// BindingsPager can be used to simplify the use of the "ListBindings" method.
//
type BindingsPager struct {
	hasNext bool
	options *ListBindingsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewBindingsPager returns a new BindingsPager instance.
func (codeEngine *CodeEngineV2) NewBindingsPager(options *ListBindingsOptions) (pager *BindingsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListBindingsOptions = *options
	pager = &BindingsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *BindingsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *BindingsPager) GetNextWithContext(ctx context.Context) (page []Binding, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListBindingsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Bindings

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *BindingsPager) GetAllWithContext(ctx context.Context) (allItems []Binding, err error) {
	for pager.HasNext() {
		var nextPage []Binding
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
func (pager *BindingsPager) GetNext() (page []Binding, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *BindingsPager) GetAll() (allItems []Binding, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// BuildsPager can be used to simplify the use of the "ListBuilds" method.
//
type BuildsPager struct {
	hasNext bool
	options *ListBuildsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewBuildsPager returns a new BuildsPager instance.
func (codeEngine *CodeEngineV2) NewBuildsPager(options *ListBuildsOptions) (pager *BuildsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListBuildsOptions = *options
	pager = &BuildsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *BuildsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *BuildsPager) GetNextWithContext(ctx context.Context) (page []Build, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListBuildsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Builds

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *BuildsPager) GetAllWithContext(ctx context.Context) (allItems []Build, err error) {
	for pager.HasNext() {
		var nextPage []Build
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
func (pager *BuildsPager) GetNext() (page []Build, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *BuildsPager) GetAll() (allItems []Build, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// BuildRunsPager can be used to simplify the use of the "ListBuildRuns" method.
//
type BuildRunsPager struct {
	hasNext bool
	options *ListBuildRunsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewBuildRunsPager returns a new BuildRunsPager instance.
func (codeEngine *CodeEngineV2) NewBuildRunsPager(options *ListBuildRunsOptions) (pager *BuildRunsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListBuildRunsOptions = *options
	pager = &BuildRunsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *BuildRunsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *BuildRunsPager) GetNextWithContext(ctx context.Context) (page []BuildRun, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListBuildRunsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.BuildRuns

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *BuildRunsPager) GetAllWithContext(ctx context.Context) (allItems []BuildRun, err error) {
	for pager.HasNext() {
		var nextPage []BuildRun
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
func (pager *BuildRunsPager) GetNext() (page []BuildRun, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *BuildRunsPager) GetAll() (allItems []BuildRun, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// DomainMappingsPager can be used to simplify the use of the "ListDomainMappings" method.
//
type DomainMappingsPager struct {
	hasNext bool
	options *ListDomainMappingsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewDomainMappingsPager returns a new DomainMappingsPager instance.
func (codeEngine *CodeEngineV2) NewDomainMappingsPager(options *ListDomainMappingsOptions) (pager *DomainMappingsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListDomainMappingsOptions = *options
	pager = &DomainMappingsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *DomainMappingsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *DomainMappingsPager) GetNextWithContext(ctx context.Context) (page []DomainMapping, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListDomainMappingsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.DomainMappings

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *DomainMappingsPager) GetAllWithContext(ctx context.Context) (allItems []DomainMapping, err error) {
	for pager.HasNext() {
		var nextPage []DomainMapping
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
func (pager *DomainMappingsPager) GetNext() (page []DomainMapping, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *DomainMappingsPager) GetAll() (allItems []DomainMapping, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// ConfigMapsPager can be used to simplify the use of the "ListConfigMaps" method.
//
type ConfigMapsPager struct {
	hasNext bool
	options *ListConfigMapsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewConfigMapsPager returns a new ConfigMapsPager instance.
func (codeEngine *CodeEngineV2) NewConfigMapsPager(options *ListConfigMapsOptions) (pager *ConfigMapsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListConfigMapsOptions = *options
	pager = &ConfigMapsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ConfigMapsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ConfigMapsPager) GetNextWithContext(ctx context.Context) (page []ConfigMap, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListConfigMapsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.ConfigMaps

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ConfigMapsPager) GetAllWithContext(ctx context.Context) (allItems []ConfigMap, err error) {
	for pager.HasNext() {
		var nextPage []ConfigMap
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
func (pager *ConfigMapsPager) GetNext() (page []ConfigMap, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ConfigMapsPager) GetAll() (allItems []ConfigMap, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// SecretsPager can be used to simplify the use of the "ListSecrets" method.
//
type SecretsPager struct {
	hasNext bool
	options *ListSecretsOptions
	client  *CodeEngineV2
	pageContext struct {
		next *string
	}
}

// NewSecretsPager returns a new SecretsPager instance.
func (codeEngine *CodeEngineV2) NewSecretsPager(options *ListSecretsOptions) (pager *SecretsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListSecretsOptions = *options
	pager = &SecretsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  codeEngine,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SecretsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SecretsPager) GetNextWithContext(ctx context.Context) (page []Secret, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListSecretsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *string
	if result.Next != nil {
		next = result.Next.Start
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Secrets

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SecretsPager) GetAllWithContext(ctx context.Context) (allItems []Secret, err error) {
	for pager.HasNext() {
		var nextPage []Secret
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
func (pager *SecretsPager) GetNext() (page []Secret, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SecretsPager) GetAll() (allItems []Secret, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
