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
 * IBM OpenAPI SDK Code Generator Version: 3.28.0-55613c9e-20210220-164656
 */

// Package satellitelinkv1 : Operations and models for the SatelliteLinkV1 service
package satellitelinkv1

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	common "github.com/IBM-Cloud/container-services-go-sdk/common"
	"github.com/IBM/go-sdk-core/v5/core"
)

// SatelliteLinkV1 : Satellite Link is a component of IBM Cloud Satellite, it provides Secure Connections between
// Locations and IBM Cloud.
//
// Version: 1.0.0
type SatelliteLinkV1 struct {
	Service *core.BaseService
}

// DefaultServiceURL is the default URL to make service requests to.
const DefaultServiceURL = "https://api.link.satellite.cloud.ibm.com"

// DefaultServiceName is the default key used to find external configuration information.
const DefaultServiceName = "satellite_link"

// SatelliteLinkV1Options : Service options
type SatelliteLinkV1Options struct {
	ServiceName   string
	URL           string
	Authenticator core.Authenticator
}

// NewSatelliteLinkV1UsingExternalConfig : constructs an instance of SatelliteLinkV1 with passed in options and external configuration.
func NewSatelliteLinkV1UsingExternalConfig(options *SatelliteLinkV1Options) (satelliteLink *SatelliteLinkV1, err error) {
	if options.ServiceName == "" {
		options.ServiceName = DefaultServiceName
	}

	if options.Authenticator == nil {
		options.Authenticator, err = core.GetAuthenticatorFromEnvironment(options.ServiceName)
		if err != nil {
			return
		}
	}

	satelliteLink, err = NewSatelliteLinkV1(options)
	if err != nil {
		return
	}

	err = satelliteLink.Service.ConfigureService(options.ServiceName)
	if err != nil {
		return
	}

	if options.URL != "" {
		err = satelliteLink.Service.SetServiceURL(options.URL)
	}
	return
}

// NewSatelliteLinkV1 : constructs an instance of SatelliteLinkV1 with passed in options.
func NewSatelliteLinkV1(options *SatelliteLinkV1Options) (service *SatelliteLinkV1, err error) {
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

	service = &SatelliteLinkV1{
		Service: baseService,
	}

	return
}

// GetServiceURLForRegion returns the service URL to be used for the specified region
func GetServiceURLForRegion(region string) (string, error) {
	return "", fmt.Errorf("service does not support regional URLs")
}

// Clone makes a copy of "satelliteLink" suitable for processing requests.
func (satelliteLink *SatelliteLinkV1) Clone() *SatelliteLinkV1 {
	if core.IsNil(satelliteLink) {
		return nil
	}
	clone := *satelliteLink
	clone.Service = satelliteLink.Service.Clone()
	return &clone
}

// SetServiceURL sets the service URL
func (satelliteLink *SatelliteLinkV1) SetServiceURL(url string) error {
	return satelliteLink.Service.SetServiceURL(url)
}

// GetServiceURL returns the service URL
func (satelliteLink *SatelliteLinkV1) GetServiceURL() string {
	return satelliteLink.Service.GetServiceURL()
}

// SetDefaultHeaders sets HTTP headers to be sent in every request
func (satelliteLink *SatelliteLinkV1) SetDefaultHeaders(headers http.Header) {
	satelliteLink.Service.SetDefaultHeaders(headers)
}

// SetEnableGzipCompression sets the service's EnableGzipCompression field
func (satelliteLink *SatelliteLinkV1) SetEnableGzipCompression(enableGzip bool) {
	satelliteLink.Service.SetEnableGzipCompression(enableGzip)
}

// GetEnableGzipCompression returns the service's EnableGzipCompression field
func (satelliteLink *SatelliteLinkV1) GetEnableGzipCompression() bool {
	return satelliteLink.Service.GetEnableGzipCompression()
}

// EnableRetries enables automatic retries for requests invoked for this service instance.
// If either parameter is specified as 0, then a default value is used instead.
func (satelliteLink *SatelliteLinkV1) EnableRetries(maxRetries int, maxRetryInterval time.Duration) {
	satelliteLink.Service.EnableRetries(maxRetries, maxRetryInterval)
}

// DisableRetries disables automatic retries for requests invoked for this service instance.
func (satelliteLink *SatelliteLinkV1) DisableRetries() {
	satelliteLink.Service.DisableRetries()
}

// CreateLink : create link [Administrator]
// Create Link for a Location.
func (satelliteLink *SatelliteLinkV1) CreateLink(createLinkOptions *CreateLinkOptions) (result *Location, response *core.DetailedResponse, err error) {
	return satelliteLink.CreateLinkWithContext(context.Background(), createLinkOptions)
}

// CreateLinkWithContext is an alternate form of the CreateLink method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) CreateLinkWithContext(ctx context.Context, createLinkOptions *CreateLinkOptions) (result *Location, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createLinkOptions, "createLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createLinkOptions, "createLinkOptions")
	if err != nil {
		return
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations`, nil)
	if err != nil {
		return
	}

	for headerName, headerValue := range createLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "CreateLink")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createLinkOptions.Crn != nil {
		body["crn"] = createLinkOptions.Crn
	}
	if createLinkOptions.LocationID != nil {
		body["location_id"] = createLinkOptions.LocationID
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocation)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetLink : read link [Administrator, Editor, Operator, Viewer, satellite-link-administrator]
// Retrieve Link information of a Location.
func (satelliteLink *SatelliteLinkV1) GetLink(getLinkOptions *GetLinkOptions) (result *Location, response *core.DetailedResponse, err error) {
	return satelliteLink.GetLinkWithContext(context.Background(), getLinkOptions)
}

// GetLinkWithContext is an alternate form of the GetLink method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) GetLinkWithContext(ctx context.Context, getLinkOptions *GetLinkOptions) (result *Location, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getLinkOptions, "getLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getLinkOptions, "getLinkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *getLinkOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "GetLink")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocation)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateLink : update link [Administrator, Editor, Operator, satellite-link-administrator]
// Update Link information of a Location.
func (satelliteLink *SatelliteLinkV1) UpdateLink(updateLinkOptions *UpdateLinkOptions) (result *Location, response *core.DetailedResponse, err error) {
	return satelliteLink.UpdateLinkWithContext(context.Background(), updateLinkOptions)
}

// UpdateLinkWithContext is an alternate form of the UpdateLink method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) UpdateLinkWithContext(ctx context.Context, updateLinkOptions *UpdateLinkOptions) (result *Location, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateLinkOptions, "updateLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateLinkOptions, "updateLinkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *updateLinkOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "UpdateLink")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateLinkOptions.WsEndpoint != nil {
		body["ws_endpoint"] = updateLinkOptions.WsEndpoint
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalLocation)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteLink : delete link [Administrator, Operator]
// Delete Link of a Location.
func (satelliteLink *SatelliteLinkV1) DeleteLink(deleteLinkOptions *DeleteLinkOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	return satelliteLink.DeleteLinkWithContext(context.Background(), deleteLinkOptions)
}

// DeleteLinkWithContext is an alternate form of the DeleteLink method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) DeleteLinkWithContext(ctx context.Context, deleteLinkOptions *DeleteLinkOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteLinkOptions, "deleteLinkOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteLinkOptions, "deleteLinkOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *deleteLinkOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteLinkOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "DeleteLink")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExecutionResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListEndpoints : list endpoints [Administrator, Editor, Viewer, Operator, satellite-link-administrator]
// List Endpoints of a Location.
func (satelliteLink *SatelliteLinkV1) ListEndpoints(listEndpointsOptions *ListEndpointsOptions) (result *Endpoints, response *core.DetailedResponse, err error) {
	return satelliteLink.ListEndpointsWithContext(context.Background(), listEndpointsOptions)
}

// ListEndpointsWithContext is an alternate form of the ListEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) ListEndpointsWithContext(ctx context.Context, listEndpointsOptions *ListEndpointsOptions) (result *Endpoints, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listEndpointsOptions, "listEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listEndpointsOptions, "listEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *listEndpointsOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "ListEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listEndpointsOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listEndpointsOptions.Type))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoints)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateEndpoints : create endpoint [Administrator, Editor, Operator, satellite-link-administrator]
// Create an endpoint.
func (satelliteLink *SatelliteLinkV1) CreateEndpoints(createEndpointsOptions *CreateEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	return satelliteLink.CreateEndpointsWithContext(context.Background(), createEndpointsOptions)
}

// CreateEndpointsWithContext is an alternate form of the CreateEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) CreateEndpointsWithContext(ctx context.Context, createEndpointsOptions *CreateEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createEndpointsOptions, "createEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createEndpointsOptions, "createEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *createEndpointsOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "CreateEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createEndpointsOptions.ConnType != nil {
		body["conn_type"] = createEndpointsOptions.ConnType
	}
	if createEndpointsOptions.DisplayName != nil {
		body["display_name"] = createEndpointsOptions.DisplayName
	}
	if createEndpointsOptions.ServerHost != nil {
		body["server_host"] = createEndpointsOptions.ServerHost
	}
	if createEndpointsOptions.ServerPort != nil {
		body["server_port"] = createEndpointsOptions.ServerPort
	}
	if createEndpointsOptions.Sni != nil {
		body["sni"] = createEndpointsOptions.Sni
	}
	if createEndpointsOptions.ClientProtocol != nil {
		body["client_protocol"] = createEndpointsOptions.ClientProtocol
	}
	if createEndpointsOptions.ClientMutualAuth != nil {
		body["client_mutual_auth"] = createEndpointsOptions.ClientMutualAuth
	}
	if createEndpointsOptions.ServerProtocol != nil {
		body["server_protocol"] = createEndpointsOptions.ServerProtocol
	}
	if createEndpointsOptions.ServerMutualAuth != nil {
		body["server_mutual_auth"] = createEndpointsOptions.ServerMutualAuth
	}
	if createEndpointsOptions.RejectUnauth != nil {
		body["reject_unauth"] = createEndpointsOptions.RejectUnauth
	}
	if createEndpointsOptions.Timeout != nil {
		body["timeout"] = createEndpointsOptions.Timeout
	}
	if createEndpointsOptions.CreatedBy != nil {
		body["created_by"] = createEndpointsOptions.CreatedBy
	}
	if createEndpointsOptions.Certs != nil {
		body["certs"] = createEndpointsOptions.Certs
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoint)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ImportEndpoints : import an endpoint [Administrator, Editor, Operator, satellite-link-administrator]
// Import an endpoint.
func (satelliteLink *SatelliteLinkV1) ImportEndpoints(importEndpointsOptions *ImportEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	return satelliteLink.ImportEndpointsWithContext(context.Background(), importEndpointsOptions)
}

// ImportEndpointsWithContext is an alternate form of the ImportEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) ImportEndpointsWithContext(ctx context.Context, importEndpointsOptions *ImportEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(importEndpointsOptions, "importEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(importEndpointsOptions, "importEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *importEndpointsOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/import`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range importEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "ImportEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	builder.AddFormData("state", "filename",
		core.StringNilMapper(importEndpointsOptions.StateContentType), importEndpointsOptions.State)

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoint)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ExportEndpoints : export an endpoint [Administrator, Editor, Operator, satellite-link-administrator]
// Export the endpoint.
func (satelliteLink *SatelliteLinkV1) ExportEndpoints(exportEndpointsOptions *ExportEndpointsOptions) (result *ExportEndpointsResponse, response *core.DetailedResponse, err error) {
	return satelliteLink.ExportEndpointsWithContext(context.Background(), exportEndpointsOptions)
}

