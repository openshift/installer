/**
 * (C) Copyright IBM Corp. 2020.
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
 * IBM OpenAPI SDK Code Generator Version: 3.20.0-debb9f29-20201203-202043
 */
 

// Package globalloadbalancerpoolsv0 : Operations and models for the GlobalLoadBalancerPoolsV0 service
package globalloadbalancerpoolsv0

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"net/http"
	"reflect"
	"time"
)

// GlobalLoadBalancerPoolsV0 : GLB Pools
//
// Version: 0.0.1
type GlobalLoadBalancerPoolsV0 struct {
	Service *core.BaseService

	// Full CRN of the service instance.
	Crn *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "global_load_balancer_pools"

// GlobalLoadBalancerPoolsV0Options : Service options
type GlobalLoadBalancerPoolsV0Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full CRN of the service instance.
	Crn *string `validate:"required"`
}

// NewGlobalLoadBalancerPoolsV0UsingExternalConfig : constructs an instance of GlobalLoadBalancerPoolsV0 with passed in options and external configuration.
func NewGlobalLoadBalancerPoolsV0UsingExternalConfig(options *GlobalLoadBalancerPoolsV0Options) (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	globalLoadBalancerPools, err = NewGlobalLoadBalancerPoolsV0(options)
	if err != nil {
		return
	}

	err = globalLoadBalancerPools.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = globalLoadBalancerPools.Service.SetServiceURL(options.URL)
	}
	return
}

// NewGlobalLoadBalancerPoolsV0 : constructs an instance of GlobalLoadBalancerPoolsV0 with passed in options.
func NewGlobalLoadBalancerPoolsV0(options *GlobalLoadBalancerPoolsV0Options) (service *GlobalLoadBalancerPoolsV0, err error) {
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

	service = &GlobalLoadBalancerPoolsV0{
		Service: baseService,
		Crn: options.Crn,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "globalLoadBalancerPools" suitable for processing requests.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) Clone() *GlobalLoadBalancerPoolsV0 {
	if core.IsNil(globalLoadBalancerPools) {
		return nil
	}
	clone := *globalLoadBalancerPools
	clone.Service = globalLoadBalancerPools.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) SetServiceURL(url string) error {
	return globalLoadBalancerPools.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) GetServiceURL() string {
	return globalLoadBalancerPools.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) SetDefaultHeaders(headers http.Header) {
	globalLoadBalancerPools.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) SetEnableGzipCompression(enableGzip bool) {
	globalLoadBalancerPools.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) GetEnableGzipCompression() bool {
	return globalLoadBalancerPools.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	globalLoadBalancerPools.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) DisableRetries() {
	globalLoadBalancerPools.Service.DisableRetries()
}

// ListAllLoadBalancerPools : List all pools
// List all configured load balancer pools.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) ListAllLoadBalancerPools(listAllLoadBalancerPoolsOptions *ListAllLoadBalancerPoolsOptions) (result *ListLoadBalancerPoolsResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerPools.ListAllLoadBalancerPoolsWithContext(context.Background(), listAllLoadBalancerPoolsOptions)
}

// ListAllLoadBalancerPoolsWithContext is an alternate form of the ListAllLoadBalancerPools method which supports a Context parameter
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) ListAllLoadBalancerPoolsWithContext(ctx context.Context, listAllLoadBalancerPoolsOptions *ListAllLoadBalancerPoolsOptions) (result *ListLoadBalancerPoolsResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllLoadBalancerPoolsOptions, "listAllLoadBalancerPoolsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerPools.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerPools.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerPools.Service.Options.URL, `/v1/{crn}/load_balancers/pools`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllLoadBalancerPoolsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_pools", "V0", "ListAllLoadBalancerPools")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancerPools.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLoadBalancerPoolsResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancerPool : Create pool
// Create a new load balancer pool.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) CreateLoadBalancerPool(createLoadBalancerPoolOptions *CreateLoadBalancerPoolOptions) (result *LoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerPools.CreateLoadBalancerPoolWithContext(context.Background(), createLoadBalancerPoolOptions)
}

