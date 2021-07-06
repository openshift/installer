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
 

// Package globalloadbalancermonitorv1 : Operations and models for the GlobalLoadBalancerMonitorV1 service
package globalloadbalancermonitorv1

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

// GlobalLoadBalancerMonitorV1 : Global Load Balancer Monitor
//
// Version: 1.0.1
type GlobalLoadBalancerMonitorV1 struct {
	Service *core.BaseService

	// Full CRN of the service instance.
	Crn *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "global_load_balancer_monitor"

// GlobalLoadBalancerMonitorV1Options : Service options
type GlobalLoadBalancerMonitorV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full CRN of the service instance.
	Crn *string `validate:"required"`
}

// NewGlobalLoadBalancerMonitorV1UsingExternalConfig : constructs an instance of GlobalLoadBalancerMonitorV1 with passed in options and external configuration.
func NewGlobalLoadBalancerMonitorV1UsingExternalConfig(options *GlobalLoadBalancerMonitorV1Options) (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	globalLoadBalancerMonitor, err = NewGlobalLoadBalancerMonitorV1(options)
	if err != nil {
		return
	}

	err = globalLoadBalancerMonitor.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = globalLoadBalancerMonitor.Service.SetServiceURL(options.URL)
	}
	return
}

// NewGlobalLoadBalancerMonitorV1 : constructs an instance of GlobalLoadBalancerMonitorV1 with passed in options.
func NewGlobalLoadBalancerMonitorV1(options *GlobalLoadBalancerMonitorV1Options) (service *GlobalLoadBalancerMonitorV1, err error) {
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

	service = &GlobalLoadBalancerMonitorV1{
		Service: baseService,
		Crn: options.Crn,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "globalLoadBalancerMonitor" suitable for processing requests.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) Clone() *GlobalLoadBalancerMonitorV1 {
	if core.IsNil(globalLoadBalancerMonitor) {
		return nil
	}
	clone := *globalLoadBalancerMonitor
	clone.Service = globalLoadBalancerMonitor.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) SetServiceURL(url string) error {
	return globalLoadBalancerMonitor.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) GetServiceURL() string {
	return globalLoadBalancerMonitor.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) SetDefaultHeaders(headers http.Header) {
	globalLoadBalancerMonitor.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) SetEnableGzipCompression(enableGzip bool) {
	globalLoadBalancerMonitor.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) GetEnableGzipCompression() bool {
	return globalLoadBalancerMonitor.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	globalLoadBalancerMonitor.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) DisableRetries() {
	globalLoadBalancerMonitor.Service.DisableRetries()
}

// ListAllLoadBalancerMonitors : List all load balancer monitors
// List configured load balancer monitors for a user.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) ListAllLoadBalancerMonitors(listAllLoadBalancerMonitorsOptions *ListAllLoadBalancerMonitorsOptions) (result *ListMonitorResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerMonitor.ListAllLoadBalancerMonitorsWithContext(context.Background(), listAllLoadBalancerMonitorsOptions)
}

// ListAllLoadBalancerMonitorsWithContext is an alternate form of the ListAllLoadBalancerMonitors method which supports a Context parameter
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) ListAllLoadBalancerMonitorsWithContext(ctx context.Context, listAllLoadBalancerMonitorsOptions *ListAllLoadBalancerMonitorsOptions) (result *ListMonitorResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listAllLoadBalancerMonitorsOptions, "listAllLoadBalancerMonitorsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerMonitor.Crn,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerMonitor.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerMonitor.Service.Options.URL, `/v1/{crn}/load_balancers/monitors`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listAllLoadBalancerMonitorsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_monitor", "V1", "ListAllLoadBalancerMonitors")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancerMonitor.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalListMonitorResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancerMonitor : Create load balancer monitor
// Create a load balancer monitor for a given service instance.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) CreateLoadBalancerMonitor(createLoadBalancerMonitorOptions *CreateLoadBalancerMonitorOptions) (result *MonitorResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerMonitor.CreateLoadBalancerMonitorWithContext(context.Background(), createLoadBalancerMonitorOptions)
}

