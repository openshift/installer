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
 * IBM OpenAPI SDK Code Generator Version: 3.97.1-d6730d2a-20241125-163317
 */

// Package appconfigurationv1 : Operations and models for the AppConfigurationV1 service
package appconfigurationv1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	common "github.com/IBM/appconfiguration-go-admin-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
)

// AppConfigurationV1 : IBM Cloud App Configuration is a centralized feature management and configuration service for
// use with web and mobile applications, microservices, and distributed environments.
//
// API Version: 1.0
// See: https://cloud.ibm.com/docs/app-configuration
type AppConfigurationV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://us-south.apprapp.cloud.ibm.com/apprapp/feature/v1/instances/provide-here-your-appconfig-instance-uuid"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "app_configuration"

const ParameterizedServiceURL = "https://{region}.apprapp.cloud.ibm.com/apprapp/feature/v1/instances/{guid}"

var defaultUrlVariables = map[string]string{
	"region": "us-south",
	"guid": "provide-here-your-appconfig-instance-uuid",
}

// AppConfigurationV1Options : Service options
type AppConfigurationV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewAppConfigurationV1UsingExternalConfig : constructs an instance of AppConfigurationV1 with passed in options and external configuration.
func NewAppConfigurationV1UsingExternalConfig(options *AppConfigurationV1Options) (appConfiguration *AppConfigurationV1, err error) {
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

	appConfiguration, err = NewAppConfigurationV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = appConfiguration.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = appConfiguration.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewAppConfigurationV1 : constructs an instance of AppConfigurationV1 with passed in options.
func NewAppConfigurationV1(options *AppConfigurationV1Options) (service *AppConfigurationV1, err error) {
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

	service = &AppConfigurationV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "appConfiguration" suitable for processing requests.
func (appConfiguration *AppConfigurationV1) Clone() *AppConfigurationV1 {
	if core.IsNil(appConfiguration) {
		return nil
	}
	clone := *appConfiguration
	clone.Service = appConfiguration.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (appConfiguration *AppConfigurationV1) SetServiceURL(url string) error {
	err := appConfiguration.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (appConfiguration *AppConfigurationV1) GetServiceURL() string {
	return appConfiguration.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (appConfiguration *AppConfigurationV1) SetDefaultHeaders(headers http.Header) {
	appConfiguration.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (appConfiguration *AppConfigurationV1) SetEnableGzipCompression(enableGzip bool) {
	appConfiguration.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (appConfiguration *AppConfigurationV1) GetEnableGzipCompression() bool {
	return appConfiguration.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (appConfiguration *AppConfigurationV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	appConfiguration.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (appConfiguration *AppConfigurationV1) DisableRetries() {
	appConfiguration.Service.DisableRetries()
}

// ListEnvironments : Get list of Environments
// List all the environments in the App Configuration service instance.
func (appConfiguration *AppConfigurationV1) ListEnvironments(listEnvironmentsOptions *ListEnvironmentsOptions) (result *EnvironmentList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListEnvironmentsWithContext(context.Background(), listEnvironmentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListEnvironmentsWithContext is an alternate form of the ListEnvironments method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListEnvironmentsWithContext(ctx context.Context, listEnvironmentsOptions *ListEnvironmentsOptions) (result *EnvironmentList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listEnvironmentsOptions, "listEnvironmentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listEnvironmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListEnvironments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listEnvironmentsOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*listEnvironmentsOptions.Expand))
	}
	if listEnvironmentsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listEnvironmentsOptions.Sort))
	}
	if listEnvironmentsOptions.Tags != nil {
		builder.AddQuery("tags", fmt.Sprint(*listEnvironmentsOptions.Tags))
	}
	if listEnvironmentsOptions.Include != nil {
		builder.AddQuery("include", strings.Join(listEnvironmentsOptions.Include, ","))
	}
	if listEnvironmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listEnvironmentsOptions.Limit))
	}
	if listEnvironmentsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listEnvironmentsOptions.Offset))
	}
	if listEnvironmentsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listEnvironmentsOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_environments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironmentList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateEnvironment : Create Environment
// Create an environment.
func (appConfiguration *AppConfigurationV1) CreateEnvironment(createEnvironmentOptions *CreateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreateEnvironmentWithContext(context.Background(), createEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateEnvironmentWithContext is an alternate form of the CreateEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateEnvironmentWithContext(ctx context.Context, createEnvironmentOptions *CreateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createEnvironmentOptions, "createEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createEnvironmentOptions, "createEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateEnvironment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createEnvironmentOptions.Name != nil {
		body["name"] = createEnvironmentOptions.Name
	}
	if createEnvironmentOptions.EnvironmentID != nil {
		body["environment_id"] = createEnvironmentOptions.EnvironmentID
	}
	if createEnvironmentOptions.Description != nil {
		body["description"] = createEnvironmentOptions.Description
	}
	if createEnvironmentOptions.Tags != nil {
		body["tags"] = createEnvironmentOptions.Tags
	}
	if createEnvironmentOptions.ColorCode != nil {
		body["color_code"] = createEnvironmentOptions.ColorCode
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateEnvironment : Update Environment
// Update an environment.
func (appConfiguration *AppConfigurationV1) UpdateEnvironment(updateEnvironmentOptions *UpdateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateEnvironmentWithContext(context.Background(), updateEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateEnvironmentWithContext is an alternate form of the UpdateEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateEnvironmentWithContext(ctx context.Context, updateEnvironmentOptions *UpdateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEnvironmentOptions, "updateEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateEnvironmentOptions, "updateEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *updateEnvironmentOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateEnvironment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateEnvironmentOptions.Name != nil {
		body["name"] = updateEnvironmentOptions.Name
	}
	if updateEnvironmentOptions.Description != nil {
		body["description"] = updateEnvironmentOptions.Description
	}
	if updateEnvironmentOptions.Tags != nil {
		body["tags"] = updateEnvironmentOptions.Tags
	}
	if updateEnvironmentOptions.ColorCode != nil {
		body["color_code"] = updateEnvironmentOptions.ColorCode
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetEnvironment : Get Environment
// Retrieve the details of the environment.
func (appConfiguration *AppConfigurationV1) GetEnvironment(getEnvironmentOptions *GetEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.GetEnvironmentWithContext(context.Background(), getEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetEnvironmentWithContext is an alternate form of the GetEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetEnvironmentWithContext(ctx context.Context, getEnvironmentOptions *GetEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEnvironmentOptions, "getEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getEnvironmentOptions, "getEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *getEnvironmentOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "GetEnvironment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getEnvironmentOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*getEnvironmentOptions.Expand))
	}
	if getEnvironmentOptions.Include != nil {
		builder.AddQuery("include", strings.Join(getEnvironmentOptions.Include, ","))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteEnvironment : Delete Environment
// Delete an Environment.
func (appConfiguration *AppConfigurationV1) DeleteEnvironment(deleteEnvironmentOptions *DeleteEnvironmentOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeleteEnvironmentWithContext(context.Background(), deleteEnvironmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteEnvironmentWithContext is an alternate form of the DeleteEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteEnvironmentWithContext(ctx context.Context, deleteEnvironmentOptions *DeleteEnvironmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEnvironmentOptions, "deleteEnvironmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteEnvironmentOptions, "deleteEnvironmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *deleteEnvironmentOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteEnvironmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteEnvironment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_environment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListCollections : Get list of Collections
// List of all the collections in the App Configuration service instance.
func (appConfiguration *AppConfigurationV1) ListCollections(listCollectionsOptions *ListCollectionsOptions) (result *CollectionList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListCollectionsWithContext(context.Background(), listCollectionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListCollectionsWithContext is an alternate form of the ListCollections method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListCollectionsWithContext(ctx context.Context, listCollectionsOptions *ListCollectionsOptions) (result *CollectionList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCollectionsOptions, "listCollectionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listCollectionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListCollections")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listCollectionsOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*listCollectionsOptions.Expand))
	}
	if listCollectionsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listCollectionsOptions.Sort))
	}
	if listCollectionsOptions.Tags != nil {
		builder.AddQuery("tags", fmt.Sprint(*listCollectionsOptions.Tags))
	}
	if listCollectionsOptions.Features != nil {
		builder.AddQuery("features", strings.Join(listCollectionsOptions.Features, ","))
	}
	if listCollectionsOptions.Properties != nil {
		builder.AddQuery("properties", strings.Join(listCollectionsOptions.Properties, ","))
	}
	if listCollectionsOptions.Include != nil {
		builder.AddQuery("include", strings.Join(listCollectionsOptions.Include, ","))
	}
	if listCollectionsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listCollectionsOptions.Limit))
	}
	if listCollectionsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listCollectionsOptions.Offset))
	}
	if listCollectionsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listCollectionsOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_collections", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectionList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateCollection : Create Collection
// Create a collection.
func (appConfiguration *AppConfigurationV1) CreateCollection(createCollectionOptions *CreateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreateCollectionWithContext(context.Background(), createCollectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateCollectionWithContext is an alternate form of the CreateCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateCollectionWithContext(ctx context.Context, createCollectionOptions *CreateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCollectionOptions, "createCollectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createCollectionOptions, "createCollectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createCollectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateCollection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createCollectionOptions.Name != nil {
		body["name"] = createCollectionOptions.Name
	}
	if createCollectionOptions.CollectionID != nil {
		body["collection_id"] = createCollectionOptions.CollectionID
	}
	if createCollectionOptions.Description != nil {
		body["description"] = createCollectionOptions.Description
	}
	if createCollectionOptions.Tags != nil {
		body["tags"] = createCollectionOptions.Tags
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_collection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectionLite)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateCollection : Update Collection
// Update the collection name, tags and description. Collection Id cannot be updated.
func (appConfiguration *AppConfigurationV1) UpdateCollection(updateCollectionOptions *UpdateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateCollectionWithContext(context.Background(), updateCollectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateCollectionWithContext is an alternate form of the UpdateCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateCollectionWithContext(ctx context.Context, updateCollectionOptions *UpdateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCollectionOptions, "updateCollectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateCollectionOptions, "updateCollectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"collection_id": *updateCollectionOptions.CollectionID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections/{collection_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateCollectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateCollection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateCollectionOptions.Name != nil {
		body["name"] = updateCollectionOptions.Name
	}
	if updateCollectionOptions.Description != nil {
		body["description"] = updateCollectionOptions.Description
	}
	if updateCollectionOptions.Tags != nil {
		body["tags"] = updateCollectionOptions.Tags
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_collection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectionLite)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetCollection : Get Collection
// Retrieve the details of the collection.
func (appConfiguration *AppConfigurationV1) GetCollection(getCollectionOptions *GetCollectionOptions) (result *Collection, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.GetCollectionWithContext(context.Background(), getCollectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCollectionWithContext is an alternate form of the GetCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetCollectionWithContext(ctx context.Context, getCollectionOptions *GetCollectionOptions) (result *Collection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCollectionOptions, "getCollectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCollectionOptions, "getCollectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"collection_id": *getCollectionOptions.CollectionID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections/{collection_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCollectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "GetCollection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getCollectionOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*getCollectionOptions.Expand))
	}
	if getCollectionOptions.Include != nil {
		builder.AddQuery("include", strings.Join(getCollectionOptions.Include, ","))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_collection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteCollection : Delete Collection
// Delete the collection.
func (appConfiguration *AppConfigurationV1) DeleteCollection(deleteCollectionOptions *DeleteCollectionOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeleteCollectionWithContext(context.Background(), deleteCollectionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteCollectionWithContext is an alternate form of the DeleteCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteCollectionWithContext(ctx context.Context, deleteCollectionOptions *DeleteCollectionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCollectionOptions, "deleteCollectionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteCollectionOptions, "deleteCollectionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"collection_id": *deleteCollectionOptions.CollectionID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections/{collection_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteCollectionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteCollection")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_collection", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListFeatures : Get list of Features
// List all the feature flags in the specified environment.
func (appConfiguration *AppConfigurationV1) ListFeatures(listFeaturesOptions *ListFeaturesOptions) (result *FeaturesList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListFeaturesWithContext(context.Background(), listFeaturesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListFeaturesWithContext is an alternate form of the ListFeatures method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListFeaturesWithContext(ctx context.Context, listFeaturesOptions *ListFeaturesOptions) (result *FeaturesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listFeaturesOptions, "listFeaturesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listFeaturesOptions, "listFeaturesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *listFeaturesOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listFeaturesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListFeatures")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listFeaturesOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*listFeaturesOptions.Expand))
	}
	if listFeaturesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listFeaturesOptions.Sort))
	}
	if listFeaturesOptions.Tags != nil {
		builder.AddQuery("tags", fmt.Sprint(*listFeaturesOptions.Tags))
	}
	if listFeaturesOptions.Collections != nil {
		builder.AddQuery("collections", strings.Join(listFeaturesOptions.Collections, ","))
	}
	if listFeaturesOptions.Segments != nil {
		builder.AddQuery("segments", strings.Join(listFeaturesOptions.Segments, ","))
	}
	if listFeaturesOptions.Include != nil {
		builder.AddQuery("include", strings.Join(listFeaturesOptions.Include, ","))
	}
	if listFeaturesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listFeaturesOptions.Limit))
	}
	if listFeaturesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listFeaturesOptions.Offset))
	}
	if listFeaturesOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listFeaturesOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_features", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeaturesList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateFeature : Create Feature
// Create a feature flag.
func (appConfiguration *AppConfigurationV1) CreateFeature(createFeatureOptions *CreateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreateFeatureWithContext(context.Background(), createFeatureOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateFeatureWithContext is an alternate form of the CreateFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateFeatureWithContext(ctx context.Context, createFeatureOptions *CreateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createFeatureOptions, "createFeatureOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createFeatureOptions, "createFeatureOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *createFeatureOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createFeatureOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateFeature")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createFeatureOptions.Name != nil {
		body["name"] = createFeatureOptions.Name
	}
	if createFeatureOptions.FeatureID != nil {
		body["feature_id"] = createFeatureOptions.FeatureID
	}
	if createFeatureOptions.Type != nil {
		body["type"] = createFeatureOptions.Type
	}
	if createFeatureOptions.EnabledValue != nil {
		body["enabled_value"] = createFeatureOptions.EnabledValue
	}
	if createFeatureOptions.DisabledValue != nil {
		body["disabled_value"] = createFeatureOptions.DisabledValue
	}
	if createFeatureOptions.Description != nil {
		body["description"] = createFeatureOptions.Description
	}
	if createFeatureOptions.Format != nil {
		body["format"] = createFeatureOptions.Format
	}
	if createFeatureOptions.Enabled != nil {
		body["enabled"] = createFeatureOptions.Enabled
	}
	if createFeatureOptions.RolloutPercentage != nil {
		body["rollout_percentage"] = createFeatureOptions.RolloutPercentage
	}
	if createFeatureOptions.Tags != nil {
		body["tags"] = createFeatureOptions.Tags
	}
	if createFeatureOptions.SegmentRules != nil {
		body["segment_rules"] = createFeatureOptions.SegmentRules
	}
	if createFeatureOptions.Collections != nil {
		body["collections"] = createFeatureOptions.Collections
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_feature", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateFeature : Update Feature
// Update a feature flag details.
func (appConfiguration *AppConfigurationV1) UpdateFeature(updateFeatureOptions *UpdateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateFeatureWithContext(context.Background(), updateFeatureOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateFeatureWithContext is an alternate form of the UpdateFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateFeatureWithContext(ctx context.Context, updateFeatureOptions *UpdateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFeatureOptions, "updateFeatureOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateFeatureOptions, "updateFeatureOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *updateFeatureOptions.EnvironmentID,
		"feature_id": *updateFeatureOptions.FeatureID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features/{feature_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateFeatureOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateFeature")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateFeatureOptions.Name != nil {
		body["name"] = updateFeatureOptions.Name
	}
	if updateFeatureOptions.Description != nil {
		body["description"] = updateFeatureOptions.Description
	}
	if updateFeatureOptions.EnabledValue != nil {
		body["enabled_value"] = updateFeatureOptions.EnabledValue
	}
	if updateFeatureOptions.DisabledValue != nil {
		body["disabled_value"] = updateFeatureOptions.DisabledValue
	}
	if updateFeatureOptions.Enabled != nil {
		body["enabled"] = updateFeatureOptions.Enabled
	}
	if updateFeatureOptions.RolloutPercentage != nil {
		body["rollout_percentage"] = updateFeatureOptions.RolloutPercentage
	}
	if updateFeatureOptions.Tags != nil {
		body["tags"] = updateFeatureOptions.Tags
	}
	if updateFeatureOptions.SegmentRules != nil {
		body["segment_rules"] = updateFeatureOptions.SegmentRules
	}
	if updateFeatureOptions.Collections != nil {
		body["collections"] = updateFeatureOptions.Collections
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_feature", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateFeatureValues : Update Feature Values
// Update the feature values. This method can be executed only by the `writer` role. This method allows the update of
// feature name, feature enabled_value, feature disabled_value, tags, description and feature segment rules, however
// this method does not allow toggling the feature flag and assigning feature to a collection.
func (appConfiguration *AppConfigurationV1) UpdateFeatureValues(updateFeatureValuesOptions *UpdateFeatureValuesOptions) (result *Feature, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateFeatureValuesWithContext(context.Background(), updateFeatureValuesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateFeatureValuesWithContext is an alternate form of the UpdateFeatureValues method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateFeatureValuesWithContext(ctx context.Context, updateFeatureValuesOptions *UpdateFeatureValuesOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFeatureValuesOptions, "updateFeatureValuesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateFeatureValuesOptions, "updateFeatureValuesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *updateFeatureValuesOptions.EnvironmentID,
		"feature_id": *updateFeatureValuesOptions.FeatureID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features/{feature_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateFeatureValuesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateFeatureValues")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateFeatureValuesOptions.Name != nil {
		body["name"] = updateFeatureValuesOptions.Name
	}
	if updateFeatureValuesOptions.Description != nil {
		body["description"] = updateFeatureValuesOptions.Description
	}
	if updateFeatureValuesOptions.Tags != nil {
		body["tags"] = updateFeatureValuesOptions.Tags
	}
	if updateFeatureValuesOptions.EnabledValue != nil {
		body["enabled_value"] = updateFeatureValuesOptions.EnabledValue
	}
	if updateFeatureValuesOptions.DisabledValue != nil {
		body["disabled_value"] = updateFeatureValuesOptions.DisabledValue
	}
	if updateFeatureValuesOptions.RolloutPercentage != nil {
		body["rollout_percentage"] = updateFeatureValuesOptions.RolloutPercentage
	}
	if updateFeatureValuesOptions.SegmentRules != nil {
		body["segment_rules"] = updateFeatureValuesOptions.SegmentRules
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_feature_values", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetFeature : Get Feature
// Retrieve details of a feature.
func (appConfiguration *AppConfigurationV1) GetFeature(getFeatureOptions *GetFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.GetFeatureWithContext(context.Background(), getFeatureOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetFeatureWithContext is an alternate form of the GetFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetFeatureWithContext(ctx context.Context, getFeatureOptions *GetFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getFeatureOptions, "getFeatureOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getFeatureOptions, "getFeatureOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *getFeatureOptions.EnvironmentID,
		"feature_id": *getFeatureOptions.FeatureID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features/{feature_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getFeatureOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "GetFeature")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getFeatureOptions.Include != nil {
		builder.AddQuery("include", strings.Join(getFeatureOptions.Include, ","))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_feature", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteFeature : Delete Feature
// Delete a feature flag.
func (appConfiguration *AppConfigurationV1) DeleteFeature(deleteFeatureOptions *DeleteFeatureOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeleteFeatureWithContext(context.Background(), deleteFeatureOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteFeatureWithContext is an alternate form of the DeleteFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteFeatureWithContext(ctx context.Context, deleteFeatureOptions *DeleteFeatureOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteFeatureOptions, "deleteFeatureOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteFeatureOptions, "deleteFeatureOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *deleteFeatureOptions.EnvironmentID,
		"feature_id": *deleteFeatureOptions.FeatureID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features/{feature_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteFeatureOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteFeature")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_feature", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ToggleFeature : Toggle Feature
// Toggle a feature.
func (appConfiguration *AppConfigurationV1) ToggleFeature(toggleFeatureOptions *ToggleFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ToggleFeatureWithContext(context.Background(), toggleFeatureOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ToggleFeatureWithContext is an alternate form of the ToggleFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ToggleFeatureWithContext(ctx context.Context, toggleFeatureOptions *ToggleFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(toggleFeatureOptions, "toggleFeatureOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(toggleFeatureOptions, "toggleFeatureOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *toggleFeatureOptions.EnvironmentID,
		"feature_id": *toggleFeatureOptions.FeatureID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/features/{feature_id}/toggle`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range toggleFeatureOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ToggleFeature")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if toggleFeatureOptions.Enabled != nil {
		body["enabled"] = toggleFeatureOptions.Enabled
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "toggle_feature", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListProperties : Get list of Properties
// List all the properties in the specified environment.
func (appConfiguration *AppConfigurationV1) ListProperties(listPropertiesOptions *ListPropertiesOptions) (result *PropertiesList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListPropertiesWithContext(context.Background(), listPropertiesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListPropertiesWithContext is an alternate form of the ListProperties method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListPropertiesWithContext(ctx context.Context, listPropertiesOptions *ListPropertiesOptions) (result *PropertiesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPropertiesOptions, "listPropertiesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listPropertiesOptions, "listPropertiesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *listPropertiesOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/properties`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listPropertiesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListProperties")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listPropertiesOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*listPropertiesOptions.Expand))
	}
	if listPropertiesOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listPropertiesOptions.Sort))
	}
	if listPropertiesOptions.Tags != nil {
		builder.AddQuery("tags", fmt.Sprint(*listPropertiesOptions.Tags))
	}
	if listPropertiesOptions.Collections != nil {
		builder.AddQuery("collections", strings.Join(listPropertiesOptions.Collections, ","))
	}
	if listPropertiesOptions.Segments != nil {
		builder.AddQuery("segments", strings.Join(listPropertiesOptions.Segments, ","))
	}
	if listPropertiesOptions.Include != nil {
		builder.AddQuery("include", strings.Join(listPropertiesOptions.Include, ","))
	}
	if listPropertiesOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listPropertiesOptions.Limit))
	}
	if listPropertiesOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listPropertiesOptions.Offset))
	}
	if listPropertiesOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listPropertiesOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_properties", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPropertiesList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateProperty : Create Property
// Create a Property.
func (appConfiguration *AppConfigurationV1) CreateProperty(createPropertyOptions *CreatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreatePropertyWithContext(context.Background(), createPropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreatePropertyWithContext is an alternate form of the CreateProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreatePropertyWithContext(ctx context.Context, createPropertyOptions *CreatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPropertyOptions, "createPropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createPropertyOptions, "createPropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *createPropertyOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/properties`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createPropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createPropertyOptions.Name != nil {
		body["name"] = createPropertyOptions.Name
	}
	if createPropertyOptions.PropertyID != nil {
		body["property_id"] = createPropertyOptions.PropertyID
	}
	if createPropertyOptions.Type != nil {
		body["type"] = createPropertyOptions.Type
	}
	if createPropertyOptions.Value != nil {
		body["value"] = createPropertyOptions.Value
	}
	if createPropertyOptions.Description != nil {
		body["description"] = createPropertyOptions.Description
	}
	if createPropertyOptions.Format != nil {
		body["format"] = createPropertyOptions.Format
	}
	if createPropertyOptions.Tags != nil {
		body["tags"] = createPropertyOptions.Tags
	}
	if createPropertyOptions.SegmentRules != nil {
		body["segment_rules"] = createPropertyOptions.SegmentRules
	}
	if createPropertyOptions.Collections != nil {
		body["collections"] = createPropertyOptions.Collections
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_property", getServiceComponentInfo())
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

// UpdateProperty : Update Property
// Update a Property.
func (appConfiguration *AppConfigurationV1) UpdateProperty(updatePropertyOptions *UpdatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdatePropertyWithContext(context.Background(), updatePropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdatePropertyWithContext is an alternate form of the UpdateProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdatePropertyWithContext(ctx context.Context, updatePropertyOptions *UpdatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePropertyOptions, "updatePropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updatePropertyOptions, "updatePropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *updatePropertyOptions.EnvironmentID,
		"property_id": *updatePropertyOptions.PropertyID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/properties/{property_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updatePropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updatePropertyOptions.Name != nil {
		body["name"] = updatePropertyOptions.Name
	}
	if updatePropertyOptions.Description != nil {
		body["description"] = updatePropertyOptions.Description
	}
	if updatePropertyOptions.Value != nil {
		body["value"] = updatePropertyOptions.Value
	}
	if updatePropertyOptions.Tags != nil {
		body["tags"] = updatePropertyOptions.Tags
	}
	if updatePropertyOptions.SegmentRules != nil {
		body["segment_rules"] = updatePropertyOptions.SegmentRules
	}
	if updatePropertyOptions.Collections != nil {
		body["collections"] = updatePropertyOptions.Collections
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_property", getServiceComponentInfo())
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

// UpdatePropertyValues : Update Property values
// Update the property values. This method can be executed by the `writer` role. Property value and targeting rules can
// be updated, however this method does not allow assigning property to a collection.
func (appConfiguration *AppConfigurationV1) UpdatePropertyValues(updatePropertyValuesOptions *UpdatePropertyValuesOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdatePropertyValuesWithContext(context.Background(), updatePropertyValuesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdatePropertyValuesWithContext is an alternate form of the UpdatePropertyValues method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdatePropertyValuesWithContext(ctx context.Context, updatePropertyValuesOptions *UpdatePropertyValuesOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePropertyValuesOptions, "updatePropertyValuesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updatePropertyValuesOptions, "updatePropertyValuesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *updatePropertyValuesOptions.EnvironmentID,
		"property_id": *updatePropertyValuesOptions.PropertyID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/properties/{property_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updatePropertyValuesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdatePropertyValues")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updatePropertyValuesOptions.Name != nil {
		body["name"] = updatePropertyValuesOptions.Name
	}
	if updatePropertyValuesOptions.Description != nil {
		body["description"] = updatePropertyValuesOptions.Description
	}
	if updatePropertyValuesOptions.Tags != nil {
		body["tags"] = updatePropertyValuesOptions.Tags
	}
	if updatePropertyValuesOptions.Value != nil {
		body["value"] = updatePropertyValuesOptions.Value
	}
	if updatePropertyValuesOptions.SegmentRules != nil {
		body["segment_rules"] = updatePropertyValuesOptions.SegmentRules
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_property_values", getServiceComponentInfo())
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

// GetProperty : Get Property
// Retrieve details of a property.
func (appConfiguration *AppConfigurationV1) GetProperty(getPropertyOptions *GetPropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.GetPropertyWithContext(context.Background(), getPropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetPropertyWithContext is an alternate form of the GetProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetPropertyWithContext(ctx context.Context, getPropertyOptions *GetPropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPropertyOptions, "getPropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getPropertyOptions, "getPropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *getPropertyOptions.EnvironmentID,
		"property_id": *getPropertyOptions.PropertyID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/properties/{property_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getPropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "GetProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getPropertyOptions.Include != nil {
		builder.AddQuery("include", strings.Join(getPropertyOptions.Include, ","))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_property", getServiceComponentInfo())
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

// DeleteProperty : Delete Property
// Delete a Property.
func (appConfiguration *AppConfigurationV1) DeleteProperty(deletePropertyOptions *DeletePropertyOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeletePropertyWithContext(context.Background(), deletePropertyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeletePropertyWithContext is an alternate form of the DeleteProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeletePropertyWithContext(ctx context.Context, deletePropertyOptions *DeletePropertyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePropertyOptions, "deletePropertyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deletePropertyOptions, "deletePropertyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *deletePropertyOptions.EnvironmentID,
		"property_id": *deletePropertyOptions.PropertyID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/properties/{property_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deletePropertyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteProperty")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_property", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListSegments : Get list of Segments
// List all the segments.
func (appConfiguration *AppConfigurationV1) ListSegments(listSegmentsOptions *ListSegmentsOptions) (result *SegmentsList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListSegmentsWithContext(context.Background(), listSegmentsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListSegmentsWithContext is an alternate form of the ListSegments method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListSegmentsWithContext(ctx context.Context, listSegmentsOptions *ListSegmentsOptions) (result *SegmentsList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSegmentsOptions, "listSegmentsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listSegmentsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListSegments")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSegmentsOptions.Expand != nil {
		builder.AddQuery("expand", fmt.Sprint(*listSegmentsOptions.Expand))
	}
	if listSegmentsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listSegmentsOptions.Sort))
	}
	if listSegmentsOptions.Tags != nil {
		builder.AddQuery("tags", fmt.Sprint(*listSegmentsOptions.Tags))
	}
	if listSegmentsOptions.Include != nil {
		builder.AddQuery("include", fmt.Sprint(*listSegmentsOptions.Include))
	}
	if listSegmentsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSegmentsOptions.Limit))
	}
	if listSegmentsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSegmentsOptions.Offset))
	}
	if listSegmentsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listSegmentsOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_segments", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegmentsList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateSegment : Create Segment
// Create a segment.
func (appConfiguration *AppConfigurationV1) CreateSegment(createSegmentOptions *CreateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreateSegmentWithContext(context.Background(), createSegmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateSegmentWithContext is an alternate form of the CreateSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateSegmentWithContext(ctx context.Context, createSegmentOptions *CreateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSegmentOptions, "createSegmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createSegmentOptions, "createSegmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createSegmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateSegment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSegmentOptions.Name != nil {
		body["name"] = createSegmentOptions.Name
	}
	if createSegmentOptions.SegmentID != nil {
		body["segment_id"] = createSegmentOptions.SegmentID
	}
	if createSegmentOptions.Rules != nil {
		body["rules"] = createSegmentOptions.Rules
	}
	if createSegmentOptions.Description != nil {
		body["description"] = createSegmentOptions.Description
	}
	if createSegmentOptions.Tags != nil {
		body["tags"] = createSegmentOptions.Tags
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_segment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateSegment : Update Segment
// Update the segment properties.
func (appConfiguration *AppConfigurationV1) UpdateSegment(updateSegmentOptions *UpdateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateSegmentWithContext(context.Background(), updateSegmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateSegmentWithContext is an alternate form of the UpdateSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateSegmentWithContext(ctx context.Context, updateSegmentOptions *UpdateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSegmentOptions, "updateSegmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateSegmentOptions, "updateSegmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"segment_id": *updateSegmentOptions.SegmentID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments/{segment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateSegmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateSegment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSegmentOptions.Name != nil {
		body["name"] = updateSegmentOptions.Name
	}
	if updateSegmentOptions.Description != nil {
		body["description"] = updateSegmentOptions.Description
	}
	if updateSegmentOptions.Tags != nil {
		body["tags"] = updateSegmentOptions.Tags
	}
	if updateSegmentOptions.Rules != nil {
		body["rules"] = updateSegmentOptions.Rules
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_segment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetSegment : Get Segment
// Retrieve details of a segment.
func (appConfiguration *AppConfigurationV1) GetSegment(getSegmentOptions *GetSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.GetSegmentWithContext(context.Background(), getSegmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetSegmentWithContext is an alternate form of the GetSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetSegmentWithContext(ctx context.Context, getSegmentOptions *GetSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSegmentOptions, "getSegmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getSegmentOptions, "getSegmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"segment_id": *getSegmentOptions.SegmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments/{segment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getSegmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "GetSegment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getSegmentOptions.Include != nil {
		builder.AddQuery("include", strings.Join(getSegmentOptions.Include, ","))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_segment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegment)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteSegment : Delete Segment
// Delete a segment.
func (appConfiguration *AppConfigurationV1) DeleteSegment(deleteSegmentOptions *DeleteSegmentOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeleteSegmentWithContext(context.Background(), deleteSegmentOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteSegmentWithContext is an alternate form of the DeleteSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteSegmentWithContext(ctx context.Context, deleteSegmentOptions *DeleteSegmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSegmentOptions, "deleteSegmentOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteSegmentOptions, "deleteSegmentOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"segment_id": *deleteSegmentOptions.SegmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments/{segment_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteSegmentOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteSegment")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_segment", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListSnapshots : Get list of Git configs
// List all the Git configs.
func (appConfiguration *AppConfigurationV1) ListSnapshots(listSnapshotsOptions *ListSnapshotsOptions) (result *GitConfigList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListSnapshotsWithContext(context.Background(), listSnapshotsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListSnapshotsWithContext is an alternate form of the ListSnapshots method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListSnapshotsWithContext(ctx context.Context, listSnapshotsOptions *ListSnapshotsOptions) (result *GitConfigList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSnapshotsOptions, "listSnapshotsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listSnapshotsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListSnapshots")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSnapshotsOptions.Sort != nil {
		builder.AddQuery("sort", fmt.Sprint(*listSnapshotsOptions.Sort))
	}
	if listSnapshotsOptions.CollectionID != nil {
		builder.AddQuery("collection_id", fmt.Sprint(*listSnapshotsOptions.CollectionID))
	}
	if listSnapshotsOptions.EnvironmentID != nil {
		builder.AddQuery("environment_id", fmt.Sprint(*listSnapshotsOptions.EnvironmentID))
	}
	if listSnapshotsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listSnapshotsOptions.Limit))
	}
	if listSnapshotsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listSnapshotsOptions.Offset))
	}
	if listSnapshotsOptions.Search != nil {
		builder.AddQuery("search", fmt.Sprint(*listSnapshotsOptions.Search))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_snapshots", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGitConfigList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateGitconfig : Create Git config
// Create a Git config.
func (appConfiguration *AppConfigurationV1) CreateGitconfig(createGitconfigOptions *CreateGitconfigOptions) (result *CreateGitConfigResponse, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreateGitconfigWithContext(context.Background(), createGitconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateGitconfigWithContext is an alternate form of the CreateGitconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateGitconfigWithContext(ctx context.Context, createGitconfigOptions *CreateGitconfigOptions) (result *CreateGitConfigResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createGitconfigOptions, "createGitconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createGitconfigOptions, "createGitconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createGitconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateGitconfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createGitconfigOptions.GitConfigName != nil {
		body["git_config_name"] = createGitconfigOptions.GitConfigName
	}
	if createGitconfigOptions.GitConfigID != nil {
		body["git_config_id"] = createGitconfigOptions.GitConfigID
	}
	if createGitconfigOptions.CollectionID != nil {
		body["collection_id"] = createGitconfigOptions.CollectionID
	}
	if createGitconfigOptions.EnvironmentID != nil {
		body["environment_id"] = createGitconfigOptions.EnvironmentID
	}
	if createGitconfigOptions.GitURL != nil {
		body["git_url"] = createGitconfigOptions.GitURL
	}
	if createGitconfigOptions.GitBranch != nil {
		body["git_branch"] = createGitconfigOptions.GitBranch
	}
	if createGitconfigOptions.GitFilePath != nil {
		body["git_file_path"] = createGitconfigOptions.GitFilePath
	}
	if createGitconfigOptions.GitToken != nil {
		body["git_token"] = createGitconfigOptions.GitToken
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_gitconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateGitConfigResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateGitconfig : Update Git Config
// Update the gitconfig properties.
func (appConfiguration *AppConfigurationV1) UpdateGitconfig(updateGitconfigOptions *UpdateGitconfigOptions) (result *GitConfig, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateGitconfigWithContext(context.Background(), updateGitconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateGitconfigWithContext is an alternate form of the UpdateGitconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateGitconfigWithContext(ctx context.Context, updateGitconfigOptions *UpdateGitconfigOptions) (result *GitConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateGitconfigOptions, "updateGitconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateGitconfigOptions, "updateGitconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"git_config_id": *updateGitconfigOptions.GitConfigID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs/{git_config_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateGitconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateGitconfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateGitconfigOptions.GitConfigName != nil {
		body["git_config_name"] = updateGitconfigOptions.GitConfigName
	}
	if updateGitconfigOptions.CollectionID != nil {
		body["collection_id"] = updateGitconfigOptions.CollectionID
	}
	if updateGitconfigOptions.EnvironmentID != nil {
		body["environment_id"] = updateGitconfigOptions.EnvironmentID
	}
	if updateGitconfigOptions.GitURL != nil {
		body["git_url"] = updateGitconfigOptions.GitURL
	}
	if updateGitconfigOptions.GitBranch != nil {
		body["git_branch"] = updateGitconfigOptions.GitBranch
	}
	if updateGitconfigOptions.GitFilePath != nil {
		body["git_file_path"] = updateGitconfigOptions.GitFilePath
	}
	if updateGitconfigOptions.GitToken != nil {
		body["git_token"] = updateGitconfigOptions.GitToken
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_gitconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGitConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetGitconfig : Get Git Config
// Retrieve details of a gitconfig.
func (appConfiguration *AppConfigurationV1) GetGitconfig(getGitconfigOptions *GetGitconfigOptions) (result *GitConfig, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.GetGitconfigWithContext(context.Background(), getGitconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetGitconfigWithContext is an alternate form of the GetGitconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetGitconfigWithContext(ctx context.Context, getGitconfigOptions *GetGitconfigOptions) (result *GitConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getGitconfigOptions, "getGitconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getGitconfigOptions, "getGitconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"git_config_id": *getGitconfigOptions.GitConfigID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs/{git_config_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getGitconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "GetGitconfig")
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_gitconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGitConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteGitconfig : Delete Git Config
// Delete a gitconfig.
func (appConfiguration *AppConfigurationV1) DeleteGitconfig(deleteGitconfigOptions *DeleteGitconfigOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeleteGitconfigWithContext(context.Background(), deleteGitconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteGitconfigWithContext is an alternate form of the DeleteGitconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteGitconfigWithContext(ctx context.Context, deleteGitconfigOptions *DeleteGitconfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteGitconfigOptions, "deleteGitconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteGitconfigOptions, "deleteGitconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"git_config_id": *deleteGitconfigOptions.GitConfigID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs/{git_config_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteGitconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteGitconfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_gitconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// PromoteGitconfig : Promote configuration
// Promote configuration, this api will write or update your chosen configuration to the GitHub based on the git url,
// file path and branch data. In simple words this api will create or updates the bootstrap json file.
// Deprecated: this method is deprecated and may be removed in a future release.
func (appConfiguration *AppConfigurationV1) PromoteGitconfig(promoteGitconfigOptions *PromoteGitconfigOptions) (result *GitConfigPromote, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.PromoteGitconfigWithContext(context.Background(), promoteGitconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PromoteGitconfigWithContext is an alternate form of the PromoteGitconfig method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (appConfiguration *AppConfigurationV1) PromoteGitconfigWithContext(ctx context.Context, promoteGitconfigOptions *PromoteGitconfigOptions) (result *GitConfigPromote, response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: PromoteGitconfig")
	err = core.ValidateNotNil(promoteGitconfigOptions, "promoteGitconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(promoteGitconfigOptions, "promoteGitconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"git_config_id": *promoteGitconfigOptions.GitConfigID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs/{git_config_id}/promote`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range promoteGitconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "PromoteGitconfig")
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "promote_gitconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGitConfigPromote)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// RestoreGitconfig : Restore configuration
// Restore configuration, this api will write or update your chosen configuration from the GitHub to App configuration
// instance. The api will read the contents in the json file that was created using promote API and recreate or updates
// the App configuration instance with the file contents like properties, features and segments.
// Deprecated: this method is deprecated and may be removed in a future release.
func (appConfiguration *AppConfigurationV1) RestoreGitconfig(restoreGitconfigOptions *RestoreGitconfigOptions) (result *GitConfigRestore, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.RestoreGitconfigWithContext(context.Background(), restoreGitconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// RestoreGitconfigWithContext is an alternate form of the RestoreGitconfig method which supports a Context parameter
// Deprecated: this method is deprecated and may be removed in a future release.
func (appConfiguration *AppConfigurationV1) RestoreGitconfigWithContext(ctx context.Context, restoreGitconfigOptions *RestoreGitconfigOptions) (result *GitConfigRestore, response *core.DetailedResponse, err error) {
	core.GetLogger().Warn("A deprecated operation has been invoked: RestoreGitconfig")
	err = core.ValidateNotNil(restoreGitconfigOptions, "restoreGitconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(restoreGitconfigOptions, "restoreGitconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"git_config_id": *restoreGitconfigOptions.GitConfigID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/gitconfigs/{git_config_id}/restore`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range restoreGitconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "RestoreGitconfig")
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "restore_gitconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalGitConfigRestore)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListOriginconfigs : Get list of Origin Configs
// List all the Origin Configs.
func (appConfiguration *AppConfigurationV1) ListOriginconfigs(listOriginconfigsOptions *ListOriginconfigsOptions) (result *OriginConfigList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListOriginconfigsWithContext(context.Background(), listOriginconfigsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListOriginconfigsWithContext is an alternate form of the ListOriginconfigs method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListOriginconfigsWithContext(ctx context.Context, listOriginconfigsOptions *ListOriginconfigsOptions) (result *OriginConfigList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listOriginconfigsOptions, "listOriginconfigsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/originconfigs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listOriginconfigsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListOriginconfigs")
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_originconfigs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOriginConfigList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateOriginconfigs : Update Origin Configs
// Update the Origin Configs.
func (appConfiguration *AppConfigurationV1) UpdateOriginconfigs(updateOriginconfigsOptions *UpdateOriginconfigsOptions) (result *OriginConfigList, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateOriginconfigsWithContext(context.Background(), updateOriginconfigsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateOriginconfigsWithContext is an alternate form of the UpdateOriginconfigs method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateOriginconfigsWithContext(ctx context.Context, updateOriginconfigsOptions *UpdateOriginconfigsOptions) (result *OriginConfigList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateOriginconfigsOptions, "updateOriginconfigsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateOriginconfigsOptions, "updateOriginconfigsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/originconfigs`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateOriginconfigsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateOriginconfigs")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateOriginconfigsOptions.AllowedOrigins != nil {
		body["allowed_origins"] = updateOriginconfigsOptions.AllowedOrigins
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_originconfigs", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalOriginConfigList)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListWorkflowconfig : Get Workflow Config
// Get the environment specific workflow configs.
func (appConfiguration *AppConfigurationV1) ListWorkflowconfig(listWorkflowconfigOptions *ListWorkflowconfigOptions) (result ListWorkflowconfigResponseIntf, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListWorkflowconfigWithContext(context.Background(), listWorkflowconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListWorkflowconfigWithContext is an alternate form of the ListWorkflowconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListWorkflowconfigWithContext(ctx context.Context, listWorkflowconfigOptions *ListWorkflowconfigOptions) (result ListWorkflowconfigResponseIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listWorkflowconfigOptions, "listWorkflowconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listWorkflowconfigOptions, "listWorkflowconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *listWorkflowconfigOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/workflowconfigs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listWorkflowconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListWorkflowconfig")
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_workflowconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListWorkflowconfigResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateWorkflowconfig : Create Workflow config
// Create a Workflow.
func (appConfiguration *AppConfigurationV1) CreateWorkflowconfig(createWorkflowconfigOptions *CreateWorkflowconfigOptions) (result CreateWorkflowconfigResponseIntf, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.CreateWorkflowconfigWithContext(context.Background(), createWorkflowconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateWorkflowconfigWithContext is an alternate form of the CreateWorkflowconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateWorkflowconfigWithContext(ctx context.Context, createWorkflowconfigOptions *CreateWorkflowconfigOptions) (result CreateWorkflowconfigResponseIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createWorkflowconfigOptions, "createWorkflowconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createWorkflowconfigOptions, "createWorkflowconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *createWorkflowconfigOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/workflowconfigs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createWorkflowconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "CreateWorkflowconfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(createWorkflowconfigOptions.WorkflowConfig)
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_Workflowconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCreateWorkflowconfigResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// UpdateWorkflowconfig : Update Workflow config
// Update a Workflow.
func (appConfiguration *AppConfigurationV1) UpdateWorkflowconfig(updateWorkflowconfigOptions *UpdateWorkflowconfigOptions) (result UpdateWorkflowconfigResponseIntf, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.UpdateWorkflowconfigWithContext(context.Background(), updateWorkflowconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UpdateWorkflowconfigWithContext is an alternate form of the UpdateWorkflowconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateWorkflowconfigWithContext(ctx context.Context, updateWorkflowconfigOptions *UpdateWorkflowconfigOptions) (result UpdateWorkflowconfigResponseIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateWorkflowconfigOptions, "updateWorkflowconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(updateWorkflowconfigOptions, "updateWorkflowconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *updateWorkflowconfigOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/workflowconfigs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range updateWorkflowconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "UpdateWorkflowconfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	_, err = builder.SetBodyContentJSON(updateWorkflowconfigOptions.UpdateWorkflowConfig)
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "update_Workflowconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUpdateWorkflowconfigResponse)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteWorkflowconfig : Delete  Workflow config
// Delete a  Workflow config.
func (appConfiguration *AppConfigurationV1) DeleteWorkflowconfig(deleteWorkflowconfigOptions *DeleteWorkflowconfigOptions) (response *core.DetailedResponse, err error) {
	response, err = appConfiguration.DeleteWorkflowconfigWithContext(context.Background(), deleteWorkflowconfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteWorkflowconfigWithContext is an alternate form of the DeleteWorkflowconfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteWorkflowconfigWithContext(ctx context.Context, deleteWorkflowconfigOptions *DeleteWorkflowconfigOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteWorkflowconfigOptions, "deleteWorkflowconfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteWorkflowconfigOptions, "deleteWorkflowconfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"environment_id": *deleteWorkflowconfigOptions.EnvironmentID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments/{environment_id}/workflowconfigs`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteWorkflowconfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "DeleteWorkflowconfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_workflowconfig", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ImportConfig : Import instance configuration
// Import configuration to the instance.
func (appConfiguration *AppConfigurationV1) ImportConfig(importConfigOptions *ImportConfigOptions) (result *ImportConfig, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ImportConfigWithContext(context.Background(), importConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ImportConfigWithContext is an alternate form of the ImportConfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ImportConfigWithContext(ctx context.Context, importConfigOptions *ImportConfigOptions) (result *ImportConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importConfigOptions, "importConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(importConfigOptions, "importConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range importConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ImportConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	if importConfigOptions.Clean != nil {
		builder.AddQuery("clean", fmt.Sprint(*importConfigOptions.Clean))
	}

	body := make(map[string]interface{})
	if importConfigOptions.Environments != nil {
		body["environments"] = importConfigOptions.Environments
	}
	if importConfigOptions.Collections != nil {
		body["collections"] = importConfigOptions.Collections
	}
	if importConfigOptions.Segments != nil {
		body["segments"] = importConfigOptions.Segments
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "import_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImportConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListInstanceConfig : Export instance configuration
// Get the instance configuration.
func (appConfiguration *AppConfigurationV1) ListInstanceConfig(listInstanceConfigOptions *ListInstanceConfigOptions) (result *ImportConfig, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.ListInstanceConfigWithContext(context.Background(), listInstanceConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListInstanceConfigWithContext is an alternate form of the ListInstanceConfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListInstanceConfigWithContext(ctx context.Context, listInstanceConfigOptions *ListInstanceConfigOptions) (result *ImportConfig, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listInstanceConfigOptions, "listInstanceConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listInstanceConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "ListInstanceConfig")
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_instance_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalImportConfig)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// PromoteRestoreConfig : Promote or Restore snapshot configuration
// This api will either promote or restore your chosen configuration from or to the GitHub based on the git url, file
// path and branch data.
func (appConfiguration *AppConfigurationV1) PromoteRestoreConfig(promoteRestoreConfigOptions *PromoteRestoreConfigOptions) (result ConfigActionIntf, response *core.DetailedResponse, err error) {
	result, response, err = appConfiguration.PromoteRestoreConfigWithContext(context.Background(), promoteRestoreConfigOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// PromoteRestoreConfigWithContext is an alternate form of the PromoteRestoreConfig method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) PromoteRestoreConfigWithContext(ctx context.Context, promoteRestoreConfigOptions *PromoteRestoreConfigOptions) (result ConfigActionIntf, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(promoteRestoreConfigOptions, "promoteRestoreConfigOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(promoteRestoreConfigOptions, "promoteRestoreConfigOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/config`, nil)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range promoteRestoreConfigOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("app_configuration", "V1", "PromoteRestoreConfig")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddQuery("git_config_id", fmt.Sprint(*promoteRestoreConfigOptions.GitConfigID))
	builder.AddQuery("action", fmt.Sprint(*promoteRestoreConfigOptions.Action))

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "promote_restore_config", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigAction)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.0")
}

// Collection : Details of the collection.
type Collection struct {
	// Collection name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Collection Id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Creation time of the collection.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last updated time of the collection data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Collection URL.
	Href *string `json:"href,omitempty"`

	// List of Features associated with the collection.
	Features []FeatureOutput `json:"features,omitempty"`

	// List of properties associated with the collection.
	Properties []PropertyOutput `json:"properties,omitempty"`

	// List of snapshots associated with the collection.
	Snapshots []SnapshotOutput `json:"snapshots,omitempty"`

	// Number of features associated with the collection.
	FeaturesCount *int64 `json:"features_count,omitempty"`

	// Number of properties associated with the collection.
	PropertiesCount *int64 `json:"properties_count,omitempty"`

	// Number of snapshot associated with the collection.
	SnapshotCount *int64 `json:"snapshot_count,omitempty"`
}

// NewCollection : Instantiate Collection (Generic Model Constructor)
func (*AppConfigurationV1) NewCollection(name string, collectionID string) (_model *Collection, err error) {
	_model = &Collection{
		Name: core.StringPtr(name),
		CollectionID: core.StringPtr(collectionID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCollection unmarshals an instance of Collection from the specified map of raw messages.
func UnmarshalCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Collection)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeatureOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "features-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "snapshots", &obj.Snapshots, UnmarshalSnapshotOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "snapshots-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "features_count", &obj.FeaturesCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "features_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "properties_count", &obj.PropertiesCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "snapshot_count", &obj.SnapshotCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "snapshot_count-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CollectionList : List of all Collections.
type CollectionList struct {
	// Array of collections.
	Collections []Collection `json:"collections" validate:"required"`

	// The number of records that are retrieved in a list.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of records that are skipped in a list.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// URL to navigate to the first page of records.
	First *PaginatedListFirst `json:"first" validate:"required"`

	// URL to navigate to the previous list of records.
	Previous *PaginatedListPrevious `json:"previous,omitempty"`

	// URL to navigate to the next list of records.
	Next *PaginatedListNext `json:"next,omitempty"`

	// URL to navigate to the last page of records.
	Last *PaginatedListLast `json:"last" validate:"required"`
}

// UnmarshalCollectionList unmarshals an instance of CollectionList from the specified map of raw messages.
func UnmarshalCollectionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectionList)
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollection)
	if err != nil {
		err = core.SDKErrorf(err, "", "collections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedListPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedListLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *CollectionList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if offset == nil {
		return nil, nil
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		err = core.SDKErrorf(err, "", "parse-int-query-error", common.GetComponentInfo())
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// CollectionLite : Details of the collection.
type CollectionLite struct {
	// Collection name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Collection Id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Creation time of the collection.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last updated time of the collection data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Collection URL.
	Href *string `json:"href,omitempty"`
}

// UnmarshalCollectionLite unmarshals an instance of CollectionLite from the specified map of raw messages.
func UnmarshalCollectionLite(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectionLite)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// CollectionRef : CollectionRef struct
type CollectionRef struct {
	// Collection id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Name of the collection.
	Name *string `json:"name,omitempty"`
}

// NewCollectionRef : Instantiate CollectionRef (Generic Model Constructor)
func (*AppConfigurationV1) NewCollectionRef(collectionID string) (_model *CollectionRef, err error) {
	_model = &CollectionRef{
		CollectionID: core.StringPtr(collectionID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalCollectionRef unmarshals an instance of CollectionRef from the specified map of raw messages.
func UnmarshalCollectionRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectionRef)
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection_id-error", common.GetComponentInfo())
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

// ConfigAction : ConfigAction struct
// Models which "extend" this model:
// - ConfigActionGitConfigPromote
// - ConfigActionGitConfigRestore
type ConfigAction struct {
	// Git commit id will be given as part of the response upon successful git operation.
	GitCommitID *string `json:"git_commit_id,omitempty"`

	// Git commit message.
	GitCommitMessage *string `json:"git_commit_message,omitempty"`

	// Latest time when the snapshot was synced to git.
	LastSyncTime *strfmt.DateTime `json:"last_sync_time,omitempty"`

	// The environments array will contain the environment data and it will also contains properties array and features
	// array that belongs to that environment.
	Environments []ImportEnvironmentSchema `json:"environments,omitempty"`

	// Segments that belongs to the features or properties.
	Segments []ImportSegmentSchema `json:"segments,omitempty"`
}
func (*ConfigAction) isaConfigAction() bool {
	return true
}

type ConfigActionIntf interface {
	isaConfigAction() bool
}

// UnmarshalConfigAction unmarshals an instance of ConfigAction from the specified map of raw messages.
func UnmarshalConfigAction(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigAction)
	err = core.UnmarshalPrimitive(m, "git_commit_id", &obj.GitCommitID)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_commit_message", &obj.GitCommitMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_sync_time", &obj.LastSyncTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_sync_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalImportEnvironmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segments", &obj.Segments, UnmarshalImportSegmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "segments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateCollectionOptions : The CreateCollection options.
type CreateCollectionOptions struct {
	// Collection name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Collection Id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateCollectionOptions : Instantiate CreateCollectionOptions
func (*AppConfigurationV1) NewCreateCollectionOptions(name string, collectionID string) *CreateCollectionOptions {
	return &CreateCollectionOptions{
		Name: core.StringPtr(name),
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetName : Allow user to set Name
func (_options *CreateCollectionOptions) SetName(name string) *CreateCollectionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetCollectionID : Allow user to set CollectionID
func (_options *CreateCollectionOptions) SetCollectionID(collectionID string) *CreateCollectionOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateCollectionOptions) SetDescription(description string) *CreateCollectionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateCollectionOptions) SetTags(tags string) *CreateCollectionOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCollectionOptions) SetHeaders(param map[string]string) *CreateCollectionOptions {
	options.Headers = param
	return options
}

// CreateEnvironmentOptions : The CreateEnvironment options.
type CreateEnvironmentOptions struct {
	// Environment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Environment id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Environment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.
	ColorCode *string `json:"color_code,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateEnvironmentOptions : Instantiate CreateEnvironmentOptions
func (*AppConfigurationV1) NewCreateEnvironmentOptions(name string, environmentID string) *CreateEnvironmentOptions {
	return &CreateEnvironmentOptions{
		Name: core.StringPtr(name),
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetName : Allow user to set Name
func (_options *CreateEnvironmentOptions) SetName(name string) *CreateEnvironmentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *CreateEnvironmentOptions) SetEnvironmentID(environmentID string) *CreateEnvironmentOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateEnvironmentOptions) SetDescription(description string) *CreateEnvironmentOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateEnvironmentOptions) SetTags(tags string) *CreateEnvironmentOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetColorCode : Allow user to set ColorCode
func (_options *CreateEnvironmentOptions) SetColorCode(colorCode string) *CreateEnvironmentOptions {
	_options.ColorCode = core.StringPtr(colorCode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateEnvironmentOptions) SetHeaders(param map[string]string) *CreateEnvironmentOptions {
	options.Headers = param
	return options
}

// CreateFeatureOptions : The CreateFeature options.
type CreateFeatureOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Feature name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Feature id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Type of the feature (BOOLEAN, STRING, NUMERIC). If `type` is `STRING`, then `format` attribute is required.
	Type *string `json:"type" validate:"required"`

	// Value of the feature when it is enabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	EnabledValue interface{} `json:"enabled_value" validate:"required"`

	// Value of the feature when it is disabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	DisabledValue interface{} `json:"disabled_value" validate:"required"`

	// Feature description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
	// for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if
	// the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.
	Format *string `json:"format,omitempty"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Rollout percentage associated with feature flag. Supported only for Lite and Enterprise plans.
	RolloutPercentage *int64 `json:"rollout_percentage,omitempty"`

	// Tags associated with the feature, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different feature flag values for different segments.
	SegmentRules []FeatureSegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified feature flag.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateFeatureOptions.Type property.
// Type of the feature (BOOLEAN, STRING, NUMERIC). If `type` is `STRING`, then `format` attribute is required.
const (
	CreateFeatureOptions_Type_Boolean = "BOOLEAN"
	CreateFeatureOptions_Type_Numeric = "NUMERIC"
	CreateFeatureOptions_Type_String = "STRING"
)

// Constants associated with the CreateFeatureOptions.Format property.
// Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
// for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if
// the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.
const (
	CreateFeatureOptions_Format_JSON = "JSON"
	CreateFeatureOptions_Format_Text = "TEXT"
	CreateFeatureOptions_Format_Yaml = "YAML"
)

// NewCreateFeatureOptions : Instantiate CreateFeatureOptions
func (*AppConfigurationV1) NewCreateFeatureOptions(environmentID string, name string, featureID string, typeVar string, enabledValue interface{}, disabledValue interface{}) *CreateFeatureOptions {
	return &CreateFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		Name: core.StringPtr(name),
		FeatureID: core.StringPtr(featureID),
		Type: core.StringPtr(typeVar),
		EnabledValue: enabledValue,
		DisabledValue: disabledValue,
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *CreateFeatureOptions) SetEnvironmentID(environmentID string) *CreateFeatureOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateFeatureOptions) SetName(name string) *CreateFeatureOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetFeatureID : Allow user to set FeatureID
func (_options *CreateFeatureOptions) SetFeatureID(featureID string) *CreateFeatureOptions {
	_options.FeatureID = core.StringPtr(featureID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreateFeatureOptions) SetType(typeVar string) *CreateFeatureOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetEnabledValue : Allow user to set EnabledValue
func (_options *CreateFeatureOptions) SetEnabledValue(enabledValue interface{}) *CreateFeatureOptions {
	_options.EnabledValue = enabledValue
	return _options
}

// SetDisabledValue : Allow user to set DisabledValue
func (_options *CreateFeatureOptions) SetDisabledValue(disabledValue interface{}) *CreateFeatureOptions {
	_options.DisabledValue = disabledValue
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateFeatureOptions) SetDescription(description string) *CreateFeatureOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *CreateFeatureOptions) SetFormat(format string) *CreateFeatureOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *CreateFeatureOptions) SetEnabled(enabled bool) *CreateFeatureOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetRolloutPercentage : Allow user to set RolloutPercentage
func (_options *CreateFeatureOptions) SetRolloutPercentage(rolloutPercentage int64) *CreateFeatureOptions {
	_options.RolloutPercentage = core.Int64Ptr(rolloutPercentage)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateFeatureOptions) SetTags(tags string) *CreateFeatureOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetSegmentRules : Allow user to set SegmentRules
func (_options *CreateFeatureOptions) SetSegmentRules(segmentRules []FeatureSegmentRule) *CreateFeatureOptions {
	_options.SegmentRules = segmentRules
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *CreateFeatureOptions) SetCollections(collections []CollectionRef) *CreateFeatureOptions {
	_options.Collections = collections
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateFeatureOptions) SetHeaders(param map[string]string) *CreateFeatureOptions {
	options.Headers = param
	return options
}

// CreateGitConfigResponse : Details of the created Git config.
type CreateGitConfigResponse struct {
	// Git config name.
	GitConfigName *string `json:"git_config_name" validate:"required"`

	// Git config Id.
	GitConfigID *string `json:"git_config_id" validate:"required"`

	// Collection Id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Git url which will be used to connect to the github account.
	GitURL *string `json:"git_url" validate:"required"`

	// Branch name to which you need to write or update the configuration.
	GitBranch *string `json:"git_branch" validate:"required"`

	// Git file path, this is a path where your configuration file will be written.
	GitFilePath *string `json:"git_file_path" validate:"required"`

	// Creation time of the git config.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the git config data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Git config URL.
	Href *string `json:"href,omitempty"`
}

// UnmarshalCreateGitConfigResponse unmarshals an instance of CreateGitConfigResponse from the specified map of raw messages.
func UnmarshalCreateGitConfigResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateGitConfigResponse)
	err = core.UnmarshalPrimitive(m, "git_config_name", &obj.GitConfigName)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_config_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_config_id", &obj.GitConfigID)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_config_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_url", &obj.GitURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_branch", &obj.GitBranch)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_branch-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_file_path", &obj.GitFilePath)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_file_path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// CreateGitconfigOptions : The CreateGitconfig options.
type CreateGitconfigOptions struct {
	// Git config name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	GitConfigName *string `json:"git_config_name" validate:"required"`

	// Git config id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	GitConfigID *string `json:"git_config_id" validate:"required"`

	// Collection Id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Git url which will be used to connect to the github account. The url must be formed in this format,
	// https://api.github.com/repos/{owner}/{repo_name} for the personal git account. If you are using the organization
	// account then url must be in this format https://github.{organization_name}.com/api/v3/repos/{owner}/{repo_name} .
	// Note do not provide /(slash) in the beginning or at the end of the url.
	GitURL *string `json:"git_url" validate:"required"`

	// Branch name to which you need to write or update the configuration. Just provide the branch name, do not provide any
	// /(slashes) in the beginning or at the end of the branch name. Note make sure branch exists in your repository.
	GitBranch *string `json:"git_branch" validate:"required"`

	// Git file path, this is a path where your configuration file will be written. The path must contain the file name
	// with `json` extension. We only create or update `json` extension file. Note do not provide any /(slashes) in the
	// beginning or at the end of the file path.
	GitFilePath *string `json:"git_file_path" validate:"required"`

	// Git token, this needs to be provided with enough permission to write and update the file.
	GitToken *string `json:"git_token" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateGitconfigOptions : Instantiate CreateGitconfigOptions
func (*AppConfigurationV1) NewCreateGitconfigOptions(gitConfigName string, gitConfigID string, collectionID string, environmentID string, gitURL string, gitBranch string, gitFilePath string, gitToken string) *CreateGitconfigOptions {
	return &CreateGitconfigOptions{
		GitConfigName: core.StringPtr(gitConfigName),
		GitConfigID: core.StringPtr(gitConfigID),
		CollectionID: core.StringPtr(collectionID),
		EnvironmentID: core.StringPtr(environmentID),
		GitURL: core.StringPtr(gitURL),
		GitBranch: core.StringPtr(gitBranch),
		GitFilePath: core.StringPtr(gitFilePath),
		GitToken: core.StringPtr(gitToken),
	}
}

// SetGitConfigName : Allow user to set GitConfigName
func (_options *CreateGitconfigOptions) SetGitConfigName(gitConfigName string) *CreateGitconfigOptions {
	_options.GitConfigName = core.StringPtr(gitConfigName)
	return _options
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *CreateGitconfigOptions) SetGitConfigID(gitConfigID string) *CreateGitconfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetCollectionID : Allow user to set CollectionID
func (_options *CreateGitconfigOptions) SetCollectionID(collectionID string) *CreateGitconfigOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *CreateGitconfigOptions) SetEnvironmentID(environmentID string) *CreateGitconfigOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetGitURL : Allow user to set GitURL
func (_options *CreateGitconfigOptions) SetGitURL(gitURL string) *CreateGitconfigOptions {
	_options.GitURL = core.StringPtr(gitURL)
	return _options
}

// SetGitBranch : Allow user to set GitBranch
func (_options *CreateGitconfigOptions) SetGitBranch(gitBranch string) *CreateGitconfigOptions {
	_options.GitBranch = core.StringPtr(gitBranch)
	return _options
}

// SetGitFilePath : Allow user to set GitFilePath
func (_options *CreateGitconfigOptions) SetGitFilePath(gitFilePath string) *CreateGitconfigOptions {
	_options.GitFilePath = core.StringPtr(gitFilePath)
	return _options
}

// SetGitToken : Allow user to set GitToken
func (_options *CreateGitconfigOptions) SetGitToken(gitToken string) *CreateGitconfigOptions {
	_options.GitToken = core.StringPtr(gitToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateGitconfigOptions) SetHeaders(param map[string]string) *CreateGitconfigOptions {
	options.Headers = param
	return options
}

// CreatePropertyOptions : The CreateProperty options.
type CreatePropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Property id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	PropertyID *string `json:"property_id" validate:"required"`

	// Type of the property (BOOLEAN, STRING, NUMERIC, SECRETREF). If `type` is `STRING`, then `format` attribute is
	// required.
	Type *string `json:"type" validate:"required"`

	// Value of the Property. The value can be Boolean, Numeric, SecretRef, String - TEXT, String - JSON, String - YAML as
	// per the `type` and `format` attributes.
	Value interface{} `json:"value" validate:"required"`

	// Property description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Format of the property (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
	// for `BOOLEAN`, `NUMERIC` or `SECRETREF` types. This attribute is populated in the response body of `POST, PUT and
	// GET` calls if the type `STRING` is used and not populated for `BOOLEAN`, `NUMERIC` and `SECRETREF` types.
	Format *string `json:"format,omitempty"`

	// Tags associated with the property, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreatePropertyOptions.Type property.
// Type of the property (BOOLEAN, STRING, NUMERIC, SECRETREF). If `type` is `STRING`, then `format` attribute is
// required.
const (
	CreatePropertyOptions_Type_Boolean = "BOOLEAN"
	CreatePropertyOptions_Type_Numeric = "NUMERIC"
	CreatePropertyOptions_Type_Secretref = "SECRETREF"
	CreatePropertyOptions_Type_String = "STRING"
)

// Constants associated with the CreatePropertyOptions.Format property.
// Format of the property (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
// for `BOOLEAN`, `NUMERIC` or `SECRETREF` types. This attribute is populated in the response body of `POST, PUT and
// GET` calls if the type `STRING` is used and not populated for `BOOLEAN`, `NUMERIC` and `SECRETREF` types.
const (
	CreatePropertyOptions_Format_JSON = "JSON"
	CreatePropertyOptions_Format_Text = "TEXT"
	CreatePropertyOptions_Format_Yaml = "YAML"
)

// NewCreatePropertyOptions : Instantiate CreatePropertyOptions
func (*AppConfigurationV1) NewCreatePropertyOptions(environmentID string, name string, propertyID string, typeVar string, value interface{}) *CreatePropertyOptions {
	return &CreatePropertyOptions{
		EnvironmentID: core.StringPtr(environmentID),
		Name: core.StringPtr(name),
		PropertyID: core.StringPtr(propertyID),
		Type: core.StringPtr(typeVar),
		Value: value,
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *CreatePropertyOptions) SetEnvironmentID(environmentID string) *CreatePropertyOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreatePropertyOptions) SetName(name string) *CreatePropertyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetPropertyID : Allow user to set PropertyID
func (_options *CreatePropertyOptions) SetPropertyID(propertyID string) *CreatePropertyOptions {
	_options.PropertyID = core.StringPtr(propertyID)
	return _options
}

// SetType : Allow user to set Type
func (_options *CreatePropertyOptions) SetType(typeVar string) *CreatePropertyOptions {
	_options.Type = core.StringPtr(typeVar)
	return _options
}

// SetValue : Allow user to set Value
func (_options *CreatePropertyOptions) SetValue(value interface{}) *CreatePropertyOptions {
	_options.Value = value
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreatePropertyOptions) SetDescription(description string) *CreatePropertyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetFormat : Allow user to set Format
func (_options *CreatePropertyOptions) SetFormat(format string) *CreatePropertyOptions {
	_options.Format = core.StringPtr(format)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreatePropertyOptions) SetTags(tags string) *CreatePropertyOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetSegmentRules : Allow user to set SegmentRules
func (_options *CreatePropertyOptions) SetSegmentRules(segmentRules []SegmentRule) *CreatePropertyOptions {
	_options.SegmentRules = segmentRules
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *CreatePropertyOptions) SetCollections(collections []CollectionRef) *CreatePropertyOptions {
	_options.Collections = collections
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePropertyOptions) SetHeaders(param map[string]string) *CreatePropertyOptions {
	options.Headers = param
	return options
}

// CreateSegmentOptions : The CreateSegment options.
type CreateSegmentOptions struct {
	// Segment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Segment id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	SegmentID *string `json:"segment_id" validate:"required"`

	// List of rules that determine if the entity belongs to the segment during feature / property evaluation. An entity is
	// identified by an unique identifier and the attributes that it defines. Any feature flag and property value
	// evaluation is performed in the context of an entity when it is targeted to segments.
	Rules []Rule `json:"rules" validate:"required"`

	// Segment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the segments, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateSegmentOptions : Instantiate CreateSegmentOptions
func (*AppConfigurationV1) NewCreateSegmentOptions(name string, segmentID string, rules []Rule) *CreateSegmentOptions {
	return &CreateSegmentOptions{
		Name: core.StringPtr(name),
		SegmentID: core.StringPtr(segmentID),
		Rules: rules,
	}
}

// SetName : Allow user to set Name
func (_options *CreateSegmentOptions) SetName(name string) *CreateSegmentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetSegmentID : Allow user to set SegmentID
func (_options *CreateSegmentOptions) SetSegmentID(segmentID string) *CreateSegmentOptions {
	_options.SegmentID = core.StringPtr(segmentID)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *CreateSegmentOptions) SetRules(rules []Rule) *CreateSegmentOptions {
	_options.Rules = rules
	return _options
}

// SetDescription : Allow user to set Description
func (_options *CreateSegmentOptions) SetDescription(description string) *CreateSegmentOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *CreateSegmentOptions) SetTags(tags string) *CreateSegmentOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSegmentOptions) SetHeaders(param map[string]string) *CreateSegmentOptions {
	options.Headers = param
	return options
}

// CreateWorkflowConfig : CreateWorkflowConfig struct
// Models which "extend" this model:
// - CreateWorkflowConfigExternalServiceNow
// - CreateWorkflowConfigIBMServiceNow
type CreateWorkflowConfig struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url,omitempty"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn,omitempty"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type,omitempty"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn,omitempty"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id,omitempty"`
}
func (*CreateWorkflowConfig) isaCreateWorkflowConfig() bool {
	return true
}

type CreateWorkflowConfigIntf interface {
	isaCreateWorkflowConfig() bool
}

// UnmarshalCreateWorkflowConfig unmarshals an instance of CreateWorkflowConfig from the specified map of raw messages.
func UnmarshalCreateWorkflowConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateWorkflowConfig)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateWorkflowconfigOptions : The CreateWorkflowconfig options.
type CreateWorkflowconfigOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// The request body to create a new workflow config.
	WorkflowConfig CreateWorkflowConfigIntf `json:"WorkflowConfig" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateWorkflowconfigOptions : Instantiate CreateWorkflowconfigOptions
func (*AppConfigurationV1) NewCreateWorkflowconfigOptions(environmentID string, workflowConfig CreateWorkflowConfigIntf) *CreateWorkflowconfigOptions {
	return &CreateWorkflowconfigOptions{
		EnvironmentID: core.StringPtr(environmentID),
		WorkflowConfig: workflowConfig,
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *CreateWorkflowconfigOptions) SetEnvironmentID(environmentID string) *CreateWorkflowconfigOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetWorkflowConfig : Allow user to set WorkflowConfig
func (_options *CreateWorkflowconfigOptions) SetWorkflowConfig(workflowConfig CreateWorkflowConfigIntf) *CreateWorkflowconfigOptions {
	_options.WorkflowConfig = workflowConfig
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateWorkflowconfigOptions) SetHeaders(param map[string]string) *CreateWorkflowconfigOptions {
	options.Headers = param
	return options
}

// CreateWorkflowconfigResponse : CreateWorkflowconfigResponse struct
// Models which "extend" this model:
// - CreateWorkflowconfigResponseExternalServiceNow
// - CreateWorkflowconfigResponseIBMServiceNow
type CreateWorkflowconfigResponse struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url,omitempty"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn,omitempty"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type,omitempty"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn,omitempty"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id,omitempty"`
}
func (*CreateWorkflowconfigResponse) isaCreateWorkflowconfigResponse() bool {
	return true
}

type CreateWorkflowconfigResponseIntf interface {
	isaCreateWorkflowconfigResponse() bool
}

// UnmarshalCreateWorkflowconfigResponse unmarshals an instance of CreateWorkflowconfigResponse from the specified map of raw messages.
func UnmarshalCreateWorkflowconfigResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateWorkflowconfigResponse)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteCollectionOptions : The DeleteCollection options.
type DeleteCollectionOptions struct {
	// Collection Id of the collection.
	CollectionID *string `json:"collection_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteCollectionOptions : Instantiate DeleteCollectionOptions
func (*AppConfigurationV1) NewDeleteCollectionOptions(collectionID string) *DeleteCollectionOptions {
	return &DeleteCollectionOptions{
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetCollectionID : Allow user to set CollectionID
func (_options *DeleteCollectionOptions) SetCollectionID(collectionID string) *DeleteCollectionOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteCollectionOptions) SetHeaders(param map[string]string) *DeleteCollectionOptions {
	options.Headers = param
	return options
}

// DeleteEnvironmentOptions : The DeleteEnvironment options.
type DeleteEnvironmentOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteEnvironmentOptions : Instantiate DeleteEnvironmentOptions
func (*AppConfigurationV1) NewDeleteEnvironmentOptions(environmentID string) *DeleteEnvironmentOptions {
	return &DeleteEnvironmentOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *DeleteEnvironmentOptions) SetEnvironmentID(environmentID string) *DeleteEnvironmentOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEnvironmentOptions) SetHeaders(param map[string]string) *DeleteEnvironmentOptions {
	options.Headers = param
	return options
}

// DeleteFeatureOptions : The DeleteFeature options.
type DeleteFeatureOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Feature Id.
	FeatureID *string `json:"feature_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteFeatureOptions : Instantiate DeleteFeatureOptions
func (*AppConfigurationV1) NewDeleteFeatureOptions(environmentID string, featureID string) *DeleteFeatureOptions {
	return &DeleteFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *DeleteFeatureOptions) SetEnvironmentID(environmentID string) *DeleteFeatureOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetFeatureID : Allow user to set FeatureID
func (_options *DeleteFeatureOptions) SetFeatureID(featureID string) *DeleteFeatureOptions {
	_options.FeatureID = core.StringPtr(featureID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteFeatureOptions) SetHeaders(param map[string]string) *DeleteFeatureOptions {
	options.Headers = param
	return options
}

// DeleteGitconfigOptions : The DeleteGitconfig options.
type DeleteGitconfigOptions struct {
	// Git Config Id.
	GitConfigID *string `json:"git_config_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteGitconfigOptions : Instantiate DeleteGitconfigOptions
func (*AppConfigurationV1) NewDeleteGitconfigOptions(gitConfigID string) *DeleteGitconfigOptions {
	return &DeleteGitconfigOptions{
		GitConfigID: core.StringPtr(gitConfigID),
	}
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *DeleteGitconfigOptions) SetGitConfigID(gitConfigID string) *DeleteGitconfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteGitconfigOptions) SetHeaders(param map[string]string) *DeleteGitconfigOptions {
	options.Headers = param
	return options
}

// DeletePropertyOptions : The DeleteProperty options.
type DeletePropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeletePropertyOptions : Instantiate DeletePropertyOptions
func (*AppConfigurationV1) NewDeletePropertyOptions(environmentID string, propertyID string) *DeletePropertyOptions {
	return &DeletePropertyOptions{
		EnvironmentID: core.StringPtr(environmentID),
		PropertyID: core.StringPtr(propertyID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *DeletePropertyOptions) SetEnvironmentID(environmentID string) *DeletePropertyOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetPropertyID : Allow user to set PropertyID
func (_options *DeletePropertyOptions) SetPropertyID(propertyID string) *DeletePropertyOptions {
	_options.PropertyID = core.StringPtr(propertyID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeletePropertyOptions) SetHeaders(param map[string]string) *DeletePropertyOptions {
	options.Headers = param
	return options
}

// DeleteSegmentOptions : The DeleteSegment options.
type DeleteSegmentOptions struct {
	// Segment Id.
	SegmentID *string `json:"segment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteSegmentOptions : Instantiate DeleteSegmentOptions
func (*AppConfigurationV1) NewDeleteSegmentOptions(segmentID string) *DeleteSegmentOptions {
	return &DeleteSegmentOptions{
		SegmentID: core.StringPtr(segmentID),
	}
}

// SetSegmentID : Allow user to set SegmentID
func (_options *DeleteSegmentOptions) SetSegmentID(segmentID string) *DeleteSegmentOptions {
	_options.SegmentID = core.StringPtr(segmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSegmentOptions) SetHeaders(param map[string]string) *DeleteSegmentOptions {
	options.Headers = param
	return options
}

// DeleteWorkflowconfigOptions : The DeleteWorkflowconfig options.
type DeleteWorkflowconfigOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteWorkflowconfigOptions : Instantiate DeleteWorkflowconfigOptions
func (*AppConfigurationV1) NewDeleteWorkflowconfigOptions(environmentID string) *DeleteWorkflowconfigOptions {
	return &DeleteWorkflowconfigOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *DeleteWorkflowconfigOptions) SetEnvironmentID(environmentID string) *DeleteWorkflowconfigOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteWorkflowconfigOptions) SetHeaders(param map[string]string) *DeleteWorkflowconfigOptions {
	options.Headers = param
	return options
}

// Environment : Details of the environment.
type Environment struct {
	// Environment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Environment id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Environment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.
	ColorCode *string `json:"color_code,omitempty"`

	// Creation time of the environment.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the environment data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Environment URL.
	Href *string `json:"href,omitempty"`

	// List of Features associated with the environment.
	Features []FeatureOutput `json:"features,omitempty"`

	// List of properties associated with the environment.
	Properties []PropertyOutput `json:"properties,omitempty"`

	// List of snapshots associated with the environment.
	Snapshots []SnapshotOutput `json:"snapshots,omitempty"`
}

// NewEnvironment : Instantiate Environment (Generic Model Constructor)
func (*AppConfigurationV1) NewEnvironment(name string, environmentID string) (_model *Environment, err error) {
	_model = &Environment{
		Name: core.StringPtr(name),
		EnvironmentID: core.StringPtr(environmentID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalEnvironment unmarshals an instance of Environment from the specified map of raw messages.
func UnmarshalEnvironment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Environment)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "color_code", &obj.ColorCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "color_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeatureOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "features-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "snapshots", &obj.Snapshots, UnmarshalSnapshotOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "snapshots-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentList : List of all environments.
type EnvironmentList struct {
	// Array of environments.
	Environments []Environment `json:"environments" validate:"required"`

	// The number of records that are retrieved in a list.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of records that are skipped in a list.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// URL to navigate to the first page of records.
	First *PaginatedListFirst `json:"first" validate:"required"`

	// URL to navigate to the previous list of records.
	Previous *PaginatedListPrevious `json:"previous,omitempty"`

	// URL to navigate to the next list of records.
	Next *PaginatedListNext `json:"next,omitempty"`

	// URL to navigate to the last page of records.
	Last *PaginatedListLast `json:"last" validate:"required"`
}

// UnmarshalEnvironmentList unmarshals an instance of EnvironmentList from the specified map of raw messages.
func UnmarshalEnvironmentList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentList)
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalEnvironment)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedListPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedListLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *EnvironmentList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if offset == nil {
		return nil, nil
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		err = core.SDKErrorf(err, "", "parse-int-query-error", common.GetComponentInfo())
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// ExternalServiceNowCredentials : The credentials of the External ServiceNow instance.
type ExternalServiceNowCredentials struct {
	// ServiceNow instance login username.
	Username *string `json:"username" validate:"required"`

	// ServiceNow instance login password.
	Password *string `json:"password" validate:"required"`

	// The auto-generated unique ID of the application in your ServiceNow instance.
	ClientID *string `json:"client_id" validate:"required"`

	// The secret string that both the ServiceNow instance and the client application use to authorize communications with
	// one another.
	ClientSecret *string `json:"client_secret" validate:"required"`
}

// NewExternalServiceNowCredentials : Instantiate ExternalServiceNowCredentials (Generic Model Constructor)
func (*AppConfigurationV1) NewExternalServiceNowCredentials(username string, password string, clientID string, clientSecret string) (_model *ExternalServiceNowCredentials, err error) {
	_model = &ExternalServiceNowCredentials{
		Username: core.StringPtr(username),
		Password: core.StringPtr(password),
		ClientID: core.StringPtr(clientID),
		ClientSecret: core.StringPtr(clientSecret),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalExternalServiceNowCredentials unmarshals an instance of ExternalServiceNowCredentials from the specified map of raw messages.
func UnmarshalExternalServiceNowCredentials(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExternalServiceNowCredentials)
	err = core.UnmarshalPrimitive(m, "username", &obj.Username)
	if err != nil {
		err = core.SDKErrorf(err, "", "username-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "password", &obj.Password)
	if err != nil {
		err = core.SDKErrorf(err, "", "password-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "client_id", &obj.ClientID)
	if err != nil {
		err = core.SDKErrorf(err, "", "client_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "client_secret", &obj.ClientSecret)
	if err != nil {
		err = core.SDKErrorf(err, "", "client_secret-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Feature : Details of the feature.
type Feature struct {
	// Feature name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Feature id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Feature description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Type of the feature (BOOLEAN, STRING, NUMERIC). If `type` is `STRING`, then `format` attribute is required.
	Type *string `json:"type" validate:"required"`

	// Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
	// for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if
	// the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.
	Format *string `json:"format,omitempty"`

	// Value of the feature when it is enabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	EnabledValue interface{} `json:"enabled_value" validate:"required"`

	// Value of the feature when it is disabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	DisabledValue interface{} `json:"disabled_value" validate:"required"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Rollout percentage associated with feature flag. Supported only for Lite and Enterprise plans.
	RolloutPercentage *int64 `json:"rollout_percentage,omitempty"`

	// Tags associated with the feature, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different feature flag values for different segments.
	SegmentRules []FeatureSegmentRule `json:"segment_rules,omitempty"`

	// Denotes if the targeting rules are specified for the feature flag.
	SegmentExists *bool `json:"segment_exists,omitempty"`

	// List of collection id representing the collections that are associated with the specified feature flag.
	Collections []CollectionRef `json:"collections,omitempty"`

	// If you have enabled the workflow configuration and have a pending CR then this provides the change_request_number.
	ChangeRequestNumber *string `json:"change_request_number,omitempty"`

	// If you have enabled the workflow configuration and have a pending CR then this provides the change_request_status.
	ChangeRequestStatus *string `json:"change_request_status,omitempty"`

	// Creation time of the feature flag.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the feature flag data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// The last occurrence of the feature flag value evaluation.
	EvaluationTime *strfmt.DateTime `json:"evaluation_time,omitempty"`

	// Feature flag URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the Feature.Type property.
// Type of the feature (BOOLEAN, STRING, NUMERIC). If `type` is `STRING`, then `format` attribute is required.
const (
	Feature_Type_Boolean = "BOOLEAN"
	Feature_Type_Numeric = "NUMERIC"
	Feature_Type_String = "STRING"
)

// Constants associated with the Feature.Format property.
// Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
// for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if
// the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.
const (
	Feature_Format_JSON = "JSON"
	Feature_Format_Text = "TEXT"
	Feature_Format_Yaml = "YAML"
)

// UnmarshalFeature unmarshals an instance of Feature from the specified map of raw messages.
func UnmarshalFeature(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Feature)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "feature_id", &obj.FeatureID)
	if err != nil {
		err = core.SDKErrorf(err, "", "feature_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		err = core.SDKErrorf(err, "", "format-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled_value", &obj.EnabledValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled_value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled_value", &obj.DisabledValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled_value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rollout_percentage", &obj.RolloutPercentage)
	if err != nil {
		err = core.SDKErrorf(err, "", "rollout_percentage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segment_rules", &obj.SegmentRules, UnmarshalFeatureSegmentRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_exists", &obj.SegmentExists)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_exists-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollectionRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "collections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "change_request_number", &obj.ChangeRequestNumber)
	if err != nil {
		err = core.SDKErrorf(err, "", "change_request_number-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "change_request_status", &obj.ChangeRequestStatus)
	if err != nil {
		err = core.SDKErrorf(err, "", "change_request_status-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluation_time", &obj.EvaluationTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "evaluation_time-error", common.GetComponentInfo())
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

// FeatureOutput : Feature object.
type FeatureOutput struct {
	// Feature id.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Feature name.
	Name *string `json:"name" validate:"required"`
}

// NewFeatureOutput : Instantiate FeatureOutput (Generic Model Constructor)
func (*AppConfigurationV1) NewFeatureOutput(featureID string, name string) (_model *FeatureOutput, err error) {
	_model = &FeatureOutput{
		FeatureID: core.StringPtr(featureID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalFeatureOutput unmarshals an instance of FeatureOutput from the specified map of raw messages.
func UnmarshalFeatureOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FeatureOutput)
	err = core.UnmarshalPrimitive(m, "feature_id", &obj.FeatureID)
	if err != nil {
		err = core.SDKErrorf(err, "", "feature_id-error", common.GetComponentInfo())
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

// FeatureSegmentRule : FeatureSegmentRule struct
type FeatureSegmentRule struct {
	// The list of targeted segments.
	Rules []TargetSegments `json:"rules" validate:"required"`

	// Value to be used for evaluation for this rule. The value can be Boolean, SecretRef, String - TEXT , String - JSON ,
	// String - YAML or a Numeric value as per the `type` and `format` attributes.
	Value interface{} `json:"value" validate:"required"`

	// Order of the rule, used during evaluation. The evaluation is performed in the order defined and the value associated
	// with the first matching rule is used for evaluation.
	Order *int64 `json:"order" validate:"required"`

	// Rollout percentage associated with feature flag. Supported only for Lite and Enterprise plans.
	RolloutPercentage *int64 `json:"rollout_percentage,omitempty"`
}

// NewFeatureSegmentRule : Instantiate FeatureSegmentRule (Generic Model Constructor)
func (*AppConfigurationV1) NewFeatureSegmentRule(rules []TargetSegments, value interface{}, order int64) (_model *FeatureSegmentRule, err error) {
	_model = &FeatureSegmentRule{
		Rules: rules,
		Value: value,
		Order: core.Int64Ptr(order),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalFeatureSegmentRule unmarshals an instance of FeatureSegmentRule from the specified map of raw messages.
func UnmarshalFeatureSegmentRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FeatureSegmentRule)
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalTargetSegments)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "order", &obj.Order)
	if err != nil {
		err = core.SDKErrorf(err, "", "order-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rollout_percentage", &obj.RolloutPercentage)
	if err != nil {
		err = core.SDKErrorf(err, "", "rollout_percentage-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// FeaturesList : List of all features.
type FeaturesList struct {
	// Array of Features.
	Features []Feature `json:"features" validate:"required"`

	// The number of records that are retrieved in a list.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of records that are skipped in a list.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// URL to navigate to the first page of records.
	First *PaginatedListFirst `json:"first" validate:"required"`

	// URL to navigate to the previous list of records.
	Previous *PaginatedListPrevious `json:"previous,omitempty"`

	// URL to navigate to the next list of records.
	Next *PaginatedListNext `json:"next,omitempty"`

	// URL to navigate to the last page of records.
	Last *PaginatedListLast `json:"last" validate:"required"`
}

// UnmarshalFeaturesList unmarshals an instance of FeaturesList from the specified map of raw messages.
func UnmarshalFeaturesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FeaturesList)
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeature)
	if err != nil {
		err = core.SDKErrorf(err, "", "features-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedListPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedListLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *FeaturesList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if offset == nil {
		return nil, nil
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		err = core.SDKErrorf(err, "", "parse-int-query-error", common.GetComponentInfo())
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// GetCollectionOptions : The GetCollection options.
type GetCollectionOptions struct {
	// Collection Id of the collection.
	CollectionID *string `json:"collection_id" validate:"required,ne="`

	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Include feature, property, snapshots details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetCollectionOptions.Include property.
const (
	GetCollectionOptions_Include_Features = "features"
	GetCollectionOptions_Include_Properties = "properties"
	GetCollectionOptions_Include_Snapshots = "snapshots"
)

// NewGetCollectionOptions : Instantiate GetCollectionOptions
func (*AppConfigurationV1) NewGetCollectionOptions(collectionID string) *GetCollectionOptions {
	return &GetCollectionOptions{
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetCollectionID : Allow user to set CollectionID
func (_options *GetCollectionOptions) SetCollectionID(collectionID string) *GetCollectionOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetExpand : Allow user to set Expand
func (_options *GetCollectionOptions) SetExpand(expand bool) *GetCollectionOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetCollectionOptions) SetInclude(include []string) *GetCollectionOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCollectionOptions) SetHeaders(param map[string]string) *GetCollectionOptions {
	options.Headers = param
	return options
}

// GetEnvironmentOptions : The GetEnvironment options.
type GetEnvironmentOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Include feature, property, snapshots details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetEnvironmentOptions.Include property.
const (
	GetEnvironmentOptions_Include_Features = "features"
	GetEnvironmentOptions_Include_Properties = "properties"
	GetEnvironmentOptions_Include_Snapshots = "snapshots"
)

// NewGetEnvironmentOptions : Instantiate GetEnvironmentOptions
func (*AppConfigurationV1) NewGetEnvironmentOptions(environmentID string) *GetEnvironmentOptions {
	return &GetEnvironmentOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *GetEnvironmentOptions) SetEnvironmentID(environmentID string) *GetEnvironmentOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetExpand : Allow user to set Expand
func (_options *GetEnvironmentOptions) SetExpand(expand bool) *GetEnvironmentOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetEnvironmentOptions) SetInclude(include []string) *GetEnvironmentOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetEnvironmentOptions) SetHeaders(param map[string]string) *GetEnvironmentOptions {
	options.Headers = param
	return options
}

// GetFeatureOptions : The GetFeature options.
type GetFeatureOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Feature Id.
	FeatureID *string `json:"feature_id" validate:"required,ne="`

	// Include the associated collections or targeting rules or change request details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetFeatureOptions.Include property.
const (
	GetFeatureOptions_Include_ChangeRequest = "change_request"
	GetFeatureOptions_Include_Collections = "collections"
	GetFeatureOptions_Include_Rules = "rules"
)

// NewGetFeatureOptions : Instantiate GetFeatureOptions
func (*AppConfigurationV1) NewGetFeatureOptions(environmentID string, featureID string) *GetFeatureOptions {
	return &GetFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *GetFeatureOptions) SetEnvironmentID(environmentID string) *GetFeatureOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetFeatureID : Allow user to set FeatureID
func (_options *GetFeatureOptions) SetFeatureID(featureID string) *GetFeatureOptions {
	_options.FeatureID = core.StringPtr(featureID)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetFeatureOptions) SetInclude(include []string) *GetFeatureOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetFeatureOptions) SetHeaders(param map[string]string) *GetFeatureOptions {
	options.Headers = param
	return options
}

// GetGitconfigOptions : The GetGitconfig options.
type GetGitconfigOptions struct {
	// Git Config Id.
	GitConfigID *string `json:"git_config_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetGitconfigOptions : Instantiate GetGitconfigOptions
func (*AppConfigurationV1) NewGetGitconfigOptions(gitConfigID string) *GetGitconfigOptions {
	return &GetGitconfigOptions{
		GitConfigID: core.StringPtr(gitConfigID),
	}
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *GetGitconfigOptions) SetGitConfigID(gitConfigID string) *GetGitconfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetGitconfigOptions) SetHeaders(param map[string]string) *GetGitconfigOptions {
	options.Headers = param
	return options
}

// GetPropertyOptions : The GetProperty options.
type GetPropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Include the associated collections or targeting rules details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetPropertyOptions.Include property.
const (
	GetPropertyOptions_Include_Collections = "collections"
	GetPropertyOptions_Include_Rules = "rules"
)

// NewGetPropertyOptions : Instantiate GetPropertyOptions
func (*AppConfigurationV1) NewGetPropertyOptions(environmentID string, propertyID string) *GetPropertyOptions {
	return &GetPropertyOptions{
		EnvironmentID: core.StringPtr(environmentID),
		PropertyID: core.StringPtr(propertyID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *GetPropertyOptions) SetEnvironmentID(environmentID string) *GetPropertyOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetPropertyID : Allow user to set PropertyID
func (_options *GetPropertyOptions) SetPropertyID(propertyID string) *GetPropertyOptions {
	_options.PropertyID = core.StringPtr(propertyID)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetPropertyOptions) SetInclude(include []string) *GetPropertyOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetPropertyOptions) SetHeaders(param map[string]string) *GetPropertyOptions {
	options.Headers = param
	return options
}

// GetSegmentOptions : The GetSegment options.
type GetSegmentOptions struct {
	// Segment Id.
	SegmentID *string `json:"segment_id" validate:"required,ne="`

	// Include feature and property details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the GetSegmentOptions.Include property.
const (
	GetSegmentOptions_Include_Features = "features"
	GetSegmentOptions_Include_Properties = "properties"
)

// NewGetSegmentOptions : Instantiate GetSegmentOptions
func (*AppConfigurationV1) NewGetSegmentOptions(segmentID string) *GetSegmentOptions {
	return &GetSegmentOptions{
		SegmentID: core.StringPtr(segmentID),
	}
}

// SetSegmentID : Allow user to set SegmentID
func (_options *GetSegmentOptions) SetSegmentID(segmentID string) *GetSegmentOptions {
	_options.SegmentID = core.StringPtr(segmentID)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *GetSegmentOptions) SetInclude(include []string) *GetSegmentOptions {
	_options.Include = include
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetSegmentOptions) SetHeaders(param map[string]string) *GetSegmentOptions {
	options.Headers = param
	return options
}

// GitConfig : Details of the Git Config.
type GitConfig struct {
	// Git config name.
	GitConfigName *string `json:"git_config_name" validate:"required"`

	// Git config id.
	GitConfigID *string `json:"git_config_id" validate:"required"`

	// Details of the collection.
	Collection *GitConfigCollection `json:"collection" validate:"required"`

	// Details of the environment.
	Environment *GitConfigEnvironment `json:"environment" validate:"required"`

	// Git url which will be used to connect to the github account.
	GitURL *string `json:"git_url" validate:"required"`

	// Branch name to which you need to write or update the configuration.
	GitBranch *string `json:"git_branch" validate:"required"`

	// Git file path, this is a path where your configuration file will be written.
	GitFilePath *string `json:"git_file_path" validate:"required"`

	// Latest time when the snapshot was synced to git.
	LastSyncTime *strfmt.DateTime `json:"last_sync_time,omitempty"`

	// Creation time of the git config.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the git config data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Git config URL.
	Href *string `json:"href,omitempty"`
}

// UnmarshalGitConfig unmarshals an instance of GitConfig from the specified map of raw messages.
func UnmarshalGitConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GitConfig)
	err = core.UnmarshalPrimitive(m, "git_config_name", &obj.GitConfigName)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_config_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_config_id", &obj.GitConfigID)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_config_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "collection", &obj.Collection, UnmarshalGitConfigCollection)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "environment", &obj.Environment, UnmarshalGitConfigEnvironment)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_url", &obj.GitURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_branch", &obj.GitBranch)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_branch-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_file_path", &obj.GitFilePath)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_file_path-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_sync_time", &obj.LastSyncTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_sync_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// GitConfigCollection : Details of the collection.
type GitConfigCollection struct {
	// Collection name.
	Name *string `json:"name,omitempty"`

	// Collection Id.
	CollectionID *string `json:"collection_id,omitempty"`
}

// UnmarshalGitConfigCollection unmarshals an instance of GitConfigCollection from the specified map of raw messages.
func UnmarshalGitConfigCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GitConfigCollection)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GitConfigEnvironment : Details of the environment.
type GitConfigEnvironment struct {
	// Environment name.
	Name *string `json:"name,omitempty"`

	// Environment Id.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Environment color code.
	ColorCode *string `json:"color_code,omitempty"`
}

// UnmarshalGitConfigEnvironment unmarshals an instance of GitConfigEnvironment from the specified map of raw messages.
func UnmarshalGitConfigEnvironment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GitConfigEnvironment)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "color_code", &obj.ColorCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "color_code-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GitConfigList : List of all Git Configs.
type GitConfigList struct {
	// Array of Git Configs.
	GitConfig []GitConfig `json:"git_config" validate:"required"`

	// The number of records that are retrieved in a list.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of records that are skipped in a list.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// URL to navigate to the first page of records.
	First *PaginatedListFirst `json:"first" validate:"required"`

	// URL to navigate to the previous list of records.
	Previous *PaginatedListPrevious `json:"previous,omitempty"`

	// URL to navigate to the next list of records.
	Next *PaginatedListNext `json:"next,omitempty"`

	// URL to navigate to the last page of records.
	Last *PaginatedListLast `json:"last" validate:"required"`
}

// UnmarshalGitConfigList unmarshals an instance of GitConfigList from the specified map of raw messages.
func UnmarshalGitConfigList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GitConfigList)
	err = core.UnmarshalModel(m, "git_config", &obj.GitConfig, UnmarshalGitConfig)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_config-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedListPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedListLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *GitConfigList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if offset == nil {
		return nil, nil
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		err = core.SDKErrorf(err, "", "parse-int-query-error", common.GetComponentInfo())
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// GitConfigPromote : Details of the promote operation.
type GitConfigPromote struct {
	// Git commit id will be given as part of the response upon successful git operation.
	GitCommitID *string `json:"git_commit_id" validate:"required"`

	// Git commit message.
	GitCommitMessage *string `json:"git_commit_message" validate:"required"`

	// Latest time when the snapshot was synced to git.
	LastSyncTime *strfmt.DateTime `json:"last_sync_time,omitempty"`
}

// UnmarshalGitConfigPromote unmarshals an instance of GitConfigPromote from the specified map of raw messages.
func UnmarshalGitConfigPromote(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GitConfigPromote)
	err = core.UnmarshalPrimitive(m, "git_commit_id", &obj.GitCommitID)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_commit_message", &obj.GitCommitMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_sync_time", &obj.LastSyncTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_sync_time-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GitConfigRestore : Details of the restore operation.
type GitConfigRestore struct {
	// The environments array will contain the environment data and it will also contains properties array and features
	// array that belongs to that environment.
	Environments []ImportEnvironmentSchema `json:"environments" validate:"required"`

	// Segments that belongs to the features or properties.
	Segments []ImportSegmentSchema `json:"segments" validate:"required"`
}

// UnmarshalGitConfigRestore unmarshals an instance of GitConfigRestore from the specified map of raw messages.
func UnmarshalGitConfigRestore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(GitConfigRestore)
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalImportEnvironmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segments", &obj.Segments, UnmarshalImportSegmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "segments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportCollectionSchema : Collection to be created.
type ImportCollectionSchema struct {
	// Collection id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Description of the collection, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`
}

// NewImportCollectionSchema : Instantiate ImportCollectionSchema (Generic Model Constructor)
func (*AppConfigurationV1) NewImportCollectionSchema(collectionID string, name string) (_model *ImportCollectionSchema, err error) {
	_model = &ImportCollectionSchema{
		CollectionID: core.StringPtr(collectionID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalImportCollectionSchema unmarshals an instance of ImportCollectionSchema from the specified map of raw messages.
func UnmarshalImportCollectionSchema(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportCollectionSchema)
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		err = core.SDKErrorf(err, "", "collection_id-error", common.GetComponentInfo())
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
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportConfig : Full instance configuration.
type ImportConfig struct {
	// Array will contain features and properties per environment.
	Environments []ImportEnvironmentSchema `json:"environments,omitempty"`

	// Array will contain collections details.
	Collections []ImportCollectionSchema `json:"collections,omitempty"`

	// Array will contain segments details.
	Segments []ImportSegmentSchema `json:"segments,omitempty"`
}

// UnmarshalImportConfig unmarshals an instance of ImportConfig from the specified map of raw messages.
func UnmarshalImportConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportConfig)
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalImportEnvironmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalImportCollectionSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "collections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segments", &obj.Segments, UnmarshalImportSegmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "segments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportConfigOptions : The ImportConfig options.
type ImportConfigOptions struct {
	// Array will contain features and properties per environment.
	Environments []ImportEnvironmentSchema `json:"environments,omitempty"`

	// Array will contain collections details.
	Collections []ImportCollectionSchema `json:"collections,omitempty"`

	// Array will contain segments details.
	Segments []ImportSegmentSchema `json:"segments,omitempty"`

	// Full instance import requires query parameter `clean=true` to perform wiping of the existing data.
	Clean *string `json:"clean,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewImportConfigOptions : Instantiate ImportConfigOptions
func (*AppConfigurationV1) NewImportConfigOptions() *ImportConfigOptions {
	return &ImportConfigOptions{}
}

// SetEnvironments : Allow user to set Environments
func (_options *ImportConfigOptions) SetEnvironments(environments []ImportEnvironmentSchema) *ImportConfigOptions {
	_options.Environments = environments
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *ImportConfigOptions) SetCollections(collections []ImportCollectionSchema) *ImportConfigOptions {
	_options.Collections = collections
	return _options
}

// SetSegments : Allow user to set Segments
func (_options *ImportConfigOptions) SetSegments(segments []ImportSegmentSchema) *ImportConfigOptions {
	_options.Segments = segments
	return _options
}

// SetClean : Allow user to set Clean
func (_options *ImportConfigOptions) SetClean(clean string) *ImportConfigOptions {
	_options.Clean = core.StringPtr(clean)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ImportConfigOptions) SetHeaders(param map[string]string) *ImportConfigOptions {
	options.Headers = param
	return options
}

// ImportEnvironmentSchema : Environment attributes to import.
type ImportEnvironmentSchema struct {
	// Environment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Environment Id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Environment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.
	ColorCode *string `json:"color_code,omitempty"`

	// Array will contain features per environment.
	Features []ImportFeatureRequestBody `json:"features,omitempty"`

	// Array will contain properties per environment.
	Properties []ImportPropertyRequestBody `json:"properties,omitempty"`
}

// NewImportEnvironmentSchema : Instantiate ImportEnvironmentSchema (Generic Model Constructor)
func (*AppConfigurationV1) NewImportEnvironmentSchema(name string, environmentID string) (_model *ImportEnvironmentSchema, err error) {
	_model = &ImportEnvironmentSchema{
		Name: core.StringPtr(name),
		EnvironmentID: core.StringPtr(environmentID),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalImportEnvironmentSchema unmarshals an instance of ImportEnvironmentSchema from the specified map of raw messages.
func UnmarshalImportEnvironmentSchema(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportEnvironmentSchema)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "color_code", &obj.ColorCode)
	if err != nil {
		err = core.SDKErrorf(err, "", "color_code-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalImportFeatureRequestBody)
	if err != nil {
		err = core.SDKErrorf(err, "", "features-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalImportPropertyRequestBody)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportFeatureRequestBody : Requset Body of the feature.
type ImportFeatureRequestBody struct {
	// Feature name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Feature id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Feature description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Type of the feature (BOOLEAN, STRING, NUMERIC). If `type` is `STRING`, then `format` attribute is required.
	Type *string `json:"type" validate:"required"`

	// Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
	// for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if
	// the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.
	Format *string `json:"format,omitempty"`

	// Value of the feature when it is enabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	EnabledValue interface{} `json:"enabled_value" validate:"required"`

	// Value of the feature when it is disabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	DisabledValue interface{} `json:"disabled_value" validate:"required"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Rollout percentage associated with feature flag. Supported only for Lite and Enterprise plans.
	RolloutPercentage *int64 `json:"rollout_percentage,omitempty"`

	// Tags associated with the feature, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different feature flag values for different segments.
	SegmentRules []FeatureSegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified feature flag.
	Collections []CollectionRef `json:"collections,omitempty"`

	// This attribute explains whether the feature has to be imported or not.
	IsOverridden *bool `json:"isOverridden" validate:"required"`
}

// Constants associated with the ImportFeatureRequestBody.Type property.
// Type of the feature (BOOLEAN, STRING, NUMERIC). If `type` is `STRING`, then `format` attribute is required.
const (
	ImportFeatureRequestBody_Type_Boolean = "BOOLEAN"
	ImportFeatureRequestBody_Type_Numeric = "NUMERIC"
	ImportFeatureRequestBody_Type_String = "STRING"
)

// Constants associated with the ImportFeatureRequestBody.Format property.
// Format of the feature (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
// for `BOOLEAN` and `NUMERIC` types. This property is populated in the response body of `POST, PUT and GET` calls if
// the type `STRING` is used and not populated for `BOOLEAN` and `NUMERIC` types.
const (
	ImportFeatureRequestBody_Format_JSON = "JSON"
	ImportFeatureRequestBody_Format_Text = "TEXT"
	ImportFeatureRequestBody_Format_Yaml = "YAML"
)

// NewImportFeatureRequestBody : Instantiate ImportFeatureRequestBody (Generic Model Constructor)
func (*AppConfigurationV1) NewImportFeatureRequestBody(name string, featureID string, typeVar string, enabledValue interface{}, disabledValue interface{}, isOverridden bool) (_model *ImportFeatureRequestBody, err error) {
	_model = &ImportFeatureRequestBody{
		Name: core.StringPtr(name),
		FeatureID: core.StringPtr(featureID),
		Type: core.StringPtr(typeVar),
		EnabledValue: enabledValue,
		DisabledValue: disabledValue,
		IsOverridden: core.BoolPtr(isOverridden),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalImportFeatureRequestBody unmarshals an instance of ImportFeatureRequestBody from the specified map of raw messages.
func UnmarshalImportFeatureRequestBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportFeatureRequestBody)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "feature_id", &obj.FeatureID)
	if err != nil {
		err = core.SDKErrorf(err, "", "feature_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		err = core.SDKErrorf(err, "", "format-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled_value", &obj.EnabledValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled_value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled_value", &obj.DisabledValue)
	if err != nil {
		err = core.SDKErrorf(err, "", "disabled_value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rollout_percentage", &obj.RolloutPercentage)
	if err != nil {
		err = core.SDKErrorf(err, "", "rollout_percentage-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segment_rules", &obj.SegmentRules, UnmarshalFeatureSegmentRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollectionRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "collections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "isOverridden", &obj.IsOverridden)
	if err != nil {
		err = core.SDKErrorf(err, "", "isOverridden-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportPropertyRequestBody : Details of the property.
type ImportPropertyRequestBody struct {
	// Property name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Property id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	PropertyID *string `json:"property_id" validate:"required"`

	// Property description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Type of the property (BOOLEAN, STRING, NUMERIC, SECRETREF). If `type` is `STRING`, then `format` attribute is
	// required.
	Type *string `json:"type" validate:"required"`

	// Format of the property (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
	// for `BOOLEAN`, `NUMERIC` or `SECRETREF` types. This attribute is populated in the response body of `POST, PUT and
	// GET` calls if the type `STRING` is used and not populated for `BOOLEAN`, `NUMERIC` and `SECRETREF` types.
	Format *string `json:"format,omitempty"`

	// Value of the Property. The value can be Boolean, Numeric, SecretRef, String - TEXT, String - JSON, String - YAML as
	// per the `type` and `format` attributes.
	Value interface{} `json:"value" validate:"required"`

	// Tags associated with the property, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// This attribute explains whether the property has to be imported or not.
	IsOverridden *bool `json:"isOverridden" validate:"required"`
}

// Constants associated with the ImportPropertyRequestBody.Type property.
// Type of the property (BOOLEAN, STRING, NUMERIC, SECRETREF). If `type` is `STRING`, then `format` attribute is
// required.
const (
	ImportPropertyRequestBody_Type_Boolean = "BOOLEAN"
	ImportPropertyRequestBody_Type_Numeric = "NUMERIC"
	ImportPropertyRequestBody_Type_Secretref = "SECRETREF"
	ImportPropertyRequestBody_Type_String = "STRING"
)

// Constants associated with the ImportPropertyRequestBody.Format property.
// Format of the property (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
// for `BOOLEAN`, `NUMERIC` or `SECRETREF` types. This attribute is populated in the response body of `POST, PUT and
// GET` calls if the type `STRING` is used and not populated for `BOOLEAN`, `NUMERIC` and `SECRETREF` types.
const (
	ImportPropertyRequestBody_Format_JSON = "JSON"
	ImportPropertyRequestBody_Format_Text = "TEXT"
	ImportPropertyRequestBody_Format_Yaml = "YAML"
)

// NewImportPropertyRequestBody : Instantiate ImportPropertyRequestBody (Generic Model Constructor)
func (*AppConfigurationV1) NewImportPropertyRequestBody(name string, propertyID string, typeVar string, value interface{}, isOverridden bool) (_model *ImportPropertyRequestBody, err error) {
	_model = &ImportPropertyRequestBody{
		Name: core.StringPtr(name),
		PropertyID: core.StringPtr(propertyID),
		Type: core.StringPtr(typeVar),
		Value: value,
		IsOverridden: core.BoolPtr(isOverridden),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalImportPropertyRequestBody unmarshals an instance of ImportPropertyRequestBody from the specified map of raw messages.
func UnmarshalImportPropertyRequestBody(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportPropertyRequestBody)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "property_id", &obj.PropertyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "property_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		err = core.SDKErrorf(err, "", "format-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segment_rules", &obj.SegmentRules, UnmarshalSegmentRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollectionRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "collections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "isOverridden", &obj.IsOverridden)
	if err != nil {
		err = core.SDKErrorf(err, "", "isOverridden-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ImportSegmentSchema : Details of the segment.
type ImportSegmentSchema struct {
	// Segment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Segment id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	SegmentID *string `json:"segment_id" validate:"required"`

	// Segment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the segments, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// List of rules that determine if the entity belongs to the segment during feature / property evaluation. An entity is
	// identified by an unique identifier and the attributes that it defines. Any feature flag and property value
	// evaluation is performed in the context of an entity when it is targeted to segments.
	Rules []Rule `json:"rules" validate:"required"`
}

// NewImportSegmentSchema : Instantiate ImportSegmentSchema (Generic Model Constructor)
func (*AppConfigurationV1) NewImportSegmentSchema(name string, segmentID string, rules []Rule) (_model *ImportSegmentSchema, err error) {
	_model = &ImportSegmentSchema{
		Name: core.StringPtr(name),
		SegmentID: core.StringPtr(segmentID),
		Rules: rules,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalImportSegmentSchema unmarshals an instance of ImportSegmentSchema from the specified map of raw messages.
func UnmarshalImportSegmentSchema(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ImportSegmentSchema)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_id", &obj.SegmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListCollectionsOptions : The ListCollections options.
type ListCollectionsOptions struct {
	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the collection details based on the specified attribute. By default, items are sorted by name.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Filter collections by a list of comma separated features.
	Features []string `json:"features,omitempty"`

	// Filter collections by a list of comma separated properties.
	Properties []string `json:"properties,omitempty"`

	// Include feature, property, snapshots details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Searches for the provided keyword and returns the appropriate row with that value. Here the search happens on the
	// '[Name OR Tag]' of the entity.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListCollectionsOptions.Sort property.
// Sort the collection details based on the specified attribute. By default, items are sorted by name.
const (
	ListCollectionsOptions_Sort_CreatedTime = "created_time"
	ListCollectionsOptions_Sort_ID = "id"
	ListCollectionsOptions_Sort_Name = "name"
	ListCollectionsOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListCollectionsOptions.Include property.
const (
	ListCollectionsOptions_Include_Features = "features"
	ListCollectionsOptions_Include_Properties = "properties"
	ListCollectionsOptions_Include_Snapshots = "snapshots"
)

// NewListCollectionsOptions : Instantiate ListCollectionsOptions
func (*AppConfigurationV1) NewListCollectionsOptions() *ListCollectionsOptions {
	return &ListCollectionsOptions{}
}

// SetExpand : Allow user to set Expand
func (_options *ListCollectionsOptions) SetExpand(expand bool) *ListCollectionsOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListCollectionsOptions) SetSort(sort string) *ListCollectionsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ListCollectionsOptions) SetTags(tags string) *ListCollectionsOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetFeatures : Allow user to set Features
func (_options *ListCollectionsOptions) SetFeatures(features []string) *ListCollectionsOptions {
	_options.Features = features
	return _options
}

// SetProperties : Allow user to set Properties
func (_options *ListCollectionsOptions) SetProperties(properties []string) *ListCollectionsOptions {
	_options.Properties = properties
	return _options
}

// SetInclude : Allow user to set Include
func (_options *ListCollectionsOptions) SetInclude(include []string) *ListCollectionsOptions {
	_options.Include = include
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListCollectionsOptions) SetLimit(limit int64) *ListCollectionsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListCollectionsOptions) SetOffset(offset int64) *ListCollectionsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListCollectionsOptions) SetSearch(search string) *ListCollectionsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListCollectionsOptions) SetHeaders(param map[string]string) *ListCollectionsOptions {
	options.Headers = param
	return options
}

// ListEnvironmentsOptions : The ListEnvironments options.
type ListEnvironmentsOptions struct {
	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the environment details based on the specified attribute. By default, items are sorted by name.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Include feature, property, snapshots details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Searches for the provided keyword and returns the appropriate row with that value. Here the search happens on the
	// '[Name OR Tag]' of the entity.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListEnvironmentsOptions.Sort property.
// Sort the environment details based on the specified attribute. By default, items are sorted by name.
const (
	ListEnvironmentsOptions_Sort_CreatedTime = "created_time"
	ListEnvironmentsOptions_Sort_ID = "id"
	ListEnvironmentsOptions_Sort_Name = "name"
	ListEnvironmentsOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListEnvironmentsOptions.Include property.
const (
	ListEnvironmentsOptions_Include_Features = "features"
	ListEnvironmentsOptions_Include_Properties = "properties"
	ListEnvironmentsOptions_Include_Snapshots = "snapshots"
)

// NewListEnvironmentsOptions : Instantiate ListEnvironmentsOptions
func (*AppConfigurationV1) NewListEnvironmentsOptions() *ListEnvironmentsOptions {
	return &ListEnvironmentsOptions{}
}

// SetExpand : Allow user to set Expand
func (_options *ListEnvironmentsOptions) SetExpand(expand bool) *ListEnvironmentsOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListEnvironmentsOptions) SetSort(sort string) *ListEnvironmentsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ListEnvironmentsOptions) SetTags(tags string) *ListEnvironmentsOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *ListEnvironmentsOptions) SetInclude(include []string) *ListEnvironmentsOptions {
	_options.Include = include
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListEnvironmentsOptions) SetLimit(limit int64) *ListEnvironmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListEnvironmentsOptions) SetOffset(offset int64) *ListEnvironmentsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListEnvironmentsOptions) SetSearch(search string) *ListEnvironmentsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListEnvironmentsOptions) SetHeaders(param map[string]string) *ListEnvironmentsOptions {
	options.Headers = param
	return options
}

// ListFeaturesOptions : The ListFeatures options.
type ListFeaturesOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the feature details based on the specified attribute. By default, items are sorted by name.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Filter features by a list of comma separated collections.
	Collections []string `json:"collections,omitempty"`

	// Filter features by a list of comma separated segments.
	Segments []string `json:"segments,omitempty"`

	// Include the associated collections or targeting rules or change request details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Searches for the provided keyword and returns the appropriate row with that value. Here the search happens on the
	// '[Name OR Tag]' of the entity.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListFeaturesOptions.Sort property.
// Sort the feature details based on the specified attribute. By default, items are sorted by name.
const (
	ListFeaturesOptions_Sort_CreatedTime = "created_time"
	ListFeaturesOptions_Sort_ID = "id"
	ListFeaturesOptions_Sort_Name = "name"
	ListFeaturesOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListFeaturesOptions.Include property.
const (
	ListFeaturesOptions_Include_ChangeRequest = "change_request"
	ListFeaturesOptions_Include_Collections = "collections"
	ListFeaturesOptions_Include_Rules = "rules"
)

// NewListFeaturesOptions : Instantiate ListFeaturesOptions
func (*AppConfigurationV1) NewListFeaturesOptions(environmentID string) *ListFeaturesOptions {
	return &ListFeaturesOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *ListFeaturesOptions) SetEnvironmentID(environmentID string) *ListFeaturesOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetExpand : Allow user to set Expand
func (_options *ListFeaturesOptions) SetExpand(expand bool) *ListFeaturesOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListFeaturesOptions) SetSort(sort string) *ListFeaturesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ListFeaturesOptions) SetTags(tags string) *ListFeaturesOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *ListFeaturesOptions) SetCollections(collections []string) *ListFeaturesOptions {
	_options.Collections = collections
	return _options
}

// SetSegments : Allow user to set Segments
func (_options *ListFeaturesOptions) SetSegments(segments []string) *ListFeaturesOptions {
	_options.Segments = segments
	return _options
}

// SetInclude : Allow user to set Include
func (_options *ListFeaturesOptions) SetInclude(include []string) *ListFeaturesOptions {
	_options.Include = include
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListFeaturesOptions) SetLimit(limit int64) *ListFeaturesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListFeaturesOptions) SetOffset(offset int64) *ListFeaturesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListFeaturesOptions) SetSearch(search string) *ListFeaturesOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListFeaturesOptions) SetHeaders(param map[string]string) *ListFeaturesOptions {
	options.Headers = param
	return options
}

// ListInstanceConfigOptions : The ListInstanceConfig options.
type ListInstanceConfigOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListInstanceConfigOptions : Instantiate ListInstanceConfigOptions
func (*AppConfigurationV1) NewListInstanceConfigOptions() *ListInstanceConfigOptions {
	return &ListInstanceConfigOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListInstanceConfigOptions) SetHeaders(param map[string]string) *ListInstanceConfigOptions {
	options.Headers = param
	return options
}

// ListOriginconfigsOptions : The ListOriginconfigs options.
type ListOriginconfigsOptions struct {

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListOriginconfigsOptions : Instantiate ListOriginconfigsOptions
func (*AppConfigurationV1) NewListOriginconfigsOptions() *ListOriginconfigsOptions {
	return &ListOriginconfigsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListOriginconfigsOptions) SetHeaders(param map[string]string) *ListOriginconfigsOptions {
	options.Headers = param
	return options
}

// ListPropertiesOptions : The ListProperties options.
type ListPropertiesOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the property details based on the specified attribute. By default, items are sorted by name.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Filter properties by a list of comma separated collections.
	Collections []string `json:"collections,omitempty"`

	// Filter properties by a list of comma separated segments.
	Segments []string `json:"segments,omitempty"`

	// Include the associated collections or targeting rules details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Searches for the provided keyword and returns the appropriate row with that value. Here the search happens on the
	// '[Name OR Tag]' of the entity.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListPropertiesOptions.Sort property.
// Sort the property details based on the specified attribute. By default, items are sorted by name.
const (
	ListPropertiesOptions_Sort_CreatedTime = "created_time"
	ListPropertiesOptions_Sort_ID = "id"
	ListPropertiesOptions_Sort_Name = "name"
	ListPropertiesOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListPropertiesOptions.Include property.
const (
	ListPropertiesOptions_Include_Collections = "collections"
	ListPropertiesOptions_Include_Rules = "rules"
)

// NewListPropertiesOptions : Instantiate ListPropertiesOptions
func (*AppConfigurationV1) NewListPropertiesOptions(environmentID string) *ListPropertiesOptions {
	return &ListPropertiesOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *ListPropertiesOptions) SetEnvironmentID(environmentID string) *ListPropertiesOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetExpand : Allow user to set Expand
func (_options *ListPropertiesOptions) SetExpand(expand bool) *ListPropertiesOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListPropertiesOptions) SetSort(sort string) *ListPropertiesOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ListPropertiesOptions) SetTags(tags string) *ListPropertiesOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *ListPropertiesOptions) SetCollections(collections []string) *ListPropertiesOptions {
	_options.Collections = collections
	return _options
}

// SetSegments : Allow user to set Segments
func (_options *ListPropertiesOptions) SetSegments(segments []string) *ListPropertiesOptions {
	_options.Segments = segments
	return _options
}

// SetInclude : Allow user to set Include
func (_options *ListPropertiesOptions) SetInclude(include []string) *ListPropertiesOptions {
	_options.Include = include
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListPropertiesOptions) SetLimit(limit int64) *ListPropertiesOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListPropertiesOptions) SetOffset(offset int64) *ListPropertiesOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListPropertiesOptions) SetSearch(search string) *ListPropertiesOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListPropertiesOptions) SetHeaders(param map[string]string) *ListPropertiesOptions {
	options.Headers = param
	return options
}

// ListSegmentsOptions : The ListSegments options.
type ListSegmentsOptions struct {
	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the segment details based on the specified attribute. By default, items are sorted by name.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Segment details to include the associated rules in the response.
	Include *string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Searches for the provided keyword and returns the appropriate row with that value. Here the search happens on the
	// '[Name OR Tag]' of the entity.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListSegmentsOptions.Sort property.
// Sort the segment details based on the specified attribute. By default, items are sorted by name.
const (
	ListSegmentsOptions_Sort_CreatedTime = "created_time"
	ListSegmentsOptions_Sort_ID = "id"
	ListSegmentsOptions_Sort_Name = "name"
	ListSegmentsOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListSegmentsOptions.Include property.
// Segment details to include the associated rules in the response.
const (
	ListSegmentsOptions_Include_Rules = "rules"
)

// NewListSegmentsOptions : Instantiate ListSegmentsOptions
func (*AppConfigurationV1) NewListSegmentsOptions() *ListSegmentsOptions {
	return &ListSegmentsOptions{}
}

// SetExpand : Allow user to set Expand
func (_options *ListSegmentsOptions) SetExpand(expand bool) *ListSegmentsOptions {
	_options.Expand = core.BoolPtr(expand)
	return _options
}

// SetSort : Allow user to set Sort
func (_options *ListSegmentsOptions) SetSort(sort string) *ListSegmentsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *ListSegmentsOptions) SetTags(tags string) *ListSegmentsOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetInclude : Allow user to set Include
func (_options *ListSegmentsOptions) SetInclude(include string) *ListSegmentsOptions {
	_options.Include = core.StringPtr(include)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSegmentsOptions) SetLimit(limit int64) *ListSegmentsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListSegmentsOptions) SetOffset(offset int64) *ListSegmentsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListSegmentsOptions) SetSearch(search string) *ListSegmentsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSegmentsOptions) SetHeaders(param map[string]string) *ListSegmentsOptions {
	options.Headers = param
	return options
}

// ListSnapshotsOptions : The ListSnapshots options.
type ListSnapshotsOptions struct {
	// Sort the git configurations details based on the specified attribute. By default, items are sorted by name.
	Sort *string `json:"sort,omitempty"`

	// Filters the response based on the specified collection_id.
	CollectionID *string `json:"collection_id,omitempty"`

	// Filters the response based on the specified environment_id.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Searches for the provided keyword and returns the appropriate row with that value. Here the search happens on the
	// '[Name]' of the entity.
	Search *string `json:"search,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the ListSnapshotsOptions.Sort property.
// Sort the git configurations details based on the specified attribute. By default, items are sorted by name.
const (
	ListSnapshotsOptions_Sort_CreatedTime = "created_time"
	ListSnapshotsOptions_Sort_ID = "id"
	ListSnapshotsOptions_Sort_Name = "name"
	ListSnapshotsOptions_Sort_UpdatedTime = "updated_time"
)

// NewListSnapshotsOptions : Instantiate ListSnapshotsOptions
func (*AppConfigurationV1) NewListSnapshotsOptions() *ListSnapshotsOptions {
	return &ListSnapshotsOptions{}
}

// SetSort : Allow user to set Sort
func (_options *ListSnapshotsOptions) SetSort(sort string) *ListSnapshotsOptions {
	_options.Sort = core.StringPtr(sort)
	return _options
}

// SetCollectionID : Allow user to set CollectionID
func (_options *ListSnapshotsOptions) SetCollectionID(collectionID string) *ListSnapshotsOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *ListSnapshotsOptions) SetEnvironmentID(environmentID string) *ListSnapshotsOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListSnapshotsOptions) SetLimit(limit int64) *ListSnapshotsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListSnapshotsOptions) SetOffset(offset int64) *ListSnapshotsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetSearch : Allow user to set Search
func (_options *ListSnapshotsOptions) SetSearch(search string) *ListSnapshotsOptions {
	_options.Search = core.StringPtr(search)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListSnapshotsOptions) SetHeaders(param map[string]string) *ListSnapshotsOptions {
	options.Headers = param
	return options
}

// ListWorkflowconfigOptions : The ListWorkflowconfig options.
type ListWorkflowconfigOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListWorkflowconfigOptions : Instantiate ListWorkflowconfigOptions
func (*AppConfigurationV1) NewListWorkflowconfigOptions(environmentID string) *ListWorkflowconfigOptions {
	return &ListWorkflowconfigOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *ListWorkflowconfigOptions) SetEnvironmentID(environmentID string) *ListWorkflowconfigOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListWorkflowconfigOptions) SetHeaders(param map[string]string) *ListWorkflowconfigOptions {
	options.Headers = param
	return options
}

// ListWorkflowconfigResponse : ListWorkflowconfigResponse struct
// Models which "extend" this model:
// - ListWorkflowconfigResponseExternalServiceNow
// - ListWorkflowconfigResponseIBMServiceNow
type ListWorkflowconfigResponse struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url,omitempty"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn,omitempty"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type,omitempty"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn,omitempty"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id,omitempty"`
}
func (*ListWorkflowconfigResponse) isaListWorkflowconfigResponse() bool {
	return true
}

type ListWorkflowconfigResponseIntf interface {
	isaListWorkflowconfigResponse() bool
}

// UnmarshalListWorkflowconfigResponse unmarshals an instance of ListWorkflowconfigResponse from the specified map of raw messages.
func UnmarshalListWorkflowconfigResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListWorkflowconfigResponse)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// OriginConfigList : List of all origin configs.
type OriginConfigList struct {
	// List of allowed origins. Specify the parameter as a list of comma separated origins.
	AllowedOrigins []string `json:"allowed_origins" validate:"required"`

	// Creation time of the origin configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the origin configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Origin Config URL.
	Href *string `json:"href,omitempty"`
}

// UnmarshalOriginConfigList unmarshals an instance of OriginConfigList from the specified map of raw messages.
func UnmarshalOriginConfigList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(OriginConfigList)
	err = core.UnmarshalPrimitive(m, "allowed_origins", &obj.AllowedOrigins)
	if err != nil {
		err = core.SDKErrorf(err, "", "allowed_origins-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// PaginatedListFirst : URL to navigate to the first page of records.
type PaginatedListFirst struct {
	// URL to the page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedListFirst unmarshals an instance of PaginatedListFirst from the specified map of raw messages.
func UnmarshalPaginatedListFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedListFirst)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedListLast : URL to navigate to the last page of records.
type PaginatedListLast struct {
	// URL to the page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedListLast unmarshals an instance of PaginatedListLast from the specified map of raw messages.
func UnmarshalPaginatedListLast(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedListLast)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedListNext : URL to navigate to the next list of records.
type PaginatedListNext struct {
	// URL to the page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedListNext unmarshals an instance of PaginatedListNext from the specified map of raw messages.
func UnmarshalPaginatedListNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedListNext)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PaginatedListPrevious : URL to navigate to the previous list of records.
type PaginatedListPrevious struct {
	// URL to the page.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalPaginatedListPrevious unmarshals an instance of PaginatedListPrevious from the specified map of raw messages.
func UnmarshalPaginatedListPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PaginatedListPrevious)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PromoteGitconfigOptions : The PromoteGitconfig options.
type PromoteGitconfigOptions struct {
	// Git Config Id.
	GitConfigID *string `json:"git_config_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewPromoteGitconfigOptions : Instantiate PromoteGitconfigOptions
func (*AppConfigurationV1) NewPromoteGitconfigOptions(gitConfigID string) *PromoteGitconfigOptions {
	return &PromoteGitconfigOptions{
		GitConfigID: core.StringPtr(gitConfigID),
	}
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *PromoteGitconfigOptions) SetGitConfigID(gitConfigID string) *PromoteGitconfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PromoteGitconfigOptions) SetHeaders(param map[string]string) *PromoteGitconfigOptions {
	options.Headers = param
	return options
}

// PromoteRestoreConfigOptions : The PromoteRestoreConfig options.
type PromoteRestoreConfigOptions struct {
	// Git Config Id.
	GitConfigID *string `json:"git_config_id" validate:"required"`

	// Promote configuration to Git or Restore configuration from Git.
	Action *string `json:"action" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the PromoteRestoreConfigOptions.Action property.
// Promote configuration to Git or Restore configuration from Git.
const (
	PromoteRestoreConfigOptions_Action_Promote = "promote"
	PromoteRestoreConfigOptions_Action_Restore = "restore"
)

// NewPromoteRestoreConfigOptions : Instantiate PromoteRestoreConfigOptions
func (*AppConfigurationV1) NewPromoteRestoreConfigOptions(gitConfigID string, action string) *PromoteRestoreConfigOptions {
	return &PromoteRestoreConfigOptions{
		GitConfigID: core.StringPtr(gitConfigID),
		Action: core.StringPtr(action),
	}
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *PromoteRestoreConfigOptions) SetGitConfigID(gitConfigID string) *PromoteRestoreConfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetAction : Allow user to set Action
func (_options *PromoteRestoreConfigOptions) SetAction(action string) *PromoteRestoreConfigOptions {
	_options.Action = core.StringPtr(action)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *PromoteRestoreConfigOptions) SetHeaders(param map[string]string) *PromoteRestoreConfigOptions {
	options.Headers = param
	return options
}

// PropertiesList : List of all properties.
type PropertiesList struct {
	// Array of properties.
	Properties []Property `json:"properties" validate:"required"`

	// The number of records that are retrieved in a list.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of records that are skipped in a list.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// URL to navigate to the first page of records.
	First *PaginatedListFirst `json:"first" validate:"required"`

	// URL to navigate to the previous list of records.
	Previous *PaginatedListPrevious `json:"previous,omitempty"`

	// URL to navigate to the next list of records.
	Next *PaginatedListNext `json:"next,omitempty"`

	// URL to navigate to the last page of records.
	Last *PaginatedListLast `json:"last" validate:"required"`
}

// UnmarshalPropertiesList unmarshals an instance of PropertiesList from the specified map of raw messages.
func UnmarshalPropertiesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PropertiesList)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedListPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedListLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *PropertiesList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if offset == nil {
		return nil, nil
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		err = core.SDKErrorf(err, "", "parse-int-query-error", common.GetComponentInfo())
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// Property : Details of the property.
type Property struct {
	// Property name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Property id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	PropertyID *string `json:"property_id" validate:"required"`

	// Property description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Type of the property (BOOLEAN, STRING, NUMERIC, SECRETREF). If `type` is `STRING`, then `format` attribute is
	// required.
	Type *string `json:"type" validate:"required"`

	// Format of the property (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
	// for `BOOLEAN`, `NUMERIC` or `SECRETREF` types. This attribute is populated in the response body of `POST, PUT and
	// GET` calls if the type `STRING` is used and not populated for `BOOLEAN`, `NUMERIC` and `SECRETREF` types.
	Format *string `json:"format,omitempty"`

	// Value of the Property. The value can be Boolean, Numeric, SecretRef, String - TEXT, String - JSON, String - YAML as
	// per the `type` and `format` attributes.
	Value interface{} `json:"value" validate:"required"`

	// Tags associated with the property, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// Denotes if the targeting rules are specified for the property.
	SegmentExists *bool `json:"segment_exists,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Creation time of the property.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the property data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// The last occurrence of the property value evaluation.
	EvaluationTime *strfmt.DateTime `json:"evaluation_time,omitempty"`

	// Property URL.
	Href *string `json:"href,omitempty"`
}

// Constants associated with the Property.Type property.
// Type of the property (BOOLEAN, STRING, NUMERIC, SECRETREF). If `type` is `STRING`, then `format` attribute is
// required.
const (
	Property_Type_Boolean = "BOOLEAN"
	Property_Type_Numeric = "NUMERIC"
	Property_Type_Secretref = "SECRETREF"
	Property_Type_String = "STRING"
)

// Constants associated with the Property.Format property.
// Format of the property (TEXT, JSON, YAML) and it is a required attribute when `type` is `STRING`. It is not required
// for `BOOLEAN`, `NUMERIC` or `SECRETREF` types. This attribute is populated in the response body of `POST, PUT and
// GET` calls if the type `STRING` is used and not populated for `BOOLEAN`, `NUMERIC` and `SECRETREF` types.
const (
	Property_Format_JSON = "JSON"
	Property_Format_Text = "TEXT"
	Property_Format_Yaml = "YAML"
)

// NewProperty : Instantiate Property (Generic Model Constructor)
func (*AppConfigurationV1) NewProperty(name string, propertyID string, typeVar string, value interface{}) (_model *Property, err error) {
	_model = &Property{
		Name: core.StringPtr(name),
		PropertyID: core.StringPtr(propertyID),
		Type: core.StringPtr(typeVar),
		Value: value,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalProperty unmarshals an instance of Property from the specified map of raw messages.
func UnmarshalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Property)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "property_id", &obj.PropertyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "property_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		err = core.SDKErrorf(err, "", "type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "format", &obj.Format)
	if err != nil {
		err = core.SDKErrorf(err, "", "format-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segment_rules", &obj.SegmentRules, UnmarshalSegmentRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_exists", &obj.SegmentExists)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_exists-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollectionRef)
	if err != nil {
		err = core.SDKErrorf(err, "", "collections-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluation_time", &obj.EvaluationTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "evaluation_time-error", common.GetComponentInfo())
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

// PropertyOutput : Property object.
type PropertyOutput struct {
	// Property id.
	PropertyID *string `json:"property_id" validate:"required"`

	// Property name.
	Name *string `json:"name" validate:"required"`
}

// NewPropertyOutput : Instantiate PropertyOutput (Generic Model Constructor)
func (*AppConfigurationV1) NewPropertyOutput(propertyID string, name string) (_model *PropertyOutput, err error) {
	_model = &PropertyOutput{
		PropertyID: core.StringPtr(propertyID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalPropertyOutput unmarshals an instance of PropertyOutput from the specified map of raw messages.
func UnmarshalPropertyOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PropertyOutput)
	err = core.UnmarshalPrimitive(m, "property_id", &obj.PropertyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "property_id-error", common.GetComponentInfo())
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

// RestoreGitconfigOptions : The RestoreGitconfig options.
type RestoreGitconfigOptions struct {
	// Git Config Id.
	GitConfigID *string `json:"git_config_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewRestoreGitconfigOptions : Instantiate RestoreGitconfigOptions
func (*AppConfigurationV1) NewRestoreGitconfigOptions(gitConfigID string) *RestoreGitconfigOptions {
	return &RestoreGitconfigOptions{
		GitConfigID: core.StringPtr(gitConfigID),
	}
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *RestoreGitconfigOptions) SetGitConfigID(gitConfigID string) *RestoreGitconfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *RestoreGitconfigOptions) SetHeaders(param map[string]string) *RestoreGitconfigOptions {
	options.Headers = param
	return options
}

// Rule : Rule is used to determine if the entity belongs to the segment during feature / property evaluation.
type Rule struct {
	// Attribute name.
	AttributeName *string `json:"attribute_name" validate:"required"`

	// Operator to be used for the evaluation if the entity belongs to the segment.
	Operator *string `json:"operator" validate:"required"`

	// List of values. Entities matching any of the given values will be considered to belong to the segment.
	Values []string `json:"values" validate:"required"`
}

// Constants associated with the Rule.Operator property.
// Operator to be used for the evaluation if the entity belongs to the segment.
const (
	Rule_Operator_Contains = "contains"
	Rule_Operator_Endswith = "endsWith"
	Rule_Operator_Greaterthan = "greaterThan"
	Rule_Operator_Greaterthanequals = "greaterThanEquals"
	Rule_Operator_Is = "is"
	Rule_Operator_Lesserthan = "lesserThan"
	Rule_Operator_Lesserthanequals = "lesserThanEquals"
	Rule_Operator_Startswith = "startsWith"
)

// NewRule : Instantiate Rule (Generic Model Constructor)
func (*AppConfigurationV1) NewRule(attributeName string, operator string, values []string) (_model *Rule, err error) {
	_model = &Rule{
		AttributeName: core.StringPtr(attributeName),
		Operator: core.StringPtr(operator),
		Values: values,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "attribute_name", &obj.AttributeName)
	if err != nil {
		err = core.SDKErrorf(err, "", "attribute_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		err = core.SDKErrorf(err, "", "operator-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "values", &obj.Values)
	if err != nil {
		err = core.SDKErrorf(err, "", "values-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Segment : Details of the segment.
type Segment struct {
	// Segment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name" validate:"required"`

	// Segment id. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	SegmentID *string `json:"segment_id" validate:"required"`

	// Segment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the segments, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// List of rules that determine if the entity belongs to the segment during feature / property evaluation. An entity is
	// identified by an unique identifier and the attributes that it defines. Any feature flag and property value
	// evaluation is performed in the context of an entity when it is targeted to segments.
	Rules []Rule `json:"rules" validate:"required"`

	// Creation time of the segment.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the segment data.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Segment URL.
	Href *string `json:"href,omitempty"`

	// List of Features associated with the segment.
	Features []FeatureOutput `json:"features,omitempty"`

	// List of properties associated with the segment.
	Properties []PropertyOutput `json:"properties,omitempty"`
}

// NewSegment : Instantiate Segment (Generic Model Constructor)
func (*AppConfigurationV1) NewSegment(name string, segmentID string, rules []Rule) (_model *Segment, err error) {
	_model = &Segment{
		Name: core.StringPtr(name),
		SegmentID: core.StringPtr(segmentID),
		Rules: rules,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalSegment unmarshals an instance of Segment from the specified map of raw messages.
func UnmarshalSegment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Segment)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_id", &obj.SegmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "segment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		err = core.SDKErrorf(err, "", "description-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		err = core.SDKErrorf(err, "", "tags-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeatureOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "features-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyOutput)
	if err != nil {
		err = core.SDKErrorf(err, "", "properties-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SegmentRule : SegmentRule struct
type SegmentRule struct {
	// The list of targeted segments.
	Rules []TargetSegments `json:"rules" validate:"required"`

	// Value to be used for evaluation for this rule. The value can be Boolean, SecretRef, String - TEXT , String - JSON ,
	// String - YAML or a Numeric value as per the `type` and `format` attributes.
	Value interface{} `json:"value" validate:"required"`

	// Order of the rule, used during evaluation. The evaluation is performed in the order defined and the value associated
	// with the first matching rule is used for evaluation.
	Order *int64 `json:"order" validate:"required"`
}

// NewSegmentRule : Instantiate SegmentRule (Generic Model Constructor)
func (*AppConfigurationV1) NewSegmentRule(rules []TargetSegments, value interface{}, order int64) (_model *SegmentRule, err error) {
	_model = &SegmentRule{
		Rules: rules,
		Value: value,
		Order: core.Int64Ptr(order),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalSegmentRule unmarshals an instance of SegmentRule from the specified map of raw messages.
func UnmarshalSegmentRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SegmentRule)
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalTargetSegments)
	if err != nil {
		err = core.SDKErrorf(err, "", "rules-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		err = core.SDKErrorf(err, "", "value-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "order", &obj.Order)
	if err != nil {
		err = core.SDKErrorf(err, "", "order-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SegmentsList : List of all segments.
type SegmentsList struct {
	// Array of Segments.
	Segments []Segment `json:"segments" validate:"required"`

	// The number of records that are retrieved in a list.
	Limit *int64 `json:"limit" validate:"required"`

	// The number of records that are skipped in a list.
	Offset *int64 `json:"offset" validate:"required"`

	// The total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// URL to navigate to the first page of records.
	First *PaginatedListFirst `json:"first" validate:"required"`

	// URL to navigate to the previous list of records.
	Previous *PaginatedListPrevious `json:"previous,omitempty"`

	// URL to navigate to the next list of records.
	Next *PaginatedListNext `json:"next,omitempty"`

	// URL to navigate to the last page of records.
	Last *PaginatedListLast `json:"last" validate:"required"`
}

// UnmarshalSegmentsList unmarshals an instance of SegmentsList from the specified map of raw messages.
func UnmarshalSegmentsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SegmentsList)
	err = core.UnmarshalModel(m, "segments", &obj.Segments, UnmarshalSegment)
	if err != nil {
		err = core.SDKErrorf(err, "", "segments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPaginatedListFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPaginatedListPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPaginatedListNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPaginatedListLast)
	if err != nil {
		err = core.SDKErrorf(err, "", "last-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *SegmentsList) GetNextOffset() (*int64, error) {
	if core.IsNil(resp.Next) {
		return nil, nil
	}
	offset, err := core.GetQueryParam(resp.Next.Href, "offset")
	if err != nil {
		err = core.SDKErrorf(err, "", "read-query-param-error", common.GetComponentInfo())
		return nil, err
	} else if offset == nil {
		return nil, nil
	}
	var offsetValue int64
	offsetValue, err = strconv.ParseInt(*offset, 10, 64)
	if err != nil {
		err = core.SDKErrorf(err, "", "parse-int-query-error", common.GetComponentInfo())
		return nil, err
	}
	return core.Int64Ptr(offsetValue), nil
}

// SnapshotOutput : Snapshot object.
type SnapshotOutput struct {
	// Git Config id.
	GitConfigID *string `json:"git_config_id" validate:"required"`

	// Git Config name.
	Name *string `json:"name" validate:"required"`
}

// NewSnapshotOutput : Instantiate SnapshotOutput (Generic Model Constructor)
func (*AppConfigurationV1) NewSnapshotOutput(gitConfigID string, name string) (_model *SnapshotOutput, err error) {
	_model = &SnapshotOutput{
		GitConfigID: core.StringPtr(gitConfigID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalSnapshotOutput unmarshals an instance of SnapshotOutput from the specified map of raw messages.
func UnmarshalSnapshotOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SnapshotOutput)
	err = core.UnmarshalPrimitive(m, "git_config_id", &obj.GitConfigID)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_config_id-error", common.GetComponentInfo())
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

// TargetSegments : TargetSegments struct
type TargetSegments struct {
	// List of segment ids that are used for targeting using the rule.
	Segments []string `json:"segments" validate:"required"`
}

// NewTargetSegments : Instantiate TargetSegments (Generic Model Constructor)
func (*AppConfigurationV1) NewTargetSegments(segments []string) (_model *TargetSegments, err error) {
	_model = &TargetSegments{
		Segments: segments,
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

// UnmarshalTargetSegments unmarshals an instance of TargetSegments from the specified map of raw messages.
func UnmarshalTargetSegments(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetSegments)
	err = core.UnmarshalPrimitive(m, "segments", &obj.Segments)
	if err != nil {
		err = core.SDKErrorf(err, "", "segments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ToggleFeatureOptions : The ToggleFeature options.
type ToggleFeatureOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Feature Id.
	FeatureID *string `json:"feature_id" validate:"required,ne="`

	// The state of the feature flag.
	Enabled *bool `json:"enabled" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewToggleFeatureOptions : Instantiate ToggleFeatureOptions
func (*AppConfigurationV1) NewToggleFeatureOptions(environmentID string, featureID string, enabled bool) *ToggleFeatureOptions {
	return &ToggleFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
		Enabled: core.BoolPtr(enabled),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *ToggleFeatureOptions) SetEnvironmentID(environmentID string) *ToggleFeatureOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetFeatureID : Allow user to set FeatureID
func (_options *ToggleFeatureOptions) SetFeatureID(featureID string) *ToggleFeatureOptions {
	_options.FeatureID = core.StringPtr(featureID)
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *ToggleFeatureOptions) SetEnabled(enabled bool) *ToggleFeatureOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ToggleFeatureOptions) SetHeaders(param map[string]string) *ToggleFeatureOptions {
	options.Headers = param
	return options
}

// UpdateCollectionOptions : The UpdateCollection options.
type UpdateCollectionOptions struct {
	// Collection Id of the collection.
	CollectionID *string `json:"collection_id" validate:"required,ne="`

	// Collection name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Description of the collection, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateCollectionOptions : Instantiate UpdateCollectionOptions
func (*AppConfigurationV1) NewUpdateCollectionOptions(collectionID string) *UpdateCollectionOptions {
	return &UpdateCollectionOptions{
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetCollectionID : Allow user to set CollectionID
func (_options *UpdateCollectionOptions) SetCollectionID(collectionID string) *UpdateCollectionOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateCollectionOptions) SetName(name string) *UpdateCollectionOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateCollectionOptions) SetDescription(description string) *UpdateCollectionOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateCollectionOptions) SetTags(tags string) *UpdateCollectionOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateCollectionOptions) SetHeaders(param map[string]string) *UpdateCollectionOptions {
	options.Headers = param
	return options
}

// UpdateEnvironmentOptions : The UpdateEnvironment options.
type UpdateEnvironmentOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Environment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Environment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.
	ColorCode *string `json:"color_code,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateEnvironmentOptions : Instantiate UpdateEnvironmentOptions
func (*AppConfigurationV1) NewUpdateEnvironmentOptions(environmentID string) *UpdateEnvironmentOptions {
	return &UpdateEnvironmentOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdateEnvironmentOptions) SetEnvironmentID(environmentID string) *UpdateEnvironmentOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateEnvironmentOptions) SetName(name string) *UpdateEnvironmentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateEnvironmentOptions) SetDescription(description string) *UpdateEnvironmentOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateEnvironmentOptions) SetTags(tags string) *UpdateEnvironmentOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetColorCode : Allow user to set ColorCode
func (_options *UpdateEnvironmentOptions) SetColorCode(colorCode string) *UpdateEnvironmentOptions {
	_options.ColorCode = core.StringPtr(colorCode)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEnvironmentOptions) SetHeaders(param map[string]string) *UpdateEnvironmentOptions {
	options.Headers = param
	return options
}

// UpdateFeatureOptions : The UpdateFeature options.
type UpdateFeatureOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Feature Id.
	FeatureID *string `json:"feature_id" validate:"required,ne="`

	// Feature name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Feature description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Value of the feature when it is enabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	EnabledValue interface{} `json:"enabled_value,omitempty"`

	// Value of the feature when it is disabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	DisabledValue interface{} `json:"disabled_value,omitempty"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Rollout percentage associated with feature flag. Supported only for Lite and Enterprise plans.
	RolloutPercentage *int64 `json:"rollout_percentage,omitempty"`

	// Tags associated with the feature, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []FeatureSegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateFeatureOptions : Instantiate UpdateFeatureOptions
func (*AppConfigurationV1) NewUpdateFeatureOptions(environmentID string, featureID string) *UpdateFeatureOptions {
	return &UpdateFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdateFeatureOptions) SetEnvironmentID(environmentID string) *UpdateFeatureOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetFeatureID : Allow user to set FeatureID
func (_options *UpdateFeatureOptions) SetFeatureID(featureID string) *UpdateFeatureOptions {
	_options.FeatureID = core.StringPtr(featureID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateFeatureOptions) SetName(name string) *UpdateFeatureOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateFeatureOptions) SetDescription(description string) *UpdateFeatureOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetEnabledValue : Allow user to set EnabledValue
func (_options *UpdateFeatureOptions) SetEnabledValue(enabledValue interface{}) *UpdateFeatureOptions {
	_options.EnabledValue = enabledValue
	return _options
}

// SetDisabledValue : Allow user to set DisabledValue
func (_options *UpdateFeatureOptions) SetDisabledValue(disabledValue interface{}) *UpdateFeatureOptions {
	_options.DisabledValue = disabledValue
	return _options
}

// SetEnabled : Allow user to set Enabled
func (_options *UpdateFeatureOptions) SetEnabled(enabled bool) *UpdateFeatureOptions {
	_options.Enabled = core.BoolPtr(enabled)
	return _options
}

// SetRolloutPercentage : Allow user to set RolloutPercentage
func (_options *UpdateFeatureOptions) SetRolloutPercentage(rolloutPercentage int64) *UpdateFeatureOptions {
	_options.RolloutPercentage = core.Int64Ptr(rolloutPercentage)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateFeatureOptions) SetTags(tags string) *UpdateFeatureOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetSegmentRules : Allow user to set SegmentRules
func (_options *UpdateFeatureOptions) SetSegmentRules(segmentRules []FeatureSegmentRule) *UpdateFeatureOptions {
	_options.SegmentRules = segmentRules
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *UpdateFeatureOptions) SetCollections(collections []CollectionRef) *UpdateFeatureOptions {
	_options.Collections = collections
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateFeatureOptions) SetHeaders(param map[string]string) *UpdateFeatureOptions {
	options.Headers = param
	return options
}

// UpdateFeatureValuesOptions : The UpdateFeatureValues options.
type UpdateFeatureValuesOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Feature Id.
	FeatureID *string `json:"feature_id" validate:"required,ne="`

	// Feature name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Feature description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the feature, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Value of the feature when it is enabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	EnabledValue interface{} `json:"enabled_value,omitempty"`

	// Value of the feature when it is disabled. The value can be Boolean, Numeric, String - TEXT, String - JSON, String -
	// YAML value as per the `type` and `format` attributes.
	DisabledValue interface{} `json:"disabled_value,omitempty"`

	// Rollout percentage associated with feature flag. Supported only for Lite and Enterprise plans.
	RolloutPercentage *int64 `json:"rollout_percentage,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []FeatureSegmentRule `json:"segment_rules,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateFeatureValuesOptions : Instantiate UpdateFeatureValuesOptions
func (*AppConfigurationV1) NewUpdateFeatureValuesOptions(environmentID string, featureID string) *UpdateFeatureValuesOptions {
	return &UpdateFeatureValuesOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdateFeatureValuesOptions) SetEnvironmentID(environmentID string) *UpdateFeatureValuesOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetFeatureID : Allow user to set FeatureID
func (_options *UpdateFeatureValuesOptions) SetFeatureID(featureID string) *UpdateFeatureValuesOptions {
	_options.FeatureID = core.StringPtr(featureID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateFeatureValuesOptions) SetName(name string) *UpdateFeatureValuesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateFeatureValuesOptions) SetDescription(description string) *UpdateFeatureValuesOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateFeatureValuesOptions) SetTags(tags string) *UpdateFeatureValuesOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetEnabledValue : Allow user to set EnabledValue
func (_options *UpdateFeatureValuesOptions) SetEnabledValue(enabledValue interface{}) *UpdateFeatureValuesOptions {
	_options.EnabledValue = enabledValue
	return _options
}

// SetDisabledValue : Allow user to set DisabledValue
func (_options *UpdateFeatureValuesOptions) SetDisabledValue(disabledValue interface{}) *UpdateFeatureValuesOptions {
	_options.DisabledValue = disabledValue
	return _options
}

// SetRolloutPercentage : Allow user to set RolloutPercentage
func (_options *UpdateFeatureValuesOptions) SetRolloutPercentage(rolloutPercentage int64) *UpdateFeatureValuesOptions {
	_options.RolloutPercentage = core.Int64Ptr(rolloutPercentage)
	return _options
}

// SetSegmentRules : Allow user to set SegmentRules
func (_options *UpdateFeatureValuesOptions) SetSegmentRules(segmentRules []FeatureSegmentRule) *UpdateFeatureValuesOptions {
	_options.SegmentRules = segmentRules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateFeatureValuesOptions) SetHeaders(param map[string]string) *UpdateFeatureValuesOptions {
	options.Headers = param
	return options
}

// UpdateGitconfigOptions : The UpdateGitconfig options.
type UpdateGitconfigOptions struct {
	// Git Config Id.
	GitConfigID *string `json:"git_config_id" validate:"required,ne="`

	// Git config name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	GitConfigName *string `json:"git_config_name,omitempty"`

	// Collection Id.
	CollectionID *string `json:"collection_id,omitempty"`

	// Environment Id.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Git url which will be used to connect to the github account. The url must be formed in this format,
	// https://api.github.com/repos/{owner}/{repo_name} for the personal git account. If you are using the organization
	// account then url must be in this format https://github.{organization_name}.com/api/v3/repos/{owner}/{repo_name} .
	// Note do not provide /(slash) in the beginning or at the end of the url.
	GitURL *string `json:"git_url,omitempty"`

	// Branch name to which you need to write or update the configuration. Just provide the branch name, do not provide any
	// /(slashes) in the beginning or at the end of the branch name. Note make sure branch exists in your repository.
	GitBranch *string `json:"git_branch,omitempty"`

	// Git file path, this is a path where your configuration file will be written. The path must contain the file name
	// with `json` extension. We only create or update `json` extension file. Note do not provide any /(slashes) in the
	// beginning or at the end of the file path.
	GitFilePath *string `json:"git_file_path,omitempty"`

	// Git token, this needs to be provided with enough permission to write and update the file.
	GitToken *string `json:"git_token,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateGitconfigOptions : Instantiate UpdateGitconfigOptions
func (*AppConfigurationV1) NewUpdateGitconfigOptions(gitConfigID string) *UpdateGitconfigOptions {
	return &UpdateGitconfigOptions{
		GitConfigID: core.StringPtr(gitConfigID),
	}
}

// SetGitConfigID : Allow user to set GitConfigID
func (_options *UpdateGitconfigOptions) SetGitConfigID(gitConfigID string) *UpdateGitconfigOptions {
	_options.GitConfigID = core.StringPtr(gitConfigID)
	return _options
}

// SetGitConfigName : Allow user to set GitConfigName
func (_options *UpdateGitconfigOptions) SetGitConfigName(gitConfigName string) *UpdateGitconfigOptions {
	_options.GitConfigName = core.StringPtr(gitConfigName)
	return _options
}

// SetCollectionID : Allow user to set CollectionID
func (_options *UpdateGitconfigOptions) SetCollectionID(collectionID string) *UpdateGitconfigOptions {
	_options.CollectionID = core.StringPtr(collectionID)
	return _options
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdateGitconfigOptions) SetEnvironmentID(environmentID string) *UpdateGitconfigOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetGitURL : Allow user to set GitURL
func (_options *UpdateGitconfigOptions) SetGitURL(gitURL string) *UpdateGitconfigOptions {
	_options.GitURL = core.StringPtr(gitURL)
	return _options
}

// SetGitBranch : Allow user to set GitBranch
func (_options *UpdateGitconfigOptions) SetGitBranch(gitBranch string) *UpdateGitconfigOptions {
	_options.GitBranch = core.StringPtr(gitBranch)
	return _options
}

// SetGitFilePath : Allow user to set GitFilePath
func (_options *UpdateGitconfigOptions) SetGitFilePath(gitFilePath string) *UpdateGitconfigOptions {
	_options.GitFilePath = core.StringPtr(gitFilePath)
	return _options
}

// SetGitToken : Allow user to set GitToken
func (_options *UpdateGitconfigOptions) SetGitToken(gitToken string) *UpdateGitconfigOptions {
	_options.GitToken = core.StringPtr(gitToken)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateGitconfigOptions) SetHeaders(param map[string]string) *UpdateGitconfigOptions {
	options.Headers = param
	return options
}

// UpdateOriginconfigsOptions : The UpdateOriginconfigs options.
type UpdateOriginconfigsOptions struct {
	// List of allowed origins. Specify the parameter as a list of comma separated origins.
	AllowedOrigins []string `json:"allowed_origins" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateOriginconfigsOptions : Instantiate UpdateOriginconfigsOptions
func (*AppConfigurationV1) NewUpdateOriginconfigsOptions(allowedOrigins []string) *UpdateOriginconfigsOptions {
	return &UpdateOriginconfigsOptions{
		AllowedOrigins: allowedOrigins,
	}
}

// SetAllowedOrigins : Allow user to set AllowedOrigins
func (_options *UpdateOriginconfigsOptions) SetAllowedOrigins(allowedOrigins []string) *UpdateOriginconfigsOptions {
	_options.AllowedOrigins = allowedOrigins
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateOriginconfigsOptions) SetHeaders(param map[string]string) *UpdateOriginconfigsOptions {
	options.Headers = param
	return options
}

// UpdatePropertyOptions : The UpdateProperty options.
type UpdatePropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Property name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Property description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Value of the Property. The value can be Boolean, Numeric, SecretRef, String - TEXT, String - JSON, String - YAML as
	// per the `type` and `format` attributes.
	Value interface{} `json:"value,omitempty"`

	// Tags associated with the property, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdatePropertyOptions : Instantiate UpdatePropertyOptions
func (*AppConfigurationV1) NewUpdatePropertyOptions(environmentID string, propertyID string) *UpdatePropertyOptions {
	return &UpdatePropertyOptions{
		EnvironmentID: core.StringPtr(environmentID),
		PropertyID: core.StringPtr(propertyID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdatePropertyOptions) SetEnvironmentID(environmentID string) *UpdatePropertyOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetPropertyID : Allow user to set PropertyID
func (_options *UpdatePropertyOptions) SetPropertyID(propertyID string) *UpdatePropertyOptions {
	_options.PropertyID = core.StringPtr(propertyID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdatePropertyOptions) SetName(name string) *UpdatePropertyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdatePropertyOptions) SetDescription(description string) *UpdatePropertyOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetValue : Allow user to set Value
func (_options *UpdatePropertyOptions) SetValue(value interface{}) *UpdatePropertyOptions {
	_options.Value = value
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdatePropertyOptions) SetTags(tags string) *UpdatePropertyOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetSegmentRules : Allow user to set SegmentRules
func (_options *UpdatePropertyOptions) SetSegmentRules(segmentRules []SegmentRule) *UpdatePropertyOptions {
	_options.SegmentRules = segmentRules
	return _options
}

// SetCollections : Allow user to set Collections
func (_options *UpdatePropertyOptions) SetCollections(collections []CollectionRef) *UpdatePropertyOptions {
	_options.Collections = collections
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePropertyOptions) SetHeaders(param map[string]string) *UpdatePropertyOptions {
	options.Headers = param
	return options
}

// UpdatePropertyValuesOptions : The UpdatePropertyValues options.
type UpdatePropertyValuesOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Property name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Property description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the property, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// Value of the Property. The value can be Boolean, Numeric, SecretRef, String - TEXT, String - JSON, String - YAML as
	// per the `type` and `format` attributes.
	Value interface{} `json:"value,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdatePropertyValuesOptions : Instantiate UpdatePropertyValuesOptions
func (*AppConfigurationV1) NewUpdatePropertyValuesOptions(environmentID string, propertyID string) *UpdatePropertyValuesOptions {
	return &UpdatePropertyValuesOptions{
		EnvironmentID: core.StringPtr(environmentID),
		PropertyID: core.StringPtr(propertyID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdatePropertyValuesOptions) SetEnvironmentID(environmentID string) *UpdatePropertyValuesOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetPropertyID : Allow user to set PropertyID
func (_options *UpdatePropertyValuesOptions) SetPropertyID(propertyID string) *UpdatePropertyValuesOptions {
	_options.PropertyID = core.StringPtr(propertyID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdatePropertyValuesOptions) SetName(name string) *UpdatePropertyValuesOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdatePropertyValuesOptions) SetDescription(description string) *UpdatePropertyValuesOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdatePropertyValuesOptions) SetTags(tags string) *UpdatePropertyValuesOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetValue : Allow user to set Value
func (_options *UpdatePropertyValuesOptions) SetValue(value interface{}) *UpdatePropertyValuesOptions {
	_options.Value = value
	return _options
}

// SetSegmentRules : Allow user to set SegmentRules
func (_options *UpdatePropertyValuesOptions) SetSegmentRules(segmentRules []SegmentRule) *UpdatePropertyValuesOptions {
	_options.SegmentRules = segmentRules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdatePropertyValuesOptions) SetHeaders(param map[string]string) *UpdatePropertyValuesOptions {
	options.Headers = param
	return options
}

// UpdateSegmentOptions : The UpdateSegment options.
type UpdateSegmentOptions struct {
	// Segment Id.
	SegmentID *string `json:"segment_id" validate:"required,ne="`

	// Segment name. Allowed special characters are dot ( . ), hyphen( - ), underscore ( _ ) only.
	Name *string `json:"name,omitempty"`

	// Segment description, allowed special characters are [.,-_ :()$&%#!].
	Description *string `json:"description,omitempty"`

	// Tags associated with the segments, allowed special characters are [_. ,-:].
	Tags *string `json:"tags,omitempty"`

	// List of rules that determine if the entity belongs to the segment during feature / property evaluation. An entity is
	// identified by an unique identifier and the attributes that it defines. Any feature flag and property value
	// evaluation is performed in the context of an entity when it is targeted to segments.
	Rules []Rule `json:"rules,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateSegmentOptions : Instantiate UpdateSegmentOptions
func (*AppConfigurationV1) NewUpdateSegmentOptions(segmentID string) *UpdateSegmentOptions {
	return &UpdateSegmentOptions{
		SegmentID: core.StringPtr(segmentID),
	}
}

// SetSegmentID : Allow user to set SegmentID
func (_options *UpdateSegmentOptions) SetSegmentID(segmentID string) *UpdateSegmentOptions {
	_options.SegmentID = core.StringPtr(segmentID)
	return _options
}

// SetName : Allow user to set Name
func (_options *UpdateSegmentOptions) SetName(name string) *UpdateSegmentOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetDescription : Allow user to set Description
func (_options *UpdateSegmentOptions) SetDescription(description string) *UpdateSegmentOptions {
	_options.Description = core.StringPtr(description)
	return _options
}

// SetTags : Allow user to set Tags
func (_options *UpdateSegmentOptions) SetTags(tags string) *UpdateSegmentOptions {
	_options.Tags = core.StringPtr(tags)
	return _options
}

// SetRules : Allow user to set Rules
func (_options *UpdateSegmentOptions) SetRules(rules []Rule) *UpdateSegmentOptions {
	_options.Rules = rules
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSegmentOptions) SetHeaders(param map[string]string) *UpdateSegmentOptions {
	options.Headers = param
	return options
}

// UpdateWorkflowConfig : UpdateWorkflowConfig struct
// Models which "extend" this model:
// - UpdateWorkflowConfigUpdateExternalServiceNow
// - UpdateWorkflowConfigUpdateIBMServiceNow
type UpdateWorkflowConfig struct {
	// ServiceNow instance URL. Only url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url,omitempty"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn,omitempty"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn,omitempty"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id,omitempty"`
}
func (*UpdateWorkflowConfig) isaUpdateWorkflowConfig() bool {
	return true
}

type UpdateWorkflowConfigIntf interface {
	isaUpdateWorkflowConfig() bool
}

// UnmarshalUpdateWorkflowConfig unmarshals an instance of UpdateWorkflowConfig from the specified map of raw messages.
func UnmarshalUpdateWorkflowConfig(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateWorkflowConfig)
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateWorkflowconfigOptions : The UpdateWorkflowconfig options.
type UpdateWorkflowconfigOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// The request body to update an existing workflow config.
	UpdateWorkflowConfig UpdateWorkflowConfigIntf `json:"UpdateWorkflowConfig" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewUpdateWorkflowconfigOptions : Instantiate UpdateWorkflowconfigOptions
func (*AppConfigurationV1) NewUpdateWorkflowconfigOptions(environmentID string, updateWorkflowConfig UpdateWorkflowConfigIntf) *UpdateWorkflowconfigOptions {
	return &UpdateWorkflowconfigOptions{
		EnvironmentID: core.StringPtr(environmentID),
		UpdateWorkflowConfig: updateWorkflowConfig,
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (_options *UpdateWorkflowconfigOptions) SetEnvironmentID(environmentID string) *UpdateWorkflowconfigOptions {
	_options.EnvironmentID = core.StringPtr(environmentID)
	return _options
}

// SetUpdateWorkflowConfig : Allow user to set UpdateWorkflowConfig
func (_options *UpdateWorkflowconfigOptions) SetUpdateWorkflowConfig(updateWorkflowConfig UpdateWorkflowConfigIntf) *UpdateWorkflowconfigOptions {
	_options.UpdateWorkflowConfig = updateWorkflowConfig
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateWorkflowconfigOptions) SetHeaders(param map[string]string) *UpdateWorkflowconfigOptions {
	options.Headers = param
	return options
}

// UpdateWorkflowconfigResponse : UpdateWorkflowconfigResponse struct
// Models which "extend" this model:
// - UpdateWorkflowconfigResponseExternalServiceNow
// - UpdateWorkflowconfigResponseIBMServiceNow
type UpdateWorkflowconfigResponse struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url,omitempty"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn,omitempty"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type,omitempty"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn,omitempty"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id,omitempty"`
}
func (*UpdateWorkflowconfigResponse) isaUpdateWorkflowconfigResponse() bool {
	return true
}

type UpdateWorkflowconfigResponseIntf interface {
	isaUpdateWorkflowconfigResponse() bool
}

// UnmarshalUpdateWorkflowconfigResponse unmarshals an instance of UpdateWorkflowconfigResponse from the specified map of raw messages.
func UnmarshalUpdateWorkflowconfigResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateWorkflowconfigResponse)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigActionGitConfigPromote : Details of the promote operation.
// This model "extends" ConfigAction
type ConfigActionGitConfigPromote struct {
	// Git commit id will be given as part of the response upon successful git operation.
	GitCommitID *string `json:"git_commit_id" validate:"required"`

	// Git commit message.
	GitCommitMessage *string `json:"git_commit_message" validate:"required"`

	// Latest time when the snapshot was synced to git.
	LastSyncTime *strfmt.DateTime `json:"last_sync_time,omitempty"`
}

func (*ConfigActionGitConfigPromote) isaConfigAction() bool {
	return true
}

// UnmarshalConfigActionGitConfigPromote unmarshals an instance of ConfigActionGitConfigPromote from the specified map of raw messages.
func UnmarshalConfigActionGitConfigPromote(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigActionGitConfigPromote)
	err = core.UnmarshalPrimitive(m, "git_commit_id", &obj.GitCommitID)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "git_commit_message", &obj.GitCommitMessage)
	if err != nil {
		err = core.SDKErrorf(err, "", "git_commit_message-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "last_sync_time", &obj.LastSyncTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "last_sync_time-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigActionGitConfigRestore : Details of the restore operation.
// This model "extends" ConfigAction
type ConfigActionGitConfigRestore struct {
	// The environments array will contain the environment data and it will also contains properties array and features
	// array that belongs to that environment.
	Environments []ImportEnvironmentSchema `json:"environments" validate:"required"`

	// Segments that belongs to the features or properties.
	Segments []ImportSegmentSchema `json:"segments" validate:"required"`
}

func (*ConfigActionGitConfigRestore) isaConfigAction() bool {
	return true
}

// UnmarshalConfigActionGitConfigRestore unmarshals an instance of ConfigActionGitConfigRestore from the specified map of raw messages.
func UnmarshalConfigActionGitConfigRestore(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigActionGitConfigRestore)
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalImportEnvironmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "environments-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "segments", &obj.Segments, UnmarshalImportSegmentSchema)
	if err != nil {
		err = core.SDKErrorf(err, "", "segments-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateWorkflowConfigExternalServiceNow : Details of the External ServiceNow workflow configuration.
// This model "extends" CreateWorkflowConfig
type CreateWorkflowConfigExternalServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url" validate:"required"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

// NewCreateWorkflowConfigExternalServiceNow : Instantiate CreateWorkflowConfigExternalServiceNow (Generic Model Constructor)
func (*AppConfigurationV1) NewCreateWorkflowConfigExternalServiceNow(workflowURL string, approvalGroupName string, approvalExpiration int64, workflowCredentials *ExternalServiceNowCredentials, enabled bool) (_model *CreateWorkflowConfigExternalServiceNow, err error) {
	_model = &CreateWorkflowConfigExternalServiceNow{
		WorkflowURL: core.StringPtr(workflowURL),
		ApprovalGroupName: core.StringPtr(approvalGroupName),
		ApprovalExpiration: core.Int64Ptr(approvalExpiration),
		WorkflowCredentials: workflowCredentials,
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*CreateWorkflowConfigExternalServiceNow) isaCreateWorkflowConfig() bool {
	return true
}

// UnmarshalCreateWorkflowConfigExternalServiceNow unmarshals an instance of CreateWorkflowConfigExternalServiceNow from the specified map of raw messages.
func UnmarshalCreateWorkflowConfigExternalServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateWorkflowConfigExternalServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// CreateWorkflowConfigIBMServiceNow : Details of the IBM ServiceNow workflow configuration.
// This model "extends" CreateWorkflowConfig
type CreateWorkflowConfigIBMServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn" validate:"required"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn" validate:"required"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

// NewCreateWorkflowConfigIBMServiceNow : Instantiate CreateWorkflowConfigIBMServiceNow (Generic Model Constructor)
func (*AppConfigurationV1) NewCreateWorkflowConfigIBMServiceNow(serviceCrn string, workflowType string, approvalExpiration int64, smInstanceCrn string, secretID string, enabled bool) (_model *CreateWorkflowConfigIBMServiceNow, err error) {
	_model = &CreateWorkflowConfigIBMServiceNow{
		ServiceCrn: core.StringPtr(serviceCrn),
		WorkflowType: core.StringPtr(workflowType),
		ApprovalExpiration: core.Int64Ptr(approvalExpiration),
		SmInstanceCrn: core.StringPtr(smInstanceCrn),
		SecretID: core.StringPtr(secretID),
		Enabled: core.BoolPtr(enabled),
	}
	err = core.ValidateStruct(_model, "required parameters")
	if err != nil {
		err = core.SDKErrorf(err, "", "model-missing-required", common.GetComponentInfo())
	}
	return
}

func (*CreateWorkflowConfigIBMServiceNow) isaCreateWorkflowConfig() bool {
	return true
}

// UnmarshalCreateWorkflowConfigIBMServiceNow unmarshals an instance of CreateWorkflowConfigIBMServiceNow from the specified map of raw messages.
func UnmarshalCreateWorkflowConfigIBMServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateWorkflowConfigIBMServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// CreateWorkflowconfigResponseExternalServiceNow : Details of the External ServiceNow workflow configuration.
// This model "extends" CreateWorkflowconfigResponse
type CreateWorkflowconfigResponseExternalServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url" validate:"required"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

func (*CreateWorkflowconfigResponseExternalServiceNow) isaCreateWorkflowconfigResponse() bool {
	return true
}

// UnmarshalCreateWorkflowconfigResponseExternalServiceNow unmarshals an instance of CreateWorkflowconfigResponseExternalServiceNow from the specified map of raw messages.
func UnmarshalCreateWorkflowconfigResponseExternalServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateWorkflowconfigResponseExternalServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// CreateWorkflowconfigResponseIBMServiceNow : Details of the IBM ServiceNow workflow configuration.
// This model "extends" CreateWorkflowconfigResponse
type CreateWorkflowconfigResponseIBMServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn" validate:"required"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn" validate:"required"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

func (*CreateWorkflowconfigResponseIBMServiceNow) isaCreateWorkflowconfigResponse() bool {
	return true
}

// UnmarshalCreateWorkflowconfigResponseIBMServiceNow unmarshals an instance of CreateWorkflowconfigResponseIBMServiceNow from the specified map of raw messages.
func UnmarshalCreateWorkflowconfigResponseIBMServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CreateWorkflowconfigResponseIBMServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// ListWorkflowconfigResponseExternalServiceNow : Details of the External ServiceNow workflow configuration.
// This model "extends" ListWorkflowconfigResponse
type ListWorkflowconfigResponseExternalServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url" validate:"required"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

func (*ListWorkflowconfigResponseExternalServiceNow) isaListWorkflowconfigResponse() bool {
	return true
}

// UnmarshalListWorkflowconfigResponseExternalServiceNow unmarshals an instance of ListWorkflowconfigResponseExternalServiceNow from the specified map of raw messages.
func UnmarshalListWorkflowconfigResponseExternalServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListWorkflowconfigResponseExternalServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// ListWorkflowconfigResponseIBMServiceNow : Details of the IBM ServiceNow workflow configuration.
// This model "extends" ListWorkflowconfigResponse
type ListWorkflowconfigResponseIBMServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn" validate:"required"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn" validate:"required"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

func (*ListWorkflowconfigResponseIBMServiceNow) isaListWorkflowconfigResponse() bool {
	return true
}

// UnmarshalListWorkflowconfigResponseIBMServiceNow unmarshals an instance of ListWorkflowconfigResponseIBMServiceNow from the specified map of raw messages.
func UnmarshalListWorkflowconfigResponseIBMServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListWorkflowconfigResponseIBMServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// UpdateWorkflowConfigUpdateExternalServiceNow : External ServiceNow workflow configuration attributes to be updated.
// This model "extends" UpdateWorkflowConfig
type UpdateWorkflowConfigUpdateExternalServiceNow struct {
	// ServiceNow instance URL. Only url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url,omitempty"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`
}

func (*UpdateWorkflowConfigUpdateExternalServiceNow) isaUpdateWorkflowConfig() bool {
	return true
}

// UnmarshalUpdateWorkflowConfigUpdateExternalServiceNow unmarshals an instance of UpdateWorkflowConfigUpdateExternalServiceNow from the specified map of raw messages.
func UnmarshalUpdateWorkflowConfigUpdateExternalServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateWorkflowConfigUpdateExternalServiceNow)
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
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

// UpdateWorkflowConfigUpdateIBMServiceNow : IBM ServiceNow workflow configuration attributes to be updated.
// This model "extends" UpdateWorkflowConfig
type UpdateWorkflowConfigUpdateIBMServiceNow struct {
	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn,omitempty"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration,omitempty"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn,omitempty"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id,omitempty"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled,omitempty"`
}

func (*UpdateWorkflowConfigUpdateIBMServiceNow) isaUpdateWorkflowConfig() bool {
	return true
}

// UnmarshalUpdateWorkflowConfigUpdateIBMServiceNow unmarshals an instance of UpdateWorkflowConfigUpdateIBMServiceNow from the specified map of raw messages.
func UnmarshalUpdateWorkflowConfigUpdateIBMServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateWorkflowConfigUpdateIBMServiceNow)
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
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

// UpdateWorkflowconfigResponseExternalServiceNow : Details of the External ServiceNow workflow configuration.
// This model "extends" UpdateWorkflowconfigResponse
type UpdateWorkflowconfigResponseExternalServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service now url https://xxxxx.service-now.com allowed, xxxxx is the service now instance id.
	WorkflowURL *string `json:"workflow_url" validate:"required"`

	// Group name of personals who can approve the Change Request on your ServiceNow. It must be first registered in your
	// ServiceNow then it must be added here.
	ApprovalGroupName *string `json:"approval_group_name" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// The credentials of the External ServiceNow instance.
	WorkflowCredentials *ExternalServiceNowCredentials `json:"workflow_credentials" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

func (*UpdateWorkflowconfigResponseExternalServiceNow) isaUpdateWorkflowconfigResponse() bool {
	return true
}

// UnmarshalUpdateWorkflowconfigResponseExternalServiceNow unmarshals an instance of UpdateWorkflowconfigResponseExternalServiceNow from the specified map of raw messages.
func UnmarshalUpdateWorkflowconfigResponseExternalServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateWorkflowconfigResponseExternalServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_url", &obj.WorkflowURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_group_name", &obj.ApprovalGroupName)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_group_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "workflow_credentials", &obj.WorkflowCredentials, UnmarshalExternalServiceNowCredentials)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_credentials-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

// UpdateWorkflowconfigResponseIBMServiceNow : Details of the IBM ServiceNow workflow configuration.
// This model "extends" UpdateWorkflowconfigResponse
type UpdateWorkflowconfigResponseIBMServiceNow struct {
	// Environment name of workflow config in which it is created.
	EnvironmentName *string `json:"environment_name,omitempty"`

	// Environment ID of workflow config in which it is created.
	EnvironmentID *string `json:"environment_id,omitempty"`

	// Only service crn will be allowed. Example: `crn:v1:staging:staging:appservice:us-south::::`.
	ServiceCrn *string `json:"service_crn" validate:"required"`

	// Allowed value is `SERVICENOW_IBM` case-sensitive.
	WorkflowType *string `json:"workflow_type" validate:"required"`

	// Integer number identifies as hours which helps in adding approval start and end time to the created Change Request.
	ApprovalExpiration *int64 `json:"approval_expiration" validate:"required"`

	// Only Secret Manager instance crn will be allowed. Example:
	// `crn:v1:staging:public:secrets-manager:eu-gb:a/3268cfe9e25d411122f9a731a:0a23274-92d0a-4d42-b1fa-d15b4293cd::`.
	SmInstanceCrn *string `json:"sm_instance_crn" validate:"required"`

	// Provide the arbitary secret key id which holds the api key to interact with service now. This is required to perform
	// action on ServiceNow like Create CR or Close CR.
	SecretID *string `json:"secret_id" validate:"required"`

	// This option enables the workflow configuration per environment. User must set it to true if they wish to create
	// Change Request for flag state changes.
	Enabled *bool `json:"enabled" validate:"required"`

	// Creation time of the workflow configs.
	CreatedTime *strfmt.DateTime `json:"created_time,omitempty"`

	// Last modified time of the workflow configs.
	UpdatedTime *strfmt.DateTime `json:"updated_time,omitempty"`

	// Workflow Config URL.
	Href *string `json:"href,omitempty"`
}

func (*UpdateWorkflowconfigResponseIBMServiceNow) isaUpdateWorkflowconfigResponse() bool {
	return true
}

// UnmarshalUpdateWorkflowconfigResponseIBMServiceNow unmarshals an instance of UpdateWorkflowconfigResponseIBMServiceNow from the specified map of raw messages.
func UnmarshalUpdateWorkflowconfigResponseIBMServiceNow(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdateWorkflowconfigResponseIBMServiceNow)
	err = core.UnmarshalPrimitive(m, "environment_name", &obj.EnvironmentName)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		err = core.SDKErrorf(err, "", "environment_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "service_crn", &obj.ServiceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "service_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "workflow_type", &obj.WorkflowType)
	if err != nil {
		err = core.SDKErrorf(err, "", "workflow_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "approval_expiration", &obj.ApprovalExpiration)
	if err != nil {
		err = core.SDKErrorf(err, "", "approval_expiration-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sm_instance_crn", &obj.SmInstanceCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "sm_instance_crn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "secret_id", &obj.SecretID)
	if err != nil {
		err = core.SDKErrorf(err, "", "secret_id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		err = core.SDKErrorf(err, "", "enabled-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "created_time-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		err = core.SDKErrorf(err, "", "updated_time-error", common.GetComponentInfo())
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

//
// EnvironmentsPager can be used to simplify the use of the "ListEnvironments" method.
//
type EnvironmentsPager struct {
	hasNext bool
	options *ListEnvironmentsOptions
	client  *AppConfigurationV1
	pageContext struct {
		next *int64
	}
}

// NewEnvironmentsPager returns a new EnvironmentsPager instance.
func (appConfiguration *AppConfigurationV1) NewEnvironmentsPager(options *ListEnvironmentsOptions) (pager *EnvironmentsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListEnvironmentsOptions = *options
	pager = &EnvironmentsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  appConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *EnvironmentsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *EnvironmentsPager) GetNextWithContext(ctx context.Context) (page []Environment, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListEnvironmentsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Environments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *EnvironmentsPager) GetAllWithContext(ctx context.Context) (allItems []Environment, err error) {
	for pager.HasNext() {
		var nextPage []Environment
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
func (pager *EnvironmentsPager) GetNext() (page []Environment, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *EnvironmentsPager) GetAll() (allItems []Environment, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// CollectionsPager can be used to simplify the use of the "ListCollections" method.
//
type CollectionsPager struct {
	hasNext bool
	options *ListCollectionsOptions
	client  *AppConfigurationV1
	pageContext struct {
		next *int64
	}
}

// NewCollectionsPager returns a new CollectionsPager instance.
func (appConfiguration *AppConfigurationV1) NewCollectionsPager(options *ListCollectionsOptions) (pager *CollectionsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListCollectionsOptions = *options
	pager = &CollectionsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  appConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *CollectionsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *CollectionsPager) GetNextWithContext(ctx context.Context) (page []Collection, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListCollectionsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Collections

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *CollectionsPager) GetAllWithContext(ctx context.Context) (allItems []Collection, err error) {
	for pager.HasNext() {
		var nextPage []Collection
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
func (pager *CollectionsPager) GetNext() (page []Collection, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *CollectionsPager) GetAll() (allItems []Collection, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// FeaturesPager can be used to simplify the use of the "ListFeatures" method.
//
type FeaturesPager struct {
	hasNext bool
	options *ListFeaturesOptions
	client  *AppConfigurationV1
	pageContext struct {
		next *int64
	}
}

// NewFeaturesPager returns a new FeaturesPager instance.
func (appConfiguration *AppConfigurationV1) NewFeaturesPager(options *ListFeaturesOptions) (pager *FeaturesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListFeaturesOptions = *options
	pager = &FeaturesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  appConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *FeaturesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *FeaturesPager) GetNextWithContext(ctx context.Context) (page []Feature, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListFeaturesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Features

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *FeaturesPager) GetAllWithContext(ctx context.Context) (allItems []Feature, err error) {
	for pager.HasNext() {
		var nextPage []Feature
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
func (pager *FeaturesPager) GetNext() (page []Feature, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *FeaturesPager) GetAll() (allItems []Feature, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// PropertiesPager can be used to simplify the use of the "ListProperties" method.
//
type PropertiesPager struct {
	hasNext bool
	options *ListPropertiesOptions
	client  *AppConfigurationV1
	pageContext struct {
		next *int64
	}
}

// NewPropertiesPager returns a new PropertiesPager instance.
func (appConfiguration *AppConfigurationV1) NewPropertiesPager(options *ListPropertiesOptions) (pager *PropertiesPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListPropertiesOptions = *options
	pager = &PropertiesPager{
		hasNext: true,
		options: &optionsCopy,
		client:  appConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *PropertiesPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *PropertiesPager) GetNextWithContext(ctx context.Context) (page []Property, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListPropertiesWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Properties

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *PropertiesPager) GetAllWithContext(ctx context.Context) (allItems []Property, err error) {
	for pager.HasNext() {
		var nextPage []Property
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
func (pager *PropertiesPager) GetNext() (page []Property, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *PropertiesPager) GetAll() (allItems []Property, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// SegmentsPager can be used to simplify the use of the "ListSegments" method.
//
type SegmentsPager struct {
	hasNext bool
	options *ListSegmentsOptions
	client  *AppConfigurationV1
	pageContext struct {
		next *int64
	}
}

// NewSegmentsPager returns a new SegmentsPager instance.
func (appConfiguration *AppConfigurationV1) NewSegmentsPager(options *ListSegmentsOptions) (pager *SegmentsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListSegmentsOptions = *options
	pager = &SegmentsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  appConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SegmentsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SegmentsPager) GetNextWithContext(ctx context.Context) (page []Segment, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSegmentsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.Segments

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SegmentsPager) GetAllWithContext(ctx context.Context) (allItems []Segment, err error) {
	for pager.HasNext() {
		var nextPage []Segment
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
func (pager *SegmentsPager) GetNext() (page []Segment, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SegmentsPager) GetAll() (allItems []Segment, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

//
// SnapshotsPager can be used to simplify the use of the "ListSnapshots" method.
//
type SnapshotsPager struct {
	hasNext bool
	options *ListSnapshotsOptions
	client  *AppConfigurationV1
	pageContext struct {
		next *int64
	}
}

// NewSnapshotsPager returns a new SnapshotsPager instance.
func (appConfiguration *AppConfigurationV1) NewSnapshotsPager(options *ListSnapshotsOptions) (pager *SnapshotsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListSnapshotsOptions = *options
	pager = &SnapshotsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  appConfiguration,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *SnapshotsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *SnapshotsPager) GetNextWithContext(ctx context.Context) (page []GitConfig, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListSnapshotsWithContext(ctx, pager.options)
	if err != nil {
		err = core.RepurposeSDKProblem(err, "error-getting-next-page")
		return
	}

	var next *int64
	if result.Next != nil {
		var offset *int64
		offset, err = core.GetQueryParamAsInt(result.Next.Href, "offset")
		if err != nil {
			errMsg := fmt.Sprintf("error retrieving 'offset' query parameter from URL '%s': %s", *result.Next.Href, err.Error())
			err = core.SDKErrorf(err, errMsg, "get-query-error", common.GetComponentInfo())
			return
		}
		next = offset
	}
	pager.pageContext.next = next
	pager.hasNext = (pager.pageContext.next != nil)
	page = result.GitConfig

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *SnapshotsPager) GetAllWithContext(ctx context.Context) (allItems []GitConfig, err error) {
	for pager.HasNext() {
		var nextPage []GitConfig
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
func (pager *SnapshotsPager) GetNext() (page []GitConfig, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *SnapshotsPager) GetAll() (allItems []GitConfig, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