// CreateLoadBalancerPoolWithContext is an alternate form of the CreateLoadBalancerPool method which supports a Context parameter
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) CreateLoadBalancerPoolWithContext(ctx context.Context, createLoadBalancerPoolOptions *CreateLoadBalancerPoolOptions) (result *LoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLoadBalancerPoolOptions, "createLoadBalancerPoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerPools.Crn,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerPools.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerPools.Service.Options.URL, `/v1/{crn}/load_balancers/pools`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createLoadBalancerPoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_pools", "V0", "CreateLoadBalancerPool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLoadBalancerPoolOptions.Name != nil {
		body["name"] = createLoadBalancerPoolOptions.Name
	}
	if createLoadBalancerPoolOptions.CheckRegions != nil {
		body["check_regions"] = createLoadBalancerPoolOptions.CheckRegions
	}
	if createLoadBalancerPoolOptions.Origins != nil {
		body["origins"] = createLoadBalancerPoolOptions.Origins
	}
	if createLoadBalancerPoolOptions.Description != nil {
		body["description"] = createLoadBalancerPoolOptions.Description
	}
	if createLoadBalancerPoolOptions.MinimumOrigins != nil {
		body["minimum_origins"] = createLoadBalancerPoolOptions.MinimumOrigins
	}
	if createLoadBalancerPoolOptions.Enabled != nil {
		body["enabled"] = createLoadBalancerPoolOptions.Enabled
	}
	if createLoadBalancerPoolOptions.Monitor != nil {
		body["monitor"] = createLoadBalancerPoolOptions.Monitor
	}
	if createLoadBalancerPoolOptions.NotificationEmail != nil {
		body["notification_email"] = createLoadBalancerPoolOptions.NotificationEmail
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
	response, err = globalLoadBalancerPools.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancerPoolResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetLoadBalancerPool : Get pool
// Get a single configured load balancer pool.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) GetLoadBalancerPool(getLoadBalancerPoolOptions *GetLoadBalancerPoolOptions) (result *LoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerPools.GetLoadBalancerPoolWithContext(context.Background(), getLoadBalancerPoolOptions)
}