// ExportEndpointsWithContext is an alternate form of the ExportEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) ExportEndpointsWithContext(ctx context.Context, exportEndpointsOptions *ExportEndpointsOptions) (result *ExportEndpointsResponse, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(exportEndpointsOptions, "exportEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(exportEndpointsOptions, "exportEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *exportEndpointsOptions.LocationID,
		"endpoint_id": *exportEndpointsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}/export`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range exportEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "ExportEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExportEndpointsResponse)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetEndpoints : read endpoint [Administrator, Editor, Viewer, Operator, satellite-link-administrator]
// Get the Endpoint's information.
func (satelliteLink *SatelliteLinkV1) GetEndpoints(getEndpointsOptions *GetEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	return satelliteLink.GetEndpointsWithContext(context.Background(), getEndpointsOptions)
}

// GetEndpointsWithContext is an alternate form of the GetEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) GetEndpointsWithContext(ctx context.Context, getEndpointsOptions *GetEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEndpointsOptions, "getEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEndpointsOptions, "getEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *getEndpointsOptions.LocationID,
		"endpoint_id": *getEndpointsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "GetEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoint)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateEndpoints : update endpoint [Administrator, Editor, Operator, satellite-link-administrator]
// Update the endpoint.
func (satelliteLink *SatelliteLinkV1) UpdateEndpoints(updateEndpointsOptions *UpdateEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	return satelliteLink.UpdateEndpointsWithContext(context.Background(), updateEndpointsOptions)
}

// UpdateEndpointsWithContext is an alternate form of the UpdateEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) UpdateEndpointsWithContext(ctx context.Context, updateEndpointsOptions *UpdateEndpointsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEndpointsOptions, "updateEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateEndpointsOptions, "updateEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *updateEndpointsOptions.LocationID,
		"endpoint_id": *updateEndpointsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "UpdateEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateEndpointsOptions.DisplayName != nil {
		body["display_name"] = updateEndpointsOptions.DisplayName
	}
	if updateEndpointsOptions.ServerHost != nil {
		body["server_host"] = updateEndpointsOptions.ServerHost
	}
	if updateEndpointsOptions.ServerPort != nil {
		body["server_port"] = updateEndpointsOptions.ServerPort
	}
	if updateEndpointsOptions.Sni != nil {
		body["sni"] = updateEndpointsOptions.Sni
	}
	if updateEndpointsOptions.ClientProtocol != nil {
		body["client_protocol"] = updateEndpointsOptions.ClientProtocol
	}
	if updateEndpointsOptions.ClientMutualAuth != nil {
		body["client_mutual_auth"] = updateEndpointsOptions.ClientMutualAuth
	}
	if updateEndpointsOptions.ServerProtocol != nil {
		body["server_protocol"] = updateEndpointsOptions.ServerProtocol
	}
	if updateEndpointsOptions.ServerMutualAuth != nil {
		body["server_mutual_auth"] = updateEndpointsOptions.ServerMutualAuth
	}
	if updateEndpointsOptions.RejectUnauth != nil {
		body["reject_unauth"] = updateEndpointsOptions.RejectUnauth
	}
	if updateEndpointsOptions.Timeout != nil {
		body["timeout"] = updateEndpointsOptions.Timeout
	}
	if updateEndpointsOptions.CreatedBy != nil {
		body["created_by"] = updateEndpointsOptions.CreatedBy
	}
	if updateEndpointsOptions.Certs != nil {
		body["certs"] = updateEndpointsOptions.Certs
	}
	if updateEndpointsOptions.Enabled != nil {
		body["enabled"] = updateEndpointsOptions.Enabled
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoint)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteEndpoints : delete endpoint [Administrator, Editor, Operator, satellite-link-administrator]
// Delete a endpoint.
func (satelliteLink *SatelliteLinkV1) DeleteEndpoints(deleteEndpointsOptions *DeleteEndpointsOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	return satelliteLink.DeleteEndpointsWithContext(context.Background(), deleteEndpointsOptions)
}

// DeleteEndpointsWithContext is an alternate form of the DeleteEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) DeleteEndpointsWithContext(ctx context.Context, deleteEndpointsOptions *DeleteEndpointsOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEndpointsOptions, "deleteEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEndpointsOptions, "deleteEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *deleteEndpointsOptions.LocationID,
		"endpoint_id": *deleteEndpointsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "DeleteEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExecutionResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// GetEndpointCerts : read endpoint certs [Administrator, Editor, Operator, satellite-link-administrator]
// Download certs for the endpoint.
func (satelliteLink *SatelliteLinkV1) GetEndpointCerts(getEndpointCertsOptions *GetEndpointCertsOptions) (result *DownloadedCerts, response *core.DetailedResponse, err error) {
	return satelliteLink.GetEndpointCertsWithContext(context.Background(), getEndpointCertsOptions)
}

// GetEndpointCertsWithContext is an alternate form of the GetEndpointCerts method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) GetEndpointCertsWithContext(ctx context.Context, getEndpointCertsOptions *GetEndpointCertsOptions) (result *DownloadedCerts, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(getEndpointCertsOptions, "getEndpointCertsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(getEndpointCertsOptions, "getEndpointCertsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *getEndpointCertsOptions.LocationID,
		"endpoint_id": *getEndpointCertsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}/cert`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range getEndpointCertsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "GetEndpointCerts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if getEndpointCertsOptions.NoZip != nil {
		builder.AddQuery("no_zip", fmt.Sprint(*getEndpointCertsOptions.NoZip))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalDownloadedCerts)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UploadEndpointCerts : upload endpoint certs [Administrator, Editor, Operator, satellite-link-administrator]
// Upload a cert for the endpoint.
func (satelliteLink *SatelliteLinkV1) UploadEndpointCerts(uploadEndpointCertsOptions *UploadEndpointCertsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	return satelliteLink.UploadEndpointCertsWithContext(context.Background(), uploadEndpointCertsOptions)
}

// UploadEndpointCertsWithContext is an alternate form of the UploadEndpointCerts method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) UploadEndpointCertsWithContext(ctx context.Context, uploadEndpointCertsOptions *UploadEndpointCertsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(uploadEndpointCertsOptions, "uploadEndpointCertsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(uploadEndpointCertsOptions, "uploadEndpointCertsOptions")
	if err != nil {
		return
	}
	if (uploadEndpointCertsOptions.ClientCert == nil) && (uploadEndpointCertsOptions.ServerCert == nil) && (uploadEndpointCertsOptions.ConnectorCert == nil) && (uploadEndpointCertsOptions.ConnectorKey == nil) {
		err = fmt.Errorf("at least one of clientCert, serverCert, connectorCert, or connectorKey must be supplied")
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *uploadEndpointCertsOptions.LocationID,
		"endpoint_id": *uploadEndpointCertsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}/cert`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range uploadEndpointCertsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "UploadEndpointCerts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if uploadEndpointCertsOptions.ClientCert != nil {
		builder.AddFormData("client_cert", "filename",
			core.StringNilMapper(uploadEndpointCertsOptions.ClientCertContentType), uploadEndpointCertsOptions.ClientCert)
	}
	if uploadEndpointCertsOptions.ServerCert != nil {
		builder.AddFormData("server_cert", "filename",
			core.StringNilMapper(uploadEndpointCertsOptions.ServerCertContentType), uploadEndpointCertsOptions.ServerCert)
	}
	if uploadEndpointCertsOptions.ConnectorCert != nil {
		builder.AddFormData("connector_cert", "filename",
			core.StringNilMapper(uploadEndpointCertsOptions.ConnectorCertContentType), uploadEndpointCertsOptions.ConnectorCert)
	}
	if uploadEndpointCertsOptions.ConnectorKey != nil {
		builder.AddFormData("connector_key", "filename",
			core.StringNilMapper(uploadEndpointCertsOptions.ConnectorKeyContentType), uploadEndpointCertsOptions.ConnectorKey)
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoint)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteEndpointCerts : delete endpoint certs [Administrator, Editor, Operator, satellite-link-administrator]
// Delete certs for the endpoint.
func (satelliteLink *SatelliteLinkV1) DeleteEndpointCerts(deleteEndpointCertsOptions *DeleteEndpointCertsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	return satelliteLink.DeleteEndpointCertsWithContext(context.Background(), deleteEndpointCertsOptions)
}

// DeleteEndpointCertsWithContext is an alternate form of the DeleteEndpointCerts method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) DeleteEndpointCertsWithContext(ctx context.Context, deleteEndpointCertsOptions *DeleteEndpointCertsOptions) (result *Endpoint, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteEndpointCertsOptions, "deleteEndpointCertsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteEndpointCertsOptions, "deleteEndpointCertsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *deleteEndpointCertsOptions.LocationID,
		"endpoint_id": *deleteEndpointCertsOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}/cert`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteEndpointCertsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "DeleteEndpointCerts")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpoint)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListEndpointSources : list endpoint sources [Administrator, Editor, Viewer, Operator, satellite-link-administrator]
// Get a list of Sources and operational status associated with given Endpoint.
func (satelliteLink *SatelliteLinkV1) ListEndpointSources(listEndpointSourcesOptions *ListEndpointSourcesOptions) (result *SourceStatus, response *core.DetailedResponse, err error) {
	return satelliteLink.ListEndpointSourcesWithContext(context.Background(), listEndpointSourcesOptions)
}

// ListEndpointSourcesWithContext is an alternate form of the ListEndpointSources method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) ListEndpointSourcesWithContext(ctx context.Context, listEndpointSourcesOptions *ListEndpointSourcesOptions) (result *SourceStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listEndpointSourcesOptions, "listEndpointSourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listEndpointSourcesOptions, "listEndpointSourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *listEndpointSourcesOptions.LocationID,
		"endpoint_id": *listEndpointSourcesOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}/sources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listEndpointSourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "ListEndpointSources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSourceStatus)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateEndpointSources : update endpoint sources [Administrator, Editor, Operator, satellite-link-source-access-controller]
// Update Source status associated with given Endpoint. If one of the status failed to update, all status will not be
// processed as well. The response will contain all source status which has been configured to this endpoint. For
// on-location endpoint, only service sources can be configured. For on-cloud endpoint, only user sources can be
// configured.
func (satelliteLink *SatelliteLinkV1) UpdateEndpointSources(updateEndpointSourcesOptions *UpdateEndpointSourcesOptions) (result *SourceStatus, response *core.DetailedResponse, err error) {
	return satelliteLink.UpdateEndpointSourcesWithContext(context.Background(), updateEndpointSourcesOptions)
}

// UpdateEndpointSourcesWithContext is an alternate form of the UpdateEndpointSources method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) UpdateEndpointSourcesWithContext(ctx context.Context, updateEndpointSourcesOptions *UpdateEndpointSourcesOptions) (result *SourceStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateEndpointSourcesOptions, "updateEndpointSourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateEndpointSourcesOptions, "updateEndpointSourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *updateEndpointSourcesOptions.LocationID,
		"endpoint_id": *updateEndpointSourcesOptions.EndpointID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/endpoints/{endpoint_id}/sources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateEndpointSourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "UpdateEndpointSources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateEndpointSourcesOptions.Sources != nil {
		body["sources"] = updateEndpointSourcesOptions.Sources
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSourceStatus)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListSources : list sources [Administrator, Editor, Viewer, Operator, satellite-link-administrator]
// List Sources associated with a Location.
func (satelliteLink *SatelliteLinkV1) ListSources(listSourcesOptions *ListSourcesOptions) (result *Sources, response *core.DetailedResponse, err error) {
	return satelliteLink.ListSourcesWithContext(context.Background(), listSourcesOptions)
}

// ListSourcesWithContext is an alternate form of the ListSources method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) ListSourcesWithContext(ctx context.Context, listSourcesOptions *ListSourcesOptions) (result *Sources, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSourcesOptions, "listSourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSourcesOptions, "listSourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *listSourcesOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/sources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "ListSources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	if listSourcesOptions.Type != nil {
		builder.AddQuery("type", fmt.Sprint(*listSourcesOptions.Type))
	}

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSources)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// CreateSources : create source [Administrator, Editor, Operator, satellite-link-administrator]
// Create a Source associated with given Location.
func (satelliteLink *SatelliteLinkV1) CreateSources(createSourcesOptions *CreateSourcesOptions) (result *Source, response *core.DetailedResponse, err error) {
	return satelliteLink.CreateSourcesWithContext(context.Background(), createSourcesOptions)
}

