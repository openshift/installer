/**
 * (C) Copyright IBM Corp. 2021.
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
 * IBM OpenAPI SDK Code Generator Version: 3.22.0-937b9a1c-20201211-223043
 */


// Package appconfigurationv1 : Operations and models for the AppConfigurationV1 service
package appconfigurationv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/appconfiguration-go-admin-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/go-openapi/strfmt"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// AppConfigurationV1 : ReST APIs for App Configuration
//
// Version: 1.0
// See: https://{DomainName}/docs/app-configuration/
type AppConfigurationV1 struct {
	Service *core.BaseService
}

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "app_configuration"

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
			return
		}
	}

	appConfiguration, err = NewAppConfigurationV1(options)
	if err != nil {
		return
	}

	err = appConfiguration.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = appConfiguration.Service.SetServiceURL(options.URL)
	}
	return
}

// NewAppConfigurationV1 : constructs an instance of AppConfigurationV1 with passed in options.
func NewAppConfigurationV1(options *AppConfigurationV1Options) (service *AppConfigurationV1, err error) {
	serviceOptions := &core.ServiceOptions{
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

	service = &AppConfigurationV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
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

// SetServiceURL sets the service URL
func (appConfiguration *AppConfigurationV1) SetServiceURL(url string) error {
	return appConfiguration.Service.SetServiceURL(url)
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
	return appConfiguration.ListEnvironmentsWithContext(context.Background(), listEnvironmentsOptions)
}

// ListEnvironmentsWithContext is an alternate form of the ListEnvironments method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListEnvironmentsWithContext(ctx context.Context, listEnvironmentsOptions *ListEnvironmentsOptions) (result *EnvironmentList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listEnvironmentsOptions, "listEnvironmentsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments`, nil)
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironmentList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateEnvironment : Create Environment
// Create an environment.
func (appConfiguration *AppConfigurationV1) CreateEnvironment(createEnvironmentOptions *CreateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	return appConfiguration.CreateEnvironmentWithContext(context.Background(), createEnvironmentOptions)
}

// CreateEnvironmentWithContext is an alternate form of the CreateEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateEnvironmentWithContext(ctx context.Context, createEnvironmentOptions *CreateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createEnvironmentOptions, "createEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createEnvironmentOptions, "createEnvironmentOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/environments`, nil)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateEnvironment : Update Environment
// Update an environment.
func (appConfiguration *AppConfigurationV1) UpdateEnvironment(updateEnvironmentOptions *UpdateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdateEnvironmentWithContext(context.Background(), updateEnvironmentOptions)
}

// UpdateEnvironmentWithContext is an alternate form of the UpdateEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateEnvironmentWithContext(ctx context.Context, updateEnvironmentOptions *UpdateEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEnvironmentOptions, "updateEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateEnvironmentOptions, "updateEnvironmentOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetEnvironment : Get Environment
// Retrieve the details of the environment.
func (appConfiguration *AppConfigurationV1) GetEnvironment(getEnvironmentOptions *GetEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	return appConfiguration.GetEnvironmentWithContext(context.Background(), getEnvironmentOptions)
}

// GetEnvironmentWithContext is an alternate form of the GetEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetEnvironmentWithContext(ctx context.Context, getEnvironmentOptions *GetEnvironmentOptions) (result *Environment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEnvironmentOptions, "getEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEnvironmentOptions, "getEnvironmentOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEnvironment)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteEnvironment : Delete Environment
// Delete an Environment.
func (appConfiguration *AppConfigurationV1) DeleteEnvironment(deleteEnvironmentOptions *DeleteEnvironmentOptions) (response *core.DetailedResponse, err error) {
	return appConfiguration.DeleteEnvironmentWithContext(context.Background(), deleteEnvironmentOptions)
}

// DeleteEnvironmentWithContext is an alternate form of the DeleteEnvironment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteEnvironmentWithContext(ctx context.Context, deleteEnvironmentOptions *DeleteEnvironmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEnvironmentOptions, "deleteEnvironmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEnvironmentOptions, "deleteEnvironmentOptions")
	if err != nil {
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
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)

	return
}

// ListCollections : Get list of Collections
// List of all the collections in the App Configuration service instance.
func (appConfiguration *AppConfigurationV1) ListCollections(listCollectionsOptions *ListCollectionsOptions) (result *CollectionList, response *core.DetailedResponse, err error) {
	return appConfiguration.ListCollectionsWithContext(context.Background(), listCollectionsOptions)
}

// ListCollectionsWithContext is an alternate form of the ListCollections method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListCollectionsWithContext(ctx context.Context, listCollectionsOptions *ListCollectionsOptions) (result *CollectionList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listCollectionsOptions, "listCollectionsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections`, nil)
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectionList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateCollection : Create Collection
// Create a collection.
func (appConfiguration *AppConfigurationV1) CreateCollection(createCollectionOptions *CreateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	return appConfiguration.CreateCollectionWithContext(context.Background(), createCollectionOptions)
}

// CreateCollectionWithContext is an alternate form of the CreateCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateCollectionWithContext(ctx context.Context, createCollectionOptions *CreateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createCollectionOptions, "createCollectionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createCollectionOptions, "createCollectionOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/collections`, nil)
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectionLite)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateCollection : Update Collection
// Update the collection name, tags and description. Collection Id cannot be updated.
func (appConfiguration *AppConfigurationV1) UpdateCollection(updateCollectionOptions *UpdateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdateCollectionWithContext(context.Background(), updateCollectionOptions)
}

// UpdateCollectionWithContext is an alternate form of the UpdateCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateCollectionWithContext(ctx context.Context, updateCollectionOptions *UpdateCollectionOptions) (result *CollectionLite, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateCollectionOptions, "updateCollectionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateCollectionOptions, "updateCollectionOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollectionLite)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetCollection : Get Collection
// Retrieve the details of the collection.
func (appConfiguration *AppConfigurationV1) GetCollection(getCollectionOptions *GetCollectionOptions) (result *Collection, response *core.DetailedResponse, err error) {
	return appConfiguration.GetCollectionWithContext(context.Background(), getCollectionOptions)
}

// GetCollectionWithContext is an alternate form of the GetCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetCollectionWithContext(ctx context.Context, getCollectionOptions *GetCollectionOptions) (result *Collection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCollectionOptions, "getCollectionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getCollectionOptions, "getCollectionOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalCollection)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteCollection : Delete Collection
// Delete the collection.
func (appConfiguration *AppConfigurationV1) DeleteCollection(deleteCollectionOptions *DeleteCollectionOptions) (response *core.DetailedResponse, err error) {
	return appConfiguration.DeleteCollectionWithContext(context.Background(), deleteCollectionOptions)
}

// DeleteCollectionWithContext is an alternate form of the DeleteCollection method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteCollectionWithContext(ctx context.Context, deleteCollectionOptions *DeleteCollectionOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteCollectionOptions, "deleteCollectionOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteCollectionOptions, "deleteCollectionOptions")
	if err != nil {
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
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)

	return
}

// ListFeatures : Get list of Features
// List all the feature flags in the specified environment.
func (appConfiguration *AppConfigurationV1) ListFeatures(listFeaturesOptions *ListFeaturesOptions) (result *FeaturesList, response *core.DetailedResponse, err error) {
	return appConfiguration.ListFeaturesWithContext(context.Background(), listFeaturesOptions)
}

// ListFeaturesWithContext is an alternate form of the ListFeatures method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListFeaturesWithContext(ctx context.Context, listFeaturesOptions *ListFeaturesOptions) (result *FeaturesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listFeaturesOptions, "listFeaturesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listFeaturesOptions, "listFeaturesOptions")
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeaturesList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateFeature : Create Feature
// Create a feature flag.
func (appConfiguration *AppConfigurationV1) CreateFeature(createFeatureOptions *CreateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	return appConfiguration.CreateFeatureWithContext(context.Background(), createFeatureOptions)
}