// CreateLoadBalancerMonitorWithContext is an alternate form of the CreateLoadBalancerMonitor method which supports a Context parameter
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) CreateLoadBalancerMonitorWithContext(ctx context.Context, createLoadBalancerMonitorOptions *CreateLoadBalancerMonitorOptions) (result *MonitorResp, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(createLoadBalancerMonitorOptions, "createLoadBalancerMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerMonitor.Crn,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerMonitor.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerMonitor.Service.Options.URL, `/v1/{crn}/load_balancers/monitors`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createLoadBalancerMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_monitor", "V1", "CreateLoadBalancerMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLoadBalancerMonitorOptions.Type != nil {
		body["type"] = createLoadBalancerMonitorOptions.Type
	}
	if createLoadBalancerMonitorOptions.Description != nil {
		body["description"] = createLoadBalancerMonitorOptions.Description
	}
	if createLoadBalancerMonitorOptions.Method != nil {
		body["method"] = createLoadBalancerMonitorOptions.Method
	}
	if createLoadBalancerMonitorOptions.Port != nil {
		body["port"] = createLoadBalancerMonitorOptions.Port
	}
	if createLoadBalancerMonitorOptions.Path != nil {
		body["path"] = createLoadBalancerMonitorOptions.Path
	}
	if createLoadBalancerMonitorOptions.Timeout != nil {
		body["timeout"] = createLoadBalancerMonitorOptions.Timeout
	}
	if createLoadBalancerMonitorOptions.Retries != nil {
		body["retries"] = createLoadBalancerMonitorOptions.Retries
	}
	if createLoadBalancerMonitorOptions.Interval != nil {
		body["interval"] = createLoadBalancerMonitorOptions.Interval
	}
	if createLoadBalancerMonitorOptions.ExpectedCodes != nil {
		body["expected_codes"] = createLoadBalancerMonitorOptions.ExpectedCodes
	}
	if createLoadBalancerMonitorOptions.FollowRedirects != nil {
		body["follow_redirects"] = createLoadBalancerMonitorOptions.FollowRedirects
	}
	if createLoadBalancerMonitorOptions.ExpectedBody != nil {
		body["expected_body"] = createLoadBalancerMonitorOptions.ExpectedBody
	}
	if createLoadBalancerMonitorOptions.AllowInsecure != nil {
		body["allow_insecure"] = createLoadBalancerMonitorOptions.AllowInsecure
	}
	if createLoadBalancerMonitorOptions.Header != nil {
		body["header"] = createLoadBalancerMonitorOptions.Header
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
	response, err = globalLoadBalancerMonitor.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitorResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// EditLoadBalancerMonitor : Edit load balancer monitor
// Edit porperties of an existing load balancer monitor.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) EditLoadBalancerMonitor(editLoadBalancerMonitorOptions *EditLoadBalancerMonitorOptions) (result *MonitorResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerMonitor.EditLoadBalancerMonitorWithContext(context.Background(), editLoadBalancerMonitorOptions)
}

// EditLoadBalancerMonitorWithContext is an alternate form of the EditLoadBalancerMonitor method which supports a Context parameter
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) EditLoadBalancerMonitorWithContext(ctx context.Context, editLoadBalancerMonitorOptions *EditLoadBalancerMonitorOptions) (result *MonitorResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(editLoadBalancerMonitorOptions, "editLoadBalancerMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(editLoadBalancerMonitorOptions, "editLoadBalancerMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerMonitor.Crn,
		"monitor_identifier": *editLoadBalancerMonitorOptions.MonitorIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerMonitor.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerMonitor.Service.Options.URL, `/v1/{crn}/load_balancers/monitors/{monitor_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range editLoadBalancerMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_monitor", "V1", "EditLoadBalancerMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if editLoadBalancerMonitorOptions.Type != nil {
		body["type"] = editLoadBalancerMonitorOptions.Type
	}
	if editLoadBalancerMonitorOptions.Description != nil {
		body["description"] = editLoadBalancerMonitorOptions.Description
	}
	if editLoadBalancerMonitorOptions.Method != nil {
		body["method"] = editLoadBalancerMonitorOptions.Method
	}
	if editLoadBalancerMonitorOptions.Port != nil {
		body["port"] = editLoadBalancerMonitorOptions.Port
	}
	if editLoadBalancerMonitorOptions.Path != nil {
		body["path"] = editLoadBalancerMonitorOptions.Path
	}
	if editLoadBalancerMonitorOptions.Timeout != nil {
		body["timeout"] = editLoadBalancerMonitorOptions.Timeout
	}
	if editLoadBalancerMonitorOptions.Retries != nil {
		body["retries"] = editLoadBalancerMonitorOptions.Retries
	}
	if editLoadBalancerMonitorOptions.Interval != nil {
		body["interval"] = editLoadBalancerMonitorOptions.Interval
	}
	if editLoadBalancerMonitorOptions.ExpectedCodes != nil {
		body["expected_codes"] = editLoadBalancerMonitorOptions.ExpectedCodes
	}
	if editLoadBalancerMonitorOptions.FollowRedirects != nil {
		body["follow_redirects"] = editLoadBalancerMonitorOptions.FollowRedirects
	}
	if editLoadBalancerMonitorOptions.ExpectedBody != nil {
		body["expected_body"] = editLoadBalancerMonitorOptions.ExpectedBody
	}
	if editLoadBalancerMonitorOptions.AllowInsecure != nil {
		body["allow_insecure"] = editLoadBalancerMonitorOptions.AllowInsecure
	}
	if editLoadBalancerMonitorOptions.Header != nil {
		body["header"] = editLoadBalancerMonitorOptions.Header
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
	response, err = globalLoadBalancerMonitor.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitorResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteLoadBalancerMonitor : Delete load balancer monitor
// Delete a load balancer monitor.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) DeleteLoadBalancerMonitor(deleteLoadBalancerMonitorOptions *DeleteLoadBalancerMonitorOptions) (result *DeleteMonitorResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerMonitor.DeleteLoadBalancerMonitorWithContext(context.Background(), deleteLoadBalancerMonitorOptions)
}

// DeleteLoadBalancerMonitorWithContext is an alternate form of the DeleteLoadBalancerMonitor method which supports a Context parameter
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) DeleteLoadBalancerMonitorWithContext(ctx context.Context, deleteLoadBalancerMonitorOptions *DeleteLoadBalancerMonitorOptions) (result *DeleteMonitorResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLoadBalancerMonitorOptions, "deleteLoadBalancerMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLoadBalancerMonitorOptions, "deleteLoadBalancerMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerMonitor.Crn,
		"monitor_identifier": *deleteLoadBalancerMonitorOptions.MonitorIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerMonitor.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerMonitor.Service.Options.URL, `/v1/{crn}/load_balancers/monitors/{monitor_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteLoadBalancerMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_monitor", "V1", "DeleteLoadBalancerMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancerMonitor.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDeleteMonitorResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetLoadBalancerMonitor : Get load balancer monitor
// For a given service instance and load balancer monitor id, get the monitor details.
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) GetLoadBalancerMonitor(getLoadBalancerMonitorOptions *GetLoadBalancerMonitorOptions) (result *MonitorResp, response *core.DetailedResponse, err error) {
	return globalLoadBalancerMonitor.GetLoadBalancerMonitorWithContext(context.Background(), getLoadBalancerMonitorOptions)
}

// GetLoadBalancerMonitorWithContext is an alternate form of the GetLoadBalancerMonitor method which supports a Context parameter
func (globalLoadBalancerMonitor *GlobalLoadBalancerMonitorV1) GetLoadBalancerMonitorWithContext(ctx context.Context, getLoadBalancerMonitorOptions *GetLoadBalancerMonitorOptions) (result *MonitorResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLoadBalancerMonitorOptions, "getLoadBalancerMonitorOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLoadBalancerMonitorOptions, "getLoadBalancerMonitorOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *globalLoadBalancerMonitor.Crn,
		"monitor_identifier": *getLoadBalancerMonitorOptions.MonitorIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = globalLoadBalancerMonitor.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(globalLoadBalancerMonitor.Service.Options.URL, `/v1/{crn}/load_balancers/monitors/{monitor_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLoadBalancerMonitorOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("global_load_balancer_monitor", "V1", "GetLoadBalancerMonitor")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = globalLoadBalancerMonitor.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalMonitorResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateLoadBalancerMonitorOptions : The CreateLoadBalancerMonitor options.
type CreateLoadBalancerMonitorOptions struct {
	// http type.
	Type *string `json:"type,omitempty"`

	// login page monitor.
	Description *string `json:"description,omitempty"`

	// method.
	Method *string `json:"method,omitempty"`

	// port number.
	Port *int64 `json:"port,omitempty"`

	// path.
	Path *string `json:"path,omitempty"`

	// timeout count.
	Timeout *int64 `json:"timeout,omitempty"`

	// retry count.
	Retries *int64 `json:"retries,omitempty"`

	// interval.
	Interval *int64 `json:"interval,omitempty"`

	// expected codes.
	ExpectedCodes *string `json:"expected_codes,omitempty"`

	// follow redirects.
	FollowRedirects *bool `json:"follow_redirects,omitempty"`

	// expected body.
	ExpectedBody *string `json:"expected_body,omitempty"`

	// allow insecure.
	AllowInsecure *bool `json:"allow_insecure,omitempty"`

	// header.
	Header map[string][]string `json:"header,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateLoadBalancerMonitorOptions : Instantiate CreateLoadBalancerMonitorOptions
func (*GlobalLoadBalancerMonitorV1) NewCreateLoadBalancerMonitorOptions() *CreateLoadBalancerMonitorOptions {
	return &CreateLoadBalancerMonitorOptions{}
}

// SetType : Allow user to set Type
func (options *CreateLoadBalancerMonitorOptions) SetType(typeVar string) *CreateLoadBalancerMonitorOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetDescription : Allow user to set Description
func (options *CreateLoadBalancerMonitorOptions) SetDescription(description string) *CreateLoadBalancerMonitorOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetMethod : Allow user to set Method
func (options *CreateLoadBalancerMonitorOptions) SetMethod(method string) *CreateLoadBalancerMonitorOptions {
	options.Method = core.StringPtr(method)
	return options
}

// SetPort : Allow user to set Port
func (options *CreateLoadBalancerMonitorOptions) SetPort(port int64) *CreateLoadBalancerMonitorOptions {
	options.Port = core.Int64Ptr(port)
	return options
}

// SetPath : Allow user to set Path
func (options *CreateLoadBalancerMonitorOptions) SetPath(path string) *CreateLoadBalancerMonitorOptions {
	options.Path = core.StringPtr(path)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *CreateLoadBalancerMonitorOptions) SetTimeout(timeout int64) *CreateLoadBalancerMonitorOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetRetries : Allow user to set Retries
func (options *CreateLoadBalancerMonitorOptions) SetRetries(retries int64) *CreateLoadBalancerMonitorOptions {
	options.Retries = core.Int64Ptr(retries)
	return options
}

// SetInterval : Allow user to set Interval
func (options *CreateLoadBalancerMonitorOptions) SetInterval(interval int64) *CreateLoadBalancerMonitorOptions {
	options.Interval = core.Int64Ptr(interval)
	return options
}

// SetExpectedCodes : Allow user to set ExpectedCodes
func (options *CreateLoadBalancerMonitorOptions) SetExpectedCodes(expectedCodes string) *CreateLoadBalancerMonitorOptions {
	options.ExpectedCodes = core.StringPtr(expectedCodes)
	return options
}

// SetFollowRedirects : Allow user to set FollowRedirects
func (options *CreateLoadBalancerMonitorOptions) SetFollowRedirects(followRedirects bool) *CreateLoadBalancerMonitorOptions {
	options.FollowRedirects = core.BoolPtr(followRedirects)
	return options
}

// SetExpectedBody : Allow user to set ExpectedBody
func (options *CreateLoadBalancerMonitorOptions) SetExpectedBody(expectedBody string) *CreateLoadBalancerMonitorOptions {
	options.ExpectedBody = core.StringPtr(expectedBody)
	return options
}

// SetAllowInsecure : Allow user to set AllowInsecure
func (options *CreateLoadBalancerMonitorOptions) SetAllowInsecure(allowInsecure bool) *CreateLoadBalancerMonitorOptions {
	options.AllowInsecure = core.BoolPtr(allowInsecure)
	return options
}

// SetHeader : Allow user to set Header
func (options *CreateLoadBalancerMonitorOptions) SetHeader(header map[string][]string) *CreateLoadBalancerMonitorOptions {
	options.Header = header
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLoadBalancerMonitorOptions) SetHeaders(param map[string]string) *CreateLoadBalancerMonitorOptions {
	options.Headers = param
	return options
}

// DeleteLoadBalancerMonitorOptions : The DeleteLoadBalancerMonitor options.
type DeleteLoadBalancerMonitorOptions struct {
	// monitor identifier.
	MonitorIdentifier *string `json:"monitor_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLoadBalancerMonitorOptions : Instantiate DeleteLoadBalancerMonitorOptions
func (*GlobalLoadBalancerMonitorV1) NewDeleteLoadBalancerMonitorOptions(monitorIdentifier string) *DeleteLoadBalancerMonitorOptions {
	return &DeleteLoadBalancerMonitorOptions{
		MonitorIdentifier: core.StringPtr(monitorIdentifier),
	}
}

// SetMonitorIdentifier : Allow user to set MonitorIdentifier
func (options *DeleteLoadBalancerMonitorOptions) SetMonitorIdentifier(monitorIdentifier string) *DeleteLoadBalancerMonitorOptions {
	options.MonitorIdentifier = core.StringPtr(monitorIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLoadBalancerMonitorOptions) SetHeaders(param map[string]string) *DeleteLoadBalancerMonitorOptions {
	options.Headers = param
	return options
}

// DeleteMonitorRespResult : result.
type DeleteMonitorRespResult struct {
	// identifier.
	ID *string `json:"id" validate:"required"`
}


// UnmarshalDeleteMonitorRespResult unmarshals an instance of DeleteMonitorRespResult from the specified map of raw messages.
func UnmarshalDeleteMonitorRespResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteMonitorRespResult)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EditLoadBalancerMonitorOptions : The EditLoadBalancerMonitor options.
type EditLoadBalancerMonitorOptions struct {
	// monitor identifier.
	MonitorIdentifier *string `json:"monitor_identifier" validate:"required,ne="`

	// http type.
	Type *string `json:"type,omitempty"`

	// login page monitor.
	Description *string `json:"description,omitempty"`

	// method.
	Method *string `json:"method,omitempty"`

	// port number.
	Port *int64 `json:"port,omitempty"`

	// path.
	Path *string `json:"path,omitempty"`

	// timeout count.
	Timeout *int64 `json:"timeout,omitempty"`

	// retry count.
	Retries *int64 `json:"retries,omitempty"`

	// interval.
	Interval *int64 `json:"interval,omitempty"`

	// expected codes.
	ExpectedCodes *string `json:"expected_codes,omitempty"`

	// follow redirects.
	FollowRedirects *bool `json:"follow_redirects,omitempty"`

	// expected body.
	ExpectedBody *string `json:"expected_body,omitempty"`

	// allow insecure.
	AllowInsecure *bool `json:"allow_insecure,omitempty"`

	// header.
	Header map[string][]string `json:"header,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewEditLoadBalancerMonitorOptions : Instantiate EditLoadBalancerMonitorOptions
func (*GlobalLoadBalancerMonitorV1) NewEditLoadBalancerMonitorOptions(monitorIdentifier string) *EditLoadBalancerMonitorOptions {
	return &EditLoadBalancerMonitorOptions{
		MonitorIdentifier: core.StringPtr(monitorIdentifier),
	}
}

// SetMonitorIdentifier : Allow user to set MonitorIdentifier
func (options *EditLoadBalancerMonitorOptions) SetMonitorIdentifier(monitorIdentifier string) *EditLoadBalancerMonitorOptions {
	options.MonitorIdentifier = core.StringPtr(monitorIdentifier)
	return options
}

// SetType : Allow user to set Type
func (options *EditLoadBalancerMonitorOptions) SetType(typeVar string) *EditLoadBalancerMonitorOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetDescription : Allow user to set Description
func (options *EditLoadBalancerMonitorOptions) SetDescription(description string) *EditLoadBalancerMonitorOptions {
	options.Description = core.StringPtr(description)
	return options
}

// SetMethod : Allow user to set Method
func (options *EditLoadBalancerMonitorOptions) SetMethod(method string) *EditLoadBalancerMonitorOptions {
	options.Method = core.StringPtr(method)
	return options
}

// SetPort : Allow user to set Port
func (options *EditLoadBalancerMonitorOptions) SetPort(port int64) *EditLoadBalancerMonitorOptions {
	options.Port = core.Int64Ptr(port)
	return options
}

// SetPath : Allow user to set Path
func (options *EditLoadBalancerMonitorOptions) SetPath(path string) *EditLoadBalancerMonitorOptions {
	options.Path = core.StringPtr(path)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *EditLoadBalancerMonitorOptions) SetTimeout(timeout int64) *EditLoadBalancerMonitorOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetRetries : Allow user to set Retries
func (options *EditLoadBalancerMonitorOptions) SetRetries(retries int64) *EditLoadBalancerMonitorOptions {
	options.Retries = core.Int64Ptr(retries)
	return options
}

// SetInterval : Allow user to set Interval
func (options *EditLoadBalancerMonitorOptions) SetInterval(interval int64) *EditLoadBalancerMonitorOptions {
	options.Interval = core.Int64Ptr(interval)
	return options
}

// SetExpectedCodes : Allow user to set ExpectedCodes
func (options *EditLoadBalancerMonitorOptions) SetExpectedCodes(expectedCodes string) *EditLoadBalancerMonitorOptions {
	options.ExpectedCodes = core.StringPtr(expectedCodes)
	return options
}

// SetFollowRedirects : Allow user to set FollowRedirects
func (options *EditLoadBalancerMonitorOptions) SetFollowRedirects(followRedirects bool) *EditLoadBalancerMonitorOptions {
	options.FollowRedirects = core.BoolPtr(followRedirects)
	return options
}

// SetExpectedBody : Allow user to set ExpectedBody
func (options *EditLoadBalancerMonitorOptions) SetExpectedBody(expectedBody string) *EditLoadBalancerMonitorOptions {
	options.ExpectedBody = core.StringPtr(expectedBody)
	return options
}

// SetAllowInsecure : Allow user to set AllowInsecure
func (options *EditLoadBalancerMonitorOptions) SetAllowInsecure(allowInsecure bool) *EditLoadBalancerMonitorOptions {
	options.AllowInsecure = core.BoolPtr(allowInsecure)
	return options
}

// SetHeader : Allow user to set Header
func (options *EditLoadBalancerMonitorOptions) SetHeader(header map[string][]string) *EditLoadBalancerMonitorOptions {
	options.Header = header
	return options
}

// SetHeaders : Allow user to set Headers
func (options *EditLoadBalancerMonitorOptions) SetHeaders(param map[string]string) *EditLoadBalancerMonitorOptions {
	options.Headers = param
	return options
}

// GetLoadBalancerMonitorOptions : The GetLoadBalancerMonitor options.
type GetLoadBalancerMonitorOptions struct {
	// monitor identifier.
	MonitorIdentifier *string `json:"monitor_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLoadBalancerMonitorOptions : Instantiate GetLoadBalancerMonitorOptions
func (*GlobalLoadBalancerMonitorV1) NewGetLoadBalancerMonitorOptions(monitorIdentifier string) *GetLoadBalancerMonitorOptions {
	return &GetLoadBalancerMonitorOptions{
		MonitorIdentifier: core.StringPtr(monitorIdentifier),
	}
}

// SetMonitorIdentifier : Allow user to set MonitorIdentifier
func (options *GetLoadBalancerMonitorOptions) SetMonitorIdentifier(monitorIdentifier string) *GetLoadBalancerMonitorOptions {
	options.MonitorIdentifier = core.StringPtr(monitorIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetLoadBalancerMonitorOptions) SetHeaders(param map[string]string) *GetLoadBalancerMonitorOptions {
	options.Headers = param
	return options
}

// ListAllLoadBalancerMonitorsOptions : The ListAllLoadBalancerMonitors options.
type ListAllLoadBalancerMonitorsOptions struct {

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListAllLoadBalancerMonitorsOptions : Instantiate ListAllLoadBalancerMonitorsOptions
func (*GlobalLoadBalancerMonitorV1) NewListAllLoadBalancerMonitorsOptions() *ListAllLoadBalancerMonitorsOptions {
	return &ListAllLoadBalancerMonitorsOptions{}
}

// SetHeaders : Allow user to set Headers
func (options *ListAllLoadBalancerMonitorsOptions) SetHeaders(param map[string]string) *ListAllLoadBalancerMonitorsOptions {
	options.Headers = param
	return options
}

// DeleteMonitorResp : delete monitor response object.
type DeleteMonitorResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result *DeleteMonitorRespResult `json:"result" validate:"required"`
}


// UnmarshalDeleteMonitorResp unmarshals an instance of DeleteMonitorResp from the specified map of raw messages.
func UnmarshalDeleteMonitorResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DeleteMonitorResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalDeleteMonitorRespResult)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ListMonitorResp : monitor list response.
type ListMonitorResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// result.
	Result []MonitorPack `json:"result" validate:"required"`

	// result information.
	ResultInfo *ResultInfo `json:"result_info" validate:"required"`
}


// UnmarshalListMonitorResp unmarshals an instance of ListMonitorResp from the specified map of raw messages.
func UnmarshalListMonitorResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ListMonitorResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalMonitorPack)
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

// MonitorPack : monitor package.
type MonitorPack struct {
	// identifier.
	ID *string `json:"id,omitempty"`

	// created date.
	CreatedOn *string `json:"created_on,omitempty"`

	// modified date.
	ModifiedOn *string `json:"modified_on,omitempty"`

	// type.
	Type *string `json:"type,omitempty"`

	// login page.
	Description *string `json:"description,omitempty"`

	// method name.
	Method *string `json:"method,omitempty"`

	// port number.
	Port *int64 `json:"port,omitempty"`

	// path.
	Path *string `json:"path,omitempty"`

	// timeout count.
	Timeout *int64 `json:"timeout,omitempty"`

	// retries count.
	Retries *int64 `json:"retries,omitempty"`

	// interval.
	Interval *int64 `json:"interval,omitempty"`

	// expected body.
	ExpectedBody *string `json:"expected_body" validate:"required"`

	// expected codes.
	ExpectedCodes *string `json:"expected_codes" validate:"required"`

	// follow redirects.
	FollowRedirects *bool `json:"follow_redirects,omitempty"`

	// allow insecure.
	AllowInsecure *bool `json:"allow_insecure,omitempty"`

	// header.
	Header map[string][]string `json:"header,omitempty"`
}


// UnmarshalMonitorPack unmarshals an instance of MonitorPack from the specified map of raw messages.
func UnmarshalMonitorPack(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MonitorPack)
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
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "description", &obj.Description)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "method", &obj.Method)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "port", &obj.Port)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "path", &obj.Path)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "retries", &obj.Retries)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "interval", &obj.Interval)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_body", &obj.ExpectedBody)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "expected_codes", &obj.ExpectedCodes)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "follow_redirects", &obj.FollowRedirects)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "allow_insecure", &obj.AllowInsecure)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "header", &obj.Header)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// MonitorResp : monitor response.
type MonitorResp struct {
	// success response.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// monitor package.
	Result *MonitorPack `json:"result" validate:"required"`
}


// UnmarshalMonitorResp unmarshals an instance of MonitorResp from the specified map of raw messages.
func UnmarshalMonitorResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(MonitorResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalMonitorPack)
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

	// per page number.
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