// CreateSourcesWithContext is an alternate form of the CreateSources method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) CreateSourcesWithContext(ctx context.Context, createSourcesOptions *CreateSourcesOptions) (result *Source, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(createSourcesOptions, "createSourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(createSourcesOptions, "createSourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *createSourcesOptions.LocationID,
	}

	builder := core.NewRequestBuilder(core.POST)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/sources`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range createSourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "CreateSources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if createSourcesOptions.Type != nil {
		body["type"] = createSourcesOptions.Type
	}
	if createSourcesOptions.SourceName != nil {
		body["source_name"] = createSourcesOptions.SourceName
	}
	if createSourcesOptions.Addresses != nil {
		body["addresses"] = createSourcesOptions.Addresses
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSource)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSources : update source [Administrator, Editor, Operator, satellite-link-administrator]
// Update a Source.
func (satelliteLink *SatelliteLinkV1) UpdateSources(updateSourcesOptions *UpdateSourcesOptions) (result *Source, response *core.DetailedResponse, err error) {
	return satelliteLink.UpdateSourcesWithContext(context.Background(), updateSourcesOptions)
}

// UpdateSourcesWithContext is an alternate form of the UpdateSources method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) UpdateSourcesWithContext(ctx context.Context, updateSourcesOptions *UpdateSourcesOptions) (result *Source, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSourcesOptions, "updateSourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSourcesOptions, "updateSourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *updateSourcesOptions.LocationID,
		"source_id":   *updateSourcesOptions.SourceID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/sources/{source_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "UpdateSources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSourcesOptions.SourceName != nil {
		body["source_name"] = updateSourcesOptions.SourceName
	}
	if updateSourcesOptions.Addresses != nil {
		body["addresses"] = updateSourcesOptions.Addresses
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalSource)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// DeleteSources : delete source [Administrator, Editor, Operator, satellite-link-administrator]
// Delete a source of a location.
func (satelliteLink *SatelliteLinkV1) DeleteSources(deleteSourcesOptions *DeleteSourcesOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	return satelliteLink.DeleteSourcesWithContext(context.Background(), deleteSourcesOptions)
}

// DeleteSourcesWithContext is an alternate form of the DeleteSources method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) DeleteSourcesWithContext(ctx context.Context, deleteSourcesOptions *DeleteSourcesOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(deleteSourcesOptions, "deleteSourcesOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(deleteSourcesOptions, "deleteSourcesOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *deleteSourcesOptions.LocationID,
		"source_id":   *deleteSourcesOptions.SourceID,
	}

	builder := core.NewRequestBuilder(core.DELETE)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/sources/{source_id}`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range deleteSourcesOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "DeleteSources")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExecutionResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// ListSourceEndpoints : list source status for all endpoints [Administrator, Editor, Viewer, Operator, satellite-link-administrator]
// List source status for multiple endpoints. For user source, only cloud endpoints status will be shown, for service
// source, only location endpoints will be shown.
func (satelliteLink *SatelliteLinkV1) ListSourceEndpoints(listSourceEndpointsOptions *ListSourceEndpointsOptions) (result *EndpointSourceStatus, response *core.DetailedResponse, err error) {
	return satelliteLink.ListSourceEndpointsWithContext(context.Background(), listSourceEndpointsOptions)
}

// ListSourceEndpointsWithContext is an alternate form of the ListSourceEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) ListSourceEndpointsWithContext(ctx context.Context, listSourceEndpointsOptions *ListSourceEndpointsOptions) (result *EndpointSourceStatus, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(listSourceEndpointsOptions, "listSourceEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(listSourceEndpointsOptions, "listSourceEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *listSourceEndpointsOptions.LocationID,
		"source_id":   *listSourceEndpointsOptions.SourceID,
	}

	builder := core.NewRequestBuilder(core.GET)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/sources/{source_id}/endpoints`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range listSourceEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "ListSourceEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")

	request, err := builder.Build()
	if err != nil {
		return
	}

	var rawResponse map[string]json.RawMessage
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalEndpointSourceStatus)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// UpdateSourceEndpoints : update source status for listed endpoints [Administrator, Editor, Operator, satellite-link-source-access-controller]
// Update source status for multiple endpoints. When getting error, the source status on some of the endpoints might
// still be updated successfully and only the first error will be returned.
func (satelliteLink *SatelliteLinkV1) UpdateSourceEndpoints(updateSourceEndpointsOptions *UpdateSourceEndpointsOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	return satelliteLink.UpdateSourceEndpointsWithContext(context.Background(), updateSourceEndpointsOptions)
}

// UpdateSourceEndpointsWithContext is an alternate form of the UpdateSourceEndpoints method which supports a Context parameter
func (satelliteLink *SatelliteLinkV1) UpdateSourceEndpointsWithContext(ctx context.Context, updateSourceEndpointsOptions *UpdateSourceEndpointsOptions) (result *ExecutionResult, response *core.DetailedResponse, err error) {
	err = core.ValidateNotNil(updateSourceEndpointsOptions, "updateSourceEndpointsOptions cannot be nil")
	if err != nil {
		return
	}
	err = core.ValidateStruct(updateSourceEndpointsOptions, "updateSourceEndpointsOptions")
	if err != nil {
		return
	}

	pathParamsMap := map[string]string{
		"location_id": *updateSourceEndpointsOptions.LocationID,
		"source_id":   *updateSourceEndpointsOptions.SourceID,
	}

	builder := core.NewRequestBuilder(core.PATCH)
	builder = builder.WithContext(ctx)
	builder.EnableGzipCompression = satelliteLink.GetEnableGzipCompression()
	_, err = builder.ResolveRequestURL(satelliteLink.Service.Options.URL, `/v1/locations/{location_id}/sources/{source_id}/endpoints`, pathParamsMap)
	if err != nil {
		return
	}

	for headerName, headerValue := range updateSourceEndpointsOptions.Headers {
		builder.AddHeader(headerName, headerValue)
	}

	sdkHeaders := common.GetSdkHeaders("satellite_link", "V1", "UpdateSourceEndpoints")
	for headerName, headerValue := range sdkHeaders {
		builder.AddHeader(headerName, headerValue)
	}
	builder.AddHeader("Accept", "application/json")
	builder.AddHeader("Content-Type", "application/json")

	body := make(map[string]interface{})
	if updateSourceEndpointsOptions.Endpoints != nil {
		body["endpoints"] = updateSourceEndpointsOptions.Endpoints
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
	response, err = satelliteLink.Service.Request(request, &rawResponse)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(rawResponse, "", &result, UnmarshalExecutionResult)
	if err != nil {
		return
	}
	response.Result = result

	return
}

// AdditionalNewEndpointRequestCerts : The certs.
type AdditionalNewEndpointRequestCerts struct {
	// The CA which Satellite Link trust when receiving the connection from the client application.
	Client *AdditionalNewEndpointRequestCertsClient `json:"client,omitempty"`

	// The CA which Satellite Link trust when sending the connection to server application.
	Server *AdditionalNewEndpointRequestCertsServer `json:"server,omitempty"`

	// The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.
	Connector *AdditionalNewEndpointRequestCertsConnector `json:"connector,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCerts unmarshals an instance of AdditionalNewEndpointRequestCerts from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCerts(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCerts)
	err = core.UnmarshalModel(m, "client", &obj.Client, UnmarshalAdditionalNewEndpointRequestCertsClient)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "server", &obj.Server, UnmarshalAdditionalNewEndpointRequestCertsServer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "connector", &obj.Connector, UnmarshalAdditionalNewEndpointRequestCertsConnector)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsClient : The CA which Satellite Link trust when receiving the connection from the client application.
