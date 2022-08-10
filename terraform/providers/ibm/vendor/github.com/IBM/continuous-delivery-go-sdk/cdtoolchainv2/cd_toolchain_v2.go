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
 * IBM OpenAPI SDK Code Generator Version: 3.48.1-52130155-20220425-145431
 */

// Package cdtoolchainv2 : Operations and models for the CdToolchainV2 service
package cdtoolchainv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"time"

	common "github.com/IBM/continuous-delivery-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

// CdToolchainV2 : This swagger document describes the options and endpoints of the Toolchain API.<br><br> All calls
// require an <strong>Authorization</strong> HTTP header to be set with a bearer token, which can be generated using the
// <a href="https://cloud.ibm.com/apidocs/iam-identity-token-api">Identity and Access Management (IAM)
// API</a>.<br><br>Note that all resources must have a corresponding <b>resource_group_id</b> to use the API, resources
// within a Cloud Foundry organization cannot be accessed or modified using the API.
//
// API Version: 2.0.0
type CdToolchainV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.us-south.devops.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "cd_toolchain"

// CdToolchainV2Options : Service options
type CdToolchainV2Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewCdToolchainV2UsingExternalConfig : constructs an instance of CdToolchainV2 with passed in options and external configuration.
func NewCdToolchainV2UsingExternalConfig(options *CdToolchainV2Options) (cdToolchain *CdToolchainV2, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	cdToolchain, err = NewCdToolchainV2(options)
	if err != nil {
		return
	}

	err = cdToolchain.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = cdToolchain.Service.SetServiceURL(options.URL)
	}
	return
}

