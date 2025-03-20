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

// Package mqcloudv1 : Operations and models for the MqcloudV1 service
package mqcloudv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/IBM/go-sdk-core/v5/core"
	common "github.com/IBM/mqcloud-go-sdk/common"
	"github.com/go-openapi/strfmt"
)

// MqcloudV1 : The MQ on Cloud API defines a REST API interface to work with MQ on Cloud service in IBM Cloud.
//
// API Version: 1.1.0
type MqcloudV1 struct {
	Service *core.BaseService

	// The acceptable list of languages supported in the client.
	AcceptLanguage *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.private.eu-de.mq2.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "mqcloud"

const ParameterizedServiceURL = "https://api.private.{region}.mq2.cloud.ibm.com"

var defaultUrlVariables = map[string]string{
	"region": "eu-de",
}

// MqcloudV1Options : Service options
type MqcloudV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// The acceptable list of languages supported in the client.
	AcceptLanguage *string
}

// NewMqcloudV1UsingExternalConfig : constructs an instance of MqcloudV1 with passed in options and external configuration.
func NewMqcloudV1UsingExternalConfig(options *MqcloudV1Options) (mqcloud *MqcloudV1, err error) {
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

	mqcloud, err = NewMqcloudV1(options)
	err = core.RepurposeSDKProblem(err, "new-client-error")
	if err != nil {
		return
	}

	err = mqcloud.Service.ConfigureService(options.ServiceName)
	if err != nil {
		err = core.SDKErrorf(err, "", "client-config-error", common.GetComponentInfo())
		return
	}

	if options.URL != "" {
		err = mqcloud.Service.SetServiceURL(options.URL)
		err = core.RepurposeSDKProblem(err, "url-set-error")
	}
	return
}

// NewMqcloudV1 : constructs an instance of MqcloudV1 with passed in options.
func NewMqcloudV1(options *MqcloudV1Options) (service *MqcloudV1, err error) {
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

	service = &MqcloudV1{
		Service:        baseService,
		AcceptLanguage: options.AcceptLanguage,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", core.SDKErrorf(nil, "service does not support regional URLs", "no-regional-support", common.GetComponentInfo())
}

// Clone makes a copy of "mqcloud" suitable for processing requests.
func (mqcloud *MqcloudV1) Clone() *MqcloudV1 {
	if core.IsNil(mqcloud) {
		return nil
	}
	clone := *mqcloud
	clone.Service = mqcloud.Service.Clone()
	return &clone
}

// ConstructServiceURL constructs a service URL from the parameterized URL.
func ConstructServiceURL(providedUrlVariables map[string]string) (string, error) {
	return core.ConstructServiceURL(ParameterizedServiceURL, defaultUrlVariables, providedUrlVariables)
}

// SetServiceURL sets the service URL
func (mqcloud *MqcloudV1) SetServiceURL(url string) error {
	err := mqcloud.Service.SetServiceURL(url)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-set-error", common.GetComponentInfo())
	}
	return err
}