type AdditionalNewEndpointRequestCertsClient struct {
	// The root cert or the self-signed cert of the client application.
	Cert *AdditionalNewEndpointRequestCertsClientCert `json:"cert,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsClient unmarshals an instance of AdditionalNewEndpointRequestCertsClient from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsClient(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsClient)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalAdditionalNewEndpointRequestCertsClientCert)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsClientCert : The root cert or the self-signed cert of the client application.
type AdditionalNewEndpointRequestCertsClientCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`

	// The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsClientCert unmarshals an instance of AdditionalNewEndpointRequestCertsClientCert from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsClientCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsClientCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsConnector : The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.
type AdditionalNewEndpointRequestCertsConnector struct {
	// The end-entity cert. This is required when the key is defined.
	Cert *AdditionalNewEndpointRequestCertsConnectorCert `json:"cert,omitempty"`

	// The private key of the end-entity certificate. This is required when the cert is defined.
	Key *AdditionalNewEndpointRequestCertsConnectorKey `json:"key,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsConnector unmarshals an instance of AdditionalNewEndpointRequestCertsConnector from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsConnector(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsConnector)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalAdditionalNewEndpointRequestCertsConnectorCert)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "key", &obj.Key, UnmarshalAdditionalNewEndpointRequestCertsConnectorKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsConnectorCert : The end-entity cert. This is required when the key is defined.
type AdditionalNewEndpointRequestCertsConnectorCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`

	// The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsConnectorCert unmarshals an instance of AdditionalNewEndpointRequestCertsConnectorCert from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsConnectorCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsConnectorCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsConnectorKey : The private key of the end-entity certificate. This is required when the cert is defined.
type AdditionalNewEndpointRequestCertsConnectorKey struct {
	// The name of the key.
	Filename *string `json:"filename,omitempty"`

	// The content of the key. The private key file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsConnectorKey unmarshals an instance of AdditionalNewEndpointRequestCertsConnectorKey from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsConnectorKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsConnectorKey)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsServer : The CA which Satellite Link trust when sending the connection to server application.
type AdditionalNewEndpointRequestCertsServer struct {
	// The root cert or the self-signed cert of the server application.
	Cert *AdditionalNewEndpointRequestCertsServerCert `json:"cert,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsServer unmarshals an instance of AdditionalNewEndpointRequestCertsServer from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsServer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsServer)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalAdditionalNewEndpointRequestCertsServerCert)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// AdditionalNewEndpointRequestCertsServerCert : The root cert or the self-signed cert of the server application.
type AdditionalNewEndpointRequestCertsServerCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`

	// The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalAdditionalNewEndpointRequestCertsServerCert unmarshals an instance of AdditionalNewEndpointRequestCertsServerCert from the specified map of raw messages.
func UnmarshalAdditionalNewEndpointRequestCertsServerCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(AdditionalNewEndpointRequestCertsServerCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// CreateEndpointsOptions : The CreateEndpoints options.
type CreateEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The type of the endpoint.
	ConnType *string

	// The display name of the endpoint. Endpoint names must start with a letter and end with an alphanumeric character,
	// can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	DisplayName *string

	// The host name or IP address of the server endpoint. For 'http-tunnel' protocol, server_host can start with '*.' ,
	// which means a wildcard to it's sub domains. Such as '*.example.com' can accept request to 'api.example.com' and
	// 'www.example.com'.
	ServerHost *string

	// The port number of the server endpoint. For 'http-tunnel' protocol, server_port can be 0, which means any port. Such
	// as 0 is good for 80 (http) and 443 (https).
	ServerPort *int64

	// The server name indicator (SNI) which used to connect to the server endpoint. Only useful if server side requires
	// SNI.
	Sni *string

	// The protocol in the client application side.
	ClientProtocol *string

	// Whether enable mutual auth in the client application side, when client_protocol is 'tls' or 'https', this field is
	// required.
	ClientMutualAuth *bool

	// The protocol in the server application side. This parameter will change to default value if it is omitted even when
	// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
	// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol
	// could be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
	ServerProtocol *string

	// Whether enable mutual auth in the server application side, when client_protocol is 'tls', this field is required.
	ServerMutualAuth *bool

	// Whether reject any connection to the server application which is not authorized with the list of supplied CAs in the
	// fields certs.server_cert.
	RejectUnauth *bool

	// The inactivity timeout in the Endpoint side.
	Timeout *int64

	// The service or person who created the endpoint. Must be 1000 characters or fewer.
	CreatedBy *string

	// The certs.
	Certs *AdditionalNewEndpointRequestCerts

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateEndpointsOptions.ConnType property.
// The type of the endpoint.
const (
	CreateEndpointsOptions_ConnType_Cloud    = "cloud"
	CreateEndpointsOptions_ConnType_Location = "location"
)

// Constants associated with the CreateEndpointsOptions.ClientProtocol property.
// The protocol in the client application side.
const (
	CreateEndpointsOptions_ClientProtocol_Http       = "http"
	CreateEndpointsOptions_ClientProtocol_HttpTunnel = "http-tunnel"
	CreateEndpointsOptions_ClientProtocol_Https      = "https"
	CreateEndpointsOptions_ClientProtocol_Tcp        = "tcp"
	CreateEndpointsOptions_ClientProtocol_Tls        = "tls"
	CreateEndpointsOptions_ClientProtocol_Udp        = "udp"
)

// Constants associated with the CreateEndpointsOptions.ServerProtocol property.
// The protocol in the server application side. This parameter will change to default value if it is omitted even when
// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol could
// be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
const (
	CreateEndpointsOptions_ServerProtocol_Tcp = "tcp"
	CreateEndpointsOptions_ServerProtocol_Tls = "tls"
	CreateEndpointsOptions_ServerProtocol_Udp = "udp"
)

// NewCreateEndpointsOptions : Instantiate CreateEndpointsOptions
func (*SatelliteLinkV1) NewCreateEndpointsOptions(locationID string) *CreateEndpointsOptions {
	return &CreateEndpointsOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *CreateEndpointsOptions) SetLocationID(locationID string) *CreateEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetConnType : Allow user to set ConnType
func (options *CreateEndpointsOptions) SetConnType(connType string) *CreateEndpointsOptions {
	options.ConnType = core.StringPtr(connType)
	return options
}

// SetDisplayName : Allow user to set DisplayName
func (options *CreateEndpointsOptions) SetDisplayName(displayName string) *CreateEndpointsOptions {
	options.DisplayName = core.StringPtr(displayName)
	return options
}

// SetServerHost : Allow user to set ServerHost
func (options *CreateEndpointsOptions) SetServerHost(serverHost string) *CreateEndpointsOptions {
	options.ServerHost = core.StringPtr(serverHost)
	return options
}

// SetServerPort : Allow user to set ServerPort
func (options *CreateEndpointsOptions) SetServerPort(serverPort int64) *CreateEndpointsOptions {
	options.ServerPort = core.Int64Ptr(serverPort)
	return options
}

// SetSni : Allow user to set Sni
func (options *CreateEndpointsOptions) SetSni(sni string) *CreateEndpointsOptions {
	options.Sni = core.StringPtr(sni)
	return options
}

// SetClientProtocol : Allow user to set ClientProtocol
func (options *CreateEndpointsOptions) SetClientProtocol(clientProtocol string) *CreateEndpointsOptions {
	options.ClientProtocol = core.StringPtr(clientProtocol)
	return options
}

// SetClientMutualAuth : Allow user to set ClientMutualAuth
func (options *CreateEndpointsOptions) SetClientMutualAuth(clientMutualAuth bool) *CreateEndpointsOptions {
	options.ClientMutualAuth = core.BoolPtr(clientMutualAuth)
	return options
}

// SetServerProtocol : Allow user to set ServerProtocol
func (options *CreateEndpointsOptions) SetServerProtocol(serverProtocol string) *CreateEndpointsOptions {
	options.ServerProtocol = core.StringPtr(serverProtocol)
	return options
}

// SetServerMutualAuth : Allow user to set ServerMutualAuth
func (options *CreateEndpointsOptions) SetServerMutualAuth(serverMutualAuth bool) *CreateEndpointsOptions {
	options.ServerMutualAuth = core.BoolPtr(serverMutualAuth)
	return options
}

// SetRejectUnauth : Allow user to set RejectUnauth
func (options *CreateEndpointsOptions) SetRejectUnauth(rejectUnauth bool) *CreateEndpointsOptions {
	options.RejectUnauth = core.BoolPtr(rejectUnauth)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *CreateEndpointsOptions) SetTimeout(timeout int64) *CreateEndpointsOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetCreatedBy : Allow user to set CreatedBy
func (options *CreateEndpointsOptions) SetCreatedBy(createdBy string) *CreateEndpointsOptions {
	options.CreatedBy = core.StringPtr(createdBy)
	return options
}

// SetCerts : Allow user to set Certs
func (options *CreateEndpointsOptions) SetCerts(certs *AdditionalNewEndpointRequestCerts) *CreateEndpointsOptions {
	options.Certs = certs
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateEndpointsOptions) SetHeaders(param map[string]string) *CreateEndpointsOptions {
	options.Headers = param
	return options
}

// CreateLinkOptions : The CreateLink options.
type CreateLinkOptions struct {
	// CRN of the Location.
	Crn *string

	// Location ID.
	LocationID *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewCreateLinkOptions : Instantiate CreateLinkOptions
func (*SatelliteLinkV1) NewCreateLinkOptions() *CreateLinkOptions {
	return &CreateLinkOptions{}
}

// SetCrn : Allow user to set Crn
func (options *CreateLinkOptions) SetCrn(crn string) *CreateLinkOptions {
	options.Crn = core.StringPtr(crn)
	return options
}

// SetLocationID : Allow user to set LocationID
func (options *CreateLinkOptions) SetLocationID(locationID string) *CreateLinkOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateLinkOptions) SetHeaders(param map[string]string) *CreateLinkOptions {
	options.Headers = param
	return options
}

// CreateSourcesOptions : The CreateSources options.
type CreateSourcesOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The type of the source.
	Type *string

	// The name of the source, should be unique under each location. Source names must start with a letter and end with an
	// alphanumeric character, can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	SourceName *string

	Addresses []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the CreateSourcesOptions.Type property.
// The type of the source.
const (
	CreateSourcesOptions_Type_Service = "service"
	CreateSourcesOptions_Type_User    = "user"
)

// NewCreateSourcesOptions : Instantiate CreateSourcesOptions
func (*SatelliteLinkV1) NewCreateSourcesOptions(locationID string) *CreateSourcesOptions {
	return &CreateSourcesOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *CreateSourcesOptions) SetLocationID(locationID string) *CreateSourcesOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetType : Allow user to set Type
func (options *CreateSourcesOptions) SetType(typeVar string) *CreateSourcesOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetSourceName : Allow user to set SourceName
func (options *CreateSourcesOptions) SetSourceName(sourceName string) *CreateSourcesOptions {
	options.SourceName = core.StringPtr(sourceName)
	return options
}

// SetAddresses : Allow user to set Addresses
func (options *CreateSourcesOptions) SetAddresses(addresses []string) *CreateSourcesOptions {
	options.Addresses = addresses
	return options
}

// SetHeaders : Allow user to set Headers
func (options *CreateSourcesOptions) SetHeaders(param map[string]string) *CreateSourcesOptions {
	options.Headers = param
	return options
}

// DeleteEndpointCertsOptions : The DeleteEndpointCerts options.
type DeleteEndpointCertsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEndpointCertsOptions : Instantiate DeleteEndpointCertsOptions
func (*SatelliteLinkV1) NewDeleteEndpointCertsOptions(locationID string, endpointID string) *DeleteEndpointCertsOptions {
	return &DeleteEndpointCertsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *DeleteEndpointCertsOptions) SetLocationID(locationID string) *DeleteEndpointCertsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *DeleteEndpointCertsOptions) SetEndpointID(endpointID string) *DeleteEndpointCertsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEndpointCertsOptions) SetHeaders(param map[string]string) *DeleteEndpointCertsOptions {
	options.Headers = param
	return options
}

// DeleteEndpointsOptions : The DeleteEndpoints options.
type DeleteEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteEndpointsOptions : Instantiate DeleteEndpointsOptions
func (*SatelliteLinkV1) NewDeleteEndpointsOptions(locationID string, endpointID string) *DeleteEndpointsOptions {
	return &DeleteEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *DeleteEndpointsOptions) SetLocationID(locationID string) *DeleteEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *DeleteEndpointsOptions) SetEndpointID(endpointID string) *DeleteEndpointsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteEndpointsOptions) SetHeaders(param map[string]string) *DeleteEndpointsOptions {
	options.Headers = param
	return options
}

// DeleteLinkOptions : The DeleteLink options.
type DeleteLinkOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteLinkOptions : Instantiate DeleteLinkOptions
func (*SatelliteLinkV1) NewDeleteLinkOptions(locationID string) *DeleteLinkOptions {
	return &DeleteLinkOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *DeleteLinkOptions) SetLocationID(locationID string) *DeleteLinkOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteLinkOptions) SetHeaders(param map[string]string) *DeleteLinkOptions {
	options.Headers = param
	return options
}

// DeleteSourcesOptions : The DeleteSources options.
type DeleteSourcesOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Source ID.
	SourceID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewDeleteSourcesOptions : Instantiate DeleteSourcesOptions
func (*SatelliteLinkV1) NewDeleteSourcesOptions(locationID string, sourceID string) *DeleteSourcesOptions {
	return &DeleteSourcesOptions{
		LocationID: core.StringPtr(locationID),
		SourceID:   core.StringPtr(sourceID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *DeleteSourcesOptions) SetLocationID(locationID string) *DeleteSourcesOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetSourceID : Allow user to set SourceID
func (options *DeleteSourcesOptions) SetSourceID(sourceID string) *DeleteSourcesOptions {
	options.SourceID = core.StringPtr(sourceID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *DeleteSourcesOptions) SetHeaders(param map[string]string) *DeleteSourcesOptions {
	options.Headers = param
	return options
}

// DownloadedCerts : The list of certs.
type DownloadedCerts struct {
	// The array of the cert(s) and key(s).
	Certs []DownloadedCertsCertsItem `json:"certs,omitempty"`
}

// UnmarshalDownloadedCerts unmarshals an instance of DownloadedCerts from the specified map of raw messages.
func UnmarshalDownloadedCerts(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DownloadedCerts)
	err = core.UnmarshalModel(m, "certs", &obj.Certs, UnmarshalDownloadedCertsCertsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// DownloadedCertsCertsItem : DownloadedCertsCertsItem struct
type DownloadedCertsCertsItem struct {
	// The filename of the cert or key.
	Name *string `json:"name,omitempty"`

	// The content of the cert or key.
	Content *string `json:"content,omitempty"`
}

// UnmarshalDownloadedCertsCertsItem unmarshals an instance of DownloadedCertsCertsItem from the specified map of raw messages.
func UnmarshalDownloadedCertsCertsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(DownloadedCertsCertsItem)
	err = core.UnmarshalPrimitive(m, "name", &obj.Name)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "content", &obj.Content)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Endpoint : The Source Status.
type Endpoint struct {
	// The type of the endpoint.
	ConnType *string `json:"conn_type,omitempty"`

	// The display name of the endpoint. Endpoint names must start with a letter and end with an alphanumeric character,
	// can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	DisplayName *string `json:"display_name,omitempty"`

	// The host name or IP address of the server endpoint. For 'http-tunnel' protocol, server_host can start with '*.' ,
	// which means a wildcard to it's sub domains. Such as '*.example.com' can accept request to 'api.example.com' and
	// 'www.example.com'.
	ServerHost *string `json:"server_host,omitempty"`

	// The port number of the server endpoint. For 'http-tunnel' protocol, server_port can be 0, which means any port. Such
	// as 0 is good for 80 (http) and 443 (https).
	ServerPort *int64 `json:"server_port,omitempty"`

	// The server name indicator (SNI) which used to connect to the server endpoint. Only useful if server side requires
	// SNI.
	Sni *string `json:"sni,omitempty"`

	// The protocol in the client application side.
	ClientProtocol *string `json:"client_protocol,omitempty"`

	// Whether enable mutual auth in the client application side, when client_protocol is 'tls' or 'https', this field is
	// required.
	ClientMutualAuth *bool `json:"client_mutual_auth,omitempty"`

	// The protocol in the server application side. This parameter will change to default value if it is omitted even when
	// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
	// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol
	// could be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
	ServerProtocol *string `json:"server_protocol,omitempty"`

	// Whether enable mutual auth in the server application side, when client_protocol is 'tls', this field is required.
	ServerMutualAuth *bool `json:"server_mutual_auth,omitempty"`

	// Whether reject any connection to the server application which is not authorized with the list of supplied CAs in the
	// fields certs.server_cert.
	RejectUnauth *bool `json:"reject_unauth,omitempty"`

	// The inactivity timeout in the Endpoint side.
	Timeout *int64 `json:"timeout,omitempty"`

	// The service or person who created the endpoint. Must be 1000 characters or fewer.
	CreatedBy *string `json:"created_by,omitempty"`

	Sources []SourceStatusObject `json:"sources,omitempty"`

	// The connector port.
	ConnectorPort *int64 `json:"connector_port,omitempty"`

	// Service instance associated with this location.
	Crn *string `json:"crn,omitempty"`

	// Unique identifier for this endpoint.
	EndpointID *string `json:"endpoint_id,omitempty"`

	// The service name of the endpoint.
	ServiceName *string `json:"service_name,omitempty"`

	// The Location ID.
	LocationID *string `json:"location_id,omitempty"`

	// The hostname which Satellite Link server listen on for the on-location endpoint, or the hostname which the connector
	// server listen on for the on-cloud endpoint destiantion.
	ClientHost *string `json:"client_host,omitempty"`

	// The port which Satellite Link server listen on for the on-location, or the port which the connector server listen on
	// for the on-cloud endpoint destiantion.
	ClientPort *int64 `json:"client_port,omitempty"`

	// The certs. Once it is generated, this field will always be defined even it is unused until the cert/key is deleted.
	Certs *EndpointCerts `json:"certs,omitempty"`

	// Whether the Endpoint is active or not.
	Status *string `json:"status,omitempty"`

	// The time when the Endpoint is created.
	CreatedAt *string `json:"created_at,omitempty"`

	// The last time modify the Endpoint configurations.
	LastChange *string `json:"last_change,omitempty"`

	// The last performance data of the endpoint.
	Performance *EndpointPerformance `json:"performance,omitempty"`
}

// Constants associated with the Endpoint.ConnType property.
// The type of the endpoint.
const (
	Endpoint_ConnType_Cloud    = "cloud"
	Endpoint_ConnType_Location = "location"
)

// Constants associated with the Endpoint.ClientProtocol property.
// The protocol in the client application side.
const (
	Endpoint_ClientProtocol_Http       = "http"
	Endpoint_ClientProtocol_HttpTunnel = "http-tunnel"
	Endpoint_ClientProtocol_Https      = "https"
	Endpoint_ClientProtocol_Tcp        = "tcp"
	Endpoint_ClientProtocol_Tls        = "tls"
	Endpoint_ClientProtocol_Udp        = "udp"
)

// Constants associated with the Endpoint.ServerProtocol property.
// The protocol in the server application side. This parameter will change to default value if it is omitted even when
// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol could
// be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
const (
	Endpoint_ServerProtocol_Tcp = "tcp"
	Endpoint_ServerProtocol_Tls = "tls"
	Endpoint_ServerProtocol_Udp = "udp"
)

// Constants associated with the Endpoint.Status property.
// Whether the Endpoint is active or not.
const (
	Endpoint_Status_Disabled = "disabled"
	Endpoint_Status_Enabled  = "enabled"
)

// UnmarshalEndpoint unmarshals an instance of Endpoint from the specified map of raw messages.
func UnmarshalEndpoint(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Endpoint)
	err = core.UnmarshalPrimitive(m, "conn_type", &obj.ConnType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_host", &obj.ServerHost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_port", &obj.ServerPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sni", &obj.Sni)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_protocol", &obj.ClientProtocol)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_mutual_auth", &obj.ClientMutualAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_protocol", &obj.ServerProtocol)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_mutual_auth", &obj.ServerMutualAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reject_unauth", &obj.RejectUnauth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_by", &obj.CreatedBy)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "sources", &obj.Sources, UnmarshalSourceStatusObject)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "connector_port", &obj.ConnectorPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "endpoint_id", &obj.EndpointID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "service_name", &obj.ServiceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location_id", &obj.LocationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_host", &obj.ClientHost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_port", &obj.ClientPort)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "certs", &obj.Certs, UnmarshalEndpointCerts)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_change", &obj.LastChange)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "performance", &obj.Performance, UnmarshalEndpointPerformance)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCerts : The certs. Once it is generated, this field will always be defined even it is unused until the cert/key is deleted.