// NewCdToolchainV2 : constructs an instance of CdToolchainV2 with passed in options.
func NewCdToolchainV2(options *CdToolchainV2Options) (service *CdToolchainV2, err error) {
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

	service = &CdToolchainV2{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://api.us-south.devops.cloud.ibm.com", // The toolchain API endpoint in the us-south region
		"us-east": "https://api.us-east.devops.cloud.ibm.com", // The toolchain API endpoint in the us-east region
		"eu-de": "https://api.eu-de.devops.cloud.ibm.com", // The toolchain API endpoint in the eu-de region
		"eu-gb": "https://api.eu-gb.devops.cloud.ibm.com", // The toolchain API endpoint in the eu-gb region
		"jp-osa": "https://api.jp-osa.devops.cloud.ibm.com", // The toolchain API endpoint in the jp-osa region
		"jp-tok": "https://api.jp-tok.devops.cloud.ibm.com", // The toolchain API endpoint in the jp-tok region
		"au-syd": "https://api.au-syd.devops.cloud.ibm.com", // The toolchain API endpoint in the au-syd region
		"ca-tor": "https://api.ca-tor.devops.cloud.ibm.com", // The toolchain API endpoint in the ca-tor region
		"br-sao": "https://api.br-sao.devops.cloud.ibm.com", // The toolchain API endpoint in the br-sao region
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", fmt.Errorf("service URL for region '%s' not found", region)
}

// Clone makes a copy of "cdToolchain" suitable for processing requests.
func (cdToolchain *CdToolchainV2) Clone() *CdToolchainV2 {
	if core.IsNil(cdToolchain) {
		return nil
	}
	clone := *cdToolchain
	clone.Service = cdToolchain.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (cdToolchain *CdToolchainV2) SetServiceURL(url string) error {
	return cdToolchain.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (cdToolchain *CdToolchainV2) GetServiceURL() string {
	return cdToolchain.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (cdToolchain *CdToolchainV2) SetDefaultHeaders(headers http.Header) {
	cdToolchain.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (cdToolchain *CdToolchainV2) SetEnableGzipCompression(enableGzip bool) {
	cdToolchain.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (cdToolchain *CdToolchainV2) GetEnableGzipCompression() bool {
	return cdToolchain.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (cdToolchain *CdToolchainV2) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	cdToolchain.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (cdToolchain *CdToolchainV2) DisableRetries() {
	cdToolchain.Service.DisableRetries()
}

// ListToolchains : Returns a list of toolchains
// Returns a list of toolchains that the caller is authorized to access and that meet the provided query parameters.
func (cdToolchain *CdToolchainV2) ListToolchains(listToolchainsOptions *ListToolchainsOptions) (result *GetToolchainsResponse, response *core.DetailedResponse, err error) {
	return cdToolchain.ListToolchainsWithContext(context.Background(), listToolchainsOptions)
}

// ListToolchainsWithContext is an alternate form of the ListToolchains method which supports a Context parameter
func (cdToolchain *CdToolchainV2) ListToolchainsWithContext(ctx context.Context, listToolchainsOptions *ListToolchainsOptions) (result *GetToolchainsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listToolchainsOptions, "listToolchainsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listToolchainsOptions, "listToolchainsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range listToolchainsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "ListToolchains")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("resource_group_id", fmt.Sprint(*listToolchainsOptions.ResourceGroupID))
	if listToolchainsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listToolchainsOptions.Limit))
	}
	if listToolchainsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listToolchainsOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetToolchainsResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateToolchain : Create a toolchain
// Creates a new toolchain based off of provided parameters in the POST body.
func (cdToolchain *CdToolchainV2) CreateToolchain(createToolchainOptions *CreateToolchainOptions) (result *PostToolchainResponse, response *core.DetailedResponse, err error) {
	return cdToolchain.CreateToolchainWithContext(context.Background(), createToolchainOptions)
}

// CreateToolchainWithContext is an alternate form of the CreateToolchain method which supports a Context parameter
func (cdToolchain *CdToolchainV2) CreateToolchainWithContext(ctx context.Context, createToolchainOptions *CreateToolchainOptions) (result *PostToolchainResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createToolchainOptions, "createToolchainOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createToolchainOptions, "createToolchainOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createToolchainOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "CreateToolchain")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createToolchainOptions.Name != nil {
		body["name"] = createToolchainOptions.Name
	}
	if createToolchainOptions.ResourceGroupID != nil {
		body["resource_group_id"] = createToolchainOptions.ResourceGroupID
	}
	if createToolchainOptions.Description != nil {
		body["description"] = createToolchainOptions.Description
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
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPostToolchainResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetToolchainByID : Fetch a toolchain
// Returns data for a single toolchain identified by id.
func (cdToolchain *CdToolchainV2) GetToolchainByID(getToolchainByIDOptions *GetToolchainByIDOptions) (result *GetToolchainByIDResponse, response *core.DetailedResponse, err error) {
	return cdToolchain.GetToolchainByIDWithContext(context.Background(), getToolchainByIDOptions)
}

// GetToolchainByIDWithContext is an alternate form of the GetToolchainByID method which supports a Context parameter
func (cdToolchain *CdToolchainV2) GetToolchainByIDWithContext(ctx context.Context, getToolchainByIDOptions *GetToolchainByIDOptions) (result *GetToolchainByIDResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getToolchainByIDOptions, "getToolchainByIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getToolchainByIDOptions, "getToolchainByIDOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *getToolchainByIDOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getToolchainByIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "GetToolchainByID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetToolchainByIDResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteToolchain : Delete a toolchain
// Delete the toolchain with the specified ID.
func (cdToolchain *CdToolchainV2) DeleteToolchain(deleteToolchainOptions *DeleteToolchainOptions) (response *core.DetailedResponse, err error) {
	return cdToolchain.DeleteToolchainWithContext(context.Background(), deleteToolchainOptions)
}

// DeleteToolchainWithContext is an alternate form of the DeleteToolchain method which supports a Context parameter
func (cdToolchain *CdToolchainV2) DeleteToolchainWithContext(ctx context.Context, deleteToolchainOptions *DeleteToolchainOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteToolchainOptions, "deleteToolchainOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteToolchainOptions, "deleteToolchainOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *deleteToolchainOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteToolchainOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "DeleteToolchain")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cdToolchain.Service.Request(request, nil)

	return
}

// UpdateToolchain : Update a toolchain
// Update the toolchain with the specified ID.
func (cdToolchain *CdToolchainV2) UpdateToolchain(updateToolchainOptions *UpdateToolchainOptions) (response *core.DetailedResponse, err error) {
	return cdToolchain.UpdateToolchainWithContext(context.Background(), updateToolchainOptions)
}

// UpdateToolchainWithContext is an alternate form of the UpdateToolchain method which supports a Context parameter
func (cdToolchain *CdToolchainV2) UpdateToolchainWithContext(ctx context.Context, updateToolchainOptions *UpdateToolchainOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateToolchainOptions, "updateToolchainOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateToolchainOptions, "updateToolchainOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *updateToolchainOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateToolchainOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "UpdateToolchain")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateToolchainOptions.Name != nil {
		body["name"] = updateToolchainOptions.Name
	}
	if updateToolchainOptions.Description != nil {
		body["description"] = updateToolchainOptions.Description
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cdToolchain.Service.Request(request, nil)

	return
}

// ListTools : Returns a list of tools bound to toolchain
// Returns a list of tools bound to toolchain that the caller is authorized to access and that meet the provided query
// parameters.
func (cdToolchain *CdToolchainV2) ListTools(listToolsOptions *ListToolsOptions) (result *GetToolsResponse, response *core.DetailedResponse, err error) {
	return cdToolchain.ListToolsWithContext(context.Background(), listToolsOptions)
}

// ListToolsWithContext is an alternate form of the ListTools method which supports a Context parameter
func (cdToolchain *CdToolchainV2) ListToolsWithContext(ctx context.Context, listToolsOptions *ListToolsOptions) (result *GetToolsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listToolsOptions, "listToolsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listToolsOptions, "listToolsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *listToolsOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}/tools`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listToolsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "ListTools")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listToolsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listToolsOptions.Limit))
	}
	if listToolsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listToolsOptions.Offset))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetToolsResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// CreateTool : Create a tool
// Provisions a new tool based off of provided parameters in the POST body and binds it to the specified toolchain.
func (cdToolchain *CdToolchainV2) CreateTool(createToolOptions *CreateToolOptions) (result *PostToolResponse, response *core.DetailedResponse, err error) {
	return cdToolchain.CreateToolWithContext(context.Background(), createToolOptions)
}

// CreateToolWithContext is an alternate form of the CreateTool method which supports a Context parameter
func (cdToolchain *CdToolchainV2) CreateToolWithContext(ctx context.Context, createToolOptions *CreateToolOptions) (result *PostToolResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createToolOptions, "createToolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createToolOptions, "createToolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *createToolOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}/tools`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createToolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "CreateTool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createToolOptions.ToolTypeID != nil {
		body["tool_type_id"] = createToolOptions.ToolTypeID
	}
	if createToolOptions.Name != nil {
		body["name"] = createToolOptions.Name
	}
	if createToolOptions.Parameters != nil {
		body["parameters"] = createToolOptions.Parameters
	}
	if createToolOptions.ParametersReferences != nil {
		body["parameters_references"] = createToolOptions.ParametersReferences
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
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPostToolResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// GetToolByID : Fetch a tool
// Returns a tool that is bound to the provided toolchain.
func (cdToolchain *CdToolchainV2) GetToolByID(getToolByIDOptions *GetToolByIDOptions) (result *GetToolByIDResponse, response *core.DetailedResponse, err error) {
	return cdToolchain.GetToolByIDWithContext(context.Background(), getToolByIDOptions)
}

// GetToolByIDWithContext is an alternate form of the GetToolByID method which supports a Context parameter
func (cdToolchain *CdToolchainV2) GetToolByIDWithContext(ctx context.Context, getToolByIDOptions *GetToolByIDOptions) (result *GetToolByIDResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getToolByIDOptions, "getToolByIDOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getToolByIDOptions, "getToolByIDOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *getToolByIDOptions.ToolchainID,
		"tool_id": *getToolByIDOptions.ToolID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}/tools/{tool_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getToolByIDOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "GetToolByID")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGetToolByIDResponse)
		if err != nil {
			return
		}
		response.Result = result
	}

	return
}

// DeleteTool : Delete a tool
// Delete the tool with the specified ID.
func (cdToolchain *CdToolchainV2) DeleteTool(deleteToolOptions *DeleteToolOptions) (response *core.DetailedResponse, err error) {
	return cdToolchain.DeleteToolWithContext(context.Background(), deleteToolOptions)
}

// DeleteToolWithContext is an alternate form of the DeleteTool method which supports a Context parameter
func (cdToolchain *CdToolchainV2) DeleteToolWithContext(ctx context.Context, deleteToolOptions *DeleteToolOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteToolOptions, "deleteToolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteToolOptions, "deleteToolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *deleteToolOptions.ToolchainID,
		"tool_id": *deleteToolOptions.ToolID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}/tools/{tool_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteToolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "DeleteTool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cdToolchain.Service.Request(request, nil)

	return
}

// UpdateTool : Update a tool
// Update the tool with the specified ID.
func (cdToolchain *CdToolchainV2) UpdateTool(updateToolOptions *UpdateToolOptions) (response *core.DetailedResponse, err error) {
	return cdToolchain.UpdateToolWithContext(context.Background(), updateToolOptions)
}

// UpdateToolWithContext is an alternate form of the UpdateTool method which supports a Context parameter
func (cdToolchain *CdToolchainV2) UpdateToolWithContext(ctx context.Context, updateToolOptions *UpdateToolOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateToolOptions, "updateToolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateToolOptions, "updateToolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *updateToolOptions.ToolchainID,
		"tool_id": *updateToolOptions.ToolID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/v2/toolchains/{toolchain_id}/tools/{tool_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateToolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "UpdateTool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateToolOptions.Name != nil {
		body["name"] = updateToolOptions.Name
	}
	if updateToolOptions.ToolTypeID != nil {
		body["tool_type_id"] = updateToolOptions.ToolTypeID
	}
	if updateToolOptions.Parameters != nil {
		body["parameters"] = updateToolOptions.Parameters
	}
	if updateToolOptions.ParametersReferences != nil {
		body["parameters_references"] = updateToolOptions.ParametersReferences
	}
	_, err = builder.SetBodyContentJSON(body)
	if err != nil {
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	response, err = cdToolchain.Service.Request(request, nil)

	return
}

// CreateToolOptions : The CreateTool options.
type CreateToolOptions struct {
	// ID of the toolchain to bind tool to.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// The unique short name of the tool that should be provisioned.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// Name of tool.
	Name *string `json:"name,omitempty"`

	// Parameters to be used to create the tool.
	Parameters map[string]interface{} `json:"parameters,omitempty"`

	// Decoded values used on provision in the broker that reference fields in the parameters.
	ParametersReferences map[string]interface{} `json:"parameters_references,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateToolOptions : Instantiate CreateToolOptions
func (*CdToolchainV2) NewCreateToolOptions(toolchainID string, toolTypeID string) *CreateToolOptions {
	return &CreateToolOptions{
		ToolchainID: core.StringPtr(toolchainID),
		ToolTypeID: core.StringPtr(toolTypeID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *CreateToolOptions) SetToolchainID(toolchainID string) *CreateToolOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetToolTypeID : Allow user to set ToolTypeID
func (_options *CreateToolOptions) SetToolTypeID(toolTypeID string) *CreateToolOptions {
	_options.ToolTypeID = core.StringPtr(toolTypeID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateToolOptions) SetName(name string) *CreateToolOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetParameters : Allow user to set Parameters
func (_options *CreateToolOptions) SetParameters(parameters map[string]interface{}) *CreateToolOptions {
	_options.Parameters = parameters
	return _options
}

// SetParametersReferences : Allow user to set ParametersReferences
func (_options *CreateToolOptions) SetParametersReferences(parametersReferences map[string]interface{}) *CreateToolOptions {
	_options.ParametersReferences = parametersReferences
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateToolOptions) SetHeaders(param map[string]string) *CreateToolOptions {
	options.Headers = param
	return options
}

// CreateToolchainOptions : The CreateToolchain options.
type CreateToolchainOptions struct {
	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Resource group where toolchain will be created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Describes the toolchain.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateToolchainOptions : Instantiate CreateToolchainOptions
func (*CdToolchainV2) NewCreateToolchainOptions(name string, resourceGroupID string) *CreateToolchainOptions {
	return &CreateToolchainOptions{
		Name: core.StringPtr(name),
		ResourceGroupID: core.StringPtr(resourceGroupID),
	}
}

// SetName : Allow user to set Name
func (_options *CreateToolchainOptions) SetName(name string) *CreateToolchainOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *CreateToolchainOptions) SetResourceGroupID(resourceGroupID string) *CreateToolchainOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateToolchainOptions) SetDescription(description string) *CreateToolchainOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateToolchainOptions) SetHeaders(param map[string]string) *CreateToolchainOptions {
	options.Headers = param
	return options
}

// DeleteToolOptions : The DeleteTool options.
type DeleteToolOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// ID of the tool bound to the toolchain.
	ToolID *string `json:"tool_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteToolOptions : Instantiate DeleteToolOptions
func (*CdToolchainV2) NewDeleteToolOptions(toolchainID string, toolID string) *DeleteToolOptions {
	return &DeleteToolOptions{
		ToolchainID: core.StringPtr(toolchainID),
		ToolID: core.StringPtr(toolID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *DeleteToolOptions) SetToolchainID(toolchainID string) *DeleteToolOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetToolID : Allow user to set ToolID
func (_options *DeleteToolOptions) SetToolID(toolID string) *DeleteToolOptions {
	_options.ToolID = core.StringPtr(toolID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteToolOptions) SetHeaders(param map[string]string) *DeleteToolOptions {
	options.Headers = param
	return options
}

// DeleteToolchainOptions : The DeleteToolchain options.
type DeleteToolchainOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteToolchainOptions : Instantiate DeleteToolchainOptions
func (*CdToolchainV2) NewDeleteToolchainOptions(toolchainID string) *DeleteToolchainOptions {
	return &DeleteToolchainOptions{
		ToolchainID: core.StringPtr(toolchainID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *DeleteToolchainOptions) SetToolchainID(toolchainID string) *DeleteToolchainOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteToolchainOptions) SetHeaders(param map[string]string) *DeleteToolchainOptions {
	options.Headers = param
	return options
}

// GetToolByIDOptions : The GetToolByID options.
type GetToolByIDOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// ID of the tool bound to the toolchain.
	ToolID *string `json:"tool_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetToolByIDOptions : Instantiate GetToolByIDOptions
func (*CdToolchainV2) NewGetToolByIDOptions(toolchainID string, toolID string) *GetToolByIDOptions {
	return &GetToolByIDOptions{
		ToolchainID: core.StringPtr(toolchainID),
		ToolID: core.StringPtr(toolID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *GetToolByIDOptions) SetToolchainID(toolchainID string) *GetToolByIDOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetToolID : Allow user to set ToolID
func (_options *GetToolByIDOptions) SetToolID(toolID string) *GetToolByIDOptions {
	_options.ToolID = core.StringPtr(toolID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetToolByIDOptions) SetHeaders(param map[string]string) *GetToolByIDOptions {
	options.Headers = param
	return options
}

// GetToolchainByIDOptions : The GetToolchainByID options.
type GetToolchainByIDOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetToolchainByIDOptions : Instantiate GetToolchainByIDOptions
func (*CdToolchainV2) NewGetToolchainByIDOptions(toolchainID string) *GetToolchainByIDOptions {
	return &GetToolchainByIDOptions{
		ToolchainID: core.StringPtr(toolchainID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *GetToolchainByIDOptions) SetToolchainID(toolchainID string) *GetToolchainByIDOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetToolchainByIDOptions) SetHeaders(param map[string]string) *GetToolchainByIDOptions {
	options.Headers = param
	return options
}

// ListToolchainsOptions : The ListToolchains options.
type ListToolchainsOptions struct {
	// The resource group ID where the toolchains exist.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Limit the number of results. Valid value 0 < limit <= 200.
	Limit *int64 `json:"limit,omitempty"`

	// Offset the number of results from the beginning of the list. Valid value 0 <= offset < 200.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListToolchainsOptions : Instantiate ListToolchainsOptions
func (*CdToolchainV2) NewListToolchainsOptions(resourceGroupID string) *ListToolchainsOptions {
	return &ListToolchainsOptions{
		ResourceGroupID: core.StringPtr(resourceGroupID),
	}
}

// SetResourceGroupID : Allow user to set ResourceGroupID
func (_options *ListToolchainsOptions) SetResourceGroupID(resourceGroupID string) *ListToolchainsOptions {
	_options.ResourceGroupID = core.StringPtr(resourceGroupID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListToolchainsOptions) SetLimit(limit int64) *ListToolchainsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListToolchainsOptions) SetOffset(offset int64) *ListToolchainsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListToolchainsOptions) SetHeaders(param map[string]string) *ListToolchainsOptions {
	options.Headers = param
	return options
}

// ListToolsOptions : The ListTools options.
type ListToolsOptions struct {
	// ID of the toolchain that tools are bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// Limit the number of results. Valid value 0 < limit <= 200.
	Limit *int64 `json:"limit,omitempty"`

	// Offset the number of results from the beginning of the list. Valid value 0 <= offset < 200.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListToolsOptions : Instantiate ListToolsOptions
func (*CdToolchainV2) NewListToolsOptions(toolchainID string) *ListToolsOptions {
	return &ListToolsOptions{
		ToolchainID: core.StringPtr(toolchainID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *ListToolsOptions) SetToolchainID(toolchainID string) *ListToolsOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListToolsOptions) SetLimit(limit int64) *ListToolsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListToolsOptions) SetOffset(offset int64) *ListToolsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListToolsOptions) SetHeaders(param map[string]string) *ListToolsOptions {
	options.Headers = param
	return options
}

// ToolReferent : Information on URIs to access this resource through the UI or API.
type ToolReferent struct {
	// URI representing the this resource through the UI.
	UIHref *string `json:"ui_href,omitempty"`

	// URI representing the this resource through an API.
	APIHref *string `json:"api_href,omitempty"`
}

// UnmarshalToolReferent unmarshals an instance of ToolReferent from the specified map of raw messages.
func UnmarshalToolReferent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolReferent)
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_href", &obj.APIHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateToolOptions : The UpdateTool options.
type UpdateToolOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// ID of the tool bound to the toolchain.
	ToolID *string `json:"tool_id" validate:"required,ne="`

	// Name of tool.
	Name *string `json:"name,omitempty"`

	// The unique short name of the tool that should be provisioned or updated.
	ToolTypeID *string `json:"tool_type_id,omitempty"`

	// Parameters to be used to create the tool.
	Parameters map[string]interface{} `json:"parameters,omitempty"`

	// Decoded values used on provision in the broker that reference fields in the parameters.
	ParametersReferences map[string]interface{} `json:"parameters_references,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateToolOptions : Instantiate UpdateToolOptions
func (*CdToolchainV2) NewUpdateToolOptions(toolchainID string, toolID string) *UpdateToolOptions {
	return &UpdateToolOptions{
		ToolchainID: core.StringPtr(toolchainID),
		ToolID: core.StringPtr(toolID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *UpdateToolOptions) SetToolchainID(toolchainID string) *UpdateToolOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetToolID : Allow user to set ToolID
func (_options *UpdateToolOptions) SetToolID(toolID string) *UpdateToolOptions {
	_options.ToolID = core.StringPtr(toolID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateToolOptions) SetName(name string) *UpdateToolOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetToolTypeID : Allow user to set ToolTypeID
func (_options *UpdateToolOptions) SetToolTypeID(toolTypeID string) *UpdateToolOptions {
	_options.ToolTypeID = core.StringPtr(toolTypeID)
	return _options
}

// SetParameters : Allow user to set Parameters
func (_options *UpdateToolOptions) SetParameters(parameters map[string]interface{}) *UpdateToolOptions {
	_options.Parameters = parameters
	return _options
}

// SetParametersReferences : Allow user to set ParametersReferences
func (_options *UpdateToolOptions) SetParametersReferences(parametersReferences map[string]interface{}) *UpdateToolOptions {
	_options.ParametersReferences = parametersReferences
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateToolOptions) SetHeaders(param map[string]string) *UpdateToolOptions {
	options.Headers = param
	return options
}

// UpdateToolchainOptions : The UpdateToolchain options.
type UpdateToolchainOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// The name of the toolchain.
	Name *string `json:"name,omitempty"`

	// An optional description.
	Description *string `json:"description,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateToolchainOptions : Instantiate UpdateToolchainOptions
func (*CdToolchainV2) NewUpdateToolchainOptions(toolchainID string) *UpdateToolchainOptions {
	return &UpdateToolchainOptions{
		ToolchainID: core.StringPtr(toolchainID),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *UpdateToolchainOptions) SetToolchainID(toolchainID string) *UpdateToolchainOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateToolchainOptions) SetName(name string) *UpdateToolchainOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateToolchainOptions) SetDescription(description string) *UpdateToolchainOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateToolchainOptions) SetHeaders(param map[string]string) *UpdateToolchainOptions {
	options.Headers = param
	return options
}

// GetToolByIDResponse : Response structure for GET tool.
type GetToolByIDResponse struct {
	// Tool ID.
	ID *string `json:"id" validate:"required"`

	// Resource group where tool can be found.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Tool CRN.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the tool is bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the tool is bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *ToolReferent `json:"referent" validate:"required"`

	// Tool name.
	Name *string `json:"name,omitempty"`

	// Latest tool update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Parameters to be used to create the tool.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the GetToolByIDResponse.State property.
// Current configuration state of the tool.
const (
	GetToolByIDResponseStateConfiguredConst = "configured"
	GetToolByIDResponseStateConfiguringConst = "configuring"
	GetToolByIDResponseStateMisconfiguredConst = "misconfigured"
	GetToolByIDResponseStateUnconfiguredConst = "unconfigured"
)

// UnmarshalGetToolByIDResponse unmarshals an instance of GetToolByIDResponse from the specified map of raw messages.
func UnmarshalGetToolByIDResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolByIDResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalToolReferent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolchainByIDResponse : Response structure for GET toolchains.
type GetToolchainByIDResponse struct {
	// Toolchain ID.
	ID *string `json:"id" validate:"required"`

	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Toolchain description.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain can be found.
	AccountID *string `json:"account_id" validate:"required"`

	// Toolchain region.
	Location *string `json:"location" validate:"required"`

	// Resource group where toolchain can be found.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain CRN.
	CRN *string `json:"crn" validate:"required"`

	// URI that can be used to retrieve toolchain.
	Href *string `json:"href" validate:"required"`

	// Toolchain creation timestamp.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Latest toolchain update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`

	// Tags associated with the toolchain.
	Tags []string `json:"tags" validate:"required"`
}

// UnmarshalGetToolchainByIDResponse unmarshals an instance of GetToolchainByIDResponse from the specified map of raw messages.
func UnmarshalGetToolchainByIDResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolchainByIDResponse)
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
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolchainsResponse : Response structure for GET toolchains.
type GetToolchainsResponse struct {
	// Maximum number of toolchains returned from collection.
	Limit *int64 `json:"limit" validate:"required"`

	// Offset applied to toolchains collection.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of toolchains found in collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Information about retrieving first toolchain results from the collection.
	First *GetToolchainsResponseFirst `json:"first" validate:"required"`

	// Information about retrieving previous toolchain results from the collection.
	Previous *GetToolchainsResponsePrevious `json:"previous,omitempty"`

	// Information about retrieving next toolchain results from the collection.
	Next *GetToolchainsResponseNext `json:"next,omitempty"`

	// Information about retrieving last toolchain results from the collection.
	Last *GetToolchainsResponseLast `json:"last" validate:"required"`

	// Toolchain results returned from the collection.
	Toolchains []Toolchain `json:"toolchains" validate:"required"`
}

// UnmarshalGetToolchainsResponse unmarshals an instance of GetToolchainsResponse from the specified map of raw messages.
func UnmarshalGetToolchainsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolchainsResponse)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalGetToolchainsResponseFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalGetToolchainsResponsePrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalGetToolchainsResponseNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalGetToolchainsResponseLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "toolchains", &obj.Toolchains, UnmarshalToolchain)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *GetToolchainsResponse) GetNextOffset() (*int64, error) {
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

// GetToolchainsResponseFirst : Information about retrieving first toolchain results from the collection.
type GetToolchainsResponseFirst struct {
	// URI that can be used to get first results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolchainsResponseFirst unmarshals an instance of GetToolchainsResponseFirst from the specified map of raw messages.
func UnmarshalGetToolchainsResponseFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolchainsResponseFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolchainsResponseLast : Information about retrieving last toolchain results from the collection.
type GetToolchainsResponseLast struct {
	// URI that can be used to get last results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolchainsResponseLast unmarshals an instance of GetToolchainsResponseLast from the specified map of raw messages.
func UnmarshalGetToolchainsResponseLast(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolchainsResponseLast)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolchainsResponseNext : Information about retrieving next toolchain results from the collection.
type GetToolchainsResponseNext struct {
	// URI that can be used to get next results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolchainsResponseNext unmarshals an instance of GetToolchainsResponseNext from the specified map of raw messages.
func UnmarshalGetToolchainsResponseNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolchainsResponseNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolchainsResponsePrevious : Information about retrieving previous toolchain results from the collection.
type GetToolchainsResponsePrevious struct {
	// URI that can be used to get previous results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolchainsResponsePrevious unmarshals an instance of GetToolchainsResponsePrevious from the specified map of raw messages.
func UnmarshalGetToolchainsResponsePrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolchainsResponsePrevious)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolsResponse : Response structure for GET tools.
type GetToolsResponse struct {
	// Maximum number of tools returned from collection.
	Limit *int64 `json:"limit" validate:"required"`

	// Offset applied to tools collection.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of tools found in collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Information about retrieving first tool results from the collection.
	First *GetToolsResponseFirst `json:"first" validate:"required"`

	// Information about retrieving previous tool results from the collection.
	Previous *GetToolsResponsePrevious `json:"previous,omitempty"`

	// Information about retrieving next tool results from the collection.
	Next *GetToolsResponseNext `json:"next,omitempty"`

	// Information about retrieving last tool results from the collection.
	Last *GetToolsResponseLast `json:"last" validate:"required"`

	// Tool results returned from the collection.
	Tools []Tool `json:"tools" validate:"required"`
}

// UnmarshalGetToolsResponse unmarshals an instance of GetToolsResponse from the specified map of raw messages.
func UnmarshalGetToolsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolsResponse)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalGetToolsResponseFirst)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalGetToolsResponsePrevious)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalGetToolsResponseNext)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalGetToolsResponseLast)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "tools", &obj.Tools, UnmarshalTool)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *GetToolsResponse) GetNextOffset() (*int64, error) {
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

// GetToolsResponseFirst : Information about retrieving first tool results from the collection.
type GetToolsResponseFirst struct {
	// URI that can be used to get first results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolsResponseFirst unmarshals an instance of GetToolsResponseFirst from the specified map of raw messages.
func UnmarshalGetToolsResponseFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolsResponseFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolsResponseLast : Information about retrieving last tool results from the collection.
type GetToolsResponseLast struct {
	// URI that can be used to get last results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolsResponseLast unmarshals an instance of GetToolsResponseLast from the specified map of raw messages.
func UnmarshalGetToolsResponseLast(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolsResponseLast)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolsResponseNext : Information about retrieving next tool results from the collection.
type GetToolsResponseNext struct {
	// URI that can be used to get next results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolsResponseNext unmarshals an instance of GetToolsResponseNext from the specified map of raw messages.
func UnmarshalGetToolsResponseNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolsResponseNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetToolsResponsePrevious : Information about retrieving previous tool results from the collection.
type GetToolsResponsePrevious struct {
	// URI that can be used to get previous results from the collection.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGetToolsResponsePrevious unmarshals an instance of GetToolsResponsePrevious from the specified map of raw messages.
func UnmarshalGetToolsResponsePrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GetToolsResponsePrevious)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostToolResponse : Response structure for POST tool.
type PostToolResponse struct {
	// ID of created tool.
	ID *string `json:"id" validate:"required"`

	// Resource group where tool was created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// CRN of created tool.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the created tool was bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the created tool was bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the created tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *Referent `json:"referent" validate:"required"`

	// Name of tool.
	Name *string `json:"name,omitempty"`

	// Parameters to be used to create the tool.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the PostToolResponse.State property.
// Current configuration state of the tool.
const (
	PostToolResponseStateConfiguredConst = "configured"
	PostToolResponseStateConfiguringConst = "configuring"
	PostToolResponseStateMisconfiguredConst = "misconfigured"
	PostToolResponseStateUnconfiguredConst = "unconfigured"
)

// UnmarshalPostToolResponse unmarshals an instance of PostToolResponse from the specified map of raw messages.
func UnmarshalPostToolResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PostToolResponse)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalReferent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PostToolchainResponse : POST toolchains response body.
type PostToolchainResponse struct {
	// ID of created toolchain.
	ID *string `json:"id" validate:"required"`

	// Name of created toolchain.
	Name *string `json:"name" validate:"required"`

	// Description of created toolchain.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain was created.
	AccountID *string `json:"account_id" validate:"required"`

	// Region where toolchain is created.
	Location *string `json:"location" validate:"required"`

	// Resource group where toolchain is created.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// CRN of created toolchain.
	CRN *string `json:"crn" validate:"required"`

	// URI representing the created toolchain.
	Href *string `json:"href" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`

	// Tags associated with the created toolchain.
	Tags []string `json:"tags" validate:"required"`
}

// UnmarshalPostToolchainResponse unmarshals an instance of PostToolchainResponse from the specified map of raw messages.
func UnmarshalPostToolchainResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PostToolchainResponse)
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
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Referent : Information on URIs to access this resource through the UI or API.
type Referent struct {
	// URI representing the this resource through the UI.
	UIHref *string `json:"ui_href,omitempty"`

	// URI representing the this resource through an API.
	APIHref *string `json:"api_href,omitempty"`
}

// UnmarshalReferent unmarshals an instance of Referent from the specified map of raw messages.
func UnmarshalReferent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Referent)
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "api_href", &obj.APIHref)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Tool : Model describing tool resource.
type Tool struct {
	// Tool ID.
	ID *string `json:"id" validate:"required"`

	// Resource group where tool can be found.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Tool CRN.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the tool is bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the tool is bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *ToolReferent `json:"referent" validate:"required"`

	// Tool name.
	Name *string `json:"name,omitempty"`

	// Latest tool update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Parameters to be used to create the tool.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the Tool.State property.
// Current configuration state of the tool.
const (
	ToolStateConfiguredConst = "configured"
	ToolStateConfiguringConst = "configuring"
	ToolStateMisconfiguredConst = "misconfigured"
	ToolStateUnconfiguredConst = "unconfigured"
)

// UnmarshalTool unmarshals an instance of Tool from the specified map of raw messages.
func UnmarshalTool(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Tool)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalToolReferent)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Toolchain : Model describing toolchain resource.
type Toolchain struct {
	// Toolchain ID.
	ID *string `json:"id" validate:"required"`

	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Toolchain description.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain can be found.
	AccountID *string `json:"account_id" validate:"required"`

	// Toolchain region.
	Location *string `json:"location" validate:"required"`

	// Resource group where toolchain can be found.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain CRN.
	CRN *string `json:"crn" validate:"required"`

	// URI that can be used to retrieve toolchain.
	Href *string `json:"href" validate:"required"`

	// Toolchain creation timestamp.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Latest toolchain update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`

	// Tags associated with the toolchain.
	Tags []string `json:"tags" validate:"required"`
}

// UnmarshalToolchain unmarshals an instance of Toolchain from the specified map of raw messages.
func UnmarshalToolchain(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Toolchain)
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
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