// GetServiceURL returns the service URL
func (mqcloud *MqcloudV1) GetServiceURL() string {
	return mqcloud.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (mqcloud *MqcloudV1) SetDefaultHeaders(headers http.Header) {
	mqcloud.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (mqcloud *MqcloudV1) SetEnableGzipCompression(enableGzip bool) {
	mqcloud.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (mqcloud *MqcloudV1) GetEnableGzipCompression() bool {
	return mqcloud.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (mqcloud *MqcloudV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	mqcloud.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (mqcloud *MqcloudV1) DisableRetries() {
	mqcloud.Service.DisableRetries()
}

// GetUsageDetails : Get the usage details
// Get the usage details.
func (mqcloud *MqcloudV1) GetUsageDetails(getUsageDetailsOptions *GetUsageDetailsOptions) (result *Usage, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetUsageDetailsWithContext(context.Background(), getUsageDetailsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetUsageDetailsWithContext is an alternate form of the GetUsageDetails method which supports a Context parameter
func (mqcloud *MqcloudV1) GetUsageDetailsWithContext(ctx context.Context, getUsageDetailsOptions *GetUsageDetailsOptions) (result *Usage, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getUsageDetailsOptions, "getUsageDetailsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getUsageDetailsOptions, "getUsageDetailsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getUsageDetailsOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/usage`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getUsageDetailsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetUsageDetails")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_usage_details", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUsage)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetOptions : Return configuration options (eg, available deployment locations, queue manager sizes)
// Return configuration options (eg, available deployment locations, queue manager sizes).
func (mqcloud *MqcloudV1) GetOptions(getOptionsOptions *GetOptionsOptions) (result *ConfigurationOptions, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetOptionsWithContext(context.Background(), getOptionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetOptionsWithContext is an alternate form of the GetOptions method which supports a Context parameter
func (mqcloud *MqcloudV1) GetOptionsWithContext(ctx context.Context, getOptionsOptions *GetOptionsOptions) (result *ConfigurationOptions, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getOptionsOptions, "getOptionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getOptionsOptions, "getOptionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getOptionsOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/options`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getOptionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetOptions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_options", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConfigurationOptions)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateQueueManager : Create a new queue manager
// Create a new queue manager.
func (mqcloud *MqcloudV1) CreateQueueManager(createQueueManagerOptions *CreateQueueManagerOptions) (result *QueueManagerTaskStatus, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateQueueManagerWithContext(context.Background(), createQueueManagerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateQueueManagerWithContext is an alternate form of the CreateQueueManager method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateQueueManagerWithContext(ctx context.Context, createQueueManagerOptions *CreateQueueManagerOptions) (result *QueueManagerTaskStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createQueueManagerOptions, "createQueueManagerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createQueueManagerOptions, "createQueueManagerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createQueueManagerOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createQueueManagerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateQueueManager")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createQueueManagerOptions.Name != nil {
		body["name"] = createQueueManagerOptions.Name
	}
	if createQueueManagerOptions.Location != nil {
		body["location"] = createQueueManagerOptions.Location
	}
	if createQueueManagerOptions.Size != nil {
		body["size"] = createQueueManagerOptions.Size
	}
	if createQueueManagerOptions.DisplayName != nil {
		body["display_name"] = createQueueManagerOptions.DisplayName
	}
	if createQueueManagerOptions.Version != nil {
		body["version"] = createQueueManagerOptions.Version
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_queue_manager", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerTaskStatus)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListQueueManagers : Get list of queue managers
// Get a list of the queue manager summaries which exist in this service instance.
func (mqcloud *MqcloudV1) ListQueueManagers(listQueueManagersOptions *ListQueueManagersOptions) (result *QueueManagerDetailsCollection, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.ListQueueManagersWithContext(context.Background(), listQueueManagersOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListQueueManagersWithContext is an alternate form of the ListQueueManagers method which supports a Context parameter
func (mqcloud *MqcloudV1) ListQueueManagersWithContext(ctx context.Context, listQueueManagersOptions *ListQueueManagersOptions) (result *QueueManagerDetailsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listQueueManagersOptions, "listQueueManagersOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listQueueManagersOptions, "listQueueManagersOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *listQueueManagersOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listQueueManagersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "ListQueueManagers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	if listQueueManagersOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listQueueManagersOptions.Offset))
	}
	if listQueueManagersOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listQueueManagersOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_queue_managers", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerDetailsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetQueueManager : Get details of a queue manager
// Get the details of a given queue manager.
func (mqcloud *MqcloudV1) GetQueueManager(getQueueManagerOptions *GetQueueManagerOptions) (result *QueueManagerDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetQueueManagerWithContext(context.Background(), getQueueManagerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetQueueManagerWithContext is an alternate form of the GetQueueManager method which supports a Context parameter
func (mqcloud *MqcloudV1) GetQueueManagerWithContext(ctx context.Context, getQueueManagerOptions *GetQueueManagerOptions) (result *QueueManagerDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getQueueManagerOptions, "getQueueManagerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getQueueManagerOptions, "getQueueManagerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getQueueManagerOptions.ServiceInstanceGuid,
		"queue_manager_id":      *getQueueManagerOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getQueueManagerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetQueueManager")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_queue_manager", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteQueueManager : Delete a queue manager
// Delete a queue manager.
func (mqcloud *MqcloudV1) DeleteQueueManager(deleteQueueManagerOptions *DeleteQueueManagerOptions) (result *QueueManagerTaskStatus, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.DeleteQueueManagerWithContext(context.Background(), deleteQueueManagerOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteQueueManagerWithContext is an alternate form of the DeleteQueueManager method which supports a Context parameter
func (mqcloud *MqcloudV1) DeleteQueueManagerWithContext(ctx context.Context, deleteQueueManagerOptions *DeleteQueueManagerOptions) (result *QueueManagerTaskStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteQueueManagerOptions, "deleteQueueManagerOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteQueueManagerOptions, "deleteQueueManagerOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *deleteQueueManagerOptions.ServiceInstanceGuid,
		"queue_manager_id":      *deleteQueueManagerOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteQueueManagerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DeleteQueueManager")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_queue_manager", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerTaskStatus)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// SetQueueManagerVersion : Upgrade a queue manager
// Upgrade a queue manager.
func (mqcloud *MqcloudV1) SetQueueManagerVersion(setQueueManagerVersionOptions *SetQueueManagerVersionOptions) (result *QueueManagerTaskStatus, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.SetQueueManagerVersionWithContext(context.Background(), setQueueManagerVersionOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetQueueManagerVersionWithContext is an alternate form of the SetQueueManagerVersion method which supports a Context parameter
func (mqcloud *MqcloudV1) SetQueueManagerVersionWithContext(ctx context.Context, setQueueManagerVersionOptions *SetQueueManagerVersionOptions) (result *QueueManagerTaskStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setQueueManagerVersionOptions, "setQueueManagerVersionOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(setQueueManagerVersionOptions, "setQueueManagerVersionOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *setQueueManagerVersionOptions.ServiceInstanceGuid,
		"queue_manager_id":      *setQueueManagerVersionOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/version`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range setQueueManagerVersionOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "SetQueueManagerVersion")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if setQueueManagerVersionOptions.Version != nil {
		body["version"] = setQueueManagerVersionOptions.Version
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_queue_manager_version", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerTaskStatus)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetQueueManagerAvailableUpgradeVersions : Get the list of available versions that this queue manager can be upgraded to
// Get the list of available versions that this queue manager can be upgraded to.
func (mqcloud *MqcloudV1) GetQueueManagerAvailableUpgradeVersions(getQueueManagerAvailableUpgradeVersionsOptions *GetQueueManagerAvailableUpgradeVersionsOptions) (result *QueueManagerVersionUpgrades, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetQueueManagerAvailableUpgradeVersionsWithContext(context.Background(), getQueueManagerAvailableUpgradeVersionsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetQueueManagerAvailableUpgradeVersionsWithContext is an alternate form of the GetQueueManagerAvailableUpgradeVersions method which supports a Context parameter
func (mqcloud *MqcloudV1) GetQueueManagerAvailableUpgradeVersionsWithContext(ctx context.Context, getQueueManagerAvailableUpgradeVersionsOptions *GetQueueManagerAvailableUpgradeVersionsOptions) (result *QueueManagerVersionUpgrades, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getQueueManagerAvailableUpgradeVersionsOptions, "getQueueManagerAvailableUpgradeVersionsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getQueueManagerAvailableUpgradeVersionsOptions, "getQueueManagerAvailableUpgradeVersionsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getQueueManagerAvailableUpgradeVersionsOptions.ServiceInstanceGuid,
		"queue_manager_id":      *getQueueManagerAvailableUpgradeVersionsOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/available_versions`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getQueueManagerAvailableUpgradeVersionsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetQueueManagerAvailableUpgradeVersions")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_queue_manager_available_upgrade_versions", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerVersionUpgrades)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetQueueManagerConnectionInfo : Get connection information for a queue manager
// Get connection information for a queue manager.
func (mqcloud *MqcloudV1) GetQueueManagerConnectionInfo(getQueueManagerConnectionInfoOptions *GetQueueManagerConnectionInfoOptions) (result *ConnectionInfo, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetQueueManagerConnectionInfoWithContext(context.Background(), getQueueManagerConnectionInfoOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetQueueManagerConnectionInfoWithContext is an alternate form of the GetQueueManagerConnectionInfo method which supports a Context parameter
func (mqcloud *MqcloudV1) GetQueueManagerConnectionInfoWithContext(ctx context.Context, getQueueManagerConnectionInfoOptions *GetQueueManagerConnectionInfoOptions) (result *ConnectionInfo, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getQueueManagerConnectionInfoOptions, "getQueueManagerConnectionInfoOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getQueueManagerConnectionInfoOptions, "getQueueManagerConnectionInfoOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getQueueManagerConnectionInfoOptions.ServiceInstanceGuid,
		"queue_manager_id":      *getQueueManagerConnectionInfoOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/connection_info`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getQueueManagerConnectionInfoOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetQueueManagerConnectionInfo")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_queue_manager_connection_info", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalConnectionInfo)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetQueueManagerStatus : Get the status of the queue manager
// Get the status of the queue manager instance.
func (mqcloud *MqcloudV1) GetQueueManagerStatus(getQueueManagerStatusOptions *GetQueueManagerStatusOptions) (result *QueueManagerStatus, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetQueueManagerStatusWithContext(context.Background(), getQueueManagerStatusOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetQueueManagerStatusWithContext is an alternate form of the GetQueueManagerStatus method which supports a Context parameter
func (mqcloud *MqcloudV1) GetQueueManagerStatusWithContext(ctx context.Context, getQueueManagerStatusOptions *GetQueueManagerStatusOptions) (result *QueueManagerStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getQueueManagerStatusOptions, "getQueueManagerStatusOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getQueueManagerStatusOptions, "getQueueManagerStatusOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getQueueManagerStatusOptions.ServiceInstanceGuid,
		"queue_manager_id":      *getQueueManagerStatusOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/status`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getQueueManagerStatusOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetQueueManagerStatus")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_queue_manager_status", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalQueueManagerStatus)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListUsers : Get a list of users for an instance
// Get a list of users for an instance.
func (mqcloud *MqcloudV1) ListUsers(listUsersOptions *ListUsersOptions) (result *UserDetailsCollection, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.ListUsersWithContext(context.Background(), listUsersOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListUsersWithContext is an alternate form of the ListUsers method which supports a Context parameter
func (mqcloud *MqcloudV1) ListUsersWithContext(ctx context.Context, listUsersOptions *ListUsersOptions) (result *UserDetailsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listUsersOptions, "listUsersOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listUsersOptions, "listUsersOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *listUsersOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/users`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listUsersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "ListUsers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	if listUsersOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listUsersOptions.Offset))
	}
	if listUsersOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listUsersOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_users", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserDetailsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateUser : Add a user to an instance
// Add a user to an instance.
func (mqcloud *MqcloudV1) CreateUser(createUserOptions *CreateUserOptions) (result *UserDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateUserWithContext(context.Background(), createUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateUserWithContext is an alternate form of the CreateUser method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateUserWithContext(ctx context.Context, createUserOptions *CreateUserOptions) (result *UserDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createUserOptions, "createUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createUserOptions, "createUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createUserOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/users`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createUserOptions.Email != nil {
		body["email"] = createUserOptions.Email
	}
	if createUserOptions.Name != nil {
		body["name"] = createUserOptions.Name
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetUser : Get a user for an instance
// Get a user for an instance.
func (mqcloud *MqcloudV1) GetUser(getUserOptions *GetUserOptions) (result *UserDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetUserWithContext(context.Background(), getUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetUserWithContext is an alternate form of the GetUser method which supports a Context parameter
func (mqcloud *MqcloudV1) GetUserWithContext(ctx context.Context, getUserOptions *GetUserOptions) (result *UserDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getUserOptions, "getUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getUserOptions, "getUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getUserOptions.ServiceInstanceGuid,
		"user_id":               *getUserOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/users/{user_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalUserDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteUser : Delete a user for an instance
// Delete a user for an instance.
func (mqcloud *MqcloudV1) DeleteUser(deleteUserOptions *DeleteUserOptions) (response *core.DetailedResponse, err error) {
	response, err = mqcloud.DeleteUserWithContext(context.Background(), deleteUserOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteUserWithContext is an alternate form of the DeleteUser method which supports a Context parameter
func (mqcloud *MqcloudV1) DeleteUserWithContext(ctx context.Context, deleteUserOptions *DeleteUserOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteUserOptions, "deleteUserOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteUserOptions, "deleteUserOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *deleteUserOptions.ServiceInstanceGuid,
		"user_id":               *deleteUserOptions.UserID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/users/{user_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteUserOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DeleteUser")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_user", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// ListApplications : Get a list of applications for an instance
// Get a list of applications for an instance.
func (mqcloud *MqcloudV1) ListApplications(listApplicationsOptions *ListApplicationsOptions) (result *ApplicationDetailsCollection, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.ListApplicationsWithContext(context.Background(), listApplicationsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListApplicationsWithContext is an alternate form of the ListApplications method which supports a Context parameter
func (mqcloud *MqcloudV1) ListApplicationsWithContext(ctx context.Context, listApplicationsOptions *ListApplicationsOptions) (result *ApplicationDetailsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listApplicationsOptions, "listApplicationsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listApplicationsOptions, "listApplicationsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *listApplicationsOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/applications`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listApplicationsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "ListApplications")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	if listApplicationsOptions.Offset != nil {
		builder.AddQuery("offset", fmt.Sprint(*listApplicationsOptions.Offset))
	}
	if listApplicationsOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listApplicationsOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_applications", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplicationDetailsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateApplication : Add an application to an instance
// Add an application to an instance.
func (mqcloud *MqcloudV1) CreateApplication(createApplicationOptions *CreateApplicationOptions) (result *ApplicationCreated, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateApplicationWithContext(context.Background(), createApplicationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateApplicationWithContext is an alternate form of the CreateApplication method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateApplicationWithContext(ctx context.Context, createApplicationOptions *CreateApplicationOptions) (result *ApplicationCreated, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createApplicationOptions, "createApplicationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createApplicationOptions, "createApplicationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createApplicationOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/applications`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createApplicationOptions.Name != nil {
		body["name"] = createApplicationOptions.Name
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_application", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplicationCreated)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetApplication : Get an application for an instance
// Get an application for an instance.
func (mqcloud *MqcloudV1) GetApplication(getApplicationOptions *GetApplicationOptions) (result *ApplicationDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetApplicationWithContext(context.Background(), getApplicationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetApplicationWithContext is an alternate form of the GetApplication method which supports a Context parameter
func (mqcloud *MqcloudV1) GetApplicationWithContext(ctx context.Context, getApplicationOptions *GetApplicationOptions) (result *ApplicationDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getApplicationOptions, "getApplicationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getApplicationOptions, "getApplicationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getApplicationOptions.ServiceInstanceGuid,
		"application_id":        *getApplicationOptions.ApplicationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/applications/{application_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_application", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplicationDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteApplication : Delete an application from an instance
// Delete an application from an instance.
func (mqcloud *MqcloudV1) DeleteApplication(deleteApplicationOptions *DeleteApplicationOptions) (response *core.DetailedResponse, err error) {
	response, err = mqcloud.DeleteApplicationWithContext(context.Background(), deleteApplicationOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteApplicationWithContext is an alternate form of the DeleteApplication method which supports a Context parameter
func (mqcloud *MqcloudV1) DeleteApplicationWithContext(ctx context.Context, deleteApplicationOptions *DeleteApplicationOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteApplicationOptions, "deleteApplicationOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteApplicationOptions, "deleteApplicationOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *deleteApplicationOptions.ServiceInstanceGuid,
		"application_id":        *deleteApplicationOptions.ApplicationID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/applications/{application_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteApplicationOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DeleteApplication")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_application", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateApplicationApikey : Create a new apikey for an application
// Create a new apikey for an application.
func (mqcloud *MqcloudV1) CreateApplicationApikey(createApplicationApikeyOptions *CreateApplicationApikeyOptions) (result *ApplicationAPIKeyCreated, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateApplicationApikeyWithContext(context.Background(), createApplicationApikeyOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateApplicationApikeyWithContext is an alternate form of the CreateApplicationApikey method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateApplicationApikeyWithContext(ctx context.Context, createApplicationApikeyOptions *CreateApplicationApikeyOptions) (result *ApplicationAPIKeyCreated, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createApplicationApikeyOptions, "createApplicationApikeyOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createApplicationApikeyOptions, "createApplicationApikeyOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createApplicationApikeyOptions.ServiceInstanceGuid,
		"application_id":        *createApplicationApikeyOptions.ApplicationID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/applications/{application_id}/api_key`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createApplicationApikeyOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateApplicationApikey")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if createApplicationApikeyOptions.Name != nil {
		body["name"] = createApplicationApikeyOptions.Name
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_application_apikey", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalApplicationAPIKeyCreated)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateTrustStorePemCertificate : Upload a trust store certificate
// Import TLS certificate from a single self-contained PEM file to the truststore.
func (mqcloud *MqcloudV1) CreateTrustStorePemCertificate(createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions) (result *TrustStoreCertificateDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateTrustStorePemCertificateWithContext(context.Background(), createTrustStorePemCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateTrustStorePemCertificateWithContext is an alternate form of the CreateTrustStorePemCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateTrustStorePemCertificateWithContext(ctx context.Context, createTrustStorePemCertificateOptions *CreateTrustStorePemCertificateOptions) (result *TrustStoreCertificateDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createTrustStorePemCertificateOptions, "createTrustStorePemCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createTrustStorePemCertificateOptions, "createTrustStorePemCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createTrustStorePemCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *createTrustStorePemCertificateOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/trust_store`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createTrustStorePemCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateTrustStorePemCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	builder.AddFormData("label", "", "", fmt.Sprint(*createTrustStorePemCertificateOptions.Label))
	builder.AddFormData("certificate_file", "filename",
		"application/octet-stream", createTrustStorePemCertificateOptions.CertificateFile)

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_trust_store_pem_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustStoreCertificateDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListTrustStoreCertificates : List trust store certificates
// Get the list of certificates in the queue manager's certificate trust store.
func (mqcloud *MqcloudV1) ListTrustStoreCertificates(listTrustStoreCertificatesOptions *ListTrustStoreCertificatesOptions) (result *TrustStoreCertificateDetailsCollection, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.ListTrustStoreCertificatesWithContext(context.Background(), listTrustStoreCertificatesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListTrustStoreCertificatesWithContext is an alternate form of the ListTrustStoreCertificates method which supports a Context parameter
func (mqcloud *MqcloudV1) ListTrustStoreCertificatesWithContext(ctx context.Context, listTrustStoreCertificatesOptions *ListTrustStoreCertificatesOptions) (result *TrustStoreCertificateDetailsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listTrustStoreCertificatesOptions, "listTrustStoreCertificatesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listTrustStoreCertificatesOptions, "listTrustStoreCertificatesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *listTrustStoreCertificatesOptions.ServiceInstanceGuid,
		"queue_manager_id":      *listTrustStoreCertificatesOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/trust_store`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listTrustStoreCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "ListTrustStoreCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_trust_store_certificates", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustStoreCertificateDetailsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetTrustStoreCertificate : Get a trust store certificate
// Get a trust store certificate from a queue manager.
func (mqcloud *MqcloudV1) GetTrustStoreCertificate(getTrustStoreCertificateOptions *GetTrustStoreCertificateOptions) (result *TrustStoreCertificateDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetTrustStoreCertificateWithContext(context.Background(), getTrustStoreCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetTrustStoreCertificateWithContext is an alternate form of the GetTrustStoreCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) GetTrustStoreCertificateWithContext(ctx context.Context, getTrustStoreCertificateOptions *GetTrustStoreCertificateOptions) (result *TrustStoreCertificateDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getTrustStoreCertificateOptions, "getTrustStoreCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getTrustStoreCertificateOptions, "getTrustStoreCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getTrustStoreCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *getTrustStoreCertificateOptions.QueueManagerID,
		"certificate_id":        *getTrustStoreCertificateOptions.CertificateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/trust_store/{certificate_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getTrustStoreCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetTrustStoreCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_trust_store_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalTrustStoreCertificateDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteTrustStoreCertificate : Delete a trust store certificate
// Delete a trust store certificate.
func (mqcloud *MqcloudV1) DeleteTrustStoreCertificate(deleteTrustStoreCertificateOptions *DeleteTrustStoreCertificateOptions) (response *core.DetailedResponse, err error) {
	response, err = mqcloud.DeleteTrustStoreCertificateWithContext(context.Background(), deleteTrustStoreCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteTrustStoreCertificateWithContext is an alternate form of the DeleteTrustStoreCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) DeleteTrustStoreCertificateWithContext(ctx context.Context, deleteTrustStoreCertificateOptions *DeleteTrustStoreCertificateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteTrustStoreCertificateOptions, "deleteTrustStoreCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteTrustStoreCertificateOptions, "deleteTrustStoreCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *deleteTrustStoreCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *deleteTrustStoreCertificateOptions.QueueManagerID,
		"certificate_id":        *deleteTrustStoreCertificateOptions.CertificateID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/trust_store/{certificate_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteTrustStoreCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DeleteTrustStoreCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_trust_store_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DownloadTrustStoreCertificate : Download a queue manager's certificate from its trust store
// Download the specified trust store certificate PEM file from the queue manager.
func (mqcloud *MqcloudV1) DownloadTrustStoreCertificate(downloadTrustStoreCertificateOptions *DownloadTrustStoreCertificateOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.DownloadTrustStoreCertificateWithContext(context.Background(), downloadTrustStoreCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DownloadTrustStoreCertificateWithContext is an alternate form of the DownloadTrustStoreCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) DownloadTrustStoreCertificateWithContext(ctx context.Context, downloadTrustStoreCertificateOptions *DownloadTrustStoreCertificateOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(downloadTrustStoreCertificateOptions, "downloadTrustStoreCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(downloadTrustStoreCertificateOptions, "downloadTrustStoreCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *downloadTrustStoreCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *downloadTrustStoreCertificateOptions.QueueManagerID,
		"certificate_id":        *downloadTrustStoreCertificateOptions.CertificateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/trust_store/{certificate_id}/download`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range downloadTrustStoreCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DownloadTrustStoreCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/octet-stream")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, &result)
	if err != nil {
		core.EnrichHTTPProblem(err, "download_trust_store_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// CreateKeyStorePemCertificate : Upload a key store certificate
// Import TLS certificate from a single self-contained PEM file into the queue manager's key store.
func (mqcloud *MqcloudV1) CreateKeyStorePemCertificate(createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions) (result *KeyStoreCertificateDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateKeyStorePemCertificateWithContext(context.Background(), createKeyStorePemCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateKeyStorePemCertificateWithContext is an alternate form of the CreateKeyStorePemCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateKeyStorePemCertificateWithContext(ctx context.Context, createKeyStorePemCertificateOptions *CreateKeyStorePemCertificateOptions) (result *KeyStoreCertificateDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createKeyStorePemCertificateOptions, "createKeyStorePemCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createKeyStorePemCertificateOptions, "createKeyStorePemCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createKeyStorePemCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *createKeyStorePemCertificateOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createKeyStorePemCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateKeyStorePemCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	builder.AddFormData("label", "", "", fmt.Sprint(*createKeyStorePemCertificateOptions.Label))
	builder.AddFormData("certificate_file", "filename",
		"application/octet-stream", createKeyStorePemCertificateOptions.CertificateFile)

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_key_store_pem_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeyStoreCertificateDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListKeyStoreCertificates : List key store certificates
// Get a list of certificates in the queue manager's certificate key store.
func (mqcloud *MqcloudV1) ListKeyStoreCertificates(listKeyStoreCertificatesOptions *ListKeyStoreCertificatesOptions) (result *KeyStoreCertificateDetailsCollection, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.ListKeyStoreCertificatesWithContext(context.Background(), listKeyStoreCertificatesOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListKeyStoreCertificatesWithContext is an alternate form of the ListKeyStoreCertificates method which supports a Context parameter
func (mqcloud *MqcloudV1) ListKeyStoreCertificatesWithContext(ctx context.Context, listKeyStoreCertificatesOptions *ListKeyStoreCertificatesOptions) (result *KeyStoreCertificateDetailsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listKeyStoreCertificatesOptions, "listKeyStoreCertificatesOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listKeyStoreCertificatesOptions, "listKeyStoreCertificatesOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *listKeyStoreCertificatesOptions.ServiceInstanceGuid,
		"queue_manager_id":      *listKeyStoreCertificatesOptions.QueueManagerID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listKeyStoreCertificatesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "ListKeyStoreCertificates")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_key_store_certificates", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeyStoreCertificateDetailsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetKeyStoreCertificate : Get a key store certificate for queue manager
// Get a key store certificate for queue manager.
func (mqcloud *MqcloudV1) GetKeyStoreCertificate(getKeyStoreCertificateOptions *GetKeyStoreCertificateOptions) (result *KeyStoreCertificateDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetKeyStoreCertificateWithContext(context.Background(), getKeyStoreCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetKeyStoreCertificateWithContext is an alternate form of the GetKeyStoreCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) GetKeyStoreCertificateWithContext(ctx context.Context, getKeyStoreCertificateOptions *GetKeyStoreCertificateOptions) (result *KeyStoreCertificateDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getKeyStoreCertificateOptions, "getKeyStoreCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getKeyStoreCertificateOptions, "getKeyStoreCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *getKeyStoreCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *getKeyStoreCertificateOptions.QueueManagerID,
		"certificate_id":        *getKeyStoreCertificateOptions.CertificateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store/{certificate_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getKeyStoreCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetKeyStoreCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_key_store_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalKeyStoreCertificateDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteKeyStoreCertificate : Delete a queue manager's key store certificate
// Delete a queue manager's key store certificate.
func (mqcloud *MqcloudV1) DeleteKeyStoreCertificate(deleteKeyStoreCertificateOptions *DeleteKeyStoreCertificateOptions) (response *core.DetailedResponse, err error) {
	response, err = mqcloud.DeleteKeyStoreCertificateWithContext(context.Background(), deleteKeyStoreCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteKeyStoreCertificateWithContext is an alternate form of the DeleteKeyStoreCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) DeleteKeyStoreCertificateWithContext(ctx context.Context, deleteKeyStoreCertificateOptions *DeleteKeyStoreCertificateOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteKeyStoreCertificateOptions, "deleteKeyStoreCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteKeyStoreCertificateOptions, "deleteKeyStoreCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *deleteKeyStoreCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *deleteKeyStoreCertificateOptions.QueueManagerID,
		"certificate_id":        *deleteKeyStoreCertificateOptions.CertificateID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store/{certificate_id}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteKeyStoreCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DeleteKeyStoreCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_key_store_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// DownloadKeyStoreCertificate : Download a queue manager's certificate from its key store
// Download the specified key store certificate PEM file from the queue manager.
func (mqcloud *MqcloudV1) DownloadKeyStoreCertificate(downloadKeyStoreCertificateOptions *DownloadKeyStoreCertificateOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.DownloadKeyStoreCertificateWithContext(context.Background(), downloadKeyStoreCertificateOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DownloadKeyStoreCertificateWithContext is an alternate form of the DownloadKeyStoreCertificate method which supports a Context parameter
func (mqcloud *MqcloudV1) DownloadKeyStoreCertificateWithContext(ctx context.Context, downloadKeyStoreCertificateOptions *DownloadKeyStoreCertificateOptions) (result io.ReadCloser, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(downloadKeyStoreCertificateOptions, "downloadKeyStoreCertificateOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(downloadKeyStoreCertificateOptions, "downloadKeyStoreCertificateOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *downloadKeyStoreCertificateOptions.ServiceInstanceGuid,
		"queue_manager_id":      *downloadKeyStoreCertificateOptions.QueueManagerID,
		"certificate_id":        *downloadKeyStoreCertificateOptions.CertificateID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store/{certificate_id}/download`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range downloadKeyStoreCertificateOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DownloadKeyStoreCertificate")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/octet-stream")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, &result)
	if err != nil {
		core.EnrichHTTPProblem(err, "download_key_store_certificate", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}

// GetCertificateAmsChannels : Get the AMS channels that are configured with this key store certificate
// Get the AMS channels that are configured with this key store certificate.
func (mqcloud *MqcloudV1) GetCertificateAmsChannels(getCertificateAmsChannelsOptions *GetCertificateAmsChannelsOptions) (result *ChannelsDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetCertificateAmsChannelsWithContext(context.Background(), getCertificateAmsChannelsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetCertificateAmsChannelsWithContext is an alternate form of the GetCertificateAmsChannels method which supports a Context parameter
func (mqcloud *MqcloudV1) GetCertificateAmsChannelsWithContext(ctx context.Context, getCertificateAmsChannelsOptions *GetCertificateAmsChannelsOptions) (result *ChannelsDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getCertificateAmsChannelsOptions, "getCertificateAmsChannelsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getCertificateAmsChannelsOptions, "getCertificateAmsChannelsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"queue_manager_id":      *getCertificateAmsChannelsOptions.QueueManagerID,
		"certificate_id":        *getCertificateAmsChannelsOptions.CertificateID,
		"service_instance_guid": *getCertificateAmsChannelsOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store/{certificate_id}/config/ams`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getCertificateAmsChannelsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetCertificateAmsChannels")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_certificate_ams_channels", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalChannelsDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// SetCertificateAmsChannels : Update the AMS channels that are configured with this key store certificate
// Update the AMS channels that are configured with this key store certificate.
func (mqcloud *MqcloudV1) SetCertificateAmsChannels(setCertificateAmsChannelsOptions *SetCertificateAmsChannelsOptions) (result *ChannelsDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.SetCertificateAmsChannelsWithContext(context.Background(), setCertificateAmsChannelsOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// SetCertificateAmsChannelsWithContext is an alternate form of the SetCertificateAmsChannels method which supports a Context parameter
func (mqcloud *MqcloudV1) SetCertificateAmsChannelsWithContext(ctx context.Context, setCertificateAmsChannelsOptions *SetCertificateAmsChannelsOptions) (result *ChannelsDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(setCertificateAmsChannelsOptions, "setCertificateAmsChannelsOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(setCertificateAmsChannelsOptions, "setCertificateAmsChannelsOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"queue_manager_id":      *setCertificateAmsChannelsOptions.QueueManagerID,
		"certificate_id":        *setCertificateAmsChannelsOptions.CertificateID,
		"service_instance_guid": *setCertificateAmsChannelsOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/queue_managers/{queue_manager_id}/certificates/key_store/{certificate_id}/config/ams`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range setCertificateAmsChannelsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "SetCertificateAmsChannels")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}

	body := make(map[string]interface{})
	if setCertificateAmsChannelsOptions.Channels != nil {
		body["channels"] = setCertificateAmsChannelsOptions.Channels
	}
	if setCertificateAmsChannelsOptions.UpdateStrategy != nil {
		body["update_strategy"] = setCertificateAmsChannelsOptions.UpdateStrategy
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "set_certificate_ams_channels", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalChannelsDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// CreateVirtualPrivateEndpointGateway : Create a new virtual private endpoint gateway
// Create a new virtual private endpoint gateway.
func (mqcloud *MqcloudV1) CreateVirtualPrivateEndpointGateway(createVirtualPrivateEndpointGatewayOptions *CreateVirtualPrivateEndpointGatewayOptions) (result *VirtualPrivateEndpointGatewayDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.CreateVirtualPrivateEndpointGatewayWithContext(context.Background(), createVirtualPrivateEndpointGatewayOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// CreateVirtualPrivateEndpointGatewayWithContext is an alternate form of the CreateVirtualPrivateEndpointGateway method which supports a Context parameter
func (mqcloud *MqcloudV1) CreateVirtualPrivateEndpointGatewayWithContext(ctx context.Context, createVirtualPrivateEndpointGatewayOptions *CreateVirtualPrivateEndpointGatewayOptions) (result *VirtualPrivateEndpointGatewayDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createVirtualPrivateEndpointGatewayOptions, "createVirtualPrivateEndpointGatewayOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(createVirtualPrivateEndpointGatewayOptions, "createVirtualPrivateEndpointGatewayOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *createVirtualPrivateEndpointGatewayOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/virtual_private_endpoint_gateway`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range createVirtualPrivateEndpointGatewayOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "CreateVirtualPrivateEndpointGateway")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}
	if createVirtualPrivateEndpointGatewayOptions.TrustedProfile != nil {
		builder.AddHeader("Trusted-Profile", fmt.Sprint(*createVirtualPrivateEndpointGatewayOptions.TrustedProfile))
	}

	body := make(map[string]interface{})
	if createVirtualPrivateEndpointGatewayOptions.Name != nil {
		body["name"] = createVirtualPrivateEndpointGatewayOptions.Name
	}
	if createVirtualPrivateEndpointGatewayOptions.TargetCrn != nil {
		body["target_crn"] = createVirtualPrivateEndpointGatewayOptions.TargetCrn
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
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "create_virtual_private_endpoint_gateway", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVirtualPrivateEndpointGatewayDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// ListVirtualPrivateEndpointGateways : Get a list of information for all virtual private endpoint gateways
// Get a list of information for all Virtual private endpoint gateways.
func (mqcloud *MqcloudV1) ListVirtualPrivateEndpointGateways(listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions) (result *VirtualPrivateEndpointGatewayDetailsCollection, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.ListVirtualPrivateEndpointGatewaysWithContext(context.Background(), listVirtualPrivateEndpointGatewaysOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ListVirtualPrivateEndpointGatewaysWithContext is an alternate form of the ListVirtualPrivateEndpointGateways method which supports a Context parameter
func (mqcloud *MqcloudV1) ListVirtualPrivateEndpointGatewaysWithContext(ctx context.Context, listVirtualPrivateEndpointGatewaysOptions *ListVirtualPrivateEndpointGatewaysOptions) (result *VirtualPrivateEndpointGatewayDetailsCollection, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listVirtualPrivateEndpointGatewaysOptions, "listVirtualPrivateEndpointGatewaysOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(listVirtualPrivateEndpointGatewaysOptions, "listVirtualPrivateEndpointGatewaysOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid": *listVirtualPrivateEndpointGatewaysOptions.ServiceInstanceGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/virtual_private_endpoint_gateway`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range listVirtualPrivateEndpointGatewaysOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "ListVirtualPrivateEndpointGateways")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}
	if listVirtualPrivateEndpointGatewaysOptions.TrustedProfile != nil {
		builder.AddHeader("Trusted-Profile", fmt.Sprint(*listVirtualPrivateEndpointGatewaysOptions.TrustedProfile))
	}

	if listVirtualPrivateEndpointGatewaysOptions.Start != nil {
		builder.AddQuery("start", fmt.Sprint(*listVirtualPrivateEndpointGatewaysOptions.Start))
	}
	if listVirtualPrivateEndpointGatewaysOptions.Limit != nil {
		builder.AddQuery("limit", fmt.Sprint(*listVirtualPrivateEndpointGatewaysOptions.Limit))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "list_virtual_private_endpoint_gateways", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVirtualPrivateEndpointGatewayDetailsCollection)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// GetVirtualPrivateEndpointGateway : Display the information for a specific virtual private endpoint gateway
// Display the information for a specific virtual private endpoint gateway.
func (mqcloud *MqcloudV1) GetVirtualPrivateEndpointGateway(getVirtualPrivateEndpointGatewayOptions *GetVirtualPrivateEndpointGatewayOptions) (result *VirtualPrivateEndpointGatewayDetails, response *core.DetailedResponse, err error) {
	result, response, err = mqcloud.GetVirtualPrivateEndpointGatewayWithContext(context.Background(), getVirtualPrivateEndpointGatewayOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetVirtualPrivateEndpointGatewayWithContext is an alternate form of the GetVirtualPrivateEndpointGateway method which supports a Context parameter
func (mqcloud *MqcloudV1) GetVirtualPrivateEndpointGatewayWithContext(ctx context.Context, getVirtualPrivateEndpointGatewayOptions *GetVirtualPrivateEndpointGatewayOptions) (result *VirtualPrivateEndpointGatewayDetails, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getVirtualPrivateEndpointGatewayOptions, "getVirtualPrivateEndpointGatewayOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(getVirtualPrivateEndpointGatewayOptions, "getVirtualPrivateEndpointGatewayOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid":                 *getVirtualPrivateEndpointGatewayOptions.ServiceInstanceGuid,
		"virtual_private_endpoint_gateway_guid": *getVirtualPrivateEndpointGatewayOptions.VirtualPrivateEndpointGatewayGuid,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/virtual_private_endpoint_gateway/{virtual_private_endpoint_gateway_guid}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range getVirtualPrivateEndpointGatewayOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "GetVirtualPrivateEndpointGateway")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}
	if getVirtualPrivateEndpointGatewayOptions.TrustedProfile != nil {
		builder.AddHeader("Trusted-Profile", fmt.Sprint(*getVirtualPrivateEndpointGatewayOptions.TrustedProfile))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = mqcloud.Service.Request(request, &rawResponse)
	if err != nil {
		core.EnrichHTTPProblem(err, "get_virtual_private_endpoint_gateway", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}
	if rawResponse != nil {
		err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalVirtualPrivateEndpointGatewayDetails)
		if err != nil {
			err = core.SDKErrorf(err, "", "unmarshal-resp-error", common.GetComponentInfo())
			return
		}
		response.Result = result
	}

	return
}

// DeleteVirtualPrivateEndpointGateway : Delete a specific virtual private endpoint gateway
// Delete a specific virtual_private_endpoint_gateway.
func (mqcloud *MqcloudV1) DeleteVirtualPrivateEndpointGateway(deleteVirtualPrivateEndpointGatewayOptions *DeleteVirtualPrivateEndpointGatewayOptions) (response *core.DetailedResponse, err error) {
	response, err = mqcloud.DeleteVirtualPrivateEndpointGatewayWithContext(context.Background(), deleteVirtualPrivateEndpointGatewayOptions)
	err = core.RepurposeSDKProblem(err, "")
	return
}

// DeleteVirtualPrivateEndpointGatewayWithContext is an alternate form of the DeleteVirtualPrivateEndpointGateway method which supports a Context parameter
func (mqcloud *MqcloudV1) DeleteVirtualPrivateEndpointGatewayWithContext(ctx context.Context, deleteVirtualPrivateEndpointGatewayOptions *DeleteVirtualPrivateEndpointGatewayOptions) (response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteVirtualPrivateEndpointGatewayOptions, "deleteVirtualPrivateEndpointGatewayOptions cannot be nil")
	if err != nil {
		err = core.SDKErrorf(err, "", "unexpected-nil-param", common.GetComponentInfo())
		return
	}
	err = core.ValidateStruct(deleteVirtualPrivateEndpointGatewayOptions, "deleteVirtualPrivateEndpointGatewayOptions")
	if err != nil {
		err = core.SDKErrorf(err, "", "struct-validation-error", common.GetComponentInfo())
		return
	}

	pathParamsMap := map[string]string{
		"service_instance_guid":                 *deleteVirtualPrivateEndpointGatewayOptions.ServiceInstanceGuid,
		"virtual_private_endpoint_gateway_guid": *deleteVirtualPrivateEndpointGatewayOptions.VirtualPrivateEndpointGatewayGuid,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = mqcloud.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(mqcloud.Service.Options.URL, `/v1/{service_instance_guid}/virtual_private_endpoint_gateway/{virtual_private_endpoint_gateway_guid}`, pathParamsMap)
	if err != nil {
		err = core.SDKErrorf(err, "", "url-resolve-error", common.GetComponentInfo())
		return
	}

	for headerName, headerValue := range deleteVirtualPrivateEndpointGatewayOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("mqcloud", "V1", "DeleteVirtualPrivateEndpointGateway")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	if mqcloud.AcceptLanguage != nil {
		builder.AddHeader("Accept-Language", fmt.Sprint(*mqcloud.AcceptLanguage))
	}
	if deleteVirtualPrivateEndpointGatewayOptions.TrustedProfile != nil {
		builder.AddHeader("Trusted-Profile", fmt.Sprint(*deleteVirtualPrivateEndpointGatewayOptions.TrustedProfile))
	}

	request, err := builder.Build()
	if err != nil {
		err = core.SDKErrorf(err, "", "build-error", common.GetComponentInfo())
		return
	}

	response, err = mqcloud.Service.Request(request, nil)
	if err != nil {
		core.EnrichHTTPProblem(err, "delete_virtual_private_endpoint_gateway", getServiceComponentInfo())
		err = core.SDKErrorf(err, "", "http-request-err", common.GetComponentInfo())
		return
	}

	return
}
func getServiceComponentInfo() *core.ProblemComponent {
	return core.NewProblemComponent(DefaultServiceName, "1.1.0")
}

// ApplicationAPIKeyCreated : A response to creating a new api key, giving the only chance to collect the new apikey.
type ApplicationAPIKeyCreated struct {
	// The name of the api key.
	ApiKeyName *string `json:"api_key_name,omitempty"`

	// The id of the api key.
	ApiKeyID *string `json:"api_key_id" validate:"required"`

	// The api key created.
	ApiKey *string `json:"api_key" validate:"required"`
}

// UnmarshalApplicationAPIKeyCreated unmarshals an instance of ApplicationAPIKeyCreated from the specified map of raw messages.
func UnmarshalApplicationAPIKeyCreated(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApplicationAPIKeyCreated)
	err = core.UnmarshalPrimitive(m, "api_key_name", &obj.ApiKeyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.ApiKeyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key_id-error", common.GetComponentInfo())
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

// ApplicationCreated : A response to creating an application, giving the only chance to collect the apikey.
type ApplicationCreated struct {
	// The ID of the application which was allocated on creation, and can be used for delete calls.
	ID *string `json:"id" validate:"required"`

	// The name of the application - conforming to MQ rules.
	Name *string `json:"name" validate:"required"`

	// The URI to create a new apikey for the application.
	CreateApiKeyURI *string `json:"create_api_key_uri" validate:"required"`

	// The URL for this application.
	Href *string `json:"href" validate:"required"`

	// The name of the api key.
	ApiKeyName *string `json:"api_key_name,omitempty"`

	// The id of the api key.
	ApiKeyID *string `json:"api_key_id,omitempty"`

	// The api key created.
	ApiKey *string `json:"api_key" validate:"required"`
}

// UnmarshalApplicationCreated unmarshals an instance of ApplicationCreated from the specified map of raw messages.
func UnmarshalApplicationCreated(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApplicationCreated)
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
	err = core.UnmarshalPrimitive(m, "create_api_key_uri", &obj.CreateApiKeyURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "create_api_key_uri-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_name", &obj.ApiKeyName)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "api_key_id", &obj.ApiKeyID)
	if err != nil {
		err = core.SDKErrorf(err, "", "api_key_id-error", common.GetComponentInfo())
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

// ApplicationDetails : A summary of the application for use in a list of applications.
type ApplicationDetails struct {
	// The ID of the application which was allocated on creation, and can be used for delete calls.
	ID *string `json:"id" validate:"required"`

	// The name of the application - conforming to MQ rules.
	Name *string `json:"name" validate:"required"`

	// The URI to create a new apikey for the application.
	CreateApiKeyURI *string `json:"create_api_key_uri" validate:"required"`

	// The URL for this application.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalApplicationDetails unmarshals an instance of ApplicationDetails from the specified map of raw messages.
func UnmarshalApplicationDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApplicationDetails)
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
	err = core.UnmarshalPrimitive(m, "create_api_key_uri", &obj.CreateApiKeyURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "create_api_key_uri-error", common.GetComponentInfo())
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

// ApplicationDetailsCollection : A list of application summaries.
type ApplicationDetailsCollection struct {
	// Pagination offset.
	Offset *int64 `json:"offset" validate:"required"`

	// Results per page, same for all collections.
	Limit *int64 `json:"limit" validate:"required"`

	// Link to first page of results.
	First *First `json:"first" validate:"required"`

	// Link to next page of results.
	Next *Next `json:"next,omitempty"`

	// Link to previous page of results.
	Previous *Previous `json:"previous,omitempty"`

	// List of applications.
	Applications []ApplicationDetails `json:"applications" validate:"required"`
}

// UnmarshalApplicationDetailsCollection unmarshals an instance of ApplicationDetailsCollection from the specified map of raw messages.
func UnmarshalApplicationDetailsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ApplicationDetailsCollection)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "applications", &obj.Applications, UnmarshalApplicationDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "applications-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *ApplicationDetailsCollection) GetNextOffset() (*int64, error) {
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

// CertificateConfiguration : The configuration details for this certificate.
type CertificateConfiguration struct {
	// A list of channels that are configured with this certificate.
	Ams *ChannelsDetails `json:"ams" validate:"required"`
}

// UnmarshalCertificateConfiguration unmarshals an instance of CertificateConfiguration from the specified map of raw messages.
func UnmarshalCertificateConfiguration(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(CertificateConfiguration)
	err = core.UnmarshalModel(m, "ams", &obj.Ams, UnmarshalChannelsDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "ams-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChannelDetails : A channel's information that is configured with this certificate.
type ChannelDetails struct {
	// The name of the channel.
	Name *string `json:"name,omitempty"`
}

// UnmarshalChannelDetails unmarshals an instance of ChannelDetails from the specified map of raw messages.
func UnmarshalChannelDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ChannelDetails)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ChannelsDetails : A list of channels that are configured with this certificate.
type ChannelsDetails struct {
	// A list of channels that are configured with this certificate.
	Channels []ChannelDetails `json:"channels" validate:"required"`
}

// UnmarshalChannelsDetails unmarshals an instance of ChannelsDetails from the specified map of raw messages.
func UnmarshalChannelsDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ChannelsDetails)
	err = core.UnmarshalModel(m, "channels", &obj.Channels, UnmarshalChannelDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "channels-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ClientConnection : Details for a client connection.
type ClientConnection struct {
	// A collection of objects with attributes that define a channel connection.
	Connection []ConnectionDetails `json:"connection,omitempty"`

	// the name of the queue_manager.
	QueueManager *string `json:"queueManager,omitempty"`
}

// UnmarshalClientConnection unmarshals an instance of ClientConnection from the specified map of raw messages.
func UnmarshalClientConnection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ClientConnection)
	err = core.UnmarshalModel(m, "connection", &obj.Connection, UnmarshalConnectionDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "connection-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "queueManager", &obj.QueueManager)
	if err != nil {
		err = core.SDKErrorf(err, "", "queueManager-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConfigurationOptions : Configuration options (eg, available deployment locations, queue manager sizes).
type ConfigurationOptions struct {
	// List of deployment locations.
	Locations []string `json:"locations,omitempty"`

	// List of queue manager sizes.
	Sizes []string `json:"sizes,omitempty"`

	// List of queue manager versions.
	Versions []string `json:"versions,omitempty"`

	// The latest Queue manager version.
	LatestVersion *string `json:"latest_version,omitempty"`
}

// Constants associated with the ConfigurationOptions.Sizes property.
// The queue manager sizes of deployment available.
const (
	ConfigurationOptions_Sizes_Large  = "large"
	ConfigurationOptions_Sizes_Medium = "medium"
	ConfigurationOptions_Sizes_Small  = "small"
	ConfigurationOptions_Sizes_Xsmall = "xsmall"
)

// UnmarshalConfigurationOptions unmarshals an instance of ConfigurationOptions from the specified map of raw messages.
func UnmarshalConfigurationOptions(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConfigurationOptions)
	err = core.UnmarshalPrimitive(m, "locations", &obj.Locations)
	if err != nil {
		err = core.SDKErrorf(err, "", "locations-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "sizes", &obj.Sizes)
	if err != nil {
		err = core.SDKErrorf(err, "", "sizes-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "versions", &obj.Versions)
	if err != nil {
		err = core.SDKErrorf(err, "", "versions-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "latest_version", &obj.LatestVersion)
	if err != nil {
		err = core.SDKErrorf(err, "", "latest_version-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConnectionDetails : Attributes that define a channel connection.
type ConnectionDetails struct {
	// Specifies the host that this channel connects to.
	Host *string `json:"host,omitempty"`

	// Specifies the port that this channel uses on this host.
	Port *int64 `json:"port,omitempty"`
}

// UnmarshalConnectionDetails unmarshals an instance of ConnectionDetails from the specified map of raw messages.
func UnmarshalConnectionDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConnectionDetails)
	err = core.UnmarshalPrimitive(m, "host", &obj.Host)
	if err != nil {
		err = core.SDKErrorf(err, "", "host-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		err = core.SDKErrorf(err, "", "port-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConnectionInfo : Responds with JSON CCDT of the connection information for the queue manager.
type ConnectionInfo struct {
	// A collection of channel connection details.
	Channel []ConnectionInfoChannel `json:"channel" validate:"required"`
}

// UnmarshalConnectionInfo unmarshals an instance of ConnectionInfo from the specified map of raw messages.
func UnmarshalConnectionInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConnectionInfo)
	err = core.UnmarshalModel(m, "channel", &obj.Channel, UnmarshalConnectionInfoChannel)
	if err != nil {
		err = core.SDKErrorf(err, "", "channel-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ConnectionInfoChannel : A subsection for channels as part of a JSON CCDT of the connection information for the queue manager.
type ConnectionInfoChannel struct {
	// Specifies the name of the channel.
	Name *string `json:"name" validate:"required"`

	// Details for a client connection.
	ClientConnection *ClientConnection `json:"clientConnection" validate:"required"`

	// An object that contains attributes that are related to security for message transmission.
	TransmissionSecurity *TransmissionSecurity `json:"transmissionSecurity" validate:"required"`

	// Specifies the type of the channel.
	Type *string `json:"type" validate:"required"`
}

// UnmarshalConnectionInfoChannel unmarshals an instance of ConnectionInfoChannel from the specified map of raw messages.
func UnmarshalConnectionInfoChannel(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ConnectionInfoChannel)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		err = core.SDKErrorf(err, "", "name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "clientConnection", &obj.ClientConnection, UnmarshalClientConnection)
	if err != nil {
		err = core.SDKErrorf(err, "", "clientConnection-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "transmissionSecurity", &obj.TransmissionSecurity, UnmarshalTransmissionSecurity)
	if err != nil {
		err = core.SDKErrorf(err, "", "transmissionSecurity-error", common.GetComponentInfo())
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

// CreateApplicationApikeyOptions : The CreateApplicationApikey options.
type CreateApplicationApikeyOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the application.
	ApplicationID *string `json:"application_id" validate:"required,ne="`

	// The short name of the application api key - conforming to MQ rules.
	Name *string `json:"name" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateApplicationApikeyOptions : Instantiate CreateApplicationApikeyOptions
func (*MqcloudV1) NewCreateApplicationApikeyOptions(serviceInstanceGuid string, applicationID string, name string) *CreateApplicationApikeyOptions {
	return &CreateApplicationApikeyOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		ApplicationID:       core.StringPtr(applicationID),
		Name:                core.StringPtr(name),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateApplicationApikeyOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateApplicationApikeyOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetApplicationID : Allow user to set ApplicationID
func (_options *CreateApplicationApikeyOptions) SetApplicationID(applicationID string) *CreateApplicationApikeyOptions {
	_options.ApplicationID = core.StringPtr(applicationID)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateApplicationApikeyOptions) SetName(name string) *CreateApplicationApikeyOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateApplicationApikeyOptions) SetHeaders(param map[string]string) *CreateApplicationApikeyOptions {
	options.Headers = param
	return options
}

// CreateApplicationOptions : The CreateApplication options.
type CreateApplicationOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The name of the application - conforming to MQ rules.
	Name *string `json:"name" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateApplicationOptions : Instantiate CreateApplicationOptions
func (*MqcloudV1) NewCreateApplicationOptions(serviceInstanceGuid string, name string) *CreateApplicationOptions {
	return &CreateApplicationOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		Name:                core.StringPtr(name),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateApplicationOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateApplicationOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateApplicationOptions) SetName(name string) *CreateApplicationOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateApplicationOptions) SetHeaders(param map[string]string) *CreateApplicationOptions {
	options.Headers = param
	return options
}

// CreateKeyStorePemCertificateOptions : The CreateKeyStorePemCertificate options.
type CreateKeyStorePemCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The label to use for the certificate to be uploaded.
	Label *string `json:"label" validate:"required"`

	// The filename and path of the certificate to be uploaded.
	CertificateFile io.ReadCloser `json:"certificate_file" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateKeyStorePemCertificateOptions : Instantiate CreateKeyStorePemCertificateOptions
func (*MqcloudV1) NewCreateKeyStorePemCertificateOptions(serviceInstanceGuid string, queueManagerID string, label string, certificateFile io.ReadCloser) *CreateKeyStorePemCertificateOptions {
	return &CreateKeyStorePemCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		Label:               core.StringPtr(label),
		CertificateFile:     certificateFile,
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateKeyStorePemCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateKeyStorePemCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *CreateKeyStorePemCertificateOptions) SetQueueManagerID(queueManagerID string) *CreateKeyStorePemCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateKeyStorePemCertificateOptions) SetLabel(label string) *CreateKeyStorePemCertificateOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetCertificateFile : Allow user to set CertificateFile
func (_options *CreateKeyStorePemCertificateOptions) SetCertificateFile(certificateFile io.ReadCloser) *CreateKeyStorePemCertificateOptions {
	_options.CertificateFile = certificateFile
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateKeyStorePemCertificateOptions) SetHeaders(param map[string]string) *CreateKeyStorePemCertificateOptions {
	options.Headers = param
	return options
}

// CreateQueueManagerOptions : The CreateQueueManager options.
type CreateQueueManagerOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The name of the queue manager - conforming to MQ rules.
	Name *string `json:"name" validate:"required"`

	// The locations in which the queue manager could be deployed.
	Location *string `json:"location" validate:"required"`

	// The queue manager sizes of deployment available.
	Size *string `json:"size" validate:"required"`

	// A displayable name for the queue manager - limited only in length.
	DisplayName *string `json:"display_name,omitempty"`

	// The IBM MQ version of the Queue Manager to deploy if not supplied the latest version will be deployed.
	Version *string `json:"version,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the CreateQueueManagerOptions.Size property.
// The queue manager sizes of deployment available.
const (
	CreateQueueManagerOptions_Size_Large  = "large"
	CreateQueueManagerOptions_Size_Medium = "medium"
	CreateQueueManagerOptions_Size_Small  = "small"
	CreateQueueManagerOptions_Size_Xsmall = "xsmall"
)

// NewCreateQueueManagerOptions : Instantiate CreateQueueManagerOptions
func (*MqcloudV1) NewCreateQueueManagerOptions(serviceInstanceGuid string, name string, location string, size string) *CreateQueueManagerOptions {
	return &CreateQueueManagerOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		Name:                core.StringPtr(name),
		Location:            core.StringPtr(location),
		Size:                core.StringPtr(size),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateQueueManagerOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateQueueManagerOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateQueueManagerOptions) SetName(name string) *CreateQueueManagerOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetLocation : Allow user to set Location
func (_options *CreateQueueManagerOptions) SetLocation(location string) *CreateQueueManagerOptions {
	_options.Location = core.StringPtr(location)
	return _options
}

// SetSize : Allow user to set Size
func (_options *CreateQueueManagerOptions) SetSize(size string) *CreateQueueManagerOptions {
	_options.Size = core.StringPtr(size)
	return _options
}

// SetDisplayName : Allow user to set DisplayName
func (_options *CreateQueueManagerOptions) SetDisplayName(displayName string) *CreateQueueManagerOptions {
	_options.DisplayName = core.StringPtr(displayName)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *CreateQueueManagerOptions) SetVersion(version string) *CreateQueueManagerOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateQueueManagerOptions) SetHeaders(param map[string]string) *CreateQueueManagerOptions {
	options.Headers = param
	return options
}

// CreateTrustStorePemCertificateOptions : The CreateTrustStorePemCertificate options.
type CreateTrustStorePemCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The label to use for the certificate to be uploaded.
	Label *string `json:"label" validate:"required"`

	// The filename and path of the certificate to be uploaded.
	CertificateFile io.ReadCloser `json:"certificate_file" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateTrustStorePemCertificateOptions : Instantiate CreateTrustStorePemCertificateOptions
func (*MqcloudV1) NewCreateTrustStorePemCertificateOptions(serviceInstanceGuid string, queueManagerID string, label string, certificateFile io.ReadCloser) *CreateTrustStorePemCertificateOptions {
	return &CreateTrustStorePemCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		Label:               core.StringPtr(label),
		CertificateFile:     certificateFile,
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateTrustStorePemCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateTrustStorePemCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *CreateTrustStorePemCertificateOptions) SetQueueManagerID(queueManagerID string) *CreateTrustStorePemCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetLabel : Allow user to set Label
func (_options *CreateTrustStorePemCertificateOptions) SetLabel(label string) *CreateTrustStorePemCertificateOptions {
	_options.Label = core.StringPtr(label)
	return _options
}

// SetCertificateFile : Allow user to set CertificateFile
func (_options *CreateTrustStorePemCertificateOptions) SetCertificateFile(certificateFile io.ReadCloser) *CreateTrustStorePemCertificateOptions {
	_options.CertificateFile = certificateFile
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateTrustStorePemCertificateOptions) SetHeaders(param map[string]string) *CreateTrustStorePemCertificateOptions {
	options.Headers = param
	return options
}

// CreateUserOptions : The CreateUser options.
type CreateUserOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The email of the user to be created.
	Email *string `json:"email" validate:"required"`

	// The shortname of the user to be created.
	Name *string `json:"name" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateUserOptions : Instantiate CreateUserOptions
func (*MqcloudV1) NewCreateUserOptions(serviceInstanceGuid string, email string, name string) *CreateUserOptions {
	return &CreateUserOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		Email:               core.StringPtr(email),
		Name:                core.StringPtr(name),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateUserOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateUserOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetEmail : Allow user to set Email
func (_options *CreateUserOptions) SetEmail(email string) *CreateUserOptions {
	_options.Email = core.StringPtr(email)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateUserOptions) SetName(name string) *CreateUserOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateUserOptions) SetHeaders(param map[string]string) *CreateUserOptions {
	options.Headers = param
	return options
}

// CreateVirtualPrivateEndpointGatewayOptions : The CreateVirtualPrivateEndpointGateway options.
type CreateVirtualPrivateEndpointGatewayOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The name of the virtual private endpoint gateway - conforming to naming rules.
	Name *string `json:"name" validate:"required"`

	// The CRN of the target reserved capacity service instance.
	TargetCrn *string `json:"target_crn" validate:"required"`

	// The CRN of the trusted profile to assume for this request.
	TrustedProfile *string `json:"Trusted-Profile,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewCreateVirtualPrivateEndpointGatewayOptions : Instantiate CreateVirtualPrivateEndpointGatewayOptions
func (*MqcloudV1) NewCreateVirtualPrivateEndpointGatewayOptions(serviceInstanceGuid string, name string, targetCrn string) *CreateVirtualPrivateEndpointGatewayOptions {
	return &CreateVirtualPrivateEndpointGatewayOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		Name:                core.StringPtr(name),
		TargetCrn:           core.StringPtr(targetCrn),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *CreateVirtualPrivateEndpointGatewayOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *CreateVirtualPrivateEndpointGatewayOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetName : Allow user to set Name
func (_options *CreateVirtualPrivateEndpointGatewayOptions) SetName(name string) *CreateVirtualPrivateEndpointGatewayOptions {
	_options.Name = core.StringPtr(name)
	return _options
}

// SetTargetCrn : Allow user to set TargetCrn
func (_options *CreateVirtualPrivateEndpointGatewayOptions) SetTargetCrn(targetCrn string) *CreateVirtualPrivateEndpointGatewayOptions {
	_options.TargetCrn = core.StringPtr(targetCrn)
	return _options
}

// SetTrustedProfile : Allow user to set TrustedProfile
func (_options *CreateVirtualPrivateEndpointGatewayOptions) SetTrustedProfile(trustedProfile string) *CreateVirtualPrivateEndpointGatewayOptions {
	_options.TrustedProfile = core.StringPtr(trustedProfile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *CreateVirtualPrivateEndpointGatewayOptions) SetHeaders(param map[string]string) *CreateVirtualPrivateEndpointGatewayOptions {
	options.Headers = param
	return options
}

// DeleteApplicationOptions : The DeleteApplication options.
type DeleteApplicationOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the application.
	ApplicationID *string `json:"application_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteApplicationOptions : Instantiate DeleteApplicationOptions
func (*MqcloudV1) NewDeleteApplicationOptions(serviceInstanceGuid string, applicationID string) *DeleteApplicationOptions {
	return &DeleteApplicationOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		ApplicationID:       core.StringPtr(applicationID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DeleteApplicationOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DeleteApplicationOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetApplicationID : Allow user to set ApplicationID
func (_options *DeleteApplicationOptions) SetApplicationID(applicationID string) *DeleteApplicationOptions {
	_options.ApplicationID = core.StringPtr(applicationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteApplicationOptions) SetHeaders(param map[string]string) *DeleteApplicationOptions {
	options.Headers = param
	return options
}

// DeleteKeyStoreCertificateOptions : The DeleteKeyStoreCertificate options.
type DeleteKeyStoreCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteKeyStoreCertificateOptions : Instantiate DeleteKeyStoreCertificateOptions
func (*MqcloudV1) NewDeleteKeyStoreCertificateOptions(serviceInstanceGuid string, queueManagerID string, certificateID string) *DeleteKeyStoreCertificateOptions {
	return &DeleteKeyStoreCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DeleteKeyStoreCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DeleteKeyStoreCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *DeleteKeyStoreCertificateOptions) SetQueueManagerID(queueManagerID string) *DeleteKeyStoreCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *DeleteKeyStoreCertificateOptions) SetCertificateID(certificateID string) *DeleteKeyStoreCertificateOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteKeyStoreCertificateOptions) SetHeaders(param map[string]string) *DeleteKeyStoreCertificateOptions {
	options.Headers = param
	return options
}

// DeleteQueueManagerOptions : The DeleteQueueManager options.
type DeleteQueueManagerOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteQueueManagerOptions : Instantiate DeleteQueueManagerOptions
func (*MqcloudV1) NewDeleteQueueManagerOptions(serviceInstanceGuid string, queueManagerID string) *DeleteQueueManagerOptions {
	return &DeleteQueueManagerOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DeleteQueueManagerOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DeleteQueueManagerOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *DeleteQueueManagerOptions) SetQueueManagerID(queueManagerID string) *DeleteQueueManagerOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteQueueManagerOptions) SetHeaders(param map[string]string) *DeleteQueueManagerOptions {
	options.Headers = param
	return options
}

// DeleteTrustStoreCertificateOptions : The DeleteTrustStoreCertificate options.
type DeleteTrustStoreCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteTrustStoreCertificateOptions : Instantiate DeleteTrustStoreCertificateOptions
func (*MqcloudV1) NewDeleteTrustStoreCertificateOptions(serviceInstanceGuid string, queueManagerID string, certificateID string) *DeleteTrustStoreCertificateOptions {
	return &DeleteTrustStoreCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DeleteTrustStoreCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DeleteTrustStoreCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *DeleteTrustStoreCertificateOptions) SetQueueManagerID(queueManagerID string) *DeleteTrustStoreCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *DeleteTrustStoreCertificateOptions) SetCertificateID(certificateID string) *DeleteTrustStoreCertificateOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteTrustStoreCertificateOptions) SetHeaders(param map[string]string) *DeleteTrustStoreCertificateOptions {
	options.Headers = param
	return options
}

// DeleteUserOptions : The DeleteUser options.
type DeleteUserOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the user.
	UserID *string `json:"user_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteUserOptions : Instantiate DeleteUserOptions
func (*MqcloudV1) NewDeleteUserOptions(serviceInstanceGuid string, userID string) *DeleteUserOptions {
	return &DeleteUserOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		UserID:              core.StringPtr(userID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DeleteUserOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DeleteUserOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetUserID : Allow user to set UserID
func (_options *DeleteUserOptions) SetUserID(userID string) *DeleteUserOptions {
	_options.UserID = core.StringPtr(userID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteUserOptions) SetHeaders(param map[string]string) *DeleteUserOptions {
	options.Headers = param
	return options
}

// DeleteVirtualPrivateEndpointGatewayOptions : The DeleteVirtualPrivateEndpointGateway options.
type DeleteVirtualPrivateEndpointGatewayOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the virtual private endpoint gateway.
	VirtualPrivateEndpointGatewayGuid *string `json:"virtual_private_endpoint_gateway_guid" validate:"required,ne="`

	// The CRN of the trusted profile to assume for this request.
	TrustedProfile *string `json:"Trusted-Profile,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDeleteVirtualPrivateEndpointGatewayOptions : Instantiate DeleteVirtualPrivateEndpointGatewayOptions
func (*MqcloudV1) NewDeleteVirtualPrivateEndpointGatewayOptions(serviceInstanceGuid string, virtualPrivateEndpointGatewayGuid string) *DeleteVirtualPrivateEndpointGatewayOptions {
	return &DeleteVirtualPrivateEndpointGatewayOptions{
		ServiceInstanceGuid:               core.StringPtr(serviceInstanceGuid),
		VirtualPrivateEndpointGatewayGuid: core.StringPtr(virtualPrivateEndpointGatewayGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DeleteVirtualPrivateEndpointGatewayOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DeleteVirtualPrivateEndpointGatewayOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetVirtualPrivateEndpointGatewayGuid : Allow user to set VirtualPrivateEndpointGatewayGuid
func (_options *DeleteVirtualPrivateEndpointGatewayOptions) SetVirtualPrivateEndpointGatewayGuid(virtualPrivateEndpointGatewayGuid string) *DeleteVirtualPrivateEndpointGatewayOptions {
	_options.VirtualPrivateEndpointGatewayGuid = core.StringPtr(virtualPrivateEndpointGatewayGuid)
	return _options
}

// SetTrustedProfile : Allow user to set TrustedProfile
func (_options *DeleteVirtualPrivateEndpointGatewayOptions) SetTrustedProfile(trustedProfile string) *DeleteVirtualPrivateEndpointGatewayOptions {
	_options.TrustedProfile = core.StringPtr(trustedProfile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteVirtualPrivateEndpointGatewayOptions) SetHeaders(param map[string]string) *DeleteVirtualPrivateEndpointGatewayOptions {
	options.Headers = param
	return options
}

// DownloadKeyStoreCertificateOptions : The DownloadKeyStoreCertificate options.
type DownloadKeyStoreCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDownloadKeyStoreCertificateOptions : Instantiate DownloadKeyStoreCertificateOptions
func (*MqcloudV1) NewDownloadKeyStoreCertificateOptions(serviceInstanceGuid string, queueManagerID string, certificateID string) *DownloadKeyStoreCertificateOptions {
	return &DownloadKeyStoreCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DownloadKeyStoreCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DownloadKeyStoreCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *DownloadKeyStoreCertificateOptions) SetQueueManagerID(queueManagerID string) *DownloadKeyStoreCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *DownloadKeyStoreCertificateOptions) SetCertificateID(certificateID string) *DownloadKeyStoreCertificateOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DownloadKeyStoreCertificateOptions) SetHeaders(param map[string]string) *DownloadKeyStoreCertificateOptions {
	options.Headers = param
	return options
}

// DownloadTrustStoreCertificateOptions : The DownloadTrustStoreCertificate options.
type DownloadTrustStoreCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewDownloadTrustStoreCertificateOptions : Instantiate DownloadTrustStoreCertificateOptions
func (*MqcloudV1) NewDownloadTrustStoreCertificateOptions(serviceInstanceGuid string, queueManagerID string, certificateID string) *DownloadTrustStoreCertificateOptions {
	return &DownloadTrustStoreCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *DownloadTrustStoreCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *DownloadTrustStoreCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *DownloadTrustStoreCertificateOptions) SetQueueManagerID(queueManagerID string) *DownloadTrustStoreCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *DownloadTrustStoreCertificateOptions) SetCertificateID(certificateID string) *DownloadTrustStoreCertificateOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *DownloadTrustStoreCertificateOptions) SetHeaders(param map[string]string) *DownloadTrustStoreCertificateOptions {
	options.Headers = param
	return options
}

// First : Link to first page of results.
type First struct {
	// The URL of the page the link goes to.
	Href *string `json:"href,omitempty"`
}

// UnmarshalFirst unmarshals an instance of First from the specified map of raw messages.
func UnmarshalFirst(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(First)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetApplicationOptions : The GetApplication options.
type GetApplicationOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the application.
	ApplicationID *string `json:"application_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetApplicationOptions : Instantiate GetApplicationOptions
func (*MqcloudV1) NewGetApplicationOptions(serviceInstanceGuid string, applicationID string) *GetApplicationOptions {
	return &GetApplicationOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		ApplicationID:       core.StringPtr(applicationID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetApplicationOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetApplicationOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetApplicationID : Allow user to set ApplicationID
func (_options *GetApplicationOptions) SetApplicationID(applicationID string) *GetApplicationOptions {
	_options.ApplicationID = core.StringPtr(applicationID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetApplicationOptions) SetHeaders(param map[string]string) *GetApplicationOptions {
	options.Headers = param
	return options
}

// GetCertificateAmsChannelsOptions : The GetCertificateAmsChannels options.
type GetCertificateAmsChannelsOptions struct {
	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetCertificateAmsChannelsOptions : Instantiate GetCertificateAmsChannelsOptions
func (*MqcloudV1) NewGetCertificateAmsChannelsOptions(queueManagerID string, certificateID string, serviceInstanceGuid string) *GetCertificateAmsChannelsOptions {
	return &GetCertificateAmsChannelsOptions{
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetCertificateAmsChannelsOptions) SetQueueManagerID(queueManagerID string) *GetCertificateAmsChannelsOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *GetCertificateAmsChannelsOptions) SetCertificateID(certificateID string) *GetCertificateAmsChannelsOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetCertificateAmsChannelsOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetCertificateAmsChannelsOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetCertificateAmsChannelsOptions) SetHeaders(param map[string]string) *GetCertificateAmsChannelsOptions {
	options.Headers = param
	return options
}

// GetKeyStoreCertificateOptions : The GetKeyStoreCertificate options.
type GetKeyStoreCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetKeyStoreCertificateOptions : Instantiate GetKeyStoreCertificateOptions
func (*MqcloudV1) NewGetKeyStoreCertificateOptions(serviceInstanceGuid string, queueManagerID string, certificateID string) *GetKeyStoreCertificateOptions {
	return &GetKeyStoreCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetKeyStoreCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetKeyStoreCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetKeyStoreCertificateOptions) SetQueueManagerID(queueManagerID string) *GetKeyStoreCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *GetKeyStoreCertificateOptions) SetCertificateID(certificateID string) *GetKeyStoreCertificateOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetKeyStoreCertificateOptions) SetHeaders(param map[string]string) *GetKeyStoreCertificateOptions {
	options.Headers = param
	return options
}

// GetOptionsOptions : The GetOptions options.
type GetOptionsOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetOptionsOptions : Instantiate GetOptionsOptions
func (*MqcloudV1) NewGetOptionsOptions(serviceInstanceGuid string) *GetOptionsOptions {
	return &GetOptionsOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetOptionsOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetOptionsOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetOptionsOptions) SetHeaders(param map[string]string) *GetOptionsOptions {
	options.Headers = param
	return options
}

// GetQueueManagerAvailableUpgradeVersionsOptions : The GetQueueManagerAvailableUpgradeVersions options.
type GetQueueManagerAvailableUpgradeVersionsOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetQueueManagerAvailableUpgradeVersionsOptions : Instantiate GetQueueManagerAvailableUpgradeVersionsOptions
func (*MqcloudV1) NewGetQueueManagerAvailableUpgradeVersionsOptions(serviceInstanceGuid string, queueManagerID string) *GetQueueManagerAvailableUpgradeVersionsOptions {
	return &GetQueueManagerAvailableUpgradeVersionsOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetQueueManagerAvailableUpgradeVersionsOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetQueueManagerAvailableUpgradeVersionsOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetQueueManagerAvailableUpgradeVersionsOptions) SetQueueManagerID(queueManagerID string) *GetQueueManagerAvailableUpgradeVersionsOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetQueueManagerAvailableUpgradeVersionsOptions) SetHeaders(param map[string]string) *GetQueueManagerAvailableUpgradeVersionsOptions {
	options.Headers = param
	return options
}

// GetQueueManagerConnectionInfoOptions : The GetQueueManagerConnectionInfo options.
type GetQueueManagerConnectionInfoOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetQueueManagerConnectionInfoOptions : Instantiate GetQueueManagerConnectionInfoOptions
func (*MqcloudV1) NewGetQueueManagerConnectionInfoOptions(serviceInstanceGuid string, queueManagerID string) *GetQueueManagerConnectionInfoOptions {
	return &GetQueueManagerConnectionInfoOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetQueueManagerConnectionInfoOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetQueueManagerConnectionInfoOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetQueueManagerConnectionInfoOptions) SetQueueManagerID(queueManagerID string) *GetQueueManagerConnectionInfoOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetQueueManagerConnectionInfoOptions) SetHeaders(param map[string]string) *GetQueueManagerConnectionInfoOptions {
	options.Headers = param
	return options
}

// GetQueueManagerOptions : The GetQueueManager options.
type GetQueueManagerOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetQueueManagerOptions : Instantiate GetQueueManagerOptions
func (*MqcloudV1) NewGetQueueManagerOptions(serviceInstanceGuid string, queueManagerID string) *GetQueueManagerOptions {
	return &GetQueueManagerOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetQueueManagerOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetQueueManagerOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetQueueManagerOptions) SetQueueManagerID(queueManagerID string) *GetQueueManagerOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetQueueManagerOptions) SetHeaders(param map[string]string) *GetQueueManagerOptions {
	options.Headers = param
	return options
}

// GetQueueManagerStatusOptions : The GetQueueManagerStatus options.
type GetQueueManagerStatusOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetQueueManagerStatusOptions : Instantiate GetQueueManagerStatusOptions
func (*MqcloudV1) NewGetQueueManagerStatusOptions(serviceInstanceGuid string, queueManagerID string) *GetQueueManagerStatusOptions {
	return &GetQueueManagerStatusOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetQueueManagerStatusOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetQueueManagerStatusOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetQueueManagerStatusOptions) SetQueueManagerID(queueManagerID string) *GetQueueManagerStatusOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetQueueManagerStatusOptions) SetHeaders(param map[string]string) *GetQueueManagerStatusOptions {
	options.Headers = param
	return options
}

// GetTrustStoreCertificateOptions : The GetTrustStoreCertificate options.
type GetTrustStoreCertificateOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetTrustStoreCertificateOptions : Instantiate GetTrustStoreCertificateOptions
func (*MqcloudV1) NewGetTrustStoreCertificateOptions(serviceInstanceGuid string, queueManagerID string, certificateID string) *GetTrustStoreCertificateOptions {
	return &GetTrustStoreCertificateOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetTrustStoreCertificateOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetTrustStoreCertificateOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *GetTrustStoreCertificateOptions) SetQueueManagerID(queueManagerID string) *GetTrustStoreCertificateOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *GetTrustStoreCertificateOptions) SetCertificateID(certificateID string) *GetTrustStoreCertificateOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetTrustStoreCertificateOptions) SetHeaders(param map[string]string) *GetTrustStoreCertificateOptions {
	options.Headers = param
	return options
}

// GetUsageDetailsOptions : The GetUsageDetails options.
type GetUsageDetailsOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetUsageDetailsOptions : Instantiate GetUsageDetailsOptions
func (*MqcloudV1) NewGetUsageDetailsOptions(serviceInstanceGuid string) *GetUsageDetailsOptions {
	return &GetUsageDetailsOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetUsageDetailsOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetUsageDetailsOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetUsageDetailsOptions) SetHeaders(param map[string]string) *GetUsageDetailsOptions {
	options.Headers = param
	return options
}

// GetUserOptions : The GetUser options.
type GetUserOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the user.
	UserID *string `json:"user_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetUserOptions : Instantiate GetUserOptions
func (*MqcloudV1) NewGetUserOptions(serviceInstanceGuid string, userID string) *GetUserOptions {
	return &GetUserOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		UserID:              core.StringPtr(userID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetUserOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetUserOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetUserID : Allow user to set UserID
func (_options *GetUserOptions) SetUserID(userID string) *GetUserOptions {
	_options.UserID = core.StringPtr(userID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetUserOptions) SetHeaders(param map[string]string) *GetUserOptions {
	options.Headers = param
	return options
}

// GetVirtualPrivateEndpointGatewayOptions : The GetVirtualPrivateEndpointGateway options.
type GetVirtualPrivateEndpointGatewayOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the virtual private endpoint gateway.
	VirtualPrivateEndpointGatewayGuid *string `json:"virtual_private_endpoint_gateway_guid" validate:"required,ne="`

	// The CRN of the trusted profile to assume for this request.
	TrustedProfile *string `json:"Trusted-Profile,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewGetVirtualPrivateEndpointGatewayOptions : Instantiate GetVirtualPrivateEndpointGatewayOptions
func (*MqcloudV1) NewGetVirtualPrivateEndpointGatewayOptions(serviceInstanceGuid string, virtualPrivateEndpointGatewayGuid string) *GetVirtualPrivateEndpointGatewayOptions {
	return &GetVirtualPrivateEndpointGatewayOptions{
		ServiceInstanceGuid:               core.StringPtr(serviceInstanceGuid),
		VirtualPrivateEndpointGatewayGuid: core.StringPtr(virtualPrivateEndpointGatewayGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *GetVirtualPrivateEndpointGatewayOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *GetVirtualPrivateEndpointGatewayOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetVirtualPrivateEndpointGatewayGuid : Allow user to set VirtualPrivateEndpointGatewayGuid
func (_options *GetVirtualPrivateEndpointGatewayOptions) SetVirtualPrivateEndpointGatewayGuid(virtualPrivateEndpointGatewayGuid string) *GetVirtualPrivateEndpointGatewayOptions {
	_options.VirtualPrivateEndpointGatewayGuid = core.StringPtr(virtualPrivateEndpointGatewayGuid)
	return _options
}

// SetTrustedProfile : Allow user to set TrustedProfile
func (_options *GetVirtualPrivateEndpointGatewayOptions) SetTrustedProfile(trustedProfile string) *GetVirtualPrivateEndpointGatewayOptions {
	_options.TrustedProfile = core.StringPtr(trustedProfile)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *GetVirtualPrivateEndpointGatewayOptions) SetHeaders(param map[string]string) *GetVirtualPrivateEndpointGatewayOptions {
	options.Headers = param
	return options
}

// KeyStoreCertificateDetails : The details of a key store certificate in a queue manager certificate key store.
type KeyStoreCertificateDetails struct {
	// ID of the certificate.
	ID *string `json:"id" validate:"required"`

	// Certificate label in queue manager store.
	Label *string `json:"label" validate:"required"`

	// The type of certificate.
	CertificateType *string `json:"certificate_type" validate:"required"`

	// Fingerprint SHA256.
	FingerprintSha256 *string `json:"fingerprint_sha256" validate:"required"`

	// Subject's Distinguished Name.
	SubjectDn *string `json:"subject_dn" validate:"required"`

	// Subject's Common Name.
	SubjectCn *string `json:"subject_cn" validate:"required"`

	// Issuer's Distinguished Name.
	IssuerDn *string `json:"issuer_dn" validate:"required"`

	// Issuer's Common Name.
	IssuerCn *string `json:"issuer_cn" validate:"required"`

	// Date certificate was issued.
	Issued *strfmt.DateTime `json:"issued" validate:"required"`

	// Expiry date for the certificate.
	Expiry *strfmt.DateTime `json:"expiry" validate:"required"`

	// Indicates whether it is the queue manager's default certificate.
	IsDefault *bool `json:"is_default" validate:"required"`

	// The total count of dns names.
	DnsNamesTotalCount *int64 `json:"dns_names_total_count" validate:"required"`

	// The list of DNS names.
	DnsNames []string `json:"dns_names" validate:"required"`

	// The URL for this key store certificate.
	Href *string `json:"href" validate:"required"`

	// The configuration details for this certificate.
	Config *CertificateConfiguration `json:"config" validate:"required"`
}

// Constants associated with the KeyStoreCertificateDetails.CertificateType property.
// The type of certificate.
const (
	KeyStoreCertificateDetails_CertificateType_KeyStore = "key_store"
)

// UnmarshalKeyStoreCertificateDetails unmarshals an instance of KeyStoreCertificateDetails from the specified map of raw messages.
func UnmarshalKeyStoreCertificateDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyStoreCertificateDetails)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_type", &obj.CertificateType)
	if err != nil {
		err = core.SDKErrorf(err, "", "certificate_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fingerprint_sha256", &obj.FingerprintSha256)
	if err != nil {
		err = core.SDKErrorf(err, "", "fingerprint_sha256-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "subject_dn", &obj.SubjectDn)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject_dn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "subject_cn", &obj.SubjectCn)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject_cn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer_dn", &obj.IssuerDn)
	if err != nil {
		err = core.SDKErrorf(err, "", "issuer_dn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer_cn", &obj.IssuerCn)
	if err != nil {
		err = core.SDKErrorf(err, "", "issuer_cn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "issued", &obj.Issued)
	if err != nil {
		err = core.SDKErrorf(err, "", "issued-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expiry", &obj.Expiry)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiry-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "is_default", &obj.IsDefault)
	if err != nil {
		err = core.SDKErrorf(err, "", "is_default-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dns_names_total_count", &obj.DnsNamesTotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "dns_names_total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "dns_names", &obj.DnsNames)
	if err != nil {
		err = core.SDKErrorf(err, "", "dns_names-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "config", &obj.Config, UnmarshalCertificateConfiguration)
	if err != nil {
		err = core.SDKErrorf(err, "", "config-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// KeyStoreCertificateDetailsCollection : A list of certificates in a queue manager's certificate key store.
type KeyStoreCertificateDetailsCollection struct {
	// The total count of key store certificates.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The list of key store certificates.
	KeyStore []KeyStoreCertificateDetails `json:"key_store,omitempty"`
}

// UnmarshalKeyStoreCertificateDetailsCollection unmarshals an instance of KeyStoreCertificateDetailsCollection from the specified map of raw messages.
func UnmarshalKeyStoreCertificateDetailsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(KeyStoreCertificateDetailsCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "key_store", &obj.KeyStore, UnmarshalKeyStoreCertificateDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "key_store-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListApplicationsOptions : The ListApplications options.
type ListApplicationsOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// Pagination offset.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of resources to return.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListApplicationsOptions : Instantiate ListApplicationsOptions
func (*MqcloudV1) NewListApplicationsOptions(serviceInstanceGuid string) *ListApplicationsOptions {
	return &ListApplicationsOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *ListApplicationsOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *ListApplicationsOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListApplicationsOptions) SetOffset(offset int64) *ListApplicationsOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListApplicationsOptions) SetLimit(limit int64) *ListApplicationsOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListApplicationsOptions) SetHeaders(param map[string]string) *ListApplicationsOptions {
	options.Headers = param
	return options
}

// ListKeyStoreCertificatesOptions : The ListKeyStoreCertificates options.
type ListKeyStoreCertificatesOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListKeyStoreCertificatesOptions : Instantiate ListKeyStoreCertificatesOptions
func (*MqcloudV1) NewListKeyStoreCertificatesOptions(serviceInstanceGuid string, queueManagerID string) *ListKeyStoreCertificatesOptions {
	return &ListKeyStoreCertificatesOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *ListKeyStoreCertificatesOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *ListKeyStoreCertificatesOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *ListKeyStoreCertificatesOptions) SetQueueManagerID(queueManagerID string) *ListKeyStoreCertificatesOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListKeyStoreCertificatesOptions) SetHeaders(param map[string]string) *ListKeyStoreCertificatesOptions {
	options.Headers = param
	return options
}

// ListQueueManagersOptions : The ListQueueManagers options.
type ListQueueManagersOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// Pagination offset.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of resources to return.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListQueueManagersOptions : Instantiate ListQueueManagersOptions
func (*MqcloudV1) NewListQueueManagersOptions(serviceInstanceGuid string) *ListQueueManagersOptions {
	return &ListQueueManagersOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *ListQueueManagersOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *ListQueueManagersOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListQueueManagersOptions) SetOffset(offset int64) *ListQueueManagersOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListQueueManagersOptions) SetLimit(limit int64) *ListQueueManagersOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListQueueManagersOptions) SetHeaders(param map[string]string) *ListQueueManagersOptions {
	options.Headers = param
	return options
}

// ListTrustStoreCertificatesOptions : The ListTrustStoreCertificates options.
type ListTrustStoreCertificatesOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListTrustStoreCertificatesOptions : Instantiate ListTrustStoreCertificatesOptions
func (*MqcloudV1) NewListTrustStoreCertificatesOptions(serviceInstanceGuid string, queueManagerID string) *ListTrustStoreCertificatesOptions {
	return &ListTrustStoreCertificatesOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *ListTrustStoreCertificatesOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *ListTrustStoreCertificatesOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *ListTrustStoreCertificatesOptions) SetQueueManagerID(queueManagerID string) *ListTrustStoreCertificatesOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListTrustStoreCertificatesOptions) SetHeaders(param map[string]string) *ListTrustStoreCertificatesOptions {
	options.Headers = param
	return options
}

// ListUsersOptions : The ListUsers options.
type ListUsersOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// Pagination offset.
	Offset *int64 `json:"offset,omitempty"`

	// The numbers of resources to return.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListUsersOptions : Instantiate ListUsersOptions
func (*MqcloudV1) NewListUsersOptions(serviceInstanceGuid string) *ListUsersOptions {
	return &ListUsersOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *ListUsersOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *ListUsersOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetOffset : Allow user to set Offset
func (_options *ListUsersOptions) SetOffset(offset int64) *ListUsersOptions {
	_options.Offset = core.Int64Ptr(offset)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListUsersOptions) SetLimit(limit int64) *ListUsersOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListUsersOptions) SetHeaders(param map[string]string) *ListUsersOptions {
	options.Headers = param
	return options
}

// ListVirtualPrivateEndpointGatewaysOptions : The ListVirtualPrivateEndpointGateways options.
type ListVirtualPrivateEndpointGatewaysOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The CRN of the trusted profile to assume for this request.
	TrustedProfile *string `json:"Trusted-Profile,omitempty"`

	// A server-provided token determining what resource to start the page on.
	Start *string `json:"start,omitempty"`

	// The numbers of resources to return.
	Limit *int64 `json:"limit,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewListVirtualPrivateEndpointGatewaysOptions : Instantiate ListVirtualPrivateEndpointGatewaysOptions
func (*MqcloudV1) NewListVirtualPrivateEndpointGatewaysOptions(serviceInstanceGuid string) *ListVirtualPrivateEndpointGatewaysOptions {
	return &ListVirtualPrivateEndpointGatewaysOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *ListVirtualPrivateEndpointGatewaysOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *ListVirtualPrivateEndpointGatewaysOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetTrustedProfile : Allow user to set TrustedProfile
func (_options *ListVirtualPrivateEndpointGatewaysOptions) SetTrustedProfile(trustedProfile string) *ListVirtualPrivateEndpointGatewaysOptions {
	_options.TrustedProfile = core.StringPtr(trustedProfile)
	return _options
}

// SetStart : Allow user to set Start
func (_options *ListVirtualPrivateEndpointGatewaysOptions) SetStart(start string) *ListVirtualPrivateEndpointGatewaysOptions {
	_options.Start = core.StringPtr(start)
	return _options
}

// SetLimit : Allow user to set Limit
func (_options *ListVirtualPrivateEndpointGatewaysOptions) SetLimit(limit int64) *ListVirtualPrivateEndpointGatewaysOptions {
	_options.Limit = core.Int64Ptr(limit)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *ListVirtualPrivateEndpointGatewaysOptions) SetHeaders(param map[string]string) *ListVirtualPrivateEndpointGatewaysOptions {
	options.Headers = param
	return options
}

// Next : Link to next page of results.
type Next struct {
	// The URL of the page the link goes to.
	Href *string `json:"href,omitempty"`
}

// UnmarshalNext unmarshals an instance of Next from the specified map of raw messages.
func UnmarshalNext(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Next)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Previous : Link to previous page of results.
type Previous struct {
	// The URL of the page the link goes to.
	Href *string `json:"href,omitempty"`
}

// UnmarshalPrevious unmarshals an instance of Previous from the specified map of raw messages.
func UnmarshalPrevious(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Previous)
	err = core.UnmarshalPrimitive(m, "href", &obj.Href)
	if err != nil {
		err = core.SDKErrorf(err, "", "href-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueueManagerDetails : The details of the queue manager.
type QueueManagerDetails struct {
	// The ID of the queue manager which was allocated on creation, and can be used for delete calls.
	ID *string `json:"id" validate:"required"`

	// A queue manager name conforming to MQ restrictions.
	Name *string `json:"name" validate:"required"`

	// A displayable name for the queue manager - limited only in length.
	DisplayName *string `json:"display_name" validate:"required"`

	// The locations in which the queue manager could be deployed.
	Location *string `json:"location" validate:"required"`

	// The queue manager sizes of deployment available.
	Size *string `json:"size" validate:"required"`

	// A reference uri to get deployment status of the queue manager.
	StatusURI *string `json:"status_uri" validate:"required"`

	// The MQ version of the queue manager.
	Version *string `json:"version" validate:"required"`

	// The url through which to access the web console for this queue manager.
	WebConsoleURL *string `json:"web_console_url" validate:"required"`

	// The url through which to access REST APIs for this queue manager.
	RestApiEndpointURL *string `json:"rest_api_endpoint_url" validate:"required"`

	// The url through which to access the Admin REST APIs for this queue manager.
	AdministratorApiEndpointURL *string `json:"administrator_api_endpoint_url" validate:"required"`

	// The uri through which the CDDT for this queue manager can be obtained.
	ConnectionInfoURI *string `json:"connection_info_uri" validate:"required"`

	// RFC3339 formatted UTC date for when the queue manager was created.
	DateCreated *strfmt.DateTime `json:"date_created" validate:"required"`

	// Describes whether an upgrade is available for this queue manager.
	UpgradeAvailable *bool `json:"upgrade_available" validate:"required"`

	// The uri through which the available versions to upgrade to can be found for this queue manager.
	AvailableUpgradeVersionsURI *string `json:"available_upgrade_versions_uri" validate:"required"`

	// The URL for this queue manager.
	Href *string `json:"href" validate:"required"`
}

// Constants associated with the QueueManagerDetails.Size property.
// The queue manager sizes of deployment available.
const (
	QueueManagerDetails_Size_Large  = "large"
	QueueManagerDetails_Size_Medium = "medium"
	QueueManagerDetails_Size_Small  = "small"
	QueueManagerDetails_Size_Xsmall = "xsmall"
)

// UnmarshalQueueManagerDetails unmarshals an instance of QueueManagerDetails from the specified map of raw messages.
func UnmarshalQueueManagerDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueueManagerDetails)
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
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		err = core.SDKErrorf(err, "", "display_name-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "location", &obj.Location)
	if err != nil {
		err = core.SDKErrorf(err, "", "location-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "size", &obj.Size)
	if err != nil {
		err = core.SDKErrorf(err, "", "size-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "status_uri", &obj.StatusURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "status_uri-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "web_console_url", &obj.WebConsoleURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "web_console_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "rest_api_endpoint_url", &obj.RestApiEndpointURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "rest_api_endpoint_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "administrator_api_endpoint_url", &obj.AdministratorApiEndpointURL)
	if err != nil {
		err = core.SDKErrorf(err, "", "administrator_api_endpoint_url-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "connection_info_uri", &obj.ConnectionInfoURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "connection_info_uri-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "date_created", &obj.DateCreated)
	if err != nil {
		err = core.SDKErrorf(err, "", "date_created-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "upgrade_available", &obj.UpgradeAvailable)
	if err != nil {
		err = core.SDKErrorf(err, "", "upgrade_available-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "available_upgrade_versions_uri", &obj.AvailableUpgradeVersionsURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "available_upgrade_versions_uri-error", common.GetComponentInfo())
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

// QueueManagerDetailsCollection : A list of queue manager summaries.
type QueueManagerDetailsCollection struct {
	// Pagination offset.
	Offset *int64 `json:"offset" validate:"required"`

	// Results per page, same for all collections.
	Limit *int64 `json:"limit" validate:"required"`

	// Link to first page of results.
	First *First `json:"first" validate:"required"`

	// Link to next page of results.
	Next *Next `json:"next,omitempty"`

	// Link to previous page of results.
	Previous *Previous `json:"previous,omitempty"`

	// List of queue managers.
	QueueManagers []QueueManagerDetails `json:"queue_managers" validate:"required"`
}

// UnmarshalQueueManagerDetailsCollection unmarshals an instance of QueueManagerDetailsCollection from the specified map of raw messages.
func UnmarshalQueueManagerDetailsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueueManagerDetailsCollection)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "queue_managers", &obj.QueueManagers, UnmarshalQueueManagerDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "queue_managers-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *QueueManagerDetailsCollection) GetNextOffset() (*int64, error) {
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

// QueueManagerStatus : Queue manager status.
type QueueManagerStatus struct {
	// The deploying and failed states are not queue manager states, they are states which can occur when the request to
	// deploy has been fired, or with that request has failed without producing a queue manager to have any state. The
	// other states map to the queue manager states. State "ending" is either quiesing or ending immediately. State "ended"
	// is either ended normally or endedimmediately. The others map one to one with queue manager states.
	Status *string `json:"status" validate:"required"`
}

// Constants associated with the QueueManagerStatus.Status property.
// The deploying and failed states are not queue manager states, they are states which can occur when the request to
// deploy has been fired, or with that request has failed without producing a queue manager to have any state. The other
// states map to the queue manager states. State "ending" is either quiesing or ending immediately. State "ended" is
// either ended normally or endedimmediately. The others map one to one with queue manager states.
const (
	QueueManagerStatus_Status_Deleting              = "deleting"
	QueueManagerStatus_Status_Deploying             = "deploying"
	QueueManagerStatus_Status_Failed                = "failed"
	QueueManagerStatus_Status_InitializationFailed  = "initialization_failed"
	QueueManagerStatus_Status_Initializing          = "initializing"
	QueueManagerStatus_Status_RestoreFailed         = "restore_failed"
	QueueManagerStatus_Status_RestoringConfig       = "restoring_config"
	QueueManagerStatus_Status_RestoringQueueManager = "restoring_queue_manager"
	QueueManagerStatus_Status_Resumable             = "resumable"
	QueueManagerStatus_Status_Running               = "running"
	QueueManagerStatus_Status_Starting              = "starting"
	QueueManagerStatus_Status_StatusNotAvailable    = "status_not_available"
	QueueManagerStatus_Status_Stopped               = "stopped"
	QueueManagerStatus_Status_Stopping              = "stopping"
	QueueManagerStatus_Status_Suspended             = "suspended"
	QueueManagerStatus_Status_UpdatingRevision      = "updating_revision"
	QueueManagerStatus_Status_UpgradingVersion      = "upgrading_version"
)

// UnmarshalQueueManagerStatus unmarshals an instance of QueueManagerStatus from the specified map of raw messages.
func UnmarshalQueueManagerStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueueManagerStatus)
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		err = core.SDKErrorf(err, "", "status-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueueManagerTaskStatus : A URI for status that can be queried periodically to get the status of the queue manager.
type QueueManagerTaskStatus struct {
	// Uri for the details of the queue manager.
	QueueManagerURI *string `json:"queue_manager_uri" validate:"required"`

	// Uri for the status of the queue manager.
	QueueManagerStatusURI *string `json:"queue_manager_status_uri" validate:"required"`

	// The queue manager id.
	QueueManagerID *string `json:"queue_manager_id" validate:"required"`
}

// UnmarshalQueueManagerTaskStatus unmarshals an instance of QueueManagerTaskStatus from the specified map of raw messages.
func UnmarshalQueueManagerTaskStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueueManagerTaskStatus)
	err = core.UnmarshalPrimitive(m, "queue_manager_uri", &obj.QueueManagerURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "queue_manager_uri-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "queue_manager_status_uri", &obj.QueueManagerStatusURI)
	if err != nil {
		err = core.SDKErrorf(err, "", "queue_manager_status_uri-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "queue_manager_id", &obj.QueueManagerID)
	if err != nil {
		err = core.SDKErrorf(err, "", "queue_manager_id-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueueManagerVersionUpgrade : An available upgrade for a queue manager.
type QueueManagerVersionUpgrade struct {
	// The target version of the queue manager upgrade.
	Version *string `json:"version" validate:"required"`

	// RFC3339 formatted UTC date for when the queue manager will automatically be updated.
	TargetDate *strfmt.DateTime `json:"target_date" validate:"required"`
}

// UnmarshalQueueManagerVersionUpgrade unmarshals an instance of QueueManagerVersionUpgrade from the specified map of raw messages.
func UnmarshalQueueManagerVersionUpgrade(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueueManagerVersionUpgrade)
	err = core.UnmarshalPrimitive(m, "version", &obj.Version)
	if err != nil {
		err = core.SDKErrorf(err, "", "version-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "target_date", &obj.TargetDate)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_date-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// QueueManagerVersionUpgrades : The list of available versions that this queue manger can upgrade to.
type QueueManagerVersionUpgrades struct {
	// Total count of versions available.
	TotalCount *int64 `json:"total_count" validate:"required"`

	// The list of available versions that this queue manger can upgrade to.
	Versions []QueueManagerVersionUpgrade `json:"versions" validate:"required"`
}

// UnmarshalQueueManagerVersionUpgrades unmarshals an instance of QueueManagerVersionUpgrades from the specified map of raw messages.
func UnmarshalQueueManagerVersionUpgrades(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(QueueManagerVersionUpgrades)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "versions", &obj.Versions, UnmarshalQueueManagerVersionUpgrade)
	if err != nil {
		err = core.SDKErrorf(err, "", "versions-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SetCertificateAmsChannelsOptions : The SetCertificateAmsChannels options.
type SetCertificateAmsChannelsOptions struct {
	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The id of the certificate.
	CertificateID *string `json:"certificate_id" validate:"required,ne="`

	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The list of AMS channels that are using this certificate.
	Channels []ChannelDetails `json:"channels" validate:"required"`

	// Strategy for how the supplied channels should be applied.
	UpdateStrategy *string `json:"update_strategy,omitempty"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// Constants associated with the SetCertificateAmsChannelsOptions.UpdateStrategy property.
// Strategy for how the supplied channels should be applied.
const (
	SetCertificateAmsChannelsOptions_UpdateStrategy_Append  = "append"
	SetCertificateAmsChannelsOptions_UpdateStrategy_Replace = "replace"
)

// NewSetCertificateAmsChannelsOptions : Instantiate SetCertificateAmsChannelsOptions
func (*MqcloudV1) NewSetCertificateAmsChannelsOptions(queueManagerID string, certificateID string, serviceInstanceGuid string, channels []ChannelDetails) *SetCertificateAmsChannelsOptions {
	return &SetCertificateAmsChannelsOptions{
		QueueManagerID:      core.StringPtr(queueManagerID),
		CertificateID:       core.StringPtr(certificateID),
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		Channels:            channels,
	}
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *SetCertificateAmsChannelsOptions) SetQueueManagerID(queueManagerID string) *SetCertificateAmsChannelsOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetCertificateID : Allow user to set CertificateID
func (_options *SetCertificateAmsChannelsOptions) SetCertificateID(certificateID string) *SetCertificateAmsChannelsOptions {
	_options.CertificateID = core.StringPtr(certificateID)
	return _options
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *SetCertificateAmsChannelsOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *SetCertificateAmsChannelsOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetChannels : Allow user to set Channels
func (_options *SetCertificateAmsChannelsOptions) SetChannels(channels []ChannelDetails) *SetCertificateAmsChannelsOptions {
	_options.Channels = channels
	return _options
}

// SetUpdateStrategy : Allow user to set UpdateStrategy
func (_options *SetCertificateAmsChannelsOptions) SetUpdateStrategy(updateStrategy string) *SetCertificateAmsChannelsOptions {
	_options.UpdateStrategy = core.StringPtr(updateStrategy)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetCertificateAmsChannelsOptions) SetHeaders(param map[string]string) *SetCertificateAmsChannelsOptions {
	options.Headers = param
	return options
}

// SetQueueManagerVersionOptions : The SetQueueManagerVersion options.
type SetQueueManagerVersionOptions struct {
	// The GUID that uniquely identifies the MQ on Cloud service instance.
	ServiceInstanceGuid *string `json:"service_instance_guid" validate:"required,ne="`

	// The id of the queue manager to retrieve its full details.
	QueueManagerID *string `json:"queue_manager_id" validate:"required,ne="`

	// The version upgrade to apply to the queue manager.
	Version *string `json:"version" validate:"required"`

	// Allows users to set headers on API requests.
	Headers map[string]string
}

// NewSetQueueManagerVersionOptions : Instantiate SetQueueManagerVersionOptions
func (*MqcloudV1) NewSetQueueManagerVersionOptions(serviceInstanceGuid string, queueManagerID string, version string) *SetQueueManagerVersionOptions {
	return &SetQueueManagerVersionOptions{
		ServiceInstanceGuid: core.StringPtr(serviceInstanceGuid),
		QueueManagerID:      core.StringPtr(queueManagerID),
		Version:             core.StringPtr(version),
	}
}

// SetServiceInstanceGuid : Allow user to set ServiceInstanceGuid
func (_options *SetQueueManagerVersionOptions) SetServiceInstanceGuid(serviceInstanceGuid string) *SetQueueManagerVersionOptions {
	_options.ServiceInstanceGuid = core.StringPtr(serviceInstanceGuid)
	return _options
}

// SetQueueManagerID : Allow user to set QueueManagerID
func (_options *SetQueueManagerVersionOptions) SetQueueManagerID(queueManagerID string) *SetQueueManagerVersionOptions {
	_options.QueueManagerID = core.StringPtr(queueManagerID)
	return _options
}

// SetVersion : Allow user to set Version
func (_options *SetQueueManagerVersionOptions) SetVersion(version string) *SetQueueManagerVersionOptions {
	_options.Version = core.StringPtr(version)
	return _options
}

// SetHeaders : Allow user to set Headers
func (options *SetQueueManagerVersionOptions) SetHeaders(param map[string]string) *SetQueueManagerVersionOptions {
	options.Headers = param
	return options
}

// TransmissionSecurity : An object that contains attributes that are related to security for message transmission.
type TransmissionSecurity struct {
	// Specifies the name of the CipherSpec for the channel to use.
	CipherSpecification *string `json:"cipherSpecification,omitempty"`
}

// UnmarshalTransmissionSecurity unmarshals an instance of TransmissionSecurity from the specified map of raw messages.
func UnmarshalTransmissionSecurity(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TransmissionSecurity)
	err = core.UnmarshalPrimitive(m, "cipherSpecification", &obj.CipherSpecification)
	if err != nil {
		err = core.SDKErrorf(err, "", "cipherSpecification-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// TrustStoreCertificateDetails : The details of a trust store certificate in a queue manager certificate trust store.
type TrustStoreCertificateDetails struct {
	// Id of the certificate.
	ID *string `json:"id" validate:"required"`

	// Certificate label in queue manager store.
	Label *string `json:"label" validate:"required"`

	// The type of certificate.
	CertificateType *string `json:"certificate_type" validate:"required"`

	// Fingerprint SHA256.
	FingerprintSha256 *string `json:"fingerprint_sha256" validate:"required"`

	// Subject's Distinguished Name.
	SubjectDn *string `json:"subject_dn" validate:"required"`

	// Subject's Common Name.
	SubjectCn *string `json:"subject_cn" validate:"required"`

	// Issuer's Distinguished Name.
	IssuerDn *string `json:"issuer_dn" validate:"required"`

	// Issuer's Common Name.
	IssuerCn *string `json:"issuer_cn" validate:"required"`

	// The Date the certificate was issued.
	Issued *strfmt.DateTime `json:"issued" validate:"required"`

	// Expiry date for the certificate.
	Expiry *strfmt.DateTime `json:"expiry" validate:"required"`

	// Indicates whether a certificate is trusted.
	Trusted *bool `json:"trusted" validate:"required"`

	// The URL for this trust store certificate.
	Href *string `json:"href" validate:"required"`
}

// Constants associated with the TrustStoreCertificateDetails.CertificateType property.
// The type of certificate.
const (
	TrustStoreCertificateDetails_CertificateType_TrustStore = "trust_store"
)

// UnmarshalTrustStoreCertificateDetails unmarshals an instance of TrustStoreCertificateDetails from the specified map of raw messages.
func UnmarshalTrustStoreCertificateDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustStoreCertificateDetails)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		err = core.SDKErrorf(err, "", "id-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "label", &obj.Label)
	if err != nil {
		err = core.SDKErrorf(err, "", "label-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "certificate_type", &obj.CertificateType)
	if err != nil {
		err = core.SDKErrorf(err, "", "certificate_type-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "fingerprint_sha256", &obj.FingerprintSha256)
	if err != nil {
		err = core.SDKErrorf(err, "", "fingerprint_sha256-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "subject_dn", &obj.SubjectDn)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject_dn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "subject_cn", &obj.SubjectCn)
	if err != nil {
		err = core.SDKErrorf(err, "", "subject_cn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer_dn", &obj.IssuerDn)
	if err != nil {
		err = core.SDKErrorf(err, "", "issuer_dn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "issuer_cn", &obj.IssuerCn)
	if err != nil {
		err = core.SDKErrorf(err, "", "issuer_cn-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "issued", &obj.Issued)
	if err != nil {
		err = core.SDKErrorf(err, "", "issued-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "expiry", &obj.Expiry)
	if err != nil {
		err = core.SDKErrorf(err, "", "expiry-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "trusted", &obj.Trusted)
	if err != nil {
		err = core.SDKErrorf(err, "", "trusted-error", common.GetComponentInfo())
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

// TrustStoreCertificateDetailsCollection : A list of certificates in a queue manager's certificate trust store.
type TrustStoreCertificateDetailsCollection struct {
	// The total count of trust store certificates.
	TotalCount *int64 `json:"total_count,omitempty"`

	// The list of trust store certificates.
	TrustStore []TrustStoreCertificateDetails `json:"trust_store,omitempty"`
}

// UnmarshalTrustStoreCertificateDetailsCollection unmarshals an instance of TrustStoreCertificateDetailsCollection from the specified map of raw messages.
func UnmarshalTrustStoreCertificateDetailsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(TrustStoreCertificateDetailsCollection)
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		err = core.SDKErrorf(err, "", "total_count-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "trust_store", &obj.TrustStore, UnmarshalTrustStoreCertificateDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "trust_store-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Usage : Usage details.
type Usage struct {
	// VPC entitlement.
	VpcEntitlement *float32 `json:"vpc_entitlement,omitempty"`

	// VPC usage.
	VpcUsage *float32 `json:"vpc_usage,omitempty"`
}

// UnmarshalUsage unmarshals an instance of Usage from the specified map of raw messages.
func UnmarshalUsage(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Usage)
	err = core.UnmarshalPrimitive(m, "vpc_entitlement", &obj.VpcEntitlement)
	if err != nil {
		err = core.SDKErrorf(err, "", "vpc_entitlement-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "vpc_usage", &obj.VpcUsage)
	if err != nil {
		err = core.SDKErrorf(err, "", "vpc_usage-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UserDetails : A summary of the user for use in a list of users.
type UserDetails struct {
	// The ID of the user which was allocated on creation, and can be used for delete calls.
	ID *string `json:"id" validate:"required"`

	// The shortname of the user that will be used as the IBM MQ administrator in interactions with a queue manager for
	// this service instance.
	Name *string `json:"name" validate:"required"`

	// The email of the user.
	Email *string `json:"email" validate:"required"`

	// The URL for the user details.
	Href *string `json:"href" validate:"required"`
}

// UnmarshalUserDetails unmarshals an instance of UserDetails from the specified map of raw messages.
func UnmarshalUserDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserDetails)
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
	err = core.UnmarshalPrimitive(m, "email", &obj.Email)
	if err != nil {
		err = core.SDKErrorf(err, "", "email-error", common.GetComponentInfo())
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

// UserDetailsCollection : A list of user summaries.
type UserDetailsCollection struct {
	// Pagination offset.
	Offset *int64 `json:"offset" validate:"required"`

	// Results per page, same for all collections.
	Limit *int64 `json:"limit" validate:"required"`

	// Link to first page of results.
	First *First `json:"first" validate:"required"`

	// Link to next page of results.
	Next *Next `json:"next,omitempty"`

	// Link to previous page of results.
	Previous *Previous `json:"previous,omitempty"`

	// List of users.
	Users []UserDetails `json:"users" validate:"required"`
}

// UnmarshalUserDetailsCollection unmarshals an instance of UserDetailsCollection from the specified map of raw messages.
func UnmarshalUserDetailsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UserDetailsCollection)
	err = core.UnmarshalPrimitive(m, "offset", &obj.Offset)
	if err != nil {
		err = core.SDKErrorf(err, "", "offset-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "previous", &obj.Previous, UnmarshalPrevious)
	if err != nil {
		err = core.SDKErrorf(err, "", "previous-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "users", &obj.Users, UnmarshalUserDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "users-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *UserDetailsCollection) GetNextOffset() (*int64, error) {
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

// VirtualPrivateEndpointGatewayDetails : The details of a specific Virtual Private Endpoint Gateway.
type VirtualPrivateEndpointGatewayDetails struct {
	// URL for the details of the virtual private endpoint gateway.
	Href *string `json:"href" validate:"required"`

	// The ID of the virtual private endpoint gateway which was allocated on creation.
	ID *string `json:"id" validate:"required"`

	// The name of the virtual private endpoint gateway, created by the user.
	Name *string `json:"name" validate:"required"`

	// The CRN of the virtual private endpoint gateway the user is trying to connect to.
	TargetCrn *string `json:"target_crn" validate:"required"`

	// The lifecycle state of this virtual privage endpoint.
	Status *string `json:"status" validate:"required"`
}

// UnmarshalVirtualPrivateEndpointGatewayDetails unmarshals an instance of VirtualPrivateEndpointGatewayDetails from the specified map of raw messages.
func UnmarshalVirtualPrivateEndpointGatewayDetails(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VirtualPrivateEndpointGatewayDetails)
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
	err = core.UnmarshalPrimitive(m, "target_crn", &obj.TargetCrn)
	if err != nil {
		err = core.SDKErrorf(err, "", "target_crn-error", common.GetComponentInfo())
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

// VirtualPrivateEndpointGatewayDetailsCollection : A list of virtual private endpoint gateway summaries.
type VirtualPrivateEndpointGatewayDetailsCollection struct {
	// Results per page, same for all collections.
	Limit *int64 `json:"limit" validate:"required"`

	// Link to first page of results.
	First *First `json:"first" validate:"required"`

	// Link to next page of results.
	Next *Next `json:"next,omitempty"`

	// List of virtual private endpoint gateways.
	VirtualPrivateEndpointGateways []VirtualPrivateEndpointGatewayDetails `json:"virtual_private_endpoint_gateways" validate:"required"`
}

// UnmarshalVirtualPrivateEndpointGatewayDetailsCollection unmarshals an instance of VirtualPrivateEndpointGatewayDetailsCollection from the specified map of raw messages.
func UnmarshalVirtualPrivateEndpointGatewayDetailsCollection(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(VirtualPrivateEndpointGatewayDetailsCollection)
	err = core.UnmarshalPrimitive(m, "limit", &obj.Limit)
	if err != nil {
		err = core.SDKErrorf(err, "", "limit-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "first", &obj.First, UnmarshalFirst)
	if err != nil {
		err = core.SDKErrorf(err, "", "first-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "next", &obj.Next, UnmarshalNext)
	if err != nil {
		err = core.SDKErrorf(err, "", "next-error", common.GetComponentInfo())
		return
	}
	err = core.UnmarshalModel(m, "virtual_private_endpoint_gateways", &obj.VirtualPrivateEndpointGateways, UnmarshalVirtualPrivateEndpointGatewayDetails)
	if err != nil {
		err = core.SDKErrorf(err, "", "virtual_private_endpoint_gateways-error", common.GetComponentInfo())
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Retrieve the value to be passed to a request to access the next page of results
func (resp *VirtualPrivateEndpointGatewayDetailsCollection) GetNextStart() (*string, error) {
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

// QueueManagersPager can be used to simplify the use of the "ListQueueManagers" method.
type QueueManagersPager struct {
	hasNext     bool
	options     *ListQueueManagersOptions
	client      *MqcloudV1
	pageContext struct {
		next *int64
	}
}

// NewQueueManagersPager returns a new QueueManagersPager instance.
func (mqcloud *MqcloudV1) NewQueueManagersPager(options *ListQueueManagersOptions) (pager *QueueManagersPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListQueueManagersOptions = *options
	pager = &QueueManagersPager{
		hasNext: true,
		options: &optionsCopy,
		client:  mqcloud,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *QueueManagersPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *QueueManagersPager) GetNextWithContext(ctx context.Context) (page []QueueManagerDetails, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListQueueManagersWithContext(ctx, pager.options)
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
	page = result.QueueManagers

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *QueueManagersPager) GetAllWithContext(ctx context.Context) (allItems []QueueManagerDetails, err error) {
	for pager.HasNext() {
		var nextPage []QueueManagerDetails
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
func (pager *QueueManagersPager) GetNext() (page []QueueManagerDetails, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *QueueManagersPager) GetAll() (allItems []QueueManagerDetails, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// UsersPager can be used to simplify the use of the "ListUsers" method.
type UsersPager struct {
	hasNext     bool
	options     *ListUsersOptions
	client      *MqcloudV1
	pageContext struct {
		next *int64
	}
}

// NewUsersPager returns a new UsersPager instance.
func (mqcloud *MqcloudV1) NewUsersPager(options *ListUsersOptions) (pager *UsersPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListUsersOptions = *options
	pager = &UsersPager{
		hasNext: true,
		options: &optionsCopy,
		client:  mqcloud,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *UsersPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *UsersPager) GetNextWithContext(ctx context.Context) (page []UserDetails, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListUsersWithContext(ctx, pager.options)
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
	page = result.Users

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *UsersPager) GetAllWithContext(ctx context.Context) (allItems []UserDetails, err error) {
	for pager.HasNext() {
		var nextPage []UserDetails
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
func (pager *UsersPager) GetNext() (page []UserDetails, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *UsersPager) GetAll() (allItems []UserDetails, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// ApplicationsPager can be used to simplify the use of the "ListApplications" method.
type ApplicationsPager struct {
	hasNext     bool
	options     *ListApplicationsOptions
	client      *MqcloudV1
	pageContext struct {
		next *int64
	}
}

// NewApplicationsPager returns a new ApplicationsPager instance.
func (mqcloud *MqcloudV1) NewApplicationsPager(options *ListApplicationsOptions) (pager *ApplicationsPager, err error) {
	if options.Offset != nil && *options.Offset != 0 {
		err = core.SDKErrorf(nil, "the 'options.Offset' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListApplicationsOptions = *options
	pager = &ApplicationsPager{
		hasNext: true,
		options: &optionsCopy,
		client:  mqcloud,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *ApplicationsPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *ApplicationsPager) GetNextWithContext(ctx context.Context) (page []ApplicationDetails, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Offset = pager.pageContext.next

	result, _, err := pager.client.ListApplicationsWithContext(ctx, pager.options)
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
	page = result.Applications

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *ApplicationsPager) GetAllWithContext(ctx context.Context) (allItems []ApplicationDetails, err error) {
	for pager.HasNext() {
		var nextPage []ApplicationDetails
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
func (pager *ApplicationsPager) GetNext() (page []ApplicationDetails, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *ApplicationsPager) GetAll() (allItems []ApplicationDetails, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// VirtualPrivateEndpointGatewaysPager can be used to simplify the use of the "ListVirtualPrivateEndpointGateways" method.
type VirtualPrivateEndpointGatewaysPager struct {
	hasNext     bool
	options     *ListVirtualPrivateEndpointGatewaysOptions
	client      *MqcloudV1
	pageContext struct {
		next *string
	}
}

// NewVirtualPrivateEndpointGatewaysPager returns a new VirtualPrivateEndpointGatewaysPager instance.
func (mqcloud *MqcloudV1) NewVirtualPrivateEndpointGatewaysPager(options *ListVirtualPrivateEndpointGatewaysOptions) (pager *VirtualPrivateEndpointGatewaysPager, err error) {
	if options.Start != nil && *options.Start != "" {
		err = core.SDKErrorf(nil, "the 'options.Start' field should not be set", "no-query-setting", common.GetComponentInfo())
		return
	}

	var optionsCopy ListVirtualPrivateEndpointGatewaysOptions = *options
	pager = &VirtualPrivateEndpointGatewaysPager{
		hasNext: true,
		options: &optionsCopy,
		client:  mqcloud,
	}
	return
}

// HasNext returns true if there are potentially more results to be retrieved.
func (pager *VirtualPrivateEndpointGatewaysPager) HasNext() bool {
	return pager.hasNext
}

// GetNextWithContext returns the next page of results using the specified Context.
func (pager *VirtualPrivateEndpointGatewaysPager) GetNextWithContext(ctx context.Context) (page []VirtualPrivateEndpointGatewayDetails, err error) {
	if !pager.HasNext() {
		return nil, fmt.Errorf("no more results available")
	}

	pager.options.Start = pager.pageContext.next

	result, _, err := pager.client.ListVirtualPrivateEndpointGatewaysWithContext(ctx, pager.options)
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
	page = result.VirtualPrivateEndpointGateways

	return
}

// GetAllWithContext returns all results by invoking GetNextWithContext() repeatedly
// until all pages of results have been retrieved.
func (pager *VirtualPrivateEndpointGatewaysPager) GetAllWithContext(ctx context.Context) (allItems []VirtualPrivateEndpointGatewayDetails, err error) {
	for pager.HasNext() {
		var nextPage []VirtualPrivateEndpointGatewayDetails
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
func (pager *VirtualPrivateEndpointGatewaysPager) GetNext() (page []VirtualPrivateEndpointGatewayDetails, err error) {
	page, err = pager.GetNextWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}

// GetAll invokes GetAllWithContext() using context.Background() as the Context parameter.
func (pager *VirtualPrivateEndpointGatewaysPager) GetAll() (allItems []VirtualPrivateEndpointGatewayDetails, err error) {
	allItems, err = pager.GetAllWithContext(context.Background())
	err = core.RepurposeSDKProblem(err, "")
	return
}
