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
 

// Package rangeapplicationsv1 : Operations and models for the RangeApplicationsV1 service
package rangeapplicationsv1

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/go-sdk-core/v4/core"
	common "github.com/IBM/networking-go-sdk/common"
	"github.com/go-openapi/strfmt"
	"net/http"
	"reflect"
	"time"
)

// RangeApplicationsV1 : Range Applications
//
// Version: 1.0.0
type RangeApplicationsV1 struct {
	Service *core.BaseService

	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string

	// zone identifier.
	ZoneIdentifier *string
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.cis.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "range_applications"

// RangeApplicationsV1Options : Service options
type RangeApplicationsV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator

	// Full url-encoded cloud resource name (CRN) of resource instance.
	Crn *string `validate:"required"`

	// zone identifier.
	ZoneIdentifier *string `validate:"required"`
}

// NewRangeApplicationsV1UsingExternalConfig : constructs an instance of RangeApplicationsV1 with passed in options and external configuration.
func NewRangeApplicationsV1UsingExternalConfig(options *RangeApplicationsV1Options) (rangeApplications *RangeApplicationsV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	rangeApplications, err = NewRangeApplicationsV1(options)
	if err != nil {
		return
	}

	err = rangeApplications.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = rangeApplications.Service.SetServiceURL(options.URL)
	}
	return
}