type EndpointCerts struct {
	// The CA which Satellite Link trust when receiving the connection from the client application.
	Client *EndpointCertsClient `json:"client,omitempty"`

	// The CA which Satellite Link trust when sending the connection to server application.
	Server *EndpointCertsServer `json:"server,omitempty"`

	// The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.
	Connector *EndpointCertsConnector `json:"connector,omitempty"`
}

// UnmarshalEndpointCerts unmarshals an instance of EndpointCerts from the specified map of raw messages.
func UnmarshalEndpointCerts(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCerts)
	err = core.UnmarshalModel(m, "client", &obj.Client, UnmarshalEndpointCertsClient)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "server", &obj.Server, UnmarshalEndpointCertsServer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "connector", &obj.Connector, UnmarshalEndpointCertsConnector)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsClient : The CA which Satellite Link trust when receiving the connection from the client application.
type EndpointCertsClient struct {
	// The root cert or the self-signed cert of the client application.
	Cert *EndpointCertsClientCert `json:"cert,omitempty"`
}

// UnmarshalEndpointCertsClient unmarshals an instance of EndpointCertsClient from the specified map of raw messages.
func UnmarshalEndpointCertsClient(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsClient)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalEndpointCertsClientCert)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsClientCert : The root cert or the self-signed cert of the client application.
type EndpointCertsClientCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`
}

// UnmarshalEndpointCertsClientCert unmarshals an instance of EndpointCertsClientCert from the specified map of raw messages.
func UnmarshalEndpointCertsClientCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsClientCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsConnector : The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.
type EndpointCertsConnector struct {
	// The end-entity cert of the connector.
	Cert *EndpointCertsConnectorCert `json:"cert,omitempty"`

	// The private key of the connector.
	Key *EndpointCertsConnectorKey `json:"key,omitempty"`
}

// UnmarshalEndpointCertsConnector unmarshals an instance of EndpointCertsConnector from the specified map of raw messages.
func UnmarshalEndpointCertsConnector(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsConnector)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalEndpointCertsConnectorCert)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "key", &obj.Key, UnmarshalEndpointCertsConnectorKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsConnectorCert : The end-entity cert of the connector.
type EndpointCertsConnectorCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`
}

// UnmarshalEndpointCertsConnectorCert unmarshals an instance of EndpointCertsConnectorCert from the specified map of raw messages.
func UnmarshalEndpointCertsConnectorCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsConnectorCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsConnectorKey : The private key of the connector.
type EndpointCertsConnectorKey struct {
	// The name of the key.
	Filename *string `json:"filename,omitempty"`
}

// UnmarshalEndpointCertsConnectorKey unmarshals an instance of EndpointCertsConnectorKey from the specified map of raw messages.
func UnmarshalEndpointCertsConnectorKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsConnectorKey)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsServer : The CA which Satellite Link trust when sending the connection to server application.
type EndpointCertsServer struct {
	// The root cert or the self-signed cert of the server application.
	Cert *EndpointCertsServerCert `json:"cert,omitempty"`
}