// GetLoadBalancerPoolWithContext is an alternate form of the GetLoadBalancerPool method which supports a Context parameter
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) GetLoadBalancerPoolWithContext(ctx context.Context, getLoadBalancerPoolOptions *GetLoadBalancerPoolOptions) (result *LoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLoadBalancerPoolOptions, "getLoadBalancerPoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLoadBalancerPoolOptions, "getLoadBalancerPoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerPools.Crn,
		"pool_identifier": *getLoadBalancerPoolOptions.PoolIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerPools.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerPools.Service.Options.URL, `/v1/{crn}/load_balancers/pools/{pool_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLoadBalancerPoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_pools", "V0", "GetLoadBalancerPool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancerPools.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancerPoolResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteLoadBalancerPool : Delete pool
// Delete a specific configured load balancer pool.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) DeleteLoadBalancerPool(deleteLoadBalancerPoolOptions *DeleteLoadBalancerPoolOptions) (result *DeleteLoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerPools.DeleteLoadBalancerPoolWithContext(context.Background(), deleteLoadBalancerPoolOptions)
}

// DeleteLoadBalancerPoolWithContext is an alternate form of the DeleteLoadBalancerPool method which supports a Context parameter
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) DeleteLoadBalancerPoolWithContext(ctx context.Context, deleteLoadBalancerPoolOptions *DeleteLoadBalancerPoolOptions) (result *DeleteLoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLoadBalancerPoolOptions, "deleteLoadBalancerPoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLoadBalancerPoolOptions, "deleteLoadBalancerPoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerPools.Crn,
		"pool_identifier": *deleteLoadBalancerPoolOptions.PoolIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerPools.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerPools.Service.Options.URL, `/v1/{crn}/load_balancers/pools/{pool_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteLoadBalancerPoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_pools", "V0", "DeleteLoadBalancerPool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancerPools.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteLoadBalancerPoolResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// EditLoadBalancerPool : Edit pool
// Edit a specific configured load balancer pool.
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) EditLoadBalancerPool(editLoadBalancerPoolOptions *EditLoadBalancerPoolOptions) (result *LoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerPools.EditLoadBalancerPoolWithContext(context.Background(), editLoadBalancerPoolOptions)
}

// EditLoadBalancerPoolWithContext is an alternate form of the EditLoadBalancerPool method which supports a Context parameter
func (globalLoadBalancerPools *GlobalLoadBalancerPoolsV0) EditLoadBalancerPoolWithContext(ctx context.Context, editLoadBalancerPoolOptions *EditLoadBalancerPoolOptions) (result *LoadBalancerPoolResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(editLoadBalancerPoolOptions, "editLoadBalancerPoolOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(editLoadBalancerPoolOptions, "editLoadBalancerPoolOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerPools.Crn,
		"pool_identifier": *editLoadBalancerPoolOptions.PoolIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerPools.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerPools.Service.Options.URL, `/v1/{crn}/load_balancers/pools/{pool_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range editLoadBalancerPoolOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_pools", "V0", "EditLoadBalancerPool")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if editLoadBalancerPoolOptions.Name != nil {
		body["name"] = editLoadBalancerPoolOptions.Name
	}
	if editLoadBalancerPoolOptions.CheckRegions != nil {
		body["check_regions"] = editLoadBalancerPoolOptions.CheckRegions
	}
	if editLoadBalancerPoolOptions.Origins != nil {
		body["origins"] = editLoadBalancerPoolOptions.Origins
	}
	if editLoadBalancerPoolOptions.Description != nil {
		body["description"] = editLoadBalancerPoolOptions.Description
	}
	if editLoadBalancerPoolOptions.MinimumOrigins != nil {
		body["minimum_origins"] = editLoadBalancerPoolOptions.MinimumOrigins
	}
	if editLoadBalancerPoolOptions.Enabled != nil {
		body["enabled"] = editLoadBalancerPoolOptions.Enabled
	}
	if editLoadBalancerPoolOptions.Monitor != nil {
		body["monitor"] = editLoadBalancerPoolOptions.Monitor
	}
	if editLoadBalancerPoolOptions.NotificationEmail != nil {
		body["notification_email"] = editLoadBalancerPoolOptions.NotificationEmail
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
	response, err = globalLoadBalancerPools.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancerPoolResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancerPoolOptions : The CreateLoadBalancerPool options.
type CreateLoadBalancerPoolOptions struct {
	// name.
	Name *string `json:"name,omitempty"`

	// regions check.
	CheckRegions []string `json:"check_regions,omitempty"`

	// origins.
	Origins []LoadBalancerPoolReqOriginsItem `json:"origins,omitempty"`

	// desc.
	Description *string `json:"description,omitempty"`

	// The minimum number of origins that must be healthy for this pool to serve traffic.
	MinimumOrigins *int64 `json:"minimum_origins,omitempty"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// monitor.
	Monitor *string `json:"monitor,omitempty"`

	// notification email.
	NotificationEmail *string `json:"notification_email,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateLoadBalancerPoolOptions : Instantiate CreateLoadBalancerPoolOptions
func (*GlobalLoadBalancerPoolsV0) NewCreateLoadBalancerPoolOptions() *CreateLoadBalancerPoolOptions {
	return &CreateLoadBalancerPoolOptions{}
}

// SetName : Allow user to set Name
func (options *CreateLoadBalancerPoolOptions) SetName(name string) *CreateLoadBalancerPoolOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetCheckRegions : Allow user to set CheckRegions
func (options *CreateLoadBalancerPoolOptions) SetCheckRegions(checkRegions []string) *CreateLoadBalancerPoolOptions {
	options.CheckRegions = checkRegions
	return options
}

// SetOrigins : Allow user to set Origins
func (options *CreateLoadBalancerPoolOptions) SetOrigins(origins []LoadBalancerPoolReqOriginsItem) *CreateLoadBalancerPoolOptions {
	options.Origins = origins
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateLoadBalancerPoolOptions) SetDescription(description string) *CreateLoadBalancerPoolOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetMinimumOrigins : Allow user to set MinimumOrigins
func (options *CreateLoadBalancerPoolOptions) SetMinimumOrigins(minimumOrigins int64) *CreateLoadBalancerPoolOptions {
	options.MinimumOrigins = core.Int64Ptr(minimumOrigins)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *CreateLoadBalancerPoolOptions) SetEnabled(enabled bool) *CreateLoadBalancerPoolOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetMonitor : Allow user to set Monitor
func (options *CreateLoadBalancerPoolOptions) SetMonitor(monitor string) *CreateLoadBalancerPoolOptions {
	options.Monitor = core.StringPtr(monitor)
	return options
}

// SetNotificationEmail : Allow user to set NotificationEmail
func (options *CreateLoadBalancerPoolOptions) SetNotificationEmail(notificationEmail string) *CreateLoadBalancerPoolOptions {
	options.NotificationEmail = core.StringPtr(notificationEmail)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLoadBalancerPoolOptions) SetHeaders(param map[string]string) *CreateLoadBalancerPoolOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancerPoolOptions : The DeleteLoadBalancerPool options.
type DeleteLoadBalancerPoolOptions struct {
	// pool identifier.
	PoolIdentifier *string `json:"pool_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLoadBalancerPoolOptions : Instantiate DeleteLoadBalancerPoolOptions
func (*GlobalLoadBalancerPoolsV0) NewDeleteLoadBalancerPoolOptions(poolIdentifier string) *DeleteLoadBalancerPoolOptions {
	return &DeleteLoadBalancerPoolOptions{
		PoolIdentifier: core.StringPtr(poolIdentifier),
	}
}

// SetPoolIdentifier : Allow user to set PoolIdentifier
func (options *DeleteLoadBalancerPoolOptions) SetPoolIdentifier(poolIdentifier string) *DeleteLoadBalancerPoolOptions {
	options.PoolIdentifier = core.StringPtr(poolIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLoadBalancerPoolOptions) SetHeaders(param map[string]string) *DeleteLoadBalancerPoolOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancerPoolRespResult : result.
type DeleteLoadBalancerPoolRespResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteLoadBalancerPoolRespResult unmarshals an instance of DeleteLoadBalancerPoolRespResult from the specified map of raw messages.
func UnmarshalDeleteLoadBalancerPoolRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLoadBalancerPoolRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EditLoadBalancerPoolOptions : The EditLoadBalancerPool options.
type EditLoadBalancerPoolOptions struct {
	// pool identifier.
	PoolIdentifier *string `json:"pool_identifier" validate:"required,ne="`

	// name.
	Name *string `json:"name,omitempty"`

	// regions check.
	CheckRegions []string `json:"check_regions,omitempty"`

	// origins.
	Origins []LoadBalancerPoolReqOriginsItem `json:"origins,omitempty"`

	// desc.
	Description *string `json:"description,omitempty"`

	// The minimum number of origins that must be healthy for this pool to serve traffic.
	MinimumOrigins *int64 `json:"minimum_origins,omitempty"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// monitor.
	Monitor *string `json:"monitor,omitempty"`

	// notification email.
	NotificationEmail *string `json:"notification_email,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewEditLoadBalancerPoolOptions : Instantiate EditLoadBalancerPoolOptions
func (*GlobalLoadBalancerPoolsV0) NewEditLoadBalancerPoolOptions(poolIdentifier string) *EditLoadBalancerPoolOptions {
	return &EditLoadBalancerPoolOptions{
		PoolIdentifier: core.StringPtr(poolIdentifier),
	}
}

// SetPoolIdentifier : Allow user to set PoolIdentifier
func (options *EditLoadBalancerPoolOptions) SetPoolIdentifier(poolIdentifier string) *EditLoadBalancerPoolOptions {
	options.PoolIdentifier = core.StringPtr(poolIdentifier)
	return options
}

// SetName : Allow user to set Name
func (options *EditLoadBalancerPoolOptions) SetName(name string) *EditLoadBalancerPoolOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetCheckRegions : Allow user to set CheckRegions
func (options *EditLoadBalancerPoolOptions) SetCheckRegions(checkRegions []string) *EditLoadBalancerPoolOptions {
	options.CheckRegions = checkRegions
	return options
}

// SetOrigins : Allow user to set Origins
func (options *EditLoadBalancerPoolOptions) SetOrigins(origins []LoadBalancerPoolReqOriginsItem) *EditLoadBalancerPoolOptions {
	options.Origins = origins
	return options
}

// SetDescription : Allow user to set Description
func (options *EditLoadBalancerPoolOptions) SetDescription(description string) *EditLoadBalancerPoolOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetMinimumOrigins : Allow user to set MinimumOrigins
func (options *EditLoadBalancerPoolOptions) SetMinimumOrigins(minimumOrigins int64) *EditLoadBalancerPoolOptions {
	options.MinimumOrigins = core.Int64Ptr(minimumOrigins)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *EditLoadBalancerPoolOptions) SetEnabled(enabled bool) *EditLoadBalancerPoolOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetMonitor : Allow user to set Monitor
func (options *EditLoadBalancerPoolOptions) SetMonitor(monitor string) *EditLoadBalancerPoolOptions {
	options.Monitor = core.StringPtr(monitor)
	return options
}

// SetNotificationEmail : Allow user to set NotificationEmail
func (options *EditLoadBalancerPoolOptions) SetNotificationEmail(notificationEmail string) *EditLoadBalancerPoolOptions {
	options.NotificationEmail = core.StringPtr(notificationEmail)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *EditLoadBalancerPoolOptions) SetHeaders(param map[string]string) *EditLoadBalancerPoolOptions {
	options.Headers = param
	return options
}

// GetLoadBalancerPoolOptions : The GetLoadBalancerPool options.
type GetLoadBalancerPoolOptions struct {
	// pool identifier.
	PoolIdentifier *string `json:"pool_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLoadBalancerPoolOptions : Instantiate GetLoadBalancerPoolOptions
func (*GlobalLoadBalancerPoolsV0) NewGetLoadBalancerPoolOptions(poolIdentifier string) *GetLoadBalancerPoolOptions {
	return &GetLoadBalancerPoolOptions{
		PoolIdentifier: core.StringPtr(poolIdentifier),
	}
}

// SetPoolIdentifier : Allow user to set PoolIdentifier
func (options *GetLoadBalancerPoolOptions) SetPoolIdentifier(poolIdentifier string) *GetLoadBalancerPoolOptions {
	options.PoolIdentifier = core.StringPtr(poolIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetLoadBalancerPoolOptions) SetHeaders(param map[string]string) *GetLoadBalancerPoolOptions {
	options.Headers = param
	return options
}

// ListAllLoadBalancerPoolsOptions : The ListAllLoadBalancerPools options.
type ListAllLoadBalancerPoolsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllLoadBalancerPoolsOptions : Instantiate ListAllLoadBalancerPoolsOptions
func (*GlobalLoadBalancerPoolsV0) NewListAllLoadBalancerPoolsOptions() *ListAllLoadBalancerPoolsOptions {
	return &ListAllLoadBalancerPoolsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListAllLoadBalancerPoolsOptions) SetHeaders(param map[string]string) *ListAllLoadBalancerPoolsOptions {
	options.Headers = param
	return options
}

// LoadBalancerPoolPackOriginsItem : LoadBalancerPoolPackOriginsItem struct
type LoadBalancerPoolPackOriginsItem struct {
	// name.
	Name *string `json:"name,omitempty"`

	// address.
	Address *string `json:"address,omitempty"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// healthy.
	Healthy *bool `json:"healthy,omitempty"`

	// weight.
	Weight *float64 `json:"weight,omitempty"`

	// Pool origin disabled date.
	DisabledAt *string `json:"disabled_at,omitempty"`

	// Reason for failure.
	FailureReason *string `json:"failure_reason,omitempty"`
}


// UnmarshalLoadBalancerPoolPackOriginsItem unmarshals an instance of LoadBalancerPoolPackOriginsItem from the specified map of raw messages.
func UnmarshalLoadBalancerPoolPackOriginsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancerPoolPackOriginsItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "address", &obj.Address)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "healthy", &obj.Healthy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "disabled_at", &obj.DisabledAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "failure_reason", &obj.FailureReason)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LoadBalancerPoolReqOriginsItem : items.
type LoadBalancerPoolReqOriginsItem struct {
	// name.
	Name *string `json:"name,omitempty"`

	// address.
	Address *string `json:"address,omitempty"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// weight.
	Weight *float64 `json:"weight,omitempty"`
}


// UnmarshalLoadBalancerPoolReqOriginsItem unmarshals an instance of LoadBalancerPoolReqOriginsItem from the specified map of raw messages.
func UnmarshalLoadBalancerPoolReqOriginsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancerPoolReqOriginsItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "address", &obj.Address)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "weight", &obj.Weight)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DeleteLoadBalancerPoolResp : load balancer pool delete response.
type DeleteLoadBalancerPoolResp struct {
	// succcess response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result *DeleteLoadBalancerPoolRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteLoadBalancerPoolResp unmarshals an instance of DeleteLoadBalancerPoolResp from the specified map of raw messages.
func UnmarshalDeleteLoadBalancerPoolResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLoadBalancerPoolResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteLoadBalancerPoolRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListLoadBalancerPoolsResp : list load balancer pools response.
type ListLoadBalancerPoolsResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result []LoadBalancerPoolPack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListLoadBalancerPoolsResp unmarshals an instance of ListLoadBalancerPoolsResp from the specified map of raw messages.
func UnmarshalListLoadBalancerPoolsResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLoadBalancerPoolsResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLoadBalancerPoolPack)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LoadBalancerPoolPack : load balancer pool pack.
type LoadBalancerPoolPack struct {
	// identifier.
	ID *string `json:"id,omitempty"`

	// created date.
	CreatedOn *string `json:"created_on,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// desc.
	Description *string `json:"description,omitempty"`

	// name.
	Name *string `json:"name" validate:"required"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// healthy.
	Healthy *bool `json:"healthy,omitempty"`

	// monitor.
	Monitor *string `json:"monitor,omitempty"`

	// Minimum origin count.
	MinimumOrigins *int64 `json:"minimum_origins,omitempty"`

	// regions check.
	CheckRegions []string `json:"check_regions,omitempty"`

	// original.
	Origins []LoadBalancerPoolPackOriginsItem `json:"origins" validate:"required"`

	// notification email.
	NotificationEmail *string `json:"notification_email,omitempty"`
}


// UnmarshalLoadBalancerPoolPack unmarshals an instance of LoadBalancerPoolPack from the specified map of raw messages.
func UnmarshalLoadBalancerPoolPack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancerPoolPack)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_on", &obj.CreatedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "modified_on", &obj.ModifiedOn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
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
	err = core.UnmarshalPrimitive(m, "healthy", &obj.Healthy)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "monitor", &obj.Monitor)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "minimum_origins", &obj.MinimumOrigins)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "check_regions", &obj.CheckRegions)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "origins", &obj.Origins, UnmarshalLoadBalancerPoolPackOriginsItem)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "notification_email", &obj.NotificationEmail)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LoadBalancerPoolResp : get load balancer pool response.
type LoadBalancerPoolResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// load balancer pool pack.
	Result *LoadBalancerPoolPack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalLoadBalancerPoolResp unmarshals an instance of LoadBalancerPoolResp from the specified map of raw messages.
func UnmarshalLoadBalancerPoolResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancerPoolResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLoadBalancerPoolPack)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "result_info", &obj.ResultInfo, UnmarshalResultInfo)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ResultInfo : result information.
type ResultInfo struct {
	// page number.
	Page *int64 `json:"page" validate:"required"`

	// per page count.
	PerPage *int64 `json:"per_page" validate:"required"`

	// count.
	Count *int64 `json:"count" validate:"required"`

	// total count.
	TotalCount *int64 `json:"total_count" validate:"required"`
}


// UnmarshalResultInfo unmarshals an instance of ResultInfo from the specified map of raw messages.
func UnmarshalResultInfo(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ResultInfo)
	err = core.UnmarshalPrimitive(m, "page", &obj.Page)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "per_page", &obj.PerPage)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "count", &obj.Count)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "total_count", &obj.TotalCount)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
