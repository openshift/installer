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
 

// Package globalloadbalancerv1 : Operations and models for the GlobalLoadBalancerV1 service
package globalloadbalancerv1

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

// GlobalLoadBalancerV1 : Global Load Balancer
//
// Version: 1.0.1
type GlobalLoadBalancerV1 struct {
	Service *core.BaseService

	// Full CRN of the service instance.
	Crn *string

	// zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "global_load_balancer"

// GlobalLoadBalancerV1Options : Service options
type GlobalLoadBalancerV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full CRN of the service instance.
	Crn *string `validate:"required"`

	// zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewGlobalLoadBalancerV1UsingExternalConfig : constructs an instance of GlobalLoadBalancerV1 with passed in options and external configuration.
func NewGlobalLoadBalancerV1UsingExternalConfig(options *GlobalLoadBalancerV1Options) (globalLoadBalancer *GlobalLoadBalancerV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	globalLoadBalancer, err = NewGlobalLoadBalancerV1(options)
	if err != nil {
		return
	}

	err = globalLoadBalancer.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = globalLoadBalancer.Service.SetServiceURL(options.URL)
	}
	return
}

// NewGlobalLoadBalancerV1 : constructs an instance of GlobalLoadBalancerV1 with passed in options.
func NewGlobalLoadBalancerV1(options *GlobalLoadBalancerV1Options) (service *GlobalLoadBalancerV1, err error) {
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

	service = &GlobalLoadBalancerV1{
		Service: baseService,
		Crn: options.Crn,
		ZoneIdentifier: options.ZoneIdentifier,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "globalLoadBalancer" suitable for processing requests.
func (globalLoadBalancer *GlobalLoadBalancerV1) Clone() *GlobalLoadBalancerV1 {
	if core.IsNil(globalLoadBalancer) {
		return nil
	}
	clone := *globalLoadBalancer
	clone.Service = globalLoadBalancer.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (globalLoadBalancer *GlobalLoadBalancerV1) SetServiceURL(url string) error {
	return globalLoadBalancer.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (globalLoadBalancer *GlobalLoadBalancerV1) GetServiceURL() string {
	return globalLoadBalancer.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (globalLoadBalancer *GlobalLoadBalancerV1) SetDefaultHeaders(headers http.Header) {
	globalLoadBalancer.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (globalLoadBalancer *GlobalLoadBalancerV1) SetEnableGzipCompression(enableGzip bool) {
	globalLoadBalancer.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (globalLoadBalancer *GlobalLoadBalancerV1) GetEnableGzipCompression() bool {
	return globalLoadBalancer.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (globalLoadBalancer *GlobalLoadBalancerV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	globalLoadBalancer.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (globalLoadBalancer *GlobalLoadBalancerV1) DisableRetries() {
	globalLoadBalancer.Service.DisableRetries()
}

// ListAllLoadBalancers : List all load balancers
// List configured load balancers.
func (globalLoadBalancer *GlobalLoadBalancerV1) ListAllLoadBalancers(listAllLoadBalancersOptions *ListAllLoadBalancersOptions) (result *ListLoadBalancersResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancer.ListAllLoadBalancersWithContext(context.Background(), listAllLoadBalancersOptions)
}

// ListAllLoadBalancersWithContext is an alternate form of the ListAllLoadBalancers method which supports a Context parameter
func (globalLoadBalancer *GlobalLoadBalancerV1) ListAllLoadBalancersWithContext(ctx context.Context, listAllLoadBalancersOptions *ListAllLoadBalancersOptions) (result *ListLoadBalancersResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllLoadBalancersOptions, "listAllLoadBalancersOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancer.Crn,
		"zone_identifier": *globalLoadBalancer.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancer.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancer.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/load_balancers`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllLoadBalancersOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer", "V1", "ListAllLoadBalancers")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancer.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListLoadBalancersResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancer : Create load balancer
// Create a load balancer for a given zone. The zone should be active before placing an order of a load balancer.
func (globalLoadBalancer *GlobalLoadBalancerV1) CreateLoadBalancer(createLoadBalancerOptions *CreateLoadBalancerOptions) (result *LoadBalancersResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancer.CreateLoadBalancerWithContext(context.Background(), createLoadBalancerOptions)
}

// CreateLoadBalancerWithContext is an alternate form of the CreateLoadBalancer method which supports a Context parameter
func (globalLoadBalancer *GlobalLoadBalancerV1) CreateLoadBalancerWithContext(ctx context.Context, createLoadBalancerOptions *CreateLoadBalancerOptions) (result *LoadBalancersResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLoadBalancerOptions, "createLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancer.Crn,
		"zone_identifier": *globalLoadBalancer.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancer.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancer.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/load_balancers`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer", "V1", "CreateLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLoadBalancerOptions.Name != nil {
		body["name"] = createLoadBalancerOptions.Name
	}
	if createLoadBalancerOptions.FallbackPool != nil {
		body["fallback_pool"] = createLoadBalancerOptions.FallbackPool
	}
	if createLoadBalancerOptions.DefaultPools != nil {
		body["default_pools"] = createLoadBalancerOptions.DefaultPools
	}
	if createLoadBalancerOptions.Description != nil {
		body["description"] = createLoadBalancerOptions.Description
	}
	if createLoadBalancerOptions.TTL != nil {
		body["ttl"] = createLoadBalancerOptions.TTL
	}
	if createLoadBalancerOptions.RegionPools != nil {
		body["region_pools"] = createLoadBalancerOptions.RegionPools
	}
	if createLoadBalancerOptions.PopPools != nil {
		body["pop_pools"] = createLoadBalancerOptions.PopPools
	}
	if createLoadBalancerOptions.Proxied != nil {
		body["proxied"] = createLoadBalancerOptions.Proxied
	}
	if createLoadBalancerOptions.Enabled != nil {
		body["enabled"] = createLoadBalancerOptions.Enabled
	}
	if createLoadBalancerOptions.SessionAffinity != nil {
		body["session_affinity"] = createLoadBalancerOptions.SessionAffinity
	}
	if createLoadBalancerOptions.SteeringPolicy != nil {
		body["steering_policy"] = createLoadBalancerOptions.SteeringPolicy
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
	response, err = globalLoadBalancer.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancersResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// EditLoadBalancer : Edit load balancer
// Edit porperties of an existing load balancer.
func (globalLoadBalancer *GlobalLoadBalancerV1) EditLoadBalancer(editLoadBalancerOptions *EditLoadBalancerOptions) (result *LoadBalancersResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancer.EditLoadBalancerWithContext(context.Background(), editLoadBalancerOptions)
}

// EditLoadBalancerWithContext is an alternate form of the EditLoadBalancer method which supports a Context parameter
func (globalLoadBalancer *GlobalLoadBalancerV1) EditLoadBalancerWithContext(ctx context.Context, editLoadBalancerOptions *EditLoadBalancerOptions) (result *LoadBalancersResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(editLoadBalancerOptions, "editLoadBalancerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(editLoadBalancerOptions, "editLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancer.Crn,
		"zone_identifier": *globalLoadBalancer.ZoneIdentifier,
		"load_balancer_identifier": *editLoadBalancerOptions.LoadBalancerIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancer.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancer.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/load_balancers/{load_balancer_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range editLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer", "V1", "EditLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if editLoadBalancerOptions.Name != nil {
		body["name"] = editLoadBalancerOptions.Name
	}
	if editLoadBalancerOptions.FallbackPool != nil {
		body["fallback_pool"] = editLoadBalancerOptions.FallbackPool
	}
	if editLoadBalancerOptions.DefaultPools != nil {
		body["default_pools"] = editLoadBalancerOptions.DefaultPools
	}
	if editLoadBalancerOptions.Description != nil {
		body["description"] = editLoadBalancerOptions.Description
	}
	if editLoadBalancerOptions.TTL != nil {
		body["ttl"] = editLoadBalancerOptions.TTL
	}
	if editLoadBalancerOptions.RegionPools != nil {
		body["region_pools"] = editLoadBalancerOptions.RegionPools
	}
	if editLoadBalancerOptions.PopPools != nil {
		body["pop_pools"] = editLoadBalancerOptions.PopPools
	}
	if editLoadBalancerOptions.Proxied != nil {
		body["proxied"] = editLoadBalancerOptions.Proxied
	}
	if editLoadBalancerOptions.Enabled != nil {
		body["enabled"] = editLoadBalancerOptions.Enabled
	}
	if editLoadBalancerOptions.SessionAffinity != nil {
		body["session_affinity"] = editLoadBalancerOptions.SessionAffinity
	}
	if editLoadBalancerOptions.SteeringPolicy != nil {
		body["steering_policy"] = editLoadBalancerOptions.SteeringPolicy
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
	response, err = globalLoadBalancer.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancersResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteLoadBalancer : Delete load balancer
// Delete a load balancer.
func (globalLoadBalancer *GlobalLoadBalancerV1) DeleteLoadBalancer(deleteLoadBalancerOptions *DeleteLoadBalancerOptions) (result *DeleteLoadBalancersResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancer.DeleteLoadBalancerWithContext(context.Background(), deleteLoadBalancerOptions)
}

// DeleteLoadBalancerWithContext is an alternate form of the DeleteLoadBalancer method which supports a Context parameter
func (globalLoadBalancer *GlobalLoadBalancerV1) DeleteLoadBalancerWithContext(ctx context.Context, deleteLoadBalancerOptions *DeleteLoadBalancerOptions) (result *DeleteLoadBalancersResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLoadBalancerOptions, "deleteLoadBalancerOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLoadBalancerOptions, "deleteLoadBalancerOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancer.Crn,
		"zone_identifier": *globalLoadBalancer.ZoneIdentifier,
		"load_balancer_identifier": *deleteLoadBalancerOptions.LoadBalancerIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancer.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancer.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/load_balancers/{load_balancer_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteLoadBalancerOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer", "V1", "DeleteLoadBalancer")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancer.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteLoadBalancersResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetLoadBalancerSettings : Get load balancer
// For a given zone identifier and load balancer id, get the load balancer settings.
func (globalLoadBalancer *GlobalLoadBalancerV1) GetLoadBalancerSettings(getLoadBalancerSettingsOptions *GetLoadBalancerSettingsOptions) (result *LoadBalancersResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancer.GetLoadBalancerSettingsWithContext(context.Background(), getLoadBalancerSettingsOptions)
}

// GetLoadBalancerSettingsWithContext is an alternate form of the GetLoadBalancerSettings method which supports a Context parameter
func (globalLoadBalancer *GlobalLoadBalancerV1) GetLoadBalancerSettingsWithContext(ctx context.Context, getLoadBalancerSettingsOptions *GetLoadBalancerSettingsOptions) (result *LoadBalancersResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLoadBalancerSettingsOptions, "getLoadBalancerSettingsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLoadBalancerSettingsOptions, "getLoadBalancerSettingsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancer.Crn,
		"zone_identifier": *globalLoadBalancer.ZoneIdentifier,
		"load_balancer_identifier": *getLoadBalancerSettingsOptions.LoadBalancerIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancer.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancer.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/load_balancers/{load_balancer_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLoadBalancerSettingsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer", "V1", "GetLoadBalancerSettings")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancer.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLoadBalancersResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancerOptions : The CreateLoadBalancer options.
type CreateLoadBalancerOptions struct {
	// name.
	Name *string `json:"name,omitempty"`

	// fallback pool.
	FallbackPool *string `json:"fallback_pool,omitempty"`

	// default pools.
	DefaultPools []string `json:"default_pools,omitempty"`

	// desc.
	Description *string `json:"description,omitempty"`

	// ttl.
	TTL *int64 `json:"ttl,omitempty"`

	// region pools.
	RegionPools interface{} `json:"region_pools,omitempty"`

	// pop pools.
	PopPools interface{} `json:"pop_pools,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// session affinity.
	SessionAffinity *string `json:"session_affinity,omitempty"`

	// steering policy.
	SteeringPolicy *string `json:"steering_policy,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateLoadBalancerOptions.SessionAffinity property.
// session affinity.
const (
	CreateLoadBalancerOptions_SessionAffinity_Cookie = "cookie"
	CreateLoadBalancerOptions_SessionAffinity_IpCookie = "ip_cookie"
	CreateLoadBalancerOptions_SessionAffinity_None = "none"
)

// Constants associated with the CreateLoadBalancerOptions.SteeringPolicy property.
// steering policy.
const (
	CreateLoadBalancerOptions_SteeringPolicy_DynamicLatency = "dynamic_latency"
	CreateLoadBalancerOptions_SteeringPolicy_Geo = "geo"
	CreateLoadBalancerOptions_SteeringPolicy_Off = "off"
	CreateLoadBalancerOptions_SteeringPolicy_Random = "random"
)

// NewCreateLoadBalancerOptions : Instantiate CreateLoadBalancerOptions
func (*GlobalLoadBalancerV1) NewCreateLoadBalancerOptions() *CreateLoadBalancerOptions {
	return &CreateLoadBalancerOptions{}
}

// SetName : Allow user to set Name
func (options *CreateLoadBalancerOptions) SetName(name string) *CreateLoadBalancerOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetFallbackPool : Allow user to set FallbackPool
func (options *CreateLoadBalancerOptions) SetFallbackPool(fallbackPool string) *CreateLoadBalancerOptions {
	options.FallbackPool = core.StringPtr(fallbackPool)
	return options
}

// SetDefaultPools : Allow user to set DefaultPools
func (options *CreateLoadBalancerOptions) SetDefaultPools(defaultPools []string) *CreateLoadBalancerOptions {
	options.DefaultPools = defaultPools
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateLoadBalancerOptions) SetDescription(description string) *CreateLoadBalancerOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTTL : Allow user to set TTL
func (options *CreateLoadBalancerOptions) SetTTL(ttl int64) *CreateLoadBalancerOptions {
	options.TTL = core.Int64Ptr(ttl)
	return options
}

// SetRegionPools : Allow user to set RegionPools
func (options *CreateLoadBalancerOptions) SetRegionPools(regionPools interface{}) *CreateLoadBalancerOptions {
	options.RegionPools = regionPools
	return options
}

// SetPopPools : Allow user to set PopPools
func (options *CreateLoadBalancerOptions) SetPopPools(popPools interface{}) *CreateLoadBalancerOptions {
	options.PopPools = popPools
	return options
}

// SetProxied : Allow user to set Proxied
func (options *CreateLoadBalancerOptions) SetProxied(proxied bool) *CreateLoadBalancerOptions {
	options.Proxied = core.BoolPtr(proxied)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *CreateLoadBalancerOptions) SetEnabled(enabled bool) *CreateLoadBalancerOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetSessionAffinity : Allow user to set SessionAffinity
func (options *CreateLoadBalancerOptions) SetSessionAffinity(sessionAffinity string) *CreateLoadBalancerOptions {
	options.SessionAffinity = core.StringPtr(sessionAffinity)
	return options
}

// SetSteeringPolicy : Allow user to set SteeringPolicy
func (options *CreateLoadBalancerOptions) SetSteeringPolicy(steeringPolicy string) *CreateLoadBalancerOptions {
	options.SteeringPolicy = core.StringPtr(steeringPolicy)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLoadBalancerOptions) SetHeaders(param map[string]string) *CreateLoadBalancerOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancerOptions : The DeleteLoadBalancer options.
type DeleteLoadBalancerOptions struct {
	// load balancer identifier.
	LoadBalancerIdentifier *string `json:"load_balancer_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLoadBalancerOptions : Instantiate DeleteLoadBalancerOptions
func (*GlobalLoadBalancerV1) NewDeleteLoadBalancerOptions(loadBalancerIdentifier string) *DeleteLoadBalancerOptions {
	return &DeleteLoadBalancerOptions{
		LoadBalancerIdentifier: core.StringPtr(loadBalancerIdentifier),
	}
}

// SetLoadBalancerIdentifier : Allow user to set LoadBalancerIdentifier
func (options *DeleteLoadBalancerOptions) SetLoadBalancerIdentifier(loadBalancerIdentifier string) *DeleteLoadBalancerOptions {
	options.LoadBalancerIdentifier = core.StringPtr(loadBalancerIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLoadBalancerOptions) SetHeaders(param map[string]string) *DeleteLoadBalancerOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancersRespResult : result.
type DeleteLoadBalancersRespResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteLoadBalancersRespResult unmarshals an instance of DeleteLoadBalancersRespResult from the specified map of raw messages.
func UnmarshalDeleteLoadBalancersRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLoadBalancersRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EditLoadBalancerOptions : The EditLoadBalancer options.
type EditLoadBalancerOptions struct {
	// load balancer identifier.
	LoadBalancerIdentifier *string `json:"load_balancer_identifier" validate:"required,ne="`

	// name.
	Name *string `json:"name,omitempty"`

	// fallback pool.
	FallbackPool *string `json:"fallback_pool,omitempty"`

	// default pools.
	DefaultPools []string `json:"default_pools,omitempty"`

	// desc.
	Description *string `json:"description,omitempty"`

	// ttl.
	TTL *int64 `json:"ttl,omitempty"`

	// region pools.
	RegionPools interface{} `json:"region_pools,omitempty"`

	// pop pools.
	PopPools interface{} `json:"pop_pools,omitempty"`

	// proxied.
	Proxied *bool `json:"proxied,omitempty"`

	// enabled/disabled.
	Enabled *bool `json:"enabled,omitempty"`

	// session affinity.
	SessionAffinity *string `json:"session_affinity,omitempty"`

	// steering policy.
	SteeringPolicy *string `json:"steering_policy,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the EditLoadBalancerOptions.SessionAffinity property.
// session affinity.
const (
	EditLoadBalancerOptions_SessionAffinity_Cookie = "cookie"
	EditLoadBalancerOptions_SessionAffinity_IpCookie = "ip_cookie"
	EditLoadBalancerOptions_SessionAffinity_None = "none"
)

// Constants associated with the EditLoadBalancerOptions.SteeringPolicy property.
// steering policy.
const (
	EditLoadBalancerOptions_SteeringPolicy_DynamicLatency = "dynamic_latency"
	EditLoadBalancerOptions_SteeringPolicy_Geo = "geo"
	EditLoadBalancerOptions_SteeringPolicy_Off = "off"
	EditLoadBalancerOptions_SteeringPolicy_Random = "random"
)

// NewEditLoadBalancerOptions : Instantiate EditLoadBalancerOptions
func (*GlobalLoadBalancerV1) NewEditLoadBalancerOptions(loadBalancerIdentifier string) *EditLoadBalancerOptions {
	return &EditLoadBalancerOptions{
		LoadBalancerIdentifier: core.StringPtr(loadBalancerIdentifier),
	}
}

// SetLoadBalancerIdentifier : Allow user to set LoadBalancerIdentifier
func (options *EditLoadBalancerOptions) SetLoadBalancerIdentifier(loadBalancerIdentifier string) *EditLoadBalancerOptions {
	options.LoadBalancerIdentifier = core.StringPtr(loadBalancerIdentifier)
	return options
}

// SetName : Allow user to set Name
func (options *EditLoadBalancerOptions) SetName(name string) *EditLoadBalancerOptions {
	options.Name = core.StringPtr(name)
	return options
}

// SetFallbackPool : Allow user to set FallbackPool
func (options *EditLoadBalancerOptions) SetFallbackPool(fallbackPool string) *EditLoadBalancerOptions {
	options.FallbackPool = core.StringPtr(fallbackPool)
	return options
}

// SetDefaultPools : Allow user to set DefaultPools
func (options *EditLoadBalancerOptions) SetDefaultPools(defaultPools []string) *EditLoadBalancerOptions {
	options.DefaultPools = defaultPools
	return options
}

// SetDescription : Allow user to set Description
func (options *EditLoadBalancerOptions) SetDescription(description string) *EditLoadBalancerOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetTTL : Allow user to set TTL
func (options *EditLoadBalancerOptions) SetTTL(ttl int64) *EditLoadBalancerOptions {
	options.TTL = core.Int64Ptr(ttl)
	return options
}

// SetRegionPools : Allow user to set RegionPools
func (options *EditLoadBalancerOptions) SetRegionPools(regionPools interface{}) *EditLoadBalancerOptions {
	options.RegionPools = regionPools
	return options
}

// SetPopPools : Allow user to set PopPools
func (options *EditLoadBalancerOptions) SetPopPools(popPools interface{}) *EditLoadBalancerOptions {
	options.PopPools = popPools
	return options
}

// SetProxied : Allow user to set Proxied
func (options *EditLoadBalancerOptions) SetProxied(proxied bool) *EditLoadBalancerOptions {
	options.Proxied = core.BoolPtr(proxied)
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *EditLoadBalancerOptions) SetEnabled(enabled bool) *EditLoadBalancerOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetSessionAffinity : Allow user to set SessionAffinity
func (options *EditLoadBalancerOptions) SetSessionAffinity(sessionAffinity string) *EditLoadBalancerOptions {
	options.SessionAffinity = core.StringPtr(sessionAffinity)
	return options
}

// SetSteeringPolicy : Allow user to set SteeringPolicy
func (options *EditLoadBalancerOptions) SetSteeringPolicy(steeringPolicy string) *EditLoadBalancerOptions {
	options.SteeringPolicy = core.StringPtr(steeringPolicy)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *EditLoadBalancerOptions) SetHeaders(param map[string]string) *EditLoadBalancerOptions {
	options.Headers = param
	return options
}

// GetLoadBalancerSettingsOptions : The GetLoadBalancerSettings options.
type GetLoadBalancerSettingsOptions struct {
	// load balancer identifier.
	LoadBalancerIdentifier *string `json:"load_balancer_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLoadBalancerSettingsOptions : Instantiate GetLoadBalancerSettingsOptions
func (*GlobalLoadBalancerV1) NewGetLoadBalancerSettingsOptions(loadBalancerIdentifier string) *GetLoadBalancerSettingsOptions {
	return &GetLoadBalancerSettingsOptions{
		LoadBalancerIdentifier: core.StringPtr(loadBalancerIdentifier),
	}
}

// SetLoadBalancerIdentifier : Allow user to set LoadBalancerIdentifier
func (options *GetLoadBalancerSettingsOptions) SetLoadBalancerIdentifier(loadBalancerIdentifier string) *GetLoadBalancerSettingsOptions {
	options.LoadBalancerIdentifier = core.StringPtr(loadBalancerIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetLoadBalancerSettingsOptions) SetHeaders(param map[string]string) *GetLoadBalancerSettingsOptions {
	options.Headers = param
	return options
}

// ListAllLoadBalancersOptions : The ListAllLoadBalancers options.
type ListAllLoadBalancersOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllLoadBalancersOptions : Instantiate ListAllLoadBalancersOptions
func (*GlobalLoadBalancerV1) NewListAllLoadBalancersOptions() *ListAllLoadBalancersOptions {
	return &ListAllLoadBalancersOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListAllLoadBalancersOptions) SetHeaders(param map[string]string) *ListAllLoadBalancersOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancersResp : delete load balancers response.
type DeleteLoadBalancersResp struct {
	// success respose.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result *DeleteLoadBalancersRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteLoadBalancersResp unmarshals an instance of DeleteLoadBalancersResp from the specified map of raw messages.
func UnmarshalDeleteLoadBalancersResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteLoadBalancersResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteLoadBalancersRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListLoadBalancersResp : load balancer list response.
type ListLoadBalancersResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result []LoadBalancerPack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListLoadBalancersResp unmarshals an instance of ListLoadBalancersResp from the specified map of raw messages.
func UnmarshalListLoadBalancersResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListLoadBalancersResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLoadBalancerPack)
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

// LoadBalancerPack : loadbalancer pack.
type LoadBalancerPack struct {
	// identifier.
	ID *string `json:"id" validate:"required"`

	// created date.
	CreatedOn *string `json:"created_on" validate:"required"`

	// modified date.
	ModifiedOn *string `json:"modified_on" validate:"required"`

	// desc.
	Description *string `json:"description" validate:"required"`

	// name.
	Name *string `json:"name" validate:"required"`

	// ttl.
	TTL *int64 `json:"ttl" validate:"required"`

	// fallback pool.
	FallbackPool *string `json:"fallback_pool" validate:"required"`

	// default pools.
	DefaultPools []string `json:"default_pools" validate:"required"`

	// region pools.
	RegionPools interface{} `json:"region_pools" validate:"required"`

	// pop pools.
	PopPools interface{} `json:"pop_pools" validate:"required"`

	// proxied.
	Proxied *bool `json:"proxied" validate:"required"`

	// enabled/disabled.
	Enabled *bool `json:"enabled" validate:"required"`

	// session affinity.
	SessionAffinity *string `json:"session_affinity" validate:"required"`

	// steering policy.
	SteeringPolicy *string `json:"steering_policy" validate:"required"`
}

// Constants associated with the LoadBalancerPack.SessionAffinity property.
// session affinity.
const (
	LoadBalancerPack_SessionAffinity_Cookie = "cookie"
	LoadBalancerPack_SessionAffinity_IpCookie = "ip_cookie"
	LoadBalancerPack_SessionAffinity_None = "none"
)

// Constants associated with the LoadBalancerPack.SteeringPolicy property.
// steering policy.
const (
	LoadBalancerPack_SteeringPolicy_DynamicLatency = "dynamic_latency"
	LoadBalancerPack_SteeringPolicy_Geo = "geo"
	LoadBalancerPack_SteeringPolicy_Off = "off"
	LoadBalancerPack_SteeringPolicy_Random = "random"
)


// UnmarshalLoadBalancerPack unmarshals an instance of LoadBalancerPack from the specified map of raw messages.
func UnmarshalLoadBalancerPack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancerPack)
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
	err = core.UnmarshalPrimitive(m, "ttl", &obj.TTL)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "fallback_pool", &obj.FallbackPool)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "default_pools", &obj.DefaultPools)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "region_pools", &obj.RegionPools)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pop_pools", &obj.PopPools)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "proxied", &obj.Proxied)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "session_affinity", &obj.SessionAffinity)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "steering_policy", &obj.SteeringPolicy)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LoadBalancersResp : load balancer response.
type LoadBalancersResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// loadbalancer pack.
	Result *LoadBalancerPack `json:"result" validate:"required"`
}


// UnmarshalLoadBalancersResp unmarshals an instance of LoadBalancersResp from the specified map of raw messages.
func UnmarshalLoadBalancersResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LoadBalancersResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalLoadBalancerPack)
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