// UnmarshalEndpointCertsServer unmarshals an instance of EndpointCertsServer from the specified map of raw messages.
func UnmarshalEndpointCertsServer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsServer)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalEndpointCertsServerCert)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointCertsServerCert : The root cert or the self-signed cert of the server application.
type EndpointCertsServerCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`
}

// UnmarshalEndpointCertsServerCert unmarshals an instance of EndpointCertsServerCert from the specified map of raw messages.
func UnmarshalEndpointCertsServerCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointCertsServerCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointPerformance : The last performance data of the endpoint.
type EndpointPerformance struct {
	// Concurrent connections number of moment when probe read the data.
	Connection *int64 `json:"connection,omitempty"`

	// Average Receive (to Cloud) Bandwidth of last two minutes, unit is Byte/s.
	RxBandwidth *int64 `json:"rx_bandwidth,omitempty"`

	// Average Transmitted (to Location) Bandwidth of last two minutes, unit is Byte/s.
	TxBandwidth *int64 `json:"tx_bandwidth,omitempty"`

	// Average Tatal Bandwidth of last two minutes, unit is Byte/s.
	Bandwidth *int64 `json:"bandwidth,omitempty"`

	// The last performance data of the endpoint from each Connector.
	Connectors []EndpointPerformanceConnectorsItem `json:"connectors,omitempty"`
}

// UnmarshalEndpointPerformance unmarshals an instance of EndpointPerformance from the specified map of raw messages.
func UnmarshalEndpointPerformance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointPerformance)
	err = core.UnmarshalPrimitive(m, "connection", &obj.Connection)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rx_bandwidth", &obj.RxBandwidth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tx_bandwidth", &obj.TxBandwidth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bandwidth", &obj.Bandwidth)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "connectors", &obj.Connectors, UnmarshalEndpointPerformanceConnectorsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointPerformanceConnectorsItem : EndpointPerformanceConnectorsItem struct
type EndpointPerformanceConnectorsItem struct {
	// The name of the connector reported the performance data.
	Connector *string `json:"connector,omitempty"`

	// Concurrent connections number of moment when probe read the data from the Connector.
	Connections *int64 `json:"connections,omitempty"`

	// Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.
	RxBW *int64 `json:"rxBW,omitempty"`

	// Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.
	TxBW *int64 `json:"txBW,omitempty"`
}

// UnmarshalEndpointPerformanceConnectorsItem unmarshals an instance of EndpointPerformanceConnectorsItem from the specified map of raw messages.
func UnmarshalEndpointPerformanceConnectorsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointPerformanceConnectorsItem)
	err = core.UnmarshalPrimitive(m, "connector", &obj.Connector)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "connections", &obj.Connections)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rxBW", &obj.RxBW)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "txBW", &obj.TxBW)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointSourceStatus : EndpointSourceStatus struct
type EndpointSourceStatus struct {
	Endpoints []EndpointSourceStatusEndpointsItem `json:"endpoints,omitempty"`
}

// UnmarshalEndpointSourceStatus unmarshals an instance of EndpointSourceStatus from the specified map of raw messages.
func UnmarshalEndpointSourceStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointSourceStatus)
	err = core.UnmarshalModel(m, "endpoints", &obj.Endpoints, UnmarshalEndpointSourceStatusEndpointsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// EndpointSourceStatusEndpointsItem : EndpointSourceStatusEndpointsItem struct
type EndpointSourceStatusEndpointsItem struct {
	// Unique identifier for this endpoint.
	EndpointID *string `json:"endpoint_id,omitempty"`

	// Whether the source is enabled for the endpoint.
	Enabled *bool `json:"enabled,omitempty"`
}

// UnmarshalEndpointSourceStatusEndpointsItem unmarshals an instance of EndpointSourceStatusEndpointsItem from the specified map of raw messages.
func UnmarshalEndpointSourceStatusEndpointsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(EndpointSourceStatusEndpointsItem)
	err = core.UnmarshalPrimitive(m, "endpoint_id", &obj.EndpointID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Endpoints : The list of the endpoint(s).
type Endpoints struct {
	// The info of the endpoint.
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// UnmarshalEndpoints unmarshals an instance of Endpoints from the specified map of raw messages.
func UnmarshalEndpoints(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Endpoints)
	err = core.UnmarshalModel(m, "endpoints", &obj.Endpoints, UnmarshalEndpoint)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExecutionResult : Result of execution.
type ExecutionResult struct {
	// Result returned.
	Message *string `json:"message,omitempty"`
}

// UnmarshalExecutionResult unmarshals an instance of ExecutionResult from the specified map of raw messages.
func UnmarshalExecutionResult(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExecutionResult)
	err = core.UnmarshalPrimitive(m, "message", &obj.Message)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// ExportEndpointsOptions : The ExportEndpoints options.
type ExportEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewExportEndpointsOptions : Instantiate ExportEndpointsOptions
func (*SatelliteLinkV1) NewExportEndpointsOptions(locationID string, endpointID string) *ExportEndpointsOptions {
	return &ExportEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *ExportEndpointsOptions) SetLocationID(locationID string) *ExportEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *ExportEndpointsOptions) SetEndpointID(endpointID string) *ExportEndpointsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ExportEndpointsOptions) SetHeaders(param map[string]string) *ExportEndpointsOptions {
	options.Headers = param
	return options
}

// ExportEndpointsResponse : The info of the endpoint.
type ExportEndpointsResponse struct {
	// The type of the endpoint.
	ConnType *string `json:"conn_type,omitempty"`

	// The display name of the endpoint. Endpoint names must start with a letter and end with an alphanumeric character,
	// can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	DisplayName *string `json:"display_name,omitempty"`

	// The host name or IP address of the server endpoint. For 'http-tunnel' protocol, server_host can start with '*.' ,
	// which means a wildcard to it's sub domains. Such as '*.example.com' can accept request to 'api.example.com' and
	// 'www.example.com'.
	ServerHost *string `json:"server_host,omitempty"`

	// The port number of the server endpoint. For 'http-tunnel' protocol, server_port can be 0, which means any port. Such
	// as 0 is good for 80 (http) and 443 (https).
	ServerPort *int64 `json:"server_port,omitempty"`

	// The server name indicator (SNI) which used to connect to the server endpoint. Only useful if server side requires
	// SNI.
	Sni *string `json:"sni,omitempty"`

	// The protocol in the client application side.
	ClientProtocol *string `json:"client_protocol,omitempty"`

	// Whether enable mutual auth in the client application side, when client_protocol is 'tls' or 'https', this field is
	// required.
	ClientMutualAuth *bool `json:"client_mutual_auth,omitempty"`

	// The protocol in the server application side. This parameter will change to default value if it is omitted even when
	// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
	// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol
	// could be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
	ServerProtocol *string `json:"server_protocol,omitempty"`

	// Whether enable mutual auth in the server application side, when client_protocol is 'tls', this field is required.
	ServerMutualAuth *bool `json:"server_mutual_auth,omitempty"`

	// Whether reject any connection to the server application which is not authorized with the list of supplied CAs in the
	// fields certs.server_cert.
	RejectUnauth *bool `json:"reject_unauth,omitempty"`

	// The inactivity timeout in the Endpoint side.
	Timeout *int64 `json:"timeout,omitempty"`

	// The certs.
	Certs *AdditionalNewEndpointRequestCerts `json:"certs,omitempty"`
}

// Constants associated with the ExportEndpointsResponse.ConnType property.
// The type of the endpoint.
const (
	ExportEndpointsResponse_ConnType_Cloud    = "cloud"
	ExportEndpointsResponse_ConnType_Location = "location"
)

// Constants associated with the ExportEndpointsResponse.ClientProtocol property.
// The protocol in the client application side.
const (
	ExportEndpointsResponse_ClientProtocol_Http       = "http"
	ExportEndpointsResponse_ClientProtocol_HttpTunnel = "http-tunnel"
	ExportEndpointsResponse_ClientProtocol_Https      = "https"
	ExportEndpointsResponse_ClientProtocol_Tcp        = "tcp"
	ExportEndpointsResponse_ClientProtocol_Tls        = "tls"
	ExportEndpointsResponse_ClientProtocol_Udp        = "udp"
)

// Constants associated with the ExportEndpointsResponse.ServerProtocol property.
// The protocol in the server application side. This parameter will change to default value if it is omitted even when
// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol could
// be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
const (
	ExportEndpointsResponse_ServerProtocol_Tcp = "tcp"
	ExportEndpointsResponse_ServerProtocol_Tls = "tls"
	ExportEndpointsResponse_ServerProtocol_Udp = "udp"
)

// UnmarshalExportEndpointsResponse unmarshals an instance of ExportEndpointsResponse from the specified map of raw messages.
func UnmarshalExportEndpointsResponse(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(ExportEndpointsResponse)
	err = core.UnmarshalPrimitive(m, "conn_type", &obj.ConnType)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "display_name", &obj.DisplayName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_host", &obj.ServerHost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_port", &obj.ServerPort)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "sni", &obj.Sni)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_protocol", &obj.ClientProtocol)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "client_mutual_auth", &obj.ClientMutualAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_protocol", &obj.ServerProtocol)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "server_mutual_auth", &obj.ServerMutualAuth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "reject_unauth", &obj.RejectUnauth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "timeout", &obj.Timeout)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "certs", &obj.Certs, UnmarshalAdditionalNewEndpointRequestCerts)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// GetEndpointCertsOptions : The GetEndpointCerts options.
type GetEndpointCertsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// Whether the result need to be packed as zip file.
	NoZip *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEndpointCertsOptions : Instantiate GetEndpointCertsOptions
func (*SatelliteLinkV1) NewGetEndpointCertsOptions(locationID string, endpointID string) *GetEndpointCertsOptions {
	return &GetEndpointCertsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *GetEndpointCertsOptions) SetLocationID(locationID string) *GetEndpointCertsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *GetEndpointCertsOptions) SetEndpointID(endpointID string) *GetEndpointCertsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetNoZip : Allow user to set NoZip
func (options *GetEndpointCertsOptions) SetNoZip(noZip bool) *GetEndpointCertsOptions {
	options.NoZip = core.BoolPtr(noZip)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetEndpointCertsOptions) SetHeaders(param map[string]string) *GetEndpointCertsOptions {
	options.Headers = param
	return options
}

// GetEndpointsOptions : The GetEndpoints options.
type GetEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetEndpointsOptions : Instantiate GetEndpointsOptions
func (*SatelliteLinkV1) NewGetEndpointsOptions(locationID string, endpointID string) *GetEndpointsOptions {
	return &GetEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *GetEndpointsOptions) SetLocationID(locationID string) *GetEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *GetEndpointsOptions) SetEndpointID(endpointID string) *GetEndpointsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetEndpointsOptions) SetHeaders(param map[string]string) *GetEndpointsOptions {
	options.Headers = param
	return options
}

// GetLinkOptions : The GetLink options.
type GetLinkOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewGetLinkOptions : Instantiate GetLinkOptions
func (*SatelliteLinkV1) NewGetLinkOptions(locationID string) *GetLinkOptions {
	return &GetLinkOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *GetLinkOptions) SetLocationID(locationID string) *GetLinkOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *GetLinkOptions) SetHeaders(param map[string]string) *GetLinkOptions {
	options.Headers = param
	return options
}

// ImportEndpointsOptions : The ImportEndpoints options.
type ImportEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The imported .endpoint file.
	State io.ReadCloser `validate:"required"`

	// The content type of state.
	StateContentType *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewImportEndpointsOptions : Instantiate ImportEndpointsOptions
func (*SatelliteLinkV1) NewImportEndpointsOptions(locationID string, state io.ReadCloser) *ImportEndpointsOptions {
	return &ImportEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		State:      state,
	}
}

// SetLocationID : Allow user to set LocationID
func (options *ImportEndpointsOptions) SetLocationID(locationID string) *ImportEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetState : Allow user to set State
func (options *ImportEndpointsOptions) SetState(state io.ReadCloser) *ImportEndpointsOptions {
	options.State = state
	return options
}

// SetStateContentType : Allow user to set StateContentType
func (options *ImportEndpointsOptions) SetStateContentType(stateContentType string) *ImportEndpointsOptions {
	options.StateContentType = core.StringPtr(stateContentType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ImportEndpointsOptions) SetHeaders(param map[string]string) *ImportEndpointsOptions {
	options.Headers = param
	return options
}

// ListEndpointSourcesOptions : The ListEndpointSources options.
type ListEndpointSourcesOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListEndpointSourcesOptions : Instantiate ListEndpointSourcesOptions
func (*SatelliteLinkV1) NewListEndpointSourcesOptions(locationID string, endpointID string) *ListEndpointSourcesOptions {
	return &ListEndpointSourcesOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *ListEndpointSourcesOptions) SetLocationID(locationID string) *ListEndpointSourcesOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *ListEndpointSourcesOptions) SetEndpointID(endpointID string) *ListEndpointSourcesOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListEndpointSourcesOptions) SetHeaders(param map[string]string) *ListEndpointSourcesOptions {
	options.Headers = param
	return options
}

// ListEndpointsOptions : The ListEndpoints options.
type ListEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// Whether to only include enabled or disabled endpoint(s). If not specified all endpoint(s) will be returned.
	Type *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListEndpointsOptions.Type property.
// Whether to only include enabled or disabled endpoint(s). If not specified all endpoint(s) will be returned.
const (
	ListEndpointsOptions_Type_Disabled = "disabled"
	ListEndpointsOptions_Type_Enabled  = "enabled"
)

// NewListEndpointsOptions : Instantiate ListEndpointsOptions
func (*SatelliteLinkV1) NewListEndpointsOptions(locationID string) *ListEndpointsOptions {
	return &ListEndpointsOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *ListEndpointsOptions) SetLocationID(locationID string) *ListEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetType : Allow user to set Type
func (options *ListEndpointsOptions) SetType(typeVar string) *ListEndpointsOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListEndpointsOptions) SetHeaders(param map[string]string) *ListEndpointsOptions {
	options.Headers = param
	return options
}

// ListSourceEndpointsOptions : The ListSourceEndpoints options.
type ListSourceEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Source ID.
	SourceID *string `validate:"required,ne="`

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewListSourceEndpointsOptions : Instantiate ListSourceEndpointsOptions
func (*SatelliteLinkV1) NewListSourceEndpointsOptions(locationID string, sourceID string) *ListSourceEndpointsOptions {
	return &ListSourceEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		SourceID:   core.StringPtr(sourceID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *ListSourceEndpointsOptions) SetLocationID(locationID string) *ListSourceEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetSourceID : Allow user to set SourceID
func (options *ListSourceEndpointsOptions) SetSourceID(sourceID string) *ListSourceEndpointsOptions {
	options.SourceID = core.StringPtr(sourceID)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListSourceEndpointsOptions) SetHeaders(param map[string]string) *ListSourceEndpointsOptions {
	options.Headers = param
	return options
}

// ListSourcesOptions : The ListSources options.
type ListSourcesOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// Type of Sources to list, all if not specified.
	Type *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the ListSourcesOptions.Type property.
// Type of Sources to list, all if not specified.
const (
	ListSourcesOptions_Type_Service = "service"
	ListSourcesOptions_Type_User    = "user"
)

// NewListSourcesOptions : Instantiate ListSourcesOptions
func (*SatelliteLinkV1) NewListSourcesOptions(locationID string) *ListSourcesOptions {
	return &ListSourcesOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *ListSourcesOptions) SetLocationID(locationID string) *ListSourcesOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetType : Allow user to set Type
func (options *ListSourcesOptions) SetType(typeVar string) *ListSourcesOptions {
	options.Type = core.StringPtr(typeVar)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *ListSourcesOptions) SetHeaders(param map[string]string) *ListSourcesOptions {
	options.Headers = param
	return options
}

// Location : The info of the location.
type Location struct {
	// The ws endpoint of the location.
	WsEndpoint *string `json:"ws_endpoint,omitempty"`

	// Unique identifier for this location.
	LocationID *string `json:"location_id,omitempty"`

	// Service instance associated with this location.
	Crn *string `json:"crn,omitempty"`

	// Description of the location.
	Desc *string `json:"desc,omitempty"`

	// Satellite Link hostname of the location.
	SatelliteLinkHost *string `json:"satellite_link_host,omitempty"`

	// Enabled/Disabled.
	Status *string `json:"status,omitempty"`

	// Timestamp of creation of location.
	CreatedAt *string `json:"created_at,omitempty"`

	// Timestamp of latest modification of location.
	LastChange *string `json:"last_change,omitempty"`

	// The last performance data of the Location.
	Performance *LocationPerformance `json:"performance,omitempty"`
}

// Constants associated with the Location.Status property.
// Enabled/Disabled.
const (
	Location_Status_Disabled = "disabled"
	Location_Status_Enabled  = "enabled"
)

// UnmarshalLocation unmarshals an instance of Location from the specified map of raw messages.
func UnmarshalLocation(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Location)
	err = core.UnmarshalPrimitive(m, "ws_endpoint", &obj.WsEndpoint)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location_id", &obj.LocationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "crn", &obj.Crn)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "desc", &obj.Desc)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "satellite_link_host", &obj.SatelliteLinkHost)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "status", &obj.Status)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_change", &obj.LastChange)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "performance", &obj.Performance, UnmarshalLocationPerformance)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LocationPerformance : The last performance data of the Location.
type LocationPerformance struct {
	// Tunnels number estbalished from the Location.
	Tunnels *int64 `json:"tunnels,omitempty"`

	// Tunnels health status based on the Tunnels number established. Down(0)/Critical(1)/Up(>=2).
	HealthStatus *string `json:"healthStatus,omitempty"`

	// Average latency calculated form latency of each Connector between Tunnel Server, unit is ms. -1 means no Connector
	// established Tunnel.
	AvgLatency *int64 `json:"avg_latency,omitempty"`

	// Average Receive (to Cloud) Bandwidth of last two minutes, unit is Byte/s.
	RxBandwidth *int64 `json:"rx_bandwidth,omitempty"`

	// Average Transmitted (to Location) Bandwidth of last two minutes, unit is Byte/s.
	TxBandwidth *int64 `json:"tx_bandwidth,omitempty"`

	// Average Tatal Bandwidth of last two minutes, unit is Byte/s.
	Bandwidth *int64 `json:"bandwidth,omitempty"`

	// The last performance data of the Location read from each Connector.
	Connectors []LocationPerformanceConnectorsItem `json:"connectors,omitempty"`
}

// UnmarshalLocationPerformance unmarshals an instance of LocationPerformance from the specified map of raw messages.
func UnmarshalLocationPerformance(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LocationPerformance)
	err = core.UnmarshalPrimitive(m, "tunnels", &obj.Tunnels)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "healthStatus", &obj.HealthStatus)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "avg_latency", &obj.AvgLatency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rx_bandwidth", &obj.RxBandwidth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "tx_bandwidth", &obj.TxBandwidth)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "bandwidth", &obj.Bandwidth)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "connectors", &obj.Connectors, UnmarshalLocationPerformanceConnectorsItem)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// LocationPerformanceConnectorsItem : LocationPerformanceConnectorsItem struct
type LocationPerformanceConnectorsItem struct {
	// The name of the connector reported the performance data.
	Connector *string `json:"connector,omitempty"`

	// Latency between Connector and the Tunnel Server it connected.
	Latency *int64 `json:"latency,omitempty"`

	// Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.
	RxBW *int64 `json:"rxBW,omitempty"`

	// Average Transmitted (to Location) Bandwidth of last two minutes read from the Connector, unit is Byte/s.
	TxBW *int64 `json:"txBW,omitempty"`
}

// UnmarshalLocationPerformanceConnectorsItem unmarshals an instance of LocationPerformanceConnectorsItem from the specified map of raw messages.
func UnmarshalLocationPerformanceConnectorsItem(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(LocationPerformanceConnectorsItem)
	err = core.UnmarshalPrimitive(m, "connector", &obj.Connector)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "latency", &obj.Latency)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "rxBW", &obj.RxBW)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "txBW", &obj.TxBW)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Source : The Source.
type Source struct {
	// The type of the source.
	Type *string `json:"type,omitempty"`

	// The name of the source, should be unique under each location. Source names must start with a letter and end with an
	// alphanumeric character, can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	SourceName *string `json:"source_name,omitempty"`

	Addresses []string `json:"addresses,omitempty"`

	// The Source ID.
	SourceID *string `json:"source_id,omitempty"`

	// The Location ID.
	LocationID *string `json:"location_id,omitempty"`

	// Timestamp of creation of location.
	CreatedAt *string `json:"created_at,omitempty"`

	// Timestamp of creation of location.
	LastChange *string `json:"last_change,omitempty"`
}

// Constants associated with the Source.Type property.
// The type of the source.
const (
	Source_Type_Service = "service"
	Source_Type_User    = "user"
)

// UnmarshalSource unmarshals an instance of Source from the specified map of raw messages.
func UnmarshalSource(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Source)
	err = core.UnmarshalPrimitive(m, "type", &obj.Type)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_name", &obj.SourceName)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "addresses", &obj.Addresses)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "source_id", &obj.SourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "location_id", &obj.LocationID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "created_at", &obj.CreatedAt)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_change", &obj.LastChange)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SourceStatus : The Source Status.
type SourceStatus struct {
	Sources []SourceStatusObject `json:"sources,omitempty"`
}

// UnmarshalSourceStatus unmarshals an instance of SourceStatus from the specified map of raw messages.
func UnmarshalSourceStatus(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SourceStatus)
	err = core.UnmarshalModel(m, "sources", &obj.Sources, UnmarshalSourceStatusObject)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SourceStatusObject : The Source Status.
type SourceStatusObject struct {
	// The Source ID.
	SourceID *string `json:"source_id,omitempty"`

	// Whether the source is enabled for the endpoint.
	Enabled *bool `json:"enabled,omitempty"`

	// The last time modify the Endpoint configurations.
	LastChange *string `json:"last_change,omitempty"`

	// Whether the source has been enabled on this endpoint.
	Pending *bool `json:"pending,omitempty"`
}

// UnmarshalSourceStatusObject unmarshals an instance of SourceStatusObject from the specified map of raw messages.
func UnmarshalSourceStatusObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SourceStatusObject)
	err = core.UnmarshalPrimitive(m, "source_id", &obj.SourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "last_change", &obj.LastChange)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "pending", &obj.Pending)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// SourceStatusRequestObject : The Source Status.
type SourceStatusRequestObject struct {
	// The Source ID.
	SourceID *string `json:"source_id,omitempty"`

	// Whether the source is enabled for the endpoint.
	Enabled *bool `json:"enabled,omitempty"`
}

// UnmarshalSourceStatusRequestObject unmarshals an instance of SourceStatusRequestObject from the specified map of raw messages.
func UnmarshalSourceStatusRequestObject(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(SourceStatusRequestObject)
	err = core.UnmarshalPrimitive(m, "source_id", &obj.SourceID)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "enabled", &obj.Enabled)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// Sources : Sources struct
type Sources struct {
	Sources []Source `json:"sources,omitempty"`
}

// UnmarshalSources unmarshals an instance of Sources from the specified map of raw messages.
func UnmarshalSources(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(Sources)
	err = core.UnmarshalModel(m, "sources", &obj.Sources, UnmarshalSource)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdateEndpointSourcesOptions : The UpdateEndpointSources options.
type UpdateEndpointSourcesOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	Sources []SourceStatusRequestObject

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateEndpointSourcesOptions : Instantiate UpdateEndpointSourcesOptions
func (*SatelliteLinkV1) NewUpdateEndpointSourcesOptions(locationID string, endpointID string) *UpdateEndpointSourcesOptions {
	return &UpdateEndpointSourcesOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *UpdateEndpointSourcesOptions) SetLocationID(locationID string) *UpdateEndpointSourcesOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *UpdateEndpointSourcesOptions) SetEndpointID(endpointID string) *UpdateEndpointSourcesOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetSources : Allow user to set Sources
func (options *UpdateEndpointSourcesOptions) SetSources(sources []SourceStatusRequestObject) *UpdateEndpointSourcesOptions {
	options.Sources = sources
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEndpointSourcesOptions) SetHeaders(param map[string]string) *UpdateEndpointSourcesOptions {
	options.Headers = param
	return options
}

// UpdateEndpointsOptions : The UpdateEndpoints options.
type UpdateEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// The display name of the endpoint. Endpoint names must start with a letter and end with an alphanumeric character,
	// can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	DisplayName *string

	// The host name or IP address of the server endpoint. For 'http-tunnel' protocol, server_host can start with '*.' ,
	// which means a wildcard to it's sub domains. Such as '*.example.com' can accept request to 'api.example.com' and
	// 'www.example.com'.
	ServerHost *string

	// The port number of the server endpoint. For 'http-tunnel' protocol, server_port can be 0, which means any port. Such
	// as 0 is good for 80 (http) and 443 (https).
	ServerPort *int64

	// The server name indicator (SNI) which used to connect to the server endpoint. Only useful if server side requires
	// SNI.
	Sni *string

	// The protocol in the client application side.
	ClientProtocol *string

	// Whether enable mutual auth in the client application side, when client_protocol is 'tls' or 'https', this field is
	// required.
	ClientMutualAuth *bool

	// The protocol in the server application side. This parameter will change to default value if it is omitted even when
	// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
	// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol
	// could be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
	ServerProtocol *string

	// Whether enable mutual auth in the server application side, when client_protocol is 'tls', this field is required.
	ServerMutualAuth *bool

	// Whether reject any connection to the server application which is not authorized with the list of supplied CAs in the
	// fields certs.server_cert.
	RejectUnauth *bool

	// The inactivity timeout in the Endpoint side.
	Timeout *int64

	// The service or person who created the endpoint. Must be 1000 characters or fewer.
	CreatedBy *string

	// The certs.
	Certs *UpdatedEndpointRequestCerts

	// Enable or disable the endpoint.
	Enabled *bool

	// Allows users to set headers on API requests
	Headers map[string]string
}

// Constants associated with the UpdateEndpointsOptions.ClientProtocol property.
// The protocol in the client application side.
const (
	UpdateEndpointsOptions_ClientProtocol_Http       = "http"
	UpdateEndpointsOptions_ClientProtocol_HttpTunnel = "http-tunnel"
	UpdateEndpointsOptions_ClientProtocol_Https      = "https"
	UpdateEndpointsOptions_ClientProtocol_Tcp        = "tcp"
	UpdateEndpointsOptions_ClientProtocol_Tls        = "tls"
	UpdateEndpointsOptions_ClientProtocol_Udp        = "udp"
)

// Constants associated with the UpdateEndpointsOptions.ServerProtocol property.
// The protocol in the server application side. This parameter will change to default value if it is omitted even when
// using PATCH API. If client_protocol is 'udp', server_protocol must be 'udp'. If client_protocol is 'tcp'/'http',
// server_protocol could be 'tcp'/'tls' and default to 'tcp'. If client_protocol is 'tls'/'https', server_protocol could
// be 'tcp'/'tls' and default to 'tls'. If client_protocol is 'http-tunnel', server_protocol must be 'tcp'.
const (
	UpdateEndpointsOptions_ServerProtocol_Tcp = "tcp"
	UpdateEndpointsOptions_ServerProtocol_Tls = "tls"
	UpdateEndpointsOptions_ServerProtocol_Udp = "udp"
)

// NewUpdateEndpointsOptions : Instantiate UpdateEndpointsOptions
func (*SatelliteLinkV1) NewUpdateEndpointsOptions(locationID string, endpointID string) *UpdateEndpointsOptions {
	return &UpdateEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *UpdateEndpointsOptions) SetLocationID(locationID string) *UpdateEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *UpdateEndpointsOptions) SetEndpointID(endpointID string) *UpdateEndpointsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetDisplayName : Allow user to set DisplayName
func (options *UpdateEndpointsOptions) SetDisplayName(displayName string) *UpdateEndpointsOptions {
	options.DisplayName = core.StringPtr(displayName)
	return options
}

// SetServerHost : Allow user to set ServerHost
func (options *UpdateEndpointsOptions) SetServerHost(serverHost string) *UpdateEndpointsOptions {
	options.ServerHost = core.StringPtr(serverHost)
	return options
}

// SetServerPort : Allow user to set ServerPort
func (options *UpdateEndpointsOptions) SetServerPort(serverPort int64) *UpdateEndpointsOptions {
	options.ServerPort = core.Int64Ptr(serverPort)
	return options
}

// SetSni : Allow user to set Sni
func (options *UpdateEndpointsOptions) SetSni(sni string) *UpdateEndpointsOptions {
	options.Sni = core.StringPtr(sni)
	return options
}

// SetClientProtocol : Allow user to set ClientProtocol
func (options *UpdateEndpointsOptions) SetClientProtocol(clientProtocol string) *UpdateEndpointsOptions {
	options.ClientProtocol = core.StringPtr(clientProtocol)
	return options
}

// SetClientMutualAuth : Allow user to set ClientMutualAuth
func (options *UpdateEndpointsOptions) SetClientMutualAuth(clientMutualAuth bool) *UpdateEndpointsOptions {
	options.ClientMutualAuth = core.BoolPtr(clientMutualAuth)
	return options
}

// SetServerProtocol : Allow user to set ServerProtocol
func (options *UpdateEndpointsOptions) SetServerProtocol(serverProtocol string) *UpdateEndpointsOptions {
	options.ServerProtocol = core.StringPtr(serverProtocol)
	return options
}

// SetServerMutualAuth : Allow user to set ServerMutualAuth
func (options *UpdateEndpointsOptions) SetServerMutualAuth(serverMutualAuth bool) *UpdateEndpointsOptions {
	options.ServerMutualAuth = core.BoolPtr(serverMutualAuth)
	return options
}

// SetRejectUnauth : Allow user to set RejectUnauth
func (options *UpdateEndpointsOptions) SetRejectUnauth(rejectUnauth bool) *UpdateEndpointsOptions {
	options.RejectUnauth = core.BoolPtr(rejectUnauth)
	return options
}

// SetTimeout : Allow user to set Timeout
func (options *UpdateEndpointsOptions) SetTimeout(timeout int64) *UpdateEndpointsOptions {
	options.Timeout = core.Int64Ptr(timeout)
	return options
}

// SetCreatedBy : Allow user to set CreatedBy
func (options *UpdateEndpointsOptions) SetCreatedBy(createdBy string) *UpdateEndpointsOptions {
	options.CreatedBy = core.StringPtr(createdBy)
	return options
}

// SetCerts : Allow user to set Certs
func (options *UpdateEndpointsOptions) SetCerts(certs *UpdatedEndpointRequestCerts) *UpdateEndpointsOptions {
	options.Certs = certs
	return options
}

// SetEnabled : Allow user to set Enabled
func (options *UpdateEndpointsOptions) SetEnabled(enabled bool) *UpdateEndpointsOptions {
	options.Enabled = core.BoolPtr(enabled)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateEndpointsOptions) SetHeaders(param map[string]string) *UpdateEndpointsOptions {
	options.Headers = param
	return options
}

// UpdateLinkOptions : The UpdateLink options.
type UpdateLinkOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The ws endpoint of the location.
	WsEndpoint *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateLinkOptions : Instantiate UpdateLinkOptions
func (*SatelliteLinkV1) NewUpdateLinkOptions(locationID string) *UpdateLinkOptions {
	return &UpdateLinkOptions{
		LocationID: core.StringPtr(locationID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *UpdateLinkOptions) SetLocationID(locationID string) *UpdateLinkOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetWsEndpoint : Allow user to set WsEndpoint
func (options *UpdateLinkOptions) SetWsEndpoint(wsEndpoint string) *UpdateLinkOptions {
	options.WsEndpoint = core.StringPtr(wsEndpoint)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateLinkOptions) SetHeaders(param map[string]string) *UpdateLinkOptions {
	options.Headers = param
	return options
}

// UpdateSourceEndpointsOptions : The UpdateSourceEndpoints options.
type UpdateSourceEndpointsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Source ID.
	SourceID *string `validate:"required,ne="`

	Endpoints []EndpointSourceStatusEndpointsItem

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSourceEndpointsOptions : Instantiate UpdateSourceEndpointsOptions
func (*SatelliteLinkV1) NewUpdateSourceEndpointsOptions(locationID string, sourceID string) *UpdateSourceEndpointsOptions {
	return &UpdateSourceEndpointsOptions{
		LocationID: core.StringPtr(locationID),
		SourceID:   core.StringPtr(sourceID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *UpdateSourceEndpointsOptions) SetLocationID(locationID string) *UpdateSourceEndpointsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetSourceID : Allow user to set SourceID
func (options *UpdateSourceEndpointsOptions) SetSourceID(sourceID string) *UpdateSourceEndpointsOptions {
	options.SourceID = core.StringPtr(sourceID)
	return options
}

// SetEndpoints : Allow user to set Endpoints
func (options *UpdateSourceEndpointsOptions) SetEndpoints(endpoints []EndpointSourceStatusEndpointsItem) *UpdateSourceEndpointsOptions {
	options.Endpoints = endpoints
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSourceEndpointsOptions) SetHeaders(param map[string]string) *UpdateSourceEndpointsOptions {
	options.Headers = param
	return options
}

// UpdateSourcesOptions : The UpdateSources options.
type UpdateSourcesOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Source ID.
	SourceID *string `validate:"required,ne="`

	// The name of the source, should be unique under each location. Source names must start with a letter and end with an
	// alphanumeric character, can contain letters, numbers, and hyphen (-), and must be 63 characters or fewer.
	SourceName *string

	Addresses []string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUpdateSourcesOptions : Instantiate UpdateSourcesOptions
func (*SatelliteLinkV1) NewUpdateSourcesOptions(locationID string, sourceID string) *UpdateSourcesOptions {
	return &UpdateSourcesOptions{
		LocationID: core.StringPtr(locationID),
		SourceID:   core.StringPtr(sourceID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *UpdateSourcesOptions) SetLocationID(locationID string) *UpdateSourcesOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetSourceID : Allow user to set SourceID
func (options *UpdateSourcesOptions) SetSourceID(sourceID string) *UpdateSourcesOptions {
	options.SourceID = core.StringPtr(sourceID)
	return options
}

// SetSourceName : Allow user to set SourceName
func (options *UpdateSourcesOptions) SetSourceName(sourceName string) *UpdateSourcesOptions {
	options.SourceName = core.StringPtr(sourceName)
	return options
}

// SetAddresses : Allow user to set Addresses
func (options *UpdateSourcesOptions) SetAddresses(addresses []string) *UpdateSourcesOptions {
	options.Addresses = addresses
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UpdateSourcesOptions) SetHeaders(param map[string]string) *UpdateSourcesOptions {
	options.Headers = param
	return options
}

// UpdatedEndpointRequestCerts : The certs.
type UpdatedEndpointRequestCerts struct {
	// The CA which Satellite Link trust when receiving the connection from the client application.
	Client *UpdatedEndpointRequestCertsClient `json:"client,omitempty"`

	// The CA which Satellite Link trust when sending the connection to server application.
	Server *UpdatedEndpointRequestCertsServer `json:"server,omitempty"`

	// The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.
	Connector *UpdatedEndpointRequestCertsConnector `json:"connector,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCerts unmarshals an instance of UpdatedEndpointRequestCerts from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCerts(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCerts)
	err = core.UnmarshalModel(m, "client", &obj.Client, UnmarshalUpdatedEndpointRequestCertsClient)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "server", &obj.Server, UnmarshalUpdatedEndpointRequestCertsServer)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "connector", &obj.Connector, UnmarshalUpdatedEndpointRequestCertsConnector)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsClient : The CA which Satellite Link trust when receiving the connection from the client application.
