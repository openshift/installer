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
 * IBM OpenAPI SDK Code Generator Version: 3.96.0-d6dec9d7-20241008-212902
 */

// Package cdtoolchainv2 : Operations and models for the CdToolchainV2 service
package cdtoolchainv2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM/continuous-delivery-go-sdk/v2/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

type CdToolchainV2 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.us-south.devops.cloud.ibm.com/toolchain/v2"

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
			err = core.SDKErrorf(err, "", "env-auth-error", common.GetComponentInfo())
			return
		}
	}

	cdToolchain, err = NewCdToolchainV2(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = cdToolchain.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = cdToolchain.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
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

	service = &CdToolchainV2{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	var endpoints = map[string]string{
		"us-south": "https://api.us-south.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the us-south region
		"us-east": "https://api.us-east.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the us-east region
		"eu-de": "https://api.eu-de.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the eu-de region
		"eu-gb": "https://api.eu-gb.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the eu-gb region
		"jp-osa": "https://api.jp-osa.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the jp-osa region
		"jp-tok": "https://api.jp-tok.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the jp-tok region
		"au-syd": "https://api.au-syd.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the au-syd region
		"ca-tor": "https://api.ca-tor.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the ca-tor region
		"br-sao": "https://api.br-sao.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the br-sao region
		"eu-es": "https://api.eu-es.devops.cloud.ibm.com/toolchain/v2", // The toolchain API endpoint in the eu-es region
	}

	if url, ok := endpoints[region]; ok {
		return url, nil
	}
	return "", core.SDKErrorf(nil, fmt.Sprintf("service URL for region '%s' not found", region), "invalid-region", common.GetComponentInfo())
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
	err := cdToolchain.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
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

// ListToolchains : Get a list of toolchains
// Returns a list of toolchains that the caller is authorized to access and that meets the provided query parameters.
func (cdToolchain *CdToolchainV2) ListToolchains(listToolchainsOptions *ListToolchainsOptions) (result *ToolchainCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.ListToolchainsWithContext(context.Background(), listToolchainsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListToolchainsWithContext is an alternate form of the ListToolchains method which supports a Context parameter
func (cdToolchain *CdToolchainV2) ListToolchainsWithContext(ctx context.Context, listToolchainsOptions *ListToolchainsOptions) (result *ToolchainCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listToolchainsOptions, "listToolchainsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listToolchainsOptions, "listToolchainsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
	if listToolchainsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listToolchainsOptions.Start))
	}
	if listToolchainsOptions.Name != nil {
		builder.AddQuery("name", fmt.Sprint(*listToolchainsOptions.Name))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_toolchains", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateToolchain : Create a toolchain
// Creates a new toolchain based off the provided parameters in the body.
func (cdToolchain *CdToolchainV2) CreateToolchain(createToolchainOptions *CreateToolchainOptions) (result *ToolchainPost, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.CreateToolchainWithContext(context.Background(), createToolchainOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateToolchainWithContext is an alternate form of the CreateToolchain method which supports a Context parameter
func (cdToolchain *CdToolchainV2) CreateToolchainWithContext(ctx context.Context, createToolchainOptions *CreateToolchainOptions) (result *ToolchainPost, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createToolchainOptions, "createToolchainOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createToolchainOptions, "createToolchainOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "set-json-body-error", common.GetComponentInfo())
		return
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_toolchain", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainPost)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetToolchainByID : Get a toolchain
// Returns data for a single toolchain identified by its ID.
func (cdToolchain *CdToolchainV2) GetToolchainByID(getToolchainByIDOptions *GetToolchainByIDOptions) (result *Toolchain, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.GetToolchainByIDWithContext(context.Background(), getToolchainByIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetToolchainByIDWithContext is an alternate form of the GetToolchainByID method which supports a Context parameter
func (cdToolchain *CdToolchainV2) GetToolchainByIDWithContext(ctx context.Context, getToolchainByIDOptions *GetToolchainByIDOptions) (result *Toolchain, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getToolchainByIDOptions, "getToolchainByIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getToolchainByIDOptions, "getToolchainByIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *getToolchainByIDOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_toolchain_by_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchain)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteToolchain : Delete a toolchain
// Delete the toolchain with the specified ID.
func (cdToolchain *CdToolchainV2) DeleteToolchain(deleteToolchainOptions *DeleteToolchainOptions) (response *core.DetailedResponse, err error) {
	response, err = cdToolchain.DeleteToolchainWithContext(context.Background(), deleteToolchainOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteToolchainWithContext is an alternate form of the DeleteToolchain method which supports a Context parameter
func (cdToolchain *CdToolchainV2) DeleteToolchainWithContext(ctx context.Context, deleteToolchainOptions *DeleteToolchainOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteToolchainOptions, "deleteToolchainOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteToolchainOptions, "deleteToolchainOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *deleteToolchainOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdToolchain.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_toolchain", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateToolchain : Update a toolchain
// Update the toolchain with the specified ID.
func (cdToolchain *CdToolchainV2) UpdateToolchain(updateToolchainOptions *UpdateToolchainOptions) (result *ToolchainPatch, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.UpdateToolchainWithContext(context.Background(), updateToolchainOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateToolchainWithContext is an alternate form of the UpdateToolchain method which supports a Context parameter
func (cdToolchain *CdToolchainV2) UpdateToolchainWithContext(ctx context.Context, updateToolchainOptions *UpdateToolchainOptions) (result *ToolchainPatch, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateToolchainOptions, "updateToolchainOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateToolchainOptions, "updateToolchainOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *updateToolchainOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateToolchainOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "UpdateToolchain")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateToolchainOptions.ToolchainPrototypePatch)
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
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_toolchain", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainPatch)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateToolchainEvent : Create a toolchain event
// Creates and sends a custom event to each Event Notifications instance configured as a tool into the toolchain. This
// operation will fail if no Event Notifications instances are configured into the toolchain.
func (cdToolchain *CdToolchainV2) CreateToolchainEvent(createToolchainEventOptions *CreateToolchainEventOptions) (result *ToolchainEventPost, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.CreateToolchainEventWithContext(context.Background(), createToolchainEventOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateToolchainEventWithContext is an alternate form of the CreateToolchainEvent method which supports a Context parameter
func (cdToolchain *CdToolchainV2) CreateToolchainEventWithContext(ctx context.Context, createToolchainEventOptions *CreateToolchainEventOptions) (result *ToolchainEventPost, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createToolchainEventOptions, "createToolchainEventOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createToolchainEventOptions, "createToolchainEventOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *createToolchainEventOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}/events`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createToolchainEventOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "CreateToolchainEvent")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createToolchainEventOptions.Title != nil {
		body["title"] = createToolchainEventOptions.Title
	}
	if createToolchainEventOptions.Description != nil {
		body["description"] = createToolchainEventOptions.Description
	}
	if createToolchainEventOptions.ContentType != nil {
		body["content_type"] = createToolchainEventOptions.ContentType
	}
	if createToolchainEventOptions.Data != nil {
		body["data"] = createToolchainEventOptions.Data
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
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_toolchain_event", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainEventPost)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListTools : Get a list of tools bound to a toolchain
// Returns a list of tools bound to a toolchain that the caller is authorized to access and that meet the provided query
// parameters.
func (cdToolchain *CdToolchainV2) ListTools(listToolsOptions *ListToolsOptions) (result *ToolchainToolCollection, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.ListToolsWithContext(context.Background(), listToolsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListToolsWithContext is an alternate form of the ListTools method which supports a Context parameter
func (cdToolchain *CdToolchainV2) ListToolsWithContext(ctx context.Context, listToolsOptions *ListToolsOptions) (result *ToolchainToolCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listToolsOptions, "listToolsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listToolsOptions, "listToolsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *listToolsOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}/tools`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
	if listToolsOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listToolsOptions.Start))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_tools", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainToolCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTool : Create a tool
// Provisions a new tool based off the provided parameters in the body and binds it to the specified toolchain.
func (cdToolchain *CdToolchainV2) CreateTool(createToolOptions *CreateToolOptions) (result *ToolchainToolPost, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.CreateToolWithContext(context.Background(), createToolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateToolWithContext is an alternate form of the CreateTool method which supports a Context parameter
func (cdToolchain *CdToolchainV2) CreateToolWithContext(ctx context.Context, createToolOptions *CreateToolOptions) (result *ToolchainToolPost, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createToolOptions, "createToolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createToolOptions, "createToolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *createToolOptions.ToolchainID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}/tools`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_tool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainToolPost)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetToolByID : Get a tool
// Returns a tool that is bound to the provided toolchain.
func (cdToolchain *CdToolchainV2) GetToolByID(getToolByIDOptions *GetToolByIDOptions) (result *ToolchainTool, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.GetToolByIDWithContext(context.Background(), getToolByIDOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetToolByIDWithContext is an alternate form of the GetToolByID method which supports a Context parameter
func (cdToolchain *CdToolchainV2) GetToolByIDWithContext(ctx context.Context, getToolByIDOptions *GetToolByIDOptions) (result *ToolchainTool, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getToolByIDOptions, "getToolByIDOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getToolByIDOptions, "getToolByIDOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *getToolByIDOptions.ToolchainID,
		"tool_id": *getToolByIDOptions.ToolID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}/tools/{tool_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_tool_by_id", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainTool)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTool : Delete a tool
// Delete the tool with the specified ID.
func (cdToolchain *CdToolchainV2) DeleteTool(deleteToolOptions *DeleteToolOptions) (response *core.DetailedResponse, err error) {
	response, err = cdToolchain.DeleteToolWithContext(context.Background(), deleteToolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteToolWithContext is an alternate form of the DeleteTool method which supports a Context parameter
func (cdToolchain *CdToolchainV2) DeleteToolWithContext(ctx context.Context, deleteToolOptions *DeleteToolOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteToolOptions, "deleteToolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteToolOptions, "deleteToolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *deleteToolOptions.ToolchainID,
		"tool_id": *deleteToolOptions.ToolID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}/tools/{tool_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
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
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = cdToolchain.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_tool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// UpdateTool : Update a tool
// Update the tool with the specified ID.
func (cdToolchain *CdToolchainV2) UpdateTool(updateToolOptions *UpdateToolOptions) (result *ToolchainToolPatch, response *core.DetailedResponse, err error) {
	result, response, err = cdToolchain.UpdateToolWithContext(context.Background(), updateToolOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateToolWithContext is an alternate form of the UpdateTool method which supports a Context parameter
func (cdToolchain *CdToolchainV2) UpdateToolWithContext(ctx context.Context, updateToolOptions *UpdateToolOptions) (result *ToolchainToolPatch, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateToolOptions, "updateToolOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateToolOptions, "updateToolOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"toolchain_id": *updateToolOptions.ToolchainID,
		"tool_id": *updateToolOptions.ToolID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = cdToolchain.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(cdToolchain.Service.Options.URL, `/toolchains/{toolchain_id}/tools/{tool_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateToolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("cd_toolchain", "V2", "UpdateTool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/merge-patch+json")

	_, err = builder.SetBodyContentJSON(updateToolOptions.ToolchainToolPrototypePatch)
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
	response, err = cdToolchain.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_tool", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalToolchainToolPatch)
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

// CreateToolOptions : The CreateTool options.
type CreateToolOptions struct {
	// ID of the toolchain to bind the tool to.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// The unique short name of the tool that should be provisioned. A table of `tool_type_id` values corresponding to each
	// tool integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// Name of the tool.
	Name *string `json:"name,omitempty"`

	// Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool
	// integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	Parameters map[string]interface{} `json:"parameters,omitempty"`

	// Allows users to set headers on API requests.
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

// SetHeaders : Allow user to set Headers
func (options *CreateToolOptions) SetHeaders(param map[string]string) *CreateToolOptions {
	options.Headers = param
	return options
}

// CreateToolchainEventOptions : The CreateToolchainEvent options.
type CreateToolchainEventOptions struct {
	// ID of the toolchain to send events from.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// Event title.
	Title *string `json:"title" validate:"required"`

	// Describes the event.
	Description *string `json:"description" validate:"required"`

	// The content type of the attached data. Supported values are `text/plain`, `application/json`, and `none`.
	ContentType *string `json:"content_type" validate:"required"`

	// Additional data to be added with the event. The format must correspond to the value of `content_type`.
	Data *ToolchainEventPrototypeData `json:"data,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateToolchainEventOptions.ContentType property.
// The content type of the attached data. Supported values are `text/plain`, `application/json`, and `none`.
const (
	CreateToolchainEventOptionsContentTypeApplicationJSONConst = "application/json"
	CreateToolchainEventOptionsContentTypeNoneConst = "none"
	CreateToolchainEventOptionsContentTypeTextPlainConst = "text/plain"
)

// NewCreateToolchainEventOptions : Instantiate CreateToolchainEventOptions
func (*CdToolchainV2) NewCreateToolchainEventOptions(toolchainID string, title string, description string, contentType string) *CreateToolchainEventOptions {
	return &CreateToolchainEventOptions{
		ToolchainID: core.StringPtr(toolchainID),
		Title: core.StringPtr(title),
		Description: core.StringPtr(description),
		ContentType: core.StringPtr(contentType),
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *CreateToolchainEventOptions) SetToolchainID(toolchainID string) *CreateToolchainEventOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetTitle : Allow user to set Title
func (_options *CreateToolchainEventOptions) SetTitle(title string) *CreateToolchainEventOptions {
	_options.Title = core.StringPtr(title)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateToolchainEventOptions) SetDescription(description string) *CreateToolchainEventOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetContentType : Allow user to set ContentType
func (_options *CreateToolchainEventOptions) SetContentType(contentType string) *CreateToolchainEventOptions {
	_options.ContentType = core.StringPtr(contentType)
	return _options
}

// SetData : Allow user to set Data
func (_options *CreateToolchainEventOptions) SetData(data *ToolchainEventPrototypeData) *CreateToolchainEventOptions {
	_options.Data = data
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateToolchainEventOptions) SetHeaders(param map[string]string) *CreateToolchainEventOptions {
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

	// Allows users to set headers on API requests.
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

	// Allows users to set headers on API requests.
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

	// Allows users to set headers on API requests.
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

	// Allows users to set headers on API requests.
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

	// Allows users to set headers on API requests.
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

	// Limit the number of results.
	Limit *int64 `json:"limit,omitempty"`

	// Pagination token.
	Start *string `json:"start,omitempty"`

	// Exact name of toolchain to look up. This parameter is case sensitive.
	Name *string `json:"name,omitempty"`

	// Allows users to set headers on API requests.
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

// SetStart : Allow user to set Start
func (_options *ListToolchainsOptions) SetStart(start string) *ListToolchainsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetName : Allow user to set Name
func (_options *ListToolchainsOptions) SetName(name string) *ListToolchainsOptions {
	_options.Name = core.StringPtr(name)
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

	// Limit the number of results.
	Limit *int64 `json:"limit,omitempty"`

	// Pagination token.
	Start *string `json:"start,omitempty"`

	// Allows users to set headers on API requests.
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

// SetStart : Allow user to set Start
func (_options *ListToolsOptions) SetStart(start string) *ListToolsOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListToolsOptions) SetHeaders(param map[string]string) *ListToolsOptions {
	options.Headers = param
	return options
}

// ToolModel : Model describing tool resource.
type ToolModel struct {
	// Tool ID.
	ID *string `json:"id" validate:"required"`

	// Resource group where the tool is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Tool CRN.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool. A table of `tool_type_id` values corresponding to each tool integration can
	// be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the tool is bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the tool is bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *ToolModelReferent `json:"referent" validate:"required"`

	// Name of the tool.
	Name *string `json:"name,omitempty"`

	// Latest tool update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool
	// integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the ToolModel.State property.
// Current configuration state of the tool.
const (
	ToolModelStateConfiguredConst = "configured"
	ToolModelStateConfiguringConst = "configuring"
	ToolModelStateMisconfiguredConst = "misconfigured"
	ToolModelStateUnconfiguredConst = "unconfigured"
)

// UnmarshalToolModel unmarshals an instance of ToolModel from the specified map of raw messages.
func UnmarshalToolModel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolModel)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool_type_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalToolModelReferent)
	if err != nil {
		err = core.SDKErrorf(err, "", "referent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		err = core.SDKErrorf(err, "", "parameters-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolModelReferent : Information on URIs to access this resource through the UI or API.
type ToolModelReferent struct {
	// URI representing this resource through the UI.
	UIHref *string `json:"ui_href,omitempty"`

	// URI representing this resource through an API.
	APIHref *string `json:"api_href,omitempty"`
}

// UnmarshalToolModelReferent unmarshals an instance of ToolModelReferent from the specified map of raw messages.
func UnmarshalToolModelReferent(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolModelReferent)
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui_href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_href", &obj.APIHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Toolchain : Response structure for GET toolchains.
type Toolchain struct {
	// Toolchain ID.
	ID *string `json:"id" validate:"required"`

	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Describes the toolchain.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain can be found.
	AccountID *string `json:"account_id" validate:"required"`

	// Toolchain region.
	Location *string `json:"location" validate:"required"`

	// Resource group where the toolchain is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain CRN.
	CRN *string `json:"crn" validate:"required"`

	// URI that can be used to retrieve toolchain.
	Href *string `json:"href" validate:"required"`

	// URL of a user-facing user interface for this toolchain.
	UIHref *string `json:"ui_href" validate:"required"`

	// Toolchain creation timestamp.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Latest toolchain update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`
}

// UnmarshalToolchain unmarshals an instance of Toolchain from the specified map of raw messages.
func UnmarshalToolchain(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Toolchain)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui_href-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainCollection : Response structure for GET toolchains.
type ToolchainCollection struct {
	// Total number of toolchains found in collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Maximum number of toolchains returned from collection.
	Limit *int64 `json:"limit" validate:"required"`

	// Information about retrieving first toolchain results from the collection.
	First *ToolchainCollectionFirst `json:"first" validate:"required"`

	// Information about retrieving previous toolchain results from the collection.
	Previous *ToolchainCollectionPrevious `json:"previous,omitempty"`

	// Information about retrieving next toolchain results from the collection.
	Next *ToolchainCollectionNext `json:"next,omitempty"`

	// Information about retrieving last toolchain results from the collection.
	Last *ToolchainCollectionLast `json:"last" validate:"required"`

	// Toolchain results returned from the collection.
	Toolchains []ToolchainModel `json:"toolchains,omitempty"`
}

// UnmarshalToolchainCollection unmarshals an instance of ToolchainCollection from the specified map of raw messages.
func UnmarshalToolchainCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalToolchainCollectionFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalToolchainCollectionPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalToolchainCollectionNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalToolchainCollectionLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "toolchains", &obj.Toolchains, UnmarshalToolchainModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchains-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ToolchainCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ToolchainCollectionFirst : Information about retrieving first toolchain results from the collection.
type ToolchainCollectionFirst struct {
	// URI that can be used to get first results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainCollectionFirst unmarshals an instance of ToolchainCollectionFirst from the specified map of raw messages.
func UnmarshalToolchainCollectionFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainCollectionFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainCollectionLast : Information about retrieving last toolchain results from the collection.
type ToolchainCollectionLast struct {
	// Cursor that can be set as the 'start' query parameter to get the last set of toolchain collections.
	Start *string `json:"start,omitempty"`

	// URI that can be used to get last results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainCollectionLast unmarshals an instance of ToolchainCollectionLast from the specified map of raw messages.
func UnmarshalToolchainCollectionLast(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainCollectionLast)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainCollectionNext : Information about retrieving next toolchain results from the collection.
type ToolchainCollectionNext struct {
	// Cursor that can be set as the 'start' query parameter to get the next set of toolchain collections.
	Start *string `json:"start,omitempty"`

	// URI that can be used to get next results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainCollectionNext unmarshals an instance of ToolchainCollectionNext from the specified map of raw messages.
func UnmarshalToolchainCollectionNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainCollectionNext)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainCollectionPrevious : Information about retrieving previous toolchain results from the collection.
type ToolchainCollectionPrevious struct {
	// Cursor that can be set as the 'start' query parameter to get the previous set of toolchain collections.
	Start *string `json:"start,omitempty"`

	// URI that can be used to get previous results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainCollectionPrevious unmarshals an instance of ToolchainCollectionPrevious from the specified map of raw messages.
func UnmarshalToolchainCollectionPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainCollectionPrevious)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainEventPost : Response structure for POST toolchain event.
type ToolchainEventPost struct {
	// Event ID.
	ID *string `json:"id" validate:"required"`
}

// UnmarshalToolchainEventPost unmarshals an instance of ToolchainEventPost from the specified map of raw messages.
func UnmarshalToolchainEventPost(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainEventPost)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainEventPrototypeData : Additional data to be added with the event. The format must correspond to the value of `content_type`.
type ToolchainEventPrototypeData struct {
	// Contains JSON data to be added with the event. `content_type` must be set to `application/json`.
	ApplicationJSON *ToolchainEventPrototypeDataApplicationJSON `json:"application_json,omitempty"`

	// Contains text data to be added with the event. `content_type` must be set to `text/plain`.
	TextPlain *ToolchainEventPrototypeDataTextPlain `json:"text_plain,omitempty"`
}

// UnmarshalToolchainEventPrototypeData unmarshals an instance of ToolchainEventPrototypeData from the specified map of raw messages.
func UnmarshalToolchainEventPrototypeData(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainEventPrototypeData)
	err = core.UnmarshalModel(m, "application_json", &obj.ApplicationJSON, UnmarshalToolchainEventPrototypeDataApplicationJSON)
	if err != nil {
		err = core.SDKErrorf(err, "", "application_json-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "text_plain", &obj.TextPlain, UnmarshalToolchainEventPrototypeDataTextPlain)
	if err != nil {
		err = core.SDKErrorf(err, "", "text_plain-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainEventPrototypeDataApplicationJSON : Contains JSON data to be added with the event. `content_type` must be set to `application/json`.
type ToolchainEventPrototypeDataApplicationJSON struct {
	// JSON-formatted key-value pairs representing any additional information to be included with the event. The payload is
	// constrained to a maximum depth of 5, and keys that must satisfy the pattern ^[a-zA-Z0-9-_]+$.
	Content map[string]interface{} `json:"content" validate:"required"`
}

// NewToolchainEventPrototypeDataApplicationJSON : Instantiate ToolchainEventPrototypeDataApplicationJSON (Generic Model Constructor)
func (*CdToolchainV2) NewToolchainEventPrototypeDataApplicationJSON(content map[string]interface{}) (_model *ToolchainEventPrototypeDataApplicationJSON, err error) {
	_model = &ToolchainEventPrototypeDataApplicationJSON{
		Content: content,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalToolchainEventPrototypeDataApplicationJSON unmarshals an instance of ToolchainEventPrototypeDataApplicationJSON from the specified map of raw messages.
func UnmarshalToolchainEventPrototypeDataApplicationJSON(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainEventPrototypeDataApplicationJSON)
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainEventPrototypeDataTextPlain : Contains text data to be added with the event. `content_type` must be set to `text/plain`.
type ToolchainEventPrototypeDataTextPlain struct {
	// The text data to send in the event.
	Content *string `json:"content" validate:"required"`
}

// NewToolchainEventPrototypeDataTextPlain : Instantiate ToolchainEventPrototypeDataTextPlain (Generic Model Constructor)
func (*CdToolchainV2) NewToolchainEventPrototypeDataTextPlain(content string) (_model *ToolchainEventPrototypeDataTextPlain, err error) {
	_model = &ToolchainEventPrototypeDataTextPlain{
		Content: core.StringPtr(content),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalToolchainEventPrototypeDataTextPlain unmarshals an instance of ToolchainEventPrototypeDataTextPlain from the specified map of raw messages.
func UnmarshalToolchainEventPrototypeDataTextPlain(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainEventPrototypeDataTextPlain)
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		err = core.SDKErrorf(err, "", "content-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainModel : Model describing toolchain resource.
type ToolchainModel struct {
	// Toolchain ID.
	ID *string `json:"id" validate:"required"`

	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Describes the toolchain.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain can be found.
	AccountID *string `json:"account_id" validate:"required"`

	// Toolchain region.
	Location *string `json:"location" validate:"required"`

	// Resource group where the toolchain is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain CRN.
	CRN *string `json:"crn" validate:"required"`

	// URI that can be used to retrieve toolchain.
	Href *string `json:"href" validate:"required"`

	// URL of a user-facing user interface for this toolchain.
	UIHref *string `json:"ui_href" validate:"required"`

	// Toolchain creation timestamp.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Latest toolchain update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`
}

// UnmarshalToolchainModel unmarshals an instance of ToolchainModel from the specified map of raw messages.
func UnmarshalToolchainModel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainModel)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui_href-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainPatch : Response structure for PATCH toolchain.
type ToolchainPatch struct {
	// Toolchain ID.
	ID *string `json:"id" validate:"required"`

	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Describes the toolchain.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain can be found.
	AccountID *string `json:"account_id" validate:"required"`

	// Toolchain region.
	Location *string `json:"location" validate:"required"`

	// Resource group where the toolchain is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain CRN.
	CRN *string `json:"crn" validate:"required"`

	// URI that can be used to retrieve toolchain.
	Href *string `json:"href" validate:"required"`

	// URL of a user-facing user interface for this toolchain.
	UIHref *string `json:"ui_href" validate:"required"`

	// Toolchain creation timestamp.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Latest toolchain update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`
}

// UnmarshalToolchainPatch unmarshals an instance of ToolchainPatch from the specified map of raw messages.
func UnmarshalToolchainPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainPatch)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui_href-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainPost : Response structure for POST toolchain.
type ToolchainPost struct {
	// Toolchain ID.
	ID *string `json:"id" validate:"required"`

	// Toolchain name.
	Name *string `json:"name" validate:"required"`

	// Describes the toolchain.
	Description *string `json:"description" validate:"required"`

	// Account ID where toolchain can be found.
	AccountID *string `json:"account_id" validate:"required"`

	// Toolchain region.
	Location *string `json:"location" validate:"required"`

	// Resource group where the toolchain is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Toolchain CRN.
	CRN *string `json:"crn" validate:"required"`

	// URI that can be used to retrieve toolchain.
	Href *string `json:"href" validate:"required"`

	// URL of a user-facing user interface for this toolchain.
	UIHref *string `json:"ui_href" validate:"required"`

	// Toolchain creation timestamp.
	CreatedAt *strfmt.DateTime `json:"created_at" validate:"required"`

	// Latest toolchain update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Identity that created the toolchain.
	CreatedBy *string `json:"created_by" validate:"required"`
}

// UnmarshalToolchainPost unmarshals an instance of ToolchainPost from the specified map of raw messages.
func UnmarshalToolchainPost(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainPost)
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
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "account_id", &obj.AccountID)
	if err != nil {
		err = core.SDKErrorf(err, "", "account_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "ui_href", &obj.UIHref)
	if err != nil {
		err = core.SDKErrorf(err, "", "ui_href-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_by-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainPrototypePatch : Body structure for the update toolchain PATCH request.
type ToolchainPrototypePatch struct {
	// The name of the toolchain.
	Name *string `json:"name,omitempty"`

	// An optional description.
	Description *string `json:"description,omitempty"`
}

// UnmarshalToolchainPrototypePatch unmarshals an instance of ToolchainPrototypePatch from the specified map of raw messages.
func UnmarshalToolchainPrototypePatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainPrototypePatch)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ToolchainPrototypePatch
func (toolchainPrototypePatch *ToolchainPrototypePatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(toolchainPrototypePatch.Name) {
		_patch["name"] = toolchainPrototypePatch.Name
	}
	if !core.IsNil(toolchainPrototypePatch.Description) {
		_patch["description"] = toolchainPrototypePatch.Description
	}

	return
}

// ToolchainTool : Response structure for GET tool.
type ToolchainTool struct {
	// Tool ID.
	ID *string `json:"id" validate:"required"`

	// Resource group where the tool is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Tool CRN.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool. A table of `tool_type_id` values corresponding to each tool integration can
	// be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the tool is bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the tool is bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *ToolModelReferent `json:"referent" validate:"required"`

	// Name of the tool.
	Name *string `json:"name,omitempty"`

	// Latest tool update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool
	// integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the ToolchainTool.State property.
// Current configuration state of the tool.
const (
	ToolchainToolStateConfiguredConst = "configured"
	ToolchainToolStateConfiguringConst = "configuring"
	ToolchainToolStateMisconfiguredConst = "misconfigured"
	ToolchainToolStateUnconfiguredConst = "unconfigured"
)

// UnmarshalToolchainTool unmarshals an instance of ToolchainTool from the specified map of raw messages.
func UnmarshalToolchainTool(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainTool)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool_type_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalToolModelReferent)
	if err != nil {
		err = core.SDKErrorf(err, "", "referent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		err = core.SDKErrorf(err, "", "parameters-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolCollection : Response structure for GET tools.
type ToolchainToolCollection struct {
	// Maximum number of tools returned from collection.
	Limit *int64 `json:"limit" validate:"required"`

	// Total number of tools found in collection.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Information about retrieving first tool results from the collection.
	First *ToolchainToolCollectionFirst `json:"first" validate:"required"`

	// Information about retrieving previous tool results from the collection.
	Previous *ToolchainToolCollectionPrevious `json:"previous,omitempty"`

	// Information about retrieving next tool results from the collection.
	Next *ToolchainToolCollectionNext `json:"next,omitempty"`

	// Information about retrieving last tool results from the collection.
	Last *ToolchainToolCollectionLast `json:"last" validate:"required"`

	// Tool results returned from the collection.
	Tools []ToolModel `json:"tools" validate:"required"`
}

// UnmarshalToolchainToolCollection unmarshals an instance of ToolchainToolCollection from the specified map of raw messages.
func UnmarshalToolchainToolCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalToolchainToolCollectionFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalToolchainToolCollectionPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalToolchainToolCollectionNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalToolchainToolCollectionLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "tools", &obj.Tools, UnmarshalToolModel)
	if err != nil {
		err = core.SDKErrorf(err, "", "tools-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ToolchainToolCollection) GetNextStart() (*string, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	return resp.Next.Start, nil
}

// ToolchainToolCollectionFirst : Information about retrieving first tool results from the collection.
type ToolchainToolCollectionFirst struct {
	// URI that can be used to get first results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainToolCollectionFirst unmarshals an instance of ToolchainToolCollectionFirst from the specified map of raw messages.
func UnmarshalToolchainToolCollectionFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolCollectionFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolCollectionLast : Information about retrieving last tool results from the collection.
type ToolchainToolCollectionLast struct {
	// Cursor that can be used to get the last set of tool collections.
	Start *string `json:"start,omitempty"`

	// URI that can be used to get last results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainToolCollectionLast unmarshals an instance of ToolchainToolCollectionLast from the specified map of raw messages.
func UnmarshalToolchainToolCollectionLast(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolCollectionLast)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolCollectionNext : Information about retrieving next tool results from the collection.
type ToolchainToolCollectionNext struct {
	// Cursor that can be used to get the next set of tool collections.
	Start *string `json:"start,omitempty"`

	// URI that can be used to get next results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainToolCollectionNext unmarshals an instance of ToolchainToolCollectionNext from the specified map of raw messages.
func UnmarshalToolchainToolCollectionNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolCollectionNext)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolCollectionPrevious : Information about retrieving previous tool results from the collection.
type ToolchainToolCollectionPrevious struct {
	// Cursor that can be used to get the previous set of tool collections.
	Start *string `json:"start,omitempty"`

	// URI that can be used to get previous results from the collection.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalToolchainToolCollectionPrevious unmarshals an instance of ToolchainToolCollectionPrevious from the specified map of raw messages.
func UnmarshalToolchainToolCollectionPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolCollectionPrevious)
	err = core.UnmarshalPrimitive(m, "start", &obj.Start)
	if err != nil {
		err = core.SDKErrorf(err, "", "start-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolPatch : Response structure for PATCH tool.
type ToolchainToolPatch struct {
	// Tool ID.
	ID *string `json:"id" validate:"required"`

	// Resource group where the tool is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Tool CRN.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool. A table of `tool_type_id` values corresponding to each tool integration can
	// be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the tool is bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the tool is bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *ToolModelReferent `json:"referent" validate:"required"`

	// Name of the tool.
	Name *string `json:"name,omitempty"`

	// Latest tool update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool
	// integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the ToolchainToolPatch.State property.
// Current configuration state of the tool.
const (
	ToolchainToolPatchStateConfiguredConst = "configured"
	ToolchainToolPatchStateConfiguringConst = "configuring"
	ToolchainToolPatchStateMisconfiguredConst = "misconfigured"
	ToolchainToolPatchStateUnconfiguredConst = "unconfigured"
)

// UnmarshalToolchainToolPatch unmarshals an instance of ToolchainToolPatch from the specified map of raw messages.
func UnmarshalToolchainToolPatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolPatch)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool_type_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalToolModelReferent)
	if err != nil {
		err = core.SDKErrorf(err, "", "referent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		err = core.SDKErrorf(err, "", "parameters-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolPost : POST tool response body.
type ToolchainToolPost struct {
	// Tool ID.
	ID *string `json:"id" validate:"required"`

	// Resource group where the tool is located.
	ResourceGroupID *string `json:"resource_group_id" validate:"required"`

	// Tool CRN.
	CRN *string `json:"crn" validate:"required"`

	// The unique name of the provisioned tool. A table of `tool_type_id` values corresponding to each tool integration can
	// be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	ToolTypeID *string `json:"tool_type_id" validate:"required"`

	// ID of toolchain which the tool is bound to.
	ToolchainID *string `json:"toolchain_id" validate:"required"`

	// CRN of toolchain which the tool is bound to.
	ToolchainCRN *string `json:"toolchain_crn" validate:"required"`

	// URI representing the tool.
	Href *string `json:"href" validate:"required"`

	// Information on URIs to access this resource through the UI or API.
	Referent *ToolModelReferent `json:"referent" validate:"required"`

	// Name of the tool.
	Name *string `json:"name,omitempty"`

	// Latest tool update timestamp.
	UpdatedAt *strfmt.DateTime `json:"updated_at" validate:"required"`

	// Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool
	// integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	Parameters map[string]interface{} `json:"parameters" validate:"required"`

	// Current configuration state of the tool.
	State *string `json:"state" validate:"required"`
}

// Constants associated with the ToolchainToolPost.State property.
// Current configuration state of the tool.
const (
	ToolchainToolPostStateConfiguredConst = "configured"
	ToolchainToolPostStateConfiguringConst = "configuring"
	ToolchainToolPostStateMisconfiguredConst = "misconfigured"
	ToolchainToolPostStateUnconfiguredConst = "unconfigured"
)

// UnmarshalToolchainToolPost unmarshals an instance of ToolchainToolPost from the specified map of raw messages.
func UnmarshalToolchainToolPost(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolPost)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "resource_group_id", &obj.ResourceGroupID)
	if err != nil {
		err = core.SDKErrorf(err, "", "resource_group_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.CRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool_type_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_id", &obj.ToolchainID)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "toolchain_crn", &obj.ToolchainCRN)
	if err != nil {
		err = core.SDKErrorf(err, "", "toolchain_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "referent", &obj.Referent, UnmarshalToolModelReferent)
	if err != nil {
		err = core.SDKErrorf(err, "", "referent-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_at", &obj.UpdatedAt)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_at-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		err = core.SDKErrorf(err, "", "parameters-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "state", &obj.State)
	if err != nil {
		err = core.SDKErrorf(err, "", "state-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToolchainToolPrototypePatch : Details on the new tool.
type ToolchainToolPrototypePatch struct {
	// Name of the tool.
	Name *string `json:"name,omitempty"`

	// The unique short name of the tool that should be provisioned or updated. A table of `tool_type_id` values
	// corresponding to each tool integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	ToolTypeID *string `json:"tool_type_id,omitempty"`

	// Unique key-value pairs representing parameters to be used to create the tool. A list of parameters for each tool
	// integration can be found in the <a
	// href="https://cloud.ibm.com/docs/ContinuousDelivery?topic=ContinuousDelivery-integrations">Configuring tool
	// integrations page</a>.
	Parameters map[string]interface{} `json:"parameters,omitempty"`
}

// UnmarshalToolchainToolPrototypePatch unmarshals an instance of ToolchainToolPrototypePatch from the specified map of raw messages.
func UnmarshalToolchainToolPrototypePatch(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ToolchainToolPrototypePatch)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tool_type_id", &obj.ToolTypeID)
	if err != nil {
		err = core.SDKErrorf(err, "", "tool_type_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "parameters", &obj.Parameters)
	if err != nil {
		err = core.SDKErrorf(err, "", "parameters-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AsPatch returns a generic map representation of the ToolchainToolPrototypePatch
func (toolchainToolPrototypePatch *ToolchainToolPrototypePatch) AsPatch() (_patch map[string]interface{}, err error) {
	_patch = map[string]interface{}{}
	if !core.IsNil(toolchainToolPrototypePatch.Name) {
		_patch["name"] = toolchainToolPrototypePatch.Name
	}
	if !core.IsNil(toolchainToolPrototypePatch.ToolTypeID) {
		_patch["tool_type_id"] = toolchainToolPrototypePatch.ToolTypeID
	}
	if !core.IsNil(toolchainToolPrototypePatch.Parameters) {
		_patch["parameters"] = toolchainToolPrototypePatch.Parameters
	}

	return
}

// UpdateToolOptions : The UpdateTool options.
type UpdateToolOptions struct {
	// ID of the toolchain.
	ToolchainID *string `json:"toolchain_id" validate:"required,ne="`

	// ID of the tool bound to the toolchain.
	ToolID *string `json:"tool_id" validate:"required,ne="`

	// JSON Merge-Patch content for update_tool.
	ToolchainToolPrototypePatch map[string]interface{} `json:"ToolchainToolPrototype_patch" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateToolOptions : Instantiate UpdateToolOptions
func (*CdToolchainV2) NewUpdateToolOptions(toolchainID string, toolID string, toolchainToolPrototypePatch map[string]interface{}) *UpdateToolOptions {
	return &UpdateToolOptions{
		ToolchainID: core.StringPtr(toolchainID),
		ToolID: core.StringPtr(toolID),
		ToolchainToolPrototypePatch: toolchainToolPrototypePatch,
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

// SetToolchainToolPrototypePatch : Allow user to set ToolchainToolPrototypePatch
func (_options *UpdateToolOptions) SetToolchainToolPrototypePatch(toolchainToolPrototypePatch map[string]interface{}) *UpdateToolOptions {
	_options.ToolchainToolPrototypePatch = toolchainToolPrototypePatch
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

	// JSON Merge-Patch content for update_toolchain.
	ToolchainPrototypePatch map[string]interface{} `json:"ToolchainPrototype_patch" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateToolchainOptions : Instantiate UpdateToolchainOptions
func (*CdToolchainV2) NewUpdateToolchainOptions(toolchainID string, toolchainPrototypePatch map[string]interface{}) *UpdateToolchainOptions {
	return &UpdateToolchainOptions{
		ToolchainID: core.StringPtr(toolchainID),
		ToolchainPrototypePatch: toolchainPrototypePatch,
	}
}

// SetToolchainID : Allow user to set ToolchainID
func (_options *UpdateToolchainOptions) SetToolchainID(toolchainID string) *UpdateToolchainOptions {
	_options.ToolchainID = core.StringPtr(toolchainID)
	return _options
}

// SetToolchainPrototypePatch : Allow user to set ToolchainPrototypePatch
func (_options *UpdateToolchainOptions) SetToolchainPrototypePatch(toolchainPrototypePatch map[string]interface{}) *UpdateToolchainOptions {
	_options.ToolchainPrototypePatch = toolchainPrototypePatch
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateToolchainOptions) SetHeaders(param map[string]string) *UpdateToolchainOptions {
	options.Headers = param
	return options
}

//
// ToolchainsPager can be used to simplify the use of the "ListToolchains" method.
//
type ToolchainsPager struct {
	hasNext bool
	options *ListToolchainsOptions
	client  *CdToolchainV2
	pageContext struct {
		next *string
	}
}

// NewToolchainsPager returns a new ToolchainsPager instance.
func (cdToolchain *CdToolchainV2) NewToolchainsPager(options *ListToolchainsOptions) (pager *ToolchainsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListToolchainsOptions = *options
	pager = &ToolchainsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  cdToolchain,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ToolchainsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ToolchainsPager) GetNextWithContext(ctx context.Context) (page []ToolchainModel, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListToolchainsWithContext(ctx, pager.options)
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
	page = result.Toolchains

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ToolchainsPager) GetAllWithContext(ctx context.Context) (allItems []ToolchainModel, err error) {
	for pager.HasNext() {
		var nextPage []ToolchainModel
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
func (pager *ToolchainsPager) GetNext() (page []ToolchainModel, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ToolchainsPager) GetAll() (allItems []ToolchainModel, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// ToolsPager can be used to simplify the use of the "ListTools" method.
//
type ToolsPager struct {
	hasNext bool
	options *ListToolsOptions
	client  *CdToolchainV2
	pageContext struct {
		next *string
	}
}

// NewToolsPager returns a new ToolsPager instance.
func (cdToolchain *CdToolchainV2) NewToolsPager(options *ListToolsOptions) (pager *ToolsPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListToolsOptions = *options
	pager = &ToolsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  cdToolchain,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ToolsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ToolsPager) GetNextWithContext(ctx context.Context) (page []ToolModel, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListToolsWithContext(ctx, pager.options)
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
	page = result.Tools

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ToolsPager) GetAllWithContext(ctx context.Context) (allItems []ToolModel, err error) {
	for pager.HasNext() {
		var nextPage []ToolModel
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
func (pager *ToolsPager) GetNext() (page []ToolModel, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ToolsPager) GetAll() (allItems []ToolModel, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