// CreateFeatureWithContext is an alternate form of the CreateFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateFeatureWithContext(ctx context.Context, createFeatureOptions *CreateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createFeatureOptions, "createFeatureOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createFeatureOptions, "createFeatureOptions")
	if err != nil {
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
	if createFeatureOptions.Enabled != nil {
		body["enabled"] = createFeatureOptions.Enabled
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateFeature : Update Feature
// Update a feature flag details.
func (appConfiguration *AppConfigurationV1) UpdateFeature(updateFeatureOptions *UpdateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdateFeatureWithContext(context.Background(), updateFeatureOptions)
}

// UpdateFeatureWithContext is an alternate form of the UpdateFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateFeatureWithContext(ctx context.Context, updateFeatureOptions *UpdateFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFeatureOptions, "updateFeatureOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateFeatureOptions, "updateFeatureOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateFeatureValues : Update Feature Values
// Update the feature values. This method can be executed only by the `writer` role. This method allows the update of
// feature name, feature enabled_value, feature disabled_value, tags, description and feature segment rules, however
// this method does not allow toggling the feature flag and assigning feature to a collection.
func (appConfiguration *AppConfigurationV1) UpdateFeatureValues(updateFeatureValuesOptions *UpdateFeatureValuesOptions) (result *Feature, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdateFeatureValuesWithContext(context.Background(), updateFeatureValuesOptions)
}

// UpdateFeatureValuesWithContext is an alternate form of the UpdateFeatureValues method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateFeatureValuesWithContext(ctx context.Context, updateFeatureValuesOptions *UpdateFeatureValuesOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateFeatureValuesOptions, "updateFeatureValuesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateFeatureValuesOptions, "updateFeatureValuesOptions")
	if err != nil {
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
	if updateFeatureValuesOptions.SegmentRules != nil {
		body["segment_rules"] = updateFeatureValuesOptions.SegmentRules
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetFeature : Get Feature
// Retrieve details of a feature.
func (appConfiguration *AppConfigurationV1) GetFeature(getFeatureOptions *GetFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	return appConfiguration.GetFeatureWithContext(context.Background(), getFeatureOptions)
}

// GetFeatureWithContext is an alternate form of the GetFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetFeatureWithContext(ctx context.Context, getFeatureOptions *GetFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getFeatureOptions, "getFeatureOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getFeatureOptions, "getFeatureOptions")
	if err != nil {
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
		builder.AddQuery("include", fmt.Sprint(*getFeatureOptions.Include))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteFeature : Delete Feature
// Delete a feature flag.
func (appConfiguration *AppConfigurationV1) DeleteFeature(deleteFeatureOptions *DeleteFeatureOptions) (response *core.DetailedResponse, err error) {
	return appConfiguration.DeleteFeatureWithContext(context.Background(), deleteFeatureOptions)
}

// DeleteFeatureWithContext is an alternate form of the DeleteFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteFeatureWithContext(ctx context.Context, deleteFeatureOptions *DeleteFeatureOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteFeatureOptions, "deleteFeatureOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteFeatureOptions, "deleteFeatureOptions")
	if err != nil {
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
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)

	return
}

// ToggleFeature : Toggle Feature
// Toggle a feature.
func (appConfiguration *AppConfigurationV1) ToggleFeature(toggleFeatureOptions *ToggleFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	return appConfiguration.ToggleFeatureWithContext(context.Background(), toggleFeatureOptions)
}

// ToggleFeatureWithContext is an alternate form of the ToggleFeature method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ToggleFeatureWithContext(ctx context.Context, toggleFeatureOptions *ToggleFeatureOptions) (result *Feature, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(toggleFeatureOptions, "toggleFeatureOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(toggleFeatureOptions, "toggleFeatureOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalFeature)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListProperties : Get list of Properties
// List all the properties in the specified environment.
func (appConfiguration *AppConfigurationV1) ListProperties(listPropertiesOptions *ListPropertiesOptions) (result *PropertiesList, response *core.DetailedResponse, err error) {
	return appConfiguration.ListPropertiesWithContext(context.Background(), listPropertiesOptions)
}

// ListPropertiesWithContext is an alternate form of the ListProperties method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListPropertiesWithContext(ctx context.Context, listPropertiesOptions *ListPropertiesOptions) (result *PropertiesList, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listPropertiesOptions, "listPropertiesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listPropertiesOptions, "listPropertiesOptions")
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalPropertiesList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateProperty : Create Property
// Create a Property.
func (appConfiguration *AppConfigurationV1) CreateProperty(createPropertyOptions *CreatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	return appConfiguration.CreatePropertyWithContext(context.Background(), createPropertyOptions)
}

// CreatePropertyWithContext is an alternate form of the CreateProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreatePropertyWithContext(ctx context.Context, createPropertyOptions *CreatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createPropertyOptions, "createPropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createPropertyOptions, "createPropertyOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateProperty : Update Property
// Update a Property.
func (appConfiguration *AppConfigurationV1) UpdateProperty(updatePropertyOptions *UpdatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdatePropertyWithContext(context.Background(), updatePropertyOptions)
}

// UpdatePropertyWithContext is an alternate form of the UpdateProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdatePropertyWithContext(ctx context.Context, updatePropertyOptions *UpdatePropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePropertyOptions, "updatePropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updatePropertyOptions, "updatePropertyOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdatePropertyValues : Update Property values
// Update the property values. This method can be executed by the `writer` role. Property value and targeting rules can
// be updated, however this method does not allow assigning property to a collection.
func (appConfiguration *AppConfigurationV1) UpdatePropertyValues(updatePropertyValuesOptions *UpdatePropertyValuesOptions) (result *Property, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdatePropertyValuesWithContext(context.Background(), updatePropertyValuesOptions)
}

// UpdatePropertyValuesWithContext is an alternate form of the UpdatePropertyValues method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdatePropertyValuesWithContext(ctx context.Context, updatePropertyValuesOptions *UpdatePropertyValuesOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updatePropertyValuesOptions, "updatePropertyValuesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updatePropertyValuesOptions, "updatePropertyValuesOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetProperty : Get Property
// Retrieve details of a property.
func (appConfiguration *AppConfigurationV1) GetProperty(getPropertyOptions *GetPropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	return appConfiguration.GetPropertyWithContext(context.Background(), getPropertyOptions)
}

// GetPropertyWithContext is an alternate form of the GetProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetPropertyWithContext(ctx context.Context, getPropertyOptions *GetPropertyOptions) (result *Property, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getPropertyOptions, "getPropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getPropertyOptions, "getPropertyOptions")
	if err != nil {
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
		builder.AddQuery("include", fmt.Sprint(*getPropertyOptions.Include))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalProperty)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteProperty : Delete Property
// Delete a Property.
func (appConfiguration *AppConfigurationV1) DeleteProperty(deletePropertyOptions *DeletePropertyOptions) (response *core.DetailedResponse, err error) {
	return appConfiguration.DeletePropertyWithContext(context.Background(), deletePropertyOptions)
}

// DeletePropertyWithContext is an alternate form of the DeleteProperty method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeletePropertyWithContext(ctx context.Context, deletePropertyOptions *DeletePropertyOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deletePropertyOptions, "deletePropertyOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deletePropertyOptions, "deletePropertyOptions")
	if err != nil {
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
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)

	return
}

// ListSegments : Get list of Segments
// List all the segments.
func (appConfiguration *AppConfigurationV1) ListSegments(listSegmentsOptions *ListSegmentsOptions) (result *SegmentsList, response *core.DetailedResponse, err error) {
	return appConfiguration.ListSegmentsWithContext(context.Background(), listSegmentsOptions)
}

// ListSegmentsWithContext is an alternate form of the ListSegments method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) ListSegmentsWithContext(ctx context.Context, listSegmentsOptions *ListSegmentsOptions) (result *SegmentsList, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listSegmentsOptions, "listSegmentsOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments`, nil)
	if err != nil {
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

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegmentsList)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateSegment : Create Segment
// Create a segment.
func (appConfiguration *AppConfigurationV1) CreateSegment(createSegmentOptions *CreateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	return appConfiguration.CreateSegmentWithContext(context.Background(), createSegmentOptions)
}

// CreateSegmentWithContext is an alternate form of the CreateSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) CreateSegmentWithContext(ctx context.Context, createSegmentOptions *CreateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createSegmentOptions, "createSegmentOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = appConfiguration.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(appConfiguration.Service.Options.URL, `/segments`, nil)
	if err != nil {
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
	if createSegmentOptions.Description != nil {
		body["description"] = createSegmentOptions.Description
	}
	if createSegmentOptions.Tags != nil {
		body["tags"] = createSegmentOptions.Tags
	}
	if createSegmentOptions.Rules != nil {
		body["rules"] = createSegmentOptions.Rules
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
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegment)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSegment : Update Segment
// Update the segment properties.
func (appConfiguration *AppConfigurationV1) UpdateSegment(updateSegmentOptions *UpdateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	return appConfiguration.UpdateSegmentWithContext(context.Background(), updateSegmentOptions)
}

// UpdateSegmentWithContext is an alternate form of the UpdateSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) UpdateSegmentWithContext(ctx context.Context, updateSegmentOptions *UpdateSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSegmentOptions, "updateSegmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSegmentOptions, "updateSegmentOptions")
	if err != nil {
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
		return
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegment)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetSegment : Get Segment
// Retrieve details of a segment.
func (appConfiguration *AppConfigurationV1) GetSegment(getSegmentOptions *GetSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	return appConfiguration.GetSegmentWithContext(context.Background(), getSegmentOptions)
}

// GetSegmentWithContext is an alternate form of the GetSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) GetSegmentWithContext(ctx context.Context, getSegmentOptions *GetSegmentOptions) (result *Segment, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getSegmentOptions, "getSegmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getSegmentOptions, "getSegmentOptions")
	if err != nil {
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
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = appConfiguration.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSegment)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteSegment : Delete Segment
// Delete a segment.
func (appConfiguration *AppConfigurationV1) DeleteSegment(deleteSegmentOptions *DeleteSegmentOptions) (response *core.DetailedResponse, err error) {
	return appConfiguration.DeleteSegmentWithContext(context.Background(), deleteSegmentOptions)
}

// DeleteSegmentWithContext is an alternate form of the DeleteSegment method which supports a Context parameter
func (appConfiguration *AppConfigurationV1) DeleteSegmentWithContext(ctx context.Context, deleteSegmentOptions *DeleteSegmentOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSegmentOptions, "deleteSegmentOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSegmentOptions, "deleteSegmentOptions")
	if err != nil {
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
		return
	}

	response, err = appConfiguration.Service.Request(request, nil)

	return
}

// Collection : Details of the collection.
type Collection struct {
	// Collection name.
	Name *string `json:"name" validate:"required"`

	// Collection Id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection.
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

	// Number of features associated with the collection.
	FeaturesCount *int64 `json:"features_count,omitempty"`

	// Number of features associated with the collection.
	PropertiesCount *int64 `json:"properties_count,omitempty"`
}


// NewCollection : Instantiate Collection (Generic Model Constructor)
func (*AppConfigurationV1) NewCollection(name string, collectionID string) (model *Collection, err error) {
	model = &Collection{
		Name: core.StringPtr(name),
		CollectionID: core.StringPtr(collectionID),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCollection unmarshals an instance of Collection from the specified map of raw messages.
func UnmarshalCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Collection)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeatureOutput)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyOutput)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "features_count", &obj.FeaturesCount)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "properties_count", &obj.PropertiesCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CollectionList : List of all Collections.
type CollectionList struct {
	// Array of collections.
	Collections []Collection `json:"collections" validate:"required"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Response having URL of the page.
	First *PageHrefResponse `json:"first,omitempty"`

	// Response having URL of the page.
	Previous *PageHrefResponse `json:"previous,omitempty"`

	// Response having URL of the page.
	Next *PageHrefResponse `json:"next,omitempty"`

	// Response having URL of the page.
	Last *PageHrefResponse `json:"last,omitempty"`
}


// UnmarshalCollectionList unmarshals an instance of CollectionList from the specified map of raw messages.
func UnmarshalCollectionList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectionList)
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollection)
	if err != nil {
		return
	}
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CollectionLite : Details of the collection.
type CollectionLite struct {
	// Collection name.
	Name *string `json:"name" validate:"required"`

	// Collection Id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection.
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
		return
	}
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
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

// CollectionRef : CollectionRef struct
type CollectionRef struct {
	// Collection id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Name of the collection.
	Name *string `json:"name,omitempty"`
}


// NewCollectionRef : Instantiate CollectionRef (Generic Model Constructor)
func (*AppConfigurationV1) NewCollectionRef(collectionID string) (model *CollectionRef, err error) {
	model = &CollectionRef{
		CollectionID: core.StringPtr(collectionID),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalCollectionRef unmarshals an instance of CollectionRef from the specified map of raw messages.
func UnmarshalCollectionRef(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CollectionRef)
	err = core.UnmarshalPrimitive(m, "collection_id", &obj.CollectionID)
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

// CreateCollectionOptions : The CreateCollection options.
type CreateCollectionOptions struct {
	// Collection name.
	Name *string `json:"name" validate:"required"`

	// Collection Id.
	CollectionID *string `json:"collection_id" validate:"required"`

	// Collection description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection.
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests
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
func (options *CreateCollectionOptions) SetName(name string) *CreateCollectionOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetCollectionID : Allow user to set CollectionID
func (options *CreateCollectionOptions) SetCollectionID(collectionID string) *CreateCollectionOptions {
	options.CollectionID = core.StringPtr(collectionID)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateCollectionOptions) SetDescription(description string) *CreateCollectionOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateCollectionOptions) SetTags(tags string) *CreateCollectionOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateCollectionOptions) SetHeaders(param map[string]string) *CreateCollectionOptions {
	options.Headers = param
	return options
}

// CreateEnvironmentOptions : The CreateEnvironment options.
type CreateEnvironmentOptions struct {
	// Environment name.
	Name *string `json:"name" validate:"required"`

	// Environment id.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Environment description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment.
	Tags *string `json:"tags,omitempty"`

	// Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.
	ColorCode *string `json:"color_code,omitempty"`

	// Allows users to set headers on API requests
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
func (options *CreateEnvironmentOptions) SetName(name string) *CreateEnvironmentOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *CreateEnvironmentOptions) SetEnvironmentID(environmentID string) *CreateEnvironmentOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateEnvironmentOptions) SetDescription(description string) *CreateEnvironmentOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateEnvironmentOptions) SetTags(tags string) *CreateEnvironmentOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetColorCode : Allow user to set ColorCode
func (options *CreateEnvironmentOptions) SetColorCode(colorCode string) *CreateEnvironmentOptions {
	options.ColorCode = core.StringPtr(colorCode)
	return options
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

	// Feature name.
	Name *string `json:"name" validate:"required"`

	// Feature id.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Type of the feature (BOOLEAN, STRING, NUMERIC).
	Type *string `json:"type" validate:"required"`

	// Value of the feature when it is enabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	EnabledValue interface{} `json:"enabled_value" validate:"required"`

	// Value of the feature when it is disabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	DisabledValue interface{} `json:"disabled_value" validate:"required"`

	// Feature description.
	Description *string `json:"description,omitempty"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Tags associated with the feature.
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different feature flag values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified feature flag.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateFeatureOptions.Type property.
// Type of the feature (BOOLEAN, STRING, NUMERIC).
const (
	CreateFeatureOptions_Type_Boolean = "BOOLEAN"
	CreateFeatureOptions_Type_Numeric = "NUMERIC"
	CreateFeatureOptions_Type_String = "STRING"
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
func (options *CreateFeatureOptions) SetEnvironmentID(environmentID string) *CreateFeatureOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetName : Allow user to set Name
func (options *CreateFeatureOptions) SetName(name string) *CreateFeatureOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetFeatureID : Allow user to set FeatureID
func (options *CreateFeatureOptions) SetFeatureID(featureID string) *CreateFeatureOptions {
	options.FeatureID = core.StringPtr(featureID)
	return options
}

// SetType : Allow user to set Type
func (options *CreateFeatureOptions) SetType(typeVar string) *CreateFeatureOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetEnabledValue : Allow user to set EnabledValue
func (options *CreateFeatureOptions) SetEnabledValue(enabledValue interface{}) *CreateFeatureOptions {
	options.EnabledValue = enabledValue
	return options
}

// SetDisabledValue : Allow user to set DisabledValue
func (options *CreateFeatureOptions) SetDisabledValue(disabledValue interface{}) *CreateFeatureOptions {
	options.DisabledValue = disabledValue
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateFeatureOptions) SetDescription(description string) *CreateFeatureOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *CreateFeatureOptions) SetEnabled(enabled bool) *CreateFeatureOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateFeatureOptions) SetTags(tags string) *CreateFeatureOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetSegmentRules : Allow user to set SegmentRules
func (options *CreateFeatureOptions) SetSegmentRules(segmentRules []SegmentRule) *CreateFeatureOptions {
	options.SegmentRules = segmentRules
	return options
}

// SetCollections : Allow user to set Collections
func (options *CreateFeatureOptions) SetCollections(collections []CollectionRef) *CreateFeatureOptions {
	options.Collections = collections
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateFeatureOptions) SetHeaders(param map[string]string) *CreateFeatureOptions {
	options.Headers = param
	return options
}

// CreatePropertyOptions : The CreateProperty options.
type CreatePropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property id.
	PropertyID *string `json:"property_id" validate:"required"`

	// Type of the Property (BOOLEAN, STRING, NUMERIC).
	Type *string `json:"type" validate:"required"`

	// Value of the Property. The value can be Boolean, String or a Numeric value as per the `type` attribute.
	Value interface{} `json:"value" validate:"required"`

	// Property description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the property.
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreatePropertyOptions.Type property.
// Type of the Property (BOOLEAN, STRING, NUMERIC).
const (
	CreatePropertyOptions_Type_Boolean = "BOOLEAN"
	CreatePropertyOptions_Type_Numeric = "NUMERIC"
	CreatePropertyOptions_Type_String = "STRING"
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
func (options *CreatePropertyOptions) SetEnvironmentID(environmentID string) *CreatePropertyOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetName : Allow user to set Name
func (options *CreatePropertyOptions) SetName(name string) *CreatePropertyOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetPropertyID : Allow user to set PropertyID
func (options *CreatePropertyOptions) SetPropertyID(propertyID string) *CreatePropertyOptions {
	options.PropertyID = core.StringPtr(propertyID)
	return options
}

// SetType : Allow user to set Type
func (options *CreatePropertyOptions) SetType(typeVar string) *CreatePropertyOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetValue : Allow user to set Value
func (options *CreatePropertyOptions) SetValue(value interface{}) *CreatePropertyOptions {
	options.Value = value
	return options
}

// SetDescription : Allow user to set Description
func (options *CreatePropertyOptions) SetDescription(description string) *CreatePropertyOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreatePropertyOptions) SetTags(tags string) *CreatePropertyOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetSegmentRules : Allow user to set SegmentRules
func (options *CreatePropertyOptions) SetSegmentRules(segmentRules []SegmentRule) *CreatePropertyOptions {
	options.SegmentRules = segmentRules
	return options
}

// SetCollections : Allow user to set Collections
func (options *CreatePropertyOptions) SetCollections(collections []CollectionRef) *CreatePropertyOptions {
	options.Collections = collections
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreatePropertyOptions) SetHeaders(param map[string]string) *CreatePropertyOptions {
	options.Headers = param
	return options
}

// CreateSegmentOptions : The CreateSegment options.
type CreateSegmentOptions struct {
	// Segment name.
	Name *string `json:"name,omitempty"`

	// Segment id.
	SegmentID *string `json:"segment_id,omitempty"`

	// Segment description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the segments.
	Tags *string `json:"tags,omitempty"`

	// List of rules that determine if the entity is part of the segment.
	Rules []Rule `json:"rules,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateSegmentOptions : Instantiate CreateSegmentOptions
func (*AppConfigurationV1) NewCreateSegmentOptions() *CreateSegmentOptions {
	return &CreateSegmentOptions{}
}

// SetName : Allow user to set Name
func (options *CreateSegmentOptions) SetName(name string) *CreateSegmentOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetSegmentID : Allow user to set SegmentID
func (options *CreateSegmentOptions) SetSegmentID(segmentID string) *CreateSegmentOptions {
	options.SegmentID = core.StringPtr(segmentID)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateSegmentOptions) SetDescription(description string) *CreateSegmentOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *CreateSegmentOptions) SetTags(tags string) *CreateSegmentOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetRules : Allow user to set Rules
func (options *CreateSegmentOptions) SetRules(rules []Rule) *CreateSegmentOptions {
	options.Rules = rules
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSegmentOptions) SetHeaders(param map[string]string) *CreateSegmentOptions {
	options.Headers = param
	return options
}

// DeleteCollectionOptions : The DeleteCollection options.
type DeleteCollectionOptions struct {
	// Collection Id of the collection.
	CollectionID *string `json:"collection_id" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteCollectionOptions : Instantiate DeleteCollectionOptions
func (*AppConfigurationV1) NewDeleteCollectionOptions(collectionID string) *DeleteCollectionOptions {
	return &DeleteCollectionOptions{
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetCollectionID : Allow user to set CollectionID
func (options *DeleteCollectionOptions) SetCollectionID(collectionID string) *DeleteCollectionOptions {
	options.CollectionID = core.StringPtr(collectionID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEnvironmentOptions : Instantiate DeleteEnvironmentOptions
func (*AppConfigurationV1) NewDeleteEnvironmentOptions(environmentID string) *DeleteEnvironmentOptions {
	return &DeleteEnvironmentOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *DeleteEnvironmentOptions) SetEnvironmentID(environmentID string) *DeleteEnvironmentOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
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

	// Allows users to set headers on API requests
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
func (options *DeleteFeatureOptions) SetEnvironmentID(environmentID string) *DeleteFeatureOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetFeatureID : Allow user to set FeatureID
func (options *DeleteFeatureOptions) SetFeatureID(featureID string) *DeleteFeatureOptions {
	options.FeatureID = core.StringPtr(featureID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteFeatureOptions) SetHeaders(param map[string]string) *DeleteFeatureOptions {
	options.Headers = param
	return options
}

// DeletePropertyOptions : The DeleteProperty options.
type DeletePropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Allows users to set headers on API requests
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
func (options *DeletePropertyOptions) SetEnvironmentID(environmentID string) *DeletePropertyOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetPropertyID : Allow user to set PropertyID
func (options *DeletePropertyOptions) SetPropertyID(propertyID string) *DeletePropertyOptions {
	options.PropertyID = core.StringPtr(propertyID)
	return options
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSegmentOptions : Instantiate DeleteSegmentOptions
func (*AppConfigurationV1) NewDeleteSegmentOptions(segmentID string) *DeleteSegmentOptions {
	return &DeleteSegmentOptions{
		SegmentID: core.StringPtr(segmentID),
	}
}

// SetSegmentID : Allow user to set SegmentID
func (options *DeleteSegmentOptions) SetSegmentID(segmentID string) *DeleteSegmentOptions {
	options.SegmentID = core.StringPtr(segmentID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSegmentOptions) SetHeaders(param map[string]string) *DeleteSegmentOptions {
	options.Headers = param
	return options
}

// Environment : Details of the environment.
type Environment struct {
	// Environment name.
	Name *string `json:"name" validate:"required"`

	// Environment id.
	EnvironmentID *string `json:"environment_id" validate:"required"`

	// Environment description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment.
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
}


// NewEnvironment : Instantiate Environment (Generic Model Constructor)
func (*AppConfigurationV1) NewEnvironment(name string, environmentID string) (model *Environment, err error) {
	model = &Environment{
		Name: core.StringPtr(name),
		EnvironmentID: core.StringPtr(environmentID),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalEnvironment unmarshals an instance of Environment from the specified map of raw messages.
func UnmarshalEnvironment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Environment)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "environment_id", &obj.EnvironmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "color_code", &obj.ColorCode)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeatureOutput)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyOutput)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EnvironmentList : List of all environments.
type EnvironmentList struct {
	// Array of environments.
	Environments []Environment `json:"environments" validate:"required"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Response having URL of the page.
	First *PageHrefResponse `json:"first,omitempty"`

	// Response having URL of the page.
	Previous *PageHrefResponse `json:"previous,omitempty"`

	// Response having URL of the page.
	Next *PageHrefResponse `json:"next,omitempty"`

	// Response having URL of the page.
	Last *PageHrefResponse `json:"last,omitempty"`
}


// UnmarshalEnvironmentList unmarshals an instance of EnvironmentList from the specified map of raw messages.
func UnmarshalEnvironmentList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EnvironmentList)
	err = core.UnmarshalModel(m, "environments", &obj.Environments, UnmarshalEnvironment)
	if err != nil {
		return
	}
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Feature : Details of the feature.
type Feature struct {
	// Feature name.
	Name *string `json:"name" validate:"required"`

	// Feature id.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Feature description.
	Description *string `json:"description,omitempty"`

	// Type of the feature (BOOLEAN, STRING, NUMERIC).
	Type *string `json:"type" validate:"required"`

	// Value of the feature when it is enabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	EnabledValue interface{} `json:"enabled_value" validate:"required"`

	// Value of the feature when it is disabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	DisabledValue interface{} `json:"disabled_value" validate:"required"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Tags associated with the feature.
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different feature flag values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// Denotes if the targeting rules are specified for the feature flag.
	SegmentExists *bool `json:"segment_exists,omitempty"`

	// List of collection id representing the collections that are associated with the specified feature flag.
	Collections []CollectionRef `json:"collections,omitempty"`

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
// Type of the feature (BOOLEAN, STRING, NUMERIC).
const (
	Feature_Type_Boolean = "BOOLEAN"
	Feature_Type_Numeric = "NUMERIC"
	Feature_Type_String = "STRING"
)


// NewFeature : Instantiate Feature (Generic Model Constructor)
func (*AppConfigurationV1) NewFeature(name string, featureID string, typeVar string, enabledValue interface{}, disabledValue interface{}) (model *Feature, err error) {
	model = &Feature{
		Name: core.StringPtr(name),
		FeatureID: core.StringPtr(featureID),
		Type: core.StringPtr(typeVar),
		EnabledValue: enabledValue,
		DisabledValue: disabledValue,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFeature unmarshals an instance of Feature from the specified map of raw messages.
func UnmarshalFeature(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Feature)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "feature_id", &obj.FeatureID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled_value", &obj.EnabledValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled_value", &obj.DisabledValue)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "segment_rules", &obj.SegmentRules, UnmarshalSegmentRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_exists", &obj.SegmentExists)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollectionRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluation_time", &obj.EvaluationTime)
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

// FeatureOutput : Feature object.
type FeatureOutput struct {
	// Feature id.
	FeatureID *string `json:"feature_id" validate:"required"`

	// Feature name.
	Name *string `json:"name" validate:"required"`
}


// NewFeatureOutput : Instantiate FeatureOutput (Generic Model Constructor)
func (*AppConfigurationV1) NewFeatureOutput(featureID string, name string) (model *FeatureOutput, err error) {
	model = &FeatureOutput{
		FeatureID: core.StringPtr(featureID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalFeatureOutput unmarshals an instance of FeatureOutput from the specified map of raw messages.
func UnmarshalFeatureOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FeatureOutput)
	err = core.UnmarshalPrimitive(m, "feature_id", &obj.FeatureID)
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

// FeaturesList : List of all features.
type FeaturesList struct {
	// Array of Features.
	Features []Feature `json:"features" validate:"required"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Response having URL of the page.
	First *PageHrefResponse `json:"first,omitempty"`

	// Response having URL of the page.
	Previous *PageHrefResponse `json:"previous,omitempty"`

	// Response having URL of the page.
	Next *PageHrefResponse `json:"next,omitempty"`

	// Response having URL of the page.
	Last *PageHrefResponse `json:"last,omitempty"`
}


// UnmarshalFeaturesList unmarshals an instance of FeaturesList from the specified map of raw messages.
func UnmarshalFeaturesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(FeaturesList)
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeature)
	if err != nil {
		return
	}
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetCollectionOptions : The GetCollection options.
type GetCollectionOptions struct {
	// Collection Id of the collection.
	CollectionID *string `json:"collection_id" validate:"required,ne="`

	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Include feature and property details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetCollectionOptions.Include property.
const (
	GetCollectionOptions_Include_Features = "features"
	GetCollectionOptions_Include_Properties = "properties"
)

// NewGetCollectionOptions : Instantiate GetCollectionOptions
func (*AppConfigurationV1) NewGetCollectionOptions(collectionID string) *GetCollectionOptions {
	return &GetCollectionOptions{
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetCollectionID : Allow user to set CollectionID
func (options *GetCollectionOptions) SetCollectionID(collectionID string) *GetCollectionOptions {
	options.CollectionID = core.StringPtr(collectionID)
	return options
}

// SetExpand : Allow user to set Expand
func (options *GetCollectionOptions) SetExpand(expand bool) *GetCollectionOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetInclude : Allow user to set Include
func (options *GetCollectionOptions) SetInclude(include []string) *GetCollectionOptions {
	options.Include = include
	return options
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

	// Include feature and property details in the response.
	Include []string `json:"include,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetEnvironmentOptions.Include property.
const (
	GetEnvironmentOptions_Include_Features = "features"
	GetEnvironmentOptions_Include_Properties = "properties"
)

// NewGetEnvironmentOptions : Instantiate GetEnvironmentOptions
func (*AppConfigurationV1) NewGetEnvironmentOptions(environmentID string) *GetEnvironmentOptions {
	return &GetEnvironmentOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *GetEnvironmentOptions) SetEnvironmentID(environmentID string) *GetEnvironmentOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetExpand : Allow user to set Expand
func (options *GetEnvironmentOptions) SetExpand(expand bool) *GetEnvironmentOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetInclude : Allow user to set Include
func (options *GetEnvironmentOptions) SetInclude(include []string) *GetEnvironmentOptions {
	options.Include = include
	return options
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

	// Include the associated collections in the response.
	Include *string `json:"include,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetFeatureOptions.Include property.
// Include the associated collections in the response.
const (
	GetFeatureOptions_Include_Collections = "collections"
)

// NewGetFeatureOptions : Instantiate GetFeatureOptions
func (*AppConfigurationV1) NewGetFeatureOptions(environmentID string, featureID string) *GetFeatureOptions {
	return &GetFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *GetFeatureOptions) SetEnvironmentID(environmentID string) *GetFeatureOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetFeatureID : Allow user to set FeatureID
func (options *GetFeatureOptions) SetFeatureID(featureID string) *GetFeatureOptions {
	options.FeatureID = core.StringPtr(featureID)
	return options
}

// SetInclude : Allow user to set Include
func (options *GetFeatureOptions) SetInclude(include string) *GetFeatureOptions {
	options.Include = core.StringPtr(include)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetFeatureOptions) SetHeaders(param map[string]string) *GetFeatureOptions {
	options.Headers = param
	return options
}

// GetPropertyOptions : The GetProperty options.
type GetPropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Include the associated collections in the response.
	Include *string `json:"include,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the GetPropertyOptions.Include property.
// Include the associated collections in the response.
const (
	GetPropertyOptions_Include_Collections = "collections"
)

// NewGetPropertyOptions : Instantiate GetPropertyOptions
func (*AppConfigurationV1) NewGetPropertyOptions(environmentID string, propertyID string) *GetPropertyOptions {
	return &GetPropertyOptions{
		EnvironmentID: core.StringPtr(environmentID),
		PropertyID: core.StringPtr(propertyID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *GetPropertyOptions) SetEnvironmentID(environmentID string) *GetPropertyOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetPropertyID : Allow user to set PropertyID
func (options *GetPropertyOptions) SetPropertyID(propertyID string) *GetPropertyOptions {
	options.PropertyID = core.StringPtr(propertyID)
	return options
}

// SetInclude : Allow user to set Include
func (options *GetPropertyOptions) SetInclude(include string) *GetPropertyOptions {
	options.Include = core.StringPtr(include)
	return options
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

	// Allows users to set headers on API requests
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
func (options *GetSegmentOptions) SetSegmentID(segmentID string) *GetSegmentOptions {
	options.SegmentID = core.StringPtr(segmentID)
	return options
}

// SetInclude : Allow user to set Include
func (options *GetSegmentOptions) SetInclude(include []string) *GetSegmentOptions {
	options.Include = include
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetSegmentOptions) SetHeaders(param map[string]string) *GetSegmentOptions {
	options.Headers = param
	return options
}

// ListCollectionsOptions : The ListCollections options.
type ListCollectionsOptions struct {
	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the collection details based on the specified attribute.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Filter collections by a list of comma separated features.
	Features []string `json:"features,omitempty"`

	// Filter collections by a list of comma separated properties.
	Properties []string `json:"properties,omitempty"`

	// Include feature and property details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListCollectionsOptions.Sort property.
// Sort the collection details based on the specified attribute.
const (
	ListCollectionsOptions_Sort_CollectionID = "collection_id"
	ListCollectionsOptions_Sort_CreatedTime = "created_time"
	ListCollectionsOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListCollectionsOptions.Include property.
const (
	ListCollectionsOptions_Include_Features = "features"
	ListCollectionsOptions_Include_Properties = "properties"
)

// NewListCollectionsOptions : Instantiate ListCollectionsOptions
func (*AppConfigurationV1) NewListCollectionsOptions() *ListCollectionsOptions {
	return &ListCollectionsOptions{}
}

// SetExpand : Allow user to set Expand
func (options *ListCollectionsOptions) SetExpand(expand bool) *ListCollectionsOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListCollectionsOptions) SetSort(sort string) *ListCollectionsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetTags : Allow user to set Tags
func (options *ListCollectionsOptions) SetTags(tags string) *ListCollectionsOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetFeatures : Allow user to set Features
func (options *ListCollectionsOptions) SetFeatures(features []string) *ListCollectionsOptions {
	options.Features = features
	return options
}

// SetProperties : Allow user to set Properties
func (options *ListCollectionsOptions) SetProperties(properties []string) *ListCollectionsOptions {
	options.Properties = properties
	return options
}

// SetInclude : Allow user to set Include
func (options *ListCollectionsOptions) SetInclude(include []string) *ListCollectionsOptions {
	options.Include = include
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListCollectionsOptions) SetLimit(limit int64) *ListCollectionsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListCollectionsOptions) SetOffset(offset int64) *ListCollectionsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
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

	// Sort the environment details based on the specified attribute.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Include feature and property details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListEnvironmentsOptions.Sort property.
// Sort the environment details based on the specified attribute.
const (
	ListEnvironmentsOptions_Sort_CreatedTime = "created_time"
	ListEnvironmentsOptions_Sort_EnvironmentID = "environment_id"
	ListEnvironmentsOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListEnvironmentsOptions.Include property.
const (
	ListEnvironmentsOptions_Include_Features = "features"
	ListEnvironmentsOptions_Include_Properties = "properties"
)

// NewListEnvironmentsOptions : Instantiate ListEnvironmentsOptions
func (*AppConfigurationV1) NewListEnvironmentsOptions() *ListEnvironmentsOptions {
	return &ListEnvironmentsOptions{}
}

// SetExpand : Allow user to set Expand
func (options *ListEnvironmentsOptions) SetExpand(expand bool) *ListEnvironmentsOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListEnvironmentsOptions) SetSort(sort string) *ListEnvironmentsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetTags : Allow user to set Tags
func (options *ListEnvironmentsOptions) SetTags(tags string) *ListEnvironmentsOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetInclude : Allow user to set Include
func (options *ListEnvironmentsOptions) SetInclude(include []string) *ListEnvironmentsOptions {
	options.Include = include
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListEnvironmentsOptions) SetLimit(limit int64) *ListEnvironmentsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListEnvironmentsOptions) SetOffset(offset int64) *ListEnvironmentsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
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

	// Sort the feature details based on the specified attribute.
	Sort *string `json:"sort,omitempty"`

	// Filter the resources to be returned based on the associated tags. Specify the parameter as a list of comma separated
	// tags. Returns resources associated with any of the specified tags.
	Tags *string `json:"tags,omitempty"`

	// Filter features by a list of comma separated collections.
	Collections []string `json:"collections,omitempty"`

	// Filter features by a list of comma separated segments.
	Segments []string `json:"segments,omitempty"`

	// Include the associated collections or targeting rules details in the response.
	Include []string `json:"include,omitempty"`

	// The number of records to retrieve. By default, the list operation return the first 10 records. To retrieve different
	// set of records, use `limit` with `offset` to page through the available records.
	Limit *int64 `json:"limit,omitempty"`

	// The number of records to skip. By specifying `offset`, you retrieve a subset of items that starts with the `offset`
	// value. Use `offset` with `limit` to page through the available records.
	Offset *int64 `json:"offset,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListFeaturesOptions.Sort property.
// Sort the feature details based on the specified attribute.
const (
	ListFeaturesOptions_Sort_CreatedTime = "created_time"
	ListFeaturesOptions_Sort_FeatureID = "feature_id"
	ListFeaturesOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListFeaturesOptions.Include property.
const (
	ListFeaturesOptions_Include_Collections = "collections"
	ListFeaturesOptions_Include_Rules = " rules"
)

// NewListFeaturesOptions : Instantiate ListFeaturesOptions
func (*AppConfigurationV1) NewListFeaturesOptions(environmentID string) *ListFeaturesOptions {
	return &ListFeaturesOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *ListFeaturesOptions) SetEnvironmentID(environmentID string) *ListFeaturesOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetExpand : Allow user to set Expand
func (options *ListFeaturesOptions) SetExpand(expand bool) *ListFeaturesOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListFeaturesOptions) SetSort(sort string) *ListFeaturesOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetTags : Allow user to set Tags
func (options *ListFeaturesOptions) SetTags(tags string) *ListFeaturesOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetCollections : Allow user to set Collections
func (options *ListFeaturesOptions) SetCollections(collections []string) *ListFeaturesOptions {
	options.Collections = collections
	return options
}

// SetSegments : Allow user to set Segments
func (options *ListFeaturesOptions) SetSegments(segments []string) *ListFeaturesOptions {
	options.Segments = segments
	return options
}

// SetInclude : Allow user to set Include
func (options *ListFeaturesOptions) SetInclude(include []string) *ListFeaturesOptions {
	options.Include = include
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListFeaturesOptions) SetLimit(limit int64) *ListFeaturesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListFeaturesOptions) SetOffset(offset int64) *ListFeaturesOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListFeaturesOptions) SetHeaders(param map[string]string) *ListFeaturesOptions {
	options.Headers = param
	return options
}

// ListPropertiesOptions : The ListProperties options.
type ListPropertiesOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// If set to `true`, returns expanded view of the resource details.
	Expand *bool `json:"expand,omitempty"`

	// Sort the property details based on the specified attribute.
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListPropertiesOptions.Sort property.
// Sort the property details based on the specified attribute.
const (
	ListPropertiesOptions_Sort_CreatedTime = "created_time"
	ListPropertiesOptions_Sort_PropertyID = "property_id"
	ListPropertiesOptions_Sort_UpdatedTime = "updated_time"
)

// Constants associated with the ListPropertiesOptions.Include property.
const (
	ListPropertiesOptions_Include_Collections = "collections"
	ListPropertiesOptions_Include_Rules = " rules"
)

// NewListPropertiesOptions : Instantiate ListPropertiesOptions
func (*AppConfigurationV1) NewListPropertiesOptions(environmentID string) *ListPropertiesOptions {
	return &ListPropertiesOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *ListPropertiesOptions) SetEnvironmentID(environmentID string) *ListPropertiesOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetExpand : Allow user to set Expand
func (options *ListPropertiesOptions) SetExpand(expand bool) *ListPropertiesOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListPropertiesOptions) SetSort(sort string) *ListPropertiesOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetTags : Allow user to set Tags
func (options *ListPropertiesOptions) SetTags(tags string) *ListPropertiesOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetCollections : Allow user to set Collections
func (options *ListPropertiesOptions) SetCollections(collections []string) *ListPropertiesOptions {
	options.Collections = collections
	return options
}

// SetSegments : Allow user to set Segments
func (options *ListPropertiesOptions) SetSegments(segments []string) *ListPropertiesOptions {
	options.Segments = segments
	return options
}

// SetInclude : Allow user to set Include
func (options *ListPropertiesOptions) SetInclude(include []string) *ListPropertiesOptions {
	options.Include = include
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListPropertiesOptions) SetLimit(limit int64) *ListPropertiesOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListPropertiesOptions) SetOffset(offset int64) *ListPropertiesOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
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

	// Sort the segment details based on the specified attribute.
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

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListSegmentsOptions.Sort property.
// Sort the segment details based on the specified attribute.
const (
	ListSegmentsOptions_Sort_CreatedTime = "created_time"
	ListSegmentsOptions_Sort_SegmentID = "segment_id"
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
func (options *ListSegmentsOptions) SetExpand(expand bool) *ListSegmentsOptions {
	options.Expand = core.BoolPtr(expand)
	return options
}

// SetSort : Allow user to set Sort
func (options *ListSegmentsOptions) SetSort(sort string) *ListSegmentsOptions {
	options.Sort = core.StringPtr(sort)
	return options
}

// SetTags : Allow user to set Tags
func (options *ListSegmentsOptions) SetTags(tags string) *ListSegmentsOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetInclude : Allow user to set Include
func (options *ListSegmentsOptions) SetInclude(include string) *ListSegmentsOptions {
	options.Include = core.StringPtr(include)
	return options
}

// SetLimit : Allow user to set Limit
func (options *ListSegmentsOptions) SetLimit(limit int64) *ListSegmentsOptions {
	options.Limit = core.Int64Ptr(limit)
	return options
}

// SetOffset : Allow user to set Offset
func (options *ListSegmentsOptions) SetOffset(offset int64) *ListSegmentsOptions {
	options.Offset = core.Int64Ptr(offset)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListSegmentsOptions) SetHeaders(param map[string]string) *ListSegmentsOptions {
	options.Headers = param
	return options
}

// PageHrefResponse : Response having URL of the page.
type PageHrefResponse struct {
	// URL of the response.
	Href *string `json:"href" validate:"required"`
}


// UnmarshalPageHrefResponse unmarshals an instance of PageHrefResponse from the specified map of raw messages.
func UnmarshalPageHrefResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PageHrefResponse)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// PropertiesList : List of all properties.
type PropertiesList struct {
	// Array of properties.
	Properties []Property `json:"properties" validate:"required"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Response having URL of the page.
	First *PageHrefResponse `json:"first,omitempty"`

	// Response having URL of the page.
	Previous *PageHrefResponse `json:"previous,omitempty"`

	// Response having URL of the page.
	Next *PageHrefResponse `json:"next,omitempty"`

	// Response having URL of the page.
	Last *PageHrefResponse `json:"last,omitempty"`
}


// UnmarshalPropertiesList unmarshals an instance of PropertiesList from the specified map of raw messages.
func UnmarshalPropertiesList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PropertiesList)
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalProperty)
	if err != nil {
		return
	}
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Property : Details of the property.
type Property struct {
	// Property name.
	Name *string `json:"name" validate:"required"`

	// Property id.
	PropertyID *string `json:"property_id" validate:"required"`

	// Property description.
	Description *string `json:"description,omitempty"`

	// Type of the Property (BOOLEAN, STRING, NUMERIC).
	Type *string `json:"type" validate:"required"`

	// Value of the Property. The value can be Boolean, String or a Numeric value as per the `type` attribute.
	Value interface{} `json:"value" validate:"required"`

	// Tags associated with the property.
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
// Type of the Property (BOOLEAN, STRING, NUMERIC).
const (
	Property_Type_Boolean = "BOOLEAN"
	Property_Type_Numeric = "NUMERIC"
	Property_Type_String = "STRING"
)


// NewProperty : Instantiate Property (Generic Model Constructor)
func (*AppConfigurationV1) NewProperty(name string, propertyID string, typeVar string, value interface{}) (model *Property, err error) {
	model = &Property{
		Name: core.StringPtr(name),
		PropertyID: core.StringPtr(propertyID),
		Type: core.StringPtr(typeVar),
		Value: value,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalProperty unmarshals an instance of Property from the specified map of raw messages.
func UnmarshalProperty(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Property)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "property_id", &obj.PropertyID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "segment_rules", &obj.SegmentRules, UnmarshalSegmentRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_exists", &obj.SegmentExists)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "collections", &obj.Collections, UnmarshalCollectionRef)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "evaluation_time", &obj.EvaluationTime)
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

// PropertyOutput : Property object.
type PropertyOutput struct {
	// Property id.
	PropertyID *string `json:"property_id" validate:"required"`

	// Property name.
	Name *string `json:"name" validate:"required"`
}


// NewPropertyOutput : Instantiate PropertyOutput (Generic Model Constructor)
func (*AppConfigurationV1) NewPropertyOutput(propertyID string, name string) (model *PropertyOutput, err error) {
	model = &PropertyOutput{
		PropertyID: core.StringPtr(propertyID),
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalPropertyOutput unmarshals an instance of PropertyOutput from the specified map of raw messages.
func UnmarshalPropertyOutput(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(PropertyOutput)
	err = core.UnmarshalPrimitive(m, "property_id", &obj.PropertyID)
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

// Rule : Rule is used to determine if the entity belongs to the segment during feature / property evaluation. An entity is
// identified by a unique identifier and the attributes that it defines. Any feature flag and property value evaluation
// is performed in the context of an entity when it is targeted to segments.
type Rule struct {
	// Attribute name.
	AttributeName *string `json:"attribute_name" validate:"required"`

	// Operator to be used for the evaluation if the entity is part of the segment.
	Operator *string `json:"operator" validate:"required"`

	// List of values. Entities matching any of the given values will be considered to be part of the segment.
	Values []string `json:"values" validate:"required"`
}

// Constants associated with the Rule.Operator property.
// Operator to be used for the evaluation if the entity is part of the segment.
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
func (*AppConfigurationV1) NewRule(attributeName string, operator string, values []string) (model *Rule, err error) {
	model = &Rule{
		AttributeName: core.StringPtr(attributeName),
		Operator: core.StringPtr(operator),
		Values: values,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRule unmarshals an instance of Rule from the specified map of raw messages.
func UnmarshalRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Rule)
	err = core.UnmarshalPrimitive(m, "attribute_name", &obj.AttributeName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "operator", &obj.Operator)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "values", &obj.Values)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Segment : Details of the segment.
type Segment struct {
	// Segment name.
	Name *string `json:"name" validate:"required"`

	// Segment id.
	SegmentID *string `json:"segment_id" validate:"required"`

	// Segment description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the segments.
	Tags *string `json:"tags,omitempty"`

	// List of rules that determine if the entity is part of the segment.
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
func (*AppConfigurationV1) NewSegment(name string, segmentID string, rules []Rule) (model *Segment, err error) {
	model = &Segment{
		Name: core.StringPtr(name),
		SegmentID: core.StringPtr(segmentID),
		Rules: rules,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalSegment unmarshals an instance of Segment from the specified map of raw messages.
func UnmarshalSegment(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Segment)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "segment_id", &obj.SegmentID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tags", &obj.Tags)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalRule)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_time", &obj.CreatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "updated_time", &obj.UpdatedTime)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "features", &obj.Features, UnmarshalFeatureOutput)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "properties", &obj.Properties, UnmarshalPropertyOutput)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SegmentRule : SegmentRule struct
type SegmentRule struct {
	// The list of targetted segments.
	Rules []TargetSegments `json:"rules" validate:"required"`

	// Value to be used for evaluation for this rule. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	Value interface{} `json:"value" validate:"required"`

	// Order of the rule, used during evaluation. The evaluation is performed in the order defined and the value associated
	// with the first matching rule is used for evaluation.
	Order *int64 `json:"order" validate:"required"`
}


// NewSegmentRule : Instantiate SegmentRule (Generic Model Constructor)
func (*AppConfigurationV1) NewSegmentRule(rules []TargetSegments, value interface{}, order int64) (model *SegmentRule, err error) {
	model = &SegmentRule{
		Rules: rules,
		Value: value,
		Order: core.Int64Ptr(order),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalSegmentRule unmarshals an instance of SegmentRule from the specified map of raw messages.
func UnmarshalSegmentRule(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SegmentRule)
	err = core.UnmarshalModel(m, "rules", &obj.Rules, UnmarshalTargetSegments)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "value", &obj.Value)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "order", &obj.Order)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SegmentsList : List of all segments.
type SegmentsList struct {
	// Array of Segments.
	Segments []Segment `json:"segments" validate:"required"`

	// Number of records returned.
	Limit *int64 `json:"limit" validate:"required"`

	// Skipped number of records.
	Offset *int64 `json:"offset" validate:"required"`

	// Total number of records.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// Response having URL of the page.
	First *PageHrefResponse `json:"first,omitempty"`

	// Response having URL of the page.
	Previous *PageHrefResponse `json:"previous,omitempty"`

	// Response having URL of the page.
	Next *PageHrefResponse `json:"next,omitempty"`

	// Response having URL of the page.
	Last *PageHrefResponse `json:"last,omitempty"`
}


// UnmarshalSegmentsList unmarshals an instance of SegmentsList from the specified map of raw messages.
func UnmarshalSegmentsList(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SegmentsList)
	err = core.UnmarshalModel(m, "segments", &obj.Segments, UnmarshalSegment)
	if err != nil {
		return
	}
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
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalPageHrefResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "last", &obj.Last, UnmarshalPageHrefResponse)
	if err != nil {
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
func (*AppConfigurationV1) NewTargetSegments(segments []string) (model *TargetSegments, err error) {
	model = &TargetSegments{
		Segments: segments,
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalTargetSegments unmarshals an instance of TargetSegments from the specified map of raw messages.
func UnmarshalTargetSegments(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TargetSegments)
	err = core.UnmarshalPrimitive(m, "segments", &obj.Segments)
	if err != nil {
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
	Enabled *bool `json:"enabled,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewToggleFeatureOptions : Instantiate ToggleFeatureOptions
func (*AppConfigurationV1) NewToggleFeatureOptions(environmentID string, featureID string) *ToggleFeatureOptions {
	return &ToggleFeatureOptions{
		EnvironmentID: core.StringPtr(environmentID),
		FeatureID: core.StringPtr(featureID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *ToggleFeatureOptions) SetEnvironmentID(environmentID string) *ToggleFeatureOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetFeatureID : Allow user to set FeatureID
func (options *ToggleFeatureOptions) SetFeatureID(featureID string) *ToggleFeatureOptions {
	options.FeatureID = core.StringPtr(featureID)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *ToggleFeatureOptions) SetEnabled(enabled bool) *ToggleFeatureOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
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

	// Collection name.
	Name *string `json:"name,omitempty"`

	// Description of the collection.
	Description *string `json:"description,omitempty"`

	// Tags associated with the collection.
	Tags *string `json:"tags,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateCollectionOptions : Instantiate UpdateCollectionOptions
func (*AppConfigurationV1) NewUpdateCollectionOptions(collectionID string) *UpdateCollectionOptions {
	return &UpdateCollectionOptions{
		CollectionID: core.StringPtr(collectionID),
	}
}

// SetCollectionID : Allow user to set CollectionID
func (options *UpdateCollectionOptions) SetCollectionID(collectionID string) *UpdateCollectionOptions {
	options.CollectionID = core.StringPtr(collectionID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateCollectionOptions) SetName(name string) *UpdateCollectionOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateCollectionOptions) SetDescription(description string) *UpdateCollectionOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateCollectionOptions) SetTags(tags string) *UpdateCollectionOptions {
	options.Tags = core.StringPtr(tags)
	return options
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

	// Environment name.
	Name *string `json:"name,omitempty"`

	// Environment description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the environment.
	Tags *string `json:"tags,omitempty"`

	// Color code to distinguish the environment. The Hex code for the color. For example `#FF0000` for `red`.
	ColorCode *string `json:"color_code,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateEnvironmentOptions : Instantiate UpdateEnvironmentOptions
func (*AppConfigurationV1) NewUpdateEnvironmentOptions(environmentID string) *UpdateEnvironmentOptions {
	return &UpdateEnvironmentOptions{
		EnvironmentID: core.StringPtr(environmentID),
	}
}

// SetEnvironmentID : Allow user to set EnvironmentID
func (options *UpdateEnvironmentOptions) SetEnvironmentID(environmentID string) *UpdateEnvironmentOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateEnvironmentOptions) SetName(name string) *UpdateEnvironmentOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateEnvironmentOptions) SetDescription(description string) *UpdateEnvironmentOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateEnvironmentOptions) SetTags(tags string) *UpdateEnvironmentOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetColorCode : Allow user to set ColorCode
func (options *UpdateEnvironmentOptions) SetColorCode(colorCode string) *UpdateEnvironmentOptions {
	options.ColorCode = core.StringPtr(colorCode)
	return options
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

	// Feature name.
	Name *string `json:"name,omitempty"`

	// Feature description.
	Description *string `json:"description,omitempty"`

	// Value of the feature when it is enabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	EnabledValue interface{} `json:"enabled_value,omitempty"`

	// Value of the feature when it is disabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	DisabledValue interface{} `json:"disabled_value,omitempty"`

	// The state of the feature flag.
	Enabled *bool `json:"enabled,omitempty"`

	// Tags associated with the feature.
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests
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
func (options *UpdateFeatureOptions) SetEnvironmentID(environmentID string) *UpdateFeatureOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetFeatureID : Allow user to set FeatureID
func (options *UpdateFeatureOptions) SetFeatureID(featureID string) *UpdateFeatureOptions {
	options.FeatureID = core.StringPtr(featureID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateFeatureOptions) SetName(name string) *UpdateFeatureOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateFeatureOptions) SetDescription(description string) *UpdateFeatureOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetEnabledValue : Allow user to set EnabledValue
func (options *UpdateFeatureOptions) SetEnabledValue(enabledValue interface{}) *UpdateFeatureOptions {
	options.EnabledValue = enabledValue
	return options
}

// SetDisabledValue : Allow user to set DisabledValue
func (options *UpdateFeatureOptions) SetDisabledValue(disabledValue interface{}) *UpdateFeatureOptions {
	options.DisabledValue = disabledValue
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *UpdateFeatureOptions) SetEnabled(enabled bool) *UpdateFeatureOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateFeatureOptions) SetTags(tags string) *UpdateFeatureOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetSegmentRules : Allow user to set SegmentRules
func (options *UpdateFeatureOptions) SetSegmentRules(segmentRules []SegmentRule) *UpdateFeatureOptions {
	options.SegmentRules = segmentRules
	return options
}

// SetCollections : Allow user to set Collections
func (options *UpdateFeatureOptions) SetCollections(collections []CollectionRef) *UpdateFeatureOptions {
	options.Collections = collections
	return options
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

	// Feature name.
	Name *string `json:"name,omitempty"`

	// Feature description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the feature.
	Tags *string `json:"tags,omitempty"`

	// Value of the feature when it is enabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	EnabledValue interface{} `json:"enabled_value,omitempty"`

	// Value of the feature when it is disabled. The value can be Boolean, String or a Numeric value as per the `type`
	// attribute.
	DisabledValue interface{} `json:"disabled_value,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// Allows users to set headers on API requests
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
func (options *UpdateFeatureValuesOptions) SetEnvironmentID(environmentID string) *UpdateFeatureValuesOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetFeatureID : Allow user to set FeatureID
func (options *UpdateFeatureValuesOptions) SetFeatureID(featureID string) *UpdateFeatureValuesOptions {
	options.FeatureID = core.StringPtr(featureID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateFeatureValuesOptions) SetName(name string) *UpdateFeatureValuesOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateFeatureValuesOptions) SetDescription(description string) *UpdateFeatureValuesOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateFeatureValuesOptions) SetTags(tags string) *UpdateFeatureValuesOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetEnabledValue : Allow user to set EnabledValue
func (options *UpdateFeatureValuesOptions) SetEnabledValue(enabledValue interface{}) *UpdateFeatureValuesOptions {
	options.EnabledValue = enabledValue
	return options
}

// SetDisabledValue : Allow user to set DisabledValue
func (options *UpdateFeatureValuesOptions) SetDisabledValue(disabledValue interface{}) *UpdateFeatureValuesOptions {
	options.DisabledValue = disabledValue
	return options
}

// SetSegmentRules : Allow user to set SegmentRules
func (options *UpdateFeatureValuesOptions) SetSegmentRules(segmentRules []SegmentRule) *UpdateFeatureValuesOptions {
	options.SegmentRules = segmentRules
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateFeatureValuesOptions) SetHeaders(param map[string]string) *UpdateFeatureValuesOptions {
	options.Headers = param
	return options
}

// UpdatePropertyOptions : The UpdateProperty options.
type UpdatePropertyOptions struct {
	// Environment Id.
	EnvironmentID *string `json:"environment_id" validate:"required,ne="`

	// Property Id.
	PropertyID *string `json:"property_id" validate:"required,ne="`

	// Property name.
	Name *string `json:"name,omitempty"`

	// Property description.
	Description *string `json:"description,omitempty"`

	// Value of the property. The value can be Boolean, String or a Numeric value as per the `type` attribute.
	Value interface{} `json:"value,omitempty"`

	// Tags associated with the property.
	Tags *string `json:"tags,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// List of collection id representing the collections that are associated with the specified property.
	Collections []CollectionRef `json:"collections,omitempty"`

	// Allows users to set headers on API requests
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
func (options *UpdatePropertyOptions) SetEnvironmentID(environmentID string) *UpdatePropertyOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetPropertyID : Allow user to set PropertyID
func (options *UpdatePropertyOptions) SetPropertyID(propertyID string) *UpdatePropertyOptions {
	options.PropertyID = core.StringPtr(propertyID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdatePropertyOptions) SetName(name string) *UpdatePropertyOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdatePropertyOptions) SetDescription(description string) *UpdatePropertyOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetValue : Allow user to set Value
func (options *UpdatePropertyOptions) SetValue(value interface{}) *UpdatePropertyOptions {
	options.Value = value
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdatePropertyOptions) SetTags(tags string) *UpdatePropertyOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetSegmentRules : Allow user to set SegmentRules
func (options *UpdatePropertyOptions) SetSegmentRules(segmentRules []SegmentRule) *UpdatePropertyOptions {
	options.SegmentRules = segmentRules
	return options
}

// SetCollections : Allow user to set Collections
func (options *UpdatePropertyOptions) SetCollections(collections []CollectionRef) *UpdatePropertyOptions {
	options.Collections = collections
	return options
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

	// Property name.
	Name *string `json:"name,omitempty"`

	// Property description.
	Description *string `json:"description,omitempty"`

	// Tags associated with the property.
	Tags *string `json:"tags,omitempty"`

	// Value of the property. The value can be Boolean, String or a Numeric value as per the `type` attribute.
	Value interface{} `json:"value,omitempty"`

	// Specify the targeting rules that is used to set different property values for different segments.
	SegmentRules []SegmentRule `json:"segment_rules,omitempty"`

	// Allows users to set headers on API requests
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
func (options *UpdatePropertyValuesOptions) SetEnvironmentID(environmentID string) *UpdatePropertyValuesOptions {
	options.EnvironmentID = core.StringPtr(environmentID)
	return options
}

// SetPropertyID : Allow user to set PropertyID
func (options *UpdatePropertyValuesOptions) SetPropertyID(propertyID string) *UpdatePropertyValuesOptions {
	options.PropertyID = core.StringPtr(propertyID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdatePropertyValuesOptions) SetName(name string) *UpdatePropertyValuesOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdatePropertyValuesOptions) SetDescription(description string) *UpdatePropertyValuesOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdatePropertyValuesOptions) SetTags(tags string) *UpdatePropertyValuesOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetValue : Allow user to set Value
func (options *UpdatePropertyValuesOptions) SetValue(value interface{}) *UpdatePropertyValuesOptions {
	options.Value = value
	return options
}

// SetSegmentRules : Allow user to set SegmentRules
func (options *UpdatePropertyValuesOptions) SetSegmentRules(segmentRules []SegmentRule) *UpdatePropertyValuesOptions {
	options.SegmentRules = segmentRules
	return options
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

	// Segment name.
	Name *string `json:"name,omitempty"`

	// Segment description.
	Description *string `json:"description,omitempty"`

	// Tags associated with segments.
	Tags *string `json:"tags,omitempty"`

	// List of rules that determine if the entity is part of the segment.
	Rules []Rule `json:"rules,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSegmentOptions : Instantiate UpdateSegmentOptions
func (*AppConfigurationV1) NewUpdateSegmentOptions(segmentID string) *UpdateSegmentOptions {
	return &UpdateSegmentOptions{
		SegmentID: core.StringPtr(segmentID),
	}
}

// SetSegmentID : Allow user to set SegmentID
func (options *UpdateSegmentOptions) SetSegmentID(segmentID string) *UpdateSegmentOptions {
	options.SegmentID = core.StringPtr(segmentID)
	return options
}

// SetName : Allow user to set Name
func (options *UpdateSegmentOptions) SetName(name string) *UpdateSegmentOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetDescription : Allow user to set Description
func (options *UpdateSegmentOptions) SetDescription(description string) *UpdateSegmentOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTags : Allow user to set Tags
func (options *UpdateSegmentOptions) SetTags(tags string) *UpdateSegmentOptions {
	options.Tags = core.StringPtr(tags)
	return options
}

// SetRules : Allow user to set Rules
func (options *UpdateSegmentOptions) SetRules(rules []Rule) *UpdateSegmentOptions {
	options.Rules = rules
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSegmentOptions) SetHeaders(param map[string]string) *UpdateSegmentOptions {
	options.Headers = param
	return options
}