type UpdatedEndpointRequestCertsClient struct {
	// The root cert or the self-signed cert of the client application.
	Cert *UpdatedEndpointRequestCertsClientCert `json:"cert,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsClient unmarshals an instance of UpdatedEndpointRequestCertsClient from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsClient(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsClient)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalUpdatedEndpointRequestCertsClientCert)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsClientCert : The root cert or the self-signed cert of the client application.
type UpdatedEndpointRequestCertsClientCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`

	// The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsClientCert unmarshals an instance of UpdatedEndpointRequestCertsClientCert from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsClientCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsClientCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsConnector : The cert which Satellite Link connector provide to identify itself for connecting to the client/server application.
type UpdatedEndpointRequestCertsConnector struct {
	// The end-entity cert. This is required when the key is defined.
	Cert *UpdatedEndpointRequestCertsConnectorCert `json:"cert,omitempty"`

	// The private key of the end-entity certificate. This is required when the cert is defined.
	Key *UpdatedEndpointRequestCertsConnectorKey `json:"key,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsConnector unmarshals an instance of UpdatedEndpointRequestCertsConnector from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsConnector(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsConnector)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalUpdatedEndpointRequestCertsConnectorCert)
	if err != nil {
		return
	}
	err = core.UnmarshalModel(m, "key", &obj.Key, UnmarshalUpdatedEndpointRequestCertsConnectorKey)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsConnectorCert : The end-entity cert. This is required when the key is defined.
type UpdatedEndpointRequestCertsConnectorCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`

	// The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsConnectorCert unmarshals an instance of UpdatedEndpointRequestCertsConnectorCert from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsConnectorCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsConnectorCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsConnectorKey : The private key of the end-entity certificate. This is required when the cert is defined.
type UpdatedEndpointRequestCertsConnectorKey struct {
	// The name of the key.
	Filename *string `json:"filename,omitempty"`

	// The content of the key. The private key file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsConnectorKey unmarshals an instance of UpdatedEndpointRequestCertsConnectorKey from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsConnectorKey(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsConnectorKey)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsServer : The CA which Satellite Link trust when sending the connection to server application.
type UpdatedEndpointRequestCertsServer struct {
	// The root cert or the self-signed cert of the server application.
	Cert *UpdatedEndpointRequestCertsServerCert `json:"cert,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsServer unmarshals an instance of UpdatedEndpointRequestCertsServer from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsServer(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsServer)
	err = core.UnmarshalModel(m, "cert", &obj.Cert, UnmarshalUpdatedEndpointRequestCertsServerCert)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UpdatedEndpointRequestCertsServerCert : The root cert or the self-signed cert of the server application.
type UpdatedEndpointRequestCertsServerCert struct {
	// The filename of the cert.
	Filename *string `json:"filename,omitempty"`

	// The content of the cert. The certificate file must be in Privacy-enhanced Electronic Mail (PEM) format.
	FileContents *string `json:"file_contents,omitempty"`
}

// UnmarshalUpdatedEndpointRequestCertsServerCert unmarshals an instance of UpdatedEndpointRequestCertsServerCert from the specified map of raw messages.
func UnmarshalUpdatedEndpointRequestCertsServerCert(m map[string]json.RawMessage, result interface{}) (err error) {
	obj := new(UpdatedEndpointRequestCertsServerCert)
	err = core.UnmarshalPrimitive(m, "filename", &obj.Filename)
	if err != nil {
		return
	}
	err = core.UnmarshalPrimitive(m, "file_contents", &obj.FileContents)
	if err != nil {
		return
	}
	reflect.ValueOf(result).Elem().Set(reflect.ValueOf(obj))
	return
}

// UploadEndpointCertsOptions : The UploadEndpointCerts options.
type UploadEndpointCertsOptions struct {
	// The Location ID.
	LocationID *string `validate:"required,ne="`

	// The Endpoint ID.
	EndpointID *string `validate:"required,ne="`

	// The cert which Satellite Link trust when receiving the connection from the client application. Up to one cert could
	// be uploaded.
	ClientCert io.ReadCloser

	// The content type of clientCert.
	ClientCertContentType *string

	// The cert which Satellite Link trust when sending the connection to server application. Up to one cert could be
	// uploaded.
	ServerCert io.ReadCloser

	// The content type of serverCert.
	ServerCertContentType *string

	// The end-entity cert which Satellite Link connector provide to identify itself for connecting to the client/server
	// application. If uploading destKey as well, this field is required. Up to one cert could be uploaded.
	ConnectorCert io.ReadCloser

	// The content type of connectorCert.
	ConnectorCertContentType *string

	// The key for the connector_cert. If uploading connector_cert as well, this field is required. Up to one key could be
	// uploaded.
	ConnectorKey io.ReadCloser

	// The content type of connectorKey.
	ConnectorKeyContentType *string

	// Allows users to set headers on API requests
	Headers map[string]string
}

// NewUploadEndpointCertsOptions : Instantiate UploadEndpointCertsOptions
func (*SatelliteLinkV1) NewUploadEndpointCertsOptions(locationID string, endpointID string) *UploadEndpointCertsOptions {
	return &UploadEndpointCertsOptions{
		LocationID: core.StringPtr(locationID),
		EndpointID: core.StringPtr(endpointID),
	}
}

// SetLocationID : Allow user to set LocationID
func (options *UploadEndpointCertsOptions) SetLocationID(locationID string) *UploadEndpointCertsOptions {
	options.LocationID = core.StringPtr(locationID)
	return options
}

// SetEndpointID : Allow user to set EndpointID
func (options *UploadEndpointCertsOptions) SetEndpointID(endpointID string) *UploadEndpointCertsOptions {
	options.EndpointID = core.StringPtr(endpointID)
	return options
}

// SetClientCert : Allow user to set ClientCert
func (options *UploadEndpointCertsOptions) SetClientCert(clientCert io.ReadCloser) *UploadEndpointCertsOptions {
	options.ClientCert = clientCert
	return options
}

// SetClientCertContentType : Allow user to set ClientCertContentType
func (options *UploadEndpointCertsOptions) SetClientCertContentType(clientCertContentType string) *UploadEndpointCertsOptions {
	options.ClientCertContentType = core.StringPtr(clientCertContentType)
	return options
}

// SetServerCert : Allow user to set ServerCert
func (options *UploadEndpointCertsOptions) SetServerCert(serverCert io.ReadCloser) *UploadEndpointCertsOptions {
	options.ServerCert = serverCert
	return options
}

// SetServerCertContentType : Allow user to set ServerCertContentType
func (options *UploadEndpointCertsOptions) SetServerCertContentType(serverCertContentType string) *UploadEndpointCertsOptions {
	options.ServerCertContentType = core.StringPtr(serverCertContentType)
	return options
}

// SetConnectorCert : Allow user to set ConnectorCert
func (options *UploadEndpointCertsOptions) SetConnectorCert(connectorCert io.ReadCloser) *UploadEndpointCertsOptions {
	options.ConnectorCert = connectorCert
	return options
}

// SetConnectorCertContentType : Allow user to set ConnectorCertContentType
func (options *UploadEndpointCertsOptions) SetConnectorCertContentType(connectorCertContentType string) *UploadEndpointCertsOptions {
	options.ConnectorCertContentType = core.StringPtr(connectorCertContentType)
	return options
}

// SetConnectorKey : Allow user to set ConnectorKey
func (options *UploadEndpointCertsOptions) SetConnectorKey(connectorKey io.ReadCloser) *UploadEndpointCertsOptions {
	options.ConnectorKey = connectorKey
	return options
}

// SetConnectorKeyContentType : Allow user to set ConnectorKeyContentType
func (options *UploadEndpointCertsOptions) SetConnectorKeyContentType(connectorKeyContentType string) *UploadEndpointCertsOptions {
	options.ConnectorKeyContentType = core.StringPtr(connectorKeyContentType)
	return options
}

// SetHeaders : Allow user to set Headers
func (options *UploadEndpointCertsOptions) SetHeaders(param map[string]string) *UploadEndpointCertsOptions {
	options.Headers = param
	return options
}