// NewRangeApplicationsV1 : constructs an instance of RangeApplicationsV1 with passed in options.
func NewRangeApplicationsV1(options *RangeApplicationsV1Options) (service *RangeApplicationsV1, err error) {
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

	service = &RangeApplicationsV1{
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

// Clone makes a copy of "rangeApplications" suitable for processing requests.
func (rangeApplications *RangeApplicationsV1) Clone() *RangeApplicationsV1 {
	if core.IsNil(rangeApplications) {
		return nil
	}
	clone := *rangeApplications
	clone.Service = rangeApplications.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (rangeApplications *RangeApplicationsV1) SetServiceURL(url string) error {
	return rangeApplications.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (rangeApplications *RangeApplicationsV1) GetServiceURL() string {
	return rangeApplications.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (rangeApplications *RangeApplicationsV1) SetDefaultHeaders(headers http.Header) {
	rangeApplications.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (rangeApplications *RangeApplicationsV1) SetEnableGzipCompression(enableGzip bool) {
	rangeApplications.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (rangeApplications *RangeApplicationsV1) GetEnableGzipCompression() bool {
	return rangeApplications.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (rangeApplications *RangeApplicationsV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	rangeApplications.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (rangeApplications *RangeApplicationsV1) DisableRetries() {
	rangeApplications.Service.DisableRetries()
}

// ListRangeApps : List range applications
// Get a list of currently existing Range Applications inside a zone.
func (rangeApplications *RangeApplicationsV1) ListRangeApps(listRangeAppsOptions *ListRangeAppsOptions) (result *RangeApplications, response *core.DetailedResponse, err error) {
	return rangeApplications.ListRangeAppsWithContext(context.Background(), listRangeAppsOptions)
}

// ListRangeAppsWithContext is an alternate form of the ListRangeApps method which supports a Context parameter
func (rangeApplications *RangeApplicationsV1) ListRangeAppsWithContext(ctx context.Context, listRangeAppsOptions *ListRangeAppsOptions) (result *RangeApplications, response *core.DetailedResponse, err error) {
	err = core.ValidateStruct(listRangeAppsOptions, "listRangeAppsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *rangeApplications.Crn,
		"zone_identifier": *rangeApplications.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rangeApplications.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rangeApplications.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/range/apps`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listRangeAppsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("range_applications", "V1", "ListRangeApps")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listRangeAppsOptions.Page != nil {
		builder.AddQuery("page", fmt.Sprint(*listRangeAppsOptions.Page))
	}
	if listRangeAppsOptions.PerPage != nil {
		builder.AddQuery("per_page", fmt.Sprint(*listRangeAppsOptions.PerPage))
	}
	if listRangeAppsOptions.Order != nil {
		builder.AddQuery("order", fmt.Sprint(*listRangeAppsOptions.Order))
	}
	if listRangeAppsOptions.Direction != nil {
		builder.AddQuery("direction", fmt.Sprint(*listRangeAppsOptions.Direction))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rangeApplications.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRangeApplications)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateRangeApp : Create Range Application
// Create a Range Applications inside a zone.
func (rangeApplications *RangeApplicationsV1) CreateRangeApp(createRangeAppOptions *CreateRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	return rangeApplications.CreateRangeAppWithContext(context.Background(), createRangeAppOptions)
}

// CreateRangeAppWithContext is an alternate form of the CreateRangeApp method which supports a Context parameter
func (rangeApplications *RangeApplicationsV1) CreateRangeAppWithContext(ctx context.Context, createRangeAppOptions *CreateRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createRangeAppOptions, "createRangeAppOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createRangeAppOptions, "createRangeAppOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *rangeApplications.Crn,
		"zone_identifier": *rangeApplications.ZoneIdentifier,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rangeApplications.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rangeApplications.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/range/apps`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createRangeAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("range_applications", "V1", "CreateRangeApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createRangeAppOptions.Protocol != nil {
		body["protocol"] = createRangeAppOptions.Protocol
	}
	if createRangeAppOptions.Dns != nil {
		body["dns"] = createRangeAppOptions.Dns
	}
	if createRangeAppOptions.OriginDirect != nil {
		body["origin_direct"] = createRangeAppOptions.OriginDirect
	}
	if createRangeAppOptions.OriginDns != nil {
		body["origin_dns"] = createRangeAppOptions.OriginDns
	}
	if createRangeAppOptions.OriginPort != nil {
		body["origin_port"] = createRangeAppOptions.OriginPort
	}
	if createRangeAppOptions.IpFirewall != nil {
		body["ip_firewall"] = createRangeAppOptions.IpFirewall
	}
	if createRangeAppOptions.ProxyProtocol != nil {
		body["proxy_protocol"] = createRangeAppOptions.ProxyProtocol
	}
	if createRangeAppOptions.EdgeIps != nil {
		body["edge_ips"] = createRangeAppOptions.EdgeIps
	}
	if createRangeAppOptions.TrafficType != nil {
		body["traffic_type"] = createRangeAppOptions.TrafficType
	}
	if createRangeAppOptions.Tls != nil {
		body["tls"] = createRangeAppOptions.Tls
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
	response, err = rangeApplications.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRangeApplicationResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetRangeApp : Get range application a zone
// Get the application configuration of a specific application inside a zone.
func (rangeApplications *RangeApplicationsV1) GetRangeApp(getRangeAppOptions *GetRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	return rangeApplications.GetRangeAppWithContext(context.Background(), getRangeAppOptions)
}

// GetRangeAppWithContext is an alternate form of the GetRangeApp method which supports a Context parameter
func (rangeApplications *RangeApplicationsV1) GetRangeAppWithContext(ctx context.Context, getRangeAppOptions *GetRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getRangeAppOptions, "getRangeAppOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getRangeAppOptions, "getRangeAppOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *rangeApplications.Crn,
		"zone_identifier": *rangeApplications.ZoneIdentifier,
		"app_identifier": *getRangeAppOptions.AppIdentifier,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rangeApplications.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rangeApplications.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/range/apps/{app_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getRangeAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("range_applications", "V1", "GetRangeApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rangeApplications.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRangeApplicationResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateRangeApp : Update range application
// Update a Range Application inside a zone.
func (rangeApplications *RangeApplicationsV1) UpdateRangeApp(updateRangeAppOptions *UpdateRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	return rangeApplications.UpdateRangeAppWithContext(context.Background(), updateRangeAppOptions)
}

// UpdateRangeAppWithContext is an alternate form of the UpdateRangeApp method which supports a Context parameter
func (rangeApplications *RangeApplicationsV1) UpdateRangeAppWithContext(ctx context.Context, updateRangeAppOptions *UpdateRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateRangeAppOptions, "updateRangeAppOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateRangeAppOptions, "updateRangeAppOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *rangeApplications.Crn,
		"zone_identifier": *rangeApplications.ZoneIdentifier,
		"app_identifier": *updateRangeAppOptions.AppIdentifier,
	}

	builder := core.NewRequestBuilder(core.PUT)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rangeApplications.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rangeApplications.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/range/apps/{app_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateRangeAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("range_applications", "V1", "UpdateRangeApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateRangeAppOptions.Protocol != nil {
		body["protocol"] = updateRangeAppOptions.Protocol
	}
	if updateRangeAppOptions.Dns != nil {
		body["dns"] = updateRangeAppOptions.Dns
	}
	if updateRangeAppOptions.OriginDirect != nil {
		body["origin_direct"] = updateRangeAppOptions.OriginDirect
	}
	if updateRangeAppOptions.OriginDns != nil {
		body["origin_dns"] = updateRangeAppOptions.OriginDns
	}
	if updateRangeAppOptions.OriginPort != nil {
		body["origin_port"] = updateRangeAppOptions.OriginPort
	}
	if updateRangeAppOptions.IpFirewall != nil {
		body["ip_firewall"] = updateRangeAppOptions.IpFirewall
	}
	if updateRangeAppOptions.ProxyProtocol != nil {
		body["proxy_protocol"] = updateRangeAppOptions.ProxyProtocol
	}
	if updateRangeAppOptions.EdgeIps != nil {
		body["edge_ips"] = updateRangeAppOptions.EdgeIps
	}
	if updateRangeAppOptions.TrafficType != nil {
		body["traffic_type"] = updateRangeAppOptions.TrafficType
	}
	if updateRangeAppOptions.Tls != nil {
		body["tls"] = updateRangeAppOptions.Tls
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
	response, err = rangeApplications.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRangeApplicationResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteRangeApp : Delete range application
// Delete a specific application configuration.
func (rangeApplications *RangeApplicationsV1) DeleteRangeApp(deleteRangeAppOptions *DeleteRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	return rangeApplications.DeleteRangeAppWithContext(context.Background(), deleteRangeAppOptions)
}

// DeleteRangeAppWithContext is an alternate form of the DeleteRangeApp method which supports a Context parameter
func (rangeApplications *RangeApplicationsV1) DeleteRangeAppWithContext(ctx context.Context, deleteRangeAppOptions *DeleteRangeAppOptions) (result *RangeApplicationResp, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteRangeAppOptions, "deleteRangeAppOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteRangeAppOptions, "deleteRangeAppOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"crn": *rangeApplications.Crn,
		"zone_identifier": *rangeApplications.ZoneIdentifier,
		"app_identifier": *deleteRangeAppOptions.AppIdentifier,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = rangeApplications.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(rangeApplications.Service.Options.URL, `/v1/{crn}/zones/{zone_identifier}/range/apps/{app_identifier}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteRangeAppOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("range_applications", "V1", "DeleteRangeApp")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = rangeApplications.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalRangeApplicationResp)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateRangeAppOptions : The CreateRangeApp options.
type CreateRangeAppOptions struct {
	// Defines the protocol and port for this application.
	Protocol *string `json:"protocol" validate:"required"`

	// Name and type of the DNS record for this application.
	Dns *RangeAppReqDns `json:"dns" validate:"required"`

	// IP address and port of the origin for this Range application. If configuring a load balancer, use 'origin_dns' and
	// 'origin_port'. This can not be combined with 'origin_dns' and 'origin_port'.
	OriginDirect []string `json:"origin_direct,omitempty"`

	// DNS record pointing to the origin for this Range application. This is used for configuring a load balancer. When
	// specifying an individual IP address, use 'origin_direct'. This requires 'origin_port' and can not be combined with
	// 'origin_direct'.
	OriginDns *RangeAppReqOriginDns `json:"origin_dns,omitempty"`

	// Port at the origin that listens to traffic from this Range application. Requires 'origin_dns' and can not be
	// combined with 'origin_direct'.
	OriginPort *int64 `json:"origin_port,omitempty"`

	// Enables the IP Firewall for this application. Only available for TCP applications.
	IpFirewall *bool `json:"ip_firewall,omitempty"`

	// Allows for the true client IP to be passed to the service.
	ProxyProtocol *string `json:"proxy_protocol,omitempty"`

	// Configures IP version for the hostname of this application. Default is {"type":"dynamic", "connectivity":"all"}.
	EdgeIps *RangeAppReqEdgeIps `json:"edge_ips,omitempty"`

	// Configure how traffic is handled at the edge. If set to "direct" traffic is passed through to the service. In the
	// case of "http" or "https" HTTP/s features at the edge are applied ot this traffic.
	TrafficType *string `json:"traffic_type,omitempty"`

	// Configure if and how TLS connections are terminated at the edge.
	Tls *string `json:"tls,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateRangeAppOptions.ProxyProtocol property.
// Allows for the true client IP to be passed to the service.
const (
	CreateRangeAppOptions_ProxyProtocol_Off = "off"
	CreateRangeAppOptions_ProxyProtocol_Simple = "simple"
	CreateRangeAppOptions_ProxyProtocol_V1 = "v1"
	CreateRangeAppOptions_ProxyProtocol_V2 = "v2"
)

// Constants associated with the CreateRangeAppOptions.TrafficType property.
// Configure how traffic is handled at the edge. If set to "direct" traffic is passed through to the service. In the
// case of "http" or "https" HTTP/s features at the edge are applied ot this traffic.
const (
	CreateRangeAppOptions_TrafficType_Direct = "direct"
	CreateRangeAppOptions_TrafficType_Http = "http"
	CreateRangeAppOptions_TrafficType_Https = "https"
)

// Constants associated with the CreateRangeAppOptions.Tls property.
// Configure if and how TLS connections are terminated at the edge.
const (
	CreateRangeAppOptions_Tls_Flexible = "flexible"
	CreateRangeAppOptions_Tls_Full = "full"
	CreateRangeAppOptions_Tls_Off = "off"
	CreateRangeAppOptions_Tls_Strict = "strict"
)

// NewCreateRangeAppOptions : Instantiate CreateRangeAppOptions
func (*RangeApplicationsV1) NewCreateRangeAppOptions(protocol string, dns *RangeAppReqDns) *CreateRangeAppOptions {
	return &CreateRangeAppOptions{
		Protocol: core.StringPtr(protocol),
		Dns: dns,
	}
}

// SetProtocol : Allow user to set Protocol
func (options *CreateRangeAppOptions) SetProtocol(protocol string) *CreateRangeAppOptions {
	options.Protocol = core.StringPtr(protocol)
	return options
}

// SetDns : Allow user to set Dns
func (options *CreateRangeAppOptions) SetDns(dns *RangeAppReqDns) *CreateRangeAppOptions {
	options.Dns = dns
	return options
}

// SetOriginDirect : Allow user to set OriginDirect
func (options *CreateRangeAppOptions) SetOriginDirect(originDirect []string) *CreateRangeAppOptions {
	options.OriginDirect = originDirect
	return options
}

// SetOriginDns : Allow user to set OriginDns
func (options *CreateRangeAppOptions) SetOriginDns(originDns *RangeAppReqOriginDns) *CreateRangeAppOptions {
	options.OriginDns = originDns
	return options
}

// SetOriginPort : Allow user to set OriginPort
func (options *CreateRangeAppOptions) SetOriginPort(originPort int64) *CreateRangeAppOptions {
	options.OriginPort = core.Int64Ptr(originPort)
	return options
}

// SetIpFirewall : Allow user to set IpFirewall
func (options *CreateRangeAppOptions) SetIpFirewall(ipFirewall bool) *CreateRangeAppOptions {
	options.IpFirewall = core.BoolPtr(ipFirewall)
	return options
}

// SetProxyProtocol : Allow user to set ProxyProtocol
func (options *CreateRangeAppOptions) SetProxyProtocol(proxyProtocol string) *CreateRangeAppOptions {
	options.ProxyProtocol = core.StringPtr(proxyProtocol)
	return options
}

// SetEdgeIps : Allow user to set EdgeIps
func (options *CreateRangeAppOptions) SetEdgeIps(edgeIps *RangeAppReqEdgeIps) *CreateRangeAppOptions {
	options.EdgeIps = edgeIps
	return options
}

// SetTrafficType : Allow user to set TrafficType
func (options *CreateRangeAppOptions) SetTrafficType(trafficType string) *CreateRangeAppOptions {
	options.TrafficType = core.StringPtr(trafficType)
	return options
}

// SetTls : Allow user to set Tls
func (options *CreateRangeAppOptions) SetTls(tls string) *CreateRangeAppOptions {
	options.Tls = core.StringPtr(tls)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateRangeAppOptions) SetHeaders(param map[string]string) *CreateRangeAppOptions {
	options.Headers = param
	return options
}

// DeleteRangeAppOptions : The DeleteRangeApp options.
type DeleteRangeAppOptions struct {
	// application identifier.
	AppIdentifier *string `json:"app_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteRangeAppOptions : Instantiate DeleteRangeAppOptions
func (*RangeApplicationsV1) NewDeleteRangeAppOptions(appIdentifier string) *DeleteRangeAppOptions {
	return &DeleteRangeAppOptions{
		AppIdentifier: core.StringPtr(appIdentifier),
	}
}

// SetAppIdentifier : Allow user to set AppIdentifier
func (options *DeleteRangeAppOptions) SetAppIdentifier(appIdentifier string) *DeleteRangeAppOptions {
	options.AppIdentifier = core.StringPtr(appIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteRangeAppOptions) SetHeaders(param map[string]string) *DeleteRangeAppOptions {
	options.Headers = param
	return options
}

// GetRangeAppOptions : The GetRangeApp options.
type GetRangeAppOptions struct {
	// application identifier.
	AppIdentifier *string `json:"app_identifier" validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetRangeAppOptions : Instantiate GetRangeAppOptions
func (*RangeApplicationsV1) NewGetRangeAppOptions(appIdentifier string) *GetRangeAppOptions {
	return &GetRangeAppOptions{
		AppIdentifier: core.StringPtr(appIdentifier),
	}
}

// SetAppIdentifier : Allow user to set AppIdentifier
func (options *GetRangeAppOptions) SetAppIdentifier(appIdentifier string) *GetRangeAppOptions {
	options.AppIdentifier = core.StringPtr(appIdentifier)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetRangeAppOptions) SetHeaders(param map[string]string) *GetRangeAppOptions {
	options.Headers = param
	return options
}

// ListRangeAppsOptions : The ListRangeApps options.
type ListRangeAppsOptions struct {
	// Page number of paginated results.
	Page *int64 `json:"page,omitempty"`

	// Maximum number of Range applications per page.
	PerPage *int64 `json:"per_page,omitempty"`

	// Field by which to order the list of Range applications.
	Order *string `json:"order,omitempty"`

	// Direction in which to order results [ascending/descending order].
	Direction *string `json:"direction,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListRangeAppsOptions.Order property.
// Field by which to order the list of Range applications.
const (
	ListRangeAppsOptions_Order_AppID = "app_id"
	ListRangeAppsOptions_Order_CreatedOn = "created_on"
	ListRangeAppsOptions_Order_Dns = "dns"
	ListRangeAppsOptions_Order_ModifiedOn = "modified_on"
	ListRangeAppsOptions_Order_Protocol = "protocol"
)

// Constants associated with the ListRangeAppsOptions.Direction property.
// Direction in which to order results [ascending/descending order].
const (
	ListRangeAppsOptions_Direction_Asc = "asc"
	ListRangeAppsOptions_Direction_Desc = "desc"
)

// NewListRangeAppsOptions : Instantiate ListRangeAppsOptions
func (*RangeApplicationsV1) NewListRangeAppsOptions() *ListRangeAppsOptions {
	return &ListRangeAppsOptions{}
}

// SetPage : Allow user to set Page
func (options *ListRangeAppsOptions) SetPage(page int64) *ListRangeAppsOptions {
	options.Page = core.Int64Ptr(page)
	return options
}

// SetPerPage : Allow user to set PerPage
func (options *ListRangeAppsOptions) SetPerPage(perPage int64) *ListRangeAppsOptions {
	options.PerPage = core.Int64Ptr(perPage)
	return options
}

// SetOrder : Allow user to set Order
func (options *ListRangeAppsOptions) SetOrder(order string) *ListRangeAppsOptions {
	options.Order = core.StringPtr(order)
	return options
}

// SetDirection : Allow user to set Direction
func (options *ListRangeAppsOptions) SetDirection(direction string) *ListRangeAppsOptions {
	options.Direction = core.StringPtr(direction)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListRangeAppsOptions) SetHeaders(param map[string]string) *ListRangeAppsOptions {
	options.Headers = param
	return options
}

// RangeAppReqDns : Name and type of the DNS record for this application.
type RangeAppReqDns struct {
	// DNS record type.
	Type *string `json:"type,omitempty"`

	// DNS record name.
	Name *string `json:"name,omitempty"`
}


// UnmarshalRangeAppReqDns unmarshals an instance of RangeAppReqDns from the specified map of raw messages.
func UnmarshalRangeAppReqDns(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeAppReqDns)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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

// RangeAppReqEdgeIps : Configures IP version for the hostname of this application. Default is {"type":"dynamic", "connectivity":"all"}.
type RangeAppReqEdgeIps struct {
	// The type of edge IP configuration.
	Type *string `json:"type,omitempty"`

	// Specifies the IP version (or all).
	Connectivity *string `json:"connectivity,omitempty"`
}

// Constants associated with the RangeAppReqEdgeIps.Type property.
// The type of edge IP configuration.
const (
	RangeAppReqEdgeIps_Type_Dynamic = "dynamic"
)

// Constants associated with the RangeAppReqEdgeIps.Connectivity property.
// Specifies the IP version (or all).
const (
	RangeAppReqEdgeIps_Connectivity_All = "all"
	RangeAppReqEdgeIps_Connectivity_Ipv4 = "ipv4"
	RangeAppReqEdgeIps_Connectivity_Ipv6 = "ipv6"
)


// UnmarshalRangeAppReqEdgeIps unmarshals an instance of RangeAppReqEdgeIps from the specified map of raw messages.
func UnmarshalRangeAppReqEdgeIps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeAppReqEdgeIps)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "connectivity", &obj.Connectivity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RangeAppReqOriginDns : DNS record pointing to the origin for this Range application. This is used for configuring a load balancer. When
// specifying an individual IP address, use 'origin_direct'. This requires 'origin_port' and can not be combined with
// 'origin_direct'.
type RangeAppReqOriginDns struct {
	// Name of the origin.
	Name *string `json:"name" validate:"required"`
}


// NewRangeAppReqOriginDns : Instantiate RangeAppReqOriginDns (Generic Model Constructor)
func (*RangeApplicationsV1) NewRangeAppReqOriginDns(name string) (model *RangeAppReqOriginDns, err error) {
	model = &RangeAppReqOriginDns{
		Name: core.StringPtr(name),
	}
	err = core.ValidateStruct(model, "required parameters")
	return
}

// UnmarshalRangeAppReqOriginDns unmarshals an instance of RangeAppReqOriginDns from the specified map of raw messages.
func UnmarshalRangeAppReqOriginDns(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeAppReqOriginDns)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RangeApplicationObjectDns : The name and type of DNS record for the Range application.
type RangeApplicationObjectDns struct {
	// The type of DNS record associated with the application.
	Type *string `json:"type,omitempty"`

	// The name of the DNS record associated with the application.
	Name *string `json:"name,omitempty"`
}

// Constants associated with the RangeApplicationObjectDns.Type property.
// The type of DNS record associated with the application.
const (
	RangeApplicationObjectDns_Type_Cname = "CNAME"
)


// UnmarshalRangeApplicationObjectDns unmarshals an instance of RangeApplicationObjectDns from the specified map of raw messages.
func UnmarshalRangeApplicationObjectDns(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeApplicationObjectDns)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
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

// RangeApplicationObjectEdgeIps : Configures IP version for the hostname of this application.
type RangeApplicationObjectEdgeIps struct {
	// The type of edge IP configuration.
	Type *string `json:"type,omitempty"`

	// Specifies the IP version (or all).
	Connectivity *string `json:"connectivity,omitempty"`
}

// Constants associated with the RangeApplicationObjectEdgeIps.Type property.
// The type of edge IP configuration.
const (
	RangeApplicationObjectEdgeIps_Type_Dynamic = "dynamic"
)

// Constants associated with the RangeApplicationObjectEdgeIps.Connectivity property.
// Specifies the IP version (or all).
const (
	RangeApplicationObjectEdgeIps_Connectivity_All = "all"
	RangeApplicationObjectEdgeIps_Connectivity_Ipv4 = "ipv4"
	RangeApplicationObjectEdgeIps_Connectivity_Ipv6 = "ipv6"
)


// UnmarshalRangeApplicationObjectEdgeIps unmarshals an instance of RangeApplicationObjectEdgeIps from the specified map of raw messages.
func UnmarshalRangeApplicationObjectEdgeIps(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeApplicationObjectEdgeIps)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "connectivity", &obj.Connectivity)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateRangeAppOptions : The UpdateRangeApp options.
type UpdateRangeAppOptions struct {
	// application identifier.
	AppIdentifier *string `json:"app_identifier" validate:"required,ne="`

	// Defines the protocol and port for this application.
	Protocol *string `json:"protocol" validate:"required"`

	// Name and type of the DNS record for this application.
	Dns *RangeAppReqDns `json:"dns" validate:"required"`

	// IP address and port of the origin for this Range application. If configuring a load balancer, use 'origin_dns' and
	// 'origin_port'. This can not be combined with 'origin_dns' and 'origin_port'.
	OriginDirect []string `json:"origin_direct,omitempty"`

	// DNS record pointing to the origin for this Range application. This is used for configuring a load balancer. When
	// specifying an individual IP address, use 'origin_direct'. This requires 'origin_port' and can not be combined with
	// 'origin_direct'.
	OriginDns *RangeAppReqOriginDns `json:"origin_dns,omitempty"`

	// Port at the origin that listens to traffic from this Range application. Requires 'origin_dns' and can not be
	// combined with 'origin_direct'.
	OriginPort *int64 `json:"origin_port,omitempty"`

	// Enables the IP Firewall for this application. Only available for TCP applications.
	IpFirewall *bool `json:"ip_firewall,omitempty"`

	// Allows for the true client IP to be passed to the service.
	ProxyProtocol *string `json:"proxy_protocol,omitempty"`

	// Configures IP version for the hostname of this application. Default is {"type":"dynamic", "connectivity":"all"}.
	EdgeIps *RangeAppReqEdgeIps `json:"edge_ips,omitempty"`

	// Configure how traffic is handled at the edge. If set to "direct" traffic is passed through to the service. In the
	// case of "http" or "https" HTTP/s features at the edge are applied ot this traffic.
	TrafficType *string `json:"traffic_type,omitempty"`

	// Configure if and how TLS connections are terminated at the edge.
	Tls *string `json:"tls,omitempty"`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateRangeAppOptions.ProxyProtocol property.
// Allows for the true client IP to be passed to the service.
const (
	UpdateRangeAppOptions_ProxyProtocol_Off = "off"
	UpdateRangeAppOptions_ProxyProtocol_Simple = "simple"
	UpdateRangeAppOptions_ProxyProtocol_V1 = "v1"
	UpdateRangeAppOptions_ProxyProtocol_V2 = "v2"
)

// Constants associated with the UpdateRangeAppOptions.TrafficType property.
// Configure how traffic is handled at the edge. If set to "direct" traffic is passed through to the service. In the
// case of "http" or "https" HTTP/s features at the edge are applied ot this traffic.
const (
	UpdateRangeAppOptions_TrafficType_Direct = "direct"
	UpdateRangeAppOptions_TrafficType_Http = "http"
	UpdateRangeAppOptions_TrafficType_Https = "https"
)

// Constants associated with the UpdateRangeAppOptions.Tls property.
// Configure if and how TLS connections are terminated at the edge.
const (
	UpdateRangeAppOptions_Tls_Flexible = "flexible"
	UpdateRangeAppOptions_Tls_Full = "full"
	UpdateRangeAppOptions_Tls_Off = "off"
	UpdateRangeAppOptions_Tls_Strict = "strict"
)

// NewUpdateRangeAppOptions : Instantiate UpdateRangeAppOptions
func (*RangeApplicationsV1) NewUpdateRangeAppOptions(appIdentifier string, protocol string, dns *RangeAppReqDns) *UpdateRangeAppOptions {
	return &UpdateRangeAppOptions{
		AppIdentifier: core.StringPtr(appIdentifier),
		Protocol: core.StringPtr(protocol),
		Dns: dns,
	}
}

// SetAppIdentifier : Allow user to set AppIdentifier
func (options *UpdateRangeAppOptions) SetAppIdentifier(appIdentifier string) *UpdateRangeAppOptions {
	options.AppIdentifier = core.StringPtr(appIdentifier)
	return options
}

// SetProtocol : Allow user to set Protocol
func (options *UpdateRangeAppOptions) SetProtocol(protocol string) *UpdateRangeAppOptions {
	options.Protocol = core.StringPtr(protocol)
	return options
}

// SetDns : Allow user to set Dns
func (options *UpdateRangeAppOptions) SetDns(dns *RangeAppReqDns) *UpdateRangeAppOptions {
	options.Dns = dns
	return options
}

// SetOriginDirect : Allow user to set OriginDirect
func (options *UpdateRangeAppOptions) SetOriginDirect(originDirect []string) *UpdateRangeAppOptions {
	options.OriginDirect = originDirect
	return options
}

// SetOriginDns : Allow user to set OriginDns
func (options *UpdateRangeAppOptions) SetOriginDns(originDns *RangeAppReqOriginDns) *UpdateRangeAppOptions {
	options.OriginDns = originDns
	return options
}

// SetOriginPort : Allow user to set OriginPort
func (options *UpdateRangeAppOptions) SetOriginPort(originPort int64) *UpdateRangeAppOptions {
	options.OriginPort = core.Int64Ptr(originPort)
	return options
}

// SetIpFirewall : Allow user to set IpFirewall
func (options *UpdateRangeAppOptions) SetIpFirewall(ipFirewall bool) *UpdateRangeAppOptions {
	options.IpFirewall = core.BoolPtr(ipFirewall)
	return options
}

// SetProxyProtocol : Allow user to set ProxyProtocol
func (options *UpdateRangeAppOptions) SetProxyProtocol(proxyProtocol string) *UpdateRangeAppOptions {
	options.ProxyProtocol = core.StringPtr(proxyProtocol)
	return options
}

// SetEdgeIps : Allow user to set EdgeIps
func (options *UpdateRangeAppOptions) SetEdgeIps(edgeIps *RangeAppReqEdgeIps) *UpdateRangeAppOptions {
	options.EdgeIps = edgeIps
	return options
}

// SetTrafficType : Allow user to set TrafficType
func (options *UpdateRangeAppOptions) SetTrafficType(trafficType string) *UpdateRangeAppOptions {
	options.TrafficType = core.StringPtr(trafficType)
	return options
}

// SetTls : Allow user to set Tls
func (options *UpdateRangeAppOptions) SetTls(tls string) *UpdateRangeAppOptions {
	options.Tls = core.StringPtr(tls)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateRangeAppOptions) SetHeaders(param map[string]string) *UpdateRangeAppOptions {
	options.Headers = param
	return options
}

// RangeApplicationObject : range application object.
type RangeApplicationObject struct {
	// Application identifier.
	ID *string `json:"id,omitempty"`

	// Port configuration.
	Protocol *string `json:"protocol,omitempty"`

	// The name and type of DNS record for the Range application.
	Dns *RangeApplicationObjectDns `json:"dns,omitempty"`

	// A list of destination addresses to the origin.
	OriginDirect []string `json:"origin_direct,omitempty"`

	// Enables the IP Firewall for this application.
	IpFirewall *bool `json:"ip_firewall,omitempty"`

	// Allows for the true client IP to be passed to the service.
	ProxyProtocol *string `json:"proxy_protocol,omitempty"`

	// Configures IP version for the hostname of this application.
	EdgeIps *RangeApplicationObjectEdgeIps `json:"edge_ips,omitempty"`

	// Specifies the TLS termination at the edge.
	Tls *string `json:"tls,omitempty"`

	// Configure how traffic is handled at the edge. If set to "direct" traffic is passed through to the service. In the
	// case of "http" or "https" HTTP/s features at the edge are applied ot this traffic.
	TrafficType *string `json:"traffic_type,omitempty"`

	// When the Application was created.
	CreatedOn *strfmt.DateTime `json:"created_on,omitempty"`

	// When the Application was last modified.
	ModifiedOn *strfmt.DateTime `json:"modified_on,omitempty"`
}

// Constants associated with the RangeApplicationObject.ProxyProtocol property.
// Allows for the true client IP to be passed to the service.
const (
	RangeApplicationObject_ProxyProtocol_Off = "off"
	RangeApplicationObject_ProxyProtocol_Simple = "simple"
	RangeApplicationObject_ProxyProtocol_V1 = "v1"
	RangeApplicationObject_ProxyProtocol_V2 = "v2"
)

// Constants associated with the RangeApplicationObject.Tls property.
// Specifies the TLS termination at the edge.
const (
	RangeApplicationObject_Tls_Flexible = "flexible"
	RangeApplicationObject_Tls_Full = "full"
	RangeApplicationObject_Tls_Off = "off"
	RangeApplicationObject_Tls_Strict = "strict"
)

// Constants associated with the RangeApplicationObject.TrafficType property.
// Configure how traffic is handled at the edge. If set to "direct" traffic is passed through to the service. In the
// case of "http" or "https" HTTP/s features at the edge are applied ot this traffic.
const (
	RangeApplicationObject_TrafficType_Direct = "direct"
	RangeApplicationObject_TrafficType_Http = "http"
	RangeApplicationObject_TrafficType_Https = "https"
)


// UnmarshalRangeApplicationObject unmarshals an instance of RangeApplicationObject from the specified map of raw messages.
func UnmarshalRangeApplicationObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeApplicationObject)
	err = core.UnmarshalPrimitive(m, "id", &obj.ID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "protocol", &obj.Protocol)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "dns", &obj.Dns, UnmarshalRangeApplicationObjectDns)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "origin_direct", &obj.OriginDirect)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "ip_firewall", &obj.IpFirewall)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "proxy_protocol", &obj.ProxyProtocol)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "edge_ips", &obj.EdgeIps, UnmarshalRangeApplicationObjectEdgeIps)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tls", &obj.Tls)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "traffic_type", &obj.TrafficType)
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
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RangeApplicationResp : range application response.
type RangeApplicationResp struct {
	// Was the get successful.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// range application object.
	Result *RangeApplicationObject `json:"result" validate:"required"`
}


// UnmarshalRangeApplicationResp unmarshals an instance of RangeApplicationResp from the specified map of raw messages.
func UnmarshalRangeApplicationResp(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeApplicationResp)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRangeApplicationObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// RangeApplications : range application.
type RangeApplications struct {
	// Was the get successful.
	Success *bool `json:"success" validate:"required"`

	// errors.
	Errors [][]string `json:"errors" validate:"required"`

	// messages.
	Messages [][]string `json:"messages" validate:"required"`

	// Container for Range application objects.
	Result []RangeApplicationObject `json:"result" validate:"required"`
}


// UnmarshalRangeApplications unmarshals an instance of RangeApplications from the specified map of raw messages.
func UnmarshalRangeApplications(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(RangeApplications)
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
	err = core.UnmarshalModel(m, "result", &obj.Result, UnmarshalRangeApplicationObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}
